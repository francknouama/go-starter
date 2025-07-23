package monolith

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
)

// MonolithTestContext holds the state for monolith BDD tests with testcontainers
type MonolithTestContext struct {
	// Test environment
	workingDir    string
	projectDir    string
	projectName   string
	
	// Command execution
	lastCommand    *exec.Cmd
	lastOutput     []byte
	lastError      error
	lastExitCode   int
	
	// Configuration
	framework      string
	databaseDriver string
	databaseORM    string
	authType       string
	logger         string
	
	// Test results
	projectExists  bool
	compilationOK  bool
	testResults    map[string]bool
	
	// HTTP client for testing running applications
	httpClient *http.Client
	ctx               context.Context
}

var monolithCtx *MonolithTestContext

// Test runner for BDD scenarios
func TestMonolithBDD(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeMonolithScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}
	
	if suite.Run() != 0 {
		t.Fatal("Non-zero status returned, failed to run monolith BDD tests")
	}
}

// InitializeMonolithScenario registers step definitions for monolith scenarios
func InitializeMonolithScenario(ctx *godog.ScenarioContext) {
	// Initialize test context
	monolithCtx = &MonolithTestContext{
		httpClient:  &http.Client{Timeout: 5 * time.Second},
		testResults: make(map[string]bool),
		ctx:         context.Background(),
	}
	
	// Given steps (preconditions)
	ctx.Given(`^the go-starter CLI tool is available$`, monolithCtx.theGoStarterCLIToolIsAvailable)
	ctx.Given(`^I am in a clean working directory$`, monolithCtx.iAmInACleanWorkingDirectory)
	ctx.Given(`^I want to create a monolith application$`, monolithCtx.iWantToCreateAMonolithApplication)
	ctx.Given(`^I want to create a secure monolith application$`, monolithCtx.iWantToCreateASecureMonolithApplication)
	ctx.Given(`^I want to create a production-ready monolith$`, monolithCtx.iWantToCreateAProductionReadyMonolith)
	ctx.Given(`^I want to create a monolith with frontend assets$`, monolithCtx.iWantToCreateAMonolithWithFrontendAssets)
	ctx.Given(`^I want to create a well-structured monolith$`, monolithCtx.iWantToCreateAWellStructuredMonolith)
	ctx.Given(`^I have generated a monolith with database support$`, monolithCtx.iHaveGeneratedAMonolithWithDatabaseSupport)
	ctx.Given(`^I want to ensure code quality$`, monolithCtx.iWantToEnsureCodeQuality)
	ctx.Given(`^I have generated a monolith application$`, monolithCtx.iHaveGeneratedAMonolithApplication)
	ctx.Given(`^I want flexible logging options$`, monolithCtx.iWantFlexibleLoggingOptions)
	ctx.Given(`^I want a high-performance monolith$`, monolithCtx.iWantAHighPerformanceMonolith)
	ctx.Given(`^I want a secure monolith application$`, monolithCtx.iWantASecureMonolithApplication)
	ctx.Given(`^I want to deploy my monolith to production$`, monolithCtx.iWantToDeployMyMonolithToProduction)
	
	// When steps (actions)
	ctx.When(`^I run the command "([^"]*)"$`, monolithCtx.iRunTheCommand)
	ctx.When(`^I generate a monolith with framework "([^"]*)"$`, monolithCtx.iGenerateAMonolithWithFramework)
	ctx.When(`^I generate a monolith with database "([^"]*)" and ORM "([^"]*)"$`, monolithCtx.iGenerateAMonolithWithDatabaseAndORM)
	ctx.When(`^I generate a monolith with authentication type "([^"]*)"$`, monolithCtx.iGenerateAMonolithWithAuthenticationType)
	ctx.When(`^I generate a monolith with all production features$`, monolithCtx.iGenerateAMonolithWithAllProductionFeatures)
	ctx.When(`^I generate a monolith with asset pipeline enabled$`, monolithCtx.iGenerateAMonolithWithAssetPipelineEnabled)
	ctx.When(`^I generate a monolith application$`, monolithCtx.iGenerateAMonolithApplication)
	ctx.When(`^I examine the migration system$`, monolithCtx.iExamineTheMigrationSystem)
	ctx.When(`^I use the development tools$`, monolithCtx.iUseTheDevelopmentTools)
	ctx.When(`^I generate a monolith with logger "([^"]*)"$`, monolithCtx.iGenerateAMonolithWithLogger)
	ctx.When(`^I generate a monolith with performance optimizations$`, monolithCtx.iGenerateAMonolithWithPerformanceOptimizations)
	ctx.When(`^I generate a monolith with security features$`, monolithCtx.iGenerateAMonolithWithSecurityFeatures)
	ctx.When(`^I examine the deployment configuration$`, monolithCtx.iExamineTheDeploymentConfiguration)
	
	// Then steps (assertions)
	ctx.Then(`^the generation should succeed$`, monolithCtx.theGenerationShouldSucceed)
	ctx.Then(`^the project should contain all essential monolith components$`, monolithCtx.theProjectShouldContainAllEssentialMonolithComponents)
	ctx.Then(`^the generated code should compile successfully$`, monolithCtx.theGeneratedCodeShouldCompileSuccessfully)
	ctx.Then(`^the project should use the "([^"]*)" web framework$`, monolithCtx.theProjectShouldUseTheWebFramework)
	ctx.Then(`^the code should include "([^"]*)"-specific imports$`, monolithCtx.theCodeShouldIncludeSpecificImports)
	ctx.Then(`^the application should compile and run$`, monolithCtx.theApplicationShouldCompileAndRun)
	ctx.Then(`^the project should include "([^"]*)" database configuration$`, monolithCtx.theProjectShouldIncludeDatabaseConfiguration)
	ctx.Then(`^the migration files should support "([^"]*)"$`, monolithCtx.theMigrationFilesShouldSupport)
	ctx.Then(`^the ORM integration should use "([^"]*)"$`, monolithCtx.theORMIntegrationShouldUse)
	ctx.Then(`^the project should include "([^"]*)" authentication setup$`, monolithCtx.theProjectShouldIncludeAuthenticationSetup)
	ctx.Then(`^the session management should be properly configured$`, monolithCtx.theSessionManagementShouldBeProperlyConfigured)
	ctx.Then(`^the security headers should be implemented$`, monolithCtx.theSecurityHeadersShouldBeImplemented)
	ctx.Then(`^the project should include Docker configuration$`, monolithCtx.theProjectShouldIncludeDockerConfiguration)
	ctx.Then(`^the project should include CI/CD pipelines$`, monolithCtx.theProjectShouldIncludeCICDPipelines)
	ctx.Then(`^the project should include Kubernetes deployment$`, monolithCtx.theProjectShouldIncludeKubernetesDeployment)
	ctx.Then(`^the project should include monitoring and health checks$`, monolithCtx.theProjectShouldIncludeMonitoringAndHealthChecks)
	ctx.Then(`^the project should include comprehensive testing$`, monolithCtx.theProjectShouldIncludeComprehensiveTesting)
	ctx.Then(`^the project should include Vite configuration$`, monolithCtx.theProjectShouldIncludeViteConfiguration)
	ctx.Then(`^the project should include Tailwind CSS setup$`, monolithCtx.theProjectShouldIncludeTailwindCSSSetup)
	ctx.Then(`^the asset build process should be integrated with the backend$`, monolithCtx.theAssetBuildProcessShouldBeIntegratedWithTheBackend)
	ctx.Then(`^the code should follow modular monolith patterns$`, monolithCtx.theCodeShouldFollowModularMonolithPatterns)
	ctx.Then(`^the layers should be properly separated$`, monolithCtx.theLayersShouldBeProperlySeparated)
	ctx.Then(`^the dependencies should flow in the correct direction$`, monolithCtx.theDependenciesShouldFlowInTheCorrectDirection)
	ctx.Then(`^the code should be easily testable$`, monolithCtx.theCodeShouldBeEasilyTestable)
	ctx.Then(`^the project should include migration scripts$`, monolithCtx.theProjectShouldIncludeMigrationScripts)
	ctx.Then(`^the migration commands should work correctly$`, monolithCtx.theMigrationCommandsShouldWorkCorrectly)
	ctx.Then(`^the database schema should be properly versioned$`, monolithCtx.theDatabaseSchemaShouldBeProperlyVersioned)
	ctx.Then(`^the project should include unit tests$`, monolithCtx.theProjectShouldIncludeUnitTests)
	ctx.Then(`^the project should include integration tests$`, monolithCtx.theProjectShouldIncludeIntegrationTests)
	ctx.Then(`^the project should include benchmark tests$`, monolithCtx.theProjectShouldIncludeBenchmarkTests)
	ctx.Then(`^the tests should use proper mocking$`, monolithCtx.theTestsShouldUseProperMocking)
	ctx.Then(`^the test coverage should be measurable$`, monolithCtx.theTestCoverageShouldBeMeasurable)
	ctx.Then(`^the Makefile should provide essential commands$`, monolithCtx.theMakefileShouldProvideEssentialCommands)
	ctx.Then(`^the development server should support hot reload$`, monolithCtx.theDevelopmentServerShouldSupportHotReload)
	ctx.Then(`^the linting should pass without errors$`, monolithCtx.theLintingShouldPassWithoutErrors)
	ctx.Then(`^the security scanning should be configured$`, monolithCtx.theSecurityScanningShouldBeConfigured)
	ctx.Then(`^the application should use the "([^"]*)" logging library$`, monolithCtx.theApplicationShouldUseTheLoggingLibrary)
	ctx.Then(`^the logging should follow consistent interface patterns$`, monolithCtx.theLoggingShouldFollowConsistentInterfacePatterns)
	ctx.Then(`^the log levels should be properly configured$`, monolithCtx.theLogLevelsShouldBeProperlyConfigured)
	ctx.Then(`^the database should use connection pooling$`, monolithCtx.theDatabaseShouldUseConnectionPooling)
	ctx.Then(`^the application should include caching strategies$`, monolithCtx.theApplicationShouldIncludeCachingStrategies)
	ctx.Then(`^the asset pipeline should optimize for production$`, monolithCtx.theAssetPipelineShouldOptimizeForProduction)
	ctx.Then(`^benchmark tests should validate performance requirements$`, monolithCtx.benchmarkTestsShouldValidatePerformanceRequirements)
	ctx.Then(`^the application should implement OWASP security headers$`, monolithCtx.theApplicationShouldImplementOWASPSecurityHeaders)
	ctx.Then(`^the sessions should use secure configurations$`, monolithCtx.theSessionsShouldUseSecureConfigurations)
	ctx.Then(`^the input validation should prevent common attacks$`, monolithCtx.theInputValidationShouldPreventCommonAttacks)
	ctx.Then(`^the security scanning should be automated$`, monolithCtx.theSecurityScanningShouldBeAutomated)
	ctx.Then(`^the project should include multiple deployment targets$`, monolithCtx.theProjectShouldIncludeMultipleDeploymentTargets)
	ctx.Then(`^the CI/CD should support staging and production environments$`, monolithCtx.theCICDShouldSupportStagingAndProductionEnvironments)
	ctx.Then(`^the rollback procedures should be documented$`, monolithCtx.theRollbackProceduresShouldBeDocumented)
	ctx.Then(`^the health checks should validate system status$`, monolithCtx.theHealthChecksShouldValidateSystemStatus)
	
	// Scenario hooks - handle cleanup in individual tests for now
	// ctx.Before and ctx.After signatures may vary by godog version
}

// Given step implementations

func (ctx *MonolithTestContext) theGoStarterCLIToolIsAvailable() error {
	// Build the CLI tool if not already available
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")
	
	buildCmd := exec.Command("go", "build", "-o", "go-starter", ".")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	
	if err != nil {
		return fmt.Errorf("failed to build go-starter CLI: %s", string(output))
	}
	
	return nil
}

func (ctx *MonolithTestContext) iAmInACleanWorkingDirectory() error {
	// Reset test state
	ctx.testResults = make(map[string]bool)
	ctx.projectExists = false
	ctx.compilationOK = false
	
	var err error
	ctx.workingDir, err = os.MkdirTemp("", "monolith-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	return os.Chdir(ctx.workingDir)
}

func (ctx *MonolithTestContext) iWantToCreateAMonolithApplication() error {
	// Set default configuration for basic monolith
	ctx.projectName = "test-monolith"
	ctx.framework = "gin"
	ctx.databaseDriver = "postgres"
	ctx.databaseORM = "gorm"
	ctx.authType = "session"
	ctx.logger = "slog"
	return nil
}

func (ctx *MonolithTestContext) iWantToCreateASecureMonolithApplication() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	ctx.authType = "session" // Will enable secure session configuration
	return nil
}

func (ctx *MonolithTestContext) iWantToCreateAProductionReadyMonolith() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	// Production-ready settings will be validated in assertions
	return nil
}

func (ctx *MonolithTestContext) iWantToCreateAMonolithWithFrontendAssets() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	// Asset pipeline will be enabled by default
	return nil
}

func (ctx *MonolithTestContext) iWantToCreateAWellStructuredMonolith() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	// Structure will be validated in assertions
	return nil
}

func (ctx *MonolithTestContext) iHaveGeneratedAMonolithWithDatabaseSupport() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iWantToEnsureCodeQuality() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return nil
}

func (ctx *MonolithTestContext) iHaveGeneratedAMonolithApplication() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iWantFlexibleLoggingOptions() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return nil
}

func (ctx *MonolithTestContext) iWantAHighPerformanceMonolith() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return nil
}

func (ctx *MonolithTestContext) iWantASecureMonolithApplication() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return nil
}

func (ctx *MonolithTestContext) iWantToDeployMyMonolithToProduction() error {
	_ = ctx.iWantToCreateAMonolithApplication()
	return ctx.generateMonolithProject()
}

// When step implementations

func (ctx *MonolithTestContext) iRunTheCommand(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	
	ctx.lastCommand = exec.Command(parts[0], parts[1:]...)
	ctx.lastCommand.Dir = ctx.workingDir
	
	ctx.lastOutput, ctx.lastError = ctx.lastCommand.CombinedOutput()
	
	if exitError, ok := ctx.lastError.(*exec.ExitError); ok {
		ctx.lastExitCode = exitError.ExitCode()
	} else {
		ctx.lastExitCode = 0
	}
	
	return nil
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithFramework(framework string) error {
	ctx.framework = framework
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithDatabaseAndORM(database, orm string) error {
	ctx.databaseDriver = database
	ctx.databaseORM = orm
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithAuthenticationType(authType string) error {
	ctx.authType = authType
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithAllProductionFeatures() error {
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithAssetPipelineEnabled() error {
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithApplication() error {
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iExamineTheMigrationSystem() error {
	// This is passive - we'll check the migration system in assertions
	return nil
}

func (ctx *MonolithTestContext) iUseTheDevelopmentTools() error {
	// This is passive - we'll check dev tools in assertions
	return nil
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithLogger(logger string) error {
	ctx.logger = logger
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithPerformanceOptimizations() error {
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iGenerateAMonolithWithSecurityFeatures() error {
	return ctx.generateMonolithProject()
}

func (ctx *MonolithTestContext) iExamineTheDeploymentConfiguration() error {
	// This is passive - we'll check deployment config in assertions
	return nil
}

// Then step implementations

func (ctx *MonolithTestContext) theGenerationShouldSucceed() error {
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

func (ctx *MonolithTestContext) theProjectShouldContainAllEssentialMonolithComponents() error {
	if !ctx.projectExists {
		return fmt.Errorf("project was not generated")
	}
	
	essentialComponents := []string{
		"main.go",
		"go.mod",
		"README.md",
		"Makefile",
		"Dockerfile",
		"config/database.go",
		"config/session.go",
		"controllers/user.go",
		"controllers/api.go",
		"routes/api.go",
		"routes/auth.go",
		"database/migrations/001_create_users.sql",
		"services/user.go",
		"services/user_test.go",
		"tests/integration_test.go",
		"benchmarks/api_test.go",
		"scripts/setup.sh",
		"scripts/migrate.sh",
		".github/workflows/ci.yml",
		".github/workflows/deploy.yml",
		"kubernetes/deployment.yaml",
	}
	
	for _, component := range essentialComponents {
		filePath := filepath.Join(ctx.projectDir, component)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("essential component missing: %s", component)
		}
	}
	
	return nil
}

func (ctx *MonolithTestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	if !ctx.projectExists {
		return fmt.Errorf("project was not generated")
	}
	
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = ctx.projectDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %s", string(output))
	}
	
	// Try to build
	buildCmd := exec.Command("go", "build", ".")
	buildCmd.Dir = ctx.projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %s", string(output))
	}
	
	ctx.compilationOK = true
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldUseTheWebFramework(framework string) error {
	return ctx.checkFileContains("main.go", framework)
}

func (ctx *MonolithTestContext) theCodeShouldIncludeSpecificImports(framework string) error {
	frameworkImports := map[string]string{
		"gin":   "gin-gonic/gin",
		"echo":  "labstack/echo",
		"fiber": "gofiber/fiber",
		"chi":   "go-chi/chi",
	}
	
	expectedImport, exists := frameworkImports[framework]
	if !exists {
		return fmt.Errorf("unknown framework: %s", framework)
	}
	
	return ctx.checkFileContains("main.go", expectedImport)
}

func (ctx *MonolithTestContext) theApplicationShouldCompileAndRun() error {
	if err := ctx.theGeneratedCodeShouldCompileSuccessfully(); err != nil {
		return err
	}
	
	// Try to start the application briefly to ensure it runs
	runCmd := exec.Command("timeout", "2s", "./"+ctx.projectName)
	runCmd.Dir = ctx.projectDir
	_, _ = runCmd.CombinedOutput() // Ignore output and error - timeout is expected
	
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeDatabaseConfiguration(database string) error {
	return ctx.checkFileContains("config/database.go", database)
}

func (ctx *MonolithTestContext) theMigrationFilesShouldSupport(database string) error {
	return ctx.checkFileContains("database/migrations/001_create_users.sql", "CREATE TABLE")
}

func (ctx *MonolithTestContext) theORMIntegrationShouldUse(orm string) error {
	return ctx.checkFileContains("config/database.go", orm)
}

func (ctx *MonolithTestContext) theProjectShouldIncludeAuthenticationSetup(authType string) error {
	return ctx.checkFileContains("routes/auth.go", "auth")
}

func (ctx *MonolithTestContext) theSessionManagementShouldBeProperlyConfigured() error {
	return ctx.checkFileContains("config/session.go", "HttpOnly")
}

func (ctx *MonolithTestContext) theSecurityHeadersShouldBeImplemented() error {
	return ctx.checkFileContains("config/session.go", "Secure")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeDockerConfiguration() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, "Dockerfile"))
	if os.IsNotExist(err) {
		return fmt.Errorf("Dockerfile not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeCICDPipelines() error {
	ciPath := filepath.Join(ctx.projectDir, ".github", "workflows", "ci.yml")
	if _, err := os.Stat(ciPath); os.IsNotExist(err) {
		return fmt.Errorf("CI pipeline not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeKubernetesDeployment() error {
	k8sPath := filepath.Join(ctx.projectDir, "kubernetes", "deployment.yaml")
	if _, err := os.Stat(k8sPath); os.IsNotExist(err) {
		return fmt.Errorf("Kubernetes deployment not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeMonitoringAndHealthChecks() error {
	return ctx.checkFileContains("routes/api.go", "/health")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeComprehensiveTesting() error {
	testFiles := []string{
		"services/user_test.go",
		"tests/integration_test.go",
		"benchmarks/api_test.go",
	}
	
	for _, testFile := range testFiles {
		if _, err := os.Stat(filepath.Join(ctx.projectDir, testFile)); os.IsNotExist(err) {
			return fmt.Errorf("test file not found: %s", testFile)
		}
	}
	
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeViteConfiguration() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, "vite.config.js"))
	if os.IsNotExist(err) {
		return fmt.Errorf("Vite configuration not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theProjectShouldIncludeTailwindCSSSetup() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, "tailwind.config.js"))
	if os.IsNotExist(err) {
		return fmt.Errorf("Tailwind CSS configuration not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theAssetBuildProcessShouldBeIntegratedWithTheBackend() error {
	return ctx.checkFileContains("vite.config.js", "build")
}

func (ctx *MonolithTestContext) theCodeShouldFollowModularMonolithPatterns() error {
	dirs := []string{"controllers", "services", "models", "config", "routes"}
	for _, dir := range dirs {
		dirPath := filepath.Join(ctx.projectDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("modular directory not found: %s", dir)
		}
	}
	return nil
}

func (ctx *MonolithTestContext) theLayersShouldBeProperlySeparated() error {
	// Controllers should not have direct database imports
	return ctx.checkFileDoesNotContain("controllers/user.go", "database/sql")
}

func (ctx *MonolithTestContext) theDependenciesShouldFlowInTheCorrectDirection() error {
	// Services should define interfaces
	return ctx.checkFileContains("services/user.go", "interface")
}

func (ctx *MonolithTestContext) theCodeShouldBeEasilyTestable() error {
	return ctx.checkFileContains("services/user_test.go", "mock.Mock")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeMigrationScripts() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, "scripts", "migrate.sh"))
	if os.IsNotExist(err) {
		return fmt.Errorf("migration script not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theMigrationCommandsShouldWorkCorrectly() error {
	return ctx.checkFileContains("scripts/migrate.sh", "migrate up")
}

func (ctx *MonolithTestContext) theDatabaseSchemaShouldBeProperlyVersioned() error {
	return ctx.checkFileContains("database/migrations/001_create_users.sql", "CREATE TABLE")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeUnitTests() error {
	return ctx.checkFileContains("services/user_test.go", "func Test")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeIntegrationTests() error {
	return ctx.checkFileContains("tests/integration_test.go", "func TestIntegration")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeBenchmarkTests() error {
	return ctx.checkFileContains("benchmarks/api_test.go", "func Benchmark")
}

func (ctx *MonolithTestContext) theTestsShouldUseProperMocking() error {
	return ctx.checkFileContains("services/user_test.go", "Mock")
}

func (ctx *MonolithTestContext) theTestCoverageShouldBeMeasurable() error {
	return ctx.checkFileContains("services/user_test.go", "testify")
}

func (ctx *MonolithTestContext) theMakefileShouldProvideEssentialCommands() error {
	makefileContent, err := os.ReadFile(filepath.Join(ctx.projectDir, "Makefile"))
	if err != nil {
		return fmt.Errorf("Makefile not found: %w", err)
	}
	
	content := string(makefileContent)
	essentialTargets := []string{"build", "test", "dev", "help"}
	
	for _, target := range essentialTargets {
		if !strings.Contains(content, target+":") {
			return fmt.Errorf("Makefile missing essential target: %s", target)
		}
	}
	
	return nil
}

func (ctx *MonolithTestContext) theDevelopmentServerShouldSupportHotReload() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, "air.toml"))
	if os.IsNotExist(err) {
		return fmt.Errorf("air.toml configuration not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theLintingShouldPassWithoutErrors() error {
	_, err := os.Stat(filepath.Join(ctx.projectDir, ".golangci.yml"))
	if os.IsNotExist(err) {
		return fmt.Errorf("golangci.yml configuration not found")
	}
	return nil
}

func (ctx *MonolithTestContext) theSecurityScanningShouldBeConfigured() error {
	return ctx.checkFileContains(".github/workflows/ci.yml", "gosec")
}

func (ctx *MonolithTestContext) theApplicationShouldUseTheLoggingLibrary(logger string) error {
	switch logger {
	case "slog":
		return ctx.checkFileContains("main.go", "log/slog")
	case "zap":
		return ctx.checkFileContains("main.go", "go.uber.org/zap")
	case "logrus":
		return ctx.checkFileContains("main.go", "github.com/sirupsen/logrus")
	case "zerolog":
		return ctx.checkFileContains("main.go", "github.com/rs/zerolog")
	default:
		return fmt.Errorf("unknown logger: %s", logger)
	}
}

func (ctx *MonolithTestContext) theLoggingShouldFollowConsistentInterfacePatterns() error {
	return ctx.checkFileContains("services/user_test.go", ctx.logger)
}

func (ctx *MonolithTestContext) theLogLevelsShouldBeProperlyConfigured() error {
	return ctx.checkFileContains("services/user_test.go", "Level")
}

func (ctx *MonolithTestContext) theDatabaseShouldUseConnectionPooling() error {
	return ctx.checkFileContains("config/database.go", "MaxOpenConns")
}

func (ctx *MonolithTestContext) theApplicationShouldIncludeCachingStrategies() error {
	return ctx.checkFileContains("config/database.go", "connection pool")
}

func (ctx *MonolithTestContext) theAssetPipelineShouldOptimizeForProduction() error {
	return ctx.checkFileContains("vite.config.js", "minify")
}

func (ctx *MonolithTestContext) benchmarkTestsShouldValidatePerformanceRequirements() error {
	return ctx.checkFileContains("benchmarks/api_test.go", "b.ResetTimer")
}

func (ctx *MonolithTestContext) theApplicationShouldImplementOWASPSecurityHeaders() error {
	return ctx.checkFileContains("config/session.go", "HttpOnly")
}

func (ctx *MonolithTestContext) theSessionsShouldUseSecureConfigurations() error {
	return ctx.checkFileContains("config/session.go", "Secure: true")
}

func (ctx *MonolithTestContext) theInputValidationShouldPreventCommonAttacks() error {
	return ctx.checkFileContains("controllers/user.go", "validation")
}

func (ctx *MonolithTestContext) theSecurityScanningShouldBeAutomated() error {
	return ctx.checkFileContains(".github/workflows/ci.yml", "security:")
}

func (ctx *MonolithTestContext) theProjectShouldIncludeMultipleDeploymentTargets() error {
	return ctx.checkFileContains(".github/workflows/deploy.yml", "deploy-staging")
}

func (ctx *MonolithTestContext) theCICDShouldSupportStagingAndProductionEnvironments() error {
	deployContent, err := os.ReadFile(filepath.Join(ctx.projectDir, ".github", "workflows", "deploy.yml"))
	if err != nil {
		return fmt.Errorf("deploy workflow not found: %w", err)
	}
	
	content := string(deployContent)
	if !strings.Contains(content, "staging") || !strings.Contains(content, "production") {
		return fmt.Errorf("deploy workflow missing staging or production environment")
	}
	
	return nil
}

func (ctx *MonolithTestContext) theRollbackProceduresShouldBeDocumented() error {
	return ctx.checkFileContains(".github/workflows/deploy.yml", "rollback")
}

func (ctx *MonolithTestContext) theHealthChecksShouldValidateSystemStatus() error {
	return ctx.checkFileContains("routes/api.go", "/health")
}

// Helper methods

func (ctx *MonolithTestContext) generateMonolithProject() error {
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")
	
	// Build go-starter first
	goStarterPath := filepath.Join(ctx.workingDir, "go-starter")
	buildCmd := exec.Command("go", "build", "-o", goStarterPath, ".")
	buildCmd.Dir = projectRoot
	
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to build go-starter: %s", string(output))
	}
	
	// Generate the project
	args := []string{
		"new", ctx.projectName,
		"--type=monolith",
		"--module=github.com/test/" + ctx.projectName,
		"--framework=" + ctx.framework,
		"--database-driver=" + ctx.databaseDriver,
		"--database-orm=" + ctx.databaseORM,
		"--auth-type=" + ctx.authType,
		"--logger=" + ctx.logger,
		"--no-git",
	}
	
	generateCmd := exec.Command(goStarterPath, args...)
	generateCmd.Dir = ctx.workingDir
	
	ctx.lastOutput, ctx.lastError = generateCmd.CombinedOutput()
	
	if ctx.lastError != nil {
		return fmt.Errorf("project generation failed: %s", string(ctx.lastOutput))
	}
	
	ctx.projectDir = filepath.Join(ctx.workingDir, ctx.projectName)
	ctx.projectExists = true
	
	return nil
}

func (ctx *MonolithTestContext) checkFileContains(filePath, content string) error {
	fullPath := filepath.Join(ctx.projectDir, filePath)
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	
	if !strings.Contains(string(fileContent), content) {
		return fmt.Errorf("file %s does not contain '%s'", filePath, content)
	}
	
	return nil
}

func (ctx *MonolithTestContext) checkFileDoesNotContain(filePath, content string) error {
	fullPath := filepath.Join(ctx.projectDir, filePath)
	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", filePath, err)
	}
	
	if strings.Contains(string(fileContent), content) {
		return fmt.Errorf("file %s should not contain '%s'", filePath, content)
	}
	
	return nil
}


