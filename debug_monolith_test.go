package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

func TestDebugMonolithGeneration(t *testing.T) {
	// Initialize the templates filesystem for testing
	templates.SetTemplatesFS(os.DirFS("blueprints"))
	// Create configuration for monolith
	config := types.ProjectConfig{
		Name:         "debug-monolith",
		Module:       "github.com/test/debug-monolith", 
		Type:         "monolith",
		Architecture: "", // Should default to standard
		GoVersion:    "1.23",
		Framework:    "gin",
		Logger:       "slog",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: "postgres",
				ORM:    "gorm",
			},
			Authentication: types.AuthConfig{
				Type:      "session",
				Providers: []string{}, // Empty slice for session auth
			},
		},
		Variables: map[string]string{
			"AssetPipeline":  "embedded",
			"TemplateEngine": "html/template", 
			"SessionStore":   "cookie",
			"DatabaseDriver": "postgres",
			"DatabaseORM":    "gorm",
			"AuthType":       "session",
			"LoggerType":     "slog",
			"AuthProviders":  "", // Empty for session auth
		},
	}

	// Create generator and test
	gen := generator.New()
	
	options := types.GenerationOptions{
		OutputPath: "/tmp/debug-monolith-test",
		DryRun:     false,
		NoGit:      true,
		Verbose:    true,
	}

	result, err := gen.Generate(config, options)
	if err != nil {
		log.Printf("Generation failed: %v", err)
		fmt.Printf("Full error details: %+v\n", err)
		t.Fatalf("Generation failed: %v", err)
	}

	fmt.Printf("Generation successful! Files created: %d\n", len(result.FilesCreated))
	for _, file := range result.FilesCreated {
		fmt.Printf("  - %s\n", file)
	}
}