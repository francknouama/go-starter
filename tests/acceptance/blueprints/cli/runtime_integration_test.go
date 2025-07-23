package cli_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Initialize templates filesystem for ATDD tests
	wd, err := os.Getwd()
	if err != nil {
		panic("Failed to get working directory: " + err.Error())
	}

	// Navigate to project root and find blueprints directory
	projectRoot := wd
	for {
		templatesDir := filepath.Join(projectRoot, "blueprints")
		if _, err := os.Stat(templatesDir); err == nil {
			entries, err := os.ReadDir(templatesDir)
			if err == nil && len(entries) > 0 {
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

		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			panic("Could not find blueprints directory")
		}
		projectRoot = parentDir
	}
}

// Feature: CLI Application Runtime Integration Testing
// As a developer
// I want to validate that generated CLI applications work correctly at runtime
// So that I can be confident the blueprint generates functional command-line tools

func TestCLI_RuntimeIntegration_ComprehensiveSubcommandTesting(t *testing.T) {
	// Scenario: Generate CLI and test all CRUD subcommands
	// Given I generate a CLI application with all subcommands
	// When I test each subcommand with various parameters
	// Then all subcommands should work correctly
	// And the CLI should handle errors gracefully
	// And output formatting should work as expected

	config := types.ProjectConfig{
		Name:      "test-cli-runtime",
		Module:    "github.com/test/test-cli-runtime",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)
	
	// Build the CLI binary
	validator := NewCLIRuntimeValidator(projectPath)
	validator.BuildCLI(t)

	// Test all subcommands comprehensively
	t.Run("create_subcommand", func(t *testing.T) {
		validator.TestCreateSubcommand(t)
	})

	t.Run("list_subcommand", func(t *testing.T) {
		validator.TestListSubcommand(t)
	})

	t.Run("delete_subcommand", func(t *testing.T) {
		validator.TestDeleteSubcommand(t)
	})

	t.Run("update_subcommand", func(t *testing.T) {
		validator.TestUpdateSubcommand(t)
	})

	t.Run("error_handling", func(t *testing.T) {
		validator.TestErrorHandling(t)
	})

	t.Run("output_formatting", func(t *testing.T) {
		validator.TestOutputFormatting(t)
	})
}

func TestCLI_MultiLoggerRuntimeValidation(t *testing.T) {
	// Scenario: Test CLI with different logger types
	// Given I generate CLI projects with different loggers
	// When I run the CLI applications
	// Then all logger types should work correctly
	// And logging output should be formatted properly

	loggers := []string{"slog", "zap", "logrus", "zerolog"}

	for _, logger := range loggers {
		t.Run("logger_"+logger+"_runtime", func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "test-cli-" + logger + "-runtime",
				Module:    "github.com/test/test-cli-" + logger + "-runtime", 
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    logger,
			}

			projectPath := helpers.GenerateProject(t, config)
			validator := NewCLIRuntimeValidator(projectPath)
			
			// Build and test basic functionality
			validator.BuildCLI(t)
			validator.TestBasicCommands(t)
			validator.TestLoggerOutput(t, logger)
		})
	}
}

func TestCLI_ConfigurationAndFlags(t *testing.T) {
	// Scenario: Test CLI configuration and flag handling
	// Given I generate a CLI application
	// When I use various flags and configuration options
	// Then the CLI should handle all combinations correctly
	// And configuration should override defaults properly

	config := types.ProjectConfig{
		Name:      "test-cli-config-runtime",
		Module:    "github.com/test/test-cli-config-runtime",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)
	validator := NewCLIRuntimeValidator(projectPath)
	validator.BuildCLI(t)

	t.Run("global_flags", func(t *testing.T) {
		validator.TestGlobalFlags(t)
	})

	t.Run("configuration_files", func(t *testing.T) {
		validator.TestConfigurationFiles(t)
	})

	t.Run("environment_variables", func(t *testing.T) {
		validator.TestEnvironmentVariables(t)
	})
}

// CLIRuntimeValidator provides comprehensive runtime testing for CLI applications
type CLIRuntimeValidator struct {
	projectPath string
	cliPath     string
}

// NewCLIRuntimeValidator creates a new CLI runtime validator
func NewCLIRuntimeValidator(projectPath string) *CLIRuntimeValidator {
	return &CLIRuntimeValidator{
		projectPath: projectPath,
		cliPath:     filepath.Join(projectPath, "test-cli"),
	}
}

// BuildCLI builds the CLI binary for testing
func (v *CLIRuntimeValidator) BuildCLI(t *testing.T) {
	t.Helper()
	
	cmd := exec.Command("go", "build", "-o", "test-cli", ".")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI should build successfully. Output: %s", string(output))
	
	// Verify binary exists
	require.FileExists(t, v.cliPath, "CLI binary should exist after build")
}

// TestCreateSubcommand tests the create subcommand comprehensively
func (v *CLIRuntimeValidator) TestCreateSubcommand(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name       string
		args       []string
		expectFail bool
		contains   []string
	}{
		{
			name: "create_project_basic",
			args: []string{"create", "project", "my-test-project"},
			contains: []string{"✅ Created project: my-test-project"},
		},
		{
			name: "create_project_with_template",
			args: []string{"create", "project", "my-project", "--template", "basic"},
			contains: []string{"✅ Created project: my-project", "Template: basic"},
		},
		{
			name: "create_config_default",
			args: []string{"create", "config"},
			contains: []string{"✅ Created config file"},
		},
		{
			name: "create_task_with_priority",
			args: []string{"create", "task", "Important task", "--priority", "high"},
			contains: []string{"✅ Created task: Important task", "Priority: high"},
		},
		{
			name: "create_dry_run",
			args: []string{"create", "project", "dry-run-project", "--dry-run"},
			contains: []string{"Would create project: dry-run-project"},
		},
		{
			name: "create_invalid_resource",
			args: []string{"create", "invalid", "name"},
			expectFail: true,
			contains: []string{"invalid resource type"},
		},
		{
			name: "create_missing_args",
			args: []string{"create"},
			expectFail: true,
			contains: []string{"requires at least 1 argument"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			
			if tc.expectFail {
				assert.Error(t, err, "Command should fail")
			} else {
				assert.NoError(t, err, "Command should succeed. Output: %s", output)
			}
			
			for _, contains := range tc.contains {
				assert.Contains(t, output, contains, "Output should contain expected text")
			}
		})
	}
}

// TestListSubcommand tests the list subcommand comprehensively
func (v *CLIRuntimeValidator) TestListSubcommand(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "list_projects_table",
			args: []string{"list", "projects"},
			contains: []string{"NAME", "PATH", "GIT"},
		},
		{
			name: "list_projects_json",
			args: []string{"list", "projects", "--format", "json"},
			contains: []string{"[", "]"}, // JSON array markers
		},
		{
			name: "list_configs_verbose",
			args: []string{"list", "configs", "--verbose"},
			contains: []string{"NAME", "TYPE", "PATH", "SIZE", "MODIFIED"},
		},
		{
			name: "list_tasks_by_priority",
			args: []string{"list", "tasks", "--sort", "priority"},
			contains: []string{"NAME", "PRIORITY", "STATUS"},
		},
		{
			name: "list_tasks_all",
			args: []string{"list", "tasks", "--all"},
			contains: []string{"completed"}, // Should show completed tasks
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			assert.NoError(t, err, "List command should succeed. Output: %s", output)
			
			for _, contains := range tc.contains {
				assert.Contains(t, output, contains, "Output should contain expected text")
			}
		})
	}
}

// TestDeleteSubcommand tests the delete subcommand comprehensively
func (v *CLIRuntimeValidator) TestDeleteSubcommand(t *testing.T) {
	t.Helper()

	// Create test files first
	testProject := filepath.Join(v.projectPath, "test-project")
	err := os.Mkdir(testProject, 0755)
	require.NoError(t, err, "Should create test project directory")

	testConfig := filepath.Join(v.projectPath, "test-config.yaml")
	err = os.WriteFile(testConfig, []byte("test: config"), 0644)
	require.NoError(t, err, "Should create test config file")

	testCases := []struct {
		name       string
		args       []string
		expectFail bool
		contains   []string
	}{
		{
			name: "delete_dry_run",
			args: []string{"delete", "project", "test-project", "--dry-run"},
			contains: []string{"Would delete project: test-project"},
		},
		{
			name: "delete_config_force",
			args: []string{"delete", "config", "test-config.yaml", "--force"},
			contains: []string{"✅ Deleted config file: test-config.yaml"},
		},
		{
			name: "delete_task_force",
			args: []string{"delete", "task", "task-001", "--force"},
			contains: []string{"✅ Deleted task: task-001"},
		},
		{
			name: "delete_nonexistent",
			args: []string{"delete", "project", "nonexistent", "--force"},
			expectFail: true,
			contains: []string{"does not exist"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			
			if tc.expectFail {
				assert.Error(t, err, "Command should fail")
			} else {
				assert.NoError(t, err, "Command should succeed. Output: %s", output)
			}
			
			for _, contains := range tc.contains {
				assert.Contains(t, output, contains, "Output should contain expected text")
			}
		})
	}
}

// TestUpdateSubcommand tests the update subcommand comprehensively
func (v *CLIRuntimeValidator) TestUpdateSubcommand(t *testing.T) {
	t.Helper()

	testCases := []struct {
		name       string
		args       []string
		expectFail bool
		contains   []string
	}{
		{
			name: "update_project_description",
			args: []string{"update", "project", "test-project", "--description", "Updated description", "--force"},
			contains: []string{"✅ Updated project: test-project", "Description: Updated description"},
		},
		{
			name: "update_config_logging",
			args: []string{"update", "config", "--logging-level", "debug", "--force"},
			contains: []string{"✅ Updated config file", "logging.level: debug"},
		},
		{
			name: "update_task_priority",
			args: []string{"update", "task", "task-001", "--priority", "urgent", "--status", "completed", "--force"},
			contains: []string{"✅ Updated task: task-001", "Priority: urgent", "Status: completed"},
		},
		{
			name: "update_dry_run",
			args: []string{"update", "task", "task-001", "--priority", "low", "--dry-run"},
			contains: []string{"Would update task: task-001", "Priority: low"},
		},
		{
			name: "update_invalid_priority",
			args: []string{"update", "task", "task-001", "--priority", "invalid", "--force"},
			expectFail: true,
			contains: []string{"Invalid priority"},
		},
		{
			name: "update_no_fields",
			args: []string{"update", "task", "task-001"},
			expectFail: true,
			contains: []string{"At least one"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			
			if tc.expectFail {
				assert.Error(t, err, "Command should fail")
			} else {
				assert.NoError(t, err, "Command should succeed. Output: %s", output)
			}
			
			for _, contains := range tc.contains {
				assert.Contains(t, output, contains, "Output should contain expected text")
			}
		})
	}
}

// TestErrorHandling tests error handling scenarios
func (v *CLIRuntimeValidator) TestErrorHandling(t *testing.T) {
	t.Helper()

	errorCases := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "invalid_command",
			args: []string{"invalid-command"},
			contains: []string{"Error:", "invalid-command"},
		},
		{
			name: "invalid_format",
			args: []string{"list", "projects", "--format", "invalid"},
			contains: []string{"Error:", "Unknown format"},
		},
		{
			name: "invalid_sort",
			args: []string{"list", "tasks", "--sort", "invalid"},
			contains: []string{}, // Should use default sort, no error
		},
	}

	for _, tc := range errorCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			
			// Most error cases should exit with non-zero code
			if len(tc.contains) > 0 {
				assert.Error(t, err, "Command should fail")
				for _, contains := range tc.contains {
					assert.Contains(t, output, contains, "Output should contain expected error text")
				}
			}
		})
	}
}

// TestOutputFormatting tests output formatting options
func (v *CLIRuntimeValidator) TestOutputFormatting(t *testing.T) {
	t.Helper()

	formatCases := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name: "table_format",
			args: []string{"list", "projects", "--format", "table"},
			contains: []string{"NAME", "PATH"}, // Table headers
		},
		{
			name: "json_format",
			args: []string{"list", "tasks", "--format", "json"},
			contains: []string{"[", "]", "name", "priority"}, // JSON structure
		},
		{
			name: "verbose_output",
			args: []string{"list", "configs", "--verbose"},
			contains: []string{"SIZE", "MODIFIED"}, // Extra verbose columns
		},
	}

	for _, tc := range formatCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := v.runCLI(tc.args...)
			assert.NoError(t, err, "Command should succeed. Output: %s", output)
			
			for _, contains := range tc.contains {
				assert.Contains(t, output, contains, "Output should contain expected formatting")
			}
		})
	}
}

// TestBasicCommands tests basic CLI functionality
func (v *CLIRuntimeValidator) TestBasicCommands(t *testing.T) {
	t.Helper()

	// Test help command
	helpOutput, err := v.runCLI("--help")
	assert.NoError(t, err, "Help command should succeed")
	assert.Contains(t, helpOutput, "Usage:")
	assert.Contains(t, helpOutput, "Resource Management:")

	// Test version command
	versionOutput, err := v.runCLI("version")
	assert.NoError(t, err, "Version command should succeed")
	assert.Contains(t, versionOutput, "Version:")
}

// TestLoggerOutput tests logger-specific output
func (v *CLIRuntimeValidator) TestLoggerOutput(t *testing.T, logger string) {
	t.Helper()

	// Run a command that should generate log output
	output, err := v.runCLI("create", "task", "test-logging", "--force")
	assert.NoError(t, err, "Create command should succeed")
	
	// The exact log format depends on the logger, but we should see structured output
	// For slog: level=INFO msg="Creating resource"
	// For zap: {"level":"info","msg":"Creating resource"}
	// For logrus: level=info msg="Creating resource"
	// For zerolog: {"level":"info","message":"Creating resource"}
	
	switch logger {
	case "slog":
		// slog typically uses key=value format
		assert.Contains(t, output, "Creating resource", "Should contain log message")
	case "zap", "zerolog":
		// These often use JSON format
		assert.Contains(t, output, "Creating resource", "Should contain log message")
	case "logrus":
		// logrus typically uses key=value format similar to slog
		assert.Contains(t, output, "Creating resource", "Should contain log message")
	}
}

// TestGlobalFlags tests global flag functionality
func (v *CLIRuntimeValidator) TestGlobalFlags(t *testing.T) {
	t.Helper()

	// Test verbose flag
	output, err := v.runCLI("--verbose", "create", "task", "verbose-test", "--force")
	assert.NoError(t, err, "Verbose command should succeed")
	// Verbose mode should include more detailed output
	assert.Contains(t, output, "Creating resource", "Verbose output should contain detailed logs")
}

// TestConfigurationFiles tests configuration file support
func (v *CLIRuntimeValidator) TestConfigurationFiles(t *testing.T) {
	t.Helper()

	// Create a test config file
	configPath := filepath.Join(v.projectPath, ".test-cli-config-runtime.yaml")
	configContent := `logging:
  level: debug
  format: json
cli:
  output_format: json
`
	err := os.WriteFile(configPath, []byte(configContent), 0644)
	require.NoError(t, err, "Should create test config file")

	// Test using the config file
	output, err := v.runCLI("--config", configPath, "create", "task", "config-test", "--force")
	assert.NoError(t, err, "Command with config should succeed. Output: %s", output)
}

// TestEnvironmentVariables tests environment variable support
func (v *CLIRuntimeValidator) TestEnvironmentVariables(t *testing.T) {
	t.Helper()

	// Set environment variables for testing
	originalEnv := os.Environ()
	defer func() {
		// Restore original environment
		os.Clearenv()
		for _, env := range originalEnv {
			parts := strings.SplitN(env, "=", 2)
			if len(parts) == 2 {
				_ = os.Setenv(parts[0], parts[1])
			}
		}
	}()

	// Set test environment variables
	_ = os.Setenv("TEST_CLI_CONFIG_RUNTIME_LOGGING_LEVEL", "debug")
	
	// Test that environment variables are used
	output, err := v.runCLI("create", "task", "env-test", "--force")
	assert.NoError(t, err, "Command with env vars should succeed. Output: %s", output)
}

// runCLI executes the CLI with given arguments and returns output
func (v *CLIRuntimeValidator) runCLI(args ...string) (string, error) {
	cmd := exec.Command(v.cliPath, args...)
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	return string(output), err
}