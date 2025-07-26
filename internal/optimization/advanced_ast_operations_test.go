package optimization

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAdvancedASTOperations(t *testing.T) {
	options := DefaultAdvancedTransformOptions()
	ops := NewAdvancedASTOperations(options)
	
	assert.NotNil(t, ops)
	assert.NotNil(t, ops.fileSet)
	assert.Equal(t, options, ops.options)
}

func TestDefaultAdvancedTransformOptions(t *testing.T) {
	options := DefaultAdvancedTransformOptions()
	
	// Conservative defaults
	assert.True(t, options.ExtractComplexExpressions)
	assert.False(t, options.InlineSimpleFunctions)
	assert.True(t, options.OptimizeControlFlow)
	assert.True(t, options.PromoteStringBuilder)
	assert.True(t, options.OptimizeLoops)
	assert.False(t, options.CacheExpensiveOperations)
	
	// Safety settings
	assert.True(t, options.PreserveSemantics)
	assert.True(t, options.RequireExplicitApproval)
	assert.Equal(t, 10, options.MaxTransformationsPerFile)
	
	// Architecture optimization
	assert.Equal(t, "standard", options.OptimizeForArchitecture)
}

func TestTransformCode_ComplexExpressions(t *testing.T) {
	testCases := []struct {
		name                 string
		code                 string
		expectedTransformations []string
	}{
		{
			name: "complex if condition",
			code: `package main

func main() {
	if (x > 0 && y < 10) && (z == 5 || w != 3) && len(items) > 0 {
		println("complex")
	}
}`,
			expectedTransformations: []string{"extract_complex_condition"},
		},
		{
			name: "deeply nested function call",
			code: `package main

func main() {
	result := process(transform(validate(getData())))
}`,
			expectedTransformations: []string{"extract_nested_call"},
		},
		{
			name: "simple condition - no transformation",
			code: `package main

func main() {
	if x > 0 {
		println("simple")
	}
}`,
			expectedTransformations: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			options.ExtractComplexExpressions = true
			// Disable other optimizations to isolate complex expression extraction
			options.OptimizeLoops = false
			options.PromoteStringBuilder = false
			options.OptimizeControlFlow = false
			options.OptimizeErrorHandling = false
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			require.NoError(t, err)
			
			var transformationTypes []string
			for _, trans := range result.Transformations {
				transformationTypes = append(transformationTypes, trans.Type)
			}
			
			// Check that the expected transformations are present
			for _, expected := range tc.expectedTransformations {
				assert.Contains(t, transformationTypes, expected, "Should contain transformation type: %s", expected)
			}
			
			// If no transformations expected, verify none were found
			if len(tc.expectedTransformations) == 0 {
				assert.Empty(t, transformationTypes, "Should not find any transformations")
			}
		})
	}
}

func TestTransformCode_ControlFlowOptimization(t *testing.T) {
	testCases := []struct {
		name                 string
		code                 string
		expectedTransformations []string
	}{
		{
			name: "early return opportunity",
			code: `package main

func validate(x int) error {
	if x > 0 {
		// do something
	} else {
		return errors.New("invalid")
	}
	return nil
}`,
			expectedTransformations: []string{"early_return"},
		},
		{
			name: "complex switch statement",
			code: `package main

func handleType(t string) {
	switch t {
	case "a": println("a")
	case "b": println("b")
	case "c": println("c")
	case "d": println("d")
	case "e": println("e")
	case "f": println("f")
	default: println("default")
	}
}`,
			expectedTransformations: []string{"optimize_switch"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			options.OptimizeControlFlow = true
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			require.NoError(t, err)
			
			var transformationTypes []string
			for _, trans := range result.Transformations {
				transformationTypes = append(transformationTypes, trans.Type)
			}
			
			assert.ElementsMatch(t, tc.expectedTransformations, transformationTypes)
		})
	}
}

func TestTransformCode_StringBuilderPromotion(t *testing.T) {
	testCases := []struct {
		name                 string
		code                 string
		expectedTransformations []string
	}{
		{
			name: "string concatenation in for loop",
			code: `package main

func buildString() string {
	var result string
	for i := 0; i < 10; i++ {
		result += "item"
	}
	return result
}`,
			expectedTransformations: []string{"promote_string_builder"},
		},
		{
			name: "string concatenation in range loop",
			code: `package main

func joinItems(items []string) string {
	var result string
	for _, item := range items {
		result += item
	}
	return result
}`,
			expectedTransformations: []string{"promote_string_builder"},
		},
		{
			name: "no string concatenation",
			code: `package main

func simple() {
	for i := 0; i < 10; i++ {
		println(i)
	}
}`,
			expectedTransformations: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			options.PromoteStringBuilder = true
			// Disable other optimizations to isolate string builder promotion
			options.OptimizeLoops = false
			options.ExtractComplexExpressions = false
			options.OptimizeControlFlow = false
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			require.NoError(t, err)
			
			var transformationTypes []string
			for _, trans := range result.Transformations {
				transformationTypes = append(transformationTypes, trans.Type)
			}
			
			// Check that the expected transformations are present
			for _, expected := range tc.expectedTransformations {
				assert.Contains(t, transformationTypes, expected, "Should contain transformation type: %s", expected)
			}
			
			// If no transformations expected, verify none were found
			if len(tc.expectedTransformations) == 0 {
				assert.Empty(t, transformationTypes, "Should not find any transformations")
			}
		})
	}
}

func TestTransformCode_LoopOptimization(t *testing.T) {
	testCases := []struct {
		name                 string
		code                 string
		expectedTransformations []string
	}{
		{
			name: "repeated len() call in loop",
			code: `package main

func process(items []string) {
	for i := 0; i < len(items); i++ {
		println(items[i])
	}
}`,
			expectedTransformations: []string{"cache_loop_len"},
		},
		{
			name: "unused range variable",
			code: `package main

func process(items []string) {
	for i, v := range items {
		println(v)
		// i is not used
	}
}`,
			expectedTransformations: []string{"optimize_range_vars"},
		},
		{
			name: "efficient loop - no optimization",
			code: `package main

func process(items []string) {
	n := len(items)
	for i := 0; i < n; i++ {
		println(items[i])
	}
}`,
			expectedTransformations: []string{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			options.OptimizeLoops = true
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			require.NoError(t, err)
			
			var transformationTypes []string
			for _, trans := range result.Transformations {
				transformationTypes = append(transformationTypes, trans.Type)
			}
			
			assert.ElementsMatch(t, tc.expectedTransformations, transformationTypes)
		})
	}
}

func TestTransformCode_ErrorHandlingOptimization(t *testing.T) {
	testCode := `package main

func example() error {
	err := doSomething()
	if err != nil {
		return err
	}
	
	err = doAnotherThing()
	if err != nil {
		return err
	}
	
	err = doThirdThing()
	if err != nil {
		return err
	}
	
	return nil
}`

	options := DefaultAdvancedTransformOptions()
	options.OptimizeErrorHandling = true
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Should find repeated error handling patterns
	var consolidationTransforms int
	for _, trans := range result.Transformations {
		if trans.Type == "consolidate_error_handling" {
			consolidationTransforms++
		}
	}
	
	assert.Greater(t, consolidationTransforms, 0, "Should find repeated error handling patterns")
}

func TestTransformCode_ArchitecturalOptimizations(t *testing.T) {
	testCases := []struct {
		name         string
		architecture string
		code         string
		expectedType string
	}{
		{
			name:         "clean architecture - direct database call",
			architecture: "clean",
			code: `package handlers

func GetUser(w http.ResponseWriter, r *http.Request) {
	user, err := db.Query("SELECT * FROM users WHERE id = ?", userID)
	// Direct database access in handler violates Clean Architecture
}`,
			expectedType: "enforce_clean_architecture",
		},
		{
			name:         "hexagonal architecture - missing port interface",
			architecture: "hexagonal",
			code: `package domain

type UserService struct {
	UserRepository UserRepository
	EmailService   EmailService
	NotificationClient NotificationClient
}`,
			expectedType: "create_port_interface",
		},
		{
			name:         "ddd - anemic domain model",
			architecture: "ddd",
			code: `package domain

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
	CreatedAt time.Time
}

// No methods - anemic domain model`,
			expectedType: "enrich_domain_model",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			options.ApplyContextualRules = true
			options.OptimizeForArchitecture = tc.architecture
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			require.NoError(t, err)
			
			var foundExpectedType bool
			for _, trans := range result.Transformations {
				if trans.Type == tc.expectedType {
					foundExpectedType = true
					break
				}
			}
			
			assert.True(t, foundExpectedType, "Should find %s transformation", tc.expectedType)
		})
	}
}

func TestTransformCode_QualityMetrics(t *testing.T) {
	testCode := `package main

func complexFunction(x, y, z int) int {
	if x > 0 {
		if y > 0 {
			if z > 0 {
				return x + y + z
			}
			return x + y
		}
		return x
	}
	return 0
}

func simpleFunction() {
	println("simple")
}`

	options := DefaultAdvancedTransformOptions()
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Verify quality metrics are calculated
	assert.Greater(t, result.QualityMetrics.CyclomaticComplexity, 0)
	assert.Greater(t, result.QualityMetrics.FunctionCount, 0)
	assert.Greater(t, result.QualityMetrics.LinesOfCode, 0)
	assert.Greater(t, result.QualityMetrics.TestCoverage, 0.0)
	assert.Greater(t, result.QualityMetrics.Maintainability, 0.0)
}

func TestTransformCode_SafetyValidation(t *testing.T) {
	testCode := `package main

func example() {
	println("test")
}`

	options := DefaultAdvancedTransformOptions()
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Verify safety validation
	assert.True(t, result.SafetyValidation.SemanticsPreserved)
	assert.True(t, result.SafetyValidation.TypeSafetyMaintained)
	assert.True(t, result.SafetyValidation.ErrorHandlingIntact)
	assert.True(t, result.SafetyValidation.TestsStillPass)
	assert.Equal(t, "neutral_or_positive", result.SafetyValidation.PerformanceImpact)
}

func TestTransformCode_ErrorHandling(t *testing.T) {
	testCases := []struct {
		name          string
		code          string
		shouldError   bool
		expectedError string
	}{
		{
			name:          "invalid Go syntax",
			code:          "invalid go code {{{",
			shouldError:   true,
			expectedError: "failed to parse source code",
		},
		{
			name: "valid code",
			code: `package main
func main() { println("valid") }`,
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			options := DefaultAdvancedTransformOptions()
			ops := NewAdvancedASTOperations(options)
			
			result, err := ops.TransformCode(tc.code)
			
			if tc.shouldError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.NotEmpty(t, result.TransformedCode)
			}
		})
	}
}

func TestCalculateExpressionComplexity(t *testing.T) {
	testCases := []struct {
		name       string
		expression string
		minComplexity int
	}{
		{
			name:       "simple variable",
			expression: "x",
			minComplexity: 0,
		},
		{
			name:       "binary expression",
			expression: "x + y",
			minComplexity: 1,
		},
		{
			name:       "nested binary expression",
			expression: "(x + y) && (z > 0)",
			minComplexity: 1,
		},
		{
			name:       "function call",
			expression: "process(x, y)",
			minComplexity: 2,
		},
		{
			name:       "complex expression",
			expression: "process(transform(x + y), validate(z && w))",
			minComplexity: 5,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test code with the expression
			testCode := `package main
func test() {
	result := ` + tc.expression + `
	_ = result
}`

			// Parse and find the expression
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
			require.NoError(t, err)
			
			ops := NewAdvancedASTOperations(DefaultAdvancedTransformOptions())
			ops.fileSet = fset
			
			// Find the assignment statement and extract the expression
			found := false
			for _, decl := range file.Decls {
				if funcDecl, ok := decl.(*ast.FuncDecl); ok {
					for _, stmt := range funcDecl.Body.List {
						if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
							if len(assignStmt.Rhs) > 0 {
								complexity := ops.calculateExpressionComplexity(assignStmt.Rhs[0])
								assert.GreaterOrEqual(t, complexity, tc.minComplexity, 
									"Expression complexity should be at least %d", tc.minComplexity)
								found = true
								break
							}
						}
					}
				}
			}
			
			assert.True(t, found, "Should find the test expression")
		})
	}
}

func TestTransformationConfidenceAndRisk(t *testing.T) {
	testCode := `package main

func example() {
	if (x > 0 && y < 10) && (z == 5 || w != 3) {
		println("complex")
	}
}`

	options := DefaultAdvancedTransformOptions()
	options.ExtractComplexExpressions = true
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Verify transformations have appropriate confidence and risk levels
	for _, trans := range result.Transformations {
		assert.GreaterOrEqual(t, trans.Confidence, 0.0)
		assert.LessOrEqual(t, trans.Confidence, 1.0)
		assert.Contains(t, []string{"safe", "moderate", "risky"}, trans.RiskLevel)
		assert.Contains(t, []string{"low", "medium", "high"}, trans.Impact)
		assert.NotEmpty(t, trans.Description)
		assert.NotEmpty(t, trans.Type)
	}
}

func TestMultipleTransformationTypes(t *testing.T) {
	testCode := `package main

func complexExample(items []string) string {
	var result string
	
	// This should trigger multiple optimizations:
	for i := 0; i < len(items); i++ {           // Loop optimization (cache len)
		if (items[i] != "" && len(items[i]) > 3) && (i % 2 == 0) { // Complex expression
			result += items[i]                   // String builder promotion
		}
	}
	
	return result
}`

	options := DefaultAdvancedTransformOptions()
	options.ExtractComplexExpressions = true
	options.OptimizeLoops = true
	options.PromoteStringBuilder = true
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Should find multiple types of transformations
	transformTypes := make(map[string]int)
	for _, trans := range result.Transformations {
		transformTypes[trans.Type]++
	}
	
	// Should have at least 2-3 different types of transformations
	assert.GreaterOrEqual(t, len(transformTypes), 2, "Should find multiple transformation types")
	
	// Verify specific transformations are found
	expectedTypes := []string{"cache_loop_len", "promote_string_builder"}
	for _, expectedType := range expectedTypes {
		assert.Contains(t, transformTypes, expectedType, "Should find %s transformation", expectedType)
	}
}

func TestTransformCode_Integration(t *testing.T) {
	// Test with a realistic code sample that should trigger multiple optimizations
	testCode := `package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"
)

func GetUsers(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var users []User
	var result string
	
	// Multiple optimization opportunities:
	rows, err := db.Query("SELECT * FROM users") // Direct DB access (Clean Arch violation)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()
	
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {  // Repeated error pattern
			http.Error(w, err.Error(), 500)
			return
		}
		users = append(users, user)
	}
	
	// String concatenation in loop
	for i := 0; i < len(users); i++ {
		if (users[i].Name != "" && len(users[i].Name) > 2) && (users[i].Email != "") {
			result += fmt.Sprintf("User: %s <%s>\n", users[i].Name, users[i].Email)
		}
	}
	
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(result))
}

type User struct {
	ID    int
	Name  string
	Email string
}`

	options := DefaultAdvancedTransformOptions()
	options.ExtractComplexExpressions = true
	options.OptimizeLoops = true
	options.PromoteStringBuilder = true
	options.OptimizeErrorHandling = true
	options.ApplyContextualRules = true
	options.OptimizeForArchitecture = "clean"
	
	ops := NewAdvancedASTOperations(options)
	
	result, err := ops.TransformCode(testCode)
	require.NoError(t, err)
	
	// Should generate a valid transformed code
	assert.NotEmpty(t, result.TransformedCode)
	assert.NotEmpty(t, result.OriginalCode)
	
	// Should have quality metrics
	assert.Greater(t, result.QualityMetrics.FunctionCount, 0)
	assert.Greater(t, result.QualityMetrics.LinesOfCode, 10)
	
	// Should have safety validation
	assert.True(t, result.SafetyValidation.SemanticsPreserved)
	
	// Should find multiple transformation opportunities
	assert.Greater(t, len(result.Transformations), 2, "Should find multiple optimization opportunities")
	
	// Verify no critical errors
	for _, err := range result.Errors {
		// Transformation errors should be non-critical
		assert.NotContains(t, err.Error(), "fatal")
		assert.NotContains(t, err.Error(), "panic")
	}
}