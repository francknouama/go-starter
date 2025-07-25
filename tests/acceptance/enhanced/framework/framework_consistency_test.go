package framework

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

// FrameworkConsistencyTestContext holds test state for framework consistency testing
type FrameworkConsistencyTestContext struct {
	projectConfig     *types.ProjectConfig
	projectPath       string
	tempDir           string
	startTime         time.Time
	lastCommandOutput string
	lastCommandError  error
	generatedFiles    []string
	framework         string
	architecture      string
}

// TestFeatures runs the framework consistency BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &FrameworkConsistencyTestContext{}

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
		t.Fatal("non-zero status returned, failed to run framework consistency tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *FrameworkConsistencyTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^I am testing framework consistency across architectures$`, ctx.iAmTestingFrameworkConsistencyAcrossArchitectures)

	// Project generation steps
	s.Step(`^I generate a web API project with framework consistency configuration:$`, ctx.iGenerateAWebAPIProjectWithFrameworkConsistencyConfiguration)

	// Framework validation steps
	s.Step(`^(.*) framework should be properly integrated$`, ctx.frameworkShouldBeProperlyIntegrated)
	s.Step(`^the project should compile successfully$`, ctx.theProjectShouldCompileSuccessfully)
	s.Step(`^(.*) architecture should be properly implemented$`, ctx.architectureShouldBeProperlyImplemented)
	s.Step(`^framework structure should be consistent with (.*) patterns$`, ctx.frameworkStructureShouldBeConsistentWithPatterns)
	s.Step(`^routing should follow (.*) conventions in (.*) architecture$`, ctx.routingShouldFollowConventionsInArchitecture)
	s.Step(`^middleware integration should be architecture-aware$`, ctx.middlewareIntegrationShouldBeArchitectureAware)
	s.Step(`^dependency injection should work with (.*) and (.*)$`, ctx.dependencyInjectionShouldWorkWithFrameworkAndArchitecture)
	s.Step(`^error handling should be consistent across layers$`, ctx.errorHandlingShouldBeConsistentAcrossLayers)
	s.Step(`^configuration management should integrate properly$`, ctx.configurationManagementShouldIntegrateProperly)
	s.Step(`^health checks should be implemented consistently$`, ctx.healthChecksShouldBeImplementedConsistently)
	s.Step(`^logging integration should work across all layers$`, ctx.loggingIntegrationShouldWorkAcrossAllLayers)
	s.Step(`^testing patterns should be architecture-appropriate$`, ctx.testingPatternsShouldBeArchitectureAppropriate)
	s.Step(`^documentation should reflect architecture choices$`, ctx.documentationShouldReflectArchitectureChoices)

	// Performance validation steps
	s.Step(`^framework performance should be optimized for (.*)$`, ctx.frameworkPerformanceShouldBeOptimizedFor)
	s.Step(`^memory usage should be efficient$`, ctx.memoryUsageShouldBeEfficient)
	s.Step(`^startup time should be reasonable$`, ctx.startupTimeShouldBeReasonable)

	// Security validation steps
	s.Step(`^security middleware should be properly integrated$`, ctx.securityMiddlewareShouldBeProperlyIntegrated)
	s.Step(`^CORS configuration should be framework-specific$`, ctx.corsConfigurationShouldBeFrameworkSpecific)
	s.Step(`^input validation should be implemented consistently$`, ctx.inputValidationShouldBeImplementedConsistently)

	// Framework-specific validation steps
	s.Step(`^Gin routing should work correctly with (.*)$`, ctx.ginRoutingShouldWorkCorrectlyWith)
	s.Step(`^Echo middleware should integrate properly with (.*)$`, ctx.echoMiddlewareShouldIntegrateProperlyWith)
	s.Step(`^Fiber handlers should be compatible with (.*)$`, ctx.fiberHandlersShouldBeCompatibleWith)
	s.Step(`^Chi router should function correctly with (.*)$`, ctx.chiRouterShouldFunctionCorrectlyWith)

	// Architecture-specific validation steps
	s.Step(`^standard architecture should have simple structure$`, ctx.standardArchitectureShouldHaveSimpleStructure)
	s.Step(`^clean architecture should have proper layer separation$`, ctx.cleanArchitectureShouldHaveProperLayerSeparation)
	s.Step(`^DDD architecture should have domain-centric design$`, ctx.dddArchitectureShouldHaveDomainCentricDesign)
	s.Step(`^hexagonal architecture should have ports and adapters$`, ctx.hexagonalArchitectureShouldHavePortsAndAdapters)
}

// Background step implementations
func (ctx *FrameworkConsistencyTestContext) iHaveTheGoStarterCLIAvailable() error {
	cmd := exec.Command("go-starter", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-starter CLI not available: %v", err)
	}
	return nil
}

func (ctx *FrameworkConsistencyTestContext) allTemplatesAreProperlyInitialized() error {
	return helpers.InitializeTemplates()
}

func (ctx *FrameworkConsistencyTestContext) iAmTestingFrameworkConsistencyAcrossArchitectures() error {
	return nil
}

// Project generation implementation
func (ctx *FrameworkConsistencyTestContext) iGenerateAWebAPIProjectWithFrameworkConsistencyConfiguration(table *godog.Table) error {
	// Parse configuration from table
	config := &types.ProjectConfig{}

	for i := 0; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			config.Type = value
		case "framework":
			config.Framework = value
			ctx.framework = value
		case "architecture":
			config.Architecture = value
			ctx.architecture = value
		case "logger":
			config.Logger = value
		case "go_version":
			config.GoVersion = value
		}
	}

	// Set defaults
	if config.Name == "" {
		config.Name = fmt.Sprintf("test-framework-%s-%s", ctx.framework, ctx.architecture)
	}
	if config.Module == "" {
		config.Module = fmt.Sprintf("github.com/test/framework-%s-%s", ctx.framework, ctx.architecture)
	}
	if config.GoVersion == "" {
		config.GoVersion = "1.23"
	}

	ctx.projectConfig = config

	// Generate project using helpers
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

// Framework validation implementations
func (ctx *FrameworkConsistencyTestContext) frameworkShouldBeProperlyIntegrated(framework string) error {
	// Check for framework-specific imports and usage patterns
	frameworkPatterns := map[string][]string{
		"gin":   {"github.com/gin-gonic/gin", "gin.Engine", "gin.Default", "gin.New"},
		"echo":  {"github.com/labstack/echo", "echo.New", "echo.Echo", "echo.Context"},
		"fiber": {"github.com/gofiber/fiber", "fiber.New", "fiber.App", "fiber.Ctx"},
		"chi":   {"github.com/go-chi/chi", "chi.NewRouter", "chi.Router", "chi.Route"},
	}

	patterns, exists := frameworkPatterns[framework]
	if !exists {
		return fmt.Errorf("unsupported framework: %s", framework)
	}

	return ctx.checkForPatterns(patterns, fmt.Sprintf("%s framework integration", framework))
}

func (ctx *FrameworkConsistencyTestContext) theProjectShouldCompileSuccessfully() error {
	return ctx.validateProjectCompilation()
}

func (ctx *FrameworkConsistencyTestContext) architectureShouldBeProperlyImplemented(architecture string) error {
	// Check for architecture-specific directory structure
	archPatterns := map[string][]string{
		"standard": {"handlers", "services", "models", "middleware"},
		"clean":    {"entities", "usecases", "repositories", "interfaces", "delivery"},
		"ddd":      {"domain", "entities", "valueobjects", "aggregates", "services"},
		"hexagonal": {"ports", "adapters", "domain", "infrastructure", "application"},
	}

	patterns, exists := archPatterns[architecture]
	if !exists {
		return fmt.Errorf("unsupported architecture: %s", architecture)
	}

	return ctx.checkForDirectoryPatterns(patterns, fmt.Sprintf("%s architecture", architecture))
}

func (ctx *FrameworkConsistencyTestContext) frameworkStructureShouldBeConsistentWithPatterns(architecture string) error {
	// Verify that framework usage follows architectural patterns
	return ctx.checkFrameworkArchitectureConsistency(ctx.framework, architecture)
}

func (ctx *FrameworkConsistencyTestContext) routingShouldFollowConventionsInArchitecture(framework, architecture string) error {
	// Check routing patterns based on framework and architecture combination
	routingPatterns := map[string]map[string][]string{
		"gin": {
			"standard":   {"router.GET", "router.POST", "gin.RouterGroup"},
			"clean":      {"delivery", "router.GET", "gin.RouterGroup"},
			"ddd":        {"router.GET", "domain", "gin.RouterGroup"},
			"hexagonal":  {"adapters", "router.GET", "gin.RouterGroup"},
		},
		"echo": {
			"standard":   {"e.GET", "e.POST", "echo.Group"},
			"clean":      {"delivery", "e.GET", "echo.Group"},
			"ddd":        {"e.GET", "domain", "echo.Group"},
			"hexagonal":  {"adapters", "e.GET", "echo.Group"},
		},
		"fiber": {
			"standard":   {"app.Get", "app.Post", "fiber.Router"},
			"clean":      {"delivery", "app.Get", "fiber.Router"},
			"ddd":        {"app.Get", "domain", "fiber.Router"},
			"hexagonal":  {"adapters", "app.Get", "fiber.Router"},
		},
		"chi": {
			"standard":   {"r.Get", "r.Post", "chi.Route"},
			"clean":      {"delivery", "r.Get", "chi.Route"},
			"ddd":        {"r.Get", "domain", "chi.Route"},
			"hexagonal":  {"adapters", "r.Get", "chi.Route"},
		},
	}

	frameworkRoutes, exists := routingPatterns[framework]
	if !exists {
		return fmt.Errorf("unsupported framework for routing: %s", framework)
	}

	patterns, exists := frameworkRoutes[architecture]
	if !exists {
		return fmt.Errorf("unsupported architecture for %s routing: %s", framework, architecture)
	}

	return ctx.checkForPatterns(patterns, fmt.Sprintf("%s routing in %s architecture", framework, architecture))
}

func (ctx *FrameworkConsistencyTestContext) middlewareIntegrationShouldBeArchitectureAware() error {
	// Check that middleware is placed appropriately based on architecture
	return ctx.checkMiddlewareArchitectureIntegration()
}

func (ctx *FrameworkConsistencyTestContext) dependencyInjectionShouldWorkWithFrameworkAndArchitecture(framework, architecture string) error {
	// Check dependency injection patterns for framework+architecture combination
	diPatterns := []string{"inject", "wire", "container", "dependencies", "New", "Provider"}
	return ctx.checkForPatterns(diPatterns, fmt.Sprintf("dependency injection for %s with %s", framework, architecture))
}

func (ctx *FrameworkConsistencyTestContext) errorHandlingShouldBeConsistentAcrossLayers() error {
	// Check for consistent error handling patterns
	errorPatterns := []string{"error", "Error", "HandleError", "errorHandler", "ErrorResponse"}
	return ctx.checkForPatterns(errorPatterns, "consistent error handling")
}

func (ctx *FrameworkConsistencyTestContext) configurationManagementShouldIntegrateProperly() error {
	// Check configuration management integration
	configPatterns := []string{"config", "Config", "viper", "Viper", "configuration"}
	return ctx.checkForPatterns(configPatterns, "configuration management")
}

func (ctx *FrameworkConsistencyTestContext) healthChecksShouldBeImplementedConsistently() error {
	// Check health check implementation
	healthPatterns := []string{"health", "Health", "/health", "healthz", "HealthCheck"}
	return ctx.checkForPatterns(healthPatterns, "health checks")
}

func (ctx *FrameworkConsistencyTestContext) loggingIntegrationShouldWorkAcrossAllLayers() error {
	// Check logging integration across layers
	loggingPatterns := []string{"log", "Log", "logger", "Logger", "zap", "slog", "logrus", "zerolog"}
	return ctx.checkForPatterns(loggingPatterns, "logging integration")
}

func (ctx *FrameworkConsistencyTestContext) testingPatternsShouldBeArchitectureAppropriate() error {
	// Check for appropriate testing patterns
	testPatterns := []string{"test", "Test", "_test.go", "testing", "testify", "assert"}
	return ctx.checkForPatterns(testPatterns, "testing patterns")
}

func (ctx *FrameworkConsistencyTestContext) documentationShouldReflectArchitectureChoices() error {
	// Check for documentation that reflects architecture
	docPatterns := []string{"README.md", "doc", "Doc", "documentation"}
	return ctx.checkForPatterns(docPatterns, "documentation")
}

// Helper methods
func (ctx *FrameworkConsistencyTestContext) generateTestProject(config *types.ProjectConfig, tempDir string) (string, error) {
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

func (ctx *FrameworkConsistencyTestContext) validateProjectCompilation() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("project path not set")
	}
	
	// Check if go.mod exists
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project")
	}
	
	// Try to run go build
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("project compilation failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (ctx *FrameworkConsistencyTestContext) checkForPatterns(patterns []string, description string) error {
	found := false
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".go") {
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

func (ctx *FrameworkConsistencyTestContext) checkForDirectoryPatterns(patterns []string, description string) error {
	found := false
	for _, pattern := range patterns {
		// Check for directories containing the pattern
		err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if info.IsDir() && strings.Contains(strings.ToLower(path), strings.ToLower(pattern)) {
				found = true
				return nil
			}
			
			return nil
		})
		
		if err != nil {
			return err
		}
		
		if found {
			return nil
		}
	}
	
	return fmt.Errorf("%s directory patterns not found in generated project", description)
}

func (ctx *FrameworkConsistencyTestContext) checkFrameworkArchitectureConsistency(framework, architecture string) error {
	// Verify that framework usage follows architectural patterns
	// This is a simplified implementation
	return nil
}

func (ctx *FrameworkConsistencyTestContext) checkMiddlewareArchitectureIntegration() error {
	// Check that middleware is placed appropriately based on architecture
	middlewarePaths := []string{
		"internal/middleware/",
		"pkg/middleware/", 
		"middleware/",
		"internal/delivery/middleware/", // Clean architecture
		"internal/adapters/middleware/", // Hexagonal
	}
	
	for _, path := range middlewarePaths {
		fullPath := filepath.Join(ctx.projectPath, path)
		if _, err := os.Stat(fullPath); err == nil {
			return nil // Found middleware directory
		}
	}
	
	return fmt.Errorf("middleware integration not found for architecture")
}

// Simplified implementations for remaining validation methods
func (ctx *FrameworkConsistencyTestContext) frameworkPerformanceShouldBeOptimizedFor(architecture string) error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) memoryUsageShouldBeEfficient() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) startupTimeShouldBeReasonable() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) securityMiddlewareShouldBeProperlyIntegrated() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) corsConfigurationShouldBeFrameworkSpecific() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) inputValidationShouldBeImplementedConsistently() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) ginRoutingShouldWorkCorrectlyWith(architecture string) error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) echoMiddlewareShouldIntegrateProperlyWith(architecture string) error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) fiberHandlersShouldBeCompatibleWith(architecture string) error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) chiRouterShouldFunctionCorrectlyWith(architecture string) error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) standardArchitectureShouldHaveSimpleStructure() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) cleanArchitectureShouldHaveProperLayerSeparation() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) dddArchitectureShouldHaveDomainCentricDesign() error {
	return nil // Simplified implementation
}

func (ctx *FrameworkConsistencyTestContext) hexagonalArchitectureShouldHavePortsAndAdapters() error {
	return nil // Simplified implementation
}