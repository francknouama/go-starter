// Package logger provides CloudWatch-optimized logging for AWS Lambda
package logger

import (
	"fmt"
	"io"
	"os"
	"log/slog"
)

// Logger defines the logging interface optimized for AWS Lambda/CloudWatch
type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
	
	DebugWith(msg string, fields Fields)
	InfoWith(msg string, fields Fields)
	WarnWith(msg string, fields Fields)
	ErrorWith(msg string, fields Fields)
}

// Fields represents structured logging fields
type Fields map[string]interface{}

// Factory creates logger instances optimized for Lambda
type Factory struct{}

// NewFactory creates a new logger factory
func NewFactory() *Factory {
	return &Factory{}
}

// CreateFromProjectConfig creates a CloudWatch-optimized logger
func (f *Factory) CreateFromProjectConfig(loggerType, level, format string, structured bool) (Logger, error) {
	// Lambda always outputs to stdout/stderr, CloudWatch captures it
	output := os.Stdout
	return f.createSlogLogger(level, output)
}

type slogLogger struct {
	logger *slog.Logger
}

func (f *Factory) createSlogLogger(level string, output io.Writer) (Logger, error) {
	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}
	
	// Always use JSON for CloudWatch
	handler := slog.NewJSONHandler(output, opts)
	
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
		return slog.LevelInfo
	}
}

// NewSimpleLogger creates a simple fallback logger
func NewSimpleLogger() Logger {
	return &simpleLogger{}
}

type simpleLogger struct{}

func (s *simpleLogger) Debug(msg string) { fmt.Printf(`{"level":"debug","msg":"%s"}`+"\n", msg) }
func (s *simpleLogger) Info(msg string)  { fmt.Printf(`{"level":"info","msg":"%s"}`+"\n", msg) }
func (s *simpleLogger) Warn(msg string)  { fmt.Printf(`{"level":"warn","msg":"%s"}`+"\n", msg) }
func (s *simpleLogger) Error(msg string) { fmt.Printf(`{"level":"error","msg":"%s"}`+"\n", msg) }

func (s *simpleLogger) DebugWith(msg string, fields Fields) { s.Debug(fmt.Sprintf("%s %+v", msg, fields)) }
func (s *simpleLogger) InfoWith(msg string, fields Fields)  { s.Info(fmt.Sprintf("%s %+v", msg, fields)) }
func (s *simpleLogger) WarnWith(msg string, fields Fields)  { s.Warn(fmt.Sprintf("%s %+v", msg, fields)) }
func (s *simpleLogger) ErrorWith(msg string, fields Fields) { s.Error(fmt.Sprintf("%s %+v", msg, fields)) }