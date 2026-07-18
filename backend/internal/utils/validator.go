package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	letterPattern   = regexp.MustCompile(`[a-zA-Z]`)
	numberPattern   = regexp.MustCompile(`[0-9]`)
	usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
)

func NewValidator() *validator.Validate {
	validate := validator.New()
	_ = validate.RegisterValidation("password", passwordValidator)
	_ = validate.RegisterValidation("timezone", timezoneValidator)
	_ = validate.RegisterValidation("username", usernameValidator)
	return validate
}

func passwordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	hasLetter := letterPattern.MatchString(password)
	hasNumber := numberPattern.MatchString(password)
	return hasLetter && hasNumber
}

func timezoneValidator(fl validator.FieldLevel) bool {
	return ValidateTimezone(fl.Field().String())
}

func usernameValidator(fl validator.FieldLevel) bool {
	return usernamePattern.MatchString(fl.Field().String())
}
