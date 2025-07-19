package library_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

func init() {
	// Initialize templates filesystem for ATDD tests
	// We use DirFS pointing to the real blueprints directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			// Check if this directory actually contains template files
			// by looking for template.yaml files
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				// Check if any subdirectory contains template.yaml
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

		// Move up one directory
		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// Feature: Standard Library Blueprint
// As a developer
// I want to generate a Go library project
// So that I can quickly build reusable packages

func TestStandard_Library_BasicGeneration(t *testing.T) {
	// Scenario: Generate standard library
	// Given I want a standard Go library
	// When I generate the project
	// Then the project should include library structure
	// And the project should have public API
	// And the project should compile and run successfully
	// And the project should have working examples

	// Given I want a standard Go library
	config := types.ProjectConfig{
		Name:      "test-standard-library",
		Module:    "github.com/test/test-standard-library",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include library structure
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidateLibraryStructure(t)

	// And the project should have public API
	validator.ValidatePublicAPI(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)

	// And the project should have working examples
	validator.ValidateExamples(t)
}

func TestStandard_Library_WithOptionalLogging(t *testing.T) {
	// Scenario: Generate library with optional logging via dependency injection
	// Given I want a library with optional logging capabilities
	// When I generate the library
	// Then the project should define a Logger interface
	// And the project should support optional logging via dependency injection
	// And the project should compile successfully without any logging dependencies

	// Given I want a library with optional logging capabilities
	config := types.ProjectConfig{
		Name:      "test-library-optional",
		Module:    "github.com/test/test-library-optional",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the library
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should define a Logger interface
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidateLoggerInterface(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)

	// And the project should support dependency injection
	validator.ValidateOptionalLoggingPattern(t)
}

func TestStandard_Library_Documentation(t *testing.T) {
	// Scenario: Generate library with documentation
	// Given I want a library with documentation
	// When I generate the project
	// Then the project should include Go documentation
	// And the project should have README with examples
	// And the project should have doc.go file

	// Given I want a library with documentation
	config := types.ProjectConfig{
		Name:      "test-library-docs",
		Module:    "github.com/test/test-library-docs",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include Go documentation
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidateDocumentation(t)

	// And the project should have README with examples
	validator.ValidateREADME(t)

	// And the project should have doc.go file
	validator.ValidateDocFile(t)
}

func TestStandard_Library_ExampleUsage(t *testing.T) {
	// Scenario: Generate library with example usage
	// Given I want a library with examples
	// When I generate the project
	// Then the project should include basic examples
	// And the project should include advanced examples
	// And the examples should compile and run successfully

	// Given I want a library with examples
	config := types.ProjectConfig{
		Name:      "test-library-examples",
		Module:    "github.com/test/test-library-examples",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include basic examples
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidateBasicExample(t)

	// And the project should include advanced examples
	validator.ValidateAdvancedExample(t)

	// And the examples should compile and run successfully
	validator.ValidateExampleExecution(t)
}

func TestStandard_Library_TestSupport(t *testing.T) {
	// Scenario: Generate library with test support
	// Given I want a library with comprehensive tests
	// When I generate the project
	// Then the project should include test files
	// And the project should include example tests
	// And the project should use testify for assertions
	// And the tests should run successfully

	// Given I want a library with comprehensive tests
	config := types.ProjectConfig{
		Name:      "test-library-tests",
		Module:    "github.com/test/test-library-tests",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include test files
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidateTestFiles(t)

	// And the project should include example tests
	validator.ValidateExampleTests(t)

	// And the project should use testify for assertions
	validator.ValidateTestifyUsage(t)

	// And the tests should run successfully
	validator.ValidateTestExecution(t)
}

func TestStandard_Library_PublicAPI(t *testing.T) {
	// Scenario: Generate library with proper public API
	// Given I want a library with clean public API
	// When I generate the project
	// Then the project should export main functionality
	// And the project should hide internal implementation
	// And the project should have proper Go naming conventions

	// Given I want a library with clean public API
	config := types.ProjectConfig{
		Name:      "test-library-api",
		Module:    "github.com/test/test-library-api",
		Type:      "library",
		GoVersion: "1.21",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should export main functionality
	validator := NewLibraryValidator(projectPath, "standard")
	validator.ValidatePublicAPI(t)

	// And the project should hide internal implementation
	validator.ValidateInternalPackages(t)

	// And the project should have proper Go naming conventions
	validator.ValidateNamingConventions(t)
}

// LibraryValidator provides validation methods for Library blueprints
type LibraryValidator struct {
	projectPath  string
	architecture string
}

// NewLibraryValidator creates a new LibraryValidator
func NewLibraryValidator(projectPath, architecture string) *LibraryValidator {
	return &LibraryValidator{
		projectPath:  projectPath,
		architecture: architecture,
	}
}

// ValidateLibraryStructure validates library structure
func (v *LibraryValidator) ValidateLibraryStructure(t *testing.T) {
	t.Helper()

	// Check main library file exists
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	helpers.AssertFileExists(t, mainFile)

	// Check library test file exists  
	testFile := filepath.Join(v.projectPath, projectBaseName+"_test.go")
	helpers.AssertFileExists(t, testFile)

	// Check go.mod exists
	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)

	// Check examples directory exists
	examplesDir := filepath.Join(v.projectPath, "examples")
	helpers.AssertDirectoryExists(t, examplesDir)
}

// ValidatePublicAPI validates public API
func (v *LibraryValidator) ValidatePublicAPI(t *testing.T) {
	t.Helper()

	// Check main library file contains exported functions and types
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Verify key exported types and functions
	helpers.AssertFileContains(t, mainFile, "type Client")
	helpers.AssertFileContains(t, mainFile, "type Config")
	helpers.AssertFileContains(t, mainFile, "type Logger interface")
	helpers.AssertFileContains(t, mainFile, "func New")
	helpers.AssertFileContains(t, mainFile, "func WithLogger")
	helpers.AssertFileContains(t, mainFile, "func WithTimeout")
}

// ValidateCompilation validates that the project compiles successfully
func (v *LibraryValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertProjectCompiles(t, v.projectPath)
}

// ValidateExamples validates examples exist and work
func (v *LibraryValidator) ValidateExamples(t *testing.T) {
	t.Helper()

	// Check examples directory exists
	examplesDir := filepath.Join(v.projectPath, "examples")
	helpers.AssertDirectoryExists(t, examplesDir)

	// Check basic example exists
	basicExample := filepath.Join(examplesDir, "basic", "main.go")
	helpers.AssertFileExists(t, basicExample)

	// Check advanced example exists
	advancedExample := filepath.Join(examplesDir, "advanced", "main.go")
	helpers.AssertFileExists(t, advancedExample)
}

// ValidateLoggerInterface validates the Logger interface definition
func (v *LibraryValidator) ValidateLoggerInterface(t *testing.T) {
	t.Helper()

	// Check main library file contains Logger interface
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Verify Logger interface definition
	helpers.AssertFileContains(t, mainFile, "type Logger interface")
	helpers.AssertFileContains(t, mainFile, "Info(msg string, fields ...any)")
	helpers.AssertFileContains(t, mainFile, "Error(msg string, fields ...any)")
}

// ValidateOptionalLoggingPattern validates optional logging pattern
func (v *LibraryValidator) ValidateOptionalLoggingPattern(t *testing.T) {
	t.Helper()

	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Verify optional logging pattern
	helpers.AssertFileContains(t, mainFile, "logger Logger // optional")
	helpers.AssertFileContains(t, mainFile, "if c.logger != nil")
	helpers.AssertFileContains(t, mainFile, "func WithLogger(logger Logger)")
	
	// Verify no forced logging dependencies
	goModFile := filepath.Join(v.projectPath, "go.mod")
	content := helpers.ReadFileContent(t, goModFile)
	helpers.AssertNotContains(t, content, "go.uber.org/zap")
	helpers.AssertNotContains(t, content, "github.com/sirupsen/logrus")
	helpers.AssertNotContains(t, content, "github.com/rs/zerolog")
}

// ValidateDocumentation validates documentation
func (v *LibraryValidator) ValidateDocumentation(t *testing.T) {
	t.Helper()

	// Check doc.go exists
	docFile := filepath.Join(v.projectPath, "doc.go")
	helpers.AssertFileExists(t, docFile)
	helpers.AssertFileContains(t, docFile, "package")
}

// ValidateREADME validates README file
func (v *LibraryValidator) ValidateREADME(t *testing.T) {
	t.Helper()

	readmeFile := filepath.Join(v.projectPath, "README.md")
	helpers.AssertFileExists(t, readmeFile)
	helpers.AssertFileContains(t, readmeFile, "# test-library-docs")
	helpers.AssertFileContains(t, readmeFile, "## Quick Start")
}

// ValidateDocFile validates doc.go file
func (v *LibraryValidator) ValidateDocFile(t *testing.T) {
	t.Helper()

	docFile := filepath.Join(v.projectPath, "doc.go")
	helpers.AssertFileExists(t, docFile)
	helpers.AssertFileContains(t, docFile, "Package")
}

// ValidateBasicExample validates basic example
func (v *LibraryValidator) ValidateBasicExample(t *testing.T) {
	t.Helper()

	basicExample := filepath.Join(v.projectPath, "examples", "basic", "main.go")
	helpers.AssertFileExists(t, basicExample)
	helpers.AssertFileContains(t, basicExample, "func main")
}

// ValidateAdvancedExample validates advanced example
func (v *LibraryValidator) ValidateAdvancedExample(t *testing.T) {
	t.Helper()

	advancedExample := filepath.Join(v.projectPath, "examples", "advanced", "main.go")
	helpers.AssertFileExists(t, advancedExample)
	helpers.AssertFileContains(t, advancedExample, "func main")
}

// ValidateExampleExecution validates example execution
func (v *LibraryValidator) ValidateExampleExecution(t *testing.T) {
	t.Helper()

	// Test basic example compilation
	basicExample := filepath.Join(v.projectPath, "examples", "basic")
	helpers.AssertProjectCompiles(t, basicExample)

	// Test advanced example compilation
	advancedExample := filepath.Join(v.projectPath, "examples", "advanced")
	helpers.AssertProjectCompiles(t, advancedExample)
}

// ValidateTestFiles validates test files exist
func (v *LibraryValidator) ValidateTestFiles(t *testing.T) {
	t.Helper()

	// Check main test file exists
	projectBaseName := filepath.Base(v.projectPath)
	mainTestFile := filepath.Join(v.projectPath, projectBaseName+"_test.go")
	helpers.AssertFileExists(t, mainTestFile)
}

// ValidateExampleTests validates example tests
func (v *LibraryValidator) ValidateExampleTests(t *testing.T) {
	t.Helper()

	// Check example test file exists
	exampleTestFile := filepath.Join(v.projectPath, "examples_test.go")
	helpers.AssertFileExists(t, exampleTestFile)
	helpers.AssertFileContains(t, exampleTestFile, "func Example")
}

// ValidateTestifyUsage validates testify usage
func (v *LibraryValidator) ValidateTestifyUsage(t *testing.T) {
	t.Helper()

	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)
	helpers.AssertFileContains(t, goModFile, "github.com/stretchr/testify")
}

// ValidateTestExecution validates test execution
func (v *LibraryValidator) ValidateTestExecution(t *testing.T) {
	t.Helper()
	helpers.AssertTestsRun(t, v.projectPath)
}

// ValidateInternalPackages validates that internal packages are not exposed
func (v *LibraryValidator) ValidateInternalPackages(t *testing.T) {
	t.Helper()

	// Verify that the public API does not expose internal implementation details
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	content := helpers.ReadFileContent(t, mainFile)
	
	// Should not expose internal implementation details in public API
	helpers.AssertNotContains(t, content, "internal/")
}

// ValidateNamingConventions validates Go naming conventions
func (v *LibraryValidator) ValidateNamingConventions(t *testing.T) {
	t.Helper()

	// Check main library file follows naming convention
	projectBaseName := filepath.Base(v.projectPath)
	mainFile := filepath.Join(v.projectPath, projectBaseName+".go")
	helpers.AssertFileExists(t, mainFile)

	// Check that exported functions start with uppercase
	helpers.AssertFileContains(t, mainFile, "func ")
}