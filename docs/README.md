# go-starter Documentation

Welcome to the comprehensive go-starter documentation. Find everything you need to master Go project generation.

## 📁 Documentation Structure

```
docs/
├── guides/           # User and developer guides
├── references/       # API and command references  
├── project-plans/    # Project roadmaps and plans
├── migration-guides/ # Migration and transition guides
└── audits/          # Blueprint audits and reports
```

## 🚀 Getting Started

**New to go-starter?** Start here:

- 📖 **[Getting Started Guide](guides/GETTING_STARTED.md)** - Your first project in 5 minutes
- ⚙️ **[Installation Guide](guides/INSTALLATION.md)** - All installation methods
- 🏃‍♂️ **[Quick Reference](references/QUICK_REFERENCE_CARD.md)** - Common commands and patterns

## 📚 User Guides

### Essential Guides
- 🏗️ **[Development Guide](guides/DEVELOPMENT.md)** - Setting up your development environment
- ⚙️ **[Configuration Guide](guides/CONFIGURATION.md)** - Global settings and profiles
- 🧪 **[Testing Guide](guides/TESTING_GUIDE.md)** - Test your generated projects
- 🔧 **[Troubleshooting Guide](guides/TROUBLESHOOTING.md)** - Solve common issues
- ❓ **[FAQ](guides/FAQ.md)** - Frequently asked questions

## 📖 References

### Project Creation
- 📖 **[Project Types Guide](references/PROJECT_TYPES.md)** - Choose the right template (Web API, CLI, Library, Lambda)
- 🏗️ **[Blueprint Guide](references/BLUEPRINTS.md)** - Deep dive into project templates
- 📊 **[Blueprint Comparison](references/BLUEPRINT_COMPARISON.md)** - Side-by-side feature comparison

### Configuration & Customization
- 📊 **[Logger Guide](references/LOGGER_GUIDE.md)** - Master the unique logger selector system
- 🗃️ **[ORM Guide](references/ORM_GUIDE.md)** - Database and ORM selection
- 📝 **[Quick Reference](references/QUICK_REFERENCE.md)** - Command cheatsheet

## 🗺️ Project Plans

- 🎯 **[Phase 2 Completion Plan](project-plans/PHASE_2_COMPLETION_PLAN.md)** - Current development status (95% complete)
- 🌐 **[Phase 3 Web UI Plan](project-plans/PHASE_3_WEB_UI_DEVELOPMENT_PLAN.md)** - Web interface development roadmap  
- 🚀 **[CI/CD Infrastructure Plan](project-plans/CI_CD_INFRASTRUCTURE_IMPROVEMENT_PLAN.md)** - Deployment automation plans

## 🔄 Migration Guides

- 🛠️ **[CLI Migration Guide](migration-guides/CLI_MIGRATION_GUIDE.md)** - Upgrading CLI blueprints
- ⚡ **[CLI Over-Engineering Resolution](migration-guides/CLI_OVER_ENGINEERING_COMPLETE_RESOLUTION.md)** - Simplification strategies
- 📝 **[TDD Enforcement](migration-guides/TDD-ENFORCEMENT.md)** - Test-driven development practices
- 📋 **[Task Master Guide](migration-guides/TASK_MASTER_GUIDE.md)** - AI-powered task management

## 🔍 Audit Reports

Comprehensive blueprint quality assessments and improvement reports are available in the [audits/](audits/) directory.

## 🆘 Help & Support

- 🐛 **[Report Issues](https://github.com/francknouama/go-starter/issues)** - Found a bug?
- 💬 **[Discussions](https://github.com/francknouama/go-starter/discussions)** - Community support

## 🚀 Quick Start

```bash
# Install go-starter (using Go install - recommended)
go install github.com/francknouama/go-starter@latest

# Alternative: Download binary from GitHub releases
# Homebrew currently unavailable due to PAT issues

# Generate your first project
go-starter new my-awesome-api

# Navigate and run
cd my-awesome-api
make run
```

## 📊 Available Templates

| Template | Use Case | Key Features |
|----------|----------|--------------|
| **Web API** | REST services, microservices | Gin framework, middleware, database |
| **CLI** | Command-line tools | Cobra framework, subcommands |
| **Library** | Reusable packages | Clean API, examples, minimal deps |
| **Lambda** | Serverless functions | AWS integration, CloudWatch logging |

## 🪵 Logger Options

| Logger | Performance | Best For |
|--------|-------------|----------|
| **slog** | Good | Standard library choice |
| **zap** | Excellent | High-performance applications |
| **logrus** | Good | Feature-rich requirements |
| **zerolog** | Excellent | JSON-heavy, cloud-native |

## 🎯 Common Use Cases

### Building a High-Performance API
```bash
go-starter new api --type=web-api --logger=zap
```

### Creating a Developer Tool
```bash
go-starter new tool --type=cli --logger=logrus
```

### Developing a Reusable Library
```bash
go-starter new sdk --type=library --logger=slog
```

### Deploying Serverless Functions
```bash
go-starter new function --type=lambda --logger=zerolog
```

## 📖 Documentation Structure

```
docs/
├── README.md                    # This file
├── GETTING_STARTED.md          # Installation and first steps
├── BLUEPRINT_COMPARISON.md      # Choosing the right blueprint
├── BLUEPRINTS.md                # Detailed blueprint documentation
├── LOGGER_GUIDE.md             # Logger selector deep dive
├── ORM_GUIDE.md                # Database interaction patterns
├── QUICK_REFERENCE_CARD.md     # Commands and patterns reference
├── FAQ.md                      # Frequently asked questions
└── TROUBLESHOOTING.md          # Problem-solving guide
```

## 🤝 Contributing

We welcome contributions! Please see:
- [Contributing Guide](../CONTRIBUTING.md)
- [Project Roadmap](../PROJECT_ROADMAP.md)

## 🔗 Links

- **Repository**: [github.com/francknouama/go-starter](https://github.com/francknouama/go-starter)
- **Issues**: [GitHub Issues](https://github.com/francknouama/go-starter/issues)
- **Releases**: [Latest Releases](https://github.com/francknouama/go-starter/releases)

---

**Need help?** Check the [FAQ](FAQ.md) or [Troubleshooting Guide](TROUBLESHOOTING.md) first. If you can't find an answer, please [open an issue](https://github.com/francknouama/go-starter/issues/new).