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
