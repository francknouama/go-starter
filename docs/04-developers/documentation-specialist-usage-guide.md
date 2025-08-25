# Documentation Specialist Agent Usage Guide

## When to Use the Documentation Specialist Agent

The `documentation-specialist` agent is designed for comprehensive documentation creation, architecture, and maintenance tasks that require technical writing excellence and multi-audience content design.

## Primary Use Cases

### ✅ Use documentation-specialist for:

#### 1. Comprehensive Documentation Creation
- **New User Guides**: Complete tutorials, getting-started content, onboarding experiences
- **Technical Reference Materials**: API documentation, configuration guides, troubleshooting guides
- **Architecture Documentation**: System design docs, decision records, technical overviews
- **Feature Documentation**: Comprehensive guides for new features or enhancements

#### 2. Information Architecture & Content Strategy
- **Documentation Restructuring**: Reorganizing content for better user experience
- **Navigation Design**: Creating logical content hierarchies and discovery paths
- **Content Gap Analysis**: Identifying missing documentation and planning content roadmaps
- **Multi-Audience Content Planning**: Designing progressive content for different user levels

#### 3. Major Documentation Updates
- **Blueprint Enhancement Documentation**: Comprehensive guides for Phase 2 features
- **Migration Guides**: Moving between blueprint types or complexity levels
- **Best Practices Documentation**: Community wisdom compilation and organization
- **Enterprise Documentation**: Advanced configuration, team setup, organizational guides

#### 4. Quality Assurance & Consistency Projects
- **Style Guide Implementation**: Establishing and enforcing documentation standards
- **Cross-Document Consistency**: Major consistency reviews and standardization
- **Content Accuracy Validation**: Comprehensive reviews against actual functionality
- **Documentation Architecture Reviews**: Evaluating and improving overall structure

### ❌ Don't use documentation-specialist for:

#### 1. Simple Status Updates
- **Use instead**: `documentation-sync-manager` for README updates and status synchronization
- **Example**: Updating blueprint counts or production status across documents

#### 2. Quick Fixes and Corrections
- **Use instead**: Direct editing for typos, small corrections, or minor updates
- **Example**: Fixing a typo, updating a single command example, correcting a link

#### 3. Issue Tracking and Project Management
- **Use instead**: `github-project-manager` for GitHub issues, project boards, milestones
- **Example**: Creating issues, updating project status, coordinating releases

#### 4. Technical Implementation
- **Use instead**: `golang-fullstack-engineer` for code changes and technical implementation
- **Example**: Fixing template code, updating blueprint logic, implementing new features

## Agent Coordination Workflows

### Documentation Enhancement Workflow
```
User Request → documentation-specialist → Comprehensive Content Creation
                                    ↓
blueprint-production-tracker → Status Data & Validation Results
                                    ↓
documentation-sync-manager → Cross-Document Consistency & Updates
                                    ↓
github-project-manager → Release Coordination & Milestone Updates
```

### Major Feature Documentation Workflow
```
Feature Release → documentation-specialist → Documentation Architecture Design
                                      ↓
golang-fullstack-engineer → Technical Details & Implementation Examples
                                      ↓
performance-security-specialist → Best Practices & Security Guidelines
                                      ↓
documentation-sync-manager → Multi-Document Synchronization
```

## Examples of Appropriate Tasks

### ✅ Perfect for documentation-specialist:

#### "Create comprehensive onboarding guide for enterprise users"
**Why**: Requires technical writing excellence, multi-audience design, and information architecture
**Scope**: Complete guide creation with progressive complexity and multiple use cases

#### "Document Phase 2 enhanced blueprint features comprehensively"
**Why**: Complex technical content requiring structured documentation architecture
**Scope**: Feature overview, implementation guides, best practices, troubleshooting

#### "Redesign blueprint selection guide for better user experience"
**Why**: Information architecture, user journey optimization, and comprehensive content strategy
**Scope**: Content restructuring, decision tree design, enhanced navigation

#### "Create migration guide between blueprint complexity levels"
**Why**: Technical writing for complex workflows, multi-scenario documentation
**Scope**: Step-by-step migration processes, validation checkpoints, troubleshooting

### ❌ Not ideal for documentation-specialist:

#### "Update README with latest blueprint count"
**Why**: Simple status update, better handled by documentation-sync-manager
**Better choice**: `documentation-sync-manager`

#### "Fix typo in quick-start guide"
**Why**: Minor correction, doesn't require comprehensive documentation expertise
**Better choice**: Direct edit

#### "Create GitHub issue for documentation bug"
**Why**: Project management task, not documentation creation
**Better choice**: `github-project-manager`

## Task Complexity Guidelines

### High Complexity (documentation-specialist)
- **Multi-document projects** spanning several guides or references
- **New content architecture** requiring information design
- **Technical content** requiring deep understanding of go-starter features
- **Multi-audience content** needing different complexity levels
- **Comprehensive guides** with examples, troubleshooting, and best practices

### Medium Complexity (coordination recommended)
- **Feature documentation** that affects multiple documents
- **Status updates** that require content restructuring
- **Migration content** between different systems or versions
- **Integration guides** requiring coordination with technical teams

### Low Complexity (other agents or direct action)
- **Single document updates** with clear scope
- **Status synchronization** across existing documents
- **Minor corrections** or small additions
- **Simple link updates** or formatting fixes

## Quality Standards

When using the documentation-specialist agent, expect:

### Content Quality
- **Technical Accuracy**: All examples tested and validated
- **Clear Communication**: Complex topics made accessible
- **Progressive Structure**: Information organized by complexity
- **Comprehensive Coverage**: Complete feature documentation

### User Experience
- **Multi-Audience Design**: Content appropriate for different expertise levels
- **Logical Navigation**: Clear paths through information
- **Actionable Content**: Users know what to do next
- **Consistent Style**: Standardized formatting and voice

### Integration Excellence
- **Agent Coordination**: Smooth workflow with other agents
- **Status Accuracy**: Current information across all documents
- **Cross-References**: Proper linking between related content
- **Maintenance Planning**: Content designed for easy updates

## Getting Started with documentation-specialist

### 1. Define Scope Clearly
- What type of content needs to be created or updated?
- Who is the target audience (beginner, intermediate, advanced)?
- What's the desired outcome for users?
- How does this fit into the broader documentation ecosystem?

### 2. Provide Context
- Current documentation status and gaps
- Related features or systems that should be covered
- Any existing content that should be referenced or updated
- Integration requirements with other documentation

### 3. Specify Quality Requirements
- Technical depth needed
- User experience goals
- Consistency requirements with existing content
- Maintenance and update considerations

### 4. Plan Coordination
- Which other agents need to be involved?
- What information or validation is needed from other systems?
- How should the work be coordinated with ongoing projects?

## Best Practices

### Effective Task Descriptions
- **Be Specific**: "Create comprehensive gRPC Gateway documentation" vs "update docs"
- **Include Context**: Mention related features, user needs, and integration points
- **Define Success**: What should users be able to do after reading the documentation?
- **Consider Maintenance**: How should the content be maintained and updated?

### Coordination Planning
- **Identify Dependencies**: What information is needed from other agents?
- **Plan Integration**: How will the content fit into existing documentation?
- **Consider Updates**: How will the content stay current with future changes?
- **Think Holistically**: How does this contribute to the overall documentation experience?

The documentation-specialist agent is designed to handle the most complex and comprehensive documentation challenges in the go-starter project, ensuring that all users—from beginners to experts—have access to excellent documentation that helps them succeed with their projects.