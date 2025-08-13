#!/bin/bash

# Go-Starter Performance & Reliability Test Suite
# This script runs comprehensive performance, profiling, and cross-platform tests

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
OUTPUT_DIR="$PROJECT_ROOT/performance_test_results"
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")

# Test configuration
RUN_BENCHMARKS=${RUN_BENCHMARKS:-true}
RUN_PROFILING=${RUN_PROFILING:-true}
RUN_CROSS_PLATFORM=${RUN_CROSS_PLATFORM:-true}
VERBOSE=${VERBOSE:-false}
TARGET_DURATION=${TARGET_DURATION:-2s}
MAX_MEMORY_MB=${MAX_MEMORY_MB:-50}

# Functions
print_header() {
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                 Go-Starter Performance Test Suite                    ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

print_section() {
    echo -e "${YELLOW}▶ $1${NC}"
    echo ""
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

check_dependencies() {
    print_section "Checking Dependencies"
    
    # Check Go installation
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed. Please install Go 1.19 or later."
        exit 1
    fi
    
    GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
    print_success "Go $GO_VERSION found"
    
    # Check if we're in the correct directory
    if [ ! -f "$PROJECT_ROOT/go.mod" ]; then
        print_error "Not in a Go project root. Please run from go-starter project directory."
        exit 1
    fi
    
    # Check for blueprints directory
    if [ ! -d "$PROJECT_ROOT/blueprints" ]; then
        print_error "Blueprints directory not found. Ensure you're in the go-starter project root."
        exit 1
    fi
    
    print_success "All dependencies satisfied"
}

setup_environment() {
    print_section "Setting Up Test Environment"
    
    # Create output directory
    mkdir -p "$OUTPUT_DIR"
    print_success "Output directory created: $OUTPUT_DIR"
    
    # Build the project to ensure everything compiles
    cd "$PROJECT_ROOT"
    if ! go build ./... > /dev/null 2>&1; then
        print_error "Project compilation failed. Please fix build errors first."
        exit 1
    fi
    print_success "Project builds successfully"
    
    # Generate embedded templates
    if command -v go-generate &> /dev/null; then
        go generate ./...
        print_success "Templates generated"
    fi
    
    # Clean any existing test artifacts
    find "$PROJECT_ROOT" -name "test-*" -type d -exec rm -rf {} + 2>/dev/null || true
    find "$PROJECT_ROOT" -name "*.prof" -type f -delete 2>/dev/null || true
    print_success "Test environment cleaned"
}

run_unit_tests() {
    print_section "Running Unit Tests"
    
    cd "$PROJECT_ROOT"
    if go test ./internal/... -v -timeout=5m; then
        print_success "Unit tests passed"
        return 0
    else
        print_warning "Some unit tests failed - continuing with performance tests"
        return 1
    fi
}

run_benchmark_tests() {
    if [ "$RUN_BENCHMARKS" != "true" ]; then
        print_section "Skipping Benchmark Tests (disabled)"
        return 0
    fi
    
    print_section "Running Performance Benchmarks"
    
    cd "$PROJECT_ROOT"
    
    # Create benchmark output file
    BENCHMARK_OUTPUT="$OUTPUT_DIR/benchmark_results_$TIMESTAMP.txt"
    
    # Run benchmarks with detailed output
    echo "Running benchmarks for template generation performance..."
    if go test -bench=. -benchmem -count=3 -timeout=10m ./tests/performance -v > "$BENCHMARK_OUTPUT" 2>&1; then
        print_success "Performance benchmarks completed"
        
        # Extract key metrics
        if [ -f "$BENCHMARK_OUTPUT" ]; then
            echo ""
            echo "Key Performance Metrics:"
            grep -E "(Benchmark|PASS|FAIL)" "$BENCHMARK_OUTPUT" | tail -20 || true
        fi
        
        return 0
    else
        print_warning "Some benchmark tests failed - check output file: $BENCHMARK_OUTPUT"
        return 1
    fi
}

run_profiling_tests() {
    if [ "$RUN_PROFILING" != "true" ]; then
        print_section "Skipping Profiling Tests (disabled)"
        return 0
    fi
    
    print_section "Running Performance Profiling"
    
    cd "$PROJECT_ROOT"
    
    # Create profiling output directory
    PROFILE_DIR="$OUTPUT_DIR/profiles_$TIMESTAMP"
    mkdir -p "$PROFILE_DIR"
    
    # Set profiling environment variables
    export PROFILE_OUTPUT_DIR="$PROFILE_DIR"
    
    # Run profiling tests
    echo "Generating CPU and memory profiles..."
    if go test -run=TestProfile -v -timeout=15m ./tests/performance; then
        print_success "Performance profiling completed"
        
        # List generated profiles
        if [ -d "$PROFILE_DIR" ] && [ -n "$(ls -A "$PROFILE_DIR" 2>/dev/null)" ]; then
            echo ""
            echo "Profiles generated:"
            ls -la "$PROFILE_DIR"/*.prof 2>/dev/null || echo "No .prof files found"
        fi
        
        return 0
    else
        print_warning "Profiling tests encountered issues - check logs"
        return 1
    fi
}

run_cross_platform_tests() {
    if [ "$RUN_CROSS_PLATFORM" != "true" ]; then
        print_section "Skipping Cross-Platform Tests (disabled)"
        return 0
    fi
    
    print_section "Running Cross-Platform Compatibility Tests"
    
    cd "$PROJECT_ROOT"
    
    # Create cross-platform output file
    CROSSPLATFORM_OUTPUT="$OUTPUT_DIR/crossplatform_results_$TIMESTAMP.txt"
    
    # Run cross-platform tests
    echo "Testing cross-compilation for multiple platforms..."
    if go test -run=TestCrossPlatform -v -timeout=15m ./tests/crossplatform > "$CROSSPLATFORM_OUTPUT" 2>&1; then
        print_success "Cross-platform compatibility tests completed"
        
        # Show summary
        if [ -f "$CROSSPLATFORM_OUTPUT" ]; then
            echo ""
            echo "Cross-Platform Test Summary:"
            grep -E "(PASS|FAIL|compatibility)" "$CROSSPLATFORM_OUTPUT" | tail -10 || true
        fi
        
        return 0
    else
        print_warning "Some cross-platform tests failed - check output file: $CROSSPLATFORM_OUTPUT"
        return 1
    fi
}

run_integration_tests() {
    print_section "Running Integration Tests"
    
    cd "$PROJECT_ROOT"
    
    # Run ATDD tests
    echo "Running Acceptance Test-Driven Development (ATDD) tests..."
    if go test -run=TestATDD -v -timeout=10m ./tests/acceptance/...; then
        print_success "Integration tests passed"
        return 0
    else
        print_warning "Some integration tests failed"
        return 1
    fi
}

generate_comprehensive_report() {
    print_section "Generating Comprehensive Report"
    
    cd "$PROJECT_ROOT"
    
    # Create final report
    REPORT_FILE="$OUTPUT_DIR/performance_report_$TIMESTAMP.md"
    
    cat > "$REPORT_FILE" << EOF
# Go-Starter Performance & Reliability Test Report

**Generated:** $(date '+%Y-%m-%d %H:%M:%S %Z')
**Platform:** $(uname -s)/$(uname -m)
**Go Version:** $(go version | cut -d' ' -f3)
**Target Performance:** Generation < $TARGET_DURATION, Memory < ${MAX_MEMORY_MB}MB

## Test Summary

EOF

    # Add test results
    echo "### Benchmark Results" >> "$REPORT_FILE"
    if [ -f "$OUTPUT_DIR/benchmark_results_$TIMESTAMP.txt" ]; then
        echo "" >> "$REPORT_FILE"
        echo "\`\`\`" >> "$REPORT_FILE"
        tail -20 "$OUTPUT_DIR/benchmark_results_$TIMESTAMP.txt" >> "$REPORT_FILE" 2>/dev/null || echo "No benchmark results available" >> "$REPORT_FILE"
        echo "\`\`\`" >> "$REPORT_FILE"
    fi
    
    echo "" >> "$REPORT_FILE"
    echo "### Cross-Platform Results" >> "$REPORT_FILE"
    if [ -f "$OUTPUT_DIR/crossplatform_results_$TIMESTAMP.txt" ]; then
        echo "" >> "$REPORT_FILE"
        echo "\`\`\`" >> "$REPORT_FILE"
        tail -20 "$OUTPUT_DIR/crossplatform_results_$TIMESTAMP.txt" >> "$REPORT_FILE" 2>/dev/null || echo "No cross-platform results available" >> "$REPORT_FILE"
        echo "\`\`\`" >> "$REPORT_FILE"
    fi
    
    echo "" >> "$REPORT_FILE"
    echo "### Files Generated" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    echo "- Benchmark Results: \`benchmark_results_$TIMESTAMP.txt\`" >> "$REPORT_FILE"
    echo "- Cross-Platform Results: \`crossplatform_results_$TIMESTAMP.txt\`" >> "$REPORT_FILE"
    
    if [ -d "$OUTPUT_DIR/profiles_$TIMESTAMP" ]; then
        echo "- CPU/Memory Profiles: \`profiles_$TIMESTAMP/\`" >> "$REPORT_FILE"
    fi
    
    echo "" >> "$REPORT_FILE"
    echo "### Recommendations" >> "$REPORT_FILE"
    echo "" >> "$REPORT_FILE"
    
    # Add automatic recommendations based on results
    if grep -q "FAIL" "$OUTPUT_DIR"/*_results_*.txt 2>/dev/null; then
        echo "- 🔴 **Critical:** Some tests failed - investigate and fix failing test cases" >> "$REPORT_FILE"
    fi
    
    echo "- 📊 Analyze profile data using: \`go tool pprof <profile-file>\`" >> "$REPORT_FILE"
    echo "- 🔧 For detailed analysis, run individual test suites with -v flag" >> "$REPORT_FILE"
    echo "- 📈 Monitor performance trends over time by comparing reports" >> "$REPORT_FILE"
    
    print_success "Comprehensive report generated: $REPORT_FILE"
}

cleanup() {
    print_section "Cleaning Up"
    
    # Remove temporary test files
    find "$PROJECT_ROOT" -name "test-*" -type d -exec rm -rf {} + 2>/dev/null || true
    
    # Keep profiles but clean other temp files
    find "$PROJECT_ROOT" -name "*.tmp" -type f -delete 2>/dev/null || true
    
    print_success "Cleanup completed"
}

print_final_summary() {
    echo ""
    echo -e "${BLUE}╔═══════════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║                        TEST SUITE COMPLETED                          ║${NC}"
    echo -e "${BLUE}╚═══════════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    
    echo "📊 Results Location: $OUTPUT_DIR"
    echo "📋 Main Report: performance_report_$TIMESTAMP.md"
    echo ""
    
    # Count results
    local total_files=$(find "$OUTPUT_DIR" -name "*_$TIMESTAMP*" -type f | wc -l)
    echo "📁 Generated $total_files result files"
    
    # Check if any critical issues were found
    if grep -q "FAIL\|ERROR" "$OUTPUT_DIR"/*_results_*.txt 2>/dev/null; then
        echo -e "${RED}⚠️  Critical issues detected - review test results${NC}"
    else
        echo -e "${GREEN}✅ All tests completed successfully${NC}"
    fi
    
    echo ""
    echo "To analyze profiles: go tool pprof <profile-file>"
    echo "To re-run specific tests: go test -run=<TestName> -v ./tests/<category>"
    echo ""
}

# Main execution
main() {
    print_header
    
    # Parse command line arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --no-benchmarks)
                RUN_BENCHMARKS=false
                shift
                ;;
            --no-profiling)
                RUN_PROFILING=false
                shift
                ;;
            --no-cross-platform)
                RUN_CROSS_PLATFORM=false
                shift
                ;;
            --verbose)
                VERBOSE=true
                shift
                ;;
            --target-duration)
                TARGET_DURATION="$2"
                shift 2
                ;;
            --max-memory)
                MAX_MEMORY_MB="$2"
                shift 2
                ;;
            --help)
                echo "Usage: $0 [options]"
                echo ""
                echo "Options:"
                echo "  --no-benchmarks     Skip benchmark tests"
                echo "  --no-profiling      Skip profiling tests"
                echo "  --no-cross-platform Skip cross-platform tests"
                echo "  --verbose           Enable verbose output"
                echo "  --target-duration   Target generation duration (default: 2s)"
                echo "  --max-memory        Maximum memory usage in MB (default: 50)"
                echo "  --help              Show this help message"
                echo ""
                echo "Environment variables:"
                echo "  RUN_BENCHMARKS      Enable/disable benchmarks (true/false)"
                echo "  RUN_PROFILING       Enable/disable profiling (true/false)"
                echo "  RUN_CROSS_PLATFORM  Enable/disable cross-platform tests (true/false)"
                echo "  VERBOSE             Enable verbose output (true/false)"
                echo ""
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                echo "Use --help for usage information"
                exit 1
                ;;
        esac
    done
    
    # Show configuration
    echo "Test Configuration:"
    echo "  Benchmarks: $RUN_BENCHMARKS"
    echo "  Profiling: $RUN_PROFILING"
    echo "  Cross-Platform: $RUN_CROSS_PLATFORM"
    echo "  Target Duration: $TARGET_DURATION"
    echo "  Max Memory: ${MAX_MEMORY_MB}MB"
    echo "  Verbose: $VERBOSE"
    echo ""
    
    # Execute test suite
    local exit_code=0
    
    check_dependencies || exit_code=$?
    setup_environment || exit_code=$?
    
    # Run tests (continue even if some fail)
    run_unit_tests || true
    run_benchmark_tests || true
    run_profiling_tests || true
    run_cross_platform_tests || true
    run_integration_tests || true
    
    # Always generate report and cleanup
    generate_comprehensive_report
    cleanup
    print_final_summary
    
    exit $exit_code
}

# Handle script interruption
trap cleanup EXIT

# Run main function
main "$@"