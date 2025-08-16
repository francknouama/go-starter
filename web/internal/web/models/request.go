package models

import (
	"errors"
	"strings"
)

// GenerationRequest represents a request to generate a project
type GenerationRequest struct {
	// Project metadata
	Name       string `json:"name" binding:"required" example:"my-awesome-api"`
	ModulePath string `json:"modulePath" binding:"required" example:"github.com/user/my-awesome-api"`
	
	// Project configuration
	Type         string `json:"type" binding:"required" example:"web-api"`
	Architecture string `json:"architecture,omitempty" example:"clean"`
	Framework    string `json:"framework,omitempty" example:"gin"`
	GoVersion    string `json:"goVersion,omitempty" example:"1.21"`
	
	// Logger configuration
	Logger string `json:"logger,omitempty" example:"slog"`
	
	// Features configuration
	Features *FeaturesConfig `json:"features,omitempty"`
	
	// Generation options
	OutputFormat string `json:"outputFormat,omitempty" example:"zip"` // "zip" or "files"
	DryRun       bool   `json:"dryRun,omitempty"`
	
	// Progressive disclosure
	Complexity string `json:"complexity,omitempty" example:"standard"` // "simple", "standard", "advanced", "expert"
	Advanced   bool   `json:"advanced,omitempty"`
}

// FeaturesConfig represents the features configuration for a project
type FeaturesConfig struct {
	Database       *DatabaseConfig       `json:"database,omitempty"`
	Authentication *AuthenticationConfig `json:"authentication,omitempty"`
	Deployment     *DeploymentConfig     `json:"deployment,omitempty"`
	Testing        *TestingConfig        `json:"testing,omitempty"`
	Logging        *LoggingConfig        `json:"logging,omitempty"`
	Monitoring     *MonitoringConfig     `json:"monitoring,omitempty"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver string `json:"driver,omitempty" example:"postgres"`    // postgres, mysql, mongodb, sqlite, redis
	ORM    string `json:"orm,omitempty" example:"gorm"`          // gorm, sqlx, sqlc, ent
}

// AuthenticationConfig represents authentication configuration
type AuthenticationConfig struct {
	Type      string   `json:"type,omitempty" example:"jwt"`                    // jwt, oauth2, session, api-key
	Providers []string `json:"providers,omitempty" example:"google,github"`     // google, github, facebook, etc.
}

// DeploymentConfig represents deployment configuration
type DeploymentConfig struct {
	Targets []string `json:"targets,omitempty" example:"docker,kubernetes"`     // docker, kubernetes, lambda, etc.
}

// TestingConfig represents testing configuration
type TestingConfig struct {
	Framework string `json:"framework,omitempty" example:"testify"`            // testify, ginkgo
	Coverage  bool   `json:"coverage,omitempty"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level  string `json:"level,omitempty" example:"info"`                      // debug, info, warn, error
	Format string `json:"format,omitempty" example:"json"`                     // json, text, console
}

// MonitoringConfig represents monitoring configuration
type MonitoringConfig struct {
	Metrics bool `json:"metrics,omitempty"`
	Tracing bool `json:"tracing,omitempty"`
}

// PreviewRequest represents a request to preview a project structure
type PreviewRequest struct {
	GenerationRequest
}

// Validate validates the generation request
func (r *GenerationRequest) Validate() error {
	// Validate required fields
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("project name is required")
	}
	
	if strings.TrimSpace(r.ModulePath) == "" {
		return errors.New("module path is required")
	}
	
	if strings.TrimSpace(r.Type) == "" {
		return errors.New("project type is required")
	}
	
	// Validate project name format
	if !isValidProjectName(r.Name) {
		return errors.New("project name must contain only alphanumeric characters, hyphens, and underscores")
	}
	
	// Validate module path format (basic check)
	if !isValidModulePath(r.ModulePath) {
		return errors.New("module path must be a valid Go module path")
	}
	
	// Validate project type
	validTypes := []string{
		"web-api", "cli", "library", "lambda", "lambda-proxy",
		"event-driven", "microservice", "monolith", "workspace",
	}
	if !contains(validTypes, r.Type) {
		return errors.New("invalid project type")
	}
	
	// Validate architecture if provided
	if r.Architecture != "" {
		validArchitectures := []string{"standard", "clean", "ddd", "hexagonal", "event-driven"}
		if !contains(validArchitectures, r.Architecture) {
			return errors.New("invalid architecture type")
		}
	}
	
	// Validate complexity if provided
	if r.Complexity != "" {
		validComplexities := []string{"simple", "standard", "advanced", "expert"}
		if !contains(validComplexities, r.Complexity) {
			return errors.New("invalid complexity level")
		}
	}
	
	// Validate logger if provided
	if r.Logger != "" {
		validLoggers := []string{"slog", "zap", "logrus", "zerolog"}
		if !contains(validLoggers, r.Logger) {
			return errors.New("invalid logger type")
		}
	}
	
	return nil
}

// SetDefaults sets default values for optional fields
func (r *GenerationRequest) SetDefaults() {
	if r.GoVersion == "" {
		r.GoVersion = "1.21"
	}
	
	if r.Logger == "" {
		r.Logger = "slog"
	}
	
	if r.OutputFormat == "" {
		r.OutputFormat = "zip"
	}
	
	if r.Complexity == "" {
		r.Complexity = "standard"
	}
	
	// Set framework defaults based on project type
	if r.Framework == "" {
		switch r.Type {
		case "web-api":
			r.Framework = "gin"
		case "cli":
			r.Framework = "cobra"
		default:
			// Keep empty for other types
		}
	}
	
	// Set architecture defaults based on project type
	if r.Architecture == "" {
		r.Architecture = "standard"
	}
}

// isValidProjectName checks if the project name is valid
func isValidProjectName(name string) bool {
	if len(name) == 0 || len(name) > 50 {
		return false
	}
	
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}
	
	return true
}

// isValidModulePath checks if the module path is valid (basic validation)
func isValidModulePath(path string) bool {
	if len(path) == 0 || len(path) > 255 {
		return false
	}
	
	// Must contain at least one slash (domain/path format)
	if !strings.Contains(path, "/") {
		return false
	}
	
	// Basic format check
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return false
	}
	
	// Domain part should contain a dot or be localhost
	domain := parts[0]
	if domain != "localhost" && !strings.Contains(domain, ".") {
		return false
	}
	
	return true
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}