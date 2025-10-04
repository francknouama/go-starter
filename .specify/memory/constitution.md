<!--
Sync Impact Report - Constitution Update

Version Change: 1.0.0 → 1.1.0
Change Type: MINOR - New core principle added (Agent-First Implementation)
Modified Principles:
  - NEW: Principle VIII - Agent-First Implementation
  - STRENGTHENED: Agent Coordination Requirements (Development Workflow)
Added Sections: Agent Selection Matrix, Implementation Delegation Rules
Removed Sections: None

Templates Requiring Updates:
- ✅ plan-template.md - UPDATED: Added agent selection in Phase 0, agent delegation in Phase 2, AGENT-FIRST in Phase 4
- ✅ tasks-template.md - UPDATED: Added Agent Delegation section, [AGENT:name] format, validation checklist items
- ⚠️  implement workflow - Should enforce agent selection before execution (documentation task)
- ✅ spec-template.md - No updates required
- ✅ agent-file-template.md - No updates required (template-only)

Follow-up TODOs:
- Update /implement command documentation to reference agent selection
- Add agent selection checklist to CONTRIBUTING.md
- Document agent capabilities matrix in docs/08-agents/
-->

# go-starter Project Constitution

## Core Principles

### I. Production-Ready First
Every blueprint MUST generate code that compiles immediately without errors or warnings. All generated projects MUST include comprehensive tests (unit, integration, ATDD), complete documentation (README, API docs, examples), and production deployment artifacts (Dockerfile, CI/CD configuration).

**Rationale**: Users expect zero-setup, production-ready code. The "30-second to working project" promise is our core value proposition and differentiator from traditional generators.

### II. Progressive Disclosure
Complexity MUST be opt-in through explicit flags (--complexity, --advanced). Default behavior MUST present the simplest viable options. Help systems MUST adapt to user experience level (basic mode shows 14 essential flags, advanced mode shows 18+ flags). Blueprint selection MUST scale from 8-file simple projects to 72-file enterprise applications.

**Rationale**: Beginners should not be overwhelmed by options they don't understand, while power users need access to all capabilities without friction.

### III. Test-First Development (NON-NEGOTIABLE)
All blueprints MUST be validated through ATDD (Acceptance Test-Driven Development) tests. Generated projects MUST include passing tests. Blueprint changes MUST NOT be merged until validation tests pass. Cross-platform compatibility (Windows, macOS, Linux) MUST be verified through automated testing.

**Rationale**: Quality gates prevent regressions and ensure production readiness. The ATDD framework is our quality insurance system.

### IV. Blueprint Quality Standards
Template syntax MUST be error-free (validated by go generate ./...). All template variables MUST be defined and used consistently. Conditional generation MUST work correctly for all configuration combinations. File counts MUST be accurate and documented. Logger integration (slog, zap, logrus, zerolog) MUST work for all applicable blueprints.

**Rationale**: Template errors break user trust and waste time. Consistent quality across all blueprints is essential for professional tooling.

### V. Dual Interface Excellence
CLI and Web UI MUST provide equivalent functionality with interface-appropriate UX. CLI MUST support both interactive prompts and direct flag-based generation. Web UI MUST provide real-time preview via WebSocket. Both interfaces MUST maintain WCAG AA accessibility standards. Mobile responsiveness MUST be verified for Web UI.

**Rationale**: Different users prefer different interfaces. Professional quality in both interfaces establishes go-starter as the industry-leading solution.

### VI. Modular Architecture
Core generation engine, CLI interface, and Web UI MUST be independently versioned and deployable. Go workspace structure MUST enable independent module development and testing. Dependencies MUST be clearly bounded between modules. Integration tests MUST verify module compatibility.

**Rationale**: Workspace migration enables independent development cycles, clearer dependency management, and flexible deployment options.

### VII. Documentation Synchronization
README, CLAUDE.md, and all status guides MUST accurately reflect current blueprint status and capabilities. Blueprint production status MUST be tracked and validated through automated processes. Agent coordination workflows MUST be documented in CLAUDE.md. Hook integration MUST be explained with concrete examples.

**Rationale**: Inaccurate documentation wastes developer time and erodes trust. Synchronized documentation ensures users and AI assistants have accurate project context.

### VIII. Agent-First Implementation (NON-NEGOTIABLE)
Task implementation MUST leverage specialized agents appropriate to the domain. General-purpose execution is ONLY permitted when no specialized agent exists for the task type. Agent selection MUST be documented in task planning phase. Complex multi-step tasks MUST use Task tool to delegate to specialized agents. Direct implementation without agent evaluation is a constitutional violation.

**Mandatory Agent Selection Rules**:
1. **Blueprint Development** → MUST use `golang-fullstack-engineer` or `blueprint-validator`
2. **gRPC/Protobuf Work** → MUST use `grpc-protobuf-specialist`
3. **Template Issues** → MUST use `template-variable-auditor`
4. **ATDD Testing** → MUST use `golang-atdd-qa-engineer`
5. **Security Vulnerabilities** → MUST use `performance-security-specialist`
6. **Documentation Updates** → MUST use `documentation-specialist` or `documentation-sync-manager`
7. **Multi-Agent Workflows** → MUST use `progress-coordinator` for orchestration
8. **Production Readiness** → MUST use `blueprint-production-tracker` for validation

**Rationale**: Specialized agents provide domain expertise that prevents errors, improves quality, and accelerates development. Bypassing agent capabilities wastes the project's investment in specialized tooling and increases risk of suboptimal implementations.

## Quality Standards

### Performance Requirements
Blueprint generation MUST complete in under 30 seconds for standard complexity projects. Generated code MUST compile in under 60 seconds. Web UI real-time preview MUST update in under 500ms. Cross-platform testing MUST complete in under 10 minutes.

**Rationale**: Developer time is valuable. Fast feedback loops improve productivity and user satisfaction.

### Security Requirements
Input validation MUST sanitize project names, module paths, and template variables. Path traversal attacks MUST be prevented in file generation. Template execution MUST run in a secure sandbox. Dependency vulnerabilities MUST be addressed within 7 days of Dependabot alerts.

**Rationale**: Security vulnerabilities can block enterprise adoption and harm users. Proactive security is non-negotiable.

### Accessibility Requirements
Web UI MUST achieve WCAG AA compliance with full keyboard navigation. CLI MUST work with screen readers and provide descriptive help text. Error messages MUST be clear, actionable, and avoid technical jargon for beginners. Documentation MUST include visual aids (screenshots, diagrams) for complex workflows.

**Rationale**: Inclusive design expands our user base and demonstrates professional quality standards.

## Development Workflow

### Blueprint Development Process
1. Create blueprint directory structure in blueprints/
2. Define template.yaml with complete metadata
3. Create all .tmpl files with Go template syntax
4. Update internal/templates/registry.go
5. Add interactive prompts in internal/prompts/
6. Write ATDD validation tests
7. Verify compilation and functionality
8. Update documentation and status tracking

**Rationale**: Systematic process ensures consistency and prevents incomplete blueprint implementations.

### Agent Coordination Requirements (STRENGTHENED)

**Pre-Implementation Agent Selection** (MANDATORY):
Before executing ANY task, implementers MUST:
1. Review task type and domain requirements
2. Consult agent capability matrix (see CLAUDE.md)
3. Select appropriate specialized agent(s)
4. Document agent selection rationale if non-obvious
5. Use Task tool to delegate to selected agent(s)

**Agent Selection Matrix**:
- **Go Code Implementation** → `golang-fullstack-engineer` (ATDD, refactoring, full-stack)
- **gRPC/Protobuf** → `grpc-protobuf-specialist` (buf config, protobuf generation, gRPC-gateway)
- **Template Variables** → `template-variable-auditor` (systematic validation, resolution)
- **Blueprint Validation** → `blueprint-validator` (compilation verification, quality checks)
- **QA Testing** → `golang-atdd-qa-engineer` (ATDD test creation, test review)
- **Security Issues** → `performance-security-specialist` (vulnerabilities, performance)
- **Documentation** → `documentation-specialist` (comprehensive guides, architecture)
- **Status Tracking** → `blueprint-production-tracker` (status management, metrics)
- **Progress Coordination** → `progress-coordinator` (multi-agent orchestration)
- **UI/UX Design** → `ux-design-expert` (interface design, usability)
- **CI/CD & DevOps** → `devops-cicd-specialist` (pipelines, automation)
- **Cross-Platform** → `cross-platform-tester` (Windows, macOS, Linux)

**Multi-Agent Workflow Orchestration**:
Complex tasks requiring multiple domains MUST use `progress-coordinator` to:
- Plan agent execution sequence
- Manage dependencies between agents
- Track progress across agents
- Aggregate results and handle errors

**Prohibited Practices**:
- ❌ Direct implementation without agent evaluation
- ❌ Using general-purpose agent when specialist exists
- ❌ Skipping agent delegation for "simple" domain-specific tasks
- ❌ Manual coordination when progress-coordinator should be used

**Hook Integration**:
Hook-enhanced workflows MUST be preferred for context-aware agent selection. Agents suggested by hooks with >90% confidence SHOULD be used unless justified otherwise.

**Rationale**: Specialized agents provide domain expertise and improve task completion quality. Coordinated workflows prevent missed steps. Mandatory agent selection ensures consistent quality and prevents ad-hoc implementations.

### Commit and Release Standards
Commits MUST include descriptive messages following conventional commits format. Generated code MUST be attributed to Claude with Co-Authored-By footer. Pre-commit hooks MUST be respected (never skip with --no-verify). Force pushes to main/master MUST be explicitly confirmed. Release preparation MUST involve github-project-manager coordination.

**Rationale**: Clear commit history aids debugging and project archaeology. Attribution acknowledges AI assistance transparently.

## Governance

### Amendment Procedure
Constitution amendments MUST increment version using semantic versioning (MAJOR for breaking changes, MINOR for new principles, PATCH for clarifications). All amendments MUST be documented in Sync Impact Report at file header. Dependent templates (plan-template.md, spec-template.md, tasks-template.md) MUST be reviewed for consistency impacts. Ratification MUST occur before amended principles take effect.

### Compliance Review
All feature implementations MUST pass Constitution Check in plan-template.md before Phase 0 research. Blueprint fixes MUST validate compliance with Principles III and IV before merging. Documentation updates MUST maintain synchronization per Principle VII. **Agent selection MUST comply with Principle VIII before task execution**. Quality gates (compilation, testing, cross-platform) MUST block non-compliant changes.

### Complexity Justification
Deviations from core principles MUST be documented in Complexity Tracking section of plan.md. Each violation MUST include rationale and explanation of why simpler alternatives are insufficient. Unjustified complexity MUST result in design revision and return to Phase 1.

### Runtime Development Guidance
Developers and AI assistants MUST consult CLAUDE.md for agent coordination workflows, hook integration patterns, and project-specific development guidelines. Constitution principles supersede ad-hoc practices. When conflicts arise, constitution takes precedence.

**Version**: 1.1.0 | **Ratified**: 2025-10-04 | **Last Amended**: 2025-10-04
