# validation-cli

A simple command-line application built with Go and Cobra.

## Features

- **Simple Structure**: Clean and minimal CLI architecture
- **Essential Flags**: `--help`, `--version`, `--quiet`, `--output`
- **Shell Completion**: Support for bash, zsh, fish, and PowerShell
- **Standard Logging**: Uses Go's standard `slog` package
- **Environment Variables**: Simple configuration via environment variables

## Installation

### Build from source

```bash
# Clone the repository
git clone <repository-url>
cd validation-cli

# Build the application
make build

# Or install directly
make install
```

### Using go install

```bash
go install github.com/test/validation-cli@latest
```

## Usage

### Basic Usage

```bash
# Show help
validation-cli --help

# Show version
validation-cli --version

# Run with quiet output
validation-cli --quiet

# JSON output format
validation-cli --output json
```

### Environment Variables

Configure the application using environment variables:

```bash
# Set log level (debug, info, warn, error)
export VALIDATION-CLI_LOG_LEVEL=debug

# Enable debug mode
export VALIDATION-CLI_DEBUG=true

# Set default output format
export VALIDATION-CLI_FORMAT=json

# Enable quiet mode by default
export VALIDATION-CLI_QUIET=true
```

### Shell Completion

Generate shell completion scripts:

```bash
# Bash
validation-cli completion bash > /etc/bash_completion.d/validation-cli

# Zsh
validation-cli completion zsh > "${fpath[1]}/_validation-cli"

# Fish
validation-cli completion fish > ~/.config/fish/completions/validation-cli.fish

# PowerShell
validation-cli completion powershell > validation-cli.ps1
```

## Development

### Prerequisites

- Go 1.23 or later
- Make (optional, for using Makefile)

### Building

```bash
# Build the application
make build

# Run tests
make test

# Format code
make fmt

# Run all checks
make check

# Build for multiple platforms
make build-all
```

### Project Structure

```
validation-cli/
├── main.go              # Application entry point
├── config.go            # Configuration management
├── cmd/
│   ├── root.go          # Root command
│   └── version.go       # Version command
├── go.mod               # Go module file
├── Makefile             # Build automation
└── README.md            # This file
```

## Architecture

This project follows a simple CLI architecture:

- **Minimal Dependencies**: Only uses Cobra for CLI framework and standard library
- **Simple Configuration**: Environment variable based configuration
- **Standard Logging**: Uses Go's built-in `slog` package
- **Essential Features**: Focus on core CLI functionality without unnecessary complexity

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Author

Created by the validation-cli team