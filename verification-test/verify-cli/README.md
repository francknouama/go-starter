# verify-cli

A command-line application built with Go and Cobra.

## Features

- 🚀 Built with [Cobra](https://cobra.dev/) CLI framework
- 📝 Structured logging with slog logger
- ⚙️ Configuration management with [Viper](https://github.com/spf13/viper)
- 🧪 Testing setup with [Testify](https://github.com/stretchr/testify)
- 🐳 Docker support
- 📦 Multi-platform builds

## Installation

### From Source

```bash
git clone github.com/verify/cli
cd verify-cli
make install
```

### Using Go

```bash
go install github.com/verify/cli@latest
```

## Usage

```bash
# Show help
verify-cli --help

# Show version
verify-cli version

# Run with verbose output
verify-cli --verbose

# Use custom config file
verify-cli --config ./my-config.yaml
```

## Configuration

verify-cli can be configured via:

1. **Configuration file** (YAML format)
2. **Environment variables** (prefixed with `VERIFY-CLI_`)
3. **Command-line flags**

### Configuration File

Create a config file at one of these locations:
- `./configs/config.yaml`
- `./config.yaml`
- `$HOME/.verify-cli.yaml`

Example configuration:

```yaml
environment: development

logging:
  level: info
  format: text
  structured: false

cli:
  output_format: text
  no_color: false
  quiet: false
```

### Environment Variables

```bash
export VERIFY-CLI_LOGGING_LEVEL=debug
export VERIFY-CLI_CLI_OUTPUT_FORMAT=json
```

## Development

### Prerequisites

- Go 1.21 or later
- Make (optional, for using Makefile commands)

### Building

```bash
# Build binary
make build

# Run tests
make test

# Run linter
make lint

# Run with coverage
make test-coverage
```

### Docker

```bash
# Build Docker image
make docker-build

# Run in Docker
make docker-run
```

## Project Structure

```
verify-cli/
├── cmd/                    # Command implementations
│   ├── root.go            # Root command
│   └── version.go         # Version command
├── internal/
│   ├── config/            # Configuration management
│   └── logger/            # Logging implementation
├── configs/               # Configuration files
├── Dockerfile            # Docker build file
├── Makefile              # Build automation
└── main.go               # Application entry point
```

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the  License - see the LICENSE file for details.

## Generated with

This project was generated using [go-starter](https://github.com/francknouama/go-starter) with the following configuration:

- **Template**: CLI Application (cobra)
- **Logger**: slog
- **Go Version**: 1.21
- **License**: 