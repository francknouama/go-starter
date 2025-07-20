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

// Feature: Simple CLI Blueprint (Issue #149)
// As a beginner Go developer
// I want to generate a simple CLI application project
// So that I can quickly build command-line tools without complexity
//
// Background: This blueprint addresses the CLI audit findings:
// - Complexity: 3/10 (vs standard 7/10)
// - Learning curve: 3/10 (vs standard 8/10)
// - File count: <10 files (vs standard 20+ files)
// - Essential features only

func TestSimple_CLI_BasicGeneration_MinimalComplexity(t *testing.T) {
	// Scenario: Generate simple CLI with minimal complexity
	// Given I want a simple CLI application with minimal complexity
	// When I generate with architecture "simple"
	// Then the project should have <10 files
	// And the project should use only essential dependencies
	// And the project should have complexity level 3/10
	// And the project should compile and run successfully

	// Given I want a simple CLI application with minimal complexity
	config := types.ProjectConfig{
		Name:         "test-simple-cli",
		Module:       "github.com/test/test-simple-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate with architecture "simple"
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should have <10 files
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateSimplicity(t)

	// And the project should use only essential dependencies
	validator.ValidateEssentialDependencies(t)

	// And the project should have complexity level 3/10
	validator.ValidateComplexityLevel(t)

	// And the project should compile and run successfully
	validator.ValidateCompilation(t)
}

func TestSimple_CLI_EssentialFeatures_Only(t *testing.T) {
	// Scenario: Simple CLI includes only essential features
	// Given I want a CLI with essential features only
	// When I generate a simple CLI project
	// Then the project should have help command support
	// And the project should have version command support
	// And the project should have quiet flag support
	// And the project should have output format support
	// And the project should NOT have complex business logic
	// And the project should NOT have factory patterns

	// Given I want a CLI with essential features only
	config := types.ProjectConfig{
		Name:         "test-essential-cli",
		Module:       "github.com/test/test-essential-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate a simple CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should have help command support
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateHelpSupport(t)

	// And the project should have version command support
	validator.ValidateVersionSupport(t)

	// And the project should have quiet flag support
	validator.ValidateQuietFlag(t)

	// And the project should have output format support
	validator.ValidateOutputFormat(t)

	// And the project should NOT have complex business logic
	validator.ValidateNoComplexBusinessLogic(t)

	// And the project should NOT have factory patterns
	validator.ValidateNoFactoryPatterns(t)
}

func TestSimple_CLI_SlogLogger_Integration(t *testing.T) {
	// Scenario: Simple CLI uses slog (standard library) only
	// Given I want a CLI with simple logging
	// When I generate a simple CLI project
	// Then the project should use slog from standard library
	// And the project should NOT have logger factory
	// And the project should NOT have logger interfaces
	// And the project should have direct slog usage
	// And logging should work correctly

	// Given I want a CLI with simple logging
	config := types.ProjectConfig{
		Name:         "test-slog-cli",
		Module:       "github.com/test/test-slog-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate a simple CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should use slog from standard library
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateSlogUsage(t)

	// And the project should NOT have logger factory
	validator.ValidateNoLoggerFactory(t)

	// And the project should NOT have logger interfaces
	validator.ValidateNoLoggerInterfaces(t)

	// And the project should have direct slog usage
	validator.ValidateDirectSlogUsage(t)

	// And logging should work correctly
	validator.ValidateLoggingFunctionality(t)
}

func TestSimple_CLI_ShellCompletion_Support(t *testing.T) {
	// Scenario: Simple CLI supports shell completion
	// Given I want a CLI with shell completion
	// When I generate a simple CLI project
	// Then the project should support bash completion
	// And the project should support zsh completion
	// And the project should support fish completion
	// And the project should support powershell completion
	// And completion should work without additional setup

	// Given I want a CLI with shell completion
	config := types.ProjectConfig{
		Name:         "test-completion-cli",
		Module:       "github.com/test/test-completion-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate a simple CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then completion should be supported by default
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateShellCompletionSupport(t)
}

func TestSimple_CLI_Configuration_Minimal(t *testing.T) {
	// Scenario: Simple CLI has minimal configuration
	// Given I want a CLI with simple configuration
	// When I generate a simple CLI project
	// Then the project should have environment variable support
	// And the project should have simple validation
	// And the project should NOT use Viper
	// And the project should NOT have complex config files
	// And configuration should be straightforward

	// Given I want a CLI with simple configuration
	config := types.ProjectConfig{
		Name:         "test-config-cli",
		Module:       "github.com/test/test-config-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate a simple CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should have environment variable support
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateEnvironmentVariables(t)

	// And the project should have simple validation
	validator.ValidateSimpleValidation(t)

	// And the project should NOT use Viper
	validator.ValidateNoViper(t)

	// And the project should NOT have complex config files
	validator.ValidateNoComplexConfigFiles(t)
}

func TestSimple_CLI_BeginnerFriendly_Structure(t *testing.T) {
	// Scenario: Simple CLI has beginner-friendly structure
	// Given I am a beginner Go developer
	// When I generate a simple CLI project
	// Then the project structure should be intuitive
	// And each file should have a clear purpose
	// And the code should be easy to understand
	// And the documentation should be comprehensive
	// And the learning curve should be minimal

	// Given I am a beginner Go developer
	config := types.ProjectConfig{
		Name:         "test-beginner-cli",
		Module:       "github.com/test/test-beginner-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I generate a simple CLI project
	projectPath := helpers.GenerateProject(t, config)

	// Then the project structure should be intuitive
	validator := NewSimpleCLIValidator(projectPath)
	validator.ValidateIntuitiveStructure(t)

	// And each file should have a clear purpose
	validator.ValidateClearFilePurpose(t)

	// And the code should be easy to understand
	validator.ValidateCodeReadability(t)

	// And the documentation should be comprehensive
	validator.ValidateDocumentation(t)
}

func TestSimple_CLI_Performance_Comparison(t *testing.T) {
	// Scenario: Simple CLI performs better than standard CLI
	// Given I want to compare simple vs standard CLI
	// When I generate both architectures
	// Then simple CLI should have fewer files
	// And simple CLI should have fewer dependencies
	// And simple CLI should compile faster
	// And simple CLI should be easier to maintain

	// Generate simple CLI
	simpleConfig := types.ProjectConfig{
		Name:         "test-simple-perf",
		Module:       "github.com/test/test-simple-perf",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	simplePath := helpers.GenerateProject(t, simpleConfig)

	// Generate standard CLI for comparison
	standardConfig := types.ProjectConfig{
		Name:         "test-standard-perf",
		Module:       "github.com/test/test-standard-perf",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	standardPath := helpers.GenerateProject(t, standardConfig)

	// Compare the architectures
	validator := NewSimpleCLIValidator(simplePath)
	validator.ValidatePerformanceComparison(t, standardPath)
}

// SimpleCLIValidator provides validation methods for simple CLI blueprints
type SimpleCLIValidator struct {
	projectPath string
}

// NewSimpleCLIValidator creates a new SimpleCLIValidator
func NewSimpleCLIValidator(projectPath string) *SimpleCLIValidator {
	return &SimpleCLIValidator{
		projectPath: projectPath,
	}
}

// ValidateSimplicity validates that the CLI is truly simple
func (v *SimpleCLIValidator) ValidateSimplicity(t *testing.T) {
	t.Helper()

	// Count total files
	fileCount := helpers.CountFiles(t, v.projectPath, ".go")
	if fileCount >= 10 {
		t.Errorf("Simple CLI should have <10 Go files, but has %d", fileCount)
	}

	// Count total files (all types)
	totalFiles := helpers.CountAllFiles(t, v.projectPath)
	if totalFiles >= 10 {
		t.Errorf("Simple CLI should have <10 total files, but has %d", totalFiles)
	}
}

// ValidateEssentialDependencies validates only essential dependencies
func (v *SimpleCLIValidator) ValidateEssentialDependencies(t *testing.T) {
	t.Helper()

	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileExists(t, goModFile)

	// Should only have cobra as direct dependency
	helpers.AssertFileContains(t, goModFile, "github.com/spf13/cobra")

	// Should NOT have Viper
	helpers.AssertFileNotContains(t, goModFile, "github.com/spf13/viper")

	// Should NOT have external loggers
	helpers.AssertFileNotContains(t, goModFile, "go.uber.org/zap")
	helpers.AssertFileNotContains(t, goModFile, "github.com/sirupsen/logrus")
	helpers.AssertFileNotContains(t, goModFile, "github.com/rs/zerolog")
}

// ValidateComplexityLevel validates complexity metrics
func (v *SimpleCLIValidator) ValidateComplexityLevel(t *testing.T) {
	t.Helper()

	// Check for absence of complex patterns
	v.ValidateNoFactoryPatterns(t)
	v.ValidateNoComplexInterfaces(t)
	v.ValidateNoDeepNesting(t)
}

// ValidateCompilation validates project compilation
func (v *SimpleCLIValidator) ValidateCompilation(t *testing.T) {
	t.Helper()
	helpers.AssertProjectCompiles(t, v.projectPath)
}

// ValidateHelpSupport validates help command functionality
func (v *SimpleCLIValidator) ValidateHelpSupport(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileExists(t, rootFile)
	helpers.AssertFileContains(t, rootFile, "cobra.Command")
	helpers.AssertFileContains(t, rootFile, "--help")
}

// ValidateVersionSupport validates version command
func (v *SimpleCLIValidator) ValidateVersionSupport(t *testing.T) {
	t.Helper()

	versionFile := filepath.Join(v.projectPath, "cmd", "version.go")
	helpers.AssertFileExists(t, versionFile)
	helpers.AssertFileContains(t, versionFile, "versionCmd")
	helpers.AssertFileContains(t, versionFile, "version")
}

// ValidateQuietFlag validates quiet flag support
func (v *SimpleCLIValidator) ValidateQuietFlag(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--quiet")
	helpers.AssertFileContains(t, rootFile, "quiet")
}

// ValidateOutputFormat validates output format support
func (v *SimpleCLIValidator) ValidateOutputFormat(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "--output")
	helpers.AssertFileContains(t, rootFile, "json")
	helpers.AssertFileContains(t, rootFile, "text")
}

// ValidateNoComplexBusinessLogic validates absence of complex business logic
func (v *SimpleCLIValidator) ValidateNoComplexBusinessLogic(t *testing.T) {
	t.Helper()

	// Should not have create/update/delete commands
	createFile := filepath.Join(v.projectPath, "cmd", "create.go")
	helpers.AssertFileNotExists(t, createFile)

	updateFile := filepath.Join(v.projectPath, "cmd", "update.go")
	helpers.AssertFileNotExists(t, updateFile)

	deleteFile := filepath.Join(v.projectPath, "cmd", "delete.go")
	helpers.AssertFileNotExists(t, deleteFile)

	listFile := filepath.Join(v.projectPath, "cmd", "list.go")
	helpers.AssertFileNotExists(t, listFile)
}

// ValidateNoFactoryPatterns validates absence of factory patterns
func (v *SimpleCLIValidator) ValidateNoFactoryPatterns(t *testing.T) {
	t.Helper()

	// Should not have logger factory
	factoryFile := filepath.Join(v.projectPath, "internal", "logger", "factory.go")
	helpers.AssertFileNotExists(t, factoryFile)

	// Should not have internal package structure
	internalDir := filepath.Join(v.projectPath, "internal")
	helpers.AssertDirectoryNotExists(t, internalDir)
}

// ValidateSlogUsage validates slog usage
func (v *SimpleCLIValidator) ValidateSlogUsage(t *testing.T) {
	t.Helper()

	mainFile := filepath.Join(v.projectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)
	helpers.AssertFileContains(t, mainFile, "log/slog")
	helpers.AssertFileContains(t, mainFile, "slog.New")
}

// ValidateNoLoggerFactory validates no logger factory
func (v *SimpleCLIValidator) ValidateNoLoggerFactory(t *testing.T) {
	t.Helper()

	loggerDir := filepath.Join(v.projectPath, "internal", "logger")
	helpers.AssertDirectoryNotExists(t, loggerDir)
}

// ValidateNoLoggerInterfaces validates no logger interfaces
func (v *SimpleCLIValidator) ValidateNoLoggerInterfaces(t *testing.T) {
	t.Helper()

	interfaceFile := filepath.Join(v.projectPath, "internal", "logger", "interface.go")
	helpers.AssertFileNotExists(t, interfaceFile)
}

// ValidateDirectSlogUsage validates direct slog usage
func (v *SimpleCLIValidator) ValidateDirectSlogUsage(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "slog.Info")
}

// ValidateLoggingFunctionality validates logging works
func (v *SimpleCLIValidator) ValidateLoggingFunctionality(t *testing.T) {
	t.Helper()
	helpers.AssertLoggerFunctionality(t, v.projectPath, "slog")
}

// ValidateShellCompletionSupport validates shell completion
func (v *SimpleCLIValidator) ValidateShellCompletionSupport(t *testing.T) {
	t.Helper()

	rootFile := filepath.Join(v.projectPath, "cmd", "root.go")
	helpers.AssertFileContains(t, rootFile, "CompletionOptions")
}

// ValidateEnvironmentVariables validates environment variable support
func (v *SimpleCLIValidator) ValidateEnvironmentVariables(t *testing.T) {
	t.Helper()

	configFile := filepath.Join(v.projectPath, "config.go")
	helpers.AssertFileExists(t, configFile)
	helpers.AssertFileContains(t, configFile, "os.Getenv")
}

// ValidateSimpleValidation validates simple validation
func (v *SimpleCLIValidator) ValidateSimpleValidation(t *testing.T) {
	t.Helper()

	configFile := filepath.Join(v.projectPath, "config.go")
	helpers.AssertFileContains(t, configFile, "Validate")
}

// ValidateNoViper validates no Viper usage
func (v *SimpleCLIValidator) ValidateNoViper(t *testing.T) {
	t.Helper()

	goModFile := filepath.Join(v.projectPath, "go.mod")
	helpers.AssertFileNotContains(t, goModFile, "github.com/spf13/viper")
}

// ValidateNoComplexConfigFiles validates no complex config files
func (v *SimpleCLIValidator) ValidateNoComplexConfigFiles(t *testing.T) {
	t.Helper()

	configsDir := filepath.Join(v.projectPath, "configs")
	helpers.AssertDirectoryNotExists(t, configsDir)
}

// ValidateIntuitiveStructure validates intuitive project structure
func (v *SimpleCLIValidator) ValidateIntuitiveStructure(t *testing.T) {
	t.Helper()

	// Should have clear main.go
	mainFile := filepath.Join(v.projectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)

	// Should have cmd directory for commands
	cmdDir := filepath.Join(v.projectPath, "cmd")
	helpers.AssertDirectoryExists(t, cmdDir)

	// Should have simple config.go
	configFile := filepath.Join(v.projectPath, "config.go")
	helpers.AssertFileExists(t, configFile)
}

// ValidateClearFilePurpose validates clear file purpose
func (v *SimpleCLIValidator) ValidateClearFilePurpose(t *testing.T) {
	t.Helper()

	// Each file should have clear comments explaining purpose
	files := []string{
		"main.go",
		"config.go",
		"cmd/root.go",
		"cmd/version.go",
	}

	for _, file := range files {
		filePath := filepath.Join(v.projectPath, file)
		helpers.AssertFileExists(t, filePath)
		// Each file should have some documentation
		helpers.AssertFileContains(t, filePath, "//")
	}
}

// ValidateCodeReadability validates code readability
func (v *SimpleCLIValidator) ValidateCodeReadability(t *testing.T) {
	t.Helper()

	// Main should be simple and clear
	mainFile := filepath.Join(v.projectPath, "main.go")
	content := helpers.ReadFile(t, mainFile)

	// Should be short and readable
	lines := helpers.CountLines(content)
	if lines > 30 {
		t.Errorf("main.go should be simple (<30 lines), but has %d lines", lines)
	}
}

// ValidateDocumentation validates documentation quality
func (v *SimpleCLIValidator) ValidateDocumentation(t *testing.T) {
	t.Helper()

	readmeFile := filepath.Join(v.projectPath, "README.md")
	helpers.AssertFileExists(t, readmeFile)
	helpers.AssertFileContains(t, readmeFile, "Simple")
	helpers.AssertFileContains(t, readmeFile, "Usage")
	helpers.AssertFileContains(t, readmeFile, "Environment Variables")
}

// ValidatePerformanceComparison validates performance vs standard
func (v *SimpleCLIValidator) ValidatePerformanceComparison(t *testing.T, standardPath string) {
	t.Helper()

	// Count files in both projects
	simpleFiles := helpers.CountAllFiles(t, v.projectPath)
	standardFiles := helpers.CountAllFiles(t, standardPath)

	if simpleFiles >= standardFiles {
		t.Errorf("Simple CLI should have fewer files than standard. Simple: %d, Standard: %d", simpleFiles, standardFiles)
	}

	// Check dependencies
	simpleMod := filepath.Join(v.projectPath, "go.mod")
	standardMod := filepath.Join(standardPath, "go.mod")

	simpleContent := helpers.ReadFile(t, simpleMod)
	standardContent := helpers.ReadFile(t, standardMod)

	simpleDeps := helpers.CountDependencies(simpleContent)
	standardDeps := helpers.CountDependencies(standardContent)

	if simpleDeps >= standardDeps {
		t.Errorf("Simple CLI should have fewer dependencies than standard. Simple: %d, Standard: %d", simpleDeps, standardDeps)
	}
}

// Helper validation methods

// ValidateNoComplexInterfaces validates absence of complex interfaces
func (v *SimpleCLIValidator) ValidateNoComplexInterfaces(t *testing.T) {
	t.Helper()

	// Should not have multiple interface definitions
	files := helpers.FindGoFiles(t, v.projectPath)
	for _, file := range files {
		content := helpers.ReadFile(t, file)
		interfaceCount := helpers.CountInterfaces(content)
		if interfaceCount > 1 {
			t.Errorf("File %s has too many interfaces (%d) for a simple CLI", file, interfaceCount)
		}
	}
}

// ValidateNoDeepNesting validates absence of deep nesting
func (v *SimpleCLIValidator) ValidateNoDeepNesting(t *testing.T) {
	t.Helper()

	// Check for deep directory structures
	maxDepth := helpers.GetMaxDirectoryDepth(t, v.projectPath)
	if maxDepth > 3 {
		t.Errorf("Simple CLI should not have deep nesting (max 3 levels), but has %d levels", maxDepth)
	}
}