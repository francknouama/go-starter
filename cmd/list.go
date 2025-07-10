package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available project blueprints",
	Long: `Display all available project blueprints with their descriptions.

This command shows all blueprints that can be used to generate new projects,
including their type, architecture, and a brief description.`,
	Run: func(cmd *cobra.Command, args []string) {
		listBlueprints()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listBlueprints() {
	registry := templates.NewRegistry()
	blueprintList := registry.List()

	if len(blueprintList) == 0 {
		noTemplatesStyle := lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("9")).
			MarginLeft(2)

		fmt.Println(noTemplatesStyle.Render("No blueprints available yet."))
		fmt.Println(noTemplatesStyle.Render("Blueprints will be added in upcoming releases."))
		return
	}

	// Create styled header
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("12")).
		BorderStyle(lipgloss.NormalBorder()).
		BorderBottom(true).
		BorderForeground(lipgloss.Color("8")).
		MarginBottom(1).
		PaddingLeft(1).
		PaddingRight(1)

	fmt.Println(headerStyle.Render("🚀 Available Go Project Blueprints"))
	fmt.Println()

	// Create table rows
	for i, blueprint := range blueprintList {
		renderBlueprint(blueprint, i == len(blueprintList)-1)
	}

	// Add total count
	totalStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		MarginLeft(2).
		MarginTop(2)

	fmt.Println(totalStyle.Render(fmt.Sprintf("📊 Total: %d blueprint(s) available", len(blueprintList))))
}

func renderBlueprint(blueprint types.Template, isLast bool) {
	// Define styles
	idStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("10")).
		MarginLeft(2)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("8")).
		MarginLeft(4)

	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		MarginLeft(1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("7")).
		MarginLeft(4).
		Width(60)

	// Print template information
	fmt.Println(idStyle.Render("● " + blueprint.ID))
	fmt.Println(labelStyle.Render("Name:") + valueStyle.Render(blueprint.Name))
	fmt.Println(labelStyle.Render("Type:") + valueStyle.Render(blueprint.Type))

	if blueprint.Architecture != "" {
		fmt.Println(labelStyle.Render("Architecture:") + valueStyle.Render(blueprint.Architecture))
	}

	// Wrap description text
	wrappedDesc := wrapText(blueprint.Description, 60)
	fmt.Println(labelStyle.Render("Description:"))
	for _, line := range strings.Split(wrappedDesc, "\n") {
		if strings.TrimSpace(line) != "" {
			fmt.Println(descStyle.Render(line))
		}
	}

	if !isLast {
		// Add separator
		separatorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("8")).
			MarginLeft(2).
			MarginTop(1).
			MarginBottom(1)
		fmt.Println(separatorStyle.Render("─────────────────────────────────────────────────────"))
	}
}

func wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}

	var result strings.Builder
	lineLength := 0

	for _, word := range words {
		if lineLength+len(word)+1 > width && lineLength > 0 {
			result.WriteString("\n")
			lineLength = 0
		}

		if lineLength > 0 {
			result.WriteString(" ")
			lineLength++
		}

		result.WriteString(word)
		lineLength += len(word)
	}

	return result.String()
}
