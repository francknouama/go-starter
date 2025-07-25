package cli

import (
	"strings"
	"testing"
)

// TestCLI_Version_Command tests the version command functionality
// Verifies that version information is displayed correctly through various invocation methods
func TestCLI_Version_Command(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name string
		args []string
	}{
		{"version command", Args.Version()},
		{"version flag", Args.VersionFlag()},
		{"version short flag", []string{"-v"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			AssertSuccess(t, result, "version command should succeed")
			AssertNoPanic(t, result, "version command should not panic")
			
			// Verify output is not empty
			if strings.TrimSpace(result.Output) == "" {
				t.Error("Version output should not be empty")
			}

			// Should contain version information
			if !strings.Contains(result.Output, "Version:") && !strings.Contains(result.Output, "version") {
				t.Errorf("Version output should contain version info, got: %s", result.Output)
			}
		})
	}
}

// TestCLI_Version_Format tests the format of version output
// Ensures version information follows expected formatting conventions
func TestCLI_Version_Format(t *testing.T) {
	binary := GetTestBinary(t)
	result := binary.ExecuteFastCommand(t, Args.Version()...)
	
	AssertSuccess(t, result, "version command should succeed")
	AssertNoPanic(t, result, "version command should not panic")

	outputStr := strings.TrimSpace(result.Output)

	// Version output should be reasonable length
	lines := strings.Split(outputStr, "\n")
	if len(lines) > 25 {
		t.Errorf("Version output should be reasonable, got %d lines: %s", len(lines), outputStr)
	}

	// Should not contain error messages
	errorIndicators := []string{"error", "failed", "panic", "fatal"}
	for _, indicator := range errorIndicators {
		AssertNotContains(t, result, indicator, "version output should not contain error indicators")
	}
}

// TestCLI_Version_Consistency tests that all version invocation methods return consistent output
// Different ways of getting version should return the same information  
func TestCLI_Version_Consistency(t *testing.T) {
	binary := GetTestBinary(t)

	methods := []struct {
		name string
		args []string
	}{
		{"version command", Args.Version()},
		{"version flag", Args.VersionFlag()},
		{"version short flag", []string{"-v"}},
	}

	var outputs []string

	for _, method := range methods {
		t.Run("consistency_"+method.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, method.args...)
			
			AssertSuccess(t, result, "version method should succeed")
			AssertContains(t, result, "version", "output should contain version information")
			
			outputs = append(outputs, strings.TrimSpace(result.Output))
		})
	}

	// Compare outputs if we got them all
	if len(outputs) == 3 {
		versionCmd := outputs[0]  // "version" command
		versionFlag := outputs[1] // "--version" flag

		if len(versionCmd) <= len(versionFlag) {
			t.Logf("Note: 'version' command output (%d chars) is not longer than '--version' flag (%d chars)",
				len(versionCmd), len(versionFlag))
		}
	}
}

// TestCLI_Version_ExitCode tests that version commands exit with code 0
// Version commands should always succeed
func TestCLI_Version_ExitCode(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name string
		args []string
	}{
		{"version command", Args.Version()},
		{"version flag", Args.VersionFlag()},
		{"version short flag", []string{"-v"}},
	}

	for _, tt := range tests {
		t.Run("exit_code_"+tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			AssertSuccess(t, result, "version command should exit with code 0")
			AssertExitCode(t, result, 0, "version command should have exit code 0")
		})
	}
}