package generator

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestErrorContext(t *testing.T) {
	t.Run("error context creation", func(t *testing.T) {
		context := &ErrorContext{
			Operation:   "test_operation",
			Component:   "test_component",
			ProjectName: "test_project",
			Template:    "test_template",
			File:        "test_file.go",
		}
		
		if context.Operation != "test_operation" {
			t.Errorf("Expected operation 'test_operation', got '%s'", context.Operation)
		}
		
		if context.Component != "test_component" {
			t.Errorf("Expected component 'test_component', got '%s'", context.Component)
		}
	})
}

func TestGenerationError(t *testing.T) {
	t.Run("generation error creation", func(t *testing.T) {
		cause := errors.New("underlying error")
		context := &ErrorContext{
			Operation: "test_operation",
			Component: "test_component",
		}
		
		genErr := NewGenerationError(types.ErrCodeValidation, "validation failed", cause, context)
		
		if genErr == nil {
			t.Error("Expected generation error to not be nil")
		}
		
		if genErr.Context == nil {
			t.Error("Expected context to not be nil")
		}
		
		if genErr.Context.Operation != "test_operation" {
			t.Errorf("Expected operation 'test_operation', got '%s'", genErr.Context.Operation)
		}
		
		if genErr.Cause != cause {
			t.Error("Expected cause to be preserved")
		}
		
		if len(genErr.Suggestions) == 0 {
			t.Error("Expected suggestions to be generated")
		}
	})
	
	t.Run("error message formatting", func(t *testing.T) {
		genErr := NewGenerationError(
			types.ErrCodeValidation,
			"test error",
			nil,
			&ErrorContext{Operation: "test_op"},
		)
		
		msg := genErr.Error()
		if msg == "" {
			t.Error("Expected error message to not be empty")
		}
		
		detailed := genErr.DetailedError()
		if detailed == "" {
			t.Error("Expected detailed error to not be empty")
		}
		
		if len(detailed) <= len(msg) {
			t.Error("Expected detailed error to be longer than basic error")
		}
	})
	
	t.Run("nil context handling", func(t *testing.T) {
		genErr := NewGenerationError(types.ErrCodeGenerationError, "test", nil, nil)
		
		if genErr.Context == nil {
			t.Error("Expected context to be created when nil is passed")
		}
		
		if genErr.Context.Operation != "unknown" {
			t.Errorf("Expected default operation 'unknown', got '%s'", genErr.Context.Operation)
		}
	})
}

func TestErrorHandler(t *testing.T) {
	t.Run("error handler creation", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		if handler == nil {
			t.Error("Expected handler to not be nil")
		}
		
		if handler.logger != logger {
			t.Error("Expected handler to use provided logger")
		}
	})
	
	t.Run("handle nil error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		err := handler.Handle(nil)
		if err != nil {
			t.Errorf("Expected nil error to remain nil, got %v", err)
		}
	})
	
	t.Run("handle generation error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		original := NewGenerationError(types.ErrCodeValidation, "test error", nil, nil)
		handled := handler.Handle(original)
		
		if handled != original {
			t.Error("Expected generation error to be returned unchanged")
		}
	})
	
	t.Run("handle go starter error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		gsErr := types.NewError(types.ErrCodeFileSystem, "file error", nil)
		handled := handler.Handle(gsErr)
		
		genErr, ok := handled.(*GenerationError)
		if !ok {
			t.Error("Expected GoStarterError to be wrapped in GenerationError")
		}
		
		if genErr.Code != types.ErrCodeFileSystem {
			t.Errorf("Expected error code %s, got %s", types.ErrCodeFileSystem, genErr.Code)
		}
	})
	
	t.Run("handle generic error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		genericErr := errors.New("generic error")
		handled := handler.Handle(genericErr)
		
		genErr, ok := handled.(*GenerationError)
		if !ok {
			t.Error("Expected generic error to be wrapped in GenerationError")
		}
		
		if genErr.Code != types.ErrCodeGenerationError {
			t.Errorf("Expected error code %s, got %s", types.ErrCodeGenerationError, genErr.Code)
		}
	})
	
	t.Run("handle with context", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		context := &ErrorContext{
			Operation: "test_op",
			Component: "test_component",
		}
		
		genericErr := errors.New("test error")
		handled := handler.HandleWithContext(genericErr, context)
		
		genErr, ok := handled.(*GenerationError)
		if !ok {
			t.Error("Expected error to be wrapped in GenerationError")
		}
		
		if genErr.Context.Operation != "test_op" {
			t.Errorf("Expected context operation 'test_op', got '%s'", genErr.Context.Operation)
		}
	})
}

func TestRecoverableOperation(t *testing.T) {
	t.Run("successful operation", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		executed := false
		op := &RecoverableOperation{
			Name:        "test_op",
			Description: "test operation",
			Execute: func() error {
				executed = true
				return nil
			},
			MaxRetries: 2,
			RetryDelay: 1 * time.Millisecond,
		}
		
		err := handler.ExecuteWithRecovery(op)
		if err != nil {
			t.Errorf("Expected successful operation to return nil, got %v", err)
		}
		
		if !executed {
			t.Error("Expected operation to be executed")
		}
	})
	
	t.Run("operation with retries", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		attempts := 0
		op := &RecoverableOperation{
			Name: "failing_op",
			Execute: func() error {
				attempts++
				if attempts < 3 {
					return errors.New("temporary failure")
				}
				return nil
			},
			MaxRetries: 3,
			RetryDelay: 1 * time.Millisecond,
		}
		
		err := handler.ExecuteWithRecovery(op)
		if err != nil {
			t.Errorf("Expected operation to succeed after retries, got %v", err)
		}
		
		if attempts != 3 {
			t.Errorf("Expected 3 attempts, got %d", attempts)
		}
	})
	
	t.Run("operation with rollback", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		rolledBack := false
		op := &RecoverableOperation{
			Name: "failing_op",
			Execute: func() error {
				return errors.New("persistent failure")
			},
			Rollback: func() error {
				rolledBack = true
				return nil
			},
			MaxRetries: 1,
			RetryDelay: 1 * time.Millisecond,
		}
		
		err := handler.ExecuteWithRecovery(op)
		if err == nil {
			t.Error("Expected operation to fail")
		}
		
		if !rolledBack {
			t.Error("Expected rollback to be executed")
		}
	})
	
	t.Run("panic recovery", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		handler := NewErrorHandler(logger)
		
		op := &RecoverableOperation{
			Name: "panicking_op",
			Execute: func() error {
				panic("test panic")
			},
			MaxRetries: 0,
		}
		
		err := handler.ExecuteWithRecovery(op)
		if err == nil {
			t.Error("Expected panic to be converted to error")
		}
		
		genErr, ok := err.(*GenerationError)
		if !ok {
			t.Error("Expected panic to be wrapped in GenerationError")
		}
		
		if genErr.Context.Component != "execution" {
			t.Errorf("Expected component 'execution', got '%s'", genErr.Context.Component)
		}
	})
}

func TestValidationErrors(t *testing.T) {
	t.Run("validation errors collection", func(t *testing.T) {
		ve := NewValidationErrors()
		if ve == nil {
			t.Error("Expected validation errors to not be nil")
		}
		
		if ve.HasErrors() {
			t.Error("Expected new validation errors to be empty")
		}
		
		if ve.First() != nil {
			t.Error("Expected first error to be nil when empty")
		}
		
		err := ve.Error()
		if err != "no validation errors" {
			t.Errorf("Expected 'no validation errors', got '%s'", err)
		}
	})
	
	t.Run("add errors", func(t *testing.T) {
		ve := NewValidationErrors()
		
		// Add nil error (should be ignored)
		ve.Add(nil)
		if ve.HasErrors() {
			t.Error("Expected nil errors to be ignored")
		}
		
		// Add real error
		ve.Add(errors.New("test error"))
		if !ve.HasErrors() {
			t.Error("Expected errors to be present")
		}
		
		first := ve.First()
		if first == nil {
			t.Error("Expected first error to not be nil")
		}
		
		if first.Error() != "test error" {
			t.Errorf("Expected 'test error', got '%s'", first.Error())
		}
	})
	
	t.Run("add field errors", func(t *testing.T) {
		ve := NewValidationErrors()
		
		ve.AddField("name", "cannot be empty")
		if !ve.HasErrors() {
			t.Error("Expected field error to be added")
		}
		
		first := ve.First()
		genErr, ok := first.(*GenerationError)
		if !ok {
			t.Error("Expected field error to be GenerationError")
		}
		
		if genErr.Code != types.ErrCodeValidation {
			t.Errorf("Expected validation error code, got %s", genErr.Code)
		}
		
		if genErr.Context.Component != "field" {
			t.Errorf("Expected component 'field', got '%s'", genErr.Context.Component)
		}
	})
	
	t.Run("multiple errors", func(t *testing.T) {
		ve := NewValidationErrors()
		
		ve.Add(errors.New("error 1"))
		ve.Add(errors.New("error 2"))
		ve.AddField("field1", "invalid")
		
		if len(ve.Errors) != 3 {
			t.Errorf("Expected 3 errors, got %d", len(ve.Errors))
		}
		
		err := ve.Error()
		if !strings.Contains(err, "multiple validation errors") {
			t.Errorf("Expected multiple validation errors message, got '%s'", err)
		}
	})
}

func TestErrorContextCreators(t *testing.T) {
	t.Run("file error context", func(t *testing.T) {
		ctx := FileErrorContext("write", "/path/to/file.go")
		
		if ctx.Operation != "write" {
			t.Errorf("Expected operation 'write', got '%s'", ctx.Operation)
		}
		
		if ctx.Component != "file" {
			t.Errorf("Expected component 'file', got '%s'", ctx.Component)
		}
		
		if ctx.File != "/path/to/file.go" {
			t.Errorf("Expected file '/path/to/file.go', got '%s'", ctx.File)
		}
	})
	
	t.Run("template error context", func(t *testing.T) {
		ctx := TemplateErrorContext("parse", "web-api", "main.go.tmpl")
		
		if ctx.Operation != "parse" {
			t.Errorf("Expected operation 'parse', got '%s'", ctx.Operation)
		}
		
		if ctx.Component != "template" {
			t.Errorf("Expected component 'template', got '%s'", ctx.Component)
		}
		
		if ctx.Template != "web-api" {
			t.Errorf("Expected template 'web-api', got '%s'", ctx.Template)
		}
		
		if ctx.File != "main.go.tmpl" {
			t.Errorf("Expected file 'main.go.tmpl', got '%s'", ctx.File)
		}
	})
	
	t.Run("project error context", func(t *testing.T) {
		ctx := ProjectErrorContext("generate", "my-project", "cli")
		
		if ctx.Operation != "generate" {
			t.Errorf("Expected operation 'generate', got '%s'", ctx.Operation)
		}
		
		if ctx.Component != "project" {
			t.Errorf("Expected component 'project', got '%s'", ctx.Component)
		}
		
		if ctx.ProjectName != "my-project" {
			t.Errorf("Expected project 'my-project', got '%s'", ctx.ProjectName)
		}
		
		if ctx.Template != "cli" {
			t.Errorf("Expected template 'cli', got '%s'", ctx.Template)
		}
	})
	
	t.Run("step error context", func(t *testing.T) {
		metadata := map[string]interface{}{
			"step_number": 1,
			"total_steps": 5,
		}
		
		ctx := StepErrorContext("validation", "validate_config", metadata)
		
		if ctx.Operation != "validation" {
			t.Errorf("Expected operation 'validation', got '%s'", ctx.Operation)
		}
		
		if ctx.Component != "step" {
			t.Errorf("Expected component 'step', got '%s'", ctx.Component)
		}
		
		if ctx.Step != "validate_config" {
			t.Errorf("Expected step 'validate_config', got '%s'", ctx.Step)
		}
		
		if ctx.Metadata["step_number"] != 1 {
			t.Errorf("Expected step_number 1, got %v", ctx.Metadata["step_number"])
		}
	})
}

func TestSuggestionGeneration(t *testing.T) {
	t.Run("validation suggestions", func(t *testing.T) {
		suggestions := generateSuggestions(types.ErrCodeValidation, "validation error", nil)
		
		if len(suggestions) == 0 {
			t.Error("Expected validation suggestions to be generated")
		}
		
		hasValidationSuggestion := false
		for _, suggestion := range suggestions {
			if strings.Contains(suggestion, "required fields") {
				hasValidationSuggestion = true
				break
			}
		}
		
		if !hasValidationSuggestion {
			t.Error("Expected at least one validation-specific suggestion")
		}
	})
	
	t.Run("filesystem suggestions", func(t *testing.T) {
		suggestions := generateSuggestions(types.ErrCodeFileSystem, "file system error", nil)
		
		hasPermissionSuggestion := false
		for _, suggestion := range suggestions {
			if strings.Contains(suggestion, "permissions") {
				hasPermissionSuggestion = true
				break
			}
		}
		
		if !hasPermissionSuggestion {
			t.Error("Expected filesystem permission suggestion")
		}
	})
	
	t.Run("context-specific suggestions", func(t *testing.T) {
		suggestions := generateSuggestions(types.ErrCodeGenerationError, "permission denied", nil)
		
		hasPermissionSuggestion := false
		for _, suggestion := range suggestions {
			if strings.Contains(suggestion, "sudo") || strings.Contains(suggestion, "permissions") {
				hasPermissionSuggestion = true
				break
			}
		}
		
		if !hasPermissionSuggestion {
			t.Error("Expected permission-specific suggestion")
		}
	})
}