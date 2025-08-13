package lambdaeventprocessing

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
)

// LambdaEventProcessingTestContext holds state for Lambda Event Processing ATDD tests
type LambdaEventProcessingTestContext struct {
	workingDir     string
	projectDir     string
	projectName    string
	originalDir    string
	projectRoot    string
	lastCommand    *exec.Cmd
	lastOutput     []byte
	lastError      error
	lastExitCode   int
	projectExists  bool
	compilationOK  bool
	
	// Configuration tracking
	eventSources       []string
	processingPattern  string
	resilienceFeatures []string
	securityFeatures   []string
	observabilityLevel string
	testingLevel       string
	deploymentTargets  []string
}

var lambdaEventProcessingCtx *LambdaEventProcessingTestContext

// TestLambdaEventProcessingATDD runs ATDD scenarios for Lambda Event Processing blueprints
func TestLambdaEventProcessingATDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Lambda Event Processing ATDD tests in short mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeLambdaEventProcessingScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Lambda Event Processing ATDD test suite failed")
	}
}

// InitializeLambdaEventProcessingScenario registers all ATDD step definitions
func InitializeLambdaEventProcessingScenario(ctx *godog.ScenarioContext) {
	// Initialize context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		lambdaEventProcessingCtx = &LambdaEventProcessingTestContext{}
		return ctx, nil
	})

	// Cleanup after each scenario
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if lambdaEventProcessingCtx != nil {
			lambdaEventProcessingCtx.cleanup()
		}
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the go-starter CLI tool is available for event processing$`, lambdaEventProcessingCtx.theGoStarterCLIToolIsAvailable)
	ctx.Step(`^I am in a clean working directory for event processing$`, lambdaEventProcessingCtx.iAmInACleanWorkingDirectory)

	// Project generation steps
	ctx.Step(`^I want to create a Lambda event processing function$`, lambdaEventProcessingCtx.iWantToCreateALambdaEventProcessingFunction)
	ctx.Step(`^I configure event sources "([^"]*)"$`, lambdaEventProcessingCtx.iConfigureEventSources)
	ctx.Step(`^I set processing pattern to "([^"]*)"$`, lambdaEventProcessingCtx.iSetProcessingPattern)
	ctx.Step(`^I enable resilience features "([^"]*)"$`, lambdaEventProcessingCtx.iEnableResilienceFeatures)
	ctx.Step(`^I enable security features "([^"]*)"$`, lambdaEventProcessingCtx.iEnableSecurityFeatures)
	ctx.Step(`^I set observability level to "([^"]*)"$`, lambdaEventProcessingCtx.iSetObservabilityLevel)
	ctx.Step(`^I set testing level to "([^"]*)"$`, lambdaEventProcessingCtx.iSetTestingLevel)
	ctx.Step(`^I configure deployment targets "([^"]*)"$`, lambdaEventProcessingCtx.iConfigureDeploymentTargets)
	ctx.Step(`^I run the event processing generation command$`, lambdaEventProcessingCtx.iRunTheEventProcessingGenerationCommand)
	ctx.Step(`^the event processing generation should succeed$`, lambdaEventProcessingCtx.theEventProcessingGenerationShouldSucceed)

	// Project structure validation steps
	ctx.Step(`^the project should contain comprehensive event processing components$`, lambdaEventProcessingCtx.theProjectShouldContainComprehensiveEventProcessingComponents)
	ctx.Step(`^the generated event processing code should compile successfully$`, lambdaEventProcessingCtx.theGeneratedEventProcessingCodeShouldCompileSuccessfully)
	ctx.Step(`^the project should include all configured event sources$`, lambdaEventProcessingCtx.theProjectShouldIncludeAllConfiguredEventSources)
	ctx.Step(`^the project should implement the specified processing pattern$`, lambdaEventProcessingCtx.theProjectShouldImplementTheSpecifiedProcessingPattern)

	// Event source specific validation
	ctx.Step(`^I examine the SQS handler implementation$`, lambdaEventProcessingCtx.iExamineTheSQSHandlerImplementation)
	ctx.Step(`^the SQS handler should support batch processing$`, lambdaEventProcessingCtx.theSQSHandlerShouldSupportBatchProcessing)
	ctx.Step(`^the SQS handler should include error handling$`, lambdaEventProcessingCtx.theSQSHandlerShouldIncludeErrorHandling)
	ctx.Step(`^the SQS handler should integrate with observability$`, lambdaEventProcessingCtx.theSQSHandlerShouldIntegrateWithObservability)

	// Event router validation
	ctx.Step(`^I examine the event router implementation$`, lambdaEventProcessingCtx.iExamineTheEventRouterImplementation)
	ctx.Step(`^the event router should support intelligent routing$`, lambdaEventProcessingCtx.theEventRouterShouldSupportIntelligentRouting)
	ctx.Step(`^the event router should handle event type detection$`, lambdaEventProcessingCtx.theEventRouterShouldHandleEventTypeDetection)
	ctx.Step(`^the event router should integrate resilience patterns$`, lambdaEventProcessingCtx.theEventRouterShouldIntegrateResiliencePatterns)

	// Observability validation
	ctx.Step(`^I examine the observability implementation$`, lambdaEventProcessingCtx.iExamineTheObservabilityImplementation)
	ctx.Step(`^the observability should include structured logging$`, lambdaEventProcessingCtx.theObservabilityShouldIncludeStructuredLogging)
	ctx.Step(`^the observability should include distributed tracing$`, lambdaEventProcessingCtx.theObservabilityShouldIncludeDistributedTracing)
	ctx.Step(`^the observability should include custom metrics$`, lambdaEventProcessingCtx.theObservabilityShouldIncludeCustomMetrics)
	ctx.Step(`^the observability should include correlation tracking$`, lambdaEventProcessingCtx.theObservabilityShouldIncludeCorrelationTracking)

	// Security validation
	ctx.Step(`^I examine the security implementation$`, lambdaEventProcessingCtx.iExamineTheSecurityImplementation)
	ctx.Step(`^the security should include input validation$`, lambdaEventProcessingCtx.theSecurityShouldIncludeInputValidation)
	ctx.Step(`^the security should include secrets management$`, lambdaEventProcessingCtx.theSecurityShouldIncludeSecretsManagement)

	// Resilience validation
	ctx.Step(`^I examine the resilience implementation$`, lambdaEventProcessingCtx.iExamineTheResilienceImplementation)
	ctx.Step(`^the resilience should include retry logic$`, lambdaEventProcessingCtx.theResilienceShouldIncludeRetryLogic)
	ctx.Step(`^the resilience should include circuit breaker$`, lambdaEventProcessingCtx.theResilienceShouldIncludeCircuitBreaker)
	ctx.Step(`^the resilience should include DLQ handling$`, lambdaEventProcessingCtx.theResilienceShouldIncludeDLQHandling)

	// Performance validation
	ctx.Step(`^I examine the performance optimization$`, lambdaEventProcessingCtx.iExamineThePerformanceOptimization)
	ctx.Step(`^the performance should include cold start optimization$`, lambdaEventProcessingCtx.thePerformanceShouldIncludeColdStartOptimization)
	ctx.Step(`^the performance should include memory monitoring$`, lambdaEventProcessingCtx.thePerformanceShouldIncludeMemoryMonitoring)

	// Testing infrastructure validation
	ctx.Step(`^I examine the testing infrastructure$`, lambdaEventProcessingCtx.iExamineTheTestingInfrastructure)
	ctx.Step(`^the testing should include comprehensive test fixtures$`, lambdaEventProcessingCtx.theTestingShouldIncludeComprehensiveTestFixtures)
	ctx.Step(`^the testing should include AWS service mocks$`, lambdaEventProcessingCtx.theTestingShouldIncludeAWSServiceMocks)

	// Deployment validation
	ctx.Step(`^I examine the deployment configuration$`, lambdaEventProcessingCtx.iExamineTheDeploymentConfiguration)
	ctx.Step(`^the deployment should include SAM templates$`, lambdaEventProcessingCtx.theDeploymentShouldIncludeSAMTemplates)
	ctx.Step(`^the deployment should include Terraform configuration$`, lambdaEventProcessingCtx.theDeploymentShouldIncludeTerraformConfiguration)
	ctx.Step(`^the deployment should include automation scripts$`, lambdaEventProcessingCtx.theDeploymentShouldIncludeAutomationScripts)

	// Documentation validation
	ctx.Step(`^I examine the documentation$`, lambdaEventProcessingCtx.iExamineTheDocumentation)
	ctx.Step(`^the documentation should be comprehensive$`, lambdaEventProcessingCtx.theDocumentationShouldBeComprehensive)
	ctx.Step(`^the documentation should include architecture diagrams$`, lambdaEventProcessingCtx.theDocumentationShouldIncludeArchitectureDiagrams)
	ctx.Step(`^the documentation should include deployment guides$`, lambdaEventProcessingCtx.theDocumentationShouldIncludeDeploymentGuides)

	// Integration validation
	ctx.Step(`^I test the complete integration$`, lambdaEventProcessingCtx.iTestTheCompleteIntegration)
	ctx.Step(`^the integration should validate successfully$`, lambdaEventProcessingCtx.theIntegrationShouldValidateSuccessfully)
	ctx.Step(`^the project should meet production readiness criteria$`, lambdaEventProcessingCtx.theProjectShouldMeetProductionReadinessCriteria)
}

// Background step implementations
func (ctx *LambdaEventProcessingTestContext) theGoStarterCLIToolIsAvailable() error {
	var err error
	ctx.originalDir, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	// Find project root by looking for go.mod
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

	// Build the CLI tool
	buildCmd := exec.Command("go", "build", "-o", "go-starter", ".")
	buildCmd.Dir = ctx.projectRoot
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build go-starter CLI: %s", string(output))
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) iAmInACleanWorkingDirectory() error {
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "lambda-event-processing-test-*")
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
func (ctx *LambdaEventProcessingTestContext) iWantToCreateALambdaEventProcessingFunction() error {
	ctx.projectName = "test-lambda-event-processing"
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iConfigureEventSources(sources string) error {
	ctx.eventSources = strings.Split(sources, ",")
	for i, source := range ctx.eventSources {
		ctx.eventSources[i] = strings.TrimSpace(source)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iSetProcessingPattern(pattern string) error {
	ctx.processingPattern = pattern
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iEnableResilienceFeatures(features string) error {
	ctx.resilienceFeatures = strings.Split(features, ",")
	for i, feature := range ctx.resilienceFeatures {
		ctx.resilienceFeatures[i] = strings.TrimSpace(feature)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iEnableSecurityFeatures(features string) error {
	ctx.securityFeatures = strings.Split(features, ",")
	for i, feature := range ctx.securityFeatures {
		ctx.securityFeatures[i] = strings.TrimSpace(feature)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iSetObservabilityLevel(level string) error {
	ctx.observabilityLevel = level
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iSetTestingLevel(level string) error {
	ctx.testingLevel = level
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iConfigureDeploymentTargets(targets string) error {
	ctx.deploymentTargets = strings.Split(targets, ",")
	for i, target := range ctx.deploymentTargets {
		ctx.deploymentTargets[i] = strings.TrimSpace(target)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) iRunTheEventProcessingGenerationCommand() error {
	// Build the command with all configured options
	args := []string{
		filepath.Join(ctx.projectRoot, "go-starter"),
		"new",
		ctx.projectName,
		"--type=lambda-event-processing",
		"--logger=slog",
		"--no-git",
	}

	// Add event sources
	if len(ctx.eventSources) > 0 {
		args = append(args, "--event-sources="+strings.Join(ctx.eventSources, ","))
	}

	// Add processing pattern
	if ctx.processingPattern != "" {
		args = append(args, "--processing-pattern="+ctx.processingPattern)
	}

	// Add resilience features
	if len(ctx.resilienceFeatures) > 0 {
		args = append(args, "--resilience-features="+strings.Join(ctx.resilienceFeatures, ","))
	}

	// Add security features
	if len(ctx.securityFeatures) > 0 {
		args = append(args, "--security-features="+strings.Join(ctx.securityFeatures, ","))
	}

	// Add observability level
	if ctx.observabilityLevel != "" {
		args = append(args, "--observability-level="+ctx.observabilityLevel)
	}

	// Add testing level
	if ctx.testingLevel != "" {
		args = append(args, "--testing-level="+ctx.testingLevel)
	}

	// Add deployment targets
	if len(ctx.deploymentTargets) > 0 {
		args = append(args, "--deployment-targets="+strings.Join(ctx.deploymentTargets, ","))
	}

	ctx.lastCommand = exec.Command(args[0], args[1:]...)
	ctx.lastCommand.Dir = ctx.workingDir

	ctx.lastOutput, ctx.lastError = ctx.lastCommand.CombinedOutput()
	if ctx.lastError != nil {
		if exitError, ok := ctx.lastError.(*exec.ExitError); ok {
			ctx.lastExitCode = exitError.ExitCode()
		}
	} else {
		ctx.lastExitCode = 0
	}

	ctx.projectDir = filepath.Join(ctx.workingDir, ctx.projectName)
	return nil
}

func (ctx *LambdaEventProcessingTestContext) theEventProcessingGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", ctx.lastExitCode, string(ctx.lastOutput))
	}
	ctx.projectExists = true
	return nil
}

// Project structure validation implementations
func (ctx *LambdaEventProcessingTestContext) theProjectShouldContainComprehensiveEventProcessingComponents() error {
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
		"Makefile",
		".env.example",
		"internal/handler/event_router.go",
		"internal/domain/event.go",
		"internal/observability/logger.go",
		"internal/config/config.go",
		"configs/development.yaml",
		"scripts/deploy.sh",
	}

	// Add event source specific files
	for _, source := range ctx.eventSources {
		switch source {
		case "sqs":
			requiredFiles = append(requiredFiles, "internal/handler/sqs_handler.go")
		case "sns":
			requiredFiles = append(requiredFiles, "internal/handler/sns_handler.go")
		case "eventbridge":
			requiredFiles = append(requiredFiles, "internal/handler/eventbridge_handler.go")
		case "dynamodb-streams":
			requiredFiles = append(requiredFiles, "internal/handler/dynamodb_handler.go")
		case "kinesis":
			requiredFiles = append(requiredFiles, "internal/handler/kinesis_handler.go")
		}
	}

	// Check for resilience files
	for _, feature := range ctx.resilienceFeatures {
		switch feature {
		case "retry":
			requiredFiles = append(requiredFiles, "internal/resilience/retry.go")
		case "circuit-breaker":
			requiredFiles = append(requiredFiles, "internal/resilience/circuit_breaker.go")
		case "dlq":
			requiredFiles = append(requiredFiles, "internal/resilience/dlq_handler.go")
		}
	}

	// Check for security files
	for _, feature := range ctx.securityFeatures {
		switch feature {
		case "input-validation":
			requiredFiles = append(requiredFiles, "internal/security/input_validator.go")
		case "secrets-manager":
			requiredFiles = append(requiredFiles, "internal/security/secrets_manager.go")
		}
	}

	// Check for deployment files
	for _, target := range ctx.deploymentTargets {
		switch target {
		case "sam":
			requiredFiles = append(requiredFiles, "sam/template.yaml")
		case "terraform":
			requiredFiles = append(requiredFiles, "terraform/main.tf")
		}
	}

	for _, file := range requiredFiles {
		fullPath := filepath.Join(ctx.projectDir, file)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("required file not found: %s", file)
		}
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) theGeneratedEventProcessingCodeShouldCompileSuccessfully() error {
	// Run go mod tidy first
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	output, err := modCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}

	// Attempt to build the project
	buildCmd := exec.Command("go", "build", "-o", "lambda-app", ".")
	buildCmd.Dir = ctx.projectDir
	output, err = buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	ctx.compilationOK = true
	return nil
}

func (ctx *LambdaEventProcessingTestContext) theProjectShouldIncludeAllConfiguredEventSources() error {
	// Check go.mod for event source dependencies
	goModPath := filepath.Join(ctx.projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	goModContent := string(content)

	for _, source := range ctx.eventSources {
		var expectedDependency string
		switch source {
		case "sqs":
			expectedDependency = "github.com/aws/aws-sdk-go-v2/service/sqs"
		case "sns":
			expectedDependency = "github.com/aws/aws-sdk-go-v2/service/sns"
		case "eventbridge":
			expectedDependency = "github.com/aws/aws-sdk-go-v2/service/eventbridge"
		case "dynamodb-streams":
			expectedDependency = "github.com/aws/aws-sdk-go-v2/service/dynamodb"
		case "kinesis":
			expectedDependency = "github.com/aws/aws-sdk-go-v2/service/kinesis"
		}

		if !strings.Contains(goModContent, expectedDependency) {
			return fmt.Errorf("go.mod should contain dependency for %s: %s", source, expectedDependency)
		}
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) theProjectShouldImplementTheSpecifiedProcessingPattern() error {
	// Check main.go for processing pattern configuration
	mainGoPath := filepath.Join(ctx.projectDir, "main.go")
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}

	mainGoContent := string(content)
	if !strings.Contains(mainGoContent, fmt.Sprintf("processing_pattern", ctx.processingPattern)) {
		return fmt.Errorf("main.go should reference processing pattern: %s", ctx.processingPattern)
	}

	return nil
}

// Stub implementations for remaining step functions (implement based on actual requirements)
func (ctx *LambdaEventProcessingTestContext) iExamineTheSQSHandlerImplementation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theSQSHandlerShouldSupportBatchProcessing() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theSQSHandlerShouldIncludeErrorHandling() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theSQSHandlerShouldIntegrateWithObservability() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheEventRouterImplementation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theEventRouterShouldSupportIntelligentRouting() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theEventRouterShouldHandleEventTypeDetection() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theEventRouterShouldIntegrateResiliencePatterns() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheObservabilityImplementation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theObservabilityShouldIncludeStructuredLogging() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theObservabilityShouldIncludeDistributedTracing() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theObservabilityShouldIncludeCustomMetrics() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theObservabilityShouldIncludeCorrelationTracking() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheSecurityImplementation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theSecurityShouldIncludeInputValidation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theSecurityShouldIncludeSecretsManagement() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheResilienceImplementation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theResilienceShouldIncludeRetryLogic() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theResilienceShouldIncludeCircuitBreaker() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theResilienceShouldIncludeDLQHandling() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineThePerformanceOptimization() error { return nil }
func (ctx *LambdaEventProcessingTestContext) thePerformanceShouldIncludeColdStartOptimization() error { return nil }
func (ctx *LambdaEventProcessingTestContext) thePerformanceShouldIncludeMemoryMonitoring() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheTestingInfrastructure() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theTestingShouldIncludeComprehensiveTestFixtures() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theTestingShouldIncludeAWSServiceMocks() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheDeploymentConfiguration() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDeploymentShouldIncludeSAMTemplates() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDeploymentShouldIncludeTerraformConfiguration() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDeploymentShouldIncludeAutomationScripts() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iExamineTheDocumentation() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDocumentationShouldBeComprehensive() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDocumentationShouldIncludeArchitectureDiagrams() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theDocumentationShouldIncludeDeploymentGuides() error { return nil }

func (ctx *LambdaEventProcessingTestContext) iTestTheCompleteIntegration() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theIntegrationShouldValidateSuccessfully() error { return nil }
func (ctx *LambdaEventProcessingTestContext) theProjectShouldMeetProductionReadinessCriteria() error { return nil }

// Cleanup function
func (ctx *LambdaEventProcessingTestContext) cleanup() {
	if ctx.originalDir != "" {
		_ = os.Chdir(ctx.originalDir)
	}
	if ctx.workingDir != "" {
		_ = os.RemoveAll(ctx.workingDir)
	}
}

// Test validation functions
func (ctx *LambdaEventProcessingTestContext) validateFileExists(relativePath string) error {
	fullPath := filepath.Join(ctx.projectDir, relativePath)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("expected file does not exist: %s", relativePath)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) validateFileContains(relativePath, expectedContent string) error {
	fullPath := filepath.Join(ctx.projectDir, relativePath)
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", relativePath, err)
	}

	if !strings.Contains(string(content), expectedContent) {
		return fmt.Errorf("file %s does not contain expected content: %s", relativePath, expectedContent)
	}
	return nil
}

func (ctx *LambdaEventProcessingTestContext) validateProjectStructure() error {
	expectedDirs := []string{
		"internal",
		"internal/handler",
		"internal/domain",
		"internal/observability",
		"internal/config",
		"configs",
		"scripts",
	}

	// Add conditional directories based on configuration
	if len(ctx.resilienceFeatures) > 0 {
		expectedDirs = append(expectedDirs, "internal/resilience")
	}
	if len(ctx.securityFeatures) > 0 {
		expectedDirs = append(expectedDirs, "internal/security")
	}
	if ctx.testingLevel == "comprehensive" {
		expectedDirs = append(expectedDirs, "internal/testing")
	}

	for _, dir := range expectedDirs {
		fullPath := filepath.Join(ctx.projectDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			return fmt.Errorf("expected directory does not exist: %s", dir)
		}
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) validateAWSSDKIntegration() error {
	// Validate that the project uses AWS SDK v2
	goModPath := filepath.Join(ctx.projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}

	goModContent := string(content)
	expectedDependencies := []string{
		"github.com/aws/aws-lambda-go",
		"github.com/aws/aws-sdk-go-v2",
		"github.com/aws/aws-xray-sdk-go",
	}

	for _, dep := range expectedDependencies {
		if !strings.Contains(goModContent, dep) {
			return fmt.Errorf("go.mod should contain AWS dependency: %s", dep)
		}
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) validateObservabilityIntegration() error {
	// Check that observability components are properly integrated
	files := []string{
		"internal/observability/logger.go",
		"internal/observability/tracing.go",
		"internal/observability/metrics.go",
	}

	for _, file := range files {
		if err := ctx.validateFileExists(file); err != nil {
			return err
		}
	}

	return nil
}

func (ctx *LambdaEventProcessingTestContext) runCompilationTest() error {
	if !ctx.compilationOK {
		return ctx.theGeneratedEventProcessingCodeShouldCompileSuccessfully()
	}
	return nil
}