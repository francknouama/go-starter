package templates

import (
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestRegistry_Register(t *testing.T) {
	registry := NewRegistry()

	template := types.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Description: "A test template",
		Type:        "test",
	}

	err := registry.Register(template)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Test registering template with empty ID
	emptyTemplate := types.Template{}
	err = registry.Register(emptyTemplate)
	if err == nil {
		t.Error("Expected error for empty template ID, got nil")
	}
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()

	template := types.Template{
		ID:          "test-template",
		Name:        "Test Template",
		Description: "A test template",
		Type:        "test",
	}

	_ = registry.Register(template)

	// Test getting existing template
	retrieved, err := registry.Get("test-template")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if retrieved.ID != template.ID {
		t.Errorf("Expected template ID %s, got %s", template.ID, retrieved.ID)
	}

	// Test getting non-existent template
	_, err = registry.Get("non-existent")
	if err == nil {
		t.Error("Expected error for non-existent template, got nil")
	}
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	// Test registry with pre-loaded templates (web-api, cli, library, lambda)
	templates := registry.List()
	if len(templates) != 4 {
		t.Errorf("Expected 4 templates (pre-loaded), got %d", len(templates))
	}

	// Add some templates
	template1 := types.Template{ID: "template1", Name: "Template 1", Type: "test"}
	template2 := types.Template{ID: "template2", Name: "Template 2", Type: "test"}

	_ = registry.Register(template1)
	_ = registry.Register(template2)

	templates = registry.List()
	if len(templates) != 6 {
		t.Errorf("Expected 6 templates (4 pre-loaded + 2 added), got %d", len(templates))
	}
}

func TestRegistry_Exists(t *testing.T) {
	registry := NewRegistry()

	template := types.Template{
		ID:   "test-template",
		Name: "Test Template",
		Type: "test",
	}

	// Test non-existent template
	if registry.Exists("test-template") {
		t.Error("Expected template to not exist")
	}

	_ = registry.Register(template)

	// Test existing template
	if !registry.Exists("test-template") {
		t.Error("Expected template to exist")
	}
}

func TestRegistry_GetByType(t *testing.T) {
	registry := NewRegistry()

	template1 := types.Template{ID: "api1", Name: "API 1", Type: "api"}
	template2 := types.Template{ID: "cli1", Name: "CLI 1", Type: "cli"}
	template3 := types.Template{ID: "api2", Name: "API 2", Type: "api"}

	_ = registry.Register(template1)
	_ = registry.Register(template2)
	_ = registry.Register(template3)

	apiTemplates := registry.GetByType("api")
	if len(apiTemplates) != 2 {
		t.Errorf("Expected 2 API templates, got %d", len(apiTemplates))
	}

	cliTemplates := registry.GetByType("cli")
	if len(cliTemplates) != 2 {
		t.Errorf("Expected 2 CLI templates (1 pre-loaded + 1 added), got %d", len(cliTemplates))
	}

	nonExistentTemplates := registry.GetByType("non-existent")
	if len(nonExistentTemplates) != 0 {
		t.Errorf("Expected 0 non-existent templates, got %d", len(nonExistentTemplates))
	}
}
