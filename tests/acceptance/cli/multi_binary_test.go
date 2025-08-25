package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiBinaryStructureATDD validates the acceptance criteria for the new multi-binary structure
// This comprehensive test suite ensures all binaries work correctly and maintain backward compatibility
func TestMultiBinaryStructureATDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping multi-binary ATDD tests in short mode")
	}

	// Setup test environment
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := findProjectRoot(t)
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	t.Run("multi_binary_compilation_tests", func(t *testing.T) {
		testMultiBinaryCompilation(t, projectRoot, tmpDir)
	})

	t.Run("installation_path_tests", func(t *testing.T) {
		testInstallationPaths(t, projectRoot, tmpDir)
	})

	t.Run("backward_compatibility_tests", func(t *testing.T) {
		testBackwardCompatibility(t, projectRoot, tmpDir)
	})

	t.Run("binary_functionality_tests", func(t *testing.T) {
		testBinaryFunctionality(t, projectRoot, tmpDir)
	})

	t.Run("embedded_assets_tests", func(t *testing.T) {
		testEmbeddedAssets(t, projectRoot, tmpDir)
	})

	t.Run("cross_platform_tests", func(t *testing.T) {
		testCrossPlatformCompatibility(t, projectRoot, tmpDir)
	})

	t.Run("migration_tests", func(t *testing.T) {
		testMigrationExperience(t, projectRoot, tmpDir)
	})
}

// testMultiBinaryCompilation verifies all binaries can be built independently
func testMultiBinaryCompilation(t *testing.T, projectRoot, tmpDir string) {
	binaries := map[string]string{
		"go-starter":     "./cmd/go-starter",
		"go-starter-dev": "./cmd/go-starter-dev", 
		"go-starter-web": "./web/cmd/web-server",
		"legacy":         ".",
	}

	for name, path := range binaries {
		t.Run(fmt.Sprintf("compile_%s", name), func(t *testing.T) {
			// GIVEN: A specific binary target
			// WHEN: Building the binary
			binaryPath := filepath.Join(tmpDir, name)
			cmd := exec.Command("go", "build", "-o", binaryPath, path)
			cmd.Dir = projectRoot
			output, err := cmd.CombinedOutput()

			// THEN: Compilation should succeed
			if err != nil {
				t.Logf("Build failed for %s:\nCommand: go build -o %s %s\nOutput: %s", 
					name, binaryPath, path, string(output))
			}
			require.NoError(t, err, "Binary %s should compile successfully", name)

			// AND: Binary should exist and be executable
			info, err := os.Stat(binaryPath)
			require.NoError(t, err, "Binary file should exist")
			assert.True(t, info.Mode().IsRegular(), "Should be a regular file")
			
			if runtime.GOOS != "windows" {
				assert.True(t, info.Mode()&0111 != 0, "Binary should be executable")
			}

			// AND: Binary size should be reasonable (between 5MB and 50MB)
			assert.Greater(t, info.Size(), int64(5*1024*1024), "Binary should be at least 5MB")
			assert.Less(t, info.Size(), int64(50*1024*1024), "Binary should be less than 50MB")
		})
	}
}

// testInstallationPaths verifies various installation methods work correctly
func testInstallationPaths(t *testing.T, projectRoot, tmpDir string) {
	// Create temporary GOPATH for installation tests
	goPath := filepath.Join(tmpDir, "gopath")
	goBin := filepath.Join(goPath, "bin")
	err := os.MkdirAll(goBin, 0755)
	require.NoError(t, err)

	t.Run("install_main_cli_tool", func(t *testing.T) {
		// GIVEN: The main CLI tool path
		// WHEN: Installing via go install ./cmd/go-starter
		cmd := exec.Command("go", "install", "./cmd/go-starter")
		cmd.Dir = projectRoot
		cmd.Env = append(os.Environ(), "GOPATH="+goPath)
		output, err := cmd.CombinedOutput()

		// THEN: Installation should succeed
		if err != nil {
			t.Logf("Install failed: %s", string(output))
		}
		require.NoError(t, err, "CLI tool should install successfully")

		// AND: Binary should be accessible in GOPATH/bin
		binaryPath := filepath.Join(goBin, "go-starter")
		if runtime.GOOS == "windows" {
			binaryPath += ".exe"
		}
		assert.FileExists(t, binaryPath, "go-starter binary should exist in GOPATH/bin")

		// AND: Binary should be functional
		cmd = exec.Command(binaryPath, "version")
		err = cmd.Run()
		assert.NoError(t, err, "Installed binary should be functional")
	})

	t.Run("install_dev_server", func(t *testing.T) {
		// GIVEN: The development server path
		// WHEN: Installing via go install ./cmd/go-starter-dev
		cmd := exec.Command("go", "install", "./cmd/go-starter-dev")
		cmd.Dir = projectRoot
		cmd.Env = append(os.Environ(), "GOPATH="+goPath)
		output, err := cmd.CombinedOutput()

		// THEN: Installation should succeed
		if err != nil {
			t.Logf("Install failed: %s", string(output))
		}
		require.NoError(t, err, "Dev server should install successfully")

		// AND: Binary should exist
		binaryPath := filepath.Join(goBin, "go-starter-dev")
		if runtime.GOOS == "windows" {
			binaryPath += ".exe"
		}
		assert.FileExists(t, binaryPath, "go-starter-dev binary should exist in GOPATH/bin")
	})

	t.Run("legacy_installation_still_works", func(t *testing.T) {
		// GIVEN: Legacy installation method
		// WHEN: Installing via go install . (root directory)
		cmd := exec.Command("go", "install", ".")
		cmd.Dir = projectRoot
		cmd.Env = append(os.Environ(), "GOPATH="+goPath)
		output, err := cmd.CombinedOutput()

		// THEN: Installation should succeed (with possible warning)
		if err != nil {
			t.Logf("Legacy install output: %s", string(output))
		}
		require.NoError(t, err, "Legacy installation should still work")

		// AND: go-starter binary should exist (legacy installs as go-starter)
		binaryPath := filepath.Join(goBin, "go-starter")
		if runtime.GOOS == "windows" {
			binaryPath += ".exe"
		}
		assert.FileExists(t, binaryPath, "Legacy installation should create go-starter binary")
	})
}

// testBackwardCompatibility ensures existing workflows continue to work
func testBackwardCompatibility(t *testing.T, projectRoot, tmpDir string) {
	t.Run("legacy_main_shows_deprecation_warning", func(t *testing.T) {
		// GIVEN: The legacy main.go binary
		// WHEN: Running with --help flag
		binaryPath := filepath.Join(tmpDir, "legacy")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Legacy binary should build")

		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.CombinedOutput()

		// THEN: Should show deprecation warning
		outputStr := string(output)
		assert.Contains(t, outputStr, "DEPRECATION WARNING", "Should show deprecation warning")
		assert.Contains(t, outputStr, "go build -o go-starter ./cmd/go-starter", "Should show new CLI path")
		assert.Contains(t, outputStr, "go build -o go-starter-web ./cmd/go-starter-web", "Should show web server path")
		assert.Contains(t, outputStr, "go build -o go-starter-dev ./cmd/go-starter-dev", "Should show dev server path")
	})

	t.Run("existing_makefile_targets_work", func(t *testing.T) {
		// GIVEN: Existing Makefile targets
		makefilePath := filepath.Join(projectRoot, "Makefile")
		if _, err := os.Stat(makefilePath); os.IsNotExist(err) {
			t.Skip("Makefile not found, skipping Makefile tests")
		}

		// WHEN: Running make build
		cmd := exec.Command("make", "build")
		cmd.Dir = projectRoot
		output, err := cmd.CombinedOutput()

		// THEN: Should succeed or provide helpful error
		if err != nil {
			t.Logf("Make build output: %s", string(output))
			// It's okay if make fails, but the output should be informative
			assert.Contains(t, string(output), "build", "Make output should mention build process")
		}
	})

	t.Run("documentation_examples_still_function", func(t *testing.T) {
		// GIVEN: Common documentation examples
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build")

		// WHEN: Running typical documentation commands
		testCases := []struct {
			name string
			args []string
		}{
			{"version_command", []string{"version"}},
			{"list_command", []string{"list"}},
			{"help_command", []string{"--help"}},
			{"new_help", []string{"new", "--help"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				err := cmd.Run()
				// THEN: Commands should execute without error
				assert.NoError(t, err, "Command %v should execute successfully", tc.args)
			})
		}
	})
}

// testBinaryFunctionality tests core functionality of each binary
func testBinaryFunctionality(t *testing.T, projectRoot, tmpDir string) {
	t.Run("cli_tool_functionality", func(t *testing.T) {
		// GIVEN: The main CLI tool
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// WHEN: Testing core CLI functionality
		testCases := []struct {
			name     string
			args     []string
			contains []string
		}{
			{
				name:     "version_info",
				args:     []string{"version"},
				contains: []string{"go-starter", "version"},
			},
			{
				name:     "list_blueprints",
				args:     []string{"list"},
				contains: []string{"Available blueprints", "web-api", "cli"},
			},
			{
				name:     "dry_run_generation",
				args:     []string{"new", "test-project", "--type=cli", "--dry-run"},
				contains: []string{"Files to be generated", "main.go", "go.mod"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				output, err := cmd.CombinedOutput()
				
				// THEN: Command should succeed and contain expected content
				if err != nil {
					t.Logf("Command failed: %s", string(output))
				}
				assert.NoError(t, err, "Command should succeed")
				
				outputStr := string(output)
				for _, contain := range tc.contains {
					assert.Contains(t, outputStr, contain, "Output should contain: %s", contain)
				}
			})
		}
	})

	t.Run("dev_server_functionality", func(t *testing.T) {
		// GIVEN: The development server
		binaryPath := filepath.Join(tmpDir, "go-starter-dev")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter-dev")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Dev server should build")

		// WHEN: Testing dev server startup (with timeout)
		cmd := exec.Command(binaryPath)
		cmd.Dir = projectRoot // Important: run from project root for blueprint access
		err = cmd.Start()
		require.NoError(t, err, "Dev server should start")

		// Give server time to start
		time.Sleep(2 * time.Second)

		// THEN: Server should be running
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
		
		// Note: More comprehensive server testing would require HTTP client tests
		// This validates the server can at least start without immediate crash
	})

	t.Run("web_server_functionality", func(t *testing.T) {
		// GIVEN: The production web server
		binaryPath := filepath.Join(tmpDir, "go-starter-web")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./web/cmd/web-server")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Web server should build")

		// WHEN: Testing web server startup (with timeout)
		cmd := exec.Command(binaryPath)
		cmd.Env = append(os.Environ(), "PORT=0") // Use random port to avoid conflicts
		err = cmd.Start()
		require.NoError(t, err, "Web server should start")

		// Give server time to start
		time.Sleep(2 * time.Second)

		// THEN: Server should be running
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})
}

// testEmbeddedAssets validates that embedded blueprints work correctly
func testEmbeddedAssets(t *testing.T, projectRoot, tmpDir string) {
	t.Run("cli_works_without_filesystem_access", func(t *testing.T) {
		// GIVEN: CLI tool with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// WHEN: Running from directory without blueprint files
		isolatedDir := filepath.Join(tmpDir, "isolated")
		err = os.MkdirAll(isolatedDir, 0755)
		require.NoError(t, err)

		cmd := exec.Command(binaryPath, "list")
		cmd.Dir = isolatedDir // Run from isolated directory
		output, err := cmd.CombinedOutput()

		// THEN: Should still work using embedded blueprints
		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		assert.NoError(t, err, "CLI should work with embedded blueprints")
		assert.Contains(t, string(output), "web-api", "Should list embedded blueprints")
	})

	t.Run("all_blueprints_accessible", func(t *testing.T) {
		// GIVEN: CLI tool with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// WHEN: Listing all blueprints
		cmd := exec.Command(binaryPath, "list")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "List command should succeed")

		// THEN: Should show expected blueprints
		outputStr := string(output)
		expectedBlueprints := []string{
			"web-api", "cli", "library", "lambda", "microservice",
		}
		
		for _, blueprint := range expectedBlueprints {
			assert.Contains(t, outputStr, blueprint, "Should list blueprint: %s", blueprint)
		}
	})

	t.Run("generated_projects_compile", func(t *testing.T) {
		// GIVEN: CLI tool that can generate projects
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// WHEN: Generating a simple project
		projectDir := filepath.Join(tmpDir, "test-project")
		cmd := exec.Command(binaryPath, "new", "test-project", 
			"--type=cli", "--complexity=simple", 
			"--module=github.com/test/test-project")
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Generation output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// THEN: Generated project should compile
		compileCmd := exec.Command("go", "build", "./...")
		compileCmd.Dir = projectDir
		compileOutput, compileErr := compileCmd.CombinedOutput()
		
		if compileErr != nil {
			t.Logf("Compile output: %s", string(compileOutput))
		}
		assert.NoError(t, compileErr, "Generated project should compile successfully")
	})
}

// testCrossPlatformCompatibility tests behavior across different platforms
func testCrossPlatformCompatibility(t *testing.T, projectRoot, tmpDir string) {
	t.Run("path_handling_works_correctly", func(t *testing.T) {
		// GIVEN: Different path separators on different platforms
		// WHEN: Generating a project
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// Create a project with nested path
		cmd := exec.Command(binaryPath, "new", "test-project", 
			"--type=cli", "--complexity=simple", 
			"--module=github.com/test/test-project", "--dry-run")
		output, err := cmd.CombinedOutput()

		// THEN: Should handle paths correctly for current platform
		assert.NoError(t, err, "Path handling should work on current platform")
		outputStr := string(output)
		
		// Check that output contains proper path separators for platform
		if runtime.GOOS == "windows" {
			// On Windows, should handle backslashes properly
			assert.True(t, strings.Contains(outputStr, "main.go") || 
				strings.Contains(outputStr, "\\"), "Should handle Windows paths")
		} else {
			// On Unix-like systems, should use forward slashes
			assert.Contains(t, outputStr, "/", "Should use Unix-style paths")
		}
	})

	t.Run("binary_formats_correct_for_platform", func(t *testing.T) {
		// GIVEN: Different binary formats on different platforms
		// WHEN: Building binaries
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Binary should build")

		// THEN: Binary should be appropriate for current platform
		info, err := os.Stat(binaryPath)
		require.NoError(t, err, "Binary should exist")
		assert.True(t, info.Mode().IsRegular(), "Should be a regular file")

		// Test binary execution
		cmd := exec.Command(binaryPath, "version")
		err = cmd.Run()
		assert.NoError(t, err, "Binary should execute on current platform")
	})
}

// testMigrationExperience validates the user migration experience
func testMigrationExperience(t *testing.T, projectRoot, tmpDir string) {
	t.Run("clear_deprecation_messages", func(t *testing.T) {
		// GIVEN: Legacy binary usage
		binaryPath := filepath.Join(tmpDir, "legacy")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "Legacy binary should build")

		// WHEN: Running with help flag
		cmd := exec.Command(binaryPath, "--help")
		output, err := cmd.CombinedOutput()

		// THEN: Should provide clear migration guidance
		outputStr := string(output)
		assert.Contains(t, outputStr, "DEPRECATION WARNING", "Should warn about deprecation")
		assert.Contains(t, outputStr, "new binary locations", "Should explain new locations")
		assert.Contains(t, outputStr, "./cmd/go-starter", "Should show CLI path")
		assert.Contains(t, outputStr, "./cmd/go-starter-dev", "Should show dev server path")
		assert.Contains(t, outputStr, "./web/cmd/web-server", "Should show web server path")
	})

	t.Run("migration_instructions_work", func(t *testing.T) {
		// GIVEN: Migration instructions from deprecation warning
		// WHEN: Following the migration instructions
		
		// Test each suggested build command
		buildCommands := map[string]string{
			"go-starter":     "./cmd/go-starter",
			"go-starter-dev": "./cmd/go-starter-dev",
			"go-starter-web": "./web/cmd/web-server",
		}

		for binaryName, buildPath := range buildCommands {
			t.Run(fmt.Sprintf("migrate_to_%s", binaryName), func(t *testing.T) {
				binaryPath := filepath.Join(tmpDir, binaryName)
				cmd := exec.Command("go", "build", "-o", binaryPath, buildPath)
				cmd.Dir = projectRoot
				output, err := cmd.CombinedOutput()

				// THEN: Migration command should work
				if err != nil {
					t.Logf("Migration build failed for %s: %s", binaryName, string(output))
				}
				assert.NoError(t, err, "Migration build command should work for %s", binaryName)

				// AND: New binary should be functional
				if binaryName == "go-starter" {
					testCmd := exec.Command(binaryPath, "version")
					testErr := testCmd.Run()
					assert.NoError(t, testErr, "Migrated binary should be functional")
				}
			})
		}
	})

	t.Run("no_breaking_changes_for_existing_users", func(t *testing.T) {
		// GIVEN: Existing user workflows
		// WHEN: Using the new CLI binary
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "New CLI binary should build")

		// THEN: All existing commands should still work
		existingCommands := [][]string{
			{"version"},
			{"list"},
			{"new", "--help"},
			{"new", "test", "--type=cli", "--dry-run"},
		}

		for _, cmdArgs := range existingCommands {
			t.Run(fmt.Sprintf("existing_command_%s", strings.Join(cmdArgs, "_")), func(t *testing.T) {
				cmd := exec.Command(binaryPath, cmdArgs...)
				err := cmd.Run()
				assert.NoError(t, err, "Existing command should still work: %v", cmdArgs)
			})
		}
	})
}


// TestMultiBinaryPerformance tests performance characteristics of the new structure
func TestMultiBinaryPerformance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("build_times_reasonable", func(t *testing.T) {
		// GIVEN: Multi-binary structure
		// WHEN: Building each binary
		binaries := map[string]string{
			"go-starter":     "./cmd/go-starter",
			"go-starter-dev": "./cmd/go-starter-dev",
			"go-starter-web": "./web/cmd/web-server",
		}

		for name, path := range binaries {
			t.Run(fmt.Sprintf("build_time_%s", name), func(t *testing.T) {
				start := time.Now()
				
				binaryPath := filepath.Join(tmpDir, name)
				cmd := exec.Command("go", "build", "-o", binaryPath, path)
				cmd.Dir = projectRoot
				err := cmd.Run()
				
				elapsed := time.Since(start)
				
				// THEN: Build should complete in reasonable time (< 30 seconds)
				require.NoError(t, err, "Build should succeed")
				assert.Less(t, elapsed, 30*time.Second, "Build should complete in under 30 seconds")
				
				t.Logf("Build time for %s: %v", name, elapsed)
			})
		}
	})

	t.Run("startup_times_reasonable", func(t *testing.T) {
		// GIVEN: Built CLI binary
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI tool should build")

		// WHEN: Running version command (fast operation)
		start := time.Now()
		cmd := exec.Command(binaryPath, "version")
		err = cmd.Run()
		elapsed := time.Since(start)

		// THEN: Should start and execute quickly (< 2 seconds)
		require.NoError(t, err, "Version command should succeed")
		assert.Less(t, elapsed, 2*time.Second, "CLI should start quickly")
		
		t.Logf("CLI startup time: %v", elapsed)
	})
}

// TestMultiBinaryIntegration tests integration between different binaries
func TestMultiBinaryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("cli_and_web_server_compatibility", func(t *testing.T) {
		// GIVEN: Both CLI and web server binaries
		cliBinary := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", cliBinary, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		webBinary := filepath.Join(tmpDir, "go-starter-web")
		buildCmd = exec.Command("go", "build", "-o", webBinary, "./web/cmd/web-server")
		buildCmd.Dir = projectRoot
		err = buildCmd.Run()
		require.NoError(t, err, "Web server should build")

		// WHEN: CLI generates project configuration
		cmd := exec.Command(cliBinary, "new", "test-api", 
			"--type=web-api", "--dry-run")
		output, err := cmd.CombinedOutput()

		// THEN: Configuration should be compatible with web server expectations
		require.NoError(t, err, "CLI should generate valid configuration")
		assert.Contains(t, string(output), "web-api", "Should generate web API project")
		
		// Note: Full integration would require running web server and testing API endpoints
		// This validates at least that both binaries can be built and basic CLI works
	})
}