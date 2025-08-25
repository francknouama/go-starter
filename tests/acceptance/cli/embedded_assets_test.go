package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestEmbeddedAssetsATDD validates that embedded blueprints work correctly in the multi-binary structure
func TestEmbeddedAssetsATDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping embedded assets tests in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("cli_binary_embeds_all_blueprints", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// WHEN: Running from isolated directory without blueprint files
		isolatedDir := filepath.Join(tmpDir, "isolated")
		err = os.MkdirAll(isolatedDir, 0755)
		require.NoError(t, err)

		cmd := exec.Command(binaryPath, "list")
		cmd.Dir = isolatedDir // Important: run from directory without blueprints
		output, err := cmd.CombinedOutput()

		// THEN: Should list all embedded blueprints
		require.NoError(t, err, "List command should work with embedded blueprints")
		
		outputStr := string(output)
		expectedBlueprints := []string{
			"web-api", "cli", "library", "lambda", "microservice", "monolith",
		}
		
		for _, blueprint := range expectedBlueprints {
			assert.Contains(t, outputStr, blueprint, "Should list embedded blueprint: %s", blueprint)
		}
	})

	t.Run("embedded_blueprints_generate_valid_projects", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// Test different blueprint types
		testCases := []struct {
			name       string
			projectType string
			complexity  string
			expectFiles []string
		}{
			{
				name:        "simple_cli",
				projectType: "cli", 
				complexity:  "simple",
				expectFiles: []string{"main.go", "go.mod", "README.md"},
			},
			{
				name:        "standard_web_api",
				projectType: "web-api",
				complexity:  "standard",
				expectFiles: []string{"main.go", "go.mod", "README.md", "Dockerfile"},
			},
			{
				name:        "library_project",
				projectType: "library",
				complexity:  "standard", 
				expectFiles: []string{"go.mod", "README.md"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				// WHEN: Generating project with embedded blueprints
				projectName := "test-" + tc.name
				cmd := exec.Command(binaryPath, "new", projectName,
					"--type="+tc.projectType,
					"--complexity="+tc.complexity,
					"--module=github.com/test/"+projectName)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()

				// THEN: Project generation should succeed
				if err != nil {
					t.Logf("Generation output: %s", string(output))
				}
				require.NoError(t, err, "Project generation should succeed")

				// AND: Expected files should exist
				projectDir := filepath.Join(tmpDir, projectName)
				for _, file := range tc.expectFiles {
					filePath := filepath.Join(projectDir, file)
					assert.FileExists(t, filePath, "Expected file should exist: %s", file)
				}

				// AND: Generated project should compile
				compileCmd := exec.Command("go", "build", "./...")
				compileCmd.Dir = projectDir
				compileOutput, compileErr := compileCmd.CombinedOutput()
				if compileErr != nil {
					t.Logf("Compile output: %s", string(compileOutput))
				}
				assert.NoError(t, compileErr, "Generated project should compile")
			})
		}
	})

	t.Run("embedded_blueprints_contain_all_template_variables", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// WHEN: Generating projects with various configurations
		testConfigs := []struct {
			name   string
			args   []string
			checks []string // Things that should be properly templated
		}{
			{
				name: "module_path_templating",
				args: []string{"new", "test-module", "--type=cli", "--complexity=simple", 
					"--module=github.com/custom/test-module"},
				checks: []string{"github.com/custom/test-module"}, 
			},
			{
				name: "project_name_templating",
				args: []string{"new", "custom-project-name", "--type=library",
					"--module=github.com/test/custom-project-name"},
				checks: []string{"custom-project-name"},
			},
			{
				name: "logger_type_templating",
				args: []string{"new", "logger-test", "--type=web-api", 
					"--logger=zap", "--module=github.com/test/logger-test"},
				checks: []string{"zap"}, // Should use zap logger
			},
		}

		for _, tc := range testConfigs {
			t.Run(tc.name, func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()

				if err != nil {
					t.Logf("Generation output: %s", string(output))
				}
				require.NoError(t, err, "Project generation should succeed")

				// Check that template variables were properly substituted
				projectName := tc.args[1] // Second argument is project name
				projectDir := filepath.Join(tmpDir, projectName)

				// Read key files and check for proper templating
				filesToCheck := []string{"go.mod", "main.go", "README.md"}
				for _, file := range filesToCheck {
					filePath := filepath.Join(projectDir, file)
					if _, err := os.Stat(filePath); err == nil {
						content, err := os.ReadFile(filePath)
						require.NoError(t, err, "Should read file %s", file)
						
						contentStr := string(content)
						for _, check := range tc.checks {
							assert.Contains(t, contentStr, check, 
								"File %s should contain templated value: %s", file, check)
						}

						// Ensure no unprocessed template variables remain
						assert.NotContains(t, contentStr, "{{", 
							"File %s should not contain unprocessed template variables", file)
						assert.NotContains(t, contentStr, "}}", 
							"File %s should not contain unprocessed template variables", file)
					}
				}
			})
		}
	})

	t.Run("embedded_blueprints_support_conditional_generation", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// WHEN: Generating projects with different logger configurations
		loggerTests := []struct {
			logger        string
			shouldContain []string
			shouldNotContain []string
		}{
			{
				logger: "slog",
				shouldContain: []string{"log/slog"},
				shouldNotContain: []string{"go.uber.org/zap", "github.com/sirupsen/logrus"},
			},
			{
				logger: "zap", 
				shouldContain: []string{"go.uber.org/zap"},
				shouldNotContain: []string{"log/slog", "github.com/sirupsen/logrus"},
			},
			{
				logger: "logrus",
				shouldContain: []string{"github.com/sirupsen/logrus"},
				shouldNotContain: []string{"log/slog", "go.uber.org/zap"},
			},
		}

		for _, tc := range loggerTests {
			t.Run("logger_"+tc.logger, func(t *testing.T) {
				projectName := "logger-" + tc.logger
				cmd := exec.Command(binaryPath, "new", projectName,
					"--type=web-api",
					"--logger="+tc.logger,
					"--module=github.com/test/"+projectName)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()

				if err != nil {
					t.Logf("Generation output: %s", string(output))
				}
				require.NoError(t, err, "Project generation should succeed")

				// Check go.mod for correct dependencies
				projectDir := filepath.Join(tmpDir, projectName)
				goModPath := filepath.Join(projectDir, "go.mod")
				goModContent, err := os.ReadFile(goModPath)
				require.NoError(t, err, "Should read go.mod")

				goModStr := string(goModContent)
				for _, shouldContain := range tc.shouldContain {
					assert.Contains(t, goModStr, shouldContain,
						"go.mod should contain logger dependency: %s", shouldContain)
				}

				for _, shouldNotContain := range tc.shouldNotContain {
					assert.NotContains(t, goModStr, shouldNotContain,
						"go.mod should not contain unwanted logger dependency: %s", shouldNotContain)
				}
			})
		}
	})

	t.Run("embedded_blueprint_validation", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// WHEN: Testing blueprint validation with invalid configurations
		invalidConfigs := []struct {
			name     string
			args     []string
			expectError bool
		}{
			{
				name: "invalid_project_type",
				args: []string{"new", "test", "--type=invalid-type"},
				expectError: true,
			},
			{
				name: "invalid_complexity",
				args: []string{"new", "test", "--type=cli", "--complexity=invalid"},
				expectError: true,
			},
			{
				name: "invalid_logger",
				args: []string{"new", "test", "--type=web-api", "--logger=invalid-logger"},
				expectError: true,
			},
		}

		for _, tc := range invalidConfigs {
			t.Run(tc.name, func(t *testing.T) {
				cmd := exec.Command(binaryPath, tc.args...)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()

				if tc.expectError {
					// THEN: Should show validation error
					assert.Error(t, err, "Invalid configuration should be rejected")
					outputStr := string(output)
					assert.True(t, 
						strings.Contains(outputStr, "invalid") || 
						strings.Contains(outputStr, "error") ||
						strings.Contains(outputStr, "unknown"),
						"Error message should indicate validation failure: %s", outputStr)
				} else {
					assert.NoError(t, err, "Valid configuration should succeed")
				}
			})
		}
	})

	t.Run("embedded_assets_file_size_reasonable", func(t *testing.T) {
		// GIVEN: CLI binary with embedded blueprints
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// WHEN: Checking binary size
		info, err := os.Stat(binaryPath)
		require.NoError(t, err, "Should get binary file info")

		// THEN: Binary size should be reasonable even with embedded assets
		size := info.Size()
		maxSizeWithAssets := int64(60 * 1024 * 1024) // 60MB max with assets
		minSizeWithAssets := int64(10 * 1024 * 1024)  // 10MB min with assets

		assert.GreaterOrEqual(t, size, minSizeWithAssets, 
			"Binary with embedded assets should be at least 10MB")
		assert.LessOrEqual(t, size, maxSizeWithAssets,
			"Binary with embedded assets should not exceed 60MB")

		t.Logf("CLI binary size with embedded assets: %.2f MB", float64(size)/(1024*1024))
	})

	t.Run("no_filesystem_dependency_for_blueprints", func(t *testing.T) {
		// GIVEN: CLI binary built and moved to isolated location
		binaryPath := filepath.Join(tmpDir, "go-starter")
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI binary should build successfully")

		// Create completely isolated environment
		isolatedDir := filepath.Join(tmpDir, "completely-isolated")
		err = os.MkdirAll(isolatedDir, 0755)
		require.NoError(t, err)

		// Copy binary to isolated location
		isolatedBinary := filepath.Join(isolatedDir, "go-starter")
		copyCmd := exec.Command("cp", binaryPath, isolatedBinary)
		err = copyCmd.Run()
		require.NoError(t, err, "Should copy binary to isolated location")

		// WHEN: Running from isolated environment with no access to source
		cmd := exec.Command(isolatedBinary, "new", "isolated-test",
			"--type=cli", "--complexity=simple",
			"--module=github.com/test/isolated-test")
		cmd.Dir = isolatedDir
		output, err := cmd.CombinedOutput()

		// THEN: Should work without any filesystem access to original blueprints
		if err != nil {
			t.Logf("Isolated generation output: %s", string(output))
		}
		assert.NoError(t, err, "Should work with embedded blueprints only")

		// AND: Should generate valid project
		projectDir := filepath.Join(isolatedDir, "isolated-test")
		assert.DirExists(t, projectDir, "Should create project directory")
		assert.FileExists(t, filepath.Join(projectDir, "main.go"), "Should create main.go")
		assert.FileExists(t, filepath.Join(projectDir, "go.mod"), "Should create go.mod")

		// AND: Generated project should compile
		compileCmd := exec.Command("go", "build", "./...")
		compileCmd.Dir = projectDir
		compileErr := compileCmd.Run()
		assert.NoError(t, compileErr, "Generated project should compile in isolation")
	})
}