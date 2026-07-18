package handlers

import (
	"errors"
	"net/http"
	"strings"

	"tarot-daily/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func bindAndValidateJSON[T any](c *gin.Context, validate *validator.Validate) (T, bool) {
	return bindJSON[T](c, validate, false)
}

func bindOptionalAndValidateJSON[T any](c *gin.Context, validate *validator.Validate) (T, bool) {
	return bindJSON[T](c, validate, true)
}

func bindJSON[T any](c *gin.Context, validate *validator.Validate, allowBodyNotAllowed bool) (T, bool) {
	var req T
	if err := c.ShouldBindJSON(&req); err != nil && !(allowBodyNotAllowed && errors.Is(err, http.ErrBodyNotAllowed)) {
		respondError(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Invalid request body")
		return req, false
	}
	if err := validate.Struct(req); err != nil {
		respondValidationError(c, err)
		return req, false
	}
	return req, true
}

func requireUserID(c *gin.Context) (primitive.ObjectID, bool) {
	id, err := primitive.ObjectIDFromHex(c.GetString("user_id"))
	if err != nil {
		respondError(c, http.StatusUnauthorized, services.ErrInvalidToken.Error(), "Authentication token is invalid")
		return primitive.NilObjectID, false
	}
	return id, true
}

func respondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   code,
		"message": message,
	})
}

func respondErrorWithDetails(c *gin.Context, status int, code, message string, details gin.H) {
	c.JSON(status, gin.H{
		"success": false,
		"error":   code,
		"message": message,
		"details": details,
	})
}

func respondValidationError(c *gin.Context, err error) {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		for _, item := range validationErrs {
			field := strings.ToLower(item.Field())
			switch item.Tag() {
			case "required":
				respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "This field is required", gin.H{"field": field, "reason": "required"})
				return
			case "email":
				respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Please enter a valid email address", gin.H{"field": field, "reason": "invalid"})
				return
			case "min", "max":
				if field == "username" {
					respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Username must contain between 3 and 20 characters", gin.H{"field": field, "reason": "invalid_format"})
					return
				}
				if field == "password" {
					respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Password must be at least 8 characters with letter and number", gin.H{"field": field, "reason": "invalid_format"})
					return
				}
			case "password":
				respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Password must be at least 8 characters with letter and number", gin.H{"field": field, "reason": "invalid_format"})
				return
			case "timezone":
				respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Please select a valid timezone", gin.H{"field": field, "reason": "invalid"})
				return
			case "username":
				respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), "Username may contain only letters, numbers, dots, underscores, and hyphens", gin.H{"field": field, "reason": "invalid_format"})
				return
			}
		}
	}

	var fieldError *services.ValidationError
	if errors.As(err, &fieldError) {
		respondErrorWithDetails(c, http.StatusBadRequest, services.ErrValidationFailed.Error(), fieldError.Message, gin.H{"field": fieldError.Field, "reason": fieldError.Reason})
		return
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
	case errors.Is(err, services.ErrRegistrationDisabled):
		respondError(c, http.StatusForbidden, services.ErrRegistrationDisabled.Error(), "Registration is currently disabled")
	case errors.Is(err, services.ErrInvalidCredentials):
		respondError(c, http.StatusUnauthorized, services.ErrInvalidCredentials.Error(), "Email or password incorrect")
	case errors.Is(err, services.ErrAccountDisabled):
		respondError(c, http.StatusForbidden, services.ErrAccountDisabled.Error(), "Your account has been suspended")
	case errors.Is(err, services.ErrForbidden):
		respondError(c, http.StatusForbidden, services.ErrForbidden.Error(), "Forbidden")
	case errors.Is(err, services.ErrInvalidGuestClaim):
		respondError(c, http.StatusBadRequest, services.ErrInvalidGuestClaim.Error(), "Guest card claim is invalid or expired")
	case errors.Is(err, services.ErrValidationFailed):
		respondValidationError(c, err)
	case errors.Is(err, mongo.ErrNoDocuments):
		respondError(c, http.StatusNotFound, "not_found", "Requested resource was not found")
	default:
		respondError(c, http.StatusInternalServerError, "internal_error", "Internal server error")
	}
}
