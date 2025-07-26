package optimization

import (
	"fmt"
	"os"
)

// DemonstratePipeline shows how to use the optimization pipeline
func DemonstratePipeline() {
	fmt.Println("=== Go Code Optimization Pipeline Demo ===")

	// 1. Show available optimization levels
	fmt.Println("1. Available Optimization Levels:")
	for level := OptimizationLevelNone; level <= OptimizationLevelExpert; level++ {
		fmt.Printf("   %s: %s\n", level.String(), level.Description())
	}
	fmt.Println()

	// 2. Show predefined profiles
	fmt.Println("2. Predefined Profiles:")
	profiles := PredefinedProfiles()
	for name, profile := range profiles {
		fmt.Printf("   %s (%s): %s\n", name, profile.Level.String(), profile.Description)
	}
	fmt.Println()

	// 3. Show recommended levels for different contexts
	fmt.Println("3. Recommended Levels by Use Case:")
	useCases := []string{"development", "testing", "production", "maintenance"}
	for _, useCase := range useCases {
		level := GetRecommendedLevel(useCase)
		fmt.Printf("   %s: %s\n", useCase, level.String())
	}
	fmt.Println()

	// 4. Demonstrate configuration management
	fmt.Println("4. Configuration Management Example:")
	
	// Create a configuration
	config := DefaultConfig()
	fmt.Printf("   Default config: Level=%s, Profile=%s\n", 
		config.Level.String(), config.ProfileName)

	// Change to aggressive optimization
	config.SetProfile("performance")
	fmt.Printf("   Performance profile: Level=%s\n", config.Level.String())
	
	// Show configuration summary
	fmt.Println("\n   Configuration Summary:")
	summary := config.ConfigSummary()
	for _, line := range []string{
		"   " + summary[:min(len(summary), 100)] + "...",
	} {
		fmt.Println(line)
	}

	fmt.Println()

	// 5. Demonstrate pipeline usage
	fmt.Println("5. Pipeline Usage Example:")
	
	// Create pipeline with safe optimizations
	safeOptions := OptimizationLevelSafe.ToPipelineOptions()
	safeOptions.DryRun = true
	safeOptions.Verbose = false
	
	_ = NewOptimizationPipeline(safeOptions) // Create pipeline (unused in demo)
	fmt.Printf("   Created pipeline with %s optimizations\n", OptimizationLevelSafe.String())
	fmt.Printf("   Dry run mode: %v\n", safeOptions.DryRun)
	fmt.Printf("   Remove unused imports: %v\n", safeOptions.RemoveUnusedImports)
	fmt.Printf("   Organize imports: %v\n", safeOptions.OrganizeImports)
	fmt.Printf("   Remove unused variables: %v\n", safeOptions.RemoveUnusedVars)

	fmt.Println("\n=== Demo Complete ===")
}

// ShowOptimizationHelp displays help information about optimization options
func ShowOptimizationHelp() {
	fmt.Println("Go Code Optimization Help")
	fmt.Println("========================")
	fmt.Println()

	fmt.Println("OPTIMIZATION LEVELS:")
	for level := OptimizationLevelNone; level <= OptimizationLevelExpert; level++ {
		fmt.Printf("  %s\n", level.String())
		fmt.Printf("    %s\n", level.Description())
		fmt.Println()
	}

	fmt.Println("PROFILES:")
	profiles := PredefinedProfiles()
	for name, profile := range profiles {
		fmt.Printf("  %s\n", name)
		fmt.Printf("    Level: %s\n", profile.Level.String())
		fmt.Printf("    %s\n", profile.Description)
		fmt.Println()
	}

	fmt.Println("USAGE EXAMPLES:")
	fmt.Println("  # Use safe optimization level")
	fmt.Println("  go-starter optimize --level=safe /path/to/project")
	fmt.Println()
	fmt.Println("  # Use performance profile")
	fmt.Println("  go-starter optimize --profile=performance /path/to/project")
	fmt.Println()
	fmt.Println("  # Dry run with verbose output")
	fmt.Println("  go-starter optimize --level=standard --dry-run --verbose /path/to/project")
	fmt.Println()
	fmt.Println("  # Custom optimization settings")
	fmt.Println("  go-starter optimize --remove-imports --organize-imports /path/to/project")
	fmt.Println()
}

// ValidateConfiguration performs comprehensive validation of a configuration
func ValidateConfiguration(config *Config) []string {
	var warnings []string

	// Check for potentially risky configurations
	options := config.GetEffectiveOptions()

	if options.RemoveUnusedVars && !options.DryRun {
		warnings = append(warnings, "Removing unused variables can break code - consider using dry-run mode first")
	}

	if options.RemoveUnusedFuncs && !options.DryRun {
		warnings = append(warnings, "Removing unused functions can break code - consider using dry-run mode first")
	}

	if config.Level >= OptimizationLevelAggressive && !options.CreateBackups && !options.DryRun {
		warnings = append(warnings, "Aggressive optimizations without backups can be dangerous - enable backups or dry-run")
	}

	if options.AddMissingImports && !options.DryRun {
		warnings = append(warnings, "Adding missing imports can introduce unintended dependencies - review changes carefully")
	}

	if options.MaxConcurrentFiles > 16 {
		warnings = append(warnings, "High concurrency settings may overwhelm system resources")
	}

	if options.MaxFileSize > 50*1024*1024 { // 50MB
		warnings = append(warnings, "Very large file size limits may cause memory issues")
	}

	return warnings
}

// min helper function since Go doesn't have a built-in min for int
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// ExampleUsage demonstrates common usage patterns
func ExampleUsage() {
	fmt.Println("=== Common Usage Patterns ===")

	// Example 1: Development workflow
	fmt.Println("1. Development Workflow (Safe Optimizations):")
	fmt.Println("   - Remove unused imports")
	fmt.Println("   - Organize imports alphabetically")
	fmt.Println("   - Use dry-run mode to preview changes")
	
	devConfig := DefaultConfig()
	devConfig.SetProfile("conservative")
	devConfig.Options.DryRun = true
	devConfig.Options.Verbose = true
	
	fmt.Printf("   Profile: %s (%s level)\n", devConfig.ProfileName, devConfig.Level.String())
	fmt.Println()

	// Example 2: CI/CD pipeline
	fmt.Println("2. CI/CD Pipeline (Standard Optimizations):")
	fmt.Println("   - All safe optimizations")
	fmt.Println("   - Carefully add missing imports")
	fmt.Println("   - Skip test files for faster processing")
	
	ciConfig := DefaultConfig()
	ciConfig.SetProfile("balanced")
	ciConfig.Options.SkipTestFiles = true
	ciConfig.Options.DryRun = false // Actually apply changes in CI
	
	fmt.Printf("   Profile: %s (%s level)\n", ciConfig.ProfileName, ciConfig.Level.String())
	fmt.Println()

	// Example 3: Production release preparation
	fmt.Println("3. Production Release (Aggressive Optimizations):")
	fmt.Println("   - All optimizations enabled")
	fmt.Println("   - Remove unused variables and functions")
	fmt.Println("   - Create backups for safety")
	
	prodConfig := DefaultConfig()
	prodConfig.SetProfile("performance")
	prodConfig.Options.CreateBackups = true
	prodConfig.Options.WriteOptimizedFiles = true
	
	fmt.Printf("   Profile: %s (%s level)\n", prodConfig.ProfileName, prodConfig.Level.String())
	fmt.Println()

	// Example 4: Code maintenance and cleanup
	fmt.Println("4. Code Maintenance (Expert Optimizations):")
	fmt.Println("   - Maximum optimizations")
	fmt.Println("   - Process test files too")
	fmt.Println("   - High performance settings")
	
	maintenanceConfig := DefaultConfig()
	maintenanceConfig.SetProfile("maximum")
	maintenanceConfig.Options.SkipTestFiles = false
	maintenanceConfig.Options.MaxConcurrentFiles = 16
	
	fmt.Printf("   Profile: %s (%s level)\n", maintenanceConfig.ProfileName, maintenanceConfig.Level.String())
	
	warnings := ValidateConfiguration(&maintenanceConfig)
	if len(warnings) > 0 {
		fmt.Println("   Warnings:")
		for _, warning := range warnings {
			fmt.Printf("   - %s\n", warning)
		}
	}
	fmt.Println()

	fmt.Println("=== End Examples ===")
}

// RunDemo runs a complete demonstration if called directly
func RunDemo() {
	if len(os.Args) > 1 && os.Args[1] == "demo" {
		DemonstratePipeline()
		fmt.Println()
		ExampleUsage()
		fmt.Println()
		ShowOptimizationHelp()
	}
}