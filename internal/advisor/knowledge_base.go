package advisor

// ArchitecturePattern represents a blueprint/architecture combination with its characteristics
type ArchitecturePattern struct {
	ID                string   `json:"id"`
	Blueprint         string   `json:"blueprint"`
	Architecture      string   `json:"architecture"`
	DisplayName       string   `json:"display_name"`
	Description       string   `json:"description"`
	
	// Scoring factors (0.0 to 1.0)
	PerformanceScore  float64  `json:"performance_score"`
	ComplexityScore   float64  `json:"complexity_score"`
	ScalabilityScore  float64  `json:"scalability_score"`
	MaintainabilityScore float64 `json:"maintainability_score"`
	
	// Characteristics
	IsDistributed     bool     `json:"is_distributed"`
	EstimatedFiles    int      `json:"estimated_files"`
	LearningCurve     string   `json:"learning_curve"` // "low", "medium", "high"
	
	// Pros and cons
	Pros              []string `json:"pros"`
	Cons              []string `json:"cons"`
	
	// Use cases
	BestFor           []string `json:"best_for"`
	AvoidIf           []string `json:"avoid_if"`
	
	// Requirements
	MinTeamSize       int      `json:"min_team_size"`
	RecommendedTeamSize int    `json:"recommended_team_size"`
}

// ArchitectureScore represents the scoring result for a pattern
type ArchitectureScore struct {
	Pattern    *ArchitecturePattern
	TotalScore float64
	Factors    map[string]float64
}

// ArchitectureKnowledgeBase contains expert knowledge about all available patterns
type ArchitectureKnowledgeBase struct {
	patterns map[string]*ArchitecturePattern
}

// NewArchitectureKnowledgeBase creates a new knowledge base with all patterns
func NewArchitectureKnowledgeBase() *ArchitectureKnowledgeBase {
	kb := &ArchitectureKnowledgeBase{
		patterns: make(map[string]*ArchitecturePattern),
	}
	kb.initializePatterns()
	return kb
}

// GetAvailablePatterns returns all available architecture patterns
func (kb *ArchitectureKnowledgeBase) GetAvailablePatterns() []*ArchitecturePattern {
	patterns := make([]*ArchitecturePattern, 0, len(kb.patterns))
	for _, pattern := range kb.patterns {
		patterns = append(patterns, pattern)
	}
	return patterns
}

// GetPattern returns a specific pattern by ID
func (kb *ArchitectureKnowledgeBase) GetPattern(id string) *ArchitecturePattern {
	return kb.patterns[id]
}

// initializePatterns populates the knowledge base with all architecture patterns
func (kb *ArchitectureKnowledgeBase) initializePatterns() {
	// CLI Patterns
	kb.addPattern(&ArchitecturePattern{
		ID:           "cli-simple",
		Blueprint:    "cli-simple",
		Architecture: "simple",
		DisplayName:  "Simple CLI",
		Description:  "Minimal command-line tool with basic functionality",
		
		PerformanceScore:     0.6,
		ComplexityScore:      0.1, // Very simple
		ScalabilityScore:     0.3,
		MaintainabilityScore: 0.7,
		
		IsDistributed:    false,
		EstimatedFiles:   8,
		LearningCurve:    "low",
		
		Pros: []string{
			"Quick to develop and deploy",
			"Minimal dependencies",
			"Easy to understand and maintain",
			"Perfect for small utilities",
		},
		Cons: []string{
			"Limited functionality",
			"Not suitable for complex workflows",
			"May need refactoring as requirements grow",
		},
		
		BestFor: []string{
			"Quick utilities and scripts",
			"Learning Go fundamentals",
			"Prototyping CLI ideas",
			"Single-purpose tools",
		},
		AvoidIf: []string{
			"Need complex command structures",
			"Require configuration management",
			"Planning multi-team development",
		},
		
		MinTeamSize:         1,
		RecommendedTeamSize: 1,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "cli",
		Blueprint:    "cli",
		Architecture: "standard",
		DisplayName:  "Standard CLI",
		Description:  "Production-ready command-line application with full features",
		
		PerformanceScore:     0.8,
		ComplexityScore:      0.4,
		ScalabilityScore:     0.6,
		MaintainabilityScore: 0.9,
		
		IsDistributed:    false,
		EstimatedFiles:   29,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Comprehensive CLI framework",
			"Excellent for complex commands",
			"Professional structure",
			"Easy to extend and maintain",
		},
		Cons: []string{
			"More files than needed for simple tools",
			"Steeper learning curve",
			"Overhead for basic utilities",
		},
		
		BestFor: []string{
			"Production CLI tools",
			"Multi-command applications",
			"Team development",
			"Long-term maintenance",
		},
		AvoidIf: []string{
			"One-off scripts",
			"Simple utilities",
			"Learning projects only",
		},
		
		MinTeamSize:         1,
		RecommendedTeamSize: 2,
	})
	
	// Web API Patterns
	kb.addPattern(&ArchitecturePattern{
		ID:           "web-api",
		Blueprint:    "web-api",
		Architecture: "standard",
		DisplayName:  "Standard Web API",
		Description:  "Traditional layered web API for rapid development",
		
		PerformanceScore:     0.7,
		ComplexityScore:      0.3,
		ScalabilityScore:     0.6,
		MaintainabilityScore: 0.7,
		
		IsDistributed:    false,
		EstimatedFiles:   35,
		LearningCurve:    "low",
		
		Pros: []string{
			"Quick to develop",
			"Familiar structure",
			"Good for CRUD operations",
			"Minimal learning curve",
		},
		Cons: []string{
			"Limited separation of concerns",
			"Can become monolithic",
			"Harder to test business logic",
		},
		
		BestFor: []string{
			"MVPs and prototypes",
			"Simple CRUD APIs",
			"Small to medium projects",
			"Rapid development needs",
		},
		AvoidIf: []string{
			"Complex business logic",
			"Multiple teams",
			"Long-term enterprise use",
		},
		
		MinTeamSize:         1,
		RecommendedTeamSize: 3,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "web-api-clean",
		Blueprint:    "web-api-clean",
		Architecture: "clean",
		DisplayName:  "Clean Architecture API",
		Description:  "Clean Architecture with strict separation of concerns",
		
		PerformanceScore:     0.8,
		ComplexityScore:      0.6,
		ScalabilityScore:     0.8,
		MaintainabilityScore: 0.9,
		
		IsDistributed:    false,
		EstimatedFiles:   68,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Excellent separation of concerns",
			"Highly testable",
			"Framework independent",
			"Scales well with team size",
		},
		Cons: []string{
			"Higher initial complexity",
			"More boilerplate code",
			"Steeper learning curve",
		},
		
		BestFor: []string{
			"Medium to large applications",
			"Multiple team development",
			"Long-term maintenance",
			"Complex business logic",
		},
		AvoidIf: []string{
			"Simple CRUD operations",
			"Tight deadlines",
			"Small utility APIs",
		},
		
		MinTeamSize:         2,
		RecommendedTeamSize: 5,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "web-api-ddd",
		Blueprint:    "web-api-ddd",
		Architecture: "ddd",
		DisplayName:  "Domain-Driven Design API",
		Description:  "DDD approach focusing on rich domain models and business logic",
		
		PerformanceScore:     0.8,
		ComplexityScore:      0.8,
		ScalabilityScore:     0.9,
		MaintainabilityScore: 0.9,
		
		IsDistributed:    false,
		EstimatedFiles:   84,
		LearningCurve:    "high",
		
		Pros: []string{
			"Perfect for complex domains",
			"Rich business logic modeling",
			"Excellent for large teams",
			"Domain expert collaboration",
		},
		Cons: []string{
			"High initial complexity",
			"Requires domain expertise",
			"Overkill for simple domains",
		},
		
		BestFor: []string{
			"Complex business domains",
			"Enterprise applications",
			"Fintech and healthcare",
			"Large development teams",
		},
		AvoidIf: []string{
			"Simple CRUD operations",
			"Prototype development",
			"Small teams without domain experts",
		},
		
		MinTeamSize:         3,
		RecommendedTeamSize: 8,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "web-api-hexagonal",
		Blueprint:    "web-api-hexagonal",
		Architecture: "hexagonal",
		DisplayName:  "Hexagonal Architecture API",
		Description:  "Ports and adapters pattern for maximum testability",
		
		PerformanceScore:     0.7,
		ComplexityScore:      0.7,
		ScalabilityScore:     0.8,
		MaintainabilityScore: 1.0,
		
		IsDistributed:    false,
		EstimatedFiles:   75,
		LearningCurve:    "high",
		
		Pros: []string{
			"Maximum testability",
			"Excellent isolation",
			"Technology agnostic",
			"Easy to change adapters",
		},
		Cons: []string{
			"Complex abstraction layers",
			"High learning curve",
			"More initial development time",
		},
		
		BestFor: []string{
			"Applications requiring high test coverage",
			"Systems with multiple integrations",
			"Long-term maintenance focus",
			"Technology migration scenarios",
		},
		AvoidIf: []string{
			"Simple applications",
			"Tight deadlines",
			"Teams new to hexagonal patterns",
		},
		
		MinTeamSize:         2,
		RecommendedTeamSize: 4,
	})
	
	// Serverless Patterns
	kb.addPattern(&ArchitecturePattern{
		ID:           "lambda",
		Blueprint:    "lambda",
		Architecture: "standard",
		DisplayName:  "AWS Lambda Function",
		Description:  "Standard serverless function for event processing",
		
		PerformanceScore:     0.6,
		ComplexityScore:      0.3,
		ScalabilityScore:     0.9,
		MaintainabilityScore: 0.7,
		
		IsDistributed:    true,
		EstimatedFiles:   15,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Automatic scaling",
			"Pay per use",
			"No server management",
			"Built-in monitoring",
		},
		Cons: []string{
			"Cold start latency",
			"Vendor lock-in",
			"Limited execution time",
		},
		
		BestFor: []string{
			"Event-driven processing",
			"Cost-conscious projects",
			"Variable workloads",
			"Quick deployment needs",
		},
		AvoidIf: []string{
			"Long-running processes",
			"Low latency requirements",
			"Complex state management",
		},
		
		MinTeamSize:         1,
		RecommendedTeamSize: 2,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "lambda-event-processing",
		Blueprint:    "lambda-event-processing",
		Architecture: "event-processing",
		DisplayName:  "Lambda Event Processing",
		Description:  "Advanced serverless event processing with comprehensive observability",
		
		PerformanceScore:     0.8,
		ComplexityScore:      0.5,
		ScalabilityScore:     1.0,
		MaintainabilityScore: 0.8,
		
		IsDistributed:    true,
		EstimatedFiles:   43,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Enterprise-grade observability",
			"Multi-source event processing",
			"Production resilience patterns",
			"Advanced monitoring and tracing",
		},
		Cons: []string{
			"Higher complexity than basic Lambda",
			"AWS-specific patterns",
			"Learning curve for event-driven architecture",
		},
		
		BestFor: []string{
			"Production event processing",
			"Multiple AWS event sources",
			"Enterprise monitoring needs",
			"High-availability requirements",
		},
		AvoidIf: []string{
			"Simple single-purpose functions",
			"Non-AWS environments",
			"Basic event processing needs",
		},
		
		MinTeamSize:         2,
		RecommendedTeamSize: 4,
	})
	
	// Distributed Patterns
	kb.addPattern(&ArchitecturePattern{
		ID:           "microservice",
		Blueprint:    "microservice",
		Architecture: "microservice",
		DisplayName:  "Microservice",
		Description:  "Distributed microservice with gRPC and containerization",
		
		PerformanceScore:     0.9,
		ComplexityScore:      0.8,
		ScalabilityScore:     1.0,
		MaintainabilityScore: 0.7,
		
		IsDistributed:    true,
		EstimatedFiles:   55,
		LearningCurve:    "high",
		
		Pros: []string{
			"Independent deployments",
			"Technology diversity",
			"Team autonomy",
			"Excellent scalability",
		},
		Cons: []string{
			"Distributed system complexity",
			"Network latency",
			"Data consistency challenges",
		},
		
		BestFor: []string{
			"Large distributed systems",
			"Multiple teams",
			"High scalability needs",
			"Independent service evolution",
		},
		AvoidIf: []string{
			"Small applications",
			"Single team projects",
			"Tight data consistency needs",
		},
		
		MinTeamSize:         3,
		RecommendedTeamSize: 6,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "monolith",
		Blueprint:    "monolith",
		Architecture: "modular",
		DisplayName:  "Modular Monolith",
		Description:  "Well-structured monolith with clear module boundaries",
		
		PerformanceScore:     0.8,
		ComplexityScore:      0.5,
		ScalabilityScore:     0.6,
		MaintainabilityScore: 0.8,
		
		IsDistributed:    false,
		EstimatedFiles:   65,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Simple deployment",
			"Easy testing",
			"Good performance",
			"Simpler debugging",
		},
		Cons: []string{
			"Scaling limitations",
			"Technology constraints",
			"Potential coupling issues",
		},
		
		BestFor: []string{
			"Traditional web applications",
			"Single team development",
			"Well-understood domains",
			"Simple deployment needs",
		},
		AvoidIf: []string{
			"Massive scale requirements",
			"Multiple independent teams",
			"Technology diversity needs",
		},
		
		MinTeamSize:         2,
		RecommendedTeamSize: 8,
	})
	
	// Specialized Patterns
	kb.addPattern(&ArchitecturePattern{
		ID:           "event-driven",
		Blueprint:    "event-driven",
		Architecture: "event-sourcing",
		DisplayName:  "Event-Driven Architecture",
		Description:  "CQRS and Event Sourcing for complex event processing",
		
		PerformanceScore:     0.9,
		ComplexityScore:      0.9,
		ScalabilityScore:     1.0,
		MaintainabilityScore: 0.6,
		
		IsDistributed:    true,
		EstimatedFiles:   78,
		LearningCurve:    "high",
		
		Pros: []string{
			"Audit trail out of the box",
			"Excellent for analytics",
			"High scalability",
			"Temporal querying",
		},
		Cons: []string{
			"High complexity",
			"Eventual consistency",
			"Storage overhead",
		},
		
		BestFor: []string{
			"Financial systems",
			"Audit-heavy applications",
			"Complex business workflows",
			"Analytics-driven systems",
		},
		AvoidIf: []string{
			"Simple CRUD operations",
			"Teams new to event sourcing",
			"Strong consistency requirements",
		},
		
		MinTeamSize:         3,
		RecommendedTeamSize: 5,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "library",
		Blueprint:    "library",
		Architecture: "package",
		DisplayName:  "Go Library/Package",
		Description:  "Reusable Go library with clean public API",
		
		PerformanceScore:     0.9,
		ComplexityScore:      0.2,
		ScalabilityScore:     0.8,
		MaintainabilityScore: 0.9,
		
		IsDistributed:    false,
		EstimatedFiles:   12,
		LearningCurve:    "low",
		
		Pros: []string{
			"Highly reusable",
			"Clean API design",
			"Excellent performance",
			"Easy to test",
		},
		Cons: []string{
			"Limited to library use cases",
			"No standalone deployment",
			"API stability concerns",
		},
		
		BestFor: []string{
			"Shared utilities",
			"SDK development",
			"Common functionality",
			"Open source packages",
		},
		AvoidIf: []string{
			"Standalone applications",
			"User interfaces",
			"Complex workflows",
		},
		
		MinTeamSize:         1,
		RecommendedTeamSize: 2,
	})
	
	kb.addPattern(&ArchitecturePattern{
		ID:           "workspace",
		Blueprint:    "workspace",
		Architecture: "monorepo",
		DisplayName:  "Go Workspace",
		Description:  "Multi-module workspace for managing related projects",
		
		PerformanceScore:     0.7,
		ComplexityScore:      0.6,
		ScalabilityScore:     0.9,
		MaintainabilityScore: 0.8,
		
		IsDistributed:    false,
		EstimatedFiles:   45,
		LearningCurve:    "medium",
		
		Pros: []string{
			"Shared dependencies",
			"Coordinated releases",
			"Code sharing",
			"Unified tooling",
		},
		Cons: []string{
			"Coordination overhead",
			"Build complexity",
			"Version management challenges",
		},
		
		BestFor: []string{
			"Related project suites",
			"Monorepo strategies",
			"Shared library development",
			"Coordinated team development",
		},
		AvoidIf: []string{
			"Independent projects",
			"Single module needs",
			"Simple applications",
		},
		
		MinTeamSize:         2,
		RecommendedTeamSize: 10,
	})
}

// addPattern adds a pattern to the knowledge base
func (kb *ArchitectureKnowledgeBase) addPattern(pattern *ArchitecturePattern) {
	kb.patterns[pattern.ID] = pattern
}

// GetPatternsByComplexity returns patterns suitable for a given complexity level
func (kb *ArchitectureKnowledgeBase) GetPatternsByComplexity(maxComplexity float64) []*ArchitecturePattern {
	var patterns []*ArchitecturePattern
	for _, pattern := range kb.patterns {
		if pattern.ComplexityScore <= maxComplexity {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

// GetPatternsByPerformance returns patterns suitable for performance requirements
func (kb *ArchitectureKnowledgeBase) GetPatternsByPerformance(minPerformance float64) []*ArchitecturePattern {
	var patterns []*ArchitecturePattern
	for _, pattern := range kb.patterns {
		if pattern.PerformanceScore >= minPerformance {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}

// GetDistributedPatterns returns all distributed architecture patterns
func (kb *ArchitectureKnowledgeBase) GetDistributedPatterns() []*ArchitecturePattern {
	var patterns []*ArchitecturePattern
	for _, pattern := range kb.patterns {
		if pattern.IsDistributed {
			patterns = append(patterns, pattern)
		}
	}
	return patterns
}