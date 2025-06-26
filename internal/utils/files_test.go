package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFile(t *testing.T) {
	tests := []struct {
		name        string
		setupSrc    func(t *testing.T, tempDir string) string
		setupDst    func(t *testing.T, tempDir string) string
		content     string
		mode        os.FileMode
		shouldError bool
		errorContains string
	}{
		{
			name: "copy file successfully",
			setupSrc: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "source.txt")
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "dest.txt")
			},
			content:     "hello world",
			mode:        0644,
			shouldError: false,
		},
		{
			name: "copy file with custom permissions",
			setupSrc: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "executable.sh")
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "dest_executable.sh")
			},
			content:     "#!/bin/bash\necho hello",
			mode:        0755,
			shouldError: false,
		},
		{
			name: "copy to nested destination directory",
			setupSrc: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "source.txt")
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nested", "subdir", "dest.txt")
			},
			content:     "nested content",
			mode:        0644,
			shouldError: false,
		},
		{
			name: "copy non-existent file should error",
			setupSrc: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "dest.txt")
			},
			content:       "",
			mode:          0644,
			shouldError:   true,
			errorContains: "failed to open source file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			srcPath := tt.setupSrc(t, tempDir)
			dstPath := tt.setupDst(t, tempDir)

			// Create source file if content is provided
			if tt.content != "" && !tt.shouldError {
				err := os.WriteFile(srcPath, []byte(tt.content), tt.mode)
				require.NoError(t, err)
			}

			// Act
			err := CopyFile(srcPath, dstPath)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.FileExists(t, dstPath)

			// Verify content
			actualContent, err := os.ReadFile(dstPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(actualContent))

			// Verify permissions
			info, err := os.Stat(dstPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.mode, info.Mode().Perm())
		})
	}
}

func TestWriteFile(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		content       string
		shouldError   bool
		errorContains string
	}{
		{
			name:        "write file successfully",
			path:        "test.txt",
			content:     "hello world",
			shouldError: false,
		},
		{
			name:        "write file with nested directory",
			path:        "nested/dir/test.txt",
			content:     "nested content",
			shouldError: false,
		},
		{
			name:        "write empty file",
			path:        "empty.txt",
			content:     "",
			shouldError: false,
		},
		{
			name:        "write file with special characters",
			path:        "special.txt",
			content:     "content with\nnewlines\tand\ttabs",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			fullPath := filepath.Join(tempDir, tt.path)

			// Act
			err := WriteFile(fullPath, tt.content)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.FileExists(t, fullPath)

			// Verify content
			actualContent, err := os.ReadFile(fullPath)
			assert.NoError(t, err)
			assert.Equal(t, tt.content, string(actualContent))

			// Verify directory was created
			dir := filepath.Dir(fullPath)
			assert.DirExists(t, dir)
		})
	}
}

func TestFileExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tempDir string) string
		expected bool
	}{
		{
			name: "existing file returns true",
			setup: func(t *testing.T, tempDir string) string {
				path := filepath.Join(tempDir, "existing.txt")
				err := os.WriteFile(path, []byte("content"), 0644)
				require.NoError(t, err)
				return path
			},
			expected: true,
		},
		{
			name: "non-existing file returns false",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			expected: false,
		},
		{
			name: "directory returns true (FileExists checks any path)",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			// Act
			result := FileExists(path)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDirExists(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(t *testing.T, tempDir string) string
		expected bool
	}{
		{
			name: "existing directory returns true",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "existing_dir")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			expected: true,
		},
		{
			name: "non-existing directory returns false",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent_dir")
			},
			expected: false,
		},
		{
			name: "file returns false for directory check",
			setup: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "testfile.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				require.NoError(t, err)
				return filePath
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
			result := DirExists(path)

			// Assert
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateDir(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		shouldError   bool
		errorContains string
	}{
		{
			name:        "create directory successfully",
			path:        "testdir",
			shouldError: false,
		},
		{
			name:        "create nested directory",
			path:        "nested/deep/directory",
			shouldError: false,
		},
		{
			name:        "create directory that already exists",
			path:        "existing",
			shouldError: false, // Should not error if dir already exists
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			fullPath := filepath.Join(tempDir, tt.path)

			// Pre-create directory for "already exists" test
			if tt.name == "create directory that already exists" {
				err := os.MkdirAll(fullPath, 0755)
				require.NoError(t, err)
			}

			// Act
			err := CreateDir(fullPath, 0755)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.DirExists(t, fullPath)
		})
	}
}

func TestCopyDir(t *testing.T) {
	tests := []struct {
		name          string
		setupSrc      func(t *testing.T, tempDir string) string
		setupDst      func(t *testing.T, tempDir string) string
		shouldError   bool
		errorContains string
	}{
		{
			name: "copy directory successfully",
			setupSrc: func(t *testing.T, tempDir string) string {
				srcDir := filepath.Join(tempDir, "source")
				err := os.MkdirAll(srcDir, 0755)
				require.NoError(t, err)
				
				// Create test files
				err = os.WriteFile(filepath.Join(srcDir, "file1.txt"), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(srcDir, "file2.txt"), []byte("content2"), 0644)
				require.NoError(t, err)
				
				// Create subdirectory
				subDir := filepath.Join(srcDir, "subdir")
				err = os.MkdirAll(subDir, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(subDir, "nested.txt"), []byte("nested"), 0644)
				require.NoError(t, err)
				
				return srcDir
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "destination")
			},
			shouldError: false,
		},
		{
			name: "copy non-existent directory should error",
			setupSrc: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "destination")
			},
			shouldError:   true,
			errorContains: "failed to stat source directory",
		},
		{
			name: "copy file instead of directory should error",
			setupSrc: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "notadir.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				require.NoError(t, err)
				return filePath
			},
			setupDst: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "destination")
			},
			shouldError:   true,
			errorContains: "source is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			srcPath := tt.setupSrc(t, tempDir)
			dstPath := tt.setupDst(t, tempDir)

			// Act
			err := CopyDir(srcPath, dstPath)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.DirExists(t, dstPath)

			// Verify files were copied
			assert.FileExists(t, filepath.Join(dstPath, "file1.txt"))
			assert.FileExists(t, filepath.Join(dstPath, "file2.txt"))
			assert.FileExists(t, filepath.Join(dstPath, "subdir", "nested.txt"))

			// Verify content
			content1, err := os.ReadFile(filepath.Join(dstPath, "file1.txt"))
			assert.NoError(t, err)
			assert.Equal(t, "content1", string(content1))

			content2, err := os.ReadFile(filepath.Join(dstPath, "file2.txt"))
			assert.NoError(t, err)
			assert.Equal(t, "content2", string(content2))

			nestedContent, err := os.ReadFile(filepath.Join(dstPath, "subdir", "nested.txt"))
			assert.NoError(t, err)
			assert.Equal(t, "nested", string(nestedContent))
		})
	}
}

func TestReadFile(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) (string, string)
		shouldError   bool
		errorContains string
	}{
		{
			name: "read file successfully",
			setup: func(t *testing.T, tempDir string) (string, string) {
				content := "hello world"
				path := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(path, []byte(content), 0644)
				require.NoError(t, err)
				return path, content
			},
			shouldError: false,
		},
		{
			name: "read empty file",
			setup: func(t *testing.T, tempDir string) (string, string) {
				content := ""
				path := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(path, []byte(content), 0644)
				require.NoError(t, err)
				return path, content
			},
			shouldError: false,
		},
		{
			name: "read non-existent file should error",
			setup: func(t *testing.T, tempDir string) (string, string) {
				path := filepath.Join(tempDir, "nonexistent.txt")
				return path, ""
			},
			shouldError:   true,
			errorContains: "failed to read file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			tempDir := t.TempDir()
			path, expectedContent := tt.setup(t, tempDir)

			// Act
			actualContent, err := ReadFile(path)

			// Assert
			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, expectedContent, actualContent)
		})
	}
}