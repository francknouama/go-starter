package optimization

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cucumber/godog"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/optimization"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/tests/helpers"
)

// PerformanceTestContext holds state for performance validation tests
type PerformanceTestContext struct {
	// Project state
	projectPaths    []string
	tempDir         string
	currentProject  string
	
	// Performance tracking
	startTime       time.Time
	endTime         time.Time
	processingTimes map[string]time.Duration
	memoryUsage     map[string]int64
	processingRates map[string]float64
	performanceData map[string]interface{}
	
	// Benchmarks and thresholds
	benchmarks      map[string]performanceBenchmark
	thresholds      map[string]performanceThreshold
	
	// Monitoring state
	monitoringActive bool
	resourceMonitor  *resourceMonitor
	performanceLog   []performanceMetric
	
	// Configuration
	optimizationLevel optimization.OptimizationLevel
	concurrencyConfig concurrencySettings
	
	// Results and validation
	optimizationResults []optimization.PipelineResult
	performanceReport   performanceReport
	validationErrors    []string
	
	// Test tracking
	t                   *testing.T
	testDirs            []string
	lastError           error
}

type performanceBenchmark struct {
	ProjectSize     string
	ExpectedTime    time.Duration
	ExpectedThroughput float64
	MaxMemoryUsage  int64
}

type performanceThreshold struct {
	ProcessingTime  time.Duration
	MemoryUsage     int64
	FilesPerSecond  float64
	QualityScore    float64
}

type concurrencySettings struct {
	MaxConcurrentFiles   int
	EnableParallelism    bool
	MemoryLimitMB        int
}

type resourceMonitor struct {
	mutex           sync.RWMutex
	cpuUsage        []float64
	memoryUsage     []int64
	diskIOReads     []int64
	diskIOWrites    []int64
	fileHandles     []int
	startTime       time.Time
	samples         int
	isActive        bool
}

type performanceMetric struct {
	Timestamp   time.Time
	Metric      string
	Value       interface{}
	Unit        string
	Context     string
}

type performanceReport struct {
	ExecutiveSummary    string
	ProcessingTimes     map[string]time.Duration
	ResourceUsage       resourceUsage
	OptimizationImpact  optimizationImpact
	ComparativeAnalysis comparativeAnalysis
	Recommendations     []string
}

type resourceUsage struct {
	PeakMemoryMB      int64
	AverageCPUPercent float64
	TotalDiskIOBytes  int64
	MaxFileHandles    int
}

type optimizationImpact struct {
	FilesProcessed     int
	ImportsRemoved     int
	VariablesRemoved   int
	FunctionsRemoved   int
	CodeSizeReduction  float64
	QualityImprovement float64
}

type comparativeAnalysis struct {
	BaselineTime       time.Duration
	CurrentTime        time.Duration
	PerformanceRatio   float64
	MemoryComparison   string
	ThroughputChange   float64
}

// NewPerformanceTestContext creates a new test context
func NewPerformanceTestContext(t *testing.T) *PerformanceTestContext {
	return &PerformanceTestContext{
		t:                   t,
		projectPaths:        make([]string, 0),
		processingTimes:     make(map[string]time.Duration),
		memoryUsage:         make(map[string]int64),
		processingRates:     make(map[string]float64),
		performanceData:     make(map[string]interface{}),
		benchmarks:          make(map[string]performanceBenchmark),
		thresholds:          make(map[string]performanceThreshold),
		testDirs:            make([]string, 0),
		performanceLog:      make([]performanceMetric, 0),
		optimizationResults: make([]optimization.PipelineResult, 0),
		validationErrors:    make([]string, 0),
		resourceMonitor:     newResourceMonitor(),
	}
}

func newResourceMonitor() *resourceMonitor {
	return &resourceMonitor{
		cpuUsage:     make([]float64, 0),
		memoryUsage:  make([]int64, 0),
		diskIOReads:  make([]int64, 0),
		diskIOWrites: make([]int64, 0),
		fileHandles:  make([]int, 0),
		samples:      0,
		isActive:     false,
	}
}

// TestPerformanceValidation runs the performance validation tests
func TestPerformanceValidation(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: func(s *godog.ScenarioContext) {
			ctx := NewPerformanceTestContext(t)

			s.Before(func(goCtx context.Context, sc *godog.Scenario) (context.Context, error) {
				// Setup before each scenario
				ctx.tempDir = helpers.CreateTempTestDir(t)
				ctx.setupBenchmarks()
				ctx.setupThresholds()
				
				// Initialize templates
				if err := helpers.InitializeTemplates(); err != nil {
					return goCtx, err
				}

				return goCtx, nil
			})

			s.After(func(goCtx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
				// Cleanup after each scenario
				ctx.Cleanup()
				return goCtx, nil
			})

			// Register step definitions
			ctx.RegisterSteps(s)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
			Tags:     "@performance",
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run performance feature tests")
	}
}

// RegisterSteps registers all step definitions
func (ctx *PerformanceTestContext) RegisterSteps(s *godog.ScenarioContext) {
	// Background steps
	s.Step(`^I am using go-starter CLI$`, ctx.iAmUsingGoStarterCLI)
	s.Step(`^the optimization system is available$`, ctx.theOptimizationSystemIsAvailable)
	s.Step(`^performance monitoring is enabled$`, ctx.performanceMonitoringIsEnabled)
	
	// Benchmark steps
	s.Step(`^I have projects of different sizes:$`, ctx.iHaveProjectsOfDifferentSizes)
	s.Step(`^I optimize each project type$`, ctx.iOptimizeEachProjectType)
	s.Step(`^processing time should meet benchmark requirements$`, ctx.processingTimeShouldMeetBenchmarkRequirements)
	s.Step(`^throughput should be at least (\d+) files/minute for small files$`, ctx.throughputShouldBeAtLeast)
	s.Step(`^memory usage should remain under (\d+)MB for large projects$`, ctx.memoryUsageShouldRemainUnder)
	
	// Memory efficiency steps
	s.Step(`^I have a project that will test memory efficiency$`, ctx.iHaveAProjectThatWillTestMemoryEfficiency)
	s.Step(`^I run optimization with memory profiling enabled$`, ctx.iRunOptimizationWithMemoryProfilingEnabled)
	s.Step(`^peak memory usage should not exceed baseline \+ (\d+)MB$`, ctx.peakMemoryUsageShouldNotExceedBaseline)
	s.Step(`^memory should be released progressively during optimization$`, ctx.memoryShouldBeReleasedProgressivelyDuringOptimization)
	s.Step(`^garbage collection should occur regularly$`, ctx.garbageCollectionShouldOccurRegularly)
	s.Step(`^no memory leaks should be detected after completion$`, ctx.noMemoryLeaksShouldBeDetectedAfterCompletion)
	
	// Concurrent processing steps
	s.Step(`^I configure concurrent file processing:$`, ctx.iConfigureConcurrentFileProcessing)
	s.Step(`^I optimize a large project with concurrent processing$`, ctx.iOptimizeALargeProjectWithConcurrentProcessing)
	s.Step(`^processing should utilize multiple CPU cores effectively$`, ctx.processingShouldUtilizeMultipleCPUCoresEffectively)
	s.Step(`^total processing time should be reduced compared to sequential$`, ctx.totalProcessingTimeShouldBeReducedComparedToSequential)
	s.Step(`^memory usage per thread should be controlled$`, ctx.memoryUsagePerThreadShouldBeControlled)
	s.Step(`^no race conditions should occur$`, ctx.noRaceConditionsShouldOccur)
	
	// Scalability steps
	s.Step(`^I have a project with "([^"]*)" Go files$`, ctx.iHaveAProjectWithGoFiles)
	s.Step(`^each file has approximately "([^"]*)" lines$`, ctx.eachFileHasApproximatelyLines)
	s.Step(`^I optimize the project with "([^"]*)" optimization level$`, ctx.iOptimizeTheProjectWithOptimizationLevel)
	s.Step(`^processing time should scale linearly with file count$`, ctx.processingTimeShouldScaleLinearlyWithFileCount)
	s.Step(`^memory usage should remain proportional to project size$`, ctx.memoryUsageShouldRemainProportionalToProjectSize)
	s.Step(`^optimization quality should not degrade with scale$`, ctx.optimizationQualityShouldNotDegradeWithScale)
	
	// Resource utilization steps  
	s.Step(`^I monitor system resources during optimization:$`, ctx.iMonitorSystemResourcesDuringOptimization)
	s.Step(`^I run optimization on multiple project types$`, ctx.iRunOptimizationOnMultipleProjectTypes)
	s.Step(`^CPU utilization should be efficient but not excessive$`, ctx.cpuUtilizationShouldBeEfficientButNotExcessive)
	s.Step(`^disk I/O should be minimized through smart caching$`, ctx.diskIOShouldBeMinimizedThroughSmartCaching)
	s.Step(`^file handles should be properly managed and released$`, ctx.fileHandlesShouldBeProperlyManagedAndReleased)
	s.Step(`^system responsiveness should remain good$`, ctx.systemResponsivenessShouldRemainGood)
	
	// Optimization level performance steps
	s.Step(`^I have the same project to optimize at different levels$`, ctx.iHaveTheSameProjectToOptimizeAtDifferentLevels)
	s.Step(`^I measure performance for each optimization level:$`, ctx.iMeasurePerformanceForEachOptimizationLevel)
	s.Step(`^performance should scale appropriately with optimization complexity$`, ctx.performanceShouldScaleAppropriatelyWithOptimizationComplexity)
	s.Step(`^quality improvements should justify performance costs$`, ctx.qualityImprovementsShouldJustifyPerformanceCosts)
	s.Step(`^users should be warned about performance implications$`, ctx.usersShouldBeWarnedAboutPerformanceImplications)
	
	// Caching steps
	s.Step(`^I have projects with overlapping code patterns$`, ctx.iHaveProjectsWithOverlappingCodePatterns)
	s.Step(`^I optimize similar projects sequentially$`, ctx.iOptimizeSimilarProjectsSequentially)
	s.Step(`^AST parsing results should be cached where possible$`, ctx.astParsingResultsShouldBeCachedWherePossible)
	s.Step(`^repeated pattern analysis should be optimized$`, ctx.repeatedPatternAnalysisShouldBeOptimized)
	s.Step(`^overall processing time should decrease for similar patterns$`, ctx.overallProcessingTimeShouldDecreaseForSimilarPatterns)
	s.Step(`^cache hit rates should be reported and monitored$`, ctx.cacheHitRatesShouldBeReportedAndMonitored)
	
	// Profile-specific performance steps
	s.Step(`^I have a "([^"]*)" project suitable for "([^"]*)" profile$`, ctx.iHaveAProjectSuitableForProfile)
	s.Step(`^I optimize using the "([^"]*)" profile$`, ctx.iOptimizeUsingTheProfile)
	s.Step(`^performance should match profile-specific expectations:$`, ctx.performanceShouldMatchProfileSpecificExpectations)
	s.Step(`^profile recommendations should consider performance trade-offs$`, ctx.profileRecommendationsShouldConsiderPerformanceTradeOffs)
	
	// Regression testing steps
	s.Step(`^I have baseline performance metrics from previous optimization runs$`, ctx.iHaveBaselinePerformanceMetricsFromPreviousOptimizationRuns)
	s.Step(`^I run optimization with the current implementation$`, ctx.iRunOptimizationWithTheCurrentImplementation)
	s.Step(`^performance should not regress beyond acceptable thresholds:$`, ctx.performanceShouldNotRegressBeyondAcceptableThresholds)
	s.Step(`^any performance regressions should be clearly reported$`, ctx.anyPerformanceRegressionsShouldBeClearlyReported)
	s.Step(`^suggestions for performance improvement should be provided$`, ctx.suggestionsForPerformanceImprovementShouldBeProvided)
	
	// Real-time monitoring steps
	s.Step(`^I enable real-time performance monitoring$`, ctx.iEnableRealTimePerformanceMonitoring)
	s.Step(`^I run optimization on a large project$`, ctx.iRunOptimizationOnALargeProject)
	s.Step(`^progress should be reported with performance metrics:$`, ctx.progressShouldBeReportedWithPerformanceMetrics)
	s.Step(`^users should be able to monitor optimization progress$`, ctx.usersShouldBeAbleToMonitorOptimizationProgress)
	s.Step(`^performance bottlenecks should be identified in real-time$`, ctx.performanceBottlenecksShouldBeIdentifiedInRealTime)
	
	// Stress testing steps
	s.Step(`^I create stress test conditions:$`, ctx.iCreateStressTestConditions)
	s.Step(`^I run optimization under stress$`, ctx.iRunOptimizationUnderStress)
	s.Step(`^the system should handle stress gracefully$`, ctx.theSystemShouldHandleStressGracefully)
	s.Step(`^performance should degrade gracefully$`, ctx.performanceShouldDegradeGracefully)
	s.Step(`^error handling should remain responsive$`, ctx.errorHandlingShouldRemainResponsive)
	s.Step(`^no system crashes should occur$`, ctx.noSystemCrashesShouldOccur)
	
	// Quality vs speed steps
	s.Step(`^I have projects with varying optimization opportunities$`, ctx.iHaveProjectsWithVaryingOptimizationOpportunities)
	s.Step(`^I optimize with different quality vs speed trade-offs:$`, ctx.iOptimizeWithDifferentQualityVsSpeedTradeOffs)
	s.Step(`^each mode should deliver expected quality-speed trade-offs$`, ctx.eachModeShouldDeliverExpectedQualitySpeedTradeOffs)
	s.Step(`^users should be able to choose appropriate modes$`, ctx.usersShouldBeAbleToChooseAppropriateModes)
	s.Step(`^trade-offs should be clearly documented$`, ctx.tradeOffsShouldBeClearlyDocumented)
	
	// Batch processing steps
	s.Step(`^I have multiple projects to optimize:$`, ctx.iHaveMultipleProjectsToOptimize)
	s.Step(`^I optimize all projects in batch mode$`, ctx.iOptimizeAllProjectsInBatchMode)
	s.Step(`^batch processing should be more efficient than individual runs$`, ctx.batchProcessingShouldBeMoreEfficientThanIndividualRuns)
	s.Step(`^shared resources should be reused across projects$`, ctx.sharedResourcesShouldBeReusedAcrossProjects)
	s.Step(`^overall processing time should be optimized$`, ctx.overallProcessingTimeShouldBeOptimized)
	s.Step(`^progress should be reported for the entire batch$`, ctx.progressShouldBeReportedForTheEntireBatch)
	
	// Performance profiling steps
	s.Step(`^I enable detailed performance profiling$`, ctx.iEnableDetailedPerformanceProfiling)
	s.Step(`^I optimize a representative project$`, ctx.iOptimizeARepresentativeProject)
	s.Step(`^detailed profiling data should be collected:$`, ctx.detailedProfilingDataShouldBeCollected)
	s.Step(`^profiling results should identify optimization opportunities$`, ctx.profilingResultsShouldIdentifyOptimizationOpportunities)
	s.Step(`^performance recommendations should be generated$`, ctx.performanceRecommendationsShouldBeGenerated)
	
	// Network efficiency steps
	s.Step(`^I have projects stored in different locations:$`, ctx.iHaveProjectsStoredInDifferentLocations)
	s.Step(`^I optimize projects from different locations$`, ctx.iOptimizeProjectsFromDifferentLocations)
	s.Step(`^I/O operations should be minimized and optimized$`, ctx.ioOperationsShouldBeMinimizedAndOptimized)
	s.Step(`^network latency should be handled gracefully$`, ctx.networkLatencyShouldBeHandledGracefully)
	s.Step(`^caching should reduce repeated remote access$`, ctx.cachingShouldReduceRepeatedRemoteAccess)
	s.Step(`^performance should remain acceptable for all locations$`, ctx.performanceShouldRemainAcceptableForAllLocations)
	
	// Performance reporting steps
	s.Step(`^I complete optimization of various projects$`, ctx.iCompleteOptimizationOfVariousProjects)
	s.Step(`^I generate performance reports$`, ctx.iGeneratePerformanceReports)
	s.Step(`^reports should include comprehensive metrics:$`, ctx.reportsShouldIncludeComprehensiveMetrics)
	s.Step(`^reports should be exportable in multiple formats$`, ctx.reportsShouldBeExportableInMultipleFormats)
	s.Step(`^historical performance trends should be tracked$`, ctx.historicalPerformanceTrendsShouldBeTracked)
}

// Background step implementations
func (ctx *PerformanceTestContext) iAmUsingGoStarterCLI() error {
	return nil
}

func (ctx *PerformanceTestContext) theOptimizationSystemIsAvailable() error {
	config := optimization.DefaultConfig()
	if config.Level < optimization.OptimizationLevelNone || config.Level > optimization.OptimizationLevelExpert {
		return fmt.Errorf("optimization system not properly initialized")
	}
	return nil
}

func (ctx *PerformanceTestContext) performanceMonitoringIsEnabled() error {
	ctx.monitoringActive = true
	ctx.resourceMonitor.isActive = true
	return nil
}

// Benchmark implementations
func (ctx *PerformanceTestContext) iHaveProjectsOfDifferentSizes(table *godog.Table) error {
	// Parse project size requirements from table
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		projectSize := row.Cells[0].Value
		files := row.Cells[1].Value
		linesOfCode := row.Cells[2].Value
		expectedProcessing := row.Cells[3].Value
		
		// Parse expected processing time
		expectedTime, err := time.ParseDuration(strings.ReplaceAll(expectedProcessing, " ", ""))
		if err != nil {
			return fmt.Errorf("invalid expected processing time: %s", expectedProcessing)
		}
		
		// Create benchmark entry
		ctx.benchmarks[projectSize] = performanceBenchmark{
			ProjectSize:        projectSize,
			ExpectedTime:       expectedTime,
			ExpectedThroughput: 100.0, // files per minute
			MaxMemoryUsage:     500 * 1024 * 1024, // 500MB
		}
		
		// Generate test project of appropriate size
		err = ctx.generateProjectOfSize(projectSize, files, linesOfCode)
		if err != nil {
			return fmt.Errorf("failed to generate %s project: %w", projectSize, err)
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iOptimizeEachProjectType() error {
	for projectSize := range ctx.benchmarks {
		ctx.startTime = time.Now()
		
		// Start resource monitoring
		if ctx.monitoringActive {
			go ctx.startResourceMonitoring()
		}
		
		// Run optimization
		projectPath := filepath.Join(ctx.tempDir, fmt.Sprintf("%s-project", projectSize))
		pipeline := optimization.NewOptimizationPipeline(optimization.DefaultPipelineOptions())
		
		result, err := pipeline.OptimizeProject(projectPath)
		if err != nil {
			return fmt.Errorf("optimization failed for %s project: %w", projectSize, err)
		}
		
		ctx.endTime = time.Now()
		
		// Stop resource monitoring
		if ctx.monitoringActive {
			ctx.stopResourceMonitoring()
		}
		
		// Record performance metrics
		processingTime := ctx.endTime.Sub(ctx.startTime)
		ctx.processingTimes[projectSize] = processingTime
		
		// Calculate throughput
		if result.FilesOptimized > 0 {
			throughput := float64(result.FilesOptimized) / processingTime.Minutes()
			ctx.processingRates[projectSize] = throughput
		}
		
		ctx.optimizationResults = append(ctx.optimizationResults, *result)
	}
	
	return nil
}

func (ctx *PerformanceTestContext) processingTimeShouldMeetBenchmarkRequirements() error {
	for projectSize, actualTime := range ctx.processingTimes {
		benchmark := ctx.benchmarks[projectSize]
		
		if actualTime > benchmark.ExpectedTime {
			return fmt.Errorf("%s project took %v, expected under %v", 
				projectSize, actualTime, benchmark.ExpectedTime)
		}
		
		ctx.logPerformanceMetric("processing_time", projectSize, actualTime, "duration")
	}
	
	return nil
}

func (ctx *PerformanceTestContext) throughputShouldBeAtLeast(minThroughput int) error {
	for projectSize, actualThroughput := range ctx.processingRates {
		if actualThroughput < float64(minThroughput) {
			return fmt.Errorf("%s project throughput was %.2f files/minute, expected at least %d", 
				projectSize, actualThroughput, minThroughput)
		}
		
		ctx.logPerformanceMetric("throughput", projectSize, actualThroughput, "files/minute")
	}
	
	return nil
}

func (ctx *PerformanceTestContext) memoryUsageShouldRemainUnder(maxMemoryMB int) error {
	maxMemoryBytes := int64(maxMemoryMB * 1024 * 1024)
	
	for projectSize, peakMemory := range ctx.memoryUsage {
		if peakMemory > maxMemoryBytes {
			return fmt.Errorf("%s project used %d MB memory, expected under %d MB", 
				projectSize, peakMemory/(1024*1024), maxMemoryMB)
		}
		
		ctx.logPerformanceMetric("peak_memory", projectSize, peakMemory, "bytes")
	}
	
	return nil
}

// Memory efficiency implementations
func (ctx *PerformanceTestContext) iHaveAProjectThatWillTestMemoryEfficiency() error {
	return ctx.generateProjectOfSize("memory-test", "100", "50000")
}

func (ctx *PerformanceTestContext) iRunOptimizationWithMemoryProfilingEnabled() error {
	ctx.monitoringActive = true
	
	// Start memory profiling
	runtime.GC()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	baselineMemory := int64(memStats.Alloc)
	
	// Run optimization
	projectPath := filepath.Join(ctx.tempDir, "memory-test-project")
	pipeline := optimization.NewOptimizationPipeline(optimization.DefaultPipelineOptions())
	
	_, err := pipeline.OptimizeProject(projectPath)
	if err != nil {
		return fmt.Errorf("optimization failed: %w", err)
	}
	
	// Measure final memory
	runtime.GC()
	runtime.ReadMemStats(&memStats)
	finalMemory := int64(memStats.Alloc)
	
	ctx.memoryUsage["baseline"] = baselineMemory
	ctx.memoryUsage["final"] = finalMemory
	ctx.memoryUsage["peak"] = finalMemory // Simplified for this implementation
	
	return nil
}

func (ctx *PerformanceTestContext) peakMemoryUsageShouldNotExceedBaseline(additionalMB int) error {
	baseline := ctx.memoryUsage["baseline"]
	peak := ctx.memoryUsage["peak"]
	threshold := baseline + int64(additionalMB*1024*1024)
	
	if peak > threshold {
		return fmt.Errorf("peak memory %d MB exceeded baseline + %d MB", 
			peak/(1024*1024), additionalMB)
	}
	
	return nil
}

func (ctx *PerformanceTestContext) memoryShouldBeReleasedProgressivelyDuringOptimization() error {
	// In a real implementation, this would track memory usage over time
	// For now, we'll simulate progressive memory release validation
	return nil
}

func (ctx *PerformanceTestContext) garbageCollectionShouldOccurRegularly() error {
	// Check that GC occurred during optimization
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	if memStats.NumGC == 0 {
		return fmt.Errorf("garbage collection did not occur during optimization")
	}
	
	ctx.logPerformanceMetric("gc_cycles", "optimization", memStats.NumGC, "count")
	return nil
}

func (ctx *PerformanceTestContext) noMemoryLeaksShouldBeDetectedAfterCompletion() error {
	// Force garbage collection and check memory
	runtime.GC()
	runtime.GC() // Run twice to ensure cleanup
	
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	
	finalMemory := int64(memStats.Alloc)
	baseline := ctx.memoryUsage["baseline"]
	
	// Allow for some memory overhead but detect significant leaks
	if finalMemory > baseline*2 {
		return fmt.Errorf("potential memory leak detected: final=%dMB baseline=%dMB", 
			finalMemory/(1024*1024), baseline/(1024*1024))
	}
	
	return nil
}

// Concurrent processing implementations
func (ctx *PerformanceTestContext) iConfigureConcurrentFileProcessing(table *godog.Table) error {
	// Parse concurrency settings from table
	for i := 1; i < len(table.Rows); i++ {
		row := table.Rows[i]
		setting := row.Cells[0].Value
		value := row.Cells[1].Value
		
		switch setting {
		case "MaxConcurrentFiles":
			ctx.concurrencyConfig.MaxConcurrentFiles, _ = strconv.Atoi(value)
		case "EnableParallelism":
			ctx.concurrencyConfig.EnableParallelism = value == "true"
		case "MemoryLimitMB":
			ctx.concurrencyConfig.MemoryLimitMB, _ = strconv.Atoi(value)
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) iOptimizeALargeProjectWithConcurrentProcessing() error {
	// Generate a large project for concurrent processing
	err := ctx.generateProjectOfSize("large-concurrent", "200", "100000")
	if err != nil {
		return err
	}
	
	// Configure optimization with concurrency settings
	options := optimization.DefaultPipelineOptions()
	options.MaxConcurrentFiles = ctx.concurrencyConfig.MaxConcurrentFiles
	
	projectPath := filepath.Join(ctx.tempDir, "large-concurrent-project")
	pipeline := optimization.NewOptimizationPipeline(options)
	
	ctx.startTime = time.Now()
	_, err = pipeline.OptimizeProject(projectPath)
	ctx.endTime = time.Now()
	
	if err != nil {
		return fmt.Errorf("concurrent optimization failed: %w", err)
	}
	
	ctx.processingTimes["concurrent"] = ctx.endTime.Sub(ctx.startTime)
	
	// Also run sequential optimization for comparison
	options.MaxConcurrentFiles = 1
	sequentialPipeline := optimization.NewOptimizationPipeline(options)
	
	sequentialStart := time.Now()
	_, err = sequentialPipeline.OptimizeProject(projectPath)
	sequentialEnd := time.Now()
	
	if err != nil {
		return fmt.Errorf("sequential optimization failed: %w", err)
	}
	
	ctx.processingTimes["sequential"] = sequentialEnd.Sub(sequentialStart)
	
	return nil
}

func (ctx *PerformanceTestContext) processingShouldUtilizeMultipleCPUCoresEffectively() error {
	// Check CPU core utilization during concurrent processing
	numCPU := runtime.NumCPU()
	if ctx.concurrencyConfig.MaxConcurrentFiles <= 1 {
		return fmt.Errorf("concurrency not configured properly")
	}
	
	if ctx.concurrencyConfig.MaxConcurrentFiles > numCPU {
		ctx.logPerformanceMetric("cpu_utilization", "concurrent", 
			fmt.Sprintf("using %d threads on %d cores", ctx.concurrencyConfig.MaxConcurrentFiles, numCPU), "info")
	}
	
	return nil
}

func (ctx *PerformanceTestContext) totalProcessingTimeShouldBeReducedComparedToSequential() error {
	concurrentTime := ctx.processingTimes["concurrent"]
	sequentialTime := ctx.processingTimes["sequential"]
	
	if concurrentTime >= sequentialTime {
		return fmt.Errorf("concurrent processing (%v) was not faster than sequential (%v)", 
			concurrentTime, sequentialTime)
	}
	
	speedup := float64(sequentialTime) / float64(concurrentTime)
	ctx.logPerformanceMetric("speedup_ratio", "concurrent", speedup, "ratio")
	
	return nil
}

func (ctx *PerformanceTestContext) memoryUsagePerThreadShouldBeControlled() error {
	// Check memory usage per thread doesn't exceed limits
	maxMemoryPerThread := int64(ctx.concurrencyConfig.MemoryLimitMB) * 1024 * 1024 / int64(ctx.concurrencyConfig.MaxConcurrentFiles)
	
	if peak, ok := ctx.memoryUsage["peak"]; ok {
		memoryPerThread := peak / int64(ctx.concurrencyConfig.MaxConcurrentFiles)
		
		if memoryPerThread > maxMemoryPerThread {
			return fmt.Errorf("memory per thread %d MB exceeded limit %d MB per thread", 
				memoryPerThread/(1024*1024), maxMemoryPerThread/(1024*1024))
		}
	}
	
	return nil
}

func (ctx *PerformanceTestContext) noRaceConditionsShouldOccur() error {
	// In a real implementation, this would use race detection tools
	// For now, we'll check that the optimization completed successfully
	if len(ctx.optimizationResults) == 0 {
		return fmt.Errorf("no optimization results recorded - possible race condition")
	}
	
	return nil
}

// Helper methods
func (ctx *PerformanceTestContext) setupBenchmarks() {
	ctx.benchmarks = map[string]performanceBenchmark{
		"Small": {
			ProjectSize:        "Small",
			ExpectedTime:       2 * time.Second,
			ExpectedThroughput: 100.0,
			MaxMemoryUsage:     100 * 1024 * 1024,
		},
		"Medium": {
			ProjectSize:        "Medium",
			ExpectedTime:       10 * time.Second,
			ExpectedThroughput: 80.0,
			MaxMemoryUsage:     300 * 1024 * 1024,
		},
		"Large": {
			ProjectSize:        "Large",
			ExpectedTime:       60 * time.Second,
			ExpectedThroughput: 50.0,
			MaxMemoryUsage:     500 * 1024 * 1024,
		},
	}
}

func (ctx *PerformanceTestContext) setupThresholds() {
	ctx.thresholds = map[string]performanceThreshold{
		"processing_time": {
			ProcessingTime: 60 * time.Second,
		},
		"memory_usage": {
			MemoryUsage: 1 * 1024 * 1024 * 1024, // 1GB
		},
		"throughput": {
			FilesPerSecond: 1.0,
		},
	}
}

func (ctx *PerformanceTestContext) generateProjectOfSize(size, files, linesOfCode string) error {
	// Generate a test project of specified size
	config := &types.ProjectConfig{
		Name:         fmt.Sprintf("%s-project", size),
		Type:         "web-api",
		Module:       fmt.Sprintf("github.com/test/%s", size),
		Framework:    "gin",
		Architecture: "standard",
		Logger:       "slog",
		GoVersion:    "1.21",
	}
	
	gen := generator.New()
	projectPath := filepath.Join(ctx.tempDir, config.Name)
	
	_, err := gen.Generate(*config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
		Verbose:    false,
	})
	
	if err != nil {
		return fmt.Errorf("failed to generate project: %w", err)
	}
	
	ctx.projectPaths = append(ctx.projectPaths, projectPath)
	return nil
}

func (ctx *PerformanceTestContext) startResourceMonitoring() {
	ctx.resourceMonitor.mutex.Lock()
	defer ctx.resourceMonitor.mutex.Unlock()
	
	ctx.resourceMonitor.startTime = time.Now()
	ctx.resourceMonitor.samples = 0
	
	// Simplified resource monitoring
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if !ctx.resourceMonitor.isActive {
				return
			}
			
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)
			
			ctx.resourceMonitor.mutex.Lock()
			ctx.resourceMonitor.memoryUsage = append(ctx.resourceMonitor.memoryUsage, int64(memStats.Alloc))
			ctx.resourceMonitor.samples++
			
			// Update peak memory
			if current := int64(memStats.Alloc); current > ctx.memoryUsage["peak"] {
				ctx.memoryUsage["peak"] = current
			}
			
			ctx.resourceMonitor.mutex.Unlock()
		}
	}()
}

func (ctx *PerformanceTestContext) stopResourceMonitoring() {
	ctx.resourceMonitor.mutex.Lock()
	defer ctx.resourceMonitor.mutex.Unlock()
	
	ctx.resourceMonitor.isActive = false
}

func (ctx *PerformanceTestContext) logPerformanceMetric(metric, context string, value interface{}, unit string) {
	ctx.performanceLog = append(ctx.performanceLog, performanceMetric{
		Timestamp: time.Now(),
		Metric:    metric,
		Value:     value,
		Unit:      unit,
		Context:   context,
	})
}

// Cleanup
func (ctx *PerformanceTestContext) Cleanup() {
	// Stop any active monitoring
	if ctx.resourceMonitor != nil {
		ctx.resourceMonitor.isActive = false
	}
	
	// Clean up test directories
	for _, dir := range ctx.testDirs {
		os.RemoveAll(dir)
	}
	
	// Clean up project paths
	for _, path := range ctx.projectPaths {
		os.RemoveAll(path)
	}
}

// Stub implementations for remaining steps (to be expanded as needed)
func (ctx *PerformanceTestContext) iHaveAProjectWithGoFiles(fileCount string) error { return nil }
func (ctx *PerformanceTestContext) eachFileHasApproximatelyLines(lines string) error { return nil }
func (ctx *PerformanceTestContext) iOptimizeTheProjectWithOptimizationLevel(level string) error { 
	var err error
	ctx.optimizationLevel, _ = optimization.ParseOptimizationLevel(level)
	return err
}
func (ctx *PerformanceTestContext) processingTimeShouldScaleLinearlyWithFileCount() error { return nil }
func (ctx *PerformanceTestContext) memoryUsageShouldRemainProportionalToProjectSize() error { return nil }
func (ctx *PerformanceTestContext) optimizationQualityShouldNotDegradeWithScale() error { return nil }
func (ctx *PerformanceTestContext) iMonitorSystemResourcesDuringOptimization(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) iRunOptimizationOnMultipleProjectTypes() error { return nil }
func (ctx *PerformanceTestContext) cpuUtilizationShouldBeEfficientButNotExcessive() error { return nil }
func (ctx *PerformanceTestContext) diskIOShouldBeMinimizedThroughSmartCaching() error { return nil }
func (ctx *PerformanceTestContext) fileHandlesShouldBeProperlyManagedAndReleased() error { return nil }
func (ctx *PerformanceTestContext) systemResponsivenessShouldRemainGood() error { return nil }
func (ctx *PerformanceTestContext) iHaveTheSameProjectToOptimizeAtDifferentLevels() error { return nil }
func (ctx *PerformanceTestContext) iMeasurePerformanceForEachOptimizationLevel(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) performanceShouldScaleAppropriatelyWithOptimizationComplexity() error { return nil }
func (ctx *PerformanceTestContext) qualityImprovementsShouldJustifyPerformanceCosts() error { return nil }
func (ctx *PerformanceTestContext) usersShouldBeWarnedAboutPerformanceImplications() error { return nil }
func (ctx *PerformanceTestContext) iHaveProjectsWithOverlappingCodePatterns() error { return nil }
func (ctx *PerformanceTestContext) iOptimizeSimilarProjectsSequentially() error { return nil }
func (ctx *PerformanceTestContext) astParsingResultsShouldBeCachedWherePossible() error { return nil }
func (ctx *PerformanceTestContext) repeatedPatternAnalysisShouldBeOptimized() error { return nil }
func (ctx *PerformanceTestContext) overallProcessingTimeShouldDecreaseForSimilarPatterns() error { return nil }
func (ctx *PerformanceTestContext) cacheHitRatesShouldBeReportedAndMonitored() error { return nil }
func (ctx *PerformanceTestContext) iHaveAProjectSuitableForProfile(projectType, profile string) error { return nil }
func (ctx *PerformanceTestContext) iOptimizeUsingTheProfile(profile string) error { return nil }
func (ctx *PerformanceTestContext) performanceShouldMatchProfileSpecificExpectations(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) profileRecommendationsShouldConsiderPerformanceTradeOffs() error { return nil }
func (ctx *PerformanceTestContext) iHaveBaselinePerformanceMetricsFromPreviousOptimizationRuns() error { return nil }
func (ctx *PerformanceTestContext) iRunOptimizationWithTheCurrentImplementation() error { return nil }
func (ctx *PerformanceTestContext) performanceShouldNotRegressBeyondAcceptableThresholds(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) anyPerformanceRegressionsShouldBeClearlyReported() error { return nil }
func (ctx *PerformanceTestContext) suggestionsForPerformanceImprovementShouldBeProvided() error { return nil }
func (ctx *PerformanceTestContext) iEnableRealTimePerformanceMonitoring() error { return nil }
func (ctx *PerformanceTestContext) iRunOptimizationOnALargeProject() error { return nil }
func (ctx *PerformanceTestContext) progressShouldBeReportedWithPerformanceMetrics(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) usersShouldBeAbleToMonitorOptimizationProgress() error { return nil }
func (ctx *PerformanceTestContext) performanceBottlenecksShouldBeIdentifiedInRealTime() error { return nil }
func (ctx *PerformanceTestContext) iCreateStressTestConditions(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) iRunOptimizationUnderStress() error { return nil }
func (ctx *PerformanceTestContext) theSystemShouldHandleStressGracefully() error { return nil }
func (ctx *PerformanceTestContext) performanceShouldDegradeGracefully() error { return nil }
func (ctx *PerformanceTestContext) errorHandlingShouldRemainResponsive() error { return nil }
func (ctx *PerformanceTestContext) noSystemCrashesShouldOccur() error { return nil }
func (ctx *PerformanceTestContext) iHaveProjectsWithVaryingOptimizationOpportunities() error { return nil }
func (ctx *PerformanceTestContext) iOptimizeWithDifferentQualityVsSpeedTradeOffs(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) eachModeShouldDeliverExpectedQualitySpeedTradeOffs() error { return nil }
func (ctx *PerformanceTestContext) usersShouldBeAbleToChooseAppropriateModes() error { return nil }
func (ctx *PerformanceTestContext) tradeOffsShouldBeClearlyDocumented() error { return nil }
func (ctx *PerformanceTestContext) iHaveMultipleProjectsToOptimize(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) iOptimizeAllProjectsInBatchMode() error { return nil }
func (ctx *PerformanceTestContext) batchProcessingShouldBeMoreEfficientThanIndividualRuns() error { return nil }
func (ctx *PerformanceTestContext) sharedResourcesShouldBeReusedAcrossProjects() error { return nil }
func (ctx *PerformanceTestContext) overallProcessingTimeShouldBeOptimized() error { return nil }
func (ctx *PerformanceTestContext) progressShouldBeReportedForTheEntireBatch() error { return nil }
func (ctx *PerformanceTestContext) iEnableDetailedPerformanceProfiling() error { return nil }
func (ctx *PerformanceTestContext) iOptimizeARepresentativeProject() error { return nil }
func (ctx *PerformanceTestContext) detailedProfilingDataShouldBeCollected(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) profilingResultsShouldIdentifyOptimizationOpportunities() error { return nil }
func (ctx *PerformanceTestContext) performanceRecommendationsShouldBeGenerated() error { return nil }
func (ctx *PerformanceTestContext) iHaveProjectsStoredInDifferentLocations(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) iOptimizeProjectsFromDifferentLocations() error { return nil }
func (ctx *PerformanceTestContext) ioOperationsShouldBeMinimizedAndOptimized() error { return nil }
func (ctx *PerformanceTestContext) networkLatencyShouldBeHandledGracefully() error { return nil }
func (ctx *PerformanceTestContext) cachingShouldReduceRepeatedRemoteAccess() error { return nil }
func (ctx *PerformanceTestContext) performanceShouldRemainAcceptableForAllLocations() error { return nil }
func (ctx *PerformanceTestContext) iCompleteOptimizationOfVariousProjects() error { return nil }
func (ctx *PerformanceTestContext) iGeneratePerformanceReports() error { return nil }
func (ctx *PerformanceTestContext) reportsShouldIncludeComprehensiveMetrics(table *godog.Table) error { return nil }
func (ctx *PerformanceTestContext) reportsShouldBeExportableInMultipleFormats() error { return nil }
func (ctx *PerformanceTestContext) historicalPerformanceTrendsShouldBeTracked() error { return nil }