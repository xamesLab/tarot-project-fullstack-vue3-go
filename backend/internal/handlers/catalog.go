package handlers

import (
	"net/http"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
)

type CatalogHandler struct {
	catalogService *services.CatalogService
}

func NewCatalogHandler(catalogService *services.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

func (h *CatalogHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/decks/available", h.listAvailableDecks)
	router.GET("/card-backs/available", h.listAvailableCardBacks)
}

func (h *CatalogHandler) listAvailableDecks(c *gin.Context) {
	decks, defaultDeckID, err := h.catalogService.ListAvailableDecks(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"decks":           decks,
		"default_deck_id": defaultDeckID,
	})
}

func (h *CatalogHandler) listAvailableCardBacks(c *gin.Context) {
	cardBacks, defaultBackID, err := h.catalogService.ListAvailableCardBacks(c.Request.Context())
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"card_backs":      cardBacks,
		"default_back_id": defaultBackID,
	})
}
