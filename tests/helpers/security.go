package helpers

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"testing"
)

// SafePath validates that a path is safe and doesn't contain path traversal attempts
// Improved for cross-platform compatibility
func SafePath(path string) bool {
	// Clean the path to resolve any .. or . components
	cleanedPath := filepath.Clean(path)
	
	// Check for path traversal attempts - after cleaning, no .. should remain
	if strings.Contains(cleanedPath, "..") {
		return false
	}
	
	// Check for null bytes (security vulnerability)
	if strings.Contains(path, "\x00") {
		return false
	}
	
	// For absolute paths, ensure they're within safe boundaries
	if filepath.IsAbs(cleanedPath) {
		// Allow common test directories across platforms
		safePrefixes := []string{
			"/tmp/",           // Unix/Linux temporary
			"/var/folders/",   // macOS temporary  
			"/private/tmp/",   // macOS alternative temporary
		}
		
		// On Windows, check for safe temporary paths
		if strings.Contains(cleanedPath, "\\") {
			windowsSafePrefixes := []string{
				"C:\\Temp\\",
				"C:\\tmp\\",
				"C:\\Users\\",
			}
			safePrefixes = append(safePrefixes, windowsSafePrefixes...)
		}
		
		// Allow if path starts with any safe prefix
		for _, prefix := range safePrefixes {
			if strings.HasPrefix(cleanedPath, prefix) {
				return true
			}
		}
		
		// Reject absolute paths that don't match safe prefixes
		return false
	}
	
	// Relative paths are generally safe after cleaning
	return true
}

// SafeProjectPath validates that a project path is safe for testing
func SafeProjectPath(t TB, path string) string {
	t.Helper()
	
	if !SafePath(path) {
		t.Fatalf("Unsafe path detected: %s", path)
	}
	
	return path
}

// AssertDirectoryExistsWithMessage asserts directory exists with detailed error message
func AssertDirectoryExistsWithMessage(t *testing.T, path, message string) {
	t.Helper()
	
	// Validate path safety first
	safePath := SafeProjectPath(t, path)
	
	if !DirExists(safePath) {
		t.Errorf("Directory does not exist at %s: %s", safePath, message)
	}
}

// AssertFileExistsWithMessage asserts file exists with detailed error message
func AssertFileExistsWithMessage(t *testing.T, path, message string) {
	t.Helper()
	
	// Validate path safety first
	safePath := SafeProjectPath(t, path)
	
	if !FileExists(safePath) {
		t.Errorf("File does not exist at %s: %s", safePath, message)
	}
}

// AssertFileContainsImport validates that a file contains a specific import statement using AST parsing
func AssertFileContainsImport(t *testing.T, filePath, importPath, message string) {
	t.Helper()
	
	// Validate path safety first
	safePath := SafeProjectPath(t, filePath)
	
	if !FileExists(safePath) {
		t.Errorf("File does not exist for import validation: %s", safePath)
		return
	}
	
	// Parse the Go file using AST
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, safePath, nil, parser.ParseComments)
	if err != nil {
		t.Errorf("Failed to parse Go file %s: %v", safePath, err)
		return
	}
	
	// Check imports in the AST
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if x.Path != nil {
				// Remove quotes from import path
				impPath := strings.Trim(x.Path.Value, `"`)
				if impPath == importPath || strings.Contains(impPath, importPath) {
					found = true
					return false // Stop inspection
				}
			}
		}
		return true
	})
	
	if !found {
		t.Errorf("Import %s not found in %s: %s", importPath, safePath, message)
	} else {
		t.Logf("✓ Import validation for %s in %s: %s", importPath, safePath, message)
	}
}

// AssertFileContainsFrameworkImport validates framework-specific imports
func AssertFileContainsFrameworkImport(t *testing.T, filePath, framework, message string) {
	t.Helper()
	
	frameworkImports := map[string][]string{
		"gin":   {"github.com/gin-gonic/gin", "gin"},
		"echo":  {"github.com/labstack/echo", "echo"},
		"fiber": {"github.com/gofiber/fiber", "fiber"},
		"chi":   {"github.com/go-chi/chi", "chi"},
	}
	
	imports, ok := frameworkImports[framework]
	if !ok {
		t.Errorf("Unknown framework %s for import validation", framework)
		return
	}
	
	// Try to find any of the framework's imports
	for _, importPath := range imports {
		// Use a non-failing version to check each import
		if fileContainsImport(t, filePath, importPath) {
			t.Logf("✓ Framework %s import found in %s: %s", framework, filePath, message)
			return
		}
	}
	
	t.Errorf("No %s framework imports found in %s: %s", framework, filePath, message)
}

// fileContainsImport is a helper that checks for import without failing the test
func fileContainsImport(t *testing.T, filePath, importPath string) bool {
	t.Helper()
	
	safePath := SafeProjectPath(t, filePath)
	if !FileExists(safePath) {
		return false
	}
	
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, safePath, nil, parser.ParseComments)
	if err != nil {
		return false
	}
	
	found := false
	ast.Inspect(node, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.ImportSpec:
			if x.Path != nil {
				impPath := strings.Trim(x.Path.Value, `"`)
				if impPath == importPath || strings.Contains(impPath, importPath) {
					found = true
					return false // Stop inspection
				}
			}
		}
		return true
	})
	
	return found
}

// AssertArchitectureStructure validates architecture-specific directory structure with security
func AssertArchitectureStructure(t *testing.T, projectPath string, expectedDirs map[string]string, architecture string) {
	t.Helper()
	
	// Validate project path safety
	safeProjectPath := SafeProjectPath(t, projectPath)
	
	for dir, purpose := range expectedDirs {
		fullPath := filepath.Join(safeProjectPath, dir)
		
		// Additional safety check for the full path
		if !SafePath(fullPath) {
			t.Errorf("Unsafe path in architecture structure: %s", fullPath)
			continue
		}
		
		if DirExists(fullPath) {
			t.Logf("✓ %s architecture directory %s exists (Purpose: %s)", architecture, dir, purpose)
		} else {
			t.Errorf("Expected directory %s not found for %s architecture (Purpose: %s)", dir, architecture, purpose)
		}
	}
}

// AssertDirectoryStructure validates expected directory structure with detailed reporting
func AssertDirectoryStructure(t *testing.T, projectPath string, expectedDirs []string, context string) {
	t.Helper()
	
	safeProjectPath := SafeProjectPath(t, projectPath)
	
	for _, dir := range expectedDirs {
		fullPath := filepath.Join(safeProjectPath, dir)
		
		if !SafePath(fullPath) {
			t.Errorf("Unsafe path in directory structure: %s", fullPath)
			continue
		}
		
		if DirExists(fullPath) {
			t.Logf("✓ %s directory %s exists", context, dir)
		} else {
			t.Errorf("Expected directory %s not found for %s", dir, context)
		}
	}
}

// AssertFileStructure validates expected file structure with detailed reporting
func AssertFileStructure(t *testing.T, projectPath string, expectedFiles []string, context string) {
	t.Helper()
	
	safeProjectPath := SafeProjectPath(t, projectPath)
	
	for _, file := range expectedFiles {
		fullPath := filepath.Join(safeProjectPath, file)
		
		if !SafePath(fullPath) {
			t.Errorf("Unsafe path in file structure: %s", fullPath)
			continue
		}
		
		if FileExists(fullPath) {
			t.Logf("✓ %s file %s exists", context, file)
		} else {
			t.Errorf("Expected file %s not found for %s", file, context)
		}
	}
}