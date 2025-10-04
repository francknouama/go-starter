# Contract: File Deletion Operations

## Purpose
Define safe deletion of build artifacts and temporary files not needed in version control.

## Safety Rules
1. **Never delete source code**
2. **Never delete unique documentation**
3. **Verify files are reproducible before deletion**
4. **Check file size before deleting (warn if >1MB)**

## Deletion Operations Contract

### Operation 1: Remove Web Build Artifacts
```bash
# Check file exists and size
ls -lh web/bin/web-server 2>/dev/null

# Verify it's a binary (reproducible)
file web/bin/web-server | grep "executable"

# Remove
rm -f web/bin/web-server

# Verify deletion
[ ! -f web/bin/web-server ] || echo "FAIL: File still exists"
```

**Success Criteria**:
- web/bin/web-server file removed
- Directory web/bin/ may be removed if empty
- Build can recreate the file

**Reproducibility Test**:
```bash
# File should be recreatable with
cd web && go build -o bin/web-server ./cmd/web-server
```

### Operation 2: Remove Coverage Reports
```bash
# List files to be deleted
ls -lh internal/monitoring/coverage-reports/

# Remove all JSON coverage reports
rm -f internal/monitoring/coverage-reports/coverage-*.json

# Verify
ls internal/monitoring/coverage-reports/ | grep "coverage-.*\.json" && echo "FAIL" || echo "PASS"
```

**Success Criteria**:
- All coverage-*.json files removed
- Directory structure preserved (may be used for future reports)
- No committed coverage data in git

**Reproducibility Test**:
```bash
# Files should be recreatable with
go test -coverprofile=coverage.out ./...
# (generates new coverage data)
```

### Operation 3: Clean Up Empty Directories
```bash
# After file deletions, remove empty directories
[ "$(ls -A .plans)" ] || rmdir .plans
[ "$(ls -A web/bin)" ] || rmdir web/bin

# Verify
[ ! -d .plans ] || [ "$(ls -A .plans)" ] || echo "FAIL: Empty .plans exists"
```

**Success Criteria**:
- Empty directories removed
- Non-empty directories preserved
- No accidental data loss

## Pre-Deletion Verification

### Checklist Contract
Before ANY deletion:

1. **Size Check**:
   ```bash
   # Warn if deleting large files
   find . -name "coverage-*.json" -size +1M -ls
   ```

2. **Uniqueness Check**:
   ```bash
   # Ensure files are not unique source/docs
   # (Coverage reports and binaries are generated, not unique)
   ```

3. **Git Status Check**:
   ```bash
   # Files should be untracked or in .gitignore
   git status --porcelain | grep "coverage-"
   ```

4. **Backup Check** (optional, for safety):
   ```bash
   # For first-time cleanup, consider creating backup
   tar -czf cleanup-backup-$(date +%Y%m%d).tar.gz \
     web/bin/ \
     internal/monitoring/coverage-reports/
   ```

## Post-Deletion Validation

### Validation Contract

1. **Verify Deletions**:
   ```bash
   # No build artifacts remain
   find . -name "*.test" -o -name "go-starter-*" -o -name "coverage-*.json" | \
     grep -v ".gitignore" | \
     wc -l
   # Expected: 0
   ```

2. **Verify Build Still Works**:
   ```bash
   # Ensure deletion didn't break build
   make build
   # Expected: Success

   make test
   # Expected: Tests pass
   ```

3. **Verify Git Clean**:
   ```bash
   # No tracked files deleted (only untracked artifacts)
   git status --porcelain | grep "^D "
   # Expected: No output (or only .plans/ if committed before)
   ```

## Rollback Contract

If deletion was premature:

```bash
# If backed up
tar -xzf cleanup-backup-YYYYMMDD.tar.gz

# If committed
git checkout HEAD -- [file-path]

# If neither, rebuild
make build
go test -cover ./...
```

## Success Metrics

- **Binary files removed**: 1 (web-server)
- **Coverage reports removed**: ~5
- **Total space saved**: ~27 MB
- **Build functionality**: Preserved (100%)
- **Test functionality**: Preserved (100%)
- **Accidental deletions**: 0

## .gitignore Update (see separate contract)

After deletion, ensure .gitignore prevents recommit:
```gitignore
# Coverage reports
**/coverage-reports/
coverage-*.json

# Web build artifacts
web/bin/

# All test binaries
*.test
```
