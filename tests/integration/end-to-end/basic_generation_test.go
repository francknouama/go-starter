package endtoend_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// setupTemplatesForTesting initializes the templates filesystem for integration tests
func setupTemplatesForTesting(t *testing.T) {
	t.Helper()
	templatesDir := "../../blueprints"
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		t.Fatalf("Blueprints directory not found: %s", templatesDir)
	}
	templates.SetTemplatesFS(os.DirFS(templatesDir))
}

// TestIntegration_BasicProjectGeneration tests the generation of a basic web-api project.
func TestIntegration_BasicProjectGeneration(t *testing.T) {
	setupTemplatesForTesting(t)

	config := types.ProjectConfig{
		Name:      "test-basic-api",
		Module:    "github.com/test/test-basic-api",
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
	}

	projectPath := helpers.GenerateProject(t, config)

	// Debug: List generated files
	if files, err := filepath.Glob(filepath.Join(projectPath, "*")); err == nil {
		t.Logf("Generated files in root: %v", files)
	}
	if files, err := filepath.Glob(filepath.Join(projectPath, "cmd/server/*")); err == nil {
		t.Logf("Generated files in cmd/server: %v", files)
	}

	// Assertions for basic project structure
	helpers.AssertProjectGenerated(t, projectPath, []string{
		"go.mod",
		"README.md",
		"cmd/server/main.go",
		"internal/handlers/health.go",
		"internal/config/config.go",
	})

	// Assert go.mod content
	helpers.AssertGoModValid(t, filepath.Join(projectPath, "go.mod"), config.Module)
	helpers.AssertGoModContainsVersion(t, filepath.Join(projectPath, "go.mod"), config.GoVersion)

	// Assert project compiles
	helpers.AssertCompilable(t, projectPath)
}
