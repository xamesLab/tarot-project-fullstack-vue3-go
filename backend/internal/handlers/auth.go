package handlers

import (
	"net/http"

	"tarot-daily/config"
	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	authService *services.AuthService
	validate    *validator.Validate
	config      *config.Config
}

func NewAuthHandler(authService *services.AuthService, validate *validator.Validate, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validate:    validate,
		config:      cfg,
	}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.register)
	router.POST("/login", h.login)
}

func (h *AuthHandler) register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, token, err := h.authService.Register(c.Request.Context(), req, h.config.JWTSecret, h.config.JWTTTL)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"user":    user.ToResponse(false),
		"token":   token,
	})
}

func (h *AuthHandler) login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	user, token, err := h.authService.Login(c.Request.Context(), req, h.config.JWTSecret, h.config.JWTTTL)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.ToResponse(false),
		"token":   token,
	})
}
