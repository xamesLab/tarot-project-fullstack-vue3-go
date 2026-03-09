package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/deckthemes"
)

type DeckThemesHandler struct {
	service *deckthemes.Service
}

type createDeckThemeRequest struct {
	Name string `json:"name"`
}

func NewDeckThemesHandler(service *deckthemes.Service) *DeckThemesHandler {
	return &DeckThemesHandler{service: service}
}

func (h *DeckThemesHandler) List(c *gin.Context) {
	themes, err := h.service.List(c.Request.Context())
	if err != nil {
		respondInternalError(c)
		return
	}

	c.JSON(http.StatusOK, themes)
}

func (h *DeckThemesHandler) Create(c *gin.Context) {
	var req createDeckThemeRequest
	if !bindJSON(c, &req) {
		return
	}

	created, err := h.service.Create(c.Request.Context(), currentUserID(c), req.Name)
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: deckthemes.ErrInvalidInput, status: http.StatusBadRequest, code: "invalid_input", message: "invalid theme name"},
			errorSpec{target: deckthemes.ErrThemeNameConflict, status: http.StatusConflict, code: "theme_exists", message: "theme name already exists"},
		)
		return
	}

	c.JSON(http.StatusCreated, created)
}
