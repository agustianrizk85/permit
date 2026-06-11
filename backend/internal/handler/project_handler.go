package handler

import (
	"errors"
	"net/http"
	"strconv"

	"legalpermit/internal/dto"
	"legalpermit/internal/middleware"
	"legalpermit/internal/repository"
	"legalpermit/internal/service"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	projects *service.ProjectService
}

func NewProjectHandler(projects *service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projects: projects}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var req dto.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := h.projects.Create(req, middleware.CurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.projects.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, projects)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	project, err := h.projects.Get(id)
	if errors.Is(err, repository.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "project not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Progress(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	progress, err := h.projects.Progress(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, progress)
}

// parseID extracts and validates the :id route parameter.
func parseID(c *gin.Context) (uint, bool) {
	raw := c.Param("id")
	n, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(n), true
}
