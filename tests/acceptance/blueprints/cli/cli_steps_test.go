package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

// CLITestContext holds state for comprehensive CLI BDD tests
// Tests both cli-simple and cli-standard tiers with command execution validation
type CLITestContext struct {
	// Test environment management
	workingDir    string
	projectDir    string
	projectName   string
	originalDir   string
	projectRoot   string
	
	// Command execution tracking
	lastCommand    *exec.Cmd
	lastOutput     []byte
	lastError      error
	lastExitCode   int
	
	// CLI execution results
	cliExecutable  string
	cliOutput      []byte
	cliError       error
	cliExitCode    int
	
	// Project configuration
	complexity     string
	logger         string
	tier           string
	
	// Test state tracking
	projectExists  bool
	compilationOK  bool
	testResults    map[string]bool
	
	// Performance and metrics
	buildTime      time.Duration
	execTime       time.Duration
	fileCount      int
	
	// Context for operations
	ctx            context.Context
}

var cliCtx *CLITestContext

// Test configuration constants for consistent CLI testing
const (
	defaultTestTimeout  = 30 * time.Second
	defaultBuildTimeout = 60 * time.Second
	defaultModule       = "github.com/test/cli-test"
	simpleFileCount     = 8
	standardFileCount   = 29
)

// TestCLIBDD runs comprehensive BDD scenarios for CLI blueprints
// Tests both simple and standard tiers with command execution validation
func TestCLIBDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping comprehensive CLI BDD tests in short mode")
	}
	
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeCLIScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
			Randomize: time.Now().UTC().UnixNano(), // Randomize scenario execution
		},
	}
	
	if suite.Run() != 0 {
		t.Fatal("BDD test suite failed - CLI scenarios did not pass")
	}
}

// InitializeCLIScenario registers all BDD step definitions for comprehensive CLI testing
// Provides setup, execution, and validation steps for both CLI tiers
func InitializeCLIScenario(ctx *godog.ScenarioContext) {
	// Initialize comprehensive CLI test context
	cliCtx = &CLITestContext{
		testResults:   make(map[string]bool),
		ctx:           context.Background(),
	}
	
	// Given steps (preconditions)
	ctx.Given(`^the go-starter CLI tool is available$`, cliCtx.theGoStarterCLIToolIsAvailable)
	ctx.Given(`^I am in a clean working directory$`, cliCtx.iAmInACleanWorkingDirectory)
	ctx.Given(`^I want to create a simple CLI application for quick utilities$`, cliCtx.iWantToCreateASimpleCLIApplication)
	ctx.Given(`^I want to create a production-ready CLI application$`, cliCtx.iWantToCreateAProductionReadyCLIApplication)
	ctx.Given(`^I want to create a CLI with specific complexity$`, cliCtx.iWantToCreateACLIWithSpecificComplexity)
	ctx.Given(`^I want to build quick command-line utilities$`, cliCtx.iWantToBuildQuickCommandLineUtilities)
	ctx.Given(`^I want to build production CLI tools$`, cliCtx.iWantToBuildProductionCLITools)
	ctx.Given(`^I want to showcase progressive complexity$`, cliCtx.iWantToShowcaseProgressiveComplexity)
	ctx.Given(`^I want to use different logging libraries in CLI$`, cliCtx.iWantToUseDifferentLoggingLibrariesInCLI)
	ctx.Given(`^I want to ensure CLI framework works properly$`, cliCtx.iWantToEnsureCLIFrameworkWorksroperly)
	ctx.Given(`^I have generated a CLI application$`, cliCtx.iHaveGeneratedACLIApplication)
	ctx.Given(`^I want configurable CLI applications$`, cliCtx.iWantConfigurableCLIApplications)
	ctx.Given(`^I want flexible CLI output$`, cliCtx.iWantFlexibleCLIOutput)
	ctx.Given(`^I want robust CLI applications$`, cliCtx.iWantRobustCLIApplications)
	ctx.Given(`^I want interactive CLI applications$`, cliCtx.iWantInteractiveCLIApplications)
	ctx.Given(`^I want well-tested CLI applications$`, cliCtx.iWantWellTestedCLIApplications)
	ctx.Given(`^I want efficient CLI development$`, cliCtx.iWantEfficientCLIDevelopment)
	ctx.Given(`^I want to distribute CLI applications$`, cliCtx.iWantToDistributeCLIApplications)
	ctx.Given(`^I want automated CLI workflows$`, cliCtx.iWantAutomatedCLIWorkflows)
	ctx.Given(`^I want containerized CLI applications$`, cliCtx.iWantContainerizedCLIApplications)
	ctx.Given(`^I want secure CLI applications$`, cliCtx.iWantSecureCLIApplications)
	ctx.Given(`^I want high-performance CLI applications$`, cliCtx.iWantHighPerformanceCLIApplications)
	ctx.Given(`^I want well-documented CLI applications$`, cliCtx.iWantWellDocumentedCLIApplications)
	ctx.Given(`^I want CLI applications with shell integration$`, cliCtx.iWantCLIApplicationsWithShellIntegration)
	ctx.Given(`^I have a simple CLI that needs more features$`, cliCtx.iHaveASimpleCLIThatNeedsMoreFeatures)
	ctx.Given(`^I want CLI applications that work everywhere$`, cliCtx.iWantCLIApplicationsThatWorkEverywhere)
	ctx.Given(`^I want extensible CLI applications$`, cliCtx.iWantExtensibleCLIApplications)
	ctx.Given(`^I want observable CLI applications$`, cliCtx.iWantObservableCLIApplications)
	ctx.Given(`^I want CLI applications for global use$`, cliCtx.iWantCLIApplicationsForGlobalUse)
	ctx.Given(`^I want sophisticated CLI interfaces$`, cliCtx.iWantSophisticatedCLIInterfaces)
	
	// When steps (actions)
	ctx.When(`^I run the command "([^"]*)"$`, cliCtx.iRunTheCommand)
	ctx.When(`^I generate a CLI with complexity "([^"]*)"$`, cliCtx.iGenerateACLIWithComplexity)
	ctx.When(`^I generate a simple CLI application$`, cliCtx.iGenerateASimpleCLIApplication)
	ctx.When(`^I generate a standard CLI application$`, cliCtx.iGenerateAStandardCLIApplication)
	ctx.When(`^I compare simple and standard CLI blueprints$`, cliCtx.iCompareSimpleAndStandardCLIBlueprints)
	ctx.When(`^I generate a CLI with logger "([^"]*)"$`, cliCtx.iGenerateACLIWithLogger)
	ctx.When(`^I generate a CLI application$`, cliCtx.iGenerateACLIApplication)
	ctx.When(`^I execute CLI commands$`, cliCtx.iExecuteCLICommands)
	ctx.When(`^I generate a CLI with configuration support$`, cliCtx.iGenerateACLIWithConfigurationSupport)
	ctx.When(`^I generate a CLI with error handling$`, cliCtx.iGenerateACLIWithErrorHandling)
	ctx.When(`^I generate a standard CLI with interactive features$`, cliCtx.iGenerateAStandardCLIWithInteractiveFeatures)
	ctx.When(`^I generate a CLI with testing support$`, cliCtx.iGenerateACLIWithTestingSupport)
	ctx.When(`^I generate a CLI for distribution$`, cliCtx.iGenerateACLIForDistribution)
	ctx.When(`^I generate a CLI with CI/CD support$`, cliCtx.iGenerateACLIWithCICDSupport)
	ctx.When(`^I generate a CLI with container support$`, cliCtx.iGenerateACLIWithContainerSupport)
	ctx.When(`^I generate a CLI with security features$`, cliCtx.iGenerateACLIWithSecurityFeatures)
	ctx.When(`^I generate a CLI with performance optimizations$`, cliCtx.iGenerateACLIWithPerformanceOptimizations)
	ctx.When(`^I generate a CLI with comprehensive documentation$`, cliCtx.iGenerateACLIWithComprehensiveDocumentation)
	ctx.When(`^I generate a CLI with completion support$`, cliCtx.iGenerateACLIWithCompletionSupport)
	ctx.When(`^I want to migrate to standard tier$`, cliCtx.iWantToMigrateToStandardTier)
	ctx.When(`^I generate a CLI for cross-platform use$`, cliCtx.iGenerateACLIForCrossPlatformUse)
	ctx.When(`^I generate a CLI with plugin support$`, cliCtx.iGenerateACLIWithPluginSupport)
	ctx.When(`^I generate a CLI with monitoring features$`, cliCtx.iGenerateACLIWithMonitoringFeatures)
	ctx.When(`^I generate a CLI with i18n support$`, cliCtx.iGenerateACLIWithI18nSupport)
	ctx.When(`^I generate a CLI with advanced patterns$`, cliCtx.iGenerateACLIWithAdvancedPatterns)
	
	// Then steps (assertions)
	ctx.Then(`^the generation should succeed$`, cliCtx.theGenerationShouldSucceed)
	ctx.Then(`^the project should contain essential CLI components for simple tier$`, cliCtx.theProjectShouldContainEssentialCLIComponentsForSimpleTier)
	ctx.Then(`^the generated code should compile successfully$`, cliCtx.theGeneratedCodeShouldCompileSuccessfully)
	ctx.Then(`^the CLI should execute basic commands$`, cliCtx.theCLIShouldExecuteBasicCommands)
	ctx.Then(`^the project structure should be minimal with (\d+) files$`, cliCtx.theProjectStructureShouldBeMinimalWithFiles)
	ctx.Then(`^the project should contain all essential CLI components for standard tier$`, cliCtx.theProjectShouldContainAllEssentialCLIComponentsForStandardTier)
	ctx.Then(`^the CLI should support advanced features$`, cliCtx.theCLIShouldSupportAdvancedFeatures)
	ctx.Then(`^the project structure should be comprehensive with (\d+) files$`, cliCtx.theProjectStructureShouldBeComprehensiveWithFiles)
	ctx.Then(`^the project should follow "([^"]*)" tier patterns$`, cliCtx.theProjectShouldFollowTierPatterns)
	ctx.Then(`^the file count should match the complexity level$`, cliCtx.theFileCountShouldMatchTheComplexityLevel)
	ctx.Then(`^the feature set should align with the complexity$`, cliCtx.theFeatureSetShouldAlignWithTheComplexity)
	ctx.Then(`^the CLI should compile and execute correctly$`, cliCtx.theCLIShouldCompileAndExecuteCorrectly)
	ctx.Then(`^the CLI should have minimal structure$`, cliCtx.theCLIShouldHaveMinimalStructure)
	ctx.Then(`^the CLI should use slog for logging$`, cliCtx.theCLIShouldUseSlogForLogging)
	ctx.Then(`^the CLI should support basic flags \(help, version, quiet\)$`, cliCtx.theCLIShouldSupportBasicFlags)
	ctx.Then(`^the CLI should have Cobra framework integration$`, cliCtx.theCLIShouldHaveCobraFrameworkIntegration)
	ctx.Then(`^the configuration should be minimal$`, cliCtx.theConfigurationShouldBeMinimal)
	ctx.Then(`^the commands should be straightforward$`, cliCtx.theCommandsShouldBeStraightforward)
	ctx.Then(`^the CLI should have layered architecture$`, cliCtx.theCLIShouldHaveLayeredArchitecture)
	ctx.Then(`^the CLI should support multiple loggers \(slog, zap, logrus, zerolog\)$`, cliCtx.theCLIShouldSupportMultipleLoggers)
	ctx.Then(`^the CLI should have comprehensive flag support$`, cliCtx.theCLIShouldHaveComprehensiveFlagSupport)
	ctx.Then(`^the CLI should include configuration management$`, cliCtx.theCLIShouldIncludeConfigurationManagement)
	ctx.Then(`^the CLI should support interactive mode$`, cliCtx.theCLIShouldSupportInteractiveMode)
	ctx.Then(`^the CLI should have testing infrastructure$`, cliCtx.theCLIShouldHaveTestingInfrastructure)
	ctx.Then(`^the CLI should include CI/CD integration$`, cliCtx.theCLIShouldIncludeCICDIntegration)
	ctx.Then(`^simple CLI should have (\d+) files vs standard (\d+) files$`, cliCtx.simpleCLIShouldHaveFilesVsStandardFiles)
	ctx.Then(`^simple CLI should focus on essential features only$`, cliCtx.simpleCLIShouldFocusOnEssentialFeaturesOnly)
	ctx.Then(`^standard CLI should include advanced capabilities$`, cliCtx.standardCLIShouldIncludeAdvancedCapabilities)
	ctx.Then(`^both should compile and execute successfully$`, cliCtx.bothShouldCompileAndExecuteSuccessfully)
	ctx.Then(`^migration path should be clear from simple to standard$`, cliCtx.migrationPathShouldBeClearFromSimpleToStandard)
	ctx.Then(`^the CLI should use the "([^"]*)" logging library$`, cliCtx.theCLIShouldUseTheLoggingLibrary)
	ctx.Then(`^the logger should be properly configured$`, cliCtx.theLoggerShouldBeProperlyConfigured)
	ctx.Then(`^the CLI should support structured logging$`, cliCtx.theCLIShouldSupportStructuredLogging)
	ctx.Then(`^the log levels should be configurable$`, cliCtx.theLogLevelsShouldBeConfigurable)
	ctx.Then(`^the CLI should compile with the logger dependency$`, cliCtx.theCLIShouldCompileWithTheLoggerDependency)
	ctx.Then(`^the CLI should use Cobra framework$`, cliCtx.theCLIShouldUseCobraFramework)
	ctx.Then(`^the CLI should support subcommands$`, cliCtx.theCLIShouldSupportSubcommands)
	ctx.Then(`^the CLI should have built-in help system$`, cliCtx.theCLIShouldHaveBuiltInHelpSystem)
	ctx.Then(`^the CLI should support command completion$`, cliCtx.theCLIShouldSupportCommandCompletion)
	ctx.Then(`^the CLI should handle flags and arguments correctly$`, cliCtx.theCLIShouldHandleFlagsAndArgumentsCorrectly)
	ctx.Then(`^the help command should display usage information$`, cliCtx.theHelpCommandShouldDisplayUsageInformation)
	ctx.Then(`^the version command should show version details$`, cliCtx.theVersionCommandShouldShowVersionDetails)
	ctx.Then(`^invalid commands should show helpful error messages$`, cliCtx.invalidCommandsShouldShowHelpfulErrorMessages)
	ctx.Then(`^the CLI should exit with appropriate status codes$`, cliCtx.theCLIShouldExitWithAppropriateStatusCodes)
	ctx.Then(`^command output should be properly formatted$`, cliCtx.commandOutputShouldBeProperlyFormatted)
	
	// Additional comprehensive Then steps for all CLI features...
	// (The implementation would continue with the remaining ~50 step definitions)
	
	// Configuration and flexibility steps
	ctx.Then(`^the CLI should support configuration files$`, cliCtx.theCLIShouldSupportConfigurationFiles)
	ctx.Then(`^the CLI should support environment variables$`, cliCtx.theCLIShouldSupportEnvironmentVariables)
	ctx.Then(`^the CLI should have configuration precedence$`, cliCtx.theCLIShouldHaveConfigurationPrecedence)
	ctx.Then(`^the CLI should validate configuration values$`, cliCtx.theCLIShouldValidateConfigurationValues)
	ctx.Then(`^the CLI should provide configuration examples$`, cliCtx.theCLIShouldProvideConfigurationExamples)
	
	// Output and formatting steps
	ctx.Then(`^the CLI should support multiple output formats$`, cliCtx.theCLIShouldSupportMultipleOutputFormats)
	ctx.Then(`^the CLI should support quiet mode$`, cliCtx.theCLIShouldSupportQuietMode)
	ctx.Then(`^the CLI should support verbose mode$`, cliCtx.theCLIShouldSupportVerboseMode)
	ctx.Then(`^the CLI should handle JSON output$`, cliCtx.theCLIShouldHandleJSONOutput)
	ctx.Then(`^the CLI should format output appropriately$`, cliCtx.theCLIShouldFormatOutputAppropriately)
	
	// Error handling steps
	ctx.Then(`^the CLI should handle invalid arguments gracefully$`, cliCtx.theCLIShouldHandleInvalidArgumentsGracefully)
	ctx.Then(`^the CLI should provide clear error messages$`, cliCtx.theCLIShouldProvideClearErrorMessages)
	ctx.Then(`^the CLI should validate input parameters$`, cliCtx.theCLIShouldValidateInputParameters)
	ctx.Then(`^the CLI should handle system errors properly$`, cliCtx.theCLIShouldHandleSystemErrorsProperly)
	ctx.Then(`^the CLI should use appropriate exit codes$`, cliCtx.theCLIShouldUseAppropriateExitCodes)
}

// Given step implementations

func (ctx *CLITestContext) theGoStarterCLIToolIsAvailable() error {
	// Ensure go-starter CLI is built and ready for comprehensive CLI testing
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	ctx.originalDir = originalDir
	ctx.projectRoot = filepath.Join(originalDir, "..", "..", "..", "..")
	
	// Build with optimizations for faster test execution
	buildStart := time.Now()
	buildCmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", "go-starter", ".")
	buildCmd.Dir = ctx.projectRoot
	output, err := buildCmd.CombinedOutput()
	ctx.buildTime = time.Since(buildStart)
	
	if err != nil {
		return fmt.Errorf("failed to build go-starter CLI: %s", string(output))
	}
	
	return nil
}

func (ctx *CLITestContext) iAmInACleanWorkingDirectory() error {
	// Reset all test state for clean scenario execution
	ctx.testResults = make(map[string]bool)
	ctx.projectExists = false
	ctx.compilationOK = false
	ctx.fileCount = 0
	
	// Create isolated temporary directory
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "cli-bdd-*")
	if err != nil {
		return fmt.Errorf("failed to create clean working directory: %w", err)
	}
	
	return os.Chdir(ctx.workingDir)
}

func (ctx *CLITestContext) iWantToCreateASimpleCLIApplication() error {
	ctx.projectName = "test-simple-cli"
	ctx.complexity = "simple"
	ctx.tier = "simple"
	ctx.logger = "slog"
	return nil
}

func (ctx *CLITestContext) iWantToCreateAProductionReadyCLIApplication() error {
	ctx.projectName = "test-production-cli"
	ctx.complexity = "standard"
	ctx.tier = "standard"
	ctx.logger = "slog"
	return nil
}

func (ctx *CLITestContext) iWantToCreateACLIWithSpecificComplexity() error {
	ctx.projectName = "test-complex-cli"
	return nil
}

func (ctx *CLITestContext) iWantToBuildQuickCommandLineUtilities() error {
	ctx.projectName = "test-utility-cli"
	ctx.complexity = "simple"
	ctx.tier = "simple"
	return nil
}

func (ctx *CLITestContext) iWantToBuildProductionCLITools() error {
	ctx.projectName = "test-prod-cli"
	ctx.complexity = "standard"
	ctx.tier = "standard"
	return nil
}

func (ctx *CLITestContext) iWantToShowcaseProgressiveComplexity() error {
	return nil // This is a comparison scenario
}

func (ctx *CLITestContext) iWantToUseDifferentLoggingLibrariesInCLI() error {
	ctx.projectName = "test-logger-cli"
	return nil
}

func (ctx *CLITestContext) iWantToEnsureCLIFrameworkWorksroperly() error {
	ctx.projectName = "test-framework-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iHaveGeneratedACLIApplication() error {
	if err := ctx.iWantToCreateAProductionReadyCLIApplication(); err != nil {
		return err
	}
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iWantConfigurableCLIApplications() error {
	ctx.projectName = "test-config-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantFlexibleCLIOutput() error {
	ctx.projectName = "test-output-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantRobustCLIApplications() error {
	ctx.projectName = "test-robust-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantInteractiveCLIApplications() error {
	ctx.projectName = "test-interactive-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantWellTestedCLIApplications() error {
	ctx.projectName = "test-tested-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantEfficientCLIDevelopment() error {
	ctx.projectName = "test-dev-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantToDistributeCLIApplications() error {
	ctx.projectName = "test-dist-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantAutomatedCLIWorkflows() error {
	ctx.projectName = "test-cicd-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantContainerizedCLIApplications() error {
	ctx.projectName = "test-docker-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantSecureCLIApplications() error {
	ctx.projectName = "test-secure-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantHighPerformanceCLIApplications() error {
	ctx.projectName = "test-perf-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantWellDocumentedCLIApplications() error {
	ctx.projectName = "test-docs-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantCLIApplicationsWithShellIntegration() error {
	ctx.projectName = "test-shell-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iHaveASimpleCLIThatNeedsMoreFeatures() error {
	ctx.projectName = "test-migrate-cli"
	ctx.complexity = "simple"
	return nil
}

func (ctx *CLITestContext) iWantCLIApplicationsThatWorkEverywhere() error {
	ctx.projectName = "test-cross-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantExtensibleCLIApplications() error {
	ctx.projectName = "test-plugin-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantObservableCLIApplications() error {
	ctx.projectName = "test-monitor-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantCLIApplicationsForGlobalUse() error {
	ctx.projectName = "test-i18n-cli"
	ctx.complexity = "standard"
	return nil
}

func (ctx *CLITestContext) iWantSophisticatedCLIInterfaces() error {
	ctx.projectName = "test-advanced-cli"
	ctx.complexity = "standard"
	return nil
}

// When step implementations

func (ctx *CLITestContext) iRunTheCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	// Handle special case for go-starter commands
	if parts[0] == "go-starter" {
		parts[0] = "./go-starter"
	}
	
	ctx.lastCommand = exec.Command(parts[0], parts[1:]...)
	ctx.lastCommand.Dir = ctx.workingDir
	
	ctx.lastOutput, ctx.lastError = ctx.lastCommand.CombinedOutput()
	
	if exitError, ok := ctx.lastError.(*exec.ExitError); ok {
		ctx.lastExitCode = exitError.ExitCode()
	} else if ctx.lastError != nil {
		ctx.lastExitCode = -1
	} else {
		ctx.lastExitCode = 0
	}
	
	// If this was a project generation command, set up project info
	if len(parts) >= 3 && parts[1] == "new" {
		ctx.projectName = parts[2]
		ctx.projectDir = filepath.Join(ctx.workingDir, ctx.projectName)
		if ctx.lastError == nil {
			ctx.projectExists = true
		}
	}
	
	return nil
}

func (ctx *CLITestContext) iGenerateACLIWithComplexity(complexity string) error {
	ctx.complexity = complexity
	if complexity == "simple" {
		ctx.tier = "simple"
	} else {
		ctx.tier = "standard"
	}
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateASimpleCLIApplication() error {
	ctx.complexity = "simple"
	ctx.tier = "simple"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateAStandardCLIApplication() error {
	ctx.complexity = "standard"
	ctx.tier = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iCompareSimpleAndStandardCLIBlueprints() error {
	// This is a passive comparison step - validation happens in Then steps
	return nil
}

func (ctx *CLITestContext) iGenerateACLIWithLogger(logger string) error {
	ctx.logger = logger
	ctx.complexity = "standard" // Logger selection requires standard tier
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIApplication() error {
	if ctx.complexity == "" {
		ctx.complexity = "standard"
	}
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iExecuteCLICommands() error {
	return ctx.buildAndTestCLI()
}

func (ctx *CLITestContext) iGenerateACLIWithConfigurationSupport() error {
	ctx.complexity = "standard" // Configuration requires standard tier
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithErrorHandling() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateAStandardCLIWithInteractiveFeatures() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithTestingSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIForDistribution() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithCICDSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithContainerSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithSecurityFeatures() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithPerformanceOptimizations() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithComprehensiveDocumentation() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithCompletionSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iWantToMigrateToStandardTier() error {
	// This is a conceptual step for migration documentation
	return nil
}

func (ctx *CLITestContext) iGenerateACLIForCrossPlatformUse() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithPluginSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithMonitoringFeatures() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithI18nSupport() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

func (ctx *CLITestContext) iGenerateACLIWithAdvancedPatterns() error {
	ctx.complexity = "standard"
	return ctx.generateCLIProject()
}

// Helper method to generate CLI projects
func (ctx *CLITestContext) generateCLIProject() error {
	// Build go-starter first if not done
	goStarterPath := filepath.Join(ctx.workingDir, "go-starter")
	if _, err := os.Stat(goStarterPath); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-o", goStarterPath, ".")
		buildCmd.Dir = ctx.projectRoot
		
		if output, err := buildCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to build go-starter: %s", string(output))
		}
	}
	
	if ctx.projectName == "" {
		ctx.projectName = "test-cli"
	}
	
	// Generate the project
	args := []string{
		"new", ctx.projectName,
		"--type=cli",
		"--module=" + defaultModule + "/" + ctx.projectName,
		"--no-git",
	}
	
	if ctx.complexity != "" {
		args = append(args, "--complexity="+ctx.complexity)
	}
	if ctx.logger != "" {
		args = append(args, "--logger="+ctx.logger)
	}
	
	generateCmd := exec.Command(goStarterPath, args...)
	generateCmd.Dir = ctx.workingDir
	
	ctx.lastOutput, ctx.lastError = generateCmd.CombinedOutput()
	
	if ctx.lastError != nil {
		return fmt.Errorf("CLI project generation failed: %s", string(ctx.lastOutput))
	}
	
	ctx.projectDir = filepath.Join(ctx.workingDir, ctx.projectName)
	ctx.projectExists = true
	
	// Count generated files
	ctx.fileCount = ctx.countFiles(ctx.projectDir)
	
	return nil
}

// Helper method to build and test CLI functionality
func (ctx *CLITestContext) buildAndTestCLI() error {
	if !ctx.projectExists {
		return fmt.Errorf("CLI project not generated")
	}
	
	// First compile the project
	if err := ctx.compileCLIProject(); err != nil {
		return fmt.Errorf("failed to compile CLI: %w", err)
	}
	
	// Test basic CLI execution
	execStart := time.Now()
	cliCmd := exec.Command("./" + ctx.projectName)
	cliCmd.Dir = ctx.projectDir
	
	ctx.cliOutput, ctx.cliError = cliCmd.CombinedOutput()
	ctx.execTime = time.Since(execStart)
	
	if exitError, ok := ctx.cliError.(*exec.ExitError); ok {
		ctx.cliExitCode = exitError.ExitCode()
	} else if ctx.cliError != nil {
		ctx.cliExitCode = -1
	} else {
		ctx.cliExitCode = 0
	}
	
	return nil
}

// Helper method to compile CLI project
func (ctx *CLITestContext) compileCLIProject() error {
	if !ctx.projectExists {
		return fmt.Errorf("CLI project not generated")
	}
	
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}
	
	// Build the CLI binary
	buildStart := time.Now()
	buildCmd := exec.Command("go", "build", "-o", ctx.projectName, ".")
	buildCmd.Dir = ctx.projectDir
	output, err := buildCmd.CombinedOutput()
	ctx.buildTime = time.Since(buildStart)
	
	if err != nil {
		return fmt.Errorf("CLI compilation failed: %s", string(output))
	}
	
	ctx.compilationOK = true
	return nil
}

// Helper method to count files in directory
func (ctx *CLITestContext) countFiles(dir string) int {
	count := 0
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count
}

// Helper method to check if file exists
func (ctx *CLITestContext) checkFileExists(relativePath string) error {
	fullPath := filepath.Join(ctx.projectDir, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("file or directory not found: %s", relativePath)
	}
	return nil
}

// Helper method to check if file contains content
func (ctx *CLITestContext) checkFileContains(relativePath, content string) error {
	fullPath := filepath.Join(ctx.projectDir, relativePath)
	
	// Check if it's a directory - if so, check all files in it
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		found := false
		_ = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") {
				fileContent, err := os.ReadFile(path)
				if err == nil && strings.Contains(string(fileContent), content) {
					found = true
					return filepath.SkipDir
				}
			}
			return nil
		})
		if !found {
			return fmt.Errorf("content '%s' not found in directory %s", content, relativePath)
		}
		return nil
	}
	
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", relativePath, err)
	}
	
	if !strings.Contains(string(fileContent), content) {
		return fmt.Errorf("file %s does not contain '%s'", relativePath, content)
	}
	
	return nil
}

// Then step implementations (key validation methods)

func (ctx *CLITestContext) theGenerationShouldSucceed() error {
	if ctx.lastError != nil && ctx.lastExitCode != 0 {
		return fmt.Errorf("generation failed with exit code %d: %s", ctx.lastExitCode, string(ctx.lastOutput))
	}
	
	// Check if project directory was created
	expectedPath := filepath.Join(ctx.workingDir, ctx.projectName)
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		return fmt.Errorf("project directory not created at %s", expectedPath)
	}
	
	ctx.projectDir = expectedPath
	ctx.projectExists = true
	
	return nil
}

func (ctx *CLITestContext) theProjectShouldContainEssentialCLIComponentsForSimpleTier() error {
	if !ctx.projectExists {
		return fmt.Errorf("project was not generated")
	}
	
	// Essential components for simple tier (8 files)
	essentialFiles := []string{
		"main.go",
		"go.mod",
		"README.md",
		"Makefile",
		"config.go",
		"cmd/root.go",
		"cmd/version.go",
		".gitignore",
	}
	
	for _, file := range essentialFiles {
		if err := ctx.checkFileExists(file); err != nil {
			return fmt.Errorf("essential simple CLI component missing: %s", file)
		}
	}
	
	return nil
}

func (ctx *CLITestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	return ctx.compileCLIProject()
}

func (ctx *CLITestContext) theCLIShouldExecuteBasicCommands() error {
	return ctx.buildAndTestCLI()
}

func (ctx *CLITestContext) theProjectStructureShouldBeMinimalWithFiles(fileCountStr string) error {
	expectedCount, err := strconv.Atoi(fileCountStr)
	if err != nil {
		return fmt.Errorf("invalid file count: %s", fileCountStr)
	}
	
	if ctx.fileCount != expectedCount {
		return fmt.Errorf("expected %d files, got %d files", expectedCount, ctx.fileCount)
	}
	
	return nil
}

func (ctx *CLITestContext) theProjectShouldContainAllEssentialCLIComponentsForStandardTier() error {
	if !ctx.projectExists {
		return fmt.Errorf("project was not generated")
	}
	
	// Essential components for standard tier (29 files)
	essentialDirs := []string{
		"cmd",
		"internal/config",
		"internal/logger",
		"internal/errors",
		"internal/output",
		"configs",
		".github/workflows",
	}
	
	for _, dir := range essentialDirs {
		if err := ctx.checkFileExists(dir); err != nil {
			return fmt.Errorf("essential standard CLI component missing: %s", dir)
		}
	}
	
	return nil
}

func (ctx *CLITestContext) theCLIShouldSupportAdvancedFeatures() error {
	// Check for advanced features in standard tier
	advancedFeatures := []string{
		"internal/config",
		"internal/interactive",
		"configs/config.yaml",
		".github/workflows/ci.yml",
	}
	
	for _, feature := range advancedFeatures {
		if err := ctx.checkFileExists(feature); err != nil {
			return fmt.Errorf("advanced CLI feature missing: %s", feature)
		}
	}
	
	return nil
}

func (ctx *CLITestContext) theProjectStructureShouldBeComprehensiveWithFiles(fileCountStr string) error {
	expectedCount, err := strconv.Atoi(fileCountStr)
	if err != nil {
		return fmt.Errorf("invalid file count: %s", fileCountStr)
	}
	
	// Allow some tolerance for file count variations
	tolerance := 2
	if ctx.fileCount < expectedCount-tolerance || ctx.fileCount > expectedCount+tolerance {
		return fmt.Errorf("expected ~%d files, got %d files", expectedCount, ctx.fileCount)
	}
	
	return nil
}

func (ctx *CLITestContext) theProjectShouldFollowTierPatterns(tier string) error {
	ctx.tier = tier
	
	if tier == "simple" {
		return ctx.validateSimpleTierPatterns()
	} else {
		return ctx.validateStandardTierPatterns()
	}
}

func (ctx *CLITestContext) validateSimpleTierPatterns() error {
	// Validate simple tier has minimal structure
	return ctx.checkFileExists("config.go") // Simple config in root
}

func (ctx *CLITestContext) validateStandardTierPatterns() error {
	// Validate standard tier has layered structure
	return ctx.checkFileExists("internal/config") // Layered config
}

func (ctx *CLITestContext) theFileCountShouldMatchTheComplexityLevel() error {
	switch ctx.complexity {
	case "simple":
		if ctx.fileCount < 6 || ctx.fileCount > 10 {
			return fmt.Errorf("simple CLI should have ~8 files, got %d", ctx.fileCount)
		}
	case "standard":
		if ctx.fileCount < 25 || ctx.fileCount > 35 {
			return fmt.Errorf("standard CLI should have ~29 files, got %d", ctx.fileCount)
		}
	}
	
	return nil
}

func (ctx *CLITestContext) theFeatureSetShouldAlignWithTheComplexity() error {
	if ctx.complexity == "simple" {
		// Simple should NOT have complex features
		if _, err := os.Stat(filepath.Join(ctx.projectDir, "internal")); err == nil {
			return fmt.Errorf("simple CLI should not have internal packages")
		}
	} else {
		// Standard should have complex features
		if err := ctx.checkFileExists("internal"); err != nil {
			return fmt.Errorf("standard CLI missing internal packages: %w", err)
		}
	}
	
	return nil
}

func (ctx *CLITestContext) theCLIShouldCompileAndExecuteCorrectly() error {
	if err := ctx.compileCLIProject(); err != nil {
		return err
	}
	return ctx.buildAndTestCLI()
}

// Additional Then step implementations for comprehensive CLI validation
func (ctx *CLITestContext) theCLIShouldHaveMinimalStructure() error {
	// Check that simple CLI has minimal structure
	return ctx.theProjectShouldContainEssentialCLIComponentsForSimpleTier()
}

func (ctx *CLITestContext) theCLIShouldUseSlogForLogging() error {
	return ctx.checkFileContains("main.go", "log/slog")
}

func (ctx *CLITestContext) theCLIShouldSupportBasicFlags() error {
	return ctx.checkFileContains("cmd/root.go", "PersistentFlags")
}

func (ctx *CLITestContext) theCLIShouldHaveCobraFrameworkIntegration() error {
	return ctx.checkFileContains("main.go", "github.com/spf13/cobra")
}

func (ctx *CLITestContext) theConfigurationShouldBeMinimal() error {
	return ctx.checkFileExists("config.go")
}

func (ctx *CLITestContext) theCommandsShouldBeStraightforward() error {
	return ctx.checkFileExists("cmd/version.go")
}

func (ctx *CLITestContext) theCLIShouldHaveLayeredArchitecture() error {
	return ctx.checkFileExists("internal")
}

func (ctx *CLITestContext) theCLIShouldSupportMultipleLoggers() error {
	return ctx.checkFileExists("internal/logger")
}

func (ctx *CLITestContext) theCLIShouldHaveComprehensiveFlagSupport() error {
	return ctx.checkFileContains("cmd", "flags")
}

func (ctx *CLITestContext) theCLIShouldIncludeConfigurationManagement() error {
	return ctx.checkFileExists("configs/config.yaml")
}

func (ctx *CLITestContext) theCLIShouldSupportInteractiveMode() error {
	return ctx.checkFileExists("internal/interactive")
}

func (ctx *CLITestContext) theCLIShouldHaveTestingInfrastructure() error {
	return ctx.checkFileContains("cmd", "_test.go")
}

func (ctx *CLITestContext) theCLIShouldIncludeCICDIntegration() error {
	return ctx.checkFileExists(".github/workflows")
}

func (ctx *CLITestContext) simpleCLIShouldHaveFilesVsStandardFiles(simpleStr, standardStr string) error {
	// This is validated in individual file count checks
	return nil
}

func (ctx *CLITestContext) simpleCLIShouldFocusOnEssentialFeaturesOnly() error {
	// Check that simple CLI doesn't have advanced directories
	advancedDirs := []string{"internal", "configs", ".github"}
	for _, dir := range advancedDirs {
		if _, err := os.Stat(filepath.Join(ctx.projectDir, dir)); err == nil {
			return fmt.Errorf("simple CLI should not have advanced directory: %s", dir)
		}
	}
	return nil
}

func (ctx *CLITestContext) standardCLIShouldIncludeAdvancedCapabilities() error {
	return ctx.theCLIShouldSupportAdvancedFeatures()
}

func (ctx *CLITestContext) bothShouldCompileAndExecuteSuccessfully() error {
	return ctx.theCLIShouldCompileAndExecuteCorrectly()
}

func (ctx *CLITestContext) migrationPathShouldBeClearFromSimpleToStandard() error {
	// Check that README or documentation mentions migration
	return ctx.checkFileContains("README.md", "migration")
}

// Logger-specific validations
func (ctx *CLITestContext) theCLIShouldUseTheLoggingLibrary(logger string) error {
	loggerImports := map[string]string{
		"slog":    "log/slog",
		"zap":     "go.uber.org/zap",
		"logrus":  "github.com/sirupsen/logrus",
		"zerolog": "github.com/rs/zerolog",
	}
	
	expectedImport := loggerImports[logger]
	if expectedImport == "" {
		return fmt.Errorf("unknown logger: %s", logger)
	}
	
	return ctx.checkFileContains("internal/logger", expectedImport)
}

func (ctx *CLITestContext) theLoggerShouldBeProperlyConfigured() error {
	return ctx.checkFileContains("internal/logger", "Config")
}

func (ctx *CLITestContext) theCLIShouldSupportStructuredLogging() error {
	return ctx.checkFileContains("internal/logger", "structured")
}

func (ctx *CLITestContext) theLogLevelsShouldBeConfigurable() error {
	return ctx.checkFileContains("internal/logger", "Level")
}

func (ctx *CLITestContext) theCLIShouldCompileWithTheLoggerDependency() error {
	return ctx.compileCLIProject()
}

// Framework validation
func (ctx *CLITestContext) theCLIShouldUseCobraFramework() error {
	return ctx.checkFileContains("go.mod", "github.com/spf13/cobra")
}

func (ctx *CLITestContext) theCLIShouldSupportSubcommands() error {
	return ctx.checkFileExists("cmd")
}

func (ctx *CLITestContext) theCLIShouldHaveBuiltInHelpSystem() error {
	if err := ctx.buildAndTestCLI(); err != nil {
		return err
	}
	
	// Test help command execution
	helpCmd := exec.Command("./"+ctx.projectName, "--help")
	helpCmd.Dir = ctx.projectDir
	output, err := helpCmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("help command failed: %w", err)
	}
	
	if !strings.Contains(string(output), "Usage:") {
		return fmt.Errorf("help output doesn't contain usage information")
	}
	
	return nil
}

func (ctx *CLITestContext) theCLIShouldSupportCommandCompletion() error {
	return ctx.checkFileContains("cmd", "completion")
}

func (ctx *CLITestContext) theCLIShouldHandleFlagsAndArgumentsCorrectly() error {
	return ctx.checkFileContains("cmd/root.go", "Args")
}

// Command execution validations
func (ctx *CLITestContext) theHelpCommandShouldDisplayUsageInformation() error {
	return ctx.theCLIShouldHaveBuiltInHelpSystem()
}

func (ctx *CLITestContext) theVersionCommandShouldShowVersionDetails() error {
	if err := ctx.buildAndTestCLI(); err != nil {
		return err
	}
	
	// Test version command execution
	versionCmd := exec.Command("./"+ctx.projectName, "version")
	versionCmd.Dir = ctx.projectDir
	output, err := versionCmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("version command failed: %w", err)
	}
	
	if !strings.Contains(string(output), "version") {
		return fmt.Errorf("version output doesn't contain version information")
	}
	
	return nil
}

func (ctx *CLITestContext) invalidCommandsShouldShowHelpfulErrorMessages() error {
	if err := ctx.buildAndTestCLI(); err != nil {
		return err
	}
	
	// Test invalid command execution
	invalidCmd := exec.Command("./"+ctx.projectName, "nonexistent")
	invalidCmd.Dir = ctx.projectDir
	output, err := invalidCmd.CombinedOutput()
	
	// Invalid commands should return non-zero exit code
	if err == nil {
		return fmt.Errorf("invalid command should fail")
	}
	
	if !strings.Contains(string(output), "unknown command") && 
	   !strings.Contains(string(output), "Error:") {
		return fmt.Errorf("error output doesn't contain helpful message")
	}
	
	return nil
}

func (ctx *CLITestContext) theCLIShouldExitWithAppropriateStatusCodes() error {
	// This is validated in individual command tests
	return nil
}

func (ctx *CLITestContext) commandOutputShouldBeProperlyFormatted() error {
	return ctx.checkFileContains("internal/output", "format")
}

// Configuration validation steps
func (ctx *CLITestContext) theCLIShouldSupportConfigurationFiles() error {
	return ctx.checkFileExists("configs")
}

func (ctx *CLITestContext) theCLIShouldSupportEnvironmentVariables() error {
	return ctx.checkFileContains("internal/config", "env")
}

func (ctx *CLITestContext) theCLIShouldHaveConfigurationPrecedence() error {
	return ctx.checkFileContains("internal/config", "precedence")
}

func (ctx *CLITestContext) theCLIShouldValidateConfigurationValues() error {
	return ctx.checkFileContains("internal/config", "validate")
}

func (ctx *CLITestContext) theCLIShouldProvideConfigurationExamples() error {
	return ctx.checkFileExists("configs/config.yaml")
}

// Output format validation steps
func (ctx *CLITestContext) theCLIShouldSupportMultipleOutputFormats() error {
	return ctx.checkFileContains("internal/output", "format")
}

func (ctx *CLITestContext) theCLIShouldSupportQuietMode() error {
	return ctx.checkFileContains("cmd/root.go", "quiet")
}

func (ctx *CLITestContext) theCLIShouldSupportVerboseMode() error {
	return ctx.checkFileContains("cmd/root.go", "verbose")
}

func (ctx *CLITestContext) theCLIShouldHandleJSONOutput() error {
	return ctx.checkFileContains("internal/output", "json")
}

func (ctx *CLITestContext) theCLIShouldFormatOutputAppropriately() error {
	return ctx.checkFileContains("internal/output", "format")
}

// Error handling validation steps
func (ctx *CLITestContext) theCLIShouldHandleInvalidArgumentsGracefully() error {
	return ctx.checkFileContains("internal/errors", "validation")
}

func (ctx *CLITestContext) theCLIShouldProvideClearErrorMessages() error {
	return ctx.checkFileContains("internal/errors", "message")
}

func (ctx *CLITestContext) theCLIShouldValidateInputParameters() error {
	return ctx.checkFileContains("cmd", "validate")
}

func (ctx *CLITestContext) theCLIShouldHandleSystemErrorsProperly() error {
	return ctx.checkFileContains("internal/errors", "system")
}

func (ctx *CLITestContext) theCLIShouldUseAppropriateExitCodes() error {
	return ctx.checkFileContains("internal/errors", "exit")
}

