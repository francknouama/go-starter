package webapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

func init() {
	// Initialize templates filesystem for tests
	// We use DirFS pointing to the real blueprints directory
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			templates.SetTemplatesFS(os.DirFS(templatesDir))
			return
		}

		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			break
		}
		projectRoot = parent
	}

	panic("Could not find blueprints directory")
}

// TestArchitecture_DDD_BasicGeneration tests basic DDD project generation
func TestArchitecture_DDD_BasicGeneration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-api",
		Module:    "github.com/test/test-ddd-api",
		Type:      "web-api-ddd",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Variables: map[string]string{
			"ProjectName":    "test-ddd-api",
			"ModulePath":     "github.com/test/test-ddd-api",
			"GoVersion":      "1.21",
			"Framework":      "gin",
			"Logger":         "slog",
			"DatabaseDriver": "postgres",
			"DomainName":     "user",
			"AuthType":       "jwt",
		},
	}

	projectPath := helpers.GenerateProject(t, config)

	// Assert core DDD files exist
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/infrastructure/config/config.go",
		"internal/infrastructure/logger/slog.go",
		"internal/shared/errors/domain_errors.go",
		"internal/shared/errors/application_errors.go",
		"internal/shared/events/domain_event.go",
		"internal/shared/events/event_dispatcher.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)
	helpers.AssertGoModValid(t, filepath.Join(projectPath, "go.mod"), config.Module)
	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_DomainStructure validates DDD domain-centric organization
func TestArchitecture_DDD_DomainStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-domain",
		Module:    "github.com/test/test-ddd-domain",
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

	// DDD should have domain-centric structure
	expectedDomainStructure := []string{
		// Domain layer - core business logic
		"internal/domain/user",
		// Application layer - use cases and services
		"internal/application/user",
		"internal/application/auth",
		// Presentation layer - HTTP handlers and DTOs
		"internal/presentation/http/handlers",
		"internal/presentation/http/dto",
		"internal/presentation/http/middleware",
		// Infrastructure layer - external concerns
		"internal/infrastructure/persistence",
		"internal/infrastructure/logger",
		"internal/infrastructure/config",
		// Shared kernel
		"internal/shared/errors",
		"internal/shared/events",
		"internal/shared/valueobjects",
	}

	for _, structure := range expectedDomainStructure {
		helpers.AssertDirExists(t, filepath.Join(projectPath, structure))
	}

	// Should NOT have generic layered structure
	unexpectedDirs := []string{
		"internal/handlers",
		"internal/models",
		"internal/repository",
		"internal/services",
		"internal/adapters",
		"internal/domain/entities", // Clean Architecture pattern
		"internal/domain/usecases", // Clean Architecture pattern
	}

	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(projectPath, dir)
		_, err := helpers.GetFileInfo(dirPath)
		assert.Error(t, err, "Directory %s should not exist in DDD architecture", dir)
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_DomainModelStructure tests domain model organization
func TestArchitecture_DDD_DomainModelStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-model",
		Module:    "github.com/test/test-ddd-model",
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

	// Assert domain model files exist
	expectedDomainFiles := []string{
		"internal/domain/user/entity.go",
		"internal/domain/user/repository.go",
		"internal/domain/user/service.go",
		"internal/domain/user/events.go",
		"internal/domain/user/value_objects.go",
		"internal/domain/user/specifications.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedDomainFiles)

	// Assert domain entity contains business logic
	entityFile := filepath.Join(projectPath, "internal/domain/user/entity.go")
	entityContent := helpers.ReadFileContent(t, entityFile)

	// Entity should have domain methods
	assert.Contains(t, entityContent, "User", "Entity should define User aggregate")

	// Domain events should be present
	eventsFile := filepath.Join(projectPath, "internal/domain/user/events.go")
	helpers.AssertFileExists(t, eventsFile)

	// Value objects should be defined
	voFile := filepath.Join(projectPath, "internal/domain/user/value_objects.go")
	helpers.AssertFileExists(t, voFile)

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_ApplicationServices tests application layer services
func TestArchitecture_DDD_ApplicationServices(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-application",
		Module:    "github.com/test/test-ddd-application",
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

	// Assert application services exist
	expectedApplicationFiles := []string{
		"internal/application/user/commands.go",
		"internal/application/user/command_handlers.go",
		"internal/application/user/queries.go",
		"internal/application/user/query_handlers.go",
		"internal/application/user/dto.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedApplicationFiles)

	// Assert CQRS pattern implementation
	commandsFile := filepath.Join(projectPath, "internal/application/user/commands.go")
	commandsContent := helpers.ReadFileContent(t, commandsFile)
	assert.Contains(t, commandsContent, "Command", "Should implement command pattern")

	queriesFile := filepath.Join(projectPath, "internal/application/user/queries.go")
	queriesContent := helpers.ReadFileContent(t, queriesFile)
	assert.Contains(t, queriesContent, "Query", "Should implement query pattern")

	// Command handlers should orchestrate domain operations
	handlersFile := filepath.Join(projectPath, "internal/application/user/command_handlers.go")
	handlersContent := helpers.ReadFileContent(t, handlersFile)
	assert.Contains(t, handlersContent, "Handler", "Should implement command handlers")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_SharedKernel tests shared kernel components
func TestArchitecture_DDD_SharedKernel(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-shared",
		Module:    "github.com/test/test-ddd-shared",
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

	// Assert shared kernel files exist
	expectedSharedFiles := []string{
		"internal/shared/errors/domain_errors.go",
		"internal/shared/errors/application_errors.go",
		"internal/shared/events/domain_event.go",
		"internal/shared/events/event_dispatcher.go",
		"internal/shared/valueobjects/email.go",
		"internal/shared/valueobjects/id.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedSharedFiles)

	// Assert domain events structure
	domainEventFile := filepath.Join(projectPath, "internal/shared/events/domain_event.go")
	domainEventContent := helpers.ReadFileContent(t, domainEventFile)
	assert.Contains(t, domainEventContent, "DomainEvent", "Should define domain event interface")

	// Assert event dispatcher
	dispatcherFile := filepath.Join(projectPath, "internal/shared/events/event_dispatcher.go")
	dispatcherContent := helpers.ReadFileContent(t, dispatcherFile)
	assert.Contains(t, dispatcherContent, "EventDispatcher", "Should implement event dispatcher")

	// Assert value objects
	emailVOFile := filepath.Join(projectPath, "internal/shared/valueobjects/email.go")
	emailVOContent := helpers.ReadFileContent(t, emailVOFile)
	assert.Contains(t, emailVOContent, "Email", "Should define Email value object")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_DomainEvents tests domain event implementation
func TestArchitecture_DDD_DomainEvents(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-events",
		Module:    "github.com/test/test-ddd-events",
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

	// Assert domain events are defined
	userEventsFile := filepath.Join(projectPath, "internal/domain/user/events.go")
	userEventsContent := helpers.ReadFileContent(t, userEventsFile)

	// Should define user-specific events
	assert.Contains(t, userEventsContent, "UserCreated", "Should define UserCreated event")
	assert.Contains(t, userEventsContent, "UserUpdated", "Should define UserUpdated event")

	// Entity should raise domain events
	entityFile := filepath.Join(projectPath, "internal/domain/user/entity.go")
	entityContent := helpers.ReadFileContent(t, entityFile)

	// Entity should have method to raise events
	assert.Contains(t, entityContent, "DomainEvent", "Entity should work with domain events")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_RepositoryPattern tests repository pattern implementation
func TestArchitecture_DDD_RepositoryPattern(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-repository",
		Module:    "github.com/test/test-ddd-repository",
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

	// Assert repository interface in domain
	domainRepoFile := filepath.Join(projectPath, "internal/domain/user/repository.go")
	domainRepoContent := helpers.ReadFileContent(t, domainRepoFile)
	assert.Contains(t, domainRepoContent, "UserRepository", "Should define repository interface in domain")

	// Assert repository implementation in infrastructure
	// Note: DDD template structure might place this in infrastructure/persistence
	persistenceFiles := []string{
		"internal/infrastructure/persistence/user_repository.go",
	}

	for _, file := range persistenceFiles {
		filePath := filepath.Join(projectPath, file)
		if helpers.FileExists(filePath) {
			content := helpers.ReadFileContent(t, filePath)
			assert.Contains(t, content, "UserRepository", "Infrastructure should implement repository interface")
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_PresentationLayer tests presentation layer organization
func TestArchitecture_DDD_PresentationLayer(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-presentation",
		Module:    "github.com/test/test-ddd-presentation",
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

	// Assert presentation layer files
	expectedPresentationFiles := []string{
		"internal/presentation/http/handlers",
		"internal/presentation/http/dto",
		"internal/presentation/http/middleware",
	}

	for _, dir := range expectedPresentationFiles {
		helpers.AssertDirExists(t, filepath.Join(projectPath, dir))
	}

	// Presentation layer should coordinate between HTTP and application layer
	if helpers.DirExists(filepath.Join(projectPath, "internal/presentation/http/handlers")) {
		// Look for handler files
		handlerFiles := helpers.FindFiles(t, filepath.Join(projectPath, "internal/presentation/http/handlers"), "*.go")
		assert.NotEmpty(t, handlerFiles, "Should have HTTP handlers")

		if len(handlerFiles) > 0 {
			handlerContent := helpers.ReadFileContent(t, handlerFiles[0])
			// Handlers should use application services
			assert.Contains(t, handlerContent, "application", "Handlers should use application layer")
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_LoggerIntegration tests logger integration in DDD
func TestArchitecture_DDD_LoggerIntegration(t *testing.T) {
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
				Variables: map[string]string{
					"ProjectName":    "test-ddd-" + logger,
					"ModulePath":     "github.com/test/test-ddd-" + logger,
					"GoVersion":      "1.21",
					"Framework":      "gin",
					"Logger":         logger,
					"DatabaseDriver": "postgres",
					"DomainName":     "user",
					"AuthType":       "jwt",
				},
			}

			projectPath := helpers.GenerateProject(t, config)

			// Assert logger files in infrastructure
			expectedFiles := []string{
				"internal/infrastructure/logger/interface.go",
				"internal/infrastructure/logger/factory.go",
				"internal/infrastructure/logger/" + logger + ".go",
			}

			helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

			// Domain should not depend on specific logger implementations
			domainFiles := helpers.FindFiles(t, filepath.Join(projectPath, "internal/domain"), "*.go")
			for _, file := range domainFiles {
				content := helpers.ReadFileContent(t, file)
				// Domain should not import specific logger implementations
				forbiddenImports := []string{
					"go.uber.org/zap",
					"github.com/sirupsen/logrus",
					"github.com/rs/zerolog",
				}
				for _, forbidden := range forbiddenImports {
					assert.NotContains(t, content, forbidden,
						"Domain should not import specific logger: %s in file: %s", forbidden, file)
				}
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_DDD_Specifications tests specification pattern implementation
func TestArchitecture_DDD_Specifications(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-specs",
		Module:    "github.com/test/test-ddd-specs",
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

	// Assert specifications file exists
	specsFile := filepath.Join(projectPath, "internal/domain/user/specifications.go")
	helpers.AssertFileExists(t, specsFile)

	// Specifications should contain business rules
	specsContent := helpers.ReadFileContent(t, specsFile)
	assert.Contains(t, specsContent, "Specification", "Should implement specification pattern")
	assert.Contains(t, specsContent, "IsSatisfiedBy", "Should have specification interface")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_ValueObjects tests value object implementation
func TestArchitecture_DDD_ValueObjects(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-valueobjects",
		Module:    "github.com/test/test-ddd-valueobjects",
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

	// Assert value objects exist
	voFiles := []string{
		"internal/domain/user/value_objects.go",
		"internal/shared/valueobjects/email.go",
		"internal/shared/valueobjects/id.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, voFiles)

	// Value objects should have validation
	voFile := filepath.Join(projectPath, "internal/domain/user/value_objects.go")
	voContent := helpers.ReadFileContent(t, voFile)
	assert.Contains(t, voContent, "Valid", "Value objects should have validation")

	// Email value object should validate email format
	emailFile := filepath.Join(projectPath, "internal/shared/valueobjects/email.go")
	emailContent := helpers.ReadFileContent(t, emailFile)
	assert.Contains(t, emailContent, "Email", "Should define Email value object")
	assert.Contains(t, emailContent, "Valid", "Email should have validation")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_DDD_TestingStrategy tests DDD testing approach
func TestArchitecture_DDD_TestingStrategy(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-ddd-testing",
		Module:    "github.com/test/test-ddd-testing",
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

	// Assert test structure
	expectedTestFiles := []string{
		"tests/unit/domain/user_test.go",
		"tests/unit/application/user_test.go",
		"tests/integration/repository_test.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedTestFiles)

	// Domain tests should focus on business logic
	domainTestFile := filepath.Join(projectPath, "tests/unit/domain/user_test.go")
	if helpers.FileExists(domainTestFile) {
		domainTestContent := helpers.ReadFileContent(t, domainTestFile)
		assert.Contains(t, domainTestContent, "User", "Domain tests should test user entity")
	}

	// Application tests should test use cases
	appTestFile := filepath.Join(projectPath, "tests/unit/application/user_test.go")
	if helpers.FileExists(appTestFile) {
		appTestContent := helpers.ReadFileContent(t, appTestFile)
		assert.Contains(t, appTestContent, "Handler", "Application tests should test handlers")
	}

	helpers.AssertCompilable(t, projectPath)
}
