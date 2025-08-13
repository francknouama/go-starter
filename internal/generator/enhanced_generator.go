package generator

import (
	"fmt"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

// GenerateWithProgressTracking generates a project with enhanced progress feedback
func (g *Generator) GenerateWithProgressTracking(config types.ProjectConfig, options types.GenerationOptions, progressOptions *ProgressOptions) (*types.GenerationResult, error) {
	startTime := time.Now()
	
	result := &types.GenerationResult{
		ProjectPath:  options.OutputPath,
		FilesCreated: []string{},
		Success:      false,
	}

	// Set up progress tracking
	if progressOptions == nil {
		progressOptions = DefaultProgressOptions()
	}
	
	// Create detailed step descriptions
	stepDetails := DetailedGenerationSteps(config, options.OutputPath, !options.NoGit)
	
	// Create progress tracker
	progressTracker := NewAdvancedProgressTracker(g.logger, GenerationSteps, progressOptions)
	for step, detail := range stepDetails {
		progressTracker.SetStepDetail(step, detail)
	}

	// Create transaction for rollback support
	tx := NewGenerationTransaction(options.OutputPath)

	// Set up recovery mechanism with progress reporting
	defer func() {
		if r := recover(); r != nil {
			progressTracker.Error("Generation panic occurred", fmt.Errorf("%v", r))
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				progressTracker.Error("Rollback failed", rollbackErr)
			}
			panic(r) // Re-panic after cleanup
		}
	}()

	// Step 1: Validate configuration
	progressTracker.NextStep()
	if err := g.validateConfig(config); err != nil {
		progressTracker.Error("Configuration validation failed", err)
		result.Error = err
		return result, err
	}

	// Step 2: Load template
	progressTracker.NextStep()
	templateID := g.getTemplateID(config)
	progressTracker.UpdateCurrentStep(fmt.Sprintf("Loading template '%s'", templateID))
	
	template, err := g.registry.Get(templateID)
	if err != nil {
		// Try fallback: look for templates by type if direct ID lookup fails
		templatesByType := g.registry.GetByType(config.Type)
		if len(templatesByType) == 0 {
			progressTracker.Error("Template not found", err)
			result.Error = err
			return result, err
		}
		template = templatesByType[0]
		templateID = template.ID
	}

	// Step 3: Parse template files
	progressTracker.NextStep()
	progressTracker.UpdateCurrentStep(fmt.Sprintf("Analyzing %d template files", len(template.Files)))
	
	// Filter files that should be generated based on conditions
	var filesToProcess []types.TemplateFile
	for _, file := range template.Files {
		shouldInclude, err := g.evaluateCondition(file.Condition, config)
		if err != nil {
			progressTracker.Error(fmt.Sprintf("Error evaluating condition for %s", file.Destination), err)
			result.Error = err
			return result, err
		}
		if shouldInclude {
			filesToProcess = append(filesToProcess, file)
		}
	}

	// Set total files for progress tracking
	progressTracker.SetTotalFiles(len(filesToProcess))
	progressTracker.UpdateCurrentStep(fmt.Sprintf("Found %d files to generate", len(filesToProcess)))

	// Step 4: Generate project structure
	progressTracker.NextStep()
	if err := g.checkOutputDirectory(options.OutputPath, options.ForceOverwrite); err != nil {
		progressTracker.Error("Output directory check failed", err)
		result.Error = err
		return result, err
	}

	// Step 5: Write files
	progressTracker.NextStep()
	
	// Process files with progress feedback
	for _, file := range filesToProcess {
		progressTracker.FileProgress(file.Destination)
		
		err := g.processTemplateFile(file, config, options.OutputPath, tx)
		if err != nil {
			progressTracker.Error(fmt.Sprintf("Failed to process file %s", file.Destination), err)
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				progressTracker.Error("Rollback failed", rollbackErr)
			}
			result.Error = err
			return result, err
		}
		
		result.FilesCreated = append(result.FilesCreated, file.Destination)
	}

	// Step 6: Initialize Git repository
	progressTracker.NextStep()
	if !options.NoGit {
		progressTracker.UpdateCurrentStep("Creating Git repository and initial commit")
		if err := g.initializeGitRepository(options.OutputPath); err != nil {
			progressTracker.Error("Git initialization failed", err)
			// Don't fail the entire generation for Git errors, just warn
			g.logger.Warning("Git initialization failed: %v", err)
		}
	} else {
		progressTracker.UpdateCurrentStep("Git initialization skipped")
	}

	// Step 7: Install dependencies
	progressTracker.NextStep()
	if true { // Always try to install dependencies unless explicitly disabled
		progressTracker.UpdateCurrentStep("Running 'go mod tidy' to install dependencies")
		if err := g.installDependencies(options.OutputPath); err != nil {
			progressTracker.Error("Dependency installation failed", err)
			// Don't fail the entire generation for dependency errors, just warn
			g.logger.Warning("Dependency installation failed: %v", err)
		}
		// Note: In a real implementation, you might want to add an InstallDeps field to GenerationOptions
		// For now, we try to install dependencies by default
	}

	// Step 8: Run post-generation hooks
	progressTracker.NextStep()
	if len(template.PostHooks) > 0 {
		progressTracker.UpdateCurrentStep(fmt.Sprintf("Running %d post-generation hooks", len(template.PostHooks)))
		for i, hook := range template.PostHooks {
			progressTracker.UpdateCurrentStep(fmt.Sprintf("Running hook %d/%d: %s", i+1, len(template.PostHooks), hook.Name))
			if err := g.executeHook(hook, options.OutputPath, config); err != nil {
				progressTracker.Error(fmt.Sprintf("Hook '%s' failed", hook.Name), err)
				g.logger.Warning("Post-generation hook failed: %v", err)
			}
		}
	} else {
		progressTracker.UpdateCurrentStep("No post-generation hooks defined")
	}

	// Complete the transaction
	if err := tx.Commit(); err != nil {
		progressTracker.Error("Transaction commit failed", err)
		result.Error = err
		return result, err
	}

	// Mark as successful
	result.Success = true
	result.GenerationTime = time.Since(startTime)

	// Complete progress tracking
	progressTracker.Complete("Project '%s' generated successfully", config.Name)
	
	return result, nil
}

// processTemplateFile processes a single template file with proper error handling
func (g *Generator) processTemplateFile(file types.TemplateFile, config types.ProjectConfig, outputPath string, tx *GenerationTransaction) error {
	return g.errorHandler.ExecuteWithRecovery(&RecoverableOperation{
		Name:        fmt.Sprintf("process_file_%s", file.Destination),
		Description: fmt.Sprintf("Process template file %s", file.Source),
		Execute: func() error {
			return g.generateFile(file, config, outputPath, tx)
		},
		Rollback: func() error {
			// File-level rollback is handled by the transaction
			return nil
		},
		MaxRetries: 2,
		RetryDelay: 100 * time.Millisecond,
	})
}

// initializeGitRepository initializes a Git repository with error recovery
func (g *Generator) initializeGitRepository(projectPath string) error {
	return g.errorHandler.ExecuteWithRecovery(&RecoverableOperation{
		Name:        "git_init",
		Description: "Initialize Git repository",
		Execute: func() error {
			return g.runGitCommand(projectPath, "init")
		},
		MaxRetries: 1,
		RetryDelay: 500 * time.Millisecond,
	})
}

// installDependencies installs Go module dependencies with error recovery
func (g *Generator) installDependencies(projectPath string) error {
	return g.errorHandler.ExecuteWithRecovery(&RecoverableOperation{
		Name:        "install_deps",
		Description: "Install Go module dependencies",
		Execute: func() error {
			return g.runGoCommand(projectPath, "mod", "tidy")
		},
		MaxRetries: 2,
		RetryDelay: time.Second,
	})
}

// executeHook executes a post-generation hook with error recovery
func (g *Generator) executeHook(hook types.Hook, projectPath string, config types.ProjectConfig) error {
	return g.errorHandler.ExecuteWithRecovery(&RecoverableOperation{
		Name:        fmt.Sprintf("hook_%s", hook.Name),
		Description: fmt.Sprintf("Execute hook: %s", hook.Description),
		Execute: func() error {
			return g.runHook(hook, projectPath, config)
		},
		MaxRetries: 1,
		RetryDelay: 500 * time.Millisecond,
	})
}

// GenerateWithQuietProgress generates a project with minimal progress output
func (g *Generator) GenerateWithQuietProgress(config types.ProjectConfig, options types.GenerationOptions) (*types.GenerationResult, error) {
	quietOptions := &ProgressOptions{
		ShowProgressBar: false,
		ShowETA:        false,
		BarWidth:       0,
		Quiet:          true,
	}
	
	return g.GenerateWithProgressTracking(config, options, quietOptions)
}

// GenerateWithVerboseProgress generates a project with detailed progress output
func (g *Generator) GenerateWithVerboseProgress(config types.ProjectConfig, options types.GenerationOptions) (*types.GenerationResult, error) {
	verboseOptions := &ProgressOptions{
		ShowProgressBar: true,
		ShowETA:        true,
		BarWidth:       50, // Wider progress bar for verbose mode
		Quiet:          false,
	}
	
	return g.GenerateWithProgressTracking(config, options, verboseOptions)
}