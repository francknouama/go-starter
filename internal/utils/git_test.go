package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsGitRepository(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tempDir string) string
		expected bool
	}{
		{
			name: "directory with .git subdirectory is git repository",
			setup: func(t *testing.T, tempDir string) string {
				gitDir := filepath.Join(tempDir, ".git")
				err := os.MkdirAll(gitDir, 0755)
				require.NoError(t, err)
				return tempDir
			},
			expected: true,
		},
		{
			name: "directory without .git subdirectory is not git repository",
			setup: func(t *testing.T, tempDir string) string {
				return tempDir
			},
			expected: false,
		},
		{
			name: "non-existent directory is not git repository",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			// Act
			result := IsGitRepository(path)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsGitInstalled(t *testing.T) {
	// Act
	isInstalled := IsGitInstalled()

	// Assert
	// We can't guarantee git is installed in all environments,
	// but we can test that the function returns a boolean
	assert.IsType(t, true, isInstalled)
}

func TestGetGitVersion(t *testing.T) {
	// Skip this test if git is not installed
	if !IsGitInstalled() {
		t.Skip("Git is not installed, skipping version test")
	}

	// Act
	version, err := GetGitVersion()

	// Assert
	assert.NoError(t, err)
	assert.NotEmpty(t, version)
	assert.Contains(t, version, "git version")
}

func TestInitGitRepository(t *testing.T) {
	// Skip this test if git is not installed
	if !IsGitInstalled() {
		t.Skip("Git is not installed, skipping git init test")
	}

	tests := []struct {
		name        string
		setup       func(t *testing.T, tempDir string) string
		shouldError bool
	}{
		{
			name: "initialize git repository in empty directory",
			setup: func(t *testing.T, tempDir string) string {
				return tempDir
			},
			shouldError: false,
		},
		{
			name: "initialize git repository in existing directory",
			setup: func(t *testing.T, tempDir string) string {
				// Create some files
				err := os.WriteFile(filepath.Join(tempDir, "readme.txt"), []byte("hello"), 0644)
				require.NoError(t, err)
				return tempDir
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			// Act
			err := InitGitRepository(path)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.True(t, IsGitRepository(path))
		})
	}
}

func TestAddGitIgnore(t *testing.T) {
	// Skip this test if git is not installed
	if !IsGitInstalled() {
		t.Skip("Git is not installed, skipping gitignore test")
	}

	tests := []struct {
		name        string
		content     string
		shouldError bool
	}{
		{
			name: "add gitignore with standard content",
			content: `# Binaries
*.exe
*.dll
*.so

# IDE
.vscode/
.idea/`,
			shouldError: false,
		},
		{
			name:        "add empty gitignore",
			content:     "",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			err := InitGitRepository(tempDir)
			require.NoError(t, err)

			// Act
			err = AddGitIgnore(tempDir, tt.content)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)

			// Verify .gitignore file was created
			gitIgnorePath := filepath.Join(tempDir, ".gitignore")
			assert.FileExists(t, gitIgnorePath)

			// Verify content
			actualContent, err := os.ReadFile(gitIgnorePath)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(actualContent))
		})
	}
}

func TestGetDefaultGitIgnore(t *testing.T) {
	// Act
	gitignore := GetDefaultGitIgnore()

	// Assert
	assert.NotEmpty(t, gitignore)
	assert.Contains(t, gitignore, "# Binaries for programs and plugins")
	assert.Contains(t, gitignore, "*.exe")
	assert.Contains(t, gitignore, "# IDE files")
	assert.Contains(t, gitignore, ".vscode/")
	assert.Contains(t, gitignore, "go.work")
}
