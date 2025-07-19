package library_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize templates filesystem for ATDD tests
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				hasTemplates := false
				for _, entry := range entries {
					if entry.IsDir() {
						templateYaml := filepath.Join(templatesDir, entry.Name(), "template.yaml")
						if _, err := os.Stat(templateYaml); err == nil {
							hasTemplates = true
							break
						}
					}
				}

				if hasTemplates {
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// Feature: Library Blueprint Runtime Integration Testing
// As a developer
// I want to validate that generated library projects work correctly at runtime
// So that I can be confident the blueprint generates functional Go libraries

func TestLibrary_RuntimeIntegration_ComprehensiveLibraryTesting(t *testing.T) {
	// Scenario: Generate library and test all major functionality
	// Given I generate a Go library with all features
	// When I test the library functionality
	// Then all library patterns should work correctly
	// And the library should follow Go best practices
	// And examples should compile and run successfully

	config := types.ProjectConfig{
		Name:      "test-library-runtime",
		Module:    "github.com/test/test-library-runtime",
		Type:      "library",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)
	
	// Build and test the library
	validator := NewLibraryRuntimeValidator(projectPath)

	// Test library structure and compilation
	t.Run("library_structure", func(t *testing.T) {
		validator.TestLibraryStructure(t)
	})

	t.Run("library_compilation", func(t *testing.T) {
		validator.TestLibraryCompilation(t)
	})

	t.Run("unit_tests", func(t *testing.T) {
		validator.TestUnitTests(t)
	})

	t.Run("benchmark_tests", func(t *testing.T) {
		validator.TestBenchmarks(t)
	})

	t.Run("examples_compilation", func(t *testing.T) {
		validator.TestExamplesCompilation(t)
	})

	t.Run("examples_execution", func(t *testing.T) {
		validator.TestExamplesExecution(t)
	})

	t.Run("library_api", func(t *testing.T) {
		validator.TestLibraryAPI(t)
	})

	t.Run("go_doc_generation", func(t *testing.T) {
		validator.TestGoDocGeneration(t)
	})

	t.Run("module_behavior", func(t *testing.T) {
		validator.TestModuleBehavior(t)
	})
}

func TestLibrary_MultiLoggerRuntimeValidation(t *testing.T) {
	// Scenario: Test library with different logger types
	// Given I generate library projects with different loggers
	// When I run the library examples
	// Then all logger types should work correctly
	// And logging behavior should be appropriate for libraries

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("logger_"+logger+"_runtime", func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-library-" + logger + "-runtime",
				Module:    "github.com/test/test-library-" + logger + "-runtime", 
				Type:      "library",
				GoVersion: "1.21",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)
			validator := NewLibraryRuntimeValidator(projectPath)
			
			// Test basic functionality with this logger
			validator.TestLibraryCompilation(t)
			validator.TestExamplesExecution(t)
			validator.TestLoggerBehavior(t, logger)
		})
	}
}

func TestLibrary_AdvancedFeatures(t *testing.T) {
	// Scenario: Test advanced library features
	// Given I generate a library with all features enabled
	// When I test advanced functionality
	// Then caching, rate limiting, metrics should work correctly
	// And error handling should be comprehensive
	// And concurrent usage should be safe

	config := types.ProjectConfig{
		Name:      "test-library-advanced",
		Module:    "github.com/test/test-library-advanced",
		Type:      "library",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)
	validator := NewLibraryRuntimeValidator(projectPath)

	t.Run("caching_functionality", func(t *testing.T) {
		validator.TestCachingFeatures(t)
	})

	t.Run("rate_limiting", func(t *testing.T) {
		validator.TestRateLimiting(t)
	})

	t.Run("metrics_collection", func(t *testing.T) {
		validator.TestMetricsCollection(t)
	})

	t.Run("error_handling", func(t *testing.T) {
		validator.TestErrorHandling(t)
	})

	t.Run("concurrent_safety", func(t *testing.T) {
		validator.TestConcurrentSafety(t)
	})
}

// LibraryRuntimeValidator provides comprehensive runtime testing for library projects
type LibraryRuntimeValidator struct {
	projectPath string
}

// NewLibraryRuntimeValidator creates a new library runtime validator
func NewLibraryRuntimeValidator(projectPath string) *LibraryRuntimeValidator {
	return &LibraryRuntimeValidator{
		projectPath: projectPath,
	}
}

// TestLibraryStructure validates the generated library structure
func (v *LibraryRuntimeValidator) TestLibraryStructure(t *testing.T) {
	t.Helper()

	// Check main library file
	mainFile := filepath.Join(v.projectPath, "test-library-runtime.go")
	assert.FileExists(t, mainFile)

	// Check pkg structure
	pkgDir := filepath.Join(v.projectPath, "pkg", "test_library_runtime")
	assert.DirExists(t, pkgDir)

	expectedFiles := []string{
		"client.go",
		"types.go", 
		"options.go",
		"client_test.go",
		"benchmark_test.go",
	}

	for _, file := range expectedFiles {
		filePath := filepath.Join(pkgDir, file)
		assert.FileExists(t, filePath, "Expected file %s should exist", file)
	}

	// Check internal structure
	internalDir := filepath.Join(v.projectPath, "internal")
	assert.DirExists(t, internalDir)

	internalSubdirs := []string{"cache", "ratelimiter", "logger"}
	for _, subdir := range internalSubdirs {
		subdirPath := filepath.Join(internalDir, subdir)
		assert.DirExists(t, subdirPath, "Expected internal subdir %s should exist", subdir)
	}

	// Check examples
	examplesDir := filepath.Join(v.projectPath, "examples")
	assert.DirExists(t, examplesDir)

	exampleSubdirs := []string{"basic", "advanced"}
	for _, subdir := range exampleSubdirs {
		subdirPath := filepath.Join(examplesDir, subdir)
		assert.DirExists(t, subdirPath, "Expected example subdir %s should exist", subdir)
		
		mainFile := filepath.Join(subdirPath, "main.go")
		assert.FileExists(t, mainFile, "Expected example main.go should exist in %s", subdir)
	}
}

// TestLibraryCompilation validates that the library compiles successfully
func (v *LibraryRuntimeValidator) TestLibraryCompilation(t *testing.T) {
	t.Helper()

	// Test main library compilation
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Library should compile successfully. Output: %s", string(output))

	// Test that library can be imported as a module
	cmd = exec.Command("go", "list", "-m", ".")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Library should be valid Go module. Output: %s", string(output))
	
	moduleName := strings.TrimSpace(string(output))
	assert.Equal(t, "github.com/test/test-library-runtime", moduleName)
}

// TestUnitTests runs all unit tests in the library
func (v *LibraryRuntimeValidator) TestUnitTests(t *testing.T) {
	t.Helper()

	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Unit tests should pass. Output: %s", string(output))

	// Check that tests actually ran
	outputStr := string(output)
	assert.Contains(t, outputStr, "PASS", "Tests should show PASS status")
	assert.Contains(t, outputStr, "ok", "Tests should show ok status")
}

// TestBenchmarks runs benchmark tests
func (v *LibraryRuntimeValidator) TestBenchmarks(t *testing.T) {
	t.Helper()

	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Benchmarks should run successfully. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Benchmark", "Should contain benchmark results")
	assert.Contains(t, outputStr, "ns/op", "Should show performance metrics")
}

// TestExamplesCompilation validates that examples compile
func (v *LibraryRuntimeValidator) TestExamplesCompilation(t *testing.T) {
	t.Helper()

	examples := []string{"basic", "advanced"}
	
	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			exampleDir := filepath.Join(v.projectPath, "examples", example)
			
			// Initialize go module for example
			cmd := exec.Command("go", "mod", "init", "example")
			cmd.Dir = exampleDir
			_, err := cmd.CombinedOutput()
			require.NoError(t, err, "Should be able to initialize go module for example")

			// Add replace directive for local library
			cmd = exec.Command("go", "mod", "edit", "-replace", "github.com/test/test-library-runtime=../..")
			cmd.Dir = exampleDir
			_, err = cmd.CombinedOutput()
			require.NoError(t, err, "Should be able to add replace directive")

			// Tidy dependencies
			cmd = exec.Command("go", "mod", "tidy")
			cmd.Dir = exampleDir
			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "Should be able to tidy dependencies. Output: %s", string(output))

			// Compile example
			cmd = exec.Command("go", "build", "main.go")
			cmd.Dir = exampleDir
			output, err = cmd.CombinedOutput()
			assert.NoError(t, err, "Example %s should compile successfully. Output: %s", example, string(output))
		})
	}
}

// TestExamplesExecution runs the examples and validates output
func (v *LibraryRuntimeValidator) TestExamplesExecution(t *testing.T) {
	t.Helper()

	examples := []string{"basic", "advanced"}
	
	for _, example := range examples {
		t.Run(example, func(t *testing.T) {
			exampleDir := filepath.Join(v.projectPath, "examples", example)
			
			// Run example with timeout
			cmd := exec.Command("timeout", "30s", "go", "run", "main.go")
			cmd.Dir = exampleDir
			output, err := cmd.CombinedOutput()
			assert.NoError(t, err, "Example %s should run successfully. Output: %s", example, string(output))

			outputStr := string(output)
			
			// Basic validations for all examples
			assert.Contains(t, outputStr, "test-library-runtime", "Output should contain library name")
			assert.Contains(t, outputStr, "Example", "Output should indicate it's an example")
			assert.Contains(t, outputStr, "completed", "Example should complete successfully")

			// Example-specific validations
			if example == "basic" {
				assert.Contains(t, outputStr, "Health Check", "Basic example should show health check")
				assert.Contains(t, outputStr, "Processing", "Basic example should show processing")
				assert.Contains(t, outputStr, "Batch Processing", "Basic example should show batch processing")
			}

			if example == "advanced" {
				assert.Contains(t, outputStr, "Cache", "Advanced example should show caching")
				assert.Contains(t, outputStr, "Rate Limiting", "Advanced example should show rate limiting")
				assert.Contains(t, outputStr, "Event", "Advanced example should show events")
				assert.Contains(t, outputStr, "Metrics", "Advanced example should show metrics")
			}
		})
	}
}

// TestLibraryAPI validates the public API structure
func (v *LibraryRuntimeValidator) TestLibraryAPI(t *testing.T) {
	t.Helper()

	// Test that main types are exported properly
	cmd := exec.Command("go", "doc", ".")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Should be able to generate docs. Output: %s", string(output))

	outputStr := string(output)
	
	// Check for main exported types
	expectedTypes := []string{
		"type Client",
		"type Input",
		"type Output", 
		"type Config",
		"type Option",
		"func New",
	}

	for _, expectedType := range expectedTypes {
		assert.Contains(t, outputStr, expectedType, "API should export %s", expectedType)
	}
}

// TestGoDocGeneration validates that godoc works correctly
func (v *LibraryRuntimeValidator) TestGoDocGeneration(t *testing.T) {
	t.Helper()

	// Test package documentation
	cmd := exec.Command("go", "doc", "-all", ".")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Should generate complete documentation. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Package", "Should contain package documentation")
	assert.Contains(t, outputStr, "functionality", "Should contain package description")
}

// TestModuleBehavior validates Go module behavior
func (v *LibraryRuntimeValidator) TestModuleBehavior(t *testing.T) {
	t.Helper()

	// Test module tidiness
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Module should tidy correctly. Output: %s", string(output))

	// Test module verification
	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Module should verify correctly. Output: %s", string(output))

	// Check go.mod file structure
	goModFile := filepath.Join(v.projectPath, "go.mod")
	content := helpers.ReadFileContent(t, goModFile)
	
	assert.Contains(t, content, "module github.com/test/test-library-runtime")
	assert.Contains(t, content, "go 1.21")
	assert.Contains(t, content, "github.com/stretchr/testify")
}

// TestLoggerBehavior validates logger-specific behavior
func (v *LibraryRuntimeValidator) TestLoggerBehavior(t *testing.T, logger string) {
	t.Helper()

	// Run basic example and check for logger output
	exampleDir := filepath.Join(v.projectPath, "examples", "basic")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Basic example should run with %s logger. Output: %s", logger, string(output))

	// Logger-specific validations would go here
	// For now, just ensure it runs without errors
}

// TestCachingFeatures validates caching functionality
func (v *LibraryRuntimeValidator) TestCachingFeatures(t *testing.T) {
	t.Helper()

	// Run advanced example which demonstrates caching
	exampleDir := filepath.Join(v.projectPath, "examples", "advanced")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Advanced example should demonstrate caching. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "cache miss", "Should show cache miss")
	assert.Contains(t, outputStr, "cache hit", "Should show cache hit")
}

// TestRateLimiting validates rate limiting functionality
func (v *LibraryRuntimeValidator) TestRateLimiting(t *testing.T) {
	t.Helper()

	// Run advanced example which demonstrates rate limiting
	exampleDir := filepath.Join(v.projectPath, "examples", "advanced")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Advanced example should demonstrate rate limiting. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Rate limited", "Should show rate limiting in action")
}

// TestMetricsCollection validates metrics functionality
func (v *LibraryRuntimeValidator) TestMetricsCollection(t *testing.T) {
	t.Helper()

	// Run basic example which shows final metrics
	exampleDir := filepath.Join(v.projectPath, "examples", "basic")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Basic example should show metrics. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Total Processed", "Should show total processed count")
	assert.Contains(t, outputStr, "Average Processing Time", "Should show average processing time")
}

// TestErrorHandling validates error handling
func (v *LibraryRuntimeValidator) TestErrorHandling(t *testing.T) {
	t.Helper()

	// Run advanced example which demonstrates error handling
	exampleDir := filepath.Join(v.projectPath, "examples", "advanced")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = exampleDir
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Advanced example should demonstrate error handling. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "Expected error", "Should show expected error handling")
	assert.Contains(t, outputStr, "Validation failed", "Should show validation errors")
}

// TestConcurrentSafety validates concurrent usage
func (v *LibraryRuntimeValidator) TestConcurrentSafety(t *testing.T) {
	t.Helper()

	// Run tests with race detector
	cmd := exec.Command("go", "test", "-race", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Tests should pass with race detector. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "PASS", "Race detector tests should pass")
}