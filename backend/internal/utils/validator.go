package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
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
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	return hasLetter && hasNumber
}

func timezoneValidator(fl validator.FieldLevel) bool {
	return ValidateTimezone(fl.Field().String())
}

func usernameValidator(fl validator.FieldLevel) bool {
	return regexp.MustCompile(`^[a-zA-Z0-9._-]+$`).MatchString(fl.Field().String())
}
