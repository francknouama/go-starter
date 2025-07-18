package webapi_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: Clean Architecture Web API
// As a software architect
// I want to generate Clean Architecture web API
// So that I can maintain separation of concerns

func TestClean_WebAPI_LayerValidation(t *testing.T) {
	// Scenario: Validate Clean Architecture layers
	// Given I generate a Clean Architecture web API
	// Then the project should have these layers:
	//   | Layer        | Directory     | Purpose                    |
	//   | entities     | entities/     | Business entities          |
	//   | usecases     | usecases/     | Business logic             |
	//   | interfaces   | interfaces/   | Contracts and adapters     |
	//   | frameworks   | frameworks/   | External frameworks        |
	// And dependencies should only point inward
	// And business logic should be framework-independent
	// And interfaces should define contracts clearly

	// Given I generate a Clean Architecture web API
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

	// Then the project should have these layers
	validator := NewCleanArchitectureValidator(projectPath)
	validator.ValidateLayerStructure(t)

	// And dependencies should only point inward
	validator.ValidateDependencyInversion(t)

	// And business logic should be framework-independent
	validator.ValidateBusinessLogicIsolation(t)

	// And interfaces should define contracts clearly
	validator.ValidateInterfaceContracts(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestClean_WebAPI_DependencyInjection(t *testing.T) {
	// Scenario: Dependency injection validation
	// Given a Clean Architecture web API
	// Then dependency injection should be configured
	// And repositories should be interfaces
	// And use cases should depend on interfaces only
	// And frameworks should implement interfaces

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

	// Then dependency injection should be configured
	validator := NewCleanArchitectureValidator(projectPath)
	validator.ValidateDependencyInjectionContainer(t)

	// And repositories should be interfaces
	validator.ValidateRepositoryInterfaces(t)

	// And use cases should depend on interfaces only
	validator.ValidateUseCaseDependencies(t)

	// And frameworks should implement interfaces
	validator.ValidateFrameworkImplementations(t)

	validator.ValidateCompilation(t)
}

func TestClean_WebAPI_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration in Clean Architecture
	// Scenario: Logger integration follows Clean Architecture patterns

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-clean-" + logger,
				Module:    "github.com/test/test-clean-" + logger,
				Type:      "web-api-clean",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewCleanArchitectureValidator(projectPath)
			
			// Logger should be in infrastructure layer
			validator.ValidateLoggerInInfrastructure(t, logger)
			
			// Logger should be injected through interfaces
			validator.ValidateLoggerInterface(t)
			
			// Business logic should not depend on concrete logger
			validator.ValidateBusinessLogicLoggerIndependence(t)
			
			validator.ValidateCompilation(t)
		})
	}
}

func TestClean_WebAPI_FrameworkAbstraction(t *testing.T) {
	// Scenario: Framework abstraction validation
	// Given I generate Clean Architecture with different frameworks
	// Then business logic should be framework-independent
	// And framework specifics should be in outer layer

	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-clean-" + framework,
				Module:    "github.com/test/test-clean-" + framework,
				Type:      "web-api-clean",
				GoVersion: "1.21",
				Framework: framework,
				Logger:    "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Driver: "postgres",
						ORM:    "gorm",
					},
				},
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewCleanArchitectureValidator(projectPath)
			
			// Framework should be abstracted in infrastructure
			validator.ValidateFrameworkAbstraction(t, framework)
			
			// Business logic should not import framework
			validator.ValidateBusinessLogicFrameworkIndependence(t, framework)
			
			validator.ValidateCompilation(t)
		})
	}
}

func TestClean_WebAPI_DatabaseIntegration(t *testing.T) {
	// Feature: Database Integration in Clean Architecture
	// Scenario: Database integration follows repository pattern

	databases := []types.DatabaseConfig{
		{Driver: "postgres", ORM: "gorm"},
		{Driver: "mysql", ORM: "sqlx"},
	}

	for _, db := range databases {
		t.Run("Database_"+db.Driver+"_"+db.ORM, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-clean-db-" + db.Driver + "-" + db.ORM,
				Module:    "github.com/test/test-clean-db-" + db.Driver + "-" + db.ORM,
				Type:      "web-api-clean",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    "slog",
				Features: &types.Features{
					Database: db,
				},
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewCleanArchitectureValidator(projectPath)
			
			// Database should be in infrastructure layer
			validator.ValidateDatabaseInInfrastructure(t, db)
			
			// Repository interfaces should be in domain
			validator.ValidateRepositoryInterfaces(t)
			
			// Business logic should not depend on database specifics
			validator.ValidateBusinessLogicDatabaseIndependence(t)
			
			validator.ValidateCompilation(t)
		})
	}
}

func TestClean_WebAPI_ArchitectureCompliance(t *testing.T) {
	// Feature: Clean Architecture Compliance
	// Scenario: Architecture compliance validation
	// Given I generate a Clean Architecture web API
	// Then the code should follow Clean Architecture principles
	// And dependency directions should be correct
	// And layer boundaries should be enforced

	config := types.ProjectConfig{
		Name:      "test-clean-compliance",
		Module:    "github.com/test/test-clean-compliance",
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

	validator := NewCleanArchitectureValidator(projectPath)
	
	// Then the code should follow Clean Architecture principles
	validator.ValidateCleanArchitecturePrinciples(t)
	
	// And dependency directions should be correct
	validator.ValidateDependencyInversion(t)
	
	// And layer boundaries should be enforced
	validator.ValidateLayerBoundaries(t)
	
	validator.ValidateCompilation(t)
}

// CleanArchitectureValidator provides validation methods specific to Clean Architecture
type CleanArchitectureValidator struct {
	ProjectPath string
}

func NewCleanArchitectureValidator(projectPath string) *CleanArchitectureValidator {
	return &CleanArchitectureValidator{
		ProjectPath: projectPath,
	}
}

func (v *CleanArchitectureValidator) ValidateLayerStructure(t *testing.T) {
	t.Helper()
	
	// Expected layers in Clean Architecture
	expectedLayers := map[string]string{
		"internal/domain/entities":              "Business entities",
		"internal/domain/usecases":              "Business logic",
		"internal/domain/ports":                 "Contracts and adapters",
		"internal/adapters/controllers":         "Interface adapters",
		"internal/adapters/presenters":          "Output adapters",
		"internal/infrastructure/persistence":   "Data access layer",
		"internal/infrastructure/web":           "Web framework layer",
		"internal/infrastructure/services":      "External services",
		"internal/infrastructure/logger":        "Logging infrastructure",
		"internal/infrastructure/config":        "Configuration",
		"internal/infrastructure/container":     "Dependency injection",
	}
	
	for layer, purpose := range expectedLayers {
		layerPath := filepath.Join(v.ProjectPath, layer)
		helpers.AssertDirExists(t, layerPath)
		t.Logf("✓ Layer %s exists (Purpose: %s)", layer, purpose)
	}
	
	// Should NOT have standard architecture structure
	unexpectedDirs := []string{
		"internal/handlers",
		"internal/models", 
		"internal/repository",
		"internal/services",
	}
	
	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		assert.False(t, helpers.DirExists(dirPath), 
			"Directory %s should not exist in Clean Architecture", dir)
	}
}

func (v *CleanArchitectureValidator) ValidateDependencyInversion(t *testing.T) {
	t.Helper()
	
	// Entities should not import anything from outer layers
	entitiesDir := filepath.Join(v.ProjectPath, "internal/domain/entities")
	if helpers.DirExists(entitiesDir) {
		entityFiles := helpers.FindFiles(t, entitiesDir, "*.go")
		for _, file := range entityFiles {
			content := helpers.ReadFileContent(t, file)
			
			// Entities should NOT import infrastructure, adapters, or external frameworks
			forbiddenImports := []string{
				"internal/infrastructure",
				"internal/adapters",
				"github.com/gin-gonic/gin",
				"gorm.io/gorm",
				"github.com/labstack/echo",
			}
			
			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Entity %s should not import %s - violates dependency rule", file, forbidden)
			}
		}
	}
	
	// Use cases should depend only on entities and ports
	usecasesDir := filepath.Join(v.ProjectPath, "internal/domain/usecases")
	if helpers.DirExists(usecasesDir) {
		usecaseFiles := helpers.FindFiles(t, usecasesDir, "*.go")
		for _, file := range usecaseFiles {
			content := helpers.ReadFileContent(t, file)
			
			// Use cases should import entities and ports
			assert.Contains(t, content, "entities", "Use case should import entities")
			assert.Contains(t, content, "ports", "Use case should import ports")
			
			// Use cases should NOT import infrastructure or adapters
			forbiddenImports := []string{"internal/infrastructure", "internal/adapters"}
			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Use case %s should not import %s - violates dependency rule", file, forbidden)
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateBusinessLogicIsolation(t *testing.T) {
	t.Helper()
	
	// Check that business logic (entities and use cases) don't depend on frameworks
	businessLogicDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
	}
	
	frameworkImports := []string{
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo",
		"github.com/gofiber/fiber",
		"github.com/go-chi/chi",
		"gorm.io/gorm",
	}
	
	for _, dir := range businessLogicDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			files := helpers.FindFiles(t, dirPath, "*.go")
			for _, file := range files {
				content := helpers.ReadFileContent(t, file)
				for _, framework := range frameworkImports {
					assert.NotContains(t, content, framework,
						"Business logic file %s should not import framework %s", file, framework)
				}
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateInterfaceContracts(t *testing.T) {
	t.Helper()
	
	// Ports directory should contain interface definitions
	portsDir := filepath.Join(v.ProjectPath, "internal/domain/ports")
	if helpers.DirExists(portsDir) {
		portFiles := helpers.FindFiles(t, portsDir, "*.go")
		assert.NotEmpty(t, portFiles, "Ports directory should contain interface files")
		
		for _, file := range portFiles {
			content := helpers.ReadFileContent(t, file)
			assert.Contains(t, content, "interface", 
				"Port file %s should define interfaces", file)
		}
	}
}

func (v *CleanArchitectureValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *CleanArchitectureValidator) ValidateDependencyInjectionContainer(t *testing.T) {
	t.Helper()
	
	containerFile := filepath.Join(v.ProjectPath, "internal/infrastructure/container/container.go")
	helpers.AssertFileExists(t, containerFile)
	
	content := helpers.ReadFileContent(t, containerFile)
	assert.Contains(t, content, "Container", "Should have dependency injection container")
}

func (v *CleanArchitectureValidator) ValidateRepositoryInterfaces(t *testing.T) {
	t.Helper()
	
	// Repository interfaces should be in ports
	portsFile := filepath.Join(v.ProjectPath, "internal/domain/ports/repositories.go")
	if helpers.FileExists(portsFile) {
		content := helpers.ReadFileContent(t, portsFile)
		assert.Contains(t, content, "Repository", "Should define repository interfaces")
		assert.Contains(t, content, "interface", "Should use interface keyword")
	}
}

func (v *CleanArchitectureValidator) ValidateUseCaseDependencies(t *testing.T) {
	t.Helper()
	
	usecasesDir := filepath.Join(v.ProjectPath, "internal/domain/usecases")
	if helpers.DirExists(usecasesDir) {
		usecaseFiles := helpers.FindFiles(t, usecasesDir, "*.go")
		for _, file := range usecaseFiles {
			content := helpers.ReadFileContent(t, file)
			
			// Use cases should depend on ports (interfaces), not concrete implementations
			if strings.Contains(content, "Repository") {
				assert.Contains(t, content, "ports", 
					"Use case should depend on repository interface from ports")
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateFrameworkImplementations(t *testing.T) {
	t.Helper()
	
	// Infrastructure should implement ports
	persistenceDir := filepath.Join(v.ProjectPath, "internal/infrastructure/persistence")
	if helpers.DirExists(persistenceDir) {
		persistenceFiles := helpers.FindFiles(t, persistenceDir, "*repository*.go")
		for _, file := range persistenceFiles {
			content := helpers.ReadFileContent(t, file)
			assert.Contains(t, content, "ports", 
				"Infrastructure should implement interfaces from ports")
		}
	}
}

func (v *CleanArchitectureValidator) ValidateLoggerInInfrastructure(t *testing.T, logger string) {
	t.Helper()
	
	loggerDir := filepath.Join(v.ProjectPath, "internal/infrastructure/logger")
	helpers.AssertDirExists(t, loggerDir)
	
	// Logger implementation should exist
	loggerFile := filepath.Join(loggerDir, logger+".go")
	helpers.AssertFileExists(t, loggerFile)
}

func (v *CleanArchitectureValidator) ValidateLoggerInterface(t *testing.T) {
	t.Helper()
	
	interfaceFile := filepath.Join(v.ProjectPath, "internal/infrastructure/logger/interface.go")
	helpers.AssertFileExists(t, interfaceFile)
	
	content := helpers.ReadFileContent(t, interfaceFile)
	assert.Contains(t, content, "interface", "Should define logger interface")
}

func (v *CleanArchitectureValidator) ValidateBusinessLogicLoggerIndependence(t *testing.T) {
	t.Helper()
	
	// Business logic should not import concrete logger implementations
	businessLogicDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
	}
	
	concreteLoggerImports := []string{
		"go.uber.org/zap",
		"github.com/sirupsen/logrus",
		"github.com/rs/zerolog",
	}
	
	for _, dir := range businessLogicDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			files := helpers.FindFiles(t, dirPath, "*.go")
			for _, file := range files {
				content := helpers.ReadFileContent(t, file)
				for _, loggerImport := range concreteLoggerImports {
					assert.NotContains(t, content, loggerImport,
						"Business logic should not import concrete logger %s", loggerImport)
				}
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateFrameworkAbstraction(t *testing.T, framework string) {
	t.Helper()
	
	// Framework should be abstracted in infrastructure/web
	webDir := filepath.Join(v.ProjectPath, "internal/infrastructure/web")
	helpers.AssertDirExists(t, webDir)
	
	// Framework-specific adapter should exist
	adapterFile := filepath.Join(webDir, "adapters", framework+"_adapter.go")
	helpers.AssertFileExists(t, adapterFile)
}

func (v *CleanArchitectureValidator) ValidateBusinessLogicFrameworkIndependence(t *testing.T, framework string) {
	t.Helper()
	
	businessLogicDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
	}
	
	frameworkImports := map[string]string{
		"gin":   "github.com/gin-gonic/gin",
		"echo":  "github.com/labstack/echo",
		"fiber": "github.com/gofiber/fiber",
		"chi":   "github.com/go-chi/chi",
	}
	
	frameworkImport := frameworkImports[framework]
	
	for _, dir := range businessLogicDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			files := helpers.FindFiles(t, dirPath, "*.go")
			for _, file := range files {
				content := helpers.ReadFileContent(t, file)
				assert.NotContains(t, content, frameworkImport,
					"Business logic should not import framework %s", framework)
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateDatabaseInInfrastructure(t *testing.T, db types.DatabaseConfig) {
	t.Helper()
	
	persistenceDir := filepath.Join(v.ProjectPath, "internal/infrastructure/persistence")
	helpers.AssertDirExists(t, persistenceDir)
	
	// Database connection should be in infrastructure
	dbFile := filepath.Join(persistenceDir, "database.go")
	helpers.AssertFileExists(t, dbFile)
}

func (v *CleanArchitectureValidator) ValidateBusinessLogicDatabaseIndependence(t *testing.T) {
	t.Helper()
	
	businessLogicDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
	}
	
	databaseImports := []string{
		"gorm.io/gorm",
		"database/sql",
		"github.com/jmoiron/sqlx",
	}
	
	for _, dir := range businessLogicDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			files := helpers.FindFiles(t, dirPath, "*.go")
			for _, file := range files {
				content := helpers.ReadFileContent(t, file)
				for _, dbImport := range databaseImports {
					assert.NotContains(t, content, dbImport,
						"Business logic should not import database %s", dbImport)
				}
			}
		}
	}
}

func (v *CleanArchitectureValidator) ValidateCleanArchitecturePrinciples(t *testing.T) {
	t.Helper()
	
	// 1. Independence of frameworks
	v.ValidateBusinessLogicFrameworkIndependence(t, "gin") // Test with gin as example
	
	// 2. Testability - business logic should be isolated
	v.ValidateBusinessLogicIsolation(t)
	
	// 3. Independence of UI - business logic should not depend on web layer
	businessLogicDirs := []string{
		"internal/domain/entities",
		"internal/domain/usecases",
	}
	
	for _, dir := range businessLogicDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		if helpers.DirExists(dirPath) {
			files := helpers.FindFiles(t, dirPath, "*.go")
			for _, file := range files {
				content := helpers.ReadFileContent(t, file)
				assert.NotContains(t, content, "internal/adapters/controllers",
					"Business logic should not depend on UI layer")
			}
		}
	}
	
	// 4. Independence of database
	v.ValidateBusinessLogicDatabaseIndependence(t)
}

func (v *CleanArchitectureValidator) ValidateLayerBoundaries(t *testing.T) {
	t.Helper()
	
	// Entities should be in the center and not import from other layers
	entitiesDir := filepath.Join(v.ProjectPath, "internal/domain/entities")
	if helpers.DirExists(entitiesDir) {
		entityFiles := helpers.FindFiles(t, entitiesDir, "*.go")
		for _, file := range entityFiles {
			content := helpers.ReadFileContent(t, file)
			
			// Entities should not import from outer layers
			outerLayerImports := []string{
				"internal/domain/usecases",
				"internal/adapters",
				"internal/infrastructure",
			}
			
			for _, outerImport := range outerLayerImports {
				assert.NotContains(t, content, outerImport,
					"Entity should not import from outer layer %s", outerImport)
			}
		}
	}
	
	// Use cases should not import from adapters or infrastructure
	usecasesDir := filepath.Join(v.ProjectPath, "internal/domain/usecases")
	if helpers.DirExists(usecasesDir) {
		usecaseFiles := helpers.FindFiles(t, usecasesDir, "*.go")
		for _, file := range usecaseFiles {
			content := helpers.ReadFileContent(t, file)
			
			outerLayerImports := []string{
				"internal/adapters",
				"internal/infrastructure",
			}
			
			for _, outerImport := range outerLayerImports {
				assert.NotContains(t, content, outerImport,
					"Use case should not import from outer layer %s", outerImport)
			}
		}
	}
}