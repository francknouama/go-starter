package prompts

import "github.com/AlecAivazis/survey/v2"

// SurveyAdapter defines an interface for the survey.AskOne function
type SurveyAdapter interface {
	AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error
}

// RealSurveyAdapter is a concrete implementation that calls the actual survey.AskOne
type RealSurveyAdapter struct{}

// AskOne calls the real survey.AskOne function
func (r *RealSurveyAdapter) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	return survey.AskOne(p, response, opts...)
}
