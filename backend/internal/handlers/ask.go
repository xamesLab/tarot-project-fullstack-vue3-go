package handlers

import (
	"errors"
	"net/http"

	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AskHandler struct {
	askService *services.AskService
	validate   *validator.Validate
}

func NewAskHandler(askService *services.AskService, validate *validator.Validate) *AskHandler {
	return &AskHandler{
		askService: askService,
		validate:   validate,
	}
}

func (h *AskHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/ask", h.createOrReuseAsk)
	router.GET("/ask/:requestId", h.getAsk)
	router.POST("/reading/card-of-day/interpretation", h.createUserCardOfDayInterpretation)
}

func (h *AskHandler) RegisterPublicRoutes(router *gin.RouterGroup) {
	router.POST("/reading/guest/card-of-day/interpretation", h.createGuestCardOfDayInterpretation)
}

func (h *AskHandler) createGuestCardOfDayInterpretation(c *gin.Context) {
	var req models.GuestCardOfDayInterpretationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	result, err := h.askService.GenerateGuestCardOfDayInterpretation(c.Request.Context(), req)
	if err != nil {
		h.respondAskError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"source":         result.Source,
		"interpretation": result.Interpretation,
	})
}

func (h *AskHandler) createUserCardOfDayInterpretation(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.UserCardOfDayInterpretationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}
	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	result, err := h.askService.GenerateUserCardOfDayInterpretation(c.Request.Context(), userID, req)
	if err != nil {
		h.respondAskError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"source":         result.Source,
		"interpretation": result.Interpretation,
	})
}

func (h *AskHandler) createOrReuseAsk(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.CreateAskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	request, err := h.askService.CreateOrReuseInterpretation(c.Request.Context(), userID, req)
	if err != nil {
		h.respondAskError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"request": request,
	})
}

func (h *AskHandler) getAsk(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	request, err := h.askService.GetRequest(c.Request.Context(), userID, c.Param("requestId"))
	if err != nil {
		h.respondAskError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"request": request,
	})
}

func (h *AskHandler) respondAskError(c *gin.Context, err error) {
	var providerErr *services.AIProviderError
	var aiLimitErr *services.AIDailyLimitReachedError
	switch {
	case errors.Is(err, services.ErrAIUnavailable):
		respondError(c, http.StatusServiceUnavailable, services.ErrAIUnavailable.Error(), "AI interpretation is currently unavailable")
	case errors.As(err, &aiLimitErr):
		c.JSON(http.StatusTooManyRequests, gin.H{
			"success":       false,
			"error":         services.ErrAIDailyLimitReached.Error(),
			"message":       "You've reached your daily AI interpretation limit",
			"next_reset_at": aiLimitErr.NextResetAt,
			"current_limit": gin.H{
				"count": aiLimitErr.Count,
				"limit": aiLimitErr.Limit,
			},
		})
	case errors.As(err, &providerErr):
		respondError(c, providerErr.StatusCode, providerErr.Code, providerErr.Message)
	default:
		respondServiceError(c, err)
	}
}
