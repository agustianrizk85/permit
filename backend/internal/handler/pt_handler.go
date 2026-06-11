package handler

import (
	"errors"
	"net/http"

	"legalpermit/internal/middleware"
	"legalpermit/internal/model"
	"legalpermit/internal/repository"
	"legalpermit/internal/service"

	"github.com/gin-gonic/gin"
)

type PTHandler struct {
	pts *service.PTService
}

func NewPTHandler(pts *service.PTService) *PTHandler {
	return &PTHandler{pts: pts}
}

func (h *PTHandler) Create(c *gin.Context) {
	var in service.CreatePTInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	pt, err := h.pts.Create(in, middleware.CurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, pt)
}

func (h *PTHandler) List(c *gin.Context) {
	pts, err := h.pts.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Expose the expected document slots so the UI can render the upload form.
	c.JSON(http.StatusOK, gin.H{"items": pts, "doc_types": model.PTDocTypes})
}

func (h *PTHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	pt, err := h.pts.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "PT not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"pt": pt, "doc_types": model.PTDocTypes})
}

func (h *PTHandler) UploadDocument(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}
	docType := c.PostForm("doc_type")
	if docType == "" {
		docType = "Lainnya"
	}
	doc, err := h.pts.UploadDocument(id, docType, fh, middleware.CurrentUserID(c), c.SaveUploadedFile)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "PT not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, doc)
}

func (h *PTHandler) DownloadDocument(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	doc, err := h.pts.GetDocument(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "document not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.FileAttachment(doc.Path, doc.OriginalName)
}
