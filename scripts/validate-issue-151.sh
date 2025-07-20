#!/bin/bash

# Validation script for Issue #151 documentation requirements
# This script checks that all required elements are present in CLAUDE.md

echo "Validating Issue #151 Documentation Requirements..."
echo "================================================"

CLAUDE_FILE="CLAUDE.md"
ERRORS=0

# Function to check if content exists
check_content() {
    local pattern="$1"
    local description="$2"
    
    if grep -q "$pattern" "$CLAUDE_FILE"; then
        echo "✓ $description"
    else
        echo "✗ MISSING: $description"
        ((ERRORS++))
    fi
}

# Check all required sections
echo -e "\nChecking Required Sections:"
check_content "## CLI Blueprint Audit Findings" "CLI Blueprint Audit Findings section"
check_content "### Two-Tier Approach" "Two-tier approach explanation"
check_content "## Blueprint Selection Guide" "Blueprint Selection Guide section"
check_content "### Blueprint Complexity Levels" "Blueprint Complexity Levels section"
check_content "### Blueprint Selection Matrix" "Blueprint Selection Matrix"
check_content "## Audit Findings Integration" "Audit Findings Integration section"
check_content "## Progressive Complexity Philosophy" "Progressive Complexity Philosophy section"

echo -e "\nChecking Blueprint Coverage:"
check_content "Simple CLI" "Simple CLI blueprint documentation"
check_content "Standard CLI" "Standard CLI blueprint documentation"

echo -e "\nChecking Complexity Levels:"
check_content "Beginner" "Beginner complexity level"
check_content "Intermediate" "Intermediate complexity level"
check_content "Advanced" "Advanced complexity level"
check_content "Expert" "Expert complexity level"

echo -e "\nChecking Issue References:"
check_content "#149" "Issue #149 reference"
check_content "#56" "Issue #56 reference"
check_content "#150" "Issue #150 reference"

echo -e "\nChecking Development Commands:"
check_content "\-\-complexity" "Complexity flag documentation"
check_content "\-\-advanced" "Advanced mode flag documentation"
check_content "Progressive Disclosure" "Progressive disclosure documentation"

echo -e "\nChecking Migration Guidance:"
check_content "Migrating Between Blueprints" "Migration guidance section"
check_content "from simple to standard" "Simple to standard migration"

# Summary
echo -e "\n================================================"
if [ $ERRORS -eq 0 ]; then
    echo "✅ All Issue #151 requirements are satisfied!"
    exit 0
else
    echo "❌ Found $ERRORS missing requirements"
    exit 1
fi