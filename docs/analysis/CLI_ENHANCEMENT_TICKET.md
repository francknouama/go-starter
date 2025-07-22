# CLI Enhancement Ticket

**Ticket ID:** CLI-001  
**Type:** Enhancement  
**Priority:** Medium  
**Status:** Pending  
**Created:** December 26, 2024  
**Assignee:** TBD  

---

## 📋 Summary

Add Go version selector to CLI interactive prompts to maintain feature parity with planned web tool and improve developer experience.

## 🎯 User Story

**As a developer using go-starter CLI**, I want to select a specific Go version for my project so that I can target specific Go environments and maintain consistency across my development workflow.

## 📝 Description

Currently, the CLI uses auto-detection for Go version or relies on system defaults. To improve developer experience and maintain parity with the planned web tool, we need to add an interactive Go version selector that offers:

- Latest 3 major Go versions (1.23, 1.22, 1.21)
- Auto-detect option (current behavior, should remain default)
- Clear indication of which version is recommended

## ✅ Acceptance Criteria

### **Interactive Prompts Enhancement**
- [ ] Add Go version selection prompt to interactive mode
- [ ] Display options: "Auto-detect (recommended)", "Go 1.23", "Go 1.22", "Go 1.21"
- [ ] Default to "Auto-detect" to maintain backward compatibility
- [ ] Show current system Go version when auto-detect is selected

### **CLI Flags Support**
- [ ] Add `--go-version` flag for non-interactive usage
- [ ] Support values: `auto`, `1.23`, `1.22`, `1.21`
- [ ] Validate Go version input and show helpful error messages

### **Template Integration**
- [ ] Update go.mod template generation to use selected version
- [ ] Ensure all 4 templates (web-api, cli, library, lambda) respect version selection
- [ ] Maintain compatibility with existing template logic

### **Documentation**
- [ ] Update help text for interactive prompts
- [ ] Add `--go-version` flag to CLI help documentation
- [ ] Update docs/QUICK_REFERENCE.md with new option

## 🛠️ Technical Implementation

### **Files to Modify:**
```
internal/prompts/interactive.go    # Add Go version prompt
cmd/new.go                        # Add --go-version flag
internal/generator/generator.go    # Use selected Go version
templates/*/go.mod.tmpl           # Update Go version in templates
```

### **Code Structure:**
```go
// internal/prompts/interactive.go
func promptGoVersion() (string, error) {
    prompt := &survey.Select{
        Message: "Select Go version:",
        Options: []string{
            "Auto-detect (recommended)",
            "Go 1.23 (latest)",
            "Go 1.22", 
            "Go 1.21",
        },
        Default: "Auto-detect (recommended)",
    }
    
    var result string
    err := survey.AskOne(prompt, &result)
    return mapSelectionToVersion(result), err
}

// cmd/new.go
var goVersionFlag string

func init() {
    newCmd.Flags().StringVar(&goVersionFlag, "go-version", "auto", 
        "Go version to use (auto, 1.23, 1.22, 1.21)")
}
```

### **Version Support Strategy:**
```go
// internal/config/versions.go
const (
    GoVersionLatest = "1.23"
    GoVersionMinus1 = "1.22"
    GoVersionMinus2 = "1.21"
)

func GetSupportedGoVersions() []string {
    return []string{"auto", GoVersionLatest, GoVersionMinus1, GoVersionMinus2}
}

func ValidateGoVersion(version string) error {
    supported := GetSupportedGoVersions()
    for _, v := range supported {
        if v == version {
            return nil
        }
    }
    return fmt.Errorf("unsupported Go version: %s. Supported versions: %v", 
        version, supported)
}
```

## 🧪 Testing Requirements

### **Unit Tests:**
- [ ] Test Go version prompt functionality
- [ ] Test CLI flag parsing and validation
- [ ] Test version selection in different scenarios

### **Integration Tests:**
- [ ] Generate projects with different Go versions
- [ ] Verify go.mod contains correct Go version
- [ ] Test both interactive and flag-based usage

### **Manual Testing:**
- [ ] Test interactive mode with all version options
- [ ] Test `--go-version` flag with valid/invalid values
- [ ] Verify generated projects compile with selected Go version

## 📊 Success Metrics

- [ ] **Backward Compatibility:** Existing workflows continue to work unchanged
- [ ] **User Experience:** Clear and intuitive version selection process
- [ ] **Consistency:** Same options available in both CLI and future web tool
- [ ] **Flexibility:** Supports both interactive and non-interactive usage

## 🔗 Related Work

### **Dependencies:**
- Must be completed before web tool development
- Should maintain consistency with WEB_TOOL_BACKLOG.md specifications

### **Future Work:**
- Quarterly updates when new Go versions are released
- Consider adding version support validation (warn if selected version not installed)

## 🏷️ Tags

`enhancement` `cli` `go-version` `user-experience` `web-tool-parity` `interactive-prompts`

## 📅 Timeline

**Estimated Effort:** 2-3 days  
**Target Completion:** Before web tool development begins  
**Blocker for:** Web tool feature parity

## 💡 Additional Notes

### **Design Considerations:**
- Maintain current auto-detect as default for backward compatibility
- Keep version list current (update quarterly)
- Provide clear feedback when selected version differs from system version

### **Edge Cases:**
- What if selected Go version is not installed on system?
- How to handle development vs. production Go version differences?
- Should we validate that selected version can compile the generated project?

---

**Status Updates:**
- [ ] 2024-12-26: Ticket created, pending assignment
- [ ] TBD: Development started
- [ ] TBD: Testing completed
- [ ] TBD: Merged and released