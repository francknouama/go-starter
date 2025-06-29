package mocks

import (
	"fmt"

	"github.com/francknouama/go-starter/pkg/types"
)

// MockPrompter implements a mock prompter for testing
type MockPrompter struct {
	responses map[string]string
	errors    map[string]error
}

// NewMockPrompter creates a new mock prompter
func NewMockPrompter() *MockPrompter {
	return &MockPrompter{
		responses: make(map[string]string),
		errors:    make(map[string]error),
	}
}

// SetResponse sets a mock response for a prompt
func (m *MockPrompter) SetResponse(promptKey string, response string) {
	m.responses[promptKey] = response
}

// SetError sets a mock error for a prompt
func (m *MockPrompter) SetError(promptKey string, err error) {
	m.errors[promptKey] = err
}

// PromptGoVersion mocks Go version selection
func (m *MockPrompter) PromptGoVersion() (string, error) {
	if err := m.errors["goVersion"]; err != nil {
		return "", err
	}
	response := m.responses["goVersion"]
	return m.mapSelectionToVersion(response), nil
}

// PromptProjectName mocks project name input
func (m *MockPrompter) PromptProjectName() (string, error) {
	if err := m.errors["projectName"]; err != nil {
		return "", err
	}
	return m.responses["projectName"], nil
}

// PromptModulePath mocks module path input
func (m *MockPrompter) PromptModulePath(projectName string) (string, error) {
	if err := m.errors["modulePath"]; err != nil {
		return "", err
	}
	response := m.responses["modulePath"]
	if response == "" {
		response = fmt.Sprintf("github.com/user/%s", projectName)
	}
	return response, nil
}

// PromptProjectType mocks project type selection
func (m *MockPrompter) PromptProjectType() (string, error) {
	if err := m.errors["projectType"]; err != nil {
		return "", err
	}
	response := m.responses["projectType"]
	if response == "" {
		response = "web-api"
	}
	return response, nil
}

// PromptLogger mocks logger selection
func (m *MockPrompter) PromptLogger() (string, error) {
	if err := m.errors["logger"]; err != nil {
		return "", err
	}
	response := m.responses["logger"]
	if response == "" {
		response = "slog"
	}
	return response, nil
}

// BuildInteractiveConfig mocks the complete interactive configuration building
func (m *MockPrompter) BuildInteractiveConfig() (*types.ProjectConfig, error) {
	projectName, err := m.PromptProjectName()
	if err != nil {
		return nil, err
	}

	modulePath, err := m.PromptModulePath(projectName)
	if err != nil {
		return nil, err
	}

	projectType, err := m.PromptProjectType()
	if err != nil {
		return nil, err
	}

	logger, err := m.PromptLogger()
	if err != nil {
		return nil, err
	}

	goVersion, err := m.PromptGoVersion()
	if err != nil {
		return nil, err
	}

	return &types.ProjectConfig{
		Name:      projectName,
		Module:    modulePath,
		Type:      projectType,
		Logger:    logger,
		GoVersion: goVersion,
	}, nil
}

// mapSelectionToVersion converts user selection to version string
func (m *MockPrompter) mapSelectionToVersion(selection string) string {
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
