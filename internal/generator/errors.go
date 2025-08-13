package generator

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

// ErrorContext provides context for generator errors
type ErrorContext struct {
	Operation   string                 `json:"operation"`
	Component   string                 `json:"component"`
	ProjectName string                 `json:"project_name,omitempty"`
	Template    string                 `json:"template,omitempty"`
	File        string                 `json:"file,omitempty"`
	Step        string                 `json:"step,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	StackTrace  []string               `json:"stack_trace,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// GenerationError represents a structured error during project generation
type GenerationError struct {
	*types.GoStarterError
	Context     *ErrorContext `json:"context"`
	Suggestions []string      `json:"suggestions,omitempty"`
}

// NewGenerationError creates a new generation error with context
func NewGenerationError(code string, message string, cause error, context *ErrorContext) *GenerationError {
	baseError := types.NewError(code, message, cause)
	
	// Add stack trace in debug mode
	var stackTrace []string
	if context != nil && context.Component == "debug" {
		stackTrace = captureStackTrace()
	}
	
	if context == nil {
		context = &ErrorContext{
			Operation:  "unknown",
			Component:  "generator",
			Timestamp:  time.Now(),
			StackTrace: stackTrace,
		}
	} else {
		context.Timestamp = time.Now()
		if len(context.StackTrace) == 0 {
			context.StackTrace = stackTrace
		}
	}
	
	return &GenerationError{
		GoStarterError: baseError,
		Context:        context,
		Suggestions:    generateSuggestions(code, message, cause),
	}
}

// Error implements the error interface
func (ge *GenerationError) Error() string {
	msg := ge.GoStarterError.Error()
	if ge.Context != nil && ge.Context.Operation != "" {
		msg = fmt.Sprintf("%s (operation: %s)", msg, ge.Context.Operation)
	}
	return msg
}

// DetailedError returns a detailed error message with context
func (ge *GenerationError) DetailedError() string {
	var parts []string
	
	parts = append(parts, fmt.Sprintf("Error: %s", ge.Message))
	
	if ge.Context != nil {
		if ge.Context.Operation != "" {
			parts = append(parts, fmt.Sprintf("Operation: %s", ge.Context.Operation))
		}
		if ge.Context.Component != "" {
			parts = append(parts, fmt.Sprintf("Component: %s", ge.Context.Component))
		}
		if ge.Context.ProjectName != "" {
			parts = append(parts, fmt.Sprintf("Project: %s", ge.Context.ProjectName))
		}
		if ge.Context.Template != "" {
			parts = append(parts, fmt.Sprintf("Template: %s", ge.Context.Template))
		}
		if ge.Context.File != "" {
			parts = append(parts, fmt.Sprintf("File: %s", ge.Context.File))
		}
		if ge.Context.Step != "" {
			parts = append(parts, fmt.Sprintf("Step: %s", ge.Context.Step))
		}
	}
	
	if ge.Cause != nil {
		parts = append(parts, fmt.Sprintf("Cause: %v", ge.Cause))
	}
	
	if len(ge.Suggestions) > 0 {
		parts = append(parts, "Suggestions:")
		for _, suggestion := range ge.Suggestions {
			parts = append(parts, fmt.Sprintf("  - %s", suggestion))
		}
	}
	
	return strings.Join(parts, "\n")
}

// ErrorHandler provides centralized error handling for the generator
type ErrorHandler struct {
	logger *GeneratorLogger
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *GeneratorLogger) *ErrorHandler {
	return &ErrorHandler{
		logger: logger,
	}
}

// Handle processes and logs an error appropriately
func (eh *ErrorHandler) Handle(err error) error {
	if err == nil {
		return nil
	}
	
	// Check if it's already a generation error
	if genErr, ok := err.(*GenerationError); ok {
		eh.handleGenerationError(genErr)
		return genErr
	}
	
	// Check if it's a GoStarterError
	if gsErr, ok := err.(*types.GoStarterError); ok {
		genErr := &GenerationError{
			GoStarterError: gsErr,
			Context: &ErrorContext{
				Operation: "unknown",
				Component: "generator",
				Timestamp: time.Now(),
			},
			Suggestions: generateSuggestions(gsErr.Code, gsErr.Message, gsErr.Cause),
		}
		eh.handleGenerationError(genErr)
		return genErr
	}
	
	// Handle generic error
	genErr := NewGenerationError(
		types.ErrCodeGenerationError,
		err.Error(),
		err,
		&ErrorContext{
			Operation: "unknown",
			Component: "generator",
		},
	)
	eh.handleGenerationError(genErr)
	return genErr
}

// HandleWithContext processes an error with specific context
func (eh *ErrorHandler) HandleWithContext(err error, context *ErrorContext) error {
	if err == nil {
		return nil
	}
	
	// Enhance existing generation error with new context
	if genErr, ok := err.(*GenerationError); ok {
		genErr = eh.enhanceErrorContext(genErr, context)
		eh.handleGenerationError(genErr)
		return genErr
	}
	
	// Create new generation error with context
	var code string
	if gsErr, ok := err.(*types.GoStarterError); ok {
		code = gsErr.Code
	} else {
		code = types.ErrCodeGenerationError
	}
	
	genErr := NewGenerationError(code, err.Error(), err, context)
	eh.handleGenerationError(genErr)
	return genErr
}

// handleGenerationError logs and processes a generation error
func (eh *ErrorHandler) handleGenerationError(genErr *GenerationError) {
	// Log based on error severity
	switch genErr.Code {
	case types.ErrCodeValidation:
		eh.logger.Warning("Validation error: %s", genErr.Message)
	case types.ErrCodeFileSystem:
		eh.logger.Error("File system error: %s", genErr.Message)
	case types.ErrCodeTemplateNotFound:
		eh.logger.Error("Template not found: %s", genErr.Message)
	case types.ErrCodeGenerationError:
		eh.logger.Error("Template parsing error: %s", genErr.Message)
	case types.ErrCodeDependency:
		eh.logger.Warning("Dependency error: %s", genErr.Message)
	default:
		eh.logger.Error("Generation error: %s", genErr.Message)
	}
	
	// Log context details in debug mode
	if eh.logger.level >= LogLevelDebug && genErr.Context != nil {
		context := map[string]interface{}{
			"operation": genErr.Context.Operation,
			"component": genErr.Context.Component,
		}
		
		if genErr.Context.ProjectName != "" {
			context["project"] = genErr.Context.ProjectName
		}
		if genErr.Context.Template != "" {
			context["template"] = genErr.Context.Template
		}
		if genErr.Context.File != "" {
			context["file"] = genErr.Context.File
		}
		
		eh.logger.Debug("Error context: %+v", context)
	}
	
	// Log suggestions
	if len(genErr.Suggestions) > 0 {
		eh.logger.Info("Suggestions:")
		for _, suggestion := range genErr.Suggestions {
			eh.logger.Info("  - %s", suggestion)
		}
	}
}

// enhanceErrorContext merges additional context into an existing error
func (eh *ErrorHandler) enhanceErrorContext(genErr *GenerationError, newContext *ErrorContext) *GenerationError {
	if newContext == nil {
		return genErr
	}
	
	if genErr.Context == nil {
		genErr.Context = newContext
		return genErr
	}
	
	// Merge contexts, preferring new values
	if newContext.Operation != "" {
		genErr.Context.Operation = newContext.Operation
	}
	if newContext.Component != "" {
		genErr.Context.Component = newContext.Component
	}
	if newContext.ProjectName != "" {
		genErr.Context.ProjectName = newContext.ProjectName
	}
	if newContext.Template != "" {
		genErr.Context.Template = newContext.Template
	}
	if newContext.File != "" {
		genErr.Context.File = newContext.File
	}
	if newContext.Step != "" {
		genErr.Context.Step = newContext.Step
	}
	
	// Merge metadata
	if genErr.Context.Metadata == nil {
		genErr.Context.Metadata = make(map[string]interface{})
	}
	for k, v := range newContext.Metadata {
		genErr.Context.Metadata[k] = v
	}
	
	return genErr
}

// RecoverableOperation represents an operation that can be recovered from errors
type RecoverableOperation struct {
	Name        string
	Description string
	Execute     func() error
	Rollback    func() error
	MaxRetries  int
	RetryDelay  time.Duration
}

// ExecuteWithRecovery executes an operation with error recovery and retry logic
func (eh *ErrorHandler) ExecuteWithRecovery(op *RecoverableOperation) error {
	var lastErr error
	
	for attempt := 0; attempt <= op.MaxRetries; attempt++ {
		// Set up panic recovery
		err := func() (err error) {
			defer func() {
				if r := recover(); r != nil {
					panicErr := fmt.Errorf("panic in %s: %v", op.Name, r)
					err = NewGenerationError(
						types.ErrCodeGenerationError,
						panicErr.Error(),
						panicErr,
						&ErrorContext{
							Operation:  op.Name,
							Component:  "recovery",
							StackTrace: captureStackTrace(),
						},
					)
				}
			}()
			
			return op.Execute()
		}()
		
		if err == nil {
			if attempt > 0 {
				eh.logger.Success("Operation '%s' succeeded after %d retries", op.Name, attempt)
			}
			return nil
		}
		
		lastErr = eh.HandleWithContext(err, &ErrorContext{
			Operation: op.Name,
			Component: "execution",
		})
		
		if attempt < op.MaxRetries {
			eh.logger.Warning("Operation '%s' failed, retrying in %v (attempt %d/%d)", 
				op.Name, op.RetryDelay, attempt+1, op.MaxRetries)
			time.Sleep(op.RetryDelay)
		}
	}
	
	// All retries failed, attempt rollback if available
	if op.Rollback != nil {
		eh.logger.Warning("All retries failed for '%s', attempting rollback", op.Name)
		if rollbackErr := op.Rollback(); rollbackErr != nil {
			eh.logger.Error("Rollback failed for '%s': %v", op.Name, rollbackErr)
		} else {
			eh.logger.Success("Rollback completed for '%s'", op.Name)
		}
	}
	
	return lastErr
}

// generateSuggestions creates helpful suggestions based on error type and context
func generateSuggestions(code string, message string, cause error) []string {
	var suggestions []string
	
	switch code {
	case types.ErrCodeValidation:
		suggestions = append(suggestions, "Check that all required fields are provided")
		suggestions = append(suggestions, "Verify the project name uses valid characters (letters, numbers, hyphens)")
		suggestions = append(suggestions, "Ensure the module path follows Go module naming conventions")
		
	case types.ErrCodeFileSystem:
		suggestions = append(suggestions, "Check that you have write permissions to the output directory")
		suggestions = append(suggestions, "Ensure the output directory exists or can be created")
		suggestions = append(suggestions, "Verify you have sufficient disk space")
		
	case types.ErrCodeTemplateNotFound:
		suggestions = append(suggestions, "Run 'go-starter list' to see available templates")
		suggestions = append(suggestions, "Check the template name spelling")
		suggestions = append(suggestions, "Ensure you're using the latest version of go-starter")
		
	case types.ErrCodeGenerationError:
		suggestions = append(suggestions, "Report this as a bug - template parsing should always work")
		suggestions = append(suggestions, "Try using a different template as a workaround")
		
	case types.ErrCodeDependency:
		suggestions = append(suggestions, "Check your internet connection")
		suggestions = append(suggestions, "Verify Go is installed and in your PATH")
		suggestions = append(suggestions, "Try running 'go mod tidy' manually in the generated project")
	}
	
	// Add context-specific suggestions
	if strings.Contains(message, "permission denied") {
		suggestions = append(suggestions, "Run with appropriate permissions (e.g., sudo on Unix systems)")
	}
	
	if strings.Contains(message, "no space left") {
		suggestions = append(suggestions, "Free up disk space and try again")
	}
	
	return suggestions
}

// captureStackTrace captures the current stack trace
func captureStackTrace() []string {
	var stackTrace []string
	
	// Skip the first few frames (this function, NewGenerationError, etc.)
	for i := 3; i < 10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		
		fn := runtime.FuncForPC(pc)
		if fn == nil {
			break
		}
		
		// Skip runtime functions
		name := fn.Name()
		if strings.Contains(name, "runtime.") {
			continue
		}
		
		stackTrace = append(stackTrace, fmt.Sprintf("%s:%d %s", file, line, name))
	}
	
	return stackTrace
}

// Common error context creators

// FileErrorContext creates error context for file operations
func FileErrorContext(operation, file string) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Component: "file",
		File:      file,
	}
}

// TemplateErrorContext creates error context for template operations
func TemplateErrorContext(operation, template, file string) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Component: "template",
		Template:  template,
		File:      file,
	}
}

// ProjectErrorContext creates error context for project operations
func ProjectErrorContext(operation, projectName, template string) *ErrorContext {
	return &ErrorContext{
		Operation:   operation,
		Component:   "project",
		ProjectName: projectName,
		Template:    template,
	}
}

// StepErrorContext creates error context for generation steps
func StepErrorContext(operation, step string, metadata map[string]interface{}) *ErrorContext {
	return &ErrorContext{
		Operation: operation,
		Component: "step",
		Step:      step,
		Metadata:  metadata,
	}
}

// ValidationErrors collects multiple validation errors
type ValidationErrors struct {
	Errors []error
}

// NewValidationErrors creates a new validation errors collection
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make([]error, 0),
	}
}

// Add adds a validation error
func (ve *ValidationErrors) Add(err error) {
	if err != nil {
		ve.Errors = append(ve.Errors, err)
	}
}

// AddField adds a field-specific validation error
func (ve *ValidationErrors) AddField(field, message string) {
	ve.Add(NewGenerationError(
		types.ErrCodeValidation,
		fmt.Sprintf("invalid %s: %s", field, message),
		nil,
		&ErrorContext{
			Operation: "validation",
			Component: "field",
			Metadata: map[string]interface{}{
				"field": field,
			},
		},
	))
}

// HasErrors returns true if there are validation errors
func (ve *ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	if len(ve.Errors) == 0 {
		return "no validation errors"
	}
	
	if len(ve.Errors) == 1 {
		return ve.Errors[0].Error()
	}
	
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, err.Error())
	}
	
	return fmt.Sprintf("multiple validation errors: %s", strings.Join(messages, "; "))
}

// First returns the first error or nil
func (ve *ValidationErrors) First() error {
	if len(ve.Errors) > 0 {
		return ve.Errors[0]
	}
	return nil
}