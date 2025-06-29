package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
)

func setupTestTemplates(t *testing.T) {
	t.Helper()

	// Get the project root for tests
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	templatesDir := filepath.Join(projectRoot, "templates")

	// Verify templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Fatalf("Templates directory not found at: %s", templatesDir)
	}

	// Set up the filesystem for tests using os.DirFS
	SetTemplatesFS(os.DirFS(templatesDir))
}

func TestRegistry_Register(t *testing.T) {
	setupTestTemplates(t)

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

	// Test registry with pre-loaded templates (web-api-standard, web-api-clean, web-api-ddd, cli-standard, library-standard, lambda-standard, microservice-standard)
	templates := registry.List()
	if len(templates) != 7 {
		t.Errorf("Expected 7 templates (pre-loaded), got %d", len(templates))
	}

	// Add some templates
	template1 := types.Template{ID: "template1", Name: "Template 1", Type: "test"}
	template2 := types.Template{ID: "template2", Name: "Template 2", Type: "test"}

	_ = registry.Register(template1)
	_ = registry.Register(template2)

	templates = registry.List()
	if len(templates) != 9 {
		t.Errorf("Expected 9 templates (7 pre-loaded + 2 added), got %d", len(templates))
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

// RED PHASE: Write tests for Remove function (0% coverage)
func TestRegistry_Remove(t *testing.T) {
	setupTestTemplates(t)
	registry := NewRegistry()

	template := types.Template{
		ID:   "test-template",
		Name: "Test Template",
		Type: "test",
	}

	// Test removing non-existent template
	err := registry.Remove("non-existent")
	if err == nil {
		t.Error("Expected error for removing non-existent template, got nil")
	}

	// Register template first
	_ = registry.Register(template)

	// Verify template exists
	if !registry.Exists("test-template") {
		t.Error("Template should exist before removal")
	}

	// Test removing existing template
	err = registry.Remove("test-template")
	if err != nil {
		t.Errorf("Expected no error removing existing template, got %v", err)
	}

	// Verify template no longer exists
	if registry.Exists("test-template") {
		t.Error("Template should not exist after removal")
	}

	// Verify template can't be retrieved
	_, err = registry.Get("test-template")
	if err == nil {
		t.Error("Expected error getting removed template, got nil")
	}
}

// RED PHASE: Write tests for GetTemplateTypes function (0% coverage)
func TestRegistry_GetTemplateTypes(t *testing.T) {
	setupTestTemplates(t)
	registry := NewRegistry()

	// Fresh registry should have pre-loaded types
	templateTypes := registry.GetTemplateTypes()
	if len(templateTypes) == 0 {
		t.Error("Expected some types from pre-loaded templates, got 0")
	}

	// Add templates with different types
	template1 := types.Template{ID: "custom1", Name: "Custom 1", Type: "custom"}
	template2 := types.Template{ID: "web1", Name: "Web 1", Type: "web"}
	template3 := types.Template{ID: "custom2", Name: "Custom 2", Type: "custom"}

	_ = registry.Register(template1)
	_ = registry.Register(template2)
	_ = registry.Register(template3)

	allTypes := registry.GetTemplateTypes()

	// Verify types include our added types
	hasCustom := false
	hasWeb := false
	for _, templateType := range allTypes {
		if templateType == "custom" {
			hasCustom = true
		}
		if templateType == "web" {
			hasWeb = true
		}
	}

	if !hasCustom {
		t.Error("Expected 'custom' type to be present")
	}
	if !hasWeb {
		t.Error("Expected 'web' type to be present")
	}

	// Each type should appear only once (deduplication test)
	typeCount := make(map[string]int)
	for _, templateType := range allTypes {
		typeCount[templateType]++
	}

	for templateType, count := range typeCount {
		if count > 1 {
			t.Errorf("Type '%s' appears %d times, should appear only once", templateType, count)
		}
	}
}

// REFACTOR PHASE: Test concurrent access to registry (additional coverage)
func TestRegistry_ConcurrentAccess(t *testing.T) {
	setupTestTemplates(t)
	registry := NewRegistry()

	// Test concurrent registration and removal
	done := make(chan bool, 2)

	// Goroutine 1: Register templates
	go func() {
		for i := 0; i < 10; i++ {
			template := types.Template{
				ID:   fmt.Sprintf("concurrent-%d", i),
				Name: fmt.Sprintf("Concurrent Template %d", i),
				Type: "concurrent",
			}
			_ = registry.Register(template)
		}
		done <- true
	}()

	// Goroutine 2: List and check existence
	go func() {
		for i := 0; i < 10; i++ {
			_ = registry.List()
			_ = registry.Exists(fmt.Sprintf("concurrent-%d", i))
			_ = registry.GetByType("concurrent")
			_ = registry.GetTemplateTypes()
		}
		done <- true
	}()

	// Wait for both goroutines to complete
	<-done
	<-done

	// Verify all templates were registered
	concurrentTemplates := registry.GetByType("concurrent")
	if len(concurrentTemplates) != 10 {
		t.Errorf("Expected 10 concurrent templates, got %d", len(concurrentTemplates))
	}
}
