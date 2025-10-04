# Tasks: Codebase Cleanup and Organization

**Input**: Design documents from `/specs/001-clean-up-and/`
**Prerequisites**: plan.md, research.md, data-model.md, contracts/, quickstart.md

## Agent Delegation (Constitution Principle VIII)

**For this cleanup feature**:
- **T001-T037**: Direct execution acceptable (simple file operations: git mv, rm, .gitignore)
  - Rationale: Well-defined, low-risk operations with clear rollback procedures
- **T038-T039**: **DELEGATE to documentation-specialist** (comprehensive documentation work)
  - Rationale: Constitutional requirement for documentation expertise
  - Method: Use Task tool with `subagent_type="documentation-specialist"`

**Agent Selection Source**: See `research.md` section "Implementation Agents"

## Path Conventions

All paths are absolute from repository root: `/Users/franck/reactive-crafters-workspace/golang-project-generator`

## Phase 3.1: Setup & Audit

- [ ] T001 [P] Verify git working directory is clean
  ```bash
  git status --porcelain | wc -l
  # Expected: 0 (no uncommitted changes)
  ```

- [ ] T002 [P] Create backup safety branch
  ```bash
  git branch backup-before-cleanup-$(date +%Y%m%d)
  git branch --list backup-before-cleanup-*
  # Verify backup branch exists
  ```

- [ ] T003 [P] Catalog current file locations for verification
  ```bash
  # Document current state
  ls -la docs/WORKSPACE_MIGRATION.md > /tmp/cleanup-audit.txt
  ls -la .plans/*.md >> /tmp/cleanup-audit.txt
  ls -la web/WEB_UI*.md >> /tmp/cleanup-audit.txt
  ls -la web/bin/ >> /tmp/cleanup-audit.txt
  ls -la internal/monitoring/coverage-reports/ >> /tmp/cleanup-audit.txt
  ```

- [ ] T004 Create required destination directories
  ```bash
  mkdir -p docs/09-workspace-migration
  mkdir -p docs/07-reports/archive
  # Verify directories created
  ls -ld docs/09-workspace-migration docs/07-reports/archive
  ```

- [ ] T005 Verify no conflicts at destination paths
  ```bash
  # Check destinations are clear
  [ ! -f docs/09-workspace-migration/README.md ] || echo "CONFLICT: docs/09-workspace-migration/README.md exists"
  [ ! -f docs/07-reports/archive/cleanup-summary-2025-08.md ] || echo "CONFLICT: archive file exists"
  # Expected: No conflicts
  ```

## Phase 3.2: File Move Operations (History Preservation)

### Move Workspace Migration Documentation

- [ ] T006 Move workspace migration doc to dedicated directory
  ```bash
  git mv docs/WORKSPACE_MIGRATION.md docs/09-workspace-migration/README.md
  # Verify move
  [ -f docs/09-workspace-migration/README.md ] && [ ! -f docs/WORKSPACE_MIGRATION.md ]
  ```

### Move .plans/ Contents to Archive

- [ ] T007 [P] Move cleanup summary to archive
  ```bash
  git mv .plans/CLEANUP_SUMMARY.md docs/07-reports/archive/cleanup-summary-2025-08.md
  [ -f docs/07-reports/archive/cleanup-summary-2025-08.md ] || echo "FAIL"
  ```

- [ ] T008 [P] Move current status to archive
  ```bash
  git mv .plans/CURRENT_STATUS.md docs/07-reports/archive/current-status-2025-08.md
  [ -f docs/07-reports/archive/current-status-2025-08.md ] || echo "FAIL"
  ```

- [ ] T009 [P] Move web UI plan to archive
  ```bash
  git mv .plans/go-starter-web-ui-plan.md docs/07-reports/archive/web-ui-plan-2025-08.md
  [ -f docs/07-reports/archive/web-ui-plan-2025-08.md ] || echo "FAIL"
  ```

- [ ] T010 [P] Move logger selector plan to archive
  ```bash
  git mv .plans/LOGGER_SELECTOR_PLAN.md docs/07-reports/archive/logger-selector-plan-2025-08.md
  [ -f docs/07-reports/archive/logger-selector-plan-2025-08.md ] || echo "FAIL"
  ```

- [ ] T011 [P] Move phase 4 validation report to archive
  ```bash
  git mv .plans/PHASE4_VALIDATION_REPORT.md docs/07-reports/archive/phase4-validation-2025-08.md
  [ -f docs/07-reports/archive/phase4-validation-2025-08.md ] || echo "FAIL"
  ```

### Move Web UI Development Reports to Archive

- [ ] T012 [P] Move web UI fix summary to archive
  ```bash
  git mv web/WEB_UI_FIX_SUMMARY.md docs/07-reports/archive/web-ui-fix-summary-2025-08.md
  [ -f docs/07-reports/archive/web-ui-fix-summary-2025-08.md ] || echo "FAIL"
  ```

- [ ] T013 [P] Move web UI rescue report to archive
  ```bash
  git mv web/WEB_UI_RESCUE_SUCCESS_REPORT.md docs/07-reports/archive/web-ui-rescue-success-2025-08.md
  [ -f docs/07-reports/archive/web-ui-rescue-success-2025-08.md ] || echo "FAIL"
  ```

### Move Development Documentation

- [ ] T014 [P] Move web E2E testing docs to development section
  ```bash
  git mv web/README-E2E-Testing.md docs/05-development/web-e2e-testing.md
  [ -f docs/05-development/web-e2e-testing.md ] || echo "FAIL"
  ```

- [ ] T015 [P] Move screenshot automation docs to development section
  ```bash
  git mv scripts/README-Screenshots.md docs/05-development/screenshot-automation.md
  [ -f docs/05-development/screenshot-automation.md ] || echo "FAIL"
  ```

- [ ] T016 Verify all file moves completed successfully
  ```bash
  # Count renamed files in git status
  git status --porcelain | grep "^R" | wc -l
  # Expected: 10 (number of moves)

  # Verify git recognizes as renames (preserves history)
  git status | grep "renamed:" | wc -l
  # Expected: 10
  ```

- [ ] T017 Commit file move operations
  ```bash
  git commit -m "refactor: reorganize documentation structure

- Move workspace migration docs to docs/09-workspace-migration/
- Archive historical reports to docs/07-reports/archive/
- Consolidate development docs in docs/05-development/

All file moves preserve git history using git mv.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

  # Verify commit succeeded
  git log --oneline -1 | grep "refactor: reorganize"
  ```

## Phase 3.3: File Deletion Operations (Build Artifacts)

- [ ] T018 Verify web server binary is reproducible
  ```bash
  # Check file exists and is executable
  file web/bin/web-server 2>/dev/null | grep "executable" || echo "File not found or not executable"

  # Document size before deletion
  ls -lh web/bin/web-server
  ```

- [ ] T019 Delete web build artifacts
  ```bash
  # Remove web server binary
  rm -f web/bin/web-server

  # Verify deletion
  [ ! -f web/bin/web-server ] || echo "FAIL: File still exists"

  # Remove empty bin directory if empty
  rmdir web/bin 2>/dev/null || true
  ```

- [ ] T020 Delete coverage report files
  ```bash
  # List files before deletion
  ls -lh internal/monitoring/coverage-reports/*.json 2>/dev/null

  # Remove all JSON coverage reports
  rm -f internal/monitoring/coverage-reports/coverage-*.json

  # Verify deletion
  ls internal/monitoring/coverage-reports/*.json 2>&1 | grep "No such file" || echo "WARNING: JSON files remain"
  ```

- [ ] T021 Clean up empty .plans directory
  ```bash
  # Check if directory is empty
  [ ! "$(ls -A .plans 2>/dev/null)" ] && rmdir .plans || echo ".plans not empty or doesn't exist"

  # Verify removal or emptiness
  [ ! -d .plans ] || [ ! "$(ls -A .plans)" ] || echo "FAIL: .plans contains files"
  ```

- [ ] T022 Verify build artifacts deleted
  ```bash
  # Check no build artifacts remain tracked
  find . -name "*.test" -o -name "web-server" -o -name "coverage-*.json" | \
    grep -v ".gitignore" | wc -l
  # Expected: 0
  ```

## Phase 3.4: .gitignore Update

- [ ] T023 Backup current .gitignore
  ```bash
  cp .gitignore .gitignore.backup-$(date +%Y%m%d)
  [ -f .gitignore.backup-* ] || echo "FAIL: Backup not created"
  ```

- [ ] T024 Add coverage reports patterns to .gitignore
  ```bash
  # Add to Testing and Coverage section
  cat >> .gitignore << 'EOF'

# Coverage reports (generated, not committed)
**/coverage-reports/
coverage-*.json
coverage-*.xml
EOF

  # Verify patterns added
  grep -q "coverage-reports" .gitignore || echo "FAIL"
  ```

- [ ] T025 Add web build artifacts patterns to .gitignore
  ```bash
  # Add Web UI Build Artifacts section
  cat >> .gitignore << 'EOF'

# =============================================================================
# Web UI Build Artifacts
# =============================================================================

# Web server binaries
web/bin/
web/dist/
web/build/

# Node/npm artifacts
web/node_modules/
web/.next/
web/.cache/
EOF

  # Verify patterns added
  grep -q "web/bin/" .gitignore || echo "FAIL"
  ```

- [ ] T026 Add specification and development tool patterns to .gitignore
  ```bash
  # Add specification and dev tools sections
  cat >> .gitignore << 'EOF'

# =============================================================================
# Specification and Planning Artifacts
# =============================================================================

# Temporary specification work
.specify/tmp/
specs/**/tmp/

# Local planning notes
*.local.md
.notes/

# =============================================================================
# Development Tools and IDE
# =============================================================================

# macOS
.DS_Store

# Vim
*.swp
*.swo
*~

# Claude Code
.claude-cache/
EOF

  # Verify patterns added
  grep -q ".specify/tmp/" .gitignore || echo "FAIL"
  ```

- [ ] T027 Validate .gitignore patterns work correctly
  ```bash
  # Test that source files are NOT ignored
  git check-ignore -v $(git ls-files | head -20) && echo "ERROR: Source files ignored" || echo "PASS"

  # Test that build artifacts WOULD BE ignored
  touch web/bin/test-server
  git check-ignore web/bin/test-server || echo "FAIL: web/bin not ignored"
  rm web/bin/test-server

  touch internal/monitoring/coverage-reports/test-coverage.json
  git check-ignore internal/monitoring/coverage-reports/test-coverage.json || echo "FAIL: coverage not ignored"
  rm internal/monitoring/coverage-reports/test-coverage.json
  ```

- [ ] T028 Commit .gitignore update
  ```bash
  git add .gitignore
  git commit -m "chore: enhance .gitignore for coverage reports and build artifacts

- Add coverage-reports/ directory pattern
- Add web/bin/ specific pattern
- Add specification temporary files pattern
- Add common development tool artifacts

Prevents accidental commit of generated files.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

  # Verify commit
  git log --oneline -1 | grep "chore: enhance"
  ```

## Phase 3.5: Validation (from quickstart.md)

- [ ] T029 [P] Verify file locations correct
  ```bash
  # Check new documentation structure
  ls docs/09-workspace-migration/README.md
  ls docs/07-reports/archive/*.md | wc -l
  # Expected: 7 files in archive
  ```

- [ ] T030 [P] Verify git history preserved
  ```bash
  # Check file history with --follow
  git log --follow --oneline docs/09-workspace-migration/README.md | head -3

  # Should show rename operations
  git log --oneline --name-status -5 | grep "^R"
  ```

- [ ] T031 [P] Verify build still works
  ```bash
  # Clean and rebuild
  make clean 2>/dev/null || true
  make build

  # Verify binary created
  ./bin/go-starter --version
  ```

- [ ] T032 [P] Verify tests still pass
  ```bash
  # Run test suite
  go test ./... -v

  # Run with coverage (should generate reports but not commit)
  go test -cover ./...

  # Verify coverage not tracked
  git status --porcelain | grep "coverage" && echo "FAIL: Coverage committed" || echo "PASS"
  ```

- [ ] T033 Check for broken documentation links
  ```bash
  # Search for references to old paths
  grep -r "docs/WORKSPACE_MIGRATION.md" docs/ README.md CLAUDE.md 2>/dev/null || echo "PASS: No old refs"

  grep -r "\.plans/" docs/ README.md CLAUDE.md 2>/dev/null | grep -v "archive" || echo "PASS: No .plans refs"
  ```

- [ ] T034 Verify .gitignore effectiveness
  ```bash
  # Rebuild to generate artifacts
  make build
  go test -cover ./...

  # Check git status - no untracked build artifacts
  git status --porcelain | grep -E "bin/|coverage" && echo "FAIL: Artifacts tracked" || echo "PASS"
  ```

- [ ] T035 [P] Verify module structure clean
  ```bash
  # Check module READMEs exist
  ls cmd/README.md web/README.md blueprints/README.md

  # Verify minimal documentation in modules
  find cmd/ web/ blueprints/ -name "*.md" -not -name "README.md" | \
    grep -v "test" | wc -l
  # Expected: Small number
  ```

- [ ] T036 Run CI checks locally
  ```bash
  # Linting
  golangci-lint run || echo "Lint issues found (check if pre-existing)"

  # Cross-platform build test
  GOOS=windows go build ./...
  GOOS=darwin go build ./...
  GOOS=linux go build ./...
  ```

## Phase 3.6: Documentation Updates

- [ ] T037 Update internal documentation links
  ```bash
  # Update any links to moved files in README.md
  # (Manual review and edit if needed)
  grep -n "WORKSPACE_MIGRATION\|\.plans/" README.md || echo "No links to update"
  ```

- [ ] T038 [AGENT:documentation-specialist] Update CONTRIBUTING.md with organization guidelines
  ```bash
  # Add file organization section based on data-model.md
  cat >> CONTRIBUTING.md << 'EOF'

## File Organization Guidelines

### Where to Place New Files

- **Documentation** → `docs/[category]/filename.md`
  - Use existing 9 categories: getting-started, user-guides, blueprints, reference, development, community, reports, agents, workspace-migration
- **Development Reports** → `docs/07-reports/` (current) or `docs/07-reports/archive/` (historical)
- **Scripts** → `scripts/script-name.sh`
- **Build Artifacts** → Never commit (ensure .gitignore coverage)
- **Module Documentation** → `[module]/README.md` for module overview

### File Naming Conventions

- **Documentation files**: `kebab-case-name.md`
- **Root essential files**: `SCREAMING_SNAKE_CASE.md` (README, CONTRIBUTING, etc.)
- **Scripts**: `kebab-case-name.sh`
- **Directories**: `kebab-case-name` or `01-numbered-category`

### Before Moving Files

1. Check for incoming links: `grep -r "filename.md" docs/`
2. Update all references
3. Use `git mv` to preserve history
4. Commit with descriptive message
5. Verify CI still passes
EOF

  git add CONTRIBUTING.md
  ```

- [ ] T039 [AGENT:documentation-specialist] Create cleanup completion report
  ```bash
  cat > docs/07-reports/cleanup-completion-$(date +%Y-%m-%d).md << 'EOF'
# Codebase Cleanup Completion Report

**Date**: $(date +%Y-%m-%d)
**Branch**: 001-clean-up-and

## Summary

Successfully reorganized go-starter codebase for improved maintainability and workspace migration readiness.

## Actions Completed

### File Moves (10 files, history preserved)
- Workspace migration docs → docs/09-workspace-migration/
- Historical reports (7 files) → docs/07-reports/archive/
- Development docs (2 files) → docs/05-development/

### File Deletions (~27MB)
- web/bin/web-server (26.8MB)
- Coverage reports (5 JSON files)
- .plans/ directory (after archiving)

### Configuration Updates
- Enhanced .gitignore with 4 new categories
- Updated CONTRIBUTING.md with organization guidelines

## Validation Results

✅ All file moves completed with git history preserved
✅ Build succeeds, all tests pass
✅ No broken documentation links
✅ .gitignore prevents artifact commits
✅ Module structure clean and organized

## Impact

- Project root decluttered
- Documentation easily discoverable
- Clear guidelines for future contributions
- Workspace migration preparation complete
EOF

  git add docs/07-reports/cleanup-completion-*.md
  ```

- [ ] T040 Final commit and verification
  ```bash
  git commit -m "docs: complete codebase cleanup and organization

- Update CONTRIBUTING.md with file organization guidelines
- Add cleanup completion report to docs/07-reports/
- Verify all documentation links working

Completes feature 001-clean-up-and.

🤖 Generated with [Claude Code](https://claude.ai/code)

Co-Authored-By: Claude <noreply@anthropic.com>"

  # Final verification
  echo "=== Cleanup Complete ==="
  echo "Files moved: $(git log --oneline --grep='refactor: reorganize' -1 --name-status | grep '^R' | wc -l)"
  echo "Git status:"
  git status --short
  echo "Branch ready for PR review"
  ```

## Dependencies

**Sequential Dependencies**:
- T001-T003 (audit) → T004-T005 (setup) → T006-T017 (moves) → T018-T022 (deletes) → T023-T028 (.gitignore) → T029-T036 (validation) → T037-T040 (docs)

**Within Phases**:
- Moves (T006-T015) can run in parallel groups, but verification (T016-T017) must follow
- Deletions (T018-T022) must be sequential (safety checks)
- .gitignore updates (T023-T028) must be sequential (build on each other)
- Validation tasks (T029-T036) mostly parallel except link checking (T033) needs all moves complete

## Parallel Execution Examples

### Audit Phase (can run together)
```bash
# T001, T002, T003 in parallel
git status --porcelain | wc -l &
git branch backup-before-cleanup-$(date +%Y%m%d) &
ls -la docs/WORKSPACE_MIGRATION.md > /tmp/cleanup-audit.txt &
wait
```

### Move .plans/ Files (can run together)
```bash
# T007-T011 in parallel
git mv .plans/CLEANUP_SUMMARY.md docs/07-reports/archive/cleanup-summary-2025-08.md &
git mv .plans/CURRENT_STATUS.md docs/07-reports/archive/current-status-2025-08.md &
git mv .plans/go-starter-web-ui-plan.md docs/07-reports/archive/web-ui-plan-2025-08.md &
git mv .plans/LOGGER_SELECTOR_PLAN.md docs/07-reports/archive/logger-selector-plan-2025-08.md &
git mv .plans/PHASE4_VALIDATION_REPORT.md docs/07-reports/archive/phase4-validation-2025-08.md &
wait
```

### Validation Phase (can run together)
```bash
# T029, T030, T031, T032, T035 in parallel
ls docs/09-workspace-migration/README.md &
git log --follow --oneline docs/09-workspace-migration/README.md | head -3 &
make build &
go test ./... &
ls cmd/README.md web/README.md blueprints/README.md &
wait
```

## Notes

- All tasks use absolute paths
- Git mv preserves history (verified in contracts)
- Deletions include safety checks (file type verification)
- .gitignore updates are incremental and tested
- Validation follows quickstart.md procedures
- All commits follow conventional commits format
- Tasks are immediately executable with provided commands
- **Agent delegation**: T038-T039 use `documentation-specialist` (Constitution Principle VIII)

## Success Metrics

- **Total tasks**: 40
- **File moves**: 10 (with history preservation)
- **File deletions**: 7 files (~27MB)
- **Configuration updates**: 1 (.gitignore enhanced)
- **Documentation updates**: 2 (CONTRIBUTING.md, completion report)
- **Estimated execution time**: 20-30 minutes (including validation)
