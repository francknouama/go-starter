package monolith

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MonolithAcceptanceTestSuite provides comprehensive ATDD coverage for monolith blueprint
type MonolithAcceptanceTestSuite struct {
	workingDir    string
	projectDir    string
	projectName   string
	originalDir   string
	projectRoot   string
	httpClient    *http.Client
}

func setupAcceptanceTest(t *testing.T) *MonolithAcceptanceTestSuite {
	suite := &MonolithAcceptanceTestSuite{
		projectName: "test-monolith-acceptance",
		httpClient:  &http.Client{Timeout: 5 * time.Second},
	}

	var err error
	suite.originalDir, err = os.Getwd()
	require.NoError(t, err)

	suite.projectRoot = filepath.Join(suite.originalDir, "..", "..", "..", "..")
	
	suite.workingDir, err = os.MkdirTemp("", "monolith-acceptance-*")
	require.NoError(t, err)

	err = os.Chdir(suite.workingDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Chdir(suite.originalDir)
		_ = os.RemoveAll(suite.workingDir)
	})

	return suite
}

func (suite *MonolithAcceptanceTestSuite) buildCLI(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-o", "go-starter", ".")
	buildCmd.Dir = suite.projectRoot
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build go-starter CLI: %s", string(output))
}

func (suite *MonolithAcceptanceTestSuite) generateMonolithProject(t *testing.T, args ...string) {
	suite.buildCLI(t)

	baseArgs := []string{
		"new", suite.projectName,
		"--type=monolith",
		"--module=github.com/test/" + suite.projectName,
		"--framework=gin",
		"--database-driver=postgres",
		"--database-orm=gorm",
		"--auth-type=session",
		"--logger=slog",
		"--no-git",
	}

	allArgs := append(baseArgs, args...)
	generateCmd := exec.Command("./go-starter", allArgs...)
	generateCmd.Dir = suite.workingDir

	output, err := generateCmd.CombinedOutput()
	require.NoError(t, err, "Project generation should succeed: %s", string(output))

	suite.projectDir = filepath.Join(suite.workingDir, suite.projectName)
	assert.DirExists(t, suite.projectDir, "Project directory should be created")
}

func (suite *MonolithAcceptanceTestSuite) checkFileExists(t *testing.T, relativePath string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	assert.FileExists(t, fullPath, "File should exist: %s", relativePath)
}

func (suite *MonolithAcceptanceTestSuite) checkFileContains(t *testing.T, relativePath, content string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	fileContent, err := os.ReadFile(fullPath)
	require.NoError(t, err, "Should be able to read file: %s", relativePath)
	assert.Contains(t, string(fileContent), content, "File %s should contain '%s'", relativePath, content)
}

func (suite *MonolithAcceptanceTestSuite) checkFileDoesNotContain(t *testing.T, relativePath, content string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	fileContent, err := os.ReadFile(fullPath)
	require.NoError(t, err, "Should be able to read file: %s", relativePath)
	assert.NotContains(t, string(fileContent), content, "File %s should not contain '%s'", relativePath, content)
}

func (suite *MonolithAcceptanceTestSuite) compileProject(t *testing.T) {
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = suite.projectDir
	output, err := modCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

	// Build the project
	buildCmd := exec.Command("go", "build", "-o", "monolith-app", ".")
	buildCmd.Dir = suite.projectDir
	output, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "Project should compile successfully: %s", string(output))
}

// ATDD Test Scenarios

func TestMonolithAcceptance_BasicGeneration(t *testing.T) {
	// GIVEN: A developer wants to create a monolith application
	// WHEN: They generate a monolith project with basic settings
	// THEN: A complete, working monolith application is generated

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Verify essential architectural components
	essentialFiles := []string{
		"main.go", "go.mod", "README.md", "Makefile", "Dockerfile",
		"config/database.go", "config/session.go",
		"controllers/user.go", "controllers/api.go",
		"routes/api.go", "routes/auth.go",
		"services/user.go", "services/user_test.go",
		"database/migrations/001_create_users.sql",
		"tests/integration_test.go", "benchmarks/api_test.go",
		"scripts/setup.sh", "scripts/migrate.sh",
		".github/workflows/ci.yml", ".github/workflows/deploy.yml",
		"kubernetes/deployment.yaml", "vite.config.js", "tailwind.config.js",
	}

	for _, file := range essentialFiles {
		suite.checkFileExists(t, file)
	}

	// Verify the project compiles
	suite.compileProject(t)
}

func TestMonolithAcceptance_MultiFrameworkSupport(t *testing.T) {
	// GIVEN: A developer wants to use different web frameworks
	// WHEN: They generate monolith projects with different frameworks
	// THEN: Each project should use the appropriate framework

	frameworks := []struct {
		name        string
		importPath  string
	}{
		{"gin", "gin-gonic/gin"},
		{"echo", "labstack/echo"},
		{"fiber", "gofiber/fiber"},
		{"chi", "go-chi/chi"},
	}

	for _, fw := range frameworks {
		t.Run(fw.name, func(t *testing.T) {
			suite := setupAcceptanceTest(t)
			suite.projectName = fmt.Sprintf("test-%s-monolith", fw.name)
			
			suite.generateMonolithProject(t, "--framework="+fw.name)
			
			// Verify framework-specific imports
			suite.checkFileContains(t, "main.go", fw.importPath)
			
			// Ensure it compiles
			suite.compileProject(t)
		})
	}
}

func TestMonolithAcceptance_DatabaseDriverSupport(t *testing.T) {
	// GIVEN: A developer wants to use different database drivers
	// WHEN: They generate monolith projects with different databases
	// THEN: Each project should have appropriate database configuration

	databases := []struct {
		driver string
		orm    string
	}{
		{"postgres", "gorm"},
		{"mysql", "gorm"},
		{"sqlite", "gorm"},
		{"postgres", "sqlx"},
	}

	for _, db := range databases {
		t.Run(fmt.Sprintf("%s_%s", db.driver, db.orm), func(t *testing.T) {
			suite := setupAcceptanceTest(t)
			suite.projectName = fmt.Sprintf("test-%s-%s", db.driver, db.orm)
			
			suite.generateMonolithProject(t, 
				"--database-driver="+db.driver,
				"--database-orm="+db.orm)
			
			// Verify database configuration
			suite.checkFileContains(t, "config/database.go", db.driver)
			suite.checkFileContains(t, "config/database.go", db.orm)
			
			// Check migration file exists
			suite.checkFileExists(t, "database/migrations/001_create_users.sql")
			
			// Ensure it compiles
			suite.compileProject(t)
		})
	}
}

func TestMonolithAcceptance_AuthenticationTypes(t *testing.T) {
	// GIVEN: A developer wants different authentication mechanisms
	// WHEN: They generate monolith projects with different auth types
	// THEN: Each project should have appropriate auth configuration

	authTypes := []string{"session", "jwt", "oauth2"}

	for _, authType := range authTypes {
		t.Run(authType, func(t *testing.T) {
			suite := setupAcceptanceTest(t)
			suite.projectName = fmt.Sprintf("test-%s-auth", authType)
			
			suite.generateMonolithProject(t, "--auth-type="+authType)
			
			// Verify auth configuration
			suite.checkFileExists(t, "routes/auth.go")
			suite.checkFileContains(t, "routes/auth.go", "auth")
			
			if authType == "session" {
				suite.checkFileExists(t, "config/session.go")
				suite.checkFileContains(t, "config/session.go", "HttpOnly")
				suite.checkFileContains(t, "config/session.go", "Secure")
			}
			
			// Ensure it compiles
			suite.compileProject(t)
		})
	}
}

func TestMonolithAcceptance_LoggerSupport(t *testing.T) {
	// GIVEN: A developer wants to use different logging libraries
	// WHEN: They generate monolith projects with different loggers
	// THEN: Each project should use the appropriate logger

	loggers := []struct {
		name       string
		importPath string
	}{
		{"slog", "log/slog"},
		{"zap", "go.uber.org/zap"},
		{"logrus", "github.com/sirupsen/logrus"},
		{"zerolog", "github.com/rs/zerolog"},
	}

	for _, logger := range loggers {
		t.Run(logger.name, func(t *testing.T) {
			suite := setupAcceptanceTest(t)
			suite.projectName = fmt.Sprintf("test-%s-logger", logger.name)
			
			suite.generateMonolithProject(t, "--logger="+logger.name)
			
			// Verify logger imports in main.go
			suite.checkFileContains(t, "main.go", logger.importPath)
			
			// Verify logger is used in tests
			suite.checkFileContains(t, "services/user_test.go", logger.name)
			
			// Ensure it compiles
			suite.compileProject(t)
		})
	}
}

func TestMonolithAcceptance_ProductionReadiness(t *testing.T) {
	// GIVEN: A developer wants a production-ready monolith
	// WHEN: They generate a monolith project
	// THEN: The project should include all production-ready features

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Docker configuration
	suite.checkFileExists(t, "Dockerfile")
	suite.checkFileExists(t, "docker-compose.yml")
	suite.checkFileExists(t, "docker-compose.prod.yml")

	// CI/CD pipelines
	suite.checkFileExists(t, ".github/workflows/ci.yml")
	suite.checkFileExists(t, ".github/workflows/deploy.yml")

	// Verify CI includes essential jobs
	suite.checkFileContains(t, ".github/workflows/ci.yml", "test:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "lint:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "security:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "build:")

	// Verify deployment includes staging and production
	suite.checkFileContains(t, ".github/workflows/deploy.yml", "deploy-staging")
	suite.checkFileContains(t, ".github/workflows/deploy.yml", "deploy-production")
	suite.checkFileContains(t, ".github/workflows/deploy.yml", "rollback")

	// Kubernetes deployment
	suite.checkFileExists(t, "kubernetes/deployment.yaml")

	// Security scanning
	suite.checkFileContains(t, ".github/workflows/ci.yml", "gosec")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "govulncheck")

	// Monitoring and health checks
	suite.checkFileContains(t, "routes/api.go", "/health")
	suite.checkFileContains(t, "routes/api.go", "/ready")

	// Development tools
	suite.checkFileExists(t, ".golangci.yml")
	suite.checkFileExists(t, "air.toml")

	// Ensure it compiles
	suite.compileProject(t)
}

func TestMonolithAcceptance_AssetPipelineIntegration(t *testing.T) {
	// GIVEN: A developer wants frontend asset management
	// WHEN: They generate a monolith project
	// THEN: The project should include modern asset pipeline

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Asset pipeline configuration
	suite.checkFileExists(t, "vite.config.js")
	suite.checkFileExists(t, "tailwind.config.js")

	// Verify Vite configuration quality
	suite.checkFileContains(t, "vite.config.js", "build:")
	suite.checkFileContains(t, "vite.config.js", "rollupOptions")
	suite.checkFileContains(t, "vite.config.js", "legacy")

	// Verify Tailwind configuration
	suite.checkFileContains(t, "tailwind.config.js", "content:")
	suite.checkFileContains(t, "tailwind.config.js", "theme:")
	suite.checkFileContains(t, "tailwind.config.js", "plugins:")
}

func TestMonolithAcceptance_ModularArchitecture(t *testing.T) {
	// GIVEN: A developer wants well-structured code
	// WHEN: They generate a monolith project
	// THEN: The code should follow modular monolith patterns

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Verify modular structure exists
	layers := []string{"controllers", "services", "models", "config", "routes"}
	for _, layer := range layers {
		suite.checkFileExists(t, layer+"/")
	}

	// Verify separation of concerns
	suite.checkFileDoesNotContain(t, "controllers/user.go", "database/sql") // Controllers shouldn't have direct DB access
	suite.checkFileContains(t, "services/user.go", "interface")             // Services should define interfaces

	// Verify testability
	suite.checkFileContains(t, "services/user_test.go", "Mock")
	suite.checkFileContains(t, "services/user_test.go", "testify")

	// Ensure it compiles
	suite.compileProject(t)
}

func TestMonolithAcceptance_DatabaseMigrationSystem(t *testing.T) {
	// GIVEN: A developer needs database management
	// WHEN: They generate a monolith project with database support
	// THEN: The project should include migration system

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Migration infrastructure
	suite.checkFileExists(t, "database/migrations/001_create_users.sql")
	suite.checkFileExists(t, "scripts/migrate.sh")

	// Verify migration script capabilities
	suite.checkFileContains(t, "scripts/migrate.sh", "migrate up")
	suite.checkFileContains(t, "scripts/migrate.sh", "migrate down")
	suite.checkFileContains(t, "scripts/migrate.sh", "create_migration")
	suite.checkFileContains(t, "scripts/migrate.sh", "test_connection")

	// Verify migration file quality
	suite.checkFileContains(t, "database/migrations/001_create_users.sql", "CREATE TABLE")
	suite.checkFileContains(t, "database/migrations/001_create_users.sql", "users")
	suite.checkFileContains(t, "database/migrations/001_create_users.sql", "INDEX")
}

func TestMonolithAcceptance_ComprehensiveTestCoverage(t *testing.T) {
	// GIVEN: A developer wants quality assurance
	// WHEN: They generate a monolith project
	// THEN: The project should include comprehensive testing

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Test files exist
	testFiles := []string{
		"services/user_test.go",        // Unit tests
		"tests/integration_test.go",    // Integration tests
		"benchmarks/api_test.go",       // Benchmark tests
	}

	for _, testFile := range testFiles {
		suite.checkFileExists(t, testFile)
	}

	// Verify unit test quality
	suite.checkFileContains(t, "services/user_test.go", "func Test")
	suite.checkFileContains(t, "services/user_test.go", "suite.Suite")
	suite.checkFileContains(t, "services/user_test.go", "mock.Mock")
	suite.checkFileContains(t, "services/user_test.go", "testify")

	// Verify integration test structure
	suite.checkFileContains(t, "tests/integration_test.go", "func TestIntegration")
	suite.checkFileContains(t, "tests/integration_test.go", "http.Client")

	// Verify benchmark tests
	suite.checkFileContains(t, "benchmarks/api_test.go", "func Benchmark")
	suite.checkFileContains(t, "benchmarks/api_test.go", "b.ResetTimer")
	suite.checkFileContains(t, "benchmarks/api_test.go", "b.RunParallel")

	// Verify CI includes testing
	suite.checkFileContains(t, ".github/workflows/ci.yml", "go test")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "coverage")
}

func TestMonolithAcceptance_SecurityHardening(t *testing.T) {
	// GIVEN: A developer wants a secure application
	// WHEN: They generate a monolith project
	// THEN: The project should implement security best practices

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Session security
	suite.checkFileContains(t, "config/session.go", "HttpOnly: true")
	suite.checkFileContains(t, "config/session.go", "Secure: true")
	suite.checkFileContains(t, "config/session.go", "SameSite")

	// Input validation
	suite.checkFileContains(t, "controllers/user.go", "validation")

	// Security scanning in CI
	suite.checkFileContains(t, ".github/workflows/ci.yml", "gosec")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "govulncheck")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "SARIF")
}

func TestMonolithAcceptance_PerformanceOptimization(t *testing.T) {
	// GIVEN: A developer wants high performance
	// WHEN: They generate a monolith project
	// THEN: The project should include performance optimizations

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Database performance
	suite.checkFileContains(t, "config/database.go", "MaxOpenConns")
	suite.checkFileContains(t, "config/database.go", "MaxIdleConns")
	suite.checkFileContains(t, "config/database.go", "connection pool")

	// Asset optimization
	suite.checkFileContains(t, "vite.config.js", "minify")

	// Performance testing
	suite.checkFileContains(t, "benchmarks/api_test.go", "b.ResetTimer")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "performance:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "go test -bench")
}

func TestMonolithAcceptance_DevelopmentWorkflow(t *testing.T) {
	// GIVEN: A developer wants efficient development workflow
	// WHEN: They generate a monolith project
	// THEN: The project should provide excellent development experience

	suite := setupAcceptanceTest(t)
	suite.generateMonolithProject(t)

	// Makefile with essential targets
	suite.checkFileExists(t, "Makefile")
	makefileContent, err := os.ReadFile(filepath.Join(suite.projectDir, "Makefile"))
	require.NoError(t, err)
	content := string(makefileContent)

	essentialTargets := []string{"build:", "test:", "dev:", "help:"}
	for _, target := range essentialTargets {
		assert.Contains(t, content, target, "Makefile should include target: %s", target)
	}

	// Development tools
	suite.checkFileExists(t, ".golangci.yml")  // Linting
	suite.checkFileExists(t, "air.toml")       // Hot reload
	suite.checkFileExists(t, "scripts/setup.sh") // Development setup

	// Development setup script
	suite.checkFileContains(t, "scripts/setup.sh", "check_requirements")
	suite.checkFileContains(t, "scripts/setup.sh", "install_go_deps")
	suite.checkFileContains(t, "scripts/setup.sh", "setup_database")

	// Environment configuration
	suite.checkFileExists(t, ".env.example")
}

// Helper function to run acceptance tests in build environment
func TestMonolithAcceptance_FullSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full acceptance test suite in short mode")
	}

	// This meta-test ensures all acceptance tests can run together
	t.Run("BasicGeneration", TestMonolithAcceptance_BasicGeneration)
	t.Run("MultiFrameworkSupport", TestMonolithAcceptance_MultiFrameworkSupport)
	t.Run("DatabaseDriverSupport", TestMonolithAcceptance_DatabaseDriverSupport)
	t.Run("ProductionReadiness", TestMonolithAcceptance_ProductionReadiness)
	t.Run("ModularArchitecture", TestMonolithAcceptance_ModularArchitecture)
	t.Run("ComprehensiveTestCoverage", TestMonolithAcceptance_ComprehensiveTestCoverage)
	t.Run("SecurityHardening", TestMonolithAcceptance_SecurityHardening)
	t.Run("PerformanceOptimization", TestMonolithAcceptance_PerformanceOptimization)
	t.Run("DevelopmentWorkflow", TestMonolithAcceptance_DevelopmentWorkflow)
}