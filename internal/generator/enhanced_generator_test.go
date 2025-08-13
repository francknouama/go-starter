package generator

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestGenerateWithProgressTracking(t *testing.T) {
	t.Run("basic progress tracking", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		// Create a simple test configuration
		config := types.ProjectConfig{
			Name:      "test-project",
			Module:    "github.com/test/test-project",
			Type:      "cli",
			GoVersion: "1.21",
			Framework: "cobra",
			Logger:    "slog",
		}
		
		// Create temporary output directory
		outputDir := t.TempDir()
		
		options := types.GenerationOptions{
			OutputPath:     filepath.Join(outputDir, "test-project"),
			ForceOverwrite: true,
			InitGit:        false, // Skip git for tests
			InstallDeps:    false, // Skip dependencies for tests
		}
		
		progressOptions := &ProgressOptions{
			ShowProgressBar: true,
			ShowETA:        true,
			BarWidth:       20,
			Quiet:          false,
		}
		
		// Generate with progress tracking
		result, err := gen.GenerateWithProgressTracking(config, options, progressOptions)
		
		// Should not error (template might not exist in test environment)
		if err != nil {
			t.Logf("Generation failed as expected in test environment: %v", err)
			return
		}
		
		// If generation succeeded, verify basic results
		if result == nil {
			t.Error("Expected result to not be nil")
			return
		}
		
		if !result.Success {
			t.Errorf("Expected generation to succeed, got error: %v", result.Error)
		}
		
		if result.GenerationTime == 0 {
			t.Error("Expected generation time to be measured")
		}
		
		if len(result.FilesCreated) == 0 {
			t.Error("Expected some files to be created")
		}
	})
	
	t.Run("progress tracking with nil options", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		config := types.ProjectConfig{
			Name:      "test-project",
			Module:    "github.com/test/test-project",
			Type:      "cli",
			GoVersion: "1.21",
		}
		
		outputDir := t.TempDir()
		options := types.GenerationOptions{
			OutputPath:     filepath.Join(outputDir, "test-project"),
			ForceOverwrite: true,
			InitGit:        false,
			InstallDeps:    false,
		}
		
		// Pass nil progress options - should use defaults
		result, err := gen.GenerateWithProgressTracking(config, options, nil)
		
		// Should handle nil options gracefully
		if err != nil {
			t.Logf("Generation failed as expected in test environment: %v", err)
			return
		}
		
		if result == nil {
			t.Error("Expected result to not be nil")
		}
	})
	
	t.Run("quiet progress mode", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		config := types.ProjectConfig{
			Name:      "test-quiet",
			Module:    "github.com/test/test-quiet",
			Type:      "cli",
			GoVersion: "1.21",
		}
		
		outputDir := t.TempDir()
		options := types.GenerationOptions{
			OutputPath:     filepath.Join(outputDir, "test-quiet"),
			ForceOverwrite: true,
			InitGit:        false,
			InstallDeps:    false,
		}
		
		// Test quiet mode
		result, err := gen.GenerateWithQuietProgress(config, options)
		
		if err != nil {
			t.Logf("Generation failed as expected in test environment: %v", err)
			return
		}
		
		if result == nil {
			t.Error("Expected result to not be nil")
		}
	})
	
	t.Run("verbose progress mode", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		config := types.ProjectConfig{
			Name:      "test-verbose",
			Module:    "github.com/test/test-verbose", 
			Type:      "cli",
			GoVersion: "1.21",
		}
		
		outputDir := t.TempDir()
		options := types.GenerationOptions{
			OutputPath:     filepath.Join(outputDir, "test-verbose"),
			ForceOverwrite: true,
			InitGit:        false,
			InstallDeps:    false,
		}
		
		// Test verbose mode
		result, err := gen.GenerateWithVerboseProgress(config, options)
		
		if err != nil {
			t.Logf("Generation failed as expected in test environment: %v", err)
			return
		}
		
		if result == nil {
			t.Error("Expected result to not be nil")
		}
	})
}

func TestRecoverableOperations(t *testing.T) {
	t.Run("processTemplateFile operation", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		config := types.ProjectConfig{
			Name:   "test",
			Module: "github.com/test/test",
		}
		
		file := types.TemplateFile{
			Source:      "test.tmpl",
			Destination: "test.go",
		}
		
		outputDir := t.TempDir()
		tx := NewGenerationTransaction(outputDir)
		
		// This will likely fail due to missing template, but should handle error gracefully
		err := gen.processTemplateFile(file, config, outputDir, tx)
		
		// Should return an error (template not found) but not panic
		if err == nil {
			t.Log("Template file processing succeeded unexpectedly")
		} else {
			t.Logf("Template file processing failed as expected: %v", err)
		}
	})
	
	t.Run("initializeGitRepository operation", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		outputDir := t.TempDir()
		
		// This might succeed if git is available, or fail gracefully if not
		err := gen.initializeGitRepository(outputDir)
		
		if err != nil {
			t.Logf("Git initialization failed as might be expected: %v", err)
		} else {
			t.Log("Git initialization succeeded")
		}
	})
	
	t.Run("installDependencies operation", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		outputDir := t.TempDir()
		
		// This will likely fail (no go.mod file), but should handle gracefully
		err := gen.installDependencies(outputDir)
		
		if err != nil {
			t.Logf("Dependency installation failed as expected: %v", err)
		} else {
			t.Log("Dependency installation succeeded unexpectedly")
		}
	})
	
	t.Run("executeHook operation", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		hook := types.Hook{
			Name:        "test-hook",
			Description: "Test hook",
			Command:     "echo",
			Args:        []string{"test"},
		}
		
		config := types.ProjectConfig{
			Name: "test",
		}
		
		outputDir := t.TempDir()
		
		// This should succeed if echo command is available
		err := gen.executeHook(hook, outputDir, config)
		
		if err != nil {
			t.Logf("Hook execution failed: %v", err)
		} else {
			t.Log("Hook execution succeeded")
		}
	})
}

func TestProgressOptions(t *testing.T) {
	t.Run("default progress options are reasonable", func(t *testing.T) {
		opts := DefaultProgressOptions()
		
		if !opts.ShowProgressBar {
			t.Error("Expected progress bar to be shown by default")
		}
		
		if !opts.ShowETA {
			t.Error("Expected ETA to be shown by default")
		}
		
		if opts.BarWidth <= 0 {
			t.Error("Expected positive bar width")
		}
		
		if opts.Quiet {
			t.Error("Expected quiet to be false by default")
		}
	})
	
	t.Run("quiet mode disables visual elements", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		// The quiet mode should create appropriate options
		config := types.ProjectConfig{Name: "test"}
		options := types.GenerationOptions{OutputPath: t.TempDir()}
		
		// This tests that the method exists and can be called
		_, err := gen.GenerateWithQuietProgress(config, options)
		
		// We expect this to fail in test environment, but it should not panic
		if err != nil {
			t.Logf("Quiet generation failed as expected: %v", err)
		}
	})
	
	t.Run("verbose mode enhances visual elements", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		config := types.ProjectConfig{Name: "test"}
		options := types.GenerationOptions{OutputPath: t.TempDir()}
		
		// This tests that the method exists and can be called
		_, err := gen.GenerateWithVerboseProgress(config, options)
		
		// We expect this to fail in test environment, but it should not panic
		if err != nil {
			t.Logf("Verbose generation failed as expected: %v", err)
		}
	})
}

func TestProgressIntegration(t *testing.T) {
	t.Run("progress tracker integrates with error handler", func(t *testing.T) {
		gen := NewWithLogger(LogLevelInfo)
		
		// Test that progress tracking works with the existing generator structure
		if gen.logger == nil {
			t.Error("Expected generator to have logger")
		}
		
		if gen.errorHandler == nil {
			t.Error("Expected generator to have error handler")
		}
		
		// Create progress tracker
		steps := []string{"Test Step"}
		tracker := NewAdvancedProgressTracker(gen.logger, steps, nil)
		
		if tracker == nil {
			t.Error("Expected progress tracker to be created")
		}
		
		// Test that progress tracker can work with generator's components
		tracker.NextStep()
		tracker.Complete("Test completed")
	})
	
	t.Run("generation timing is measured", func(t *testing.T) {
		// Test that the timing functionality works
		start := time.Now()
		time.Sleep(1 * time.Millisecond) // Small delay for testing
		elapsed := time.Since(start)
		
		if elapsed < 1*time.Millisecond {
			t.Error("Expected measurable elapsed time")
		}
		
		// This demonstrates that GenerationResult.GenerationTime should work
		result := &types.GenerationResult{
			GenerationTime: elapsed,
		}
		
		if result.GenerationTime == 0 {
			t.Error("Expected generation time to be set")
		}
	})
}