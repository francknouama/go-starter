#!/bin/bash

# Screenshot Generation Pipeline for go-starter Web UI Documentation
# 
# This script provides a complete automation pipeline for generating 
# high-quality screenshots of the Web UI for documentation purposes.
#
# Usage:
#   ./screenshot-pipeline.sh [command] [options]
#
# Commands:
#   setup     - Install dependencies and prepare environment
#   generate  - Generate all screenshots 
#   desktop   - Desktop screenshots only
#   mobile    - Mobile/responsive screenshots only
#   features  - Feature demonstration screenshots only
#   clean     - Clean up generated files
#   validate  - Validate screenshot quality and completeness
#   deploy    - Deploy screenshots to documentation

set -e  # Exit on error
set -u  # Exit on undefined variable

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
WEB_DIR="$PROJECT_ROOT/web"
DOCS_DIR="$PROJECT_ROOT/docs"
SCREENSHOT_DIR="$DOCS_DIR/screenshots"
SERVER_PORT=3000
BACKEND_PORT=8080
LOG_FILE="$SCRIPT_DIR/screenshot-generation.log"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "[$(date '+%Y-%m-%d %H:%M:%S')] $1" | tee -a "$LOG_FILE"
}

info() {
    log "${BLUE}INFO${NC}: $1"
}

success() {
    log "${GREEN}SUCCESS${NC}: $1"
}

warning() {
    log "${YELLOW}WARNING${NC}: $1"
}

error() {
    log "${RED}ERROR${NC}: $1"
}

# Check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
check_prerequisites() {
    info "Checking prerequisites..."
    
    # Check Node.js
    if ! command_exists node; then
        error "Node.js is not installed. Please install Node.js 18+ to continue."
        exit 1
    fi
    
    local node_version=$(node -v | sed 's/v//')
    local major_version=$(echo "$node_version" | cut -d. -f1)
    if [ "$major_version" -lt 18 ]; then
        error "Node.js version 18+ is required. Current version: $node_version"
        exit 1
    fi
    
    # Check npm
    if ! command_exists npm; then
        error "npm is not installed. Please install npm to continue."
        exit 1
    fi
    
    # Check Go (for backend server)
    if ! command_exists go; then
        error "Go is not installed. Please install Go 1.21+ to continue."
        exit 1
    fi
    
    # Check if we're in the right directory
    if [ ! -f "$WEB_DIR/package.json" ]; then
        error "Web directory not found. Please run this script from the project root."
        exit 1
    fi
    
    success "All prerequisites satisfied"
}

# Setup environment
setup_environment() {
    info "Setting up environment..."
    
    # Create necessary directories
    mkdir -p "$SCREENSHOT_DIR"
    mkdir -p "$SCRIPT_DIR/logs"
    
    # Install web dependencies
    cd "$WEB_DIR"
    info "Installing web dependencies..."
    npm ci
    
    # Install Playwright browsers if needed
    info "Installing Playwright browsers..."
    npx playwright install --with-deps chromium firefox webkit
    
    success "Environment setup complete"
}

# Start servers
start_servers() {
    info "Starting servers..."
    
    # Kill any existing processes on our ports
    pkill -f ":$SERVER_PORT" || true
    pkill -f ":$BACKEND_PORT" || true
    sleep 2
    
    # Start backend server
    cd "$PROJECT_ROOT"
    info "Starting backend server on port $BACKEND_PORT..."
    go run cmd/web-server/main.go > "$SCRIPT_DIR/logs/backend.log" 2>&1 &
    BACKEND_PID=$!
    
    # Wait for backend to be ready
    local backend_ready=false
    for i in {1..30}; do
        if curl -s "http://localhost:$BACKEND_PORT/api/health" >/dev/null 2>&1; then
            backend_ready=true
            break
        fi
        sleep 1
    done
    
    if [ "$backend_ready" = false ]; then
        warning "Backend server may not be ready, proceeding anyway"
    fi
    
    # Start frontend server  
    cd "$WEB_DIR"
    info "Starting frontend server on port $SERVER_PORT..."
    npm run dev > "$SCRIPT_DIR/logs/frontend.log" 2>&1 &
    FRONTEND_PID=$!
    
    # Wait for frontend to be ready
    local frontend_ready=false
    for i in {1..60}; do
        if curl -s "http://localhost:$SERVER_PORT" >/dev/null 2>&1; then
            frontend_ready=true
            break
        fi
        sleep 1
    done
    
    if [ "$frontend_ready" = false ]; then
        error "Frontend server failed to start"
        cleanup_servers
        exit 1
    fi
    
    success "Servers are running (Frontend: $SERVER_PORT, Backend: $BACKEND_PORT)"
}

# Stop servers
cleanup_servers() {
    info "Stopping servers..."
    
    if [ -n "${FRONTEND_PID:-}" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    
    if [ -n "${BACKEND_PID:-}" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    
    # Kill any remaining processes
    pkill -f ":$SERVER_PORT" 2>/dev/null || true
    pkill -f ":$BACKEND_PORT" 2>/dev/null || true
    
    success "Servers stopped"
}

# Generate screenshots
generate_screenshots() {
    local command=${1:-"all"}
    
    info "Generating screenshots: $command"
    
    cd "$WEB_DIR"
    
    # Ensure output directory exists
    mkdir -p "$SCREENSHOT_DIR"
    
    case "$command" in
        "desktop")
            npm run screenshots:desktop
            ;;
        "mobile") 
            npm run screenshots:mobile
            ;;
        "features")
            npm run screenshots:features
            ;;
        "blueprints")
            npm run screenshots:blueprints
            ;;
        "states")
            npm run screenshots:states
            ;;
        "all"|*)
            npm run screenshots:all
            ;;
    esac
    
    success "Screenshot generation completed"
}

# Validate screenshots
validate_screenshots() {
    info "Validating generated screenshots..."
    
    local total_screenshots=0
    local valid_screenshots=0
    
    if [ ! -d "$SCREENSHOT_DIR" ]; then
        error "Screenshot directory not found: $SCREENSHOT_DIR"
        return 1
    fi
    
    # Count total screenshots
    total_screenshots=$(find "$SCREENSHOT_DIR" -name "*.png" | wc -l | tr -d ' ')
    
    if [ "$total_screenshots" -eq 0 ]; then
        error "No screenshots found in $SCREENSHOT_DIR"
        return 1
    fi
    
    info "Found $total_screenshots screenshots"
    
    # Validate each screenshot
    while IFS= read -r -d '' file; do
        local filename=$(basename "$file")
        local filesize=$(stat -f%z "$file" 2>/dev/null || stat -c%s "$file" 2>/dev/null || echo 0)
        
        # Check file size (minimum 10KB for valid screenshots)
        if [ "$filesize" -gt 10240 ]; then
            valid_screenshots=$((valid_screenshots + 1))
        else
            warning "Screenshot may be corrupted (too small): $filename ($filesize bytes)"
        fi
    done < <(find "$SCREENSHOT_DIR" -name "*.png" -print0)
    
    info "Valid screenshots: $valid_screenshots/$total_screenshots"
    
    # Check for essential screenshots
    local essential_files=(
        "web-ui/01-landing-page.png"
        "web-ui/02-blueprint-gallery.png"  
        "workflows/01-initial-state.png"
        "responsive/desktop/full-interface.png"
        "responsive/mobile/standard.png"
    )
    
    local missing_essential=0
    for file in "${essential_files[@]}"; do
        if [ ! -f "$SCREENSHOT_DIR/$file" ]; then
            warning "Missing essential screenshot: $file"
            missing_essential=$((missing_essential + 1))
        fi
    done
    
    if [ "$missing_essential" -eq 0 ]; then
        success "All essential screenshots are present"
    else
        warning "$missing_essential essential screenshots are missing"
    fi
    
    # Generate validation report
    cat > "$SCREENSHOT_DIR/validation-report.json" << EOF
{
    "timestamp": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "total_screenshots": $total_screenshots,
    "valid_screenshots": $valid_screenshots, 
    "missing_essential": $missing_essential,
    "validation_passed": $([ "$valid_screenshots" -gt 0 ] && [ "$missing_essential" -eq 0 ] && echo "true" || echo "false")
}
EOF
    
    success "Validation complete - report saved to validation-report.json"
}

# Clean up generated files
clean_screenshots() {
    info "Cleaning up generated screenshots..."
    
    if [ -d "$SCREENSHOT_DIR" ]; then
        rm -rf "$SCREENSHOT_DIR"/*
        success "Screenshot directory cleaned"
    else
        info "No screenshots to clean"
    fi
}

# Deploy screenshots to documentation
deploy_screenshots() {
    info "Deploying screenshots to documentation..."
    
    # This could be extended to copy to docs, upload to CDN, etc.
    # For now, we'll just ensure they're in the right place
    
    if [ ! -d "$SCREENSHOT_DIR" ]; then
        error "No screenshots to deploy. Run generate first."
        return 1
    fi
    
    # Count screenshots
    local screenshot_count=$(find "$SCREENSHOT_DIR" -name "*.png" | wc -l | tr -d ' ')
    
    if [ "$screenshot_count" -eq 0 ]; then
        error "No screenshots found to deploy"
        return 1
    fi
    
    # Create documentation links
    info "Creating documentation index..."
    
    # The screenshot script already generates README.md, so we just verify it exists
    if [ -f "$SCREENSHOT_DIR/README.md" ]; then
        success "Screenshot documentation is ready"
    else
        warning "Screenshot index not found"
    fi
    
    success "Deployment complete - $screenshot_count screenshots ready for documentation"
}

# Show usage
show_usage() {
    cat << EOF
Screenshot Generation Pipeline for go-starter Web UI

Usage: $0 [COMMAND] [OPTIONS]

Commands:
    setup       Install dependencies and prepare environment
    generate    Generate all screenshots (default)
    desktop     Generate desktop screenshots only
    mobile      Generate mobile/responsive screenshots only
    features    Generate feature demonstration screenshots only
    blueprints  Generate blueprint showcase screenshots only
    states      Generate application state screenshots only
    clean       Clean up generated screenshot files
    validate    Validate screenshot quality and completeness
    deploy      Deploy screenshots to documentation
    help        Show this help message

Options:
    --no-servers    Don't start/stop servers (assume they're running)
    --verbose       Enable verbose logging
    --dry-run       Show what would be done without executing

Examples:
    $0 setup                    # Setup environment
    $0 generate                 # Generate all screenshots
    $0 desktop --no-servers     # Generate desktop screenshots with external servers
    $0 validate                 # Validate generated screenshots
    $0 clean                    # Clean up all generated files

Servers:
    Frontend: http://localhost:$SERVER_PORT
    Backend:  http://localhost:$BACKEND_PORT

Output:
    Screenshots: $SCREENSHOT_DIR
    Logs:        $SCRIPT_DIR/logs/

EOF
}

# Main execution
main() {
    local command=${1:-"generate"}
    local no_servers=false
    local verbose=false
    local dry_run=false
    
    # Parse options
    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-servers)
                no_servers=true
                shift
                ;;
            --verbose)
                verbose=true
                shift
                ;;
            --dry-run)
                dry_run=true
                shift
                ;;
            help|--help|-h)
                show_usage
                exit 0
                ;;
            *)
                if [ -z "${command_set:-}" ]; then
                    command="$1"
                    command_set=true
                fi
                shift
                ;;
        esac
    done
    
    # Setup logging
    if [ "$verbose" = true ]; then
        set -x
    fi
    
    echo "# go-starter Screenshot Generation Pipeline"
    echo "=========================================="
    echo ""
    
    # Initialize log file
    echo "# Screenshot Generation Log - $(date)" > "$LOG_FILE"
    
    # Execute command
    case "$command" in
        "setup")
            check_prerequisites
            setup_environment
            success "Setup complete"
            ;;
            
        "generate"|"desktop"|"mobile"|"features"|"blueprints"|"states")
            if [ "$dry_run" = true ]; then
                info "DRY RUN: Would generate $command screenshots"
                exit 0
            fi
            
            check_prerequisites
            
            if [ "$no_servers" = false ]; then
                start_servers
                trap cleanup_servers EXIT
            fi
            
            generate_screenshots "$command"
            validate_screenshots
            
            if [ "$no_servers" = false ]; then
                cleanup_servers
            fi
            
            success "Screenshot generation pipeline completed successfully"
            ;;
            
        "clean")
            clean_screenshots
            ;;
            
        "validate")
            validate_screenshots
            ;;
            
        "deploy")
            validate_screenshots
            deploy_screenshots
            ;;
            
        *)
            error "Unknown command: $command"
            show_usage
            exit 1
            ;;
    esac
}

# Trap for cleanup on script exit
trap 'cleanup_servers' INT TERM

# Execute main function
main "$@"