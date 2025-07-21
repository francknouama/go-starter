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
			// Check if this directory actually contains template files
			// by looking for template.yaml files
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
				// Check if any subdirectory contains template.yaml
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
					println("Found blueprints directory at:", templatesDir)
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			break
		}
		projectRoot = parent
	}

	panic("Could not find blueprints directory with actual blueprints from working directory: " + wd)
}

// TestArchitecture_Standard_BasicGeneration tests basic standard architecture project generation
func TestArchitecture_Standard_BasicGeneration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-standard-api",
		Module:    "github.com/test/test-standard-api",
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Assert core files exist
	expectedFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/config/config.go",
		"internal/handlers/health_gin.go",
		"internal/logger/interface.go",
		"internal/logger/factory.go",
		"internal/logger/slog.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)
	helpers.AssertGoModValid(t, filepath.Join(projectPath, "go.mod"), config.Module)
	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Standard_FileStructure validates standard architecture file organization
func TestArchitecture_Standard_FileStructure(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-standard-structure",
		Module:    "github.com/test/test-standard-structure",
		Type:      "web-api-standard",
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

	// Standard architecture should have flat structure in internal/
	expectedDirs := []string{
		"internal/config",
		"internal/handlers",
		"internal/logger",
		"internal/models",
		"internal/repository",
		"internal/services",
		"internal/middleware",
		"internal/database",
	}

	for _, dir := range expectedDirs {
		helpers.AssertDirExists(t, filepath.Join(projectPath, dir))
	}

	// Should NOT have layered architecture directories
	unexpectedDirs := []string{
		"internal/domain",
		"internal/adapters",
		"internal/infrastructure",
		"internal/application",
		"internal/presentation",
	}

	for _, dir := range unexpectedDirs {
		dirPath := filepath.Join(projectPath, dir)
		_, err := helpers.GetFileInfo(dirPath)
		assert.Error(t, err, "Directory %s should not exist in standard architecture", dir)
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Standard_DatabaseIntegration tests database integration in standard architecture
func TestArchitecture_Standard_DatabaseIntegration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-standard-db",
		Module:    "github.com/test/test-standard-db",
		Type:      "web-api-standard",
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

	// Validate database-specific files
	expectedFiles := []string{
		"internal/database/connection.go",
		"internal/database/migrations.go",
		"internal/models/base.go",
		"internal/models/user.go",
		"internal/repository/interfaces.go",
		"internal/repository/user.go",
		"internal/services/user.go",
		"internal/handlers/users_gin.go",
		"migrations/001_create_users.up.sql",
		"migrations/001_create_users.down.sql",
		"migrations/embed.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

	// Validate repository pattern implementation
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/repository/interfaces.go"), "UserRepository")
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/repository/user.go"), "userRepository")
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/services/user.go"), "UserService")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Standard_LoggerVariations tests different logger implementations
func TestArchitecture_Standard_LoggerVariations(t *testing.T) {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-standard-" + logger,
				Module:       "github.com/test/test-standard-" + logger,
				Type:         "web-api",
				Architecture: "standard",
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Assert logger-specific files
			expectedFiles := []string{
				"internal/logger/interface.go",
				"internal/logger/factory.go",
				"internal/logger/" + logger + ".go",
			}

			helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

			// Assert factory creates correct logger
			helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/logger/factory.go"), logger)

			// Assert go.mod contains correct dependency
			goModPath := filepath.Join(projectPath, "go.mod")
			switch logger {
			case "zap":
				helpers.AssertFileContains(t, goModPath, "go.uber.org/zap")
			case "logrus":
				helpers.AssertFileContains(t, goModPath, "github.com/sirupsen/logrus")
			case "zerolog":
				helpers.AssertFileContains(t, goModPath, "github.com/rs/zerolog")
			case "slog":
				// slog is standard library, no external dependency
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Standard_FrameworkIntegration tests framework-specific handlers
func TestArchitecture_Standard_FrameworkIntegration(t *testing.T) {
	frameworks := []string{"gin", "echo", "fiber", "chi"}

	for _, framework := range frameworks {
		t.Run("Framework_"+framework, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "test-standard-" + framework,
				Module:       "github.com/test/test-standard-" + framework,
				Type:         "web-api",
				Architecture: "standard",
				GoVersion:    "1.21",
				Framework:    framework,
				Logger:       "slog",
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

			// Assert framework-specific handlers
			expectedFiles := []string{
				"internal/handlers/health_" + framework + ".go",
				"internal/handlers/users_" + framework + ".go",
				"internal/handlers/auth_" + framework + ".go",
			}

			helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

			// Assert go.mod contains framework dependency
			goModPath := filepath.Join(projectPath, "go.mod")
			switch framework {
			case "gin":
				helpers.AssertFileContains(t, goModPath, "github.com/gin-gonic/gin")
			case "echo":
				helpers.AssertFileContains(t, goModPath, "github.com/labstack/echo")
			case "fiber":
				helpers.AssertFileContains(t, goModPath, "github.com/gofiber/fiber")
			case "chi":
				helpers.AssertFileContains(t, goModPath, "github.com/go-chi/chi")
			}

			helpers.AssertCompilable(t, projectPath)
		})
	}
}

// TestArchitecture_Standard_AuthenticationIntegration tests authentication features
func TestArchitecture_Standard_AuthenticationIntegration(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-standard-auth",
		Module:    "github.com/test/test-standard-auth",
		Type:      "web-api-standard",
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

	// Assert authentication files
	expectedFiles := []string{
		"internal/handlers/auth_gin.go",
		"internal/services/auth.go",
		"internal/middleware/auth.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

	// Assert JWT dependency
	helpers.AssertFileContains(t, filepath.Join(projectPath, "go.mod"), "github.com/golang-jwt/jwt")

	// Assert authentication service
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/services/auth.go"), "AuthService")

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Standard_WithoutOptionalFeatures tests minimal standard architecture
func TestArchitecture_Standard_WithoutOptionalFeatures(t *testing.T) {
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

	// Assert only basic files exist
	expectedFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
		"Makefile",
		"cmd/server/main.go",
		"internal/config/config.go",
		"internal/handlers/health_gin.go",
		"internal/logger/interface.go",
		"internal/logger/factory.go",
		"internal/logger/slog.go",
	}

	helpers.AssertProjectGenerated(t, projectPath, expectedFiles)

	// Assert database files don't exist
	unexpectedFiles := []string{
		"internal/database/connection.go",
		"internal/models/user.go",
		"internal/repository/user.go",
		"internal/handlers/users_gin.go",
		"internal/handlers/auth_gin.go",
		"migrations/001_create_users.up.sql",
	}

	for _, file := range unexpectedFiles {
		filePath := filepath.Join(projectPath, file)
		_, err := helpers.GetFileInfo(filePath)
		assert.Error(t, err, "File %s should not exist in minimal configuration", file)
	}

	helpers.AssertCompilable(t, projectPath)
}

// TestArchitecture_Standard_DependencyManagement tests dependency flow in standard architecture
func TestArchitecture_Standard_DependencyManagement(t *testing.T) {
	config := types.ProjectConfig{
		Name:      "test-standard-deps",
		Module:    "github.com/test/test-standard-deps",
		Type:      "web-api-standard",
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

	// In standard architecture, dependencies are straightforward:
	// handlers -> services -> repository -> database

	// Assert handlers import services
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/handlers/users_gin.go"), "services")

	// Assert services import repository
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/services/user.go"), "repository")

	// Assert repository imports models
	helpers.AssertFileContains(t, filepath.Join(projectPath, "internal/repository/user.go"), "models")

	helpers.AssertCompilable(t, projectPath)
}
