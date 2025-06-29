package ui

import (
	"errors"
	"testing"

	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
)

// MockSurveyAdapter for testing Survey fallbacks
type MockSurveyAdapter struct {
	responses map[string]interface{}
	errors    map[string]error
}

func NewMockSurveyAdapter() *MockSurveyAdapter {
	return &MockSurveyAdapter{
		responses: make(map[string]interface{}),
		errors:    make(map[string]error),
	}
}

func (m *MockSurveyAdapter) SetResponse(promptType string, response interface{}) {
	m.responses[promptType] = response
}

func (m *MockSurveyAdapter) SetError(promptType string, err error) {
	m.errors[promptType] = err
}

func (m *MockSurveyAdapter) AskOne(p survey.Prompt, response interface{}, opts ...survey.AskOpt) error {
	var promptType string
	switch p := p.(type) {
	case *survey.Input:
		promptType = p.Message
	case *survey.Select:
		promptType = p.Message
	case *survey.Confirm:
		promptType = p.Message
	case *survey.MultiSelect:
		promptType = p.Message
	default:
		promptType = "unknown"
	}

	if err, exists := m.errors[promptType]; exists {
		return err
	}

	if resp, exists := m.responses[promptType]; exists {
		switch v := response.(type) {
		case *string:
			if str, ok := resp.(string); ok {
				*v = str
			}
		case *[]string:
			if strs, ok := resp.([]string); ok {
				*v = strs
			}
		case *bool:
			if b, ok := resp.(bool); ok {
				*v = b
			}
		}
		return nil
	}

	// Default responses for common prompts
	switch response := response.(type) {
	case *string:
		*response = "default"
	case *[]string:
		*response = []string{"default"}
	case *bool:
		*response = false
	}
	return nil
}

func TestNewFangPrompter(t *testing.T) {
	prompter := NewFangPrompter()
	if prompter == nil {
		t.Fatal("Expected prompter to be created")
	}
	if !prompter.useEnhancedUI {
		t.Error("Expected enhanced UI to be enabled by default")
	}
	if prompter.surveyAdapter == nil {
		t.Error("Expected survey adapter to be initialized")
	}
}

func TestNewFangPrompterWithSurvey(t *testing.T) {
	prompter := NewFangPrompterWithSurvey()
	if prompter == nil {
		t.Fatal("Expected prompter to be created")
	}
	if prompter.useEnhancedUI {
		t.Error("Expected enhanced UI to be disabled for survey mode")
	}
	if prompter.surveyAdapter == nil {
		t.Error("Expected survey adapter to be initialized")
	}
}

func TestFangPrompter_GetProjectConfig_SurveyFallback(t *testing.T) {
	// Test with Survey fallback (no enhanced UI)
	prompter := NewFangPrompterWithSurvey()
	mockAdapter := NewMockSurveyAdapter()

	// Set up mock responses
	mockAdapter.SetResponse("What's your project name?", "test-project")
	mockAdapter.SetResponse("Module path:", "github.com/user/test-project")
	mockAdapter.SetResponse("What type of project?", "Web API - REST API or web service")
	mockAdapter.SetResponse("Which web framework?", "Gin (recommended)")
	mockAdapter.SetResponse("Which Go version?", "1.21 - Stable LTS release (recommended)")
	mockAdapter.SetResponse("Which logger?", "slog - Go built-in structured logging (recommended)")

	prompter.surveyAdapter = mockAdapter

	initial := types.ProjectConfig{}
	config, err := prompter.GetProjectConfig(initial, false)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if config.Name != "test-project" {
		t.Errorf("Expected name 'test-project', got '%s'", config.Name)
	}

	if config.Module != "github.com/user/test-project" {
		t.Errorf("Expected module 'github.com/user/test-project', got '%s'", config.Module)
	}

	if config.Type != "web-api" {
		t.Errorf("Expected type 'web-api', got '%s'", config.Type)
	}

	if config.Framework != "gin" {
		t.Errorf("Expected framework 'gin', got '%s'", config.Framework)
	}

	// Go version should be set (the actual version may vary based on optimal detection)
	if config.GoVersion == "" {
		t.Error("Expected Go version to be set")
	}

	if config.Logger != "slog" {
		t.Errorf("Expected logger 'slog', got '%s'", config.Logger)
	}
}

func TestFangPrompter_PromptProjectNameSurvey(t *testing.T) {
	prompter := NewFangPrompterWithSurvey()
	mockAdapter := NewMockSurveyAdapter()
	mockAdapter.SetResponse("What's your project name?", "my-awesome-project")
	prompter.surveyAdapter = mockAdapter

	result, err := prompter.promptProjectNameSurvey()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "my-awesome-project" {
		t.Errorf("Expected 'my-awesome-project', got '%s'", result)
	}
}

func TestFangPrompter_PromptModulePathSurvey(t *testing.T) {
	prompter := NewFangPrompterWithSurvey()
	mockAdapter := NewMockSurveyAdapter()
	mockAdapter.SetResponse("Module path:", "github.com/myuser/myproject")
	prompter.surveyAdapter = mockAdapter

	result, err := prompter.promptModulePathSurvey("myproject")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result != "github.com/myuser/myproject" {
		t.Errorf("Expected 'github.com/myuser/myproject', got '%s'", result)
	}
}

func TestFangPrompter_PromptProjectTypeSurvey(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		expected string
	}{
		{"Web API", "Web API - REST API or web service", "web-api"},
		{"CLI", "CLI Application - Command-line tool", "cli"},
		{"Library", "Library - Reusable Go package", "library"},
		{"Lambda", "AWS Lambda - Serverless function", "lambda"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompter := NewFangPrompterWithSurvey()
			mockAdapter := NewMockSurveyAdapter()
			mockAdapter.SetResponse("What type of project?", tc.response)
			prompter.surveyAdapter = mockAdapter

			result, err := prompter.promptProjectTypeSurvey()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFangPrompter_PromptFrameworkSurvey(t *testing.T) {
	testCases := []struct {
		name        string
		projectType string
		response    string
		expected    string
	}{
		{"Gin for web-api", "web-api", "Gin (recommended)", "gin"},
		{"Echo for web-api", "web-api", "Echo", "echo"},
		{"Cobra for cli", "cli", "Cobra (recommended)", "cobra"},
		{"No framework for library", "library", "", ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompter := NewFangPrompterWithSurvey()
			mockAdapter := NewMockSurveyAdapter()

			switch tc.projectType {
			case "web-api":
				mockAdapter.SetResponse("Which web framework?", tc.response)
			case "cli":
				mockAdapter.SetResponse("Which CLI framework?", tc.response)
			}

			prompter.surveyAdapter = mockAdapter

			result, err := prompter.promptFrameworkSurvey(tc.projectType)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFangPrompter_PromptLoggerSurvey(t *testing.T) {
	testCases := []struct {
		name        string
		projectType string
		response    string
		expected    string
	}{
		{"slog for web-api", "web-api", "slog - Go built-in structured logging (recommended)", "slog"},
		{"zap for web-api", "web-api", "zap - High-performance, zero-allocation logging", "zap"},
		{"Default slog for library", "library", "", "slog"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompter := NewFangPrompterWithSurvey()
			mockAdapter := NewMockSurveyAdapter()

			if tc.projectType != "library" {
				mockAdapter.SetResponse("Which logger?", tc.response)
			}

			prompter.surveyAdapter = mockAdapter

			result, err := prompter.promptLoggerSurvey(tc.projectType)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFangPrompter_PromptGoVersionSurvey(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		expected string
	}{
		{"Go 1.21", "1.21 - Stable LTS release (recommended)", "1.21"},
		{"Go 1.22", "1.22 - Latest stable release", "1.22"},
		{"Go 1.23", "1.23 - Latest release", "1.23"},
		{"Go 1.24", "1.24 - Development release", "1.24"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompter := NewFangPrompterWithSurvey()
			mockAdapter := NewMockSurveyAdapter()
			mockAdapter.SetResponse("Which Go version?", tc.response)
			prompter.surveyAdapter = mockAdapter

			result, err := prompter.promptGoVersionSurvey()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if result != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, result)
			}
		})
	}
}

func TestFangPrompter_PromptArchitectureSurvey(t *testing.T) {
	testCases := []struct {
		name     string
		response string
		expected string
	}{
		{"Standard", "Standard - Simple structure", "standard"},
		{"Clean", "Clean Architecture - Uncle Bob's principles", "clean"},
		{"DDD", "Domain-Driven Design - Business-focused", "ddd"},
		{"Hexagonal", "Hexagonal - Ports and adapters", "hexagonal"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			prompter := NewFangPrompterWithSurvey()
			mockAdapter := NewMockSurveyAdapter()
			mockAdapter.SetResponse("Architecture pattern?", tc.response)
			prompter.surveyAdapter = mockAdapter

			var config types.ProjectConfig
			err := prompter.promptArchitectureSurvey(&config)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}

			if config.Architecture != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, config.Architecture)
			}
		})
	}
}

func TestFangPrompter_PromptArchitectureSurvey_Error(t *testing.T) {
	prompter := NewFangPrompterWithSurvey()
	mockAdapter := NewMockSurveyAdapter()
	mockAdapter.SetError("Architecture pattern?", errors.New("survey error"))
	prompter.surveyAdapter = mockAdapter

	var config types.ProjectConfig
	err := prompter.promptArchitectureSurvey(&config)
	if err != nil {
		t.Fatalf("Expected no error on survey failure, got %v", err)
	}

	// Should default to standard architecture on error
	if config.Architecture != "standard" {
		t.Errorf("Expected 'standard' as default, got '%s'", config.Architecture)
	}
}

func TestFangPrompter_IsInteractiveMode(t *testing.T) {
	prompter := NewFangPrompter()

	testCases := []struct {
		name     string
		config   types.ProjectConfig
		expected bool
	}{
		{"Empty config", types.ProjectConfig{}, true},
		{"Only name", types.ProjectConfig{Name: "test"}, true},
		{"Name and module", types.ProjectConfig{Name: "test", Module: "github.com/user/test"}, true},
		{"Complete config", types.ProjectConfig{Name: "test", Module: "github.com/user/test", Type: "web-api"}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := prompter.isInteractiveMode(tc.config)
			if result != tc.expected {
				t.Errorf("Expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestFangPrompter_ErrorHandling(t *testing.T) {
	prompter := NewFangPrompterWithSurvey()
	mockAdapter := NewMockSurveyAdapter()

	// Test error propagation
	mockAdapter.SetError("What's your project name?", errors.New("input error"))
	prompter.surveyAdapter = mockAdapter

	_, err := prompter.promptProjectNameSurvey()
	if err == nil {
		t.Error("Expected error to be propagated")
	}

	if err.Error() != "input error" {
		t.Errorf("Expected 'input error', got '%s'", err.Error())
	}
}

func TestRealSurveyAdapter(t *testing.T) {
	adapter := &RealSurveyAdapter{}

	// Test that the adapter exists and has the correct interface
	assert.NotNil(t, adapter)

	// We can't easily test the actual survey interaction, but we can verify the method exists
	// This is mainly for coverage and interface verification
	_ = adapter.AskOne
}
