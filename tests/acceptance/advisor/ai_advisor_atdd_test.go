package advisor

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/francknouama/go-starter/internal/advisor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ATDD Test Suite for AI-Powered Architecture Advisor
// These tests validate the advisor from a user acceptance perspective

func TestAIAdvisor_ATDD_QuickModeRecommendations(t *testing.T) {
	tests := []struct {
		name           string
		projectType    string
		domain         string
		teamExperience string
		expectations   AdvisorExpectations
	}{
		{
			name:           "E-commerce API for Senior Team",
			projectType:    "api",
			domain:         "e-commerce",
			teamExperience: "senior",
			expectations: AdvisorExpectations{
				shouldRecommendBlueprint: []string{"web-api-ddd", "web-api-clean", "web-api"},
				shouldRecommendFramework: []string{"gin", "echo", "fiber"},
				minConfidence:           0.6,
				shouldProvideReasoning:  true,
				shouldProvideAlternatives: true,
				maxEstimatedFiles:       100,
				minEstimatedFiles:       20,
			},
		},
		{
			name:           "CLI Tool for Junior Team",
			projectType:    "cli",
			domain:         "devtools",
			teamExperience: "junior",
			expectations: AdvisorExpectations{
				shouldRecommendBlueprint: []string{"cli-simple", "cli"},
				shouldRecommendFramework: []string{"cobra"},
				minConfidence:           0.5,
				shouldProvideReasoning:  true,
				shouldProvideAlternatives: true,
				maxEstimatedFiles:       35,
				minEstimatedFiles:       5,
			},
		},
		{
			name:           "Fintech API for Expert Team",
			projectType:    "api",
			domain:         "fintech",
			teamExperience: "expert",
			expectations: AdvisorExpectations{
				shouldRecommendBlueprint: []string{"web-api-ddd", "web-api-hexagonal"},
				shouldRecommendFramework: []string{"gin", "echo"},
				shouldRecommendLogger:    []string{"logrus", "zap"},
				minConfidence:           0.7,
				shouldProvideReasoning:  true,
				shouldProvideAlternatives: true,
				maxEstimatedFiles:       150,
				minEstimatedFiles:       50,
			},
		},
		{
			name:           "IoT Lambda for Mixed Team",
			projectType:    "lambda",
			domain:         "iot",
			teamExperience: "mixed",
			expectations: AdvisorExpectations{
				shouldRecommendBlueprint: []string{"lambda-event-processing", "lambda"},
				shouldRecommendLogger:    []string{"slog", "zap"},
				minConfidence:           0.5,
				shouldProvideReasoning:  true,
				shouldProvideAlternatives: true,
				maxEstimatedFiles:       50,
				minEstimatedFiles:       10,
			},
		},
		{
			name:           "Microservice for Expert Team",
			projectType:    "microservice",
			domain:         "e-commerce",
			teamExperience: "expert",
			expectations: AdvisorExpectations{
				shouldRecommendBlueprint: []string{"microservice"},
				shouldRecommendFramework: []string{"gin", "fiber"},
				minConfidence:           0.6,
				shouldProvideReasoning:  true,
				shouldProvideAlternatives: true,
				maxEstimatedFiles:       80,
				minEstimatedFiles:       40,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: Getting a quick recommendation
			recommendation := whenGettingQuickRecommendation(t, tt.projectType, tt.domain, tt.teamExperience)

			// Then: Validate all expectations
			thenValidateRecommendation(t, recommendation, tt.expectations)
		})
	}
}

func TestAIAdvisor_ATDD_CLIIntegration(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		expectError bool
		validations []func(*testing.T, string)
	}{
		{
			name: "Quick mode with all required flags",
			args: []string{"advisor", "--quick", "--type=api", "--domain=e-commerce", "--team=senior", "--format=json"},
			validations: []func(*testing.T, string){
				validateJSONOutput,
				validateRecommendationFields,
				validateConfidenceLevel,
			},
		},
		{
			name: "Quick mode with summary format",
			args: []string{"advisor", "--quick", "--type=cli", "--domain=devtools", "--team=junior", "--format=summary"},
			validations: []func(*testing.T, string){
				validateSummaryOutput,
				validateBlueprintRecommendation,
			},
		},
		{
			name:        "Quick mode missing required flags",
			args:        []string{"advisor", "--quick", "--type=api"},
			expectError: true,
		},
		{
			name:        "Invalid project type",
			args:        []string{"advisor", "--quick", "--type=invalid", "--domain=e-commerce", "--team=senior"},
			expectError: true,
		},
		{
			name:        "Invalid team experience",
			args:        []string{"advisor", "--quick", "--type=api", "--domain=e-commerce", "--team=invalid"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// When: Running the CLI command
			output, err := whenRunningAdvisorCLI(t, tt.args)

			// Then: Validate expectations
			if tt.expectError {
				assert.Error(t, err, "Expected command to fail")
			} else {
				require.NoError(t, err, "Command should succeed")
				
				// Run all validations
				for _, validation := range tt.validations {
					validation(t, output)
				}
			}
		})
	}
}

func TestAIAdvisor_ATDD_EdgeCases(t *testing.T) {
	t.Run("Empty requirements should provide default recommendation", func(t *testing.T) {
		// Given: Minimal requirements
		req := advisor.ProjectRequirements{
			Domain:         "other",
			TeamExperience: "mixed",
		}

		// When: Getting recommendation
		adv := advisor.NewArchitectureAdvisor()
		recommendation, err := adv.AnalyzeRequirements(req)

		// Then: Should still provide a recommendation
		require.NoError(t, err)
		assert.NotEmpty(t, recommendation.Blueprint)
		assert.Greater(t, recommendation.Confidence, 0.0)
	})

	t.Run("Conflicting requirements should handle gracefully", func(t *testing.T) {
		// Given: Conflicting requirements (urgent timeline but complex needs)
		req := advisor.ProjectRequirements{
			Domain:               "fintech",
			TeamExperience:       "junior",
			TimeToMarket:        "mvp",
			ExpectedLoad:        "massive",
			DatabaseRequirements: "distributed",
			ComplianceNeeds:     []string{"sox", "hipaa"},
		}

		// When: Getting recommendation
		adv := advisor.NewArchitectureAdvisor()
		recommendation, err := adv.AnalyzeRequirements(req)

		// Then: Should provide reasonable compromise
		require.NoError(t, err)
		assert.NotEmpty(t, recommendation.Blueprint)
		assert.Contains(t, strings.Join(recommendation.Reasoning, " "), "compromise",
			"Should acknowledge conflicting requirements")
	})
}

func TestAIAdvisor_ATDD_PerformanceRequirements(t *testing.T) {
	advisorInstance := advisor.NewArchitectureAdvisor()

	t.Run("Recommendation generation should be fast", func(t *testing.T) {
		// Given: Standard requirements
		req := advisor.ProjectRequirements{
			Domain:         "e-commerce",
			TeamExperience: "mixed",
			ExpectedLoad:   "medium",
		}

		// When: Measuring time for 100 recommendations
		start := time.Now()
		for i := 0; i < 100; i++ {
			_, err := advisorInstance.AnalyzeRequirements(req)
			require.NoError(t, err)
		}
		duration := time.Since(start)

		// Then: Should be fast enough for interactive use
		assert.Less(t, duration.Seconds(), 1.0, "100 recommendations should take less than 1 second")
	})
}

// Test Helper Types and Functions

type AdvisorExpectations struct {
	shouldRecommendBlueprint  []string
	shouldRecommendFramework  []string
	shouldRecommendLogger     []string
	minConfidence            float64
	shouldProvideReasoning   bool
	shouldProvideAlternatives bool
	maxEstimatedFiles        int
	minEstimatedFiles        int
}

func whenGettingQuickRecommendation(t *testing.T, projectType, domain, teamExperience string) *advisor.ArchitectureRecommendation {
	interactiveAdvisor := advisor.NewInteractiveAdvisor()
	recommendation, err := interactiveAdvisor.QuickRecommendation(projectType, domain, teamExperience)
	require.NoError(t, err, "Quick recommendation should not fail")
	require.NotNil(t, recommendation, "Recommendation should not be nil")
	return recommendation
}

func thenValidateRecommendation(t *testing.T, rec *advisor.ArchitectureRecommendation, exp AdvisorExpectations) {
	// Validate blueprint recommendation
	if len(exp.shouldRecommendBlueprint) > 0 {
		assert.Contains(t, exp.shouldRecommendBlueprint, rec.Blueprint,
			"Blueprint %s should be one of %v", rec.Blueprint, exp.shouldRecommendBlueprint)
	}

	// Validate framework recommendation
	if len(exp.shouldRecommendFramework) > 0 {
		assert.Contains(t, exp.shouldRecommendFramework, rec.Framework,
			"Framework %s should be one of %v", rec.Framework, exp.shouldRecommendFramework)
	}

	// Validate logger recommendation
	if len(exp.shouldRecommendLogger) > 0 {
		assert.Contains(t, exp.shouldRecommendLogger, rec.Logger,
			"Logger %s should be one of %v", rec.Logger, exp.shouldRecommendLogger)
	}

	// Validate confidence
	assert.GreaterOrEqual(t, rec.Confidence, exp.minConfidence,
		"Confidence %.2f should be at least %.2f", rec.Confidence, exp.minConfidence)

	// Validate reasoning
	if exp.shouldProvideReasoning {
		assert.NotEmpty(t, rec.Reasoning, "Should provide reasoning")
	}

	// Validate alternatives
	if exp.shouldProvideAlternatives {
		assert.NotEmpty(t, rec.Alternatives, "Should provide alternatives")
		assert.LessOrEqual(t, len(rec.Alternatives), 3, "Should provide at most 3 alternatives")
	}

	// Validate file estimates
	if exp.maxEstimatedFiles > 0 {
		assert.LessOrEqual(t, rec.EstimatedFiles, exp.maxEstimatedFiles,
			"Estimated files %d should be <= %d", rec.EstimatedFiles, exp.maxEstimatedFiles)
	}
	if exp.minEstimatedFiles > 0 {
		assert.GreaterOrEqual(t, rec.EstimatedFiles, exp.minEstimatedFiles,
			"Estimated files %d should be >= %d", rec.EstimatedFiles, exp.minEstimatedFiles)
	}

	// Validate basic structure
	assert.NotEmpty(t, rec.Blueprint, "Blueprint should not be empty")
	assert.NotEmpty(t, rec.Architecture, "Architecture should not be empty")
	assert.NotEmpty(t, rec.Framework, "Framework should not be empty")
	assert.NotEmpty(t, rec.Logger, "Logger should not be empty")
	assert.NotEmpty(t, rec.DevelopmentTime, "Development time should not be empty")
	assert.Greater(t, rec.EstimatedFiles, 0, "Should estimate positive file count")
}

func whenRunningAdvisorCLI(t *testing.T, args []string) (string, error) {
	// Build the go-starter binary if needed
	binaryPath := "../../../bin/go-starter"
	
	// Run the command
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	
	return string(output), err
}

// CLI Output Validators

func validateJSONOutput(t *testing.T, output string) {
	var rec advisor.ArchitectureRecommendation
	err := json.Unmarshal([]byte(output), &rec)
	require.NoError(t, err, "Output should be valid JSON")
	assert.NotEmpty(t, rec.Blueprint, "JSON should contain blueprint")
}

func validateRecommendationFields(t *testing.T, output string) {
	var rec advisor.ArchitectureRecommendation
	json.Unmarshal([]byte(output), &rec)
	
	assert.NotEmpty(t, rec.Blueprint)
	assert.NotEmpty(t, rec.Architecture)
	assert.NotEmpty(t, rec.Framework)
	assert.NotEmpty(t, rec.Logger)
	assert.Greater(t, rec.EstimatedFiles, 0)
}

func validateConfidenceLevel(t *testing.T, output string) {
	var rec advisor.ArchitectureRecommendation
	json.Unmarshal([]byte(output), &rec)
	
	assert.GreaterOrEqual(t, rec.Confidence, 0.0)
	assert.LessOrEqual(t, rec.Confidence, 1.0)
}

func validateSummaryOutput(t *testing.T, output string) {
	assert.Contains(t, output, "Blueprint:")
	assert.Contains(t, output, "Architecture:")
	assert.Contains(t, output, "Confidence:")
	assert.Contains(t, output, "Framework:")
}

func validateBlueprintRecommendation(t *testing.T, output string) {
	// Should contain a valid blueprint name
	validBlueprints := []string{
		"cli-simple", "cli", "web-api", "web-api-clean", 
		"web-api-ddd", "web-api-hexagonal", "lambda", 
		"lambda-event-processing", "microservice", "library",
	}
	
	foundValidBlueprint := false
	for _, blueprint := range validBlueprints {
		if strings.Contains(output, blueprint) {
			foundValidBlueprint = true
			break
		}
	}
	
	assert.True(t, foundValidBlueprint, "Output should contain a valid blueprint recommendation")
}

// Integration Test with Real CLI Binary

func TestAIAdvisor_ATDD_RealCLIIntegration(t *testing.T) {
	// Skip if binary doesn't exist
	binaryPath := "../../../bin/go-starter"
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		t.Skip("go-starter binary not found, run 'make build' first")
	}

	t.Run("Real CLI advisor command execution", func(t *testing.T) {
		// When: Running real CLI command
		cmd := exec.Command(binaryPath, "advisor", "--quick", 
			"--type=api", "--domain=e-commerce", "--team=senior", "--format=json")
		
		output, err := cmd.CombinedOutput()
		
		// Then: Should succeed and provide valid recommendation
		require.NoError(t, err, "CLI command should succeed: %s", string(output))
		
		var rec advisor.ArchitectureRecommendation
		err = json.Unmarshal(output, &rec)
		require.NoError(t, err, "Should output valid JSON")
		
		// Validate the recommendation
		assert.NotEmpty(t, rec.Blueprint)
		assert.GreaterOrEqual(t, rec.Confidence, 0.0)
		assert.LessOrEqual(t, rec.Confidence, 1.0)
		assert.NotEmpty(t, rec.Reasoning)
	})
}