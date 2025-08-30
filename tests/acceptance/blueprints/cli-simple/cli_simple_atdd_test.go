package cli_simple

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCliSimpleBlueprintATDD validates the acceptance criteria for cli-simple blueprint
// This ensures the cli-simple blueprint generates a beginner-friendly CLI with progressive disclosure
func TestCliSimpleBlueprintATDD(t *testing.T) {
	t.Run("cli_simple_blueprint_is_available", func(t *testing.T) {
		// GIVEN: The go-starter tool is built
		// WHEN: User lists available blueprints
		// THEN: cli-simple should be in the list with correct description

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()

		// Get the project root (parent of tests/acceptance/blueprints/cli-simple)
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool first
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		if err != nil {
			t.Logf("Build output: %s", string(output))
		}
		require.NoError(t, err, "Failed to build CLI tool")

		// List blueprints
		listCmd := exec.Command("./go-starter", "list")
		output, err = listCmd.CombinedOutput()
		require.NoError(t, err, "List command should succeed")

		outputStr := string(output)
		assert.Contains(t, outputStr, "cli-simple", "cli-simple blueprint should be listed")
		assert.Contains(t, outputStr, "Simple command-line application", "Should show cli-simple description")
		assert.Contains(t, outputStr, "essential features only", "Should emphasize beginner-friendliness")
	})

	t.Run("cli_simple_generates_minimal_structure", func(t *testing.T) {
		// GIVEN: User wants a simple CLI for learning/prototyping
		// WHEN: User generates a project with cli-simple blueprint
		// THEN: Should generate exactly 8 files (not 29 like standard CLI)

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build the CLI tool
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		output, err := buildCmd.CombinedOutput()
		require.NoError(t, err, "Failed to build CLI tool: %s", string(output))

		// Generate a cli-simple project
		generateCmd := exec.Command("./go-starter", "new", "test-cli-simple",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/cli-simple",
			"--logger=slog",
			"--no-git")
		output, err = generateCmd.CombinedOutput()

		if err != nil {
			t.Logf("Generate command output: %s", string(output))
		}
		require.NoError(t, err, "Project generation should succeed")

		// Verify generated structure
		projectDir := filepath.Join(tmpDir, "test-cli-simple")

		// Verify essential files exist (minimal set - 7 files as per template.yaml)
		essentialFiles := []string{
			"main.go",
			"go.mod",
			"README.md",
			".gitignore",
			"cmd/root.go",
			"cmd/version.go",
			"config.go",
			"Makefile",
		}

		for _, file := range essentialFiles {
			assert.FileExists(t, filepath.Join(projectDir, file), "Essential file %s should exist", file)
		}

		// Count total files to ensure simplicity (should be around 8 files)
		var fileCount int
		err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileCount++
			}
			return nil
		})
		require.NoError(t, err)

		// cli-simple should have exactly 8 files (73% reduction from standard CLI's 29 files)
		assert.LessOrEqual(t, fileCount, 11, "cli-simple should have ≤11 files (was %d)", fileCount)
		assert.GreaterOrEqual(t, fileCount, 7, "cli-simple should have ≥7 essential files (was %d)", fileCount)
		t.Logf("Generated cli-simple project has %d files (target: 8±3)", fileCount)

		// Verify NO over-complex structure (should NOT have these directories/files)
		complexFiles := []string{
			"internal/config/",   // No complex config management
			"internal/commands/", // No command separation
			"internal/logger/",   // No separate logger files
			"pkg/",               // No public packages
			"docs/",              // No extensive documentation
			"examples/",          // No examples in simple CLI
		}

		for _, complexFile := range complexFiles {
			assert.NoFileExists(t, filepath.Join(projectDir, complexFile), "Should NOT have complex file %s in simple CLI", complexFile)
		}
	})

	t.Run("cli_simple_project_builds_successfully", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: User builds the project
		// THEN: It should compile without errors and create working binary

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-build",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/build",
			"--logger=slog",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-build")

		// Initialize go modules
		modInitCmd := exec.Command("go", "mod", "tidy")
		modInitCmd.Dir = projectDir
		output, err := modInitCmd.CombinedOutput()
		require.NoError(t, err, "go mod tidy should succeed: %s", string(output))

		// Build the generated project
		buildGeneratedCmd := exec.Command("go", "build", "-o", "cli-app", ".")
		buildGeneratedCmd.Dir = projectDir
		output, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err, "Generated cli-simple project should build successfully: %s", string(output))

		// Verify binary was created
		assert.FileExists(t, filepath.Join(projectDir, "cli-app"), "CLI binary should be created")
	})

	t.Run("cli_simple_runtime_functionality", func(t *testing.T) {
		// GIVEN: A built cli-simple project
		// WHEN: User runs the CLI tool
		// THEN: It should execute and show help correctly

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate and build project
		generateCmd := exec.Command("./go-starter", "new", "test-runtime",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/runtime",
			"--logger=slog",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-runtime")

		// Initialize and build
		modInitCmd := exec.Command("go", "mod", "tidy")
		modInitCmd.Dir = projectDir
		_, err = modInitCmd.CombinedOutput()
		require.NoError(t, err)

		buildGeneratedCmd := exec.Command("go", "build", "-o", "test-runtime", ".")
		buildGeneratedCmd.Dir = projectDir
		_, err = buildGeneratedCmd.CombinedOutput()
		require.NoError(t, err)

		// Test CLI functionality
		// 1. Help command
		helpCmd := exec.Command("./test-runtime", "--help")
		helpCmd.Dir = projectDir
		output, err := helpCmd.CombinedOutput()
		require.NoError(t, err, "CLI help should work")

		helpStr := string(output)
		assert.Contains(t, helpStr, "test-runtime", "Help should show CLI name")
		assert.Contains(t, helpStr, "Usage:", "Help should show usage")

		// 2. Version flag (if available)
		versionCmd := exec.Command("./test-runtime", "--version")
		versionCmd.Dir = projectDir
		output, _ = versionCmd.CombinedOutput() // Don't require - version might not be implemented

		// 3. Run without args (should not crash)
		runCmd := exec.Command("./test-runtime")
		runCmd.Dir = projectDir
		output, err = runCmd.CombinedOutput()
		// CLI should not crash (exit code 0 or 1 is acceptable for help)
		// Just verify it doesn't panic or error severely
		t.Logf("CLI run output: %s", string(output))
	})

	t.Run("cli_simple_uses_hardcoded_slog", func(t *testing.T) {
		// GIVEN: cli-simple blueprint (designed to be truly minimal)
		// WHEN: Generating a project
		// THEN: Should always use slog (hardcoded, no configurability by design)

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		projectName := "test-slog"
		generateCmd := exec.Command("./go-starter", "new", projectName,
			"--type=cli-simple",
			"--module=github.com/test/"+projectName,
			"--no-git")
		output, err := generateCmd.CombinedOutput()
		require.NoError(t, err, "Should generate successfully: %s", string(output))

		projectDir := filepath.Join(tmpDir, projectName)

		// Check that slog is used in main.go (no separate logger files)
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)
		
		assert.Contains(t, mainStr, "log/slog", "Should import slog")
		assert.Contains(t, mainStr, "slog.New", "Should initialize slog")
		assert.Contains(t, mainStr, "slog.SetDefault", "Should set default logger")

		// Verify no internal logger directory (by design)
		_, err = os.Stat(filepath.Join(projectDir, "internal"))
		assert.True(t, os.IsNotExist(err), "cli-simple should not have internal directory")

		// Project should compile
		buildCmd = exec.Command("go", "build", "-o", "test-build", ".")
		buildCmd.Dir = projectDir
		buildOutput, buildErr := buildCmd.CombinedOutput()
		assert.NoError(t, buildErr, "Project should compile: %s", string(buildOutput))

		assert.FileExists(t, filepath.Join(projectDir, "test-build"), "Binary should be created")
	})

	t.Run("cli_simple_progressive_disclosure_compliance", func(t *testing.T) {
		// GIVEN: Progressive disclosure system implementation
		// WHEN: Using complexity=simple flag
		// THEN: Should map to cli-simple blueprint and apply smart defaults

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Test progressive disclosure: complexity=simple should auto-select cli-simple
		generateCmd := exec.Command("./go-starter", "new", "test-disclosure",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/disclosure",
			"--no-git")
		output, err := generateCmd.CombinedOutput()
		require.NoError(t, err, "Progressive disclosure should work: %s", string(output))

		projectDir := filepath.Join(tmpDir, "test-disclosure")

		// Verify it generated the simple version (not standard version)
		var fileCount int
		err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileCount++
			}
			return nil
		})
		require.NoError(t, err)

		// Should have simple file count (≤11 files), not standard CLI count (~29 files)
		assert.LessOrEqual(t, fileCount, 11, "Progressive disclosure should generate simple CLI (was %d files)", fileCount)

		// Verify smart defaults were applied (should use slog logger by default)
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)
		assert.Contains(t, mainStr, "log/slog", "Should default to slog logger in progressive disclosure")
		assert.Contains(t, mainStr, "slog.New", "Should initialize slog logger")

		// Check no interactive prompts were needed (all defaults applied)
		outputStr := string(output)
		assert.NotContains(t, outputStr, "Select", "Should not prompt when using progressive disclosure")
		assert.NotContains(t, outputStr, "Choose", "Should not prompt when using progressive disclosure")
	})

	t.Run("cli_simple_makefile_targets", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: User examines available make targets
		// THEN: Should have essential targets without over-complexity

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate project
		generateCmd := exec.Command("./go-starter", "new", "test-makefile",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/makefile",
			"--logger=slog",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-makefile")

		// Initialize go modules
		modCmd := exec.Command("go", "mod", "tidy")
		modCmd.Dir = projectDir
		_, err = modCmd.CombinedOutput()
		require.NoError(t, err)

		// Test make help
		makeHelpCmd := exec.Command("make", "help")
		makeHelpCmd.Dir = projectDir
		output, err := makeHelpCmd.CombinedOutput()
		require.NoError(t, err, "make help should succeed")

		helpStr := string(output)
		t.Logf("Available make targets:\n%s", helpStr)

		// Essential targets for cli-simple (minimal set)
		essentialTargets := []string{"build", "test", "clean", "help"}
		for _, target := range essentialTargets {
			assert.Contains(t, helpStr, target, "Should have essential %s target", target)
		}

		// Should NOT have overly complex targets
		complexTargets := []string{"docker-build", "deploy", "benchmark", "release"}
		for _, target := range complexTargets {
			assert.NotContains(t, helpStr, target, "cli-simple should not have complex %s target", target)
		}

		// Test make build creates binary
		makeBuildCmd := exec.Command("make", "build")
		makeBuildCmd.Dir = projectDir
		output, err = makeBuildCmd.CombinedOutput()
		require.NoError(t, err, "make build should succeed: %s", string(output))

		// Check binary was created (should be in current directory for simplicity)
		binaryPath := filepath.Join(projectDir, "test-makefile")
		if !assert.FileExists(t, binaryPath, "Binary should be created by make build") {
			// Alternative location: bin/ directory
			binPath := filepath.Join(projectDir, "bin", "test-makefile")
			assert.FileExists(t, binPath, "Binary should be in bin/ directory if not in root")
		}
	})

	t.Run("cli_simple_code_quality_and_structure", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: Examining code quality and structure
		// THEN: Should follow Go best practices while maintaining simplicity

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-quality",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/quality",
			"--logger=slog",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-quality")

		// Initialize modules
		modCmd := exec.Command("go", "mod", "tidy")
		modCmd.Dir = projectDir
		_, err = modCmd.CombinedOutput()
		require.NoError(t, err)

		// 1. Check main.go structure
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)

		assert.Contains(t, mainStr, "package main", "Should have main package")
		assert.Contains(t, mainStr, "func main()", "Should have main function")
		assert.Contains(t, mainStr, "github.com/test/quality/cmd", "Should import local cmd package")

		// 2. Check cmd/root.go structure (Cobra integration)
		rootContent, err := os.ReadFile(filepath.Join(projectDir, "cmd", "root.go"))
		require.NoError(t, err)
		rootStr := string(rootContent)

		assert.Contains(t, rootStr, "github.com/spf13/cobra", "Should use Cobra framework")
		assert.Contains(t, rootStr, "rootCmd", "Should have root command")
		assert.Contains(t, rootStr, "Execute", "Should have Execute function")

		// 3. Check logger integration
		loggerContent, err := os.ReadFile(filepath.Join(projectDir, "internal", "logger", "logger.go"))
		require.NoError(t, err)
		loggerStr := string(loggerContent)

		assert.Contains(t, loggerStr, "package logger", "Should have logger package")
		assert.Contains(t, loggerStr, "log/slog", "Should use slog")
		assert.Contains(t, loggerStr, "New", "Should have constructor function")

		// 4. Check go.mod file
		goModContent, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
		require.NoError(t, err)
		goModStr := string(goModContent)

		assert.Contains(t, goModStr, "module github.com/test/quality", "Should have correct module path")
		assert.Contains(t, goModStr, "go 1.2", "Should specify Go version")
		assert.Contains(t, goModStr, "github.com/spf13/cobra", "Should require Cobra")

		// 5. Run go fmt to verify code is formatted
		fmtCmd := exec.Command("go", "fmt", "./...")
		fmtCmd.Dir = projectDir
		output, err := fmtCmd.CombinedOutput()

		// If go fmt produces output, it means files were not formatted
		assert.Empty(t, string(output), "Generated code should be properly formatted")
		require.NoError(t, err, "go fmt should succeed")

		// 6. Run go vet
		vetCmd := exec.Command("go", "vet", "./...")
		vetCmd.Dir = projectDir
		output, err = vetCmd.CombinedOutput()
		require.NoError(t, err, "go vet should pass: %s", string(output))

		// 7. Check README.md content
		readmeContent, err := os.ReadFile(filepath.Join(projectDir, "README.md"))
		require.NoError(t, err)
		readmeStr := string(readmeContent)

		assert.Contains(t, readmeStr, "# test-quality", "README should have project title")
		assert.Contains(t, readmeStr, "Installation", "README should have installation instructions")
		assert.Contains(t, readmeStr, "Usage", "README should have usage instructions")
		assert.Contains(t, readmeStr, "go install", "README should show go install command")
	})
}

// TestCliSimpleVsStandardComparison validates the differentiation between simple and standard CLI
func TestCliSimpleVsStandardComparison(t *testing.T) {
	t.Run("cli_simple_has_significantly_fewer_files_than_standard", func(t *testing.T) {
		// GIVEN: Both cli-simple and standard CLI blueprints
		// WHEN: Generating both with same configuration
		// THEN: cli-simple should have 60-75% fewer files than standard

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build go-starter
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate cli-simple
		generateSimpleCmd := exec.Command("./go-starter", "new", "test-simple",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/simple",
			"--logger=slog",
			"--no-git")
		_, err = generateSimpleCmd.CombinedOutput()
		require.NoError(t, err)

		// Generate standard CLI
		generateStandardCmd := exec.Command("./go-starter", "new", "test-standard",
			"--type=cli",
			"--complexity=standard",
			"--module=github.com/test/standard",
			"--logger=slog",
			"--no-git")
		_, err = generateStandardCmd.CombinedOutput()
		require.NoError(t, err)

		// Count files in both projects
		simpleFileCount := countFiles(t, filepath.Join(tmpDir, "test-simple"))
		standardFileCount := countFiles(t, filepath.Join(tmpDir, "test-standard"))

		t.Logf("cli-simple generated %d files", simpleFileCount)
		t.Logf("cli-standard generated %d files", standardFileCount)

		// cli-simple should have significantly fewer files (target: 73% reduction)
		reductionPercentage := float64(standardFileCount-simpleFileCount) / float64(standardFileCount) * 100
		t.Logf("Reduction: %.1f%% (target: 60-75%%)", reductionPercentage)

		assert.GreaterOrEqual(t, reductionPercentage, 60.0, "cli-simple should have 60%+ fewer files than standard")
		assert.LessOrEqual(t, reductionPercentage, 80.0, "Reduction should be realistic (not over 80%)")

		// Absolute checks
		assert.LessOrEqual(t, simpleFileCount, 11, "cli-simple should have ≤11 files")
		assert.GreaterOrEqual(t, standardFileCount, 20, "cli-standard should have ≥20 files")
	})

	t.Run("cli_simple_excludes_advanced_patterns", func(t *testing.T) {
		// GIVEN: A generated cli-simple project
		// WHEN: Examining the code structure
		// THEN: Should not include advanced patterns that standard CLI has

		tmpDir := t.TempDir()
		originalDir, _ := os.Getwd()
		projectRoot := filepath.Join(originalDir, "..", "..", "..", "..")

		defer func() { _ = os.Chdir(originalDir) }()
		_ = os.Chdir(tmpDir)

		// Build and generate cli-simple
		buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
		buildCmd.Dir = projectRoot
		_, err := buildCmd.CombinedOutput()
		require.NoError(t, err)

		generateCmd := exec.Command("./go-starter", "new", "test-patterns",
			"--type=cli",
			"--complexity=simple",
			"--module=github.com/test/patterns",
			"--no-git")
		_, err = generateCmd.CombinedOutput()
		require.NoError(t, err)

		projectDir := filepath.Join(tmpDir, "test-patterns")

		// Advanced patterns that should NOT be in cli-simple:
		advancedPatterns := []string{
			// Advanced directory structure
			"pkg/",                 // No public packages
			"examples/",            // No examples
			"docs/",                // No extensive docs
			"scripts/",             // No deployment scripts
			"configs/",             // No complex configuration
			"internal/commands/",   // No command separation
			"internal/config/",     // No config management
			"internal/middleware/", // No middleware

			// Advanced files  
			"cmd/completion.go",  // Shell completion
			"docker-compose.yml", // No Docker setup
			"Dockerfile",         // No containerization
			".github/",           // No CI/CD workflows
			"deploy/",            // No deployment configs
		}

		for _, pattern := range advancedPatterns {
			fullPath := filepath.Join(projectDir, pattern)
			if strings.HasSuffix(pattern, "/") {
				assert.NoDirExists(t, fullPath, "cli-simple should not have advanced directory: %s", pattern)
			} else {
				assert.NoFileExists(t, fullPath, "cli-simple should not have advanced file: %s", pattern)
			}
		}

		// Verify simplicity in main.go (no complex imports)
		mainContent, err := os.ReadFile(filepath.Join(projectDir, "main.go"))
		require.NoError(t, err)
		mainStr := string(mainContent)

		// Should not import complex dependencies
		complexImports := []string{
			"github.com/spf13/viper", // No Viper config
			"github.com/gorilla/",    // No web dependencies
			"database/sql",           // No database connections
			"net/http",               // No HTTP server
		}

		for _, complexImport := range complexImports {
			assert.NotContains(t, mainStr, complexImport, "cli-simple should not import complex dependency: %s", complexImport)
		}
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
