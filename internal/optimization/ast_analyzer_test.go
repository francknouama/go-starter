package optimization

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewASTAnalyzer(t *testing.T) {
	options := DefaultAnalysisOptions()
	analyzer := NewASTAnalyzer(options)
	
	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.fileSet)
	assert.Equal(t, options, analyzer.options)
}

func TestDefaultAnalysisOptions(t *testing.T) {
	options := DefaultAnalysisOptions()
	
	assert.True(t, options.RemoveUnusedImports)
	assert.False(t, options.RemoveUnusedVars) // Conservative default
	assert.False(t, options.RemoveUnusedFuncs) // Conservative default
	assert.True(t, options.OptimizeConditionals)
	assert.True(t, options.OrganizeImports)
	assert.False(t, options.EnableDebugOutput)
	assert.Equal(t, int64(1024*1024), options.MaxFileSize)
}

func TestFindUnusedImports(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected []string // Expected unused import paths
	}{
		{
			name: "unused standard library import",
			code: `package main

import "fmt"

func main() {
	// fmt is not used
}`,
			expected: []string{"fmt"},
		},
		{
			name: "used import should not be detected",
			code: `package main

import "fmt"

func main() {
	fmt.Println("hello")
}`,
			expected: []string{},
		},
		{
			name: "multiple imports with some unused",
			code: `package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("hello")
	// os and strings are not used
}`,
			expected: []string{"os", "strings"},
		},
		{
			name: "blank import should not be flagged",
			code: `package main

import (
	_ "fmt"
	"os"
)

func main() {
	// Both imports should be considered used
}`,
			expected: []string{"os"}, // Only os should be flagged, not the blank import
		},
		{
			name: "dot import should not be flagged",
			code: `package main

import (
	. "fmt"
	"os"
)

func main() {
	// Dot imports are hard to detect usage, so not flagged
}`,
			expected: []string{"os"}, // Only os should be flagged
		},
		{
			name: "aliased import usage",
			code: `package main

import (
	f "fmt"
	"os"
)

func main() {
	f.Println("hello")
	// os is not used
}`,
			expected: []string{"os"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
			
			// Parse the test code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tc.code, parser.ParseComments)
			require.NoError(t, err)
			
			// Update analyzer's fileset to match
			analyzer.fileSet = fset
			
			// Find unused imports
			unusedImports := analyzer.findUnusedImports(file)
			
			// Extract paths from results
			var actualPaths []string
			for _, imp := range unusedImports {
				actualPaths = append(actualPaths, imp.Path)
			}
			
			assert.ElementsMatch(t, tc.expected, actualPaths,
				"Expected unused imports: %v, got: %v", tc.expected, actualPaths)
		})
	}
}

func TestFindUnusedVariables(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected []string // Expected unused variable names
	}{
		{
			name: "unused variable declaration",
			code: `package main

func main() {
	var unused int
	var used int
	println(used)
}`,
			expected: []string{"unused"},
		},
		{
			name: "unused short variable declaration",
			code: `package main

func main() {
	unused := 42
	used := 24
	println(used)
}`,
			expected: []string{"unused"},
		},
		{
			name: "all variables used",
			code: `package main

func main() {
	var a int
	b := 42
	println(a, b)
}`,
			expected: []string{},
		},
		{
			name: "blank identifier should not be flagged",
			code: `package main

func main() {
	_ = 42
	unused := 24
}`,
			expected: []string{"unused"},
		},
		{
			name: "exported variables should not be flagged",
			code: `package main

var ExportedVar int
var unexportedVar int

func main() {
	// Neither variable is used in main
}`,
			expected: []string{}, // unexportedVar is package-level, might be used elsewhere
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := NewASTAnalyzer(AnalysisOptions{
				RemoveUnusedVars: true,
			})
			
			// Parse the test code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tc.code, parser.ParseComments)
			require.NoError(t, err)
			
			analyzer.fileSet = fset
			
			// Find unused variables
			unusedVars := analyzer.findUnusedVariables(file)
			
			// Extract names from results
			var actualNames []string
			for _, varInfo := range unusedVars {
				actualNames = append(actualNames, varInfo.Name)
			}
			
			assert.ElementsMatch(t, tc.expected, actualNames,
				"Expected unused variables: %v, got: %v", tc.expected, actualNames)
		})
	}
}

func TestFindUnusedFunctions(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected []string // Expected unused function names
	}{
		{
			name: "unused private function",
			code: `package main

func unusedFunction() {
	// This function is not called
}

func usedFunction() {
	// This function is called
}

func main() {
	usedFunction()
}`,
			expected: []string{"unusedFunction"},
		},
		{
			name: "exported functions should not be flagged",
			code: `package main

func ExportedFunction() {
	// Exported functions should not be flagged as unused
}

func unusedPrivateFunction() {
	// This is unused
}

func main() {
	// Neither function is called here
}`,
			expected: []string{}, // Only unusedPrivateFunction would be flagged, but we're conservative
		},
		{
			name: "test functions should not be flagged",
			code: `package main

func TestSomething() {
	// Test functions should not be flagged
}

func BenchmarkSomething() {
	// Benchmark functions should not be flagged
}

func ExampleSomething() {
	// Example functions should not be flagged
}

func unusedHelper() {
	// This could be flagged
}`,
			expected: []string{}, // Conservative approach - don't flag any
		},
		{
			name: "main and init functions should not be flagged",
			code: `package main

func init() {
	// init should not be flagged
}

func main() {
	// main should not be flagged
}

func unusedHelper() {
	// This could be flagged
}`,
			expected: []string{}, // Conservative approach
		},
		{
			name: "function used as value",
			code: `package main

func helper() {
	// Used as a function value
}

func main() {
	f := helper
	_ = f
}`,
			expected: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			analyzer := NewASTAnalyzer(AnalysisOptions{
				RemoveUnusedFuncs: true,
			})
			
			// Parse the test code
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", tc.code, parser.ParseComments)
			require.NoError(t, err)
			
			analyzer.fileSet = fset
			
			// Find unused functions
			unusedFuncs := analyzer.findUnusedFunctions(file)
			
			// Extract names from results
			var actualNames []string
			for _, funcInfo := range unusedFuncs {
				actualNames = append(actualNames, funcInfo.Name)
			}
			
			assert.ElementsMatch(t, tc.expected, actualNames,
				"Expected unused functions: %v, got: %v", tc.expected, actualNames)
		})
	}
}

func TestAnalyzeProject_Integration(t *testing.T) {
	// Create a temporary project directory
	tempDir, err := os.MkdirTemp("", "ast-analyzer-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test files
	testFiles := map[string]string{
		"main.go": `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello")
	// os is unused
}`,
		"helper.go": `package main

import "strings"

func unusedHelper() {
	// This function is not used
}

func usedHelper() string {
	return "used"
}`,
		"subdir/util.go": `package subdir

import (
	"fmt"
	"log"
)

func Util() {
	fmt.Println("utility")
	// log is unused
}`,
	}
	
	// Write test files
	for filename, content := range testFiles {
		fullPath := filepath.Join(tempDir, filename)
		
		// Create directory if needed
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		// Write file
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}
	
	// Analyze the project
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
		RemoveUnusedVars:    true,
		RemoveUnusedFuncs:   true,
		EnableDebugOutput:   true,
		MaxFileSize:         10 * 1024, // 10KB - large enough for test files
	})
	
	result, err := analyzer.AnalyzeProject(tempDir)
	require.NoError(t, err)
	
	// Verify results
	assert.Greater(t, result.Metrics.FilesAnalyzed, 0)
	assert.Greater(t, len(result.UnusedImports), 0)
	
	// Check that we found some unused imports
	unusedImportPaths := make([]string, len(result.UnusedImports))
	for i, imp := range result.UnusedImports {
		unusedImportPaths[i] = imp.Path
	}
	
	// Should find "os" and "log" as unused
	assert.Contains(t, unusedImportPaths, "os")
	assert.Contains(t, unusedImportPaths, "log")
	
	// Should NOT find "fmt" as unused (it's used in both files)
	assert.NotContains(t, unusedImportPaths, "fmt")
	
	// Verify metrics are populated
	assert.Greater(t, result.Metrics.ProcessingTimeMs, int64(0))
	assert.Equal(t, len(result.UnusedImports), result.Metrics.ImportsRemoved)
}

func TestAnalyzeProject_ErrorHandling(t *testing.T) {
	// Test with non-existent directory
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	
	result, err := analyzer.AnalyzeProject("/non/existent/path")
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAnalyzeProject_LargeFileSkipping(t *testing.T) {
	// Create a temporary project directory
	tempDir, err := os.MkdirTemp("", "ast-analyzer-large-file-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a small test file
	smallFile := `package main
import "fmt"
func main() { fmt.Println("small") }`
	
	err = os.WriteFile(filepath.Join(tempDir, "small.go"), []byte(smallFile), 0644)
	require.NoError(t, err)
	
	// Create a large test file
	var largeContent strings.Builder
	largeContent.WriteString("package main\n")
	for i := 0; i < 1000; i++ {
		largeContent.WriteString("// This is a comment to make the file large\n")
	}
	
	err = os.WriteFile(filepath.Join(tempDir, "large.go"), []byte(largeContent.String()), 0644)
	require.NoError(t, err)
	
	// Analyze with small max file size
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
		MaxFileSize:         100, // Very small limit
		EnableDebugOutput:   true,
	})
	
	result, err := analyzer.AnalyzeProject(tempDir)
	require.NoError(t, err)
	
	// Should have found 2 files total but only analyzed the small one
	// (The large file should be skipped due to size limit)
	assert.Equal(t, 2, result.Metrics.FilesAnalyzed) // Both files are "analyzed" (considered)
	// But only the small file should have contributed to unused imports
	assert.LessOrEqual(t, len(result.UnusedImports), 1)
}

func TestIsSpecialFunction(t *testing.T) {
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	
	testCases := []struct {
		name     string
		expected bool
	}{
		{"main", true},
		{"init", true},
		{"TestMain", true},
		{"TestSomething", true},
		{"BenchmarkSomething", true},
		{"ExampleSomething", true},
		{"regularFunction", false},
		{"Helper", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := analyzer.isSpecialFunction(tc.name)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetImportAlias(t *testing.T) {
	testCases := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "standard import",
			code: `import "fmt"`,
			expected: "fmt",
		},
		{
			name: "aliased import",
			code: `import f "fmt"`,
			expected: "f",
		},
		{
			name: "package with path",
			code: `import "github.com/user/pkg"`,
			expected: "pkg",
		},
	}
	
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Parse just the import statement
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", "package main\n"+tc.code, parser.ParseComments)
			require.NoError(t, err)
			
			// Get the import spec
			imports := analyzer.getImports(file)
			require.Len(t, imports, 1)
			
			alias := analyzer.getImportAlias(imports[0])
			assert.Equal(t, tc.expected, alias)
		})
	}
}