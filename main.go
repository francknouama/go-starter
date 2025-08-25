package main

import (
	"embed"
	"fmt"
	"io/fs"
	"os"

	"github.com/francknouama/go-starter/cmd"
	"github.com/francknouama/go-starter/pkg/embedfs"
)

//go:embed all:blueprints
var blueprintsFS embed.FS

func main() {
	// Show deprecation warning for root main.go usage
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] == "-h") {
		fmt.Fprintf(os.Stderr, "⚠️  DEPRECATION WARNING: Using the root binary is deprecated.\n")
		fmt.Fprintf(os.Stderr, "   Please use the new binary locations:\n")
		fmt.Fprintf(os.Stderr, "   • CLI tool: go build -o go-starter ./cmd/go-starter\n")
		fmt.Fprintf(os.Stderr, "   • Web server (prod): go build -o go-starter-web ./cmd/go-starter-web\n")
		fmt.Fprintf(os.Stderr, "   • Web server (dev): go build -o go-starter-dev ./cmd/go-starter-dev\n\n")
	}

	// Set the embedded filesystem for blueprints (strip "blueprints" prefix)
	subFS, err := fs.Sub(blueprintsFS, "blueprints")
	if err != nil {
		panic("failed to create sub-filesystem for blueprints: " + err.Error())
	}
	embedfs.SetBlueprintsFS(subFS)

	// Execute the CLI
	cmd.Execute()
}
