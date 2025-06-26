package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
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