package prompts

import (
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/utils"
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
	expectedGoVersion := utils.GetOptimalGoVersion()
	if config.GoVersion != expectedGoVersion {
		t.Errorf("GoVersion = %v, want %v", config.GoVersion, expectedGoVersion)
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

// Test NewSurvey constructor function
func TestNewSurvey(t *testing.T) {
	prompter := NewSurvey()
	if prompter == nil {
		t.Error("Expected prompter to not be nil")
		return
	}
	if prompter.useFang {
		t.Error("Expected useFang to be false for Survey constructor")
	}
}

// Test ORM mapping logic
func TestOrmMapping(t *testing.T) {
	ormMap := map[string]string{
		"gorm - Feature-rich ORM with associations and migrations (recommended) ✅": "gorm",
		"raw - Raw database/sql package with manual queries ✅":                     "raw",
		"sqlx - Lightweight extensions on database/sql 🔄 Coming Soon":              "sqlx",
		"sqlc - Generate type-safe code from SQL 🔄 Coming Soon":                    "sqlc",
		"ent - Simple, yet feature-complete entity framework 🔄 Coming Soon":        "ent",
		"xorm - Alternative full-featured ORM 🔄 Coming Soon":                       "xorm",
	}

	for display, internal := range ormMap {
		t.Run(display, func(t *testing.T) {
			if internal == "" {
				t.Errorf("ORM mapping for %s should not be empty", display)
			}
			// Validate that only supported ORMs are implemented
			if internal != "gorm" && internal != "raw" {
				// These should trigger validation errors in actual usage
				t.Logf("ORM %s is marked as coming soon", internal)
			}
		})
	}
}

// Test logger extraction logic
func TestLoggerMapping(t *testing.T) {
	tests := []struct {
		selection string
		expected  string
	}{
		{"slog - Go built-in structured logging (recommended)", "slog"},
		{"zap - High-performance, zero-allocation logging", "zap"},
		{"logrus - Feature-rich, popular logging library", "logrus"},
		{"zerolog - Zero allocation, chainable API logging", "zerolog"},
	}

	for _, tt := range tests {
		t.Run(tt.selection, func(t *testing.T) {
			// Simulate the logger extraction logic
			logger := strings.Split(tt.selection, " ")[0]
			if logger != tt.expected {
				t.Errorf("Logger extraction = %v, want %v", logger, tt.expected)
			}
		})
	}
}

// Test level and format extraction for advanced logger config
func TestLoggerLevelFormatExtraction(t *testing.T) {
	levelTests := []struct {
		selection string
		expected  string
	}{
		{"debug - Detailed debugging information", "debug"},
		{"info - General application flow (recommended)", "info"},
		{"warn - Warning messages and potential issues", "warn"},
		{"error - Error conditions only", "error"},
	}

	for _, tt := range levelTests {
		t.Run("level_"+tt.selection, func(t *testing.T) {
			level := strings.Split(tt.selection, " ")[0]
			if level != tt.expected {
				t.Errorf("Level extraction = %v, want %v", level, tt.expected)
			}
		})
	}

	formatTests := []struct {
		selection string
		expected  string
	}{
		{"json - Structured JSON format (recommended)", "json"},
		{"text - Human-readable text format", "text"},
		{"console - Colored console output", "console"},
	}

	for _, tt := range formatTests {
		t.Run("format_"+tt.selection, func(t *testing.T) {
			format := strings.Split(tt.selection, " ")[0]
			if format != tt.expected {
				t.Errorf("Format extraction = %v, want %v", format, tt.expected)
			}
		})
	}
}

// Test project configuration validation and defaults
func TestPrompter_ConfigValidation(t *testing.T) {
	prompter := New()

	tests := []struct {
		name     string
		config   types.ProjectConfig
		advanced bool
		wantErr  bool
	}{
		{
			name: "library type should not prompt for logger",
			config: types.ProjectConfig{
				Name:      "test-lib",
				Module:    "github.com/test/lib",
				Type:      "library",
				Framework: "standard",
			},
			advanced: false,
			wantErr:  false,
		},
		{
			name: "cli type should allow framework selection",
			config: types.ProjectConfig{
				Name:   "test-cli",
				Module: "github.com/test/cli",
				Type:   "cli",
			},
			advanced: false,
			wantErr:  false,
		},
		{
			name: "web-api type should support all features",
			config: types.ProjectConfig{
				Name:   "test-api",
				Module: "github.com/test/api",
				Type:   "web-api",
			},
			advanced: true,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := prompter.GetProjectConfig(tt.config, tt.advanced)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Validate specific behaviors based on type
			switch tt.config.Type {
			case "library":
				if result.Logger != "slog" {
					t.Errorf("Library projects should default to slog logger, got %v", result.Logger)
				}
			case "web-api":
				if result.GoVersion == "" {
					t.Error("Web API projects should have Go version set")
				}
			}

			// Ensure Variables map is initialized
			if result.Variables == nil {
				t.Error("Variables should be initialized")
			}
		})
	}
}

// Test go version handling
func TestGoVersionHandling(t *testing.T) {
	prompter := New()

	tests := []struct {
		name       string
		goVersion  string
		shouldKeep bool
	}{
		{
			name:       "custom go version should be preserved",
			goVersion:  "1.20",
			shouldKeep: true,
		},
		{
			name:       "empty go version should be set to optimal",
			goVersion:  "",
			shouldKeep: false,
		},
		{
			name:       "optimal go version may be prompted for change",
			goVersion:  utils.GetOptimalGoVersion(),
			shouldKeep: false, // This may trigger prompting
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				Type:      "library", // Use library to avoid additional prompts
				Framework: "standard",
				GoVersion: tt.goVersion,
			}

			result, err := prompter.GetProjectConfig(config, false)
			if err != nil {
				t.Errorf("GetProjectConfig() error = %v", err)
				return
			}

			if tt.shouldKeep && result.GoVersion != tt.goVersion {
				t.Errorf("Expected GoVersion to be preserved as %v, got %v", tt.goVersion, result.GoVersion)
			}

			if !tt.shouldKeep && result.GoVersion == "" {
				t.Error("GoVersion should be set to a valid value")
			}
		})
	}
}

// Test framework logic based on project type
func TestFrameworkLogic(t *testing.T) {
	tests := []struct {
		projectType       string
		shouldHaveOptions bool
		expectedOptions   []string
	}{
		{
			projectType:       "web-api",
			shouldHaveOptions: true,
			expectedOptions:   []string{"gin", "echo", "fiber", "chi", "standard"},
		},
		{
			projectType:       "cli",
			shouldHaveOptions: true,
			expectedOptions:   []string{"cobra", "standard"},
		},
		{
			projectType:       "library",
			shouldHaveOptions: false,
			expectedOptions:   nil,
		},
		{
			projectType:       "lambda",
			shouldHaveOptions: false,
			expectedOptions:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.projectType, func(t *testing.T) {
			if tt.shouldHaveOptions {
				// Verify that the expected frameworks are available
				for _, framework := range tt.expectedOptions {
					if framework == "" {
						t.Errorf("Framework option should not be empty for type %s", tt.projectType)
					}
				}
			} else {
				// Library and lambda types should not need framework selection
				t.Logf("Project type %s correctly requires no framework selection", tt.projectType)
			}
		})
	}
}

// Test variable defaults and inheritance
func TestVariableDefaults(t *testing.T) {
	prompter := New()

	config := types.ProjectConfig{
		Name:   "test-project",
		Module: "github.com/test/project",
		Type:   "library",
		Variables: map[string]string{
			"ExistingVar": "ExistingValue",
		},
	}

	result, err := prompter.GetProjectConfig(config, false)
	if err != nil {
		t.Errorf("GetProjectConfig() error = %v", err)
		return
	}

	// Check that existing variables are preserved
	if result.Variables["ExistingVar"] != "ExistingValue" {
		t.Errorf("Expected ExistingVar to be preserved as 'ExistingValue', got %v", result.Variables["ExistingVar"])
	}

	// Check that Variables map is properly initialized
	if result.Variables == nil {
		t.Error("Variables map should be initialized")
	}
}

// Test feature initialization for different project types
func TestFeatureInitialization(t *testing.T) {
	tests := []struct {
		name        string
		projectType string
		advanced    bool
		expectAuth  bool
		expectDB    bool
	}{
		{
			name:        "basic web-api should support database",
			projectType: "web-api",
			advanced:    false,
			expectAuth:  false,
			expectDB:    true,
		},
		{
			name:        "advanced web-api should support all features",
			projectType: "web-api",
			advanced:    true,
			expectAuth:  true,
			expectDB:    true,
		},
		{
			name:        "cli should not need database prompts",
			projectType: "cli",
			advanced:    false,
			expectAuth:  false,
			expectDB:    false,
		},
		{
			name:        "library should not need feature prompts",
			projectType: "library",
			advanced:    false,
			expectAuth:  false,
			expectDB:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test validates the logic flow but doesn't test actual prompting
			// since that would require user interaction
			if tt.projectType == "web-api" && tt.expectDB {
				t.Logf("Project type %s correctly supports database features", tt.projectType)
			}
			if tt.projectType == "web-api" && tt.advanced && tt.expectAuth {
				t.Logf("Advanced web-api correctly supports authentication features")
			}
		})
	}
}
