package cli

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// CLIAcceptanceTestSuite provides comprehensive ATDD coverage for CLI blueprints
// Tests both simple and standard tiers with command execution validation
type CLIAcceptanceTestSuite struct {
	workingDir    string
	projectDir    string
	projectName   string
	originalDir   string
	projectRoot   string
	httpClient    *http.Client
	
	// CLI specific fields
	tier           string
	complexity     string
	logger         string
	fileCount      int
	
	// Execution tracking
	buildTime      time.Duration
	execTime       time.Duration
}

func setupCLIAcceptanceTest(t *testing.T, tier string) *CLIAcceptanceTestSuite {
	suite := &CLIAcceptanceTestSuite{
		projectName:   "test-cli-" + tier,
		tier:          tier,
		httpClient:    &http.Client{Timeout: 10 * time.Second},
	}

	// Set complexity based on tier
	if tier == "simple" {
		suite.complexity = "simple"
	} else {
		suite.complexity = "standard"
	}

	var err error
	suite.originalDir, err = os.Getwd()
	require.NoError(t, err)

	// Find project root by looking for go.mod file
	projectRoot := suite.originalDir
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			// Reached filesystem root without finding go.mod
			projectRoot = filepath.Join(suite.originalDir, "..", "..", "..", "..")
			break
		}
		projectRoot = parent
	}
	suite.projectRoot = projectRoot
	
	suite.workingDir, err = os.MkdirTemp("", "cli-acceptance-*")
	require.NoError(t, err)

	err = os.Chdir(suite.workingDir)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = os.Chdir(suite.originalDir)
		_ = os.RemoveAll(suite.workingDir)
	})

	return suite
}

func (suite *CLIAcceptanceTestSuite) buildCLI(t *testing.T) {
	// Use the original directory (where tests were started) to find the binary
	srcBinary := filepath.Join(suite.originalDir, "go-starter")
	
	// If binary doesn't exist, build it in the original directory
	if _, err := os.Stat(srcBinary); os.IsNotExist(err) {
		buildCmd := exec.Command("go", "build", "-ldflags", "-s -w", "-o", "go-starter", ".")
		buildCmd.Dir = suite.originalDir
		output, err := buildCmd.CombinedOutput()
		require.NoError(t, err, "Failed to build go-starter CLI in %s: %s", suite.originalDir, string(output))
	}

	// Copy binary to working directory for test execution
	dstBinary := filepath.Join(suite.workingDir, "go-starter")
	
	data, err := os.ReadFile(srcBinary)
	require.NoError(t, err, "Failed to read built binary from %s", srcBinary)
	
	err = os.WriteFile(dstBinary, data, 0755)
	require.NoError(t, err, "Failed to copy binary to working directory")
}

func (suite *CLIAcceptanceTestSuite) generateCLIProject(t *testing.T, args ...string) {
	suite.buildCLI(t)

	baseArgs := []string{
		"new", suite.projectName,
		"--type=cli",
		"--module=github.com/test/" + suite.projectName,
		"--complexity=" + suite.complexity,
		"--logger=" + getDefault(suite.logger, "slog"),
		"--no-git",
	}

	allArgs := append(baseArgs, args...)
	goStarterPath := filepath.Join(suite.workingDir, "go-starter")
	generateCmd := exec.Command(goStarterPath, allArgs...)
	generateCmd.Dir = suite.workingDir

	output, err := generateCmd.CombinedOutput()
	require.NoError(t, err, "CLI project generation should succeed: %s", string(output))

	suite.projectDir = filepath.Join(suite.workingDir, suite.projectName)
	assert.DirExists(t, suite.projectDir, "CLI project directory should be created")
	
	// Count generated files
	suite.fileCount = suite.countFiles(suite.projectDir)
}

func (suite *CLIAcceptanceTestSuite) countFiles(dir string) int {
	count := 0
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && !strings.HasPrefix(filepath.Base(path), ".") {
			count++
		}
		return nil
	})
	return count
}

func (suite *CLIAcceptanceTestSuite) checkFileExists(t *testing.T, relativePath string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	
	// Check if path exists (file or directory)
	_, err := os.Stat(fullPath)
	assert.NoError(t, err, "Path should exist: %s", relativePath)
}

func (suite *CLIAcceptanceTestSuite) checkFileContains(t *testing.T, relativePath, content string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	
	// Handle directory searches
	if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
		found := false
		_ = filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if strings.HasSuffix(path, ".go") {
				fileContent, err := os.ReadFile(path)
				if err == nil && strings.Contains(string(fileContent), content) {
					found = true
					return filepath.SkipDir
				}
			}
			return nil
		})
		assert.True(t, found, "Content '%s' not found in directory %s", content, relativePath)
		return
	}

	fileContent, err := os.ReadFile(fullPath)
	require.NoError(t, err, "Should be able to read file: %s", relativePath)
	assert.Contains(t, string(fileContent), content, "File %s should contain '%s'", relativePath, content)
}

func (suite *CLIAcceptanceTestSuite) checkFileDoesNotExist(t *testing.T, relativePath string) {
	fullPath := filepath.Join(suite.projectDir, relativePath)
	assert.NoFileExists(t, fullPath, "File should not exist: %s", relativePath)
}

func (suite *CLIAcceptanceTestSuite) compileCLIProject(t *testing.T) {
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = suite.projectDir
	output, err := modCmd.CombinedOutput()
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

	// Build the CLI binary
	buildStart := time.Now()
	buildCmd := exec.Command("go", "build", "-o", suite.projectName, ".")
	buildCmd.Dir = suite.projectDir
	output, err = buildCmd.CombinedOutput()
	suite.buildTime = time.Since(buildStart)
	require.NoError(t, err, "CLI project should compile successfully: %s", string(output))
	
}

func (suite *CLIAcceptanceTestSuite) testCLIExecution(t *testing.T, args ...string) ([]byte, error) {
	cmd := exec.Command("./"+suite.projectName, args...)
	cmd.Dir = suite.projectDir
	
	execStart := time.Now()
	output, err := cmd.CombinedOutput()
	suite.execTime = time.Since(execStart)
	
	return output, err
}

// ATDD Test Scenarios for CLI Blueprints

func TestCLIAcceptance_SimpleTierGeneration(t *testing.T) {
	// GIVEN: A developer wants to create a simple CLI for quick utilities
	// WHEN: They generate a CLI with simple complexity
	// THEN: A minimal, functional CLI application is generated

	suite := setupCLIAcceptanceTest(t, "simple")
	suite.generateCLIProject(t)

	// Verify simple tier file count (11 files, excluding .gitignore)
	expectedFiles := 11
	assert.Equal(t, expectedFiles, suite.fileCount, 
		"Simple CLI should have exactly %d files, got %d", expectedFiles, suite.fileCount)

	// Verify essential simple CLI components
	essentialFiles := []string{
		"main.go",
		"go.mod", 
		"README.md",
		"Makefile",
		"config.go",
		"cmd/root.go",
		"cmd/version.go",
		".gitignore",
	}

	for _, file := range essentialFiles {
		suite.checkFileExists(t, file)
	}

	// Verify simple CLI should NOT have advanced directories
	advancedDirs := []string{
		"internal",
		"configs",
		".github",
	}

	for _, dir := range advancedDirs {
		suite.checkFileDoesNotExist(t, dir)
	}

	// Verify simple CLI characteristics
	suite.checkFileContains(t, "main.go", "log/slog")        // Uses slog by default
	suite.checkFileContains(t, "config.go", "Config")       // Simple config in root
	suite.checkFileContains(t, "go.mod", "github.com/spf13/cobra") // Uses Cobra

	// Verify the CLI compiles and runs
	suite.compileCLIProject(t)
	
	// Test basic CLI execution
	output, err := suite.testCLIExecution(t, "--help")
	require.NoError(t, err, "CLI help command should work")
	assert.Contains(t, string(output), "Usage:", "Help output should contain usage information")
	
	// Test version command
	output, err = suite.testCLIExecution(t, "version")
	require.NoError(t, err, "CLI version command should work")
	assert.Contains(t, string(output), "version", "Version output should contain version information")
	
	// Verify fast startup time for simple CLI
	assert.Less(t, suite.execTime.Milliseconds(), int64(100), 
		"Simple CLI should have fast startup time")
}

func TestCLIAcceptance_StandardTierGeneration(t *testing.T) {
	// GIVEN: A developer wants to create a production-ready CLI application
	// WHEN: They generate a CLI with standard complexity
	// THEN: A comprehensive, production-ready CLI application is generated

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify standard tier file count (~29 files, allow some tolerance)
	expectedMinFiles := 25
	expectedMaxFiles := 35
	assert.GreaterOrEqual(t, suite.fileCount, expectedMinFiles, 
		"Standard CLI should have at least %d files, got %d", expectedMinFiles, suite.fileCount)
	assert.LessOrEqual(t, suite.fileCount, expectedMaxFiles,
		"Standard CLI should have at most %d files, got %d", expectedMaxFiles, suite.fileCount)

	// Verify standard CLI has layered architecture
	essentialDirs := []string{
		"cmd",
		"internal/config",
		"internal/logger",
		"internal/errors", 
		"internal/output",
		"internal/version",
		"configs",
		".github/workflows",
	}

	for _, dir := range essentialDirs {
		suite.checkFileExists(t, dir)
	}

	// Verify advanced features exist
	advancedFiles := []string{
		"configs/config.yaml",
		".github/workflows/ci.yml",
		".github/workflows/release.yml",
		"Dockerfile",
		".env.example",
	}

	for _, file := range advancedFiles {
		suite.checkFileExists(t, file)
	}

	// Verify standard CLI characteristics
	suite.checkFileContains(t, "internal/config", "viper")   // Uses Viper for config
	suite.checkFileContains(t, "internal/logger", "interface") // Logger abstraction
	suite.checkFileContains(t, "internal/interactive", "survey") // Interactive features

	// Verify the CLI compiles and runs
	suite.compileCLIProject(t)
	
	// Test CLI execution with various commands
	testCommands := []struct {
		args     []string
		contains string
	}{
		{[]string{"--help"}, "Usage:"},
		{[]string{"version"}, "version"},
		{[]string{"completion", "--help"}, "completion"},
	}

	for _, test := range testCommands {
		output, err := suite.testCLIExecution(t, test.args...)
		require.NoError(t, err, "CLI command %v should work", test.args)
		assert.Contains(t, string(output), test.contains, 
			"Command %v output should contain '%s'", test.args, test.contains)
	}
}

func TestCLIAcceptance_ProgressiveComplexityComparison(t *testing.T) {
	// GIVEN: A developer wants to understand the difference between CLI tiers
	// WHEN: They compare simple and standard CLI blueprints
	// THEN: The differences should be clear and documented

	// Generate simple CLI
	simpleSuite := setupCLIAcceptanceTest(t, "simple")
	simpleSuite.generateCLIProject(t)
	simpleSuite.compileCLIProject(t)

	// Generate standard CLI in separate directory  
	standardSuite := setupCLIAcceptanceTest(t, "standard")
	standardSuite.generateCLIProject(t)
	standardSuite.compileCLIProject(t)

	// Compare file counts
	assert.Equal(t, 11, simpleSuite.fileCount, "Simple CLI should have 11 files")
	assert.GreaterOrEqual(t, standardSuite.fileCount, 25, "Standard CLI should have 25+ files")
	
	fileReduction := float64(simpleSuite.fileCount) / float64(standardSuite.fileCount)
	assert.Less(t, fileReduction, 0.35, "Simple CLI should be 73%% reduction in files")

	// Verify simple CLI focuses on essentials only
	simpleSuite.checkFileExists(t, "config.go") // Simple config
	simpleSuite.checkFileDoesNotExist(t, "internal") // No internal packages

	// Verify standard CLI has advanced features
	standardSuite.checkFileExists(t, "internal") // Internal packages
	standardSuite.checkFileExists(t, "configs") // Advanced config
	standardSuite.checkFileExists(t, ".github") // CI/CD

	// Compare performance
	simpleOutput, err := simpleSuite.testCLIExecution(t, "--version")
	require.NoError(t, err)
	
	standardOutput, err := standardSuite.testCLIExecution(t, "--version")
	require.NoError(t, err)

	// Both should work
	assert.NotEmpty(t, simpleOutput)
	assert.NotEmpty(t, standardOutput)

	// Simple should be faster
	assert.Less(t, simpleSuite.buildTime, standardSuite.buildTime,
		"Simple CLI should build faster than standard CLI")

	// Verify migration path documentation
	simpleSuite.checkFileContains(t, "README.md", "upgrade")
	standardSuite.checkFileContains(t, "README.md", "migration")
}

func TestCLIAcceptance_MultiLoggerSupport(t *testing.T) {
	// GIVEN: A developer wants to use different logging libraries
	// WHEN: They generate CLI applications with different loggers
	// THEN: Each CLI should use the appropriate logger

	loggers := []struct {
		name       string
		importPath string
		tier       string
	}{
		{"slog", "log/slog", "simple"},
		{"slog", "log/slog", "standard"},
		{"zap", "go.uber.org/zap", "standard"},
		{"logrus", "github.com/sirupsen/logrus", "standard"},
		{"zerolog", "github.com/rs/zerolog", "standard"},
	}

	for _, logger := range loggers {
		t.Run(fmt.Sprintf("%s_%s", logger.name, logger.tier), func(t *testing.T) {
			suite := setupCLIAcceptanceTest(t, logger.tier)
			suite.logger = logger.name
			
			suite.generateCLIProject(t)
			
			// Verify logger import exists
			if logger.tier == "simple" {
				suite.checkFileContains(t, "main.go", logger.importPath)
			} else {
				suite.checkFileContains(t, "internal/logger", logger.importPath)
			}
			
			// Ensure it compiles with the logger dependency
			suite.compileCLIProject(t)
			
			// Test CLI execution
			output, err := suite.testCLIExecution(t, "--help")
			require.NoError(t, err, "CLI with %s logger should work", logger.name)
			assert.Contains(t, string(output), "Usage:", "CLI should display help")
		})
	}
}

func TestCLIAcceptance_CommandLineInterface(t *testing.T) {
	// GIVEN: A developer wants comprehensive CLI functionality
	// WHEN: They generate a standard CLI application
	// THEN: The CLI should support all expected command patterns

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)
	suite.compileCLIProject(t)

	// Test various CLI patterns and flags
	testScenarios := []struct {
		name     string
		args     []string
		expectError bool
		contains []string
	}{
		{
			name:     "help_command",
			args:     []string{"--help"},
			contains: []string{"Usage:", "Available Commands:", "Flags:"},
		},
		{
			name:     "version_command",
			args:     []string{"version"},
			contains: []string{"version"},
		},
		{
			name:     "quiet_flag",
			args:     []string{"--quiet", "version"},
			contains: []string{}, // Should suppress output
		},
		{
			name:     "verbose_flag", 
			args:     []string{"--verbose", "version"},
			contains: []string{"version"},
		},
		{
			name:     "invalid_command",
			args:     []string{"nonexistent"},
			expectError: true,
			contains: []string{"unknown command", "Error:"},
		},
		{
			name:     "completion_help",
			args:     []string{"completion", "--help"},
			contains: []string{"completion", "bash", "zsh"},
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			output, err := suite.testCLIExecution(t, scenario.args...)
			
			if scenario.expectError {
				assert.Error(t, err, "Command %v should fail", scenario.args)
			} else {
				assert.NoError(t, err, "Command %v should succeed", scenario.args)
			}
			
			for _, expectedContent := range scenario.contains {
				assert.Contains(t, string(output), expectedContent,
					"Output should contain '%s' for command %v", expectedContent, scenario.args)
			}
		})
	}
}

func TestCLIAcceptance_ConfigurationManagement(t *testing.T) {
	// GIVEN: A developer wants configurable CLI applications
	// WHEN: They generate a standard CLI with configuration support
	// THEN: The CLI should support various configuration methods

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify configuration infrastructure
	suite.checkFileExists(t, "configs/config.yaml")
	suite.checkFileExists(t, ".env.example")
	suite.checkFileExists(t, "internal/config")

	// Verify configuration features
	suite.checkFileContains(t, "internal/config", "viper")
	suite.checkFileContains(t, "internal/config", "environment")
	suite.checkFileContains(t, "internal/config", "yaml")

	// Compile and test configuration
	suite.compileCLIProject(t)
	
	// Test that CLI works with default configuration
	output, err := suite.testCLIExecution(t, "--help")
	require.NoError(t, err)
	assert.Contains(t, string(output), "config", "CLI should mention configuration")
}

func TestCLIAcceptance_ProductionReadiness(t *testing.T) {
	// GIVEN: A developer wants production-ready CLI tools
	// WHEN: They generate a standard CLI application
	// THEN: All production features should be included

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Docker and containerization
	suite.checkFileExists(t, "Dockerfile")
	suite.checkFileContains(t, "Dockerfile", "FROM golang:")
	suite.checkFileContains(t, "Dockerfile", "ENTRYPOINT")

	// CI/CD pipelines
	suite.checkFileExists(t, ".github/workflows/ci.yml")
	suite.checkFileExists(t, ".github/workflows/release.yml")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "test:")
	suite.checkFileContains(t, ".github/workflows/ci.yml", "lint:")
	suite.checkFileContains(t, ".github/workflows/release.yml", "goreleaser")

	// Development workflow
	suite.checkFileExists(t, "Makefile")
	makefileContent, err := os.ReadFile(filepath.Join(suite.projectDir, "Makefile"))
	require.NoError(t, err)
	content := string(makefileContent)

	essentialTargets := []string{"build:", "test:", "lint:", "clean:", "help:"}
	for _, target := range essentialTargets {
		assert.Contains(t, content, target, "Makefile should include target: %s", target)
	}

	// Environment configuration
	suite.checkFileExists(t, ".env.example")

	// Testing infrastructure
	suite.checkFileExists(t, "cmd")
	suite.checkFileContains(t, "cmd", "_test.go")

	// Verify compilation works
	suite.compileCLIProject(t)
}

func TestCLIAcceptance_DevelopmentWorkflow(t *testing.T) {
	// GIVEN: A developer wants efficient CLI development workflow
	// WHEN: They generate a CLI application
	// THEN: Development tools and workflows should be ready

	suite := setupCLIAcceptanceTest(t, "standard") 
	suite.generateCLIProject(t)

	// Verify development tools
	suite.checkFileExists(t, "Makefile")
	suite.checkFileExists(t, ".github/workflows")
	
	// Test Makefile targets by examining content
	makefileContent, err := os.ReadFile(filepath.Join(suite.projectDir, "Makefile"))
	require.NoError(t, err)
	content := string(makefileContent)

	// Essential development targets
	devTargets := []string{
		"build:",     // Building the CLI
		"test:",      // Running tests
		"lint:",      // Code linting
		"clean:",     // Cleanup
		"install:",   // Installing dependencies
		"dev:",       // Development mode
		"help:",      // Help information
	}

	for _, target := range devTargets {
		assert.Contains(t, content, target, "Makefile should include development target: %s", target)
	}

	// Verify README contains development instructions
	suite.checkFileExists(t, "README.md")
	readmeContent, err := os.ReadFile(filepath.Join(suite.projectDir, "README.md"))
	require.NoError(t, err)
	readme := string(readmeContent)

	devInstructions := []string{
		"installation", "usage", "development", "build", "test",
	}

	for _, instruction := range devInstructions {
		assert.Contains(t, strings.ToLower(readme), instruction, 
			"README should include development instruction: %s", instruction)
	}
}

func TestCLIAcceptance_TestingInfrastructure(t *testing.T) {
	// GIVEN: A developer wants well-tested CLI applications
	// WHEN: They generate a standard CLI with testing support
	// THEN: Complete testing infrastructure should be included

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify test files exist
	testFiles := []string{
		"cmd",        // Command tests
	}

	for _, testPath := range testFiles {
		// Check for test files in the directory
		testDir := filepath.Join(suite.projectDir, testPath)
		found := false
		_ = filepath.Walk(testDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if strings.HasSuffix(info.Name(), "_test.go") {
				found = true
				return filepath.SkipDir
			}
			return nil
		})
		assert.True(t, found, "Should find test files in %s", testPath)
	}

	// Verify testing framework
	suite.checkFileContains(t, "go.mod", "testify")

	// Compile and run tests
	suite.compileCLIProject(t)
	
	// Test that go test works
	testCmd := exec.Command("go", "test", "./cmd/...")
	testCmd.Dir = suite.projectDir
	output, err := testCmd.CombinedOutput()
	
	// Tests should either pass or indicate no tests (but not fail to compile)
	if err != nil {
		// Check if it's just "no tests to run"
		assert.Contains(t, string(output), "no tests", 
			"Test command should either pass or have no tests: %s", string(output))
	}
}

func TestCLIAcceptance_ErrorHandlingAndValidation(t *testing.T) {
	// GIVEN: A developer wants robust CLI applications
	// WHEN: They generate a CLI with error handling
	// THEN: Comprehensive error handling should be implemented

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify error handling infrastructure
	suite.checkFileExists(t, "internal/errors")
	suite.checkFileContains(t, "internal/errors", "Error")
	suite.checkFileContains(t, "internal/errors", "exit")

	// Compile and test error scenarios
	suite.compileCLIProject(t)

	// Test invalid command handling
	output, err := suite.testCLIExecution(t, "invalid-command")
	assert.Error(t, err, "Invalid command should return error")
	assert.Contains(t, string(output), "unknown command", 
		"Should provide helpful error message for invalid commands")

	// Test help for error recovery
	output, err = suite.testCLIExecution(t, "--help")
	require.NoError(t, err, "Help should always work")
	assert.Contains(t, string(output), "Usage:", "Help should show usage")
}

func TestCLIAcceptance_PerformanceAndEfficiency(t *testing.T) {
	// GIVEN: A developer wants high-performance CLI applications
	// WHEN: They generate CLI applications
	// THEN: Performance characteristics should be optimized

	// Test both tiers for performance comparison
	simpleSuite := setupCLIAcceptanceTest(t, "simple")
	simpleSuite.generateCLIProject(t)
	simpleSuite.compileCLIProject(t)

	standardSuite := setupCLIAcceptanceTest(t, "standard") 
	standardSuite.generateCLIProject(t)
	standardSuite.compileCLIProject(t)

	// Performance expectations
	maxBuildTime := 30 * time.Second
	maxExecTime := 200 * time.Millisecond

	// Build time should be reasonable
	assert.Less(t, simpleSuite.buildTime, maxBuildTime, 
		"Simple CLI build time should be under %v", maxBuildTime)
	assert.Less(t, standardSuite.buildTime, maxBuildTime,
		"Standard CLI build time should be under %v", maxBuildTime)

	// Execution should be fast
	_, err := simpleSuite.testCLIExecution(t, "version")
	require.NoError(t, err)
	assert.Less(t, simpleSuite.execTime, maxExecTime,
		"Simple CLI execution should be under %v", maxExecTime)

	_, err = standardSuite.testCLIExecution(t, "version")
	require.NoError(t, err) 
	assert.Less(t, standardSuite.execTime, maxExecTime,
		"Standard CLI execution should be under %v", maxExecTime)

	// Simple should be faster than standard
	assert.Less(t, simpleSuite.buildTime, standardSuite.buildTime,
		"Simple CLI should build faster than standard")
}

func TestCLIAcceptance_CrossPlatformSupport(t *testing.T) {
	// GIVEN: A developer wants CLI applications that work everywhere
	// WHEN: They generate a CLI for cross-platform use
	// THEN: The CLI should be platform-agnostic

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify cross-platform considerations
	suite.checkFileContains(t, "go.mod", "go 1.21") // Modern Go version
	
	// Check for cross-platform build configuration
	if suite.checkFileExists(t, ".github/workflows/ci.yml"); true {
		suite.checkFileContains(t, ".github/workflows/ci.yml", "matrix:")
		suite.checkFileContains(t, ".github/workflows/ci.yml", "os:")
	}

	// Verify CLI compiles (on current platform)
	suite.compileCLIProject(t)
	
	// Test basic functionality
	output, err := suite.testCLIExecution(t, "--help")
	require.NoError(t, err)
	assert.Contains(t, string(output), "Usage:")
}

func TestCLIAcceptance_DocumentationAndUsability(t *testing.T) {
	// GIVEN: A developer wants well-documented CLI applications
	// WHEN: They generate a CLI with comprehensive documentation
	// THEN: Documentation should be complete and helpful

	suite := setupCLIAcceptanceTest(t, "standard")
	suite.generateCLIProject(t)

	// Verify documentation files
	suite.checkFileExists(t, "README.md")
	
	readmeContent, err := os.ReadFile(filepath.Join(suite.projectDir, "README.md"))
	require.NoError(t, err)
	readme := strings.ToLower(string(readmeContent))

	// Essential documentation sections
	docSections := []string{
		"installation", "usage", "examples", "development", 
		"build", "test", "contributing", "license",
	}

	for _, section := range docSections {
		assert.Contains(t, readme, section, 
			"README should include section about: %s", section)
	}

	// Verify built-in help system
	suite.compileCLIProject(t)
	
	output, err := suite.testCLIExecution(t, "--help")
	require.NoError(t, err)
	help := string(output)

	helpElements := []string{
		"Usage:", "Available Commands:", "Flags:", 
		"Global Flags:", "Use", "--help",
	}

	for _, element := range helpElements {
		assert.Contains(t, help, element,
			"Help output should include: %s", element)
	}

	// Test command-specific help
	output, err = suite.testCLIExecution(t, "version", "--help")
	require.NoError(t, err)
	assert.Contains(t, string(output), "version", "Version command should have help")
}

// Helper function to get default values
func getDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

// Meta-test to ensure all CLI acceptance tests can run together
func TestCLIAcceptance_FullSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping full CLI acceptance test suite in short mode")
	}

	// This meta-test ensures all acceptance tests can run together
	t.Run("SimpleTierGeneration", TestCLIAcceptance_SimpleTierGeneration)
	t.Run("StandardTierGeneration", TestCLIAcceptance_StandardTierGeneration)
	t.Run("ProgressiveComplexityComparison", TestCLIAcceptance_ProgressiveComplexityComparison)
	t.Run("MultiLoggerSupport", TestCLIAcceptance_MultiLoggerSupport)
	t.Run("CommandLineInterface", TestCLIAcceptance_CommandLineInterface)
	t.Run("ConfigurationManagement", TestCLIAcceptance_ConfigurationManagement)
	t.Run("ProductionReadiness", TestCLIAcceptance_ProductionReadiness)
	t.Run("TestingInfrastructure", TestCLIAcceptance_TestingInfrastructure)
	t.Run("ErrorHandlingAndValidation", TestCLIAcceptance_ErrorHandlingAndValidation)
	t.Run("PerformanceAndEfficiency", TestCLIAcceptance_PerformanceAndEfficiency)
}