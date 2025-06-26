package cmd

import (
	"fmt"
	"os/exec"
	"path/filepath"

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
	prompter := prompts.New()

	// Get project configuration through interactive prompts or flags
	config, err := prompter.GetProjectConfig(types.ProjectConfig{
		Name:      projectName,
		Module:    projectModule,
		Type:      projectType,
		Framework: framework,
		Logger:    logger,
	}, advanced)
	if err != nil {
		return fmt.Errorf("failed to get project configuration: %w", err)
	}

	// Validate the configuration
	if err := validateConfig(config); err != nil {
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

	result, err := gen.Generate(config, options)
	if err != nil {
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

	return nil
}

func printSuccessMessage(config types.ProjectConfig, result *types.GenerationResult) {
	fmt.Printf("✓ Project '%s' created successfully!\n", config.Name)
	fmt.Printf("✓ Type: %s\n", config.Type)
	if config.Framework != "" {
		fmt.Printf("✓ Framework: %s\n", config.Framework)
	}
	if config.Logger != "" {
		fmt.Printf("✓ Logger: %s\n", config.Logger)
	}
	fmt.Printf("✓ Go module: %s\n", config.Module)
	fmt.Printf("✓ Files created: %d\n", len(result.FilesCreated))
	if !noGit {
		fmt.Printf("✓ Git repository initialized\n")
	}
	fmt.Printf("✓ Generation completed in %v\n", result.Duration)

	fmt.Printf("\nGet started:\n")
	fmt.Printf("  cd %s\n", config.Name)

	// Check if Go is available and provide appropriate next steps
	if isGoAvailable() {
		fmt.Printf("  make run\n")
	} else {
		fmt.Printf("  # Install Go first, then run:\n")
		fmt.Printf("  go mod tidy\n")
		fmt.Printf("  make run\n")
	}
}

// isGoAvailable checks if Go is installed and available in PATH
func isGoAvailable() bool {
	cmd := exec.Command("go", "version")
	return cmd.Run() == nil
}
