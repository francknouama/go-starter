#!/bin/bash

# gRPC-Pure Template Progress Monitor
# Real-time monitoring script for the Distinguished Engineer

set -e

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

PROJECT_ROOT=$(cd "$(dirname "$0")/.." && pwd)
BLUEPRINT_DIR="$PROJECT_ROOT/blueprints/grpc-pure"

echo -e "${BLUE}🔍 gRPC-Pure Template Progress Monitor${NC}"
echo "=================================================="

# Count templates
CURRENT_COUNT=$(find "$BLUEPRINT_DIR" -name "*.tmpl" | wc -l | tr -d ' ')
TARGET_COUNT=57
MISSING_COUNT=$((TARGET_COUNT - CURRENT_COUNT))

echo -e "Progress: ${GREEN}$CURRENT_COUNT${NC}/${TARGET_COUNT} templates completed"
echo -e "Missing: ${RED}$MISSING_COUNT${NC} templates"

if [ $MISSING_COUNT -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TEMPLATES COMPLETED!${NC}"
    
    echo -e "\n${BLUE}Running final validation...${NC}"
    cd "$PROJECT_ROOT/playground"
    
    # Test generation
    if OUTPUT=$("$PROJECT_ROOT/bin/go-starter" new final-test --type=grpc-pure --module=github.com/test/final --no-git --quiet 2>&1); then
        FILES_CREATED=$(echo "$OUTPUT" | grep "Files created:" | sed 's/.*Files created: *//' | sed 's/ .*//')
        if [ "$FILES_CREATED" -eq 57 ]; then
            echo -e "${GREEN}✅ All 57 files generated successfully!${NC}"
            
            # Test compilation
            cd final-test
            if go mod tidy && go build ./...; then
                echo -e "${GREEN}✅ Project compiles successfully!${NC}"
                echo -e "\n${GREEN}🚀 READY FOR FULL ATDD TESTING!${NC}"
            else
                echo -e "${RED}❌ Compilation failed${NC}"
            fi
        else
            echo -e "${RED}❌ Expected 57 files, got $FILES_CREATED${NC}"
        fi
    else
        echo -e "${RED}❌ Generation failed${NC}"
    fi
    
    exit 0
fi

# Show progress bar
PROGRESS=$((CURRENT_COUNT * 100 / TARGET_COUNT))
echo -e "\nProgress: $PROGRESS%"

# Show next missing files to work on
echo -e "\n${YELLOW}🎯 Priority Missing Templates:${NC}"

# High priority files
PRIORITY_FILES=(
    "configs/config.prod.yaml.tmpl"
    "configs/config.test.yaml.tmpl"
    "internal/auth/interface.go.tmpl"
    "internal/auth/oauth.go.tmpl"
    "internal/models/user.go.tmpl"
    "internal/repository/interface.go.tmpl"
    "internal/repository/user.go.tmpl"
    "internal/tls/config.go.tmpl"
)

for file in "${PRIORITY_FILES[@]}"; do
    if [ ! -f "$BLUEPRINT_DIR/$file" ]; then
        echo -e "  ${RED}❌${NC} $file"
    else
        echo -e "  ${GREEN}✅${NC} $file"
    fi
done

echo -e "\n${BLUE}Quick Commands:${NC}"
echo "  Test current progress: ./validate_grpc_pure.sh"
echo "  Test generation: cd playground && ../bin/go-starter new test --type=grpc-pure --module=github.com/test/grpc --no-git"
echo "  Monitor again: ./monitor_progress.sh"

echo -e "\n${BLUE}💡 Tips for DE:${NC}"
echo "  1. Focus on High Priority files first"
echo "  2. Test generation after each batch of 3-5 files"
echo "  3. Use existing templates as reference for syntax"
echo "  4. Keep templates simple and focused"