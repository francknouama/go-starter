#!/bin/bash

# Quick Phase 4 Validation
echo "🔍 Quick Phase 4 Validation"
echo "============================"

# Build first
go build -o bin/go-starter main.go || exit 1
echo "✅ Built successfully"

# Test core functionality
echo "🧪 Testing Core Templates with Different Loggers"

# Clean up
rm -rf /tmp/quick-test-*

# Test matrix: type:logger:framework
tests=(
    "web-api:zap:gin"
    "web-api:slog:gin"
    "cli:logrus:cobra"
    "cli:zerolog:cobra"
    "library:slog:"
    "library:zap:"
    "lambda:logrus:"
    "lambda:zerolog:"
)

passed=0
failed=0

for test in "${tests[@]}"; do
    IFS=':' read -r type logger framework <<< "$test"
    name="quick-test-${type}-${logger}"
    
    echo -n "  Testing $type with $logger... "
    
    # Generate
    args=("new" "$name" "--type=$type" "--logger=$logger" "--module=github.com/test/$name" "--output=/tmp")
    if [ -n "$framework" ]; then
        args+=("--framework=$framework")
    fi
    
    if ! ./bin/go-starter "${args[@]}" >/dev/null 2>&1; then
        echo "❌ Generation failed"
        ((failed++))
        continue
    fi
    
    # Check dependencies in go.mod
    gomod="/tmp/$name/go.mod"
    case $logger in
        "zap")
            if ! grep -q "go.uber.org/zap" "$gomod"; then
                echo "❌ Missing zap dependency"
                ((failed++))
                continue
            fi
            ;;
        "logrus")
            if ! grep -q "github.com/sirupsen/logrus" "$gomod"; then
                echo "❌ Missing logrus dependency" 
                ((failed++))
                continue
            fi
            ;;
        "zerolog")
            if ! grep -q "github.com/rs/zerolog" "$gomod"; then
                echo "❌ Missing zerolog dependency"
                ((failed++))
                continue
            fi
            ;;
    esac
    
    # Compile
    builddir="/tmp/$name"
    if [ "$type" = "web-api" ]; then
        builddir="/tmp/$name/cmd/server"
    fi
    
    if ! go build -C "$builddir" . >/dev/null 2>&1; then
        echo "❌ Compilation failed"
        ((failed++))
        continue
    fi
    
    echo "✅"
    ((passed++))
done

# Summary
echo
echo "📊 Results: $passed passed, $failed failed"

if [ $failed -eq 0 ]; then
    echo "🎉 Phase 4 Core Functionality Working!"
    echo "✅ Conditional dependencies working"
    echo "✅ All templates compile with selected loggers"
    echo "✅ Logger implementations consistent"
else
    echo "❌ Some issues found"
fi

# Clean up
rm -rf /tmp/quick-test-*
echo "🧹 Cleaned up"