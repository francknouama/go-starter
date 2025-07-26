package optimization

// OptimizationLevel defines different levels of code optimization
type OptimizationLevel int

const (
	// OptimizationLevelNone performs no optimizations
	OptimizationLevelNone OptimizationLevel = iota
	
	// OptimizationLevelSafe performs only safe, non-destructive optimizations
	// - Remove unused imports
	// - Organize imports alphabetically
	OptimizationLevelSafe
	
	// OptimizationLevelStandard performs common optimizations that are generally safe
	// - All safe optimizations
	// - Add missing imports (with caution)
	// - Basic code organization
	OptimizationLevelStandard
	
	// OptimizationLevelAggressive performs more comprehensive optimizations
	// - All standard optimizations  
	// - Remove unused variables (local scope only)
	// - Remove unused private functions
	// - Optimize conditional expressions
	OptimizationLevelAggressive
	
	// OptimizationLevelExpert performs all available optimizations
	// - All aggressive optimizations
	// - Advanced AST transformations
	// - Maximum performance settings
	OptimizationLevelExpert
)

// String returns a human-readable name for the optimization level
func (level OptimizationLevel) String() string {
	switch level {
	case OptimizationLevelNone:
		return "none"
	case OptimizationLevelSafe:
		return "safe"
	case OptimizationLevelStandard:
		return "standard"
	case OptimizationLevelAggressive:
		return "aggressive"
	case OptimizationLevelExpert:
		return "expert"
	default:
		return "unknown"
	}
}

// ParseOptimizationLevel parses a string into an OptimizationLevel
func ParseOptimizationLevel(s string) (OptimizationLevel, bool) {
	switch s {
	case "none", "off", "0":
		return OptimizationLevelNone, true
	case "safe", "basic", "1":
		return OptimizationLevelSafe, true
	case "standard", "normal", "default", "2":
		return OptimizationLevelStandard, true
	case "aggressive", "advanced", "3":
		return OptimizationLevelAggressive, true
	case "expert", "maximum", "max", "4":
		return OptimizationLevelExpert, true
	default:
		return OptimizationLevelNone, false
	}
}

// Description returns a detailed description of what the optimization level includes
func (level OptimizationLevel) Description() string {
	switch level {
	case OptimizationLevelNone:
		return "No optimizations performed. Files are left unchanged."
	case OptimizationLevelSafe:
		return "Safe optimizations only: remove unused imports, organize imports alphabetically. No code logic changes."
	case OptimizationLevelStandard:
		return "Standard optimizations: all safe optimizations plus careful addition of missing imports and basic code organization."
	case OptimizationLevelAggressive:
		return "Aggressive optimizations: all standard optimizations plus removal of unused variables and private functions, conditional optimization."
	case OptimizationLevelExpert:
		return "Expert optimizations: all available optimizations with maximum performance settings. Use with caution."
	default:
		return "Unknown optimization level."
	}
}

// ToPipelineOptions converts an OptimizationLevel to PipelineOptions
func (level OptimizationLevel) ToPipelineOptions() PipelineOptions {
	base := DefaultPipelineOptions()
	
	switch level {
	case OptimizationLevelNone:
		// Disable all optimizations
		base.RemoveUnusedImports = false
		base.OrganizeImports = false
		base.AddMissingImports = false
		base.RemoveUnusedVars = false
		base.RemoveUnusedFuncs = false
		base.OptimizeConditionals = false
		
	case OptimizationLevelSafe:
		// Only safe import optimizations
		base.RemoveUnusedImports = true
		base.OrganizeImports = true
		base.AddMissingImports = false  // Can be risky
		base.RemoveUnusedVars = false   // Can break code
		base.RemoveUnusedFuncs = false  // Can break code
		base.OptimizeConditionals = false
		
	case OptimizationLevelStandard:
		// Common optimizations with some import additions
		base.RemoveUnusedImports = true
		base.OrganizeImports = true
		base.AddMissingImports = true   // Carefully add missing imports
		base.RemoveUnusedVars = false   // Still conservative
		base.RemoveUnusedFuncs = false  // Still conservative
		base.OptimizeConditionals = false
		
	case OptimizationLevelAggressive:
		// More comprehensive optimizations
		base.RemoveUnusedImports = true
		base.OrganizeImports = true
		base.AddMissingImports = true
		base.RemoveUnusedVars = true    // Remove unused local variables
		base.RemoveUnusedFuncs = true   // Remove unused private functions
		base.OptimizeConditionals = true
		base.MaxConcurrentFiles = 8     // Higher performance
		
	case OptimizationLevelExpert:
		// All optimizations enabled
		base.RemoveUnusedImports = true
		base.OrganizeImports = true
		base.AddMissingImports = true
		base.RemoveUnusedVars = true
		base.RemoveUnusedFuncs = true
		base.OptimizeConditionals = true
		base.MaxConcurrentFiles = 16    // Maximum performance
		base.MaxFileSize = 10 * 1024 * 1024  // 10MB limit
		base.SkipTestFiles = false      // Process test files too
		
	default:
		// Use safe defaults for unknown levels
		return OptimizationLevelSafe.ToPipelineOptions()
	}
	
	return base
}

// GetRecommendedLevel returns the recommended optimization level based on use case
func GetRecommendedLevel(useCase string) OptimizationLevel {
	switch useCase {
	case "development", "dev", "local":
		return OptimizationLevelSafe
	case "testing", "test", "ci":
		return OptimizationLevelStandard
	case "production", "prod", "release":
		return OptimizationLevelAggressive
	case "maintenance", "refactor", "cleanup":
		return OptimizationLevelExpert
	default:
		return OptimizationLevelStandard
	}
}

// OptimizationProfile represents a named configuration profile
type OptimizationProfile struct {
	Name        string
	Level       OptimizationLevel
	Description string
	Options     PipelineOptions
}

// PredefinedProfiles returns a set of predefined optimization profiles
func PredefinedProfiles() map[string]OptimizationProfile {
	return map[string]OptimizationProfile{
		"conservative": {
			Name:        "Conservative",
			Level:       OptimizationLevelSafe,
			Description: "Safe optimizations for development and experimentation",
			Options:     OptimizationLevelSafe.ToPipelineOptions(),
		},
		"balanced": {
			Name:        "Balanced",
			Level:       OptimizationLevelStandard,
			Description: "Balanced optimizations for CI/CD and testing",
			Options:     OptimizationLevelStandard.ToPipelineOptions(),
		},
		"performance": {
			Name:        "Performance",
			Level:       OptimizationLevelAggressive,
			Description: "Aggressive optimizations for production code",
			Options:     OptimizationLevelAggressive.ToPipelineOptions(),
		},
		"maximum": {
			Name:        "Maximum",
			Level:       OptimizationLevelExpert,
			Description: "All optimizations for code cleanup and refactoring",
			Options:     OptimizationLevelExpert.ToPipelineOptions(),
		},
	}
}

// CustomProfile creates a custom optimization profile
func CustomProfile(name, description string, options PipelineOptions) OptimizationProfile {
	// Determine the closest matching level based on options
	level := determineOptimizationLevel(options)
	
	return OptimizationProfile{
		Name:        name,
		Level:       level,
		Description: description,
		Options:     options,
	}
}

// determineOptimizationLevel attempts to determine the optimization level from options
func determineOptimizationLevel(options PipelineOptions) OptimizationLevel {
	// Count enabled optimizations to determine level
	optimizationsEnabled := 0
	
	if options.RemoveUnusedImports {
		optimizationsEnabled++
	}
	if options.OrganizeImports {
		optimizationsEnabled++
	}
	if options.AddMissingImports {
		optimizationsEnabled++
	}
	if options.RemoveUnusedVars {
		optimizationsEnabled++
	}
	if options.RemoveUnusedFuncs {
		optimizationsEnabled++
	}
	if options.OptimizeConditionals {
		optimizationsEnabled++
	}
	
	// Determine level based on optimization count and aggressiveness
	switch {
	case optimizationsEnabled == 0:
		return OptimizationLevelNone
	case optimizationsEnabled <= 2 && !options.RemoveUnusedVars && !options.RemoveUnusedFuncs:
		return OptimizationLevelSafe
	case optimizationsEnabled <= 3 && !options.RemoveUnusedVars && !options.RemoveUnusedFuncs:
		return OptimizationLevelStandard
	case optimizationsEnabled <= 5:
		return OptimizationLevelAggressive
	default:
		return OptimizationLevelExpert
	}
}

// ValidateOptimizationLevel checks if the optimization level is appropriate for the context
func ValidateOptimizationLevel(level OptimizationLevel, context string) (bool, string) {
	switch context {
	case "development":
		if level > OptimizationLevelStandard {
			return false, "aggressive optimizations not recommended for development - use 'safe' or 'standard'"
		}
	case "ci", "testing":
		if level > OptimizationLevelAggressive {
			return false, "expert optimizations may be too aggressive for CI/testing - consider 'standard' or 'aggressive'"
		}
	case "production":
		// All levels are acceptable for production
		return true, ""
	case "maintenance":
		// All levels are acceptable for maintenance
		return true, ""
	}
	
	return true, ""
}