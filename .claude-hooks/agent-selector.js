/**
 * Claude Code Agent Selector Hook
 * Automatically suggests the most appropriate specialized agent based on the current task context
 */

const SPECIALIZED_AGENTS = {
  accessibility: {
    file: 'ACCESSIBILITY_UX_AGENT.md',
    name: 'Accessibility & UX Compliance Agent',
    triggers: [
      'accessibility', 'a11y', 'wcag', 'screen reader', 'color contrast',
      'keyboard navigation', 'aria', 'mobile responsiveness', 'usability',
      'web/src', 'components', 'forms', 'ui', 'ux'
    ],
    icon: '🎯',
    description: 'WCAG 2.1 AA compliance, usability testing, progressive enhancement'
  },
  
  devops: {
    file: 'DEVOPS_DEPLOYMENT_AGENT.md',
    name: 'DevOps & Production Deployment Agent',
    triggers: [
      'deployment', 'docker', 'kubernetes', 'k8s', 'ci/cd', 'terraform',
      'infrastructure', 'monitoring', 'prometheus', 'grafana', 'helm',
      'production', 'staging', 'pipeline', 'automation', 'cloud'
    ],
    icon: '⚡',
    description: 'Multi-cloud deployment, infrastructure as code, CI/CD automation'
  },
  
  performance: {
    file: 'PERFORMANCE_SECURITY_AGENT.md',
    name: 'Performance & Security Optimization Agent',
    triggers: [
      'performance', 'security', 'optimization', 'speed', 'bundle',
      'lighthouse', 'benchmark', 'profiling', 'vulnerability', 'csrf',
      'xss', 'auth', 'jwt', 'encryption', 'scan', 'load test'
    ],
    icon: '🚀',
    description: 'Sub-2 second loading, security hardening, vulnerability scanning'
  },
  
  documentation: {
    file: 'DOCUMENTATION_COMMUNITY_AGENT.md',
    name: 'Documentation & Community Agent',
    triggers: [
      'documentation', 'docs', 'readme', 'api docs', 'tutorial',
      'community', 'blog', 'video', 'guide', 'examples', 'faq',
      'contributing', 'onboarding', 'changelog', 'release notes'
    ],
    icon: '📚',
    description: 'Technical writing, interactive tutorials, community engagement'
  }
}

/**
 * Analyzes the current context and suggests the most appropriate agent
 * @param {Object} context - Current Claude Code context
 * @returns {Object} Suggested agent and confidence score
 */
function suggestAgent(context) {
  const {
    currentFile = '',
    recentFiles = [],
    userQuery = '',
    gitStatus = '',
    projectStructure = ''
  } = context
  
  const analysisText = [
    currentFile,
    recentFiles.join(' '),
    userQuery,
    gitStatus,
    projectStructure
  ].join(' ').toLowerCase()
  
  const agentScores = {}
  
  // Calculate relevance scores for each agent
  Object.entries(SPECIALIZED_AGENTS).forEach(([key, agent]) => {
    let score = 0
    let matchedTriggers = []
    
    agent.triggers.forEach(trigger => {
      const triggerCount = (analysisText.match(new RegExp(trigger, 'g')) || []).length
      if (triggerCount > 0) {
        score += triggerCount * (trigger.length > 3 ? 2 : 1) // Longer triggers get higher weight
        matchedTriggers.push(trigger)
      }
    })
    
    // Bonus for file path relevance
    if (agent.triggers.some(trigger => currentFile.includes(trigger))) {
      score += 5
    }
    
    // Context-specific bonuses
    if (key === 'accessibility' && (currentFile.includes('web/') || currentFile.includes('components'))) {
      score += 3
    }
    
    if (key === 'devops' && (currentFile.includes('docker') || currentFile.includes('.yml') || currentFile.includes('k8s'))) {
      score += 3
    }
    
    if (key === 'performance' && (currentFile.includes('test') || userQuery.includes('optimize') || userQuery.includes('slow'))) {
      score += 3
    }
    
    if (key === 'documentation' && (currentFile.includes('docs/') || currentFile.includes('README') || currentFile.includes('.md'))) {
      score += 3
    }
    
    agentScores[key] = {
      score,
      matchedTriggers,
      agent
    }
  })
  
  // Find the highest scoring agent
  const sortedAgents = Object.entries(agentScores)
    .filter(([_, data]) => data.score > 0)
    .sort(([_, a], [__, b]) => b.score - a.score)
  
  if (sortedAgents.length === 0) {
    return {
      suggested: null,
      confidence: 0,
      reason: 'No specialized agent triggers detected'
    }
  }
  
  const [topAgentKey, topAgentData] = sortedAgents[0]
  const confidence = Math.min(100, (topAgentData.score / 10) * 100)
  
  return {
    suggested: {
      key: topAgentKey,
      ...topAgentData.agent
    },
    confidence: Math.round(confidence),
    matchedTriggers: topAgentData.matchedTriggers,
    alternatives: sortedAgents.slice(1, 3).map(([key, data]) => ({
      key,
      name: data.agent.name,
      icon: data.agent.icon,
      score: data.score
    })),
    reason: `Detected ${topAgentData.matchedTriggers.length} relevant triggers: ${topAgentData.matchedTriggers.slice(0, 3).join(', ')}`
  }
}

/**
 * Generates a formatted agent suggestion for Claude
 * @param {Object} context - Current context
 * @returns {String} Formatted suggestion
 */
function generateAgentSuggestion(context) {
  const suggestion = suggestAgent(context)
  
  if (!suggestion.suggested) {
    return `
## 🤖 No Specialized Agent Required

The current task appears to be general development work. All four specialized agents are available if needed:

${Object.entries(SPECIALIZED_AGENTS).map(([key, agent]) => 
  `- ${agent.icon} **${agent.name}**: ${agent.description}`
).join('\n')}

To activate a specific agent, mention relevant keywords or specify the agent explicitly.
`
  }
  
  const { suggested, confidence, matchedTriggers, alternatives, reason } = suggestion
  
  return `
## 🎯 Recommended Specialized Agent

${suggested.icon} **${suggested.name}** (${confidence}% confidence)

**Why this agent?** ${reason}

**Agent Capabilities:**
${suggested.description}

**Matched Context:** ${matchedTriggers.slice(0, 5).join(', ')}

### Quick Actions:
- \`Load ${suggested.name}\` - Activate this specialized agent
- \`Switch to ${suggested.key} mode\` - Focus on ${suggested.key} concerns
- \`Run ${suggested.key} analysis\` - Perform specialized analysis

${alternatives.length > 0 ? `
### Alternative Agents:
${alternatives.map(alt => `- ${alt.icon} ${alt.name} (${alt.score} relevance)`).join('\n')}
` : ''}

### All Available Agents:
${Object.entries(SPECIALIZED_AGENTS).map(([key, agent]) => 
  `${key === suggested.key ? '**' : ''}${agent.icon} ${agent.name}${key === suggested.key ? '** ← Recommended' : ''}`
).join('\n')}
`
}

/**
 * Hook execution function called by Claude Code
 */
function executeAgentSelectorHook(context) {
  console.log('🔍 Analyzing context for specialized agent suggestion...')
  
  try {
    const suggestion = generateAgentSuggestion(context)
    
    // Log for debugging
    console.log('Agent analysis completed:', {
      currentFile: context.currentFile,
      suggestedAgent: suggestAgent(context).suggested?.name || 'None',
      confidence: suggestAgent(context).confidence
    })
    
    return {
      success: true,
      suggestion,
      metadata: {
        timestamp: new Date().toISOString(),
        context: {
          file: context.currentFile,
          hasQuery: !!context.userQuery,
          filesAnalyzed: context.recentFiles?.length || 0
        }
      }
    }
  } catch (error) {
    console.error('Agent selector hook error:', error)
    return {
      success: false,
      error: error.message,
      fallback: 'All specialized agents are available. Use the file list to select the appropriate agent manually.'
    }
  }
}

// Export for Claude Code integration
module.exports = {
  SPECIALIZED_AGENTS,
  suggestAgent,
  generateAgentSuggestion,
  executeAgentSelectorHook
}

// Example usage contexts for testing:
/*
Test contexts:
1. Accessibility: { currentFile: 'web/src/components/Form.tsx', userQuery: 'make this accessible' }
2. DevOps: { currentFile: 'Dockerfile', userQuery: 'deploy to production' }
3. Performance: { currentFile: 'web/src/App.tsx', userQuery: 'this is loading slowly' }
4. Documentation: { currentFile: 'docs/api.md', userQuery: 'update the documentation' }
*/