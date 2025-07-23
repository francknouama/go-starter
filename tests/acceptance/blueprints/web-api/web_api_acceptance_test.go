package webapi

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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// WebAPIAcceptanceTestSuite provides comprehensive ATDD coverage for web API blueprints
// Tests all architectures: standard, clean, ddd, hexagonal with multi-framework support
type WebAPIAcceptanceTestSuite struct {
	workingDir    string
	projectDir    string
	projectName   string
	originalDir   string
	projectRoot   string
	httpClient    *http.Client
	
	// Configuration
	blueprintType  string
	framework      string
	architecture   string
	databaseDriver string
	databaseORM    string
	authType       string
	logger         string
	
	// Testcontainers
	postgresContainer testcontainers.Container
	mysqlContainer    testcontainers.Container
	databaseURL       string
	ctx               context.Context
}

func setupWebAPIAcceptanceTest(t *testing.T, blueprintType string) *WebAPIAcceptanceTestSuite {
	suite := &WebAPIAcceptanceTestSuite{
		projectName:   "test-web-api-" + blueprintType,
		blueprintType: blueprintType,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
		ctx:           context.Background(),
	}

	var err error
	suite.originalDir, err = os.Getwd()
	require.NoError(t, err)

	suite.projectRoot = filepath.Join(suite.originalDir, "..", "..", "..", "..")
	
	suite.workingDir, err = os.MkdirTemp("", "web-api-acceptance-*")
	require.NoError(t, err)

	err = os.Chdir(suite.workingDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		suite.cleanupContainers()
		_ = os.Chdir(suite.originalDir)
		_ = os.RemoveAll(suite.workingDir)
	})

	return suite
}

func (suite *WebAPIAcceptanceTestSuite) buildCLI(t *testing.T) {
	buildCmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", "go-starter", ".")
	buildCmd.Dir = suite.projectRoot
	output, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "Failed to build go-starter CLI: %s", string(output))
}

func (suite *WebAPIAcceptanceTestSuite) generateWebAPIProject(t *testing.T, args ...string) {
	suite.buildCLI(t)

	// Determine blueprint type based on architecture
	blueprintType := "web-api"
	if suite.architecture != "" && suite.architecture != "standard" {
		blueprintType = "web-api-" + suite.architecture
	}

	baseArgs := []string{
		"new", suite.projectName,
		"--type=" + blueprintType,
		"--module=github.com/test/" + suite.projectName,
		"--framework=" + getDefault(suite.framework, "gin"),
		"--database-driver=" + getDefault(suite.databaseDriver, "postgres"),
		"--database-orm=" + getDefault(suite.databaseORM, "gorm"),
		"--auth-type=" + getDefault(suite.authType, "jwt"),
		"--logger=" + getDefault(suite.logger, "slog"),
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

func (suite *WebAPIAcceptanceTestSuite) checkFileExists(t *testing.T, relativePath string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	assert.FileExists(t, fullPath, "File should exist: %s", relativePath)
}

func (suite *WebAPIAcceptanceTestSuite) checkFileContains(t *testing.T, relativePath, content string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	
	// Handle directory searches
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		found := false
		_ = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") {
				fileContent, err := os.ReadFile(path)
				if err == nil && strings.Contains(string(fileContent), content) {
					found = true
					return filepath.SkipDir
				}
			}
			return nil
		})
		assert.True(t, found, "Content '%s' not found in directory %s", content, relativePath)
		return
	}

	fileContent, err := os.ReadFile(fullPath)
	require.NoError(t, err, "Should be able to read file: %s", relativePath)
	assert.Contains(t, string(fileContent), content, "File %s should contain '%s'", relativePath, content)
}

func (suite *WebAPIAcceptanceTestSuite) checkFileDoesNotContain(t *testing.T, relativePath, content string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	
	// Handle directory searches
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		_ = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") {
				fileContent, err := os.ReadFile(path)
				if err == nil {
					assert.NotContains(t, string(fileContent), content, 
						"File %s should not contain '%s'", path, content)
				}
			}
			return nil
		})
		return
	}

	fileContent, err := os.ReadFile(fullPath)
	if err != nil {
		return // File doesn't exist, so it doesn't contain the content
	}
	assert.NotContains(t, string(fileContent), content, "File %s should not contain '%s'", relativePath, content)
}

func (suite *WebAPIAcceptanceTestSuite) compileProject(t *testing.T) {
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = suite.projectDir
	output, err := modCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

	// Build the project
	buildCmd := exec.Command("go", "build", "-o", "web-api-app", ".")
	buildCmd.Dir = suite.projectDir
	output, err = buildCmd.CombinedOutput()
	require.NoError(t, err, "Project should compile successfully: %s", string(output))
}

func (suite *WebAPIAcceptanceTestSuite) cleanupContainers() {
	if suite.postgresContainer != nil {
		_ = suite.postgresContainer.Terminate(suite.ctx)
	}
	if suite.mysqlContainer != nil {
		_ = suite.mysqlContainer.Terminate(suite.ctx)
	}
}

// ATDD Test Scenarios for Web API Blueprints

func TestWebAPIAcceptance_StandardArchitecture(t *testing.T) {
	// GIVEN: A developer wants to create a standard web API
	// WHEN: They generate a web API project with standard architecture
	// THEN: A complete, working web API application is generated

	suite := setupWebAPIAcceptanceTest(t, "standard")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// Verify essential web API components
	essentialFiles := []string{
		"main.go", "go.mod", "README.md", "Makefile", "Dockerfile",
		"internal/handlers", "internal/middleware", "internal/services",
		"internal/models", "internal/repository", "internal/config",
		"api/openapi.yaml", "api/docs",
		"tests/integration", "tests/unit",
		".github/workflows/ci.yml", ".github/workflows/deploy.yml",
		"configs/config.yaml", "configs/config.prod.yaml",
	}

	for _, file := range essentialFiles {
		suite.checkFileExists(t, file)
	}

	// Verify standard architecture structure
	suite.checkFileExists(t, "internal/handlers/user.go")
	suite.checkFileExists(t, "internal/middleware/auth.go")
	suite.checkFileExists(t, "internal/services/user.go")
	suite.checkFileExists(t, "internal/repository/user.go")

	// Verify the project compiles
	suite.compileProject(t)
}

func TestWebAPIAcceptance_CleanArchitecture(t *testing.T) {
	// GIVEN: A developer wants clean architecture patterns
	// WHEN: They generate a web API with clean architecture
	// THEN: The project follows clean architecture principles

	suite := setupWebAPIAcceptanceTest(t, "clean")
	suite.architecture = "clean"
	suite.generateWebAPIProject(t)

	// Verify clean architecture structure
	cleanArchitectureDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases", 
		"internal/domain/ports",
		"internal/adapters/controllers",
		"internal/adapters/presenters",
		"internal/infrastructure/persistence",
		"internal/infrastructure/web",
		"internal/infrastructure/logger",
		"internal/infrastructure/config",
		"internal/infrastructure/container",
	}

	for _, dir := range cleanArchitectureDirs {
		suite.checkFileExists(t, dir)
	}

	// Verify dependency inversion - domain shouldn't depend on infrastructure
	suite.checkFileDoesNotContain(t, "internal/domain", "internal/infrastructure")
	suite.checkFileDoesNotContain(t, "internal/domain", "gin-gonic")
	suite.checkFileDoesNotContain(t, "internal/domain", "gorm.io")

	// Verify interfaces are defined in ports
	suite.checkFileContains(t, "internal/domain/ports", "interface")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_DDDArchitecture(t *testing.T) {
	// GIVEN: A developer wants domain-driven design patterns
	// WHEN: They generate a web API with DDD architecture  
	// THEN: The project follows DDD principles

	suite := setupWebAPIAcceptanceTest(t, "ddd")
	suite.architecture = "ddd"
	suite.generateWebAPIProject(t)

	// Verify DDD structure
	dddDirs := []string{
		"internal/domain/aggregates",
		"internal/domain/entities", 
		"internal/domain/valueobjects",
		"internal/domain/services",
		"internal/domain/events",
		"internal/application/commands",
		"internal/application/queries",
		"internal/application/handlers",
		"internal/infrastructure/persistence",
		"internal/infrastructure/messaging",
	}

	for _, dir := range dddDirs {
		suite.checkFileExists(t, dir)
	}

	// Verify DDD patterns
	suite.checkFileContains(t, "internal/domain/entities", "Entity")
	suite.checkFileContains(t, "internal/domain/valueobjects", "ValueObject") 
	suite.checkFileContains(t, "internal/domain/events", "DomainEvent")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_HexagonalArchitecture(t *testing.T) {
	// GIVEN: A developer wants hexagonal architecture patterns
	// WHEN: They generate a web API with hexagonal architecture
	// THEN: The project follows ports and adapters pattern

	suite := setupWebAPIAcceptanceTest(t, "hexagonal")
	suite.architecture = "hexagonal"
	suite.generateWebAPIProject(t)

	// Verify hexagonal structure
	hexagonalDirs := []string{
		"internal/application/core",
		"internal/application/ports/input",
		"internal/application/ports/output", 
		"internal/adapters/primary/http",
		"internal/adapters/primary/grpc",
		"internal/adapters/secondary/database",
		"internal/adapters/secondary/messaging",
		"internal/adapters/secondary/cache",
	}

	for _, dir := range hexagonalDirs {
		suite.checkFileExists(t, dir)
	}

	// Verify ports and adapters pattern
	suite.checkFileContains(t, "internal/application/ports/input", "interface")
	suite.checkFileContains(t, "internal/application/ports/output", "interface")
	suite.checkFileContains(t, "internal/adapters/primary", "Port")
	suite.checkFileContains(t, "internal/adapters/secondary", "Port")

	// Application core should be isolated
	suite.checkFileDoesNotContain(t, "internal/application/core", "gin-gonic")
	suite.checkFileDoesNotContain(t, "internal/application/core", "gorm.io")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_MultiFrameworkSupport(t *testing.T) {
	// GIVEN: A developer wants to use different web frameworks
	// WHEN: They generate web API projects with different frameworks
	// THEN: Each project should use the appropriate framework

	frameworks := []struct {
		name        string
		importPath  string
		handlerType string
	}{
		{"gin", "gin-gonic/gin", "*gin.Context"},
		{"echo", "labstack/echo", "echo.Context"},
		{"fiber", "gofiber/fiber", "*fiber.Ctx"},
		{"chi", "go-chi/chi", "http.ResponseWriter"},
	}

	for _, fw := range frameworks {
		t.Run(fw.name, func(t *testing.T) {
			suite := setupWebAPIAcceptanceTest(t, fw.name+"-standard")
			suite.framework = fw.name
			suite.architecture = "standard"
			
			suite.generateWebAPIProject(t)
			
			// Verify framework-specific imports and patterns
			suite.checkFileContains(t, "main.go", fw.importPath)
			suite.checkFileContains(t, "internal/handlers", fw.handlerType)
			
			// Verify middleware is framework-compatible
			suite.checkFileExists(t, "internal/middleware")
			
			suite.compileProject(t)
		})
	}
}

func TestWebAPIAcceptance_DatabaseIntegration(t *testing.T) {
	// GIVEN: A developer wants database integration
	// WHEN: They generate web API with different databases and ORMs
	// THEN: Each configuration should work correctly

	if testing.Short() {
		t.Skip("Skipping database integration tests in short mode")
	}

	databases := []struct {
		driver string
		orm    string
	}{
		{"postgres", "gorm"},
		{"postgres", "sqlx"}, 
		{"mysql", "gorm"},
		{"sqlite", "gorm"},
	}

	for _, db := range databases {
		t.Run(fmt.Sprintf("%s_%s", db.driver, db.orm), func(t *testing.T) {
			suite := setupWebAPIAcceptanceTest(t, fmt.Sprintf("db-%s-%s", db.driver, db.orm))
			suite.databaseDriver = db.driver
			suite.databaseORM = db.orm
			suite.architecture = "standard"
			
			suite.generateWebAPIProject(t)
			
			// Verify database configuration
			suite.checkFileExists(t, "internal/database")
			suite.checkFileContains(t, "internal/repository", db.orm)
			
			// Check migration files
			suite.checkFileExists(t, "migrations")
			
			// Test with real database container
			if db.driver == "postgres" {
				err := suite.setupPostgresContainer(t)
				require.NoError(t, err)
			} else if db.driver == "mysql" {
				err := suite.setupMySQLContainer(t)
				require.NoError(t, err)
			}
			
			suite.compileProject(t)
		})
	}
}

func TestWebAPIAcceptance_AuthenticationTypes(t *testing.T) {
	// GIVEN: A developer wants different authentication mechanisms
	// WHEN: They generate web API with different auth types
	// THEN: Each auth type should be properly configured

	authTypes := []string{"jwt", "session", "oauth2", "api-key"}

	for _, authType := range authTypes {
		t.Run(authType, func(t *testing.T) {
			suite := setupWebAPIAcceptanceTest(t, "auth-"+authType)
			suite.authType = authType
			suite.architecture = "standard"
			
			suite.generateWebAPIProject(t)
			
			// Verify authentication configuration
			suite.checkFileExists(t, "internal/middleware/auth.go")
			suite.checkFileContains(t, "internal/handlers", "auth")
			
			if authType == "jwt" {
				suite.checkFileContains(t, "internal/services", "jwt")
			} else if authType == "session" {
				suite.checkFileExists(t, "internal/session")
			}
			
			suite.compileProject(t)
		})
	}
}

func TestWebAPIAcceptance_LoggerIntegration(t *testing.T) {
	// GIVEN: A developer wants different logging libraries
	// WHEN: They generate web API with different loggers
	// THEN: Each logger should be properly integrated

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
			suite := setupWebAPIAcceptanceTest(t, "logger-"+logger.name)
			suite.logger = logger.name
			suite.architecture = "standard"
			
			suite.generateWebAPIProject(t)
			
			// Verify logger integration
			suite.checkFileContains(t, "main.go", logger.importPath)
			suite.checkFileExists(t, "internal/logger")
			
			// Check structured logging in handlers
			suite.checkFileContains(t, "internal/handlers", "logger")
			
			suite.compileProject(t)
		})
	}
}

func TestWebAPIAcceptance_ProductionReadiness(t *testing.T) {
	// GIVEN: A developer wants production-ready web API
	// WHEN: They generate a web API project
	// THEN: All production features should be included

	suite := setupWebAPIAcceptanceTest(t, "production")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// Docker and containerization
	suite.checkFileExists(t, "Dockerfile")
	suite.checkFileExists(t, "docker-compose.yml")
	suite.checkFileExists(t, "docker-compose.prod.yml")
	suite.checkFileContains(t, "Dockerfile", "FROM")
	suite.checkFileContains(t, "Dockerfile", "HEALTHCHECK")

	// CI/CD pipelines
	suite.checkFileExists(t, ".github/workflows/ci.yml")
	suite.checkFileExists(t, ".github/workflows/deploy.yml")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "test:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "security:")
	suite.checkFileContains(t, ".github/workflows/deploy.yml", "deploy:")

	// API documentation
	suite.checkFileExists(t, "api/openapi.yaml")
	suite.checkFileContains(t, "api/openapi.yaml", "openapi: 3.0")
	suite.checkFileExists(t, "api/docs")

	// Monitoring and observability
	suite.checkFileContains(t, "internal/handlers", "/health")
	suite.checkFileContains(t, "internal/handlers", "/metrics")
	suite.checkFileExists(t, "internal/middleware/metrics.go")

	// Security
	suite.checkFileExists(t, "internal/middleware/cors.go")
	suite.checkFileExists(t, "internal/middleware/security_headers.go")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "gosec")

	// Configuration management
	suite.checkFileExists(t, "configs")
	suite.checkFileContains(t, "configs", "dev")
	suite.checkFileContains(t, "configs", "prod")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_APIDocumentation(t *testing.T) {
	// GIVEN: A developer wants comprehensive API documentation
	// WHEN: They generate a web API project
	// THEN: OpenAPI documentation should be complete

	suite := setupWebAPIAcceptanceTest(t, "docs")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// OpenAPI specification
	suite.checkFileExists(t, "api/openapi.yaml")
	suite.checkFileContains(t, "api/openapi.yaml", "openapi: 3.0")
	suite.checkFileContains(t, "api/openapi.yaml", "info:")
	suite.checkFileContains(t, "api/openapi.yaml", "paths:")
	suite.checkFileContains(t, "api/openapi.yaml", "components:")
	suite.checkFileContains(t, "api/openapi.yaml", "schemas:")

	// API documentation generation
	suite.checkFileExists(t, "api/docs")
	suite.checkFileContains(t, "internal/handlers", "swagger")

	// RESTful endpoint patterns
	suite.checkFileContains(t, "api/openapi.yaml", "/api/v1")
	suite.checkFileContains(t, "api/openapi.yaml", "get:")
	suite.checkFileContains(t, "api/openapi.yaml", "post:")
	suite.checkFileContains(t, "api/openapi.yaml", "put:")
	suite.checkFileContains(t, "api/openapi.yaml", "delete:")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_SecurityHardening(t *testing.T) {
	// GIVEN: A developer wants a secure web API
	// WHEN: They generate a web API with security features
	// THEN: Security best practices should be implemented

	suite := setupWebAPIAcceptanceTest(t, "security")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// Security middleware
	suite.checkFileExists(t, "internal/middleware/cors.go")
	suite.checkFileExists(t, "internal/middleware/security_headers.go")
	suite.checkFileExists(t, "internal/middleware/rate_limiter.go")
	suite.checkFileExists(t, "internal/middleware/validation.go")

	// Security headers
	suite.checkFileContains(t, "internal/middleware/security_headers.go", "X-Content-Type-Options")
	suite.checkFileContains(t, "internal/middleware/security_headers.go", "X-Frame-Options")
	suite.checkFileContains(t, "internal/middleware/security_headers.go", "X-XSS-Protection")

	// Input validation
	suite.checkFileContains(t, "internal/middleware/validation.go", "validate")
	suite.checkFileContains(t, "internal/handlers", "validation")

	// Security scanning in CI
	suite.checkFileContains(t, ".github/workflows/ci.yml", "gosec")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "govulncheck")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "SARIF")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_TestingInfrastructure(t *testing.T) {
	// GIVEN: A developer wants comprehensive testing
	// WHEN: They generate a web API project
	// THEN: Complete testing infrastructure should be included

	suite := setupWebAPIAcceptanceTest(t, "testing")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// Test directories and files
	testFiles := []string{
		"tests/unit",
		"tests/integration", 
		"tests/e2e",
		"tests/benchmarks",
	}

	for _, testFile := range testFiles {
		suite.checkFileExists(t, testFile)
	}

	// Unit tests
	suite.checkFileContains(t, "tests/unit", "func Test")
	suite.checkFileContains(t, "tests/unit", "testify")

	// Integration tests with testcontainers
	suite.checkFileContains(t, "tests/integration", "testcontainers")
	suite.checkFileContains(t, "tests/integration", "TestIntegration")

	// Benchmark tests  
	suite.checkFileContains(t, "tests/benchmarks", "func Benchmark")
	suite.checkFileContains(t, "tests/benchmarks", "b.ResetTimer")

	// Test coverage in CI
	suite.checkFileContains(t, ".github/workflows/ci.yml", "coverage")
	suite.checkFileContains(t, "Makefile", "test-coverage")

	suite.compileProject(t)
}

func TestWebAPIAcceptance_PerformanceOptimization(t *testing.T) {
	// GIVEN: A developer wants optimized performance
	// WHEN: They generate a web API project
	// THEN: Performance optimizations should be included

	suite := setupWebAPIAcceptanceTest(t, "performance")
	suite.architecture = "standard"
	suite.generateWebAPIProject(t)

	// Database connection pooling
	suite.checkFileContains(t, "internal/database", "MaxOpenConns")
	suite.checkFileContains(t, "internal/database", "MaxIdleConns")
	suite.checkFileContains(t, "internal/database", "ConnMaxLifetime")

	// Caching middleware
	suite.checkFileExists(t, "internal/middleware/cache.go")
	suite.checkFileContains(t, "internal/middleware/cache.go", "cache")

	// Performance monitoring
	suite.checkFileExists(t, "internal/middleware/metrics.go")
	suite.checkFileContains(t, "internal/handlers", "/metrics")

	// Benchmark tests
	suite.checkFileContains(t, "tests/benchmarks", "b.RunParallel")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "go test -bench")

	suite.compileProject(t)
}

// Helper methods

func (suite *WebAPIAcceptanceTestSuite) setupPostgresContainer(t *testing.T) error {
	postgresReq := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "testdb",
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(60 * time.Second),
	}

	var err error
	suite.postgresContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: postgresReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start postgres container: %w", err)
	}

	// Get connection details
	host, err := suite.postgresContainer.Host(suite.ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := suite.postgresContainer.MappedPort(suite.ctx, "5432")
	if err != nil {
		return fmt.Errorf("failed to get container port: %w", err)
	}

	suite.databaseURL = fmt.Sprintf("postgres://testuser:testpass@%s:%s/testdb?sslmode=disable", host, port.Port())
	return nil
}

func (suite *WebAPIAcceptanceTestSuite) setupMySQLContainer(t *testing.T) error {
	mysqlReq := testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "rootpass",
			"MYSQL_DATABASE":      "testdb",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForLog("ready for connections").
			WithStartupTimeout(120 * time.Second),
	}

	var err error
	suite.mysqlContainer, err = testcontainers.GenericContainer(suite.ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: mysqlReq,
		Started:          true,
	})
	if err != nil {
		return fmt.Errorf("failed to start mysql container: %w", err)
	}

	// Get connection details
	host, err := suite.mysqlContainer.Host(suite.ctx)
	if err != nil {
		return fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := suite.mysqlContainer.MappedPort(suite.ctx, "3306")
	if err != nil {
		return fmt.Errorf("failed to get container port: %w", err)
	}

	suite.databaseURL = fmt.Sprintf("testuser:testpass@tcp(%s:%s)/testdb", host, port.Port())
	return nil
}

// Helper function to get default values
func getDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Meta-test to ensure all acceptance tests can run together
func TestWebAPIAcceptance_FullSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full web API acceptance test suite in short mode")
	}

	// This meta-test ensures all acceptance tests can run together  
	t.Run("StandardArchitecture", TestWebAPIAcceptance_StandardArchitecture)
	t.Run("CleanArchitecture", TestWebAPIAcceptance_CleanArchitecture)
	t.Run("DDDArchitecture", TestWebAPIAcceptance_DDDArchitecture)
	t.Run("HexagonalArchitecture", TestWebAPIAcceptance_HexagonalArchitecture)
	t.Run("ProductionReadiness", TestWebAPIAcceptance_ProductionReadiness)
	t.Run("SecurityHardening", TestWebAPIAcceptance_SecurityHardening)
	t.Run("TestingInfrastructure", TestWebAPIAcceptance_TestingInfrastructure)
	t.Run("PerformanceOptimization", TestWebAPIAcceptance_PerformanceOptimization)
}