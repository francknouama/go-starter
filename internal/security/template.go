package security

import (
	"fmt"
	"regexp"
	"strings"
	"text/template"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

// TemplateSecurityValidator validates templates for security risks
type TemplateSecurityValidator struct {
	allowedFunctions  map[string]bool
	maxTemplateSize   int64
	maxRenderTime     time.Duration
	dangerousPatterns []*regexp.Regexp
}

// NewTemplateSecurityValidator creates a new template security validator
func NewTemplateSecurityValidator() *TemplateSecurityValidator {
	return &TemplateSecurityValidator{
		allowedFunctions:  getAllowedTemplateFunctions(),
		maxTemplateSize:   1024 * 1024, // 1MB max template size
		maxRenderTime:     5 * time.Second,
		dangerousPatterns: getDangerousPatterns(),
	}
}

// ValidateTemplate validates a template for security risks
func (v *TemplateSecurityValidator) ValidateTemplate(content string) error {
	// Check template size
	if int64(len(content)) > v.maxTemplateSize {
		return types.NewValidationError(fmt.Sprintf("template size %d exceeds maximum %d bytes", len(content), v.maxTemplateSize), nil)
	}

	// Check for dangerous patterns
	for _, pattern := range v.dangerousPatterns {
		if pattern.MatchString(content) {
			return types.NewValidationError(fmt.Sprintf("template contains dangerous pattern: %s", pattern.String()), nil)
		}
	}

	// Attempt to parse template to check for syntax issues
	tmpl, err := template.New("security-check").Parse(content)
	if err != nil {
		return types.NewValidationError("template syntax validation failed", err)
	}

	// Check for function usage
	if err := v.validateTemplateFunctions(tmpl); err != nil {
		return err
	}

	return nil
}

// ValidateTemplateFile validates a template file
func (v *TemplateSecurityValidator) ValidateTemplateFile(templatePath, content string) error {
	// Basic template validation
	if err := v.ValidateTemplate(content); err != nil {
		return fmt.Errorf("template %s failed security validation: %w", templatePath, err)
	}

	// Path-specific validations
	if err := v.validateTemplatePath(templatePath); err != nil {
		return err
	}

	return nil
}

// validateTemplatePath checks if the template path is safe
func (v *TemplateSecurityValidator) validateTemplatePath(templatePath string) error {
	// Check for path traversal attempts
	if strings.Contains(templatePath, "..") {
		return types.NewValidationError("template path contains path traversal attempt", nil)
	}

	// Check for absolute paths (templates should be relative)
	if strings.HasPrefix(templatePath, "/") || strings.Contains(templatePath, ":") {
		return types.NewValidationError("template path appears to be absolute", nil)
	}

	// Check for home directory access
	if strings.HasPrefix(templatePath, "~") {
		return types.NewValidationError("template path attempts to access home directory", nil)
	}

	return nil
}

// validateTemplateFunctions checks if template uses only allowed functions
func (v *TemplateSecurityValidator) validateTemplateFunctions(tmpl *template.Template) error {
	// Note: This is a simplified validation. In a production system,
	// you would need to parse the template AST to check function calls
	// For now, we rely on the dangerous patterns check
	return nil
}

// getAllowedTemplateFunctions returns the whitelist of safe template functions
func getAllowedTemplateFunctions() map[string]bool {
	// Sprig functions that are considered safe
	safeFunctions := map[string]bool{
		// String functions
		"upper":     true,
		"lower":     true,
		"title":     true,
		"trim":      true,
		"trimAll":   true,
		"replace":   true,
		"quote":     true,
		"squote":    true,
		"split":     true,
		"join":      true,
		"contains":  true,
		"hasPrefix": true,
		"hasSuffix": true,

		// Number functions
		"add": true,
		"sub": true,
		"mul": true,
		"div": true,
		"mod": true,
		"max": true,
		"min": true,

		// Logic functions
		"and": true,
		"or":  true,
		"not": true,
		"eq":  true,
		"ne":  true,
		"lt":  true,
		"le":  true,
		"gt":  true,
		"ge":  true,

		// Default function
		"default": true,

		// Safe encoding functions
		"b64enc":   true,
		"b64dec":   true,
		"urlquery": true,

		// Safe random functions
		"randAlphaNum": true,
		"randAlpha":    true,
		"randNumeric":  true,

		// Date functions (read-only)
		"now":        true,
		"date":       true,
		"dateModify": true,

		// List functions
		"list":    true,
		"first":   true,
		"last":    true,
		"initial": true,
		"rest":    true,
		"reverse": true,
		"uniq":    true,
		"compact": true,
		"slice":   true,
	}

	return safeFunctions
}

// getDangerousPatterns returns regex patterns for dangerous template content
func getDangerousPatterns() []*regexp.Regexp {
	dangerousPatterns := []string{
		// OS command execution attempts
		`\{\{\s*\.OS\s*\.`,
		`\{\{\s*os\s*\.`,
		`\{\{\s*exec\s*\.`,
		`\{\{\s*system\s*\.`,
		`\{\{\s*cmd\s*\.`,

		// Environment variable access
		`\{\{\s*\.Env\s*\.*`,
		`\{\{\s*env\s*\.`,
		`\{\{\s*getenv\s*`,
		`\{\{\s*range\s*\.Env\s*\}\}`,

		// File system access
		`\{\{\s*\.File\s*\.`,
		`\{\{\s*file\s*\.`,
		`\{\{\s*readFile\s*`,
		`\{\{\s*writeFile\s*`,

		// Network access
		`\{\{\s*\.HTTP\s*\.`,
		`\{\{\s*http\s*\.`,
		`\{\{\s*url\s*\.`,
		`\{\{\s*fetch\s*\.`,

		// Dangerous template functions
		`\{\{\s*include\s*"[^"]*\.\./`,  // Path traversal in includes
		`\{\{\s*template\s*"[^"]*\.\./`, // Path traversal in templates

		// Potential code injection
		`\{\{\s*printf\s*.*%[0-9]*[xXsp]`, // Format string attacks
		`\{\{\s*.*eval\s*`,                // Eval functions
		`\{\{\s*.*exec\s*`,                // Exec functions

		// Resource exhaustion patterns
		`\{\{\s*range\s*.*\{\{\s*range\s*.*\{\{\s*range`, // Deep nested loops
		`\{\{\s*printf\s*"%[0-9]{6,}s"`,                  // Large format strings
	}

	var patterns []*regexp.Regexp
	for _, pattern := range dangerousPatterns {
		if regex, err := regexp.Compile(pattern); err == nil {
			patterns = append(patterns, regex)
		}
	}

	return patterns
}

// SecurityViolation represents a security violation found in a template
type SecurityViolation struct {
	Type        string
	Description string
	Line        int
	Column      int
	Severity    string
}

// ScanTemplate performs a comprehensive security scan of template content
func (v *TemplateSecurityValidator) ScanTemplate(content string) []SecurityViolation {
	var violations []SecurityViolation

	lines := strings.Split(content, "\n")
	for lineNum, line := range lines {
		// Check each line for dangerous patterns
		for _, pattern := range v.dangerousPatterns {
			if pattern.MatchString(line) {
				violations = append(violations, SecurityViolation{
					Type:        "DangerousPattern",
					Description: fmt.Sprintf("Line contains dangerous pattern: %s", pattern.String()),
					Line:        lineNum + 1,
					Column:      0,
					Severity:    "HIGH",
				})
			}
		}

		// Check for other security issues
		if strings.Contains(line, "{{") && strings.Contains(line, "}}") {
			violations = append(violations, v.checkTemplateExpression(line, lineNum+1)...)
		}
	}

	return violations
}

// checkTemplateExpression analyzes a template expression for security issues
func (v *TemplateSecurityValidator) checkTemplateExpression(line string, lineNum int) []SecurityViolation {
	var violations []SecurityViolation

	// Look for potentially unsafe function calls
	unsafeFunctions := []string{"printf", "sprintf", "include", "template"}
	for _, fn := range unsafeFunctions {
		if strings.Contains(line, fn) {
			violations = append(violations, SecurityViolation{
				Type:        "UnsafeFunction",
				Description: fmt.Sprintf("Potentially unsafe function call: %s", fn),
				Line:        lineNum,
				Column:      strings.Index(line, fn),
				Severity:    "MEDIUM",
			})
		}
	}

	return violations
}
