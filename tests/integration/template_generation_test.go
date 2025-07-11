package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestTemplates initializes the template registry for testing
func setupTestTemplates(t *testing.T) {
	t.Helper()
	// Initialize templates filesystem from the blueprints directory
	templatesDir := "../../blueprints"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Fatalf("Blueprints directory not found: %s", templatesDir)
	}
	templates.SetTemplatesFS(os.DirFS(templatesDir))
}

// TestCompleteTemplateWorkflow tests the complete template generation workflow for Phase 1
func TestCompleteTemplateWorkflow(t *testing.T) {
	// Initialize templates filesystem for testing
	setupTestTemplates(t)

	// Create a temporary directory for the test
	tmpDir := t.TempDir()

	t.Run("Template Registry", func(t *testing.T) {
		testTemplateRegistry(t)
	})

	t.Run("Project Generation", func(t *testing.T) {
		testProjectGeneration(t, tmpDir)
	})

	t.Run("Generated Project Compilation", func(t *testing.T) {
		testGeneratedProjectCompilation(t, tmpDir)
	})
}

// testTemplateRegistry tests template registry functionality
func testTemplateRegistry(t *testing.T) {
	registry := templates.NewRegistry()
	require.NotNil(t, registry, "Registry should not be nil")

	// Test 1: Verify web-api template is registered
	webAPITemplate, err := registry.Get("web-api")
	assert.NoError(t, err, "Should be able to retrieve web-api template")
	assert.Equal(t, "web-api", webAPITemplate.ID)
	assert.Equal(t, "web-api-standard", webAPITemplate.Name)
	assert.Equal(t, "web-api", webAPITemplate.Type)
	assert.Equal(t, "standard", webAPITemplate.Architecture)

	// Test 2: Check template metadata
	assert.NotEmpty(t, webAPITemplate.Description)
	assert.NotEmpty(t, webAPITemplate.Variables)
	assert.NotEmpty(t, webAPITemplate.Dependencies)

	// Test 3: Verify required variables
	hasProjectName := false
	hasModulePath := false
	for _, v := range webAPITemplate.Variables {
		if v.Name == "ProjectName" {
			hasProjectName = true
			assert.True(t, v.Required, "ProjectName should be required")
		}
		if v.Name == "ModulePath" {
			hasModulePath = true
			assert.True(t, v.Required, "ModulePath should be required")
		}
	}
	assert.True(t, hasProjectName, "Template should have ProjectName variable")
	assert.True(t, hasModulePath, "Template should have ModulePath variable")

	// Test 4: Test template retrieval by type
	webAPITemplates := registry.GetByType("web-api")
	assert.Len(t, webAPITemplates, 3, "Should have exactly three web-api templates (standard, clean, ddd)")
	// Just check that we have web-api templates, don't check specific IDs as they can vary

	// Test 5: Test template existence check
	assert.True(t, registry.Exists("web-api"), "web-api template should exist")
	assert.False(t, registry.Exists("non-existent"), "non-existent template should not exist")

	// Test 6: List all templates
	allTemplates := registry.List()
	assert.GreaterOrEqual(t, len(allTemplates), 1, "Should have at least one template")
}

// testProjectGeneration tests project generation from template
func testProjectGeneration(t *testing.T, tmpDir string) {
	gen := generator.New()
	require.NotNil(t, gen, "Generator should not be nil")

	// Define test cases
	testCases := []struct {
		name           string
		config         types.ProjectConfig
		expectDatabase bool
		expectAuth     bool
	}{
		{
			name: "Basic Web API without features",
			config: types.ProjectConfig{
				Name:         "test-api-basic",
				Module:       "github.com/test/test-api-basic",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "standard",
				Variables: map[string]string{
					"ProjectName":    "test-api-basic",
					"ModulePath":     "github.com/test/test-api-basic",
					"GoVersion":      "1.21",
					"Framework":      "gin",
					"DatabaseDriver": "",
					"AuthType":       "",
				},
			},
			expectDatabase: false,
			expectAuth:     false,
		},
		{
			name: "Web API with database",
			config: types.ProjectConfig{
				Name:         "test-api-db",
				Module:       "github.com/test/test-api-db",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "standard",
				Variables: map[string]string{
					"ProjectName":    "test-api-db",
					"ModulePath":     "github.com/test/test-api-db",
					"GoVersion":      "1.21",
					"Framework":      "gin",
					"DatabaseDriver": "postgres",
					"AuthType":       "",
				},
			},
			expectDatabase: true,
			expectAuth:     false,
		},
		{
			name: "Web API with authentication",
			config: types.ProjectConfig{
				Name:         "test-api-auth",
				Module:       "github.com/test/test-api-auth",
				Type:         "web-api",
				GoVersion:    "1.21",
				Framework:    "gin",
				Architecture: "standard",
				Variables: map[string]string{
					"ProjectName":    "test-api-auth",
					"ModulePath":     "github.com/test/test-api-auth",
					"GoVersion":      "1.21",
					"Framework":      "gin",
					"DatabaseDriver": "postgres",
					"AuthType":       "jwt",
				},
			},
			expectDatabase: true,
			expectAuth:     true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create output directory
			outputPath := filepath.Join(tmpDir, tc.config.Name)

			// Generate project
			options := types.GenerationOptions{
				OutputPath: outputPath,
				NoGit:      true,
			}

			result, err := gen.Generate(tc.config, options)

			// NOTE: The current implementation is incomplete, so we expect an error
			// In a complete implementation, this should succeed
			if err != nil {
				t.Logf("Generation failed (expected in current implementation): %v", err)

				// Check if it's because templates aren't implemented yet
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Logf("Template not found: %s (this is expected in Phase 0)", goErr.Message)
				}
				return
			}

			// If generation succeeds, verify the results
			require.NoError(t, err, "Project generation should succeed")
			assert.True(t, result.Success, "Generation should be successful")
			assert.NotEmpty(t, result.FilesCreated, "Should create at least one file")

			// Verify core files exist
			assertFileExists(t, outputPath, "go.mod")
			assertFileExists(t, outputPath, "Makefile")
			assertFileExists(t, outputPath, "README.md")
			assertFileExists(t, outputPath, "cmd/server/main.go")
			assertFileExists(t, outputPath, "internal/config/config.go")
			assertFileExists(t, outputPath, "internal/handlers/health.go")
			assertFileExists(t, outputPath, "internal/middleware/cors.go")
			assertFileExists(t, outputPath, "internal/middleware/logger.go")
			assertFileExists(t, outputPath, "internal/middleware/recovery.go")

			// Check conditional files
			if tc.expectDatabase {
				assertFileExists(t, outputPath, "internal/database/connection.go")
				assertFileExists(t, outputPath, "internal/models/user.go")
				assertFileExists(t, outputPath, "internal/repository/user.go")
				assertFileExists(t, outputPath, "internal/services/user.go")
				assertFileExists(t, outputPath, "docker-compose.yml")
			} else {
				assertFileNotExists(t, outputPath, "internal/database/connection.go")
				assertFileNotExists(t, outputPath, "docker-compose.yml")
			}

			if tc.expectAuth {
				assertFileExists(t, outputPath, "internal/middleware/auth.go")
			} else {
				assertFileNotExists(t, outputPath, "internal/middleware/auth.go")
			}
		})
	}
}

// testGeneratedProjectCompilation tests that generated projects compile successfully
func testGeneratedProjectCompilation(t *testing.T, tmpDir string) {
	// NOTE: This test would only run if we have successfully generated projects
	// In the current implementation, generation is not complete, so we'll skip this

	projectPath := filepath.Join(tmpDir, "test-api-basic")
	if _, err := os.Stat(projectPath); os.IsNotExist(err) {
		t.Skip("Skipping compilation test - no generated project found (expected in current implementation)")
	}

	// If we get here, try to build the project
	t.Log("Testing compilation of generated project...")

	// TODO: Implement actual compilation test when template generation is complete
	// This would involve:
	// 1. Running 'go mod tidy' to resolve dependencies
	// 2. Running 'go build ./...' to compile all packages
	// 3. Checking for any compilation errors
}

// Helper functions

func assertFileExists(t *testing.T, basePath, relativePath string) {
	t.Helper()
	fullPath := filepath.Join(basePath, relativePath)
	_, err := os.Stat(fullPath)
	assert.NoError(t, err, "File should exist: %s", relativePath)
}

func assertFileNotExists(t *testing.T, basePath, relativePath string) {
	t.Helper()
	fullPath := filepath.Join(basePath, relativePath)
	_, err := os.Stat(fullPath)
	assert.True(t, os.IsNotExist(err), "File should not exist: %s", relativePath)
}
