package optimization

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// OptimizationPipeline orchestrates comprehensive code optimization
type OptimizationPipeline struct {
	analyzer *ASTAnalyzer
	manager  *ImportManager
	options  PipelineOptions
}

// PipelineOptions configures the optimization pipeline
type PipelineOptions struct {
	// Import optimizations
	RemoveUnusedImports bool
	OrganizeImports     bool
	AddMissingImports   bool

	// Code optimizations
	RemoveUnusedVars      bool
	RemoveUnusedFuncs     bool
	OptimizeConditionals  bool

	// Output options
	WriteOptimizedFiles bool
	CreateBackups       bool
	DryRun             bool
	Verbose            bool

	// Performance limits
	MaxFileSize         int64
	MaxConcurrentFiles  int
	SkipTestFiles      bool
	SkipVendorDirs     bool

	// File patterns
	IncludePatterns []string
	ExcludePatterns []string
}

// DefaultPipelineOptions returns sensible defaults for the optimization pipeline
func DefaultPipelineOptions() PipelineOptions {
	return PipelineOptions{
		// Import optimizations - conservative defaults
		RemoveUnusedImports: true,
		OrganizeImports:     true,
		AddMissingImports:   false, // Can be risky without proper analysis

		// Code optimizations - very conservative defaults
		RemoveUnusedVars:      false,
		RemoveUnusedFuncs:     false,
		OptimizeConditionals:  false,

		// Output options
		WriteOptimizedFiles: false, // Default to dry run
		CreateBackups:       true,
		DryRun:             true,
		Verbose:            false,

		// Performance limits
		MaxFileSize:        1024 * 1024, // 1MB
		MaxConcurrentFiles: 4,
		SkipTestFiles:      true,
		SkipVendorDirs:     true,

		// File patterns
		IncludePatterns: []string{"*.go"},
		ExcludePatterns: []string{"vendor/"}, // .git/ and node_modules/ are handled by findGoFiles
	}
}

// PipelineResult contains the results of running the optimization pipeline
type PipelineResult struct {
	// Summary statistics
	TotalFiles         int
	FilesProcessed     int
	FilesOptimized     int
	FilesSkipped       int
	FilesWithErrors    int

	// Optimization metrics
	ImportsRemoved     int
	ImportsAdded       int
	ImportsOrganized   int
	VariablesRemoved   int
	FunctionsRemoved   int

	// Performance metrics
	ProcessingTimeMs   int64
	SizeBeforeBytes    int64
	SizeAfterBytes     int64
	SizeReductionBytes int64

	// Detailed results per file
	FileResults map[string]*FileOptimizationResult

	// Errors encountered
	Errors []error
}

// FileOptimizationResult contains optimization results for a single file
type FileOptimizationResult struct {
	FilePath           string
	OriginalSize       int64
	OptimizedSize      int64
	OptimizationApplied bool
	ImportsResult      *ImportOptimizationResult
	BackupPath         string
	Errors             []error
}

// NewOptimizationPipeline creates a new optimization pipeline
func NewOptimizationPipeline(options PipelineOptions) *OptimizationPipeline {
	// Create analyzer with options from pipeline
	analysisOptions := AnalysisOptions{
		RemoveUnusedImports:   options.RemoveUnusedImports,
		RemoveUnusedVars:      options.RemoveUnusedVars,
		RemoveUnusedFuncs:     options.RemoveUnusedFuncs,
		OptimizeConditionals:  options.OptimizeConditionals,
		OrganizeImports:       options.OrganizeImports,
		EnableDebugOutput:     options.Verbose,
		MaxFileSize:           options.MaxFileSize,
	}

	analyzer := NewASTAnalyzer(analysisOptions)
	manager := NewImportManager(analyzer)

	return &OptimizationPipeline{
		analyzer: analyzer,
		manager:  manager,
		options:  options,
	}
}

// OptimizeProject optimizes all Go files in a project directory
func (p *OptimizationPipeline) OptimizeProject(projectPath string) (*PipelineResult, error) {
	startTime := time.Now()

	result := &PipelineResult{
		FileResults: make(map[string]*FileOptimizationResult),
		Errors:      make([]error, 0),
	}

	p.log("Starting optimization pipeline for project: %s", projectPath)

	// Find all Go files in the project
	goFiles, err := p.findGoFiles(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to find Go files: %w", err)
	}

	result.TotalFiles = len(goFiles)
	p.log("Found %d Go files to analyze", result.TotalFiles)

	// Process each file
	for _, filePath := range goFiles {
		fileResult := p.optimizeFile(filePath)
		result.FileResults[filePath] = fileResult

		// Update summary statistics
		result.FilesProcessed++
		
		if fileResult.OptimizationApplied {
			result.FilesOptimized++
		}

		if len(fileResult.Errors) > 0 {
			result.FilesWithErrors++
			result.Errors = append(result.Errors, fileResult.Errors...)
		}

		// Update optimization metrics
		if fileResult.ImportsResult != nil {
			result.ImportsRemoved += len(fileResult.ImportsResult.RemovedImports)
			result.ImportsAdded += len(fileResult.ImportsResult.AddedImports)
			if fileResult.ImportsResult.OptimizationApplied {
				result.ImportsOrganized++
			}
		}

		// Update size metrics
		result.SizeBeforeBytes += fileResult.OriginalSize
		result.SizeAfterBytes += fileResult.OptimizedSize
	}

	// Calculate final metrics
	result.ProcessingTimeMs = time.Since(startTime).Milliseconds()
	result.SizeReductionBytes = result.SizeBeforeBytes - result.SizeAfterBytes
	result.FilesSkipped = result.TotalFiles - result.FilesProcessed

	p.log("Pipeline completed in %dms", result.ProcessingTimeMs)
	p.logSummary(result)

	return result, nil
}

// OptimizeFile optimizes a single Go file
func (p *OptimizationPipeline) OptimizeFile(filePath string) (*FileOptimizationResult, error) {
	result := p.optimizeFile(filePath)
	if len(result.Errors) > 0 {
		return result, result.Errors[0]
	}
	return result, nil
}

// optimizeFile is the internal implementation for optimizing a single file
func (p *OptimizationPipeline) optimizeFile(filePath string) *FileOptimizationResult {
	result := &FileOptimizationResult{
		FilePath: filePath,
		Errors:   make([]error, 0),
	}

	p.log("Processing file: %s", filePath)

	// Check if file should be skipped
	if p.shouldSkipFile(filePath) {
		p.log("Skipping file: %s", filePath)
		return result
	}

	// Get original file size
	info, err := os.Stat(filePath)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to stat file: %w", err))
		return result
	}

	result.OriginalSize = info.Size()
	result.OptimizedSize = result.OriginalSize // Default to no change

	// Check file size limits
	if result.OriginalSize > p.options.MaxFileSize {
		result.Errors = append(result.Errors, fmt.Errorf("file too large: %d bytes (max: %d)", result.OriginalSize, p.options.MaxFileSize))
		return result
	}

	// Parse the file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("failed to parse file: %w", err))
		return result
	}

	p.analyzer.fileSet = fset

	// Optimize imports if requested
	if p.options.RemoveUnusedImports || p.options.OrganizeImports || p.options.AddMissingImports {
		importsResult, err := p.manager.OptimizeImports(file)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to optimize imports: %w", err))
			return result
		}

		result.ImportsResult = importsResult
		if importsResult.OptimizationApplied {
			result.OptimizationApplied = true
		}
	}

	// Write optimized file if requested and optimization was applied
	if result.OptimizationApplied && p.options.WriteOptimizedFiles && !p.options.DryRun {
		err = p.writeOptimizedFile(file, filePath, result)
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("failed to write optimized file: %w", err))
			return result
		}
	}

	return result
}

// writeOptimizedFile writes the optimized file to disk
func (p *OptimizationPipeline) writeOptimizedFile(file *ast.File, filePath string, result *FileOptimizationResult) error {
	// Create backup if requested
	if p.options.CreateBackups {
		backupPath := filePath + ".bak"
		err := p.createBackup(filePath, backupPath)
		if err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		result.BackupPath = backupPath
		p.log("Created backup: %s", backupPath)
	}

	// Get optimized content
	optimizedContent, err := p.manager.GetOptimizedFileContent(file)
	if err != nil {
		return fmt.Errorf("failed to get optimized content: %w", err)
	}

	// Write optimized file
	err = os.WriteFile(filePath, []byte(optimizedContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	result.OptimizedSize = int64(len(optimizedContent))
	p.log("Wrote optimized file: %s (%d -> %d bytes)", filePath, result.OriginalSize, result.OptimizedSize)

	return nil
}

// createBackup creates a backup of the original file
func (p *OptimizationPipeline) createBackup(originalPath, backupPath string) error {
	content, err := os.ReadFile(originalPath)
	if err != nil {
		return err
	}

	return os.WriteFile(backupPath, content, 0644)
}

// shouldSkipFile determines if a file should be skipped based on options
func (p *OptimizationPipeline) shouldSkipFile(filePath string) bool {
	// Skip test files if requested
	if p.options.SkipTestFiles && strings.HasSuffix(filePath, "_test.go") {
		return true
	}

	// Skip vendor directories if requested
	if p.options.SkipVendorDirs && strings.Contains(filePath, "/vendor/") {
		return true
	}

	// Check exclude patterns
	for _, pattern := range p.options.ExcludePatterns {
		if strings.Contains(filePath, pattern) {
			return true
		}
	}

	// Check include patterns (if any specified)
	if len(p.options.IncludePatterns) > 0 {
		matched := false
		for _, pattern := range p.options.IncludePatterns {
			// Simple pattern matching - in a real implementation would use filepath.Match
			if strings.HasSuffix(filePath, strings.TrimPrefix(pattern, "*")) {
				matched = true
				break
			}
		}
		if !matched {
			return true
		}
	}

	return false
}

// findGoFiles finds all Go files in the project directory
func (p *OptimizationPipeline) findGoFiles(projectPath string) ([]string, error) {
	var goFiles []string

	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			// Skip hidden directories and common non-source directories
			name := info.Name()
			if strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" {
				return filepath.SkipDir
			}
			return nil
		}

		// Include Go files that match our criteria
		if strings.HasSuffix(path, ".go") && !p.shouldSkipFile(path) {
			goFiles = append(goFiles, path)
		}

		return nil
	})

	return goFiles, err
}

// log outputs a message if verbose mode is enabled
func (p *OptimizationPipeline) log(format string, args ...interface{}) {
	if p.options.Verbose {
		fmt.Printf("[PIPELINE] "+format+"\n", args...)
	}
}

// logSummary outputs a summary of the optimization results
func (p *OptimizationPipeline) logSummary(result *PipelineResult) {
	if !p.options.Verbose {
		return
	}

	fmt.Printf("\n=== Optimization Pipeline Summary ===\n")
	fmt.Printf("Files processed: %d/%d\n", result.FilesProcessed, result.TotalFiles)
	fmt.Printf("Files optimized: %d\n", result.FilesOptimized)
	fmt.Printf("Files with errors: %d\n", result.FilesWithErrors)
	fmt.Printf("Imports removed: %d\n", result.ImportsRemoved)
	fmt.Printf("Imports added: %d\n", result.ImportsAdded)
	fmt.Printf("Files with organized imports: %d\n", result.ImportsOrganized)
	
	if result.SizeReductionBytes > 0 {
		reductionPercent := float64(result.SizeReductionBytes) / float64(result.SizeBeforeBytes) * 100
		fmt.Printf("Size reduction: %d bytes (%.2f%%)\n", result.SizeReductionBytes, reductionPercent)
	}
	
	fmt.Printf("Processing time: %dms\n", result.ProcessingTimeMs)
	fmt.Printf("=====================================\n")
}