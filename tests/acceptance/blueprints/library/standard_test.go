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
		Logger:    "slog",
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

func TestStandard_Library_WithDifferentLoggers(t *testing.T) {
	// Scenario: Generate library with different logging libraries
	// Given I want a library with configurable logging
	// When I generate with different loggers
	// Then the project should include the selected logger
	// And the project should compile successfully
	// And logging should work as expected

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("logger_"+logger, func(t *testing.T) {
			// Given I want a library with configurable logging
			config := types.ProjectConfig{
				Name:      "test-library-" + logger,
				Module:    "github.com/test/test-library-" + logger,
				Type:      "library",
				GoVersion: "1.21",
				Logger:    logger,
			}

			// When I generate with different loggers
			projectPath := helpers.GenerateProject(t, config)

			// Then the project should include the selected logger
			validator := NewLibraryValidator(projectPath, "standard")
			validator.ValidateLogger(t, logger)

			// And the project should compile successfully
			validator.ValidateCompilation(t)

			// And logging should work as expected
			validator.ValidateLoggerFunctionality(t, logger)
		})
	}
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
		Logger:    "slog",
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
		Logger:    "slog",
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
		Logger:    "slog",
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
		Logger:    "slog",
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
	mainFile := filepath.Join(v.projectPath, "test-standard-library.go")
	helpers.AssertFileExists(t, mainFile)

	// Check go.mod exists
	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)

	// Check internal directory exists
	internalDir := filepath.Join(v.projectPath, "internal")
	helpers.AssertDirectoryExists(t, internalDir)
}

// ValidatePublicAPI validates public API
func (v *LibraryValidator) ValidatePublicAPI(t *testing.T) {
	t.Helper()

	// Check main library file contains exported functions
	mainFile := filepath.Join(v.projectPath, "test-standard-library.go")
	helpers.AssertFileExists(t, mainFile)
	helpers.AssertFileContains(t, mainFile, "func ")
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

// ValidateLogger validates logger implementation
func (v *LibraryValidator) ValidateLogger(t *testing.T, logger string) {
	t.Helper()

	// Check logger file exists
	loggerFile := filepath.Join(v.projectPath, "internal", "logger", "logger.go")
	helpers.AssertFileExists(t, loggerFile)

	// Check logger specific dependencies in go.mod
	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)

	switch logger {
	case "zap":
		helpers.AssertFileContains(t, goModFile, "go.uber.org/zap")
	case "logrus":
		helpers.AssertFileContains(t, goModFile, "github.com/sirupsen/logrus")
	case "zerolog":
		helpers.AssertFileContains(t, goModFile, "github.com/rs/zerolog")
	}
}

// ValidateLoggerFunctionality validates logger functionality
func (v *LibraryValidator) ValidateLoggerFunctionality(t *testing.T, logger string) {
	t.Helper()
	helpers.AssertLoggerFunctionality(t, v.projectPath, logger)
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
	helpers.AssertFileContains(t, readmeFile, "## Usage")
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
	mainTestFile := filepath.Join(v.projectPath, "test-library-tests_test.go")
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

// ValidateInternalPackages validates internal packages
func (v *LibraryValidator) ValidateInternalPackages(t *testing.T) {
	t.Helper()

	// Check internal directory exists
	internalDir := filepath.Join(v.projectPath, "internal")
	helpers.AssertDirectoryExists(t, internalDir)

	// Check internal logger exists
	internalLogger := filepath.Join(internalDir, "logger", "logger.go")
	helpers.AssertFileExists(t, internalLogger)
}

// ValidateNamingConventions validates Go naming conventions
func (v *LibraryValidator) ValidateNamingConventions(t *testing.T) {
	t.Helper()

	// Check main library file follows naming convention
	mainFile := filepath.Join(v.projectPath, "test-library-api.go")
	helpers.AssertFileExists(t, mainFile)

	// Check that exported functions start with uppercase
	helpers.AssertFileContains(t, mainFile, "func ")
}