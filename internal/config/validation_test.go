package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name          string
		projectName   string
		shouldError   bool
		errorContains string
	}{
		// Valid project names
		{
			name:        "valid simple name",
			projectName: "myproject",
			shouldError: false,
		},
		{
			name:        "valid name with hyphens",
			projectName: "my-awesome-project",
			shouldError: false,
		},
		{
			name:        "valid name with underscores",
			projectName: "my_awesome_project",
			shouldError: false,
		},
		{
			name:        "valid name with numbers",
			projectName: "project123",
			shouldError: false,
		},

		// Invalid project names
		{
			name:          "empty name should error",
			projectName:   "",
			shouldError:   true,
			errorContains: "project name cannot be empty",
		},
		{
			name:          "name starting with hyphen should error",
			projectName:   "-myproject",
			shouldError:   true,
			errorContains: "cannot start or end with hyphen or underscore",
		},
		{
			name:          "name ending with hyphen should error",
			projectName:   "myproject-",
			shouldError:   true,
			errorContains: "cannot start or end with hyphen or underscore",
		},
		{
			name:          "name with spaces should error",
			projectName:   "my project",
			shouldError:   true,
			errorContains: "can only contain letters, numbers, hyphens, and underscores",
		},
		{
			name:          "reserved name should error",
			projectName:   "con",
			shouldError:   true,
			errorContains: "is reserved",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateProjectName(tt.projectName)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateModulePath(t *testing.T) {
	tests := []struct {
		name          string
		modulePath    string
		shouldError   bool
		errorContains string
	}{
		// Valid module paths
		{
			name:        "valid github path",
			modulePath:  "github.com/user/repo",
			shouldError: false,
		},
		{
			name:        "valid gitlab path",
			modulePath:  "gitlab.com/user/repo",
			shouldError: false,
		},
		{
			name:        "valid domain with subpath",
			modulePath:  "example.com/my/project",
			shouldError: false,
		},

		// Invalid module paths
		{
			name:          "empty path should error",
			modulePath:    "",
			shouldError:   true,
			errorContains: "module path cannot be empty",
		},
		{
			name:          "single word should error",
			modulePath:    "myproject",
			shouldError:   true,
			errorContains: "should contain at least domain and path",
		},
		{
			name:          "no domain should error",
			modulePath:    "user/repo",
			shouldError:   true,
			errorContains: "should start with a domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateModulePath(tt.modulePath)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		shouldError   bool
		errorContains string
	}{
		// Valid emails
		{
			name:        "valid email",
			email:       "user@example.com",
			shouldError: false,
		},
		{
			name:        "empty email is allowed (optional)",
			email:       "",
			shouldError: false,
		},
		{
			name:        "email with subdomain",
			email:       "user@mail.example.com",
			shouldError: false,
		},

		// Invalid emails
		{
			name:          "malformed email should error",
			email:         "notanemail",
			shouldError:   true,
			errorContains: "invalid email address",
		},
		{
			name:          "email without domain should error",
			email:         "user@",
			shouldError:   true,
			errorContains: "invalid email address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateEmail(tt.email)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateLogger(t *testing.T) {
	tests := []struct {
		name          string
		logger        string
		shouldError   bool
		errorContains string
	}{
		// Valid loggers
		{
			name:        "valid slog logger",
			logger:      "slog",
			shouldError: false,
		},
		{
			name:        "valid zap logger",
			logger:      "zap",
			shouldError: false,
		},
		{
			name:        "valid logrus logger",
			logger:      "logrus",
			shouldError: false,
		},
		{
			name:        "valid zerolog logger",
			logger:      "zerolog",
			shouldError: false,
		},
		{
			name:        "empty logger is allowed (defaults to slog)",
			logger:      "",
			shouldError: false,
		},

		// Invalid loggers
		{
			name:          "unsupported logger should error",
			logger:        "winston",
			shouldError:   true,
			errorContains: "invalid logger",
		},
		{
			name:          "case sensitive logger should error",
			logger:        "ZAP",
			shouldError:   true,
			errorContains: "invalid logger",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateLogger(tt.logger)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateAuthor(t *testing.T) {
	tests := []struct {
		name          string
		author        string
		shouldError   bool
		errorContains string
	}{
		// Valid authors
		{
			name:        "valid author name",
			author:      "John Doe",
			shouldError: false,
		},
		{
			name:        "empty author is allowed (optional)",
			author:      "",
			shouldError: false,
		},
		{
			name:        "author with special characters",
			author:      "José María-García",
			shouldError: false,
		},

		// Invalid authors  
		{
			name:          "author with control characters should error",
			author:        "John\x00Doe",
			shouldError:   true,
			errorContains: "invalid characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateAuthor(tt.author)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateOutputPath(t *testing.T) {
	tests := []struct {
		name          string
		outputPath    string
		shouldError   bool
		errorContains string
	}{
		// Valid output paths
		{
			name:        "valid relative path",
			outputPath:  "./my-project",
			shouldError: false,
		},
		{
			name:        "valid absolute Unix path",
			outputPath:  "/tmp/my-project",
			shouldError: false,
		},
		{
			name:        "simple directory name",
			outputPath:  "my-project",
			shouldError: false,
		},

		// Invalid output paths
		{
			name:          "empty path should error",
			outputPath:    "",
			shouldError:   true,
			errorContains: "output path cannot be empty",
		},
		{
			name:          "path traversal should error",
			outputPath:    "../../../etc/passwd",
			shouldError:   true,
			errorContains: "path traversal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateOutputPath(tt.outputPath)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateGoVersion(t *testing.T) {
	tests := []struct {
		name          string
		goVersion     string
		shouldError   bool
		errorContains string
	}{
		// Valid Go versions
		{
			name:        "valid Go 1.21",
			goVersion:   "1.21",
			shouldError: false,
		},
		{
			name:        "valid Go 1.22",
			goVersion:   "1.22",
			shouldError: false,
		},
		{
			name:        "valid Go 1.23",
			goVersion:   "1.23",
			shouldError: false,
		},
		{
			name:        "valid Go with patch version",
			goVersion:   "1.21.5",
			shouldError: false,
		},

		// Invalid Go versions
		{
			name:          "empty version should error",
			goVersion:     "",
			shouldError:   true,
			errorContains: "go version cannot be empty",
		},
		{
			name:          "unsupported old version should error",
			goVersion:     "1.17",
			shouldError:   true,
			errorContains: "unsupported Go version",
		},
		{
			name:          "invalid format should error",
			goVersion:     "2.0",
			shouldError:   true,
			errorContains: "invalid Go version format",
		},
		{
			name:          "malformed version should error",
			goVersion:     "1.x",
			shouldError:   true,
			errorContains: "invalid Go version format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateGoVersion(tt.goVersion)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateTemplateType(t *testing.T) {
	tests := []struct {
		name          string
		templateType  string
		shouldError   bool
		errorContains string
	}{
		// Valid template types
		{
			name:         "valid web-api template",
			templateType: "web-api",
			shouldError:  false,
		},
		{
			name:         "valid cli template",
			templateType: "cli",
			shouldError:  false,
		},
		{
			name:         "valid library template",
			templateType: "library",
			shouldError:  false,
		},
		{
			name:         "valid lambda template",
			templateType: "lambda",
			shouldError:  false,
		},

		// Invalid template types
		{
			name:          "empty template type should error",
			templateType:  "",
			shouldError:   true,
			errorContains: "invalid template type",
		},
		{
			name:          "unsupported template type should error",
			templateType:  "desktop",
			shouldError:   true,
			errorContains: "invalid template type",
		},
		{
			name:          "case sensitive template type should error",
			templateType:  "WEB-API",
			shouldError:   true,
			errorContains: "invalid template type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateTemplateType(tt.templateType)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateFramework(t *testing.T) {
	tests := []struct {
		name          string
		framework     string
		shouldError   bool
		errorContains string
	}{
		// Valid frameworks
		{
			name:        "valid gin framework",
			framework:   "gin",
			shouldError: false,
		},
		{
			name:        "valid echo framework",
			framework:   "echo",
			shouldError: false,
		},
		{
			name:        "valid fiber framework",
			framework:   "fiber",
			shouldError: false,
		},
		{
			name:        "valid chi framework",
			framework:   "chi",
			shouldError: false,
		},
		{
			name:        "valid cobra framework",
			framework:   "cobra",
			shouldError: false,
		},

		// Empty framework (allowed)
		{
			name:        "empty framework is allowed",
			framework:   "",
			shouldError: false,
		},

		// Invalid frameworks
		{
			name:          "unsupported framework",
			framework:     "express",
			shouldError:   true,
			errorContains: "invalid framework",
		},
		{
			name:          "case sensitive framework",
			framework:     "GIN",
			shouldError:   true,
			errorContains: "invalid framework",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Act
			err := ValidateFramework(tt.framework)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}