package services

import "errors"

var (
	ErrValidationFailed   = errors.New("validation_failed")
	ErrInvalidCredentials = errors.New("invalid_credentials")
	ErrMissingToken       = errors.New("missing_token")
	ErrInvalidToken       = errors.New("invalid_token")
	ErrForbidden          = errors.New("forbidden")
	ErrAccountDisabled    = errors.New("account_disabled")
)
