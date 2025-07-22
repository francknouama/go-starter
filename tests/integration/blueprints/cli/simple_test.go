package cli_blueprints_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCLISimpleBlueprintIntegration validates the integration aspects of cli-simple blueprint
// This ensures the simplified CLI blueprint generates correctly and produces working applications
func TestCLISimpleBlueprintIntegration(t *testing.T) {
	t.Run("cli_simple_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: cli-simple should be in the list

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/cli)
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool first
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot // Run build from project root
		output, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Logf("Build output: %s", string(output))
			t.Logf("Project root: %s", projectRoot)
			t.Logf("Original dir: %s", originalDir)
		}
		require.NoError(t, err, "Failed to build CLI tool")

		// List blueprints
		listCmd := exec.Command("./go-starter", "list")
		output, err = listCmd.CombinedOutput()
		require.NoError(t, err, "List command should succeed")

		outputStr := string(output)
		assert.Contains(t, outputStr, "cli-simple", "cli-simple blueprint should be listed")
		assert.Contains(t, outputStr, "Simple command-line application", "Should show cli-simple description")
	})

	t.Run("cli_simple_generates_minimal_structure", func(t *testing.T) {
		// GIVEN: User wants a simple CLI with minimal complexity
		// WHEN: User generates a project with cli-simple blueprint
		// THEN: Should generate exactly 8 files with minimal structure

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		require.NoError(t, err, "Failed to build CLI tool: %s", string(output))

		// Generate a cli-simple project
		generateCmd := exec.Command("./go-starter", "new", "test-simple-cli",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/simple-cli",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-simple-cli")

		// Check exact file count (should be 8 files)
		var fileCount int
		err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileCount++
				t.Logf("Found file: %s", strings.TrimPrefix(path, projectDir+"/"))
			}
			return nil
		})
		require.NoError(t, err)
		assert.Equal(t, 8, fileCount, "Should have exactly 8 files")

		// Verify essential files exist
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
			assert.FileExists(t, filepath.Join(projectDir, file), "Essential file %s should exist", file)
		}

		// Verify no over-engineered structures exist
		overEngineeredDirs := []string{
			"internal/logger",
			"internal/config",
			"internal/errors",
			"internal/interactive",
			"internal/output",
			".github",
			"configs",
		}

		for _, dir := range overEngineeredDirs {
			assert.NoDirExists(t, filepath.Join(projectDir, dir), "Over-engineered directory %s should not exist", dir)
		}
	})

	t.Run("cli_simple_generated_project_builds_and_runs", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: User builds and runs the project
		// THEN: It should compile and execute successfully

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-cli",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/cli",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-cli")

		// Build the generated project
		buildGeneratedCmd := exec.Command("go", "build", "-o", "test-cli-binary", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err := buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated project should build successfully: %s", string(output))

		// Run version command
		versionCmd := exec.Command("./test-cli-binary", "version")
		versionCmd.Dir = projectDir
		output, err = versionCmd.CombinedOutput()
		require.NoError(t, err, "Version command should run successfully")

		outputStr := string(output)
		assert.Contains(t, outputStr, "test-cli", "Should show project name")
		assert.Contains(t, outputStr, "version", "Should show version info")

		// Run help command
		helpCmd := exec.Command("./test-cli-binary", "--help")
		helpCmd.Dir = projectDir
		output, err = helpCmd.CombinedOutput()
		require.NoError(t, err, "Help command should run successfully")

		outputStr = string(output)
		assert.Contains(t, outputStr, "Usage:", "Should show usage information")
		assert.Contains(t, outputStr, "Available Commands:", "Should show available commands")
		assert.Contains(t, outputStr, "version", "Should list version command")
	})

	t.Run("cli_simple_has_working_makefile", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: User runs make commands
		// THEN: Common make targets should work

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-make",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/make",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-make")

		// Test make build
		makeBuildCmd := exec.Command("make", "build")
		makeBuildCmd.Dir = projectDir
		output, err := makeBuildCmd.CombinedOutput()
		require.NoError(t, err, "make build should succeed: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "bin", "test-make"), "Binary should be created in bin/")

		// Test make test
		makeTestCmd := exec.Command("make", "test")
		makeTestCmd.Dir = projectDir
		output, err = makeTestCmd.CombinedOutput()
		// It's ok if there are no tests, but the command should run
		if err != nil {
			assert.Contains(t, string(output), "no test files", "Should indicate no test files")
		}

		// Test make help
		makeHelpCmd := exec.Command("make", "help")
		makeHelpCmd.Dir = projectDir
		output, err = makeHelpCmd.CombinedOutput()
		require.NoError(t, err, "make help should succeed")

		outputStr := string(output)
		assert.Contains(t, outputStr, "build", "Should show build target")
		assert.Contains(t, outputStr, "test", "Should show test target")
		assert.Contains(t, outputStr, "run", "Should show run target")
	})

	t.Run("cli_simple_config_is_minimal", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: Examining the configuration approach
		// THEN: It should use simple flags/env vars, not complex config files

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-config",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/config",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-config")

		// Read config.go to verify it's simple
		configContent, err := os.ReadFile(filepath.Join(projectDir, "config.go"))
		require.NoError(t, err)

		configStr := string(configContent)

		// Should use simple configuration approach
		assert.Contains(t, configStr, "type Config struct", "Should have Config struct")
		assert.NotContains(t, configStr, "viper", "Should not use Viper for simple config")
		assert.NotContains(t, configStr, "yaml", "Should not use YAML config files")

		// Should support basic env vars
		assert.Contains(t, configStr, "os.Getenv", "Should support environment variables")
	})
}

// TestCLISimpleVsStandardIntegration validates the integration differences between cli-simple and cli-standard
func TestCLISimpleVsStandardIntegration(t *testing.T) {
	t.Run("file_count_comparison", func(t *testing.T) {
		// GIVEN: Both cli-simple and cli-standard blueprints
		// WHEN: Generating projects with each
		// THEN: cli-simple should have ~70% fewer files than cli-standard

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate cli-simple project
		generateSimpleCmd := exec.Command("./go-starter", "new", "simple-project",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/simple",
			"--no-git")
		_, err = generateSimpleCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate cli-standard project
		generateStandardCmd := exec.Command("./go-starter", "new", "standard-project",
			"--type=cli",
			"--complexity=standard",
			"--module=github.com/test/standard",
			"--no-git")
		_, err = generateStandardCmd.CombinedOutput()
		require.NoError(t, err)

		// Count files in each project
		simpleFiles := countFiles(t, filepath.Join(tmpDir, "simple-project"))
		standardFiles := countFiles(t, filepath.Join(tmpDir, "standard-project"))

		t.Logf("cli-simple files: %d", simpleFiles)
		t.Logf("cli-standard files: %d", standardFiles)

		// Verify significant reduction
		reduction := float64(standardFiles-simpleFiles) / float64(standardFiles) * 100
		assert.Greater(t, reduction, 60.0, "cli-simple should have at least 60%% fewer files")
		assert.Less(t, simpleFiles, 10, "cli-simple should have less than 10 files")
		assert.Greater(t, standardFiles, 20, "cli-standard should have more than 20 files")
	})
}

// Helper function to count files in a directory
func countFiles(t *testing.T, dir string) int {
	var count int
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	require.NoError(t, err)
	return count
}
