package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/auth"
)

type SettingsHandler struct {
	service *auth.Service
}

type patchSettingsRequest struct {
	DisplayName     *string `json:"displayName"`
	Timezone        *string `json:"timezone"`
	PreferredDeckID *string `json:"preferredDeckId"`
	DeckThemeID     *string `json:"deckThemeId"`
}

func NewSettingsHandler(service *auth.Service) *SettingsHandler {
	return &SettingsHandler{service: service}
}

func (h *SettingsHandler) PatchMeSettings(c *gin.Context) {
	var req patchSettingsRequest
	if !bindJSON(c, &req) {
		return
	}

	updated, err := h.service.UpdateSettings(c.Request.Context(), currentUserID(c), auth.UpdateSettingsInput{
		DisplayName:     req.DisplayName,
		Timezone:        req.Timezone,
		PreferredDeckID: req.PreferredDeckID,
		DeckThemeID:     req.DeckThemeID,
	})
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: auth.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid displayName, timezone, preferredDeckId or deckThemeId"},
			errorSpec{target: auth.ErrPreferredDeckNotUsable, status: http.StatusBadRequest, code: "invalid_preferred_deck", message: "preferred deck must exist and be done"},
			errorSpec{target: auth.ErrDeckThemeNotFound, status: http.StatusBadRequest, code: "invalid_deck_theme", message: "deck theme not found"},
			errorSpec{target: auth.ErrUserNotFound, status: http.StatusUnauthorized, code: "unauthorized", message: "user not found"},
		)
		return
	}

	c.JSON(http.StatusOK, updated)
}
