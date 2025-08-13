#!/bin/bash

# gRPC-Pure Template Validation Script
# Provides continuous validation and rapid feedback for template completion

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMP_DIR=""
VALIDATION_LOG=""

# Color output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Cleanup function
cleanup() {
    if [[ -n "$TEMP_DIR" && -d "$TEMP_DIR" ]]; then
        rm -rf "$TEMP_DIR"
    fi
}
trap cleanup EXIT

log() {
    local msg="${BLUE}[$(date '+%H:%M:%S')]${NC} $1"
    echo -e "$msg"
    if [[ -n "$VALIDATION_LOG" ]]; then
        echo "$1" >> "$VALIDATION_LOG"
    fi
}

success() {
    local msg="${GREEN}✅ $1${NC}"
    echo -e "$msg"
    if [[ -n "$VALIDATION_LOG" ]]; then
        echo "✅ $1" >> "$VALIDATION_LOG"
    fi
}

warning() {
    local msg="${YELLOW}⚠️  $1${NC}"
    echo -e "$msg"
    if [[ -n "$VALIDATION_LOG" ]]; then
        echo "⚠️  $1" >> "$VALIDATION_LOG"
    fi
}

error() {
    local msg="${RED}❌ $1${NC}"
    echo -e "$msg"
    if [[ -n "$VALIDATION_LOG" ]]; then
        echo "❌ $1" >> "$VALIDATION_LOG"
    fi
}

# Initialize
init_validation() {
    log "Initializing gRPC-Pure template validation..."
    
    TEMP_DIR=$(mktemp -d -t grpc_pure_validation_XXXXXX)
    VALIDATION_LOG="$PROJECT_ROOT/grpc_pure_validation.log"
    
    log "Temp directory: $TEMP_DIR"
    log "Validation log: $VALIDATION_LOG"
    
    # Build go-starter
    cd "$PROJECT_ROOT"
    log "Building go-starter..."
    go build -o bin/go-starter main.go
    success "go-starter built successfully"
}

# Count existing vs required templates
check_template_count() {
    log "Checking template count..."
    
    local blueprint_dir="$PROJECT_ROOT/blueprints/grpc-pure"
    local existing_count=$(find "$blueprint_dir" -name "*.tmpl" | wc -l | tr -d ' ')
    local required_count=$(grep -c "source:" "$blueprint_dir/template.yaml")
    local missing_count=$((required_count - existing_count))
    
    log "Template count status:"
    log "  - Required: $required_count"
    log "  - Existing: $existing_count"
    log "  - Missing: $missing_count"
    
    if [[ $missing_count -eq 0 ]]; then
        success "All templates are present!"
        return 0
    else
        warning "$missing_count templates are missing"
        return 1
    fi
}

# Validate template syntax
validate_template_syntax() {
    log "Validating template syntax..."
    
    local blueprint_dir="$PROJECT_ROOT/blueprints/grpc-pure"
    local errors=0
    
    while IFS= read -r -d '' template_file; do
        local filename=$(basename "$template_file")
        log "Checking syntax: $filename"
        
        # Basic template syntax check
        if ! go run -tags=dev "$PROJECT_ROOT/scripts/template_validator.go" "$template_file"; then
            error "Template syntax error in $filename"
            ((errors++))
        else
            success "Template syntax OK: $filename"
        fi
    done < <(find "$blueprint_dir" -name "*.tmpl" -print0)
    
    if [[ $errors -eq 0 ]]; then
        success "All template syntax validation passed"
        return 0
    else
        error "$errors template syntax errors found"
        return 1
    fi
}

# Test basic generation
test_basic_generation() {
    log "Testing basic gRPC-Pure generation..."
    
    cd "$TEMP_DIR"
    
    # Test dry run first
    log "Running dry-run generation..."
    if "$PROJECT_ROOT/bin/go-starter" new grpc-basic-test \
        --type=grpc-pure \
        --architecture=microservice \
        --module=github.com/example/grpc-basic-test \
        --no-git \
        --logger=slog \
        --dry-run > generation_output.log 2>&1; then
        success "Dry-run generation successful"
    else
        error "Dry-run generation failed"
        cat generation_output.log
        return 1
    fi
    
    # Test actual generation
    log "Running actual generation..."
    if "$PROJECT_ROOT/bin/go-starter" new grpc-basic-test \
        --type=grpc-pure \
        --architecture=microservice \
        --module=github.com/example/grpc-basic-test \
        --no-git \
        --logger=slog > generation_output.log 2>&1; then
        success "Basic generation successful"
    else
        error "Basic generation failed"
        cat generation_output.log
        return 1
    fi
    
    # Count generated files
    local generated_count=$(find grpc-basic-test -type f | wc -l | tr -d ' ')
    log "Generated $generated_count files"
    
    return 0
}

# Test compilation
test_compilation() {
    log "Testing project compilation..."
    
    local project_dir="$TEMP_DIR/grpc-basic-test"
    
    if [[ ! -d "$project_dir" ]]; then
        warning "Project directory not found, skipping compilation test"
        return 1
    fi
    
    cd "$project_dir"
    
    # Run go mod tidy
    log "Running go mod tidy..."
    if go mod tidy > ../compilation.log 2>&1; then
        success "go mod tidy successful"
    else
        error "go mod tidy failed"
        cat ../compilation.log
        return 1
    fi
    
    # Test compilation
    log "Testing compilation..."
    if go build ./... > ../compilation.log 2>&1; then
        success "Project compilation successful"
        return 0
    else
        error "Project compilation failed"
        cat ../compilation.log
        return 1
    fi
}

# Test protobuf generation
test_protobuf_generation() {
    log "Testing protobuf file validation..."
    
    local project_dir="$TEMP_DIR/grpc-basic-test"
    
    if [[ ! -d "$project_dir" ]]; then
        warning "Project directory not found, skipping protobuf test"
        return 1
    fi
    
    cd "$project_dir"
    
    # Check for required proto files
    local proto_files=(
        "proto/grpc-basic-test/v1/service.proto"
        "proto/health/v1/health.proto"
        "proto/common/v1/common.proto"
    )
    
    for proto_file in "${proto_files[@]}"; do
        if [[ -f "$proto_file" ]]; then
            success "Proto file exists: $proto_file"
        else
            error "Proto file missing: $proto_file"
            return 1
        fi
    done
    
    # Check for buf configuration
    if [[ -f "buf.yaml" && -f "buf.gen.yaml" ]]; then
        success "Buf configuration files present"
    else
        warning "Buf configuration files missing"
        return 1
    fi
    
    return 0
}

# Run ATDD tests
run_atdd_tests() {
    log "Running gRPC-Pure ATDD tests..."
    
    cd "$PROJECT_ROOT"
    
    if go test -v ./tests/acceptance/blueprints/grpc-pure/... > atdd_test_output.log 2>&1; then
        success "ATDD tests passed"
        return 0
    else
        error "ATDD tests failed"
        cat atdd_test_output.log
        return 1
    fi
}

# Generate validation report
generate_report() {
    local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
    local report_file="$PROJECT_ROOT/grpc_pure_validation_report.md"
    
    log "Generating validation report..."
    
    cat > "$report_file" << EOF
# gRPC-Pure Template Validation Report

**Generated:** $timestamp

## Summary

EOF

    # Count current templates
    local blueprint_dir="$PROJECT_ROOT/blueprints/grpc-pure"
    local existing_count=$(find "$blueprint_dir" -name "*.tmpl" | wc -l | tr -d ' ')
    local required_count=$(grep -c "source:" "$blueprint_dir/template.yaml")
    local missing_count=$((required_count - existing_count))
    local completion_percent=$((existing_count * 100 / required_count))
    
    cat >> "$report_file" << EOF
- **Template Completion:** $existing_count/$required_count ($completion_percent%)
- **Missing Templates:** $missing_count
- **Validation Status:** $(if [[ $missing_count -eq 0 ]]; then echo "✅ Complete"; else echo "🔄 In Progress"; fi)

## Template Status

### Existing Templates ($existing_count)
EOF

    # List existing templates
    find "$blueprint_dir" -name "*.tmpl" | sort | while read -r template; do
        local rel_path=${template#$blueprint_dir/}
        echo "- ✅ $rel_path" >> "$report_file"
    done
    
    if [[ $missing_count -gt 0 ]]; then
        cat >> "$report_file" << EOF

### Missing Templates ($missing_count)
EOF
        
        # List missing templates
        cd "$blueprint_dir"
        grep "source:" template.yaml | sed 's/.*source: "//' | sed 's/".*//' | while read -r file; do
            if [[ ! -f "$file" ]]; then
                echo "- ❌ $file" >> "$report_file"
            fi
        done
    fi
    
    cat >> "$report_file" << EOF

## Recent Validation Log

\`\`\`
EOF
    
    # Add last 50 lines of validation log
    if [[ -f "$VALIDATION_LOG" ]]; then
        tail -50 "$VALIDATION_LOG" >> "$report_file"
    fi
    
    cat >> "$report_file" << EOF
\`\`\`

## Next Steps

EOF

    if [[ $missing_count -gt 0 ]]; then
        cat >> "$report_file" << EOF
1. **Priority:** Complete missing template files
2. **Validation:** Run validation after each template completion
3. **Testing:** Ensure all ATDD scenarios pass
4. **Compilation:** Verify generated projects compile successfully

### Immediate Action Items

- Focus on core infrastructure templates first (Dockerfile, Makefile, buf configs)
- Validate template syntax as files are added
- Test incremental generation with each completed phase
EOF
    else
        cat >> "$report_file" << EOF
1. **Run comprehensive ATDD tests**
2. **Test all feature combinations**
3. **Validate compilation across scenarios**
4. **Prepare for GraphQL API blueprint framework**

### All Templates Complete! 🎉
EOF
    fi
    
    success "Validation report generated: $report_file"
}

# Main validation flow
main() {
    log "=== gRPC-Pure Template Validation Started ==="
    
    init_validation
    
    local validation_passed=true
    
    # Phase 1: Template inventory and syntax
    if ! check_template_count; then
        validation_passed=false
    fi
    
    # Validate existing templates only (skip syntax validation for missing files)
    log "Validating existing template syntax..."
    local syntax_errors=0
    local blueprint_dir="$PROJECT_ROOT/blueprints/grpc-pure"
    
    while IFS= read -r -d '' template_file; do
        local filename=$(basename "$template_file")
        log "Checking template: $filename"
        
        # Simple syntax check - look for obvious template errors
        if grep -q "{{.*}}" "$template_file"; then
            success "Template markers found in $filename"
        else
            warning "No template markers in $filename (may be static file)"
        fi
        
        # Check for unterminated template actions
        if grep -q "{{[^}]*$" "$template_file"; then
            error "Unterminated template action in $filename"
            ((syntax_errors++))
        fi
        
    done < <(find "$blueprint_dir" -name "*.tmpl" -print0)
    
    if [[ $syntax_errors -gt 0 ]]; then
        validation_passed=false
    fi
    
    # Phase 2: Generation testing (only if we have core templates)
    if [[ -f "$blueprint_dir/cmd/server/main.go.tmpl" && -f "$blueprint_dir/internal/config/config.go.tmpl" ]]; then
        if ! test_basic_generation; then
            validation_passed=false
        fi
        
        if ! test_compilation; then
            validation_passed=false
        fi
        
        if ! test_protobuf_generation; then
            validation_passed=false
        fi
    else
        warning "Core templates missing, skipping generation tests"
    fi
    
    # Phase 3: ATDD testing (only if templates are mostly complete)
    local existing_count=$(find "$blueprint_dir" -name "*.tmpl" | wc -l | tr -d ' ')
    local required_count=$(grep -c "source:" "$blueprint_dir/template.yaml")
    if [[ $existing_count -gt $((required_count * 80 / 100)) ]]; then
        if ! run_atdd_tests; then
            validation_passed=false
        fi
    else
        warning "Less than 80% templates complete, skipping ATDD tests"
    fi
    
    # Generate report
    generate_report
    
    if [[ $validation_passed == true ]]; then
        success "=== All validations passed! ==="
        log "Ready for next phase of development"
        exit 0
    else
        error "=== Some validations failed ==="
        log "Check validation log for details"
        exit 1
    fi
}

# Run main function
main "$@"