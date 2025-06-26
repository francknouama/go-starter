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

func TestGenerator_getTemplateID(t *testing.T) {
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
	generator := New()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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
	defer os.RemoveAll(tempDir)

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
	generator := New()

	// Test with a directory that definitely doesn't have git
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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
	generator := New()

	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "go-starter-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

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
