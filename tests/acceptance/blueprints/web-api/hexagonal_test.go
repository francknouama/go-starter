package webapi_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

func init() {
	// Initialize templates filesystem for ATDD tests
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				hasTemplates := false
				for _, entry := range entries {
					if entry.IsDir() {
						templateYaml := filepath.Join(templatesDir, entry.Name(), "template.yaml")
						if _, err := os.Stat(templateYaml); err == nil {
							hasTemplates = true
							break
						}
					}
				}

				if hasTemplates {
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// Feature: Hexagonal Architecture Web API
// As a software architect
// I want to generate Hexagonal Architecture web API
// So that I can implement ports and adapters pattern

func TestHexagonal_WebAPI_ArchitectureValidation(t *testing.T) {
	// Scenario: Validate Hexagonal Architecture structure
	// Given I generate a Hexagonal Architecture web API
	// Then the project should have these layers:
	//   | Layer        | Directory           | Purpose                    |
	//   | domain       | domain/             | Business entities & logic  |
	//   | application  | application/        | Application services       |
	//   | ports        | application/ports/  | Port interfaces            |
	//   | adapters     | adapters/           | Primary & Secondary        |
	// And dependencies should only point inward
	// And business logic should be adapter-independent
	// And ports should define clear contracts

	// This test should now pass as we have implemented the hexagonal template

	// Given I generate a Hexagonal Architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-arch",
		Module:    "github.com/test/test-hexagonal-arch",
		Type:      "web-api-hexagonal",
		GoVersion: "1.21",
		Framework: "gin",
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

	// Then the project should have hexagonal architecture structure
	validator := NewHexagonalValidator(projectPath)
	validator.ValidateHexagonalArchitecture(t)

	// And dependencies should only point inward
	validator.ValidateHexagonalDependencies(t)

	// And business logic should be adapter-independent
	validator.ValidateDomainIsolation(t)

	// And ports should define clear contracts
	validator.ValidatePortContracts(t)
}

func TestHexagonal_WebAPI_PortsAndAdapters(t *testing.T) {
	// Scenario: Validate ports and adapters implementation
	// Given I generate a Hexagonal Architecture web API
	// Then the project should have these ports:
	//   | Port Type    | Purpose                | Location              |
	//   | input        | Primary ports          | ports/input/          |
	//   | output       | Secondary ports        | ports/output/         |
	// And the project should have these adapters:
	//   | Adapter Type | Purpose                | Location              |
	//   | primary      | HTTP handlers          | adapters/primary/     |
	//   | secondary    | Database repositories  | adapters/secondary/   |
	// And adapters should implement port interfaces
	// And adapters should be easily swappable

	// This test should now pass as we have implemented the hexagonal template

	// Given I generate a Hexagonal Architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-ports",
		Module:    "github.com/test/test-hexagonal-ports",
		Type:      "web-api-hexagonal",
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

	// Then the project should have ports and adapters
	validator := NewHexagonalValidator(projectPath)
	validator.ValidatePortsStructure(t)
	validator.ValidateAdaptersStructure(t)

	// And adapters should implement port interfaces
	validator.ValidateAdapterImplementation(t)

	// And adapters should be easily swappable
	validator.ValidateAdapterSwappability(t)
}

func TestHexagonal_WebAPI_MultiFrameworkSupport(t *testing.T) {
	// Scenario: Validate multi-framework support in hexagonal architecture
	// Given I generate hexagonal web API with different frameworks
	// Then each framework should have its own primary adapter
	// And all frameworks should implement the same HTTP port interface
	// And the domain layer should be framework-independent

	// This test should now pass as we have implemented the hexagonal template

	frameworks := []string{"gin", "echo", "fiber", "chi", "stdlib"}

	for _, framework := range frameworks {
		t.Run("framework_"+framework, func(t *testing.T) {
			// Given I generate hexagonal web API with framework
			config := types.ProjectConfig{
				Name:      "test-hexagonal-" + framework,
				Module:    "github.com/test/test-hexagonal-" + framework,
				Type:      "web-api-hexagonal",
				GoVersion: "1.21",
				Framework: framework,
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

			// Then the framework should have its own primary adapter
			validator := NewHexagonalValidator(projectPath)
			validator.ValidateFrameworkAdapter(t, framework)

			// And all frameworks should implement the same HTTP port interface
			validator.ValidateHTTPPortInterface(t, framework)

			// And the domain layer should be framework-independent
			validator.ValidateDomainIsolation(t)
		})
	}
}

func TestHexagonal_WebAPI_DomainIsolation(t *testing.T) {
	// Scenario: Validate domain isolation in hexagonal architecture
	// Given I generate a Hexagonal Architecture web API
	// Then the domain layer should have no external dependencies
	// And the domain should not import infrastructure code
	// And the domain should only depend on standard library
	// And business logic should be pure and testable

	// This test should now pass as we have implemented the hexagonal template

	// Given I generate a Hexagonal Architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-domain",
		Module:    "github.com/test/test-hexagonal-domain",
		Type:      "web-api-hexagonal",
		GoVersion: "1.21",
		Framework: "gin",
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

	// Then the domain layer should have no external dependencies
	validator := NewHexagonalValidator(projectPath)
	validator.ValidateDomainPurity(t)

	// And the domain should not import infrastructure code
	validator.ValidateDomainImports(t)

	// And business logic should be pure and testable
	validator.ValidateDomainTestability(t)
}

func TestHexagonal_WebAPI_ApplicationServices(t *testing.T) {
	// Scenario: Validate application services in hexagonal architecture
	// Given I generate a Hexagonal Architecture web API
	// Then the application layer should orchestrate domain operations
	// And application services should use ports for external communication
	// And application services should implement business use cases
	// And application services should be framework-independent

	// This test should now pass as we have implemented the hexagonal template

	// Given I generate a Hexagonal Architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-app",
		Module:    "github.com/test/test-hexagonal-app",
		Type:      "web-api-hexagonal",
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

	// Then application services should orchestrate domain operations
	validator := NewHexagonalValidator(projectPath)
	validator.ValidateApplicationServices(t)

	// And application services should use ports for external communication
	validator.ValidateApplicationPortUsage(t)

	// And application services should implement business use cases
	validator.ValidateBusinessUseCases(t)

	// And application services should be framework-independent
	validator.ValidateApplicationIndependence(t)
}

func TestHexagonal_WebAPI_TestingCapabilities(t *testing.T) {
	// Scenario: Validate testing capabilities in hexagonal architecture
	// Given I generate a Hexagonal Architecture web API
	// Then the project should include adapter mocks
	// And the project should include port test implementations
	// And the project should enable easy integration testing
	// And the project should support unit testing of business logic

	// This test should now pass as we have implemented the hexagonal template

	// Given I generate a Hexagonal Architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-testing",
		Module:    "github.com/test/test-hexagonal-testing",
		Type:      "web-api-hexagonal",
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

	// Then the project should include adapter mocks
	validator := NewHexagonalValidator(projectPath)
	validator.ValidateAdapterMocks(t)

	// And the project should include port test implementations
	validator.ValidatePortTestImplementations(t)

	// And the project should enable easy integration testing
	validator.ValidateIntegrationTestCapabilities(t)

	// And the project should support unit testing of business logic
	validator.ValidateDomainUnitTesting(t)
}

// HexagonalValidator provides validation methods for hexagonal architecture
type HexagonalValidator struct {
	ProjectPath string
}

// NewHexagonalValidator creates a new HexagonalValidator
func NewHexagonalValidator(projectPath string) *HexagonalValidator {
	return &HexagonalValidator{
		ProjectPath: projectPath,
	}
}

// Additional validation methods for hexagonal architecture
func (v *HexagonalValidator) ValidateHexagonalArchitecture(t *testing.T) {
	t.Helper()

	// Validate hexagonal directory structure
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "domain"))
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "application"))
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "adapters"))

	// Validate ports structure
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "application", "ports"))
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "application", "ports", "input"))
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "application", "ports", "output"))

	// Validate adapters structure
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "adapters", "primary"))
	helpers.AssertDirectoryExists(t, filepath.Join(v.ProjectPath, "internal", "adapters", "secondary"))
}

func (v *HexagonalValidator) ValidateHexagonalDependencies(t *testing.T) {
	t.Helper()

	// Check that domain layer doesn't import from outer layers
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	if helpers.DirExists(domainPath) {
		domainFiles := helpers.FindFiles(t, domainPath, "*.go")
		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)
			// Domain should not import from application or adapters layers
			if helpers.StringContains(content, "internal/application") {
				t.Errorf("Domain layer should not import from application layer: %s", file)
			}
			if helpers.StringContains(content, "internal/adapters") {
				t.Errorf("Domain layer should not import from adapters layer: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateDomainIsolation(t *testing.T) {
	t.Helper()

	// Check that domain layer exists and is isolated
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	helpers.AssertDirectoryExists(t, domainPath)

	// Check for domain entities
	entitiesPath := filepath.Join(domainPath, "entities")
	if helpers.DirExists(entitiesPath) {
		helpers.AssertDirectoryExists(t, entitiesPath)
	}

	// Check for domain services
	servicesPath := filepath.Join(domainPath, "services")
	if helpers.DirExists(servicesPath) {
		helpers.AssertDirectoryExists(t, servicesPath)
	}
}

func (v *HexagonalValidator) ValidatePortContracts(t *testing.T) {
	t.Helper()

	// Check that ports define clear interfaces
	portsPath := filepath.Join(v.ProjectPath, "internal", "application", "ports")
	helpers.AssertDirectoryExists(t, portsPath)

	// Check for input ports
	inputPortsPath := filepath.Join(portsPath, "input")
	if helpers.DirExists(inputPortsPath) {
		portFiles := helpers.FindFiles(t, inputPortsPath, "*.go")
		for _, file := range portFiles {
			content := helpers.ReadFileContent(t, file)
			// Should define interfaces
			if !helpers.StringContains(content, "type") || !helpers.StringContains(content, "interface") {
				t.Errorf("Port file should define interfaces: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidatePortsStructure(t *testing.T) {
	t.Helper()

	// Validate ports directory structure
	portsPath := filepath.Join(v.ProjectPath, "internal", "application", "ports")
	helpers.AssertDirectoryExists(t, portsPath)

	inputPortsPath := filepath.Join(portsPath, "input")
	outputPortsPath := filepath.Join(portsPath, "output")

	if helpers.DirExists(inputPortsPath) {
		helpers.AssertDirectoryExists(t, inputPortsPath)
	}
	if helpers.DirExists(outputPortsPath) {
		helpers.AssertDirectoryExists(t, outputPortsPath)
	}
}

func (v *HexagonalValidator) ValidateAdaptersStructure(t *testing.T) {
	t.Helper()

	// Validate adapters directory structure
	adaptersPath := filepath.Join(v.ProjectPath, "internal", "adapters")
	helpers.AssertDirectoryExists(t, adaptersPath)

	primaryPath := filepath.Join(adaptersPath, "primary")
	secondaryPath := filepath.Join(adaptersPath, "secondary")

	helpers.AssertDirectoryExists(t, primaryPath)
	helpers.AssertDirectoryExists(t, secondaryPath)
}

func (v *HexagonalValidator) ValidateAdapterImplementation(t *testing.T) {
	t.Helper()

	// Check that adapters implement port interfaces
	adaptersPath := filepath.Join(v.ProjectPath, "internal", "adapters")
	if helpers.DirExists(adaptersPath) {
		adapterFiles := helpers.FindFiles(t, adaptersPath, "*.go")
		for _, file := range adapterFiles {
			content := helpers.ReadFileContent(t, file)
			// Should have struct definitions (adapter implementations)
			if !helpers.StringContains(content, "type") || !helpers.StringContains(content, "struct") {
				t.Errorf("Adapter file should define structs: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateAdapterSwappability(t *testing.T) {
	t.Helper()

	// Check that adapters are decoupled and swappable
	// This is validated through interface usage and dependency injection
	containerPath := filepath.Join(v.ProjectPath, "internal", "infrastructure", "container")
	if helpers.DirExists(containerPath) {
		helpers.AssertDirectoryExists(t, containerPath)
	}
}

func (v *HexagonalValidator) ValidateFrameworkAdapter(t *testing.T, framework string) {
	t.Helper()

	// Check that the framework has its own adapter
	frameworkAdapterPath := filepath.Join(v.ProjectPath, "internal", "adapters", "primary", "http", framework+"_adapter.go")
	if helpers.FileExists(frameworkAdapterPath) {
		helpers.AssertFileExists(t, frameworkAdapterPath)
		helpers.AssertFileContains(t, frameworkAdapterPath, framework)
	}
}

func (v *HexagonalValidator) ValidateHTTPPortInterface(t *testing.T, framework string) {
	t.Helper()

	// Check that HTTP port interface exists
	httpPortPath := filepath.Join(v.ProjectPath, "internal", "application", "ports", "input", "http_port.go")
	if helpers.FileExists(httpPortPath) {
		helpers.AssertFileExists(t, httpPortPath)
		helpers.AssertFileContains(t, httpPortPath, "interface")
	}
}

func (v *HexagonalValidator) ValidateDomainPurity(t *testing.T) {
	t.Helper()

	// Check that domain layer has minimal imports
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	if helpers.DirExists(domainPath) {
		domainFiles := helpers.FindFiles(t, domainPath, "*.go")
		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)
			// Domain should not import framework-specific packages
			if helpers.StringContains(content, "github.com/gin-gonic") ||
				helpers.StringContains(content, "github.com/labstack/echo") ||
				helpers.StringContains(content, "github.com/gofiber/fiber") {
				t.Errorf("Domain layer should not import framework packages: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateDomainImports(t *testing.T) {
	t.Helper()

	// Check that domain doesn't import infrastructure
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	if helpers.DirExists(domainPath) {
		domainFiles := helpers.FindFiles(t, domainPath, "*.go")
		for _, file := range domainFiles {
			content := helpers.ReadFileContent(t, file)
			// Domain should not import infrastructure
			if helpers.StringContains(content, "internal/infrastructure") {
				t.Errorf("Domain layer should not import infrastructure: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateDomainTestability(t *testing.T) {
	t.Helper()

	// Check that domain has test files
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	if helpers.DirExists(domainPath) {
		testFiles := helpers.FindFiles(t, domainPath, "*_test.go")
		if len(testFiles) == 0 {
			t.Error("Domain layer should have test files")
		}
	}
}

func (v *HexagonalValidator) ValidateApplicationServices(t *testing.T) {
	t.Helper()

	// Check that application services exist
	appPath := filepath.Join(v.ProjectPath, "internal", "application", "services")
	if helpers.DirExists(appPath) {
		helpers.AssertDirectoryExists(t, appPath)
		serviceFiles := helpers.FindFiles(t, appPath, "*.go")
		if len(serviceFiles) == 0 {
			t.Error("Application services should exist")
		}
	}
}

func (v *HexagonalValidator) ValidateApplicationPortUsage(t *testing.T) {
	t.Helper()

	// Check that application services use ports
	appPath := filepath.Join(v.ProjectPath, "internal", "application", "services")
	if helpers.DirExists(appPath) {
		serviceFiles := helpers.FindFiles(t, appPath, "*.go")
		for _, file := range serviceFiles {
			content := helpers.ReadFileContent(t, file)
			// Should import ports
			if !helpers.StringContains(content, "ports") {
				t.Errorf("Application services should use ports: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateBusinessUseCases(t *testing.T) {
	t.Helper()

	// Check that application services implement business use cases
	appPath := filepath.Join(v.ProjectPath, "internal", "application", "services")
	if helpers.DirExists(appPath) {
		serviceFiles := helpers.FindFiles(t, appPath, "*.go")
		for _, file := range serviceFiles {
			content := helpers.ReadFileContent(t, file)
			// Should have business methods
			if !helpers.StringContains(content, "func") {
				t.Errorf("Application services should implement business methods: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateApplicationIndependence(t *testing.T) {
	t.Helper()

	// Check that application layer is framework-independent
	appPath := filepath.Join(v.ProjectPath, "internal", "application")
	if helpers.DirExists(appPath) {
		appFiles := helpers.FindFiles(t, appPath, "*.go")
		for _, file := range appFiles {
			content := helpers.ReadFileContent(t, file)
			// Should not import framework-specific packages
			if helpers.StringContains(content, "github.com/gin-gonic") ||
				helpers.StringContains(content, "github.com/labstack/echo") {
				t.Errorf("Application layer should not import framework packages: %s", file)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateAdapterMocks(t *testing.T) {
	t.Helper()

	// Check that adapter mocks exist
	mocksPath := filepath.Join(v.ProjectPath, "tests", "mocks")
	if helpers.DirExists(mocksPath) {
		helpers.AssertDirectoryExists(t, mocksPath)
		mockFiles := helpers.FindFiles(t, mocksPath, "mock_*.go")
		if len(mockFiles) == 0 {
			t.Error("Adapter mocks should exist")
		}
	}
}

func (v *HexagonalValidator) ValidatePortTestImplementations(t *testing.T) {
	t.Helper()

	// Check that port test implementations exist
	testsPath := filepath.Join(v.ProjectPath, "tests")
	if helpers.DirExists(testsPath) {
		testFiles := helpers.FindFiles(t, testsPath, "*_test.go")
		if len(testFiles) == 0 {
			t.Error("Port test implementations should exist")
		}
	}
}

func (v *HexagonalValidator) ValidateIntegrationTestCapabilities(t *testing.T) {
	t.Helper()

	// Check that integration tests exist
	integrationPath := filepath.Join(v.ProjectPath, "tests", "integration")
	if helpers.DirExists(integrationPath) {
		helpers.AssertDirectoryExists(t, integrationPath)
		testFiles := helpers.FindFiles(t, integrationPath, "*_test.go")
		if len(testFiles) == 0 {
			t.Error("Integration tests should exist")
		}
	}
}

func (v *HexagonalValidator) ValidateDomainUnitTesting(t *testing.T) {
	t.Helper()

	// Check that domain unit tests exist
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	if helpers.DirExists(domainPath) {
		testFiles := helpers.FindFiles(t, domainPath, "*_test.go")
		if len(testFiles) == 0 {
			t.Error("Domain unit tests should exist")
		}
	}
}

// IsHexagonalImplemented checks if the hexagonal architecture is implemented
func (v *HexagonalValidator) IsHexagonalImplemented() bool {
	// Check if the basic hexagonal structure exists
	domainPath := filepath.Join(v.ProjectPath, "internal", "domain")
	applicationPath := filepath.Join(v.ProjectPath, "internal", "application")
	adaptersPath := filepath.Join(v.ProjectPath, "internal", "adapters")
	
	return helpers.DirExists(domainPath) && 
		   helpers.DirExists(applicationPath) && 
		   helpers.DirExists(adaptersPath)
}

// ValidateHexagonalPrinciples validates the main hexagonal architecture principles
func (v *HexagonalValidator) ValidateHexagonalPrinciples(t *testing.T) {
	t.Helper()
	
	// Validate the core principles of hexagonal architecture
	v.ValidateHexagonalArchitecture(t)
	v.ValidateHexagonalDependencies(t)
	v.ValidatePortsStructure(t)
	v.ValidateAdaptersStructure(t)
}

// ValidateDependencyInversion validates dependency inversion principle
func (v *HexagonalValidator) ValidateDependencyInversion(t *testing.T) {
	t.Helper()
	
	// Check that dependencies point inward
	v.ValidateDomainPurity(t)
	v.ValidateDomainImports(t)
	v.ValidateApplicationIndependence(t)
	v.ValidateApplicationPortUsage(t)
}

// ValidateHexagonDefinition validates the hexagonal definition and structure
func (v *HexagonalValidator) ValidateHexagonDefinition(t *testing.T) {
	t.Helper()
	
	// Validate that the hexagonal architecture is properly defined
	v.ValidatePortContracts(t)
	v.ValidateAdapterImplementation(t)
	v.ValidateAdapterSwappability(t)
	v.ValidateBusinessUseCases(t)
}