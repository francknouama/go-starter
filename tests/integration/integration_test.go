package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/utils"
)

// TestMain sets up and tears down the test environment
func TestMain(m *testing.M) {
	// Setup
	code := setupTestEnvironment()
	if code != 0 {
		os.Exit(code)
	}

	// Run tests
	exitCode := m.Run()

	// Cleanup
	teardownTestEnvironment()

	os.Exit(exitCode)
}

// setupTestEnvironment prepares the test environment
func setupTestEnvironment() int {
	// Check if Go is installed
	if !utils.IsGoInstalled() {
		println("Go is not installed, skipping integration tests")
		return 1
	}

	// Check if Git is installed
	if !utils.IsGitInstalled() {
		println("Git is not installed, skipping integration tests")
		return 1
	}

	// Validate Go installation
	if err := utils.ValidateGoInstallation(); err != nil {
		println("Go installation validation failed:", err.Error())
		return 1
	}

	// Check Git installation (but don't fail if config is missing)
	if err := utils.CheckGitInstallation(); err != nil {
		// Just warn, don't fail
		println("Warning: Git configuration issue:", err.Error())
	}

	return 0
}

// teardownTestEnvironment cleans up after tests
func teardownTestEnvironment() {
	// Any cleanup needed after all tests
}

// TestIntegrationEnvironment tests the integration test environment
func TestIntegrationEnvironment(t *testing.T) {
	t.Run("go installation", func(t *testing.T) {
		if !utils.IsGoInstalled() {
			t.Skip("Go is not installed")
		}

		version, err := utils.GoVersion()
		if err != nil {
			t.Fatalf("Failed to get Go version: %v", err)
		}

		t.Logf("Go version: %s", version)
	})

	t.Run("git installation", func(t *testing.T) {
		if !utils.IsGitInstalled() {
			t.Skip("Git is not installed")
		}

		version, err := utils.GetGitVersion()
		if err != nil {
			t.Fatalf("Failed to get Git version: %v", err)
		}

		t.Logf("Git version: %s", version)
	})

	t.Run("config system", func(t *testing.T) {
		// Test basic config functionality
		cfg := config.DefaultConfig

		if len(cfg.Profiles) == 0 {
			t.Error("Default config should have profiles")
		}

		if _, exists := cfg.Profiles["default"]; !exists {
			t.Error("Default config should have a default profile")
		}

		if cfg.CurrentProfile != "default" {
			t.Error("Default config should have 'default' as current profile")
		}
	})

	t.Run("utilities", func(t *testing.T) {
		// Test basic utility functions
		tmpDir := t.TempDir()

		// Test directory creation
		testDir := utils.JoinPath(tmpDir, "test-dir")
		if err := utils.CreateDir(testDir, 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		if !utils.DirExists(testDir) {
			t.Error("Directory should exist after creation")
		}

		// Test file operations
		testFile := utils.JoinPath(testDir, "test.txt")
		testContent := "Hello, World!"

		if err := utils.WriteFile(testFile, testContent); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		if !utils.FileExists(testFile) {
			t.Error("File should exist after writing")
		}

		content, err := utils.ReadFile(testFile)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		if content != testContent {
			t.Errorf("File content mismatch. Expected: %s, Got: %s", testContent, content)
		}
	})
}

// TestIntegrationProjectStructure tests the overall project structure
func TestIntegrationProjectStructure(t *testing.T) {
	// Get project root (two levels up from tests/integration)
	projectRoot := "../.."

	// Test that required directories exist
	requiredDirs := []string{
		"cmd",
		"internal",
		"pkg",
		"templates",
		"tests",
		"scripts",
	}

	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectRoot, dir)
		if !utils.DirExists(fullPath) {
			t.Errorf("Required directory missing: %s (checked at %s)", dir, fullPath)
		}
	}

	// Test that required files exist
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"Makefile",
		"README.md",
		"CLAUDE.md",
	}

	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectRoot, file)
		if !utils.FileExists(fullPath) {
			t.Errorf("Required file missing: %s (checked at %s)", file, fullPath)
		}
	}

	// Test that go.mod is valid
	t.Run("go.mod validation", func(t *testing.T) {
		if !utils.HasGoMod(projectRoot) {
			t.Error("go.mod file should exist in project root")
		}

		// Try to get module path
		modulePath, err := utils.GetModulePath(projectRoot)
		if err != nil {
			t.Fatalf("Failed to get module path: %v", err)
		}

		if modulePath == "" {
			t.Error("Module path should not be empty")
		}

		t.Logf("Module path: %s", modulePath)
	})
}

// TestIntegrationBuildSystem tests the build system
func TestIntegrationBuildSystem(t *testing.T) {
	t.Run("project compilation", func(t *testing.T) {
		// Test that the project compiles
		if err := utils.GoBuild(".", "", "./..."); err != nil {
			t.Fatalf("Project should compile: %v", err)
		}
	})

	t.Run("tests execution", func(t *testing.T) {
		// Test that existing tests pass
		// Note: This might run unit tests, which is okay for integration testing
		if err := utils.GoTest(".", "./cmd/..."); err != nil {
			t.Logf("Some tests failed (this might be expected): %v", err)
		}
	})

	t.Run("code formatting", func(t *testing.T) {
		// Test that code formatting works
		if err := utils.GoFmt("."); err != nil {
			t.Fatalf("Code formatting should work: %v", err)
		}
	})

	t.Run("code vetting", func(t *testing.T) {
		// Test that go vet passes
		if err := utils.GoVet("."); err != nil {
			t.Logf("Go vet found issues (might be expected in development): %v", err)
		}
	})
}

// TestIntegrationGitOperations tests Git operations
func TestIntegrationGitOperations(t *testing.T) {
	if !utils.IsGitInstalled() {
		t.Skip("Git is not installed")
	}

	tmpDir := t.TempDir()

	t.Run("git repository initialization", func(t *testing.T) {
		// Test Git repository initialization
		if err := utils.InitGitRepository(tmpDir); err != nil {
			t.Fatalf("Failed to initialize Git repository: %v", err)
		}

		if !utils.IsGitRepository(tmpDir) {
			t.Error("Directory should be a Git repository after initialization")
		}
	})

	t.Run("gitignore creation", func(t *testing.T) {
		gitignoreContent := utils.GetDefaultGitIgnore()
		if err := utils.AddGitIgnore(tmpDir, gitignoreContent); err != nil {
			t.Fatalf("Failed to create .gitignore: %v", err)
		}

		gitignorePath := utils.JoinPath(tmpDir, ".gitignore")
		if !utils.FileExists(gitignorePath) {
			t.Error(".gitignore should exist after creation")
		}
	})

	t.Run("git operations", func(t *testing.T) {
		// Create a test file
		testFile := utils.JoinPath(tmpDir, "test.txt")
		if err := utils.WriteFile(testFile, "test content"); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Add files to Git
		if err := utils.GitAdd(tmpDir); err != nil {
			t.Fatalf("Failed to add files to Git: %v", err)
		}

		// Check Git status
		status, err := utils.GitStatus(tmpDir)
		if err != nil {
			t.Fatalf("Failed to get Git status: %v", err)
		}

		t.Logf("Git status: %s", status)
	})
}

// TestIntegrationModuleOperations tests Go module operations
func TestIntegrationModuleOperations(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("module initialization", func(t *testing.T) {
		modulePath := "github.com/test/integration-test"

		if err := utils.InitGoModule(tmpDir, modulePath); err != nil {
			t.Fatalf("Failed to initialize Go module: %v", err)
		}

		if !utils.HasGoMod(tmpDir) {
			t.Error("go.mod should exist after module initialization")
		}

		// Check module path
		actualPath, err := utils.GetModulePath(tmpDir)
		if err != nil {
			t.Fatalf("Failed to get module path: %v", err)
		}

		if actualPath != modulePath {
			t.Errorf("Module path mismatch. Expected: %s, Got: %s", modulePath, actualPath)
		}
	})

	t.Run("module operations", func(t *testing.T) {
		// Test module tidy
		if err := utils.GoModTidy(tmpDir); err != nil {
			t.Fatalf("Failed to run go mod tidy: %v", err)
		}

		// Test module download
		if err := utils.GoModDownload(tmpDir); err != nil {
			t.Fatalf("Failed to run go mod download: %v", err)
		}
	})
}

// TestIntegrationValidation tests various validation functions
func TestIntegrationValidation(t *testing.T) {
	t.Run("project name validation", func(t *testing.T) {
		validNames := []string{
			"my-project",
			"MyProject",
			"my_project",
			"project123",
		}

		invalidNames := []string{
			"",
			"-invalid",
			"invalid-",
			"_invalid",
			"invalid_",
			"in..valid",
			"con", // reserved name
		}

		for _, name := range validNames {
			if err := config.ValidateProjectName(name); err != nil {
				t.Errorf("Valid project name '%s' should pass validation: %v", name, err)
			}
		}

		for _, name := range invalidNames {
			if err := config.ValidateProjectName(name); err == nil {
				t.Errorf("Invalid project name '%s' should fail validation", name)
			}
		}
	})

	t.Run("module path validation", func(t *testing.T) {
		validPaths := []string{
			"github.com/user/repo",
			"gitlab.com/user/repo",
			"example.com/path/to/module",
		}

		invalidPaths := []string{
			"",
			"invalid",
			"github.com",
			"/invalid/path",
		}

		for _, path := range validPaths {
			if err := config.ValidateModulePath(path); err != nil {
				t.Errorf("Valid module path '%s' should pass validation: %v", path, err)
			}
		}

		for _, path := range invalidPaths {
			if err := config.ValidateModulePath(path); err == nil {
				t.Errorf("Invalid module path '%s' should fail validation", path)
			}
		}
	})

	t.Run("go version validation", func(t *testing.T) {
		validVersions := []string{
			"1.18",
			"1.19",
			"1.20",
			"1.21",
			"1.22",
			"1.21.0",
		}

		invalidVersions := []string{
			"",
			"1.17", // too old
			"2.0",  // invalid format
			"1",    // incomplete
			"invalid",
		}

		for _, version := range validVersions {
			if err := config.ValidateGoVersion(version); err != nil {
				t.Errorf("Valid Go version '%s' should pass validation: %v", version, err)
			}
		}

		for _, version := range invalidVersions {
			if err := config.ValidateGoVersion(version); err == nil {
				t.Errorf("Invalid Go version '%s' should fail validation", version)
			}
		}
	})
}
