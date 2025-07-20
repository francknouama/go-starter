package prompts

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestComplexityLevel tests the ComplexityLevel type and its methods
func TestComplexityLevel(t *testing.T) {
	t.Run("valid_complexity_levels", func(t *testing.T) {
		tests := []struct {
			input    string
			expected ComplexityLevel
		}{
			{"simple", ComplexitySimple},
			{"standard", ComplexityStandard},
			{"advanced", ComplexityAdvanced},
			{"expert", ComplexityExpert},
		}

		for _, tt := range tests {
			t.Run(tt.input, func(t *testing.T) {
				level, err := ParseComplexityLevel(tt.input)
				require.NoError(t, err)
				assert.Equal(t, tt.expected, level)
			})
		}
	})

	t.Run("invalid_complexity_level", func(t *testing.T) {
		_, err := ParseComplexityLevel("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid complexity level")
	})

	t.Run("complexity_level_string", func(t *testing.T) {
		tests := []struct {
			level    ComplexityLevel
			expected string
		}{
			{ComplexitySimple, "simple"},
			{ComplexityStandard, "standard"},
			{ComplexityAdvanced, "advanced"},
			{ComplexityExpert, "expert"},
		}

		for _, tt := range tests {
			t.Run(tt.expected, func(t *testing.T) {
				assert.Equal(t, tt.expected, tt.level.String())
			})
		}
	})
}

// TestProgressiveDisclosureMode tests the progressive disclosure mode logic
func TestProgressiveDisclosureMode(t *testing.T) {
	t.Run("determine_mode_from_flags", func(t *testing.T) {
		tests := []struct {
			name           string
			basicFlag      bool
			advancedFlag   bool
			complexityFlag string
			expected       DisclosureMode
		}{
			{"default_mode", false, false, "", DisclosureModeBasic},
			{"basic_flag_set", true, false, "", DisclosureModeBasic},
			{"advanced_flag_set", false, true, "", DisclosureModeAdvanced},
			{"complexity_simple", false, false, "simple", DisclosureModeBasic},
			{"complexity_standard", false, false, "standard", DisclosureModeBasic},
			{"complexity_advanced", false, false, "advanced", DisclosureModeAdvanced},
			{"complexity_expert", false, false, "expert", DisclosureModeAdvanced},
			{"advanced_overrides_basic", true, true, "", DisclosureModeAdvanced},
			{"advanced_overrides_complexity", false, true, "simple", DisclosureModeAdvanced},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				mode := DetermineDisclosureMode(tt.basicFlag, tt.advancedFlag, tt.complexityFlag)
				assert.Equal(t, tt.expected, mode)
			})
		}
	})
}

// TestFlagFiltering tests the flag filtering based on disclosure mode
func TestFlagFiltering(t *testing.T) {
	// Create a test command with various flags
	cmd := &cobra.Command{
		Use: "test",
	}
	
	// Add essential flags
	cmd.Flags().String("name", "", "Project name")
	cmd.Flags().String("type", "", "Project type")
	cmd.Flags().String("module", "", "Module path")
	cmd.Flags().String("framework", "", "Framework")
	cmd.Flags().String("logger", "", "Logger")
	
	// Add advanced flags
	cmd.Flags().String("database-driver", "", "Database driver")
	cmd.Flags().String("database-orm", "", "Database ORM")
	cmd.Flags().String("auth-type", "", "Authentication type")
	cmd.Flags().String("banner-style", "", "Banner style")
	cmd.Flags().Bool("dry-run", false, "Dry run")

	t.Run("basic_mode_filters_advanced_flags", func(t *testing.T) {
		filtered := FilterFlagsForDisclosureMode(cmd, DisclosureModeBasic)
		
		// Should include essential flags
		assert.True(t, hasFlag(filtered, "name"))
		assert.True(t, hasFlag(filtered, "type"))
		assert.True(t, hasFlag(filtered, "module"))
		assert.True(t, hasFlag(filtered, "framework"))
		assert.True(t, hasFlag(filtered, "logger"))
		
		// Should exclude advanced flags
		assert.False(t, hasFlag(filtered, "database-driver"))
		assert.False(t, hasFlag(filtered, "database-orm"))
		assert.False(t, hasFlag(filtered, "auth-type"))
		assert.False(t, hasFlag(filtered, "banner-style"))
		
		// Should include utility flags
		assert.True(t, hasFlag(filtered, "dry-run"))
	})

	t.Run("advanced_mode_includes_all_flags", func(t *testing.T) {
		filtered := FilterFlagsForDisclosureMode(cmd, DisclosureModeAdvanced)
		
		// Should include all flags
		assert.True(t, hasFlag(filtered, "name"))
		assert.True(t, hasFlag(filtered, "type"))
		assert.True(t, hasFlag(filtered, "database-driver"))
		assert.True(t, hasFlag(filtered, "database-orm"))
		assert.True(t, hasFlag(filtered, "auth-type"))
		assert.True(t, hasFlag(filtered, "banner-style"))
		assert.True(t, hasFlag(filtered, "dry-run"))
	})
}

// TestBlueprintSelection tests blueprint selection based on complexity
func TestBlueprintSelection(t *testing.T) {
	t.Run("cli_blueprint_selection", func(t *testing.T) {
		tests := []struct {
			complexity ComplexityLevel
			expected   string
		}{
			{ComplexitySimple, "cli-simple"},
			{ComplexityStandard, "cli"},
			{ComplexityAdvanced, "cli"},
			{ComplexityExpert, "cli"},
		}

		for _, tt := range tests {
			t.Run(tt.complexity.String(), func(t *testing.T) {
				blueprint := SelectBlueprintForComplexity("cli", tt.complexity)
				assert.Equal(t, tt.expected, blueprint)
			})
		}
	})

	t.Run("other_blueprint_types_unchanged", func(t *testing.T) {
		// Non-CLI blueprint types should not be affected by complexity
		tests := []string{"web-api", "library", "lambda", "microservice"}
		
		for _, blueprintType := range tests {
			t.Run(blueprintType, func(t *testing.T) {
				blueprint := SelectBlueprintForComplexity(blueprintType, ComplexitySimple)
				assert.Equal(t, blueprintType, blueprint)
				
				blueprint = SelectBlueprintForComplexity(blueprintType, ComplexityExpert)
				assert.Equal(t, blueprintType, blueprint)
			})
		}
	})
}

// TestPromptFiltering tests prompt filtering based on disclosure mode
func TestPromptFiltering(t *testing.T) {
	t.Run("basic_mode_essential_prompts_only", func(t *testing.T) {
		prompts := GetPromptsForDisclosureMode(DisclosureModeBasic)
		
		// Should include essential prompts
		assert.Contains(t, prompts, "project_name")
		assert.Contains(t, prompts, "project_type")
		assert.Contains(t, prompts, "module_path")
		assert.Contains(t, prompts, "framework")
		assert.Contains(t, prompts, "logger")
		
		// Should exclude advanced prompts
		assert.NotContains(t, prompts, "database_driver")
		assert.NotContains(t, prompts, "database_orm")
		assert.NotContains(t, prompts, "auth_type")
		assert.NotContains(t, prompts, "deployment_targets")
	})

	t.Run("advanced_mode_all_prompts", func(t *testing.T) {
		prompts := GetPromptsForDisclosureMode(DisclosureModeAdvanced)
		
		// Should include all prompts
		assert.Contains(t, prompts, "project_name")
		assert.Contains(t, prompts, "project_type")
		assert.Contains(t, prompts, "database_driver")
		assert.Contains(t, prompts, "database_orm")
		assert.Contains(t, prompts, "auth_type")
		assert.Contains(t, prompts, "deployment_targets")
	})
}

// TestComplexityValidation tests complexity level validation
func TestComplexityValidation(t *testing.T) {
	t.Run("validate_complexity_for_blueprint", func(t *testing.T) {
		tests := []struct {
			blueprint  string
			complexity ComplexityLevel
			shouldErr  bool
		}{
			{"cli", ComplexitySimple, false},
			{"cli", ComplexityStandard, false},
			{"cli", ComplexityAdvanced, false},
			{"cli", ComplexityExpert, false},
			{"web-api", ComplexitySimple, false},
			{"web-api", ComplexityStandard, false},
			{"web-api", ComplexityAdvanced, false},
			{"web-api", ComplexityExpert, false},
			{"library", ComplexitySimple, false},
			{"library", ComplexityExpert, false},
		}

		for _, tt := range tests {
			t.Run(tt.blueprint+"_"+tt.complexity.String(), func(t *testing.T) {
				err := ValidateComplexityForBlueprint(tt.blueprint, tt.complexity)
				if tt.shouldErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// Helper function to check if a command has a specific flag
func hasFlag(cmd *cobra.Command, flagName string) bool {
	flag := cmd.Flags().Lookup(flagName)
	return flag != nil
}