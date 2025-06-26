package cmd

import (
	"fmt"
	"runtime"

	"github.com/fatih/color"
	"github.com/francknouama/go-starter/internal/ascii"
	"github.com/spf13/cobra"
)

var (
	// Version is set by build flags
	Version = "1.1.1"
	// Commit is set by build flags
	Commit = "none"
	// Date is set by build flags
	Date = "unknown"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  `Display version information for go-starter including version, commit, and build date.`,
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func showVersion() {
	// Color functions
	cyan := color.New(color.FgCyan, color.Bold).SprintFunc()
	green := color.New(color.FgGreen, color.Bold).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()

	// Display logo and version information
	fmt.Print(ascii.Logo())
	fmt.Printf("\n%s %s\n", green("Version:"), cyan(Version))
	fmt.Printf("%s %s\n", yellow("Commit: "), blue(Commit))
	fmt.Printf("%s %s\n", yellow("Built:  "), blue(Date))
	fmt.Printf("%s %s\n", yellow("Go:     "), magenta(runtime.Version()))
	fmt.Printf("%s %s\n", yellow("OS:     "), magenta(runtime.GOOS))
	fmt.Printf("%s %s\n", yellow("Arch:   "), magenta(runtime.GOARCH))

	// Add some decoration
	fmt.Printf("\n%s\n", green("🚀 Ready to generate awesome Go projects!"))
}
