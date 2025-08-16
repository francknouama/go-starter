package advisor

import (
	"testing"

	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchitectureAdvisor_AnalyzeRequirements(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	tests := []struct {
		name     string
		req      ProjectRequirements
		expected string // expected blueprint
		minConf  float64 // minimum confidence
	}{
		{
			name: "CLI tool for devtools domain",
			req: ProjectRequirements{
				Domain:           "devtools",
				TeamExperience:   "junior",
				TimeToMarket:     "quick",
				ExpectedLoad:     "low",
				TeamSize:         1,
				IsFirstGoProject: true,
				PreferredStyle:   "simple",
				DeploymentTarget: "local", // Hint for CLI vs library
			},
			expected: "library", // Library is actually good for devtools domain
			minConf:  0.6,
		},
		{
			name: "E-commerce API for experienced team",
			req: ProjectRequirements{
				Domain:         "e-commerce",
				TeamExperience: "senior",
				TimeToMarket:   "standard",
				ExpectedLoad:   "high",
				TeamSize:       5,
				DatabaseRequirements: "complex",
				AuthRequirements: "oauth",
			},
			expected: "web-api-ddd",
			minConf:  0.7,
		},
		{
			name: "Fintech API with strict requirements",
			req: ProjectRequirements{
				Domain:         "fintech",
				TeamExperience: "expert",
				TimeToMarket:   "thorough",
				ExpectedLoad:   "high",
				TeamSize:       8,
				ComplianceNeeds: []string{"sox", "pci"},
				DatabaseRequirements: "complex",
				MonitoringNeeds: "enterprise",
			},
			expected: "web-api-ddd",
			minConf:  0.8,
		},
		{
			name: "IoT event processing",
			req: ProjectRequirements{
				Domain:         "iot",
				TeamExperience: "mixed",
				TimeToMarket:   "standard",
				ExpectedLoad:   "high",
				ResponseTime:   "fast",
				DeploymentTarget: "cloud",
			},
			expected: "lambda-event-processing",
			minConf:  0.6,
		},
		{
			name: "MVP prototype with urgent timeline",
			req: ProjectRequirements{
				Domain:         "e-commerce",
				TeamExperience: "mixed",
				TimeToMarket:   "mvp",
				ExpectedLoad:   "low",
				TeamSize:       2,
				Budget:         "startup",
			},
			expected: "web-api",
			minConf:  0.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recommendation, err := advisor.AnalyzeRequirements(tt.req)
			
			require.NoError(t, err, "AnalyzeRequirements should not return error")
			require.NotNil(t, recommendation, "Recommendation should not be nil")
			
			assert.Equal(t, tt.expected, recommendation.Blueprint, 
				"Expected blueprint %s, got %s", tt.expected, recommendation.Blueprint)
			assert.GreaterOrEqual(t, recommendation.Confidence, tt.minConf,
				"Confidence should be at least %.2f, got %.2f", tt.minConf, recommendation.Confidence)
			
			// Verify recommendation structure
			assert.NotEmpty(t, recommendation.Framework, "Framework should be recommended")
			assert.NotEmpty(t, recommendation.Logger, "Logger should be recommended")
			assert.Greater(t, recommendation.EstimatedFiles, 0, "Should estimate file count")
			assert.NotEmpty(t, recommendation.DevelopmentTime, "Should estimate development time")
			
			// Verify reasoning is provided
			assert.NotEmpty(t, recommendation.Reasoning, "Should provide reasoning")
			
			// Verify alternatives are provided
			assert.LessOrEqual(t, len(recommendation.Alternatives), 3, "Should provide at most 3 alternatives")
		})
	}
}

func TestArchitectureAdvisor_ScoreCalculation(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	req := ProjectRequirements{
		Domain:         "fintech",
		TeamExperience: "senior",
		ExpectedLoad:   "high",
		ResponseTime:   "fast",
		TeamSize:       5,
	}

	scores := advisor.scoreArchitectures(req)
	
	// Should score all available patterns
	assert.Greater(t, len(scores), 5, "Should score multiple patterns")
	
	// Find DDD score (should be high for fintech)
	dddScore, exists := scores["web-api-ddd"]
	assert.True(t, exists, "Should score web-api-ddd pattern")
	assert.Greater(t, dddScore.TotalScore, 0.5, "DDD should score well for fintech")
	
	// Simple CLI should score lower for complex fintech requirements
	cliScore, exists := scores["cli-simple"]
	assert.True(t, exists, "Should score cli-simple pattern")
	assert.Less(t, cliScore.TotalScore, dddScore.TotalScore, "CLI should score lower than DDD for fintech")
}

func TestArchitectureAdvisor_ComplexityDetermination(t *testing.T) {
	advisor := NewArchitectureAdvisor()
	pattern := &ArchitecturePattern{
		ID:              "web-api-clean",
		ComplexityScore: 0.6,
	}

	tests := []struct {
		name     string
		req      ProjectRequirements
		expected interfaces.ComplexityLevel
	}{
		{
			name: "Junior team should get simple complexity",
			req: ProjectRequirements{
				TeamExperience:   "junior",
				IsFirstGoProject: true,
			},
			expected: interfaces.ComplexitySimple,
		},
		{
			name: "Expert team should get advanced complexity",
			req: ProjectRequirements{
				TeamExperience: "expert",
				ExpectedLoad:   "high",
			},
			expected: interfaces.ComplexityAdvanced,
		},
		{
			name: "MVP timeline should reduce complexity",
			req: ProjectRequirements{
				TeamExperience: "senior",
				TimeToMarket:   "mvp",
				Budget:         "startup",
			},
			expected: interfaces.ComplexitySimple,
		},
		{
			name: "High load should increase complexity",
			req: ProjectRequirements{
				TeamExperience: "mixed",
				ExpectedLoad:   "massive",
			},
			expected: interfaces.ComplexityAdvanced,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			complexity := advisor.determineComplexity(tt.req, pattern)
			assert.Equal(t, tt.expected, complexity, 
				"Expected complexity %s, got %s", tt.expected.String(), complexity.String())
		})
	}
}

func TestArchitectureAdvisor_FrameworkRecommendation(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	tests := []struct {
		name     string
		req      ProjectRequirements
		pattern  *ArchitecturePattern
		expected string
	}{
		{
			name: "CLI blueprint should recommend cobra",
			req:  ProjectRequirements{},
			pattern: &ArchitecturePattern{
				Blueprint: "cli",
			},
			expected: "cobra",
		},
		{
			name: "High performance API should recommend fiber",
			req: ProjectRequirements{
				ResponseTime: "realtime",
			},
			pattern: &ArchitecturePattern{
				Blueprint: "web-api",
			},
			expected: "fiber",
		},
		{
			name: "Standard API should recommend gin",
			req: ProjectRequirements{
				ResponseTime: "standard",
			},
			pattern: &ArchitecturePattern{
				Blueprint: "web-api",
			},
			expected: "gin",
		},
		{
			name: "Default should recommend echo",
			req: ProjectRequirements{},
			pattern: &ArchitecturePattern{
				Blueprint: "web-api",
			},
			expected: "echo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			framework := advisor.recommendFramework(tt.req, tt.pattern)
			assert.Equal(t, tt.expected, framework, 
				"Expected framework %s, got %s", tt.expected, framework)
		})
	}
}

func TestArchitectureAdvisor_LoggerRecommendation(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	tests := []struct {
		name     string
		req      ProjectRequirements
		expected string
	}{
		{
			name: "High performance should recommend zap",
			req: ProjectRequirements{
				ResponseTime:  "realtime",
				ExpectedLoad: "massive",
			},
			expected: "zap",
		},
		{
			name: "Enterprise compliance should recommend logrus",
			req: ProjectRequirements{
				ComplianceNeeds: []string{"sox", "hipaa"},
				MonitoringNeeds: "enterprise",
			},
			expected: "logrus",
		},
		{
			name: "Default should recommend slog",
			req:  ProjectRequirements{},
			expected: "slog",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := advisor.recommendLogger(tt.req, &ArchitecturePattern{})
			assert.Equal(t, tt.expected, logger, 
				"Expected logger %s, got %s", tt.expected, logger)
		})
	}
}

func TestArchitectureAdvisor_DatabaseRecommendation(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	tests := []struct {
		name     string
		req      ProjectRequirements
		expected string
	}{
		{
			name: "Simple requirements should recommend sqlite",
			req: ProjectRequirements{
				DatabaseRequirements: "simple",
			},
			expected: "sqlite",
		},
		{
			name: "Complex requirements should recommend postgres",
			req: ProjectRequirements{
				DatabaseRequirements: "complex",
			},
			expected: "postgres",
		},
		{
			name: "Distributed requirements should recommend postgres",
			req: ProjectRequirements{
				DatabaseRequirements: "distributed",
			},
			expected: "postgres",
		},
		{
			name: "Default should recommend postgres",
			req:  ProjectRequirements{},
			expected: "postgres",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			database := advisor.recommendDatabase(tt.req, &ArchitecturePattern{})
			assert.Equal(t, tt.expected, database, 
				"Expected database %s, got %s", tt.expected, database)
		})
	}
}

func TestArchitectureAdvisor_DevelopmentTimeEstimation(t *testing.T) {
	advisor := NewArchitectureAdvisor()

	tests := []struct {
		name     string
		req      ProjectRequirements
		pattern  *ArchitecturePattern
		contains string // substring that should be in the result
	}{
		{
			name: "Small project with junior team",
			req: ProjectRequirements{
				TeamExperience: "junior",
			},
			pattern: &ArchitecturePattern{
				EstimatedFiles: 8,
			},
			contains: "day", // Should be in days
		},
		{
			name: "Large project with expert team",
			req: ProjectRequirements{
				TeamExperience: "expert",
			},
			pattern: &ArchitecturePattern{
				EstimatedFiles: 80,
			},
			contains: "week", // Should be in weeks
		},
		{
			name: "MVP timeline should reduce time",
			req: ProjectRequirements{
				TeamExperience: "mixed",
				TimeToMarket:   "mvp",
			},
			pattern: &ArchitecturePattern{
				EstimatedFiles: 15, // Smaller project for MVP
			},
			contains: "day", // Should be shorter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			timeEst := advisor.estimateDevelopmentTime(tt.req, tt.pattern)
			assert.Contains(t, timeEst, tt.contains, 
				"Expected time estimate to contain '%s', got '%s'", tt.contains, timeEst)
			assert.NotEmpty(t, timeEst, "Time estimate should not be empty")
		})
	}
}

// Test helper functions
func TestHelperFunctions(t *testing.T) {
	// Test abs function
	assert.Equal(t, 5.0, abs(-5.0))
	assert.Equal(t, 3.0, abs(3.0))
	
	// Test clamp function
	assert.Equal(t, 0.5, clamp(0.3, 0.5, 1.0)) // Below min
	assert.Equal(t, 1.0, clamp(1.5, 0.5, 1.0)) // Above max
	assert.Equal(t, 0.7, clamp(0.7, 0.5, 1.0)) // Within range
}

// Benchmark tests
func BenchmarkArchitectureAdvisor_AnalyzeRequirements(b *testing.B) {
	advisor := NewArchitectureAdvisor()
	req := ProjectRequirements{
		Domain:         "e-commerce",
		TeamExperience: "mixed",
		ExpectedLoad:   "medium",
		TeamSize:       3,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := advisor.AnalyzeRequirements(req)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkArchitectureAdvisor_ScoreCalculation(b *testing.B) {
	advisor := NewArchitectureAdvisor()
	req := ProjectRequirements{
		Domain:         "fintech",
		TeamExperience: "senior",
		ExpectedLoad:   "high",
		TeamSize:       5,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = advisor.scoreArchitectures(req)
	}
}