# Go-Starter Project Agents

This directory contains specialized AI agents designed to assist with various aspects of the go-starter project development and maintenance.

## Available Agents

### 🔍 blueprint-validator
**Purpose**: Validates Go blueprint templates and ensures generated code quality

**Use when**:
- Adding or modifying blueprint templates
- Testing that generated projects compile
- Validating template syntax and structure
- Ensuring consistency across blueprints

**Example**: `/agent blueprint-validator validate the new grpc blueprint`

### 🔧 grpc-protobuf-specialist
**Purpose**: Expert in gRPC service development, protobuf generation, and buf configuration

**Use when**:
- Working with gRPC or grpc-gateway blueprints
- Fixing protobuf compilation issues
- Configuring buf tools and googleapis dependencies
- Implementing gRPC streaming or advanced patterns
- Resolving gRPC-related template or generation issues

**Example**: `/agent grpc-protobuf-specialist fix the grpc-gateway blueprint buf configuration`

### 🔍 template-variable-auditor
**Purpose**: Systematic template variable analysis, resolution, and validation across blueprints

**Use when**:
- Fixing template variable mismatches or undefined variables
- Adding new variables consistently across multiple blueprints
- Conducting comprehensive template audits
- Resolving Go template syntax errors
- Ensuring template variable consistency

**Example**: `/agent template-variable-auditor audit template variables across all web-api blueprints`

### 🚀 blueprint-production-pipeline
**Purpose**: End-to-end blueprint production readiness pipeline automation

**Use when**:
- Taking blueprints from development to production-ready status
- Coordinating fix-validate-document-update workflows
- Managing blueprint lifecycle and quality gates
- Ensuring comprehensive production readiness validation
- Automating blueprint status tracking and reporting

**Example**: `/agent blueprint-production-pipeline move grpc-gateway blueprint to production-ready status`

### ☁️ serverless-specialist
**Purpose**: AWS Lambda, serverless patterns, and cloud-native architecture expert

**Use when**:
- Working with Lambda blueprints (standard, proxy, event-processing)
- Implementing serverless microservice patterns
- Optimizing Lambda performance and cost
- Integrating with AWS services (API Gateway, EventBridge, DynamoDB)
- Designing event-driven architectures

**Example**: `/agent serverless-specialist enhance lambda-event-processing with advanced event patterns`

### 🏗️ microservice-orchestrator
**Purpose**: Microservice patterns, service mesh, containerization, and distributed systems

**Use when**:
- Working with microservice blueprints
- Implementing service mesh integration (Istio, Linkerd)
- Designing distributed system patterns (circuit breakers, saga)
- Kubernetes orchestration and deployment
- Inter-service communication and resilience patterns

**Example**: `/agent microservice-orchestrator add service mesh and resilience patterns to microservice blueprint`

### 🧪 atdd-test-creator
**Purpose**: Creates comprehensive ATDD tests for new features

**Use when**:
- Implementing new features that need acceptance tests
- Improving test coverage for existing features
- Creating BDD-style test scenarios
- Enhancing the test infrastructure

**Example**: `/agent atdd-test-creator create tests for the new database migration feature`

### 🎯 progressive-disclosure-optimizer
**Purpose**: Optimizes the progressive disclosure system for better UX

**Use when**:
- Adding new CLI flags or options
- Improving help documentation
- Refining complexity levels
- Enhancing the beginner experience

**Example**: `/agent progressive-disclosure-optimizer review the new deployment flags`

### ⚡ performance-auditor
**Purpose**: Analyzes and optimizes performance across the project

**Use when**:
- Blueprint generation is slow
- Generated projects have performance issues
- Test suite execution needs optimization
- Build times need improvement

**Example**: `/agent performance-auditor analyze the workspace blueprint generation time`

### 🌍 cross-platform-tester
**Purpose**: Ensures compatibility across Windows, macOS, and Linux

**Use when**:
- Adding file system operations
- Implementing shell commands
- Working with paths or process execution
- Testing platform-specific features

**Example**: `/agent cross-platform-tester test the new file watcher feature`

### 📊 complexity-analyzer
**Purpose**: Analyzes requirements and recommends appropriate complexity levels

**Use when**:
- Users need help choosing a blueprint
- Designing new blueprint variations
- Planning complexity progression paths
- Documenting use cases

**Example**: `/agent complexity-analyzer recommend blueprint for a GraphQL API with authentication`

### 🎨 web-ui-designer
**Purpose**: Designs and implements the React-based web interface for go-starter (Phase 3/4)

**Use when**:
- Building web UI components
- Implementing progressive disclosure in web interface
- Creating live preview functionality
- Designing responsive layouts
- Optimizing web UX/UI

**Example**: `/agent web-ui-designer create the project configuration wizard component`

### 🎭 ascii-art-designer
**Purpose**: Creates professional ASCII art logos, banners, and decorative elements for CLI applications

**Use when**:
- Designing ASCII art logos for CLI tools
- Creating banner headers for documentation
- Building terminal UI elements and decorative components
- Developing themed artwork (e.g., Gopher for Go projects)
- Optimizing ASCII art for cross-platform terminal compatibility

**Example**: `/agent ascii-art-designer create a professional logo for go-starter CLI`

### 📋 product-owner
**Purpose**: Manages GitHub issues, projects, and roadmap for strategic development planning

**Use when**:
- Prioritizing the development backlog
- Creating or refining GitHub issues
- Organizing project boards and milestones
- Planning releases and sprints
- Analyzing feature requests and bug reports

**Example**: `/agent product-owner organize the Phase 3 web UI milestone and prioritize issues`

### 📈 marketing-specialist
**Purpose**: Positions go-starter as the leading AI-powered Go development platform and drives user acquisition

**Use when**:
- Developing market positioning and messaging
- Creating go-to-market strategies for AI features
- Planning content marketing and thought leadership
- Analyzing competitive landscape
- Designing user acquisition campaigns

**Example**: `/agent marketing-specialist create launch strategy for AI Template Generator feature`

### 💼 sales-specialist
**Purpose**: Converts users into paying customers and develops enterprise AI service offerings

**Use when**:
- Designing freemium conversion funnels
- Creating enterprise sales processes
- Developing pricing and packaging strategies
- Building customer success frameworks
- Planning revenue expansion strategies

**Example**: `/agent sales-specialist design enterprise sales process for AI Migration Assistant`

## Agent Collaboration

Agents can work together for comprehensive solutions:

1. **Blueprint Production Pipeline** (Enhanced Workflow):
   - `template-variable-auditor` → systematic variable analysis and fixes
   - `grpc-protobuf-specialist` → fix gRPC/protobuf specific issues
   - `blueprint-validator` → comprehensive blueprint validation
   - `golang-atdd-qa-engineer` → ATDD test creation and execution
   - `blueprint-production-pipeline` → coordinate entire workflow and status tracking

2. **gRPC/Microservice Development**:
   - `grpc-protobuf-specialist` → protobuf design and buf configuration
   - `microservice-orchestrator` → service mesh and distributed patterns
   - `template-variable-auditor` → ensure template consistency
   - `blueprint-validator` → validate generated microservice projects

3. **Serverless Architecture Development**:
   - `serverless-specialist` → AWS Lambda and event-driven patterns
   - `microservice-orchestrator` → distributed system resilience patterns
   - `template-variable-auditor` → serverless-specific template variables
   - `blueprint-production-pipeline` → production readiness validation

4. **Template System Maintenance**:
   - `template-variable-auditor` → systematic template analysis
   - `blueprint-validator` → template syntax and compilation validation
   - `cross-platform-tester` → multi-platform template compatibility
   - `blueprint-production-pipeline` → automated quality gates

5. **New Blueprint Development**:
   - `complexity-analyzer` → determine complexity level and requirements
   - `template-variable-auditor` → establish consistent variable patterns
   - `blueprint-validator` → validate implementation and compilation
   - `golang-atdd-qa-engineer` → create comprehensive acceptance tests
   - `blueprint-production-pipeline` → manage development-to-production pipeline

6. **Performance Optimization**:
   - `performance-auditor` → identify bottlenecks
   - `blueprint-validator` → ensure optimizations don't break functionality
   - `golang-atdd-qa-engineer` → add performance benchmarks

7. **User Experience Enhancement**:
   - `progressive-disclosure-optimizer` → improve flag organization
   - `complexity-analyzer` → ensure appropriate defaults
   - `cross-platform-tester` → verify consistent behavior

8. **Web Interface Development** (Phase 3/4):
   - `web-ui-designer` → create React components and UX flows
   - `ascii-art-designer` → coordinate CLI and web visual branding
   - `progressive-disclosure-optimizer` → ensure web/CLI consistency
   - `performance-auditor` → optimize web performance
   - `cross-platform-tester` → test across browsers

9. **Project Management & Planning**:
   - `product-owner` → prioritize features and manage backlog
   - `complexity-analyzer` → assess feature complexity for estimation
   - `golang-atdd-qa-engineer` → define acceptance criteria for issues
   - `blueprint-production-pipeline` → production readiness tracking

10. **Business & Revenue Growth**:
   - `marketing-specialist` → position product and drive acquisition
   - `sales-specialist` → convert users and expand revenue
   - `product-owner` → align business strategy with development
   - `complexity-analyzer` → inform pricing and packaging decisions
   - `web-ui-designer` → create conversion-optimized interfaces

## Best Practices

1. **Use agents for their specialized expertise** - Each agent has deep knowledge in their domain
2. **Provide context** - Give agents relevant information about your specific task
3. **Review agent suggestions** - Agents provide recommendations, but human judgment is important
4. **Combine agent insights** - Multiple perspectives often lead to better solutions

## Adding New Agents

To add a new agent:
1. Create a new `.md` file in this directory
2. Include front matter with name, description, and tools
3. Write a detailed system prompt explaining the agent's expertise
4. Document when and how to use the agent
5. Update this README with the new agent information

## Maintenance

These agents should be updated when:
- Project structure changes significantly
- New features are added that affect their domain
- Best practices evolve
- User feedback suggests improvements