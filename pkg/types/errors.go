package types

import "fmt"

// Error codes for go-starter
const (
	ErrCodeUnknown          = "UNKNOWN"
	ErrCodeValidation       = "VALIDATION_ERROR"
	ErrCodeTemplateNotFound = "TEMPLATE_NOT_FOUND"
	ErrCodeGenerationError  = "GENERATION_ERROR"
	ErrCodeFileSystem       = "FILESYSTEM_ERROR"
	ErrCodeConfigError      = "CONFIG_ERROR"
)

// GoStarterError represents a go-starter specific error
type GoStarterError struct {
	Code    string
	Message string
	Cause   error
}

func (e *GoStarterError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *GoStarterError) Unwrap() error {
	return e.Cause
}

// NewError creates a new GoStarterError
func NewError(code, message string, cause error) *GoStarterError {
	return &GoStarterError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewValidationError creates a validation error
func NewValidationError(message string, cause error) *GoStarterError {
	return NewError(ErrCodeValidation, message, cause)
}

// NewTemplateNotFoundError creates a template not found error
func NewTemplateNotFoundError(templateID string) *GoStarterError {
	return NewError(ErrCodeTemplateNotFound, fmt.Sprintf("template '%s' not found", templateID), nil)
}

// NewGenerationError creates a generation error
func NewGenerationError(message string, cause error) *GoStarterError {
	return NewError(ErrCodeGenerationError, message, cause)
}

// NewFileSystemError creates a filesystem error
func NewFileSystemError(message string, cause error) *GoStarterError {
	return NewError(ErrCodeFileSystem, message, cause)
}

// NewConfigError creates a configuration error
func NewConfigError(message string, cause error) *GoStarterError {
	return NewError(ErrCodeConfigError, message, cause)
}
