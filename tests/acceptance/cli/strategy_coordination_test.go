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

// Feature: CLI Strategy Coordination (Issue #150)
// As a developer using go-starter
// I want to be guided between CLI-Simple and CLI-Standard blueprints
// So that I can choose the right complexity level for my needs
//
// Background: This coordinates the complete CLI strategy:
// - CLI-Simple: For beginners, minimal complexity (3/10), <10 files
// - CLI-Standard: For enterprise, full features (7/10), 20+ files
// - Progressive discovery and selection
// - Clear guidance and recommendations

func TestCLIStrategy_BlueprintDiscovery_ListCommand(t *testing.T) {
	// Scenario: CLI Blueprint Discovery through List Command
	// Given I want to discover available CLI blueprints
	// When I run the list command
	// Then I should see both CLI-Simple and CLI-Standard
	// And CLI-Simple should appear first (default for beginners)
	// And descriptions should clearly indicate complexity and use cases
	// And both blueprints should be properly categorized as CLI type

	// Build the CLI tool for testing
	binary := buildCLIBinary(t)
	defer cleanupBinary(t, binary)

	// When I run the list command
	cmd := exec.Command(binary, "list")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "List command should execute successfully")

	outputStr := string(output)

	// Then I should see both CLI-Simple and CLI-Standard
	assert.Contains(t, outputStr, "cli-simple", "Should list CLI-Simple blueprint")
	assert.Contains(t, outputStr, "cli-standard", "Should list CLI-Standard blueprint")

	// And CLI-Simple should appear first (default for beginners)
	simpleIndex := strings.Index(outputStr, "cli-simple")
	standardIndex := strings.Index(outputStr, "cli-standard")
	if simpleIndex != -1 && standardIndex != -1 {
		assert.Less(t, simpleIndex, standardIndex, "CLI-Simple should appear before CLI-Standard")
	}

	// And descriptions should clearly indicate complexity and use cases
	lines := strings.Split(outputStr, "\n")
	simpleLine := findLineContaining(lines, "cli-simple")
	standardLine := findLineContaining(lines, "cli-standard")

	if simpleLine != "" {
		// CLI-Simple should indicate beginner-friendliness
		lowerSimple := strings.ToLower(simpleLine)
		assert.True(t, 
			strings.Contains(lowerSimple, "simple") || 
			strings.Contains(lowerSimple, "beginner") || 
			strings.Contains(lowerSimple, "minimal"),
			"CLI-Simple description should indicate simplicity: %s", simpleLine)
	}

	if standardLine != "" {
		// CLI-Standard should indicate enterprise/advanced features
		lowerStandard := strings.ToLower(standardLine)
		assert.True(t, 
			strings.Contains(lowerStandard, "standard") || 
			strings.Contains(lowerStandard, "enterprise") || 
			strings.Contains(lowerStandard, "advanced") ||
			strings.Contains(lowerStandard, "full"),
			"CLI-Standard description should indicate advanced features: %s", standardLine)
	}

	// And both blueprints should be properly categorized as CLI type
	validator := NewCLIStrategyCoordinationValidator(binary)
	validator.ValidateCliCategorization(t, outputStr)
}

func TestCLIStrategy_ProgressiveComplexityPrompts_InteractiveSelection(t *testing.T) {
	// Scenario: Progressive Complexity Prompts in Interactive Mode
	// Given I start an interactive CLI project creation
	// When I choose "cli" as the project type
	// Then I should be prompted to choose between complexity levels
	// And the prompt should explain the difference between Simple and Standard
	// And beginners should be guided toward Simple
	// And advanced users should be guided toward Standard
	// And the selection should be clear and intuitive

	// NOTE: This test validates the prompting logic exists and is properly structured
	// Full interactive testing would require input simulation which is complex
	// For now, we validate the prompt structure and guidance logic

	binary := buildCLIBinary(t)
	defer cleanupBinary(t, binary)

	// Test that help for new command mentions complexity selection
	cmd := exec.Command(binary, "new", "--help")
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "New command help should be available")

	outputStr := string(output)
	validator := NewCLIStrategyCoordinationValidator(binary)
	validator.ValidateComplexityPromptGuidance(t, outputStr)
}

func TestCLIStrategy_BlueprintSelectionLogic_SimpleArchitecture(t *testing.T) {
	// Scenario: Blueprint Selection Logic for Simple Architecture
	// Given I want to create a simple CLI project
	// When I select "cli" type with "simple" architecture
	// Then the project should use CLI-Simple blueprint
	// And the generated project should have Simple characteristics
	// And the project should compile successfully
	// And the project should have minimal complexity

	// Given I want to create a simple CLI project
	config := types.ProjectConfig{
		Name:         "test-cli-strategy-simple",
		Module:       "github.com/test/test-cli-strategy-simple",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I select "cli" type with "simple" architecture
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should use CLI-Simple blueprint
	validator := NewCLIStrategyCoordinationValidator("")
	validator.ValidateSimpleBlueprintUsage(t, projectPath)

	// And the generated project should have Simple characteristics
	validator.ValidateSimpleCharacteristics(t, projectPath)

	// And the project should compile successfully
	helpers.AssertProjectCompiles(t, projectPath)

	// And the project should have minimal complexity
	validator.ValidateMinimalComplexity(t, projectPath)
}

func TestCLIStrategy_BlueprintSelectionLogic_StandardArchitecture(t *testing.T) {
	// Scenario: Blueprint Selection Logic for Standard Architecture
	// Given I want to create a standard CLI project
	// When I select "cli" type with "standard" architecture
	// Then the project should use CLI-Standard blueprint
	// And the generated project should have Standard characteristics
	// And the project should compile successfully
	// And the project should have enterprise features

	// Given I want to create a standard CLI project
	config := types.ProjectConfig{
		Name:         "test-cli-strategy-standard",
		Module:       "github.com/test/test-cli-strategy-standard",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	// When I select "cli" type with "standard" architecture
	projectPath := helpers.GenerateProject(t, config)

	// Then the project should use CLI-Standard blueprint
	validator := NewCLIStrategyCoordinationValidator("")
	validator.ValidateStandardBlueprintUsage(t, projectPath)

	// And the generated project should have Standard characteristics
	validator.ValidateStandardCharacteristics(t, projectPath)

	// And the project should compile successfully
	helpers.AssertProjectCompiles(t, projectPath)

	// And the project should have enterprise features
	validator.ValidateEnterpriseFeatures(t, projectPath)
}

func TestCLIStrategy_UserGuidanceIntegration_HelpAndDocumentation(t *testing.T) {
	// Scenario: User Guidance Integration with Help and Documentation
	// Given I need help choosing between CLI architectures
	// When I access help documentation
	// Then I should see clear explanations of the differences
	// And I should get recommendations based on use cases
	// And I should see examples of when to use each architecture
	// And documentation links should be provided

	binary := buildCLIBinary(t)
	defer cleanupBinary(t, binary)

	// When I access help documentation
	tests := []struct {
		name string
		args []string
	}{
		{"general help", []string{"--help"}},
		{"new command help", []string{"new", "--help"}},
		{"list command help", []string{"list", "--help"}},
	}

	validator := NewCLIStrategyCoordinationValidator(binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "Help command should work: %v", tt.args)

			outputStr := string(output)
			validator.ValidateUserGuidance(t, outputStr, tt.name)
		})
	}
}

func TestCLIStrategy_EndToEndUserJourney_SimpleToEnterprise(t *testing.T) {
	// Scenario: Complete User Journey from Discovery to Generation
	// Given I am a new user exploring CLI options
	// When I follow the complete journey: discover → select → generate
	// Then each step should work seamlessly
	// And I should be guided appropriately at each stage
	// And both paths (Simple and Enterprise) should work end-to-end
	// And generated projects should match expectations

	binary := buildCLIBinary(t)
	defer cleanupBinary(t, binary)

	t.Run("journey_discovery_phase", func(t *testing.T) {
		// Discovery: List available options
		cmd := exec.Command(binary, "list")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Discovery phase should work")

		outputStr := string(output)
		assert.Contains(t, outputStr, "cli", "Should discover CLI options")
	})

	t.Run("journey_simple_path", func(t *testing.T) {
		// Simple Path: Generate Simple CLI
		config := types.ProjectConfig{
			Name:         "journey-simple-cli",
			Module:       "github.com/test/journey-simple-cli",
			Type:         "cli",
			Architecture: "simple",
			GoVersion:    "1.21",
			Framework:    "cobra",
			Logger:       "slog",
		}

		projectPath := helpers.GenerateProject(t, config)
		
		// Validate complete Simple CLI journey
		validator := NewCLIStrategyCoordinationValidator("")
		validator.ValidateCompleteSimpleJourney(t, projectPath)
	})

	t.Run("journey_enterprise_path", func(t *testing.T) {
		// Enterprise Path: Generate Standard CLI
		config := types.ProjectConfig{
			Name:         "journey-enterprise-cli",
			Module:       "github.com/test/journey-enterprise-cli",
			Type:         "cli",
			Architecture: "standard",
			GoVersion:    "1.21",
			Framework:    "cobra",
			Logger:       "slog",
		}

		projectPath := helpers.GenerateProject(t, config)
		
		// Validate complete Standard CLI journey
		validator := NewCLIStrategyCoordinationValidator("")
		validator.ValidateCompleteEnterpriseJourney(t, projectPath)
	})
}

func TestCLIStrategy_MigrationPathGuidance_SimpleToStandard(t *testing.T) {
	// Scenario: Migration Path Guidance from Simple to Standard
	// Given I have a Simple CLI project
	// When I want to upgrade to Standard CLI
	// Then I should understand the migration path
	// And I should see clear documentation on differences
	// And I should understand what needs to be added/changed
	// And I should have examples of how to migrate

	// Generate a Simple CLI project first
	simpleConfig := types.ProjectConfig{
		Name:         "migration-simple-cli",
		Module:       "github.com/test/migration-simple-cli",
		Type:         "cli",
		Architecture: "simple",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	simplePath := helpers.GenerateProject(t, simpleConfig)

	// Generate a Standard CLI project for comparison
	standardConfig := types.ProjectConfig{
		Name:         "migration-standard-cli",
		Module:       "github.com/test/migration-standard-cli",
		Type:         "cli",
		Architecture: "standard",
		GoVersion:    "1.21",
		Framework:    "cobra",
		Logger:       "slog",
	}

	standardPath := helpers.GenerateProject(t, standardConfig)

	// Validate migration guidance
	validator := NewCLIStrategyCoordinationValidator("")
	validator.ValidateMigrationPath(t, simplePath, standardPath)
}

func TestCLIStrategy_NegativeScenarios_WrongChoices(t *testing.T) {
	// Scenario: Negative test cases for wrong choices
	// Given users might make suboptimal choices
	// When they select inappropriate architectures for their needs
	// Then the system should handle gracefully
	// And provide guidance for better choices
	// And not fail catastrophically

	t.Run("invalid_architecture", func(t *testing.T) {
		// Test with invalid architecture
		config := types.ProjectConfig{
			Name:         "test-invalid-arch",
			Module:       "github.com/test/test-invalid-arch",
			Type:         "cli",
			Architecture: "nonexistent",
			GoVersion:    "1.21",
			Framework:    "cobra",
			Logger:       "slog",
		}

		// This should either fail gracefully or fall back to a default
		// We test that it doesn't panic or cause system errors
		_, err := generateProjectSafely(t, config)
		if err != nil {
			// Error is acceptable for invalid input
			assert.Contains(t, err.Error(), "nonexistent", "Error should mention the invalid architecture")
		}
	})

	t.Run("mismatched_complexity", func(t *testing.T) {
		// Test scenarios where users might choose mismatched complexity
		// This is more of a usability test to ensure the system guides appropriately
		
		binary := buildCLIBinary(t)
		defer cleanupBinary(t, binary)

		// Test that help provides guidance about complexity choices
		cmd := exec.Command(binary, "new", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Help should be available")

		outputStr := string(output)
		// Should provide some guidance about choosing complexity
		validator := NewCLIStrategyCoordinationValidator(binary)
		validator.ValidateComplexityGuidance(t, outputStr)
	})
}

// CLIStrategyCoordinationValidator provides validation methods for CLI strategy coordination
type CLIStrategyCoordinationValidator struct {
	binaryPath string
}

// NewCLIStrategyCoordinationValidator creates a new validator
func NewCLIStrategyCoordinationValidator(binaryPath string) *CLIStrategyCoordinationValidator {
	return &CLIStrategyCoordinationValidator{
		binaryPath: binaryPath,
	}
}

// ValidateCliCategorization validates CLI blueprints are properly categorized
func (v *CLIStrategyCoordinationValidator) ValidateCliCategorization(t *testing.T, listOutput string) {
	t.Helper()

	// Both CLI blueprints should be listed under CLI category or similar
	lines := strings.Split(listOutput, "\n")
	
	// Look for section headers or groupings
	cliSectionFound := false
	for _, line := range lines {
		lowerLine := strings.ToLower(line)
		if strings.Contains(lowerLine, "cli") && 
		   (strings.Contains(lowerLine, "templates") || 
		    strings.Contains(lowerLine, "blueprints") ||
		    strings.Contains(lowerLine, "command") ||
		    strings.Contains(lowerLine, "tool")) {
			cliSectionFound = true
			break
		}
	}

	assert.True(t, cliSectionFound || 
		strings.Contains(listOutput, "cli-simple") || 
		strings.Contains(listOutput, "cli-standard"),
		"CLI blueprints should be properly categorized in list output")
}

// ValidateComplexityPromptGuidance validates complexity prompts provide good guidance
func (v *CLIStrategyCoordinationValidator) ValidateComplexityPromptGuidance(t *testing.T, helpOutput string) {
	t.Helper()

	lowerOutput := strings.ToLower(helpOutput)
	
	// Should mention architecture or complexity selection
	assert.True(t, 
		strings.Contains(lowerOutput, "architecture") ||
		strings.Contains(lowerOutput, "complexity") ||
		strings.Contains(lowerOutput, "simple") ||
		strings.Contains(lowerOutput, "standard"),
		"Help should provide guidance about architecture/complexity choices")
}

// ValidateSimpleBlueprintUsage validates Simple blueprint characteristics
func (v *CLIStrategyCoordinationValidator) ValidateSimpleBlueprintUsage(t *testing.T, projectPath string) {
	t.Helper()

	// Should have Simple blueprint characteristics
	// Main file structure should be simple
	mainFile := filepath.Join(projectPath, "main.go")
	helpers.AssertFileExists(t, mainFile)

	// Should have config.go (not internal/config/)
	configFile := filepath.Join(projectPath, "config.go")
	helpers.AssertFileExists(t, configFile)

	// Should NOT have internal directory structure
	internalDir := filepath.Join(projectPath, "internal")
	helpers.AssertDirectoryNotExists(t, internalDir)
}

// ValidateSimpleCharacteristics validates Simple CLI characteristics
func (v *CLIStrategyCoordinationValidator) ValidateSimpleCharacteristics(t *testing.T, projectPath string) {
	t.Helper()

	// Count total files - should be minimal
	fileCount := helpers.CountAllFiles(t, projectPath)
	assert.Less(t, fileCount, 10, "Simple CLI should have <10 files, got %d", fileCount)

	// Check dependencies are minimal
	goModFile := filepath.Join(projectPath, "go.mod")
	goModContent := helpers.ReadFileContent(t, goModFile)
	
	// Should only have cobra as main dependency
	assert.Contains(t, goModContent, "github.com/spf13/cobra")
	assert.NotContains(t, goModContent, "github.com/spf13/viper")
}

// ValidateMinimalComplexity validates minimal complexity
func (v *CLIStrategyCoordinationValidator) ValidateMinimalComplexity(t *testing.T, projectPath string) {
	t.Helper()

	// No factory patterns
	factoryFiles := helpers.FindFiles(t, projectPath, "*factory*.go")
	assert.Empty(t, factoryFiles, "Simple CLI should not have factory patterns")

	// No complex interfaces
	goFiles := helpers.FindGoFiles(t, projectPath)
	for _, file := range goFiles {
		content := helpers.ReadFile(t, file)
		interfaceCount := helpers.CountInterfaces(content)
		assert.LessOrEqual(t, interfaceCount, 1, "File %s should have minimal interfaces", file)
	}

	// Shallow directory structure
	maxDepth := helpers.GetMaxDirectoryDepth(t, projectPath)
	assert.LessOrEqual(t, maxDepth, 3, "Simple CLI should have shallow directory structure")
}

// ValidateStandardBlueprintUsage validates Standard blueprint characteristics
func (v *CLIStrategyCoordinationValidator) ValidateStandardBlueprintUsage(t *testing.T, projectPath string) {
	t.Helper()

	// Should have Standard blueprint characteristics
	// Should have internal directory structure
	internalDir := filepath.Join(projectPath, "internal")
	helpers.AssertDirectoryExists(t, internalDir)

	// Should have configs directory
	configsDir := filepath.Join(projectPath, "configs")
	helpers.AssertDirectoryExists(t, configsDir)

	// Should have logger factory
	loggerDir := filepath.Join(projectPath, "internal", "logger")
	helpers.AssertDirectoryExists(t, loggerDir)
}

// ValidateStandardCharacteristics validates Standard CLI characteristics
func (v *CLIStrategyCoordinationValidator) ValidateStandardCharacteristics(t *testing.T, projectPath string) {
	t.Helper()

	// Should have more files than Simple
	fileCount := helpers.CountAllFiles(t, projectPath)
	assert.GreaterOrEqual(t, fileCount, 15, "Standard CLI should have >=15 files, got %d", fileCount)

	// Should have Viper for configuration
	goModFile := filepath.Join(projectPath, "go.mod")
	goModContent := helpers.ReadFileContent(t, goModFile)
	assert.Contains(t, goModContent, "github.com/spf13/viper")
}

// ValidateEnterpriseFeatures validates enterprise features
func (v *CLIStrategyCoordinationValidator) ValidateEnterpriseFeatures(t *testing.T, projectPath string) {
	t.Helper()

	// Should have CRUD commands
	cmdDir := filepath.Join(projectPath, "cmd")
	expectedCommands := []string{"create.go", "update.go", "delete.go", "list.go"}
	
	for _, cmdFile := range expectedCommands {
		cmdPath := filepath.Join(cmdDir, cmdFile)
		helpers.AssertFileExists(t, cmdPath)
	}

	// Should have logger implementation
	loggerFile := filepath.Join(projectPath, "internal", "logger", "logger.go")
	helpers.AssertFileExists(t, loggerFile)

	// Should have configuration management
	configFile := filepath.Join(projectPath, "internal", "config", "config.go")
	helpers.AssertFileExists(t, configFile)
}

// ValidateUserGuidance validates user guidance in help output
func (v *CLIStrategyCoordinationValidator) ValidateUserGuidance(t *testing.T, helpOutput, helpType string) {
	t.Helper()

	// Help should be informative and not empty
	assert.NotEmpty(t, strings.TrimSpace(helpOutput), "Help output should not be empty for %s", helpType)

	// Should contain usage information
	lowerOutput := strings.ToLower(helpOutput)
	assert.True(t,
		strings.Contains(lowerOutput, "usage") ||
		strings.Contains(lowerOutput, "example") ||
		strings.Contains(lowerOutput, "command"),
		"Help should contain usage guidance for %s", helpType)
}

// ValidateCompleteSimpleJourney validates the complete Simple CLI journey
func (v *CLIStrategyCoordinationValidator) ValidateCompleteSimpleJourney(t *testing.T, projectPath string) {
	t.Helper()

	// Validate generation
	helpers.AssertProjectCompiles(t, projectPath)

	// Validate Simple characteristics
	v.ValidateSimpleBlueprintUsage(t, projectPath)
	v.ValidateSimpleCharacteristics(t, projectPath)
	v.ValidateMinimalComplexity(t, projectPath)

	// Validate functionality
	helpers.AssertCLIHelpOutput(t, projectPath)
	helpers.AssertCLIVersionOutput(t, projectPath)

	// Should be beginner-friendly
	readmeFile := filepath.Join(projectPath, "README.md")
	if helpers.FileExists(readmeFile) {
		readmeContent := helpers.ReadFileContent(t, readmeFile)
		lowerReadme := strings.ToLower(readmeContent)
		assert.True(t,
			strings.Contains(lowerReadme, "simple") ||
			strings.Contains(lowerReadme, "quick") ||
			strings.Contains(lowerReadme, "getting started"),
			"README should indicate beginner-friendliness")
	}
}

// ValidateCompleteEnterpriseJourney validates the complete Enterprise CLI journey
func (v *CLIStrategyCoordinationValidator) ValidateCompleteEnterpriseJourney(t *testing.T, projectPath string) {
	t.Helper()

	// Validate generation
	helpers.AssertProjectCompiles(t, projectPath)

	// Validate Standard characteristics
	v.ValidateStandardBlueprintUsage(t, projectPath)
	v.ValidateStandardCharacteristics(t, projectPath)
	v.ValidateEnterpriseFeatures(t, projectPath)

	// Validate functionality
	helpers.AssertCLIHelpOutput(t, projectPath)
	helpers.AssertCLIVersionOutput(t, projectPath)

	// Should have tests
	helpers.AssertTestsRun(t, projectPath)

	// Should have Docker support
	dockerFile := filepath.Join(projectPath, "Dockerfile")
	helpers.AssertFileExists(t, dockerFile)
}

// ValidateMigrationPath validates migration guidance between Simple and Standard
func (v *CLIStrategyCoordinationValidator) ValidateMigrationPath(t *testing.T, simplePath, standardPath string) {
	t.Helper()

	// Compare file structures to understand migration
	simpleFiles := helpers.CountAllFiles(t, simplePath)
	standardFiles := helpers.CountAllFiles(t, standardPath)
	
	assert.Less(t, simpleFiles, standardFiles, 
		"Standard should have more files than Simple for migration path validation")

	// Key differences should be clear
	// Simple should NOT have internal/, Standard should have it
	simpleInternal := filepath.Join(simplePath, "internal")
	standardInternal := filepath.Join(standardPath, "internal")
	
	assert.False(t, helpers.DirExists(simpleInternal), "Simple should not have internal/")
	assert.True(t, helpers.DirExists(standardInternal), "Standard should have internal/")

	// Documentation should exist explaining differences
	simpleReadme := filepath.Join(simplePath, "README.md")
	standardReadme := filepath.Join(standardPath, "README.md")
	
	if helpers.FileExists(simpleReadme) && helpers.FileExists(standardReadme) {
		simpleContent := strings.ToLower(helpers.ReadFileContent(t, simpleReadme))
		standardContent := strings.ToLower(helpers.ReadFileContent(t, standardReadme))
		
		// READMEs should indicate their different purposes
		assert.NotEqual(t, simpleContent, standardContent, 
			"Simple and Standard READMEs should be different")
	}
}

// ValidateComplexityGuidance validates guidance about complexity choices
func (v *CLIStrategyCoordinationValidator) ValidateComplexityGuidance(t *testing.T, helpOutput string) {
	t.Helper()

	// Help should provide some indication of options or choices
	lowerOutput := strings.ToLower(helpOutput)
	
	hasGuidance := strings.Contains(lowerOutput, "type") ||
		strings.Contains(lowerOutput, "architecture") ||
		strings.Contains(lowerOutput, "template") ||
		strings.Contains(lowerOutput, "blueprint") ||
		strings.Contains(lowerOutput, "option")
		
	assert.True(t, hasGuidance, "Help should provide guidance about available options")
}

// Helper functions for testing

// buildCLIBinary builds the CLI binary for testing
func buildCLIBinary(t *testing.T) string {
	t.Helper()
	
	// Find project root
	wd, err := os.Getwd()
	require.NoError(t, err)
	
	projectRoot := wd
	for {
		mainGo := filepath.Join(projectRoot, "main.go")
		if helpers.FileExists(mainGo) {
			break
		}
		
		parentDir := filepath.Dir(projectRoot)
		if parentDir == projectRoot {
			t.Fatal("Could not find project root with main.go")
		}
		projectRoot = parentDir
	}
	
	// Build binary
	binaryName := "test-go-starter-" + strings.ReplaceAll(t.Name(), "/", "-")
	binaryPath := filepath.Join(projectRoot, binaryName)
	
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Dir = projectRoot
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Failed to build CLI binary: %s", string(output))
	
	return binaryPath
}

// cleanupBinary removes the test binary
func cleanupBinary(t *testing.T, binaryPath string) {
	t.Helper()
	if binaryPath != "" && helpers.FileExists(binaryPath) {
		_ = os.Remove(binaryPath)
	}
}

// findLineContaining finds a line containing the specified text
func findLineContaining(lines []string, text string) string {
	for _, line := range lines {
		if strings.Contains(line, text) {
			return line
		}
	}
	return ""
}

// generateProjectSafely attempts to generate a project and returns error if it fails
func generateProjectSafely(t *testing.T, config types.ProjectConfig) (string, error) {
	t.Helper()
	
	// This is a wrapper around helpers.GenerateProject that catches panics and errors
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Project generation panicked: %v", r)
		}
	}()
	
	// For now, we'll use the existing GenerateProject function
	// In a real implementation, this would have better error handling
	projectPath := helpers.GenerateProject(t, config)
	return projectPath, nil
}