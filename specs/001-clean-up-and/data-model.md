# Data Model: File Inventory and Organization Structure

## File Categories

### 1. Root Essential Files (Keep at Root)
**Location**: `/`
**Naming Convention**: SCREAMING_SNAKE_CASE
**Purpose**: Critical project documentation visible immediately

**Files**:
- README.md - Project overview and quick start
- CLAUDE.md - AI assistant guidance and agent coordination
- CONTRIBUTING.md - Contribution guidelines (will be updated with organization rules)
- CHANGELOG.md - Version history
- SECURITY.md - Security policies
- LICENSE - Project license

**Validation Rules**:
- Must exist at root level
- Must be in SCREAMING_SNAKE_CASE (except LICENSE)
- Content must be up-to-date and accurate

### 2. Documentation Files (Consolidate to docs/)
**Location**: `/docs/`
**Naming Convention**: kebab-case for files, Title Case for directories
**Purpose**: All project documentation in hierarchical structure

**Hierarchy**:
```
docs/
├── 01-getting-started/        # Installation, quick start, first project
├── 02-user-guides/            # How-to guides, configuration, troubleshooting
├── 03-blueprints/             # Blueprint documentation and catalog
├── 04-reference/              # API reference, CLI commands
├── 05-development/            # Developer docs, architecture, testing
├── 06-community/              # Examples, showcases, contributing
├── 07-reports/                # Status reports, milestones
│   └── archive/               # Historical development reports
├── 08-agents/                 # Agent coordination documentation
└── 09-workspace-migration/    # Go workspace migration tracking (NEW)
```

**Current Files to Consolidate**:
- ./docs/WORKSPACE_MIGRATION.md → docs/09-workspace-migration/README.md
- ./.plans/* → docs/07-reports/archive/ (timestamped)
- ./web/WEB_UI_FIX_SUMMARY.md → docs/07-reports/archive/
- ./web/WEB_UI_RESCUE_SUCCESS_REPORT.md → docs/07-reports/archive/
- ./web/README-E2E-Testing.md → docs/05-development/web-e2e-testing.md
- ./scripts/README-Screenshots.md → docs/05-development/screenshot-automation.md

**Validation Rules**:
- All .md files (except root essentials) must be in docs/
- Files must follow kebab-case naming
- No duplicate or conflicting documentation
- Internal links must work after reorganization

### 3. Build Artifacts (Remove from Repository)
**Current Locations**: Multiple
**Action**: Delete, ensure .gitignore coverage
**Purpose**: Remove non-source files from version control

**Files to Remove**:
- web/bin/web-server (26.8 MB binary)
- internal/monitoring/coverage-reports/*.json (5 files)
- Any *.test binaries
- Any go-starter-* compiled binaries outside /bin/

**Gitignore Additions Needed**:
```gitignore
# Coverage reports (should be generated, not committed)
**/coverage-reports/
coverage-*.json

# Web build artifacts
web/bin/

# Test binaries
*.test
```

**Validation Rules**:
- No binary files in git (except documentation images)
- No generated coverage reports committed
- .gitignore must prevent future commits

### 4. Configuration Files (Organize by Purpose)
**Locations**: Remain at root or move to appropriate subdirs
**Naming Convention**: Tool-specific (respect conventions)

**Root Configuration** (Keep):
- .gitignore - Git exclusions
- .golangci.yml - Linter config
- Makefile - Build automation
- go.mod, go.sum - Go dependencies
- package.json - Web dependencies

**Subdirectory Configuration**:
- .claude-hooks/* - Claude Code hooks (keep as-is)
- .gemini/* - Gemini context (keep as-is)
- .specify/* - Specification workflow (keep as-is)

**Validation Rules**:
- Essential configs remain at root
- Tool-specific configs in appropriate directories
- No orphaned config files

### 5. Script Files (Consolidate in scripts/)
**Location**: `/scripts/`
**Purpose**: Automation and utility scripts

**Current State**: Already well-organized
**New Addition**:
- Move workspace-migration-tracker.sh to scripts/ (currently at root)

**Validation Rules**:
- All executable scripts in scripts/
- Clear naming indicating purpose
- README.md in scripts/ explaining usage

### 6. Module-Specific Files (Respect Module Boundaries)
**Locations**: `/cmd/`, `/internal/`, `/web/`, `/blueprints/`, `/tests/`
**Purpose**: Maintain separation of concerns for workspace migration

**Web Module**:
- web/README.md - Keep (module documentation)
- web/TESTING.md - Keep (module testing guide)
- web/WEB_UI_FIX_SUMMARY.md → Archive
- web/WEB_UI_RESCUE_SUCCESS_REPORT.md → Archive

**Other Modules**:
- cmd/README.md - Keep (explains CLI structure)
- blueprints/README.md - Keep (blueprint overview)
- package-managers/README.md - Keep (package manager docs)

**Validation Rules**:
- Each module can have its own README.md
- Module READMEs explain module-specific concerns
- No cross-module documentation duplication
- Historical reports move to docs/07-reports/archive/

## File State Transitions

### Files to Move (with git mv)
| Source | Destination | Reason |
|--------|-------------|--------|
| docs/WORKSPACE_MIGRATION.md | docs/09-workspace-migration/README.md | New category |
| .plans/CLEANUP_SUMMARY.md | docs/07-reports/archive/cleanup-summary-2025-08.md | Historical archive |
| .plans/CURRENT_STATUS.md | docs/07-reports/archive/current-status-2025-08.md | Historical archive |
| .plans/go-starter-web-ui-plan.md | docs/07-reports/archive/web-ui-plan-2025-08.md | Historical archive |
| .plans/LOGGER_SELECTOR_PLAN.md | docs/07-reports/archive/logger-selector-plan-2025-08.md | Historical archive |
| .plans/PHASE4_VALIDATION_REPORT.md | docs/07-reports/archive/phase4-validation-2025-08.md | Historical archive |
| web/WEB_UI_FIX_SUMMARY.md | docs/07-reports/archive/web-ui-fix-summary-2025-08.md | Historical archive |
| web/WEB_UI_RESCUE_SUCCESS_REPORT.md | docs/07-reports/archive/web-ui-rescue-success-2025-08.md | Historical archive |
| web/README-E2E-Testing.md | docs/05-development/web-e2e-testing.md | Development docs |
| scripts/README-Screenshots.md | docs/05-development/screenshot-automation.md | Development docs |
| scripts/workspace-migration-tracker.sh | scripts/workspace-migration-tracker.sh | Already correct (if at root, move) |

### Files to Delete
| File | Reason |
|------|--------|
| web/bin/web-server | Build artifact, reproducible |
| internal/monitoring/coverage-reports/coverage-2025-*.json | Generated reports, not source |
| .plans/ (directory) | After moving contents to archive |

### Directories to Create
| Path | Purpose |
|------|---------|
| docs/09-workspace-migration/ | Workspace migration tracking |
| docs/07-reports/archive/ | Historical development reports |

## Organization Rules (for CONTRIBUTING.md)

### Where to Place New Files:

1. **Documentation** → `docs/[appropriate-category]/filename.md`
2. **Development Reports** → `docs/07-reports/` (current) or `docs/07-reports/archive/` (historical)
3. **Scripts** → `scripts/script-name.sh`
4. **Build Artifacts** → Never commit (ensure .gitignore coverage)
5. **Module Documentation** → `[module]/README.md` for module overview
6. **Root Files** → Only essential project-level docs (README, CONTRIBUTING, etc.)

### File Naming Rules:

- Documentation files: `kebab-case-name.md`
- Root essential files: `SCREAMING_SNAKE_CASE.md`
- Scripts: `kebab-case-name.sh` or `snake_case_name.sh`
- Directories: `kebab-case-name` or numbered with leading zero (`01-category`)

### Before Moving Files:

1. Check for incoming links with: `grep -r "filename.md" docs/`
2. Update all references
3. Use `git mv` to preserve history
4. Commit with descriptive message explaining move
5. Verify CI still passes
