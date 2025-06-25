package prompts

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/internal/utils"
)

// Prompter handles interactive prompts for project configuration
type Prompter struct{}

// New creates a new Prompter instance
func New() *Prompter {
	return &Prompter{}
}

// GetProjectConfig prompts the user for project configuration
func (p *Prompter) GetProjectConfig(initial types.ProjectConfig, advanced bool) (types.ProjectConfig, error) {
	config := initial

	// Set defaults
	if config.GoVersion == "" {
		config.GoVersion = utils.GetOptimalGoVersion()
	}
	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}

	// Project name
	if config.Name == "" {
		if err := p.promptProjectName(&config); err != nil {
			return config, err
		}
	}

	// Module path
	if config.Module == "" {
		if err := p.promptModulePath(&config); err != nil {
			return config, err
		}
	}

	// Project type
	if config.Type == "" {
		if err := p.promptProjectType(&config); err != nil {
			return config, err
		}
	}

	// Framework selection based on project type
	if config.Framework == "" {
		if err := p.promptFramework(&config); err != nil {
			return config, err
		}
	}

	// Logger selection (always prompt for applicable project types)
	if config.Logger == "" {
		if err := p.promptLogger(&config); err != nil {
			return config, err
		}
	}

	// Advanced options (only if in interactive mode)
	if advanced {
		if err := p.promptAdvancedOptions(&config); err != nil {
			return config, err
		}
	} else if p.isInteractiveMode(initial) {
		if err := p.promptBasicOptions(&config); err != nil {
			return config, err
		}
	}

	return config, nil
}

// isInteractiveMode determines if we need to prompt for additional options
func (p *Prompter) isInteractiveMode(initial types.ProjectConfig) bool {
	// If most fields are already provided, we're in non-interactive mode
	return initial.Name == "" || initial.Module == "" || initial.Type == ""
}

func (p *Prompter) promptProjectName(config *types.ProjectConfig) error {
	prompt := &survey.Input{
		Message: "What's your project name?",
		Help:    "This will be used as the directory name and default module path",
	}
	return survey.AskOne(prompt, &config.Name, survey.WithValidator(survey.Required))
}

func (p *Prompter) promptModulePath(config *types.ProjectConfig) error {
	defaultModule := fmt.Sprintf("github.com/username/%s", config.Name)

	prompt := &survey.Input{
		Message: "Module path:",
		Default: defaultModule,
		Help:    "Go module path for imports (e.g., github.com/username/project)",
	}
	return survey.AskOne(prompt, &config.Module, survey.WithValidator(survey.Required))
}

func (p *Prompter) promptProjectType(config *types.ProjectConfig) error {
	// For Phase 0, we only show available types
	// In Phase 1+, this will be dynamic based on available templates
	options := []string{
		"Web API - REST API or web service",
		"CLI Application - Command-line tool",
		"Library - Reusable Go package",
		"AWS Lambda - Serverless function",
	}

	prompt := &survey.Select{
		Message: "What type of project?",
		Options: options,
		Help:    "Choose the type of Go project you want to create",
	}

	var selection string
	if err := survey.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Map display names to internal types
	typeMap := map[string]string{
		"Web API - REST API or web service":   "web-api",
		"CLI Application - Command-line tool": "cli",
		"Library - Reusable Go package":       "library",
		"AWS Lambda - Serverless function":    "lambda",
	}

	config.Type = typeMap[selection]
	return nil
}

func (p *Prompter) promptFramework(config *types.ProjectConfig) error {
	var options []string
	var message string

	switch config.Type {
	case "web-api":
		message = "Which web framework?"
		options = []string{"Gin (recommended)", "Echo", "Fiber", "Chi", "Standard library"}
	case "cli":
		message = "Which CLI framework?"
		options = []string{"Cobra (recommended)", "Standard library"}
	default:
		// No framework selection needed for library or lambda
		return nil
	}

	prompt := &survey.Select{
		Message: message,
		Options: options,
		Help:    "Choose the framework for your project",
	}

	var selection string
	if err := survey.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Extract framework name (remove description)
	config.Framework = strings.ToLower(strings.Split(selection, " ")[0])
	return nil
}

func (p *Prompter) promptBasicOptions(config *types.ProjectConfig) error {
	// Basic options for non-advanced mode
	if config.Type == "web-api" {
		return p.promptDatabaseSupport(config)
	}
	return nil
}

func (p *Prompter) promptAdvancedOptions(config *types.ProjectConfig) error {
	// Advanced configuration options
	if config.Type == "web-api" {
		if err := p.promptArchitecture(config); err != nil {
			return err
		}
		if err := p.promptDatabaseSupport(config); err != nil {
			return err
		}
		if err := p.promptAuthentication(config); err != nil {
			return err
		}
	}

	// Advanced logger configuration (for all project types except library)
	if config.Type != "library" {
		if err := p.promptAdvancedLogger(config); err != nil {
			return err
		}
	}

	return nil
}

func (p *Prompter) promptArchitecture(config *types.ProjectConfig) error {
	prompt := &survey.Select{
		Message: "Architecture pattern?",
		Options: []string{
			"Standard - Simple structure",
			"Clean Architecture - Uncle Bob's principles",
			"Domain-Driven Design - Business-focused",
			"Hexagonal - Ports and adapters",
		},
		Default: "Standard - Simple structure",
		Help:    "Choose the architectural pattern for your project",
	}

	var selection string
	if err := survey.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Map selection to architecture
	archMap := map[string]string{
		"Standard - Simple structure":                 "standard",
		"Clean Architecture - Uncle Bob's principles": "clean",
		"Domain-Driven Design - Business-focused":     "ddd",
		"Hexagonal - Ports and adapters":              "hexagonal",
	}

	config.Architecture = archMap[selection]
	return nil
}

func (p *Prompter) promptDatabaseSupport(config *types.ProjectConfig) error {
	// Initialize Features if nil
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	addDB := false
	prompt := &survey.Confirm{
		Message: "Add database support?",
		Default: true,
		Help:    "Include database configuration and basic setup",
	}

	if err := survey.AskOne(prompt, &addDB); err != nil {
		return err
	}

	if addDB {
		// Use MultiSelect for multiple database selection
		dbPrompt := &survey.MultiSelect{
			Message: "Which databases do you want to use? (Space to select, Enter to confirm)",
			Options: []string{"PostgreSQL", "MySQL", "MongoDB", "SQLite", "Redis"},
			Default: []string{"PostgreSQL"},
			Help:    "Select one or more databases for your project. PostgreSQL for main data, Redis for caching, etc.",
		}

		var selectedDBs []string
		if err := survey.AskOne(dbPrompt, &selectedDBs); err != nil {
			return err
		}

		// Convert to lowercase for consistency
		var drivers []string
		for _, db := range selectedDBs {
			drivers = append(drivers, strings.ToLower(db))
		}

		config.Features.Database.Drivers = drivers
		config.Features.Database.ORM = "gorm" // Default ORM
		
		// For backward compatibility, set Driver to the first selected database
		if len(drivers) > 0 {
			config.Features.Database.Driver = drivers[0]
		}
	}

	return nil
}

func (p *Prompter) promptAuthentication(config *types.ProjectConfig) error {
	// Initialize Features if nil
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	addAuth := false
	prompt := &survey.Confirm{
		Message: "Add authentication?",
		Default: false,
		Help:    "Include authentication setup (JWT, OAuth, etc.)",
	}

	if err := survey.AskOne(prompt, &addAuth); err != nil {
		return err
	}

	if addAuth {
		authPrompt := &survey.Select{
			Message: "Authentication type?",
			Options: []string{"JWT", "OAuth2", "Session-based", "API Key"},
			Default: "JWT",
			Help:    "Choose the authentication method",
		}

		var authType string
		if err := survey.AskOne(authPrompt, &authType); err != nil {
			return err
		}

		config.Features.Authentication.Type = strings.ToLower(authType)
	}

	return nil
}

func (p *Prompter) promptAdvancedLogger(config *types.ProjectConfig) error {
	// Initialize Features if nil
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	// Set logger defaults if not already set
	if config.Features.Logging.Type == "" {
		config.Features.Logging.Type = config.Logger
	}

	// Log level configuration
	levelPrompt := &survey.Select{
		Message: "Log level?",
		Options: []string{
			"debug - Detailed debugging information",
			"info - General application flow (recommended)",
			"warn - Warning messages and potential issues",
			"error - Error conditions only",
		},
		Default: "info - General application flow (recommended)",
		Help:    "Choose the default log level for your application",
	}

	var levelSelection string
	if err := survey.AskOne(levelPrompt, &levelSelection); err != nil {
		return err
	}
	config.Features.Logging.Level = strings.Split(levelSelection, " ")[0]

	// Log format configuration
	formatPrompt := &survey.Select{
		Message: "Log format?",
		Options: []string{
			"json - Structured JSON format (recommended)",
			"text - Human-readable text format",
			"console - Colored console output",
		},
		Default: "json - Structured JSON format (recommended)",
		Help:    "Choose the log output format. JSON is recommended for production.",
	}

	var formatSelection string
	if err := survey.AskOne(formatPrompt, &formatSelection); err != nil {
		return err
	}
	config.Features.Logging.Format = strings.Split(formatSelection, " ")[0]

	// Structured logging (always enabled for consistency)
	config.Features.Logging.Structured = true

	return nil
}

func (p *Prompter) promptLogger(config *types.ProjectConfig) error {
	// Skip logger selection for library projects (they typically don't need logging setup)
	if config.Type == "library" {
		config.Logger = "slog" // Default for libraries
		return nil
	}

	options := []string{
		"slog - Go built-in structured logging (recommended)",
		"zap - High-performance, zero-allocation logging",
		"logrus - Feature-rich, popular logging library",
		"zerolog - Zero allocation, chainable API logging",
	}

	prompt := &survey.Select{
		Message: "Which logger?",
		Options: options,
		Default: "slog - Go built-in structured logging (recommended)",
		Help:    "Choose the logging library for your project. slog is built into Go 1.21+ and recommended for most projects.",
	}

	var selection string
	if err := survey.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Extract logger name (first word before the dash)
	config.Logger = strings.Split(selection, " ")[0]

	return nil
}
