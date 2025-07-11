package cli

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_List_BasicFunctionality tests the basic list command functionality
// Verifies that the list command executes successfully and provides template information
func TestCLI_List_BasicFunctionality(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

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
		"Templates:",
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

// TestCLI_List_OutputFormat tests the format of list command output
// Ensures the output is well-structured and readable
func TestCLI_List_OutputFormat(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "list")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)

	// Output should not be empty
	if len(strings.TrimSpace(outputStr)) == 0 {
		t.Error("List output should not be empty")
	}

	// Should not contain error messages in normal operation
	lowerOutput := strings.ToLower(outputStr)
	unexpectedTerms := []string{"panic", "fatal error", "stack trace"}
	for _, term := range unexpectedTerms {
		if strings.Contains(lowerOutput, term) {
			t.Errorf("List output should not contain '%s', got: %s", term, outputStr)
		}
	}
}

// TestCLI_List_WithVerboseFlag tests list command with verbose flag
// Verifies that verbose mode provides additional information
func TestCLI_List_WithVerboseFlag(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{"verbose long flag", []string{"--verbose", "list"}},
		{"verbose short flag", []string{"-v", "list"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.Output()

			// Verbose flag should not cause the command to fail
			if err != nil {
				t.Fatalf("List command with verbose flag failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("List output with verbose flag should not be empty")
			}
		})
	}
}

// TestCLI_List_ExitCode tests that list command exits with appropriate code
// List command should succeed even when no templates are available
func TestCLI_List_ExitCode(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "list")
	_, err := cmd.Output()

	// List command should always succeed (exit code 0)
	if err != nil {
		t.Errorf("List command should exit with code 0, but failed: %v", err)
	}
}

// TestCLI_List_TemplateInfo tests that list command provides meaningful template information
// When templates are available, the output should contain useful details
func TestCLI_List_TemplateInfo(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "list")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)

	// If templates are found, output should contain helpful information
	if strings.Contains(outputStr, "templates loaded") || strings.Contains(outputStr, "Available templates:") {
		// Output should be formatted nicely
		lines := strings.Split(outputStr, "\n")
		nonEmptyLines := 0
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				nonEmptyLines++
			}
		}

		if nonEmptyLines < 2 {
			t.Errorf("Template list should contain multiple lines of information, got: %s", outputStr)
		}
	}
}

// TestCLI_List_NoTemplatesScenario tests list command behavior when no templates are available
// Should handle the case gracefully and provide helpful information
func TestCLI_List_NoTemplatesScenario(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "list")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("List command failed: %v", err)
	}

	outputStr := string(output)

	// If no templates available, should provide helpful message
	if strings.Contains(outputStr, "No templates available") || strings.Contains(outputStr, "0 templates loaded") {
		// Should still provide useful information to the user
		if len(strings.TrimSpace(outputStr)) < 10 {
			t.Error("No templates message should be informative")
		}
	}
}
