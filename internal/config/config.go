package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	// Profiles contain user configuration profiles
	Profiles map[string]Profile `yaml:"profiles" mapstructure:"profiles"`
	// CurrentProfile is the active profile name
	CurrentProfile string `yaml:"current_profile" mapstructure:"current_profile"`
}

// Profile represents a user configuration profile
type Profile struct {
	// Author is the default project author
	Author string `yaml:"author" mapstructure:"author"`
	// Email is the default author email
	Email string `yaml:"email" mapstructure:"email"`
	// License is the default license type
	License string `yaml:"license" mapstructure:"license"`
	// Defaults contains default values for project generation
	Defaults ProfileDefaults `yaml:"defaults" mapstructure:"defaults"`
}

// ProfileDefaults contains default values for project generation
type ProfileDefaults struct {
	// GoVersion is the default Go version to use
	GoVersion string `yaml:"go_version" mapstructure:"go_version"`
	// Framework is the default framework
	Framework string `yaml:"framework" mapstructure:"framework"`
	// Architecture is the default architecture pattern
	Architecture string `yaml:"architecture" mapstructure:"architecture"`
	// Logger is the default logger type
	Logger string `yaml:"logger" mapstructure:"logger"`
	// Database is the default database configuration
	Database DatabaseDefaults `yaml:"database" mapstructure:"database"`
	// Auth is the default authentication configuration
	Auth AuthDefaults `yaml:"auth" mapstructure:"auth"`
	// Logging is the default logging configuration
	Logging LoggingDefaults `yaml:"logging" mapstructure:"logging"`
}

// DatabaseDefaults contains default database configuration
type DatabaseDefaults struct {
	// Driver is the default database driver
	Driver string `yaml:"driver" mapstructure:"driver"`
	// ORM is the default ORM
	ORM string `yaml:"orm" mapstructure:"orm"`
}

// AuthDefaults contains default authentication configuration
type AuthDefaults struct {
	// Type is the default authentication type
	Type string `yaml:"type" mapstructure:"type"`
}

// LoggingDefaults contains default logging configuration
type LoggingDefaults struct {
	// Type is the default logger type
	Type string `yaml:"type" mapstructure:"type"`
	// Level is the default log level
	Level string `yaml:"level" mapstructure:"level"`
	// Format is the default log format
	Format string `yaml:"format" mapstructure:"format"`
	// Structured indicates if structured logging is enabled by default
	Structured bool `yaml:"structured" mapstructure:"structured"`
}

var (
	// DefaultConfig is the default configuration
	DefaultConfig = Config{
		Profiles: map[string]Profile{
			"default": {
				Author:  "",
				Email:   "",
				License: "MIT",
				Defaults: ProfileDefaults{
					GoVersion:    "1.21",
					Framework:    "gin",
					Architecture: "standard",
					Logger:       "slog",
					Database: DatabaseDefaults{
						Driver: "",
						ORM:    "gorm",
					},
					Auth: AuthDefaults{
						Type: "",
					},
					Logging: LoggingDefaults{
						Type:       "slog",
						Level:      "info",
						Format:     "json",
						Structured: true,
					},
				},
			},
		},
		CurrentProfile: "default",
	}
)

// Load loads the configuration from the config file
func Load(configFile string) (*Config, error) {
	v := viper.New()

	// Set config file if provided
	if configFile != "" {
		v.SetConfigFile(configFile)
	} else {
		// Set default config locations
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}

		v.SetConfigName(".go-starter")
		v.SetConfigType("yaml")
		v.AddConfigPath(home)
		v.AddConfigPath(".")
	}

	// Set environment variable prefix
	v.SetEnvPrefix("GO_STARTER")
	v.AutomaticEnv()

	config := &Config{}

	// Try to read config file
	if err := v.ReadInConfig(); err != nil {
		// If config file doesn't exist, use default config
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			*config = DefaultConfig
			return config, nil
		}
		// For other file-not-found errors (like invalid path), also use default
		if os.IsNotExist(err) {
			*config = DefaultConfig
			return config, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Unmarshal into config struct
	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate and apply defaults
	config.applyDefaults()

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves the configuration to the specified file
func (c *Config) Save(configFile string) error {
	// Determine config file path
	var filePath string
	if configFile != "" {
		filePath = configFile
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get user home directory: %w", err)
		}
		filePath = filepath.Join(home, ".go-starter.yaml")
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	v := viper.New()
	v.SetConfigFile(filePath)
	v.SetConfigType("yaml")

	// Set all config values
	v.Set("profiles", c.Profiles)
	v.Set("current_profile", c.CurrentProfile)

	if err := v.WriteConfig(); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetCurrentProfile returns the current active profile
func (c *Config) GetCurrentProfile() (*Profile, error) {
	profile, exists := c.Profiles[c.CurrentProfile]
	if !exists {
		return nil, fmt.Errorf("profile '%s' not found", c.CurrentProfile)
	}
	return &profile, nil
}

// SetCurrentProfile sets the current active profile
func (c *Config) SetCurrentProfile(name string) error {
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}
	c.CurrentProfile = name
	return nil
}

// AddProfile adds a new profile to the configuration
func (c *Config) AddProfile(name string, profile Profile) {
	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}
	c.Profiles[name] = profile
}

// RemoveProfile removes a profile from the configuration
func (c *Config) RemoveProfile(name string) error {
	if name == "default" {
		return fmt.Errorf("cannot remove default profile")
	}
	if _, exists := c.Profiles[name]; !exists {
		return fmt.Errorf("profile '%s' not found", name)
	}
	delete(c.Profiles, name)

	// If we're removing the current profile, switch to default
	if c.CurrentProfile == name {
		c.CurrentProfile = "default"
	}

	return nil
}

// applyDefaults applies default values to the configuration
func (c *Config) applyDefaults() {
	// Ensure we have at least a default profile
	if c.Profiles == nil {
		c.Profiles = make(map[string]Profile)
	}

	if _, exists := c.Profiles["default"]; !exists {
		c.Profiles["default"] = DefaultConfig.Profiles["default"]
	}

	// Set current profile to default if empty
	if c.CurrentProfile == "" {
		c.CurrentProfile = "default"
	}

	// Apply default values to each profile
	for name, profile := range c.Profiles {
		defaultProfile := DefaultConfig.Profiles["default"]

		// Apply defaults if values are empty
		if profile.License == "" {
			profile.License = defaultProfile.License
		}
		if profile.Defaults.GoVersion == "" {
			profile.Defaults.GoVersion = defaultProfile.Defaults.GoVersion
		}
		if profile.Defaults.Framework == "" {
			profile.Defaults.Framework = defaultProfile.Defaults.Framework
		}
		if profile.Defaults.Architecture == "" {
			profile.Defaults.Architecture = defaultProfile.Defaults.Architecture
		}
		if profile.Defaults.Logger == "" {
			profile.Defaults.Logger = defaultProfile.Defaults.Logger
		}
		if profile.Defaults.Database.ORM == "" {
			profile.Defaults.Database.ORM = defaultProfile.Defaults.Database.ORM
		}
		if profile.Defaults.Logging.Type == "" {
			profile.Defaults.Logging.Type = defaultProfile.Defaults.Logging.Type
		}
		if profile.Defaults.Logging.Level == "" {
			profile.Defaults.Logging.Level = defaultProfile.Defaults.Logging.Level
		}
		if profile.Defaults.Logging.Format == "" {
			profile.Defaults.Logging.Format = defaultProfile.Defaults.Logging.Format
		}

		c.Profiles[name] = profile
	}
}

// validate validates the configuration
func (c *Config) validate() error {
	if len(c.Profiles) == 0 {
		return fmt.Errorf("no profiles defined")
	}

	if c.CurrentProfile == "" {
		return fmt.Errorf("current profile not set")
	}

	if _, exists := c.Profiles[c.CurrentProfile]; !exists {
		return fmt.Errorf("current profile '%s' not found", c.CurrentProfile)
	}

	// Validate each profile
	for name, profile := range c.Profiles {
		if err := validateProfile(name, profile); err != nil {
			return fmt.Errorf("invalid profile '%s': %w", name, err)
		}
	}

	return nil
}

// validateProfile validates a single profile
func validateProfile(_ string, profile Profile) error {
	// Validate Go version format
	if profile.Defaults.GoVersion != "" {
		// Simple validation - should be in format "1.xx"
		if len(profile.Defaults.GoVersion) < 3 || profile.Defaults.GoVersion[0:2] != "1." {
			return fmt.Errorf("invalid Go version format: %s", profile.Defaults.GoVersion)
		}
	}

	// Validate framework options
	validFrameworks := map[string]bool{
		"gin":   true,
		"echo":  true,
		"fiber": true,
		"chi":   true,
		"cobra": true,
		"":      true, // empty is allowed
	}
	if !validFrameworks[profile.Defaults.Framework] {
		return fmt.Errorf("invalid framework: %s", profile.Defaults.Framework)
	}

	// Validate architecture options
	validArchitectures := map[string]bool{
		"standard":     true,
		"clean":        true,
		"ddd":          true,
		"hexagonal":    true,
		"event-driven": true,
		"":             true, // empty is allowed
	}
	if !validArchitectures[profile.Defaults.Architecture] {
		return fmt.Errorf("invalid architecture: %s", profile.Defaults.Architecture)
	}

	// Validate database driver options
	validDrivers := map[string]bool{
		"postgres": true,
		"mysql":    true,
		"mongodb":  true,
		"sqlite":   true,
		"redis":    true,
		"":         true, // empty is allowed
	}
	if !validDrivers[profile.Defaults.Database.Driver] {
		return fmt.Errorf("invalid database driver: %s", profile.Defaults.Database.Driver)
	}

	// Validate ORM options
	validORMs := map[string]bool{
		"gorm": true,
		"sqlx": true,
		"sqlc": true,
		"ent":  true,
		"":     true, // empty is allowed
	}
	if !validORMs[profile.Defaults.Database.ORM] {
		return fmt.Errorf("invalid ORM: %s", profile.Defaults.Database.ORM)
	}

	// Validate auth type options
	validAuthTypes := map[string]bool{
		"jwt":     true,
		"oauth2":  true,
		"session": true,
		"api-key": true,
		"":        true, // empty is allowed
	}
	if !validAuthTypes[profile.Defaults.Auth.Type] {
		return fmt.Errorf("invalid auth type: %s", profile.Defaults.Auth.Type)
	}

	// Validate logger options
	validLoggers := map[string]bool{
		"slog":    true,
		"zap":     true,
		"logrus":  true,
		"zerolog": true,
		"":        true, // empty is allowed
	}
	if !validLoggers[profile.Defaults.Logger] {
		return fmt.Errorf("invalid logger: %s", profile.Defaults.Logger)
	}
	if !validLoggers[profile.Defaults.Logging.Type] {
		return fmt.Errorf("invalid logging type: %s", profile.Defaults.Logging.Type)
	}

	// Validate log level options
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"":      true, // empty is allowed
	}
	if !validLogLevels[profile.Defaults.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", profile.Defaults.Logging.Level)
	}

	// Validate log format options
	validLogFormats := map[string]bool{
		"json":    true,
		"text":    true,
		"console": true,
		"":        true, // empty is allowed
	}
	if !validLogFormats[profile.Defaults.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", profile.Defaults.Logging.Format)
	}

	return nil
}
