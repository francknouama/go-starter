package quality

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestComprehensiveATDDRunner runs all comprehensive ATDD tests and provides summary
// This is the main entry point for validating blueprint quality acceptance criteria
func TestComprehensiveATDDRunner(t *testing.T) {
	startTime := time.Now()
	
	// Ensure we have a clean environment
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	projectRoot := filepath.Join(originalDir, "..", "..", "..")
	
	defer func() { _ = os.Chdir(originalDir) }()
	_ = os.Chdir(tmpDir)

	// Build go-starter once for all tests
	buildGoStarter(t, tmpDir, projectRoot)
	
	// Test suite configuration
	testSuites := []ATDDTestSuite{
		{
			Name:        "Blueprint Generation Quality",
			Description: "Validates that all high-priority blueprints generate correctly and compile successfully",
			TestFunc:    runBlueprintGenerationQualityTests,
			Priority:    "CRITICAL",
		},
		{
			Name:        "Logger Integration",
			Description: "Validates simplified logger system works across all blueprints with 60-90% code reduction",
			TestFunc:    runLoggerIntegrationTests,
			Priority:    "HIGH",
		},
		{
			Name:        "Progressive Complexity",
			Description: "Validates complexity tiers work correctly with 60-75% file reduction from simple to standard",
			TestFunc:    runProgressiveComplexityTests,
			Priority:    "HIGH",
		},
		{
			Name:        "Cross-Blueprint Validation",
			Description: "Validates consistency and quality across all high-priority blueprints",
			TestFunc:    runCrossBlueprintValidationTests,
			Priority:    "HIGH",
		},
		{
			Name:        "Template Variable Resolution",
			Description: "Validates all template variables resolve correctly with no unresolved {{.Variable}} remnants",
			TestFunc:    runTemplateVariableResolutionTests,
			Priority:    "CRITICAL",
		},
	}
	
	results := make([]ATDDTestResult, 0, len(testSuites))
	
	t.Run("comprehensive_atdd_validation", func(t *testing.T) {
		for _, suite := range testSuites {
			t.Run(suite.Name, func(t *testing.T) {
				result := ATDDTestResult{
					SuiteName:   suite.Name,
					Description: suite.Description,
					Priority:    suite.Priority,
					StartTime:   time.Now(),
				}
				
				// Run the test suite
				passed := t.Run("execute_"+suite.Name, func(t *testing.T) {
					suite.TestFunc(t, tmpDir, projectRoot)
				})
				
				result.EndTime = time.Now()
				result.Duration = result.EndTime.Sub(result.StartTime)
				result.Passed = passed
				
				results = append(results, result)
				
				if passed {
					t.Logf("✅ %s: PASSED (%v)", suite.Name, result.Duration)
				} else {
					t.Errorf("❌ %s: FAILED (%v)", suite.Name, result.Duration)
				}
			})
		}
	})
	
	// Generate comprehensive summary
	generateATDDSummary(t, results, time.Since(startTime))
}

// ATDDTestSuite represents a test suite configuration
type ATDDTestSuite struct {
	Name        string
	Description string
	TestFunc    func(t *testing.T, tmpDir, projectRoot string)
	Priority    string
}

// ATDDTestResult represents the result of running a test suite
type ATDDTestResult struct {
	SuiteName   string
	Description string
	Priority    string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Passed      bool
}

// Test suite runners

func runBlueprintGenerationQualityTests(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	
	// Key validation points from comprehensive_blueprint_quality_atdd_test.go
	priorityBlueprints := map[string]BlueprintConfig{
		"cli-simple": {
			Type:           "cli",
			Complexity:     "simple",
			ExpectedFiles:  8,
			MaxFiles:       11,
			MinFiles:       7,
			Description:    "Simple CLI with minimal structure",
		},
		"cli-standard": {
			Type:           "cli",
			Complexity:     "standard",
			ExpectedFiles:  25,
			MaxFiles:       35,
			MinFiles:       20,
			Description:    "Full-featured CLI application",
		},
		"web-api-standard": {
			Type:           "web-api",
			Architecture:   "standard",
			ExpectedFiles:  20,
			MaxFiles:       30,
			MinFiles:       15,
			Description:    "Standard web API with REST endpoints",
		},
		"lambda-standard": {
			Type:           "lambda",
			Architecture:   "standard", 
			ExpectedFiles:  10,
			MaxFiles:       15,
			MinFiles:       8,
			Description:    "AWS Lambda serverless function",
		},
		"library-standard": {
			Type:           "library",
			Architecture:   "standard",
			ExpectedFiles:  8,
			MaxFiles:       12,
			MinFiles:       6,
			Description:    "Go library package",
		},
	}
	
	t.Log("Running blueprint generation quality validation...")
	
	for blueprintName, config := range priorityBlueprints {
		t.Run(blueprintName, func(t *testing.T) {
			projectName := fmt.Sprintf("quality-%s", blueprintName)
			
			generateProject(t, tmpDir, projectName, config)
			projectDir := filepath.Join(tmpDir, projectName)
			
			// Core validations
			validateProjectCompilation(t, projectDir, blueprintName)
			validateFileCount(t, projectDir, config)
			validateBasicStructure(t, projectDir, config)
		})
	}
	
	t.Log("✓ Blueprint generation quality validation completed")
}

func runLoggerIntegrationTests(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	
	t.Log("Running logger integration validation...")
	
	loggerTypes := []string{"slog", "zap", "logrus", "zerolog"}
	blueprints := []string{"cli-simple", "cli-standard", "web-api-standard"}
	
	for _, logger := range loggerTypes {
		for _, blueprint := range blueprints {
			t.Run(fmt.Sprintf("%s_with_%s", blueprint, logger), func(t *testing.T) {
				projectName := fmt.Sprintf("logger-%s-%s", blueprint, logger)
				
				config := BlueprintConfig{
					Type:       strings.Split(blueprint, "-")[0],
					Logger:     logger,
					MaxFiles:   50, // Generous limit for testing
					MinFiles:   5,
				}
				
				if strings.Contains(blueprint, "simple") {
					config.Complexity = "simple"
				} else if strings.Contains(blueprint, "standard") && strings.HasPrefix(blueprint, "cli") {
					config.Complexity = "standard"
				} else if strings.Contains(blueprint, "standard") {
					config.Architecture = "standard"
				}
				
				generateProject(t, tmpDir, projectName, config)
				projectDir := filepath.Join(tmpDir, projectName)
				
				// Validate logger integration and compilation
				validateLoggerIntegration(t, projectDir, logger, blueprint)
				validateProjectCompilation(t, projectDir, fmt.Sprintf("%s-%s", blueprint, logger))
			})
		}
	}
	
	t.Log("✓ Logger integration validation completed")
}

func runProgressiveComplexityTests(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	
	t.Log("Running progressive complexity validation...")
	
	// Test complexity reduction
	simpleConfig := BlueprintConfig{
		Type:           "cli",
		Complexity:     "simple", 
		ExpectedFiles:  8,
		MaxFiles:       11,
		MinFiles:       7,
	}
	
	standardConfig := BlueprintConfig{
		Type:           "cli",
		Complexity:     "standard",
		ExpectedFiles:  25,
		MaxFiles:       35,
		MinFiles:       20,
	}
	
	// Generate both variants
	generateProject(t, tmpDir, "complexity-simple", simpleConfig)
	generateProject(t, tmpDir, "complexity-standard", standardConfig)
	
	simpleDir := filepath.Join(tmpDir, "complexity-simple")
	standardDir := filepath.Join(tmpDir, "complexity-standard")
	
	simpleCount := countProjectFiles(t, simpleDir)
	standardCount := countProjectFiles(t, standardDir)
	
	// Validate complexity reduction (60-75% reduction)
	validateComplexityReduction(t, simpleCount, standardCount)
	
	// Both should compile
	validateProjectCompilation(t, simpleDir, "cli-simple-complexity")
	validateProjectCompilation(t, standardDir, "cli-standard-complexity")
	
	t.Log("✓ Progressive complexity validation completed")
}

func runCrossBlueprintValidationTests(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	
	t.Log("Running cross-blueprint validation...")
	
	blueprints := []string{"cli-simple", "cli-standard", "web-api-standard", "lambda-standard", "library-standard"}
	
	for _, blueprint := range blueprints {
		t.Run(blueprint, func(t *testing.T) {
			projectName := fmt.Sprintf("cross-%s", blueprint)
			
			config := BlueprintConfig{
				Type:     strings.Split(blueprint, "-")[0],
				MaxFiles: 50,
				MinFiles: 5,
			}
			
			if strings.Contains(blueprint, "simple") {
				config.Complexity = "simple"
			} else if strings.Contains(blueprint, "standard") && strings.HasPrefix(blueprint, "cli") {
				config.Complexity = "standard"
			} else if strings.Contains(blueprint, "standard") {
				config.Architecture = "standard"
			}
			
			generateProject(t, tmpDir, projectName, config)
			projectDir := filepath.Join(tmpDir, projectName)
			
			// Cross-validation checks
			validateProjectCompilation(t, projectDir, blueprint)
			validateBasicStructure(t, projectDir, config)
			validateGoDependencies(t, projectDir, blueprint)
		})
	}
	
	t.Log("✓ Cross-blueprint validation completed")
}

func runTemplateVariableResolutionTests(t *testing.T, tmpDir, projectRoot string) {
	t.Helper()
	
	t.Log("Running template variable resolution validation...")
	
	blueprints := []string{"cli-simple", "cli-standard", "web-api-standard", "lambda-standard", "library-standard"}
	
	for _, blueprint := range blueprints {
		t.Run(blueprint, func(t *testing.T) {
			projectName := fmt.Sprintf("template-%s", blueprint)
			
			config := BlueprintConfig{
				Type:     strings.Split(blueprint, "-")[0],
				MaxFiles: 50,
				MinFiles: 5,
			}
			
			if strings.Contains(blueprint, "simple") {
				config.Complexity = "simple"
			} else if strings.Contains(blueprint, "standard") && strings.HasPrefix(blueprint, "cli") {
				config.Complexity = "standard"
			} else if strings.Contains(blueprint, "standard") {
				config.Architecture = "standard"
			}
			
			generateProject(t, tmpDir, projectName, config)
			projectDir := filepath.Join(tmpDir, projectName)
			
			// Validate no unresolved template variables
			validateTemplateVariableResolution(t, projectDir, blueprint)
		})
	}
	
	t.Log("✓ Template variable resolution validation completed")
}

// generateATDDSummary generates a comprehensive summary of all test results
func generateATDDSummary(t *testing.T, results []ATDDTestResult, totalDuration time.Duration) {
	t.Helper()
	
	passedCount := 0
	failedCount := 0
	criticalPassed := 0
	criticalTotal := 0
	
	t.Log(strings.Repeat("=", 80))
	t.Log("COMPREHENSIVE ATDD QUALITY VALIDATION SUMMARY")
	t.Log(strings.Repeat("=", 80))
	
	for _, result := range results {
		status := "❌ FAILED"
		if result.Passed {
			status = "✅ PASSED"
			passedCount++
		} else {
			failedCount++
		}
		
		if result.Priority == "CRITICAL" {
			criticalTotal++
			if result.Passed {
				criticalPassed++
			}
		}
		
		t.Logf("%-30s %s (%v) [%s]", result.SuiteName, status, result.Duration, result.Priority)
		t.Logf("  Description: %s", result.Description)
	}
	
	t.Log(strings.Repeat("-", 80))
	t.Logf("TOTAL TESTS: %d", len(results))
	t.Logf("PASSED: %d", passedCount)
	t.Logf("FAILED: %d", failedCount)
	t.Logf("SUCCESS RATE: %.1f%%", float64(passedCount)/float64(len(results))*100)
	t.Logf("CRITICAL TESTS: %d/%d passed", criticalPassed, criticalTotal)
	t.Logf("TOTAL DURATION: %v", totalDuration)
	t.Log(strings.Repeat("-", 80))
	
	// Production readiness assessment
	if failedCount == 0 {
		t.Log("🎉 ALL TESTS PASSED - BLUEPRINT QUALITY MEETS ACCEPTANCE CRITERIA")
		t.Log("✓ Blueprint generation quality validated")
		t.Log("✓ Logger integration with simplified architecture validated")
		t.Log("✓ Progressive complexity with 60-75% reduction validated")
		t.Log("✓ Cross-blueprint consistency validated")
		t.Log("✓ Template variable resolution validated")
	} else if criticalPassed == criticalTotal && failedCount <= 2 {
		t.Log("⚠️  MOSTLY PASSED - MINOR ISSUES DETECTED")
		t.Log("Critical tests passed, but some non-critical issues need attention")
	} else {
		t.Error("❌ SIGNIFICANT FAILURES - BLUEPRINT QUALITY NEEDS IMPROVEMENT")
		t.Error("Critical acceptance criteria not met - review failed tests above")
	}
	
	t.Log(strings.Repeat("=", 80))
	
	// Assert overall success for CI/CD
	assert.Equal(t, 0, failedCount, "All ATDD quality tests should pass for production readiness")
	assert.Equal(t, criticalTotal, criticalPassed, "All critical tests must pass")
}