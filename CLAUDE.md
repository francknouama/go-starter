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

### Progressive Disclosure System ✨

**NEW**: go-starter now features a comprehensive progressive disclosure system that adapts the interface based on user experience level and project complexity needs.

#### Basic Usage (Beginner-Friendly)
- `go-starter new` - Interactive mode with basic options only
- `go-starter new my-app --type=cli` - Simple direct generation
- `go-starter new --help` - Shows essential flags only (14 flags)

#### Advanced Usage (Expert Mode)
- `go-starter new --advanced` - Interactive mode with all options
- `go-starter new my-app --type=web-api --advanced` - Advanced direct generation
- `go-starter new --advanced --help` - Shows all flags (18+ flags)

#### Complexity-Aware Generation
- `go-starter new --type=cli --complexity=simple` - Generate simple CLI (8 files)
- `go-starter new --type=cli --complexity=standard` - Generate standard CLI (29 files)
- `go-starter new --type=web-api --architecture=clean` - Generate Clean Architecture API
- `go-starter new --complexity=advanced` - Advanced project structure

#### Progressive Disclosure Flags
- `--basic` - Show only essential options (default for new users)
- `--advanced` - Enable advanced configuration with all options
- `--complexity [simple|standard|advanced|expert]` - Specify project complexity level
- `--dry-run` - Preview project structure without creating files
- `--help` - Context-aware help (basic vs advanced mode)

### Testing
- `go test -v ./...` - Run all tests (unit, integration, acceptance)
- `go test ./internal/prompts/... -v` - Run progressive disclosure tests
- `go test ./tests/acceptance/cli/progressive_disclosure_test.go -v` - ATDD tests
- `make test` - Run tests via Makefile

#### Progressive Disclosure Specific Testing
- **Unit Tests**: `go test ./internal/prompts/progressive_test.go -v`
- **Blueprint Selection**: Test complexity-aware blueprint mapping
- **Help Filtering**: Validate basic vs advanced flag visibility  
- **Interactive Prevention**: Ensure no prompts when flags provided
- **Default Application**: Test smart defaults for CLI blueprints

#### Integration Testing Requirements
- **Critical**: All complexity levels must generate working code
- **Compilation**: Generated projects must compile with `go build`
- **File Counts**: Simple CLI (8 files), Standard CLI (29 files)
- **Blueprint Validation**: All blueprints parse without errors
- **Cross-platform**: Test on Windows, macOS, Linux

#### Test Categories and Coverage
- Test blueprint generation with various configurations  
- Test conditional file generation logic
- Test help filtering and flag visibility
- Test CLI defaults and non-interactive mode
- Test complexity flag validation and error handling

### Code Quality
- `golangci-lint run` - Run linting (essential before commits)
- `make lint` - Run linting via Makefile
- `go generate ./...` - Generate embedded blueprints
- `make generate` - Generate embedded blueprints via Makefile

### Web Development (Phase 3)
- `npm run dev` - Start React development server (Vite)
- `npm run build` - Build web UI for production
- `vite preview` - Preview production build

## Architecture Overview

### Four-Phase Implementation Strategy

#### Phase 1: Core CLI Tool
- **Goal**: Functional CLI with 4 basic blueprints (Web API, CLI, Library, Lambda)
- **Key Components**: Cobra framework, interactive prompts, blueprint engine, basic generation
- **Blueprints**: Standard Web API, CLI Application, Library, AWS Lambda
- **Architecture Patterns**: Standard only

#### Phase 2: Complete Blueprint System  
- **Goal**: All 12 project types with multiple architecture patterns
- **Blueprints**: web-api (4 architectures), cli, library, lambda, lambda-proxy, event-driven, microservice, monolith, workspace
- **Architecture Patterns**: Standard, Clean Architecture, DDD, Hexagonal, Event-driven
- **Features**: Conditional generation, blueprint validation, enhanced prompts

#### Phase 3: Web UI
- **Goal**: React-based web interface with live preview
- **Tech Stack**: React + Vite, WebSocket for real-time updates
- **Features**: Progressive disclosure, file structure visualization, live preview, project download
- **API**: RESTful endpoints + WebSocket for real-time features

#### Phase 4: Advanced Features
- **Goal**: GitHub integration, blueprint marketplace, deployment platforms
- **Features**: OAuth authentication, community blueprints, Vercel/Railway/Netlify deployment
- **Enterprise**: Template marketplace, analytics, monitoring

### Core Components Architecture

#### CLI Framework (`cmd/`)
- **Framework**: Cobra for command structure
- **Commands**: `new`, `list`, `version`, `config`, `completion`  
- **Interactive Mode**: AlecAivazis/survey for user prompts
- **Progressive Disclosure**: Basic mode (essential options) vs Advanced mode (all features)

#### Template Engine (`internal/templates/`)
- **Storage**: Embedded blueprints using `embed.FS`
- **Processing**: Go `text/template` with Sprig functions
- **Registry**: Blueprint loading and management system
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

## Progressive Disclosure System Deep Dive ✨

> **NEW**: Full progressive disclosure implementation with basic/advanced modes, complexity-aware generation, and intelligent help filtering.

### Progressive Disclosure Features

#### 🎯 Smart Help System
- **Basic Mode**: Shows 14 essential flags for beginners
- **Advanced Mode**: Shows 18+ flags for power users  
- **Context-Aware**: Automatically detects user intent from flags
- **No Duplicates**: Intelligent flag deduplication between local and persistent flags
- **Visual Design**: Styled output with hints and tips

#### 🚀 Complexity-Aware Blueprint Selection
- **Automatic Selection**: `--complexity=simple` → `cli-simple` blueprint
- **Smart Defaults**: CLI blueprints auto-set cobra + slog when complexity is specified
- **No Prompting**: Sufficient flags prevent interactive prompts
- **Default Module**: Auto-generates module path for testing scenarios

#### 🛠️ Two-Tier CLI System
- **Simple CLI**: 8 files, minimal structure, perfect for learning
- **Standard CLI**: 29 files, production-ready, full feature set
- **Clear Distinction**: Different use cases and complexity levels
- **Migration Path**: Easy upgrade from simple to standard

#### 📚 Progressive Learning Philosophy
- **Start Simple**: Minimal viable structure for beginners
- **Grow Organically**: Add complexity only when needed
- **Clear Path**: Obvious progression from simple to advanced
- **No Overwhelm**: Hide complexity until users are ready

## Blueprint System Deep Dive

> **Updated v2.1**: Blueprint Selection Guide, Complexity Levels, CLI Two-Tier Approach, Progressive Complexity Philosophy, and Progressive Disclosure System. See sections below for detailed guidance on choosing the right blueprint for your project.

### Blueprint Structure
Each blueprint consists of:
1. **template.yaml** - Metadata, variables, file definitions, dependencies
2. **Template files** - `.tmpl` files with Go template syntax
3. **Conditional logic** - Files generated based on configuration

### 12 Core Blueprints

| Blueprint | Use Case | Key Features |
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

## Blueprint Selection Guide

Choosing the right blueprint is crucial for project success. This guide helps you make informed decisions based on your project requirements and team expertise.

### Blueprint Complexity Levels

We categorize blueprints into four complexity levels to help you choose appropriately:

| Complexity | Description | Typical Characteristics |
|------------|-------------|------------------------|
| **Beginner** | Simple, straightforward projects | < 10 files, minimal dependencies, standard patterns |
| **Intermediate** | Moderate complexity with some patterns | 10-25 files, common patterns, standard dependencies |
| **Advanced** | Complex architecture with multiple patterns | 25-50 files, advanced patterns, multiple dependencies |
| **Expert** | Highly complex, enterprise-grade | 50+ files, complex patterns, extensive dependencies |

### Blueprint Selection Matrix

| Blueprint | Complexity | Use Case | When to Use | Architecture Pattern |
|-----------|------------|----------|-------------|---------------------|
| **Simple CLI** | Beginner | Basic command-line tools | Quick scripts, simple utilities | None (procedural) |
| **Standard CLI** | Intermediate | Full-featured CLI apps | Complex CLI tools, multiple commands | MVC-lite |
| **Library** | Beginner | Reusable packages | Shared code, utilities | Public API pattern |
| **Standard Web API** | Intermediate | Basic REST APIs | Simple CRUD, microservices | Standard layered |
| **Clean Architecture** | Advanced | Enterprise APIs | Complex business logic, testability | Clean Architecture |
| **DDD Web API** | Advanced | Domain-complex APIs | Rich domain models, business rules | Domain-Driven Design |
| **Hexagonal Architecture** | Expert | Highly testable APIs | Multiple adapters, high testability | Ports & Adapters |
| **AWS Lambda** | Beginner | Simple serverless | Event handlers, webhooks | Functional |
| **Lambda API Proxy** | Intermediate | API Gateway integration | RESTful serverless APIs | Proxy pattern |
| **Event-Driven** | Expert | CQRS/Event Sourcing | Event-based systems, audit trails | Event Sourcing |
| **Microservice** | Advanced | Distributed services | Service mesh, gRPC communication | Microservice patterns |
| **Monolith** | Intermediate | Traditional web apps | All-in-one deployment | Modular monolith |
| **Go Workspace** | Advanced | Multi-module projects | Monorepos, shared libraries | Workspace pattern |

## CLI Blueprint Audit Findings

Based on comprehensive analysis (Issue [#149](https://github.com/francknouama/go-starter/issues/149)), we've identified key improvements for CLI blueprints:

### Two-Tier Approach

The audit revealed a complexity mismatch between user expectations and the standard CLI blueprint. We've implemented a two-tier approach:

1. **Simple CLI** (NEW)
   - **Complexity**: Beginner
   - **Files**: Reduced from 25 to 8 files
   - **Purpose**: Quick command-line utilities
   - **Features**: Basic flags, simple output, minimal structure

2. **Standard CLI**
   - **Complexity**: Intermediate
   - **Files**: 25 files with full structure
   - **Purpose**: Production-ready CLI applications
   - **Features**: Subcommands, configuration, advanced patterns

### When to use Simple CLI

Choose the Simple CLI blueprint when:
- Building quick utilities or scripts
- Learning Go or exploring ideas
- Creating internal tools with minimal requirements
- Prototyping command-line interfaces
- You need < 3 commands and < 5 flags

### When to use Standard CLI

Choose the Standard CLI blueprint when:
- Building production CLI tools
- Requiring multiple subcommands
- Needing configuration file support
- Implementing complex business logic
- Building tools for distribution

## Audit Findings Integration

The blueprint audit (Issue [#149](https://github.com/francknouama/go-starter/issues/149)) revealed several critical insights:

### Key Findings

1. **Complexity mismatch**: The standard CLI blueprint was overly complex for simple use cases
2. **File count concerns**: 25 files for a basic CLI tool discouraged beginners
3. **Progressive learning**: Users need a stepping stone between "hello world" and production CLIs

### Improvements Made

1. **Two-tier approach**: Introduced Simple CLI blueprint (Issue [#56](https://github.com/francknouama/go-starter/issues/56))
2. **Progressive disclosure**: ✅ **COMPLETED** - Implemented full progressive disclosure system (Issue [#150](https://github.com/francknouama/go-starter/issues/150))
3. **Complexity reduction**: Simple CLI reduced from 29 to 8 files (73% reduction)
4. **Clear guidance**: Added selection matrix and use case documentation
5. **Smart help filtering**: Basic vs advanced help modes with context awareness
6. **Interactive prompting fixes**: Eliminate prompts when sufficient flags provided
7. **Default handling**: Smart defaults for CLI blueprints to improve UX
8. **Comprehensive testing**: Unit, integration, and ATDD test coverage

## Progressive Disclosure System Implementation 🔧

### Architecture Overview

The progressive disclosure system is implemented across multiple layers:

#### Core Components

1. **Complexity Level System** (`internal/prompts/progressive.go`)
   ```go
   type ComplexityLevel int
   const (
       ComplexitySimple   // Beginner-friendly, minimal structure
       ComplexityStandard // Balanced, production-ready
       ComplexityAdvanced // Enterprise patterns
       ComplexityExpert   // Full-featured, all options
   )
   ```

2. **Disclosure Mode System** (`internal/prompts/interfaces/types.go`)
   ```go
   type DisclosureMode int
   const (
       DisclosureModeBasic    // Essential options only
       DisclosureModeAdvanced // All available options
   )
   ```

3. **Custom Help Rendering** (`cmd/new.go`)
   - Smart flag filtering based on disclosure mode
   - Duplicate flag elimination
   - Context-aware help hints
   - Visual styling with lipgloss

#### Progressive Help System Details

**Basic Mode Help (Default)**:
- Shows 14 essential flags: `--type`, `--name`, `--module`, `--framework`, `--logger`, etc.
- Hides advanced flags: `--database-driver`, `--auth-type`, `--banner-style`, etc.
- Includes hint: "💡 Use --advanced to see all available options"

**Advanced Mode Help**:
- Shows all 18+ flags including database, authentication, and deployment options
- Includes hint: "💡 Use --basic to see only essential options"
- No flag filtering - full feature exposure

#### Blueprint Selection Logic

```go
func SelectBlueprintForComplexity(blueprintType string, complexity ComplexityLevel) string {
    if blueprintType == "cli" {
        switch complexity {
        case ComplexitySimple:
            return "cli-simple"  // 8 files, minimal structure
        default:
            return "cli"         // 29 files, full structure
        }
    }
    return blueprintType // Other blueprints unchanged
}
```

#### Smart Defaults System

When complexity and type are specified together, the system automatically sets defaults to prevent unnecessary prompting:

- **CLI Projects**: 
  - Framework: `cobra` (industry standard)
  - Logger: `slog` (Go standard library)
  - Module: `github.com/username/{project-name}` (for testing)

#### Interactive Prompting Prevention

The system uses a two-stage approach:

1. **Pre-Prompt Analysis**: Check if sufficient flags are provided
2. **Default Application**: Apply smart defaults based on blueprint type and complexity
3. **Prompt Bypass**: Skip interactive prompts when configuration is complete

### Implementation Files

#### Core Logic Files
- `internal/prompts/progressive.go` - Core progressive disclosure logic
- `internal/prompts/interfaces/types.go` - Type definitions and interfaces
- `cmd/new.go` - CLI command integration and help rendering

#### Prompter Integration
- `internal/prompts/bubbletea/prompter.go` - BubbleTea prompter with disclosure support
- `internal/prompts/survey/prompter.go` - Survey prompter with disclosure support
- Both implement `GetProjectConfigWithDisclosure` method

#### Test Coverage
- `internal/prompts/progressive_test.go` - Unit tests for core logic
- `tests/acceptance/cli/progressive_disclosure_test.go` - ATDD tests
- `tests/acceptance/blueprints/cli/cli_simple_atdd_test.go` - Simple CLI validation

### Technical Decision Points

#### Help System Architecture
**Decision**: Custom help function vs. Cobra's built-in help
**Rationale**: Cobra's built-in help doesn't support dynamic flag filtering
**Implementation**: 
```go
newCmd.SetHelpFunc(progressiveHelpFunc)
```

#### Blueprint Naming Strategy
**Decision**: `cli-simple` vs `cli` (not `cli-standard`)
**Rationale**: Maintain existing blueprint registry while adding simple variant
**Implementation**: Complexity-aware selection function maps appropriately

#### Default Module Path
**Decision**: Generate `github.com/username/{project}` for tests
**Rationale**: Prevents test hangs while providing valid Go module syntax
**Implementation**: Applied only when `--complexity` and `--type=cli` are specified

#### Flag Categorization
**Essential Flags** (Basic Mode):
```go
var essentialFlags = map[string]bool{
    "name": true, "type": true, "module": true,
    "framework": true, "logger": true, "go-version": true,
    "output": true, "help": true, "quiet": true,
    "dry-run": true, "no-git": true, "random-name": true,
    "advanced": true, "basic": true, "complexity": true,
}
```

**Advanced Flags** (Advanced Mode Only):
- `--database-driver`, `--database-orm`
- `--auth-type`, `--banner-style`
- `--no-banner`, `--architecture`

### Performance Characteristics

- **Help Rendering**: O(n) where n = number of flags (~20ms for 18 flags)
- **Blueprint Selection**: O(1) constant time lookup
- **Default Application**: O(1) constant time assignments
- **Flag Filtering**: O(n) linear scan with early termination

### Error Handling

- **Invalid Complexity**: Clear error message with valid options
- **Missing Dependencies**: Smart defaults prevent most missing value errors
- **Help Rendering Failures**: Graceful fallback to standard Cobra help
- **Blueprint Loading Errors**: Detailed error messages with troubleshooting hints

### Future Extensions

1. **Blueprint-Specific Complexity**: Different complexity models per blueprint type
2. **User Preference Learning**: Remember user's preferred disclosure mode
3. **Dynamic Complexity Assessment**: Suggest complexity based on project requirements
4. **Progressive Feature Unlocking**: Gradually expose features as users gain experience

### Usage Examples and Workflows

#### Beginner Workflow (Basic Mode)
```bash
# Step 1: See only essential help
go-starter new --help
# Shows: 14 essential flags with beginner-friendly descriptions

# Step 2: Generate simple CLI
go-starter new my-tool --type=cli --complexity=simple
# Result: 8 files, minimal structure, no prompts

# Step 3: Preview before creating
go-starter new my-tool --type=cli --complexity=simple --dry-run
# Shows: File list and structure preview
```

#### Expert Workflow (Advanced Mode)
```bash
# Step 1: See all available options
go-starter new --advanced --help
# Shows: 18+ flags including database, auth, deployment options

# Step 2: Generate complex project with all options
go-starter new enterprise-api \
  --type=web-api \
  --architecture=hexagonal \
  --database-driver=postgres \
  --database-orm=gorm \
  --auth-type=jwt \
  --logger=zap \
  --advanced
# Result: Full enterprise structure with all requested features
```

#### Progressive Learning Path
```bash
# Week 1: Start simple
go-starter new hello-cli --complexity=simple

# Week 2: Try standard CLI
go-starter new todo-cli --complexity=standard

# Week 3: Explore web APIs
go-starter new rest-api --type=web-api

# Week 4: Advanced patterns
go-starter new enterprise-api --type=web-api --architecture=clean --advanced
```

#### Development and Testing Workflows
```bash
# Quick prototype testing
go-starter new test-{1..5} --type=cli --complexity=simple --dry-run

# Validate all complexity levels work
for complexity in simple standard advanced expert; do
  go-starter new test-$complexity --type=cli --complexity=$complexity --dry-run
done

# Test different architectures
for arch in standard clean ddd hexagonal; do
  go-starter new api-$arch --type=web-api --architecture=$arch --dry-run
done
```

### Troubleshooting Progressive Disclosure

#### Common Issues and Solutions

**Problem**: Help shows too many options for beginners
```bash
# Wrong
go-starter new --help  # Shows advanced flags

# Right
go-starter new --basic --help  # Shows only essential flags
# Or just rely on default basic mode
go-starter new --help  # Already defaults to basic
```

**Problem**: Interactive prompts appear even with flags
```bash
# Wrong - missing required flags
go-starter new --type=cli --complexity=simple
# Prompts for module path

# Right - include all required flags
go-starter new my-app --type=cli --complexity=simple --module=github.com/user/my-app
```

**Problem**: Can't find advanced options
```bash
# Solution: Use advanced mode
go-starter new --advanced --help
# Shows: database, auth, deployment options
```

**Problem**: Too many files generated for simple project
```bash
# Wrong - generates 29 files
go-starter new my-tool --type=cli

# Right - generates 8 files  
go-starter new my-tool --type=cli --complexity=simple
```

**Problem**: Complexity level not working as expected
```bash
# Check valid complexity levels
go-starter new --complexity=invalid  # Shows error with valid options
go-starter new --complexity=simple   # ✓ Valid
go-starter new --complexity=standard # ✓ Valid
```

#### Testing and Validation

**Verify Progressive Disclosure Works**:
```bash
# Test basic help (should show ~14 flags)
go-starter new --help | grep -c "^    --"

# Test advanced help (should show ~18+ flags)
go-starter new --advanced --help | grep -c "^    --"

# Test complexity generation
go-starter new test-simple --type=cli --complexity=simple --dry-run | grep -c "Files to be generated"
# Should show 8 files

go-starter new test-standard --type=cli --complexity=standard --dry-run | grep -c "Files to be generated"  
# Should show 29 files
```

**Debug Help System**:
```bash
# Check flag filtering
go-starter new --basic --help | grep "database-driver"
# Should be empty (filtered out)

go-starter new --advanced --help | grep "database-driver"
# Should show the flag
```

## Progressive Complexity Philosophy

go-starter embraces a "Start simple, grow as needed" philosophy:

### Core Principles

1. **Start Simple**: Begin with the minimal viable structure
2. **Grow Organically**: Add complexity only when needed
3. **Clear Progression**: Obvious path from simple to advanced
4. **No Premature Optimization**: Avoid over-engineering from the start

### Learning Path

The recommended progression for CLI development:

1. **Simple CLI** → Learn basics, understand flags and output
2. **Standard CLI** → Add subcommands, configuration, and structure
3. **Advanced patterns** → Integrate with web APIs, databases, or complex logic

### Migrating Between Blueprints

When your simple CLI outgrows its initial structure:

1. **Identify pain points**: Multiple commands needed? Configuration required?
2. **Plan the migration**: Map simple structure to standard structure
3. **Incremental refactoring**: Move code piece by piece
4. **Maintain compatibility**: Ensure existing functionality remains intact

Example migration from simple to standard:
- `main.go` → Split into `cmd/root.go` and command files
- Inline logic → Move to `internal/` packages
- Direct flag access → Configuration struct pattern
- Simple output → Structured logging with levels

### Complexity Criteria

When evaluating blueprint complexity, we consider:

1. **Number of files**: More files = higher complexity
2. **Architecture patterns**: Advanced patterns increase complexity
3. **Dependencies**: External dependencies add complexity
4. **Abstraction levels**: More layers = higher complexity
5. **Configuration options**: Flexibility adds complexity

### Blueprint Variables
Standard variables available in all blueprints:
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
- **Blueprint Integration**: All generated blueprints use the logger interface consistently

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
Blueprints use Go template conditions for optional files:
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

### Adding New Blueprints
1. Create directory in `blueprints/` (e.g., `blueprints/new-type/`)
2. Add `template.yaml` with metadata and file definitions
3. Create template files with `.tmpl` extension using Go template syntax
4. Update blueprint registry in `internal/templates/registry.go`
5. Add prompts for new blueprint in `internal/prompts/interactive.go`
6. Add tests to validate blueprint generation

### Blueprint Development Best Practices
- Use descriptive variable names in templates
- Include conditional logic for optional features
- Provide sensible defaults for all variables
- Test blueprints with various configuration combinations
- Include proper error handling in generated code
- Follow Go best practices in generated project structure

## Testing Strategy

### Critical Testing Requirements
- **Blueprint Validation**: All blueprints must parse without errors
- **Generation Testing**: Generated projects must compile successfully with `go build`
- **Logger Testing**: All logger types must generate working implementations and compile
- **Integration Testing**: End-to-end CLI workflow testing
- **Web UI Testing**: API endpoints and WebSocket functionality
- **Cross-platform Testing**: Windows, macOS, Linux compatibility

### Test Categories
1. **Unit Tests**: Individual component testing
   - Progressive disclosure logic (`internal/prompts/progressive_test.go`)
   - Complexity level parsing and blueprint selection
   - Help filtering and flag management
2. **Integration Tests**: Full project generation workflow
   - End-to-end CLI generation with different complexity levels
   - Blueprint compilation validation
3. **Blueprint Tests**: Validate all blueprints generate working code
   - CLI-simple (8 files) vs CLI-standard (29 files) generation
   - Logger integration testing across all complexity levels
4. **Logger Tests**: Test each logger implementation with various configurations
5. **Acceptance Tests (ATDD)**: User-focused behavior validation
   - Progressive disclosure acceptance criteria (`tests/acceptance/cli/progressive_disclosure_test.go`)
   - Help system behavior validation
   - Complexity flag validation and blueprint selection
6. **API Tests**: Web interface endpoint testing (Phase 3)
7. **CLI Tests**: Command-line interface testing
   - Flag parsing and validation
   - Help output formatting
   - Exit code verification

## Database Schema (Phases 3-4)

### Core Tables
- `users` - Authentication and user profiles
- `projects` - Generated project metadata and sharing
- `marketplace_blueprints` - Community-contributed blueprints
- `blueprint_ratings` - Blueprint reviews and ratings
- `analytics_events` - Usage analytics and metrics
- `api_keys` - API access management

### Key Relationships
- Users can create multiple projects
- Projects reference blueprints used for generation
- Blueprints can be rated and reviewed by users
- Analytics track blueprint usage patterns

## Security Considerations

### Input Validation
- **Sanitize all inputs**: Project names, module paths, template variables
- **Validate module paths**: Ensure proper Go module syntax
- **Blueprint security**: Scan community blueprints for malicious code
- **Path traversal protection**: Prevent blueprint files from accessing unauthorized paths

### Web Interface Security
- **Authentication**: JWT-based authentication for user sessions
- **CORS**: Proper cross-origin request handling
- **Rate limiting**: Prevent abuse of generation endpoints
- **Input sanitization**: All user inputs must be validated

## Performance Considerations

### Blueprint Engine Optimization
- **Blueprint caching**: Cache parsed blueprints in memory
- **Parallel generation**: Generate multiple files concurrently
- **Memory management**: Efficient handling of large projects
- **File I/O optimization**: Batch file operations when possible

### Web Interface Performance
- **Lazy loading**: Load blueprints and previews on demand
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

### Blueprint Generation Failures
- **Blueprint syntax errors**: Use `go generate ./...` to validate
- **Missing variables**: Ensure all template variables are defined
- **File conflicts**: Check for overlapping file destinations
- **Rollback mechanism**: Use recovery system for partial failures

### Performance Issues
- **Large projects**: Use streaming for file downloads
- **Memory usage**: Implement garbage collection for long-running sessions
- **Database queries**: Use connection pooling and query optimization
- **Blueprint loading**: Cache blueprints to avoid repeated parsing

### Cross-platform Compatibility
- **File paths**: Use `filepath.Join()` for cross-platform paths
- **Line endings**: Handle CRLF vs LF appropriately
- **Permissions**: Set proper file permissions on Unix systems
- **Shell commands**: Provide alternatives for different platforms

## Future Extensibility

### Plugin System (Phase 4+)
- **HashiCorp go-plugin**: Process isolation for custom generators
- **Blueprint marketplace**: Community blueprint sharing and discovery
- **Organization blueprints**: Private blueprint repositories for enterprises
- **Custom hooks**: Pre/post generation custom scripts

### Integration Opportunities
- **IDE Extensions**: VS Code, GoLand plugins
- **CI/CD Integration**: GitHub Actions, GitLab CI blueprints
- **Cloud Platforms**: AWS, GCP, Azure integration
- **Monitoring Tools**: Observability and metrics integration

This comprehensive project represents a significant advancement in Go developer tooling, filling critical gaps in the ecosystem while providing a foundation for future innovation.