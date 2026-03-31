package handlers

import (
	"errors"
	"net/http"
	"strings"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   code,
		"message": message,
	})
}

func respondValidationError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, item := range validationErrs {
			switch item.Tag() {
			case "password":
				respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Password must be at least 8 characters with letter and number")
				return
			case "timezone":
				respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Please select a valid timezone")
				return
			case "alphanum":
				respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Username must be alphanumeric")
				return
			}
		}
	}

	message := err.Error()
	if errors.Is(err, services.ErrValidationFailed) {
		parts := strings.SplitN(message, "\n", 2)
		if len(parts) == 2 && strings.TrimSpace(parts[1]) != "" {
			message = strings.TrimSpace(parts[1])
		}
	}
	respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), message)
}

func respondServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidCredentials):
		respondError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error(), "Email or password incorrect")
	case errors.Is(err, services.ErrAccountDisabled):
		respondError(c, http.StatusForbidden, services.ErrAccountDisabled.Error(), "Your account has been suspended")
	case errors.Is(err, services.ErrValidationFailed):
		respondValidationError(c, err)
	case errors.Is(err, mongo.ErrNoDocuments):
		respondError(c, http.StatusNotFound, "not_found", "Requested resource was not found")
	default:
		respondError(c, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}
