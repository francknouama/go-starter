package ui

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
)

// Styles for the Fang UI
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12")).
			Padding(1, 2).
			MarginBottom(1)

	questionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10"))

	optionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("12")).
			Padding(0, 1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			Italic(true).
			MarginTop(1)

)

// FangPrompter is a Fang-based UI prompter with Survey fallback
type FangPrompter struct {
	useEnhancedUI bool
	surveyAdapter SurveyAdapter
}

// SurveyAdapter defines an interface for survey operations
type SurveyAdapter interface {
	AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
}

// RealSurveyAdapter implements the SurveyAdapter interface
type RealSurveyAdapter struct{}

// AskOne calls the real survey.AskOne function
func (r *RealSurveyAdapter) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}

// NewFangPrompter creates a new Fang-based prompter with enhanced UI
func NewFangPrompter() *FangPrompter {
	return &FangPrompter{
		useEnhancedUI: true,
		surveyAdapter: &RealSurveyAdapter{},
	}
}

// NewFangPrompterWithSurvey creates a new Fang-based prompter using Survey fallback
func NewFangPrompterWithSurvey() *FangPrompter {
	return &FangPrompter{
		useEnhancedUI: false,
		surveyAdapter: &RealSurveyAdapter{},
	}
}

// GetProjectConfig prompts the user for project configuration using Fang UI
func (f *FangPrompter) GetProjectConfig(initial types.ProjectConfig, advanced bool) (types.ProjectConfig, error) {
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
		name, err := f.promptProjectName()
		if err != nil {
			return config, err
		}
		config.Name = name
	}

	// Module path
	if config.Module == "" {
		module, err := f.promptModulePath(config.Name)
		if err != nil {
			return config, err
		}
		config.Module = module
	}

	// Project type
	if config.Type == "" {
		projectType, err := f.promptProjectType()
		if err != nil {
			return config, err
		}
		config.Type = projectType
	}

	// Framework selection based on project type
	if config.Framework == "" {
		framework, err := f.promptFramework(config.Type)
		if err != nil {
			return config, err
		}
		config.Framework = framework
	}

	// Go version selection  
	if config.GoVersion == "" {
		goVersion, err := f.promptGoVersion()
		if err != nil {
			return config, err
		}
		config.GoVersion = goVersion
	}

	// Logger selection
	if config.Logger == "" {
		logger, err := f.promptLogger(config.Type)
		if err != nil {
			return config, err
		}
		config.Logger = logger
	}

	// Advanced options
	if advanced {
		if err := f.promptAdvancedOptions(&config); err != nil {
			return config, err
		}
	} else if f.isInteractiveMode(initial) {
		if err := f.promptBasicOptions(&config); err != nil {
			return config, err
		}
	}

	return config, nil
}

// isInteractiveMode determines if we need to prompt for additional options
func (f *FangPrompter) isInteractiveMode(initial types.ProjectConfig) bool {
	return initial.Name == "" || initial.Module == "" || initial.Type == ""
}

// promptProjectName prompts for project name with suggestions using interactive UI
func (f *FangPrompter) promptProjectName() (string, error) {
	suggestion := utils.GenerateRandomProjectName()
	alternatives := utils.GenerateMultipleNames(3)

	help := fmt.Sprintf("Press Enter for: %s\nOther suggestions: %s", 
		suggestion, strings.Join(alternatives, ", "))

	if f.useEnhancedUI {
		result, err := RunTextInput("🚀 What's your project name?", help, suggestion)
		if err != nil {
			// Fallback to survey on error
			return f.promptProjectNameSurvey()
		}
		return result, nil
	}
	return f.promptProjectNameSurvey()
}

// promptModulePath prompts for Go module path using interactive UI
func (f *FangPrompter) promptModulePath(projectName string) (string, error) {
	defaultModule := fmt.Sprintf("github.com/username/%s", projectName)
	help := "Go module path for imports (e.g., github.com/username/project)"

	if f.useEnhancedUI {
		result, err := RunTextInput("Module path:", help, defaultModule)
		if err != nil {
			// Fallback to survey on error
			return f.promptModulePathSurvey(projectName)
		}
		return result, nil
	}
	return f.promptModulePathSurvey(projectName)
}

// promptProjectType prompts for project type selection using interactive UI
func (f *FangPrompter) promptProjectType() (string, error) {
	items := []SelectionItem{
		{title: "Web API", description: "REST API or web service", value: "web-api"},
		{title: "CLI Application", description: "Command-line tool", value: "cli"},
		{title: "Library", description: "Reusable Go package", value: "library"},
		{title: "AWS Lambda", description: "Serverless function", value: "lambda"},
	}

	if f.useEnhancedUI {
		result, err := RunSelection("What type of project?", items)
		if err != nil {
			// Fallback to survey on error
			return f.promptProjectTypeSurvey()
		}
		return result, nil
	}
	return f.promptProjectTypeSurvey()
}

// promptFramework prompts for framework selection using interactive UI
func (f *FangPrompter) promptFramework(projectType string) (string, error) {
	var items []SelectionItem

	switch projectType {
	case "web-api":
		items = []SelectionItem{
			{title: "Gin", description: "Fast HTTP web framework (recommended)", value: "gin"},
			{title: "Echo", description: "High performance, minimalist web framework", value: "echo"},
			{title: "Fiber", description: "Express inspired web framework", value: "fiber"},
			{title: "Chi", description: "Lightweight, idiomatic router", value: "chi"},
			{title: "Standard library", description: "Built-in net/http package", value: "standard"},
		}
	case "cli":
		items = []SelectionItem{
			{title: "Cobra", description: "Powerful CLI framework (recommended)", value: "cobra"},
			{title: "Standard library", description: "Built-in flag package", value: "standard"},
		}
	default:
		// No framework selection needed
		return "", nil
	}

	if f.useEnhancedUI {
		result, err := RunSelection("Which framework?", items)
		if err != nil {
			// Fallback to survey on error
			return f.promptFrameworkSurvey(projectType)
		}
		return result, nil
	}
	return f.promptFrameworkSurvey(projectType)
}

// promptLogger prompts for logger selection using interactive UI
func (f *FangPrompter) promptLogger(projectType string) (string, error) {
	// Skip logger selection for library projects
	if projectType == "library" {
		return "slog", nil
	}

	items := []SelectionItem{
		{title: "slog", description: "Go built-in structured logging (recommended)", value: "slog"},
		{title: "zap", description: "High-performance, zero-allocation logging", value: "zap"},
		{title: "logrus", description: "Feature-rich, popular logging library", value: "logrus"},
		{title: "zerolog", description: "Zero allocation, chainable API logging", value: "zerolog"},
	}

	if f.useEnhancedUI {
		result, err := RunSelection("Which logger?", items)
		if err != nil {
			// Fallback to survey on error
			return f.promptLoggerSurvey(projectType)
		}
		return result, nil
	}
	return f.promptLoggerSurvey(projectType)
}

// promptBasicOptions prompts for basic configuration options
func (f *FangPrompter) promptBasicOptions(config *types.ProjectConfig) error {
	if config.Type == "web-api" {
		return f.promptDatabaseSupport(config)
	}
	return nil
}

// promptAdvancedOptions prompts for advanced configuration options
func (f *FangPrompter) promptAdvancedOptions(config *types.ProjectConfig) error {
	if config.Type == "web-api" {
		if err := f.promptArchitecture(config); err != nil {
			return err
		}
		if err := f.promptDatabaseSupport(config); err != nil {
			return err
		}
		if err := f.promptAuthentication(config); err != nil {
			return err
		}
	}
	return nil
}

// promptArchitecture prompts for architecture pattern using interactive UI
func (f *FangPrompter) promptArchitecture(config *types.ProjectConfig) error {
	items := []SelectionItem{
		{title: "Standard", description: "Simple, straightforward structure", value: "standard"},
		{title: "Clean Architecture", description: "Uncle Bob's principles", value: "clean"},
		{title: "Domain-Driven Design", description: "Business-focused design", value: "ddd"},
		{title: "Hexagonal", description: "Ports and adapters pattern", value: "hexagonal"},
	}

	if f.useEnhancedUI {
		choice, err := RunSelection("Architecture pattern?", items)
		if err != nil {
			// Fallback to survey on error
			return f.promptArchitectureSurvey(config)
		}
		config.Architecture = choice
		return nil
	}
	return f.promptArchitectureSurvey(config)
}

// promptDatabaseSupport prompts for database configuration
func (f *FangPrompter) promptDatabaseSupport(config *types.ProjectConfig) error {
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	fmt.Print(questionStyle.Render("Add database support? (y/N): "))
	fmt.Print(optionStyle.Render(" "))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// Default to no if input fails
		response = "n"
	}

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		// Simple database selection for now
		config.Features.Database.Drivers = []string{"postgres"}
		config.Features.Database.ORM = "gorm"
		fmt.Println(selectedStyle.Render("✓ Database support enabled with PostgreSQL"))
	}

	return nil
}

// promptAuthentication prompts for authentication configuration
func (f *FangPrompter) promptAuthentication(config *types.ProjectConfig) error {
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	fmt.Print(questionStyle.Render("Add authentication? (y/N): "))
	fmt.Print(optionStyle.Render(" "))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		// Default to no if input fails
		response = "n"
	}

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		config.Features.Authentication.Type = "jwt"
		fmt.Println(selectedStyle.Render("✓ JWT authentication enabled"))
	}

	return nil
}

// Option types for better organization
type ProjectTypeOption struct {
	Value       string
	Display     string
	Description string
}

type FrameworkOption struct {
	Value       string
	Display     string
	Description string
}

type LoggerOption struct {
	Value       string
	Display     string
	Description string
}

type ArchitectureOption struct {
	Value       string
	Display     string
	Description string
}

// promptProjectNameSurvey provides Survey fallback for project name
func (f *FangPrompter) promptProjectNameSurvey() (string, error) {
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
	err := f.surveyAdapter.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	return result, err
}

// promptModulePathSurvey provides Survey fallback for module path
func (f *FangPrompter) promptModulePathSurvey(projectName string) (string, error) {
	defaultModule := fmt.Sprintf("github.com/username/%s", projectName)
	prompt := &survey.Input{
		Message: "Module path:",
		Default: defaultModule,
		Help:    "Go module path for imports (e.g., github.com/username/project)",
	}
	var result string
	err := f.surveyAdapter.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	return result, err
}

// promptProjectTypeSurvey provides Survey fallback for project type
func (f *FangPrompter) promptProjectTypeSurvey() (string, error) {
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
	if err := f.surveyAdapter.AskOne(prompt, &selection); err != nil {
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

// promptFrameworkSurvey provides Survey fallback for framework selection
func (f *FangPrompter) promptFrameworkSurvey(projectType string) (string, error) {
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
	if err := f.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract framework name (remove description)
	return strings.ToLower(strings.Split(selection, " ")[0]), nil
}

// promptLoggerSurvey provides Survey fallback for logger selection
func (f *FangPrompter) promptLoggerSurvey(projectType string) (string, error) {
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
	if err := f.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract logger name (first word before the dash)
	return strings.Split(selection, " ")[0], nil
}

// promptGoVersion prompts for Go version selection
func (f *FangPrompter) promptGoVersion() (string, error) {
	// Go version options
	items := []SelectionItem{
		{title: "1.21", description: "Stable LTS release (recommended)", value: "1.21"},
		{title: "1.22", description: "Latest stable release", value: "1.22"},
		{title: "1.23", description: "Latest release", value: "1.23"},
		{title: "1.24", description: "Development release", value: "1.24"},
	}

	if f.useEnhancedUI {
		result, err := RunSelection("Which Go version?", items)
		if err != nil {
			// Fallback to survey on error
			return f.promptGoVersionSurvey()
		}
		return result, nil
	}
	return f.promptGoVersionSurvey()
}

// promptGoVersionSurvey provides Survey fallback for Go version selection
func (f *FangPrompter) promptGoVersionSurvey() (string, error) {
	options := []string{
		"1.21 - Stable LTS release (recommended)",
		"1.22 - Latest stable release", 
		"1.23 - Latest release",
		"1.24 - Development release",
	}

	prompt := &survey.Select{
		Message: "Which Go version?",
		Options: options,
		Default: "1.21 - Stable LTS release (recommended)",
		Help:    "Choose the Go version for your project. 1.21 is recommended for most projects.",
	}

	var selection string
	if err := f.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	// Extract version number (first part before the dash)
	return strings.Split(selection, " ")[0], nil
}

// promptArchitectureSurvey provides Survey fallback for architecture selection
func (f *FangPrompter) promptArchitectureSurvey(config *types.ProjectConfig) error {
	options := []string{
		"Standard - Simple structure",
		"Clean Architecture - Uncle Bob's principles",
		"Domain-Driven Design - Business-focused",
		"Hexagonal - Ports and adapters",
	}

	prompt := &survey.Select{
		Message: "Architecture pattern?",
		Options: options,
		Default: "Standard - Simple structure",
		Help:    "Choose the architectural pattern for your project",
	}

	var selection string
	if err := f.surveyAdapter.AskOne(prompt, &selection); err != nil {
		// Default to standard architecture on error
		config.Architecture = "standard"
		return nil
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

