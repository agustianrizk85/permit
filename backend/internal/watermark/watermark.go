// Package watermark renders a "CONFIDENTIAL" watermark and converts images to
// black-and-white, per the spec: IMB/PBG, SHM and PBB shared with Sales must be
// watermarked Confidential and edited to black/white.
package watermark

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// IsImage reports whether the mime type is a supported raster image.
func IsImage(mime string) bool {
	mime = strings.ToLower(mime)
	return strings.HasPrefix(mime, "image/jpeg") ||
		strings.HasPrefix(mime, "image/jpg") ||
		strings.HasPrefix(mime, "image/png")
}

// Apply converts the image to grayscale and tiles a diagonal "CONFIDENTIAL"
// watermark across it. Returns JPEG bytes. Non-image inputs return an error so
// callers can fall back to serving the original file.
func Apply(data []byte, mime, label string) ([]byte, error) {
	src, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	bounds := src.Bounds()

	// 1) Grayscale (black & white) base.
	gray := image.NewGray(bounds)
	draw.Draw(gray, bounds, src, bounds.Min, draw.Src)

	// 2) Promote to RGBA so we can draw colored watermark text.
	out := image.NewRGBA(bounds)
	draw.Draw(out, bounds, gray, bounds.Min, draw.Src)

	if label == "" {
		label = "CONFIDENTIAL"
	}
	watermarkColor := color.RGBA{R: 200, G: 0, B: 0, A: 90}
	face := basicfont.Face7x13

	// Tile the label across the image on a diagonal grid.
	stepX, stepY := 220, 120
	for y := bounds.Min.Y; y < bounds.Max.Y; y += stepY {
		offset := 0
		if (y/stepY)%2 == 1 {
			offset = stepX / 2
		}
		for x := bounds.Min.X - stepX; x < bounds.Max.X; x += stepX {
			drawer := &font.Drawer{
				Dst:  out,
				Src:  image.NewUniform(watermarkColor),
				Face: face,
				Dot:  fixed.P(x+offset, y),
			}
			drawer.DrawString(label)
		}
	}

	var buf bytes.Buffer
	if strings.Contains(strings.ToLower(mime), "png") {
		if err := png.Encode(&buf, out); err != nil {
			return nil, err
		}
	} else {
		if err := jpeg.Encode(&buf, out, &jpeg.Options{Quality: 85}); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}
