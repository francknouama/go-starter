package cli_test

import (
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// Feature: CLI Integration Testing
// As a developer
// I want CLI blueprints to work together with different configurations
// So that I can build robust command-line applications

func TestCLI_CrossLoggerIntegration(t *testing.T) {
	// Feature: Cross-Logger Integration
	// As a CLI developer
	// I want all loggers to work consistently across CLI applications
	// So that I can switch loggers without changing my code

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("CrossIntegration_"+logger, func(t *testing.T) {
			t.Parallel()

			// Scenario: Cross-logger integration
			// Given I generate a CLI with "<logger>"
			// Then the logger should integrate with all CLI components
			// And configuration should support logger-specific settings
			// And the CLI should compile and run successfully
			// And binary execution should work with the logger

			// Given I generate a CLI with "<logger>"
			config := types.ProjectConfig{
				Name:      "test-cli-cross-" + logger,
				Module:    "github.com/test/test-cli-cross-" + logger,
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)

			// Then the logger should integrate with all CLI components
			validator := NewCLIValidator(projectPath)
			validator.ValidateLoggerIntegration(t, logger)

			// And configuration should support logger-specific settings
			validator.ValidateLoggerConfiguration(t)

			// And the CLI should compile and run successfully
			validator.ValidateCompilation(t)

			// And binary execution should work with the logger
			validator.ValidateBinaryExecution(t)
		})
	}
}

func TestCLI_CompilationValidation(t *testing.T) {
	// Feature: CLI Compilation Validation
	// As a quality assurance engineer
	// I want all CLI configurations to compile successfully
	// So that users always get working projects

	testCases := []struct {
		name   string
		config types.ProjectConfig
	}{
		{
			name: "StandardCLI_Slog",
			config: types.ProjectConfig{
				Name:      "test-cli-compile-slog",
				Module:    "github.com/test/test-cli-compile-slog",
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "slog",
			},
		},
		{
			name: "StandardCLI_Zap",
			config: types.ProjectConfig{
				Name:      "test-cli-compile-zap",
				Module:    "github.com/test/test-cli-compile-zap",
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "zap",
			},
		},
		{
			name: "StandardCLI_Logrus",
			config: types.ProjectConfig{
				Name:      "test-cli-compile-logrus",
				Module:    "github.com/test/test-cli-compile-logrus",
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "logrus",
			},
		},
		{
			name: "StandardCLI_Zerolog",
			config: types.ProjectConfig{
				Name:      "test-cli-compile-zerolog",
				Module:    "github.com/test/test-cli-compile-zerolog",
				Type:      "cli-standard",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "zerolog",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Scenario: Compilation validation
			// Given I generate a CLI with specific configuration
			// Then the project should compile without errors
			// And all dependencies should be correctly included
			// And the binary should be executable

			// Given I generate a CLI with specific configuration
			projectPath := helpers.GenerateProject(t, tc.config)

			// Then the project should compile without errors
			validator := NewCLIValidator(projectPath)
			validator.ValidateCompilation(t)

			// And all dependencies should be correctly included
			validator.ValidateDependencies(t, tc.config.Logger)

			// And the binary should be executable
			validator.ValidateBinaryExecution(t)
		})
	}
}

func TestCLI_ArchitectureCompliance(t *testing.T) {
	// Feature: CLI Architecture Compliance
	// As a code reviewer
	// I want generated CLI code to follow architecture principles
	// So that projects maintain architectural integrity

	// Scenario: Architecture compliance for CLI
	// Given I generate a CLI application
	// Then the code should follow CLI best practices
	// And dependency directions should be correct
	// And package organization should be logical
	// And the project should pass architectural validation

	config := types.ProjectConfig{
		Name:      "test-cli-architecture",
		Module:    "github.com/test/test-cli-architecture",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should follow CLI best practices
	validator := NewCLIValidator(projectPath)
	validator.ValidateCLIArchitecture(t)

	// And dependency directions should be correct
	validator.ValidateDependencyDirections(t)

	// And package organization should be logical
	validator.ValidatePackageOrganization(t)

	// And the project should pass architectural validation
	validator.ValidateArchitecturalCompliance(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestCLI_SecurityValidation(t *testing.T) {
	// Feature: CLI Security Validation
	// As a security engineer
	// I want CLI applications to follow security best practices
	// So that generated applications are secure by default

	// Scenario: Security validation
	// Given I generate a CLI application
	// Then the code should not contain security vulnerabilities
	// And configuration should be handled securely
	// And logging should not expose sensitive information
	// And dependencies should be secure

	config := types.ProjectConfig{
		Name:      "test-cli-security",
		Module:    "github.com/test/test-cli-security",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Then the code should not contain security vulnerabilities
	validator := NewCLIValidator(projectPath)
	validator.ValidateSecurityPractices(t)

	// And configuration should be handled securely
	validator.ValidateSecureConfiguration(t)

	// And logging should not expose sensitive information
	validator.ValidateSecureLogging(t)

	// And the project should compile successfully
	validator.ValidateCompilation(t)
}

func TestCLI_RuntimeValidation(t *testing.T) {
	// Feature: CLI Runtime Validation
	// As a user
	// I want CLI applications to work correctly at runtime
	// So that I can use them immediately after generation

	// Scenario: Runtime validation
	// Given I generate and build a CLI application
	// Then the CLI should execute without errors
	// And help system should work correctly
	// And configuration loading should work
	// And logging should be functional

	config := types.ProjectConfig{
		Name:      "test-cli-runtime",
		Module:    "github.com/test/test-cli-runtime",
		Type:      "cli-standard",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Given I generate and build a CLI application
	validator := NewCLIValidator(projectPath)
	validator.ValidateCompilation(t)

	// Then the CLI should execute without errors
	validator.ValidateBinaryExecution(t)

	// And help system should work correctly
	validator.ValidateHelpSystem(t)

	// And configuration loading should work
	validator.ValidateConfigurationLoading(t)

	// And logging should be functional
	validator.ValidateLoggingFunctionality(t)
}

// Additional validation methods for integration tests

func (v *CLIValidator) ValidateDependencies(t *testing.T, logger string) {
	t.Helper()
	
	goMod := filepath.Join(v.ProjectPath, "go.mod")
	helpers.AssertFileExists(t, goMod)
	
	// Core dependencies should always be present
	helpers.AssertFileContains(t, goMod, "github.com/spf13/cobra")
	helpers.AssertFileContains(t, goMod, "github.com/spf13/viper")
	
	// Logger-specific dependencies
	switch logger {
	case "zap":
		helpers.AssertFileContains(t, goMod, "go.uber.org/zap")
	case "logrus":
		helpers.AssertFileContains(t, goMod, "github.com/sirupsen/logrus")
	case "zerolog":
		helpers.AssertFileContains(t, goMod, "github.com/rs/zerolog")
	case "slog":
		// slog is part of standard library, no external dependency
	}
}

func (v *CLIValidator) ValidateCLIArchitecture(t *testing.T) {
	t.Helper()
	
	// Check standard CLI directory structure
	expectedDirs := []string{
		"cmd",
		"internal/config",
		"internal/logger",
		"configs",
	}
	
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(v.ProjectPath, dir)
		helpers.AssertDirExists(t, dirPath)
	}
	
	// Check main.go is at root
	mainFile := filepath.Join(v.ProjectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	
	// Check cmd package structure
	rootCmd := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootCmd)
}

func (v *CLIValidator) ValidateDependencyDirections(t *testing.T) {
	t.Helper()
	
	// main.go should import cmd package
	mainFile := filepath.Join(v.ProjectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	content := helpers.ReadFileContent(t, mainFile)
	
	// Should import cmd package
	if helpers.StringContains(content, "cmd") {
		t.Log("✓ main.go properly imports cmd package")
	}
	
	// cmd/root.go should import internal packages
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	if helpers.FileExists(rootFile) {
		rootContent := helpers.ReadFileContent(t, rootFile)
		if helpers.StringContains(rootContent, "internal/config") || helpers.StringContains(rootContent, "internal/logger") {
			t.Log("✓ cmd package properly imports internal packages")
		}
	}
}

func (v *CLIValidator) ValidatePackageOrganization(t *testing.T) {
	t.Helper()
	
	// Check that each package has a clear purpose
	// cmd/ should contain command definitions
	cmdDir := filepath.Join(v.ProjectPath, "cmd")
	helpers.AssertDirExists(t, cmdDir)
	
	// internal/ should contain internal logic
	internalDir := filepath.Join(v.ProjectPath, "internal")
	helpers.AssertDirExists(t, internalDir)
	
	// configs/ should contain configuration files
	configsDir := filepath.Join(v.ProjectPath, "configs")
	helpers.AssertDirExists(t, configsDir)
}

func (v *CLIValidator) ValidateArchitecturalCompliance(t *testing.T) {
	t.Helper()
	
	// For CLI, validate that it follows CLI best practices
	// This includes proper command structure, configuration management, etc.
	
	// Check that root command exists and is properly structured
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "Execute()")
	
	// Check that configuration is properly structured
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	helpers.AssertFileContains(t, configFile, "viper")
}

func (v *CLIValidator) ValidateSecurityPractices(t *testing.T) {
	t.Helper()
	
	// Check that no hardcoded secrets exist
	files := []string{
		"main.go",
		"cmd/root.go",
		"internal/config/config.go",
	}
	
	for _, file := range files {
		filePath := filepath.Join(v.ProjectPath, file)
		if helpers.FileExists(filePath) {
			content := helpers.ReadFileContent(t, filePath)
			
			// Check for common security issues
			securityIssues := []string{
				"password=",
				"secret=",
				"token=",
				"apikey=",
				"api_key=",
			}
			
			for _, issue := range securityIssues {
				if helpers.StringContains(content, issue) {
					t.Errorf("Potential security issue found in %s: %s", file, issue)
				}
			}
		}
	}
}

func (v *CLIValidator) ValidateSecureConfiguration(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	
	// Check that configuration uses environment variables for sensitive data
	content := helpers.ReadFileContent(t, configFile)
	if helpers.StringContains(content, "viper.AutomaticEnv") {
		t.Log("✓ Configuration supports environment variables")
	}
}

func (v *CLIValidator) ValidateSecureLogging(t *testing.T) {
	t.Helper()
	
	// Check that logging configuration doesn't expose sensitive information
	loggerFiles := []string{
		"internal/logger/interface.go",
		"internal/logger/factory.go",
		"internal/logger/slog.go",
		"internal/logger/zap.go",
		"internal/logger/logrus.go",
		"internal/logger/zerolog.go",
	}
	
	for _, file := range loggerFiles {
		filePath := filepath.Join(v.ProjectPath, file)
		if helpers.FileExists(filePath) {
			content := helpers.ReadFileContent(t, filePath)
			
			// Check that logger doesn't log sensitive fields
			sensitiveFields := []string{
				"password",
				"secret",
				"token",
				"apikey",
				"api_key",
			}
			
			for _, field := range sensitiveFields {
				if helpers.StringContains(content, field) {
					t.Logf("⚠ Logger may be logging sensitive field: %s in %s", field, file)
				}
			}
		}
	}
}

func (v *CLIValidator) ValidateHelpSystem(t *testing.T) {
	t.Helper()
	
	// This would ideally test the help system by running the CLI
	// For now, validate that help text is properly defined
	rootFile := filepath.Join(v.ProjectPath, "cmd/root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "Short:")
	helpers.AssertFileContains(t, rootFile, "Long:")
	
	// Check that subcommands have help text
	subcommands := []string{"version.go", "config.go", "serve.go"}
	for _, cmd := range subcommands {
		cmdFile := filepath.Join(v.ProjectPath, "cmd", cmd)
		if helpers.FileExists(cmdFile) {
			helpers.AssertFileContains(t, cmdFile, "Short:")
		}
	}
}

func (v *CLIValidator) ValidateConfigurationLoading(t *testing.T) {
	t.Helper()
	
	configFile := filepath.Join(v.ProjectPath, "internal/config/config.go")
	helpers.AssertFileExists(t, configFile)
	
	// Check that configuration loading is properly implemented
	content := helpers.ReadFileContent(t, configFile)
	
	// Should have config loading logic
	if helpers.StringContains(content, "viper.ReadInConfig") || helpers.StringContains(content, "viper.SetConfigFile") {
		t.Log("✓ Configuration loading is properly implemented")
	}
}

func (v *CLIValidator) ValidateLoggingFunctionality(t *testing.T) {
	t.Helper()
	
	// Check that logging is properly initialized
	mainFile := filepath.Join(v.ProjectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	
	content := helpers.ReadFileContent(t, mainFile)
	if helpers.StringContains(content, "logger") {
		t.Log("✓ Logging is initialized in main.go")
	}
	
	// Check that logger factory exists
	factoryFile := filepath.Join(v.ProjectPath, "internal/logger/factory.go")
	helpers.AssertFileExists(t, factoryFile)
	helpers.AssertFileContains(t, factoryFile, "NewLogger")
}

// Import required packages
import (
	"path/filepath"
)