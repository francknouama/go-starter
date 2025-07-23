package workspace_test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"gopkg.in/yaml.v3"
)

// TestContext holds the test execution context
type TestContext struct {
	t               *testing.T
	tempDir         string
	projectDir      string
	projectName     string
	generationError error
	parameters      map[string]string
	frameworks      []string
	databases       []string
	messageQueues   []string
	loggers         []string
	currentFramework string
}

// NewTestContext creates a new test context
func NewTestContext(t *testing.T) *TestContext {
	return &TestContext{
		t:             t,
		parameters:    make(map[string]string),
		frameworks:    []string{"gin", "echo", "fiber", "chi"},
		databases:     []string{"postgres", "mysql", "mongodb", "sqlite"},
		messageQueues: []string{"redis", "nats", "kafka", "rabbitmq"},
		loggers:       []string{"slog", "zap", "logrus", "zerolog"},
	}
}

// Cleanup removes temporary directories
func (tc *TestContext) Cleanup() {
	if tc.tempDir != "" {
		_ = os.RemoveAll(tc.tempDir)
	}
}

// Workspace ATDD Step Definitions

func (tc *TestContext) iHaveTheGoStarterCLIToolAvailable() error {
	// Find project root by looking for go.mod file
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	
	projectRoot := cwd
	for {
		if _, err := os.Stat(filepath.Join(projectRoot, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(projectRoot)
		if parent == projectRoot {
			return fmt.Errorf("could not find project root with go.mod")
		}
		projectRoot = parent
	}
	
	// Check if binary already exists in project root
	existingBinary := filepath.Join(projectRoot, "go-starter")
	cliPath := filepath.Join(tc.tempDir, "go-starter")
	
	if _, err := os.Stat(existingBinary); err == nil {
		// Copy existing binary
		data, err := os.ReadFile(existingBinary)
		if err != nil {
			return fmt.Errorf("failed to read existing binary: %w", err)
		}
		
		err = os.WriteFile(cliPath, data, 0755)
		if err != nil {
			return fmt.Errorf("failed to copy binary: %w", err)
		}
	} else {
		// Build the CLI tool
		buildCmd := exec.Command("go", "build", "-o", cliPath, ".")
		buildCmd.Dir = projectRoot
		
		if output, err := buildCmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to build CLI tool: %w\nOutput: %s", err, output)
		}
	}
	
	// Verify the tool is executable
	if _, err := os.Stat(cliPath); err != nil {
		return fmt.Errorf("CLI tool not found at %s: %w", cliPath, err)
	}
	
	return nil
}

func (tc *TestContext) iAmInATemporaryWorkingDirectory() error {
	var err error
	tc.tempDir, err = os.MkdirTemp("", "workspace-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	if err := os.Chdir(tc.tempDir); err != nil {
		return fmt.Errorf("failed to change to temp directory: %w", err)
	}
	
	return nil
}

func (tc *TestContext) iWantToGenerateAWorkspaceBlueprint() error {
	tc.parameters["type"] = "workspace"
	return nil
}

func (tc *TestContext) iRunTheGeneratorWith(table *godog.Table) error {
	// Parse parameters from table
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			tc.parameters[row.Cells[0].Value] = row.Cells[1].Value
		}
	}
	
	// Set project name for directory creation
	if name, exists := tc.parameters["name"]; exists {
		tc.projectName = name
		tc.projectDir = filepath.Join(tc.tempDir, name)
	}
	
	// Build CLI command
	cliPath := filepath.Join(tc.tempDir, "go-starter")
	args := []string{"new"}
	
	// Add parameters as arguments
	for key, value := range tc.parameters {
		switch key {
		case "type":
			args = append(args, "--type", value)
		case "name":
			args = append(args, value) // Project name is positional
		case "module":
			args = append(args, "--module", value)
		case "go_version":
			args = append(args, "--go-version", value)
		case "framework":
			args = append(args, "--framework", value)
		case "database_type":
			args = append(args, "--database-type", value)
		case "message_queue":
			args = append(args, "--message-queue", value)
		case "logger_type":
			args = append(args, "--logger", value)
		case "enable_web_api":
			args = append(args, "--enable-web-api", value)
		case "enable_cli":
			args = append(args, "--enable-cli", value)
		case "enable_worker":
			args = append(args, "--enable-worker", value)
		case "enable_microservices":
			args = append(args, "--enable-microservices", value)
		case "enable_docker":
			args = append(args, "--enable-docker", value)
		case "enable_kubernetes":
			args = append(args, "--enable-kubernetes", value)
		}
	}
	
	// Run the generator
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	
	cmd := exec.CommandContext(ctx, cliPath, args...)
	cmd.Dir = tc.tempDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		tc.generationError = fmt.Errorf("generation failed: %w\nOutput: %s", err, output)
		return nil // Don't return error here, let subsequent steps check
	}
	
	return nil
}

func (tc *TestContext) theGenerationShouldSucceed() error {
	if tc.generationError != nil {
		return tc.generationError
	}
	
	// Verify project directory was created
	if tc.projectDir == "" {
		return fmt.Errorf("project directory not set")
	}
	
	if _, err := os.Stat(tc.projectDir); err != nil {
		return fmt.Errorf("project directory not created: %w", err)
	}
	
	return nil
}

func (tc *TestContext) theGeneratedWorkspaceShouldHaveTheGoWorkspaceStructure(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			filePath := row.Cells[0].Value
			fileType := row.Cells[1].Value
			
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			switch fileType {
			case "file":
				if _, err := os.Stat(fullPath); err != nil {
					return fmt.Errorf("expected file %s not found: %w", filePath, err)
				}
			case "directory":
				if info, err := os.Stat(fullPath); err != nil {
					return fmt.Errorf("expected directory %s not found: %w", filePath, err)
				} else if !info.IsDir() {
					return fmt.Errorf("%s is not a directory", filePath)
				}
			}
		}
	}
	
	return nil
}

func (tc *TestContext) allModulesShouldCompileSuccessfully() error {
	// Run go work sync first
	syncCmd := exec.Command("go", "work", "sync")
	syncCmd.Dir = tc.projectDir
	if output, err := syncCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go work sync failed: %w\nOutput: %s", err, output)
	}
	
	// Build all modules
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = tc.projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (tc *TestContext) theWorkspaceShouldSyncWithoutErrors() error {
	cmd := exec.Command("go", "work", "sync")
	cmd.Dir = tc.projectDir
	
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("workspace sync failed: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (tc *TestContext) theGeneratedWorkspaceShouldHaveMinimalStructure(table *godog.Table) error {
	return tc.theGeneratedWorkspaceShouldHaveTheGoWorkspaceStructure(table)
}

func (tc *TestContext) theWorkspaceShouldNotInclude(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			filePath := row.Cells[0].Value
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			if _, err := os.Stat(fullPath); err == nil {
				return fmt.Errorf("file/directory %s should not exist in minimal workspace", filePath)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) iRunTheGeneratorForEachFramework(table *godog.Table) error {
	frameworks := []string{}
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			frameworks = append(frameworks, row.Cells[0].Value)
		}
	}
	
	for _, framework := range frameworks {
		tc.currentFramework = framework
		tc.parameters["framework"] = framework
		tc.parameters["name"] = fmt.Sprintf("test-workspace-%s", framework)
		
		if err := tc.iRunTheGeneratorWith(nil); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for framework %s: %w", framework, err)
		}
	}
	
	return nil
}

func (tc *TestContext) iRunTheGeneratorForEachDatabase(table *godog.Table) error {
	databases := []string{}
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			databases = append(databases, row.Cells[0].Value)
		}
	}
	
	for _, database := range databases {
		tc.parameters["database_type"] = database
		tc.parameters["name"] = fmt.Sprintf("test-workspace-db-%s", database)
		
		if err := tc.iRunTheGeneratorWith(nil); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for database %s: %w", database, err)
		}
	}
	
	return nil
}

func (tc *TestContext) iRunTheGeneratorForEachMessageQueue(table *godog.Table) error {
	messageQueues := []string{}
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			messageQueues = append(messageQueues, row.Cells[0].Value)
		}
	}
	
	for _, mq := range messageQueues {
		tc.parameters["message_queue"] = mq
		tc.parameters["name"] = fmt.Sprintf("test-workspace-mq-%s", mq)
		
		if err := tc.iRunTheGeneratorWith(nil); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for message queue %s: %w", mq, err)
		}
	}
	
	return nil
}

func (tc *TestContext) iRunTheGeneratorForEachLogger(table *godog.Table) error {
	loggers := []string{}
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			loggers = append(loggers, row.Cells[0].Value)
		}
	}
	
	for _, logger := range loggers {
		tc.parameters["logger_type"] = logger
		tc.parameters["name"] = fmt.Sprintf("test-workspace-logger-%s", logger)
		
		if err := tc.iRunTheGeneratorWith(nil); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for logger %s: %w", logger, err)
		}
	}
	
	return nil
}

func (tc *TestContext) eachGeneratedWorkspaceShouldWithTable(table *godog.Table) error {
	for _, framework := range tc.frameworks {
		projectDir := filepath.Join(tc.tempDir, fmt.Sprintf("test-workspace-%s", framework))
		
		for _, row := range table.Rows {
			if len(row.Cells) >= 1 {
				validation := row.Cells[0].Value
				
				switch validation {
				case "compile successfully":
					if err := tc.verifyCompilation(projectDir); err != nil {
						return fmt.Errorf("compilation failed for %s: %w", framework, err)
					}
					
				case "contain framework-specific imports":
					if err := tc.verifyFrameworkSpecificImports(projectDir, framework); err != nil {
						return err
					}
					
				case "have consistent module structure":
					if err := tc.verifyConsistentModuleStructure(projectDir); err != nil {
						return err
					}
					
				case "include proper dependency management":
					if err := tc.verifyDependencyManagement(projectDir); err != nil {
						return err
					}
					
				case "include database-specific storage package":
					if err := tc.verifyDatabaseSpecificStorage(projectDir); err != nil {
						return err
					}
					
				case "contain appropriate database drivers":
					if err := tc.verifyDatabaseDrivers(projectDir); err != nil {
						return err
					}
					
				case "have correct connection configuration":
					if err := tc.verifyConnectionConfiguration(projectDir); err != nil {
						return err
					}
					
				case "include message queue events package":
					if err := tc.verifyMessageQueueEventsPackage(projectDir); err != nil {
						return err
					}
					
				case "contain appropriate MQ client libraries":
					if err := tc.verifyMQClientLibraries(projectDir); err != nil {
						return err
					}
					
				case "include logger-specific implementation":
					if err := tc.verifyLoggerImplementation(projectDir); err != nil {
						return err
					}
					
				case "have consistent logging interface":
					if err := tc.verifyConsistentLoggingInterface(projectDir); err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theGeneratedWorkspaceShouldIncludeDockerFiles(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			filePath := row.Cells[0].Value
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			if _, err := os.Stat(fullPath); err != nil {
				return fmt.Errorf("expected Docker file %s not found: %w", filePath, err)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theDockerConfigurationsShouldBeValid() error {
	dockerFiles := []string{
		"docker-compose.yml",
		"docker-compose.dev.yml",
	}
	
	for _, file := range dockerFiles {
		fullPath := filepath.Join(tc.projectDir, file)
		content, err := os.ReadFile(fullPath)
		if err != nil {
			return fmt.Errorf("failed to read Docker file %s: %w", file, err)
		}
		
		var dockerCompose interface{}
		if err := yaml.Unmarshal(content, &dockerCompose); err != nil {
			return fmt.Errorf("Docker file %s is not valid YAML: %w", file, err)
		}
	}
	
	return nil
}

func (tc *TestContext) theGeneratedWorkspaceShouldIncludeKubernetesManifests(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			filePath := row.Cells[0].Value
			fullPath := filepath.Join(tc.projectDir, filePath)
			
			if _, err := os.Stat(fullPath); err != nil {
				return fmt.Errorf("expected Kubernetes manifest %s not found: %w", filePath, err)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theKubernetesManifestsShouldBeValidYAML() error {
	k8sDir := filepath.Join(tc.projectDir, "deployments/k8s")
	
	// Walk through all YAML files in k8s directory
	return filepath.Walk(k8sDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && (strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml")) {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return fmt.Errorf("failed to read K8s manifest %s: %w", path, readErr)
			}
			
			var manifest interface{}
			if yamlErr := yaml.Unmarshal(content, &manifest); yamlErr != nil {
				return fmt.Errorf("K8s manifest %s is not valid YAML: %w", path, yamlErr)
			}
		}
		
		return nil
	})
}

func (tc *TestContext) iHaveGeneratedACompleteWorkspace() error {
	// Generate a complete workspace for testing
	tc.parameters = map[string]string{
		"type":                  "workspace",
		"name":                  "complete-workspace",
		"module":                "github.com/test/complete-workspace",
		"enable_web_api":        "true",
		"enable_cli":            "true",
		"enable_worker":         "true",
		"enable_microservices":  "true",
		"database_type":         "postgres",
		"message_queue":         "redis",
		"enable_docker":         "true",
		"enable_kubernetes":     "true",
	}
	
	return tc.iRunTheGeneratorWith(nil)
}

func (tc *TestContext) iTestTheBuildSystem() error {
	// This will be verified in the next step
	return nil
}

func (tc *TestContext) theMakefileShouldProvideAllExpectedTargets(table *godog.Table) error {
	makefilePath := filepath.Join(tc.projectDir, "Makefile")
	content, err := os.ReadFile(makefilePath)
	if err != nil {
		return fmt.Errorf("failed to read Makefile: %w", err)
	}
	
	makefileStr := string(content)
	
	for _, row := range table.Rows {
		if len(row.Cells) >= 1 {
			target := row.Cells[0].Value
			if !strings.Contains(makefileStr, target+":") {
				return fmt.Errorf("Makefile missing target: %s", target)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) allBuildTargetsShouldExecuteSuccessfully() error {
	targets := []string{"build-all", "test-all", "lint-all", "clean-all"}
	
	for _, target := range targets {
		// Set a timeout for each target
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		
		cmd := exec.CommandContext(ctx, "make", target)
		cmd.Dir = tc.projectDir
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("make target %s failed: %w\nOutput: %s", target, err, output)
		}
	}
	
	return nil
}

func (tc *TestContext) theBuildShouldRespectDependencyOrder() error {
	// Verify that dependencies are built before dependents
	buildScript := filepath.Join(tc.projectDir, "scripts/build-all.sh")
	content, err := os.ReadFile(buildScript)
	if err != nil {
		return fmt.Errorf("failed to read build script: %w", err)
	}
	
	contentStr := string(content)
	
	// Check that pkg/shared is built before other packages
	sharedIndex := strings.Index(contentStr, "pkg/shared")
	modelsIndex := strings.Index(contentStr, "pkg/models")
	
	if sharedIndex == -1 || modelsIndex == -1 {
		return fmt.Errorf("build script missing shared or models packages")
	}
	
	if sharedIndex > modelsIndex {
		return fmt.Errorf("build script should build pkg/shared before pkg/models")
	}
	
	return nil
}

// Helper methods

func (tc *TestContext) verifyCompilation(projectDir string) error {
	// First sync the workspace
	syncCmd := exec.Command("go", "work", "sync")
	syncCmd.Dir = projectDir
	if output, err := syncCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go work sync failed: %w\nOutput: %s", err, output)
	}
	
	// Then build all modules
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (tc *TestContext) verifyFrameworkSpecificImports(projectDir, framework string) error {
	if tc.parameters["enable_web_api"] != "true" {
		return nil // Skip if web API is not enabled
	}
	
	apiMainPath := filepath.Join(projectDir, "cmd/api/main.go")
	content, err := os.ReadFile(apiMainPath)
	if err != nil {
		return fmt.Errorf("failed to read API main.go for %s: %w", framework, err)
	}
	
	expectedImports := map[string]string{
		"gin":   "github.com/gin-gonic/gin",
		"echo":  "github.com/labstack/echo",
		"fiber": "github.com/gofiber/fiber",
		"chi":   "github.com/go-chi/chi",
	}
	
	if expectedImport, exists := expectedImports[framework]; exists {
		if !strings.Contains(string(content), expectedImport) {
			return fmt.Errorf("framework %s missing expected import %s", framework, expectedImport)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyConsistentModuleStructure(projectDir string) error {
	requiredDirs := []string{
		"pkg/shared",
		"pkg/models",
	}
	
	for _, dir := range requiredDirs {
		fullPath := filepath.Join(projectDir, dir)
		if info, err := os.Stat(fullPath); err != nil || !info.IsDir() {
			return fmt.Errorf("required module directory missing: %s", dir)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyDependencyManagement(projectDir string) error {
	goWorkPath := filepath.Join(projectDir, "go.work")
	if _, err := os.Stat(goWorkPath); err != nil {
		return fmt.Errorf("go.work file missing: %w", err)
	}
	
	content, err := os.ReadFile(goWorkPath)
	if err != nil {
		return fmt.Errorf("failed to read go.work: %w", err)
	}
	
	if !strings.Contains(string(content), "use") {
		return fmt.Errorf("go.work should contain 'use' directives")
	}
	
	return nil
}

func (tc *TestContext) verifyDatabaseSpecificStorage(projectDir string) error {
	if tc.parameters["database_type"] == "none" {
		return nil
	}
	
	storageDir := filepath.Join(projectDir, "pkg/storage")
	if info, err := os.Stat(storageDir); err != nil || !info.IsDir() {
		return fmt.Errorf("database storage package missing")
	}
	
	return nil
}

func (tc *TestContext) verifyDatabaseDrivers(projectDir string) error {
	if tc.parameters["database_type"] == "none" {
		return nil
	}
	
	goModPath := filepath.Join(projectDir, "pkg/storage/go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read storage go.mod: %w", err)
	}
	
	database := tc.parameters["database_type"]
	expectedDrivers := map[string]string{
		"postgres": "github.com/lib/pq",
		"mysql":    "github.com/go-sql-driver/mysql",
		"mongodb":  "go.mongodb.org/mongo-driver",
		"sqlite":   "github.com/mattn/go-sqlite3",
	}
	
	if expectedDriver, exists := expectedDrivers[database]; exists {
		if !strings.Contains(string(content), expectedDriver) {
			return fmt.Errorf("database driver %s not found for %s", expectedDriver, database)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyConnectionConfiguration(projectDir string) error {
	configPath := filepath.Join(projectDir, "pkg/shared/config/config.go")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if !strings.Contains(string(content), "Database") {
		return fmt.Errorf("config should contain database configuration")
	}
	
	return nil
}

func (tc *TestContext) verifyMessageQueueEventsPackage(projectDir string) error {
	if tc.parameters["message_queue"] == "none" {
		return nil
	}
	
	eventsDir := filepath.Join(projectDir, "pkg/events")
	if info, err := os.Stat(eventsDir); err != nil || !info.IsDir() {
		return fmt.Errorf("message queue events package missing")
	}
	
	return nil
}

func (tc *TestContext) verifyMQClientLibraries(projectDir string) error {
	if tc.parameters["message_queue"] == "none" {
		return nil
	}
	
	goModPath := filepath.Join(projectDir, "pkg/events/go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read events go.mod: %w", err)
	}
	
	mq := tc.parameters["message_queue"]
	expectedLibraries := map[string]string{
		"redis":    "github.com/go-redis/redis",
		"nats":     "github.com/nats-io/nats.go",
		"kafka":    "github.com/segmentio/kafka-go",
		"rabbitmq": "github.com/streadway/amqp",
	}
	
	if expectedLib, exists := expectedLibraries[mq]; exists {
		if !strings.Contains(string(content), expectedLib) {
			return fmt.Errorf("message queue library %s not found for %s", expectedLib, mq)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyLoggerImplementation(projectDir string) error {
	loggerPath := filepath.Join(projectDir, "pkg/shared/logger/logger.go")
	content, err := os.ReadFile(loggerPath)
	if err != nil {
		return fmt.Errorf("failed to read logger file: %w", err)
	}
	
	logger := tc.parameters["logger_type"]
	expectedPackages := map[string]string{
		"slog":    "log/slog",
		"zap":     "go.uber.org/zap",
		"logrus":  "github.com/sirupsen/logrus",
		"zerolog": "github.com/rs/zerolog",
	}
	
	if expectedPackage, exists := expectedPackages[logger]; exists {
		if !strings.Contains(string(content), expectedPackage) {
			return fmt.Errorf("logger package %s not found for %s", expectedPackage, logger)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyConsistentLoggingInterface(projectDir string) error {
	loggerPath := filepath.Join(projectDir, "pkg/shared/logger/interface.go")
	content, err := os.ReadFile(loggerPath)
	if err != nil {
		return fmt.Errorf("failed to read logger interface: %w", err)
	}
	
	if !strings.Contains(string(content), "Logger interface") {
		return fmt.Errorf("logger interface not found")
	}
	
	return nil
}

// Test runner function
func TestWorkspaceBlueprintATDD(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			tc := NewTestContext(t)
			
			// Register step definitions
			ctx.Step(`^I have the go-starter CLI tool available$`, tc.iHaveTheGoStarterCLIToolAvailable)
			ctx.Step(`^I am in a temporary working directory$`, tc.iAmInATemporaryWorkingDirectory)
			ctx.Step(`^I want to generate a workspace blueprint$`, tc.iWantToGenerateAWorkspaceBlueprint)
			ctx.Step(`^I run the generator with:$`, tc.iRunTheGeneratorWith)
			ctx.Step(`^the generation should succeed$`, tc.theGenerationShouldSucceed)
			ctx.Step(`^the generated workspace should have the Go workspace structure:$`, tc.theGeneratedWorkspaceShouldHaveTheGoWorkspaceStructure)
			ctx.Step(`^all modules should compile successfully$`, tc.allModulesShouldCompileSuccessfully)
			ctx.Step(`^the workspace should sync without errors$`, tc.theWorkspaceShouldSyncWithoutErrors)
			ctx.Step(`^the generated workspace should have minimal structure:$`, tc.theGeneratedWorkspaceShouldHaveMinimalStructure)
			ctx.Step(`^the workspace should not include:$`, tc.theWorkspaceShouldNotInclude)
			ctx.Step(`^I run the generator for each framework:$`, tc.iRunTheGeneratorForEachFramework)
			ctx.Step(`^I run the generator for each database:$`, tc.iRunTheGeneratorForEachDatabase)
			ctx.Step(`^I run the generator for each message queue:$`, tc.iRunTheGeneratorForEachMessageQueue)
			ctx.Step(`^I run the generator for each logger:$`, tc.iRunTheGeneratorForEachLogger)
			ctx.Step(`^each generated workspace should:$`, tc.eachGeneratedWorkspaceShouldWithTable)
			ctx.Step(`^the generated workspace should include Docker files:$`, tc.theGeneratedWorkspaceShouldIncludeDockerFiles)
			ctx.Step(`^the Docker configurations should be valid$`, tc.theDockerConfigurationsShouldBeValid)
			ctx.Step(`^the generated workspace should include Kubernetes manifests:$`, tc.theGeneratedWorkspaceShouldIncludeKubernetesManifests)
			ctx.Step(`^the Kubernetes manifests should be valid YAML$`, tc.theKubernetesManifestsShouldBeValidYAML)
			ctx.Step(`^I have generated a complete workspace$`, tc.iHaveGeneratedACompleteWorkspace)
			ctx.Step(`^I test the build system$`, tc.iTestTheBuildSystem)
			ctx.Step(`^the Makefile should provide all expected targets:$`, tc.theMakefileShouldProvideAllExpectedTargets)
			ctx.Step(`^all build targets should execute successfully$`, tc.allBuildTargetsShouldExecuteSuccessfully)
			ctx.Step(`^the build should respect dependency order$`, tc.theBuildShouldRespectDependencyOrder)
			
			// Cleanup after each scenario
			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				tc.Cleanup()
				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/workspace.feature"},
			TestingT: t,
			Output:   colors.Colored(os.Stdout),
		},
	}
	
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run workspace ATDD tests")
	}
}