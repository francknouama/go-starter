package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/verify/cli/internal/logger"
)

var (
	// These will be set by build flags
	version   = "dev"
	commit    = "none"
	date      = "unknown"
	goVersion = runtime.Version()
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version information",
	Long:  `Print the version information for verify-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
		appLogger.InfoWith("Version command executed", logger.Fields{
			"version": version,
			"commit":  commit,
			"date":    date,
		})
		
		fmt.Printf("verify-cli version information:\n")
		fmt.Printf("  Version:    %s\n", version)
		fmt.Printf("  Commit:     %s\n", commit)
		fmt.Printf("  Build Date: %s\n", date)
		fmt.Printf("  Go Version: %s\n", goVersion)
		fmt.Printf("  Logger:     slog\n")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}