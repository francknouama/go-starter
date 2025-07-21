package library

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/suite"
)

// LibraryAcceptanceTestSuite provides acceptance tests for library blueprints
type LibraryAcceptanceTestSuite struct {
	suite.Suite
	ctx *LibraryTestContext
}

// SetupSuite runs before all tests in the suite
func (suite *LibraryAcceptanceTestSuite) SetupSuite() {
	suite.ctx = InitializeLibraryContext()
}

// TearDownSuite runs after all tests in the suite
func (suite *LibraryAcceptanceTestSuite) TearDownSuite() {
	if suite.ctx != nil {
		suite.ctx.cleanup()
	}
}

// TestLibraryBlueprintGeneration tests library blueprint generation with godog
func (suite *LibraryAcceptanceTestSuite) TestLibraryBlueprintGeneration() {
	godogSuite := godog.TestSuite{
		ScenarioInitializer: suite.InitializeLibraryScenarios,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: suite.T(),
		},
	}

	if godogSuite.Run() != 0 {
		suite.T().Fatal("Non-zero status returned, failed to run library acceptance tests")
	}
}

// InitializeLibraryScenarios registers step definitions for library scenarios
func (suite *LibraryAcceptanceTestSuite) InitializeLibraryScenarios(sc *godog.ScenarioContext) {
	ctx := suite.ctx

	// Background steps
	sc.Given(`^the go-starter CLI tool is available$`, func() error {
		// Check if go-starter is available
		cmd := exec.Command("go-starter", "--help")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("go-starter CLI not available: %s", err.Error())
		}
		return nil
	})

	sc.Given(`^I am in a clean working directory$`, func() error {
		// Ensure we're in a clean state
		ctx.cleanup()
		return nil
	})

	// Library generation steps
	sc.Given(`^I want to create a reusable Go library$`, ctx.iWantToCreateAReusableGoLibrary)
	sc.When(`^I run the command "([^"]*)"$`, ctx.iRunTheCommand)
	sc.Then(`^the generation should succeed$`, ctx.theGenerationShouldSucceed)
	sc.Then(`^the project should contain all essential library components$`, ctx.theProjectShouldContainAllEssentialLibraryComponents)
	sc.Then(`^the generated code should compile successfully$`, ctx.theGeneratedCodeShouldCompileSuccessfully)
	sc.Then(`^the library should follow Go package best practices$`, ctx.theLibraryShouldFollowGoPackageBestPractices)
	sc.Then(`^the library should include comprehensive documentation$`, ctx.theLibraryShouldIncludeComprehensiveDocumentation)

	// Logging implementation steps
	sc.Given(`^I want to create libraries with various logging options$`, func() error {
		ctx.scenarios["logging_test"] = true
		return nil
	})
	sc.When(`^I generate a library with logger "([^"]*)"$`, ctx.iGenerateALibraryWithLogger)
	sc.Then(`^the project should support the "([^"]*)" logging interface$`, ctx.theProjectShouldSupportTheLoggingInterface)
	sc.Then(`^the library should use dependency injection for logging$`, ctx.theLibraryShouldUseDependencyInjectionForLogging)
	sc.Then(`^the logger should be optional and not forced on consumers$`, ctx.theLoggerShouldBeOptionalAndNotForcedOnConsumers)
	sc.Then(`^the library should compile without the logger dependency$`, ctx.theLibraryShouldCompileWithoutTheLoggerDependency)

	// Documentation steps
	sc.Given(`^I want a well-documented library$`, func() error {
		ctx.scenarios["documentation_focus"] = true
		return nil
	})
	sc.When(`^I generate a library with comprehensive documentation$`, ctx.iGenerateALibraryWithComprehensiveDocumentation)
	sc.Then(`^the project should include a detailed README$`, ctx.theProjectShouldIncludeADetailedREADME)
	sc.Then(`^the project should include a package documentation file$`, ctx.theProjectShouldIncludeAPackageDocumentationFile)
	sc.Then(`^the project should include usage examples$`, ctx.theProjectShouldIncludeUsageExamples)
	sc.Then(`^the examples should be executable and testable$`, ctx.theExamplesShouldBeExecutableAndTestable)
	sc.Then(`^the documentation should follow Go documentation standards$`, ctx.theDocumentationShouldFollowGoDocumentationStandards)

	// Testing structure steps
	sc.Given(`^I want a thoroughly tested library$`, func() error {
		ctx.scenarios["testing_focus"] = true
		return nil
	})
	sc.When(`^I generate a library with test infrastructure$`, ctx.iGenerateALibraryWithTestInfrastructure)
	sc.Then(`^the project should include unit tests$`, ctx.theProjectShouldIncludeUnitTests)
	sc.Then(`^the project should include example tests$`, ctx.theProjectShouldIncludeExampleTests)
	sc.Then(`^the project should include benchmark tests$`, ctx.theProjectShouldIncludeBenchmarkTests)
	sc.Then(`^the project should include test coverage configuration$`, ctx.theProjectShouldIncludeTestCoverageConfiguration)
	sc.Then(`^the tests should follow Go testing conventions$`, ctx.theTestsShouldFollowGoTestingConventions)

	// API design pattern steps
	sc.Given(`^I want a library with clean API design$`, func() error {
		ctx.scenarios["api_design"] = true
		return nil
	})
	sc.When(`^I generate a library following best practices$`, ctx.iGenerateALibraryFollowingBestPractices)
	sc.Then(`^the library should use functional options pattern$`, ctx.theLibraryShouldUseFunctionalOptionsPattern)
	sc.Then(`^the library should have minimal public API surface$`, ctx.theLibraryShouldHaveMinimalPublicAPISurface)
	sc.Then(`^the library should provide clear error types$`, ctx.theLibraryShouldProvideClearErrorTypes)
	sc.Then(`^the library should support context for cancellation$`, ctx.theLibraryShouldSupportContextForCancellation)
	sc.Then(`^the library should be thread-safe$`, ctx.theLibraryShouldBeThreadSafe)

	// Additional step definitions would be added here for remaining scenarios...

	// Scenario hooks - commented out due to version compatibility
	// sc.Before(ctx.beforeScenario)
	// sc.After(ctx.afterScenario)
}

// TestBasicLibraryGeneration tests basic library generation
func (suite *LibraryAcceptanceTestSuite) TestBasicLibraryGeneration() {
	ctx := suite.ctx
	
	// Test basic library generation
	err := ctx.iRunTheCommand("go-starter new test-library --type=library-standard --module=github.com/test/library --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldContainAllEssentialLibraryComponents()
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
	
	// Verify library structure
	suite.Assert().True(len(ctx.generatedFiles) > 15, "Should generate at least 15 files")
	suite.Assert().Contains(ctx.generatedFiles, "library.go")
	suite.Assert().Contains(ctx.generatedFiles, "library_test.go")
	suite.Assert().Contains(ctx.generatedFiles, "doc.go")
	suite.Assert().Contains(ctx.generatedFiles, "examples/basic/main.go")
	suite.Assert().Contains(ctx.generatedFiles, "examples/advanced/main.go")
}

// TestLoggerIntegration tests different logger integrations
func (suite *LibraryAcceptanceTestSuite) TestLoggerIntegration() {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}
	
	for _, logger := range loggers {
		suite.Run(fmt.Sprintf("Logger_%s", logger), func() {
			ctx := InitializeLibraryContext()
			defer ctx.cleanup()
			
			err := ctx.iGenerateALibraryWithLogger(logger)
			suite.Require().NoError(err)
			
			err = ctx.theGenerationShouldSucceed()
			suite.Require().NoError(err)
			
			err = ctx.theProjectShouldSupportTheLoggingInterface(logger)
			suite.Require().NoError(err)
			
			err = ctx.theLibraryShouldUseDependencyInjectionForLogging()
			suite.Require().NoError(err)
			
			err = ctx.theLoggerShouldBeOptionalAndNotForcedOnConsumers()
			suite.Require().NoError(err)
			
			err = ctx.theGeneratedCodeShouldCompileSuccessfully()
			suite.Require().NoError(err)
		})
	}
}

// TestDocumentationQuality tests documentation generation
func (suite *LibraryAcceptanceTestSuite) TestDocumentationQuality() {
	ctx := suite.ctx
	
	err := ctx.iGenerateALibraryWithComprehensiveDocumentation()
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Test documentation components
	err = ctx.theProjectShouldIncludeADetailedREADME()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldIncludeAPackageDocumentationFile()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldIncludeUsageExamples()
	suite.Require().NoError(err)
	
	err = ctx.theExamplesShouldBeExecutableAndTestable()
	suite.Require().NoError(err)
	
	err = ctx.theDocumentationShouldFollowGoDocumentationStandards()
	suite.Require().NoError(err)
}

// TestTestingInfrastructure tests the testing setup
func (suite *LibraryAcceptanceTestSuite) TestTestingInfrastructure() {
	ctx := suite.ctx
	
	err := ctx.iGenerateALibraryWithTestInfrastructure()
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Test all testing components
	err = ctx.theProjectShouldIncludeUnitTests()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldIncludeExampleTests()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldIncludeBenchmarkTests()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldIncludeTestCoverageConfiguration()
	suite.Require().NoError(err)
	
	err = ctx.theTestsShouldFollowGoTestingConventions()
	suite.Require().NoError(err)
	
	// Verify test execution time
	suite.Assert().Less(ctx.testTime.Seconds(), float64(30), "Tests should complete within 30 seconds")
}

// TestAPIDesignPatterns tests API design best practices
func (suite *LibraryAcceptanceTestSuite) TestAPIDesignPatterns() {
	ctx := suite.ctx
	
	err := ctx.iGenerateALibraryFollowingBestPractices()
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Test API design patterns
	err = ctx.theLibraryShouldUseFunctionalOptionsPattern()
	suite.Require().NoError(err)
	
	err = ctx.theLibraryShouldHaveMinimalPublicAPISurface()
	suite.Require().NoError(err)
	
	err = ctx.theLibraryShouldProvideClearErrorTypes()
	suite.Require().NoError(err)
	
	err = ctx.theLibraryShouldSupportContextForCancellation()
	suite.Require().NoError(err)
	
	err = ctx.theLibraryShouldBeThreadSafe()
	suite.Require().NoError(err)
}

// TestLicenseOptions tests different license configurations
func (suite *LibraryAcceptanceTestSuite) TestLicenseOptions() {
	licenses := []string{"MIT", "Apache-2.0", "GPL-3.0", "BSD-3-Clause"}
	
	for _, license := range licenses {
		suite.Run(fmt.Sprintf("License_%s", license), func() {
			ctx := InitializeLibraryContext()
			defer ctx.cleanup()
			
			command := fmt.Sprintf("go-starter new test-license-%s --type=library-standard --license=%s --module=github.com/test/license-%s --no-git",
				strings.ToLower(license), license, strings.ToLower(license))
			
			err := ctx.iRunTheCommand(command)
			suite.Require().NoError(err)
			
			err = ctx.theGenerationShouldSucceed()
			suite.Require().NoError(err)
			
			// Check LICENSE file exists and contains correct license
			licensePath := filepath.Join(ctx.projectPath, "LICENSE")
			content, err := os.ReadFile(licensePath)
			suite.Require().NoError(err)
			
			licenseContent := string(content)
			
			// Verify license type
			switch license {
			case "MIT":
				suite.Contains(licenseContent, "MIT License")
			case "Apache-2.0":
				suite.Contains(licenseContent, "Apache License, Version 2.0")
			case "GPL-3.0":
				suite.Contains(licenseContent, "GNU GENERAL PUBLIC LICENSE")
			case "BSD-3-Clause":
				suite.Contains(licenseContent, "BSD 3-Clause License")
			}
		})
	}
}

// TestCIConfiguration tests CI/CD setup
func (suite *LibraryAcceptanceTestSuite) TestCIConfiguration() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-ci --type=library-standard --module=github.com/test/ci --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check CI workflow files
	ciFiles := []string{
		".github/workflows/ci.yml",
		".github/workflows/release.yml",
	}
	
	err = ctx.checkRequiredFiles(ciFiles)
	suite.Require().NoError(err)
	
	// Verify CI workflow content
	ciPath := filepath.Join(ctx.projectPath, ".github/workflows/ci.yml")
	content, err := os.ReadFile(ciPath)
	suite.Require().NoError(err)
	
	ciContent := string(content)
	
	// Check for essential CI steps
	suite.Contains(ciContent, "go test")
	suite.Contains(ciContent, "go fmt")
	suite.Contains(ciContent, "golangci-lint")
	suite.Contains(ciContent, "coverage")
	
	// Check for multiple Go versions
	suite.Contains(ciContent, "matrix")
	suite.Contains(ciContent, "go-version")
}

// TestExamplesCompilation tests that all examples compile and run
func (suite *LibraryAcceptanceTestSuite) TestExamplesCompilation() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-examples --type=library-standard --module=github.com/test/examples --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Test basic example
	basicDir := filepath.Join(ctx.projectPath, "examples", "basic")
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = basicDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		suite.T().Logf("Basic example output: %s", string(output))
	}
	suite.Require().NoError(err, "Basic example should run successfully")
	
	// Test advanced example
	advancedDir := filepath.Join(ctx.projectPath, "examples", "advanced")
	cmd = exec.Command("go", "run", "main.go")
	cmd.Dir = advancedDir
	
	output, err = cmd.CombinedOutput()
	if err != nil {
		suite.T().Logf("Advanced example output: %s", string(output))
	}
	suite.Require().NoError(err, "Advanced example should run successfully")
}

// TestPerformanceGeneration tests generation performance
func (suite *LibraryAcceptanceTestSuite) TestPerformanceGeneration() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-perf --type=library-standard --module=github.com/test/perf --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check performance metrics
	suite.Assert().Less(ctx.generationTime.Seconds(), float64(5), "Generation should complete within 5 seconds")
	suite.Assert().Less(ctx.buildTime.Seconds(), float64(10), "Build should complete within 10 seconds")
	
	suite.T().Logf("Performance metrics - Generation: %v, Build: %v", ctx.generationTime, ctx.buildTime)
}

// TestMinimalDependencies tests that libraries have minimal dependencies
func (suite *LibraryAcceptanceTestSuite) TestMinimalDependencies() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-deps --type=library-standard --module=github.com/test/deps --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check go.mod for minimal dependencies
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	suite.Require().NoError(err)
	
	goModContent := string(content)
	
	// Count direct dependencies (excluding test dependencies)
	requireBlock := strings.Split(goModContent, "require (")[1]
	requireBlock = strings.Split(requireBlock, ")")[0]
	dependencies := strings.Split(requireBlock, "\n")
	
	directDeps := 0
	for _, dep := range dependencies {
		dep = strings.TrimSpace(dep)
		if dep != "" && !strings.Contains(dep, "// indirect") {
			directDeps++
		}
	}
	
	// Should have minimal direct dependencies (typically just testify for testing)
	suite.Assert().LessOrEqual(directDeps, 2, "Library should have minimal direct dependencies")
}

// Run the library acceptance test suite
func TestLibraryAcceptanceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping library acceptance tests in short mode")
	}
	
	suite.Run(t, new(LibraryAcceptanceTestSuite))
}

// Benchmark tests for library generation performance
func BenchmarkLibraryGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := InitializeLibraryContext()
		
		err := ctx.iRunTheCommand(fmt.Sprintf("go-starter new bench-lib-%d --type=library-standard --module=github.com/bench/lib-%d --no-git", i, i))
		if err != nil {
			b.Fatalf("Library generation failed: %v", err)
		}
		
		ctx.cleanup()
	}
}

// Benchmark test for compilation speed
func BenchmarkLibraryCompilation(b *testing.B) {
	// Generate a library once
	ctx := InitializeLibraryContext()
	defer ctx.cleanup()
	
	err := ctx.iRunTheCommand("go-starter new bench-compile --type=library-standard --module=github.com/bench/compile --no-git")
	if err != nil {
		b.Fatalf("Library generation failed: %v", err)
	}
	
	// Benchmark compilation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("go", "build", "-v", "./...")
		cmd.Dir = ctx.projectPath
		err := cmd.Run()
		if err != nil {
			b.Fatalf("Compilation failed: %v", err)
		}
	}
}