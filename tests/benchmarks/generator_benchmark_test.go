package benchmarks

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// setupBenchmarkTemplates sets up templates for benchmark tests
func setupBenchmarkTemplates(b *testing.B) {
	b.Helper()

	// Get the project root for tests
	_, file, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(file)))
	templatesDir := filepath.Join(projectRoot, "templates")

	// Verify templates directory exists
	if _, err := os.Stat(templatesDir); os.IsNotExist(err) {
		b.Fatalf("Templates directory not found at: %s", templatesDir)
	}

	// Set up the filesystem for tests using os.DirFS
	templates.SetTemplatesFS(os.DirFS(templatesDir))
}

// BenchmarkGenerator_GenerateWebAPI benchmarks the generation of a web-api project.
func BenchmarkGenerator_GenerateWebAPI(b *testing.B) {
	// Initialize templates
	setupBenchmarkTemplates(b)
	
	config := types.ProjectConfig{
		Name:      "benchmark-web-api",
		Module:    "github.com/benchmark/web-api",
		Type:      "web-api",
		GoVersion: "1.21",
		Framework: "gin",
		Logger:    "slog",
		Features: &types.Features{
			Database:       types.DatabaseConfig{Driver: "postgres", ORM: "gorm"},
			Authentication: types.AuthConfig{Type: "jwt"},
		},
	}

	gen := generator.New()

	for b.Loop() {
		outputDir := helpers.CreateTempDir(b)
		projectPath := filepath.Join(outputDir, config.Name)

		options := types.GenerationOptions{
			OutputPath: projectPath,
			DryRun:     false,
			NoGit:      true,
			Verbose:    false,
		}

		_, err := gen.Generate(config, options)
		if err != nil {
			b.Fatal(err)
		}
	}
}
