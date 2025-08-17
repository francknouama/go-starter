package generator

import (
	"fmt"
	"testing"
	"time"

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
		progress = pool.GetProgress()
		if progress != 0 {
			t.Errorf("Expected 0%% progress with no processed jobs, got %f", progress)
		}
	})
}

func TestParallelTemplateProcessor(t *testing.T) {
	t.Run("processor creation", func(t *testing.T) {
		// Create a mock transaction for testing
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(2, tx)
		if processor == nil {
			t.Error("Expected processor to not be nil")
		}
		
		if processor.pool == nil {
			t.Error("Expected processor to have worker pool")
		}
		
		// Clean up
		processor.pool.Stop()
	})
	
	t.Run("processor lifecycle", func(t *testing.T) {
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(1, tx)
		
		// Should start successfully (no explicit Start method needed)
		// Pool starts when ProcessTemplates is called
		
		// Should shutdown gracefully
		processor.pool.Stop()
		
		// Should handle multiple shutdowns
		processor.pool.Stop()
	})
	
	t.Run("empty job processing", func(t *testing.T) {
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(1, tx)
		defer processor.pool.Stop()
		
		// Process empty job list
		jobs := []types.TemplateFile{}
		config := types.ProjectConfig{Name: "test"}
		outputDir := t.TempDir()
		context := make(map[string]any)
		tmpl := &types.Template{}
		
		results, err := processor.ProcessTemplates(jobs, outputDir, outputDir, context, config, tmpl)
		if err != nil {
			t.Errorf("Expected no error for empty jobs, got %v", err)
		}
		
		if len(results) != 0 {
			t.Errorf("Expected no results for empty jobs, got %d", len(results))
		}
	})
	
	t.Run("progress monitoring", func(t *testing.T) {
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(1, tx)
		defer processor.pool.Stop()
		
		progress := processor.pool.GetProgress()
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
		
		outputDir := t.TempDir()
		destPath := outputDir + "/test.go"
		context := make(map[string]any)
		tx := &GenerationTransaction{}
		
		job := &templateJob{
			templateFile: file,
			templateDir:  "/templates",
			destPath:     destPath,
			context:      context,
			jobID:        1,
			transaction:  tx,
		}
		
		if job.templateFile.Source != "test.tmpl" {
			t.Errorf("Expected source 'test.tmpl', got '%s'", job.templateFile.Source)
		}
		
		if job.destPath != destPath {
			t.Errorf("Expected dest path '%s', got '%s'", destPath, job.destPath)
		}
		
		if job.jobID != 1 {
			t.Errorf("Expected job ID 1, got %d", job.jobID)
		}
	})
}

func TestTemplateResult(t *testing.T) {
	t.Run("successful result", func(t *testing.T) {
		result := &templateResult{
			jobID:      1,
			filePath:   "/path/to/test.go",
			error:      nil,
			wasSkipped: false,
		}
		
		if result.error != nil {
			t.Error("Expected successful result to have no error")
		}
		
		if result.filePath != "/path/to/test.go" {
			t.Errorf("Expected file path '/path/to/test.go', got '%s'", result.filePath)
		}
		
		if result.jobID != 1 {
			t.Errorf("Expected job ID 1, got %d", result.jobID)
		}
	})
	
	t.Run("error result", func(t *testing.T) {
		testErr := fmt.Errorf("test error")
		result := &templateResult{
			jobID: 2,
			error: testErr,
		}
		
		if result.error == nil {
			t.Error("Expected error result to have error")
		}
		
		if result.error.Error() != "test error" {
			t.Errorf("Expected error 'test error', got '%v'", result.error)
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
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(1, tx)
		defer processor.pool.Stop()
		
		// Create a temporary output directory
		outputDir := t.TempDir()
		
		// Simple config for testing
		config := types.ProjectConfig{
			Name:   "integration-test",
			Module: "github.com/test/integration-test",
		}
		
		// Test with empty template files list
		files := []types.TemplateFile{}
		context := make(map[string]any)
		tmpl := &types.Template{}
		
		results, err := processor.ProcessTemplates(files, outputDir, outputDir, context, config, tmpl)
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
		
		tx := &GenerationTransaction{}
		processor := NewParallelTemplateProcessor(1, tx)
		processor.pool.Stop()
		
		elapsed := time.Since(start)
		if elapsed > time.Second {
			t.Errorf("Processor lifecycle took too long: %v", elapsed)
		}
	})
}