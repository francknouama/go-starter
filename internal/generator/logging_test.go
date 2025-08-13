package generator

import (
	"context"
	"testing"
	"time"
)

func TestGeneratorLogger(t *testing.T) {
	t.Run("logger creation", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		if logger == nil {
			t.Error("Expected logger to not be nil")
		}
		
		if logger.level != LogLevelInfo {
			t.Errorf("Expected level %v, got %v", LogLevelInfo, logger.level)
		}
		
		if logger.context == nil {
			t.Error("Expected context map to be initialized")
		}
	})
	
	t.Run("log levels", func(t *testing.T) {
		tests := []struct {
			level    LogLevel
			expected string
		}{
			{LogLevelSilent, "silent"},
			{LogLevelError, "error"},
			{LogLevelWarn, "warn"},
			{LogLevelInfo, "info"},
			{LogLevelDebug, "debug"},
		}
		
		for _, tt := range tests {
			if tt.level.String() != tt.expected {
				t.Errorf("Level %d: expected %q, got %q", tt.level, tt.expected, tt.level.String())
			}
		}
	})
	
	t.Run("context manipulation", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		
		// Test WithField
		loggerWithField := logger.WithField("key", "value")
		if loggerWithField == logger {
			t.Error("Expected new logger instance")
		}
		
		if loggerWithField.context["key"] != "value" {
			t.Error("Expected field to be set in context")
		}
		
		// Test WithFields
		fields := map[string]interface{}{
			"field1": "value1",
			"field2": 42,
		}
		loggerWithFields := logger.WithFields(fields)
		
		if loggerWithFields.context["field1"] != "value1" {
			t.Error("Expected field1 to be set")
		}
		if loggerWithFields.context["field2"] != 42 {
			t.Error("Expected field2 to be set")
		}
	})
	
	t.Run("duration tracking", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		time.Sleep(10 * time.Millisecond) // Small delay for testing
		
		duration := logger.Duration()
		if duration < 10*time.Millisecond {
			t.Error("Expected duration to be at least 10ms")
		}
	})
	
	t.Run("context logger", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		ctx := context.Background()
		
		contextLogger := logger.WithContext(ctx)
		if contextLogger == nil {
			t.Error("Expected context logger to not be nil")
		}
	})
}

func TestProgressTracker(t *testing.T) {
	t.Run("progress tracking", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1", "Step 2", "Step 3"}
		
		tracker := NewProgressTracker(logger, steps)
		if tracker == nil {
			t.Error("Expected tracker to not be nil")
		}
		
		if tracker.totalSteps != 3 {
			t.Errorf("Expected 3 total steps, got %d", tracker.totalSteps)
		}
		
		if tracker.currentStep != 0 {
			t.Errorf("Expected current step to be 0, got %d", tracker.currentStep)
		}
		
		// Test step progression
		tracker.NextStep()
		if tracker.currentStep != 1 {
			t.Errorf("Expected current step to be 1, got %d", tracker.currentStep)
		}
		
		tracker.NextStep()
		tracker.NextStep()
		if tracker.currentStep != 3 {
			t.Errorf("Expected current step to be 3, got %d", tracker.currentStep)
		}
		
		// Test beyond bounds
		tracker.NextStep() // Should not increase beyond totalSteps
		if tracker.currentStep != 3 {
			t.Errorf("Expected current step to remain 3, got %d", tracker.currentStep)
		}
	})
	
	t.Run("completion and error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Test Step"}
		tracker := NewProgressTracker(logger, steps)
		
		// These should not panic
		tracker.Complete("Test completed")
		tracker.Error("Test error", nil)
	})
}

func TestLoggerOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		opts := DefaultLoggerOptions()
		if opts == nil {
			t.Error("Expected options to not be nil")
		}
		
		if opts.Level != LogLevelInfo {
			t.Errorf("Expected default level %v, got %v", LogLevelInfo, opts.Level)
		}
		
		if !opts.ShowColors {
			t.Error("Expected colors to be enabled by default")
		}
	})
	
	t.Run("quiet mode", func(t *testing.T) {
		opts := DefaultLoggerOptions().SetQuietMode()
		if opts.Level != LogLevelWarn {
			t.Errorf("Expected quiet mode level %v, got %v", LogLevelWarn, opts.Level)
		}
	})
	
	t.Run("verbose mode", func(t *testing.T) {
		opts := DefaultLoggerOptions().SetVerboseMode()
		if opts.Level != LogLevelDebug {
			t.Errorf("Expected verbose mode level %v, got %v", LogLevelDebug, opts.Level)
		}
		
		if !opts.ShowTimestamps {
			t.Error("Expected timestamps to be enabled in verbose mode")
		}
	})
	
	t.Run("silent mode", func(t *testing.T) {
		opts := DefaultLoggerOptions().SetSilentMode()
		if opts.Level != LogLevelSilent {
			t.Errorf("Expected silent mode level %v, got %v", LogLevelSilent, opts.Level)
		}
	})
}

func TestFallbackLogger(t *testing.T) {
	t.Run("fallback logger methods", func(t *testing.T) {
		logger := &fallbackLogger{level: LogLevelInfo}
		
		// These should not panic
		logger.Debug("debug message")
		logger.Info("info message") 
		logger.Warn("warn message")
		logger.Error("error message")
		
		logger.DebugWith("debug", nil)
		logger.InfoWith("info", nil)
		logger.WarnWith("warn", nil)
		logger.ErrorWith("error", nil)
		
		ctxLogger := logger.WithContext(context.Background())
		if ctxLogger != logger {
			t.Error("Expected context logger to return same instance")
		}
		
		fieldsLogger := logger.WithFields(nil)
		if fieldsLogger != logger {
			t.Error("Expected fields logger to return same instance")
		}
		
		logger.SetLevel(0)
		logger.SetOutput(nil)
		logger.DisableColor()
		
		err := logger.Sync()
		if err != nil {
			t.Errorf("Expected sync to return nil, got %v", err)
		}
	})
	
	t.Run("level filtering", func(t *testing.T) {
		// Test with different log levels
		levels := []LogLevel{LogLevelSilent, LogLevelError, LogLevelWarn, LogLevelInfo, LogLevelDebug}
		
		for _, level := range levels {
			logger := &fallbackLogger{level: level}
			
			// These should not panic regardless of level
			logger.Debug("debug")
			logger.Info("info")
			logger.Warn("warn")
			logger.Error("error")
		}
	})
}

func TestGeneratorWithLogging(t *testing.T) {
	t.Run("generator with default logging", func(t *testing.T) {
		gen := New()
		if gen.logger == nil {
			t.Error("Expected logger to be initialized")
		}
		
		if gen.errorHandler == nil {
			t.Error("Expected error handler to be initialized")
		}
		
		if gen.logger.level != LogLevelInfo {
			t.Errorf("Expected default log level %v, got %v", LogLevelInfo, gen.logger.level)
		}
	})
	
	t.Run("generator with custom logging", func(t *testing.T) {
		gen := NewWithLogger(LogLevelDebug)
		if gen.logger == nil {
			t.Error("Expected logger to be initialized")
		}
		
		if gen.logger.level != LogLevelDebug {
			t.Errorf("Expected custom log level %v, got %v", LogLevelDebug, gen.logger.level)
		}
	})
	
	t.Run("logger methods", func(t *testing.T) {
		gen := NewWithLogger(LogLevelDebug)
		
		// These should not panic
		gen.logger.Debug("debug message")
		gen.logger.Info("info message")
		gen.logger.Warn("warn message")
		gen.logger.Error("error message")
		gen.logger.Success("success message")
		gen.logger.Progress("progress message")
		gen.logger.Warning("warning message")
		
		gen.logger.Step(1, 3, "step message")
		gen.logger.ErrorWithDetails("error with details", nil, nil)
		
		duration := gen.logger.Duration()
		if duration < 0 {
			t.Error("Expected duration to be non-negative")
		}
	})
}