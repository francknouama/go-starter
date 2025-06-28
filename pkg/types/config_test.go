package types

import (
	"testing"
)

func TestProjectConfig_BasicFields(t *testing.T) {
	config := ProjectConfig{
		Name:         "test-project",
		Module:       "github.com/user/test-project",
		Type:         "web-api",
		GoVersion:    "1.21",
		Framework:    "gin",
		Architecture: "clean",
		Logger:       "slog",
		Author:       "John Doe",
		Email:        "john@example.com",
		License:      "MIT",
	}

	// Test that all fields are accessible
	if config.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", config.Name)
	}
	
	if config.Module != "github.com/user/test-project" {
		t.Errorf("Expected module 'github.com/user/test-project', got '%s'", config.Module)
	}
	
	if config.Type != "web-api" {
		t.Errorf("Expected type 'web-api', got '%s'", config.Type)
	}
	
	if config.GoVersion != "1.21" {
		t.Errorf("Expected Go version '1.21', got '%s'", config.GoVersion)
	}
	
	if config.Framework != "gin" {
		t.Errorf("Expected framework 'gin', got '%s'", config.Framework)
	}
}

func TestDatabaseConfig_GetDrivers(t *testing.T) {
	testCases := []struct {
		name     string
		config   DatabaseConfig
		expected []string
	}{
		{
			name: "multiple drivers",
			config: DatabaseConfig{
				Drivers: []string{"postgres", "redis", "mongo"},
			},
			expected: []string{"postgres", "redis", "mongo"},
		},
		{
			name: "single driver via Drivers",
			config: DatabaseConfig{
				Drivers: []string{"postgres"},
			},
			expected: []string{"postgres"},
		},
		{
			name: "legacy Driver field only",
			config: DatabaseConfig{
				Driver: "mysql", //nolint:staticcheck
			},
			expected: []string{"mysql"},
		},
		{
			name: "both Drivers and Driver set",
			config: DatabaseConfig{
				Drivers: []string{"postgres", "redis"},
				Driver:  "mysql", //nolint:staticcheck
			},
			expected: []string{"postgres", "redis"},
		},
		{
			name:     "empty config",
			config:   DatabaseConfig{},
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.config.GetDrivers()
			if len(result) != len(tc.expected) {
				t.Errorf("Expected %d drivers, got %d", len(tc.expected), len(result))
			}

			for i, expected := range tc.expected {
				if i >= len(result) || result[i] != expected {
					t.Errorf("Expected driver %d to be '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}

func TestDatabaseConfig_HasDriver(t *testing.T) {
	config := DatabaseConfig{
		Drivers: []string{"postgres", "redis", "mongo"},
	}

	testCases := []struct {
		driver   string
		expected bool
	}{
		{"postgres", true},
		{"redis", true},
		{"mongo", true},
		{"mysql", false},
		{"sqlite", false},
		{"", false},
	}

	for _, tc := range testCases {
		t.Run(tc.driver, func(t *testing.T) {
			result := config.HasDriver(tc.driver)
			if result != tc.expected {
				t.Errorf("HasDriver(%q) = %v, expected %v", tc.driver, result, tc.expected)
			}
		})
	}
}

func TestDatabaseConfig_HasDriver_Legacy(t *testing.T) {
	// Test legacy Driver field support
	config := DatabaseConfig{
		Driver: "mysql", //nolint:staticcheck
	}

	if !config.HasDriver("mysql") {
		t.Error("Should support legacy Driver field")
	}

	if config.HasDriver("postgres") {
		t.Error("Should not match non-existent drivers")
	}
}

func TestDatabaseConfig_Initialization(t *testing.T) {
	config := DatabaseConfig{
		Drivers: []string{"postgres", "redis"},
		ORM:     "gorm",
	}

	drivers := config.GetDrivers()
	if len(drivers) != 2 {
		t.Errorf("Expected 2 drivers, got %d", len(drivers))
	}

	if config.ORM != "gorm" {
		t.Errorf("Expected ORM 'gorm', got '%s'", config.ORM)
	}
}

func TestAuthConfig_BasicFields(t *testing.T) {
	config := AuthConfig{
		Type:      "jwt",
		Providers: []string{"google", "github"},
	}

	if config.Type != "jwt" {
		t.Errorf("Expected type 'jwt', got '%s'", config.Type)
	}

	if len(config.Providers) != 2 {
		t.Errorf("Expected 2 providers, got %d", len(config.Providers))
	}
}

func TestDeployConfig_BasicFields(t *testing.T) {
	config := DeployConfig{
		Targets: []string{"docker", "kubernetes", "lambda"},
	}

	if len(config.Targets) != 3 {
		t.Errorf("Expected 3 targets, got %d", len(config.Targets))
	}

	expectedTargets := []string{"docker", "kubernetes", "lambda"}
	for i, expected := range expectedTargets {
		if i >= len(config.Targets) || config.Targets[i] != expected {
			t.Errorf("Expected target %d to be '%s', got '%s'", i, expected, config.Targets[i])
		}
	}
}

func TestLoggingConfig_BasicFields(t *testing.T) {
	config := LoggingConfig{
		Type:       "slog",
		Level:      "info",
		Format:     "json",
		Structured: true,
	}

	if config.Type != "slog" {
		t.Errorf("Expected type 'slog', got '%s'", config.Type)
	}

	if config.Level != "info" {
		t.Errorf("Expected level 'info', got '%s'", config.Level)
	}

	if config.Format != "json" {
		t.Errorf("Expected format 'json', got '%s'", config.Format)
	}

	if !config.Structured {
		t.Error("Expected Structured to be true")
	}
}

func TestGenerationOptions_BasicFields(t *testing.T) {
	options := GenerationOptions{
		OutputPath: "/tmp/test-project",
		NoGit:      true,
	}

	if options.OutputPath != "/tmp/test-project" {
		t.Errorf("Expected OutputPath '/tmp/test-project', got '%s'", options.OutputPath)
	}

	if !options.NoGit {
		t.Error("Expected NoGit to be true")
	}
}

func TestGenerationResult_BasicFields(t *testing.T) {
	result := GenerationResult{
		ProjectPath:  "/tmp/test-project",
		FilesCreated: []string{"go.mod", "main.go"},
		Success:      true,
		Error:        nil,
	}

	if result.ProjectPath != "/tmp/test-project" {
		t.Errorf("Expected ProjectPath '/tmp/test-project', got '%s'", result.ProjectPath)
	}

	if len(result.FilesCreated) != 2 {
		t.Errorf("Expected 2 files created, got %d", len(result.FilesCreated))
	}

	if !result.Success {
		t.Error("Expected Success to be true")
	}
}

func TestTemplate_BasicFields(t *testing.T) {
	template := Template{
		Name:        "web-api",
		Description: "REST API template",
		Dependencies: []Dependency{
			{Module: "github.com/gin-gonic/gin", Version: "v1.9.0"},
			{Module: "github.com/spf13/cobra", Version: "v1.7.0"},
		},
	}

	if template.Name != "web-api" {
		t.Errorf("Expected Name 'web-api', got '%s'", template.Name)
	}

	if template.Description != "REST API template" {
		t.Errorf("Expected Description 'REST API template', got '%s'", template.Description)
	}

	if len(template.Dependencies) != 2 {
		t.Errorf("Expected 2 dependencies, got %d", len(template.Dependencies))
	}
}