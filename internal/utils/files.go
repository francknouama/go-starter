package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// CopyFile copies a file from src to dst
func CopyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Printf("Warning: failed to close source file: %v\n", err)
		}
	}()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			fmt.Printf("Warning: failed to close destination file: %v\n", err)
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	// Copy file permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to get source file info: %w", err)
	}

	return os.Chmod(dst, srcInfo.Mode())
}

// CopyDir copies a directory recursively from src to dst
func CopyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("failed to stat source directory: %w", err)
	}

	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("failed to read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := CopyDir(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy subdirectory %s: %w", entry.Name(), err)
			}
		} else {
			if err := CopyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("failed to copy file %s: %w", entry.Name(), err)
			}
		}
	}

	return nil
}

// WriteFile writes content to a file, creating directories as needed
func WriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// ReadFile reads content from a file
func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(content), nil
}

// FileExists checks if a file exists and is accessible
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// FileExistsWithError checks if a file exists and returns any access errors
func FileExistsWithError(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file existence for %q: %w", path, err)
}

// DirExists checks if a directory exists and is accessible
func DirExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// DirExistsWithError checks if a directory exists and returns any access errors
func DirExistsWithError(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check directory existence for %q: %w", path, err)
}

// IsEmptyDir checks if a directory is empty
func IsEmptyDir(path string) (bool, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return false, fmt.Errorf("failed to read directory: %w", err)
	}
	return len(entries) == 0, nil
}

// RemoveDir removes a directory and all its contents
func RemoveDir(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("failed to remove directory: %w", err)
	}
	return nil
}

// CreateDir creates a directory with the specified permissions
func CreateDir(path string, perm os.FileMode) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}
	return nil
}

// WalkDir walks through a directory and calls walkFn for each file/directory
func WalkDir(root string, walkFn filepath.WalkFunc) error {
	return filepath.Walk(root, walkFn)
}

// FindFiles finds files matching a pattern in a directory
func FindFiles(root, pattern string) ([]string, error) {
	var matches []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err != nil {
			return err
		}

		if matched {
			matches = append(matches, path)
		}

		return nil
	})

	return matches, err
}

// GetRelativePath returns the relative path from base to target
func GetRelativePath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	return rel, nil
}

// CleanPath cleans and normalizes a file path
func CleanPath(path string) string {
	return filepath.Clean(path)
}

// JoinPath joins path elements
func JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

// SplitPath splits a path into directory and file
func SplitPath(path string) (dir, file string) {
	return filepath.Split(path)
}

// GetFileExt returns the file extension
func GetFileExt(path string) string {
	return filepath.Ext(path)
}

// GetBaseName returns the base name of a path (without extension)
func GetBaseName(path string) string {
	base := filepath.Base(path)
	ext := filepath.Ext(base)
	return strings.TrimSuffix(base, ext)
}

// IsHidden checks if a file or directory is hidden (starts with .)
func IsHidden(path string) bool {
	return strings.HasPrefix(filepath.Base(path), ".")
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.Size(), nil
}

// GetFileMode returns the file permissions
func GetFileMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("failed to stat file: %w", err)
	}
	return info.Mode(), nil
}

// SetFileMode sets the file permissions
func SetFileMode(path string, mode os.FileMode) error {
	if err := os.Chmod(path, mode); err != nil {
		return fmt.Errorf("failed to set file mode: %w", err)
	}
	return nil
}

// CreateTempDir creates a temporary directory
func CreateTempDir(prefix string) (string, error) {
	dir, err := os.MkdirTemp("", prefix)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	return dir, nil
}

// CreateTempFile creates a temporary file
func CreateTempFile(prefix string) (*os.File, error) {
	file, err := os.CreateTemp("", prefix)
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	return file, nil
}

// EnsureDir ensures a directory exists, creating it if necessary
func EnsureDir(path string) error {
	if !DirExists(path) {
		if err := CreateDir(path, 0755); err != nil {
			return err
		}
	}
	return nil
}

// SafeWriteFile writes content to a file safely (atomic write)
func SafeWriteFile(path, content string) error {
	dir := filepath.Dir(path)
	if err := EnsureDir(dir); err != nil {
		return err
	}

	// Create temporary file in the same directory
	tmpFile, err := os.CreateTemp(dir, ".tmp-"+filepath.Base(path))
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Printf("Warning: failed to remove temporary file: %v\n", err)
		}
	}() // Clean up if we fail

	// Write content to temp file
	if _, err := tmpFile.WriteString(content); err != nil {
		if closeErr := tmpFile.Close(); closeErr != nil {
			fmt.Printf("Warning: failed to close temp file: %v\n", closeErr)
		}
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	// Close temp file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomically replace the target file
	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("failed to replace file: %w", err)
	}

	return nil
}

// ListFiles lists all files in a directory (non-recursive)
func ListFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var files []string
	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}

	return files, nil
}

// ListDirs lists all directories in a directory (non-recursive)
func ListDirs(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory: %w", err)
	}

	var dirs []string
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry.Name())
		}
	}

	return dirs, nil
}

// ValidateDirectoryStructure validates that a directory structure matches expectations
func ValidateDirectoryStructure(root string, expectedDirs []string, expectedFiles []string) error {
	// Check expected directories exist
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(root, dir)
		if !DirExists(dirPath) {
			return fmt.Errorf("expected directory %s does not exist", dir)
		}
	}

	// Check expected files exist
	for _, file := range expectedFiles {
		filePath := filepath.Join(root, file)
		if !FileExists(filePath) {
			return fmt.Errorf("expected file %s does not exist", file)
		}
	}

	return nil
}

// CopyFileFromFS copies a file from an embedded filesystem to disk
func CopyFileFromFS(fsys fs.FS, src, dst string) error {
	srcFile, err := fsys.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open source file from FS: %w", err)
	}
	defer func() {
		if err := srcFile.Close(); err != nil {
			fmt.Printf("Warning: failed to close source file from FS: %v\n", err)
		}
	}()

	// Create destination directory if it doesn't exist
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func() {
		if err := dstFile.Close(); err != nil {
			fmt.Printf("Warning: failed to close destination file: %v\n", err)
		}
	}()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file from FS: %w", err)
	}

	return nil
}
