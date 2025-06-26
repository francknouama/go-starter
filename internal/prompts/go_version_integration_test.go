package prompts

import (
	"path/filepath"
	"testing"

	"github.com/francknouama/go-starter/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestGoVersionIntegration(t *testing.T) {
	tests := []struct {
		name          string
		goVersion     string
		expectedError bool
	}{
		{
			name:          "valid auto version",
			goVersion:     "auto",
			expectedError: false,
		},
		{
			name:          "valid specific version 1.23",
			goVersion:     "1.23",
			expectedError: false,
		},
		{
			name:          "valid specific version 1.22",
			goVersion:     "1.22",
			expectedError: false,
		},
		{
			name:          "valid specific version 1.21",
			goVersion:     "1.21",
			expectedError: false,
		},
		{
			name:          "invalid version",
			goVersion:     "1.20",
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			config := types.ProjectConfig{
				Name:       "test-project",
				Module:     "github.com/test/project",
				Type:       "web-api",
				GoVersion:  tt.goVersion,
				Framework:  "gin",
				Logger:     "slog",
			}

			// Act
			err := ValidateGoVersion(config.GoVersion)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "unsupported Go version")
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGoVersionTemplateIntegration(t *testing.T) {
	// This test verifies that Go version is properly integrated into templates
	// For now, it's a placeholder for future template integration tests
	t.Run("go version used in template generation", func(t *testing.T) {
		// Arrange
		config := types.ProjectConfig{
			Name:       "test-project",
			Module:     "github.com/test/project",
			Type:       "web-api",
			GoVersion:  "1.23",
			Framework:  "gin",
			Logger:     "slog",
		}

		tempDir := t.TempDir()

		// This would require the actual generator to be available
		// For now, we'll just test the configuration validation
		err := ValidateGoVersion(config.GoVersion)
		assert.NoError(t, err)

		// When template integration is complete, we would:
		// 1. Generate project with specific Go version
		// 2. Verify go.mod contains correct version
		// 3. Verify project compiles

		// Placeholder assertions
		goModPath := filepath.Join(tempDir, "go.mod")
		expectedFiles := []string{"go.mod", "main.go", "Makefile"}

		// These would be used when generator integration is complete:
		_ = goModPath
		_ = expectedFiles
		// helpers.AssertGoModContainsVersion(t, goModPath, "1.23")
		// helpers.AssertProjectGenerated(t, tempDir, expectedFiles)
	})
}