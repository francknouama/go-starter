# Implementation Plan: Codebase Cleanup and Organization

**Branch**: `001-clean-up-and` | **Date**: 2025-10-04 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-clean-up-and/spec.md`

## Summary

Reorganize the go-starter codebase to improve maintainability and prepare for Go workspace migration. Key actions include: consolidating documentation into docs/ hierarchy, removing build artifacts and coverage reports, archiving historical development reports, standardizing file naming conventions, and updating .gitignore to prevent future artifact commits. This will reduce clutter, improve discoverability, and establish clear organizational guidelines.

## Technical Context

**Language/Version**: Go 1.21+, Shell scripting (bash), Markdown
**Primary Dependencies**: Git (2.x+), Make, golangci-lint
**Storage**: File system (git repository)
**Testing**: Manual validation checklist, git history verification, CI/CD pipeline testing
**Target Platform**: Cross-platform (repository organization is platform-agnostic)
**Project Type**: Single repository with multiple modules (preparing for Go workspace structure)
**Performance Goals**: File operations complete in under 10 minutes, validation under 30 minutes
**Constraints**: Must preserve git history (use git mv only), must not break existing CI/CD, must be reversible
**Scale/Scope**: ~200 markdown files, ~10 file moves, ~7 file deletions, 1 .gitignore update

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Alignment with Core Principles

**VII. Documentation Synchronization** (DIRECTLY APPLICABLE):
- ✅ FR-001: Consolidate documentation into clear hierarchy
- ✅ FR-008: Establish documentation organization guidelines
- ✅ FR-009: Consolidate duplicate documentation
- ✅ FR-010: Verify no broken links

**Rationale**: This cleanup directly supports Constitution Principle VII by ensuring documentation is organized, synchronized, and maintainable.

**VI. Modular Architecture** (SUPPORTING):
- ✅ NFR-005: File organization must support Go workspace migration
- ✅ Module-specific documentation preserved in each module

**Rationale**: Clean module boundaries prepare for workspace migration per Constitution Principle VI.

**Commit and Release Standards** (DIRECTLY APPLICABLE):
- ✅ NFR-001: Use git mv to preserve history (Constitution requirement)
- ✅ NFR-002: Changes must be reversible through git
- ✅ Commits will follow conventional commits format

**Rationale**: Adheres to Development Workflow standards from Constitution.

### Quality Standards Alignment

**Performance Requirements**:
- ✅ File operations under 10 minutes (well under Constitution's 30-second generation threshold analogue)

**Security Requirements**:
- ✅ FR-006: .gitignore prevents accidental commit of sensitive/build files
- ✅ Input validation through manual review (FR-010)

**Accessibility Requirements**:
- ✅ Documentation structure improves discoverability for all users
- ✅ Organization guidelines (FR-008) help new contributors

### Gate Status: PASS ✅

No violations. This feature aligns with and supports constitutional principles, particularly Documentation Synchronization and Modular Architecture.

## Project Structure

### Documentation (this feature)
```
specs/001-clean-up-and/
├── plan.md              # This file (/plan command output)
├── research.md          # Phase 0 output (/plan command) ✅
├── data-model.md        # Phase 1 output (/plan command) ✅
├── quickstart.md        # Phase 1 output (/plan command) ✅
├── contracts/           # Phase 1 output (/plan command) ✅
│   ├── file-move-operations.md
│   ├── file-deletion-operations.md
│   └── gitignore-update.md
└── tasks.md             # Phase 2 output (/tasks command - NOT created by /plan)
```

### Source Code (repository root)

**Current State** (before cleanup):
```
go-starter/
├── README.md, CLAUDE.md, CONTRIBUTING.md, etc. (root docs)
├── .plans/               # 5 historical reports TO ARCHIVE
├── docs/                 # Main documentation directory
│   └── WORKSPACE_MIGRATION.md  # TO MOVE to 09-workspace-migration/
├── web/
│   ├── README.md (keep)
│   ├── TESTING.md (keep)
│   ├── WEB_UI_FIX_SUMMARY.md  # TO ARCHIVE
│   ├── WEB_UI_RESCUE_SUCCESS_REPORT.md  # TO ARCHIVE
│   ├── README-E2E-Testing.md  # TO MOVE to docs/05-development/
│   └── bin/web-server  # TO DELETE (26MB binary)
├── internal/monitoring/coverage-reports/  # TO DELETE (5 JSON files)
├── scripts/
│   └── README-Screenshots.md  # TO MOVE to docs/05-development/
└── [other modules remain unchanged]
```

**Target State** (after cleanup):
```
go-starter/
├── README.md, CLAUDE.md, CONTRIBUTING.md, etc. (root docs, updated guidelines)
├── docs/
│   ├── 01-getting-started/
│   ├── 02-user-guides/
│   ├── 03-blueprints/
│   ├── 04-reference/
│   ├── 05-development/
│   │   ├── web-e2e-testing.md  # MOVED from web/
│   │   └── screenshot-automation.md  # MOVED from scripts/
│   ├── 06-community/
│   ├── 07-reports/
│   │   └── archive/  # NEW - historical reports
│   │       ├── cleanup-summary-2025-08.md
│   │       ├── current-status-2025-08.md
│   │       ├── web-ui-plan-2025-08.md
│   │       ├── logger-selector-plan-2025-08.md
│   │       ├── phase4-validation-2025-08.md
│   │       ├── web-ui-fix-summary-2025-08.md
│   │       └── web-ui-rescue-success-2025-08.md
│   ├── 08-agents/
│   └── 09-workspace-migration/  # NEW
│       └── README.md  # MOVED from docs/WORKSPACE_MIGRATION.md
├── web/
│   ├── README.md (module doc)
│   └── TESTING.md (module testing guide)
├── internal/monitoring/
│   └── coverage-reports/  # Directory preserved, contents deleted
└── [other modules clean and organized]
```

**Structure Decision**: Single repository with modular organization. This cleanup establishes clear boundaries between root-level essentials, documentation hierarchy, and module-specific files, preparing for future Go workspace migration where each module (core, CLI, web) will have independent versioning.

## Phase 0: Outline & Research

See [research.md](./research.md) for complete research findings.

**Key Decisions**:

1. **File Organization Strategy**: Hierarchical documentation structure with 8+1 categories (added 09-workspace-migration)
2. **Build Artifact Handling**: Remove from repository, enhance .gitignore
3. **Workspace Migration Files**: Dedicated tracking in docs/09-workspace-migration/
4. **Development Reports**: Archive to docs/07-reports/archive/ with timestamps
5. **File Naming**: kebab-case for docs, SCREAMING_SNAKE_CASE for root critical files
6. **Validation Approach**: Script-based link validation with manual review

**Research Highlights**:
- Git mv preserves history (verified)
- Existing docs/ has 8-category structure (will add 9th for workspace migration)
- ~27MB of deletable build artifacts identified
- Shell scripting and Make automation available for validation

**Implementation Agents Identified** (Constitution Principle VIII):
- **documentation-specialist** → CONTRIBUTING.md updates, completion report (T038-T039)
- **general-purpose** → Available for coordination if needed
- Direct execution acceptable for simple file operations (git mv, rm, .gitignore)

**Output**: ✅ research.md complete, all technical decisions resolved, agents identified

## Phase 1: Design & Contracts

*Prerequisites: research.md complete* ✅

See detailed artifacts:
- [data-model.md](./data-model.md) - File inventory and categorization
- [contracts/file-move-operations.md](./contracts/file-move-operations.md) - Move operation specifications
- [contracts/file-deletion-operations.md](./contracts/file-deletion-operations.md) - Deletion operation specifications
- [contracts/gitignore-update.md](./contracts/gitignore-update.md) - .gitignore update contract
- [quickstart.md](./quickstart.md) - Validation procedures

**Design Summary**:

### File Categorization

1. **Root Essential Files** (6 files): README.md, CLAUDE.md, CONTRIBUTING.md, CHANGELOG.md, SECURITY.md, LICENSE
2. **Documentation Files** (~200 files): Organized in docs/ with 9 categories
3. **Build Artifacts** (7 files): For deletion
4. **Configuration Files** (multiple): Organized by tool/purpose
5. **Script Files** (multiple): Consolidated in scripts/
6. **Module-Specific Files** (per module): Each module keeps README.md

### Move Operations (10 files)

| Source | Destination | Method |
|--------|-------------|--------|
| docs/WORKSPACE_MIGRATION.md | docs/09-workspace-migration/README.md | git mv |
| .plans/CLEANUP_SUMMARY.md | docs/07-reports/archive/cleanup-summary-2025-08.md | git mv |
| .plans/CURRENT_STATUS.md | docs/07-reports/archive/current-status-2025-08.md | git mv |
| .plans/go-starter-web-ui-plan.md | docs/07-reports/archive/web-ui-plan-2025-08.md | git mv |
| .plans/LOGGER_SELECTOR_PLAN.md | docs/07-reports/archive/logger-selector-plan-2025-08.md | git mv |
| .plans/PHASE4_VALIDATION_REPORT.md | docs/07-reports/archive/phase4-validation-2025-08.md | git mv |
| web/WEB_UI_FIX_SUMMARY.md | docs/07-reports/archive/web-ui-fix-summary-2025-08.md | git mv |
| web/WEB_UI_RESCUE_SUCCESS_REPORT.md | docs/07-reports/archive/web-ui-rescue-success-2025-08.md | git mv |
| web/README-E2E-Testing.md | docs/05-development/web-e2e-testing.md | git mv |
| scripts/README-Screenshots.md | docs/05-development/screenshot-automation.md | git mv |

### Deletion Operations (7 files)

| File/Pattern | Size | Reason |
|--------------|------|--------|
| web/bin/web-server | 26.8 MB | Build artifact (reproducible) |
| internal/monitoring/coverage-reports/*.json | ~100KB | Generated reports (not source) |
| .plans/ (directory) | N/A | After archiving contents |

### .gitignore Additions

Four new categories:
1. Coverage reports (`**/coverage-reports/`, `coverage-*.json`)
2. Web build artifacts (`web/bin/`, `web/dist/`)
3. Specification artifacts (`.specify/tmp/`, `*.local.md`)
4. Development tools (`.DS_Store`, `*.swp`, `.claude-cache/`)

### Organization Guidelines (for CONTRIBUTING.md)

**Where to Place New Files**:
- Documentation → `docs/[category]/filename.md`
- Development Reports → `docs/07-reports/` or `docs/07-reports/archive/`
- Scripts → `scripts/script-name.sh`
- Build Artifacts → Never commit (gitignore)
- Module Docs → `[module]/README.md`

**File Naming**:
- Documentation: `kebab-case-name.md`
- Root essentials: `SCREAMING_SNAKE_CASE.md`
- Scripts: `kebab-case-name.sh`
- Directories: `kebab-case-name` or `01-numbered-category`

**Output**: ✅ data-model.md, contracts/*, quickstart.md, CLAUDE.md updated

## Phase 2: Task Planning Approach

*This section describes what the /tasks command will do - DO NOT execute during /plan*

**Task Generation Strategy**:

1. **From data-model.md**:
   - File inventory → audit tasks (verify current state)
   - File categories → organization tasks (group by operation type)
   - Move operations → git mv tasks with verification
   - Delete operations → safe deletion tasks with backup

2. **From contracts/**:
   - file-move-operations.md → 10 move tasks (5 operations, each with verify step)
   - file-deletion-operations.md → 3 deletion tasks (with safety checks)
   - gitignore-update.md → 1 update task with validation

3. **From quickstart.md**:
   - Validation steps → verification tasks
   - Success criteria → acceptance test tasks
   - Each numbered step → corresponding validation task

**Ordering Strategy**:

1. **Audit Phase** (parallel tasks):
   - [P] Catalog current file locations
   - [P] Verify git status clean
   - [P] Create backup branch

2. **Setup Phase** (sequential):
   - Create new directories (docs/09-workspace-migration/, docs/07-reports/archive/)
   - Verify no conflicts at destination paths

3. **Move Phase** (parallel by category, sequential within category):
   - [P] Move workspace migration docs (1 file)
   - [P] Move .plans/ contents to archive (5 files)
   - [P] Move web reports to archive (2 files)
   - [P] Move development docs (2 files)
   - Sequential verification after each group

4. **Delete Phase** (sequential with safety checks):
   - Verify files are reproducible
   - Delete web build artifacts
   - Delete coverage reports
   - Clean up empty directories

5. **Config Phase** (sequential):
   - Update .gitignore
   - Verify patterns work
   - Test with rebuild

6. **Validation Phase** (mostly parallel):
   - [P] Verify file locations
   - [P] Verify git history preserved
   - [P] Run build
   - [P] Run tests
   - Sequential: Link validation (depends on all moves)
   - Sequential: Update CONTRIBUTING.md (depends on validation)

7. **Documentation Phase** (parallel):
   - [P] Update internal links
   - [P] Update CONTRIBUTING.md with guidelines
   - [P] Create completion report

**Estimated Output**: ~40 numbered, ordered tasks in tasks.md

**Task Breakdown Estimate**:
- Audit: 3 tasks
- Setup: 2 tasks
- Move: 12 tasks (10 moves + 2 verifications)
- Delete: 4 tasks
- Config: 3 tasks
- Validation: 10 tasks (from quickstart.md)
- Documentation: 6 tasks (T038-T039 delegate to `documentation-specialist`)
- Total: ~40 tasks

**Agent Delegation** (from research.md):
- T038-T039: [AGENT:documentation-specialist] for CONTRIBUTING.md and completion report
- T001-T037: Direct execution (simple file operations)

**IMPORTANT**: This phase is executed by the /tasks command, NOT by /plan

## Phase 3+: Future Implementation

*These phases are beyond the scope of the /plan command*

**Phase 3**: Task execution (/tasks command creates tasks.md)
**Phase 4**: Implementation (AGENT-FIRST - Constitution Principle VIII)
- T001-T037: Direct execution (simple file operations with clear rollback)
- T038-T039: **DELEGATE to documentation-specialist** (comprehensive documentation work)
- Use Task tool for agent delegation: `Task("Update CONTRIBUTING.md with organization guidelines", agent="documentation-specialist")`
**Phase 5**: Validation (run quickstart.md, verify all success criteria, update docs)

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

No violations identified. This feature fully aligns with constitutional principles.

## Progress Tracking

*This checklist is updated during execution flow*

**Phase Status**:
- [x] Phase 0: Research complete (/plan command)
- [x] Phase 1: Design complete (/plan command)
- [x] Phase 2: Task planning complete (/plan command - describe approach only)
- [ ] Phase 3: Tasks generated (/tasks command)
- [ ] Phase 4: Implementation complete
- [ ] Phase 5: Validation passed

**Gate Status**:
- [x] Initial Constitution Check: PASS
- [x] Post-Design Constitution Check: PASS
- [x] All NEEDS CLARIFICATION resolved
- [x] Complexity deviations documented (N/A - no deviations)

**Artifacts Generated**:
- [x] research.md (Phase 0)
- [x] data-model.md (Phase 1)
- [x] contracts/file-move-operations.md (Phase 1)
- [x] contracts/file-deletion-operations.md (Phase 1)
- [x] contracts/gitignore-update.md (Phase 1)
- [x] quickstart.md (Phase 1)
- [x] CLAUDE.md updated (Phase 1)

---
*Based on Constitution v1.0.0 - See `.specify/memory/constitution.md`*
