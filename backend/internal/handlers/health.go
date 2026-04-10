package handlers

import (
	"net/http"

	"tarot-daily/config"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	config *config.Config
}

func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{config: cfg}
}

func (h *HealthHandler) Register(router gin.IRoutes) {
	router.GET("/health", h.health)
	router.GET("/ai/health", h.aiHealth)
	router.GET("/gateway/health", h.gatewayHealth)
}

func (h *HealthHandler) health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"app":    "tarot-daily-api",
		"env":    h.config.AppEnv,
	})
}

func (h *HealthHandler) gatewayHealth(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "stub",
		"message": "Gateway is not implemented",
	})
}

func (h *HealthHandler) aiHealth(c *gin.Context) {
	configured := h.config.OpenRouterKey != "" && h.config.OpenRouterModelMain != ""
	status := "disabled"
	message := "Direct AI integration is disabled"

	switch {
	case h.config.AIDirectEnabled && configured:
		status = "configured"
		message = "Direct AI integration is configured"
	case h.config.AIDirectEnabled:
		status = "misconfigured"
		message = "Direct AI integration is enabled but missing required configuration"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   status,
		"provider": "openrouter",
		"enabled":  h.config.AIDirectEnabled,
		"message":  message,
		"model":    h.config.OpenRouterModelMain,
	})
}
