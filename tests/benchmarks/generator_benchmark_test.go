package benchmarks

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// BenchmarkGenerator_GenerateWebAPI benchmarks the generation of a web-api project.
func BenchmarkGenerator_GenerateWebAPI(b *testing.B) {
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
