package cmd

import (
	"fmt"

	"github.com/francknouama/go-starter/internal/templates"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available project templates",
	Long: `Display all available project templates with their descriptions.

This command shows all templates that can be used to generate new projects,
including their type, architecture, and a brief description.`,
	Run: func(cmd *cobra.Command, args []string) {
		listTemplates()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func listTemplates() {
	registry := templates.NewRegistry()
	templateList := registry.List()

	if len(templateList) == 0 {
		fmt.Println("No templates available yet.")
		fmt.Println("Templates will be added in upcoming releases.")
		return
	}

	fmt.Println("Available templates:")
	fmt.Println()

	for _, template := range templateList {
		fmt.Printf("  %s\n", template.ID)
		fmt.Printf("    Name: %s\n", template.Name)
		fmt.Printf("    Type: %s\n", template.Type)
		if template.Architecture != "" {
			fmt.Printf("    Architecture: %s\n", template.Architecture)
		}
		fmt.Printf("    Description: %s\n", template.Description)
		fmt.Println()
	}

	fmt.Printf("Total: %d template(s)\n", len(templateList))
}
