package cli_test

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
					templates.SetTemplatesFS(os.DirFS(templatesDir))
					return
				}
			}
		}

		// Move up one directory
		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// Feature: Standard CLI Blueprint
// As a developer
// I want to generate a CLI application project
// So that I can quickly build command-line tools

func TestStandard_CLI_BasicGeneration_WithCobra(t *testing.T) {
	// Scenario: Generate standard CLI with Cobra
	// Given I want a standard CLI application
	// When I generate with framework "cobra"
	// Then the project should include Cobra command setup
	// And the project should have root and version commands
	// And the project should compile and run successfully
	// And CLI should show help and version output

	// Given I want a standard CLI application
	config := types.ProjectConfig{
		Name:      "test-standard-cli",
		Module:    "github.com/test/test-standard-cli",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	// When I generate with framework "cobra"
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include Cobra command setup
	validator := NewCLIValidator(projectPath, "standard")
	validator.ValidateCobraSetup(t)

	// And the project should have root and version commands
	validator.ValidateRootCommand(t)
	validator.ValidateVersionCommand(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)

	// And CLI should show help and version output
	validator.ValidateHelpOutput(t)
	validator.ValidateVersionOutput(t)
}

func TestStandard_CLI_WithDifferentLoggers(t *testing.T) {
	// Scenario: Generate CLI with different logging libraries
	// Given I want a CLI application with configurable logging
	// When I generate with different loggers
	// Then the project should include the selected logger
	// And the project should compile successfully
	// And logging should work as expected

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("logger_"+logger, func(t *testing.T) {
			// Given I want a CLI application with configurable logging
			config := types.ProjectConfig{
				Name:      "test-cli-" + logger,
				Module:    "github.com/test/test-cli-" + logger,
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    logger,
			}

			// When I generate with different loggers
			projectPath := helpers.GenerateProject(t, config)

			// Then the project should include the selected logger
			validator := NewCLIValidator(projectPath, "standard")
			validator.ValidateLogger(t, logger)

			// And the project should compile successfully
			validator.ValidateCompilation(t)

			// And logging should work as expected
			validator.ValidateLoggerFunctionality(t, logger)
		})
	}
}

func TestStandard_CLI_ConfigurationSupport(t *testing.T) {
	// Scenario: Generate CLI with configuration support
	// Given I want a CLI application with configuration
	// When I generate the project
	// Then the project should include Viper configuration
	// And the project should have config file support
	// And the project should compile and run successfully

	// Given I want a CLI application with configuration
	config := types.ProjectConfig{
		Name:      "test-cli-config",
		Module:    "github.com/test/test-cli-config",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include Viper configuration
	validator := NewCLIValidator(projectPath, "standard")
	validator.ValidateViperConfiguration(t)

	// And the project should have config file support
	validator.ValidateConfigFileSupport(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)
}

func TestStandard_CLI_DockerSupport(t *testing.T) {
	// Scenario: Generate CLI with Docker support
	// Given I want a CLI application with Docker
	// When I generate the project
	// Then the project should include Dockerfile
	// And the project should have Makefile with Docker targets
	// And the project should compile and run successfully

	// Given I want a CLI application with Docker
	config := types.ProjectConfig{
		Name:      "test-cli-docker",
		Module:    "github.com/test/test-cli-docker",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include Dockerfile
	validator := NewCLIValidator(projectPath, "standard")
	validator.ValidateDockerSupport(t)

	// And the project should have Makefile with Docker targets
	validator.ValidateMakefileTargets(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)
}

func TestStandard_CLI_TestSupport(t *testing.T) {
	// Scenario: Generate CLI with test support
	// Given I want a CLI application with tests
	// When I generate the project
	// Then the project should include test files
	// And the project should use testify for assertions
	// And the tests should run successfully

	// Given I want a CLI application with tests
	config := types.ProjectConfig{
		Name:      "test-cli-tests",
		Module:    "github.com/test/test-cli-tests",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	// When I generate the project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should include test files
	validator := NewCLIValidator(projectPath, "standard")
	validator.ValidateTestFiles(t)

	// And the project should use testify for assertions
	validator.ValidateTestifyUsage(t)

	// And the tests should run successfully
	validator.ValidateTestExecution(t)
}

// CLIValidator provides validation methods for CLI blueprints
type CLIValidator struct {
	projectPath  string
	architecture string
}

// NewCLIValidator creates a new CLIValidator
func NewCLIValidator(projectPath, architecture string) *CLIValidator {
	return &CLIValidator{
		projectPath:  projectPath,
		architecture: architecture,
	}
}

// ValidateCobraSetup validates Cobra framework setup
func (v *CLIValidator) ValidateCobraSetup(t *testing.T) {
	t.Helper()

	// Check if main.go exists and contains Cobra imports
	mainFile := filepath.Join(v.projectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	helpers.AssertFileContains(t, mainFile, "github.com/spf13/cobra")

	// Check if cmd/root.go exists
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "cobra.Command")
}

// ValidateRootCommand validates root command implementation
func (v *CLIValidator) ValidateRootCommand(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "rootCmd")
	helpers.AssertFileContains(t, rootFile, "Execute")
}

// ValidateVersionCommand validates version command implementation
func (v *CLIValidator) ValidateVersionCommand(t *testing.T) {
	t.Helper()

	versionFile := filepath.Join(v.projectPath, "cmd", "version.go")
	helpers.AssertFileExists(t, versionFile)
	helpers.AssertFileContains(t, versionFile, "versionCmd")
	helpers.AssertFileContains(t, versionFile, "Version")
}

// ValidateCompilation validates that the project compiles successfully
func (v *CLIValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertProjectCompiles(t, v.projectPath)
}

// ValidateHelpOutput validates CLI help output
func (v *CLIValidator) ValidateHelpOutput(t *testing.T) {
	t.Helper()
	helpers.AssertCLIHelpOutput(t, v.projectPath)
}

// ValidateVersionOutput validates CLI version output
func (v *CLIValidator) ValidateVersionOutput(t *testing.T) {
	t.Helper()
	helpers.AssertCLIVersionOutput(t, v.projectPath)
}

// ValidateLogger validates logger implementation
func (v *CLIValidator) ValidateLogger(t *testing.T, logger string) {
	t.Helper()

	loggerFile := filepath.Join(v.projectPath, "internal", "logger", logger+".go")
	helpers.AssertFileExists(t, loggerFile)

	// Check logger interface
	interfaceFile := filepath.Join(v.projectPath, "internal", "logger", "interface.go")
	helpers.AssertFileExists(t, interfaceFile)
	helpers.AssertFileContains(t, interfaceFile, "Logger interface")
}

// ValidateLoggerFunctionality validates logger functionality
func (v *CLIValidator) ValidateLoggerFunctionality(t *testing.T, logger string) {
	t.Helper()
	helpers.AssertLoggerFunctionality(t, v.projectPath, logger)
}

// ValidateViperConfiguration validates Viper configuration setup
func (v *CLIValidator) ValidateViperConfiguration(t *testing.T) {
	t.Helper()

	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	helpers.AssertFileExists(t, configFile)
	helpers.AssertFileContains(t, configFile, "github.com/spf13/viper")
}

// ValidateConfigFileSupport validates config file support
func (v *CLIValidator) ValidateConfigFileSupport(t *testing.T) {
	t.Helper()

	configYaml := filepath.Join(v.projectPath, "configs", "config.yaml")
	helpers.AssertFileExists(t, configYaml)
}

// ValidateDockerSupport validates Docker support
func (v *CLIValidator) ValidateDockerSupport(t *testing.T) {
	t.Helper()

	dockerfile := filepath.Join(v.projectPath, "Dockerfile")
	helpers.AssertFileExists(t, dockerfile)
	helpers.AssertFileContains(t, dockerfile, "FROM golang:")
}

// ValidateMakefileTargets validates Makefile targets
func (v *CLIValidator) ValidateMakefileTargets(t *testing.T) {
	t.Helper()

	makefile := filepath.Join(v.projectPath, "Makefile")
	helpers.AssertFileExists(t, makefile)
	helpers.AssertFileContains(t, makefile, "build")
	helpers.AssertFileContains(t, makefile, "docker")
}

// ValidateTestFiles validates test files exist
func (v *CLIValidator) ValidateTestFiles(t *testing.T) {
	t.Helper()

	rootTestFile := filepath.Join(v.projectPath, "cmd", "root_test.go")
	helpers.AssertFileExists(t, rootTestFile)

	configTestFile := filepath.Join(v.projectPath, "internal", "config", "config_test.go")
	helpers.AssertFileExists(t, configTestFile)
}

// ValidateTestifyUsage validates testify usage
func (v *CLIValidator) ValidateTestifyUsage(t *testing.T) {
	t.Helper()

	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)
	helpers.AssertFileContains(t, goModFile, "github.com/stretchr/testify")
}

// ValidateTestExecution validates test execution
func (v *CLIValidator) ValidateTestExecution(t *testing.T) {
	t.Helper()
	helpers.AssertTestsRun(t, v.projectPath)
}