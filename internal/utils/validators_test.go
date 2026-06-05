package utils

import (
	"testing"
)

func TestValidateNonEmpty(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		fieldName string
		wantErr   bool
	}{
		{"Valid non-empty string", "test", "field", false},
		{"Empty string", "", "field", true},
		{"Whitespace only", "   ", "field", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonEmpty(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonEmpty() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidatePositive(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		fieldName string
		wantErr   bool
	}{
		{"Positive number", 100.50, "price", false},
		{"Zero (accepted)", 0.0, "price", false},
		{"Negative number", -50.0, "price", true},
		{"Small positive", 0.01, "price", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePositive(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePositive() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateNonNegative(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		fieldName string
		wantErr   bool
	}{
		{"Positive number", 100, "stock", false},
		{"Zero", 0, "stock", false},
		{"Negative number", -1, "stock", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNonNegative(tt.value, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNonNegative() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateMinLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		minLength int
		fieldName string
		wantErr   bool
	}{
		{"Valid length", "password123", 6, "password", false},
		{"Exact minimum length", "pass", 4, "password", false},
		{"Below minimum", "pass", 5, "password", true},
		{"Empty string", "", 1, "password", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMinLength(tt.value, tt.minLength, tt.fieldName)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMinLength() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationError(t *testing.T) {
	err := NewValidationError("test error message")
	if err == nil {
		t.Fatal("NewValidationError returned nil")
	}

	if err.Error() != "test error message" {
		t.Errorf("Error message mismatch: got %s, want 'test error message'", err.Error())
	}
}
