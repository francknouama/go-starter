package templates

import (
	"io/fs"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// RED PHASE: Write failing tests for walkTemplatesFromRoot
func TestTemplateLoader_walkTemplatesFromRoot(t *testing.T) {
	tests := []struct {
		name          string
		setupFS       func() fs.FS
		expectedCount int
		shouldError   bool
		errorContains string
	}{
		{
			name: "should find templates in filesystem walk",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: web-api
name: Web API
type: api
description: REST API with Go
version: "1.0.0"
architecture: standard
`),
					},
					"cli/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: cli
name: CLI Tool  
type: cli
description: Command line tool
version: "1.0.0"
architecture: standard
`),
					},
				}
			},
			expectedCount: 2,
			shouldError:   false,
		},
		{
			name: "should handle no templates found",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"README.md": &fstest.MapFile{Data: []byte("# Project")},
				}
			},
			expectedCount: 0,
			shouldError:   false,
		},
		{
			name: "should handle nested template directories",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"apis/web-api/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: web-api
name: Web API
type: api
description: REST API
version: "1.0.0"
`),
					},
				}
			},
			expectedCount: 1,
			shouldError:   false,
		},
		{
			name: "should handle malformed template.yaml",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"bad-template/template.yaml": &fstest.MapFile{
						Data: []byte("invalid: yaml: content: ["),
					},
				}
			},
			expectedCount: 0,
			shouldError:   true,
			errorContains: "failed to load template",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{fs: tt.setupFS()}

			// Act
			templates, err := loader.walkTemplatesFromRoot()

			// Assert
			if tt.shouldError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.Len(t, templates, tt.expectedCount)
			}
		})
	}
}

// RED PHASE: Write failing tests for LoadTemplateFile
func TestTemplateLoader_LoadTemplateFile(t *testing.T) {
	tests := []struct {
		name            string
		setupFS         func() fs.FS
		templateDir     string
		filePath        string
		expectedContent string
		shouldError     bool
		errorContains   string
	}{
		{
			name: "should load template file content successfully",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/main.go.tmpl": &fstest.MapFile{
						Data: []byte("package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"{{.ProjectName}}\")\n}"),
					},
				}
			},
			templateDir:     "web-api",
			filePath:        "main.go.tmpl",
			expectedContent: "package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"{{.ProjectName}}\")\n}",
			shouldError:     false,
		},
		{
			name: "should handle file not found error",
			setupFS: func() fs.FS {
				return fstest.MapFS{}
			},
			templateDir:   "web-api",
			filePath:      "nonexistent.tmpl",
			shouldError:   true,
			errorContains: "failed to open template file",
		},
		{
			name: "should load empty file content",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"template/empty.tmpl": &fstest.MapFile{
						Data: []byte(""),
					},
				}
			},
			templateDir:     "template",
			filePath:        "empty.tmpl",
			expectedContent: "",
			shouldError:     false,
		},
		{
			name: "should handle nested file paths",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/internal/handler/user.go.tmpl": &fstest.MapFile{
						Data: []byte("package handler\n\n// User handler for {{.ProjectName}}"),
					},
				}
			},
			templateDir:     "web-api",
			filePath:        "internal/handler/user.go.tmpl",
			expectedContent: "package handler\n\n// User handler for {{.ProjectName}}",
			shouldError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{fs: tt.setupFS()}

			// Act
			content, err := loader.LoadTemplateFile(tt.templateDir, tt.filePath)

			// Assert
			if tt.shouldError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				assert.Empty(t, content)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedContent, content)
			}
		})
	}
}

// RED PHASE: Write failing tests for GetTemplatePath
func TestTemplateLoader_GetTemplatePath(t *testing.T) {
	tests := []struct {
		name         string
		templateDir  string
		filePath     string
		expectedPath string
	}{
		{
			name:         "should construct path without .tmpl extension",
			templateDir:  "web-api",
			filePath:     "main.go.tmpl",
			expectedPath: filepath.Join("web-api", "main.go"),
		},
		{
			name:         "should handle path without .tmpl extension",
			templateDir:  "cli",
			filePath:     "cmd/root.go",
			expectedPath: filepath.Join("cli", "cmd", "root.go"),
		},
		{
			name:         "should handle nested paths with .tmpl extension",
			templateDir:  "web-api",
			filePath:     "internal/handler/user.go.tmpl",
			expectedPath: filepath.Join("web-api", "internal", "handler", "user.go"),
		},
		{
			name:         "should handle empty template directory",
			templateDir:  "",
			filePath:     "main.go.tmpl",
			expectedPath: "main.go",
		},
		{
			name:         "should handle empty file path",
			templateDir:  "web-api",
			filePath:     "",
			expectedPath: "web-api",
		},
		{
			name:         "should handle multiple .tmpl extensions",
			templateDir:  "test",
			filePath:     "file.tmpl.tmpl",
			expectedPath: filepath.Join("test", "file.tmpl"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{}

			// Act
			result := loader.GetTemplatePath(tt.templateDir, tt.filePath)

			// Assert
			assert.Equal(t, tt.expectedPath, result)
		})
	}
}

// RED PHASE: Write failing tests for FileExists
func TestTemplateLoader_FileExists(t *testing.T) {
	tests := []struct {
		name        string
		setupFS     func() fs.FS
		templateDir string
		filePath    string
		expected    bool
	}{
		{
			name: "should return true for existing file",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/main.go.tmpl": &fstest.MapFile{
						Data: []byte("package main"),
					},
				}
			},
			templateDir: "web-api",
			filePath:    "main.go.tmpl",
			expected:    true,
		},
		{
			name: "should return false for non-existing file",
			setupFS: func() fs.FS {
				return fstest.MapFS{}
			},
			templateDir: "web-api",
			filePath:    "nonexistent.tmpl",
			expected:    false,
		},
		{
			name: "should return true for existing directory",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/internal/dummy": &fstest.MapFile{
						Data: []byte(""),
					},
				}
			},
			templateDir: "web-api",
			filePath:    "internal",
			expected:    true,
		},
		{
			name: "should handle nested file paths",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"template/internal/handler/user.go.tmpl": &fstest.MapFile{
						Data: []byte("package handler"),
					},
				}
			},
			templateDir: "template",
			filePath:    "internal/handler/user.go.tmpl",
			expected:    true,
		},
		{
			name: "should handle empty directory and path",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"file.txt": &fstest.MapFile{
						Data: []byte("content"),
					},
				}
			},
			templateDir: "",
			filePath:    "file.txt",
			expected:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{fs: tt.setupFS()}

			// Act
			result := loader.FileExists(tt.templateDir, tt.filePath)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Additional tests for better coverage of existing functionality
func TestTemplateLoader_Integration(t *testing.T) {
	tests := []struct {
		name          string
		setupFS       func() fs.FS
		expectedCount int
		shouldError   bool
		errorContains string
	}{
		{
			name: "should load multiple templates with different architectures",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: web-api
name: Web API
type: api
description: REST API
version: "1.0.0"
architecture: standard
`),
					},
					"web-api-clean/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: web-api-clean
name: Web API Clean
type: api
description: Clean Architecture API
version: "1.0.0"
architecture: clean
`),
					},
					"cli-tool/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: cli-tool
name: CLI Tool
type: cli
description: Command line tool
version: "1.0.0"
`),
					},
				}
			},
			expectedCount: 3,
			shouldError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{fs: tt.setupFS()}

			// Act
			templates, err := loader.LoadAll()

			// Assert
			if tt.shouldError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				assert.Len(t, templates, tt.expectedCount)

				// Verify each template has proper metadata
				for _, template := range templates {
					assert.NotEmpty(t, template.ID)
					assert.NotEmpty(t, template.Name)
					assert.NotEmpty(t, template.Type)
					assert.NotNil(t, template.Metadata)
					assert.Contains(t, template.Metadata, "path")
				}
			}
		})
	}
}

// RED PHASE: Write tests for LoadTemplate edge cases not covered
func TestTemplateLoader_LoadTemplate_EdgeCases(t *testing.T) {
	tests := []struct {
		name          string
		setupFS       func() fs.FS
		templateDir   string
		shouldError   bool
		errorContains string
		validateFn    func(t *testing.T, template types.Template)
	}{
		{
			name: "should generate ID from type when ID is empty",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api/template.yaml": &fstest.MapFile{
						Data: []byte(`
name: Web API
type: api
description: REST API
version: "1.0.0"
architecture: standard
`),
					},
				}
			},
			templateDir: "web-api",
			shouldError: false,
			validateFn: func(t *testing.T, template types.Template) {
				assert.Equal(t, "api", template.ID)
			},
		},
		{
			name: "should generate ID with architecture when not standard",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"web-api-clean/template.yaml": &fstest.MapFile{
						Data: []byte(`
name: Web API Clean
type: api
description: Clean Architecture API
version: "1.0.0"
architecture: clean
`),
					},
				}
			},
			templateDir: "web-api-clean",
			shouldError: false,
			validateFn: func(t *testing.T, template types.Template) {
				assert.Equal(t, "api-clean", template.ID)
			},
		},
		{
			name: "should preserve existing ID when provided",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"custom/template.yaml": &fstest.MapFile{
						Data: []byte(`
id: custom-template-id
name: Custom Template
type: custom
description: Custom template
version: "1.0.0"
`),
					},
				}
			},
			templateDir: "custom",
			shouldError: false,
			validateFn: func(t *testing.T, template types.Template) {
				assert.Equal(t, "custom-template-id", template.ID)
			},
		},
		{
			name: "should add template path to metadata",
			setupFS: func() fs.FS {
				return fstest.MapFS{
					"api/template.yaml": &fstest.MapFile{
						Data: []byte(`
name: API Template
type: api
description: API template
version: "1.0.0"
`),
					},
				}
			},
			templateDir: "api",
			shouldError: false,
			validateFn: func(t *testing.T, template types.Template) {
				require.NotNil(t, template.Metadata)
				assert.Equal(t, "api", template.Metadata["path"])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			loader := &TemplateLoader{fs: tt.setupFS()}

			// Act
			template, err := loader.LoadTemplate(tt.templateDir)

			// Assert
			if tt.shouldError {
				require.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				require.NoError(t, err)
				if tt.validateFn != nil {
					tt.validateFn(t, template)
				}
			}
		})
	}
}
