# Quickstart: Codebase Cleanup Validation

## Purpose
Verify the codebase reorganization was successful and all systems still function correctly.

## Prerequisites
- All file move operations completed
- All file deletion operations completed
- .gitignore updated
- Changes committed to git

## Validation Steps

### Step 1: Verify File Locations (2 minutes)

**Check documentation consolidation**:
```bash
# Verify new structure exists
ls docs/09-workspace-migration/README.md
ls docs/07-reports/archive/

# Count moved files in archive
ls docs/07-reports/archive/*.md | wc -l
# Expected: 7 files (5 from .plans/ + 2 from web/)
```

**Expected Outcome**: ✅ All files in correct locations

**Troubleshooting**: If files missing, check git log for rename operations

### Step 2: Verify Deleted Artifacts (1 minute)

**Check build artifacts removed**:
```bash
# Verify no web binaries
ls web/bin/ 2>/dev/null
# Expected: "No such file or directory" OR empty directory

# Verify no coverage reports
ls internal/monitoring/coverage-reports/*.json 2>/dev/null
# Expected: "No such file or directory" OR no JSON files

# Verify .plans directory cleaned
ls .plans 2>/dev/null
# Expected: "No such file or directory" OR empty
```

**Expected Outcome**: ✅ All build artifacts removed

**Troubleshooting**: If artifacts remain, re-run deletion operations

### Step 3: Verify Git History Preserved (2 minutes)

**Check file history for moved files**:
```bash
# Should show history from old location
git log --follow --oneline docs/09-workspace-migration/README.md | head -5

# Should show rename operation
git log --oneline --name-status HEAD~5..HEAD | grep "R[0-9]"
```

**Expected Outcome**: ✅ File history accessible with `--follow` flag

**Troubleshooting**: If history lost, files were copied instead of moved (use `git mv`)

### Step 4: Verify Build Still Works (3 minutes)

**Run full build**:
```bash
# Clean build from scratch
make clean 2>/dev/null || true
make build

# Verify binary created
./bin/go-starter --version
```

**Expected Outcome**: ✅ Build succeeds, binary runs

**Troubleshooting**: If build fails, check for broken import paths or missing files

### Step 5: Verify Tests Still Pass (5 minutes)

**Run test suite**:
```bash
# Run all tests
go test ./... -v

# Check coverage (should generate new reports, not commit them)
go test -cover ./...

# Verify coverage not committed
git status --porcelain | grep "coverage"
# Expected: No output (coverage files ignored)
```

**Expected Outcome**: ✅ All tests pass, coverage reports generated but not tracked

**Troubleshooting**: If tests fail, check for broken file references in test code

### Step 6: Verify Documentation Links (3 minutes)

**Check internal links**:
```bash
# Search for broken links to old paths
grep -r "docs/WORKSPACE_MIGRATION.md" docs/ README.md CLAUDE.md 2>/dev/null
# Expected: No matches (or updated references)

grep -r "\.plans/" docs/ README.md CLAUDE.md 2>/dev/null | grep -v "archive"
# Expected: No matches (except in archive documenting the move)

# Check for broken relative links
grep -r "\](\.\./" docs/ | grep -v "http" | head -20
# Manual review: Verify these links still work
```

**Expected Outcome**: ✅ No references to old file locations

**Troubleshooting**: Update any found references to new paths

### Step 7: Verify .gitignore Effectiveness (2 minutes)

**Test ignore patterns**:
```bash
# Rebuild to generate artifacts
make build
go test -cover ./...

# Check git status
git status --porcelain
# Expected: No untracked build artifacts (bin/, coverage reports)

# Verify specific patterns
git check-ignore web/bin/web-server || echo "WARN: web binaries not ignored"
git check-ignore internal/monitoring/coverage-reports/coverage-test.json || echo "WARN: coverage not ignored"
```

**Expected Outcome**: ✅ Build artifacts automatically ignored

**Troubleshooting**: If artifacts show as untracked, check .gitignore patterns

### Step 8: Verify Module Structure (2 minutes)

**Check module boundaries still clear**:
```bash
# Each module should have its own README
ls cmd/README.md
ls web/README.md
ls blueprints/README.md

# Verify no cross-module documentation clutter
find cmd/ web/ blueprints/ -name "*.md" -not -name "README.md" | \
  grep -v "test" | wc -l
# Expected: Small number (only essential module docs)
```

**Expected Outcome**: ✅ Module structure clean and clear

**Troubleshooting**: Review and relocate any misplaced documentation files

### Step 9: Verify CI/CD Pipeline (5 minutes)

**Run CI checks locally**:
```bash
# Linting
golangci-lint run
# Expected: Pass (or same issues as before cleanup)

# Cross-platform test (if available)
GOOS=windows go build ./...
GOOS=darwin go build ./...
GOOS=linux go build ./...
# Expected: All platforms build successfully
```

**Expected Outcome**: ✅ All CI checks pass

**Troubleshooting**: If CI fails, check for path dependencies in CI config

### Step 10: Manual Review Checklist (5 minutes)

**Developer experience check**:

- [ ] Can I find the workspace migration docs easily?
  - Navigate to docs/09-workspace-migration/README.md

- [ ] Is the project root clean?
  - `ls` shows only essential files
  - No clutter from build artifacts or old plans

- [ ] Can I understand where to put new files?
  - Review updated CONTRIBUTING.md
  - Guidelines are clear

- [ ] Are historical reports accessible but out of the way?
  - Check docs/07-reports/archive/
  - Files are named clearly with dates

- [ ] Do module READMEs make sense?
  - Read cmd/README.md, web/README.md
  - Context is clear for each module

## Success Criteria

All 10 steps must pass:

1. ✅ Files in correct locations
2. ✅ Artifacts deleted
3. ✅ Git history preserved
4. ✅ Build works
5. ✅ Tests pass
6. ✅ Links updated
7. ✅ .gitignore effective
8. ✅ Modules clean
9. ✅ CI passes
10. ✅ Manual review positive

## Time Estimate

- **Quick validation**: 10 minutes (steps 1-5)
- **Thorough validation**: 30 minutes (all steps)

## Rollback Procedure

If validation fails and issues cannot be quickly resolved:

```bash
# View recent commits
git log --oneline -10

# Identify cleanup commits
# Revert to before cleanup (replace <commit> with actual hash)
git revert <commit>..HEAD

# Or reset if cleanup was in single branch
git reset --hard <commit-before-cleanup>
```

## Post-Validation Actions

After successful validation:

1. **Update CONTRIBUTING.md** with organization guidelines (from data-model.md)
2. **Create PR** for review
3. **Run CI** on PR to verify in clean environment
4. **Document completion** in docs/07-reports/
5. **Archive this spec** to specs/001-clean-up-and/ (completed)

## Future Maintenance

**Weekly check** (automated):
```bash
# Verify no build artifacts committed
git diff --name-only main | grep -E "bin/|coverage-.*\.json|\.test$"
# Expected: No matches
```

**Monthly review**:
- Review docs/ structure for any new clutter
- Check root directory for creeping files
- Verify .gitignore still comprehensive
