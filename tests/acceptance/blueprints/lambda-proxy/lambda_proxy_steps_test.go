package lambda_proxy_test

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
	"github.com/cucumber/godog/colors"
	messages "github.com/cucumber/messages/go/v21"
	"gopkg.in/yaml.v3"
)

// TestContext holds the test execution context
type TestContext struct {
	t               *testing.T
	tempDir         string
	projectDir      string
	projectName     string
	generationError error
	parameters      map[string]string
	frameworks      []string
	currentFramework string
}

// NewTestContext creates a new test context
func NewTestContext(t *testing.T) *TestContext {
	return &TestContext{
		t:          t,
		parameters: make(map[string]string),
		frameworks: []string{"gin", "echo", "fiber", "chi", "stdlib"},
	}
}

// Cleanup removes temporary directories
func (tc *TestContext) Cleanup() {
	if tc.tempDir != "" {
		_ = os.RemoveAll(tc.tempDir)
	}
}

// Lambda Proxy ATDD Step Definitions

func (tc *TestContext) iHaveTheGoStarterCLIToolAvailable() error {
	// Build the CLI tool if needed
	cliPath := filepath.Join(tc.tempDir, "go-starter")
	buildCmd := exec.Command("go", "build", "-o", cliPath, "../../../main.go")
	buildCmd.Dir = "../../../../"
	
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build CLI tool: %w\nOutput: %s", err, output)
	}
	
	// Verify the tool is executable
	if _, err := os.Stat(cliPath); err != nil {
		return fmt.Errorf("CLI tool not found at %s: %w", cliPath, err)
	}
	
	return nil
}

func (tc *TestContext) iAmInATemporaryWorkingDirectory() error {
	var err error
	tc.tempDir, err = os.MkdirTemp("", "lambda-proxy-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	if err := os.Chdir(tc.tempDir); err != nil {
		return fmt.Errorf("failed to change to temp directory: %w", err)
	}
	
	return nil
}

func (tc *TestContext) iWantToGenerateALambdaProxyBlueprint() error {
	tc.parameters["type"] = "lambda-proxy"
	return nil
}

func (tc *TestContext) iRunTheGeneratorWith(table *godog.Table) error {
	// Parse parameters from table
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			tc.parameters[row.Cells[0].Value] = row.Cells[1].Value
		}
	}
	
	// Set project name for directory creation
	if name, exists := tc.parameters["name"]; exists {
		tc.projectName = name
		tc.projectDir = filepath.Join(tc.tempDir, name)
	}
	
	// Build CLI command
	cliPath := filepath.Join(tc.tempDir, "go-starter")
	args := []string{"new"}
	
	// Add parameters as arguments
	for key, value := range tc.parameters {
		switch key {
		case "type":
			args = append(args, "--type", value)
		case "name":
			args = append(args, value) // Project name is positional
		case "module":
			args = append(args, "--module", value)
		case "framework":
			args = append(args, "--framework", value)
		case "auth_type":
			args = append(args, "--auth-type", value)
		case "logger_type":
			args = append(args, "--logger", value)
		}
	}
	
	// Run the generator
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, cliPath, args...)
	cmd.Dir = tc.tempDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		tc.generationError = fmt.Errorf("generation failed: %w\nOutput: %s", err, output)
		return nil // Don't return error here, let subsequent steps check
	}
	
	return nil
}

func (tc *TestContext) theGenerationShouldSucceed() error {
	if tc.generationError != nil {
		return tc.generationError
	}
	
	// Verify project directory was created
	if tc.projectDir == "" {
		return fmt.Errorf("project directory not set")
	}
	
	if _, err := os.Stat(tc.projectDir); err != nil {
		return fmt.Errorf("project directory not created: %w", err)
	}
	
	return nil
}

func (tc *TestContext) theGeneratedProjectShouldHaveTheFollowingStructure(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			filePath := row.Cells[0].Value
			fileType := row.Cells[1].Value
			
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			switch fileType {
			case "file":
				if _, err := os.Stat(fullPath); err != nil {
					return fmt.Errorf("expected file %s not found: %w", filePath, err)
				}
			case "directory":
				if info, err := os.Stat(fullPath); err != nil {
					return fmt.Errorf("expected directory %s not found: %w", filePath, err)
				} else if !info.IsDir() {
					return fmt.Errorf("%s is not a directory", filePath)
				}
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theProjectShouldCompileSuccessfully() error {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = tc.projectDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (tc *TestContext) theMainGoShouldContainGinFrameworkImports() error {
	return tc.verifyFrameworkImports("gin", "github.com/gin-gonic/gin")
}

func (tc *TestContext) verifyFrameworkImports(framework, expectedImport string) error {
	mainGoPath := filepath.Join(tc.projectDir, "main.go")
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go: %w", err)
	}
	
	if framework != "stdlib" && !strings.Contains(string(content), expectedImport) {
		return fmt.Errorf("expected import %s not found in main.go", expectedImport)
	}
	
	return nil
}

func (tc *TestContext) theSAMTemplateShouldBeValid() error {
	templatePath := filepath.Join(tc.projectDir, "template.yaml")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read SAM template: %w", err)
	}
	
	var samTemplate interface{}
	if err := yaml.Unmarshal(content, &samTemplate); err != nil {
		return fmt.Errorf("SAM template is not valid YAML: %w", err)
	}
	
	// Basic SAM template validation
	templateStr := string(content)
	requiredSections := []string{"AWSTemplateFormatVersion", "Transform", "Resources"}
	for _, section := range requiredSections {
		if !strings.Contains(templateStr, section) {
			return fmt.Errorf("SAM template missing required section: %s", section)
		}
	}
	
	return nil
}

func (tc *TestContext) theGeneratedProjectShouldIncludeAuthenticationFiles(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			filePath := row.Cells[0].Value
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			if _, err := os.Stat(fullPath); err != nil {
				return fmt.Errorf("expected authentication file %s not found: %w", filePath, err)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theAuthMiddlewareShouldContainJWTValidationLogic() error {
	authPath := filepath.Join(tc.projectDir, "internal/middleware/auth.go")
	content, err := os.ReadFile(authPath)
	if err != nil {
		return fmt.Errorf("failed to read auth middleware: %w", err)
	}
	
	contentStr := string(content)
	jwtKeywords := []string{"jwt", "token", "authorization", "Bearer"}
	found := false
	for _, keyword := range jwtKeywords {
		if strings.Contains(strings.ToLower(contentStr), strings.ToLower(keyword)) {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("auth middleware does not contain JWT validation logic")
	}
	
	return nil
}

func (tc *TestContext) theSAMTemplateShouldIncludeCustomAuthorizers() error {
	templatePath := filepath.Join(tc.projectDir, "template.yaml")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read SAM template: %w", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "Authorizer") && !strings.Contains(contentStr, "Auth") {
		return fmt.Errorf("SAM template does not include custom authorizers")
	}
	
	return nil
}

func (tc *TestContext) theGeneratedProjectShouldIncludeCognitoIntegration(table *godog.Table) error {
	return tc.theGeneratedProjectShouldIncludeAuthenticationFiles(table)
}

func (tc *TestContext) theAuthServiceShouldContainCognitoUserPoolValidation() error {
	authPath := filepath.Join(tc.projectDir, "internal/services/auth.go")
	content, err := os.ReadFile(authPath)
	if err != nil {
		return fmt.Errorf("failed to read auth service: %w", err)
	}
	
	contentStr := string(content)
	cognitoKeywords := []string{"cognito", "user pool", "UserPool"}
	found := false
	for _, keyword := range cognitoKeywords {
		if strings.Contains(strings.ToLower(contentStr), strings.ToLower(keyword)) {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("auth service does not contain Cognito user pool validation")
	}
	
	return nil
}

func (tc *TestContext) theSAMTemplateShouldReferenceCognitoUserPools() error {
	templatePath := filepath.Join(tc.projectDir, "template.yaml")
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read SAM template: %w", err)
	}
	
	contentStr := string(content)
	if !strings.Contains(contentStr, "UserPool") && !strings.Contains(contentStr, "Cognito") {
		return fmt.Errorf("SAM template does not reference Cognito user pools")
	}
	
	return nil
}

func (tc *TestContext) iRunTheGeneratorForEachFramework(table *godog.Table) error {
	frameworks := []string{}
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			frameworks = append(frameworks, row.Cells[0].Value)
		}
	}
	
	for _, framework := range frameworks {
		tc.currentFramework = framework
		tc.parameters["framework"] = framework
		tc.parameters["name"] = fmt.Sprintf("test-lambda-%s", framework)
		
		if err := tc.iRunTheGeneratorWith(&godog.Table{
			Rows: []*messages.PickleTableRow{
				{Cells: []*messages.PickleTableCell{{Value: "framework"}, {Value: framework}}},
			},
		}); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for framework %s: %w", framework, err)
		}
	}
	
	return nil
}

func (tc *TestContext) eachGeneratedProjectShould(table *godog.Table) error {
	for _, framework := range tc.frameworks {
		projectDir := filepath.Join(tc.tempDir, fmt.Sprintf("test-lambda-%s", framework))
		
		for _, row := range table.Rows {
			if len(row.Cells) >= 1 {
				validation := row.Cells[0].Value
				
				switch validation {
				case "compile successfully":
					cmd := exec.Command("go", "build", "./...")
					cmd.Dir = projectDir
					if output, err := cmd.CombinedOutput(); err != nil {
						return fmt.Errorf("compilation failed for %s: %w\nOutput: %s", framework, err, output)
					}
					
				case "contain framework-specific imports":
					if err := tc.verifyFrameworkSpecificImports(projectDir, framework); err != nil {
						return err
					}
					
				case "have consistent API structure":
					if err := tc.verifyConsistentAPIStructure(projectDir); err != nil {
						return err
					}
					
				case "include proper middleware integration":
					if err := tc.verifyMiddlewareIntegration(projectDir); err != nil {
						return err
					}
					
				case "generate valid SAM templates":
					templatePath := filepath.Join(projectDir, "template.yaml")
					if err := tc.validateSAMTemplate(templatePath); err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

func (tc *TestContext) verifyFrameworkSpecificImports(projectDir, framework string) error {
	mainGoPath := filepath.Join(projectDir, "main.go")
	content, err := os.ReadFile(mainGoPath)
	if err != nil {
		return fmt.Errorf("failed to read main.go for %s: %w", framework, err)
	}
	
	expectedImports := map[string]string{
		"gin":    "github.com/gin-gonic/gin",
		"echo":   "github.com/labstack/echo",
		"fiber":  "github.com/gofiber/fiber",
		"chi":    "github.com/go-chi/chi",
		"stdlib": "", // stdlib doesn't need external imports
	}
	
	if expectedImport, exists := expectedImports[framework]; exists && expectedImport != "" {
		if !strings.Contains(string(content), expectedImport) {
			return fmt.Errorf("framework %s missing expected import %s", framework, expectedImport)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyConsistentAPIStructure(projectDir string) error {
	requiredFiles := []string{
		"internal/handlers/health.go",
		"internal/handlers/api.go",
		"internal/middleware/cors.go",
		"internal/middleware/logging.go",
		"internal/middleware/recovery.go",
	}
	
	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectDir, file)
		if _, err := os.Stat(fullPath); err != nil {
			return fmt.Errorf("consistent API structure file missing: %s", file)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyMiddlewareIntegration(projectDir string) error {
	handlerPath := filepath.Join(projectDir, "handler.go")
	if _, err := os.Stat(handlerPath); err != nil {
		// Try alternative handler file location
		handlerPath = filepath.Join(projectDir, "internal/handlers/api.go")
	}
	
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		return fmt.Errorf("failed to read handler file: %w", err)
	}
	
	middlewareKeywords := []string{"middleware", "CORS", "logging", "recovery"}
	for _, keyword := range middlewareKeywords {
		if !strings.Contains(strings.ToLower(string(content)), strings.ToLower(keyword)) {
			return fmt.Errorf("middleware integration missing keyword: %s", keyword)
		}
	}
	
	return nil
}

func (tc *TestContext) validateSAMTemplate(templatePath string) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read SAM template: %w", err)
	}
	
	var samTemplate interface{}
	if err := yaml.Unmarshal(content, &samTemplate); err != nil {
		return fmt.Errorf("SAM template is not valid YAML: %w", err)
	}
	
	return nil
}

// Additional step definitions for remaining scenarios...

func (tc *TestContext) theGeneratedProjectShouldIncludeTerraformFiles(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			filePath := row.Cells[0].Value
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			if _, err := os.Stat(fullPath); err != nil {
				return fmt.Errorf("expected Terraform file %s not found: %w", filePath, err)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theTerraformConfigurationShouldBeValid() error {
	terraformFiles := []string{
		"terraform/main.tf",
		"terraform/variables.tf",
		"terraform/outputs.tf",
	}
	
	for _, file := range terraformFiles {
		fullPath := filepath.Join(tc.projectDir, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read Terraform file %s: %w", file, err)
		}
		
		// Basic Terraform syntax validation
		contentStr := string(content)
		if !strings.Contains(contentStr, "resource") && !strings.Contains(contentStr, "variable") && !strings.Contains(contentStr, "output") {
			return fmt.Errorf("Terraform file %s appears to be empty or invalid", file)
		}
	}
	
	return nil
}

// Test runner function
func TestLambdaProxyBlueprintATDD(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			tc := NewTestContext(t)
			
			// Register step definitions
			ctx.Step(`^I have the go-starter CLI tool available$`, tc.iHaveTheGoStarterCLIToolAvailable)
			ctx.Step(`^I am in a temporary working directory$`, tc.iAmInATemporaryWorkingDirectory)
			ctx.Step(`^I want to generate a lambda-proxy blueprint$`, tc.iWantToGenerateALambdaProxyBlueprint)
			ctx.Step(`^I run the generator with:$`, tc.iRunTheGeneratorWith)
			ctx.Step(`^the generation should succeed$`, tc.theGenerationShouldSucceed)
			ctx.Step(`^the generated project should have the following structure:$`, tc.theGeneratedProjectShouldHaveTheFollowingStructure)
			ctx.Step(`^the project should compile successfully$`, tc.theProjectShouldCompileSuccessfully)
			ctx.Step(`^the main\.go should contain Gin framework imports$`, tc.theMainGoShouldContainGinFrameworkImports)
			ctx.Step(`^the SAM template should be valid$`, tc.theSAMTemplateShouldBeValid)
			ctx.Step(`^the generated project should include authentication files:$`, tc.theGeneratedProjectShouldIncludeAuthenticationFiles)
			ctx.Step(`^the auth middleware should contain JWT validation logic$`, tc.theAuthMiddlewareShouldContainJWTValidationLogic)
			ctx.Step(`^the SAM template should include custom authorizers$`, tc.theSAMTemplateShouldIncludeCustomAuthorizers)
			ctx.Step(`^the generated project should include Cognito integration:$`, tc.theGeneratedProjectShouldIncludeCognitoIntegration)
			ctx.Step(`^the auth service should contain Cognito user pool validation$`, tc.theAuthServiceShouldContainCognitoUserPoolValidation)
			ctx.Step(`^the SAM template should reference Cognito user pools$`, tc.theSAMTemplateShouldReferenceCognitoUserPools)
			ctx.Step(`^I run the generator for each framework:$`, tc.iRunTheGeneratorForEachFramework)
			ctx.Step(`^each generated project should:$`, tc.eachGeneratedProjectShould)
			ctx.Step(`^the generated project should include Terraform files:$`, tc.theGeneratedProjectShouldIncludeTerraformFiles)
			ctx.Step(`^the Terraform configuration should be valid$`, tc.theTerraformConfigurationShouldBeValid)
			
			// Cleanup after each scenario
			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				tc.Cleanup()
				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"lambda-proxy.feature"},
			TestingT: t,
			Output:   colors.Colored(os.Stdout),
		},
	}
	
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run lambda-proxy ATDD tests")
	}
}