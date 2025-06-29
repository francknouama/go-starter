#!/bin/bash

# Template Analysis Script
# Analyzes all .tmpl files for syntax errors, variable consistency, and security issues

TEMPLATE_DIR="/Users/franck/reactive-crafters-workspace/golang-project-generator/templates"
LOG_FILE="/tmp/template_analysis.log"
ERROR_LOG="/tmp/template_errors.log"

echo "=== TEMPLATE ANALYSIS REPORT ===" > "$LOG_FILE"
echo "Generated: $(date)" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Find all template files
echo "Finding all template files..." >> "$LOG_FILE"
TEMPLATE_FILES=$(find "$TEMPLATE_DIR" -name "*.tmpl" -type f)
TOTAL_FILES=$(echo "$TEMPLATE_FILES" | wc -l)
echo "Total template files found: $TOTAL_FILES" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Initialize error tracking
> "$ERROR_LOG"

# 1. Check for basic template syntax errors
echo "=== 1. TEMPLATE SYNTAX ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

syntax_errors=0
for file in $TEMPLATE_FILES; do
    # Check for unmatched braces
    unmatched_open=$(grep -o '{{' "$file" | wc -l)
    unmatched_close=$(grep -o '}}' "$file" | wc -l)
    
    if [ "$unmatched_open" -ne "$unmatched_close" ]; then
        echo "SYNTAX ERROR: Unmatched template braces in $file" >> "$LOG_FILE"
        echo "  Open braces: $unmatched_open, Close braces: $unmatched_close" >> "$LOG_FILE"
        echo "$file: Unmatched template braces" >> "$ERROR_LOG"
        ((syntax_errors++))
    fi
    
    # Check for invalid template constructs
    if grep -q '{{[^}]*{{' "$file"; then
        echo "SYNTAX ERROR: Nested template braces in $file" >> "$LOG_FILE"
        grep -n '{{[^}]*{{' "$file" | head -5 >> "$LOG_FILE"
        echo "$file: Nested template braces" >> "$ERROR_LOG"
        ((syntax_errors++))
    fi
    
    # Check for template injection vulnerabilities
    if grep -qE '\{\{.*\|.*exec.*\}\}|\{\{.*\|.*system.*\}\}|\{\{.*\|.*shell.*\}\}' "$file"; then
        echo "SECURITY WARNING: Potential template injection in $file" >> "$LOG_FILE"
        grep -nE '\{\{.*\|.*exec.*\}\}|\{\{.*\|.*system.*\}\}|\{\{.*\|.*shell.*\}\}' "$file" >> "$LOG_FILE"
        echo "$file: Potential template injection" >> "$ERROR_LOG"
        ((syntax_errors++))
    fi
done

echo "Total syntax errors found: $syntax_errors" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 2. Variable consistency analysis
echo "=== 2. VARIABLE CONSISTENCY ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Extract all template variables from all files
echo "Extracting template variables..." >> "$LOG_FILE"
ALL_VARIABLES=$(grep -ho '\{\{[^}]*\}\}' $TEMPLATE_FILES | sort | uniq)
echo "Unique template constructs found:" >> "$LOG_FILE"
echo "$ALL_VARIABLES" | head -20 >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Check for common variable name inconsistencies
inconsistency_errors=0

# Extract just the variable names (without template syntax)
VARIABLE_NAMES=$(echo "$ALL_VARIABLES" | sed 's/{{\.//g' | sed 's/}}//g' | sed 's/{{//g' | grep -v '^[[:space:]]*$' | sort | uniq)

echo "Common variables used across templates:" >> "$LOG_FILE"
for var in ProjectName ModulePath Author Email License GoVersion Framework Logger DatabaseDriver DatabaseORM AuthType DomainName; do
    count=$(echo "$VARIABLE_NAMES" | grep -c "^$var\$" || echo "0")
    echo "  $var: $count occurrences" >> "$LOG_FILE"
done
echo "" >> "$LOG_FILE"

# Check for variable naming inconsistencies
echo "Checking for variable naming inconsistencies..." >> "$LOG_FILE"
for file in $TEMPLATE_FILES; do
    # Check for inconsistent project name variables
    if grep -q '{{\.Project}}' "$file" && grep -q '{{\.ProjectName}}' "$file"; then
        echo "INCONSISTENCY: Mixed Project/ProjectName variables in $file" >> "$LOG_FILE"
        echo "$file: Mixed Project/ProjectName variables" >> "$ERROR_LOG"
        ((inconsistency_errors++))
    fi
    
    # Check for inconsistent module path variables
    if grep -q '{{\.Module}}' "$file" && grep -q '{{\.ModulePath}}' "$file"; then
        echo "INCONSISTENCY: Mixed Module/ModulePath variables in $file" >> "$LOG_FILE"
        echo "$file: Mixed Module/ModulePath variables" >> "$ERROR_LOG"
        ((inconsistency_errors++))
    fi
done

echo "Total inconsistency errors found: $inconsistency_errors" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 3. Conditional logic analysis
echo "=== 3. CONDITIONAL LOGIC ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

conditional_errors=0
for file in $TEMPLATE_FILES; do
    # Check for unmatched if/end pairs
    if_count=$(grep -c '{{if' "$file" || echo "0")
    end_count=$(grep -c '{{end}}' "$file" || echo "0")
    
    if [ "$if_count" -ne "$end_count" ]; then
        echo "CONDITIONAL ERROR: Unmatched if/end in $file" >> "$LOG_FILE"
        echo "  {{if}} count: $if_count, {{end}} count: $end_count" >> "$LOG_FILE"
        echo "$file: Unmatched if/end blocks" >> "$ERROR_LOG"
        ((conditional_errors++))
    fi
    
    # Check for complex conditional logic that might be error-prone
    if grep -qE '\{\{if.*and.*or.*\}\}' "$file"; then
        echo "COMPLEXITY WARNING: Complex conditional logic in $file" >> "$LOG_FILE"
        grep -nE '\{\{if.*and.*or.*\}\}' "$file" >> "$LOG_FILE"
    fi
done

echo "Total conditional logic errors found: $conditional_errors" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 4. Check for undefined variables
echo "=== 4. UNDEFINED VARIABLE ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Read template.yaml files to get defined variables
DEFINED_VARS=""
for yaml_file in $(find "$TEMPLATE_DIR" -name "template.yaml"); do
    if [ -f "$yaml_file" ]; then
        yaml_vars=$(grep -A 100 "variables:" "$yaml_file" | grep "name:" | sed 's/.*name: *"\([^"]*\)".*/\1/')
        DEFINED_VARS="$DEFINED_VARS $yaml_vars"
    fi
done

undefined_errors=0
echo "Checking for undefined variables..." >> "$LOG_FILE"

# Common template variables that should be defined
EXPECTED_VARS="ProjectName ModulePath Author Email License GoVersion Framework Logger DatabaseDriver DatabaseORM AuthType"

for file in $TEMPLATE_FILES; do
    for expected_var in $EXPECTED_VARS; do
        if grep -q "{{.*\.$expected_var" "$file"; then
            # Check if this variable is defined in the corresponding template.yaml
            template_dir=$(dirname "$file")
            while [ "$template_dir" != "$TEMPLATE_DIR" ] && [ "$template_dir" != "/" ]; do
                if [ -f "$template_dir/template.yaml" ]; then
                    if ! grep -q "name: *\"$expected_var\"" "$template_dir/template.yaml"; then
                        echo "UNDEFINED VARIABLE: $expected_var used in $file but not defined in $template_dir/template.yaml" >> "$LOG_FILE"
                        echo "$file: Undefined variable $expected_var" >> "$ERROR_LOG"
                        ((undefined_errors++))
                    fi
                    break
                fi
                template_dir=$(dirname "$template_dir")
            done
        fi
    done
done

echo "Total undefined variable errors found: $undefined_errors" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 5. File-specific checks
echo "=== 5. FILE-SPECIFIC ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

file_errors=0

# Check go.mod.tmpl files
echo "Checking go.mod.tmpl files..." >> "$LOG_FILE"
for file in $(find "$TEMPLATE_DIR" -name "go.mod.tmpl"); do
    if ! grep -q '{{\.ModulePath}}' "$file"; then
        echo "ERROR: go.mod.tmpl missing ModulePath variable: $file" >> "$LOG_FILE"
        echo "$file: Missing ModulePath in go.mod.tmpl" >> "$ERROR_LOG"
        ((file_errors++))
    fi
    
    if ! grep -q '{{\.GoVersion}}' "$file"; then
        echo "ERROR: go.mod.tmpl missing GoVersion variable: $file" >> "$LOG_FILE"
        echo "$file: Missing GoVersion in go.mod.tmpl" >> "$ERROR_LOG"
        ((file_errors++))
    fi
done

# Check main.go.tmpl files for import consistency
echo "Checking main.go.tmpl files..." >> "$LOG_FILE"
for file in $(find "$TEMPLATE_DIR" -name "main.go.tmpl"); do
    if ! grep -q '{{\.ModulePath}}' "$file"; then
        echo "ERROR: main.go.tmpl missing ModulePath in imports: $file" >> "$LOG_FILE"
        echo "$file: Missing ModulePath in main.go.tmpl imports" >> "$ERROR_LOG"
        ((file_errors++))
    fi
done

echo "Total file-specific errors found: $file_errors" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# 6. Security analysis
echo "=== 6. SECURITY ANALYSIS ===" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

security_issues=0

for file in $TEMPLATE_FILES; do
    # Check for potential code injection
    if grep -qE '\{\{.*\|.*html.*\}\}|\{\{.*\|.*js.*\}\}|\{\{.*\|.*css.*\}\}' "$file"; then
        echo "SECURITY: Potential XSS vulnerability in $file" >> "$LOG_FILE"
        grep -nE '\{\{.*\|.*html.*\}\}|\{\{.*\|.*js.*\}\}|\{\{.*\|.*css.*\}\}' "$file" >> "$LOG_FILE"
        ((security_issues++))
    fi
    
    # Check for file path traversal
    if grep -qE '\{\{.*\.\./.*\}\}|\{\{.*/\.\./.*\}\}' "$file"; then
        echo "SECURITY: Potential path traversal in $file" >> "$LOG_FILE"
        grep -nE '\{\{.*\.\./.*\}\}|\{\{.*/\.\./.*\}\}' "$file" >> "$LOG_FILE"
        ((security_issues++))
    fi
done

echo "Total security issues found: $security_issues" >> "$LOG_FILE"
echo "" >> "$LOG_FILE"

# Summary
echo "=== SUMMARY ===" >> "$LOG_FILE"
echo "Total template files analyzed: $TOTAL_FILES" >> "$LOG_FILE"
echo "Syntax errors: $syntax_errors" >> "$LOG_FILE"
echo "Variable inconsistencies: $inconsistency_errors" >> "$LOG_FILE"
echo "Conditional logic errors: $conditional_errors" >> "$LOG_FILE"
echo "Undefined variable errors: $undefined_errors" >> "$LOG_FILE"
echo "File-specific errors: $file_errors" >> "$LOG_FILE"
echo "Security issues: $security_issues" >> "$LOG_FILE"

total_errors=$((syntax_errors + inconsistency_errors + conditional_errors + undefined_errors + file_errors + security_issues))
echo "TOTAL ERRORS FOUND: $total_errors" >> "$LOG_FILE"

echo "Analysis complete. Results saved to $LOG_FILE"
echo "Error summary saved to $ERROR_LOG"

# Output to console as well
cat "$LOG_FILE"