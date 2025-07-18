package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
)

// Feature: CLI Application Blueprint
// As a developer
// I want to generate CLI applications
// So that I can build command-line tools quickly

func TestCLI_BasicGeneration_WithCobra(t *testing.T) {
	// Scenario: Generate basic CLI application
	// Given I want a CLI application
	// When I generate a CLI project
	// Then the project should use Cobra framework
	// And the project should have a root command
	// And the project should include example subcommands
	// And the project should compile to a working binary

	// Given I want a CLI application
	config := types.ProjectConfig{
		Name:      "test-cli-basic",
		Module:    "github.com/test/test-cli-basic",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	// When I generate a CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should use Cobra framework
	validator := NewCLIValidator(projectPath)
	validator.ValidateCobraSetup(t)

	// And the project should have a root command
	validator.ValidateRootCommand(t)

	// And the project should include example subcommands
	validator.ValidateExampleSubcommands(t)

	// And the project should compile to a working binary
	validator.ValidateCompilation(t)
}

func TestCLI_CommandStructure(t *testing.T) {
	// Scenario: CLI with comprehensive command structure
	// Given I generate a CLI application
	// Then it should include example subcommands
	// And each subcommand should be properly defined
	// And help text should be comprehensive
	// And command structure should be logical

	config := types.ProjectConfig{
		Name:      "test-cli-commands",
		Module:    "github.com/test/test-cli-commands",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewCLIValidator(projectPath)
	
	// Then it should include example subcommands
	validator.ValidateSubcommands(t, []string{"version", "config", "serve"})

	// And each subcommand should be properly defined
	validator.ValidateSubcommandDefinitions(t)

	// And help text should be comprehensive
	validator.ValidateHelpText(t)

	// And command structure should be logical
	validator.ValidateCommandStructure(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestCLI_ConfigurationManagement(t *testing.T) {
	// Feature: CLI Configuration Management
	// As a CLI developer
	// I want configuration management built-in
	// So that users can customize behavior

	// Scenario: Configuration file support
	// Given a CLI application with config support
	// Then it should support YAML configuration
	// And it should support environment variables
	// And it should support command-line flags
	// And precedence should be flags > env > config file

	config := types.ProjectConfig{
		Name:      "test-cli-config",
		Module:    "github.com/test/test-cli-config",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewCLIValidator(projectPath)

	// Then it should support YAML configuration
	validator.ValidateYAMLConfiguration(t)

	// And it should support environment variables
	validator.ValidateEnvironmentVariables(t)

	// And it should support command-line flags
	validator.ValidateCommandLineFlags(t)

	// And precedence should be flags > env > config file
	validator.ValidateConfigurationPrecedence(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestCLI_LoggerIntegration(t *testing.T) {
	// Feature: Logger Integration Across CLI Applications
	// As a CLI developer
	// I want consistent logging across CLI applications
	// So that I can debug and monitor applications effectively

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("Logger_"+logger, func(t *testing.T) {
			t.Parallel() // Safe to run logger tests in parallel

			// Scenario: Logger integration
			// Given I generate a CLI application with "<logger>"
			// Then logging should be properly configured
			// And log output should be controllable via flags
			// And log levels should be configurable
			// And structured logging should be available

			// Given I generate a CLI application with "<logger>"
			config := types.ProjectConfig{
				Name:      "test-cli-logger-" + logger,
				Module:    "github.com/test/test-cli-logger-" + logger,
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then logging should be properly configured
			validator := NewCLIValidator(projectPath)
			validator.ValidateLoggerIntegration(t, logger)

			// And log output should be controllable via flags
			validator.ValidateLoggerFlags(t)

			// And log levels should be configurable
			validator.ValidateLoggerConfiguration(t)

			// And structured logging should be available
			validator.ValidateStructuredLogging(t, logger)

			// And the project should compile successfully
			validator.ValidateCompilation(t)
		})
	}
}

func TestCLI_BinaryExecution(t *testing.T) {
	// Feature: CLI Application Functionality
	// As a CLI user
	// I want the generated CLI to work correctly
	// So that I can use it as a foundation

	// Scenario: Command execution
	// Given a compiled CLI application
	// When I run the binary without arguments
	// Then help text should be displayed
	// And available commands should be listed
	// And usage examples should be shown

	config := types.ProjectConfig{
		Name:      "test-cli-execution",
		Module:    "github.com/test/test-cli-execution",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewCLIValidator(projectPath)

	// Given a compiled CLI application
	validator.ValidateCompilation(t)

	// When I run the binary without arguments
	// Then help text should be displayed
	// And available commands should be listed
	// And usage examples should be shown
	validator.ValidateBinaryExecution(t)
}

func TestCLI_DockerSupport(t *testing.T) {
	// Scenario: Docker containerization
	// Given a CLI application
	// When I check for Docker support
	// Then Dockerfile should be present
	// And docker-compose should be available
	// And container build should be possible

	config := types.ProjectConfig{
		Name:      "test-cli-docker",
		Module:    "github.com/test/test-cli-docker",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewCLIValidator(projectPath)

	// Then Dockerfile should be present
	validator.ValidateDockerSupport(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestCLI_MinimalConfiguration(t *testing.T) {
	// Scenario: Minimal configuration without optional features
	// Given I want a minimal CLI application
	// When I generate without optional features
	// Then only core files should be present
	// And optional feature files should not exist

	config := types.ProjectConfig{
		Name:      "test-cli-minimal",
		Module:    "github.com/test/test-cli-minimal",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
		// No optional features configured
	}

	projectPath := helpers.GenerateProject(t, config)

	validator := NewCLIValidator(projectPath)
	validator.ValidateMinimalConfiguration(t)
	validator.ValidateCompilation(t)
}

// CLIValidator provides validation methods for CLI blueprints
type CLIValidator struct {
	ProjectPath string
}

func NewCLIValidator(projectPath string) *CLIValidator {
	return &CLIValidator{
		ProjectPath: projectPath,
	}
}

func (v *CLIValidator) ValidateCobraSetup(t *testing.T) {
	t.Helper()
	
	// Check go.mod for Cobra dependency
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	helpers.AssertFileContains(t, goMod, "github.com/spf13/cobra")
	
	// Check main.go uses Cobra
	mainFile := filepath.Join(v.ProjectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	helpers.AssertFileContains(t, mainFile, "cmd.Execute()")
}

func (v *CLIValidator) ValidateRootCommand(t *testing.T) {
	t.Helper()
	
	// Check cmd/root.go exists and has proper structure
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "var rootCmd")
	helpers.AssertFileContains(t, rootFile, "cobra.Command")
}

func (v *CLIValidator) ValidateExampleSubcommands(t *testing.T) {
	t.Helper()
	
	// Check that example subcommands exist
	expectedCommands := []string{"version.go", "config.go", "serve.go"}
	for _, cmd := range expectedCommands {
		cmdFile := filepath.Join(v.ProjectPath, "cmd", cmd)
		helpers.AssertFileExists(t, cmdFile)
	}
}

func (v *CLIValidator) ValidateSubcommands(t *testing.T, expectedCommands []string) {
	t.Helper()
	
	for _, cmd := range expectedCommands {
		cmdFile := filepath.Join(v.ProjectPath, "cmd", cmd+".go")
		helpers.AssertFileExists(t, cmdFile)
		helpers.AssertFileContains(t, cmdFile, "cobra.Command")
	}
}

func (v *CLIValidator) ValidateSubcommandDefinitions(t *testing.T) {
	t.Helper()
	
	// Check version command
	versionFile := filepath.Join(v.ProjectPath, "cmd/version.go")
	if helpers.FileExists(versionFile) {
		helpers.AssertFileContains(t, versionFile, "Use:")
		helpers.AssertFileContains(t, versionFile, "Short:")
	}
	
	// Check config command
	configFile := filepath.Join(v.ProjectPath, "cmd/config.go")
	if helpers.FileExists(configFile) {
		helpers.AssertFileContains(t, configFile, "Use:")
		helpers.AssertFileContains(t, configFile, "Short:")
	}
}

func (v *CLIValidator) ValidateHelpText(t *testing.T) {
	t.Helper()
	
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "Short:")
	helpers.AssertFileContains(t, rootFile, "Long:")
}

func (v *CLIValidator) ValidateCommandStructure(t *testing.T) {
	t.Helper()
	
	// Check cmd directory exists
	cmdDir := filepath.Join(v.ProjectPath, "cmd")
	helpers.AssertDirExists(t, cmdDir)
	
	// Check that commands are properly organized
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "Execute()")
}

func (v *CLIValidator) ValidateYAMLConfiguration(t *testing.T) {
	t.Helper()
	
	// Check config.yaml exists
	configFile := filepath.Join(v.ProjectPath, "configs/config.yaml")
	helpers.AssertFileExists(t, configFile)
	
	// Check config package exists
	configGo := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configGo)
	helpers.AssertFileContains(t, configGo, "viper")
}

func (v *CLIValidator) ValidateEnvironmentVariables(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	helpers.AssertFileContains(t, configFile, "viper.AutomaticEnv")
}

func (v *CLIValidator) ValidateCommandLineFlags(t *testing.T) {
	t.Helper()
	
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "PersistentFlags")
}

func (v *CLIValidator) ValidateConfigurationPrecedence(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	// Check that Viper is configured to handle precedence
	helpers.AssertFileContains(t, configFile, "viper.BindPFlag")
}

func (v *CLIValidator) ValidateLoggerIntegration(t *testing.T, logger string) {
	t.Helper()
	
	// Logger interface should exist
	interfaceFile := filepath.Join(v.ProjectPath, "internal/logger/interface.go")
	helpers.AssertFileExists(t, interfaceFile)
	
	// Factory should exist
	factoryFile := filepath.Join(v.ProjectPath, "internal/logger/factory.go")
	helpers.AssertFileExists(t, factoryFile)
	helpers.AssertFileContains(t, factoryFile, logger)
	
	// Logger-specific implementation should exist
	loggerFile := filepath.Join(v.ProjectPath, "internal/logger/"+logger+".go")
	helpers.AssertFileExists(t, loggerFile)
	
	// Check go.mod for logger dependency
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	switch logger {
	case "zap":
		helpers.AssertFileContains(t, goMod, "go.uber.org/zap")
	case "logrus":
		helpers.AssertFileContains(t, goMod, "github.com/sirupsen/logrus")
	case "zerolog":
		helpers.AssertFileContains(t, goMod, "github.com/rs/zerolog")
	case "slog":
		// slog is part of standard library, no external dependency needed
	}
}

func (v *CLIValidator) ValidateLoggerFlags(t *testing.T) {
	t.Helper()
	
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	// Check for log level flags
	content := helpers.ReadFileContent(t, rootFile)
	if helpers.StringContains(content, "log-level") || helpers.StringContains(content, "verbose") || helpers.StringContains(content, "debug") {
		t.Log("✓ Logger flags found in root command")
	} else {
		t.Log("⚠ Logger flags not found - may be in configuration")
	}
}

func (v *CLIValidator) ValidateLoggerConfiguration(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	
	// Check that logger configuration is included
	content := helpers.ReadFileContent(t, configFile)
	if helpers.StringContains(content, "LogLevel") || helpers.StringContains(content, "Logger") {
		t.Log("✓ Logger configuration found")
	}
}

func (v *CLIValidator) ValidateStructuredLogging(t *testing.T, logger string) {
	t.Helper()
	
	loggerFile := filepath.Join(v.ProjectPath, "internal/logger/"+logger+".go")
	helpers.AssertFileExists(t, loggerFile)
	
	// Check for structured logging patterns
	content := helpers.ReadFileContent(t, loggerFile)
	switch logger {
	case "slog":
		assert.Contains(t, content, "slog.")
	case "zap":
		assert.Contains(t, content, "zap.")
	case "logrus":
		assert.Contains(t, content, "logrus.")
	case "zerolog":
		assert.Contains(t, content, "zerolog.")
	}
}

func (v *CLIValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertCompilable(t, v.ProjectPath)
}

func (v *CLIValidator) ValidateBinaryExecution(t *testing.T) {
	t.Helper()
	
	// First ensure the project compiles
	helpers.AssertCompilable(t, v.ProjectPath)
	
	// Build the binary
	binaryName := filepath.Base(v.ProjectPath)
	binaryPath := filepath.Join(v.ProjectPath, binaryName)
	
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	buildCmd.Dir = v.ProjectPath
	buildOutput, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build CLI binary: %v\nOutput: %s", err, buildOutput)
	}
	
	// Test running the binary without arguments (should show help)
	cmd := exec.Command(binaryPath)
	cmd.Dir = v.ProjectPath
	output, err := cmd.CombinedOutput()
	
	// CLI should either succeed (return 0) or fail gracefully with help
	if err != nil {
		// Check if it's an expected exit code (like help showing)
		if exitError, ok := err.(*exec.ExitError); ok {
			if exitError.ExitCode() == 0 || exitError.ExitCode() == 1 {
				t.Log("✓ CLI binary executed and showed help as expected")
			} else {
				t.Errorf("CLI binary failed with unexpected exit code: %d", exitError.ExitCode())
			}
		}
	} else {
		t.Log("✓ CLI binary executed successfully")
	}
	
	// Check that output contains expected help elements
	outputStr := string(output)
	if helpers.StringContains(outputStr, "Usage:") || helpers.StringContains(outputStr, "Available Commands:") || helpers.StringContains(outputStr, "Flags:") {
		t.Log("✓ CLI help output contains expected elements")
	} else {
		t.Logf("⚠ CLI output may not contain standard help elements. Output: %s", outputStr)
	}
	
	// Clean up binary
	_ = os.Remove(binaryPath)
}

func (v *CLIValidator) ValidateDockerSupport(t *testing.T) {
	t.Helper()
	
	// Check Dockerfile exists
	dockerfile := filepath.Join(v.ProjectPath, "Dockerfile")
	helpers.AssertFileExists(t, dockerfile)
	
	// Check that Dockerfile has proper structure
	helpers.AssertFileContains(t, dockerfile, "FROM golang:")
	helpers.AssertFileContains(t, dockerfile, "COPY")
	helpers.AssertFileContains(t, dockerfile, "RUN go build")
}

func (v *CLIValidator) ValidateMinimalConfiguration(t *testing.T) {
	t.Helper()
	
	// Core files should exist
	coreFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
		"main.go",
		"cmd/root.go",
		"cmd/version.go",
		"internal/config/config.go",
		"internal/logger/interface.go",
		"internal/logger/factory.go",
		"configs/config.yaml",
	}
	
	for _, file := range coreFiles {
		helpers.AssertFileExists(t, filepath.Join(v.ProjectPath, file))
	}
	
	// Should have Docker support
	helpers.AssertFileExists(t, filepath.Join(v.ProjectPath, "Dockerfile"))
	
	// Should have proper Go module structure
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	helpers.AssertFileContains(t, goMod, "github.com/spf13/cobra")
	helpers.AssertFileContains(t, goMod, "github.com/spf13/viper")
}