package generator

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

// AdvancedProgressTracker provides enhanced progress tracking with visual feedback
type AdvancedProgressTracker struct {
	logger *GeneratorLogger
	
	// Progress state
	totalSteps    int
	currentStep   int
	stepNames     []string
	stepDetails   map[int]string
	
	// File tracking
	totalFiles    int
	processedFiles int
	
	// Timing
	startTime     time.Time
	stepStartTime time.Time
	
	// Display options
	showProgressBar bool
	showETA        bool
	barWidth       int
	
	// Thread safety
	mu sync.RWMutex
}

// ProgressOptions configures progress display
type ProgressOptions struct {
	ShowProgressBar bool
	ShowETA        bool
	BarWidth       int
	Quiet          bool
}

// DefaultProgressOptions returns sensible defaults
func DefaultProgressOptions() *ProgressOptions {
	return &ProgressOptions{
		ShowProgressBar: true,
		ShowETA:        true,
		BarWidth:       40,
		Quiet:          false,
	}
}

// NewAdvancedProgressTracker creates an enhanced progress tracker
func NewAdvancedProgressTracker(logger *GeneratorLogger, stepNames []string, opts *ProgressOptions) *AdvancedProgressTracker {
	if opts == nil {
		opts = DefaultProgressOptions()
	}
	
	return &AdvancedProgressTracker{
		logger:          logger,
		totalSteps:      len(stepNames),
		currentStep:     0,
		stepNames:       stepNames,
		stepDetails:     make(map[int]string),
		showProgressBar: opts.ShowProgressBar,
		showETA:        opts.ShowETA,
		barWidth:       opts.BarWidth,
		startTime:      time.Now(),
		stepStartTime:  time.Now(),
	}
}

// SetStepDetail adds detailed description for a step
func (apt *AdvancedProgressTracker) SetStepDetail(step int, detail string) {
	apt.mu.Lock()
	defer apt.mu.Unlock()
	
	apt.stepDetails[step] = detail
}

// SetTotalFiles sets the expected number of files to process
func (apt *AdvancedProgressTracker) SetTotalFiles(total int) {
	apt.mu.Lock()
	defer apt.mu.Unlock()
	
	apt.totalFiles = total
	apt.processedFiles = 0
}

// AddProcessedFiles increments the processed file count
func (apt *AdvancedProgressTracker) AddProcessedFiles(count int) {
	apt.mu.Lock()
	defer apt.mu.Unlock()
	
	apt.processedFiles += count
}

// NextStep advances to the next step with enhanced logging
func (apt *AdvancedProgressTracker) NextStep() {
	apt.mu.Lock()
	defer apt.mu.Unlock()
	
	if apt.currentStep < apt.totalSteps {
		apt.currentStep++
		apt.stepStartTime = time.Now()
		
		if apt.currentStep <= len(apt.stepNames) {
			stepName := apt.stepNames[apt.currentStep-1]
			detail := apt.stepDetails[apt.currentStep]
			
			// Create progress message with visual elements
			progress := float64(apt.currentStep) / float64(apt.totalSteps) * 100
			
			var msg strings.Builder
			
			// Progress bar
			if apt.showProgressBar {
				bar := apt.createProgressBar(progress)
				msg.WriteString(fmt.Sprintf("[%s] ", bar))
			}
			
			// Step info
			msg.WriteString(fmt.Sprintf("Step %d/%d: %s", apt.currentStep, apt.totalSteps, stepName))
			
			// Add detail if available
			if detail != "" {
				msg.WriteString(fmt.Sprintf(" - %s", detail))
			}
			
			// ETA calculation
			if apt.showETA && apt.currentStep > 1 {
				eta := apt.calculateETA()
				if eta > 0 {
					msg.WriteString(fmt.Sprintf(" (ETA: %v)", eta.Round(time.Second)))
				}
			}
			
			apt.logger.Step(apt.currentStep, apt.totalSteps, "%s", msg.String())
		}
	}
}

// NextStepWithDetail advances to next step and sets detail
func (apt *AdvancedProgressTracker) NextStepWithDetail(detail string) {
	apt.SetStepDetail(apt.currentStep+1, detail)
	apt.NextStep()
}

// UpdateCurrentStep updates the current step with new information
func (apt *AdvancedProgressTracker) UpdateCurrentStep(info string) {
	apt.mu.RLock()
	defer apt.mu.RUnlock()
	
	if apt.currentStep > 0 && apt.currentStep <= len(apt.stepNames) {
		stepName := apt.stepNames[apt.currentStep-1]
		progress := float64(apt.currentStep) / float64(apt.totalSteps) * 100
		
		var msg strings.Builder
		
		if apt.showProgressBar {
			bar := apt.createProgressBar(progress)
			msg.WriteString(fmt.Sprintf("[%s] ", bar))
		}
		
		msg.WriteString(fmt.Sprintf("%s - %s", stepName, info))
		
		apt.logger.Progress("%s", msg.String())
	}
}

// FileProgress shows progress for file operations
func (apt *AdvancedProgressTracker) FileProgress(fileName string) {
	apt.mu.Lock()
	apt.processedFiles++
	current := apt.processedFiles
	total := apt.totalFiles
	apt.mu.Unlock()
	
	if total > 0 {
		progress := float64(current) / float64(total) * 100
		
		var msg strings.Builder
		msg.WriteString(fmt.Sprintf("   Processing: %s (%d/%d files", fileName, current, total))
		
		if progress > 0 {
			msg.WriteString(fmt.Sprintf(", %.1f%%", progress))
		}
		msg.WriteString(")")
		
		apt.logger.Progress("%s", msg.String())
	} else {
		apt.logger.Progress("   Processing: %s", fileName)
	}
}

// Complete logs completion with summary
func (apt *AdvancedProgressTracker) Complete(msg string, args ...interface{}) {
	apt.mu.RLock()
	duration := time.Since(apt.startTime)
	filesCount := apt.processedFiles
	apt.mu.RUnlock()
	
	summary := fmt.Sprintf(msg, args...)
	if filesCount > 0 {
		summary += fmt.Sprintf(" (%d files processed in %v)", filesCount, duration.Round(time.Millisecond))
	} else {
		summary += fmt.Sprintf(" (completed in %v)", duration.Round(time.Millisecond))
	}
	
	apt.logger.Success("%s", summary)
}

// Error logs error with context
func (apt *AdvancedProgressTracker) Error(msg string, err error) {
	apt.mu.RLock()
	step := apt.currentStep
	stepName := ""
	if step > 0 && step <= len(apt.stepNames) {
		stepName = apt.stepNames[step-1]
	}
	apt.mu.RUnlock()
	
	apt.logger.ErrorWithDetails(msg, err, map[string]interface{}{
		"step":        step,
		"step_name":   stepName,
		"total_steps": apt.totalSteps,
	})
}

// createProgressBar creates a visual progress bar
func (apt *AdvancedProgressTracker) createProgressBar(progress float64) string {
	if apt.barWidth <= 0 {
		return ""
	}
	
	filled := int(progress * float64(apt.barWidth) / 100)
	if filled > apt.barWidth {
		filled = apt.barWidth
	}
	
	bar := strings.Repeat("█", filled) + strings.Repeat("░", apt.barWidth-filled)
	return fmt.Sprintf("%s %3.0f%%", bar, progress)
}

// calculateETA estimates time to completion
func (apt *AdvancedProgressTracker) calculateETA() time.Duration {
	if apt.currentStep <= 1 {
		return 0
	}
	
	elapsed := time.Since(apt.startTime)
	avgTimePerStep := elapsed / time.Duration(apt.currentStep)
	remaining := apt.totalSteps - apt.currentStep
	
	return avgTimePerStep * time.Duration(remaining)
}

// GetProgress returns current progress information
func (apt *AdvancedProgressTracker) GetProgress() ProgressInfo {
	apt.mu.RLock()
	defer apt.mu.RUnlock()
	
	return ProgressInfo{
		CurrentStep:    apt.currentStep,
		TotalSteps:     apt.totalSteps,
		ProcessedFiles: apt.processedFiles,
		TotalFiles:     apt.totalFiles,
		ElapsedTime:    time.Since(apt.startTime),
		ETA:           apt.calculateETA(),
		Percentage:    float64(apt.currentStep) / float64(apt.totalSteps) * 100,
	}
}

// ProgressInfo contains comprehensive progress information
type ProgressInfo struct {
	CurrentStep    int
	TotalSteps     int
	ProcessedFiles int
	TotalFiles     int
	ElapsedTime    time.Duration
	ETA           time.Duration
	Percentage    float64
}

// String returns a formatted progress summary
func (pi ProgressInfo) String() string {
	var parts []string
	
	parts = append(parts, fmt.Sprintf("Step %d/%d (%.1f%%)", pi.CurrentStep, pi.TotalSteps, pi.Percentage))
	
	if pi.TotalFiles > 0 {
		fileProgress := float64(pi.ProcessedFiles) / float64(pi.TotalFiles) * 100
		parts = append(parts, fmt.Sprintf("%d/%d files (%.1f%%)", pi.ProcessedFiles, pi.TotalFiles, fileProgress))
	}
	
	parts = append(parts, fmt.Sprintf("Elapsed: %v", pi.ElapsedTime.Round(time.Second)))
	
	if pi.ETA > 0 {
		parts = append(parts, fmt.Sprintf("ETA: %v", pi.ETA.Round(time.Second)))
	}
	
	return strings.Join(parts, " | ")
}

// GenerationSteps defines the standard steps for project generation
var GenerationSteps = []string{
	"Validating configuration",
	"Loading template",
	"Parsing template files", 
	"Generating project structure",
	"Writing files",
	"Initializing Git repository",
	"Installing dependencies",
	"Running post-generation hooks",
}

// DetailedGenerationSteps provides detailed descriptions for each step
func DetailedGenerationSteps(config types.ProjectConfig, outputDir string, initGit bool) map[int]string {
	details := make(map[int]string)
	
	details[1] = fmt.Sprintf("Validating project name '%s' and module path", config.Name)
	details[2] = fmt.Sprintf("Loading template '%s'", config.Type)
	details[3] = "Analyzing template files and dependencies"
	details[4] = fmt.Sprintf("Creating directory structure in %s", outputDir)
	details[5] = "Generating and writing project files"
	
	if initGit {
		details[6] = "Initializing Git repository with initial commit"
	} else {
		details[6] = "Skipping Git initialization"
	}
	
	details[7] = "Installing Go module dependencies"
	details[8] = "Running template-specific setup tasks"
	
	return details
}