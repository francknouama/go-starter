package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	// Test default configuration
	config, err := Load()
	require.NoError(t, err)
	
	// Verify defaults
	assert.Equal(t, "development", config.Environment)
	assert.Equal(t, "info", config.Logging.Level)
	assert.Equal(t, "text", config.Logging.Format)
	assert.False(t, config.Logging.Structured)
	assert.Equal(t, "text", config.CLI.OutputFormat)
	assert.False(t, config.CLI.NoColor)
	assert.False(t, config.CLI.Quiet)
}

func TestLoadWithEnvironmentVariables(t *testing.T) {
	// Set environment variables
	os.Setenv("VERIFY-CLI_LOGGING_LEVEL", "debug")
	os.Setenv("VERIFY-CLI_LOGGING_FORMAT", "json")
	os.Setenv("VERIFY-CLI_CLI_OUTPUT_FORMAT", "json")
	defer func() {
		os.Unsetenv("VERIFY-CLI_LOGGING_LEVEL")
		os.Unsetenv("VERIFY-CLI_LOGGING_FORMAT")
		os.Unsetenv("VERIFY-CLI_CLI_OUTPUT_FORMAT")
	}()

	config, err := Load()
	require.NoError(t, err)

	// Verify environment variables are applied
	assert.Equal(t, "debug", config.Logging.Level)
	assert.Equal(t, "json", config.Logging.Format)
	assert.Equal(t, "json", config.CLI.OutputFormat)
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Logging: LoggingConfig{
					Level:  "info",
					Format: "text",
				},
				CLI: CLIConfig{
					OutputFormat: "text",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid logging level",
			config: &Config{
				Logging: LoggingConfig{
					Level:  "invalid",
					Format: "text",
				},
				CLI: CLIConfig{
					OutputFormat: "text",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid logging format",
			config: &Config{
				Logging: LoggingConfig{
					Level:  "info",
					Format: "invalid",
				},
				CLI: CLIConfig{
					OutputFormat: "text",
				},
			},
			wantErr: true,
		},
		{
			name: "invalid CLI output format",
			config: &Config{
				Logging: LoggingConfig{
					Level:  "info",
					Format: "text",
				},
				CLI: CLIConfig{
					OutputFormat: "invalid",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateConfig(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}