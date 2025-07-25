package logger

import (
	"context"
	"io"
)

// Logger defines a common interface for all supported logger types
type Logger interface {
	// Log level methods
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})

	// Structured logging with key-value pairs
	DebugWith(msg string, fields Fields)
	InfoWith(msg string, fields Fields)
	WarnWith(msg string, fields Fields)
	ErrorWith(msg string, fields Fields)

	// Context-aware logging
	WithContext(ctx context.Context) Logger
	WithFields(fields Fields) Logger

	// Configuration
	SetLevel(level Level)
	SetOutput(w io.Writer)
	DisableColor()

	// Cleanup
	Sync() error
}

// Fields represents structured logging fields
type Fields map[string]interface{}

// Level represents log levels
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// String returns the string representation of the log level
func (l Level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	default:
		return "info"
	}
}

// ParseLevel parses a string level to Level type
func ParseLevel(level string) Level {
	switch level {
	case "debug":
		return DebugLevel
	case "info":
		return InfoLevel
	case "warn":
		return WarnLevel
	case "error":
		return ErrorLevel
	default:
		return InfoLevel
	}
}

// Format represents log output formats
type Format int

const (
	JSONFormat Format = iota
	TextFormat
	ConsoleFormat
)

// String returns the string representation of the log format
func (f Format) String() string {
	switch f {
	case JSONFormat:
		return "json"
	case TextFormat:
		return "text"
	case ConsoleFormat:
		return "console"
	default:
		return "json"
	}
}

// ParseFormat parses a string format to Format type
func ParseFormat(format string) Format {
	switch format {
	case "json":
		return JSONFormat
	case "text":
		return TextFormat
	case "console":
		return ConsoleFormat
	default:
		return JSONFormat
	}
}

// Config represents logger configuration
type Config struct {
	Type       string `yaml:"type" json:"type"`             // slog, zap, logrus, zerolog
	Level      string `yaml:"level" json:"level"`           // debug, info, warn, error
	Format     string `yaml:"format" json:"format"`         // json, text, console
	Structured bool   `yaml:"structured" json:"structured"` // enable structured logging
	Output     string `yaml:"output" json:"output"`         // stdout, stderr, file path
}

// DefaultConfig returns a default logger configuration
func DefaultConfig() *Config {
	return &Config{
		Type:       "slog",
		Level:      "info",
		Format:     "json",
		Structured: true,
		Output:     "stdout",
	}
}
