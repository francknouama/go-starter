package quality

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestProgressiveComplexityATDD validates progressive complexity acceptance criteria
// This ensures the complexity tier system works correctly and provides appropriate progression
func TestProgressiveComplexityATDD(t *testing.T) {
	// Complexity tier definitions based on documentation
	complexityTiers := map[string]ComplexityTier{
		"simple": {
			Level:       "simple",
			Description: "Beginner-friendly, minimal structure",
			FileRange:   FileRange{Min: 7, Max: 11, Target: 8},
			Features:    []string{"basic structure", "minimal dependencies", "single-file patterns"},
			NotFeatures: []string{"complex config", "middleware", "advanced patterns"},
		},
		"standard": {
			Level:       "standard", 
			Description: "Balanced, production-ready",
			FileRange:   FileRange{Min: 20, Max: 35, Target: 25},
			Features:    []string{"full structure", "configuration", "production patterns"},
			NotFeatures: []string{"over-engineering", "unnecessary complexity"},
		},
		"advanced": {
			Level:       "advanced",
			Description: "Enterprise patterns, full features",
			FileRange:   FileRange{Min: 30, Max: 50, Target: 40},
			Features:    []string{"enterprise patterns", "advanced config", "monitoring"},
			NotFeatures: []string{"experimental features"},
		},
		"expert": {
			Level:       "expert",
			Description: "Full-featured, all options",
			FileRange:   FileRange{Min: 45, Max: 80, Target: 60},
			Features:    []string{"all patterns", "full configuration", "advanced monitoring"},
			NotFeatures: []string{},
		},
	}

	// Blueprint types that support complexity levels
	blueprintTypes := []string{"cli", "web-api", "library"}

	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once
	buildGoStarter(t, tmpDir, projectRoot)

	t.Run("complexity_tier_file_count_validation", func(t *testing.T) {
		// GIVEN: Different complexity levels
		// WHEN: Generating projects with each complexity level
		// THEN: File counts should match expected ranges for each tier
		
		for complexityLevel, tier := range complexityTiers {
			for _, blueprintType := range blueprintTypes {
				t.Run(fmt.Sprintf("%s_complexity_%s_blueprint", complexityLevel, blueprintType), func(t *testing.T) {
					projectName := fmt.Sprintf("tier-%s-%s", complexityLevel, blueprintType)
					
					generateComplexityProject(t, tmpDir, projectName, blueprintType, complexityLevel)
					projectDir := filepath.Join(tmpDir, projectName)
					
					fileCount := countProjectFiles(t, projectDir)
					
					// Validate file count is within expected range
					assert.GreaterOrEqual(t, fileCount, tier.FileRange.Min, 
						"File count should be ≥%d for %s complexity (was %d)", 
						tier.FileRange.Min, complexityLevel, fileCount)
					assert.LessOrEqual(t, fileCount, tier.FileRange.Max,
						"File count should be ≤%d for %s complexity (was %d)", 
						tier.FileRange.Max, complexityLevel, fileCount)
					
					t.Logf("✓ %s %s: %d files (target: %d±%d)", 
						complexityLevel, blueprintType, fileCount, 
						tier.FileRange.Target, tier.FileRange.Max-tier.FileRange.Target)
				})
			}
		}
	})

	t.Run("progressive_complexity_reduction_validation", func(t *testing.T) {
		// GIVEN: CLI blueprints with simple vs standard complexity
		// WHEN: Comparing file counts
		// THEN: Simple should have 60-75% fewer files than standard
		
		// Generate simple CLI
		generateComplexityProject(t, tmpDir, "reduction-simple", "cli", "simple")
		simpleDir := filepath.Join(tmpDir, "reduction-simple")
		simpleCount := countProjectFiles(t, simpleDir)
		
		// Generate standard CLI
		generateComplexityProject(t, tmpDir, "reduction-standard", "cli", "standard")
		standardDir := filepath.Join(tmpDir, "reduction-standard")
		standardCount := countProjectFiles(t, standardDir)
		
		// Calculate reduction percentage
		reductionPercentage := float64(standardCount-simpleCount) / float64(standardCount) * 100
		
		t.Logf("File count comparison: simple=%d, standard=%d, reduction=%.1f%%", 
			simpleCount, standardCount, reductionPercentage)
		
		// Validate reduction is in expected range (60-75%)
		assert.GreaterOrEqual(t, reductionPercentage, 60.0, 
			"Simple CLI should have 60%+ fewer files than standard (was %.1f%%)", reductionPercentage)
		assert.LessOrEqual(t, reductionPercentage, 80.0,
			"Reduction should be realistic, not over 80%% (was %.1f%%)", reductionPercentage)
		
		// Both should compile successfully
		validateProjectCompilation(t, simpleDir, "simple-cli-complexity")
		validateProjectCompilation(t, standardDir, "standard-cli-complexity")
		
		t.Logf("✓ Progressive complexity reduction validated: %.1f%% reduction", reductionPercentage)
	})

	t.Run("complexity_feature_inclusion_validation", func(t *testing.T) {
		// GIVEN: Different complexity levels with specific features
		// WHEN: Examining generated project structure
		// THEN: Features should match complexity level expectations
		
		t.Run("simple_complexity_excludes_advanced_features", func(t *testing.T) {
			generateComplexityProject(t, tmpDir, "features-simple", "cli", "simple")
			projectDir := filepath.Join(tmpDir, "features-simple")
			
			// Simple should NOT have these advanced features
			advancedFeatures := []string{
				"internal/middleware/",
				"internal/config/",
				"cmd/completion.go",
				"configs/",
				"docker-compose.yml",
				"Dockerfile",
				"scripts/",
			}
			
			for _, feature := range advancedFeatures {
				featurePath := filepath.Join(projectDir, feature)
				if strings.HasSuffix(feature, "/") {
					assert.NoDirExists(t, featurePath, "Simple complexity should not have: %s", feature)
				} else {
					assert.NoFileExists(t, featurePath, "Simple complexity should not have: %s", feature)
				}
			}
			
			// Simple should have these basic features
			basicFeatures := []string{
				"main.go",
				"cmd/root.go",
				"go.mod",
				"README.md",
			}
			
			for _, feature := range basicFeatures {
				featurePath := filepath.Join(projectDir, feature)
				assert.FileExists(t, featurePath, "Simple complexity should have: %s", feature)
			}
		})
		
		t.Run("standard_complexity_includes_production_features", func(t *testing.T) {
			generateComplexityProject(t, tmpDir, "features-standard", "cli", "standard")
			projectDir := filepath.Join(tmpDir, "features-standard")
			
			// Standard should have these production features
			productionFeatures := []string{
				"internal/",
				"cmd/",
				"Makefile",
				"go.mod",
				"README.md",
			}
			
			for _, feature := range productionFeatures {
				featurePath := filepath.Join(projectDir, feature)
				if strings.HasSuffix(feature, "/") {
					assert.DirExists(t, featurePath, "Standard complexity should have directory: %s", feature)
				} else {
					assert.FileExists(t, featurePath, "Standard complexity should have file: %s", feature)
				}
			}
		})
	})

	t.Run("complexity_compilation_validation", func(t *testing.T) {
		// GIVEN: Projects generated with different complexity levels
		// WHEN: Compiling each project
		// THEN: All complexity levels should compile successfully
		
		for complexityLevel := range complexityTiers {
			for _, blueprintType := range blueprintTypes {
				t.Run(fmt.Sprintf("compile_%s_%s", complexityLevel, blueprintType), func(t *testing.T) {
					projectName := fmt.Sprintf("compile-%s-%s", complexityLevel, blueprintType)
					
					generateComplexityProject(t, tmpDir, projectName, blueprintType, complexityLevel)
					projectDir := filepath.Join(tmpDir, projectName)
					
					// Validate compilation
					validateProjectCompilation(t, projectDir, fmt.Sprintf("%s-%s", complexityLevel, blueprintType))
				})
			}
		}
	})

	t.Run("progressive_disclosure_integration", func(t *testing.T) {
		// GIVEN: Progressive disclosure flags
		// WHEN: Using --complexity flag with CLI generation
		// THEN: Should automatically select appropriate blueprint and apply defaults
		
		t.Run("complexity_simple_selects_cli_simple_blueprint", func(t *testing.T) {
			// Generate with complexity=simple
			args := []string{"new", "disclosure-simple", 
				"--type=cli", 
				"--complexity=simple", 
				"--module=github.com/test/disclosure-simple",
				"--dry-run"}
			
			goStarterPath := filepath.Join(tmpDir, "go-starter")
			cmd := exec.Command(goStarterPath, args...)
			output, err := cmd.CombinedOutput()
			
			require.NoError(t, err, "Should generate with complexity=simple")
			outputStr := string(output)
			
			// Should indicate cli-simple blueprint usage
			assert.Contains(t, outputStr, "cli-simple", "Should use cli-simple blueprint for simple complexity")
		})
		
		t.Run("complexity_standard_selects_standard_cli_blueprint", func(t *testing.T) {
			// Generate with complexity=standard
			args := []string{"new", "disclosure-standard",
				"--type=cli",
				"--complexity=standard", 
				"--module=github.com/test/disclosure-standard",
				"--dry-run"}
			
			goStarterPath := filepath.Join(tmpDir, "go-starter")
			cmd := exec.Command(goStarterPath, args...)
			output, err := cmd.CombinedOutput()
			
			require.NoError(t, err, "Should generate with complexity=standard")
			outputStr := string(output)
			
			// Should use standard CLI blueprint (not cli-simple)
			assert.NotContains(t, outputStr, "cli-simple", "Should not use cli-simple for standard complexity")
		})
		
		t.Run("smart_defaults_applied_with_complexity", func(t *testing.T) {
			// When complexity is specified, smart defaults should prevent prompting
			projectName := "smart-defaults-test"
			
			generateComplexityProject(t, tmpDir, projectName, "cli", "simple")
			projectDir := filepath.Join(tmpDir, projectName)
			
			// Verify logger default was applied (should be slog)
			loggerFiles := findLoggerFiles(t, projectDir)
			require.NotEmpty(t, loggerFiles, "Should have logger files")
			
			for _, loggerFile := range loggerFiles {
				content, err := os.ReadFile(loggerFile)
				require.NoError(t, err)
				contentStr := string(content)
				
				// Should default to slog
				assert.Contains(t, contentStr, "log/slog", "Should default to slog logger")
			}
			
			// Verify go.mod has correct module path
			goModContent, err := os.ReadFile(filepath.Join(projectDir, "go.mod"))
			require.NoError(t, err)
			goModStr := string(goModContent)
			
			assert.Contains(t, goModStr, "github.com/test/smart-defaults-test", "Should have correct module path")
		})
	})

	t.Run("complexity_validation_error_handling", func(t *testing.T) {
		// GIVEN: Invalid complexity levels
		// WHEN: Attempting to generate projects
		// THEN: Should provide clear error messages with valid options
		
		invalidComplexityLevels := []string{"invalid", "beginner", "pro", "enterprise"}
		
		for _, invalidLevel := range invalidComplexityLevels {
			t.Run("invalid_complexity_"+invalidLevel, func(t *testing.T) {
				args := []string{"new", "invalid-test",
					"--type=cli",
					"--complexity=" + invalidLevel,
					"--dry-run"}
				
				goStarterPath := filepath.Join(tmpDir, "go-starter")
				cmd := exec.Command(goStarterPath, args...)
				output, err := cmd.CombinedOutput()
				
				assert.Error(t, err, "Invalid complexity level should be rejected")
				outputStr := string(output)
				assert.Contains(t, outputStr, "invalid complexity", "Should show validation error")
				assert.Contains(t, outputStr, "simple", "Should show valid complexity options")
				assert.Contains(t, outputStr, "standard", "Should show valid complexity options")
			})
		}
	})

	t.Run("complexity_tier_learning_progression", func(t *testing.T) {
		// GIVEN: A learning progression scenario
		// WHEN: User progresses from simple → standard → advanced
		// THEN: Each tier should build on the previous with clear upgrade paths
		
		// Generate all tiers for CLI
		tiers := []string{"simple", "standard", "advanced"}
		tierProjects := make(map[string]string)
		tierFileCounts := make(map[string]int)
		
		for _, tier := range tiers {
			projectName := fmt.Sprintf("progression-%s", tier)
			generateComplexityProject(t, tmpDir, projectName, "cli", tier)
			
			projectDir := filepath.Join(tmpDir, projectName)
			tierProjects[tier] = projectDir
			tierFileCounts[tier] = countProjectFiles(t, projectDir)
			
			// All tiers should compile
			validateProjectCompilation(t, projectDir, fmt.Sprintf("progression-%s", tier))
		}
		
		// Validate progression: each tier should have more files than the previous
		assert.Less(t, tierFileCounts["simple"], tierFileCounts["standard"], 
			"Standard should have more files than simple")
		assert.Less(t, tierFileCounts["standard"], tierFileCounts["advanced"], 
			"Advanced should have more files than standard")
		
		// Validate common core exists in all tiers
		coreFiles := []string{"main.go", "go.mod", "README.md"}
		for _, tier := range tiers {
			for _, coreFile := range coreFiles {
				filePath := filepath.Join(tierProjects[tier], coreFile)
				assert.FileExists(t, filePath, "Core file %s should exist in %s tier", coreFile, tier)
			}
		}
		
		t.Logf("✓ Learning progression validated: simple(%d) < standard(%d) < advanced(%d) files",
			tierFileCounts["simple"], tierFileCounts["standard"], tierFileCounts["advanced"])
	})
}

// ComplexityTier represents a complexity tier configuration
type ComplexityTier struct {
	Level       string
	Description string
	FileRange   FileRange
	Features    []string
	NotFeatures []string
}

// FileRange represents the expected file count range for a complexity tier
type FileRange struct {
	Min    int
	Max    int
	Target int
}

// generateComplexityProject generates a project with specific complexity level
func generateComplexityProject(t *testing.T, tmpDir, projectName, blueprintType, complexity string) {
	t.Helper()
	
	args := []string{"new", projectName}
	
	// Add blueprint type
	args = append(args, "--type="+blueprintType)
	
	// Add complexity level
	args = append(args, "--complexity="+complexity)
	
	// Add common flags
	args = append(args, 
		"--module=github.com/test/"+projectName,
		"--logger=slog",  // Use consistent logger for comparison
		"--no-git",
	)
	
	goStarterPath := filepath.Join(tmpDir, "go-starter")
	cmd := exec.Command(goStarterPath, args...)
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		t.Logf("Generate command: %s %s", goStarterPath, strings.Join(args, " "))
		t.Logf("Generate output: %s", string(output))
	}
	require.NoError(t, err, "Should generate project with %s complexity", complexity)
	
	t.Logf("Successfully generated %s project with %s complexity", blueprintType, complexity)
}