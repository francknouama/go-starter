# CLI Blueprint Recommendations

## Executive Summary

Based on the comprehensive analysis of the cli-standard and cli-simple blueprints, this document provides actionable recommendations for the go-starter project to better serve users with appropriate complexity levels.

## Key Findings

1. **cli-standard is over-engineered** for 80% of CLI use cases
2. **cli-simple successfully reduces complexity** by 73% (8 files vs 30 files)
3. **Users need clear guidance** on which blueprint to choose
4. **A middle tier might be beneficial** for growing projects

## Immediate Recommendations

### 1. Update Default Selection Logic

Modify the CLI blueprint selection to default to cli-simple:

```go
// In internal/prompts/interactive.go
func selectCLIComplexity() string {
    if userIsNew || !advancedMode {
        return "simple"  // Default to simple
    }
    // Show complexity options
}
```

### 2. Improve Blueprint Descriptions

Update template.yaml descriptions to be more prescriptive:

**cli-simple**:
```yaml
description: "Simple CLI for scripts, utilities, and prototypes (< 100 lines, 1-3 commands)"
```

**cli-standard**:
```yaml
description: "Production CLI with config files, testing, and CI/CD (> 5 commands, team development)"
```

### 3. Add Selection Wizard

Implement a blueprint selection helper:

```go
func recommendCLIBlueprint() string {
    questions := []string{
        "How many commands will your CLI have? (1-3, 4-10, 10+)",
        "Do you need configuration files? (y/n)",
        "Will this be distributed publicly? (y/n)",
        "Do you need CI/CD pipelines? (y/n)",
    }
    // Logic to recommend based on answers
}
```

## Medium-term Recommendations

### 1. Create cli-medium Blueprint

Fill the gap between simple and standard:

```yaml
name: "cli-medium"
description: "Growing CLI with config and testing (15-18 files)"
features:
  - Multiple commands (5-8)
  - Simple config file support
  - Basic testing structure
  - Standard library logging (slog)
  - No Docker/CI by default
```

### 2. Progressive Enhancement Documentation

Create a migration guide showing how to evolve from simple to standard:

```markdown
# CLI Evolution Guide

## Stage 1: cli-simple (Start here)
- Single purpose tool
- 1-3 commands
- No configuration needed

## Stage 2: Add Configuration
- When: Need persistent settings
- How: Add internal/config package
- Example: User preferences, API endpoints

## Stage 3: Add Testing
- When: Critical business logic
- How: Add *_test.go files
- Example: Data processing, calculations

## Stage 4: Add CI/CD
- When: Team development
- How: Add .github/workflows
- Example: Automated releases
```

### 3. Interactive Complexity Assessment

Add to the generator:

```go
type CLIComplexityAssessment struct {
    CommandCount      int
    NeedsConfig      bool
    TeamSize         int
    PublicDistribution bool
    ExpectedUsers    int
}

func (a CLIComplexityAssessment) RecommendBlueprint() string {
    score := 0
    if a.CommandCount > 3 { score += 2 }
    if a.NeedsConfig { score += 2 }
    if a.TeamSize > 1 { score += 1 }
    if a.PublicDistribution { score += 3 }
    if a.ExpectedUsers > 100 { score += 2 }
    
    switch {
    case score <= 2:
        return "cli-simple"
    case score <= 6:
        return "cli-medium" // future
    default:
        return "cli-standard"
    }
}
```

## Long-term Recommendations

### 1. Blueprint Analytics

Track which blueprints users select and their success:

```go
type BlueprintMetrics struct {
    Selected   map[string]int
    Completed  map[string]int
    Abandoned  map[string]int
    Migrated   map[string]map[string]int // from -> to
}
```

### 2. Adaptive Blueprints

Allow blueprints to grow with user needs:

```bash
go-starter enhance --add-config
go-starter enhance --add-testing
go-starter enhance --add-docker
```

### 3. Community Feedback Loop

- Survey users after 30 days
- Track migration patterns
- Identify missing features
- Adjust blueprints based on data

## Implementation Priority

### Phase 1 (Immediate)
1. ✅ cli-simple blueprint exists
2. Update selection logic to favor simple
3. Improve blueprint descriptions
4. Add selection guidance to README

### Phase 2 (Next Sprint)
1. Create cli-medium blueprint
2. Write migration documentation
3. Add complexity assessment
4. Update web UI to reflect tiers

### Phase 3 (Future)
1. Implement blueprint analytics
2. Add enhancement commands
3. Create feedback system
4. Build adaptive features

## Success Metrics

Track these metrics to validate the approach:

1. **Adoption Rate**: % choosing cli-simple vs cli-standard
2. **Success Rate**: % completing project generation
3. **Migration Rate**: % upgrading from simple to standard
4. **Time to First Commit**: How quickly users start coding
5. **User Satisfaction**: Survey scores and feedback

## Risk Mitigation

### Risk: Users Choose Wrong Blueprint
**Mitigation**: Clear descriptions, selection wizard, easy migration path

### Risk: cli-simple Too Limited
**Mitigation**: Document growth path, provide cli-medium option

### Risk: Fragmentation
**Mitigation**: Limit to 3 CLI blueprints max, clear use cases

## Conclusion

The two-tier CLI blueprint approach successfully addresses the over-engineering problem. With these recommendations, go-starter can better serve both beginners and advanced users while maintaining a clear growth path. The key is making the simple path the default while keeping advanced features accessible when needed.

The philosophy should be: **"Start simple, grow as needed, never more complex than necessary."**