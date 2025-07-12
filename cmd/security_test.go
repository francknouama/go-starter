package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/francknouama/go-starter/internal/security"
)

func setupSecurityTestBlueprints(t *testing.T) string {
	t.Helper()

	// Get the project root for tests
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(file))
	blueprintsDir := filepath.Join(projectRoot, "blueprints")

	// Verify blueprints directory exists
	if _, err := os.Stat(blueprintsDir); os.IsNotExist(err) {
		t.Fatalf("Blueprints directory not found at: %s", blueprintsDir)
	}

	return blueprintsDir
}


func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "config.yaml")
	
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	require.NoError(t, err)
	
	return tmpFile
}

func TestSecurityCmd_Configuration(t *testing.T) {
	// Test that the security command is properly configured
	assert.Equal(t, "security", securityCmd.Use)
	assert.Equal(t, "Security scanning and validation tools", securityCmd.Short)
	assert.NotEmpty(t, securityCmd.Long)
	assert.Contains(t, securityCmd.Long, "scan-blueprints")
	assert.Contains(t, securityCmd.Long, "scan-config")
}

func TestScanBlueprintsCmd_Configuration(t *testing.T) {
	// Test that the scan-blueprints command is properly configured
	assert.Equal(t, "scan-blueprints [path]", scanBlueprintsCmd.Use)
	assert.Equal(t, "Scan blueprint files for security vulnerabilities", scanBlueprintsCmd.Short)
	assert.NotEmpty(t, scanBlueprintsCmd.Long)
	assert.Contains(t, scanBlueprintsCmd.Long, "blueprint files")
	assert.NotNil(t, scanBlueprintsCmd.RunE)
}

func TestScanConfigCmd_Configuration(t *testing.T) {
	// Test that the scan-config command is properly configured
	assert.Equal(t, "scan-config [config-file]", scanConfigCmd.Use)
	assert.Equal(t, "Validate project configuration for security issues", scanConfigCmd.Short)
	assert.NotEmpty(t, scanConfigCmd.Long)
	assert.NotNil(t, scanConfigCmd.RunE)
}

func TestScanBlueprints_WithValidBlueprints(t *testing.T) {
	blueprintsDir := setupSecurityTestBlueprints(t)
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run scan with actual blueprints directory
	err := scanBlueprints(blueprintsDir, false, "console")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error with valid blueprints
	assert.NoError(t, err)
	
	// Should produce some output
	assert.NotEmpty(t, output)
	
	// Should either show security violations or no violations message
	assert.True(t, 
		strings.Contains(output, "No security violations found") ||
		strings.Contains(output, "security violations") ||
		strings.Contains(output, "Scanning:"),
		"Output should contain security scan results or scanning messages")
}

func TestScanBlueprints_WithNonExistentPath(t *testing.T) {
	nonExistentPath := "/path/that/does/not/exist"
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run scan with non-existent path
	err := scanBlueprints(nonExistentPath, false, "console")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)

	// Should return an error for non-existent path
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error scanning blueprints")
}

func TestScanBlueprints_VerboseMode(t *testing.T) {
	// Create a temporary directory with a test blueprint file
	tmpDir := t.TempDir()
	blueprintFile := filepath.Join(tmpDir, "test.tmpl")
	err := os.WriteFile(blueprintFile, []byte("{{.ProjectName}}"), 0644)
	require.NoError(t, err)
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run scan in verbose mode
	err = scanBlueprints(tmpDir, true, "console")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Verbose mode should show scanning messages
	assert.Contains(t, output, "Scanning:")
}

func TestScanBlueprints_JSONOutput(t *testing.T) {
	// Create a temporary directory with a test blueprint file
	tmpDir := t.TempDir()
	blueprintFile := filepath.Join(tmpDir, "test.tmpl")
	err := os.WriteFile(blueprintFile, []byte("{{.ProjectName}}"), 0644)
	require.NoError(t, err)
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run scan with JSON output
	err = scanBlueprints(tmpDir, false, "json")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Should produce some output
	assert.NotEmpty(t, output)
}

func TestScanConfig_WithValidConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
name: test-project
module: github.com/test/project
type: web-api
`
	configFile := createTempConfigFile(t, configContent)
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run config scan
	err := scanConfig(configFile, false, "console")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error with valid config
	assert.NoError(t, err)
	
	// Should produce validation output
	assert.Contains(t, output, "Configuration file validation completed")
}

func TestScanConfig_WithNonExistentFile(t *testing.T) {
	nonExistentFile := "/path/that/does/not/exist.yaml"
	
	// Run config scan with non-existent file
	err := scanConfig(nonExistentFile, false, "console")

	// Should return an error for non-existent file
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "configuration file not found")
}

func TestScanConfig_VerboseMode(t *testing.T) {
	// Create a temporary config file
	configFile := createTempConfigFile(t, "name: test")
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run config scan in verbose mode
	err := scanConfig(configFile, true, "console")

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Verbose mode should show scanning message
	assert.Contains(t, output, "Scanning configuration:")
}

func TestOutputSecurityResults_NoViolations(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with no violations
	err := outputSecurityResults([]security.SecurityViolation{}, "console", false)

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Should show no violations message
	assert.Contains(t, output, "No security violations found")
}

func TestOutputSecurityResults_WithViolations(t *testing.T) {
	// Mock security violations (we'll need to import the security package properly)
	// For now, we'll test with a simplified approach
	violations := []security.SecurityViolation{
		{
			Type:        "test-violation",
			Description: "Test security violation",
			Severity:    "HIGH",
			Line:        10,
		},
	}
	
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test with violations
	err := outputSecurityResults(violations, "console", false)

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Should show violations
	assert.Contains(t, output, "security violations")
	assert.Contains(t, output, "test-violation")
}

func TestOutputSecurityResults_JSONFormat(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Test JSON output with no violations
	err := outputSecurityResults([]security.SecurityViolation{}, "json", false)

	// Restore stdout and capture output
	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error
	assert.NoError(t, err)
	
	// Should produce JSON-like output
	assert.Contains(t, output, "violations")
}

func TestScanBlueprintsCmd_Execution(t *testing.T) {
	// Test that the scan-blueprints command can be executed without panicking
	assert.NotPanics(t, func() {
		// Create a buffer to capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the command function with default blueprints path
		err := scanBlueprintsCmd.RunE(scanBlueprintsCmd, []string{})

		// Restore stdout
		closeErr := w.Close()
		require.NoError(t, closeErr)
		os.Stdout = oldStdout

		// Read output
		var buf bytes.Buffer
		_, readErr := buf.ReadFrom(r)
		require.NoError(t, readErr)

		// Should handle the execution (may error due to missing blueprints in test env)
		// but shouldn't panic
		_ = err // We expect this might error in test environment
	})
}

func TestScanConfigCmd_Execution(t *testing.T) {
	// Test that the scan-config command can be executed without panicking
	assert.NotPanics(t, func() {
		// Create a temporary config file
		configFile := createTempConfigFile(t, "name: test")
		
		// Create a buffer to capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the command function
		err := scanConfigCmd.RunE(scanConfigCmd, []string{configFile})

		// Restore stdout
		closeErr := w.Close()
		require.NoError(t, closeErr)
		os.Stdout = oldStdout

		// Read output
		var buf bytes.Buffer
		_, readErr := buf.ReadFrom(r)
		require.NoError(t, readErr)

		// Should not error with valid config file
		assert.NoError(t, err)
	})
}

func TestSecurityCmd_Flags(t *testing.T) {
	// Test that flags are properly configured
	verboseFlag := scanBlueprintsCmd.Flags().Lookup("verbose")
	assert.NotNil(t, verboseFlag)
	assert.Equal(t, "bool", verboseFlag.Value.Type())
	
	outputFlag := scanBlueprintsCmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlag)
	assert.Equal(t, "string", outputFlag.Value.Type())
	
	// Test same flags on config command
	verboseFlagConfig := scanConfigCmd.Flags().Lookup("verbose")
	assert.NotNil(t, verboseFlagConfig)
	
	outputFlagConfig := scanConfigCmd.Flags().Lookup("output")
	assert.NotNil(t, outputFlagConfig)
}

func TestSecurityCmd_SubcommandRelations(t *testing.T) {
	// Test that security command has the expected subcommands
	subCommands := securityCmd.Commands()
	assert.Len(t, subCommands, 2, "Security command should have exactly 2 subcommands")
	
	// Find scan-blueprints and scan-config commands
	var foundScanBlueprints, foundScanConfig bool
	for _, cmd := range subCommands {
		switch cmd.Use {
		case "scan-blueprints [path]":
			foundScanBlueprints = true
		case "scan-config [config-file]":
			foundScanConfig = true
		}
	}
	
	assert.True(t, foundScanBlueprints, "Should have scan-blueprints subcommand")
	assert.True(t, foundScanConfig, "Should have scan-config subcommand")
}

func TestScanBlueprints_EdgeCases(t *testing.T) {
	// Test with empty directory
	emptyDir := t.TempDir()
	
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	err := scanBlueprints(emptyDir, false, "console")

	closeErr := w.Close()
	require.NoError(t, closeErr)
	os.Stdout = oldStdout
	_, readErr := buf.ReadFrom(r)
	require.NoError(t, readErr)
	output := buf.String()

	// Should not error with empty directory
	assert.NoError(t, err)
	assert.Contains(t, output, "No security violations found")
}