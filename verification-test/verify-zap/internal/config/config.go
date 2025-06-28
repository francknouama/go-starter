package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Environment string         `mapstructure:"environment"`
	Server      ServerConfig   `mapstructure:"server"`
	Logging     LoggingConfig  `mapstructure:"logging"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int `mapstructure:"port"`
	ReadTimeout  int `mapstructure:"read_timeout"`
	WriteTimeout int `mapstructure:"write_timeout"`
	IdleTimeout  int `mapstructure:"idle_timeout"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Structured bool   `mapstructure:"structured"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set default values
	setDefaults(v)

	// Set config file name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// Set environment variable prefix
	v.SetEnvPrefix("VERIFY-ZAP")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Determine environment
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	// Try to read environment-specific config file
	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.ReadInConfig(); err != nil {
		// If environment-specific config doesn't exist, try default
		v.SetConfigName("config")
		if err := v.ReadInConfig(); err != nil {
			// If no config file found, use defaults with env vars
			if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
				return nil, fmt.Errorf("failed to read config file: %w", err)
			}
		}
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Set environment
	config.Environment = env

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	// Server defaults
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)
	v.SetDefault("server.idle_timeout", 60)

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.structured", true)
}

// validateConfig validates the configuration
func validateConfig(config *Config) error {
	// Validate server port
	if config.Server.Port <= 0 || config.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", config.Server.Port)
	}

	// Validate logging level
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}
	if !validLevels[config.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s", config.Logging.Level)
	}

	return nil
}