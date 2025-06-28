package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestCLIGoVersionFlag tests the --go-version flag for the new command.
func TestCLIGoVersionFlag(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	tests := []struct {
		name        string
		args        []string
		shouldFail  bool
		expectedMsg string
		checkGoMod  bool
		goVersion   string
	}{
		{
			name:       "valid go version",
			args:       []string{"new", "test-project", "--type=cli", "--go-version=1.23", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail: false,
			checkGoMod: true,
			goVersion:  "1.23",
		},
		{
			name:        "invalid go version",
			args:        []string{"new", "test-project", "--type=cli", "--go-version=1.10", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail:  true,
			expectedMsg: "invalid Go version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			cmd := exec.Command(binary, tt.args...)
			cmd.Dir = tmpDir
			// Prepare input for interactive prompts
			input := strings.NewReader(`github.com/test/test-project
Cobra (recommended)
slog - Go built-in structured logging (recommended)
`)
			cmd.Stdin = input

			output, err := cmd.CombinedOutput()

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Expected command to fail, but it succeeded.")
				}
				if !strings.Contains(string(output), tt.expectedMsg) {
					t.Errorf("Expected output to contain '%s', but got '%s'", tt.expectedMsg, string(output))
				}
			} else {
				if err != nil {
					t.Errorf("Expected command to succeed, but it failed with error: %v. Output: %s", err, string(output))
				}

				if tt.checkGoMod {
					goModPath := filepath.Join(tmpDir, "test-project", "go.mod")
					content, err := os.ReadFile(goModPath)
					if err != nil {
						t.Fatalf("Failed to read go.mod file: %v", err)
					}
					expectedGoVersion := "go " + tt.goVersion
					if !strings.Contains(string(content), expectedGoVersion) {
						t.Errorf("go.mod does not contain the correct go version. Expected '%s', Content: %s", expectedGoVersion, string(content))
					}
				}
			}
		})
	}
}
