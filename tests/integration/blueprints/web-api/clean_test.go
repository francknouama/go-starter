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

// TestArchitecture_Clean_BasicGeneration tests basic Clean Architecture project generation
func TestArchitecture_Clean_BasicGeneration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-api",
		Module:    "github.com/test/test-clean-api",
		Type:      "web-api-clean",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Assert core files exist
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/infrastructure/config/config.go",
		"internal/adapters/controllers/health_controller.go",
		"internal/infrastructure/logger/interface.go",
		"internal/infrastructure/logger/factory.go",
		"internal/infrastructure/logger/slog.go",
		"internal/infrastructure/container/container.go",
		"internal/infrastructure/web/router.go",
		"internal/infrastructure/web/factory.go",
		"internal/infrastructure/web/adapters/gin_adapter.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)
	helpers.AssertGoModValid(t, filepath.Join(projectPath, "go.mod"), config.Module)
	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_LayeredStructure validates Clean Architecture layered organization
func TestArchitecture_Clean_LayeredStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-layers",
		Module:    "github.com/test/test-clean-layers",
		Type:      "web-api-clean",
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

	// Clean Architecture should have distinct layers
	expectedLayers := []string{
		// Entities layer (innermost)
		"internal/domain/entities",
		// Use cases layer
		"internal/domain/usecases",
		"internal/domain/ports",
		// Interface adapters layer
		"internal/adapters/controllers",
		"internal/adapters/presenters",
		// Frameworks & drivers layer (outermost)
		"internal/infrastructure/persistence",
		"internal/infrastructure/web",
		"internal/infrastructure/services",
		"internal/infrastructure/logger",
		"internal/infrastructure/config",
		"internal/infrastructure/container",
	}

	for _, layer := range expectedLayers {
		helpers.AssertDirExists(t, filepath.Join(projectPath, layer))
	}

	// Should NOT have standard architecture structure
	unexpectedDirs := []string{
		"internal/handlers",
		"internal/models",
		"internal/repository",
		"internal/services",
	}

	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(projectPath, dir)
		_, err := helpers.GetFileInfo(dirPath)
		assert.Error(t, err, "Directory %s should not exist in Clean Architecture", dir)
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_DependencyInversionPrinciple tests proper dependency directions
func TestArchitecture_Clean_DependencyInversionPrinciple(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-dip",
		Module:    "github.com/test/test-clean-dip",
		Type:      "web-api-clean",
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

	// Entities should not import anything from outer layers
	entityFile := filepath.Join(projectPath, "internal/domain/entities/user.go")
	entityContent := helpers.ReadFileContent(t, entityFile)

	// Entities should NOT import infrastructure, adapters, or external frameworks
	forbiddenImports := []string{
		"internal/infrastructure",
		"internal/adapters",
		"github.com/gin-gonic/gin",
		"gorm.io/gorm",
	}

	for _, forbidden := range forbiddenImports {
		assert.NotContains(t, entityContent, forbidden,
			"Entity should not import %s - violates dependency rule", forbidden)
	}

	// Use cases should depend only on entities and ports (interfaces)
	usecaseFile := filepath.Join(projectPath, "internal/domain/usecases/user_usecase.go")
	usecaseContent := helpers.ReadFileContent(t, usecaseFile)

	// Use cases should import entities and ports
	assert.Contains(t, usecaseContent, "entities", "Use case should import entities")
	assert.Contains(t, usecaseContent, "ports", "Use case should import ports")

	// Use cases should NOT import infrastructure or adapters
	for _, forbidden := range []string{"internal/infrastructure", "internal/adapters"} {
		assert.NotContains(t, usecaseContent, forbidden,
			"Use case should not import %s - violates dependency rule", forbidden)
	}

	// Infrastructure implements ports/interfaces defined in domain
	repositoryFile := filepath.Join(projectPath, "internal/infrastructure/persistence/user_repository.go")
	repositoryContent := helpers.ReadFileContent(t, repositoryFile)

	// Repository should implement interfaces from ports
	assert.Contains(t, repositoryContent, "ports", "Repository should implement ports")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_PortsAndAdapters tests interfaces and their implementations
func TestArchitecture_Clean_PortsAndAdapters(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-ports",
		Module:    "github.com/test/test-clean-ports",
		Type:      "web-api-clean",
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

	// Assert ports (interfaces) exist
	expectedPorts := []string{
		"internal/domain/ports/repositories.go",
		"internal/domain/ports/services.go",
		"internal/domain/ports/web.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedPorts)

	// Assert repository interface is defined
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/domain/ports/repositories.go"), "UserRepository")

	// Assert repository implementation exists
	helpers.AssertFileExists(t, filepath.Join(projectPath, "internal/infrastructure/persistence/user_repository.go"))

	// Assert implementation references the interface
	helpers.AssertFileContains(t,
		filepath.Join(projectPath, "internal/infrastructure/persistence/user_repository.go"),
		"ports.UserRepository")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_DependencyInjection tests dependency injection container
func TestArchitecture_Clean_DependencyInjection(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-di",
		Module:    "github.com/test/test-clean-di",
		Type:      "web-api-clean",
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

	// Assert container exists
	containerFile := filepath.Join(projectPath, "internal/infrastructure/container/container.go")
	helpers.AssertFileExists(t, containerFile)

	// Container should wire dependencies
	helpers.AssertFileContains(t, containerFile, "Container")
	helpers.AssertFileContains(t, containerFile, "NewContainer")

	// Main should use container
	mainFile := filepath.Join(projectPath, "cmd/server/main.go")
	helpers.AssertFileContains(t, mainFile, "container")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_LoggerVariations tests different logger implementations in Clean Architecture
func TestArchitecture_Clean_LoggerVariations(t *testing.T) {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-clean-" + logger,
				Module:       "github.com/test/test-clean-" + logger,
				Type:         "web-api",
				Architecture: "clean",
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Assert logger files in infrastructure layer
			expectedFiles := []string{
				"internal/infrastructure/logger/interface.go",
				"internal/infrastructure/logger/factory.go",
				"internal/infrastructure/logger/" + logger + ".go",
			}

			helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

			// Logger should be in infrastructure layer, not leaking to domain
			domainFiles := []string{
				"internal/domain/entities/user.go",
				"internal/domain/usecases/user_usecase.go",
			}

			for _, file := range domainFiles {
				filePath := filepath.Join(projectPath, file)
				if helpers.FileExists(filePath) {
					content := helpers.ReadFileContent(t, filePath)
					// Domain should not directly import specific logger implementations
					loggerImports := []string{
						"go.uber.org/zap",
						"github.com/sirupsen/logrus",
						"github.com/rs/zerolog",
					}
					for _, loggerImport := range loggerImports {
						assert.NotContains(t, content, loggerImport,
							"Domain layer should not import specific logger: %s", loggerImport)
					}
				}
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Clean_WebFrameworkAbstraction tests framework abstraction
func TestArchitecture_Clean_WebFrameworkAbstraction(t *testing.T) {
	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-clean-" + framework,
				Module:       "github.com/test/test-clean-" + framework,
				Type:         "web-api",
				Architecture: "clean",
				GoVersion:    "1.21",
				Framework:    framework,
				Logger:       "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Driver: "postgres",
						ORM:    "gorm",
					},
				},
			}

			projectPath := helpers.GenerateProject(t, config)

			// Assert framework adapter exists
			adapterFile := filepath.Join(projectPath, "internal/infrastructure/web/adapters/"+framework+"_adapter.go")
			helpers.AssertFileExists(t, adapterFile)

			// Assert web factory abstracts framework selection
			factoryFile := filepath.Join(projectPath, "internal/infrastructure/web/factory.go")
			helpers.AssertFileExists(t, factoryFile)
			helpers.AssertFileContains(t, factoryFile, framework)

			// Controllers should not directly import framework
			controllerFile := filepath.Join(projectPath, "internal/adapters/controllers/user_controller.go")
			if helpers.FileExists(controllerFile) {
				controllerContent := helpers.ReadFileContent(t, controllerFile)
				frameworkImports := []string{
					"github.com/gin-gonic/gin",
					"github.com/labstack/echo",
					"github.com/gofiber/fiber",
					"github.com/go-chi/chi",
				}
				for _, frameworkImport := range frameworkImports {
					assert.NotContains(t, controllerContent, frameworkImport,
						"Controller should not directly import framework: %s", frameworkImport)
				}
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Clean_TestingStructure tests testing organization in Clean Architecture
func TestArchitecture_Clean_TestingStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-testing",
		Module:    "github.com/test/test-clean-testing",
		Type:      "web-api-clean",
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
		"tests/unit/entities_test.go",
		"tests/unit/usecases_test.go",
		"tests/integration/api_test.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedTestFiles)

	// Assert mocks for interfaces
	mockFiles := []string{
		"tests/mocks/MockUserRepository.go",
		"tests/mocks/MockLogger.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, mockFiles)

	// Unit tests should test domain logic in isolation
	usecaseTestFile := filepath.Join(projectPath, "tests/unit/usecases_test.go")
	usecaseTestContent := helpers.ReadFileContent(t, usecaseTestFile)

	// Should import mocks for external dependencies
	assert.Contains(t, usecaseTestContent, "mocks", "Use case tests should use mocks")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Clean_BusinessLogicIsolation tests business logic separation
func TestArchitecture_Clean_BusinessLogicIsolation(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-clean-business",
		Module:    "github.com/test/test-clean-business",
		Type:      "web-api-clean",
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

	// Business logic should be in entities and use cases
	entityFile := filepath.Join(projectPath, "internal/domain/entities/user.go")
	usecaseFile := filepath.Join(projectPath, "internal/domain/usecases/user_usecase.go")

	helpers.AssertFileExists(t, entityFile)
	helpers.AssertFileExists(t, usecaseFile)

	// Use cases should contain business logic methods
	usecaseContent := helpers.ReadFileContent(t, usecaseFile)
	assert.Contains(t, usecaseContent, "UseCase", "Use case should define business operations")

	// Controllers should be thin - just HTTP concerns
	controllerFile := filepath.Join(projectPath, "internal/adapters/controllers/user_controller.go")
	if helpers.FileExists(controllerFile) {
		controllerContent := helpers.ReadFileContent(t, controllerFile)
		// Controller should delegate to use cases
		assert.Contains(t, controllerContent, "usecase", "Controller should delegate to use cases")
	}

	helpers.AssertCompilable(t, projectPath)
}
