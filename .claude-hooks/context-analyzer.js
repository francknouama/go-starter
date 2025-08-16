/**
 * Claude Code Context Analyzer Hook
 * Provides intelligent context analysis for specialized agent activation
 */

const fs = require('fs')
const path = require('path')

/**
 * Analyzes the current project context to provide rich information for agent selection
 */
class ContextAnalyzer {
  constructor(projectRoot = process.cwd()) {
    this.projectRoot = projectRoot
    this.cache = new Map()
    this.cacheTimeout = 60000 // 1 minute cache
  }

  /**
   * Get comprehensive project context
   */
  async getProjectContext() {
    const cacheKey = 'project-context'
    const cached = this.cache.get(cacheKey)
    
    if (cached && Date.now() - cached.timestamp < this.cacheTimeout) {
      return cached.data
    }

    try {
      const context = {
        phase: await this.detectProjectPhase(),
        technologies: await this.detectTechnologies(),
        structure: await this.analyzeProjectStructure(),
        recentActivity: await this.analyzeRecentActivity(),
        healthCheck: await this.performHealthCheck(),
        priorities: await this.calculatePriorities()
      }

      this.cache.set(cacheKey, {
        data: context,
        timestamp: Date.now()
      })

      return context
    } catch (error) {
      console.error('Context analysis error:', error)
      return this.getMinimalContext()
    }
  }

  /**
   * Detect current project phase based on code and documentation
   */
  async detectProjectPhase() {
    const indicators = {
      'Phase 1': ['cmd/new.go', 'blueprints/web-api/', 'blueprints/cli/'],
      'Phase 2': ['blueprints/web-api-clean/', 'blueprints/web-api-ddd/', 'tests/acceptance/'],
      'Phase 3': ['web/src/', 'cmd/web-server/', 'web/package.json'],
      'Phase 4': ['marketplace/', 'auth/', 'github-integration/']
    }

    const detectedPhases = []
    
    for (const [phase, files] of Object.entries(indicators)) {
      const phaseScore = files.reduce((score, file) => {
        return score + (this.fileExists(file) ? 1 : 0)
      }, 0)
      
      if (phaseScore > 0) {
        detectedPhases.push({
          phase,
          score: phaseScore,
          completion: (phaseScore / files.length) * 100
        })
      }
    }

    // Sort by score and return current phase
    detectedPhases.sort((a, b) => b.score - a.score)
    
    const currentPhase = detectedPhases[0]
    return {
      current: currentPhase?.phase || 'Phase 1',
      completion: currentPhase?.completion || 0,
      evidence: detectedPhases
    }
  }

  /**
   * Detect technologies and frameworks in use
   */
  async detectTechnologies() {
    const techIndicators = {
      // Backend
      'Go': ['go.mod', 'main.go', '*.go'],
      'Gin Framework': ['gin-gonic', 'gin.New()'],
      'Cobra CLI': ['cobra', 'cmd/root.go'],
      
      // Frontend
      'React': ['web/package.json', 'react'],
      'TypeScript': ['tsconfig.json', '.tsx', '.ts'],
      'Vite': ['vite.config', 'web/vite.config'],
      'Tailwind CSS': ['tailwind.config', '@tailwind'],
      
      // Infrastructure
      'Docker': ['Dockerfile', 'docker-compose'],
      'Kubernetes': ['kubernetes/', '*.yaml'],
      'GitHub Actions': ['.github/workflows/'],
      
      // Testing
      'Go Testing': ['*_test.go', 'testing'],
      'Playwright': ['playwright.config', '@playwright'],
      'Jest': ['jest.config', 'jest'],
      
      // Specialized Agents
      'Accessibility Agent': ['ACCESSIBILITY_UX_AGENT.md'],
      'DevOps Agent': ['DEVOPS_DEPLOYMENT_AGENT.md'],
      'Performance Agent': ['PERFORMANCE_SECURITY_AGENT.md'],
      'Documentation Agent': ['DOCUMENTATION_COMMUNITY_AGENT.md']
    }

    const detectedTech = {}
    
    for (const [tech, indicators] of Object.entries(techIndicators)) {
      const detected = indicators.some(indicator => {
        if (indicator.includes('*')) {
          return this.globSearch(indicator)
        }
        return this.fileExists(indicator) || this.contentSearch(indicator)
      })
      
      if (detected) {
        detectedTech[tech] = true
      }
    }

    return detectedTech
  }

  /**
   * Analyze project structure and architecture
   */
  async analyzeProjectStructure() {
    const structure = {
      architecture: 'unknown',
      directories: [],
      keyFiles: [],
      codeComplexity: 'medium'
    }

    try {
      // Detect architecture pattern
      if (this.fileExists('internal/domain/') && this.fileExists('internal/application/')) {
        structure.architecture = 'clean'
      } else if (this.fileExists('internal/entities/') && this.fileExists('internal/repositories/')) {
        structure.architecture = 'ddd'
      } else if (this.fileExists('internal/ports/') && this.fileExists('internal/adapters/')) {
        structure.architecture = 'hexagonal'
      } else if (this.fileExists('internal/')) {
        structure.architecture = 'standard'
      }

      // Count directories and files
      structure.directories = this.getDirectoryList()
      structure.keyFiles = this.getKeyFiles()
      
      // Estimate complexity
      const goFiles = this.globSearch('**/*.go').length
      const jsFiles = this.globSearch('**/*.{js,ts,tsx}').length
      const totalFiles = goFiles + jsFiles
      
      if (totalFiles < 50) structure.codeComplexity = 'simple'
      else if (totalFiles < 200) structure.codeComplexity = 'medium'
      else structure.codeComplexity = 'complex'

    } catch (error) {
      console.warn('Structure analysis failed:', error)
    }

    return structure
  }

  /**
   * Analyze recent git activity to understand current focus
   */
  async analyzeRecentActivity() {
    try {
      const { execSync } = require('child_process')
      
      // Get recent commits
      const recentCommits = execSync('git log --oneline -10', { encoding: 'utf8' })
        .split('\n')
        .filter(line => line.trim())
        .map(line => {
          const [hash, ...messageParts] = line.split(' ')
          return {
            hash: hash,
            message: messageParts.join(' ')
          }
        })

      // Get modified files
      const modifiedFiles = execSync('git status --porcelain', { encoding: 'utf8' })
        .split('\n')
        .filter(line => line.trim())
        .map(line => line.substring(3))

      // Categorize activity
      const activityPatterns = {
        'web-ui': /web\/|frontend|react|ui|interface/i,
        'backend': /internal\/|cmd\/|api|server/i,
        'testing': /test|spec|__tests__|\.test\./i,
        'docs': /docs\/|readme|documentation|\.md$/i,
        'infrastructure': /docker|kubernetes|k8s|deploy|ci|cd/i,
        'security': /security|auth|jwt|encryption/i,
        'performance': /performance|optimization|speed|benchmark/i
      }

      const activityFocus = {}
      
      // Analyze commit messages and modified files
      const allActivity = [
        ...recentCommits.map(c => c.message),
        ...modifiedFiles
      ].join(' ')

      Object.entries(activityPatterns).forEach(([category, pattern]) => {
        const matches = allActivity.match(pattern) || []
        activityFocus[category] = matches.length
      })

      return {
        recentCommits: recentCommits.slice(0, 5),
        modifiedFiles: modifiedFiles.slice(0, 10),
        activityFocus,
        primaryFocus: Object.entries(activityFocus)
          .sort(([,a], [,b]) => b - a)[0]?.[0] || 'general'
      }
    } catch (error) {
      return {
        recentCommits: [],
        modifiedFiles: [],
        activityFocus: {},
        primaryFocus: 'general',
        error: 'Git analysis failed'
      }
    }
  }

  /**
   * Perform basic health check of the project
   */
  async performHealthCheck() {
    const health = {
      score: 0,
      issues: [],
      strengths: [],
      recommendations: []
    }

    const checks = [
      {
        name: 'Go modules',
        check: () => this.fileExists('go.mod'),
        weight: 10,
        issue: 'go.mod missing'
      },
      {
        name: 'README documentation',
        check: () => this.fileExists('README.md'),
        weight: 5,
        issue: 'README.md missing'
      },
      {
        name: 'Tests present',
        check: () => this.globSearch('**/*_test.go').length > 0,
        weight: 15,
        issue: 'No test files found'
      },
      {
        name: 'CI/CD configuration',
        check: () => this.fileExists('.github/workflows/'),
        weight: 10,
        issue: 'No CI/CD workflows found'
      },
      {
        name: 'Docker configuration',
        check: () => this.fileExists('Dockerfile'),
        weight: 8,
        issue: 'Dockerfile missing'
      },
      {
        name: 'Security file',
        check: () => this.fileExists('SECURITY.md'),
        weight: 5,
        issue: 'SECURITY.md missing'
      },
      {
        name: 'Contributing guide',
        check: () => this.fileExists('CONTRIBUTING.md'),
        weight: 5,
        issue: 'CONTRIBUTING.md missing'
      },
      {
        name: 'Specialized agents',
        check: () => this.globSearch('*_AGENT.md').length >= 3,
        weight: 12,
        issue: 'Specialized agents not fully configured'
      }
    ]

    let totalScore = 0
    let maxScore = 0

    checks.forEach(({ name, check, weight, issue }) => {
      maxScore += weight
      if (check()) {
        totalScore += weight
        health.strengths.push(name)
      } else {
        health.issues.push(issue)
      }
    })

    health.score = Math.round((totalScore / maxScore) * 100)

    // Generate recommendations based on score
    if (health.score < 60) {
      health.recommendations.push('Focus on basic project setup and documentation')
    } else if (health.score < 80) {
      health.recommendations.push('Improve testing and CI/CD automation')
    } else {
      health.recommendations.push('Consider advanced optimizations and community features')
    }

    return health
  }

  /**
   * Calculate current priorities based on project state
   */
  async calculatePriorities() {
    const phase = await this.detectProjectPhase()
    const tech = await this.detectTechnologies()
    const activity = await this.analyzeRecentActivity()
    const health = await this.performHealthCheck()

    const priorities = []

    // Phase-based priorities
    if (phase.current === 'Phase 3' && tech['React']) {
      priorities.push({
        priority: 'high',
        agent: 'accessibility',
        reason: 'Phase 3 Web UI requires accessibility compliance',
        action: 'Run WCAG 2.1 AA compliance audit'
      })
      
      priorities.push({
        priority: 'high',
        agent: 'performance',
        reason: 'Web UI needs performance optimization',
        action: 'Optimize bundle size and loading times'
      })
    }

    // Health-based priorities
    if (health.score < 70) {
      priorities.push({
        priority: 'medium',
        agent: 'documentation',
        reason: 'Project health score indicates documentation gaps',
        action: 'Improve documentation completeness'
      })
    }

    // Activity-based priorities
    if (activity.primaryFocus === 'infrastructure') {
      priorities.push({
        priority: 'high',
        agent: 'devops',
        reason: 'Recent infrastructure activity detected',
        action: 'Review deployment and infrastructure changes'
      })
    }

    // Production readiness check
    if (tech['Docker'] && tech['Kubernetes'] && health.score > 80) {
      priorities.push({
        priority: 'high',
        agent: 'devops',
        reason: 'Project appears ready for production deployment',
        action: 'Execute production deployment checklist'
      })
    }

    return priorities.sort((a, b) => {
      const priorityOrder = { high: 3, medium: 2, low: 1 }
      return priorityOrder[b.priority] - priorityOrder[a.priority]
    })
  }

  /**
   * Helper methods
   */
  fileExists(filePath) {
    try {
      return fs.existsSync(path.join(this.projectRoot, filePath))
    } catch {
      return false
    }
  }

  contentSearch(searchTerm) {
    try {
      const { execSync } = require('child_process')
      const result = execSync(`grep -r "${searchTerm}" . --exclude-dir=node_modules --exclude-dir=.git`, 
        { encoding: 'utf8', cwd: this.projectRoot })
      return result.length > 0
    } catch {
      return false
    }
  }

  globSearch(pattern) {
    try {
      const glob = require('glob')
      return glob.sync(pattern, { cwd: this.projectRoot })
    } catch {
      return []
    }
  }

  getDirectoryList() {
    try {
      return fs.readdirSync(this.projectRoot, { withFileTypes: true })
        .filter(dirent => dirent.isDirectory())
        .map(dirent => dirent.name)
        .filter(name => !name.startsWith('.') && name !== 'node_modules')
    } catch {
      return []
    }
  }

  getKeyFiles() {
    const keyPatterns = [
      'go.mod', 'package.json', 'Dockerfile', 'README.md',
      'main.go', 'cmd/**/*.go', 'internal/**/*.go',
      'web/src/**/*.tsx', 'web/src/**/*.ts'
    ]

    const keyFiles = []
    keyPatterns.forEach(pattern => {
      const matches = this.globSearch(pattern)
      keyFiles.push(...matches.slice(0, 3)) // Limit per pattern
    })

    return [...new Set(keyFiles)].slice(0, 20) // Remove duplicates and limit total
  }

  getMinimalContext() {
    return {
      phase: { current: 'Unknown', completion: 0 },
      technologies: {},
      structure: { architecture: 'unknown', directories: [], keyFiles: [] },
      recentActivity: { primaryFocus: 'general' },
      healthCheck: { score: 50, issues: ['Context analysis failed'] },
      priorities: []
    }
  }
}

/**
 * Claude Code hook execution function
 */
async function executeContextAnalyzer(options = {}) {
  console.log('🔍 Analyzing project context for specialized agents...')
  
  try {
    const analyzer = new ContextAnalyzer(options.projectRoot)
    const context = await analyzer.getProjectContext()
    
    const summary = {
      phase: context.phase.current,
      completion: `${context.phase.completion}%`,
      architecture: context.structure.architecture,
      primaryFocus: context.recentActivity.primaryFocus,
      healthScore: `${context.healthCheck.score}%`,
      topPriorities: context.priorities.slice(0, 3),
      availableAgents: Object.keys(context.technologies)
        .filter(tech => tech.includes('Agent'))
        .length
    }

    console.log('Context analysis completed:', summary)

    return {
      success: true,
      context,
      summary,
      recommendations: generateContextRecommendations(context)
    }
  } catch (error) {
    console.error('Context analyzer error:', error)
    return {
      success: false,
      error: error.message,
      context: new ContextAnalyzer().getMinimalContext()
    }
  }
}

/**
 * Generate actionable recommendations based on context
 */
function generateContextRecommendations(context) {
  const recommendations = []

  // Phase-specific recommendations
  if (context.phase.current === 'Phase 3' && context.phase.completion < 100) {
    recommendations.push({
      type: 'phase-completion',
      priority: 'high',
      message: 'Complete Phase 3 Web UI development',
      agents: ['accessibility', 'performance'],
      actions: ['WCAG compliance audit', 'Performance optimization']
    })
  }

  // Health-based recommendations
  if (context.healthCheck.score < 70) {
    recommendations.push({
      type: 'health-improvement',
      priority: 'medium',
      message: 'Improve project health score',
      agents: ['documentation'],
      actions: context.healthCheck.issues.slice(0, 3)
    })
  }

  // Priority-based recommendations
  context.priorities.slice(0, 2).forEach(priority => {
    recommendations.push({
      type: 'priority-action',
      priority: priority.priority,
      message: priority.reason,
      agents: [priority.agent],
      actions: [priority.action]
    })
  })

  return recommendations
}

module.exports = {
  ContextAnalyzer,
  executeContextAnalyzer,
  generateContextRecommendations
}