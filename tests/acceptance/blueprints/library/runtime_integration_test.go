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

func TestLibrary_OptionalLoggingRuntimeValidation(t *testing.T) {
	// Scenario: Test library with optional logging via dependency injection
	// Given I generate a library project without forced logging dependencies
	// When I run the library examples
	// Then the library should work correctly without any logging
	// And the library should support optional logging via dependency injection

	config := types.ProjectConfig{
		Name:      "test-library-optional-logging",
		Module:    "github.com/test/test-library-optional-logging", 
		Type:      "library",
		GoVersion: "1.21",
	}

	projectPath := helpers.GenerateProject(t, config)
	validator := NewLibraryRuntimeValidator(projectPath)
	
	// Test basic functionality without logging
	validator.TestLibraryCompilation(t)
	validator.TestExamplesExecution(t)
	validator.TestOptionalLoggingPattern(t)
}

func TestLibrary_CoreFeatures(t *testing.T) {
	// Scenario: Test core library features
	// Given I generate a library with core functionality
	// When I test the library functionality
	// Then error handling should be comprehensive
	// And concurrent usage should be safe
	// And the API should follow Go best practices

	config := types.ProjectConfig{
		Name:      "test-library-core",
		Module:    "github.com/test/test-library-core",
		Type:      "library",
		GoVersion: "1.21",
	}

	projectPath := helpers.GenerateProject(t, config)
	validator := NewLibraryRuntimeValidator(projectPath)

	t.Run("error_handling", func(t *testing.T) {
		validator.TestErrorHandling(t)
	})

	t.Run("concurrent_safety", func(t *testing.T) {
		validator.TestConcurrentSafety(t)
	})

	t.Run("api_design", func(t *testing.T) {
		validator.TestAPIDesign(t)
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
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	assert.FileExists(t, mainFile)

	// Check library test file
	testFile := filepath.Join(v.projectPath, projectBaseName+"_test.go")
	assert.FileExists(t, testFile)

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

	// Check that no forced dependencies exist
	goModFile := filepath.Join(v.projectPath, "go.mod")
	assert.FileExists(t, goModFile)
	
	// Should only have testify as dependency (no forced logging deps)
	goModContent := helpers.ReadFileContent(t, goModFile)
	assert.Contains(t, goModContent, "github.com/stretchr/testify")
	
	// Verify no forced logging dependencies
	assert.NotContains(t, goModContent, "go.uber.org/zap")
	assert.NotContains(t, goModContent, "github.com/sirupsen/logrus") 
	assert.NotContains(t, goModContent, "github.com/rs/zerolog")
}

// TestLibraryCompilation validates that the library compiles successfully
func (v *LibraryRuntimeValidator) TestLibraryCompilation(t *testing.T) {
	t.Helper()

	// Ensure module is tidy first
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "go mod tidy should succeed. Output: %s", string(output))

	// Test main library compilation
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Library should compile successfully. Output: %s", string(output))

	// Test that library can be imported as a module
	cmd = exec.Command("go", "list", "-m")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Library should be valid Go module. Output: %s", string(output))
	
	moduleName := strings.TrimSpace(string(output))
	projectBaseName := filepath.Base(v.projectPath)
	expectedModule := "github.com/test/" + projectBaseName
	assert.Equal(t, expectedModule, moduleName)
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
			projectBaseName := filepath.Base(v.projectPath)
			replaceDirective := "github.com/test/" + projectBaseName + "=../.."
			cmd = exec.Command("go", "mod", "edit", "-replace", replaceDirective)
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
			projectBaseName := filepath.Base(v.projectPath)
			assert.Contains(t, outputStr, projectBaseName, "Output should contain project name")
			assert.Contains(t, outputStr, "Example", "Output should indicate it's an example")
			assert.Contains(t, outputStr, "completed", "Example should complete successfully")

			// Example-specific validations for our new simple library
			if example == "basic" {
				assert.Contains(t, outputStr, "default configuration", "Basic example should show default configuration")
				assert.Contains(t, outputStr, "Result:", "Basic example should show result")
			}

			if example == "advanced" {
				assert.Contains(t, outputStr, "logging enabled", "Advanced example should show logging enabled")
				assert.Contains(t, outputStr, "Processing:", "Advanced example should show processing")
				assert.Contains(t, outputStr, "Expected error", "Advanced example should show error handling")
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
	
	// Check for main exported types in our new library structure
	expectedTypes := []string{
		"type Client",
		"type Config",
		"type Logger", // Optional logger interface
		"type Option",
		"func New",
		"func WithLogger",
		"func WithTimeout",
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

	// Run tests first to ensure testify dependency is used
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Tests should run successfully. Output: %s", string(output))

	// Test module tidiness (should pull in testify now)
	cmd = exec.Command("go", "mod", "tidy")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Module should tidy correctly. Output: %s", string(output))

	// Test module verification
	cmd = exec.Command("go", "mod", "verify")
	cmd.Dir = v.projectPath
	output, err = cmd.CombinedOutput()
	assert.NoError(t, err, "Module should verify correctly. Output: %s", string(output))

	// Check go.mod file structure
	goModFile := filepath.Join(v.projectPath, "go.mod")
	content := helpers.ReadFileContent(t, goModFile)
	
	projectBaseName := filepath.Base(v.projectPath)
	expectedModule := "module github.com/test/" + projectBaseName
	assert.Contains(t, content, expectedModule)
	assert.Contains(t, content, "go 1.")
	// Note: testify dependency will be added automatically when tests use it
	// For now, we focus on verifying no forced logging dependencies
	
	// Verify no forced logging dependencies (key aspect of our refactor)
	assert.NotContains(t, content, "go.uber.org/zap")
	assert.NotContains(t, content, "github.com/sirupsen/logrus") 
	assert.NotContains(t, content, "github.com/rs/zerolog")
}

// TestOptionalLoggingPattern validates the optional logging pattern
func (v *LibraryRuntimeValidator) TestOptionalLoggingPattern(t *testing.T) {
	t.Helper()

	// Check that Logger interface exists in the library
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	content := helpers.ReadFileContent(t, mainFile)
	
	// Verify Logger interface is defined
	assert.Contains(t, content, "type Logger interface", "Should define Logger interface")
	assert.Contains(t, content, "Info(msg string, fields ...any)", "Should have Info method")
	assert.Contains(t, content, "Error(msg string, fields ...any)", "Should have Error method")
	
	// Verify WithLogger option exists
	assert.Contains(t, content, "func WithLogger", "Should have WithLogger option")
	assert.Contains(t, content, "logger Logger", "Should accept Logger parameter")
	
	// Verify optional logging pattern in main functionality
	assert.Contains(t, content, "if c.logger != nil", "Should check for logger before logging")
}

// TestAPIDesign validates the library follows Go API design best practices
func (v *LibraryRuntimeValidator) TestAPIDesign(t *testing.T) {
	t.Helper()

	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	content := helpers.ReadFileContent(t, mainFile)
	
	// Verify functional options pattern
	assert.Contains(t, content, "type Option func(*Client)", "Should use functional options pattern")
	assert.Contains(t, content, "func New(opts ...Option)", "Should accept variadic options")
	
	// Verify proper error handling (returns errors instead of logging)
	assert.Contains(t, content, "return \"\", fmt.Errorf", "Should return errors instead of logging them")
	
	// Verify clean API design
	assert.Contains(t, content, "func (c *Client) Process", "Should have main Process method")
	assert.Contains(t, content, "func (c *Client) Close", "Should have Close method for cleanup")
}

// TestErrorHandling validates error handling
func (v *LibraryRuntimeValidator) TestErrorHandling(t *testing.T) {
	t.Helper()

	// Run tests to verify error handling behavior
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	assert.NoError(t, err, "Tests should pass and demonstrate error handling. Output: %s", string(output))

	outputStr := string(output)
	assert.Contains(t, outputStr, "PASS", "Error handling tests should pass")
	// Verify tests include error cases (our tests should test both success and error cases)
	assert.True(t, strings.Contains(outputStr, "empty") || strings.Contains(outputStr, "error") || 
		strings.Contains(outputStr, "invalid"), "Should test error cases")
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
	assert.True(t, strings.Contains(outputStr, "PASS") || strings.Contains(outputStr, "ok"), "Race detector tests should pass")
}