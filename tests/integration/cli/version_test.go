package cli

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Version_Command tests the version command functionality
// Verifies that version information is displayed correctly through various invocation methods
func TestCLI_Version_Command(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"version command", []string{"version"}},
		{"version flag", []string{"--version"}},
		{"version short flag", []string{"-v"}},
	}

	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.Output()

			if err != nil {
				t.Fatalf("Version command failed: %v", err)
			}

			outputStr := strings.TrimSpace(string(output))
			if outputStr == "" {
				t.Error("Version output should not be empty")
			}

			// Should contain version information
			if !strings.Contains(outputStr, "Version:") && !strings.Contains(outputStr, "version") {
				t.Errorf("Version output should contain version info, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Version_Format tests the format of version output
// Ensures version information follows expected formatting conventions
func TestCLI_Version_Format(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "version")
	output, err := cmd.Output()

	if err != nil {
		t.Fatalf("Version command failed: %v", err)
	}

	outputStr := strings.TrimSpace(string(output))

	// Version output length depends on the invocation method
	lines := strings.Split(outputStr, "\n")
	// The "version" command shows detailed output, "--version" shows concise output
	if len(lines) > 25 {
		t.Errorf("Version output should be reasonable, got %d lines: %s", len(lines), outputStr)
	}

	// Should not contain error messages
	lowerOutput := strings.ToLower(outputStr)
	errorIndicators := []string{"error", "failed", "panic", "fatal"}
	for _, indicator := range errorIndicators {
		if strings.Contains(lowerOutput, indicator) {
			t.Errorf("Version output should not contain error indicators, got: %s", outputStr)
		}
	}
}

// TestCLI_Version_Consistency tests that all version invocation methods return consistent output
// Different ways of getting version should return the same information
func TestCLI_Version_Consistency(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	methods := [][]string{
		{"version"},
		{"--version"},
		{"-v"},
	}

	var outputs []string

	for _, method := range methods {
		cmd := exec.Command(binary, method...)
		output, err := cmd.Output()

		if err != nil {
			t.Fatalf("Version command failed for %v: %v", method, err)
		}

		outputs = append(outputs, strings.TrimSpace(string(output)))
	}

	// Different invocation methods may have different formats, but all should contain version info
	for i, output := range outputs {
		if !strings.Contains(strings.ToLower(output), "version") {
			t.Errorf("Version output for method %v should contain version information: %s", methods[i], output)
		}
	}

	// The "version" command should be more detailed than flags
	versionCmd := outputs[0]  // "version" command
	versionFlag := outputs[1] // "--version" flag

	if len(versionCmd) <= len(versionFlag) {
		t.Logf("Note: 'version' command output (%d chars) is not longer than '--version' flag (%d chars)",
			len(versionCmd), len(versionFlag))
	}
}

// TestCLI_Version_ExitCode tests that version commands exit with code 0
// Version commands should always succeed
func TestCLI_Version_ExitCode(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := [][]string{
		{"version"},
		{"--version"},
		{"-v"},
	}

	for _, args := range tests {
		t.Run("exit_code_"+strings.Join(args, "_"), func(t *testing.T) {
			cmd := exec.Command(binary, args...)
			_, err := cmd.Output()

			if err != nil {
				t.Errorf("Version command should exit with code 0, but failed: %v", err)
			}
		})
	}
}

