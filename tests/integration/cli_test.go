package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestCLIHelp tests the CLI help command
func TestCLIHelp(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "root help",
			args:     []string{"--help"},
			contains: []string{"go-starter", "Available Commands", "Flags"},
		},
		{
			name:     "short help",
			args:     []string{"-h"},
			contains: []string{"go-starter"},
		},
		{
			name:     "new command help",
			args:     []string{"new", "--help"},
			contains: []string{"Create a new Go project", "Usage"},
		},
		{
			name:     "list command help",
			args:     []string{"list", "--help"},
			contains: []string{"Display all available project templates", "Usage"},
		},
		{
			name:     "version command help",
			args:     []string{"version", "--help"},
			contains: []string{"Display version information", "Usage"},
		},
	}

	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			// Help commands should exit with code 0
			if err != nil {
				t.Fatalf("Command failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Output does not contain '%s'\nOutput: %s", contain, outputStr)
				}
			}
		})
	}
}

// TestCLIVersion tests the version command
func TestCLIVersion(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	tests := []struct {
		name string
		args []string
	}{
		{"version command", []string{"version"}},
		{"version flag", []string{"--version"}},
		{"version short flag", []string{"-v"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Command failed: %v", err)
			}

			outputStr := strings.TrimSpace(string(output))
			if outputStr == "" {
				t.Error("Version output is empty")
			}

			// Should contain version information
			if !strings.Contains(outputStr, "Version:") && !strings.Contains(outputStr, "version") {
				t.Errorf("Version output should contain version info, got: %s", outputStr)
			}
		})
	}
}

// TestCLIList tests the list command
func TestCLIList(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	cmd := exec.Command(binary, "list")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)

	// Should contain some indication of available templates
	// Even if no templates are loaded, it should show a message
	expectedMessages := []string{
		"Available templates:",
		"No templates available",
		"templates loaded",
	}

	found := false
	for _, msg := range expectedMessages {
		if strings.Contains(outputStr, msg) {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("List output should contain template information, got: %s", outputStr)
	}
}

// TestCLINewCommand tests the new command basic functionality
func TestCLINewCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	// Test new command without arguments
	cmd := exec.Command(binary, "new")
	output, err := cmd.CombinedOutput()

	// Should fail because project name is required
	if err == nil {
		t.Error("New command should fail without project name")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "project name") && !strings.Contains(outputStr, "Usage") {
		t.Errorf("Error message should mention project name or show usage, got: %s", outputStr)
	}
}

// TestCLIInvalidCommand tests invalid command handling
func TestCLIInvalidCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	cmd := exec.Command(binary, "invalid-command")
	output, err := cmd.CombinedOutput()

	// Should fail with non-zero exit code
	if err == nil {
		t.Error("Invalid command should fail")
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "unknown command") && !strings.Contains(outputStr, "invalid-command") {
		t.Errorf("Error message should mention unknown/invalid command, got: %s", outputStr)
	}
}

// TestCLIFlags tests global flags
func TestCLIFlags(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	tests := []struct {
		name       string
		args       []string
		shouldPass bool
		contains   []string
	}{
		{
			name:       "verbose flag",
			args:       []string{"--verbose", "list"},
			shouldPass: true,
			contains:   []string{}, // Just test it doesn't crash
		},
		{
			name:       "config flag",
			args:       []string{"--config", "/tmp/nonexistent.yaml", "list"},
			shouldPass: true, // Should handle missing config gracefully
			contains:   []string{},
		},
		{
			name:       "invalid flag",
			args:       []string{"--invalid-flag"},
			shouldPass: false,
			contains:   []string{"unknown flag"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.shouldPass && err != nil {
				t.Errorf("Command should pass but failed: %v\nOutput: %s", err, output)
			}

			if !tt.shouldPass && err == nil {
				t.Errorf("Command should fail but passed. Output: %s", output)
			}

			outputStr := string(output)
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Output should contain '%s', got: %s", contain, outputStr)
				}
			}
		})
	}
}

// TestCLICompletion tests shell completion
func TestCLICompletion(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run("completion "+shell, func(t *testing.T) {
			cmd := exec.Command(binary, "completion", shell)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Completion command failed for %s: %v", shell, err)
			}

			outputStr := string(output)
			if len(outputStr) == 0 {
				t.Errorf("Completion output is empty for %s", shell)
			}

			// Basic check that it looks like shell completion
			if shell == "bash" && !strings.Contains(outputStr, "complete") {
				t.Errorf("Bash completion should contain 'complete', got: %s", outputStr)
			}
		})
	}
}

// buildTestBinary builds the CLI binary for testing
func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create temporary binary
	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "go-starter-test")

	// Build the binary - need to build the whole package to include embed.go
	cmd := exec.Command("go", "build", "-o", binary, "../..")
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}

	// Verify binary was created
	if _, err := os.Stat(binary); err != nil {
		t.Fatalf("Test binary was not created: %v", err)
	}

	return binary
}

// TestCLITimeout tests that CLI commands don't hang
func TestCLITimeout(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	commands := [][]string{
		{"--help"},
		{"version"},
		{"list"},
		{"new", "--help"},
	}

	for _, args := range commands {
		t.Run("timeout "+strings.Join(args, " "), func(t *testing.T) {
			cmd := exec.Command(binary, args...)

			// Set a reasonable timeout
			timer := time.AfterFunc(30*time.Second, func() {
				if cmd.Process != nil {
					if err := cmd.Process.Kill(); err != nil {
						t.Logf("Warning: failed to kill process: %v", err)
					}
				}
			})
			defer timer.Stop()

			_, err := cmd.CombinedOutput()

			// We don't care about the exit code here, just that it doesn't hang
			select {
			case <-timer.C:
				t.Error("Command timed out")
			default:
				// Command completed within timeout
			}

			// Check if process was killed due to timeout
			if err != nil {
				if strings.Contains(err.Error(), "killed") {
					t.Error("Command was killed due to timeout")
				}
			}
		})
	}
}

// TestCLIEnvironment tests CLI behavior with different environment variables
func TestCLIEnvironment(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	tests := []struct {
		name     string
		env      map[string]string
		args     []string
		contains []string
	}{
		{
			name: "custom config via env",
			env:  map[string]string{"GO_STARTER_CONFIG": "/tmp/nonexistent.yaml"},
			args: []string{"list"},
			// Should handle missing config gracefully
		},
		{
			name: "verbose via env",
			env:  map[string]string{"GO_STARTER_VERBOSE": "true"},
			args: []string{"list"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			// Set environment variables
			env := os.Environ()
			for k, v := range tt.env {
				env = append(env, k+"="+v)
			}
			cmd.Env = env

			output, err := cmd.CombinedOutput()

			// Commands should not fail due to environment variables
			if err != nil {
				// Only fail if it's not a graceful handling of missing config
				outputStr := string(output)
				if !strings.Contains(outputStr, "config") && !strings.Contains(outputStr, "not found") {
					t.Errorf("Command failed unexpectedly: %v\nOutput: %s", err, output)
				}
			}

			outputStr := string(output)
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Output should contain '%s', got: %s", contain, outputStr)
				}
			}
		})
	}
}
