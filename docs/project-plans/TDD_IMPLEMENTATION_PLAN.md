# TDD Implementation Plan for go-starter

**Philosophy:** Red → Green → Refactor  
**Goal:** Embrace TDD to naturally achieve 85%+ coverage and improve code quality  
**Timeline:** Immediate adoption for all new development + systematic retrofit of existing code

---

## 🔴 TDD Fundamentals for go-starter

### **The TDD Cycle**
1. **🔴 Red:** Write a failing test first
2. **🟢 Green:** Write minimal code to make test pass
3. **🔵 Refactor:** Improve code while keeping tests green
4. **Repeat:** Continue cycle for each new requirement

### **TDD Benefits for Our Project**
- **Natural 85%+ coverage** - Tests drive all code creation
- **Better design** - Forces thinking about interfaces first
- **Regression protection** - Prevents breaking existing functionality
- **Documentation** - Tests serve as living specification
- **Confidence** - Safe to refactor and add features

---

## 📋 TDD Implementation Strategy

### **Phase 1: TDD for New Features (Immediate)**
- **All new development** must follow TDD
- **Go version selector** (GitHub issue #26) - perfect TDD candidate
- **Future enterprise templates** - built with TDD from start
- **Web tool backend** - API endpoints developed with TDD

### **Phase 2: Retrofit Existing Code (4 weeks)**
- **Systematic TDD retrofit** of existing untested code
- **Focus on critical gaps** identified in coverage analysis
- **Write tests first** for missing functionality
- **Refactor** existing code to be more testable

---

## 🛠️ TDD Infrastructure Setup

### **Testing Framework Enhancement**
```go
// tests/tdd/
├── helpers/
│   ├── assertions.go       # Custom assertion helpers
│   ├── mocks/             # Mock implementations
│   │   ├── file_system.go
│   │   ├── git_client.go
│   │   ├── prompter.go
│   │   └── template_engine.go
│   ├── builders/          # Test data builders
│   │   ├── project_builder.go
│   │   ├── config_builder.go
│   │   └── template_builder.go
│   └── fixtures/          # Test fixtures and data
└── examples/              # TDD example implementations
```

### **Custom Test Assertions**
```go
// tests/tdd/helpers/assertions.go
package helpers

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

// AssertProjectGenerated validates complete project generation
func AssertProjectGenerated(t *testing.T, outputDir string, expectedFiles []string) {
    t.Helper()
    for _, file := range expectedFiles {
        assert.FileExists(t, filepath.Join(outputDir, file), 
            "Expected file %s should exist", file)
    }
}

// AssertGoModValid validates go.mod file structure
func AssertGoModValid(t *testing.T, goModPath string, expectedModule string) {
    t.Helper()
    content, err := os.ReadFile(goModPath)
    assert.NoError(t, err)
    assert.Contains(t, string(content), expectedModule)
}

// AssertCompilable validates generated project compiles
func AssertCompilable(t *testing.T, projectDir string) {
    t.Helper()
    cmd := exec.Command("go", "build", "./...")
    cmd.Dir = projectDir
    err := cmd.Run()
    assert.NoError(t, err, "Generated project should compile successfully")
}
```

### **Mock Interfaces**
```go
// tests/tdd/helpers/mocks/file_system.go
type MockFileSystem struct {
    files map[string]string
    dirs  map[string]bool
}

func (m *MockFileSystem) WriteFile(path string, content []byte) error {
    if m.shouldFailWrite(path) {
        return errors.New("mock write failure")
    }
    m.files[path] = string(content)
    return nil
}

func (m *MockFileSystem) ReadFile(path string) ([]byte, error) {
    content, exists := m.files[path]
    if !exists {
        return nil, os.ErrNotExist
    }
    return []byte(content), nil
}

// Configure mock behaviors for testing
func (m *MockFileSystem) ShouldFailWrite(path string) {
    // Configure mock to simulate failures
}
```

---

## 🎯 TDD Example: Go Version Selector (Issue #26)

### **Step 1: Write Failing Test First**
```go
// internal/prompts/go_version_test.go
func TestGoVersionPrompt(t *testing.T) {
    tests := []struct {
        name           string
        userSelection  string
        expectedResult string
    }{
        {
            name:           "auto-detect selection",
            userSelection:  "Auto-detect (recommended)",
            expectedResult: "auto",
        },
        {
            name:           "go 1.23 selection",
            userSelection:  "Go 1.23 (latest)",
            expectedResult: "1.23",
        },
        {
            name:           "go 1.22 selection", 
            userSelection:  "Go 1.22",
            expectedResult: "1.22",
        },
        {
            name:           "go 1.21 selection",
            userSelection:  "Go 1.21", 
            expectedResult: "1.21",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockPrompter := &MockPrompter{
                responses: map[string]string{
                    "goVersion": tt.userSelection,
                },
            }
            
            // Act
            result, err := mockPrompter.PromptGoVersion()
            
            // Assert
            assert.NoError(t, err)
            assert.Equal(t, tt.expectedResult, result)
        })
    }
}

func TestGoVersionValidation(t *testing.T) {
    tests := []struct {
        name        string
        version     string
        shouldError bool
    }{
        {"valid auto", "auto", false},
        {"valid 1.23", "1.23", false}, 
        {"valid 1.22", "1.22", false},
        {"valid 1.21", "1.21", false},
        {"invalid 1.20", "1.20", true},
        {"invalid 2.0", "2.0", true},
        {"empty version", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Act
            err := ValidateGoVersion(tt.version)
            
            // Assert
            if tt.shouldError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

**🔴 Run Test:** `go test ./internal/prompts -run TestGoVersion`  
**Result:** Test fails (functions don't exist yet)

### **Step 2: Write Minimal Code to Pass**
```go
// internal/prompts/go_version.go
package prompts

import (
    "fmt"
    "strings"
)

var supportedGoVersions = []string{"auto", "1.23", "1.22", "1.21"}

func (p *Prompter) PromptGoVersion() (string, error) {
    options := []string{
        "Auto-detect (recommended)",
        "Go 1.23 (latest)",
        "Go 1.22",
        "Go 1.21",
    }
    
    prompt := &survey.Select{
        Message: "Select Go version:",
        Options: options,
        Default: "Auto-detect (recommended)",
    }
    
    var selection string
    if err := survey.AskOne(prompt, &selection); err != nil {
        return "", err
    }
    
    return mapSelectionToVersion(selection), nil
}

func mapSelectionToVersion(selection string) string {
    switch selection {
    case "Auto-detect (recommended)":
        return "auto"
    case "Go 1.23 (latest)":
        return "1.23"
    case "Go 1.22":
        return "1.22" 
    case "Go 1.21":
        return "1.21"
    default:
        return "auto"
    }
}

func ValidateGoVersion(version string) error {
    for _, supported := range supportedGoVersions {
        if supported == version {
            return nil
        }
    }
    return fmt.Errorf("unsupported Go version: %s", version)
}
```

**🟢 Run Test:** `go test ./internal/prompts -run TestGoVersion`  
**Result:** Tests pass

### **Step 3: Refactor & Add Integration Tests**
```go
// internal/prompts/go_version_integration_test.go
func TestGoVersionIntegration(t *testing.T) {
    t.Run("go version used in template generation", func(t *testing.T) {
        // Arrange
        config := &types.ProjectConfig{
            Name:       "test-project",
            ModulePath: "github.com/test/project",
            Type:       "web-api",
            GoVersion:  "1.23",
        }
        
        generator := generator.New()
        tempDir := t.TempDir()
        
        // Act
        err := generator.Generate(config, tempDir)
        
        // Assert
        assert.NoError(t, err)
        
        // Verify go.mod contains correct version
        goModPath := filepath.Join(tempDir, "go.mod")
        helpers.AssertGoModValid(t, goModPath, "1.23")
    })
}
```

---

## 📈 TDD for Existing Code Retrofit

### **Example: File Operations (Currently 0% covered)**

#### **Step 1: Write Tests for Desired Behavior**
```go
// internal/utils/files_test.go
func TestWriteFile(t *testing.T) {
    tests := []struct {
        name        string
        path        string
        content     []byte
        permissions os.FileMode
        shouldError bool
    }{
        {
            name:        "write file successfully",
            path:        "test.txt",
            content:     []byte("hello world"),
            permissions: 0644,
            shouldError: false,
        },
        {
            name:        "write file with custom permissions",
            path:        "executable.sh", 
            content:     []byte("#!/bin/bash\necho hello"),
            permissions: 0755,
            shouldError: false,
        },
        {
            name:        "write to invalid path should error",
            path:        "/invalid/path/file.txt",
            content:     []byte("content"),
            permissions: 0644,
            shouldError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            tempDir := t.TempDir()
            fullPath := filepath.Join(tempDir, tt.path)
            
            // Act
            err := WriteFileWithPermissions(fullPath, tt.content, tt.permissions)
            
            // Assert
            if tt.shouldError {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.FileExists(t, fullPath)
            
            // Verify content
            actualContent, err := os.ReadFile(fullPath)
            assert.NoError(t, err)
            assert.Equal(t, tt.content, actualContent)
            
            // Verify permissions
            info, err := os.Stat(fullPath)
            assert.NoError(t, err)
            assert.Equal(t, tt.permissions, info.Mode().Perm())
        })
    }
}

func TestCopyFile(t *testing.T) {
    t.Run("copy file successfully", func(t *testing.T) {
        // Arrange
        tempDir := t.TempDir()
        srcPath := filepath.Join(tempDir, "source.txt")
        destPath := filepath.Join(tempDir, "dest.txt")
        content := []byte("test content")
        
        err := os.WriteFile(srcPath, content, 0644)
        require.NoError(t, err)
        
        // Act
        err = CopyFile(srcPath, destPath)
        
        // Assert
        assert.NoError(t, err)
        assert.FileExists(t, destPath)
        
        actualContent, err := os.ReadFile(destPath)
        assert.NoError(t, err)
        assert.Equal(t, content, actualContent)
    })
    
    t.Run("copy non-existent file should error", func(t *testing.T) {
        // Act
        err := CopyFile("/non/existent/file.txt", "/dest/path.txt")
        
        // Assert
        assert.Error(t, err)
    })
}
```

#### **Step 2: Implement/Fix Code to Pass Tests**
```go
// internal/utils/files.go - Enhanced implementation
func WriteFileWithPermissions(path string, content []byte, perm os.FileMode) error {
    // Ensure directory exists
    dir := filepath.Dir(path)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return fmt.Errorf("failed to create directory %s: %w", dir, err)
    }
    
    // Write file with specified permissions
    if err := os.WriteFile(path, content, perm); err != nil {
        return fmt.Errorf("failed to write file %s: %w", path, err)
    }
    
    return nil
}

func CopyFile(src, dest string) error {
    sourceFile, err := os.Open(src)
    if err != nil {
        return fmt.Errorf("failed to open source file %s: %w", src, err)
    }
    defer sourceFile.Close()
    
    // Get source file info for permissions
    sourceInfo, err := sourceFile.Stat()
    if err != nil {
        return fmt.Errorf("failed to get source file info: %w", err)
    }
    
    // Create destination file
    destFile, err := os.Create(dest)
    if err != nil {
        return fmt.Errorf("failed to create destination file %s: %w", dest, err)
    }
    defer destFile.Close()
    
    // Copy content
    if _, err := io.Copy(destFile, sourceFile); err != nil {
        return fmt.Errorf("failed to copy file content: %w", err)
    }
    
    // Set same permissions as source
    if err := os.Chmod(dest, sourceInfo.Mode()); err != nil {
        return fmt.Errorf("failed to set file permissions: %w", err)
    }
    
    return nil
}
```

---

## 🏗️ TDD Workflow Integration

### **Development Workflow**
```bash
# 1. Start new feature with failing test
git checkout -b feature/go-version-selector
cd internal/prompts
vim go_version_test.go  # Write failing test

# 2. Run test to confirm failure  
go test -run TestGoVersion
# Expected: FAIL (function doesn't exist)

# 3. Write minimal implementation
vim go_version.go  # Implement function

# 4. Run test until green
go test -run TestGoVersion  
# Expected: PASS

# 5. Run all tests to ensure no regression
go test ./...

# 6. Check coverage improvement
go test -coverprofile=coverage.out ./internal/prompts
go tool cover -func=coverage.out

# 7. Refactor if needed while keeping tests green
# 8. Commit when cycle complete
git add .
git commit -m "feat: add Go version selector with TDD

- Write failing tests for Go version prompt
- Implement minimal PromptGoVersion function  
- Add validation for supported versions
- Coverage: internal/prompts improved from 3.2% to 45%"
```

### **Pre-commit Hooks**
```bash
#!/bin/sh
# .git/hooks/pre-commit

echo "Running TDD validation..."

# 1. All tests must pass
if ! go test ./...; then
    echo "❌ Tests failing - commit rejected"
    exit 1
fi

# 2. Coverage must not decrease
CURRENT_COVERAGE=$(go test -coverprofile=coverage.out ./... 2>/dev/null | \
    go tool cover -func=coverage.out | grep "total:" | \
    awk '{print $3}' | sed 's/%//')

if (( $(echo "${CURRENT_COVERAGE} < 85" | bc -l) )); then
    echo "❌ Coverage ${CURRENT_COVERAGE}% below 85% threshold"
    exit 1
fi

echo "✅ TDD validation passed - coverage: ${CURRENT_COVERAGE}%"
```

---

## 📊 TDD Success Metrics

### **Quality Indicators**
- **Test-to-Code Ratio:** Aim for 1:1 or better (lines of test : lines of code)
- **Cycle Time:** Red → Green → Refactor cycles under 10 minutes
- **Coverage Growth:** Natural increase without explicit targeting
- **Bug Reduction:** Fewer production issues due to better design

### **Weekly TDD Tracking**
```bash
# Track TDD adoption and quality
make tdd-metrics

# Output example:
TDD Metrics for Week 1:
- New functions written: 12
- Functions with tests first: 12 (100%)
- Average cycle time: 6.2 minutes  
- Coverage increase: +15.3%
- Test lines added: 847
- Code lines added: 623
- Test/Code ratio: 1.36:1 ✅
```

### **Team TDD Adoption**
- **Code Reviews:** Require tests in all PRs
- **TDD Pairing:** Pair programming with TDD focus
- **TDD Training:** Internal workshops and knowledge sharing
- **Tool Integration:** IDE plugins for TDD workflow

---

## 🎯 Immediate Next Steps

### **Week 1: TDD Infrastructure**
1. **Set up TDD helpers** (assertions, mocks, builders)
2. **Configure pre-commit hooks** for coverage enforcement
3. **Start Go version selector** with pure TDD approach
4. **Document TDD patterns** for team adoption

### **Week 2-4: Systematic Retrofit**
1. **Retrofit internal/utils** with TDD approach
2. **Retrofit internal/prompts** with comprehensive tests
3. **Retrofit internal/generator** core logic
4. **Validate 85% coverage** achievement

### **Ongoing: TDD Culture**
- **All new features** must use TDD
- **Zero tolerance** for untested code
- **Regular TDD retrospectives** and improvement
- **Celebrate TDD wins** and learning

---

**TDD will naturally drive us to 85%+ coverage while improving code quality, design, and confidence. It's the perfect approach for go-starter's next evolution!** 🚀