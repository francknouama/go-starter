package prompts

import (
	"github.com/AlecAivazis/survey/v2"
)

// MockSurveyAdapter is a mock implementation of SurveyAdapter for testing
type MockSurveyAdapter struct {
	responses map[string]interface{}
}

// NewMockSurveyAdapter creates a new mock adapter with predefined responses
func NewMockSurveyAdapter(responses map[string]interface{}) *MockSurveyAdapter {
	return &MockSurveyAdapter{
		responses: responses,
	}
}

// AskOne mocks the survey.AskOne function
func (m *MockSurveyAdapter) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	switch prompt := p.(type) {
	case *survey.Input:
		if val, ok := m.responses[prompt.Message]; ok {
			*response.(*string) = val.(string)
		}
	case *survey.Select:
		if val, ok := m.responses[prompt.Message]; ok {
			*response.(*string) = val.(string)
		}
	case *survey.MultiSelect:
		if val, ok := m.responses[prompt.Message]; ok {
			*response.(*[]string) = val.([]string)
		}
	case *survey.Confirm:
		if val, ok := m.responses[prompt.Message]; ok {
			*response.(*bool) = val.(bool)
		}
	}
	return nil
}