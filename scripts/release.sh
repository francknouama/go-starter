#!/bin/bash

# Go-Starter Release Script
# This script handles the release process for the go-starter project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_NAME="go-starter"
BINARY_NAME="go-starter"
RELEASE_DIR="dist"
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

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

# Show usage
show_usage() {
    echo "Usage: $0 [OPTIONS] <version>"
    echo
    echo "Release script for $PROJECT_NAME"
    echo
    echo "Arguments:"
    echo "  version     Version to release (e.g., v1.0.0, v1.2.3-beta.1)"
    echo
    echo "Options:"
    echo "  -h, --help     Show this help message"
    echo "  -d, --dry-run  Perform a dry run without making changes"
    echo "  -p, --pre      Create a pre-release"
    echo "  -f, --force    Force release even if working directory is dirty"
    echo
    echo "Examples:"
    echo "  $0 v1.0.0                    # Create release v1.0.0"
    echo "  $0 --pre v1.0.0-rc.1        # Create pre-release v1.0.0-rc.1"
    echo "  $0 --dry-run v1.0.0         # Dry run for v1.0.0"
}

# Parse command line arguments
parse_args() {
    DRY_RUN=false
    PRE_RELEASE=false
    FORCE=false
    VERSION=""
    
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_usage
                exit 0
                ;;
            -d|--dry-run)
                DRY_RUN=true
                shift
                ;;
            -p|--pre)
                PRE_RELEASE=true
                shift
                ;;
            -f|--force)
                FORCE=true
                shift
                ;;
            -*)
                log_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
            *)
                if [ -z "$VERSION" ]; then
                    VERSION="$1"
                else
                    log_error "Multiple versions specified"
                    show_usage
                    exit 1
                fi
                shift
                ;;
        esac
    done
    
    if [ -z "$VERSION" ]; then
        log_error "Version is required"
        show_usage
        exit 1
    fi
}

# Validate version format
validate_version() {
    log_info "Validating version format..."
    
    # Check if version starts with 'v'
    if [[ ! "$VERSION" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?$ ]]; then
        log_error "Invalid version format: $VERSION"
        log_info "Version must be in format: v1.2.3 or v1.2.3-beta.1"
        exit 1
    fi
    
    log_success "Version format is valid: $VERSION"
}

# Check if working directory is clean
check_working_directory() {
    log_info "Checking working directory status..."
    
    if [ "$FORCE" = true ]; then
        log_warning "Forcing release with potentially dirty working directory"
        return
    fi
    
    if ! git diff --quiet || ! git diff --cached --quiet; then
        log_error "Working directory is not clean"
        log_info "Please commit or stash your changes before releasing"
        log_info "Use --force to override this check"
        exit 1
    fi
    
    log_success "Working directory is clean"
}

# Check if version tag already exists
check_version_tag() {
    log_info "Checking if version tag exists..."
    
    if git tag --list | grep -q "^$VERSION$"; then
        log_error "Version tag $VERSION already exists"
        exit 1
    fi
    
    log_success "Version tag $VERSION is available"
}

# Update version in code
update_version() {
    log_info "Updating version in code..."
    
    # Update version in main.go or version file
    if [ -f "cmd/version.go" ]; then
        if [ "$DRY_RUN" = false ]; then
            sed -i.bak "s/Version = \".*\"/Version = \"${VERSION#v}\"/" cmd/version.go
            rm -f cmd/version.go.bak
        fi
        log_success "Version updated in cmd/version.go"
    else
        log_warning "Version file not found, skipping version update"
    fi
}

# Run pre-release checks
run_pre_release_checks() {
    log_info "Running pre-release checks..."
    
    # Run tests
    log_info "Running tests..."
    if ! make test; then
        log_error "Tests failed"
        exit 1
    fi
    log_success "Tests passed"
    
    # Run linting
    log_info "Running linting..."
    if ! make lint; then
        log_error "Linting failed"
        exit 1
    fi
    log_success "Linting passed"
    
    # Run security checks
    if command -v govulncheck &> /dev/null; then
        log_info "Running security checks..."
        if ! govulncheck ./...; then
            log_error "Security checks failed"
            exit 1
        fi
        log_success "Security checks passed"
    else
        log_warning "govulncheck not found, skipping security checks"
    fi
}

# Build binaries for all platforms
build_binaries() {
    log_info "Building binaries for all platforms..."
    
    # Clean and create release directory
    if [ "$DRY_RUN" = false ]; then
        rm -rf "$RELEASE_DIR"
        mkdir -p "$RELEASE_DIR"
    fi
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r -a platform_split <<< "$platform"
        GOOS="${platform_split[0]}"
        GOARCH="${platform_split[1]}"
        
        output_name="$BINARY_NAME"
        if [ "$GOOS" = "windows" ]; then
            output_name="$output_name.exe"
        fi
        
        archive_name="$PROJECT_NAME-${VERSION#v}-$GOOS-$GOARCH"
        if [ "$GOOS" = "windows" ]; then
            archive_name="$archive_name.zip"
        else
            archive_name="$archive_name.tar.gz"
        fi
        
        log_info "Building for $GOOS/$GOARCH..."
        
        if [ "$DRY_RUN" = false ]; then
            env GOOS="$GOOS" GOARCH="$GOARCH" go build \
                -ldflags="-s -w -X main.version=${VERSION#v}" \
                -o "$RELEASE_DIR/$output_name" \
                ./main.go
            
            # Create archive
            if [ "$GOOS" = "windows" ]; then
                (cd "$RELEASE_DIR" && zip "$archive_name" "$output_name")
            else
                (cd "$RELEASE_DIR" && tar -czf "$archive_name" "$output_name")
            fi
            
            # Remove binary (keep only archive)
            rm "$RELEASE_DIR/$output_name"
        fi
        
        log_success "Built $archive_name"
    done
    
    log_success "All binaries built successfully"
}

# Generate checksums
generate_checksums() {
    log_info "Generating checksums..."
    
    if [ "$DRY_RUN" = false ]; then
        (cd "$RELEASE_DIR" && shasum -a 256 *.tar.gz *.zip > checksums.txt)
    fi
    
    log_success "Checksums generated"
}

# Create and push git tag
create_git_tag() {
    log_info "Creating git tag..."
    
    if [ "$DRY_RUN" = false ]; then
        git add -A
        git commit -m "Release $VERSION" || log_warning "No changes to commit"
        git tag -a "$VERSION" -m "Release $VERSION"
        
        log_info "Pushing tag to origin..."
        git push origin "$VERSION"
        git push origin main
    fi
    
    log_success "Git tag $VERSION created and pushed"
}

# Create GitHub release
create_github_release() {
    log_info "Creating GitHub release..."
    
    if ! command -v gh &> /dev/null; then
        log_warning "GitHub CLI not found, skipping GitHub release creation"
        log_info "Please create the release manually at https://github.com/francknouama/go-starter/releases"
        return
    fi
    
    # Check if gh is authenticated
    if ! gh auth status &> /dev/null; then
        log_warning "GitHub CLI not authenticated, skipping GitHub release creation"
        log_info "Run 'gh auth login' to authenticate"
        return
    fi
    
    if [ "$DRY_RUN" = false ]; then
        RELEASE_NOTES="Release $VERSION

## What's Changed
- Updated version to $VERSION
- Bug fixes and improvements

## Installation
Download the appropriate binary for your platform from the assets below.

## Checksums
See checksums.txt for file verification."
        
        gh_args=("release" "create" "$VERSION" "--title" "Release $VERSION" "--notes" "$RELEASE_NOTES")
        
        if [ "$PRE_RELEASE" = true ]; then
            gh_args+=("--prerelease")
        fi
        
        # Add all files in release directory
        if [ -d "$RELEASE_DIR" ]; then
            for file in "$RELEASE_DIR"/*; do
                if [ -f "$file" ]; then
                    gh_args+=("$file")
                fi
            done
        fi
        
        gh "${gh_args[@]}"
    fi
    
    log_success "GitHub release created"
}

# Cleanup
cleanup() {
    log_info "Cleaning up..."
    
    if [ "$DRY_RUN" = false ]; then
        # Remove release directory if it exists
        if [ -d "$RELEASE_DIR" ]; then
            rm -rf "$RELEASE_DIR"
        fi
    fi
    
    log_success "Cleanup completed"
}

# Print release summary
print_summary() {
    echo
    log_success "=== Release Summary ==="
    echo
    log_info "Version: $VERSION"
    log_info "Pre-release: $PRE_RELEASE"
    log_info "Dry run: $DRY_RUN"
    echo
    
    if [ "$DRY_RUN" = false ]; then
        log_success "Release $VERSION has been successfully created!"
        echo
        log_info "Next steps:"
        echo "  1. Verify the release on GitHub"
        echo "  2. Test the release binaries"
        echo "  3. Update documentation if needed"
        echo "  4. Announce the release"
    else
        log_info "This was a dry run. No changes were made."
    fi
    
    echo
}

# Main release function
main() {
    echo "=== Go-Starter Release Script ==="
    echo
    
    parse_args "$@"
    
    if [ "$DRY_RUN" = true ]; then
        log_warning "DRY RUN MODE - No changes will be made"
        echo
    fi
    
    validate_version
    check_working_directory
    check_version_tag
    update_version
    run_pre_release_checks
    build_binaries
    generate_checksums
    create_git_tag
    create_github_release
    cleanup
    print_summary
}

# Handle script interruption
trap 'log_error "Release interrupted"; cleanup; exit 1' INT TERM

# Run main function
main "$@"