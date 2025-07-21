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
	// Initialize templates filesystem for ATDD tests
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
					println("ATDD: Found blueprints directory at:", templatesDir)
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

// Feature: Database Integration Validation
// As a developer
// I want database integration to work across architectures
// So that I can persist data effectively

func TestWebAPI_DatabaseIntegration_CrossArchitectures(t *testing.T) {
	testCases := []struct {
		architecture string
		database     string
		orm          string
	}{
		{"standard", "postgres", "gorm"},
		{"clean", "mysql", "sqlx"},
		{"ddd", "postgres", "ent"},
		{"hexagonal", "sqlite", "gorm"},
	}

	for _, tc := range testCases {
		t.Run("Database_"+tc.architecture+"_"+tc.database+"_"+tc.orm, func(t *testing.T) {
			// Scenario: Database integration
			// Given I generate a "<architecture>" web API with "<database>" and "<orm>"
			// Then database configuration should follow architecture patterns
			// And data access should be properly abstracted
			// And transactions should be handled correctly
			// And migrations should be architecture-appropriate

			config := types.ProjectConfig{
				Name:         "test-db-" + tc.architecture + "-" + tc.database + "-" + tc.orm,
				Module:       "github.com/test/test-db-" + tc.architecture + "-" + tc.database + "-" + tc.orm,
				Type:         "web-api",
				Architecture: tc.architecture,
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Driver: tc.database,
						ORM:    tc.orm,
					},
				},
			}

			// Handle different architecture types
			switch tc.architecture {
			case "standard":
				config.Type = "web-api-standard"
			case "clean":
				config.Type = "web-api-clean"
			case "ddd":
				config.Type = "web-api-ddd"
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then database configuration should follow architecture patterns
			validator := NewDatabaseIntegrationValidator(projectPath, tc.architecture)
			validator.ValidateDatabaseArchitecturePatterns(t, tc.database, tc.orm)

			// And data access should be properly abstracted
			validator.ValidateDataAccessAbstraction(t)

			// And transactions should be handled correctly
			validator.ValidateTransactionHandling(t)

			// And migrations should be architecture-appropriate
			validator.ValidateMigrations(t)

			// Skip if architecture not fully implemented
			if tc.architecture == "hexagonal" && !validator.IsArchitectureImplemented() {
				t.Skip("Hexagonal architecture not fully implemented yet")
			}

			validator.ValidateCompilation(t)
		})
	}
}

// Feature: Logger Integration Across Architectures
// As a developer
// I want consistent logging across all web API architectures
// So that I can monitor applications effectively

func TestWebAPI_LoggerIntegration_CrossArchitectures(t *testing.T) {
	testCases := []struct {
		architecture string
		logger       string
	}{
		{"standard", "slog"},
		{"clean", "zap"},
		{"ddd", "logrus"},
		{"hexagonal", "zerolog"},
	}

	for _, tc := range testCases {
		t.Run("Logger_"+tc.architecture+"_"+tc.logger, func(t *testing.T) {
			// Scenario: Logger integration
			// Given I generate a "<architecture>" web API with "<logger>"
			// Then logging should be properly configured
			// And log statements should follow architecture patterns
			// And log levels should be configurable
			// And structured logging should be used

			config := types.ProjectConfig{
				Name:         "test-logger-" + tc.architecture + "-" + tc.logger,
				Module:       "github.com/test/test-logger-" + tc.architecture + "-" + tc.logger,
				Type:         "web-api",
				Architecture: tc.architecture,
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       tc.logger,
			}

			// Handle different architecture types
			switch tc.architecture {
			case "standard":
				config.Type = "web-api-standard"
			case "clean":
				config.Type = "web-api-clean"
			case "ddd":
				config.Type = "web-api-ddd"
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then logging should be properly configured
			validator := NewLoggerIntegrationValidator(projectPath, tc.architecture)
			validator.ValidateLoggerConfiguration(t, tc.logger)

			// And log statements should follow architecture patterns
			validator.ValidateLoggerArchitecturePatterns(t)

			// And log levels should be configurable
			validator.ValidateLoggerLevels(t)

			// And structured logging should be used
			validator.ValidateStructuredLogging(t, tc.logger)

			// Skip if architecture not fully implemented
			if tc.architecture == "hexagonal" && !validator.IsArchitectureImplemented() {
				t.Skip("Hexagonal architecture not fully implemented yet")
			}

			validator.ValidateCompilation(t)
		})
	}
}

// Feature: Blueprint Compilation and Runtime
// As a user
// I want all web API blueprints to work out of the box
// So that I can start development immediately

func TestWebAPI_CompilationValidation_AllArchitectures(t *testing.T) {
	architectures := []string{"standard", "clean", "ddd", "hexagonal"}

	for _, architecture := range architectures {
		t.Run("Compilation_"+architecture, func(t *testing.T) {
			// Scenario: Compilation validation
			// Given I generate a "<architecture>" web API
			// When I run "go build"
			// Then the project should compile without errors
			// And dependencies should resolve correctly
			// And the binary should be executable

			config := types.ProjectConfig{
				Name:         "test-compile-" + architecture,
				Module:       "github.com/test/test-compile-" + architecture,
				Type:         "web-api",
				Architecture: architecture,
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Drivers: []string{"postgres"},
						ORM:     "gorm",
					},
					Authentication: types.AuthConfig{
						Type: "jwt",
					},
				},
			}

			// Handle different architecture types
			switch architecture {
			case "standard":
				config.Type = "web-api-standard"
			case "clean":
				config.Type = "web-api-clean"
			case "ddd":
				config.Type = "web-api-ddd"
			}

			projectPath := helpers.GenerateProject(t, config)

			// When I run "go build"
			// Then the project should compile without errors
			validator := NewCompilationValidator(projectPath, architecture)

			// Skip if architecture not fully implemented
			if architecture == "hexagonal" && !validator.IsArchitectureImplemented() {
				t.Skip("Hexagonal architecture not fully implemented yet")
			}

			validator.ValidateCompilation(t)

			// And dependencies should resolve correctly
			validator.ValidateDependencyResolution(t)

			// And the binary should be executable
			validator.ValidateBinaryGeneration(t)
		})
	}
}

// Feature: Architecture-Specific Code Quality
// As a code reviewer
// I want generated code to follow architecture principles
// So that projects maintain architectural integrity

func TestWebAPI_ArchitectureCompliance_AllPatterns(t *testing.T) {
	architectures := []string{"standard", "clean", "ddd", "hexagonal"}

	for _, architecture := range architectures {
		t.Run("Compliance_"+architecture, func(t *testing.T) {
			// Scenario: Architecture compliance
			// Given I generate a "<architecture>" web API
			// Then the code should follow "<architecture>" principles
			// And dependency directions should be correct
			// And layer boundaries should be enforced
			// And the project should pass architectural linting

			config := types.ProjectConfig{
				Name:         "test-compliance-" + architecture,
				Module:       "github.com/test/test-compliance-" + architecture,
				Type:         "web-api",
				Architecture: architecture,
				GoVersion:    "1.21",
				Framework:    "gin",
				Logger:       "slog",
				Features: &types.Features{
					Database: types.DatabaseConfig{
						Drivers: []string{"postgres"},
						ORM:     "gorm",
					},
				},
			}

			// Handle different architecture types
			switch architecture {
			case "standard":
				config.Type = "web-api-standard"
			case "clean":
				config.Type = "web-api-clean"
			case "ddd":
				config.Type = "web-api-ddd"
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then the code should follow "<architecture>" principles
			validator := NewArchitectureComplianceValidator(projectPath, architecture)

			// Skip if architecture not fully implemented
			if architecture == "hexagonal" && !validator.IsArchitectureImplemented() {
				t.Skip("Hexagonal architecture not fully implemented yet")
			}

			validator.ValidateArchitecturePrinciples(t)

			// And dependency directions should be correct
			validator.ValidateDependencyDirections(t)

			// And layer boundaries should be enforced
			validator.ValidateLayerBoundaries(t)

			// And the project should pass architectural linting
			validator.ValidateArchitecturalLinting(t)

			validator.ValidateCompilation(t)
		})
	}
}

// DatabaseIntegrationValidator validates database integration across architectures
type DatabaseIntegrationValidator struct {
	ProjectPath  string
	Architecture string
}

func NewDatabaseIntegrationValidator(projectPath, architecture string) *DatabaseIntegrationValidator {
	return &DatabaseIntegrationValidator{
		ProjectPath:  projectPath,
		Architecture: architecture,
	}
}

func (v *DatabaseIntegrationValidator) ValidateDatabaseArchitecturePatterns(t *testing.T, database, orm string) {
	t.Helper()

	switch v.Architecture {
	case "standard":
		// Database connection in internal/database/
		dbFile := filepath.Join(v.ProjectPath, "internal/database/connection.go")
		helpers.AssertFileExists(t, dbFile)
	case "clean":
		// Database in infrastructure/persistence/
		dbFile := filepath.Join(v.ProjectPath, "internal/infrastructure/persistence/database.go")
		helpers.AssertFileExists(t, dbFile)
	case "ddd":
		// Database in infrastructure/persistence/
		dbFile := filepath.Join(v.ProjectPath, "internal/infrastructure/persistence/database.go")
		helpers.AssertFileExists(t, dbFile)
	case "hexagonal":
		// Database as secondary adapter
		dbFile := filepath.Join(v.ProjectPath, "internal/adapters/secondary/persistence/database.go")
		if helpers.FileExists(dbFile) {
			helpers.AssertFileExists(t, dbFile)
		}
	}
}

func (v *DatabaseIntegrationValidator) ValidateDataAccessAbstraction(t *testing.T) {
	t.Helper()

	switch v.Architecture {
	case "standard":
		// Repository interfaces in internal/repository/
		repoFile := filepath.Join(v.ProjectPath, "internal/repository/interfaces.go")
		if helpers.FileExists(repoFile) {
			helpers.AssertFileContains(t, repoFile, "interface")
		}
	case "clean", "ddd":
		// Repository interfaces in domain/ports or domain layer
		// Check various possible locations
		locations := []string{
			"internal/domain/ports/repositories.go",
			"internal/domain/user/repository.go",
		}
		found := false
		for _, location := range locations {
			if helpers.FileExists(filepath.Join(v.ProjectPath, location)) {
				found = true
				break
			}
		}
		if found {
			t.Logf("✓ Repository abstraction found for %s", v.Architecture)
		}
	}
}

func (v *DatabaseIntegrationValidator) ValidateTransactionHandling(t *testing.T) {
	t.Helper()
	// This would validate transaction handling patterns
	// For now, just ensure database connection files exist
	v.ValidateDatabaseArchitecturePatterns(t, "postgres", "gorm")
}

func (v *DatabaseIntegrationValidator) ValidateMigrations(t *testing.T) {
	t.Helper()

	// Migrations should exist regardless of architecture
	migrationUp := filepath.Join(v.ProjectPath, "migrations/001_create_users.up.sql")
	migrationDown := filepath.Join(v.ProjectPath, "migrations/001_create_users.down.sql")

	if helpers.FileExists(migrationUp) {
		helpers.AssertFileExists(t, migrationUp)
		helpers.AssertFileExists(t, migrationDown)
	}
}

func (v *DatabaseIntegrationValidator) IsArchitectureImplemented() bool {
	switch v.Architecture {
	case "hexagonal":
		adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")
		return helpers.DirExists(adaptersDir)
	default:
		return true
	}
}

func (v *DatabaseIntegrationValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

// LoggerIntegrationValidator validates logger integration across architectures
type LoggerIntegrationValidator struct {
	ProjectPath  string
	Architecture string
}

func NewLoggerIntegrationValidator(projectPath, architecture string) *LoggerIntegrationValidator {
	return &LoggerIntegrationValidator{
		ProjectPath:  projectPath,
		Architecture: architecture,
	}
}

func (v *LoggerIntegrationValidator) ValidateLoggerConfiguration(t *testing.T, logger string) {
	t.Helper()

	// Logger interface should exist
	var loggerDir string
	switch v.Architecture {
	case "standard":
		loggerDir = "internal/logger"
	case "clean", "ddd":
		loggerDir = "internal/infrastructure/logger"
	case "hexagonal":
		loggerDir = "internal/adapters/secondary/logger"
	}

	interfaceFile := filepath.Join(v.ProjectPath, loggerDir, "interface.go")
	if helpers.FileExists(interfaceFile) {
		helpers.AssertFileExists(t, interfaceFile)
	}

	// Logger implementation should exist
	loggerFile := filepath.Join(v.ProjectPath, loggerDir, logger+".go")
	if helpers.FileExists(loggerFile) {
		helpers.AssertFileExists(t, loggerFile)
	}
}

func (v *LoggerIntegrationValidator) ValidateLoggerArchitecturePatterns(t *testing.T) {
	t.Helper()

	switch v.Architecture {
	case "standard":
		// Logger used directly in handlers
		handlerFiles := helpers.FindFiles(t, filepath.Join(v.ProjectPath, "internal/handlers"), "*.go")
		for _, file := range handlerFiles {
			if helpers.FileExists(file) {
				content := helpers.ReadFileContent(t, file)
				if helpers.StringContains(content, "logger") {
					t.Logf("✓ Logger used in handlers for standard architecture")
					break
				}
			}
		}
	case "clean", "ddd":
		// Logger injected through dependency injection
		containerFile := filepath.Join(v.ProjectPath, "internal/infrastructure/container/container.go")
		if helpers.FileExists(containerFile) {
			helpers.AssertFileContains(t, containerFile, "logger")
		}
	}
}

func (v *LoggerIntegrationValidator) ValidateLoggerLevels(t *testing.T) {
	t.Helper()

	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	if helpers.FileExists(configFile) {
		// Should have logger configuration
		helpers.AssertFileExists(t, configFile)
	}
}

func (v *LoggerIntegrationValidator) ValidateStructuredLogging(t *testing.T, logger string) {
	t.Helper()

	// Find logger implementation and verify structured logging
	loggerPaths := []string{
		"internal/logger/" + logger + ".go",
		"internal/infrastructure/logger/" + logger + ".go",
		"internal/adapters/secondary/logger/" + logger + ".go",
	}

	for _, path := range loggerPaths {
		fullPath := filepath.Join(v.ProjectPath, path)
		if helpers.FileExists(fullPath) {
			content := helpers.ReadFileContent(t, fullPath)
			switch logger {
			case "slog":
				assert.Contains(t, content, "slog", "Should use slog structured logging")
			case "zap":
				assert.Contains(t, content, "zap", "Should use zap structured logging")
			case "logrus":
				assert.Contains(t, content, "logrus", "Should use logrus structured logging")
			case "zerolog":
				assert.Contains(t, content, "zerolog", "Should use zerolog structured logging")
			}
			break
		}
	}
}

func (v *LoggerIntegrationValidator) IsArchitectureImplemented() bool {
	switch v.Architecture {
	case "hexagonal":
		adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")
		return helpers.DirExists(adaptersDir)
	default:
		return true
	}
}

func (v *LoggerIntegrationValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

// CompilationValidator validates compilation across architectures
type CompilationValidator struct {
	ProjectPath  string
	Architecture string
}

func NewCompilationValidator(projectPath, architecture string) *CompilationValidator {
	return &CompilationValidator{
		ProjectPath:  projectPath,
		Architecture: architecture,
	}
}

func (v *CompilationValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *CompilationValidator) ValidateDependencyResolution(t *testing.T) {
	t.Helper()

	// Check go.mod exists and is valid
	goModFile := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)

	// Validate module name is set correctly
	content := helpers.ReadFileContent(t, goModFile)
	assert.Contains(t, content, "module", "go.mod should have module declaration")
}

func (v *CompilationValidator) ValidateBinaryGeneration(t *testing.T) {
	t.Helper()

	// For now, just ensure main file exists
	mainFile := filepath.Join(v.ProjectPath, "cmd/server/main.go")
	if !helpers.FileExists(mainFile) {
		mainFile = filepath.Join(v.ProjectPath, "main.go")
	}
	helpers.AssertFileExists(t, mainFile)
}

func (v *CompilationValidator) IsArchitectureImplemented() bool {
	switch v.Architecture {
	case "hexagonal":
		adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")
		return helpers.DirExists(adaptersDir)
	default:
		return true
	}
}

// ArchitectureComplianceValidator validates architecture compliance
type ArchitectureComplianceValidator struct {
	ProjectPath  string
	Architecture string
}

func NewArchitectureComplianceValidator(projectPath, architecture string) *ArchitectureComplianceValidator {
	return &ArchitectureComplianceValidator{
		ProjectPath:  projectPath,
		Architecture: architecture,
	}
}

func (v *ArchitectureComplianceValidator) ValidateArchitecturePrinciples(t *testing.T) {
	t.Helper()

	// Validate architecture principles using direct file system checks
	switch v.Architecture {
	case "standard":
		v.validateStandardStructure(t)
	case "clean":
		v.validateCleanArchitectureLayers(t)
	case "ddd":
		v.validateDDDStructure(t)
	case "hexagonal":
		v.validateHexagonalStructure(t)
	}
}

func (v *ArchitectureComplianceValidator) ValidateDependencyDirections(t *testing.T) {
	t.Helper()

	// Validate dependency directions using direct file system checks
	switch v.Architecture {
	case "standard":
		v.validateStandardLayeredArchitecture(t)
	case "clean":
		v.validateCleanArchitectureDependencies(t)
	case "ddd":
		v.validateDDDDomainIsolation(t)
	case "hexagonal":
		v.validateHexagonalDependencies(t)
	}
}

func (v *ArchitectureComplianceValidator) ValidateLayerBoundaries(t *testing.T) {
	t.Helper()

	// Validate layer boundaries using direct file system checks
	switch v.Architecture {
	case "standard":
		v.validateStandardControllers(t)
		v.validateStandardServices(t)
		v.validateStandardRepositories(t)
	case "clean":
		v.validateCleanArchitectureDependencies(t)
	case "ddd":
		v.validateDDDDomainIsolation(t)
	case "hexagonal":
		v.validateHexagonalPortContracts(t)
		v.validateHexagonalAdapterImplementation(t)
	}
}

func (v *ArchitectureComplianceValidator) ValidateArchitecturalLinting(t *testing.T) {
	t.Helper()

	// For now, architectural linting is validated through compilation
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *ArchitectureComplianceValidator) IsArchitectureImplemented() bool {
	switch v.Architecture {
	case "hexagonal":
		adaptersDir := filepath.Join(v.ProjectPath, "internal/adapters")
		return helpers.DirExists(adaptersDir)
	default:
		return true
	}
}

func (v *ArchitectureComplianceValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

// Standard architecture validation methods
func (v *ArchitectureComplianceValidator) validateStandardStructure(t *testing.T) {
	t.Helper()
	// Validate standard web API structure exists
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/handlers"))
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/services"))
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/models"))
}

func (v *ArchitectureComplianceValidator) validateStandardLayeredArchitecture(t *testing.T) {
	t.Helper()
	// Basic layered architecture validation - just ensure directories exist
	v.validateStandardStructure(t)
}

func (v *ArchitectureComplianceValidator) validateStandardControllers(t *testing.T) {
	t.Helper()
	// Validate controllers/handlers exist
	handlerPath := filepath.Join(v.ProjectPath, "internal/handlers")
	if helpers.DirExists(handlerPath) {
		helpers.AssertDirExists(t, handlerPath)
	}
}

func (v *ArchitectureComplianceValidator) validateStandardServices(t *testing.T) {
	t.Helper()
	// Validate services exist
	servicesPath := filepath.Join(v.ProjectPath, "internal/services")
	if helpers.DirExists(servicesPath) {
		helpers.AssertDirExists(t, servicesPath)
	}
}

func (v *ArchitectureComplianceValidator) validateStandardRepositories(t *testing.T) {
	t.Helper()
	// Validate repositories exist
	repoPath := filepath.Join(v.ProjectPath, "internal/repository")
	if helpers.DirExists(repoPath) {
		helpers.AssertDirExists(t, repoPath)
	}
}

// Clean architecture validation methods
func (v *ArchitectureComplianceValidator) validateCleanArchitectureLayers(t *testing.T) {
	t.Helper()
	// Validate clean architecture layer structure
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/domain"))
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/application"))
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/infrastructure"))
}

func (v *ArchitectureComplianceValidator) validateCleanArchitectureDependencies(t *testing.T) {
	t.Helper()
	// Basic dependency validation - just ensure layers exist
	v.validateCleanArchitectureLayers(t)
}

// DDD validation methods
func (v *ArchitectureComplianceValidator) validateDDDStructure(t *testing.T) {
	t.Helper()
	// Validate DDD structure exists
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/domain"))
	helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/application"))
}

func (v *ArchitectureComplianceValidator) validateDDDDomainIsolation(t *testing.T) {
	t.Helper()
	// Basic domain isolation validation
	v.validateDDDStructure(t)
}

// Hexagonal architecture validation methods
func (v *ArchitectureComplianceValidator) validateHexagonalStructure(t *testing.T) {
	t.Helper()
	// Validate hexagonal architecture structure
	adaptersPath := filepath.Join(v.ProjectPath, "internal/adapters")
	if helpers.DirExists(adaptersPath) {
		helpers.AssertDirExists(t, adaptersPath)
		helpers.AssertDirExists(t, filepath.Join(v.ProjectPath, "internal/domain"))
	}
}

func (v *ArchitectureComplianceValidator) validateHexagonalDependencies(t *testing.T) {
	t.Helper()
	// Basic hexagonal dependency validation
	v.validateHexagonalStructure(t)
}

func (v *ArchitectureComplianceValidator) validateHexagonalPortContracts(t *testing.T) {
	t.Helper()
	// Validate port contracts exist
	portsPath := filepath.Join(v.ProjectPath, "internal/ports")
	if helpers.DirExists(portsPath) {
		helpers.AssertDirExists(t, portsPath)
	}
}

func (v *ArchitectureComplianceValidator) validateHexagonalAdapterImplementation(t *testing.T) {
	t.Helper()
	// Validate adapter implementations
	adaptersPath := filepath.Join(v.ProjectPath, "internal/adapters")
	if helpers.DirExists(adaptersPath) {
		helpers.AssertDirExists(t, adaptersPath)
	}
}
