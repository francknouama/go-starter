# cmd Package

This package contains the CLI commands and binary entry points for the go-starter project generator.

## Overview

The cmd package implements the command-line interface using the Cobra framework and provides multiple binary entry points for different use cases. It includes all user-facing commands for generating Go projects with various blueprints and architectures.

## Binary Structure

The project now follows Go best practices for multi-binary projects with separate entry points:

### go-starter (CLI Tool)
- **Location**: `cmd/go-starter/main.go`
- **Purpose**: Main CLI tool for project generation
- **Build**: `make build-cli` or `go build -o bin/go-starter ./cmd/go-starter`
- **Usage**: Primary interface for developers

### go-starter-web (Production Web Server)
- **Location**: `web/cmd/web-server/main.go` (built via `make build-web`)
- **Purpose**: Production web server with embedded assets
- **Build**: `make build-web` (builds from web module)
- **Usage**: Production deployment of web interface

### go-starter-dev (Development Web Server)
- **Location**: `cmd/go-starter-dev/main.go`
- **Purpose**: Development web server with filesystem access
- **Build**: `make build-dev` or `go build -o bin/go-starter-dev ./cmd/go-starter-dev`
- **Usage**: Local development with hot reloading

### Legacy Binary (Deprecated)
- **Location**: `main.go` (root level)
- **Purpose**: Backward compatibility during transition
- **Build**: `make build-legacy` (shows deprecation warning)
- **Status**: Will be removed in future version

## Commands

### Root Command
- **File**: `root.go`
- **Description**: Main entry point for the CLI, sets up the command structure and global flags
- **Usage**: `go-starter [command]`

### New Command
- **File**: `new.go`
- **Description**: Creates a new Go project from available blueprints with progressive disclosure
- **Usage**: `go-starter new [project-name]`
- **Features**:
  - **Progressive Disclosure System** ✨: Smart help filtering and complexity-aware generation
  - **Two-Tier Help**: Basic mode (14 flags) vs Advanced mode (18+ flags)
  - **Complexity Levels**: Simple, Standard, Advanced, Expert project complexity
  - **Interactive Prompts**: Context-aware prompts with smart defaults
  - **Blueprint Selection**: Complexity-driven blueprint mapping (cli → cli-simple)
  - **Architecture Support**: Standard, Clean Architecture, DDD, Hexagonal, Event-driven
  - **Framework Selection**: gin, echo, fiber, chi with context-aware options
  - **Simplified Logger Selection**: slog, zap, logrus, zerolog with 60-90% code reduction

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

## Progressive Disclosure System ✨

The CLI implements a comprehensive progressive disclosure system that adapts the interface based on user experience level and project complexity needs.

### Two-Tier Help System

**Basic Mode (Default)**:
```bash
go-starter new --help
# Shows 14 essential flags: --type, --name, --module, --framework, --logger, etc.
# Includes hint: "💡 Use --advanced to see all available options"
```

**Advanced Mode**:
```bash
go-starter new --advanced --help
# Shows all 18+ flags including database, authentication, deployment options
# Includes hint: "💡 Use --basic to see only essential options"
```

### Complexity-Aware Generation

**Blueprint Selection Logic**:
- `--complexity=simple --type=cli` → Selects `cli-simple` blueprint (8 files)
- `--complexity=standard --type=cli` → Selects `cli` blueprint (29 files)
- Other blueprints remain unchanged by complexity flag

**Smart Defaults Application**:
When `--complexity` and `--type=cli` are specified:
- **Framework**: Automatically set to `cobra` (industry standard)
- **Logger**: Automatically set to `slog` (Go standard library)
- **Module**: Auto-generates `github.com/username/{project-name}` for testing

### Interactive Prompting Prevention

The system prevents unnecessary prompting through:
1. **Pre-Prompt Analysis**: Check if sufficient flags are provided
2. **Default Application**: Apply smart defaults based on blueprint type and complexity  
3. **Prompt Bypass**: Skip interactive prompts when configuration is complete

### Help System Architecture

**Custom Help Function**: 
- Overrides Cobra's built-in help to support dynamic flag filtering
- Eliminates duplicate flags between local and persistent flags
- Provides context-aware hints and styling with lipgloss

**Flag Categorization**:
- **Essential Flags** (Basic Mode): name, type, module, framework, logger, go-version, etc.
- **Advanced Flags** (Advanced Mode): database-driver, auth-type, banner-style, etc.

## Global Flags

### Essential Flags (Always Shown)
- `--config` - Config file path (default: $HOME/.go-starter.yaml)
- `--debug` - Enable debug logging
- `--no-color` - Disable colored output
- `--basic` - Force basic mode help (14 essential flags)
- `--advanced` - Force advanced mode help (18+ flags)
- `--complexity` - Set project complexity level (simple|standard|advanced|expert)

### Advanced Flags (Advanced Mode Only)
- `--database-driver` - Database driver selection
- `--database-orm` - ORM framework selection
- `--auth-type` - Authentication method selection
- `--banner-style` - CLI banner styling options
- `--no-banner` - Disable banner display

## Environment Variables

- `GO_STARTER_CONFIG` - Override default config location
- `GO_STARTER_DEBUG` - Enable debug mode
- `NO_COLOR` - Disable colored output globally

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra) - Command framework
- [spf13/viper](https://github.com/spf13/viper) - Configuration management
- Internal packages for prompts, templates, and generation logic