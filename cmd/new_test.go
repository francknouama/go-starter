package cmd

import (
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  types.ProjectConfig
		wantErr bool
	}{
		{
			name: "valid config",
			config: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			config: types.ProjectConfig{
				Module: "github.com/test/project",
				Type:   "web-api",
			},
			wantErr: true,
		},
		{
			name: "missing module",
			config: types.ProjectConfig{
				Name: "test-project",
				Type: "web-api",
			},
			wantErr: true,
		},
		{
			name: "missing type",
			config: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
			},
			wantErr: true,
		},
		{
			name:    "empty config",
			config:  types.ProjectConfig{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Check error type for validation errors
			if err != nil {
				if goStarterErr, ok := err.(*types.GoStarterError); ok {
					if goStarterErr.Code != types.ErrCodeValidation {
						t.Errorf("Expected validation error, got error code: %s", goStarterErr.Code)
					}
				} else {
					t.Error("Expected GoStarterError type for validation error")
				}
			}
		})
	}
}

func TestPrintSuccessMessage(t *testing.T) {
	// Test that printSuccessMessage doesn't panic
	// This is mainly a smoke test since the function prints to stdout
	config := types.ProjectConfig{
		Name:      "test-project",
		Type:      "web-api",
		Framework: "gin",
		Module:    "github.com/test/project",
	}

	result := &types.GenerationResult{
		ProjectPath:  "/tmp/test-project",
		FilesCreated: []string{"go.mod", "main.go", "README.md"},
		Success:      true,
	}

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printSuccessMessage() panicked: %v", r)
		}
	}()

	printSuccessMessage(config, result)
}

func TestPrintSuccessMessage_WithoutFramework(t *testing.T) {
	// Test printSuccessMessage with empty framework
	config := types.ProjectConfig{
		Name:   "test-project",
		Type:   "library",
		Module: "github.com/test/project",
		// Framework is empty
	}

	result := &types.GenerationResult{
		ProjectPath:  "/tmp/test-project",
		FilesCreated: []string{"go.mod", "README.md"},
		Success:      true,
	}

	// This should not panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("printSuccessMessage() panicked: %v", r)
		}
	}()

	printSuccessMessage(config, result)
}

func TestNewCommandFlags(t *testing.T) {
	// Test that the new command has the expected flags
	expectedFlags := []string{
		"name",
		"module",
		"type",
		"framework",
		"output",
		"advanced",
		"dry-run",
		"no-git",
	}

	for _, flagName := range expectedFlags {
		flag := newCmd.Flags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected flag %s to exist", flagName)
		}
	}
}

func TestNewCommandUsage(t *testing.T) {
	// Test basic command properties
	if newCmd.Use != "new [project-name]" {
		t.Errorf("Expected Use to be 'new [project-name]', got %s", newCmd.Use)
	}

	if newCmd.Short == "" {
		t.Error("Expected Short description to not be empty")
	}

	if newCmd.Long == "" {
		t.Error("Expected Long description to not be empty")
	}
}

func TestGlobalVariables(t *testing.T) {
	// Test that global variables are properly declared
	// This is mainly to ensure they exist and have expected types

	// Test string variables
	stringVars := map[string]*string{
		"projectName":   &projectName,
		"projectModule": &projectModule,
		"projectType":   &projectType,
		"outputDir":     &outputDir,
		"framework":     &framework,
	}

	for name, ptr := range stringVars {
		if ptr == nil {
			t.Errorf("String variable %s should not be nil", name)
		}
	}

	// Test boolean variables
	boolVars := map[string]*bool{
		"advanced": &advanced,
		"dryRun":   &dryRun,
		"noGit":    &noGit,
	}

	for name, ptr := range boolVars {
		if ptr == nil {
			t.Errorf("Boolean variable %s should not be nil", name)
		}
	}
}

// Tests for helper functions would go here if they were exported

func TestFlagBinding(t *testing.T) {
	// Test that flags are properly bound and accessible
	flags := newCmd.Flags()

	flagTests := []struct {
		name string
		typ  string
	}{
		{"name", "string"},
		{"module", "string"},
		{"type", "string"},
		{"framework", "string"},
		{"output", "string"},
		{"advanced", "bool"},
		{"dry-run", "bool"},
		{"no-git", "bool"},
	}

	for _, ft := range flagTests {
		t.Run(ft.name, func(t *testing.T) {
			flag := flags.Lookup(ft.name)
			if flag == nil {
				t.Fatalf("Flag %s should exist", ft.name)
			}

			if flag.Value.Type() != ft.typ {
				t.Errorf("Flag %s should be %s, got %s", ft.name, ft.typ, flag.Value.Type())
			}
		})
	}
}
