#!/bin/bash

# Phase 4 Test Script
# Tests conditional logger dependencies and logger integration

echo "🎯 Phase 4 Test: Conditional Dependencies & Logger Integration"
echo "============================================================"
echo

# Build the tool first
echo "📦 Building go-starter..."
go build -o bin/go-starter main.go || exit 1
echo "✅ Build successful"
echo

# Test function
test_project() {
    local name=$1
    local type=$2
    local logger=$3
    local framework=$4
    
    echo "🧪 Testing: $name"
    echo "   Type: $type, Logger: $logger, Framework: $framework"
    
    # Clean up any existing project
    rm -rf "/tmp/$name"
    
    # Generate project
    echo -n "   ⚙️  Generating... "
    if [ -n "$framework" ]; then
        ./bin/go-starter new "$name" --type="$type" --logger="$logger" --framework="$framework" --module="github.com/test/$name" --output=/tmp >/dev/null 2>&1
    else
        ./bin/go-starter new "$name" --type="$type" --logger="$logger" --module="github.com/test/$name" --output=/tmp >/dev/null 2>&1
    fi
    
    if [ $? -ne 0 ]; then
        echo "❌ Generation failed"
        return 1
    fi
    echo "✅"
    
    # Check go.mod dependencies
    echo -n "   📋 Checking dependencies... "
    gomod="/tmp/$name/go.mod"
    
    # Expected logger dependencies
    case $logger in
        "zap")
            if ! grep -q "go.uber.org/zap" "$gomod"; then
                echo "❌ Missing zap dependency"
                return 1
            fi
            ;;
        "logrus")
            if ! grep -q "github.com/sirupsen/logrus" "$gomod"; then
                echo "❌ Missing logrus dependency"
                return 1
            fi
            ;;
        "zerolog")
            if ! grep -q "github.com/rs/zerolog" "$gomod"; then
                echo "❌ Missing zerolog dependency"
                return 1
            fi
            ;;
        "slog")
            # slog is built-in, no external dependency
            ;;
    esac
    
    # Check that other loggers are NOT included
    for other_logger in zap logrus zerolog; do
        if [ "$other_logger" != "$logger" ]; then
            case $other_logger in
                "zap")
                    grep -q "go.uber.org/zap" "$gomod" && echo "❌ Unexpected zap dependency" && return 1
                    ;;
                "logrus")
                    grep -q "github.com/sirupsen/logrus" "$gomod" && echo "❌ Unexpected logrus dependency" && return 1
                    ;;
                "zerolog")
                    grep -q "github.com/rs/zerolog" "$gomod" && echo "❌ Unexpected zerolog dependency" && return 1
                    ;;
            esac
        fi
    done
    echo "✅"
    
    # Compile project
    echo -n "   🔨 Compiling... "
    cd "/tmp/$name"
    if [ "$type" = "web-api" ]; then
        cd cmd/server
    fi
    
    if ! go build . >/dev/null 2>&1; then
        echo "❌ Compilation failed"
        cd - >/dev/null
        return 1
    fi
    cd - >/dev/null
    echo "✅"
    
    echo "   ✅ Test passed!"
    echo
    
    # Clean up
    rm -rf "/tmp/$name"
    
    return 0
}

# Run tests
passed=0
failed=0

# Test cases
test_project "test-api-zap" "web-api" "zap" "gin" && ((passed++)) || ((failed++))
test_project "test-api-slog" "web-api" "slog" "gin" && ((passed++)) || ((failed++))
test_project "test-cli-logrus" "cli" "logrus" "cobra" && ((passed++)) || ((failed++))
test_project "test-cli-zerolog" "cli" "zerolog" "cobra" && ((passed++)) || ((failed++))
test_project "test-lib-slog" "library" "slog" "" && ((passed++)) || ((failed++))
test_project "test-lib-zap" "library" "zap" "" && ((passed++)) || ((failed++))
test_project "test-lambda-slog" "lambda" "slog" "" && ((passed++)) || ((failed++))
test_project "test-lambda-zerolog" "lambda" "zerolog" "" && ((passed++)) || ((failed++))

# Summary
echo "============================================================"
echo "📊 Test Summary: $passed passed, $failed failed"
echo

if [ $failed -eq 0 ]; then
    echo "🎉 Phase 4 Complete!"
    echo "✅ Conditional logger dependencies working correctly"
    echo "✅ All templates compile with selected loggers"
    echo "✅ No unnecessary dependencies included"
else
    echo "❌ Some tests failed. Please review the errors above."
    exit 1
fi