package security

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/francknouama/go-starter/pkg/types"
)

// InputSanitizer handles security validation and sanitization of user inputs
type InputSanitizer struct {
	pathValidator   *PathValidator
	moduleValidator *ModulePathValidator
	resourceLimiter *ResourceLimiter
}

// NewInputSanitizer creates a new input sanitizer
func NewInputSanitizer() *InputSanitizer {
	return &InputSanitizer{
		pathValidator:   NewPathValidator(),
		moduleValidator: NewModulePathValidator(),
		resourceLimiter: NewResourceLimiter(),
	}
}

// SanitizeProjectConfig validates and sanitizes a project configuration
func (s *InputSanitizer) SanitizeProjectConfig(config *types.ProjectConfig) error {
	// Validate project name
	if err := s.validateProjectName(config.Name); err != nil {
		return fmt.Errorf("invalid project name: %w", err)
	}

	// Validate module path
	if err := s.moduleValidator.ValidateModulePath(config.Module); err != nil {
		return fmt.Errorf("invalid module path: %w", err)
	}

	// Validate project type
	if err := s.validateProjectType(config.Type); err != nil {
		return fmt.Errorf("invalid project type: %w", err)
	}

	// Sanitize text fields
	config.Author = s.sanitizeTextField(config.Author)
	config.Email = s.sanitizeTextField(config.Email)
	config.License = s.sanitizeTextField(config.License)

	// Validate variables map
	if config.Variables != nil {
		for key, value := range config.Variables {
			sanitizedKey := s.sanitizeVariableName(key)
			sanitizedValue := s.sanitizeTextField(value)
			
			if sanitizedKey != key {
				delete(config.Variables, key)
				config.Variables[sanitizedKey] = sanitizedValue
			} else {
				config.Variables[key] = sanitizedValue
			}
		}
	}

	return nil
}

// ValidateOutputPath validates the output path for security issues
func (s *InputSanitizer) ValidateOutputPath(outputPath string) error {
	return s.pathValidator.ValidateOutputPath(outputPath)
}

// validateProjectName validates project name for security and format issues
func (s *InputSanitizer) validateProjectName(name string) error {
	if name == "" {
		return types.NewValidationError("project name cannot be empty", nil)
	}

	// Check length limits
	if len(name) > 100 {
		return types.NewValidationError("project name too long (max 100 characters)", nil)
	}

	if len(name) < 2 {
		return types.NewValidationError("project name too short (min 2 characters)", nil)
	}

	// Check for dangerous characters
	dangerousChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\x00"}
	for _, char := range dangerousChars {
		if strings.Contains(name, char) {
			return types.NewValidationError(fmt.Sprintf("project name contains dangerous character: %s", char), nil)
		}
	}

	// Check for path traversal attempts
	if strings.Contains(name, "..") || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return types.NewValidationError("project name contains path characters", nil)
	}

	// Check for reserved names
	reservedNames := []string{"CON", "PRN", "AUX", "NUL", "COM1", "COM2", "COM3", "COM4", "COM5", "COM6", "COM7", "COM8", "COM9", "LPT1", "LPT2", "LPT3", "LPT4", "LPT5", "LPT6", "LPT7", "LPT8", "LPT9"}
	upperName := strings.ToUpper(name)
	for _, reserved := range reservedNames {
		if upperName == reserved {
			return types.NewValidationError(fmt.Sprintf("project name '%s' is reserved", name), nil)
		}
	}

	// Must contain at least one letter
	hasLetter := false
	for _, r := range name {
		if unicode.IsLetter(r) {
			hasLetter = true
			break
		}
	}
	if !hasLetter {
		return types.NewValidationError("project name must contain at least one letter", nil)
	}

	return nil
}

// validateProjectType validates the project type
func (s *InputSanitizer) validateProjectType(projectType string) error {
	if projectType == "" {
		return types.NewValidationError("project type cannot be empty", nil)
	}

	// Whitelist of allowed project types
	allowedTypes := map[string]bool{
		"web-api":      true,
		"cli":          true,
		"library":      true,
		"lambda":       true,
		"lambda-proxy": true,
		"event-driven": true,
		"microservice": true,
		"monolith":     true,
		"workspace":    true,
	}

	if !allowedTypes[projectType] {
		return types.NewValidationError(fmt.Sprintf("unknown project type: %s", projectType), nil)
	}

	return nil
}

// sanitizeTextField sanitizes a text field by removing dangerous characters
func (s *InputSanitizer) sanitizeTextField(text string) string {
	if text == "" {
		return text
	}

	// Remove null bytes and other control characters
	cleaned := strings.ReplaceAll(text, "\x00", "")
	
	// Remove common injection patterns (case insensitive)
	patterns := []string{
		"<script",
		"</script>",
		"javascript:",
		"data:",
		"vbscript:",
	}

	lowerCleaned := strings.ToLower(cleaned)
	for _, pattern := range patterns {
		// Find and remove pattern while preserving case for non-matching parts
		for {
			index := strings.Index(lowerCleaned, pattern)
			if index == -1 {
				break
			}
			// Remove the pattern from both original and lowercase versions
			cleaned = cleaned[:index] + cleaned[index+len(pattern):]
			lowerCleaned = lowerCleaned[:index] + lowerCleaned[index+len(pattern):]
		}
	}

	// Limit length
	if len(cleaned) > 255 {
		cleaned = cleaned[:255]
	}

	return strings.TrimSpace(cleaned)
}

// sanitizeVariableName sanitizes a variable name
func (s *InputSanitizer) sanitizeVariableName(name string) string {
	// Remove dangerous characters, keep only alphanumeric and underscore
	var result strings.Builder
	for _, r := range name {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' {
			result.WriteRune(r)
		}
	}
	
	sanitized := result.String()
	
	// Ensure it starts with a letter
	if len(sanitized) > 0 && !unicode.IsLetter(rune(sanitized[0])) {
		sanitized = "var_" + sanitized
	}
	
	// Limit length
	if len(sanitized) > 64 {
		sanitized = sanitized[:64]
	}
	
	return sanitized
}

// PathValidator validates file and directory paths
type PathValidator struct {
	maxPathLength int
}

// NewPathValidator creates a new path validator
func NewPathValidator() *PathValidator {
	return &PathValidator{
		maxPathLength: 260, // Windows MAX_PATH limit
	}
}

// ValidateOutputPath validates an output path for security issues
func (p *PathValidator) ValidateOutputPath(outputPath string) error {
	if outputPath == "" {
		return types.NewValidationError("output path cannot be empty", nil)
	}

	// Check path length
	if len(outputPath) > p.maxPathLength {
		return types.NewValidationError(fmt.Sprintf("output path too long (max %d characters)", p.maxPathLength), nil)
	}

	// Check for URL-encoded path traversal attempts
	urlEncodedPatterns := []string{
		"%2e%2e%2f", "%2e%2e/", "..%2f", "%2e%2e%5c", "%2e%2e\\",
		"%252e%252e%252f", "....//", "..../", "....\\\\",
	}
	lowerPath := strings.ToLower(outputPath)
	for _, pattern := range urlEncodedPatterns {
		if strings.Contains(lowerPath, pattern) {
			return types.NewValidationError("output path contains path traversal attempt", nil)
		}
	}

	// Clean and normalize the path
	cleanPath := filepath.Clean(outputPath)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return types.NewValidationError("output path contains path traversal attempt", nil)
	}

	// Check for dangerous characters
	dangerousChars := []string{"<", ">", ":", "\"", "|", "?", "*", "\x00"}
	for _, char := range dangerousChars {
		if strings.Contains(outputPath, char) {
			return types.NewValidationError(fmt.Sprintf("output path contains dangerous character: %s", char), nil)
		}
	}

	return nil
}

// ModulePathValidator validates Go module paths
type ModulePathValidator struct {
	modulePathRegex *regexp.Regexp
}

// NewModulePathValidator creates a new module path validator
func NewModulePathValidator() *ModulePathValidator {
	// Basic regex for Go module paths
	moduleRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-._/])*[a-zA-Z0-9]$`)
	
	return &ModulePathValidator{
		modulePathRegex: moduleRegex,
	}
}

// ValidateModulePath validates a Go module path
func (m *ModulePathValidator) ValidateModulePath(modulePath string) error {
	if modulePath == "" {
		return types.NewValidationError("module path cannot be empty", nil)
	}

	// Check length limits
	if len(modulePath) > 255 {
		return types.NewValidationError("module path too long (max 255 characters)", nil)
	}

	// Check basic format
	if !m.modulePathRegex.MatchString(modulePath) {
		return types.NewValidationError("module path format is invalid", nil)
	}

	// Check for dangerous patterns
	dangerousPatterns := []string{
		"..",
		"//",
		"\\",
		"@",
		" ",
	}

	for _, pattern := range dangerousPatterns {
		if strings.Contains(modulePath, pattern) {
			return types.NewValidationError(fmt.Sprintf("module path contains dangerous pattern: %s", pattern), nil)
		}
	}

	// Check for common malicious domains
	maliciousDomains := []string{
		"localhost",
		"127.0.0.1",
		"0.0.0.0",
		"169.254.169.254", // AWS metadata service
		"metadata.google.internal", // GCP metadata service
	}

	lowerPath := strings.ToLower(modulePath)
	for _, domain := range maliciousDomains {
		if strings.HasPrefix(lowerPath, domain) {
			return types.NewValidationError(fmt.Sprintf("module path uses potentially dangerous domain: %s", domain), nil)
		}
	}

	return nil
}

// ResourceLimiter enforces resource usage limits
type ResourceLimiter struct {
	maxFileSize     int64
	maxFiles        int
	maxDirectories  int
}

// NewResourceLimiter creates a new resource limiter
func NewResourceLimiter() *ResourceLimiter {
	return &ResourceLimiter{
		maxFileSize:    10 * 1024 * 1024, // 10MB per file
		maxFiles:       1000,              // Max files per project
		maxDirectories: 100,               // Max directories per project
	}
}

// ValidateResourceUsage validates resource usage limits
func (r *ResourceLimiter) ValidateResourceUsage(fileCount, dirCount int, totalSize int64) error {
	if fileCount > r.maxFiles {
		return types.NewValidationError(fmt.Sprintf("too many files: %d (max %d)", fileCount, r.maxFiles), nil)
	}

	if dirCount > r.maxDirectories {
		return types.NewValidationError(fmt.Sprintf("too many directories: %d (max %d)", dirCount, r.maxDirectories), nil)
	}

	if totalSize > r.maxFileSize*int64(r.maxFiles) {
		return types.NewValidationError("total project size exceeds limits", nil)
	}

	return nil
}

// ValidateFileSize validates a single file size
func (r *ResourceLimiter) ValidateFileSize(size int64) error {
	if size > r.maxFileSize {
		return types.NewValidationError(fmt.Sprintf("file size %d exceeds maximum %d bytes", size, r.maxFileSize), nil)
	}
	return nil
}