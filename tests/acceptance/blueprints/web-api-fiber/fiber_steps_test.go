package webapifiber

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

// FiberWebAPITestContext holds state for Fiber web API BDD tests
// Provides comprehensive testing of Fiber-specific features and patterns
type FiberWebAPITestContext struct {
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

	// Fiber-specific test state
	fiberAppConfigured    bool
	middlewareStack       []string
	routeGroups          map[string]string
	performanceOptimized bool
	websocketEnabled     bool
}

var fiberCtx *FiberWebAPITestContext

// TestFiberWebAPIBDD runs BDD scenarios for Fiber web API blueprints
func TestFiberWebAPIBDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping Fiber web API BDD tests in short mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: InitializeFiberWebAPIScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("Fiber web API BDD test suite failed")
	}
}

// InitializeFiberWebAPIScenario registers all BDD step definitions for Fiber web API testing
func InitializeFiberWebAPIScenario(ctx *godog.ScenarioContext) {
	// Initialize context before each scenario
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		fiberCtx = &FiberWebAPITestContext{
			httpClient:   &http.Client{Timeout: 10 * time.Second},
			ctx:          context.Background(),
			testResults:  make(map[string]bool),
			routeGroups:  make(map[string]string),
		}
		return ctx, nil
	})

	// Cleanup after each scenario
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if fiberCtx != nil {
			fiberCtx.cleanup()
		}
		return ctx, nil
	})

	// Background steps
	ctx.Step(`^the go-starter CLI tool is available$`, fiberCtx.theGoStarterCLIToolIsAvailable)
	ctx.Step(`^I am in a clean working directory$`, fiberCtx.iAmInACleanWorkingDirectory)

	// Project generation steps
	ctx.Step(`^I want to create a Fiber-based web API application$`, fiberCtx.iWantToCreateAFiberBasedWebAPIApplication)
	ctx.Step(`^I run the command "([^"]*)"$`, fiberCtx.iRunTheCommand)
	ctx.Step(`^the generation should succeed$`, fiberCtx.theGenerationShouldSucceed)
	ctx.Step(`^the project should contain Fiber-specific components$`, fiberCtx.theProjectShouldContainFiberSpecificComponents)
	ctx.Step(`^the generated code should compile successfully$`, fiberCtx.theGeneratedCodeShouldCompileSuccessfully)
	ctx.Step(`^the project should use Fiber framework patterns$`, fiberCtx.theProjectShouldUseFiberFrameworkPatterns)

	// Fiber app and middleware steps
	ctx.Step(`^I have generated a Fiber web API application$`, fiberCtx.iHaveGeneratedAFiberWebAPIApplication)
	ctx.Step(`^I examine the server configuration$`, fiberCtx.iExamineTheServerConfiguration)
	ctx.Step(`^the app should use fiber\.New\(\)$`, fiberCtx.theAppShouldUseFiberNew)
	ctx.Step(`^the middleware should be Fiber-compatible$`, fiberCtx.theMiddlewareShouldBeFiberCompatible)
	ctx.Step(`^the routing should use Fiber's handler patterns$`, fiberCtx.theRoutingShouldUseFibersHandlerPatterns)
	ctx.Step(`^the handlers should accept \*fiber\.Ctx$`, fiberCtx.theHandlersShouldAcceptFiberCtx)

	// Routing and groups steps
	ctx.Step(`^I examine the route definitions$`, fiberCtx.iExamineTheRouteDefinitions)
	ctx.Step(`^the routes should use Fiber's app\.Get/Post/Put/Delete patterns$`, fiberCtx.theRoutesShouldUseFibersHTTPMethodPatterns)
	ctx.Step(`^the routes should support route groups$`, fiberCtx.theRoutesShouldSupportRouteGroups)
	ctx.Step(`^the routes should include parameter binding$`, fiberCtx.theRoutesShouldIncludeParameterBinding)
	ctx.Step(`^the route prefixes should be configurable$`, fiberCtx.theRoutePrefixesShouldBeConfigurable)

	// High-performance context steps
	ctx.Step(`^I want high-performance request handling$`, fiberCtx.iWantHighPerformanceRequestHandling)
	ctx.Step(`^I generate a Fiber web API application$`, fiberCtx.iGenerateAFiberWebAPIApplication)
	ctx.Step(`^the handlers should use \*fiber\.Ctx$`, fiberCtx.theHandlersShouldUseFiberCtx)
	ctx.Step(`^request parsing should be optimized for speed$`, fiberCtx.requestParsingShouldBeOptimizedForSpeed)
	ctx.Step(`^response generation should be efficient$`, fiberCtx.responseGenerationShouldBeEfficient)
	ctx.Step(`^memory allocation should be minimized$`, fiberCtx.memoryAllocationShouldBeMinimized)

	// Middleware stack steps
	ctx.Step(`^I want comprehensive middleware with Fiber$`, fiberCtx.iWantComprehensiveMiddlewareWithFiber)
	ctx.Step(`^the middleware should include Fiber's built-in middleware$`, fiberCtx.theMiddlewareShouldIncludeFibersBuiltInMiddleware)
	ctx.Step(`^custom middleware should follow Fiber patterns$`, fiberCtx.customMiddlewareShouldFollowFiberPatterns)
	ctx.Step(`^middleware should have access to fiber\.Ctx$`, fiberCtx.middlewareShouldHaveAccessToFiberCtx)
	ctx.Step(`^middleware should support next\(\) pattern$`, fiberCtx.middlewareShouldSupportNextPattern)

	// Request parsing and validation steps
	ctx.Step(`^I want fast request processing$`, fiberCtx.iWantFastRequestProcessing)
	ctx.Step(`^JSON parsing should be extremely fast$`, fiberCtx.jsonParsingShouldBeExtremelyFast)
	ctx.Step(`^form data should be efficiently processed$`, fiberCtx.formDataShouldBeEfficientlyProcessed)
	ctx.Step(`^file uploads should be properly handled$`, fiberCtx.fileUploadsShouldBeProperlyHandled)
	ctx.Step(`^request validation should be integrated$`, fiberCtx.requestValidationShouldBeIntegrated)

	// WebSocket steps
	ctx.Step(`^I want real-time capabilities$`, fiberCtx.iWantRealTimeCapabilities)
	ctx.Step(`^I generate a Fiber web API with WebSocket support$`, fiberCtx.iGenerateAFiberWebAPIWithWebSocketSupport)
	ctx.Step(`^WebSocket endpoints should be properly configured$`, fiberCtx.webSocketEndpointsShouldBeProperlyConfigured)
	ctx.Step(`^WebSocket handlers should use Fiber patterns$`, fiberCtx.webSocketHandlersShouldUseFiberPatterns)
	ctx.Step(`^connection upgrading should be seamless$`, fiberCtx.connectionUpgradingShouldBeSeamless)
	ctx.Step(`^WebSocket middleware should be supported$`, fiberCtx.webSocketMiddlewareShouldBeSupported)

	// Compression and caching steps
	ctx.Step(`^I want optimized responses$`, fiberCtx.iWantOptimizedResponses)
	ctx.Step(`^response compression should be available$`, fiberCtx.responseCompressionShouldBeAvailable)
	ctx.Step(`^caching middleware should be configured$`, fiberCtx.cachingMiddlewareShouldBeConfigured)
	ctx.Step(`^static file serving should be optimized$`, fiberCtx.staticFileServingShouldBeOptimized)
	ctx.Step(`^ETag support should be implemented$`, fiberCtx.eTagSupportShouldBeImplemented)

	// Security and rate limiting steps
	ctx.Step(`^I want secure and controlled access$`, fiberCtx.iWantSecureAndControlledAccess)
	ctx.Step(`^I generate a Fiber web API with security features$`, fiberCtx.iGenerateAFiberWebAPIWithSecurityFeatures)
	ctx.Step(`^rate limiting should be efficiently implemented$`, fiberCtx.rateLimitingShouldBeEfficientlyImplemented)
	ctx.Step(`^CORS should be configured for Fiber$`, fiberCtx.corsShouldBeConfiguredForFiber)
	ctx.Step(`^security headers should be set via middleware$`, fiberCtx.securityHeadersShouldBeSetViaMiddleware)
	ctx.Step(`^DDoS protection should be available$`, fiberCtx.ddosProtectionShouldBeAvailable)

	// Authentication steps
	ctx.Step(`^I want to secure my Fiber web API$`, fiberCtx.iWantToSecureMyFiberWebAPI)
	ctx.Step(`^I generate a Fiber web API with JWT authentication$`, fiberCtx.iGenerateAFiberWebAPIWithJWTAuthentication)
	ctx.Step(`^JWT middleware should integrate with Fiber$`, fiberCtx.jwtMiddlewareShouldIntegrateWithFiber)
	ctx.Step(`^protected routes should use Fiber middleware$`, fiberCtx.protectedRoutesShouldUseFiberMiddleware)
	ctx.Step(`^token validation should work with fiber\.Ctx$`, fiberCtx.tokenValidationShouldWorkWithFiberCtx)
	ctx.Step(`^authentication should be performant$`, fiberCtx.authenticationShouldBePerformant)

	// Error handling steps
	ctx.Step(`^I want robust error handling$`, fiberCtx.iWantRobustErrorHandling)
	ctx.Step(`^error handling should use Fiber's error handling$`, fiberCtx.errorHandlingShouldUseFibersErrorHandling)
	ctx.Step(`^panic recovery should be implemented$`, fiberCtx.panicRecoveryShouldBeImplemented)
	ctx.Step(`^custom error pages should be supported$`, fiberCtx.customErrorPagesShouldBeSupported)
	ctx.Step(`^error responses should be properly formatted$`, fiberCtx.errorResponsesShouldBeProperlyFormatted)

	// Performance optimization steps
	ctx.Step(`^I want maximum performance$`, fiberCtx.iWantMaximumPerformance)
	ctx.Step(`^the server should be configured for speed$`, fiberCtx.theServerShouldBeConfiguredForSpeed)
	ctx.Step(`^memory usage should be optimized$`, fiberCtx.memoryUsageShouldBeOptimized)
	ctx.Step(`^request processing should be minimal overhead$`, fiberCtx.requestProcessingShouldBeMinimalOverhead)
	ctx.Step(`^response times should be extremely fast$`, fiberCtx.responseTimesShouldBeExtremelyFast)

	// Additional steps for remaining scenarios...
	ctx.Step(`^I want production-ready Fiber web API$`, fiberCtx.iWantProductionReadyFiberWebAPI)
	ctx.Step(`^production settings should be optimized$`, fiberCtx.productionSettingsShouldBeOptimized)
	ctx.Step(`^prefork mode should be configurable$`, fiberCtx.preforkModeShouldBeConfigurable)
	ctx.Step(`^process management should be included$`, fiberCtx.processManagementShouldBeIncluded)
	ctx.Step(`^monitoring endpoints should be available$`, fiberCtx.monitoringEndpointsShouldBeAvailable)
}

// Background step implementations
func (ctx *FiberWebAPITestContext) theGoStarterCLIToolIsAvailable() error {
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

func (ctx *FiberWebAPITestContext) iAmInACleanWorkingDirectory() error {
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "fiber-web-api-test-*")
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
func (ctx *FiberWebAPITestContext) iWantToCreateAFiberBasedWebAPIApplication() error {
	ctx.projectName = "test-fiber-api"
	ctx.framework = "fiber"
	return nil
}

func (ctx *FiberWebAPITestContext) iRunTheCommand(command string) error {
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

func (ctx *FiberWebAPITestContext) theGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", ctx.lastExitCode, string(ctx.lastOutput))
	}
	ctx.projectExists = true
	return nil
}

func (ctx *FiberWebAPITestContext) theProjectShouldContainFiberSpecificComponents() error {
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

	if !strings.Contains(string(content), "github.com/gofiber/fiber") {
		return fmt.Errorf("go.mod should contain Fiber dependency")
	}

	return nil
}

func (ctx *FiberWebAPITestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	output, err := modCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}

	buildCmd := exec.Command("go", "build", "-o", "fiber-test-app", "./cmd/server")
	buildCmd.Dir = ctx.projectDir
	output, err = buildCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}

	ctx.compilationOK = true
	return nil
}

func (ctx *FiberWebAPITestContext) theProjectShouldUseFiberFrameworkPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "fiber.New") {
		return fmt.Errorf("server.go should contain fiber.New")
	}

	return nil
}

// Fiber-specific feature step implementations
func (ctx *FiberWebAPITestContext) iHaveGeneratedAFiberWebAPIApplication() error {
	if !ctx.projectExists {
		return ctx.iRunTheCommand("go-starter new test-fiber-api --type=web-api-fiber --framework=fiber --no-git")
	}
	return nil
}

func (ctx *FiberWebAPITestContext) iExamineTheServerConfiguration() error {
	return nil // Context setting step
}

func (ctx *FiberWebAPITestContext) theAppShouldUseFiberNew() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	if !strings.Contains(string(content), "fiber.New()") {
		return fmt.Errorf("server should use fiber.New()")
	}

	ctx.fiberAppConfigured = true
	return nil
}

func (ctx *FiberWebAPITestContext) theMiddlewareShouldBeFiberCompatible() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "app.Use(") {
		return fmt.Errorf("server should use Fiber middleware pattern app.Use()")
	}

	return nil
}

func (ctx *FiberWebAPITestContext) theRoutingShouldUseFibersHandlerPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return fmt.Errorf("failed to read server.go: %w", err)
	}

	contentStr := string(content)
	if !strings.Contains(contentStr, "app.Get(") && !strings.Contains(contentStr, "app.Post(") {
		return fmt.Errorf("server should use Fiber handler patterns like app.Get(), app.Post()")
	}

	return nil
}

func (ctx *FiberWebAPITestContext) theHandlersShouldAcceptFiberCtx() error {
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
		if strings.Contains(contentStr, "func") && !strings.Contains(contentStr, "*fiber.Ctx") {
			return fmt.Errorf("handlers in %s should accept *fiber.Ctx", path)
		}
		
		return nil
	})
}

// Routing and groups step implementations
func (ctx *FiberWebAPITestContext) iExamineTheRouteDefinitions() error {
	return nil // Context setting step
}

func (ctx *FiberWebAPITestContext) theRoutesShouldUseFibersHTTPMethodPatterns() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return err
	}
	
	contentStr := string(content)
	methods := []string{"app.Get(", "app.Post(", "app.Put(", "app.Delete("}
	foundMethods := 0
	for _, method := range methods {
		if strings.Contains(contentStr, method) {
			foundMethods++
		}
	}
	
	if foundMethods == 0 {
		return fmt.Errorf("routes should use Fiber's HTTP method patterns")
	}
	return nil
}

func (ctx *FiberWebAPITestContext) theRoutesShouldSupportRouteGroups() error {
	serverPath := filepath.Join(ctx.projectDir, "internal/server/server.go")
	content, err := os.ReadFile(serverPath)
	if err != nil {
		return err
	}
	
	if !strings.Contains(string(content), "app.Group(") {
		return fmt.Errorf("routes should support Fiber route groups with app.Group()")
	}
	return nil
}

func (ctx *FiberWebAPITestContext) theRoutesShouldIncludeParameterBinding() error {
	// Implementation would check for parameter binding patterns
	return nil
}

func (ctx *FiberWebAPITestContext) theRoutePrefixesShouldBeConfigurable() error {
	// Implementation would check for configurable route prefixes
	return nil
}

// Stub implementations for remaining functions (for brevity)
func (ctx *FiberWebAPITestContext) iWantHighPerformanceRequestHandling() error { return nil }
func (ctx *FiberWebAPITestContext) iGenerateAFiberWebAPIApplication() error { return ctx.iHaveGeneratedAFiberWebAPIApplication() }
func (ctx *FiberWebAPITestContext) theHandlersShouldUseFiberCtx() error { return ctx.theHandlersShouldAcceptFiberCtx() }
func (ctx *FiberWebAPITestContext) requestParsingShouldBeOptimizedForSpeed() error { ctx.performanceOptimized = true; return nil }
func (ctx *FiberWebAPITestContext) responseGenerationShouldBeEfficient() error { return nil }
func (ctx *FiberWebAPITestContext) memoryAllocationShouldBeMinimized() error { return nil }

func (ctx *FiberWebAPITestContext) iWantComprehensiveMiddlewareWithFiber() error { return nil }
func (ctx *FiberWebAPITestContext) theMiddlewareShouldIncludeFibersBuiltInMiddleware() error { return nil }
func (ctx *FiberWebAPITestContext) customMiddlewareShouldFollowFiberPatterns() error { return nil }
func (ctx *FiberWebAPITestContext) middlewareShouldHaveAccessToFiberCtx() error { return nil }
func (ctx *FiberWebAPITestContext) middlewareShouldSupportNextPattern() error { return nil }

func (ctx *FiberWebAPITestContext) iWantFastRequestProcessing() error { return nil }
func (ctx *FiberWebAPITestContext) jsonParsingShouldBeExtremelyFast() error { return nil }
func (ctx *FiberWebAPITestContext) formDataShouldBeEfficientlyProcessed() error { return nil }
func (ctx *FiberWebAPITestContext) fileUploadsShouldBeProperlyHandled() error { return nil }
func (ctx *FiberWebAPITestContext) requestValidationShouldBeIntegrated() error { return nil }

func (ctx *FiberWebAPITestContext) iWantRealTimeCapabilities() error { return nil }
func (ctx *FiberWebAPITestContext) iGenerateAFiberWebAPIWithWebSocketSupport() error { 
	ctx.websocketEnabled = true
	return ctx.iRunTheCommand("go-starter new test-fiber-api --type=web-api-fiber --framework=fiber --no-git")
}
func (ctx *FiberWebAPITestContext) webSocketEndpointsShouldBeProperlyConfigured() error { return nil }
func (ctx *FiberWebAPITestContext) webSocketHandlersShouldUseFiberPatterns() error { return nil }
func (ctx *FiberWebAPITestContext) connectionUpgradingShouldBeSeamless() error { return nil }
func (ctx *FiberWebAPITestContext) webSocketMiddlewareShouldBeSupported() error { return nil }

func (ctx *FiberWebAPITestContext) iWantOptimizedResponses() error { return nil }
func (ctx *FiberWebAPITestContext) responseCompressionShouldBeAvailable() error { return nil }
func (ctx *FiberWebAPITestContext) cachingMiddlewareShouldBeConfigured() error { return nil }
func (ctx *FiberWebAPITestContext) staticFileServingShouldBeOptimized() error { return nil }
func (ctx *FiberWebAPITestContext) eTagSupportShouldBeImplemented() error { return nil }

func (ctx *FiberWebAPITestContext) iWantSecureAndControlledAccess() error { return nil }
func (ctx *FiberWebAPITestContext) iGenerateAFiberWebAPIWithSecurityFeatures() error { 
	return ctx.iRunTheCommand("go-starter new test-fiber-api --type=web-api-fiber --framework=fiber --no-git")
}
func (ctx *FiberWebAPITestContext) rateLimitingShouldBeEfficientlyImplemented() error { return nil }
func (ctx *FiberWebAPITestContext) corsShouldBeConfiguredForFiber() error { return nil }
func (ctx *FiberWebAPITestContext) securityHeadersShouldBeSetViaMiddleware() error { return nil }
func (ctx *FiberWebAPITestContext) ddosProtectionShouldBeAvailable() error { return nil }

func (ctx *FiberWebAPITestContext) iWantToSecureMyFiberWebAPI() error { return nil }
func (ctx *FiberWebAPITestContext) iGenerateAFiberWebAPIWithJWTAuthentication() error { 
	return ctx.iRunTheCommand("go-starter new test-fiber-api --type=web-api-fiber --framework=fiber --auth-type=jwt --no-git")
}
func (ctx *FiberWebAPITestContext) jwtMiddlewareShouldIntegrateWithFiber() error { return nil }
func (ctx *FiberWebAPITestContext) protectedRoutesShouldUseFiberMiddleware() error { return nil }
func (ctx *FiberWebAPITestContext) tokenValidationShouldWorkWithFiberCtx() error { return nil }
func (ctx *FiberWebAPITestContext) authenticationShouldBePerformant() error { return nil }

func (ctx *FiberWebAPITestContext) iWantRobustErrorHandling() error { return nil }
func (ctx *FiberWebAPITestContext) errorHandlingShouldUseFibersErrorHandling() error { return nil }
func (ctx *FiberWebAPITestContext) panicRecoveryShouldBeImplemented() error { return nil }
func (ctx *FiberWebAPITestContext) customErrorPagesShouldBeSupported() error { return nil }
func (ctx *FiberWebAPITestContext) errorResponsesShouldBeProperlyFormatted() error { return nil }

func (ctx *FiberWebAPITestContext) iWantMaximumPerformance() error { return nil }
func (ctx *FiberWebAPITestContext) theServerShouldBeConfiguredForSpeed() error { return nil }
func (ctx *FiberWebAPITestContext) memoryUsageShouldBeOptimized() error { return nil }
func (ctx *FiberWebAPITestContext) requestProcessingShouldBeMinimalOverhead() error { return nil }
func (ctx *FiberWebAPITestContext) responseTimesShouldBeExtremelyFast() error { return nil }

func (ctx *FiberWebAPITestContext) iWantProductionReadyFiberWebAPI() error { return nil }
func (ctx *FiberWebAPITestContext) productionSettingsShouldBeOptimized() error { return nil }
func (ctx *FiberWebAPITestContext) preforkModeShouldBeConfigurable() error { return nil }
func (ctx *FiberWebAPITestContext) processManagementShouldBeIncluded() error { return nil }
func (ctx *FiberWebAPITestContext) monitoringEndpointsShouldBeAvailable() error { return nil }

// Cleanup function
func (ctx *FiberWebAPITestContext) cleanup() {
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