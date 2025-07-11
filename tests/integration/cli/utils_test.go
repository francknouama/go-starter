package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// buildTestBinary builds the CLI binary for testing
// This function creates a temporary binary that can be used across all CLI tests
func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create temporary binary in test temp directory
	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "go-starter-test")

	// Build the binary - need to build the whole package to include embed.go
	// We build from the root of the project (../../..)
	cmd := exec.Command("go", "build", "-o", binary, "../../..")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}

	// Verify binary was created and is executable
	if stat, err := os.Stat(binary); err != nil {
		t.Fatalf("Test binary was not created: %v", err)
	} else if stat.Mode().Perm()&0111 == 0 {
		// Make binary executable if it isn't already
		if err := os.Chmod(binary, 0755); err != nil {
			t.Fatalf("Failed to make test binary executable: %v", err)
		}
	}

	return binary
}

// cleanupBinary removes the test binary and handles any cleanup errors
// This function provides consistent cleanup across all test files
func cleanupBinary(t *testing.T, binary string) {
	t.Helper()

	if binary == "" {
		return
	}

	if err := os.Remove(binary); err != nil {
		// Log warning but don't fail the test for cleanup issues
		t.Logf("Warning: failed to remove test binary %s: %v", binary, err)
	}
}













