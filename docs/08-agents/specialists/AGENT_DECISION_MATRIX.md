# Agent Decision Matrix for Claude Code

## 🎯 Quick Decision Tree

Use this matrix to quickly select the optimal agent(s) for any go-starter task:

## 📋 Task Type → Agent Selection

### **Compilation/Build Issues**
| Issue Type | Primary Agent | Supporting Agents | Validation |
|------------|---------------|-------------------|------------|
| gRPC/protobuf errors | `grpc-protobuf-specialist` | `template-variable-auditor` | `blueprint-validator` |
| Template syntax errors | `template-variable-auditor` | `senior-bug-resolver` | `blueprint-validator` |
| Go compilation errors | `senior-bug-resolver` | `golang-fullstack-engineer` | `cross-platform-tester` |
| Cross-platform build issues | `cross-platform-tester` | `senior-bug-resolver` | `blueprint-validator` |

### **Blueprint Development**
| Task | Primary Agent | Supporting Agents | Validation |
|------|---------------|-------------------|------------|
| New blueprint creation | `golang-fullstack-engineer` | `template-variable-auditor`, `complexity-analyzer` | `blueprint-validator` |
| Blueprint enhancement | `blueprint-production-pipeline` | `golang-fullstack-engineer`, `template-variable-auditor` | `golang-atdd-qa-engineer` |
| Template variable fixes | `template-variable-auditor` | `blueprint-validator` | `cross-platform-tester` |
| Production readiness | `blueprint-production-pipeline` | All quality agents | `performance-auditor` |

### **Infrastructure & Deployment**
| Task | Primary Agent | Supporting Agents | Coordination |
|------|---------------|-------------------|--------------|
| AWS infrastructure setup | `terraform-infrastructure-specialist` | `aws-deployment-specialist` | `phase-committer-specialist` |
| Deployment automation | `ansible-automation-specialist` | `terraform-infrastructure-specialist` | `phase-committer-specialist` |
| Cloud optimization | `aws-deployment-specialist` | `performance-auditor` | General purpose |
| Multi-cloud strategy | `terraform-infrastructure-specialist` | `ansible-automation-specialist`, `aws-deployment-specialist` | `phase-committer-specialist` |
| Complete infrastructure pipeline | `phase-committer-specialist` | All infrastructure agents | `cross-platform-tester` |

### **Testing & Quality Assurance**
| Task | Primary Agent | Supporting Agents | Validation |
|------|---------------|-------------------|------------|
| ATDD test creation | `golang-atdd-qa-engineer` | `gherkin-bdd-specialist` | `golang-fullstack-engineer` |
| BDD scenario design | `gherkin-bdd-specialist` | `golang-atdd-qa-engineer` | `blueprint-validator` |
| Performance testing | `performance-auditor` | `golang-atdd-qa-engineer` | `cross-platform-tester` |
| Quality assurance | `blueprint-validator` | `golang-atdd-qa-engineer`, `cross-platform-tester` | `performance-auditor` |

### **Feature Development**
| Feature Type | Primary Agent | Supporting Agents | Testing |
|--------------|---------------|-------------------|---------|
| Core Go features | `golang-fullstack-engineer` | `senior-bug-resolver` | `golang-atdd-qa-engineer` |
| gRPC services | `grpc-protobuf-specialist` | `microservice-orchestrator` | `blueprint-validator` |
| Microservices | `microservice-orchestrator` | `grpc-protobuf-specialist`, `serverless-specialist` | `performance-auditor` |
| Serverless functions | `serverless-specialist` | `aws-deployment-specialist` | `golang-atdd-qa-engineer` |
| Web interfaces | `web-ui-designer` | `ux-design-expert` | `cross-platform-tester` |

### **Documentation & Management**
| Task | Primary Agent | Supporting Agents | Review |
|------|---------------|-------------------|--------|
| Technical documentation | General purpose | `ux-design-expert` | `open-source-community-manager` |
| Project management | `product-owner` | `complexity-analyzer` | General purpose |
| Community management | `open-source-community-manager` | `product-owner` | General purpose |
| Requirements analysis | `complexity-analyzer` | `product-owner` | General purpose |

## 🔄 Multi-Agent Orchestration Patterns

### **Pattern 1: Production Pipeline (5-6 agents)**
For taking blueprints from development to production:
```
1. blueprint-validator → identify issues
2. template-variable-auditor → fix template problems  
3. golang-fullstack-engineer → implement fixes
4. golang-atdd-qa-engineer → add test coverage
5. cross-platform-tester → validate compatibility
6. blueprint-production-pipeline → coordinate final validation
```

### **Pattern 2: Infrastructure Setup (3-4 agents)**
For complete application deployment:
```
1. terraform-infrastructure-specialist → provision resources
2. ansible-automation-specialist → configure and deploy
3. aws-deployment-specialist → optimize and secure
4. cross-platform-tester → validate deployment (optional)
```

### **Pattern 3: Bug Resolution (2-3 agents)**
For critical bug fixes:
```
1. senior-bug-resolver → analyze and fix
2. blueprint-validator → validate fix doesn't break other functionality
3. cross-platform-tester → ensure cross-platform compatibility
```

### **Pattern 4: Feature Development (3-5 agents)**
For new feature implementation:
```
1. complexity-analyzer → requirements analysis (optional)
2. golang-fullstack-engineer → core implementation
3. golang-atdd-qa-engineer → acceptance test creation
4. gherkin-bdd-specialist → BDD scenario development
5. performance-auditor → performance validation (if needed)
```

### **Pattern 5: Phase/Milestone Completion (5-8 agents)**
For strategic phase deliveries:
```
1. blueprint-production-pipeline → coordinate blueprint readiness
2. golang-fullstack-engineer → final implementation fixes
3. terraform-infrastructure-specialist → infrastructure readiness
4. ansible-automation-specialist → deployment automation validation
5. aws-deployment-specialist → cloud optimization validation
6. cross-platform-tester → comprehensive platform validation
7. performance-auditor → phase performance validation
8. phase-committer-specialist → strategic coordination and final commit
```

## ⚡ Performance Guidelines

### **✅ DO: Parallel Execution**
- Launch independent agents simultaneously
- Use single message with multiple Task calls
- Don't wait for one agent before starting another if tasks are independent

### **✅ DO: Right-Sized Agent Selection**
- Use most specialized agent for the domain
- Don't use heavyweight agents for simple tasks
- Consider coordination overhead vs task complexity

### **❌ DON'T: Over-Engineering**
- Don't use 5 agents to fix a typo
- Don't use production pipeline for simple documentation updates
- Don't use infrastructure agents for pure development tasks

### **❌ DON'T: Wrong Agent Domain**
- Don't use documentation agents for compilation fixes
- Don't use infrastructure agents for Go development
- Don't use quality agents for project management tasks

## 🚨 Critical Decision Points

### **When User Mentions "Production Ready"**
→ ALWAYS use `blueprint-production-pipeline` as coordinator

### **When User Mentions "gRPC" or "protobuf"**  
→ ALWAYS include `grpc-protobuf-specialist` as primary or supporting

### **When User Mentions "AWS" or "cloud deployment"**
→ ALWAYS include `aws-deployment-specialist` 

### **When User Mentions "Infrastructure" or "Terraform"**
→ ALWAYS include `terraform-infrastructure-specialist`

### **When User Mentions "Testing" or "ATDD"**  
→ ALWAYS include `golang-atdd-qa-engineer` or `gherkin-bdd-specialist`

### **When User Mentions "Cross-platform" or "Windows/Mac/Linux"**
→ ALWAYS include `cross-platform-tester`

### **When User Reports "Bug" or "Error"**
→ PRIORITIZE `senior-bug-resolver` for complex issues

### **When User Mentions "Phase" or "Milestone" or "Release"**
→ ALWAYS use `phase-committer-specialist` for strategic coordination

### **When User Mentions "Infrastructure" AND "Deployment"**
→ COORDINATE: `terraform-infrastructure-specialist` + `ansible-automation-specialist` + `aws-deployment-specialist`

### **When User Mentions "Documentation" (comprehensive)**
→ USE `general-purpose` agent (documentation-specialist not available in Task tool)

### **When User Says "Complete [major feature/phase]"**
→ ALWAYS use `phase-committer-specialist` as lead coordinator

## 🎯 Success Metrics

### **Optimal Orchestration Indicators:**
- ✅ Right agent expertise matches task domain
- ✅ Parallel execution when tasks are independent  
- ✅ Appropriate number of agents for task complexity
- ✅ Clear coordination pattern with lead agent
- ✅ Validation agents included for quality assurance

### **Sub-Optimal Orchestration Indicators:**
- ❌ Single general-purpose agent for specialized tasks
- ❌ Sequential execution when parallel is possible
- ❌ Wrong agent domain for task type
- ❌ Over-engineering with too many agents
- ❌ Missing validation for quality-critical tasks

Use this decision matrix to ensure optimal agent selection and coordination for all go-starter project tasks.