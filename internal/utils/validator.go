package utils

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateStruct(s interface{}) error {
	// Register custom validations if needed
	// validate.RegisterValidation("custom", func(fl validator.FieldLevel) bool { ... })

	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	// Convert validation errors to user-friendly messages
	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		field := strings.ToLower(err.Field())
		switch err.Tag() {
		case "required":
			errors = append(errors, field+" is required")
		case "email":
			errors = append(errors, field+" must be a valid email")
		case "min":
			errors = append(errors, field+" must be at least "+err.Param()+" characters")
		default:
			errors = append(errors, field+" is invalid")
		}
	}

	return &ValidationError{Errors: errors}
}

type ValidationError struct {
	Errors []string
}

func (e *ValidationError) Error() string {
	return strings.Join(e.Errors, ", ")
}

func IsZero(value interface{}) bool {
	return reflect.ValueOf(value).IsZero()
}
