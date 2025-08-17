package generator

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
)

// WorkerPool manages parallel template processing
type WorkerPool struct {
	numWorkers      int
	jobs            chan *templateJob
	results         chan *templateResult
	errors          chan error
	ctx             context.Context
	cancel          context.CancelFunc
	wg              sync.WaitGroup
	dirCache        *dirCache
	progressCounter *int64
	totalJobs       int64
}

// templateJob represents a single template file processing job
type templateJob struct {
	templateFile types.TemplateFile
	templateDir  string
	destPath     string
	context      map[string]any
	jobID        int
	transaction  *GenerationTransaction
}

// templateResult represents the result of processing a template job
type templateResult struct {
	jobID       int
	filePath    string
	error       error
	wasSkipped  bool
	skipReason  string
}

// dirCache provides thread-safe directory creation with caching
type dirCache struct {
	created map[string]bool
	mu      sync.RWMutex
}

// newDirCache creates a new directory cache
func newDirCache() *dirCache {
	return &dirCache{
		created: make(map[string]bool),
	}
}

// ensureDir creates a directory if it doesn't exist, using cache to avoid redundant calls
func (dc *dirCache) ensureDir(path string, tx *GenerationTransaction) error {
	// First check cache with read lock
	dc.mu.RLock()
	if dc.created[path] {
		dc.mu.RUnlock()
		return nil
	}
	dc.mu.RUnlock()

	// Check if directory actually exists
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		dc.mu.Lock()
		dc.created[path] = true
		dc.mu.Unlock()
		return nil
	}

	// Need to create directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	// Cache the creation and track for rollback
	dc.mu.Lock()
	dc.created[path] = true
	if tx != nil {
		tx.AddDirectory(path)
	}
	dc.mu.Unlock()

	return nil
}

// NewWorkerPool creates a new worker pool for parallel template processing
func NewWorkerPool(numWorkers int) *WorkerPool {
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
	}

	ctx, cancel := context.WithCancel(context.Background())
	
	var counter int64 = 0

	return &WorkerPool{
		numWorkers:      numWorkers,
		jobs:            make(chan *templateJob, numWorkers*2), // Buffered channel for better performance
		results:         make(chan *templateResult, numWorkers*2),
		errors:          make(chan error, numWorkers),
		ctx:             ctx,
		cancel:          cancel,
		dirCache:        newDirCache(),
		progressCounter: &counter,
	}
}

// Start initializes and starts all workers
func (wp *WorkerPool) Start() {
	for i := 0; i < wp.numWorkers; i++ {
		wp.wg.Add(1)
		go wp.worker(i)
	}
}

// Stop gracefully shuts down the worker pool
func (wp *WorkerPool) Stop() {
	wp.cancel()
	close(wp.jobs)
	wp.wg.Wait()
	close(wp.results)
	close(wp.errors)
}

// SubmitJob adds a template processing job to the queue
func (wp *WorkerPool) SubmitJob(job *templateJob) {
	select {
	case wp.jobs <- job:
		atomic.AddInt64(&wp.totalJobs, 1)
	case <-wp.ctx.Done():
		return
	}
}

// GetResult retrieves a processed result
func (wp *WorkerPool) GetResult() *templateResult {
	select {
	case result := <-wp.results:
		return result
	case <-wp.ctx.Done():
		return nil
	}
}

// GetError retrieves any worker error
func (wp *WorkerPool) GetError() error {
	select {
	case err := <-wp.errors:
		return err
	default:
		return nil
	}
}

// GetProgress returns the current progress as a percentage
func (wp *WorkerPool) GetProgress() float64 {
	if wp.totalJobs == 0 {
		return 0
	}
	processed := atomic.LoadInt64(wp.progressCounter)
	return float64(processed) / float64(wp.totalJobs) * 100
}

// GetProcessedCount returns the number of processed jobs
func (wp *WorkerPool) GetProcessedCount() int64 {
	return atomic.LoadInt64(wp.progressCounter)
}

// GetOptimalWorkerCount returns the optimal number of workers for the current system
func GetOptimalWorkerCount() int {
	cpuCount := runtime.NumCPU()
	// Use CPU count for I/O-bound tasks like template processing
	// Cap at 8 for reasonable resource usage
	if cpuCount > 8 {
		return 8
	}
	if cpuCount < 1 {
		return 1
	}
	return cpuCount
}

// worker processes template jobs
func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()

	for {
		select {
		case job := <-wp.jobs:
			if job == nil {
				return
			}
			wp.processJob(workerID, job)
		case <-wp.ctx.Done():
			return
		}
	}
}

// processJob processes a single template job
func (wp *WorkerPool) processJob(workerID int, job *templateJob) {
	defer atomic.AddInt64(wp.progressCounter, 1)

	result := &templateResult{
		jobID: job.jobID,
	}

	// Ensure destination directory exists
	destDir := filepath.Dir(job.destPath)
	if err := wp.dirCache.ensureDir(destDir, job.transaction); err != nil {
		result.error = fmt.Errorf("worker %d: failed to create directory: %w", workerID, err)
		wp.results <- result
		return
	}

	// Create a temporary generator instance for processing
	tempGen := &Generator{
		loader: templates.NewTemplateLoader(), // Create a new instance for thread safety
		currentTransaction: job.transaction, // Pass the transaction for tracking
	}

	// Process the template file (will use cached templates)
	if err := tempGen.processTemplateFile(job.templateDir, job.templateFile.Source, job.destPath, job.context); err != nil {
		result.error = fmt.Errorf("worker %d: failed to process template %s: %w", workerID, job.templateFile.Source, err)
		wp.results <- result
		return
	}

	// Set executable permission if needed
	if job.templateFile.Executable {
		if err := os.Chmod(job.destPath, 0755); err != nil {
			result.error = fmt.Errorf("worker %d: failed to set executable permission: %w", workerID, err)
			wp.results <- result
			return
		}
	}

	result.filePath = job.destPath
	wp.results <- result
}

// ParallelTemplateProcessor handles parallel processing of multiple template files
type ParallelTemplateProcessor struct {
	pool        *WorkerPool
	transaction *GenerationTransaction
	loader      *templates.TemplateLoader
}

// NewParallelTemplateProcessor creates a new parallel template processor
func NewParallelTemplateProcessor(numWorkers int, transaction *GenerationTransaction) *ParallelTemplateProcessor {
	return &ParallelTemplateProcessor{
		pool:        NewWorkerPool(numWorkers),
		transaction: transaction,
		loader:      templates.NewTemplateLoader(),
	}
}

// ProcessTemplates processes multiple template files in parallel
func (ptp *ParallelTemplateProcessor) ProcessTemplates(
	templateFiles []types.TemplateFile,
	templateDir string,
	outputPath string,
	context map[string]any,
	config types.ProjectConfig,
	tmpl *types.Template,
) ([]string, error) {
	if len(templateFiles) == 0 {
		return []string{}, nil
	}

	// Start the worker pool
	ptp.pool.Start()
	defer ptp.pool.Stop()

	var filesCreated []string
	var mu sync.Mutex // Protect filesCreated slice

	// Submit all jobs
	jobCount := 0
	for _, templateFile := range templateFiles {
		// Process template path with variables for destination only
		destPath := ptp.processTemplatePath(templateFile.Destination, config, tmpl)
		fullDestPath := filepath.Join(outputPath, destPath)

		job := &templateJob{
			templateFile: templateFile,
			templateDir:  templateDir,
			destPath:     fullDestPath,
			context:      context,
			jobID:        jobCount,
			transaction:  ptp.transaction,
		}

		ptp.pool.SubmitJob(job)
		jobCount++
	}

	// Collect results and track directories for rollback
	directorySet := make(map[string]bool)
	
	// Collect results
	completedJobs := 0
	var firstError error
	var allErrors []error

	for completedJobs < jobCount {
		// Check for worker errors first
		if err := ptp.pool.GetError(); err != nil && firstError == nil {
			firstError = err
			allErrors = append(allErrors, err)
			ptp.pool.cancel() // Stop processing on first error
			break
		}

		// Get result
		result := ptp.pool.GetResult()
		if result == nil {
			break // Pool was stopped
		}

		completedJobs++

		// Show progress
		progress := ptp.pool.GetProgress()
		fmt.Printf("\r   [%.0f%%] Processing file %d/%d", progress, completedJobs, jobCount)

		if result.error != nil {
			if firstError == nil {
				firstError = result.error
			}
			allErrors = append(allErrors, result.error)
			ptp.pool.cancel() // Stop processing on first error
			continue
		}

		if !result.wasSkipped {
			// Thread-safe append to filesCreated
			mu.Lock()
			filesCreated = append(filesCreated, result.filePath)
			
			// Track file creation for rollback if transaction is active
			if ptp.transaction != nil {
				ptp.transaction.AddFile(result.filePath)
				
				// Track all parent directories for this file
				dir := filepath.Dir(result.filePath)
				for dir != filepath.Dir(dir) { // Stop at root
					if !directorySet[dir] {
						ptp.transaction.AddDirectory(dir)
						directorySet[dir] = true
					}
					dir = filepath.Dir(dir)
				}
			}
			mu.Unlock()
		}
	}

	// Clear the progress line and show completion
	fmt.Printf("\r\033[K") // Clear the line
	if firstError != nil {
		fmt.Printf("✗ Processing stopped due to error after %d/%d files\n", completedJobs, jobCount)
		
		// If we have multiple errors, show a summary
		if len(allErrors) > 1 {
			fmt.Printf("   Multiple errors occurred:\n")
			for i, err := range allErrors {
				if i >= 3 { // Limit to first 3 errors to avoid spam
					fmt.Printf("   ... and %d more errors\n", len(allErrors)-3)
					break
				}
				fmt.Printf("   - %v\n", err)
			}
		}
		
		return filesCreated, firstError
	}

	fmt.Printf("✓ Generated %d files successfully\n", len(filesCreated))
	return filesCreated, nil
}

// processTemplatePath processes template variables in file paths (thread-safe copy)
func (ptp *ParallelTemplateProcessor) processTemplatePath(path string, config types.ProjectConfig, tmplObj *types.Template) string {
	// This is a simplified version of the path processing logic
	// We'll reuse the main generator's method through a temporary instance
	tempGen := &Generator{}
	return tempGen.processTemplatePath(path, config, tmplObj)
}