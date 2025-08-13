package generator

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

func TestAdvancedProgressTracker(t *testing.T) {
	t.Run("creation and basic properties", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1", "Step 2", "Step 3"}
		opts := DefaultProgressOptions()
		
		tracker := NewAdvancedProgressTracker(logger, steps, opts)
		
		if tracker == nil {
			t.Error("Expected tracker to not be nil")
		}
		
		if tracker.totalSteps != 3 {
			t.Errorf("Expected 3 total steps, got %d", tracker.totalSteps)
		}
		
		if tracker.currentStep != 0 {
			t.Errorf("Expected current step to be 0, got %d", tracker.currentStep)
		}
		
		if !tracker.showProgressBar {
			t.Error("Expected progress bar to be enabled by default")
		}
		
		if !tracker.showETA {
			t.Error("Expected ETA to be enabled by default")
		}
	})
	
	t.Run("custom options", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1"}
		opts := &ProgressOptions{
			ShowProgressBar: false,
			ShowETA:        false,
			BarWidth:       20,
			Quiet:          true,
		}
		
		tracker := NewAdvancedProgressTracker(logger, steps, opts)
		
		if tracker.showProgressBar {
			t.Error("Expected progress bar to be disabled")
		}
		
		if tracker.showETA {
			t.Error("Expected ETA to be disabled")
		}
		
		if tracker.barWidth != 20 {
			t.Errorf("Expected bar width 20, got %d", tracker.barWidth)
		}
	})
	
	t.Run("step progression", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1", "Step 2", "Step 3"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		// Initial state
		progress := tracker.GetProgress()
		if progress.CurrentStep != 0 {
			t.Errorf("Expected current step 0, got %d", progress.CurrentStep)
		}
		if progress.Percentage != 0 {
			t.Errorf("Expected 0%% progress, got %.1f%%", progress.Percentage)
		}
		
		// First step
		tracker.NextStep()
		progress = tracker.GetProgress()
		if progress.CurrentStep != 1 {
			t.Errorf("Expected current step 1, got %d", progress.CurrentStep)
		}
		if progress.Percentage < 33 || progress.Percentage > 34 {
			t.Errorf("Expected ~33%% progress, got %.1f%%", progress.Percentage)
		}
		
		// Second step
		tracker.NextStep()
		progress = tracker.GetProgress()
		if progress.CurrentStep != 2 {
			t.Errorf("Expected current step 2, got %d", progress.CurrentStep)
		}
		
		// Final step
		tracker.NextStep()
		progress = tracker.GetProgress()
		if progress.CurrentStep != 3 {
			t.Errorf("Expected current step 3, got %d", progress.CurrentStep)
		}
		if progress.Percentage != 100 {
			t.Errorf("Expected 100%% progress, got %.1f%%", progress.Percentage)
		}
		
		// Beyond bounds
		tracker.NextStep()
		progress = tracker.GetProgress()
		if progress.CurrentStep != 3 {
			t.Errorf("Expected current step to remain 3, got %d", progress.CurrentStep)
		}
	})
	
	t.Run("step details", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Validation", "Generation"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		// Set step details
		tracker.SetStepDetail(1, "Validating project configuration")
		tracker.SetStepDetail(2, "Generating 25 files")
		
		// Check details are stored
		if len(tracker.stepDetails) != 2 {
			t.Errorf("Expected 2 step details, got %d", len(tracker.stepDetails))
		}
		
		if tracker.stepDetails[1] != "Validating project configuration" {
			t.Errorf("Expected step 1 detail, got '%s'", tracker.stepDetails[1])
		}
	})
	
	t.Run("file tracking", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Processing"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		// Set total files
		tracker.SetTotalFiles(10)
		progress := tracker.GetProgress()
		if progress.TotalFiles != 10 {
			t.Errorf("Expected 10 total files, got %d", progress.TotalFiles)
		}
		if progress.ProcessedFiles != 0 {
			t.Errorf("Expected 0 processed files, got %d", progress.ProcessedFiles)
		}
		
		// Add processed files
		tracker.AddProcessedFiles(3)
		progress = tracker.GetProgress()
		if progress.ProcessedFiles != 3 {
			t.Errorf("Expected 3 processed files, got %d", progress.ProcessedFiles)
		}
		
		// Test file progress
		tracker.FileProgress("main.go")
		progress = tracker.GetProgress()
		if progress.ProcessedFiles != 4 {
			t.Errorf("Expected 4 processed files after FileProgress, got %d", progress.ProcessedFiles)
		}
	})
	
	t.Run("next step with detail", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1", "Step 2"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		tracker.NextStepWithDetail("Custom detail for step 1")
		
		if tracker.currentStep != 1 {
			t.Errorf("Expected current step 1, got %d", tracker.currentStep)
		}
		
		if tracker.stepDetails[1] != "Custom detail for step 1" {
			t.Errorf("Expected step detail, got '%s'", tracker.stepDetails[1])
		}
	})
	
	t.Run("timing and ETA", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Step 1", "Step 2", "Step 3", "Step 4"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		// First step
		tracker.NextStep()
		time.Sleep(10 * time.Millisecond)
		
		// Second step
		tracker.NextStep()
		time.Sleep(10 * time.Millisecond)
		
		progress := tracker.GetProgress()
		
		// Should have elapsed time
		if progress.ElapsedTime < 20*time.Millisecond {
			t.Error("Expected measurable elapsed time")
		}
		
		// Should have ETA estimate (though it may be very small due to fast execution)
		if progress.ETA < 0 {
			t.Error("Expected non-negative ETA")
		}
	})
	
	t.Run("update current step", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Processing"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		tracker.NextStep()
		
		// This should not panic and should provide meaningful output
		tracker.UpdateCurrentStep("Processing configuration files")
	})
	
	t.Run("completion and error", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Test Step"}
		tracker := NewAdvancedProgressTracker(logger, steps, nil)
		
		tracker.SetTotalFiles(5)
		tracker.AddProcessedFiles(5)
		
		// These should not panic
		tracker.Complete("Test completed successfully")
		tracker.Error("Test error occurred", nil)
	})
}

func TestProgressBar(t *testing.T) {
	t.Run("progress bar creation", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Test"}
		opts := &ProgressOptions{
			ShowProgressBar: true,
			BarWidth:       10,
		}
		tracker := NewAdvancedProgressTracker(logger, steps, opts)
		
		// Test progress bar at different percentages
		tests := []struct {
			progress float64
			expected int // expected number of filled characters
		}{
			{0, 0},
			{25, 2},  // 2.5 rounded down
			{50, 5},
			{75, 7},  // 7.5 rounded down  
			{100, 10},
		}
		
		for _, tt := range tests {
			bar := tracker.createProgressBar(tt.progress)
			
			// Count filled characters (█)
			filled := strings.Count(bar, "█")
			if filled != tt.expected {
				t.Errorf("Progress %.0f%%: expected %d filled chars, got %d in bar: %s", 
					tt.progress, tt.expected, filled, bar)
			}
			
			// Should contain percentage
			if !strings.Contains(bar, fmt.Sprintf("%.0f%%", tt.progress)) {
				t.Errorf("Progress bar should contain percentage: %s", bar)
			}
		}
	})
	
	t.Run("zero width bar", func(t *testing.T) {
		logger := NewGeneratorLogger(LogLevelInfo)
		steps := []string{"Test"}
		opts := &ProgressOptions{
			ShowProgressBar: true,
			BarWidth:       0,
		}
		tracker := NewAdvancedProgressTracker(logger, steps, opts)
		
		bar := tracker.createProgressBar(50)
		if bar != "" {
			t.Errorf("Expected empty bar for zero width, got: %s", bar)
		}
	})
}

func TestProgressInfo(t *testing.T) {
	t.Run("progress info string representation", func(t *testing.T) {
		info := ProgressInfo{
			CurrentStep:    2,
			TotalSteps:     4,
			ProcessedFiles: 15,
			TotalFiles:     30,
			ElapsedTime:    45 * time.Second,
			ETA:           90 * time.Second,
			Percentage:    50.0,
		}
		
		str := info.String()
		
		// Should contain step progress
		if !strings.Contains(str, "Step 2/4") {
			t.Errorf("String should contain step progress: %s", str)
		}
		
		// Should contain percentage
		if !strings.Contains(str, "50.0%") {
			t.Errorf("String should contain percentage: %s", str)
		}
		
		// Should contain file progress
		if !strings.Contains(str, "15/30 files") {
			t.Errorf("String should contain file progress: %s", str)
		}
		
		// Should contain elapsed time
		if !strings.Contains(str, "Elapsed: 45s") {
			t.Errorf("String should contain elapsed time: %s", str)
		}
		
		// Should contain ETA
		if !strings.Contains(str, "ETA: 1m30s") {
			t.Errorf("String should contain ETA: %s", str)
		}
	})
	
	t.Run("progress info without files", func(t *testing.T) {
		info := ProgressInfo{
			CurrentStep: 1,
			TotalSteps:  3,
			ElapsedTime: 10 * time.Second,
			Percentage:  33.3,
		}
		
		str := info.String()
		
		// Should not contain file information
		if strings.Contains(str, "files") {
			t.Errorf("String should not contain file info when no files: %s", str)
		}
		
		// Should contain step and time info
		if !strings.Contains(str, "Step 1/3") || !strings.Contains(str, "Elapsed: 10s") {
			t.Errorf("String should contain basic progress info: %s", str)
		}
	})
}

func TestGenerationSteps(t *testing.T) {
	t.Run("standard generation steps", func(t *testing.T) {
		if len(GenerationSteps) != 8 {
			t.Errorf("Expected 8 generation steps, got %d", len(GenerationSteps))
		}
		
		expectedSteps := []string{
			"Validating configuration",
			"Loading template", 
			"Parsing template files",
			"Generating project structure",
			"Writing files",
			"Initializing Git repository",
			"Installing dependencies",
			"Running post-generation hooks",
		}
		
		for i, expected := range expectedSteps {
			if i >= len(GenerationSteps) {
				t.Errorf("Missing step at index %d: %s", i, expected)
				continue
			}
			
			if GenerationSteps[i] != expected {
				t.Errorf("Step %d: expected '%s', got '%s'", i, expected, GenerationSteps[i])
			}
		}
	})
	
	t.Run("detailed generation steps", func(t *testing.T) {
		config := types.ProjectConfig{
			Name: "test-project",
			Type: "web-api",
		}
		
		details := DetailedGenerationSteps(config, "/tmp/test", true)
		
		if len(details) != 8 {
			t.Errorf("Expected 8 detailed steps, got %d", len(details))
		}
		
		// Check specific details
		if !strings.Contains(details[1], "test-project") {
			t.Errorf("Step 1 should contain project name: %s", details[1])
		}
		
		if !strings.Contains(details[2], "web-api") {
			t.Errorf("Step 2 should contain template ID: %s", details[2])
		}
		
		if !strings.Contains(details[4], "/tmp/test") {
			t.Errorf("Step 4 should contain output directory: %s", details[4])
		}
		
		if !strings.Contains(details[6], "Initializing Git") {
			t.Errorf("Step 6 should mention Git initialization: %s", details[6])
		}
	})
	
	t.Run("detailed steps without git", func(t *testing.T) {
		config := types.ProjectConfig{
			Name: "test-project",
			Type: "cli",
		}
		
		details := DetailedGenerationSteps(config, "/tmp/test", false)
		
		if !strings.Contains(details[6], "Skipping Git") {
			t.Errorf("Step 6 should mention skipping Git: %s", details[6])
		}
	})
}

func TestProgressOptions(t *testing.T) {
	t.Run("default options", func(t *testing.T) {
		opts := DefaultProgressOptions()
		
		if !opts.ShowProgressBar {
			t.Error("Expected progress bar to be enabled by default")
		}
		
		if !opts.ShowETA {
			t.Error("Expected ETA to be enabled by default")
		}
		
		if opts.BarWidth != 40 {
			t.Errorf("Expected default bar width 40, got %d", opts.BarWidth)
		}
		
		if opts.Quiet {
			t.Error("Expected quiet to be false by default")
		}
	})
}