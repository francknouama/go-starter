package types

import (
	"errors"
	"testing"
)

func TestGoStarterError(t *testing.T) {
	tests := []struct {
		name     string
		error    *GoStarterError
		expected string
	}{
		{
			name:     "error without cause",
			error:    NewError("TEST_CODE", "test message", nil),
			expected: "[TEST_CODE] test message",
		},
		{
			name:     "error with cause",
			error:    NewError("TEST_CODE", "test message", errors.New("underlying error")),
			expected: "[TEST_CODE] test message: underlying error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.error.Error(); got != tt.expected {
				t.Errorf("Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	err := NewValidationError("invalid input", nil)
	if err.Code != ErrCodeValidation {
		t.Errorf("Expected code %s, got %s", ErrCodeValidation, err.Code)
	}
	if err.Message != "invalid input" {
		t.Errorf("Expected message 'invalid input', got %s", err.Message)
	}
}

func TestNewTemplateNotFoundError(t *testing.T) {
	err := NewTemplateNotFoundError("web-api")
	if err.Code != ErrCodeTemplateNotFound {
		t.Errorf("Expected code %s, got %s", ErrCodeTemplateNotFound, err.Code)
	}
	expected := "template 'web-api' not found"
	if err.Message != expected {
		t.Errorf("Expected message '%s', got %s", expected, err.Message)
	}
}

func TestErrorUnwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := NewError("TEST_CODE", "test message", cause)

	if unwrapped := errors.Unwrap(err); unwrapped != cause {
		t.Errorf("Expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}
