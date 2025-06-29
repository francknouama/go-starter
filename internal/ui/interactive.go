package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Interactive selection styles
var (
	listStyle = lipgloss.NewStyle().
			Margin(1, 2)

	selectedItemStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("12")).
				Background(lipgloss.Color("0")).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("15"))

	titleListStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("10")).
			MarginBottom(1)
)

// SelectionItem represents an item in a selection list
type SelectionItem struct {
	title       string
	description string
	value       string
}

func (i SelectionItem) FilterValue() string { return i.title }
func (i SelectionItem) Title() string       { return i.title }
func (i SelectionItem) Description() string { return i.description }

// SelectionModel represents an interactive selection model
type SelectionModel struct {
	list     list.Model
	choice   string
	quitting bool
	title    string
}

// NewSelectionModel creates a new selection model
func NewSelectionModel(title string, items []SelectionItem) SelectionModel {
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
			i, ok := m.list.SelectedItem().(SelectionItem)
			if ok {
				m.choice = i.value
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

// RunSelection runs an interactive selection and returns the chosen value
func RunSelection(title string, items []SelectionItem) (string, error) {
	model := NewSelectionModel(title, items)

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		// Fallback to numbered selection if TTY is not available
		return runNumberedSelection(title, items)
	}

	if m, ok := finalModel.(SelectionModel); ok {
		return m.choice, nil
	}

	return "", fmt.Errorf("unexpected model type")
}

// runNumberedSelection provides fallback numbered selection when TTY is not available
func runNumberedSelection(title string, items []SelectionItem) (string, error) {
	fmt.Print(titleListStyle.Render(title))
	fmt.Print("\n\n")

	for i, item := range items {
		fmt.Printf("  %d) %s - %s\n", i+1,
			normalItemStyle.Render(item.title),
			item.description)
	}

	fmt.Printf("\nSelect option (1-%d): ", len(items))

	var choice int
	n, err := fmt.Scanf("%d", &choice)
	if err != nil || n == 0 || choice < 1 || choice > len(items) {
		// Default to first option
		fmt.Printf("Using default: %s\n", selectedItemStyle.Render(items[0].title))
		return items[0].value, nil
	}

	fmt.Printf("Selected: %s\n", selectedItemStyle.Render(items[choice-1].title))
	return items[choice-1].value, nil
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

// RunTextInput runs an interactive text input and returns the entered value
func RunTextInput(title, help, defaultValue string) (string, error) {
	model := NewTextInputModel(title, help, defaultValue)

	p := tea.NewProgram(model)
	finalModel, err := p.Run()
	if err != nil {
		// Fallback to simple input if TTY is not available
		return runSimpleTextInput(title, help, defaultValue)
	}

	if m, ok := finalModel.(TextInputModel); ok {
		return m.value, nil
	}

	return "", fmt.Errorf("unexpected model type")
}

// runSimpleTextInput provides fallback text input when TTY is not available
func runSimpleTextInput(title, help, defaultValue string) (string, error) {
	fmt.Print(titleListStyle.Render(title))
	fmt.Print("\n")
	fmt.Print(normalItemStyle.Render(help))
	fmt.Printf("\nPress Enter for: %s\n> ", defaultValue)

	var input string
	n, err := fmt.Scanln(&input)
	if err != nil || n == 0 || input == "" {
		input = defaultValue
	}

	return input, nil
}
