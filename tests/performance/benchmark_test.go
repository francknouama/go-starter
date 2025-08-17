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


// Blueprint configurations for testing different project types
var blueprintConfigs = []struct {
	name   string
	config types.ProjectConfig
}{
	{
		name: "cli-simple",
		config: types.ProjectConfig{
			Name:      "test-cli-simple",
			Module:    "github.com/test/cli-simple",
			Type:      "cli",
			GoVersion: "1.21",
			Framework: "cobra",
			Logger:    "slog",
			Variables: map[string]string{
				"blueprint_id": "cli-simple",
			},
			Features: &types.Features{},
		},
	},
	{
		name: "cli-standard",
		config: types.ProjectConfig{
			Name:      "test-cli-standard",
			Module:    "github.com/test/cli-standard",
			Type:      "cli",
			GoVersion: "1.21",
			Framework: "cobra",
			Logger:    "slog",
			Features:  &types.Features{},
		},
	},
	{
		name: "web-api-standard",
		config: types.ProjectConfig{
			Name:         "test-web-api",
			Module:       "github.com/test/web-api",
			Type:         "web-api",
			GoVersion:    "1.21",
			Framework:    "gin",
			Architecture: "standard",
			Logger:       "slog",
			Features:     &types.Features{},
		},
	},
	{
		name: "web-api-clean",
		config: types.ProjectConfig{
			Name:         "test-web-api-clean",
			Module:       "github.com/test/web-api-clean",
			Type:         "web-api",
			GoVersion:    "1.21",
			Framework:    "gin",
			Architecture: "clean",
			Logger:       "zap",
			Features:     &types.Features{},
		},
	},
	{
		name: "lambda-standard",
		config: types.ProjectConfig{
			Name:      "test-lambda",
			Module:    "github.com/test/lambda",
			Type:      "lambda",
			GoVersion: "1.21",
			Logger:    "slog",
			Features:  &types.Features{},
		},
	},
	{
		name: "library-standard",
		config: types.ProjectConfig{
			Name:      "test-library",
			Module:    "github.com/test/library",
			Type:      "library",
			GoVersion: "1.21",
			Logger:    "slog",
			Features:  &types.Features{},
		},
	},
}

// Initialize the test environment
func init() {
	// Set up embedded templates for testing
	blueprintsPath := findBlueprintsPath()
	if blueprintsPath != "" {
		templates.SetTemplatesFS(os.DirFS(blueprintsPath))
	}
}

// findBlueprintsPath locates the blueprints directory
func findBlueprintsPath() string {
	// Start from current directory and work up to find blueprints
	currentDir, _ := os.Getwd()
	for currentDir != "/" && currentDir != "" {
		blueprintsPath := filepath.Join(currentDir, "blueprints")
		if _, err := os.Stat(blueprintsPath); err == nil {
			return blueprintsPath
		}
		currentDir = filepath.Dir(currentDir)
	}
	return ""
}

// BenchmarkTemplateGeneration benchmarks all blueprint types
func BenchmarkTemplateGeneration(b *testing.B) {
	for _, blueprint := range blueprintConfigs {
		b.Run(blueprint.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				benchmarkSingleGeneration(b, blueprint.name, blueprint.config)
			}
		})
	}
}

// BenchmarkParallelGeneration tests concurrent template generation
func BenchmarkParallelGeneration(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		blueprintIndex := 0
		for pb.Next() {
			blueprint := blueprintConfigs[blueprintIndex%len(blueprintConfigs)]
			config := blueprint.config
			// Make unique names for parallel runs
			config.Name = fmt.Sprintf("%s-parallel-%d", blueprint.config.Name, blueprintIndex)
			config.Module = fmt.Sprintf("%s-parallel-%d", blueprint.config.Module, blueprintIndex)
			
			benchmarkSingleGeneration(b, blueprint.name, config)
			blueprintIndex++
		}
	})
}

// BenchmarkMemoryGeneration benchmarks in-memory generation (no file I/O)
func BenchmarkMemoryGeneration(b *testing.B) {
	gen := generator.New()
	
	for _, blueprint := range blueprintConfigs {
		b.Run(fmt.Sprintf("%s-memory", blueprint.name), func(b *testing.B) {
			var memStart, memEnd runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&memStart)
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				config := blueprint.config
				config.Name = fmt.Sprintf("%s-mem-%d", config.Name, i)
				config.Module = fmt.Sprintf("%s-mem-%d", config.Module, i)
				
				blueprintID := getBlueprintID(config)
				_, err := gen.GenerateInMemory(&config, blueprintID)
				if err != nil {
					b.Errorf("Memory generation failed for %s: %v", blueprint.name, err)
				}
			}
			b.StopTimer()
			
			runtime.GC()
			runtime.ReadMemStats(&memEnd)
			
			memUsed := int64(memEnd.TotalAlloc - memStart.TotalAlloc)
			recordResult(BenchmarkResult{
				Blueprint:   fmt.Sprintf("%s-memory", blueprint.name),
				Duration:    time.Duration(b.Elapsed().Nanoseconds() / int64(b.N)),
				MemoryUsed:  memUsed / int64(b.N),
				Platform:    runtime.GOOS,
				GoVersion:   runtime.Version(),
				Timestamp:   time.Now(),
				Success:     true,
			})
		})
	}
}

// BenchmarkFileIO specifically tests file I/O performance
func BenchmarkFileIO(b *testing.B) {
	for _, blueprint := range blueprintConfigs {
		b.Run(fmt.Sprintf("%s-fileio", blueprint.name), func(b *testing.B) {
			gen := generator.New()
			config := blueprint.config
			blueprintID := getBlueprintID(config)
			
			// Generate in memory first to get file content
			files, err := gen.GenerateInMemory(&config, blueprintID)
			if err != nil {
				b.Fatalf("Failed to generate in memory: %v", err)
			}
			
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				tempDir, err := os.MkdirTemp("", fmt.Sprintf("benchmark-fileio-%s-*", blueprint.name))
				if err != nil {
					b.Fatalf("Failed to create temp dir: %v", err)
				}
				
				// Measure file writing time
				startTime := time.Now()
				for filePath, content := range files {
					fullPath := filepath.Join(tempDir, filePath)
					if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
						b.Errorf("Failed to create directory: %v", err)
						continue
					}
					if err := os.WriteFile(fullPath, content, 0644); err != nil {
						b.Errorf("Failed to write file %s: %v", filePath, err)
					}
				}
				fileIODuration := time.Since(startTime)
				
				// Clean up
				os.RemoveAll(tempDir)
				
				recordResult(BenchmarkResult{
					Blueprint:    fmt.Sprintf("%s-fileio", blueprint.name),
					Duration:     fileIODuration,
					FilesCreated: len(files),
					Platform:     runtime.GOOS,
					GoVersion:    runtime.Version(),
					Timestamp:    time.Now(),
					Success:      true,
				})
			}
		})
	}
}

// benchmarkSingleGeneration performs a complete generation benchmark
func benchmarkSingleGeneration(b *testing.B, blueprintName string, config types.ProjectConfig) {
	var memStart, memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)
	
	gen := generator.New()
	
	// Create unique temporary directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("benchmark-%s-*", blueprintName))
	if err != nil {
		b.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	projectPath := filepath.Join(tempDir, config.Name)
	
	startTime := time.Now()
	
	result, err := gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true, // Skip git for benchmarking
	})
	
	duration := time.Since(startTime)
	
	runtime.GC()
	runtime.ReadMemStats(&memEnd)
	
	// Calculate directory size
	dirSize := calculateDirectorySize(projectPath)
	
	benchResult := BenchmarkResult{
		Blueprint:   blueprintName,
		Duration:    duration,
		FilesCreated: len(result.FilesCreated),
		MemoryUsed:  int64(memEnd.TotalAlloc - memStart.TotalAlloc),
		TempDirSize: dirSize,
		Platform:    runtime.GOOS,
		GoVersion:   runtime.Version(),
		Timestamp:   time.Now(),
		Success:     err == nil,
	}
	
	if err != nil {
		benchResult.ErrorMessage = err.Error()
		b.Errorf("Generation failed for %s: %v", blueprintName, err)
	}
	
	recordResult(benchResult)
}


// calculateDirectorySize calculates the total size of a directory
func calculateDirectorySize(dirPath string) int64 {
	var size int64
	filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Ignore errors for size calculation
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})
	return size
}

// recordResult safely adds a benchmark result to the suite
func recordResult(result BenchmarkResult) {
	resultsMutex.Lock()
	defer resultsMutex.Unlock()
	benchmarkSuite.Results = append(benchmarkSuite.Results, result)
}

// TestGeneratePerformanceReport generates a comprehensive performance report
func TestGeneratePerformanceReport(t *testing.T) {
	// Run a subset of benchmarks for the report
	runPerformanceTests(t)
	
	// Calculate summary statistics
	calculateSummary()
	
	// Generate reports
	if err := generateJSONReport(); err != nil {
		t.Errorf("Failed to generate JSON report: %v", err)
	}
	
	if err := generateMarkdownReport(); err != nil {
		t.Errorf("Failed to generate Markdown report: %v", err)
	}
	
	// Validate performance goals
	validatePerformanceGoals(t)
}

// runPerformanceTests executes a comprehensive set of performance tests
func runPerformanceTests(t *testing.T) {
	for _, blueprint := range blueprintConfigs {
		t.Run(fmt.Sprintf("Performance_%s", blueprint.name), func(t *testing.T) {
			// Run multiple iterations for statistical significance
			for i := 0; i < 3; i++ {
				config := blueprint.config
				config.Name = fmt.Sprintf("%s-perf-%d", config.Name, i)
				config.Module = fmt.Sprintf("%s-perf-%d", config.Module, i)
				
				benchmarkSingleGeneration(&testing.B{}, blueprint.name, config)
			}
		})
	}
}

// calculateSummary computes aggregate statistics
func calculateSummary() {
	resultsMutex.Lock()
	defer resultsMutex.Unlock()
	
	if len(benchmarkSuite.Results) == 0 {
		return
	}
	
	var totalDuration time.Duration
	var minDuration = time.Hour // Start with a large value
	var maxDuration time.Duration
	var totalFiles int
	var totalMemory int64
	var successCount int
	
	for _, result := range benchmarkSuite.Results {
		if result.Success {
			successCount++
			totalDuration += result.Duration
			totalFiles += result.FilesCreated
			totalMemory += result.MemoryUsed
			
			if result.Duration < minDuration {
				minDuration = result.Duration
			}
			if result.Duration > maxDuration {
				maxDuration = result.Duration
			}
		}
	}
	
	benchmarkSuite.Summary = BenchmarkSummary{
		TotalRuns:         len(benchmarkSuite.Results),
		SuccessfulRuns:    successCount,
		FailedRuns:        len(benchmarkSuite.Results) - successCount,
		TotalFilesCreated: totalFiles,
		TotalMemoryUsed:   totalMemory,
		MinDuration:       minDuration,
		MaxDuration:       maxDuration,
	}
	
	if successCount > 0 {
		benchmarkSuite.Summary.AverageDuration = totalDuration / time.Duration(successCount)
	}
	
	// Calculate performance grade
	avgMs := benchmarkSuite.Summary.AverageDuration.Milliseconds()
	switch {
	case avgMs < 500:
		benchmarkSuite.Summary.PerformanceGrade = "A+ (Excellent)"
	case avgMs < 1000:
		benchmarkSuite.Summary.PerformanceGrade = "A (Very Good)"
	case avgMs < 2000:
		benchmarkSuite.Summary.PerformanceGrade = "B (Good)"
	case avgMs < 5000:
		benchmarkSuite.Summary.PerformanceGrade = "C (Acceptable)"
	default:
		benchmarkSuite.Summary.PerformanceGrade = "D (Needs Improvement)"
	}
}

// validatePerformanceGoals checks if performance targets are met
func validatePerformanceGoals(t *testing.T) {
	summary := benchmarkSuite.Summary
	
	// Goal: <2s generation time
	if summary.AverageDuration > 2*time.Second {
		t.Errorf("Performance goal not met: Average generation time %.2fs exceeds 2s target", 
			summary.AverageDuration.Seconds())
	}
	
	// Goal: 95% success rate
	successRate := float64(summary.SuccessfulRuns) / float64(summary.TotalRuns) * 100
	if successRate < 95.0 {
		t.Errorf("Reliability goal not met: Success rate %.1f%% is below 95%% target", successRate)
	}
	
	// Goal: Memory usage under 50MB per generation
	avgMemoryMB := float64(summary.TotalMemoryUsed) / float64(summary.SuccessfulRuns) / 1024 / 1024
	if avgMemoryMB > 50.0 {
		t.Errorf("Memory goal not met: Average memory usage %.1fMB exceeds 50MB target", avgMemoryMB)
	}
	
	t.Logf("Performance Summary:")
	t.Logf("  Average Duration: %v", summary.AverageDuration)
	t.Logf("  Success Rate: %.1f%%", successRate)
	t.Logf("  Average Memory: %.1fMB", avgMemoryMB)
	t.Logf("  Performance Grade: %s", summary.PerformanceGrade)
}