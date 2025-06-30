package prompts

import (
	"fmt"
	"regexp"
	"strings"
)

// Validation functions for user input

// ValidateProjectName validates the project name
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	
	// Check length
	if len(name) < 2 {
		return fmt.Errorf("project name must be at least 2 characters")
	}
	if len(name) > 50 {
		return fmt.Errorf("project name must be less than 50 characters")
	}
	
	// Check format (letters, numbers, hyphens, underscores)
	validName := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
	if !validName.MatchString(name) {
		return fmt.Errorf("project name must start with a letter and contain only letters, numbers, hyphens, and underscores")
	}
	
	// Check for reserved names
	reserved := []string{"test", "main", "init", "vendor", "internal"}
	for _, r := range reserved {
		if strings.EqualFold(name, r) {
			return fmt.Errorf("'%s' is a reserved name", name)
		}
	}
	
	return nil
}

// ValidateModulePath validates the Go module path
func ValidateModulePath(path string) error {
	if path == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	
	// Basic module path validation
	// Should be something like: github.com/user/project
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return fmt.Errorf("module path should include at least domain and project (e.g., github.com/user/project)")
	}
	
	// Check for common patterns
	validModule := regexp.MustCompile(`^[a-zA-Z0-9.-]+(/[a-zA-Z0-9._-]+)+$`)
	if !validModule.MatchString(path) {
		return fmt.Errorf("invalid module path format")
	}
	
	return nil
}

// ValidateGoVersion is defined in go_version.go

