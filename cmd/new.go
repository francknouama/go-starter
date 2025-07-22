package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/ascii"
	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/prompts"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	projectName    string
	projectModule  string
	projectType    string
	architecture   string
	goVersion      string
	outputDir      string
	framework      string
	logger         string
	databaseDriver string
	databaseORM    string
	authType       string
	advanced       bool
	basic          bool
	complexity     string
	dryRun         bool
	noGit          bool
	randomName     bool
	quiet          bool
	noBanner       bool
	bannerStyle    string
	assetPipeline  string
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with the specified blueprint and configuration.

Examples:
  # Basic usage (beginner-friendly)
  go-starter new my-api                                          # Interactive mode (basic)
  go-starter new my-api --type=web-api --framework=gin           # Direct mode (basic)
  go-starter new my-cli --type=cli --complexity=simple           # Simple CLI project
  
  # Advanced usage (all options)
  go-starter new my-api --advanced                               # Interactive mode (advanced)
  go-starter new my-api --type=web-api --logger=zap --advanced   # Direct mode (advanced)
  go-starter new my-cli --complexity=standard                    # Standard CLI project
  
  # Complexity-based generation
  go-starter new my-proto --complexity=simple                    # Minimal structure
  go-starter new my-app --complexity=standard                    # Balanced structure
  go-starter new my-enterprise --complexity=advanced             # Enterprise structure

The command will guide you through the project configuration process
or use the provided flags for direct project generation.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)
	
	// Set custom help function for progressive disclosure
	newCmd.SetHelpFunc(progressiveHelpFunc)

	// Project configuration flags
	newCmd.Flags().StringVar(&projectName, "name", "", "Project name")
	newCmd.Flags().StringVar(&projectModule, "module", "", "Go module path (e.g., github.com/user/project)")
	newCmd.Flags().StringVar(&projectType, "type", "", "Project type (web-api, cli, library, lambda)")
	newCmd.Flags().StringVar(&architecture, "architecture", "", "Architecture pattern (standard, clean, ddd, hexagonal)")
	newCmd.Flags().StringVarP(&goVersion, "go-version", "g", "", "Go version to use (auto, 1.23, 1.22, 1.21)")
	newCmd.Flags().StringVar(&framework, "framework", "", "Framework to use (gin, echo, cobra, etc.)")
	newCmd.Flags().StringVar(&logger, "logger", "", "Logger to use (slog, zap, logrus, zerolog)")
	newCmd.Flags().StringVar(&outputDir, "output", ".", "Output directory")
	newCmd.Flags().StringVar(&databaseDriver, "database-driver", "", "Database driver (postgres, mysql, sqlite)")
	newCmd.Flags().StringVar(&databaseORM, "database-orm", "", "Database ORM/query builder (gorm, sqlx)")
	newCmd.Flags().StringVar(&authType, "auth-type", "", "Authentication type (jwt, oauth2, session)")
	newCmd.Flags().StringVar(&assetPipeline, "asset-pipeline", "", "Asset build system (embedded, webpack, vite, esbuild)")

	// Progressive disclosure options
	newCmd.Flags().BoolVar(&basic, "basic", false, "Show only essential options (default)")
	newCmd.Flags().BoolVar(&advanced, "advanced", false, "Enable advanced configuration")
	newCmd.Flags().StringVar(&complexity, "complexity", "", "Complexity level (simple, standard, advanced, expert)")
	
	// Generation options
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview project structure without creating files")
	newCmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git repository initialization")
	newCmd.Flags().BoolVar(&randomName, "random-name", false, "Generate a random project name (GitHub-style)")
	
	// Banner control options
	newCmd.Flags().BoolVarP(&quiet, "quiet", "q", false, "Suppress all output except errors")
	newCmd.Flags().BoolVar(&noBanner, "no-banner", false, "Disable banner display")
	newCmd.Flags().StringVar(&bannerStyle, "banner-style", "", "Banner style (full, minimal, none)")
}

func runNew(cmd *cobra.Command, args []string) error {
	// Validate complexity flag if provided
	if complexity != "" {
		if _, err := prompts.ParseComplexityLevel(complexity); err != nil {
			return fmt.Errorf("invalid complexity level: %w", err)
		}
	}
	
	// Determine disclosure mode based on flags
	disclosureMode := prompts.DetermineDisclosureMode(basic, advanced, complexity)
	
	// Configure banner display
	bannerConfig := ascii.GetBannerConfig(quiet, noBanner, bannerStyle)
	
	// Show welcome banner for new project generation
	if !bannerConfig.Quiet && bannerConfig.Enabled {
		welcomeStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true).
			MarginBottom(1)
		
		if !bannerConfig.Colors {
			welcomeStyle = lipgloss.NewStyle().Bold(true)
		}
		
		// Show minimal banner for new command
		if projectName == "" && projectType == "" {
			// Interactive mode - show full welcome
			ascii.PrintWelcomeWithConfig(bannerConfig)
		} else {
			// Direct mode - show minimal banner
			minimalConfig := *bannerConfig
			minimalConfig.Style = ascii.StyleMinimal
			fmt.Print(ascii.BannerWithConfig(&minimalConfig))
			fmt.Println(welcomeStyle.Render("🚀 Generating new Go project..."))
			fmt.Println()
		}
	}
	
	// Get project name from args if provided
	if len(args) > 0 {
		projectName = args[0]
	}

	// Generate random name if requested and no name provided
	if randomName && projectName == "" {
		projectName = utils.GenerateRandomProjectName()
		if !quiet {
			fmt.Printf("🎲 Generated random project name: %s\n", projectName)
		}
	}

	// Initialize the prompter for interactive configuration
	// Use the new factory pattern with Bubble Tea UI and Survey fallback
	prompter := prompts.NewDefault()

	// Parse complexity level if provided
	var complexityLevel prompts.ComplexityLevel
	if complexity != "" {
		complexityLevel, _ = prompts.ParseComplexityLevel(complexity)
	}

	// Adjust blueprint type based on complexity level BEFORE prompting
	actualProjectType := projectType
	actualFramework := framework
	if complexity != "" && projectType != "" {
		actualProjectType = prompts.SelectBlueprintForComplexity(projectType, complexityLevel)
		
		// Set defaults for CLI blueprints to avoid unnecessary prompts
		if projectType == "cli" {
			if framework == "" {
				actualFramework = "cobra" // CLI blueprints use Cobra
			}
			if logger == "" {
				logger = "slog" // Default to slog
			}
			// Set default module path if not provided
			if projectModule == "" && projectName != "" {
				projectModule = "github.com/username/" + projectName
			}
		}
	}

	// Apply defaults for all project types when sufficient flags are provided (regardless of complexity)
	if projectType != "" {
		// Set framework defaults for each project type
		if framework == "" {
			switch projectType {
			case "cli":
				actualFramework = "cobra"
			case "web-api", "monolith":
				actualFramework = "gin"
			case "microservice":
				actualFramework = "gin"
			default:
				// Library and lambda don't need frameworks
			}
		}

		// Set logger defaults for all project types that use logging
		if logger == "" && projectType != "library" {
			logger = "slog" // Default to Go's standard structured logging
		}

		// Set default module path for testing scenarios
		if projectModule == "" && projectName != "" {
			projectModule = "github.com/username/" + projectName
		}
	}

	// Get project configuration through interactive prompts or flags
	initialConfig := types.ProjectConfig{
		Name:         projectName,
		Module:       projectModule,
		Type:         actualProjectType,
		Architecture: architecture,
		GoVersion:    goVersion,
		Framework:    actualFramework,
		Logger:       logger,
		Features: &types.Features{
			Database: types.DatabaseConfig{
				Driver: databaseDriver,
				ORM:    databaseORM,
			},
			Authentication: types.AuthConfig{
				Type: authType,
			},
		},
		Variables: map[string]string{
			"AssetPipeline":  assetPipeline,
			"TemplateEngine": "", // Will use template default
			"SessionStore":   "", // Will use template default
			"DatabaseDriver": databaseDriver,
			"DatabaseORM":    databaseORM,
			"AuthType":       authType,
			"LoggerType":     logger,
		},
	}

	// Use new disclosure-aware method if available, fallback to old method
	var config types.ProjectConfig
	var err error
	if disclosurePrompter, ok := prompter.(interface {
		GetProjectConfigWithDisclosure(types.ProjectConfig, prompts.DisclosureMode, prompts.ComplexityLevel) (types.ProjectConfig, error)
	}); ok {
		config, err = disclosurePrompter.GetProjectConfigWithDisclosure(initialConfig, disclosureMode, complexityLevel)
	} else {
		// Fallback to old method
		config, err = prompter.GetProjectConfig(initialConfig, advanced)
	}
	if err != nil {
		printErrorMessage("Failed to get project configuration", err)
		return fmt.Errorf("failed to get project configuration: %w", err)
	}

	// Blueprint type was already adjusted before prompting

	// Validate the configuration
	if err := validateConfig(config); err != nil {
		printErrorMessage("Invalid configuration", err)
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Initialize the generator
	gen := generator.New()

	// Handle dry run mode
	if dryRun {
		return gen.Preview(config, outputDir)
	}

	// Determine output path
	projectPath := filepath.Join(outputDir, config.Name)

	// Generate the project
	options := types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     dryRun,
		NoGit:      noGit,
		Verbose:    cmd.Flag("verbose").Changed,
	}

	// Generate the project with spinner
	var result *types.GenerationResult
	err = spinner.New().
		Title("🚀 Generating your Go project...").
		Action(func() {
			result, err = gen.Generate(config, options)
		}).
		Run()

	if err != nil {
		printErrorMessage("Failed to generate project", err)
		return fmt.Errorf("failed to generate project: %w", err)
	}

	// Print success message
	printSuccessMessage(config, result)
	return nil
}

func validateConfig(cfg types.ProjectConfig) error {
	if cfg.Name == "" {
		return types.NewValidationError("project name is required", nil)
	}
	if cfg.Module == "" {
		return types.NewValidationError("module path is required", nil)
	}
	if cfg.Type == "" {
		return types.NewValidationError("project type is required", nil)
	}

	// Validate logger if provided
	if cfg.Logger != "" {
		if err := config.ValidateLogger(cfg.Logger); err != nil {
			return types.NewValidationError(fmt.Sprintf("invalid logger: %v", err), nil)
		}
	}

	// Validate Go version if provided
	if cfg.GoVersion != "" {
		if err := prompts.ValidateGoVersion(cfg.GoVersion); err != nil {
			return types.NewValidationError(fmt.Sprintf("invalid Go version: %v", err), nil)
		}
	}

	return nil
}

func printSuccessMessage(config types.ProjectConfig, result *types.GenerationResult) {
	// Define beautiful styles
	successStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("8")).
		MarginBottom(1)

	checkStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10"))

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("8"))

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	commandStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("11")).
		Background(lipgloss.Color("0")).
		Padding(0, 1).
		MarginLeft(2)

	// Print success header
	fmt.Println(successStyle.Render("🎉 Project Created Successfully!"))

	// Print project details
	fmt.Println(headerStyle.Render("📋 Project Details"))
	fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Name:") + " " + valueStyle.Render(config.Name))
	fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Type:") + " " + valueStyle.Render(config.Type))

	if config.GoVersion != "" {
		fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Go Version:") + " " + valueStyle.Render(config.GoVersion))
	}
	if config.Framework != "" {
		fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Framework:") + " " + valueStyle.Render(config.Framework))
	}
	if config.Logger != "" {
		fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Logger:") + " " + valueStyle.Render(config.Logger))
	}

	fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Module:") + " " + valueStyle.Render(config.Module))
	fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Files created:") + " " + valueStyle.Render(fmt.Sprintf("%d", len(result.FilesCreated))))

	if !noGit {
		fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Git repository:") + " " + valueStyle.Render("Initialized"))
	}

	fmt.Println(checkStyle.Render("✓") + " " + labelStyle.Render("Duration:") + " " + valueStyle.Render(result.Duration.String()))

	// Print next steps
	fmt.Println()
	fmt.Println(headerStyle.Render("🚀 Next Steps"))

	// Check if Go is available and provide appropriate next steps
	if isGoAvailable() {
		fmt.Println(commandStyle.Render("cd " + config.Name))
		fmt.Println(commandStyle.Render("make run"))
	} else {
		fmt.Println(labelStyle.Render("# Install Go first, then run:"))
		fmt.Println(commandStyle.Render("cd " + config.Name))
		fmt.Println(commandStyle.Render("go mod tidy"))
		fmt.Println(commandStyle.Render("make run"))
	}

	// Add helpful tips
	fmt.Println()
	tipStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Italic(true).
		MarginLeft(2)

	fmt.Println(tipStyle.Render("💡 Tip: Run 'make help' inside your project to see all available commands"))
}

// progressiveHelpFunc provides progressive disclosure for help output
func progressiveHelpFunc(cmd *cobra.Command, args []string) {
	// Parse disclosure mode from command line arguments
	disclosureMode := parseDisclosureModeFromArgs(os.Args)
	
	if disclosureMode == prompts.DisclosureModeBasic {
		printBasicHelp(cmd)
	} else {
		printAdvancedHelp(cmd)
	}
}

// parseDisclosureModeFromArgs extracts disclosure mode from raw command line arguments
func parseDisclosureModeFromArgs(args []string) prompts.DisclosureMode {
	for i, arg := range args {
		switch arg {
		case "--advanced":
			return prompts.DisclosureModeAdvanced
		case "--basic":
			return prompts.DisclosureModeBasic
		case "--complexity":
			// Check next argument for complexity value
			if i+1 < len(args) {
				switch args[i+1] {
				case "advanced", "expert":
					return prompts.DisclosureModeAdvanced
				}
			}
		}
		// Handle --complexity=value format
		if strings.HasPrefix(arg, "--complexity=") {
			value := strings.SplitN(arg, "=", 2)[1]
			switch value {
			case "advanced", "expert":
				return prompts.DisclosureModeAdvanced
			}
		}
	}
	
	// Default to basic mode for new users
	return prompts.DisclosureModeBasic
}

// printBasicHelp renders help with only essential flags visible
func printBasicHelp(cmd *cobra.Command) {
	// Essential flags that beginners need to see
	essentialFlags := map[string]bool{
		"name":       true,
		"type":       true,
		"module":     true,
		"framework":  true,
		"logger":     true,
		"go-version": true,
		"output":     true,
		"complexity": true,
		"basic":      true,
		"advanced":   true,
		"dry-run":    true,
		"help":       true,
		"quiet":      true,
		"no-git":     true,
		"random-name": true,
	}
	
	fmt.Print(buildCustomHelp(cmd, essentialFlags, true))
}

// printAdvancedHelp renders help with all flags visible
func printAdvancedHelp(cmd *cobra.Command) {
	fmt.Print(buildCustomHelp(cmd, nil, false))
}

// buildCustomHelp creates the help text with filtered flags
func buildCustomHelp(cmd *cobra.Command, essentialFlags map[string]bool, isBasic bool) string {
	// Use lipgloss for styling
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("11"))
	flagStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	descStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	hintStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Italic(true)
	
	var help strings.Builder
	
	// Title and description
	help.WriteString(titleStyle.Render(cmd.Short) + "\n\n")
	help.WriteString(cmd.Long + "\n\n")
	
	// Usage
	help.WriteString(sectionStyle.Render("USAGE") + "\n")
	help.WriteString(fmt.Sprintf("  %s [flags]\n", cmd.Use) + "\n")
	
	// Examples
	if cmd.Example != "" {
		help.WriteString(sectionStyle.Render("EXAMPLES") + "\n")
		help.WriteString(cmd.Example + "\n\n")
	}
	
	// Flags section
	help.WriteString(sectionStyle.Render("FLAGS") + "\n")
	
	// Collect and filter flags
	var flagLines []string
	seenFlags := make(map[string]bool)
	
	// Process local flags first
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		// Skip if in basic mode and flag is not essential
		if isBasic && essentialFlags != nil && !essentialFlags[flag.Name] {
			return
		}
		
		flagName := flagStyle.Render("--" + flag.Name)
		if flag.Shorthand != "" {
			flagName = flagStyle.Render("-" + flag.Shorthand + " --" + flag.Name)
		}
		
		// Format flag line with proper spacing
		flagLine := fmt.Sprintf("    %-20s %s", flagName, descStyle.Render(flag.Usage))
		if flag.DefValue != "" && flag.DefValue != "false" {
			flagLine += fmt.Sprintf(" (%s)", flag.DefValue)
		}
		flagLines = append(flagLines, flagLine)
		seenFlags[flag.Name] = true
	})
	
	// Add persistent flags that haven't been seen yet
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		// Skip if already processed or if in basic mode and flag is not essential
		if seenFlags[flag.Name] || (isBasic && essentialFlags != nil && !essentialFlags[flag.Name]) {
			return
		}
		
		flagName := flagStyle.Render("--" + flag.Name)
		if flag.Shorthand != "" {
			flagName = flagStyle.Render("-" + flag.Shorthand + " --" + flag.Name)
		}
		
		flagLine := fmt.Sprintf("    %-20s %s", flagName, descStyle.Render(flag.Usage))
		if flag.DefValue != "" && flag.DefValue != "false" {
			flagLine += fmt.Sprintf(" (%s)", flag.DefValue)
		}
		flagLines = append(flagLines, flagLine)
		seenFlags[flag.Name] = true
	})
	
	// Sort and display flags
	for _, line := range flagLines {
		help.WriteString(line + "\n")
	}
	
	// Add progressive disclosure hint
	if isBasic {
		help.WriteString("\n" + hintStyle.Render("💡 Use --advanced to see all available options") + "\n")
	} else {
		help.WriteString("\n" + hintStyle.Render("💡 Use --basic to see only essential options") + "\n")
	}
	
	return help.String()
}

// isGoAvailable checks if Go is installed and available in PATH
func isGoAvailable() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}

// printErrorMessage prints a beautiful error message using lipgloss styling
func printErrorMessage(title string, err error) {
	// Define error styles
	errorStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("9")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("9")).
		Padding(1, 2).
		MarginTop(1).
		MarginBottom(1)

	iconStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9"))

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("15"))

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		MarginLeft(2)

	// Print error message
	fmt.Println()
	fmt.Println(errorStyle.Render(iconStyle.Render("❌") + " " + titleStyle.Render(title)))
	if err != nil {
		fmt.Println(messageStyle.Render("Error: " + err.Error()))
	}
	fmt.Println()
}
