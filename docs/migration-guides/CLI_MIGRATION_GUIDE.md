# CLI Migration Guide: From Simple to Standard

## Overview

This guide helps you migrate from CLI-Simple to CLI-Standard when your project outgrows its initial simple structure. The migration path is designed to be incremental and non-disruptive.

## When to Migrate

### Migration Indicators

Consider migrating from CLI-Simple to CLI-Standard when you encounter these scenarios:

#### **Immediate Migration Triggers**
- ✅ **Need for 5+ commands**: CLI-Simple is optimal for 2-3 commands
- ✅ **Configuration file requirements**: Need for YAML/JSON configuration
- ✅ **Multiple environments**: Dev, staging, production configurations
- ✅ **Team collaboration**: Multiple developers requiring standards
- ✅ **Distribution needs**: Packaging for external users

#### **Growth Indicators**
- ✅ **Complex business logic**: Logic outgrows simple structure
- ✅ **Advanced logging**: Need for structured logging with multiple outputs
- ✅ **Interactive features**: Survey-style prompts and wizards
- ✅ **CI/CD requirements**: Automated testing and deployment
- ✅ **Container deployment**: Docker and Kubernetes integration

#### **Complexity Thresholds**
- **Files**: When you need more than 10 files
- **Commands**: When you have more than 3 subcommands  
- **Flags**: When you need more than 5 persistent flags
- **Dependencies**: When you need specialized libraries
- **Team Size**: When more than 2 developers work on the CLI

## Migration Strategy

### Phase 1: Assessment and Planning

#### 1.1 Current State Analysis
```bash
# Analyze your current CLI-Simple project
find . -name "*.go" | wc -l          # Count Go files
grep -r "cobra.Command" . | wc -l    # Count commands
grep -r "Flag" cmd/ | wc -l          # Count flags
```

#### 1.2 Feature Gap Analysis
Compare what you have vs. what CLI-Standard offers:

| Feature | CLI-Simple | CLI-Standard | Migration Benefit |
|---------|------------|--------------|------------------|
| Commands | 2-3 basic | 8+ with groups | Better organization |
| Configuration | Environment only | YAML + Viper | Complex config support |
| Logging | slog only | 4 logger options | Production logging |
| Testing | Basic | Comprehensive | Better test coverage |
| CI/CD | None | GitHub Actions | Automated workflows |
| Docker | None | Multi-stage build | Container deployment |

#### 1.3 Migration Timeline
- **Small projects**: 1-2 days
- **Medium projects**: 3-5 days  
- **Large projects**: 1-2 weeks

### Phase 2: Preparation

#### 2.1 Backup Current Project
```bash
# Create a backup before migration
cp -r my-cli-project my-cli-project-backup
git tag pre-migration-backup
```

#### 2.2 Generate CLI-Standard Reference
```bash
# Generate a CLI-Standard project for reference
go-starter new my-cli-reference --type cli --complexity standard \
  --module github.com/yourorg/reference --output ./reference
```

#### 2.3 Dependency Analysis
```bash
# Review current dependencies
go mod graph
go mod why -m github.com/spf13/cobra
```

### Phase 3: Incremental Migration

#### 3.1 Project Structure Migration

**Step 1: Create internal packages**
```bash
mkdir -p internal/config
mkdir -p internal/logger  
mkdir -p internal/errors
mkdir -p internal/version
```

**Step 2: Migrate configuration**
```go
// Before (config.go)
func LoadConfig() *Config {
    return &Config{
        AppName: "my-cli",
        Debug:   getEnvBool("MY_CLI_DEBUG", false),
    }
}

// After (internal/config/config.go) 
func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./configs")
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, fmt.Errorf("failed to read config: %w", err)
    }
    
    var cfg Config
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, fmt.Errorf("failed to unmarshal config: %w", err)
    }
    
    return &cfg, nil
}
```

**Step 3: Migrate logging system**
```go
// Before (main.go)
logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: slog.LevelInfo,
}))

// After (internal/logger/factory.go)
loggerFactory := logger.NewFactory()
appLogger, err := loggerFactory.CreateFromProjectConfig(
    "slog",
    cfg.Logging.Level,
    cfg.Logging.Format,
    cfg.Logging.Structured,
)
```

#### 3.2 Command Structure Migration

**Step 1: Enhance root command**
```go
// Before (cmd/root.go) - Simple
var rootCmd = &cobra.Command{
    Use:   "my-cli",
    Short: "A simple CLI",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Welcome!")
    },
}

// After (cmd/root.go) - Standard
var rootCmd = &cobra.Command{
    Use:   "my-cli",
    Short: "A professional CLI tool",
    Long: `Professional CLI application with:
- Configuration management
- Structured logging
- Multiple output formats
- Interactive mode`,
    RunE: func(cmd *cobra.Command, args []string) error {
        return runRootCommand(cmd, args)
    },
}
```

**Step 2: Add new commands incrementally**
```bash
# Copy command templates from CLI-Standard reference
cp reference/cmd/create.go cmd/
cp reference/cmd/list.go cmd/
cp reference/cmd/delete.go cmd/
cp reference/cmd/update.go cmd/
```

**Step 3: Update command registration**
```go
// Add to cmd/root.go init()
func init() {
    rootCmd.AddCommand(createCmd)
    rootCmd.AddCommand(listCmd)
    rootCmd.AddCommand(deleteCmd)
    rootCmd.AddCommand(updateCmd)
    rootCmd.AddCommand(completionCmd)
}
```

#### 3.3 Configuration System Migration

**Step 1: Add configuration files**
```bash
mkdir configs
```

```yaml
# configs/config.yaml
app:
  name: "my-cli"
  version: "1.0.0"

logging:
  level: "info"
  format: "text"
  structured: true

output:
  format: "table"
  quiet: false
  no_color: false
```

**Step 2: Update main.go**
```go
// Before (simple main.go)
func main() {
    if err := cmd.Execute(); err != nil {
        os.Exit(1)
    }
}

// After (standard main.go)
func main() {
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
        os.Exit(1)
    }

    loggerFactory := logger.NewFactory()
    appLogger, err := loggerFactory.CreateFromProjectConfig(
        "slog",
        cfg.Logging.Level,
        cfg.Logging.Format,
        cfg.Logging.Structured,
    )
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
        os.Exit(1)
    }

    if err := cmd.Execute(appLogger); err != nil {
        appLogger.ErrorWith("Command execution failed", logger.Fields{"error": err})
        os.Exit(1)
    }
}
```

#### 3.4 Dependency Migration

**Step 1: Update go.mod**
```bash
go get github.com/spf13/viper@v1.16.0
go get github.com/AlecAivazis/survey/v2@v2.3.7
go get gopkg.in/yaml.v3@v3.0.1
go get github.com/stretchr/testify@v1.8.4
```

**Step 2: Update imports**
```go
// Add new imports to relevant files
import (
    "github.com/spf13/viper"
    "github.com/AlecAivazis/survey/v2"
    "gopkg.in/yaml.v3"
)
```

### Phase 4: Testing and Validation

#### 4.1 Build Verification
```bash
# Ensure project builds successfully
go build -o bin/my-cli .
go test ./...
```

#### 4.2 Feature Validation
```bash
# Test all commands work
./bin/my-cli --help
./bin/my-cli version
./bin/my-cli create --help
./bin/my-cli list --help
```

#### 4.3 Configuration Testing
```bash
# Test configuration loading
./bin/my-cli --config configs/config.yaml --verbose
./bin/my-cli --output json
```

### Phase 5: Enhancement

#### 5.1 Add CI/CD Pipeline
```bash
mkdir -p .github/workflows
cp reference/.github/workflows/ci.yml .github/workflows/
cp reference/.github/workflows/release.yml .github/workflows/
```

#### 5.2 Add Docker Support
```bash
cp reference/Dockerfile .
cp reference/docker-compose.yml .
```

#### 5.3 Enhanced Documentation
```bash
cp reference/README.md ./README-standard.md
# Merge with your existing README.md
```

## Troubleshooting Common Issues

### Issue 1: Configuration Not Found
```bash
# Error: Config File "my-cli" Not Found
# Solution: Ensure configs/config.yaml exists
mkdir -p configs
echo "app:\n  name: my-cli" > configs/config.yaml
```

### Issue 2: Logger Factory Errors
```bash
# Error: Failed to create logger
# Solution: Check logger configuration
# Ensure cfg.Logging is properly initialized
```

### Issue 3: Command Conflicts
```bash
# Error: Command already exists
# Solution: Remove duplicate command registrations
# Check cmd/root.go init() function
```

### Issue 4: Import Cycle
```bash
# Error: import cycle not allowed
# Solution: Move shared types to separate package
mkdir internal/types
```

## Migration Checklist

### Pre-Migration
- [ ] Project backup created
- [ ] CLI-Standard reference generated
- [ ] Dependencies analyzed
- [ ] Migration timeline planned

### Structure Migration
- [ ] internal/ packages created
- [ ] Configuration migrated to internal/config/
- [ ] Logging migrated to internal/logger/
- [ ] Version info migrated to internal/version/

### Feature Migration  
- [ ] Commands enhanced with error handling
- [ ] Configuration file support added
- [ ] Interactive prompts implemented
- [ ] Output formatting enhanced

### Testing
- [ ] All commands build successfully
- [ ] Configuration loading works
- [ ] Logging outputs correctly
- [ ] Interactive features functional

### Enhancement
- [ ] CI/CD pipeline added
- [ ] Docker support implemented
- [ ] Documentation updated
- [ ] Tests passing

## Post-Migration Benefits

### Immediate Benefits
- ✅ **Professional CLI structure**
- ✅ **Configuration file support**
- ✅ **Structured logging**
- ✅ **Better error handling**
- ✅ **Interactive features**

### Long-term Benefits
- ✅ **Team collaboration support**
- ✅ **CI/CD integration**
- ✅ **Container deployment**
- ✅ **Extensible architecture**
- ✅ **Production readiness**

## When NOT to Migrate

Sometimes CLI-Simple is the better choice:

### Stay with CLI-Simple if:
- ✅ **Single purpose tool**: Does one thing well
- ✅ **Personal use**: Only you use the CLI
- ✅ **Stable requirements**: No new features planned
- ✅ **Performance critical**: Startup time matters
- ✅ **Simplicity preferred**: KISS principle applies

### Alternative: CLI-Medium
Consider a custom middle-ground approach:
- **15-18 files** (between Simple's 8 and Standard's 30)
- **Basic configuration** support
- **Single logger** (slog only)
- **Essential middleware**
- **No Docker/CI by default**

## Support and Resources

### Documentation
- [CLI Blueprint Comparison](./CLI_FEATURE_COMPARISON.md)
- [Progressive Disclosure Guide](./PROGRESSIVE_DISCLOSURE.md)
- [Best Practices](./CLI_BEST_PRACTICES.md)

### Community
- GitHub Issues: Report migration problems
- Discussions: Share migration experiences
- Examples: See successful migration examples

### Tools
- **Migration Script**: Automated migration assistance (planned)
- **Validation Tool**: Check migration completeness
- **Diff Tool**: Compare Simple vs Standard structures

---

**Remember**: Migration should be driven by actual needs, not perceived complexity. CLI-Simple serves 80% of use cases effectively. Migrate only when you've outgrown its capabilities.