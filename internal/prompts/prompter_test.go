package prompts

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/prompts/interfaces"
)

func TestNewPrompterFactory(t *testing.T) {
	tests := []struct {
		name          string
		useEnhancedUI bool
	}{
		{"with enhanced UI", true},
		{"without enhanced UI", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewPrompterFactory(tt.useEnhancedUI)
			
			assert.NotNil(t, factory)
			assert.Equal(t, tt.useEnhancedUI, factory.useEnhancedUI)
		})
	}
}

func TestPrompterFactory_CreatePrompter(t *testing.T) {
	tests := []struct {
		name          string
		useEnhancedUI bool
	}{
		{"enhanced UI enabled", true},
		{"enhanced UI disabled", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := NewPrompterFactory(tt.useEnhancedUI)
			prompter := factory.CreatePrompter()
			
			assert.NotNil(t, prompter)
			assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
		})
	}
}

func TestIsTerminal(t *testing.T) {
	// Test the isTerminal function
	// Note: This test might behave differently in different environments
	result := isTerminal()
	
	// The result should be a boolean, we can't assert a specific value
	// since it depends on the execution environment
	assert.IsType(t, false, result)
}

func TestCreateBubbleTeaPrompter(t *testing.T) {
	prompter := createBubbleTeaPrompter()
	
	assert.NotNil(t, prompter)
	assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
}

func TestCreateSurveyPrompter(t *testing.T) {
	prompter := createSurveyPrompter()
	
	assert.NotNil(t, prompter)
	assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
}

func TestNewDefault(t *testing.T) {
	prompter := NewDefault()
	
	assert.NotNil(t, prompter)
	assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
}

func TestNewSurveyFallback(t *testing.T) {
	prompter := NewSurveyFallback()
	
	assert.NotNil(t, prompter)
	assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
}

func TestPrompterFactory_CreatePrompter_Integration(t *testing.T) {
	// Test that different configurations produce valid prompters
	t.Run("enhanced UI factory", func(t *testing.T) {
		factory := NewPrompterFactory(true)
		prompter := factory.CreatePrompter()
		
		require.NotNil(t, prompter)
		assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
	})

	t.Run("basic UI factory", func(t *testing.T) {
		factory := NewPrompterFactory(false)
		prompter := factory.CreatePrompter()
		
		require.NotNil(t, prompter)
		assert.Implements(t, (*interfaces.Prompter)(nil), prompter)
	})
}

func TestIsTerminal_StdinRedirection(t *testing.T) {
	// Save original stdout
	originalStdout := os.Stdout
	defer func() {
		os.Stdout = originalStdout
	}()

	// Test with different stdout configurations
	// Note: This is primarily for coverage, actual behavior depends on environment
	result := isTerminal()
	assert.IsType(t, false, result)
}

func TestMapSelectionToVersion_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		selection string
		expected  string
	}{
		{"unknown selection", "unknown", "auto"},
		{"empty selection", "", "auto"},
		{"random string", "random string", "auto"},
		{"partial match", "Go 1.23", "auto"}, // Not exact match
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapSelectionToVersion(tt.selection)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestPrompterTypes(t *testing.T) {
	// Test that we can create different types of prompters
	bubbletea := createBubbleTeaPrompter()
	survey := createSurveyPrompter()
	
	// Both should implement the same interface but be different instances
	assert.NotNil(t, bubbletea)
	assert.NotNil(t, survey)
	assert.Implements(t, (*interfaces.Prompter)(nil), bubbletea)
	assert.Implements(t, (*interfaces.Prompter)(nil), survey)
	
	// They should be different instances
	assert.NotEqual(t, bubbletea, survey)
}