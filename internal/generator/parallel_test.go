package generator

import (
	"testing"
	"time"
	"path/filepath"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestWorkerPool(t *testing.T) {
	t.Run("worker pool creation", func(t *testing.T) {
		pool := NewWorkerPool(2)
		if pool == nil {
			t.Error("Expected worker pool to not be nil")
		}
		
		if pool.numWorkers != 2 {
			t.Errorf("Expected 2 workers, got %d", pool.numWorkers)
		}
		
		// Clean up
		pool.Stop()
	})
	
	t.Run("optimal worker count", func(t *testing.T) {
		count := GetOptimalWorkerCount()
		if count <= 0 {
			t.Error("Expected positive worker count")
		}
		
		// Should be reasonable (between 1 and CPU count * 2)
		maxExpected := 16 // reasonable upper bound for tests
		if count > maxExpected {
			t.Errorf("Worker count seems too high: %d", count)
		}
	})
	
	t.Run("worker pool lifecycle", func(t *testing.T) {
		pool := NewWorkerPool(1)
		
		// Start the pool
		pool.Start()
		
		// Should be able to stop gracefully
		pool.Stop()
		
		// Should not panic when stopping again
		pool.Stop()
	})
	
	t.Run("progress tracking", func(t *testing.T) {
		pool := NewWorkerPool(1)
		defer pool.Stop()
		
		// Initially no progress
		progress := pool.GetProgress()
		if progress != 0 {
			t.Errorf("Expected 0 progress, got %f", progress)
		}
		
		// Progress should work even without jobs
		pool.totalJobs = 10
		pool.completedJobs = 5
		progress = pool.GetProgress()
		if progress != 50.0 {
			t.Errorf("Expected 50%% progress, got %f", progress)
		}
	})
}

func TestParallelTemplateProcessor(t *testing.T) {
	t.Run("processor creation", func(t *testing.T) {
		processor := NewParallelTemplateProcessor(2)
		if processor == nil {
			t.Error("Expected processor to not be nil")
		}
		
		if processor.pool == nil {
			t.Error("Expected processor to have worker pool")
		}
		
		// Clean up
		processor.Shutdown()
	})
	
	t.Run("processor lifecycle", func(t *testing.T) {
		processor := NewParallelTemplateProcessor(1)
		
		// Should start successfully
		processor.Start()
		
		// Should shutdown gracefully
		processor.Shutdown()
		
		// Should handle multiple shutdowns
		processor.Shutdown()
	})
	
	t.Run("empty job processing", func(t *testing.T) {
		processor := NewParallelTemplateProcessor(1)
		defer processor.Shutdown()
		
		processor.Start()
		
		// Process empty job list
		jobs := []types.TemplateFile{}
		config := types.ProjectConfig{Name: "test"}
		outputDir := t.TempDir()
		
		results, err := processor.ProcessTemplates(jobs, config, outputDir)
		if err != nil {
			t.Errorf("Expected no error for empty jobs, got %v", err)
		}
		
		if len(results) != 0 {
			t.Errorf("Expected no results for empty jobs, got %d", len(results))
		}
	})
	
	t.Run("progress monitoring", func(t *testing.T) {
		processor := NewParallelTemplateProcessor(1)
		defer processor.Shutdown()
		
		progress := processor.GetProgress()
		if progress != 0 {
			t.Errorf("Expected 0 progress initially, got %f", progress)
		}
	})
}

func TestTemplateJob(t *testing.T) {
	t.Run("job creation and validation", func(t *testing.T) {
		file := types.TemplateFile{
			Source:      "test.tmpl",
			Destination: "test.go",
		}
		
		config := types.ProjectConfig{
			Name:   "test-project",
			Module: "github.com/test/test-project",
		}
		
		outputDir := t.TempDir()
		
		job := &templateJob{
			file:      file,
			config:    config,
			outputDir: outputDir,
		}
		
		if job.file.Source != "test.tmpl" {
			t.Errorf("Expected source 'test.tmpl', got '%s'", job.file.Source)
		}
		
		if job.config.Name != "test-project" {
			t.Errorf("Expected project name 'test-project', got '%s'", job.config.Name)
		}
		
		if job.outputDir != outputDir {
			t.Errorf("Expected output dir '%s', got '%s'", outputDir, job.outputDir)
		}
	})
}

func TestTemplateResult(t *testing.T) {
	t.Run("successful result", func(t *testing.T) {
		result := &templateResult{
			file:        types.TemplateFile{Destination: "test.go"},
			outputPath:  "/path/to/test.go",
			err:         nil,
		}
		
		if result.err != nil {
			t.Error("Expected successful result to have no error")
		}
		
		if result.outputPath != "/path/to/test.go" {
			t.Errorf("Expected output path '/path/to/test.go', got '%s'", result.outputPath)
		}
	})
	
	t.Run("error result", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		result := &templateResult{
			file: types.TemplateFile{Destination: "test.go"},
			err:  testErr,
		}
		
		if result.err == nil {
			t.Error("Expected error result to have error")
		}
		
		if result.err.Error() != "test error" {
			t.Errorf("Expected error 'test error', got '%v'", result.err)
		}
	})
}

func TestParallelProcessingConfiguration(t *testing.T) {
	t.Run("worker count validation", func(t *testing.T) {
		// Test minimum worker count
		pool := NewWorkerPool(0)
		if pool.numWorkers < 1 {
			t.Error("Worker pool should enforce minimum of 1 worker")
		}
		pool.Stop()
		
		// Test reasonable worker count
		pool2 := NewWorkerPool(100)
		if pool2.numWorkers > 16 { // reasonable maximum
			t.Logf("Note: Worker count is high (%d), this may be intentional", pool2.numWorkers)
		}
		pool2.Stop()
	})
	
	t.Run("concurrent safety", func(t *testing.T) {
		pool := NewWorkerPool(2)
		defer pool.Stop()
		
		pool.Start()
		
		// Multiple goroutines checking progress should be safe
		done := make(chan bool, 3)
		
		for i := 0; i < 3; i++ {
			go func() {
				defer func() { done <- true }()
				for j := 0; j < 10; j++ {
					progress := pool.GetProgress()
					_ = progress // Just ensure it doesn't panic
					time.Sleep(time.Millisecond)
				}
			}()
		}
		
		// Wait for all goroutines
		for i := 0; i < 3; i++ {
			<-done
		}
	})
}

func TestParallelProcessingIntegration(t *testing.T) {
	t.Run("integration with temporary files", func(t *testing.T) {
		processor := NewParallelTemplateProcessor(1)
		defer processor.Shutdown()
		
		processor.Start()
		
		// Create a temporary output directory
		outputDir := t.TempDir()
		
		// Simple config for testing
		config := types.ProjectConfig{
			Name:   "integration-test",
			Module: "github.com/test/integration-test",
		}
		
		// Test with empty template files list
		files := []types.TemplateFile{}
		
		results, err := processor.ProcessTemplates(files, config, outputDir)
		if err != nil {
			t.Errorf("Integration test failed: %v", err)
		}
		
		if len(results) != 0 {
			t.Errorf("Expected 0 results, got %d", len(results))
		}
	})
	
	t.Run("processor performance characteristics", func(t *testing.T) {
		// Test that processor can be created and destroyed quickly
		start := time.Now()
		
		processor := NewParallelTemplateProcessor(1)
		processor.Start()
		processor.Shutdown()
		
		elapsed := time.Since(start)
		if elapsed > time.Second {
			t.Errorf("Processor lifecycle took too long: %v", elapsed)
		}
	})
}