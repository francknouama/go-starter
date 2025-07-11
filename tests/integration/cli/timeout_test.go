package cli

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

// TestCLI_Timeout_NonInteractiveCommands tests that non-interactive commands complete within reasonable time
// Verifies that commands that should not require user interaction don't hang
func TestCLI_Timeout_NonInteractiveCommands(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	commands := []struct {
		name    string
		args    []string
		timeout time.Duration
	}{
		{"help", []string{"--help"}, 10 * time.Second},
		{"version", []string{"version"}, 5 * time.Second},
		{"version flag", []string{"--version"}, 5 * time.Second},
		{"list", []string{"list"}, 15 * time.Second},
		{"new help", []string{"new", "--help"}, 10 * time.Second},
		{"completion help", []string{"completion", "--help"}, 10 * time.Second},
		{"completion bash", []string{"completion", "bash"}, 10 * time.Second},
	}

	for _, tt := range commands {
		t.Run("timeout_"+tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			// Create a channel to signal command completion
			done := make(chan error, 1)

			go func() {
				_, err := cmd.CombinedOutput()
				done <- err
			}()

			// Set up timeout
			timer := time.NewTimer(tt.timeout)
			defer timer.Stop()

			select {
			case err := <-done:
				// Command completed within timeout
				if err != nil {
					// Command failure is acceptable, hanging is not
					t.Logf("Command failed (acceptable): %v", err)
				}
			case <-timer.C:
				// Command timed out
				if cmd.Process != nil {
					if killErr := cmd.Process.Kill(); killErr != nil {
						t.Logf("Warning: failed to kill timed-out process: %v", killErr)
					}
				}
				t.Errorf("Command '%s' timed out after %v", strings.Join(tt.args, " "), tt.timeout)
			}
		})
	}
}

// TestCLI_Timeout_InteractiveCommands tests timeout behavior for potentially interactive commands
// Verifies that interactive commands either complete quickly or handle timeout gracefully
func TestCLI_Timeout_InteractiveCommands(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	commands := []struct {
		name          string
		args          []string
		timeout       time.Duration
		expectTimeout bool
		description   string
	}{
		{
			name:          "new without args",
			args:          []string{"new"},
			timeout:       5 * time.Second,
			expectTimeout: true, // Might wait for input
			description:   "new command without arguments might wait for input",
		},
		{
			name:          "new with name only",
			args:          []string{"new", "--name", "test-project"},
			timeout:       10 * time.Second,
			expectTimeout: true, // Might prompt for additional info
			description:   "new command with partial args might prompt for more info",
		},
	}

	for _, tt := range commands {
		t.Run("interactive_timeout_"+tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			// Provide empty input to avoid hanging on prompts
			cmd.Stdin = strings.NewReader("")

			done := make(chan error, 1)
			var output []byte

			go func() {
				var err error
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			timer := time.NewTimer(tt.timeout)
			defer timer.Stop()

			select {
			case err := <-done:
				// Command completed
				outputStr := string(output)

				if !tt.expectTimeout {
					// Should have completed successfully
					if err != nil && !strings.Contains(outputStr, "EOF") && !strings.Contains(outputStr, "input") {
						t.Errorf("%s, but command failed unexpectedly: %v\nOutput: %s", tt.description, err, outputStr)
					}
				} else {
					// It's OK if it completed or failed, just shouldn't hang
					t.Logf("Interactive command completed (with or without error): %v", err)
				}

				// Should not panic regardless
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Command should not panic, got: %s", outputStr)
				}

			case <-timer.C:
				// Command timed out
				if cmd.Process != nil {
					if killErr := cmd.Process.Kill(); killErr != nil {
						t.Logf("Warning: failed to kill timed-out process: %v", killErr)
					}
				}

				if !tt.expectTimeout {
					t.Errorf("%s, but it timed out after %v", tt.description, tt.timeout)
				} else {
					// Expected timeout for interactive commands is acceptable
					t.Logf("Interactive command timed out as expected: %s", tt.description)
				}
			}
		})
	}
}

// TestCLI_Timeout_WithStdin tests commands with stdin input
// Verifies that commands handle stdin input without hanging
func TestCLI_Timeout_WithStdin(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		args        []string
		stdinInput  string
		timeout     time.Duration
		description string
	}{
		{
			name:        "new with empty stdin",
			args:        []string{"new"},
			stdinInput:  "",
			timeout:     5 * time.Second,
			description: "new command with empty stdin should not hang indefinitely",
		},
		{
			name:        "new with newline stdin",
			args:        []string{"new"},
			stdinInput:  "\n",
			timeout:     5 * time.Second,
			description: "new command with newline should handle input",
		},
		{
			name:        "new with multiple newlines",
			args:        []string{"new"},
			stdinInput:  "\n\n\n\n\n",
			timeout:     10 * time.Second,
			description: "new command with multiple newlines should not hang",
		},
	}

	for _, tt := range tests {
		t.Run("stdin_"+tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			cmd.Stdin = strings.NewReader(tt.stdinInput)

			done := make(chan error, 1)
			var output []byte

			go func() {
				var err error
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			timer := time.NewTimer(tt.timeout)
			defer timer.Stop()

			select {
			case err := <-done:
				// Command completed
				outputStr := string(output)

				// Any completion is acceptable, hanging is not
				t.Logf("%s - completed with error: %v", tt.description, err)

				// Should not panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Command should not panic, got: %s", outputStr)
				}

			case <-timer.C:
				// Command timed out
				if cmd.Process != nil {
					if killErr := cmd.Process.Kill(); killErr != nil {
						t.Logf("Warning: failed to kill timed-out process: %v", killErr)
					}
				}
				t.Errorf("%s, but it timed out after %v", tt.description, tt.timeout)
			}
		})
	}
}

// TestCLI_Timeout_ResourceIntensive tests potentially resource-intensive operations
// Verifies that operations that might take time still complete within reasonable bounds
func TestCLI_Timeout_ResourceIntensive(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		args        []string
		timeout     time.Duration
		description string
	}{
		{
			name:        "list with verbose",
			args:        []string{"--verbose", "list"},
			timeout:     20 * time.Second,
			description: "verbose list should complete within reasonable time",
		},
		{
			name:        "multiple flags",
			args:        []string{"--verbose", "--config", "/tmp/nonexistent.yaml", "list"},
			timeout:     15 * time.Second,
			description: "multiple flags should not cause excessive delays",
		},
	}

	for _, tt := range tests {
		t.Run("resource_"+tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			done := make(chan error, 1)
			var output []byte

			go func() {
				var err error
				output, err = cmd.CombinedOutput()
				done <- err
			}()

			timer := time.NewTimer(tt.timeout)
			defer timer.Stop()

			select {
			case err := <-done:
				// Command completed
				outputStr := string(output)

				// Command may fail, but should not hang
				t.Logf("%s - completed: %v", tt.description, err)

				// Should not panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Command should not panic, got: %s", outputStr)
				}

			case <-timer.C:
				// Command timed out
				if cmd.Process != nil {
					if killErr := cmd.Process.Kill(); killErr != nil {
						t.Logf("Warning: failed to kill timed-out process: %v", killErr)
					}
				}
				t.Errorf("%s, but it timed out after %v", tt.description, tt.timeout)
			}
		})
	}
}

// TestCLI_Timeout_GracefulShutdown tests that commands can be interrupted gracefully
// Verifies that commands respond appropriately to interruption signals
func TestCLI_Timeout_GracefulShutdown(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Test a potentially long-running command
	cmd := exec.Command(binary, "new")
	cmd.Stdin = strings.NewReader("") // Empty input might cause it to wait

	// Start the command
	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start command: %v", err)
	}

	// Let it run briefly
	time.Sleep(100 * time.Millisecond)

	// Kill the process
	if cmd.Process != nil {
		killErr := cmd.Process.Kill()
		if killErr != nil {
			t.Errorf("Failed to kill process: %v", killErr)
		}
	}

	// Wait for it to finish
	waitErr := cmd.Wait()

	// Process should have been killed (exit status will indicate this)
	if waitErr == nil {
		t.Log("Process completed before kill signal")
	} else {
		// This is expected when process is killed
		t.Logf("Process was terminated: %v", waitErr)
	}
}
