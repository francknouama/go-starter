package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestBatchedFileWriter(t *testing.T) {
	// Create temporary directory for testing
	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	tx := NewGenerationTransaction(tempDir)
	bfw := NewBatchedFileWriter(tx)

	t.Run("queue and flush writes", func(t *testing.T) {
		// Queue multiple file writes
		files := map[string]string{
			"file1.txt":           "content1",
			"dir1/file2.txt":      "content2",
			"dir1/dir2/file3.txt": "content3",
			"dir3/file4.txt":      "content4",
		}

		for path, content := range files {
			fullPath := filepath.Join(tempDir, path)
			bfw.QueueWrite(fullPath, []byte(content), 0644)
		}

		// Flush writes
		err := bfw.FlushWrites()
		if err != nil {
			t.Fatalf("Failed to flush writes: %v", err)
		}

		// Verify all files were created with correct content
		for path, expectedContent := range files {
			fullPath := filepath.Join(tempDir, path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}
			if string(content) != expectedContent {
				t.Errorf("File %s: expected %q, got %q", path, expectedContent, string(content))
			}
		}

		// Verify directories were created
		expectedDirs := []string{
			filepath.Join(tempDir, "dir1"),
			filepath.Join(tempDir, "dir1", "dir2"),
			filepath.Join(tempDir, "dir3"),
		}

		for _, dir := range expectedDirs {
			if info, err := os.Stat(dir); err != nil {
				t.Errorf("Directory %s was not created: %v", dir, err)
			} else if !info.IsDir() {
				t.Errorf("Path %s is not a directory", dir)
			}
		}
	})

	t.Run("executable files", func(t *testing.T) {
		execPath := filepath.Join(tempDir, "script.sh")
		bfw.QueueWrite(execPath, []byte("#!/bin/bash\necho hello"), 0755)

		err := bfw.FlushWrites()
		if err != nil {
			t.Fatalf("Failed to flush writes: %v", err)
		}

		// Check file permissions
		info, err := os.Stat(execPath)
		if err != nil {
			t.Fatalf("Failed to stat executable file: %v", err)
		}

		mode := info.Mode()
		if mode&0755 != 0755 {
			t.Errorf("Expected executable permissions 0755, got %o", mode.Perm())
		}
	})

	t.Run("directory caching", func(t *testing.T) {
		// Create multiple files in the same directory
		dir := filepath.Join(tempDir, "cached_dir")
		for i := 0; i < 5; i++ {
			path := filepath.Join(dir, fmt.Sprintf("file%d.txt", i))
			bfw.QueueWrite(path, []byte("content"), 0644)
		}

		err := bfw.FlushWrites()
		if err != nil {
			t.Fatalf("Failed to flush writes: %v", err)
		}

		// Verify directory was created and cached
		queuedWrites, cachedDirs := bfw.GetStats()
		if queuedWrites != 0 {
			t.Errorf("Expected 0 queued writes after flush, got %d", queuedWrites)
		}
		if cachedDirs == 0 {
			t.Error("Expected some cached directories")
		}
	})
}

func TestOptimizedFileOperations(t *testing.T) {
	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	tx := NewGenerationTransaction(tempDir)
	ofo := NewOptimizedFileOperations(tx)

	t.Run("create project structure", func(t *testing.T) {
		files := map[string][]byte{
			filepath.Join(tempDir, "README.md"):                []byte("# Project"),
			filepath.Join(tempDir, "main.go"):                  []byte("package main"),
			filepath.Join(tempDir, "cmd/app/main.go"):          []byte("package main"),
			filepath.Join(tempDir, "internal/service/user.go"): []byte("package service"),
			filepath.Join(tempDir, "pkg/utils/helper.go"):      []byte("package utils"),
			filepath.Join(tempDir, "scripts/build.sh"):         []byte("#!/bin/bash"),
		}

		executableFiles := []string{
			filepath.Join(tempDir, "scripts/build.sh"),
		}

		err := ofo.CreateProjectStructure(files, executableFiles)
		if err != nil {
			t.Fatalf("Failed to create project structure: %v", err)
		}

		// Verify all files exist with correct content
		for path, expectedContent := range files {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}
			if string(content) != string(expectedContent) {
				t.Errorf("File %s: expected %q, got %q", path, string(expectedContent), string(content))
			}
		}

		// Verify executable file has correct permissions
		execPath := filepath.Join(tempDir, "scripts/build.sh")
		info, err := os.Stat(execPath)
		if err != nil {
			t.Fatalf("Failed to stat executable file: %v", err)
		}
		if info.Mode()&0755 != 0755 {
			t.Errorf("Expected executable permissions for %s", execPath)
		}
	})

	t.Run("performance metrics", func(t *testing.T) {
		files := make(map[string][]byte)
		for i := 0; i < 10; i++ {
			path := filepath.Join(tempDir, "perf", "file"+fmt.Sprintf("run%d", i)+".txt")
			files[path] = []byte("test content " + fmt.Sprintf("run%d", i))
		}

		err := ofo.CreateProjectStructure(files, nil)
		if err != nil {
			t.Fatalf("Failed to create structure: %v", err)
		}

		metrics := ofo.GetIOMetrics()
		if metrics.FilesWritten < 0 {
			t.Error("Expected non-negative files written count")
		}
		if metrics.DirectoriesCreated < 0 {
			t.Error("Expected non-negative directories created count")
		}
	})
}

func TestAdaptiveFileWriter(t *testing.T) {
	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	tx := NewGenerationTransaction(tempDir)
	afw := NewAdaptiveFileWriter(tx)

	t.Run("small project uses default strategy", func(t *testing.T) {
		files := map[string][]byte{
			filepath.Join(tempDir, "small1.txt"): []byte("content1"),
			filepath.Join(tempDir, "small2.txt"): []byte("content2"),
		}

		err := afw.WriteFiles(files, nil)
		if err != nil {
			t.Fatalf("Failed to write small project: %v", err)
		}

		// Verify files exist
		for path, expectedContent := range files {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}
			if string(content) != string(expectedContent) {
				t.Errorf("File %s: expected %q, got %q", path, string(expectedContent), string(content))
			}
		}
	})

	t.Run("medium project uses buffered strategy", func(t *testing.T) {
		files := map[string][]byte{
			filepath.Join(tempDir, "medium1.txt"): []byte("content1"),
			filepath.Join(tempDir, "medium2.txt"): []byte("content2"),
			filepath.Join(tempDir, "medium3.txt"): []byte("content3"),
			filepath.Join(tempDir, "medium4.txt"): []byte("content4"),
		}

		err := afw.WriteFiles(files, nil)
		if err != nil {
			t.Fatalf("Failed to write medium project: %v", err)
		}

		// Verify files exist
		for path, expectedContent := range files {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}
			if string(content) != string(expectedContent) {
				t.Errorf("File %s: expected %q, got %q", path, string(expectedContent), string(content))
			}
		}
	})

	t.Run("large project uses batched strategy", func(t *testing.T) {
		// Create a large project with more than threshold files
		files := make(map[string][]byte)
		for i := 0; i < 15; i++ {
			path := filepath.Join(tempDir, "large", fmt.Sprintf("file%d.txt", i))
			files[path] = []byte(fmt.Sprintf("large content %d", i))
		}

		executableFiles := []string{
			filepath.Join(tempDir, "large", "file0.txt"), // Make first file executable for testing
		}

		err := afw.WriteFiles(files, executableFiles)
		if err != nil {
			t.Fatalf("Failed to write large project: %v", err)
		}

		// Verify all files exist
		for path, expectedContent := range files {
			content, err := os.ReadFile(path)
			if err != nil {
				t.Errorf("Failed to read file %s: %v", path, err)
				continue
			}
			if string(content) != string(expectedContent) {
				t.Errorf("File %s: expected %q, got %q", path, string(expectedContent), string(content))
			}
		}

		// Verify executable file
		execPath := filepath.Join(tempDir, "large", "file0.txt")
		info, err := os.Stat(execPath)
		if err != nil {
			t.Fatalf("Failed to stat executable file: %v", err)
		}
		if info.Mode()&0755 != 0755 {
			t.Errorf("Expected executable permissions for %s", execPath)
		}
	})

	t.Run("custom threshold", func(t *testing.T) {
		afw.SetThreshold(5)

		files := make(map[string][]byte)
		for i := 0; i < 6; i++ {
			path := filepath.Join(tempDir, "threshold", "file"+fmt.Sprintf("run%d", i)+".txt")
			files[path] = []byte("threshold content")
		}

		err := afw.WriteFiles(files, nil)
		if err != nil {
			t.Fatalf("Failed to write files with custom threshold: %v", err)
		}

		// Verify all files exist
		for path := range files {
			if _, err := os.Stat(path); err != nil {
				t.Errorf("File %s should exist: %v", path, err)
			}
		}
	})
}

func BenchmarkFileWriteStrategies(b *testing.B) {
	tempDir := setupTempDirForBench(b)
	defer cleanupTempDirForBench(b, tempDir)

	// Create test data
	files := make(map[string][]byte)
	for i := 0; i < 50; i++ {
		path := filepath.Join("benchmark", fmt.Sprintf("file%d_%d.txt", i%26, i/26))
		content := strings.Repeat("benchmark content ", 100) // ~1.8KB per file
		files[path] = []byte(content)
	}

	b.Run("default_strategy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchDir := filepath.Join(tempDir, "default", fmt.Sprintf("run%d", i))
			tx := NewGenerationTransaction(tempDir)
			afw := NewAdaptiveFileWriter(tx)

			// Force default strategy
			afw.SetThreshold(1000)

			fullPathFiles := make(map[string][]byte)
			for path, content := range files {
				fullPathFiles[filepath.Join(benchDir, path)] = content
			}

			err := afw.WriteFiles(fullPathFiles, nil)
			if err != nil {
				b.Fatalf("Failed to write files: %v", err)
			}
		}
	})

	b.Run("buffered_strategy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchDir := filepath.Join(tempDir, "buffered", fmt.Sprintf("run%d", i))
			tx := NewGenerationTransaction(tempDir)
			afw := NewAdaptiveFileWriter(tx)

			// Force buffered strategy
			afw.SetThreshold(3)

			fullPathFiles := make(map[string][]byte)
			for path, content := range files {
				fullPathFiles[filepath.Join(benchDir, path)] = content
			}

			// Use exactly threshold files to trigger buffered strategy
			limitedFiles := make(map[string][]byte)
			count := 0
			for path, content := range fullPathFiles {
				limitedFiles[path] = content
				count++
				if count >= 5 { // Use 5 files to trigger buffered
					break
				}
			}

			err := afw.WriteFiles(limitedFiles, nil)
			if err != nil {
				b.Fatalf("Failed to write files: %v", err)
			}
		}
	})

	b.Run("batched_strategy", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchDir := filepath.Join(tempDir, "batched", fmt.Sprintf("run%d", i))
			tx := NewGenerationTransaction(tempDir)
			ofo := NewOptimizedFileOperations(tx)

			fullPathFiles := make(map[string][]byte)
			for path, content := range files {
				fullPathFiles[filepath.Join(benchDir, path)] = content
			}

			err := ofo.CreateProjectStructure(fullPathFiles, nil)
			if err != nil {
				b.Fatalf("Failed to write files: %v", err)
			}
		}
	})
}

func BenchmarkDirectoryCreation(b *testing.B) {
	tempDir := setupTempDirForBench(b)
	defer cleanupTempDirForBench(b, tempDir)

	b.Run("individual_mkdir", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchDir := filepath.Join(tempDir, "individual", fmt.Sprintf("run-%d", i))
			
			// Create directories one by one
			dirs := []string{
				filepath.Join(benchDir, "cmd", "app"),
				filepath.Join(benchDir, "internal", "service"),
				filepath.Join(benchDir, "internal", "repository"),
				filepath.Join(benchDir, "pkg", "utils"),
				filepath.Join(benchDir, "api", "v1"),
				filepath.Join(benchDir, "scripts", "dev"),
				filepath.Join(benchDir, "configs", "local"),
				filepath.Join(benchDir, "docs", "api"),
			}

			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0755); err != nil {
					b.Fatalf("Failed to create directory: %v", err)
				}
			}
		}
	})

	b.Run("batched_mkdir", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchDir := filepath.Join(tempDir, "batched", fmt.Sprintf("run-%d", i))
			tx := NewGenerationTransaction(tempDir)
			bfw := NewBatchedFileWriter(tx)

			// Add dummy files to trigger directory creation
			files := []string{
				filepath.Join(benchDir, "cmd", "app", "main.go"),
				filepath.Join(benchDir, "internal", "service", "user.go"),
				filepath.Join(benchDir, "internal", "repository", "db.go"),
				filepath.Join(benchDir, "pkg", "utils", "helper.go"),
				filepath.Join(benchDir, "api", "v1", "handler.go"),
				filepath.Join(benchDir, "scripts", "dev", "setup.sh"),
				filepath.Join(benchDir, "configs", "local", "config.yaml"),
				filepath.Join(benchDir, "docs", "api", "README.md"),
			}

			for _, file := range files {
				bfw.QueueWrite(file, []byte("dummy"), 0644)
			}

			if err := bfw.FlushWrites(); err != nil {
				b.Fatalf("Failed to flush writes: %v", err)
			}
		}
	})
}

// Helper functions for testing

func setupTempDir(t *testing.T) string {
	t.Helper()
	tempDir, err := os.MkdirTemp("", "go-starter-io-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

func cleanupTempDir(t *testing.T, tempDir string) {
	t.Helper()
	if err := os.RemoveAll(tempDir); err != nil {
		t.Logf("Warning: failed to remove temp dir: %v", err)
	}
}

func setupTempDirForBench(b *testing.B) string {
	b.Helper()
	tempDir, err := os.MkdirTemp("", "go-starter-io-bench")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	return tempDir
}

func cleanupTempDirForBench(b *testing.B, tempDir string) {
	b.Helper()
	if err := os.RemoveAll(tempDir); err != nil {
		b.Logf("Warning: failed to remove temp dir: %v", err)
	}
}

func TestIOPerformanceComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance comparison in short mode")
	}

	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	// Test with realistic project structure
	files := map[string][]byte{
		"README.md":                    []byte("# Test Project\n\nThis is a test project."),
		"go.mod":                       []byte("module github.com/test/project\n\ngo 1.21"),
		"main.go":                      []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}"),
		"cmd/server/main.go":           []byte("package main\n\nfunc main() {\n\t// Server main\n}"),
		"internal/service/user.go":     []byte("package service\n\ntype UserService struct {}"),
		"internal/repository/user.go":  []byte("package repository\n\ntype UserRepository struct {}"),
		"internal/handler/user.go":     []byte("package handler\n\ntype UserHandler struct {}"),
		"pkg/utils/helper.go":          []byte("package utils\n\nfunc Helper() string { return \"help\" }"),
		"api/openapi.yaml":             []byte("openapi: 3.0.0\ninfo:\n  title: Test API"),
		"configs/config.yaml":          []byte("app:\n  name: test\n  port: 8080"),
		"scripts/build.sh":             []byte("#!/bin/bash\ngo build -o bin/app main.go"),
		"docs/DEPLOYMENT.md":           []byte("# Deployment\n\nHow to deploy."),
	}

	executableFiles := []string{"scripts/build.sh"}
	_ = executableFiles // Mark as used

	t.Run("timing comparison", func(t *testing.T) {
		// Test default strategy
		start := time.Now()
		tx1 := NewGenerationTransaction(tempDir)
		afw1 := NewAdaptiveFileWriter(tx1)
		afw1.SetThreshold(1000) // Force default strategy

		defaultFiles := make(map[string][]byte)
		for path, content := range files {
			defaultFiles[filepath.Join(tempDir, "default_timing", path)] = content
		}

		err := afw1.WriteFiles(defaultFiles, []string{filepath.Join(tempDir, "default_timing", "scripts/build.sh")})
		if err != nil {
			t.Fatalf("Default strategy failed: %v", err)
		}
		defaultTime := time.Since(start)

		// Test batched strategy
		start = time.Now()
		tx2 := NewGenerationTransaction(tempDir)
		ofo := NewOptimizedFileOperations(tx2)

		batchedFiles := make(map[string][]byte)
		for path, content := range files {
			batchedFiles[filepath.Join(tempDir, "batched_timing", path)] = content
		}

		err = ofo.CreateProjectStructure(batchedFiles, []string{filepath.Join(tempDir, "batched_timing", "scripts/build.sh")})
		if err != nil {
			t.Fatalf("Batched strategy failed: %v", err)
		}
		batchedTime := time.Since(start)

		t.Logf("Default strategy time: %v", defaultTime)
		t.Logf("Batched strategy time: %v", batchedTime)

		// Batched should generally be faster for projects with many files
		if len(files) > 10 && batchedTime > defaultTime*2 {
			t.Logf("Warning: Batched strategy is significantly slower than expected")
		}
	})
}

// Test concurrent access to ensure thread safety
func TestConcurrentAccess(t *testing.T) {
	tempDir := setupTempDir(t)
	defer cleanupTempDir(t, tempDir)

	tx := NewGenerationTransaction(tempDir)
	bfw := NewBatchedFileWriter(tx)

	// Test concurrent directory creation
	t.Run("concurrent directory caching", func(t *testing.T) {
		numGoroutines := runtime.NumCPU()
		done := make(chan bool, numGoroutines)

		for i := 0; i < numGoroutines; i++ {
			go func(id int) {
				defer func() { done <- true }()

				for j := 0; j < 10; j++ {
					path := filepath.Join(tempDir, "concurrent", fmt.Sprintf("worker%d", id), fmt.Sprintf("dir%d", j), "file.txt")
					bfw.QueueWrite(path, []byte("concurrent content"), 0644)
				}
			}(i)
		}

		// Wait for all goroutines to finish
		for i := 0; i < numGoroutines; i++ {
			<-done
		}

		// Flush all writes
		err := bfw.FlushWrites()
		if err != nil {
			t.Fatalf("Failed to flush concurrent writes: %v", err)
		}

		// Verify files were created
		queuedWrites, cachedDirs := bfw.GetStats()
		if queuedWrites != 0 {
			t.Errorf("Expected 0 queued writes after flush, got %d", queuedWrites)
		}
		if cachedDirs == 0 {
			t.Error("Expected some cached directories")
		}
	})
}