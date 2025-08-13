#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ROOT="/Users/franck/reactive-crafters-workspace/golang-project-generator"
TEST_DIR="/tmp/go-starter-blueprint-test"
BLUEPRINTS_DIR="$PROJECT_ROOT/blueprints"

# Build the go-starter CLI first
echo -e "${BLUE}🔨 Building go-starter CLI...${NC}"
cd "$PROJECT_ROOT"
go build -o bin/go-starter main.go
CLI_PATH="$PROJECT_ROOT/bin/go-starter"

# Cleanup and create test directory
echo -e "${BLUE}🧹 Setting up test environment...${NC}"
rm -rf "$TEST_DIR"
mkdir -p "$TEST_DIR"
cd "$TEST_DIR"

# Array to track results
declare -A results

# List of blueprints to test
blueprints=(
    "cli-simple"
    "cli-standard" 
    "event-driven"
    "graphql-api"
    "grpc-gateway"
    "grpc-pure"
    "lambda-proxy"
    "lambda-standard"
    "library-standard"
    "microservice-standard"
    "monolith"
    "web-api-chi"
    "web-api-clean"
    "web-api-ddd"
    "web-api-echo"
    "web-api-fiber"
    "web-api-hexagonal"
    "web-api-standard"
    "workspace"
)

# Function to test a blueprint
test_blueprint() {
    local blueprint="$1"
    local test_name="test-${blueprint}"
    
    echo -e "\n${YELLOW}📋 Testing blueprint: ${blueprint}${NC}"
    
    # Skip if blueprint directory doesn't exist
    if [ ! -d "$BLUEPRINTS_DIR/$blueprint" ]; then
        echo -e "${RED}❌ Blueprint directory not found: $blueprint${NC}"
        results["$blueprint"]="MISSING"
        return 1
    fi
    
    # Create test project directory
    local project_dir="$TEST_DIR/$test_name"
    mkdir -p "$project_dir"
    cd "$project_dir"
    
    # Determine blueprint type and set appropriate flags
    local blueprint_type=""
    local extra_flags=""
    
    case "$blueprint" in
        "cli-"*)
            blueprint_type="cli"
            extra_flags="--complexity=simple"
            ;;
        "web-api-"*)
            blueprint_type="web-api"
            if [[ "$blueprint" == "web-api-clean" ]]; then
                extra_flags="--architecture=clean"
            elif [[ "$blueprint" == "web-api-ddd" ]]; then
                extra_flags="--architecture=ddd"
            elif [[ "$blueprint" == "web-api-hexagonal" ]]; then
                extra_flags="--architecture=hexagonal"
            elif [[ "$blueprint" == "web-api-chi" ]]; then
                extra_flags="--framework=chi"
            elif [[ "$blueprint" == "web-api-echo" ]]; then
                extra_flags="--framework=echo"
            elif [[ "$blueprint" == "web-api-fiber" ]]; then
                extra_flags="--framework=fiber"
            fi
            ;;
        "lambda-"*)
            blueprint_type="lambda"
            ;;
        "grpc-"*)
            blueprint_type="microservice"
            extra_flags="--communication-protocol=grpc"
            ;;
        "microservice-"*)
            blueprint_type="microservice"
            ;;
        "library-"*)
            blueprint_type="library"
            ;;
        "graphql-"*)
            blueprint_type="web-api"
            extra_flags="--framework=gin"
            ;;
        "event-driven")
            blueprint_type="event-driven"
            ;;
        "monolith")
            blueprint_type="monolith"
            ;;
        "workspace")
            blueprint_type="workspace"
            ;;
    esac
    
    # Generate project
    echo "🔧 Generating project with type: $blueprint_type"
    if ! timeout 60 "$CLI_PATH" new "$test_name" \
        --type="$blueprint_type" \
        --module="github.com/test/$test_name" \
        --go-version="1.21" \
        --logger="slog" \
        --quiet \
        $extra_flags 2>&1; then
        echo -e "${RED}❌ Failed to generate project for $blueprint${NC}"
        results["$blueprint"]="GENERATE_FAILED"
        return 1
    fi
    
    # Check if project was created
    if [ ! -f "go.mod" ]; then
        echo -e "${RED}❌ go.mod not found for $blueprint${NC}"
        results["$blueprint"]="NO_GO_MOD"
        return 1
    fi
    
    # Try to build the project
    echo "🔨 Building project..."
    if ! timeout 120 go mod tidy 2>&1; then
        echo -e "${RED}❌ go mod tidy failed for $blueprint${NC}"
        results["$blueprint"]="MOD_TIDY_FAILED"
        return 1
    fi
    
    if ! timeout 120 go build . 2>&1; then
        echo -e "${RED}❌ Build failed for $blueprint${NC}"
        results["$blueprint"]="BUILD_FAILED"
        return 1
    fi
    
    # Count generated files
    local file_count=$(find . -type f -name "*.go" | wc -l)
    echo "📁 Generated $file_count Go files"
    
    echo -e "${GREEN}✅ $blueprint: SUCCESS${NC}"
    results["$blueprint"]="SUCCESS"
    
    # Cleanup
    cd "$TEST_DIR"
    rm -rf "$project_dir"
    
    return 0
}

# Test all blueprints
echo -e "${BLUE}🚀 Starting blueprint validation tests...${NC}"
echo -e "${BLUE}Testing ${#blueprints[@]} blueprints${NC}"

total_count=${#blueprints[@]}
success_count=0
failed_count=0

for blueprint in "${blueprints[@]}"; do
    if test_blueprint "$blueprint"; then
        ((success_count++))
    else
        ((failed_count++))
    fi
done

# Print summary
echo -e "\n${BLUE}📊 VALIDATION SUMMARY${NC}"
echo "======================================"
echo -e "${GREEN}✅ Successful: $success_count/$total_count${NC}"
echo -e "${RED}❌ Failed: $failed_count/$total_count${NC}"
echo ""

# Print detailed results
echo -e "${BLUE}📋 DETAILED RESULTS${NC}"
echo "======================================"
for blueprint in "${blueprints[@]}"; do
    result="${results[$blueprint]}"
    case "$result" in
        "SUCCESS")
            echo -e "${GREEN}✅ $blueprint: $result${NC}"
            ;;
        *)
            echo -e "${RED}❌ $blueprint: $result${NC}"
            ;;
    esac
done

# Cleanup
echo -e "\n${BLUE}🧹 Cleaning up test environment...${NC}"
rm -rf "$TEST_DIR"

# Exit with appropriate code
if [ $failed_count -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All blueprints validated successfully!${NC}"
    exit 0
else
    echo -e "\n${RED}💥 $failed_count blueprint(s) failed validation${NC}"
    exit 1
fi