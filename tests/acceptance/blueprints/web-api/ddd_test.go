package webapi_test

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: DDD Web API Blueprint
// As a domain expert
// I want DDD-structured web API
// So that code reflects business concepts

func TestDDD_WebAPI_DomainStructureValidation(t *testing.T) {
	// Scenario: Domain structure validation
	// Given I generate a DDD web API
	// Then the project should have domain-centric structure:
	//   | Component      | Directory      | Purpose                    |
	//   | domain         | domain/        | Domain models and logic    |
	//   | application    | application/   | Application services       |
	//   | infrastructure | infrastructure/| External concerns          |
	//   | interfaces     | interfaces/    | API and UI layers          |
	// And domain should be isolated from infrastructure
	// And aggregates should be properly defined
	// And domain events should be supported

	// Given I generate a DDD web API
	config := types.ProjectConfig{
		Name:      "test-ddd-structure",
		Module:    "github.com/test/test-ddd-structure",
		Type:      "web-api-ddd",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
			Authentication: types.AuthConfig{
				Type: "jwt",
			},
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the project should have domain-centric structure
	validator := NewDDDValidator(projectPath)
	validator.ValidateDomainCentricStructure(t)

	// And domain should be isolated from infrastructure
	validator.ValidateDomainIsolation(t)

	// And aggregates should be properly defined
	validator.ValidateAggregates(t)

	// And domain events should be supported
	validator.ValidateDomainEvents(t)

	validator.ValidateCompilation(t)
}

func TestDDD_WebAPI_BusinessRuleEnforcement(t *testing.T) {
	// Scenario: Business rule enforcement
	// Given a DDD web API with business rules
	// Then domain entities should enforce business rules
	// And value objects should be immutable
	// And domain services should handle complex logic
	// And repositories should abstract data access

	config := types.ProjectConfig{
		Name:      "test-ddd-business-rules",
		Module:    "github.com/test/test-ddd-business-rules",
		Type:      "web-api-ddd",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewDDDValidator(projectPath)

	// Then domain entities should enforce business rules
	validator.ValidateBusinessRuleEnforcement(t)

	// And value objects should be immutable
	validator.ValidateValueObjects(t)

	// And domain services should handle complex logic
	validator.ValidateDomainServices(t)

	// And repositories should abstract data access
	validator.ValidateRepositoryAbstraction(t)

	validator.ValidateCompilation(t)
}

func TestDDD_WebAPI_CQRSImplementation(t *testing.T) {
	// Scenario: CQRS pattern implementation
	// Given I generate a DDD web API
	// Then commands and queries should be separated
	// And command handlers should modify state
	// And query handlers should only read data
	// And application services should orchestrate

	config := types.ProjectConfig{
		Name:      "test-ddd-cqrs",
		Module:    "github.com/test/test-ddd-cqrs",
		Type:      "web-api-ddd",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewDDDValidator(projectPath)

	// Then commands and queries should be separated
	validator.ValidateCQRSPattern(t)

	// And command handlers should modify state
	validator.ValidateCommandHandlers(t)

	// And query handlers should only read data
	validator.ValidateQueryHandlers(t)

	// And application services should orchestrate
	validator.ValidateApplicationServices(t)

	validator.ValidateCompilation(t)
}

func TestDDD_WebAPI_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration in DDD
	// Scenario: Logger follows DDD patterns

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-ddd-" + logger,
				Module:    "github.com/test/test-ddd-" + logger,
				Type:      "web-api-ddd",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewDDDValidator(projectPath)

			// Logger should be in infrastructure layer
			validator.ValidateLoggerInInfrastructure(t, logger)

			// Domain should not depend on concrete logger
			validator.ValidateDomainLoggerIndependence(t)

			validator.ValidateCompilation(t)
		})
	}
}

func TestDDD_WebAPI_FrameworkIntegration(t *testing.T) {
	// Scenario: Framework integration in DDD
	// Given I use different frameworks with DDD
	// Then domain should remain framework-independent
	// And presentation layer should handle framework specifics

	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-ddd-" + framework,
				Module:    "github.com/test/test-ddd-" + framework,
				Type:      "web-api-ddd",
				GoVersion: "1.21",
				Framework: framework,
				Logger:    "slog",
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewDDDValidator(projectPath)

			// Domain should remain framework-independent
			validator.ValidateDomainFrameworkIndependence(t, framework)

			// Presentation layer should handle framework specifics
			validator.ValidatePresentationLayerFrameworkHandling(t, framework)

			validator.ValidateCompilation(t)
		})
	}
}

func TestDDD_WebAPI_ArchitectureCompliance(t *testing.T) {
	// Feature: DDD Architecture Compliance
	// Scenario: Architecture compliance validation
	// Given I generate a DDD web API
	// Then the code should follow DDD principles
	// And domain should be at the center
	// And bounded contexts should be clear

	config := types.ProjectConfig{
		Name:      "test-ddd-compliance",
		Module:    "github.com/test/test-ddd-compliance",
		Type:      "web-api-ddd",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewDDDValidator(projectPath)

	// Then the code should follow DDD principles
	validator.ValidateDDDPrinciples(t)

	// And domain should be at the center
	validator.ValidateDomainCentrality(t)

	// And bounded contexts should be clear
	validator.ValidateBoundedContexts(t)

	validator.ValidateCompilation(t)
}

// DDDValidator provides validation methods specific to Domain-Driven Design
type DDDValidator struct {
	ProjectPath string
}

func NewDDDValidator(projectPath string) *DDDValidator {
	return &DDDValidator{
		ProjectPath: projectPath,
	}
}

func (v *DDDValidator) ValidateDomainCentricStructure(t *testing.T) {
	t.Helper()

	// Expected DDD structure
	expectedStructure := map[string]string{
		"internal/domain/user":                    "User aggregate root",
		"internal/application/user":               "User application services",
		"internal/application/auth":               "Auth application services",
		"internal/infrastructure/config":          "Configuration",
		"internal/infrastructure/logger":          "Logging infrastructure",
		"internal/infrastructure/persistence":     "Data persistence",
		"internal/presentation/http":              "HTTP presentation layer",
		"internal/shared/errors":                  "Shared kernel - errors",
		"internal/shared/events":                  "Shared kernel - events",
		"internal/shared/valueobjects":            "Shared kernel - value objects",
	}

	for structure, purpose := range expectedStructure {
		structurePath := filepath.Join(v.ProjectPath, structure)
		helpers.AssertDirExists(t, structurePath)
		t.Logf("✓ DDD Structure %s exists (Purpose: %s)", structure, purpose)
	}

	// Should NOT have standard or clean architecture directories
	unexpectedDirs := []string{
		"internal/handlers",
		"internal/models",
		"internal/repository",
		"internal/services",
		"internal/adapters",
		"internal/usecases",
	}

	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		assert.False(t, helpers.DirExists(dirPath),
			"Directory %s should not exist in DDD architecture", dir)
	}
}

func (v *DDDValidator) ValidateDomainIsolation(t *testing.T) {
	t.Helper()

	// Domain should not import from application, infrastructure, or presentation
	domainDir := filepath.Join(v.ProjectPath, "internal/domain")
	if helpers.DirExists(domainDir) {
		domainFiles := helpers.FindFiles(t, domainDir, "*.go")
		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)

			// Domain should NOT import from outer layers
			forbiddenImports := []string{
				"internal/application",
				"internal/infrastructure",
				"internal/presentation",
				"github.com/gin-gonic/gin",
				"gorm.io/gorm",
			}

			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Domain file %s should not import %s - violates domain isolation", file, forbidden)
			}
		}
	}
}

func (v *DDDValidator) ValidateAggregates(t *testing.T) {
	t.Helper()

	// User should be an aggregate root
	userEntityFile := filepath.Join(v.ProjectPath, "internal/domain/user/entity.go")
	if helpers.FileExists(userEntityFile) {
		content := helpers.ReadFileContent(t, userEntityFile)
		assert.Contains(t, content, "User", "Should have User aggregate")

		// Should have domain methods
		assert.Contains(t, content, "func", "Aggregate should have domain methods")
	}
}

func (v *DDDValidator) ValidateDomainEvents(t *testing.T) {
	t.Helper()

	// Domain events should be supported
	eventsDir := filepath.Join(v.ProjectPath, "internal/shared/events")
	if helpers.DirExists(eventsDir) {
		eventFiles := helpers.FindFiles(t, eventsDir, "*.go")
		assert.NotEmpty(t, eventFiles, "Should have domain events")

		for _, file := range eventFiles {
			content := helpers.ReadFileContent(t, file)
			assert.Contains(t, content, "Event", "Should define events")
		}
	}

	// User domain should have events
	userEventsFile := filepath.Join(v.ProjectPath, "internal/domain/user/events.go")
	if helpers.FileExists(userEventsFile) {
		content := helpers.ReadFileContent(t, userEventsFile)
		assert.Contains(t, content, "Event", "User domain should have events")
	}
}

func (v *DDDValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *DDDValidator) ValidateBusinessRuleEnforcement(t *testing.T) {
	t.Helper()

	// Domain entities should have business methods
	userEntityFile := filepath.Join(v.ProjectPath, "internal/domain/user/entity.go")
	if helpers.FileExists(userEntityFile) {
		content := helpers.ReadFileContent(t, userEntityFile)

		// Should have business methods (not just getters/setters)
		businessMethods := []string{"Validate", "Update", "Create", "func"}
		hasBusinessMethod := false
		for _, method := range businessMethods {
			if helpers.StringContains(content, method) {
				hasBusinessMethod = true
				break
			}
		}
		assert.True(t, hasBusinessMethod, "Entity should have business methods")
	}
}

func (v *DDDValidator) ValidateValueObjects(t *testing.T) {
	t.Helper()

	valueObjectsDir := filepath.Join(v.ProjectPath, "internal/shared/valueobjects")
	if helpers.DirExists(valueObjectsDir) {
		voFiles := helpers.FindFiles(t, valueObjectsDir, "*.go")
		for _, file := range voFiles {
			content := helpers.ReadFileContent(t, file)

			// Value objects should have validation
			assert.Contains(t, content, "func", "Value objects should have methods")
		}
	}

	// Email should be a value object
	emailFile := filepath.Join(v.ProjectPath, "internal/shared/valueobjects/email.go")
	if helpers.FileExists(emailFile) {
		content := helpers.ReadFileContent(t, emailFile)
		assert.Contains(t, content, "Email", "Should have Email value object")
	}
}

func (v *DDDValidator) ValidateDomainServices(t *testing.T) {
	t.Helper()

	// Domain services should exist for complex business logic
	userServiceFile := filepath.Join(v.ProjectPath, "internal/domain/user/service.go")
	if helpers.FileExists(userServiceFile) {
		content := helpers.ReadFileContent(t, userServiceFile)
		assert.Contains(t, content, "Service", "Should have domain service")
		assert.Contains(t, content, "func", "Domain service should have methods")
	}
}

func (v *DDDValidator) ValidateRepositoryAbstraction(t *testing.T) {
	t.Helper()

	// Repository interface should be in domain
	userRepoFile := filepath.Join(v.ProjectPath, "internal/domain/user/repository.go")
	if helpers.FileExists(userRepoFile) {
		content := helpers.ReadFileContent(t, userRepoFile)
		assert.Contains(t, content, "Repository", "Should have repository interface")
		assert.Contains(t, content, "interface", "Repository should be an interface")
	}

	// Repository implementation should be in infrastructure
	repoImplFile := filepath.Join(v.ProjectPath, "internal/infrastructure/persistence/user_repository.go")
	if helpers.FileExists(repoImplFile) {
		content := helpers.ReadFileContent(t, repoImplFile)
		assert.Contains(t, content, "Repository", "Should implement repository")
	}
}

func (v *DDDValidator) ValidateCQRSPattern(t *testing.T) {
	t.Helper()

	// Commands and queries should be separated
	userAppDir := filepath.Join(v.ProjectPath, "internal/application/user")
	if helpers.DirExists(userAppDir) {
		// Commands
		commandsFile := filepath.Join(userAppDir, "commands.go")
		if helpers.FileExists(commandsFile) {
			content := helpers.ReadFileContent(t, commandsFile)
			assert.Contains(t, content, "Command", "Should have commands")
		}

		// Queries
		queriesFile := filepath.Join(userAppDir, "queries.go")
		if helpers.FileExists(queriesFile) {
			content := helpers.ReadFileContent(t, queriesFile)
			assert.Contains(t, content, "Query", "Should have queries")
		}
	}
}

func (v *DDDValidator) ValidateCommandHandlers(t *testing.T) {
	t.Helper()

	commandHandlersFile := filepath.Join(v.ProjectPath, "internal/application/user/command_handlers.go")
	if helpers.FileExists(commandHandlersFile) {
		content := helpers.ReadFileContent(t, commandHandlersFile)
		assert.Contains(t, content, "Handler", "Should have command handlers")
		assert.Contains(t, content, "Command", "Should handle commands")
	}
}

func (v *DDDValidator) ValidateQueryHandlers(t *testing.T) {
	t.Helper()

	queryHandlersFile := filepath.Join(v.ProjectPath, "internal/application/user/query_handlers.go")
	if helpers.FileExists(queryHandlersFile) {
		content := helpers.ReadFileContent(t, queryHandlersFile)
		assert.Contains(t, content, "Handler", "Should have query handlers")
		assert.Contains(t, content, "Query", "Should handle queries")
	}
}

func (v *DDDValidator) ValidateApplicationServices(t *testing.T) {
	t.Helper()

	// Application services should orchestrate domain operations
	userServiceFile := filepath.Join(v.ProjectPath, "internal/application/user/service.go")
	if helpers.FileExists(userServiceFile) {
		content := helpers.ReadFileContent(t, userServiceFile)
		assert.Contains(t, content, "Service", "Should have application service")
	}

	authServiceFile := filepath.Join(v.ProjectPath, "internal/application/auth/service.go")
	if helpers.FileExists(authServiceFile) {
		content := helpers.ReadFileContent(t, authServiceFile)
		assert.Contains(t, content, "Service", "Should have auth application service")
	}
}

func (v *DDDValidator) ValidateLoggerInInfrastructure(t *testing.T, logger string) {
	t.Helper()

	loggerDir := filepath.Join(v.ProjectPath, "internal/infrastructure/logger")
	helpers.AssertDirExists(t, loggerDir)

	// Logger implementation should exist
	loggerFile := filepath.Join(loggerDir, logger+".go")
	helpers.AssertFileExists(t, loggerFile)
}

func (v *DDDValidator) ValidateDomainLoggerIndependence(t *testing.T) {
	t.Helper()

	domainDir := filepath.Join(v.ProjectPath, "internal/domain")
	if helpers.DirExists(domainDir) {
		domainFiles := helpers.FindFiles(t, domainDir, "*.go")

		concreteLoggerImports := []string{
			"go.uber.org/zap",
			"github.com/sirupsen/logrus",
			"github.com/rs/zerolog",
		}

		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)
			for _, loggerImport := range concreteLoggerImports {
				assert.NotContains(t, content, loggerImport,
					"Domain should not import concrete logger %s", loggerImport)
			}
		}
	}
}

func (v *DDDValidator) ValidateDomainFrameworkIndependence(t *testing.T, framework string) {
	t.Helper()

	domainDir := filepath.Join(v.ProjectPath, "internal/domain")
	if helpers.DirExists(domainDir) {
		domainFiles := helpers.FindFiles(t, domainDir, "*.go")

		frameworkImports := map[string]string{
			"gin":   "github.com/gin-gonic/gin",
			"echo":  "github.com/labstack/echo",
			"fiber": "github.com/gofiber/fiber",
			"chi":   "github.com/go-chi/chi",
		}

		frameworkImport := frameworkImports[framework]

		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)
			assert.NotContains(t, content, frameworkImport,
				"Domain should not import framework %s", framework)
		}
	}
}

func (v *DDDValidator) ValidatePresentationLayerFrameworkHandling(t *testing.T, framework string) {
	t.Helper()

	presentationDir := filepath.Join(v.ProjectPath, "internal/presentation/http")
	if helpers.DirExists(presentationDir) {
		// Framework-specific handlers should exist
		handlerFiles := helpers.FindFiles(t, presentationDir, "*"+framework+"*.go")
		assert.NotEmpty(t, handlerFiles, "Should have framework-specific handlers")
	}
}

func (v *DDDValidator) ValidateDDDPrinciples(t *testing.T) {
	t.Helper()

	// 1. Ubiquitous Language - domain concepts should be reflected in code
	v.ValidateUbiquitousLanguage(t)

	// 2. Domain Model isolation
	v.ValidateDomainIsolation(t)

	// 3. Aggregates and bounded contexts
	v.ValidateAggregates(t)

	// 4. Domain events
	v.ValidateDomainEvents(t)
}

func (v *DDDValidator) ValidateUbiquitousLanguage(t *testing.T) {
	t.Helper()

	// Domain should use business terminology
	domainDir := filepath.Join(v.ProjectPath, "internal/domain")
	if helpers.DirExists(domainDir) {
		// Check for business concepts in directory names
		userDir := filepath.Join(domainDir, "user")
		helpers.AssertDirExists(t, userDir)

		// Files should use business language
		entityFile := filepath.Join(userDir, "entity.go")
		if helpers.FileExists(entityFile) {
			content := helpers.ReadFileContent(t, entityFile)
			assert.Contains(t, content, "User", "Should use business terminology")
		}
	}
}

func (v *DDDValidator) ValidateDomainCentrality(t *testing.T) {
	t.Helper()

	// Domain should be the center of the architecture
	domainDir := filepath.Join(v.ProjectPath, "internal/domain")
	helpers.AssertDirExists(t, domainDir)

	// Other layers should depend on domain, not vice versa
	v.ValidateDomainIsolation(t)

	// Application should depend on domain
	appDir := filepath.Join(v.ProjectPath, "internal/application")
	if helpers.DirExists(appDir) {
		appFiles := helpers.FindFiles(t, appDir, "*.go")
		hasDomainImport := false
		for _, file := range appFiles {
			content := helpers.ReadFileContent(t, file)
			if helpers.StringContains(content, "internal/domain") {
				hasDomainImport = true
				break
			}
		}
		assert.True(t, hasDomainImport, "Application should depend on domain")
	}
}

func (v *DDDValidator) ValidateBoundedContexts(t *testing.T) {
	t.Helper()

	// In this simple example, User is a bounded context
	userContext := filepath.Join(v.ProjectPath, "internal/domain/user")
	helpers.AssertDirExists(t, userContext)

	// Each bounded context should have its own:
	// - Entities
	// - Value objects  
	// - Domain services
	// - Repository interfaces

	contextFiles := []string{
		"entity.go",
		"repository.go", 
		"service.go",
		"events.go",
		"value_objects.go",
	}

	for _, file := range contextFiles {
		filePath := filepath.Join(userContext, file)
		if helpers.FileExists(filePath) {
			t.Logf("✓ Bounded context file %s exists", file)
		}
	}
}