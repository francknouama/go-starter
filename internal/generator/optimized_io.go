package generator

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
)

// BatchedFileWriter provides optimized file I/O operations for project generation
type BatchedFileWriter struct {
	dirCache     map[string]bool
	dirCacheMu   sync.RWMutex
	writeJobs    []WriteJob
	transaction  *GenerationTransaction
	bufferSize   int
}

// WriteJob represents a file write operation to be batched
type WriteJob struct {
	Path    string
	Content []byte
	Mode    os.FileMode
}

// NewBatchedFileWriter creates a new batched file writer
func NewBatchedFileWriter(transaction *GenerationTransaction) *BatchedFileWriter {
	return &BatchedFileWriter{
		dirCache:    make(map[string]bool),
		writeJobs:   make([]WriteJob, 0),
		transaction: transaction,
		bufferSize:  64 * 1024, // 64KB buffer for file operations
	}
}

// QueueWrite adds a file write operation to the batch
func (bfw *BatchedFileWriter) QueueWrite(path string, content []byte, mode os.FileMode) {
	bfw.writeJobs = append(bfw.writeJobs, WriteJob{
		Path:    path,
		Content: content,
		Mode:    mode,
	})
}

// FlushWrites executes all queued write operations with optimizations
func (bfw *BatchedFileWriter) FlushWrites() error {
	if len(bfw.writeJobs) == 0 {
		return nil
	}

	// Step 1: Batch create all directories first
	if err := bfw.batchCreateDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Step 2: Sort writes by directory to improve filesystem locality
	bfw.sortWritesByDirectory()

	// Step 3: Execute writes with buffered I/O
	for _, job := range bfw.writeJobs {
		if err := bfw.writeFileBuffered(job.Path, job.Content, job.Mode); err != nil {
			return fmt.Errorf("failed to write file %s: %w", job.Path, err)
		}
		
		// Track file for rollback
		if bfw.transaction != nil {
			bfw.transaction.AddFile(job.Path)
		}
	}

	// Clear the batch
	bfw.writeJobs = bfw.writeJobs[:0]
	return nil
}

// batchCreateDirectories creates all required directories in a single pass
func (bfw *BatchedFileWriter) batchCreateDirectories() error {
	// Collect all unique directories
	dirSet := make(map[string]bool)
	for _, job := range bfw.writeJobs {
		dir := filepath.Dir(job.Path)
		
		// Add all parent directories
		for {
			if dir == "" || dir == "." || dir == "/" {
				break
			}
			dirSet[dir] = true
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}

	// Convert to sorted slice for consistent creation order
	dirs := make([]string, 0, len(dirSet))
	for dir := range dirSet {
		dirs = append(dirs, dir)
	}
	sort.Strings(dirs)

	// Create directories in order (deepest first is handled by sorting)
	for _, dir := range dirs {
		if err := bfw.ensureDirectoryExists(dir); err != nil {
			return err
		}
	}

	return nil
}

// ensureDirectoryExists creates a directory if it doesn't exist, with caching
func (bfw *BatchedFileWriter) ensureDirectoryExists(path string) error {
	// Check cache first
	bfw.dirCacheMu.RLock()
	if bfw.dirCache[path] {
		bfw.dirCacheMu.RUnlock()
		return nil
	}
	bfw.dirCacheMu.RUnlock()

	// Check if directory exists
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		bfw.dirCacheMu.Lock()
		bfw.dirCache[path] = true
		bfw.dirCacheMu.Unlock()
		return nil
	}

	// Create directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", path, err)
	}

	// Cache the creation and track for rollback
	bfw.dirCacheMu.Lock()
	bfw.dirCache[path] = true
	if bfw.transaction != nil {
		bfw.transaction.AddDirectory(path)
	}
	bfw.dirCacheMu.Unlock()

	return nil
}

// sortWritesByDirectory sorts write operations by directory for better filesystem locality
func (bfw *BatchedFileWriter) sortWritesByDirectory() {
	sort.Slice(bfw.writeJobs, func(i, j int) bool {
		dirI := filepath.Dir(bfw.writeJobs[i].Path)
		dirJ := filepath.Dir(bfw.writeJobs[j].Path)
		
		if dirI != dirJ {
			return dirI < dirJ
		}
		
		// Within same directory, sort by filename
		return filepath.Base(bfw.writeJobs[i].Path) < filepath.Base(bfw.writeJobs[j].Path)
	})
}

// writeFileBuffered writes a file using buffered I/O for better performance
func (bfw *BatchedFileWriter) writeFileBuffered(path string, content []byte, mode os.FileMode) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use buffered writer for better performance
	writer := bufio.NewWriterSize(file, bfw.bufferSize)
	
	if _, err := writer.Write(content); err != nil {
		return err
	}
	
	if err := writer.Flush(); err != nil {
		return err
	}
	
	return file.Sync() // Ensure data is written to disk
}

// SetBufferSize sets the buffer size for file operations
func (bfw *BatchedFileWriter) SetBufferSize(size int) {
	if size > 0 {
		bfw.bufferSize = size
	}
}

// GetStats returns statistics about the batched operations
func (bfw *BatchedFileWriter) GetStats() (queuedWrites int, cachedDirs int) {
	bfw.dirCacheMu.RLock()
	defer bfw.dirCacheMu.RUnlock()
	return len(bfw.writeJobs), len(bfw.dirCache)
}

// OptimizedFileOperations provides high-level optimized file operations
type OptimizedFileOperations struct {
	batchWriter *BatchedFileWriter
}

// NewOptimizedFileOperations creates a new optimized file operations instance
func NewOptimizedFileOperations(transaction *GenerationTransaction) *OptimizedFileOperations {
	return &OptimizedFileOperations{
		batchWriter: NewBatchedFileWriter(transaction),
	}
}

// WriteFileOptimized queues a file for optimized writing
func (ofo *OptimizedFileOperations) WriteFileOptimized(path string, content []byte) error {
	ofo.batchWriter.QueueWrite(path, content, 0644)
	return nil
}

// WriteExecutableOptimized queues an executable file for optimized writing
func (ofo *OptimizedFileOperations) WriteExecutableOptimized(path string, content []byte) error {
	ofo.batchWriter.QueueWrite(path, content, 0755)
	return nil
}

// FlushAll executes all queued operations
func (ofo *OptimizedFileOperations) FlushAll() error {
	return ofo.batchWriter.FlushWrites()
}

// CreateProjectStructure optimally creates a complete project structure
func (ofo *OptimizedFileOperations) CreateProjectStructure(files map[string][]byte, executableFiles []string) error {
	// Mark executable files
	execSet := make(map[string]bool)
	for _, path := range executableFiles {
		execSet[path] = true
	}

	// Queue all file writes
	for path, content := range files {
		if execSet[path] {
			ofo.WriteExecutableOptimized(path, content)
		} else {
			ofo.WriteFileOptimized(path, content)
		}
	}

	// Execute all writes
	return ofo.FlushAll()
}

// Performance monitoring
type IOPerformanceMetrics struct {
	FilesWritten     int
	DirectoriesCreated int
	TotalBytes       int64
	WriteTime        int64 // in nanoseconds
	DirectoryTime    int64 // in nanoseconds
}

// GetIOMetrics returns performance metrics for the last operation
func (ofo *OptimizedFileOperations) GetIOMetrics() IOPerformanceMetrics {
	queuedWrites, cachedDirs := ofo.batchWriter.GetStats()
	
	return IOPerformanceMetrics{
		FilesWritten:       queuedWrites,
		DirectoriesCreated: cachedDirs,
		// TODO: Add timing metrics in future iterations
	}
}

// FileWriteStrategy defines different strategies for file writing
type FileWriteStrategy int

const (
	// StrategyDefault uses standard os.WriteFile
	StrategyDefault FileWriteStrategy = iota
	// StrategyBuffered uses buffered I/O
	StrategyBuffered
	// StrategyBatched uses batched operations
	StrategyBatched
)

// AdaptiveFileWriter adapts the write strategy based on workload
type AdaptiveFileWriter struct {
	strategy    FileWriteStrategy
	threshold   int // files threshold for switching to batched mode
	transaction *GenerationTransaction
}

// NewAdaptiveFileWriter creates a new adaptive file writer
func NewAdaptiveFileWriter(transaction *GenerationTransaction) *AdaptiveFileWriter {
	return &AdaptiveFileWriter{
		strategy:    StrategyDefault,
		threshold:   10, // Switch to batched mode for 10+ files
		transaction: transaction,
	}
}

// WriteFiles writes files using the optimal strategy for the workload
func (afw *AdaptiveFileWriter) WriteFiles(files map[string][]byte, executableFiles []string) error {
	fileCount := len(files)
	
	// Choose strategy based on file count
	if fileCount >= afw.threshold {
		// Use batched strategy for large projects
		ofo := NewOptimizedFileOperations(afw.transaction)
		return ofo.CreateProjectStructure(files, executableFiles)
	} else if fileCount >= 3 {
		// Use buffered strategy for medium projects
		return afw.writeFilesBuffered(files, executableFiles)
	} else {
		// Use default strategy for small projects
		return afw.writeFilesDefault(files, executableFiles)
	}
}

// writeFilesBuffered writes files using buffered I/O
func (afw *AdaptiveFileWriter) writeFilesBuffered(files map[string][]byte, executableFiles []string) error {
	execSet := make(map[string]bool)
	for _, path := range executableFiles {
		execSet[path] = true
	}

	// Create directories first
	dirs := make(map[string]bool)
	for path := range files {
		dir := filepath.Dir(path)
		dirs[dir] = true
	}

	for dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		if afw.transaction != nil {
			afw.transaction.AddDirectory(dir)
		}
	}

	// Write files with buffering
	for path, content := range files {
		mode := os.FileMode(0644)
		if execSet[path] {
			mode = 0755
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", path, err)
		}

		writer := bufio.NewWriter(file)
		if _, err := writer.Write(content); err != nil {
			file.Close()
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}

		if err := writer.Flush(); err != nil {
			file.Close()
			return fmt.Errorf("failed to flush file %s: %w", path, err)
		}

		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close file %s: %w", path, err)
		}

		if afw.transaction != nil {
			afw.transaction.AddFile(path)
		}
	}

	return nil
}

// writeFilesDefault writes files using standard os.WriteFile
func (afw *AdaptiveFileWriter) writeFilesDefault(files map[string][]byte, executableFiles []string) error {
	execSet := make(map[string]bool)
	for _, path := range executableFiles {
		execSet[path] = true
	}

	for path, content := range files {
		// Create directory if needed
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		if afw.transaction != nil {
			afw.transaction.AddDirectory(dir)
		}

		// Write file
		mode := os.FileMode(0644)
		if execSet[path] {
			mode = 0755
		}

		if err := os.WriteFile(path, content, mode); err != nil {
			return fmt.Errorf("failed to write file %s: %w", path, err)
		}

		if afw.transaction != nil {
			afw.transaction.AddFile(path)
		}
	}

	return nil
}

// SetThreshold sets the file count threshold for strategy switching
func (afw *AdaptiveFileWriter) SetThreshold(threshold int) {
	if threshold > 0 {
		afw.threshold = threshold
	}
}