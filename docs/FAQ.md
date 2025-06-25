# Troubleshooting & FAQ

Common issues, solutions, and frequently asked questions for go-starter.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Project Generation Issues](#project-generation-issues)
- [Compilation Issues](#compilation-issues)
- [Logger Issues](#logger-issues)
- [Template Issues](#template-issues)
- [Performance Issues](#performance-issues)
- [Frequently Asked Questions](#frequently-asked-questions)

---

## Installation Issues

### Q: `go install` fails with permission error

**Problem**: 
```bash
go install github.com/francknouama/go-starter@latest
# Error: permission denied
```

**Solutions**:

1. **Check GOPATH and GOBIN**:
   ```bash
   echo $GOPATH
   echo $GOBIN
   mkdir -p $GOPATH/bin
   ```

2. **Use proper Go version**:
   ```bash
   go version  # Should be Go 1.19+
   ```

3. **Install to specific directory**:
   ```bash
   GOBIN=/usr/local/bin go install github.com/francknouama/go-starter@latest
   ```

4. **Use sudo if necessary** (macOS/Linux):
   ```bash
   sudo go install github.com/francknouama/go-starter@latest
   ```

### Q: Binary not found in PATH

**Problem**: 
```bash
go-starter: command not found
```

**Solutions**:

1. **Add GOBIN to PATH**:
   ```bash
   # Add to ~/.bashrc, ~/.zshrc, or ~/.profile
   export PATH=$PATH:$(go env GOPATH)/bin
   source ~/.bashrc  # or ~/.zshrc
   ```

2. **Verify installation location**:
   ```bash
   ls $(go env GOPATH)/bin/go-starter
   ```

3. **Use full path temporarily**:
   ```bash
   $(go env GOPATH)/bin/go-starter new my-project
   ```

### Q: Download fails with network error

**Problem**: 
```bash
go install: module github.com/francknouama/go-starter@latest: 
Get "https://proxy.golang.org/...": dial tcp: lookup proxy.golang.org: no such host
```

**Solutions**:

1. **Configure Go proxy**:
   ```bash
   go env -w GOPROXY=https://proxy.golang.org,direct
   go env -w GOSUMDB=sum.golang.org
   ```

2. **Use direct download**:
   ```bash
   go env -w GOPROXY=direct
   go install github.com/francknouama/go-starter@latest
   ```

3. **Corporate proxy setup**:
   ```bash
   export GOPROXY=https://your-corporate-proxy.com
   export GOSUMDB=off
   ```

---

## Project Generation Issues

### Q: Generation fails with "template not found" 

**Problem**:
```bash
go-starter new my-project --type=web-api
# Error: template 'web-api' not found
```

**Solutions**:

1. **Check available templates**:
   ```bash
   go-starter list  # Shows all available templates
   ```

2. **Use correct template name**:
   ```bash
   go-starter new my-project --type=web-api    # Correct
   go-starter new my-project --type=api        # Incorrect
   ```

3. **Update to latest version**:
   ```bash
   go install github.com/francknouama/go-starter@latest
   ```

### Q: Generation fails with "invalid module path"

**Problem**:
```bash
go-starter new my-project --module=invalid-path
# Error: invalid module path
```

**Solutions**:

1. **Use valid Go module path**:
   ```bash
   # ✅ Good
   go-starter new my-project --module=github.com/username/my-project
   go-starter new my-project --module=example.com/my-project
   
   # ❌ Bad
   go-starter new my-project --module=my-project
   go-starter new my-project --module=../my-project
   ```

2. **Interactive mode for guidance**:
   ```bash
   go-starter new my-project
   # Follow prompts for module path
   ```

### Q: Generation partially completes then fails

**Problem**: Some files are created but generation fails mid-way.

**Solutions**:

1. **Check disk space**:
   ```bash
   df -h .  # Check available space
   ```

2. **Check permissions**:
   ```bash
   ls -la .  # Check directory permissions
   mkdir test-dir  # Test write permissions
   ```

3. **Clean and retry**:
   ```bash
   rm -rf my-project  # Remove partial generation
   go-starter new my-project --type=web-api
   ```

4. **Use different directory**:
   ```bash
   cd /tmp
   go-starter new my-project --type=web-api
   ```

---

## Compilation Issues

### Q: Generated project doesn't compile

**Problem**:
```bash
cd my-project
go build
# Error: cannot find module providing package
```

**Solutions**:

1. **Initialize and download dependencies**:
   ```bash
   go mod tidy
   go mod download
   ```

2. **Check Go version**:
   ```bash
   go version  # Should match go.mod requirements
   ```

3. **Clear module cache**:
   ```bash
   go clean -modcache
   go mod download
   ```

4. **Check network connectivity**:
   ```bash
   go env GOPROXY
   curl -I https://proxy.golang.org
   ```

### Q: Logger-related compilation errors

**Problem**:
```bash
go build
# Error: undefined: zap.Logger
```

**Solutions**:

1. **Verify logger dependencies**:
   ```bash
   go mod tidy
   cat go.mod  # Check if zap is listed
   ```

2. **Reinstall dependencies**:
   ```bash
   go mod download
   go get go.uber.org/zap  # If using zap
   ```

3. **Check logger factory**:
   ```bash
   cat internal/logger/factory.go  # Verify correct implementation
   ```

### Q: Test compilation fails

**Problem**:
```bash
go test ./...
# Error: cannot find package in any of:
```

**Solutions**:

1. **Update test dependencies**:
   ```bash
   go mod tidy
   go get github.com/stretchr/testify/assert
   ```

2. **Check test imports**:
   ```bash
   grep -r "import" . --include="*_test.go"
   ```

3. **Run specific test**:
   ```bash
   go test -v ./internal/handlers/
   ```

---

## Logger Issues

### Q: No log output visible

**Problem**: Application runs but no logs appear.

**Solutions**:

1. **Check log level**:
   ```bash
   LOG_LEVEL=debug go run cmd/server/main.go
   ```

2. **Check log output destination**:
   ```bash
   LOG_OUTPUT=stdout go run cmd/server/main.go
   ```

3. **Verify logger initialization**:
   ```go
   logger := logger.New()
   logger.Info("Test message")  // Should appear
   ```

4. **Check configuration**:
   ```bash
   cat configs/config.yaml  # Verify logger config
   ```

### Q: Log format is incorrect

**Problem**: Logs appear in wrong format (e.g., text instead of JSON).

**Solutions**:

1. **Set format explicitly**:
   ```bash
   LOG_FORMAT=json go run cmd/server/main.go
   ```

2. **Check configuration loading**:
   ```go
   // In main.go, add debug output
   fmt.Printf("Config: %+v\n", config.Logger)
   ```

3. **Verify logger implementation**:
   ```bash
   cat internal/logger/factory.go  # Check format settings
   ```

### Q: Logger performance issues

**Problem**: Application slows down significantly with logging.

**Solutions**:

1. **Use appropriate log level**:
   ```bash
   LOG_LEVEL=info go run cmd/server/main.go  # Avoid debug in production
   ```

2. **Switch to high-performance logger**:
   ```bash
   go-starter new my-project --logger=zap  # Fastest option
   ```

3. **Use conditional logging**:
   ```go
   if logger.IsDebugEnabled() {
       logger.Debug("Expensive operation", "data", formatData(data))
   }
   ```

---

## Template Issues

### Q: Template customization not working

**Problem**: Changes to templates don't appear in generated projects.

**Solutions**:

1. **Templates are embedded**: go-starter uses embedded templates, so you can't modify them directly.

2. **Generate fresh project**:
   ```bash
   go-starter new my-project --type=web-api
   # Then customize the generated project
   ```

3. **Fork and modify** (advanced):
   ```bash
   git clone https://github.com/francknouama/go-starter.git
   # Modify templates/ directory
   # Build your own version
   ```

### Q: Missing files in generated project

**Problem**: Expected files are missing from generated project.

**Solutions**:

1. **Check template completeness**:
   ```bash
   go-starter new test-project --type=web-api
   find test-project -type f  # List all generated files
   ```

2. **Compare with documentation**:
   ```bash
   # Check docs/TEMPLATES.md for expected structure
   ```

3. **Verify flags**:
   ```bash
   go-starter new my-project --type=web-api --logger=zap --verbose
   ```

---

## Performance Issues

### Q: Project generation is slow

**Problem**: `go-starter new` takes a long time to complete.

**Solutions**:

1. **Check disk I/O**:
   ```bash
   time go-starter new test-project --type=web-api
   ```

2. **Use faster disk** (SSD vs HDD):
   ```bash
   cd /tmp  # Often faster than home directory
   go-starter new test-project --type=web-api
   ```

3. **Update to latest version**:
   ```bash
   go install github.com/francknouama/go-starter@latest
   ```

### Q: Generated project has slow build times

**Problem**: `go build` takes too long in generated project.

**Solutions**:

1. **Enable build cache**:
   ```bash
   go env GOCACHE
   go clean -cache  # If cache is corrupted
   ```

2. **Use Go modules properly**:
   ```bash
   go mod tidy
   go mod download
   ```

3. **Parallel builds**:
   ```bash
   go build -p 4  # Use 4 parallel builds
   ```

---

## Frequently Asked Questions

### General Questions

#### Q: What's the difference between go-starter and other Go generators?

**A**: go-starter uniquely combines:
- **Logger selector system**: Choose from 4 logging libraries with consistent interface
- **Production-ready templates**: All templates are fully tested and production-ready
- **Best practices built-in**: Following Go community standards and patterns
- **Zero vendor lock-in**: Switch loggers without changing code

#### Q: Which project type should I choose?

**A**: Choose based on your use case:
- **Web API**: REST services, microservices, HTTP APIs
- **CLI**: Command-line tools, utilities, automation scripts
- **Library**: Reusable packages, SDKs, shared functionality
- **Lambda**: AWS serverless functions, event-driven processing

#### Q: Which logger should I choose?

**A**: Choose based on your requirements:
- **slog**: General purpose, no dependencies, Go 1.21+
- **zap**: High performance, zero allocation, production apps
- **logrus**: Feature-rich, large ecosystem, enterprise apps
- **zerolog**: Cloud-native, zero allocation, clean API

#### Q: Can I switch loggers after project generation?

**A**: The logger interface is consistent, but the implementation files are different. You can:
1. Generate a new project with the desired logger
2. Copy your business logic to the new project
3. Manually replace the logger implementation (advanced)

#### Q: Is go-starter production-ready?

**A**: Yes! All templates are:
- ✅ Fully tested with comprehensive test suites
- ✅ Follow Go best practices and idioms
- ✅ Include production configurations
- ✅ Have proper error handling and logging
- ✅ Include Docker and CI/CD configurations

### Technical Questions

#### Q: What Go version is required?

**A**: 
- **Minimum**: Go 1.19 (for go-starter tool)
- **Generated projects**: Go 1.19+ (Go 1.21+ for slog)
- **Recommended**: Go 1.21+ for best experience

#### Q: Can I use go-starter in corporate environments?

**A**: Yes! go-starter is designed for professional use:
- MIT license allows commercial use
- No external dependencies in generated projects (except chosen logger)
- Configurable for corporate proxies and air-gapped environments
- Follows security best practices

#### Q: How do I customize templates?

**A**: Templates are embedded in the binary. For customization:
1. **Post-generation**: Modify the generated project
2. **Fork approach**: Fork the repository and modify templates
3. **Configuration**: Use config files to customize behavior

#### Q: Does go-starter support databases?

**A**: Current templates include:
- ✅ Database configuration structure
- ✅ GORM integration examples
- ✅ Migration setup
- ✅ Connection pooling

Future versions will include:
- Database driver selection
- Multiple ORM options
- Migration tools

#### Q: Can I add more templates?

**A**: Currently, go-starter includes 4 core templates. Future versions will include:
- Clean Architecture patterns
- Domain-Driven Design (DDD)
- Hexagonal Architecture
- Microservice templates
- Event-driven architectures

### Troubleshooting Questions

#### Q: My project won't start after generation

**A**: Check the following:
```bash
# 1. Dependencies installed
go mod tidy

# 2. Configuration valid
cat configs/config.yaml

# 3. Run with debug logging
LOG_LEVEL=debug make run

# 4. Check port availability
netstat -tulpn | grep :8080
```

#### Q: Tests are failing in generated project

**A**: Common solutions:
```bash
# 1. Update test dependencies
go mod tidy

# 2. Run specific test
go test -v ./internal/handlers/

# 3. Check test environment
go test -v -race ./...

# 4. Update assertions
go get github.com/stretchr/testify/assert@latest
```

#### Q: Docker build fails

**A**: Check the following:
```bash
# 1. Verify Dockerfile exists
ls -la Dockerfile

# 2. Check Go version in Dockerfile
grep "FROM golang" Dockerfile

# 3. Build with more verbose output
docker build -t my-app . --no-cache --progress=plain

# 4. Check multi-stage build
docker build --target builder -t my-app-builder .
```

### Best Practices Questions

#### Q: How should I structure my generated project?

**A**: The generated structure follows Go best practices:
- `cmd/`: Application entry points
- `internal/`: Private application code
- `pkg/`: Public library code (if applicable)
- `configs/`: Configuration files
- `tests/`: Test files and test data

#### Q: What's the recommended development workflow?

**A**: Follow this workflow:
1. **Generate project**: `go-starter new my-project`
2. **Install dependencies**: `go mod tidy`
3. **Run tests**: `make test`
4. **Start development**: `make run`
5. **Make changes**: Edit code and re-run tests
6. **Build**: `make build`
7. **Deploy**: Use provided Docker/CI configurations

#### Q: How do I deploy generated projects?

**A**: Multiple deployment options:
```bash
# Docker
make docker
docker run -p 8080:8080 my-app:latest

# Binary
make build
./bin/server --config=configs/config.prod.yaml

# Cloud platforms
# (AWS Lambda, Google Cloud Run, etc.)
```

#### Q: How do I handle secrets and configuration?

**A**: Use environment variables and configuration files:
```yaml
# configs/config.prod.yaml
database:
  host: ${DB_HOST}
  password: ${DB_PASSWORD}
  
# Environment variables
export DB_HOST=prod-db.example.com
export DB_PASSWORD=secure-password
```

---

## Getting Help

### Community Support

- **GitHub Issues**: [Report bugs and request features](https://github.com/francknouama/go-starter/issues)
- **GitHub Discussions**: [Ask questions and share ideas](https://github.com/francknouama/go-starter/discussions)
- **Documentation**: [Read comprehensive guides](https://github.com/francknouama/go-starter/tree/main/docs)

### Reporting Issues

When reporting issues, please include:

1. **go-starter version**: `go-starter version`
2. **Go version**: `go version`
3. **Operating system**: `uname -a` (Linux/Mac) or `ver` (Windows)
4. **Command used**: Exact command that failed
5. **Error output**: Complete error message
6. **Expected behavior**: What you expected to happen

### Contributing

We welcome contributions! See:
- [Contributing Guidelines](https://github.com/francknouama/go-starter/blob/main/CONTRIBUTING.md)
- [Code of Conduct](https://github.com/francknouama/go-starter/blob/main/CODE_OF_CONDUCT.md)
- [Development Guide](https://github.com/francknouama/go-starter/blob/main/docs/DEVELOPMENT.md)

---

*This FAQ is regularly updated. If you don't find your question here, please check our [GitHub Discussions](https://github.com/francknouama/go-starter/discussions) or [create a new issue](https://github.com/francknouama/go-starter/issues/new).*