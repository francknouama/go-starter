package cli

import (
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Help_RootCommand tests the root help command functionality
// Verifies that the main help command displays expected content including
// tool name, examples section, and supported templates section
func TestCLI_Help_RootCommand(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "full help flag",
			args:     []string{"--help"},
			contains: []string{"go-starter", "EXAMPLES:", "SUPPORTED TEMPLATES:"},
		},
		{
			name:     "short help flag",
			args:     []string{"-h"},
			contains: []string{"go-starter"},
		},
	}

	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("Help command failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Help output does not contain '%s'\nOutput: %s", contain, outputStr)
				}
			}
		})
	}
}

// TestCLI_Help_SubCommands tests help functionality for individual subcommands
// Ensures each subcommand provides appropriate help information
func TestCLI_Help_SubCommands(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "new command help",
			args:     []string{"new", "--help"},
			contains: []string{"new", "project"},
		},
		{
			name:     "list command help",
			args:     []string{"list", "--help"},
			contains: []string{"list", "templates"},
		},
		{
			name:     "version command help",
			args:     []string{"version", "--help"},
			contains: []string{"version"},
		},
	}

	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			if err != nil {
				t.Fatalf("Subcommand help failed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Subcommand help does not contain '%s'\nOutput: %s", contain, outputStr)
				}
			}
		})
	}
}

// TestCLI_Help_CommandStructure tests the overall help command structure
// Verifies that help commands provide consistent formatting and expected sections
func TestCLI_Help_CommandStructure(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	cmd := exec.Command(binary, "--help")
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Fatalf("Help command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)

	// Test for consistent help structure
	expectedSections := []string{
		"USAGE",
		"COMMANDS",
		"FLAGS",
	}

	for _, section := range expectedSections {
		if !strings.Contains(outputStr, section) {
			t.Errorf("Help output should contain '%s' section\nOutput: %s", section, outputStr)
		}
	}

	// Test that help output is not empty
	if len(strings.TrimSpace(outputStr)) == 0 {
		t.Error("Help output should not be empty")
	}
}

// TestCLI_Help_ExitCodes tests that help commands exit with appropriate codes
// Help commands should always exit with code 0 (success)
func TestCLI_Help_ExitCodes(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{"root help", []string{"--help"}},
		{"short help", []string{"-h"}},
		{"new help", []string{"new", "--help"}},
		{"list help", []string{"list", "--help"}},
		{"version help", []string{"version", "--help"}},
	}

	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			_, err := cmd.CombinedOutput()

			// Help commands should always succeed (exit code 0)
			if err != nil {
				t.Errorf("Help command should exit with code 0, but failed: %v", err)
			}
		})
	}
}
