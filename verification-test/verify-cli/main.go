package main

import (
	"fmt"
	"os"

	"github.com/verify/cli/cmd"
	"github.com/verify/cli/internal/config"
	"github.com/verify/cli/internal/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger factory
	loggerFactory := logger.NewFactory()
	
	// Create logger from configuration
	appLogger, err := loggerFactory.CreateFromProjectConfig(
		"slog",
		cfg.Logging.Level,
		cfg.Logging.Format,
		cfg.Logging.Structured,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		os.Exit(1)
	}

	// Execute the root command with logger context
	if err := cmd.Execute(appLogger); err != nil {
		appLogger.ErrorWith("Command execution failed", logger.Fields{"error": err})
		os.Exit(1)
	}
}