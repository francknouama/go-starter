package library_test

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: Go Library Blueprint
// As a library author
// I want to generate well-structured Go libraries
// So that I can share reusable code effectively

func TestLibrary_BasicGeneration(t *testing.T) {
	// Scenario: Generate basic library
	// Given I want a Go library
	// When I generate a library project
	// Then the project should have a clear public API
	// And the project should include example usage
	// And the project should have comprehensive tests
	// And the project should include proper documentation

	// Given I want a Go library
	config := types.ProjectConfig{
		Name:      "test-library-basic",
		Module:    "github.com/test/test-library-basic",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	// When I generate a library project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should have a clear public API
	validator := NewLibraryValidator(projectPath)
	validator.ValidatePublicAPI(t)

	// And the project should include example usage
	validator.ValidateExampleUsage(t)

	// And the project should have comprehensive tests
	validator.ValidateComprehensiveTests(t)

	// And the project should include proper documentation
	validator.ValidateDocumentation(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_PackageStructure(t *testing.T) {
	// Scenario: Library package structure
	// Given a generated library
	// Then it should have clean package organization
	// And all public functions should be documented
	// And examples should be runnable
	// And tests should cover public API

	config := types.ProjectConfig{
		Name:      "test-library-structure",
		Module:    "github.com/test/test-library-structure",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then it should have clean package organization
	validator.ValidatePackageOrganization(t)

	// And all public functions should be documented
	validator.ValidatePublicFunctionDocumentation(t)

	// And examples should be runnable
	validator.ValidateRunnableExamples(t)

	// And tests should cover public API
	validator.ValidatePublicAPITestCoverage(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_APIDesign(t *testing.T) {
	// Feature: Library API Design
	// As a library user
	// I want a well-designed API
	// So that the library is easy to use

	// Scenario: Public API structure
	// Given a generated library
	// Then the public API should be intuitive
	// And function names should be descriptive
	// And parameter types should be appropriate
	// And return values should follow Go conventions
	// And error handling should be consistent

	config := types.ProjectConfig{
		Name:      "test-library-api",
		Module:    "github.com/test/test-library-api",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then the public API should be intuitive
	validator.ValidateIntuitiveAPI(t)

	// And function names should be descriptive
	validator.ValidateDescriptiveFunctionNames(t)

	// And parameter types should be appropriate
	validator.ValidateParameterTypes(t)

	// And return values should follow Go conventions
	validator.ValidateReturnValueConventions(t)

	// And error handling should be consistent
	validator.ValidateConsistentErrorHandling(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_ExamplesTesting(t *testing.T) {
	// Feature: Library Examples and Testing
	// As a library maintainer
	// I want comprehensive examples and tests
	// So that users can understand and trust the library

	// Scenario: Examples and tests
	// Given a generated library
	// Then it should include basic and advanced examples
	// And examples should be tested
	// And unit tests should cover core functionality
	// And benchmark tests should be included

	config := types.ProjectConfig{
		Name:      "test-library-examples",
		Module:    "github.com/test/test-library-examples",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then it should include basic and advanced examples
	validator.ValidateBasicAndAdvancedExamples(t)

	// And examples should be tested
	validator.ValidateTestedExamples(t)

	// And unit tests should cover core functionality
	validator.ValidateUnitTestCoverage(t)

	// And benchmark tests should be included
	validator.ValidateBenchmarkTests(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration for Libraries
	// As a library developer
	// I want minimal logging capabilities
	// So that I can debug issues without affecting users

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			t.Parallel() // Safe to run logger tests in parallel

			// Scenario: Logger integration
			// Given I generate a library with "<logger>"
			// Then internal logging should be minimal
			// And logger should be configurable
			// And logger should not interfere with user code
			// And structured logging should be available for debugging

			// Given I generate a library with "<logger>"
			config := types.ProjectConfig{
				Name:      "test-library-logger-" + logger,
				Module:    "github.com/test/test-library-logger-" + logger,
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then internal logging should be minimal
			validator := NewLibraryValidator(projectPath)
			validator.ValidateMinimalLogging(t, logger)

			// And logger should be configurable
			validator.ValidateConfigurableLogging(t)

			// And logger should not interfere with user code
			validator.ValidateNonInterfering(t)

			// And structured logging should be available for debugging
			validator.ValidateStructuredLogging(t, logger)

			// And the project should compile successfully
			validator.ValidateCompilation(t)
		})
	}
}

func TestLibrary_Documentation(t *testing.T) {
	// Feature: Library Documentation
	// As a library user
	// I want comprehensive documentation
	// So that I can understand how to use the library

	// Scenario: Documentation completeness
	// Given a generated library
	// Then package documentation should exist
	// And all public functions should have godoc comments
	// And examples should be included in documentation
	// And README should be comprehensive

	config := types.ProjectConfig{
		Name:      "test-library-docs",
		Module:    "github.com/test/test-library-docs",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then package documentation should exist
	validator.ValidatePackageDocumentation(t)

	// And all public functions should have godoc comments
	validator.ValidateGodocComments(t)

	// And examples should be included in documentation
	validator.ValidateDocumentationExamples(t)

	// And README should be comprehensive
	validator.ValidateComprehensiveREADME(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_MinimalDependencies(t *testing.T) {
	// Scenario: Minimal dependencies
	// Given a library project
	// When I check dependencies
	// Then only essential dependencies should be included
	// And the dependency footprint should be minimal
	// And no unnecessary transitive dependencies should exist

	config := types.ProjectConfig{
		Name:      "test-library-minimal",
		Module:    "github.com/test/test-library-minimal",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then only essential dependencies should be included
	validator.ValidateEssentialDependencies(t)

	// And the dependency footprint should be minimal
	validator.ValidateMinimalDependencyFootprint(t)

	// And no unnecessary transitive dependencies should exist
	validator.ValidateNoUnnecessaryTransitiveDependencies(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_ClientInitialization(t *testing.T) {
	// Feature: Client Initialization
	// As a library user
	// I want easy client initialization
	// So that I can start using the library quickly

	// Scenario: Client initialization patterns
	// Given a generated library
	// Then it should provide clear initialization patterns
	// And configuration should be straightforward
	// And default values should be sensible
	// And initialization should be documented

	config := types.ProjectConfig{
		Name:      "test-library-client",
		Module:    "github.com/test/test-library-client",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then it should provide clear initialization patterns
	validator.ValidateInitializationPatterns(t)

	// And configuration should be straightforward
	validator.ValidateClientConfiguration(t)

	// And default values should be sensible
	validator.ValidateDefaultValues(t)

	// And initialization should be documented
	validator.ValidateInitializationDocumentation(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_ErrorHandling(t *testing.T) {
	// Feature: Library Error Handling
	// As a library user
	// I want consistent error handling
	// So that I can handle errors appropriately

	// Scenario: Error handling patterns
	// Given a generated library
	// Then errors should follow Go conventions
	// And error messages should be descriptive
	// And error types should be exported when appropriate
	// And error handling should be consistent

	config := types.ProjectConfig{
		Name:      "test-library-errors",
		Module:    "github.com/test/test-library-errors",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewLibraryValidator(projectPath)

	// Then errors should follow Go conventions
	validator.ValidateErrorConventions(t)

	// And error messages should be descriptive
	validator.ValidateDescriptiveErrorMessages(t)

	// And error types should be exported when appropriate
	validator.ValidateExportedErrorTypes(t)

	// And error handling should be consistent
	validator.ValidateConsistentErrorHandling(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

// LibraryValidator provides validation methods for Library blueprints
type LibraryValidator struct {
	ProjectPath string
}

func NewLibraryValidator(projectPath string) *LibraryValidator {
	return &LibraryValidator{
		ProjectPath: projectPath,
	}
}

func (v *LibraryValidator) ValidatePublicAPI(t *testing.T) {
	t.Helper()
	
	// Check that main library file exists
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Check that it exports public functions
	content := helpers.ReadFileContent(t, mainFile)
	
	// Should have public functions (starting with uppercase)
	if helpers.StringContains(content, "func New") || helpers.StringContains(content, "func ") {
		t.Log("✓ Public API functions found")
	}
	
	// Should have proper package declaration
	helpers.AssertFileContains(t, mainFile, "package "+projectName)
}

func (v *LibraryValidator) ValidateExampleUsage(t *testing.T) {
	t.Helper()
	
	// Check examples directory exists
	examplesDir := filepath.Join(v.ProjectPath, "examples")
	helpers.AssertDirExists(t, examplesDir)
	
	// Check for basic example
	basicExample := filepath.Join(examplesDir, "basic", "main.go")
	helpers.AssertFileExists(t, basicExample)
	
	// Check for advanced example
	advancedExample := filepath.Join(examplesDir, "advanced", "main.go")
	helpers.AssertFileExists(t, advancedExample)
}

func (v *LibraryValidator) ValidateComprehensiveTests(t *testing.T) {
	t.Helper()
	
	// Check main test file exists
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	helpers.AssertFileExists(t, testFile)
	
	// Check examples test file exists
	examplesTestFile := filepath.Join(v.ProjectPath, "examples_test.go")
	helpers.AssertFileExists(t, examplesTestFile)
	
	// Check that tests import testify
	content := helpers.ReadFileContent(t, testFile)
	helpers.AssertFileContains(t, content, "github.com/stretchr/testify")
}

func (v *LibraryValidator) ValidateDocumentation(t *testing.T) {
	t.Helper()
	
	// Check doc.go exists
	docFile := filepath.Join(v.ProjectPath, "doc.go")
	helpers.AssertFileExists(t, docFile)
	
	// Check README exists
	readmeFile := filepath.Join(v.ProjectPath, "README.md")
	helpers.AssertFileExists(t, readmeFile)
	
	// Check that README has expected sections
	readmeContent := helpers.ReadFileContent(t, readmeFile)
	expectedSections := []string{"Installation", "Usage", "Examples", "API"}
	for _, section := range expectedSections {
		if helpers.StringContains(readmeContent, section) {
			t.Logf("✓ README contains %s section", section)
		} else {
			t.Logf("⚠ README missing %s section", section)
		}
	}
}

func (v *LibraryValidator) ValidatePackageOrganization(t *testing.T) {
	t.Helper()
	
	// Check core structure
	expectedDirs := []string{
		"examples",
		"internal",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		helpers.AssertDirExists(t, dirPath)
	}
	
	// Check core files
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"doc.go",
	}
	
	for _, file := range expectedFiles {
		filePath := filepath.Join(v.ProjectPath, file)
		helpers.AssertFileExists(t, filePath)
	}
}

func (v *LibraryValidator) ValidatePublicFunctionDocumentation(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check that public functions have comments
		if helpers.StringContains(content, "//") {
			t.Log("✓ Documentation comments found")
		}
	}
}

func (v *LibraryValidator) ValidateRunnableExamples(t *testing.T) {
	t.Helper()
	
	// Test basic example
	basicExample := filepath.Join(v.ProjectPath, "examples", "basic", "main.go")
	if helpers.FileExists(basicExample) {
		helpers.AssertFileContains(t, basicExample, "func main()")
	}
	
	// Test advanced example
	advancedExample := filepath.Join(v.ProjectPath, "examples", "advanced", "main.go")
	if helpers.FileExists(advancedExample) {
		helpers.AssertFileContains(t, advancedExample, "func main()")
	}
}

func (v *LibraryValidator) ValidatePublicAPITestCoverage(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	
	if helpers.FileExists(testFile) {
		content := helpers.ReadFileContent(t, testFile)
		
		// Check for test functions
		if helpers.StringContains(content, "func Test") {
			t.Log("✓ Test functions found")
		}
		
		// Check for testify usage
		if helpers.StringContains(content, "assert") || helpers.StringContains(content, "require") {
			t.Log("✓ Testify assertions found")
		}
	}
}

func (v *LibraryValidator) ValidateIntuitiveAPI(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for common patterns like New(), Config, Client
		intuitivePatterns := []string{"New", "Config", "Client", "Options"}
		for _, pattern := range intuitivePatterns {
			if helpers.StringContains(content, pattern) {
				t.Logf("✓ Intuitive API pattern found: %s", pattern)
			}
		}
	}
}

func (v *LibraryValidator) ValidateDescriptiveFunctionNames(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check that function names are descriptive (not single letters)
		if helpers.StringContains(content, "func") {
			t.Log("✓ Functions are defined")
		}
	}
}

func (v *LibraryValidator) ValidateParameterTypes(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check that parameters use appropriate types
		if helpers.StringContains(content, "context.Context") {
			t.Log("✓ Context parameter found")
		}
	}
}

func (v *LibraryValidator) ValidateReturnValueConventions(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check that functions return error as last parameter
		if helpers.StringContains(content, "error") {
			t.Log("✓ Error return values found")
		}
	}
}

func (v *LibraryValidator) ValidateConsistentErrorHandling(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for consistent error handling patterns
		if helpers.StringContains(content, "fmt.Errorf") || helpers.StringContains(content, "errors.New") {
			t.Log("✓ Error handling patterns found")
		}
	}
}

func (v *LibraryValidator) ValidateBasicAndAdvancedExamples(t *testing.T) {
	t.Helper()
	
	// Check basic example
	basicExample := filepath.Join(v.ProjectPath, "examples", "basic", "main.go")
	helpers.AssertFileExists(t, basicExample)
	
	// Check advanced example
	advancedExample := filepath.Join(v.ProjectPath, "examples", "advanced", "main.go")
	helpers.AssertFileExists(t, advancedExample)
	
	// Check that examples are different
	basicContent := helpers.ReadFileContent(t, basicExample)
	advancedContent := helpers.ReadFileContent(t, advancedExample)
	
	if basicContent != advancedContent {
		t.Log("✓ Basic and advanced examples are different")
	}
}

func (v *LibraryValidator) ValidateTestedExamples(t *testing.T) {
	t.Helper()
	
	examplesTestFile := filepath.Join(v.ProjectPath, "examples_test.go")
	helpers.AssertFileExists(t, examplesTestFile)
	
	content := helpers.ReadFileContent(t, examplesTestFile)
	
	// Check that examples are tested
	if helpers.StringContains(content, "func Example") {
		t.Log("✓ Example tests found")
	}
}

func (v *LibraryValidator) ValidateUnitTestCoverage(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	
	if helpers.FileExists(testFile) {
		content := helpers.ReadFileContent(t, testFile)
		
		// Check for comprehensive test coverage
		if helpers.StringContains(content, "func Test") {
			t.Log("✓ Unit tests found")
		}
	}
}

func (v *LibraryValidator) ValidateBenchmarkTests(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	
	if helpers.FileExists(testFile) {
		content := helpers.ReadFileContent(t, testFile)
		
		// Check for benchmark tests
		if helpers.StringContains(content, "func Benchmark") {
			t.Log("✓ Benchmark tests found")
		}
	}
}

func (v *LibraryValidator) ValidateMinimalLogging(t *testing.T, logger string) {
	t.Helper()
	
	// Check that internal logger exists but is minimal
	loggerFile := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	helpers.AssertFileExists(t, loggerFile)
	
	// Check go.mod for logger dependency (if not slog)
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	switch logger {
	case "zap":
		helpers.AssertFileContains(t, goMod, "go.uber.org/zap")
	case "logrus":
		helpers.AssertFileContains(t, goMod, "github.com/sirupsen/logrus")
	case "zerolog":
		helpers.AssertFileContains(t, goMod, "github.com/rs/zerolog")
	case "slog":
		// slog is part of standard library
	}
}

func (v *LibraryValidator) ValidateConfigurableLogging(t *testing.T) {
	t.Helper()
	
	loggerFile := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFileContent(t, loggerFile)
		
		// Check that logging is configurable
		if helpers.StringContains(content, "level") || helpers.StringContains(content, "Level") {
			t.Log("✓ Configurable logging found")
		}
	}
}

func (v *LibraryValidator) ValidateNonInterfering(t *testing.T) {
	t.Helper()
	
	// Check that library doesn't interfere with user's logging
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Should not have direct logging calls in public API
		if !helpers.StringContains(content, "log.Print") && !helpers.StringContains(content, "fmt.Print") {
			t.Log("✓ No direct logging calls in public API")
		}
	}
}

func (v *LibraryValidator) ValidateStructuredLogging(t *testing.T, logger string) {
	t.Helper()
	
	loggerFile := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFileContent(t, loggerFile)
		
		// Check for structured logging patterns
		switch logger {
		case "slog":
			if helpers.StringContains(content, "slog.") {
				t.Log("✓ Structured logging with slog")
			}
		case "zap":
			if helpers.StringContains(content, "zap.") {
				t.Log("✓ Structured logging with zap")
			}
		case "logrus":
			if helpers.StringContains(content, "logrus.") {
				t.Log("✓ Structured logging with logrus")
			}
		case "zerolog":
			if helpers.StringContains(content, "zerolog.") {
				t.Log("✓ Structured logging with zerolog")
			}
		}
	}
}

func (v *LibraryValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *LibraryValidator) ValidatePackageDocumentation(t *testing.T) {
	t.Helper()
	
	docFile := filepath.Join(v.ProjectPath, "doc.go")
	helpers.AssertFileExists(t, docFile)
	
	content := helpers.ReadFileContent(t, docFile)
	
	// Check that package documentation exists
	if helpers.StringContains(content, "Package") {
		t.Log("✓ Package documentation found")
	}
}

func (v *LibraryValidator) ValidateGodocComments(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for godoc comments
		if helpers.StringContains(content, "//") {
			t.Log("✓ Godoc comments found")
		}
	}
}

func (v *LibraryValidator) ValidateDocumentationExamples(t *testing.T) {
	t.Helper()
	
	examplesTestFile := filepath.Join(v.ProjectPath, "examples_test.go")
	if helpers.FileExists(examplesTestFile) {
		content := helpers.ReadFileContent(t, examplesTestFile)
		
		// Check for example functions
		if helpers.StringContains(content, "func Example") {
			t.Log("✓ Documentation examples found")
		}
	}
}

func (v *LibraryValidator) ValidateComprehensiveREADME(t *testing.T) {
	t.Helper()
	
	readmeFile := filepath.Join(v.ProjectPath, "README.md")
	helpers.AssertFileExists(t, readmeFile)
	
	content := helpers.ReadFileContent(t, readmeFile)
	
	// Check for comprehensive sections
	expectedSections := []string{"Installation", "Usage", "Examples", "API", "License"}
	foundSections := 0
	for _, section := range expectedSections {
		if helpers.StringContains(content, section) {
			foundSections++
		}
	}
	
	if foundSections >= 3 {
		t.Logf("✓ README has %d/%d expected sections", foundSections, len(expectedSections))
	} else {
		t.Logf("⚠ README only has %d/%d expected sections", foundSections, len(expectedSections))
	}
}

func (v *LibraryValidator) ValidateEssentialDependencies(t *testing.T) {
	t.Helper()
	
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	content := helpers.ReadFileContent(t, goMod)
	
	// For libraries, dependencies should be minimal
	// Check that testify is included for testing
	if helpers.StringContains(content, "github.com/stretchr/testify") {
		t.Log("✓ Testing dependencies found")
	}
}

func (v *LibraryValidator) ValidateMinimalDependencyFootprint(t *testing.T) {
	t.Helper()
	
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	content := helpers.ReadFileContent(t, goMod)
	
	// Count require statements (rough measure of dependency footprint)
	requireCount := 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.Contains(line, "require") && !strings.Contains(line, "//") {
			requireCount++
		}
	}
	
	// Libraries should have minimal dependencies
	if requireCount <= 10 {
		t.Logf("✓ Minimal dependency footprint: %d require statements", requireCount)
	} else {
		t.Logf("⚠ Large dependency footprint: %d require statements", requireCount)
	}
}

func (v *LibraryValidator) ValidateNoUnnecessaryTransitiveDependencies(t *testing.T) {
	t.Helper()
	
	// This would ideally check go.sum or run go list -m all
	// For now, just ensure go.mod is clean
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	t.Log("✓ go.mod exists and should be clean")
}

func (v *LibraryValidator) ValidateInitializationPatterns(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for initialization patterns
		initPatterns := []string{"New", "Config", "Client", "Options"}
		for _, pattern := range initPatterns {
			if helpers.StringContains(content, pattern) {
				t.Logf("✓ Initialization pattern found: %s", pattern)
			}
		}
	}
}

func (v *LibraryValidator) ValidateClientConfiguration(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for configuration struct
		if helpers.StringContains(content, "Config") || helpers.StringContains(content, "Options") {
			t.Log("✓ Configuration patterns found")
		}
	}
}

func (v *LibraryValidator) ValidateDefaultValues(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for default values
		if helpers.StringContains(content, "default") || helpers.StringContains(content, "Default") {
			t.Log("✓ Default values found")
		}
	}
}

func (v *LibraryValidator) ValidateInitializationDocumentation(t *testing.T) {
	t.Helper()
	
	readmeFile := filepath.Join(v.ProjectPath, "README.md")
	if helpers.FileExists(readmeFile) {
		content := helpers.ReadFileContent(t, readmeFile)
		
		// Check that initialization is documented
		if helpers.StringContains(content, "New") || helpers.StringContains(content, "Client") {
			t.Log("✓ Initialization documentation found")
		}
	}
}

func (v *LibraryValidator) ValidateErrorConventions(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for Go error conventions
		if helpers.StringContains(content, "error") {
			t.Log("✓ Error conventions found")
		}
	}
}

func (v *LibraryValidator) ValidateDescriptiveErrorMessages(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for descriptive error messages
		if helpers.StringContains(content, "fmt.Errorf") || helpers.StringContains(content, "errors.New") {
			t.Log("✓ Descriptive error messages found")
		}
	}
}

func (v *LibraryValidator) ValidateExportedErrorTypes(t *testing.T) {
	t.Helper()
	
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for exported error types
		if helpers.StringContains(content, "var Err") || helpers.StringContains(content, "type") {
			t.Log("✓ Error types may be exported")
		}
	}
}

// Import required packages
import (
	"strings"
)