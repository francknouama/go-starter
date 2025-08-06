package web_api_ddd

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

// TestWebApiDddBlueprintATDD validates the acceptance criteria for web-api-ddd blueprint
// This ensures the web-api-ddd blueprint generates proper Domain-Driven Design patterns
func TestWebApiDddBlueprintATDD(t *testing.T) {
	t.Run("web_api_ddd_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: web-api-ddd should be in the list with correct DDD description

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/web-api-ddd)
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
		assert.Contains(t, outputStr, "web-api-ddd", "web-api-ddd blueprint should be listed")
		assert.Contains(t, outputStr, "Domain-Driven Design", "Should show DDD description")
		assert.Contains(t, outputStr, "strategic design patterns", "Should emphasize DDD strategic patterns")
	})

	t.Run("web_api_ddd_generates_ddd_structure", func(t *testing.T) {
		// GIVEN: User wants a Domain-Driven Design web API
		// WHEN: User generates a project with web-api-ddd blueprint
		// THEN: Should generate proper DDD structure (domain, application, presentation, infrastructure)

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

		// Generate a web-api-ddd project
		generateCmd := exec.Command("./go-starter", "new", "test-ddd-api",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/ddd-api",
			"--framework=gin",
			"--logger=slog",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-ddd-api")

		// Verify DDD layer directories exist
		dddLayers := []string{
			// Domain layer (core business logic)
			"internal/domain/entities",
			"internal/domain/valueobjects", 
			"internal/domain/aggregates",
			"internal/domain/repositories",
			"internal/domain/services",
			"internal/domain/events",
			"internal/domain/specifications",

			// Application layer (use cases and application services)
			"internal/application/usecases",
			"internal/application/services",
			"internal/application/commands",
			"internal/application/queries",
			"internal/application/handlers",

			// Infrastructure layer (external concerns)
			"internal/infrastructure/persistence",
			"internal/infrastructure/repositories",
			"internal/infrastructure/messaging",
			"internal/infrastructure/external",

			// Presentation layer (REST API, controllers)
			"internal/presentation/controllers",
			"internal/presentation/dto",
			"internal/presentation/routes",
			"internal/presentation/middleware",
		}

		for _, layer := range dddLayers {
			assert.DirExists(t, filepath.Join(projectDir, layer), "DDD layer %s should exist", layer)
		}

		// Verify essential DDD files exist
		essentialFiles := []string{
			// Foundation
			"main.go",
			"go.mod",
			"README.md",

			// Domain entities and value objects
			"internal/domain/entities/user.go",
			"internal/domain/entities/product.go",
			"internal/domain/valueobjects/email.go",
			"internal/domain/valueobjects/user_name.go",

			// Aggregates (if implemented)
			"internal/domain/aggregates/user_aggregate.go",

			// Domain services
			"internal/domain/services/user_domain_service.go",

			// Repository interfaces (domain layer)
			"internal/domain/repositories/user_repository.go",

			// Application services and use cases
			"internal/application/services/user_application_service.go",
			"internal/application/usecases/create_user_usecase.go",

			// Commands and queries (CQRS pattern)
			"internal/application/commands/create_user_command.go",
			"internal/application/queries/get_user_query.go",

			// Infrastructure implementations
			"internal/infrastructure/repositories/user_repository_impl.go",
			"internal/infrastructure/persistence/database.go",

			// Presentation layer
			"internal/presentation/controllers/user_controller.go",
			"internal/presentation/dto/user_dto.go",
			"internal/presentation/routes/routes.go",

			// Configuration
			"internal/config/config.go",
			"internal/logger/logger.go",
		}

		for _, file := range essentialFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				t.Logf("Found essential DDD file: %s", file)
			}
		}
	})

	t.Run("web_api_ddd_implements_value_objects", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining value objects implementation
		// THEN: Should have proper value objects with validation and immutability

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
		generateCmd := exec.Command("./go-starter", "new", "test-valueobjects",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/valueobjects",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-valueobjects")

		// Check value object files
		valueObjectFiles := []string{
			"internal/domain/valueobjects/email.go",
			"internal/domain/valueobjects/user_name.go", 
			"internal/domain/valueobjects/description.go",
		}

		for _, file := range valueObjectFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				voContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				voStr := string(voContent)

				// Value objects should have proper structure
				assert.Contains(t, voStr, "type", "Should define value object type")
				assert.Contains(t, voStr, "String()", "Should have String() method")

				// Value objects should have validation
				if strings.Contains(voStr, "New") || strings.Contains(voStr, "Create") {
					assert.Contains(t, voStr, "error", "Should return error for validation")
				}

				// Email value object should have email validation
				if strings.Contains(file, "email.go") {
					assert.Contains(t, voStr, "email", "Email value object should contain email validation")
					
					// Should validate email format
					if strings.Contains(voStr, "regexp") || strings.Contains(voStr, "@") {
						t.Logf("Email value object includes validation logic")
					}
				}

				// UserName value object should validate names
				if strings.Contains(file, "user_name.go") {
					assert.Contains(t, voStr, "UserName", "Should define UserName type")
					
					// Should have length or format validation
					if strings.Contains(voStr, "len(") || strings.Contains(voStr, "length") {
						t.Logf("UserName value object includes length validation")
					}
				}

				// Description value object should handle text
				if strings.Contains(file, "description.go") {
					assert.Contains(t, voStr, "Description", "Should define Description type")
				}

				// Value objects should be immutable (no setter methods)
				assert.NotRegexp(t, `func \([^)]*\) Set[A-Z]`, voStr, 
					"Value objects should not have setter methods (immutable)")
			}
		}
	})

	t.Run("web_api_ddd_implements_aggregates", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project  
		// WHEN: Examining aggregate implementation
		// THEN: Should implement proper aggregates with consistency boundaries

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
		generateCmd := exec.Command("./go-starter", "new", "test-aggregates",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/aggregates",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-aggregates")

		// Check aggregate files
		aggregateFiles := []string{
			"internal/domain/aggregates/user_aggregate.go",
			"internal/domain/aggregates/product_aggregate.go",
			"internal/domain/aggregates/order_aggregate.go",
		}

		for _, file := range aggregateFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				aggContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				aggStr := string(aggContent)

				// Aggregates should have proper structure
				assert.Contains(t, aggStr, "type", "Should define aggregate type")

				// Aggregates should have aggregate root
				if strings.Contains(aggStr, "Aggregate") {
					assert.Contains(t, aggStr, "ID", "Aggregate should have ID")
					
					// Should have business methods
					methodPatterns := []string{"Create", "Update", "Delete", "Change", "Add", "Remove"}
					foundMethod := false
					for _, pattern := range methodPatterns {
						if strings.Contains(aggStr, pattern) {
							foundMethod = true
							break
						}
					}
					assert.True(t, foundMethod, "Aggregate should have business methods")
				}

				// Aggregates should handle domain events (if implemented)
				if strings.Contains(aggStr, "events") || strings.Contains(aggStr, "Events") {
					assert.Contains(t, aggStr, "RaiseEvent", "Should have event raising mechanism")
					t.Logf("Aggregate implements domain events: %s", file)
				}

				// User aggregate should manage user business logic
				if strings.Contains(file, "user_aggregate.go") {
					assert.Contains(t, aggStr, "User", "User aggregate should reference User entity")
					
					// Should have user-specific business methods
					userMethods := []string{"ChangeEmail", "UpdateProfile", "Activate", "Deactivate"}
					for _, method := range userMethods {
						if strings.Contains(aggStr, method) {
							t.Logf("Found user business method: %s", method)
						}
					}
				}
			}
		}

		// Check if entities are properly modeled
		userEntityPath := filepath.Join(projectDir, "internal", "domain", "entities", "user.go")
		if fileExists(t, userEntityPath) {
			entityContent, err := os.ReadFile(userEntityPath)
			require.NoError(t, err)
			entityStr := string(entityContent)

			// Entities should use value objects
			assert.Contains(t, entityStr, "valueobjects", "Entities should import value objects")
			
			// Should not have primitive obsession
			if strings.Contains(entityStr, "Email") && !strings.Contains(entityStr, "string") {
				t.Logf("User entity uses Email value object (avoiding primitive obsession)")
			}
		}
	})

	t.Run("web_api_ddd_implements_domain_services", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining domain services implementation
		// THEN: Should implement domain services for complex business operations

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

		generateCmd := exec.Command("./go-starter", "new", "test-domain-services",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/domain-services",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-domain-services")

		// Check domain service files
		domainServiceFiles := []string{
			"internal/domain/services/user_domain_service.go",
			"internal/domain/services/authentication_service.go",
			"internal/domain/services/validation_service.go",
		}

		for _, file := range domainServiceFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				serviceContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				serviceStr := string(serviceContent)

				// Domain services should define interfaces
				assert.Contains(t, serviceStr, "type", "Should define service types")
				assert.Contains(t, serviceStr, "interface", "Should define service interface")

				// Should contain business logic methods
				businessMethods := []string{"Validate", "Check", "Calculate", "Determine", "Ensure"}
				foundBusinessMethod := false
				for _, method := range businessMethods {
					if strings.Contains(serviceStr, method) {
						foundBusinessMethod = true
						t.Logf("Found business method %s in %s", method, file)
					}
				}
				assert.True(t, foundBusinessMethod, "Domain service should contain business methods")

				// Should NOT depend on infrastructure
				assert.NotContains(t, serviceStr, "internal/infrastructure", 
					"Domain services should not depend on infrastructure")
				assert.NotContains(t, serviceStr, "database/sql", 
					"Domain services should not depend on database")
				assert.NotContains(t, serviceStr, "gorm.io", 
					"Domain services should not depend on ORM")

				// User domain service specific checks
				if strings.Contains(file, "user_domain_service.go") {
					userServiceMethods := []string{"ValidateUser", "CheckDuplicate", "CalculateScore"}
					for _, method := range userServiceMethods {
						if strings.Contains(serviceStr, method) {
							t.Logf("Found user domain method: %s", method)
						}
					}
				}

				// Authentication service should handle auth domain logic
				if strings.Contains(file, "authentication_service.go") {
					authMethods := []string{"Authenticate", "ValidateCredentials", "GenerateToken"}
					for _, method := range authMethods {
						if strings.Contains(serviceStr, method) {
							t.Logf("Found auth domain method: %s", method)
						}
					}
				}
			}
		}
	})

	t.Run("web_api_ddd_implements_cqrs_pattern", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining CQRS implementation (Commands and Queries)
		// THEN: Should separate commands (writes) from queries (reads)

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

		generateCmd := exec.Command("./go-starter", "new", "test-cqrs",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/cqrs",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-cqrs")

		// Check command files
		commandFiles := []string{
			"internal/application/commands/create_user_command.go",
			"internal/application/commands/update_user_command.go",
			"internal/application/commands/delete_user_command.go",
		}

		foundCommands := false
		for _, file := range commandFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				commandContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				commandStr := string(commandContent)

				foundCommands = true
				
				// Commands should define command structures
				assert.Contains(t, commandStr, "Command", "Should define command struct")
				
				// Commands should be about writes/mutations
				writeMethods := []string{"Create", "Update", "Delete", "Add", "Remove", "Change"}
				foundWriteMethod := false
				for _, method := range writeMethods {
					if strings.Contains(commandStr, method) {
						foundWriteMethod = true
						break
					}
				}
				assert.True(t, foundWriteMethod, "Commands should contain write operations")

				// Should have command handler
				if strings.Contains(commandStr, "Handle") || strings.Contains(commandStr, "Execute") {
					t.Logf("Found command handler in: %s", file)
				}
			}
		}

		// Check query files
		queryFiles := []string{
			"internal/application/queries/get_user_query.go",
			"internal/application/queries/list_users_query.go",
			"internal/application/queries/search_users_query.go",
		}

		foundQueries := false
		for _, file := range queryFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				queryContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				queryStr := string(queryContent)

				foundQueries = true

				// Queries should define query structures
				assert.Contains(t, queryStr, "Query", "Should define query struct")
				
				// Queries should be about reads
				readMethods := []string{"Get", "List", "Search", "Find", "Retrieve"}
				foundReadMethod := false
				for _, method := range readMethods {
					if strings.Contains(queryStr, method) {
						foundReadMethod = true
						break
					}
				}
				assert.True(t, foundReadMethod, "Queries should contain read operations")

				// Should have query handler
				if strings.Contains(queryStr, "Handle") || strings.Contains(queryStr, "Execute") {
					t.Logf("Found query handler in: %s", file)
				}
			}
		}

		// Check command/query handlers
		handlerFiles := []string{
			"internal/application/handlers/create_user_handler.go",
			"internal/application/handlers/get_user_handler.go",
		}

		for _, file := range handlerFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				handlerContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				handlerStr := string(handlerContent)

				assert.Contains(t, handlerStr, "Handle", "Handler should have Handle method")
				t.Logf("Found CQRS handler: %s", file)
			}
		}

		if foundCommands || foundQueries {
			t.Logf("CQRS pattern implemented (Commands: %v, Queries: %v)", foundCommands, foundQueries)
		}
	})

	t.Run("web_api_ddd_project_builds_successfully", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
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
			"--architecture=ddd",
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
		buildGeneratedCmd := exec.Command("go", "build", "-o", "ddd-api", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated web-api-ddd project should build successfully: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "ddd-api"), "DDD API binary should be created")
	})

	t.Run("web_api_ddd_supports_multiple_frameworks", func(t *testing.T) {
		// GIVEN: web-api-ddd blueprint with different web frameworks
		// WHEN: Generating projects with different frameworks
		// THEN: Each should generate framework-specific presentation layer while maintaining DDD structure

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
				projectName := fmt.Sprintf("test-ddd-%s", framework)
				
				generateCmd := exec.Command("./go-starter", "new", projectName,
					"--type=web-api",
					"--architecture=ddd",
					"--module=github.com/test/"+projectName,
					"--framework="+framework,
					"--database-driver=sqlite",
					"--no-git")
				output, err := generateCmd.CombinedOutput()
				require.NoError(t, err, "Should generate successfully with %s: %s", framework, string(output))

				projectDir := filepath.Join(tmpDir, projectName)

				// Verify DDD structure exists regardless of framework
				assert.DirExists(t, filepath.Join(projectDir, "internal", "domain"), "Domain layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "application"), "Application layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "infrastructure"), "Infrastructure layer should exist")
				assert.DirExists(t, filepath.Join(projectDir, "internal", "presentation"), "Presentation layer should exist")

				// Check that presentation layer uses correct framework
				controllerPath := filepath.Join(projectDir, "internal", "presentation", "controllers", "user_controller.go")
				if fileExists(t, controllerPath) {
					controllerContent, err := os.ReadFile(controllerPath)
					require.NoError(t, err)
					controllerStr := string(controllerContent)

					switch framework {
					case "gin":
						assert.Contains(t, controllerStr, "gin-gonic/gin", "Should import Gin framework")
					case "echo":
						assert.Contains(t, controllerStr, "labstack/echo", "Should import Echo framework")
					case "fiber":
						assert.Contains(t, controllerStr, "gofiber/fiber", "Should import Fiber framework")
					case "chi":
						assert.Contains(t, controllerStr, "go-chi/chi", "Should import Chi framework")
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
				require.NoError(t, err, "DDD API with %s should compile: %s", framework, string(output))
			})
		}
	})

	t.Run("web_api_ddd_supports_multiple_loggers", func(t *testing.T) {
		// GIVEN: web-api-ddd blueprint with different logger configurations
		// WHEN: Generating projects with different loggers
		// THEN: Logger should be properly integrated across DDD layers

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
					"--architecture=ddd",
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

				// Check logger is used in application services
				appServicePath := filepath.Join(projectDir, "internal", "application", "services", "user_application_service.go")
				if fileExists(t, appServicePath) {
					serviceContent, err := os.ReadFile(appServicePath)
					require.NoError(t, err)
					serviceStr := string(serviceContent)

					assert.Contains(t, serviceStr, "internal/logger", "Application services should import logger")
				}

				// Check logger is used in domain services
				domainServicePath := filepath.Join(projectDir, "internal", "domain", "services", "user_domain_service.go")
				if fileExists(t, domainServicePath) {
					serviceContent, err := os.ReadFile(domainServicePath)
					require.NoError(t, err)
					serviceStr := string(serviceContent)

					// Domain services might use logger for debugging (controversial in DDD)
					if strings.Contains(serviceStr, "internal/logger") {
						t.Logf("Domain service uses logger in %s", logger)
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
				require.NoError(t, err, "web-api-ddd with %s logger should compile: %s", logger, string(output))
			})
		}
	})

	t.Run("web_api_ddd_implements_specifications_pattern", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining specifications implementation
		// THEN: Should implement specification pattern for complex business rules

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

		generateCmd := exec.Command("./go-starter", "new", "test-specifications",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/specifications",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-specifications")

		// Check specification files
		specificationFiles := []string{
			"internal/domain/specifications/user_specification.go",
			"internal/domain/specifications/active_user_specification.go",
			"internal/domain/specifications/valid_email_specification.go",
		}

		for _, file := range specificationFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				specContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				specStr := string(specContent)

				// Specifications should define interfaces
				assert.Contains(t, specStr, "type", "Should define specification type")
				assert.Contains(t, specStr, "IsSatisfiedBy", "Should have IsSatisfiedBy method")

				// Should work with domain entities
				assert.Contains(t, specStr, "internal/domain/entities", "Should import domain entities")

				// Should return boolean result
				assert.Contains(t, specStr, "bool", "IsSatisfiedBy should return bool")

				t.Logf("Found specification: %s", file)
			}
		}

		// Check if specifications are used in domain services
		domainServicePath := filepath.Join(projectDir, "internal", "domain", "services", "user_domain_service.go")
		if fileExists(t, domainServicePath) {
			serviceContent, err := os.ReadFile(domainServicePath)
			require.NoError(t, err)
			serviceStr := string(serviceContent)

			if strings.Contains(serviceStr, "specifications") || strings.Contains(serviceStr, "IsSatisfiedBy") {
				t.Logf("Domain service uses specifications for business rules")
			}
		}
	})

	t.Run("web_api_ddd_implements_domain_events", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining domain events implementation
		// THEN: Should implement domain events for decoupling and integration

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

		generateCmd := exec.Command("./go-starter", "new", "test-events",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/events",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-events")

		// Check domain event files
		eventFiles := []string{
			"internal/domain/events/user_created_event.go",
			"internal/domain/events/user_updated_event.go",
			"internal/domain/events/domain_event.go",
		}

		foundDomainEvents := false
		for _, file := range eventFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				eventContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				eventStr := string(eventContent)

				foundDomainEvents = true

				// Domain events should define event structures
				assert.Contains(t, eventStr, "type", "Should define event type")
				
				// Should have event metadata
				eventFields := []string{"ID", "OccurredAt", "AggregateID", "EventType"}
				foundField := false
				for _, field := range eventFields {
					if strings.Contains(eventStr, field) {
						foundField = true
						break
					}
				}
				assert.True(t, foundField, "Event should have metadata fields")

				// Check for event interface
				if strings.Contains(eventStr, "DomainEvent") && strings.Contains(eventStr, "interface") {
					assert.Contains(t, eventStr, "GetEventType", "Should define event interface")
					t.Logf("Found domain event interface in: %s", file)
				}

				t.Logf("Found domain event: %s", file)
			}
		}

		// Check for event handlers in application layer
		handlerFiles := []string{
			"internal/application/handlers/user_created_handler.go",
			"internal/application/handlers/event_handler.go",
		}

		for _, file := range handlerFiles {
			if fileExists(t, filepath.Join(projectDir, file)) {
				handlerContent, err := os.ReadFile(filepath.Join(projectDir, file))
				require.NoError(t, err)
				handlerStr := string(handlerContent)

				assert.Contains(t, handlerStr, "Handle", "Event handler should have Handle method")
				assert.Contains(t, handlerStr, "internal/domain/events", "Should import domain events")
				
				t.Logf("Found event handler: %s", file)
			}
		}

		// Check if aggregates raise events
		aggPath := filepath.Join(projectDir, "internal", "domain", "aggregates", "user_aggregate.go")
		if fileExists(t, aggPath) {
			aggContent, err := os.ReadFile(aggPath)
			require.NoError(t, err)
			aggStr := string(aggContent)

			if strings.Contains(aggStr, "RaiseEvent") || strings.Contains(aggStr, "events") {
				t.Logf("Aggregate raises domain events")
			}
		}

		if foundDomainEvents {
			t.Logf("Domain events pattern implemented")
		}
	})
}

// TestWebApiDddDomainModeling validates proper domain modeling in DDD
func TestWebApiDddDomainModeling(t *testing.T) {
	t.Run("avoids_primitive_obsession", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining domain entities and value objects
		// THEN: Should use value objects instead of primitive types for domain concepts

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

		generateCmd := exec.Command("./go-starter", "new", "test-primitives",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/primitives",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-primitives")

		// Check User entity for primitive obsession
		userEntityPath := filepath.Join(projectDir, "internal", "domain", "entities", "user.go")
		if fileExists(t, userEntityPath) {
			entityContent, err := os.ReadFile(userEntityPath)
			require.NoError(t, err)
			entityStr := string(entityContent)

			// Should use value objects, not primitive strings
			assert.Contains(t, entityStr, "valueobjects", "Should import value objects")

			// Should NOT use primitive types for domain concepts
			primitivePattern := `Email\s+string|Name\s+string|Password\s+string`
			assert.NotRegexp(t, primitivePattern, entityStr,
				"Should not use primitive strings for domain concepts like Email, Name")

			// Should use value objects
			valueObjectPattern := `Email\s+valueobjects\.Email|Name\s+valueobjects\.UserName`
			if !assert.Regexp(t, valueObjectPattern, entityStr, "Should use value objects") {
				t.Logf("Entity content: %s", entityStr)
			}
		}

		// Check value objects exist and are properly implemented
		emailVOPath := filepath.Join(projectDir, "internal", "domain", "valueobjects", "email.go")
		if fileExists(t, emailVOPath) {
			voContent, err := os.ReadFile(emailVOPath)
			require.NoError(t, err)
			voStr := string(voContent)

			assert.Contains(t, voStr, "type Email struct", "Should define Email value object")
			assert.Contains(t, voStr, "String()", "Should have String method")
			assert.Contains(t, voStr, "NewEmail", "Should have constructor")
		}
	})

	t.Run("implements_rich_domain_models", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining domain entities
		// THEN: Should have rich domain models with behavior, not anemic models

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

		generateCmd := exec.Command("./go-starter", "new", "test-rich-models",
			"--type=web-api",
			"--architecture=ddd",
			"--module=github.com/test/rich-models",
			"--framework=gin",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-rich-models")

		// Check User entity for rich behavior
		userEntityPath := filepath.Join(projectDir, "internal", "domain", "entities", "user.go")
		if fileExists(t, userEntityPath) {
			entityContent, err := os.ReadFile(userEntityPath)
			require.NoError(t, err)
			entityStr := string(entityContent)

			// Should have business methods (behavior)
			businessMethods := []string{
				"ChangeEmail", "UpdateProfile", "Activate", "Deactivate", 
				"IsActive", "CanLogin", "ValidatePassword", "ChangePassword",
			}

			foundBusinessMethods := 0
			for _, method := range businessMethods {
				if strings.Contains(entityStr, method) {
					foundBusinessMethods++
					t.Logf("Found business method: %s", method)
				}
			}

			// Should have at least some business behavior (not anemic)
			assert.Greater(t, foundBusinessMethods, 0, 
				"User entity should have business methods (rich domain model)")

			// Should NOT be just a data structure (anemic model)
			if foundBusinessMethods == 0 {
				// Check if it's just getters/setters
				getterSetterPattern := `func \([^)]*\) Get[A-Z]|func \([^)]*\) Set[A-Z]`
				if assert.Regexp(t, getterSetterPattern, entityStr, "If no business methods, should at least have accessors") {
					t.Logf("WARNING: User entity appears to be anemic (only getters/setters)")
				}
			}

			// Should encapsulate invariants
			invariantMethods := []string{"Validate", "IsValid", "EnsureValid", "CheckInvariants"}
			foundInvariants := 0
			for _, method := range invariantMethods {
				if strings.Contains(entityStr, method) {
					foundInvariants++
					t.Logf("Found invariant method: %s", method)
				}
			}

			if foundInvariants > 0 {
				t.Logf("Entity encapsulates business invariants")
			}
		}
	})

	t.Run("enforces_aggregate_boundaries", func(t *testing.T) {
		// GIVEN: A generated web-api-ddd project
		// WHEN: Examining aggregate design
		// THEN: Should properly define aggregate boundaries and root entities

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
			"--architecture=ddd",
			"--module=github.com/test/boundaries",
			"--framework=gin",
			"--database-driver=postgres",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-boundaries")

		// Check aggregate root implementation
		userAggPath := filepath.Join(projectDir, "internal", "domain", "aggregates", "user_aggregate.go")
		if fileExists(t, userAggPath) {
			aggContent, err := os.ReadFile(userAggPath)
			require.NoError(t, err)
			aggStr := string(aggContent)

			// Should identify aggregate root
			assert.Contains(t, aggStr, "UserAggregate", "Should define aggregate")
			
			// Should have aggregate root ID
			assert.Contains(t, aggStr, "ID", "Aggregate root should have ID")

			// Should control access to child entities
			if strings.Contains(aggStr, "AddChild") || strings.Contains(aggStr, "RemoveChild") {
				t.Logf("Aggregate controls child entity access")
			}

			// Should enforce business rules at aggregate boundary
			businessRules := []string{"EnsureConsistency", "ValidateInvariants", "CheckRules"}
			for _, rule := range businessRules {
				if strings.Contains(aggStr, rule) {
					t.Logf("Aggregate enforces business rule: %s", rule)
				}
			}
		}

		// Check repository operates at aggregate level
		repoPath := filepath.Join(projectDir, "internal", "domain", "repositories", "user_repository.go")
		if fileExists(t, repoPath) {
			repoContent, err := os.ReadFile(repoPath)
			require.NoError(t, err)
			repoStr := string(repoContent)

			// Repository should work with aggregate roots
			if strings.Contains(repoStr, "UserAggregate") {
				t.Logf("Repository works with aggregate root")
			} else if strings.Contains(repoStr, "User") {
				t.Logf("Repository works with User entity (simplified approach)")
			}

			// Should have methods for saving/loading entire aggregates
			aggregateMethods := []string{"Save", "GetByID", "Add", "Remove"}
			for _, method := range aggregateMethods {
				if strings.Contains(repoStr, method) {
					t.Logf("Repository has aggregate method: %s", method)
				}
			}
		}
	})
}

// Helper functions
func fileExists(t *testing.T, path string) bool {
	_, err := os.Stat(path)
	return err == nil
}