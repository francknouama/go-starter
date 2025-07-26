#!/bin/bash

# Package Distribution Script for go-starter
# Creates distribution packages for multiple platforms and package managers

set -euo pipefail

# Configuration
PROJECT_NAME="go-starter"
GITHUB_REPO="francknouama/go-starter"
VERSION="${1:-}"
DIST_DIR="dist"
FORMULA_DIR="homebrew-formula"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}INFO:${NC} $1"
}

log_success() {
    echo -e "${GREEN}SUCCESS:${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}WARNING:${NC} $1"
}

log_error() {
    echo -e "${RED}ERROR:${NC} $1"
}

# Check if version is provided
if [ -z "$VERSION" ]; then
    log_error "Version is required. Usage: $0 <version>"
    echo "Example: $0 v1.4.0"
    exit 1
fi

# Validate version format
if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    log_error "Invalid version format. Expected: vX.Y.Z (e.g., v1.4.0)"
    exit 1
fi

log_info "Starting package distribution for $PROJECT_NAME $VERSION"

# Clean and create distribution directory
rm -rf "$DIST_DIR"
mkdir -p "$DIST_DIR"

# Verify Go module is properly configured
log_info "Verifying Go module configuration..."
if ! go mod verify; then
    log_error "Go module verification failed"
    exit 1
fi

if ! go mod tidy; then
    log_error "Go mod tidy failed"
    exit 1
fi

log_success "Go module verification passed"

# Build for multiple platforms
log_info "Building binaries for multiple platforms..."

platforms=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    output_name="${PROJECT_NAME}-${VERSION}-${os}-${arch}"
    
    if [[ "$os" == "windows" ]]; then
        output_name="${output_name}.exe"
    fi
    
    log_info "Building for $os/$arch..."
    
    if GOOS="$os" GOARCH="$arch" CGO_ENABLED=0 go build \
        -ldflags="-s -w -X 'github.com/$GITHUB_REPO/internal/version.Version=$VERSION' -X 'github.com/$GITHUB_REPO/internal/version.BuildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)'" \
        -o "$DIST_DIR/$output_name" \
        ./cmd/go-starter; then
        
        log_success "Built $output_name"
        
        # Create tar.gz for non-Windows platforms
        if [[ "$os" != "windows" ]]; then
            tar_name="${PROJECT_NAME}-${VERSION}-${os}-${arch}.tar.gz"
            (cd "$DIST_DIR" && tar -czf "$tar_name" "$output_name")
            rm "$DIST_DIR/$output_name"
            log_success "Created archive $tar_name"
        else
            # Create zip for Windows
            zip_name="${PROJECT_NAME}-${VERSION}-${os}-${arch}.zip"
            (cd "$DIST_DIR" && zip -q "$zip_name" "$output_name")
            rm "$DIST_DIR/$output_name"
            log_success "Created archive $zip_name"
        fi
    else
        log_error "Failed to build for $os/$arch"
        exit 1
    fi
done

# Generate checksums
log_info "Generating checksums..."
(cd "$DIST_DIR" && sha256sum * > checksums.txt)
log_success "Generated checksums.txt"

# Create Homebrew formula
log_info "Creating Homebrew formula..."
mkdir -p "$FORMULA_DIR"

# Get checksums for macOS binaries
macos_amd64_sha256=$(cd "$DIST_DIR" && sha256sum "${PROJECT_NAME}-${VERSION}-darwin-amd64.tar.gz" | cut -d' ' -f1)
macos_arm64_sha256=$(cd "$DIST_DIR" && sha256sum "${PROJECT_NAME}-${VERSION}-darwin-arm64.tar.gz" | cut -d' ' -f1)

cat > "$FORMULA_DIR/${PROJECT_NAME}.rb" << EOF
class GoStarter < Formula
  desc "Comprehensive Go project generator with modern best practices"
  homepage "https://github.com/$GITHUB_REPO"
  version "$VERSION"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/$GITHUB_REPO/releases/download/$VERSION/${PROJECT_NAME}-${VERSION}-darwin-amd64.tar.gz"
      sha256 "$macos_amd64_sha256"

      def install
        bin.install "${PROJECT_NAME}-${VERSION}-darwin-amd64" => "$PROJECT_NAME"
      end
    end

    if Hardware::CPU.arm?
      url "https://github.com/$GITHUB_REPO/releases/download/$VERSION/${PROJECT_NAME}-${VERSION}-darwin-arm64.tar.gz"
      sha256 "$macos_arm64_sha256"

      def install
        bin.install "${PROJECT_NAME}-${VERSION}-darwin-arm64" => "$PROJECT_NAME"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/$GITHUB_REPO/releases/download/$VERSION/${PROJECT_NAME}-${VERSION}-linux-amd64.tar.gz"
      sha256 "$(cd "$DIST_DIR" && sha256sum "${PROJECT_NAME}-${VERSION}-linux-amd64.tar.gz" | cut -d' ' -f1)"

      def install
        bin.install "${PROJECT_NAME}-${VERSION}-linux-amd64" => "$PROJECT_NAME"
      end
    end

    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/$GITHUB_REPO/releases/download/$VERSION/${PROJECT_NAME}-${VERSION}-linux-arm64.tar.gz"
      sha256 "$(cd "$DIST_DIR" && sha256sum "${PROJECT_NAME}-${VERSION}-linux-arm64.tar.gz" | cut -d' ' -f1)"

      def install
        bin.install "${PROJECT_NAME}-${VERSION}-linux-arm64" => "$PROJECT_NAME"
      end
    end
  end

  test do
    system "#{bin}/$PROJECT_NAME", "version"
    system "#{bin}/$PROJECT_NAME", "--help"
  end
end
EOF

log_success "Created Homebrew formula: $FORMULA_DIR/${PROJECT_NAME}.rb"

# Create installation script
log_info "Creating installation script..."
cat > "$DIST_DIR/install.sh" << 'EOF'
#!/bin/bash
# Installation script for go-starter

set -euo pipefail

# Configuration
PROJECT_NAME="go-starter"
GITHUB_REPO="francknouama/go-starter"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log_info() { echo -e "${BLUE}INFO:${NC} $1"; }
log_success() { echo -e "${GREEN}SUCCESS:${NC} $1"; }
log_warning() { echo -e "${YELLOW}WARNING:${NC} $1"; }
log_error() { echo -e "${RED}ERROR:${NC} $1"; }

# Detect platform
detect_platform() {
    local os arch
    
    case "$(uname -s)" in
        Darwin*) os="darwin" ;;
        Linux*) os="linux" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *) log_error "Unsupported operating system: $(uname -s)"; exit 1 ;;
    esac
    
    case "$(uname -m)" in
        x86_64) arch="amd64" ;;
        arm64|aarch64) arch="arm64" ;;
        *) log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac
    
    echo "${os}-${arch}"
}

# Get latest version from GitHub
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "https://api.github.com/repos/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        log_error "curl or wget is required for installation"
        exit 1
    fi
}

# Download and install
install_go_starter() {
    local version="${1:-$(get_latest_version)}"
    local platform
    platform=$(detect_platform)
    
    log_info "Installing $PROJECT_NAME $version for $platform..."
    
    # Create temporary directory
    local temp_dir
    temp_dir=$(mktemp -d)
    trap "rm -rf $temp_dir" EXIT
    
    # Download URL
    local download_url="https://github.com/$GITHUB_REPO/releases/download/$version/${PROJECT_NAME}-${version}-${platform}.tar.gz"
    local archive_file="$temp_dir/${PROJECT_NAME}-${version}-${platform}.tar.gz"
    
    log_info "Downloading from $download_url..."
    
    if command -v curl >/dev/null 2>&1; then
        if ! curl -L -o "$archive_file" "$download_url"; then
            log_error "Failed to download $PROJECT_NAME"
            exit 1
        fi
    elif command -v wget >/dev/null 2>&1; then
        if ! wget -O "$archive_file" "$download_url"; then
            log_error "Failed to download $PROJECT_NAME"
            exit 1
        fi
    fi
    
    # Extract archive
    log_info "Extracting archive..."
    tar -xzf "$archive_file" -C "$temp_dir"
    
    # Find the binary
    local binary_path
    binary_path=$(find "$temp_dir" -name "${PROJECT_NAME}-${version}-${platform}" -o -name "$PROJECT_NAME" | head -1)
    
    if [[ ! -f "$binary_path" ]]; then
        log_error "Binary not found in archive"
        exit 1
    fi
    
    # Install binary
    log_info "Installing to $INSTALL_DIR..."
    
    # Check if we need sudo
    if [[ ! -w "$INSTALL_DIR" ]]; then
        if command -v sudo >/dev/null 2>&1; then
            sudo cp "$binary_path" "$INSTALL_DIR/$PROJECT_NAME"
            sudo chmod +x "$INSTALL_DIR/$PROJECT_NAME"
        else
            log_error "No write permission to $INSTALL_DIR and sudo not available"
            exit 1
        fi
    else
        cp "$binary_path" "$INSTALL_DIR/$PROJECT_NAME"
        chmod +x "$INSTALL_DIR/$PROJECT_NAME"
    fi
    
    log_success "$PROJECT_NAME $version installed successfully!"
    log_info "Run '$PROJECT_NAME --help' to get started"
}

# Main execution
main() {
    local version="${1:-}"
    
    echo "go-starter Installation Script"
    echo "=============================="
    echo
    
    # Check dependencies
    if ! command -v tar >/dev/null 2>&1; then
        log_error "tar is required for installation"
        exit 1
    fi
    
    install_go_starter "$version"
}

main "$@"
EOF

chmod +x "$DIST_DIR/install.sh"
log_success "Created installation script: $DIST_DIR/install.sh"

# Create package metadata
log_info "Creating package metadata..."
cat > "$DIST_DIR/package-info.json" << EOF
{
  "name": "$PROJECT_NAME",
  "version": "$VERSION",
  "description": "Comprehensive Go project generator with modern best practices",
  "repository": "https://github.com/$GITHUB_REPO",
  "license": "MIT",
  "author": "Franck Nouama",
  "homepage": "https://github.com/$GITHUB_REPO",
  "bugs": "https://github.com/$GITHUB_REPO/issues",
  "build_info": {
    "build_date": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
    "go_version": "$(go version | cut -d' ' -f3)",
    "platforms": [
      "linux/amd64",
      "linux/arm64", 
      "darwin/amd64",
      "darwin/arm64",
      "windows/amd64",
      "windows/arm64"
    ]
  },
  "installation": {
    "homebrew": "brew install $GITHUB_REPO",
    "script": "curl -sSL https://github.com/$GITHUB_REPO/releases/download/$VERSION/install.sh | bash",
    "manual": "Download binary from https://github.com/$GITHUB_REPO/releases/tag/$VERSION"
  }
}
EOF

log_success "Created package metadata: $DIST_DIR/package-info.json"

# Generate release notes template
log_info "Generating release notes template..."
cat > "$DIST_DIR/release-notes.md" << EOF
# go-starter $VERSION Release Notes

## 🚀 What's New

### Phase 4D: Production Hardening (Major Update)
- **Advanced AST Operations**: Sophisticated code transformation with 91% reduction capability
- **Automated Test Generation**: AI-powered test creation from blueprint analysis  
- **Continuous Coverage Monitoring**: Real-time quality gates and regression tracking
- **Self-Maintaining Test Infrastructure**: Autonomous test management and optimization

### Enhanced Features
- **Progressive Disclosure System**: Smart help filtering and complexity-aware generation
- **Simplified Logger System**: 60-90% code reduction across all blueprints
- **Comprehensive ATDD Quality Tests**: 60% performance improvement with parallel execution
- **GitHub Actions Integration**: Full CI/CD pipeline with release automation

## 📦 Installation

### Homebrew (macOS/Linux)
\`\`\`bash
brew tap $GITHUB_REPO
brew install go-starter
\`\`\`

### Quick Install Script
\`\`\`bash
curl -sSL https://github.com/$GITHUB_REPO/releases/download/$VERSION/install.sh | bash
\`\`\`

### Go Install
\`\`\`bash
go install github.com/$GITHUB_REPO/cmd/go-starter@$VERSION
\`\`\`

### Manual Download
Download the appropriate binary for your platform from the [releases page](https://github.com/$GITHUB_REPO/releases/tag/$VERSION).

## 🔧 Usage

### Basic Usage
\`\`\`bash
# Interactive mode (beginner-friendly)
go-starter new

# Direct generation with progressive disclosure
go-starter new my-app --type=cli --complexity=simple
\`\`\`

### Advanced Usage
\`\`\`bash
# Advanced mode with all options
go-starter new --advanced

# Complex project generation
go-starter new enterprise-api \\
  --type=web-api \\
  --architecture=hexagonal \\
  --database-driver=postgres \\
  --logger=zap \\
  --advanced
\`\`\`

## 📋 Supported Project Types

| Type | Complexity | Use Case |
|------|------------|----------|
| **Simple CLI** | Beginner | Quick utilities (8 files) |
| **Standard CLI** | Intermediate | Production CLIs (29 files) |
| **Web API** | Intermediate | REST APIs with multiple architectures |
| **Lambda** | Beginner-Advanced | AWS serverless functions |
| **Microservice** | Advanced | Distributed systems with gRPC |
| **Workspace** | Advanced | Multi-module monorepos |

## 🏗️ Architecture Patterns

- **Standard**: Traditional layered architecture
- **Clean Architecture**: Enterprise-grade separation of concerns
- **DDD**: Domain-driven design with rich models
- **Hexagonal**: Ports & adapters for maximum testability
- **Event-driven**: CQRS and event sourcing patterns

## 🔍 Quality Assurance

This release includes comprehensive testing infrastructure:
- **70 Test Functions**: Across 4 major infrastructure components
- **Multi-platform Builds**: Linux, macOS, Windows (AMD64 & ARM64)
- **Security Scanning**: CodeQL analysis and dependency checks
- **Performance Monitoring**: Automated regression detection

## 📈 Performance Improvements

- **Logger System**: 60-90% code reduction
- **ATDD Tests**: 60% performance improvement with parallel execution
- **Blueprint Generation**: Optimized file creation and template processing
- **Memory Usage**: Reduced memory footprint in code generation

## 🐛 Bug Fixes

- Fixed errcheck linting violations across all components
- Resolved template generation issues with unused imports
- Corrected AST operation complexity calculations
- Fixed concurrency issues in test infrastructure

## 🔒 Security

- Enhanced input validation for all user inputs
- Secure template processing with path traversal protection
- Dependency scanning with automated updates
- CodeQL security analysis integration

## 📚 Documentation

- Complete blueprint selection guide with complexity matrix
- Progressive disclosure system documentation
- Logger comparison and selection guide  
- ATDD strategy and implementation details

## 🙏 Acknowledgments

Special thanks to all contributors and the Go community for their continued support and feedback.

## 🔗 Links

- **Repository**: https://github.com/$GITHUB_REPO
- **Documentation**: https://github.com/$GITHUB_REPO/blob/main/README.md
- **Issues**: https://github.com/$GITHUB_REPO/issues
- **Discussions**: https://github.com/$GITHUB_REPO/discussions

---

**Full Changelog**: https://github.com/$GITHUB_REPO/compare/v1.3.0...$VERSION
EOF

log_success "Generated release notes template: $DIST_DIR/release-notes.md"

# Create Go module proxy verification
log_info "Verifying Go module proxy compatibility..."
if go list -m "github.com/$GITHUB_REPO@$VERSION" >/dev/null 2>&1; then
    log_success "Go module proxy verification passed"
else
    log_warning "Go module proxy verification failed (version may not be published yet)"
fi

# Summary
log_info "Package distribution summary:"
echo "  - Built binaries for 6 platforms"
echo "  - Created Homebrew formula"
echo "  - Generated installation script"
echo "  - Created package metadata"
echo "  - Generated release notes template"
echo "  - Calculated checksums"
echo
log_success "Package distribution completed successfully!"
echo
log_info "Next steps:"
echo "  1. Review generated files in $DIST_DIR/ and $FORMULA_DIR/"
echo "  2. Test installation script on different platforms"
echo "  3. Submit Homebrew formula to tap repository"
echo "  4. Upload binaries to GitHub releases"
echo "  5. Update package managers and distribution channels"