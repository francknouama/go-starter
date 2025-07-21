package cli

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"
)

// Test configuration constants
const (
	// Timeout durations for different command types
	FastCommandTimeout    = 5 * time.Second  // help, version, etc.
	SlowCommandTimeout    = 15 * time.Second // list, completion, etc.
	InteractiveTimeout    = 10 * time.Second // new command with prompts
	ResourceIntensiveTimeout = 20 * time.Second // verbose operations

	// Common test paths and values
	NonexistentConfigPath = "/tmp/nonexistent.yaml"
	TestProjectName      = "test-project"
	TestModulePath       = "github.com/test/example"
)

// CommandResult holds the result of a command execution
type CommandResult struct {
	Output   string
	Error    error
	ExitCode int
	TimedOut bool
}

// TestBinary manages a test binary instance with caching
type TestBinary struct {
	path string
	mu   sync.RWMutex
}

var (
	globalTestBinary *TestBinary
	binaryOnce       sync.Once
)

// GetTestBinary returns a cached test binary instance
func GetTestBinary(t *testing.T) *TestBinary {
	binaryOnce.Do(func() {
		globalTestBinary = &TestBinary{}
	})
	
	// Build binary if not already built
	globalTestBinary.mu.RLock()
	needsBuild := globalTestBinary.path == ""
	globalTestBinary.mu.RUnlock()
	
	if needsBuild {
		globalTestBinary.build(t)
	}
	
	return globalTestBinary
}

// build builds the CLI binary for testing
func (tb *TestBinary) build(t *testing.T) {
	t.Helper()
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.path != "" {
		return // Already built
	}

	// Create temporary binary in system temp directory
	tmpDir, err := os.MkdirTemp("", "go-starter-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	
	binary := filepath.Join(tmpDir, "go-starter-test")

	// Build the binary - need to build the whole package to include embed.go
	// We build from the root of the project (../../..)
	cmd := exec.Command("go", "build", "-o", binary, "../../..")
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

	tb.path = binary
	
	// Clean up the temp directory when tests finish
	// Note: This is a simple approach - in production we might want more sophisticated cleanup
}

// Path returns the path to the test binary
func (tb *TestBinary) Path() string {
	tb.mu.RLock()
	defer tb.mu.RUnlock()
	return tb.path
}

// ExecuteCommand executes a command with the given arguments and options
func (tb *TestBinary) ExecuteCommand(t *testing.T, timeout time.Duration, stdin string, args ...string) *CommandResult {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, tb.Path(), args...)
	if stdin != "" {
		cmd.Stdin = strings.NewReader(stdin)
	}

	output, err := cmd.CombinedOutput()
	result := &CommandResult{
		Output: string(output),
		Error:  err,
	}

	// Check if the command timed out
	if ctx.Err() == context.DeadlineExceeded {
		result.TimedOut = true
	}

	// Extract exit code if available
	if exitErr, ok := err.(*exec.ExitError); ok {
		result.ExitCode = exitErr.ExitCode()
	}

	return result
}

// ExecuteFastCommand executes a command expected to complete quickly
func (tb *TestBinary) ExecuteFastCommand(t *testing.T, args ...string) *CommandResult {
	return tb.ExecuteCommand(t, FastCommandTimeout, "", args...)
}

// ExecuteSlowCommand executes a command that might take longer
func (tb *TestBinary) ExecuteSlowCommand(t *testing.T, args ...string) *CommandResult {
	return tb.ExecuteCommand(t, SlowCommandTimeout, "", args...)
}

// ExecuteInteractiveCommand executes an interactive command with stdin input
func (tb *TestBinary) ExecuteInteractiveCommand(t *testing.T, stdin string, args ...string) *CommandResult {
	return tb.ExecuteCommand(t, InteractiveTimeout, stdin, args...)
}

// Test assertion utilities

// AssertSuccess verifies that a command executed successfully
func AssertSuccess(t *testing.T, result *CommandResult, context string) {
	t.Helper()
	if result.TimedOut {
		t.Fatalf("%s: command timed out", context)
	}
	if result.Error != nil {
		t.Fatalf("%s: command failed: %v\nOutput: %s", context, result.Error, result.Output)
	}
}

// AssertFailure verifies that a command failed as expected
func AssertFailure(t *testing.T, result *CommandResult, context string) {
	t.Helper()
	if result.TimedOut {
		t.Fatalf("%s: command timed out when expecting failure", context)
	}
	if result.Error == nil {
		t.Fatalf("%s: command should have failed but succeeded\nOutput: %s", context, result.Output)
	}
}

// AssertContains verifies that output contains expected content
func AssertContains(t *testing.T, result *CommandResult, expected string, context string) {
	t.Helper()
	if !strings.Contains(result.Output, expected) {
		t.Errorf("%s: output should contain '%s'\nActual output: %s", context, expected, result.Output)
	}
}

// AssertNotContains verifies that output does not contain specific content
func AssertNotContains(t *testing.T, result *CommandResult, unexpected string, context string) {
	t.Helper()
	if strings.Contains(result.Output, unexpected) {
		t.Errorf("%s: output should not contain '%s'\nActual output: %s", context, unexpected, result.Output)
	}
}

// AssertNoPanic verifies that the output doesn't contain panic messages
func AssertNoPanic(t *testing.T, result *CommandResult, context string) {
	t.Helper()
	panicIndicators := []string{"panic", "runtime error", "stack trace"}
	for _, indicator := range panicIndicators {
		if strings.Contains(result.Output, indicator) {
			t.Errorf("%s: command should not panic\nOutput: %s", context, result.Output)
			break
		}
	}
}

// AssertExitCode verifies the exit code
func AssertExitCode(t *testing.T, result *CommandResult, expectedCode int, context string) {
	t.Helper()
	if result.ExitCode != expectedCode {
		t.Errorf("%s: expected exit code %d, got %d\nOutput: %s", 
			context, expectedCode, result.ExitCode, result.Output)
	}
}

// AssertNotTimedOut verifies that the command completed within timeout
func AssertNotTimedOut(t *testing.T, result *CommandResult, context string) {
	t.Helper()
	if result.TimedOut {
		t.Errorf("%s: command should not have timed out", context)
	}
}

// Test data utilities

// CommonArgs provides common argument combinations for testing
type CommonArgs struct{}

func (CommonArgs) Help() []string           { return []string{"--help"} }
func (CommonArgs) Version() []string        { return []string{"version"} }
func (CommonArgs) VersionFlag() []string    { return []string{"--version"} }
func (CommonArgs) List() []string           { return []string{"list"} }
func (CommonArgs) NewHelp() []string        { return []string{"new", "--help"} }
func (CommonArgs) Completion(shell string) []string { return []string{"completion", shell} }
func (CommonArgs) VerboseList() []string    { return []string{"--verbose", "list"} }

var Args = CommonArgs{}

// Common stdin inputs for interactive commands
type CommonInputs struct{}

func (CommonInputs) Empty() string          { return "" }
func (CommonInputs) Newline() string        { return "\n" }
func (CommonInputs) MultipleNewlines() string { return "\n\n\n\n\n" }
func (CommonInputs) ProjectName() string    { return TestProjectName + "\n" }
func (CommonInputs) SkipPrompts() string    { return "\n\n\n\n\n" }

var Inputs = CommonInputs{}

// Legacy functions for backward compatibility
func buildTestBinary(t *testing.T) string {
	t.Helper()
	return GetTestBinary(t).Path()
}

func cleanupBinary(t *testing.T, _ string) {
	t.Helper()
	// No-op since we're using cached binaries now
	// The temp directory cleanup will handle binary removal
}
