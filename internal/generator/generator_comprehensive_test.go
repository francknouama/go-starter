package generator

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGeneratorComprehensive tests the generator with various scenarios
func TestGeneratorComprehensive(t *testing.T) {
	// Set up mock filesystem with templates
	setupMockTemplateSystem()

	t.Run("Generate_WebAPIStandard", func(t *testing.T) {
		tempDir := t.TempDir()
		outputPath := filepath.Join(tempDir, "test-project")

		generator := New()
		config := types.ProjectConfig{
			Name:         "test-project",
			Module:       "github.com/test/project",
			Type:         "web-api",
			Architecture: "standard",
			Framework:    "gin",
			Logger:       "slog",
			GoVersion:    "1.21",
			Features: &types.Features{
				Database: types.DatabaseConfig{
					Drivers: []string{"postgres"},
					ORM:     "gorm",
				},
				Authentication: types.AuthConfig{
					Type: "jwt",
				},
			},
		}
		options := types.GenerationOptions{
			OutputPath: outputPath,
			DryRun:     false,
			NoGit:      true,
			Verbose:    false,
		}

		result, err := generator.Generate(config, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Greater(t, len(result.FilesCreated), 0)

		// Verify key files were created
		assert.FileExists(t, filepath.Join(outputPath, "go.mod"))
		assert.FileExists(t, filepath.Join(outputPath, "cmd/server/main.go"))
		assert.FileExists(t, filepath.Join(outputPath, "internal/config/config.go"))
	})

	t.Run("Generate_CleanArchitecture", func(t *testing.T) {
		tempDir := t.TempDir()
		outputPath := filepath.Join(tempDir, "clean-project")

		generator := New()
		config := types.ProjectConfig{
			Name:         "clean-project",
			Module:       "github.com/test/clean",
			Type:         "web-api",
			Architecture: "clean",
			Framework:    "echo",
			Logger:       "zap",
			GoVersion:    "1.21",
		}
		options := types.GenerationOptions{
			OutputPath: outputPath,
			NoGit:      true,
		}

		result, err := generator.Generate(config, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Greater(t, len(result.FilesCreated), 0)
	})

	t.Run("Generate_DryRun", func(t *testing.T) {
		generator := New()
		config := types.ProjectConfig{
			Name:      "dry-run-project",
			Module:    "github.com/test/dryrun",
			Type:      "web-api",
			Framework: "gin",
			Logger:    "slog",
			GoVersion: "1.21",
		}
		options := types.GenerationOptions{
			OutputPath: "/tmp/nonexistent",
			DryRun:     true,
		}

		result, err := generator.Generate(config, options)

		require.NoError(t, err)
		assert.True(t, result.Success)
		assert.Len(t, result.FilesCreated, 0) // No files should be created in dry run
	})

	t.Run("Generate_InvalidConfig", func(t *testing.T) {
		generator := New()
		config := types.ProjectConfig{
			// Missing required fields
			Name: "", // Empty name should cause validation error
		}
		options := types.GenerationOptions{
			OutputPath: t.TempDir(),
		}

		result, err := generator.Generate(config, options)

		assert.Error(t, err)
		assert.False(t, result.Success)
		assert.Contains(t, err.Error(), "project name is required")
	})

	t.Run("Generate_TemplateNotFound", func(t *testing.T) {
		generator := New()
		config := types.ProjectConfig{
			Name:      "test-project",
			Module:    "github.com/test/project",
			Type:      "nonexistent-type",
			GoVersion: "1.21",
		}
		options := types.GenerationOptions{
			OutputPath: t.TempDir(),
		}

		result, err := generator.Generate(config, options)

		assert.Error(t, err)
		assert.False(t, result.Success)
	})
}

// TestGeneratorTemplateID tests template ID generation logic
func TestGeneratorTemplateID(t *testing.T) {
	generator := New()

	testCases := []struct {
		name       string
		config     types.ProjectConfig
		expectedID string
	}{
		{
			name: "Standard web-api",
			config: types.ProjectConfig{
				Type:         "web-api",
				Architecture: "standard",
			},
			expectedID: "web-api",
		},
		{
			name: "Clean Architecture",
			config: types.ProjectConfig{
				Type:         "web-api",
				Architecture: "clean",
			},
			expectedID: "web-api-clean",
		},
		{
			name: "DDD Architecture",
			config: types.ProjectConfig{
				Type:         "web-api",
				Architecture: "ddd",
			},
			expectedID: "web-api-ddd",
		},
		{
			name: "Empty architecture defaults to type",
			config: types.ProjectConfig{
				Type:         "cli",
				Architecture: "",
			},
			expectedID: "cli",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			templateID := generator.getTemplateID(tc.config)
			assert.Equal(t, tc.expectedID, templateID)
		})
	}
}

// TestGeneratorValidation tests configuration validation
func TestGeneratorValidation(t *testing.T) {
	generator := New()

	testCases := []struct {
		name        string
		config      types.ProjectConfig
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid config",
			config: types.ProjectConfig{
				Name:      "valid-project",
				Module:    "github.com/test/valid",
				Type:      "web-api",
				GoVersion: "1.21",
			},
			expectError: false,
		},
		{
			name: "Missing name",
			config: types.ProjectConfig{
				Module:    "github.com/test/project",
				Type:      "web-api",
				GoVersion: "1.21",
			},
			expectError: true,
			errorMsg:    "project name is required",
		},
		{
			name: "Missing module",
			config: types.ProjectConfig{
				Name:      "test-project",
				Type:      "web-api",
				GoVersion: "1.21",
			},
			expectError: true,
			errorMsg:    "module path is required",
		},
		{
			name: "Missing type",
			config: types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/test/project",
				GoVersion: "1.21",
			},
			expectError: true,
			errorMsg:    "project type is required",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := generator.validateConfig(tc.config)

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestGeneratorContextCreation tests template context creation
func TestGeneratorContextCreation(t *testing.T) {
	generator := New()

	t.Run("CreateTemplateContext_BasicFields", func(t *testing.T) {
		config := types.ProjectConfig{
			Name:         "test-project",
			Module:       "github.com/test/project",
			Type:         "web-api",
			Architecture: "clean",
			Framework:    "gin",
			Logger:       "slog",
			GoVersion:    "1.21",
			Author:       "Test Author",
			Email:        "test@example.com",
			License:      "MIT",
		}

		template := types.Template{
			Variables: []types.TemplateVariable{
				{Name: "CustomVar", Default: "defaultValue"},
			},
		}

		context := generator.createTemplateContext(config, template)

		assert.Equal(t, "test-project", context["ProjectName"])
		assert.Equal(t, "github.com/test/project", context["ModulePath"])
		assert.Equal(t, "web-api", context["Type"])
		assert.Equal(t, "clean", context["Architecture"])
		assert.Equal(t, "gin", context["Framework"])
		assert.Equal(t, "slog", context["Logger"])
		assert.Equal(t, "1.21", context["GoVersion"])
		assert.Equal(t, "Test Author", context["Author"])
		assert.Equal(t, "test@example.com", context["Email"])
		assert.Equal(t, "MIT", context["License"])
		assert.Equal(t, "defaultValue", context["CustomVar"])
	})

	t.Run("CreateTemplateContext_DatabaseFeatures", func(t *testing.T) {
		config := types.ProjectConfig{
			Name:   "test-project",
			Module: "github.com/test/project",
			Type:   "web-api",
			Features: &types.Features{
				Database: types.DatabaseConfig{
					Drivers: []string{"postgres", "redis"},
					ORM:     "gorm",
				},
			},
		}

		template := types.Template{}
		context := generator.createTemplateContext(config, template)

		assert.Equal(t, "postgres", context["DatabaseDriver"])
		assert.Equal(t, []string{"postgres", "redis"}, context["DatabaseDrivers"])
		assert.True(t, context["HasDatabase"].(bool))
		assert.True(t, context["HasPostgreSQL"].(bool))
		assert.True(t, context["HasRedis"].(bool))
		assert.Equal(t, "gorm", context["ORM"])
	})

	t.Run("CreateTemplateContext_LoggerConfiguration", func(t *testing.T) {
		config := types.ProjectConfig{
			Name:   "test-project",
			Module: "github.com/test/project",
			Type:   "web-api",
			Logger: "zap",
			Features: &types.Features{
				Logging: types.LoggingConfig{
					Type:       "zap",
					Level:      "debug",
					Format:     "json",
					Structured: true,
				},
			},
		}

		template := types.Template{}
		context := generator.createTemplateContext(config, template)

		assert.Equal(t, "zap", context["LoggerType"])
		assert.True(t, context["UseZap"].(bool))

		loggerConfig := context["LoggerConfig"].(map[string]interface{})
		assert.Equal(t, "zap", loggerConfig["Type"])
		assert.Equal(t, "debug", loggerConfig["Level"])
		assert.Equal(t, "json", loggerConfig["Format"])
		assert.True(t, loggerConfig["Structured"].(bool))
	})
}

// TestGeneratorRollback tests the rollback mechanism
func TestGeneratorRollback(t *testing.T) {
	t.Run("GenerationTransaction_AddAndRollback", func(t *testing.T) {
		tempDir := t.TempDir()
		tx := NewGenerationTransaction(tempDir)

		// Create some test files
		testFile1 := filepath.Join(tempDir, "test1.txt")
		testFile2 := filepath.Join(tempDir, "test2.txt")
		testDir := filepath.Join(tempDir, "testdir")

		err := os.WriteFile(testFile1, []byte("test content"), 0644)
		require.NoError(t, err)
		err = os.WriteFile(testFile2, []byte("test content"), 0644)
		require.NoError(t, err)
		err = os.MkdirAll(testDir, 0755)
		require.NoError(t, err)

		// Track files in transaction
		tx.AddFile(testFile1)
		tx.AddFile(testFile2)
		tx.AddDirectory(testDir)

		// Verify files exist
		assert.FileExists(t, testFile1)
		assert.FileExists(t, testFile2)
		assert.DirExists(t, testDir)

		// Rollback
		err = tx.Rollback()
		require.NoError(t, err)

		// Verify files were removed
		assert.NoFileExists(t, testFile1)
		assert.NoFileExists(t, testFile2)
		// Directory might not be removed if not empty, which is expected behavior
	})

	t.Run("GenerationTransaction_EmptyRollback", func(t *testing.T) {
		tempDir := t.TempDir()
		tx := NewGenerationTransaction(tempDir)

		// Rollback with no files tracked should succeed
		err := tx.Rollback()
		assert.NoError(t, err)
	})
}

// setupMockTemplateSystem sets up a mock template system for testing
func setupMockTemplateSystem() {
	mockFS := fstest.MapFS{
		"web-api/template.yaml": &fstest.MapFile{
			Data: []byte(`
name: "web-api-standard"
description: "Standard Web API template"
type: "web-api"
architecture: "standard"
variables:
  - name: "ProjectName"
    type: "string"
    required: true
files:
  - source: "go.mod.tmpl"
    destination: "go.mod"
  - source: "cmd/server/main.go.tmpl"
    destination: "cmd/server/main.go"
  - source: "internal/config/config.go.tmpl"
    destination: "internal/config/config.go"
dependencies:
  - module: "github.com/gin-gonic/gin"
    version: "v1.9.1"
    condition: "{{eq .Framework \"gin\"}}"
`),
		},
		"web-api/go.mod.tmpl": &fstest.MapFile{
			Data: []byte(`module {{.ModulePath}}

go {{.GoVersion}}
`),
		},
		"web-api/cmd/server/main.go.tmpl": &fstest.MapFile{
			Data: []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello {{.ProjectName}}")
}
`),
		},
		"web-api/internal/config/config.go.tmpl": &fstest.MapFile{
			Data: []byte(`package config

type Config struct {
	ProjectName string
}
`),
		},
		"web-api-clean/template.yaml": &fstest.MapFile{
			Data: []byte(`
name: "web-api-clean"
description: "Clean Architecture Web API template"
type: "web-api"
architecture: "clean"
files:
  - source: "go.mod.tmpl"
    destination: "go.mod"
`),
		},
		"web-api-clean/go.mod.tmpl": &fstest.MapFile{
			Data: []byte(`module {{.ModulePath}}

go {{.GoVersion}}
`),
		},
	}

	templates.SetTemplatesFS(mockFS)
}
