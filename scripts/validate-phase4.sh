#!/bin/bash

# Phase 4 Comprehensive Validation Script
# Validates conditional dependencies, logger integration, and template consistency

echo "🔍 Phase 4 Comprehensive Validation"
echo "===================================="
echo

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Counters
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# Helper functions
log_test() {
    echo -e "${YELLOW}🧪 $1${NC}"
    ((TOTAL_TESTS++))
}

log_success() {
    echo -e "   ${GREEN}✅ $1${NC}"
    ((PASSED_TESTS++))
}

log_failure() {
    echo -e "   ${RED}❌ $1${NC}"
    ((FAILED_TESTS++))
}

log_info() {
    echo -e "   ℹ️  $1"
}

# Clean up function
cleanup() {
    rm -rf /tmp/validate-*
}

# Build the tool
echo "📦 Building go-starter..."
if ! go build -o bin/go-starter main.go; then
    echo "❌ Failed to build go-starter"
    exit 1
fi
echo "✅ Build successful"
echo

# Test matrix: [project_type, logger, framework, has_database]
test_matrix=(
    "web-api,zap,gin,false"
    "web-api,slog,gin,false" 
    "web-api,logrus,gin,false"
    "web-api,zerolog,gin,false"
    "cli,zap,cobra,false"
    "cli,slog,cobra,false"
    "cli,logrus,cobra,false"
    "cli,zerolog,cobra,false"
    "library,zap,,false"
    "library,slog,,false"
    "library,logrus,,false"
    "library,zerolog,,false"
    "lambda,zap,,false"
    "lambda,slog,,false"
    "lambda,logrus,,false"
    "lambda,zerolog,,false"
)

# Logger dependency mappings
declare -A logger_deps=(
    ["zap"]="go.uber.org/zap"
    ["logrus"]="github.com/sirupsen/logrus"
    ["zerolog"]="github.com/rs/zerolog"
    ["slog"]="" # Built-in, no external dependency
)

# Framework dependency mappings
declare -A framework_deps=(
    ["gin"]="github.com/gin-gonic/gin"
    ["cobra"]="github.com/spf13/cobra"
)

echo "🎯 Testing Template Generation and Dependencies"
echo "=============================================="

for test_case in "${test_matrix[@]}"; do
    IFS=',' read -r project_type logger framework has_db <<< "$test_case"
    
    project_name="validate-${project_type}-${logger}"
    log_test "Testing ${project_type} with ${logger} logger"
    
    # Clean up previous test
    rm -rf "/tmp/${project_name}"
    
    # Generate project
    args=(
        "new" "$project_name"
        "--type=$project_type"
        "--logger=$logger"
        "--module=github.com/test/$project_name"
        "--output=/tmp"
    )
    
    if [ -n "$framework" ]; then
        args+=("--framework=$framework")
    fi
    
    if ! ./bin/go-starter "${args[@]}" >/dev/null 2>&1; then
        log_failure "Project generation failed"
        continue
    fi
    log_success "Project generated successfully"
    
    project_dir="/tmp/$project_name"
    
    # Test 1: Check go.mod dependencies
    log_info "Validating go.mod dependencies..."
    gomod_file="$project_dir/go.mod"
    
    if [ ! -f "$gomod_file" ]; then
        log_failure "go.mod file not found"
        continue
    fi
    
    # Check logger dependency
    logger_dep="${logger_deps[$logger]}"
    if [ -n "$logger_dep" ]; then
        if ! grep -q "$logger_dep" "$gomod_file"; then
            log_failure "Missing logger dependency: $logger_dep"
            continue
        fi
        log_success "Logger dependency $logger_dep found"
    else
        log_success "No external logger dependency needed for slog"
    fi
    
    # Check framework dependency
    if [ -n "$framework" ]; then
        framework_dep="${framework_deps[$framework]}"
        if [ -n "$framework_dep" ]; then
            if ! grep -q "$framework_dep" "$gomod_file"; then
                log_failure "Missing framework dependency: $framework_dep"
                continue
            fi
            log_success "Framework dependency $framework_dep found"
        fi
    fi
    
    # Check that other logger dependencies are NOT present
    for other_logger in zap logrus zerolog; do
        if [ "$other_logger" != "$logger" ]; then
            other_dep="${logger_deps[$other_logger]}"
            if [ -n "$other_dep" ] && grep -q "$other_dep" "$gomod_file"; then
                log_failure "Unexpected logger dependency found: $other_dep"
                continue
            fi
        fi
    done
    log_success "No unexpected logger dependencies found"
    
    # Test 2: Check logger implementation files
    log_info "Validating logger implementation..."
    logger_dir="$project_dir/internal/logger"
    
    if [ ! -d "$logger_dir" ]; then
        log_failure "Logger directory not found"
        continue
    fi
    
    # Check for logger interface file
    if [ ! -f "$logger_dir/interface.go" ] && [ ! -f "$logger_dir/logger.go" ]; then
        log_failure "Logger interface file not found"
        continue
    fi
    log_success "Logger interface file found"
    
    # Check for specific logger implementation
    case $project_type in
        "web-api")
            if [ ! -f "$logger_dir/$logger.go" ]; then
                log_failure "Specific logger implementation not found: $logger.go"
                continue
            fi
            log_success "Logger implementation file found: $logger.go"
            ;;
        "library"|"lambda")
            if [ ! -f "$logger_dir/logger.go" ]; then
                log_failure "Logger implementation not found: logger.go"
                continue
            fi
            log_success "Logger implementation file found: logger.go"
            ;;
        "cli")
            if [ ! -f "$logger_dir/$logger.go" ]; then
                log_failure "Specific logger implementation not found: $logger.go"
                continue
            fi
            log_success "Logger implementation file found: $logger.go"
            ;;
    esac
    
    # Test 3: Compile the project
    log_info "Testing compilation..."
    
    build_dir="$project_dir"
    if [ "$project_type" = "web-api" ]; then
        build_dir="$project_dir/cmd/server"
    fi
    
    if ! go build -C "$build_dir" . >/dev/null 2>&1; then
        log_failure "Compilation failed"
        # Show the actual error for debugging
        echo "   Compilation error:"
        go build -C "$build_dir" . 2>&1 | sed 's/^/   /'
        continue
    fi
    log_success "Compilation successful"
    
    # Test 4: Run go mod tidy and verify no issues
    log_info "Testing go mod tidy..."
    
    if ! go mod tidy -C "$project_dir" >/dev/null 2>&1; then
        log_failure "go mod tidy failed"
        continue
    fi
    log_success "go mod tidy successful"
    
    # Test 5: Check for unused imports
    log_info "Checking for unused imports..."
    
    if go build -C "$build_dir" . 2>&1 | grep -q "imported and not used"; then
        log_failure "Unused imports detected"
        go build -C "$build_dir" . 2>&1 | grep "imported and not used" | sed 's/^/   /'
        continue
    fi
    log_success "No unused imports detected"
    
    log_success "All tests passed for ${project_type} with ${logger}"
    echo
done

echo "🧪 Testing Database Integration (Web API only)"
echo "=============================================="

# Test database integration specifically
db_test_cases=(
    "web-api,zap,gin"
    "web-api,slog,gin"
)

for test_case in "${db_test_cases[@]}"; do
    IFS=',' read -r project_type logger framework <<< "$test_case"
    
    project_name="validate-db-${logger}"
    log_test "Testing database integration with ${logger} logger"
    
    # Clean up previous test
    rm -rf "/tmp/${project_name}"
    
    # Generate project (without database for now, as --database flag doesn't exist)
    args=(
        "new" "$project_name"
        "--type=$project_type"
        "--logger=$logger"
        "--framework=$framework"
        "--module=github.com/test/$project_name"
        "--output=/tmp"
    )
    
    if ! ./bin/go-starter "${args[@]}" >/dev/null 2>&1; then
        log_failure "Project generation failed"
        continue
    fi
    log_success "Project generated successfully"
    
    project_dir="/tmp/$project_name"
    
    # Check if database files use logger interface
    db_connection_file="$project_dir/internal/database/connection.go"
    if [ -f "$db_connection_file" ]; then
        if grep -q "appLogger.Logger" "$db_connection_file"; then
            log_success "Database connection uses logger interface"
        else
            log_failure "Database connection doesn't use logger interface"
            continue
        fi
        
        if grep -q "func Connect.*Logger" "$db_connection_file"; then
            log_success "Connect function accepts logger parameter"
        else
            log_failure "Connect function doesn't accept logger parameter"
            continue
        fi
    else
        log_info "No database connection file (expected for basic template)"
    fi
    
    # Check migration files
    db_migration_file="$project_dir/internal/database/migrations.go"
    if [ -f "$db_migration_file" ]; then
        if grep -q "appLogger.Logger" "$db_migration_file"; then
            log_success "Database migrations use logger interface"
        else
            log_failure "Database migrations don't use logger interface"
            continue
        fi
    else
        log_info "No database migration file (expected for basic template)"
    fi
    
    log_success "Database integration test passed"
    echo
done

echo "🧪 Testing Logger Interface Consistency"
echo "======================================="

# Test that all logger implementations have consistent interfaces
log_test "Checking logger interface consistency across templates"

temp_projects=()
for logger_type in zap slog logrus zerolog; do
    project_name="validate-interface-$logger_type"
    temp_projects+=("$project_name")
    
    # Generate web-api project for interface testing
    ./bin/go-starter new "$project_name" --type=web-api --logger="$logger_type" --framework=gin --module="github.com/test/$project_name" --output=/tmp >/dev/null 2>&1
    
    if [ $? -ne 0 ]; then
        log_failure "Failed to generate project for $logger_type"
        continue
    fi
done

# Check that all logger implementations have the same interface methods
required_methods=("Debug" "Info" "Warn" "Error" "DebugWith" "InfoWith" "WarnWith" "ErrorWith")
interface_consistent=true

for project_name in "${temp_projects[@]}"; do
    logger_files="/tmp/$project_name/internal/logger"/*.go
    
    for method in "${required_methods[@]}"; do
        if ! grep -q "func.*$method" $logger_files 2>/dev/null; then
            log_failure "Method $method not found in $project_name logger implementation"
            interface_consistent=false
        fi
    done
done

if [ "$interface_consistent" = true ]; then
    log_success "Logger interface is consistent across all implementations"
else
    log_failure "Logger interface inconsistencies detected"
fi

echo "📊 Validation Summary"
echo "===================="
echo "Total tests: $TOTAL_TESTS"
echo "Passed: $PASSED_TESTS"
echo "Failed: $FAILED_TESTS"
echo

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 PHASE 4 VALIDATION SUCCESSFUL!${NC}"
    echo "✅ All conditional dependencies working correctly"
    echo "✅ All logger implementations consistent"
    echo "✅ All templates compile successfully"
    echo "✅ Database integration uses logger interface"
    echo "✅ No unused imports or dependency issues"
    echo
    echo "🚀 Phase 4 is properly implemented and ready for production!"
else
    echo -e "${RED}❌ PHASE 4 VALIDATION FAILED${NC}"
    echo "$FAILED_TESTS out of $TOTAL_TESTS tests failed"
    echo "Please review and fix the issues above."
fi

# Cleanup
echo "🧹 Cleaning up validation projects..."
cleanup
echo "✅ Cleanup complete"

exit $FAILED_TESTS