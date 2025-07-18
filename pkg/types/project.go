package types

import "time"

// ProjectConfig represents the configuration for generating a Go project
type ProjectConfig struct {
	Name         string            `yaml:"name" json:"name"`
	Module       string            `yaml:"module" json:"module"`
	Type         string            `yaml:"type" json:"type"`
	GoVersion    string            `yaml:"go_version" json:"go_version"`
	Architecture string            `yaml:"architecture" json:"architecture"`
	Framework    string            `yaml:"framework" json:"framework"`
	Logger       string            `yaml:"logger" json:"logger"`
	Author       string            `yaml:"author" json:"author"`
	Email        string            `yaml:"email" json:"email"`
	License      string            `yaml:"license" json:"license"`
	Features     *Features         `yaml:"features" json:"features"`
	Variables    map[string]string `yaml:"variables" json:"variables"`
}

// Features represents optional features for the project
type Features struct {
	Database       DatabaseConfig `yaml:"database" json:"database"`
	Authentication AuthConfig     `yaml:"authentication" json:"authentication"`
	Deployment     DeployConfig   `yaml:"deployment" json:"deployment"`
	Testing        TestConfig     `yaml:"testing" json:"testing"`
	Monitoring     MonitorConfig  `yaml:"monitoring" json:"monitoring"`
	Logging        LoggingConfig  `yaml:"logging" json:"logging"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Drivers []string `yaml:"drivers" json:"drivers"` // Support multiple databases
	ORM     string   `yaml:"orm" json:"orm"`

	// Deprecated: use Drivers instead. Kept for backward compatibility
	Driver string `yaml:"driver,omitempty" json:"driver,omitempty"`
}

// HasDatabase returns true if any database is configured
func (dc *DatabaseConfig) HasDatabase() bool {
	return len(dc.Drivers) > 0 || dc.Driver != ""
}

// GetDrivers returns all configured database drivers
func (dc *DatabaseConfig) GetDrivers() []string {
	if len(dc.Drivers) > 0 {
		return dc.Drivers
	}
	if dc.Driver != "" {
		return []string{dc.Driver}
	}
	return []string{}
}

// HasDriver returns true if the specified driver is configured
func (dc *DatabaseConfig) HasDriver(driver string) bool {
	drivers := dc.GetDrivers()
	for _, d := range drivers {
		if d == driver {
			return true
		}
	}
	return false
}

// PrimaryDriver returns the first/primary database driver
func (dc *DatabaseConfig) PrimaryDriver() string {
	drivers := dc.GetDrivers()
	if len(drivers) > 0 {
		return drivers[0]
	}
	return ""
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Type      string   `yaml:"type" json:"type"`
	Providers []string `yaml:"providers" json:"providers"`
}

// DeployConfig represents deployment configuration
type DeployConfig struct {
	Targets []string `yaml:"targets" json:"targets"`
}

// TestConfig represents testing configuration
type TestConfig struct {
	Framework string `yaml:"framework" json:"framework"`
	Coverage  bool   `yaml:"coverage" json:"coverage"`
}

// MonitorConfig represents monitoring configuration
type MonitorConfig struct {
	Metrics bool `yaml:"metrics" json:"metrics"`
	Logging bool `yaml:"logging" json:"logging"`
	Tracing bool `yaml:"tracing" json:"tracing"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Type       string `yaml:"type" json:"type"`             // slog, zap, logrus, zerolog
	Level      string `yaml:"level" json:"level"`           // debug, info, warn, error
	Format     string `yaml:"format" json:"format"`         // json, text, console
	Structured bool   `yaml:"structured" json:"structured"` // structured logging enabled
}

// GenerationOptions represents options for the generation process
type GenerationOptions struct {
	OutputPath string
	DryRun     bool
	NoGit      bool
	Verbose    bool
}

// GenerationResult represents the result of a project generation
type GenerationResult struct {
	ProjectPath  string
	FilesCreated []string
	Duration     time.Duration
	Success      bool
	Error        error
}
