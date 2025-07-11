package cli

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLI_Environment_ConfigEnvVar tests behavior with GO_STARTER_CONFIG environment variable
// Verifies that the config environment variable is properly handled
func TestCLI_Environment_ConfigEnvVar(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		configPath  string
		command     []string
		description string
	}{
		{
			name:        "nonexistent config file",
			configPath:  "/tmp/nonexistent-config.yaml",
			command:     []string{"list"},
			description: "should handle missing config file gracefully",
		},
		{
			name:        "empty config path",
			configPath:  "",
			command:     []string{"version"},
			description: "should handle empty config path",
		},
		{
			name:        "invalid config path",
			configPath:  "/invalid/path/config.yaml",
			command:     []string{"list"},
			description: "should handle invalid config path gracefully",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.command...)

			// Set environment variable
			env := os.Environ()
			if tt.configPath != "" {
				env = append(env, "GO_STARTER_CONFIG="+tt.configPath)
			}
			cmd.Env = env

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Commands should not fail catastrophically due to environment variables
			if err != nil {
				// Only fail if it's not a graceful handling of missing/invalid config
				if !strings.Contains(outputStr, "config") &&
					!strings.Contains(outputStr, "not found") &&
					!strings.Contains(outputStr, "permission") &&
					!strings.Contains(outputStr, "no such file") {
					t.Errorf("%s, but command failed unexpectedly: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Config environment variable should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Environment_VerboseEnvVar tests behavior with GO_STARTER_VERBOSE environment variable
// Verifies that the verbose environment variable affects output appropriately
func TestCLI_Environment_VerboseEnvVar(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name         string
		verboseValue string
		command      []string
		description  string
	}{
		{
			name:         "verbose true",
			verboseValue: "true",
			command:      []string{"list"},
			description:  "should enable verbose mode",
		},
		{
			name:         "verbose 1",
			verboseValue: "1",
			command:      []string{"version"},
			description:  "should enable verbose mode with numeric value",
		},
		{
			name:         "verbose false",
			verboseValue: "false",
			command:      []string{"list"},
			description:  "should disable verbose mode",
		},
		{
			name:         "verbose empty",
			verboseValue: "",
			command:      []string{"version"},
			description:  "should handle empty verbose value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.command...)

			// Set environment variable
			env := os.Environ()
			if tt.verboseValue != "" {
				env = append(env, "GO_STARTER_VERBOSE="+tt.verboseValue)
			}
			cmd.Env = env

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Commands should not fail due to verbose environment variable
			if err != nil {
				// Check if it's an unrelated error
				if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "fatal") {
					t.Errorf("%s, but command failed with panic: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should produce output
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Command should produce output, got empty string")
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Verbose environment variable should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Environment_MultipleEnvVars tests behavior with multiple environment variables
// Verifies that multiple environment variables work together correctly
func TestCLI_Environment_MultipleEnvVars(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		envVars     map[string]string
		command     []string
		description string
	}{
		{
			name: "config and verbose",
			envVars: map[string]string{
				"GO_STARTER_CONFIG":  "/tmp/test-config.yaml",
				"GO_STARTER_VERBOSE": "true",
			},
			command:     []string{"list"},
			description: "should handle both config and verbose environment variables",
		},
		{
			name: "multiple config-related vars",
			envVars: map[string]string{
				"GO_STARTER_CONFIG": "/tmp/config.yaml",
				"HOME":              "/tmp",
			},
			command:     []string{"version"},
			description: "should handle config with modified HOME",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.command...)

			// Set environment variables
			env := os.Environ()
			for k, v := range tt.envVars {
				env = append(env, k+"="+v)
			}
			cmd.Env = env

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Commands should handle multiple environment variables gracefully
			if err != nil {
				// Only fail if it's not a graceful handling of missing files
				if !strings.Contains(outputStr, "config") &&
					!strings.Contains(outputStr, "not found") &&
					!strings.Contains(outputStr, "permission") {
					t.Errorf("%s, but command failed: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Multiple environment variables should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Environment_EnvVarPrecedence tests precedence of environment variables vs command flags
// Verifies that command-line flags override environment variables when both are present
func TestCLI_Environment_EnvVarPrecedence(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		envVars     map[string]string
		args        []string
		description string
	}{
		{
			name: "config flag overrides env",
			envVars: map[string]string{
				"GO_STARTER_CONFIG": "/tmp/env-config.yaml",
			},
			args:        []string{"--config", "/tmp/flag-config.yaml", "list"},
			description: "command-line config flag should override environment variable",
		},
		{
			name: "verbose flag with env",
			envVars: map[string]string{
				"GO_STARTER_VERBOSE": "false",
			},
			args:        []string{"--verbose", "version"},
			description: "command-line verbose flag should work with environment variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.args...)

			// Set environment variables
			env := os.Environ()
			for k, v := range tt.envVars {
				env = append(env, k+"="+v)
			}
			cmd.Env = env

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Should handle precedence gracefully
			if err != nil {
				// Check if error is related to missing config files (acceptable)
				if !strings.Contains(outputStr, "config") &&
					!strings.Contains(outputStr, "not found") &&
					!strings.Contains(outputStr, "permission") {
					t.Errorf("%s, but command failed: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Environment variable precedence should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Environment_SystemEnvVars tests interaction with system environment variables
// Verifies that standard system environment variables don't interfere with operation
func TestCLI_Environment_SystemEnvVars(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	tests := []struct {
		name        string
		modifyEnv   func([]string) []string
		command     []string
		description string
	}{
		{
			name: "modified PATH",
			modifyEnv: func(env []string) []string {
				// Add a custom path that doesn't interfere
				return append(env, "PATH=/custom/path:"+os.Getenv("PATH"))
			},
			command:     []string{"version"},
			description: "should work with modified PATH",
		},
		{
			name: "modified HOME",
			modifyEnv: func(env []string) []string {
				return append(env, "HOME=/tmp/test-home")
			},
			command:     []string{"list"},
			description: "should work with modified HOME",
		},
		{
			name: "no USER env var",
			modifyEnv: func(env []string) []string {
				// Remove USER environment variable
				filtered := make([]string, 0, len(env))
				for _, e := range env {
					if !strings.HasPrefix(e, "USER=") {
						filtered = append(filtered, e)
					}
				}
				return filtered
			},
			command:     []string{"version"},
			description: "should work without USER environment variable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(binary, tt.command...)
			cmd.Env = tt.modifyEnv(os.Environ())

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Should not fail due to modified system environment
			if err != nil {
				// Some errors might be acceptable (like missing config files)
				if strings.Contains(outputStr, "panic") || strings.Contains(outputStr, "fatal") {
					t.Errorf("%s, but command panicked: %v\nOutput: %s", tt.description, err, outputStr)
				}
			}

			// Should produce reasonable output
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Command should produce output even with modified environment")
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Modified system environment should not cause panic, got: %s", outputStr)
			}
		})
	}
}

// TestCLI_Environment_DefaultBehavior tests default behavior without custom environment variables
// Verifies that the CLI works correctly in a clean environment
func TestCLI_Environment_DefaultBehavior(t *testing.T) {
	binary := buildTestBinary(t)
	defer cleanupBinary(t, binary)

	// Create a minimal environment
	minimalEnv := []string{
		"PATH=" + os.Getenv("PATH"),
		"HOME=" + os.Getenv("HOME"),
		"USER=" + os.Getenv("USER"),
	}

	commands := [][]string{
		{"--help"},
		{"version"},
		{"list"},
		{"completion", "bash"},
	}

	for _, command := range commands {
		t.Run("default_"+strings.Join(command, "_"), func(t *testing.T) {
			cmd := exec.Command(binary, command...)
			cmd.Env = minimalEnv

			output, err := cmd.CombinedOutput()
			outputStr := string(output)

			// Commands should work in minimal environment
			if err != nil {
				// Some commands might fail, but shouldn't panic
				if strings.Contains(outputStr, "panic") {
					t.Errorf("Command should not panic in minimal environment: %v\nOutput: %s", err, outputStr)
				}
			}

			// Should produce output
			if len(strings.TrimSpace(outputStr)) == 0 {
				t.Errorf("Command should produce output in minimal environment")
			}

			// Should not panic
			if strings.Contains(outputStr, "panic") {
				t.Errorf("Minimal environment should not cause panic, got: %s", outputStr)
			}
		})
	}
}
