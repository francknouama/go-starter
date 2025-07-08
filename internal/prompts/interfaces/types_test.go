package interfaces

import (
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelectionItem(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		value       string
	}{
		{
			name:        "basic item",
			title:       "Option 1",
			description: "First option",
			value:       "opt1",
		},
		{
			name:        "empty description",
			title:       "Option 2",
			description: "",
			value:       "opt2",
		},
		{
			name:        "empty value",
			title:       "Option 3",
			description: "Third option",
			value:       "",
		},
		{
			name:        "all empty",
			title:       "",
			description: "",
			value:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewSelectionItem(tt.title, tt.description, tt.value)
			
			assert.Equal(t, tt.title, item.title)
			assert.Equal(t, tt.description, item.description)
			assert.Equal(t, tt.value, item.value)
		})
	}
}

func TestSelectionItem_Methods(t *testing.T) {
	item := NewSelectionItem("Test Title", "Test Description", "test_value")
	
	t.Run("FilterValue", func(t *testing.T) {
		assert.Equal(t, "Test Title", item.FilterValue())
	})
	
	t.Run("Title", func(t *testing.T) {
		assert.Equal(t, "Test Title", item.Title())
	})
	
	t.Run("Description", func(t *testing.T) {
		assert.Equal(t, "Test Description", item.Description())
	})
	
	t.Run("Value", func(t *testing.T) {
		assert.Equal(t, "test_value", item.Value())
	})
}

func TestSelectionItem_EmptyValues(t *testing.T) {
	item := NewSelectionItem("", "", "")
	
	assert.Equal(t, "", item.FilterValue())
	assert.Equal(t, "", item.Title())
	assert.Equal(t, "", item.Description())
	assert.Equal(t, "", item.Value())
}

func TestSelectionItem_SpecialCharacters(t *testing.T) {
	item := NewSelectionItem(
		"Title with spaces & symbols!",
		"Description with émojis 🚀 and unicode ñ",
		"value-with_special.chars",
	)
	
	assert.Equal(t, "Title with spaces & symbols!", item.Title())
	assert.Equal(t, "Description with émojis 🚀 and unicode ñ", item.Description())
	assert.Equal(t, "value-with_special.chars", item.Value())
	assert.Equal(t, "Title with spaces & symbols!", item.FilterValue())
}

func TestRealSurveyAdapter(t *testing.T) {
	adapter := &RealSurveyAdapter{}
	
	// Test that the adapter implements the interface
	var _ SurveyAdapter = adapter
	
	// Test AskOne with a simple prompt
	// Note: This test just verifies the method exists and has the right signature
	// We can't actually test the interactive behavior in a unit test
	prompt := &survey.Input{Message: "Test question"}
	var response string
	
	// We expect this to fail in a non-interactive environment, but that's okay
	// The important thing is that the method signature is correct
	err := adapter.AskOne(prompt, &response)
	
	// We don't assert on the error since it depends on the execution environment
	// The test just verifies the method can be called
	_ = err
}

func TestSurveyAdapter_Interface(t *testing.T) {
	// Test that RealSurveyAdapter implements SurveyAdapter
	var adapter SurveyAdapter = &RealSurveyAdapter{}
	require.NotNil(t, adapter)
	
	// Verify the interface method exists
	prompt := &survey.Input{Message: "Test"}
	var response string
	
	// Call the method (may fail in test environment, but that's expected)
	_ = adapter.AskOne(prompt, &response)
}

func TestSelectionItem_ComprehensiveCoverage(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		value       string
	}{
		{"normal case", "Title", "Desc", "val"},
		{"long strings", "Very Long Title That Goes On And On", "Very Long Description That Provides Lots Of Detail", "very_long_value_name"},
		{"numeric content", "123", "456", "789"},
		{"mixed content", "Title123", "Desc with numbers 456", "val_789"},
		{"single characters", "T", "D", "V"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item := NewSelectionItem(tt.title, tt.description, tt.value)
			
			// Test all methods return expected values
			assert.Equal(t, tt.title, item.Title())
			assert.Equal(t, tt.title, item.FilterValue()) // FilterValue returns title
			assert.Equal(t, tt.description, item.Description())
			assert.Equal(t, tt.value, item.Value())
			
			// Test that the item is properly constructed
			assert.Equal(t, tt.title, item.title)
			assert.Equal(t, tt.description, item.description)
			assert.Equal(t, tt.value, item.value)
		})
	}
}