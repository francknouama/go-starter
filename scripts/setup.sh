#!/bin/bash

# Go-Starter Development Environment Setup Script
# This script sets up the development environment for the go-starter project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running on supported OS
check_os() {
    log_info "Checking operating system..."
    
    case "$(uname -s)" in
        Linux*)
            OS="Linux"
            ;;
        Darwin*)
            OS="macOS"
            ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*)
            OS="Windows"
            ;;
        *)
            log_error "Unsupported operating system: $(uname -s)"
            exit 1
            ;;
    esac
    
    log_success "Operating system: $OS"
}

# Check if Go is installed and meets minimum version
check_go() {
    log_info "Checking Go installation..."
    
    if ! command -v go &> /dev/null; then
        log_error "Go is not installed. Please install Go 1.18 or later."
        log_info "Visit https://golang.org/doc/install for installation instructions."
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | sed 's/go//')
    GO_MAJOR=$(echo $GO_VERSION | cut -d. -f1)
    GO_MINOR=$(echo $GO_VERSION | cut -d. -f2)
    
    if [ "$GO_MAJOR" -lt 1 ] || ([ "$GO_MAJOR" -eq 1 ] && [ "$GO_MINOR" -lt 18 ]); then
        log_error "Go version $GO_VERSION is not supported. Please install Go 1.18 or later."
        exit 1
    fi
    
    log_success "Go version: $GO_VERSION"
}

# Check if Git is installed and configured
check_git() {
    log_info "Checking Git installation..."
    
    if ! command -v git &> /dev/null; then
        log_error "Git is not installed. Please install Git."
        exit 1
    fi
    
    GIT_VERSION=$(git --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
    log_success "Git version: $GIT_VERSION"
    
    # Check Git configuration
    if ! git config --global user.name &> /dev/null; then
        log_warning "Git user.name is not configured."
        log_info "Please run: git config --global user.name \"Your Name\""
    fi
    
    if ! git config --global user.email &> /dev/null; then
        log_warning "Git user.email is not configured."
        log_info "Please run: git config --global user.email \"your.email@example.com\""
    fi
}

# Install golangci-lint if not present
install_golangci_lint() {
    log_info "Checking golangci-lint installation..."
    
    if command -v golangci-lint &> /dev/null; then
        LINT_VERSION=$(golangci-lint --version | grep -oE '[0-9]+\.[0-9]+\.[0-9]+')
        log_success "golangci-lint version: $LINT_VERSION"
        return
    fi
    
    log_info "Installing golangci-lint..."
    
    case "$OS" in
        Linux|macOS)
            curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
            ;;
        Windows)
            log_info "Please install golangci-lint manually from: https://golangci-lint.run/usage/install/"
            log_warning "Skipping golangci-lint installation on Windows"
            return
            ;;
    esac
    
    if command -v golangci-lint &> /dev/null; then
        log_success "golangci-lint installed successfully"
    else
        log_error "Failed to install golangci-lint"
        exit 1
    fi
}

# Install additional development tools
install_dev_tools() {
    log_info "Installing development tools..."
    
    # Install goimports for code formatting
    if ! command -v goimports &> /dev/null; then
        log_info "Installing goimports..."
        go install golang.org/x/tools/cmd/goimports@latest
        log_success "goimports installed"
    else
        log_success "goimports already installed"
    fi
    
    # Install gotests for test generation
    if ! command -v gotests &> /dev/null; then
        log_info "Installing gotests..."
        go install github.com/cweill/gotests/gotests@latest
        log_success "gotests installed"
    else
        log_success "gotests already installed"
    fi
    
    # Install govulncheck for security scanning
    if ! command -v govulncheck &> /dev/null; then
        log_info "Installing govulncheck..."
        go install golang.org/x/vuln/cmd/govulncheck@latest
        log_success "govulncheck installed"
    else
        log_success "govulncheck already installed"
    fi
}

# Set up pre-commit hooks
setup_git_hooks() {
    log_info "Setting up Git hooks..."
    
    HOOKS_DIR=".git/hooks"
    
    if [ ! -d "$HOOKS_DIR" ]; then
        log_warning "Git hooks directory not found. Make sure you're in a Git repository."
        return
    fi
    
    # Create pre-commit hook
    cat > "$HOOKS_DIR/pre-commit" << 'EOF'
#!/bin/bash

# Pre-commit hook for go-starter
# Runs linting, formatting, and tests before commit

set -e

echo "Running pre-commit checks..."

# Check if there are any Go files to check
if ! git diff --cached --name-only --diff-filter=ACM | grep '\.go$' > /dev/null; then
    echo "No Go files to check"
    exit 0
fi

# Run go mod tidy
echo "Running go mod tidy..."
go mod tidy

# Run go fmt
echo "Running go fmt..."
go fmt ./...

# Run golangci-lint
echo "Running golangci-lint..."
if command -v golangci-lint &> /dev/null; then
    golangci-lint run
else
    echo "Warning: golangci-lint not found, skipping linting"
fi

# Run tests
echo "Running tests..."
go test -v ./...

echo "Pre-commit checks passed!"
EOF
    
    chmod +x "$HOOKS_DIR/pre-commit"
    log_success "Pre-commit hook installed"
}

# Verify project structure
verify_project_structure() {
    log_info "Verifying project structure..."
    
    REQUIRED_DIRS=(
        "cmd"
        "internal"
        "pkg"
        "templates"
        "tests"
        "scripts"
    )
    
    REQUIRED_FILES=(
        "go.mod"
        "Makefile"
        "README.md"
        "CLAUDE.md"
    )
    
    for dir in "${REQUIRED_DIRS[@]}"; do
        if [ ! -d "$dir" ]; then
            log_error "Required directory missing: $dir"
            exit 1
        fi
    done
    
    for file in "${REQUIRED_FILES[@]}"; do
        if [ ! -f "$file" ]; then
            log_error "Required file missing: $file"
            exit 1
        fi
    done
    
    log_success "Project structure verified"
}

# Download dependencies
download_dependencies() {
    log_info "Downloading Go dependencies..."
    
    go mod download
    go mod tidy
    
    log_success "Dependencies downloaded"
}

# Run initial build and tests
initial_build_test() {
    log_info "Running initial build..."
    
    if ! make build; then
        log_error "Build failed"
        exit 1
    fi
    
    log_success "Build successful"
    
    log_info "Running tests..."
    
    if ! make test; then
        log_error "Tests failed"
        exit 1
    fi
    
    log_success "Tests passed"
}

# Create development configuration
create_dev_config() {
    log_info "Creating development configuration..."
    
    CONFIG_FILE="$HOME/.go-starter.yaml"
    
    if [ ! -f "$CONFIG_FILE" ]; then
        cat > "$CONFIG_FILE" << EOF
profiles:
  default:
    author: ""
    email: ""
    license: "MIT"
    defaults:
      go_version: "1.21"
      framework: "gin"
      architecture: "standard"
      database:
        driver: ""
        orm: "gorm"
      auth:
        type: ""
current_profile: "default"
EOF
        log_success "Development configuration created at $CONFIG_FILE"
        log_info "Please edit $CONFIG_FILE to set your default author and email"
    else
        log_success "Development configuration already exists"
    fi
}

# Print setup summary
print_summary() {
    echo
    log_success "=== Development Environment Setup Complete ==="
    echo
    log_info "Next steps:"
    echo "  1. Edit ~/.go-starter.yaml to set your defaults"
    echo "  2. Run 'make build' to build the project"
    echo "  3. Run 'make test' to run tests"
    echo "  4. Run 'make lint' to run linting"
    echo "  5. Run './bin/go-starter --help' to see available commands"
    echo
    log_info "Available make targets:"
    echo "  make build      - Build the CLI binary"
    echo "  make test       - Run all tests"
    echo "  make lint       - Run linting"
    echo "  make install    - Install CLI globally"
    echo "  make clean      - Clean build artifacts"
    echo
    log_info "Happy coding! 🚀"
}

# Main setup function
main() {
    echo "=== Go-Starter Development Environment Setup ==="
    echo
    
    check_os
    check_go
    check_git
    verify_project_structure
    install_golangci_lint
    install_dev_tools
    download_dependencies
    setup_git_hooks
    initial_build_test
    create_dev_config
    print_summary
}

# Run main function
main "$@"