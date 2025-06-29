package mocks

import (
	"github.com/stretchr/testify/mock"
	"{{.ModulePath}}/internal/domain/ports"
)

// MockLogger is a mock implementation of ports.Logger
type MockLogger struct {
	mock.Mock
}

// Debug provides a mock function with given fields: msg, fields
func (m *MockLogger) Debug(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// Info provides a mock function with given fields: msg, fields
func (m *MockLogger) Info(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// Warn provides a mock function with given fields: msg, fields
func (m *MockLogger) Warn(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// Error provides a mock function with given fields: msg, fields
func (m *MockLogger) Error(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// Fatal provides a mock function with given fields: msg, fields
func (m *MockLogger) Fatal(msg string, fields ...interface{}) {
	m.Called(msg, fields)
}

// With provides a mock function with given fields: fields
func (m *MockLogger) With(fields ...interface{}) ports.Logger {
	args := m.Called(fields)
	return args.Get(0).(ports.Logger)
}
