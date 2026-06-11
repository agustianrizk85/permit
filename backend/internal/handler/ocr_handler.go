package handler

import (
	"io"
	"net/http"

	"legalpermit/internal/ocr"

	"github.com/gin-gonic/gin"
)

type OCRHandler struct {
	provider ocr.Provider
}

func NewOCRHandler(provider ocr.Provider) *OCRHandler {
	return &OCRHandler{provider: provider}
}

// Extract runs OCR/AI extraction on an uploaded document and returns structured
// fields the UI can drop into step metadata. Does not persist the file.
func (h *OCRHandler) Extract(c *gin.Context) {
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	docType := c.PostForm("doc_type")
	if docType == "" {
		docType = "KTP"
	}

	f, err := fh.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer f.Close()
	data, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	result, err := h.provider.Extract(c.Request.Context(), docType, fh.Filename, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}
