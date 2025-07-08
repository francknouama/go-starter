package cmd

import (
	"bytes"
	"os"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShowVersion(t *testing.T) {
	// Save original values
	originalVersion := Version
	originalCommit := Commit
	originalDate := Date
	
	// Set test values
	Version = "1.0.0-test"
	Commit = "abc123"
	Date = "2024-01-01"
	
	// Restore original values after test
	defer func() {
		Version = originalVersion
		Commit = originalCommit
		Date = originalDate
	}()
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	showVersion()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Verify output contains expected information
	assert.Contains(t, output, "1.0.0-test", "Should contain version")
	assert.Contains(t, output, "abc123", "Should contain commit")
	assert.Contains(t, output, "2024-01-01", "Should contain date")
	assert.Contains(t, output, runtime.Version(), "Should contain Go version")
	assert.Contains(t, output, runtime.GOOS, "Should contain OS")
	assert.Contains(t, output, runtime.GOARCH, "Should contain architecture")
	
	// Check for expected sections
	assert.Contains(t, output, "Version Information", "Should contain header")
	assert.Contains(t, output, "Ready to generate", "Should contain footer message")
	
	// Should not be empty
	assert.NotEmpty(t, output)
}

func TestVersionCmd_Configuration(t *testing.T) {
	// Test that the version command is properly configured
	assert.Equal(t, "version", versionCmd.Use)
	assert.Equal(t, "Show version information", versionCmd.Short)
	assert.NotEmpty(t, versionCmd.Long)
	assert.NotNil(t, versionCmd.Run)
}

func TestVersionCmd_Execution(t *testing.T) {
	// Test that the version command can be executed without panicking
	assert.NotPanics(t, func() {
		// Create a buffer to capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the command function
		versionCmd.Run(versionCmd, []string{})

		// Restore stdout
		_ = w.Close()
		os.Stdout = oldStdout

		// Read and verify we got some output
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		output := buf.String()
		
		// Should produce output
		assert.NotEmpty(t, output)
	})
}

func TestVersionVariables(t *testing.T) {
	// Test that version variables are set to reasonable defaults
	assert.NotEmpty(t, Version, "Version should not be empty")
	assert.NotEmpty(t, Commit, "Commit should not be empty")
	assert.NotEmpty(t, Date, "Date should not be empty")
	
	// Test that they are strings
	assert.IsType(t, "", Version)
	assert.IsType(t, "", Commit)
	assert.IsType(t, "", Date)
}

func TestShowVersion_WithDifferentValues(t *testing.T) {
	tests := []struct {
		name    string
		version string
		commit  string
		date    string
	}{
		{
			name:    "normal values",
			version: "1.2.3",
			commit:  "def456",
			date:    "2024-12-01",
		},
		{
			name:    "dev version",
			version: "dev",
			commit:  "local",
			date:    "now",
		},
		{
			name:    "empty values",
			version: "",
			commit:  "",
			date:    "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original values
			originalVersion := Version
			originalCommit := Commit
			originalDate := Date
			
			// Set test values
			Version = tt.version
			Commit = tt.commit
			Date = tt.date
			
			// Restore original values after test
			defer func() {
				Version = originalVersion
				Commit = originalCommit
				Date = originalDate
			}()
			
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			assert.NotPanics(t, func() {
				showVersion()
			})

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout

			// Read the output
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			// Should always produce some output
			assert.NotEmpty(t, output)
			
			// Should contain runtime information regardless of input
			assert.Contains(t, output, runtime.Version())
			assert.Contains(t, output, runtime.GOOS)
			assert.Contains(t, output, runtime.GOARCH)
		})
	}
}

func TestShowVersion_OutputFormat(t *testing.T) {
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	showVersion()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// Check that output is formatted properly (contains key sections)
	lines := strings.Split(output, "\n")
	assert.Greater(t, len(lines), 5, "Should have multiple lines of output")
	
	// Should contain structured information
	hasVersionLine := false
	hasCommitLine := false
	hasBuiltLine := false
	hasGoLine := false
	
	for _, line := range lines {
		if strings.Contains(line, "Version:") {
			hasVersionLine = true
		}
		if strings.Contains(line, "Commit:") {
			hasCommitLine = true
		}
		if strings.Contains(line, "Built:") {
			hasBuiltLine = true
		}
		if strings.Contains(line, "Go:") {
			hasGoLine = true
		}
	}
	
	assert.True(t, hasVersionLine, "Should have version line")
	assert.True(t, hasCommitLine, "Should have commit line")
	assert.True(t, hasBuiltLine, "Should have built line")
	assert.True(t, hasGoLine, "Should have Go version line")
}