package survey

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
)


// SurveyPrompter implements the Prompter interface using AlecAivazis/survey
type SurveyPrompter struct {
	surveyAdapter interfaces.SurveyAdapter
}

// NewPrompter creates a new SurveyPrompter
func NewPrompter() interfaces.Prompter {
	return &SurveyPrompter{
		surveyAdapter: &interfaces.RealSurveyAdapter{},
	}
}

// NewWithAdapter creates a new SurveyPrompter with a custom survey adapter
func NewWithAdapter(adapter interfaces.SurveyAdapter) interfaces.Prompter {
	return &SurveyPrompter{
		surveyAdapter: adapter,
	}
}

// GetProjectConfig prompts the user for project configuration using Survey
func (p *SurveyPrompter) GetProjectConfig(initial types.ProjectConfig, advanced bool) (types.ProjectConfig, error) {
	config := initial

	// Set defaults
	if config.GoVersion == "" {
		config.GoVersion = utils.GetOptimalGoVersion()
	}
	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}
	if config.Features == nil {
		config.Features = &types.Features{
			Database:       types.DatabaseConfig{},
			Authentication: types.AuthConfig{},
			Deployment:     types.DeployConfig{},
			Testing:        types.TestConfig{},
			Monitoring:     types.MonitorConfig{},
			Logging:        types.LoggingConfig{},
		}
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

	// Go version selection
	if config.GoVersion == "" {
		if err := p.promptGoVersionSelection(&config); err != nil {
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
func (p *SurveyPrompter) isInteractiveMode(initial types.ProjectConfig) bool {
	// If most fields are already provided, we're in non-interactive mode
	return initial.Name == "" || initial.Module == "" || initial.Type == ""
}

func (p *SurveyPrompter) promptProjectName(config *types.ProjectConfig) error {
	name, err := p.promptProjectNameSurvey()
	if err != nil {
		return err
	}
	config.Name = name
	return nil
}

func (p *SurveyPrompter) promptModulePath(config *types.ProjectConfig) error {
	module, err := p.promptModulePathSurvey(config.Name)
	if err != nil {
		return err
	}
	config.Module = module
	return nil
}

func (p *SurveyPrompter) promptProjectType(config *types.ProjectConfig) error {
	projectType, err := p.promptProjectTypeSurvey()
	if err != nil {
		return err
	}
	config.Type = projectType
	return nil
}

func (p *SurveyPrompter) promptFramework(config *types.ProjectConfig) error {
	framework, err := p.promptFrameworkSurvey(config.Type)
	if err != nil {
		return err
	}
	config.Framework = framework
	return nil
}

func (p *SurveyPrompter) promptBasicOptions(config *types.ProjectConfig) error {
	// Basic options for non-advanced mode
	if config.Type == "web-api" {
		return p.promptDatabaseSupport(config)
	}
	return nil
}

func (p *SurveyPrompter) promptAdvancedOptions(config *types.ProjectConfig) error {
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

func (p *SurveyPrompter) promptArchitecture(config *types.ProjectConfig) error {
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
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
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

func (p *SurveyPrompter) promptDatabaseSupport(config *types.ProjectConfig) error {
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

	if err := p.surveyAdapter.AskOne(prompt, &addDB); err != nil {
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
		if err := p.surveyAdapter.AskOne(dbPrompt, &selectedDBs); err != nil {
			return err
		}

		// Convert to lowercase for consistency
		var drivers []string
		for _, db := range selectedDBs {
			drivers = append(drivers, strings.ToLower(db))
		}

		config.Features.Database.Drivers = drivers

		// Prompt for ORM selection
		if err := p.promptORM(config); err != nil {
			return err
		}

		// For backward compatibility, set Driver to the first selected database
		if len(drivers) > 0 {
			config.Features.Database.Driver = drivers[0] //nolint:staticcheck // kept for backward compatibility
		}
	}

	return nil
}

func (p *SurveyPrompter) promptORM(config *types.ProjectConfig) error {
	ormPrompt := &survey.Select{
		Message: "Which ORM/database abstraction do you prefer?",
		Options: []string{
			"gorm - Feature-rich ORM with associations and migrations (recommended) ✅",
			"raw - Raw database/sql package with manual queries ✅",
			"sqlx - Lightweight extensions on database/sql 🔄 Coming Soon",
			"sqlc - Generate type-safe code from SQL 🔄 Coming Soon",
			"ent - Simple, yet feature-complete entity framework 🔄 Coming Soon",
			"xorm - Alternative full-featured ORM 🔄 Coming Soon",
		},
		Default: "gorm - Feature-rich ORM with associations and migrations (recommended) ✅",
		Help:    "✅ = Currently supported | 🔄 = Coming soon in future releases. GORM provides rich ORM features, while raw gives full control over SQL.",
	}

	var selection string
	if err := p.surveyAdapter.AskOne(ormPrompt, &selection); err != nil {
		return err
	}

	// Map selection to ORM
	ormMap := map[string]string{
		"gorm - Feature-rich ORM with associations and migrations (recommended) ✅": "gorm",
		"raw - Raw database/sql package with manual queries ✅":                     "raw",
		"sqlx - Lightweight extensions on database/sql 🔄 Coming Soon":              "sqlx",
		"sqlc - Generate type-safe code from SQL 🔄 Coming Soon":                    "sqlc",
		"ent - Simple, yet feature-complete entity framework 🔄 Coming Soon":        "ent",
		"xorm - Alternative full-featured ORM 🔄 Coming Soon":                       "xorm",
	}

	selectedORM := ormMap[selection]

	// Check if the selected ORM is implemented
	if selectedORM != "gorm" && selectedORM != "raw" {
		message := fmt.Sprintf("ORM '%s' is not yet implemented. Currently supported: gorm, raw", selectedORM)
		return types.NewValidationError(message, nil)
	}

	config.Features.Database.ORM = selectedORM
	return nil
}

func (p *SurveyPrompter) promptAuthentication(config *types.ProjectConfig) error {
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

	if err := p.surveyAdapter.AskOne(prompt, &addAuth); err != nil {
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
		if err := p.surveyAdapter.AskOne(authPrompt, &authType); err != nil {
			return err
		}

		config.Features.Authentication.Type = strings.ToLower(authType)
	}

	return nil
}

func (p *SurveyPrompter) promptAdvancedLogger(config *types.ProjectConfig) error {
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
	if err := p.surveyAdapter.AskOne(levelPrompt, &levelSelection); err != nil {
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
	if err := p.surveyAdapter.AskOne(formatPrompt, &formatSelection); err != nil {
		return err
	}
	config.Features.Logging.Format = strings.Split(formatSelection, " ")[0]

	// Structured logging (always enabled for consistency)
	config.Features.Logging.Structured = true

	return nil
}

func (p *SurveyPrompter) promptGoVersionSelection(config *types.ProjectConfig) error {
	goVersion, err := p.promptGoVersionSurvey()
	if err != nil {
		return err
	}
	config.GoVersion = goVersion
	return nil
}

func (p *SurveyPrompter) promptLogger(config *types.ProjectConfig) error {
	logger, err := p.promptLoggerSurvey(config.Type)
	if err != nil {
		return err
	}
	config.Logger = logger
	return nil
}

// promptProjectNameSurvey prompts for project name with suggestions
func (p *SurveyPrompter) promptProjectNameSurvey() (string, error) {
	suggestion := utils.GenerateRandomProjectName()
	alternatives := utils.GenerateMultipleNames(3)

	helpText := fmt.Sprintf("This will be used as the directory name and default module path.\n"+
		"Press Enter to use: %s\n"+
		"Other suggestions: %s",
		suggestion,
		strings.Join(alternatives, ", "))

	prompt := &survey.Input{
		Message: "What's your project name?",
		Default: suggestion,
		Help:    helpText,
	}
	var result string
	err := p.surveyAdapter.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	return result, err
}

func (p *SurveyPrompter) promptModulePathSurvey(projectName string) (string, error) {
	defaultModule := fmt.Sprintf("github.com/username/%s", projectName)
	prompt := &survey.Input{
		Message: "Module path:",
		Default: defaultModule,
		Help:    "Go module path for imports (e.g., github.com/username/project)",
	}
	var result string
	err := p.surveyAdapter.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	return result, err
}

func (p *SurveyPrompter) promptProjectTypeSurvey() (string, error) {
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
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Map display names to internal types
	typeMap := map[string]string{
		"Web API - REST API or web service":   "web-api",
		"CLI Application - Command-line tool": "cli",
		"Library - Reusable Go package":       "library",
		"AWS Lambda - Serverless function":    "lambda",
	}

	return typeMap[selection], nil
}

func (p *SurveyPrompter) promptFrameworkSurvey(projectType string) (string, error) {
	var options []string
	var message string

	switch projectType {
	case "web-api":
		message = "Which web framework?"
		options = []string{"Gin (recommended)", "Echo", "Fiber", "Chi", "Standard library"}
	case "cli":
		message = "Which CLI framework?"
		options = []string{"Cobra (recommended)", "Standard library"}
	default:
		// No framework selection needed for library or lambda
		return "", nil
	}

	prompt := &survey.Select{
		Message: message,
		Options: options,
		Help:    "Choose the framework for your project",
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract framework name (remove description)
	return strings.ToLower(strings.Split(selection, " ")[0]), nil
}

func (p *SurveyPrompter) promptLoggerSurvey(projectType string) (string, error) {
	// Skip logger selection for library projects
	if projectType == "library" {
		return "slog", nil
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
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract logger name (first word before the dash)
	return strings.Split(selection, " ")[0], nil
}

func (p *SurveyPrompter) promptGoVersionSurvey() (string, error) {
	options := []string{
		"1.21 - Stable LTS release (recommended)",
		"1.22 - Latest stable release",
		"1.23 - Latest release",
	}

	prompt := &survey.Select{
		Message: "Which Go version?",
		Options: options,
		Default: "1.21 - Stable LTS release (recommended)",
		Help:    "Choose the Go version for your project. 1.21 is recommended for most projects.",
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract version number (first part before the dash)
	return strings.Split(selection, " ")[0], nil
}