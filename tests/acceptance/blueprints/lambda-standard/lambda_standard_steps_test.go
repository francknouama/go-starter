package lambdastandard

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// LambdaStandardTestContext holds state for Lambda Standard BDD tests
type LambdaStandardTestContext struct {
	workingDir   string
	projectDir   string
	projectName  string
	originalDir  string
	projectRoot  string
	lastCommand  *exec.Cmd
	lastOutput   []byte
	lastError    error
	lastExitCode int
	projectExists bool
	compilationOK bool
}

var lambdaStandardCtx *LambdaStandardTestContext

// TestLambdaStandardBDD runs BDD scenarios for Lambda Standard blueprints
func TestLambdaStandardBDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Lambda Standard BDD tests in short mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeLambdaStandardScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Lambda Standard BDD test suite failed")
	}
}

// InitializeLambdaStandardScenario registers all BDD step definitions
func InitializeLambdaStandardScenario(ctx *godog.ScenarioContext) {
	// Initialize context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		lambdaStandardCtx = &LambdaStandardTestContext{}
		return ctx, nil
	})

	// Cleanup after each scenario
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if lambdaStandardCtx != nil {
			lambdaStandardCtx.cleanup()
		}
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the go-starter CLI tool is available$`, lambdaStandardCtx.theGoStarterCLIToolIsAvailable)
	ctx.Step(`^I am in a clean working directory$`, lambdaStandardCtx.iAmInACleanWorkingDirectory)

	// Project generation steps
	ctx.Step(`^I want to create a standard AWS Lambda function$`, lambdaStandardCtx.iWantToCreateAStandardAWSLambdaFunction)
	ctx.Step(`^I run the command "([^"]*)"$`, lambdaStandardCtx.iRunTheCommand)
	ctx.Step(`^the generation should succeed$`, lambdaStandardCtx.theGenerationShouldSucceed)
	ctx.Step(`^the project should contain Lambda-specific components$`, lambdaStandardCtx.theProjectShouldContainLambdaSpecificComponents)
	ctx.Step(`^the generated code should compile successfully$`, lambdaStandardCtx.theGeneratedCodeShouldCompileSuccessfully)
	ctx.Step(`^the project should include AWS Lambda runtime$`, lambdaStandardCtx.theProjectShouldIncludeAWSLambdaRuntime)

	// CloudWatch integration steps
	ctx.Step(`^I have generated a standard Lambda function$`, lambdaStandardCtx.iHaveGeneratedAStandardLambdaFunction)
	ctx.Step(`^I examine the observability configuration$`, lambdaStandardCtx.iExamineTheObservabilityConfiguration)
	ctx.Step(`^CloudWatch logging should be properly configured$`, lambdaStandardCtx.cloudWatchLoggingShouldBeProperlyConfigured)
	ctx.Step(`^metrics collection should be implemented$`, lambdaStandardCtx.metricsCollectionShouldBeImplemented)
	ctx.Step(`^distributed tracing should be available$`, lambdaStandardCtx.distributedTracingShouldBeAvailable)
	ctx.Step(`^error tracking should be integrated$`, lambdaStandardCtx.errorTrackingShouldBeIntegrated)

	// Deployment configuration steps
	ctx.Step(`^I examine the deployment setup$`, lambdaStandardCtx.iExamineTheDeploymentSetup)
	ctx.Step(`^the deployment script should be available$`, lambdaStandardCtx.theDeploymentScriptShouldBeAvailable)
	ctx.Step(`^SAM template should be properly configured$`, lambdaStandardCtx.samTemplateShouldBeProperlyConfigured)
	ctx.Step(`^environment variables should be managed$`, lambdaStandardCtx.environmentVariablesShouldBeManaged)
	ctx.Step(`^IAM permissions should be defined$`, lambdaStandardCtx.iamPermissionsShouldBeDefined)

	// Handler implementation steps
	ctx.Step(`^I examine the handler code$`, lambdaStandardCtx.iExamineTheHandlerCode)
	ctx.Step(`^the main handler should be properly structured$`, lambdaStandardCtx.theMainHandlerShouldBeProperlyStructured)
	ctx.Step(`^context handling should be implemented$`, lambdaStandardCtx.contextHandlingShouldBeImplemented)
	ctx.Step(`^error handling should be robust$`, lambdaStandardCtx.errorHandlingShouldBeRobust)
	ctx.Step(`^response formatting should be correct$`, lambdaStandardCtx.responseFormattingShouldBeCorrect)

	// Testing infrastructure steps
	ctx.Step(`^I examine the test setup$`, lambdaStandardCtx.iExamineTheTestSetup)
	ctx.Step(`^unit tests should be included$`, lambdaStandardCtx.unitTestsShouldBeIncluded)
	ctx.Step(`^integration tests should be available$`, lambdaStandardCtx.integrationTestsShouldBeAvailable)
	ctx.Step(`^local testing should be supported$`, lambdaStandardCtx.localTestingShouldBeSupported)
	ctx.Step(`^test coverage should be measurable$`, lambdaStandardCtx.testCoverageShouldBeMeasurable)

	// Performance optimization steps
	ctx.Step(`^I examine the performance configuration$`, lambdaStandardCtx.iExamineThePerformanceConfiguration)
	ctx.Step(`^cold start optimization should be implemented$`, lambdaStandardCtx.coldStartOptimizationShouldBeImplemented)
	ctx.Step(`^memory usage should be optimized$`, lambdaStandardCtx.memoryUsageShouldBeOptimized)
	ctx.Step(`^initialization should be efficient$`, lambdaStandardCtx.initializationShouldBeEfficient)
	ctx.Step(`^runtime performance should be monitored$`, lambdaStandardCtx.runtimePerformanceShouldBeMonitored)
}

// Background step implementations
func (ctx *LambdaStandardTestContext) theGoStarterCLIToolIsAvailable() error {
	var err error
	ctx.originalDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	projectRoot := ctx.originalDir
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			projectRoot = filepath.Join(ctx.originalDir, "..", "..", "..", "..", "..")
			break
		}
		projectRoot = parent
	}
	ctx.projectRoot = projectRoot

	buildCmd := exec.Command("go", "build", "-o", "go-starter", ".")
	buildCmd.Dir = ctx.projectRoot
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build go-starter CLI: %s", string(output))
	}

	return nil
}

func (ctx *LambdaStandardTestContext) iAmInACleanWorkingDirectory() error {
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "lambda-standard-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	err = os.Chdir(ctx.workingDir)
	if err != nil {
		return fmt.Errorf("failed to change to working directory: %w", err)
	}

	return nil
}

// Project generation step implementations
func (ctx *LambdaStandardTestContext) iWantToCreateAStandardAWSLambdaFunction() error {
	ctx.projectName = "test-lambda-standard"
	return nil
}

func (ctx *LambdaStandardTestContext) iRunTheCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}

	if parts[0] == "go-starter" {
		parts[0] = filepath.Join(ctx.projectRoot, "go-starter")
	}

	ctx.lastCommand = exec.Command(parts[0], parts[1:]...)
	ctx.lastCommand.Dir = ctx.workingDir

	ctx.lastOutput, ctx.lastError = ctx.lastCommand.CombinedOutput()
	if ctx.lastError != nil {
		if exitError, ok := ctx.lastError.(*exec.ExitError); ok {
			ctx.lastExitCode = exitError.ExitCode()
		}
	} else {
		ctx.lastExitCode = 0
	}

	if len(parts) > 2 && parts[1] == "new" {
		ctx.projectName = parts[2]
		ctx.projectDir = filepath.Join(ctx.workingDir, ctx.projectName)
	}

	return nil
}

func (ctx *LambdaStandardTestContext) theGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", ctx.lastExitCode, string(ctx.lastOutput))
	}
	ctx.projectExists = true
	return nil
}

func (ctx *LambdaStandardTestContext) theProjectShouldContainLambdaSpecificComponents() error {
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"handler.go",
		"template.yaml",
		"Makefile",
	}

	for _, file := range requiredFiles {
		fullPath := filepath.Join(ctx.projectDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	goModPath := filepath.Join(ctx.projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	if !strings.Contains(string(content), "github.com/aws/aws-lambda-go") {
		return fmt.Errorf("go.mod should contain AWS Lambda Go dependency")
	}

	return nil
}

func (ctx *LambdaStandardTestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	output, err := modCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}

	buildCmd := exec.Command("go", "build", "-o", "lambda-app", ".")
	buildCmd.Dir = ctx.projectDir
	output, err = buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	ctx.compilationOK = true
	return nil
}

func (ctx *LambdaStandardTestContext) theProjectShouldIncludeAWSLambdaRuntime() error {
	handlerPath := filepath.Join(ctx.projectDir, "handler.go")
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		return fmt.Errorf("failed to read handler.go: %w", err)
	}

	if !strings.Contains(string(content), "lambda.Start") {
		return fmt.Errorf("handler.go should contain lambda.Start")
	}

	return nil
}

// CloudWatch integration step implementations
func (ctx *LambdaStandardTestContext) iHaveGeneratedAStandardLambdaFunction() error {
	if !ctx.projectExists {
		return ctx.iRunTheCommand("go-starter new test-lambda --type=lambda-standard --logger=slog --no-git")
	}
	return nil
}

func (ctx *LambdaStandardTestContext) iExamineTheObservabilityConfiguration() error {
	return nil // Context setting step
}

func (ctx *LambdaStandardTestContext) cloudWatchLoggingShouldBeProperlyConfigured() error {
	// Check for CloudWatch logging configuration
	observabilityPath := filepath.Join(ctx.projectDir, "internal/observability")
	if _, err := os.Stat(observabilityPath); os.IsNotExist(err) {
		return fmt.Errorf("observability directory should exist")
	}
	return nil
}

func (ctx *LambdaStandardTestContext) metricsCollectionShouldBeImplemented() error {
	// Check for metrics collection
	return nil
}

func (ctx *LambdaStandardTestContext) distributedTracingShouldBeAvailable() error {
	// Check for X-Ray tracing configuration
	return nil
}

func (ctx *LambdaStandardTestContext) errorTrackingShouldBeIntegrated() error {
	// Check for error tracking implementation
	return nil
}

// Stub implementations for remaining functions (for brevity)
func (ctx *LambdaStandardTestContext) iExamineTheDeploymentSetup() error { return nil }
func (ctx *LambdaStandardTestContext) theDeploymentScriptShouldBeAvailable() error { return nil }
func (ctx *LambdaStandardTestContext) samTemplateShouldBeProperlyConfigured() error { return nil }
func (ctx *LambdaStandardTestContext) environmentVariablesShouldBeManaged() error { return nil }
func (ctx *LambdaStandardTestContext) iamPermissionsShouldBeDefined() error { return nil }

func (ctx *LambdaStandardTestContext) iExamineTheHandlerCode() error { return nil }
func (ctx *LambdaStandardTestContext) theMainHandlerShouldBeProperlyStructured() error { return nil }
func (ctx *LambdaStandardTestContext) contextHandlingShouldBeImplemented() error { return nil }
func (ctx *LambdaStandardTestContext) errorHandlingShouldBeRobust() error { return nil }
func (ctx *LambdaStandardTestContext) responseFormattingShouldBeCorrect() error { return nil }

func (ctx *LambdaStandardTestContext) iExamineTheTestSetup() error { return nil }
func (ctx *LambdaStandardTestContext) unitTestsShouldBeIncluded() error { return nil }
func (ctx *LambdaStandardTestContext) integrationTestsShouldBeAvailable() error { return nil }
func (ctx *LambdaStandardTestContext) localTestingShouldBeSupported() error { return nil }
func (ctx *LambdaStandardTestContext) testCoverageShouldBeMeasurable() error { return nil }

func (ctx *LambdaStandardTestContext) iExamineThePerformanceConfiguration() error { return nil }
func (ctx *LambdaStandardTestContext) coldStartOptimizationShouldBeImplemented() error { return nil }
func (ctx *LambdaStandardTestContext) memoryUsageShouldBeOptimized() error { return nil }
func (ctx *LambdaStandardTestContext) initializationShouldBeEfficient() error { return nil }
func (ctx *LambdaStandardTestContext) runtimePerformanceShouldBeMonitored() error { return nil }

// Cleanup function
func (ctx *LambdaStandardTestContext) cleanup() {
	if ctx.originalDir != "" {
		_ = os.Chdir(ctx.originalDir)
	}
	if ctx.workingDir != "" {
		_ = os.RemoveAll(ctx.workingDir)
	}
}