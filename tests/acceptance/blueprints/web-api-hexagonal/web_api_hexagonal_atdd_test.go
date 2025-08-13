package web_api_hexagonal

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

// TestWebApiHexagonalBlueprintATDD validates the acceptance criteria for web-api-hexagonal blueprint
// This ensures the web-api-hexagonal blueprint generates proper Ports & Adapters (Hexagonal Architecture) patterns
func TestWebApiHexagonalBlueprintATDD(t *testing.T) {
	t.Run("web_api_hexagonal_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: web-api-hexagonal should be in the list with correct Hexagonal Architecture description

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/web-api-hexagonal)
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
		assert.Contains(t, outputStr, "web-api-hexagonal", "web-api-hexagonal blueprint should be listed")
		assert.Contains(t, outputStr, "Hexagonal Architecture", "Should show Hexagonal Architecture description")
		assert.Contains(t, outputStr, "ports and adapters", "Should emphasize ports and adapters pattern")
	})

	t.Run("web_api_hexagonal_generates_hexagonal_structure", func(t *testing.T) {
		// GIVEN: User wants a Hexagonal Architecture web API
		// WHEN: User generates a project with web-api-hexagonal blueprint
		// THEN: Should generate proper hexagonal structure (core domain + ports + adapters + infrastructure)

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

		// Generate a web-api-hexagonal project
		generateCmd := exec.Command("./go-starter", "new", "test-hexagonal-api",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/hexagonal-api",
			"--framework=gin",
			"--logger=slog",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-hexagonal-api")

		// Verify Hexagonal Architecture layer directories exist
		hexagonalLayers := []string{
			// Core Domain (hexagon center - business logic)
			"internal/core/domain/entities",
			"internal/core/domain/valueobjects",
			"internal/core/domain/services",
			"internal/core/usecases",
			"internal/core/services",

			// Ports (interfaces - hexagon boundaries)
			"internal/core/ports/input",   // Primary ports (driving adapters)
			"internal/core/ports/output",  // Secondary ports (driven adapters)
			"internal/ports/primary",      // Alternative naming
			"internal/ports/secondary",    // Alternative naming

			// Primary Adapters (drivers - external world calling in)
			"internal/adapters/primary/http",
			"internal/adapters/primary/rest",
			"internal/adapters/primary/controllers",

			// Secondary Adapters (driven - called by business logic)
			"internal/adapters/secondary/database",
			"internal/adapters/secondary/repositories",
			"internal/adapters/secondary/external",
			"internal/adapters/secondary/persistence",

			// Infrastructure (cross-cutting concerns)
			"internal/infrastructure/config",
			"internal/infrastructure/logging",
			"internal/infrastructure/middleware",
		}

		existingLayers := []string{}
		for _, layer := range hexagonalLayers {
			if dirExists(t, filepath.Join(projectDir, layer)) {
				existingLayers = append(existingLayers, layer)
				t.Logf("Found hexagonal layer: %s", layer)
			}
		}

		// Should have at least core structure
		assert.True(t, len(existingLayers) > 0, "Should have hexagonal architecture layers")
		
		// Core should always exist
		coreExists := dirExists(t, filepath.Join(projectDir, "internal", "core")) ||
					 dirExists(t, filepath.Join(projectDir, "internal", "domain"))
		assert.True(t, coreExists, "Core domain should exist")

		// Ports should exist (either style)
		portsExist := dirExists(t, filepath.Join(projectDir, "internal", "ports")) ||
					 dirExists(t, filepath.Join(projectDir, "internal", "core", "ports"))
		assert.True(t, portsExist, "Ports should exist")

		// Adapters should exist
		adaptersExist := dirExists(t, filepath.Join(projectDir, "internal", "adapters")) ||
					    dirExists(t, filepath.Join(projectDir, "internal", "infrastructure"))
		assert.True(t, adaptersExist, "Adapters should exist")

		// Verify essential Hexagonal files exist
		essentialFiles := []string{
			// Foundation
			"main.go",
			"go.mod",
			"README.md",

			// Core domain entities
			"internal/core/domain/entities/user.go",
			"internal/core/domain/entities/product.go",

			// Core use cases (application layer)
			"internal/core/usecases/user_usecase.go",
			"internal/core/usecases/create_user_usecase.go",

			// Primary ports (interfaces for use cases)
			"internal/core/ports/input/user_service_port.go",
			"internal/ports/primary/user_port.go",

			// Secondary ports (interfaces for repositories)
			"internal/core/ports/output/user_repository_port.go",
			"internal/ports/secondary/user_repository_port.go",

			// Primary adapters (HTTP controllers)
			"internal/adapters/primary/http/user_handler.go",
			"internal/adapters/primary/rest/user_controller.go",

			// Secondary adapters (database implementation)
			"internal/adapters/secondary/database/user_repository.go",
			"internal/adapters/secondary/repositories/user_repository_impl.go",

			// Configuration
			"internal/infrastructure/config/config.go",
			"internal/infrastructure/logging/logger.go",
		}

		existingFiles := []string{}
		for _, file := range essentialFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				existingFiles = append(existingFiles, file)
				t.Logf("Found essential file: %s", file)
			}
		}

		assert.True(t, len(existingFiles) > 5, "Should have essential hexagonal files (found %d)", len(existingFiles))
	})

	t.Run("web_api_hexagonal_implements_ports_correctly", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining ports implementation
		// THEN: Should define proper input ports (primary) and output ports (secondary) as interfaces

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
		generateCmd := exec.Command("./go-starter", "new", "test-ports",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/ports",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-ports")

		// Check Primary Ports (Input Ports - for use cases)
		primaryPortFiles := []string{
			"internal/core/ports/input/user_service_port.go",
			"internal/ports/primary/user_port.go", 
			"internal/ports/input/user_usecase_port.go",
		}

		foundPrimaryPorts := false
		for _, file := range primaryPortFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				portContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				portStr := string(portContent)

				foundPrimaryPorts = true

				// Primary ports should be interfaces
				assert.Contains(t, portStr, "interface", "Primary port should be interface")
				
				// Should define business operations
				businessOps := []string{"Create", "Get", "Update", "Delete", "List"}
				foundOp := false
				for _, op := range businessOps {
					if strings.Contains(portStr, op) {
						foundOp = true
						t.Logf("Found primary port operation: %s in %s", op, file)
					}
				}
				assert.True(t, foundOp, "Primary port should define business operations")

				// Should NOT depend on external frameworks
				assert.NotContains(t, portStr, "gin-gonic", "Primary port should not depend on web framework")
				assert.NotContains(t, portStr, "gorm.io", "Primary port should not depend on ORM")
				assert.NotContains(t, portStr, "database/sql", "Primary port should not depend on database")

				t.Logf("Found primary port: %s", file)
			}
		}

		// Check Secondary Ports (Output Ports - for repositories/external services)
		secondaryPortFiles := []string{
			"internal/core/ports/output/user_repository_port.go",
			"internal/ports/secondary/user_repository_port.go",
			"internal/ports/output/user_persistence_port.go",
		}

		foundSecondaryPorts := false
		for _, file := range secondaryPortFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				portContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				portStr := string(portContent)

				foundSecondaryPorts = true

				// Secondary ports should be interfaces
				assert.Contains(t, portStr, "interface", "Secondary port should be interface")
				
				// Should define persistence/external operations
				persistenceOps := []string{"Save", "FindByID", "FindAll", "Delete", "Update"}
				foundOp := false
				for _, op := range persistenceOps {
					if strings.Contains(portStr, op) {
						foundOp = true
						t.Logf("Found secondary port operation: %s in %s", op, file)
					}
				}
				assert.True(t, foundOp, "Secondary port should define persistence operations")

				// Should work with domain entities
				assert.Contains(t, portStr, "internal/core/domain/entities", "Secondary port should reference domain entities")

				t.Logf("Found secondary port: %s", file)
			}
		}

		assert.True(t, foundPrimaryPorts || foundSecondaryPorts, "Should implement either primary or secondary ports")
	})

	t.Run("web_api_hexagonal_implements_adapters_correctly", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining adapters implementation
		// THEN: Should implement proper primary adapters (HTTP) and secondary adapters (Database)

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
		generateCmd := exec.Command("./go-starter", "new", "test-adapters",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/adapters",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-adapters")

		// Check Primary Adapters (HTTP controllers/handlers)
		primaryAdapterFiles := []string{
			"internal/adapters/primary/http/user_handler.go",
			"internal/adapters/primary/rest/user_controller.go",
			"internal/adapters/primary/controllers/user_controller.go",
		}

		foundPrimaryAdapters := false
		for _, file := range primaryAdapterFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				adapterContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				adapterStr := string(adapterContent)

				foundPrimaryAdapters = true

				// Primary adapters should implement HTTP handling
				assert.Contains(t, adapterStr, "gin-gonic/gin", "Primary adapter should use web framework")
				
				// Should depend on primary ports (use cases)
				portImports := []string{
					"internal/core/ports/input",
					"internal/ports/primary", 
					"internal/core/usecases",
				}
				foundPortImport := false
				for _, imp := range portImports {
					if strings.Contains(adapterStr, imp) {
						foundPortImport = true
						t.Logf("Primary adapter imports port: %s", imp)
					}
				}
				assert.True(t, foundPortImport, "Primary adapter should depend on ports/use cases")

				// Should handle HTTP concerns (not business logic)
				httpConcerns := []string{"http.StatusOK", "c.JSON", "c.Param", "c.Bind"}
				foundHttpConcern := false
				for _, concern := range httpConcerns {
					if strings.Contains(adapterStr, concern) {
						foundHttpConcern = true
					}
				}
				assert.True(t, foundHttpConcern, "Primary adapter should handle HTTP concerns")

				// Should NOT contain business logic
				assert.NotContains(t, adapterStr, "business validation", "Primary adapter should not contain business logic")

				t.Logf("Found primary adapter: %s", file)
			}
		}

		// Check Secondary Adapters (Database repositories)
		secondaryAdapterFiles := []string{
			"internal/adapters/secondary/database/user_repository.go",
			"internal/adapters/secondary/repositories/user_repository_impl.go",
			"internal/adapters/secondary/persistence/user_persistence.go",
		}

		foundSecondaryAdapters := false
		for _, file := range secondaryAdapterFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				adapterContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				adapterStr := string(adapterContent)

				foundSecondaryAdapters = true

				// Secondary adapters should implement database/external concerns
				dbConcerns := []string{"gorm.io/gorm", "database/sql", "sql.DB"}
				foundDbConcern := false
				for _, concern := range dbConcerns {
					if strings.Contains(adapterStr, concern) {
						foundDbConcern = true
						t.Logf("Secondary adapter uses database: %s", concern)
					}
				}
				assert.True(t, foundDbConcern, "Secondary adapter should handle database concerns")

				// Should implement secondary ports
				portImports := []string{
					"internal/core/ports/output",
					"internal/ports/secondary",
				}
				for _, imp := range portImports {
					if strings.Contains(adapterStr, imp) {
						t.Logf("Secondary adapter implements port: %s", imp)
					}
				}
				
				// Should work with domain entities
				assert.Contains(t, adapterStr, "internal/core/domain/entities", 
					"Secondary adapter should work with domain entities")

				// Should implement repository pattern methods
				repoMethods := []string{"Save", "FindByID", "FindAll", "Delete", "Update"}
				foundRepoMethod := false
				for _, method := range repoMethods {
					if strings.Contains(adapterStr, method) {
						foundRepoMethod = true
					}
				}
				assert.True(t, foundRepoMethod, "Secondary adapter should implement repository methods")

				t.Logf("Found secondary adapter: %s", file)
			}
		}

		assert.True(t, foundPrimaryAdapters || foundSecondaryAdapters, "Should implement primary or secondary adapters")
	})

	t.Run("web_api_hexagonal_enforces_dependency_inversion", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining dependency relationships
		// THEN: Should enforce dependency inversion (core doesn't depend on adapters/infrastructure)

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
		generateCmd := exec.Command("./go-starter", "new", "test-dependencies",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/dependencies",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-dependencies")

		// 1. Core domain should NOT depend on adapters or infrastructure
		coreDirs := []string{
			"internal/core/domain",
			"internal/core/usecases", 
			"internal/core/services",
		}

		for _, coreDir := range coreDirs {
			fullDir := filepath.Join(projectDir, coreDir)
			if dirExists(t, fullDir) {
				analyzeCoreDependencies(t, fullDir, projectDir)
			}
		}

		// 2. Ports should NOT depend on adapters or infrastructure
		portDirs := []string{
			"internal/core/ports",
			"internal/ports",
		}

		for _, portDir := range portDirs {
			fullDir := filepath.Join(projectDir, portDir)
			if dirExists(t, fullDir) {
				analyzePortDependencies(t, fullDir, projectDir)
			}
		}

		// 3. Primary adapters should depend on core (use cases/ports)
		primaryAdapterDirs := []string{
			"internal/adapters/primary",
		}

		for _, adapterDir := range primaryAdapterDirs {
			fullDir := filepath.Join(projectDir, adapterDir)
			if dirExists(t, fullDir) {
				analyzePrimaryAdapterDependencies(t, fullDir, projectDir)
			}
		}

		// 4. Secondary adapters should implement secondary ports
		secondaryAdapterDirs := []string{
			"internal/adapters/secondary",
		}

		for _, adapterDir := range secondaryAdapterDirs {
			fullDir := filepath.Join(projectDir, adapterDir)
			if dirExists(t, fullDir) {
				analyzeSecondaryAdapterDependencies(t, fullDir, projectDir)
			}
		}
	})

	t.Run("web_api_hexagonal_project_builds_successfully", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
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
			"--architecture=hexagonal",
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
		buildGeneratedCmd := exec.Command("go", "build", "-o", "hexagonal-api", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated web-api-hexagonal project should build successfully: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "hexagonal-api"), "Hexagonal API binary should be created")
	})

	t.Run("web_api_hexagonal_supports_multiple_frameworks", func(t *testing.T) {
		// GIVEN: web-api-hexagonal blueprint with different web frameworks
		// WHEN: Generating projects with different frameworks
		// THEN: Each should generate framework-specific primary adapters while maintaining hexagonal structure

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
				projectName := fmt.Sprintf("test-hex-%s", framework)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=web-api",
					"--architecture=hexagonal",
					"--module=github.com/test/"+projectName,
					"--framework="+framework,
					"--database-driver=sqlite",
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s: %s", framework, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Verify hexagonal structure exists regardless of framework
				assert.True(t, dirExists(t, filepath.Join(projectDir, "internal", "core")) ||
						   dirExists(t, filepath.Join(projectDir, "internal", "domain")), 
					"Core domain should exist")
				assert.True(t, dirExists(t, filepath.Join(projectDir, "internal", "ports")) ||
						   dirExists(t, filepath.Join(projectDir, "internal", "adapters")),
					"Ports/Adapters should exist")

				// Check that primary adapters use correct framework
				primaryAdapterPaths := []string{
					"internal/adapters/primary/http/user_handler.go",
					"internal/adapters/primary/rest/user_controller.go",
					"internal/presentation/handlers/user_handler.go",
					"internal/controllers/user_controller.go",
				}

				frameworkImportFound := false
				for _, path := range primaryAdapterPaths {
					if fileExists(t, filepath.Join(projectDir, path)) {
						adapterContent, err := os.ReadFile(filepath.Join(projectDir, path))
						require.NoError(t, err)
						adapterStr := string(adapterContent)

						switch framework {
						case "gin":
							if strings.Contains(adapterStr, "gin-gonic/gin") {
								frameworkImportFound = true
							}
						case "echo":
							if strings.Contains(adapterStr, "labstack/echo") {
								frameworkImportFound = true
							}
						case "fiber":
							if strings.Contains(adapterStr, "gofiber/fiber") {
								frameworkImportFound = true
							}
						case "chi":
							if strings.Contains(adapterStr, "go-chi/chi") {
								frameworkImportFound = true
							}
						}

						if frameworkImportFound {
							t.Logf("Found %s framework import in %s", framework, path)
							break
						}
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
				require.NoError(t, err, "Hexagonal API with %s should compile: %s", framework, string(output))
			})
		}
	})

	t.Run("web_api_hexagonal_supports_multiple_loggers", func(t *testing.T) {
		// GIVEN: web-api-hexagonal blueprint with different logger configurations
		// WHEN: Generating projects with different loggers
		// THEN: Logger should be properly injected as infrastructure concern

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
					"--architecture=hexagonal",
					"--module=github.com/test/"+projectName,
					"--framework=gin",
					"--logger="+logger,
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s logger: %s", logger, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Check that logger setup exists in infrastructure
				loggerPaths := []string{
					"internal/infrastructure/logging/logger.go",
					"internal/infrastructure/logger/logger.go",
					"internal/logger/logger.go",
				}

				loggerFound := false
				for _, path := range loggerPaths {
					if fileExists(t, filepath.Join(projectDir, path)) {
						loggerFound = true
						t.Logf("Found logger at: %s", path)
						break
					}
				}
				assert.True(t, loggerFound, "Logger setup should exist for %s", logger)

				// Check logger is used in adapters (not core)
				primaryAdapterPaths := []string{
					"internal/adapters/primary/http/user_handler.go",
					"internal/adapters/primary/rest/user_controller.go",
				}

				for _, path := range primaryAdapterPaths {
					if fileExists(t, filepath.Join(projectDir, path)) {
						adapterContent, err := os.ReadFile(filepath.Join(projectDir, path))
						require.NoError(t, err)
						adapterStr := string(adapterContent)

						// Adapters can use logger (infrastructure concern)
						if strings.Contains(adapterStr, "internal/infrastructure/log") ||
						   strings.Contains(adapterStr, "internal/logger") {
							t.Logf("Primary adapter uses logger")
						}
					}
				}

				// Verify project compiles with the logger
				modTidyCmd := exec.Command("go", "mod", "tidy")
				modTidyCmd.Dir = projectDir
				_, err = modTidyCmd.CombinedOutput()
				require.NoError(t, err)

				buildGeneratedCmd := exec.Command("go", "build", "-o", "test-"+logger, ".")
				buildGeneratedCmd.Dir = projectDir
				output, err = buildGeneratedCmd.CombinedOutput()
				require.NoError(t, err, "web-api-hexagonal with %s logger should compile: %s", logger, string(output))
			})
		}
	})

	t.Run("web_api_hexagonal_enables_testability", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining testability features
		// THEN: Should enable easy testing through ports and dependency injection

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
			"--architecture=hexagonal",
			"--module=github.com/test/testability",
			"--framework=gin",
			"--database-driver=sqlite",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-testability")

		// Look for test files
		testFiles := []string{
			"internal/core/usecases/user_usecase_test.go",
			"internal/core/domain/entities/user_test.go",
			"internal/adapters/primary/http/user_handler_test.go",
			"internal/adapters/secondary/database/user_repository_test.go",
			"tests/unit/core/user_test.go",
		}

		foundTests := false
		for _, file := range testFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				testContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				testStr := string(testContent)

				foundTests = true
				assert.Contains(t, testStr, "func Test", "Should contain test functions")
				
				// Check for mocking patterns (enabled by ports)
				if strings.Contains(testStr, "mock") || strings.Contains(testStr, "Mock") {
					assert.Contains(t, testStr, "testify", "Should use testify for mocking")
					t.Logf("Found mocking in: %s", file)
				}

				t.Logf("Found test file: %s", file)
			}
		}

		// Check ports enable easy mocking
		portFiles := []string{
			"internal/core/ports/output/user_repository_port.go",
			"internal/ports/secondary/user_repository_port.go",
		}

		for _, file := range portFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				portContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				portStr := string(portContent)

				if strings.Contains(portStr, "interface") {
					t.Logf("Port interface enables mocking: %s", file)
				}
			}
		}

		// Check use cases can be tested in isolation (depend only on interfaces)
		usecaseFiles := []string{
			"internal/core/usecases/user_usecase.go",
			"internal/core/usecases/create_user_usecase.go",
		}

		for _, file := range usecaseFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				usecaseContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				usecaseStr := string(usecaseContent)

				// Use cases should depend on ports (interfaces) not implementations
				portImports := []string{
					"internal/core/ports/output",
					"internal/ports/secondary",
				}
				
				for _, imp := range portImports {
					if strings.Contains(usecaseStr, imp) {
						t.Logf("Use case depends on port interface (testable): %s", imp)
					}
				}

				// Should NOT depend on concrete implementations
				assert.NotContains(t, usecaseStr, "internal/adapters/secondary", 
					"Use case should not depend on concrete adapter")
				assert.NotContains(t, usecaseStr, "gorm.io/gorm", 
					"Use case should not depend on ORM")
			}
		}

		if foundTests {
			t.Logf("Hexagonal architecture enables testability through ports")
		}
	})

	t.Run("web_api_hexagonal_supports_multiple_adapters", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining adapter implementation
		// THEN: Should support multiple adapters for same port (adapter interchangeability)

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

		generateCmd := exec.Command("./go-starter", "new", "test-multiple-adapters",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/multiple-adapters",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-multiple-adapters")

		// Check for multiple secondary adapters for same port
		repositoryAdapters := []string{
			"internal/adapters/secondary/database/user_repository.go",        // Database adapter
			"internal/adapters/secondary/memory/user_repository.go",          // In-memory adapter
			"internal/adapters/secondary/redis/user_repository.go",           // Redis adapter
			"internal/adapters/secondary/external/user_api_repository.go",    // External API adapter
		}

		foundAdapters := 0
		for _, adapter := range repositoryAdapters {
			if fileExists(t, filepath.Join(projectDir, adapter)) {
				adapterContent, err := os.ReadFile(filepath.Join(projectDir, adapter))
				require.NoError(t, err)
				adapterStr := string(adapterContent)

				// All adapters should implement same port interface
				portImports := []string{
					"internal/core/ports/output",
					"internal/ports/secondary",
				}

				implementsPort := false
				for _, imp := range portImports {
					if strings.Contains(adapterStr, imp) {
						implementsPort = true
						break
					}
				}

				if implementsPort {
					foundAdapters++
					t.Logf("Found adapter implementing port: %s", adapter)
				}
			}
		}

		// Check for multiple primary adapters
		primaryAdapters := []string{
			"internal/adapters/primary/http/user_handler.go",    // HTTP REST
			"internal/adapters/primary/grpc/user_service.go",   // gRPC
			"internal/adapters/primary/cli/user_commands.go",   // CLI
			"internal/adapters/primary/graphql/user_resolver.go", // GraphQL
		}

		for _, adapter := range primaryAdapters {
			if fileExists(t, filepath.Join(projectDir, adapter)) {
				foundAdapters++
				t.Logf("Found primary adapter: %s", adapter)
			}
		}

		// Check main.go for dependency injection that supports adapter swapping
		if fileExists(t, filepath.Join(projectDir, "main.go")) {
			mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
			require.NoError(t, err)
			mainStr := string(mainContent)

			// Look for dependency injection patterns
			diPatterns := []string{"New", "wire", "container", "inject"}
			for _, pattern := range diPatterns {
				if strings.Contains(mainStr, pattern) {
					t.Logf("Found DI pattern in main.go: %s", pattern)
				}
			}

			// Should wire adapters to ports
			if strings.Contains(mainStr, "adapters") && strings.Contains(mainStr, "ports") {
				t.Logf("Main.go wires adapters to ports")
			}
		}

		if foundAdapters > 1 {
			t.Logf("Multiple adapters found (%d) - supports adapter interchangeability", foundAdapters)
		}
	})
}

// TestWebApiHexagonalCompliance validates Hexagonal Architecture compliance
func TestWebApiHexagonalCompliance(t *testing.T) {
	t.Run("isolates_business_logic_in_core", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Examining core domain and use cases
		// THEN: Should contain pure business logic without external dependencies

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

		generateCmd := exec.Command("./go-starter", "new", "test-isolation",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/isolation",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-isolation")

		// Analyze core domain for purity
		coreDirs := []string{
			"internal/core/domain",
			"internal/core/usecases",
			"internal/core/services",
		}

		for _, coreDir := range coreDirs {
			fullDir := filepath.Join(projectDir, coreDir)
			if dirExists(t, fullDir) {
				analyzeBusinessLogicPurity(t, fullDir, projectDir)
			}
		}
	})

	t.Run("enforces_hexagonal_dependency_rules", func(t *testing.T) {
		// GIVEN: A generated web-api-hexagonal project
		// WHEN: Analyzing import dependencies
		// THEN: Should follow hexagonal dependency rules strictly

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

		generateCmd := exec.Command("./go-starter", "new", "test-hex-rules",
			"--type=web-api",
			"--architecture=hexagonal",
			"--module=github.com/test/hex-rules",
			"--framework=gin",
			"--database-driver=postgres",
			"--database-orm=gorm",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-hex-rules")

		// Validate dependency rules
		validateHexagonalDependencyRules(t, projectDir)
	})
}

// Helper functions for hexagonal architecture analysis

func fileExists(t *testing.T, path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func dirExists(t *testing.T, path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func analyzeCoreDependencies(t *testing.T, coreDir string, projectRoot string) {
	err := filepath.Walk(coreDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Core should NEVER depend on adapters or infrastructure
			forbiddenImports := []string{
				"internal/adapters",
				"internal/infrastructure",
				"gin-gonic/gin",
				"labstack/echo",
				"gofiber/fiber",
				"go-chi/chi",
				"gorm.io/gorm",
				"database/sql",
			}

			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, contentStr, forbidden,
					"Core should not import %s in %s", forbidden, path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func analyzePortDependencies(t *testing.T, portDir string, projectRoot string) {
	err := filepath.Walk(portDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Ports should only define interfaces
			assert.Contains(t, contentStr, "interface", "Ports should define interfaces in %s", path)
			
			// Ports should NOT depend on implementations
			forbiddenImports := []string{
				"internal/adapters",
				"gin-gonic/gin",
				"gorm.io/gorm",
				"database/sql",
			}

			for _, forbidden := range forbiddenImports {
				assert.NotContains(t, contentStr, forbidden,
					"Port should not import %s in %s", forbidden, path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func analyzePrimaryAdapterDependencies(t *testing.T, adapterDir string, projectRoot string) {
	err := filepath.Walk(adapterDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Primary adapters should depend on core/ports
			coreImports := []string{
				"internal/core",
				"internal/ports",
			}

			foundCoreImport := false
			for _, imp := range coreImports {
				if strings.Contains(contentStr, imp) {
					foundCoreImport = true
					break
				}
			}

			if !foundCoreImport {
				t.Logf("WARNING: Primary adapter may not depend on core: %s", path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func analyzeSecondaryAdapterDependencies(t *testing.T, adapterDir string, projectRoot string) {
	err := filepath.Walk(adapterDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Secondary adapters should implement secondary ports
			portImports := []string{
				"internal/core/ports/output",
				"internal/ports/secondary",
			}

			for _, imp := range portImports {
				if strings.Contains(contentStr, imp) {
					t.Logf("Secondary adapter implements port: %s", path)
					break
				}
			}

			// Should work with domain entities
			if strings.Contains(contentStr, "internal/core/domain/entities") {
				t.Logf("Secondary adapter works with domain entities: %s", path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func analyzeBusinessLogicPurity(t *testing.T, coreDir string, projectRoot string) {
	err := filepath.Walk(coreDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			// Core should contain business logic
			businessPatterns := []string{
				"business", "validate", "calculate", "process", "rule", "policy",
			}

			for _, pattern := range businessPatterns {
				if strings.Contains(strings.ToLower(contentStr), pattern) {
					t.Logf("Found business logic in core: %s", path)
					break
				}
			}

			// Should NOT contain infrastructure concerns
			infraConcerns := []string{
				"http.StatusOK", "c.JSON", "db.Create", "SELECT", "INSERT",
				"gin.Context", "echo.Context",
			}

			for _, concern := range infraConcerns {
				assert.NotContains(t, contentStr, concern,
					"Core should not contain infrastructure concern %s in %s", concern, path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}

func validateHexagonalDependencyRules(t *testing.T, projectDir string) {
	rules := map[string][]string{
		"internal/core":                      {"internal/adapters", "internal/infrastructure"}, // Core cannot depend on outer layers
		"internal/ports":                     {"internal/adapters"},                             // Ports cannot depend on adapters  
		"internal/adapters/secondary":        {},                                                // Secondary adapters have no restrictions
		"internal/adapters/primary":          {},                                                // Primary adapters have no restrictions
	}

	for layer, forbiddenDeps := range rules {
		layerPath := filepath.Join(projectDir, layer)
		if dirExists(t, layerPath) {
			analyzeDependencyViolations(t, layerPath, forbiddenDeps, projectDir)
		}
	}
}

func analyzeDependencyViolations(t *testing.T, layerDir string, forbiddenDeps []string, projectRoot string) {
	err := filepath.Walk(layerDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			contentStr := string(content)
			
			for _, forbidden := range forbiddenDeps {
				assert.NotContains(t, contentStr, forbidden,
					"Layer %s should not depend on %s in %s", layerDir, forbidden, path)
			}
		}
		return nil
	})
	require.NoError(t, err)
}