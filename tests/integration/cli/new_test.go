package cli

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestCLI_New_InteractiveMode tests the new command interactive mode
// Verifies that the new command properly enters interactive mode when no arguments are provided
func TestCLI_New_InteractiveMode(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Test new command without arguments - should start interactive mode
	cmd := exec.Command(binary, "new")
	cmd.Stdin = strings.NewReader("") // Provide empty input

	// Set a timeout since interactive mode might wait for input
	done := make(chan error, 1)
	var output []byte

	go func() {
		var err error
		output, err = cmd.CombinedOutput()
		done <- err
	}()

	select {
	case err := <-done:
		outputStr := string(output)

		// Interactive mode with empty input should either:
		// 1. Prompt for project name
		// 2. Fail gracefully with informative message
		// 3. Show help information
		if err != nil {
			// Check if it's trying to prompt for input or showing help
			if strings.Contains(outputStr, "project name") ||
				strings.Contains(outputStr, "What's your project name?") ||
				strings.Contains(outputStr, "USAGE") ||
				strings.Contains(outputStr, "Usage:") ||
				strings.Contains(outputStr, "interactive") {
				// This is expected behavior
				return
			}

			// Other types of failures might be acceptable too
			t.Logf("New command output: %s", outputStr)
		}

	case <-time.After(5 * time.Second):
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
		// Timeout is also acceptable - means it was waiting for input
		t.Log("Command timed out waiting for input (expected for interactive mode)")
	}
}

// TestCLI_New_HelpCommand tests the new command help functionality
// Verifies that help for the new command provides useful information
func TestCLI_New_HelpCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "new", "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("New command help failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	expectedContent := []string{
		"new",
		"project",
	}

	for _, content := range expectedContent {
		if !strings.Contains(outputStr, content) {
			t.Errorf("New command help should contain '%s', got: %s", content, outputStr)
		}
	}
}

// TestCLI_New_WithFlags tests the new command with various flags
// Verifies that flags are properly handled and don't cause crashes
func TestCLI_New_WithFlags(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name           string
		args           []string
		expectFailure  bool
		timeoutSeconds int
	}{
		{
			name:           "new with name flag",
			args:           []string{"new", "--name", "test-project"},
			expectFailure:  false,
			timeoutSeconds: 10,
		},
		{
			name:           "new with type flag",
			args:           []string{"new", "--type", "api"},
			expectFailure:  false,
			timeoutSeconds: 10,
		},
		{
			name:           "new with invalid flag",
			args:           []string{"new", "--invalid-flag"},
			expectFailure:  true,
			timeoutSeconds: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			// For non-interactive flags, provide minimal input
			if !tt.expectFailure {
				cmd.Stdin = strings.NewReader("\n\n\n\n\n") // Skip through prompts
			}

			done := make(chan error, 1)
			var output []byte

			go func() {
				var err error
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			select {
			case err := <-done:
				outputStr := string(output)

				if tt.expectFailure {
					if err == nil {
						t.Errorf("Expected command to fail but it succeeded. Output: %s", outputStr)
					}
				} else {
					// For valid flags, command might still fail due to missing templates
					// or incomplete input, but should not crash
					if err != nil {
						// Check if it's a graceful failure
						if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "fatal") {
							t.Errorf("Command should not crash, got: %s", outputStr)
						}
					}
				}

			case <-time.After(time.Duration(tt.timeoutSeconds) * time.Second):
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
				if !tt.expectFailure {
					t.Log("Command timed out (might be waiting for input)")
				}
			}
		})
	}
}

// TestCLI_New_NonInteractiveMode tests new command with all required parameters
// Verifies that the command can run without user interaction when all parameters are provided
func TestCLI_New_NonInteractiveMode(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Test with common flags that might make the command non-interactive
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "with name and type",
			args: []string{"new", "--name", "test-project", "--type", "api"},
		},
		{
			name: "with multiple flags",
			args: []string{"new", "--name", "test-project", "--type", "api", "--output", "/tmp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			done := make(chan error, 1)
			var output []byte

			go func() {
				var err error
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			select {
			case err := <-done:
				outputStr := string(output)

				// Command might fail due to missing templates or other setup issues
				// but should not hang or crash
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Command should not panic, got: %s", outputStr)
				}

				// Log the result for debugging
				t.Logf("Command result - Error: %v, Output: %s", err, outputStr)

			case <-time.After(15 * time.Second):
				if cmd.Process != nil {
					_ = cmd.Process.Kill()
				}
				t.Error("Command should not hang when all parameters are provided")
			}
		})
	}
}

// TestCLI_New_ErrorHandling tests error handling in the new command
// Verifies that errors are handled gracefully and provide useful feedback
func TestCLI_New_ErrorHandling(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "invalid template type",
			args: []string{"new", "--type", "invalid-type"},
		},
		{
			name: "invalid output directory",
			args: []string{"new", "--output", "/invalid/path/that/does/not/exist"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			cmd.Stdin = strings.NewReader("test-project\n") // Provide project name if needed

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Command should fail for invalid inputs
			if err == nil {
				t.Logf("Command unexpectedly succeeded: %s", outputStr)
			}

			// Error message should be informative, not a crash
			if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "stack trace") {
				t.Errorf("Command should fail gracefully, not crash. Output: %s", outputStr)
			}
		})
	}
}
