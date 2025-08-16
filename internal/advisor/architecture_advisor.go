package advisor

import (
	"fmt"
	"sort"
	"strings"

	"github.com/francknouama/go-starter/internal/prompts/interfaces"
)

// ProjectRequirements represents the requirements gathered from the user
type ProjectRequirements struct {
	// Basic project info
	Name        string
	Description string
	Domain      string // e.g., "e-commerce", "fintech", "iot", "content-management"
	
	// Technical requirements
	ExpectedLoad       string   // "low", "medium", "high", "massive"
	ConcurrentUsers    int      // Expected concurrent users
	DataVolume         string   // "small", "medium", "large", "big-data"
	ResponseTime       string   // "relaxed", "standard", "fast", "realtime"
	Availability       string   // "standard", "high", "critical"
	
	// Business requirements
	TimeToMarket       string   // "mvp", "quick", "standard", "thorough"
	Budget             string   // "startup", "small", "medium", "enterprise"
	TeamSize           int      // Number of developers
	TeamExperience     string   // "junior", "mixed", "senior", "expert"
	
	// Integration requirements
	ThirdPartyServices []string // APIs, payment processors, etc.
	DatabaseRequirements string // "simple", "complex", "distributed"
	AuthRequirements   string   // "none", "basic", "oauth", "enterprise"
	
	// Operational requirements
	DeploymentTarget   string   // "local", "cloud", "hybrid", "edge"
	ComplianceNeeds    []string // "gdpr", "hipaa", "sox", "pci"
	MonitoringNeeds    string   // "basic", "standard", "advanced", "enterprise"
	
	// User indicators
	IsFirstGoProject   bool
	HasMicroserviceExp bool
	PreferredStyle     string // "simple", "structured", "enterprise"
}

// ArchitectureRecommendation represents a recommended architecture
type ArchitectureRecommendation struct {
	Blueprint    string                      `json:"blueprint"`
	Architecture string                      `json:"architecture"`
	Complexity   interfaces.ComplexityLevel  `json:"complexity"`
	Confidence   float64                     `json:"confidence"` // 0.0 to 1.0
	Reasoning    []string                    `json:"reasoning"`
	Pros         []string                    `json:"pros"`
	Cons         []string                    `json:"cons"`
	Alternatives []AlternativeRecommendation `json:"alternatives"`
	
	// Specific recommendations
	Framework        string            `json:"framework"`
	Logger           string            `json:"logger"`
	Database         string            `json:"database"`
	Features         map[string]string `json:"features"`
	EstimatedFiles   int              `json:"estimated_files"`
	DevelopmentTime  string           `json:"development_time"`
}

// AlternativeRecommendation represents an alternative architecture choice
type AlternativeRecommendation struct {
	Blueprint  string  `json:"blueprint"`
	Confidence float64 `json:"confidence"`
	Reason     string  `json:"reason"`
}

// ArchitectureAdvisor provides intelligent architecture recommendations
type ArchitectureAdvisor struct {
	knowledgeBase *ArchitectureKnowledgeBase
}

// NewArchitectureAdvisor creates a new architecture advisor
func NewArchitectureAdvisor() *ArchitectureAdvisor {
	return &ArchitectureAdvisor{
		knowledgeBase: NewArchitectureKnowledgeBase(),
	}
}

// AnalyzeRequirements analyzes project requirements and provides architecture recommendations
func (a *ArchitectureAdvisor) AnalyzeRequirements(req ProjectRequirements) (*ArchitectureRecommendation, error) {
	// Score all available architectures
	scores := a.scoreArchitectures(req)
	
	// Get the best recommendation
	bestMatch := a.getBestMatch(scores)
	if bestMatch == nil {
		return nil, fmt.Errorf("no suitable architecture found for requirements")
	}
	
	// Generate detailed recommendation
	recommendation := a.generateRecommendation(req, bestMatch, scores)
	
	return recommendation, nil
}

// scoreArchitectures scores all available architectures against the requirements
func (a *ArchitectureAdvisor) scoreArchitectures(req ProjectRequirements) map[string]*ArchitectureScore {
	scores := make(map[string]*ArchitectureScore)
	
	for _, pattern := range a.knowledgeBase.GetAvailablePatterns() {
		score := a.calculateScore(req, pattern)
		scores[pattern.ID] = score
	}
	
	return scores
}

// calculateScore calculates a score for a specific architecture pattern
func (a *ArchitectureAdvisor) calculateScore(req ProjectRequirements, pattern *ArchitecturePattern) *ArchitectureScore {
	score := &ArchitectureScore{
		Pattern:    pattern,
		TotalScore: 0.0,
		Factors:    make(map[string]float64),
	}
	
	// Performance requirements scoring
	performanceScore := a.scorePerformanceRequirements(req, pattern)
	score.Factors["performance"] = performanceScore
	
	// Complexity requirements scoring
	complexityScore := a.scoreComplexityRequirements(req, pattern)
	score.Factors["complexity"] = complexityScore
	
	// Team experience scoring
	teamScore := a.scoreTeamExperience(req, pattern)
	score.Factors["team"] = teamScore
	
	// Time to market scoring
	timeScore := a.scoreTimeToMarket(req, pattern)
	score.Factors["time"] = timeScore
	
	// Domain fit scoring
	domainScore := a.scoreDomainFit(req, pattern)
	score.Factors["domain"] = domainScore
	
	// Scalability scoring
	scalabilityScore := a.scoreScalabilityRequirements(req, pattern)
	score.Factors["scalability"] = scalabilityScore
	
	// Calculate weighted total score
	weights := map[string]float64{
		"performance":   0.20,
		"complexity":    0.25,
		"team":         0.15,
		"time":         0.15,
		"domain":       0.15,
		"scalability":  0.10,
	}
	
	for factor, weight := range weights {
		score.TotalScore += score.Factors[factor] * weight
	}
	
	return score
}

// scorePerformanceRequirements scores based on performance requirements
func (a *ArchitectureAdvisor) scorePerformanceRequirements(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	performanceMap := map[string]float64{
		"relaxed":  0.2,
		"standard": 0.5,
		"fast":     0.8,
		"realtime": 1.0,
	}
	
	requiredPerf := performanceMap[req.ResponseTime]
	if requiredPerf == 0 {
		requiredPerf = 0.5 // default
	}
	
	// Pattern performance capabilities
	patternPerf := pattern.PerformanceScore
	
	// If pattern can handle the requirement, score well
	if patternPerf >= requiredPerf {
		return 1.0 - (patternPerf-requiredPerf)*0.1 // Slight penalty for overkill
	}
	
	// If pattern can't handle requirement, score poorly
	return patternPerf / requiredPerf * 0.5
}

// scoreComplexityRequirements scores based on complexity requirements
func (a *ArchitectureAdvisor) scoreComplexityRequirements(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	// Match pattern complexity with team experience and project needs
	teamExperienceMap := map[string]float64{
		"junior": 0.2,
		"mixed":  0.5,
		"senior": 0.8,
		"expert": 1.0,
	}
	
	teamLevel := teamExperienceMap[req.TeamExperience]
	if teamLevel == 0 {
		teamLevel = 0.5
	}
	
	// If pattern complexity matches team capability
	complexityDiff := abs(pattern.ComplexityScore - teamLevel)
	return 1.0 - complexityDiff
}

// scoreTeamExperience scores based on team experience with the pattern
func (a *ArchitectureAdvisor) scoreTeamExperience(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	baseScore := 0.5
	
	// Bonus for Go experience with simpler patterns
	if !req.IsFirstGoProject {
		baseScore += 0.2
	}
	
	// Bonus for microservice experience with complex patterns
	if req.HasMicroserviceExp && pattern.IsDistributed {
		baseScore += 0.3
	}
	
	// Penalty for complex patterns with inexperienced teams
	if req.TeamExperience == "junior" && pattern.ComplexityScore > 0.7 {
		baseScore -= 0.4
	}
	
	return clamp(baseScore, 0.0, 1.0)
}

// scoreTimeToMarket scores based on time to market requirements
func (a *ArchitectureAdvisor) scoreTimeToMarket(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	timeMap := map[string]float64{
		"mvp":      1.0, // Need fastest
		"quick":    0.8,
		"standard": 0.5,
		"thorough": 0.2, // Can take time
	}
	
	urgency := timeMap[req.TimeToMarket]
	if urgency == 0 {
		urgency = 0.5
	}
	
	// Simpler patterns score better for urgent projects
	return 1.0 - (pattern.ComplexityScore * urgency)
}

// scoreDomainFit scores based on domain-specific requirements
func (a *ArchitectureAdvisor) scoreDomainFit(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	// Domain-specific preferences
	domainPreferences := map[string]map[string]float64{
		"e-commerce": {
			"cli":                0.1,
			"web-api":           0.9,
			"web-api-clean":     0.8,
			"web-api-ddd":       0.9,
			"web-api-hexagonal": 0.8,
			"microservice":      0.9,
			"monolith":          0.7,
			"lambda":            0.6,
		},
		"fintech": {
			"cli":                0.2,
			"web-api":           0.7,
			"web-api-clean":     0.9,
			"web-api-ddd":       1.0,
			"web-api-hexagonal": 0.9,
			"microservice":      0.8,
			"monolith":          0.6,
			"lambda":            0.5,
		},
		"iot": {
			"cli":                0.8,
			"web-api":           0.6,
			"lambda":            0.9,
			"microservice":      0.8,
			"event-driven":      1.0,
		},
		"content-management": {
			"web-api":           0.9,
			"web-api-clean":     0.8,
			"monolith":          0.9,
			"microservice":      0.7,
		},
	}
	
	if prefs, exists := domainPreferences[req.Domain]; exists {
		if score, exists := prefs[pattern.ID]; exists {
			return score
		}
	}
	
	return 0.5 // neutral score for unknown domains
}

// scoreScalabilityRequirements scores based on scalability needs
func (a *ArchitectureAdvisor) scoreScalabilityRequirements(req ProjectRequirements, pattern *ArchitecturePattern) float64 {
	loadMap := map[string]float64{
		"low":     0.2,
		"medium":  0.5,
		"high":    0.8,
		"massive": 1.0,
	}
	
	requiredScale := loadMap[req.ExpectedLoad]
	if requiredScale == 0 {
		requiredScale = 0.5
	}
	
	// Pattern scalability capabilities
	patternScale := pattern.ScalabilityScore
	
	// Similar logic to performance scoring
	if patternScale >= requiredScale {
		return 1.0 - (patternScale-requiredScale)*0.1
	}
	
	return patternScale / requiredScale * 0.6
}

// getBestMatch finds the highest scoring architecture
func (a *ArchitectureAdvisor) getBestMatch(scores map[string]*ArchitectureScore) *ArchitectureScore {
	var best *ArchitectureScore
	bestScore := 0.0
	
	for _, score := range scores {
		if score.TotalScore > bestScore {
			bestScore = score.TotalScore
			best = score
		}
	}
	
	return best
}

// generateRecommendation creates a detailed recommendation
func (a *ArchitectureAdvisor) generateRecommendation(req ProjectRequirements, best *ArchitectureScore, allScores map[string]*ArchitectureScore) *ArchitectureRecommendation {
	pattern := best.Pattern
	
	recommendation := &ArchitectureRecommendation{
		Blueprint:    pattern.Blueprint,
		Architecture: pattern.Architecture,
		Complexity:   a.determineComplexity(req, pattern),
		Confidence:   best.TotalScore,
		Framework:    a.recommendFramework(req, pattern),
		Logger:       a.recommendLogger(req, pattern),
		Database:     a.recommendDatabase(req, pattern),
		Features:     a.recommendFeatures(req, pattern),
		EstimatedFiles: pattern.EstimatedFiles,
		DevelopmentTime: a.estimateDevelopmentTime(req, pattern),
	}
	
	// Generate reasoning
	recommendation.Reasoning = a.generateReasoning(req, pattern, best.Factors)
	recommendation.Pros = pattern.Pros
	recommendation.Cons = pattern.Cons
	
	// Generate alternatives
	recommendation.Alternatives = a.generateAlternatives(allScores, best.Pattern.ID)
	
	return recommendation
}

// determineComplexity determines the appropriate complexity level
func (a *ArchitectureAdvisor) determineComplexity(req ProjectRequirements, pattern *ArchitecturePattern) interfaces.ComplexityLevel {
	// Start with team experience
	baseComplexity := interfaces.ComplexityStandard
	
	if req.TeamExperience == "junior" || req.IsFirstGoProject {
		baseComplexity = interfaces.ComplexitySimple
	} else if req.TeamExperience == "expert" {
		baseComplexity = interfaces.ComplexityAdvanced
	}
	
	// Adjust based on project requirements
	if req.TimeToMarket == "mvp" || req.Budget == "startup" {
		if baseComplexity > interfaces.ComplexitySimple {
			baseComplexity--
		}
	}
	
	if req.ExpectedLoad == "high" || req.ExpectedLoad == "massive" {
		if baseComplexity < interfaces.ComplexityAdvanced {
			baseComplexity++
		}
	}
	
	return baseComplexity
}

// Helper functions
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Additional recommendation methods
func (a *ArchitectureAdvisor) recommendFramework(req ProjectRequirements, pattern *ArchitecturePattern) string {
	if pattern.Blueprint == "cli" {
		return "cobra"
	}
	
	// Web frameworks based on performance needs
	switch req.ResponseTime {
	case "realtime", "fast":
		return "fiber" // Fastest
	case "standard":
		return "gin"   // Balanced
	default:
		return "echo"  // Simple and reliable
	}
}

func (a *ArchitectureAdvisor) recommendLogger(req ProjectRequirements, pattern *ArchitecturePattern) string {
	// High performance needs
	if req.ResponseTime == "realtime" || req.ExpectedLoad == "massive" {
		return "zap"
	}
	
	// Enterprise/compliance needs
	if len(req.ComplianceNeeds) > 0 || req.MonitoringNeeds == "enterprise" {
		return "logrus"
	}
	
	// Default to standard library
	return "slog"
}

func (a *ArchitectureAdvisor) recommendDatabase(req ProjectRequirements, pattern *ArchitecturePattern) string {
	switch req.DatabaseRequirements {
	case "simple":
		return "sqlite"
	case "complex":
		return "postgres"
	case "distributed":
		return "postgres" // With potential for sharding
	default:
		return "postgres" // Safe default
	}
}

func (a *ArchitectureAdvisor) recommendFeatures(req ProjectRequirements, pattern *ArchitecturePattern) map[string]string {
	features := make(map[string]string)
	
	// Authentication
	if req.AuthRequirements != "none" {
		features["auth"] = req.AuthRequirements
	}
	
	// Monitoring
	if req.MonitoringNeeds != "basic" {
		features["monitoring"] = req.MonitoringNeeds
	}
	
	// Compliance
	if len(req.ComplianceNeeds) > 0 {
		features["compliance"] = strings.Join(req.ComplianceNeeds, ",")
	}
	
	return features
}

func (a *ArchitectureAdvisor) estimateDevelopmentTime(req ProjectRequirements, pattern *ArchitecturePattern) string {
	baseHours := pattern.EstimatedFiles * 2 // 2 hours per file as baseline
	
	// Adjust for team experience
	switch req.TeamExperience {
	case "junior":
		baseHours = int(float64(baseHours) * 1.5)
	case "expert":
		baseHours = int(float64(baseHours) * 0.7)
	}
	
	// Adjust for complexity
	if req.TimeToMarket == "mvp" {
		baseHours = int(float64(baseHours) * 0.8)
	}
	
	days := baseHours / 8
	if days < 1 {
		return "< 1 day"
	} else if days < 5 {
		return fmt.Sprintf("%d days", days)
	} else if days < 20 {
		return fmt.Sprintf("%.1f weeks", float64(days)/5)
	} else {
		return fmt.Sprintf("%.1f months", float64(days)/20)
	}
}

func (a *ArchitectureAdvisor) generateReasoning(req ProjectRequirements, pattern *ArchitecturePattern, factors map[string]float64) []string {
	reasoning := []string{}
	
	// Performance reasoning
	if factors["performance"] > 0.8 {
		reasoning = append(reasoning, fmt.Sprintf("Excellent performance match for %s response time requirements", req.ResponseTime))
	} else if factors["performance"] < 0.4 {
		reasoning = append(reasoning, fmt.Sprintf("May need optimization for %s response time requirements", req.ResponseTime))
	}
	
	// Team experience reasoning
	if factors["team"] > 0.7 {
		reasoning = append(reasoning, fmt.Sprintf("Well-suited for %s team experience level", req.TeamExperience))
	}
	
	// Time to market reasoning
	if factors["time"] > 0.7 {
		reasoning = append(reasoning, fmt.Sprintf("Good fit for %s time to market", req.TimeToMarket))
	}
	
	// Domain reasoning
	if factors["domain"] > 0.8 {
		reasoning = append(reasoning, fmt.Sprintf("Excellent match for %s domain requirements", req.Domain))
	}
	
	return reasoning
}

func (a *ArchitectureAdvisor) generateAlternatives(allScores map[string]*ArchitectureScore, excludeID string) []AlternativeRecommendation {
	type scoreEntry struct {
		id    string
		score *ArchitectureScore
	}
	
	var entries []scoreEntry
	for id, score := range allScores {
		if id != excludeID {
			entries = append(entries, scoreEntry{id, score})
		}
	}
	
	// Sort by score
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].score.TotalScore > entries[j].score.TotalScore
	})
	
	// Take top 3 alternatives
	alternatives := []AlternativeRecommendation{}
	for i := 0; i < 3 && i < len(entries); i++ {
		entry := entries[i]
		alternatives = append(alternatives, AlternativeRecommendation{
			Blueprint:  entry.score.Pattern.Blueprint,
			Confidence: entry.score.TotalScore,
			Reason:     fmt.Sprintf("Alternative with %.1f%% confidence", entry.score.TotalScore*100),
		})
	}
	
	return alternatives
}