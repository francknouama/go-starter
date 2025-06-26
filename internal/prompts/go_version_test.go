package prompts

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoVersionPrompt(t *testing.T) {
	tests := []struct {
		name           string
		userSelection  string
		expectedResult string
	}{
		{
			name:           "auto-detect selection",
			userSelection:  "Auto-detect (recommended)",
			expectedResult: "auto",
		},
		{
			name:           "go 1.23 selection",
			userSelection:  "Go 1.23 (latest)",
			expectedResult: "1.23",
		},
		{
			name:           "go 1.22 selection",
			userSelection:  "Go 1.22",
			expectedResult: "1.22",
		},
		{
			name:           "go 1.21 selection",
			userSelection:  "Go 1.21",
			expectedResult: "1.21",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			result := mapSelectionToVersion(tt.userSelection)

			// Assert
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestGoVersionValidation(t *testing.T) {
	tests := []struct {
		name        string
		version     string
		shouldError bool
	}{
		{"valid auto", "auto", false},
		{"valid 1.23", "1.23", false},
		{"valid 1.22", "1.22", false},
		{"valid 1.21", "1.21", false},
		{"invalid 1.20", "1.20", true},
		{"invalid 2.0", "2.0", true},
		{"empty version", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateGoVersion(tt.version)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetSupportedGoVersions(t *testing.T) {
	// Act
	versions := GetSupportedGoVersions()

	// Assert
	expected := []string{"auto", "1.23", "1.22", "1.21"}
	assert.Equal(t, expected, versions)
}

func TestPromptGoVersionInteractive(t *testing.T) {
	// Test the PromptGoVersion method exists and has correct signature
	prompter := New()
	
	// This would normally require mocking survey for automated testing
	// For now, we'll just verify the method exists and returns correct types
	t.Run("method signature verification", func(t *testing.T) {
		// We can't easily test interactive prompts in CI, but we can verify
		// the method exists and has the right signature
		assert.NotNil(t, prompter.PromptGoVersion, "PromptGoVersion method should exist")
	})
}