package workspace_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

// WorkspaceIntegrationTestSuite contains integration tests for workspace blueprint
type WorkspaceIntegrationTestSuite struct {
	suite.Suite
	tempDir           string
	projectDir        string
	cliPath           string
	generatedProjects map[string]string // configuration -> project path
}

// SetupSuite prepares the test suite
func (s *WorkspaceIntegrationTestSuite) SetupSuite() {
	var err error
	
	// Create temporary directory
	s.tempDir, err = ioutil.TempDir("", "workspace-integration-*")
	s.Require().NoError(err)
	
	// Build CLI tool
	s.cliPath = filepath.Join(s.tempDir, "go-starter")
	buildCmd := exec.Command("go", "build", "-o", s.cliPath, "../../../../main.go")
	err = buildCmd.Run()
	s.Require().NoError(err, "Failed to build CLI tool")
	
	// Initialize generated projects map
	s.generatedProjects = make(map[string]string)
}

// TearDownSuite cleans up after tests
func (s *WorkspaceIntegrationTestSuite) TearDownSuite() {
	if s.tempDir != "" {
		os.RemoveAll(s.tempDir)
	}
}

// SetupTest prepares each test
func (s *WorkspaceIntegrationTestSuite) SetupTest() {
	// Change to temp directory for each test
	err := os.Chdir(s.tempDir)
	s.Require().NoError(err)
}

// TestWorkspaceBasicGeneration tests basic workspace generation
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceBasicGeneration() {
	projectName := "test-workspace-basic"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate workspace project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "workspace",
		"--module", "github.com/test/workspace-basic",
		"--framework", "gin",
		"--logger", "slog",
		"--enable-web-api", "true",
		"--enable-cli", "true",
		"--database-type", "postgres",
		"--message-queue", "redis",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify project structure
	s.verifyBasicWorkspaceStructure(projectPath)
	
	// Verify workspace synchronization
	s.verifyWorkspaceSync(projectPath)
	
	// Verify compilation
	s.verifyWorkspaceCompilation(projectPath)
	
	// Store for later use
	s.generatedProjects["basic"] = projectPath
}

// TestWorkspaceMinimalGeneration tests minimal workspace generation
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceMinimalGeneration() {
	projectName := "test-workspace-minimal"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate minimal workspace project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "workspace",
		"--module", "github.com/test/workspace-minimal",
		"--enable-web-api", "false",
		"--enable-cli", "true",
		"--enable-worker", "false",
		"--enable-microservices", "false",
		"--database-type", "none",
		"--message-queue", "none",
		"--enable-docker", "false",
		"--enable-kubernetes", "false",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify minimal structure
	s.verifyMinimalWorkspaceStructure(projectPath)
	
	// Verify unwanted components are not present
	s.verifyUnwantedComponentsAbsent(projectPath)
	
	// Verify compilation
	s.verifyWorkspaceCompilation(projectPath)
	
	// Store for later use
	s.generatedProjects["minimal"] = projectPath
}

// TestWorkspaceWithAllModules tests workspace with all modules enabled
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceWithAllModules() {
	projectName := "test-workspace-full"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate full workspace project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "workspace",
		"--module", "github.com/test/workspace-full",
		"--framework", "echo",
		"--logger", "zap",
		"--enable-web-api", "true",
		"--enable-cli", "true",
		"--enable-worker", "true",
		"--enable-microservices", "true",
		"--database-type", "postgres",
		"--message-queue", "redis",
		"--enable-docker", "true",
		"--enable-kubernetes", "true",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify full structure
	s.verifyFullWorkspaceStructure(projectPath)
	
	// Verify Docker configuration
	s.verifyDockerConfiguration(projectPath)
	
	// Verify Kubernetes configuration
	s.verifyKubernetesConfiguration(projectPath)
	
	// Verify compilation
	s.verifyWorkspaceCompilation(projectPath)
	
	// Store for later use
	s.generatedProjects["full"] = projectPath
}

// TestWorkspaceMultiFramework tests workspace generation with different frameworks
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceMultiFramework() {
	frameworks := []string{"gin", "echo", "fiber", "chi"}
	
	for _, framework := range frameworks {
		s.Run(fmt.Sprintf("Framework_%s", framework), func() {
			projectName := fmt.Sprintf("test-workspace-%s", framework)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "workspace",
				"--module", fmt.Sprintf("github.com/test/workspace-%s", framework),
				"--framework", framework,
				"--enable-web-api", "true",
				"--enable-cli", "true",
				"--logger", "slog",
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for %s: %s", framework, string(output))
			
			// Verify framework-specific imports
			s.verifyFrameworkImports(projectPath, framework)
			
			// Verify compilation
			s.verifyWorkspaceCompilation(projectPath)
			
			// Store project path
			s.generatedProjects[framework] = projectPath
		})
	}
}

// TestWorkspaceMultiDatabase tests workspace generation with different databases
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceMultiDatabase() {
	databases := []string{"postgres", "mysql", "mongodb", "sqlite"}
	
	for _, database := range databases {
		s.Run(fmt.Sprintf("Database_%s", database), func() {
			projectName := fmt.Sprintf("test-workspace-db-%s", database)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "workspace",
				"--module", fmt.Sprintf("github.com/test/workspace-db-%s", database),
				"--framework", "gin",
				"--database-type", database,
				"--enable-web-api", "true",
				"--logger", "slog",
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for database %s: %s", database, string(output))
			
			// Verify database-specific configuration
			s.verifyDatabaseConfiguration(projectPath, database)
			
			// Verify compilation
			s.verifyWorkspaceCompilation(projectPath)
		})
	}
}

// TestWorkspaceMultiMessageQueue tests workspace generation with different message queues
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceMultiMessageQueue() {
	messageQueues := []string{"redis", "nats", "kafka", "rabbitmq"}
	
	for _, mq := range messageQueues {
		s.Run(fmt.Sprintf("MessageQueue_%s", mq), func() {
			projectName := fmt.Sprintf("test-workspace-mq-%s", mq)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "workspace",
				"--module", fmt.Sprintf("github.com/test/workspace-mq-%s", mq),
				"--framework", "gin",
				"--message-queue", mq,
				"--enable-worker", "true",
				"--logger", "slog",
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for message queue %s: %s", mq, string(output))
			
			// Verify message queue configuration
			s.verifyMessageQueueConfiguration(projectPath, mq)
			
			// Verify compilation
			s.verifyWorkspaceCompilation(projectPath)
		})
	}
}

// TestWorkspaceMultiLogger tests workspace generation with different loggers
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceMultiLogger() {
	loggers := []string{"slog", "zap", "logrus", "zerolog"}
	
	for _, logger := range loggers {
		s.Run(fmt.Sprintf("Logger_%s", logger), func() {
			projectName := fmt.Sprintf("test-workspace-logger-%s", logger)
			projectPath := filepath.Join(s.tempDir, projectName)
			
			// Generate project
			cmd := exec.Command(s.cliPath, "new", projectName,
				"--type", "workspace",
				"--module", fmt.Sprintf("github.com/test/workspace-logger-%s", logger),
				"--framework", "gin",
				"--logger", logger,
				"--enable-web-api", "true",
			)
			
			output, err := cmd.CombinedOutput()
			s.Require().NoError(err, "Generation failed for logger %s: %s", logger, string(output))
			
			// Verify logger implementation
			s.verifyLoggerImplementation(projectPath, logger)
			
			// Verify compilation
			s.verifyWorkspaceCompilation(projectPath)
		})
	}
}

// TestWorkspaceBuildSystem tests the workspace build system
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceBuildSystem() {
	projectName := "test-workspace-build"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate workspace project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "workspace",
		"--module", "github.com/test/workspace-build",
		"--framework", "gin",
		"--enable-web-api", "true",
		"--enable-cli", "true",
		"--enable-worker", "true",
		"--database-type", "postgres",
		"--message-queue", "redis",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify Makefile targets
	s.verifyMakefileTargets(projectPath)
	
	// Test build targets
	s.testBuildTargets(projectPath)
	
	// Test build scripts
	s.testBuildScripts(projectPath)
}

// TestWorkspaceDependencyManagement tests module dependency management
func (s *WorkspaceIntegrationTestSuite) TestWorkspaceDependencyManagement() {
	projectName := "test-workspace-deps"
	projectPath := filepath.Join(s.tempDir, projectName)
	
	// Generate workspace project
	cmd := exec.Command(s.cliPath, "new", projectName,
		"--type", "workspace",
		"--module", "github.com/test/workspace-deps",
		"--framework", "gin",
		"--enable-web-api", "true",
		"--enable-cli", "true",
		"--database-type", "postgres",
		"--message-queue", "redis",
	)
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Generation failed: %s", string(output))
	
	// Verify go.work configuration
	s.verifyGoWorkConfiguration(projectPath)
	
	// Verify module dependencies
	s.verifyModuleDependencies(projectPath)
	
	// Test dependency resolution
	s.testDependencyResolution(projectPath)
}

// Helper methods

func (s *WorkspaceIntegrationTestSuite) verifyBasicWorkspaceStructure(projectPath string) {
	requiredFiles := []string{
		"go.work",
		"workspace.yaml",
		"Makefile",
		"pkg/shared/go.mod",
		"pkg/models/go.mod",
		"pkg/storage/go.mod",
		"pkg/events/go.mod",
		"cmd/api/go.mod",
		"cmd/cli/go.mod",
		"scripts/build-all.sh",
		"scripts/test-all.sh",
		"scripts/lint-all.sh",
	}
	
	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Required file should exist: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyMinimalWorkspaceStructure(projectPath string) {
	requiredFiles := []string{
		"go.work",
		"workspace.yaml",
		"Makefile",
		"pkg/shared/go.mod",
		"pkg/models/go.mod",
		"cmd/cli/go.mod",
	}
	
	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Required minimal file should exist: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyFullWorkspaceStructure(projectPath string) {
	requiredFiles := []string{
		"go.work",
		"workspace.yaml",
		"Makefile",
		"pkg/shared/go.mod",
		"pkg/models/go.mod",
		"pkg/storage/go.mod",
		"pkg/events/go.mod",
		"cmd/api/go.mod",
		"cmd/cli/go.mod",
		"cmd/worker/go.mod",
		"services/user-service/go.mod",
		"services/notification-service/go.mod",
		"docker-compose.yml",
		"docker-compose.dev.yml",
		"deployments/k8s/namespace.yaml",
		"deployments/k8s/configmap.yaml",
	}
	
	for _, file := range requiredFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Required full file should exist: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyUnwantedComponentsAbsent(projectPath string) {
	unwantedFiles := []string{
		"pkg/storage",
		"pkg/events",
		"cmd/api",
		"cmd/worker",
		"services",
		"docker-compose.yml",
		"deployments/k8s",
	}
	
	for _, file := range unwantedFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().NoFileExists(fullPath, "Unwanted file should not exist: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyWorkspaceSync(projectPath string) {
	cmd := exec.Command("go", "work", "sync")
	cmd.Dir = projectPath
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Workspace sync should succeed: %s", string(output))
}

func (s *WorkspaceIntegrationTestSuite) verifyWorkspaceCompilation(projectPath string) {
	// First sync the workspace
	syncCmd := exec.Command("go", "work", "sync")
	syncCmd.Dir = projectPath
	err := syncCmd.Run()
	s.Require().NoError(err, "go work sync should succeed")
	
	// Then build all modules
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectPath
	output, err := buildCmd.CombinedOutput()
	s.Require().NoError(err, "Workspace should compile: %s", string(output))
}

func (s *WorkspaceIntegrationTestSuite) verifyFrameworkImports(projectPath, framework string) {
	if !s.fileExists(filepath.Join(projectPath, "cmd/api")) {
		return // Skip if API module not enabled
	}
	
	apiMainPath := filepath.Join(projectPath, "cmd/api/main.go")
	content, err := ioutil.ReadFile(apiMainPath)
	s.Require().NoError(err)
	
	expectedImports := map[string]string{
		"gin":   "github.com/gin-gonic/gin",
		"echo":  "github.com/labstack/echo",
		"fiber": "github.com/gofiber/fiber",
		"chi":   "github.com/go-chi/chi",
	}
	
	if expectedImport, exists := expectedImports[framework]; exists {
		s.Contains(string(content), expectedImport, "Should contain framework import for %s", framework)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyDatabaseConfiguration(projectPath, database string) {
	// Verify storage package exists
	storageDir := filepath.Join(projectPath, "pkg/storage")
	s.Require().DirExists(storageDir, "Storage package should exist for database %s", database)
	
	// Verify database-specific dependencies
	goModPath := filepath.Join(storageDir, "go.mod")
	content, err := ioutil.ReadFile(goModPath)
	s.Require().NoError(err)
	
	expectedDrivers := map[string]string{
		"postgres": "github.com/lib/pq",
		"mysql":    "github.com/go-sql-driver/mysql",
		"mongodb":  "go.mongodb.org/mongo-driver",
		"sqlite":   "github.com/mattn/go-sqlite3",
	}
	
	if expectedDriver, exists := expectedDrivers[database]; exists {
		s.Contains(string(content), expectedDriver, "Should contain database driver for %s", database)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyMessageQueueConfiguration(projectPath, mq string) {
	// Verify events package exists
	eventsDir := filepath.Join(projectPath, "pkg/events")
	s.Require().DirExists(eventsDir, "Events package should exist for message queue %s", mq)
	
	// Verify message queue specific dependencies
	goModPath := filepath.Join(eventsDir, "go.mod")
	content, err := ioutil.ReadFile(goModPath)
	s.Require().NoError(err)
	
	expectedLibraries := map[string]string{
		"redis":    "github.com/go-redis/redis",
		"nats":     "github.com/nats-io/nats.go",
		"kafka":    "github.com/segmentio/kafka-go",
		"rabbitmq": "github.com/streadway/amqp",
	}
	
	if expectedLib, exists := expectedLibraries[mq]; exists {
		s.Contains(string(content), expectedLib, "Should contain message queue library for %s", mq)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyLoggerImplementation(projectPath, logger string) {
	loggerPath := filepath.Join(projectPath, "pkg/shared/logger/logger.go")
	content, err := ioutil.ReadFile(loggerPath)
	s.Require().NoError(err)
	
	expectedPackages := map[string]string{
		"slog":    "log/slog",
		"zap":     "go.uber.org/zap",
		"logrus":  "github.com/sirupsen/logrus",
		"zerolog": "github.com/rs/zerolog",
	}
	
	if expectedPackage, exists := expectedPackages[logger]; exists {
		s.Contains(string(content), expectedPackage, "Should contain logger package for %s", logger)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyDockerConfiguration(projectPath string) {
	dockerFiles := []string{
		"docker-compose.yml",
		"docker-compose.dev.yml",
	}
	
	for _, file := range dockerFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Docker file should exist: %s", file)
		
		// Verify YAML syntax
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		
		var dockerCompose interface{}
		err = yaml.Unmarshal(content, &dockerCompose)
		s.Require().NoError(err, "Docker file should be valid YAML: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyKubernetesConfiguration(projectPath string) {
	k8sFiles := []string{
		"deployments/k8s/namespace.yaml",
		"deployments/k8s/configmap.yaml",
		"deployments/k8s/secrets.yaml",
	}
	
	for _, file := range k8sFiles {
		fullPath := filepath.Join(projectPath, file)
		s.Require().FileExists(fullPath, "Kubernetes file should exist: %s", file)
		
		// Verify YAML syntax
		content, err := ioutil.ReadFile(fullPath)
		s.Require().NoError(err)
		
		var k8sManifest interface{}
		err = yaml.Unmarshal(content, &k8sManifest)
		s.Require().NoError(err, "Kubernetes file should be valid YAML: %s", file)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyMakefileTargets(projectPath string) {
	makefilePath := filepath.Join(projectPath, "Makefile")
	content, err := ioutil.ReadFile(makefilePath)
	s.Require().NoError(err)
	
	makefileStr := string(content)
	expectedTargets := []string{
		"build-all:",
		"test-all:",
		"lint-all:",
		"clean-all:",
		"deps-update:",
	}
	
	for _, target := range expectedTargets {
		s.Contains(makefileStr, target, "Makefile should contain target: %s", target)
	}
}

func (s *WorkspaceIntegrationTestSuite) testBuildTargets(projectPath string) {
	targets := []string{"build-all", "clean-all"}
	
	for _, target := range targets {
		cmd := exec.Command("make", target)
		cmd.Dir = projectPath
		
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		
		output, err := cmd.CombinedOutput()
		s.Require().NoError(err, "Make target %s should succeed: %s", target, string(output))
	}
}

func (s *WorkspaceIntegrationTestSuite) testBuildScripts(projectPath string) {
	scripts := []string{
		"scripts/build-all.sh",
		"scripts/test-all.sh",
		"scripts/lint-all.sh",
	}
	
	for _, script := range scripts {
		fullPath := filepath.Join(projectPath, script)
		
		// Verify script is executable
		info, err := os.Stat(fullPath)
		s.Require().NoError(err)
		s.True(info.Mode()&0111 != 0, "Script should be executable: %s", script)
	}
}

func (s *WorkspaceIntegrationTestSuite) verifyGoWorkConfiguration(projectPath string) {
	goWorkPath := filepath.Join(projectPath, "go.work")
	content, err := ioutil.ReadFile(goWorkPath)
	s.Require().NoError(err)
	
	goWorkStr := string(content)
	s.Contains(goWorkStr, "go ")
	s.Contains(goWorkStr, "use (")
	s.Contains(goWorkStr, "./pkg/shared")
	s.Contains(goWorkStr, "./pkg/models")
}

func (s *WorkspaceIntegrationTestSuite) verifyModuleDependencies(projectPath string) {
	// Verify that dependent modules reference shared modules properly
	apiGoModPath := filepath.Join(projectPath, "cmd/api/go.mod")
	if s.fileExists(apiGoModPath) {
		content, err := ioutil.ReadFile(apiGoModPath)
		s.Require().NoError(err)
		
		// Should contain replace directives for local modules
		goModStr := string(content)
		s.Contains(goModStr, "replace ")
	}
}

func (s *WorkspaceIntegrationTestSuite) testDependencyResolution(projectPath string) {
	// Test that go work sync resolves dependencies correctly
	cmd := exec.Command("go", "work", "sync")
	cmd.Dir = projectPath
	
	output, err := cmd.CombinedOutput()
	s.Require().NoError(err, "Dependency resolution should succeed: %s", string(output))
}

func (s *WorkspaceIntegrationTestSuite) fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// TestSuite runner
func TestWorkspaceIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	
	suite.Run(t, new(WorkspaceIntegrationTestSuite))
}