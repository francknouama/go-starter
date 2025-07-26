package optimization

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFullIntegration_ConfigurationAndPipeline(t *testing.T) {
	// Create a temporary project and config directory
	tempDir, err := os.MkdirTemp("", "full-integration-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	projectDir := filepath.Join(tempDir, "project")
	configDir := filepath.Join(tempDir, "config")
	
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err)

	// Create a test project with optimization opportunities
	testFiles := map[string]string{
		"main.go": `package main

import (
	"fmt"
	"os"        // unused
	"strings"   // unused  
	"github.com/gin-gonic/gin"
	"encoding/json" // unused
)

func main() {
	fmt.Println("Hello World")
	r := gin.Default()
	r.Run()
}`,

		"helper.go": `package main

import (
	"fmt"
	"log"    // unused
	"time"   // unused
)

func helper() {
	fmt.Println("Helper function")
}`,
	}

	// Write test files
	for filename, content := range testFiles {
		filePath := filepath.Join(projectDir, filename)
		err = os.WriteFile(filePath, []byte(content), 0644)
		require.NoError(t, err)
	}

	// Test different optimization levels
	testScenarios := []struct {
		name           string
		level          OptimizationLevel
		profile        string
		expectedRemoved int
		shouldCreateBackups bool
	}{
		{
			name:           "safe optimization",
			level:          OptimizationLevelSafe,
			profile:        "conservative",
			expectedRemoved: 5, // All unused imports should be removed
			shouldCreateBackups: true,
		},
		{
			name:           "standard optimization", 
			level:          OptimizationLevelStandard,
			profile:        "balanced",
			expectedRemoved: 5,
			shouldCreateBackups: true,
		},
		{
			name:           "aggressive optimization",
			level:          OptimizationLevelAggressive,
			profile:        "performance", 
			expectedRemoved: 5,
			shouldCreateBackups: true,
		},
	}

	for _, scenario := range testScenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Create configuration for this scenario
			configPath := filepath.Join(configDir, scenario.name+".json")
			cm := NewConfigManager(configPath)

			// Load or create default config
			err := cm.Load()
			require.NoError(t, err)

			// Set the desired profile
			err = cm.UpdateProfile(scenario.profile)
			require.NoError(t, err)

			// Configure for testing
			config := cm.GetConfig()
			config.Options.DryRun = true // Don't actually modify files
			config.Options.Verbose = false // Reduce test output
			config.Options.CreateBackups = scenario.shouldCreateBackups

			// Save the configuration
			err = cm.Save()
			require.NoError(t, err)

			// Create pipeline from configuration
			pipeline := NewOptimizationPipeline(config.GetEffectiveOptions())

			// Run optimization
			result, err := pipeline.OptimizeProject(projectDir)
			require.NoError(t, err)

			// Verify results
			assert.Equal(t, 2, result.TotalFiles, "Should find 2 Go files")
			assert.Equal(t, 2, result.FilesProcessed, "Should process 2 files")
			assert.Equal(t, 2, result.FilesOptimized, "Both files should be optimized")
			assert.Equal(t, 0, result.FilesWithErrors, "Should have no errors")
			assert.Equal(t, scenario.expectedRemoved, result.ImportsRemoved, "Should remove expected number of imports")
			assert.Equal(t, 0, result.ImportsAdded, "Should not add any imports")

			// Verify configuration was saved and can be reloaded
			cm2 := NewConfigManager(configPath)
			err = cm2.Load()
			require.NoError(t, err)

			config2 := cm2.GetConfig()
			assert.Equal(t, config.Level, config2.Level)
			assert.Equal(t, config.ProfileName, config2.ProfileName)

			t.Logf("Scenario %s: Level=%s, Profile=%s, Removed=%d imports", 
				scenario.name, config.Level.String(), config.ProfileName, result.ImportsRemoved)
		})
	}
}

func TestConfigurationLevels_BehaviorVerification(t *testing.T) {
	// Verify that different optimization levels behave as expected
	testCases := []struct {
		level              OptimizationLevel
		expectRemoveImports bool
		expectOrganizeImports bool
		expectAddMissing   bool
		expectRemoveVars   bool
		expectRemoveFuncs  bool
	}{
		{
			level:              OptimizationLevelNone,
			expectRemoveImports: false,
			expectOrganizeImports: false,
			expectAddMissing:   false,
			expectRemoveVars:   false,
			expectRemoveFuncs:  false,
		},
		{
			level:              OptimizationLevelSafe,
			expectRemoveImports: true,
			expectOrganizeImports: true,
			expectAddMissing:   false,
			expectRemoveVars:   false,
			expectRemoveFuncs:  false,
		},
		{
			level:              OptimizationLevelStandard,
			expectRemoveImports: true,
			expectOrganizeImports: true,
			expectAddMissing:   true,
			expectRemoveVars:   false,
			expectRemoveFuncs:  false,
		},
		{
			level:              OptimizationLevelAggressive,
			expectRemoveImports: true,
			expectOrganizeImports: true,
			expectAddMissing:   true,
			expectRemoveVars:   true,
			expectRemoveFuncs:  true,
		},
		{
			level:              OptimizationLevelExpert,
			expectRemoveImports: true,
			expectOrganizeImports: true,
			expectAddMissing:   true,
			expectRemoveVars:   true,
			expectRemoveFuncs:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			options := tc.level.ToPipelineOptions()

			assert.Equal(t, tc.expectRemoveImports, options.RemoveUnusedImports)
			assert.Equal(t, tc.expectOrganizeImports, options.OrganizeImports)
			assert.Equal(t, tc.expectAddMissing, options.AddMissingImports)
			assert.Equal(t, tc.expectRemoveVars, options.RemoveUnusedVars)
			assert.Equal(t, tc.expectRemoveFuncs, options.RemoveUnusedFuncs)

			// Verify performance settings scale appropriately
			if tc.level >= OptimizationLevelAggressive {
				assert.GreaterOrEqual(t, options.MaxConcurrentFiles, 8)
			}

			if tc.level == OptimizationLevelExpert {
				assert.GreaterOrEqual(t, options.MaxConcurrentFiles, 16)
				assert.GreaterOrEqual(t, options.MaxFileSize, int64(10*1024*1024)) // 10MB
			}
		})
	}
}

func TestProfileConsistency(t *testing.T) {
	// Verify that profiles are consistent with their stated optimization levels
	profiles := PredefinedProfiles()

	for name, profile := range profiles {
		t.Run(name, func(t *testing.T) {
			// The profile's options should match what the level would produce
			expectedOptions := profile.Level.ToPipelineOptions()

			assert.Equal(t, expectedOptions.RemoveUnusedImports, profile.Options.RemoveUnusedImports,
				"Profile %s should match level %s for RemoveUnusedImports", name, profile.Level.String())
			assert.Equal(t, expectedOptions.OrganizeImports, profile.Options.OrganizeImports,
				"Profile %s should match level %s for OrganizeImports", name, profile.Level.String())
			assert.Equal(t, expectedOptions.AddMissingImports, profile.Options.AddMissingImports,
				"Profile %s should match level %s for AddMissingImports", name, profile.Level.String())
			assert.Equal(t, expectedOptions.RemoveUnusedVars, profile.Options.RemoveUnusedVars,
				"Profile %s should match level %s for RemoveUnusedVars", name, profile.Level.String())
			assert.Equal(t, expectedOptions.RemoveUnusedFuncs, profile.Options.RemoveUnusedFuncs,
				"Profile %s should match level %s for RemoveUnusedFuncs", name, profile.Level.String())
		})
	}
}

func TestValidationIntegration(t *testing.T) {
	// Test the complete validation workflow
	testCases := []struct {
		name        string
		level       OptimizationLevel
		context     string
		shouldPass  bool
		expectWarnings bool
	}{
		{
			name:        "safe development",
			level:       OptimizationLevelSafe,
			context:     "development",
			shouldPass:  true,
			expectWarnings: false,
		},
		{
			name:        "aggressive development",
			level:       OptimizationLevelAggressive,
			context:     "development",
			shouldPass:  false,
			expectWarnings: true,
		},
		{
			name:        "expert production",
			level:       OptimizationLevelExpert,
			context:     "production",
			shouldPass:  true,
			expectWarnings: true, // Still has warnings about risk
		},
		{
			name:        "standard ci",
			level:       OptimizationLevelStandard,
			context:     "ci",
			shouldPass:  true,
			expectWarnings: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test level validation
			valid, message := ValidateOptimizationLevel(tc.level, tc.context)
			assert.Equal(t, tc.shouldPass, valid)
			if !tc.shouldPass {
				assert.NotEmpty(t, message)
			}

			// Test configuration warnings
			config := DefaultConfig()
			config.Level = tc.level
			config.Options = tc.level.ToPipelineOptions()
			config.Options.DryRun = false // Enable warnings about non-dry-run

			warnings := ValidateConfiguration(&config)
			if tc.expectWarnings {
				assert.NotEmpty(t, warnings, "Expected warnings for %s level", tc.level.String())
			}

			// Log results for manual verification
			t.Logf("Level: %s, Context: %s, Valid: %v, Warnings: %d", 
				tc.level.String(), tc.context, valid, len(warnings))
			for _, warning := range warnings {
				t.Logf("  Warning: %s", warning)
			}
		})
	}
}

func TestRecommendationSystem(t *testing.T) {
	// Test that the recommendation system provides sensible defaults
	useCases := map[string]OptimizationLevel{
		"development": OptimizationLevelSafe,
		"testing":     OptimizationLevelStandard,
		"production":  OptimizationLevelAggressive,
		"maintenance": OptimizationLevelExpert,
	}

	for useCase, expectedLevel := range useCases {
		t.Run(useCase, func(t *testing.T) {
			recommendedLevel := GetRecommendedLevel(useCase)
			assert.Equal(t, expectedLevel, recommendedLevel)

			// Verify the recommended level is appropriate for the context
			valid, message := ValidateOptimizationLevel(recommendedLevel, useCase)
			assert.True(t, valid, "Recommended level should be valid for context: %s", message)

			t.Logf("Use case: %s -> Recommended level: %s", useCase, recommendedLevel.String())
		})
	}
}

func TestEndToEndWorkflow(t *testing.T) {
	// Simulate a complete workflow from configuration to optimization
	tempDir, err := os.MkdirTemp("", "e2e-workflow-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Step 1: Create a configuration manager
	configPath := filepath.Join(tempDir, "optimization.json")
	cm := NewConfigManager(configPath)

	// Step 2: Load default configuration
	err = cm.Load()
	require.NoError(t, err)

	// Step 3: Customize configuration for development workflow
	err = cm.UpdateProfile("conservative")
	require.NoError(t, err)

	config := cm.GetConfig()
	config.Options.DryRun = true
	config.Options.Verbose = false
	config.Options.CreateBackups = true

	// Step 4: Save configuration
	err = cm.Save()
	require.NoError(t, err)

	// Step 5: Create a test project
	projectDir := filepath.Join(tempDir, "test-project")
	err = os.MkdirAll(projectDir, 0755)
	require.NoError(t, err)

	testCode := `package main

import (
	"fmt"
	"os"      // unused
	"strings" // unused
)

func main() {
	fmt.Println("Hello World")
}`

	testFile := filepath.Join(projectDir, "main.go")
	err = os.WriteFile(testFile, []byte(testCode), 0644)
	require.NoError(t, err)

	// Step 6: Create pipeline from configuration
	pipeline := NewOptimizationPipeline(config.GetEffectiveOptions())

	// Step 7: Run optimization
	result, err := pipeline.OptimizeProject(projectDir)
	require.NoError(t, err)

	// Step 8: Verify results
	assert.Equal(t, 1, result.TotalFiles)
	assert.Equal(t, 1, result.FilesProcessed)
	assert.Equal(t, 1, result.FilesOptimized)
	assert.Equal(t, 2, result.ImportsRemoved) // "os" and "strings"
	assert.Equal(t, 0, result.ImportsAdded)
	assert.Equal(t, 0, result.FilesWithErrors)

	// Step 9: Verify configuration persistence
	cm2 := NewConfigManager(configPath)
	err = cm2.Load()
	require.NoError(t, err)

	config2 := cm2.GetConfig()
	assert.Equal(t, "conservative", config2.ProfileName)
	assert.Equal(t, OptimizationLevelSafe, config2.Level)

	t.Logf("End-to-end workflow completed successfully!")
	t.Logf("Configuration: Profile=%s, Level=%s", config2.ProfileName, config2.Level.String())
	t.Logf("Results: Processed=%d files, Removed=%d imports", 
		result.FilesProcessed, result.ImportsRemoved)
}