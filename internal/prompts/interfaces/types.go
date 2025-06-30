package interfaces

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/pkg/types"
)

// Prompter defines the interface for collecting user input during project generation
type Prompter interface {
	GetProjectConfig(initial types.ProjectConfig, advanced bool) (types.ProjectConfig, error)
}

// SelectionItem represents an item in a selection list
type SelectionItem struct {
	title       string
	description string
	value       string
}

// NewSelectionItem creates a new SelectionItem
func NewSelectionItem(title, description, value string) SelectionItem {
	return SelectionItem{
		title:       title,
		description: description,
		value:       value,
	}
}

func (i SelectionItem) FilterValue() string { return i.title }
func (i SelectionItem) Title() string       { return i.title }
func (i SelectionItem) Description() string { return i.description }
func (i SelectionItem) Value() string       { return i.value }

// SurveyAdapter interface allows for testable survey interactions
type SurveyAdapter interface {
	AskOne(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) error
}

// RealSurveyAdapter implements SurveyAdapter using the actual survey library
type RealSurveyAdapter struct{}

func (r *RealSurveyAdapter) AskOne(prompt survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(prompt, response, opts...)
}