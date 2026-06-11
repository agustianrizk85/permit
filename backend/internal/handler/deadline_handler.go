package handler

import (
	"net/http"

	"legalpermit/internal/model"
	"legalpermit/internal/service"

	"github.com/gin-gonic/gin"
)

type DeadlineHandler struct {
	deadlines *service.DeadlineService
}

func NewDeadlineHandler(deadlines *service.DeadlineService) *DeadlineHandler {
	return &DeadlineHandler{deadlines: deadlines}
}

// List returns the Master Deadline rules for every process step.
func (h *DeadlineHandler) List(c *gin.Context) {
	rules, err := h.deadlines.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rules})
}

// Update saves edited Master Deadline rules (KADEP/DIROPS only).
func (h *DeadlineHandler) Update(c *gin.Context) {
	var req struct {
		Items []model.DeadlineRule `json:"items"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	rules, err := h.deadlines.Update(req.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": rules})
}
