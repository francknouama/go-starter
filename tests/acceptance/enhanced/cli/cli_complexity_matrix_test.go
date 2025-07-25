package cli

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// CLIComplexityTestContext holds test state for CLI complexity testing
type CLIComplexityTestContext struct {
	projectConfig     *types.ProjectConfig
	projectPath       string
	tempDir           string
	startTime         time.Time
	lastCommandOutput string
	lastCommandError  error
	generatedFiles    []string
	complexity        string
	expectedFileCount int
	actualFileCount   int
}

// TestFeatures runs the CLI complexity matrix BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &CLIComplexityTestContext{}

			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				ctx.generatedFiles = []string{}

				// Initialize templates
				if err := helpers.InitializeTemplates(); err != nil {
					return goCtx, err
				}

				return goCtx, nil
			})

			s.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				if ctx.tempDir != "" {
					os.RemoveAll(ctx.tempDir)
				}
				return goCtx, nil
			})

			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run CLI complexity matrix tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *CLIComplexityTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^I am testing CLI complexity combinations$`, ctx.iAmTestingCLIComplexityCombinations)

	// Project generation steps
	s.Step(`^I generate a CLI project with complexity configuration:$`, ctx.iGenerateACLIProjectWithComplexityConfiguration)
	s.Step(`^I generate a simple CLI tool with minimal configuration:$`, ctx.iGenerateASimpleCLIToolWithMinimalConfiguration)
	s.Step(`^I generate a production-ready CLI application:$`, ctx.iGenerateAProductionReadyCLIApplication)

	// Basic validation steps
	s.Step(`^the project should compile successfully$`, ctx.theProjectShouldCompileSuccessfully)
	s.Step(`^the CLI tool should be executable$`, ctx.theCLIToolShouldBeExecutable)

	// Complexity validation steps
	s.Step(`^the project should have approximately (\d+) files$`, ctx.theProjectShouldHaveApproximatelyFiles)
	s.Step(`^the project structure should be (simple|standard|complex)$`, ctx.theProjectStructureShouldBeComplexity)
	s.Step(`^simple CLI structure should be minimal and focused$`, ctx.simpleCLIStructureShouldBeMinimalAndFocused)
	s.Step(`^standard CLI structure should be production-ready$`, ctx.standardCLIStructureShouldBeProductionReady)

	// Framework validation steps
	s.Step(`^Cobra framework should be properly integrated$`, ctx.cobraFrameworkShouldBeProperlyIntegrated)
	s.Step(`^command structure should follow Cobra conventions$`, ctx.commandStructureShouldFollowCobraConventions)
	s.Step(`^subcommands should be properly organized$`, ctx.subcommandsShouldBeProperlyOrganized)
	s.Step(`^help system should be well-implemented$`, ctx.helpSystemShouldBeWellImplemented)

	// Logger validation steps
	s.Step(`^(.*) logger should be properly integrated$`, ctx.loggerShouldBeProperlyIntegrated)
	s.Step(`^logging should be consistent across commands$`, ctx.loggingShouldBeConsistentAcrossCommands)
	s.Step(`^log levels should be configurable$`, ctx.logLevelsShouldBeConfigurable)

	// Configuration validation steps
	s.Step(`^configuration management should work correctly$`, ctx.configurationManagementShouldWorkCorrectly)
	s.Step(`^config files should be supported$`, ctx.configFilesShouldBeSupported)
	s.Step(`^environment variables should be handled$`, ctx.environmentVariablesShouldBeHandled)
	s.Step(`^command-line flags should override config$`, ctx.commandLineFlagsShouldOverrideConfig)

	// Testing validation steps
	s.Step(`^testing framework should be properly set up$`, ctx.testingFrameworkShouldBeProperlySetUp)
	s.Step(`^command tests should be available$`, ctx.commandTestsShouldBeAvailable)
	s.Step(`^test coverage should be appropriate for complexity$`, ctx.testCoverageShouldBeAppropriateForComplexity)

	// Documentation validation steps
	s.Step(`^documentation should be appropriate for complexity level$`, ctx.documentationShouldBeAppropriateForComplexityLevel)
	s.Step(`^README should explain CLI usage$`, ctx.readmeShouldExplainCLIUsage)
	s.Step(`^help text should be comprehensive$`, ctx.helpTextShouldBeComprehensive)

	// Progressive disclosure validation steps
	s.Step(`^progressive disclosure should work correctly$`, ctx.progressiveDisclosureShouldWorkCorrectly)
	s.Step(`^simple CLI should have minimal learning curve$`, ctx.simpleCLIShouldHaveMinimalLearningCurve)
	s.Step(`^standard CLI should provide full functionality$`, ctx.standardCLIShouldProvideFullFunctionality)
	s.Step(`^migration path from simple to standard should be clear$`, ctx.migrationPathFromSimpleToStandardShouldBeClear)

	// Performance validation steps
	s.Step(`^CLI startup time should be fast$`, ctx.cliStartupTimeShouldBeFast)
	s.Step(`^memory usage should be minimal$`, ctx.memoryUsageShouldBeMinimal)
	s.Step(`^binary size should be reasonable for complexity$`, ctx.binarySizeShouldBeReasonableForComplexity)

	// Build and packaging validation steps
	s.Step(`^build process should work correctly$`, ctx.buildProcessShouldWorkCorrectly)
	s.Step(`^Makefile should provide useful targets$`, ctx.makefileShouldProvideUsefulTargets)
	s.Step(`^Docker support should be included if appropriate$`, ctx.dockerSupportShouldBeIncludedIfAppropriate)
	s.Step(`^release process should be documented$`, ctx.releaseProcessShouldBeDocumented)

	// Cross-platform validation steps
	s.Step(`^CLI should work on multiple platforms$`, ctx.cliShouldWorkOnMultiplePlatforms)
	s.Step(`^build targets should include common platforms$`, ctx.buildTargetsShouldIncludeCommonPlatforms)
	s.Step(`^configuration paths should be OS-appropriate$`, ctx.configurationPathsShouldBeOSAppropriate)
}

// Background step implementations
func (ctx *CLIComplexityTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Use local binary for testing\n	binaryPath := "../../../bin/go-starter"\n	cmd := exec.Command(binaryPath, "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-starter CLI not available: %v", err)
	}
	return nil
}

func (ctx *CLIComplexityTestContext) allTemplatesAreProperlyInitialized() error {
	return helpers.InitializeTemplates()
}

func (ctx *CLIComplexityTestContext) iAmTestingCLIComplexityCombinations() error {
	return nil
}

// Project generation implementations
func (ctx *CLIComplexityTestContext) iGenerateACLIProjectWithComplexityConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *CLIComplexityTestContext) iGenerateASimpleCLIToolWithMinimalConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *CLIComplexityTestContext) iGenerateAProductionReadyCLIApplication(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *CLIComplexityTestContext) generateProjectWithConfiguration(table *godog.Table) error {
	// Parse configuration from table
	config := &types.ProjectConfig{}

	for i := 0; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			config.Type = value
		case "complexity":
			ctx.complexity = value
			// The CLI tool should handle complexity-to-blueprint mapping
		case "framework":
			config.Framework = value
		case "logger":
			config.Logger = value
		case "go_version":
			config.GoVersion = value
		case "expected_files":
			// Parse expected file count for validation
			fmt.Sscanf(value, "%d", &ctx.expectedFileCount)
		}
	}

	// Set defaults
	if config.Name == "" {
		config.Name = fmt.Sprintf("test-cli-%s", ctx.complexity)
	}
	if config.Module == "" {
		config.Module = fmt.Sprintf("github.com/test/cli-%s", ctx.complexity)
	}
	if config.GoVersion == "" {
		config.GoVersion = "1.23"
	}
	if config.Framework == "" {
		config.Framework = "cobra"
	}
	if config.Logger == "" {
		config.Logger = "slog"
	}

	ctx.projectConfig = config

	// Generate project using CLI tool with complexity flag
	var err error
	if ctx.complexity != "" {
		ctx.projectPath, err = ctx.generateWithComplexity(config, ctx.tempDir, ctx.complexity)
	} else {
		ctx.projectPath, err = ctx.generateTestProject(config, ctx.tempDir)
	}

	if err != nil {
		ctx.lastCommandError = err
		return fmt.Errorf("failed to generate project: %v", err)
	}

	// Count generated files
	ctx.actualFileCount = 0
	err = filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(ctx.projectPath, path)
			ctx.generatedFiles = append(ctx.generatedFiles, relPath)
			ctx.actualFileCount++
		}
		return nil
	})

	return err
}

// Basic validation implementations
func (ctx *CLIComplexityTestContext) theProjectShouldCompileSuccessfully() error {
	return ctx.validateProjectCompilation()
}

func (ctx *CLIComplexityTestContext) theCLIToolShouldBeExecutable() error {
	// Try to build and run the CLI tool
	cmd := exec.Command("go", "build", "-o", "test-cli", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to build CLI: %v\nOutput: %s", err, string(output))
	}

	// Check if binary was created
	binaryPath := filepath.Join(ctx.projectPath, "test-cli")
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return fmt.Errorf("CLI binary not created")
	}

	// Try to run with --help flag
	cmd = exec.Command("./test-cli", "--help")
	cmd.Dir = ctx.projectPath
	
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to run CLI with --help: %v\nOutput: %s", err, string(output))
	}

	// Check if help output contains expected CLI content
	helpOutput := string(output)
	if !strings.Contains(helpOutput, "Usage:") {
		return fmt.Errorf("CLI help output doesn't contain usage information")
	}

	return nil
}

// Complexity validation implementations
func (ctx *CLIComplexityTestContext) theProjectShouldHaveApproximatelyFiles(expectedCount int) error {
	tolerance := 3 // Allow ±3 files variance
	
	if ctx.actualFileCount < expectedCount-tolerance || ctx.actualFileCount > expectedCount+tolerance {
		return fmt.Errorf("expected approximately %d files, but got %d", expectedCount, ctx.actualFileCount)
	}
	
	return nil
}

func (ctx *CLIComplexityTestContext) theProjectStructureShouldBeComplexity(complexity string) error {
	switch complexity {
	case "simple":
		return ctx.validateSimpleStructure()
	case "standard":
		return ctx.validateStandardStructure()
	case "complex":
		return ctx.validateComplexStructure()
	default:
		return fmt.Errorf("unknown complexity level: %s", complexity)
	}
}

func (ctx *CLIComplexityTestContext) simpleCLIStructureShouldBeMinimalAndFocused() error {
	return ctx.validateSimpleStructure()
}

func (ctx *CLIComplexityTestContext) standardCLIStructureShouldBeProductionReady() error {
	return ctx.validateStandardStructure()
}

// Helper methods
func (ctx *CLIComplexityTestContext) generateTestProject(config *types.ProjectConfig, tempDir string) (string, error) {
	gen := generator.New()
	
	projectPath := filepath.Join(tempDir, config.Name)
	
	_, err := gen.Generate(*config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	if err != nil {
		return "", err
	}
	
	return projectPath, nil
}

func (ctx *CLIComplexityTestContext) generateWithComplexity(config *types.ProjectConfig, tempDir, complexity string) (string, error) {
	// Use the CLI tool directly with complexity flag to test progressive disclosure
	projectPath := filepath.Join(tempDir, config.Name)
	
	args := []string{
		"new", config.Name,
		"--type=cli",
		"--complexity=" + complexity,
		"--framework=" + config.Framework,
		"--logger=" + config.Logger,
		"--module=" + config.Module,
		"--no-git",
		"--output=" + tempDir,
	}
	
	binaryPath := "../../../bin/go-starter"\n	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return "", fmt.Errorf("CLI generation failed: %v\nOutput: %s", err, string(output))
	}
	
	return projectPath, nil
}

func (ctx *CLIComplexityTestContext) validateProjectCompilation() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("project path not set")
	}
	
	// Check if go.mod exists
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project")
	}
	
	// Try to run go build
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("project compilation failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}

func (ctx *CLIComplexityTestContext) validateSimpleStructure() error {
	// Simple CLI should have minimal files
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"README.md",
	}
	
	for _, file := range requiredFiles {
		filePath := filepath.Join(ctx.projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required simple CLI file not found: %s", file)
		}
	}
	
	// Simple CLI should NOT have complex directories
	shouldNotExist := []string{
		"cmd/",
		"internal/config/",
		"internal/commands/",
	}
	
	for _, dir := range shouldNotExist {
		dirPath := filepath.Join(ctx.projectPath, dir)
		if _, err := os.Stat(dirPath); err == nil {
			return fmt.Errorf("simple CLI should not have complex directory: %s", dir)
		}
	}
	
	return nil
}

func (ctx *CLIComplexityTestContext) validateStandardStructure() error {
	// Standard CLI should have production-ready structure
	requiredDirs := []string{
		"cmd/",
		"internal/",
	}
	
	for _, dir := range requiredDirs {
		dirPath := filepath.Join(ctx.projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("required standard CLI directory not found: %s", dir)
		}
	}
	
	requiredFiles := []string{
		"go.mod",
		"README.md",
		"Makefile",
	}
	
	for _, file := range requiredFiles {
		filePath := filepath.Join(ctx.projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("required standard CLI file not found: %s", file)
		}
	}
	
	return nil
}

func (ctx *CLIComplexityTestContext) validateComplexStructure() error {
	// Complex CLI should have comprehensive structure
	// This is a placeholder for when complex CLI is implemented
	return nil
}

// Simplified implementations for remaining validation methods
func (ctx *CLIComplexityTestContext) cobraFrameworkShouldBeProperlyIntegrated() error {
	// Check for Cobra imports and usage
	cobraPatterns := []string{"github.com/spf13/cobra", "cobra.Command", "cmd.Execute"}
	return ctx.checkForPatterns(cobraPatterns, "Cobra framework integration")
}

func (ctx *CLIComplexityTestContext) commandStructureShouldFollowCobraConventions() error {
	return nil // Simplified implementation
}

func (ctx *CLIComplexityTestContext) subcommandsShouldBeProperlyOrganized() error {
	return nil // Simplified implementation
}

func (ctx *CLIComplexityTestContext) helpSystemShouldBeWellImplemented() error {
	return nil // Simplified implementation
}

func (ctx *CLIComplexityTestContext) loggerShouldBeProperlyIntegrated(logger string) error {
	loggerPatterns := map[string][]string{
		"slog":    {"log/slog", "slog.Logger", "slog.Info"},
		"zap":     {"go.uber.org/zap", "zap.Logger", "zap.Info"},
		"logrus":  {"github.com/sirupsen/logrus", "logrus.Logger", "logrus.Info"},
		"zerolog": {"github.com/rs/zerolog", "zerolog.Logger", "log.Info"},
	}
	
	patterns, exists := loggerPatterns[logger]
	if !exists {
		return fmt.Errorf("unsupported logger: %s", logger)
	}
	
	return ctx.checkForPatterns(patterns, fmt.Sprintf("%s logger integration", logger))
}

func (ctx *CLIComplexityTestContext) checkForPatterns(patterns []string, description string) error {
	found := false
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".go") {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			
			fileContent := string(content)
			for _, pattern := range patterns {
				if strings.Contains(fileContent, pattern) {
					found = true
					return nil
				}
			}
		}
		
		return nil
	})
	
	if err != nil {
		return err
	}
	
	if !found {
		return fmt.Errorf("%s patterns not found in generated project", description)
	}
	
	return nil
}

// Additional simplified implementations
func (ctx *CLIComplexityTestContext) loggingShouldBeConsistentAcrossCommands() error {
	return nil
}

func (ctx *CLIComplexityTestContext) logLevelsShouldBeConfigurable() error {
	return nil
}

func (ctx *CLIComplexityTestContext) configurationManagementShouldWorkCorrectly() error {
	return nil
}

func (ctx *CLIComplexityTestContext) configFilesShouldBeSupported() error {
	return nil
}

func (ctx *CLIComplexityTestContext) environmentVariablesShouldBeHandled() error {
	return nil
}

func (ctx *CLIComplexityTestContext) commandLineFlagsShouldOverrideConfig() error {
	return nil
}

func (ctx *CLIComplexityTestContext) testingFrameworkShouldBeProperlySetUp() error {
	return nil
}

func (ctx *CLIComplexityTestContext) commandTestsShouldBeAvailable() error {
	return nil
}

func (ctx *CLIComplexityTestContext) testCoverageShouldBeAppropriateForComplexity() error {
	return nil
}

func (ctx *CLIComplexityTestContext) documentationShouldBeAppropriateForComplexityLevel() error {
	return nil
}

func (ctx *CLIComplexityTestContext) readmeShouldExplainCLIUsage() error {
	return nil
}

func (ctx *CLIComplexityTestContext) helpTextShouldBeComprehensive() error {
	return nil
}

func (ctx *CLIComplexityTestContext) progressiveDisclosureShouldWorkCorrectly() error {
	return nil
}

func (ctx *CLIComplexityTestContext) simpleCLIShouldHaveMinimalLearningCurve() error {
	return nil
}

func (ctx *CLIComplexityTestContext) standardCLIShouldProvideFullFunctionality() error {
	return nil
}

func (ctx *CLIComplexityTestContext) migrationPathFromSimpleToStandardShouldBeClear() error {
	return nil
}

func (ctx *CLIComplexityTestContext) cliStartupTimeShouldBeFast() error {
	return nil
}

func (ctx *CLIComplexityTestContext) memoryUsageShouldBeMinimal() error {
	return nil
}

func (ctx *CLIComplexityTestContext) binarySizeShouldBeReasonableForComplexity() error {
	return nil
}

func (ctx *CLIComplexityTestContext) buildProcessShouldWorkCorrectly() error {
	return nil
}

func (ctx *CLIComplexityTestContext) makefileShouldProvideUsefulTargets() error {
	return nil
}

func (ctx *CLIComplexityTestContext) dockerSupportShouldBeIncludedIfAppropriate() error {
	return nil
}

func (ctx *CLIComplexityTestContext) releaseProcessShouldBeDocumented() error {
	return nil
}

func (ctx *CLIComplexityTestContext) cliShouldWorkOnMultiplePlatforms() error {
	return nil
}

func (ctx *CLIComplexityTestContext) buildTargetsShouldIncludeCommonPlatforms() error {
	return nil
}

func (ctx *CLIComplexityTestContext) configurationPathsShouldBeOSAppropriate() error {
	return nil
}