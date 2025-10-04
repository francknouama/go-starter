# Research: Codebase Cleanup and Organization

## Decision: File Organization Strategy

**Chosen Approach**: Hierarchical documentation structure with separation of concerns

### Rationale:
1. **Existing Structure Alignment**: The docs/ directory already has an 8-category structure defined in README (getting-started, user-guides, blueprints, reference, development, community, reports, agents)
2. **Workspace Migration Readiness**: Preparation for Go workspace structure requires clear module boundaries
3. **Developer Experience**: Reduces cognitive load by grouping related files and removing clutter
4. **Git History Preservation**: Using `git mv` maintains file history for archaeology and blame functionality

### Alternatives Considered:
- **Flat Structure**: Rejected - doesn't scale with 200+ documentation files
- **Feature-based Organization**: Rejected - cross-cutting concerns make this complex
- **Keep As-Is**: Rejected - current scattered state hinders development velocity

## Decision: Build Artifact Handling

**Chosen Approach**: Remove from repository, enhance .gitignore, document in README

### Rationale:
1. **Repository Size**: Coverage reports and binaries add unnecessary weight to git history
2. **CI/CD Best Practice**: Build artifacts should be generated fresh in CI, not committed
3. **Cross-Platform Concerns**: Binaries for one platform may not work on another

### Alternatives Considered:
- **Git LFS for Binaries**: Rejected - adds complexity for no benefit (binaries are reproducible)
- **Keep Coverage Reports**: Rejected - historical coverage data belongs in CI artifacts or dedicated tooling

## Decision: Workspace Migration File Organization

**Chosen Approach**: Dedicated tracking in docs/09-workspace-migration/ with script in scripts/

### Rationale:
1. **Discoverability**: Migration documentation should be easy to find
2. **Separation**: Active migration tracking separate from historical reports
3. **Script Location**: Migration tracker script logically belongs with other automation scripts

### Alternatives Considered:
- **Keep at Root**: Rejected - clutters root directory
- **Archive Immediately**: Rejected - migration is active/in-progress work

## Decision: Development Reports Handling

**Chosen Approach**: Move to docs/07-reports/archive/ with timestamped naming

### Rationale:
1. **Historical Value**: Reports document decisions and progress, useful for project archaeology
2. **Discoverability**: Archive structure makes it clear these are historical
3. **Timestamps**: Preserve chronological context

### Alternatives Considered:
- **Delete Entirely**: Rejected - loses valuable project history
- **Keep in .plans/**: Rejected - hidden directory convention is non-standard for Go projects

## Decision: File Naming Conventions

**Chosen Approach**: kebab-case for docs, SCREAMING_SNAKE_CASE for root critical files

### Rationale:
1. **Go Convention**: Go community prefers kebab-case for documentation
2. **Visibility**: SCREAMING_SNAKE_CASE (README.md, CONTRIBUTING.md, CLAUDE.md) signals importance
3. **URL-Friendly**: kebab-case works well in URLs for documentation sites
4. **Consistency**: Single standard across all documentation

### Alternatives Considered:
- **snake_case**: Rejected - less common in Go ecosystem
- **camelCase**: Rejected - poor readability for file names
- **All UPPERCASE**: Rejected - too aggressive, reserve for truly critical root files

## Decision: Documentation Validation Approach

**Chosen Approach**: Script-based link validation with manual review checklist

### Rationale:
1. **Automation**: Script can check internal links programmatically
2. **Manual Review**: Some context requires human judgment (outdated content, conflicting information)
3. **CI Integration**: Can be added to pre-commit or CI pipeline later

### Alternatives Considered:
- **Manual Only**: Rejected - error-prone for 200+ files
- **Full Automation**: Rejected - some checks require semantic understanding

## Technical Constraints

### Must Preserve:
- Git history for all moved files (use `git mv` exclusively)
- Existing docs/ hierarchy categories (8 categories established)
- Working CI/CD pipelines (no path breakage)
- External documentation links (those pointing into repo)

### Must Not:
- Break existing workflows that depend on file paths
- Remove files without verification they're obsolete
- Change file content during moves (separate concern)
- Add build artifacts to version control

## Implementation Approach Summary

1. **Audit Phase**: Catalog all files to be moved/deleted
2. **Validation Phase**: Verify no dependencies will break
3. **Execution Phase**: Systematic `git mv` operations with commits per category
4. **Verification Phase**: Link checking, CI validation, manual review
5. **Documentation Phase**: Update CONTRIBUTING.md with new guidelines

## Implementation Agents (Constitution Principle VIII)

**Agent Selection for Phase 4 Execution**:

### Primary Agents Required:
- **documentation-specialist** → For CONTRIBUTING.md updates with organization guidelines (T038, T039)
  - Rationale: Comprehensive documentation creation requires specialized expertise
  - Tasks: File organization guidelines, completion report drafting

### Secondary/Optional Agents:
- **general-purpose** → For file operations coordination if complexity increases
  - Rationale: Basic file moves and deletions are straightforward bash operations
  - Note: Most tasks (T001-T037) are simple enough for direct execution

### Agent Delegation Decision:
- **Simple Operations** (git mv, rm, .gitignore updates): Direct execution acceptable
  - These are well-defined, low-risk operations with clear rollback procedures
- **Documentation Tasks** (T038-T039): MUST use `documentation-specialist`
  - Constitutional requirement for comprehensive documentation work
- **Workflow Coordination**: `progress-coordinator` available if multi-agent coordination needed
  - Currently not required as tasks are sequential and well-defined

**Constitutional Compliance**: This cleanup feature primarily involves infrastructure operations rather than domain-specific development. Agent delegation applies to documentation tasks per Principle VIII requirements.
