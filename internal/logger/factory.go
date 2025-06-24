package logger

import (
	"fmt"
	"io"
	"os"
)

// Factory creates logger instances based on configuration
type Factory struct{}

// New creates a new Factory instance
func NewFactory() *Factory {
	return &Factory{}
}

// Create creates a logger instance based on the provided configuration
func (f *Factory) Create(config *Config) (Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Get output writer
	output, err := f.getOutput(config.Output)
	if err != nil {
		return nil, fmt.Errorf("failed to get output writer: %w", err)
	}

	// Parse level and format
	level := ParseLevel(config.Level)
	format := ParseFormat(config.Format)

	// Create logger based on type
	switch config.Type {
	case "slog", "":
		return f.createSlogLogger(level, format, output, config.Structured)
	case "zap":
		return f.createZapLogger(level, format, output, config.Structured)
	case "logrus":
		return f.createLogrusLogger(level, format, output, config.Structured)
	case "zerolog":
		return f.createZerologLogger(level, format, output, config.Structured)
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", config.Type)
	}
}

// CreateDefault creates a logger with default slog configuration
func (f *Factory) CreateDefault() Logger {
	logger, _ := f.Create(DefaultConfig())
	return logger
}

// CreateFromProjectConfig creates a logger from project configuration
func (f *Factory) CreateFromProjectConfig(loggerType, level, format string, structured bool) (Logger, error) {
	config := &Config{
		Type:       loggerType,
		Level:      level,
		Format:     format,
		Structured: structured,
		Output:     "stdout",
	}
	
	return f.Create(config)
}

// getOutput returns the appropriate io.Writer based on output configuration
func (f *Factory) getOutput(output string) (io.Writer, error) {
	switch output {
	case "stdout", "":
		return os.Stdout, nil
	case "stderr":
		return os.Stderr, nil
	default:
		// Assume it's a file path
		file, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", output, err)
		}
		return file, nil
	}
}

// createSlogLogger creates an slog-based logger
func (f *Factory) createSlogLogger(level Level, format Format, output io.Writer, structured bool) (Logger, error) {
	return NewSlogLogger(level, format, output, structured)
}

// createZapLogger creates a zap-based logger (placeholder for Phase 3.1)
func (f *Factory) createZapLogger(level Level, format Format, output io.Writer, structured bool) (Logger, error) {
	// TODO: Implement in Phase 3.1
	return nil, fmt.Errorf("zap logger implementation pending (Phase 3.1)")
}

// createLogrusLogger creates a logrus-based logger (placeholder for Phase 3.1)
func (f *Factory) createLogrusLogger(level Level, format Format, output io.Writer, structured bool) (Logger, error) {
	// TODO: Implement in Phase 3.1
	return nil, fmt.Errorf("logrus logger implementation pending (Phase 3.1)")
}

// createZerologLogger creates a zerolog-based logger (placeholder for Phase 3.1)
func (f *Factory) createZerologLogger(level Level, format Format, output io.Writer, structured bool) (Logger, error) {
	// TODO: Implement in Phase 3.1
	return nil, fmt.Errorf("zerolog logger implementation pending (Phase 3.1)")
}

// GetSupportedTypes returns a list of supported logger types
func GetSupportedTypes() []string {
	return []string{"slog", "zap", "logrus", "zerolog"}
}

// GetSupportedLevels returns a list of supported log levels
func GetSupportedLevels() []string {
	return []string{"debug", "info", "warn", "error"}
}

// GetSupportedFormats returns a list of supported log formats
func GetSupportedFormats() []string {
	return []string{"json", "text", "console"}
}