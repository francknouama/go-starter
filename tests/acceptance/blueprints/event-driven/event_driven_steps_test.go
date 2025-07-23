package event_driven_test

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
)

// TestContext holds the test execution context for event-driven blueprint testing
type TestContext struct {
	t               *testing.T
	tempDir         string
	projectDir      string
	projectName     string
	generationError error
	parameters      map[string]string
	
	// Event-driven specific fields
	databases       []string
	messageQueues   []string
	frameworks      []string
	loggers         []string
}

// NewTestContext creates a new test context for event-driven blueprint testing
func NewTestContext(t *testing.T) *TestContext {
	return &TestContext{
		t:             t,
		parameters:    make(map[string]string),
		databases:     []string{"postgres", "mysql", "mongodb"},
		messageQueues: []string{"redis", "nats", "kafka", "rabbitmq"},
		frameworks:    []string{"gin", "echo", "fiber", "chi"},
		loggers:       []string{"slog", "zap", "logrus", "zerolog"},
	}
}

// Cleanup removes temporary directories
func (tc *TestContext) Cleanup() {
	if tc.tempDir != "" {
		_ = os.RemoveAll(tc.tempDir)
	}
}

// Event-Driven ATDD Step Definitions

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
	tc.tempDir, err = os.MkdirTemp("", "event-driven-test-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	
	if err := os.Chdir(tc.tempDir); err != nil {
		return fmt.Errorf("failed to change to temp directory: %w", err)
	}
	
	return nil
}

func (tc *TestContext) iWantToGenerateAnEventDrivenBlueprint() error {
	tc.parameters["type"] = "event-driven"
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
		case "framework":
			args = append(args, "--framework", value)
		case "logger":
			args = append(args, "--logger", value)
		case "database_type":
			args = append(args, "--database-type", value)
		case "message_queue":
			args = append(args, "--message-queue", value)
		case "architecture":
			args = append(args, "--architecture", value)
		case "services":
			args = append(args, "--services", value)
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

func (tc *TestContext) theGeneratedProjectShouldHaveTheEventDrivenStructure(table *godog.Table) error {
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
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = tc.projectDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w\nOutput: %s", err, output)
	}
	
	// Build all modules
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = tc.projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %w\nOutput: %s", err, output)
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
		tc.parameters["name"] = fmt.Sprintf("test-event-driven-db-%s", database)
		
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
		tc.parameters["name"] = fmt.Sprintf("test-event-driven-mq-%s", mq)
		
		if err := tc.iRunTheGeneratorWith(nil); err != nil {
			return err
		}
		
		if err := tc.theGenerationShouldSucceed(); err != nil {
			return fmt.Errorf("generation failed for message queue %s: %w", mq, err)
		}
	}
	
	return nil
}

func (tc *TestContext) eachGeneratedProjectShouldWithTable(table *godog.Table) error {
	for _, database := range tc.databases {
		projectDir := filepath.Join(tc.tempDir, fmt.Sprintf("test-event-driven-db-%s", database))
		
		for _, row := range table.Rows {
			if len(row.Cells) >= 1 {
				validation := row.Cells[0].Value
				
				switch validation {
				case "compile successfully":
					if err := tc.verifyCompilation(projectDir); err != nil {
						return fmt.Errorf("compilation failed for %s: %w", database, err)
					}
					
				case "include database-specific event store":
					if err := tc.verifyEventStore(projectDir, database); err != nil {
						return err
					}
					
				case "contain appropriate database drivers":
					if err := tc.verifyDatabaseDrivers(projectDir, database); err != nil {
						return err
					}
					
				case "have correct connection configuration":
					if err := tc.verifyConnectionConfiguration(projectDir); err != nil {
						return err
					}
				
				case "include message queue event bus":
					if err := tc.verifyEventBus(projectDir); err != nil {
						return err
					}
					
				case "contain appropriate MQ client libraries":
					if err := tc.verifyMQClientLibraries(projectDir); err != nil {
						return err
					}
					
				case "have correct event publishing configuration":
					if err := tc.verifyEventPublishingConfiguration(projectDir); err != nil {
						return err
					}
				}
			}
		}
	}
	
	return nil
}

func (tc *TestContext) iHaveGeneratedAnEventDrivenProject() error {
	// Generate a complete event-driven project for testing
	tc.parameters = map[string]string{
		"type":           "event-driven",
		"name":           "complete-event-driven",
		"module":         "github.com/test/complete-event-driven",
		"framework":      "gin",
		"logger":         "slog",
		"database_type":  "postgres",
		"message_queue":  "redis",
	}
	
	return tc.iRunTheGeneratorWith(nil)
}

func (tc *TestContext) iExamineTheCQRSImplementation() error {
	// This will be verified in the next step
	return nil
}

func (tc *TestContext) theProjectShouldInclude(table *godog.Table) error {
	for _, row := range table.Rows {
		if len(row.Cells) >= 2 {
			component := row.Cells[0].Value
			path := row.Cells[1].Value
			
			fullPath := filepath.Join(tc.projectDir, path)
			if _, err := os.Stat(fullPath); err != nil {
				return fmt.Errorf("expected component %s at path %s not found: %w", component, path, err)
			}
		}
	}
	
	return nil
}

func (tc *TestContext) theCommandBusShouldHandleCommandDispatching() error {
	commandBusPath := filepath.Join(tc.projectDir, "internal/cqrs/command_bus.go")
	content, err := os.ReadFile(commandBusPath)
	if err != nil {
		return fmt.Errorf("failed to read command bus file: %w", err)
	}
	
	requiredFunctions := []string{"Dispatch", "RegisterHandler", "AddMiddleware"}
	for _, function := range requiredFunctions {
		if !strings.Contains(string(content), function) {
			return fmt.Errorf("command bus missing required function: %s", function)
		}
	}
	
	return nil
}

func (tc *TestContext) theQueryBusShouldHandleQueryExecution() error {
	queryBusPath := filepath.Join(tc.projectDir, "internal/cqrs/query_bus.go")
	content, err := os.ReadFile(queryBusPath)
	if err != nil {
		return fmt.Errorf("failed to read query bus file: %w", err)
	}
	
	requiredFunctions := []string{"Execute", "RegisterHandler", "AddMiddleware"}
	for _, function := range requiredFunctions {
		if !strings.Contains(string(content), function) {
			return fmt.Errorf("query bus missing required function: %s", function)
		}
	}
	
	return nil
}

func (tc *TestContext) commandsAndQueriesShouldBeProperlySeparated() error {
	// Verify command and query separation
	commandPath := filepath.Join(tc.projectDir, "internal/cqrs/command.go")
	queryPath := filepath.Join(tc.projectDir, "internal/cqrs/query.go")
	
	// Both files should exist
	if _, err := os.Stat(commandPath); err != nil {
		return fmt.Errorf("command interface not found: %w", err)
	}
	
	if _, err := os.Stat(queryPath); err != nil {
		return fmt.Errorf("query interface not found: %w", err)
	}
	
	// Verify they define different interfaces
	commandContent, err := os.ReadFile(commandPath)
	if err != nil {
		return fmt.Errorf("failed to read command file: %w", err)
	}
	
	queryContent, err := os.ReadFile(queryPath)
	if err != nil {
		return fmt.Errorf("failed to read query file: %w", err)
	}
	
	if !strings.Contains(string(commandContent), "Command interface") {
		return fmt.Errorf("command interface not properly defined")
	}
	
	if !strings.Contains(string(queryContent), "Query interface") {
		return fmt.Errorf("query interface not properly defined")
	}
	
	return nil
}

func (tc *TestContext) iExamineTheEventSourcingImplementation() error {
	return nil
}

func (tc *TestContext) aggregatesShouldApplyEventsCorrectly() error {
	aggregatePath := filepath.Join(tc.projectDir, "internal/domain/aggregate.go")
	content, err := os.ReadFile(aggregatePath)
	if err != nil {
		return fmt.Errorf("failed to read aggregate file: %w", err)
	}
	
	requiredMethods := []string{"ApplyEvent", "GetUncommittedEvents", "MarkEventsAsCommitted"}
	for _, method := range requiredMethods {
		if !strings.Contains(string(content), method) {
			return fmt.Errorf("aggregate missing required method: %s", method)
		}
	}
	
	return nil
}

func (tc *TestContext) eventsShouldBePersistedInTheEventStore() error {
	eventStorePath := filepath.Join(tc.projectDir, "internal/eventstore/store.go")
	content, err := os.ReadFile(eventStorePath)
	if err != nil {
		return fmt.Errorf("failed to read event store file: %w", err)
	}
	
	requiredMethods := []string{"SaveEvents", "GetEvents", "GetEventsFromSnapshot"}
	for _, method := range requiredMethods {
		if !strings.Contains(string(content), method) {
			return fmt.Errorf("event store missing required method: %s", method)
		}
	}
	
	return nil
}

func (tc *TestContext) snapshotsShouldBeSupportedForPerformance() error {
	snapshotPath := filepath.Join(tc.projectDir, "internal/eventstore/snapshots.go")
	if _, err := os.Stat(snapshotPath); err != nil {
		return fmt.Errorf("snapshot support not found: %w", err)
	}
	
	content, err := os.ReadFile(snapshotPath)
	if err != nil {
		return fmt.Errorf("failed to read snapshot file: %w", err)
	}
	
	requiredMethods := []string{"SaveSnapshot", "GetLatestSnapshot"}
	for _, method := range requiredMethods {
		if !strings.Contains(string(content), method) {
			return fmt.Errorf("snapshot store missing required method: %s", method)
		}
	}
	
	return nil
}

// Helper methods

func (tc *TestContext) verifyCompilation(projectDir string) error {
	// Initialize go modules
	modCmd := exec.Command("go", "mod", "tidy")
	modCmd.Dir = projectDir
	if output, err := modCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w\nOutput: %s", err, output)
	}
	
	// Build all modules
	buildCmd := exec.Command("go", "build", "./...")
	buildCmd.Dir = projectDir
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("compilation failed: %w\nOutput: %s", err, output)
	}
	
	return nil
}

func (tc *TestContext) verifyEventStore(projectDir, database string) error {
	eventStorePath := filepath.Join(projectDir, "internal/eventstore")
	if info, err := os.Stat(eventStorePath); err != nil || !info.IsDir() {
		return fmt.Errorf("event store directory missing")
	}
	
	// Verify database-specific implementation
	dbSpecificFile := filepath.Join(eventStorePath, fmt.Sprintf("%s.go", database))
	if _, err := os.Stat(dbSpecificFile); err != nil {
		return fmt.Errorf("database-specific event store implementation missing for %s", database)
	}
	
	return nil
}

func (tc *TestContext) verifyDatabaseDrivers(projectDir, database string) error {
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
	}
	
	expectedDrivers := map[string]string{
		"postgres": "github.com/lib/pq",
		"mysql":    "github.com/go-sql-driver/mysql",
		"mongodb":  "go.mongodb.org/mongo-driver",
	}
	
	if expectedDriver, exists := expectedDrivers[database]; exists {
		if !strings.Contains(string(content), expectedDriver) {
			return fmt.Errorf("database driver %s not found for %s", expectedDriver, database)
		}
	}
	
	return nil
}

func (tc *TestContext) verifyConnectionConfiguration(projectDir string) error {
	configPath := filepath.Join(projectDir, "internal/config/config.go")
	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	
	if !strings.Contains(string(content), "Database") {
		return fmt.Errorf("config should contain database configuration")
	}
	
	return nil
}

func (tc *TestContext) verifyEventBus(projectDir string) error {
	eventBusPath := filepath.Join(projectDir, "internal/eventbus")
	if info, err := os.Stat(eventBusPath); err != nil || !info.IsDir() {
		return fmt.Errorf("event bus directory missing")
	}
	
	return nil
}

func (tc *TestContext) verifyMQClientLibraries(projectDir string) error {
	goModPath := filepath.Join(projectDir, "go.mod")
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("failed to read go.mod: %w", err)
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

func (tc *TestContext) verifyEventPublishingConfiguration(projectDir string) error {
	eventBusPath := filepath.Join(projectDir, "internal/eventbus/eventbus.go")
	content, err := os.ReadFile(eventBusPath)
	if err != nil {
		return fmt.Errorf("failed to read event bus file: %w", err)
	}
	
	if !strings.Contains(string(content), "Publish") {
		return fmt.Errorf("event bus should contain publish method")
	}
	
	return nil
}

// Test runner function
func TestEventDrivenBlueprintATDD(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(ctx *godog.ScenarioContext) {
			tc := NewTestContext(t)
			
			// Register step definitions
			ctx.Step(`^I have the go-starter CLI tool available$`, tc.iHaveTheGoStarterCLIToolAvailable)
			ctx.Step(`^I am in a temporary working directory$`, tc.iAmInATemporaryWorkingDirectory)
			ctx.Step(`^I want to generate an event-driven blueprint$`, tc.iWantToGenerateAnEventDrivenBlueprint)
			ctx.Step(`^I run the generator with:$`, tc.iRunTheGeneratorWith)
			ctx.Step(`^the generation should succeed$`, tc.theGenerationShouldSucceed)
			ctx.Step(`^the generated project should have the event-driven structure:$`, tc.theGeneratedProjectShouldHaveTheEventDrivenStructure)
			ctx.Step(`^all modules should compile successfully$`, tc.allModulesShouldCompileSuccessfully)
			ctx.Step(`^I run the generator for each database:$`, tc.iRunTheGeneratorForEachDatabase)
			ctx.Step(`^I run the generator for each message queue:$`, tc.iRunTheGeneratorForEachMessageQueue)
			ctx.Step(`^each generated project should:$`, tc.eachGeneratedProjectShouldWithTable)
			ctx.Step(`^I have generated an event-driven project$`, tc.iHaveGeneratedAnEventDrivenProject)
			ctx.Step(`^I examine the CQRS implementation$`, tc.iExamineTheCQRSImplementation)
			ctx.Step(`^the project should include:$`, tc.theProjectShouldInclude)
			ctx.Step(`^the command bus should handle command dispatching$`, tc.theCommandBusShouldHandleCommandDispatching)
			ctx.Step(`^the query bus should handle query execution$`, tc.theQueryBusShouldHandleQueryExecution)
			ctx.Step(`^commands and queries should be properly separated$`, tc.commandsAndQueriesShouldBeProperlySeparated)
			ctx.Step(`^I examine the Event Sourcing implementation$`, tc.iExamineTheEventSourcingImplementation)
			ctx.Step(`^aggregates should apply events correctly$`, tc.aggregatesShouldApplyEventsCorrectly)
			ctx.Step(`^events should be persisted in the event store$`, tc.eventsShouldBePersistedInTheEventStore)
			ctx.Step(`^snapshots should be supported for performance$`, tc.snapshotsShouldBeSupportedForPerformance)
			
			// Cleanup after each scenario
			ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				tc.Cleanup()
				return ctx, nil
			})
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features/event-driven.feature"},
			TestingT: t,
			Output:   colors.Colored(os.Stdout),
		},
	}
	
	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run event-driven ATDD tests")
	}
}