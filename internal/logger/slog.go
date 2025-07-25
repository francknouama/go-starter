package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
)

// SlogLogger wraps Go's built-in slog logger to implement our Logger interface
type SlogLogger struct {
	logger *slog.Logger
	level  Level
	output io.Writer
}

// NewSlogLogger creates a new slog-based logger
func NewSlogLogger(level Level, format Format, output io.Writer, structured bool) (*SlogLogger, error) {
	if output == nil {
		output = os.Stdout
	}

	// Create handler based on format
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: convertToSlogLevel(level),
	}

	switch format {
	case JSONFormat:
		handler = slog.NewJSONHandler(output, opts)
	case TextFormat, ConsoleFormat:
		handler = slog.NewTextHandler(output, opts)
	default:
		handler = slog.NewJSONHandler(output, opts)
	}

	logger := slog.New(handler)

	return &SlogLogger{
		logger: logger,
		level:  level,
		output: output,
	}, nil
}

// Debug logs a debug message
func (s *SlogLogger) Debug(msg string, args ...interface{}) {
	s.logger.Debug(msg, args...)
}

// Info logs an info message
func (s *SlogLogger) Info(msg string, args ...interface{}) {
	s.logger.Info(msg, args...)
}

// Warn logs a warning message
func (s *SlogLogger) Warn(msg string, args ...interface{}) {
	s.logger.Warn(msg, args...)
}

// Error logs an error message
func (s *SlogLogger) Error(msg string, args ...interface{}) {
	s.logger.Error(msg, args...)
}

// DebugWith logs a debug message with structured fields
func (s *SlogLogger) DebugWith(msg string, fields Fields) {
	s.logger.Debug(msg, fieldsToSlogArgs(fields)...)
}

// InfoWith logs an info message with structured fields
func (s *SlogLogger) InfoWith(msg string, fields Fields) {
	s.logger.Info(msg, fieldsToSlogArgs(fields)...)
}

// WarnWith logs a warning message with structured fields
func (s *SlogLogger) WarnWith(msg string, fields Fields) {
	s.logger.Warn(msg, fieldsToSlogArgs(fields)...)
}

// ErrorWith logs an error message with structured fields
func (s *SlogLogger) ErrorWith(msg string, fields Fields) {
	s.logger.Error(msg, fieldsToSlogArgs(fields)...)
}

// WithContext returns a logger with context
func (s *SlogLogger) WithContext(ctx context.Context) Logger {
	return &SlogLogger{
		logger: s.logger.With(),
		level:  s.level,
		output: s.output,
	}
}

// WithFields returns a logger with additional fields
func (s *SlogLogger) WithFields(fields Fields) Logger {
	return &SlogLogger{
		logger: s.logger.With(fieldsToSlogArgs(fields)...),
		level:  s.level,
		output: s.output,
	}
}

// SetLevel sets the log level (note: slog level is set at handler creation)
func (s *SlogLogger) SetLevel(level Level) {
	s.level = level
	// Note: For slog, we'd need to recreate the handler to change level
	// This is a simplified implementation
}

// SetOutput sets the output writer (note: slog output is set at handler creation)
func (s *SlogLogger) SetOutput(w io.Writer) {
	s.output = w
	// Note: For slog, we'd need to recreate the handler to change output
	// This is a simplified implementation
}

// DisableColor disables color output (no-op for slog)
func (s *SlogLogger) DisableColor() {
	// slog doesn't have built-in color support to disable
	// This is a no-op for compatibility with the Logger interface
}

// Sync flushes any buffered log entries (no-op for slog)
func (s *SlogLogger) Sync() error {
	// slog doesn't require explicit syncing
	return nil
}

// Helper functions

// convertToSlogLevel converts our Level type to slog.Level
func convertToSlogLevel(level Level) slog.Level {
	switch level {
	case DebugLevel:
		return slog.LevelDebug
	case InfoLevel:
		return slog.LevelInfo
	case WarnLevel:
		return slog.LevelWarn
	case ErrorLevel:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// fieldsToSlogArgs converts Fields to slog arguments
func fieldsToSlogArgs(fields Fields) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}
	return args
}
