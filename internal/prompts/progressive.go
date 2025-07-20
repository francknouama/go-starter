package prompts

import (
	"fmt"
	"strings"

	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Type aliases for convenience
type ComplexityLevel = interfaces.ComplexityLevel
type DisclosureMode = interfaces.DisclosureMode

const (
	ComplexitySimple   = interfaces.ComplexitySimple
	ComplexityStandard = interfaces.ComplexityStandard
	ComplexityAdvanced = interfaces.ComplexityAdvanced
	ComplexityExpert   = interfaces.ComplexityExpert
	
	DisclosureModeBasic    = interfaces.DisclosureModeBasic
	DisclosureModeAdvanced = interfaces.DisclosureModeAdvanced
)

// ParseComplexityLevel parses a string into a ComplexityLevel
func ParseComplexityLevel(s string) (ComplexityLevel, error) {
	switch strings.ToLower(s) {
	case "simple":
		return ComplexitySimple, nil
	case "standard":
		return ComplexityStandard, nil
	case "advanced":
		return ComplexityAdvanced, nil
	case "expert":
		return ComplexityExpert, nil
	default:
		return ComplexitySimple, fmt.Errorf("invalid complexity level: %s (valid options: simple, standard, advanced, expert)", s)
	}
}


// DetermineDisclosureMode determines the disclosure mode based on flags
func DetermineDisclosureMode(basicFlag, advancedFlag bool, complexityFlag string) DisclosureMode {
	// Advanced flag takes precedence
	if advancedFlag {
		return DisclosureModeAdvanced
	}
	
	// If basic flag is explicitly set, use basic mode
	if basicFlag {
		return DisclosureModeBasic
	}
	
	// If complexity is specified, determine mode based on complexity level
	if complexityFlag != "" {
		if complexity, err := ParseComplexityLevel(complexityFlag); err == nil {
			switch complexity {
			case ComplexityAdvanced, ComplexityExpert:
				return DisclosureModeAdvanced
			default:
				return DisclosureModeBasic
			}
		}
	}
	
	// Default to basic mode for new users
	return DisclosureModeBasic
}

// Essential flags that should always be shown in basic mode
var essentialFlags = map[string]bool{
	"name":       true,
	"type":       true,
	"module":     true,
	"framework":  true,
	"logger":     true,
	"go-version": true,
	"output":     true,
	"help":       true,
	"quiet":      true,
	"dry-run":    true,
	"no-git":     true,
	"random-name": true,
	"no-banner":  true,
	"advanced":   true, // Always show the hint to advanced mode
	"basic":      true, // Always show basic mode flag
	"complexity": true, // Always show complexity flag
}

// FilterFlagsForDisclosureMode filters command flags based on disclosure mode
func FilterFlagsForDisclosureMode(cmd *cobra.Command, mode DisclosureMode) *cobra.Command {
	if mode == DisclosureModeAdvanced {
		// In advanced mode, show all flags
		return cmd
	}
	
	// In basic mode, create a filtered command with only essential flags
	filteredCmd := &cobra.Command{
		Use:     cmd.Use,
		Short:   cmd.Short,
		Long:    cmd.Long + "\n\n💡 Tip: Use --advanced to see all available options",
		Example: cmd.Example,
		Args:    cmd.Args,
		RunE:    cmd.RunE,
	}
	
	// Copy only essential flags
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if essentialFlags[flag.Name] {
			filteredCmd.Flags().AddFlag(flag)
		}
	})
	
	// Copy persistent flags
	cmd.PersistentFlags().VisitAll(func(flag *pflag.Flag) {
		filteredCmd.PersistentFlags().AddFlag(flag)
	})
	
	return filteredCmd
}

// SelectBlueprintForComplexity selects the appropriate blueprint based on type and complexity
func SelectBlueprintForComplexity(blueprintType string, complexity ComplexityLevel) string {
	// Only CLI blueprints have complexity variants for now
	if blueprintType == "cli" {
		switch complexity {
		case ComplexitySimple:
			return "cli-simple"
		default:
			return "cli" // Standard CLI uses just "cli" (architecture: standard)
		}
	}
	
	// For other blueprint types, return as-is
	return blueprintType
}

// Essential prompts that should be shown in basic mode
var essentialPrompts = []string{
	"project_name",
	"project_type",
	"module_path",
	"framework",
	"logger",
	"go_version",
}

// Advanced prompts that should only be shown in advanced mode
var advancedPrompts = []string{
	"database_driver",
	"database_orm",
	"auth_type",
	"deployment_targets",
	"testing_framework",
	"api_documentation",
	"monitoring",
	"caching",
	"message_queue",
}

// GetPromptsForDisclosureMode returns the list of prompts to show based on disclosure mode
func GetPromptsForDisclosureMode(mode DisclosureMode) []string {
	if mode == DisclosureModeAdvanced {
		// Return all prompts
		allPrompts := make([]string, 0, len(essentialPrompts)+len(advancedPrompts))
		allPrompts = append(allPrompts, essentialPrompts...)
		allPrompts = append(allPrompts, advancedPrompts...)
		return allPrompts
	}
	
	// Return only essential prompts for basic mode
	return essentialPrompts
}

// ValidateComplexityForBlueprint validates that the complexity level is appropriate for the blueprint
func ValidateComplexityForBlueprint(blueprintType string, complexity ComplexityLevel) error {
	// Currently, all blueprints support all complexity levels
	// This function is here for future validation logic if needed
	return nil
}

// GetRecommendedComplexity suggests a complexity level based on user experience indicators
func GetRecommendedComplexity(isFirstTime bool, hasGoExperience bool, projectSize string) ComplexityLevel {
	if isFirstTime || !hasGoExperience {
		return ComplexitySimple
	}
	
	switch projectSize {
	case "small", "prototype":
		return ComplexitySimple
	case "medium", "production":
		return ComplexityStandard
	case "large", "enterprise":
		return ComplexityAdvanced
	default:
		return ComplexityStandard
	}
}

// GetComplexityDescription returns a user-friendly description of the complexity level
func GetComplexityDescription(complexity ComplexityLevel) string {
	switch complexity {
	case ComplexitySimple:
		return "Simple - Minimal structure for quick prototypes and learning"
	case ComplexityStandard:
		return "Standard - Balanced structure for most production applications"
	case ComplexityAdvanced:
		return "Advanced - Comprehensive structure with enterprise patterns"
	case ComplexityExpert:
		return "Expert - Full-featured structure with all advanced options"
	default:
		return "Unknown complexity level"
	}
}