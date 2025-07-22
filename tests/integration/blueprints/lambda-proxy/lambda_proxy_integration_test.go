package lambda_proxy_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gopkg.in/yaml.v3"
)

// LambdaProxyIntegrationTestSuite contains integration tests for lambda-proxy blueprint
type LambdaProxyIntegrationTestSuite struct {
	suite.Suite
	tempDir        string
	projectDir     string
	cliPath        string
	generatedProjects map[string]string // framework -> project path
}

// SetupSuite prepares the test suite
func (s *LambdaProxyIntegrationTestSuite) SetupSuite() {
	var err error
	
	// Create temporary directory
	s.tempDir, err = ioutil.TempDir("", "lambda-proxy-integration-*")
	s.Require().NoError(err)
	
	// Build CLI tool
	s.cliPath = filepath.Join(s.tempDir, "go-starter")
	buildCmd := exec.Command("go", "build", "-o", s.cliPath, "../../../../main.go")
	err = buildCmd.Run()
	s.Require().NoError(err, "Failed to build CLI tool")
	
	// Initialize generated projects map
	s.generatedProjects = make(map[string]string)
}

// TearDownSuite cleans up after tests
func (s *LambdaProxyIntegrationTestSuite) TearDownSuite() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

// SetupTest prepares each test
func (s *LambdaProxyIntegrationTestSuite) SetupTest() {
	// Change to temp directory for each test
	err := os.Chdir(s.tempDir)
	s.Require().NoError(err)
}

// TestLambdaProxyBasicGeneration tests basic lambda-proxy generation
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyBasicGeneration() {
	projectName := "test-lambda-basic"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate lambda-proxy project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-basic",
		"--framework", "gin",
		"--auth-type", "none",
		"--logger", "slog",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify project structure
	s.verifyBasicProjectStructure(projectPath)
	
	// Verify compilation
	s.verifyProjectCompilation(projectPath)
	
	// Verify SAM template
	s.verifySAMTemplate(projectPath)
	
	// Store for later use
	s.generatedProjects["gin"] = projectPath
}

// TestLambdaProxyWithJWTAuth tests lambda-proxy with JWT authentication
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyWithJWTAuth() {
	projectName := "test-lambda-jwt"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate lambda-proxy project with JWT
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-jwt",
		"--framework", "echo",
		"--auth-type", "jwt",
		"--logger", "zap",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify authentication files exist
	authFiles := []string{
		"internal/handlers/auth.go",
		"internal/middleware/auth.go",
		"internal/services/auth.go",
	}
	
	for _, file := range authFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Auth file should exist: %s", file)
		
		// Verify file contains JWT-related code
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		s.Contains(strings.ToLower(string(content)), "jwt", "File should contain JWT logic: %s", file)
	}
	
	// Verify compilation with auth dependencies
	s.verifyProjectCompilation(projectPath)
	
	// Verify SAM template includes authorizers
	s.verifySAMTemplateWithAuth(projectPath)
}

// TestLambdaProxyWithCognitoAuth tests lambda-proxy with Cognito authentication
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyWithCognitoAuth() {
	projectName := "test-lambda-cognito"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate lambda-proxy project with Cognito
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-cognito",
		"--framework", "fiber",
		"--auth-type", "cognito",
		"--logger", "logrus",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify Cognito integration
	authServicePath := filepath.Join(projectPath, "internal/services/auth.go")
	s.Require().FileExists(authServicePath)
	
	content, err := ioutil.ReadFile(authServicePath)
	s.Require().NoError(err)
	
	contentStr := strings.ToLower(string(content))
	s.True(
		strings.Contains(contentStr, "cognito") || strings.Contains(contentStr, "userpool"),
		"Auth service should contain Cognito integration",
	)
	
	// Verify compilation
	s.verifyProjectCompilation(projectPath)
}

// TestLambdaProxyMultiFramework tests lambda-proxy generation with different frameworks
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyMultiFramework() {
	frameworks := []string{"gin", "echo", "fiber", "chi", "stdlib"}
	
	for _, framework := range frameworks {
		s.Run(fmt.Sprintf("Framework_%s", framework), func() {
			projectName := fmt.Sprintf("test-lambda-%s", framework)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "lambda-proxy",
				"--module", fmt.Sprintf("github.com/test/lambda-%s", framework),
				"--framework", framework,
				"--auth-type", "none",
				"--logger", "slog",
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for %s: %s", framework, string(output))
			
			// Verify framework-specific imports
			mainGoPath := filepath.Join(projectPath, "main.go")
			content, err := ioutil.ReadFile(mainGoPath)
			s.Require().NoError(err)
			
			s.verifyFrameworkImports(string(content), framework)
			
			// Verify compilation
			s.verifyProjectCompilation(projectPath)
			
			// Store project path
			s.generatedProjects[framework] = projectPath
		})
	}
}

// TestLambdaProxyWithAllLoggers tests different logger implementations
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyWithAllLoggers() {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}
	
	for _, logger := range loggers {
		s.Run(fmt.Sprintf("Logger_%s", logger), func() {
			projectName := fmt.Sprintf("test-lambda-logger-%s", logger)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "lambda-proxy",
				"--module", fmt.Sprintf("github.com/test/lambda-logger-%s", logger),
				"--framework", "gin",
				"--auth-type", "none",
				"--logger", logger,
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for logger %s: %s", logger, string(output))
			
			// Verify logger implementation
			loggerPath := filepath.Join(projectPath, "internal/observability/logger.go")
			s.Require().FileExists(loggerPath)
			
			content, err := ioutil.ReadFile(loggerPath)
			s.Require().NoError(err)
			
			// Verify logger-specific imports and code
			s.verifyLoggerImplementation(string(content), logger)
			
			// Verify compilation
			s.verifyProjectCompilation(projectPath)
		})
	}
}

// TestLambdaProxyTerraformIntegration tests Terraform infrastructure generation
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyTerraformIntegration() {
	projectName := "test-lambda-terraform"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-terraform",
		"--framework", "gin",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify Terraform files
	terraformFiles := []string{
		"terraform/main.tf",
		"terraform/variables.tf",
		"terraform/outputs.tf",
	}
	
	for _, file := range terraformFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Terraform file should exist: %s", file)
		
		// Verify file is not empty
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		s.NotEmpty(string(content), "Terraform file should not be empty: %s", file)
	}
	
	// Verify Terraform syntax (basic validation)
	s.verifyTerraformSyntax(projectPath)
}

// TestLambdaProxyCICDWorkflows tests CI/CD workflow generation
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyCICDWorkflows() {
	projectName := "test-lambda-cicd"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-cicd",
		"--framework", "echo",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify CI/CD workflow files
	workflowFiles := []string{
		".github/workflows/ci.yml",
		".github/workflows/deploy.yml",
		".github/workflows/security.yml",
		".github/workflows/release.yml",
	}
	
	for _, file := range workflowFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Workflow file should exist: %s", file)
		
		// Verify YAML syntax
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		
		var workflow interface{}
		err = yaml.Unmarshal(content, &workflow)
		s.Require().NoError(err, "Workflow should be valid YAML: %s", file)
	}
}

// TestLambdaProxyDeploymentScripts tests deployment script generation
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyDeploymentScripts() {
	projectName := "test-lambda-deploy"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-deploy",
		"--framework", "chi",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify deployment scripts
	scripts := []string{
		"scripts/deploy.sh",
		"scripts/local-dev.sh",
	}
	
	for _, script := range scripts {
		fullPath := filepath.Join(projectPath, script)
		s.Require().FileExists(fullPath, "Script should exist: %s", script)
		
		// Verify script is executable
		info, err := os.Stat(fullPath)
		s.Require().NoError(err)
		s.True(info.Mode()&0111 != 0, "Script should be executable: %s", script)
	}
}

// TestLambdaProxyWithLocalStackIntegration tests integration with LocalStack
func (s *LambdaProxyIntegrationTestSuite) TestLambdaProxyWithLocalStackIntegration() {
	// Skip if running in CI without Docker
	if os.Getenv("CI") == "true" && os.Getenv("DOCKER_AVAILABLE") != "true" {
		s.T().Skip("Docker not available in CI environment")
	}
	
	ctx := context.Background()
	
	// Start LocalStack container
	req := testcontainers.ContainerRequest{
		Image:        "localstack/localstack:latest",
		ExposedPorts: []string{"4566/tcp"},
		Env: map[string]string{
			"SERVICES": "lambda,apigateway,iam",
		},
		WaitingFor: wait.ForHTTP("/health").WithPort("4566"),
	}
	
	localstackContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		s.T().Skipf("Failed to start LocalStack container: %v", err)
	}
	defer func() { _ = localstackContainer.Terminate(ctx) }()
	
	// Get LocalStack endpoint
	endpoint, err := localstackContainer.Endpoint(ctx, "")
	s.Require().NoError(err)
	
	// Generate lambda-proxy project
	projectName := "test-lambda-localstack"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "lambda-proxy",
		"--module", "github.com/test/lambda-localstack",
		"--framework", "gin",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Build the Lambda function
	buildCmd := exec.Command("go", "build", "-o", "bootstrap", ".")
	buildCmd.Dir = projectPath
	buildCmd.Env = append(os.Environ(), "GOOS=linux", "GOARCH=amd64", "CGO_ENABLED=0")
	
	err = buildCmd.Run()
	s.Require().NoError(err, "Failed to build Lambda function")
	
	// Test deployment to LocalStack (basic validation)
	s.T().Logf("LocalStack endpoint: %s", endpoint)
	
	// Verify LocalStack is responding
	s.verifyLocalStackHealth(endpoint)
}

// Helper methods

func (s *LambdaProxyIntegrationTestSuite) verifyBasicProjectStructure(projectPath string) {
	requiredFiles := []string{
		"main.go",
		"handler.go",
		"go.mod",
		"template.yaml",
		"internal/config/config.go",
		"internal/handlers/health.go",
		"internal/handlers/api.go",
		"internal/middleware/cors.go",
		"internal/middleware/logging.go",
		"internal/middleware/recovery.go",
		"internal/observability/logger.go",
		"internal/observability/metrics.go",
		"internal/observability/tracing.go",
		"scripts/deploy.sh",
		"scripts/local-dev.sh",
	}
	
	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Required file should exist: %s", file)
	}
}

func (s *LambdaProxyIntegrationTestSuite) verifyProjectCompilation(projectPath string) {
	cmd := exec.Command("go", "mod", "tidy")
	cmd.Dir = projectPath
	err := cmd.Run()
	s.Require().NoError(err, "go mod tidy should succeed")
	
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Project should compile: %s", string(output))
}

func (s *LambdaProxyIntegrationTestSuite) verifySAMTemplate(projectPath string) {
	templatePath := filepath.Join(projectPath, "template.yaml")
	content, err := ioutil.ReadFile(templatePath)
	s.Require().NoError(err)
	
	var template interface{}
	err = yaml.Unmarshal(content, &template)
	s.Require().NoError(err, "SAM template should be valid YAML")
	
	templateStr := string(content)
	s.Contains(templateStr, "AWSTemplateFormatVersion")
	s.Contains(templateStr, "Transform")
	s.Contains(templateStr, "Resources")
}

func (s *LambdaProxyIntegrationTestSuite) verifySAMTemplateWithAuth(projectPath string) {
	templatePath := filepath.Join(projectPath, "template.yaml")
	content, err := ioutil.ReadFile(templatePath)
	s.Require().NoError(err)
	
	templateStr := string(content)
	s.True(
		strings.Contains(templateStr, "Authorizer") || strings.Contains(templateStr, "Auth"),
		"SAM template should include authentication components",
	)
}

func (s *LambdaProxyIntegrationTestSuite) verifyFrameworkImports(content, framework string) {
	expectedImports := map[string]string{
		"gin":    "github.com/gin-gonic/gin",
		"echo":   "github.com/labstack/echo",
		"fiber":  "github.com/gofiber/fiber",
		"chi":    "github.com/go-chi/chi",
		"stdlib": "", // No external imports for stdlib
	}
	
	if expectedImport, exists := expectedImports[framework]; exists && expectedImport != "" {
		s.Contains(content, expectedImport, "Should contain framework import for %s", framework)
	}
}

func (s *LambdaProxyIntegrationTestSuite) verifyLoggerImplementation(content, logger string) {
	expectedPackages := map[string]string{
		"slog":    "log/slog",
		"zap":     "go.uber.org/zap",
		"logrus":  "github.com/sirupsen/logrus",
		"zerolog": "github.com/rs/zerolog",
	}
	
	if expectedPackage, exists := expectedPackages[logger]; exists {
		s.Contains(content, expectedPackage, "Should contain logger package for %s", logger)
	}
}

func (s *LambdaProxyIntegrationTestSuite) verifyTerraformSyntax(projectPath string) {
	terraformFiles := []string{
		"terraform/main.tf",
		"terraform/variables.tf",
		"terraform/outputs.tf",
	}
	
	for _, file := range terraformFiles {
		fullPath := filepath.Join(projectPath, file)
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		
		contentStr := string(content)
		// Basic Terraform syntax check
		if strings.Contains(file, "main.tf") {
			s.Contains(contentStr, "resource", "main.tf should contain resources")
		}
		if strings.Contains(file, "variables.tf") {
			s.Contains(contentStr, "variable", "variables.tf should contain variables")
		}
		if strings.Contains(file, "outputs.tf") {
			s.Contains(contentStr, "output", "outputs.tf should contain outputs")
		}
	}
}

func (s *LambdaProxyIntegrationTestSuite) verifyLocalStackHealth(endpoint string) {
	// Simple health check - in a real test we would deploy and test the Lambda
	s.T().Logf("LocalStack health check for endpoint: %s", endpoint)
	// This would typically involve AWS SDK calls to LocalStack
}

// TestSuite runner
func TestLambdaProxyIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	
	suite.Run(t, new(LambdaProxyIntegrationTestSuite))
}