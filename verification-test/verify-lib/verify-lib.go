// Package verify_lib provides verify-lib functionality.
package verify_lib

import (
	"context"
	"fmt"

	"github.com/verify/lib/internal/logger"
)

// Client represents a verify-lib client
type Client struct {
	logger logger.Logger
	config *Config
}

// Config holds configuration for the verify-lib client
type Config struct {
	// Logger configuration
	Logger struct {
		Level  string `json:"level"`
		Format string `json:"format"`
	} `json:"logger"`
	
	// Add your configuration fields here
	Debug bool `json:"debug"`
}

// DefaultConfig returns a default configuration
func DefaultConfig() *Config {
	return &Config{
		Debug: false,
	}
}

// New creates a new verify-lib client with the given configuration
func New(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// Initialize internal logger
	loggerFactory := logger.NewFactory()
	internalLogger, err := loggerFactory.CreateFromProjectConfig(
		"slog",
		getLogLevel(config),
		getLogFormat(config),
		true,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	client := &Client{
		logger: internalLogger,
		config: config,
	}

	client.logger.InfoWith("verify-lib client initialized", logger.Fields{
		"logger": "slog",
		"debug":  config.Debug,
	})

	return client, nil
}

// Process demonstrates the main functionality of the library
func (c *Client) Process(ctx context.Context, input string) (string, error) {
	c.logger.DebugWith("Processing input", logger.Fields{
		"input_length": len(input),
	})

	if input == "" {
		c.logger.Warn("Empty input provided")
		return "", fmt.Errorf("input cannot be empty")
	}

	// Add your processing logic here
	result := fmt.Sprintf("Processed: %s", input)

	c.logger.InfoWith("Processing completed", logger.Fields{
		"result_length": len(result),
	})

	return result, nil
}

// Close gracefully shuts down the client
func (c *Client) Close() error {
	c.logger.Info("Closing verify-lib client")
	
	// Sync logger before closing
	if err := c.logger.Sync(); err != nil {
		c.logger.ErrorWith("Failed to sync logger", logger.Fields{"error": err})
		return err
	}
	
	return nil
}

// Helper functions

func getLogLevel(config *Config) string {
	if config.Logger.Level != "" {
		return config.Logger.Level
	}
	if config.Debug {
		return "debug"
	}
	return "warn" // Libraries should be quiet by default
}

func getLogFormat(config *Config) string {
	if config.Logger.Format != "" {
		return config.Logger.Format
	}
	return "text" // Human-readable for library debugging
}