package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/config"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/prompts"
	"github.com/francknouama/go-starter/internal/utils"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/spf13/cobra"
)

var (
	projectName   string
	projectModule string
	projectType   string
	goVersion     string
	outputDir     string
	framework     string
	logger        string
	advanced      bool
	dryRun        bool
	noGit         bool
	randomName    bool
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new Go project",
	Long: `Create a new Go project with the specified template and configuration.

Examples:
  go-starter new my-api                                          # Interactive mode
  go-starter new my-api --type=web-api --framework=gin           # Direct mode
  go-starter new my-api --type=web-api --logger=zap              # With specific logger
  go-starter new my-cli --type=cli --logger=slog                 # CLI application
  go-starter new my-lib --type=library                           # Go library
  go-starter new --random-name --type=web-api                    # Generate random project name
  go-starter new --random-name                                   # Fully interactive with random name

The command will guide you through the project configuration process
or use the provided flags for direct project generation.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runNew,
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Project configuration flags
	newCmd.Flags().StringVar(&projectName, "name", "", "Project name")
	newCmd.Flags().StringVar(&projectModule, "module", "", "Go module path (e.g., github.com/user/project)")
	newCmd.Flags().StringVar(&projectType, "type", "", "Project type (web-api, cli, library, lambda)")
	newCmd.Flags().StringVar(&goVersion, "go-version", "", "Go version to use (auto, 1.23, 1.22, 1.21)")
	newCmd.Flags().StringVar(&framework, "framework", "", "Framework to use (gin, echo, cobra, etc.)")
	newCmd.Flags().StringVar(&logger, "logger", "", "Logger to use (slog, zap, logrus, zerolog)")
	newCmd.Flags().StringVar(&outputDir, "output", ".", "Output directory")

	// Generation options
	newCmd.Flags().BoolVar(&advanced, "advanced", false, "Enable advanced configuration mode")
	newCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Preview project structure without creating files")
	newCmd.Flags().BoolVar(&noGit, "no-git", false, "Skip git repository initialization")
	newCmd.Flags().BoolVar(&randomName, "random-name", false, "Generate a random project name (GitHub-style)")
}

func runNew(cmd *cobra.Command, args []string) error {
	// Get project name from args if provided
	if len(args) > 0 {
		projectName = args[0]
	}

	// Generate random name if requested and no name provided
	if randomName && projectName == "" {
		projectName = utils.GenerateRandomProjectName()
		fmt.Printf("🎲 Generated random project name: %s\n", projectName)
	}

	// Initialize the prompter for interactive configuration
	// Use Fang UI by default, set to prompts.NewSurvey() for fallback
	prompter := prompts.New()

	// Get project configuration through interactive prompts or flags
	config, err := prompter.GetProjectConfig(types.ProjectConfig{
		Name:      projectName,
		Module:    projectModule,
		Type:      projectType,
		GoVersion: goVersion,
		Framework: framework,
		Logger:    logger,
	}, advanced)
	if err != nil {
		printErrorMessage("Failed to get project configuration", err)
		return fmt.Errorf("failed to get project configuration: %w", err)
	}

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
