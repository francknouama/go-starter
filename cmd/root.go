package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/francknouama/go-starter/internal/ascii"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-starter",
	Short: "Generate Go project structures with best practices",
	Long:  buildLongDescription(),
	Version: Version,
}

// buildLongDescription creates a colorized long description with ASCII art
func buildLongDescription() string {
	// Color functions
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow, color.Bold).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	return ascii.Banner() + "\n" +
		cyan("🚀 A comprehensive Go project generator") + "\n\n" +
		"Generate Go project structures with modern best practices,\n" +
		"multiple architecture patterns, and deployment configurations.\n\n" +
		yellow("📚 EXAMPLES:") + "\n" +
		green("  go-starter new my-api                    ") + "# Interactive project creation\n" +
		green("  go-starter new my-api --type=web-api     ") + "# Direct project creation\n" +
		green("  go-starter list                          ") + "# List available templates\n" +
		green("  go-starter version                       ") + "# Show version information\n\n" +
		yellow("🏗️  SUPPORTED TEMPLATES:") + "\n" +
		blue("  • web-api       ") + "- REST APIs with multiple architectures\n" +
		blue("  • cli           ") + "- Command-line applications\n" +
		blue("  • library       ") + "- Reusable Go packages\n" +
		blue("  • lambda        ") + "- AWS Lambda functions\n" +
		blue("  • microservice  ") + "- Distributed systems\n" +
		blue("  • monolith      ") + "- Traditional web applications\n\n" +
		yellow("🎨 ARCHITECTURE PATTERNS:") + "\n" +
		magenta("  • Standard      ") + "- Simple, straightforward structure\n" +
		magenta("  • Clean         ") + "- Clean Architecture principles\n" +
		magenta("  • DDD           ") + "- Domain-Driven Design\n" +
		magenta("  • Hexagonal     ") + "- Ports and Adapters pattern\n\n" +
		"For more information, visit: " + cyan("https://github.com/francknouama/go-starter")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go-starter.yaml)")
	rootCmd.PersistentFlags().Bool("verbose", false, "verbose output")

	// Bind flags to viper
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		// Search config in home directory with name ".go-starter" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".go-starter")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}
