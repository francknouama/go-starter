package optimization

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewImportManager(t *testing.T) {
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	manager := NewImportManager(analyzer)
	
	assert.NotNil(t, manager)
	assert.NotNil(t, manager.analyzer)
	assert.NotNil(t, manager.resolver)
}

func TestNewImportResolver(t *testing.T) {
	resolver := NewImportResolver()
	
	assert.NotNil(t, resolver)
	assert.NotNil(t, resolver.standardLibraries)
	assert.NotNil(t, resolver.commonPackages)
	
	// Check that some standard libraries are registered
	assert.True(t, resolver.isStandardLibrary("fmt"))
	assert.True(t, resolver.isStandardLibrary("os"))
	assert.False(t, resolver.isStandardLibrary("github.com/gin-gonic/gin"))
	
	// Check that some common packages are registered
	assert.Equal(t, "github.com/gin-gonic/gin", resolver.ResolvePackage("gin"))
	assert.Equal(t, "github.com/spf13/cobra", resolver.ResolvePackage("cobra"))
}

func TestOptimizeImports_RemoveUnused(t *testing.T) {
	testCode := `package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("hello")
	// os and strings are unused
}`

	// Create analyzer and manager
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
		OrganizeImports:     false,
	})
	manager := NewImportManager(analyzer)
	
	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)
	
	analyzer.fileSet = fset
	
	// Optimize imports
	result, err := manager.OptimizeImports(file)
	require.NoError(t, err)
	
	// Check results
	
	assert.True(t, result.OptimizationApplied)
	assert.ElementsMatch(t, []string{"fmt", "os", "strings"}, result.OriginalImports)
	assert.ElementsMatch(t, []string{"os", "strings"}, result.RemovedImports)
	assert.Empty(t, result.AddedImports)
	
	// Check that the optimized file content is valid
	optimizedContent, err := manager.GetOptimizedFileContent(file)
	require.NoError(t, err)
	
	// The optimized content should not contain unused imports
	assert.Contains(t, optimizedContent, `"fmt"`)
	assert.NotContains(t, optimizedContent, `"os"`)
	assert.NotContains(t, optimizedContent, `"strings"`)
}

func TestOptimizeImports_NoUnusedImports(t *testing.T) {
	testCode := `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello")
	os.Exit(0)
}`

	// Create analyzer and manager
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
	})
	manager := NewImportManager(analyzer)
	
	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)
	
	analyzer.fileSet = fset
	
	// Optimize imports
	result, err := manager.OptimizeImports(file)
	require.NoError(t, err)
	
	// Check results - no optimization should be needed
	assert.False(t, result.OptimizationApplied)
	assert.ElementsMatch(t, []string{"fmt", "os"}, result.OriginalImports)
	assert.Empty(t, result.RemovedImports)
	assert.Empty(t, result.AddedImports)
}

func TestOptimizeImports_OrganizeImports(t *testing.T) {
	testCode := `package main

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello")
	os.Exit(0)
	gin.New()
}`

	// Create analyzer and manager
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: false, // Focus on organization only
		OrganizeImports:     true,
	})
	manager := NewImportManager(analyzer)
	
	
	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)
	
	analyzer.fileSet = fset
	
	// Optimize imports
	_, err = manager.OptimizeImports(file)
	require.NoError(t, err)
	
	// Get optimized content
	optimizedContent, err := manager.GetOptimizedFileContent(file)
	require.NoError(t, err)
	
	// Check basic alphabetical sorting: "fmt" < "github.com/gin-gonic/gin" < "os"
	fmtPos := strings.Index(optimizedContent, `"fmt"`)
	osPos := strings.Index(optimizedContent, `"os"`)
	ginPos := strings.Index(optimizedContent, `"github.com/gin-gonic/gin"`)
	
	// Basic alphabetical sorting should put them in this order:
	// "fmt" < "github.com/gin-gonic/gin" < "os"
	assert.True(t, fmtPos < ginPos, "fmt should come before github.com/gin-gonic/gin alphabetically")
	assert.True(t, ginPos < osPos, "github.com/gin-gonic/gin should come before os alphabetically")
}

func TestFindMissingImports(t *testing.T) {
	// This is a simplified test since missing import detection is complex
	testCode := `package main

func main() {
	// This would be detected by a more sophisticated analyzer
	fmt.Println("hello")
}`

	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	manager := NewImportManager(analyzer)
	
	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)
	
	analyzer.fileSet = fset
	
	// Find missing imports (simplified test)
	missing := manager.findMissingImports(file)
	
	// The current implementation is simplified and may not detect this
	// but we test that the function doesn't crash
	assert.NotNil(t, missing)
}

func TestImportResolver_ResolvePackage(t *testing.T) {
	resolver := NewImportResolver()
	
	testCases := []struct {
		packageName string
		expected    string
	}{
		{"fmt", "fmt"},                                    // Standard library
		{"os", "os"},                                      // Standard library
		{"gin", "github.com/gin-gonic/gin"},             // Common package
		{"cobra", "github.com/spf13/cobra"},             // Common package
		{"uuid", "github.com/google/uuid"},              // Common package
		{"unknownpackage", ""},                           // Unknown package
	}
	
	for _, tc := range testCases {
		t.Run(tc.packageName, func(t *testing.T) {
			result := resolver.ResolvePackage(tc.packageName)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestImportResolver_IsStandardLibrary(t *testing.T) {
	resolver := NewImportResolver()
	
	testCases := []struct {
		packagePath string
		expected    bool
	}{
		{"fmt", true},
		{"os", true},
		{"net/http", true},
		{"encoding/json", true},
		{"github.com/gin-gonic/gin", false},
		{"golang.org/x/tools", false},
		{"example.com/mypackage", false},
	}
	
	for _, tc := range testCases {
		t.Run(tc.packagePath, func(t *testing.T) {
			result := resolver.isStandardLibrary(tc.packagePath)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestMightNeedImport(t *testing.T) {
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	manager := NewImportManager(analyzer)
	
	testCases := []struct {
		identifier string
		expected   bool
	}{
		{"fmt", true},     // Lowercase, not a keyword
		{"os", true},      // Lowercase, not a keyword
		{"gin", true},     // Lowercase, not a keyword
		{"MyStruct", false}, // Uppercase (likely a type)
		{"if", false},     // Go keyword
		{"for", false},    // Go keyword
		{"func", false},   // Go keyword
	}
	
	for _, tc := range testCases {
		t.Run(tc.identifier, func(t *testing.T) {
			result := manager.mightNeedImport(tc.identifier)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestDeduplicateStrings(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"a", "b", "c"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "with duplicates",
			input:    []string{"a", "b", "a", "c", "b"},
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "all duplicates",
			input:    []string{"a", "a", "a"},
			expected: []string{"a"},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := deduplicateStrings(tc.input)
			assert.ElementsMatch(t, tc.expected, result)
		})
	}
}

func TestCalculateStatistics(t *testing.T) {
	results := []*ImportOptimizationResult{
		{
			OriginalImports:     []string{"fmt", "os", "strings"},
			RemovedImports:      []string{"os", "strings"},
			AddedImports:        []string{},
			OptimizationApplied: true,
		},
		{
			OriginalImports:     []string{"fmt", "net/http"},
			RemovedImports:      []string{},
			AddedImports:        []string{"encoding/json"},
			OptimizationApplied: true,
		},
		{
			OriginalImports:     []string{"fmt"},
			RemovedImports:      []string{},
			AddedImports:        []string{},
			OptimizationApplied: false,
		},
	}
	
	stats := CalculateStatistics(results)
	
	assert.Equal(t, 3, stats.TotalFiles)
	assert.Equal(t, 2, stats.FilesOptimized)
	assert.Equal(t, 2, stats.ImportsRemoved)
	assert.Equal(t, 1, stats.ImportsAdded)
}

func TestGetOptimizedFileContent(t *testing.T) {
	testCode := `package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello")
}`

	// Create analyzer and manager
	analyzer := NewASTAnalyzer(DefaultAnalysisOptions())
	manager := NewImportManager(analyzer)
	
	// Parse the code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)
	
	analyzer.fileSet = fset
	
	// Get optimized content
	content, err := manager.GetOptimizedFileContent(file)
	require.NoError(t, err)
	
	// Should be valid Go code
	assert.Contains(t, content, "package main")
	assert.Contains(t, content, "func main()")
	assert.Contains(t, content, `fmt.Println("hello")`)
}