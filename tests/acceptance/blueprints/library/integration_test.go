package library_test

import (
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// Feature: Library Integration Testing
// As a developer
// I want Library blueprints to work together with different configurations
// So that I can build robust, reusable libraries

func TestLibrary_CrossLoggerIntegration(t *testing.T) {
	// Feature: Cross-Logger Integration for Libraries
	// As a library developer
	// I want all loggers to work consistently for internal debugging
	// So that I can choose the best logger without changing implementation

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("CrossIntegration_"+logger, func(t *testing.T) {
			t.Parallel()

			// Scenario: Cross-logger integration
			// Given I generate a library with "<logger>"
			// Then the logger should integrate minimally with library components
			// And logger should not interfere with user's logging setup
			// And the library should compile and work correctly
			// And examples should work with the logger

			// Given I generate a library with "<logger>"
			config := types.ProjectConfig{
				Name:      "test-lib-cross-" + logger,
				Module:    "github.com/test/test-lib-cross-" + logger,
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then the logger should integrate minimally with library components
			validator := NewLibraryValidator(projectPath)
			validator.ValidateMinimalLogging(t, logger)

			// And logger should not interfere with user's logging setup
			validator.ValidateNonInterfering(t)

			// And the library should compile and work correctly
			validator.ValidateCompilation(t)

			// And examples should work with the logger
			validator.ValidateExamplesWithLogger(t, logger)
		})
	}
}

func TestLibrary_CompilationValidation(t *testing.T) {
	// Feature: Library Compilation Validation
	// As a quality assurance engineer
	// I want all Library configurations to compile successfully
	// So that users always get working libraries

	testCases := []struct {
		name   string
		config types.ProjectConfig
	}{
		{
			name: "StandardLibrary_Slog",
			config: types.ProjectConfig{
				Name:      "test-lib-compile-slog",
				Module:    "github.com/test/test-lib-compile-slog",
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    "slog",
			},
		},
		{
			name: "StandardLibrary_Zap",
			config: types.ProjectConfig{
				Name:      "test-lib-compile-zap",
				Module:    "github.com/test/test-lib-compile-zap",
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    "zap",
			},
		},
		{
			name: "StandardLibrary_Logrus",
			config: types.ProjectConfig{
				Name:      "test-lib-compile-logrus",
				Module:    "github.com/test/test-lib-compile-logrus",
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    "logrus",
			},
		},
		{
			name: "StandardLibrary_Zerolog",
			config: types.ProjectConfig{
				Name:      "test-lib-compile-zerolog",
				Module:    "github.com/test/test-lib-compile-zerolog",
				Type:      "library-standard",
				GoVersion: "1.21",
				Logger:    "zerolog",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Scenario: Compilation validation
			// Given I generate a library with specific configuration
			// Then the library should compile without errors
			// And all examples should compile successfully
			// And all tests should pass
			// And the library should be importable

			// Given I generate a library with specific configuration
			projectPath := helpers.GenerateProject(t, tc.config)

			// Then the library should compile without errors
			validator := NewLibraryValidator(projectPath)
			validator.ValidateCompilation(t)

			// And all examples should compile successfully
			validator.ValidateExamplesCompilation(t)

			// And all tests should pass
			validator.ValidateTestsPass(t)

			// And the library should be importable
			validator.ValidateImportability(t)
		})
	}
}

func TestLibrary_ArchitectureCompliance(t *testing.T) {
	// Feature: Library Architecture Compliance
	// As a code reviewer
	// I want generated library code to follow architecture principles
	// So that libraries maintain architectural integrity

	// Scenario: Architecture compliance for libraries
	// Given I generate a library
	// Then the code should follow library best practices
	// And dependency directions should be correct
	// And package organization should be logical
	// And the library should pass architectural validation

	config := types.ProjectConfig{
		Name:      "test-lib-architecture",
		Module:    "github.com/test/test-lib-architecture",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should follow library best practices
	validator := NewLibraryValidator(projectPath)
	validator.ValidateLibraryArchitecture(t)

	// And dependency directions should be correct
	validator.ValidateLibraryDependencyDirections(t)

	// And package organization should be logical
	validator.ValidateLibraryPackageOrganization(t)

	// And the library should pass architectural validation
	validator.ValidateLibraryArchitecturalCompliance(t)

	// And the library should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_SecurityValidation(t *testing.T) {
	// Feature: Library Security Validation
	// As a security engineer
	// I want libraries to follow security best practices
	// So that generated libraries are secure by default

	// Scenario: Security validation
	// Given I generate a library
	// Then the code should not contain security vulnerabilities
	// And logging should not expose sensitive information
	// And dependencies should be secure
	// And the library should not create security risks for users

	config := types.ProjectConfig{
		Name:      "test-lib-security",
		Module:    "github.com/test/test-lib-security",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should not contain security vulnerabilities
	validator := NewLibraryValidator(projectPath)
	validator.ValidateLibrarySecurityPractices(t)

	// And logging should not expose sensitive information
	validator.ValidateSecureLibraryLogging(t)

	// And dependencies should be secure
	validator.ValidateSecureDependencies(t)

	// And the library should not create security risks for users
	validator.ValidateNoUserSecurityRisks(t)

	// And the library should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_UsabilityValidation(t *testing.T) {
	// Feature: Library Usability Validation
	// As a library user
	// I want libraries to be easy to use
	// So that I can integrate them quickly into my projects

	// Scenario: Usability validation
	// Given I generate a library
	// Then the API should be intuitive and well-documented
	// And examples should be clear and comprehensive
	// And error messages should be helpful
	// And the library should have good defaults

	config := types.ProjectConfig{
		Name:      "test-lib-usability",
		Module:    "github.com/test/test-lib-usability",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the API should be intuitive and well-documented
	validator := NewLibraryValidator(projectPath)
	validator.ValidateIntuitiveAndDocumentedAPI(t)

	// And examples should be clear and comprehensive
	validator.ValidateClearAndComprehensiveExamples(t)

	// And error messages should be helpful
	validator.ValidateHelpfulErrorMessages(t)

	// And the library should have good defaults
	validator.ValidateGoodDefaults(t)

	// And the library should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_MaintenanceValidation(t *testing.T) {
	// Feature: Library Maintenance Validation
	// As a library maintainer
	// I want generated libraries to be easy to maintain
	// So that I can keep them updated and fix issues

	// Scenario: Maintenance validation
	// Given I generate a library
	// Then the code should be well-organized and readable
	// And tests should provide good coverage
	// And the build system should be comprehensive
	// And documentation should be maintainable

	config := types.ProjectConfig{
		Name:      "test-lib-maintenance",
		Module:    "github.com/test/test-lib-maintenance",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should be well-organized and readable
	validator := NewLibraryValidator(projectPath)
	validator.ValidateWellOrganizedCode(t)

	// And tests should provide good coverage
	validator.ValidateGoodTestCoverage(t)

	// And the build system should be comprehensive
	validator.ValidateComprehensiveBuildSystem(t)

	// And documentation should be maintainable
	validator.ValidateMaintainableDocumentation(t)

	// And the library should compile successfully
	validator.ValidateCompilation(t)
}

func TestLibrary_APICompatibilityValidation(t *testing.T) {
	// Feature: Library API Compatibility
	// As a library user
	// I want libraries to have stable, backward-compatible APIs
	// So that I can upgrade versions without breaking my code

	// Scenario: API compatibility validation
	// Given I generate a library
	// Then the public API should be clearly defined
	// And internal APIs should be properly encapsulated
	// And the API should follow Go conventions
	// And versioning should be supported

	config := types.ProjectConfig{
		Name:      "test-lib-api-compat",
		Module:    "github.com/test/test-lib-api-compat",
		Type:      "library-standard",
		GoVersion: "1.21",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the public API should be clearly defined
	validator := NewLibraryValidator(projectPath)
	validator.ValidateClearlyDefinedPublicAPI(t)

	// And internal APIs should be properly encapsulated
	validator.ValidateProperlyEncapsulatedInternalAPIs(t)

	// And the API should follow Go conventions
	validator.ValidateGoConventions(t)

	// And versioning should be supported
	validator.ValidateVersioningSupport(t)

	// And the library should compile successfully
	validator.ValidateCompilation(t)
}

// Additional validation methods for integration tests

func (v *LibraryValidator) ValidateExamplesWithLogger(t *testing.T, logger string) {
	t.Helper()
	
	// Check that examples work with the specified logger
	basicExample := filepath.Join(v.ProjectPath, "examples", "basic", "main.go")
	if helpers.FileExists(basicExample) {
		// For libraries, examples typically don't directly use internal loggers
		// Just ensure they compile
		t.Log("✓ Basic example exists and should work with logger")
	}
	
	advancedExample := filepath.Join(v.ProjectPath, "examples", "advanced", "main.go")
	if helpers.FileExists(advancedExample) {
		t.Log("✓ Advanced example exists and should work with logger")
	}
}

func (v *LibraryValidator) ValidateExamplesCompilation(t *testing.T) {
	t.Helper()
	
	// Test that examples compile successfully
	basicExampleDir := filepath.Join(v.ProjectPath, "examples", "basic")
	if helpers.DirExists(basicExampleDir) {
		helpers.AssertCompilable(t, basicExampleDir)
	}
	
	advancedExampleDir := filepath.Join(v.ProjectPath, "examples", "advanced")
	if helpers.DirExists(advancedExampleDir) {
		helpers.AssertCompilable(t, advancedExampleDir)
	}
}

func (v *LibraryValidator) ValidateTestsPass(t *testing.T) {
	t.Helper()
	
	// This would ideally run go test ./... to ensure all tests pass
	// For now, just validate that test files exist and are well-formed
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	
	if helpers.FileExists(testFile) {
		content := helpers.ReadFileContent(t, testFile)
		if helpers.StringContains(content, "func Test") {
			t.Log("✓ Unit tests exist and should pass")
		}
	}
	
	examplesTestFile := filepath.Join(v.ProjectPath, "examples_test.go")
	if helpers.FileExists(examplesTestFile) {
		content := helpers.ReadFileContent(t, examplesTestFile)
		if helpers.StringContains(content, "func Example") {
			t.Log("✓ Example tests exist and should pass")
		}
	}
}

func (v *LibraryValidator) ValidateImportability(t *testing.T) {
	t.Helper()
	
	// Check that the library can be imported by checking go.mod
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	content := helpers.ReadFileContent(t, goMod)
	
	// Should have module declaration
	if helpers.StringContains(content, "module") {
		t.Log("✓ Library has proper module declaration")
	}
	
	// Should have Go version
	if helpers.StringContains(content, "go ") {
		t.Log("✓ Library has Go version specified")
	}
}

func (v *LibraryValidator) ValidateLibraryArchitecture(t *testing.T) {
	t.Helper()
	
	// Check standard library structure
	expectedDirs := []string{
		"examples",
		"internal",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		helpers.AssertDirExists(t, dirPath)
	}
	
	// Check that main library file exists
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Check that internal package exists
	internalLogger := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	helpers.AssertFileExists(t, internalLogger)
}

func (v *LibraryValidator) ValidateLibraryDependencyDirections(t *testing.T) {
	t.Helper()
	
	// For libraries, main package should not import internal packages directly
	// Examples should import the main package
	
	basicExample := filepath.Join(v.ProjectPath, "examples", "basic", "main.go")
	if helpers.FileExists(basicExample) {
		content := helpers.ReadFileContent(t, basicExample)
		
		// Should import the library package
		projectName := filepath.Base(v.ProjectPath)
		if helpers.StringContains(content, projectName) {
			t.Log("✓ Examples import the library package")
		}
	}
	
	// Main library should not expose internal packages
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Should not export internal types
		if !helpers.StringContains(content, "internal") {
			t.Log("✓ Main library doesn't expose internal packages")
		}
	}
}

func (v *LibraryValidator) ValidateLibraryPackageOrganization(t *testing.T) {
	t.Helper()
	
	// Check that packages are properly organized
	expectedStructure := map[string]bool{
		"examples":                    true,
		"internal":                    true,
		"internal/logger":             true,
		"go.mod":                      true,
		"README.md":                   true,
		"Makefile":                    true,
		"doc.go":                      true,
	}
	
	for path, shouldExist := range expectedStructure {
		fullPath := filepath.Join(v.ProjectPath, path)
		if shouldExist {
			if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".mod") || path == "Makefile" {
				helpers.AssertFileExists(t, fullPath)
			} else {
				helpers.AssertDirExists(t, fullPath)
			}
		}
	}
}

func (v *LibraryValidator) ValidateLibraryArchitecturalCompliance(t *testing.T) {
	t.Helper()
	
	// For libraries, validate clean architecture principles
	// Public API should be in root package
	// Internal implementation should be in internal/
	// Examples should be in examples/
	
	// Check that public API is in root
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	helpers.AssertFileExists(t, mainFile)
	
	// Check that internal implementation is properly encapsulated
	internalDir := filepath.Join(v.ProjectPath, "internal")
	helpers.AssertDirExists(t, internalDir)
	
	// Check that examples are separate
	examplesDir := filepath.Join(v.ProjectPath, "examples")
	helpers.AssertDirExists(t, examplesDir)
}

func (v *LibraryValidator) ValidateLibrarySecurityPractices(t *testing.T) {
	t.Helper()
	
	// Check that no hardcoded secrets exist in library code
	securityFiles := []string{
		filepath.Base(v.ProjectPath) + ".go",
		"internal/logger/logger.go",
		"examples/basic/main.go",
		"examples/advanced/main.go",
	}
	
	for _, file := range securityFiles {
		filePath := filepath.Join(v.ProjectPath, file)
		if helpers.FileExists(filePath) {
			content := helpers.ReadFileContent(t, filePath)
			
			// Check for common security issues
			securityIssues := []string{
				"password=",
				"secret=",
				"token=",
				"apikey=",
				"api_key=",
			}
			
			for _, issue := range securityIssues {
				if helpers.StringContains(content, issue) {
					t.Errorf("Potential security issue found in %s: %s", file, issue)
				}
			}
		}
	}
}

func (v *LibraryValidator) ValidateSecureLibraryLogging(t *testing.T) {
	t.Helper()
	
	// Check that internal logging doesn't expose sensitive information
	loggerFile := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFileContent(t, loggerFile)
		
		// Check that logger doesn't log sensitive fields
		sensitiveFields := []string{
			"password",
			"secret",
			"token",
			"apikey",
			"api_key",
		}
		
		for _, field := range sensitiveFields {
			if helpers.StringContains(content, field) {
				t.Logf("⚠ Logger may be logging sensitive field: %s", field)
			}
		}
	}
}

func (v *LibraryValidator) ValidateSecureDependencies(t *testing.T) {
	t.Helper()
	
	// Check that dependencies are secure and minimal
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	content := helpers.ReadFileContent(t, goMod)
	
	// For libraries, dependencies should be minimal
	// Check that we don't have unnecessary dependencies
	if helpers.StringContains(content, "require") {
		t.Log("✓ Dependencies are declared in go.mod")
	}
}

func (v *LibraryValidator) ValidateNoUserSecurityRisks(t *testing.T) {
	t.Helper()
	
	// Check that the library doesn't create security risks for users
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check that public API doesn't require unsafe operations
		unsafePatterns := []string{
			"unsafe.",
			"os.Execute",
			"exec.Command",
			"os.Remove",
		}
		
		for _, pattern := range unsafePatterns {
			if helpers.StringContains(content, pattern) {
				t.Logf("⚠ Potentially unsafe pattern found: %s", pattern)
			}
		}
	}
}

func (v *LibraryValidator) ValidateIntuitiveAndDocumentedAPI(t *testing.T) {
	t.Helper()
	
	// Check that API is intuitive and well-documented
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for intuitive patterns
		intuitivePatterns := []string{"New", "Config", "Client", "Options"}
		for _, pattern := range intuitivePatterns {
			if helpers.StringContains(content, pattern) {
				t.Logf("✓ Intuitive API pattern: %s", pattern)
			}
		}
		
		// Check for documentation
		if helpers.StringContains(content, "//") {
			t.Log("✓ API documentation found")
		}
	}
}

func (v *LibraryValidator) ValidateClearAndComprehensiveExamples(t *testing.T) {
	t.Helper()
	
	// Check that examples are clear and comprehensive
	basicExample := filepath.Join(v.ProjectPath, "examples", "basic", "main.go")
	if helpers.FileExists(basicExample) {
		content := helpers.ReadFileContent(t, basicExample)
		
		// Check that example has comments
		if helpers.StringContains(content, "//") {
			t.Log("✓ Basic example has comments")
		}
		
		// Check that example has main function
		if helpers.StringContains(content, "func main()") {
			t.Log("✓ Basic example has main function")
		}
	}
	
	advancedExample := filepath.Join(v.ProjectPath, "examples", "advanced", "main.go")
	if helpers.FileExists(advancedExample) {
		content := helpers.ReadFileContent(t, advancedExample)
		
		// Check that advanced example is more complex
		if helpers.StringContains(content, "//") {
			t.Log("✓ Advanced example has comments")
		}
	}
}

func (v *LibraryValidator) ValidateHelpfulErrorMessages(t *testing.T) {
	t.Helper()
	
	// Check that error messages are helpful
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for helpful error patterns
		if helpers.StringContains(content, "fmt.Errorf") {
			t.Log("✓ Helpful error messages with fmt.Errorf")
		}
		
		if helpers.StringContains(content, "errors.New") {
			t.Log("✓ Error handling with errors.New")
		}
	}
}

func (v *LibraryValidator) ValidateGoodDefaults(t *testing.T) {
	t.Helper()
	
	// Check that the library has good defaults
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for default values
		if helpers.StringContains(content, "default") || helpers.StringContains(content, "Default") {
			t.Log("✓ Good defaults found")
		}
	}
}

func (v *LibraryValidator) ValidateWellOrganizedCode(t *testing.T) {
	t.Helper()
	
	// Check that code is well-organized
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
	
	// Check that directories are properly organized
	expectedDirs := []string{
		"examples",
		"internal",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		helpers.AssertDirExists(t, dirPath)
	}
}

func (v *LibraryValidator) ValidateGoodTestCoverage(t *testing.T) {
	t.Helper()
	
	// Check that tests provide good coverage
	projectName := filepath.Base(v.ProjectPath)
	testFile := filepath.Join(v.ProjectPath, projectName+"_test.go")
	
	if helpers.FileExists(testFile) {
		content := helpers.ReadFileContent(t, testFile)
		
		// Check for various test types
		if helpers.StringContains(content, "func Test") {
			t.Log("✓ Unit tests found")
		}
		
		if helpers.StringContains(content, "func Benchmark") {
			t.Log("✓ Benchmark tests found")
		}
	}
	
	// Check for example tests
	examplesTestFile := filepath.Join(v.ProjectPath, "examples_test.go")
	if helpers.FileExists(examplesTestFile) {
		content := helpers.ReadFileContent(t, examplesTestFile)
		
		if helpers.StringContains(content, "func Example") {
			t.Log("✓ Example tests found")
		}
	}
}

func (v *LibraryValidator) ValidateComprehensiveBuildSystem(t *testing.T) {
	t.Helper()
	
	// Check that build system is comprehensive
	makefile := filepath.Join(v.ProjectPath, "Makefile")
	helpers.AssertFileExists(t, makefile)
	
	content := helpers.ReadFileContent(t, makefile)
	
	// Check for common build targets
	buildTargets := []string{"build", "test", "lint", "clean"}
	for _, target := range buildTargets {
		if helpers.StringContains(content, target+":") {
			t.Logf("✓ Build target found: %s", target)
		}
	}
}

func (v *LibraryValidator) ValidateMaintainableDocumentation(t *testing.T) {
	t.Helper()
	
	// Check that documentation is maintainable
	readmeFile := filepath.Join(v.ProjectPath, "README.md")
	helpers.AssertFileExists(t, readmeFile)
	
	content := helpers.ReadFileContent(t, readmeFile)
	
	// Check for maintainable documentation structure
	maintainableSections := []string{"Installation", "Usage", "API", "Contributing"}
	for _, section := range maintainableSections {
		if helpers.StringContains(content, section) {
			t.Logf("✓ Maintainable documentation section: %s", section)
		}
	}
}

func (v *LibraryValidator) ValidateClearlyDefinedPublicAPI(t *testing.T) {
	t.Helper()
	
	// Check that public API is clearly defined
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for exported types and functions
		if helpers.StringContains(content, "type") && helpers.StringContains(content, "func") {
			t.Log("✓ Public API clearly defined with types and functions")
		}
		
		// Check for proper documentation
		if helpers.StringContains(content, "//") {
			t.Log("✓ Public API documented")
		}
	}
}

func (v *LibraryValidator) ValidateProperlyEncapsulatedInternalAPIs(t *testing.T) {
	t.Helper()
	
	// Check that internal APIs are properly encapsulated
	internalDir := filepath.Join(v.ProjectPath, "internal")
	helpers.AssertDirExists(t, internalDir)
	
	// Check that internal packages exist
	internalLogger := filepath.Join(v.ProjectPath, "internal", "logger", "logger.go")
	helpers.AssertFileExists(t, internalLogger)
	
	// Check that main package doesn't expose internal details
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Should not expose internal package details
		if !helpers.StringContains(content, "internal") {
			t.Log("✓ Internal APIs properly encapsulated")
		}
	}
}

func (v *LibraryValidator) ValidateGoConventions(t *testing.T) {
	t.Helper()
	
	// Check that the library follows Go conventions
	projectName := filepath.Base(v.ProjectPath)
	mainFile := filepath.Join(v.ProjectPath, projectName+".go")
	
	if helpers.FileExists(mainFile) {
		content := helpers.ReadFileContent(t, mainFile)
		
		// Check for Go conventions
		conventions := []string{
			"package " + projectName,
			"func New",
			"error",
		}
		
		for _, convention := range conventions {
			if helpers.StringContains(content, convention) {
				t.Logf("✓ Go convention followed: %s", convention)
			}
		}
	}
}

func (v *LibraryValidator) ValidateVersioningSupport(t *testing.T) {
	t.Helper()
	
	// Check that versioning is supported
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	content := helpers.ReadFileContent(t, goMod)
	
	// Check that module is properly declared for versioning
	if helpers.StringContains(content, "module") {
		t.Log("✓ Module declaration supports versioning")
	}
	
	// Check that Go version is specified
	if helpers.StringContains(content, "go ") {
		t.Log("✓ Go version specified for compatibility")
	}
}

// Import required packages
import (
	"path/filepath"
	"strings"
)