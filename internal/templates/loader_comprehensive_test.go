package templates

import (
	"testing"
	"testing/fstest"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTemplateLoaderComprehensive tests the template loader with various scenarios
func TestTemplateLoaderComprehensive(t *testing.T) {
	t.Run("LoadTemplate_WithIncludes", func(t *testing.T) {
		// Create a mock filesystem with template and include files
		mockFS := fstest.MapFS{
			"web-api-test/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "web-api-test"
description: "Test template with includes"
type: "web-api"
architecture: "standard"
version: "1.0.0"
author: "Test"
license: "MIT"

include:
  variables: "config/variables.yaml"
  dependencies: "config/dependencies.yaml"
  features: "config/features.yaml"

files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
			},
			"web-api-test/config/variables.yaml": &fstest.MapFile{
				Data: []byte(`
variables:
  - name: "ProjectName"
    description: "Name of the project"
    type: "string"
    required: true
  - name: "Framework"
    description: "Web framework"
    type: "string"
    default: "gin"
    choices:
      - "gin"
      - "echo"
`),
			},
			"web-api-test/config/dependencies.yaml": &fstest.MapFile{
				Data: []byte(`
dependencies:
  - module: "github.com/gin-gonic/gin"
    version: "v1.9.1"
    condition: "{{eq .Framework \"gin\"}}"
  - module: "github.com/labstack/echo/v4"
    version: "v4.11.3"
    condition: "{{eq .Framework \"echo\"}}"
`),
			},
			"web-api-test/config/features.yaml": &fstest.MapFile{
				Data: []byte(`
features:
  - name: "logging"
    description: "Structured logging"
    enabled_when: "true"

validation:
  - name: "project_name_format"
    description: "Validate project name"

post_hooks:
  - name: "format_code"
    command: "go fmt ./..."
`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		template, err := loader.LoadTemplate("web-api-test")

		require.NoError(t, err)
		assert.Equal(t, "web-api-test", template.Name)
		assert.Equal(t, "web-api", template.Type)
		assert.Equal(t, "standard", template.Architecture)

		// Verify includes were processed
		assert.Len(t, template.Variables, 2)
		assert.Equal(t, "ProjectName", template.Variables[0].Name)
		assert.Equal(t, "Framework", template.Variables[1].Name)

		assert.Len(t, template.Dependencies, 2)
		assert.Equal(t, "github.com/gin-gonic/gin", template.Dependencies[0].Module)

		assert.Len(t, template.Features, 1)
		assert.Equal(t, "logging", template.Features[0].Name)

		assert.Len(t, template.Validation, 1)
		assert.Equal(t, "project_name_format", template.Validation[0].Name)

		assert.Len(t, template.PostHooks, 1)
		assert.Equal(t, "format_code", template.PostHooks[0].Name)
	})

	t.Run("LoadTemplate_WithoutIncludes", func(t *testing.T) {
		mockFS := fstest.MapFS{
			"simple-template/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "simple-template"
description: "Simple template without includes"
type: "cli"
variables:
  - name: "ProjectName"
    type: "string"
    required: true
files:
  - source: "main.go.tmpl"
    destination: "main.go"
dependencies:
  - module: "github.com/spf13/cobra"
    version: "v1.7.0"
`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		template, err := loader.LoadTemplate("simple-template")

		require.NoError(t, err)
		assert.Equal(t, "simple-template", template.Name)
		assert.Equal(t, "cli", template.Type)
		assert.Len(t, template.Variables, 1)
		assert.Len(t, template.Dependencies, 1)
		assert.Nil(t, template.Include)
	})

	t.Run("LoadTemplate_ArchitectureID", func(t *testing.T) {
		testCases := []struct {
			name         string
			architecture string
			expectedID   string
		}{
			{"Standard architecture", "standard", "web-api"},
			{"Clean architecture", "clean", "web-api-clean"},
			{"DDD architecture", "ddd", "web-api-ddd"},
			{"Empty architecture", "", "web-api"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockFS := fstest.MapFS{
					"test-template/template.yaml": &fstest.MapFile{
						Data: []byte(`
name: "test-template"
type: "web-api"
architecture: "` + tc.architecture + `"
files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
					},
				}

				loader := &TemplateLoader{fs: mockFS}
				template, err := loader.LoadTemplate("test-template")

				require.NoError(t, err)
				assert.Equal(t, tc.expectedID, template.ID)
			})
		}
	})

	t.Run("LoadTemplate_InvalidYAML", func(t *testing.T) {
		mockFS := fstest.MapFS{
			"invalid-template/template.yaml": &fstest.MapFile{
				Data: []byte(`invalid yaml: [[[`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		_, err := loader.LoadTemplate("invalid-template")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to parse template.yaml")
	})

	t.Run("LoadTemplate_MissingIncludeFile", func(t *testing.T) {
		mockFS := fstest.MapFS{
			"missing-include/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "missing-include"
type: "web-api"
include:
  variables: "config/missing.yaml"
files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		_, err := loader.LoadTemplate("missing-include")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to process includes")
	})
}

// TestTemplateLoaderAll tests loading all templates
func TestTemplateLoaderAll(t *testing.T) {
	t.Run("LoadAll_MultipleTemplates", func(t *testing.T) {
		mockFS := fstest.MapFS{
			"web-api/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "web-api"
type: "web-api"
files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
			},
			"cli/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "cli"
type: "cli"
files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		templates, err := loader.LoadAll()

		require.NoError(t, err)
		assert.Len(t, templates, 2)

		// Check that both templates were loaded
		templateNames := make(map[string]bool)
		for _, tmpl := range templates {
			templateNames[tmpl.Name] = true
		}

		assert.True(t, templateNames["web-api"])
		assert.True(t, templateNames["cli"])
	})

	t.Run("LoadAll_EmptyFileSystem", func(t *testing.T) {
		mockFS := fstest.MapFS{}
		loader := &TemplateLoader{fs: mockFS}
		templates, err := loader.LoadAll()

		require.NoError(t, err)
		assert.Len(t, templates, 0)
	})
}

// TestTemplateLoaderFileOperations tests file operation methods
func TestTemplateLoaderFileOperations(t *testing.T) {
	mockFS := fstest.MapFS{
		"test-template/main.go.tmpl": &fstest.MapFile{
			Data: []byte(`package main

import "fmt"

func main() {
	fmt.Println("Hello {{.ProjectName}}")
}`),
		},
		"test-template/README.md.tmpl": &fstest.MapFile{
			Data: []byte(`# {{.ProjectName}}

This is a test project.`),
		},
	}

	loader := &TemplateLoader{fs: mockFS}

	t.Run("LoadTemplateFile_Success", func(t *testing.T) {
		content, err := loader.LoadTemplateFile("test-template", "main.go.tmpl")

		require.NoError(t, err)
		assert.Contains(t, content, "Hello {{.ProjectName}}")
		assert.Contains(t, content, "package main")
	})

	t.Run("LoadTemplateFile_NotFound", func(t *testing.T) {
		_, err := loader.LoadTemplateFile("test-template", "nonexistent.tmpl")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to open template file")
	})

	t.Run("FileExists_True", func(t *testing.T) {
		exists := loader.FileExists("test-template", "main.go.tmpl")
		assert.True(t, exists)
	})

	t.Run("FileExists_False", func(t *testing.T) {
		exists := loader.FileExists("test-template", "nonexistent.tmpl")
		assert.False(t, exists)
	})

	t.Run("GetTemplatePath", func(t *testing.T) {
		path := loader.GetTemplatePath("test-template", "main.go.tmpl")
		assert.Equal(t, "test-template/main.go", path)

		// Test .tmpl extension removal
		path2 := loader.GetTemplatePath("test-template", "README.md.tmpl")
		assert.Equal(t, "test-template/README.md", path2)
	})
}

// TestIncludeProcessing tests the include processing functionality
func TestIncludeProcessing(t *testing.T) {
	t.Run("ProcessIncludes_AllTypes", func(t *testing.T) {
		mockFS := fstest.MapFS{
			"config/variables.yaml": &fstest.MapFile{
				Data: []byte(`
variables:
  - name: "TestVar"
    type: "string"
`),
			},
			"config/dependencies.yaml": &fstest.MapFile{
				Data: []byte(`
dependencies:
  - module: "test/module"
    version: "v1.0.0"
`),
			},
			"config/features.yaml": &fstest.MapFile{
				Data: []byte(`
features:
  - name: "test-feature"
    description: "Test feature"

validation:
  - name: "test-validation"
    description: "Test validation"

post_hooks:
  - name: "test-hook"
    command: "echo test"
`),
			},
		}

		loader := &TemplateLoader{fs: mockFS}
		template := &types.Template{
			Include: &types.TemplateIncludes{
				Variables:    "config/variables.yaml",
				Dependencies: "config/dependencies.yaml",
				Features:     "config/features.yaml",
			},
		}

		err := loader.processIncludes(template, ".")

		require.NoError(t, err)
		assert.Len(t, template.Variables, 1)
		assert.Equal(t, "TestVar", template.Variables[0].Name)

		assert.Len(t, template.Dependencies, 1)
		assert.Equal(t, "test/module", template.Dependencies[0].Module)

		assert.Len(t, template.Features, 1)
		assert.Equal(t, "test-feature", template.Features[0].Name)

		assert.Len(t, template.Validation, 1)
		assert.Equal(t, "test-validation", template.Validation[0].Name)

		assert.Len(t, template.PostHooks, 1)
		assert.Equal(t, "test-hook", template.PostHooks[0].Name)
	})

	t.Run("ProcessIncludes_NoIncludes", func(t *testing.T) {
		loader := &TemplateLoader{fs: fstest.MapFS{}}
		template := &types.Template{Include: nil}

		err := loader.processIncludes(template, ".")

		require.NoError(t, err)
		assert.Len(t, template.Variables, 0)
		assert.Len(t, template.Dependencies, 0)
		assert.Len(t, template.Features, 0)
	})
}