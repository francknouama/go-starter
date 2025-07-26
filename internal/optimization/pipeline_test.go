package optimization

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOptimizationPipeline(t *testing.T) {
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	assert.NotNil(t, pipeline)
	assert.NotNil(t, pipeline.analyzer)
	assert.NotNil(t, pipeline.manager)
	assert.Equal(t, options, pipeline.options)
}

func TestDefaultPipelineOptions(t *testing.T) {
	options := DefaultPipelineOptions()

	// Check conservative defaults
	assert.True(t, options.RemoveUnusedImports)
	assert.True(t, options.OrganizeImports)
	assert.False(t, options.AddMissingImports) // Risky default
	assert.False(t, options.RemoveUnusedVars) // Conservative
	assert.False(t, options.RemoveUnusedFuncs) // Conservative
	assert.False(t, options.WriteOptimizedFiles) // Safe default
	assert.True(t, options.CreateBackups)
	assert.True(t, options.DryRun) // Safe default
	assert.True(t, options.SkipTestFiles)
	assert.True(t, options.SkipVendorDirs)
}

func TestOptimizeFile_SingleFile(t *testing.T) {
	// Create a temporary file with optimization opportunities
	tempDir, err := os.MkdirTemp("", "pipeline-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testCode := `package main

import (
	"fmt"
	"os"     // unused
	"strings" // unused
)

func main() {
	fmt.Println("Hello World")
}`

	testFile := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Create pipeline with optimization enabled
	options := DefaultPipelineOptions()
	options.RemoveUnusedImports = true
	options.OrganizeImports = true
	options.DryRun = true // Don't write files in test

	pipeline := NewOptimizationPipeline(options)

	// Optimize the file
	result, err := pipeline.OptimizeFile(testFile)
	require.NoError(t, err)

	// Verify optimization results
	assert.True(t, result.OptimizationApplied)
	assert.NotNil(t, result.ImportsResult)
	assert.ElementsMatch(t, []string{"os", "strings"}, result.ImportsResult.RemovedImports)
	assert.Empty(t, result.ImportsResult.AddedImports)
	assert.Greater(t, result.OriginalSize, int64(0))
	assert.Empty(t, result.Errors)
}

func TestOptimizeFile_NoOptimizationNeeded(t *testing.T) {
	// Create a temporary file that doesn't need optimization
	tempDir, err := os.MkdirTemp("", "pipeline-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
}`

	testFile := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Create pipeline
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	// Optimize the file
	result, err := pipeline.OptimizeFile(testFile)
	require.NoError(t, err)

	// Verify no optimization was needed
	assert.False(t, result.OptimizationApplied)
	assert.NotNil(t, result.ImportsResult)
	assert.Empty(t, result.ImportsResult.RemovedImports)
	assert.Empty(t, result.ImportsResult.AddedImports)
	assert.Empty(t, result.Errors)
}

func TestOptimizeProject_MultipleFiles(t *testing.T) {
	// Create a temporary project directory with multiple files
	tempDir, err := os.MkdirTemp("", "pipeline-project-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"main.go": `package main

import (
	"fmt"
	"os"     // unused
)

func main() {
	fmt.Println("Hello from main")
}`,
		"helper.go": `package main

import (
	"fmt"
	"strings" // unused
)

func helper() {
	fmt.Println("Helper function")
}`,
		"utils_test.go": `package main

import (
	"testing"
	"os" // unused in test file
)

func TestSomething(t *testing.T) {
	// Test code
}`,
	}

	for filename, content := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Create pipeline with optimization enabled
	options := DefaultPipelineOptions()
	options.RemoveUnusedImports = true
	options.OrganizeImports = true
	options.SkipTestFiles = true // Should skip utils_test.go
	options.DryRun = true
	options.Verbose = true

	pipeline := NewOptimizationPipeline(options)

	// Optimize the project
	result, err := pipeline.OptimizeProject(tempDir)
	require.NoError(t, err)

	// Verify project optimization results
	assert.Equal(t, 2, result.TotalFiles) // Found 2 non-test files (test file was skipped during findGoFiles)
	assert.Equal(t, 2, result.FilesProcessed) // Processed both non-test files
	assert.Equal(t, 0, result.FilesSkipped) // No files skipped during processing (test file was excluded during finding)
	assert.Equal(t, 2, result.FilesOptimized) // Both non-test files had optimizations
	assert.Equal(t, 0, result.FilesWithErrors)
	assert.Equal(t, 2, result.ImportsRemoved) // "os" and "strings"
	assert.Equal(t, 0, result.ImportsAdded)
	assert.GreaterOrEqual(t, result.ProcessingTimeMs, int64(0)) // Allow for 0ms on fast systems

	// Verify individual file results
	mainResult, exists := result.FileResults[filepath.Join(tempDir, "main.go")]
	assert.True(t, exists)
	assert.True(t, mainResult.OptimizationApplied)
	assert.ElementsMatch(t, []string{"os"}, mainResult.ImportsResult.RemovedImports)

	helperResult, exists := result.FileResults[filepath.Join(tempDir, "helper.go")]
	assert.True(t, exists)
	assert.True(t, helperResult.OptimizationApplied)
	assert.ElementsMatch(t, []string{"strings"}, helperResult.ImportsResult.RemovedImports)

	// Test file should have been skipped (no result)
	_, exists = result.FileResults[filepath.Join(tempDir, "utils_test.go")]
	assert.False(t, exists)
}

func TestShouldSkipFile(t *testing.T) {
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	testCases := []struct {
		name     string
		filePath string
		expected bool
	}{
		{"regular go file", "/path/to/main.go", false},
		{"test file with SkipTestFiles=true", "/path/to/main_test.go", true},
		{"vendor directory", "/path/vendor/pkg/file.go", true},
		{"hidden directory", "/path/.git/file.go", false}, // .git is excluded by findGoFiles, not shouldSkipFile
		{"node_modules", "/path/node_modules/file.go", false}, // Same as above
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := pipeline.shouldSkipFile(tc.filePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestShouldSkipFile_WithTestFilesEnabled(t *testing.T) {
	options := DefaultPipelineOptions()
	options.SkipTestFiles = false // Enable processing test files
	pipeline := NewOptimizationPipeline(options)

	assert.False(t, pipeline.shouldSkipFile("/path/to/main_test.go"))
}

func TestPipelineOptions_CustomPatterns(t *testing.T) {
	options := DefaultPipelineOptions()
	options.ExcludePatterns = []string{"generated/", "_gen.go"}
	pipeline := NewOptimizationPipeline(options)

	assert.True(t, pipeline.shouldSkipFile("/path/generated/file.go"))
	assert.True(t, pipeline.shouldSkipFile("/path/to/model_gen.go"))
	assert.False(t, pipeline.shouldSkipFile("/path/to/regular.go"))
}

func TestOptimizeFile_ErrorHandling(t *testing.T) {
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	// Test with non-existent file
	result, err := pipeline.OptimizeFile("/non/existent/file.go")
	assert.Error(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Errors, 1)
}

func TestOptimizeFile_LargeFile(t *testing.T) {
	// Create a temporary file that's too large
	tempDir, err := os.MkdirTemp("", "pipeline-large-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create content that exceeds the size limit
	var largeContent strings.Builder
	largeContent.WriteString("package main\n")
	for i := 0; i < 1000; i++ {
		largeContent.WriteString("// This is a very long comment to make the file large\n")
	}

	testFile := filepath.Join(tempDir, "large.go")
	err = os.WriteFile(testFile, []byte(largeContent.String()), 0644)
	require.NoError(t, err)

	// Create pipeline with small size limit
	options := DefaultPipelineOptions()
	options.MaxFileSize = 100 // Very small limit
	pipeline := NewOptimizationPipeline(options)

	// Optimize the file - should return an error about file being too large
	result, err := pipeline.OptimizeFile(testFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file too large")
	assert.NotNil(t, result)
	assert.Len(t, result.Errors, 1) // Error should also be in the result
}

func TestOptimizeFile_InvalidGoSyntax(t *testing.T) {
	// Create a temporary file with invalid Go syntax
	tempDir, err := os.MkdirTemp("", "pipeline-invalid-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	invalidCode := `package main

import "fmt"

func main() {
	fmt.Println("Hello World"
	// Missing closing parenthesis - invalid syntax
}`

	testFile := filepath.Join(tempDir, "invalid.go")
	err = os.WriteFile(testFile, []byte(invalidCode), 0644)
	require.NoError(t, err)

	// Create pipeline
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	// Optimize the file - should return a parsing error
	result, err := pipeline.OptimizeFile(testFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse file")
	assert.NotNil(t, result)
	assert.Len(t, result.Errors, 1) // Error should also be in the result
}

func TestFindGoFiles(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "pipeline-find-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create directory structure with various files
	files := []string{
		"main.go",
		"helper.go",
		"main_test.go",
		"subdir/utils.go",
		"subdir/utils_test.go",
		"vendor/pkg/vendored.go", // Should be skipped
		".git/hooks/pre-commit",  // Should be skipped
		"README.md",              // Not a Go file
	}

	for _, file := range files {
		fullPath := filepath.Join(tempDir, file)
		
		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err = os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		// Create file
		err = os.WriteFile(fullPath, []byte("package main"), 0644)
		require.NoError(t, err)
	}

	// Test finding Go files
	options := DefaultPipelineOptions()
	pipeline := NewOptimizationPipeline(options)

	goFiles, err := pipeline.findGoFiles(tempDir)
	require.NoError(t, err)

	// Should find non-test Go files (test files are skipped by shouldSkipFile during findGoFiles)  
	expected := []string{
		filepath.Join(tempDir, "main.go"),
		filepath.Join(tempDir, "helper.go"),
		filepath.Join(tempDir, "subdir/utils.go"),
	}

	assert.ElementsMatch(t, expected, goFiles)
}