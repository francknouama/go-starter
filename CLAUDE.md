# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with this repository.

## Project Overview

**go-starter** - A comprehensive Go project generator combining create-react-app simplicity with Spring Initializr flexibility. Features progressive disclosure for beginners and advanced developers.

## 🎯 Interactive Mode (Primary Interface)

go-starter features a **powerful interactive mode** that guides you through project creation with intelligent prompts:

### Interactive Features
- **Smart prompts** - Only asks relevant questions based on your choices
- **Progressive disclosure** - Shows basic options by default, advanced when needed
- **Contextual help** - Each prompt includes helpful descriptions
- **Validation** - Ensures valid project names, module paths, and configurations
- **Preview mode** - See what will be generated before creation

### Interactive Usage
```bash
# Start interactive mode (recommended for all users)
go-starter new

# Interactive with advanced options
go-starter new --advanced

# Interactive for specific blueprint type
go-starter new --type=web-api  # Still prompts for remaining options
```

### Interactive Flow Example
1. **Project name** → Enter name or use random generator
2. **Blueprint type** → Choose from 12 templates with descriptions
3. **Complexity level** → Simple, Standard, Advanced, or Expert
4. **Framework selection** → Context-aware options (e.g., Gin/Echo for web, Cobra for CLI)
5. **Logger choice** → slog (default), zap, logrus, or zerolog
6. **Additional features** → Database, authentication, deployment options (advanced mode)

## Quick Reference

### Essential Commands
```bash
# Build & Install
go build -o bin/go-starter main.go
make build && make install

# Interactive Mode (Recommended)
go-starter new                          # Start guided project creation
go-starter new --advanced               # Interactive with all options

# Direct Mode (Skip Interactive)
go-starter new my-app --type=cli        # Non-interactive generation
go-starter new --help                   # Essential flags (14)
go-starter new --advanced --help        # All flags (18+)

# Testing
go test -v ./...                        # All tests
make test                               # Via Makefile

# Code Quality
golangci-lint run                       # Before commits
go generate ./...                       # Embed blueprints
```

### Progressive Disclosure Flags
- `--complexity [simple|standard|advanced|expert]` - Project complexity
- `--basic/--advanced` - Help mode toggle
- `--dry-run` - Preview without creating

### File Count by Complexity
- **Simple CLI**: 8 files (learning/prototypes)
- **Standard CLI**: 29 files (production-ready)
- **Web API**: Varies by architecture pattern

## 12 Core Blueprints

| Blueprint | Complexity | Files | Use Case |
|-----------|------------|-------|----------|
| **Simple CLI** | Beginner | 8 | Quick utilities, learning |
| **Standard CLI** | Intermediate | 29 | Production CLIs |
| **Library** | Beginner | ~10 | Reusable packages |
| **Standard Web API** | Intermediate | ~25 | REST APIs, CRUD |
| **Clean Architecture** | Advanced | ~40 | Enterprise APIs |
| **DDD Web API** | Advanced | ~40 | Complex domains |
| **Hexagonal** | Expert | 50+ | High testability |
| **AWS Lambda** | Beginner | ~15 | Serverless functions |
| **Lambda Proxy** | Intermediate | ~20 | API Gateway |
| **Event-Driven** | Expert | 50+ | CQRS/Event Sourcing |
| **Microservice** | Advanced | ~35 | gRPC services |
| **Monolith** | Intermediate | ~30 | Traditional apps |
| **Workspace** | Advanced | Varies | Multi-module |

## Architecture Components

### Project Structure
```
/cmd            # CLI commands (Cobra)
/internal       # Core logic
  /generator    # Project generation engine
  /prompts      # Interactive UI & progressive disclosure
  /templates    # Blueprint engine (embed.FS)
/blueprints     # Template definitions
/web            # React UI (Phase 3)
/tests          # ATDD & integration tests
```

### Progressive Disclosure System

#### Core Implementation
- **Location**: `internal/prompts/progressive.go`
- **Complexity Levels**: Simple → Standard → Advanced → Expert
- **Disclosure Modes**: Basic (14 flags) → Advanced (18+ flags)
- **Smart Defaults**: Auto-sets framework & logger for CLI projects

#### Blueprint Selection Logic
```go
// CLI blueprints adapt to complexity
--type=cli --complexity=simple  → cli-simple (8 files)
--type=cli --complexity=standard → cli (29 files)
```

## Development Workflow

### Adding New Blueprints
1. Create `blueprints/new-type/` directory
2. Add `template.yaml` with metadata
3. Create `.tmpl` files with Go templates
4. Update `internal/templates/registry.go`
5. Add prompts in `internal/prompts/interactive.go`
6. Write tests for validation

### Blueprint Variables
- `{{.ProjectName}}` - Project name
- `{{.ModulePath}}` - Go module path
- `{{.GoVersion}}` - Go version (default: 1.21)
- `{{.Framework}}` - Framework (gin, echo, cobra, etc.)
- `{{.LoggerType}}` - Logger (slog, zap, logrus, zerolog)
- `{{.Architecture}}` - Pattern (standard, clean, ddd, hexagonal)

### Conditional Generation
```yaml
files:
  - source: database.go.tmpl
    destination: internal/database/database.go
    condition: "{{ne .Features.Database.Driver \"\"}}"
```

## Logger System

### Supported Loggers
| Logger | Performance | Default | Use Case |
|--------|-------------|---------|----------|
| **slog** | Good | ✅ | Standard library |
| **zap** | Excellent | | High performance |
| **logrus** | Good | | Feature-rich |
| **zerolog** | Excellent | | Zero allocation |

### Complexity Reduction
- CLI-Standard: 91% reduction (1,051 → 98 lines)
- Web-API: 72% reduction (398 → 110 lines)

## Testing Strategy

### Test Categories
1. **Unit Tests** - Component logic
2. **Integration Tests** - Full workflow
3. **Blueprint Tests** - Generation validation
4. **ATDD Tests** - User behavior validation
5. **Cross-platform** - Windows/macOS/Linux

### Critical Requirements
- All blueprints must parse without errors
- Generated projects must compile (`go build`)
- All logger types must work
- Tests run on all platforms

### ATDD Infrastructure
- Dynamic project root detection
- Cross-platform compatibility
- Template discovery from `blueprints/`
- Compilation validation
- Performance monitoring

## Common Issues & Solutions

### Progressive Disclosure
```bash
# Too many options for beginners?
go-starter new --help  # Already defaults to basic

# Can't find advanced options?
go-starter new --advanced --help

# Too many files for simple project?
go-starter new my-tool --type=cli --complexity=simple
```

### Blueprint Generation
- **Syntax errors**: Run `go generate ./...`
- **Missing variables**: Check template definitions
- **File conflicts**: Verify destinations don't overlap

## Agent Ecosystem

### Core Development
- **golang-fullstack-engineer** - Go development, ATDD, refactoring
- **ux-design-expert** - Web/mobile UI, user experience
- **senior-bug-resolver** - QA bugs, critical issues
- **cross-platform-tester** - Platform compatibility

### Production Readiness
- **accessibility-ux-specialist** - WCAG compliance
- **devops-deployment-specialist** - Infrastructure, CI/CD
- **performance-security-specialist** - Optimization, security
- **documentation-community-specialist** - Docs, tutorials

### Project Management
- **product-owner** - GitHub issues, roadmap
- **general-purpose** - Multi-domain tasks

### Agent Selection
1. Match primary skill area
2. Consider production requirements
3. Use general-purpose for complex multi-domain tasks
4. Prioritize senior-bug-resolver for critical issues

## Phase Implementation

### Phase 1: Core CLI ✅
- Basic blueprints (Web API, CLI, Library, Lambda)
- Cobra framework, interactive prompts
- Blueprint engine

### Phase 2: Complete System ✅
- All 12 project types
- Multiple architecture patterns
- Progressive disclosure

### Phase 3: Web UI (In Progress)
- React + Vite frontend
- WebSocket real-time preview
- RESTful API

### Phase 4: Advanced Features (Future)
- GitHub integration
- Blueprint marketplace
- Cloud deployments

## Configuration

### CLI Config (`~/.go-starter.yaml`)
```yaml
profiles:
  default:
    author: "John Doe"
    defaults:
      goVersion: "1.21"
      framework: "gin"
      logger: "slog"
```

### Project Config (`project.yaml`)
```yaml
name: my-api
module: github.com/user/my-api
type: web-api
architecture: hexagonal
framework: gin
logger: slog
features:
  database:
    driver: postgres
    orm: gorm
```

## Security & Performance

### Security
- Input sanitization (names, paths, templates)
- Module path validation
- Blueprint security scanning
- Path traversal protection

### Performance
- Blueprint caching in memory
- Parallel file generation
- Efficient I/O batching
- WebSocket debouncing

## Quick Troubleshooting

| Problem | Solution |
|---------|----------|
| Interactive prompts with flags | Include all required flags |
| Help shows too many options | Use default basic mode |
| Too many files generated | Use `--complexity=simple` |
| Blueprint syntax errors | Run `go generate ./...` |
| Cross-platform issues | Use `filepath.Join()` |

## Important Reminders

- NEVER create files unless necessary
- ALWAYS prefer editing over creating new files
- NEVER create docs unless explicitly requested
- Follow existing code conventions
- Check for existing libraries before adding new ones
- Test with `golangci-lint` before commits
- Mark todos as completed immediately in TodoWrite