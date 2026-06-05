package utils

import "fmt"

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) error {
	return &ValidationError{Message: message}
}

func ValidateNonEmpty(value, fieldName string) error {
	if value == "" {
		return NewValidationError(fieldName + " is required")
	}
	return nil
}

func ValidatePositive(value float64, fieldName string) error {
	if value < 0 {
		return NewValidationError(fmt.Sprintf("%s must be positive", fieldName))
	}
	return nil
}

func ValidateNonNegative(value int, fieldName string) error {
	if value < 0 {
		return NewValidationError(fmt.Sprintf("%s must be non-negative", fieldName))
	}
	return nil
}

func ValidateMinLength(value string, minLength int, fieldName string) error {
	if len(value) < minLength {
		return NewValidationError(fmt.Sprintf("%s must be at least %d characters", fieldName, minLength))
	}
	return nil
}

func ValidateUUIDArray(arr []interface{}, fieldName string) error {
	if len(arr) == 0 {
		return NewValidationError(fmt.Sprintf("at least one %s is required", fieldName))
	}
	return nil
}
