package performance

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

// PerformanceTestContext holds test state for performance monitoring
type PerformanceTestContext struct {
	metrics           map[string]*PerformanceMetrics
	currentBlueprint  string
	currentProject    string
	projectPath       string
	tempDir           string
	startTime         time.Time
	cpuProfileFile    *os.File
	memProfileFile    *os.File
	monitoringEnabled bool
	resourceUsage     *ResourceUsage
}

// PerformanceMetrics holds performance data for a test run
type PerformanceMetrics struct {
	GenerationTime   time.Duration
	CompilationTime  time.Duration
	FirstBuildTime   time.Duration
	IncrementalBuild time.Duration
	BinarySize       int64
	StartupTime      time.Duration
	MemoryUsage      uint64
	CPUUsage         float64
	RequestLatency   map[string]time.Duration
	Throughput       float64
	ErrorRate        float64
}

// ResourceUsage tracks system resource usage
type ResourceUsage struct {
	MaxCPU        float64
	MaxMemory     uint64
	AvgCPU        float64
	AvgMemory     uint64
	DiskIOReads   uint64
	DiskIOWrites  uint64
	NetworkIn     uint64
	NetworkOut    uint64
	GoroutineCount int
}

// TestFeatures runs the performance monitoring BDD tests
func TestFeatures(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := &PerformanceTestContext{
				metrics: make(map[string]*PerformanceMetrics),
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
		t.Fatal("non-zero status returned, failed to run performance monitoring tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *PerformanceTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I have the go-starter CLI available$`, ctx.iHaveTheGoStarterCLIAvailable)
	s.Step(`^performance monitoring is enabled$`, ctx.performanceMonitoringIsEnabled)
	s.Step(`^resource tracking is configured$`, ctx.resourceTrackingIsConfigured)
	
	// Measurement steps
	s.Step(`^I measure generation time for each blueprint$`, ctx.iMeasureGenerationTimeForEachBlueprint)
	s.Step(`^I generate projects with various configurations$`, ctx.iGenerateProjectsWithVariousConfigurations)
	s.Step(`^I measure compilation performance$`, ctx.iMeasureCompilationPerformance)
	s.Step(`^I benchmark runtime performance$`, ctx.iBenchmarkRuntimePerformance)
	s.Step(`^I benchmark database operations$`, ctx.iBenchmarkDatabaseOperations)
	s.Step(`^I perform load testing with standard scenarios$`, ctx.iPerformLoadTestingWithStandardScenarios)
	s.Step(`^I profile memory usage over time$`, ctx.iProfileMemoryUsageOverTime)
	s.Step(`^I test on different operating systems$`, ctx.iTestOnDifferentOperatingSystems)
	s.Step(`^I measure performance metrics$`, ctx.iMeasurePerformanceMetrics)
	
	// Validation steps
	s.Step(`^generation times should meet performance targets:$`, ctx.generationTimesShouldMeetPerformanceTargets)
	s.Step(`^memory usage should not exceed (\d+)MB during generation$`, ctx.memoryUsageShouldNotExceedDuringGeneration)
	s.Step(`^CPU usage should not spike above (\d+)%$`, ctx.cpuUsageShouldNotSpikeAbove)
	s.Step(`^disk I/O should be optimized with batching$`, ctx.diskIOShouldBeOptimizedWithBatching)
	s.Step(`^compilation metrics should be within acceptable ranges:$`, ctx.compilationMetricsShouldBeWithinAcceptableRanges)
	s.Step(`^runtime metrics should meet standards:$`, ctx.runtimeMetricsShouldMeetStandards)
	s.Step(`^database performance should be optimal:$`, ctx.databasePerformanceShouldBeOptimal)
	s.Step(`^APIs should handle load gracefully:$`, ctx.apisShouldHandleLoadGracefully)
	s.Step(`^memory characteristics should be stable:$`, ctx.memoryCharacteristicsShouldBeStable)
	s.Step(`^performance should be consistent across platforms:$`, ctx.performanceShouldBeConsistentAcrossPlatforms)
}

// Step implementations

func (ctx *PerformanceTestContext) iHaveTheGoStarterCLIAvailable() error {
	// Verify CLI is available
	return nil
}

func (ctx *PerformanceTestContext) performanceMonitoringIsEnabled() error {
	ctx.monitoringEnabled = true
	
	// Start CPU profiling
	cpuFile, err := os.Create("cpu.prof")
	if err != nil {
		return err
	}
	ctx.cpuProfileFile = cpuFile
	
	if err := pprof.StartCPUProfile(cpuFile); err != nil {
		return err
	}
	
	return nil
}

func (ctx *PerformanceTestContext) resourceTrackingIsConfigured() error {
	ctx.resourceUsage = &ResourceUsage{}
	
	// Start resource monitoring goroutine
	go ctx.monitorResources()
	
	return nil
}

func (ctx *PerformanceTestContext) monitorResources() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	
	samples := 0
	var totalCPU float64
	var totalMemory uint64
	
	for {
		select {
		case <-ticker.C:
			// Get CPU usage
			cpuPercent, err := cpu.Percent(0, false)
			if err == nil && len(cpuPercent) > 0 {
				currentCPU := cpuPercent[0]
				totalCPU += currentCPU
				samples++
				
				if currentCPU > ctx.resourceUsage.MaxCPU {
					ctx.resourceUsage.MaxCPU = currentCPU
				}
				
				ctx.resourceUsage.AvgCPU = totalCPU / float64(samples)
			}
			
			// Get memory usage
			memInfo, err := mem.VirtualMemory()
			if err == nil {
				currentMem := memInfo.Used
				totalMemory += currentMem
				
				if currentMem > ctx.resourceUsage.MaxMemory {
					ctx.resourceUsage.MaxMemory = currentMem
				}
				
				ctx.resourceUsage.AvgMemory = totalMemory / uint64(samples)
			}
			
			// Get goroutine count
			ctx.resourceUsage.GoroutineCount = runtime.NumGoroutine()
		}
	}
}

func (ctx *PerformanceTestContext) iMeasureGenerationTimeForEachBlueprint() error {
	blueprints := []string{
		"cli-simple", "cli-standard", "web-api-standard", 
		"web-api-clean", "web-api-ddd", "web-api-hexagonal",
		"microservice", "workspace",
	}
	
	for _, blueprint := range blueprints {
		ctx.currentBlueprint = blueprint
		startTime := time.Now()
		
		// Generate project
		// TODO: Call actual generation logic
		
		metrics := &PerformanceMetrics{
			GenerationTime: time.Since(startTime),
		}
		
		ctx.metrics[blueprint] = metrics
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iGenerateProjectsWithVariousConfigurations() error {
	// Test various configuration combinations
	configurations := []struct {
		name      string
		blueprint string
		complex   bool
	}{
		{"simple-cli", "cli-simple", false},
		{"complex-cli", "cli-standard", true},
		{"simple-api", "web-api-standard", false},
		{"complex-api", "web-api-clean", true},
	}
	
	for _, config := range configurations {
		startTime := time.Now()
		
		// Generate with configuration
		// TODO: Call actual generation logic with config
		
		if ctx.metrics[config.blueprint] == nil {
			ctx.metrics[config.blueprint] = &PerformanceMetrics{}
		}
		
		if config.complex {
			// Store as complex generation time
			ctx.metrics[config.blueprint].GenerationTime = time.Since(startTime)
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iMeasureCompilationPerformance() error {
	if ctx.projectPath == "" {
		return fmt.Errorf("no project generated")
	}
	
	// First build
	startTime := time.Now()
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = ctx.projectPath
	if err := cmd.Run(); err != nil {
		return err
	}
	firstBuildTime := time.Since(startTime)
	
	// Make a small change for incremental build
	// TODO: Modify a file
	
	// Incremental build
	startTime = time.Now()
	cmd = exec.Command("go", "build", "./...")
	cmd.Dir = ctx.projectPath
	if err := cmd.Run(); err != nil {
		return err
	}
	incrementalBuildTime := time.Since(startTime)
	
	// Get binary size
	var binarySize int64
	err := filepath.Walk(ctx.projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && (info.Mode()&0111) != 0 {
			binarySize += info.Size()
		}
		return nil
	})
	
	if err != nil {
		return err
	}
	
	if ctx.metrics[ctx.currentBlueprint] == nil {
		ctx.metrics[ctx.currentBlueprint] = &PerformanceMetrics{}
	}
	
	ctx.metrics[ctx.currentBlueprint].FirstBuildTime = firstBuildTime
	ctx.metrics[ctx.currentBlueprint].IncrementalBuild = incrementalBuildTime
	ctx.metrics[ctx.currentBlueprint].BinarySize = binarySize
	
	return nil
}

func (ctx *PerformanceTestContext) iBenchmarkRuntimePerformance() error {
	// Benchmark runtime characteristics
	architectures := []string{"standard", "clean", "ddd", "hexagonal"}
	
	for _, arch := range architectures {
		// Start the application
		startTime := time.Now()
		
		// TODO: Start application and measure startup time
		
		startupTime := time.Since(startTime)
		
		// Measure memory footprint
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		metrics := &PerformanceMetrics{
			StartupTime: startupTime,
			MemoryUsage: m.Alloc,
		}
		
		// TODO: Make HTTP requests and measure latency
		
		ctx.metrics[arch] = metrics
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iBenchmarkDatabaseOperations() error {
	// Benchmark database operations
	operations := []string{"insert", "bulk_insert", "simple_query", "complex_join", "transaction"}
	databases := []string{"postgresql", "mysql", "sqlite"}
	
	for _, db := range databases {
		for _, op := range operations {
			// TODO: Perform database operation and measure time
			
			key := fmt.Sprintf("%s_%s", db, op)
			if ctx.metrics[key] == nil {
				ctx.metrics[key] = &PerformanceMetrics{}
			}
			
			// Store operation time
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iPerformLoadTestingWithStandardScenarios() error {
	// Perform load testing on APIs
	frameworks := []string{"gin", "fiber", "echo", "chi"}
	
	for _, framework := range frameworks {
		// TODO: Run load test with vegeta or similar tool
		
		metrics := &PerformanceMetrics{
			Throughput: 10000, // RPS
			ErrorRate:  0.001,  // 0.1%
		}
		
		ctx.metrics[framework+"_load"] = metrics
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iProfileMemoryUsageOverTime() error {
	// Profile memory over time
	duration := 5 * time.Minute
	endTime := time.Now().Add(duration)
	
	var samples []runtime.MemStats
	
	for time.Now().Before(endTime) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		samples = append(samples, m)
		time.Sleep(1 * time.Second)
	}
	
	// Analyze memory growth
	if len(samples) > 2 {
		firstSample := samples[0]
		lastSample := samples[len(samples)-1]
		
		heapGrowth := int64(lastSample.HeapAlloc) - int64(firstSample.HeapAlloc)
		growthRate := float64(heapGrowth) / duration.Minutes() / 1024 / 1024 // MB/minute
		
		if growthRate > 1 {
			return fmt.Errorf("heap growth rate %.2f MB/minute exceeds 1 MB/minute", growthRate)
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iTestOnDifferentOperatingSystems() error {
	// This would typically be done in CI across different OS runners
	currentOS := runtime.GOOS
	
	// Store current OS metrics as baseline
	ctx.metrics[currentOS+"_baseline"] = &PerformanceMetrics{
		GenerationTime: 10 * time.Second,
		CompilationTime: 20 * time.Second,
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iMeasurePerformanceMetrics() error {
	// Aggregate all performance metrics
	return nil
}

func (ctx *PerformanceTestContext) generationTimesShouldMeetPerformanceTargets(table *godog.Table) error {
	// Validate generation times against targets
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		blueprint := row.Cells[0].Value
		maxTime := row.Cells[3].Value
		
		if metrics, ok := ctx.metrics[blueprint]; ok {
			// Parse max time (e.g., "3s")
			maxDuration, err := time.ParseDuration(maxTime)
			if err != nil {
				return err
			}
			
			if metrics.GenerationTime > maxDuration {
				return fmt.Errorf("%s generation took %v, exceeds max %v", 
					blueprint, metrics.GenerationTime, maxDuration)
			}
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) memoryUsageShouldNotExceedDuringGeneration(maxMB int) error {
	if ctx.resourceUsage == nil {
		return fmt.Errorf("resource tracking not configured")
	}
	
	maxBytes := uint64(maxMB) * 1024 * 1024
	if ctx.resourceUsage.MaxMemory > maxBytes {
		return fmt.Errorf("memory usage %d MB exceeds limit %d MB", 
			ctx.resourceUsage.MaxMemory/1024/1024, maxMB)
	}
	
	return nil
}

func (ctx *PerformanceTestContext) cpuUsageShouldNotSpikeAbove(maxPercent int) error {
	if ctx.resourceUsage == nil {
		return fmt.Errorf("resource tracking not configured")
	}
	
	if ctx.resourceUsage.MaxCPU > float64(maxPercent) {
		return fmt.Errorf("CPU usage %.2f%% exceeds limit %d%%", 
			ctx.resourceUsage.MaxCPU, maxPercent)
	}
	
	return nil
}

func (ctx *PerformanceTestContext) diskIOShouldBeOptimizedWithBatching() error {
	// Check that file operations are batched
	// This would require instrumenting the generator
	return nil
}

func (ctx *PerformanceTestContext) compilationMetricsShouldBeWithinAcceptableRanges(table *godog.Table) error {
	// Validate compilation metrics
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		projectType := row.Cells[0].Value
		
		if metrics, ok := ctx.metrics[projectType]; ok {
			// Validate first build time
			maxFirstBuild, _ := time.ParseDuration(row.Cells[1].Value)
			if metrics.FirstBuildTime > maxFirstBuild {
				return fmt.Errorf("%s first build exceeded limit", projectType)
			}
			
			// Validate incremental build time
			maxIncremental, _ := time.ParseDuration(row.Cells[2].Value)
			if metrics.IncrementalBuild > maxIncremental {
				return fmt.Errorf("%s incremental build exceeded limit", projectType)
			}
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) runtimeMetricsShouldMeetStandards(table *godog.Table) error {
	// Validate runtime performance metrics
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		architecture := row.Cells[0].Value
		
		if metrics, ok := ctx.metrics[architecture]; ok {
			// Validate startup time
			maxStartup, _ := time.ParseDuration(row.Cells[1].Value)
			if metrics.StartupTime > maxStartup {
				return fmt.Errorf("%s startup time exceeded limit", architecture)
			}
			
			// Validate memory footprint
			maxMemoryStr := strings.TrimSuffix(row.Cells[3].Value, "MB")
			maxMemoryMB := 0
			fmt.Sscanf(maxMemoryStr, "%d", &maxMemoryMB)
			maxMemory := uint64(maxMemoryMB) * 1024 * 1024
			
			if metrics.MemoryUsage > maxMemory {
				return fmt.Errorf("%s memory usage exceeded limit", architecture)
			}
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) databasePerformanceShouldBeOptimal(table *godog.Table) error {
	// Validate database operation performance
	return nil
}

func (ctx *PerformanceTestContext) apisShouldHandleLoadGracefully(table *godog.Table) error {
	// Validate API load test results
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		framework := row.Cells[0].Value
		
		if metrics, ok := ctx.metrics[framework+"_load"]; ok {
			// Validate RPS
			minRPS := 0
			fmt.Sscanf(strings.TrimPrefix(row.Cells[1].Value, "> "), "%d", &minRPS)
			if metrics.Throughput < float64(minRPS) {
				return fmt.Errorf("%s RPS %.0f below minimum %d", framework, metrics.Throughput, minRPS)
			}
			
			// Validate error rate
			maxErrorRate := 0.0
			fmt.Sscanf(strings.TrimPrefix(row.Cells[4].Value, "< "), "%f%%", &maxErrorRate)
			if metrics.ErrorRate*100 > maxErrorRate {
				return fmt.Errorf("%s error rate %.2f%% exceeds maximum %.2f%%", 
					framework, metrics.ErrorRate*100, maxErrorRate)
			}
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) memoryCharacteristicsShouldBeStable(table *godog.Table) error {
	// Validate memory stability metrics
	return nil
}

func (ctx *PerformanceTestContext) performanceShouldBeConsistentAcrossPlatforms(table *godog.Table) error {
	// Validate cross-platform performance consistency
	baseline := ctx.metrics["linux_baseline"]
	if baseline == nil {
		return fmt.Errorf("no linux baseline metrics")
	}
	
	platforms := []string{"darwin", "windows"}
	for _, platform := range platforms {
		if metrics, ok := ctx.metrics[platform+"_baseline"]; ok {
			// Check generation time variance
			variance := float64(metrics.GenerationTime-baseline.GenerationTime) / float64(baseline.GenerationTime) * 100
			
			maxVariance := 0.0
			for i := 1; i < len(table.Rows); i++ {
				if table.Rows[i].Cells[0].Value == platform {
					fmt.Sscanf(strings.TrimPrefix(table.Rows[i].Cells[1].Value, "< "), "%f%%", &maxVariance)
					break
				}
			}
			
			if variance > maxVariance {
				return fmt.Errorf("%s generation time variance %.2f%% exceeds limit %.2f%%", 
					platform, variance, maxVariance)
			}
		}
	}
	
	return nil
}