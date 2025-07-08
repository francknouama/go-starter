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
		name          string
		setupSrc      func(t *testing.T, tempDir string) string
		setupDst      func(t *testing.T, tempDir string) string
		content       string
		mode          os.FileMode
		shouldError   bool
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

func TestFileExistsWithError(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedExists bool
		shouldError   bool
	}{
		{
			name: "existing file returns true with no error",
			setup: func(t *testing.T, tempDir string) string {
				path := filepath.Join(tempDir, "exists.txt")
				err := os.WriteFile(path, []byte("content"), 0644)
				require.NoError(t, err)
				return path
			},
			expectedExists: true,
			shouldError:    false,
		},
		{
			name: "non-existing file returns false with no error",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			expectedExists: false,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			exists, err := FileExistsWithError(path)

			assert.Equal(t, tt.expectedExists, exists)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDirExistsWithError(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedExists bool
		shouldError   bool
	}{
		{
			name: "existing directory returns true with no error",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			expectedExists: true,
			shouldError:    false,
		},
		{
			name: "non-existing directory returns false with no error",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectedExists: false,
			shouldError:    false,
		},
		{
			name: "file path returns false with no error",
			setup: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "file.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				require.NoError(t, err)
				return filePath
			},
			expectedExists: false,
			shouldError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			exists, err := DirExistsWithError(path)

			assert.Equal(t, tt.expectedExists, exists)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestIsEmptyDir(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedEmpty bool
		shouldError   bool
	}{
		{
			name: "empty directory returns true",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "empty")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			expectedEmpty: true,
			shouldError:   false,
		},
		{
			name: "directory with files returns false",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "nonempty")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(dirPath, "file.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				return dirPath
			},
			expectedEmpty: false,
			shouldError:   false,
		},
		{
			name: "non-existent directory returns error",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectedEmpty: false,
			shouldError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			isEmpty, err := IsEmptyDir(path)

			assert.Equal(t, tt.expectedEmpty, isEmpty)
			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRemoveDir(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		shouldError   bool
		errorContains string
	}{
		{
			name: "remove empty directory successfully",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "toremove")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			shouldError: false,
		},
		{
			name: "remove directory with contents successfully",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "toremove")
				err := os.MkdirAll(filepath.Join(dirPath, "subdir"), 0755)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(dirPath, "file.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				return dirPath
			},
			shouldError: false,
		},
		{
			name: "remove non-existent directory succeeds",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			dirPath := tt.setup(t, tempDir)

			err := RemoveDir(dirPath)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.NoDirExists(t, dirPath)
			}
		})
	}
}

func TestGetFileExt(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "file with extension",
			path:     "file.txt",
			expected: ".txt",
		},
		{
			name:     "file with multiple dots",
			path:     "file.test.go",
			expected: ".go",
		},
		{
			name:     "file without extension",
			path:     "README",
			expected: "",
		},
		{
			name:     "path with directories",
			path:     "/path/to/file.json",
			expected: ".json",
		},
		{
			name:     "hidden file with extension",
			path:     ".gitignore",
			expected: ".gitignore",
		},
		{
			name:     "hidden file with extension after dot",
			path:     ".config.yaml",
			expected: ".yaml",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExt(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetBaseName(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple filename",
			path:     "file.txt",
			expected: "file",
		},
		{
			name:     "path with directories",
			path:     "/path/to/file.go",
			expected: "file",
		},
		{
			name:     "file without extension",
			path:     "README",
			expected: "README",
		},
		{
			name:     "multiple dots",
			path:     "file.test.json",
			expected: "file.test",
		},
		{
			name:     "hidden file",
			path:     ".gitignore",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetBaseName(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCleanPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "simple/path",
			expected: "simple/path",
		},
		{
			name:     "path with dot notation",
			path:     "path/./to/file",
			expected: "path/to/file",
		},
		{
			name:     "path with double dots",
			path:     "path/../to/file",
			expected: "to/file",
		},
		{
			name:     "complex path",
			path:     "/path/./to/../from/file",
			expected: "/path/from/file",
		},
		{
			name:     "empty path",
			path:     "",
			expected: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanPath(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestJoinPath(t *testing.T) {
	tests := []struct {
		name     string
		paths    []string
		expected string
	}{
		{
			name:     "join two paths",
			paths:    []string{"path", "to"},
			expected: filepath.Join("path", "to"),
		},
		{
			name:     "join multiple paths",
			paths:    []string{"path", "to", "file.txt"},
			expected: filepath.Join("path", "to", "file.txt"),
		},
		{
			name:     "join with empty string",
			paths:    []string{"path", "", "file"},
			expected: filepath.Join("path", "", "file"),
		},
		{
			name:     "single path",
			paths:    []string{"single"},
			expected: "single",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := JoinPath(tt.paths...)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		wantDir string
		wantFile string
	}{
		{
			name:     "simple path",
			path:     "path/to/file.txt",
			wantDir:  "path/to/",
			wantFile: "file.txt",
		},
		{
			name:     "root file",
			path:     "/file.txt",
			wantDir:  "/",
			wantFile: "file.txt",
		},
		{
			name:     "current directory file",
			path:     "file.txt",
			wantDir:  "",
			wantFile: "file.txt",
		},
		{
			name:     "directory path",
			path:     "path/to/",
			wantDir:  "path/to/",
			wantFile: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, file := SplitPath(tt.path)
			assert.Equal(t, tt.wantDir, dir)
			assert.Equal(t, tt.wantFile, file)
		})
	}
}

func TestIsHidden(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "hidden file",
			path:     ".hidden",
			expected: true,
		},
		{
			name:     "hidden file in path",
			path:     "path/to/.hidden",
			expected: true,
		},
		{
			name:     "regular file",
			path:     "visible.txt",
			expected: false,
		},
		{
			name:     "regular file in path",
			path:     "path/to/visible.txt",
			expected: false,
		},
		{
			name:     "file starting with dot but not hidden",
			path:     "path/.not/visible.txt",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsHidden(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetFileSize(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedSize  int64
		shouldError   bool
		errorContains string
	}{
		{
			name: "get size of existing file",
			setup: func(t *testing.T, tempDir string) string {
				path := filepath.Join(tempDir, "test.txt")
				content := "hello world"
				err := os.WriteFile(path, []byte(content), 0644)
				require.NoError(t, err)
				return path
			},
			expectedSize: 11, // "hello world" is 11 bytes
			shouldError:  false,
		},
		{
			name: "get size of empty file",
			setup: func(t *testing.T, tempDir string) string {
				path := filepath.Join(tempDir, "empty.txt")
				err := os.WriteFile(path, []byte(""), 0644)
				require.NoError(t, err)
				return path
			},
			expectedSize: 0,
			shouldError:  false,
		},
		{
			name: "get size of non-existent file returns error",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			expectedSize:  0,
			shouldError:   true,
			errorContains: "failed to stat file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			size, err := GetFileSize(path)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSize, size)
			}
		})
	}
}

func TestGetFileMode(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) (string, os.FileMode)
		shouldError   bool
		errorContains string
	}{
		{
			name: "get mode of existing file",
			setup: func(t *testing.T, tempDir string) (string, os.FileMode) {
				path := filepath.Join(tempDir, "test.txt")
				mode := os.FileMode(0644)
				err := os.WriteFile(path, []byte("content"), mode)
				require.NoError(t, err)
				return path, mode
			},
			shouldError: false,
		},
		{
			name: "get mode of executable file",
			setup: func(t *testing.T, tempDir string) (string, os.FileMode) {
				path := filepath.Join(tempDir, "script.sh")
				mode := os.FileMode(0755)
				err := os.WriteFile(path, []byte("#!/bin/bash"), mode)
				require.NoError(t, err)
				return path, mode
			},
			shouldError: false,
		},
		{
			name: "get mode of non-existent file returns error",
			setup: func(t *testing.T, tempDir string) (string, os.FileMode) {
				return filepath.Join(tempDir, "nonexistent.txt"), 0
			},
			shouldError:   true,
			errorContains: "failed to stat file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path, expectedMode := tt.setup(t, tempDir)

			mode, err := GetFileMode(path)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expectedMode, mode.Perm())
			}
		})
	}
}

func TestSetFileMode(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		newMode       os.FileMode
		shouldError   bool
		errorContains string
	}{
		{
			name: "set mode of existing file",
			setup: func(t *testing.T, tempDir string) string {
				path := filepath.Join(tempDir, "test.txt")
				err := os.WriteFile(path, []byte("content"), 0644)
				require.NoError(t, err)
				return path
			},
			newMode:     0755,
			shouldError: false,
		},
		{
			name: "set mode of non-existent file returns error",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent.txt")
			},
			newMode:       0755,
			shouldError:   true,
			errorContains: "failed to set file mode",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setup(t, tempDir)

			err := SetFileMode(path, tt.newMode)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				// Verify the mode was set correctly
				info, statErr := os.Stat(path)
				assert.NoError(t, statErr)
				assert.Equal(t, tt.newMode, info.Mode().Perm())
			}
		})
	}
}

func TestCreateTempDir(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{
			name:   "create temp directory with prefix",
			prefix: "test_",
		},
		{
			name:   "create temp directory with empty prefix",
			prefix: "",
		},
		{
			name:   "create temp directory with special chars in prefix",
			prefix: "test-prefix_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dirPath, err := CreateTempDir(tt.prefix)

			assert.NoError(t, err)
			assert.NotEmpty(t, dirPath)
			assert.DirExists(t, dirPath)
			
			if tt.prefix != "" {
				assert.Contains(t, filepath.Base(dirPath), tt.prefix)
			}

			// Clean up
			defer func() { _ = os.RemoveAll(dirPath) }()
		})
	}
}

func TestCreateTempFile(t *testing.T) {
	tests := []struct {
		name   string
		prefix string
	}{
		{
			name:   "create temp file with prefix",
			prefix: "test_",
		},
		{
			name:   "create temp file with empty prefix",
			prefix: "",
		},
		{
			name:   "create temp file with special chars in prefix",
			prefix: "test-prefix_",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file, err := CreateTempFile(tt.prefix)

			assert.NoError(t, err)
			assert.NotNil(t, file)
			
			filePath := file.Name()
			assert.NotEmpty(t, filePath)
			assert.FileExists(t, filePath)
			
			if tt.prefix != "" {
				assert.Contains(t, filepath.Base(filePath), tt.prefix)
			}

			// Clean up
			_ = file.Close()
			defer func() { _ = os.Remove(filePath) }()
		})
	}
}

func TestEnsureDir(t *testing.T) {
	tests := []struct {
		name          string
		setupPath     func(t *testing.T, tempDir string) string
		shouldError   bool
		errorContains string
	}{
		{
			name: "ensure non-existent directory",
			setupPath: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "new", "nested", "dir")
			},
			shouldError: false,
		},
		{
			name: "ensure existing directory",
			setupPath: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "existing")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			shouldError: false,
		},
		{
			name: "ensure directory where file exists",
			setupPath: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "file.txt")
				err := os.WriteFile(filePath, []byte("content"), 0644)
				require.NoError(t, err)
				return filePath
			},
			shouldError:   true,
			errorContains: "failed to create directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			path := tt.setupPath(t, tempDir)

			err := EnsureDir(path)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.DirExists(t, path)
			}
		})
	}
}

func TestGetRelativePath(t *testing.T) {
	tests := []struct {
		name          string
		basePath      string
		targetPath    string
		expected      string
		shouldError   bool
		errorContains string
	}{
		{
			name:        "relative path from base to target",
			basePath:    "/home/user",
			targetPath:  "/home/user/projects/test",
			expected:    "projects/test",
			shouldError: false,
		},
		{
			name:        "relative path to parent directory",
			basePath:    "/home/user/projects",
			targetPath:  "/home/user",
			expected:    "..",
			shouldError: false,
		},
		{
			name:        "same directory",
			basePath:    "/home/user",
			targetPath:  "/home/user",
			expected:    ".",
			shouldError: false,
		},
		{
			name:        "relative path with common ancestor",
			basePath:    "/home/user/projects/app1",
			targetPath:  "/home/user/projects/app2",
			expected:    "../app2",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetRelativePath(tt.basePath, tt.targetPath)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestListFiles(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedFiles []string
		shouldError   bool
		errorContains string
	}{
		{
			name: "list files in directory with files",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				
				// Create files
				err = os.WriteFile(filepath.Join(dirPath, "file1.txt"), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(dirPath, "file2.go"), []byte("content2"), 0644)
				require.NoError(t, err)
				
				// Create subdirectory (should not be included in files list)
				err = os.Mkdir(filepath.Join(dirPath, "subdir"), 0755)
				require.NoError(t, err)
				
				return dirPath
			},
			expectedFiles: []string{"file1.txt", "file2.go"},
			shouldError:   false,
		},
		{
			name: "list files in empty directory",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "empty")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				return dirPath
			},
			expectedFiles: []string{},
			shouldError:   false,
		},
		{
			name: "list files in non-existent directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectedFiles: nil,
			shouldError:   true,
			errorContains: "failed to read directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			dirPath := tt.setup(t, tempDir)

			files, err := ListFiles(dirPath)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedFiles, files)
			}
		})
	}
}

func TestListDirs(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(t *testing.T, tempDir string) string
		expectedDirs []string
		shouldError  bool
		errorContains string
	}{
		{
			name: "list directories in directory with subdirs",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "testdir")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				
				// Create subdirectories
				err = os.Mkdir(filepath.Join(dirPath, "subdir1"), 0755)
				require.NoError(t, err)
				err = os.Mkdir(filepath.Join(dirPath, "subdir2"), 0755)
				require.NoError(t, err)
				
				// Create file (should not be included in dirs list)
				err = os.WriteFile(filepath.Join(dirPath, "file.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return dirPath
			},
			expectedDirs: []string{"subdir1", "subdir2"},
			shouldError:  false,
		},
		{
			name: "list directories in directory with no subdirs",
			setup: func(t *testing.T, tempDir string) string {
				dirPath := filepath.Join(tempDir, "nosubdirs")
				err := os.Mkdir(dirPath, 0755)
				require.NoError(t, err)
				
				// Create only files
				err = os.WriteFile(filepath.Join(dirPath, "file1.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(dirPath, "file2.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return dirPath
			},
			expectedDirs: []string{},
			shouldError:  false,
		},
		{
			name: "list directories in non-existent directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectedDirs:  nil,
			shouldError:   true,
			errorContains: "failed to read directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			dirPath := tt.setup(t, tempDir)

			dirs, err := ListDirs(dirPath)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedDirs, dirs)
			}
		})
	}
}

func TestSafeWriteFile(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		content       string
		shouldError   bool
		errorContains string
	}{
		{
			name: "safely write file to new location",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "safe.txt")
			},
			content:     "safe content",
			shouldError: false,
		},
		{
			name: "safely write file to nested directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nested", "deep", "safe.txt")
			},
			content:     "nested safe content",
			shouldError: false,
		},
		{
			name: "safely overwrite existing file",
			setup: func(t *testing.T, tempDir string) string {
				filePath := filepath.Join(tempDir, "existing.txt")
				err := os.WriteFile(filePath, []byte("old content"), 0644)
				require.NoError(t, err)
				return filePath
			},
			content:     "new content",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			filePath := tt.setup(t, tempDir)

			err := SafeWriteFile(filePath, tt.content)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.FileExists(t, filePath)
				
				// Verify content
				actualContent, readErr := os.ReadFile(filePath)
				assert.NoError(t, readErr)
				assert.Equal(t, tt.content, string(actualContent))
			}
		})
	}
}

func TestWalkDir(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		expectedPaths []string
		shouldError   bool
		errorContains string
	}{
		{
			name: "walk directory with files and subdirs",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "walktest")
				err := os.MkdirAll(filepath.Join(rootDir, "subdir"), 0755)
				require.NoError(t, err)
				
				// Create files
				err = os.WriteFile(filepath.Join(rootDir, "file1.txt"), []byte("content1"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(rootDir, "subdir", "file2.txt"), []byte("content2"), 0644)
				require.NoError(t, err)
				
				return rootDir
			},
			expectedPaths: []string{"file1.txt", "subdir", "subdir/file2.txt"},
			shouldError:   false,
		},
		{
			name: "walk empty directory",
			setup: func(t *testing.T, tempDir string) string {
				emptyDir := filepath.Join(tempDir, "empty")
				err := os.Mkdir(emptyDir, 0755)
				require.NoError(t, err)
				return emptyDir
			},
			expectedPaths: []string{},
			shouldError:   false,
		},
		{
			name: "walk non-existent directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			expectedPaths: nil,
			shouldError:   true,
			errorContains: "lstat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			rootDir := tt.setup(t, tempDir)

			var foundPaths []string
			walkFunc := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				// Get relative path from root
				relPath, _ := filepath.Rel(rootDir, path)
				if relPath != "." {
					foundPaths = append(foundPaths, relPath)
				}
				return nil
			}

			err := WalkDir(rootDir, walkFunc)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				assert.ElementsMatch(t, tt.expectedPaths, foundPaths)
			}
		})
	}
}

func TestFindFiles(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T, tempDir string) string
		pattern       string
		expectedFiles []string
		shouldError   bool
		errorContains string
	}{
		{
			name: "find text files",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "findtest")
				err := os.MkdirAll(filepath.Join(rootDir, "subdir"), 0755)
				require.NoError(t, err)
				
				// Create various files
				err = os.WriteFile(filepath.Join(rootDir, "file1.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(rootDir, "file2.go"), []byte("content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(rootDir, "subdir", "file3.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return rootDir
			},
			pattern:       "*.txt",
			expectedFiles: []string{"file1.txt", "subdir/file3.txt"},
			shouldError:   false,
		},
		{
			name: "find go files",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "findtest")
				err := os.MkdirAll(rootDir, 0755)
				require.NoError(t, err)
				
				err = os.WriteFile(filepath.Join(rootDir, "main.go"), []byte("content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(rootDir, "readme.txt"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return rootDir
			},
			pattern:       "*.go",
			expectedFiles: []string{"main.go"},
			shouldError:   false,
		},
		{
			name: "find files in non-existent directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			pattern:       "*.txt",
			expectedFiles: nil,
			shouldError:   true,
			errorContains: "lstat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			rootDir := tt.setup(t, tempDir)

			files, err := FindFiles(rootDir, tt.pattern)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
				
				// Convert absolute paths to relative for comparison
				var relFiles []string
				for _, file := range files {
					relPath, _ := filepath.Rel(rootDir, file)
					relFiles = append(relFiles, relPath)
				}
				assert.ElementsMatch(t, tt.expectedFiles, relFiles)
			}
		})
	}
}

func TestValidateDirectoryStructure(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(t *testing.T, tempDir string) string
		requiredDirs    []string
		requiredFiles   []string
		shouldError     bool
		errorContains   string
	}{
		{
			name: "validate complete structure",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "project")
				err := os.MkdirAll(filepath.Join(rootDir, "src"), 0755)
				require.NoError(t, err)
				err = os.MkdirAll(filepath.Join(rootDir, "tests"), 0755)
				require.NoError(t, err)
				
				err = os.WriteFile(filepath.Join(rootDir, "main.go"), []byte("content"), 0644)
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(rootDir, "go.mod"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return rootDir
			},
			requiredDirs:  []string{"src", "tests"},
			requiredFiles: []string{"main.go", "go.mod"},
			shouldError:   false,
		},
		{
			name: "validate with missing directory",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "project")
				err := os.MkdirAll(filepath.Join(rootDir, "src"), 0755)
				require.NoError(t, err)
				// Missing "tests" directory
				
				err = os.WriteFile(filepath.Join(rootDir, "main.go"), []byte("content"), 0644)
				require.NoError(t, err)
				
				return rootDir
			},
			requiredDirs:    []string{"src", "tests"},
			requiredFiles:   []string{"main.go"},
			shouldError:     true,
			errorContains:   "expected directory",
		},
		{
			name: "validate with missing file",
			setup: func(t *testing.T, tempDir string) string {
				rootDir := filepath.Join(tempDir, "project")
				err := os.MkdirAll(filepath.Join(rootDir, "src"), 0755)
				require.NoError(t, err)
				
				err = os.WriteFile(filepath.Join(rootDir, "main.go"), []byte("content"), 0644)
				require.NoError(t, err)
				// Missing "go.mod" file
				
				return rootDir
			},
			requiredDirs:    []string{"src"},
			requiredFiles:   []string{"main.go", "go.mod"},
			shouldError:     true,
			errorContains:   "expected file",
		},
		{
			name: "validate non-existent root directory",
			setup: func(t *testing.T, tempDir string) string {
				return filepath.Join(tempDir, "nonexistent")
			},
			requiredDirs:    []string{"src"},
			requiredFiles:   []string{"main.go"},
			shouldError:     true,
			errorContains:   "expected directory src does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()
			rootDir := tt.setup(t, tempDir)

			err := ValidateDirectoryStructure(rootDir, tt.requiredDirs, tt.requiredFiles)

			if tt.shouldError {
				assert.Error(t, err)
				if tt.errorContains != "" {
					assert.Contains(t, err.Error(), tt.errorContains)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
