package optimization

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOptimizationPipeline_EndToEnd(t *testing.T) {
	// Create a complete test project with various optimization opportunities
	tempDir, err := os.MkdirTemp("", "e2e-optimization-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a realistic project structure
	projectFiles := map[string]string{
		"main.go": `package main

import (
	"fmt"
	"os"     // unused
	"log"    // unused
	"github.com/gin-gonic/gin"
	"strings" // unused
)

func main() {
	fmt.Println("Starting server...")
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		fmt.Println("Hello World")
	})
	r.Run(":8080")
}`,

		"handlers/user.go": `package handlers

import (
	"fmt"
	"net/http"
	"encoding/json" // unused
	"time"         // unused
	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	fmt.Println("Getting user")
	c.JSON(http.StatusOK, gin.H{"message": "user"})
}`,

		"models/user.go": `package models

import (
	"time"
	"fmt" // unused
)

type User struct {
	ID        int       ` + "`json:\"id\"`" + `
	Name      string    ` + "`json:\"name\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
}`,

		"utils/helper.go": `package utils

import (
	"strings"
	"fmt"    // unused
	"os"     // unused
)

func ProcessString(s string) string {
	return strings.ToUpper(s)
}`,

		// Test files that should be skipped
		"handlers/user_test.go": `package handlers

import (
	"testing"
	"net/http/httptest" // unused in test - should be skipped
	"github.com/gin-gonic/gin"
)

func TestGetUser(t *testing.T) {
	r := gin.Default()
	// Test code here
}`,
	}

	// Create all test files
	for relativePath, content := range projectFiles {
		fullPath := filepath.Join(tempDir, relativePath)
		
		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err = os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		// Write file
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}

	t.Logf("Created test project in: %s", tempDir)

	// Create pipeline with comprehensive optimization
	options := PipelineOptions{
		// Enable import optimizations
		RemoveUnusedImports: true,
		OrganizeImports:     true,
		AddMissingImports:   false,

		// Conservative code optimizations (disabled for now)
		RemoveUnusedVars:    false,
		RemoveUnusedFuncs:   false,
		OptimizeConditionals: false,

		// Output options - dry run for testing
		WriteOptimizedFiles: false,
		CreateBackups:       true,
		DryRun:             true,
		Verbose:            true,

		// Performance settings
		MaxFileSize:        1024 * 1024,
		MaxConcurrentFiles: 4,
		SkipTestFiles:      true,  // Should skip user_test.go
		SkipVendorDirs:     true,

		// File patterns
		IncludePatterns: []string{"*.go"},
		ExcludePatterns: []string{"vendor/", ".git/"},
	}

	pipeline := NewOptimizationPipeline(options)

	// Optimize the entire project
	result, err := pipeline.OptimizeProject(tempDir)
	require.NoError(t, err)

	// Verify overall project statistics
	assert.Equal(t, 4, result.TotalFiles)     // 4 non-test .go files
	assert.Equal(t, 4, result.FilesProcessed) // All 4 processed
	assert.Equal(t, 4, result.FilesOptimized) // All had unused imports
	assert.Equal(t, 0, result.FilesWithErrors)
	assert.Equal(t, 0, result.FilesSkipped)   // No files skipped during processing (test files excluded during finding)

	// Verify import optimization statistics
	expectedRemovedImports := 8 // os, log, strings (main) + encoding/json, time (handlers/user) + fmt (models/user) + fmt, os (utils/helper)
	assert.Equal(t, expectedRemovedImports, result.ImportsRemoved)
	assert.Equal(t, 0, result.ImportsAdded)
	assert.Equal(t, 4, result.ImportsOrganized) // All files had imports organized

	// Verify performance metrics
	assert.Greater(t, result.SizeBeforeBytes, int64(0))
	assert.Greater(t, result.SizeAfterBytes, int64(0))
	assert.GreaterOrEqual(t, result.ProcessingTimeMs, int64(0))

	// Verify specific file optimizations
	testCases := []struct {
		file                string
		expectedRemovedImports []string
	}{
		{
			file:                filepath.Join(tempDir, "main.go"),
			expectedRemovedImports: []string{"os", "log", "strings"},
		},
		{
			file:                filepath.Join(tempDir, "handlers/user.go"),
			expectedRemovedImports: []string{"encoding/json", "time"},
		},
		{
			file:                filepath.Join(tempDir, "models/user.go"),
			expectedRemovedImports: []string{"fmt"},
		},
		{
			file:                filepath.Join(tempDir, "utils/helper.go"),
			expectedRemovedImports: []string{"fmt", "os"},
		},
	}

	for _, tc := range testCases {
		t.Run(filepath.Base(tc.file), func(t *testing.T) {
			fileResult, exists := result.FileResults[tc.file]
			require.True(t, exists, "File result should exist for %s", tc.file)

			assert.True(t, fileResult.OptimizationApplied, "Optimization should have been applied")
			assert.NotNil(t, fileResult.ImportsResult, "Import result should exist")
			assert.ElementsMatch(t, tc.expectedRemovedImports, fileResult.ImportsResult.RemovedImports,
				"Removed imports should match expected for %s", tc.file)
			assert.Empty(t, fileResult.ImportsResult.AddedImports, "No imports should have been added")
			assert.Empty(t, fileResult.Errors, "Should have no errors")
			assert.Greater(t, fileResult.OriginalSize, int64(0), "Should have original size")
		})
	}

	// Verify test file was excluded (not in results)
	testFile := filepath.Join(tempDir, "handlers/user_test.go")
	_, exists := result.FileResults[testFile]
	assert.False(t, exists, "Test file should have been excluded from processing")

	t.Logf("End-to-end test completed successfully!")
	t.Logf("Processed %d files, optimized %d files, removed %d imports", 
		result.FilesProcessed, result.FilesOptimized, result.ImportsRemoved)
}

func TestPipelineStatistics_Calculation(t *testing.T) {
	// Test the statistics calculation with known data
	tempDir, err := os.MkdirTemp("", "stats-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create files with known sizes and optimization potential
	testFiles := map[string]string{
		"small.go": `package main
import "fmt"
func main() { fmt.Println("hello") }`,

		"medium.go": `package main

import (
	"fmt"
	"os"      // unused - will be removed
	"strings" // unused - will be removed  
)

func main() {
	fmt.Println("Hello World")
}`,

		"large.go": `package main

import (
	"fmt"
	"log"     // unused
	"os"      // unused
	"strings" // unused
	"time"    // unused
	"net/http" // unused
)

func main() {
	fmt.Println("This is a larger file with more content")
	fmt.Println("It has multiple lines and more unused imports")
	fmt.Println("This helps test the size reduction calculation")
}`,
	}

	var originalTotalSize int64
	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
		
		// Calculate original size
		originalTotalSize += int64(len(content))
	}

	// Run optimization
	options := DefaultPipelineOptions()
	options.RemoveUnusedImports = true
	options.OrganizeImports = true
	options.DryRun = true
	options.Verbose = false // Reduce test output

	pipeline := NewOptimizationPipeline(options)
	result, err := pipeline.OptimizeProject(tempDir)
	require.NoError(t, err)

	// Verify statistics make sense
	assert.Equal(t, 3, result.TotalFiles)
	assert.Equal(t, 3, result.FilesProcessed)
	assert.Equal(t, 2, result.FilesOptimized)  // small.go has no unused imports
	assert.Equal(t, 0, result.FilesWithErrors)

	// medium.go should remove 2 imports, large.go should remove 5 imports
	assert.Equal(t, 7, result.ImportsRemoved)
	assert.Equal(t, 0, result.ImportsAdded)

	// Size calculations - in dry run mode, sizes are not actually changed
	assert.Equal(t, originalTotalSize, result.SizeBeforeBytes)
	assert.Greater(t, result.SizeAfterBytes, int64(0))
	
	// In dry run mode, the files aren't actually optimized, so sizes remain the same
	// But we can verify that optimizations were detected
	assert.GreaterOrEqual(t, result.SizeAfterBytes, int64(0))
	assert.GreaterOrEqual(t, result.SizeReductionBytes, int64(0))

	// Performance metrics
	assert.GreaterOrEqual(t, result.ProcessingTimeMs, int64(0))

	t.Logf("Size reduction: %d bytes (%.2f%% reduction)",
		result.SizeReductionBytes,
		float64(result.SizeReductionBytes)/float64(result.SizeBeforeBytes)*100)
}

func TestPipelineConfiguration_Scenarios(t *testing.T) {
	// Test different pipeline configuration scenarios
	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testCode := `package main

import (
	"fmt"
	"os"     // unused
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Hello")
	gin.New()
}`

	testFile := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	scenarios := []struct {
		name                  string
		removeUnusedImports   bool
		organizeImports       bool
		expectedRemovedCount  int
		expectedOptimization  bool
	}{
		{
			name:                  "imports removal enabled",
			removeUnusedImports:   true,
			organizeImports:       false,
			expectedRemovedCount:  1, // "os"
			expectedOptimization:  true,
		},
		{
			name:                  "only organization enabled",
			removeUnusedImports:   false,
			organizeImports:       true,
			expectedRemovedCount:  0,
			expectedOptimization:  false, // Organization alone doesn't count without actual changes
		},
		{
			name:                  "both optimizations enabled",
			removeUnusedImports:   true,
			organizeImports:       true,
			expectedRemovedCount:  1,
			expectedOptimization:  true,
		},
		{
			name:                  "all optimizations disabled",
			removeUnusedImports:   false,
			organizeImports:       false,
			expectedRemovedCount:  0,
			expectedOptimization:  false,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			options := DefaultPipelineOptions()
			options.RemoveUnusedImports = scenario.removeUnusedImports
			options.OrganizeImports = scenario.organizeImports
			options.DryRun = true
			options.Verbose = false

			pipeline := NewOptimizationPipeline(options)
			result, err := pipeline.OptimizeProject(tempDir)
			require.NoError(t, err)

			assert.Equal(t, scenario.expectedRemovedCount, result.ImportsRemoved)
			
			if scenario.expectedOptimization {
				assert.Equal(t, 1, result.FilesOptimized)
			} else {
				assert.Equal(t, 0, result.FilesOptimized)
			}
		})
	}
}