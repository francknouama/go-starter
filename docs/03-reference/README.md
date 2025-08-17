# Reference Documentation

Detailed technical documentation for go-starter. Complete specifications, API references, and configuration options.

## 📋 Available References

### 🪵 [Logger Comparison Guide](LOGGER_COMPARISON_GUIDE.md)
**Deep dive into logger selection and optimization**
- Performance benchmarks and comparisons
- Configuration options for each logger
- Best practices and optimization techniques
- Migration strategies between loggers

### 🏗️ [Blueprint Reference](../references/BLUEPRINTS.md)
**Complete specifications for all blueprints**
- Detailed file structures for all 12 blueprints
- Template variables and conditional logic
- Dependencies and requirements for each blueprint
- Generated code examples and patterns

### 📊 [Blueprint Comparison](../references/BLUEPRINT_COMPARISON.md)
**Feature and architecture comparison**
- Side-by-side blueprint comparisons
- Architecture pattern analysis
- Use case recommendations

### 📝 [Quick Reference](../references/QUICK_REFERENCE.md)
**Commands and options cheatsheet**
- Essential commands at a glance
- Common usage patterns
- Quick troubleshooting tips

## 🎯 Quick Navigation

### Looking for specific information?

| Need | Go To |
|------|-------|
| **Project structure details** | [Blueprint Reference](blueprints.md) |
| **Command syntax** | [CLI Reference](cli-reference.md) |
| **Logger performance** | [Logger Guide](logger-guide.md) |
| **Config file options** | [Configuration Reference](configuration-reference.md) |
| **Troubleshooting** | [User Guides](../02-user-guides/troubleshooting.md) |

### By Experience Level

**New Users**:
- Start with [User Guides](../02-user-guides/) for practical guidance
- Use reference docs for specific questions

**Experienced Users**:
- [CLI Reference](cli-reference.md) for command mastery
- [Configuration Reference](configuration-reference.md) for advanced setups

**Contributors**:
- [Developer Docs](../04-developers/) for internal architecture
- [Blueprint Reference](blueprints.md) for template specifications

## 📊 Reference Quick Facts

### Blueprint Statistics
- **12 Total Blueprints**: Covering all Go project types
- **4 Architecture Patterns**: Standard, Clean, DDD, Hexagonal
- **File Range**: 8 files (CLI Simple) to 60+ files (Enterprise)
- **Production Ready**: All blueprints tested and deployment-ready

### CLI Capabilities
- **2 Interface Modes**: Basic (14 flags) and Advanced (18+ flags)
- **Progressive Disclosure**: Adapts to user experience level
- **Cross-Platform**: Windows, macOS, Linux support
- **Multiple Install Methods**: Go install, Homebrew, package managers

### Logger Performance
- **4 Logger Options**: slog, zap, logrus, zerolog
- **Unified Interface**: Same API across all implementations
- **Performance Range**: Good (slog/logrus) to Excellent (zap/zerolog)
- **Zero Dependencies**: slog has no external dependencies

### Configuration Features
- **Profile System**: Multiple configuration contexts
- **Environment Variables**: Override any configuration option
- **Team Standards**: Shared configuration for consistency
- **Validation**: Built-in validation with helpful error messages

---

**Need practical guidance?** Check our [User Guides](../02-user-guides/) for step-by-step instructions and examples.

**Want to contribute?** See [Developer Documentation](../04-developers/) for internal architecture and contribution guidelines.