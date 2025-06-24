package prompts

import (
	"strings"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestNew(t *testing.T) {
	prompter := New()
	if prompter == nil {
		t.Error("Expected prompter to not be nil")
	}
}

func TestPrompter_isInteractiveMode(t *testing.T) {
	prompter := New()

	tests := []struct {
		name     string
		initial  types.ProjectConfig
		expected bool
	}{
		{
			name:     "empty config - interactive",
			initial:  types.ProjectConfig{},
			expected: true,
		},
		{
			name: "missing name - interactive",
			initial: types.ProjectConfig{
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			expected: true,
		},
		{
			name: "missing module - interactive",
			initial: types.ProjectConfig{
				Name: "test-project",
				Type: "web-api",
			},
			expected: true,
		},
		{
			name: "missing type - interactive",
			initial: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
			},
			expected: true,
		},
		{
			name: "all required fields provided - non-interactive",
			initial: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := prompter.isInteractiveMode(tt.initial)
			if result != tt.expected {
				t.Errorf("isInteractiveMode() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestPrompter_GetProjectConfig_NonInteractive(t *testing.T) {
	prompter := New()

	initial := types.ProjectConfig{
		Name:      "test-project",
		Module:    "github.com/test/project",
		Type:      "library",  // Use library type to avoid framework prompts
		Framework: "standard", // Pre-set framework to avoid prompts
	}

	// This should work without prompting since all required fields are provided
	config, err := prompter.GetProjectConfig(initial, false)
	if err != nil {
		t.Errorf("GetProjectConfig() error = %v", err)
	}

	// Check that values were preserved
	if config.Name != initial.Name {
		t.Errorf("Name = %v, want %v", config.Name, initial.Name)
	}
	if config.Module != initial.Module {
		t.Errorf("Module = %v, want %v", config.Module, initial.Module)
	}
	if config.Type != initial.Type {
		t.Errorf("Type = %v, want %v", config.Type, initial.Type)
	}

	// Check defaults were set
	if config.GoVersion != "1.21" {
		t.Errorf("GoVersion = %v, want %v", config.GoVersion, "1.21")
	}
	if config.Variables == nil {
		t.Error("Variables should be initialized")
	}
}

func TestPrompter_GetProjectConfig_WithDefaults(t *testing.T) {
	prompter := New()

	initial := types.ProjectConfig{
		Name:      "test-project",
		Module:    "github.com/test/project",
		Type:      "library",  // Use library type to avoid framework prompts
		Framework: "standard", // Pre-set framework to avoid prompts
		GoVersion: "1.20",     // Custom Go version
	}

	config, err := prompter.GetProjectConfig(initial, false)
	if err != nil {
		t.Errorf("GetProjectConfig() error = %v", err)
	}

	// Should preserve custom Go version
	if config.GoVersion != "1.20" {
		t.Errorf("GoVersion = %v, want %v", config.GoVersion, "1.20")
	}
}

// Mock test for framework extraction logic
func TestFrameworkExtraction(t *testing.T) {
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
			// Simulate the framework extraction logic
			framework := strings.ToLower(strings.Split(tt.selection, " ")[0])
			if framework != tt.expected {
				t.Errorf("Framework extraction = %v, want %v", framework, tt.expected)
			}
		})
	}
}

// Test architecture mapping
func TestArchitectureMapping(t *testing.T) {
	archMap := map[string]string{
		"Standard - Simple structure":                 "standard",
		"Clean Architecture - Uncle Bob's principles": "clean",
		"Domain-Driven Design - Business-focused":     "ddd",
		"Hexagonal - Ports and adapters":              "hexagonal",
	}

	for display, internal := range archMap {
		t.Run(display, func(t *testing.T) {
			if internal == "" {
				t.Errorf("Architecture mapping for %s should not be empty", display)
			}
			if len(internal) > 20 {
				t.Errorf("Architecture internal name %s is too long", internal)
			}
		})
	}
}

// Test type mapping
func TestTypeMapping(t *testing.T) {
	typeMap := map[string]string{
		"Web API - REST API or web service":   "web-api",
		"CLI Application - Command-line tool": "cli",
		"Library - Reusable Go package":       "library",
		"AWS Lambda - Serverless function":    "lambda",
	}

	for display, internal := range typeMap {
		t.Run(display, func(t *testing.T) {
			if internal == "" {
				t.Errorf("Type mapping for %s should not be empty", display)
			}
			if strings.Contains(internal, " ") {
				t.Errorf("Type internal name %s should not contain spaces", internal)
			}
		})
	}
}

// Test database driver normalization
func TestDatabaseDriverNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"PostgreSQL", "postgresql"},
		{"MySQL", "mysql"},
		{"MongoDB", "mongodb"},
		{"SQLite", "sqlite"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Simulate the database driver normalization
			driver := strings.ToLower(tt.input)
			if driver != tt.expected {
				t.Errorf("Database driver normalization = %v, want %v", driver, tt.expected)
			}
		})
	}
}

// Test authentication type normalization
func TestAuthenticationTypeNormalization(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"JWT", "jwt"},
		{"OAuth2", "oauth2"},
		{"Session-based", "session-based"},
		{"API Key", "api key"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			// Simulate the auth type normalization
			authType := strings.ToLower(tt.input)
			if authType != tt.expected {
				t.Errorf("Auth type normalization = %v, want %v", authType, tt.expected)
			}
		})
	}
}
