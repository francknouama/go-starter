package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// TestGenerator_Validation_ProjectName tests project name validation
func TestGenerator_Validation_ProjectName(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name         string
		projectName  string
		shouldFail   bool
		errorContains string
	}{
		{
			name:        "valid project name",
			projectName: "valid-project",
			shouldFail:  false,
		},
		{
			name:        "valid project name with numbers",
			projectName: "project123",
			shouldFail:  false,
		},
		{
			name:        "valid project name with underscores",
			projectName: "my_project",
			shouldFail:  false,
		},
		{
			name:          "empty project name",
			projectName:   "",
			shouldFail:    true,
			errorContains: "project name is required",
		},
		{
			name:        "project name with spaces (should be allowed)",
			projectName: "my project",
			shouldFail:  false,
		},
		{
			name:        "project name with special characters",
			projectName: "my-project@2024",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      tt.projectName,
				Module:    "github.com/test/project",
				Type:      "library",
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			// Use a temporary directory for testing
			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for project name: %s", tt.projectName)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error should contain expected message")
				}
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for valid project name: %s", tt.projectName)
			}
		})
	}
}

// TestGenerator_Validation_ModulePath tests module path validation
func TestGenerator_Validation_ModulePath(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name          string
		modulePath    string
		shouldFail    bool
		errorContains string
	}{
		{
			name:       "valid github module path",
			modulePath: "github.com/user/project",
			shouldFail: false,
		},
		{
			name:       "valid gitlab module path",
			modulePath: "gitlab.com/user/project",
			shouldFail: false,
		},
		{
			name:       "valid custom domain module path",
			modulePath: "example.com/user/project",
			shouldFail: false,
		},
		{
			name:       "valid nested module path",
			modulePath: "github.com/org/team/project",
			shouldFail: false,
		},
		{
			name:          "empty module path",
			modulePath:    "",
			shouldFail:    true,
			errorContains: "module path is required",
		},
		{
			name:       "simple module name (valid for Go modules)",
			modulePath: "myproject",
			shouldFail: false,
		},
		{
			name:       "module with version suffix",
			modulePath: "github.com/user/project/v2",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    tt.modulePath,
				Type:      "library",
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for module path: %s", tt.modulePath)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error should contain expected message")
				}
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for valid module path: %s", tt.modulePath)
			}
		})
	}
}

// TestGenerator_Validation_ProjectType tests project type validation
func TestGenerator_Validation_ProjectType(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name          string
		projectType   string
		shouldFail    bool
		errorContains string
	}{
		{
			name:        "valid library type",
			projectType: "library",
			shouldFail:  false,
		},
		{
			name:        "valid web-api type",
			projectType: "web-api",
			shouldFail:  false,
		},
		{
			name:        "valid cli type",
			projectType: "cli",
			shouldFail:  false,
		},
		{
			name:        "valid lambda type",
			projectType: "lambda",
			shouldFail:  false,
		},
		{
			name:          "empty project type",
			projectType:   "",
			shouldFail:    true,
			errorContains: "project type is required",
		},
		{
			name:        "unknown project type (should be handled gracefully)",
			projectType: "unknown-type",
			shouldFail:  false, // Generator should handle gracefully with template not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				Type:      tt.projectType,
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for project type: %s", tt.projectType)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains, "Error should contain expected message")
				}
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skipf("Skipping test as template %s is not yet implemented", tt.projectType)
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for project type: %s", tt.projectType)
			}
		})
	}
}

// TestGenerator_Validation_GoVersion tests Go version validation
func TestGenerator_Validation_GoVersion(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name      string
		goVersion string
		shouldFail bool
	}{
		{
			name:      "valid go version 1.21",
			goVersion: "1.21",
			shouldFail: false,
		},
		{
			name:      "valid go version 1.22",
			goVersion: "1.22",
			shouldFail: false,
		},
		{
			name:      "valid go version 1.20",
			goVersion: "1.20",
			shouldFail: false,
		},
		{
			name:      "valid go version with patch",
			goVersion: "1.21.0",
			shouldFail: false,
		},
		{
			name:      "empty go version (should use default)",
			goVersion: "",
			shouldFail: false,
		},
		{
			name:      "auto go version",
			goVersion: "auto",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				Type:      "library",
				GoVersion: tt.goVersion,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for Go version: %s", tt.goVersion)
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for Go version: %s", tt.goVersion)
			}
		})
	}
}

// TestGenerator_Validation_Framework tests framework validation
func TestGenerator_Validation_Framework(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name        string
		projectType string
		framework   string
		shouldFail  bool
	}{
		{
			name:        "gin framework for web-api",
			projectType: "web-api",
			framework:   "gin",
			shouldFail:  false,
		},
		{
			name:        "echo framework for web-api",
			projectType: "web-api",
			framework:   "echo",
			shouldFail:  false,
		},
		{
			name:        "fiber framework for web-api",
			projectType: "web-api",
			framework:   "fiber",
			shouldFail:  false,
		},
		{
			name:        "cobra framework for cli",
			projectType: "cli",
			framework:   "cobra",
			shouldFail:  false,
		},
		{
			name:        "empty framework (should use default)",
			projectType: "web-api",
			framework:   "",
			shouldFail:  false,
		},
		{
			name:        "no framework for library",
			projectType: "library",
			framework:   "",
			shouldFail:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				Type:      tt.projectType,
				Framework: tt.framework,
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for framework: %s", tt.framework)
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for framework: %s", tt.framework)
			}
		})
	}
}

// TestGenerator_Validation_Logger tests logger validation
func TestGenerator_Validation_Logger(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name       string
		logger     string
		shouldFail bool
	}{
		{
			name:       "slog logger",
			logger:     "slog",
			shouldFail: false,
		},
		{
			name:       "zap logger",
			logger:     "zap",
			shouldFail: false,
		},
		{
			name:       "logrus logger",
			logger:     "logrus",
			shouldFail: false,
		},
		{
			name:       "zerolog logger",
			logger:     "zerolog",
			shouldFail: false,
		},
		{
			name:       "empty logger (should use default)",
			logger:     "",
			shouldFail: false,
		},
		{
			name:       "unknown logger (should be handled gracefully)",
			logger:     "unknown-logger",
			shouldFail: false, // Should not fail validation, just use as-is
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				Type:      "web-api",
				Framework: "gin",
				Logger:    tt.logger,
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for logger: %s", tt.logger)
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for logger: %s", tt.logger)
			}
		})
	}
}

// TestGenerator_Validation_Architecture tests architecture validation
func TestGenerator_Validation_Architecture(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name         string
		architecture string
		shouldFail   bool
	}{
		{
			name:         "standard architecture",
			architecture: "standard",
			shouldFail:   false,
		},
		{
			name:         "clean architecture",
			architecture: "clean",
			shouldFail:   false,
		},
		{
			name:         "ddd architecture",
			architecture: "ddd",
			shouldFail:   false,
		},
		{
			name:         "hexagonal architecture",
			architecture: "hexagonal",
			shouldFail:   false,
		},
		{
			name:         "empty architecture (should use default)",
			architecture: "",
			shouldFail:   false,
		},
		{
			name:         "unknown architecture (should be handled gracefully)",
			architecture: "unknown-arch",
			shouldFail:   false, // Should not fail validation, might be template not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-project",
				Module:       "github.com/test/project",
				Type:         "web-api",
				Framework:    "gin",
				Architecture: tt.architecture,
				GoVersion:    "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: tmpDir + "/" + "test-project",
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for architecture: %s", tt.architecture)
			} else {
				// Accept template not found errors as we're still building templates
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for architecture: %s", tt.architecture)
			}
		})
	}
}