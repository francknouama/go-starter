package bubbletea

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
)


// Styles for the enhanced UI
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

	listStyle = lipgloss.NewStyle().
		Margin(1, 2)

	titleListStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		MarginBottom(1)
)

// SelectionModel represents an interactive selection model
type SelectionModel struct {
	list     list.Model
	choice   string
	quitting bool
	title    string
}

func NewSelectionModel(title string, items []interfaces.SelectionItem) SelectionModel {
	listItems := make([]list.Item, len(items))
	for i, item := range items {
		listItems[i] = item
	}

	l := list.New(listItems, list.NewDefaultDelegate(), 80, 14)
	l.Title = title
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	l.Styles.Title = titleListStyle

	return SelectionModel{
		list:  l,
		title: title,
	}
}

func (m SelectionModel) Init() tea.Cmd {
	return nil
}

func (m SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(interfaces.SelectionItem)
			if ok {
				m.choice = i.Value()
			}
			m.quitting = true
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m SelectionModel) View() string {
	if m.quitting {
		return ""
	}
	return listStyle.Render(m.list.View())
}

// TextInputModel for text input prompts
type TextInputModel struct {
	textInput textinput.Model
	title     string
	help      string
	value     string
	quitting  bool
}

func NewTextInputModel(title, help, defaultValue string) TextInputModel {
	ti := textinput.New()
	ti.Placeholder = defaultValue
	ti.Focus()
	ti.CharLimit = 256
	ti.Width = 60

	return TextInputModel{
		textInput: ti,
		title:     title,
		help:      help,
	}
}

func (m TextInputModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m TextInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.value = m.textInput.Value()
			if m.value == "" {
				m.value = m.textInput.Placeholder
			}
			m.quitting = true
			return m, tea.Quit
		case tea.KeyCtrlC, tea.KeyEsc:
			m.quitting = true
			return m, tea.Quit
		}
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m TextInputModel) View() string {
	if m.quitting {
		return ""
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n",
		titleStyle.Render(m.title),
		m.textInput.View(),
		helpStyle.Render(m.help),
	)
}

// BubbleTeaPrompter implements the Prompter interface using Bubble Tea UI
type BubbleTeaPrompter struct{}


// NewPrompter creates a new BubbleTeaPrompter
func NewPrompter() interfaces.Prompter {
	return &BubbleTeaPrompter{}
}

// GetProjectConfigWithDisclosure prompts the user for project configuration using disclosure mode
func (p *BubbleTeaPrompter) GetProjectConfigWithDisclosure(initial types.ProjectConfig, mode interfaces.DisclosureMode, complexity interfaces.ComplexityLevel) (types.ProjectConfig, error) {
	// Convert disclosure mode to advanced boolean for compatibility
	advanced := mode == interfaces.DisclosureModeAdvanced
	return p.GetProjectConfig(initial, advanced)
}

// GetProjectConfig prompts the user for project configuration using Bubble Tea UI
func (p *BubbleTeaPrompter) GetProjectConfig(initial types.ProjectConfig, advanced bool) (types.ProjectConfig, error) {
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
		name, err := p.promptProjectName()
		if err != nil {
			return config, err
		}
		config.Name = name
	}

	// Module path
	if config.Module == "" {
		module, err := p.promptModulePath(config.Name)
		if err != nil {
			return config, err
		}
		config.Module = module
	}

	// Project type
	if config.Type == "" {
		projectType, err := p.promptProjectType()
		if err != nil {
			return config, err
		}
		config.Type = projectType
	}

	// Framework selection based on project type
	if config.Framework == "" {
		framework, err := p.promptFramework(config.Type)
		if err != nil {
			return config, err
		}
		config.Framework = framework
	}

	// Go version selection
	if config.GoVersion == "" {
		goVersion, err := p.promptGoVersion()
		if err != nil {
			return config, err
		}
		config.GoVersion = goVersion
	}

	// Logger selection
	if config.Logger == "" {
		logger, err := p.promptLogger(config.Type)
		if err != nil {
			return config, err
		}
		config.Logger = logger
	}

	// Advanced options
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

// RunSelection runs an interactive selection and returns the chosen value
func (p *BubbleTeaPrompter) RunSelection(title string, items []interfaces.SelectionItem) (string, error) {
	model := NewSelectionModel(title, items)

	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(SelectionModel); ok {
		return m.choice, nil
	}

	return "", fmt.Errorf("unexpected model type")
}

// RunTextInput runs an interactive text input and returns the entered value
func (p *BubbleTeaPrompter) RunTextInput(title, help, defaultValue string) (string, error) {
	model := NewTextInputModel(title, help, defaultValue)

	program := tea.NewProgram(model)
	finalModel, err := program.Run()
	if err != nil {
		return "", err
	}

	if m, ok := finalModel.(TextInputModel); ok {
		return m.value, nil
	}

	return "", fmt.Errorf("unexpected model type")
}

// isInteractiveMode determines if we need to prompt for additional options
func (p *BubbleTeaPrompter) isInteractiveMode(initial types.ProjectConfig) bool {
	return initial.Name == "" || initial.Module == "" || initial.Type == ""
}

// promptProjectName prompts for project name with suggestions using Bubble Tea UI
func (p *BubbleTeaPrompter) promptProjectName() (string, error) {
	suggestion := utils.GenerateRandomProjectName()
	alternatives := utils.GenerateMultipleNames(3)

	help := fmt.Sprintf("Press Enter for: %s\nOther suggestions: %s",
		suggestion, strings.Join(alternatives, ", "))

	return p.RunTextInput("🚀 What's your project name?", help, suggestion)
}

// promptModulePath prompts for Go module path using Bubble Tea UI
func (p *BubbleTeaPrompter) promptModulePath(projectName string) (string, error) {
	defaultModule := fmt.Sprintf("github.com/username/%s", projectName)
	help := "Go module path for imports (e.g., github.com/username/project)"

	return p.RunTextInput("Module path:", help, defaultModule)
}

// promptProjectType prompts for project type selection using Bubble Tea UI
func (p *BubbleTeaPrompter) promptProjectType() (string, error) {
	items := []interfaces.SelectionItem{
		interfaces.NewSelectionItem("Web API", "REST API or web service", "web-api"),
		interfaces.NewSelectionItem("CLI Application", "Command-line tool", "cli"),
		interfaces.NewSelectionItem("Library", "Reusable Go package", "library"),
		interfaces.NewSelectionItem("AWS Lambda", "Serverless function", "lambda"),
	}

	return p.RunSelection("What type of project?", items)
}

// promptFramework prompts for framework selection using Bubble Tea UI
func (p *BubbleTeaPrompter) promptFramework(projectType string) (string, error) {
	var items []interfaces.SelectionItem

	switch projectType {
	case "web-api":
		items = []interfaces.SelectionItem{
			interfaces.NewSelectionItem("Gin", "Fast HTTP web framework (recommended)", "gin"),
			interfaces.NewSelectionItem("Echo", "High performance, minimalist web framework", "echo"),
			interfaces.NewSelectionItem("Fiber", "Express inspired web framework", "fiber"),
			interfaces.NewSelectionItem("Chi", "Lightweight, idiomatic router", "chi"),
			interfaces.NewSelectionItem("Standard library", "Built-in net/http package", "standard"),
		}
	case "cli":
		items = []interfaces.SelectionItem{
			interfaces.NewSelectionItem("Cobra", "Powerful CLI framework (recommended)", "cobra"),
			interfaces.NewSelectionItem("Standard library", "Built-in flag package", "standard"),
		}
	default:
		// No framework selection needed
		return "", nil
	}

	return p.RunSelection("Which framework?", items)
}

// promptLogger prompts for logger selection using Bubble Tea UI
func (p *BubbleTeaPrompter) promptLogger(projectType string) (string, error) {
	// Skip logger selection for library projects
	if projectType == "library" {
		return "slog", nil
	}

	items := []interfaces.SelectionItem{
		interfaces.NewSelectionItem("slog", "Go built-in structured logging (recommended)", "slog"),
		interfaces.NewSelectionItem("zap", "High-performance, zero-allocation logging", "zap"),
		interfaces.NewSelectionItem("logrus", "Feature-rich, popular logging library", "logrus"),
		interfaces.NewSelectionItem("zerolog", "Zero allocation, chainable API logging", "zerolog"),
	}

	return p.RunSelection("Which logger?", items)
}

// promptGoVersion prompts for Go version selection
func (p *BubbleTeaPrompter) promptGoVersion() (string, error) {
	items := []interfaces.SelectionItem{
		interfaces.NewSelectionItem("1.21", "Stable LTS release (recommended)", "1.21"),
		interfaces.NewSelectionItem("1.22", "Latest stable release", "1.22"),
		interfaces.NewSelectionItem("1.23", "Latest release", "1.23"),
	}

	return p.RunSelection("Which Go version?", items)
}

// promptBasicOptions prompts for basic configuration options
func (p *BubbleTeaPrompter) promptBasicOptions(config *types.ProjectConfig) error {
	if config.Type == "web-api" {
		return p.promptDatabaseSupport(config)
	}
	return nil
}

// promptAdvancedOptions prompts for advanced configuration options
func (p *BubbleTeaPrompter) promptAdvancedOptions(config *types.ProjectConfig) error {
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
	return nil
}

// promptArchitecture prompts for architecture pattern using Bubble Tea UI
func (p *BubbleTeaPrompter) promptArchitecture(config *types.ProjectConfig) error {
	items := []interfaces.SelectionItem{
		interfaces.NewSelectionItem("Standard", "Simple, straightforward structure", "standard"),
		interfaces.NewSelectionItem("Clean Architecture", "Uncle Bob's principles", "clean"),
		interfaces.NewSelectionItem("Domain-Driven Design", "Business-focused design", "ddd"),
		interfaces.NewSelectionItem("Hexagonal", "Ports and adapters pattern", "hexagonal"),
	}

	choice, err := p.RunSelection("Architecture pattern?", items)
	if err != nil {
		return err
	}
	config.Architecture = choice
	return nil
}

// promptDatabaseSupport prompts for database configuration
func (p *BubbleTeaPrompter) promptDatabaseSupport(config *types.ProjectConfig) error {
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	fmt.Print(questionStyle.Render("Add database support? (y/N): "))
	fmt.Print(optionStyle.Render(" "))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		response = "n"
	}

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		config.Features.Database.Drivers = []string{"postgres"}
		config.Features.Database.ORM = ""
		fmt.Println(selectedStyle.Render("✓ Database support enabled with PostgreSQL"))
	}

	return nil
}

// promptAuthentication prompts for authentication configuration
func (p *BubbleTeaPrompter) promptAuthentication(config *types.ProjectConfig) error {
	if config.Features == nil {
		config.Features = &types.Features{}
	}

	fmt.Print(questionStyle.Render("Add authentication? (y/N): "))
	fmt.Print(optionStyle.Render(" "))
	var response string
	if _, err := fmt.Scanln(&response); err != nil {
		response = "n"
	}

	if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
		config.Features.Authentication.Type = "jwt"
		fmt.Println(selectedStyle.Render("✓ JWT authentication enabled"))
	}

	return nil
}