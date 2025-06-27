# go-starter Quick Reference

## Installation

```bash
# Homebrew
brew tap francknouama/tap && brew install go-starter

# Go Install  
go install github.com/francknouama/go-starter@latest
```

## Basic Commands

```bash
# Interactive mode (recommended)
go-starter new my-project

# Direct mode
go-starter new my-project --type web-api --logger zap

# List available templates
go-starter list

# Show version
go-starter version

# Generate shell completion
go-starter completion bash    # or zsh, fish, powershell

# Configuration commands
go-starter config get         # Show current configuration
go-starter config set key value  # Set configuration value
```

## Project Types

| Type | Command | Use Case |
|------|---------|----------|
| Web API | `--type web-api` | REST APIs, web services |
| CLI | `--type cli` | Command-line tools |
| Library | `--type library` | Reusable packages |
| Lambda | `--type lambda` | AWS serverless functions |

## Logger Options

| Logger | Flag | Best For | Performance |
|--------|------|----------|-------------|
| slog | `--logger slog` | Standard choice | Good |
| zap | `--logger zap` | High-performance APIs | Excellent |
| logrus | `--logger logrus` | Feature-rich apps | Moderate |
| zerolog | `--logger zerolog` | JSON-heavy services | Excellent |

## Common Flags

```bash
--type        # Project type (web-api, cli, library, lambda)
--go-version  # Go version to use (auto, 1.23, 1.22, 1.21)
--logger      # Logger type (slog, zap, logrus, zerolog)
--module      # Module path (e.g., github.com/user/project)
--database    # Database driver (postgres, mysql, mongodb, sqlite, redis)
--force       # Overwrite existing directory
--verbose     # Enable verbose output
--config      # Custom config file location
```

## Go Version Options

| Version | Flag Value | Description |
|---------|------------|-------------|
| Auto-detect | `auto` | Automatically detects the Go version from your system (default) |
| Go 1.23 | `1.23` | Uses Go version 1.23 |
| Go 1.22 | `1.22` | Uses Go version 1.22 |
| Go 1.21 | `1.21` | Uses Go version 1.21 |

## Database Options

| Database | Flag Value | Use Case |
|----------|------------|----------|
| PostgreSQL | `postgres` | ACID-compliant relational DB |
| MySQL | `mysql` | Popular relational DB |
| MongoDB | `mongodb` | Document database |
| SQLite | `sqlite` | File-based, great for dev |
| Redis | `redis` | In-memory cache/store |

## Multiple Databases

```bash
# Interactive mode - select multiple when prompted
go-starter new my-api

# Direct mode - comma-separated
go-starter new my-api --type web-api --database postgres,redis,mongodb
```

## Generated Project Commands

```bash
make help       # Show all available commands
make run        # Start application
make test       # Run tests with coverage  
make lint       # Run linter
make build      # Build production binary
make docker     # Build Docker image
make clean      # Clean build artifacts
```

## Examples

```bash
# High-performance API
go-starter new fast-api --type web-api --logger zap --database postgres

# Feature-rich CLI tool
go-starter new mytool --type cli --logger logrus

# Simple library with standard logger
go-starter new mylib --type library --logger slog

# Serverless function
go-starter new mylambda --type lambda --logger zerolog

# Multi-database microservice
go-starter new user-service --type web-api --logger zap --database postgres,redis
```

## Configuration File

Create `~/.go-starter.yaml`:

```yaml
profiles:
  default:
    author: "Your Name"
    email: "your.email@example.com"
    defaults:
      logger: zap
      goVersion: "1.22"
current_profile: default
```

## Tips

1. **Start simple**: Use interactive mode first
2. **Logger choice**: When in doubt, use `slog`
3. **Performance**: Use `zap` or `zerolog` for high-traffic services
4. **Features**: Use `logrus` when you need hooks and formatters
5. **Regenerate**: Use `--force` to try different configurations
6. **Go version**: go-starter automatically detects your Go version
7. **Databases**: You can select multiple databases for a single project
8. **Shell completion**: Generate completion scripts for faster CLI usage

---

**Need help?** Run `go-starter help` or visit [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)