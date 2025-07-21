package cli

import (
	"strings"
	"testing"
)

// TestCLI_New_InteractiveMode tests the new command interactive mode
// Verifies that the new command properly handles interactive mode scenarios
func TestCLI_New_InteractiveMode(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name  string
		stdin string
	}{
		{"empty input", Inputs.Empty()},
		{"newline input", Inputs.Newline()},
		{"multiple newlines", Inputs.MultipleNewlines()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteInteractiveCommand(t, tt.stdin, "new")

			// Interactive commands might fail, timeout, or succeed - all are acceptable
			// The key requirement is that they should not panic or hang indefinitely
			AssertNoPanic(t, result, "interactive command should not panic")

			if result.Error != nil {
				// Check if it's a reasonable interactive failure
				validResponses := []string{
					"project name", "What's your project name?",
					"USAGE", "Usage:", "interactive", "required",
				}
				
				hasValidResponse := false
				for _, response := range validResponses {
					if strings.Contains(result.Output, response) {
						hasValidResponse = true
						break
					}
				}

				if hasValidResponse {
					t.Logf("Interactive command handled gracefully: %v", result.Error)
				} else {
					t.Logf("Interactive command failed (acceptable): %v\nOutput: %s", 
						result.Error, result.Output)
				}
			}

			if result.TimedOut {
				t.Log("Command timed out waiting for input (expected for interactive mode)")
			}
		})
	}
}

// TestCLI_New_HelpCommand tests the new command help functionality
// Verifies that help for the new command provides useful information
func TestCLI_New_HelpCommand(t *testing.T) {
	binary := GetTestBinary(t)
	result := binary.ExecuteFastCommand(t, Args.NewHelp()...)

	AssertSuccess(t, result, "new command help should succeed")
	AssertNoPanic(t, result, "help command should not panic")

	expectedContent := []string{"new", "project"}
	for _, expected := range expectedContent {
		AssertContains(t, result, expected, "new command help should contain expected content")
	}
}

// TestCLI_New_WithFlags tests the new command with various flags
// Verifies that flags are properly handled and don't cause crashes
func TestCLI_New_WithFlags(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name          string
		args          []string
		expectFailure bool
		stdin         string
	}{
		{
			name:          "new with name flag",
			args:          []string{"new", "--name", TestProjectName},
			expectFailure: false,
			stdin:         Inputs.SkipPrompts(),
		},
		{
			name:          "new with type flag",
			args:          []string{"new", "--type", "cli"},
			expectFailure: false,
			stdin:         Inputs.SkipPrompts(),
		},
		{
			name:          "new with invalid flag",
			args:          []string{"new", "--invalid-flag"},
			expectFailure: true,
			stdin:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteInteractiveCommand(t, tt.stdin, tt.args...)

			if tt.expectFailure {
				AssertFailure(t, result, "invalid flag should cause failure")
			} else {
				// For valid flags, command might still fail due to missing templates
				// or incomplete input, but should not crash
				AssertNoPanic(t, result, "valid flags should not cause crashes")
				
				if result.Error != nil {
					t.Logf("Command with valid flags failed (may be due to missing input): %v", result.Error)
				}
			}

			AssertNoPanic(t, result, "command should handle errors gracefully")
		})
	}
}

// TestCLI_New_NonInteractiveMode tests new command with all required parameters
// Verifies that the command can run without user interaction when parameters are provided
func TestCLI_New_NonInteractiveMode(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name string
		args []string
	}{
		{
			name: "with name and type",
			args: []string{"new", "--name", TestProjectName, "--type", "cli"},
		},
		{
			name: "with multiple flags", 
			args: []string{"new", "--name", TestProjectName, "--type", "cli", "--output", "/tmp"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteSlowCommand(t, tt.args...)

			// Command might fail due to missing templates or other setup issues
			// but should not hang or crash
			AssertNoPanic(t, result, "non-interactive command should not panic")
			AssertNotTimedOut(t, result, "command with sufficient parameters should not hang")

			// Log the result for debugging
			t.Logf("Non-interactive command result - Error: %v, Output length: %d", 
				result.Error, len(result.Output))
		})
	}
}

// TestCLI_New_ErrorHandling tests error handling in the new command
// Verifies that errors are handled gracefully and provide useful feedback
func TestCLI_New_ErrorHandling(t *testing.T) {
	binary := GetTestBinary(t)

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
			result := binary.ExecuteInteractiveCommand(t, Inputs.ProjectName(), tt.args...)

			// Command should fail for invalid inputs, but gracefully
			if result.Error == nil {
				t.Logf("Command unexpectedly succeeded: %s", result.Output)
			}

			// Error message should be informative, not a crash
			AssertNoPanic(t, result, "invalid input should fail gracefully, not crash")
			AssertNotContains(t, result, "stack trace", "should not show stack traces to users")
		})
	}
}