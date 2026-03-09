package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/decks"
)

type DecksHandler struct {
	service *decks.Service
}

func NewDecksHandler(service *decks.Service) *DecksHandler {
	return &DecksHandler{service: service}
}

func (h *DecksHandler) ListAvailable(c *gin.Context) {
	items, err := h.service.ListAvailable(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, items)
}
