# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [2.2.0] - 2025-08-25 🎉 MAJOR MILESTONE ACHIEVED

### Added - **7 Production-Ready Blueprints Milestone**

#### 🚀 gRPC Gateway Blueprint - Now Production Ready
- **Complete dual-protocol API support**: HTTP REST + gRPC endpoints
- **Enhanced interceptors**: Monitoring, security, and performance interceptors  
- **Unified middleware**: Consistent middleware across HTTP and gRPC protocols
- **Production features**: Rate limiting, metrics collection, OpenTelemetry tracing
- **Protocol buffer integration**: Full protobuf support with code generation
- **Gateway pattern implementation**: Production-ready API gateway architecture

#### 🎯 Milestone Achievement
- **7 Production-Ready Blueprints**: CLI-Simple, CLI-Standard, Web-API-Standard, Lambda-Standard, Library-Standard, Web-API-Clean, gRPC-Gateway
- **Comprehensive ATDD Validation**: All production blueprints validated with cross-platform testing
- **Hook-Enhanced Agent Workflows**: Proved highly effective for complex technical issue resolution
- **Blueprint Status Tracking**: Comprehensive status guide for user transparency

### Enhanced
- **Documentation Ecosystem**: Updated all user-facing guides to reflect milestone achievement
- **Community Messaging**: Enhanced communication highlighting 7 production-ready blueprints
- **Blueprint Selection Guide**: Updated with comprehensive gRPC Gateway documentation
- **Status Guide Accuracy**: Real-time status tracking across all documentation

### Technical Achievements
- **gRPC Dependencies Resolution**: Fixed critical Go module dependency issues
- **Protobuf Integration**: Resolved buf.yaml version compatibility problems  
- **Template Processing**: Enhanced template engine for complex blueprint generation
- **Cross-Platform Validation**: Windows, macOS, Linux compatibility confirmed

### Community Impact
- **Production Readiness**: Users can now build enterprise dual HTTP/gRPC APIs
- **Microservice Support**: Complete enterprise microservice development stack
- **Gateway Services**: Automated protocol translation and unified endpoints
- **Development Velocity**: Hook-enhanced workflows significantly improved issue resolution time

## [1.1.0] - 2025-06-25

### Added
- **Multi-Database Selection Support**: Projects can now select multiple databases (PostgreSQL, MySQL, MongoDB, SQLite, Redis) during generation
- **Dynamic Go Version Detection**: Automatically detects and uses the developer's current Go version for generated projects
- **Improved CLI Behavior**: Graceful handling when Go is not installed, with helpful warnings and installation instructions

### Changed
- Docker Compose generation now supports multi-service configurations for selected databases
- Go version in generated projects now defaults to the detected version instead of hardcoded 1.21

### Fixed
- Prevented hard failures during dependency installation when Go is not installed
- Improved backward compatibility for existing single-database projects

## [1.0.0] - 2024-06-24

### Added

#### 🚀 Core Features
- **Complete Logger Selector System**: Choose from 4 logger types (slog, zap, logrus, zerolog)
- **4 Production-Ready Templates**: Web API, CLI, Library, AWS Lambda
- **16 Template+Logger Combinations**: All tested and validated
- **Interactive CLI Interface**: User-friendly prompts with progressive disclosure
- **Direct Command Mode**: Non-interactive mode for automation and CI/CD

#### 🏗️ Template System
- **Web API Template**: Gin framework with database integration, authentication, Docker support
- **CLI Application Template**: Cobra framework with configuration management
- **Go Library Template**: Clean public API with comprehensive documentation and examples
- **AWS Lambda Template**: Serverless functions with CloudWatch-optimized logging

#### 📦 Logger Integration
- **slog**: Go standard library structured logging (default)
- **zap**: Uber's high-performance logger with zero allocations
- **logrus**: Feature-rich structured logger with hooks
- **zerolog**: Zero allocation JSON logger with chainable API

#### 🔧 Development Tools
- **Comprehensive CLI**: Full command-line interface with help and validation
- **Blueprint Validation**: All blueprints compile and run successfully
- **Integration Testing**: End-to-end testing of all blueprint+logger combinations
- **Configuration Management**: Support for user profiles and defaults

#### 📋 Project Generation Features
- **Conditional File Generation**: Smart dependency management based on selected options
- **Git Integration**: Automatic repository initialization and .gitignore creation
- **Module Management**: Proper Go module setup with correct import paths
- **Development Scripts**: Makefile with common development commands
- **Docker Support**: Multi-stage Dockerfiles for production deployment

#### 🧪 Quality Assurance
- **100% Test Coverage**: Comprehensive unit and integration tests
- **Linting**: golangci-lint integration for code quality
- **Security Scanning**: govulncheck integration for vulnerability detection
- **Cross-platform Support**: Windows, macOS, and Linux compatibility

#### 📚 Documentation
- **Comprehensive README**: Getting started guide and usage examples
- **Template Documentation**: Detailed docs for each project type
- **Architecture Guide**: Explanation of logger selector design
- **Development Guide**: Contributing and development workflow

### Technical Details

#### Template Architecture
- **Embedded Blueprints**: All blueprints embedded in binary using Go embed
- **Blueprint Registry**: Centralized blueprint management and loading
- **Variable System**: Consistent variable system across all blueprints
- **Conditional Logic**: Smart file generation based on user selections

#### Logger Architecture
- **Interface-Based Design**: Consistent logging interface across all implementations
- **Factory Pattern**: Logger creation through factory methods
- **Configuration-Driven**: Logger behavior controlled through config files
- **Performance Optimized**: Each logger tuned for its specific use case

#### Build System
- **Multi-platform Builds**: Automated builds for Linux, macOS, Windows
- **Release Automation**: GitHub Actions for automated releases
- **Package Distribution**: Support for Homebrew, APT, RPM packages
- **Docker Images**: Multi-arch Docker images for containerized usage

### Project Statistics
- **Templates Implemented**: 4/12 (33% of planned scope)
- **Total Combinations**: 16 (4 templates × 4 loggers)
- **Test Coverage**: 100% for implemented features
- **Lines of Code**: ~15,000 lines including templates
- **Development Time**: 6-7 weeks

### Breaking Changes
- This is the initial release, no breaking changes

### Migration Guide
- This is the first stable release (v1.0.0) of go-starter
- No migration needed from previous versions

### Known Issues
- None at release time

### Future Roadmap
- See [PROJECT_ROADMAP.md](PROJECT_ROADMAP.md) for planned features
- Next phase: Additional templates (Clean Architecture, DDD, Hexagonal)
- Future phases: Web UI, advanced features, enterprise templates

---

## Development Notes

### Version Numbering
- This project follows [Semantic Versioning](https://semver.org/)
- Version format: MAJOR.MINOR.PATCH
- Pre-release versions: MAJOR.MINOR.PATCH-alpha.N, MAJOR.MINOR.PATCH-beta.N, MAJOR.MINOR.PATCH-rc.N

### Release Process
- Releases are automated through GitHub Actions
- Tags trigger automatic builds and releases
- Manual releases possible through workflow dispatch

### Contributors
- Primary Developer: Franck Nouama
- Special thanks to the Go community for inspiration and feedback

[Unreleased]: https://github.com/francknouama/go-starter/compare/v2.2.0...HEAD
[2.2.0]: https://github.com/francknouama/go-starter/compare/v1.1.0...v2.2.0
[1.1.0]: https://github.com/francknouama/go-starter/compare/v1.0.0...v1.1.0
[1.0.0]: https://github.com/francknouama/go-starter/releases/tag/v1.0.0