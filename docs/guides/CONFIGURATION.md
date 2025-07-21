# Configuration Guide

## Global Configuration

Create a global configuration file to set your preferences:

```yaml
# ~/.go-starter.yaml
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

Override configuration with environment variables:

```bash
export GO_STARTER_AUTHOR="Jane Doe"
export GO_STARTER_EMAIL="jane@example.com"
export GO_STARTER_LICENSE="Apache-2.0"
export GO_STARTER_GO_VERSION="1.23"
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

## Next Steps

- 🚀 **[Quick Start Guide](GETTING_STARTED.md)** - Create your first project
- 📊 **[Logger Guide](LOGGER_GUIDE.md)** - Choose the right logging strategy
- 📖 **[Project Types](BLUEPRINTS.md)** - Explore available templates