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

type SPKHandler struct {
	spks *service.SPKService
}

func NewSPKHandler(spks *service.SPKService) *SPKHandler {
	return &SPKHandler{spks: spks}
}

// Types returns the SPK catalog (Proses J-1..J-8) for the UI form.
func (h *SPKHandler) Types(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"types": model.SPKTypes})
}

func (h *SPKHandler) List(c *gin.Context) {
	items, err := h.spks.List(c.Query("status"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SPKHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	spk, err := h.spks.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SPK not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spk)
}

// Create issues a draft SPK. Restricted to KADEP at the route.
func (h *SPKHandler) Create(c *gin.Context) {
	var in service.CreateSPKInput
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	spk, err := h.spks.Create(in, middleware.CurrentUserID(c))
	if errors.Is(err, service.ErrInvalidSPKType) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "vendor tidak ditemukan"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, spk)
}

type spkDecisionInput struct {
	Note string `json:"note"`
}

// Approve / Reject — restricted to DIROPS at the route.
func (h *SPKHandler) Approve(c *gin.Context) { h.decide(c, true) }
func (h *SPKHandler) Reject(c *gin.Context)  { h.decide(c, false) }

func (h *SPKHandler) decide(c *gin.Context, approve bool) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var in spkDecisionInput
	_ = c.ShouldBindJSON(&in)
	spk, err := h.spks.Decide(id, approve, middleware.CurrentUserID(c), in.Note)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "SPK not found"})
		return
	}
	if errors.Is(err, service.ErrSPKNotDraft) {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, spk)
}
