package cmd

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/verify/cli/internal/logger"
)

func TestRootCommand(t *testing.T) {
	// Create a test logger
	factory := logger.NewFactory()
	testLogger, err := factory.CreateFromProjectConfig("slog", "info", "text", false)
	assert.NoError(t, err)

	// Capture output
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	// Test root command execution
	rootCmd.SetArgs([]string{})
	err = Execute(testLogger)
	assert.NoError(t, err)

	// Verify output contains expected content
	outputStr := output.String()
	assert.Contains(t, outputStr, "Welcome to verify-cli!")
	assert.Contains(t, outputStr, "Use --help to see available commands")
}

func TestVersionCommand(t *testing.T) {
	// Create a test logger
	factory := logger.NewFactory()
	testLogger, err := factory.CreateFromProjectConfig("slog", "info", "text", false)
	assert.NoError(t, err)

	// Capture output
	var output bytes.Buffer
	rootCmd.SetOut(&output)
	rootCmd.SetErr(&output)

	// Test version command
	rootCmd.SetArgs([]string{"version"})
	err = Execute(testLogger)
	assert.NoError(t, err)

	// Verify output contains version information
	outputStr := output.String()
	assert.Contains(t, outputStr, "verify-cli version information")
	assert.Contains(t, outputStr, "Version:")
	assert.Contains(t, outputStr, "Logger:     slog")
}