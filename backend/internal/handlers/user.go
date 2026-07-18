package handlers

import (
	"net/http"

	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
	userID, ok := requireUserID(c)
	if !ok {
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
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	req, ok := bindAndValidateJSON[models.UpdateProfileRequest](c, h.validate)
	if !ok {
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
	userID, ok := requireUserID(c)
	if !ok {
		return
	}

	req, ok := bindAndValidateJSON[models.UpdateSettingsRequest](c, h.validate)
	if !ok {
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
