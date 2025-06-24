# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a comprehensive Go project generator called "go-starter" that generates Go project structures with modern best practices, multiple architecture patterns, and deployment configurations. It combines the simplicity of create-react-app with the flexibility of Spring Initializr, offering both CLI and web interfaces with progressive disclosure for beginners and advanced developers.

## Development Commands

### Building and Running
- `go build -o bin/go-starter main.go` - Build the CLI tool
- `go install .` - Install the CLI tool globally
- `make build` - Build using Makefile
- `make install` - Install the CLI tool
- `make dev-build` - Development build with race detection
- `make run` - Start development server (for web UI phase)

### Testing
- `go test -v ./...` - Run all tests
- `make test` - Run tests via Makefile
- **Critical**: Integration tests must validate that all generated projects compile successfully
- Test template generation with various configurations
- Test conditional file generation logic

### Code Quality
- `golangci-lint run` - Run linting (essential before commits)
- `make lint` - Run linting via Makefile
- `go generate ./...` - Generate embedded templates
- `make generate` - Generate embedded templates via Makefile

### Web Development (Phase 3)
- `npm run dev` - Start React development server (Vite)
- `npm run build` - Build web UI for production
- `vite preview` - Preview production build

## Architecture Overview

### Four-Phase Implementation Strategy

#### Phase 1: Core CLI Tool
- **Goal**: Functional CLI with 4 basic templates (Web API, CLI, Library, Lambda)
- **Key Components**: Cobra framework, interactive prompts, template engine, basic generation
- **Templates**: Standard Web API, CLI Application, Library, AWS Lambda
- **Architecture Patterns**: Standard only

#### Phase 2: Complete Template System  
- **Goal**: All 12 project types with multiple architecture patterns
- **Templates**: web-api (4 architectures), cli, library, lambda, lambda-proxy, event-driven, microservice, monolith, workspace
- **Architecture Patterns**: Standard, Clean Architecture, DDD, Hexagonal, Event-driven
- **Features**: Conditional generation, template validation, enhanced prompts

#### Phase 3: Web UI
- **Goal**: React-based web interface with live preview
- **Tech Stack**: React + Vite, WebSocket for real-time updates
- **Features**: Progressive disclosure, file structure visualization, live preview, project download
- **API**: RESTful endpoints + WebSocket for real-time features

#### Phase 4: Advanced Features
- **Goal**: GitHub integration, template marketplace, deployment platforms
- **Features**: OAuth authentication, community templates, Vercel/Railway/Netlify deployment
- **Enterprise**: Template marketplace, analytics, monitoring

### Core Components Architecture

#### CLI Framework (`cmd/`)
- **Framework**: Cobra for command structure
- **Commands**: `new`, `list`, `version`, `config`, `completion`  
- **Interactive Mode**: AlecAivazis/survey for user prompts
- **Progressive Disclosure**: Basic mode (essential options) vs Advanced mode (all features)

#### Template Engine (`internal/templates/`)
- **Storage**: Embedded templates using `embed.FS`
- **Processing**: Go `text/template` with Sprig functions
- **Registry**: Template loading and management system
- **Validation**: Template syntax and variable validation

#### Generator Service (`internal/generator/`)
- **Core Logic**: Project generation with conditional file creation
- **Recovery**: Rollback mechanism for failed generations
- **Hooks**: Post-generation scripts and commands
- **Memory Mode**: In-memory generation for web interface

#### Web Interface (`web/`)
- **Frontend**: React with TypeScript, Tailwind CSS
- **Build Tool**: Vite for fast development and building
- **Real-time**: WebSocket for live preview updates
- **API**: RESTful backend with Gin framework

## Template System Deep Dive

### Template Structure
Each template consists of:
1. **template.yaml** - Metadata, variables, file definitions, dependencies
2. **Template files** - `.tmpl` files with Go template syntax
3. **Conditional logic** - Files generated based on configuration

### 12 Core Templates

| Template | Use Case | Key Features |
|----------|----------|-------------|
| **Standard Web API** | Basic REST APIs | Simple structure, fast setup |
| **Clean Architecture** | Enterprise apps | Layered architecture, separation of concerns |
| **DDD Web API** | Complex domains | Domain-focused, business logic emphasis |
| **Hexagonal Architecture** | Testable apps | Ports & adapters, dependency inversion |
| **CLI Application** | Command tools | Cobra framework, subcommands |
| **Library** | Reusable packages | Public API focus, examples |
| **AWS Lambda** | Serverless functions | Event processing, cloud-native |
| **Lambda API Proxy** | API Gateway | Request routing, proxy patterns |
| **Event-Driven** | CQRS/Event Sourcing | Scalable, distributed systems |
| **Microservice** | Distributed systems | gRPC, containerized, service mesh ready |
| **Monolith** | Traditional web apps | Modular structure, all-in-one deployment |
| **Go Workspace** | Multi-module projects | Monorepo, shared dependencies |

### Template Variables
Standard variables available in all templates:
- `{{.ProjectName}}` - Project name
- `{{.ModulePath}}` - Go module path (e.g., github.com/user/project)
- `{{.GoVersion}}` - Go version (default: 1.21)
- `{{.Framework}}` - Selected framework (gin, echo, cobra, etc.)
- `{{.Architecture}}` - Architecture pattern (standard, clean, ddd, hexagonal)
- `{{.LoggerType}}` - Logging library (slog, zap, logrus, zerolog)
- `{{.Features}}` - Feature configuration object

### Logger Selector System

The go-starter generator includes a sophisticated logger selector that allows developers to choose from multiple high-quality logging libraries while maintaining a consistent interface across all generated code.

#### Supported Loggers

| Logger | Package | Performance | Use Case | Default |
|--------|---------|-------------|----------|---------|
| **slog** | `log/slog` | Good | Standard library, structured logging | ✅ |
| **zap** | `go.uber.org/zap` | Excellent | High performance, zero allocation | |
| **logrus** | `github.com/sirupsen/logrus` | Good | Feature-rich, popular choice | |
| **zerolog** | `github.com/rs/zerolog` | Excellent | Zero allocation, chainable API | |

#### Logger Selection Benefits

- **Consistent Interface**: All loggers implement the same interface for seamless switching
- **Conditional Dependencies**: Only the selected logger's dependencies are included
- **Performance Optimization**: Choose the logger that best fits your performance requirements
- **Configuration Driven**: Logger behavior controlled through configuration files
- **Template Integration**: All generated templates use the logger interface consistently

#### Usage Examples

**slog (Default)**:
```go
logger.Info("Server starting", "port", 8080, "env", "production")
```

**zap**:
```go  
logger.Info("Server starting", zap.Int("port", 8080), zap.String("env", "production"))
```

**logrus**:
```go
logger.WithFields(logrus.Fields{"port": 8080, "env": "production"}).Info("Server starting")
```

**zerolog**:
```go
logger.Info().Int("port", 8080).Str("env", "production").Msg("Server starting")
```

### Conditional Generation
Templates use Go template conditions for optional files:
```yaml
files:
  - source: database.go.tmpl
    destination: internal/database/database.go
    condition: "{{ne .Features.Database.Driver \"\"}}"
  - source: auth.go.tmpl
    destination: internal/middleware/auth.go
    condition: "{{eq .Features.Authentication.Type \"jwt\"}}"
  - source: logger/zap.go.tmpl
    destination: internal/logger/zap.go
    condition: "{{eq .LoggerType \"zap\"}}"
```

## Configuration Management

### CLI Configuration
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
current_profile: "default"
```

### Project Configuration
```yaml
# project.yaml
name: my-awesome-api
module: github.com/myuser/my-awesome-api
type: api                    # api, cli, library, lambda, lambda-proxy, microservice, monolith, workspace
goVersion: "1.21"
architecture: hexagonal     # standard, clean, ddd, hexagonal, event-driven
framework: gin              # gin, echo, fiber, chi, cobra
logger: slog                # slog, zap, logrus, zerolog

features:
  database:
    driver: postgres         # postgres, mysql, mongodb, sqlite, redis
    orm: gorm               # gorm, sqlx, sqlc, ent
  authentication:
    type: jwt               # jwt, oauth2, session, api-key
    providers: ["google", "github"]
  deployment:
    targets: ["docker", "kubernetes", "lambda"]
  testing:
    framework: testify      # testify, ginkgo
    coverage: true
  logging:
    level: info             # debug, info, warn, error
    format: json            # json, text, console
```

## Development Workflow

### Adding New Templates
1. Create directory in `templates/` (e.g., `templates/new-type/`)
2. Add `template.yaml` with metadata and file definitions
3. Create template files with `.tmpl` extension using Go template syntax
4. Update template registry in `internal/templates/registry.go`
5. Add prompts for new template in `internal/prompts/interactive.go`
6. Add tests to validate template generation

### Template Development Best Practices
- Use descriptive variable names in templates
- Include conditional logic for optional features
- Provide sensible defaults for all variables
- Test templates with various configuration combinations
- Include proper error handling in generated code
- Follow Go best practices in generated project structure

## Testing Strategy

### Critical Testing Requirements
- **Template Validation**: All templates must parse without errors
- **Generation Testing**: Generated projects must compile successfully with `go build`
- **Logger Testing**: All logger types must generate working implementations and compile
- **Integration Testing**: End-to-end CLI workflow testing
- **Web UI Testing**: API endpoints and WebSocket functionality
- **Cross-platform Testing**: Windows, macOS, Linux compatibility

### Test Categories
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: Full project generation workflow
3. **Template Tests**: Validate all templates generate working code
4. **Logger Tests**: Test each logger implementation with various configurations
5. **API Tests**: Web interface endpoint testing
6. **CLI Tests**: Command-line interface testing

## Database Schema (Phases 3-4)

### Core Tables
- `users` - Authentication and user profiles
- `projects` - Generated project metadata and sharing
- `marketplace_templates` - Community-contributed templates
- `template_ratings` - Template reviews and ratings
- `analytics_events` - Usage analytics and metrics
- `api_keys` - API access management

### Key Relationships
- Users can create multiple projects
- Projects reference templates used for generation
- Templates can be rated and reviewed by users
- Analytics track template usage patterns

## Security Considerations

### Input Validation
- **Sanitize all inputs**: Project names, module paths, template variables
- **Validate module paths**: Ensure proper Go module syntax
- **Template security**: Scan community templates for malicious code
- **Path traversal protection**: Prevent template files from accessing unauthorized paths

### Web Interface Security
- **Authentication**: JWT-based authentication for user sessions
- **CORS**: Proper cross-origin request handling
- **Rate limiting**: Prevent abuse of generation endpoints
- **Input sanitization**: All user inputs must be validated

## Performance Considerations

### Template Engine Optimization
- **Template caching**: Cache parsed templates in memory
- **Parallel generation**: Generate multiple files concurrently
- **Memory management**: Efficient handling of large projects
- **File I/O optimization**: Batch file operations when possible

### Web Interface Performance
- **Lazy loading**: Load templates and previews on demand
- **WebSocket efficiency**: Debounce real-time preview updates
- **Caching**: Cache generated previews and project structures
- **Bundle optimization**: Minimize JavaScript bundle size

## Deployment and Production

### CLI Distribution
- **GitHub Releases**: Binary distribution for multiple platforms
- **Package managers**: Homebrew, Chocolatey, APT/YUM packages
- **Auto-update**: Self-updating mechanism for CLI tool
- **Docker**: Containerized version for CI/CD usage

### Web Interface Deployment
- **Static hosting**: Netlify, Vercel for frontend
- **Backend hosting**: Railway, Render for API server
- **Database**: PostgreSQL with connection pooling
- **Monitoring**: Health checks, metrics, error tracking

## Common Issues and Solutions

### Template Generation Failures
- **Template syntax errors**: Use `go generate ./...` to validate
- **Missing variables**: Ensure all template variables are defined
- **File conflicts**: Check for overlapping file destinations
- **Rollback mechanism**: Use recovery system for partial failures

### Performance Issues
- **Large projects**: Use streaming for file downloads
- **Memory usage**: Implement garbage collection for long-running sessions
- **Database queries**: Use connection pooling and query optimization
- **Template loading**: Cache templates to avoid repeated parsing

### Cross-platform Compatibility
- **File paths**: Use `filepath.Join()` for cross-platform paths
- **Line endings**: Handle CRLF vs LF appropriately
- **Permissions**: Set proper file permissions on Unix systems
- **Shell commands**: Provide alternatives for different platforms

## Future Extensibility

### Plugin System (Phase 4+)
- **HashiCorp go-plugin**: Process isolation for custom generators
- **Template marketplace**: Community template sharing and discovery
- **Organization templates**: Private template repositories for enterprises
- **Custom hooks**: Pre/post generation custom scripts

### Integration Opportunities
- **IDE Extensions**: VS Code, GoLand plugins
- **CI/CD Integration**: GitHub Actions, GitLab CI templates
- **Cloud Platforms**: AWS, GCP, Azure integration
- **Monitoring Tools**: Observability and metrics integration

This comprehensive project represents a significant advancement in Go developer tooling, filling critical gaps in the ecosystem while providing a foundation for future innovation.