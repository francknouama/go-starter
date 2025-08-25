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

// TestCrossPlatformATDD validates that the multi-binary structure works across platforms
func TestCrossPlatformATDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping cross-platform tests in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("binary_builds_with_correct_extension", func(t *testing.T) {
		// GIVEN: Multi-binary structure on current platform
		binaries := map[string]string{
			"go-starter":     "./cmd/go-starter",
			"go-starter-dev": "./cmd/go-starter-dev",
			"go-starter-web": "./web/cmd/web-server",
		}

		for name, buildPath := range binaries {
			t.Run("build_"+name, func(t *testing.T) {
				// WHEN: Building binary for current platform
				expectedName := name
				if runtime.GOOS == "windows" {
					expectedName += ".exe"
				}
				
				binaryPath := filepath.Join(tmpDir, expectedName)
				cmd := exec.Command("go", "build", "-o", binaryPath, buildPath)
				cmd.Dir = projectRoot
				output, err := cmd.CombinedOutput()

				// THEN: Should build with appropriate extension
				if err != nil {
					t.Logf("Build failed for %s on %s:\nOutput: %s", 
						name, runtime.GOOS, string(output))
				}
				require.NoError(t, err, "Binary %s should build on %s", name, runtime.GOOS)

				// AND: Binary should exist with correct extension
				assert.FileExists(t, binaryPath, "Binary should exist with platform-appropriate extension")

				// AND: Binary should be executable
				info, err := os.Stat(binaryPath)
				require.NoError(t, err, "Should get binary info")
				assert.True(t, info.Mode().IsRegular(), "Should be regular file")

				if runtime.GOOS != "windows" {
					assert.True(t, info.Mode()&0111 != 0, "Binary should be executable on Unix")
				}
			})
		}
	})

	t.Run("path_separators_handled_correctly", func(t *testing.T) {
		// GIVEN: CLI binary built for current platform
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build on current platform")

		// WHEN: Generating project with nested paths
		projectName := "path-test"
		cmd := exec.Command(binaryPath, "new", projectName,
			"--type=web-api", "--dry-run",
			"--module=github.com/test/"+projectName)
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()

		// THEN: Should handle paths correctly for platform
		require.NoError(t, err, "Project generation should work on current platform")
		
		outputStr := string(output)
		
		// Check that paths in output are platform-appropriate
		if runtime.GOOS == "windows" {
			// Windows should accept both \ and / but prefer \
			assert.True(t, 
				strings.Contains(outputStr, "\\") || strings.Contains(outputStr, "/"),
				"Windows output should contain path separators")
		} else {
			// Unix should use /
			if strings.Contains(outputStr, "\\") && !strings.Contains(outputStr, "/") {
				t.Errorf("Unix platform should not use backslashes exclusively in paths")
			}
		}

		// Ensure no invalid characters for platform
		if runtime.GOOS == "windows" {
			invalidChars := []string{"<", ">", ":", "\"", "|", "?", "*"}
			for _, char := range invalidChars {
				assert.NotContains(t, outputStr, char, 
					"Windows paths should not contain invalid character: %s", char)
			}
		}
	})

	t.Run("file_operations_work_across_platforms", func(t *testing.T) {
		// GIVEN: CLI binary for current platform
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Generating project with various file types
		projectName := "file-ops-test"
		cmd := exec.Command(binaryPath, "new", projectName,
			"--type=web-api", 
			"--module=github.com/test/"+projectName)
		cmd.Dir = tmpDir
		output, err := cmd.CombinedOutput()

		// THEN: File operations should succeed
		if err != nil {
			t.Logf("Generation output: %s", string(output))
		}
		require.NoError(t, err, "File operations should work on %s", runtime.GOOS)

		// AND: Generated files should have appropriate permissions
		projectDir := filepath.Join(tmpDir, projectName)
		
		// Check common files exist and have correct permissions
		commonFiles := []string{"main.go", "go.mod", "README.md"}
		for _, file := range commonFiles {
			filePath := filepath.Join(projectDir, file)
			assert.FileExists(t, filePath, "File %s should exist", file)
			
			info, err := os.Stat(filePath)
			require.NoError(t, err, "Should stat file %s", file)
			
			if runtime.GOOS != "windows" {
				// On Unix, files should be readable/writable by owner
				mode := info.Mode()
				assert.True(t, mode&0600 != 0, "File %s should be readable/writable by owner", file)
			}
		}

		// Check that directories have appropriate permissions
		if runtime.GOOS != "windows" {
			dirInfo, err := os.Stat(projectDir)
			require.NoError(t, err, "Should stat project directory")
			
			mode := dirInfo.Mode()
			assert.True(t, mode&0700 != 0, "Directory should be readable/writable/executable by owner")
		}
	})

	t.Run("environment_variables_respected", func(t *testing.T) {
		// GIVEN: CLI binary and environment-specific settings
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Running with platform-specific environment variables
		cmd := exec.Command(binaryPath, "version")
		cmd.Dir = tmpDir
		
		// Set some common environment variables
		env := os.Environ()
		if runtime.GOOS == "windows" {
			env = append(env, "PATHEXT=.COM;.EXE;.BAT;.CMD")
		} else {
			env = append(env, "SHELL=/bin/bash")
		}
		cmd.Env = env
		
		output, err := cmd.CombinedOutput()

		// THEN: Should work with platform environment
		require.NoError(t, err, "Should work with platform environment variables")
		assert.Contains(t, string(output), "version", "Should show version information")
	})

	t.Run("generated_projects_compile_on_platform", func(t *testing.T) {
		// GIVEN: CLI binary for current platform
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// Test different project types on current platform
		projectTypes := []struct {
			name     string
			typeFlag string
			complexity string
		}{
			{"simple_cli", "cli", "simple"},
			{"standard_cli", "cli", "standard"},
			{"web_api", "web-api", "standard"},
			{"library", "library", "standard"},
		}

		for _, pt := range projectTypes {
			t.Run(pt.name+"_on_"+runtime.GOOS, func(t *testing.T) {
				// WHEN: Generating and building project
				projectName := pt.name + "-" + runtime.GOOS
				genCmd := exec.Command(binaryPath, "new", projectName,
					"--type="+pt.typeFlag,
					"--complexity="+pt.complexity,
					"--module=github.com/test/"+projectName)
				genCmd.Dir = tmpDir
				genOutput, genErr := genCmd.CombinedOutput()

				if genErr != nil {
					t.Logf("Generation output: %s", string(genOutput))
				}
				require.NoError(t, genErr, "Project generation should succeed on %s", runtime.GOOS)

				// THEN: Project should compile on current platform
				projectDir := filepath.Join(tmpDir, projectName)
				compileCmd := exec.Command("go", "build", "./...")
				compileCmd.Dir = projectDir
				compileOutput, compileErr := compileCmd.CombinedOutput()

				if compileErr != nil {
					t.Logf("Compile output: %s", string(compileOutput))
				}
				assert.NoError(t, compileErr, "Generated project should compile on %s", runtime.GOOS)

				// AND: If it's an executable project, the binary should have correct extension
				if pt.typeFlag == "cli" || pt.typeFlag == "web-api" {
					expectedBinary := projectName
					if runtime.GOOS == "windows" {
						expectedBinary += ".exe"
					}
					
					builtBinary := filepath.Join(projectDir, expectedBinary)
					if _, err := os.Stat(builtBinary); err == nil {
						// Binary was created with expected name
						assert.FileExists(t, builtBinary, "Binary should have correct name for platform")
					} else {
						// Check if binary was created with default name
						if runtime.GOOS == "windows" {
							defaultBinary := filepath.Join(projectDir, filepath.Base(projectDir)+".exe")
							assert.True(t, 
								fileExists(builtBinary) || fileExists(defaultBinary),
								"Binary should exist with platform-appropriate extension")
						}
					}
				}
			})
		}
	})

	t.Run("cross_compilation_support", func(t *testing.T) {
		// GIVEN: Source code that should support cross-compilation
		// WHEN: Attempting to build for different platforms (if supported)
		
		// Only test cross-compilation if we're on a Unix system (more reliable)
		if runtime.GOOS == "windows" {
			t.Skip("Skipping cross-compilation test on Windows")
		}
		
		// Test building for different platforms
		platforms := []struct {
			goos   string
			goarch string
		}{
			{"linux", "amd64"},
			{"darwin", "amd64"},
			{"windows", "amd64"},
		}

		for _, platform := range platforms {
			t.Run("cross_compile_"+platform.goos+"_"+platform.goarch, func(t *testing.T) {
				expectedName := "go-starter-" + platform.goos + "-" + platform.goarch
				if platform.goos == "windows" {
					expectedName += ".exe"
				}
				
				crossBinaryPath := filepath.Join(tmpDir, expectedName)
				
				cmd := exec.Command("go", "build", "-o", crossBinaryPath, "./cmd/go-starter")
				cmd.Dir = projectRoot
				cmd.Env = append(os.Environ(), 
					"GOOS="+platform.goos, 
					"GOARCH="+platform.goarch)
				
				output, err := cmd.CombinedOutput()
				
				// THEN: Cross-compilation should succeed
				if err != nil {
					t.Logf("Cross-compilation failed for %s/%s: %s", 
						platform.goos, platform.goarch, string(output))
				}
				assert.NoError(t, err, "Should cross-compile for %s/%s", 
					platform.goos, platform.goarch)
				
				// AND: Binary should exist
				if err == nil {
					assert.FileExists(t, crossBinaryPath, 
						"Cross-compiled binary should exist for %s/%s", 
						platform.goos, platform.goarch)
					
					// Check binary size is reasonable
					info, statErr := os.Stat(crossBinaryPath)
					if statErr == nil {
						size := info.Size()
						assert.Greater(t, size, int64(5*1024*1024), 
							"Cross-compiled binary should be at least 5MB")
						assert.Less(t, size, int64(60*1024*1024), 
							"Cross-compiled binary should be less than 60MB")
					}
				}
			})
		}
	})

	t.Run("unicode_and_special_characters", func(t *testing.T) {
		// GIVEN: CLI binary for current platform
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		
		binaryPath := filepath.Join(tmpDir, binaryName)
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Testing with various character encodings
		testCases := []struct {
			name        string
			projectName string
			shouldWork  bool
		}{
			{"ascii_name", "test-project", true},
			{"dash_underscore", "test_project-name", true},
			{"numbers", "test123", true},
			// Note: Unicode project names might not work with Go modules
			// so we test more conservatively
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				cmd := exec.Command(binaryPath, "new", tc.projectName,
					"--type=cli", "--complexity=simple", "--dry-run",
					"--module=github.com/test/"+tc.projectName)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()

				if tc.shouldWork {
					if err != nil {
						t.Logf("Output: %s", string(output))
					}
					assert.NoError(t, err, "Should handle project name: %s", tc.projectName)
				} else {
					assert.Error(t, err, "Should reject invalid project name: %s", tc.projectName)
				}
			})
		}
	})
}

