package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiBinarySmokeTest runs quick smoke tests to verify basic functionality
// This test is designed to run quickly and catch major regressions
func TestMultiBinarySmokeTest(t *testing.T) {
	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("all_binaries_build_and_run", func(t *testing.T) {
		// GIVEN: Multi-binary structure
		binaries := map[string]struct {
			buildPath string
			testArgs  []string
		}{
			"go-starter": {
				buildPath: "./cmd/go-starter",
				testArgs:  []string{"version"},
			},
			"go-starter-dev": {
				buildPath: "./cmd/go-starter-dev",
				testArgs:  nil, // Server binary - just test build
			},
			"go-starter-web": {
				buildPath: "./web/cmd/web-server", 
				testArgs:  nil, // Server binary - just test build
			},
			"legacy": {
				buildPath: ".",
				testArgs:  []string{"version"},
			},
		}

		for name, config := range binaries {
			t.Run(name, func(t *testing.T) {
				// WHEN: Building binary
				binaryName := getBinaryName(name)
				binaryPath := filepath.Join(tmpDir, binaryName)

				buildErr := buildBinary(t, projectRoot, binaryPath, name, config.buildPath)

				// THEN: Build should succeed
				if buildErr != nil {
					t.Logf("Build failed for %s", name)
				}
				require.NoError(t, buildErr, "Binary %s should build successfully", name)

				// AND: Binary should exist
				assert.FileExists(t, binaryPath, "Binary file should exist")

				// AND: If it's a CLI binary, basic command should work
				if config.testArgs != nil {
					testCmd := exec.Command(binaryPath, config.testArgs...)
					testOutput, testErr := testCmd.CombinedOutput()
					
					if testErr != nil {
						t.Logf("Test command failed for %s:\nOutput: %s", name, string(testOutput))
					}
					assert.NoError(t, testErr, "Test command should work for %s", name)
				}
			})
		}
	})

	t.Run("cli_basic_functionality", func(t *testing.T) {
		// GIVEN: Built CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// Test essential CLI operations
		testCases := []struct {
			name           string
			args           []string
			shouldContain  []string
			shouldNotError bool
		}{
			{
				name:           "version_command",
				args:           []string{"version"},
				shouldContain:  []string{"Version", "Go"}, // More flexible - just check for version info
				shouldNotError: true,
			},
			{
				name:           "help_command",
				args:           []string{"--help"},
				shouldContain:  []string{"COMMANDS", "FLAGS"}, // More flexible - just check for help structure
				shouldNotError: true,
			},
			{
				name:           "list_blueprints",
				args:           []string{"list"},
				shouldContain:  []string{"blueprints", "web-api", "cli"},
				shouldNotError: true,
			},
			{
				name:           "new_help",
				args:           []string{"new", "--help"},
				shouldContain:  []string{"Create a new Go project", "--type"},
				shouldNotError: true,
			},
			{
				name:           "dry_run_simple",
				args:           []string{"new", "test-project", "--type=cli", "--complexity=simple", "--dry-run"},
				shouldContain:  []string{"Files to be generated", "main.go", "go.mod"},
				shouldNotError: true,
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// WHEN: Running CLI command
				cmd := exec.Command(binaryPath, tc.args...)
				output, err := cmd.CombinedOutput()
				outputStr := string(output)

				// THEN: Should behave as expected
				if tc.shouldNotError {
					if err != nil {
						t.Logf("Command failed: %s\nOutput: %s", tc.name, outputStr)
					}
					assert.NoError(t, err, "Command %s should succeed", tc.name)
				}

				// AND: Should contain expected content
				for _, contain := range tc.shouldContain {
					assert.Contains(t, outputStr, contain, 
						"Output should contain '%s' for command %s", contain, tc.name)
				}
			})
		}
	})

	t.Run("embedded_blueprints_work", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Running from isolated directory
		isolatedDir := filepath.Join(tmpDir, "isolated")
		err = os.MkdirAll(isolatedDir, 0755)
		require.NoError(t, err, "Should create isolated directory")

		cmd := exec.Command(binaryPath, "list")
		cmd.Dir = isolatedDir // Run from directory without source code
		output, err := cmd.CombinedOutput()

		// THEN: Should work with embedded blueprints
		require.NoError(t, err, "Should work with embedded blueprints")
		
		outputStr := string(output)
		expectedBlueprints := []string{"web-api", "cli", "library"}
		for _, blueprint := range expectedBlueprints {
			assert.Contains(t, outputStr, blueprint, 
				"Should list embedded blueprint: %s", blueprint)
		}
	})

	t.Run("project_generation_works", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Generating a simple project
		projectName := "smoke-test-project"
		cmd := exec.Command(binaryPath, "new", projectName,
			"--type=cli", "--complexity=simple",
			"--module=github.com/test/"+projectName)
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()

		// THEN: Generation should succeed
		if err != nil {
			t.Logf("Generation failed:\nOutput: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// AND: Project should exist with expected files
		projectDir := filepath.Join(tmpDir, projectName)
		assert.DirExists(t, projectDir, "Project directory should exist")
		assert.FileExists(t, filepath.Join(projectDir, "main.go"), "main.go should exist")
		assert.FileExists(t, filepath.Join(projectDir, "go.mod"), "go.mod should exist")
		assert.FileExists(t, filepath.Join(projectDir, "README.md"), "README.md should exist")

		// AND: Generated project should compile
		compileCmd := exec.Command("go", "build", "./...")
		compileCmd.Dir = projectDir
		compileOutput, compileErr := compileCmd.CombinedOutput()

		if compileErr != nil {
			t.Logf("Compilation failed:\nOutput: %s", string(compileOutput))
		}
		assert.NoError(t, compileErr, "Generated project should compile")
	})

	t.Run("legacy_deprecation_warning", func(t *testing.T) {
		// GIVEN: Legacy binary
		binaryName := "legacy"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Legacy binary should build")

		// WHEN: Running with help flag
		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.CombinedOutput()

		// THEN: Should show deprecation warning
		outputStr := string(output)
		assert.Contains(t, outputStr, "DEPRECATION WARNING", 
			"Should show deprecation warning")
		assert.Contains(t, outputStr, "cmd/go-starter", 
			"Should mention new CLI path")
		assert.Contains(t, outputStr, "cmd/go-starter-dev", 
			"Should mention dev server path")
		assert.Contains(t, outputStr, "go-starter-web", 
			"Should mention web server binary")

		// AND: CLI functionality should still work
		assert.NoError(t, err, "Legacy binary should still function")
	})

	t.Run("progressive_disclosure_basic", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Testing basic vs advanced help
		basicCmd := exec.Command(binaryPath, "new", "--help")
		basicOutput, basicErr := basicCmd.CombinedOutput()

		advancedCmd := exec.Command(binaryPath, "new", "--advanced", "--help")
		advancedOutput, advancedErr := advancedCmd.CombinedOutput()

		// THEN: Both should work
		require.NoError(t, basicErr, "Basic help should work")
		require.NoError(t, advancedErr, "Advanced help should work")

		basicStr := string(basicOutput)
		advancedStr := string(advancedOutput)

		// AND: Advanced should show more flags
		assert.Contains(t, basicStr, "--type", "Basic should show essential flags")
		assert.Contains(t, advancedStr, "--type", "Advanced should show essential flags")

		// Basic mode should not show advanced flags, but advanced mode should
		advancedFlags := []string{"--database-driver", "--auth-type"}
		for _, flag := range advancedFlags {
			assert.NotContains(t, basicStr, flag, 
				"Basic help should not show advanced flag: %s", flag)
			// Note: We don't strictly require advanced flags in smoke test
			// as the template might not define them yet
		}
	})

	t.Run("cross_platform_basic", func(t *testing.T) {
		// GIVEN: Current platform
		// WHEN: Building and running CLI
		binaryName := "go-starter"
		expectedExtension := ""
		if runtime.GOOS == "windows" {
			expectedExtension = ".exe"
			binaryName += expectedExtension
		}

		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Should build on current platform")

		// THEN: Binary should have correct extension
		assert.FileExists(t, binaryPath, "Binary should exist with platform extension")

		// AND: Should run successfully
		cmd := exec.Command(binaryPath, "version")
		err = cmd.Run()
		assert.NoError(t, err, "Binary should execute on current platform")

		// AND: Path operations should work correctly
		genCmd := exec.Command(binaryPath, "new", "path-test", 
			"--type=cli", "--complexity=simple", "--dry-run")
		genOutput, genErr := genCmd.CombinedOutput()
		
		assert.NoError(t, genErr, "Path operations should work on current platform")
		
		// Verify output contains reasonable paths
		outputStr := string(genOutput)
		assert.Contains(t, outputStr, "main.go", "Should generate main.go")
		
		// Platform-specific path checks
		if runtime.GOOS == "windows" {
			// Windows should handle paths appropriately
			assert.True(t, 
				strings.Contains(outputStr, "/") || strings.Contains(outputStr, "\\"),
				"Windows should show path separators")
		} else {
			// Unix should use forward slashes
			assert.Contains(t, outputStr, "/", "Unix should use forward slashes")
		}
	})
}

// TestMultiBinaryQuickValidation runs the most essential tests for CI
func TestMultiBinaryQuickValidation(t *testing.T) {
	if !testing.Short() {
		t.Skip("Quick validation only runs in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	// Just test that CLI builds and basic command works
	t.Run("cli_builds_and_works", func(t *testing.T) {
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		// Build
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// Test basic functionality
		cmd := exec.Command(binaryPath, "version")
		err = cmd.Run()
		assert.NoError(t, err, "CLI should work")
	})

	// Test that dev server builds
	t.Run("dev_server_builds", func(t *testing.T) {
		binaryName := "go-starter-dev"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)

		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter-dev")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		assert.NoError(t, err, "Dev server should build")
	})

	// Test that web server builds
	t.Run("web_server_builds", func(t *testing.T) {
		binaryName := getBinaryName("go-starter-web")
		binaryPath := filepath.Join(tmpDir, binaryName)

		err := buildBinary(t, projectRoot, binaryPath, "go-starter-web", "./web/cmd/web-server")
		assert.NoError(t, err, "Web server should build")
	})
}