package optimization

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestImportOptimizationIntegration(t *testing.T) {
	// Test code with unused imports and disorganized order
	testCode := `package main

import (
	"github.com/gin-gonic/gin"
	"fmt" 
	"os"
	"strings"
	"net/http"
)

func main() {
	fmt.Println("Hello World")
	gin.New()
	// os, strings, and net/http are unused
}`

	// Create analyzer with all import optimizations enabled
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
		OrganizeImports:     true,
	})
	manager := NewImportManager(analyzer)

	// Parse the test code
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)

	analyzer.fileSet = fset

	// Optimize imports
	result, err := manager.OptimizeImports(file)
	require.NoError(t, err)

	// Verify optimization was applied
	assert.True(t, result.OptimizationApplied)
	
	// Should have removed unused imports
	assert.ElementsMatch(t, []string{"os", "strings", "net/http"}, result.RemovedImports)
	assert.Empty(t, result.AddedImports)

	// Get optimized content
	optimizedContent, err := manager.GetOptimizedFileContent(file)
	require.NoError(t, err)

	// Should contain only used imports
	assert.Contains(t, optimizedContent, `"fmt"`)
	assert.Contains(t, optimizedContent, `"github.com/gin-gonic/gin"`)
	
	// Should not contain unused imports
	assert.NotContains(t, optimizedContent, `"os"`)
	assert.NotContains(t, optimizedContent, `"strings"`)
	assert.NotContains(t, optimizedContent, `"net/http"`)

	// Should be organized alphabetically
	fmtPos := strings.Index(optimizedContent, `"fmt"`)
	ginPos := strings.Index(optimizedContent, `"github.com/gin-gonic/gin"`)
	assert.True(t, fmtPos < ginPos, "Imports should be organized alphabetically")

	t.Logf("Optimized code:\n%s", optimizedContent)
}

func TestNoOptimizationNeeded(t *testing.T) {
	// Test code that doesn't need optimization
	testCode := `package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello World")
}`

	// Create analyzer with all optimizations enabled
	analyzer := NewASTAnalyzer(AnalysisOptions{
		RemoveUnusedImports: true,
		OrganizeImports:     true,
	})
	manager := NewImportManager(analyzer)

	// Parse and optimize
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
	require.NoError(t, err)

	analyzer.fileSet = fset

	result, err := manager.OptimizeImports(file)
	require.NoError(t, err)

	// No optimization should be needed
	assert.False(t, result.OptimizationApplied)
	assert.Empty(t, result.RemovedImports)
	assert.Empty(t, result.AddedImports)
	assert.Equal(t, []string{"fmt"}, result.OriginalImports)
}