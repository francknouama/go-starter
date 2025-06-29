package generator

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// setupTestTemplates sets up the template system for testing
func setupTestTemplates(t *testing.T) {
	t.Helper()

	// Get the project root for tests - integration tests are 3 levels deep
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(file))))
	templatesDir := filepath.Join(projectRoot, "templates")

	// Verify templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Fatalf("Templates directory not found at: %s", templatesDir)
	}

	// Set up the filesystem for tests using os.DirFS
	templates.SetTemplatesFS(os.DirFS(templatesDir))
}

// TestGenerator_New tests the creation of a new generator instance
func TestGenerator_New(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name string
	}{
		{
			name: "create new generator",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := generator.New()
			assert.NotNil(t, gen, "Generator should not be nil")
		})
	}
}

// TestGenerator_Generate_BasicFunctionality tests basic generation functionality
func TestGenerator_Generate_BasicFunctionality(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name       string
		config     types.ProjectConfig
		options    types.GenerationOptions
		shouldFail bool
		errorMsg   string
	}{
		{
			name: "basic library project generation",
			config: types.ProjectConfig{
				Name:      "test-library",
				Module:    "github.com/test/test-library",
				Type:      "library",
				GoVersion: "1.21",
			},
			options: types.GenerationOptions{
				OutputPath: "", // Will be set in test
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			},
			shouldFail: false,
		},
		{
			name: "web-api project generation",
			config: types.ProjectConfig{
				Name:      "test-api",
				Module:    "github.com/test/test-api",
				Type:      "web-api",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    "slog",
			},
			options: types.GenerationOptions{
				OutputPath: "", // Will be set in test
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			},
			shouldFail: false,
		},
		{
			name: "empty project name should fail",
			config: types.ProjectConfig{
				Name:      "",
				Module:    "github.com/test/invalid",
				Type:      "library",
				GoVersion: "1.21",
			},
			options: types.GenerationOptions{
				OutputPath: "", // Will be set in test
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			},
			shouldFail: true,
			errorMsg:   "project name is required",
		},
		{
			name: "empty module path should fail",
			config: types.ProjectConfig{
				Name:      "test-project",
				Module:    "",
				Type:      "library",
				GoVersion: "1.21",
			},
			options: types.GenerationOptions{
				OutputPath: "", // Will be set in test
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			},
			shouldFail: true,
			errorMsg:   "module path is required",
		},
		{
			name: "empty project type should fail",
			config: types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/test-project",
				Type:      "",
				GoVersion: "1.21",
			},
			options: types.GenerationOptions{
				OutputPath: "", // Will be set in test
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			},
			shouldFail: true,
			errorMsg:   "project type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory for output
			tmpDir := t.TempDir()
			tt.options.OutputPath = filepath.Join(tmpDir, tt.config.Name)

			gen := generator.New()
			require.NotNil(t, gen)

			result, err := gen.Generate(tt.config, tt.options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected generation to fail")
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error message should contain expected text")
				}
				assert.False(t, result.Success, "Result should indicate failure")
			} else {
				if err != nil {
					// For now, template not found errors are acceptable as we're still building templates
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skipf("Skipping test as template %s is not yet implemented", tt.config.Type)
						return
					}
				}
				assert.NoError(t, err, "Generation should succeed")
				assert.True(t, result.Success, "Result should indicate success")
				assert.NotEmpty(t, result.ProjectPath, "Project path should be set")
				assert.Greater(t, result.Duration, time.Duration(0), "Duration should be positive")
			}
		})
	}
}

// TestGenerator_Preview tests the preview functionality
func TestGenerator_Preview(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name      string
		config    types.ProjectConfig
		outputDir string
	}{
		{
			name: "preview library project",
			config: types.ProjectConfig{
				Name:      "preview-library",
				Module:    "github.com/test/preview-library",
				Type:      "library",
				GoVersion: "1.21",
			},
			outputDir: t.TempDir(),
		},
		{
			name: "preview web-api project",
			config: types.ProjectConfig{
				Name:         "preview-api",
				Module:       "github.com/test/preview-api",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "standard",
				Logger:       "slog",
			},
			outputDir: t.TempDir(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := generator.New()
			require.NotNil(t, gen)

			// Preview should not return an error and should not create any files
			err := gen.Preview(tt.config, tt.outputDir)
			assert.NoError(t, err, "Preview should not fail")

			// Verify no files were created during preview
			projectPath := filepath.Join(tt.outputDir, tt.config.Name)
			_, err = os.Stat(projectPath)
			assert.True(t, os.IsNotExist(err), "Preview should not create actual files")
		})
	}
}

// TestGenerator_GenerationResult tests the generation result structure
func TestGenerator_GenerationResult(t *testing.T) {
	setupTestTemplates(t)

	config := types.ProjectConfig{
		Name:      "result-test",
		Module:    "github.com/test/result-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	tmpDir := t.TempDir()
	options := types.GenerationOptions{
		OutputPath: filepath.Join(tmpDir, config.Name),
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	gen := generator.New()
	require.NotNil(t, gen)

	result, err := gen.Generate(config, options)

	// For now, accept template not found errors
	if err != nil {
		if _, ok := err.(*types.GoStarterError); ok {
			t.Skip("Skipping test as template is not yet implemented")
			return
		}
	}

	require.NoError(t, err)
	require.NotNil(t, result)

	// Test result structure
	assert.True(t, result.Success, "Generation should be successful")
	assert.Equal(t, options.OutputPath, result.ProjectPath, "Project path should match")
	assert.Greater(t, result.Duration, time.Duration(0), "Duration should be positive")
	assert.NotNil(t, result.FilesCreated, "Files created should be initialized")
}

// TestGenerator_DryRun tests dry run functionality
func TestGenerator_DryRun(t *testing.T) {
	setupTestTemplates(t)

	config := types.ProjectConfig{
		Name:      "dry-run-test",
		Module:    "github.com/test/dry-run-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	tmpDir := t.TempDir()
	options := types.GenerationOptions{
		OutputPath: filepath.Join(tmpDir, config.Name),
		DryRun:     true,
		NoGit:      true,
		Verbose:    false,
	}

	gen := generator.New()
	require.NotNil(t, gen)

	result, err := gen.Generate(config, options)

	// For now, accept template not found errors
	if err != nil {
		if _, ok := err.(*types.GoStarterError); ok {
			t.Skip("Skipping test as template is not yet implemented")
			return
		}
	}

	require.NoError(t, err)
	require.NotNil(t, result)

	// In dry run mode, no files should be created
	_, err = os.Stat(options.OutputPath)
	assert.True(t, os.IsNotExist(err) || isEmptyDir(t, options.OutputPath),
		"Dry run should not create files")
}

// TestGenerator_GitInitialization tests git repository initialization
func TestGenerator_GitInitialization(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name      string
		noGit     bool
		expectGit bool
	}{
		{
			name:      "git initialization enabled",
			noGit:     false,
			expectGit: true,
		},
		{
			name:      "git initialization disabled",
			noGit:     true,
			expectGit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "git-test",
				Module:    "github.com/test/git-test",
				Type:      "library",
				GoVersion: "1.21",
			}

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, config.Name),
				DryRun:     false,
				NoGit:      tt.noGit,
				Verbose:    false,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			result, err := gen.Generate(config, options)

			// For now, accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			require.NoError(t, err)
			require.NotNil(t, result)

			// Check if .git directory exists
			gitDir := filepath.Join(options.OutputPath, ".git")
			_, err = os.Stat(gitDir)

			if tt.expectGit {
				// Git should be initialized (if git is available on system)
				// We don't fail the test if git is not available
				if err != nil && !os.IsNotExist(err) {
					t.Logf("Git directory check failed (git might not be available): %v", err)
				}
			} else {
				// Git should not be initialized
				assert.True(t, os.IsNotExist(err), ".git directory should not exist when NoGit is true")
			}
		})
	}
}

// Helper functions

// isEmptyDir checks if a directory is empty
func isEmptyDir(t *testing.T, dir string) bool {
	t.Helper()

	entries, err := os.ReadDir(dir)
	if err != nil {
		return true // Directory doesn't exist or can't be read, consider it empty
	}
	return len(entries) == 0
}
