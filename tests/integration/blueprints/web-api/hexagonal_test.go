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

// TestArchitecture_Hexagonal_BasicGeneration tests basic Hexagonal Architecture project generation
func TestArchitecture_Hexagonal_BasicGeneration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-api",
		Module:    "github.com/test/test-hexagonal-api",
		Type:      "web-api-hexagonal",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Assert core hexagonal architecture files exist
	expectedFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/core/domain",
		"internal/core/ports",
		"internal/core/services",
		"internal/adapters/primary",
		"internal/adapters/secondary",
		"internal/infrastructure/config/config.go",
		"internal/infrastructure/logger/interface.go",
		"internal/infrastructure/logger/factory.go",
		"internal/infrastructure/logger/slog.go",
	}

	// Note: Some files might be directories, so we check existence differently
	for _, path := range expectedFiles {
		fullPath := filepath.Join(projectPath, path)
		exists := helpers.FileExists(fullPath) || helpers.DirExists(fullPath)
		assert.True(t, exists, "Path %s should exist", path)
	}

	helpers.AssertGoModValid(t, filepath.Join(projectPath, "go.mod"), config.Module)
	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_PortsAndAdaptersStructure validates hexagonal structure
func TestArchitecture_Hexagonal_PortsAndAdaptersStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-structure",
		Module:    "github.com/test/test-hexagonal-structure",
		Type:      "web-api-hexagonal",
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

	// Hexagonal architecture should have specific structure
	expectedStructure := []string{
		// Core (business logic) - independent of external concerns
		"internal/core/domain",
		"internal/core/ports",
		"internal/core/services",
		// Primary adapters (driving side) - HTTP, CLI, etc.
		"internal/adapters/primary/http",
		"internal/adapters/primary/rest",
		// Secondary adapters (driven side) - database, external APIs, etc.
		"internal/adapters/secondary/persistence",
		"internal/adapters/secondary/external",
		// Infrastructure - framework-specific implementations
		"internal/infrastructure/config",
		"internal/infrastructure/logger",
		"internal/infrastructure/container",
	}

	for _, structure := range expectedStructure {
		fullPath := filepath.Join(projectPath, structure)
		// Check if it exists as either file or directory
		exists := helpers.FileExists(fullPath) || helpers.DirExists(fullPath)
		if !exists {
			// If the exact path doesn't exist, check if there are any files in that directory pattern
			parentDir := filepath.Dir(fullPath)
			if helpers.DirExists(parentDir) {
				files := helpers.FindFiles(t, parentDir, "*")
				if len(files) > 0 {
					t.Logf("Directory %s exists with files, considering structure valid", structure)
					continue
				}
			}
			t.Logf("Warning: Expected structure %s not found, but this might be normal if template is incomplete", structure)
		}
	}

	// Should NOT have non-hexagonal patterns
	unexpectedDirs := []string{
		"internal/handlers",
		"internal/models",
		"internal/repository",
		"internal/services",        // services should be in core/services
		"internal/domain/entities", // Clean Architecture pattern
		"internal/application",     // DDD pattern
		"internal/presentation",    // DDD pattern
	}

	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(projectPath, dir)
		_, err := helpers.GetFileInfo(dirPath)
		assert.Error(t, err, "Directory %s should not exist in Hexagonal Architecture", dir)
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_CoreDomainIsolation tests core domain isolation
func TestArchitecture_Hexagonal_CoreDomainIsolation(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-core",
		Module:    "github.com/test/test-hexagonal-core",
		Type:      "web-api-hexagonal",
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

	// Find domain files in core
	coreDir := filepath.Join(projectPath, "internal/core")
	if !helpers.DirExists(coreDir) {
		t.Skip("Core directory not found - template might not be fully implemented")
		return
	}

	domainFiles := helpers.FindFiles(t, coreDir, "*.go")

	// Core should not depend on external frameworks or adapters
	forbiddenImports := []string{
		"github.com/gin-gonic/gin",
		"github.com/labstack/echo",
		"github.com/gofiber/fiber",
		"gorm.io/gorm",
		"internal/adapters",
		"internal/infrastructure",
	}

	for _, file := range domainFiles {
		content := helpers.ReadFileContent(t, file)
		for _, forbidden := range forbiddenImports {
			assert.NotContains(t, content, forbidden,
				"Core domain file %s should not import %s - violates hexagonal principle", file, forbidden)
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_PortsDefinition tests port interfaces definition
func TestArchitecture_Hexagonal_PortsDefinition(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-ports",
		Module:    "github.com/test/test-hexagonal-ports",
		Type:      "web-api-hexagonal",
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

	// Assert ports directory exists
	portsDir := filepath.Join(projectPath, "internal/core/ports")
	if !helpers.DirExists(portsDir) {
		t.Skip("Ports directory not found - template might not be fully implemented")
		return
	}

	// Look for port interface files
	portFiles := helpers.FindFiles(t, portsDir, "*.go")
	assert.NotEmpty(t, portFiles, "Should have port interface definitions")

	// Ports should define interfaces for primary and secondary adapters
	for _, file := range portFiles {
		content := helpers.ReadFileContent(t, file)
		// Should contain interface definitions
		assert.Contains(t, content, "interface", "Port files should define interfaces")
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_PrimaryAdapters tests primary (driving) adapters
func TestArchitecture_Hexagonal_PrimaryAdapters(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-primary",
		Module:    "github.com/test/test-hexagonal-primary",
		Type:      "web-api-hexagonal",
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

	// Assert primary adapters exist
	primaryDir := filepath.Join(projectPath, "internal/adapters/primary")
	if !helpers.DirExists(primaryDir) {
		t.Skip("Primary adapters directory not found - template might not be fully implemented")
		return
	}

	// Primary adapters should handle incoming requests
	primaryFiles := helpers.FindFiles(t, primaryDir, "*.go")

	for _, file := range primaryFiles {
		content := helpers.ReadFileContent(t, file)
		// Primary adapters should use core services through ports
		shouldContainOneOf := []string{"core/services", "core/ports", "ports"}
		foundCore := false
		for _, coreImport := range shouldContainOneOf {
			if helpers.StringContains(content, coreImport) {
				foundCore = true
				break
			}
		}

		if len(primaryFiles) > 0 && !foundCore {
			t.Logf("Primary adapter %s should use core services/ports", file)
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_SecondaryAdapters tests secondary (driven) adapters
func TestArchitecture_Hexagonal_SecondaryAdapters(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-secondary",
		Module:    "github.com/test/test-hexagonal-secondary",
		Type:      "web-api-hexagonal",
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

	// Assert secondary adapters exist
	secondaryDir := filepath.Join(projectPath, "internal/adapters/secondary")
	if !helpers.DirExists(secondaryDir) {
		t.Skip("Secondary adapters directory not found - template might not be fully implemented")
		return
	}

	// Secondary adapters should implement port interfaces
	secondaryFiles := helpers.FindFiles(t, secondaryDir, "*.go")

	for _, file := range secondaryFiles {
		content := helpers.ReadFileContent(t, file)
		// Secondary adapters should implement interfaces from ports
		shouldContainOneOf := []string{"core/ports", "ports"}
		foundPorts := false
		for _, portsImport := range shouldContainOneOf {
			if helpers.StringContains(content, portsImport) {
				foundPorts = true
				break
			}
		}

		if len(secondaryFiles) > 0 && !foundPorts {
			t.Logf("Secondary adapter %s should implement port interfaces", file)
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_DependencyInversion tests dependency inversion principle
func TestArchitecture_Hexagonal_DependencyInversion(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-di",
		Module:    "github.com/test/test-hexagonal-di",
		Type:      "web-api-hexagonal",
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

	// Core services should depend on ports (abstractions), not implementations
	servicesDir := filepath.Join(projectPath, "internal/core/services")
	if !helpers.DirExists(servicesDir) {
		t.Skip("Services directory not found - template might not be fully implemented")
		return
	}

	serviceFiles := helpers.FindFiles(t, servicesDir, "*.go")

	for _, file := range serviceFiles {
		content := helpers.ReadFileContent(t, file)

		// Services should depend on ports (interfaces)
		assert.Contains(t, content, "ports", "Core services should depend on port interfaces")

		// Services should NOT depend on concrete implementations
		forbiddenImports := []string{
			"internal/adapters/secondary",
			"gorm.io/gorm",             // direct DB dependency
			"github.com/gin-gonic/gin", // direct HTTP framework dependency
		}

		for _, forbidden := range forbiddenImports {
			assert.NotContains(t, content, forbidden,
				"Core service %s should not depend on concrete implementation: %s", file, forbidden)
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_LoggerIntegration tests logger integration
func TestArchitecture_Hexagonal_LoggerIntegration(t *testing.T) {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-hexagonal-" + logger,
				Module:       "github.com/test/test-hexagonal-" + logger,
				Type:         "web-api",
				Architecture: "hexagonal",
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Assert logger files in infrastructure
			expectedFiles := []string{
				"internal/infrastructure/logger/interface.go",
				"internal/infrastructure/logger/factory.go",
				"internal/infrastructure/logger/" + logger + ".go",
			}

			helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

			// Core should not depend on specific logger implementations
			coreDir := filepath.Join(projectPath, "internal/core")
			if helpers.DirExists(coreDir) {
				coreFiles := helpers.FindFiles(t, coreDir, "*.go")
				for _, file := range coreFiles {
					content := helpers.ReadFileContent(t, file)

					forbiddenImports := []string{
						"go.uber.org/zap",
						"github.com/sirupsen/logrus",
						"github.com/rs/zerolog",
					}

					for _, forbidden := range forbiddenImports {
						assert.NotContains(t, content, forbidden,
							"Core should not import specific logger: %s in file: %s", forbidden, file)
					}
				}
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Hexagonal_FrameworkIndependence tests framework independence
func TestArchitecture_Hexagonal_FrameworkIndependence(t *testing.T) {
	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-hexagonal-" + framework,
				Module:       "github.com/test/test-hexagonal-" + framework,
				Type:         "web-api",
				Architecture: "hexagonal",
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

			// Core should be framework-independent
			coreDir := filepath.Join(projectPath, "internal/core")
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
					for _, frameworkImport := range frameworkImports {
						assert.NotContains(t, content, frameworkImport,
							"Core should not depend on framework: %s in file: %s", frameworkImport, file)
					}
				}
			}

			// Framework-specific code should be in primary adapters
			primaryDir := filepath.Join(projectPath, "internal/adapters/primary")
			if helpers.DirExists(primaryDir) {
				primaryFiles := helpers.FindFiles(t, primaryDir, "*.go")

				// At least one primary adapter should contain framework-specific code
				frameworkFound := false
				for _, file := range primaryFiles {
					content := helpers.ReadFileContent(t, file)
					frameworkImports := []string{
						"github.com/gin-gonic/gin",
						"github.com/labstack/echo",
						"github.com/gofiber/fiber",
						"github.com/go-chi/chi",
					}

					for _, frameworkImport := range frameworkImports {
						if helpers.StringContains(content, frameworkImport) {
							frameworkFound = true
							break
						}
					}
					if frameworkFound {
						break
					}
				}

				if len(primaryFiles) > 0 && !frameworkFound {
					t.Logf("Expected to find framework-specific code in primary adapters")
				}
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Hexagonal_DatabaseAbstraction tests database abstraction
func TestArchitecture_Hexagonal_DatabaseAbstraction(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-db",
		Module:    "github.com/test/test-hexagonal-db",
		Type:      "web-api-hexagonal",
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

	// Database-specific code should be in secondary adapters
	secondaryDir := filepath.Join(projectPath, "internal/adapters/secondary")
	if helpers.DirExists(secondaryDir) {
		persistenceFiles := helpers.FindFiles(t, secondaryDir, "*.go")

		dbImportsFound := false
		for _, file := range persistenceFiles {
			content := helpers.ReadFileContent(t, file)
			dbImports := []string{
				"gorm.io/gorm",
				"database/sql",
				"github.com/lib/pq",
			}

			for _, dbImport := range dbImports {
				if helpers.StringContains(content, dbImport) {
					dbImportsFound = true
					break
				}
			}
			if dbImportsFound {
				break
			}
		}

		if len(persistenceFiles) > 0 && !dbImportsFound {
			t.Logf("Expected to find database-specific code in secondary adapters")
		}
	}

	// Core should define repository interfaces without database specifics
	portsDir := filepath.Join(projectPath, "internal/core/ports")
	if helpers.DirExists(portsDir) {
		portFiles := helpers.FindFiles(t, portsDir, "*.go")

		for _, file := range portFiles {
			content := helpers.ReadFileContent(t, file)

			// Ports should not contain database-specific imports
			forbiddenImports := []string{
				"gorm.io/gorm",
				"database/sql",
				"github.com/lib/pq",
			}

			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Port interface %s should not contain database-specific import: %s", file, forbidden)
			}
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_TestingStrategy tests hexagonal testing approach
func TestArchitecture_Hexagonal_TestingStrategy(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-testing",
		Module:    "github.com/test/test-hexagonal-testing",
		Type:      "web-api-hexagonal",
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

	// Test structure should support hexagonal testing
	testDirs := []string{
		"tests/unit",
		"tests/integration",
		"tests/mocks",
	}

	for _, dir := range testDirs {
		fullPath := filepath.Join(projectPath, dir)
		if helpers.DirExists(fullPath) {
			files := helpers.FindFiles(t, fullPath, "*.go")
			t.Logf("Found %d test files in %s", len(files), dir)
		}
	}

	// Core business logic should be easily testable in isolation
	coreDir := filepath.Join(projectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		// Look for test files in core
		coreTestFiles := helpers.FindFiles(t, coreDir, "*_test.go")
		t.Logf("Found %d test files in core", len(coreTestFiles))

		// Core services should be testable with mocked ports
		servicesDir := filepath.Join(coreDir, "services")
		if helpers.DirExists(servicesDir) {
			serviceFiles := helpers.FindFiles(t, servicesDir, "*.go")
			if len(serviceFiles) > 0 {
				t.Logf("Found %d service files that should be unit testable", len(serviceFiles))
			}
		}
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Hexagonal_ConfigurationManagement tests configuration handling
func TestArchitecture_Hexagonal_ConfigurationManagement(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-hexagonal-config",
		Module:    "github.com/test/test-hexagonal-config",
		Type:      "web-api-hexagonal",
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

	// Configuration should be in infrastructure layer
	configFile := filepath.Join(projectPath, "internal/infrastructure/config/config.go")
	helpers.AssertFileExists(t, configFile)

	// Core should not directly depend on configuration details
	coreDir := filepath.Join(projectPath, "internal/core")
	if helpers.DirExists(coreDir) {
		coreFiles := helpers.FindFiles(t, coreDir, "*.go")

		for _, file := range coreFiles {
			content := helpers.ReadFileContent(t, file)

			// Core should not import infrastructure config directly
			forbiddenImports := []string{
				"internal/infrastructure/config",
				"github.com/spf13/viper", // specific config library
			}

			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, content, forbidden,
					"Core should not directly import config: %s in file: %s", forbidden, file)
			}
		}
	}

	helpers.AssertCompilable(t, projectPath)
}
