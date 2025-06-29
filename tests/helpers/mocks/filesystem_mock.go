package mocks

import (
	"io/fs"

	"github.com/stretchr/testify/mock"
)

// MockFileSystem is a mock implementation of a file system interface.
type MockFileSystem struct {
	mock.Mock
}

// ReadFile provides a mock function with given fields: name
func (m *MockFileSystem) ReadFile(name string) ([]byte, error) {
	args := m.Called(name)
	return args.Get(0).([]byte), args.Error(1)
}

// WriteFile provides a mock function with given fields: name, data, perm
func (m *MockFileSystem) WriteFile(name string, data []byte, perm fs.FileMode) error {
	args := m.Called(name, data, perm)
	return args.Error(0)
}

// MkdirAll provides a mock function with given fields: path, perm
func (m *MockFileSystem) MkdirAll(path string, perm fs.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}

// Stat provides a mock function with given fields: name
func (m *MockFileSystem) Stat(name string) (fs.FileInfo, error) {
	args := m.Called(name)
	return args.Get(0).(fs.FileInfo), args.Error(1)
}

// RemoveAll provides a mock function with given fields: path
func (m *MockFileSystem) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}
