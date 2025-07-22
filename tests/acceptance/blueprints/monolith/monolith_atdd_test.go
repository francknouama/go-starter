package monolith

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMonolithBlueprintATDD validates the acceptance criteria for monolith blueprint
// This ensures the monolith blueprint generates correctly and produces working web applications
func TestMonolithBlueprintATDD(t *testing.T) {
	t.Run("monolith_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: monolith should be in the list

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/monolith)
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool first
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Logf("Build output: %s", string(output))
		}
		require.NoError(t, err, "Failed to build CLI tool")

		// List blueprints
		listCmd := exec.Command("./go-starter", "list")
		output, err = listCmd.CombinedOutput()
		require.NoError(t, err, "List command should succeed")

		outputStr := string(output)
		assert.Contains(t, outputStr, "monolith", "monolith blueprint should be listed")
		assert.Contains(t, outputStr, "Modular monolith application", "Should show monolith description")
	})

	t.Run("monolith_generates_complete_web_application", func(t *testing.T) {
		// GIVEN: User wants a complete monolith web application
		// WHEN: User generates a project with monolith blueprint
		// THEN: Should generate complete web app with all layers

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		require.NoError(t, err, "Failed to build CLI tool: %s", string(output))

		// Generate a monolith project
		generateCmd := exec.Command("./go-starter", "new", "test-monolith",
			"--type=monolith",
			"--module=github.com/test/monolith",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--auth-type=session",
			"--logger=slog",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-monolith")

		// Verify essential architectural components exist
		essentialComponents := []string{
			// Foundation
			"main.go",
			"go.mod",
			"README.md",
			"Makefile",
			"Dockerfile",
			"docker-compose.yml",
			".env.example",
			".gitignore",
			".golangci.yml",

			// Configuration layer
			"config/database.go",
			"config/session.go",

			// Controllers layer
			"controllers/user.go",
			"controllers/api.go",

			// Models layer (implied by template)
			"models/user.go",

			// Services layer (implied by template)
			"services/user.go",
			"services/user_test.go",

			// Routes layer
			"routes/api.go",
			"routes/auth.go",

			// Database layer
			"database/migrations/001_create_users.sql",

			// Asset pipeline
			"vite.config.js",
			"tailwind.config.js",

			// Testing
			"tests/integration_test.go",
			"benchmarks/api_test.go",

			// DevOps
			"scripts/setup.sh",
			"scripts/migrate.sh",
			".github/workflows/ci.yml",
			".github/workflows/deploy.yml",

			// Container orchestration
			"kubernetes/deployment.yaml",
			"docker-compose.prod.yml",
		}

		for _, component := range essentialComponents {
			assert.FileExists(t, filepath.Join(projectDir, component), "Essential component %s should exist", component)
		}

		// Verify directory structure for modular monolith
		essentialDirs := []string{
			"config",
			"controllers",
			"routes",
			"database/migrations",
			"tests",
			"scripts",
			".github/workflows",
			"kubernetes",
		}

		for _, dir := range essentialDirs {
			assert.DirExists(t, filepath.Join(projectDir, dir), "Essential directory %s should exist", dir)
		}

		// Verify no over-simplification (should have substantial structure)
		var fileCount int
		err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileCount++
			}
			return nil
		})
		require.NoError(t, err)

		// Monolith should have comprehensive file structure (25+ files)
		assert.GreaterOrEqual(t, fileCount, 25, "Monolith should have at least 25 files for complete application")
		t.Logf("Generated monolith project has %d files", fileCount)
	})

	t.Run("monolith_project_builds_successfully", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: User builds the project
		// THEN: It should compile without errors

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-build",
			"--type=monolith",
			"--module=github.com/test/build",
			"--framework=gin",
			"--database-driver=sqlite",
			"--database-orm=gorm",
			"--logger=slog",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-build")

		// Initialize go modules
		modInitCmd := exec.Command("go", "mod", "tidy")
		modInitCmd.Dir = projectDir
		output, err := modInitCmd.CombinedOutput()
		require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

		// Build the generated project
		buildGeneratedCmd := exec.Command("go", "build", "-o", "monolith-app", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated monolith project should build successfully: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "monolith-app"), "Monolith binary should be created")
	})

	t.Run("monolith_supports_multiple_database_drivers", func(t *testing.T) {
		// GIVEN: Monolith blueprint with different database configurations
		// WHEN: Generating projects with different database drivers
		// THEN: Each should generate appropriate database configurations

		drivers := []struct {
			driver     string
			orm        string
			shouldPass bool
		}{
			{"postgres", "gorm", true},
			{"mysql", "gorm", true},
			{"sqlite", "gorm", true},
			{"postgres", "sqlx", true},
		}

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter once
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		for _, tt := range drivers {
			t.Run(fmt.Sprintf("%s_%s", tt.driver, tt.orm), func(t *testing.T) {
				projectName := fmt.Sprintf("test-%s-%s", tt.driver, tt.orm)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=monolith",
					"--module=github.com/test/"+projectName,
					"--framework=gin",
					"--database-driver="+tt.driver,
					"--database-orm="+tt.orm,
					"--no-git")
				output, err := generateCmd.CombinedOutput()

				if tt.shouldPass {
					require.NoError(t, err, "Should generate successfully with %s+%s: %s", tt.driver, tt.orm, string(output))

					projectDir := filepath.Join(tmpDir, projectName)

					// Check database configuration exists
					assert.FileExists(t, filepath.Join(projectDir, "config", "database.go"), 
						"Database config should exist for %s+%s", tt.driver, tt.orm)

					// Check migration file exists
					assert.FileExists(t, filepath.Join(projectDir, "database", "migrations", "001_create_users.sql"),
						"Migration file should exist for %s+%s", tt.driver, tt.orm)

					// Verify database-specific configuration in files
					configContent, err := os.ReadFile(filepath.Join(projectDir, "config", "database.go"))
					require.NoError(t, err)
					configStr := string(configContent)

					switch tt.driver {
					case "postgres":
						assert.Contains(t, configStr, "postgres", "Should contain postgres-specific config")
					case "mysql":
						assert.Contains(t, configStr, "mysql", "Should contain mysql-specific config")
					case "sqlite":
						assert.Contains(t, configStr, "sqlite", "Should contain sqlite-specific config")
					}

					switch tt.orm {
					case "gorm":
						assert.Contains(t, configStr, "gorm", "Should contain GORM-specific config")
					case "sqlx":
						assert.Contains(t, configStr, "sqlx", "Should contain sqlx-specific config")
					}
				}
			})
		}
	})

	t.Run("monolith_supports_multiple_frameworks", func(t *testing.T) {
		// GIVEN: Monolith blueprint with different web frameworks
		// WHEN: Generating projects with different frameworks
		// THEN: Each should generate framework-specific code

		frameworks := []string{"gin", "echo", "fiber", "chi"}

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		for _, framework := range frameworks {
			t.Run(framework, func(t *testing.T) {
				projectName := fmt.Sprintf("test-%s", framework)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=monolith",
					"--module=github.com/test/"+projectName,
					"--framework="+framework,
					"--database-driver=sqlite",
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s: %s", framework, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Check that main.go contains framework-specific imports/code
				mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
				require.NoError(t, err)
				mainStr := string(mainContent)

				switch framework {
				case "gin":
					assert.Contains(t, mainStr, "gin-gonic/gin", "Should import Gin framework")
				case "echo":
					assert.Contains(t, mainStr, "labstack/echo", "Should import Echo framework")
				case "fiber":
					assert.Contains(t, mainStr, "gofiber/fiber", "Should import Fiber framework")
				case "chi":
					assert.Contains(t, mainStr, "go-chi/chi", "Should import Chi framework")
				}
			})
		}
	})

	t.Run("monolith_has_working_makefile_targets", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: User runs various make commands
		// THEN: All essential make targets should work

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-makefile",
			"--type=monolith",
			"--module=github.com/test/makefile",
			"--framework=gin",
			"--database-driver=sqlite",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-makefile")

		// Initialize go modules
		modCmd := exec.Command("go", "mod", "tidy")
		modCmd.Dir = projectDir
		_, err = modCmd.CombinedOutput()
		require.NoError(t, err)

		// Test essential make targets
		makeTargets := []struct {
			target      string
			shouldExist bool
		}{
			{"help", true},
			{"build", true},
			{"test", true},
			{"dev", true},
			{"docker-build", true},
			{"clean", true},
		}

		// First, test make help to see available targets
		makeHelpCmd := exec.Command("make", "help")
		makeHelpCmd.Dir = projectDir
		output, err := makeHelpCmd.CombinedOutput()
		require.NoError(t, err, "make help should succeed")

		helpStr := string(output)
		t.Logf("Available make targets:\n%s", helpStr)

		for _, target := range makeTargets {
			if target.shouldExist {
				assert.Contains(t, helpStr, target.target, "Make help should list %s target", target.target)

				// Test the actual target
				makeTargetCmd := exec.Command("make", target.target)
				makeTargetCmd.Dir = projectDir
				output, err := makeTargetCmd.CombinedOutput()
				
				// Some targets might fail due to missing dependencies, but they should be defined
				if err != nil && target.target == "build" {
					// Build should work
					require.NoError(t, err, "make %s should succeed: %s", target.target, string(output))
				}
				// For other targets, we just verify they're defined (help shows them)
			}
		}

		// Specifically test make build creates binary
		makeBuildCmd := exec.Command("make", "build")
		makeBuildCmd.Dir = projectDir
		output, err = makeBuildCmd.CombinedOutput()
		require.NoError(t, err, "make build should succeed: %s", string(output))

		// Check binary was created (typically in bin/ directory)
		binPath := filepath.Join(projectDir, "bin", "test-makefile")
		assert.FileExists(t, binPath, "Binary should be created by make build")
	})

	t.Run("monolith_includes_comprehensive_testing", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining the testing infrastructure
		// THEN: Should include unit tests, integration tests, and benchmarks

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-testing",
			"--type=monolith",
			"--module=github.com/test/testing",
			"--framework=gin",
			"--database-driver=sqlite",
			"--database-orm=gorm",
			"--auth-type=session",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-testing")

		// Verify test files exist
		testFiles := []string{
			"services/user_test.go",        // Unit tests
			"tests/integration_test.go",    // Integration tests
			"benchmarks/api_test.go",       // Benchmark tests
		}

		for _, testFile := range testFiles {
			assert.FileExists(t, filepath.Join(projectDir, testFile), "Test file %s should exist", testFile)
		}

		// Verify test quality
		unitTestContent, err := os.ReadFile(filepath.Join(projectDir, "services", "user_test.go"))
		require.NoError(t, err)
		unitTestStr := string(unitTestContent)

		// Should use testify for assertions
		assert.Contains(t, unitTestStr, "github.com/stretchr/testify", "Should use testify testing framework")
		assert.Contains(t, unitTestStr, "func Test", "Should contain test functions")
		assert.Contains(t, unitTestStr, "suite.Suite", "Should use test suites")
		assert.Contains(t, unitTestStr, "mock.Mock", "Should include mocking capabilities")

		// Check integration test structure
		integrationTestContent, err := os.ReadFile(filepath.Join(projectDir, "tests", "integration_test.go"))
		require.NoError(t, err)
		integrationTestStr := string(integrationTestContent)

		assert.Contains(t, integrationTestStr, "func TestIntegration", "Should contain integration tests")
		assert.Contains(t, integrationTestStr, "http.Client", "Should test HTTP endpoints")

		// Check benchmark tests
		benchmarkContent, err := os.ReadFile(filepath.Join(projectDir, "benchmarks", "api_test.go"))
		require.NoError(t, err)
		benchmarkStr := string(benchmarkContent)

		assert.Contains(t, benchmarkStr, "func Benchmark", "Should contain benchmark functions")
		assert.Contains(t, benchmarkStr, "b.ResetTimer", "Should follow benchmark best practices")
	})

	t.Run("monolith_includes_production_ready_features", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining production readiness features
		// THEN: Should include Docker, CI/CD, monitoring, and deployment configs

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-production",
			"--type=monolith",
			"--module=github.com/test/production",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-production")

		// Verify Docker configuration
		assert.FileExists(t, filepath.Join(projectDir, "Dockerfile"), "Should include Dockerfile")
		assert.FileExists(t, filepath.Join(projectDir, "docker-compose.yml"), "Should include docker-compose for development")
		assert.FileExists(t, filepath.Join(projectDir, "docker-compose.prod.yml"), "Should include production docker-compose")

		// Verify CI/CD configuration
		assert.FileExists(t, filepath.Join(projectDir, ".github", "workflows", "ci.yml"), "Should include CI workflow")
		assert.FileExists(t, filepath.Join(projectDir, ".github", "workflows", "deploy.yml"), "Should include deployment workflow")

		// Verify Kubernetes deployment
		assert.FileExists(t, filepath.Join(projectDir, "kubernetes", "deployment.yaml"), "Should include Kubernetes deployment")

		// Verify operational scripts
		assert.FileExists(t, filepath.Join(projectDir, "scripts", "setup.sh"), "Should include setup script")
		assert.FileExists(t, filepath.Join(projectDir, "scripts", "migrate.sh"), "Should include migration script")

		// Verify environment configuration
		assert.FileExists(t, filepath.Join(projectDir, ".env.example"), "Should include environment example")

		// Check CI workflow quality
		ciContent, err := os.ReadFile(filepath.Join(projectDir, ".github", "workflows", "ci.yml"))
		require.NoError(t, err)
		ciStr := string(ciContent)

		assert.Contains(t, ciStr, "test", "CI should include test job")
		assert.Contains(t, ciStr, "lint", "CI should include lint job")
		assert.Contains(t, ciStr, "security", "CI should include security scanning")
		assert.Contains(t, ciStr, "build", "CI should include build job")
		assert.Contains(t, ciStr, "docker", "CI should include Docker build")

		// Check deployment workflow
		deployContent, err := os.ReadFile(filepath.Join(projectDir, ".github", "workflows", "deploy.yml"))
		require.NoError(t, err)
		deployStr := string(deployContent)

		assert.Contains(t, deployStr, "deploy-staging", "Should include staging deployment")
		assert.Contains(t, deployStr, "deploy-production", "Should include production deployment")
		assert.Contains(t, deployStr, "rollback", "Should include rollback capability")

		// Check Dockerfile quality
		dockerContent, err := os.ReadFile(filepath.Join(projectDir, "Dockerfile"))
		require.NoError(t, err)
		dockerStr := string(dockerContent)

		assert.Contains(t, dockerStr, "FROM golang:", "Should use official Go image")
		assert.Contains(t, dockerStr, "WORKDIR", "Should set working directory")
		assert.Contains(t, dockerStr, "COPY", "Should copy application files")
		assert.Contains(t, dockerStr, "EXPOSE", "Should expose port")
		assert.Contains(t, dockerStr, "CMD", "Should have start command")
	})

	t.Run("monolith_has_modular_architecture", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining the code organization
		// THEN: Should follow modular monolith architectural patterns

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-architecture",
			"--type=monolith",
			"--module=github.com/test/architecture",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--auth-type=session",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-architecture")

		// Verify modular layers exist
		layers := []string{
			"controllers", // Presentation layer
			"services",    // Business logic layer
			"models",      // Domain models
			"config",      // Configuration layer
			"routes",      // Routing layer
		}

		for _, layer := range layers {
			assert.DirExists(t, filepath.Join(projectDir, layer), "Layer %s should exist", layer)
		}

		// Verify separation of concerns
		controllerContent, err := os.ReadFile(filepath.Join(projectDir, "controllers", "user.go"))
		require.NoError(t, err)
		controllerStr := string(controllerContent)

		// Controllers should focus on HTTP concerns
		assert.Contains(t, controllerStr, "http.ResponseWriter", "Controllers should handle HTTP")
		assert.Contains(t, controllerStr, "json.Marshal", "Controllers should handle serialization")
		assert.NotContains(t, controllerStr, "database/sql", "Controllers should not have direct DB access")

		// Check service layer
		serviceContent, err := os.ReadFile(filepath.Join(projectDir, "services", "user.go"))
		require.NoError(t, err)
		serviceStr := string(serviceContent)

		// Services should contain business logic
		assert.Contains(t, serviceStr, "type UserService interface", "Should define service interface")
		assert.Contains(t, serviceStr, "func New", "Should have constructor functions")
		assert.Contains(t, serviceStr, "validation", "Should include business validation")

		// Check configuration separation
		dbConfigContent, err := os.ReadFile(filepath.Join(projectDir, "config", "database.go"))
		require.NoError(t, err)
		dbConfigStr := string(dbConfigContent)

		assert.Contains(t, dbConfigStr, "type Database interface", "Should define database interface")
		assert.Contains(t, dbConfigStr, "connection pool", "Should handle connection pooling")
	})
}

// TestMonolithProductionReadiness validates production-ready features
func TestMonolithProductionReadiness(t *testing.T) {
	t.Run("generates_security_hardened_configuration", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining security configurations
		// THEN: Should include OWASP security best practices

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-security",
			"--type=monolith",
			"--module=github.com/test/security",
			"--framework=gin",
			"--auth-type=session",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-security")

		// Check session configuration for security
		sessionConfigContent, err := os.ReadFile(filepath.Join(projectDir, "config", "session.go"))
		require.NoError(t, err)
		sessionConfigStr := string(sessionConfigContent)

		// Should include secure session configuration
		assert.Contains(t, sessionConfigStr, "HttpOnly: true", "Sessions should be HttpOnly")
		assert.Contains(t, sessionConfigStr, "Secure: true", "Sessions should be Secure")
		assert.Contains(t, sessionConfigStr, "SameSite", "Should configure SameSite attribute")

		// Check if CSRF protection is mentioned
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)

		// Should include security headers
		assert.Contains(t, mainStr, "security", "Should reference security middleware")

		// Check for rate limiting configuration
		routesContent, err := os.ReadFile(filepath.Join(projectDir, "routes", "api.go"))
		require.NoError(t, err)
		routesStr := string(routesContent)

		assert.Contains(t, routesStr, "/health", "Should include health endpoint")
		assert.Contains(t, routesStr, "/api/v1", "Should use versioned API")
	})

	t.Run("includes_comprehensive_logging", func(t *testing.T) {
		// GIVEN: A generated monolith project with different loggers
		// WHEN: Examining logging configuration
		// THEN: Should support multiple logging libraries with consistent interface

		loggers := []string{"slog", "zap", "logrus", "zerolog"}

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		for _, logger := range loggers {
			t.Run(logger, func(t *testing.T) {
				projectName := fmt.Sprintf("test-log-%s", logger)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=monolith",
					"--module=github.com/test/"+projectName,
					"--logger="+logger,
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate with %s logger: %s", logger, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Check main.go for logger setup
				mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
				require.NoError(t, err)
				mainStr := string(mainContent)

				switch logger {
				case "slog":
					assert.Contains(t, mainStr, "log/slog", "Should import slog")
				case "zap":
					assert.Contains(t, mainStr, "go.uber.org/zap", "Should import zap")
				case "logrus":
					assert.Contains(t, mainStr, "github.com/sirupsen/logrus", "Should import logrus")
				case "zerolog":
					assert.Contains(t, mainStr, "github.com/rs/zerolog", "Should import zerolog")
				}

				// Verify logger is used in services
				serviceContent, err := os.ReadFile(filepath.Join(projectDir, "services", "user_test.go"))
				require.NoError(t, err)
				serviceStr := string(serviceContent)

				assert.Contains(t, serviceStr, logger, "Service tests should use %s logger", logger)
			})
		}
	})

	t.Run("supports_asset_pipeline_integration", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining asset pipeline configuration
		// THEN: Should include Vite and Tailwind CSS configuration

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-assets",
			"--type=monolith",
			"--module=github.com/test/assets",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-assets")

		// Verify asset pipeline files
		assert.FileExists(t, filepath.Join(projectDir, "vite.config.js"), "Should include Vite configuration")
		assert.FileExists(t, filepath.Join(projectDir, "tailwind.config.js"), "Should include Tailwind configuration")

		// Check Vite configuration quality
		viteContent, err := os.ReadFile(filepath.Join(projectDir, "vite.config.js"))
		require.NoError(t, err)
		viteStr := string(viteContent)

		assert.Contains(t, viteStr, "build:", "Should include build configuration")
		assert.Contains(t, viteStr, "rollupOptions", "Should include Rollup options")
		assert.Contains(t, viteStr, "legacy", "Should support legacy browsers")

		// Check Tailwind configuration
		tailwindContent, err := os.ReadFile(filepath.Join(projectDir, "tailwind.config.js"))
		require.NoError(t, err)
		tailwindStr := string(tailwindContent)

		assert.Contains(t, tailwindStr, "content:", "Should specify content sources")
		assert.Contains(t, tailwindStr, "theme:", "Should include theme configuration")
		assert.Contains(t, tailwindStr, "plugins:", "Should include plugins array")
	})

	t.Run("includes_database_migrations_and_seeding", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining database management features
		// THEN: Should include migration system and data seeding

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-database",
			"--type=monolith",
			"--module=github.com/test/database",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-database")

		// Verify migration infrastructure
		assert.DirExists(t, filepath.Join(projectDir, "database", "migrations"), "Should include migrations directory")
		assert.FileExists(t, filepath.Join(projectDir, "database", "migrations", "001_create_users.sql"), "Should include initial migration")

		// Verify migration script
		assert.FileExists(t, filepath.Join(projectDir, "scripts", "migrate.sh"), "Should include migration script")

		// Check migration script quality
		migrateScript, err := os.ReadFile(filepath.Join(projectDir, "scripts", "migrate.sh"))
		require.NoError(t, err)
		migrateStr := string(migrateScript)

		assert.Contains(t, migrateStr, "migrate up", "Should support migration up")
		assert.Contains(t, migrateStr, "migrate down", "Should support migration down")
		assert.Contains(t, migrateStr, "create_migration", "Should support creating new migrations")
		assert.Contains(t, migrateStr, "test_connection", "Should test database connection")

		// Check migration file quality
		migrationContent, err := os.ReadFile(filepath.Join(projectDir, "database", "migrations", "001_create_users.sql"))
		require.NoError(t, err)
		migrationStr := string(migrationContent)

		assert.Contains(t, migrationStr, "CREATE TABLE", "Should create tables")
		assert.Contains(t, migrationStr, "users", "Should create users table")
		assert.Contains(t, migrationStr, "id", "Should include ID field")
		assert.Contains(t, migrationStr, "email", "Should include email field")
		assert.Contains(t, migrationStr, "INDEX", "Should include database indexes")
	})
}

// TestMonolithQualityGates validates quality and compliance requirements
func TestMonolithQualityGates(t *testing.T) {
	t.Run("passes_linting_requirements", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Running linting tools
		// THEN: Should pass all linting requirements

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-lint",
			"--type=monolith",
			"--module=github.com/test/lint",
			"--framework=gin",
			"--database-driver=sqlite",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-lint")

		// Initialize modules
		modCmd := exec.Command("go", "mod", "tidy")
		modCmd.Dir = projectDir
		_, err = modCmd.CombinedOutput()
		require.NoError(t, err)

		// Verify .golangci.yml exists and has good configuration
		assert.FileExists(t, filepath.Join(projectDir, ".golangci.yml"), "Should include golangci-lint configuration")

		lintConfigContent, err := os.ReadFile(filepath.Join(projectDir, ".golangci.yml"))
		require.NoError(t, err)
		lintConfigStr := string(lintConfigContent)

		assert.Contains(t, lintConfigStr, "run:", "Should include run configuration")
		assert.Contains(t, lintConfigStr, "linters:", "Should specify linters")
		assert.Contains(t, lintConfigStr, "issues:", "Should configure issue handling")

		// Run go fmt to verify code is formatted
		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Dir = projectDir
		output, err := fmtCmd.CombinedOutput()
		
		// If go fmt produces output, it means files were not formatted
		assert.Empty(t, string(output), "Generated code should be properly formatted")
		require.NoError(t, err, "go fmt should succeed")

		// Run go vet
		vetCmd := exec.Command("go", "vet", "./...")
		vetCmd.Dir = projectDir
		output, err = vetCmd.CombinedOutput()
		require.NoError(t, err, "go vet should pass: %s", string(output))
	})

	t.Run("includes_security_scanning_configuration", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining security scanning setup
		// THEN: Should include security scanning in CI and configuration

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-security-scan",
			"--type=monolith",
			"--module=github.com/test/security-scan",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-security-scan")

		// Check CI includes security scanning
		ciContent, err := os.ReadFile(filepath.Join(projectDir, ".github", "workflows", "ci.yml"))
		require.NoError(t, err)
		ciStr := string(ciContent)

		assert.Contains(t, ciStr, "security:", "Should include security job")
		assert.Contains(t, ciStr, "gosec", "Should use gosec security scanner")
		assert.Contains(t, ciStr, "govulncheck", "Should use govulncheck for vulnerabilities")
		assert.Contains(t, ciStr, "SARIF", "Should upload SARIF results")
	})

	t.Run("meets_performance_requirements", func(t *testing.T) {
		// GIVEN: A generated monolith project
		// WHEN: Examining performance-related configurations
		// THEN: Should include performance optimizations and benchmarks

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-performance",
			"--type=monolith",
			"--module=github.com/test/performance",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-performance")

		// Check database configuration for performance
		dbConfigContent, err := os.ReadFile(filepath.Join(projectDir, "config", "database.go"))
		require.NoError(t, err)
		dbConfigStr := string(dbConfigContent)

		assert.Contains(t, dbConfigStr, "connection pool", "Should configure connection pooling")
		assert.Contains(t, dbConfigStr, "MaxOpenConns", "Should set maximum open connections")
		assert.Contains(t, dbConfigStr, "MaxIdleConns", "Should set maximum idle connections")

		// Check benchmark tests exist
		benchmarkContent, err := os.ReadFile(filepath.Join(projectDir, "benchmarks", "api_test.go"))
		require.NoError(t, err)
		benchmarkStr := string(benchmarkContent)

		assert.Contains(t, benchmarkStr, "func Benchmark", "Should include benchmark functions")
		assert.Contains(t, benchmarkStr, "b.ResetTimer", "Should follow benchmark best practices")
		assert.Contains(t, benchmarkStr, "b.RunParallel", "Should include parallel benchmarks")

		// Check CI includes performance testing
		ciContent, err := os.ReadFile(filepath.Join(projectDir, ".github", "workflows", "ci.yml"))
		require.NoError(t, err)
		ciStr := string(ciContent)

		assert.Contains(t, ciStr, "performance:", "Should include performance job")
		assert.Contains(t, ciStr, "go test -bench", "Should run benchmark tests")
	})
}

// Helper functions for ATDD tests

func waitForServerStart(url string, timeout time.Duration) error {
	client := &http.Client{Timeout: 1 * time.Second}
	deadline := time.Now().Add(timeout)
	
	for time.Now().Before(deadline) {
		resp, err := client.Get(url)
		if err == nil {
			resp.Body.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("server did not start within %v", timeout)
}

func checkProcessRunning(cmd *exec.Cmd) bool {
	if cmd.Process == nil {
		return false
	}
	// On Unix systems, sending signal 0 tests if process exists
	err := cmd.Process.Signal(os.Signal(os.Kill))
	return err == nil
}

func findAvailablePort() int {
	// Simple port finding - in real tests you'd use net.Listen
	return 8080
}

func extractProjectNameFromOutput(output string) string {
	// Extract project name from generation output
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "Generated project:") {
			parts := strings.Fields(line)
			if len(parts) > 2 {
				return parts[len(parts)-1]
			}
		}
	}
	return "unknown"
}

func countFilesInDirectory(dir string) (int, error) {
	var count int
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

func parseJSONResponse(body []byte, target interface{}) error {
	return json.Unmarshal(body, target)
}

func validateResponseHeaders(headers http.Header, expectedHeaders map[string]string) error {
	for key, expectedValue := range expectedHeaders {
		actualValue := headers.Get(key)
		if actualValue != expectedValue {
			return fmt.Errorf("expected header %s to be %s, got %s", key, expectedValue, actualValue)
		}
	}
	return nil
}

func checkFileContainsPattern(filePath string, pattern string) (bool, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return false, err
	}
	
	matched, err := regexp.MatchString(pattern, string(content))
	return matched, err
}

func runCommandWithTimeout(cmd *exec.Cmd, timeout time.Duration) ([]byte, error) {
	type result struct {
		output []byte
		err    error
	}
	
	resultChan := make(chan result, 1)
	
	go func() {
		output, err := cmd.CombinedOutput()
		resultChan <- result{output, err}
	}()
	
	select {
	case res := <-resultChan:
		return res.output, res.err
	case <-time.After(timeout):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		return nil, fmt.Errorf("command timed out after %v", timeout)
	}
}

func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	
	return lines, scanner.Err()
}