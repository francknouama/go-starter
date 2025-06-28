package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupCompilationTestTemplates sets up templates for compilation tests
// (renamed to avoid conflict with existing setupTestTemplates)
func setupCompilationTestTemplates(t *testing.T) {
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

// TestTemplateCompilation tests that all generated projects compile successfully
func TestTemplateCompilation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping template compilation tests in short mode")
	}

	// Initialize templates
	setupCompilationTestTemplates(t)
	gen := generator.New()

	// Create temporary directory for test projects
	tmpDir, err := os.MkdirTemp("", "go-starter-compilation-test-*")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to clean up temp dir: %v", err)
		}
	}()

	templateTests := []struct {
		name     string
		config   types.ProjectConfig
		validate func(t *testing.T, projectPath string)
	}{
		{
			name: "web-api-standard",
			config: types.ProjectConfig{
				Name:      "test-web-api",
				Module:    "github.com/test/web-api",
				Type:      "web-api",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    "slog",
			},
			validate: func(t *testing.T, projectPath string) {
				// Check main.go exists
				mainPath := filepath.Join(projectPath, "cmd", "server", "main.go")
				assert.FileExists(t, mainPath)

				// Check go.mod exists and has correct module
				goModPath := filepath.Join(projectPath, "go.mod")
				assert.FileExists(t, goModPath)

				// Check Makefile exists
				makefilePath := filepath.Join(projectPath, "Makefile")
				assert.FileExists(t, makefilePath)

				// Check key directories exist
				assert.DirExists(t, filepath.Join(projectPath, "internal", "handlers"))
				assert.DirExists(t, filepath.Join(projectPath, "internal", "middleware"))
				assert.DirExists(t, filepath.Join(projectPath, "internal", "logger"))
			},
		},
		{
			name: "cli-standard",
			config: types.ProjectConfig{
				Name:      "test-cli",
				Module:    "github.com/test/cli",
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "slog",
			},
			validate: func(t *testing.T, projectPath string) {
				// Check main.go exists
				mainPath := filepath.Join(projectPath, "main.go")
				assert.FileExists(t, mainPath)

				// Check cmd directory structure
				assert.DirExists(t, filepath.Join(projectPath, "cmd"))
				assert.FileExists(t, filepath.Join(projectPath, "cmd", "root.go"))
			},
		},
		{
			name: "library-standard",
			config: types.ProjectConfig{
				Name:      "test-library",
				Module:    "github.com/test/library",
				Type:      "library",
				GoVersion: "1.21",
				Logger:    "slog",
			},
			validate: func(t *testing.T, projectPath string) {
				// Check main library file exists
				libPath := filepath.Join(projectPath, "test-library.go")
				assert.FileExists(t, libPath)

				// Check test file exists
				testPath := filepath.Join(projectPath, "test-library_test.go")
				assert.FileExists(t, testPath)

				// Check examples directory
				assert.DirExists(t, filepath.Join(projectPath, "examples"))
			},
		},
		{
			name: "lambda-standard",
			config: types.ProjectConfig{
				Name:      "test-lambda",
				Module:    "github.com/test/lambda",
				Type:      "lambda",
				GoVersion: "1.21",
				Logger:    "slog",
			},
			validate: func(t *testing.T, projectPath string) {
				// Check main.go and handler.go exist
				assert.FileExists(t, filepath.Join(projectPath, "main.go"))
				assert.FileExists(t, filepath.Join(projectPath, "handler.go"))

				// Check template.yaml for SAM deployment
				assert.FileExists(t, filepath.Join(projectPath, "template.yaml"))
			},
		},
	}

	for _, tt := range templateTests {
		t.Run(tt.name, func(t *testing.T) {
			// Set output path for this test
			projectPath := filepath.Join(tmpDir, tt.config.Name)
			options := types.GenerationOptions{
				OutputPath: projectPath,
				NoGit:      true,
			}

			// Generate the project
			result, err := gen.Generate(tt.config, options)
			require.NoError(t, err, "Failed to generate project %s", tt.name)
			require.True(t, result.Success, "Generation should be successful")

			// Run template-specific validations
			tt.validate(t, projectPath)

			// Test compilation
			t.Run("compilation", func(t *testing.T) {
				testCompilation(t, projectPath, tt.config.Name)
			})

			// Test go mod tidy works
			t.Run("go_mod_tidy", func(t *testing.T) {
				testGoModTidy(t, projectPath)
			})

			// Test go vet passes
			t.Run("go_vet", func(t *testing.T) {
				testGoVet(t, projectPath)
			})
		})
	}
}

// testCompilation tests that the generated project compiles
func testCompilation(t *testing.T, projectPath, projectName string) {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Go build output for %s:\n%s", projectName, string(output))
		t.Fatalf("Failed to compile project %s: %v", projectName, err)
	}

	t.Logf("✓ Project %s compiled successfully", projectName)
}

// testGoModTidy tests that go mod tidy works on the generated project
func testGoModTidy(t *testing.T, projectPath string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Go mod tidy output:\n%s", string(output))
		t.Fatalf("Failed to run go mod tidy: %v", err)
	}
}

// testGoVet tests that go vet passes on the generated project
func testGoVet(t *testing.T, projectPath string) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Logf("Go vet output:\n%s", string(output))
		t.Fatalf("go vet failed: %v", err)
	}
}

// TestTemplateWithDifferentLoggers tests templates with different logger configurations
func TestTemplateWithDifferentLoggers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping logger variation tests in short mode")
	}

	// Initialize templates
	setupCompilationTestTemplates(t)
	gen := generator.New()

	// Create temporary directory for test projects
	tmpDir, err := os.MkdirTemp("", "go-starter-logger-test-*")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to clean up temp dir: %v", err)
		}
	}()

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("web-api-with-"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-api-" + logger,
				Module:    "github.com/test/api-" + logger,
				Type:      "web-api",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    logger,
			}

			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, "test-api-"+logger),
				NoGit:      true,
			}

			// Generate the project
			result, err := gen.Generate(config, options)
			require.NoError(t, err, "Failed to generate project with %s logger", logger)
			require.True(t, result.Success, "Generation should be successful")

			// Test compilation
			testCompilation(t, options.OutputPath, config.Name)

			// Verify logger-specific files are generated
			loggerFilePath := filepath.Join(options.OutputPath, "internal", "logger", logger+".go")
			assert.FileExists(t, loggerFilePath, "Logger-specific file should exist for %s", logger)
		})
	}
}

// TestTemplateWithDatabaseOptions tests templates with different database configurations
func TestTemplateWithDatabaseOptions(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping database variation tests in short mode")
	}

	// Initialize templates
	setupCompilationTestTemplates(t)
	gen := generator.New()

	// Create temporary directory for test projects
	tmpDir, err := os.MkdirTemp("", "go-starter-db-test-*")
	require.NoError(t, err)
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Failed to clean up temp dir: %v", err)
		}
	}()

	databases := []string{"postgres", "mysql", "sqlite"}

	for _, db := range databases {
		t.Run("web-api-with-"+db, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-api-" + db,
				Module:    "github.com/test/api-" + db,
				Type:      "web-api",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Driver: db,
						ORM:    "gorm",
					},
				},
			}

			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, "test-api-"+db),
				NoGit:      true,
			}

			// Generate the project
			result, err := gen.Generate(config, options)
			require.NoError(t, err, "Failed to generate project with %s database", db)
			require.True(t, result.Success, "Generation should be successful")

			// Test compilation
			testCompilation(t, options.OutputPath, config.Name)

			// Verify database-specific files are generated
			assert.FileExists(t, filepath.Join(options.OutputPath, "internal", "database", "connection.go"))
			assert.FileExists(t, filepath.Join(options.OutputPath, "internal", "models", "user.go"))
			assert.FileExists(t, filepath.Join(options.OutputPath, "docker-compose.yml"))
		})
	}
}

// TestTemplateValidation tests that all templates pass validation
func TestTemplateValidation(t *testing.T) {
	// Initialize templates
	setupCompilationTestTemplates(t)
	registry := templates.NewRegistry()

	// Get all templates
	allTemplates := registry.List()
	require.NotEmpty(t, allTemplates, "No templates found in registry")

	for _, template := range allTemplates {
		t.Run("validate_"+template.ID, func(t *testing.T) {
			// Basic validation
			assert.NotEmpty(t, template.ID, "Template ID should not be empty")
			assert.NotEmpty(t, template.Type, "Template type should not be empty")
			assert.NotEmpty(t, template.Files, "Template should have files defined")

			// Validate that required files exist in template.yaml
			hasMainFile := false
			hasGoMod := false
			for _, file := range template.Files {
				if file.Source == "main.go.tmpl" || file.Source == "cmd/server/main.go.tmpl" {
					hasMainFile = true
				}
				if file.Source == "go.mod.tmpl" {
					hasGoMod = true
				}
			}

			if template.Type != "library" {
				assert.True(t, hasMainFile, "Template %s should have a main file", template.ID)
			}
			assert.True(t, hasGoMod, "Template %s should have go.mod", template.ID)
		})
	}
}
