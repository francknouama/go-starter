# Claude Code Agent Orchestration Master Guide

## 🎯 CRITICAL: Agent Selection Rules

This is the definitive guide for Claude Code to properly orchestrate agents in the go-starter project. **ALWAYS** refer to this guide for agent selection.

## 📋 Agent Availability in Task Tool

### **Available Agents for Task Tool**
These agents are available via the Task tool and should be used for specialized tasks:

- `general-purpose` - Multi-domain tasks, documentation (use when documentation-specialist needed)
- `golang-fullstack-engineer` - Go development, ATDD, refactoring
- `senior-bug-resolver` - Critical bugs, complex debugging
- `golang-atdd-qa-engineer` - ATDD tests, acceptance criteria
- `grpc-protobuf-specialist` - gRPC, protobuf, buf configuration
- `template-variable-auditor` - Template syntax, variable validation
- `blueprint-validator` - Blueprint compilation, quality assurance
- `microservice-orchestrator` - Distributed systems, service mesh
- `serverless-specialist` - AWS Lambda, serverless patterns
- `cross-platform-tester` - Windows/macOS/Linux compatibility
- `performance-auditor` - Performance analysis, optimization
- `blueprint-production-pipeline` - End-to-end production readiness
- `web-ui-designer` - React/TypeScript, web UX
- `product-owner` - GitHub issues, project management
- `open-source-community-manager` - Community engagement
- `complexity-analyzer` - Requirements analysis
- `marketing-specialist` - Product positioning, user acquisition
- `sales-specialist` - Enterprise solutions, developer tools
- `progressive-disclosure-optimizer` - UX complexity management
- `ux-design-expert` - Mobile/web UX design, accessibility

### **IMPORTANT: Infrastructure Agents Not Available in Task Tool**
These agents exist as documentation but are **NOT available** in the Task tool:
- ~~`terraform-infrastructure-specialist`~~ → Use `general-purpose` for infrastructure
- ~~`ansible-automation-specialist`~~ → Use `general-purpose` for deployment
- ~~`aws-deployment-specialist`~~ → Use `general-purpose` for AWS tasks  
- ~~`documentation-specialist`~~ → Use `general-purpose` for documentation
- ~~`phase-committer-specialist`~~ → Use `general-purpose` for phase coordination
- ~~`gherkin-bdd-specialist`~~ → Use `golang-atdd-qa-engineer` for BDD

## 🚨 CRITICAL DECISION RULES

### **Rule 1: Infrastructure & Cloud Tasks**
```
User mentions: "Infrastructure", "Terraform", "AWS", "Deployment", "Cloud"
→ PRIMARY: general-purpose (Infrastructure agents not available)
→ SUPPORTING: performance-auditor (if performance-related)
→ VALIDATION: cross-platform-tester (if platform compatibility needed)
```

### **Rule 2: Phase/Milestone/Release Tasks**
```
User mentions: "Phase", "Milestone", "Release", "Complete [major feature]"
→ PRIMARY: general-purpose (Phase-committer not available)
→ COORDINATION: blueprint-production-pipeline (for blueprint readiness)
→ VALIDATION: All relevant quality agents
```

### **Rule 3: Documentation Tasks**
```
User mentions: "Documentation", "Docs", "README", "Guides"
→ PRIMARY: general-purpose (Documentation-specialist not available)
→ SUPPORTING: ux-design-expert (for user experience)
→ REVIEW: open-source-community-manager (for community feedback)
```

### **Rule 4: gRPC/Protobuf Tasks**
```
User mentions: "gRPC", "protobuf", "buf", "gateway"
→ PRIMARY: grpc-protobuf-specialist ✅ AVAILABLE
→ SUPPORTING: template-variable-auditor (for templates)
→ VALIDATION: blueprint-validator (for compilation)
```

### **Rule 5: Blueprint Development**
```
User mentions: "Blueprint", "Production ready", "Enhancement"
→ PRIMARY: blueprint-production-pipeline ✅ AVAILABLE
→ SUPPORTING: golang-fullstack-engineer, template-variable-auditor
→ VALIDATION: blueprint-validator, cross-platform-tester
```

### **Rule 6: Bug Fixes**
```
User reports: "Bug", "Error", "Compilation issue"
→ PRIMARY: senior-bug-resolver ✅ AVAILABLE
→ SUPPORTING: blueprint-validator (for validation)
→ TESTING: cross-platform-tester (for compatibility)
```

### **Rule 7: Testing & Quality**
```
User mentions: "Test", "ATDD", "BDD", "Quality"
→ PRIMARY: golang-atdd-qa-engineer ✅ AVAILABLE (handles both ATDD and BDD)
→ SUPPORTING: blueprint-validator (for quality assurance)
→ VALIDATION: performance-auditor (for performance testing)
```

## 🔄 Optimal Orchestration Patterns

### **Pattern A: Blueprint Production (Available Agents)**
```
1. blueprint-validator → identify compilation/quality issues
2. template-variable-auditor → fix template syntax problems
3. golang-fullstack-engineer → implement fixes and enhancements
4. golang-atdd-qa-engineer → create/validate test coverage
5. cross-platform-tester → validate cross-platform compatibility
6. blueprint-production-pipeline → coordinate final validation
```

### **Pattern B: Infrastructure Tasks (Fallback Pattern)**
```
1. general-purpose → design/implement infrastructure (Terraform/AWS/Ansible)
2. performance-auditor → validate performance implications
3. cross-platform-tester → validate deployment compatibility
```

### **Pattern C: Bug Resolution (Available Agents)**
```
1. senior-bug-resolver → analyze and fix critical issues
2. blueprint-validator → validate fixes don't break functionality
3. cross-platform-tester → ensure cross-platform compatibility
```

### **Pattern D: Feature Development (Available Agents)**
```
1. complexity-analyzer → analyze requirements (optional)
2. golang-fullstack-engineer → implement core functionality
3. golang-atdd-qa-engineer → create acceptance tests and BDD scenarios
4. blueprint-validator → validate implementation quality
5. performance-auditor → validate performance impact
```

### **Pattern E: Phase/Milestone Completion (Fallback Pattern)**
```
1. general-purpose → coordinate phase completion strategy
2. blueprint-production-pipeline → validate blueprint readiness
3. golang-fullstack-engineer → final implementation validation
4. cross-platform-tester → comprehensive platform validation
5. performance-auditor → phase performance validation
```

## ⚡ Performance Best Practices

### **✅ DO: Use Available Agents**
```javascript
// ✅ CORRECT: Use available agents
<function_calls>
  <invoke name="Task">
    <parameter name="subagent_type">golang-fullstack-engineer</parameter>
    <parameter name="description">Blueprint enhancement</parameter>
    <parameter name="prompt">Enhance the lambda blueprint with advanced patterns</parameter>
  </invoke>
  <invoke name="Task">
    <parameter name="subagent_type">blueprint-validator</parameter>
    <parameter name="description">Validation testing</parameter>
    <parameter name="prompt">Validate the enhanced blueprint compiles correctly</parameter>
  </invoke>
</function_calls>
```

### **❌ DON'T: Use Unavailable Agents**
```javascript
// ❌ WRONG: These agents are not available in Task tool
<invoke name="Task">
  <parameter name="subagent_type">terraform-infrastructure-specialist</parameter> // NOT AVAILABLE
</invoke>
<invoke name="Task">
  <parameter name="subagent_type">documentation-specialist</parameter> // NOT AVAILABLE
</invoke>
```

### **✅ DO: Use Fallback Pattern for Infrastructure**
```javascript
// ✅ CORRECT: Use general-purpose for infrastructure tasks
<function_calls>
  <invoke name="Task">
    <parameter name="subagent_type">general-purpose</parameter>
    <parameter name="description">AWS infrastructure setup</parameter>
    <parameter name="prompt">Create Terraform modules and AWS deployment strategy for go-starter applications</parameter>
  </invoke>
</function_calls>
```

### **✅ DO: Parallel Execution When Possible**
```javascript
// ✅ CORRECT: Launch multiple available agents in parallel
<function_calls>
  <invoke name="Task">
    <parameter name="subagent_type">golang-fullstack-engineer</parameter>
    <parameter name="description">Core implementation</parameter>
    <parameter name="prompt">Fix blueprint implementation issues</parameter>
  </invoke>
  <invoke name="Task">
    <parameter name="subagent_type">cross-platform-tester</parameter>
    <parameter name="description">Platform validation</parameter>
    <parameter name="prompt">Validate cross-platform compatibility</parameter>
  </invoke>
  <invoke name="Task">
    <parameter name="subagent_type">performance-auditor</parameter>
    <parameter name="description">Performance analysis</parameter>
    <parameter name="prompt">Analyze performance implications</parameter>
  </invoke>
</function_calls>
```

## 🎯 Task-to-Agent Quick Reference

| User Request Pattern | Primary Agent | Supporting Agents | Notes |
|---------------------|---------------|-------------------|-------|
| "Fix [blueprint] compilation" | `senior-bug-resolver` | `blueprint-validator`, `template-variable-auditor` | ✅ Available |
| "Make [blueprint] production ready" | `blueprint-production-pipeline` | `golang-fullstack-engineer`, quality agents | ✅ Available |
| "Create infrastructure for [app]" | `general-purpose` | `performance-auditor` | ⚠️ Fallback pattern |
| "Deploy to AWS/cloud" | `general-purpose` | `cross-platform-tester` | ⚠️ Fallback pattern |
| "Write comprehensive docs" | `general-purpose` | `ux-design-expert` | ⚠️ Fallback pattern |
| "Implement ATDD tests" | `golang-atdd-qa-engineer` | `blueprint-validator` | ✅ Available |
| "Create BDD scenarios" | `golang-atdd-qa-engineer` | `blueprint-validator` | ✅ Available (handles BDD) |
| "Complete Phase [X]" | `general-purpose` | `blueprint-production-pipeline`, quality agents | ⚠️ Fallback pattern |
| "Fix gRPC/protobuf issues" | `grpc-protobuf-specialist` | `template-variable-auditor` | ✅ Available |
| "Cross-platform testing" | `cross-platform-tester` | `blueprint-validator` | ✅ Available |

## 🚨 Common Mistakes to Avoid

### **❌ MISTAKE 1: Using Unavailable Agents**
Don't try to use terraform-infrastructure-specialist, aws-deployment-specialist, documentation-specialist, or phase-committer-specialist in Task calls.

### **❌ MISTAKE 2: Wrong Agent for Task Domain**
Don't use web-ui-designer for backend Go development or marketing-specialist for technical implementation.

### **❌ MISTAKE 3: Sequential When Parallel is Possible**
Don't wait for one agent to complete before starting independent tasks.

### **❌ MISTAKE 4: Over-Engineering Simple Tasks**
Don't use 5 agents to fix a typo or simple documentation update.

## ✅ Success Indicators

### **Optimal Agent Orchestration:**
- ✅ All selected agents are available in Task tool
- ✅ Right agent expertise matches task domain
- ✅ Parallel execution for independent tasks
- ✅ Appropriate fallback patterns for unavailable agents
- ✅ Clear coordination with lead agent

### **Sub-Optimal Orchestration:**
- ❌ Attempting to use unavailable agents
- ❌ Wrong agent domain for task type
- ❌ Sequential execution when parallel is possible
- ❌ Over-engineering with too many agents

**Use this guide as the definitive reference for all go-starter agent orchestration decisions.**