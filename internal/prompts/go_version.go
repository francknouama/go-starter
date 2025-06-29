package prompts

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"slices"
)

// Supported Go versions (latest 3 major versions)
var supportedGoVersions = []string{"auto", "1.23", "1.22", "1.21"}

// GetSupportedGoVersions returns the list of supported Go versions
func GetSupportedGoVersions() []string {
	return supportedGoVersions
}

// mapSelectionToVersion converts user selection to version string
func mapSelectionToVersion(selection string) string {
	switch selection {
	case "Auto-detect (recommended)":
		return "auto"
	case "Go 1.23 (latest)":
		return "1.23"
	case "Go 1.22":
		return "1.22"
	case "Go 1.21":
		return "1.21"
	default:
		return "auto"
	}
}

// PromptGoVersion prompts the user to select a Go version
func (p *SurveyPrompter) PromptGoVersion() (string, error) {
	options := []string{
		"Auto-detect (recommended)",
		"Go 1.23 (latest)",
		"Go 1.22",
		"Go 1.21",
	}

	prompt := &survey.Select{
		Message: "Select Go version:",
		Options: options,
		Default: "Auto-detect (recommended)",
		Help:    "Choose the Go version for your project. Auto-detect uses your system's current Go version.",
	}

	var selection string
	if err := p.surveyAdapter.AskOne(prompt, &selection); err != nil {
		return "", err
	}

	return mapSelectionToVersion(selection), nil
}

// ValidateGoVersion validates if the provided Go version is supported
func ValidateGoVersion(version string) error {
	if slices.Contains(supportedGoVersions, version) {
		return nil
	}
	return fmt.Errorf("unsupported Go version: %s. Supported versions: %v", version, supportedGoVersions)
}
