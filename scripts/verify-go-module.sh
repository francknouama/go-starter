#!/bin/bash

# Go Module Publishing Verification Script
# Verifies that the Go module is properly configured for public distribution

set -euo pipefail

# Configuration
MODULE_PATH="github.com/francknouama/go-starter"
VERSION="${1:-latest}"

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

echo "Go Module Publishing Verification"
echo "================================="
echo

# 1. Verify go.mod file
log_info "Verifying go.mod file..."
if [[ ! -f "go.mod" ]]; then
    log_error "go.mod file not found"
    exit 1
fi

# Check module path
if ! grep -q "module $MODULE_PATH" go.mod; then
    log_error "Module path mismatch in go.mod. Expected: $MODULE_PATH"
    exit 1
fi

# Check Go version
go_version=$(grep "^go " go.mod | cut -d' ' -f2)
if [[ -z "$go_version" ]]; then
    log_error "Go version not specified in go.mod"
    exit 1
fi

log_success "go.mod file is valid"
echo "  Module: $MODULE_PATH"
echo "  Go version: $go_version"

# 2. Verify module dependencies
log_info "Verifying module dependencies..."
if ! go mod verify; then
    log_error "Module verification failed"
    exit 1
fi

if ! go mod tidy; then
    log_error "go mod tidy failed"
    exit 1
fi

log_success "Module dependencies verified"

# 3. Check for proper module structure
log_info "Checking module structure..."

required_files=(
    "main.go"
    "go.mod"
    "go.sum"
    "README.md"
    "LICENSE"
)

missing_files=()
for file in "${required_files[@]}"; do
    if [[ ! -f "$file" ]]; then
        missing_files+=("$file")
    fi
done

if [[ ${#missing_files[@]} -gt 0 ]]; then
    log_warning "Missing recommended files: ${missing_files[*]}"
else
    log_success "All recommended files present"
fi

# 4. Verify CLI structure
log_info "Verifying CLI structure..."
if [[ -f "main.go" ]]; then
    log_success "Main CLI structure verified (root main.go found)"
elif [[ -f "cmd/go-starter/main.go" ]]; then
    log_success "Main CLI structure verified (cmd/go-starter/main.go found)"
else
    log_error "No main.go found in expected locations (root or cmd/go-starter/)"
    exit 1
fi

# 5. Check for proper package documentation
log_info "Checking package documentation..."

# Check for package comments in main packages
main_packages=$(find . -name "*.go" -path "*/cmd/*" | head -5)
documented_packages=0
total_packages=0

for pkg in $main_packages; do
    total_packages=$((total_packages + 1))
    if grep -q "^// Package\|^//go:build\|^// Command" "$pkg"; then
        documented_packages=$((documented_packages + 1))
    fi
done

if [[ $documented_packages -eq 0 ]]; then
    log_warning "No package documentation found in main packages"
else
    log_success "Package documentation found"
fi

# 6. Verify build for multiple platforms
log_info "Verifying cross-platform builds..."

platforms=(
    "linux/amd64"
    "darwin/amd64"
    "windows/amd64"
)

build_failures=0
for platform in "${platforms[@]}"; do
    IFS='/' read -r os arch <<< "$platform"
    # Try different main entry points
    if GOOS="$os" GOARCH="$arch" go build -o /dev/null . 2>/dev/null; then
        log_success "Build successful for $platform (root main.go)"
    elif GOOS="$os" GOARCH="$arch" go build -o /dev/null ./cmd/go-starter 2>/dev/null; then
        log_success "Build successful for $platform (cmd/go-starter/main.go)"
    else
        log_error "Build failed for $platform"
        build_failures=$((build_failures + 1))
    fi
done

if [[ $build_failures -gt 0 ]]; then
    log_error "$build_failures platform builds failed"
    exit 1
fi

# 7. Check version consistency
log_info "Checking version consistency..."

# Check if version is defined in code
version_files=(
    "internal/version/version.go"
    "cmd/go-starter/main.go"
)

version_defined=false
for file in "${version_files[@]}"; do
    if [[ -f "$file" ]] && grep -q -i "version\|Version" "$file"; then
        version_defined=true
        break
    fi
done

if $version_defined; then
    log_success "Version information found in source code"
else
    log_warning "No version information found in source code"
fi

# 8. Verify module can be fetched
if [[ "$VERSION" != "latest" ]]; then
    log_info "Verifying module can be fetched from proxy..."
    
    # Test if module can be fetched (this will work only if already published)
    temp_dir=$(mktemp -d)
    (
        cd "$temp_dir"
        if go mod init test-fetch && go get "$MODULE_PATH@$VERSION" 2>/dev/null; then
            log_success "Module can be fetched from Go proxy"
        else
            log_warning "Module cannot be fetched from Go proxy (may not be published yet)"
        fi
    )
    rm -rf "$temp_dir"
fi

# 9. Check for security vulnerabilities
log_info "Checking for security vulnerabilities..."
if command -v govulncheck >/dev/null 2>&1; then
    if govulncheck ./...; then
        log_success "No security vulnerabilities found"
    else
        log_warning "Security vulnerabilities detected"
    fi
else
    log_warning "govulncheck not available, skipping vulnerability check"
    echo "  Install with: go install golang.org/x/vuln/cmd/govulncheck@latest"
fi

# 10. Verify license
log_info "Verifying license..."
if [[ -f "LICENSE" ]]; then
    license_type=$(head -n 5 LICENSE | grep -i -E "(MIT|Apache|BSD|GPL)" | head -1 || echo "Unknown")
    log_success "License file found: $license_type"
else
    log_warning "No LICENSE file found"
fi

# 11. Check for CI/CD configuration
log_info "Checking CI/CD configuration..."
ci_files=(
    ".github/workflows/ci.yml"
    ".github/workflows/release.yml"
    ".goreleaser.yml"
)

ci_found=false
for file in "${ci_files[@]}"; do
    if [[ -f "$file" ]]; then
        ci_found=true
        log_success "CI/CD configuration found: $file"
    fi
done

if ! $ci_found; then
    log_warning "No CI/CD configuration found"
fi

# 12. Generate module information
log_info "Generating module information..."
cat > module-info.txt << EOF
Go Module Information
====================

Module Path: $MODULE_PATH
Go Version: $go_version
Generated: $(date -u +"%Y-%m-%d %H:%M:%S UTC")

Dependencies:
$(go list -m all | head -10)

Build Information:
- Cross-platform builds: Verified
- Package documentation: $([ $documented_packages -gt 0 ] && echo "Present" || echo "Missing")
- Security scan: $(command -v govulncheck >/dev/null 2>&1 && echo "Completed" || echo "Skipped")
- License: $([ -f LICENSE ] && echo "Present" || echo "Missing")

Module Ready for Publishing: $([ $build_failures -eq 0 ] && echo "YES" || echo "NO")
EOF

log_success "Module information saved to module-info.txt"

# Summary
echo
log_info "Verification Summary:"
echo "  ✓ Module structure verified"
echo "  ✓ Dependencies validated"
echo "  ✓ Cross-platform builds successful"
echo "  $([ $documented_packages -gt 0 ] && echo "✓" || echo "⚠") Package documentation"
echo "  $([ -f LICENSE ] && echo "✓" || echo "⚠") License file"
echo "  $([ $ci_found == true ] && echo "✓" || echo "⚠") CI/CD configuration"

if [[ $build_failures -eq 0 && -f "LICENSE" ]]; then
    echo
    log_success "Module is ready for publishing!"
    echo
    log_info "Publishing steps:"
    echo "  1. Tag version: git tag $VERSION"
    echo "  2. Push tags: git push origin --tags"
    echo "  3. Go proxy will automatically index the module"
    echo "  4. Verify with: go get $MODULE_PATH@$VERSION"
else
    echo
    log_warning "Module needs attention before publishing"
    echo "  - Fix build failures if any"
    echo "  - Add LICENSE file if missing"
    echo "  - Consider adding package documentation"
fi