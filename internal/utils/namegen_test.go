package utils

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateRandomProjectName(t *testing.T) {
	// Test basic generation
	name := GenerateRandomProjectName()
	
	// Should not be empty
	assert.NotEmpty(t, name)
	
	// Should contain a hyphen (adjective-noun format)
	assert.Contains(t, name, "-")
	
	// Should have exactly one hyphen
	parts := strings.Split(name, "-")
	assert.Len(t, parts, 2, "Name should have exactly one hyphen separating adjective and noun")
	
	// Both parts should be non-empty
	assert.NotEmpty(t, parts[0], "Adjective part should not be empty")
	assert.NotEmpty(t, parts[1], "Noun part should not be empty")
	
	// Should match valid project name pattern (alphanumeric and hyphens only)
	validPattern := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	assert.True(t, validPattern.MatchString(name), "Generated name should only contain alphanumeric characters and hyphens")
	
	// Should not start or end with hyphen
	assert.False(t, strings.HasPrefix(name, "-"), "Name should not start with hyphen")
	assert.False(t, strings.HasSuffix(name, "-"), "Name should not end with hyphen")
}

func TestGenerateRandomProjectNameUniqueness(t *testing.T) {
	// Generate multiple names and check for some diversity
	names := make(map[string]bool)
	const iterations = 100
	
	for i := 0; i < iterations; i++ {
		name := GenerateRandomProjectName()
		names[name] = true
	}
	
	// Should have generated some variety (at least 80% unique)
	uniqueCount := len(names)
	expectedMinUnique := int(float64(iterations) * 0.8)
	assert.GreaterOrEqual(t, uniqueCount, expectedMinUnique, 
		"Should generate reasonable variety in names (at least 80%% unique)")
}

func TestGenerateMultipleNames(t *testing.T) {
	// Test with default count
	names := GenerateMultipleNames(0)
	assert.Len(t, names, 3, "Should default to 3 names when count is 0")
	
	// Test with specific count
	count := 5
	names = GenerateMultipleNames(count)
	assert.Len(t, names, count, "Should generate exactly the requested number of names")
	
	// All names should be unique within the set
	nameSet := make(map[string]bool)
	for _, name := range names {
		assert.False(t, nameSet[name], "Each generated name should be unique within the set")
		nameSet[name] = true
		
		// Each name should be valid
		assert.NotEmpty(t, name)
		assert.Contains(t, name, "-")
	}
}

func TestGenerateWithPrefix(t *testing.T) {
	prefix := "myapp"
	name := GenerateWithPrefix(prefix)
	
	// Should start with the prefix
	assert.True(t, strings.HasPrefix(name, prefix+"-"), 
		"Generated name should start with prefix followed by hyphen")
	
	// Should have the expected format: prefix-noun
	parts := strings.Split(name, "-")
	assert.Len(t, parts, 2, "Should have prefix and noun separated by hyphen")
	assert.Equal(t, prefix, parts[0], "First part should be the prefix")
	assert.NotEmpty(t, parts[1], "Second part (noun) should not be empty")
}

func TestGenerateWithSuffix(t *testing.T) {
	suffix := "api"
	name := GenerateWithSuffix(suffix)
	
	// Should end with the suffix
	assert.True(t, strings.HasSuffix(name, "-"+suffix), 
		"Generated name should end with hyphen followed by suffix")
	
	// Should have the expected format: adjective-suffix
	parts := strings.Split(name, "-")
	assert.Len(t, parts, 2, "Should have adjective and suffix separated by hyphen")
	assert.NotEmpty(t, parts[0], "First part (adjective) should not be empty")
	assert.Equal(t, suffix, parts[1], "Second part should be the suffix")
}

func TestGenerateWithEmptyPrefixSuffix(t *testing.T) {
	// Empty prefix should behave like normal generation
	nameWithEmptyPrefix := GenerateWithPrefix("")
	normalName := GenerateRandomProjectName()
	
	// Both should have the same format (adjective-noun)
	assert.Equal(t, 
		len(strings.Split(nameWithEmptyPrefix, "-")), 
		len(strings.Split(normalName, "-")),
		"Empty prefix should generate normal adjective-noun format")
	
	// Empty suffix should behave like normal generation
	nameWithEmptySuffix := GenerateWithSuffix("")
	assert.Equal(t, 
		len(strings.Split(nameWithEmptySuffix, "-")), 
		len(strings.Split(normalName, "-")),
		"Empty suffix should generate normal adjective-noun format")
}

func TestIsValidProjectNameChar(t *testing.T) {
	// Valid characters
	validChars := []rune{'a', 'Z', '5', '-', '_'}
	for _, ch := range validChars {
		assert.True(t, IsValidProjectNameChar(ch), 
			"Character '%c' should be valid", ch)
	}
	
	// Invalid characters
	invalidChars := []rune{' ', '!', '@', '#', '$', '%', '^', '&', '*', '(', ')', '=', '+', '[', ']', '{', '}', '|', '\\', ':', ';', '"', '\'', '<', '>', ',', '.', '?', '/'}
	for _, ch := range invalidChars {
		assert.False(t, IsValidProjectNameChar(ch), 
			"Character '%c' should be invalid", ch)
	}
}

func TestSanitizeProjectName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "hello-world",
			expected: "hello-world",
			name:     "already valid name should remain unchanged",
		},
		{
			input:    "hello world!",
			expected: "helloworld",
			name:     "spaces and special chars should be removed",
		},
		{
			input:    "my@awesome#project",
			expected: "myawesomeproject",
			name:     "special characters should be removed",
		},
		{
			input:    "test_project-123",
			expected: "test_project-123",
			name:     "underscores, hyphens, and numbers should be preserved",
		},
		{
			input:    "",
			expected: "",
			name:     "empty string should remain empty",
		},
		{
			input:    "!!!@@@###",
			expected: "",
			name:     "string with only invalid chars should become empty",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeProjectName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestNameGeneratorWithActualWordLists(t *testing.T) {
	// Test that our word lists contain expected words
	name := GenerateRandomProjectName()
	parts := strings.Split(name, "-")
	require.Len(t, parts, 2)
	
	adjective := parts[0]
	noun := parts[1]
	
	// Check that adjective exists in our list
	found := false
	for _, adj := range adjectives {
		if adj == adjective {
			found = true
			break
		}
	}
	assert.True(t, found, "Generated adjective should exist in adjectives list")
	
	// Check that noun exists in our list
	found = false
	for _, n := range nouns {
		if n == noun {
			found = true
			break
		}
	}
	assert.True(t, found, "Generated noun should exist in nouns list")
}

func BenchmarkGenerateRandomProjectName(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateRandomProjectName()
	}
}

func BenchmarkGenerateMultipleNames(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GenerateMultipleNames(5)
	}
}