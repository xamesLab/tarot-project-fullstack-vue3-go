package services

import "errors"

var (
	ErrValidationFailed     = errors.New("validation_failed")
	ErrInvalidCredentials   = errors.New("invalid_credentials")
	ErrRegistrationDisabled = errors.New("registration_disabled")
	ErrMissingToken         = errors.New("missing_token")
	ErrInvalidToken         = errors.New("invalid_token")
	ErrForbidden            = errors.New("forbidden")
	ErrAccountDisabled      = errors.New("account_disabled")
	ErrCardOfDayExists      = errors.New("card_of_day_already_generated")
	ErrInvalidGuestClaim    = errors.New("invalid_guest_claim")
	ErrDailyLimitReached    = errors.New("daily_limit_reached")
	ErrAIUnavailable        = errors.New("ai_unavailable")
	ErrAIDailyLimitReached  = errors.New("ai_daily_limit_reached")
)

type ValidationError struct {
	Field   string
	Reason  string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func (e *ValidationError) Unwrap() error {
	return ErrValidationFailed
}

func newValidationError(field, reason, message string) error {
	return &ValidationError{Field: field, Reason: reason, Message: message}
}
