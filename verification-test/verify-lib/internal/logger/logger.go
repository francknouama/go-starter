// Package logger provides minimal logging functionality for the verify_lib library.
// This is kept simple to avoid forcing logging dependencies on library users.
package logger

import (
	"io"
	"os"
	"log/slog"
)

// Logger defines a minimal logging interface for library use
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	
	DebugWith(msg string, fields Fields)
	InfoWith(msg string, fields Fields)
	WarnWith(msg string, fields Fields)
	ErrorWith(msg string, fields Fields)
	
	Sync() error
}

// Fields represents structured logging fields
type Fields map[string]interface{}

// Factory creates logger instances
type Factory struct{}

// NewFactory creates a new logger factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateFromProjectConfig creates a logger from project configuration
func (f *Factory) CreateFromProjectConfig(loggerType, level, format string, structured bool) (Logger, error) {
	output := os.Stderr // Libraries should log to stderr by default
	return f.createSlogLogger(level, format, output)
}

// slogLogger implements Logger using Go's built-in slog
type slogLogger struct {
	logger *slog.Logger
}

func (f *Factory) createSlogLogger(level, format string, output io.Writer) (Logger, error) {
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}
	
	var handler slog.Handler
	if format == "json" {
		handler = slog.NewJSONHandler(output, opts)
	} else {
		handler = slog.NewTextHandler(output, opts)
	}
	
	return &slogLogger{
		logger: slog.New(handler),
	}, nil
}

func (s *slogLogger) Debug(msg string) { s.logger.Debug(msg) }
func (s *slogLogger) Info(msg string)  { s.logger.Info(msg) }
func (s *slogLogger) Warn(msg string)  { s.logger.Warn(msg) }
func (s *slogLogger) Error(msg string) { s.logger.Error(msg) }

func (s *slogLogger) DebugWith(msg string, fields Fields) {
	s.logger.Debug(msg, fieldsToSlogArgs(fields)...)
}
func (s *slogLogger) InfoWith(msg string, fields Fields) {
	s.logger.Info(msg, fieldsToSlogArgs(fields)...)
}
func (s *slogLogger) WarnWith(msg string, fields Fields) {
	s.logger.Warn(msg, fieldsToSlogArgs(fields)...)
}
func (s *slogLogger) ErrorWith(msg string, fields Fields) {
	s.logger.Error(msg, fieldsToSlogArgs(fields)...)
}

func (s *slogLogger) Sync() error { return nil }

func fieldsToSlogArgs(fields Fields) []interface{} {
	args := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		args = append(args, k, v)
	}
	return args
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelWarn
	}
}