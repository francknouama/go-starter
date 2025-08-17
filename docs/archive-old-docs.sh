#!/bin/bash
# Script to archive old documentation files

DOCS_DIR="/Users/franck/reactive-crafters-workspace/golang-project-generator/docs"
ARCHIVE_DIR="$DOCS_DIR/archive"

# Function to archive a file with a timestamp
archive_file() {
    local src="$1"
    local filename=$(basename "$src")
    local dest="$ARCHIVE_DIR/${filename%.md}_archived_$(date +%Y%m%d).md"
    
    if [ -f "$src" ]; then
        echo "Archiving: $src -> $dest"
        mv "$src" "$dest"
    else
        echo "File not found: $src"
    fi
}

# Archive uppercase documentation files that have lowercase versions
echo "Starting documentation archival process..."

# From 01-getting-started
if [ -f "$DOCS_DIR/01-getting-started/getting-started.md" ] && [ -f "$DOCS_DIR/01-getting-started/GETTING_STARTED.md" ]; then
    archive_file "$DOCS_DIR/01-getting-started/GETTING_STARTED.md"
fi

# From 02-user-guides
if [ -f "$DOCS_DIR/02-user-guides/blueprint-selection.md" ] && [ -f "$DOCS_DIR/02-user-guides/BLUEPRINT_SELECTION_GUIDE.md" ]; then
    archive_file "$DOCS_DIR/02-user-guides/BLUEPRINT_SELECTION_GUIDE.md"
fi

if [ -f "$DOCS_DIR/02-user-guides/USER_GUIDE.md" ]; then
    archive_file "$DOCS_DIR/02-user-guides/USER_GUIDE.md"
fi

# From 03-reference
if [ -f "$DOCS_DIR/03-reference/LOGGER_COMPARISON_GUIDE.md" ]; then
    archive_file "$DOCS_DIR/03-reference/LOGGER_COMPARISON_GUIDE.md"
fi

# From 04-developers
if [ -f "$DOCS_DIR/04-developers/CI_INTEGRATION.md" ]; then
    archive_file "$DOCS_DIR/04-developers/CI_INTEGRATION.md"
fi

if [ -f "$DOCS_DIR/04-developers/DEVELOPMENT.md" ]; then
    archive_file "$DOCS_DIR/04-developers/DEVELOPMENT.md"
fi

if [ -f "$DOCS_DIR/04-developers/TEMPLATE_DOCUMENTATION.md" ]; then
    archive_file "$DOCS_DIR/04-developers/TEMPLATE_DOCUMENTATION.md"
fi

if [ -f "$DOCS_DIR/04-developers/TESTING_GUIDE.md" ]; then
    archive_file "$DOCS_DIR/04-developers/TESTING_GUIDE.md"
fi

# Move old reference files that are now in different sections
OLD_FILES=(
    "$DOCS_DIR/references/BLUEPRINTS.md"
    "$DOCS_DIR/references/BLUEPRINT_COMPARISON.md"
    "$DOCS_DIR/references/LOGGER_GUIDE.md"
    "$DOCS_DIR/references/ORM_GUIDE.md"
    "$DOCS_DIR/references/PROJECT_TYPES.md"
    "$DOCS_DIR/references/QUICK_REFERENCE.md"
    "$DOCS_DIR/references/QUICK_REFERENCE_CARD.md"
)

for file in "${OLD_FILES[@]}"; do
    if [ -f "$file" ]; then
        archive_file "$file"
    fi
done

# Archive old testing documentation
OLD_TESTING_FILES=(
    "$DOCS_DIR/testing/API_REFERENCE.md"
    "$DOCS_DIR/testing/ATDD_ARCHITECTURE.md"
    "$DOCS_DIR/testing/CONTRIBUTING_TESTS.md"
    "$DOCS_DIR/testing/ENHANCED_TESTING_GUIDE.md"
    "$DOCS_DIR/testing/MAINTENANCE.md"
)

for file in "${OLD_TESTING_FILES[@]}"; do
    if [ -f "$file" ]; then
        archive_file "$file"
    fi
done

echo "Documentation archival complete!"
echo "Archived files are in: $ARCHIVE_DIR"