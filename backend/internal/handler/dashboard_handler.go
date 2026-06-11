package handler

import (
	"net/http"
	"strconv"

	"legalpermit/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboard *service.DashboardService
	docs      *service.DocumentService
}

func NewDashboardHandler(dashboard *service.DashboardService, docs *service.DocumentService) *DashboardHandler {
	return &DashboardHandler{dashboard: dashboard, docs: docs}
}

// EarlyWarnings returns the AI-style early-warning feed across all projects.
func (h *DashboardHandler) EarlyWarnings(c *gin.Context) {
	warnings, err := h.dashboard.EarlyWarnings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"warnings": warnings, "count": len(warnings)})
}

// SearchDocuments powers the "Search All Dokumen" feature.
func (h *DashboardHandler) SearchDocuments(c *gin.Context) {
	query := c.Query("q")
	var projectID *uint
	if raw := c.Query("project_id"); raw != "" {
		if n, err := strconv.ParseUint(raw, 10, 64); err == nil {
			id := uint(n)
			projectID = &id
		}
	}
	docs, err := h.docs.Search(query, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": docs, "count": len(docs)})
}
