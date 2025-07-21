package cli

import (
	"testing"
)

// TestCLI_Help_RootCommand tests the root help command functionality
// Verifies that the main help command displays expected content including
// tool name, examples section, and supported blueprints section
func TestCLI_Help_RootCommand(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "full help flag",
			args:     Args.Help(),
			contains: []string{"go-starter", "EXAMPLES:", "SUPPORTED BLUEPRINTS:"},
		},
		{
			name:     "short help flag",
			args:     []string{"-h"},
			contains: []string{"go-starter"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			AssertSuccess(t, result, "help command should succeed")
			AssertNoPanic(t, result, "help command should not panic")

			for _, expected := range tt.contains {
				AssertContains(t, result, expected, "help output should contain expected content")
			}
		})
	}
}

// TestCLI_Help_SubCommands tests help functionality for individual subcommands
// Ensures each subcommand provides appropriate help information
func TestCLI_Help_SubCommands(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name     string
		args     []string
		contains []string
	}{
		{
			name:     "new command help",
			args:     Args.NewHelp(),
			contains: []string{"new", "project"},
		},
		{
			name:     "list command help",
			args:     []string{"list", "--help"},
			contains: []string{"list", "blueprints"},
		},
		{
			name:     "version command help",
			args:     []string{"version", "--help"},
			contains: []string{"version"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			AssertSuccess(t, result, "subcommand help should succeed")
			AssertNoPanic(t, result, "help command should not panic")

			for _, expected := range tt.contains {
				AssertContains(t, result, expected, "subcommand help should contain expected content")
			}
		})
	}
}

// TestCLI_Help_CommandStructure tests the overall help command structure
// Verifies that help commands provide consistent formatting and expected sections
func TestCLI_Help_CommandStructure(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name     string
		args     []string
		sections []string
	}{
		{
			name:     "root help structure",
			args:     Args.Help(),
			sections: []string{"USAGE", "COMMANDS", "FLAGS"},
		},
		{
			name:     "new command structure", 
			args:     Args.NewHelp(),
			sections: []string{"USAGE", "FLAGS"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			AssertSuccess(t, result, "help command should succeed")
			AssertNoPanic(t, result, "help command should not panic")

			for _, section := range tt.sections {
				AssertContains(t, result, section, "help should contain standard sections")
			}
		})
	}
}

// TestCLI_Help_ErrorCases tests help command error scenarios
// Verifies graceful handling of invalid help requests
func TestCLI_Help_ErrorCases(t *testing.T) {
	binary := GetTestBinary(t)

	tests := []struct {
		name string
		args []string
	}{
		{"invalid command help", []string{"invalid-command", "--help"}},
		{"malformed help flag", []string{"---help"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := binary.ExecuteFastCommand(t, tt.args...)
			
			// These should fail gracefully
			AssertFailure(t, result, "invalid help command should fail")
			AssertNoPanic(t, result, "invalid command should not panic")
		})
	}
}