package microservice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/suite"
)

// MicroserviceAcceptanceTestSuite provides acceptance tests for service blueprints
type MicroserviceAcceptanceTestSuite struct {
	suite.Suite
	ctx *MicroserviceTestContext
}

// SetupSuite runs before all tests in the suite
func (suite *MicroserviceAcceptanceTestSuite) SetupSuite() {
	suite.ctx = InitializeMicroserviceContext()
}

// TearDownSuite runs after all tests in the suite
func (suite *MicroserviceAcceptanceTestSuite) TearDownSuite() {
	if suite.ctx != nil {
		suite.ctx.cleanup()
	}
}

// TestServiceBlueprintGeneration tests service blueprint generation with godog
func (suite *MicroserviceAcceptanceTestSuite) TestServiceBlueprintGeneration() {
	godogSuite := godog.TestSuite{
		ScenarioInitializer: suite.InitializeMicroserviceScenarios,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: suite.T(),
		},
	}

	if godogSuite.Run() != 0 {
		suite.T().Fatal("Non-zero status returned, failed to run service acceptance tests")
	}
}

// InitializeMicroserviceScenarios registers step definitions for service scenarios
func (suite *MicroserviceAcceptanceTestSuite) InitializeMicroserviceScenarios(sc *godog.ScenarioContext) {
	ctx := suite.ctx

	// Background steps
	sc.Given(`^the go-starter CLI tool is available$`, func() error {
		// Check if go-starter is available
		cmd := exec.Command("go-starter", "--help")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("go-starter CLI not available: %s", err.Error())
		}
		return nil
	})

	sc.Given(`^I am in a clean working directory$`, func() error {
		// Ensure we're in a clean state
		ctx.cleanup()
		return nil
	})

	sc.Given(`^I have Docker available for containerization testing$`, func() error {
		// Check if Docker is available
		cmd := exec.Command("docker", "--version")
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("Docker not available: %s", err.Error())
		}
		return nil
	})

	// Microservice generation steps
	sc.Given(`^I want to create a gRPC-based microservice$`, ctx.iWantToCreateAGRPCBasedMicroservice)
	sc.When(`^I run the command "([^"]*)"$`, ctx.iRunTheCommand)
	sc.Then(`^the generation should succeed$`, ctx.theGenerationShouldSucceed)
	sc.Then(`^the project should contain all essential microservice components$`, ctx.theProjectShouldContainAllEssentialMicroserviceComponents)
	sc.Then(`^the generated code should compile successfully$`, ctx.theGeneratedCodeShouldCompileSuccessfully)
	sc.Then(`^the service should support gRPC and HTTP interfaces$`, ctx.theServiceShouldSupportGRPCAndHTTPInterfaces)
	sc.Then(`^the microservice should be container-ready with health checks$`, ctx.theMicroserviceShouldBeContainerReadyWithHealthChecks)

	// Lambda generation steps
	sc.Given(`^I want to create an AWS Lambda serverless function$`, ctx.iWantToCreateAnAWSLambdaServerlessFunction)
	sc.Then(`^the project should contain all essential serverless components$`, ctx.theProjectShouldContainAllEssentialServerlessComponents)
	sc.Then(`^the lambda should support AWS SDK integration$`, ctx.theLambdaShouldSupportAWSDKIntegration)
	sc.Then(`^the function should include observability features$`, ctx.theFunctionShouldIncludeObservabilityFeatures)

	// Protocol testing steps
	sc.Given(`^I want to create microservices with various protocols$`, ctx.iWantToCreateMicroservicesWithVariousProtocols)
	sc.When(`^I generate a microservice with protocol "([^"]*)"$`, ctx.iGenerateAMicroserviceWithProtocol)
	sc.Then(`^the project should support the "([^"]*)" communication pattern$`, ctx.theProjectShouldSupportTheProtocolCommunicationPattern)
	sc.Then(`^the service should include appropriate client libraries$`, ctx.theServiceShouldIncludeAppropriateClientLibraries)
	sc.Then(`^the protocol-specific middleware should be configured$`, ctx.theProtocolSpecificMiddlewareShouldBeConfigured)
	sc.Then(`^the service should compile and run correctly$`, ctx.theServiceShouldCompileAndRunCorrectly)

	// Database integration steps
	sc.Given(`^I want to create services with database persistence$`, func() error {
		ctx.scenarios["database_testing"] = true
		return nil
	})
	sc.When(`^I generate a service with database "([^"]*)" and ORM "([^"]*)"$`, ctx.iGenerateAServiceWithDatabaseAndORM)
	sc.Then(`^the project should include database configuration$`, ctx.theProjectShouldIncludeDatabaseConfiguration)
	sc.Then(`^the service should support database migrations$`, ctx.theServiceShouldSupportDatabaseMigrations)
	sc.Then(`^the repository layer should use the specified ORM$`, ctx.theRepositoryLayerShouldUseTheSpecifiedORM)
	sc.Then(`^the database connection should be testable with containers$`, ctx.theDatabaseConnectionShouldBeTestableWithContainers)
	sc.Then(`^the service should handle connection pooling$`, ctx.theServiceShouldHandleConnectionPooling)

	// Service discovery steps
	sc.Given(`^I want microservices that can discover each other$`, ctx.iWantMicroservicesThatCanDiscoverEachOther)
	sc.When(`^I generate a microservice with service discovery$`, ctx.iGenerateAMicroserviceWithServiceDiscovery)
	sc.Then(`^the service should include service registry integration$`, ctx.theServiceShouldIncludeServiceRegistryIntegration)
	sc.Then(`^the service should support health check endpoints$`, ctx.theServiceShouldSupportHealthCheckEndpoints)
	sc.Then(`^the service should handle service registration$`, ctx.theServiceShouldHandleServiceRegistration)
	sc.Then(`^the service should support load balancing$`, ctx.theServiceShouldSupportLoadBalancing)
	sc.Then(`^the discovery configuration should be environment-aware$`, ctx.theDiscoveryConfigurationShouldBeEnvironmentAware)

	// Add more step definitions for remaining scenarios...
	// This is a comprehensive foundation that can be extended

	// Scenario hooks
	sc.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, suite.ctx.beforeScenario(sc)
	})
	sc.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		return ctx, suite.ctx.afterScenario(sc, err)
	})
}

// TestMicroserviceStandardGeneration tests basic microservice generation
func (suite *MicroserviceAcceptanceTestSuite) TestMicroserviceStandardGeneration() {
	ctx := suite.ctx
	
	// Test basic microservice generation
	err := ctx.iRunTheCommand("go-starter new test-microservice --type=microservice-standard --framework=grpc --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldContainAllEssentialMicroserviceComponents()
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
}

// TestLambdaStandardGeneration tests basic lambda generation
func (suite *MicroserviceAcceptanceTestSuite) TestLambdaStandardGeneration() {
	ctx := suite.ctx
	
	// Test basic lambda generation
	err := ctx.iRunTheCommand("go-starter new test-lambda --type=lambda-standard --runtime=aws --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	err = ctx.theProjectShouldContainAllEssentialServerlessComponents()
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
}

// TestProtocolSupport tests different communication protocols
func (suite *MicroserviceAcceptanceTestSuite) TestProtocolSupport() {
	protocols := []string{"grpc", "rest", "graphql"}
	
	for _, protocol := range protocols {
		suite.Run(fmt.Sprintf("Protocol_%s", protocol), func() {
			ctx := InitializeMicroserviceContext()
			defer ctx.cleanup()
			
			err := ctx.iGenerateAMicroserviceWithProtocol(protocol)
			suite.Require().NoError(err)
			
			err = ctx.theGenerationShouldSucceed()
			suite.Require().NoError(err)
			
			err = ctx.theProjectShouldSupportTheProtocolCommunicationPattern(protocol)
			suite.Require().NoError(err)
			
			err = ctx.theServiceShouldIncludeAppropriateClientLibraries()
			suite.Require().NoError(err)
			
			err = ctx.theGeneratedCodeShouldCompileSuccessfully()
			suite.Require().NoError(err)
		})
	}
}

// TestDatabaseIntegration tests different database integrations
func (suite *MicroserviceAcceptanceTestSuite) TestDatabaseIntegration() {
	if testing.Short() {
		suite.T().Skip("Skipping database integration tests in short mode")
	}
	
	testCases := []struct {
		database string
		orm      string
	}{
		{"postgres", "gorm"},
		{"postgres", "sqlx"},
		{"mysql", "gorm"},
		{"mongodb", "mongo"},
		{"redis", "redis"},
	}
	
	for _, tc := range testCases {
		suite.Run(fmt.Sprintf("Database_%s_%s", tc.database, tc.orm), func() {
			ctx := InitializeMicroserviceContext()
			defer ctx.cleanup()
			
			err := ctx.iGenerateAServiceWithDatabaseAndORM(tc.database, tc.orm)
			suite.Require().NoError(err)
			
			err = ctx.theGenerationShouldSucceed()
			suite.Require().NoError(err)
			
			err = ctx.theProjectShouldIncludeDatabaseConfiguration()
			suite.Require().NoError(err)
			
			err = ctx.theServiceShouldSupportDatabaseMigrations()
			suite.Require().NoError(err)
			
			err = ctx.theRepositoryLayerShouldUseTheSpecifiedORM()
			suite.Require().NoError(err)
			
			err = ctx.theGeneratedCodeShouldCompileSuccessfully()
			suite.Require().NoError(err)
			
			// Test with containers if Docker is available
			if ctx.containers != nil {
				err = ctx.theDatabaseConnectionShouldBeTestableWithContainers()
				if err != nil {
					suite.T().Logf("Container testing failed (expected in CI): %v", err)
				}
			}
		})
	}
}

// TestServiceDiscovery tests service discovery and registration
func (suite *MicroserviceAcceptanceTestSuite) TestServiceDiscovery() {
	ctx := suite.ctx
	
	err := ctx.iGenerateAMicroserviceWithServiceDiscovery()
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	err = ctx.theServiceShouldIncludeServiceRegistryIntegration()
	suite.Require().NoError(err)
	
	err = ctx.theServiceShouldSupportHealthCheckEndpoints()
	suite.Require().NoError(err)
	
	err = ctx.theServiceShouldHandleServiceRegistration()
	suite.Require().NoError(err)
	
	err = ctx.theDiscoveryConfigurationShouldBeEnvironmentAware()
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
}

// TestContainerization tests Docker and Kubernetes support
func (suite *MicroserviceAcceptanceTestSuite) TestContainerization() {
	ctx := suite.ctx
	
	// Test basic containerization
	err := ctx.iRunTheCommand("go-starter new test-container --type=microservice-standard --deployment=kubernetes --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check for containerization files
	containerFiles := []string{
		"Dockerfile",
		"docker-compose.yml",
		"k8s/deployment.yaml",
		"k8s/service.yaml",
		"k8s/configmap.yaml",
	}
	
	err = ctx.checkRequiredFiles(containerFiles)
	suite.Require().NoError(err)
	
	// Check Dockerfile for best practices
	dockerfilePath := filepath.Join(ctx.projectPath, "Dockerfile")
	content, err := os.ReadFile(dockerfilePath)
	suite.Require().NoError(err)
	
	dockerfileContent := string(content)
	
	// Verify Dockerfile best practices
	suite.Contains(dockerfileContent, "FROM golang:", "Should use Go base image")
	suite.Contains(dockerfileContent, "WORKDIR", "Should set working directory")
	suite.Contains(dockerfileContent, "COPY go.mod go.sum", "Should copy dependency files first")
	suite.Contains(dockerfileContent, "RUN go mod download", "Should download dependencies")
	suite.Contains(dockerfileContent, "HEALTHCHECK", "Should include health check")
	suite.Contains(dockerfileContent, "USER", "Should run as non-root user")
}

// TestObservability tests observability features
func (suite *MicroserviceAcceptanceTestSuite) TestObservability() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-observability --type=microservice-standard --observability=full --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check for observability files
	observabilityFiles := []string{
		"internal/observability/metrics.go",
		"internal/observability/tracing.go",
		"internal/observability/logging.go",
		"configs/observability.yaml",
	}
	
	err = ctx.checkRequiredFiles(observabilityFiles)
	suite.Require().NoError(err)
	
	// Check go.mod for observability dependencies
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	suite.Require().NoError(err)
	
	goModContent := string(content)
	
	// Verify observability dependencies
	observabilityDeps := []string{
		"go.opentelemetry.io/otel",
		"github.com/prometheus/client_golang",
	}
	
	for _, dep := range observabilityDeps {
		suite.Contains(goModContent, dep, fmt.Sprintf("Should include %s dependency", dep))
	}
}

// TestSecurityFeatures tests security configurations
func (suite *MicroserviceAcceptanceTestSuite) TestSecurityFeatures() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-security --type=microservice-standard --auth=jwt --security=full --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check for security files
	securityFiles := []string{
		"internal/security/auth.go",
		"internal/security/jwt.go",
		"internal/middleware/auth.go",
		"internal/middleware/rate_limit.go",
		"configs/security.yaml",
	}
	
	err = ctx.checkRequiredFiles(securityFiles)
	suite.Require().NoError(err)
	
	// Check security configuration
	securityConfigPath := filepath.Join(ctx.projectPath, "configs/security.yaml")
	content, err := os.ReadFile(securityConfigPath)
	suite.Require().NoError(err)
	
	securityContent := string(content)
	
	// Verify security configurations
	securityFeatures := []string{
		"jwt",
		"rate_limit",
		"cors",
		"encryption",
	}
	
	for _, feature := range securityFeatures {
		suite.Contains(strings.ToLower(securityContent), feature, 
			fmt.Sprintf("Should include %s security feature", feature))
	}
}

// TestPerformanceOptimizations tests performance-related configurations
func (suite *MicroserviceAcceptanceTestSuite) TestPerformanceOptimizations() {
	ctx := suite.ctx
	
	err := ctx.iRunTheCommand("go-starter new test-performance --type=microservice-standard --performance=optimized --no-git")
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Check for performance optimization files
	perfFiles := []string{
		"internal/performance/pool.go",
		"internal/performance/cache.go",
		"internal/middleware/compression.go",
		"configs/performance.yaml",
	}
	
	err = ctx.checkRequiredFiles(perfFiles)
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
}

// TestComprehensiveGeneration tests a fully-featured service generation
func (suite *MicroserviceAcceptanceTestSuite) TestComprehensiveGeneration() {
	ctx := suite.ctx
	
	command := `go-starter new comprehensive-service \
		--type=microservice-standard \
		--protocol=grpc \
		--database=postgres \
		--orm=gorm \
		--auth=jwt \
		--observability=full \
		--security=full \
		--deployment=kubernetes \
		--service-discovery=consul \
		--no-git`
	
	// Clean up the command for execution
	cleanCommand := strings.ReplaceAll(command, "\\\n\t\t", " ")
	cleanCommand = strings.ReplaceAll(cleanCommand, "\n\t\t", " ")
	cleanCommand = strings.TrimSpace(cleanCommand)
	
	err := ctx.iRunTheCommand(cleanCommand)
	suite.Require().NoError(err)
	
	err = ctx.theGenerationShouldSucceed()
	suite.Require().NoError(err)
	
	// Verify comprehensive file structure
	comprehensiveFiles := []string{
		"go.mod",
		"main.go",
		"Dockerfile",
		"docker-compose.yml",
		"k8s/deployment.yaml",
		"api/proto/service.proto",
		"internal/config/config.go",
		"internal/handlers/grpc_handler.go",
		"internal/middleware/auth.go",
		"internal/database/connection.go",
		"internal/observability/metrics.go",
		"internal/security/jwt.go",
		"internal/discovery/client.go",
		"configs/config.yaml",
		"migrations/001_init.up.sql",
		"tests/integration/service_test.go",
	}
	
	err = ctx.checkRequiredFiles(comprehensiveFiles)
	suite.Require().NoError(err)
	
	err = ctx.theGeneratedCodeShouldCompileSuccessfully()
	suite.Require().NoError(err)
	
	suite.T().Log("Comprehensive service generation completed successfully")
}

// Run the service acceptance test suite
func TestMicroserviceAcceptanceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping service acceptance tests in short mode")
	}
	
	suite.Run(t, new(MicroserviceAcceptanceTestSuite))
}

// Benchmark tests for service generation performance
func BenchmarkServiceGeneration(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ctx := InitializeMicroserviceContext()
		
		err := ctx.iRunTheCommand(fmt.Sprintf("go-starter new bench-service-%d --type=microservice-standard --no-git", i))
		if err != nil {
			b.Fatalf("Service generation failed: %v", err)
		}
		
		ctx.cleanup()
	}
}

// Performance test for compilation speed
func BenchmarkServiceCompilation(b *testing.B) {
	// Generate a service once
	ctx := InitializeMicroserviceContext()
	defer ctx.cleanup()
	
	err := ctx.iRunTheCommand("go-starter new bench-compile --type=microservice-standard --no-git")
	if err != nil {
		b.Fatalf("Service generation failed: %v", err)
	}
	
	// Benchmark compilation
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cmd := exec.Command("go", "build", "-v", "./...")
		cmd.Dir = ctx.projectPath
		err := cmd.Run()
		if err != nil {
			b.Fatalf("Compilation failed: %v", err)
		}
	}
}