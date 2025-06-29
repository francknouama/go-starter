package templates

import (
	"testing"
	"testing/fstest"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRegistryComprehensive tests registry operations comprehensively
func TestRegistryComprehensive(t *testing.T) {
	t.Run("NewRegistry_Initialization", func(t *testing.T) {
		// Set up mock filesystem for template loading
		SetTemplatesFS(fstest.MapFS{
			"web-api/template.yaml": &fstest.MapFile{
				Data: []byte(`
name: "web-api"
type: "web-api"
architecture: "standard"
files:
  - source: "main.go.tmpl"
    destination: "main.go"
`),
			},
		})

		registry := NewRegistry()

		assert.NotNil(t, registry)
		assert.NotNil(t, registry.templates)

		// Check that templates were loaded during initialization
		templates := registry.List()
		assert.GreaterOrEqual(t, len(templates), 1)
	})

	t.Run("Register_Success", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		template := types.Template{
			ID:   "test-template",
			Name: "Test Template",
			Type: "test",
		}

		err := registry.Register(template)

		require.NoError(t, err)

		// Verify template was registered
		retrieved, err := registry.Get("test-template")
		require.NoError(t, err)
		assert.Equal(t, "test-template", retrieved.ID)
		assert.Equal(t, "Test Template", retrieved.Name)
	})

	t.Run("Register_EmptyID", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		template := types.Template{
			ID:   "", // Empty ID should cause error
			Name: "Test Template",
			Type: "test",
		}

		err := registry.Register(template)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template ID cannot be empty")
	})

	t.Run("Get_Success", func(t *testing.T) {
		registry := &Registry{
			templates: map[string]types.Template{
				"existing": {
					ID:   "existing",
					Name: "Existing Template",
					Type: "test",
				},
			},
		}

		template, err := registry.Get("existing")

		require.NoError(t, err)
		assert.Equal(t, "existing", template.ID)
		assert.Equal(t, "Existing Template", template.Name)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		_, err := registry.Get("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template not found: nonexistent")
	})

	t.Run("List_MultipleTemplates", func(t *testing.T) {
		templates := map[string]types.Template{
			"web-api": {
				ID:   "web-api",
				Name: "Web API",
				Type: "web-api",
			},
			"cli": {
				ID:   "cli",
				Name: "CLI App",
				Type: "cli",
			},
		}

		registry := &Registry{templates: templates}
		result := registry.List()

		assert.Len(t, result, 2)

		// Check that both templates are present
		templateIDs := make(map[string]bool)
		for _, tmpl := range result {
			templateIDs[tmpl.ID] = true
		}

		assert.True(t, templateIDs["web-api"])
		assert.True(t, templateIDs["cli"])
	})

	t.Run("GetByType_Success", func(t *testing.T) {
		templates := map[string]types.Template{
			"web-api-standard": {
				ID:   "web-api-standard",
				Type: "web-api",
			},
			"web-api-clean": {
				ID:   "web-api-clean",
				Type: "web-api",
			},
			"cli": {
				ID:   "cli",
				Type: "cli",
			},
		}

		registry := &Registry{templates: templates}
		webAPITemplates := registry.GetByType("web-api")

		assert.Len(t, webAPITemplates, 2)

		// Verify only web-api templates are returned
		for _, tmpl := range webAPITemplates {
			assert.Equal(t, "web-api", tmpl.Type)
		}
	})

	t.Run("GetByType_NoMatches", func(t *testing.T) {
		registry := &Registry{
			templates: map[string]types.Template{
				"cli": {
					ID:   "cli",
					Type: "cli",
				},
			},
		}

		result := registry.GetByType("web-api")

		assert.Len(t, result, 0)
	})

	t.Run("Exists_True", func(t *testing.T) {
		registry := &Registry{
			templates: map[string]types.Template{
				"existing": {ID: "existing"},
			},
		}

		exists := registry.Exists("existing")
		assert.True(t, exists)
	})

	t.Run("Exists_False", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		exists := registry.Exists("nonexistent")
		assert.False(t, exists)
	})

	t.Run("Remove_Success", func(t *testing.T) {
		registry := &Registry{
			templates: map[string]types.Template{
				"to-remove": {ID: "to-remove"},
			},
		}

		err := registry.Remove("to-remove")

		require.NoError(t, err)
		assert.False(t, registry.Exists("to-remove"))
	})

	t.Run("Remove_NotFound", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		err := registry.Remove("nonexistent")

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "template not found: nonexistent")
	})

	t.Run("GetTemplateTypes_MultipleTypes", func(t *testing.T) {
		templates := map[string]types.Template{
			"web-api-1": {Type: "web-api"},
			"web-api-2": {Type: "web-api"},
			"cli":       {Type: "cli"},
			"library":   {Type: "library"},
		}

		registry := &Registry{templates: templates}
		types := registry.GetTemplateTypes()

		assert.Len(t, types, 3) // web-api, cli, library

		typeSet := make(map[string]bool)
		for _, t := range types {
			typeSet[t] = true
		}

		assert.True(t, typeSet["web-api"])
		assert.True(t, typeSet["cli"])
		assert.True(t, typeSet["library"])
	})

	t.Run("GetTemplateTypes_Empty", func(t *testing.T) {
		registry := &Registry{
			templates: make(map[string]types.Template),
		}

		types := registry.GetTemplateTypes()
		assert.Len(t, types, 0)
	})
}

// TestRegistryArchitectureSupport tests architecture-specific template handling
func TestRegistryArchitectureSupport(t *testing.T) {
	t.Run("ArchitectureTemplates_Separation", func(t *testing.T) {
		templates := map[string]types.Template{
			"web-api": {
				ID:           "web-api",
				Type:         "web-api",
				Architecture: "standard",
			},
			"web-api-clean": {
				ID:           "web-api-clean",
				Type:         "web-api",
				Architecture: "clean",
			},
			"web-api-ddd": {
				ID:           "web-api-ddd",
				Type:         "web-api",
				Architecture: "ddd",
			},
		}

		registry := &Registry{templates: templates}

		// Verify we can get all web-api templates
		webAPITemplates := registry.GetByType("web-api")
		assert.Len(t, webAPITemplates, 3)

		// Verify each specific architecture template exists
		standardTemplate, err := registry.Get("web-api")
		require.NoError(t, err)
		assert.Equal(t, "standard", standardTemplate.Architecture)

		cleanTemplate, err := registry.Get("web-api-clean")
		require.NoError(t, err)
		assert.Equal(t, "clean", cleanTemplate.Architecture)

		dddTemplate, err := registry.Get("web-api-ddd")
		require.NoError(t, err)
		assert.Equal(t, "ddd", dddTemplate.Architecture)
	})
}

// TestRegistryConcurrency tests thread safety of registry operations
func TestRegistryConcurrency(t *testing.T) {
	t.Run("ConcurrentReadWrite", func(t *testing.T) {
		registry := NewRegistry()

		// Test concurrent reads and writes
		done := make(chan bool, 4)

		// Concurrent readers
		go func() {
			for i := 0; i < 100; i++ {
				registry.List()
				registry.GetTemplateTypes()
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 100; i++ {
				registry.Exists("web-api")
				registry.GetByType("web-api")
			}
			done <- true
		}()

		// Concurrent writers (using different IDs to avoid conflicts)
		go func() {
			for i := 0; i < 50; i++ {
				template := types.Template{
					ID:   "test-template-1",
					Name: "Test Template 1",
					Type: "test",
				}
				registry.Register(template)
				registry.Remove("test-template-1")
			}
			done <- true
		}()

		go func() {
			for i := 0; i < 50; i++ {
				template := types.Template{
					ID:   "test-template-2",
					Name: "Test Template 2",
					Type: "test",
				}
				registry.Register(template)
				registry.Remove("test-template-2")
			}
			done <- true
		}()

		// Wait for all goroutines to complete
		for i := 0; i < 4; i++ {
			<-done
		}

		// Registry should still be functional
		templates := registry.List()
		assert.NotNil(t, templates)
	})
}
