package cmd

import (
	"fmt"
	"runtime"

	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/ascii"
	"github.com/spf13/cobra"
)

var (
	// Version is set by build flags
	Version = "1.3.1"
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
	// Define beautiful styles
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("8")).
		MarginBottom(1).
		PaddingLeft(1).
		PaddingRight(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		MarginLeft(2)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	footerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("10")).
		Padding(0, 2).
		MarginTop(2)

	// Display logo and version information
	fmt.Print(ascii.Logo())
	fmt.Println()
	fmt.Println(headerStyle.Render("🚀 Go-Starter Version Information"))
	fmt.Println()

	fmt.Println(labelStyle.Render("Version:") + " " + valueStyle.Render(Version))
	fmt.Println(labelStyle.Render("Commit: ") + " " + valueStyle.Render(Commit))
	fmt.Println(labelStyle.Render("Built:  ") + " " + valueStyle.Render(Date))
	fmt.Println(labelStyle.Render("Go:     ") + " " + valueStyle.Render(runtime.Version()))
	fmt.Println(labelStyle.Render("OS:     ") + " " + valueStyle.Render(runtime.GOOS))
	fmt.Println(labelStyle.Render("Arch:   ") + " " + valueStyle.Render(runtime.GOARCH))

	// Add footer message
	fmt.Println()
	fmt.Println(footerStyle.Render("🎉 Ready to generate awesome Go projects!"))
}
