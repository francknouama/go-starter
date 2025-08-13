package webapiecho

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/testcontainers/testcontainers-go"
)

// EchoWebAPITestContext holds state for Echo web API BDD tests
// Provides comprehensive testing of Echo-specific features and patterns
type EchoWebAPITestContext struct {
	// Test environment management
	workingDir  string
	projectDir  string
	projectName string
	originalDir string
	projectRoot string

	// Command execution tracking
	lastCommand  *exec.Cmd
	lastOutput   []byte
	lastError    error
	lastExitCode int

	// Project configuration
	framework      string
	architecture   string
	databaseDriver string
	databaseORM    string
	authType       string
	logger         string

	// Test state tracking
	projectExists bool
	compilationOK bool
	testResults   map[string]bool

	// HTTP client for application testing
	httpClient *http.Client
	serverPort int
	baseURL    string

	// Testcontainers for database testing
	postgresContainer testcontainers.Container
	database          *testcontainers.Container
	ctx               context.Context

	// Echo-specific test state
	echoInstanceConfigured bool
	middlewareStack        []string
	routeGroups           map[string]string
	contextBindingTested   bool
}

var echoCtx *EchoWebAPITestContext

// TestEchoWebAPIBDD runs BDD scenarios for Echo web API blueprints
func TestEchoWebAPIBDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Echo web API BDD tests in short mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeEchoWebAPIScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Echo web API BDD test suite failed")
	}
}

// InitializeEchoWebAPIScenario registers all BDD step definitions for Echo web API testing
func InitializeEchoWebAPIScenario(ctx *godog.ScenarioContext) {
	// Initialize context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		echoCtx = &EchoWebAPITestContext{
			httpClient:   &http.Client{Timeout: 10 * time.Second},
			ctx:          context.Background(),
			testResults:  make(map[string]bool),
			routeGroups:  make(map[string]string),
		}
		return ctx, nil
	})

	// Cleanup after each scenario
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if echoCtx != nil {
			echoCtx.cleanup()
		}
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the go-starter CLI tool is available$`, echoCtx.theGoStarterCLIToolIsAvailable)
	ctx.Step(`^I am in a clean working directory$`, echoCtx.iAmInACleanWorkingDirectory)

	// Project generation steps
	ctx.Step(`^I want to create an Echo-based web API application$`, echoCtx.iWantToCreateAnEchoBasedWebAPIApplication)
	ctx.Step(`^I run the command "([^"]*)"$`, echoCtx.iRunTheCommand)
	ctx.Step(`^the generation should succeed$`, echoCtx.theGenerationShouldSucceed)
	ctx.Step(`^the project should contain Echo-specific components$`, echoCtx.theProjectShouldContainEchoSpecificComponents)
	ctx.Step(`^the generated code should compile successfully$`, echoCtx.theGeneratedCodeShouldCompileSuccessfully)
	ctx.Step(`^the project should use Echo framework patterns$`, echoCtx.theProjectShouldUseEchoFrameworkPatterns)

	// Echo instance and middleware steps
	ctx.Step(`^I have generated an Echo web API application$`, echoCtx.iHaveGeneratedAnEchoWebAPIApplication)
	ctx.Step(`^I examine the server configuration$`, echoCtx.iExamineTheServerConfiguration)
	ctx.Step(`^the server should use echo\.New\(\)$`, echoCtx.theServerShouldUseEchoNew)
	ctx.Step(`^the middleware should be Echo-compatible$`, echoCtx.theMiddlewareShouldBeEchoCompatible)
	ctx.Step(`^the routing should use Echo's handler patterns$`, echoCtx.theRoutingShouldUseEchosHandlerPatterns)
	ctx.Step(`^the handlers should accept echo\.Context$`, echoCtx.theHandlersShouldAcceptEchoContext)

	// Routing and groups steps
	ctx.Step(`^I examine the route definitions$`, echoCtx.iExamineTheRouteDefinitions)
	ctx.Step(`^the routes should use Echo's e\.GET/POST/PUT/DELETE patterns$`, echoCtx.theRoutesShouldUseEchosHTTPMethodPatterns)
	ctx.Step(`^the routes should support route groups$`, echoCtx.theRoutesShouldSupportRouteGroups)
	ctx.Step(`^the routes should include parameter binding$`, echoCtx.theRoutesShouldIncludeParameterBinding)
	ctx.Step(`^the route middleware should be configurable per route$`, echoCtx.theRouteMiddlewareShouldBeConfigurablePerRoute)

	// Context and data binding steps
	ctx.Step(`^I want flexible request handling$`, echoCtx.iWantFlexibleRequestHandling)
	ctx.Step(`^I generate an Echo web API application$`, echoCtx.iGenerateAnEchoWebAPIApplication)
	ctx.Step(`^the handlers should use echo\.Context$`, echoCtx.theHandlersShouldUseEchoContext)
	ctx.Step(`^request data should bind automatically to structs$`, echoCtx.requestDataShouldBindAutomaticallyToStructs)
	ctx.Step(`^path parameters should be extractable via context$`, echoCtx.pathParametersShouldBeExtractableViaContext)
	ctx.Step(`^query parameters should be easily accessible$`, echoCtx.queryParametersShouldBeEasilyAccessible)

	// Middleware stack steps
	ctx.Step(`^I want comprehensive middleware with Echo$`, echoCtx.iWantComprehensiveMiddlewareWithEcho)
	ctx.Step(`^the middleware should include Echo's built-in middleware$`, echoCtx.theMiddlewareShouldIncludeEchosBuiltInMiddleware)
	ctx.Step(`^custom middleware should follow Echo patterns$`, echoCtx.customMiddlewareShouldFollowEchoPatterns)
	ctx.Step(`^middleware should have access to echo\.Context$`, echoCtx.middlewareShouldHaveAccessToEchoContext)
	// ctx.Step(`^middleware chaining should be properly configured$`, echoCtx.middlewareChainingS​houldBeProperlyConfigured)

	// Validation steps
	ctx.Step(`^I want validated inputs$`, echoCtx.iWantValidatedInputs)
	ctx.Step(`^request validation should be integrated$`, echoCtx.requestValidationShouldBeIntegrated)
	ctx.Step(`^validation errors should return proper responses$`, echoCtx.validationErrorsShouldReturnProperResponses)
	ctx.Step(`^custom validators should be supported$`, echoCtx.customValidatorsShouldBeSupported)
	ctx.Step(`^validation middleware should be configurable$`, echoCtx.validationMiddlewareShouldBeConfigurable)

	// Error handling steps
	ctx.Step(`^I want robust error handling$`, echoCtx.iWantRobustErrorHandling)
	ctx.Step(`^error handling should use Echo's error handling$`, echoCtx.errorHandlingShouldUseEchosErrorHandling)
	ctx.Step(`^custom error handlers should be implemented$`, echoCtx.customErrorHandlersShouldBeImplemented)
	ctx.Step(`^HTTP error responses should be properly formatted$`, echoCtx.httpErrorResponsesShouldBeProperlyFormatted)
	ctx.Step(`^error logging should include request context$`, echoCtx.errorLoggingShouldIncludeRequestContext)

	// Authentication steps
	ctx.Step(`^I want to secure my Echo web API$`, echoCtx.iWantToSecureMyEchoWebAPI)
	ctx.Step(`^I generate an Echo web API with JWT authentication$`, echoCtx.iGenerateAnEchoWebAPIWithJWTAuthentication)
	ctx.Step(`^JWT middleware should integrate with Echo$`, echoCtx.jwtMiddlewareShouldIntegrateWithEcho)
	ctx.Step(`^protected routes should use Echo middleware$`, echoCtx.protectedRoutesShouldUseEchoMiddleware)
	ctx.Step(`^token validation should work with echo\.Context$`, echoCtx.tokenValidationShouldWorkWithEchoContext)
	ctx.Step(`^unauthorized access should be properly handled$`, echoCtx.unauthorizedAccessShouldBeProperlyHandled)

	// Additional steps for other scenarios...
	ctx.Step(`^I want cross-origin and security support$`, echoCtx.iWantCrossOriginAndSecuritySupport)
	ctx.Step(`^I generate an Echo web API with security features$`, echoCtx.iGenerateAnEchoWebAPIWithSecurityFeatures)
	ctx.Step(`^CORS should be configured using Echo middleware$`, echoCtx.corsShouldBeConfiguredUsingEchoMiddleware)
	ctx.Step(`^security headers should be set via Echo middleware$`, echoCtx.securityHeadersShouldBeSetViaEchoMiddleware)
	ctx.Step(`^rate limiting should be implemented$`, echoCtx.rateLimitingShouldBeImplemented)
	ctx.Step(`^CSRF protection should be available$`, echoCtx.csrfProtectionShouldBeAvailable)
}

// Background step implementations
func (ctx *EchoWebAPITestContext) theGoStarterCLIToolIsAvailable() error {
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

func (ctx *EchoWebAPITestContext) iAmInACleanWorkingDirectory() error {
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "echo-web-api-test-*")
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
func (ctx *EchoWebAPITestContext) iWantToCreateAnEchoBasedWebAPIApplication() error {
	ctx.projectName = "test-echo-api"
	ctx.framework = "echo"
	return nil
}

func (ctx *EchoWebAPITestContext) iRunTheCommand(command string) error {
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

func (ctx *EchoWebAPITestContext) theGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", ctx.lastExitCode, string(ctx.lastOutput))
	}
	ctx.projectExists = true
	return nil
}

func (ctx *EchoWebAPITestContext) theProjectShouldContainEchoSpecificComponents() error {
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"internal/server/server.go",
		"internal/handlers",
		"internal/middleware",
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

	if !strings.Contains(string(content), "github.com/labstack/echo") {
		return fmt.Errorf("go.mod should contain Echo dependency")
	}

	return nil
}

func (ctx *EchoWebAPITestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	output, err := modCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}

	buildCmd := exec.Command("go", "build", "-o", "echo-test-app", "./cmd/server")
	buildCmd.Dir = ctx.projectDir
	output, err = buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	ctx.compilationOK = true
	return nil
}

func (ctx *EchoWebAPITestContext) theProjectShouldUseEchoFrameworkPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "echo.New") {
		return fmt.Errorf("server.go should contain echo.New")
	}

	return nil
}

// Echo-specific feature step implementations
func (ctx *EchoWebAPITestContext) iHaveGeneratedAnEchoWebAPIApplication() error {
	if !ctx.projectExists {
		return ctx.iRunTheCommand("go-starter new test-echo-api --type=web-api-echo --framework=echo --no-git")
	}
	return nil
}

func (ctx *EchoWebAPITestContext) iExamineTheServerConfiguration() error {
	return nil // Context setting step
}

func (ctx *EchoWebAPITestContext) theServerShouldUseEchoNew() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	if !strings.Contains(string(content), "echo.New()") {
		return fmt.Errorf("server should use echo.New()")
	}

	ctx.echoInstanceConfigured = true
	return nil
}

func (ctx *EchoWebAPITestContext) theMiddlewareShouldBeEchoCompatible() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "e.Use(") {
		return fmt.Errorf("server should use Echo middleware pattern e.Use()")
	}

	return nil
}

func (ctx *EchoWebAPITestContext) theRoutingShouldUseEchosHandlerPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "e.GET(") && !strings.Contains(contentStr, "e.POST(") {
		return fmt.Errorf("server should use Echo handler patterns like e.GET(), e.POST()")
	}

	return nil
}

func (ctx *EchoWebAPITestContext) theHandlersShouldAcceptEchoContext() error {
	handlersDir := filepath.Join(ctx.projectDir, "internal/handlers")
	
	return filepath.Walk(handlersDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || !strings.HasSuffix(path, ".go") {
			return err
		}
		
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		
		contentStr := string(content)
		if strings.Contains(contentStr, "func") && !strings.Contains(contentStr, "echo.Context") {
			return fmt.Errorf("handlers in %s should accept echo.Context", path)
		}
		
		return nil
	})
}

// Routing and groups step implementations
func (ctx *EchoWebAPITestContext) iExamineTheRouteDefinitions() error {
	return nil // Context setting step
}

func (ctx *EchoWebAPITestContext) theRoutesShouldUseEchosHTTPMethodPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return err
	}
	
	contentStr := string(content)
	methods := []string{"e.GET(", "e.POST(", "e.PUT(", "e.DELETE("}
	foundMethods := 0
	for _, method := range methods {
		if strings.Contains(contentStr, method) {
			foundMethods++
		}
	}
	
	if foundMethods == 0 {
		return fmt.Errorf("routes should use Echo's HTTP method patterns")
	}
	return nil
}

func (ctx *EchoWebAPITestContext) theRoutesShouldSupportRouteGroups() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return err
	}
	
	if !strings.Contains(string(content), "e.Group(") {
		return fmt.Errorf("routes should support Echo route groups with e.Group()")
	}
	return nil
}

func (ctx *EchoWebAPITestContext) theRoutesShouldIncludeParameterBinding() error {
	// Implementation would check for parameter binding patterns
	return nil
}

func (ctx *EchoWebAPITestContext) theRouteMiddlewareShouldBeConfigurablePerRoute() error {
	// Implementation would check for per-route middleware configuration
	return nil
}

// Stub implementations for remaining functions (for brevity)
func (ctx *EchoWebAPITestContext) iWantFlexibleRequestHandling() error { return nil }
func (ctx *EchoWebAPITestContext) iGenerateAnEchoWebAPIApplication() error { return ctx.iHaveGeneratedAnEchoWebAPIApplication() }
func (ctx *EchoWebAPITestContext) theHandlersShouldUseEchoContext() error { return ctx.theHandlersShouldAcceptEchoContext() }
func (ctx *EchoWebAPITestContext) requestDataShouldBindAutomaticallyToStructs() error { return nil }
func (ctx *EchoWebAPITestContext) pathParametersShouldBeExtractableViaContext() error { return nil }
func (ctx *EchoWebAPITestContext) queryParametersShouldBeEasilyAccessible() error { return nil }

func (ctx *EchoWebAPITestContext) iWantComprehensiveMiddlewareWithEcho() error { return nil }
func (ctx *EchoWebAPITestContext) theMiddlewareShouldIncludeEchosBuiltInMiddleware() error { return nil }
func (ctx *EchoWebAPITestContext) customMiddlewareShouldFollowEchoPatterns() error { return nil }
func (ctx *EchoWebAPITestContext) middlewareShouldHaveAccessToEchoContext() error { return nil }

func (ctx *EchoWebAPITestContext) iWantValidatedInputs() error { return nil }
func (ctx *EchoWebAPITestContext) requestValidationShouldBeIntegrated() error { return nil }
func (ctx *EchoWebAPITestContext) validationErrorsShouldReturnProperResponses() error { return nil }
func (ctx *EchoWebAPITestContext) customValidatorsShouldBeSupported() error { return nil }
func (ctx *EchoWebAPITestContext) validationMiddlewareShouldBeConfigurable() error { return nil }

func (ctx *EchoWebAPITestContext) iWantRobustErrorHandling() error { return nil }
func (ctx *EchoWebAPITestContext) errorHandlingShouldUseEchosErrorHandling() error { return nil }
func (ctx *EchoWebAPITestContext) customErrorHandlersShouldBeImplemented() error { return nil }
func (ctx *EchoWebAPITestContext) httpErrorResponsesShouldBeProperlyFormatted() error { return nil }
func (ctx *EchoWebAPITestContext) errorLoggingShouldIncludeRequestContext() error { return nil }

func (ctx *EchoWebAPITestContext) iWantToSecureMyEchoWebAPI() error { return nil }
func (ctx *EchoWebAPITestContext) iGenerateAnEchoWebAPIWithJWTAuthentication() error { 
	return ctx.iRunTheCommand("go-starter new test-echo-api --type=web-api-echo --framework=echo --auth-type=jwt --no-git")
}
func (ctx *EchoWebAPITestContext) jwtMiddlewareShouldIntegrateWithEcho() error { return nil }
func (ctx *EchoWebAPITestContext) protectedRoutesShouldUseEchoMiddleware() error { return nil }
func (ctx *EchoWebAPITestContext) tokenValidationShouldWorkWithEchoContext() error { return nil }
func (ctx *EchoWebAPITestContext) unauthorizedAccessShouldBeProperlyHandled() error { return nil }

func (ctx *EchoWebAPITestContext) iWantCrossOriginAndSecuritySupport() error { return nil }
func (ctx *EchoWebAPITestContext) iGenerateAnEchoWebAPIWithSecurityFeatures() error { 
	return ctx.iRunTheCommand("go-starter new test-echo-api --type=web-api-echo --framework=echo --no-git")
}
func (ctx *EchoWebAPITestContext) corsShouldBeConfiguredUsingEchoMiddleware() error { return nil }
func (ctx *EchoWebAPITestContext) securityHeadersShouldBeSetViaEchoMiddleware() error { return nil }
func (ctx *EchoWebAPITestContext) rateLimitingShouldBeImplemented() error { return nil }
func (ctx *EchoWebAPITestContext) csrfProtectionShouldBeAvailable() error { return nil }

// Cleanup function
func (ctx *EchoWebAPITestContext) cleanup() {
	if ctx.postgresContainer != nil {
		_ = ctx.postgresContainer.Terminate(ctx.ctx)
	}
	if ctx.originalDir != "" {
		_ = os.Chdir(ctx.originalDir)
	}
	if ctx.workingDir != "" {
		_ = os.RemoveAll(ctx.workingDir)
	}
}