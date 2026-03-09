package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/reading"
)

type ReadingHandler struct {
	service *reading.Service
}

func NewReadingHandler(service *reading.Service) *ReadingHandler {
	return &ReadingHandler{service: service}
}

func (h *ReadingHandler) GetToday(c *gin.Context) {
	result, err := h.service.GetToday(c.Request.Context(), currentUserID(c))
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReadingHandler) GenerateSpreads(c *gin.Context) {
	result, err := h.service.GenerateSpreads(c.Request.Context(), currentUserID(c))
	if err != nil {
		h.mapError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ReadingHandler) mapError(c *gin.Context, err error) {
	respondByError(
		c,
		err,
		errorSpec{target: reading.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid reading input"},
		errorSpec{target: reading.ErrUserNotFound, status: http.StatusUnauthorized, code: "unauthorized", message: "user not found"},
		errorSpec{target: reading.ErrRolePolicyNotFound, status: http.StatusForbidden, code: "forbidden", message: "role policy not found"},
		errorSpec{target: reading.ErrDailySpreadLimitReached, status: http.StatusTooManyRequests, code: "daily_spread_limit_reached", message: "daily spread limit reached"},
		errorSpec{target: reading.ErrCardsCatalogNotAvailable, status: http.StatusServiceUnavailable, code: "cards_unavailable", message: "cards catalog is not available"},
	)
}
