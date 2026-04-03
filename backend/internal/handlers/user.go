package handlers

import (
	"net/http"

	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserHandler struct {
	userService *services.UserService
	validate    *validator.Validate
}

func NewUserHandler(userService *services.UserService, validate *validator.Validate) *UserHandler {
	return &UserHandler{userService: userService, validate: validate}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/me", h.getMe)
	router.PATCH("/me/profile", h.updateProfile)
	router.PATCH("/me/settings", h.updateSettings)
}

func (h *UserHandler) getMe(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	user, err := h.userService.GetUserProfile(c.Request.Context(), userID)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(true),
	})
}

func (h *UserHandler) updateProfile(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateProfile(c.Request.Context(), userID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(true),
	})
}

func (h *UserHandler) updateSettings(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, err := h.userService.UpdateSettings(c.Request.Context(), userID, req)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(true),
	})
}

func parseUserID(c *gin.Context) (primitive.ObjectID, bool) {
	raw := c.GetString("user_id")
	id, err := primitive.ObjectIDFromHex(raw)
	if err != nil {
		return primitive.NilObjectID, false
	}
	return id, true
}
