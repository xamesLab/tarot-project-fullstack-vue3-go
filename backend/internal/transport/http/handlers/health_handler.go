package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	PingDB      func(ctx context.Context) error
	PingTimeout time.Duration
}

func NewHealthHandler(pingDB func(ctx context.Context) error, pingTimeout time.Duration) *HealthHandler {
	return &HealthHandler{
		PingDB:      pingDB,
		PingTimeout: pingTimeout,
	}
}

func (h *HealthHandler) Live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"type":   "liveness",
	})
}

func (h *HealthHandler) Ready(c *gin.Context) {
	if h.PingDB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"type":     "readiness",
			"database": "unavailable",
		})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), h.PingTimeout)
	defer cancel()

	if err := h.PingDB(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":   "degraded",
			"type":     "readiness",
			"database": "disconnected",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   "ok",
		"type":     "readiness",
		"database": "connected",
	})
}
