package handler

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"legalpermit/internal/dto"
	"legalpermit/internal/middleware"
	"legalpermit/internal/repository"
	"legalpermit/internal/service"
	"legalpermit/internal/watermark"

	"github.com/gin-gonic/gin"
)

type StepHandler struct {
	steps *service.StepService
	docs  *service.DocumentService
}

func NewStepHandler(steps *service.StepService, docs *service.DocumentService) *StepHandler {
	return &StepHandler{steps: steps, docs: docs}
}

func (h *StepHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	step, err := h.steps.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, step)
}

func (h *StepHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var req dto.UpdateStepRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	step, err := h.steps.Update(id, req, middleware.CurrentUserID(c))
	switch {
	case errors.Is(err, repository.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
	case errors.Is(err, service.ErrPriceRequired), errors.Is(err, service.ErrSPKRequired):
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusOK, step)
	}
}

// UploadDocument handles multipart upload of a file attached to a step.
func (h *StepHandler) UploadDocument(c *gin.Context) {
	stepID, ok := parseID(c)
	if !ok {
		return
	}
	step, err := h.steps.Get(stepID)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "step not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	docType := c.PostForm("doc_type")
	if docType == "" {
		docType = step.Code
	}
	confidential, _ := strconv.ParseBool(c.PostForm("confidential"))

	doc, err := h.docs.Upload(fh, service.UploadParams{
		ProjectID:     step.ProjectID,
		ProcessStepID: &step.ID,
		DocType:       docType,
		Confidential:  confidential,
		UploadedBy:    middleware.CurrentUserID(c),
	}, c.SaveUploadedFile)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

// DownloadDocument streams a stored file back to the client. With ?watermark=1
// (and for image files) it returns a black-and-white, "CONFIDENTIAL"-watermarked
// rendition for Sales / cross-department sharing per the spec.
func (h *StepHandler) DownloadDocument(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	doc, err := h.docs.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	wantWatermark, _ := strconv.ParseBool(c.Query("watermark"))
	if wantWatermark && watermark.IsImage(doc.MimeType) {
		raw, err := os.ReadFile(doc.Path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		out, err := watermark.Apply(raw, doc.MimeType, "CONFIDENTIAL")
		if err != nil {
			// Fall back to the original file if watermarking fails.
			c.FileAttachment(doc.Path, doc.OriginalName)
			return
		}
		c.Header("Content-Disposition", "attachment; filename=confidential_"+doc.OriginalName)
		c.Data(http.StatusOK, doc.MimeType, out)
		return
	}

	c.FileAttachment(doc.Path, doc.OriginalName)
}
