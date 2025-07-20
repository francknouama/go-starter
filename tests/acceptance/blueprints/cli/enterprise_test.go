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

// Feature: CLI-Enterprise Enhancement (Issue #56)
// As an enterprise Go developer
// I want enhanced CLI capabilities with professional standards
// So that I can build enterprise-grade command-line applications
//
// Background: This enhancement upgrades the cli-standard blueprint:
// - Complexity: 7/10 (enterprise-level)
// - Professional CLI standards: --quiet, --no-color, --output
// - Better command organization with groups
// - Improved error handling and validation
// - Shell completion support
// - Interactive mode for complex commands
// - Progressive disclosure for advanced flags

func TestEnterprise_CLI_Standards_Compliance(t *testing.T) {
	// Scenario: CLI supports standard professional flags
	// Given I want a CLI with professional standards
	// When I generate an enterprise CLI project
	// Then the project should support --quiet flag globally
	// And the project should support --no-color flag globally
	// And the project should support --output flag with table|json|yaml formats
	// And these flags should be persistent across all commands
	// And the help output should show these global flags

	// Given I want a CLI with professional standards
	config := types.ProjectConfig{
		Name:         "test-enterprise-cli",
		Module:       "github.com/test/test-enterprise-cli",
		Type:         "cli",
		Architecture: "standard", // Enhanced standard becomes enterprise
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should support standard flags
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateQuietFlag(t)
	validator.ValidateNoColorFlag(t)
	validator.ValidateOutputFlag(t)
	validator.ValidatePersistentFlags(t)
	validator.ValidateGlobalFlagsInHelp(t)
}

func TestEnterprise_CLI_Command_Organization(t *testing.T) {
	// Scenario: Commands are organized into logical groups
	// Given I want a CLI with organized commands
	// When I generate an enterprise CLI project
	// Then commands should be grouped logically (e.g., "manage" group)
	// And the help output should show organized command structure
	// And subcommands should be properly nested
	// And each group should have clear descriptions
	// And navigation should be intuitive

	// Given I want a CLI with organized commands
	config := types.ProjectConfig{
		Name:         "test-organized-cli",
		Module:       "github.com/test/test-organized-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then commands should be organized properly
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateCommandGroups(t)
	validator.ValidateHelpStructure(t)
	validator.ValidateSubcommandOrganization(t)
	validator.ValidateGroupDescriptions(t)
}

func TestEnterprise_CLI_Error_Handling(t *testing.T) {
	// Scenario: Improved error handling and validation
	// Given I want a CLI with professional error handling
	// When I generate an enterprise CLI project
	// Then the project should use cobra.ExactArgs for validation
	// And error messages should be graceful and informative
	// And errors should propagate correctly through command chain
	// And validation should happen before execution
	// And error formatting should be consistent

	// Given I want a CLI with professional error handling
	config := types.ProjectConfig{
		Name:         "test-error-cli",
		Module:       "github.com/test/test-error-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then error handling should be professional
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateExactArgsUsage(t)
	validator.ValidateGracefulErrors(t)
	validator.ValidateErrorPropagation(t)
	validator.ValidatePreExecutionValidation(t)
	validator.ValidateConsistentErrorFormat(t)
}

func TestEnterprise_CLI_Shell_Completion(t *testing.T) {
	// Scenario: Advanced shell completion support
	// Given I want a CLI with shell completion
	// When I generate an enterprise CLI project
	// Then the project should have a completion command
	// And bash completion should be supported
	// And zsh completion should be supported
	// And fish completion should be supported
	// And completion should be context-aware
	// And custom completion functions should be available

	// Given I want a CLI with shell completion
	config := types.ProjectConfig{
		Name:         "test-completion-cli",
		Module:       "github.com/test/test-completion-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then shell completion should be enterprise-grade
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateCompletionCommand(t)
	validator.ValidateBashCompletion(t)
	validator.ValidateZshCompletion(t)
	validator.ValidateFishCompletion(t)
	validator.ValidateContextAwareCompletion(t)
	validator.ValidateCustomCompletionFunctions(t)
}

func TestEnterprise_CLI_Interactive_Mode(t *testing.T) {
	// Scenario: Interactive mode for complex commands
	// Given I want a CLI with interactive mode
	// When I generate an enterprise CLI project
	// Then the project should support interactive prompts
	// And complex commands should have interactive alternatives
	// And validation should work in interactive mode
	// And the user experience should be smooth
	// And interactive mode should be optional

	// Given I want a CLI with interactive mode
	config := types.ProjectConfig{
		Name:         "test-interactive-cli",
		Module:       "github.com/test/test-interactive-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
		Features: &types.Features{
			// Interactive mode will be handled via variables
		},
		Variables: map[string]string{
			"Interactive": "true", // Enable interactive mode
		},
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then interactive mode should be available
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateInteractiveSupport(t)
	validator.ValidateInteractivePrompts(t)
	validator.ValidateInteractiveValidation(t)
	validator.ValidateInteractiveUX(t)
	validator.ValidateInteractiveOptional(t)
}

func TestEnterprise_CLI_Progressive_Disclosure(t *testing.T) {
	// Scenario: Progressive disclosure for advanced features
	// Given I want a CLI with progressive disclosure
	// When I generate an enterprise CLI project
	// Then basic usage should hide complex features
	// And --advanced flag should reveal additional options
	// And help should adapt based on user level
	// And documentation should be layered
	// And the learning curve should be manageable

	// Given I want a CLI with progressive disclosure
	config := types.ProjectConfig{
		Name:         "test-progressive-cli",
		Module:       "github.com/test/test-progressive-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then progressive disclosure should work
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateBasicUsageSimplicity(t)
	validator.ValidateAdvancedFlag(t)
	validator.ValidateAdaptiveHelp(t)
	validator.ValidateLayeredDocumentation(t)
	validator.ValidateManageableLearningCurve(t)
}

func TestEnterprise_CLI_Enterprise_Positioning(t *testing.T) {
	// Scenario: Enterprise-grade positioning and capabilities
	// Given I want an enterprise-grade CLI
	// When I generate an enterprise CLI project
	// Then the project should indicate 7/10 complexity
	// And the project should handle 10+ commands well
	// And configuration management should be robust
	// And the architecture should be scalable
	// And professional patterns should be used

	// Given I want an enterprise-grade CLI
	config := types.ProjectConfig{
		Name:         "test-enterprise-grade-cli",
		Module:       "github.com/test/test-enterprise-grade-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then it should be enterprise-grade
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateComplexityLevel(t)
	validator.ValidateLargeCommandHandling(t)
	validator.ValidateRobustConfiguration(t)
	validator.ValidateScalableArchitecture(t)
	validator.ValidateProfessionalPatterns(t)
}

func TestEnterprise_CLI_Configuration_Management(t *testing.T) {
	// Scenario: Enterprise configuration management
	// Given I want enterprise-level configuration
	// When I generate an enterprise CLI project
	// Then the project should support multiple config sources
	// And configuration should support environments
	// And secrets should be handled securely
	// And configuration validation should be comprehensive
	// And hot reloading should be supported

	// Given I want enterprise-level configuration
	config := types.ProjectConfig{
		Name:         "test-config-enterprise-cli",
		Module:       "github.com/test/test-config-enterprise-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
		Features: &types.Features{
			// Configuration will be handled via variables
		},
		Variables: map[string]string{
			"Configuration": "viper", // Enterprise uses Viper
		},
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then configuration should be enterprise-grade
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateMultipleConfigSources(t)
	validator.ValidateEnvironmentSupport(t)
	validator.ValidateSecureSecretsHandling(t)
	validator.ValidateComprehensiveValidation(t)
	validator.ValidateHotReloading(t)
}

func TestEnterprise_CLI_Logging_Standards(t *testing.T) {
	// Scenario: Enterprise logging standards
	// Given I want enterprise logging capabilities
	// When I generate an enterprise CLI project
	// Then logging should support structured formats
	// And log levels should be configurable
	// And log output should respect --quiet flag
	// And colors should respect --no-color flag
	// And logging should integrate with output formats

	// Given I want enterprise logging capabilities
	config := types.ProjectConfig{
		Name:         "test-logging-enterprise-cli",
		Module:       "github.com/test/test-logging-enterprise-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "zap", // Enterprise might prefer zap
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then logging should be enterprise-grade
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateStructuredLogging(t)
	validator.ValidateConfigurableLogLevels(t)
	validator.ValidateQuietModeIntegration(t)
	validator.ValidateNoColorIntegration(t)
	validator.ValidateLoggingOutputIntegration(t)
}

func TestEnterprise_CLI_Build_Performance(t *testing.T) {
	// Scenario: Optimized build and runtime performance
	// Given I want optimized CLI performance
	// When I generate an enterprise CLI project
	// Then the project should compile efficiently
	// And binary size should be reasonable
	// And startup time should be fast
	// And memory usage should be optimized
	// And the project should scale to many commands

	// Given I want optimized CLI performance
	config := types.ProjectConfig{
		Name:         "test-performance-cli",
		Module:       "github.com/test/test-performance-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate an enterprise CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then performance should be optimized
	validator := NewEnterpriseCLIValidator(projectPath)
	validator.ValidateEfficientCompilation(t)
	validator.ValidateReasonableBinarySize(t)
	validator.ValidateFastStartup(t)
	validator.ValidateMemoryOptimization(t)
	validator.ValidateCommandScalability(t)
}

// EnterpriseCLIValidator provides validation methods for enterprise CLI blueprints
type EnterpriseCLIValidator struct {
	projectPath string
}

// NewEnterpriseCLIValidator creates a new EnterpriseCLIValidator
func NewEnterpriseCLIValidator(projectPath string) *EnterpriseCLIValidator {
	return &EnterpriseCLIValidator{
		projectPath: projectPath,
	}
}

// CLI Standards Compliance validations

func (v *EnterpriseCLIValidator) ValidateQuietFlag(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "--quiet")
	helpers.AssertFileContains(t, rootFile, "\"Suppress all output\"")
}

func (v *EnterpriseCLIValidator) ValidateNoColorFlag(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--no-color")
	helpers.AssertFileContains(t, rootFile, "\"Disable colored output\"")
}

func (v *EnterpriseCLIValidator) ValidateOutputFlag(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--output")
	helpers.AssertFileContains(t, rootFile, "table")
	helpers.AssertFileContains(t, rootFile, "json")
	helpers.AssertFileContains(t, rootFile, "yaml")
}

func (v *EnterpriseCLIValidator) ValidatePersistentFlags(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "PersistentFlags()")
	helpers.AssertFileContains(t, rootFile, "rootCmd.PersistentFlags()")
}

func (v *EnterpriseCLIValidator) ValidateGlobalFlagsInHelp(t *testing.T) {
	t.Helper()
	// Build and run help to verify global flags appear
	helpers.AssertCLIHelpOutput(t, v.projectPath)
}

// Command Organization validations

func (v *EnterpriseCLIValidator) ValidateCommandGroups(t *testing.T) {
	t.Helper()
	// Check for command group organization
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "GroupID")
	helpers.AssertFileContains(t, rootFile, "manage")
}

func (v *EnterpriseCLIValidator) ValidateHelpStructure(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "SetHelpCommandGroupID")
}

func (v *EnterpriseCLIValidator) ValidateSubcommandOrganization(t *testing.T) {
	t.Helper()
	// Check for well-organized subcommands
	manageDir := filepath.Join(v.projectPath, "cmd", "manage")
	if helpers.DirExists(manageDir) {
		// Validate subcommand structure
		entries, _ := os.ReadDir(manageDir)
		if len(entries) < 2 {
			t.Error("Manage group should have multiple subcommands")
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateGroupDescriptions(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "Title:")
}

// Error Handling validations

func (v *EnterpriseCLIValidator) ValidateExactArgsUsage(t *testing.T) {
	t.Helper()
	// Find command files and check for ExactArgs usage
	cmdFiles := helpers.FindFiles(t, filepath.Join(v.projectPath, "cmd"), "*.go")
	foundExactArgs := false
	for _, file := range cmdFiles {
		content := helpers.ReadFile(t, file)
		if strings.Contains(content, "cobra.ExactArgs") {
			foundExactArgs = true
			break
		}
	}
	if !foundExactArgs {
		t.Error("Enterprise CLI should use cobra.ExactArgs for validation")
	}
}

func (v *EnterpriseCLIValidator) ValidateGracefulErrors(t *testing.T) {
	t.Helper()
	// Check for error handling patterns
	mainFile := filepath.Join(v.projectPath, "main.go")
	helpers.AssertFileContains(t, mainFile, "if err != nil")
	helpers.AssertFileContains(t, mainFile, "os.Exit")
}

func (v *EnterpriseCLIValidator) ValidateErrorPropagation(t *testing.T) {
	t.Helper()
	// Check for proper error returns in command files
	cmdFiles := helpers.FindFiles(t, filepath.Join(v.projectPath, "cmd"), "*.go")
	for _, file := range cmdFiles {
		content := helpers.ReadFile(t, file)
		if strings.Contains(content, "RunE:") {
			// Commands should use RunE for error propagation
			return
		}
	}
	t.Error("Commands should use RunE for proper error propagation")
}

func (v *EnterpriseCLIValidator) ValidatePreExecutionValidation(t *testing.T) {
	t.Helper()
	// Check for PreRunE or validation functions
	cmdFiles := helpers.FindFiles(t, filepath.Join(v.projectPath, "cmd"), "*.go")
	foundValidation := false
	for _, file := range cmdFiles {
		content := helpers.ReadFile(t, file)
		if strings.Contains(content, "PreRunE:") || strings.Contains(content, "validate") {
			foundValidation = true
			break
		}
	}
	if !foundValidation {
		t.Error("Enterprise CLI should have pre-execution validation")
	}
}

func (v *EnterpriseCLIValidator) ValidateConsistentErrorFormat(t *testing.T) {
	t.Helper()
	// Check for error formatting utilities
	errorsFile := filepath.Join(v.projectPath, "internal", "errors", "errors.go")
	if helpers.FileExists(errorsFile) {
		helpers.AssertFileContains(t, errorsFile, "Format")
	}
}

// Shell Completion validations

func (v *EnterpriseCLIValidator) ValidateCompletionCommand(t *testing.T) {
	t.Helper()
	completionFile := filepath.Join(v.projectPath, "cmd", "completion.go")
	helpers.AssertFileExists(t, completionFile)
	helpers.AssertFileContains(t, completionFile, "completionCmd")
}

func (v *EnterpriseCLIValidator) ValidateBashCompletion(t *testing.T) {
	t.Helper()
	completionFile := filepath.Join(v.projectPath, "cmd", "completion.go")
	helpers.AssertFileContains(t, completionFile, "bash")
	helpers.AssertFileContains(t, completionFile, "GenBashCompletion")
}

func (v *EnterpriseCLIValidator) ValidateZshCompletion(t *testing.T) {
	t.Helper()
	completionFile := filepath.Join(v.projectPath, "cmd", "completion.go")
	helpers.AssertFileContains(t, completionFile, "zsh")
	helpers.AssertFileContains(t, completionFile, "GenZshCompletion")
}

func (v *EnterpriseCLIValidator) ValidateFishCompletion(t *testing.T) {
	t.Helper()
	completionFile := filepath.Join(v.projectPath, "cmd", "completion.go")
	helpers.AssertFileContains(t, completionFile, "fish")
	helpers.AssertFileContains(t, completionFile, "GenFishCompletion")
}

func (v *EnterpriseCLIValidator) ValidateContextAwareCompletion(t *testing.T) {
	t.Helper()
	// Check for custom completion functions
	cmdFiles := helpers.FindFiles(t, filepath.Join(v.projectPath, "cmd"), "*.go")
	foundCompletion := false
	for _, file := range cmdFiles {
		content := helpers.ReadFile(t, file)
		if strings.Contains(content, "ValidArgsFunction") || strings.Contains(content, "RegisterFlagCompletionFunc") {
			foundCompletion = true
			break
		}
	}
	if !foundCompletion {
		t.Error("Enterprise CLI should have context-aware completion")
	}
}

func (v *EnterpriseCLIValidator) ValidateCustomCompletionFunctions(t *testing.T) {
	t.Helper()
	completionFile := filepath.Join(v.projectPath, "internal", "completion", "completion.go")
	if helpers.FileExists(completionFile) {
		helpers.AssertFileContains(t, completionFile, "CompleteFunc")
	}
}

// Interactive Mode validations

func (v *EnterpriseCLIValidator) ValidateInteractiveSupport(t *testing.T) {
	t.Helper()
	// Check for interactive package
	interactiveDir := filepath.Join(v.projectPath, "internal", "interactive")
	if !helpers.DirExists(interactiveDir) {
		t.Error("Enterprise CLI with interactive feature should have interactive package")
	}
}

func (v *EnterpriseCLIValidator) ValidateInteractivePrompts(t *testing.T) {
	t.Helper()
	// Check for prompt implementations
	promptFile := filepath.Join(v.projectPath, "internal", "interactive", "prompt.go")
	if helpers.FileExists(promptFile) {
		helpers.AssertFileContains(t, promptFile, "Prompt")
		helpers.AssertFileContains(t, promptFile, "survey")
	}
}

func (v *EnterpriseCLIValidator) ValidateInteractiveValidation(t *testing.T) {
	t.Helper()
	// Check for validation in interactive mode
	promptFile := filepath.Join(v.projectPath, "internal", "interactive", "prompt.go")
	if helpers.FileExists(promptFile) {
		helpers.AssertFileContains(t, promptFile, "Validate")
	}
}

func (v *EnterpriseCLIValidator) ValidateInteractiveUX(t *testing.T) {
	t.Helper()
	// Check for good UX patterns
	promptFile := filepath.Join(v.projectPath, "internal", "interactive", "prompt.go")
	if helpers.FileExists(promptFile) {
		helpers.AssertFileContains(t, promptFile, "Help")
		helpers.AssertFileContains(t, promptFile, "Default")
	}
}

func (v *EnterpriseCLIValidator) ValidateInteractiveOptional(t *testing.T) {
	t.Helper()
	// Check that interactive mode is optional
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--interactive")
}

// Progressive Disclosure validations

func (v *EnterpriseCLIValidator) ValidateBasicUsageSimplicity(t *testing.T) {
	t.Helper()
	// Check that basic commands are simple
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	content := helpers.ReadFile(t, rootFile)
	// Basic usage should not show all flags
	if strings.Count(content, "Flags()") > 10 {
		t.Error("Basic usage should hide complex features")
	}
}

func (v *EnterpriseCLIValidator) ValidateAdvancedFlag(t *testing.T) {
	t.Helper()
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--advanced")
	helpers.AssertFileContains(t, rootFile, "\"Show advanced options\"")
}

func (v *EnterpriseCLIValidator) ValidateAdaptiveHelp(t *testing.T) {
	t.Helper()
	// Check for adaptive help implementation
	helpFile := filepath.Join(v.projectPath, "internal", "help", "help.go")
	if helpers.FileExists(helpFile) {
		helpers.AssertFileContains(t, helpFile, "Advanced")
	}
}

func (v *EnterpriseCLIValidator) ValidateLayeredDocumentation(t *testing.T) {
	t.Helper()
	// Check README for layered docs
	readmeFile := filepath.Join(v.projectPath, "README.md")
	helpers.AssertFileContains(t, readmeFile, "## Basic Usage")
	helpers.AssertFileContains(t, readmeFile, "## Advanced Usage")
}

func (v *EnterpriseCLIValidator) ValidateManageableLearningCurve(t *testing.T) {
	t.Helper()
	// Check for gradual complexity increase
	examplesDir := filepath.Join(v.projectPath, "examples")
	if helpers.DirExists(examplesDir) {
		entries, _ := os.ReadDir(examplesDir)
		if len(entries) < 2 {
			t.Error("Should have multiple examples showing progression")
		}
	}
}

// Enterprise Positioning validations

func (v *EnterpriseCLIValidator) ValidateComplexityLevel(t *testing.T) {
	t.Helper()
	// Check file count and structure complexity
	fileCount := helpers.CountFiles(t, v.projectPath, ".go")
	if fileCount < 20 {
		t.Errorf("Enterprise CLI should have 20+ Go files, but has %d", fileCount)
	}
	
	// Check for enterprise patterns
	internalDir := filepath.Join(v.projectPath, "internal")
	helpers.AssertDirectoryExists(t, internalDir)
	
	// Should have multiple internal packages
	entries, _ := os.ReadDir(internalDir)
	if len(entries) < 5 {
		t.Error("Enterprise CLI should have multiple internal packages")
	}
}

func (v *EnterpriseCLIValidator) ValidateLargeCommandHandling(t *testing.T) {
	t.Helper()
	// Check for command organization that scales
	cmdDir := filepath.Join(v.projectPath, "cmd")
	entries, _ := os.ReadDir(cmdDir)
	
	// Should support many commands
	if len(entries) < 5 {
		t.Error("Enterprise CLI should demonstrate handling of multiple commands")
	}
}

func (v *EnterpriseCLIValidator) ValidateRobustConfiguration(t *testing.T) {
	t.Helper()
	// Check for robust configuration
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	helpers.AssertFileExists(t, configFile)
	helpers.AssertFileContains(t, configFile, "Config")
	helpers.AssertFileContains(t, configFile, "Load")
	helpers.AssertFileContains(t, configFile, "Validate")
}

func (v *EnterpriseCLIValidator) ValidateScalableArchitecture(t *testing.T) {
	t.Helper()
	// Check for scalable patterns
	// Should have clear separation of concerns
	dirs := []string{
		"internal/config",
		"internal/logger",
		"internal/errors",
		"internal/version",
	}
	
	for _, dir := range dirs {
		dirPath := filepath.Join(v.projectPath, dir)
		if !helpers.DirExists(dirPath) {
			t.Errorf("Enterprise CLI should have %s for scalable architecture", dir)
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateProfessionalPatterns(t *testing.T) {
	t.Helper()
	// Check for professional Go patterns
	
	// Should use interfaces
	files := helpers.FindGoFiles(t, v.projectPath)
	foundInterface := false
	for _, file := range files {
		content := helpers.ReadFile(t, file)
		if strings.Contains(content, "type") && strings.Contains(content, "interface") {
			foundInterface = true
			break
		}
	}
	if !foundInterface {
		t.Error("Enterprise CLI should use interfaces for abstraction")
	}
	
	// Should have proper error types
	errorsDir := filepath.Join(v.projectPath, "internal", "errors")
	if helpers.DirExists(errorsDir) {
		files := helpers.FindGoFiles(t, errorsDir)
		if len(files) == 0 {
			t.Error("Enterprise CLI should have custom error types")
		}
	}
}

// Configuration Management validations

func (v *EnterpriseCLIValidator) ValidateMultipleConfigSources(t *testing.T) {
	t.Helper()
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	if helpers.FileExists(configFile) {
		helpers.AssertFileContains(t, configFile, "viper")
		helpers.AssertFileContains(t, configFile, "SetConfigFile")
		helpers.AssertFileContains(t, configFile, "SetEnvPrefix")
	}
}

func (v *EnterpriseCLIValidator) ValidateEnvironmentSupport(t *testing.T) {
	t.Helper()
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	if helpers.FileExists(configFile) {
		helpers.AssertFileContains(t, configFile, "AutomaticEnv")
		helpers.AssertFileContains(t, configFile, "BindEnv")
	}
}

func (v *EnterpriseCLIValidator) ValidateSecureSecretsHandling(t *testing.T) {
	t.Helper()
	// Check for secure handling patterns
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	if helpers.FileExists(configFile) {
		content := helpers.ReadFile(t, configFile)
		// Should not log secrets
		if strings.Contains(content, "Secret") || strings.Contains(content, "Password") {
			if !strings.Contains(content, "redact") && !strings.Contains(content, "mask") {
				t.Error("Enterprise CLI should handle secrets securely")
			}
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateComprehensiveValidation(t *testing.T) {
	t.Helper()
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	if helpers.FileExists(configFile) {
		helpers.AssertFileContains(t, configFile, "Validate")
		// Should validate required fields
		helpers.AssertFileContains(t, configFile, "required")
	}
}

func (v *EnterpriseCLIValidator) ValidateHotReloading(t *testing.T) {
	t.Helper()
	configFile := filepath.Join(v.projectPath, "internal", "config", "config.go")
	if helpers.FileExists(configFile) {
		helpers.AssertFileContains(t, configFile, "WatchConfig")
		helpers.AssertFileContains(t, configFile, "OnConfigChange")
	}
}

// Logging Standards validations

func (v *EnterpriseCLIValidator) ValidateStructuredLogging(t *testing.T) {
	t.Helper()
	loggerFile := filepath.Join(v.projectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFile(t, loggerFile)
		// Should support structured logging
		if !strings.Contains(content, "WithField") && !strings.Contains(content, "With") {
			t.Error("Enterprise CLI should support structured logging")
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateConfigurableLogLevels(t *testing.T) {
	t.Helper()
	loggerFile := filepath.Join(v.projectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		helpers.AssertFileContains(t, loggerFile, "SetLevel")
		helpers.AssertFileContains(t, loggerFile, "Debug")
		helpers.AssertFileContains(t, loggerFile, "Info")
		helpers.AssertFileContains(t, loggerFile, "Warn")
		helpers.AssertFileContains(t, loggerFile, "Error")
	}
}

func (v *EnterpriseCLIValidator) ValidateQuietModeIntegration(t *testing.T) {
	t.Helper()
	// Check that logger respects quiet mode
	loggerFile := filepath.Join(v.projectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFile(t, loggerFile)
		if !strings.Contains(content, "quiet") && !strings.Contains(content, "Quiet") {
			t.Error("Logger should integrate with --quiet flag")
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateNoColorIntegration(t *testing.T) {
	t.Helper()
	// Check that logger respects no-color mode
	loggerFile := filepath.Join(v.projectPath, "internal", "logger", "logger.go")
	if helpers.FileExists(loggerFile) {
		content := helpers.ReadFile(t, loggerFile)
		if !strings.Contains(content, "color") && !strings.Contains(content, "Color") {
			t.Error("Logger should integrate with --no-color flag")
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateLoggingOutputIntegration(t *testing.T) {
	t.Helper()
	// Check that logger can output in different formats
	outputFile := filepath.Join(v.projectPath, "internal", "output", "output.go")
	if helpers.FileExists(outputFile) {
		helpers.AssertFileContains(t, outputFile, "Format")
		helpers.AssertFileContains(t, outputFile, "json")
		helpers.AssertFileContains(t, outputFile, "yaml")
	}
}

// Build Performance validations

func (v *EnterpriseCLIValidator) ValidateEfficientCompilation(t *testing.T) {
	t.Helper()
	// Ensure project compiles
	helpers.AssertProjectCompiles(t, v.projectPath)
}

func (v *EnterpriseCLIValidator) ValidateReasonableBinarySize(t *testing.T) {
	t.Helper()
	// Build the binary
	cmd := exec.Command("go", "build", "-o", "test-cli", ".")
	cmd.Dir = v.projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Failed to build: %v\nOutput: %s", err, string(output))
		return
	}
	
	// Check binary size (should be reasonable for enterprise CLI)
	binaryPath := filepath.Join(v.projectPath, "test-cli")
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Errorf("Failed to stat binary: %v", err)
		return
	}
	
	// Enterprise CLI might be larger but should still be reasonable
	maxSize := int64(30 * 1024 * 1024) // 30MB
	if info.Size() > maxSize {
		t.Errorf("Binary size %d bytes exceeds reasonable limit of %d bytes", info.Size(), maxSize)
	}
}

func (v *EnterpriseCLIValidator) ValidateFastStartup(t *testing.T) {
	t.Helper()
	// Check for lazy loading patterns
	mainFile := filepath.Join(v.projectPath, "main.go")
	content := helpers.ReadFile(t, mainFile)
	
	// Should not initialize everything in main
	if strings.Count(content, "init") > 5 {
		t.Error("Main should use lazy initialization for fast startup")
	}
}

func (v *EnterpriseCLIValidator) ValidateMemoryOptimization(t *testing.T) {
	t.Helper()
	// Check for memory-efficient patterns
	files := helpers.FindGoFiles(t, v.projectPath)
	
	for _, file := range files {
		content := helpers.ReadFile(t, file)
		// Check for common memory issues
		if strings.Contains(content, "make(") && strings.Contains(content, "1000000") {
			t.Error("Should avoid pre-allocating large slices")
		}
	}
}

func (v *EnterpriseCLIValidator) ValidateCommandScalability(t *testing.T) {
	t.Helper()
	// Check that command registration is scalable
	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	content := helpers.ReadFile(t, rootFile)
	
	// Should use init() functions for command registration
	if !strings.Contains(content, "init()") {
		t.Error("Should use init() for scalable command registration")
	}
	
	// Should not have all commands in one file
	cmdDir := filepath.Join(v.projectPath, "cmd")
	files, _ := os.ReadDir(cmdDir)
	if len(files) < 5 {
		t.Error("Commands should be organized in separate files for scalability")
	}
}