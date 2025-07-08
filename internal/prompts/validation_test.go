package prompts

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name        string
		projectName string
		shouldError bool
		errorMsg    string
	}{
		// Valid cases
		{"valid simple name", "myproject", false, ""},
		{"valid with numbers", "project123", false, ""},
		{"valid with hyphens", "my-project", false, ""},
		{"valid with underscores", "my_project", false, ""},
		{"valid mixed", "my-project_123", false, ""},
		{"minimum length", "ab", false, ""},
		{"maximum length", "a123456789012345678901234567890123456789012345678", false, ""}, // 49 chars
		
		// Invalid cases - empty/length
		{"empty name", "", true, "project name cannot be empty"},
		{"too short", "a", true, "project name must be at least 2 characters"},
		{"too long", "a12345678901234567890123456789012345678901234567890", true, "project name must be less than 50 characters"}, // 51 chars
		
		// Invalid cases - format
		{"starts with number", "123project", true, "project name must start with a letter"},
		{"starts with hyphen", "-project", true, "project name must start with a letter"},
		{"starts with underscore", "_project", true, "project name must start with a letter"},
		{"contains spaces", "my project", true, "project name must start with a letter and contain only letters, numbers, hyphens, and underscores"},
		{"contains special chars", "my@project", true, "project name must start with a letter and contain only letters, numbers, hyphens, and underscores"},
		{"contains dots", "my.project", true, "project name must start with a letter and contain only letters, numbers, hyphens, and underscores"},
		
		// Invalid cases - reserved names
		{"reserved test", "test", true, "'test' is a reserved name"},
		{"reserved main", "main", true, "'main' is a reserved name"},
		{"reserved init", "init", true, "'init' is a reserved name"},
		{"reserved vendor", "vendor", true, "'vendor' is a reserved name"},
		{"reserved internal", "internal", true, "'internal' is a reserved name"},
		{"reserved case insensitive", "TEST", true, "'TEST' is a reserved name"},
		{"reserved case insensitive", "Main", true, "'Main' is a reserved name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.projectName)
			
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateModulePath(t *testing.T) {
	tests := []struct {
		name        string
		modulePath  string
		shouldError bool
		errorMsg    string
	}{
		// Valid cases
		{"github path", "github.com/user/project", false, ""},
		{"gitlab path", "gitlab.com/user/project", false, ""},
		{"custom domain", "example.com/user/project", false, ""},
		{"subdomain", "api.example.com/user/project", false, ""},
		{"deep path", "github.com/user/organization/project", false, ""},
		{"with numbers", "github.com/user123/project456", false, ""},
		{"with hyphens", "github.com/my-user/my-project", false, ""},
		{"with underscores", "github.com/my_user/my_project", false, ""},
		{"with dots in domain", "my.custom.domain.com/user/project", false, ""},
		
		// Invalid cases
		{"empty path", "", true, "module path cannot be empty"},
		{"single part", "project", true, "module path should include at least domain and project"},
		{"only domain", "github.com", true, "module path should include at least domain and project"},
		{"invalid format with spaces", "github.com/user name/project", true, "invalid module path format"},
		{"invalid format with @", "github.com/user@domain/project", true, "invalid module path format"},
		{"invalid format with special chars", "github.com/user$/project", true, "invalid module path format"},
		{"starts with slash", "/github.com/user/project", true, "invalid module path format"},
		{"ends with slash", "github.com/user/project/", true, "invalid module path format"},
		{"double slash", "github.com//user/project", true, "invalid module path format"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateModulePath(tt.modulePath)
			
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidationEdgeCases(t *testing.T) {
	t.Run("project name boundary length", func(t *testing.T) {
		// Test exact boundary conditions
		assert.NoError(t, ValidateProjectName("ab")) // exactly 2 chars
		assert.Error(t, ValidateProjectName("a"))   // exactly 1 char
		
		// 49 characters (should pass)
		longName49 := "a123456789012345678901234567890123456789012345678"
		assert.NoError(t, ValidateProjectName(longName49))
		
		// 50 characters (should fail)
		longName50 := "a12345678901234567890123456789012345678901234567890"
		assert.Error(t, ValidateProjectName(longName50))
	})

	t.Run("module path edge cases", func(t *testing.T) {
		// Test minimum valid path
		assert.NoError(t, ValidateModulePath("a.b/c"))
		
		// Test with numeric TLD
		assert.NoError(t, ValidateModulePath("example.123/user/project"))
		
		// Test with single character parts
		assert.NoError(t, ValidateModulePath("a.b/c/d"))
	})

	t.Run("reserved names case sensitivity", func(t *testing.T) {
		reserved := []string{"test", "main", "init", "vendor", "internal"}
		for _, name := range reserved {
			// Test lowercase
			assert.Error(t, ValidateProjectName(name))
			// Test uppercase
			assert.Error(t, ValidateProjectName(strings.ToUpper(name)))
			// Test mixed case
			if len(name) > 1 {
				mixed := strings.ToUpper(name[:1]) + name[1:]
				assert.Error(t, ValidateProjectName(mixed))
			}
		}
	})
}

func TestValidationRegexPatterns(t *testing.T) {
	t.Run("project name regex", func(t *testing.T) {
		// Test various combinations that should be valid
		validNames := []string{
			"a", // Would fail length check, but regex should pass
			"abc",
			"a1",
			"a_",
			"a-",
			"a1b2c3",
			"a_b_c",
			"a-b-c",
			"a1_b2-c3",
		}
		
		// Test that regex allows these patterns (ignoring other validation)
		for _, name := range validNames {
			if len(name) >= 2 && len(name) <= 50 {
				// Should pass if not reserved
				if !isReservedName(name) {
					assert.NoError(t, ValidateProjectName(name), "Should be valid: %s", name)
				}
			}
		}
	})

	t.Run("module path regex", func(t *testing.T) {
		// Test edge cases for module path regex
		validPaths := []string{
			"a.b/c",
			"a-b.c-d/e_f/g.h",
			"123.456/789",
			"a.b.c.d/e/f/g/h",
		}
		
		for _, path := range validPaths {
			assert.NoError(t, ValidateModulePath(path), "Should be valid: %s", path)
		}
	})
}

// Helper function to check if a name is reserved (extracted for testing)
func isReservedName(name string) bool {
	reserved := []string{"test", "main", "init", "vendor", "internal"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return true
		}
	}
	return false
}