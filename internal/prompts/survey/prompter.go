package survey

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	
	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
)


// SurveyPrompter implements the Prompter interface using AlecAivazis/survey
type SurveyPrompter struct {
	surveyAdapter interfaces.SurveyAdapter
	registry      *templates.Registry
}

// NewPrompter creates a new SurveyPrompter
func NewPrompter() interfaces.Prompter {
	return &SurveyPrompter{
		surveyAdapter: &interfaces.RealSurveyAdapter{},
		registry:      templates.NewRegistry(),
	}
}

// NewWithAdapter creates a new SurveyPrompter with a custom survey adapter
func NewWithAdapter(adapter interfaces.SurveyAdapter) interfaces.Prompter {
	return &SurveyPrompter{
		surveyAdapter: adapter,
		registry:      templates.NewRegistry(),
	}
}

// GetProjectConfigWithDisclosure prompts the user for project configuration using disclosure mode
func (p *SurveyPrompter) GetProjectConfigWithDisclosure(initial types.ProjectConfig, mode interfaces.DisclosureMode, complexity interfaces.ComplexityLevel) (types.ProjectConfig, error) {
	// Convert disclosure mode to advanced boolean for compatibility
	advanced := mode == interfaces.DisclosureModeAdvanced
	return p.GetProjectConfig(initial, advanced)
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
		if err := p.promptProjectType(&config, advanced); err != nil {
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

func (p *SurveyPrompter) promptProjectType(config *types.ProjectConfig, advanced bool) error {
	projectType, blueprintID, err := p.promptProjectTypeSurvey(advanced)
	if err != nil {
		return err
	}
	config.Type = projectType
	
	// Store the specific blueprint ID if different from type
	if blueprintID != projectType {
		if config.Variables == nil {
			config.Variables = make(map[string]string)
		}
		config.Variables["blueprint_id"] = blueprintID
	}
	
	// If CLI type selected, prompt for complexity level (unless already handled)
	if projectType == "cli" && blueprintID == "cli" {
		if err := p.promptCLIComplexity(config); err != nil {
			return err
		}
	}
	
	// If web-api type selected, prompt for architecture (unless already handled)
	if projectType == "web-api" && blueprintID == "web-api" {
		if err := p.promptWebAPIArchitecture(config, advanced); err != nil {
			return err
		}
	}
	
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
		Default: "raw - Raw database/sql package with manual queries ✅",
		Help:    "✅ = Currently supported | 🔄 = Coming soon in future releases. GORM provides rich ORM features, while raw gives full control over SQL.",
	}

	var selection string
	if err := p.surveyAdapter.AskOne(ormPrompt, &selection); err != nil {
		return err
	}

	// Map selection to ORM
	ormMap := map[string]string{
		"gorm - Feature-rich ORM with associations and migrations (recommended) ✅": "gorm",
		"raw - Raw database/sql package with manual queries ✅":                     "",
		"sqlx - Lightweight extensions on database/sql 🔄 Coming Soon":              "sqlx",
		"sqlc - Generate type-safe code from SQL 🔄 Coming Soon":                    "sqlc",
		"ent - Simple, yet feature-complete entity framework 🔄 Coming Soon":        "ent",
		"xorm - Alternative full-featured ORM 🔄 Coming Soon":                       "xorm",
	}

	selectedORM := ormMap[selection]

	// Check if the selected ORM is implemented
	if selectedORM != "gorm" && selectedORM != "" {
		message := fmt.Sprintf("ORM '%s' is not yet implemented. Currently supported: gorm, raw (empty)", selectedORM)
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

func (p *SurveyPrompter) promptProjectTypeSurvey(advanced bool) (string, string, error) {
	// Get all available blueprints from registry
	allBlueprints := p.registry.List()
	if len(allBlueprints) == 0 {
		return "", "", fmt.Errorf("no blueprints available in registry")
	}

	// Categorize blueprints for display
	categories := p.categorizeBlueprints(allBlueprints, advanced)
	
	// Build options list
	var options []string
	var blueprintMap = make(map[string]BlueprintSelection)
	
	for _, category := range categories {
		for _, item := range category.Items {
			displayName := item.DisplayName
			if category.ShowCategory && len(categories) > 1 {
				displayName = fmt.Sprintf("%-12s │ %s", category.Name, item.DisplayName)
			}
			options = append(options, displayName)
			blueprintMap[displayName] = item
		}
		
		// Add separator between categories (except for last)
		if category.ShowSeparator {
			separator := strings.Repeat("─", 50)
			options = append(options, separator)
			blueprintMap[separator] = BlueprintSelection{} // dummy entry
		}
	}

	// Determine help text based on mode
	helpText := "Choose the type of Go project you want to create"
	if !advanced {
		helpText += "\n💡 Use --advanced to see all available blueprints including advanced architectures"
	}

	prompt := &survey.Select{
		Message: "What type of project?",
		Options: options,
		Help:    helpText,
		Filter:  func(filter string, value string, index int) bool {
			// Hide separators from filtering
			if strings.Contains(value, "─") {
				return false
			}
			return strings.Contains(strings.ToLower(value), strings.ToLower(filter))
		},
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", "", err
	}

	// Handle separator selection (should not happen with proper filtering)
	if strings.Contains(selection, "─") {
		return "", "", fmt.Errorf("invalid selection")
	}

	selected := blueprintMap[selection]
	return selected.Type, selected.BlueprintID, nil
}

// BlueprintSelection represents a selectable blueprint option
type BlueprintSelection struct {
	Type        string // The base type (web-api, cli, etc.)
	BlueprintID string // The specific blueprint ID (web-api-clean, cli-simple, etc.)
	DisplayName string // User-friendly display name
}

// BlueprintCategory represents a category of blueprints
type BlueprintCategory struct {
	Name          string
	Items         []BlueprintSelection
	ShowCategory  bool
	ShowSeparator bool
}

// categorizeBlueprints organizes blueprints into categories based on disclosure mode
func (p *SurveyPrompter) categorizeBlueprints(blueprints []types.Template, advanced bool) []BlueprintCategory {
	if advanced {
		return p.getAdvancedBlueprintCategories(blueprints)
	}
	return p.getBasicBlueprintCategories(blueprints)
}

// getBasicBlueprintCategories returns essential blueprints for basic mode
func (p *SurveyPrompter) getBasicBlueprintCategories(blueprints []types.Template) []BlueprintCategory {
	// Create a map for quick lookup
	blueprintMap := make(map[string]types.Template)
	for _, bp := range blueprints {
		blueprintMap[bp.ID] = bp
	}

	// Define the 6 essential project types for basic mode
	essentialItems := []BlueprintSelection{
		{Type: "web-api", BlueprintID: "web-api-standard", DisplayName: "🌐 Web API - REST API or web service"},
		{Type: "cli", BlueprintID: "cli", DisplayName: "⚡ CLI Application - Command-line tool"},
		{Type: "library", BlueprintID: "library-standard", DisplayName: "📦 Library - Reusable Go package"},
		{Type: "lambda", BlueprintID: "lambda-standard", DisplayName: "☁️  AWS Lambda - Serverless function"},
		{Type: "monolith", BlueprintID: "monolith", DisplayName: "🏢 Monolith - Traditional web application"},
		{Type: "microservice", BlueprintID: "microservice-standard", DisplayName: "🔗 Microservice - Distributed service"},
	}

	// Filter out items that don't exist in the registry
	var availableItems []BlueprintSelection
	for _, item := range essentialItems {
		if _, exists := blueprintMap[item.BlueprintID]; exists {
			availableItems = append(availableItems, item)
		}
	}

	return []BlueprintCategory{
		{
			Name:          "Essential",
			Items:         availableItems,
			ShowCategory:  false,
			ShowSeparator: false,
		},
	}
}

// getAdvancedBlueprintCategories returns all blueprints organized by category
func (p *SurveyPrompter) getAdvancedBlueprintCategories(blueprints []types.Template) []BlueprintCategory {
	// Organize blueprints by type
	typeGroups := make(map[string][]types.Template)
	for _, bp := range blueprints {
		typeGroups[bp.Type] = append(typeGroups[bp.Type], bp)
	}

	// Sort blueprints within each group by architecture/complexity
	for _, group := range typeGroups {
		sort.Slice(group, func(i, j int) bool {
			// Standard architecture first, then alphabetical
			if group[i].Architecture == "standard" && group[j].Architecture != "standard" {
				return true
			}
			if group[i].Architecture != "standard" && group[j].Architecture == "standard" {
				return false
			}
			return group[i].Architecture < group[j].Architecture
		})
	}

	var categories []BlueprintCategory

	// Web APIs category
	if webAPIs, exists := typeGroups["web-api"]; exists {
		var items []BlueprintSelection
		for _, bp := range webAPIs {
			archName := cases.Title(language.English).String(strings.ReplaceAll(bp.Architecture, "-", " "))
			if bp.Architecture == "standard" {
				archName = "Standard"
			}
			items = append(items, BlueprintSelection{
				Type:        "web-api",
				BlueprintID: bp.ID,
				DisplayName: fmt.Sprintf("🌐 %s - %s", archName, getArchitectureDescription(bp.Architecture)),
			})
		}
		categories = append(categories, BlueprintCategory{
			Name:          "Web APIs",
			Items:         items,
			ShowCategory:  true,
			ShowSeparator: true,
		})
	}

	// CLI Tools category
	if cliTools, exists := typeGroups["cli"]; exists {
		var items []BlueprintSelection
		for _, bp := range cliTools {
			complexity := cases.Title(language.English).String(bp.Architecture)
			if bp.Architecture == "simple" {
				complexity = "Simple"
			}
			items = append(items, BlueprintSelection{
				Type:        "cli",
				BlueprintID: bp.ID,
				DisplayName: fmt.Sprintf("⚡ %s CLI - %s", complexity, getComplexityDescription(bp.Architecture)),
			})
		}
		categories = append(categories, BlueprintCategory{
			Name:          "CLI Tools",
			Items:         items,
			ShowCategory:  true,
			ShowSeparator: true,
		})
	}

	// Serverless category
	var serverlessItems []BlueprintSelection
	if lambdas, exists := typeGroups["lambda"]; exists {
		for _, bp := range lambdas {
			serverlessItems = append(serverlessItems, BlueprintSelection{
				Type:        "lambda",
				BlueprintID: bp.ID,
				DisplayName: "☁️  AWS Lambda - Event-driven serverless",
			})
		}
	}
	if lambdaProxy, exists := typeGroups["lambda-proxy"]; exists {
		for _, bp := range lambdaProxy {
			serverlessItems = append(serverlessItems, BlueprintSelection{
				Type:        "lambda-proxy",
				BlueprintID: bp.ID,
				DisplayName: "🔗 Lambda API Proxy - API Gateway integration",
			})
		}
	}
	if len(serverlessItems) > 0 {
		categories = append(categories, BlueprintCategory{
			Name:          "Serverless",
			Items:         serverlessItems,
			ShowCategory:  true,
			ShowSeparator: true,
		})
	}

	// Architecture Patterns category
	var architectureItems []BlueprintSelection
	if eventDriven, exists := typeGroups["event-driven"]; exists {
		for _, bp := range eventDriven {
			architectureItems = append(architectureItems, BlueprintSelection{
				Type:        "event-driven",
				BlueprintID: bp.ID,
				DisplayName: "📡 Event-Driven - CQRS & Event Sourcing",
			})
		}
	}
	if microservices, exists := typeGroups["microservice"]; exists {
		for _, bp := range microservices {
			architectureItems = append(architectureItems, BlueprintSelection{
				Type:        "microservice",
				BlueprintID: bp.ID,
				DisplayName: "🔗 Microservice - Distributed service patterns",
			})
		}
	}
	if monoliths, exists := typeGroups["monolith"]; exists {
		for _, bp := range monoliths {
			architectureItems = append(architectureItems, BlueprintSelection{
				Type:        "monolith",
				BlueprintID: bp.ID,
				DisplayName: "🏢 Monolith - Traditional web application",
			})
		}
	}
	if len(architectureItems) > 0 {
		categories = append(categories, BlueprintCategory{
			Name:          "Architecture",
			Items:         architectureItems,
			ShowCategory:  true,
			ShowSeparator: true,
		})
	}

	// Infrastructure category
	var infraItems []BlueprintSelection
	if grpcGateway, exists := typeGroups["grpc-gateway"]; exists {
		for _, bp := range grpcGateway {
			infraItems = append(infraItems, BlueprintSelection{
				Type:        "grpc-gateway",
				BlueprintID: bp.ID,
				DisplayName: "⚙️  gRPC Gateway - gRPC with REST gateway",
			})
		}
	}
	if workspaces, exists := typeGroups["workspace"]; exists {
		for _, bp := range workspaces {
			infraItems = append(infraItems, BlueprintSelection{
				Type:        "workspace",
				BlueprintID: bp.ID,
				DisplayName: "📁 Go Workspace - Multi-module projects",
			})
		}
	}
	if len(infraItems) > 0 {
		categories = append(categories, BlueprintCategory{
			Name:          "Infrastructure",
			Items:         infraItems,
			ShowCategory:  true,
			ShowSeparator: true,
		})
	}

	// Packages category
	if libraries, exists := typeGroups["library"]; exists {
		var items []BlueprintSelection
		for _, bp := range libraries {
			items = append(items, BlueprintSelection{
				Type:        "library",
				BlueprintID: bp.ID,
				DisplayName: "📦 Library - Reusable Go package",
			})
		}
		categories = append(categories, BlueprintCategory{
			Name:          "Packages",
			Items:         items,
			ShowCategory:  true,
			ShowSeparator: false, // Last category
		})
	}

	return categories
}

// Helper functions for descriptions
func getArchitectureDescription(arch string) string {
	switch arch {
	case "standard":
		return "Simple layered structure"
	case "clean":
		return "Uncle Bob's Clean Architecture"
	case "ddd":
		return "Domain-Driven Design"
	case "hexagonal":
		return "Ports & Adapters pattern"
	default:
		return "Advanced architecture pattern"
	}
}

func getComplexityDescription(complexity string) string {
	switch complexity {
	case "simple":
		return "Quick scripts & utilities"
	case "standard":
		return "Production-ready CLI tools"
	default:
		return "Command-line application"
	}
}

// promptWebAPIArchitecture prompts the user to choose web API architecture when needed
func (p *SurveyPrompter) promptWebAPIArchitecture(config *types.ProjectConfig, advanced bool) error {
	// Get all web-api blueprints from registry
	webAPIBlueprints := p.registry.GetByType("web-api")
	if len(webAPIBlueprints) <= 1 {
		// Only one web-api blueprint, no need to prompt
		return nil
	}

	// Sort architectures: standard first, then alphabetical
	sort.Slice(webAPIBlueprints, func(i, j int) bool {
		if webAPIBlueprints[i].Architecture == "standard" && webAPIBlueprints[j].Architecture != "standard" {
			return true
		}
		if webAPIBlueprints[i].Architecture != "standard" && webAPIBlueprints[j].Architecture == "standard" {
			return false
		}
		return webAPIBlueprints[i].Architecture < webAPIBlueprints[j].Architecture
	})

	// Build architecture options
	var options []string
	var archMap = make(map[string]string)

	for _, bp := range webAPIBlueprints {
		archName := cases.Title(language.English).String(strings.ReplaceAll(bp.Architecture, "-", " "))
		if bp.Architecture == "standard" {
			archName = "Standard"
		}
		
		displayName := fmt.Sprintf("%s - %s", archName, getArchitectureDescription(bp.Architecture))
		options = append(options, displayName)
		archMap[displayName] = bp.ID
	}

	helpText := `Web API Architecture Guide:

• Standard: Simple layered structure, great for most APIs
• Clean Architecture: Separation of concerns, highly testable  
• Domain-Driven Design: Business logic focused, complex domains
• Hexagonal: Ports & adapters, maximum testability

💡 Tip: Start with Standard, upgrade when complexity grows`

	prompt := &survey.Select{
		Message: "Which Web API architecture?",
		Options: options,
		Help:    helpText,
		Default: options[0], // Standard is first
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Update config with selected blueprint ID
	blueprintID := archMap[selection]
	if config.Variables == nil {
		config.Variables = make(map[string]string)
	}
	config.Variables["blueprint_id"] = blueprintID

	return nil
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

// promptCLIComplexity prompts the user to choose CLI complexity level with clear guidance
func (p *SurveyPrompter) promptCLIComplexity(config *types.ProjectConfig) error {
	options := []string{
		"Simple - Quick scripts & utilities (8 files, minimal deps)",
		"Standard - Production CLIs (30 files, full features)",
	}

	prompt := &survey.Select{
		Message: "Choose CLI complexity level:",
		Options: options,
		Help: `CLI Complexity Guide:

• Simple CLI (Recommended for 80% of use cases):
  - Quick utilities and scripts
  - Learning Go CLI development  
  - Internal tools with minimal requirements
  - Prototyping command-line interfaces
  - Projects needing < 3 commands
  - 8 files, single dependency (cobra)

• Standard CLI (Enterprise-grade):
  - Production CLI tools for distribution
  - Multiple subcommands (5+)
  - Configuration file support
  - Complex business logic
  - Team collaboration with CI/CD
  - 30 files, multiple dependencies

💡 Tip: Start simple, migrate to standard when needed`,
		Default: options[0], // Default to Simple for most users
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return err
	}

	// Map selection to complexity and update config
	if strings.HasPrefix(selection, "Simple") {
		config.Variables["complexity"] = "simple"
		config.Variables["blueprint_hint"] = "cli-simple"
	} else {
		config.Variables["complexity"] = "standard"
		config.Variables["blueprint_hint"] = "cli-standard"
	}

	return nil
}