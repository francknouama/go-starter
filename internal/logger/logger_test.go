package logger

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestFactory_Create(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "default slog logger",
			config: &Config{
				Type:       "slog",
				Level:      "info",
				Format:     "json",
				Structured: true,
				Output:     "stdout",
			},
			wantErr: false,
		},
		{
			name: "empty type defaults to slog",
			config: &Config{
				Type:       "",
				Level:      "debug",
				Format:     "text",
				Structured: true,
				Output:     "stdout",
			},
			wantErr: false,
		},
		{
			name: "unsupported logger type",
			config: &Config{
				Type:       "invalid",
				Level:      "info",
				Format:     "json",
				Structured: true,
				Output:     "stdout",
			},
			wantErr: true,
		},
		{
			name:    "nil config uses defaults",
			config:  nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger, err := factory.Create(tt.config)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			if logger == nil {
				t.Error("expected logger but got nil")
			}
		})
	}
}

func TestSlogLogger_Logging(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewSlogLogger(InfoLevel, JSONFormat, &buf, true)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Test basic logging
	logger.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "test message") {
		t.Error("expected log message not found in output")
	}

	// Reset buffer
	buf.Reset()

	// Test structured logging
	fields := Fields{
		"key1": "value1",
		"key2": 42,
	}
	logger.InfoWith("structured message", fields)
	output = buf.String()
	if !strings.Contains(output, "structured message") {
		t.Error("expected structured message not found in output")
	}
	if !strings.Contains(output, "key1") || !strings.Contains(output, "value1") {
		t.Error("expected structured fields not found in output")
	}
}

func TestSlogLogger_WithFields(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewSlogLogger(InfoLevel, JSONFormat, &buf, true)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// Create logger with fields
	fieldsLogger := logger.WithFields(Fields{
		"service": "test-service",
		"version": "1.0.0",
	})

	fieldsLogger.Info("test message")
	output := buf.String()
	if !strings.Contains(output, "service") || !strings.Contains(output, "test-service") {
		t.Error("expected service field not found in output")
	}
}

func TestSlogLogger_WithContext(t *testing.T) {
	var buf bytes.Buffer
	logger, err := NewSlogLogger(InfoLevel, JSONFormat, &buf, true)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	ctx := context.Background()
	contextLogger := logger.WithContext(ctx)
	if contextLogger == nil {
		t.Error("expected context logger but got nil")
	}

	contextLogger.Info("context message")
	output := buf.String()
	if !strings.Contains(output, "context message") {
		t.Error("expected context message not found in output")
	}
}

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level Level
		want  string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{Level(999), "info"}, // Invalid level defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.level.String(); got != tt.want {
				t.Errorf("Level.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input string
		want  Level
	}{
		{"debug", DebugLevel},
		{"info", InfoLevel},
		{"warn", WarnLevel},
		{"error", ErrorLevel},
		{"invalid", InfoLevel}, // Invalid input defaults to info
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseLevel(tt.input); got != tt.want {
				t.Errorf("ParseLevel(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestFormat_String(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{JSONFormat, "json"},
		{TextFormat, "text"},
		{ConsoleFormat, "console"},
		{Format(999), "json"}, // Invalid format defaults to json
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.format.String(); got != tt.want {
				t.Errorf("Format.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFormat(t *testing.T) {
	tests := []struct {
		input string
		want  Format
	}{
		{"json", JSONFormat},
		{"text", TextFormat},
		{"console", ConsoleFormat},
		{"invalid", JSONFormat}, // Invalid input defaults to json
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if got := ParseFormat(tt.input); got != tt.want {
				t.Errorf("ParseFormat(%v) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetSupportedTypes(t *testing.T) {
	types := GetSupportedTypes()
	expected := []string{"slog", "zap", "logrus", "zerolog"}
	
	if len(types) != len(expected) {
		t.Errorf("expected %d types, got %d", len(expected), len(types))
	}
	
	for i, expectedType := range expected {
		if types[i] != expectedType {
			t.Errorf("expected type %s at index %d, got %s", expectedType, i, types[i])
		}
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Type != "slog" {
		t.Errorf("expected default type 'slog', got '%s'", config.Type)
	}
	if config.Level != "info" {
		t.Errorf("expected default level 'info', got '%s'", config.Level)
	}
	if config.Format != "json" {
		t.Errorf("expected default format 'json', got '%s'", config.Format)
	}
	if !config.Structured {
		t.Error("expected default structured to be true")
	}
	if config.Output != "stdout" {
		t.Errorf("expected default output 'stdout', got '%s'", config.Output)
	}
}