package microservice

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// MicroserviceTestContext holds the test execution context for microservice blueprints
type MicroserviceTestContext struct {
	// Project generation
	projectName     string
	projectType     string
	projectPath     string
	generatedFiles  []string
	lastCommand     *exec.Cmd
	lastOutput      string
	lastError       string
	lastExitCode    int

	// Service runtime
	serviceURL      string
	httpClient      *http.Client
	lastResponse    *http.Response
	lastResponseBody []byte
	serviceProcess  *exec.Cmd
	servicePID      int

	// Container testing
	containers      map[string]testcontainers.Container
	dbHost          string
	dbPort          string
	messageQueueURL string
	tracingURL      string

	// Test data
	testData        map[string]interface{}
	scenarios       map[string]interface{}
	loadTestResults map[string]interface{}

	// Performance tracking
	responseTime    time.Duration
	requestStart    time.Time
}

// Global test context
var microserviceCtx *MicroserviceTestContext

// Initialize context for microservice blueprint testing
func InitializeMicroserviceContext() *MicroserviceTestContext {
	if microserviceCtx == nil {
		microserviceCtx = &MicroserviceTestContext{
			httpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
			containers:      make(map[string]testcontainers.Container),
			testData:        make(map[string]interface{}),
			scenarios:       make(map[string]interface{}),
			loadTestResults: make(map[string]interface{}),
		}
	}
	return microserviceCtx
}

// Service Generation Steps

func (ctx *MicroserviceTestContext) iWantToCreateAGRPCBasedMicroservice() error {
	ctx.projectType = "microservice-standard"
	ctx.scenarios["communication"] = "grpc"
	return nil
}

func (ctx *MicroserviceTestContext) iWantToCreateAnAWSLambdaServerlessFunction() error {
	ctx.projectType = "lambda-standard"
	ctx.scenarios["runtime"] = "aws"
	return nil
}

func (ctx *MicroserviceTestContext) iWantToCreateMicroservicesWithVariousProtocols() error {
	ctx.projectType = "microservice-standard"
	ctx.scenarios["multi_protocol"] = true
	return nil
}

func (ctx *MicroserviceTestContext) iRunTheCommand(command string) error {
	// Parse the command and extract project name
	parts := strings.Fields(command)
	if len(parts) < 3 {
		return fmt.Errorf("invalid command format: %s", command)
	}

	ctx.projectName = parts[2] // Extract project name from "go-starter new my-service ..."
	
	// Create temporary directory for project
	tempDir := os.TempDir()
	ctx.projectPath = filepath.Join(tempDir, ctx.projectName)
	
	// Remove existing directory if it exists
	os.RemoveAll(ctx.projectPath)

	// Prepare command with working directory
	fullCommand := strings.Join(parts, " ")
	cmd := exec.Command("sh", "-c", fullCommand)
	cmd.Dir = tempDir

	// Capture output
	output, err := cmd.CombinedOutput()
	ctx.lastOutput = string(output)
	ctx.lastExitCode = cmd.ProcessState.ExitCode()

	if err != nil {
		ctx.lastError = err.Error()
		return fmt.Errorf("command failed: %s, output: %s", err.Error(), ctx.lastOutput)
	}

	return nil
}

func (ctx *MicroserviceTestContext) theGenerationShouldSucceed() error {
	if ctx.lastExitCode != 0 {
		return fmt.Errorf("generation failed with exit code %d: %s", ctx.lastExitCode, ctx.lastOutput)
	}
	return nil
}

func (ctx *MicroserviceTestContext) theProjectShouldContainAllEssentialMicroserviceComponents() error {
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"internal/app/app.go",
		"internal/config/config.go",
		"internal/handlers/grpc_handler.go",
		"internal/middleware/auth.go",
		"internal/middleware/logging.go",
		"internal/services/health.go",
		"Dockerfile",
		"docker-compose.yml",
		"k8s/deployment.yaml",
		"k8s/service.yaml",
	}

	return ctx.checkRequiredFiles(requiredFiles)
}

func (ctx *MicroserviceTestContext) theProjectShouldContainAllEssentialServerlessComponents() error {
	requiredFiles := []string{
		"go.mod",
		"main.go",
		"handler.go",
		"internal/config/config.go",
		"internal/services/processor.go",
		"template.yaml", // SAM template
		"Dockerfile",
		"event-examples/api-gateway.json",
		"event-examples/s3.json",
	}

	return ctx.checkRequiredFiles(requiredFiles)
}

func (ctx *MicroserviceTestContext) theGeneratedCodeShouldCompileSuccessfully() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("no project path available")
	}

	// Change to project directory and run go build
	cmd := exec.Command("go", "build", "-v", "./...")
	cmd.Dir = ctx.projectPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("compilation failed: %s, output: %s", err.Error(), string(output))
	}

	return nil
}

// Service Runtime Steps

func (ctx *MicroserviceTestContext) theServiceShouldSupportGRPCAndHTTPInterfaces() error {
	// Check for gRPC and HTTP server implementations
	grpcFile := filepath.Join(ctx.projectPath, "internal/grpc/server.go")
	httpFile := filepath.Join(ctx.projectPath, "internal/http/server.go")

	if _, err := os.Stat(grpcFile); os.IsNotExist(err) {
		return fmt.Errorf("gRPC server implementation not found: %s", grpcFile)
	}

	if _, err := os.Stat(httpFile); os.IsNotExist(err) {
		return fmt.Errorf("HTTP server implementation not found: %s", httpFile)
	}

	return nil
}

func (ctx *MicroserviceTestContext) theMicroserviceShouldBeContainerReadyWithHealthChecks() error {
	// Check Dockerfile exists and contains health check
	dockerfilePath := filepath.Join(ctx.projectPath, "Dockerfile")
	content, err := os.ReadFile(dockerfilePath)
	if err != nil {
		return fmt.Errorf("Dockerfile not found: %s", err.Error())
	}

	dockerfileContent := string(content)
	if !strings.Contains(dockerfileContent, "HEALTHCHECK") {
		return fmt.Errorf("Dockerfile missing HEALTHCHECK instruction")
	}

	// Check for health endpoint implementation
	healthFile := filepath.Join(ctx.projectPath, "internal/handlers/health.go")
	if _, err := os.Stat(healthFile); os.IsNotExist(err) {
		return fmt.Errorf("health handler not found: %s", healthFile)
	}

	return nil
}

func (ctx *MicroserviceTestContext) theLambdaShouldSupportAWSDKIntegration() error {
	// Check for AWS SDK imports and usage
	mainFile := filepath.Join(ctx.projectPath, "main.go")
	content, err := os.ReadFile(mainFile)
	if err != nil {
		return fmt.Errorf("main.go not found: %s", err.Error())
	}

	fileContent := string(content)
	if !strings.Contains(fileContent, "github.com/aws/aws-lambda-go") {
		return fmt.Errorf("AWS Lambda Go SDK not imported")
	}

	return nil
}

func (ctx *MicroserviceTestContext) theFunctionShouldIncludeObservabilityFeatures() error {
	// Check for observability configuration
	configFile := filepath.Join(ctx.projectPath, "internal/config/observability.go")
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("observability configuration not found")
	}

	// Check for metrics and tracing
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read observability config: %s", err.Error())
	}

	configContent := string(content)
	requiredFeatures := []string{"metrics", "tracing", "logging"}
	for _, feature := range requiredFeatures {
		if !strings.Contains(strings.ToLower(configContent), feature) {
			return fmt.Errorf("observability feature missing: %s", feature)
		}
	}

	return nil
}

// Protocol Testing Steps

func (ctx *MicroserviceTestContext) iGenerateAMicroserviceWithProtocol(protocol string) error {
	ctx.scenarios["protocol"] = protocol
	
	command := fmt.Sprintf("go-starter new test-%s --type=microservice-standard --protocol=%s --no-git", 
		protocol, protocol)
	
	return ctx.iRunTheCommand(command)
}

func (ctx *MicroserviceTestContext) theProjectShouldSupportTheProtocolCommunicationPattern(protocol string) error {
	var expectedFiles []string
	
	switch protocol {
	case "grpc":
		expectedFiles = []string{
			"api/proto/service.proto",
			"internal/grpc/server.go",
			"internal/grpc/client.go",
		}
	case "rest":
		expectedFiles = []string{
			"internal/http/server.go",
			"internal/http/router.go",
			"api/openapi.yaml",
		}
	case "graphql":
		expectedFiles = []string{
			"internal/graphql/schema.go",
			"internal/graphql/resolver.go",
			"schema.graphql",
		}
	default:
		return fmt.Errorf("unsupported protocol: %s", protocol)
	}

	return ctx.checkRequiredFiles(expectedFiles)
}

func (ctx *MicroserviceTestContext) theServiceShouldIncludeAppropriateClientLibraries() error {
	// Check go.mod for protocol-specific dependencies
	goModPath := filepath.Join(ctx.projectPath, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("go.mod not found: %s", err.Error())
	}

	goModContent := string(content)
	protocol := ctx.scenarios["protocol"].(string)

	switch protocol {
	case "grpc":
		if !strings.Contains(goModContent, "google.golang.org/grpc") {
			return fmt.Errorf("gRPC client library not found in go.mod")
		}
	case "rest":
		// REST typically uses standard HTTP client
		if !strings.Contains(goModContent, "github.com/gin-gonic/gin") &&
			!strings.Contains(goModContent, "github.com/labstack/echo") {
			return fmt.Errorf("HTTP framework not found in go.mod")
		}
	case "graphql":
		if !strings.Contains(goModContent, "github.com/99designs/gqlgen") {
			return fmt.Errorf("GraphQL library not found in go.mod")
		}
	}

	return nil
}

func (ctx *MicroserviceTestContext) theProtocolSpecificMiddlewareShouldBeConfigured() error {
	middlewareDir := filepath.Join(ctx.projectPath, "internal/middleware")
	
	// Check if middleware directory exists
	if _, err := os.Stat(middlewareDir); os.IsNotExist(err) {
		return fmt.Errorf("middleware directory not found")
	}

	// Check for protocol-specific middleware files
	protocol := ctx.scenarios["protocol"].(string)
	var expectedMiddleware []string

	switch protocol {
	case "grpc":
		expectedMiddleware = []string{"grpc_auth.go", "grpc_logging.go", "grpc_recovery.go"}
	case "rest":
		expectedMiddleware = []string{"cors.go", "auth.go", "logging.go"}
	case "graphql":
		expectedMiddleware = []string{"graphql_auth.go", "complexity.go", "query_cache.go"}
	}

	for _, middleware := range expectedMiddleware {
		middlewarePath := filepath.Join(middlewareDir, middleware)
		if _, err := os.Stat(middlewarePath); os.IsNotExist(err) {
			return fmt.Errorf("middleware file not found: %s", middleware)
		}
	}

	return nil
}

func (ctx *MicroserviceTestContext) theServiceShouldCompileAndRunCorrectly() error {
	return ctx.theGeneratedCodeShouldCompileSuccessfully()
}

// Database Integration Steps

func (ctx *MicroserviceTestContext) iGenerateAServiceWithDatabaseAndORM(database, orm string) error {
	ctx.scenarios["database"] = database
	ctx.scenarios["orm"] = orm
	
	command := fmt.Sprintf("go-starter new test-db-%s --type=microservice-standard --database=%s --orm=%s --no-git", 
		database, database, orm)
	
	return ctx.iRunTheCommand(command)
}

func (ctx *MicroserviceTestContext) theProjectShouldIncludeDatabaseConfiguration() error {
	configFiles := []string{
		"internal/config/database.go",
		"internal/database/connection.go",
	}
	
	return ctx.checkRequiredFiles(configFiles)
}

func (ctx *MicroserviceTestContext) theServiceShouldSupportDatabaseMigrations() error {
	migrationFiles := []string{
		"migrations/001_init.up.sql",
		"migrations/001_init.down.sql",
		"internal/database/migrations.go",
	}
	
	return ctx.checkRequiredFiles(migrationFiles)
}

func (ctx *MicroserviceTestContext) theRepositoryLayerShouldUseTheSpecifiedORM() error {
	repoDir := filepath.Join(ctx.projectPath, "internal/repository")
	if _, err := os.Stat(repoDir); os.IsNotExist(err) {
		return fmt.Errorf("repository directory not found")
	}

	orm := ctx.scenarios["orm"].(string)
	
	// Check for ORM-specific implementation
	var expectedFile string
	switch orm {
	case "gorm":
		expectedFile = "gorm_repository.go"
	case "sqlx":
		expectedFile = "sqlx_repository.go"
	default:
		expectedFile = "repository.go"
	}

	repoPath := filepath.Join(repoDir, expectedFile)
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("ORM-specific repository not found: %s", expectedFile)
	}

	return nil
}

func (ctx *MicroserviceTestContext) theDatabaseConnectionShouldBeTestableWithContainers() error {
	// Start database container for testing
	database := ctx.scenarios["database"].(string)
	
	container, err := ctx.startDatabaseContainer(database)
	if err != nil {
		return fmt.Errorf("failed to start database container: %s", err.Error())
	}
	
	ctx.containers["database"] = container
	
	// Test database connectivity
	return ctx.testDatabaseConnectivity(database)
}

func (ctx *MicroserviceTestContext) theServiceShouldHandleConnectionPooling() error {
	// Check for connection pooling configuration
	configFile := filepath.Join(ctx.projectPath, "internal/config/database.go")
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("database config not found: %s", err.Error())
	}

	configContent := string(content)
	poolingFeatures := []string{"MaxOpenConns", "MaxIdleConns", "ConnMaxLifetime"}
	
	for _, feature := range poolingFeatures {
		if !strings.Contains(configContent, feature) {
			return fmt.Errorf("connection pooling feature missing: %s", feature)
		}
	}

	return nil
}

// Service Discovery and Mesh Steps

func (ctx *MicroserviceTestContext) iWantMicroservicesThatCanDiscoverEachOther() error {
	ctx.scenarios["service_discovery"] = true
	return nil
}

func (ctx *MicroserviceTestContext) iGenerateAMicroserviceWithServiceDiscovery() error {
	command := "go-starter new test-discovery --type=microservice-standard --service-discovery=consul --no-git"
	return ctx.iRunTheCommand(command)
}

func (ctx *MicroserviceTestContext) theServiceShouldIncludeServiceRegistryIntegration() error {
	discoveryFiles := []string{
		"internal/discovery/client.go",
		"internal/discovery/registry.go",
		"internal/config/discovery.go",
	}
	
	return ctx.checkRequiredFiles(discoveryFiles)
}

func (ctx *MicroserviceTestContext) theServiceShouldSupportHealthCheckEndpoints() error {
	healthFiles := []string{
		"internal/handlers/health.go",
		"internal/health/checker.go",
	}
	
	for _, file := range healthFiles {
		filePath := filepath.Join(ctx.projectPath, file)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("health file not found: %s", file)
		}
		
		// Check for health check endpoints
		content, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read health file: %s", err.Error())
		}
		
		healthContent := string(content)
		endpoints := []string{"/health", "/health/live", "/health/ready"}
		
		found := false
		for _, endpoint := range endpoints {
			if strings.Contains(healthContent, endpoint) {
				found = true
				break
			}
		}
		
		if !found {
			return fmt.Errorf("health endpoints not found in %s", file)
		}
	}
	
	return nil
}

func (ctx *MicroserviceTestContext) theServiceShouldHandleServiceRegistration() error {
	registryFile := filepath.Join(ctx.projectPath, "internal/discovery/registry.go")
	content, err := os.ReadFile(registryFile)
	if err != nil {
		return fmt.Errorf("registry file not found: %s", err.Error())
	}

	registryContent := string(content)
	registrationFeatures := []string{"Register", "Deregister", "Health"}
	
	for _, feature := range registrationFeatures {
		if !strings.Contains(registryContent, feature) {
			return fmt.Errorf("service registration feature missing: %s", feature)
		}
	}

	return nil
}

func (ctx *MicroserviceTestContext) theServiceShouldSupportLoadBalancing() error {
	lbFile := filepath.Join(ctx.projectPath, "internal/loadbalancer/client.go")
	if _, err := os.Stat(lbFile); os.IsNotExist(err) {
		return fmt.Errorf("load balancer implementation not found")
	}

	return nil
}

func (ctx *MicroserviceTestContext) theDiscoveryConfigurationShouldBeEnvironmentAware() error {
	configFile := filepath.Join(ctx.projectPath, "configs/config.yaml")
	content, err := os.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("config file not found: %s", err.Error())
	}

	configContent := string(content)
	envFeatures := []string{"development:", "staging:", "production:"}
	
	found := false
	for _, env := range envFeatures {
		if strings.Contains(configContent, env) {
			found = true
			break
		}
	}
	
	if !found {
		return fmt.Errorf("environment-aware configuration not found")
	}

	return nil
}

// Container and Infrastructure Steps

func (ctx *MicroserviceTestContext) startDatabaseContainer(dbType string) (testcontainers.Container, error) {
	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	var req testcontainers.ContainerRequest

	switch dbType {
	case "postgres":
		req = testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_DB":       "testdb",
				"POSTGRES_USER":     "testuser",
				"POSTGRES_PASSWORD": "testpass",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp"),
		}
	case "mysql":
		req = testcontainers.ContainerRequest{
			Image:        "mysql:8.0",
			ExposedPorts: []string{"3306/tcp"},
			Env: map[string]string{
				"MYSQL_DATABASE":      "testdb",
				"MYSQL_USER":          "testuser",
				"MYSQL_PASSWORD":      "testpass",
				"MYSQL_ROOT_PASSWORD": "rootpass",
			},
			WaitingFor: wait.ForListeningPort("3306/tcp"),
		}
	case "mongodb":
		req = testcontainers.ContainerRequest{
			Image:        "mongo:6",
			ExposedPorts: []string{"27017/tcp"},
			Env: map[string]string{
				"MONGO_INITDB_DATABASE": "testdb",
			},
			WaitingFor: wait.ForListeningPort("27017/tcp"),
		}
	case "redis":
		req = testcontainers.ContainerRequest{
			Image:        "redis:7-alpine",
			ExposedPorts: []string{"6379/tcp"},
			WaitingFor: wait.ForListeningPort("6379/tcp"),
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	container, err := testcontainers.GenericContainer(ctxTimeout, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err == nil {
		// Store connection details
		host, _ := container.Host(ctxTimeout)
		var port string
		switch dbType {
		case "postgres":
			mappedPort, _ := container.MappedPort(ctxTimeout, "5432")
			port = mappedPort.Port()
		case "mysql":
			mappedPort, _ := container.MappedPort(ctxTimeout, "3306")
			port = mappedPort.Port()
		case "mongodb":
			mappedPort, _ := container.MappedPort(ctxTimeout, "27017")
			port = mappedPort.Port()
		case "redis":
			mappedPort, _ := container.MappedPort(ctxTimeout, "6379")
			port = mappedPort.Port()
		}
		
		ctx.dbHost = host
		ctx.dbPort = port
	}

	return container, err
}

func (ctx *MicroserviceTestContext) testDatabaseConnectivity(dbType string) error {
	// Test basic connectivity to the database container
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(ctx.dbHost, ctx.dbPort), 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %s", err.Error())
	}
	defer conn.Close()

	return nil
}

// Utility Functions

func (ctx *MicroserviceTestContext) checkRequiredFiles(files []string) error {
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

// HTTP Request Helpers

func (ctx *MicroserviceTestContext) makeHTTPRequest(method, path string, body []byte) error {
	url := ctx.serviceURL + path
	
	var reqBody io.Reader
	if body != nil {
		reqBody = strings.NewReader(string(body))
	}
	
	ctx.requestStart = time.Now()
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return fmt.Errorf("failed to create request: %s", err.Error())
	}
	
	ctx.lastResponse, err = ctx.httpClient.Do(req)
	ctx.responseTime = time.Since(ctx.requestStart)
	
	if err != nil {
		return fmt.Errorf("request failed: %s", err.Error())
	}
	
	// Read response body
	ctx.lastResponseBody, err = io.ReadAll(ctx.lastResponse.Body)
	ctx.lastResponse.Body.Close()
	
	return err
}

// Cleanup Functions

func (ctx *MicroserviceTestContext) cleanup() {
	// Stop service process
	if ctx.serviceProcess != nil {
		_ = ctx.serviceProcess.Process.Kill()
		_ = ctx.serviceProcess.Wait()
	}
	
	// Stop all containers
	for name, container := range ctx.containers {
		if container != nil {
			container.Terminate(context.Background())
		}
		delete(ctx.containers, name)
	}
	
	// Clean up temporary files
	if ctx.projectPath != "" {
		os.RemoveAll(ctx.projectPath)
	}
	
	// Reset context
	ctx.projectName = ""
	ctx.projectPath = ""
	ctx.serviceURL = ""
	ctx.lastResponse = nil
	ctx.lastResponseBody = nil
	ctx.testData = make(map[string]interface{})
	ctx.scenarios = make(map[string]interface{})
}

// Scenario hooks

func (ctx *MicroserviceTestContext) beforeScenario(sc *godog.Scenario) error {
	// Reset state before each scenario
	ctx.testData = make(map[string]interface{})
	ctx.scenarios = make(map[string]interface{})
	ctx.lastResponse = nil
	ctx.lastResponseBody = nil
	return nil
}

func (ctx *MicroserviceTestContext) afterScenario(sc *godog.Scenario, err error) error {
	// Cleanup after each scenario
	if err != nil {
		fmt.Printf("Scenario failed: %s - %v\n", sc.Name, err)
	}
	
	// Cleanup containers and processes for this scenario
	for name, container := range ctx.containers {
		if container != nil {
			container.Terminate(context.Background())
		}
		delete(ctx.containers, name)
	}
	
	if ctx.serviceProcess != nil {
		ctx.serviceProcess.Process.Kill()
		ctx.serviceProcess.Wait()
		ctx.serviceProcess = nil
	}
	return nil
}