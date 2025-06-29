package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// AssertProjectGenerated validates complete project generation
func AssertProjectGenerated(t *testing.T, outputDir string, expectedFiles []string) {
	t.Helper()
	for _, file := range expectedFiles {
		assert.FileExists(t, filepath.Join(outputDir, file),
			"Expected file %s should exist", file)
	}
}

// AssertGoModValid validates go.mod file structure
func AssertGoModValid(t *testing.T, goModPath string, expectedModule string) {
	t.Helper()
	content, err := os.ReadFile(goModPath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), expectedModule)
}

// AssertGoModContainsVersion validates go.mod contains specific Go version
func AssertGoModContainsVersion(t *testing.T, goModPath string, expectedVersion string) {
	t.Helper()
	content, err := os.ReadFile(goModPath)
	assert.NoError(t, err)
	if expectedVersion == "auto" {
		// For auto, just verify go directive exists
		assert.Contains(t, string(content), "go ")
	} else {
		assert.Contains(t, string(content), "go "+expectedVersion)
	}
}

// AssertCompilable validates generated project compiles
func AssertCompilable(t *testing.T, projectDir string) {
	t.Helper()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = projectDir
	err := cmd.Run()
	assert.NoError(t, err, "Generated project should compile successfully")
}

// AssertFileContains validates file contains expected content
func AssertFileContains(t *testing.T, filePath string, expectedContent string) {
	t.Helper()
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), expectedContent)
}

// AssertDirExists validates directory exists
func AssertDirExists(t *testing.T, dirPath string) {
	t.Helper()
	info, err := os.Stat(dirPath)
	assert.NoError(t, err)
	assert.True(t, info.IsDir(), "Path %s should be a directory", dirPath)
}

// AssertFileExists validates file exists
func AssertFileExists(t *testing.T, filePath string) {
	t.Helper()
	_, err := os.Stat(filePath)
	assert.NoError(t, err, "File %s should exist", filePath)
}

// GetFileInfo returns file info or error
func GetFileInfo(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// FileExists checks if file exists
func FileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// DirExists checks if directory exists
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// ReadFileContent reads file content as string
func ReadFileContent(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read file %s", filePath)
	return string(content)
}

// FindFiles finds files matching pattern in directory
func FindFiles(t *testing.T, dir string, pattern string) []string {
	t.Helper()

	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			matched, err := filepath.Match(pattern, info.Name())
			if err != nil {
				return err
			}
			if matched {
				files = append(files, path)
			}
		}
		return nil
	})

	if err != nil {
		t.Logf("Warning: Could not walk directory %s: %v", dir, err)
		return []string{}
	}

	return files
}

// StringContains checks if string contains substring
func StringContains(s, substr string) bool {
	return strings.Contains(s, substr)
}
