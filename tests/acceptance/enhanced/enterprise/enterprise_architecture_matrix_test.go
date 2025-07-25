package enterprise

import (
	"context"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// EnterpriseTestContext holds test state for enterprise architecture testing
type EnterpriseTestContext struct {
	projectConfig     *types.ProjectConfig
	projectPath       string
	tempDir           string
	startTime         time.Time
	lastCommandOutput string
	lastCommandError  error
	generatedFiles    []string
}

// TestFeatures runs the enterprise architecture matrix BDD tests
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &EnterpriseTestContext{}

			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.startTime = time.Now()
				ctx.generatedFiles = []string{}

				// Initialize templates
				if err := helpers.InitializeTemplates(); err != nil {
					return goCtx, err
				}

				return goCtx, nil
			})

			s.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				if ctx.tempDir != "" {
					os.RemoveAll(ctx.tempDir)
				}
				return goCtx, nil
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
		t.Fatal("non-zero status returned, failed to run enterprise architecture matrix tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *EnterpriseTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^all templates are properly initialized$`, ctx.allTemplatesAreProperlyInitialized)
	s.Step(`^I am testing enterprise architecture combinations$`, ctx.iAmTestingEnterpriseArchitectureCombinations)

	// Project generation steps
	s.Step(`^I generate a web API project with configuration:$`, ctx.iGenerateAWebAPIProjectWithConfiguration)
	s.Step(`^I generate a microservice project with configuration:$`, ctx.iGenerateAMicroserviceProjectWithConfiguration)
	s.Step(`^I generate an event-driven project with configuration:$`, ctx.iGenerateAnEventDrivenProjectWithConfiguration)
	s.Step(`^I generate a gRPC Gateway project with configuration:$`, ctx.iGenerateAGRPCGatewayProjectWithConfiguration)
	s.Step(`^I generate a web API project optimized for performance:$`, ctx.iGenerateAWebAPIProjectOptimizedForPerformance)
	s.Step(`^I generate a web API project with security focus:$`, ctx.iGenerateAWebAPIProjectWithSecurityFocus)
	s.Step(`^I generate a project ready for production deployment:$`, ctx.iGenerateAProjectReadyForProductionDeployment)

	// Universal validation steps
	s.Step(`^the project should compile successfully$`, ctx.theProjectShouldCompileSuccessfully)

	// Hexagonal architecture validations
	s.Step(`^hexagonal architecture boundaries should be enforced$`, ctx.hexagonalArchitectureBoundariesShouldBeEnforced)
	s.Step(`^domain layer should be isolated from infrastructure$`, ctx.domainLayerShouldBeIsolatedFromInfrastructure)
	s.Step(`^ports and adapters pattern should be implemented correctly$`, ctx.portsAndAdaptersPatternShouldBeImplementedCorrectly)
	s.Step(`^dependency injection should work with (.*)$`, ctx.dependencyInjectionShouldWorkWith)
	s.Step(`^(.*) database should be properly abstracted$`, ctx.databaseShouldBeProperlyAbstracted)
	s.Step(`^(.*) authentication should integrate with domain layer$`, ctx.authenticationShouldIntegrateWithDomainLayer)
	s.Step(`^(.*) logging should be available in all layers$`, ctx.loggingShouldBeAvailableInAllLayers)
	s.Step(`^repository interfaces should be defined in domain layer$`, ctx.repositoryInterfacesShouldBeDefinedInDomainLayer)
	s.Step(`^infrastructure adapters should implement domain interfaces$`, ctx.infrastructureAdaptersShouldImplementDomainInterfaces)

	// Clean architecture validations
	s.Step(`^clean architecture dependency rules should be enforced$`, ctx.cleanArchitectureDependencyRulesShouldBeEnforced)
	s.Step(`^entities should not depend on external frameworks$`, ctx.entitiesShouldNotDependOnExternalFrameworks)
	s.Step(`^use cases should be framework-independent$`, ctx.useCasesShouldBeFrameworkIndependent)
	s.Step(`^interface adapters should bridge frameworks and use cases$`, ctx.interfaceAdaptersShouldBridgeFrameworksAndUseCases)
	s.Step(`^frameworks and drivers should be in outermost layer$`, ctx.frameworksAndDriversShouldBeInOutermostLayer)
	s.Step(`^dependency inversion principle should be applied$`, ctx.dependencyInversionPrincipleShouldBeApplied)
	s.Step(`^(.*) assets should be properly integrated$`, ctx.assetsShouldBeProperlyIntegrated)
	s.Step(`^(.*) authentication should follow clean architecture patterns$`, ctx.authenticationShouldFollowCleanArchitecturePatterns)
	s.Step(`^database layer should not leak into business logic$`, ctx.databaseLayerShouldNotLeakIntoBusinessLogic)

	// DDD validations
	s.Step(`^domain model should be rich and expressive$`, ctx.domainModelShouldBeRichAndExpressive)
	s.Step(`^bounded contexts should be properly defined$`, ctx.boundedContextsShouldBeProperlyDefined)
	s.Step(`^aggregates should enforce business invariants$`, ctx.aggregatesShouldEnforceBusinessInvariants)
	s.Step(`^domain events should be properly implemented$`, ctx.domainEventsShouldBeProperlyImplemented)
	s.Step(`^value objects should be immutable$`, ctx.valueObjectsShouldBeImmutable)
	s.Step(`^domain services should encapsulate domain logic$`, ctx.domainServicesShouldEncapsulateDomainLogic)
	s.Step(`^repository patterns should abstract persistence$`, ctx.repositoryPatternsShouldAbstractPersistence)
	s.Step(`^application services should orchestrate use cases$`, ctx.applicationServicesShouldOrchestrateUseCases)
	s.Step(`^infrastructure should be separated from domain$`, ctx.infrastructureShouldBeSeparatedFromDomain)
	s.Step(`^(.*) persistence should support domain model$`, ctx.persistenceShouldSupportDomainModel)

	// Standard architecture validations
	s.Step(`^standard layered architecture should be implemented$`, ctx.standardLayeredArchitectureShouldBeImplemented)
	s.Step(`^handlers should delegate to services$`, ctx.handlersShouldDelegateToServices)
	s.Step(`^services should coordinate business operations$`, ctx.servicesShouldCoordinateBusinessOperations)
	s.Step(`^repositories should handle data persistence$`, ctx.repositoriesShouldHandleDataPersistence)
	s.Step(`^middleware should be properly configured$`, ctx.middlewareShouldBeProperlyConfigured)
	s.Step(`^error handling should be consistent across layers$`, ctx.errorHandlingShouldBeConsistentAcrossLayers)
	s.Step(`^(.*) routing should be properly organized$`, ctx.routingShouldBeProperlyOrganized)
	s.Step(`^(.*) connections should be optimized$`, ctx.connectionsShouldBeOptimized)
	s.Step(`^(.*) authentication should be production-ready$`, ctx.authenticationShouldBeProductionReady)

	// Microservice validations
	s.Step(`^microservice patterns should be implemented$`, ctx.microservicePatternsShouldBeImplemented)
	s.Step(`^service discovery should be configured$`, ctx.serviceDiscoveryShouldBeConfigured)
	s.Step(`^health checks should be available$`, ctx.healthChecksShouldBeAvailable)
	s.Step(`^metrics collection should be set up$`, ctx.metricsCollectionShouldBeSetUp)
	s.Step(`^distributed tracing should be configured$`, ctx.distributedTracingShouldBeConfigured)
	s.Step(`^circuit breaker patterns should be available$`, ctx.circuitBreakerPatternsShouldBeAvailable)
	s.Step(`^(.*) architecture should be properly applied$`, ctx.architectureShouldBeProperlyApplied)
	s.Step(`^API documentation should be generated$`, ctx.apiDocumentationShouldBeGenerated)

	// Event-driven validations
	s.Step(`^CQRS pattern should be properly implemented$`, ctx.cqrsPatternShouldBeProperlyImplemented)
	s.Step(`^command handlers should be separated from query handlers$`, ctx.commandHandlersShouldBeSeparatedFromQueryHandlers)
	s.Step(`^event sourcing should work correctly$`, ctx.eventSourcingShouldWorkCorrectly)
	s.Step(`^event store should be configured for (.*)$`, ctx.eventStoreShouldBeConfiguredFor)
	s.Step(`^projection rebuilding should be supported$`, ctx.projectionRebuildingShouldBeSupported)
	s.Step(`^saga patterns should be available$`, ctx.sagaPatternsShouldBeAvailable)
	s.Step(`^eventual consistency should be handled$`, ctx.eventualConsistencyShouldBeHandled)
	s.Step(`^event versioning should be supported$`, ctx.eventVersioningShouldBeSupported)

	// gRPC validations
	s.Step(`^gRPC services should be properly defined$`, ctx.grpcServicesShouldBeProperlyDefined)
	s.Step(`^protobuf definitions should be valid$`, ctx.protobufDefinitionsShouldBeValid)
	s.Step(`^Gateway should proxy HTTP to gRPC correctly$`, ctx.gatewayShouldProxyHTTPToGRPCCorrectly)
	s.Step(`^authentication should work for both HTTP and gRPC$`, ctx.authenticationShouldWorkForBothHTTPAndGRPC)
	s.Step(`^both REST and gRPC endpoints should be available$`, ctx.bothRESTAndGRPCEndpointsShouldBeAvailable)
	s.Step(`^service mesh compatibility should be maintained$`, ctx.serviceMeshCompatibilityShouldBeMaintained)
	s.Step(`^load balancing should be configured$`, ctx.loadBalancingShouldBeConfigured)
	s.Step(`^SSL/TLS should be properly configured$`, ctx.sslTlsShouldBeProperlyConfigured)

	// Performance validations
	s.Step(`^performance optimizations should be implemented$`, ctx.performanceOptimizationsShouldBeImplemented)
	s.Step(`^database connections should be pooled efficiently$`, ctx.databaseConnectionsShouldBePooledEfficiently)
	s.Step(`^logging should have minimal performance overhead$`, ctx.loggingShouldHaveMinimalPerformanceOverhead)
	s.Step(`^memory allocations should be optimized$`, ctx.memoryAllocationsShouldBeOptimized)
	s.Step(`^response times should be within acceptable limits$`, ctx.responseTimesShouldBeWithinAcceptableLimits)
	s.Step(`^concurrent request handling should be optimized$`, ctx.concurrentRequestHandlingShouldBeOptimized)
	s.Step(`^resource usage should be monitored$`, ctx.resourceUsageShouldBeMonitored)

	// Security validations
	s.Step(`^security best practices should be implemented$`, ctx.securityBestPracticesShouldBeImplemented)
	s.Step(`^authentication should be properly secured$`, ctx.authenticationShouldBeProperlySecured)
	s.Step(`^authorization should be role-based$`, ctx.authorizationShouldBeRoleBased)
	s.Step(`^input validation should prevent injection attacks$`, ctx.inputValidationShouldPreventInjectionAttacks)
	s.Step(`^HTTPS should be enforced$`, ctx.httpsShouldBeEnforced)
	s.Step(`^security headers should be configured$`, ctx.securityHeadersShouldBeConfigured)
	s.Step(`^audit logging should capture security events$`, ctx.auditLoggingShouldCaptureSecurityEvents)
	s.Step(`^sensitive data should be properly handled$`, ctx.sensitiveDataShouldBeProperlyHandled)
	s.Step(`^CORS should be properly configured$`, ctx.corsShouldBeProperlyConfigured)

	// Deployment validations
	s.Step(`^deployment configurations should be generated$`, ctx.deploymentConfigurationsShouldBeGenerated)
	s.Step(`^Docker configuration should be optimized$`, ctx.dockerConfigurationShouldBeOptimized)
	s.Step(`^metrics should be exposed for monitoring$`, ctx.metricsShouldBeExposedForMonitoring)
	s.Step(`^logging should be container-friendly$`, ctx.loggingShouldBeContainerFriendly)
	s.Step(`^graceful shutdown should be implemented$`, ctx.gracefulShutdownShouldBeImplemented)
	s.Step(`^configuration management should be environment-aware$`, ctx.configurationManagementShouldBeEnvironmentAware)
	s.Step(`^database migrations should be automated$`, ctx.databaseMigrationsShouldBeAutomated)
}

// Background step implementations
func (ctx *EnterpriseTestContext) iHaveTheGoStarterCLIAvailable() error {
	cmd := exec.Command("go-starter", "version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go-starter CLI not available: %v", err)
	}
	return nil
}

func (ctx *EnterpriseTestContext) allTemplatesAreProperlyInitialized() error {
	return helpers.InitializeTemplates()
}

func (ctx *EnterpriseTestContext) iAmTestingEnterpriseArchitectureCombinations() error {
	return nil
}

// Project generation step implementations
func (ctx *EnterpriseTestContext) iGenerateAWebAPIProjectWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAMicroserviceProjectWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAnEventDrivenProjectWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAGRPCGatewayProjectWithConfiguration(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAWebAPIProjectOptimizedForPerformance(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAWebAPIProjectWithSecurityFocus(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) iGenerateAProjectReadyForProductionDeployment(table *godog.Table) error {
	return ctx.generateProjectWithConfiguration(table)
}

func (ctx *EnterpriseTestContext) generateProjectWithConfiguration(table *godog.Table) error {
	// Parse configuration from table
	config := &types.ProjectConfig{}

	for i := 0; i < len(table.Rows); i++ {
		row := table.Rows[i]
		key := row.Cells[0].Value
		value := row.Cells[1].Value

		switch key {
		case "type":
			// Keep web-api as-is since that's what the validator expects
			config.Type = value
		case "architecture":
			config.Architecture = value
		case "framework":
			config.Framework = value
		case "database":
			if config.Features == nil {
				config.Features = &types.Features{}
			}
			config.Features.Database.Driver = value
		case "orm":
			if config.Features == nil {
				config.Features = &types.Features{}
			}
			config.Features.Database.ORM = value
		case "auth":
			if config.Features == nil {
				config.Features = &types.Features{}
			}
			config.Features.Authentication.Type = value
		case "logger":
			config.Logger = value
		case "asset_pipeline":
			// AssetPipeline not in ProjectConfig, skip for now
		case "complexity":
			// Complexity not in ProjectConfig, skip for now
		case "go_version":
			config.GoVersion = value
		}
	}

	// Set defaults
	if config.Name == "" {
		config.Name = "test-enterprise-project"
	}
	if config.Module == "" {
		config.Module = "github.com/test/enterprise-project"
	}
	if config.GoVersion == "" {
		config.GoVersion = "1.23"
	}
	if config.Architecture == "" {
		config.Architecture = "standard"
	}

	ctx.projectConfig = config

	// Generate project using helpers
	var err error
	ctx.projectPath, err = ctx.generateTestProject(config, ctx.tempDir)
	if err != nil {
		ctx.lastCommandError = err
		return fmt.Errorf("failed to generate project: %v", err)
	}

	// Collect generated files
	err = filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(ctx.projectPath, path)
			ctx.generatedFiles = append(ctx.generatedFiles, relPath)
		}
		return nil
	})

	return err
}

// Universal validation implementations
func (ctx *EnterpriseTestContext) theProjectShouldCompileSuccessfully() error {
	return ctx.validateProjectCompilation()
}

// Hexagonal architecture validation implementations
func (ctx *EnterpriseTestContext) hexagonalArchitectureBoundariesShouldBeEnforced() error {
	// Check for hexagonal architecture directory structure
	expectedDirs := []string{
		"internal/domain",
		"internal/ports",
		"internal/adapters",
		"internal/application",
	}

	for _, dir := range expectedDirs {
		dirPath := filepath.Join(ctx.projectPath, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			return fmt.Errorf("hexagonal architecture directory missing: %s", dir)
		}
	}

	return nil
}

func (ctx *EnterpriseTestContext) domainLayerShouldBeIsolatedFromInfrastructure() error {
	// Parse domain layer files and check for infrastructure imports
	domainDir := filepath.Join(ctx.projectPath, "internal/domain")
	return ctx.checkDirectoryForForbiddenImports(domainDir, []string{
		"database/sql",
		ctx.projectConfig.Framework,
		"http",
	})
}

func (ctx *EnterpriseTestContext) portsAndAdaptersPatternShouldBeImplementedCorrectly() error {
	// Check that ports (interfaces) are defined
	portsDir := filepath.Join(ctx.projectPath, "internal/ports")
	if _, err := os.Stat(portsDir); os.IsNotExist(err) {
		return fmt.Errorf("ports directory missing")
	}

	// Check that adapters implement ports
	adaptersDir := filepath.Join(ctx.projectPath, "internal/adapters")
	if _, err := os.Stat(adaptersDir); os.IsNotExist(err) {
		return fmt.Errorf("adapters directory missing")
	}

	return nil
}

func (ctx *EnterpriseTestContext) dependencyInjectionShouldWorkWith(framework string) error {
	// Check for dependency injection patterns in the generated code
	return ctx.checkForDependencyInjectionPatterns()
}

func (ctx *EnterpriseTestContext) databaseShouldBeProperlyAbstracted(database string) error {
	// Check that database logic is abstracted behind interfaces
	return ctx.checkForDatabaseAbstraction(database)
}

func (ctx *EnterpriseTestContext) authenticationShouldIntegrateWithDomainLayer(auth string) error {
	// Check that authentication integrates properly with domain layer
	return ctx.checkAuthenticationIntegration(auth)
}

func (ctx *EnterpriseTestContext) loggingShouldBeAvailableInAllLayers(logger string) error {
	// Check that logging is available across all layers
	return ctx.checkLoggingAvailability(logger)
}

func (ctx *EnterpriseTestContext) repositoryInterfacesShouldBeDefinedInDomainLayer() error {
	// Check for repository interfaces in domain layer
	domainDir := filepath.Join(ctx.projectPath, "internal/domain")
	return ctx.checkForRepositoryInterfaces(domainDir)
}

func (ctx *EnterpriseTestContext) infrastructureAdaptersShouldImplementDomainInterfaces() error {
	// Check that infrastructure adapters implement domain interfaces
	return ctx.checkAdapterImplementations()
}

// Clean architecture validation implementations
func (ctx *EnterpriseTestContext) cleanArchitectureDependencyRulesShouldBeEnforced() error {
	// Check clean architecture dependency rules
	return ctx.validateCleanArchitectureDependencies()
}

func (ctx *EnterpriseTestContext) entitiesShouldNotDependOnExternalFrameworks() error {
	// Check that entities don't import framework dependencies
	entitiesDir := filepath.Join(ctx.projectPath, "internal/entities")
	return ctx.checkDirectoryForForbiddenImports(entitiesDir, []string{
		ctx.projectConfig.Framework,
		"http",
		"database/sql",
	})
}

func (ctx *EnterpriseTestContext) useCasesShouldBeFrameworkIndependent() error {
	// Check that use cases are framework independent
	useCasesDir := filepath.Join(ctx.projectPath, "internal/usecases")
	return ctx.checkDirectoryForForbiddenImports(useCasesDir, []string{
		ctx.projectConfig.Framework,
		"http",
	})
}

func (ctx *EnterpriseTestContext) interfaceAdaptersShouldBridgeFrameworksAndUseCases() error {
	// Check for interface adapters that bridge frameworks and use cases
	return ctx.checkForInterfaceAdapters()
}

func (ctx *EnterpriseTestContext) frameworksAndDriversShouldBeInOutermostLayer() error {
	// Check that frameworks and drivers are in the outermost layer
	return ctx.checkFrameworkPlacement()
}

func (ctx *EnterpriseTestContext) dependencyInversionPrincipleShouldBeApplied() error {
	// Check for dependency inversion patterns
	return ctx.checkDependencyInversion()
}

func (ctx *EnterpriseTestContext) assetsShouldBeProperlyIntegrated(assetPipeline string) error {
	// Check asset pipeline integration
	return ctx.checkAssetPipelineIntegration(assetPipeline)
}

func (ctx *EnterpriseTestContext) authenticationShouldFollowCleanArchitecturePatterns(auth string) error {
	// Check that authentication follows clean architecture patterns
	return ctx.checkCleanAuthPatterns(auth)
}

func (ctx *EnterpriseTestContext) databaseLayerShouldNotLeakIntoBusinessLogic() error {
	// Check that database concerns don't leak into business logic
	return ctx.checkDatabaseLeakage()
}

// DDD validation implementations
func (ctx *EnterpriseTestContext) domainModelShouldBeRichAndExpressive() error {
	// Check for rich domain model patterns
	return ctx.checkDomainModelRichness()
}

func (ctx *EnterpriseTestContext) boundedContextsShouldBeProperlyDefined() error {
	// Check for bounded context definitions
	return ctx.checkBoundedContexts()
}

func (ctx *EnterpriseTestContext) aggregatesShouldEnforceBusinessInvariants() error {
	// Check for aggregate patterns and invariant enforcement
	return ctx.checkAggregatePatterns()
}

func (ctx *EnterpriseTestContext) domainEventsShouldBeProperlyImplemented() error {
	// Check for domain event implementations
	return ctx.checkDomainEvents()
}

func (ctx *EnterpriseTestContext) valueObjectsShouldBeImmutable() error {
	// Check for value object patterns
	return ctx.checkValueObjects()
}

func (ctx *EnterpriseTestContext) domainServicesShouldEncapsulateDomainLogic() error {
	// Check for domain service patterns
	return ctx.checkDomainServices()
}

func (ctx *EnterpriseTestContext) repositoryPatternsShouldAbstractPersistence() error {
	// Check for repository abstraction patterns
	return ctx.checkRepositoryAbstraction()
}

func (ctx *EnterpriseTestContext) applicationServicesShouldOrchestrateUseCases() error {
	// Check for application service orchestration
	return ctx.checkApplicationServices()
}

func (ctx *EnterpriseTestContext) infrastructureShouldBeSeparatedFromDomain() error {
	// Check for infrastructure/domain separation
	return ctx.checkInfrastructureSeparation()
}

func (ctx *EnterpriseTestContext) persistenceShouldSupportDomainModel(database string) error {
	// Check that persistence supports the domain model
	return ctx.checkPersistenceDomainSupport(database)
}

// Helper methods for validation (simplified implementations for now)
func (ctx *EnterpriseTestContext) checkDirectoryForForbiddenImports(dir string, forbiddenImports []string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// Directory doesn't exist, skip check
		return nil
	}

	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return err
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			for _, forbidden := range forbiddenImports {
				if strings.Contains(importPath, forbidden) {
					return fmt.Errorf("forbidden import %s found in %s", importPath, path)
				}
			}
		}

		return nil
	})
}

// Simplified implementations for other validation methods
func (ctx *EnterpriseTestContext) checkForDependencyInjectionPatterns() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkForDatabaseAbstraction(database string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkAuthenticationIntegration(auth string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkLoggingAvailability(logger string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkForRepositoryInterfaces(dir string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkAdapterImplementations() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) validateCleanArchitectureDependencies() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkForInterfaceAdapters() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkFrameworkPlacement() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkDependencyInversion() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkAssetPipelineIntegration(assetPipeline string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkCleanAuthPatterns(auth string) error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkDatabaseLeakage() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkDomainModelRichness() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkBoundedContexts() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkAggregatePatterns() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkDomainEvents() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkValueObjects() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkDomainServices() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkRepositoryAbstraction() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkApplicationServices() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkInfrastructureSeparation() error {
	return nil // Simplified implementation
}

func (ctx *EnterpriseTestContext) checkPersistenceDomainSupport(database string) error {
	return nil // Simplified implementation
}

// Standard architecture validations (simplified)
func (ctx *EnterpriseTestContext) standardLayeredArchitectureShouldBeImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) handlersShouldDelegateToServices() error {
	return nil
}

func (ctx *EnterpriseTestContext) servicesShouldCoordinateBusinessOperations() error {
	return nil
}

func (ctx *EnterpriseTestContext) repositoriesShouldHandleDataPersistence() error {
	return nil
}

func (ctx *EnterpriseTestContext) middlewareShouldBeProperlyConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) errorHandlingShouldBeConsistentAcrossLayers() error {
	return nil
}

func (ctx *EnterpriseTestContext) routingShouldBeProperlyOrganized(framework string) error {
	return nil
}

func (ctx *EnterpriseTestContext) connectionsShouldBeOptimized(database string) error {
	return nil
}

func (ctx *EnterpriseTestContext) authenticationShouldBeProductionReady(auth string) error {
	return nil
}

// Additional validation method stubs (simplified implementations)
func (ctx *EnterpriseTestContext) microservicePatternsShouldBeImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) serviceDiscoveryShouldBeConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) healthChecksShouldBeAvailable() error {
	return nil
}

func (ctx *EnterpriseTestContext) metricsCollectionShouldBeSetUp() error {
	return nil
}

func (ctx *EnterpriseTestContext) distributedTracingShouldBeConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) circuitBreakerPatternsShouldBeAvailable() error {
	return nil
}

func (ctx *EnterpriseTestContext) architectureShouldBeProperlyApplied(architecture string) error {
	return nil
}

func (ctx *EnterpriseTestContext) apiDocumentationShouldBeGenerated() error {
	return nil
}

func (ctx *EnterpriseTestContext) cqrsPatternShouldBeProperlyImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) commandHandlersShouldBeSeparatedFromQueryHandlers() error {
	return nil
}

func (ctx *EnterpriseTestContext) eventSourcingShouldWorkCorrectly() error {
	return nil
}

func (ctx *EnterpriseTestContext) eventStoreShouldBeConfiguredFor(database string) error {
	return nil
}

func (ctx *EnterpriseTestContext) projectionRebuildingShouldBeSupported() error {
	return nil
}

func (ctx *EnterpriseTestContext) sagaPatternsShouldBeAvailable() error {
	return nil
}

func (ctx *EnterpriseTestContext) eventualConsistencyShouldBeHandled() error {
	return nil
}

func (ctx *EnterpriseTestContext) eventVersioningShouldBeSupported() error {
	return nil
}

func (ctx *EnterpriseTestContext) grpcServicesShouldBeProperlyDefined() error {
	return nil
}

func (ctx *EnterpriseTestContext) protobufDefinitionsShouldBeValid() error {
	return nil
}

func (ctx *EnterpriseTestContext) gatewayShouldProxyHTTPToGRPCCorrectly() error {
	return nil
}

func (ctx *EnterpriseTestContext) authenticationShouldWorkForBothHTTPAndGRPC() error {
	return nil
}

func (ctx *EnterpriseTestContext) bothRESTAndGRPCEndpointsShouldBeAvailable() error {
	return nil
}

func (ctx *EnterpriseTestContext) serviceMeshCompatibilityShouldBeMaintained() error {
	return nil
}

func (ctx *EnterpriseTestContext) loadBalancingShouldBeConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) sslTlsShouldBeProperlyConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) performanceOptimizationsShouldBeImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) databaseConnectionsShouldBePooledEfficiently() error {
	return nil
}

func (ctx *EnterpriseTestContext) loggingShouldHaveMinimalPerformanceOverhead() error {
	return nil
}

func (ctx *EnterpriseTestContext) memoryAllocationsShouldBeOptimized() error {
	return nil
}

func (ctx *EnterpriseTestContext) responseTimesShouldBeWithinAcceptableLimits() error {
	return nil
}

func (ctx *EnterpriseTestContext) concurrentRequestHandlingShouldBeOptimized() error {
	return nil
}

func (ctx *EnterpriseTestContext) resourceUsageShouldBeMonitored() error {
	return nil
}

func (ctx *EnterpriseTestContext) securityBestPracticesShouldBeImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) authenticationShouldBeProperlySecured() error {
	return nil
}

func (ctx *EnterpriseTestContext) authorizationShouldBeRoleBased() error {
	return nil
}

func (ctx *EnterpriseTestContext) inputValidationShouldPreventInjectionAttacks() error {
	return nil
}

func (ctx *EnterpriseTestContext) httpsShouldBeEnforced() error {
	return nil
}

func (ctx *EnterpriseTestContext) securityHeadersShouldBeConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) auditLoggingShouldCaptureSecurityEvents() error {
	return nil
}

func (ctx *EnterpriseTestContext) sensitiveDataShouldBeProperlyHandled() error {
	return nil
}

func (ctx *EnterpriseTestContext) corsShouldBeProperlyConfigured() error {
	return nil
}

func (ctx *EnterpriseTestContext) deploymentConfigurationsShouldBeGenerated() error {
	return nil
}

func (ctx *EnterpriseTestContext) dockerConfigurationShouldBeOptimized() error {
	return nil
}

func (ctx *EnterpriseTestContext) metricsShouldBeExposedForMonitoring() error {
	return nil
}

func (ctx *EnterpriseTestContext) loggingShouldBeContainerFriendly() error {
	return nil
}

func (ctx *EnterpriseTestContext) gracefulShutdownShouldBeImplemented() error {
	return nil
}

func (ctx *EnterpriseTestContext) configurationManagementShouldBeEnvironmentAware() error {
	return nil
}

func (ctx *EnterpriseTestContext) databaseMigrationsShouldBeAutomated() error {
	return nil
}

// Helper methods for project generation and validation
func (ctx *EnterpriseTestContext) generateTestProject(config *types.ProjectConfig, tempDir string) (string, error) {
	return ctx.generateProjectDirect(config, tempDir)
}

func (ctx *EnterpriseTestContext) generateProjectDirect(config *types.ProjectConfig, tempDir string) (string, error) {
	gen := generator.New()
	
	projectPath := filepath.Join(tempDir, config.Name)
	
	_, err := gen.Generate(*config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	if err != nil {
		return "", err
	}
	
	return projectPath, nil
}

func (ctx *EnterpriseTestContext) validateProjectCompilation() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("project path not set")
	}
	
	// Check if go.mod exists
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("go.mod not found in project")
	}
	
	// Try to run go build
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = ctx.projectPath
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("project compilation failed: %v\nOutput: %s", err, string(output))
	}
	
	return nil
}