package ui

import (
	"fmt"
	"strings"

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

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Padding(1, 2)
)

// FangPrompter is a Fang-based UI prompter
type FangPrompter struct{}

// NewFangPrompter creates a new Fang-based prompter
func NewFangPrompter() *FangPrompter {
	return &FangPrompter{}
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

	result, err := RunTextInput("🚀 What's your project name?", help, suggestion)
	if err != nil {
		fmt.Println(errorStyle.Render("❌ Failed to get project name: " + err.Error()))
		return "", err
	}
	return result, nil
}

// promptModulePath prompts for Go module path using interactive UI
func (f *FangPrompter) promptModulePath(projectName string) (string, error) {
	defaultModule := fmt.Sprintf("github.com/username/%s", projectName)
	help := "Go module path for imports (e.g., github.com/username/project)"

	return RunTextInput("Module path:", help, defaultModule)
}

// promptProjectType prompts for project type selection using interactive UI
func (f *FangPrompter) promptProjectType() (string, error) {
	items := []SelectionItem{
		{title: "Web API", description: "REST API or web service", value: "web-api"},
		{title: "CLI Application", description: "Command-line tool", value: "cli"},
		{title: "Library", description: "Reusable Go package", value: "library"},
		{title: "AWS Lambda", description: "Serverless function", value: "lambda"},
	}

	return RunSelection("What type of project?", items)
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

	return RunSelection("Which framework?", items)
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

	return RunSelection("Which logger?", items)
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

	choice, err := RunSelection("Architecture pattern?", items)
	if err != nil {
		// Default to standard architecture
		config.Architecture = "standard"
		return nil
	}

	config.Architecture = choice
	return nil
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

