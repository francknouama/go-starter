# Contributing to Go-Starter

Thank you for your interest in contributing to go-starter! This document provides guidelines and workflows for contributing to the project.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [GitHub Project Management](#github-project-management)
- [Commit Guidelines](#commit-guidelines)
- [Pull Request Process](#pull-request-process)
- [Testing Requirements](#testing-requirements)
- [Documentation](#documentation)

## Code of Conduct

Please note that this project is released with a Contributor Code of Conduct. By participating in this project you agree to abide by its terms. Be kind, respectful, and professional in all interactions.

## Getting Started

1. **Fork the repository** and clone your fork:
   ```bash
   git clone https://github.com/YOUR_USERNAME/go-starter.git
   cd go-starter
   ```

2. **Set up your development environment**:
   ```bash
   # Install dependencies
   go mod download
   
   # Run tests to ensure everything works
   make test
   
   # Build the project
   make build
   ```

3. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

### 1. Check Existing Issues

Before starting work:
- Check the [GitHub Issues](https://github.com/francknouama/go-starter/issues) for existing issues
- Look at the [Project Board](https://github.com/users/francknouama/projects/4) for issue status
- If no issue exists for your work, create one first

### 2. Assign Yourself to an Issue

When you decide to work on an issue:
- Comment on the issue to let others know you're working on it
- If you have access, assign yourself to the issue

### 3. Update Project Status

**IMPORTANT**: Keep the GitHub Project board updated with your progress.

When starting work:
```bash
# Move issue to "In Progress" column
gh project item-edit --project-id PVT_kwHOAAjzD84A8ata --id <ITEM_ID> --field-id PVTSSF_lAHOAAjzD84A8atazgwbEdk --single-select-option-id 47fc9ee4
```

When completing work:
```bash
# Move issue to "Done" column
gh project item-edit --project-id PVT_kwHOAAjzD84A8ata --id <ITEM_ID> --field-id PVTSSF_lAHOAAjzD84A8atazgwbEdk --single-select-option-id 98236657
```

## GitHub Project Management

We use GitHub Projects to track progress. The project board has three columns:

| Column | ID | Purpose |
|--------|-----|---------|
| **Todo** | `f75ad846` | Issues ready to be worked on |
| **In Progress** | `47fc9ee4` | Issues currently being worked on |
| **Done** | `98236657` | Completed issues |

### Finding Item IDs

To find the item ID for an issue:
```bash
gh project item-list 4 --owner francknouama --format json | jq '.items[] | select(.content.url | contains("/issues/YOUR_ISSUE_NUMBER")) | {id, title}'
```

## Commit Guidelines

### Commit Message Format

**IMPORTANT**: All commits MUST reference the related issue number.

Format:
```
<type>(<scope>): <subject> (#<issue-number>)

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, missing semicolons, etc.)
- **refactor**: Code refactoring without changing functionality
- **perf**: Performance improvements
- **test**: Adding or updating tests
- **chore**: Maintenance tasks, dependency updates
- **ci**: CI/CD configuration changes

### Examples

```bash
# Feature commit
git commit -m "feat(cli): add progress indicators with Charm's Fang (#22)

- Implement spinner for project generation
- Add progress bar for file creation
- Include styled success/error messages

Resolves #22"

# Bug fix commit
git commit -m "fix(generator): correct template path resolution (#45)

- Fix relative path issues on Windows
- Add proper path normalization
- Update tests to cover edge cases

Fixes #45"

# Documentation commit
git commit -m "docs(readme): update installation instructions (#12)

- Add Homebrew installation method
- Include Windows installation guide
- Fix broken links

Updates #12"
```

### Commit Best Practices

1. **Reference issues**: Always include `#<issue-number>` in your commits
2. **Be descriptive**: Explain what and why, not just how
3. **Keep it focused**: One logical change per commit
4. **Use present tense**: "add feature" not "added feature"
5. **Line length**: Keep the subject line under 50 characters

## Pull Request Process

1. **Ensure your branch is up to date**:
   ```bash
   git checkout main
   git pull upstream main
   git checkout feature/your-feature
   git rebase main
   ```

2. **Run all tests and checks**:
   ```bash
   make test
   make lint
   ```

3. **Create Pull Request**:
   - Use a clear title that includes the issue number
   - Reference the issue in the PR description using "Fixes #XX" or "Resolves #XX"
   - Fill out the PR template completely
   - Add relevant labels

### PR Title Format
```
<type>(<scope>): <description> (#<issue-number>)
```

Example: `feat(cli): integrate Charm's Fang for UI enhancement (#22)`

### PR Description Template
```markdown
## Summary
Brief description of the changes

## Related Issue
Fixes #<issue-number>

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing completed

## Checklist
- [ ] My code follows the project's style guidelines
- [ ] I have performed a self-review
- [ ] I have added tests that prove my fix/feature works
- [ ] I have updated the documentation accordingly
```

## Testing Requirements

### Running Tests

```bash
# Run all tests
make test

# Run with race detection
make test-race

# Run specific tests
go test -v ./internal/generator/...

# Run integration tests
make test-integration
```

### Test Coverage

- Aim for at least 80% test coverage
- All new features must include tests
- Bug fixes should include regression tests

### Template Testing

When modifying templates:
1. Test all template+logger combinations
2. Ensure generated projects compile: `go build`
3. Run generated project tests: `go test ./...`
4. Validate against the test matrix in `scripts/test_all_combinations.sh`

## Documentation

### Code Documentation

- Add godoc comments to all exported functions, types, and packages
- Include examples in godoc comments where appropriate
- Keep comments up to date with code changes

### Project Documentation

When adding features:
1. Update README.md if it affects usage
2. Update CLAUDE.md if it affects AI assistance
3. Update PROJECT_ROADMAP.md for significant features
4. Add entries to CHANGELOG.md (if present)

### Template Documentation

For new templates:
1. Include comprehensive README.md in the template
2. Add usage examples
3. Document configuration options
4. Provide troubleshooting guide

## Development Guidelines

### Code Style

- Follow standard Go formatting (use `gofmt`)
- Use meaningful variable and function names
- Keep functions small and focused
- Handle errors explicitly
- Add comments for complex logic

### Dependencies

- Minimize external dependencies
- Document why each dependency is needed
- Keep dependencies up to date
- Use go modules for version management

### Performance

- Profile before optimizing
- Consider memory allocations in hot paths
- Use benchmarks to validate improvements
- Document performance-critical code

## Getting Help

- **Questions**: Open a [Discussion](https://github.com/francknouama/go-starter/discussions)
- **Bugs**: Open an [Issue](https://github.com/francknouama/go-starter/issues)
- **Security**: Email security concerns privately

## Recognition

Contributors will be:
- Listed in CONTRIBUTORS.md
- Mentioned in release notes
- Given credit in commit messages

Thank you for contributing to go-starter! 🚀