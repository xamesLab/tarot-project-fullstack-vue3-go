package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/gatewayclient"
	"tarotdaily/backend/internal/transport/http/response"
)

type gatewayHealthChecker interface {
	Health(ctx context.Context) (*gatewayclient.HealthResponse, error)
}

type GatewayHealthHandler struct {
	client gatewayHealthChecker
}

func NewGatewayHealthHandler(client gatewayHealthChecker) *GatewayHealthHandler {
	return &GatewayHealthHandler{client: client}
}

func (h *GatewayHealthHandler) Get(c *gin.Context) {
	if h == nil || h.client == nil {
		response.JSONError(c, http.StatusInternalServerError, "invalid_config", "gateway integration is not configured")
		return
	}

	out, err := h.client.Health(c.Request.Context())
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *GatewayHealthHandler) mapError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, context.DeadlineExceeded):
		response.JSONError(c, http.StatusGatewayTimeout, "gateway_timeout", "gateway health request timed out")
	case errors.Is(err, gatewayclient.ErrInvalidInput):
		response.JSONError(c, http.StatusInternalServerError, "invalid_config", "gateway integration is not configured")
	case errors.Is(err, gatewayclient.ErrInvalidResponse):
		response.JSONError(c, http.StatusBadGateway, "gateway_bad_response", "gateway returned invalid response")
	case errors.Is(err, gatewayclient.ErrUnexpectedStatus):
		response.JSONError(c, http.StatusBadGateway, "gateway_unavailable", "gateway is unavailable")
	default:
		response.JSONError(c, http.StatusBadGateway, "gateway_unavailable", "gateway is unavailable")
	}
}
