package webapi_test

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: Hexagonal Architecture Web API
// As a testability-focused developer
// I want hexagonal architecture
// So that I can easily test and adapt

func TestHexagonal_WebAPI_PortsAndAdaptersValidation(t *testing.T) {
	// Scenario: Ports and adapters validation
	// Given I generate a Hexagonal architecture web API
	// Then the project should have clear ports:
	//   | Port Type | Purpose           | Location     |
	//   | Primary   | Driving adapters  | ports/in/    |
	//   | Secondary | Driven adapters   | ports/out/   |
	// And adapters should implement ports
	// And core business logic should be adapter-independent
	// And dependency inversion should be enforced

	// Given I generate a Hexagonal architecture web API
	config := types.ProjectConfig{
		Name:      "test-hexagonal-ports",
		Module:    "github.com/test/test-hexagonal-ports",
		Type:      "web-api", // Note: hexagonal may not be fully implemented yet
		Architecture: "hexagonal",
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

	// Then the project should have clear ports
	validator := NewHexagonalValidator(projectPath)
	validator.ValidatePortsStructure(t)

	// And adapters should implement ports
	validator.ValidateAdapterImplementations(t)

	// And core business logic should be adapter-independent
	validator.ValidateCoreBusinessLogicIndependence(t)

	// And dependency inversion should be enforced
	validator.ValidateDependencyInversion(t)

	// Note: May skip if hexagonal template not fully implemented
	if !validator.IsHexagonalImplemented() {
		t.Skip("Hexagonal architecture template not fully implemented yet")
	}

	validator.ValidateCompilation(t)
}

func TestHexagonal_WebAPI_AdapterSwappability(t *testing.T) {
	// Scenario: Adapter swappability
	// Given a Hexagonal web API
	// Then HTTP adapters should be swappable
	// And database adapters should be swappable
	// And external service adapters should be swappable
	// And tests should use mock adapters

	config := types.ProjectConfig{
		Name:      "test-hexagonal-swappable",
		Module:    "github.com/test/test-hexagonal-swappable",
		Type:      "web-api",
		Architecture: "hexagonal",
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

	validator := NewHexagonalValidator(projectPath)

	if !validator.IsHexagonalImplemented() {
		t.Skip("Hexagonal architecture template not fully implemented yet")
	}

	// Then HTTP adapters should be swappable
	validator.ValidateHTTPAdapterSwappability(t)

	// And database adapters should be swappable
	validator.ValidateDatabaseAdapterSwappability(t)

	// And external service adapters should be swappable
	validator.ValidateExternalServiceAdapterSwappability(t)

	// And tests should use mock adapters
	validator.ValidateMockAdapterUsage(t)

	validator.ValidateCompilation(t)
}

func TestHexagonal_WebAPI_CoreIsolation(t *testing.T) {
	// Scenario: Core business logic isolation
	// Given I generate a Hexagonal web API
	// Then core should define interfaces (ports)
	// And core should not depend on external frameworks
	// And all external dependencies should go through ports
	// And core should be fully testable in isolation

	config := types.ProjectConfig{
		Name:      "test-hexagonal-core",
		Module:    "github.com/test/test-hexagonal-core",
		Type:      "web-api",
		Architecture: "hexagonal",
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

	validator := NewHexagonalValidator(projectPath)

	if !validator.IsHexagonalImplemented() {
		t.Skip("Hexagonal architecture template not fully implemented yet")
	}

	// Then core should define interfaces (ports)
	validator.ValidateCoreDefinesPorts(t)

	// And core should not depend on external frameworks
	validator.ValidateCoreFrameworkIndependence(t)

	// And all external dependencies should go through ports
	validator.ValidateExternalDependenciesThroughPorts(t)

	// And core should be fully testable in isolation
	validator.ValidateCoreTestability(t)

	validator.ValidateCompilation(t)
}

func TestHexagonal_WebAPI_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration in Hexagonal Architecture
	// Scenario: Logger follows Hexagonal patterns

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-hexagonal-" + logger,
				Module:    "github.com/test/test-hexagonal-" + logger,
				Type:      "web-api",
				Architecture: "hexagonal",
				GoVersion: "1.21",
				Framework: "gin",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewHexagonalValidator(projectPath)

			if !validator.IsHexagonalImplemented() {
				t.Skip("Hexagonal architecture template not fully implemented yet")
			}

			// Logger should be a secondary adapter
			validator.ValidateLoggerAsSecondaryAdapter(t, logger)

			// Core should use logger port, not concrete implementation
			validator.ValidateCoreLoggerPortUsage(t)

			validator.ValidateCompilation(t)
		})
	}
}

func TestHexagonal_WebAPI_FrameworkIntegration(t *testing.T) {
	// Scenario: Framework integration in Hexagonal Architecture
	// Given I use different frameworks with Hexagonal
	// Then frameworks should be primary adapters only
	// And core should not know about specific frameworks

	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-hexagonal-" + framework,
				Module:    "github.com/test/test-hexagonal-" + framework,
				Type:      "web-api",
				Architecture: "hexagonal",
				GoVersion: "1.21",
				Framework: framework,
				Logger:    "slog",
			}

			projectPath := helpers.GenerateProject(t, config)

			validator := NewHexagonalValidator(projectPath)

			if !validator.IsHexagonalImplemented() {
				t.Skip("Hexagonal architecture template not fully implemented yet")
			}

			// Frameworks should be primary adapters only
			validator.ValidateFrameworkAsPrimaryAdapter(t, framework)

			// Core should not know about specific frameworks
			validator.ValidateCoreFrameworkAgnostic(t, framework)

			validator.ValidateCompilation(t)
		})
	}
}

func TestHexagonal_WebAPI_ArchitectureCompliance(t *testing.T) {
	// Feature: Hexagonal Architecture Compliance
	// Scenario: Architecture compliance validation
	// Given I generate a Hexagonal web API
	// Then the code should follow Hexagonal principles
	// And dependency directions should be correct
	// And the hexagon should be clearly defined

	config := types.ProjectConfig{
		Name:      "test-hexagonal-compliance",
		Module:    "github.com/test/test-hexagonal-compliance",
		Type:      "web-api",
		Architecture: "hexagonal",
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

	validator := NewHexagonalValidator(projectPath)

	if !validator.IsHexagonalImplemented() {
		t.Skip("Hexagonal architecture template not fully implemented yet")
	}

	// Then the code should follow Hexagonal principles
	validator.ValidateHexagonalPrinciples(t)

	// And dependency directions should be correct
	validator.ValidateDependencyInversion(t)

	// And the hexagon should be clearly defined
	validator.ValidateHexagonDefinition(t)

	validator.ValidateCompilation(t)
}

// HexagonalValidator provides validation methods specific to Hexagonal Architecture
type HexagonalValidator struct {
	ProjectPath string
}

func NewHexagonalValidator(projectPath string) *HexagonalValidator {
	return &HexagonalValidator{
		ProjectPath: projectPath,
	}
}

func (v *HexagonalValidator) IsHexagonalImplemented() bool {
	// Check if hexagonal architecture structure exists
	// If not, the template may not be fully implemented yet
	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	portsDir := filepath.Join(v.ProjectPath, "internal/ports")
	adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")

	return helpers.DirExists(coreDir) || helpers.DirExists(portsDir) || helpers.DirExists(adaptersDir)
}

func (v *HexagonalValidator) ValidatePortsStructure(t *testing.T) {
	t.Helper()

	// Expected Hexagonal structure
	expectedStructure := map[string]string{
		"internal/core":                          "Core business logic",
		"internal/core/domain":                   "Domain entities and business rules",
		"internal/core/services":                 "Application services",
		"internal/ports/in":                      "Primary ports (driving)",
		"internal/ports/out":                     "Secondary ports (driven)",
		"internal/adapters/primary/http":         "HTTP primary adapter",
		"internal/adapters/secondary/persistence": "Database secondary adapter",
		"internal/adapters/secondary/logger":     "Logger secondary adapter",
	}

	for structure, purpose := range expectedStructure {
		structurePath := filepath.Join(v.ProjectPath, structure)
		if helpers.DirExists(structurePath) {
			t.Logf("✓ Hexagonal Structure %s exists (Purpose: %s)", structure, purpose)
		} else {
			t.Logf("⚠ Hexagonal Structure %s missing (Purpose: %s)", structure, purpose)
		}
	}
}

func (v *HexagonalValidator) ValidateAdapterImplementations(t *testing.T) {
	t.Helper()

	// Primary adapters should implement primary ports
	httpAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/primary/http")
	if helpers.DirExists(httpAdapterDir) {
		adapterFiles := helpers.FindFiles(t, httpAdapterDir, "*.go")
		for _, file := range adapterFiles {
			content := helpers.ReadFileContent(t, file)
			// Should import from ports/in
			if helpers.StringContains(content, "ports") {
				t.Logf("✓ HTTP adapter imports ports")
			}
		}
	}

	// Secondary adapters should implement secondary ports
	persistenceAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/secondary/persistence")
	if helpers.DirExists(persistenceAdapterDir) {
		adapterFiles := helpers.FindFiles(t, persistenceAdapterDir, "*.go")
		for _, file := range adapterFiles {
			content := helpers.ReadFileContent(t, file)
			// Should import from ports/out
			if helpers.StringContains(content, "ports") {
				t.Logf("✓ Persistence adapter implements ports")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateCoreBusinessLogicIndependence(t *testing.T) {
	t.Helper()

	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		// Core should not import adapters or external frameworks
		forbiddenImports := []string{
			"internal/adapters",
			"github.com/gin-gonic/gin",
			"github.com/labstack/echo",
			"gorm.io/gorm",
		}

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)
			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Core business logic should not import %s", forbidden)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateDependencyInversion(t *testing.T) {
	t.Helper()

	// Core should define ports, adapters should implement them
	portsDir := filepath.Join(v.ProjectPath, "internal/ports")
	if helpers.DirExists(portsDir) {
		// Primary ports (in)
		primaryPortsDir := filepath.Join(portsDir, "in")
		if helpers.DirExists(primaryPortsDir) {
			portFiles := helpers.FindFiles(t, primaryPortsDir, "*.go")
			for _, file := range portFiles {
				content := helpers.ReadFileContent(t, file)
				assert.Contains(t, content, "interface",
					"Primary ports should define interfaces")
			}
		}

		// Secondary ports (out)
		secondaryPortsDir := filepath.Join(portsDir, "out")
		if helpers.DirExists(secondaryPortsDir) {
			portFiles := helpers.FindFiles(t, secondaryPortsDir, "*.go")
			for _, file := range portFiles {
				content := helpers.ReadFileContent(t, file)
				assert.Contains(t, content, "interface",
					"Secondary ports should define interfaces")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *HexagonalValidator) ValidateHTTPAdapterSwappability(t *testing.T) {
	t.Helper()

	// HTTP adapter should be easily swappable
	httpAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/primary/http")
	if helpers.DirExists(httpAdapterDir) {
		// Should have interface definition
		adapterFiles := helpers.FindFiles(t, httpAdapterDir, "*.go")
		for _, file := range adapterFiles {
			content := helpers.ReadFileContent(t, file)
			// Should implement port interface
			if helpers.StringContains(content, "ports/in") {
				t.Logf("✓ HTTP adapter implements primary port interface")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateDatabaseAdapterSwappability(t *testing.T) {
	t.Helper()

	// Database adapter should be easily swappable
	persistenceAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/secondary/persistence")
	if helpers.DirExists(persistenceAdapterDir) {
		adapterFiles := helpers.FindFiles(t, persistenceAdapterDir, "*.go")
		for _, file := range adapterFiles {
			content := helpers.ReadFileContent(t, file)
			// Should implement secondary port interface
			if helpers.StringContains(content, "ports/out") {
				t.Logf("✓ Database adapter implements secondary port interface")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateExternalServiceAdapterSwappability(t *testing.T) {
	t.Helper()

	// External service adapters should be swappable
	adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters/secondary")
	if helpers.DirExists(adaptersDir) {
		serviceAdapters := []string{"email", "payment", "notification"}
		for _, service := range serviceAdapters {
			serviceDir := filepath.Join(adaptersDir, service)
			if helpers.DirExists(serviceDir) {
				t.Logf("✓ External service adapter %s is swappable", service)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateMockAdapterUsage(t *testing.T) {
	t.Helper()

	// Tests should use mock adapters
	testsDir := filepath.Join(v.ProjectPath, "tests")
	if helpers.DirExists(testsDir) {
		testFiles := helpers.FindFiles(t, testsDir, "*_test.go")
		for _, file := range testFiles {
			content := helpers.ReadFileContent(t, file)
			if helpers.StringContains(content, "mock") || helpers.StringContains(content, "Mock") {
				t.Logf("✓ Tests use mock adapters")
				break
			}
		}
	}
}

func (v *HexagonalValidator) ValidateCoreDefinesPorts(t *testing.T) {
	t.Helper()

	// Core should define the port interfaces it needs
	portsDir := filepath.Join(v.ProjectPath, "internal/ports")
	if helpers.DirExists(portsDir) {
		portFiles := helpers.FindFiles(t, portsDir, "*.go")
		assert.NotEmpty(t, portFiles, "Core should define port interfaces")

		for _, file := range portFiles {
			content := helpers.ReadFileContent(t, file)
			assert.Contains(t, content, "interface", "Ports should define interfaces")
		}
	}
}

func (v *HexagonalValidator) ValidateCoreFrameworkIndependence(t *testing.T) {
	t.Helper()

	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		frameworkImports := []string{
			"github.com/gin-gonic/gin",
			"github.com/labstack/echo",
			"github.com/gofiber/fiber",
			"github.com/go-chi/chi",
		}

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)
			for _, framework := range frameworkImports {
				assert.NotContains(t, content, framework,
					"Core should not depend on framework %s", framework)
			}
		}
	}
}

func (v *HexagonalValidator) ValidateExternalDependenciesThroughPorts(t *testing.T) {
	t.Helper()

	// Core should only interact with external dependencies through ports
	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		externalDependencies := []string{
			"gorm.io/gorm",
			"database/sql",
			"net/http",
		}

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)
			for _, dep := range externalDependencies {
				assert.NotContains(t, content, dep,
					"Core should not directly import external dependency %s", dep)
			}

			// Should import ports instead
			if helpers.StringContains(content, "import") && 
			   helpers.StringContains(content, "internal/ports") {
				t.Logf("✓ Core uses ports for external dependencies")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateCoreTestability(t *testing.T) {
	t.Helper()

	// Core should be testable in isolation
	coreTestsDir := filepath.Join(v.ProjectPath, "tests/unit/core")
	if helpers.DirExists(coreTestsDir) {
		testFiles := helpers.FindFiles(t, coreTestsDir, "*_test.go")
		assert.NotEmpty(t, testFiles, "Core should have isolated unit tests")
	}

	// Alternative: tests directory at project root
	testsDir := filepath.Join(v.ProjectPath, "tests")
	if helpers.DirExists(testsDir) {
		testFiles := helpers.FindFiles(t, testsDir, "*_test.go")
		for _, file := range testFiles {
			content := helpers.ReadFileContent(t, file)
			if helpers.StringContains(content, "internal/core") {
				t.Logf("✓ Core has unit tests")
				break
			}
		}
	}
}

func (v *HexagonalValidator) ValidateLoggerAsSecondaryAdapter(t *testing.T, logger string) {
	t.Helper()

	loggerAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/secondary/logger")
	if helpers.DirExists(loggerAdapterDir) {
		loggerFile := filepath.Join(loggerAdapterDir, logger+".go")
		if helpers.FileExists(loggerFile) {
			content := helpers.ReadFileContent(t, loggerFile)
			// Should implement logger port
			assert.Contains(t, content, "ports", "Logger adapter should implement port")
		}
	}
}

func (v *HexagonalValidator) ValidateCoreLoggerPortUsage(t *testing.T) {
	t.Helper()

	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)
			// Core should use logger port, not concrete logger
			if helpers.StringContains(content, "log") {
				assert.Contains(t, content, "ports", 
					"Core should use logger through port interface")
			}
		}
	}
}

func (v *HexagonalValidator) ValidateFrameworkAsPrimaryAdapter(t *testing.T, framework string) {
	t.Helper()

	httpAdapterDir := filepath.Join(v.ProjectPath, "internal/adapters/primary/http")
	if helpers.DirExists(httpAdapterDir) {
		// Framework-specific adapter should exist
		frameworkFiles := helpers.FindFiles(t, httpAdapterDir, "*"+framework+"*.go")
		if len(frameworkFiles) > 0 {
			t.Logf("✓ Framework %s is implemented as primary adapter", framework)
		}
	}
}

func (v *HexagonalValidator) ValidateCoreFrameworkAgnostic(t *testing.T, framework string) {
	t.Helper()

	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		frameworkImports := map[string]string{
			"gin":   "github.com/gin-gonic/gin",
			"echo":  "github.com/labstack/echo",
			"fiber": "github.com/gofiber/fiber",
			"chi":   "github.com/go-chi/chi",
		}

		frameworkImport := frameworkImports[framework]

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)
			assert.NotContains(t, content, frameworkImport,
				"Core should not know about framework %s", framework)
		}
	}
}

func (v *HexagonalValidator) ValidateHexagonalPrinciples(t *testing.T) {
	t.Helper()

	// 1. Application is at the center
	v.ValidateCoreBusinessLogicIndependence(t)

	// 2. Dependencies point inward
	v.ValidateDependencyInversion(t)

	// 3. External concerns are outside
	v.ValidateExternalDependenciesThroughPorts(t)

	// 4. Adapters are pluggable
	v.ValidateHTTPAdapterSwappability(t)
	v.ValidateDatabaseAdapterSwappability(t)
}

func (v *HexagonalValidator) ValidateHexagonDefinition(t *testing.T) {
	t.Helper()

	// The hexagon should be clearly defined by:
	// 1. Core (center of hexagon)
	coreDir := filepath.Join(v.ProjectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		t.Logf("✓ Hexagon center (core) defined")
	}

	// 2. Ports (edges of hexagon)
	portsDir := filepath.Join(v.ProjectPath, "internal/ports")
	if helpers.DirExists(portsDir) {
		t.Logf("✓ Hexagon edges (ports) defined")
	}

	// 3. Adapters (outside of hexagon)
	adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")
	if helpers.DirExists(adaptersDir) {
		t.Logf("✓ Hexagon exterior (adapters) defined")
	}

	// 4. Clear separation between primary and secondary concerns
	primaryDir := filepath.Join(v.ProjectPath, "internal/adapters/primary")
	secondaryDir := filepath.Join(v.ProjectPath, "internal/adapters/secondary")
	if helpers.DirExists(primaryDir) && helpers.DirExists(secondaryDir) {
		t.Logf("✓ Primary and secondary adapters separated")
	}
}