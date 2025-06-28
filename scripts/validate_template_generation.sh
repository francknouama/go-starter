#!/bin/bash

# Template Generation Validation Script
# This script validates that all template+logger combinations generate and compile successfully
# Used in CI/CD to prevent regressions

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

echo "🧪 Validating go-starter template generation..."
echo "Root directory: $ROOT_DIR"

# Ensure we're in the right directory
cd "$ROOT_DIR"

# Check if go-starter binary exists
if [ ! -f "./bin/go-starter" ]; then
    echo "❌ go-starter binary not found. Please run 'make build' first."
    exit 1
fi

# Create temporary directory for testing
TEST_DIR=$(mktemp -d)
echo "📁 Using temporary directory: $TEST_DIR"

# Cleanup function
cleanup() {
    echo "🧹 Cleaning up temporary directory..."
    rm -rf "$TEST_DIR"
}

# Set up trap to clean up on exit
trap cleanup EXIT

# Change to test directory
cd "$TEST_DIR"

# Template+logger combinations to test
declare -a combinations=(
    "web-api:slog:gin"
    "web-api:zap:gin"
    "web-api:logrus:gin" 
    "web-api:zerolog:gin"
    "cli:slog:cobra"
    "cli:zap:cobra"
    "cli:logrus:cobra"
    "cli:zerolog:cobra"
    "library:slog:"
    "library:zap:"
    "library:logrus:"
    "library:zerolog:"
    "lambda:slog:lambda"
    "lambda:zap:lambda"
    "lambda:logrus:lambda"
    "lambda:zerolog:lambda"
)

success_count=0
total_count=${#combinations[@]}
failed_combinations=()

echo "📋 Testing ${total_count} template+logger combinations..."
echo

for combo in "${combinations[@]}"; do
    IFS=':' read -r template logger framework <<< "$combo"
    project_name="validate-${template}-${logger}"
    
    echo "🔄 Testing ${template} + ${logger}..."
    
    # Build command arguments
    cmd_args=(
        "$project_name"
        "--type=$template"
        "--logger=$logger" 
        "--module=github.com/test/$project_name"
        "--go-version=1.21"
        "--no-git"
    )
    
    # Add framework if specified
    if [ -n "$framework" ]; then
        cmd_args+=("--framework=$framework")
    fi
    
    # Generate project with timeout - provide all inputs via echo
    if echo -e "$project_name\ngithub.com/test/$project_name" | timeout 30s "$ROOT_DIR/bin/go-starter" new "${cmd_args[@]}" > /dev/null 2>&1; then
        # Check if project compiles
        if cd "$project_name" && go build ./... > /dev/null 2>&1; then
            echo "  ✅ ${template}:${logger} - OK"
            success_count=$((success_count + 1))
        else
            echo "  ❌ ${template}:${logger} - Generated but failed to compile"
            failed_combinations+=("${template}:${logger} (compilation)")
            # Show compilation errors for debugging
            echo "     Compilation errors:"
            go build ./... 2>&1 | head -5 | sed 's/^/     /'
        fi
        cd "$TEST_DIR"
    else
        echo "  ❌ ${template}:${logger} - Failed to generate"
        failed_combinations+=("${template}:${logger} (generation)")
    fi
done

echo
echo "📊 Validation Results:"
echo "  ✅ Successful: ${success_count}/${total_count}"
echo "  ❌ Failed: $((total_count - success_count))/${total_count}"

if [ "$success_count" -eq "$total_count" ]; then
    echo
    echo "🎉 All template+logger combinations validated successfully!"
    echo "✨ Template generation is working correctly."
    exit 0
else
    echo
    echo "💥 Some combinations failed validation:"
    for failed in "${failed_combinations[@]}"; do
        echo "  - $failed"
    done
    echo
    echo "❌ Template generation validation failed."
    echo "Please fix the failing combinations before merging."
    exit 1
fi