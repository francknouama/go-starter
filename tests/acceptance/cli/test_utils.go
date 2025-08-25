package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// findProjectRoot locates the project root by finding go.mod
func findProjectRoot(t *testing.T) string {
	currentDir, err := os.Getwd()
	require.NoError(t, err, "Should get current directory")

	// Start from current directory and walk up
	dir := currentDir
	for {
		goModPath := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(goModPath); err == nil {
			return dir
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Reached root directory
		}
		dir = parent
	}

	// Fallback: assume we're in tests/acceptance/cli and go up 3 levels
	projectRoot := filepath.Join(currentDir, "..", "..", "..")
	goModPath := filepath.Join(projectRoot, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		return projectRoot
	}

	t.Fatal("Could not find project root (go.mod not found)")
	return ""
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// getBinaryName returns the appropriate binary name for the current platform
func getBinaryName(baseName string) string {
	if filepath.Ext(baseName) != "" {
		return baseName // Already has extension
	}
	
	// Add .exe extension on Windows
	if isWindows() {
		return baseName + ".exe"
	}
	return baseName
}

// isWindows checks if running on Windows
func isWindows() bool {
	return filepath.Separator == '\\'
}

// getProjectTypes returns common project types for testing
func getProjectTypes() []string {
	return []string{"cli", "web-api", "library", "lambda"}
}

// getComplexityLevels returns supported complexity levels
func getComplexityLevels() []string {
	return []string{"simple", "standard", "advanced", "expert"}
}

// getLoggerTypes returns supported logger types
func getLoggerTypes() []string {
	return []string{"slog", "zap", "logrus", "zerolog"}
}

// buildBinary builds a binary with proper directory handling for different binary types
func buildBinary(t *testing.T, projectRoot, outputPath, binaryName, buildPath string) error {
	var cmd *exec.Cmd
	
	if binaryName == "go-starter-web" {
		// Web server has its own go.mod, build from web directory
		cmd = exec.Command("go", "build", "-o", outputPath, "./cmd/web-server")
		cmd.Dir = filepath.Join(projectRoot, "web")
	} else {
		// Other binaries build from project root
		cmd = exec.Command("go", "build", "-o", outputPath, buildPath)
		cmd.Dir = projectRoot
	}
	
	return cmd.Run()
}