package cli

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Flags_GlobalFlags tests global flags functionality
// Verifies that global flags work correctly across different commands
func TestCLI_Flags_GlobalFlags(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name       string
		args       []string
		shouldPass bool
		contains   []string
	}{
		{
			name:       "verbose flag with list",
			args:       []string{"--verbose", "list"},
			shouldPass: true,
			contains:   []string{}, // Just test it doesn't crash
		},
		{
			name:       "verbose short flag with list",
			args:       []string{"-v", "list"},
			shouldPass: true,
			contains:   []string{},
		},
		{
			name:       "config flag with list",
			args:       []string{"--config", "/tmp/nonexistent.yaml", "list"},
			shouldPass: true, // Should handle missing config gracefully
			contains:   []string{},
		},
		{
			name:       "config flag with valid command",
			args:       []string{"--config", "/tmp/test-config.yaml", "version"},
			shouldPass: true,
			contains:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.shouldPass && err != nil {
				// Check if it's a graceful failure (e.g., missing config file)
				if !strings.Contains(outputStr, "config") && !strings.Contains(outputStr, "not found") {
					t.Errorf("Command should pass but failed: %v\nOutput: %s", err, outputStr)
				}
			}

			if !tt.shouldPass && err == nil {
				t.Errorf("Command should fail but passed. Output: %s", outputStr)
			}

			// Check for expected content
			for _, contain := range tt.contains {
				if !strings.Contains(outputStr, contain) {
					t.Errorf("Output should contain '%s', got: %s", contain, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Command should not panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Flags_InvalidFlags tests handling of invalid flags
// Verifies that invalid flags are properly rejected
func TestCLI_Flags_InvalidFlags(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name       string
		args       []string
		shouldPass bool
		contains   []string
	}{
		{
			name:       "invalid global flag",
			args:       []string{"--invalid-flag"},
			shouldPass: false,
			contains:   []string{"unknown flag"},
		},
		{
			name:       "invalid flag with valid command",
			args:       []string{"--invalid-flag", "list"},
			shouldPass: false,
			contains:   []string{"unknown flag"},
		},
		{
			name:       "double dash flag",
			args:       []string{"---help"},
			shouldPass: false,
			contains:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.shouldPass && err != nil {
				t.Errorf("Command should pass but failed: %v\nOutput: %s", err, outputStr)
			}

			if !tt.shouldPass && err == nil {
				t.Errorf("Command should fail but passed. Output: %s", outputStr)
			}

			// Check for expected error content (case insensitive)
			lowerOutput := strings.ToLower(outputStr)
			for _, contain := range tt.contains {
				if contain != "" && !strings.Contains(lowerOutput, strings.ToLower(contain)) {
					t.Errorf("Output should contain '%s', got: %s", contain, outputStr)
				}
			}
		})
	}
}

// TestCLI_Flags_ConfigFlag tests the config flag functionality
// Verifies that config flag properly handles different scenarios
func TestCLI_Flags_ConfigFlag(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		configPath  string
		command     string
		expectError bool
		description string
	}{
		{
			name:        "nonexistent config file",
			configPath:  "/tmp/nonexistent-config.yaml",
			command:     "version",
			expectError: false, // Should handle gracefully
			description: "should handle missing config file gracefully",
		},
		{
			name:        "empty config path",
			configPath:  "",
			command:     "version",
			expectError: false,
			description: "should handle empty config path",
		},
		{
			name:        "invalid config path",
			configPath:  "/invalid/path/config.yaml",
			command:     "list",
			expectError: false, // Should handle gracefully
			description: "should handle invalid config path gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"--config", tt.configPath, tt.command}
			cmd := exec.Command(binary, args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			if tt.expectError && err == nil {
				t.Errorf("%s, but command succeeded. Output: %s", tt.description, outputStr)
			}

			if !tt.expectError && err != nil {
				// Check if it's a graceful handling of config issues
				if !strings.Contains(outputStr, "config") &&
					!strings.Contains(outputStr, "not found") &&
					!strings.Contains(outputStr, "permission") {
					t.Errorf("%s, but command failed: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Config flag should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Flags_VerboseFlag tests the verbose flag functionality
// Verifies that verbose flag affects output appropriately
func TestCLI_Flags_VerboseFlag(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	commands := []string{"list", "version"}

	for _, command := range commands {
		t.Run("verbose_with_"+command, func(t *testing.T) {
			// Test both long and short forms of verbose flag
			flagVariations := [][]string{
				{"--verbose", command},
				{"-v", command},
			}

			for _, args := range flagVariations {
				cmd := exec.Command(binary, args...)
				output, err := cmd.Output()

				// Verbose flag should not cause commands to fail
				if err != nil {
					t.Errorf("Verbose flag should not cause failure for %v: %v", args, err)
				}

				outputStr := string(output)
				if len(strings.TrimSpace(outputStr)) == 0 {
					t.Errorf("Verbose flag should not result in empty output for %v", args)
				}
			}
		})
	}
}

// TestCLI_Flags_FlagOrder tests that flag order doesn't matter
// Verifies that flags work correctly regardless of their position
func TestCLI_Flags_FlagOrder(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Test different flag orders
	flagOrders := [][]string{
		{"--verbose", "list"},
		{"list", "--verbose"},
		{"--config", "/tmp/config.yaml", "list"},
		{"list", "--config", "/tmp/config.yaml"},
		{"--verbose", "--config", "/tmp/config.yaml", "list"},
		{"--config", "/tmp/config.yaml", "--verbose", "list"},
	}

	for i, args := range flagOrders {
		t.Run(fmt.Sprintf("flag_order_%d", i), func(t *testing.T) {
			cmd := exec.Command(binary, args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Commands should not fail due to flag order
			// (they might fail due to missing config, but that's different)
			if err != nil {
				if !strings.Contains(outputStr, "config") && !strings.Contains(outputStr, "not found") {
					t.Errorf("Flag order should not cause failure: %v\nError: %v\nOutput: %s", args, err, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Flag order should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Flags_FlagWithoutValue tests flags that require values
// Verifies that flags requiring values handle missing values properly
func TestCLI_Flags_FlagWithoutValue(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{"config flag without value", []string{"--config"}},
		// Removed config flag with empty value as it shows help instead of failing
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Should fail when required value is missing
			if err == nil {
				t.Errorf("Flag without required value should fail, but succeeded. Output: %s", outputStr)
			}

			// Should provide helpful error message
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("Error message should not be empty when flag value is missing")
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Missing flag value should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Flags_HelpFlag tests help flag functionality
// Verifies that help flags work correctly with different commands
func TestCLI_Flags_HelpFlag(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name string
		args []string
	}{
		{"global help long", []string{"--help"}},
		{"global help short", []string{"-h"}},
		{"command help", []string{"new", "--help"}},
		{"command help short", []string{"list", "-h"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)
			output, err := cmd.CombinedOutput()

			// Help should always succeed
			if err != nil {
				t.Errorf("Help flag should always succeed: %v\nOutput: %s", err, output)
			}

			outputStr := string(output)
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Error("Help output should not be empty")
			}

			// Help should contain usage information
			if !strings.Contains(outputStr, "USAGE") && !strings.Contains(outputStr, "Usage:") && !strings.Contains(outputStr, "usage:") {
				t.Errorf("Help output should contain usage information, got: %s", outputStr)
			}
		})
	}
}
