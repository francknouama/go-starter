# Go-Starter Web Server

A comprehensive backend API for the go-starter web UI, providing REST endpoints for project generation, blueprint management, and real-time WebSocket communication.

## Features

- **RESTful API** - Complete REST API for project generation
- **Blueprint Management** - List, filter, and validate project blueprints
- **Real-time WebSockets** - Live preview updates and generation progress
- **ZIP Downloads** - Generate and download projects as ZIP files
- **Health Monitoring** - Comprehensive health checks and metrics
- **CORS Support** - Configured for React frontend integration
- **Error Handling** - Robust error handling with proper HTTP status codes
- **Static File Serving** - Serves React frontend from embedded assets

## API Endpoints

### Health and System

- `GET /api/v1/health` - Comprehensive health check
- `GET /api/v1/health/simple` - Simple health check for load balancers
- `GET /api/v1/health/readiness` - Readiness probe
- `GET /api/v1/health/liveness` - Liveness probe
- `GET /api/v1/metrics` - System metrics
- `GET /api/v1/status` - Detailed system status
- `GET /api/v1/version` - Version information

### Configuration

- `GET /api/v1/config` - Get default configuration options
- `GET /api/v1/config/types/:type` - Get project type details
- `GET /api/v1/config/frameworks?type=web-api` - Get frameworks for project type
- `GET /api/v1/config/architectures?type=web-api` - Get architectures for project type

### Blueprints

- `GET /api/v1/blueprints` - List all blueprints with filtering
- `GET /api/v1/blueprints/:id` - Get specific blueprint details
- `GET /api/v1/blueprints/category/:category` - Get blueprints by category
- `POST /api/v1/blueprints/:id/validate` - Validate blueprint configuration

### Project Generation

- `POST /api/v1/preview` - Generate project preview (structure only)
- `POST /api/v1/generate` - Generate project with options
- `POST /api/v1/generate/download` - Generate and immediately download

### Downloads

- `GET /api/v1/download/:token` - Download generated ZIP file
- `GET /api/v1/download/:token/status` - Check download status

### WebSocket

- `GET /api/v1/ws` - WebSocket endpoint for real-time communication
- `GET /api/v1/ws/info` - WebSocket endpoint information
- `GET /api/v1/ws/clients` - List connected clients
- `GET /api/v1/ws/stats` - WebSocket hub statistics
- `GET /api/v1/ws/test` - WebSocket test page

## Quick Start

### Prerequisites

- Go 1.21 or later
- Node.js 18 or later (for frontend)
- npm or yarn

### Setup

```bash
# Install dependencies
make setup

# Build frontend and backend
make build-all

# Run development server
make dev
```

### Development

```bash
# Start development server with hot reload
make dev

# Build and run
make run

# Build frontend only
make build-frontend

# Build backend only
make build
```

### Production

```bash
# Production build
make build-all-prod

# Run production server
./bin/go-starter-web
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (debug, test, release) | `release` |
| `CORS_ALLOWED_ORIGINS` | Additional CORS origins (comma-separated) | - |

## Configuration

### CORS Configuration

The server is pre-configured to accept requests from common development origins:

- `http://localhost:3000` (React dev server)
- `http://localhost:5173` (Vite dev server)
- `http://localhost:8080` (Production server)

Additional origins can be added via the `CORS_ALLOWED_ORIGINS` environment variable.

### WebSocket Configuration

WebSocket connections support:

- Real-time preview updates
- Generation progress notifications
- Client connection management
- Automatic reconnection handling

## API Usage Examples

### Generate Project Preview

```bash
curl -X POST http://localhost:8080/api/v1/preview \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-api",
    "modulePath": "github.com/user/my-api",
    "type": "web-api",
    "architecture": "clean",
    "framework": "gin",
    "logger": "slog"
  }'
```

### Generate and Download Project

```bash
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "my-api",
    "modulePath": "github.com/user/my-api",
    "type": "web-api",
    "outputFormat": "zip"
  }' | jq '.data.downloadUrl'
```

### List Available Blueprints

```bash
curl http://localhost:8080/api/v1/blueprints
```

### Filter Blueprints by Category

```bash
curl "http://localhost:8080/api/v1/blueprints?category=web-api&complexity=standard"
```

## WebSocket Usage

### Connect to WebSocket

```javascript
const ws = new WebSocket('ws://localhost:8080/api/v1/ws');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message.type, message.data);
};

// Send ping
ws.send(JSON.stringify({
  type: 'ping',
  data: 'Hello Server'
}));
```

### Subscribe to Events

```javascript
// Subscribe to progress updates
ws.send(JSON.stringify({
  type: 'subscribe',
  channel: 'progress'
}));
```

## Development Tools

### Hot Reloading

```bash
# Install air for hot reloading
go install github.com/cosmtrek/air@latest

# Start with hot reload
make dev
```

### Testing

```bash
# Run tests
make test

# Run with coverage
make test-coverage
```

### Linting

```bash
# Install and run linter
make lint

# Format code
make fmt
```

## Docker Support

### Build Docker Image

```bash
make docker-build
```

### Run Docker Container

```bash
make docker-run
```

### Docker Compose

```yaml
version: '3.8'
services:
  go-starter-web:
    build: .
    ports:
      - "8080:8080"
    environment:
      - GIN_MODE=release
      - CORS_ALLOWED_ORIGINS=https://mydomain.com
```

## Architecture

### Project Structure

```
web/
├── cmd/web-server/          # Main application entry point
├── internal/web/            # Internal web packages
│   ├── handlers/           # HTTP handlers
│   ├── middleware/         # Middleware components
│   ├── models/            # Request/response models
│   └── websocket/         # WebSocket implementation
├── dist/                  # Frontend build output (embedded)
└── bin/                   # Compiled binaries
```

### Key Components

- **Gin Framework** - HTTP web framework
- **WebSocket Support** - Real-time communication
- **Template Integration** - Uses existing CLI generator
- **Embedded Assets** - Serves React frontend
- **Error Handling** - Comprehensive error management
- **CORS Middleware** - Cross-origin request support

## Integration with CLI

The web server reuses the existing CLI generator logic:

- **Templates Registry** - Same blueprint system
- **Generator Engine** - Identical generation logic
- **Validation** - Same validation rules
- **File Processing** - Consistent file handling

## Security Features

- **CORS Protection** - Configurable origin restrictions
- **Security Headers** - CSP, XSS protection, etc.
- **Request Validation** - Input sanitization and validation
- **Error Sanitization** - Safe error message exposure
- **Rate Limiting** - Built-in request limiting (planned)

## Performance

- **Parallel Processing** - Concurrent file generation
- **Caching** - Template parsing cache
- **Streaming** - Efficient file serving
- **Compression** - Gzip response compression (planned)
- **Connection Pooling** - Efficient resource usage

## Monitoring

### Health Checks

- **Liveness** - Server responsiveness
- **Readiness** - Service availability
- **Dependencies** - External service status

### Metrics

- **Request Count** - HTTP request metrics
- **Response Time** - Request latency tracking
- **Error Rate** - Error occurrence monitoring
- **WebSocket Connections** - Real-time connection stats

## Troubleshooting

### Common Issues

1. **Frontend Assets Not Found**
   - Ensure `npm run build` was executed
   - Check `dist/` directory exists

2. **CORS Errors**
   - Verify frontend origin in CORS config
   - Check `CORS_ALLOWED_ORIGINS` environment variable

3. **WebSocket Connection Failed**
   - Ensure WebSocket endpoint is accessible
   - Check firewall and proxy settings

4. **Generation Errors**
   - Verify blueprint exists
   - Check request payload validation

### Debug Mode

```bash
# Enable debug logging
GIN_MODE=debug make run
```

### Log Analysis

```bash
# View structured logs
make run 2>&1 | jq '.'

# Filter error logs
make run 2>&1 | grep "ERROR"
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make lint` and `make test`
5. Submit pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
