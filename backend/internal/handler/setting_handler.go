package handler

import (
	"net/http"

	"legalpermit/internal/model"
	"legalpermit/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	settings *service.SettingService
}

func NewSettingHandler(settings *service.SettingService) *SettingHandler {
	return &SettingHandler{settings: settings}
}

func (h *SettingHandler) GetDACI(c *gin.Context) {
	cfg, err := h.settings.GetDACI()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *SettingHandler) SetDACI(c *gin.Context) {
	var cfg model.DACIConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.settings.SetDACI(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *SettingHandler) GetNotification(c *gin.Context) {
	cfg, err := h.settings.GetNotification()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}

func (h *SettingHandler) SetNotification(c *gin.Context) {
	var cfg model.NotificationConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.settings.SetNotification(cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cfg)
}
