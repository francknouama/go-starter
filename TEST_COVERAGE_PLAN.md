# Test Coverage Improvement Plan

**Current Coverage:** 21.1%  
**Target Coverage:** 85%  
**Gap to Close:** +63.9%  
**Priority:** HIGH - Blocking for new feature development

---

## 🎯 Coverage Goals by Package

### **Phase 1: Critical Foundation (Weeks 1-2)**
| Package | Current | Target | Priority | Effort |
|---------|---------|--------|----------|--------|
| **internal/utils** | 5.7% | 85% | 🔴 Critical | High |
| **internal/prompts** | 3.2% | 80% | 🔴 Critical | High |
| **internal/generator** | 20.6% | 85% | 🔴 Critical | High |

### **Phase 2: Core Functionality (Week 3)**
| Package | Current | Target | Priority | Effort |
|---------|---------|--------|----------|--------|
| **cmd** | 36.2% | 85% | 🟡 High | Medium |
| **pkg/types** | 28.0% | 85% | 🟡 High | Medium |

### **Phase 3: Enhancement (Week 4)**
| Package | Current | Target | Priority | Effort |
|---------|---------|--------|----------|--------|
| **internal/config** | 41.8% | 85% | 🟢 Medium | Low |
| **internal/templates** | 51.1% | 85% | 🟢 Medium | Low |
| **internal/logger** | 63.6% | 85% | 🟢 Low | Low |

**Target Overall Coverage:** 85%+

---

## 📋 Detailed Implementation Plan

### **Week 1: Utils & File Operations (Priority 1)**

#### 🔧 internal/utils Package (5.7% → 85%)
**Current Gap:** 79.3% - **Highest Priority**

**Missing Coverage Areas:**
- **File Operations (0% covered)**
  - CopyFile, CopyDir, WriteFile, ReadFile
  - FileExists, DirExists, CreateDir
  - File permissions and validation
  
- **Git Operations (0% covered)**
  - InitGitRepository, GitAdd, GitCommit
  - Git configuration and status checks
  - Repository validation
  
- **Module Operations (0% covered)**
  - InitGoModule, GoModTidy, GoBuild
  - Go version detection and validation
  - Module path validation

**Test Files to Create/Enhance:**
```
internal/utils/
├── files_test.go        # File operations testing
├── git_test.go          # Git operations testing  
├── modules_test.go      # Go module operations testing
└── namegen_test.go      # ✅ Already good (93%+)
```

**Critical Tests Needed:**
```go
// files_test.go
func TestFileOperations(t *testing.T) {
    // Test file creation, reading, writing
    // Test directory operations
    // Test file permissions and validation
    // Test error conditions (permissions, disk space)
}

// git_test.go  
func TestGitOperations(t *testing.T) {
    // Test git repository initialization
    // Test git configuration
    // Test commit and status operations
    // Test error conditions (no git, invalid repo)
}

// modules_test.go
func TestGoModuleOperations(t *testing.T) {
    // Test go.mod creation and management
    // Test Go version detection
    // Test module building and testing
    // Test dependency management
}
```

---

### **Week 1: Interactive Prompts (Priority 2)**

#### 🎨 internal/prompts Package (3.2% → 80%)
**Current Gap:** 76.8% - **Critical for UX**

**Missing Coverage Areas:**
- **Interactive Mode Detection (0% covered)**
  - isInteractiveMode logic
  - Configuration completeness checks
  
- **All Prompt Functions (0% covered)**
  - promptProjectName, promptModulePath
  - promptProjectType, promptFramework
  - promptLogger, promptArchitecture
  - promptDatabaseSupport, promptORM
  - promptAuthentication

**Test Strategy:**
- **Mock user input** for automated testing
- **Test validation logic** without actual prompts
- **Test configuration building** from prompt responses
- **Test error handling** for invalid inputs

**Test Files to Create:**
```
internal/prompts/
├── interactive_test.go     # ✅ Exists but minimal
├── mock_prompter_test.go   # New: Mock implementation
└── validation_test.go      # New: Input validation tests
```

**Critical Tests Needed:**
```go
// mock_prompter_test.go
type MockPrompter struct {
    responses map[string]string
}

func TestPromptWorkflows(t *testing.T) {
    // Test complete prompt workflows
    // Test validation logic
    // Test configuration building
    // Test error handling
}

// validation_test.go
func TestInputValidation(t *testing.T) {
    // Test project name validation
    // Test module path validation  
    // Test framework/logger selection validation
    // Test database configuration validation
}
```

---

### **Week 2: Core Generation Logic (Priority 3)**

#### ⚙️ internal/generator Package (20.6% → 85%)
**Current Gap:** 64.4% - **Core Functionality**

**Missing Coverage Areas:**
- **Project Generation (0% covered)**
  - generateProjectFiles - main generation logic
  - processTemplatePath - template processing
  - processTemplateFile - file creation
  - createTemplateContext - context building

- **Template Processing (0% covered)**
  - evaluateCondition - conditional file generation
  - processDependencies - go.mod management
  - addDependencies - dependency injection

- **Hooks & Post-Processing (0% covered)**
  - executeHooks - post-generation scripts
  - Git operations integration
  - Go module operations

**Test Strategy:**
- **Integration tests** with real template generation
- **Unit tests** for individual functions
- **Mock file system** for isolated testing
- **Template validation** tests

**Test Files to Enhance:**
```
internal/generator/
├── generator_test.go      # ✅ Exists but minimal (20.6%)
├── template_test.go       # New: Template processing tests
├── context_test.go        # New: Context building tests
├── hooks_test.go          # New: Post-generation hook tests
└── integration_test.go    # New: End-to-end generation tests
```

**Critical Tests Needed:**
```go
// template_test.go
func TestTemplateProcessing(t *testing.T) {
    // Test template file processing
    // Test conditional file generation
    // Test context variable substitution
    // Test error handling for invalid templates
}

// context_test.go
func TestContextBuilding(t *testing.T) {
    // Test template context creation
    // Test variable mapping and validation
    // Test feature flag evaluation
    // Test database configuration context
}

// integration_test.go
func TestEndToEndGeneration(t *testing.T) {
    // Test complete project generation workflows
    // Test all 4 templates with all 4 loggers (16 combinations)
    // Test generated project compilation
    // Test hooks execution
}
```

---

### **Week 3: CLI & Types (Priority 4-5)**

#### 💻 cmd Package (36.2% → 85%)
**Current Gap:** 48.8%

**Missing Coverage Areas:**
- **Command Execution (0% covered)**
  - runNew - main command logic
  - Command flag parsing and validation
  - Error handling and user feedback

- **List Command (0% covered)**
  - listTemplates - template listing
  - Template rendering and display

- **Version Command (0% covered)**
  - showVersion - version display

**Test Files to Enhance:**
```
cmd/
├── new_test.go      # ✅ Exists but minimal
├── root_test.go     # ✅ Exists but minimal  
├── list_test.go     # New: List command tests
└── version_test.go  # New: Version command tests
```

#### 📦 pkg/types Package (28.0% → 85%)
**Current Gap:** 57%

**Missing Coverage Areas:**
- **Project Configuration (0% covered)**
  - HasDatabase, GetDrivers, HasDriver
  - PrimaryDriver, configuration validation
  
- **Error Types (0% covered)**
  - NewGenerationError, NewFileSystemError
  - NewConfigError, error wrapping

---

### **Week 4: Configuration & Templates (Priority 6-7)**

#### ⚙️ internal/config Package (41.8% → 85%)
**Current Gap:** 43.2%

**Missing Coverage Areas:**
- **Validation Functions (0% covered)**
  - All ValidateXxx functions in validation.go
  - Input sanitization and format checking
  - Cross-field validation logic

#### 📄 internal/templates Package (51.1% → 85%)
**Current Gap:** 33.9%

**Missing Coverage Areas:**
- **Template Loading (0% covered)**
  - LoadTemplateFile, GetTemplatePath
  - FileExists, template validation
  - Error handling for missing templates

---

## 🧪 Testing Infrastructure Improvements

### **Enhanced Test Utilities**
```go
// tests/testutils/
├── mock_fs.go           # Mock file system for testing
├── mock_git.go          # Mock git operations
├── mock_prompter.go     # Mock user input
├── project_builder.go   # Test project generation helper
└── assertion_helpers.go # Custom test assertions
```

### **Coverage Enforcement**
```go
// Add to Makefile
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out | grep "total:" | \
	awk '{if ($$3+0 < 85) {print "Coverage " $$3 " is below 85% threshold"; exit 1}}'

test-coverage-html:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: $(PWD)/coverage.html"
```

### **CI/CD Integration**
```yaml
# .github/workflows/ci.yml enhancement
- name: Test Coverage Check
  run: |
    go test -v -coverprofile=coverage.out ./...
    COVERAGE=$(go tool cover -func=coverage.out | grep "total:" | awk '{print $3}' | sed 's/%//')
    echo "Current coverage: ${COVERAGE}%"
    if (( $(echo "${COVERAGE} < 85" | bc -l) )); then
      echo "❌ Coverage ${COVERAGE}% is below 85% threshold"
      exit 1
    fi
    echo "✅ Coverage ${COVERAGE}% meets 85% threshold"
```

---

## 📊 Success Metrics & Milestones

### **Weekly Targets:**
- **Week 1:** 21.1% → 45% (Focus: utils + prompts)
- **Week 2:** 45% → 65% (Focus: generator)
- **Week 3:** 65% → 80% (Focus: cmd + types)
- **Week 4:** 80% → 85%+ (Focus: config + templates)

### **Quality Gates:**
- **No new code without tests** (enforce with CI)
- **All critical paths covered** (file operations, generation, prompts)
- **Integration tests pass** (all template+logger combinations)
- **Performance regression tests** (generation time < 10s)

### **Milestone Validation:**
```bash
# After each week, validate:
make test-coverage-html
open coverage.html

# Check specific package coverage:
go test -coverprofile=coverage.out ./internal/utils
go tool cover -func=coverage.out
```

---

## 🚫 Coverage Blockers & Risks

### **High-Risk Areas:**
1. **File System Operations** - Hard to test, require careful mocking
2. **Interactive Prompts** - Need mock user input strategies
3. **Git Operations** - Require repository setup/teardown
4. **Template Generation** - Complex integration testing

### **Mitigation Strategies:**
1. **Mock Interfaces** - Create mockable interfaces for external dependencies
2. **Test Utilities** - Build robust test helper functions
3. **Temporary Environments** - Use temp directories for file operations
4. **Parallel Testing** - Ensure tests can run concurrently

---

## 🎯 Next Steps

1. **Week 1 Priority:** Start with `internal/utils` package (biggest gap)
2. **Setup Test Infrastructure:** Mock interfaces and test utilities
3. **Establish CI Gates:** Enforce 85% coverage threshold
4. **Document Patterns:** Create testing best practices guide

The 85% coverage target is ambitious but achievable with focused effort on the critical gaps. This foundation will make future development (Go version selector, web tool, enterprise templates) much more reliable and maintainable.