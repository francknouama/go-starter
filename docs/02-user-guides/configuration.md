# Configuration Guide

Comprehensive guide to configuring go-starter for your workflow, team standards, and project requirements.

## Table of Contents

- [Global Configuration](#global-configuration)
- [Profile Management](#profile-management)
- [Environment Variables](#environment-variables)
- [Project-Specific Configuration](#project-specific-configuration)
- [Advanced Configuration](#advanced-configuration)
- [Configuration Commands](#configuration-commands)
- [Team Setup](#team-setup)

## Global Configuration

### Basic Configuration File

Create a global configuration file to set your preferences:

```yaml
# ~/.go-starter.yaml
profiles:
  default:
    author: "Your Name"
    email: "your.email@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
      complexity: "standard"
current_profile: "default"
```

### Configuration Location

The configuration file is automatically created in:
- **macOS/Linux**: `~/.go-starter.yaml`
- **Windows**: `%USERPROFILE%\.go-starter.yaml`

You can also specify a custom location:
```bash
export GO_STARTER_CONFIG="/path/to/custom/config.yaml"
```

## Profile Management

### Multiple Profiles

Create different profiles for different contexts:

```yaml
# ~/.go-starter.yaml
profiles:
  personal:
    author: "Your Name"
    email: "personal@example.com"
    license: "MIT"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
      complexity: "simple"
      
  work:
    author: "Your Name"
    email: "work@company.com"
    license: "Apache-2.0"
    defaults:
      goVersion: "1.22"
      framework: "gin"
      logger: "zap"
      complexity: "standard"
      architecture: "clean"
      database: "postgres"
      auth: "jwt"
      
  enterprise:
    author: "Your Name"
    email: "enterprise@bigcorp.com"
    license: "Proprietary"
    defaults:
      goVersion: "1.22"
      framework: "gin"
      logger: "zap"
      complexity: "advanced"
      architecture: "hexagonal"
      database: "postgres,redis"
      auth: "oauth2"
      
current_profile: "personal"
```

### Using Profiles

```bash
# Use specific profile for one project
go-starter new my-project --profile work

# Switch default profile
go-starter config set-profile enterprise

# View current profile
go-starter config current-profile

# List all profiles
go-starter config list-profiles
```

### Profile-Specific Defaults

When you generate a project, go-starter uses profile defaults as starting values, which you can still override:

```bash
# Uses 'work' profile defaults, but overrides logger
go-starter new my-api --profile work --logger logrus
```

## Environment Variables

Override any configuration with environment variables:

### Basic Variables
```bash
export GO_STARTER_AUTHOR="Jane Doe"
export GO_STARTER_EMAIL="jane@example.com"
export GO_STARTER_LICENSE="Apache-2.0"
export GO_STARTER_GO_VERSION="1.22"
```

### Framework and Logger Defaults
```bash
export GO_STARTER_FRAMEWORK="echo"
export GO_STARTER_LOGGER="zap"
export GO_STARTER_COMPLEXITY="standard"
```

### Database and Authentication
```bash
export GO_STARTER_DATABASE="postgres"
export GO_STARTER_AUTH="jwt"
export GO_STARTER_ARCHITECTURE="clean"
```

### Output and Behavior
```bash
export GO_STARTER_OUTPUT_DIR="./projects"
export GO_STARTER_QUIET="true"
export GO_STARTER_ADVANCED="true"
```

### Full Example
```bash
# Set up environment for enterprise development
export GO_STARTER_AUTHOR="Enterprise Team"
export GO_STARTER_EMAIL="team@enterprise.com"
export GO_STARTER_LICENSE="Proprietary"
export GO_STARTER_GO_VERSION="1.22"
export GO_STARTER_FRAMEWORK="gin"
export GO_STARTER_LOGGER="zap"
export GO_STARTER_DATABASE="postgres"
export GO_STARTER_AUTH="oauth2"
export GO_STARTER_ARCHITECTURE="clean"
export GO_STARTER_ADVANCED="true"

# Now all projects will use these defaults
go-starter new enterprise-api --type=web-api
```

## Project-Specific Configuration

### Override for Single Project

```bash
# Override specific settings for this project only
go-starter new special-project \
  --author "Special Author" \
  --email "special@example.com" \
  --license "GPL-3.0" \
  --go-version "1.21" \
  --logger "logrus" \
  --framework "echo"
```

### Consistent Team Configuration

For consistent team setup, you can specify all options explicitly:

```bash
# Enterprise API with specific configuration
go-starter new enterprise-api \
  --type=web-api \
  --architecture=clean \
  --framework=gin \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap \
  --go-version="1.22" \
  --author="Enterprise Team" \
  --license="Proprietary" \
  --advanced
```

## Advanced Configuration

### Template Customization

```yaml
# ~/.go-starter.yaml
templates:
  customPath: "/path/to/custom/templates"
  overrides:
    web-api: "custom-web-api"
    cli: "custom-cli"
  
  # Custom blueprint locations
  blueprints:
    - path: "/company/blueprints"
      priority: 1
    - path: "/team/blueprints" 
      priority: 2
```

### Advanced Project Defaults

```yaml
# ~/.go-starter.yaml
profiles:
  enterprise:
    defaults:
      # Core settings
      goVersion: "1.22"
      framework: "gin"
      logger: "zap"
      architecture: "clean"
      
      # Database configuration
      database:
        driver: "postgres"
        orm: "gorm"
        migrations: true
        seeders: true
        
      # Authentication
      authentication:
        type: "jwt"
        providers: ["google", "github"]
        middleware: true
        
      # Features
      features:
        docker: true
        kubernetes: true
        cicd: true
        monitoring: true
        testing: true
        
      # Observability
      observability:
        metrics: "prometheus"
        tracing: "jaeger"
        logging: "structured"
        
      # Deployment
      deployment:
        platforms: ["docker", "kubernetes"]
        environments: ["dev", "staging", "prod"]
```

### CI/CD Integration

```yaml
# ~/.go-starter.yaml
cicd:
  github:
    enabled: true
    workflows:
      - "ci"
      - "release"
      - "security"
    secrets:
      - "DATABASE_URL"
      - "JWT_SECRET"
      
  gitlab:
    enabled: false
    
  jenkins:
    enabled: false
```

## Configuration Commands

### View Configuration

```bash
# Show current configuration
go-starter config show

# Show specific profile
go-starter config show --profile work

# Show configuration in different formats
go-starter config show --format json
go-starter config show --format yaml
```

### Modify Configuration

```bash
# Set global defaults
go-starter config set author "Jane Doe"
go-starter config set email "jane@example.com"
go-starter config set license "Apache-2.0"
go-starter config set defaults.goVersion "1.22"
go-starter config set defaults.logger "zap"

# Set profile-specific values
go-starter config set --profile work author "Work Jane"
go-starter config set --profile work email "jane@company.com"

# Set complex values
go-starter config set defaults.database.driver "postgres"
go-starter config set defaults.authentication.type "jwt"
```

### Profile Management Commands

```bash
# Create new profile
go-starter config create-profile startup \
  --author "Startup Team" \
  --email "team@startup.com" \
  --license "MIT"

# Copy profile
go-starter config copy-profile work enterprise

# Delete profile
go-starter config delete-profile old-profile

# Set current profile
go-starter config set-profile enterprise
```

### Import/Export Configuration

```bash
# Export current configuration
go-starter config export > my-config.yaml

# Export specific profile
go-starter config export --profile work > work-config.yaml

# Import configuration
go-starter config import my-config.yaml

# Import and merge with existing
go-starter config import --merge company-standards.yaml
```

### Reset Configuration

```bash
# Reset all configuration to defaults
go-starter config reset

# Reset specific profile
go-starter config reset --profile work

# Reset specific values
go-starter config reset defaults.logger
go-starter config reset author email
```

## Team Setup

### Company-Wide Standards

Create a shared configuration file for your team:

```yaml
# company-standards.yaml
profiles:
  company:
    author: "{{DEVELOPER_NAME}}"  # Template to be filled
    email: "{{DEVELOPER_EMAIL}}"
    license: "Proprietary"
    defaults:
      goVersion: "1.22"
      framework: "gin"
      logger: "zap"
      architecture: "clean"
      database: "postgres"
      auth: "oauth2"
      features:
        docker: true
        kubernetes: true
        monitoring: true
        testing: true
```

### Team Setup Script

```bash
#!/bin/bash
# setup-go-starter.sh

# Download company standards
curl -o /tmp/company-standards.yaml https://company.com/go-starter-config.yaml

# Import configuration
go-starter config import /tmp/company-standards.yaml

# Set company profile as default
go-starter config set-profile company

# Set developer-specific information
echo "Setting up go-starter for $(whoami)"
read -p "Enter your full name: " developer_name
read -p "Enter your company email: " developer_email

go-starter config set --profile company author "$developer_name"
go-starter config set --profile company email "$developer_email"

echo "✅ go-starter configured for $developer_name"
```

### Project Templates for Teams

```yaml
# team-templates.yaml
project_templates:
  microservice:
    type: "web-api"
    architecture: "clean"
    framework: "gin"
    database: "postgres,redis"
    auth: "jwt"
    logger: "zap"
    features:
      docker: true
      kubernetes: true
      monitoring: true
      
  cli_tool:
    type: "cli"
    complexity: "standard"
    logger: "logrus"
    features:
      docker: true
      
  library:
    type: "library"
    logger: "slog"
    features:
      ci: true
      docs: true
```

Usage:
```bash
# Use team template
go-starter new user-service --template microservice

# Override template values
go-starter new special-service --template microservice --logger slog
```

### Docker Configuration

For consistent development environments:

```dockerfile
# Dockerfile.go-starter
FROM golang:1.22-alpine

RUN go install github.com/francknouama/go-starter@latest

# Copy team configuration
COPY company-standards.yaml /tmp/
RUN go-starter config import /tmp/company-standards.yaml

# Set up workspace
WORKDIR /workspace
VOLUME ["/workspace"]

ENTRYPOINT ["go-starter"]
```

Usage:
```bash
# Build team image
docker build -t company/go-starter .

# Use in projects
docker run -v $(pwd):/workspace company/go-starter new my-project --type=web-api
```

## Validation and Testing

### Validate Configuration

```bash
# Check configuration syntax
go-starter config validate

# Check specific profile
go-starter config validate --profile work

# Dry run with configuration
go-starter new test-project --type=web-api --dry-run
```

### Test Configuration

```bash
# Test project generation with current config
go-starter new test-$(date +%s) --type=cli --complexity=simple --dry-run

# Test all profiles
for profile in $(go-starter config list-profiles); do
  echo "Testing profile: $profile"
  go-starter new test-$profile --profile $profile --type=web-api --dry-run
done
```

## Troubleshooting Configuration

### Common Issues

#### Configuration Not Loading
```bash
# Check configuration file location
go-starter config show --debug

# Verify file permissions
ls -la ~/.go-starter.yaml

# Check environment variables
env | grep GO_STARTER
```

#### Invalid Configuration
```bash
# Validate syntax
go-starter config validate

# Check for common issues
go-starter config doctor

# Reset if corrupted
go-starter config reset
mv ~/.go-starter.yaml ~/.go-starter.yaml.backup
```

#### Profile Issues
```bash
# List available profiles
go-starter config list-profiles

# Check current profile
go-starter config current-profile

# Switch to default profile
go-starter config set-profile default
```

### Debug Mode

Enable debug mode for configuration troubleshooting:

```bash
# Enable debug output
export GO_STARTER_DEBUG=true

# Or use debug flag
go-starter --debug config show
go-starter --debug new my-project --type=web-api
```

## Best Practices

### 1. Use Profiles for Different Contexts
- **Personal**: Simple projects, learning
- **Work**: Company standards, team consistency
- **Open Source**: Public repositories, MIT license

### 2. Environment Variables for CI/CD
Use environment variables in CI/CD pipelines for consistency:

```yaml
# .github/workflows/generate.yml
env:
  GO_STARTER_AUTHOR: "CI Bot"
  GO_STARTER_EMAIL: "ci@company.com"
  GO_STARTER_LICENSE: "Proprietary"
  GO_STARTER_LOGGER: "zap"
```

### 3. Version Control Configuration
Include team configuration in your repository:

```bash
# .go-starter.yaml (project-specific)
profiles:
  project:
    defaults:
      logger: "zap"
      database: "postgres"
      auth: "jwt"
```

### 4. Regular Configuration Audits
```bash
# Monthly configuration review
go-starter config show > config-$(date +%Y-%m).yaml
```

---

**Next Steps**: 
- Set up your [first profile](quick-start.md) and generate a project
- Learn about [blueprint selection](blueprint-selection.md) for different project types
- Check [troubleshooting](troubleshooting.md) if you encounter issues