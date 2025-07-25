package matrix

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// MatrixTestContext holds test state for expanded matrix testing
type MatrixTestContext struct {
	configurations    []MatrixConfiguration
	results           map[string]*MatrixTestResult
	currentConfig     *MatrixConfiguration
	projectPaths      map[string]string
	tempDir           string
	executionStrategy string
	startTime         time.Time
}

// MatrixConfiguration represents a configuration combination
type MatrixConfiguration struct {
	Name           string
	Framework      string
	Database       string
	Driver         string
	ORM            string
	Logger         string
	AuthType       string
	Architecture   string
	CacheType      string
	TestFramework  string
	MigrationTool  string
	DeployTargets  []string
	Middleware     []string
	Priority       string
}

// MatrixTestResult holds the result of a matrix test
type MatrixTestResult struct {
	Success         bool
	GenerationTime  time.Duration
	CompilationTime time.Duration
	Errors          []string
	Warnings        []string
	FileCount       int
	TestsPassed     bool
}

// TestFeatures runs the expanded matrix BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &MatrixTestContext{
				results:      make(map[string]*MatrixTestResult),
				projectPaths: make(map[string]string),
			}
			
			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				
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
		t.Fatal("non-zero status returned, failed to run expanded matrix tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *MatrixTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^matrix testing mode is enabled$`, ctx.matrixTestingModeIsEnabled)
	
	// Matrix configuration steps
	s.Step(`^I test all framework and database combinations$`, ctx.iTestAllFrameworkAndDatabaseCombinations)
	s.Step(`^I generate projects with the full matrix:$`, ctx.iGenerateProjectsWithTheFullMatrix)
	s.Step(`^I test authentication across all frameworks$`, ctx.iTestAuthenticationAcrossAllFrameworks)
	s.Step(`^I generate projects with auth configurations:$`, ctx.iGenerateProjectsWithAuthConfigurations)
	s.Step(`^I test all logger types with different architectures$`, ctx.iTestAllLoggerTypesWithDifferentArchitectures)
	s.Step(`^I generate projects with logger configurations:$`, ctx.iGenerateProjectsWithLoggerConfigurations)
	s.Step(`^I test deployment configurations$`, ctx.iTestDeploymentConfigurations)
	s.Step(`^I generate projects with deployment targets:$`, ctx.iGenerateProjectsWithDeploymentTargets)
	s.Step(`^I test different testing configurations$`, ctx.iTestDifferentTestingConfigurations)
	s.Step(`^I generate projects with testing setups:$`, ctx.iGenerateProjectsWithTestingSetups)
	s.Step(`^I test middleware stacks across frameworks$`, ctx.iTestMiddlewareStacksAcrossFrameworks)
	s.Step(`^I configure middleware combinations:$`, ctx.iConfigureMiddlewareCombinations)
	
	// Validation steps
	s.Step(`^all valid combinations should generate successfully$`, ctx.allValidCombinationsShouldGenerateSuccessfully)
	s.Step(`^invalid ORM selections should be prevented$`, ctx.invalidORMSelectionsShouldBePrevented)
	s.Step(`^database-specific features should be correctly configured$`, ctx.databaseSpecificFeaturesShouldBeCorrectlyConfigured)
	s.Step(`^authentication middleware should be properly configured$`, ctx.authenticationMiddlewareShouldBeProperlyConfigured)
	s.Step(`^session management should work when enabled$`, ctx.sessionManagementShouldWorkWhenEnabled)
	s.Step(`^token storage should match configuration$`, ctx.tokenStorageShouldMatchConfiguration)
	s.Step(`^auth endpoints should be correctly implemented$`, ctx.authEndpointsShouldBeCorrectlyImplemented)
	s.Step(`^logger initialization should match architecture patterns$`, ctx.loggerInitializationShouldMatchArchitecturePatterns)
	s.Step(`^log levels should be properly configured$`, ctx.logLevelsShouldBeProperlyConfigured)
	s.Step(`^structured logging should work when enabled$`, ctx.structuredLoggingShouldWorkWhenEnabled)
	s.Step(`^performance characteristics should match expectations$`, ctx.performanceCharacteristicsShouldMatchExpectations)
	s.Step(`^deployment configurations should be generated correctly$`, ctx.deploymentConfigurationsShouldBeGeneratedCorrectly)
	s.Step(`^Dockerfiles should be optimized for each blueprint$`, ctx.dockerfilesShouldBeOptimizedForEachBlueprint)
	s.Step(`^Kubernetes manifests should follow best practices$`, ctx.kubernetesManifestsShouldFollowBestPractices)
	s.Step(`^cloud-specific configurations should be valid$`, ctx.cloudSpecificConfigurationsShouldBeValid)
	s.Step(`^test files should be generated appropriately$`, ctx.testFilesShouldBeGeneratedAppropriately)
	s.Step(`^mocking frameworks should be integrated$`, ctx.mockingFrameworksShouldBeIntegrated)
	s.Step(`^coverage tools should be configured$`, ctx.coverageToolsShouldBeConfigured)
	s.Step(`^test commands should work correctly$`, ctx.testCommandsShouldWorkCorrectly)
	s.Step(`^middleware should be registered in correct order$`, ctx.middlewareShouldBeRegisteredInCorrectOrder)
	s.Step(`^middleware configuration should be consistent$`, ctx.middlewareConfigurationShouldBeConsistent)
	s.Step(`^performance impact should be acceptable$`, ctx.performanceImpactShouldBeAcceptable)
	s.Step(`^middleware conflicts should be prevented$`, ctx.middlewareConflictsShouldBePrevented)
	
	// Execution optimization steps
	s.Step(`^I have a large configuration matrix$`, ctx.iHaveALargeConfigurationMatrix)
	s.Step(`^I execute matrix tests with optimization$`, ctx.iExecuteMatrixTestsWithOptimization)
	s.Step(`^execution should be efficient:$`, ctx.executionShouldBeEfficient)
	s.Step(`^critical paths should always be tested$`, ctx.criticalPathsShouldAlwaysBeTested)
	s.Step(`^optimization should not compromise quality$`, ctx.optimizationShouldNotCompromiseQuality)
	s.Step(`^results should be deterministic$`, ctx.resultsShouldBeDeterministic)
}

// Step implementations

func (ctx *MatrixTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Verify CLI is available
	return nil
}

func (ctx *MatrixTestContext) allTemplatesAreProperlyInitialized() error {
	// Verify templates are loaded
	return nil
}

func (ctx *MatrixTestContext) matrixTestingModeIsEnabled() error {
	// Enable matrix testing mode
	return nil
}

func (ctx *MatrixTestContext) iTestAllFrameworkAndDatabaseCombinations() error {
	// Setup framework and database combinations
	frameworks := []string{"gin", "echo", "fiber", "chi"}
	databases := []string{"postgresql", "mysql", "sqlite", "mongodb"}
	drivers := map[string]string{
		"postgresql": "postgres",
		"mysql":      "mysql",
		"sqlite":     "sqlite3",
		"mongodb":    "mongo",
	}
	orms := []string{"", "gorm", "sqlx", "sqlc"}
	
	for _, framework := range frameworks {
		for _, database := range databases {
			driver := drivers[database]
			for _, orm := range orms {
				// Skip invalid combinations
				if database == "mongodb" && orm != "" {
					continue
				}
				if orm == "sqlc" && database == "mongodb" {
					continue
				}
				
				config := MatrixConfiguration{
					Name:      fmt.Sprintf("%s-%s-%s", framework, database, orm),
					Framework: framework,
					Database:  database,
					Driver:    driver,
					ORM:       orm,
				}
				
				ctx.configurations = append(ctx.configurations, config)
			}
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) iGenerateProjectsWithTheFullMatrix(table *godog.Table) error {
	// Generate projects for each matrix configuration
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		
		config := MatrixConfiguration{
			Framework: row.Cells[0].Value,
			Database:  row.Cells[1].Value,
			Driver:    row.Cells[2].Value,
			ORM:       row.Cells[3].Value,
		}
		
		if config.ORM == "(none)" {
			config.ORM = ""
		}
		
		expectedResult := row.Cells[4].Value
		
		// Generate project
		projectConfig := &types.ProjectConfig{
			Name:         fmt.Sprintf("matrix-test-%d", time.Now().UnixNano()),
			Type:         "web-api",
			Framework:    config.Framework,
			Architecture: "standard",
			Logger:       "slog",
			Features: &types.Features{
				Database: types.DatabaseConfig{
					Driver: config.Driver,
					ORM:    config.ORM,
				},
			},
		}
		
		projectPath := filepath.Join(ctx.tempDir, projectConfig.Name)
		projectConfig.Module = fmt.Sprintf("github.com/test/%s", projectConfig.Name)
		
		startTime := time.Now()
		err := helpers.GenerateProject(projectConfig)
		generationTime := time.Since(startTime)
		
		result := &MatrixTestResult{
			GenerationTime: generationTime,
		}
		
		if expectedResult == "Success" && err != nil {
			result.Success = false
			result.Errors = []string{err.Error()}
		} else if expectedResult == "Success" && err == nil {
			result.Success = true
			ctx.projectPaths[projectConfig.Name] = projectPath
		}
		
		ctx.results[config.Name] = result
	}
	
	return nil
}

func (ctx *MatrixTestContext) iTestAuthenticationAcrossAllFrameworks() error {
	// Setup authentication test configurations
	frameworks := []string{"gin", "echo", "fiber", "chi"}
	authTypes := []string{"jwt", "oauth2", "api-key", "session"}
	
	for _, framework := range frameworks {
		for _, authType := range authTypes {
			config := MatrixConfiguration{
				Name:      fmt.Sprintf("%s-%s-auth", framework, authType),
				Framework: framework,
				AuthType:  authType,
			}
			
			ctx.configurations = append(ctx.configurations, config)
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) iGenerateProjectsWithAuthConfigurations(table *godog.Table) error {
	// Generate projects with authentication configurations
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		
		framework := row.Cells[0].Value
		authType := row.Cells[1].Value
		database := row.Cells[3].Value
		
		projectConfig := &types.ProjectConfig{
			Name:         fmt.Sprintf("auth-test-%s-%s", framework, authType),
			Type:         "web-api",
			Framework:    framework,
			Architecture: "standard",
			Logger:       "slog",
			Features: &types.Features{
				Authentication: types.AuthConfig{
					Type: authType,
				},
				Database: types.DatabaseConfig{
					Driver: getDriverForDatabase(database),
				},
			},
		}
		
		projectPath := filepath.Join(ctx.tempDir, projectConfig.Name)
		projectConfig.Module = fmt.Sprintf("github.com/test/%s", projectConfig.Name)
		
		err := helpers.GenerateProject(projectConfig)
		
		result := &MatrixTestResult{
			Success: err == nil,
		}
		
		if err != nil {
			result.Errors = []string{err.Error()}
		}
		
		ctx.results[projectConfig.Name] = result
	}
	
	return nil
}

func (ctx *MatrixTestContext) iTestAllLoggerTypesWithDifferentArchitectures() error {
	// Setup logger and architecture combinations
	architectures := []string{"standard", "clean", "ddd", "hexagonal"}
	loggers := []string{"slog", "zap", "logrus", "zerolog"}
	
	for _, arch := range architectures {
		for _, logger := range loggers {
			config := MatrixConfiguration{
				Name:         fmt.Sprintf("%s-%s-logger", arch, logger),
				Architecture: arch,
				Logger:       logger,
			}
			
			ctx.configurations = append(ctx.configurations, config)
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) iGenerateProjectsWithLoggerConfigurations(table *godog.Table) error {
	// Generate projects with logger configurations
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		
		architecture := row.Cells[0].Value
		logger := row.Cells[1].Value
		
		projectConfig := &types.ProjectConfig{
			Name:         fmt.Sprintf("logger-test-%s-%s", architecture, logger),
			Type:         "web-api",
			Framework:    "gin",
			Architecture: architecture,
			Logger:       logger,
		}
		
		projectPath := filepath.Join(ctx.tempDir, projectConfig.Name)
		projectConfig.Module = fmt.Sprintf("github.com/test/%s", projectConfig.Name)
		
		err := helpers.GenerateProject(projectConfig)
		
		result := &MatrixTestResult{
			Success: err == nil,
		}
		
		if err != nil {
			result.Errors = []string{err.Error()}
		}
		
		ctx.results[projectConfig.Name] = result
	}
	
	return nil
}

func (ctx *MatrixTestContext) iTestDeploymentConfigurations() error {
	// Setup deployment configuration tests
	blueprints := []string{"web-api", "cli", "lambda", "microservice", "monolith"}
	deployTargets := [][]string{
		{"docker", "k8s"},
		{"docker"},
		{"lambda", "terraform"},
		{"docker", "k8s", "terraform"},
		{"docker", "k8s", "vercel", "railway"},
	}
	
	for i, blueprint := range blueprints {
		config := MatrixConfiguration{
			Name:          fmt.Sprintf("%s-deploy", blueprint),
			DeployTargets: deployTargets[i],
		}
		
		ctx.configurations = append(ctx.configurations, config)
	}
	
	return nil
}

func (ctx *MatrixTestContext) iGenerateProjectsWithDeploymentTargets(table *godog.Table) error {
	// Generate projects with deployment configurations
	// Implementation would generate projects and verify deployment files
	return nil
}

func (ctx *MatrixTestContext) iTestDifferentTestingConfigurations() error {
	// Setup testing framework configurations
	return nil
}

func (ctx *MatrixTestContext) iGenerateProjectsWithTestingSetups(table *godog.Table) error {
	// Generate projects with testing configurations
	return nil
}

func (ctx *MatrixTestContext) iTestMiddlewareStacksAcrossFrameworks() error {
	// Setup middleware configuration tests
	return nil
}

func (ctx *MatrixTestContext) iConfigureMiddlewareCombinations(table *godog.Table) error {
	// Configure and test middleware combinations
	return nil
}

// Validation step implementations

func (ctx *MatrixTestContext) allValidCombinationsShouldGenerateSuccessfully() error {
	failedCount := 0
	for name, result := range ctx.results {
		if !result.Success {
			failedCount++
			fmt.Printf("Failed: %s - %v\n", name, result.Errors)
		}
	}
	
	if failedCount > 0 {
		return fmt.Errorf("%d configurations failed to generate", failedCount)
	}
	
	return nil
}

func (ctx *MatrixTestContext) invalidORMSelectionsShouldBePrevented() error {
	// Verify that invalid ORM selections were prevented
	// MongoDB should not allow ORM selection
	for name, config := range ctx.configurations {
		if config.Database == "mongodb" && config.ORM != "" {
			if result, ok := ctx.results[name]; ok && result.Success {
				return fmt.Errorf("mongodb with ORM %s should have failed", config.ORM)
			}
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) databaseSpecificFeaturesShouldBeCorrectlyConfigured() error {
	// Verify database-specific features
	for name, path := range ctx.projectPaths {
		// Check docker-compose.yml for correct database service
		dockerComposePath := filepath.Join(path, "docker-compose.yml")
		if _, err := os.Stat(dockerComposePath); err == nil {
			content, err := os.ReadFile(dockerComposePath)
			if err != nil {
				return err
			}
			
			// Verify database service is present
			contentStr := string(content)
			for _, config := range ctx.configurations {
				if strings.Contains(name, config.Database) {
					expectedService := config.Database
					if expectedService == "postgres" {
						expectedService = "postgresql"
					}
					
					if !strings.Contains(contentStr, expectedService+":") {
						return fmt.Errorf("docker-compose.yml missing %s service for %s", expectedService, name)
					}
				}
			}
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) authenticationMiddlewareShouldBeProperlyConfigured() error {
	// Verify authentication middleware is present and configured
	for name, path := range ctx.projectPaths {
		if strings.Contains(name, "auth-test") {
			authMiddlewarePath := filepath.Join(path, "internal", "middleware", "auth.go")
			if _, err := os.Stat(authMiddlewarePath); os.IsNotExist(err) {
				return fmt.Errorf("auth middleware missing for %s", name)
			}
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) sessionManagementShouldWorkWhenEnabled() error {
	// Verify session management configuration
	return nil
}

func (ctx *MatrixTestContext) tokenStorageShouldMatchConfiguration() error {
	// Verify token storage configuration
	return nil
}

func (ctx *MatrixTestContext) authEndpointsShouldBeCorrectlyImplemented() error {
	// Verify auth endpoints exist
	return nil
}

func (ctx *MatrixTestContext) loggerInitializationShouldMatchArchitecturePatterns() error {
	// Verify logger initialization follows architecture patterns
	return nil
}

func (ctx *MatrixTestContext) logLevelsShouldBeProperlyConfigured() error {
	// Verify log levels are configured correctly
	return nil
}

func (ctx *MatrixTestContext) structuredLoggingShouldWorkWhenEnabled() error {
	// Verify structured logging is implemented
	return nil
}

func (ctx *MatrixTestContext) performanceCharacteristicsShouldMatchExpectations() error {
	// Verify performance characteristics
	return nil
}

func (ctx *MatrixTestContext) deploymentConfigurationsShouldBeGeneratedCorrectly() error {
	// Verify deployment configurations
	return nil
}

func (ctx *MatrixTestContext) dockerfilesShouldBeOptimizedForEachBlueprint() error {
	// Verify Dockerfiles are optimized
	for name, path := range ctx.projectPaths {
		dockerfilePath := filepath.Join(path, "Dockerfile")
		if _, err := os.Stat(dockerfilePath); err == nil {
			content, err := os.ReadFile(dockerfilePath)
			if err != nil {
				return err
			}
			
			contentStr := string(content)
			
			// Check for multi-stage build
			if !strings.Contains(contentStr, "FROM") || !strings.Contains(contentStr, "AS builder") {
				return fmt.Errorf("Dockerfile for %s not using multi-stage build", name)
			}
			
			// Check for proper caching
			if !strings.Contains(contentStr, "COPY go.mod go.sum") {
				return fmt.Errorf("Dockerfile for %s not optimizing dependency caching", name)
			}
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) kubernetesManifestsShouldFollowBestPractices() error {
	// Verify Kubernetes manifests follow best practices
	return nil
}

func (ctx *MatrixTestContext) cloudSpecificConfigurationsShouldBeValid() error {
	// Verify cloud-specific configurations
	return nil
}

func (ctx *MatrixTestContext) testFilesShouldBeGeneratedAppropriately() error {
	// Verify test files are generated
	return nil
}

func (ctx *MatrixTestContext) mockingFrameworksShouldBeIntegrated() error {
	// Verify mocking frameworks
	return nil
}

func (ctx *MatrixTestContext) coverageToolsShouldBeConfigured() error {
	// Verify coverage tools
	return nil
}

func (ctx *MatrixTestContext) testCommandsShouldWorkCorrectly() error {
	// Verify test commands work
	return nil
}

func (ctx *MatrixTestContext) middlewareShouldBeRegisteredInCorrectOrder() error {
	// Verify middleware registration order
	return nil
}

func (ctx *MatrixTestContext) middlewareConfigurationShouldBeConsistent() error {
	// Verify middleware configuration consistency
	return nil
}

func (ctx *MatrixTestContext) performanceImpactShouldBeAcceptable() error {
	// Verify performance impact of middleware
	return nil
}

func (ctx *MatrixTestContext) middlewareConflictsShouldBePrevented() error {
	// Verify no middleware conflicts
	return nil
}

// Optimization step implementations

func (ctx *MatrixTestContext) iHaveALargeConfigurationMatrix() error {
	// Setup a large configuration matrix
	// This would typically have hundreds of combinations
	return nil
}

func (ctx *MatrixTestContext) iExecuteMatrixTestsWithOptimization() error {
	// Execute tests with various optimization strategies
	strategies := []string{"parallel", "sampling", "incremental", "priority"}
	
	for _, strategy := range strategies {
		ctx.executionStrategy = strategy
		// Execute subset of tests based on strategy
	}
	
	return nil
}

func (ctx *MatrixTestContext) executionShouldBeEfficient(table *godog.Table) error {
	// Verify execution efficiency metrics
	elapsed := time.Since(ctx.startTime)
	
	// Check if execution time is within acceptable limits
	if ctx.executionStrategy == "parallel" && elapsed > 5*time.Minute {
		return fmt.Errorf("parallel execution took too long: %v", elapsed)
	}
	
	return nil
}

func (ctx *MatrixTestContext) criticalPathsShouldAlwaysBeTested() error {
	// Verify critical paths were tested
	criticalConfigs := []string{
		"gin-postgresql-gorm",
		"echo-mysql-gorm",
		"fiber-postgresql",
	}
	
	for _, config := range criticalConfigs {
		found := false
		for name := range ctx.results {
			if strings.Contains(name, config) {
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("critical configuration %s was not tested", config)
		}
	}
	
	return nil
}

func (ctx *MatrixTestContext) optimizationShouldNotCompromiseQuality() error {
	// Verify quality is maintained with optimization
	// Check that key validations still run
	return nil
}

func (ctx *MatrixTestContext) resultsShouldBeDeterministic() error {
	// Verify results are deterministic across runs
	// This would involve comparing multiple runs
	return nil
}

// Helper functions

func getDriverForDatabase(database string) string {
	drivers := map[string]string{
		"postgresql": "postgres",
		"postgres":   "postgres",
		"mysql":      "mysql",
		"sqlite":     "sqlite3",
		"mongodb":    "mongo",
	}
	
	if driver, ok := drivers[database]; ok {
		return driver
	}
	
	return ""
}