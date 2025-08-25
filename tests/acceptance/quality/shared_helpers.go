package quality

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// buildGoStarter builds the go-starter CLI tool once for all tests
func buildGoStarter(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	buildCmd := exec.Command("go", "build", "-o", filepath.Join(tmpDir, "go-starter"), ".")
	buildCmd.Dir = projectRoot
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Logf("Build output: %s", string(output))
		t.Logf("Project root: %s", projectRoot)
	}
	require.NoError(t, err, "Failed to build CLI tool")
	t.Logf("Successfully built go-starter CLI tool")
}

// countProjectFiles counts all files in a project directory
func countProjectFiles(t *testing.T, projectDir string) int {
	t.Helper()
	var count int
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	require.NoError(t, err, "Should be able to count files")
	return count
}

// findLoggerFiles finds logger-related files in a project
func findLoggerFiles(t *testing.T, projectDir string) []string {
	t.Helper()
	var loggerFiles []string
	
	err := filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Look for logger-related files
		if strings.Contains(path, "logger") && strings.HasSuffix(path, ".go") {
			loggerFiles = append(loggerFiles, path)
		}
		
		return nil
	})
	
	require.NoError(t, err, "Should be able to find logger files")
	return loggerFiles
}

// initializeGoModules runs go mod tidy to initialize modules
func initializeGoModules(t *testing.T, projectDir string) {
	t.Helper()
	
	modTidyCmd := exec.Command("go", "mod", "tidy")
	modTidyCmd.Dir = projectDir
	output, err := modTidyCmd.CombinedOutput()
	
	require.NoError(t, err, "go mod tidy should succeed: %s", string(output))
}

// validateProjectCompilation ensures the generated project compiles successfully
func validateProjectCompilation(t *testing.T, projectDir, blueprintName string) {
	t.Helper()
	
	// Initialize go modules first
	initializeGoModules(t, projectDir)
	
	// Determine the correct build path based on project type
	buildPath := "."
	binaryName := "test-binary"
	
	// Web API projects have main.go in cmd/server/
	if strings.Contains(blueprintName, "web-api") {
		buildPath = "./cmd/server"
	}
	
	// Build the project
	buildCmd := exec.Command("go", "build", "-o", binaryName, buildPath)
	buildCmd.Dir = projectDir
	output, err := buildCmd.CombinedOutput()
	
	require.NoError(t, err, "Project %s should compile successfully: %s", blueprintName, string(output))
	
	// Verify binary was created
	binaryPath := filepath.Join(projectDir, binaryName)
	require.FileExists(t, binaryPath, "Binary should be created for %s", blueprintName)
	
	t.Logf("✓ %s compiles successfully", blueprintName)
}

// isBinaryFile checks if a file is binary (to avoid reading it)
func isBinaryFile(path string) bool {
	// Simple heuristic to avoid reading binary files
	ext := strings.ToLower(filepath.Ext(path))
	binaryExts := []string{".exe", ".bin", ".so", ".dylib", ".dll", ".png", ".jpg", ".jpeg", ".gif", ".pdf"}
	
	for _, bExt := range binaryExts {
		if ext == bExt {
			return true
		}
	}
	
	return false
}