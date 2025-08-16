#!/bin/bash

# End-to-End Integration Test for go-starter Web UI
# Tests the complete backend API and frontend integration

set -e

echo "🚀 Starting go-starter Web UI Integration Tests"
echo "=============================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test configuration
WEB_SERVER_PORT=${WEB_SERVER_PORT:-8080}
BASE_URL="http://localhost:$WEB_SERVER_PORT"
API_URL="$BASE_URL/api/v1"

# Test status tracking
TESTS_PASSED=0
TESTS_FAILED=0

# Helper functions
print_success() {
    echo -e "${GREEN}✅ $1${NC}"
    ((TESTS_PASSED++))
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
    ((TESTS_FAILED++))
}

print_info() {
    echo -e "${YELLOW}ℹ️  $1${NC}"
}

# Check if server is running
check_server() {
    print_info "Checking if web server is running on port $WEB_SERVER_PORT..."
    
    if curl -s "$BASE_URL/api/v1/health/simple" > /dev/null; then
        print_success "Web server is running"
        return 0
    else
        print_error "Web server is not running. Please start it first with: make dev"
        return 1
    fi
}

# Test health endpoints
test_health_endpoints() {
    print_info "Testing health endpoints..."
    
    # Simple health check
    if curl -s "$API_URL/health/simple" | grep -q "ok"; then
        print_success "Simple health check passed"
    else
        print_error "Simple health check failed"
    fi
    
    # Full health check
    if curl -s "$API_URL/health" | grep -q "status"; then
        print_success "Full health check passed"
    else
        print_error "Full health check failed"
    fi
    
    # Metrics endpoint
    if curl -s "$API_URL/metrics" > /dev/null; then
        print_success "Metrics endpoint accessible"
    else
        print_error "Metrics endpoint failed"
    fi
}

# Test configuration endpoints
test_config_endpoints() {
    print_info "Testing configuration endpoints..."
    
    # Default config
    if curl -s "$API_URL/config" | grep -q "projectType"; then
        print_success "Default config endpoint passed"
    else
        print_error "Default config endpoint failed"
    fi
    
    # Frameworks
    if curl -s "$API_URL/config/frameworks" | grep -q "gin"; then
        print_success "Frameworks endpoint passed"
    else
        print_error "Frameworks endpoint failed"
    fi
    
    # Architectures
    if curl -s "$API_URL/config/architectures" | grep -q "standard"; then
        print_success "Architectures endpoint passed"
    else
        print_error "Architectures endpoint failed"
    fi
}

# Test blueprint endpoints
test_blueprint_endpoints() {
    print_info "Testing blueprint endpoints..."
    
    # List blueprints
    if curl -s "$API_URL/blueprints" | grep -q "web-api"; then
        print_success "Blueprints list endpoint passed"
    else
        print_error "Blueprints list endpoint failed"
    fi
    
    # Get specific blueprint
    if curl -s "$API_URL/blueprints/web-api" | grep -q "name"; then
        print_success "Specific blueprint endpoint passed"
    else
        print_error "Specific blueprint endpoint failed"
    fi
}

# Test project preview
test_project_preview() {
    print_info "Testing project preview..."
    
    local preview_data='{
        "projectName": "test-api",
        "moduleUrl": "github.com/test/test-api",
        "projectType": "web-api",
        "architecture": "standard",
        "framework": "gin",
        "logger": "slog",
        "goVersion": "1.21"
    }'
    
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$preview_data" \
        "$API_URL/preview")
    
    if echo "$response" | grep -q "fileStructure"; then
        print_success "Project preview generation passed"
    else
        print_error "Project preview generation failed"
        echo "Response: $response"
    fi
}

# Test project generation
test_project_generation() {
    print_info "Testing project generation..."
    
    local generation_data='{
        "projectName": "test-cli",
        "moduleUrl": "github.com/test/test-cli",
        "projectType": "cli",
        "architecture": "standard",
        "framework": "cobra",
        "logger": "slog",
        "goVersion": "1.21",
        "outputFormat": "zip",
        "includeTests": true,
        "includeDocs": true
    }'
    
    local response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "$generation_data" \
        "$API_URL/generate")
    
    if echo "$response" | grep -q "projectId"; then
        print_success "Project generation passed"
    else
        print_error "Project generation failed"
        echo "Response: $response"
    fi
}

# Test WebSocket endpoint
test_websocket() {
    print_info "Testing WebSocket endpoint..."
    
    # Check if WebSocket endpoint is accessible
    local ws_url="ws://localhost:$WEB_SERVER_PORT/api/v1/ws"
    
    # Simple WebSocket test using curl to check if the endpoint exists
    if curl -s -H "Upgrade: websocket" -H "Connection: Upgrade" "$API_URL/ws" 2>&1 | grep -q "101\|400\|426"; then
        print_success "WebSocket endpoint is accessible"
    else
        print_error "WebSocket endpoint failed"
    fi
}

# Test frontend static files
test_frontend_files() {
    print_info "Testing frontend static files..."
    
    # Test main HTML file
    if curl -s "$BASE_URL/" | grep -q "go-starter"; then
        print_success "Frontend HTML served correctly"
    else
        print_error "Frontend HTML not accessible"
    fi
    
    # Test if React app loads
    if curl -s "$BASE_URL/" | grep -q "div id=\"root\""; then
        print_success "React app structure found"
    else
        print_error "React app structure not found"
    fi
}

# Test CORS headers
test_cors() {
    print_info "Testing CORS headers..."
    
    local cors_response=$(curl -s -H "Origin: http://localhost:3000" \
        -H "Access-Control-Request-Method: POST" \
        -H "Access-Control-Request-Headers: Content-Type" \
        -X OPTIONS "$API_URL/blueprints")
    
    if echo "$cors_response" | grep -q "Access-Control-Allow"; then
        print_success "CORS headers present"
    else
        print_error "CORS headers missing"
    fi
}

# Test error handling
test_error_handling() {
    print_info "Testing error handling..."
    
    # Test invalid endpoint
    local error_response=$(curl -s "$API_URL/invalid-endpoint")
    if echo "$error_response" | grep -q "error\|not found"; then
        print_success "404 error handling works"
    else
        print_error "404 error handling failed"
    fi
    
    # Test invalid JSON
    local invalid_json_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d "invalid json" \
        "$API_URL/preview")
    
    if echo "$invalid_json_response" | grep -q "error"; then
        print_success "Invalid JSON error handling works"
    else
        print_error "Invalid JSON error handling failed"
    fi
}

# Main test execution
main() {
    echo "Starting integration tests at $(date)"
    echo ""
    
    # Check if server is running
    if ! check_server; then
        echo ""
        echo "Please start the web server first:"
        echo "  cd web/"
        echo "  make dev"
        exit 1
    fi
    
    echo ""
    
    # Run all tests
    test_health_endpoints
    test_config_endpoints
    test_blueprint_endpoints
    test_project_preview
    test_project_generation
    test_websocket
    test_frontend_files
    test_cors
    test_error_handling
    
    # Print summary
    echo ""
    echo "=============================================="
    echo "🏁 Integration Test Results"
    echo "=============================================="
    echo -e "${GREEN}✅ Tests Passed: $TESTS_PASSED${NC}"
    echo -e "${RED}❌ Tests Failed: $TESTS_FAILED${NC}"
    
    if [ $TESTS_FAILED -eq 0 ]; then
        echo -e "${GREEN}🎉 All tests passed! Web UI integration is working correctly.${NC}"
        exit 0
    else
        echo -e "${RED}💥 Some tests failed. Please check the issues above.${NC}"
        exit 1
    fi
}

# Run main function
main "$@"