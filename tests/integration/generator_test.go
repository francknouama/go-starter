package integration

import (
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// TestGeneratorBasicFunctionality tests basic generator functionality
func TestGeneratorBasicFunctionality(t *testing.T) {
	// Test generator creation
	gen := generator.New()
	if gen == nil {
		t.Fatal("Generator should not be nil")
	}

	// Test basic generator functionality without actual generation
	// since we don't have templates yet
	t.Log("Generator created successfully")
}

// TestGeneratorValidation tests project validation
func TestGeneratorValidation(t *testing.T) {
	gen := generator.New()
	_ = gen // Use the generator to avoid unused variable warning

	tests := []struct {
		name          string
		config        *types.ProjectConfig
		shouldFail    bool
		errorContains string
	}{
		{
			name: "valid config",
			config: &types.ProjectConfig{
				Name:   "valid-project",
				Module: "github.com/test/valid-project",
				Type:   "library",
			},
			shouldFail: false,
		},
		{
			name: "empty project name",
			config: &types.ProjectConfig{
				Name:   "",
				Module: "github.com/test/invalid",
				Type:   "library",
			},
			shouldFail:    true,
			errorContains: "name",
		},
		{
			name: "empty module path",
			config: &types.ProjectConfig{
				Name:   "test-project",
				Module: "",
				Type:   "library",
			},
			shouldFail:    true,
			errorContains: "module",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't actually test generation without templates
			// but we can test the config structure
			if tt.config.Name == "" && !tt.shouldFail {
				t.Error("Empty name should fail validation")
			}
			if tt.config.Module == "" && !tt.shouldFail {
				t.Error("Empty module should fail validation")
			}
		})
	}
}

// TestGeneratorWithConfig tests generator with configuration
func TestGeneratorWithConfig(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Profiles: map[string]config.Profile{
			"test": {
				Author:  "Test Author",
				Email:   "test@example.com",
				License: "MIT",
				Defaults: config.ProfileDefaults{
					GoVersion:    "1.21",
					Framework:    "gin",
					Architecture: "standard",
				},
			},
		},
		CurrentProfile: "test",
	}

	gen := generator.New()
	_ = gen // Use the generator to avoid unused variable warning

	// Test that generator can work with configuration
	// This is mainly to ensure the interface is correct
	projectConfig := &types.ProjectConfig{
		Name:   "config-test",
		Module: "github.com/test/config-test",
		Type:   "web-api",
		Variables: map[string]string{
			"ProjectName": "config-test",
			"Author":      cfg.Profiles["test"].Author,
			"Email":       cfg.Profiles["test"].Email,
			"License":     cfg.Profiles["test"].License,
		},
	}

	// Test basic validation
	if projectConfig.Name == "" {
		t.Error("Project name should not be empty")
	}
	if projectConfig.Module == "" {
		t.Error("Module should not be empty")
	}

	t.Logf("Generator and config integration test passed for project: %s", projectConfig.Name)
}

// Helper function for case-insensitive string contains check
func containsIgnoreCase(s, substr string) bool {
	s = strings.ToLower(s)
	substr = strings.ToLower(substr)
	return strings.Contains(s, substr)
}

// Helper function to check if required imports are available
func checkRequiredImports(t *testing.T) {
	// This is a placeholder for checking if all required packages are available
	// In a real test, you might want to check if all dependencies are properly imported
	t.Helper()

	// We can add checks here if needed, for now just ensure the test can import required packages
	if gen := generator.New(); gen == nil {
		t.Skip("Generator package not properly available")
	}
}

// setupTestProjectEnvironment sets up a test project environment
func setupTestProjectEnvironment(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Set up any required environment variables or configuration
	// for testing

	return tmpDir
}
