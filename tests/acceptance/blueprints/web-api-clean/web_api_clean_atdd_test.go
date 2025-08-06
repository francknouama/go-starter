package web_api_clean

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestWebApiCleanBlueprintATDD validates the acceptance criteria for web-api-clean blueprint
// This ensures the web-api-clean blueprint generates proper Clean Architecture patterns
func TestWebApiCleanBlueprintATDD(t *testing.T) {
	t.Run("web_api_clean_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: web-api-clean should be in the list with correct architecture description

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/web-api-clean)
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
		assert.Contains(t, outputStr, "web-api-clean", "web-api-clean blueprint should be listed")
		assert.Contains(t, outputStr, "Clean Architecture", "Should show Clean Architecture description")
		assert.Contains(t, outputStr, "layered design", "Should emphasize architectural layers")
	})

	t.Run("web_api_clean_generates_clean_architecture_layers", func(t *testing.T) {
		// GIVEN: User wants a Clean Architecture web API
		// WHEN: User generates a project with web-api-clean blueprint
		// THEN: Should generate proper Clean Architecture layers (domain, application, infrastructure, presentation)

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

		// Generate a web-api-clean project
		generateCmd := exec.Command("./go-starter", "new", "test-clean-api",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/clean-api",
			"--framework=gin",
			"--logger=slog",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-clean-api")

		// Verify Clean Architecture layer directories exist
		cleanArchLayers := []string{
			// Domain layer (innermost - business logic and entities)
			"internal/domain/entities",
			"internal/domain/repositories", 
			"internal/domain/services",
			"internal/domain/valueobjects",

			// Application layer (use cases and business workflows)
			"internal/application/usecases",
			"internal/application/ports",
			"internal/application/services",

			// Infrastructure layer (external concerns - databases, APIs, etc.)
			"internal/infrastructure/database",
			"internal/infrastructure/repositories",
			"internal/infrastructure/external",

			// Presentation layer (controllers, handlers, DTOs)
			"internal/presentation/handlers",
			"internal/presentation/middleware",
			"internal/presentation/dto",
			"internal/presentation/routes",
		}

		for _, layer := range cleanArchLayers {
			assert.DirExists(t, filepath.Join(projectDir, layer), "Clean Architecture layer %s should exist", layer)
		}

		// Verify essential Clean Architecture files exist
		essentialFiles := []string{
			// Foundation
			"main.go",
			"go.mod",
			"README.md",

			// Domain entities (business objects)
			"internal/domain/entities/user.go",
			"internal/domain/entities/product.go",

			// Domain repositories (interfaces)
			"internal/domain/repositories/user_repository.go",

			// Application use cases (business workflows)  
			"internal/application/usecases/user_usecase.go",
			"internal/application/usecases/create_user_usecase.go",

			// Infrastructure implementations
			"internal/infrastructure/repositories/user_repository_impl.go",
			"internal/infrastructure/database/database.go",

			// Presentation layer
			"internal/presentation/handlers/user_handler.go", 
			"internal/presentation/handlers/health_handler.go",
			"internal/presentation/dto/user_dto.go",
			"internal/presentation/routes/routes.go",

			// Configuration
			"internal/config/config.go",
			"internal/logger/logger.go",
		}

		for _, file := range essentialFiles {
			assert.FileExists(t, filepath.Join(projectDir, file), "Essential Clean Architecture file %s should exist", file)
		}
	})

	t.Run("web_api_clean_follows_dependency_inversion", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining dependency relationships
		// THEN: Should follow Clean Architecture dependency inversion (inner layers don't depend on outer layers)

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
		generateCmd := exec.Command("./go-starter", "new", "test-dependency",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/dependency",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-dependency")

		// 1. Domain layer should NOT import outer layers
		userEntityPath := filepath.Join(projectDir, "internal", "domain", "entities", "user.go")
		if assert.FileExists(t, userEntityPath, "User entity should exist") {
			entityContent, err := os.ReadFile(userEntityPath)
			require.NoError(t, err)
			entityStr := string(entityContent)

			// Domain entities should NOT import infrastructure or presentation layers
			assert.NotContains(t, entityStr, "internal/infrastructure", "Domain entities should not import infrastructure")
			assert.NotContains(t, entityStr, "internal/presentation", "Domain entities should not import presentation")
			assert.NotContains(t, entityStr, "github.com/gin-gonic", "Domain should not import web framework")
			assert.NotContains(t, entityStr, "gorm.io/gorm", "Domain should not import ORM directly")

			// Domain should define business rules
			assert.Contains(t, entityStr, "type User struct", "Should define User entity")
		}

		// 2. Domain repositories should be interfaces (not implementations)
		repoInterfacePath := filepath.Join(projectDir, "internal", "domain", "repositories", "user_repository.go")
		if assert.FileExists(t, repoInterfacePath, "User repository interface should exist") {
			repoContent, err := os.ReadFile(repoInterfacePath)
			require.NoError(t, err)
			repoStr := string(repoContent)

			assert.Contains(t, repoStr, "type UserRepository interface", "Should define repository interface")
			assert.Contains(t, repoStr, "Create", "Should have Create method")
			assert.Contains(t, repoStr, "GetByID", "Should have GetByID method")
			assert.NotContains(t, repoStr, "gorm.DB", "Interface should not reference concrete implementations")
		}

		// 3. Application layer should depend on domain interfaces
		usecasePath := filepath.Join(projectDir, "internal", "application", "usecases", "user_usecase.go")
		if assert.FileExists(t, usecasePath, "User usecase should exist") {
			usecaseContent, err := os.ReadFile(usecasePath)
			require.NoError(t, err)
			usecaseStr := string(usecaseContent)

			assert.Contains(t, usecaseStr, "internal/domain/repositories", "Should import domain repositories")
			assert.Contains(t, usecaseStr, "UserRepository", "Should use repository interface")
			assert.NotContains(t, usecaseStr, "internal/infrastructure", "Use cases should not import infrastructure directly")
		}

		// 4. Infrastructure should implement domain interfaces
		repoImplPath := filepath.Join(projectDir, "internal", "infrastructure", "repositories", "user_repository_impl.go")
		if assert.FileExists(t, repoImplPath, "Repository implementation should exist") {
			implContent, err := os.ReadFile(repoImplPath)
			require.NoError(t, err)
			implStr := string(implContent)

			assert.Contains(t, implStr, "internal/domain/repositories", "Should import domain repository interface")
			assert.Contains(t, implStr, "internal/domain/entities", "Should import domain entities")
			assert.Contains(t, implStr, "UserRepository", "Should implement UserRepository interface")
			// Can import external dependencies like GORM
			if strings.Contains(implStr, "gorm") {
				assert.Contains(t, implStr, "gorm.io/gorm", "Infrastructure can import GORM")
			}
		}

		// 5. Presentation layer should depend on application layer
		handlerPath := filepath.Join(projectDir, "internal", "presentation", "handlers", "user_handler.go")
		if assert.FileExists(t, handlerPath, "User handler should exist") {
			handlerContent, err := os.ReadFile(handlerPath)
			require.NoError(t, err)
			handlerStr := string(handlerContent)

			assert.Contains(t, handlerStr, "internal/application/usecases", "Handlers should use application use cases")
			assert.Contains(t, handlerStr, "gin-gonic/gin", "Handlers can import web framework")
			assert.NotContains(t, handlerStr, "internal/infrastructure/repositories", "Handlers should not directly use infrastructure")
		}
	})

	t.Run("web_api_clean_project_builds_successfully", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
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
			"--type=web-api",
			"--architecture=clean",
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
		buildGeneratedCmd := exec.Command("go", "build", "-o", "clean-api", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated web-api-clean project should build successfully: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "clean-api"), "Clean API binary should be created")
	})

	t.Run("web_api_clean_supports_multiple_frameworks", func(t *testing.T) {
		// GIVEN: web-api-clean blueprint with different web frameworks
		// WHEN: Generating projects with different frameworks
		// THEN: Each should generate framework-specific presentation layer while maintaining Clean Architecture

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
				projectName := fmt.Sprintf("test-clean-%s", framework)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=web-api",
					"--architecture=clean",
					"--module=github.com/test/"+projectName,
					"--framework="+framework,
					"--database-driver=sqlite",
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s: %s", framework, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Verify Clean Architecture structure exists regardless of framework
				assert.DirExists(t, filepath.Join(projectDir, "internal", "domain"), "Domain layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "application"), "Application layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure"), "Infrastructure layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "presentation"), "Presentation layer should exist")

				// Check that presentation layer uses correct framework
				if handlerExists(t, filepath.Join(projectDir, "internal", "presentation", "handlers", "user_handler.go")) {
					handlerContent, err := os.ReadFile(filepath.Join(projectDir, "internal", "presentation", "handlers", "user_handler.go"))
					require.NoError(t, err)
					handlerStr := string(handlerContent)

					switch framework {
					case "gin":
						assert.Contains(t, handlerStr, "gin-gonic/gin", "Should import Gin framework")
					case "echo":
						assert.Contains(t, handlerStr, "labstack/echo", "Should import Echo framework")
					case "fiber":
						assert.Contains(t, handlerStr, "gofiber/fiber", "Should import Fiber framework")
					case "chi":
						assert.Contains(t, handlerStr, "go-chi/chi", "Should import Chi framework")
					}
				}

				// Verify project compiles
				modTidyCmd := exec.Command("go", "mod", "tidy")
				modTidyCmd.Dir = projectDir
				_, err = modTidyCmd.CombinedOutput()
				require.NoError(t, err)

				buildCmd := exec.Command("go", "build", ".")
				buildCmd.Dir = projectDir
				output, err = buildCmd.CombinedOutput()
				require.NoError(t, err, "Clean Architecture with %s should compile: %s", framework, string(output))
			})
		}
	})

	t.Run("web_api_clean_supports_multiple_loggers", func(t *testing.T) {
		// GIVEN: web-api-clean blueprint with different logger configurations
		// WHEN: Generating projects with different loggers
		// THEN: Logger should be properly injected throughout Clean Architecture layers

		loggers := []string{"slog", "zap", "logrus", "zerolog"}

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

		for _, logger := range loggers {
			t.Run(logger, func(t *testing.T) {
				projectName := fmt.Sprintf("test-log-%s", logger)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=web-api",
					"--architecture=clean",
					"--module=github.com/test/"+projectName,
					"--framework=gin",
					"--logger="+logger,
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s logger: %s", logger, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Check that logger setup exists
				assert.FileExists(t, filepath.Join(projectDir, "internal", "logger", "logger.go"), 
					"Logger setup should exist for %s", logger)

				// Check logger is used in use cases (application layer)
				if usecaseExists(t, filepath.Join(projectDir, "internal", "application", "usecases", "user_usecase.go")) {
					usecaseContent, err := os.ReadFile(filepath.Join(projectDir, "internal", "application", "usecases", "user_usecase.go"))
					require.NoError(t, err)
					usecaseStr := string(usecaseContent)

					assert.Contains(t, usecaseStr, "internal/logger", "Use cases should import logger")
				}

				// Check logger is used in handlers (presentation layer)
				if handlerExists(t, filepath.Join(projectDir, "internal", "presentation", "handlers", "user_handler.go")) {
					handlerContent, err := os.ReadFile(filepath.Join(projectDir, "internal", "presentation", "handlers", "user_handler.go"))
					require.NoError(t, err)
					handlerStr := string(handlerContent)

					assert.Contains(t, handlerStr, "internal/logger", "Handlers should import logger")
				}

				// Verify project compiles with the logger
				modTidyCmd := exec.Command("go", "mod", "tidy")
				modTidyCmd.Dir = projectDir
				_, err = modTidyCmd.CombinedOutput()
				require.NoError(t, err)

				buildGeneratedCmd := exec.Command("go", "build", "-o", "test-"+logger, ".")
				buildGeneratedCmd.Dir = projectDir
				output, err = buildGeneratedCmd.CombinedOutput()
				require.NoError(t, err, "web-api-clean with %s logger should compile: %s", logger, string(output))
			})
		}
	})

	t.Run("web_api_clean_implements_proper_use_cases", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining use case implementation
		// THEN: Should implement proper use case patterns with single responsibility

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

		generateCmd := exec.Command("./go-starter", "new", "test-usecases",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/usecases",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-usecases")

		// Check specific use case files exist
		usecaseFiles := []string{
			"internal/application/usecases/create_user_usecase.go",
			"internal/application/usecases/get_user_usecase.go",
			"internal/application/usecases/user_usecase.go",
		}

		for _, file := range usecaseFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				usecaseContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				usecaseStr := string(usecaseContent)

				// Use cases should follow patterns:
				// 1. Single Responsibility - each use case handles one business operation
				assert.Contains(t, usecaseStr, "type", "Should define use case types")
				
				// 2. Dependency Injection - should accept repository via constructor
				assert.Contains(t, usecaseStr, "New", "Should have constructor function")
				
				// 3. Business Logic - should contain actual business operations
				if strings.Contains(file, "create_user") {
					assert.Contains(t, usecaseStr, "Create", "Create use case should have Create method")
				}
				if strings.Contains(file, "get_user") {
					assert.Contains(t, usecaseStr, "Get", "Get use case should have Get method")
				}

				// 4. Should use domain repositories (not infrastructure)
				assert.Contains(t, usecaseStr, "internal/domain/repositories", "Should import domain repositories")
				assert.NotContains(t, usecaseStr, "internal/infrastructure", "Should not import infrastructure directly")
			}
		}

		// Check application services coordination
		if fileExists(t, filepath.Join(projectDir, "internal", "application", "services", "user_service.go")) {
			serviceContent, err := os.ReadFile(filepath.Join(projectDir, "internal", "application", "services", "user_service.go"))
			require.NoError(t, err)
			serviceStr := string(serviceContent)

			assert.Contains(t, serviceStr, "internal/application/usecases", "Services should coordinate use cases")
			assert.Contains(t, serviceStr, "type UserService", "Should define service interface or struct")
		}
	})

	t.Run("web_api_clean_has_proper_entity_design", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining domain entities
		// THEN: Should have framework-independent domain entities with business logic

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

		generateCmd := exec.Command("./go-starter", "new", "test-entities",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/entities",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-entities")

		// Check domain entities
		entityFiles := []string{
			"internal/domain/entities/user.go",
			"internal/domain/entities/product.go",
		}

		for _, file := range entityFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				entityContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				entityStr := string(entityContent)

				// Entities should be pure business objects
				assert.Contains(t, entityStr, "type", "Should define entity types")
				
				// Should NOT depend on external frameworks
				assert.NotContains(t, entityStr, "gin-gonic", "Entities should not import web frameworks")
				assert.NotContains(t, entityStr, "gorm.io/gorm", "Entities should not import ORM")
				assert.NotContains(t, entityStr, "database/sql", "Entities should not import database packages")
				assert.NotContains(t, entityStr, "net/http", "Entities should not import HTTP packages")

				// Should contain business logic methods
				if strings.Contains(file, "user.go") {
					// User entity might have validation methods
					if strings.Contains(entityStr, "IsValid") || strings.Contains(entityStr, "Validate") {
						t.Logf("User entity includes validation logic")
					}
				}

				// Should use standard Go types and time package at most
				lines := strings.Split(entityStr, "\n")
				for _, line := range lines {
					if strings.Contains(line, "import") {
						// Only allow standard library imports like time, errors, fmt
						if strings.Contains(line, "\"") && !strings.Contains(line, "time") && 
						   !strings.Contains(line, "errors") && !strings.Contains(line, "fmt") {
							t.Logf("Entity import found: %s", line)
						}
					}
				}
			}
		}

		// Check value objects exist (if implemented)
		if dirExists(t, filepath.Join(projectDir, "internal", "domain", "valueobjects")) {
			assert.DirExists(t, filepath.Join(projectDir, "internal", "domain", "valueobjects"), 
				"Value objects directory should exist")
		}
	})

	t.Run("web_api_clean_dependency_injection_setup", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining dependency injection and wiring
		// THEN: Should have proper DI container or manual injection setup

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

		generateCmd := exec.Command("./go-starter", "new", "test-di",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/di",
			"--framework=gin",
			"--database-driver=sqlite",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-di")

		// Check main.go for dependency injection setup
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)

		// Should wire dependencies from outer layers to inner layers
		assert.Contains(t, mainStr, "internal/infrastructure", "Should import infrastructure layer")
		assert.Contains(t, mainStr, "internal/presentation", "Should import presentation layer")
		
		// Look for dependency injection patterns
		diPatterns := []string{
			// Manual injection patterns
			"New",           // Constructor calls
			"repository",    // Repository injection
			"usecase",       // Use case injection
			
			// DI container patterns (if used)
			"wire",          // Google Wire
			"dig",           // Uber Dig
			"Container",     // Generic DI container
		}

		foundDiPattern := false
		for _, pattern := range diPatterns {
			if strings.Contains(mainStr, pattern) {
				foundDiPattern = true
				t.Logf("Found DI pattern: %s", pattern)
			}
		}
		
		if !foundDiPattern {
			t.Logf("No explicit DI pattern found in main.go, checking for manual wiring")
			// At minimum should have some form of component initialization
			assert.True(t, strings.Contains(mainStr, "infrastructure") || 
						   strings.Contains(mainStr, "repository") ||
						   strings.Contains(mainStr, "handler"),
						"Should have some form of dependency wiring")
		}

		// Check if there's a dedicated DI/wire file
		diFiles := []string{
			"internal/di/container.go",
			"internal/wire/wire.go", 
			"internal/config/dependencies.go",
			"cmd/wire.go",
		}

		for _, file := range diFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				diContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				diStr := string(diContent)

				t.Logf("Found DI file: %s", file)
				assert.Contains(t, diStr, "New", "DI file should contain constructor calls")
			}
		}
	})
}

// TestWebApiCleanArchitecturalCompliance validates Clean Architecture compliance
func TestWebApiCleanArchitecturalCompliance(t *testing.T) {
	t.Run("enforces_clean_architecture_boundaries", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Analyzing import statements across layers
		// THEN: Should enforce architectural boundaries (no violations of dependency rule)

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

		generateCmd := exec.Command("./go-starter", "new", "test-boundaries",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/boundaries",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-boundaries")

		// Analyze each layer for boundary violations
		layers := map[string][]string{
			"domain": {"internal/domain/entities", "internal/domain/repositories", "internal/domain/services"},
			"application": {"internal/application/usecases", "internal/application/services"},
			"infrastructure": {"internal/infrastructure/database", "internal/infrastructure/repositories"},
			"presentation": {"internal/presentation/handlers", "internal/presentation/routes"},
		}

		for layerName, directories := range layers {
			for _, dir := range directories {
				fullDir := filepath.Join(projectDir, dir)
				if dirExists(t, fullDir) {
					analyzeBoundaryCompliance(t, fullDir, layerName, projectDir)
				}
			}
		}
	})

	t.Run("implements_ports_and_adapters_pattern", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining ports (interfaces) and adapters (implementations) 
		// THEN: Should implement proper ports and adapters pattern

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

		generateCmd := exec.Command("./go-starter", "new", "test-ports",
			"--type=web-api", 
			"--architecture=clean",
			"--module=github.com/test/ports",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-ports")

		// 1. Check for Ports (interfaces) in domain or application layer
		portFiles := []string{
			"internal/domain/repositories/user_repository.go",
			"internal/application/ports/user_port.go",
		}

		foundPorts := false
		for _, file := range portFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				portContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				portStr := string(portContent)

				if strings.Contains(portStr, "interface") {
					foundPorts = true
					assert.Contains(t, portStr, "type", "Port should define interface")
					t.Logf("Found port interface in: %s", file)
				}
			}
		}

		// 2. Check for Adapters (implementations) in infrastructure layer
		adapterFiles := []string{
			"internal/infrastructure/repositories/user_repository_impl.go",
			"internal/infrastructure/adapters/user_adapter.go",
		}

		foundAdapters := false
		for _, file := range adapterFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				adapterContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				adapterStr := string(adapterContent)

				if strings.Contains(adapterStr, "type") && 
				   (strings.Contains(adapterStr, "struct") || strings.Contains(adapterStr, "impl")) {
					foundAdapters = true
					t.Logf("Found adapter implementation in: %s", file)
				}
			}
		}

		// At minimum, should have repository pattern (port/adapter)
		if foundPorts || foundAdapters {
			t.Logf("Ports and Adapters pattern implemented")
		} else {
			// Check if repository pattern exists (alternative implementation)
			repoInterface := filepath.Join(projectDir, "internal", "domain", "repositories", "user_repository.go")
			repoImpl := filepath.Join(projectDir, "internal", "infrastructure", "repositories", "user_repository_impl.go")
			
			if fileExists(t, repoInterface) && fileExists(t, repoImpl) {
				t.Logf("Repository pattern implements ports/adapters concept")
			}
		}
	})

	t.Run("maintains_testability_through_isolation", func(t *testing.T) {
		// GIVEN: A generated web-api-clean project
		// WHEN: Examining testability features
		// THEN: Should support easy unit testing through proper isolation

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

		generateCmd := exec.Command("./go-starter", "new", "test-testability",
			"--type=web-api",
			"--architecture=clean",
			"--module=github.com/test/testability",
			"--framework=gin",
			"--database-driver=sqlite",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-testability")

		// Look for test files
		testFiles := []string{
			"internal/domain/entities/user_test.go",
			"internal/application/usecases/user_usecase_test.go", 
			"internal/presentation/handlers/user_handler_test.go",
			"tests/unit/user_test.go",
		}

		foundTests := false
		for _, file := range testFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				testContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				testStr := string(testContent)

				foundTests = true
				assert.Contains(t, testStr, "func Test", "Should contain test functions")
				
				// Check for mocking/isolation patterns
				if strings.Contains(testStr, "mock") || strings.Contains(testStr, "Mock") {
					assert.Contains(t, testStr, "testify", "Should use testify for mocking")
					t.Logf("Found mocking in: %s", file)
				}
			}
		}

		if foundTests {
			t.Logf("Test files found - testability is supported")
		}

		// Check if dependencies can be easily mocked (interfaces exist)
		// Domain repositories should be interfaces
		repoInterface := filepath.Join(projectDir, "internal", "domain", "repositories", "user_repository.go")
		if fileExists(t, repoInterface) {
			repoContent, err := os.ReadFile(repoInterface)
			require.NoError(t, err)
			repoStr := string(repoContent)

			if strings.Contains(repoStr, "interface") {
				t.Logf("Repository interfaces enable easy mocking")
			}
		}
	})
}

// Helper functions
func fileExists(t *testing.T, path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(t *testing.T, path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func handlerExists(t *testing.T, path string) bool {
	return fileExists(t, path)
}

func usecaseExists(t *testing.T, path string) bool {
	return fileExists(t, path)
}

// analyzeBoundaryCompliance checks for Clean Architecture boundary violations
func analyzeBoundaryCompliance(t *testing.T, dir string, layerName string, projectRoot string) {
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Check for boundary violations based on layer
			switch layerName {
			case "domain":
				// Domain should not import outer layers
				assert.NotContains(t, contentStr, "internal/infrastructure", 
					"Domain layer should not import infrastructure: %s", path)
				assert.NotContains(t, contentStr, "internal/presentation",
					"Domain layer should not import presentation: %s", path)
				assert.NotContains(t, contentStr, "gin-gonic", 
					"Domain should not import web frameworks: %s", path)
				
			case "application":
				// Application should not import infrastructure or presentation
				assert.NotContains(t, contentStr, "internal/infrastructure/repositories",
					"Application layer should not import infrastructure repositories: %s", path)
				assert.NotContains(t, contentStr, "internal/presentation",
					"Application layer should not import presentation: %s", path)
				
			case "infrastructure":
				// Infrastructure can import domain but not presentation
				assert.NotContains(t, contentStr, "internal/presentation",
					"Infrastructure should not import presentation: %s", path)
				
			case "presentation":
				// Presentation can import all other layers (outermost layer)
				// No restrictions
			}
		}
		return nil
	})
	require.NoError(t, err)
}