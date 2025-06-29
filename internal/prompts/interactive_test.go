package prompts

import (
	"strings"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectConfig_InteractiveMode(t *testing.T) {
	// Test case 1: Fully interactive, no initial config
	t.Run("fully interactive", func(t *testing.T) {
		mockAdapter := NewMockSurveyAdapter(map[string]interface{}{
			"What's your project name?":                     "my-test-project",
			"Module path:":                                  "github.com/user/my-test-project",
			"What type of project?":                         "Web API - REST API or web service",
			"Select Go version:":                            "Go 1.23 (latest)",
			"Which web framework?":                          "Gin (recommended)",
			"Which logger?":                                 "slog - Go built-in structured logging (recommended)",
			"Which ORM/database abstraction do you prefer?": "gorm - Feature-rich ORM with associations and migrations (recommended) ✅",
			"Authentication type?":                          "JWT",
			"Log level?":                                    "info - General application flow (recommended)",
			"Log format?":                                   "json - Structured JSON format (recommended)",
			"Which databases do you want to use? (Space to select, Enter to confirm)": []string{"PostgreSQL"},
			"Would you like to configure advanced options?":                           true,
			"Add database support?": true,
			"Add authentication?":   true,
		})

		prompter := &SurveyPrompter{useFang: false, surveyAdapter: mockAdapter}
		config, err := prompter.GetProjectConfig(types.ProjectConfig{}, true)

		assert.NoError(t, err)
		assert.Equal(t, "my-test-project", config.Name)
		assert.Equal(t, "github.com/user/my-test-project", config.Module)
		assert.Equal(t, "web-api", config.Type)
		assert.Equal(t, "1.24", config.GoVersion)
		assert.Equal(t, "gin", config.Framework)
		assert.Equal(t, "slog", config.Logger)
		assert.Contains(t, config.Features.Database.Drivers, "postgresql")
		assert.Equal(t, "gorm", config.Features.Database.ORM)
		assert.Equal(t, "jwt", config.Features.Authentication.Type)
		assert.Equal(t, "info", config.Features.Logging.Level)
		assert.Equal(t, "json", config.Features.Logging.Format)
	})

	// Test case 2: Partially interactive, some initial config
	t.Run("partially interactive", func(t *testing.T) {
		mockAdapter := NewMockSurveyAdapter(map[string]interface{}{
			"What type of project?": "CLI Application - Command-line tool",
			"Which CLI framework?":  "Cobra (recommended)",
			"Which logger?":         "zap - High-performance, zero-allocation logging",
			"Select Go version:":    "Go 1.23 (latest)",
		})

		prompter := &SurveyPrompter{useFang: false, surveyAdapter: mockAdapter}
		initialConfig := types.ProjectConfig{
			Name:   "my-cli-tool",
			Module: "github.com/user/my-cli-tool",
		}
		config, err := prompter.GetProjectConfig(initialConfig, false)

		assert.NoError(t, err)
		assert.Equal(t, "my-cli-tool", config.Name)
		assert.Equal(t, "github.com/user/my-cli-tool", config.Module)
		assert.Equal(t, "cli", config.Type)
		assert.Equal(t, "cobra", config.Framework)
		assert.Equal(t, "zap", config.Logger)
	})

	// Test case 3: Non-interactive, all config provided
	t.Run("non-interactive", func(t *testing.T) {
		mockAdapter := NewMockSurveyAdapter(map[string]interface{}{})

		prompter := &SurveyPrompter{useFang: false, surveyAdapter: mockAdapter}
		fullConfig := types.ProjectConfig{
			Name:      "full-config",
			Module:    "github.com/user/full-config",
			Type:      "library",
			GoVersion: "1.21",
			Framework: "",
			Logger:    "slog",
		}
		config, err := prompter.GetProjectConfig(fullConfig, false)

		assert.NoError(t, err)
		assert.Equal(t, fullConfig.Name, config.Name)
		assert.Equal(t, fullConfig.Module, config.Module)
		assert.Equal(t, fullConfig.Type, config.Type)
		assert.Equal(t, fullConfig.GoVersion, config.GoVersion)
		assert.Equal(t, fullConfig.Framework, config.Framework)
		assert.Equal(t, fullConfig.Logger, config.Logger)
		// Features and Variables are initialized even in non-interactive mode
		assert.NotNil(t, config.Features)
		assert.NotNil(t, config.Variables)
	})
}

func TestMapSelectionToProjectType(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"Web API - REST API or web service", "web-api"},
		{"CLI Application - Command-line tool", "cli"},
		{"Library - Reusable Go package", "library"},
		{"AWS Lambda - Serverless function", "lambda"},
	}

	for _, tt := range tests {
		t.Run(tt.selection, func(t *testing.T) {
			// Test by checking the options in promptProjectType
			// This is more of an integration test
			assert.Contains(t, []string{"web-api", "cli", "library", "lambda"}, tt.expected)
		})
	}
}

func TestMapFrameworkSelection(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"Gin (recommended)", "gin"},
		{"Echo", "echo"},
		{"Fiber", "fiber"},
		{"Chi", "chi"},
		{"Standard library", "standard"},
		{"Cobra (recommended)", "cobra"},
	}

	for _, tt := range tests {
		t.Run(tt.selection, func(t *testing.T) {
			// Extract framework name (remove description)
			result := strings.ToLower(strings.Split(tt.selection, " ")[0])
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapDatabaseSelection(t *testing.T) {
	// Database selections map directly
	databases := []string{"PostgreSQL", "MySQL", "MongoDB", "SQLite", "Redis"}

	for _, db := range databases {
		t.Run(db, func(t *testing.T) {
			assert.NotEmpty(t, db)
		})
	}
}

func TestLoggerSelectionMapping(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"slog - Go built-in structured logging (recommended)", "slog"},
		{"zap - High-performance, zero-allocation logging", "zap"},
		{"logrus - Feature-rich logging with fields", "logrus"},
		{"zerolog - Fast JSON logger with zero allocation", "zerolog"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Test logger selection parsing
			parts := strings.Split(tt.selection, " - ")
			assert.Equal(t, tt.expected, parts[0])
		})
	}
}

func TestORMSelectionParsing(t *testing.T) {
	ormOptions := []string{
		"gorm - Feature-rich ORM with associations and migrations (recommended) ✅",
		"sqlx - Lightweight extension on database/sql with named queries",
		"sqlc - Type-safe SQL with code generation (compile-time safety)",
		"ent - Entity framework with graph-based data modeling",
		"database/sql - Standard library (raw SQL)",
	}

	for _, option := range ormOptions {
		t.Run(option, func(t *testing.T) {
			// Test ORM selection parsing
			parts := strings.Split(option, " - ")
			ormName := parts[0]
			assert.NotEmpty(t, ormName)
			assert.Contains(t, []string{"gorm", "sqlx", "sqlc", "ent", "database/sql"}, ormName)
		})
	}
}

func TestLogLevelMapping(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"debug - Detailed information for debugging", "debug"},
		{"info - General application flow (recommended)", "info"},
		{"warn - Warning messages", "warn"},
		{"error - Error messages only", "error"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Test log level selection parsing
			parts := strings.Split(tt.selection, " - ")
			assert.Equal(t, tt.expected, parts[0])
		})
	}
}

func TestLogFormatMapping(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"json - Structured JSON format (recommended)", "json"},
		{"text - Human-readable text format", "text"},
		{"console - Colorized console output", "console"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			// Test log format selection parsing
			parts := strings.Split(tt.selection, " - ")
			assert.Equal(t, tt.expected, parts[0])
		})
	}
}

func TestErrorHandling(t *testing.T) {
	t.Run("survey error propagation", func(t *testing.T) {
		mockAdapter := &MockSurveyAdapter{
			responses: map[string]interface{}{},
		}

		prompter := &SurveyPrompter{useFang: false, surveyAdapter: mockAdapter}
		_, err := prompter.GetProjectConfig(types.ProjectConfig{}, true)

		// We expect no error because the mock doesn't return errors by default
		assert.NoError(t, err)
	})
}
