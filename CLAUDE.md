# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with this repository.

## Project Overview

**go-starter** - A comprehensive Go project generator combining create-react-app simplicity with Spring Initializr flexibility. Features progressive disclosure for beginners and advanced developers with an interactive CLI interface.

## 🚀 Go Workspace Migration (In Progress)

### Migration Overview
The project is undergoing a migration to Go workspace structure to separate core functionalities, CLI, and Web UI into independent modules. This will improve modularity, testing, and allow independent development of each component.

**GitHub Project**: [#10 - Go Workspace Migration](https://github.com/users/francknouama/projects/10)

### Target Architecture
```
go-starter/
├── go.work                    # Workspace configuration
├── modules/
│   ├── core/                  # Core functionality
│   │   ├── go.mod
│   │   ├── generator/         # Generation engine
│   │   ├── templates/         # Template system
│   │   └── blueprints/        # Blueprint definitions
│   ├── cli/                   # CLI implementation
│   │   ├── go.mod
│   │   ├── cmd/               # Cobra commands
│   │   └── prompts/           # Interactive system
│   └── web/                   # Web UI & API
│       ├── go.mod
│       ├── frontend/          # React application
│       └── api/               # Go API server
└── tests/                      # Integration tests
```

### Migration Phases & Issues
- **Phase 1: Preparation** - #227, #228, #229
- **Phase 2: Core Module** - #230, #231
- **Phase 3: CLI Module** - #232, #233
- **Phase 4: Web Module** - #234
- **Phase 5: Workspace** - #235, #240
- **Phase 6: Validation** - #236, #237, #238, #239

**Timeline**: 10-15 days (estimated)

## 🎯 Interactive CLI System

### 💻 CLI Interactive Mode

go-starter features a **powerful interactive CLI mode** that guides you through project creation with intelligent prompts:

### Interactive Features
- **Smart prompts** - Only asks relevant questions based on your choices
- **Progressive disclosure** - Shows basic options by default, advanced when needed
- **Contextual help** - Each prompt includes helpful descriptions
- **Validation** - Ensures valid project names, module paths, and configurations
- **Preview mode** - See what will be generated before creation

### Usage
```bash
# Interactive Mode (Recommended)
go-starter new                              # Start guided project creation
go-starter new --advanced                   # Interactive with all options
go-starter new --type=web-api               # Interactive for specific blueprint

# Direct Mode (Skip interactive)
go-starter new my-app --type=cli            # Non-interactive generation
```

### CLI Interactive Flow
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

# Interactive Mode
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

### CLI Features
- **Simple CLI**: 8 files (learning/prototypes)
- **Standard CLI**: 29 files (production-ready)
- **Web API**: Varies by architecture pattern

## Production-Ready Blueprints - 12 PRODUCTION-READY ✨ **100% COVERAGE!**

**ACHIEVEMENT**: Complete CLI system with all 12 production-ready blueprints

| Blueprint | Complexity | Files | Use Case | Production Features |
|-----------|------------|-------|----------|------------------------------|
| **CLI Simple** ✅ | Beginner | 8 | Quick utilities, learning | Basic logging, Makefile |
| **CLI Standard** ✅ | Intermediate | 29 | Production CLIs | Full CLI framework, tests |
| **Library Standard** ✅ | Beginner | ~10 | Reusable packages | API design, examples |
| **Web API Standard** ✅ | Intermediate | ~25 | REST APIs, CRUD | HTTP framework, middleware |
| **Web API Clean** ✅ | Advanced | ~40 | Enterprise APIs | Layered architecture, DI |
| **Web API Echo** ✅ **PHASE 3A** | Intermediate | ~25 | Echo REST APIs | High-performance middleware |
| **Web API Fiber** ✅ **PHASE 3A** | Intermediate | ~25 | Fiber REST APIs | Ultra-fast performance |
| **Lambda Standard** ✅ | Beginner | ~15 | Serverless functions | AWS SDK, X-Ray tracing |
| **Lambda Proxy** ✅ **PHASE 3A** | Intermediate | ~20 | API Gateway | HTTP routing, serverless |
| **gRPC Gateway** ✅ | Advanced | 45 | **Dual HTTP/gRPC APIs** | **🚀 Enhanced interceptors, unified middleware** |
| **Monolith** ✅ **PHASE 3B** | Intermediate | 72 | **Production web apps** | **🚀 Background jobs, multi-layer caching, performance monitoring** |
| **Microservice** ✅ **PHASE 3B** | Advanced | 47 | **Enterprise gRPC services** | **🚀 OpenTelemetry, rate limiting, resilience patterns** |

**ACHIEVEMENT**: Achieved **12 production-ready blueprints** with complete CLI system. This represents the completion of the Extended Blueprint System (Phase 3) with progressive disclosure capabilities.

## Architecture Components

### Project Structure
```
/cmd            # CLI commands (Cobra)
/internal       # Core logic
  /generator    # Project generation engine
  /prompts      # Interactive UI & progressive disclosure
  /templates    # Blueprint engine (embed.FS)
/blueprints     # Template definitions
/web            # React UI (Under Development)
/tests          # ATDD & integration tests
```

### CLI Progressive Disclosure
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

### Specialized Agent Definitions

#### blueprint-production-tracker
**Purpose**: Track production readiness status across all blueprints with precision
**When to Use**: After blueprint fixes, during validation, for status updates, milestone tracking

**Core Expertise**:
- Blueprint status management and classification (Production Ready, Enhancement Ready, Development Phase)
- Production readiness metrics and criteria validation
- Progress tracking across multiple blueprint fixes
- ATDD test result interpretation and status updates
- File count validation and architectural compliance checking

**Primary Functions**:
- Update blueprint status in all tracking documents
- Validate production readiness criteria (compilation, logger integration, template variables)
- Track fix completion and validation results
- Maintain accurate counts of production-ready vs needs-fixes blueprints
- Generate status reports and progress summaries

**Key Files to Monitor/Update**:
- `BLUEPRINT_PRODUCTION_READINESS_REPORT.md`
- `docs/02-user-guides/BLUEPRINT_STATUS_GUIDE.md`
- `CLAUDE.md` (blueprint status table)

#### documentation-sync-manager
**Purpose**: Keep documentation synchronized with actual project status and progress
**When to Use**: After major changes, before releases, during status updates, when documentation drift is detected

**Core Expertise**:
- README.md maintenance and status synchronization
- Status guide updates to reflect current blueprint states
- Changelog management and version documentation
- Cross-document consistency validation
- Documentation accuracy verification against actual code/tests

**Primary Functions**:
- Update README.md to reflect current blueprint counts and status
- Synchronize status guides with actual validation results
- Maintain consistent messaging across all documentation files
- Generate changelog entries for significant progress
- Validate documentation accuracy against actual project state

**Key Files to Monitor/Update**:
- `README.md`
- `docs/02-user-guides/BLUEPRINT_STATUS_GUIDE.md`
- `docs/02-user-guides/BLUEPRINT_SELECTION_GUIDE.md`
- `CHANGELOG.md` (if exists)

#### github-project-manager
**Purpose**: Manage GitHub issues, projects, and release coordination activities
**When to Use**: For issue tracking, milestone management, release preparation, project board updates

**Core Expertise**:
- GitHub Issues API and project board management
- Release preparation and coordination
- Milestone tracking and issue lifecycle management
- Label management and issue categorization
- Project board automation and workflow management

**Primary Functions**:
- Create and update GitHub issues for blueprint fixes
- Manage project board status and progress tracking
- Coordinate release preparation activities
- Track issue resolution and milestone completion
- Generate release notes and status updates

**Key Responsibilities**:
- Blueprint fix tracking issues
- Production readiness milestone management
- Release preparation coordination
- Community communication preparation

#### progress-coordinator
**Purpose**: Coordinate multi-agent workflows and track overall project progress
**When to Use**: For complex multi-step tasks, workflow orchestration, cross-agent coordination

**Core Expertise**:
- Agent workflow coordination and task delegation
- Multi-step process management and tracking
- Progress reporting and status aggregation
- Workflow optimization and bottleneck identification
- Cross-functional task coordination

**Primary Functions**:
- Orchestrate multi-agent workflows for complex tasks
- Coordinate between development and tracking agents
- Aggregate progress from multiple work streams
- Identify and resolve workflow bottlenecks
- Generate comprehensive progress reports

**Key Coordination Areas**:
- Blueprint fix workflows (development → validation → documentation → release)
- Multi-agent collaboration for complex features
- Release preparation workflows
- Quality assurance coordination

#### performance-security-specialist
**Purpose**: Security vulnerability remediation, performance optimization, and code security analysis
**When to Use**: For Dependabot alerts, security vulnerabilities, performance bottlenecks, security code reviews, dependency updates, and security compliance assessments

**Core Expertise**:
- Security vulnerability assessment and remediation
- Dependency management and security updates
- Performance profiling and optimization
- Security code review and static analysis
- OWASP security best practices implementation
- CVE analysis and mitigation strategies

**Primary Functions**:
- Analyze and fix Dependabot security alerts
- Update vulnerable dependencies to secure versions
- Perform security code reviews and audits
- Optimize application performance and resource usage
- Implement security headers and middleware
- Configure secure deployment practices

**Security Specialties**:
- **Dependency Security**: Go modules, npm packages, security updates
- **Code Security**: Input validation, injection prevention, secrets management
- **Infrastructure Security**: Docker security, CI/CD pipeline security
- **Performance**: Memory optimization, CPU profiling, load testing
- **Compliance**: Security standards, vulnerability scanning

**Technical Focus Areas**:
- Go security best practices and secure coding
- Web application security (OWASP Top 10)
- Container and deployment security
- Performance monitoring and optimization
- Security automation in CI/CD pipelines

#### ascii-art-designer
**Purpose**: Create professional ASCII art logos, banners, and decorative elements for CLI applications
**When to Use**: When creating ASCII art for terminal display, logo design, splash screens, documentation headers, or visual branding elements

**Core Expertise**:
- ASCII art creation for various terminal widths (80-120 characters)
- Logo and banner design using standard ASCII characters
- Terminal-safe character selection (no extended/Unicode unless requested)
- Multiple style variations (minimalist, elaborate, themed)
- Go string literal formatting for easy code embedding

**Primary Functions**:
- Create main ASCII logos for CLI applications
- Design banner versions for documentation headers
- Generate compact versions for constrained spaces
- Develop themed variants (e.g., Gopher for Go projects)
- Provide multiple style options for user selection

**Design Specialties**:
- **Main Logos**: Professional branding for CLI tools
- **Banner Headers**: Wide-format designs for documentation
- **Compact Versions**: Single/multi-line for prompts
- **Themed Art**: Language/framework specific designs
- **Box Drawing**: Terminal UI elements and borders

**Technical Considerations**:
- Cross-platform terminal compatibility
- Proper escaping for embedding in code
- Monospace font optimization
- Width constraints for different terminal sizes
- Raw string formatting for Go embedding

### Core Development
- **golang-fullstack-engineer** - Go development, ATDD, refactoring
- **ux-design-expert** - Web/mobile UI, user experience
- **senior-bug-resolver** - QA bugs, critical issues
- **cross-platform-tester** - Platform compatibility
- **devops-cicd-specialist** - CI/CD pipelines, GitHub Actions, automation

### Production Readiness
- **accessibility-ux-specialist** - WCAG compliance
- **devops-deployment-specialist** - Infrastructure, CI/CD
- **performance-security-specialist** - Optimization, security
- **documentation-community-specialist** - Docs, tutorials

### Project Management
- **product-owner** - GitHub issues, roadmap
- **general-purpose** - Multi-domain tasks

### Design & UX
- **ascii-art-designer** - ASCII art logos, banners, terminal UI elements

### Documentation & Content Management
- **documentation-specialist** - Comprehensive documentation architecture, technical writing excellence, multi-audience content design

### Progress Tracking & Coordination
- **blueprint-production-tracker** - Blueprint status management, production readiness metrics, progress tracking
- **documentation-sync-manager** - README updates, status guides, changelog management
- **github-project-manager** - GitHub issues, project boards, release coordination
- **progress-coordinator** - Multi-agent workflows, workflow orchestration, progress reporting

### Agent Selection Strategy
1. **Match primary skill area** - Choose agents whose expertise directly addresses the task
2. **Consider production requirements** - Use production-focused agents for critical systems
3. **Use general-purpose for multi-domain tasks** - Complex tasks requiring multiple skillsets
4. **Prioritize senior-bug-resolver for critical issues** - Emergency fixes and production problems
5. **Leverage hooks for intelligent selection** - Use Claude Code hooks for context-aware recommendations
6. **Consider workflow coordination** - Some agents work well together in sequences

### Hook Integration & Automation

#### Available Hooks System
The project includes sophisticated hooks in `.claude-hooks/` that enhance agent coordination:

**Agent Selector Hook** (`agent-selector.js`):
- **Intelligent Selection**: Automatically suggests optimal agents based on file context, user queries, and project phase
- **Confidence Scoring**: Provides confidence ratings and alternative agent suggestions  
- **Trigger Patterns**: Context-aware matching for common scenarios

**Context Analyzer Hook** (`context-analyzer.js`):
- **Project Phase Detection**: Identifies current phase (1-4) and completion status
- **Technology Stack Analysis**: Detects frameworks, tools, and architectural patterns
- **Health Assessment**: Provides project health scores and priority recommendations
- **Architecture Detection**: Recognizes Clean Architecture, DDD, Hexagonal patterns

**Agent Coordinator Hook** (`agent-coordinator.js`):
- **Multi-Agent Orchestration**: Coordinates complex workflows requiring multiple agents
- **Task Planning**: Creates execution plans with dependencies and risk assessment
- **Parallel Processing**: Identifies tasks that can run concurrently
- **Progress Tracking**: Monitors multi-agent workflows and results aggregation

#### Hook-Enhanced Agent Workflows

**Context-Aware Agent Selection**:
```javascript
// Hook automatically detects:
// - File: blueprints/grpc-gateway/internal/server/grpc.go.tmpl
// - Query: "fix protobuf compilation errors"
// → Suggests: grpc-protobuf-specialist (95% confidence)
// → Alternative: template-variable-auditor (78% confidence)
```

**Multi-Agent Coordination**:
```javascript
// Hook detects complex blueprint fix workflow:
// → Phase 1: grpc-protobuf-specialist (protobuf issues)
// → Phase 2: template-variable-auditor (template validation) 
// → Phase 3: blueprint-validator (compilation check)
// → Phase 4: general-purpose (documentation update)
```

#### Hook Configuration Triggers

**File Pattern Triggers**:
- `blueprints/**/*.tmpl` → `template-variable-auditor` or `blueprint-validator`
- `blueprints/grpc-*/**` → `grpc-protobuf-specialist`
- `web/src/**/*.tsx` → `ux-design-expert`
- `tests/**/*_test.go` → `golang-atdd-qa-engineer`
- `docs/**/*.md` → `general-purpose` (documentation)

**Query Pattern Triggers**:
- "fix compilation", "template errors" → `template-variable-auditor`
- "protobuf", "grpc", "buf.yaml" → `grpc-protobuf-specialist`
- "performance", "optimization" → `performance-auditor`
- "cross platform", "windows", "macos" → `cross-platform-tester`
- "production ready", "pipeline" → `blueprint-production-pipeline`

#### Hook-Driven Workflow Examples

**Blueprint Fix with Hook Intelligence**:
```bash
# User: "Fix the grpc-gateway compilation issues"
# Hook Analysis:
# - Context: grpc files, compilation errors, template issues
# - Confidence: grpc-protobuf-specialist (95%)
# - Workflow: 4-phase coordinated approach

1. grpc-protobuf-specialist → Fix protobuf and buf configuration  
2. template-variable-auditor → Validate template variables
3. blueprint-validator → Confirm compilation success
4. general-purpose → Update status documentation
```

**Performance Analysis with Context Awareness**:
```bash
# User: "The blueprint generation is slow"
# Hook Analysis:
# - Context: performance issue, generation pipeline
# - Phase Detection: Phase 2 (production features)
# - Health Score: 8.5/10 (good, but needs optimization)

1. performance-auditor → Analyze blueprint generation performance
2. blueprint-production-pipeline → Review production pipeline efficiency
3. cross-platform-tester → Validate performance across platforms
```

### Agent Usage Guidelines

#### Core Development Agents
Use for implementation, refactoring, and technical work:
- **golang-fullstack-engineer**: Complex Go implementation, ATDD testing, architectural decisions
- **senior-bug-resolver**: Critical bug fixes, production issues, emergency fixes
- **cross-platform-tester**: Platform compatibility, CI/CD pipeline issues
- **performance-security-specialist**: Security vulnerabilities, performance optimization

#### Progress Tracking Agents
Use for coordination, status tracking, and documentation:
- **blueprint-production-tracker**: When you need to update blueprint status, track fixes, or assess production readiness
- **documentation-sync-manager**: When documentation needs to be updated to reflect current project status
- **github-project-manager**: When you need to create/update GitHub issues, manage project boards, or coordinate releases
- **progress-coordinator**: For complex multi-step tasks requiring coordination between multiple agents

#### Design & UX Agents
Use for visual design, branding, and user interface elements:
- **ascii-art-designer**: When you need ASCII art for CLI logos, banners, splash screens, or terminal UI elements

#### Documentation & Content Management Agents
Use for comprehensive documentation creation, architecture, and maintenance:
- **documentation-specialist**: When you need comprehensive documentation creation, technical writing excellence, multi-audience content design, information architecture, or major documentation updates

### Current Blueprint Status (as of 2025-08-25)

#### Production Ready Blueprints (12 TOTAL - HISTORIC 100% COVERAGE ACHIEVED! 🎉🎉🎉)
- **cli-simple**: 8 files, basic CLI with logging
- **cli-standard**: 29 files, full CLI framework with tests
- **web-api-standard**: ~25 files, REST APIs with middleware
- **lambda-standard**: ~15 files, AWS serverless functions
- **library-standard**: ~10 files, reusable Go packages
- **web-api-clean**: ~40 files, clean architecture pattern
- **grpc-gateway**: 45 files, dual HTTP/gRPC APIs
- **lambda-proxy**: ~20 files, API Gateway integration (Phase 3A)
- **web-api-echo**: ~25 files, Echo framework REST APIs (Phase 3A)
- **web-api-fiber**: ~25 files, Fiber framework REST APIs (Phase 3A)
- **monolith**: 72 files, full-stack web applications ✅ **PHASE 3B NEWLY PRODUCTION READY**
- **microservice-standard**: 47 files, enterprise gRPC services ✅ **PHASE 3B NEWLY PRODUCTION READY**

#### Enhancement-Ready Blueprints
✅ **ALL ENHANCEMENT-READY BLUEPRINTS NOW PRODUCTION-READY!**
- **lambda-event-processing**: ✅ Template syntax fixes completed
- **monolith**: ✅ Module resolution fixes completed  
- **microservice-standard**: ✅ Template processing fixes completed

#### Agent-Managed Tracking Strategy

**Status Tracking Agents**:
- `blueprint-production-pipeline` → Orchestrate production readiness workflows
- `blueprint-validator` → Validate individual blueprint compilation
- `general-purpose` → Aggregate status across multiple blueprints and update tracking documents

**Quality Assurance Agents**:
- `performance-auditor` → Monitor blueprint generation performance and generated code quality
- `cross-platform-tester` → Ensure compatibility across all target platforms
- `atdd-test-creator` → Maintain comprehensive test coverage

**Documentation Maintenance Agents**:
- `general-purpose` → Keep README.md and status guides synchronized
- `product-owner` → Update roadmap and milestone documentation
- `open-source-community-manager` → Maintain user-facing guides and selection documentation

#### Key Tracking Documents & Responsible Agents
- `README.md` → `general-purpose` (status updates), `product-owner` (roadmap)
- `docs/02-user-guides/BLUEPRINT_STATUS_GUIDE.md` → `general-purpose` (status sync), `open-source-community-manager` (user guidance)
- `docs/02-user-guides/BLUEPRINT_SELECTION_GUIDE.md` → `complexity-analyzer` (recommendations), `open-source-community-manager` (user experience)
- `CLAUDE.md` → `general-purpose` (coordination updates), `product-owner` (workflow documentation)

#### Critical grpc-gateway Issues (RESOLVED ✅)
1. **Go Module Dependencies** (High Priority) → ✅ FIXED by `grpc-protobuf-specialist`
2. **Protobuf Integration** (High Priority) → ✅ FIXED by `grpc-protobuf-specialist` + `template-variable-auditor`
3. **Template Dependencies** (Medium Priority) → ✅ FIXED by `template-variable-auditor` + `blueprint-validator`
4. **Compilation Validation** → ✅ CONFIRMED by `blueprint-validator`

### Recent Achievements & Hook-Enhanced Agent Workflow Victory 🚀
- **Phase 2 Completion**: Enhanced microservice, monolith, and gRPC gateway blueprints with enterprise features
- **Production Features**: OpenTelemetry tracing, security middleware, resilience patterns
- **Quality Improvements**: Comprehensive testing and performance optimization across blueprints
- **Agent Integration**: Mapped project-specific functions to available Claude Code agents for improved workflow coordination
- **Hook-Enhanced Agent Workflow Success**: Multi-specialist coordination successfully resolved complex grpc-gateway technical issues
  - **grpc-protobuf-specialist** (Round 1): Fixed major gRPC dependencies and protobuf configuration
  - **template-variable-auditor**: Validated all template variables are correct and syntactically sound  
  - **golang-fullstack-engineer**: Fixed 6 critical template compilation bugs
  - **grpc-protobuf-specialist** (Round 2): Definitively resolved remaining gRPC dependency issues
  - **blueprint-validator**: Confirmed the blueprint now generates projects that can compile
- **HISTORIC MILESTONE ACHIEVEMENT**: ALL BLUEPRINTS NOW PRODUCTION-READY ✅ (12/12 total production-ready blueprints - 100% COVERAGE ACHIEVED!)

## Phase Implementation

### Phase 1: Core CLI ✅
- Basic blueprints (Web API, CLI, Library, Lambda)
- Cobra framework, interactive prompts
- Blueprint engine

### Phase 2: Enterprise Production Features ✅
- **Enhanced Blueprints**: Microservice, Monolith, gRPC Gateway now production-ready
- **Observability**: OpenTelemetry tracing, Prometheus metrics, structured logging
- **Security**: Input validation, rate limiting, security headers, CORS
- **Resilience**: Circuit breakers, retry logic, graceful error handling
- **Performance**: Multi-layer caching, connection pooling, resource management
- **Background Processing**: Comprehensive job manager with queuing and monitoring
- **Enterprise Middleware**: Enhanced interceptors with monitoring and security

### Phase 3: Extended Blueprint System (✅ Complete)
- React + Vite frontend
- WebSocket real-time preview
- RESTful API

### Phase 4: Web UI (Under Development)
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

## Agent Coordination Workflows

### Blueprint Fix Workflow
When fixing a blueprint (e.g., grpc-gateway):
1. **golang-fullstack-engineer**: Implement the technical fixes
2. **blueprint-production-tracker**: Update status tracking and validate fixes
3. **documentation-sync-manager**: Update all documentation to reflect current status
4. **github-project-manager**: Update relevant issues and project boards

### Status Update Workflow
When project status changes:
1. **blueprint-production-tracker**: Update core status tracking documents
2. **documentation-sync-manager**: Ensure README and guides are synchronized
3. **progress-coordinator**: Generate comprehensive progress reports if needed

### Release Preparation Workflow
When preparing for releases:
1. **blueprint-production-tracker**: Validate production readiness status
2. **documentation-sync-manager**: Ensure all documentation is current
3. **github-project-manager**: Coordinate release issues and milestones
4. **progress-coordinator**: Orchestrate the full release workflow

### Multi-Agent Task Coordination
For complex tasks requiring multiple agents:
1. **progress-coordinator**: Plan and orchestrate the workflow
2. Specialized agents: Execute their specific responsibilities
3. **progress-coordinator**: Track progress and generate status reports
4. All tracking agents: Update their respective tracking documents

### Documentation Enhancement Workflow
When comprehensive documentation updates are needed:
1. **documentation-specialist**: Create comprehensive guides, technical content, and information architecture
2. **blueprint-production-tracker**: Provide current status data and validation results
3. **documentation-sync-manager**: Ensure cross-document consistency and status synchronization
4. **github-project-manager**: Coordinate documentation releases with project milestones

### Major Feature Documentation Workflow
When documenting new features (e.g., Phase 2 enhancements):
1. **documentation-specialist**: Design comprehensive documentation structure and create detailed guides
2. **golang-fullstack-engineer**: Provide technical implementation details and examples
3. **performance-security-specialist**: Contribute security and performance best practices
4. **documentation-sync-manager**: Ensure all affected documents are updated consistently

## Important Reminders

### Core Development Guidelines
- NEVER create files unless necessary
- ALWAYS prefer editing over creating new files
- NEVER create docs unless explicitly requested
- Follow existing code conventions
- Check for existing libraries before adding new ones
- Test with `golangci-lint` before commits
- Mark todos as completed immediately in TodoWrite

### Hook-Enhanced Workflows
- **Leverage hooks for intelligent agent selection** - Let the hook system suggest optimal agents based on context
- **Use context-aware workflows** - Hooks provide project phase detection and health analysis
- **Trust hook confidence scores** - Higher confidence suggestions (>90%) are typically most accurate
- **Consider hook alternatives** - Review alternative agent suggestions for complex scenarios
- **Coordinate multi-agent tasks** - Use agent coordinator hook for workflows requiring multiple specialists
- **Monitor hook feedback** - Pay attention to hook suggestions to improve workflow efficiency

### Agent Coordination Best Practices  
- Use tracking agents to maintain accurate status across all documents
- Coordinate between development and documentation agents for comprehensive updates
- Leverage specialized agents (grpc-protobuf-specialist, template-variable-auditor) for specific domains
- Use general-purpose agent for complex multi-domain coordination tasks
- Follow hook-suggested workflows for optimal agent orchestration
