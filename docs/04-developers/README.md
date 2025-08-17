# Developer Documentation

Welcome to the go-starter developer documentation. Here you'll find everything needed to contribute to, extend, or understand the internal workings of go-starter.

## 🛠️ Developer Guides

### 🏗️ [Development Guide](DEVELOPMENT.md)
**Local development environment and setup**
- Required tools and dependencies
- Building from source
- Running tests and benchmarks
- Debugging and profiling

### 📋 [Template Documentation](TEMPLATE_DOCUMENTATION.md)
**Blueprint development and template system**
- Creating new blueprints
- Template syntax and variables
- Conditional file generation
- Best practices for templates

### 🧪 [Testing Guide](TESTING_GUIDE.md)
**Testing strategy and implementation**
- Unit testing patterns
- Integration test framework
- ATDD (Acceptance Test-Driven Development)
- Performance testing and benchmarks

### 🔧 [CI Integration](CI_INTEGRATION.md)
**Continuous integration and deployment**
- CI/CD pipeline setup
- Automated testing strategies
- Release automation
- Quality gates and validation

### 🚀 [Phase 2 Enhancements](PHASE_2_ENHANCEMENTS.md) ✨ **NEW**
**Production-ready enterprise features**
- Enhanced blueprints (Microservice, Monolith, gRPC Gateway)
- Enterprise observability and monitoring
- Security and resilience patterns
- Performance optimization features

## 🚀 Quick Start for Contributors

### 1. Fork and Clone
```bash
# Fork the repository on GitHub, then:
git clone https://github.com/YOUR_USERNAME/go-starter.git
cd go-starter
```

### 2. Set Up Development Environment
```bash
# Install dependencies
go mod download

# Run tests to verify setup
make test

# Build the binary
make build
```

### 3. Make Your Changes
```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Make your changes
# Add tests for your changes
# Ensure all tests pass
make test
```

### 4. Submit Pull Request
```bash
# Push your changes
git push origin feature/your-feature-name

# Create pull request on GitHub
```

## 📁 Project Structure

```
go-starter/
├── cmd/                    # CLI commands and main entry points
├── internal/               # Private application code
│   ├── generator/          # Core project generation logic
│   ├── templates/          # Template loading and processing
│   ├── prompts/            # Interactive prompt system
│   ├── config/             # Configuration management
│   └── logger/             # Logging infrastructure
├── blueprints/             # Project templates
│   ├── cli-simple/         # Simple CLI blueprint
│   ├── cli-standard/       # Standard CLI blueprint
│   ├── web-api-*/          # Web API blueprints (various architectures)
│   └── lambda-*/           # Lambda blueprints
├── tests/                  # Test files and fixtures
│   ├── acceptance/         # ATDD tests
│   ├── integration/        # Integration tests
│   └── unit/               # Unit tests
├── docs/                   # Documentation
└── scripts/                # Development and CI scripts
```

## 🎯 Key Components

### Template Engine
**Location**: `internal/templates/`
**Purpose**: Load, parse, and process blueprint templates
**Key Files**:
- `loader.go` - Template loading from embedded filesystem
- `registry.go` - Blueprint registry and management
- `cache.go` - Template caching for performance

### Project Generator
**Location**: `internal/generator/`
**Purpose**: Core project generation logic
**Key Files**:
- `generator.go` - Main generation orchestration
- `parallel.go` - Parallel file generation
- `logging.go` - Generation progress and logging

### Progressive Disclosure System
**Location**: `internal/prompts/`
**Purpose**: Adaptive user interface based on experience level
**Key Files**:
- `progressive.go` - Core progressive disclosure logic
- `interfaces/types.go` - Type definitions and interfaces
- `bubbletea/` - BubbleTea prompter implementation
- `survey/` - Survey prompter implementation

### Configuration System
**Location**: `internal/config/`
**Purpose**: Configuration loading, validation, and management
**Key Files**:
- `config.go` - Configuration structures and loading
- `validation.go` - Configuration validation logic

## 🔧 Development Workflows

### Adding a New Blueprint

1. **Create Blueprint Directory**:
   ```bash
   mkdir blueprints/new-blueprint
   cd blueprints/new-blueprint
   ```

2. **Add Template Files**:
   ```
   blueprints/new-blueprint/
   ├── template.yaml        # Blueprint metadata
   ├── main.go.tmpl        # Template files
   ├── README.md.tmpl      # Documentation template
   └── Makefile.tmpl       # Build automation
   ```

3. **Update Registry**:
   ```go
   // internal/templates/registry.go
   // Add new blueprint to registry
   ```

4. **Add Tests**:
   ```bash
   # Create ATDD tests
   tests/acceptance/blueprints/new-blueprint/
   ```

5. **Update Documentation**:
   ```bash
   # Update blueprint guides
   docs/02-user-guides/blueprint-selection.md
   ```

### Adding a New Logger

1. **Create Logger Implementation**:
   ```go
   // blueprints/shared/logger/new-logger.go.tmpl
   ```

2. **Update Logger Factory**:
   ```go
   // blueprints/shared/logger/factory.go.tmpl
   ```

3. **Add Tests**:
   ```bash
   # Test logger functionality
   tests/unit/logger/
   ```

4. **Update Documentation**:
   ```bash
   # Update logger comparison
   docs/03-reference/logger-guide.md
   ```

### Modifying Progressive Disclosure

1. **Update Core Logic**:
   ```go
   // internal/prompts/progressive.go
   ```

2. **Modify Flag Categorization**:
   ```go
   // cmd/new.go - essentialFlags map
   ```

3. **Test Changes**:
   ```bash
   # Test progressive disclosure
   go test ./internal/prompts/
   tests/acceptance/cli/progressive_disclosure_test.go
   ```

## 🧪 Testing Strategy

### Test Categories

1. **Unit Tests** (`tests/unit/`):
   - Individual component testing
   - Fast execution, isolated tests
   - Mock external dependencies

2. **Integration Tests** (`tests/integration/`):
   - Component interaction testing
   - End-to-end generation workflows
   - Real filesystem operations

3. **Acceptance Tests** (`tests/acceptance/`):
   - ATDD (Acceptance Test-Driven Development)
   - User-focused behavior validation
   - Blueprint compilation verification

4. **Performance Tests** (`tests/performance/`):
   - Generation speed benchmarks
   - Memory usage profiling
   - Scalability testing

### Running Tests

```bash
# All tests
make test

# Specific test categories
go test ./internal/...           # Unit tests
go test ./tests/integration/...  # Integration tests
go test ./tests/acceptance/...   # ATDD tests

# With coverage
make test-coverage

# Watch mode for development
make test-watch
```

## 🏗️ Architecture Principles

### Design Principles

1. **Progressive Disclosure**: Adapt interface to user expertise
2. **Performance First**: Fast generation and startup times
3. **Production Ready**: All outputs must be deployment-ready
4. **Zero Dependencies**: Minimize external dependencies in generated code
5. **Cross-Platform**: Support Windows, macOS, Linux consistently

### Code Organization

- **Clear Separation**: Public API vs internal implementation
- **Interface Design**: Use interfaces for testability and extension
- **Error Handling**: Comprehensive error handling with helpful messages
- **Documentation**: Code should be self-documenting with examples

### Performance Considerations

- **Template Caching**: Cache parsed templates for repeated use
- **Parallel Generation**: Generate files concurrently when possible
- **Memory Efficiency**: Stream large operations, avoid loading everything in memory
- **Build Speed**: Optimize build times for development workflow

## 📚 Resources for Contributors

### Required Reading
- [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Testing](https://go.dev/doc/tutorial/add-a-test)

### Useful Tools
- **golangci-lint**: Code linting and formatting
- **dlv**: Go debugger for troubleshooting
- **go tool pprof**: Performance profiling
- **go tool trace**: Execution tracing

### Community Guidelines
- **Be Respectful**: Follow our code of conduct
- **Ask Questions**: Use discussions for clarification
- **Start Small**: Begin with small contributions
- **Test Thoroughly**: Include tests with all changes

## 🚀 Release Process

### Version Management
- **Semantic Versioning**: Major.Minor.Patch format
- **Release Notes**: Detailed changelog for each release
- **Backward Compatibility**: Maintain compatibility within major versions

### Release Steps
1. **Update Version**: Bump version in appropriate files
2. **Update Changelog**: Document all changes
3. **Run Full Test Suite**: Ensure all tests pass
4. **Create Release**: Tag and create GitHub release
5. **Update Documentation**: Reflect any changes in docs

---

**Ready to contribute?** Start with the [Contributing Guide](contributing.md) and join our community of developers making Go development better! 🚀

**Questions?** Join our [GitHub Discussions](https://github.com/francknouama/go-starter/discussions) or check the [FAQ](../02-user-guides/faq.md).