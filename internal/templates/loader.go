package templates

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/francknouama/go-starter/pkg/types"
	"gopkg.in/yaml.v3"
)

// TemplateLoader loads templates from the embedded filesystem
type TemplateLoader struct {
	fs fs.FS
}

// NewTemplateLoader creates a new template loader
func NewTemplateLoader() *TemplateLoader {
	return &TemplateLoader{
		fs: GetTemplatesFS(),
	}
}

// LoadAll loads all templates from the embedded filesystem
func (l *TemplateLoader) LoadAll() ([]types.Template, error) {
	var templates []types.Template

	// Walk through the filesystem to find all template.yaml files
	err := fs.WalkDir(l.fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip if not a template.yaml file
		if d.IsDir() || filepath.Base(path) != "template.yaml" {
			return nil
		}

		// Load the template
		template, err := l.LoadTemplate(filepath.Dir(path))
		if err != nil {
			return fmt.Errorf("failed to load template from %s: %w", path, err)
		}

		templates = append(templates, template)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk templates directory: %w", err)
	}

	return templates, nil
}

// LoadTemplate loads a single template from a directory
func (l *TemplateLoader) LoadTemplate(templateDir string) (types.Template, error) {
	templatePath := filepath.Join(templateDir, "template.yaml")
	
	// Read template.yaml
	file, err := l.fs.Open(templatePath)
	if err != nil {
		return types.Template{}, fmt.Errorf("failed to open template.yaml: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return types.Template{}, fmt.Errorf("failed to read template.yaml: %w", err)
	}

	// Parse template metadata
	var template types.Template
	if err := yaml.Unmarshal(data, &template); err != nil {
		return types.Template{}, fmt.Errorf("failed to parse template.yaml: %w", err)
	}

	// Set the template ID based on type and architecture
	if template.ID == "" {
		if template.Architecture != "" && template.Architecture != "standard" {
			template.ID = fmt.Sprintf("%s-%s", template.Type, template.Architecture)
		} else {
			template.ID = template.Type
		}
	}

	// Add template directory to metadata
	if template.Metadata == nil {
		template.Metadata = make(map[string]any)
	}
	template.Metadata["path"] = templateDir

	return template, nil
}

// LoadTemplateFile loads a template file content
func (l *TemplateLoader) LoadTemplateFile(templateDir, filePath string) (string, error) {
	fullPath := filepath.Join(templateDir, filePath)
	
	file, err := l.fs.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to open template file %s: %w", filePath, err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read template file %s: %w", filePath, err)
	}

	return string(content), nil
}

// GetTemplatePath returns the full path for a template file
func (l *TemplateLoader) GetTemplatePath(templateDir, filePath string) string {
	// Remove .tmpl extension if present
	filePath = strings.TrimSuffix(filePath, ".tmpl")
	return filepath.Join(templateDir, filePath)
}

// FileExists checks if a template file exists
func (l *TemplateLoader) FileExists(templateDir, filePath string) bool {
	fullPath := filepath.Join(templateDir, filePath)
	file, err := l.fs.Open(fullPath)
	if err != nil {
		return false
	}
	file.Close()
	return true
}