package optimization

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.Equal(t, OptimizationLevelStandard, config.Level)
	assert.Equal(t, "balanced", config.ProfileName)
	assert.Equal(t, "1.0.0", config.Version)
	assert.NotEmpty(t, config.Description)
	assert.NotNil(t, config.Options)
}

func TestConfig_SaveAndLoad(t *testing.T) {
	// Create a temporary file for testing
	tempDir, err := os.MkdirTemp("", "config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "test-config.json")

	// Create and save a config
	originalConfig := DefaultConfig()
	originalConfig.Level = OptimizationLevelAggressive
	originalConfig.ProfileName = "performance"
	originalConfig.Description = "Test configuration"

	err = originalConfig.SaveConfig(configPath)
	require.NoError(t, err)

	// Verify file was created
	_, err = os.Stat(configPath)
	require.NoError(t, err)

	// Load the config back
	loadedConfig, err := LoadConfig(configPath)
	require.NoError(t, err)

	// Verify the loaded config matches the original
	assert.Equal(t, originalConfig.Level, loadedConfig.Level)
	assert.Equal(t, originalConfig.ProfileName, loadedConfig.ProfileName)
	assert.Equal(t, originalConfig.Description, loadedConfig.Description)
	assert.Equal(t, originalConfig.Version, loadedConfig.Version)
}

func TestConfig_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		config      Config
		shouldPass  bool
		errorString string
	}{
		{
			name:       "valid default config",
			config:     DefaultConfig(),
			shouldPass: true,
		},
		{
			name: "invalid optimization level",
			config: Config{
				Level:   OptimizationLevel(999),
				Options: DefaultPipelineOptions(),
				Version: "1.0.0",
			},
			shouldPass:  false,
			errorString: "invalid optimization level",
		},
		{
			name: "unknown profile",
			config: Config{
				Level:       OptimizationLevelStandard,
				ProfileName: "nonexistent-profile",
				Options:     DefaultPipelineOptions(),
				Version:     "1.0.0",
			},
			shouldPass:  false,
			errorString: "unknown profile",
		},
		{
			name: "invalid pipeline options",
			config: Config{
				Level: OptimizationLevelStandard,
				Options: PipelineOptions{
					MaxFileSize:        -1, // Invalid
					MaxConcurrentFiles: 1,
				},
				Version: "1.0.0",
			},
			shouldPass:  false,
			errorString: "invalid pipeline options",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.Validate()
			if tc.shouldPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorString)
			}
		})
	}
}

func TestConfig_Normalize(t *testing.T) {
	config := Config{
		Level: OptimizationLevelSafe,
		// Missing version and other fields
	}

	config.Normalize()

	assert.Equal(t, "1.0.0", config.Version)
	assert.NotNil(t, config.CustomProfiles)
	
	// Should have options matching the level since no profile was specified
	expectedOptions := OptimizationLevelSafe.ToPipelineOptions()
	assert.Equal(t, expectedOptions.RemoveUnusedImports, config.Options.RemoveUnusedImports)
	assert.Equal(t, expectedOptions.OrganizeImports, config.Options.OrganizeImports)
}

func TestConfig_SetProfile(t *testing.T) {
	config := DefaultConfig()

	// Test setting a valid predefined profile
	err := config.SetProfile("performance")
	require.NoError(t, err)

	assert.Equal(t, "performance", config.ProfileName)
	assert.Equal(t, OptimizationLevelAggressive, config.Level)
	assert.True(t, config.Options.RemoveUnusedVars) // Performance profile should enable this

	// Test setting an invalid profile
	err = config.SetProfile("nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown profile")
}

func TestConfig_AddCustomProfile(t *testing.T) {
	config := DefaultConfig()

	// Create a custom profile
	customOptions := DefaultPipelineOptions()
	customOptions.RemoveUnusedImports = true
	customOptions.OrganizeImports = false
	customProfile := CustomProfile("my-custom", "My custom profile", customOptions)

	// Add the custom profile
	config.AddCustomProfile("my-custom", customProfile)

	// Verify it was added
	profiles := config.ListProfiles()
	profile, exists := profiles["my-custom"]
	assert.True(t, exists)
	assert.Equal(t, "my-custom", profile.Name)
	assert.Equal(t, "My custom profile", profile.Description)

	// Test setting the custom profile
	err := config.SetProfile("my-custom")
	require.NoError(t, err)
	assert.Equal(t, "my-custom", config.ProfileName)
}

func TestConfig_ListProfiles(t *testing.T) {
	config := DefaultConfig()

	// Add a custom profile
	customProfile := CustomProfile("test-profile", "Test profile", DefaultPipelineOptions())
	config.AddCustomProfile("test-profile", customProfile)

	profiles := config.ListProfiles()

	// Should contain predefined profiles
	expectedPredefined := []string{"conservative", "balanced", "performance", "maximum"}
	for _, name := range expectedPredefined {
		_, exists := profiles[name]
		assert.True(t, exists, "Should contain predefined profile: %s", name)
	}

	// Should contain custom profile
	_, exists := profiles["test-profile"]
	assert.True(t, exists, "Should contain custom profile")
}

func TestPipelineOptions_Validate(t *testing.T) {
	testCases := []struct {
		name        string
		options     PipelineOptions
		shouldPass  bool
		errorString string
	}{
		{
			name:       "valid default options",
			options:    DefaultPipelineOptions(),
			shouldPass: true,
		},
		{
			name: "invalid max file size",
			options: PipelineOptions{
				MaxFileSize:        -1,
				MaxConcurrentFiles: 1,
			},
			shouldPass:  false,
			errorString: "max file size must be positive",
		},
		{
			name: "invalid max concurrent files",
			options: PipelineOptions{
				MaxFileSize:        1024,
				MaxConcurrentFiles: 0,
			},
			shouldPass:  false,
			errorString: "max concurrent files must be positive",
		},
		{
			name: "empty include pattern",
			options: PipelineOptions{
				MaxFileSize:        1024,
				MaxConcurrentFiles: 1,
				IncludePatterns:    []string{""},
			},
			shouldPass:  false,
			errorString: "include patterns cannot be empty",
		},
		{
			name: "conflicting write and dry run",
			options: PipelineOptions{
				MaxFileSize:         1024,
				MaxConcurrentFiles:  1,
				WriteOptimizedFiles: true,
				DryRun:              true,
			},
			shouldPass:  false,
			errorString: "cannot write files in dry run mode",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.options.Validate()
			if tc.shouldPass {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.errorString)
			}
		})
	}
}

func TestConfigManager(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "config-manager-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.json")
	cm := NewConfigManager(configPath)

	// Test loading non-existent config (should create default)
	err = cm.Load()
	require.NoError(t, err)

	config := cm.GetConfig()
	assert.NotNil(t, config)
	assert.Equal(t, OptimizationLevelStandard, config.Level)

	// Test updating level
	err = cm.UpdateLevel(OptimizationLevelExpert)
	require.NoError(t, err)

	config = cm.GetConfig()
	assert.Equal(t, OptimizationLevelExpert, config.Level)
	assert.Empty(t, config.ProfileName) // Should clear profile when setting level directly

	// Test updating profile
	err = cm.UpdateProfile("conservative")
	require.NoError(t, err)

	config = cm.GetConfig()
	assert.Equal(t, "conservative", config.ProfileName)
	assert.Equal(t, OptimizationLevelSafe, config.Level)

	// Test saving and loading
	err = cm.Save()
	require.NoError(t, err)

	// Create a new manager and load the saved config
	cm2 := NewConfigManager(configPath)
	err = cm2.Load()
	require.NoError(t, err)

	config2 := cm2.GetConfig()
	assert.Equal(t, config.Level, config2.Level)
	assert.Equal(t, config.ProfileName, config2.ProfileName)
}

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()
	assert.NotEmpty(t, path)
	assert.Contains(t, path, ".go-starter-optimization.json")
}

func TestConfig_ConfigSummary(t *testing.T) {
	config := DefaultConfig()
	config.Level = OptimizationLevelAggressive
	config.ProfileName = "performance"

	summary := config.ConfigSummary()

	assert.Contains(t, summary, "aggressive")
	assert.Contains(t, summary, "performance")
	assert.Contains(t, summary, "Remove unused imports")
	assert.Contains(t, summary, "Organize imports")
	assert.Contains(t, summary, "Dry run mode")

	// Verify it contains actual values
	assert.Contains(t, summary, "true")
	assert.Contains(t, summary, "false")
}

func TestConfig_GetEffectiveOptions(t *testing.T) {
	// Test with profile-based options
	config := DefaultConfig()
	err := config.SetProfile("performance")
	require.NoError(t, err)

	options := config.GetEffectiveOptions()
	
	// Performance profile should enable aggressive optimizations
	assert.True(t, options.RemoveUnusedImports)
	assert.True(t, options.OrganizeImports)
	assert.True(t, options.RemoveUnusedVars)

	// Test with level-based options (no profile)
	config2 := Config{
		Level:   OptimizationLevelSafe,
		Options: OptimizationLevelSafe.ToPipelineOptions(),
		Version: "1.0.0",
	}

	options2 := config2.GetEffectiveOptions()
	assert.True(t, options2.RemoveUnusedImports)
	assert.True(t, options2.OrganizeImports)
	assert.False(t, options2.RemoveUnusedVars) // Safe level shouldn't enable this
}

func TestLoadConfig_InvalidFile(t *testing.T) {
	// Test loading from non-existent file
	_, err := LoadConfig("/nonexistent/path/config.json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")

	// Test loading invalid JSON
	tempDir, err := os.MkdirTemp("", "invalid-config-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	invalidConfigPath := filepath.Join(tempDir, "invalid.json")
	err = os.WriteFile(invalidConfigPath, []byte("invalid json"), 0644)
	require.NoError(t, err)

	_, err = LoadConfig(invalidConfigPath)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}