package cli

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Completion_SupportedShells tests shell completion for supported shells
// Verifies that completion scripts are generated for all supported shell types
func TestCLI_Completion_SupportedShells(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	supportedShells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range supportedShells {
		t.Run("completion_"+shell, func(t *testing.T) {
			cmd := exec.Command(binary, "completion", shell)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Completion command failed for %s: %v", shell, err)
			}

			outputStr := string(output)
			if len(outputStr) == 0 {
				t.Errorf("Completion output is empty for %s", shell)
			}

			// Basic validation that output looks like shell completion
			switch shell {
			case "bash":
				if !strings.Contains(outputStr, "complete") && !strings.Contains(outputStr, "_go-starter") {
					t.Errorf("Bash completion should contain completion functions, got: %s", outputStr)
				}
			case "zsh":
				if !strings.Contains(outputStr, "compdef") && !strings.Contains(outputStr, "_go-starter") {
					t.Errorf("Zsh completion should contain compdef, got: %s", outputStr)
				}
			case "fish":
				if !strings.Contains(outputStr, "complete") && !strings.Contains(outputStr, "go-starter") {
					t.Errorf("Fish completion should contain complete commands, got: %s", outputStr)
				}
			case "powershell":
				if !strings.Contains(outputStr, "Register-ArgumentCompleter") && !strings.Contains(outputStr, "go-starter") {
					t.Errorf("PowerShell completion should contain Register-ArgumentCompleter, got: %s", outputStr)
				}
			}
		})
	}
}

// TestCLI_Completion_UnsupportedShell tests completion with unsupported shell
// Verifies that unsupported shells are handled gracefully
func TestCLI_Completion_UnsupportedShell(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	unsupportedShells := []string{"csh", "tcsh", "ksh", "sh", "invalid-shell"}

	for _, shell := range unsupportedShells {
		t.Run("unsupported_"+shell, func(t *testing.T) {
			cmd := exec.Command(binary, "completion", shell)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Should fail for unsupported shells OR show help
			if err == nil {
				// Check if it shows help instead of failing
				if !strings.Contains(outputStr, "USAGE") && !strings.Contains(outputStr, "COMMANDS") {
					t.Errorf("Completion should fail or show help for unsupported shell %s, but got unexpected output: %s", shell, outputStr)
				}
			}
			// Should provide helpful error message
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Error message should not be empty for unsupported shell %s", shell)
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Unsupported shell should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Completion_NoArgument tests completion command without shell argument
// Verifies that missing shell argument is handled properly
func TestCLI_Completion_NoArgument(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "completion")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	// Should fail when no shell is specified OR show help
	if err == nil {
		// Check if it shows help instead of failing
		if !strings.Contains(outputStr, "USAGE") && !strings.Contains(outputStr, "COMMANDS") {
			t.Errorf("Completion without shell argument should fail or show help, but got unexpected output: %s", outputStr)
		}
	}

	// Should provide helpful error message
	if len(strings.TrimSpace(outputStr)) == 0 {
		t.Error("Error message should not be empty when shell argument is missing")
	}

	// Error message should be informative
	lowerOutput := strings.ToLower(outputStr)
	if !strings.Contains(lowerOutput, "shell") && !strings.Contains(lowerOutput, "argument") && !strings.Contains(lowerOutput, "required") {
		t.Errorf("Error message should mention missing shell argument, got: %s", outputStr)
	}

	// Should not panic
	if strings.Contains(outputStr, "panic") {
		t.Errorf("Missing shell argument should not cause panic, got: %s", outputStr)
	}
}

// TestCLI_Completion_Help tests completion command help functionality
// Verifies that help for completion command is available and useful
func TestCLI_Completion_Help(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{"completion help long", []string{"completion", "--help"}},
		{"completion help short", []string{"completion", "-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("Completion help should not fail: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)

			// Should contain help information
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("Completion help output should not be empty")
			}

			// Should mention supported shells
			expectedTerms := []string{"completion", "shell"}
			for _, term := range expectedTerms {
				if !strings.Contains(strings.ToLower(outputStr), term) {
					t.Errorf("Completion help should mention '%s', got: %s", term, outputStr)
				}
			}
		})
	}
}

// TestCLI_Completion_OutputFormat tests the format of completion output
// Verifies that completion scripts are properly formatted and executable
func TestCLI_Completion_OutputFormat(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	shells := []string{"bash", "zsh", "fish", "powershell"}

	for _, shell := range shells {
		t.Run("format_"+shell, func(t *testing.T) {
			cmd := exec.Command(binary, "completion", shell)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Completion command failed for %s: %v", shell, err)
			}

			outputStr := string(output)

			// Should not be empty
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Completion output should not be empty for %s", shell)
			}

			// Should not contain actual error messages (but bash completion may contain "error" in function names)
			lowerOutput := strings.ToLower(outputStr)
			errorTerms := []string{"failed", "panic", "fatal"}
			for _, term := range errorTerms {
				if strings.Contains(lowerOutput, term) {
					t.Errorf("Completion output should not contain error terms for %s, got: %s", shell, outputStr)
				}
			}

			// For bash, the word "error" might appear in function names, so check for actual error messages
			if shell == "bash" && strings.Contains(lowerOutput, "error:") {
				t.Errorf("Bash completion should not contain error messages, got: %s", truncateString(outputStr, 200))
			}

			// Should contain shell-specific completion syntax
			switch shell {
			case "bash":
				// Bash completions typically use 'complete' or function definitions
				if !strings.Contains(outputStr, "complete") && !strings.Contains(outputStr, "function") && !strings.Contains(outputStr, "_go-starter") {
					t.Errorf("Bash completion should contain completion syntax, got first 200 chars: %s", truncateString(outputStr, 200))
				}
			case "zsh":
				// Zsh completions use compdef and function definitions
				if !strings.Contains(outputStr, "compdef") && !strings.Contains(outputStr, "function") && !strings.Contains(outputStr, "_go-starter") {
					t.Errorf("Zsh completion should contain zsh syntax, got first 200 chars: %s", truncateString(outputStr, 200))
				}
			case "fish":
				// Fish uses 'complete' commands
				if !strings.Contains(outputStr, "complete") {
					t.Errorf("Fish completion should contain complete commands, got first 200 chars: %s", truncateString(outputStr, 200))
				}
			case "powershell":
				// PowerShell uses Register-ArgumentCompleter
				if !strings.Contains(outputStr, "Register-ArgumentCompleter") {
					t.Errorf("PowerShell completion should contain Register-ArgumentCompleter, got first 200 chars: %s", truncateString(outputStr, 200))
				}
			}
		})
	}
}

// TestCLI_Completion_ExitCodes tests exit codes for completion commands
// Verifies that completion commands use appropriate exit codes
func TestCLI_Completion_ExitCodes(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name          string
		args          []string
		expectSuccess bool
	}{
		{"valid shell", []string{"completion", "bash"}, true},
		{"invalid shell", []string{"completion", "invalid"}, true}, // Shows help instead of failing
		{"no shell argument", []string{"completion"}, true},        // Shows help instead of failing
		{"help", []string{"completion", "--help"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			_, err := cmd.CombinedOutput()

			if tt.expectSuccess && err != nil {
				t.Errorf("Command should succeed but failed: %v", err)
			}

			if !tt.expectSuccess && err == nil {
				t.Errorf("Command should fail but succeeded")
			}
		})
	}
}

// TestCLI_Completion_CommandCompletion tests that completion includes available commands
// Verifies that generated completions know about available commands
func TestCLI_Completion_CommandCompletion(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Test bash completion as it's usually the most readable
	cmd := exec.Command(binary, "completion", "bash")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Completion command failed: %v", err)
	}

	outputStr := string(output)

	// Should include available commands in completion
	expectedCommands := []string{"new", "list", "version", "completion"}

	// Check if commands are mentioned in the completion script
	// This is a basic check - real completion testing would require shell interaction
	foundCommands := 0
	for _, command := range expectedCommands {
		if strings.Contains(outputStr, command) {
			foundCommands++
		}
	}

	// At least some commands should be mentioned in the completion
	if foundCommands == 0 {
		t.Errorf("Completion script should reference available commands, got: %s", truncateString(outputStr, 300))
	}
}

// truncateString truncates a string to a maximum length for readable test output
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
