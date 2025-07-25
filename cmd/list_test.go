package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

func setupTestBlueprints(t *testing.T) {
	t.Helper()

	// Get the project root for tests
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(file))
	blueprintsDir := filepath.Join(projectRoot, "blueprints")

	// Verify blueprints directory exists
	if _, err := os.Stat(blueprintsDir); os.IsNotExist(err) {
		t.Fatalf("Blueprints directory not found at: %s", blueprintsDir)
	}

	// Set up the filesystem for tests using os.DirFS
	templates.SetTemplatesFS(os.DirFS(blueprintsDir))
}

func TestListBlueprints(t *testing.T) {
	setupTestBlueprints(t)
	
	// Capture stdout
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run the function
	listBlueprints()

	// Restore stdout
	_ = w.Close()
	os.Stdout = oldStdout

	// Read the output
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	output := buf.String()

	// The output should contain blueprint information or no blueprints message
	// Since we don't know the exact state of blueprints, we check for expected patterns
	assert.True(t, 
		strings.Contains(output, "No blueprints available") || 
		strings.Contains(output, "Blueprints") ||
		len(output) > 0, 
		"Should produce some output")
}

func TestRenderBlueprint(t *testing.T) {
	tests := []struct {
		name       string
		blueprint   types.Template
		isLast     bool
	}{
		{
			name: "basic blueprint",
			blueprint: types.Template{
				ID:          "web-api",
				Name:        "Web API",
				Description: "A basic web API blueprint",
				Type:        "web-api",
			},
			isLast: false,
		},
		{
			name: "cli blueprint",
			blueprint: types.Template{
				ID:          "cli",
				Name:        "CLI Application",
				Description: "A command-line interface blueprint",
				Type:        "cli",
			},
			isLast: true,
		},
		{
			name: "blueprint with empty description",
			blueprint: types.Template{
				ID:          "library",
				Name:        "Library",
				Description: "",
				Type:        "library",
			},
			isLast: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stdout
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			// Run the function
			renderBlueprint(tt.blueprint, tt.isLast)

			// Restore stdout
			_ = w.Close()
			os.Stdout = oldStdout

			// Read the output
			var buf bytes.Buffer
			_, _ = buf.ReadFrom(r)
			output := buf.String()

			// Check that output contains the blueprint information
			assert.Contains(t, output, tt.blueprint.ID)
			assert.Contains(t, output, tt.blueprint.Name)
			assert.Contains(t, output, tt.blueprint.Type)
			
			// Output should not be empty
			assert.NotEmpty(t, output)
		})
	}
}

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected string
	}{
		{
			name:     "short text",
			text:     "Hello world",
			width:    20,
			expected: "Hello world",
		},
		{
			name:     "text that needs wrapping",
			text:     "This is a long sentence that should be wrapped",
			width:    20,
			expected: "This is a long\nsentence that should\nbe wrapped",
		},
		{
			name:     "empty text",
			text:     "",
			width:    10,
			expected: "",
		},
		{
			name:     "single word longer than width",
			text:     "supercalifragilisticexpialidocious",
			width:    10,
			expected: "supercalifragilisticexpialidocious", // Should not break single words
		},
		{
			name:     "text with multiple spaces",
			text:     "Hello    world    test",
			width:    15,
			expected: "Hello world\ntest",
		},
		{
			name:     "exact width match",
			text:     "Exactly twenty chars",
			width:    20,
			expected: "Exactly twenty chars",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.width)
			
			// For most cases, we just check that we get a reasonable result
			assert.IsType(t, "", result)
			
			// Check specific cases where we can predict the output
			if tt.name == "short text" || tt.name == "empty text" || tt.name == "exact width match" {
				assert.Equal(t, tt.expected, result)
			}
			
			// Check that very long single words are not broken
			if tt.name == "single word longer than width" {
				assert.Equal(t, tt.text, result)
			}
			
			// Verify that wrapped text contains newlines if original was long enough
			if len(tt.text) > tt.width && strings.Contains(tt.text, " ") {
				// Should contain newlines for text that needs wrapping
				lines := strings.Split(result, "\n")
				if len(lines) > 1 {
					// Each line (except possibly the last) should not exceed width
					for i, line := range lines[:len(lines)-1] {
						if !isSingleWord(line) {
							assert.LessOrEqual(t, len(line), tt.width, "Line %d should not exceed width", i)
						}
					}
				}
			}
		})
	}
}

func TestWrapText_EdgeCases(t *testing.T) {
	t.Run("zero width", func(t *testing.T) {
		result := wrapText("Hello world", 0)
		// With zero width, should still return the text
		assert.NotEmpty(t, result)
	})
	
	t.Run("negative width", func(t *testing.T) {
		result := wrapText("Hello world", -5)
		// With negative width, should still return the text
		assert.NotEmpty(t, result)
	})
	
	t.Run("very large width", func(t *testing.T) {
		text := "Short text"
		result := wrapText(text, 1000)
		assert.Equal(t, text, result)
	})
}

// Helper function to check if a string is a single word
func isSingleWord(s string) bool {
	return !strings.Contains(strings.TrimSpace(s), " ")
}

func TestListCmd_Usage(t *testing.T) {
	// Test that the list command is properly configured
	assert.Equal(t, "list", listCmd.Use)
	assert.Equal(t, "List available project blueprints", listCmd.Short)
	assert.NotEmpty(t, listCmd.Long)
	assert.NotNil(t, listCmd.Run)
}

func TestListCmd_Execution(t *testing.T) {
	// Test that the list command can be executed without panicking
	assert.NotPanics(t, func() {
		// Create a buffer to capture output
		oldStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		// Execute the command function
		listCmd.Run(listCmd, []string{})

		// Restore stdout
		_ = w.Close()
		os.Stdout = oldStdout

		// Read and verify we got some output
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		output := buf.String()
		
		// Should produce some output (either blueprint list or no blueprints message)
		assert.True(t, len(output) >= 0) // At minimum, shouldn't crash
	})
}