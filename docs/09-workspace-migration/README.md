# Go Workspace Migration Plan

## Executive Summary

This document outlines the migration plan for transitioning go-starter from a monolithic structure to a modular Go workspace architecture. The migration will separate core functionalities, CLI implementation, and Web UI into independent modules while maintaining backward compatibility.

**Project Tracking**: [GitHub Project #10](https://github.com/users/francknouama/projects/10)

## Motivation

### Current Challenges
- **Tight Coupling**: Core logic, CLI, and Web UI are intertwined
- **Testing Complexity**: Difficult to test components in isolation
- **Development Friction**: Changes in one area affect others
- **Deployment Constraints**: Cannot use components independently

### Benefits of Migration
- **Independent Development**: Teams can work on modules separately
- **Better Testing**: Module-level and integration tests
- **Clear Dependencies**: Explicit module boundaries
- **Flexible Deployment**: Use modules independently or together
- **Improved Maintainability**: Easier to understand and modify

## Target Architecture

```
go-starter/
├── go.work                    # Workspace configuration
├── go.work.sum               # Workspace checksums
├── modules/
│   ├── core/                 # Core functionality module
│   │   ├── go.mod           # github.com/francknouama/go-starter/core
│   │   ├── go.sum
│   │   ├── generator/        # Generation engine
│   │   ├── templates/        # Template system
│   │   ├── blueprints/       # Blueprint definitions
│   │   └── api/             # Public API interfaces
│   ├── cli/                  # CLI implementation module
│   │   ├── go.mod           # github.com/francknouama/go-starter/cli
│   │   ├── go.sum
│   │   ├── cmd/             # Cobra commands
│   │   ├── prompts/         # Interactive system
│   │   └── main.go         # CLI entry point
│   └── web/                  # Web UI module
│       ├── go.mod           # github.com/francknouama/go-starter/web
│       ├── go.sum
│       ├── frontend/        # React application
│       │   ├── package.json
│       │   └── src/
│       └── api/            # Go API server
│           └── server.go
├── tests/                    # Integration tests
│   ├── integration/         # Cross-module tests
│   └── e2e/                # End-to-end tests
└── scripts/                 # Build and deployment scripts
```

## Migration Phases

### Phase 1: Preparation (Issues #227-229)
**Duration**: 2-3 days

#### Tasks
1. **Architecture Analysis** (#227)
   - Analyze current codebase structure
   - Document all dependencies
   - Identify circular dependency risks
   - Create architecture decision records

2. **Dependency Mapping** (#228)
   - Generate complete dependency graph
   - Map internal package dependencies
   - Identify shared dependencies
   - Document version constraints

3. **Migration Strategy** (#229)
   - Develop phase-by-phase plan
   - Create rollback procedures
   - Define success criteria
   - Establish timeline with gates

### Phase 2: Core Module (Issues #230-231)
**Duration**: 2-3 days

#### Tasks
1. **Foundation Setup** (#230)
   - Create `modules/core` structure
   - Initialize go.mod
   - Move generator engine
   - Move template system
   - Define public API

2. **Blueprint System** (#231)
   - Move blueprints directory
   - Refactor blueprint registry
   - Implement validation framework
   - Add versioning support

#### Module Structure
```go
// modules/core/api/generator.go
package api

type Generator interface {
    Generate(config GeneratorConfig) error
}

type GeneratorConfig struct {
    ProjectName string
    BlueprintType string
    Options map[string]interface{}
}
```

### Phase 3: CLI Module (Issues #232-233)
**Duration**: 1-2 days

#### Tasks
1. **Module Separation** (#232)
   - Create `modules/cli` structure
   - Move cmd directory
   - Extract CLI-specific logic
   - Implement core dependency

2. **Interactive System** (#233)
   - Move prompts package
   - Refactor progressive disclosure
   - Update validation logic
   - Test interactive workflows

### Phase 4: Web Module (Issue #234)
**Duration**: 1-2 days

#### Tasks
1. **Frontend/Backend Separation** (#234)
   - Create `modules/web` structure
   - Move React application
   - Extract API handlers
   - Setup build pipeline

### Phase 5: Workspace Configuration (Issues #235, #240)
**Duration**: 1 day

#### Tasks
1. **Workspace Setup** (#235)
   - Create go.work file
   - Configure all modules
   - Setup development tools
   - Update Makefile

2. **CI/CD Update** (#240)
   - Update GitHub Actions
   - Configure matrix testing
   - Setup module deployments
   - Add integration tests

#### go.work Configuration
```go
go 1.23

use (
    ./modules/core
    ./modules/cli
    ./modules/web
)
```

### Phase 6: Validation (Issues #236-239)
**Duration**: 2-3 days

#### Tasks
1. **Backward Compatibility** (#236)
   - Create compatibility shim
   - Implement module aliases
   - Add deprecation warnings

2. **Testing Suite** (#237)
   - Module-level unit tests
   - Cross-module integration
   - Workspace validation

3. **Performance** (#238)
   - Benchmark comparisons
   - Memory usage analysis
   - Build time impacts

4. **Documentation** (#239)
   - Update README
   - Create migration guide
   - Update architecture docs

## Implementation Guidelines

### Module Boundaries
- **Core**: Pure business logic, no UI dependencies
- **CLI**: Cobra, prompts, terminal UI only
- **Web**: React, HTTP handlers, WebSocket

### Dependency Rules
- CLI and Web depend on Core
- Core has no dependencies on CLI or Web
- Shared types defined in Core
- No circular dependencies allowed

### Testing Strategy
```bash
# Module tests
cd modules/core && go test ./...
cd modules/cli && go test ./...
cd modules/web && go test ./...

# Integration tests
go test ./tests/integration/...

# Workspace validation
go work sync
go mod graph | check-circular-deps
```

### Backward Compatibility
```go
// Compatibility shim for old imports
// OLD: github.com/francknouama/go-starter/internal/generator
// NEW: github.com/francknouama/go-starter/core/generator

// go.mod replace directive during transition
replace github.com/francknouama/go-starter/internal => ./modules/core
```

## Success Criteria

### Phase Gates
Each phase must meet criteria before proceeding:

1. **Preparation Complete**
   - All dependencies documented
   - Migration plan approved
   - Rollback procedures defined

2. **Core Module Complete**
   - Core compiles independently
   - All tests pass
   - API documented

3. **CLI Module Complete**
   - CLI works with Core
   - Interactive mode functional
   - All commands tested

4. **Web Module Complete**
   - Web UI functional
   - API server working
   - Real-time preview operational

5. **Workspace Complete**
   - All modules integrated
   - CI/CD updated
   - Development workflow smooth

6. **Validation Complete**
   - Backward compatibility verified
   - Performance acceptable
   - Documentation updated

## Risk Mitigation

### Identified Risks
1. **Breaking Changes**: Mitigated by compatibility layer
2. **Performance Regression**: Continuous benchmarking
3. **Development Disruption**: Phased approach
4. **User Impact**: Backward compatibility maintained

### Rollback Plan
- Each phase can be rolled back independently
- Git branches for each phase
- Compatibility layer allows gradual migration
- Feature flags for new structure

## Timeline

| Phase | Start | End | Duration | Dependencies |
|-------|-------|-----|----------|--------------|
| Preparation | Day 1 | Day 3 | 2-3 days | None |
| Core Module | Day 4 | Day 6 | 2-3 days | Preparation |
| CLI Module | Day 7 | Day 8 | 1-2 days | Core Module |
| Web Module | Day 9 | Day 10 | 1-2 days | Core Module |
| Workspace | Day 11 | Day 11 | 1 day | All Modules |
| Validation | Day 12 | Day 15 | 2-3 days | Workspace |

**Total Duration**: 10-15 days

## Monitoring & Tracking

### Progress Tracking
- GitHub Project #10 for issue tracking
- Daily status updates
- Phase completion reports
- Blocker identification

### Quality Metrics
- Test coverage maintained >80%
- Zero breaking changes
- Performance within 5% of baseline
- All blueprints functional

## Post-Migration

### Documentation Updates
- Architecture guides
- API documentation
- Migration guide for users
- Contribution guidelines

### Communication Plan
- Announcement of migration start
- Phase completion updates
- Migration complete announcement
- User migration guide

## Appendix

### Related Issues
- #227: Architecture Analysis and Design
- #228: Dependency Mapping and Analysis
- #229: Migration Strategy and Rollback Plan
- #230: Core Module - Foundation Setup
- #231: Core Module - Blueprint System
- #232: CLI Module - Separation and Setup
- #233: CLI Module - Interactive System
- #234: Web Module - Frontend and Backend Separation
- #235: Workspace Configuration
- #236: Backward Compatibility Layer
- #237: Integration Testing Suite
- #238: Performance Validation
- #239: Documentation Update
- #240: CI/CD Pipeline Update

### Resources
- [Go Workspaces Documentation](https://go.dev/doc/tutorial/workspaces)
- [Go Module Migration Guide](https://go.dev/doc/modules/migration)
- [GitHub Project #10](https://github.com/users/francknouama/projects/10)