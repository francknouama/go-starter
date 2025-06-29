package types

import (
	"errors"
	"strings"
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

func TestNewFileSystemError(t *testing.T) {
	// Test without cause
	err := NewFileSystemError("permission denied", nil)
	if err.Code != ErrCodeFileSystem {
		t.Errorf("Expected code %s, got %s", ErrCodeFileSystem, err.Code)
	}
	if err.Message != "permission denied" {
		t.Errorf("Expected message 'permission denied', got %s", err.Message)
	}

	// Test with cause
	cause := errors.New("underlying fs error")
	err = NewFileSystemError("failed to write", cause)
	expectedMsg := "[FILESYSTEM_ERROR] failed to write: underlying fs error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestNewGenerationError(t *testing.T) {
	err := NewGenerationError("template compilation failed", nil)
	if err.Code != ErrCodeGenerationError {
		t.Errorf("Expected code %s, got %s", ErrCodeGenerationError, err.Code)
	}
	if err.Message != "template compilation failed" {
		t.Errorf("Expected message 'template compilation failed', got %s", err.Message)
	}
}

func TestErrorConstants(t *testing.T) {
	// Test that error codes are properly defined
	expectedCodes := []string{
		ErrCodeValidation,
		ErrCodeTemplateNotFound,
		ErrCodeFileSystem,
		ErrCodeGenerationError,
	}

	for _, code := range expectedCodes {
		if code == "" {
			t.Errorf("Error code should not be empty")
		}
	}

	// Test uniqueness
	codeMap := make(map[string]bool)
	for _, code := range expectedCodes {
		if codeMap[code] {
			t.Errorf("Duplicate error code: %s", code)
		}
		codeMap[code] = true
	}
}

func TestGoStarterError_Is(t *testing.T) {
	err1 := NewValidationError("test", nil)
	err2 := NewValidationError("test", nil)
	err3 := NewFileSystemError("test", nil)

	// Note: errors.Is checks for equality, not just same type
	// These are different instances, so they won't be equal
	// We test that the errors have the expected codes instead
	if err1.Code != err2.Code {
		t.Error("Errors of same type should have same code")
	}

	// Errors with different codes should not be equal
	if err1.Code == err3.Code {
		t.Error("Errors with different types should have different codes")
	}
}

func TestGoStarterError_As(t *testing.T) {
	originalErr := NewValidationError("test validation", nil)
	wrappedErr := NewGenerationError("generation failed", originalErr)

	// Should be able to extract a GoStarterError from wrapped error
	var goStarterErr *GoStarterError
	if !errors.As(wrappedErr, &goStarterErr) {
		t.Error("Should be able to extract GoStarterError from wrapped error")
	}

	// The extracted error will be the wrapping error, not the original
	if goStarterErr.Code != ErrCodeGenerationError {
		t.Errorf("Expected generation error code, got %s", goStarterErr.Code)
	}

	// Test that we can unwrap to get the original error
	unwrapped := errors.Unwrap(wrappedErr)
	if unwrapped != originalErr {
		t.Error("Should be able to unwrap to get original validation error")
	}
}

func TestErrorChaining(t *testing.T) {
	// Test deep error chaining
	root := errors.New("root cause")
	middle := NewFileSystemError("middle error", root)
	top := NewGenerationError("top error", middle)

	// Test that the full chain is preserved
	errMsg := top.Error()
	if !containsAll(errMsg, "top error", "middle error", "root cause") {
		t.Errorf("Error message should contain all levels: %s", errMsg)
	}

	// Test unwrapping through the chain
	if errors.Unwrap(top) != middle {
		t.Error("Should unwrap to middle error")
	}

	if errors.Unwrap(middle) != root {
		t.Error("Should unwrap to root error")
	}
}

func TestErrorWithNilCause(t *testing.T) {
	err := NewError("TEST", "test message", nil)

	if err.Cause != nil {
		t.Error("Cause should be nil")
	}

	if errors.Unwrap(err) != nil {
		t.Error("Unwrap should return nil for error without cause")
	}
}

func TestErrorMessageFormatting(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		message  string
		cause    error
		expected string
	}{
		{
			name:     "simple error",
			code:     "TEST",
			message:  "test message",
			cause:    nil,
			expected: "[TEST] test message",
		},
		{
			name:     "error with cause",
			code:     "TEST",
			message:  "test message",
			cause:    errors.New("cause"),
			expected: "[TEST] test message: cause",
		},
		{
			name:     "empty message",
			code:     "TEST",
			message:  "",
			cause:    nil,
			expected: "[TEST] ",
		},
		{
			name:     "empty code",
			code:     "",
			message:  "test message",
			cause:    nil,
			expected: "[] test message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := NewError(tc.code, tc.message, tc.cause)
			if err.Error() != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, err.Error())
			}
		})
	}
}

func TestSpecificErrorTypes(t *testing.T) {
	// Test that each error constructor creates the right type
	validationErr := NewValidationError("validation failed", nil)
	if validationErr.Code != ErrCodeValidation {
		t.Errorf("Expected validation error code, got %s", validationErr.Code)
	}

	templateErr := NewTemplateNotFoundError("missing-template")
	if templateErr.Code != ErrCodeTemplateNotFound {
		t.Errorf("Expected template not found error code, got %s", templateErr.Code)
	}
	expectedTemplateMsg := "template 'missing-template' not found"
	if templateErr.Message != expectedTemplateMsg {
		t.Errorf("Expected '%s', got '%s'", expectedTemplateMsg, templateErr.Message)
	}

	fsErr := NewFileSystemError("file operation failed", nil)
	if fsErr.Code != ErrCodeFileSystem {
		t.Errorf("Expected filesystem error code, got %s", fsErr.Code)
	}

	genErr := NewGenerationError("generation failed", nil)
	if genErr.Code != ErrCodeGenerationError {
		t.Errorf("Expected generation error code, got %s", genErr.Code)
	}
}

// Helper function to check if a string contains all substrings
func containsAll(s string, substrings ...string) bool {
	for _, substr := range substrings {
		if !strings.Contains(s, substr) {
			return false
		}
	}
	return true
}
