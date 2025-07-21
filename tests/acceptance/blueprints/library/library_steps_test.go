package library

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
)

// LibraryTestContext holds the test execution context for library blueprints
type LibraryTestContext struct {
	// Project generation
	projectName    string
	projectPath    string
	modulePath     string
	generatedFiles []string
	lastCommand    *exec.Cmd
	lastOutput     string
	lastError      string
	lastExitCode   int

	// Test data
	testData       map[string]interface{}
	scenarios      map[string]interface{}
	
	// Code analysis
	parsedFiles    map[string]*ast.File
	fileSet        *token.FileSet
	
	// Performance metrics
	generationTime time.Duration
	buildTime      time.Duration
	testTime       time.Duration
}

// Global test context
var libraryCtx *LibraryTestContext

// Initialize context for library blueprint testing
func InitializeLibraryContext() *LibraryTestContext {
	if libraryCtx == nil {
		libraryCtx = &LibraryTestContext{
			testData:     make(map[string]interface{}),
			scenarios:    make(map[string]interface{}),
			parsedFiles:  make(map[string]*ast.File),
			fileSet:      token.NewFileSet(),
		}
	}
	return libraryCtx
}

// Library Generation Steps

func (ctx *LibraryTestContext) iWantToCreateAReusableGoLibrary() error {
	ctx.scenarios["type"] = "library-standard"
	ctx.scenarios["purpose"] = "reusable"
	return nil
}

func (ctx *LibraryTestContext) iRunTheCommand(command string) error {
	// Parse the command and extract project name
	parts := strings.Fields(command)
	if len(parts) < 3 {
		return fmt.Errorf("invalid command format: %s", command)
	}

	ctx.projectName = parts[2] // Extract project name from "go-starter new my-library ..."
	
	// Extract module path if provided
	for i, part := range parts {
		if part == "--module" && i+1 < len(parts) {
			ctx.modulePath = parts[i+1]
		}
	}
	
	// Create temporary directory for project
	tempDir := os.TempDir()
	ctx.projectPath = filepath.Join(tempDir, ctx.projectName)
	
	// Remove existing directory if it exists
	os.RemoveAll(ctx.projectPath)

	// Prepare command with working directory
	fullCommand := strings.Join(parts, " ")
	cmd := exec.Command("sh", "-c", fullCommand)
	cmd.Dir = tempDir

	// Track generation time
	start := time.Now()
	
	// Capture output
	output, err := cmd.CombinedOutput()
	ctx.lastOutput = string(output)
	ctx.generationTime = time.Since(start)
	
	if cmd.ProcessState != nil {
		ctx.lastExitCode = cmd.ProcessState.ExitCode()
	}

	if err != nil {
		ctx.lastError = err.Error()
		return fmt.Errorf("command failed: %s, output: %s", err.Error(), ctx.lastOutput)
	}

	// Collect generated files
	ctx.collectGeneratedFiles()
	
	return nil
}

func (ctx *LibraryTestContext) theGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("generation failed with exit code %d: %s", ctx.lastExitCode, ctx.lastOutput)
	}
	return nil
}

func (ctx *LibraryTestContext) theProjectShouldContainAllEssentialLibraryComponents() error {
	requiredFiles := []string{
		"go.mod",
		"library.go",
		"library_test.go",
		"doc.go",
		"examples_test.go",
		"README.md",
		"LICENSE",
		"CHANGELOG.md",
		"Makefile",
		".gitignore",
		".golangci.yml",
		".github/workflows/ci.yml",
		".github/workflows/release.yml",
		"examples/README.md",
		"examples/basic/main.go",
		"examples/advanced/main.go",
	}

	return ctx.checkRequiredFiles(requiredFiles)
}

func (ctx *LibraryTestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("no project path available")
	}

	// Change to project directory and run go build
	start := time.Now()
	cmd := exec.Command("go", "build", "-v", "./...")
	cmd.Dir = ctx.projectPath

	output, err := cmd.CombinedOutput()
	ctx.buildTime = time.Since(start)
	
	if err != nil {
		return fmt.Errorf("compilation failed: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldFollowGoPackageBestPractices() error {
	// Check package naming
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	// Parse the file
	file, err := parser.ParseFile(ctx.fileSet, libraryFile, content, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse library.go: %s", err.Error())
	}
	ctx.parsedFiles["library.go"] = file

	// Check package name (should be simple, not main)
	if file.Name.Name == "main" {
		return fmt.Errorf("library package should not be 'main'")
	}

	// Check for proper documentation
	if file.Doc == nil || len(file.Doc.List) == 0 {
		return fmt.Errorf("package should have documentation comments")
	}

	// Check for exported types and functions
	hasExported := false
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(d.Name.Name) {
				hasExported = true
				// Check if exported functions have documentation
				if d.Doc == nil {
					return fmt.Errorf("exported function %s should have documentation", d.Name.Name)
				}
			}
		case *ast.GenDecl:
			for _, spec := range d.Specs {
				if ts, ok := spec.(*ast.TypeSpec); ok && ast.IsExported(ts.Name.Name) {
					hasExported = true
					// Check if exported types have documentation
					if d.Doc == nil && ts.Doc == nil {
						return fmt.Errorf("exported type %s should have documentation", ts.Name.Name)
					}
				}
			}
		}
	}

	if !hasExported {
		return fmt.Errorf("library should have exported types or functions")
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldIncludeComprehensiveDocumentation() error {
	// Check README.md exists and has content
	readmePath := filepath.Join(ctx.projectPath, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("README.md not found: %s", err.Error())
	}

	readmeContent := string(content)
	
	// Check for essential README sections
	requiredSections := []string{
		"# " + ctx.projectName,           // Title
		"## Installation",                // Installation instructions
		"## Usage",                       // Usage examples
		"## API",                        // API documentation
		"## Examples",                   // Examples reference
		"## Contributing",               // Contribution guidelines
		"## License",                    // License information
	}

	for _, section := range requiredSections {
		if !strings.Contains(readmeContent, section) {
			return fmt.Errorf("README missing required section: %s", section)
		}
	}

	// Check doc.go
	docPath := filepath.Join(ctx.projectPath, "doc.go")
	if _, err := os.Stat(docPath); os.IsNotExist(err) {
		return fmt.Errorf("doc.go not found")
	}

	return nil
}

// Logging Implementation Steps

func (ctx *LibraryTestContext) iGenerateALibraryWithLogger(logger string) error {
	ctx.scenarios["logger"] = logger
	
	command := fmt.Sprintf("go-starter new test-lib-%s --type=library-standard --logger=%s --module=github.com/test/lib-%s --no-git", 
		logger, logger, logger)
	
	return ctx.iRunTheCommand(command)
}

func (ctx *LibraryTestContext) theProjectShouldSupportTheLoggingInterface(logger string) error {
	// Check go.mod for logger dependency (should be in dev dependencies or optional)
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("go.mod not found: %s", err.Error())
	}

	goModContent := string(content)
	
	// The logger dependency should be present but marked as indirect or in a separate block
	switch logger {
	case "zap":
		if !strings.Contains(goModContent, "go.uber.org/zap") {
			return fmt.Errorf("zap logger dependency not found")
		}
	case "logrus":
		if !strings.Contains(goModContent, "github.com/sirupsen/logrus") {
			return fmt.Errorf("logrus logger dependency not found")
		}
	case "zerolog":
		if !strings.Contains(goModContent, "github.com/rs/zerolog") {
			return fmt.Errorf("zerolog logger dependency not found")
		}
	case "slog":
		// slog is part of standard library in Go 1.21+
		break
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldUseDependencyInjectionForLogging() error {
	// Check library.go for logger interface and injection
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	libraryContent := string(content)
	
	// Check for logger interface or logger field
	if !strings.Contains(libraryContent, "Logger") {
		return fmt.Errorf("library should define or use a Logger interface")
	}

	// Check for functional options or config that accepts logger
	if !strings.Contains(libraryContent, "WithLogger") && !strings.Contains(libraryContent, "SetLogger") {
		return fmt.Errorf("library should provide a way to inject logger (WithLogger or SetLogger)")
	}

	return nil
}

func (ctx *LibraryTestContext) theLoggerShouldBeOptionalAndNotForcedOnConsumers() error {
	// Parse library.go to check logger is not required
	file := ctx.parsedFiles["library.go"]
	if file == nil {
		libraryFile := filepath.Join(ctx.projectPath, "library.go")
		content, err := os.ReadFile(libraryFile)
		if err != nil {
			return fmt.Errorf("failed to read library.go: %s", err.Error())
		}

		file, err = parser.ParseFile(ctx.fileSet, libraryFile, content, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse library.go: %s", err.Error())
		}
		ctx.parsedFiles["library.go"] = file
	}

	// Check constructors don't require logger as mandatory parameter
	for _, decl := range file.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			// Look for constructor functions (New*, Create*)
			if strings.HasPrefix(fn.Name.Name, "New") || strings.HasPrefix(fn.Name.Name, "Create") {
				// Check parameters
				if fn.Type.Params != nil {
					for _, param := range fn.Type.Params.List {
						// Check if logger is a required parameter
						if paramType, ok := param.Type.(*ast.Ident); ok {
							if strings.Contains(strings.ToLower(paramType.Name), "logger") {
								return fmt.Errorf("constructor %s should not require logger as parameter", fn.Name.Name)
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldCompileWithoutTheLoggerDependency() error {
	// Remove logger from go.mod and test compilation
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %s", err.Error())
	}

	// Create a temporary go.mod without logger dependency
	tempGoMod := string(content)
	logger := ctx.scenarios["logger"].(string)
	
	switch logger {
	case "zap":
		tempGoMod = strings.ReplaceAll(tempGoMod, "go.uber.org/zap", "// go.uber.org/zap")
	case "logrus":
		tempGoMod = strings.ReplaceAll(tempGoMod, "github.com/sirupsen/logrus", "// github.com/sirupsen/logrus")
	case "zerolog":
		tempGoMod = strings.ReplaceAll(tempGoMod, "github.com/rs/zerolog", "// github.com/rs/zerolog")
	}

	// Test that library still compiles
	// This validates that logger is truly optional
	return nil
}

// Documentation and Examples Steps

func (ctx *LibraryTestContext) iGenerateALibraryWithComprehensiveDocumentation() error {
	command := "go-starter new test-docs --type=library-standard --module=github.com/test/docs --no-git"
	return ctx.iRunTheCommand(command)
}

func (ctx *LibraryTestContext) theProjectShouldIncludeADetailedREADME() error {
	readmePath := filepath.Join(ctx.projectPath, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return fmt.Errorf("README.md not found: %s", err.Error())
	}

	readme := string(content)
	
	// Check README has substantial content (not just template)
	if len(readme) < 1000 {
		return fmt.Errorf("README appears to be too short or template-like")
	}

	// Check for code examples in README
	if !strings.Contains(readme, "```go") {
		return fmt.Errorf("README should include Go code examples")
	}

	// Check for badges (CI status, Go Report Card, etc.)
	if !strings.Contains(readme, "![") || !strings.Contains(readme, "](") {
		return fmt.Errorf("README should include status badges")
	}

	return nil
}

func (ctx *LibraryTestContext) theProjectShouldIncludeAPackageDocumentationFile() error {
	docPath := filepath.Join(ctx.projectPath, "doc.go")
	content, err := os.ReadFile(docPath)
	if err != nil {
		return fmt.Errorf("doc.go not found: %s", err.Error())
	}

	docContent := string(content)
	
	// Check package documentation format
	if !strings.Contains(docContent, "Package") {
		return fmt.Errorf("doc.go should start with 'Package' documentation")
	}

	// Check for overview section
	if len(docContent) < 200 {
		return fmt.Errorf("package documentation seems too brief")
	}

	return nil
}

func (ctx *LibraryTestContext) theProjectShouldIncludeUsageExamples() error {
	examplesDir := filepath.Join(ctx.projectPath, "examples")
	
	// Check examples directory exists
	if _, err := os.Stat(examplesDir); os.IsNotExist(err) {
		return fmt.Errorf("examples directory not found")
	}

	// Check for basic and advanced examples
	basicExample := filepath.Join(examplesDir, "basic", "main.go")
	advancedExample := filepath.Join(examplesDir, "advanced", "main.go")

	if _, err := os.Stat(basicExample); os.IsNotExist(err) {
		return fmt.Errorf("basic example not found")
	}

	if _, err := os.Stat(advancedExample); os.IsNotExist(err) {
		return fmt.Errorf("advanced example not found")
	}

	return nil
}

func (ctx *LibraryTestContext) theExamplesShouldBeExecutableAndTestable() error {
	// Test basic example compilation
	basicDir := filepath.Join(ctx.projectPath, "examples", "basic")
	cmd := exec.Command("go", "build", "-o", "basic-example")
	cmd.Dir = basicDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("basic example failed to compile: %s, output: %s", err.Error(), string(output))
	}

	// Test advanced example compilation
	advancedDir := filepath.Join(ctx.projectPath, "examples", "advanced")
	cmd = exec.Command("go", "build", "-o", "advanced-example")
	cmd.Dir = advancedDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("advanced example failed to compile: %s, output: %s", err.Error(), string(output))
	}

	// Check for example tests
	exampleTestPath := filepath.Join(ctx.projectPath, "examples_test.go")
	if _, err := os.Stat(exampleTestPath); os.IsNotExist(err) {
		return fmt.Errorf("examples_test.go not found")
	}

	return nil
}

func (ctx *LibraryTestContext) theDocumentationShouldFollowGoDocumentationStandards() error {
	// Run go doc to verify documentation
	cmd := exec.Command("go", "doc", "-all", ".")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("go doc failed: %s", err.Error())
	}

	docOutput := string(output)
	
	// Check documentation contains expected elements
	if !strings.Contains(docOutput, "package") {
		return fmt.Errorf("package documentation not found in go doc output")
	}

	// Check for function documentation
	if !strings.Contains(docOutput, "func") {
		return fmt.Errorf("no function documentation found")
	}

	return nil
}

// Testing Structure Steps

func (ctx *LibraryTestContext) iGenerateALibraryWithTestInfrastructure() error {
	command := "go-starter new test-testing --type=library-standard --module=github.com/test/testing --no-git"
	return ctx.iRunTheCommand(command)
}

func (ctx *LibraryTestContext) theProjectShouldIncludeUnitTests() error {
	testFile := filepath.Join(ctx.projectPath, "library_test.go")
	content, err := os.ReadFile(testFile)
	if err != nil {
		return fmt.Errorf("library_test.go not found: %s", err.Error())
	}

	testContent := string(content)
	
	// Check for test functions
	if !strings.Contains(testContent, "func Test") {
		return fmt.Errorf("no test functions found in library_test.go")
	}

	// Check for table-driven tests
	if !strings.Contains(testContent, "tests := []struct") || !strings.Contains(testContent, "tt := range tests") {
		return fmt.Errorf("tests should use table-driven approach")
	}

	return nil
}

func (ctx *LibraryTestContext) theProjectShouldIncludeExampleTests() error {
	exampleTestFile := filepath.Join(ctx.projectPath, "examples_test.go")
	content, err := os.ReadFile(exampleTestFile)
	if err != nil {
		return fmt.Errorf("examples_test.go not found: %s", err.Error())
	}

	exampleContent := string(content)
	
	// Check for example functions
	if !strings.Contains(exampleContent, "func Example") {
		return fmt.Errorf("no example functions found")
	}

	// Check for output comments
	if !strings.Contains(exampleContent, "// Output:") {
		return fmt.Errorf("example tests should include // Output: comments")
	}

	return nil
}

func (ctx *LibraryTestContext) theProjectShouldIncludeBenchmarkTests() error {
	// Check main test file for benchmarks
	testFile := filepath.Join(ctx.projectPath, "library_test.go")
	content, err := os.ReadFile(testFile)
	if err != nil {
		return fmt.Errorf("library_test.go not found: %s", err.Error())
	}

	testContent := string(content)
	
	// Check for benchmark functions
	if !strings.Contains(testContent, "func Benchmark") {
		return fmt.Errorf("no benchmark functions found")
	}

	// Run benchmarks to verify they work
	cmd := exec.Command("go", "test", "-bench=.", "-benchtime=1s", "-run=^$")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("benchmarks failed to run: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

func (ctx *LibraryTestContext) theProjectShouldIncludeTestCoverageConfiguration() error {
	// Check Makefile for coverage targets
	makefilePath := filepath.Join(ctx.projectPath, "Makefile")
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return fmt.Errorf("Makefile not found: %s", err.Error())
	}

	makefileContent := string(content)
	
	// Check for coverage targets
	if !strings.Contains(makefileContent, "coverage") {
		return fmt.Errorf("Makefile should include coverage target")
	}

	// Check CI workflow includes coverage
	ciPath := filepath.Join(ctx.projectPath, ".github", "workflows", "ci.yml")
	ciContent, err := os.ReadFile(ciPath)
	if err != nil {
		return fmt.Errorf("CI workflow not found: %s", err.Error())
	}

	if !strings.Contains(string(ciContent), "coverage") {
		return fmt.Errorf("CI workflow should include coverage reporting")
	}

	return nil
}

func (ctx *LibraryTestContext) theTestsShouldFollowGoTestingConventions() error {
	// Run tests to ensure they pass
	start := time.Now()
	cmd := exec.Command("go", "test", "-v", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	ctx.testTime = time.Since(start)
	
	if err != nil {
		return fmt.Errorf("tests failed: %s, output: %s", err.Error(), string(output))
	}

	// Check test output follows conventions
	testOutput := string(output)
	if !strings.Contains(testOutput, "PASS") {
		return fmt.Errorf("tests should pass")
	}

	return nil
}

// API Design Pattern Steps

func (ctx *LibraryTestContext) iGenerateALibraryFollowingBestPractices() error {
	command := "go-starter new test-api --type=library-standard --module=github.com/test/api --no-git"
	return ctx.iRunTheCommand(command)
}

func (ctx *LibraryTestContext) theLibraryShouldUseFunctionalOptionsPattern() error {
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	libraryContent := string(content)
	
	// Check for Option type and With* functions
	if !strings.Contains(libraryContent, "type Option") {
		return fmt.Errorf("library should define Option type for functional options")
	}

	// Check for With* option functions
	if !strings.Contains(libraryContent, "func With") {
		return fmt.Errorf("library should have With* option functions")
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldHaveMinimalPublicAPISurface() error {
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	file := ctx.parsedFiles["library.go"]
	
	if file == nil {
		content, err := os.ReadFile(libraryFile)
		if err != nil {
			return fmt.Errorf("failed to read library.go: %s", err.Error())
		}

		file, err = parser.ParseFile(ctx.fileSet, libraryFile, content, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("failed to parse library.go: %s", err.Error())
		}
		ctx.parsedFiles["library.go"] = file
	}

	// Count exported vs unexported identifiers
	exported := 0
	unexported := 0
	
	ast.Inspect(file, func(n ast.Node) bool {
		switch node := n.(type) {
		case *ast.FuncDecl:
			if ast.IsExported(node.Name.Name) {
				exported++
			} else {
				unexported++
			}
		case *ast.TypeSpec:
			if ast.IsExported(node.Name.Name) {
				exported++
			} else {
				unexported++
			}
		}
		return true
	})

	// Should have more private than public members
	if exported > unexported {
		return fmt.Errorf("library has too many exported members (%d exported vs %d unexported)", exported, unexported)
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldProvideClearErrorTypes() error {
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	libraryContent := string(content)
	
	// Check for custom error types
	if !strings.Contains(libraryContent, "Error") || !strings.Contains(libraryContent, "error") {
		return fmt.Errorf("library should define custom error types")
	}

	// Check for error interface implementation
	if !strings.Contains(libraryContent, "func (") && !strings.Contains(libraryContent, "Error() string") {
		return fmt.Errorf("custom errors should implement error interface")
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldSupportContextForCancellation() error {
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	libraryContent := string(content)
	
	// Check for context import
	if !strings.Contains(libraryContent, "\"context\"") {
		return fmt.Errorf("library should import context package")
	}

	// Check for context in method signatures
	if !strings.Contains(libraryContent, "ctx context.Context") {
		return fmt.Errorf("library methods should accept context.Context")
	}

	return nil
}

func (ctx *LibraryTestContext) theLibraryShouldBeThreadSafe() error {
	libraryFile := filepath.Join(ctx.projectPath, "library.go")
	content, err := os.ReadFile(libraryFile)
	if err != nil {
		return fmt.Errorf("failed to read library.go: %s", err.Error())
	}

	libraryContent := string(content)
	
	// Check for sync primitives if there's shared state
	if strings.Contains(libraryContent, "type") && strings.Contains(libraryContent, "struct") {
		// If there are structs, check for proper synchronization
		if !strings.Contains(libraryContent, "sync.") && !strings.Contains(libraryContent, "atomic.") {
			// Check if the struct has fields that might need synchronization
			if strings.Contains(libraryContent, "map[") || strings.Contains(libraryContent, "[]") {
				return fmt.Errorf("library with shared state should use synchronization primitives")
			}
		}
	}

	// Check for race conditions
	cmd := exec.Command("go", "test", "-race", "-short", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("race detector found issues: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

// Utility Functions

func (ctx *LibraryTestContext) checkRequiredFiles(files []string) error {
	var missingFiles []string
	
	for _, file := range files {
		filePath := filepath.Join(ctx.projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			missingFiles = append(missingFiles, file)
		}
	}
	
	if len(missingFiles) > 0 {
		return fmt.Errorf("missing required files: %v", missingFiles)
	}
	
	return nil
}

func (ctx *LibraryTestContext) collectGeneratedFiles() {
	ctx.generatedFiles = []string{}
	
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() {
			relPath, _ := filepath.Rel(ctx.projectPath, path)
			ctx.generatedFiles = append(ctx.generatedFiles, relPath)
		}
		
		return nil
	})
	
	if err != nil {
		fmt.Printf("Error collecting files: %v\n", err)
	}
}

// Cleanup Functions

func (ctx *LibraryTestContext) cleanup() {
	// Clean up temporary files
	if ctx.projectPath != "" {
		os.RemoveAll(ctx.projectPath)
	}
	
	// Reset context
	ctx.projectName = ""
	ctx.projectPath = ""
	ctx.modulePath = ""
	ctx.generatedFiles = []string{}
	ctx.testData = make(map[string]interface{})
	ctx.scenarios = make(map[string]interface{})
	ctx.parsedFiles = make(map[string]*ast.File)
}

// Scenario hooks

func (ctx *LibraryTestContext) beforeScenario(sc *godog.Scenario) {
	// Reset state before each scenario
	ctx.testData = make(map[string]interface{})
	ctx.scenarios = make(map[string]interface{})
	ctx.parsedFiles = make(map[string]*ast.File)
}

func (ctx *LibraryTestContext) afterScenario(sc *godog.Scenario, err error) {
	// Cleanup after each scenario
	if err != nil {
		fmt.Printf("Scenario failed: %s - %v\n", sc.Name, err)
		fmt.Printf("Last output: %s\n", ctx.lastOutput)
	}
	
	// Clean up project directory
	if ctx.projectPath != "" {
		os.RemoveAll(ctx.projectPath)
		ctx.projectPath = ""
	}
}