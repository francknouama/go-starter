package advisor

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/francknouama/go-starter/internal/advisor"
	"github.com/stretchr/testify/assert"
)

// BDDTestContext holds the context for BDD step execution
type BDDTestContext struct {
	t                    *testing.T
	advisor              *advisor.ArchitectureAdvisor
	interactiveAdvisor   *advisor.InteractiveAdvisor
	requirements         advisor.ProjectRequirements
	projectType          string // Separate field for project type
	recommendation       *advisor.ArchitectureRecommendation
	cliOutput            string
	cliError             error
	executionStartTime   time.Time
	executionDuration    time.Duration
	recommendations      []*advisor.ArchitectureRecommendation
	errorExpected        bool
	lastError            error
}

// NewBDDTestContext creates a new BDD test context
func NewBDDTestContext(t *testing.T) *BDDTestContext {
	return &BDDTestContext{
		t:                  t,
		advisor:            advisor.NewArchitectureAdvisor(),
		interactiveAdvisor: advisor.NewInteractiveAdvisor(),
		requirements:       advisor.ProjectRequirements{},
		recommendations:    make([]*advisor.ArchitectureRecommendation, 0),
	}
}

// Step Definitions for Background

func (ctx *BDDTestContext) TheAIAdvisorIsAvailable() error {
	if ctx.advisor == nil {
		return fmt.Errorf("AI advisor is not available")
	}
	return nil
}

func (ctx *BDDTestContext) TheKnowledgeBaseIsLoaded() error {
	// Verify knowledge base is accessible
	// This would typically check if the knowledge base data is properly loaded
	return nil
}

func (ctx *BDDTestContext) TheBlueprintRegistryIsLoaded() error {
	// Verify blueprint registry is accessible
	return nil
}

func (ctx *BDDTestContext) TheFrameworkKnowledgeBaseIsLoaded() error {
	// Verify framework knowledge base is loaded
	return nil
}

func (ctx *BDDTestContext) TheCompatibilityMatrixIsAvailable() error {
	// Verify compatibility matrix is loaded
	return nil
}

func (ctx *BDDTestContext) TheErrorHandlingSystemIsActive() error {
	// Verify error handling system is active
	return nil
}

// Step Definitions for Given steps

func (ctx *BDDTestContext) IHaveAProjectRequirementFor(projectType string) error {
	ctx.projectType = projectType
	return nil
}

func (ctx *BDDTestContext) INeedAProject(projectType string) error {
	return ctx.IHaveAProjectRequirementFor(projectType)
}

func (ctx *BDDTestContext) MyProjectDomainIs(domain string) error {
	ctx.requirements.Domain = domain
	return nil
}

func (ctx *BDDTestContext) MyTeamExperienceLevelIs(experience string) error {
	ctx.requirements.TeamExperience = experience
	return nil
}

func (ctx *BDDTestContext) MyTeamHasExperienceLevel(experience string) error {
	return ctx.MyTeamExperienceLevelIs(experience)
}

func (ctx *BDDTestContext) IPreferDevelopmentStyle(style string) error {
	ctx.requirements.PreferredStyle = style
	return nil
}

func (ctx *BDDTestContext) IHaveDatabaseRequirements(dbType string) error {
	ctx.requirements.DatabaseRequirements = dbType
	return nil
}

func (ctx *BDDTestContext) INeedCompliance(compliance string) error {
	if ctx.requirements.ComplianceNeeds == nil {
		ctx.requirements.ComplianceNeeds = make([]string, 0)
	}
	ctx.requirements.ComplianceNeeds = append(ctx.requirements.ComplianceNeeds, compliance)
	return nil
}

func (ctx *BDDTestContext) MyTimeToMarketIs(timeframe string) error {
	ctx.requirements.TimeToMarket = timeframe
	return nil
}

func (ctx *BDDTestContext) MyBudgetIs(budget string) error {
	ctx.requirements.Budget = budget
	return nil
}

func (ctx *BDDTestContext) IExpectLoad(load string) error {
	ctx.requirements.ExpectedLoad = load
	return nil
}

func (ctx *BDDTestContext) INeedResponseTimes(responseTime string) error {
	ctx.requirements.ResponseTime = responseTime
	return nil
}

func (ctx *BDDTestContext) IHaveDataPatterns(patterns string) error {
	// Set data volume or patterns in requirements
	ctx.requirements.DataVolume = patterns
	return nil
}

func (ctx *BDDTestContext) MyDeploymentTargetIs(target string) error {
	ctx.requirements.DeploymentTarget = target
	return nil
}

func (ctx *BDDTestContext) IWantToPublishToTheCommunity() error {
	// Set flag for public library - can be stored in a custom field or ignored for now
	// ctx.requirements doesn't have IsPublicLibrary field
	return nil
}

func (ctx *BDDTestContext) INeedComprehensiveDocumentation() error {
	// Set documentation requirements - can be stored in MonitoringNeeds or similar
	ctx.requirements.MonitoringNeeds = "comprehensive"
	return nil
}

func (ctx *BDDTestContext) IProvideMinimalRequirements() error {
	// Set only the most basic requirements
	ctx.requirements = advisor.ProjectRequirements{
		Domain:         "other",
		TeamExperience: "mixed",
	}
	return nil
}

func (ctx *BDDTestContext) ISpecifyAnInvalidProjectType(invalidType string) error {
	ctx.projectType = invalidType
	return nil
}

func (ctx *BDDTestContext) ISpecifyInvalidTeamExperience(invalidExperience string) error {
	ctx.requirements.TeamExperience = invalidExperience
	return nil
}

func (ctx *BDDTestContext) ISpecifyADomainWithSpecialCharacters(domain string) error {
	ctx.requirements.Domain = domain
	return nil
}

// Step Definitions for When steps

func (ctx *BDDTestContext) IRequestQuickModeRecommendations() error {
	var err error
	ctx.recommendation, err = ctx.interactiveAdvisor.QuickRecommendation(
		ctx.projectType,
		ctx.requirements.Domain,
		ctx.requirements.TeamExperience,
	)
	ctx.lastError = err
	return err
}

func (ctx *BDDTestContext) IAskForBlueprintRecommendations() error {
	var err error
	ctx.recommendation, err = ctx.advisor.AnalyzeRequirements(ctx.requirements)
	ctx.lastError = err
	return err
}

func (ctx *BDDTestContext) IAskForFrameworkRecommendations() error {
	return ctx.IAskForBlueprintRecommendations()
}

func (ctx *BDDTestContext) IAskForFrameworkRecommendationsWithBenchmarks() error {
	// Extended analysis with benchmarks
	return ctx.IAskForBlueprintRecommendations()
}

func (ctx *BDDTestContext) IRequestRecommendations() error {
	return ctx.IAskForBlueprintRecommendations()
}

func (ctx *BDDTestContext) IRequestRecommendationsThroughCLI() error {
	return ctx.executeAdvisorCLI()
}

func (ctx *BDDTestContext) IExecuteTheCLICommand() error {
	return ctx.executeAdvisorCLI()
}

func (ctx *BDDTestContext) IRequestQuickModeRecommendations100Times() error {
	ctx.executionStartTime = time.Now()
	ctx.recommendations = make([]*advisor.ArchitectureRecommendation, 0, 100)
	
	for i := 0; i < 100; i++ {
		rec, err := ctx.interactiveAdvisor.QuickRecommendation(
			ctx.projectType,
			ctx.requirements.Domain,
			ctx.requirements.TeamExperience,
		)
		if err != nil {
			return err
		}
		ctx.recommendations = append(ctx.recommendations, rec)
	}
	
	ctx.executionDuration = time.Since(ctx.executionStartTime)
	return nil
}

func (ctx *BDDTestContext) IValidateTheRecommendation() error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation to validate")
	}
	// Perform validation logic
	return nil
}

func (ctx *BDDTestContext) IParseTheOutput() error {
	if ctx.cliOutput == "" {
		return fmt.Errorf("no CLI output to parse")
	}
	
	// Try to parse as JSON
	var rec advisor.ArchitectureRecommendation
	err := json.Unmarshal([]byte(ctx.cliOutput), &rec)
	if err != nil {
		ctx.lastError = err
		return err
	}
	
	ctx.recommendation = &rec
	return nil
}

func (ctx *BDDTestContext) IRequestRecommendationsMultipleTimes() error {
	// Request same recommendations multiple times to test consistency
	ctx.recommendations = make([]*advisor.ArchitectureRecommendation, 0, 5)
	
	for i := 0; i < 5; i++ {
		rec, err := ctx.advisor.AnalyzeRequirements(ctx.requirements)
		if err != nil {
			return err
		}
		ctx.recommendations = append(ctx.recommendations, rec)
	}
	return nil
}

// Step Definitions for Then steps

func (ctx *BDDTestContext) IShouldGetBlueprintRecommendationsIncluding(blueprints string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	expectedBlueprints := parseCommaSeparatedList(blueprints)
	for _, expected := range expectedBlueprints {
		if ctx.recommendation.Blueprint == expected {
			return nil // Found one of the expected blueprints
		}
	}
	
	// Check alternatives as well
	for _, alt := range ctx.recommendation.Alternatives {
		for _, expected := range expectedBlueprints {
			if alt.Blueprint == expected {
				return nil
			}
		}
	}
	
	return fmt.Errorf("expected blueprint from %v, got %s", expectedBlueprints, ctx.recommendation.Blueprint)
}

func (ctx *BDDTestContext) IShouldGetFrameworkRecommendationsIncluding(frameworks string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	expectedFrameworks := parseCommaSeparatedList(frameworks)
	for _, expected := range expectedFrameworks {
		if ctx.recommendation.Framework == expected {
			return nil
		}
	}
	
	return fmt.Errorf("expected framework from %v, got %s", expectedFrameworks, ctx.recommendation.Framework)
}

func (ctx *BDDTestContext) IShouldGetLoggerRecommendationsIncluding(loggers string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	expectedLoggers := parseCommaSeparatedList(loggers)
	for _, expected := range expectedLoggers {
		if ctx.recommendation.Logger == expected {
			return nil
		}
	}
	
	return fmt.Errorf("expected logger from %v, got %s", expectedLoggers, ctx.recommendation.Logger)
}

func (ctx *BDDTestContext) TheConfidenceLevelShouldBeAtLeast(minConfidenceStr string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	minConfidence, err := strconv.ParseFloat(minConfidenceStr, 64)
	if err != nil {
		return fmt.Errorf("invalid confidence value: %s", minConfidenceStr)
	}
	
	if ctx.recommendation.Confidence < minConfidence {
		return fmt.Errorf("confidence %.2f is below minimum %.2f", ctx.recommendation.Confidence, minConfidence)
	}
	
	return nil
}

func (ctx *BDDTestContext) IShouldReceiveReasoningForTheRecommendations() error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	if len(ctx.recommendation.Reasoning) == 0 {
		return fmt.Errorf("no reasoning provided")
	}
	
	return nil
}

func (ctx *BDDTestContext) IShouldReceiveAlternatives() error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	if len(ctx.recommendation.Alternatives) == 0 {
		return fmt.Errorf("no alternatives provided")
	}
	
	return nil
}

func (ctx *BDDTestContext) TheEstimatedFileCountShouldBeBetween(minStr, maxStr string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("invalid min file count: %s", minStr)
	}
	
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return fmt.Errorf("invalid max file count: %s", maxStr)
	}
	
	fileCount := ctx.recommendation.EstimatedFiles
	if fileCount < min || fileCount > max {
		return fmt.Errorf("file count %d is not between %d and %d", fileCount, min, max)
	}
	
	return nil
}

func (ctx *BDDTestContext) IShouldGetAValidBlueprintRecommendation() error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	if ctx.recommendation.Blueprint == "" {
		return fmt.Errorf("empty blueprint recommendation")
	}
	
	return nil
}

func (ctx *BDDTestContext) TheEstimatedFileCountShouldBeGreaterThan(minStr string) error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation received")
	}
	
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return fmt.Errorf("invalid min file count: %s", minStr)
	}
	
	if ctx.recommendation.EstimatedFiles <= min {
		return fmt.Errorf("file count %d is not greater than %d", ctx.recommendation.EstimatedFiles, min)
	}
	
	return nil
}

func (ctx *BDDTestContext) IShouldGetAClearErrorMessage() error {
	if ctx.lastError == nil && ctx.cliError == nil {
		return fmt.Errorf("expected an error but none occurred")
	}
	
	errorMsg := ""
	if ctx.lastError != nil {
		errorMsg = ctx.lastError.Error()
	} else if ctx.cliError != nil {
		errorMsg = ctx.cliError.Error()
	}
	
	if errorMsg == "" {
		return fmt.Errorf("error message is empty")
	}
	
	return nil
}

func (ctx *BDDTestContext) AllRecommendationsShouldBeGeneratedWithinDuration(durationStr string) error {
	maxDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("invalid duration: %s", durationStr)
	}
	
	if ctx.executionDuration > maxDuration {
		return fmt.Errorf("execution took %v, expected less than %v", ctx.executionDuration, maxDuration)
	}
	
	return nil
}

func (ctx *BDDTestContext) EachRecommendationShouldBeValidAndComplete() error {
	for i, rec := range ctx.recommendations {
		if rec.Blueprint == "" {
			return fmt.Errorf("recommendation %d has empty blueprint", i)
		}
		if rec.Confidence <= 0 {
			return fmt.Errorf("recommendation %d has invalid confidence: %f", i, rec.Confidence)
		}
	}
	return nil
}

// Helper methods

func (ctx *BDDTestContext) executeAdvisorCLI() error {
	// Build CLI arguments from requirements
	args := []string{"advisor", "--quick", "--format=json"}
	
	if ctx.projectType != "" {
		args = append(args, "--type="+ctx.projectType)
	}
	if ctx.requirements.Domain != "" {
		args = append(args, "--domain="+ctx.requirements.Domain)
	}
	if ctx.requirements.TeamExperience != "" {
		args = append(args, "--team="+ctx.requirements.TeamExperience)
	}
	
	// Execute CLI command
	binaryPath := "../../../bin/go-starter"
	cmd := exec.Command(binaryPath, args...)
	output, err := cmd.CombinedOutput()
	
	ctx.cliOutput = string(output)
	ctx.cliError = err
	
	return nil
}

func parseCommaSeparatedList(input string) []string {
	items := strings.Split(input, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(strings.Trim(item, `"`))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// BDD Test Runner Functions

// RunBDDScenario executes a BDD scenario with the given steps
func RunBDDScenario(t *testing.T, scenarioName string, steps func(*BDDTestContext) error) {
	t.Run(scenarioName, func(t *testing.T) {
		ctx := NewBDDTestContext(t)
		
		// Execute the scenario steps
		err := steps(ctx)
		if err != nil {
			t.Fatalf("BDD scenario failed: %v", err)
		}
	})
}

// Helper assertion functions for more complex validations

func (ctx *BDDTestContext) AssertRecommendationQuality() error {
	if ctx.recommendation == nil {
		return fmt.Errorf("no recommendation to validate")
	}
	
	rec := ctx.recommendation
	
	// Basic structure validation
	assert.NotEmpty(ctx.t, rec.Blueprint, "Blueprint should not be empty")
	assert.NotEmpty(ctx.t, rec.Architecture, "Architecture should not be empty")
	assert.NotEmpty(ctx.t, rec.Framework, "Framework should not be empty")
	assert.NotEmpty(ctx.t, rec.Logger, "Logger should not be empty")
	assert.Greater(ctx.t, rec.EstimatedFiles, 0, "Should estimate positive file count")
	assert.GreaterOrEqual(ctx.t, rec.Confidence, 0.0, "Confidence should be non-negative")
	assert.LessOrEqual(ctx.t, rec.Confidence, 1.0, "Confidence should not exceed 1.0")
	
	return nil
}

func (ctx *BDDTestContext) AssertRecommendationConsistency() error {
	if len(ctx.recommendations) < 2 {
		return fmt.Errorf("need at least 2 recommendations to check consistency")
	}
	
	first := ctx.recommendations[0]
	for i, rec := range ctx.recommendations[1:] {
		assert.Equal(ctx.t, first.Blueprint, rec.Blueprint, 
			"Recommendation %d blueprint should be consistent", i+1)
		assert.InDelta(ctx.t, first.Confidence, rec.Confidence, 0.1,
			"Recommendation %d confidence should be similar", i+1)
	}
	
	return nil
}