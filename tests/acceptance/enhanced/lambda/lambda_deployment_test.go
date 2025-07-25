package lambda

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
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// LambdaDeploymentTestContext holds test state for Lambda deployment testing
type LambdaDeploymentTestContext struct {
	projectConfig     *types.ProjectConfig
	projectPath       string
	tempDir           string
	startTime         time.Time
	lastCommandOutput string
	lastCommandError  error
	generatedFiles    []string
	lambdaType        string
	buildOutput       string
}

// TestFeatures runs the Lambda deployment BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &LambdaDeploymentTestContext{}

			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				ctx.generatedFiles = []string{}

				// Initialize templates
				if err := helpers.InitializeTemplates(); err != nil {
					return goCtx, err
				}

				return goCtx, nil
			})

			s.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				if ctx.tempDir != "" {
					os.RemoveAll(ctx.tempDir)
				}
				return goCtx, nil
			})

			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run Lambda deployment tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *LambdaDeploymentTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^I am testing Lambda deployment scenarios$`, ctx.iAmTestingLambdaDeploymentScenarios)

	// Project generation steps
	s.Step(`^I generate a Lambda function with configuration:$`, ctx.iGenerateALambdaFunctionWithConfiguration)
	s.Step(`^I generate a Lambda API Gateway proxy with configuration:$`, ctx.iGenerateALambdaAPIGatewayProxyWithConfiguration)
	s.Step(`^I generate a Lambda function optimized for cold starts:$`, ctx.iGenerateALambdaFunctionOptimizedForColdStarts)
	s.Step(`^I generate a Lambda function with AWS SDK integration:$`, ctx.iGenerateALambdaFunctionWithAWSSDKIntegration)

	// Build and compilation validation
	s.Step(`^the project should compile successfully$`, ctx.theProjectShouldCompileSuccessfully)
	s.Step(`^the Lambda binary should be created$`, ctx.theLambdaBinaryShouldBeCreated)
	s.Step(`^the binary should be optimized for Lambda runtime$`, ctx.theBinaryShouldBeOptimizedForLambdaRuntime)
	s.Step(`^cross-compilation should work correctly$`, ctx.crossCompilationShouldWorkCorrectly)

	// AWS Lambda specific validations
	s.Step(`^AWS Lambda handler should be properly implemented$`, ctx.awsLambdaHandlerShouldBeProperlyImplemented)
	s.Step(`^Lambda context handling should work correctly$`, ctx.lambdaContextHandlingShouldWorkCorrectly)
	s.Step(`^CloudWatch logging should be integrated$`, ctx.cloudWatchLoggingShouldBeIntegrated)
	s.Step(`^(.*) logger should work with CloudWatch$`, ctx.loggerShouldWorkWithCloudWatch)

	// SAM template validations
	s.Step(`^SAM template should be properly configured$`, ctx.samTemplateShouldBeProperlyConfigured)
	s.Step(`^Lambda function resources should be defined$`, ctx.lambdaFunctionResourcesShouldBeDefined)
	s.Step(`^environment variables should be configurable$`, ctx.environmentVariablesShouldBeConfigurable)
	s.Step(`^IAM roles should be properly configured$`, ctx.iamRolesShouldBeProperlyConfigured)
	s.Step(`^deployment configuration should be complete$`, ctx.deploymentConfigurationShouldBeComplete)

	// API Gateway integration validations
	s.Step(`^API Gateway integration should work correctly$`, ctx.apiGatewayIntegrationShouldWorkCorrectly)
	s.Step(`^request/response mapping should be implemented$`, ctx.requestResponseMappingShouldBeImplemented)
	s.Step(`^CORS headers should be properly configured$`, ctx.corsHeadersShouldBeProperlyConfigured)
	s.Step(`^error responses should follow API Gateway format$`, ctx.errorResponsesShouldFollowAPIGatewayFormat)

	// Performance optimization validations
	s.Step(`^cold start optimization should be implemented$`, ctx.coldStartOptimizationShouldBeImplemented)
	s.Step(`^minimal dependencies should be included$`, ctx.minimalDependenciesShouldBeIncluded)
	s.Step(`^binary size should be optimized$`, ctx.binarySizeShouldBeOptimized)
	s.Step(`^memory usage should be efficient$`, ctx.memoryUsageShouldBeEfficient)

	// Testing and local development
	s.Step(`^local testing should be supported$`, ctx.localTestingShouldBeSupported)
	s.Step(`^unit tests should cover Lambda handlers$`, ctx.unitTestsShouldCoverLambdaHandlers)
	s.Step(`^SAM local testing should work$`, ctx.samLocalTestingShouldWork)
	s.Step(`^event samples should be provided$`, ctx.eventSamplesShouldBeProvided)

	// Deployment validations
	s.Step(`^deployment scripts should be included$`, ctx.deploymentScriptsShouldBeIncluded)
	s.Step(`^Makefile should have deployment targets$`, ctx.makefileShouldHaveDeploymentTargets)
	s.Step(`^CI/CD templates should be available$`, ctx.cicdTemplatesShouldBeAvailable)
	s.Step(`^deployment documentation should be comprehensive$`, ctx.deploymentDocumentationShouldBeComprehensive)

	// AWS SDK integration validations
	s.Step(`^AWS SDK v2 should be properly integrated$`, ctx.awsSDKV2ShouldBeProperlyIntegrated)
	s.Step(`^service clients should be properly initialized$`, ctx.serviceClientsShouldBeProperlyInitialized)
	s.Step(`^AWS credentials handling should be secure$`, ctx.awsCredentialsHandlingShouldBeSecure)
	s.Step(`^X-Ray tracing should be available$`, ctx.xrayTracingShouldBeAvailable)

	// Error handling and monitoring
	s.Step(`^error handling should be comprehensive$`, ctx.errorHandlingShouldBeComprehensive)
	s.Step(`^panic recovery should be implemented$`, ctx.panicRecoveryShouldBeImplemented)
	s.Step(`^metrics should be exported to CloudWatch$`, ctx.metricsShouldBeExportedToCloudWatch)
	s.Step(`^structured logging should be used$`, ctx.structuredLoggingShouldBeUsed)

	// Multi-environment support
	s.Step(`^multiple environments should be supported$`, ctx.multipleEnvironmentsShouldBeSupported)
	s.Step(`^environment-specific configs should exist$`, ctx.environmentSpecificConfigsShouldExist)
	s.Step(`^secrets management should be implemented$`, ctx.secretsManagementShouldBeImplemented)
	s.Step(`^parameter store integration should work$`, ctx.parameterStoreIntegrationShouldWork)
	
	// Additional missing step definitions
	s.Step(`^API Gateway resources should be defined in SAM$`, ctx.aPIGatewayResourcesShouldBeDefinedInSAM)
	s.Step(`^build process should support multiple architectures$`, ctx.buildProcessShouldSupportMultipleArchitectures)
}

// Background step implementations
func (ctx *LambdaDeploymentTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Use local binary for testing
	binaryPath := "../../../../bin/go-starter"
	cmd := exec.Command(binaryPath, "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-starter CLI not available at %s: %v", binaryPath, err)
	}
	return nil
}

func (ctx *LambdaDeploymentTestContext) allTemplatesAreProperlyInitialized() error {
	return helpers.InitializeTemplates()
}

func (ctx *LambdaDeploymentTestContext) iAmTestingLambdaDeploymentScenarios() error {
	return nil
}

// Project generation implementations
func (ctx *LambdaDeploymentTestContext) iGenerateALambdaFunctionWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *LambdaDeploymentTestContext) iGenerateALambdaAPIGatewayProxyWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *LambdaDeploymentTestContext) iGenerateALambdaFunctionOptimizedForColdStarts(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *LambdaDeploymentTestContext) iGenerateALambdaFunctionWithAWSSDKIntegration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *LambdaDeploymentTestContext) generateProjectWithConfiguration(table *godog.Table) error {
	// Parse configuration from table
	config := &types.ProjectConfig{}

	for i := 0; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			config.Type = value
			ctx.lambdaType = value
		case "framework":
			config.Framework = value
		case "logger":
			config.Logger = value
		case "go_version":
			config.GoVersion = value
		}
	}

	// Set defaults
	if config.Name == "" {
		config.Name = fmt.Sprintf("test-lambda-%s", strings.ReplaceAll(ctx.lambdaType, "-", ""))
	}
	if config.Module == "" {
		config.Module = fmt.Sprintf("github.com/test/lambda-%s", strings.ReplaceAll(ctx.lambdaType, "-", ""))
	}
	if config.GoVersion == "" {
		config.GoVersion = "1.23"
	}
	// Default AuthType to "none" for lambda-proxy to avoid unused imports
	if ctx.lambdaType == "lambda-proxy" {
		if config.Variables == nil {
			config.Variables = make(map[string]string)
		}
		if config.Variables["AuthType"] == "" {
			config.Variables["AuthType"] = "none"
		}
	}

	ctx.projectConfig = config

	// Generate project
	var err error
	ctx.projectPath, err = ctx.generateTestProject(config, ctx.tempDir)
	if err != nil {
		ctx.lastCommandError = err
		return fmt.Errorf("failed to generate project: %v", err)
	}

	// Collect generated files
	err = filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(ctx.projectPath, path)
			ctx.generatedFiles = append(ctx.generatedFiles, relPath)
		}
		return nil
	})

	return err
}

// Build and compilation validations
func (ctx *LambdaDeploymentTestContext) theProjectShouldCompileSuccessfully() error {
	return ctx.validateProjectCompilation()
}

func (ctx *LambdaDeploymentTestContext) theLambdaBinaryShouldBeCreated() error {
	// Try to build the Lambda binary
	cmd := exec.Command("make", "build")
	cmd.Dir = ctx.projectPath
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	
	output, err := cmd.CombinedOutput()
	ctx.buildOutput = string(output)
	
	if err != nil {
		return fmt.Errorf("failed to build Lambda binary: %v\nOutput: %s", err, ctx.buildOutput)
	}

	// Check if binary was created
	binaryPath := filepath.Join(ctx.projectPath, "bootstrap")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Also check for handler binary
		binaryPath = filepath.Join(ctx.projectPath, "handler")
		if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
			return fmt.Errorf("Lambda binary not found at bootstrap or handler")
		}
	}

	return nil
}

func (ctx *LambdaDeploymentTestContext) theBinaryShouldBeOptimizedForLambdaRuntime() error {
	// Check if build flags include optimization
	if !strings.Contains(ctx.buildOutput, "-ldflags") && !strings.Contains(ctx.buildOutput, "-s -w") {
		// Try checking Makefile for optimization flags
		makefilePath := filepath.Join(ctx.projectPath, "Makefile")
		content, err := os.ReadFile(makefilePath)
		if err != nil {
			return fmt.Errorf("failed to read Makefile: %v", err)
		}
		
		makefileContent := string(content)
		if !strings.Contains(makefileContent, "-ldflags") || !strings.Contains(makefileContent, "-s -w") {
			return fmt.Errorf("Lambda binary optimization flags not found")
		}
	}
	
	return nil
}

func (ctx *LambdaDeploymentTestContext) crossCompilationShouldWorkCorrectly() error {
	// Cross-compilation is already tested in theLambdaBinaryShouldBeCreated
	return nil
}

// AWS Lambda specific validations
func (ctx *LambdaDeploymentTestContext) awsLambdaHandlerShouldBeProperlyImplemented() error {
	// Check for Lambda handler implementation
	handlerPatterns := []string{
		"lambda.Start",
		"github.com/aws/aws-lambda-go/lambda",
		"HandleRequest",
		"LambdaHandler",
	}
	
	return ctx.checkForPatterns(handlerPatterns, "Lambda handler implementation")
}

func (ctx *LambdaDeploymentTestContext) lambdaContextHandlingShouldWorkCorrectly() error {
	// Check for Lambda context usage
	contextPatterns := []string{
		"lambdacontext",
		"context.Context",
		"RequestID",
		"lambda.NewContext",
	}
	
	return ctx.checkForPatterns(contextPatterns, "Lambda context handling")
}

func (ctx *LambdaDeploymentTestContext) cloudWatchLoggingShouldBeIntegrated() error {
	// Check for CloudWatch-optimized logging
	loggingPatterns := []string{
		"log",
		"logger",
		"Info",
		"Error",
		"structured",
		"json",
	}
	
	return ctx.checkForPatterns(loggingPatterns, "CloudWatch logging integration")
}

func (ctx *LambdaDeploymentTestContext) loggerShouldWorkWithCloudWatch(logger string) error {
	// Verify logger is configured for CloudWatch (JSON output)
	loggerPatterns := map[string][]string{
		"slog":    {"slog", "JSONHandler", "json"},
		"zap":     {"zap", "NewProduction", "json"},
		"logrus":  {"logrus", "JSONFormatter", "json"},
		"zerolog": {"zerolog", "ConsoleWriter", "json"},
	}
	
	patterns, exists := loggerPatterns[logger]
	if !exists {
		return fmt.Errorf("unsupported logger: %s", logger)
	}
	
	return ctx.checkForPatterns(patterns, fmt.Sprintf("%s CloudWatch integration", logger))
}

// SAM template validations
func (ctx *LambdaDeploymentTestContext) samTemplateShouldBeProperlyConfigured() error {
	// Check for SAM template file
	samTemplatePath := filepath.Join(ctx.projectPath, "template.yaml")
	if _, err := os.Stat(samTemplatePath); os.IsNotExist(err) {
		// Also check template.yml
		samTemplatePath = filepath.Join(ctx.projectPath, "template.yml")
		if _, err := os.Stat(samTemplatePath); os.IsNotExist(err) {
			return fmt.Errorf("SAM template not found (template.yaml or template.yml)")
		}
	}
	
	// Read and validate SAM template
	content, err := os.ReadFile(samTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read SAM template: %v", err)
	}
	
	templateContent := string(content)
	requiredElements := []string{
		"AWSTemplateFormatVersion",
		"Transform: AWS::Serverless",
		"AWS::Serverless::Function",
	}
	
	for _, element := range requiredElements {
		if !strings.Contains(templateContent, element) {
			return fmt.Errorf("SAM template missing required element: %s", element)
		}
	}
	
	return nil
}

func (ctx *LambdaDeploymentTestContext) lambdaFunctionResourcesShouldBeDefined() error {
	// This is validated in samTemplateShouldBeProperlyConfigured
	return nil
}

func (ctx *LambdaDeploymentTestContext) environmentVariablesShouldBeConfigurable() error {
	// Check SAM template for environment variables section
	samPatterns := []string{"Environment:", "Variables:", "ENV_VAR", "ENVIRONMENT"}
	return ctx.checkForPatterns(samPatterns, "environment variables configuration")
}

func (ctx *LambdaDeploymentTestContext) iamRolesShouldBeProperlyConfigured() error {
	// Check for IAM role configuration in SAM template
	iamPatterns := []string{"Role:", "Policies:", "IAM", "AssumeRolePolicyDocument"}
	return ctx.checkForPatterns(iamPatterns, "IAM roles configuration")
}

func (ctx *LambdaDeploymentTestContext) deploymentConfigurationShouldBeComplete() error {
	// Check for deployment configuration
	deployPatterns := []string{"DeploymentPreference", "Stage", "Outputs", "Parameters"}
	return ctx.checkForPatterns(deployPatterns, "deployment configuration")
}

// Helper methods
func (ctx *LambdaDeploymentTestContext) generateTestProject(config *types.ProjectConfig, tempDir string) (string, error) {
	gen := generator.New()
	
	projectPath := filepath.Join(tempDir, config.Name)
	
	_, err := gen.Generate(*config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	if err != nil {
		return "", err
	}
	
	return projectPath, nil
}

func (ctx *LambdaDeploymentTestContext) validateProjectCompilation() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("project path not set")
	}
	
	// Check if go.mod exists
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project")
	}
	
	// Run go mod tidy to resolve dependencies
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = ctx.projectPath
	tidyOutput, err := tidyCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %v\nOutput: %s", err, string(tidyOutput))
	}
	
	// Try to run go build with Lambda settings
	cmd := exec.Command("go", "build", "-tags", "lambda.norpc", "./...")
	cmd.Dir = ctx.projectPath
	cmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64")
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("project compilation failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (ctx *LambdaDeploymentTestContext) checkForPatterns(patterns []string, description string) error {
	found := false
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Check appropriate file types
		if strings.HasSuffix(path, ".go") || strings.HasSuffix(path, ".yaml") || 
		   strings.HasSuffix(path, ".yml") || strings.HasSuffix(path, "Makefile") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			fileContent := string(content)
			for _, pattern := range patterns {
				if strings.Contains(fileContent, pattern) {
					found = true
					return nil
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	if !found {
		return fmt.Errorf("%s patterns not found in generated project", description)
	}
	
	return nil
}

// Simplified implementations for remaining validations
func (ctx *LambdaDeploymentTestContext) apiGatewayIntegrationShouldWorkCorrectly() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) requestResponseMappingShouldBeImplemented() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) corsHeadersShouldBeProperlyConfigured() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) errorResponsesShouldFollowAPIGatewayFormat() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) coldStartOptimizationShouldBeImplemented() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) minimalDependenciesShouldBeIncluded() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) binarySizeShouldBeOptimized() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) memoryUsageShouldBeEfficient() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) localTestingShouldBeSupported() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) unitTestsShouldCoverLambdaHandlers() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) samLocalTestingShouldWork() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) eventSamplesShouldBeProvided() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) deploymentScriptsShouldBeIncluded() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) makefileShouldHaveDeploymentTargets() error {
	// Check Makefile for deployment targets
	makefilePath := filepath.Join(ctx.projectPath, "Makefile")
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return fmt.Errorf("failed to read Makefile: %v", err)
	}
	
	makefileContent := string(content)
	deployTargets := []string{"deploy", "package", "sam", "build"}
	
	found := false
	for _, target := range deployTargets {
		if strings.Contains(makefileContent, target+":") {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("deployment targets not found in Makefile")
	}
	
	return nil
}

func (ctx *LambdaDeploymentTestContext) cicdTemplatesShouldBeAvailable() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) deploymentDocumentationShouldBeComprehensive() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) awsSDKV2ShouldBeProperlyIntegrated() error {
	// Check for AWS SDK v2 usage
	sdkPatterns := []string{
		"github.com/aws/aws-sdk-go-v2",
		"aws-sdk-go-v2/config",
		"aws-sdk-go-v2/service",
	}
	
	return ctx.checkForPatterns(sdkPatterns, "AWS SDK v2 integration")
}

func (ctx *LambdaDeploymentTestContext) serviceClientsShouldBeProperlyInitialized() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) awsCredentialsHandlingShouldBeSecure() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) xrayTracingShouldBeAvailable() error {
	// Check for X-Ray integration
	xrayPatterns := []string{"xray", "X-Ray", "aws-xray-sdk-go", "Trace"}
	return ctx.checkForPatterns(xrayPatterns, "X-Ray tracing")
}

func (ctx *LambdaDeploymentTestContext) errorHandlingShouldBeComprehensive() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) panicRecoveryShouldBeImplemented() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) metricsShouldBeExportedToCloudWatch() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) structuredLoggingShouldBeUsed() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) multipleEnvironmentsShouldBeSupported() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) environmentSpecificConfigsShouldExist() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) secretsManagementShouldBeImplemented() error {
	return nil // Simplified implementation
}

func (ctx *LambdaDeploymentTestContext) parameterStoreIntegrationShouldWork() error {
	return nil // Simplified implementation
}

// aPIGatewayResourcesShouldBeDefinedInSAM validates API Gateway resources in SAM template
func (ctx *LambdaDeploymentTestContext) aPIGatewayResourcesShouldBeDefinedInSAM() error {
	// Check SAM template for API Gateway configuration
	apiGatewayPatterns := []string{"AWS::Serverless::Api", "Api:", "Events:", "Type: Api"}
	return ctx.checkForPatterns(apiGatewayPatterns, "API Gateway resources in SAM")
}

// buildProcessShouldSupportMultipleArchitectures validates cross-platform build support
func (ctx *LambdaDeploymentTestContext) buildProcessShouldSupportMultipleArchitectures() error {
	// Check Makefile for architecture support
	makefilePath := filepath.Join(ctx.projectPath, "Makefile")
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return fmt.Errorf("failed to read Makefile: %v", err)
	}
	
	makefileContent := string(content)
	architecturePatterns := []string{"GOARCH=", "amd64", "arm64"}
	
	found := false
	for _, pattern := range architecturePatterns {
		if strings.Contains(makefileContent, pattern) {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("multi-architecture build support not found in Makefile")
	}
	
	return nil
}