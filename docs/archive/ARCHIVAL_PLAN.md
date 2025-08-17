# Documentation Archival Plan

## Overview
This document outlines the files that need to be archived as part of the documentation reorganization.

## Files to Archive

### From 01-getting-started/
- `GETTING_STARTED.md` → Archive (duplicate of getting-started.md)

### From 02-user-guides/
- `BLUEPRINT_SELECTION_GUIDE.md` → Archive (duplicate of blueprint-selection.md)
- `USER_GUIDE.md` → Archive (replaced by new structure)

### From 03-reference/
- `LOGGER_COMPARISON_GUIDE.md` → Archive (content integrated into other docs)

### From 04-developers/
- `CI_INTEGRATION.md` → Keep (still relevant)
- `DEVELOPMENT.md` → Keep (still relevant)
- `TEMPLATE_DOCUMENTATION.md` → Keep (still relevant)
- `TESTING_GUIDE.md` → Keep (still relevant)

### From references/ (old directory)
- All files should be archived as this directory structure is deprecated:
  - `BLUEPRINTS.md`
  - `BLUEPRINT_COMPARISON.md`
  - `LOGGER_GUIDE.md`
  - `ORM_GUIDE.md`
  - `PROJECT_TYPES.md`
  - `QUICK_REFERENCE.md`
  - `QUICK_REFERENCE_CARD.md`

### From testing/ (old directory)
- All files should be archived as this directory structure is deprecated:
  - `API_REFERENCE.md`
  - `ATDD_ARCHITECTURE.md`
  - `CONTRIBUTING_TESTS.md`
  - `ENHANCED_TESTING_GUIDE.md`
  - `MAINTENANCE.md`

## Archival Strategy

1. **Preserve History**: Use git mv to maintain file history
2. **Add Timestamps**: Append dates to archived files for clarity
3. **Create Index**: Update archive/README.md with archived file descriptions
4. **Update Links**: Ensure no broken links remain in active documentation

## Manual Archival Commands

```bash
# From project root
cd /Users/franck/reactive-crafters-workspace/golang-project-generator

# Archive duplicates from numbered directories
git mv docs/01-getting-started/GETTING_STARTED.md docs/archive/GETTING_STARTED_20250116.md
git mv docs/02-user-guides/BLUEPRINT_SELECTION_GUIDE.md docs/archive/BLUEPRINT_SELECTION_GUIDE_20250116.md
git mv docs/02-user-guides/USER_GUIDE.md docs/archive/USER_GUIDE_20250116.md
git mv docs/03-reference/LOGGER_COMPARISON_GUIDE.md docs/archive/LOGGER_COMPARISON_GUIDE_20250116.md

# Archive old references directory
git mv docs/references/*.md docs/archive/references/

# Archive old testing directory
git mv docs/testing/*.md docs/archive/testing/
```

## Status
- [ ] Archive duplicate files from numbered directories
- [ ] Archive old references directory
- [ ] Archive old testing directory
- [ ] Update archive/README.md with descriptions
- [ ] Verify no broken links remain

## Notes
- Some files may already be in archive based on previous cleanup efforts
- Use git status to verify which files still need to be moved
- Consider creating subdirectories in archive/ for better organization