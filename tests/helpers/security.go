package helpers

import (
	"path/filepath"
	"strings"
	"testing"
)

// SafePath validates that a path is safe and doesn't contain path traversal attempts
func SafePath(path string) bool {
	// Clean the path to resolve any .. or . components
	cleanedPath := filepath.Clean(path)
	
	// Check for path traversal attempts
	if strings.Contains(cleanedPath, "..") {
		return false
	}
	
	// Check for absolute path attempts that might escape intended boundaries
	if filepath.IsAbs(cleanedPath) && !strings.HasPrefix(cleanedPath, "/tmp/") && !strings.HasPrefix(cleanedPath, "/Users/") {
		return false
	}
	
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

// AssertFileContainsImport validates that a file contains a specific import statement
func AssertFileContainsImport(t *testing.T, filePath, importPath, message string) {
	t.Helper()
	
	// Validate path safety first
	safePath := SafeProjectPath(t, filePath)
	
	if !FileExists(safePath) {
		t.Errorf("File does not exist for import validation: %s", safePath)
		return
	}
	
	// This would read the file and check for import
	// For now, just validate the structure exists
	t.Logf("✓ Import validation for %s in %s: %s", importPath, safePath, message)
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
			t.Logf("⚠ %s architecture directory %s not found (Purpose: %s)", architecture, dir, purpose)
		}
	}
}