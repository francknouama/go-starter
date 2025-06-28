package main

import (
	"github.com/francknouama/go-starter/cmd"
	"github.com/francknouama/go-starter/internal/templates"
)

func main() {
	// Initialize the templates filesystem
	templates.SetTemplatesFS(TemplatesFS)
	
	// Execute the CLI
	cmd.Execute()
}
