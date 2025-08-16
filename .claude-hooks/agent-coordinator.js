/**
 * Claude Code Agent Coordinator Hook
 * Orchestrates multiple specialized agents and manages their interactions
 */

const { SPECIALIZED_AGENTS } = require('./agent-selector')
const { ContextAnalyzer } = require('./context-analyzer')

/**
 * Coordinates multiple specialized agents for complex tasks
 */
class AgentCoordinator {
  constructor(projectRoot = process.cwd()) {
    this.projectRoot = projectRoot
    this.activeAgents = new Set()
    this.taskQueue = []
    this.results = new Map()
    this.contextAnalyzer = new ContextAnalyzer(projectRoot)
  }

  /**
   * Analyze task and determine which agents should collaborate
   */
  async planTaskExecution(task) {
    const context = await this.contextAnalyzer.getProjectContext()
    const agentPlan = this.createExecutionPlan(task, context)
    
    return {
      task,
      executionPlan: agentPlan,
      estimatedDuration: this.estimateDuration(agentPlan),
      dependencies: this.analyzeDependencies(agentPlan),
      riskFactors: this.identifyRisks(agentPlan, context)
    }
  }

  /**
   * Create execution plan with agent coordination
   */
  createExecutionPlan(task, context) {
    const plan = {
      phases: [],
      parallelTasks: [],
      sequentialTasks: [],
      requiredAgents: new Set()
    }

    // Analyze task complexity and requirements
    const taskAnalysis = this.analyzeTask(task)
    
    // Production deployment coordination
    if (taskAnalysis.type === 'production-deployment') {
      plan.phases = [
        {
          name: 'Pre-deployment Validation',
          agents: ['performance', 'accessibility', 'devops'],
          parallel: true,
          tasks: [
            'Performance benchmarking',
            'Accessibility compliance audit', 
            'Infrastructure validation'
          ]
        },
        {
          name: 'Security & Quality Assurance',
          agents: ['performance'],
          parallel: false,
          tasks: [
            'Security vulnerability scan',
            'Load testing',
            'Code quality validation'
          ]
        },
        {
          name: 'Deployment Execution',
          agents: ['devops'],
          parallel: false,
          tasks: [
            'Blue-green deployment',
            'Health checks',
            'Monitoring setup'
          ]
        },
        {
          name: 'Post-deployment Activities',
          agents: ['documentation', 'performance'],
          parallel: true,
          tasks: [
            'Documentation updates',
            'Performance monitoring',
            'Community notifications'
          ]
        }
      ]
    }
    
    // Web UI optimization coordination
    else if (taskAnalysis.type === 'web-ui-optimization') {
      plan.phases = [
        {
          name: 'Analysis Phase',
          agents: ['performance', 'accessibility'],
          parallel: true,
          tasks: [
            'Performance profiling',
            'Accessibility audit',
            'Bundle analysis'
          ]
        },
        {
          name: 'Implementation Phase',
          agents: ['performance', 'accessibility'],
          parallel: false,
          tasks: [
            'Code optimization',
            'Accessibility improvements',
            'Bundle size reduction'
          ]
        },
        {
          name: 'Validation Phase',
          agents: ['performance', 'accessibility', 'documentation'],
          parallel: true,
          tasks: [
            'Performance testing',
            'Accessibility validation',
            'Documentation updates'
          ]
        }
      ]
    }
    
    // Documentation overhaul coordination
    else if (taskAnalysis.type === 'documentation-overhaul') {
      plan.phases = [
        {
          name: 'Content Audit',
          agents: ['documentation'],
          parallel: false,
          tasks: [
            'Existing content analysis',
            'Gap identification',
            'User journey mapping'
          ]
        },
        {
          name: 'Content Creation',
          agents: ['documentation', 'accessibility'],
          parallel: true,
          tasks: [
            'API documentation generation',
            'Tutorial creation',
            'Accessibility documentation'
          ]
        },
        {
          name: 'Integration & Testing',
          agents: ['documentation', 'devops'],
          parallel: true,
          tasks: [
            'Documentation site deployment',
            'Link validation',
            'Search optimization'
          ]
        }
      ]
    }
    
    // Security hardening coordination
    else if (taskAnalysis.type === 'security-hardening') {
      plan.phases = [
        {
          name: 'Security Assessment',
          agents: ['performance', 'devops'],
          parallel: true,
          tasks: [
            'Vulnerability scanning',
            'Infrastructure security audit',
            'Code security analysis'
          ]
        },
        {
          name: 'Implementation',
          agents: ['performance', 'devops'],
          parallel: false,
          tasks: [
            'Security controls implementation',
            'Infrastructure hardening',
            'Monitoring setup'
          ]
        },
        {
          name: 'Validation & Documentation',
          agents: ['performance', 'documentation'],
          parallel: true,
          tasks: [
            'Security testing',
            'Security documentation',
            'Incident response procedures'
          ]
        }
      ]
    }
    
    // Default: Single agent task
    else {
      const primaryAgent = this.selectPrimaryAgent(taskAnalysis)
      plan.phases = [
        {
          name: 'Task Execution',
          agents: [primaryAgent],
          parallel: false,
          tasks: [taskAnalysis.description]
        }
      ]
    }

    // Extract required agents
    plan.phases.forEach(phase => {
      phase.agents.forEach(agent => plan.requiredAgents.add(agent))
    })

    return plan
  }

  /**
   * Analyze task to determine type and complexity
   */
  analyzeTask(task) {
    const taskLower = task.toLowerCase()
    
    const patterns = {
      'production-deployment': /deploy|production|launch|release|go-live/,
      'web-ui-optimization': /web.*ui|interface|frontend|react|performance.*ui|ui.*performance/,
      'documentation-overhaul': /documentation|docs|api.*docs|tutorial|guide/,
      'security-hardening': /security|vulnerability|harden|penetration|audit.*security/,
      'accessibility-compliance': /accessibility|a11y|wcag|screen.*reader|compliance/,
      'performance-optimization': /performance|optimization|speed|slow|benchmark/,
      'infrastructure-setup': /infrastructure|kubernetes|docker|ci.*cd|devops/
    }

    for (const [type, pattern] of Object.entries(patterns)) {
      if (pattern.test(taskLower)) {
        return {
          type,
          description: task,
          complexity: this.estimateComplexity(type),
          priority: this.determinePriority(type)
        }
      }
    }

    return {
      type: 'general',
      description: task,
      complexity: 'medium',
      priority: 'medium'
    }
  }

  /**
   * Select primary agent for single-agent tasks
   */
  selectPrimaryAgent(taskAnalysis) {
    const agentMapping = {
      'accessibility-compliance': 'accessibility',
      'performance-optimization': 'performance',
      'infrastructure-setup': 'devops',
      'documentation-overhaul': 'documentation'
    }

    return agentMapping[taskAnalysis.type] || 'documentation'
  }

  /**
   * Estimate task duration based on complexity and agents involved
   */
  estimateDuration(plan) {
    const baseTime = {
      simple: 30,    // 30 minutes
      medium: 120,   // 2 hours
      complex: 480,  // 8 hours
      expert: 1440   // 24 hours
    }

    let totalTime = 0
    let maxParallelTime = 0

    plan.phases.forEach(phase => {
      const phaseTime = phase.tasks.length * (baseTime.medium / phase.agents.length)
      
      if (phase.parallel) {
        maxParallelTime = Math.max(maxParallelTime, phaseTime)
      } else {
        totalTime += phaseTime
      }
    })

    return {
      sequential: totalTime,
      parallel: maxParallelTime,
      total: totalTime + maxParallelTime,
      formatted: this.formatDuration(totalTime + maxParallelTime)
    }
  }

  /**
   * Analyze dependencies between agents and tasks
   */
  analyzeDependencies(plan) {
    const dependencies = []

    // Cross-phase dependencies
    for (let i = 1; i < plan.phases.length; i++) {
      const currentPhase = plan.phases[i]
      const previousPhase = plan.phases[i - 1]
      
      dependencies.push({
        type: 'phase-dependency',
        requires: previousPhase.name,
        blocks: currentPhase.name,
        reason: 'Sequential phase execution'
      })
    }

    // Agent-specific dependencies
    const agentDependencies = {
      'accessibility': {
        requires: ['performance'],
        reason: 'Performance optimization affects accessibility testing'
      },
      'devops': {
        requires: ['performance', 'accessibility'],
        reason: 'Deployment requires quality validation'
      },
      'documentation': {
        requires: ['accessibility', 'performance', 'devops'],
        reason: 'Documentation should reflect final implementations'
      }
    }

    plan.requiredAgents.forEach(agent => {
      const deps = agentDependencies[agent]
      if (deps) {
        deps.requires.forEach(requiredAgent => {
          if (plan.requiredAgents.has(requiredAgent)) {
            dependencies.push({
              type: 'agent-dependency',
              requires: requiredAgent,
              blocks: agent,
              reason: deps.reason
            })
          }
        })
      }
    })

    return dependencies
  }

  /**
   * Identify potential risks in the execution plan
   */
  identifyRisks(plan, context) {
    const risks = []

    // Resource contention risk
    const simultaneousAgents = Math.max(...plan.phases.map(p => p.agents.length))
    if (simultaneousAgents > 2) {
      risks.push({
        type: 'resource-contention',
        severity: 'medium',
        description: 'Multiple agents running simultaneously may cause resource contention',
        mitigation: 'Consider sequential execution for resource-intensive tasks'
      })
    }

    // Complexity risk
    const totalTasks = plan.phases.reduce((sum, phase) => sum + phase.tasks.length, 0)
    if (totalTasks > 10) {
      risks.push({
        type: 'complexity-risk',
        severity: 'high',
        description: 'High number of tasks increases coordination complexity',
        mitigation: 'Break down into smaller, focused phases'
      })
    }

    // Context-specific risks
    if (context.healthCheck.score < 70) {
      risks.push({
        type: 'project-health-risk',
        severity: 'medium',
        description: 'Low project health score may affect task execution',
        mitigation: 'Address health issues before complex tasks'
      })
    }

    if (context.phase.completion < 80) {
      risks.push({
        type: 'incomplete-phase-risk',
        severity: 'low',
        description: 'Current phase not fully complete',
        mitigation: 'Ensure phase completion before advanced tasks'
      })
    }

    return risks
  }

  /**
   * Execute coordinated task with multiple agents
   */
  async executeCoordinatedTask(taskPlan) {
    console.log(`🚀 Starting coordinated task execution: ${taskPlan.task}`)
    console.log(`📋 Execution plan: ${taskPlan.executionPlan.phases.length} phases`)
    
    const results = {
      startTime: Date.now(),
      phases: [],
      success: true,
      errors: []
    }

    try {
      for (const [index, phase] of taskPlan.executionPlan.phases.entries()) {
        console.log(`⚡ Phase ${index + 1}: ${phase.name}`)
        
        const phaseResult = await this.executePhase(phase)
        results.phases.push(phaseResult)
        
        if (!phaseResult.success) {
          results.success = false
          results.errors.push(`Phase ${phase.name} failed: ${phaseResult.error}`)
          
          // Decide whether to continue or abort
          if (phaseResult.critical) {
            console.log(`❌ Critical failure in phase ${phase.name}, aborting execution`)
            break
          }
        }
      }
      
      results.duration = Date.now() - results.startTime
      results.summary = this.generateExecutionSummary(results)
      
      console.log(`✅ Task execution completed in ${this.formatDuration(results.duration)}`)
      return results
      
    } catch (error) {
      console.error('❌ Coordinated task execution failed:', error)
      return {
        ...results,
        success: false,
        error: error.message,
        duration: Date.now() - results.startTime
      }
    }
  }

  /**
   * Execute a single phase with potentially multiple agents
   */
  async executePhase(phase) {
    const phaseStartTime = Date.now()
    
    try {
      if (phase.parallel) {
        // Execute agent tasks in parallel
        const agentPromises = phase.agents.map(agentKey => 
          this.executeAgentTasks(agentKey, phase.tasks)
        )
        
        const agentResults = await Promise.allSettled(agentPromises)
        
        return {
          success: agentResults.every(r => r.status === 'fulfilled'),
          agents: phase.agents,
          results: agentResults,
          duration: Date.now() - phaseStartTime,
          parallel: true
        }
      } else {
        // Execute agent tasks sequentially
        const agentResults = []
        
        for (const agentKey of phase.agents) {
          const result = await this.executeAgentTasks(agentKey, phase.tasks)
          agentResults.push(result)
          
          if (!result.success && result.critical) {
            break // Stop on critical failure
          }
        }
        
        return {
          success: agentResults.every(r => r.success),
          agents: phase.agents,
          results: agentResults,
          duration: Date.now() - phaseStartTime,
          parallel: false
        }
      }
    } catch (error) {
      return {
        success: false,
        error: error.message,
        critical: true,
        duration: Date.now() - phaseStartTime
      }
    }
  }

  /**
   * Execute tasks for a specific agent
   */
  async executeAgentTasks(agentKey, tasks) {
    const agent = SPECIALIZED_AGENTS[agentKey]
    if (!agent) {
      return {
        success: false,
        error: `Agent ${agentKey} not found`,
        critical: true
      }
    }

    console.log(`  🎯 ${agent.icon} ${agent.name}: ${tasks.length} tasks`)
    
    // Simulate agent task execution
    // In a real implementation, this would invoke the actual agent
    const taskResults = tasks.map(task => ({
      task,
      success: true, // Simulate success
      duration: Math.random() * 1000 + 500, // Random duration
      output: `${agent.name} completed: ${task}`
    }))

    return {
      success: taskResults.every(t => t.success),
      agent: agent.name,
      tasks: taskResults,
      summary: `Completed ${taskResults.length} tasks for ${agent.name}`
    }
  }

  /**
   * Generate execution summary
   */
  generateExecutionSummary(results) {
    const totalTasks = results.phases.reduce((sum, phase) => 
      sum + (phase.results?.reduce((taskSum, agent) => 
        taskSum + (agent.value?.tasks?.length || 0), 0) || 0), 0)
    
    const successfulPhases = results.phases.filter(p => p.success).length
    const totalPhases = results.phases.length

    return {
      totalPhases,
      successfulPhases,
      totalTasks,
      overallSuccess: results.success,
      duration: this.formatDuration(results.duration),
      efficiency: Math.round((successfulPhases / totalPhases) * 100)
    }
  }

  /**
   * Helper methods
   */
  estimateComplexity(taskType) {
    const complexityMap = {
      'production-deployment': 'expert',
      'web-ui-optimization': 'complex',
      'security-hardening': 'complex',
      'documentation-overhaul': 'medium',
      'accessibility-compliance': 'medium',
      'performance-optimization': 'medium'
    }
    
    return complexityMap[taskType] || 'medium'
  }

  determinePriority(taskType) {
    const priorityMap = {
      'production-deployment': 'high',
      'security-hardening': 'high',
      'accessibility-compliance': 'high',
      'web-ui-optimization': 'medium',
      'performance-optimization': 'medium',
      'documentation-overhaul': 'medium'
    }
    
    return priorityMap[taskType] || 'medium'
  }

  formatDuration(milliseconds) {
    const minutes = Math.floor(milliseconds / 60000)
    const seconds = Math.floor((milliseconds % 60000) / 1000)
    
    if (minutes > 0) {
      return `${minutes}m ${seconds}s`
    }
    return `${seconds}s`
  }
}

/**
 * Claude Code hook execution function
 */
async function executeAgentCoordination(task, options = {}) {
  console.log('🤝 Initializing agent coordination...')
  
  try {
    const coordinator = new AgentCoordinator(options.projectRoot)
    
    // Plan the task execution
    const taskPlan = await coordinator.planTaskExecution(task)
    
    // Execute if auto-execute is enabled
    let executionResults = null
    if (options.autoExecute) {
      executionResults = await coordinator.executeCoordinatedTask(taskPlan)
    }

    return {
      success: true,
      taskPlan,
      executionResults,
      recommendations: generateCoordinationRecommendations(taskPlan)
    }
  } catch (error) {
    console.error('Agent coordination error:', error)
    return {
      success: false,
      error: error.message,
      fallback: 'Consider executing agents individually'
    }
  }
}

/**
 * Generate recommendations for agent coordination
 */
function generateCoordinationRecommendations(taskPlan) {
  const recommendations = []

  if (taskPlan.riskFactors.length > 0) {
    recommendations.push({
      type: 'risk-mitigation',
      message: 'High-risk coordination detected',
      actions: taskPlan.riskFactors.map(risk => risk.mitigation)
    })
  }

  if (taskPlan.estimatedDuration.total > 240) { // > 4 hours
    recommendations.push({
      type: 'time-management',
      message: 'Long execution time estimated',
      actions: ['Consider breaking into smaller tasks', 'Schedule during low-activity periods']
    })
  }

  if (taskPlan.executionPlan.requiredAgents.size > 3) {
    recommendations.push({
      type: 'complexity-management', 
      message: 'Multiple agents required',
      actions: ['Ensure agent availability', 'Consider phased approach']
    })
  }

  return recommendations
}

module.exports = {
  AgentCoordinator,
  executeAgentCoordination,
  generateCoordinationRecommendations
}