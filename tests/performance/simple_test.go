package performance

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// TestSimplePerformance runs basic performance validation
func TestSimplePerformance(t *testing.T) {
	// Initialize templates
	blueprintsPath := findBlueprintsPath()
	if blueprintsPath == "" {
		t.Skip("Blueprints directory not found - skipping performance tests")
	}
	templates.SetTemplatesFS(os.DirFS(blueprintsPath))

	// Test simple CLI generation
	config := types.ProjectConfig{
		Name:      "test-simple-perf",
		Module:    "github.com/test/simple-perf",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
		Variables: map[string]string{
			"blueprint_id": "cli-simple",
		},
		Features: &types.Features{}, // Initialize Features to avoid nil pointer
	}

	tempDir, err := os.MkdirTemp("", "simple-perf-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, config.Name)
	
	// Measure generation time
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)
	
	startTime := time.Now()
	
	gen := generator.New()
	result, err := gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
	})
	
	duration := time.Since(startTime)
	
	runtime.GC()
	runtime.ReadMemStats(&memEnd)
	
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Validate performance goals
	t.Logf("Generation completed in %v", duration)
	t.Logf("Files created: %d", len(result.FilesCreated))
	
	memoryUsedMB := float64(memEnd.TotalAlloc-memStart.TotalAlloc) / 1024 / 1024
	t.Logf("Memory used: %.2f MB", memoryUsedMB)

	// Check performance targets
	if duration > 5*time.Second {
		t.Errorf("Generation took too long: %v (target: <5s for basic test)", duration)
	}
	
	if memoryUsedMB > 100 {
		t.Errorf("Memory usage too high: %.2f MB (target: <100MB for basic test)", memoryUsedMB)
	}
	
	if len(result.FilesCreated) == 0 {
		t.Error("No files were generated")
	}

	// Check if files exist and are readable
	for _, filePath := range result.FilesCreated {
		if _, err := os.Stat(filePath); err != nil {
			t.Errorf("Generated file does not exist or is not readable: %s", filePath)
		}
	}

	t.Logf("✅ Simple performance test passed: %v duration, %.2f MB memory, %d files", 
		duration, memoryUsedMB, len(result.FilesCreated))
}

// TestMemoryGeneration tests in-memory generation performance
func TestMemoryGeneration(t *testing.T) {
	blueprintsPath := findBlueprintsPath()
	if blueprintsPath == "" {
		t.Skip("Blueprints directory not found - skipping memory generation test")
	}
	templates.SetTemplatesFS(os.DirFS(blueprintsPath))

	config := types.ProjectConfig{
		Name:      "test-memory-gen",
		Module:    "github.com/test/memory-gen",
		Type:      "cli",
		GoVersion: "1.21",
		Framework: "cobra",
		Logger:    "slog",
		Variables: map[string]string{
			"blueprint_id": "cli-simple",
		},
		Features: &types.Features{},
	}

	gen := generator.New()
	
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)
	
	startTime := time.Now()
	
	files, err := gen.GenerateInMemory(&config, "cli-simple")
	
	duration := time.Since(startTime)
	
	runtime.GC()
	runtime.ReadMemStats(&memEnd)
	
	if err != nil {
		t.Fatalf("Memory generation failed: %v", err)
	}

	memoryUsedMB := float64(memEnd.TotalAlloc-memStart.TotalAlloc) / 1024 / 1024
	
	t.Logf("Memory generation completed in %v", duration)
	t.Logf("Files generated: %d", len(files))
	t.Logf("Memory used: %.2f MB", memoryUsedMB)
	
	// Memory generation should be faster than file generation
	if duration > 2*time.Second {
		t.Errorf("Memory generation took too long: %v (target: <2s)", duration)
	}
	
	if len(files) == 0 {
		t.Error("No files were generated in memory")
	}
	
	// Validate file content is not empty
	for filePath, content := range files {
		if len(content) == 0 {
			t.Errorf("Generated file is empty: %s", filePath)
		}
	}

	t.Logf("✅ Memory generation test passed: %v duration, %.2f MB memory, %d files", 
		duration, memoryUsedMB, len(files))
}

// TestConcurrentGeneration tests concurrent generation performance
func TestConcurrentGeneration(t *testing.T) {
	blueprintsPath := findBlueprintsPath()
	if blueprintsPath == "" {
		t.Skip("Blueprints directory not found - skipping concurrent generation test")
	}
	templates.SetTemplatesFS(os.DirFS(blueprintsPath))

	// Test concurrent generation of multiple projects
	numProjects := 3
	results := make(chan error, numProjects)
	
	startTime := time.Now()
	
	for i := 0; i < numProjects; i++ {
		go func(index int) {
			config := types.ProjectConfig{
				Name:      fmt.Sprintf("test-concurrent-%d", index),
				Module:    fmt.Sprintf("github.com/test/concurrent-%d", index),
				Type:      "cli",
				GoVersion: "1.21",
				Framework: "cobra",
				Logger:    "slog",
				Variables: map[string]string{
					"blueprint_id": "cli-simple",
				},
				Features: &types.Features{},
			}

			tempDir, err := os.MkdirTemp("", fmt.Sprintf("concurrent-%d-*", index))
			if err != nil {
				results <- err
				return
			}
			defer os.RemoveAll(tempDir)

			projectPath := filepath.Join(tempDir, config.Name)
			
			gen := generator.New()
			_, err = gen.Generate(config, types.GenerationOptions{
				OutputPath: projectPath,
				DryRun:     false,
				NoGit:      true,
			})
			
			results <- err
		}(i)
	}
	
	// Wait for all goroutines to complete
	var errors []error
	for i := 0; i < numProjects; i++ {
		if err := <-results; err != nil {
			errors = append(errors, err)
		}
	}
	
	duration := time.Since(startTime)
	
	if len(errors) > 0 {
		t.Errorf("Concurrent generation failed with %d errors: %v", len(errors), errors[0])
	}
	
	t.Logf("Concurrent generation of %d projects completed in %v", numProjects, duration)
	
	// Concurrent generation should not be significantly slower than sequential
	if duration > 10*time.Second {
		t.Errorf("Concurrent generation took too long: %v (target: <10s for %d projects)", 
			duration, numProjects)
	}

	t.Logf("✅ Concurrent generation test passed: %v duration for %d projects", 
		duration, numProjects)
}

