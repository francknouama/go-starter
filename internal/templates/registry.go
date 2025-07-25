package templates

import (
	"fmt"
	"sort"
	"sync"

	"github.com/francknouama/go-starter/pkg/types"
)

// Registry manages all available project templates
type Registry struct {
	templates map[string]types.Template
	mutex     sync.RWMutex
}

// NewRegistry creates a new template registry
func NewRegistry() *Registry {
	r := &Registry{
		templates: make(map[string]types.Template),
	}
	// Load embedded blueprints when they're available
	r.loadEmbeddedTemplates()
	return r
}

// Register adds a template to the registry
func (r *Registry) Register(template types.Template) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if template.ID == "" {
		return types.NewValidationError("template ID cannot be empty", nil)
	}

	r.templates[template.ID] = template
	return nil
}

// Get retrieves a template by ID
func (r *Registry) Get(templateID string) (types.Template, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	template, exists := r.templates[templateID]
	if !exists {
		return types.Template{}, types.NewTemplateNotFoundError(templateID)
	}

	return template, nil
}

// List returns all available templates
func (r *Registry) List() []types.Template {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	templates := make([]types.Template, 0, len(r.templates))
	for _, template := range r.templates {
		templates = append(templates, template)
	}
	
	// Sort templates to ensure consistent ordering
	// Priority: cli-simple first, then by type, then by name
	sort.Slice(templates, func(i, j int) bool {
		// cli-simple always comes first
		if templates[i].ID == "cli-simple" {
			return true
		}
		if templates[j].ID == "cli-simple" {
			return false
		}
		
		// Then sort by type
		if templates[i].Type != templates[j].Type {
			return templates[i].Type < templates[j].Type
		}
		
		// Finally sort by ID
		return templates[i].ID < templates[j].ID
	})

	return templates
}

// GetByType returns all templates of a specific type
func (r *Registry) GetByType(templateType string) []types.Template {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var templates []types.Template
	for _, template := range r.templates {
		if template.Type == templateType {
			templates = append(templates, template)
		}
	}

	return templates
}

// Exists checks if a template exists
func (r *Registry) Exists(templateID string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, exists := r.templates[templateID]
	return exists
}

// Remove removes a template from the registry
func (r *Registry) Remove(templateID string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if !r.exists(templateID) {
		return types.NewTemplateNotFoundError(templateID)
	}

	delete(r.templates, templateID)
	return nil
}

// GetTemplateTypes returns all unique template types
func (r *Registry) GetTemplateTypes() []string {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	typeSet := make(map[string]bool)
	for _, template := range r.templates {
		typeSet[template.Type] = true
	}

	types := make([]string, 0, len(typeSet))
	for templateType := range typeSet {
		types = append(types, templateType)
	}

	return types
}

// exists is an internal helper that doesn't lock (assumes caller has lock)
func (r *Registry) exists(templateID string) bool {
	_, exists := r.templates[templateID]
	return exists
}

// loadEmbeddedTemplates loads templates from embedded blueprint files
func (r *Registry) loadEmbeddedTemplates() {
	loader := NewTemplateLoader()

	templates, err := loader.LoadAll()
	if err != nil {
		fmt.Printf("Warning: Failed to load blueprints: %v\n", err)
		return
	}

	for _, template := range templates {
		if err := r.Register(template); err != nil {
			fmt.Printf("Warning: Failed to register template %s: %v\n", template.ID, err)
		}
	}

	if len(templates) > 0 {
		fmt.Printf("Template registry initialized (%d templates loaded)\n", len(templates))
	} else {
		fmt.Println("Warning: No blueprints found in embedded filesystem")
	}
}
