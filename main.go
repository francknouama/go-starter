package main

import (
	"embed"

	"github.com/francknouama/go-starter/cmd"
	"github.com/francknouama/go-starter/internal/templates"
)

//go:embed all:blueprints
var templatesFS embed.FS

func main() {
	// Initialize the templates filesystem
	templates.SetTemplatesFS(templatesFS)

	// Execute the CLI
	cmd.Execute()
}
