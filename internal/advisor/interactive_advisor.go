package advisor

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/francknouama/go-starter/internal/prompts/interfaces"
	"github.com/francknouama/go-starter/pkg/types"
)

// InteractiveAdvisor provides an interactive interface for gathering requirements
type InteractiveAdvisor struct {
	advisor       *ArchitectureAdvisor
	surveyAdapter interfaces.SurveyAdapter
}

// NewInteractiveAdvisor creates a new interactive advisor
func NewInteractiveAdvisor() *InteractiveAdvisor {
	return &InteractiveAdvisor{
		advisor:       NewArchitectureAdvisor(),
		surveyAdapter: &interfaces.RealSurveyAdapter{},
	}
}

// NewInteractiveAdvisorWithAdapter creates an advisor with a custom survey adapter (for testing)
func NewInteractiveAdvisorWithAdapter(adapter interfaces.SurveyAdapter) *InteractiveAdvisor {
	return &InteractiveAdvisor{
		advisor:       NewArchitectureAdvisor(),
		surveyAdapter: adapter,
	}
}

// GatherRequirementsAndRecommend gathers requirements interactively and provides recommendations
func (ia *InteractiveAdvisor) GatherRequirementsAndRecommend() (*ArchitectureRecommendation, error) {
	fmt.Println("🤖 Welcome to the AI-Powered Architecture Advisor!")
	fmt.Println("I'll help you choose the perfect architecture for your Go project.")
	
	requirements, err := ia.gatherRequirements()
	if err != nil {
		return nil, fmt.Errorf("failed to gather requirements: %w", err)
	}
	
	recommendation, err := ia.advisor.AnalyzeRequirements(*requirements)
	if err != nil {
		return nil, fmt.Errorf("failed to generate recommendation: %w", err)
	}
	
	ia.presentRecommendation(recommendation)
	
	return recommendation, nil
}

// gatherRequirements collects project requirements through interactive prompts
func (ia *InteractiveAdvisor) gatherRequirements() (*ProjectRequirements, error) {
	req := &ProjectRequirements{}
	
	// Basic project information
	if err := ia.gatherBasicInfo(req); err != nil {
		return nil, err
	}
	
	// Technical requirements
	if err := ia.gatherTechnicalRequirements(req); err != nil {
		return nil, err
	}
	
	// Business requirements
	if err := ia.gatherBusinessRequirements(req); err != nil {
		return nil, err
	}
	
	// Integration and operational requirements
	if err := ia.gatherIntegrationRequirements(req); err != nil {
		return nil, err
	}
	
	// User experience indicators
	if err := ia.gatherUserExperience(req); err != nil {
		return nil, err
	}
	
	return req, nil
}

// gatherBasicInfo collects basic project information
func (ia *InteractiveAdvisor) gatherBasicInfo(req *ProjectRequirements) error {
	questions := []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "What's your project name?",
				Help:    "A descriptive name for your Go project",
			},
			Validate: survey.Required,
		},
		{
			Name: "description",
			Prompt: &survey.Input{
				Message: "Briefly describe your project:",
				Help:    "What does your project do? (e.g., 'REST API for e-commerce', 'CLI tool for deployment')",
			},
		},
		{
			Name: "domain",
			Prompt: &survey.Select{
				Message: "What domain does your project belong to?",
				Options: []string{
					"e-commerce",
					"fintech",
					"healthcare",
					"iot",
					"content-management",
					"devtools",
					"gaming",
					"education",
					"enterprise",
					"other",
				},
				Help: "This helps me recommend domain-specific patterns",
			},
		},
	}
	
	answers := struct {
		Name        string
		Description string
		Domain      string
	}{}
	
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}
	
	req.Name = answers.Name
	req.Description = answers.Description
	req.Domain = answers.Domain
	
	return nil
}

// gatherTechnicalRequirements collects technical requirements
func (ia *InteractiveAdvisor) gatherTechnicalRequirements(req *ProjectRequirements) error {
	questions := []*survey.Question{
		{
			Name: "expected_load",
			Prompt: &survey.Select{
				Message: "What's your expected load?",
				Options: []string{"low", "medium", "high", "massive"},
				Default: "medium",
				Help:    "Low: <1000 users, Medium: <10K users, High: <100K users, Massive: >100K users",
			},
		},
		{
			Name: "concurrent_users",
			Prompt: &survey.Input{
				Message: "Expected concurrent users (approximate):",
				Default: "100",
				Help:    "How many users will use your system simultaneously?",
			},
			Validate: func(input interface{}) error {
				if str, ok := input.(string); ok {
					if _, err := strconv.Atoi(str); err != nil {
						return fmt.Errorf("please enter a valid number")
					}
				}
				return nil
			},
		},
		{
			Name: "data_volume",
			Prompt: &survey.Select{
				Message: "Expected data volume?",
				Options: []string{"small", "medium", "large", "big-data"},
				Default: "medium",
				Help:    "Small: <1GB, Medium: <100GB, Large: <10TB, Big-data: >10TB",
			},
		},
		{
			Name: "response_time",
			Prompt: &survey.Select{
				Message: "Response time requirements?",
				Options: []string{"relaxed", "standard", "fast", "realtime"},
				Default: "standard",
				Help:    "Relaxed: >1s, Standard: <500ms, Fast: <100ms, Realtime: <10ms",
			},
		},
		{
			Name: "availability",
			Prompt: &survey.Select{
				Message: "Availability requirements?",
				Options: []string{"standard", "high", "critical"},
				Default: "standard",
				Help:    "Standard: 99%, High: 99.9%, Critical: 99.99%+",
			},
		},
	}
	
	answers := struct {
		ExpectedLoad     string
		ConcurrentUsers  string
		DataVolume       string
		ResponseTime     string
		Availability     string
	}{}
	
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}
	
	req.ExpectedLoad = answers.ExpectedLoad
	req.DataVolume = answers.DataVolume
	req.ResponseTime = answers.ResponseTime
	req.Availability = answers.Availability
	
	// Convert concurrent users to int
	if users, err := strconv.Atoi(answers.ConcurrentUsers); err == nil {
		req.ConcurrentUsers = users
	}
	
	return nil
}

// gatherBusinessRequirements collects business requirements
func (ia *InteractiveAdvisor) gatherBusinessRequirements(req *ProjectRequirements) error {
	questions := []*survey.Question{
		{
			Name: "time_to_market",
			Prompt: &survey.Select{
				Message: "Time to market priority?",
				Options: []string{"mvp", "quick", "standard", "thorough"},
				Default: "standard",
				Help:    "MVP: weeks, Quick: 1-2 months, Standard: 3-6 months, Thorough: >6 months",
			},
		},
		{
			Name: "budget",
			Prompt: &survey.Select{
				Message: "Project budget category?",
				Options: []string{"startup", "small", "medium", "enterprise"},
				Default: "medium",
				Help:    "This affects complexity and team size recommendations",
			},
		},
		{
			Name: "team_size",
			Prompt: &survey.Input{
				Message: "Development team size:",
				Default: "3",
				Help:    "Number of developers working on this project",
			},
			Validate: func(input interface{}) error {
				if str, ok := input.(string); ok {
					if size, err := strconv.Atoi(str); err != nil || size < 1 {
						return fmt.Errorf("please enter a valid team size (>= 1)")
					}
				}
				return nil
			},
		},
		{
			Name: "team_experience",
			Prompt: &survey.Select{
				Message: "Team's Go experience level?",
				Options: []string{"junior", "mixed", "senior", "expert"},
				Default: "mixed",
				Help:    "This affects complexity recommendations",
			},
		},
	}
	
	answers := struct {
		TimeToMarket   string
		Budget         string
		TeamSize       string
		TeamExperience string
	}{}
	
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}
	
	req.TimeToMarket = answers.TimeToMarket
	req.Budget = answers.Budget
	req.TeamExperience = answers.TeamExperience
	
	// Convert team size to int
	if size, err := strconv.Atoi(answers.TeamSize); err == nil {
		req.TeamSize = size
	}
	
	return nil
}

// gatherIntegrationRequirements collects integration and operational requirements
func (ia *InteractiveAdvisor) gatherIntegrationRequirements(req *ProjectRequirements) error {
	// Third-party services
	thirdPartyOptions := []string{
		"Payment processors (Stripe, PayPal)",
		"Email services (SendGrid, Mailgun)",
		"SMS services (Twilio, AWS SNS)",
		"Social media APIs",
		"Cloud storage (S3, GCS)",
		"Analytics (Google Analytics, Mixpanel)",
		"Monitoring (New Relic, DataDog)",
		"Other APIs",
		"None",
	}
	
	var selectedServices []string
	if err := ia.surveyAdapter.AskOne(&survey.MultiSelect{
		Message: "Select third-party integrations needed:",
		Options: thirdPartyOptions,
		Help:    "Choose all that apply",
	}, &selectedServices); err != nil {
		return err
	}
	req.ThirdPartyServices = selectedServices
	
	// Database requirements
	var dbRequirements string
	if err := ia.surveyAdapter.AskOne(&survey.Select{
		Message: "Database complexity?",
		Options: []string{"simple", "complex", "distributed"},
		Default: "simple",
		Help:    "Simple: single table/collection, Complex: multiple tables with relations, Distributed: sharding/replication",
	}, &dbRequirements); err != nil {
		return err
	}
	req.DatabaseRequirements = dbRequirements
	
	// Authentication requirements
	var authRequirements string
	if err := ia.surveyAdapter.AskOne(&survey.Select{
		Message: "Authentication requirements?",
		Options: []string{"none", "basic", "oauth", "enterprise"},
		Default: "basic",
		Help:    "None: no auth, Basic: username/password, OAuth: social login, Enterprise: SSO/LDAP",
	}, &authRequirements); err != nil {
		return err
	}
	req.AuthRequirements = authRequirements
	
	// Deployment target
	var deploymentTarget string
	if err := ia.surveyAdapter.AskOne(&survey.Select{
		Message: "Primary deployment target?",
		Options: []string{"local", "cloud", "hybrid", "edge"},
		Default: "cloud",
		Help:    "Where will your application primarily run?",
	}, &deploymentTarget); err != nil {
		return err
	}
	req.DeploymentTarget = deploymentTarget
	
	// Compliance needs
	complianceOptions := []string{
		"GDPR (EU data protection)",
		"HIPAA (Healthcare)",
		"SOX (Financial reporting)",
		"PCI DSS (Payment cards)",
		"ISO 27001",
		"None",
	}
	
	var complianceNeeds []string
	if err := ia.surveyAdapter.AskOne(&survey.MultiSelect{
		Message: "Compliance requirements:",
		Options: complianceOptions,
		Help:    "Select all applicable compliance standards",
	}, &complianceNeeds); err != nil {
		return err
	}
	req.ComplianceNeeds = complianceNeeds
	
	// Monitoring needs
	var monitoringNeeds string
	if err := ia.surveyAdapter.AskOne(&survey.Select{
		Message: "Monitoring and observability needs?",
		Options: []string{"basic", "standard", "advanced", "enterprise"},
		Default: "standard",
		Help:    "Basic: logs only, Standard: logs + metrics, Advanced: + tracing, Enterprise: + custom dashboards",
	}, &monitoringNeeds); err != nil {
		return err
	}
	req.MonitoringNeeds = monitoringNeeds
	
	return nil
}

// gatherUserExperience collects user experience indicators
func (ia *InteractiveAdvisor) gatherUserExperience(req *ProjectRequirements) error {
	questions := []*survey.Question{
		{
			Name: "first_go_project",
			Prompt: &survey.Confirm{
				Message: "Is this your first Go project?",
				Default: false,
				Help:    "This helps me recommend appropriate complexity levels",
			},
		},
		{
			Name: "microservice_experience",
			Prompt: &survey.Confirm{
				Message: "Do you have microservice architecture experience?",
				Default: false,
				Help:    "Experience with distributed systems and microservices",
			},
		},
		{
			Name: "preferred_style",
			Prompt: &survey.Select{
				Message: "Preferred development style?",
				Options: []string{"simple", "structured", "enterprise"},
				Default: "structured",
				Help:    "Simple: minimal files, Structured: organized layers, Enterprise: comprehensive patterns",
			},
		},
	}
	
	answers := struct {
		FirstGoProject        bool
		MicroserviceExperience bool
		PreferredStyle        string
	}{}
	
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}
	
	req.IsFirstGoProject = answers.FirstGoProject
	req.HasMicroserviceExp = answers.MicroserviceExperience
	req.PreferredStyle = answers.PreferredStyle
	
	return nil
}

// presentRecommendation displays the recommendation in a user-friendly format
func (ia *InteractiveAdvisor) presentRecommendation(rec *ArchitectureRecommendation) {
	fmt.Println("\n🎯 ARCHITECTURE RECOMMENDATION")
	fmt.Println(strings.Repeat("=", 50))
	
	// Main recommendation
	fmt.Printf("📋 Recommended Blueprint: %s\n", rec.Blueprint)
	fmt.Printf("🏗️  Architecture Pattern: %s\n", rec.Architecture)
	fmt.Printf("⚡ Complexity Level: %s\n", rec.Complexity.String())
	fmt.Printf("🎯 Confidence: %.1f%%\n\n", rec.Confidence*100)
	
	// Reasoning
	if len(rec.Reasoning) > 0 {
		fmt.Println("💡 Why this recommendation:")
		for _, reason := range rec.Reasoning {
			fmt.Printf("   • %s\n", reason)
		}
		fmt.Println()
	}
	
	// Pros and Cons
	if len(rec.Pros) > 0 {
		fmt.Println("✅ Advantages:")
		for _, pro := range rec.Pros {
			fmt.Printf("   • %s\n", pro)
		}
		fmt.Println()
	}
	
	if len(rec.Cons) > 0 {
		fmt.Println("⚠️  Considerations:")
		for _, con := range rec.Cons {
			fmt.Printf("   • %s\n", con)
		}
		fmt.Println()
	}
	
	// Technical recommendations
	fmt.Println("🔧 Technical Recommendations:")
	fmt.Printf("   • Framework: %s\n", rec.Framework)
	fmt.Printf("   • Logger: %s\n", rec.Logger)
	fmt.Printf("   • Database: %s\n", rec.Database)
	fmt.Printf("   • Estimated Files: %d\n", rec.EstimatedFiles)
	fmt.Printf("   • Development Time: %s\n\n", rec.DevelopmentTime)
	
	// Alternatives
	if len(rec.Alternatives) > 0 {
		fmt.Println("🔄 Alternative Options:")
		for i, alt := range rec.Alternatives {
			fmt.Printf("   %d. %s (%.1f%% confidence) - %s\n", 
				i+1, alt.Blueprint, alt.Confidence*100, alt.Reason)
		}
		fmt.Println()
	}
	
	fmt.Println("🚀 Ready to generate? Use the recommended settings above!")
}

// GetProjectConfigFromRecommendation converts a recommendation to ProjectConfig
func (ia *InteractiveAdvisor) GetProjectConfigFromRecommendation(rec *ArchitectureRecommendation, projectName, modulePath string) types.ProjectConfig {
	config := types.ProjectConfig{
		Name:      projectName,
		Type:      rec.Blueprint,
		Module:    modulePath,
		GoVersion: "1.21", // Default, can be customized
		Framework: rec.Framework,
		Logger:    rec.Logger,
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: rec.Database,
			},
		},
	}
	
	// Apply features from recommendation
	for key, value := range rec.Features {
		switch key {
		case "auth":
			if config.Features.Authentication.Type == "" {
				config.Features.Authentication.Type = value
			}
		case "monitoring":
			if value == "enterprise" || value == "advanced" {
				config.Features.Monitoring.Metrics = true
				config.Features.Monitoring.Logging = true
				config.Features.Monitoring.Tracing = true
			}
		case "compliance":
			// Set compliance-related features
		}
	}
	
	return config
}

// QuickRecommendation provides a quick recommendation based on minimal input
func (ia *InteractiveAdvisor) QuickRecommendation(projectType, domain, teamExperience string) (*ArchitectureRecommendation, error) {
	// Create minimal requirements for quick recommendation
	req := ProjectRequirements{
		Domain:         domain,
		TeamExperience: teamExperience,
		TimeToMarket:   "standard",
		ExpectedLoad:   "medium",
		ResponseTime:   "standard",
		Budget:         "medium",
		TeamSize:       3,
	}
	
	// Map project type to more specific requirements
	switch projectType {
	case "api", "web-api":
		req.DeploymentTarget = "cloud"
		req.DatabaseRequirements = "simple"
	case "cli":
		req.DeploymentTarget = "local"
		req.DatabaseRequirements = "simple"
		req.PreferredStyle = "simple" // CLI tools are often simpler
	case "microservice":
		req.ExpectedLoad = "high"
		req.DeploymentTarget = "cloud"
		req.HasMicroserviceExp = true
	case "lambda":
		req.DeploymentTarget = "cloud"
		req.PreferredStyle = "simple"
	}
	
	// Get the recommendation from the algorithm
	recommendation, err := ia.advisor.AnalyzeRequirements(req)
	if err != nil {
		return nil, err
	}
	
	// For CLI type, if the algorithm didn't recommend CLI, provide a CLI alternative
	if projectType == "cli" && recommendation.Blueprint != "cli" {
		// Create a CLI-focused recommendation as an override
		cliRecommendation := &ArchitectureRecommendation{
			Blueprint:    "cli",
			Architecture: "standard",
			Complexity:   1,
			Confidence:   0.9, // High confidence since user explicitly requested CLI
			Reasoning: []string{
				"User explicitly requested CLI tool type",
				"CLI architecture matches devtools domain requirements",
				"Suitable for team experience level: " + teamExperience,
			},
			Pros: []string{
				"Direct command-line interface",
				"Easy to deploy and distribute",
				"Perfect for developer tooling",
				"Simple to test and maintain",
			},
			Cons: []string{
				"Limited to command-line usage",
				"No web interface",
				"Platform-specific distribution",
			},
			Alternatives: []AlternativeRecommendation{
				{
					Blueprint:  recommendation.Blueprint,
					Confidence: recommendation.Confidence,
					Reason:     fmt.Sprintf("Algorithm suggested %s with %.1f%% confidence", recommendation.Blueprint, recommendation.Confidence*100),
				},
			},
			Framework:        "cobra",
			Logger:           "slog",
			Database:         "",
			Features:         recommendation.Features,
			EstimatedFiles:   15,
			DevelopmentTime:  "2-3 days",
		}
		return cliRecommendation, nil
	}
	
	return recommendation, nil
}