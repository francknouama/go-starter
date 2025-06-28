package security

import (
	"strings"
	"testing"

	"github.com/francknouama/go-starter/internal/security"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPathTraversalPrevention(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	maliciousPaths := []string{
		"../../../etc/passwd",
		"..\\..\\windows\\system32",
		"%2e%2e%2f%2e%2e%2f",
		"..%252f..%252f",
		"....//....//",
		"..../....//",
		"\\..\\..\\system32",
		"./../../etc/shadow",
		"~/../../etc/hosts",
	}

	for _, path := range maliciousPaths {
		t.Run("blocks_"+path, func(t *testing.T) {
			err := sanitizer.ValidateOutputPath(path)
			assert.Error(t, err, "Should block path traversal attempt: %s", path)
			if err != nil {
				assert.Contains(t, err.Error(), "traversal", "Error should mention path traversal")
			}
		})
	}
}

func TestProjectNameValidation(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	testCases := []struct {
		name        string
		projectName string
		expectError bool
		reason      string
	}{
		{
			name:        "Valid simple name",
			projectName: "my-project",
			expectError: false,
		},
		{
			name:        "Valid with numbers",
			projectName: "project123",
			expectError: false,
		},
		{
			name:        "Empty name",
			projectName: "",
			expectError: true,
			reason:      "Empty names should be rejected",
		},
		{
			name:        "Too long name",
			projectName: strings.Repeat("a", 101),
			expectError: true,
			reason:      "Names over 100 chars should be rejected",
		},
		{
			name:        "Too short name",
			projectName: "a",
			expectError: true,
			reason:      "Names under 2 chars should be rejected",
		},
		{
			name:        "Dangerous characters",
			projectName: "project<script>",
			expectError: true,
			reason:      "Names with HTML tags should be rejected",
		},
		{
			name:        "Path traversal",
			projectName: "../etc",
			expectError: true,
			reason:      "Names with path characters should be rejected",
		},
		{
			name:        "Reserved name CON",
			projectName: "CON",
			expectError: true,
			reason:      "Windows reserved names should be rejected",
		},
		{
			name:        "Reserved name PRN",
			projectName: "PRN",
			expectError: true,
			reason:      "Windows reserved names should be rejected",
		},
		{
			name:        "Numbers only",
			projectName: "12345",
			expectError: true,
			reason:      "Names with only numbers should be rejected",
		},
		{
			name:        "Null byte",
			projectName: "project\x00name",
			expectError: true,
			reason:      "Names with null bytes should be rejected",
		},
		{
			name:        "Pipe character",
			projectName: "project|name",
			expectError: true,
			reason:      "Names with pipe characters should be rejected",
		},
		{
			name:        "Question mark",
			projectName: "project?name",
			expectError: true,
			reason:      "Names with question marks should be rejected",
		},
		{
			name:        "Asterisk",
			projectName: "project*name",
			expectError: true,
			reason:      "Names with asterisks should be rejected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:   tc.projectName,
				Module: "github.com/user/test",
				Type:   "web-api",
			}

			err := sanitizer.SanitizeProjectConfig(config)
			if tc.expectError {
				assert.Error(t, err, tc.reason)
			} else {
				assert.NoError(t, err, "Valid project name should be accepted")
			}
		})
	}
}

func TestModulePathValidation(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	testCases := []struct {
		name        string
		modulePath  string
		expectError bool
		reason      string
	}{
		{
			name:        "Valid GitHub module",
			modulePath:  "github.com/user/project",
			expectError: false,
		},
		{
			name:        "Valid GitLab module",
			modulePath:  "gitlab.com/user/project",
			expectError: false,
		},
		{
			name:        "Valid private module",
			modulePath:  "company.com/internal/project",
			expectError: false,
		},
		{
			name:        "Empty module path",
			modulePath:  "",
			expectError: true,
			reason:      "Empty module paths should be rejected",
		},
		{
			name:        "Too long module path",
			modulePath:  "github.com/" + strings.Repeat("a", 250),
			expectError: true,
			reason:      "Module paths over 255 chars should be rejected",
		},
		{
			name:        "Path traversal in module",
			modulePath:  "github.com/../../../etc/passwd",
			expectError: true,
			reason:      "Module paths with traversal should be rejected",
		},
		{
			name:        "Double slash",
			modulePath:  "github.com//user//project",
			expectError: true,
			reason:      "Module paths with double slashes should be rejected",
		},
		{
			name:        "Backslash",
			modulePath:  "github.com\\user\\project",
			expectError: true,
			reason:      "Module paths with backslashes should be rejected",
		},
		{
			name:        "At symbol",
			modulePath:  "github.com/user@evil.com/project",
			expectError: true,
			reason:      "Module paths with @ symbols should be rejected",
		},
		{
			name:        "Space character",
			modulePath:  "github.com/user name/project",
			expectError: true,
			reason:      "Module paths with spaces should be rejected",
		},
		{
			name:        "Localhost domain",
			modulePath:  "localhost/user/project",
			expectError: true,
			reason:      "Localhost module paths should be rejected",
		},
		{
			name:        "Local IP address",
			modulePath:  "127.0.0.1/user/project",
			expectError: true,
			reason:      "Local IP module paths should be rejected",
		},
		{
			name:        "AWS metadata service",
			modulePath:  "169.254.169.254/user/project",
			expectError: true,
			reason:      "AWS metadata service paths should be rejected",
		},
		{
			name:        "GCP metadata service",
			modulePath:  "metadata.google.internal/user/project",
			expectError: true,
			reason:      "GCP metadata service paths should be rejected",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:   "test-project",
				Module: tc.modulePath,
				Type:   "web-api",
			}

			err := sanitizer.SanitizeProjectConfig(config)
			if tc.expectError {
				assert.Error(t, err, tc.reason)
			} else {
				assert.NoError(t, err, "Valid module path should be accepted")
			}
		})
	}
}

func TestInputSanitization(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	testCases := []struct {
		field    string
		input    string
		expected string
	}{
		{
			field:    "Author",
			input:    "John Doe",
			expected: "John Doe",
		},
		{
			field:    "Author with script",
			input:    "<script>alert('xss')</script>John",
			expected: "alert('xss')John", // Script tags removed
		},
		{
			field:    "Email with null byte",
			input:    "user\x00@example.com",
			expected: "user@example.com", // Null byte removed
		},
		{
			field:    "License with javascript",
			input:    "javascript:alert('xss')",
			expected: "alert('xss')", // javascript: removed
		},
		{
			field:    "Long text",
			input:    strings.Repeat("a", 300),
			expected: strings.Repeat("a", 255), // Truncated to 255 chars
		},
		{
			field:    "Text with vbscript",
			input:    "vbscript:msgbox('xss')",
			expected: "msgbox('xss')", // vbscript: removed
		},
		{
			field:    "Text with data URI",
			input:    "data:text/html,<script>alert('xss')</script>",
			expected: "text/html,alert('xss')", // data: and script tags removed
		},
	}

	for _, tc := range testCases {
		t.Run(tc.field, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:    "test-project",
				Module:  "github.com/user/test",
				Type:    "web-api",
				Author:  tc.input,
				Email:   tc.input,
				License: tc.input,
			}

			err := sanitizer.SanitizeProjectConfig(config)
			require.NoError(t, err, "Sanitization should not fail")

			// Check that dangerous patterns are removed
			assert.NotContains(t, config.Author, "<script", "Author should not contain script tags")
			assert.NotContains(t, config.Email, "<script", "Email should not contain script tags")
			assert.NotContains(t, config.License, "<script", "License should not contain script tags")
			
			assert.NotContains(t, config.Author, "javascript:", "Author should not contain javascript:")
			assert.NotContains(t, config.Email, "javascript:", "Email should not contain javascript:")
			assert.NotContains(t, config.License, "javascript:", "License should not contain javascript:")
			
			assert.NotContains(t, config.Author, "\x00", "Author should not contain null bytes")
			assert.NotContains(t, config.Email, "\x00", "Email should not contain null bytes")
			assert.NotContains(t, config.License, "\x00", "License should not contain null bytes")
			
			// Check length limits
			assert.LessOrEqual(t, len(config.Author), 255, "Author should be truncated to 255 chars")
			assert.LessOrEqual(t, len(config.Email), 255, "Email should be truncated to 255 chars")
			assert.LessOrEqual(t, len(config.License), 255, "License should be truncated to 255 chars")
		})
	}
}

func TestVariableNameSanitization(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	testCases := []struct {
		name     string
		input    map[string]string
		expected map[string]string
	}{
		{
			name: "Valid variable names",
			input: map[string]string{
				"valid_name":  "value1",
				"another-var": "value2",
				"CamelCase":   "value3",
			},
			expected: map[string]string{
				"valid_name":  "value1",
				"another-var": "value2",
				"CamelCase":   "value3",
			},
		},
		{
			name: "Invalid characters",
			input: map[string]string{
				"var@name":    "value1",
				"var name":    "value2",
				"var#name":    "value3",
				"var$name":    "value4",
			},
			// Note: All these inputs will sanitize to the same key "varname"
			// so we expect only the last one to remain in the map
			expected: map[string]string{
				"varname": "value4", // Only the last value remains due to key collision
			},
		},
		{
			name: "Starting with number",
			input: map[string]string{
				"123var": "value1",
				"9name":  "value2",
			},
			expected: map[string]string{
				"var_123var": "value1",
				"var_9name":  "value2",
			},
		},
		{
			name: "Too long variable name",
			input: map[string]string{
				strings.Repeat("a", 100): "value1",
			},
			expected: map[string]string{
				strings.Repeat("a", 64): "value1", // Truncated to 64 chars
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:      "test-project",
				Module:    "github.com/user/test",
				Type:      "web-api",
				Variables: tc.input,
			}

			err := sanitizer.SanitizeProjectConfig(config)
			require.NoError(t, err, "Sanitization should not fail")

			// For the "Invalid characters" test case, we expect key collisions
			if tc.name == "Invalid characters" {
				assert.NotNil(t, config.Variables, "Variables should exist")
				assert.Len(t, config.Variables, 1, "Should have only one key due to sanitization collisions")
				assert.Contains(t, config.Variables, "varname", "Should contain sanitized key")
			} else {
				// For other test cases, check expected results
				for expectedKey, expectedValue := range tc.expected {
					assert.Equal(t, expectedValue, config.Variables[expectedKey], "Variable value should match expected")
				}
			}
			
			// Check that no dangerous characters remain in keys
			for key := range config.Variables {
				assert.NotContains(t, key, "@", "Sanitized key should not contain @")
				assert.NotContains(t, key, " ", "Sanitized key should not contain spaces")
				assert.NotContains(t, key, "#", "Sanitized key should not contain #")
				assert.NotContains(t, key, "$", "Sanitized key should not contain $")
			}
		})
	}
}

func TestProjectTypeValidation(t *testing.T) {
	sanitizer := security.NewInputSanitizer()

	validTypes := []string{
		"web-api",
		"cli",
		"library",
		"lambda",
		"lambda-proxy",
		"event-driven",
		"microservice",
		"monolith",
		"workspace",
	}

	for _, validType := range validTypes {
		t.Run("valid_"+validType, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/user/test",
				Type:   validType,
			}

			err := sanitizer.SanitizeProjectConfig(config)
			assert.NoError(t, err, "Valid project type should be accepted")
		})
	}

	invalidTypes := []string{
		"",
		"invalid-type",
		"web api", // with space
		"WEB-API", // wrong case
		"unknown",
		"hacker-type",
		"../traversal",
		"<script>",
	}

	for _, invalidType := range invalidTypes {
		t.Run("invalid_"+invalidType, func(t *testing.T) {
			config := &types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/user/test",
				Type:   invalidType,
			}

			err := sanitizer.SanitizeProjectConfig(config)
			assert.Error(t, err, "Invalid project type should be rejected: %s", invalidType)
		})
	}
}

func TestResourceLimits(t *testing.T) {
	limiter := security.NewResourceLimiter()

	testCases := []struct {
		name        string
		fileCount   int
		dirCount    int
		totalSize   int64
		expectError bool
		reason      string
	}{
		{
			name:        "Within limits",
			fileCount:   10,
			dirCount:    5,
			totalSize:   1024 * 1024, // 1MB
			expectError: false,
		},
		{
			name:        "Too many files",
			fileCount:   1001,
			dirCount:    5,
			totalSize:   1024,
			expectError: true,
			reason:      "Should reject projects with too many files",
		},
		{
			name:        "Too many directories",
			fileCount:   10,
			dirCount:    101,
			totalSize:   1024,
			expectError: true,
			reason:      "Should reject projects with too many directories",
		},
		{
			name:        "Total size too large",
			fileCount:   10,
			dirCount:    5,
			totalSize:   100 * 1024 * 1024 * 1024, // 100GB
			expectError: true,
			reason:      "Should reject projects that are too large",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := limiter.ValidateResourceUsage(tc.fileCount, tc.dirCount, tc.totalSize)
			if tc.expectError {
				assert.Error(t, err, tc.reason)
			} else {
				assert.NoError(t, err, "Valid resource usage should be accepted")
			}
		})
	}
}

func TestFileSizeValidation(t *testing.T) {
	limiter := security.NewResourceLimiter()

	testCases := []struct {
		name        string
		fileSize    int64
		expectError bool
	}{
		{
			name:        "Small file",
			fileSize:    1024,
			expectError: false,
		},
		{
			name:        "Medium file",
			fileSize:    1024 * 1024, // 1MB
			expectError: false,
		},
		{
			name:        "Large valid file",
			fileSize:    10 * 1024 * 1024, // 10MB (at limit)
			expectError: false,
		},
		{
			name:        "Too large file",
			fileSize:    11 * 1024 * 1024, // 11MB (over limit)
			expectError: true,
		},
		{
			name:        "Extremely large file",
			fileSize:    100 * 1024 * 1024, // 100MB
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := limiter.ValidateFileSize(tc.fileSize)
			if tc.expectError {
				assert.Error(t, err, "Large files should be rejected")
			} else {
				assert.NoError(t, err, "Valid file sizes should be accepted")
			}
		})
	}
}

func BenchmarkInputSanitization(b *testing.B) {
	sanitizer := security.NewInputSanitizer()
	config := &types.ProjectConfig{
		Name:    "test-project",
		Module:  "github.com/user/test-project",
		Type:    "web-api",
		Author:  "John Doe <john@example.com>",
		Email:   "john@example.com",
		License: "MIT",
		Variables: map[string]string{
			"var1": "value1",
			"var2": "value2",
			"var3": "value3",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create a copy to avoid modifying the original
		configCopy := *config
		if config.Variables != nil {
			configCopy.Variables = make(map[string]string)
			for k, v := range config.Variables {
				configCopy.Variables[k] = v
			}
		}
		_ = sanitizer.SanitizeProjectConfig(&configCopy)
	}
}

func BenchmarkPathValidation(b *testing.B) {
	sanitizer := security.NewInputSanitizer()
	path := "/home/user/projects/my-awesome-project"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = sanitizer.ValidateOutputPath(path)
	}
}