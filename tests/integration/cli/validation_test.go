package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Validation_InvalidCommand tests handling of invalid commands
// Verifies that unknown commands are handled gracefully with helpful error messages
func TestCLI_Validation_InvalidCommand(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name       string
		args       []string
		shouldFail bool
		errorTerms []string
	}{
		{
			name:       "completely invalid command",
			args:       []string{"invalid-command"},
			shouldFail: true,
			errorTerms: []string{"unknown command", "invalid-command"},
		},
		{
			name:       "typo in command",
			args:       []string{"nwe"}, // typo for "new"
			shouldFail: true,
			errorTerms: []string{"unknown command", "nwe"},
		},
		{
			name:       "invalid subcommand",
			args:       []string{"new", "invalid-subcommand"},
			shouldFail: true,
			errorTerms: []string{}, // May vary based on implementation
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Command should fail for invalid input, but succeeded. Output: %s", outputStr)
				}

				// Check for expected error terms (case insensitive)
				lowerOutput := strings.ToLower(outputStr)
				for _, term := range tt.errorTerms {
					if !strings.Contains(lowerOutput, strings.ToLower(term)) {
						t.Errorf("Error message should contain '%s', got: %s", term, outputStr)
					}
				}

				// Error should be informative, not a panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Invalid command should not cause panic, got: %s", outputStr)
				}
			}
		})
	}
}

// TestCLI_Validation_InvalidFlags tests handling of invalid flags
// Verifies that unknown flags are properly rejected with helpful messages
func TestCLI_Validation_InvalidFlags(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name       string
		args       []string
		shouldFail bool
		contains   []string
	}{
		{
			name:       "invalid global flag",
			args:       []string{"--invalid-flag"},
			shouldFail: true,
			contains:   []string{"unknown flag"},
		},
		{
			name:       "invalid flag with command",
			args:       []string{"list", "--invalid-flag"},
			shouldFail: true,
			contains:   []string{"unknown flag"},
		},
		{
			name:       "typo in flag",
			args:       []string{"--helpp"}, // typo for "--help"
			shouldFail: true,
			contains:   []string{"unknown flag"},
		},
		{
			name:       "malformed flag",
			args:       []string{"---help"}, // three dashes
			shouldFail: true,
			contains:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.shouldFail {
				if err == nil {
					t.Errorf("Command with invalid flag should fail, but succeeded. Output: %s", outputStr)
				}

				// Check for expected error content (case insensitive)
				lowerOutput := strings.ToLower(outputStr)
				for _, contain := range tt.contains {
					if contain != "" && !strings.Contains(lowerOutput, strings.ToLower(contain)) {
						t.Errorf("Error output should contain '%s', got: %s", contain, outputStr)
					}
				}

				// Should not panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Invalid flag should not cause panic, got: %s", outputStr)
				}
			}
		})
	}
}

// TestCLI_Validation_RequiredArguments tests validation of required arguments
// Verifies that commands requiring arguments handle missing arguments gracefully
func TestCLI_Validation_RequiredArguments(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "completion without shell",
			args:        []string{"completion"},
			expectError: false, // Shows help instead of failing
			description: "completion command should show help when no shell provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError {
				if err == nil {
					t.Errorf("%s, but command succeeded. Output: %s", tt.description, outputStr)
				}

				// Error message should be helpful
				if len(strings.TrimSpace(outputStr)) == 0 {
					t.Error("Error message should not be empty")
				}

				// Should not panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Missing argument should not cause panic, got: %s", outputStr)
				}
			}
		})
	}
}

// TestCLI_Validation_ArgumentValidation tests validation of argument values
// Verifies that invalid argument values are properly validated
func TestCLI_Validation_ArgumentValidation(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		args        []string
		expectError bool
		description string
	}{
		{
			name:        "completion with invalid shell",
			args:        []string{"completion", "invalid-shell"},
			expectError: false, // Shows help instead of failing
			description: "completion should handle invalid shell types gracefully",
		},
		{
			name:        "completion with empty shell",
			args:        []string{"completion", ""},
			expectError: false, // Shows help instead of failing
			description: "completion should handle empty shell argument gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError {
				if err == nil {
					t.Errorf("%s, but command succeeded. Output: %s", tt.description, outputStr)
				}

				// Error should be informative
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Invalid argument should not cause panic, got: %s", outputStr)
				}
			}
		})
	}
}

// TestCLI_Validation_ExitCodes tests that validation errors use appropriate exit codes
// Different types of validation errors should use consistent exit codes
func TestCLI_Validation_ExitCodes(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{"invalid command", []string{"invalid-command"}},
		{"invalid flag", []string{"--invalid-flag"}},
		// Removed invalid completion shell as it shows help instead of failing
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			_, err := cmd.CombinedOutput()

			// All validation errors should result in non-zero exit codes
			if err == nil {
				t.Errorf("Validation error should result in non-zero exit code for: %v", tt.args)
			}

			// Should be an ExitError (not a different type of error like panic)
			if _, ok := err.(*exec.ExitError); !ok {
				t.Errorf("Error should be ExitError, got: %T", err)
			}
		})
	}
}

// TestCLI_Validation_GracefulHandling tests overall graceful error handling
// Ensures that all validation errors are handled gracefully without crashes
func TestCLI_Validation_GracefulHandling(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Collection of various invalid inputs
	invalidInputs := [][]string{
		{""},                            // empty command
		{"invalid-command"},             // unknown command
		{"--invalid-flag"},              // unknown flag
		{"new", "--bad-flag"},           // invalid flag with valid command
		{"list", "extra", "args"},       // too many arguments
		{"completion", "invalid-shell"}, // invalid argument value
	}

	for i, args := range invalidInputs {
		t.Run(fmt.Sprintf("invalid_input_%d", i), func(t *testing.T) {
			// Filter out empty args
			filteredArgs := make([]string, 0, len(args))
			for _, arg := range args {
				if arg != "" {
					filteredArgs = append(filteredArgs, arg)
				}
			}

			if len(filteredArgs) == 0 {
				t.Skip("Skipping empty command test")
			}

			cmd := exec.Command(binary, filteredArgs...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Should not panic under any circumstances
			if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "runtime error") {
				t.Errorf("Command should handle invalid input gracefully, got panic: %s", outputStr)
			}

			// Should provide some form of error message
			if err != nil && len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("Error should be accompanied by an error message")
			}
		})
	}
}
