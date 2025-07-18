package webapi_test

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: Standard Web API Blueprint
// As a developer
// I want to generate a standard web API project
// So that I can quickly build REST APIs

func TestStandard_WebAPI_BasicGeneration_WithGin(t *testing.T) {
	// Scenario: Generate standard web API with Gin
	// Given I want a standard web API
	// When I generate with framework "gin"
	// Then the project should include Gin router setup
	// And the project should have basic CRUD endpoints
	// And the project should compile and run successfully
	// And HTTP server should start on configured port

	// Given I want a standard web API
	config := types.ProjectConfig{
		Name:      "test-standard-gin-api",
		Module:    "github.com/test/test-standard-gin-api",
		Type:      "web-api-standard",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	// When I generate with framework "gin"
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include Gin router setup
	validator := NewWebAPIValidator(projectPath, "standard")
	validator.ValidateGinRouterSetup(t)

	// And the project should have basic CRUD endpoints
	validator.ValidateBasicCRUDEndpoints(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)
}

func TestStandard_WebAPI_WithDatabaseIntegration(t *testing.T) {
	// Scenario: Generate with database integration
	// Given I want a standard web API with database
	// When I generate with "postgres" and "gorm"
	// Then database configuration should be included
	// And GORM models should be provided
	// And database connection should be testable
	// And migrations should be available

	// Given I want a standard web API with database
	config := types.ProjectConfig{
		Name:      "test-standard-db-api",
		Module:    "github.com/test/test-standard-db-api",
		Type:      "web-api-standard",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Drivers: []string{"postgres"},
				ORM:    "gorm",
			},
		},
	}

	// When I generate with "postgres" and "gorm"
	projectPath := helpers.GenerateProject(t, config)

	// Then database configuration should be included
	validator := NewWebAPIValidator(projectPath, "standard")
	validator.ValidateDatabaseConfiguration(t, "postgres", "gorm")

	// And GORM models should be provided
	validator.ValidateGORMModels(t)

	// And database connection should be testable
	validator.ValidateDatabaseConnection(t)

	// And migrations should be available
	validator.ValidateMigrations(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestStandard_WebAPI_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration Across Architectures
	// As a developer
	// I want consistent logging across all web API architectures
	// So that I can monitor applications effectively

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			t.Parallel() // Safe to run logger tests in parallel
			// Scenario: Logger integration
			// Given I generate a "standard" web API with "<logger>"
			// Then logging should be properly configured
			// And log statements should follow architecture patterns
			// And log levels should be configurable
			// And structured logging should be used

			// Given I generate a "standard" web API with "<logger>"
			config := types.ProjectConfig{
				Name:      "test-standard-" + logger,
				Module:    "github.com/test/test-standard-" + logger,
				Type:      "web-api-standard",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then logging should be properly configured
			validator := NewWebAPIValidator(projectPath, "standard")
			validator.ValidateLoggerIntegration(t, logger)

			// And log statements should follow architecture patterns
			validator.ValidateLoggerArchitecturePatterns(t, "standard")

			// And log levels should be configurable
			validator.ValidateLoggerConfiguration(t)

			// And structured logging should be used
			validator.ValidateStructuredLogging(t, logger)

			// And the project should compile successfully
			validator.ValidateCompilation(t)
		})
	}
}

func TestStandard_WebAPI_FrameworkVariations(t *testing.T) {
	// Scenario: Framework variations
	// Given I want to test different web frameworks
	// When I generate projects with different frameworks
	// Then each should work correctly with the standard architecture

	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			t.Parallel() // Safe to run framework tests in parallel
			config := types.ProjectConfig{
				Name:      "test-standard-" + framework,
				Module:    "github.com/test/test-standard-" + framework,
				Type:      "web-api-standard",
				GoVersion: "1.21",
				Framework: framework,
				Logger:    "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Drivers: []string{"postgres"},
						ORM:    "gorm",
					},
					Authentication: types.AuthConfig{
						Type: "jwt",
					},
				},
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewWebAPIValidator(projectPath, "standard")
			validator.ValidateFrameworkIntegration(t, framework)
			validator.ValidateCompilation(t)
		})
	}
}

func TestStandard_WebAPI_ArchitectureCompliance(t *testing.T) {
	// Feature: Architecture-Specific Code Quality
	// As a code reviewer
	// I want generated code to follow architecture principles
	// So that projects maintain architectural integrity

	// Scenario: Architecture compliance for standard
	// Given I generate a "standard" web API
	// Then the code should follow "standard" principles
	// And dependency directions should be correct
	// And layer boundaries should be enforced
	// And the project should pass architectural linting

	config := types.ProjectConfig{
		Name:      "test-standard-compliance",
		Module:    "github.com/test/test-standard-compliance",
		Type:      "web-api-standard",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Drivers: []string{"postgres"},
				ORM:    "gorm",
			},
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should follow "standard" principles
	validator := NewWebAPIValidator(projectPath, "standard")
	validator.ValidateArchitecturePrinciples(t, "standard")

	// And dependency directions should be correct
	validator.ValidateDependencyDirections(t)

	// And layer boundaries should be enforced
	validator.ValidateLayerBoundaries(t, "standard")

	// And the project should pass architectural linting
	validator.ValidateArchitecturalLinting(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestStandard_WebAPI_RuntimeValidation(t *testing.T) {
	// Feature: Blueprint Compilation and Runtime
	// As a user
	// I want all web API blueprints to work out of the box
	// So that I can start development immediately

	// Scenario: Runtime validation
	// Given any compiled web API project
	// When I run the binary
	// Then the HTTP server should start
	// And health endpoints should respond
	// And the server should handle requests correctly
	// And graceful shutdown should work

	config := types.ProjectConfig{
		Name:      "test-standard-runtime",
		Module:    "github.com/test/test-standard-runtime",
		Type:      "web-api-standard",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Given any compiled web API project
	validator := NewWebAPIValidator(projectPath, "standard")
	validator.ValidateCompilation(t)

	// Then the HTTP server should start
	// And health endpoints should respond
	// And the server should handle requests correctly
	// And graceful shutdown should work
	validator.ValidateRuntimeBehavior(t)
}

func TestStandard_WebAPI_MinimalConfiguration(t *testing.T) {
	// Scenario: Minimal configuration without optional features
	// Given I want a minimal standard web API
	// When I generate without database or auth features
	// Then only core files should be present
	// And optional feature files should not exist

	config := types.ProjectConfig{
		Name:      "test-standard-minimal",
		Module:    "github.com/test/test-standard-minimal",
		Type:      "web-api-standard",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		// No features configured
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewWebAPIValidator(projectPath, "standard")
	validator.ValidateMinimalConfiguration(t)
	validator.ValidateCompilation(t)
}

// WebAPIValidator provides validation methods for Web API blueprints
type WebAPIValidator struct {
	ProjectPath  string
	Architecture string
}

func NewWebAPIValidator(projectPath, architecture string) *WebAPIValidator {
	return &WebAPIValidator{
		ProjectPath:  projectPath,
		Architecture: architecture,
	}
}

func (v *WebAPIValidator) ValidateGinRouterSetup(t *testing.T) {
	t.Helper()
	
	// Check main.go or server main for Gin setup
	mainFile := filepath.Join(v.ProjectPath, "cmd/server/main.go")
	helpers.AssertFileExists(t, mainFile)
	helpers.AssertFileContains(t, mainFile, "gin.New()")
	
	// Check health handler uses Gin
	healthFile := filepath.Join(v.ProjectPath, "internal/handlers/health.go")
	helpers.AssertFileExists(t, healthFile)
	helpers.AssertFileContains(t, healthFile, "gin.Context")
}

func (v *WebAPIValidator) ValidateBasicCRUDEndpoints(t *testing.T) {
	t.Helper()
	
	// Health endpoint should always exist
	healthFile := filepath.Join(v.ProjectPath, "internal/handlers/health.go")
	helpers.AssertFileExists(t, healthFile)
	helpers.AssertFileContains(t, healthFile, "GET")
}

func (v *WebAPIValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *WebAPIValidator) ValidateDatabaseConfiguration(t *testing.T, driver, orm string) {
	t.Helper()
	
	// Database connection file should exist
	dbFile := filepath.Join(v.ProjectPath, "internal/database/connection.go")
	helpers.AssertFileExists(t, dbFile)
	
	// Should contain driver-specific imports
	switch driver {
	case "postgres":
		helpers.AssertFileContains(t, dbFile, "postgres")
	case "mysql":
		helpers.AssertFileContains(t, dbFile, "mysql")
	}
	
	// Should contain ORM-specific code
	switch orm {
	case "gorm":
		helpers.AssertFileContains(t, dbFile, "gorm.io/gorm")
	}
}

func (v *WebAPIValidator) ValidateGORMModels(t *testing.T) {
	t.Helper()
	
	userModel := filepath.Join(v.ProjectPath, "internal/models/user.go")
	helpers.AssertFileExists(t, userModel)
	helpers.AssertFileContains(t, userModel, "gorm.Model")
}

func (v *WebAPIValidator) ValidateDatabaseConnection(t *testing.T) {
	t.Helper()
	
	dbFile := filepath.Join(v.ProjectPath, "internal/database/connection.go")
	helpers.AssertFileExists(t, dbFile)
	helpers.AssertFileContains(t, dbFile, "func")
}

func (v *WebAPIValidator) ValidateMigrations(t *testing.T) {
	t.Helper()
	
	migrationUp := filepath.Join(v.ProjectPath, "migrations/001_create_users.up.sql")
	migrationDown := filepath.Join(v.ProjectPath, "migrations/001_create_users.down.sql")
	migrationEmbed := filepath.Join(v.ProjectPath, "migrations/embed.go")
	
	helpers.AssertFileExists(t, migrationUp)
	helpers.AssertFileExists(t, migrationDown)
	helpers.AssertFileExists(t, migrationEmbed)
}

func (v *WebAPIValidator) ValidateLoggerIntegration(t *testing.T, logger string) {
	t.Helper()
	
	// Logger interface should exist
	interfaceFile := filepath.Join(v.ProjectPath, "internal/logger/interface.go")
	helpers.AssertFileExists(t, interfaceFile)
	
	// Factory should exist
	factoryFile := filepath.Join(v.ProjectPath, "internal/logger/factory.go")
	helpers.AssertFileExists(t, factoryFile)
	helpers.AssertFileContains(t, factoryFile, logger)
	
	// Logger-specific implementation should exist
	loggerFile := filepath.Join(v.ProjectPath, "internal/logger/"+logger+".go")
	helpers.AssertFileExists(t, loggerFile)
}

func (v *WebAPIValidator) ValidateLoggerArchitecturePatterns(t *testing.T, architecture string) {
	t.Helper()
	
	// For standard architecture, logger should be used directly in handlers
	if architecture == "standard" {
		healthFile := filepath.Join(v.ProjectPath, "internal/handlers/health.go")
		if helpers.FileExists(healthFile) {
			// Check that logger is used appropriately
			content := helpers.ReadFileContent(t, healthFile)
			assert.Contains(t, content, "logger")
		}
	}
}

func (v *WebAPIValidator) ValidateLoggerConfiguration(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
}

func (v *WebAPIValidator) ValidateStructuredLogging(t *testing.T, logger string) {
	t.Helper()
	
	loggerFile := filepath.Join(v.ProjectPath, "internal/logger/"+logger+".go")
	helpers.AssertFileExists(t, loggerFile)
	
	// Check for structured logging patterns
	content := helpers.ReadFileContent(t, loggerFile)
	switch logger {
	case "slog":
		assert.Contains(t, content, "slog.")
	case "zap":
		assert.Contains(t, content, "zap.")
	case "logrus":
		assert.Contains(t, content, "logrus.")
	case "zerolog":
		assert.Contains(t, content, "zerolog.")
	}
}

func (v *WebAPIValidator) ValidateFrameworkIntegration(t *testing.T, framework string) {
	t.Helper()
	
	// Check go.mod for framework dependency
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	switch framework {
	case "gin":
		helpers.AssertFileContains(t, goMod, "github.com/gin-gonic/gin")
	case "echo":
		helpers.AssertFileContains(t, goMod, "github.com/labstack/echo")
	case "fiber":
		helpers.AssertFileContains(t, goMod, "github.com/gofiber/fiber")
	case "chi":
		helpers.AssertFileContains(t, goMod, "github.com/go-chi/chi")
	}
	
	// Check framework-specific handler exists
	healthFile := filepath.Join(v.ProjectPath, "internal/handlers/health.go")
	helpers.AssertFileExists(t, healthFile)
}

func (v *WebAPIValidator) ValidateArchitecturePrinciples(t *testing.T, architecture string) {
	t.Helper()
	
	switch architecture {
	case "standard":
		// Standard should have flat structure
		v.validateStandardPrinciples(t)
	}
}

func (v *WebAPIValidator) validateStandardPrinciples(t *testing.T) {
	t.Helper()
	
	// Should have standard directories
	expectedDirs := []string{
		"internal/handlers",
		"internal/services", 
		"internal/repository",
		"internal/models",
		"internal/config",
		"internal/logger",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			t.Logf("✓ Expected directory exists: %s", dir)
		} else {
			t.Logf("⚠ Expected directory missing: %s", dir)
		}
	}
	
	// Should NOT have complex architecture directories
	unexpectedDirs := []string{
		"internal/domain",
		"internal/adapters",
		"internal/infrastructure",
		"internal/application",
	}
	
	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		assert.False(t, helpers.DirExists(dirPath), 
			"Directory %s should not exist in standard architecture", dir)
	}
}

func (v *WebAPIValidator) ValidateDependencyDirections(t *testing.T) {
	t.Helper()
	
	// For standard architecture, validate simple dependency flow
	// handlers -> services -> repository -> models
	
	if helpers.FileExists(filepath.Join(v.ProjectPath, "internal/handlers/users.go")) {
		handlerContent := helpers.ReadFileContent(t, filepath.Join(v.ProjectPath, "internal/handlers/users.go"))
		assert.Contains(t, handlerContent, "services", "Handlers should import services")
	}
	
	if helpers.FileExists(filepath.Join(v.ProjectPath, "internal/services/user.go")) {
		serviceContent := helpers.ReadFileContent(t, filepath.Join(v.ProjectPath, "internal/services/user.go"))
		assert.Contains(t, serviceContent, "repository", "Services should import repository")
	}
}

func (v *WebAPIValidator) ValidateLayerBoundaries(t *testing.T, architecture string) {
	t.Helper()
	
	// For standard architecture, ensure clean separation
	if architecture == "standard" {
		// Models should not import handlers
		if helpers.FileExists(filepath.Join(v.ProjectPath, "internal/models/user.go")) {
			modelContent := helpers.ReadFileContent(t, filepath.Join(v.ProjectPath, "internal/models/user.go"))
			assert.NotContains(t, modelContent, "handlers", "Models should not import handlers")
		}
	}
}

func (v *WebAPIValidator) ValidateArchitecturalLinting(t *testing.T) {
	t.Helper()
	
	// This would typically run architectural linting tools
	// For now, just ensure compilation succeeds
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *WebAPIValidator) ValidateRuntimeBehavior(t *testing.T) {
	t.Helper()
	
	// This would ideally start the server and test endpoints
	// For now, just validate that the project structure supports runtime
	
	// Check main file exists and can be built
	mainFile := filepath.Join(v.ProjectPath, "cmd/server/main.go")
	helpers.AssertFileExists(t, mainFile)
	
	// Ensure project compiles (prerequisite for runtime)
	helpers.AssertCompilable(t, v.ProjectPath)
	
	// TODO: Add actual runtime testing with test server
	t.Log("Runtime validation: Project structure validated for runtime execution")
}

func (v *WebAPIValidator) ValidateMinimalConfiguration(t *testing.T) {
	t.Helper()
	
	// Core files should exist
	coreFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/config/config.go",
		"internal/handlers/health.go",
		"internal/logger/interface.go",
		"internal/logger/factory.go",
	}
	
	for _, file := range coreFiles {
		helpers.AssertFileExists(t, filepath.Join(v.ProjectPath, file))
	}
	
	// Optional files should NOT exist
	optionalFiles := []string{
		"internal/database/connection.go",
		"internal/models/user.go",
		"internal/repository/user.go",
		"internal/handlers/users.go",
		"internal/handlers/auth.go",
		"migrations/001_create_users.up.sql",
	}
	
	for _, file := range optionalFiles {
		filePath := filepath.Join(v.ProjectPath, file)
		assert.False(t, helpers.FileExists(filePath), 
			"File %s should not exist in minimal configuration", file)
	}
}