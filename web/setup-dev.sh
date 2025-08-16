#!/bin/bash

# Development Setup Script for go-starter Web UI
# Sets up the complete development environment with backend API and frontend

set -e

echo "🚀 Setting up go-starter Web UI Development Environment"
echo "======================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Check dependencies
check_dependencies() {
    print_info "Checking dependencies..."
    
    # Check Go
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.21 or later."
        exit 1
    fi
    
    local go_version=$(go version | awk '{print $3}' | sed 's/go//')
    print_success "Go $go_version found"
    
    # Check Node.js
    if ! command -v node &> /dev/null; then
        print_error "Node.js is not installed. Please install Node.js 18 or later."
        exit 1
    fi
    
    local node_version=$(node --version)
    print_success "Node.js $node_version found"
    
    # Check npm
    if ! command -v npm &> /dev/null; then
        print_error "npm is not installed. Please install npm."
        exit 1
    fi
    
    local npm_version=$(npm --version)
    print_success "npm $npm_version found"
    
    # Check if we're in the web directory
    if [ ! -f "package.json" ]; then
        print_error "Please run this script from the web/ directory"
        exit 1
    fi
}

# Install dependencies
install_dependencies() {
    print_info "Installing dependencies..."
    
    # Install Go dependencies
    if [ ! -f "go.mod" ]; then
        print_info "Initializing Go module..."
        go mod init github.com/francknouama/go-starter/web
    fi
    
    print_info "Installing Go dependencies..."
    go mod tidy
    
    # Install Node.js dependencies
    print_info "Installing Node.js dependencies..."
    npm install
    
    print_success "Dependencies installed"
}

# Build frontend
build_frontend() {
    print_info "Building frontend..."
    
    # Build React app
    npm run build
    
    if [ -d "dist" ]; then
        print_success "Frontend built successfully"
    else
        print_error "Frontend build failed"
        exit 1
    fi
}

# Generate Go embedded files
generate_go_files() {
    print_info "Generating Go embedded files..."
    
    # Generate embedded static files
    go generate ./...
    
    print_success "Go files generated"
}

# Build backend
build_backend() {
    print_info "Building backend..."
    
    # Build web server
    go build -o bin/web-server ./cmd/web-server
    
    if [ -f "bin/web-server" ]; then
        print_success "Backend built successfully"
    else
        print_error "Backend build failed"
        exit 1
    fi
}

# Create startup scripts
create_scripts() {
    print_info "Creating startup scripts..."
    
    # Development script
    cat > start-dev.sh << 'EOF'
#!/bin/bash

# Start development server with hot reloading

echo "🚀 Starting go-starter Web UI Development Server"
echo "================================================"

# Check if air is installed
if ! command -v air &> /dev/null; then
    echo "Installing air for hot reloading..."
    go install github.com/cosmtrek/air@latest
fi

# Start with air for hot reloading
air -c .air.toml
EOF

    # Production script
    cat > start-prod.sh << 'EOF'
#!/bin/bash

# Start production server

echo "🚀 Starting go-starter Web UI Production Server"
echo "==============================================="

# Build if needed
if [ ! -f "bin/web-server" ]; then
    echo "Building production server..."
    go build -o bin/web-server ./cmd/web-server
fi

# Start server
./bin/web-server
EOF

    chmod +x start-dev.sh start-prod.sh
    print_success "Startup scripts created"
}

# Create air configuration
create_air_config() {
    print_info "Creating air configuration..."
    
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/web-server"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", "dist"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  poll = false
  poll_interval = 0
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
EOF

    print_success "Air configuration created"
}

# Test the setup
test_setup() {
    print_info "Testing the setup..."
    
    # Start server in background for testing
    print_info "Starting server for testing..."
    ./bin/web-server &
    SERVER_PID=$!
    
    # Wait for server to start
    sleep 3
    
    # Test health endpoint
    if curl -s http://localhost:8080/api/v1/health/simple | grep -q "ok"; then
        print_success "Server is responding correctly"
    else
        print_warning "Server health check failed - may need manual verification"
    fi
    
    # Stop test server
    kill $SERVER_PID 2>/dev/null || true
    wait $SERVER_PID 2>/dev/null || true
}

# Create README for web development
create_readme() {
    print_info "Creating development README..."
    
    cat > README-DEV.md << 'EOF'
# go-starter Web UI Development

## Quick Start

1. **Development Mode (with hot reloading):**
   ```bash
   ./start-dev.sh
   ```
   
2. **Production Mode:**
   ```bash
   ./start-prod.sh
   ```

3. **Run Integration Tests:**
   ```bash
   ./test-integration.sh
   ```

## Development Workflow

### Frontend Development
- React + TypeScript + Tailwind CSS
- Vite for build tooling
- Components in `src/components/`
- API integration in `src/services/api.ts`

### Backend Development
- Gin web framework
- API handlers in `internal/web/handlers/`
- WebSocket support in `internal/web/websocket/`
- Middleware in `internal/web/middleware/`

### API Endpoints
- Health: `GET /api/v1/health`
- Config: `GET /api/v1/config`
- Blueprints: `GET /api/v1/blueprints`
- Preview: `POST /api/v1/preview`
- Generate: `POST /api/v1/generate`
- WebSocket: `GET /api/v1/ws`

### File Structure
```
web/
├── cmd/web-server/          # Main server
├── internal/web/            # Backend code
│   ├── handlers/           # API handlers
│   ├── middleware/         # HTTP middleware
│   ├── models/            # Request/response models
│   └── websocket/         # WebSocket support
├── src/                    # Frontend React app
│   ├── components/        # React components
│   ├── hooks/            # Custom React hooks
│   ├── services/         # API service layer
│   └── utils/            # Utilities
└── dist/                  # Built frontend
```

## Environment Variables

- `PORT`: Server port (default: 8080)
- `NODE_ENV`: Environment (development/production)
- `VITE_API_BASE_URL`: API base URL for frontend

## Troubleshooting

1. **Server won't start:**
   - Check if port 8080 is available
   - Ensure Go dependencies are installed: `go mod tidy`

2. **Frontend build fails:**
   - Ensure Node.js dependencies: `npm install`
   - Clear cache: `npm run clean`

3. **API errors:**
   - Check server logs
   - Verify API endpoints with: `curl http://localhost:8080/api/v1/health`

4. **WebSocket issues:**
   - Ensure WebSocket endpoint is accessible
   - Check browser console for connection errors
EOF

    print_success "Development README created"
}

# Main setup function
main() {
    echo "Starting setup at $(date)"
    echo ""
    
    check_dependencies
    install_dependencies
    build_frontend
    generate_go_files
    build_backend
    create_scripts
    create_air_config
    test_setup
    create_readme
    
    echo ""
    echo "=============================================="
    echo "🎉 Setup Complete!"
    echo "=============================================="
    echo ""
    echo "Next steps:"
    echo "1. Start development server:"
    echo "   ./start-dev.sh"
    echo ""
    echo "2. Open your browser to:"
    echo "   http://localhost:8080"
    echo ""
    echo "3. Run integration tests:"
    echo "   ./test-integration.sh"
    echo ""
    echo "4. Read development guide:"
    echo "   cat README-DEV.md"
    echo ""
    print_success "go-starter Web UI is ready for development!"
}

# Run main function
main "$@"