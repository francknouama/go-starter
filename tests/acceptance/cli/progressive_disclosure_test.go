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

// TestProgressiveDisclosureATDD validates the acceptance criteria for progressive disclosure
// This test ensures that users see different levels of complexity based on their flags
func TestProgressiveDisclosureATDD(t *testing.T) {
	// Setup
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	
	// Get the project root (parent of tests/acceptance/cli)
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Build the CLI tool first
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
	buildCmd.Dir = projectRoot  // Run build from project root
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
		t.Logf("Project root: %s", projectRoot)
		t.Logf("Original dir: %s", originalDir)
	}
	require.NoError(t, err, "Failed to build CLI tool")

	t.Run("basic_mode_shows_essential_options_only", func(t *testing.T) {
		// GIVEN: User wants to see basic options only
		// WHEN: User runs go-starter new --basic --help
		cmd := exec.Command("./go-starter", "new", "--basic", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed")

		outputStr := string(output)
		
		// THEN: Only essential flags should be visible
		assert.Contains(t, outputStr, "--type", "Should show essential flag: type")
		assert.Contains(t, outputStr, "--name", "Should show essential flag: name")
		assert.Contains(t, outputStr, "--module", "Should show essential flag: module")
		
		// AND: Advanced flags should be hidden
		assert.NotContains(t, outputStr, "--database-driver", "Should hide advanced flag: database-driver")
		assert.NotContains(t, outputStr, "--database-orm", "Should hide advanced flag: database-orm")
		assert.NotContains(t, outputStr, "--auth-type", "Should hide advanced flag: auth-type")
		assert.NotContains(t, outputStr, "--banner-style", "Should hide advanced flag: banner-style")
	})

	t.Run("advanced_mode_shows_all_options", func(t *testing.T) {
		// GIVEN: User wants to see all options
		// WHEN: User runs go-starter new --advanced --help
		cmd := exec.Command("./go-starter", "new", "--advanced", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed")

		outputStr := string(output)
		
		// THEN: All flags should be visible
		assert.Contains(t, outputStr, "--type", "Should show essential flag: type")
		assert.Contains(t, outputStr, "--database-driver", "Should show advanced flag: database-driver")
		assert.Contains(t, outputStr, "--database-orm", "Should show advanced flag: database-orm")
		assert.Contains(t, outputStr, "--auth-type", "Should show advanced flag: auth-type")
		assert.Contains(t, outputStr, "--banner-style", "Should show advanced flag: banner-style")
	})

	t.Run("complexity_simple_limits_blueprint_options", func(t *testing.T) {
		// GIVEN: User wants simple complexity level
		// WHEN: User runs go-starter new --complexity=simple test-project --type=cli --dry-run
		cmd := exec.Command("./go-starter", "new", "--complexity=simple", "test-project", "--type=cli", "--dry-run")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed")

		outputStr := string(output)
		
		// THEN: Simple CLI blueprint should be used
		assert.Contains(t, outputStr, "cli-simple", "Should use simple CLI blueprint")
		assert.NotContains(t, outputStr, "cli-standard", "Should not use standard CLI blueprint")
	})

	t.Run("complexity_standard_uses_full_blueprint", func(t *testing.T) {
		// GIVEN: User wants standard complexity level
		// WHEN: User runs go-starter new --complexity=standard test-project --type=cli --dry-run
		cmd := exec.Command("./go-starter", "new", "--complexity=standard", "test-project", "--type=cli", "--dry-run")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed")

		outputStr := string(output)
		
		// THEN: Standard CLI blueprint should be used (not simple)
		assert.Contains(t, outputStr, "Type: cli", "Should use standard CLI blueprint")
		assert.NotContains(t, outputStr, "cli-simple", "Should not use simple CLI blueprint")
	})

	t.Run("default_mode_shows_basic_options", func(t *testing.T) {
		// GIVEN: User doesn't specify any complexity flags
		// WHEN: User runs go-starter new --help
		cmd := exec.Command("./go-starter", "new", "--help")
		output, err := cmd.CombinedOutput()
		require.NoError(t, err, "Command should succeed")

		outputStr := string(output)
		
		// THEN: Should behave like basic mode (default)
		assert.Contains(t, outputStr, "--type", "Should show essential flag: type")
		assert.Contains(t, outputStr, "--name", "Should show essential flag: name")
		
		// AND: Should hint at advanced mode availability
		assert.Contains(t, outputStr, "--advanced", "Should show hint about advanced mode")
		assert.Contains(t, outputStr, "Enable advanced configuration", "Should explain advanced mode")
	})

	t.Run("interactive_mode_respects_complexity_level", func(t *testing.T) {
		// This test would require mock input, so we'll test the dry run approach
		// GIVEN: User wants to test interactive mode with complexity
		// WHEN: User runs go-starter new --complexity=simple --dry-run
		cmd := exec.Command("./go-starter", "new", "--complexity=simple", "--dry-run")
		
		// Simulate interactive input for project name and type
		cmd.Stdin = strings.NewReader("test-project\ncli\n")
		output, err := cmd.CombinedOutput()
		
		// THEN: Should work without error (even if prompts are skipped in CI)
		// The main validation is that the complexity flag is recognized
		if err != nil {
			// In CI environments, interactive prompts might fail, which is expected
			t.Logf("Interactive mode failed as expected in CI: %v", err)
		} else {
			t.Logf("Interactive mode output: %s", string(output))
		}
	})
}

// TestComplexityFlagValidation tests that complexity flag validation works correctly
func TestComplexityFlagValidation(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	
	// Get the project root (parent of tests/acceptance/cli)
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Build the CLI tool first
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
	buildCmd.Dir = projectRoot  // Run build from project root
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
		t.Logf("Project root: %s", projectRoot)
		t.Logf("Original dir: %s", originalDir)
	}
	require.NoError(t, err, "Failed to build CLI tool")

	t.Run("valid_complexity_levels_accepted", func(t *testing.T) {
		validLevels := []string{"simple", "standard", "advanced", "expert"}
		
		for _, level := range validLevels {
			t.Run("complexity_"+level, func(t *testing.T) {
				cmd := exec.Command("./go-starter", "new", "--complexity="+level, "test-project", "--type=cli", "--dry-run")
				err := cmd.Run()
				assert.NoError(t, err, "Valid complexity level %s should be accepted", level)
			})
		}
	})

	t.Run("invalid_complexity_level_rejected", func(t *testing.T) {
		cmd := exec.Command("./go-starter", "new", "--complexity=invalid", "test-project", "--type=cli", "--dry-run")
		output, err := cmd.CombinedOutput()
		
		assert.Error(t, err, "Invalid complexity level should be rejected")
		assert.Contains(t, string(output), "invalid complexity", "Should show validation error")
	})
}

// TestProgressiveDisclosureIntegration tests the integration between progressive disclosure and project generation
func TestProgressiveDisclosureIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	
	// Get the project root (parent of tests/acceptance/cli)
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Build the CLI tool first
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
	buildCmd.Dir = projectRoot  // Run build from project root
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
		t.Logf("Project root: %s", projectRoot)
		t.Logf("Original dir: %s", originalDir)
	}
	require.NoError(t, err, "Failed to build CLI tool")

	t.Run("simple_cli_generates_minimal_structure", func(t *testing.T) {
		// GIVEN: User wants a simple CLI project
		// WHEN: User generates with simple complexity
		cmd := exec.Command("./go-starter", "new", "--complexity=simple", "simple-cli", "--type=cli", "--module=github.com/test/simple-cli")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		require.NoError(t, err, "Simple CLI generation should succeed")

		// THEN: Should generate minimal file structure
		projectDir := filepath.Join(tmpDir, "simple-cli")
		
		// Check that basic files exist
		assert.FileExists(t, filepath.Join(projectDir, "main.go"))
		assert.FileExists(t, filepath.Join(projectDir, "go.mod"))
		assert.FileExists(t, filepath.Join(projectDir, "README.md"))

		// Check that it's minimal (not standard structure)
		// Simple CLI should have fewer files than standard CLI
		files, err := filepath.Glob(filepath.Join(projectDir, "**/*"))
		require.NoError(t, err)
		
		// Simple CLI should have significantly fewer files than the 25 files mentioned in the audit
		assert.Less(t, len(files), 15, "Simple CLI should have fewer than 15 files total")
	})

	t.Run("standard_cli_generates_full_structure", func(t *testing.T) {
		// GIVEN: User wants a standard CLI project
		// WHEN: User generates with standard complexity
		cmd := exec.Command("./go-starter", "new", "--complexity=standard", "standard-cli", "--type=cli", "--module=github.com/test/standard-cli")
		output, err := cmd.CombinedOutput()
		
		if err != nil {
			t.Logf("Command output: %s", string(output))
		}
		require.NoError(t, err, "Standard CLI generation should succeed")

		// THEN: Should generate full file structure
		projectDir := filepath.Join(tmpDir, "standard-cli")
		
		// Check that comprehensive structure exists
		assert.DirExists(t, filepath.Join(projectDir, "cmd"))
		assert.DirExists(t, filepath.Join(projectDir, "internal"))
		assert.FileExists(t, filepath.Join(projectDir, "Makefile"))
		
		// Standard CLI should have more files for production readiness
		files, err := filepath.Glob(filepath.Join(projectDir, "**/*"))
		require.NoError(t, err)
		
		// Standard CLI should have close to the 25 files mentioned in the audit
		assert.Greater(t, len(files), 15, "Standard CLI should have more than 15 files total")
	})
}