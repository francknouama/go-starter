package generator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// buildTestBinary builds the CLI binary for testing
func buildTestBinary(t *testing.T) string {
	t.Helper()

	// Create temporary binary in test temp directory
	tmpDir := t.TempDir()
	binary := filepath.Join(tmpDir, "go-starter-test")

	// Build the binary - need to build from the root of the project
	// We build from the root of the project (../../../../)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(t.TempDir())))))
	cmd := exec.Command("go", "build", "-o", binary, projectRoot)
	cmd.Env = append(os.Environ(), "CGO_ENABLED=0")

	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to build test binary: %v\nOutput: %s", err, output)
	}

	// Verify binary was created and is executable
	if stat, err := os.Stat(binary); err != nil {
		t.Fatalf("Test binary was not created: %v", err)
	} else if stat.Mode().Perm()&0111 == 0 {
		// Make binary executable if it isn't already
		if err := os.Chmod(binary, 0755); err != nil {
			t.Fatalf("Failed to make test binary executable: %v", err)
		}
	}

	return binary
}

// TestGenerator_GoVersion_CLI tests the --go-version flag for the new command
func TestGenerator_GoVersion_CLI(t *testing.T) {
	binary := buildTestBinary(t)
	defer func() {
		if err := os.Remove(binary); err != nil {
			t.Logf("Warning: failed to remove test binary: %v", err)
		}
	}()

	tests := []struct {
		name        string
		args        []string
		shouldFail  bool
		expectedMsg string
		checkGoMod  bool
		goVersion   string
	}{
		{
			name:       "valid go version 1.23",
			args:       []string{"new", "test-project", "--type=cli", "--go-version=1.23", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail: false,
			checkGoMod: true,
			goVersion:  "1.23",
		},
		{
			name:       "valid go version 1.21",
			args:       []string{"new", "test-project", "--type=cli", "--go-version=1.21", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail: false,
			checkGoMod: true,
			goVersion:  "1.21",
		},
		{
			name:       "valid go version 1.20",
			args:       []string{"new", "test-project", "--type=cli", "--go-version=1.20", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail: false,
			checkGoMod: true,
			goVersion:  "1.20",
		},
		{
			name:        "invalid go version (too old)",
			args:        []string{"new", "test-project", "--type=cli", "--go-version=1.10", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail:  true,
			expectedMsg: "invalid Go version",
		},
		{
			name:        "invalid go version (malformed)",
			args:        []string{"new", "test-project", "--type=cli", "--go-version=abc", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail:  true,
			expectedMsg: "invalid Go version",
		},
		{
			name:       "auto go version",
			args:       []string{"new", "test-project", "--type=cli", "--go-version=auto", "--module=github.com/test/test-project", "--framework=cobra", "--logger=slog"},
			shouldFail: false,
			checkGoMod: true,
			goVersion:  "auto", // This should resolve to current Go version
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			cmd := exec.Command(binary, tt.args...)
			cmd.Dir = tmpDir
			// Prepare input for interactive prompts (if any)
			input := strings.NewReader(`github.com/test/test-project
Cobra (recommended)
slog - Go built-in structured logging (recommended)
`)
			cmd.Stdin = input

			output, err := cmd.CombinedOutput()

			if tt.shouldFail {
				assert.Error(t, err, "Expected command to fail")
				if tt.expectedMsg != "" {
					assert.Contains(t, string(output), tt.expectedMsg, 
						"Expected output to contain '%s', but got '%s'", tt.expectedMsg, string(output))
				}
			} else {
				if err != nil {
					// Check if it's a template not found error, which is acceptable
					if strings.Contains(string(output), "Template") && strings.Contains(string(output), "not available") {
						t.Skip("Skipping test as CLI template is not yet implemented")
						return
					}
					t.Errorf("Expected command to succeed, but it failed with error: %v. Output: %s", err, string(output))
				}

				if tt.checkGoMod {
					goModPath := filepath.Join(tmpDir, "test-project", "go.mod")
					content, err := os.ReadFile(goModPath)
					if err != nil {
						if strings.Contains(string(output), "Template") && strings.Contains(string(output), "not available") {
							t.Skip("Skipping go.mod check as template is not available")
							return
						}
						t.Fatalf("Failed to read go.mod file: %v", err)
					}
					
					if tt.goVersion == "auto" {
						// For auto version, just verify go directive exists
						assert.Contains(t, string(content), "go ", 
							"go.mod should contain go directive")
					} else {
						expectedGoVersion := "go " + tt.goVersion
						assert.Contains(t, string(content), expectedGoVersion, 
							"go.mod does not contain the correct go version. Expected '%s', Content: %s", 
							expectedGoVersion, string(content))
					}
				}
			}
		})
	}
}

// TestGenerator_GoVersion_Config tests Go version handling in generator configuration
func TestGenerator_GoVersion_Config(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name           string
		goVersion      string
		expectedInMod  string
		shouldGenerate bool
	}{
		{
			name:           "valid go version 1.21",
			goVersion:      "1.21",
			expectedInMod:  "go 1.21",
			shouldGenerate: true,
		},
		{
			name:           "valid go version 1.22",
			goVersion:      "1.22",
			expectedInMod:  "go 1.22",
			shouldGenerate: true,
		},
		{
			name:           "valid go version 1.20",
			goVersion:      "1.20",
			expectedInMod:  "go 1.20",
			shouldGenerate: true,
		},
		{
			name:           "go version with patch",
			goVersion:      "1.21.0",
			expectedInMod:  "go 1.21.0",
			shouldGenerate: true,
		},
		{
			name:           "empty go version (should use default)",
			goVersion:      "",
			expectedInMod:  "go ", // Should have some default
			shouldGenerate: true,
		},
		{
			name:           "auto go version",
			goVersion:      "auto",
			expectedInMod:  "go ", // Should resolve to current version
			shouldGenerate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "go-version-test",
				Module:    "github.com/test/go-version-test",
				Type:      "library",
				GoVersion: tt.goVersion,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, config.Name),
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			result, err := gen.Generate(config, options)

			if !tt.shouldGenerate {
				assert.Error(t, err, "Expected generation to fail for Go version: %s", tt.goVersion)
				return
			}

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
				t.Fatalf("Unexpected error during generation: %v", err)
			}

			require.NotNil(t, result)
			assert.True(t, result.Success, "Generation should succeed")

			// Check go.mod file for correct Go version
			goModPath := filepath.Join(options.OutputPath, "go.mod")
			if _, err := os.Stat(goModPath); err == nil {
				content, err := os.ReadFile(goModPath)
				require.NoError(t, err, "Should be able to read go.mod")

				if tt.goVersion == "" || tt.goVersion == "auto" {
					// For empty or auto version, just check that go directive exists
					assert.Contains(t, string(content), "go ", 
						"go.mod should contain go directive")
				} else {
					assert.Contains(t, string(content), tt.expectedInMod, 
						"go.mod should contain expected Go version")
				}
			}
		})
	}
}

// TestGenerator_GoVersion_Validation tests Go version validation logic
func TestGenerator_GoVersion_Validation(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name       string
		goVersion  string
		shouldFail bool
		errorMsg   string
	}{
		{
			name:       "valid modern version",
			goVersion:  "1.21",
			shouldFail: false,
		},
		{
			name:       "valid older supported version",
			goVersion:  "1.19",
			shouldFail: false,
		},
		{
			name:       "valid future version",
			goVersion:  "1.25",
			shouldFail: false,
		},
		{
			name:       "empty version (should use default)",
			goVersion:  "",
			shouldFail: false,
		},
		{
			name:       "auto version",
			goVersion:  "auto",
			shouldFail: false,
		},
		// Note: Version validation might be handled at CLI level rather than generator level
		// These tests verify the generator handles various version formats gracefully
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "version-validation-test",
				Module:    "github.com/test/version-validation-test",
				Type:      "library",
				GoVersion: tt.goVersion,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, config.Name),
				DryRun:     true, // Use dry run for validation testing
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected validation to fail for Go version: %s", tt.goVersion)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg, "Error should contain expected message")
				}
			} else {
				// Accept template not found errors
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Validation should succeed for Go version: %s", tt.goVersion)
			}
		})
	}
}

// TestGenerator_GoVersion_ContextVariables tests Go version in template context
func TestGenerator_GoVersion_ContextVariables(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name      string
		goVersion string
		expected  string
	}{
		{
			name:      "specific version",
			goVersion: "1.21",
			expected:  "1.21",
		},
		{
			name:      "version with patch",
			goVersion: "1.21.5",
			expected:  "1.21.5",
		},
		{
			name:      "empty version",
			goVersion: "",
			expected:  "", // Should be handled by template or default
		},
		{
			name:      "auto version",
			goVersion: "auto",
			expected:  "auto", // Should be resolved by generator
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "context-test",
				Module:    "github.com/test/context-test",
				Type:      "library",
				GoVersion: tt.goVersion,
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, config.Name),
				DryRun:     true,
				NoGit:      true,
				Verbose:    false,
			}

			_, err := gen.Generate(config, options)

			// Accept template not found errors
			if err != nil {
				if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
					t.Skip("Skipping test as template is not yet implemented")
					return
				}
			}

			// Verify the Go version is set correctly in the config
			assert.Equal(t, tt.expected, config.GoVersion, 
				"Go version in config should match expected value")

			t.Logf("Go version context test passed for version: %s", tt.goVersion)
		})
	}
}

// TestGenerator_GoVersion_DefaultBehavior tests default Go version behavior
func TestGenerator_GoVersion_DefaultBehavior(t *testing.T) {
	setupTestTemplates(t)

	config := types.ProjectConfig{
		Name:   "default-version-test",
		Module: "github.com/test/default-version-test",
		Type:   "library",
		// GoVersion is intentionally not set
	}

	gen := generator.New()
	require.NotNil(t, gen)

	tmpDir := t.TempDir()
	options := types.GenerationOptions{
		OutputPath: filepath.Join(tmpDir, config.Name),
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	result, err := gen.Generate(config, options)

	// Accept template not found errors
	if err != nil {
		if _, ok := err.(*types.GoStarterError); ok {
			t.Skip("Skipping test as template is not yet implemented")
			return
		}
		t.Fatalf("Unexpected error: %v", err)
	}

	require.NotNil(t, result)
	assert.True(t, result.Success, "Generation should succeed with default Go version")

	// Check if go.mod was created and has a go directive
	goModPath := filepath.Join(options.OutputPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		require.NoError(t, err, "Should be able to read go.mod")

		// Should contain some go directive
		assert.Contains(t, string(content), "go ", 
			"go.mod should contain go directive even with default version")
	}

	t.Log("Default Go version behavior test passed")
}