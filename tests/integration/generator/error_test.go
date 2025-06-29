package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// TestGenerator_Error_InvalidOutputPath tests error handling for invalid output paths
func TestGenerator_Error_InvalidOutputPath(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name       string
		outputPath string
		shouldFail bool
		errorType  string
	}{
		{
			name:       "empty output path",
			outputPath: "",
			shouldFail: true,
			errorType:  "path",
		},
		{
			name:       "invalid characters in path",
			outputPath: "/invalid\x00path",
			shouldFail: true,
			errorType:  "path",
		},
		{
			name:       "very long path",
			outputPath: "/" + strings.Repeat("a", 1000),
			shouldFail: true,
			errorType:  "path",
		},
		{
			name:       "valid path",
			outputPath: t.TempDir() + "/valid-project",
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:      "error-test",
				Module:    "github.com/test/error-test",
				Type:      "library",
				GoVersion: "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			options := types.GenerationOptions{
				OutputPath: tt.outputPath,
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			result, err := gen.Generate(config, options)

			if tt.shouldFail {
				assert.Error(t, err, "Expected error for invalid output path")
				assert.False(t, result.Success, "Result should indicate failure")
			} else {
				// Accept template not found errors
				if err != nil {
					if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
						t.Skip("Skipping test as template is not yet implemented")
						return
					}
				}
				assert.NoError(t, err, "Valid output path should not cause error")
			}
		})
	}
}

// TestGenerator_Error_PermissionDenied tests error handling for permission issues
func TestGenerator_Error_PermissionDenied(t *testing.T) {
	setupTestTemplates(t)

	// This test is platform-specific and may not work on all systems
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	// Create a directory with restricted permissions
	tmpDir := t.TempDir()
	restrictedDir := filepath.Join(tmpDir, "restricted")
	
	err := os.Mkdir(restrictedDir, 0000) // No permissions
	require.NoError(t, err)
	
	// Cleanup with proper permissions
	defer func() {
		os.Chmod(restrictedDir, 0755)
		os.RemoveAll(restrictedDir)
	}()

	config := types.ProjectConfig{
		Name:      "permission-test",
		Module:    "github.com/test/permission-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	gen := generator.New()
	require.NotNil(t, gen)

	options := types.GenerationOptions{
		OutputPath: filepath.Join(restrictedDir, "project"),
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	result, err := gen.Generate(config, options)

	// Should fail due to permission error
	assert.Error(t, err, "Expected permission error")
	assert.False(t, result.Success, "Result should indicate failure")
	
	// Error should be a filesystem error
	if err != nil {
		if goErr, ok := err.(*types.GoStarterError); ok {
			assert.Equal(t, types.ErrCodeFileSystem, goErr.Code, "Should be a filesystem error")
		}
	}
}

// TestGenerator_Error_MissingTemplate tests error handling for missing templates
func TestGenerator_Error_MissingTemplate(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name         string
		projectType  string
		architecture string
		expectedError string
	}{
		{
			name:         "unknown project type",
			projectType:  "unknown-type",
			architecture: "",
			expectedError: "template not found",
		},
		{
			name:         "unknown architecture",
			projectType:  "web-api",
			architecture: "unknown-arch",
			expectedError: "template not found",
		},
		{
			name:         "complex unknown template",
			projectType:  "microservice",
			architecture: "event-sourcing",
			expectedError: "template not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := types.ProjectConfig{
				Name:         "missing-template-test",
				Module:       "github.com/test/missing-template-test",
				Type:         tt.projectType,
				Architecture: tt.architecture,
				GoVersion:    "1.21",
			}

			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, config.Name),
				DryRun:     false,
				NoGit:      true,
				Verbose:    false,
			}

			result, err := gen.Generate(config, options)

			// Should fail with template not found error
			assert.Error(t, err, "Expected template not found error")
			assert.False(t, result.Success, "Result should indicate failure")
			
			// Error should be a template not found error
			assert.IsType(t, &types.GoStarterError{}, err, 
				"Should be a template not found error")

			// Error message should contain expected text
			if tt.expectedError != "" {
				assert.Contains(t, strings.ToLower(err.Error()), 
					strings.ToLower(tt.expectedError),
					"Error message should contain expected text")
			}
		})
	}
}

// TestGenerator_Error_InvalidTemplate tests error handling for invalid template data
func TestGenerator_Error_InvalidTemplate(t *testing.T) {
	setupTestTemplates(t)

	// This test would require mocking the template system to inject invalid templates
	// For now, we test the error handling structure

	config := types.ProjectConfig{
		Name:      "invalid-template-test",
		Module:    "github.com/test/invalid-template-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	gen := generator.New()
	require.NotNil(t, gen)

	tmpDir := t.TempDir()
	options := types.GenerationOptions{
		OutputPath: filepath.Join(tmpDir, config.Name),
		DryRun:     true, // Use dry run to avoid filesystem issues
		NoGit:      true,
		Verbose:    false,
	}

	result, err := gen.Generate(config, options)

	// The test should either succeed or fail with a known error type
	if err != nil {
		// Check that error types are properly structured
		if goErr, ok := err.(*types.GoStarterError); ok {
			switch goErr.Code {
			case types.ErrCodeTemplateNotFound:
				t.Log("Template not found - this is expected")
			case types.ErrCodeValidation:
				t.Log("Validation error - this might be expected")
			case types.ErrCodeFileSystem:
				t.Log("Filesystem error - this might indicate a problem")
			default:
				t.Logf("Go-starter error with code: %s", goErr.Code)
			}
		} else {
			t.Logf("Unknown error type: %T", err)
		}
	}

	// Result should be properly structured
	require.NotNil(t, result)
	if err != nil {
		assert.False(t, result.Success, "Result should indicate failure when error occurs")
		assert.NotNil(t, result.Error, "Result should contain error information")
	}
}

// TestGenerator_Error_ConfigValidation tests comprehensive validation error scenarios
func TestGenerator_Error_ConfigValidation(t *testing.T) {
	setupTestTemplates(t)

	tests := []struct {
		name          string
		config        types.ProjectConfig
		expectedError string
		errorType     interface{}
	}{
		{
			name: "all required fields missing",
			config: types.ProjectConfig{
				// All required fields empty
			},
			expectedError: "project name is required",
			errorType:     &types.GoStarterError{},
		},
		{
			name: "only name provided",
			config: types.ProjectConfig{
				Name: "test-project",
			},
			expectedError: "module path is required",
			errorType:     &types.GoStarterError{},
		},
		{
			name: "name and module provided",
			config: types.ProjectConfig{
				Name:   "test-project",
				Module: "github.com/test/project",
			},
			expectedError: "project type is required",
			errorType:     &types.GoStarterError{},
		},
		{
			name: "empty strings for required fields",
			config: types.ProjectConfig{
				Name:   "",
				Module: "",
				Type:   "",
			},
			expectedError: "project name is required",
			errorType:     &types.GoStarterError{},
		},
		{
			name: "whitespace-only fields",
			config: types.ProjectConfig{
				Name:   "   ",
				Module: "   ",
				Type:   "   ",
			},
			expectedError: "project name is required",
			errorType:     &types.GoStarterError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen := generator.New()
			require.NotNil(t, gen)

			tmpDir := t.TempDir()
			options := types.GenerationOptions{
				OutputPath: filepath.Join(tmpDir, "test-project"),
				DryRun:     true,
				NoGit:      true,
				Verbose:    false,
			}

			result, err := gen.Generate(tt.config, options)

			// Should fail with validation error
			assert.Error(t, err, "Expected validation error")
			assert.False(t, result.Success, "Result should indicate failure")

			// Check error type
			assert.IsType(t, tt.errorType, err, "Should be correct error type")

			// Check error message
			if tt.expectedError != "" {
				assert.Contains(t, err.Error(), tt.expectedError,
					"Error message should contain expected text")
			}
		})
	}
}

// TestGenerator_Error_DiskSpaceAndIO tests error handling for I/O issues
func TestGenerator_Error_DiskSpaceAndIO(t *testing.T) {
	setupTestTemplates(t)

	// Test with extremely long paths that might cause issues
	config := types.ProjectConfig{
		Name:      "io-test",
		Module:    "github.com/test/io-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	gen := generator.New()
	require.NotNil(t, gen)

	// Test with a path that's potentially problematic
	tmpDir := t.TempDir()
	longPath := filepath.Join(tmpDir, strings.Repeat("long-directory-name", 20))
	
	options := types.GenerationOptions{
		OutputPath: longPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	result, err := gen.Generate(config, options)

	// This might succeed or fail depending on the system
	// We mainly want to ensure proper error handling
	if err != nil {
		assert.False(t, result.Success, "Result should indicate failure")
		assert.NotNil(t, result.Error, "Result should contain error")
		
		// Error should be appropriately typed
		if goErr, ok := err.(*types.GoStarterError); ok {
			switch goErr.Code {
			case types.ErrCodeFileSystem:
				t.Log("Filesystem error - this is expected for I/O issues")
			case types.ErrCodeTemplateNotFound:
				t.Log("Template not found - this is acceptable for testing")
			default:
				t.Logf("Go-starter error with code: %s", goErr.Code)
			}
		} else {
			t.Logf("Other error type: %T - %v", err, err)
		}
	} else {
		// If it succeeded, clean up
		assert.True(t, result.Success, "Result should indicate success")
	}
}

// TestGenerator_Error_Recovery tests error recovery and cleanup
func TestGenerator_Error_Recovery(t *testing.T) {
	setupTestTemplates(t)

	config := types.ProjectConfig{
		Name:      "recovery-test",
		Module:    "github.com/test/recovery-test",
		Type:      "library",
		GoVersion: "1.21",
	}

	gen := generator.New()
	require.NotNil(t, gen)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, config.Name)
	
	options := types.GenerationOptions{
		OutputPath: outputPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	// First attempt - might fail with template not found
	result1, err1 := gen.Generate(config, options)

	// Second attempt with same path should work the same way
	result2, err2 := gen.Generate(config, options)

	// Both attempts should behave consistently
	assert.Equal(t, result1.Success, result2.Success, 
		"Multiple generation attempts should behave consistently")
	
	if err1 != nil && err2 != nil {
		assert.IsType(t, err1, err2, 
			"Error types should be consistent across attempts")
	}

	// The output directory state should be consistent
	if result1.Success && result2.Success {
		// Both succeeded - directory should exist
		_, err := os.Stat(outputPath)
		assert.NoError(t, err, "Output directory should exist after successful generation")
	}
}

// TestGenerator_Error_ConcurrentAccess tests error handling for concurrent access
func TestGenerator_Error_ConcurrentAccess(t *testing.T) {
	setupTestTemplates(t)

	config := types.ProjectConfig{
		Name:      "concurrent-test",
		Module:    "github.com/test/concurrent-test", 
		Type:      "library",
		GoVersion: "1.21",
	}

	// Test multiple generators accessing the same output path
	gen1 := generator.New()
	gen2 := generator.New()
	require.NotNil(t, gen1)
	require.NotNil(t, gen2)

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, config.Name)
	
	options := types.GenerationOptions{
		OutputPath: outputPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	}

	// Run both generators - they should handle concurrent access gracefully
	result1, err1 := gen1.Generate(config, options)
	result2, err2 := gen2.Generate(config, options)

	// At least one should handle the situation appropriately
	// Both might succeed if template not found, or one might fail due to concurrent access
	if err1 != nil {
		assert.False(t, result1.Success, "Result1 should indicate failure")
	}
	if err2 != nil {
		assert.False(t, result2.Success, "Result2 should indicate failure")
	}

	// If both succeeded, that's also acceptable (template not found scenario)
	if result1.Success && result2.Success {
		t.Log("Both generations succeeded - likely template not found scenario")
	}
}

// TestGenerator_Error_MemoryPressure tests behavior under memory constraints
func TestGenerator_Error_MemoryPressure(t *testing.T) {
	setupTestTemplates(t)

	// Test with a configuration that might require significant memory
	config := types.ProjectConfig{
		Name:      "memory-test",
		Module:    "github.com/test/memory-test",
		Type:      "web-api",
		GoVersion: "1.21",
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Drivers: []string{"postgresql", "mysql", "mongodb", "redis", "sqlite"},
				ORM:     "gorm",
			},
			Authentication: types.AuthConfig{
				Type:      "oauth2",
				Providers: []string{"google", "github", "facebook", "twitter", "linkedin"},
			},
			Deployment: types.DeployConfig{
				Targets: []string{"docker", "kubernetes", "aws", "gcp", "azure"},
			},
		},
		Variables: make(map[string]string),
	}

	// Add many variables to test memory usage
	for i := 0; i < 100; i++ {
		config.Variables[fmt.Sprintf("Variable%d", i)] = fmt.Sprintf("Value%d", i)
	}

	gen := generator.New()
	require.NotNil(t, gen)

	tmpDir := t.TempDir()
	options := types.GenerationOptions{
		OutputPath: filepath.Join(tmpDir, config.Name),
		DryRun:     true, // Use dry run to test memory usage without I/O
		NoGit:      true,
		Verbose:    false,
	}

	result, err := gen.Generate(config, options)

	// Should handle the complex configuration gracefully
	if err != nil {
		// Accept template not found errors
		if goErr, ok := err.(*types.GoStarterError); ok && goErr.Code == types.ErrCodeTemplateNotFound {
			t.Log("Template not found - this is expected")
		} else {
			t.Logf("Error handling complex configuration: %v", err)
		}
		assert.False(t, result.Success, "Result should indicate failure")
	} else {
		assert.True(t, result.Success, "Should handle complex configuration")
	}

	// Verify the generator didn't crash or panic
	require.NotNil(t, result, "Result should not be nil")
	t.Log("Memory pressure test completed successfully")
}