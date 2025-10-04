# Contract: File Move Operations

## Purpose
Define the exact git mv operations required for file reorganization, ensuring git history preservation.

## Prerequisites
- All files exist at source paths
- Destination directories exist
- No conflicts at destination paths
- Working directory is clean (no uncommitted changes)

## Move Operations Contract

### Operation 1: Workspace Migration Documentation
```bash
# Create destination directory
mkdir -p docs/09-workspace-migration

# Move file with git
git mv docs/WORKSPACE_MIGRATION.md docs/09-workspace-migration/README.md

# Verify
[ -f docs/09-workspace-migration/README.md ] || echo "FAIL: File not at destination"
[ ! -f docs/WORKSPACE_MIGRATION.md ] || echo "FAIL: Source file still exists"
```

**Success Criteria**:
- File exists at docs/09-workspace-migration/README.md
- File does not exist at docs/WORKSPACE_MIGRATION.md
- Git log shows file history preserved

### Operation 2: Archive .plans/ Directory Contents
```bash
# Create destination directory
mkdir -p docs/07-reports/archive

# Move each file with timestamp
git mv .plans/CLEANUP_SUMMARY.md docs/07-reports/archive/cleanup-summary-2025-08.md
git mv .plans/CURRENT_STATUS.md docs/07-reports/archive/current-status-2025-08.md
git mv .plans/go-starter-web-ui-plan.md docs/07-reports/archive/web-ui-plan-2025-08.md
git mv .plans/LOGGER_SELECTOR_PLAN.md docs/07-reports/archive/logger-selector-plan-2025-08.md
git mv .plans/PHASE4_VALIDATION_REPORT.md docs/07-reports/archive/phase4-validation-2025-08.md

# Verify directory is empty
[ ! "$(ls -A .plans)" ] || echo "WARN: .plans directory not empty after moves"
```

**Success Criteria**:
- All 5 files exist in docs/07-reports/archive/ with timestamped names
- .plans/ directory contains no .md files
- Git log shows file history preserved for each file

### Operation 3: Archive Web UI Development Reports
```bash
# Destination already exists from Operation 2
git mv web/WEB_UI_FIX_SUMMARY.md docs/07-reports/archive/web-ui-fix-summary-2025-08.md
git mv web/WEB_UI_RESCUE_SUCCESS_REPORT.md docs/07-reports/archive/web-ui-rescue-success-2025-08.md

# Verify
[ -f docs/07-reports/archive/web-ui-fix-summary-2025-08.md ] || echo "FAIL"
[ -f docs/07-reports/archive/web-ui-rescue-success-2025-08.md ] || echo "FAIL"
```

**Success Criteria**:
- Both files exist in archive with timestamped names
- Original files removed from web/
- Git history preserved

### Operation 4: Relocate Web E2E Testing Documentation
```bash
# Move to development documentation
git mv web/README-E2E-Testing.md docs/05-development/web-e2e-testing.md

# Verify
[ -f docs/05-development/web-e2e-testing.md ] || echo "FAIL"
[ ! -f web/README-E2E-Testing.md ] || echo "FAIL"
```

**Success Criteria**:
- File exists at docs/05-development/web-e2e-testing.md
- Original file removed from web/
- Git history preserved

### Operation 5: Relocate Screenshot Automation Documentation
```bash
# Move to development documentation
git mv scripts/README-Screenshots.md docs/05-development/screenshot-automation.md

# Verify
[ -f docs/05-development/screenshot-automation.md ] || echo "FAIL"
[ ! -f scripts/README-Screenshots.md ] || echo "FAIL"
```

**Success Criteria**:
- File exists at docs/05-development/screenshot-automation.md
- Original file removed from scripts/
- Git history preserved

## Post-Move Validation

### Validation Contract
After all move operations complete:

1. **Git Status Check**:
   ```bash
   # Should show renamed files, not added/deleted
   git status | grep "renamed:" | wc -l
   # Expected: 10 (number of moves)
   ```

2. **History Preservation Check**:
   ```bash
   # For each moved file, verify history exists
   git log --follow docs/09-workspace-migration/README.md | head -n 10
   # Expected: Shows commits from when file was docs/WORKSPACE_MIGRATION.md
   ```

3. **Link Validation**:
   ```bash
   # Check for broken links to old paths
   grep -r "docs/WORKSPACE_MIGRATION.md" docs/ README.md CLAUDE.md
   # Expected: No matches (or updated to new path)
   ```

4. **Directory Cleanup**:
   ```bash
   # .plans directory should be empty or removed
   [ ! -d .plans ] || [ ! "$(ls -A .plans)" ]
   # Expected: Success (directory gone or empty)
   ```

## Rollback Contract

If any operation fails:

```bash
# Git provides automatic rollback
git reset --hard HEAD

# Or for individual files
git checkout HEAD -- [file-path]
```

## Success Metrics

- **Total files moved**: 10
- **Total directories created**: 2 (docs/09-workspace-migration, docs/07-reports/archive)
- **Git history preservation**: 100% (all moves use git mv)
- **Broken links**: 0 (after updates)
- **CI pipeline**: PASS (no path breakage)
