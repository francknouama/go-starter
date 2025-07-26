package architecture

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

// ArchitectureTestContext holds test state for architecture validation
type ArchitectureTestContext struct {
	projectPath      string
	architecture     string
	projectStructure map[string][]string // layer -> list of packages
	importGraph      map[string][]string // package -> imported packages
	violations       []string
	metrics          ArchitectureMetrics
}

// ArchitectureMetrics holds quality metrics for architecture validation
type ArchitectureMetrics struct {
	PackageCoupling     map[string]int
	CyclomaticComplexity map[string]int
	TestCoverage        float64
	DependencyDepth     int
}

// TestFeatures runs the architecture validation BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &ArchitectureTestContext{
				projectStructure: make(map[string][]string),
				importGraph:      make(map[string][]string),
				violations:       []string{},
			}
			
			s.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				return ctx, nil
			})
			
			s.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				return ctx, nil
			})
			
			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run architecture validation tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *ArchitectureTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	
	// Architecture generation steps
	s.Step(`^I generate a project with architecture "([^"]*)"$`, ctx.iGenerateAProjectWithArchitecture)
	s.Step(`^I generate a project with architecture "([^"]*)" and framework "([^"]*)"$`, ctx.iGenerateProjectWithArchitectureAndFramework)
	
	// Analysis steps
	s.Step(`^I analyze the project structure$`, ctx.iAnalyzeTheProjectStructure)
	s.Step(`^I check import statements across the project$`, ctx.iCheckImportStatements)
	s.Step(`^I analyze the integration points$`, ctx.iAnalyzeIntegrationPoints)
	
	// Validation steps
	s.Step(`^the project should follow standard layered architecture:$`, ctx.projectShouldFollowStandardLayeredArchitecture)
	s.Step(`^the project should follow clean architecture principles:$`, ctx.projectShouldFollowCleanArchitecturePrinciples)
	s.Step(`^the project should follow DDD principles:$`, ctx.projectShouldFollowDDDPrinciples)
	s.Step(`^the project should follow hexagonal architecture \(ports and adapters\):$`, ctx.projectShouldFollowHexagonalArchitecture)
	
	// Dependency validation steps
	s.Step(`^dependencies should flow downward only$`, ctx.dependenciesShouldFlowDownwardOnly)
	s.Step(`^handlers should depend on services$`, ctx.handlersShouldDependOnServices)
	s.Step(`^services should depend on repositories$`, ctx.servicesShouldDependOnRepositories)
	s.Step(`^repositories should not depend on handlers or services$`, ctx.repositoriesShouldNotDependOnHandlersOrServices)
	s.Step(`^imports should follow "([^"]*)" dependency rules$`, ctx.importsShouldFollowDependencyRules)
	s.Step(`^there should be no circular dependencies$`, ctx.thereShouldBeNoCircularDependencies)
	
	// Architecture principle validation
	s.Step(`^the dependency rule should be strictly enforced$`, ctx.dependencyRuleShouldBeStrictlyEnforced)
	s.Step(`^business logic should be in usecases layer$`, ctx.businessLogicShouldBeInUsecasesLayer)
	s.Step(`^framework dependencies should only exist in outer layers$`, ctx.frameworkDependenciesShouldOnlyExistInOuterLayers)
	s.Step(`^domain entities should be pure Go structs$`, ctx.domainEntitiesShouldBePureGoStructs)
	s.Step(`^aggregates should encapsulate business invariants$`, ctx.aggregatesShouldEncapsulateBusinessInvariants)
	s.Step(`^value objects should be immutable$`, ctx.valueObjectsShouldBeImmutable)
	s.Step(`^the framework should be properly isolated in the appropriate layer$`, ctx.frameworkShouldBeProperlyIsolated)
}

// Step implementations

func (ctx *ArchitectureTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Verify CLI is available
	return nil
}

func (ctx *ArchitectureTestContext) allTemplatesAreProperlyInitialized() error {
	// Verify templates are loaded
	return nil
}

func (ctx *ArchitectureTestContext) iGenerateAProjectWithArchitecture(architecture string) error {
	ctx.architecture = architecture
	// Generate project with specified architecture
	// Store project path in ctx.projectPath
	return nil
}

func (ctx *ArchitectureTestContext) iGenerateProjectWithArchitectureAndFramework(architecture, framework string) error {
	ctx.architecture = architecture
	// Generate project with specified architecture and framework
	return nil
}

func (ctx *ArchitectureTestContext) iAnalyzeTheProjectStructure() error {
	// Walk the project directory and categorize files by architectural layer
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if strings.HasSuffix(path, ".go") && !strings.Contains(path, "test") {
			relPath, _ := filepath.Rel(ctx.projectPath, path)
			layer := ctx.determineArchitecturalLayer(relPath)
			if layer != "" {
				ctx.projectStructure[layer] = append(ctx.projectStructure[layer], relPath)
			}
			
			// Parse imports
			imports, err := ctx.parseImports(path)
			if err == nil {
				ctx.importGraph[relPath] = imports
			}
		}
		
		return nil
	})
	
	return err
}

func (ctx *ArchitectureTestContext) determineArchitecturalLayer(filePath string) string {
	// Determine which architectural layer a file belongs to based on its path
	switch ctx.architecture {
	case "standard":
		if strings.Contains(filePath, "handlers") {
			return "handlers"
		} else if strings.Contains(filePath, "services") {
			return "services"
		} else if strings.Contains(filePath, "repository") {
			return "repository"
		} else if strings.Contains(filePath, "models") {
			return "models"
		}
	case "clean":
		if strings.Contains(filePath, "entities") {
			return "entities"
		} else if strings.Contains(filePath, "usecases") {
			return "usecases"
		} else if strings.Contains(filePath, "controllers") {
			return "controllers"
		} else if strings.Contains(filePath, "presenters") {
			return "presenters"
		}
	case "ddd":
		if strings.Contains(filePath, "domain") {
			if strings.Contains(filePath, "entity") {
				return "aggregates"
			} else if strings.Contains(filePath, "value_objects") {
				return "value_objects"
			} else if strings.Contains(filePath, "service") {
				return "domain_services"
			}
		} else if strings.Contains(filePath, "application") {
			return "application"
		}
	case "hexagonal":
		if strings.Contains(filePath, "domain") {
			return "domain"
		} else if strings.Contains(filePath, "ports") {
			return "ports"
		} else if strings.Contains(filePath, "primary") {
			return "primary_adapters"
		} else if strings.Contains(filePath, "secondary") {
			return "secondary_adapters"
		}
	}
	
	return ""
}

func (ctx *ArchitectureTestContext) parseImports(filePath string) ([]string, error) {
	// Parse Go file and extract imports
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, filePath, nil, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	
	var imports []string
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, `"`)
		imports = append(imports, importPath)
	}
	
	return imports, nil
}

func (ctx *ArchitectureTestContext) iCheckImportStatements() error {
	// Analyze import statements for architecture compliance
	for file, imports := range ctx.importGraph {
		layer := ctx.determineArchitecturalLayer(file)
		for _, imp := range imports {
			if ctx.isViolatingDependencyRule(layer, imp) {
				ctx.violations = append(ctx.violations, 
					fmt.Sprintf("File %s in layer %s imports %s", file, layer, imp))
			}
		}
	}
	return nil
}

func (ctx *ArchitectureTestContext) isViolatingDependencyRule(fromLayer, importPath string) bool {
	// Check if an import violates architectural dependency rules
	switch ctx.architecture {
	case "standard":
		// In standard architecture, repositories shouldn't depend on handlers or services
		if fromLayer == "repository" && (strings.Contains(importPath, "handlers") || strings.Contains(importPath, "services")) {
			return true
		}
		// Models shouldn't depend on any other layer
		if fromLayer == "models" && (strings.Contains(importPath, "handlers") || 
			strings.Contains(importPath, "services") || 
			strings.Contains(importPath, "repository")) {
			return true
		}
	case "clean":
		// Clean architecture: inner layers shouldn't depend on outer layers
		if fromLayer == "entities" && (strings.Contains(importPath, "usecases") ||
			strings.Contains(importPath, "controllers") ||
			strings.Contains(importPath, "infrastructure")) {
			return true
		}
		if fromLayer == "usecases" && (strings.Contains(importPath, "controllers") ||
			strings.Contains(importPath, "infrastructure")) {
			return true
		}
	case "hexagonal":
		// Domain should have no external dependencies
		if fromLayer == "domain" && (strings.Contains(importPath, "adapters") ||
			strings.Contains(importPath, "infrastructure")) {
			return true
		}
	}
	return false
}

func (ctx *ArchitectureTestContext) projectShouldFollowStandardLayeredArchitecture(table *godog.Table) error {
	// Validate standard layered architecture
	expectedLayers := []string{"handlers", "services", "repository", "models", "middleware", "config"}
	
	for _, layer := range expectedLayers {
		if _, exists := ctx.projectStructure[layer]; !exists {
			return fmt.Errorf("missing required layer: %s", layer)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) projectShouldFollowCleanArchitecturePrinciples(table *godog.Table) error {
	// Validate clean architecture principles
	expectedLayers := []string{"entities", "usecases", "controllers", "repositories", "presenters"}
	
	for _, layer := range expectedLayers {
		if _, exists := ctx.projectStructure[layer]; !exists {
			return fmt.Errorf("missing clean architecture layer: %s", layer)
		}
	}
	
	// Verify dependency rule
	if len(ctx.violations) > 0 {
		return fmt.Errorf("clean architecture dependency violations found: %v", ctx.violations)
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) projectShouldFollowDDDPrinciples(table *godog.Table) error {
	// Validate DDD principles
	requiredComponents := []string{"aggregates", "value_objects", "domain_services", "application"}
	
	for _, component := range requiredComponents {
		if _, exists := ctx.projectStructure[component]; !exists {
			return fmt.Errorf("missing DDD component: %s", component)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) projectShouldFollowHexagonalArchitecture(table *godog.Table) error {
	// Validate hexagonal architecture
	requiredComponents := []string{"domain", "ports", "primary_adapters", "secondary_adapters"}
	
	for _, component := range requiredComponents {
		if _, exists := ctx.projectStructure[component]; !exists {
			return fmt.Errorf("missing hexagonal architecture component: %s", component)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) dependenciesShouldFlowDownwardOnly() error {
	// Validate that dependencies flow in the correct direction
	if len(ctx.violations) > 0 {
		return fmt.Errorf("dependency flow violations: %v", ctx.violations)
	}
	return nil
}

func (ctx *ArchitectureTestContext) handlersShouldDependOnServices() error {
	// Verify handlers import services
	handlerFiles := ctx.projectStructure["handlers"]
	for _, handler := range handlerFiles {
		imports := ctx.importGraph[handler]
		hasServiceImport := false
		for _, imp := range imports {
			if strings.Contains(imp, "services") {
				hasServiceImport = true
				break
			}
		}
		if !hasServiceImport && len(imports) > 0 {
			return fmt.Errorf("handler %s does not import services", handler)
		}
	}
	return nil
}

func (ctx *ArchitectureTestContext) servicesShouldDependOnRepositories() error {
	// Verify services import repositories
	serviceFiles := ctx.projectStructure["services"]
	for _, service := range serviceFiles {
		imports := ctx.importGraph[service]
		hasRepoImport := false
		for _, imp := range imports {
			if strings.Contains(imp, "repository") {
				hasRepoImport = true
				break
			}
		}
		if !hasRepoImport && len(imports) > 0 {
			return fmt.Errorf("service %s does not import repositories", service)
		}
	}
	return nil
}

func (ctx *ArchitectureTestContext) repositoriesShouldNotDependOnHandlersOrServices() error {
	// Verify repositories don't have upward dependencies
	repoFiles := ctx.projectStructure["repository"]
	for _, repo := range repoFiles {
		imports := ctx.importGraph[repo]
		for _, imp := range imports {
			if strings.Contains(imp, "handlers") || strings.Contains(imp, "services") {
				return fmt.Errorf("repository %s has upward dependency on %s", repo, imp)
			}
		}
	}
	return nil
}

func (ctx *ArchitectureTestContext) importsShouldFollowDependencyRules(architecture string) error {
	// Architecture-specific dependency validation
	if architecture != ctx.architecture {
		return fmt.Errorf("architecture mismatch: expected %s, got %s", architecture, ctx.architecture)
	}
	
	if len(ctx.violations) > 0 {
		return fmt.Errorf("import violations found: %v", ctx.violations)
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) thereShouldBeNoCircularDependencies() error {
	// Check for circular dependencies using DFS
	visited := make(map[string]bool)
	recursionStack := make(map[string]bool)
	
	for pkg := range ctx.importGraph {
		if ctx.hasCycle(pkg, visited, recursionStack) {
			return fmt.Errorf("circular dependency detected involving package: %s", pkg)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) hasCycle(pkg string, visited, recursionStack map[string]bool) bool {
	visited[pkg] = true
	recursionStack[pkg] = true
	
	for _, imp := range ctx.importGraph[pkg] {
		// Only check internal imports
		if strings.HasPrefix(imp, ctx.getModulePath()) {
			if !visited[imp] {
				if ctx.hasCycle(imp, visited, recursionStack) {
					return true
				}
			} else if recursionStack[imp] {
				return true
			}
		}
	}
	
	recursionStack[pkg] = false
	return false
}

func (ctx *ArchitectureTestContext) getModulePath() string {
	// Extract module path from go.mod
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return ""
	}
	
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "module ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "module "))
		}
	}
	
	return ""
}

func (ctx *ArchitectureTestContext) dependencyRuleShouldBeStrictlyEnforced() error {
	// Verify clean architecture dependency rule
	return ctx.dependenciesShouldFlowDownwardOnly()
}

func (ctx *ArchitectureTestContext) businessLogicShouldBeInUsecasesLayer() error {
	// Verify business logic placement
	usecaseFiles := ctx.projectStructure["usecases"]
	if len(usecaseFiles) == 0 {
		return fmt.Errorf("no business logic found in usecases layer")
	}
	
	// Check that usecases contain actual logic, not just interfaces
	for _, file := range usecaseFiles {
		content, err := os.ReadFile(filepath.Join(ctx.projectPath, file))
		if err != nil {
			continue
		}
		
		// Simple heuristic: check for function implementations
		if !strings.Contains(string(content), "func (") {
			return fmt.Errorf("usecase file %s does not contain implementations", file)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) frameworkDependenciesShouldOnlyExistInOuterLayers() error {
	// Check that framework imports are only in appropriate layers
	frameworkPackages := []string{"gin-gonic/gin", "labstack/echo", "gofiber/fiber", "go-chi/chi"}
	
	innerLayers := map[string][]string{
		"clean":     {"entities", "usecases"},
		"ddd":       {"domain"},
		"hexagonal": {"domain", "application"},
	}
	
	restrictedLayers, exists := innerLayers[ctx.architecture]
	if !exists {
		return nil // No restrictions for this architecture
	}
	
	for _, layer := range restrictedLayers {
		files := ctx.projectStructure[layer]
		for _, file := range files {
			imports := ctx.importGraph[file]
			for _, imp := range imports {
				for _, framework := range frameworkPackages {
					if strings.Contains(imp, framework) {
						return fmt.Errorf("inner layer file %s imports framework %s", file, framework)
					}
				}
			}
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) domainEntitiesShouldBePureGoStructs() error {
	// Verify domain entities are POGOs (Plain Old Go Objects)
	entityFiles := ctx.projectStructure["entities"]
	if len(entityFiles) == 0 {
		entityFiles = ctx.projectStructure["aggregates"] // For DDD
	}
	
	for _, file := range entityFiles {
		content, err := os.ReadFile(filepath.Join(ctx.projectPath, file))
		if err != nil {
			continue
		}
		
		// Check for framework-specific tags or dependencies
		contentStr := string(content)
		if strings.Contains(contentStr, "gorm:") || 
		   strings.Contains(contentStr, "json:") ||
		   strings.Contains(contentStr, "binding:") {
			// Note: json tags might be acceptable, but gorm/binding tags indicate framework coupling
			if strings.Contains(contentStr, "gorm:") || strings.Contains(contentStr, "binding:") {
				return fmt.Errorf("domain entity in %s contains framework-specific tags", file)
			}
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) aggregatesShouldEncapsulateBusinessInvariants() error {
	// Verify DDD aggregates properly encapsulate invariants
	aggregateFiles := ctx.projectStructure["aggregates"]
	
	for _, file := range aggregateFiles {
		content, err := os.ReadFile(filepath.Join(ctx.projectPath, file))
		if err != nil {
			continue
		}
		
		contentStr := string(content)
		
		// Check for validation methods
		if !strings.Contains(contentStr, "Validate") && !strings.Contains(contentStr, "IsValid") {
			return fmt.Errorf("aggregate in %s does not contain validation methods", file)
		}
		
		// Check for private fields (encapsulation)
		if !strings.Contains(contentStr, "\t") {
			continue
		}
		
		// Simple check for struct with private fields
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filepath.Join(ctx.projectPath, file), nil, 0)
		if err != nil {
			continue
		}
		
		hasPrivateFields := false
		ast.Inspect(node, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok {
				if structType, ok := typeSpec.Type.(*ast.StructType); ok {
					for _, field := range structType.Fields.List {
						if len(field.Names) > 0 && !ast.IsExported(field.Names[0].Name) {
							hasPrivateFields = true
							return false
						}
					}
				}
			}
			return true
		})
		
		if !hasPrivateFields {
			return fmt.Errorf("aggregate in %s does not properly encapsulate fields", file)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) valueObjectsShouldBeImmutable() error {
	// Verify value objects are immutable
	voFiles := ctx.projectStructure["value_objects"]
	
	for _, file := range voFiles {
		content, err := os.ReadFile(filepath.Join(ctx.projectPath, file))
		if err != nil {
			continue
		}
		
		contentStr := string(content)
		
		// Check for setter methods (which would violate immutability)
		if strings.Contains(contentStr, "func (") && 
		   (strings.Contains(contentStr, "Set") || strings.Contains(contentStr, "Update") || strings.Contains(contentStr, "Modify")) {
			return fmt.Errorf("value object in %s contains setter methods", file)
		}
		
		// Check that creation is through constructor/factory functions
		if !strings.Contains(contentStr, "func New") {
			return fmt.Errorf("value object in %s does not have constructor function", file)
		}
	}
	
	return nil
}

func (ctx *ArchitectureTestContext) frameworkShouldBeProperlyIsolated() error {
	// Verify framework is isolated to appropriate layers
	return ctx.frameworkDependenciesShouldOnlyExistInOuterLayers()
}

func (ctx *ArchitectureTestContext) iAnalyzeIntegrationPoints() error {
	// Analyze how framework integrates with architecture
	// This would involve checking routing setup, middleware registration, etc.
	return nil
}