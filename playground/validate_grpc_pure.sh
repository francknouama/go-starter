#!/bin/bash

# gRPC-Pure Blueprint Validation Script
# This script validates the progress of gRPC-Pure template file completion

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
BLUEPRINT_DIR="$PROJECT_ROOT/blueprints/grpc-pure"
PLAYGROUND_DIR="$PROJECT_ROOT/playground"
TEST_PROJECT="grpc-validation-test"

echo -e "${BLUE}=== gRPC-Pure Blueprint Validation ===${NC}"
echo "Project Root: $PROJECT_ROOT"
echo "Blueprint Dir: $BLUEPRINT_DIR"
echo ""

# Function to count template files
count_template_files() {
    find "$BLUEPRINT_DIR" -name "*.tmpl" | wc -l | tr -d ' '
}

# Function to list missing template files
list_missing_templates() {
    echo -e "${YELLOW}Expected template files (from template.yaml):${NC}"
    
    # Extract expected template files from template.yaml
    grep -A 200 "^files:" "$BLUEPRINT_DIR/template.yaml" | \
    grep "source:" | \
    sed 's/.*source: *"//' | sed 's/".*//' | \
    sort | while read -r template_file; do
        template_path="$BLUEPRINT_DIR/$template_file"
        if [ -f "$template_path" ]; then
            echo -e "  ${GREEN}✓${NC} $template_file"
        else
            echo -e "  ${RED}✗${NC} $template_file"
        fi
    done
}

# Function to test project generation
test_generation() {
    echo -e "\n${BLUE}=== Testing Project Generation ===${NC}"
    
    cd "$PLAYGROUND_DIR"
    
    # Clean up any existing test project
    rm -rf "$TEST_PROJECT"
    
    echo "Running: go-starter new $TEST_PROJECT --type=grpc-pure --module=github.com/test/grpc --no-git --quiet"
    
    # Capture the output
    if OUTPUT=$("$PROJECT_ROOT/bin/go-starter" new "$TEST_PROJECT" --type=grpc-pure --module=github.com/test/grpc --no-git --quiet 2>&1); then
        echo -e "${GREEN}Generation completed successfully${NC}"
        
        # Extract files created count
        FILES_CREATED=$(echo "$OUTPUT" | grep "Files created:" | sed 's/.*Files created: *//' | sed 's/ .*//')
        echo "Files created: $FILES_CREATED"
        
        if [ -d "$TEST_PROJECT" ]; then
            ACTUAL_FILES=$(find "$TEST_PROJECT" -type f | wc -l | tr -d ' ')
            echo "Actual files found: $ACTUAL_FILES"
            
            if [ "$FILES_CREATED" -gt 0 ]; then
                echo -e "${GREEN}✓ Files were successfully generated${NC}"
                test_compilation
            else
                echo -e "${RED}✗ No files were generated despite success message${NC}"
            fi
        else
            echo -e "${RED}✗ Project directory was not created${NC}"
        fi
    else
        echo -e "${RED}Generation failed:${NC}"
        echo "$OUTPUT"
    fi
}

# Function to test compilation
test_compilation() {
    echo -e "\n${BLUE}=== Testing Compilation ===${NC}"
    
    if [ -d "$TEST_PROJECT" ]; then
        cd "$TEST_PROJECT"
        
        echo "Running: go mod tidy"
        if go mod tidy 2>/dev/null; then
            echo -e "${GREEN}✓ go mod tidy succeeded${NC}"
            
            echo "Running: go build ./..."
            if go build ./... 2>/dev/null; then
                echo -e "${GREEN}✓ Project compiles successfully${NC}"
            else
                echo -e "${RED}✗ Project compilation failed${NC}"
            fi
        else
            echo -e "${RED}✗ go mod tidy failed${NC}"
        fi
        
        cd "$PLAYGROUND_DIR"
    fi
}

# Function to run ATDD tests
run_atdd_tests() {
    echo -e "\n${BLUE}=== Running ATDD Tests ===${NC}"
    
    cd "$PROJECT_ROOT"
    
    echo "Running: go test ./tests/acceptance/blueprints/grpc-pure/ -v -short"
    if go test ./tests/acceptance/blueprints/grpc-pure/ -v -short; then
        echo -e "${GREEN}✓ ATDD tests passed${NC}"
    else
        echo -e "${RED}✗ ATDD tests failed${NC}"
    fi
}

# Main validation flow
main() {
    # Count current template files
    TEMPLATE_COUNT=$(count_template_files)
    echo -e "${BLUE}Current template files: $TEMPLATE_COUNT${NC}"
    
    # List missing templates
    echo ""
    list_missing_templates
    
    # Test generation
    test_generation
    
    # Run ATDD tests (with short flag to skip long-running tests)
    run_atdd_tests
    
    echo ""
    echo -e "${BLUE}=== Validation Summary ===${NC}"
    echo "Template files present: $TEMPLATE_COUNT"
    echo "Expected total files: 57"
    
    if [ "$TEMPLATE_COUNT" -eq 57 ]; then
        echo -e "${GREEN}✓ All template files are present${NC}"
    else
        MISSING=$((57 - TEMPLATE_COUNT))
        echo -e "${YELLOW}⚠ $MISSING template files still missing${NC}"
    fi
}

# Run main function
main "$@"