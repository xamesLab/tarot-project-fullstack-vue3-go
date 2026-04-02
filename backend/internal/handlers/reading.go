package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"tarot-daily/internal/models"
	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ReadingHandler struct {
	readingService *services.ReadingService
	validate       *validator.Validate
}

func NewReadingHandler(readingService *services.ReadingService, validate *validator.Validate) *ReadingHandler {
	return &ReadingHandler{
		readingService: readingService,
		validate:       validate,
	}
}

func (h *ReadingHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/reading/generate/card-of-day", h.generateCardOfDay)
	router.POST("/reading/generate/three-card", h.generateThreeCardSpread)
	router.POST("/reading/generate/spread-a", h.generateThreeCardSpread)
	router.GET("/reading/history", h.getHistory)
}

func (h *ReadingHandler) generateCardOfDay(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.GenerateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	reading, err := h.readingService.GenerateCardOfDay(c.Request.Context(), userID, req, c.ClientIP())
	if err != nil {
		var alreadyExists *services.CardOfDayAlreadyGeneratedError
		if errors.As(err, &alreadyExists) {
			c.JSON(http.StatusBadRequest, gin.H{
				"success":           false,
				"error":             services.ErrCardOfDayExists.Error(),
				"message":           "You already received your Card of the Day",
				"existing_reading":  alreadyExists.Existing,
				"next_available_at": alreadyExists.NextAvailableAt,
			})
			return
		}

		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"reading": reading,
	})
}

func (h *ReadingHandler) generateThreeCardSpread(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	var req models.GenerateReadingRequest
	if err := c.ShouldBindJSON(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return
	}

	if err := h.validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return
	}

	reading, err := h.readingService.GenerateThreeCardSpread(c.Request.Context(), userID, req, c.ClientIP())
	if err != nil {
		var limitReached *services.DailyLimitReachedError
		if errors.As(err, &limitReached) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success":       false,
				"error":         services.ErrDailyLimitReached.Error(),
				"message":       "You've reached your daily spread limit",
				"next_reset_at": limitReached.NextResetAt,
				"current_limit": models.RemainingToday{
					Count: limitReached.Count,
					Limit: limitReached.Limit,
				},
			})
			return
		}

		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"reading": reading,
	})
}

func (h *ReadingHandler) getHistory(c *gin.Context) {
	userID, ok := parseUserID(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return
	}

	page, err := parseIntQuery(c.Query("page"), 1)
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Page must be a positive integer")
		return
	}

	limit, err := parseIntQuery(c.Query("limit"), 10)
	if err != nil {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Limit must be a positive integer")
		return
	}

	result, err := h.readingService.GetReadingHistory(c.Request.Context(), userID, c.DefaultQuery("type", "all"), page, limit)
	if err != nil {
		respondServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":        true,
		"readings":       result.Readings,
		"pagination":     result.Pagination,
		"retention_info": result.RetentionInfo,
	})
}

func parseIntQuery(raw string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, errors.New("invalid integer")
	}

	return value, nil
}
