package advisor

import (
	"fmt"
	"testing"

	"github.com/francknouama/go-starter/internal/advisor"
)

// TestFixtures provides pre-configured test data for BDD scenarios
type TestFixtures struct {
	// Common project requirements
	StandardRequirements     advisor.ProjectRequirements
	MinimalRequirements      advisor.ProjectRequirements
	ComplexRequirements      advisor.ProjectRequirements
	ConflictingRequirements  advisor.ProjectRequirements
	
	// Expected recommendation patterns
	ExpectedCLIBlueprints    []string
	ExpectedAPIBlueprints    []string
	ExpectedFrameworks       map[string][]string // project type -> frameworks
	ExpectedLoggers          []string
	
	// Performance expectations
	MaxRecommendationTime    int64 // milliseconds
	MinConfidenceThreshold   float64
}

// NewTestFixtures creates standardized test fixtures
func NewTestFixtures() *TestFixtures {
	return &TestFixtures{
		StandardRequirements: advisor.ProjectRequirements{
			Domain:         "e-commerce",
			TeamExperience: "senior",
			TimeToMarket:   "standard",
			ExpectedLoad:   "medium",
			ResponseTime:   "standard",
			Budget:         "medium",
			TeamSize:       3,
		},
		
		MinimalRequirements: advisor.ProjectRequirements{
			Domain:         "other",
			TeamExperience: "mixed",
		},
		
		ComplexRequirements: advisor.ProjectRequirements{
			Domain:               "fintech",
			TeamExperience:       "expert",
			TimeToMarket:        "thorough",
			ExpectedLoad:        "massive",
			ResponseTime:        "realtime",
			DatabaseRequirements: "distributed",
			ComplianceNeeds:     []string{"sox", "pci"},
			AuthRequirements:    "enterprise",
			DeploymentTarget:    "cloud",
			MonitoringNeeds:     "enterprise",
		},
		
		ConflictingRequirements: advisor.ProjectRequirements{
			Domain:               "fintech",
			TeamExperience:       "junior",
			TimeToMarket:        "mvp",
			ExpectedLoad:        "massive",
			DatabaseRequirements: "distributed",
			ComplianceNeeds:     []string{"sox", "hipaa"},
		},
		
		ExpectedCLIBlueprints: []string{"cli-simple", "cli"},
		ExpectedAPIBlueprints: []string{"web-api", "web-api-clean", "web-api-ddd", "web-api-hexagonal"},
		
		ExpectedFrameworks: map[string][]string{
			"api": {"gin", "echo", "fiber", "chi"},
			"cli": {"cobra"},
			"microservice": {"gin", "fiber"},
		},
		
		ExpectedLoggers: []string{"slog", "zap", "logrus", "zerolog"},
		
		MaxRecommendationTime:  100, // 100ms per recommendation
		MinConfidenceThreshold: 0.3, // Minimum acceptable confidence
	}
}

// BDDTestHelpers provides utility functions for BDD testing
type BDDTestHelpers struct {
	fixtures *TestFixtures
}

// NewBDDTestHelpers creates new test helpers
func NewBDDTestHelpers() *BDDTestHelpers {
	return &BDDTestHelpers{
		fixtures: NewTestFixtures(),
	}
}

// ValidateRecommendationStructure checks basic recommendation structure
func (h *BDDTestHelpers) ValidateRecommendationStructure(t *testing.T, rec *advisor.ArchitectureRecommendation) error {
	if rec == nil {
		return fmt.Errorf("recommendation is nil")
	}
	
	if rec.Blueprint == "" {
		return fmt.Errorf("blueprint is empty")
	}
	
	if rec.Architecture == "" {
		return fmt.Errorf("architecture is empty")
	}
	
	if rec.Framework == "" {
		return fmt.Errorf("framework is empty")
	}
	
	if rec.Logger == "" {
		return fmt.Errorf("logger is empty")
	}
	
	if rec.Confidence < 0 || rec.Confidence > 1 {
		return fmt.Errorf("confidence %f is out of range [0,1]", rec.Confidence)
	}
	
	if rec.EstimatedFiles <= 0 {
		return fmt.Errorf("estimated files %d should be positive", rec.EstimatedFiles)
	}
	
	if len(rec.Reasoning) == 0 {
		return fmt.Errorf("reasoning is empty")
	}
	
	return nil
}

// ValidateBlueprintExists checks if blueprint is in expected list
func (h *BDDTestHelpers) ValidateBlueprintExists(blueprint string, expectedList []string) error {
	for _, expected := range expectedList {
		if blueprint == expected {
			return nil
		}
	}
	return fmt.Errorf("blueprint %s not found in expected list %v", blueprint, expectedList)
}

// ValidateFrameworkCompatibility checks framework compatibility with project type
func (h *BDDTestHelpers) ValidateFrameworkCompatibility(projectType, framework string) error {
	expectedFrameworks, exists := h.fixtures.ExpectedFrameworks[projectType]
	if !exists {
		// If we don't have specific expectations, any framework is acceptable
		return nil
	}
	
	for _, expected := range expectedFrameworks {
		if framework == expected {
			return nil
		}
	}
	
	return fmt.Errorf("framework %s not compatible with project type %s, expected one of %v", 
		framework, projectType, expectedFrameworks)
}

// ValidateLoggerChoice checks if logger is in acceptable list
func (h *BDDTestHelpers) ValidateLoggerChoice(logger string) error {
	for _, expected := range h.fixtures.ExpectedLoggers {
		if logger == expected {
			return nil
		}
	}
	return fmt.Errorf("logger %s not in expected list %v", logger, h.fixtures.ExpectedLoggers)
}

// ValidateConfidenceLevel checks confidence is within acceptable range
func (h *BDDTestHelpers) ValidateConfidenceLevel(confidence float64, minExpected float64) error {
	if confidence < minExpected {
		return fmt.Errorf("confidence %f below minimum expected %f", confidence, minExpected)
	}
	
	if confidence < h.fixtures.MinConfidenceThreshold {
		return fmt.Errorf("confidence %f below global minimum threshold %f", 
			confidence, h.fixtures.MinConfidenceThreshold)
	}
	
	return nil
}

// ValidateFileCountRange checks if file count is in reasonable range
func (h *BDDTestHelpers) ValidateFileCountRange(fileCount int, projectType string, complexity string) error {
	// Define reasonable ranges based on project type and complexity
	ranges := map[string]map[string][2]int{
		"cli": {
			"simple":   {5, 15},
			"standard": {15, 40},
			"advanced": {30, 60},
		},
		"api": {
			"simple":   {15, 30},
			"standard": {25, 60},
			"advanced": {50, 120},
		},
		"microservice": {
			"standard": {30, 80},
			"advanced": {60, 150},
		},
		"lambda": {
			"simple":   {5, 20},
			"standard": {10, 40},
		},
	}
	
	if projectRanges, exists := ranges[projectType]; exists {
		if complexityRange, exists := projectRanges[complexity]; exists {
			min, max := complexityRange[0], complexityRange[1]
			if fileCount < min || fileCount > max {
				return fmt.Errorf("file count %d for %s/%s project outside expected range [%d, %d]",
					fileCount, projectType, complexity, min, max)
			}
		}
	}
	
	// Global sanity check
	if fileCount < 1 || fileCount > 500 {
		return fmt.Errorf("file count %d outside global reasonable range [1, 500]", fileCount)
	}
	
	return nil
}

// ValidateReasoningQuality checks if reasoning is meaningful
func (h *BDDTestHelpers) ValidateReasoningQuality(reasoning []string) error {
	if len(reasoning) == 0 {
		return fmt.Errorf("reasoning is empty")
	}
	
	// Check for meaningful content
	totalLength := 0
	for _, reason := range reasoning {
		if len(reason) < 10 {
			return fmt.Errorf("reasoning item too short: '%s'", reason)
		}
		totalLength += len(reason)
	}
	
	if totalLength < 50 {
		return fmt.Errorf("total reasoning length %d too short for meaningful explanation", totalLength)
	}
	
	return nil
}

// ValidateAlternatives checks if alternatives are meaningful
func (h *BDDTestHelpers) ValidateAlternatives(alternatives []advisor.AlternativeRecommendation) error {
	if len(alternatives) == 0 {
		return fmt.Errorf("no alternatives provided")
	}
	
	if len(alternatives) > 5 {
		return fmt.Errorf("too many alternatives: %d (max 5)", len(alternatives))
	}
	
	for i, alt := range alternatives {
		if alt.Blueprint == "" {
			return fmt.Errorf("alternative %d has empty blueprint", i)
		}
		if alt.Confidence < 0 || alt.Confidence > 1 {
			return fmt.Errorf("alternative %d confidence %f out of range", i, alt.Confidence)
		}
		if alt.Reason == "" {
			return fmt.Errorf("alternative %d has empty reason", i)
		}
	}
	
	return nil
}

// CreateTestScenario creates a standardized test scenario
func (h *BDDTestHelpers) CreateTestScenario(
	scenarioName string,
	projectType string,
	requirements advisor.ProjectRequirements,
) *BDDTestScenario {
	return &BDDTestScenario{
		Name:         scenarioName,
		ProjectType:  projectType,
		Requirements: requirements,
		Helpers:      h,
	}
}

// BDDTestScenario represents a complete test scenario
type BDDTestScenario struct {
	Name         string
	ProjectType  string
	Requirements advisor.ProjectRequirements
	Helpers      *BDDTestHelpers
}

// Execute runs the scenario with the given context
func (s *BDDTestScenario) Execute(t *testing.T, ctx *BDDTestContext) error {
	// Set up the scenario
	ctx.projectType = s.ProjectType
	ctx.requirements = s.Requirements
	
	// Execute recommendation
	err := ctx.IRequestQuickModeRecommendations()
	if err != nil {
		return fmt.Errorf("scenario '%s' failed to get recommendation: %w", s.Name, err)
	}
	
	// Validate the recommendation
	if err := s.Helpers.ValidateRecommendationStructure(t, ctx.recommendation); err != nil {
		return fmt.Errorf("scenario '%s' recommendation validation failed: %w", s.Name, err)
	}
	
	return nil
}

// CommonScenarios provides pre-built common test scenarios
func (h *BDDTestHelpers) CommonScenarios() []*BDDTestScenario {
	return []*BDDTestScenario{
		h.CreateTestScenario(
			"Standard E-commerce API",
			"api",
			h.fixtures.StandardRequirements,
		),
		h.CreateTestScenario(
			"Simple CLI Tool",
			"cli",
			advisor.ProjectRequirements{
				Domain:         "devtools",
				TeamExperience: "junior",
				PreferredStyle: "simple",
			},
		),
		h.CreateTestScenario(
			"Complex Fintech API",
			"api",
			h.fixtures.ComplexRequirements,
		),
		h.CreateTestScenario(
			"Minimal Requirements",
			"api",
			h.fixtures.MinimalRequirements,
		),
	}
}

// BDDAssertions provides fluent assertion interface
type BDDAssertions struct {
	t   *testing.T
	rec *advisor.ArchitectureRecommendation
	err error
}

// NewBDDAssertions creates a new assertion helper
func NewBDDAssertions(t *testing.T, rec *advisor.ArchitectureRecommendation) *BDDAssertions {
	return &BDDAssertions{
		t:   t,
		rec: rec,
	}
}

// ShouldHaveBlueprint asserts blueprint matches expectation
func (a *BDDAssertions) ShouldHaveBlueprint(expected string) *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	if a.rec.Blueprint != expected {
		a.err = fmt.Errorf("expected blueprint '%s', got '%s'", expected, a.rec.Blueprint)
	}
	return a
}

// ShouldHaveBlueprintIn asserts blueprint is in expected list
func (a *BDDAssertions) ShouldHaveBlueprintIn(expected []string) *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	for _, exp := range expected {
		if a.rec.Blueprint == exp {
			return a
		}
	}
	
	a.err = fmt.Errorf("blueprint '%s' not in expected list %v", a.rec.Blueprint, expected)
	return a
}

// ShouldHaveConfidenceAtLeast asserts minimum confidence level
func (a *BDDAssertions) ShouldHaveConfidenceAtLeast(min float64) *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	if a.rec.Confidence < min {
		a.err = fmt.Errorf("confidence %f below minimum %f", a.rec.Confidence, min)
	}
	return a
}

// ShouldHaveReasoning asserts reasoning is provided
func (a *BDDAssertions) ShouldHaveReasoning() *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	if len(a.rec.Reasoning) == 0 {
		a.err = fmt.Errorf("no reasoning provided")
	}
	return a
}

// ShouldHaveAlternatives asserts alternatives are provided
func (a *BDDAssertions) ShouldHaveAlternatives() *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	if len(a.rec.Alternatives) == 0 {
		a.err = fmt.Errorf("no alternatives provided")
	}
	return a
}

// ShouldHaveFileCountBetween asserts file count range
func (a *BDDAssertions) ShouldHaveFileCountBetween(min, max int) *BDDAssertions {
	if a.err != nil {
		return a
	}
	
	if a.rec.EstimatedFiles < min || a.rec.EstimatedFiles > max {
		a.err = fmt.Errorf("file count %d not between %d and %d", 
			a.rec.EstimatedFiles, min, max)
	}
	return a
}

// Assert performs the final assertion check
func (a *BDDAssertions) Assert() error {
	if a.err != nil {
		a.t.Error(a.err)
		return a.err
	}
	return nil
}

// Usage example for fluent assertions:
/*
err := NewBDDAssertions(t, recommendation).
	ShouldHaveBlueprintIn([]string{"web-api", "web-api-clean"}).
	ShouldHaveConfidenceAtLeast(0.6).
	ShouldHaveReasoning().
	ShouldHaveAlternatives().
	ShouldHaveFileCountBetween(20, 100).
	Assert()
*/