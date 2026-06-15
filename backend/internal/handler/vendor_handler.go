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

type VendorHandler struct {
	vendors *service.VendorService
}

func NewVendorHandler(vendors *service.VendorService) *VendorHandler {
	return &VendorHandler{vendors: vendors}
}

func (h *VendorHandler) List(c *gin.Context) {
	items, err := h.vendors.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "categories": model.VendorCategories})
}

func (h *VendorHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	v, err := h.vendors.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "vendor not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}

func (h *VendorHandler) Create(c *gin.Context) {
	var in service.VendorInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := h.vendors.Create(in, middleware.CurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, v)
}

func (h *VendorHandler) Update(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in service.VendorInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	v, err := h.vendors.Update(id, in)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "vendor not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, v)
}
