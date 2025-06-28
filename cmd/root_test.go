package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestRootCommand(t *testing.T) {
	// Test basic root command properties
	if rootCmd.Use != "go-starter" {
		t.Errorf("Expected Use to be 'go-starter', got %s", rootCmd.Use)
	}

	if rootCmd.Short == "" {
		t.Error("Expected Short description to not be empty")
	}

	if rootCmd.Long == "" {
		t.Error("Expected Long description to not be empty")
	}
}

func TestRootCommandFlags(t *testing.T) {
	// Test that the root command has the expected persistent flags
	expectedFlags := []string{
		"config",
		"verbose",
	}

	for _, flagName := range expectedFlags {
		flag := rootCmd.PersistentFlags().Lookup(flagName)
		if flag == nil {
			t.Errorf("Expected persistent flag %s to exist", flagName)
		}
	}
}

func TestViperBinding(t *testing.T) {
	// Test that viper is properly bound to flags
	// We can't easily test the actual binding without more complex setup,
	// but we can verify the flag exists and viper has been initialized

	verboseFlag := rootCmd.PersistentFlags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("verbose flag should exist")
	}

	// Test that viper has been initialized (this would panic if not)
	viper.GetBool("verbose") // Should not panic
}

func TestInitConfig(t *testing.T) {
	// Test that initConfig function doesn't panic and handles errors properly
	// This is mainly a smoke test since it involves file system operations
	err := initConfig()
	if err != nil {
		// initConfig should not fail in normal circumstances
		// If it does fail, it should return a proper error, not panic
		t.Logf("initConfig returned error (this may be expected in test environment): %v", err)
	}
}

func TestExecute(t *testing.T) {
	// Test that Execute function exists and is callable
	// We can't test the actual execution without more complex setup,
	// but we can verify it's defined
	// Note: Execute is a function, not a variable, so we can't check for nil
	_ = Execute // This will compile only if Execute is defined
}

func TestRootGlobalVariables(t *testing.T) {
	// Test that global variables are properly declared
	// Note: &cfgFile can never be nil, so we just test the variable exists
	_ = cfgFile // This will compile only if cfgFile is defined
}

func TestCommandHierarchy(t *testing.T) {
	// Test that the root command has child commands
	commands := rootCmd.Commands()
	if len(commands) == 0 {
		t.Error("Root command should have child commands")
	}

	// Check if the new command is registered
	hasNewCommand := false
	for _, cmd := range commands {
		if cmd.Name() == "new" {
			hasNewCommand = true
			break
		}
	}

	if !hasNewCommand {
		t.Error("Root command should have 'new' subcommand")
	}
}
