# cmd Package

This package contains the CLI commands for the go-starter project generator.

## Overview

The cmd package implements the command-line interface using the Cobra framework. It provides all user-facing commands for generating Go projects with various blueprints and architectures.

## Commands

### Root Command
- **File**: `root.go`
- **Description**: Main entry point for the CLI, sets up the command structure and global flags
- **Usage**: `go-starter [command]`

### New Command
- **File**: `new.go`
- **Description**: Creates a new Go project from available blueprints
- **Usage**: `go-starter new [project-name]`
- **Features**:
  - Interactive prompts for project configuration
  - Support for multiple architectures (standard, clean, DDD, hexagonal)
  - Framework selection (gin, echo, fiber, chi)
  - Logger selection (slog, zap, logrus, zerolog)

### List Command
- **File**: `list.go`
- **Description**: Lists all available project blueprints
- **Usage**: `go-starter list`
- **Output**: Displays blueprint names, descriptions, and supported features

### Security Command
- **File**: `security.go`
- **Description**: Security-related operations and checks
- **Usage**: `go-starter security [subcommand]`
- **Features**: Input validation, template security checks

### Version Command
- **File**: `version.go`
- **Description**: Displays the version information
- **Usage**: `go-starter version`
- **Output**: Shows version, build date, and commit hash

## Testing

Each command has corresponding test files:
- `root_test.go` - Tests for root command setup
- `new_test.go` - Tests for project creation workflow
- `list_test.go` - Tests for blueprint listing
- `security_test.go` - Tests for security features
- `version_test.go` - Tests for version display

## Command Structure

```
go-starter
├── new          # Create a new project
├── list         # List available blueprints
├── security     # Security operations
├── version      # Show version info
└── help         # Show help for any command
```

## Global Flags

- `--config` - Config file path (default: $HOME/.go-starter.yaml)
- `--debug` - Enable debug logging
- `--no-color` - Disable colored output

## Environment Variables

- `GO_STARTER_CONFIG` - Override default config location
- `GO_STARTER_DEBUG` - Enable debug mode
- `NO_COLOR` - Disable colored output globally

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - Command framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
- Internal packages for prompts, templates, and generation logic