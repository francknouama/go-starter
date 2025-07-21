# Development Guide

This guide covers setting up and running the go-starter web interface in development mode.

## Quick Start

### Prerequisites
- Docker and Docker Compose
- Node.js 20+ (for local development)
- Go 1.21+ (for local development)

### Option 1: Full Docker Development (Recommended)

Start the entire development stack with one command:

```bash
make dev-setup
```

This will:
1. Build the React app
2. Build Docker images for web server and React dev server
3. Start all services (web server, React dev server, PostgreSQL, Redis)

Access the application:
- **Web Interface**: http://localhost:8080
- **React Dev Server**: http://localhost:5173 (with hot reloading)
- **API Endpoints**: http://localhost:8080/api/v1/
- **PostgreSQL**: localhost:5432 (user: dev, password: dev_password)
- **Redis**: localhost:6379

### Option 2: Hybrid Development

For faster Go development with file watching:

```bash
# Start supporting services
docker-compose -f docker-compose.dev.yml up postgres redis -d

# Run Go backend locally
make backend-dev

# In another terminal, run React dev server
make web-dev
```

### Option 3: Pure Local Development

```bash
# Build React app
make web-build

# Start Go backend
go run ./cmd/web-server/main.go

# In another terminal, start React dev server
cd web && npm run dev
```

## Development Commands

### Docker Commands
```bash
# Start development environment
make dev-up

# Stop development environment  
make dev-down

# View logs from all services
make dev-logs

# Rebuild images
make dev-build

# Clean everything (containers, volumes, images)
make dev-clean

# Quick restart
make dev-restart

# Check service health
make dev-health

# Show service status
make dev-status
```

### Local Development Commands
```bash
# Build React app
make web-build

# Start React dev server (with hot reloading)
make web-dev

# Start Go backend server
make backend-dev
```

## Architecture Overview

### Services

1. **web-server** (Go + Gin)
   - Port: 8080
   - Serves React app and API endpoints
   - Uses filesystem access to blueprints for development

2. **web-ui** (React + Vite)
   - Port: 5173
   - Development server with hot reloading
   - Proxies API calls to web-server

3. **postgres** (PostgreSQL 16)
   - Port: 5432
   - For testing database-enabled blueprints
   - Database: go_starter_dev

4. **redis** (Redis 7)
   - Port: 6379
   - For testing cache-enabled blueprints

### Directory Structure

```
├── cmd/web-server/          # Go web server source
├── web/                     # React application
│   ├── src/                 # React source code
│   ├── dist/                # Built React app (served by Go)
│   └── Dockerfile.dev       # React dev container
├── blueprints/              # Go project templates
├── docker-compose.dev.yml   # Development services
├── Dockerfile.web-server    # Go web server container
└── Makefile.dev            # Development commands
```

## Development Workflow

### 1. Making Changes to React App

With Docker development:
1. Edit files in `web/src/`
2. Changes are automatically hot-reloaded at http://localhost:5173
3. API calls are proxied to the Go backend at http://localhost:8080

### 2. Making Changes to Go Backend

With Docker development:
1. Edit files in `cmd/web-server/` or `internal/`
2. Rebuild and restart: `make dev-restart`

For faster iteration, use hybrid development:
1. Stop the web-server container: `docker-compose -f docker-compose.dev.yml stop web-server`
2. Run locally: `make backend-dev`
3. Changes require manual restart

### 3. Adding New API Endpoints

1. Add handler in `internal/web/handlers/`
2. Register route in `cmd/web-server/main.go`
3. Update React services in `web/src/services/`
4. Restart backend (Docker or local)

### 4. Testing Blueprint Generation

The development environment includes PostgreSQL and Redis for testing blueprints that require databases:

```bash
# Test a web API blueprint with PostgreSQL
curl -X POST http://localhost:8080/api/v1/generate \
  -H "Content-Type: application/json" \
  -d '{
    "name": "test-api",
    "type": "web-api",
    "module": "github.com/test/test-api",
    "features": {
      "database": {
        "driver": "postgres"
      }
    }
  }'
```

## Environment Variables

### Go Backend
- `GIN_MODE`: gin mode (debug/release)
- `LOG_LEVEL`: logging level (debug/info/warn/error)

### React App
- `NODE_ENV`: development/production
- `VITE_API_URL`: Backend API URL (http://localhost:8080)

## Troubleshooting

### Common Issues

**Port conflicts:**
```bash
# Check what's using port 8080
lsof -i :8080

# Kill existing process
make dev-down
```

**Docker build failures:**
```bash
# Clean and rebuild
make dev-clean
make dev-build
```

**React app not loading:**
```bash
# Rebuild React app
make web-build
make dev-restart
```

**API not responding:**
```bash
# Check backend logs
docker-compose -f docker-compose.dev.yml logs web-server

# Check if blueprints are loaded
curl http://localhost:8080/api/v1/health
```

### Debugging

**Backend debugging:**
1. Add `fmt.Printf` or `slog.Debug` statements
2. View logs: `make dev-logs`

**Frontend debugging:**
1. Use browser dev tools
2. React dev server provides detailed error messages
3. Check network tab for API call failures

**Database debugging:**
```bash
# Connect to PostgreSQL
docker-compose -f docker-compose.dev.yml exec postgres psql -U dev -d go_starter_dev

# Connect to Redis
docker-compose -f docker-compose.dev.yml exec redis redis-cli
```

## Production Build

To test the production build locally:

```bash
# Build React app for production
make web-build

# Build Go binary
go build -o bin/web-server ./cmd/web-server

# Run production mode
GIN_MODE=release ./bin/web-server
```

The production build serves the React app from the embedded filesystem instead of the filesystem, providing a single binary deployment.