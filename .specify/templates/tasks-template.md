# Tasks: [FEATURE NAME]

**Input**: Design documents from `/specs/[###-feature-name]/`
**Prerequisites**: plan.md (required), research.md, data-model.md, contracts/

## Execution Flow (main)
```
1. Load plan.md from feature directory
   → If not found: ERROR "No implementation plan found"
   → Extract: tech stack, libraries, structure, implementation agents
2. Load optional design documents:
   → data-model.md: Extract entities → model tasks
   → contracts/: Each file → contract test task
   → research.md: Extract decisions → setup tasks, agent assignments
3. Generate tasks by category:
   → Setup: project init, dependencies, linting
   → Tests: contract tests, integration tests
   → Core: models, services, CLI commands
   → Integration: DB, middleware, logging
   → Polish: unit tests, performance, docs
   → **Agent Delegation**: Domain-specific tasks requiring specialists
4. Apply task rules:
   → Different files = mark [P] for parallel
   → Same file = sequential (no [P])
   → Tests before implementation (TDD)
   → **Complex tasks → Delegate to specialized agent (Principle VIII)**
5. Number tasks sequentially (T001, T002...)
6. Generate dependency graph
7. Create parallel execution examples
8. Add agent delegation guidance for complex tasks
9. Validate task completeness:
   → All contracts have tests?
   → All entities have models?
   → All endpoints implemented?
   → **Domain-specific tasks have agent assignments?**
10. Return: SUCCESS (tasks ready for execution)
```

## Format: `[ID] [P?] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[AGENT:name]**: Delegate to specialized agent (Constitution Principle VIII)
- Include exact file paths in descriptions

## Agent Delegation (Constitution Principle VIII - MANDATORY)

**Before implementation, identify tasks requiring specialized agents**:

### Agent Assignment Patterns
- **Go Implementation** → `[AGENT:golang-fullstack-engineer]` for complex Go code
- **gRPC/Protobuf** → `[AGENT:grpc-protobuf-specialist]` for buf config, protobuf
- **Template Work** → `[AGENT:template-variable-auditor]` for template validation
- **Blueprint Tasks** → `[AGENT:blueprint-validator]` for compilation checks
- **ATDD Testing** → `[AGENT:golang-atdd-qa-engineer]` for test creation
- **Security Issues** → `[AGENT:performance-security-specialist]` for vulnerabilities
- **Documentation** → `[AGENT:documentation-specialist]` for comprehensive docs
- **Multi-Step** → `[AGENT:progress-coordinator]` for workflow orchestration

### Example Agent Tasks
```
- [ ] T015 [AGENT:golang-fullstack-engineer] Implement user service with CRUD operations in src/services/user_service.go
- [ ] T020 [AGENT:grpc-protobuf-specialist] Generate protobuf files and configure buf.yaml
- [ ] T025 [AGENT:template-variable-auditor] Validate all template variables in blueprints/
- [ ] T030 [AGENT:progress-coordinator] Coordinate multi-agent deployment workflow
```

**Constitutional Requirement**: Domain-specific tasks MUST be delegated to appropriate specialized agents. Direct implementation without agent evaluation is a violation.

## Path Conventions
- **Single project**: `src/`, `tests/` at repository root
- **Web app**: `backend/src/`, `frontend/src/`
- **Mobile**: `api/src/`, `ios/src/` or `android/src/`
- Paths shown below assume single project - adjust based on plan.md structure

## Phase 3.1: Setup
- [ ] T001 Create project structure per implementation plan
- [ ] T002 Initialize [language] project with [framework] dependencies
- [ ] T003 [P] Configure linting and formatting tools

## Phase 3.2: Tests First (TDD) ⚠️ MUST COMPLETE BEFORE 3.3
**CRITICAL: These tests MUST be written and MUST FAIL before ANY implementation**
- [ ] T004 [P] Contract test POST /api/users in tests/contract/test_users_post.py
- [ ] T005 [P] Contract test GET /api/users/{id} in tests/contract/test_users_get.py
- [ ] T006 [P] Integration test user registration in tests/integration/test_registration.py
- [ ] T007 [P] Integration test auth flow in tests/integration/test_auth.py

## Phase 3.3: Core Implementation (ONLY after tests are failing)
- [ ] T008 [P] User model in src/models/user.py
- [ ] T009 [P] UserService CRUD in src/services/user_service.py
- [ ] T010 [P] CLI --create-user in src/cli/user_commands.py
- [ ] T011 POST /api/users endpoint
- [ ] T012 GET /api/users/{id} endpoint
- [ ] T013 Input validation
- [ ] T014 Error handling and logging

## Phase 3.4: Integration
- [ ] T015 Connect UserService to DB
- [ ] T016 Auth middleware
- [ ] T017 Request/response logging
- [ ] T018 CORS and security headers

## Phase 3.5: Polish
- [ ] T019 [P] Unit tests for validation in tests/unit/test_validation.py
- [ ] T020 Performance tests (<200ms)
- [ ] T021 [P] Update docs/api.md
- [ ] T022 Remove duplication
- [ ] T023 Run manual-testing.md

## Dependencies
- Tests (T004-T007) before implementation (T008-T014)
- T008 blocks T009, T015
- T016 blocks T018
- Implementation before polish (T019-T023)

## Parallel Example
```
# Launch T004-T007 together:
Task: "Contract test POST /api/users in tests/contract/test_users_post.py"
Task: "Contract test GET /api/users/{id} in tests/contract/test_users_get.py"
Task: "Integration test registration in tests/integration/test_registration.py"
Task: "Integration test auth in tests/integration/test_auth.py"
```

## Notes
- [P] tasks = different files, no dependencies
- Verify tests fail before implementing
- Commit after each task
- Avoid: vague tasks, same file conflicts

## Task Generation Rules
*Applied during main() execution*

1. **From Contracts**:
   - Each contract file → contract test task [P]
   - Each endpoint → implementation task
   - Complex endpoints → Add [AGENT:golang-fullstack-engineer]

2. **From Data Model**:
   - Each entity → model creation task [P]
   - Relationships → service layer tasks
   - Complex models → Add [AGENT:golang-fullstack-engineer]

3. **From User Stories**:
   - Each story → integration test [P]
   - Quickstart scenarios → validation tasks
   - Test creation → Add [AGENT:golang-atdd-qa-engineer]

4. **From Research (Agent Assignments)**:
   - Load agent list from research.md "Implementation Agents" section
   - Map agents to appropriate task categories
   - Add [AGENT:name] marker to domain-specific tasks

5. **Ordering**:
   - Setup → Tests → Models → Services → Endpoints → Polish
   - Dependencies block parallel execution
   - Agent tasks can run in parallel if independent domains

## Validation Checklist
*GATE: Checked by main() before returning*

- [ ] All contracts have corresponding tests
- [ ] All entities have model tasks
- [ ] All tests come before implementation
- [ ] Parallel tasks truly independent
- [ ] Each task specifies exact file path
- [ ] No task modifies same file as another [P] task
- [ ] **Domain-specific tasks have [AGENT:name] assignments (Principle VIII)**
- [ ] **Agent assignments match capability matrix from constitution**
- [ ] **Multi-domain workflows delegate to progress-coordinator**