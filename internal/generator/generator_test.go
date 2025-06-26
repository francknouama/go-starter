package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

func setupTestTemplates(t *testing.T) {
	t.Helper()

	// Get the project root for tests
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	templatesDir := filepath.Join(projectRoot, "templates")

	// Verify templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Fatalf("Templates directory not found at: %s", templatesDir)
	}

	// Set up the filesystem for tests using os.DirFS
	templates.SetTemplatesFS(os.DirFS(templatesDir))
}

func TestNew(t *testing.T) {
	setupTestTemplates(t)

	generator := New()
	if generator == nil {
		t.Error("Expected generator to not be nil")
		return
	}
	if generator.registry == nil {
		t.Error("Expected registry to be initialized")
	}
}

func TestGenerator_validateConfig(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name    string
		config  types.ProjectConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: types.ProjectConfig{
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			wantErr: true,
		},
		{
			name: "missing module",
			config: types.ProjectConfig{
				Name: "test-project",
				Type: "web-api",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			config: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_validateORM(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name    string
		orm     string
		wantErr bool
	}{
		{
			name:    "valid gorm ORM",
			orm:     "gorm",
			wantErr: false,
		},
		{
			name:    "valid raw ORM",
			orm:     "raw",
			wantErr: false,
		},
		{
			name:    "empty ORM (default)",
			orm:     "",
			wantErr: false,
		},
		{
			name:    "unsupported sqlx ORM",
			orm:     "sqlx",
			wantErr: true,
		},
		{
			name:    "unsupported sqlc ORM",
			orm:     "sqlc",
			wantErr: true,
		},
		{
			name:    "unsupported ent ORM",
			orm:     "ent",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.validateORM(tt.orm)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateORM() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerator_getTemplateID(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name     string
		config   types.ProjectConfig
		expected string
	}{
		{
			name: "standard architecture",
			config: types.ProjectConfig{
				Type:         "web-api",
				Architecture: "standard",
			},
			expected: "web-api",
		},
		{
			name: "no architecture",
			config: types.ProjectConfig{
				Type: "web-api",
			},
			expected: "web-api",
		},
		{
			name: "clean architecture",
			config: types.ProjectConfig{
				Type:         "web-api",
				Architecture: "clean",
			},
			expected: "web-api-clean",
		},
		{
			name: "hexagonal architecture",
			config: types.ProjectConfig{
				Type:         "cli",
				Architecture: "hexagonal",
			},
			expected: "cli-hexagonal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.getTemplateID(tt.config)
			if result != tt.expected {
				t.Errorf("getTemplateID() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_createGoMod(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	config := types.ProjectConfig{
		Module:    "github.com/test/project",
		GoVersion: "1.21",
	}

	goModPath := filepath.Join(tempDir, "go.mod")
	err = generator.createGoMod(config, goModPath)
	if err != nil {
		t.Errorf("createGoMod() error = %v", err)
	}

	// Check if file was created
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		t.Error("go.mod file was not created")
	}

	// Check file content
	content, err := os.ReadFile(goModPath)
	if err != nil {
		t.Errorf("Failed to read go.mod file: %v", err)
	}

	expected := "module github.com/test/project\n\ngo 1.21\n"
	if string(content) != expected {
		t.Errorf("go.mod content = %q, want %q", string(content), expected)
	}
}

func TestGenerator_Preview(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	config := types.ProjectConfig{
		Name:         "test-project",
		Module:       "github.com/test/project",
		Type:         "web-api",
		Framework:    "gin",
		Architecture: "clean",
	}

	// Preview should not fail even when template doesn't exist
	err := generator.Preview(config, "/tmp")
	if err != nil {
		t.Errorf("Preview() error = %v", err)
	}
}

func TestGenerator_Generate_ValidationError(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	invalidConfig := types.ProjectConfig{
		Name: "", // Missing required field
	}

	options := types.GenerationOptions{
		OutputPath: "/tmp/test",
		NoGit:      true,
	}

	result, err := generator.Generate(invalidConfig, options)
	if err == nil {
		t.Error("Expected validation error, got nil")
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	if result.Success {
		t.Error("Expected result.Success to be false")
	}
}

func TestGenerator_Generate_TemplateNotFound(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	config := types.ProjectConfig{
		Name:   "test-project",
		Module: "github.com/test/project",
		Type:   "non-existent-template", // Template doesn't exist
	}

	// Create temporary directory for output
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	options := types.GenerationOptions{
		OutputPath: tempDir,
		NoGit:      true,
	}

	result, err := generator.Generate(config, options)
	if err == nil {
		t.Error("Expected template not found error, got nil")
	}

	if result == nil {
		t.Error("Expected result to not be nil")
		return
	}

	if result.Success {
		t.Error("Expected result.Success to be false")
	}

	// Check that the error is a template not found error
	if goStarterErr, ok := err.(*types.GoStarterError); ok {
		if goStarterErr.Code != types.ErrCodeTemplateNotFound {
			t.Errorf("Expected error code %s, got %s", types.ErrCodeTemplateNotFound, goStarterErr.Code)
		}
	} else {
		t.Error("Expected GoStarterError type")
	}
}

func TestGenerator_isGitAvailable(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// This test depends on the system having git installed
	// In most development environments, git should be available
	available := generator.isGitAvailable()

	// We can't assume git is always available, so we just test the function doesn't panic
	if available {
		t.Log("Git is available on this system")
	} else {
		t.Log("Git is not available on this system")
	}
}

func TestGenerator_hasGitRepository(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// Test with a directory that definitely doesn't have git
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// Should return false for directory without git
	if generator.hasGitRepository(tempDir) {
		t.Error("Expected hasGitRepository to return false for directory without git")
	}

	// Create a fake .git directory
	gitDir := filepath.Join(tempDir, ".git")
	if err := os.Mkdir(gitDir, 0755); err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}

	// Should return true for directory with .git
	if !generator.hasGitRepository(tempDir) {
		t.Error("Expected hasGitRepository to return true for directory with .git")
	}
}

func TestGenerator_createGitignore(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	// Create gitignore
	err = generator.createGitignore(tempDir)
	if err != nil {
		t.Errorf("createGitignore() error = %v", err)
	}

	// Check if file was created
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if _, err := os.Stat(gitignorePath); os.IsNotExist(err) {
		t.Error(".gitignore file was not created")
	}

	// Check file content
	content, err := os.ReadFile(gitignorePath)
	if err != nil {
		t.Errorf("Failed to read .gitignore file: %v", err)
	}

	contentStr := string(content)
	expectedPatterns := []string{
		"*.exe",
		"*.test",
		"*.out",
		"vendor/",
		".env",
		".DS_Store",
		"dist/",
		"*.log",
	}

	for _, pattern := range expectedPatterns {
		if !strings.Contains(contentStr, pattern) {
			t.Errorf(".gitignore should contain pattern: %s", pattern)
		}
	}
}

func TestGenerator_processTemplatePath(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name     string
		path     string
		config   types.ProjectConfig
		expected string
	}{
		{
			name: "simple path without variables",
			path: "cmd/main.go",
			config: types.ProjectConfig{
				Name: "test-project",
			},
			expected: "cmd/main.go",
		},
		{
			name: "path with project name variable",
			path: "cmd/{{.ProjectName}}/main.go",
			config: types.ProjectConfig{
				Name: "my-app",
			},
			expected: "cmd/my-app/main.go",
		},
		{
			name: "path with multiple variables",
			path: "internal/{{.Framework}}/{{.ProjectName}}.go",
			config: types.ProjectConfig{
				Name:      "my-api",
				Framework: "gin",
			},
			expected: "internal/gin/my-api.go",
		},
		{
			name: "path with go version",
			path: "go.mod.{{.GoVersion}}",
			config: types.ProjectConfig{
				GoVersion: "1.21",
			},
			expected: "go.mod.1.21",
		},
		{
			name: "path with invalid template (fallback to original)",
			path: "cmd/{{.InvalidVar", // Missing closing braces
			config: types.ProjectConfig{
				Name: "test",
			},
			expected: "cmd/{{.InvalidVar", // Should return original path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.processTemplatePath(tt.path, tt.config)
			if result != tt.expected {
				t.Errorf("processTemplatePath() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_createTemplateContext(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	config := types.ProjectConfig{
		Name:         "test-project",
		Module:       "github.com/test/project",
		Type:         "web-api",
		Framework:    "gin",
		Architecture: "clean",
		GoVersion:    "1.21",
		Author:       "John Doe",
		Email:        "john@example.com",
		License:      "MIT",
		Logger:       "zap",
		Variables: map[string]string{
			"CustomVar": "CustomValue",
		},
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
			Authentication: types.AuthConfig{
				Type: "jwt",
			},
			Logging: types.LoggingConfig{
				Type:       "zap",
				Level:      "debug",
				Format:     "json",
				Structured: true,
			},
		},
	}

	tmpl := types.Template{
		Variables: []types.TemplateVariable{
			{
				Name:    "TemplateVar",
				Default: "DefaultValue",
			},
			{
				Name: "ExistingVar", // Should not override config variable
				Default: "ShouldNotUse",
			},
		},
	}

	context := generator.createTemplateContext(config, tmpl)

	// Test basic variables
	if context["ProjectName"] != "test-project" {
		t.Errorf("Expected ProjectName to be 'test-project', got %v", context["ProjectName"])
	}
	if context["ModulePath"] != "github.com/test/project" {
		t.Errorf("Expected ModulePath to be 'github.com/test/project', got %v", context["ModulePath"])
	}
	if context["Framework"] != "gin" {
		t.Errorf("Expected Framework to be 'gin', got %v", context["Framework"])
	}
	if context["Logger"] != "zap" {
		t.Errorf("Expected Logger to be 'zap', got %v", context["Logger"])
	}

	// Test logger-specific flags
	if !context["UseZap"].(bool) {
		t.Error("Expected UseZap to be true")
	}
	if context["UseSlog"] != nil {
		t.Error("Expected UseSlog to be nil")
	}

	// Test database features
	if context["DatabaseDriver"] != "postgres" {
		t.Errorf("Expected DatabaseDriver to be 'postgres', got %v", context["DatabaseDriver"])
	}
	if context["HasDatabase"] != true {
		t.Errorf("Expected HasDatabase to be true, got %v", context["HasDatabase"])
	}
	if context["HasPostgreSQL"] != true {
		t.Errorf("Expected HasPostgreSQL to be true, got %v", context["HasPostgreSQL"])
	}

	// Test authentication
	if context["AuthType"] != "jwt" {
		t.Errorf("Expected AuthType to be 'jwt', got %v", context["AuthType"])
	}

	// Test ORM
	if context["ORM"] != "gorm" {
		t.Errorf("Expected ORM to be 'gorm', got %v", context["ORM"])
	}

	// Test custom variables
	if context["CustomVar"] != "CustomValue" {
		t.Errorf("Expected CustomVar to be 'CustomValue', got %v", context["CustomVar"])
	}

	// Test template variables
	if context["TemplateVar"] != "DefaultValue" {
		t.Errorf("Expected TemplateVar to be 'DefaultValue', got %v", context["TemplateVar"])
	}

	// Test logger config
	loggerConfig, ok := context["LoggerConfig"].(map[string]interface{})
	if !ok {
		t.Error("Expected LoggerConfig to be a map")
	} else {
		if loggerConfig["Type"] != "zap" {
			t.Errorf("Expected LoggerConfig.Type to be 'zap', got %v", loggerConfig["Type"])
		}
		if loggerConfig["Level"] != "debug" {
			t.Errorf("Expected LoggerConfig.Level to be 'debug', got %v", loggerConfig["Level"])
		}
	}
}

func TestGenerator_getFeatureValue(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name         string
		config       types.ProjectConfig
		feature      string
		key          string
		defaultValue string
		expected     string
	}{
		{
			name: "database driver from features",
			config: types.ProjectConfig{
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Driver: "mysql",
					},
				},
			},
			feature:      "database",
			key:          "driver",
			defaultValue: "postgres",
			expected:     "mysql",
		},
		{
			name: "database ORM from features",
			config: types.ProjectConfig{
				Features: &types.Features{
					Database: types.DatabaseConfig{
						ORM: "raw",
					},
				},
			},
			feature:      "database",
			key:          "orm",
			defaultValue: "gorm",
			expected:     "raw",
		},
		{
			name: "authentication type from features",
			config: types.ProjectConfig{
				Features: &types.Features{
					Authentication: types.AuthConfig{
						Type: "oauth2",
					},
				},
			},
			feature:      "authentication",
			key:          "type",
			defaultValue: "jwt",
			expected:     "oauth2",
		},
		{
			name: "default value when features is nil",
			config: types.ProjectConfig{
				Features: nil,
			},
			feature:      "database",
			key:          "driver",
			defaultValue: "sqlite",
			expected:     "sqlite",
		},
		{
			name: "default value when feature not found",
			config: types.ProjectConfig{
				Features: &types.Features{},
			},
			feature:      "unknown",
			key:          "key",
			defaultValue: "default",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.getFeatureValue(tt.config, tt.feature, tt.key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("getFeatureValue() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_getDatabaseDrivers(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	tests := []struct {
		name     string
		config   types.ProjectConfig
		expected []string
	}{
		{
			name: "multiple database drivers",
			config: types.ProjectConfig{
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Drivers: []string{"postgres", "redis", "mongodb"},
					},
				},
			},
			expected: []string{"postgres", "redis", "mongodb"},
		},
		{
			name: "no features",
			config: types.ProjectConfig{
				Features: nil,
			},
			expected: []string{},
		},
		{
			name: "empty database drivers",
			config: types.ProjectConfig{
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Drivers: []string{},
					},
				},
			},
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.getDatabaseDrivers(tt.config)
			if len(result) != len(tt.expected) {
				t.Errorf("getDatabaseDrivers() length = %v, want %v", len(result), len(tt.expected))
				return
			}
			for i, driver := range result {
				if driver != tt.expected[i] {
					t.Errorf("getDatabaseDrivers()[%d] = %v, want %v", i, driver, tt.expected[i])
				}
			}
		})
	}
}

func TestGenerator_evaluateCondition(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	context := map[string]any{
		"DatabaseDriver": "postgres",
		"Framework":      "gin",
		"HasDatabase":    true,
		"Count":          5,
		"EmptyString":    "",
	}

	tests := []struct {
		name      string
		condition string
		expected  bool
		wantErr   bool
	}{
		{
			name:      "literal true",
			condition: "true",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "literal false",
			condition: "false",
			expected:  false,
			wantErr:   false,
		},
		{
			name:      "string equality check (true)",
			condition: "{{eq .DatabaseDriver \"postgres\"}}",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "string equality check (false)",
			condition: "{{eq .DatabaseDriver \"mysql\"}}",
			expected:  false,
			wantErr:   false,
		},
		{
			name:      "boolean variable (true)",
			condition: "{{.HasDatabase}}",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "numeric non-zero (true)",
			condition: "{{.Count}}",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "numeric zero (false)",
			condition: "0",
			expected:  false,
			wantErr:   false,
		},
		{
			name:      "empty string (false)",
			condition: "{{.EmptyString}}",
			expected:  false,
			wantErr:   false,
		},
		{
			name:      "non-empty string (true)",
			condition: "{{.Framework}}",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "complex condition with ne function",
			condition: "{{ne .DatabaseDriver \"\"}}",
			expected:  true,
			wantErr:   false,
		},
		{
			name:      "invalid template syntax",
			condition: "{{.Invalid}",
			expected:  false,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := generator.evaluateCondition(tt.condition, context)
			if (err != nil) != tt.wantErr {
				t.Errorf("evaluateCondition() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("evaluateCondition() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGenerator_isGoAvailable(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// Test that the function returns a boolean without error
	available := generator.isGoAvailable()
	if available {
		t.Log("Go is available on this system")
	} else {
		t.Log("Go is not available on this system")
	}

	// We can't assume Go is always available, so we just test the function works
	// The important thing is that it doesn't panic and returns a boolean
	if available != true && available != false {
		t.Error("isGoAvailable() should return a boolean value")
	}
}

func TestGenerator_runGitCommand(t *testing.T) {
	setupTestTemplates(t)

	generator := New()

	// Skip test if git is not available
	if !generator.isGitAvailable() {
		t.Skip("Git is not available, skipping git command test")
	}

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer func() {
		if err := os.RemoveAll(tempDir); err != nil {
			t.Logf("Warning: failed to remove temp dir: %v", err)
		}
	}()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "git init command",
			args:    []string{"init"},
			wantErr: false,
		},
		{
			name:    "git status command",
			args:    []string{"status", "--porcelain"},
			wantErr: false,
		},
		{
			name:    "invalid git command",
			args:    []string{"invalid-command"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := generator.runGitCommand(tempDir, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("runGitCommand() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
