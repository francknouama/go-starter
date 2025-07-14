# config Package

This package manages configuration for the go-starter CLI tool, including user preferences, default values, and profile management.

## Overview

The config package provides a robust configuration system that supports multiple profiles, default values, and persistent storage of user preferences.

## Key Components

### Types

- **`Config`** - Main configuration structure
- **`Profile`** - User profile with defaults and preferences
- **`ProjectConfig`** - Project-specific configuration

### Functions

- **`Load() (*Config, error)`** - Load configuration from disk
- **`Save(cfg *Config) error`** - Save configuration to disk
- **`GetDefaultProfile() *Profile`** - Get current default profile
- **`SetProfile(name string) error`** - Switch active profile

## Configuration File Location

- **Unix/Linux**: `~/.config/go-starter/config.yaml`
- **macOS**: `~/Library/Application Support/go-starter/config.yaml`
- **Windows**: `%APPDATA%\go-starter\config.yaml`

## Configuration Structure

```yaml
profiles:
  default:
    author: "John Doe"
    email: "john@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
      architecture: "standard"
  work:
    author: "John Doe"
    email: "john@company.com"
    license: "Proprietary"
    defaults:
      goVersion: "1.21"
      framework: "echo"
      logger: "zap"
current_profile: "default"
```

## Features

- Multiple profile support
- Environment variable overrides
- Validation of configuration values
- Migration support for config updates
- Thread-safe operations

## Usage Example

```go
import "github.com/yourusername/go-starter/internal/config"

// Load configuration
cfg, err := config.Load()
if err != nil {
    log.Fatal(err)
}

// Get current profile
profile := cfg.GetDefaultProfile()
fmt.Printf("Using profile: %s\n", cfg.CurrentProfile)

// Update and save
cfg.CurrentProfile = "work"
err = config.Save(cfg)
```

## Dependencies

- **github.com/spf13/viper** - Configuration management
- **gopkg.in/yaml.v3** - YAML parsing