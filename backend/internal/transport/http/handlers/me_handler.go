package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"tarotdaily/backend/internal/auth"
)

type MeHandler struct {
	service *auth.Service
}

func NewMeHandler(service *auth.Service) *MeHandler {
	return &MeHandler{service: service}
}

func (h *MeHandler) GetMe(c *gin.Context) {
	me, err := h.service.Me(c.Request.Context(), currentUserID(c))
	if err != nil {
		respondByError(
			c,
			err,
			errorSpec{target: auth.ErrInvalidInput, status: http.StatusUnauthorized, code: "unauthorized", message: "invalid token subject"},
			errorSpec{target: auth.ErrUserNotFound, status: http.StatusUnauthorized, code: "unauthorized", message: "user not found"},
		)
		return
	}

	c.JSON(http.StatusOK, me)
}
