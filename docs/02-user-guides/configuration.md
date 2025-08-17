# Configuration Guide

This guide covers how to configure go-starter for team environments, standardize project generation, and customize default settings.

## Table of Contents
- [Configuration Files](#configuration-files)
- [Team Configuration](#team-configuration)
- [Default Settings](#default-settings)
- [Environment Variables](#environment-variables)
- [Blueprint Customization](#blueprint-customization)

## Configuration Files

go-starter supports configuration at multiple levels:

### 1. Global Configuration
Located at `~/.go-starter/config.yaml`:

```yaml
# Global configuration
profiles:
  default:
    author: "John Doe"
    email: "john@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
```

### 2. Project Configuration
Located in project root as `.go-starter.yaml`:

```yaml
# Project-specific overrides
project:
  type: "web-api"
  architecture: "clean"
  defaults:
    framework: "gin"
    logger: "zap"
    database:
      driver: "postgres"
      orm: "gorm"
```

### 3. Team Configuration
Share configuration via Git:

```yaml
# team-config.yaml - commit to your repo
team:
  standards:
    goVersion: "1.21"
    linter: "golangci-lint"
    testFramework: "testify"
  
  blueprints:
    preferred:
      - "web-api-clean"
      - "cli-standard"
    forbidden:
      - "monolith"
  
  naming:
    projectPrefix: "svc-"
    modulePrefix: "github.com/yourorg"
```

## Team Configuration

### Setting Up Team Standards

1. **Create Team Configuration File**:
```bash
# In your team's shared repository
touch .go-starter/team-config.yaml
```

2. **Define Standards**:
```yaml
standards:
  # Go version requirement
  goVersion: "1.21"
  
  # Required tools
  tools:
    - golangci-lint
    - gosec
    - gofumpt
  
  # Project structure
  structure:
    testsLocation: "tests/"
    docsLocation: "docs/"
    scriptsLocation: "scripts/"
  
  # Code style
  style:
    lineLength: 120
    importGroups: 
      - standard
      - third-party
      - company
      - project
```

3. **Apply Team Configuration**:
```bash
# Import team configuration
go-starter config import https://github.com/yourorg/standards/.go-starter/team-config.yaml

# Or from local file
go-starter config import ./team-config.yaml
```

### Enforcing Standards

```yaml
# Enforcement rules
enforcement:
  required:
    - "README.md must exist"
    - "LICENSE file required"
    - "CI/CD configuration"
  
  blueprints:
    allowed:
      - pattern: "web-api-*"
        reason: "Only web APIs allowed"
    
    forbidden:
      - pattern: "monolith"
        reason: "Monoliths not allowed per architecture decision"
  
  validation:
    modulePrefix: "github.com/yourorg"
    projectNamePattern: "^[a-z][-a-z0-9]*$"
```

## Default Settings

### Setting User Defaults

```bash
# Set default author
go-starter config set author "John Doe"

# Set default organization
go-starter config set organization "Acme Corp"

# Set default Go version
go-starter config set defaults.goVersion "1.21"

# Set default framework
go-starter config set defaults.framework "gin"

# Set default complexity for CLI projects
go-starter config set defaults.cli.complexity "simple"
```

## Profile Management

### Multiple Profiles
```yaml
# ~/.go-starter.yaml
profiles:
  default:
    author: "John Doe"
    email: "john@example.com"
    license: "MIT"
  work:
    author: "John Doe"
    email: "john@company.com"
    license: "Apache-2.0"
  personal:
    author: "John Doe"
    email: "john@personal.com"
    license: "MIT"
current_profile: "default"
```

### Switching Profiles
```bash
# Use specific profile
go-starter new my-project --profile work

# Set default profile
go-starter config set-profile work
```

## Advanced Configuration

### Project Defaults
```yaml
# ~/.go-starter.yaml
profiles:
  default:
    defaults:
      goVersion: "1.23"
      framework: "gin"
      logger: "zap"
      database: "postgres"
      authentication: "jwt"
      enableDocker: true
      enableCICD: true
```

### Template Customization
```yaml
# ~/.go-starter.yaml
templates:
  customPath: "/path/to/custom/templates"
  overrides:
    web-api: "custom-web-api"
    cli: "custom-cli"
```

## Environment Variables

go-starter respects the following environment variables:

```bash
# Override config file location
export GO_STARTER_CONFIG="$HOME/.config/go-starter/config.yaml"

# Set default profile
export GO_STARTER_PROFILE="work"

# Disable interactive mode
export GO_STARTER_NON_INTERACTIVE="true"

# Set default output directory
export GO_STARTER_OUTPUT_DIR="$HOME/projects"

# Enable debug logging
export GO_STARTER_DEBUG="true"

# Skip git initialization
export GO_STARTER_SKIP_GIT="true"

# Custom templates directory
export GO_STARTER_TEMPLATES_DIR="$HOME/my-templates"

# Override specific values
export GO_STARTER_AUTHOR="Jane Doe"
export GO_STARTER_EMAIL="jane@example.com"
export GO_STARTER_LICENSE="Apache-2.0"
export GO_STARTER_GO_VERSION="1.23"
```

### Docker Environment

```dockerfile
# Dockerfile example
FROM golang:1.21-alpine

# Install go-starter
RUN go install github.com/francknouama/go-starter@latest

# Configure defaults
ENV GO_STARTER_NON_INTERACTIVE=true
ENV GO_STARTER_SKIP_GIT=true

# Set team configuration
COPY .go-starter/team-config.yaml /root/.go-starter/config.yaml
```

## Blueprint Customization

### Custom Blueprint Directory

```bash
# Set custom blueprint directory
export GO_STARTER_CUSTOM_BLUEPRINTS="$HOME/my-blueprints"

# Or in config
go-starter config set customBlueprintsDir "$HOME/my-blueprints"
```

### Blueprint Override

Create custom versions of existing blueprints:

```
my-blueprints/
├── web-api-standard/      # Overrides built-in web-api-standard
│   ├── template.yaml
│   └── custom-files/
└── my-custom-blueprint/   # New custom blueprint
    ├── template.yaml
    └── files/
```

### Team Blueprint Repository

```yaml
# team-config.yaml
blueprints:
  repositories:
    - url: "https://github.com/yourorg/go-blueprints"
      prefix: "org-"
    - url: "https://github.com/community/blueprints"
      prefix: "community-"
  
  # Override built-in blueprints
  overrides:
    "web-api-standard": "org-web-api"
    "cli-standard": "org-cli"
```

## Advanced Configuration

### Hooks and Scripts

```yaml
# .go-starter.yaml
hooks:
  pre-generate:
    - script: "./scripts/validate-env.sh"
      description: "Validate environment"
  
  post-generate:
    - script: "./scripts/setup-git-hooks.sh"
      description: "Install git hooks"
    - command: "make setup"
      description: "Run initial setup"
```

### Conditional Configuration

```yaml
# Conditional defaults based on project type
conditionals:
  - when:
      projectType: "web-api"
    then:
      defaults:
        framework: "gin"
        database:
          driver: "postgres"
  
  - when:
      projectType: "cli"
      complexity: "simple"
    then:
      defaults:
        framework: "cobra"
        logger: "slog"
```

### Security Configuration

```yaml
security:
  # Require security scan
  requireSecurityScan: true
  
  # Allowed licenses
  allowedLicenses:
    - "MIT"
    - "Apache-2.0"
    - "BSD-3-Clause"
  
  # Required files
  requiredFiles:
    - "SECURITY.md"
    - ".gitignore"
    - "LICENSE"
```

## Advanced Mode

Enable advanced configuration options for complex projects:

```bash
go-starter new my-project --advanced
```

Advanced mode includes:
- **Database selection**: PostgreSQL, MySQL, MongoDB, SQLite, Redis
- **Authentication methods**: JWT, OAuth2, API Key, Session
- **Message queues**: RabbitMQ, Kafka, Redis Streams
- **Observability**: Prometheus metrics, Jaeger tracing, OpenTelemetry
- **Deployment platforms**: Docker, Kubernetes, AWS Lambda, Google Cloud Run

## Configuration Commands

```bash
# View current configuration
go-starter config show

# Set configuration values
go-starter config set author "Jane Doe"
go-starter config set email "jane@example.com"
go-starter config set license "Apache-2.0"

# Reset to defaults
go-starter config reset

# Export configuration
go-starter config export > my-config.yaml

# Import configuration
go-starter config import my-config.yaml
```

## Project-Specific Configuration

Override global settings per project:

```bash
# Use different settings for this project only
go-starter new my-project \
  --author "Special Author" \
  --email "special@example.com" \
  --license "GPL-3.0" \
  --go-version "1.22"
```

## Configuration Best Practices

### 1. **Version Control Your Configuration**
```bash
# Add to your dotfiles repo
cd ~/dotfiles
mkdir -p .go-starter
cp ~/.go-starter/config.yaml .go-starter/
git add .go-starter/config.yaml
git commit -m "Add go-starter configuration"
```

### 2. **Use Profiles for Different Contexts**
- `default` - Personal projects
- `work` - Work projects
- `oss` - Open source projects
- `client-x` - Client-specific settings

### 3. **Document Team Standards**
```markdown
# Team go-starter Standards

## Required Configuration
All team members must import the team configuration:
`go-starter config import https://github.com/ourorg/standards/go-starter.yaml`

## Approved Blueprints
- `web-api-clean` - For all REST APIs
- `cli-standard` - For CLI tools
- `library-standard` - For shared libraries

## Naming Conventions
- Services: `svc-{name}`
- Libraries: `lib-{name}`
- Tools: `tool-{name}`
```

### 4. **Automate Configuration**
```bash
#!/bin/bash
# setup-go-starter.sh

# Install go-starter
go install github.com/francknouama/go-starter@latest

# Import team configuration
go-starter config import https://github.com/yourorg/standards/go-starter.yaml

# Set user-specific values
read -p "Enter your name: " name
read -p "Enter your email: " email

go-starter config set author "$name"
go-starter config set email "$email"

echo "go-starter configured successfully!"
```

## Troubleshooting Configuration

### Configuration Not Loading

```bash
# Check configuration location
go-starter config path

# Validate configuration
go-starter config validate

# Reset to defaults
go-starter config reset
```

### Profile Issues

```bash
# Check current profile
go-starter config profile current

# List all profiles
go-starter config profile list

# Fix profile issues
go-starter config profile repair
```

### Environment Variable Conflicts

```bash
# Show all environment variables
go-starter config env

# Debug configuration loading
GO_STARTER_DEBUG=true go-starter config show
```

## Next Steps

- Learn about [Blueprint Selection](blueprint-selection.md) to choose the right project type
- Read the [Troubleshooting Guide](troubleshooting.md) for common issues
- Check the [FAQ](faq.md) for quick answers

---
*Configuration guide for go-starter v2.0*