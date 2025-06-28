# verify-zap

verify-zap is a modern Go web API built with gin framework.

## Features

- 🚀 **Fast & Lightweight**: Built with gin framework for high performance
- 🔒 **Secure**: Ready for authentication integration
- 📝 **API Documentation**: OpenAPI/Swagger specification
- 🐳 **Docker Ready**: Multi-stage Docker builds
- 🧪 **Testing**: Comprehensive test suite
- 📊 **Monitoring**: Health checks and logging
- 🔧 **Development**: Hot reload and development tools

## Quick Start

### Prerequisites

- Go 1.21 or later
- Docker (optional)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd verify-zap
```

2. Install dependencies:
```bash
go mod download
```
3. Set up environment variables:
```bash
cp .env.example .env
# Edit .env with your configuration
```

### Development

Start the development server:
```bash
make dev
```

The API will be available at `http://localhost:8080`.

### Building

Build the application:
```bash
make build
```

Run the built binary:
```bash
make run
```

## API Documentation

The API documentation is available at:
- OpenAPI spec: `/api/openapi.yaml`
- When running: `http://localhost:8080/api/openapi.yaml`

## Available Endpoints

### Health Checks
- `GET /health` - Basic health check
- `GET /ready` - Readiness check

## Configuration

The application can be configured using environment variables or YAML files.

### Environment Variables
- `SERVER_PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment (development, production)

### Configuration Files

Configuration files are located in the `configs/` directory:
- `config.dev.yaml` - Development configuration
- `config.prod.yaml` - Production configuration
- `config.test.yaml` - Test configuration

## Development Commands

```bash
# Development
make dev          # Start development server with hot reload
make build        # Build the application
make run          # Run the built application

# Testing
make test         # Run tests with coverage
make lint         # Run linter

# Database (if enabled)

# Docker
make docker-build # Build Docker image
make docker-run   # Run with Docker

# Utilities
make clean        # Clean build artifacts
make fmt          # Format code
make tidy         # Tidy dependencies
make help         # Show available commands
```

## Project Structure

```
verify-zap/
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # HTTP middleware
├── api/                 # API documentation
├── configs/             # Configuration files
├── tests/               # Test files
├── scripts/             # Development scripts
├── Dockerfile           # Docker configuration
└── Makefile            # Development commands
```

## Testing

Run the test suite:
```bash
make test
```

This will run all tests and generate a coverage report.

## Docker

### Build and run with Docker:
```bash
make docker-run
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the  License - see the LICENSE file for details.

## Support

If you have any questions or need help, please:
- Check the [API documentation](api/openapi.yaml)
- Review the configuration files in `configs/`
- Look at the example requests in `tests/`

---

Generated with [go-starter](https://github.com/your-org/go-starter) 🚀