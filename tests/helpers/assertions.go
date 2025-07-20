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
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Generated project should compile successfully.\nBuild output:\n%s\nError: %v", string(output), err)
	}
}

// AssertFileContains validates file contains expected content
func AssertFileContains(t *testing.T, filePath string, expectedContent string) {
	t.Helper()
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Contains(t, string(content), expectedContent)
}

// AssertNotContains validates that content does not contain specific text
func AssertNotContains(t *testing.T, content string, unwantedText string) {
	t.Helper()
	assert.NotContains(t, content, unwantedText,
		"Content should not contain '%s'", unwantedText)
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

// AssertProjectCompiles validates that the project compiles successfully
func AssertProjectCompiles(t *testing.T, projectPath string) {
	t.Helper()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Generated project should compile successfully.\nBuild output:\n%s\nError: %v", string(output), err)
	}
}

// AssertDirectoryExists validates directory exists
func AssertDirectoryExists(t *testing.T, dirPath string) {
	t.Helper()
	info, err := os.Stat(dirPath)
	assert.NoError(t, err, "Directory %s should exist", dirPath)
	assert.True(t, info.IsDir(), "Path %s should be a directory", dirPath)
}

// AssertCLIHelpOutput validates CLI help output
func AssertCLIHelpOutput(t *testing.T, projectPath string) {
	t.Helper()
	
	// Build the project first
	buildCmd := exec.Command("go", "build", "-o", "test-cli", ".")
	buildCmd.Dir = projectPath
	buildOutput, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		t.Errorf("Failed to build CLI: %v\nOutput: %s", buildErr, string(buildOutput))
		return
	}
	
	// Run help command
	helpCmd := exec.Command("./test-cli", "--help")
	helpCmd.Dir = projectPath
	helpOutput, helpErr := helpCmd.CombinedOutput()
	if helpErr != nil {
		t.Errorf("Failed to run help command: %v\nOutput: %s", helpErr, string(helpOutput))
		return
	}
	
	// Verify help output contains expected content
	helpText := string(helpOutput)
	assert.Contains(t, helpText, "Usage:")
	assert.Contains(t, helpText, "Available Commands:")
}

// AssertCLIVersionOutput validates CLI version output
func AssertCLIVersionOutput(t *testing.T, projectPath string) {
	t.Helper()
	
	// Build the project first
	buildCmd := exec.Command("go", "build", "-o", "test-cli", ".")
	buildCmd.Dir = projectPath
	buildOutput, buildErr := buildCmd.CombinedOutput()
	if buildErr != nil {
		t.Errorf("Failed to build CLI: %v\nOutput: %s", buildErr, string(buildOutput))
		return
	}
	
	// Run version command
	versionCmd := exec.Command("./test-cli", "version")
	versionCmd.Dir = projectPath
	versionOutput, versionErr := versionCmd.CombinedOutput()
	if versionErr != nil {
		t.Errorf("Failed to run version command: %v\nOutput: %s", versionErr, string(versionOutput))
		return
	}
	
	// Verify version output contains expected content
	versionText := string(versionOutput)
	assert.Contains(t, versionText, "version")
}

// AssertLoggerFunctionality validates logger functionality
func AssertLoggerFunctionality(t *testing.T, projectPath string, logger string) {
	t.Helper()
	
	// Check if logger specific files exist and contain expected content
	loggerFile := filepath.Join(projectPath, "internal", "logger", logger+".go")
	if FileExists(loggerFile) {
		content := ReadFileContent(t, loggerFile)
		
		// Check for common logger patterns
		switch logger {
		case "slog":
			assert.Contains(t, content, "log/slog")
		case "zap":
			assert.Contains(t, content, "go.uber.org/zap")
		case "logrus":
			assert.Contains(t, content, "github.com/sirupsen/logrus")
		case "zerolog":
			assert.Contains(t, content, "github.com/rs/zerolog")
		}
	}
}

// AssertTestsRun validates that tests run successfully
func AssertTestsRun(t *testing.T, projectPath string) {
	t.Helper()
	
	cmd := exec.Command("go", "test", "./...")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Tests should run successfully.\nTest output:\n%s\nError: %v", string(output), err)
	}
}

// Additional helper functions for CLI-simple testing

// AssertFileNotExists validates file does not exist
func AssertFileNotExists(t *testing.T, filePath string) {
	t.Helper()
	_, err := os.Stat(filePath)
	assert.True(t, os.IsNotExist(err), "File %s should not exist", filePath)
}

// AssertDirectoryNotExists validates directory does not exist
func AssertDirectoryNotExists(t *testing.T, dirPath string) {
	t.Helper()
	_, err := os.Stat(dirPath)
	assert.True(t, os.IsNotExist(err), "Directory %s should not exist", dirPath)
}

// AssertFileNotContains validates file does not contain specific content
func AssertFileNotContains(t *testing.T, filePath string, unwantedContent string) {
	t.Helper()
	if !FileExists(filePath) {
		return // If file doesn't exist, it doesn't contain the content
	}
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.NotContains(t, string(content), unwantedContent, 
		"File %s should not contain '%s'", filePath, unwantedContent)
}

// CountFiles counts files with specific extension in directory
func CountFiles(t *testing.T, dirPath string, extension string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == extension {
			count++
		}
		return nil
	})
	assert.NoError(t, err, "Failed to count files in %s", dirPath)
	return count
}

// CountAllFiles counts all files in directory (excluding directories)
func CountAllFiles(t *testing.T, dirPath string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	assert.NoError(t, err, "Failed to count files in %s", dirPath)
	return count
}

// ReadFile reads file content as string
func ReadFile(t *testing.T, filePath string) string {
	t.Helper()
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err, "Failed to read file %s", filePath)
	return string(content)
}

// CountLines counts lines in text content
func CountLines(content string) int {
	if content == "" {
		return 0
	}
	return strings.Count(content, "\n") + 1
}

// CountDependencies counts dependencies in go.mod content
func CountDependencies(goModContent string) int {
	lines := strings.Split(goModContent, "\n")
	count := 0
	inRequireBlock := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "require (" {
			inRequireBlock = true
			continue
		}
		if line == ")" && inRequireBlock {
			inRequireBlock = false
			continue
		}
		if inRequireBlock && line != "" && !strings.Contains(line, "// indirect") {
			count++
		} else if strings.HasPrefix(line, "require ") && !strings.Contains(line, "// indirect") {
			count++
		}
	}
	return count
}

// FindGoFiles finds all .go files in directory
func FindGoFiles(t *testing.T, dirPath string) []string {
	t.Helper()
	var files []string
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			files = append(files, path)
		}
		return nil
	})
	assert.NoError(t, err, "Failed to find Go files in %s", dirPath)
	return files
}

// CountInterfaces counts interface definitions in Go code
func CountInterfaces(content string) int {
	// Simple count of "type ... interface" patterns
	return strings.Count(content, "interface{")
}

// GetMaxDirectoryDepth gets maximum directory depth in project
func GetMaxDirectoryDepth(t *testing.T, dirPath string) int {
	t.Helper()
	maxDepth := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			relPath, err := filepath.Rel(dirPath, path)
			if err != nil {
				return err
			}
			depth := strings.Count(relPath, string(filepath.Separator))
			if depth > maxDepth {
				maxDepth = depth
			}
		}
		return nil
	})
	assert.NoError(t, err, "Failed to calculate directory depth in %s", dirPath)
	return maxDepth
}
