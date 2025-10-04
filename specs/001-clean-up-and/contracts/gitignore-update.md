# Contract: .gitignore Update

## Purpose
Ensure build artifacts, coverage reports, and temporary files are excluded from version control.

## Current State Analysis

The existing .gitignore covers:
- ✅ Basic Go binaries (*.exe, *.dll, *.so)
- ✅ /bin/, /dist/, /build/ directories
- ✅ go-starter and go-starter-* binaries
- ❌ Coverage JSON reports (MISSING)
- ❌ Web-specific build artifacts (PARTIAL)
- ❌ Specification artifacts (MISSING)

## Required Additions

### Addition 1: Coverage Reports
```gitignore
# =============================================================================
# Testing and Coverage
# =============================================================================

# Coverage data files
*.out
*.prof
coverage.txt
coverage.html

# Coverage reports directory (NEW)
**/coverage-reports/
coverage-*.json
coverage-*.xml
```

**Rationale**: Coverage reports are generated during testing and should not be committed.

**Validation**:
```bash
# Test pattern matches
echo "internal/monitoring/coverage-reports/coverage-2025-10-04.json" | \
  git check-ignore --stdin
# Expected: Match (file would be ignored)
```

### Addition 2: Web Build Artifacts
```gitignore
# =============================================================================
# Web UI Build Artifacts
# =============================================================================

# Web server binaries (NEW - make more specific)
web/bin/
web/dist/
web/build/

# Node modules and package locks (already covered elsewhere, confirm)
web/node_modules/
web/.next/
web/.cache/
```

**Rationale**: Web binaries and build outputs should be generated fresh, not committed.

**Validation**:
```bash
# Test pattern matches
echo "web/bin/web-server" | git check-ignore --stdin
# Expected: Match
```

### Addition 3: Specification Artifacts
```gitignore
# =============================================================================
# Specification and Planning Artifacts
# =============================================================================

# Temporary specification work
.specify/tmp/
specs/**/tmp/

# Local planning notes (not shared)
*.local.md
.notes/
```

**Rationale**: Temporary and local-only specification files should not clutter the repository.

**Validation**:
```bash
# Test pattern matches
echo "specs/001-feature/tmp/draft.md" | git check-ignore --stdin
# Expected: Match
```

### Addition 4: Development Tools
```gitignore
# =============================================================================
# Development Tools and IDE
# =============================================================================

# Claude Code artifacts (if any)
.claude-cache/

# Temporary development files
.DS_Store
*.swp
*.swo
*~
```

**Rationale**: IDE and development tool artifacts are user-specific.

**Validation**:
```bash
# Test pattern matches
echo ".DS_Store" | git check-ignore --stdin
# Expected: Match
```

## Complete Updated Section

```gitignore
# =============================================================================
# Testing and Coverage
# =============================================================================

# Test binary, built with `go test -c`
*.test

# Output of the go coverage tool
*.out
*.prof
coverage.txt
coverage.html

# Coverage reports (generated, not committed)
**/coverage-reports/
coverage-*.json
coverage-*.xml

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

# Emacs
\#*\#
.\#*

# VSCode
.vscode/
!.vscode/settings.json
!.vscode/tasks.json
!.vscode/launch.json
!.vscode/extensions.json

# JetBrains IDEs
.idea/
*.iml
*.iws

# Claude Code
.claude-cache/
```

## Implementation Contract

### Step 1: Backup Current .gitignore
```bash
cp .gitignore .gitignore.backup
```

### Step 2: Update .gitignore
```bash
# Append new sections to .gitignore
# (Manual edit or script-based)
```

### Step 3: Verify No Wanted Files Ignored
```bash
# Check that source files are not accidentally ignored
git check-ignore -v $(git ls-files) | wc -l
# Expected: 0 (no currently tracked files should be ignored)
```

### Step 4: Verify Unwanted Files Are Ignored
```bash
# Create test files to verify patterns
touch web/bin/test-server
touch internal/monitoring/coverage-reports/test-coverage.json

# Check they're ignored
git check-ignore web/bin/test-server
git check-ignore internal/monitoring/coverage-reports/test-coverage.json
# Expected: Both match (are ignored)

# Clean up test files
rm web/bin/test-server
rm internal/monitoring/coverage-reports/test-coverage.json
```

### Step 5: Commit Update
```bash
git add .gitignore
git commit -m "chore: enhance .gitignore for coverage reports and build artifacts

- Add coverage-reports/ directory pattern
- Add web/bin/ specific pattern
- Add specification temporary files pattern
- Add common development tool artifacts

Prevents accidental commit of generated files."
```

## Success Criteria

- [ ] .gitignore updated with new patterns
- [ ] No existing tracked files match new patterns (verified)
- [ ] Build artifacts and coverage reports match patterns (verified)
- [ ] Commit includes descriptive message
- [ ] .gitignore.backup exists (for rollback)

## Rollback Contract

If update causes issues:

```bash
# Restore backup
cp .gitignore.backup .gitignore

# Verify restore
git diff .gitignore
# Should show no changes from pre-update state
```

## Validation After Cleanup

After file deletions and .gitignore update:

```bash
# Verify clean working directory
git status --porcelain
# Expected: No untracked build artifacts

# Rebuild project
make build
go test ./...

# Verify new artifacts are ignored
git status --porcelain | grep "web/bin\\|coverage-"
# Expected: No output (ignored)
```

## Success Metrics

- **Patterns added**: 4 categories
- **False positives**: 0 (no source files ignored)
- **Coverage**: 100% (all artifact types covered)
- **Backwards compatible**: Yes (only additions, no removals)
