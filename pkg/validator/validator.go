package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type CustomValidator struct {
	validator *validator.Validate
}

type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value"`
	Message string `json:"message"`
}

func NewValidator() *CustomValidator {
	v := validator.New()
	
	// Register custom validations here if needed
	// v.RegisterValidation("customtag", customValidationFunc)
	
	return &CustomValidator{
		validator: v,
	}
}

// Alias for compatibility
func NewCustomValidator() *CustomValidator {
	return NewValidator()
}

// GetValidator returns the underlying validator.Validate instance
func (cv *CustomValidator) GetValidator() *validator.Validate {
	return cv.validator
}

func (cv *CustomValidator) Validate(i interface{}) []ValidationError {
	var errors []ValidationError
	
	err := cv.validator.Struct(i)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   err.Field(),
				Tag:     err.Tag(),
				Value:   fmt.Sprintf("%v", err.Value()),
				Message: getErrorMessage(err),
			})
		}
	}
	
	return errors
}

func getErrorMessage(err validator.FieldError) string {
	field := strings.ToLower(err.Field())
	
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}