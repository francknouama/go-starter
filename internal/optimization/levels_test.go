package optimization

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptimizationLevel_String(t *testing.T) {
	testCases := []struct {
		level    OptimizationLevel
		expected string
	}{
		{OptimizationLevelNone, "none"},
		{OptimizationLevelSafe, "safe"},
		{OptimizationLevelStandard, "standard"},
		{OptimizationLevelAggressive, "aggressive"},
		{OptimizationLevelExpert, "expert"},
		{OptimizationLevel(999), "unknown"},
	}

	for _, tc := range testCases {
		t.Run(tc.expected, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.level.String())
		})
	}
}

func TestParseOptimizationLevel(t *testing.T) {
	testCases := []struct {
		input       string
		expected    OptimizationLevel
		shouldParse bool
	}{
		{"none", OptimizationLevelNone, true},
		{"off", OptimizationLevelNone, true},
		{"0", OptimizationLevelNone, true},
		{"safe", OptimizationLevelSafe, true},
		{"basic", OptimizationLevelSafe, true},
		{"1", OptimizationLevelSafe, true},
		{"standard", OptimizationLevelStandard, true},
		{"normal", OptimizationLevelStandard, true},
		{"default", OptimizationLevelStandard, true},
		{"2", OptimizationLevelStandard, true},
		{"aggressive", OptimizationLevelAggressive, true},
		{"advanced", OptimizationLevelAggressive, true},
		{"3", OptimizationLevelAggressive, true},
		{"expert", OptimizationLevelExpert, true},
		{"maximum", OptimizationLevelExpert, true},
		{"max", OptimizationLevelExpert, true},
		{"4", OptimizationLevelExpert, true},
		{"invalid", OptimizationLevelNone, false},
		{"", OptimizationLevelNone, false},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			level, ok := ParseOptimizationLevel(tc.input)
			assert.Equal(t, tc.shouldParse, ok)
			if tc.shouldParse {
				assert.Equal(t, tc.expected, level)
			}
		})
	}
}

func TestOptimizationLevel_Description(t *testing.T) {
	levels := []OptimizationLevel{
		OptimizationLevelNone,
		OptimizationLevelSafe,
		OptimizationLevelStandard,
		OptimizationLevelAggressive,
		OptimizationLevelExpert,
	}

	for _, level := range levels {
		t.Run(level.String(), func(t *testing.T) {
			desc := level.Description()
			assert.NotEmpty(t, desc)
			assert.Contains(t, desc, "optimization")
		})
	}
}

func TestOptimizationLevel_ToPipelineOptions(t *testing.T) {
	testCases := []struct {
		level                   OptimizationLevel
		expectedRemoveImports   bool
		expectedOrganizeImports bool
		expectedAddMissing      bool
		expectedRemoveVars      bool
		expectedRemoveFuncs     bool
		expectedOptimizeConditionals bool
	}{
		{
			level:                   OptimizationLevelNone,
			expectedRemoveImports:   false,
			expectedOrganizeImports: false,
			expectedAddMissing:      false,
			expectedRemoveVars:      false,
			expectedRemoveFuncs:     false,
			expectedOptimizeConditionals: false,
		},
		{
			level:                   OptimizationLevelSafe,
			expectedRemoveImports:   true,
			expectedOrganizeImports: true,
			expectedAddMissing:      false,
			expectedRemoveVars:      false,
			expectedRemoveFuncs:     false,
			expectedOptimizeConditionals: false,
		},
		{
			level:                   OptimizationLevelStandard,
			expectedRemoveImports:   true,
			expectedOrganizeImports: true,
			expectedAddMissing:      true,
			expectedRemoveVars:      false,
			expectedRemoveFuncs:     false,
			expectedOptimizeConditionals: false,
		},
		{
			level:                   OptimizationLevelAggressive,
			expectedRemoveImports:   true,
			expectedOrganizeImports: true,
			expectedAddMissing:      true,
			expectedRemoveVars:      true,
			expectedRemoveFuncs:     true,
			expectedOptimizeConditionals: true,
		},
		{
			level:                   OptimizationLevelExpert,
			expectedRemoveImports:   true,
			expectedOrganizeImports: true,
			expectedAddMissing:      true,
			expectedRemoveVars:      true,
			expectedRemoveFuncs:     true,
			expectedOptimizeConditionals: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.level.String(), func(t *testing.T) {
			options := tc.level.ToPipelineOptions()

			assert.Equal(t, tc.expectedRemoveImports, options.RemoveUnusedImports)
			assert.Equal(t, tc.expectedOrganizeImports, options.OrganizeImports)
			assert.Equal(t, tc.expectedAddMissing, options.AddMissingImports)
			assert.Equal(t, tc.expectedRemoveVars, options.RemoveUnusedVars)
			assert.Equal(t, tc.expectedRemoveFuncs, options.RemoveUnusedFuncs)
			assert.Equal(t, tc.expectedOptimizeConditionals, options.OptimizeConditionals)

			// Verify performance settings scale with level
			if tc.level >= OptimizationLevelAggressive {
				assert.GreaterOrEqual(t, options.MaxConcurrentFiles, 8)
			}
			if tc.level == OptimizationLevelExpert {
				assert.GreaterOrEqual(t, options.MaxConcurrentFiles, 16)
				assert.False(t, options.SkipTestFiles) // Expert level processes test files
			}
		})
	}
}

func TestGetRecommendedLevel(t *testing.T) {
	testCases := []struct {
		useCase  string
		expected OptimizationLevel
	}{
		{"development", OptimizationLevelSafe},
		{"dev", OptimizationLevelSafe},
		{"local", OptimizationLevelSafe},
		{"testing", OptimizationLevelStandard},
		{"test", OptimizationLevelStandard},
		{"ci", OptimizationLevelStandard},
		{"production", OptimizationLevelAggressive},
		{"prod", OptimizationLevelAggressive},
		{"release", OptimizationLevelAggressive},
		{"maintenance", OptimizationLevelExpert},
		{"refactor", OptimizationLevelExpert},
		{"cleanup", OptimizationLevelExpert},
		{"unknown", OptimizationLevelStandard},
		{"", OptimizationLevelStandard},
	}

	for _, tc := range testCases {
		t.Run(tc.useCase, func(t *testing.T) {
			level := GetRecommendedLevel(tc.useCase)
			assert.Equal(t, tc.expected, level)
		})
	}
}

func TestPredefinedProfiles(t *testing.T) {
	profiles := PredefinedProfiles()

	// Check that all expected profiles exist
	expectedProfiles := []string{"conservative", "balanced", "performance", "maximum"}
	for _, profileName := range expectedProfiles {
		t.Run(profileName, func(t *testing.T) {
			profile, exists := profiles[profileName]
			assert.True(t, exists, "Profile %s should exist", profileName)
			assert.NotEmpty(t, profile.Name)
			assert.NotEmpty(t, profile.Description)
			assert.NotNil(t, profile.Options)

			// Verify the profile level matches expected behavior
			switch profileName {
			case "conservative":
				assert.Equal(t, OptimizationLevelSafe, profile.Level)
			case "balanced":
				assert.Equal(t, OptimizationLevelStandard, profile.Level)
			case "performance":
				assert.Equal(t, OptimizationLevelAggressive, profile.Level)
			case "maximum":
				assert.Equal(t, OptimizationLevelExpert, profile.Level)
			}
		})
	}
}

func TestCustomProfile(t *testing.T) {
	options := DefaultPipelineOptions()
	options.RemoveUnusedImports = true
	options.OrganizeImports = false
	options.RemoveUnusedVars = false

	profile := CustomProfile("test-profile", "Test profile for unit tests", options)

	assert.Equal(t, "test-profile", profile.Name)
	assert.Equal(t, "Test profile for unit tests", profile.Description)
	assert.Equal(t, options, profile.Options)

	// The level should be determined based on the options
	assert.Equal(t, OptimizationLevelSafe, profile.Level)
}

func TestDetermineOptimizationLevel(t *testing.T) {
	testCases := []struct {
		name     string
		options  PipelineOptions  
		expected OptimizationLevel
	}{
		{
			name:     "no optimizations",
			options:  PipelineOptions{},
			expected: OptimizationLevelNone,
		},
		{
			name: "safe optimizations only",
			options: PipelineOptions{
				RemoveUnusedImports: true,
				OrganizeImports:     true,
			},
			expected: OptimizationLevelSafe,
		},
		{
			name: "standard optimizations",
			options: PipelineOptions{
				RemoveUnusedImports: true,
				OrganizeImports:     true,
				AddMissingImports:   true,
			},
			expected: OptimizationLevelStandard,
		},
		{
			name: "aggressive optimizations",
			options: PipelineOptions{
				RemoveUnusedImports:   true,
				OrganizeImports:       true,
				AddMissingImports:     true,
				RemoveUnusedVars:      true,
				RemoveUnusedFuncs:     true,
			},
			expected: OptimizationLevelAggressive,
		},
		{
			name: "expert optimizations",
			options: PipelineOptions{
				RemoveUnusedImports:   true,
				OrganizeImports:       true,
				AddMissingImports:     true,
				RemoveUnusedVars:      true,
				RemoveUnusedFuncs:     true,
				OptimizeConditionals:  true,
			},
			expected: OptimizationLevelExpert,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			level := determineOptimizationLevel(tc.options)
			assert.Equal(t, tc.expected, level)
		})
	}
}

func TestValidateOptimizationLevel(t *testing.T) {
	testCases := []struct {
		level     OptimizationLevel
		context   string
		shouldPass bool
		message   string
	}{
		{OptimizationLevelSafe, "development", true, ""},
		{OptimizationLevelStandard, "development", true, ""},
		{OptimizationLevelAggressive, "development", false, "not recommended for development"},
		{OptimizationLevelExpert, "development", false, "not recommended for development"},
		
		{OptimizationLevelSafe, "ci", true, ""},
		{OptimizationLevelStandard, "ci", true, ""},
		{OptimizationLevelAggressive, "ci", true, ""},
		{OptimizationLevelExpert, "ci", false, "may be too aggressive for CI"},
		
		{OptimizationLevelNone, "production", true, ""},
		{OptimizationLevelSafe, "production", true, ""},
		{OptimizationLevelStandard, "production", true, ""},
		{OptimizationLevelAggressive, "production", true, ""},
		{OptimizationLevelExpert, "production", true, ""},
		
		{OptimizationLevelExpert, "maintenance", true, ""},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s_%s", tc.level.String(), tc.context), func(t *testing.T) {
			valid, message := ValidateOptimizationLevel(tc.level, tc.context)
			assert.Equal(t, tc.shouldPass, valid)
			if !tc.shouldPass {
				assert.Contains(t, message, tc.message)
			}
		})
	}
}