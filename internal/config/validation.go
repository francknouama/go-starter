package config

import (
	"fmt"
	"net/mail"
	"path/filepath"
	"regexp"
	"strings"
)

// ValidateProjectName validates a project name
func ValidateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}

	// Check length
	if len(name) > 214 {
		return fmt.Errorf("project name too long (max 214 characters)")
	}

	// Check for valid characters (alphanumeric, hyphens, underscores)
	validNameRegex := regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)
	if !validNameRegex.MatchString(name) {
		return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
	}

	// Cannot start or end with hyphen or underscore
	if strings.HasPrefix(name, "-") || strings.HasPrefix(name, "_") ||
		strings.HasSuffix(name, "-") || strings.HasSuffix(name, "_") {
		return fmt.Errorf("project name cannot start or end with hyphen or underscore")
	}

	// Check for reserved names
	reservedNames := map[string]bool{
		"con": true, "prn": true, "aux": true, "nul": true,
		"com1": true, "com2": true, "com3": true, "com4": true,
		"com5": true, "com6": true, "com7": true, "com8": true, "com9": true,
		"lpt1": true, "lpt2": true, "lpt3": true, "lpt4": true,
		"lpt5": true, "lpt6": true, "lpt7": true, "lpt8": true, "lpt9": true,
	}
	if reservedNames[strings.ToLower(name)] {
		return fmt.Errorf("project name '%s' is reserved", name)
	}

	return nil
}

// ValidateModulePath validates a Go module path
func ValidateModulePath(path string) error {
	if path == "" {
		return fmt.Errorf("module path cannot be empty")
	}

	// Check length
	if len(path) > 500 {
		return fmt.Errorf("module path too long (max 500 characters)")
	}

	// Basic module path validation
	// Should be in format like: github.com/user/repo or domain.com/path
	validModuleRegex := regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]*[a-zA-Z0-9])?)*(/[a-zA-Z0-9]([a-zA-Z0-9\-_]*[a-zA-Z0-9])?)*$`)
	if !validModuleRegex.MatchString(path) {
		return fmt.Errorf("invalid module path format")
	}

	// Check for common patterns
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return fmt.Errorf("module path should contain at least domain and path (e.g., github.com/user/repo)")
	}

	// First part should look like a domain
	domain := parts[0]
	if !strings.Contains(domain, ".") {
		return fmt.Errorf("module path should start with a domain (e.g., github.com, gitlab.com)")
	}

	return nil
}

// ValidateEmail validates an email address
func ValidateEmail(email string) error {
	if email == "" {
		return nil // Email is optional
	}

	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	return nil
}

// ValidateAuthor validates an author name
func ValidateAuthor(author string) error {
	if author == "" {
		return nil // Author is optional
	}

	// Check length
	if len(author) > 100 {
		return fmt.Errorf("author name too long (max 100 characters)")
	}

	// Author can contain most characters, but not control characters
	for _, char := range author {
		if char < 32 || char == 127 {
			return fmt.Errorf("author name contains invalid characters")
		}
	}

	return nil
}

// ValidateOutputPath validates an output directory path
func ValidateOutputPath(path string) error {
	if path == "" {
		return fmt.Errorf("output path cannot be empty")
	}

	// Clean the path
	cleanPath := filepath.Clean(path)

	// Check for path traversal attempts
	if strings.Contains(cleanPath, "..") {
		return fmt.Errorf("output path cannot contain '..' (path traversal)")
	}

	// Check for absolute paths on Windows
	if filepath.IsAbs(cleanPath) && len(cleanPath) > 1 && cleanPath[1] == ':' {
		// Windows absolute path - allow it
		return nil
	}

	// Check for Unix absolute paths
	if filepath.IsAbs(cleanPath) {
		// Unix absolute path - allow it
		return nil
	}

	// Relative paths are OK
	return nil
}

// ValidateGoVersion validates a Go version string
func ValidateGoVersion(version string) error {
	if version == "" {
		return fmt.Errorf("Go version cannot be empty")
	}

	// Should be in format "1.xx" or "1.xx.x"
	validVersionRegex := regexp.MustCompile(`^1\.\d+(\.\d+)?$`)
	if !validVersionRegex.MatchString(version) {
		return fmt.Errorf("invalid Go version format (expected format: 1.xx or 1.xx.x)")
	}

	// Extract major version number
	parts := strings.Split(version, ".")
	if len(parts) < 2 {
		return fmt.Errorf("invalid Go version format")
	}

	// We support Go 1.18 and later
	minorVersion := parts[1]
	switch minorVersion {
	case "18", "19", "20", "21", "22", "23":
		return nil
	default:
		// Check if it's a numeric value >= 18
		if len(minorVersion) == 2 && minorVersion >= "18" {
			return nil
		}
		if len(minorVersion) == 1 && minorVersion >= "18"[1:] {
			return nil
		}
		return fmt.Errorf("unsupported Go version (minimum supported: 1.18)")
	}
}

// ValidateTemplateType validates a template type
func ValidateTemplateType(templateType string) error {
	validTypes := map[string]bool{
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

	if !validTypes[templateType] {
		return fmt.Errorf("invalid template type '%s'", templateType)
	}

	return nil
}

// ValidateFramework validates a framework choice
func ValidateFramework(framework string) error {
	validFrameworks := map[string]bool{
		"gin":   true,
		"echo":  true,
		"fiber": true,
		"chi":   true,
		"cobra": true,
		"":      true, // empty is allowed for some templates
	}

	if !validFrameworks[framework] {
		return fmt.Errorf("invalid framework '%s'", framework)
	}

	return nil
}

// ValidateArchitecture validates an architecture pattern
func ValidateArchitecture(architecture string) error {
	validArchitectures := map[string]bool{
		"standard":     true,
		"clean":        true,
		"ddd":          true,
		"hexagonal":    true,
		"event-driven": true,
		"":             true, // empty is allowed
	}

	if !validArchitectures[architecture] {
		return fmt.Errorf("invalid architecture '%s'", architecture)
	}

	return nil
}

// ValidateDatabaseDriver validates a database driver
func ValidateDatabaseDriver(driver string) error {
	validDrivers := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"mongodb":  true,
		"sqlite":   true,
		"redis":    true,
		"":         true, // empty is allowed
	}

	if !validDrivers[driver] {
		return fmt.Errorf("invalid database driver '%s'", driver)
	}

	return nil
}

// ValidateORM validates an ORM choice
func ValidateORM(orm string) error {
	validORMs := map[string]bool{
		"gorm": true,
		"sqlx": true,
		"sqlc": true,
		"ent":  true,
		"":     true, // empty is allowed
	}

	if !validORMs[orm] {
		return fmt.Errorf("invalid ORM '%s'", orm)
	}

	return nil
}

// ValidateAuthType validates an authentication type
func ValidateAuthType(authType string) error {
	validAuthTypes := map[string]bool{
		"jwt":     true,
		"oauth2":  true,
		"session": true,
		"api-key": true,
		"":        true, // empty is allowed
	}

	if !validAuthTypes[authType] {
		return fmt.Errorf("invalid authentication type '%s'", authType)
	}

	return nil
}

// ValidateLogger validates a logger type
func ValidateLogger(logger string) error {
	validLoggers := map[string]bool{
		"slog":    true,
		"zap":     true,
		"logrus":  true,
		"zerolog": true,
		"":        true, // empty is allowed (will default to slog)
	}

	if !validLoggers[logger] {
		return fmt.Errorf("invalid logger '%s' (supported: slog, zap, logrus, zerolog)", logger)
	}

	return nil
}

// ValidateLogLevel validates a log level
func ValidateLogLevel(level string) error {
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"":      true, // empty is allowed (will default to info)
	}

	if !validLevels[level] {
		return fmt.Errorf("invalid log level '%s' (supported: debug, info, warn, error)", level)
	}

	return nil
}

// ValidateLogFormat validates a log format
func ValidateLogFormat(format string) error {
	validFormats := map[string]bool{
		"json":    true,
		"text":    true,
		"console": true,
		"":        true, // empty is allowed (will default to json)
	}

	if !validFormats[format] {
		return fmt.Errorf("invalid log format '%s' (supported: json, text, console)", format)
	}

	return nil
}

// ValidateFeatures validates project features configuration
func ValidateFeatures(features map[string]interface{}) error {
	// This is a flexible validation function that can be extended
	// as we add more feature types

	// For now, just check that it's not nil if provided
	if features == nil {
		return nil
	}

	// Could add specific feature validation here in the future
	// For example:
	// - Validate database configuration
	// - Validate authentication providers
	// - Validate deployment targets

	return nil
}
