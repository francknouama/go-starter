package mocks

import (
	"github.com/stretchr/testify/mock"
	"github.com/francknouama/go-starter/pkg/types"
)

// MockTemplateRegistry is a mock implementation of the template registry.
type MockTemplateRegistry struct {
	mock.Mock
}

// Register provides a mock function with given fields: template
func (m *MockTemplateRegistry) Register(template types.Template) error {
	args := m.Called(template)
	return args.Error(0)
}

// Get provides a mock function with given fields: templateID
func (m *MockTemplateRegistry) Get(templateID string) (types.Template, error) {
	args := m.Called(templateID)
	return args.Get(0).(types.Template), args.Error(1)
}

// List provides a mock function with no fields.
func (m *MockTemplateRegistry) List() []types.Template {
	args := m.Called()
	return args.Get(0).([]types.Template)
}

// GetByType provides a mock function with given fields: templateType
func (m *MockTemplateRegistry) GetByType(templateType string) []types.Template {
	args := m.Called(templateType)
	return args.Get(0).([]types.Template)
}

// Exists provides a mock function with given fields: templateID
func (m *MockTemplateRegistry) Exists(templateID string) bool {
	args := m.Called(templateID)
	return args.Bool(0)
}

// Remove provides a mock function with given fields: templateID
func (m *MockTemplateRegistry) Remove(templateID string) error {
	args := m.Called(templateID)
	return args.Error(0)
}

// GetTemplateTypes provides a mock function with no fields.
func (m *MockTemplateRegistry) GetTemplateTypes() []string {
	args := m.Called()
	return args.Get(0).([]string)
}
