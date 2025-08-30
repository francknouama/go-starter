package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/francknouama/go-starter/internal/advisor"
)

var (
	advisorOutputFormat string
	advisorQuickMode    bool
	advisorProjectType  string
	advisorDomain       string
	advisorTeamExp      string
)

// advisorCmd represents the advisor command
var advisorCmd = &cobra.Command{
	Use:   "advisor",
	Short: "AI-powered architecture advisor for Go projects",
	Long: `The AI-powered architecture advisor analyzes your project requirements 
and recommends the most suitable blueprint and architecture pattern.

It considers factors like:
• Team experience and size
• Performance requirements  
• Business constraints
• Domain-specific needs
• Technical complexity

Examples:
  # Interactive mode (recommended)
  go-starter advisor

  # Quick recommendation
  go-starter advisor --quick --type=api --domain=e-commerce --team=mixed

  # Output as JSON for automation
  go-starter advisor --quick --format=json --type=cli --domain=devtools`,
	RunE: runAdvisor,
}

func init() {
	rootCmd.AddCommand(advisorCmd)

	// Output format flag
	advisorCmd.Flags().StringVarP(&advisorOutputFormat, "format", "f", "interactive", 
		"Output format: interactive, json, summary")
	
	// Quick mode flags
	advisorCmd.Flags().BoolVarP(&advisorQuickMode, "quick", "q", false,
		"Quick recommendation mode (requires --type, --domain, --team)")
	advisorCmd.Flags().StringVar(&advisorProjectType, "type", "",
		"Project type for quick mode: api, cli, microservice, lambda")
	advisorCmd.Flags().StringVar(&advisorDomain, "domain", "",
		"Domain for quick mode: e-commerce, fintech, iot, devtools, etc.")
	advisorCmd.Flags().StringVar(&advisorTeamExp, "team", "",
		"Team experience for quick mode: junior, mixed, senior, expert")

	// Mark dependencies for quick mode
	advisorCmd.MarkFlagsRequiredTogether("quick", "type", "domain", "team")
}

func runAdvisor(cmd *cobra.Command, args []string) error {
	// Initialize interactive advisor
	interactiveAdvisor := advisor.NewInteractiveAdvisor()
	
	var recommendation *advisor.ArchitectureRecommendation
	var err error
	
	if advisorQuickMode {
		// Quick recommendation mode
		// Only show progress message for non-JSON output to avoid parsing issues
		if advisorOutputFormat != "json" {
			fmt.Printf("🤖 Generating quick architecture recommendation...\n")
		}
		
		recommendation, err = interactiveAdvisor.QuickRecommendation(
			advisorProjectType, 
			advisorDomain, 
			advisorTeamExp,
		)
		if err != nil {
			return fmt.Errorf("failed to generate quick recommendation: %w", err)
		}
	} else {
		// Interactive mode
		if advisorOutputFormat == "json" {
			return fmt.Errorf("interactive mode cannot use JSON format, use --quick for JSON output")
		}
		
		recommendation, err = interactiveAdvisor.GatherRequirementsAndRecommend()
		if err != nil {
			return fmt.Errorf("interactive advisor failed: %w", err)
		}
	}
	
	// Output the recommendation
	return outputRecommendation(recommendation, advisorOutputFormat)
}

func outputRecommendation(rec *advisor.ArchitectureRecommendation, format string) error {
	switch format {
	case "json":
		return outputJSON(rec)
	case "summary":
		return outputSummary(rec)
	case "interactive":
		// Interactive format already displayed by the advisor
		return outputGenerationCommand(rec)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}

func outputJSON(rec *advisor.ArchitectureRecommendation) error {
	jsonData, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal recommendation to JSON: %w", err)
	}
	
	fmt.Println(string(jsonData))
	return nil
}

func outputSummary(rec *advisor.ArchitectureRecommendation) error {
	fmt.Printf("Blueprint: %s\n", rec.Blueprint)
	fmt.Printf("Architecture: %s\n", rec.Architecture)
	fmt.Printf("Complexity: %s\n", rec.Complexity.String())
	fmt.Printf("Confidence: %.1f%%\n", rec.Confidence*100)
	fmt.Printf("Framework: %s\n", rec.Framework)
	fmt.Printf("Logger: %s\n", rec.Logger)
	fmt.Printf("Database: %s\n", rec.Database)
	fmt.Printf("Estimated Files: %d\n", rec.EstimatedFiles)
	fmt.Printf("Development Time: %s\n", rec.DevelopmentTime)
	
	return nil
}

func outputGenerationCommand(rec *advisor.ArchitectureRecommendation) error {
	fmt.Println("\n🚀 READY TO GENERATE?")
	fmt.Println("Use this command to generate your project:")
	fmt.Println()
	
	// Build the recommended command
	cmd := "go-starter new my-project"
	cmd += fmt.Sprintf(" --type=%s", rec.Blueprint)
	cmd += fmt.Sprintf(" --framework=%s", rec.Framework)
	cmd += fmt.Sprintf(" --logger=%s", rec.Logger)
	cmd += fmt.Sprintf(" --complexity=%s", rec.Complexity.String())
	
	if rec.Database != "" && rec.Database != "none" {
		cmd += fmt.Sprintf(" --database-driver=%s", rec.Database)
	}
	
	if rec.Architecture != "standard" {
		cmd += fmt.Sprintf(" --architecture=%s", rec.Architecture)
	}
	
	// Add module path placeholder
	cmd += " --module=example.com/my-project"
	
	fmt.Printf("  %s\n", cmd)
	fmt.Println()
	fmt.Println("💡 Customize the project name, module path, and any other options as needed!")
	
	return nil
}

// Helper functions for advisor integration

// GetAdvisorRecommendation provides a programmatic interface for other commands
func GetAdvisorRecommendation(projectType, domain, teamExperience string) (*advisor.ArchitectureRecommendation, error) {
	interactiveAdvisor := advisor.NewInteractiveAdvisor()
	return interactiveAdvisor.QuickRecommendation(projectType, domain, teamExperience)
}

// ValidateAdvisorFlags validates advisor command flags
func ValidateAdvisorFlags() error {
	if advisorQuickMode {
		if advisorProjectType == "" {
			return fmt.Errorf("--type is required in quick mode")
		}
		if advisorDomain == "" {
			return fmt.Errorf("--domain is required in quick mode")
		}
		if advisorTeamExp == "" {
			return fmt.Errorf("--team is required in quick mode")
		}
		
		// Validate project type
		validTypes := []string{"api", "web-api", "cli", "microservice", "lambda", "library"}
		if !contains(validTypes, advisorProjectType) {
			return fmt.Errorf("invalid project type: %s (valid options: %v)", advisorProjectType, validTypes)
		}
		
		// Validate team experience
		validTeamExp := []string{"junior", "mixed", "senior", "expert"}
		if !contains(validTeamExp, advisorTeamExp) {
			return fmt.Errorf("invalid team experience: %s (valid options: %v)", advisorTeamExp, validTeamExp)
		}
	}
	
	return nil
}

// contains checks if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Integration helper for the new command to optionally use advisor
func SuggestArchitectureForNewCommand(projectType string) {
	if projectType == "" {
		return
	}
	
	// Only suggest if the user hasn't specified detailed options
	fmt.Printf("\n💡 Tip: For personalized architecture recommendations, try:\n")
	fmt.Printf("   go-starter advisor --quick --type=%s --domain=<your-domain> --team=<your-experience>\n", projectType)
	fmt.Println()
}