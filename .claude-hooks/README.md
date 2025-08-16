# Claude Code Hooks for go-starter Specialized Agents

This directory contains Claude Code hooks that provide intelligent agent selection and coordination for the go-starter project. These hooks enhance the Claude Code experience by automatically suggesting the most appropriate specialized agent based on context and coordinating multi-agent workflows.

## Available Hooks

### 1. Agent Selector (`agent-selector.js`)
**Purpose**: Intelligently suggests the most appropriate specialized agent based on current context.

**Features**:
- 🎯 Context-aware agent recommendation
- 📊 Confidence scoring and alternative suggestions  
- 🔍 Trigger-based pattern matching
- 📋 Quick action recommendations

**Usage**:
```javascript
// Automatically triggered when Claude analyzes context
const suggestion = suggestAgent({
  currentFile: 'web/src/components/Form.tsx',
  userQuery: 'make this accessible',
  recentFiles: ['web/src/App.tsx', 'web/src/index.tsx']
})
```

**Trigger Patterns**:
- **Accessibility Agent**: `accessibility`, `a11y`, `wcag`, `web/src`, `components`
- **DevOps Agent**: `deployment`, `docker`, `kubernetes`, `infrastructure`, `production`
- **Performance Agent**: `performance`, `optimization`, `security`, `bundle`, `benchmark`
- **Documentation Agent**: `documentation`, `docs`, `readme`, `tutorial`, `guide`

### 2. Context Analyzer (`context-analyzer.js`)
**Purpose**: Provides comprehensive project analysis for informed agent selection.

**Features**:
- 🏗️ Project phase detection (Phase 1-4)
- 🔧 Technology stack analysis
- 📊 Project health assessment
- 📈 Priority calculation
- 🎯 Architecture pattern detection

**Analysis Areas**:
```javascript
const context = await analyzer.getProjectContext()
// Returns:
// - phase: Current project phase and completion
// - technologies: Detected frameworks and tools
// - structure: Architecture and code complexity
// - recentActivity: Git activity analysis
// - healthCheck: Project health score
// - priorities: Recommended focus areas
```

**Health Checks**:
- ✅ Go modules configuration
- ✅ Documentation completeness
- ✅ Test coverage
- ✅ CI/CD setup
- ✅ Security configuration
- ✅ Specialized agent availability

### 3. Agent Coordinator (`agent-coordinator.js`)
**Purpose**: Orchestrates complex tasks requiring multiple specialized agents.

**Features**:
- 🤝 Multi-agent workflow planning
- ⚡ Parallel and sequential task execution
- 🔗 Dependency management
- ⚠️ Risk assessment
- 📊 Progress tracking

**Coordination Scenarios**:

#### Production Deployment
```javascript
// Coordinates: Performance + Accessibility + DevOps + Documentation
const plan = await coordinator.planTaskExecution('deploy to production')
```

**Execution Phases**:
1. **Pre-deployment Validation** (Parallel)
   - Performance benchmarking
   - Accessibility compliance audit
   - Infrastructure validation

2. **Security & Quality Assurance** (Sequential)
   - Security vulnerability scan
   - Load testing
   - Code quality validation

3. **Deployment Execution** (Sequential)
   - Blue-green deployment
   - Health checks
   - Monitoring setup

4. **Post-deployment Activities** (Parallel)
   - Documentation updates
   - Performance monitoring
   - Community notifications

#### Web UI Optimization
```javascript
// Coordinates: Performance + Accessibility
const plan = await coordinator.planTaskExecution('optimize web interface')
```

## Integration with Claude Code

### Automatic Hook Execution

The hooks are automatically executed by Claude Code when:

1. **File Context Changes**: When switching between files
2. **Query Analysis**: When processing user requests
3. **Project Navigation**: When exploring project structure
4. **Task Planning**: When complex tasks are identified

### Manual Hook Invocation

You can also manually invoke hooks:

```bash
# In Claude Code terminal
node .claude-hooks/agent-selector.js
node .claude-hooks/context-analyzer.js
node .claude-hooks/agent-coordinator.js
```

### Hook Output Integration

The hooks provide rich output that Claude Code integrates into responses:

```markdown
## 🎯 Recommended Specialized Agent

⚡ **DevOps & Production Deployment Agent** (85% confidence)

**Why this agent?** Detected infrastructure deployment triggers: docker, kubernetes, production

**Agent Capabilities:**
Multi-cloud deployment, infrastructure as code, CI/CD automation

**Quick Actions:**
- Load DevOps & Production Deployment Agent
- Switch to devops mode
- Run devops analysis
```

## Hook Configuration

### Environment Variables

```bash
# Project root (auto-detected)
CLAUDE_PROJECT_ROOT=/path/to/go-starter

# Hook execution timeout
CLAUDE_HOOK_TIMEOUT=30000

# Cache duration for context analysis
CLAUDE_CONTEXT_CACHE_TTL=60000

# Enable debug logging
CLAUDE_HOOKS_DEBUG=true
```

### Customization

Each hook supports customization through options:

```javascript
// Agent Selector customization
const customAgents = {
  myAgent: {
    file: 'MY_CUSTOM_AGENT.md',
    name: 'Custom Agent',
    triggers: ['custom', 'special'],
    icon: '🔧'
  }
}

// Context Analyzer customization
const analyzer = new ContextAnalyzer(projectRoot, {
  cacheTimeout: 120000,
  includeGitAnalysis: true,
  healthChecks: customHealthChecks
})

// Agent Coordinator customization
const coordinator = new AgentCoordinator(projectRoot, {
  maxParallelAgents: 3,
  defaultTimeout: 30000,
  riskToleranceLevel: 'medium'
})
```

## Performance Considerations

### Caching Strategy
- **Context Analysis**: 1-minute cache for project context
- **Agent Suggestions**: Real-time processing (fast)
- **Coordination Plans**: Cached based on task similarity

### Resource Usage
- **Memory**: ~10MB for full context analysis
- **CPU**: Lightweight processing, mostly I/O bound
- **Disk**: Minimal temporary file usage

### Optimization Features
- Lazy loading of expensive operations
- Parallel execution where possible
- Intelligent caching with TTL
- Graceful degradation on errors

## Error Handling

### Graceful Degradation
All hooks implement graceful degradation:

```javascript
// If hook fails, provide fallback response
if (!hookResult.success) {
  return {
    fallback: 'All specialized agents are available for manual selection',
    error: hookResult.error,
    agents: Object.keys(SPECIALIZED_AGENTS)
  }
}
```

### Common Error Scenarios
1. **Git Not Available**: Falls back to file-based analysis
2. **Node Modules Missing**: Provides basic functionality
3. **File System Permissions**: Uses read-only operations
4. **Network Issues**: Uses local analysis only

## Development and Testing

### Running Tests
```bash
# Test individual hooks
npm test agent-selector
npm test context-analyzer  
npm test agent-coordinator

# Integration tests
npm test hooks-integration

# Performance benchmarks
npm run benchmark:hooks
```

### Debug Mode
```bash
# Enable debug logging
export CLAUDE_HOOKS_DEBUG=true

# Run with verbose output
node .claude-hooks/agent-selector.js --debug

# Trace execution
node .claude-hooks/context-analyzer.js --trace
```

### Hook Development Guidelines

1. **Idempotent Operations**: Hooks should be safe to run multiple times
2. **Fast Execution**: Target <1 second execution time
3. **Error Resilience**: Handle missing dependencies gracefully
4. **Rich Output**: Provide actionable recommendations
5. **Context Awareness**: Use all available project information

## Examples

### Example 1: Accessibility Task Detection
```javascript
// Context
const context = {
  currentFile: 'web/src/components/LoginForm.tsx',
  userQuery: 'add keyboard navigation support',
  recentFiles: ['web/src/styles/forms.css']
}

// Hook Output
{
  suggested: {
    key: 'accessibility',
    name: 'Accessibility & UX Compliance Agent',
    icon: '🎯'
  },
  confidence: 92,
  matchedTriggers: ['keyboard navigation', 'web/src', 'components'],
  reason: 'Detected 3 relevant triggers: keyboard navigation, web/src, components'
}
```

### Example 2: Production Deployment Coordination
```javascript
// Task Planning
const task = 'prepare for production launch'
const plan = await coordinator.planTaskExecution(task)

// Generated Plan
{
  phases: [
    {
      name: 'Pre-deployment Validation',
      agents: ['performance', 'accessibility', 'devops'],
      parallel: true,
      tasks: ['Performance benchmarking', 'WCAG compliance', 'Infrastructure check']
    }
    // ... additional phases
  ],
  estimatedDuration: { total: 480, formatted: '8h 0m' },
  riskFactors: [
    {
      type: 'complexity-risk',
      severity: 'medium',
      mitigation: 'Break down into focused phases'
    }
  ]
}
```

## Agent Integration Matrix

| Hook | Accessibility Agent | DevOps Agent | Performance Agent | Documentation Agent |
|------|-------------------|--------------|------------------|-------------------|
| **Agent Selector** | ✅ Pattern matching | ✅ Pattern matching | ✅ Pattern matching | ✅ Pattern matching |
| **Context Analyzer** | ✅ WCAG compliance check | ✅ Infrastructure analysis | ✅ Performance profiling | ✅ Content audit |
| **Agent Coordinator** | ✅ Multi-phase coordination | ✅ Deployment orchestration | ✅ Optimization workflows | ✅ Content strategy |

## Future Enhancements

### Planned Features
- 🤖 Machine learning for improved agent selection
- 📊 Historical performance tracking
- 🔄 Workflow templates and presets
- 🌐 Remote agent execution
- 📱 Mobile-optimized hook interfaces

### Integration Roadmap
- **VS Code Extension**: Direct integration with VS Code Claude extension
- **GitHub Actions**: Automated hook execution in CI/CD
- **Slack Bot**: Agent coordination through Slack commands
- **Web Dashboard**: Visual hook management interface

---

These Claude Code hooks transform the go-starter development experience by providing intelligent, context-aware agent assistance that adapts to your current focus and project needs.