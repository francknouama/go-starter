package infrastructure

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelfMaintainingTestInfrastructure(t *testing.T) {
	config := DefaultInfrastructureConfig()
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	assert.NotNil(t, infra)
	assert.Equal(t, config, infra.config)
	assert.NotNil(t, infra.performanceMonitor)
	assert.NotNil(t, infra.regressionDetector)
	assert.NotNil(t, infra.testMaintainer)
	assert.NotNil(t, infra.dependencyAnalyzer)
	assert.NotNil(t, infra.optimizationEngine)
	assert.NotNil(t, infra.reportGenerator)
	assert.NotNil(t, infra.healthChecker)
	assert.NotNil(t, infra.automationScheduler)
	assert.False(t, infra.running)
}

func TestDefaultInfrastructureConfig(t *testing.T) {
	config := DefaultInfrastructureConfig()
	
	// Core settings
	assert.Equal(t, ".", config.ProjectRoot)
	assert.Equal(t, "./tests", config.TestDirectory)
	assert.Equal(t, time.Hour*6, config.MaintenanceInterval)
	assert.Equal(t, time.Minute*15, config.HealthCheckInterval)
	
	// Performance thresholds
	assert.Equal(t, time.Minute*5, config.MaxTestDuration)
	assert.Equal(t, time.Minute*30, config.MaxSuiteDuration)
	assert.Equal(t, int64(512), config.MaxMemoryUsage)
	assert.Equal(t, 80.0, config.MaxCPUUsage)
	
	// Regression detection
	assert.True(t, config.EnableRegressionDetection)
	assert.Equal(t, 10.0, config.RegressionThreshold)
	assert.Equal(t, 30, config.PerformanceHistoryDays)
	assert.Equal(t, time.Hour*24, config.BaselineUpdateFrequency)
	
	// Maintenance actions
	assert.False(t, config.AutoFixFailingTests) // Conservative default
	assert.True(t, config.AutoOptimizeSlowTests)
	assert.False(t, config.AutoUpdateDependencies) // Conservative default
	assert.True(t, config.AutoCleanupObsoleteTests)
	assert.True(t, config.AutoGenerateMissingTests)
	
	// Quality gates
	assert.Equal(t, 80.0, config.MinCodeCoverage)
	assert.Equal(t, 5.0, config.MaxFailureRate)
	assert.Equal(t, 2.0, config.MaxFlakyTestRate)
	assert.True(t, config.RequireDocumentation)
	
	// Integration settings
	assert.True(t, config.ContinuousMonitoring)
	assert.False(t, config.IntegrateWithCI)
	assert.Equal(t, []string{"console"}, config.NotificationChannels)
	assert.True(t, config.BackupBeforeChanges)
	
	// Advanced features
	assert.False(t, config.PredictiveOptimization)
	assert.True(t, config.LearningEnabled)
	assert.False(t, config.ExperimentalFeatures)
}

func TestInfrastructureStartStop(t *testing.T) {
	config := DefaultInfrastructureConfig()
	config.ContinuousMonitoring = false // Disable to avoid background processes
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	// Test start
	err := infra.Start()
	assert.NoError(t, err)
	assert.True(t, infra.running)
	
	// Test start when already running
	err = infra.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
	
	// Test stop
	err = infra.Stop()
	assert.NoError(t, err)
	assert.False(t, infra.running)
	
	// Test stop when not running
	err = infra.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestPerformanceMetricSerialization(t *testing.T) {
	metric := PerformanceMetric{
		TestName:     "TestExample",
		TestSuite:    "main",
		Duration:     time.Millisecond * 150,
		MemoryUsage:  1024 * 1024, // 1MB
		CPUUsage:     45.5,
		Timestamp:    time.Now(),
		Passed:       true,
		ErrorMessage: "",
	}
	
	// Test JSON serialization
	data, err := json.Marshal(metric)
	require.NoError(t, err)
	assert.Contains(t, string(data), "TestExample")
	assert.Contains(t, string(data), "150000000") // 150ms in nanoseconds
	
	// Test JSON deserialization
	var decoded PerformanceMetric
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, metric.TestName, decoded.TestName)
	assert.Equal(t, metric.TestSuite, decoded.TestSuite)
	assert.Equal(t, metric.Duration, decoded.Duration)
	assert.Equal(t, metric.MemoryUsage, decoded.MemoryUsage)
	assert.Equal(t, metric.CPUUsage, decoded.CPUUsage)
	assert.Equal(t, metric.Passed, decoded.Passed)
}

func TestPerformanceMonitor(t *testing.T) {
	config := DefaultInfrastructureConfig()
	monitor := NewPerformanceMonitor(config)
	
	assert.NotNil(t, monitor)
	assert.NotNil(t, monitor.metrics)
	assert.NotNil(t, monitor.baselines)
	
	// Test adding metrics
	metrics := []PerformanceMetric{
		{
			TestName:  "TestFast",
			TestSuite: "unit",
			Duration:  time.Millisecond * 50,
			Timestamp: time.Now(),
			Passed:    true,
		},
		{
			TestName:  "TestSlow",
			TestSuite: "integration",
			Duration:  time.Second * 3,
			Timestamp: time.Now(),
			Passed:    true,
		},
		{
			TestName:  "TestFailing",
			TestSuite: "unit",
			Duration:  time.Millisecond * 100,
			Timestamp: time.Now(),
			Passed:    false,
		},
	}
	
	monitor.AddMetrics(metrics)
	
	// Verify metrics were added
	assert.Len(t, monitor.metrics, 3)
	assert.Contains(t, monitor.metrics, "TestFast")
	assert.Contains(t, monitor.metrics, "TestSlow")
	assert.Contains(t, monitor.metrics, "TestFailing")
	
	// Test summary generation
	summary := monitor.GenerateSummary()
	assert.Equal(t, 3, summary.TotalTests)
	assert.Equal(t, 2, summary.PassingTests)
	assert.Equal(t, 1, summary.FailingTests)
	assert.Greater(t, summary.AverageDuration, time.Duration(0))
	
	// Should identify slow test
	assert.Contains(t, summary.SlowestTests, "TestSlow")
}

func TestRegressionDetector(t *testing.T) {
	config := DefaultInfrastructureConfig()
	detector := NewRegressionDetector(config)
	
	assert.NotNil(t, detector)
	assert.NotNil(t, detector.history)
	assert.NotNil(t, detector.alerts)
	
	// Test regression detection
	regressions := detector.DetectRegressions()
	
	// Initial call might have placeholder regressions
	assert.NotNil(t, regressions)
	
	// Test with actual history (simplified)
	snapshot1 := PerformanceSnapshot{
		Timestamp: time.Now().Add(-time.Hour),
		Metrics: map[string]PerformanceMetric{
			"TestExample": {
				TestName: "TestExample",
				Duration: time.Second,
			},
		},
	}
	
	snapshot2 := PerformanceSnapshot{
		Timestamp: time.Now(),
		Metrics: map[string]PerformanceMetric{
			"TestExample": {
				TestName: "TestExample",
				Duration: time.Second * 2, // 100% regression
			},
		},
	}
	
	detector.mu.Lock()
	detector.history = []PerformanceSnapshot{snapshot1, snapshot2}
	detector.mu.Unlock()
	
	regressions = detector.DetectRegressions()
	assert.Greater(t, len(regressions), 0)
}

func TestTestMaintainer(t *testing.T) {
	config := DefaultInfrastructureConfig()
	maintainer := NewTestMaintainer(config)
	
	assert.NotNil(t, maintainer)
	assert.NotNil(t, maintainer.maintenanceHistory)
	
	// Test optimizing slow test
	action := maintainer.OptimizeSlowTest("SlowTest")
	require.NotNil(t, action)
	
	assert.Equal(t, "optimize", action.Type)
	assert.Equal(t, "SlowTest", action.Target)
	assert.Equal(t, "completed", action.Status)
	assert.NotEmpty(t, action.Description)
	assert.Greater(t, len(action.Changes), 0)
	assert.False(t, action.StartTime.IsZero())
	assert.False(t, action.EndTime.IsZero())
	
	// Verify action was recorded
	assert.Len(t, maintainer.maintenanceHistory, 1)
	
	// Test cleanup obsolete tests
	cleanupActions := maintainer.CleanupObsoleteTests()
	assert.Greater(t, len(cleanupActions), 0)
	
	for _, action := range cleanupActions {
		assert.Equal(t, "cleanup", action.Type)
		assert.Equal(t, "completed", action.Status)
		assert.NotEmpty(t, action.Description)
	}
	
	// Test generate missing tests
	generateActions := maintainer.GenerateMissingTests()
	assert.Greater(t, len(generateActions), 0)
	
	for _, action := range generateActions {
		assert.Equal(t, "generate", action.Type)
		assert.Equal(t, "completed", action.Status)
		assert.NotEmpty(t, action.Description)
	}
	
	// Verify all actions were recorded
	assert.Greater(t, len(maintainer.maintenanceHistory), 1)
}

func TestDependencyAnalyzer(t *testing.T) {
	config := DefaultInfrastructureConfig()
	analyzer := NewDependencyAnalyzer(config)
	
	assert.NotNil(t, analyzer)
	assert.NotNil(t, analyzer.dependencies)
	assert.NotNil(t, analyzer.conflicts)
	
	// Test dependency analysis
	err := analyzer.AnalyzeDependencies()
	assert.NoError(t, err)
	
	// Should have found some dependencies
	assert.Greater(t, len(analyzer.dependencies), 0)
	
	// Check for expected dependency
	testify, exists := analyzer.dependencies["testify"]
	assert.True(t, exists)
	assert.Equal(t, "github.com/stretchr/testify", testify.Name)
	assert.Equal(t, "test", testify.Type)
	assert.True(t, testify.Required)
	assert.False(t, testify.LastChecked.IsZero())
}

func TestOptimizationEngine(t *testing.T) {
	config := DefaultInfrastructureConfig()
	engine := NewOptimizationEngine(config)
	
	assert.NotNil(t, engine)
	assert.Greater(t, len(engine.optimizations), 0)
	assert.NotNil(t, engine.learned)
	
	// Test running optimizations
	results := engine.RunOptimizations()
	
	// Should have applied some optimizations
	assert.Greater(t, len(results), 0)
	
	for _, result := range results {
		assert.NotEmpty(t, result.RuleID)
		assert.NotEmpty(t, result.RuleName)
		assert.NotEmpty(t, result.Target)
		assert.True(t, result.Applied)
		assert.NotEmpty(t, result.ImprovementType)
		assert.Greater(t, result.BeforeValue, 0.0)
		assert.Greater(t, result.AfterValue, 0.0)
		assert.Greater(t, result.ImprovementPercent, 0.0)
		assert.False(t, result.Timestamp.IsZero())
	}
	
	// Test optimization rules
	assert.Greater(t, len(engine.optimizations), 0)
	for _, rule := range engine.optimizations {
		assert.NotEmpty(t, rule.ID)
		assert.NotEmpty(t, rule.Name)
		assert.NotEmpty(t, rule.Condition)
		assert.NotEmpty(t, rule.Action)
		assert.Greater(t, rule.Priority, 0)
		assert.GreaterOrEqual(t, rule.SuccessRate, 0.0)
		assert.LessOrEqual(t, rule.SuccessRate, 100.0)
	}
}

func TestReportGenerator(t *testing.T) {
	config := DefaultInfrastructureConfig()
	generator := NewReportGenerator(config)
	
	assert.NotNil(t, generator)
	assert.NotNil(t, generator.reports)
	
	// Create test report
	report := InfrastructureReport{
		ID:            "test_report_1",
		Timestamp:     time.Now(),
		Period:        "test_period",
		OverallHealth: "healthy",
		PerformanceSummary: PerformanceSummary{
			TotalTests:      10,
			PassingTests:    9,
			FailingTests:    1,
			AverageDuration: time.Millisecond * 100,
		},
		MaintenanceActions: []MaintenanceAction{
			{
				ID:          "action_1",
				Type:        "optimize",
				Target:      "TestSlow",
				Description: "Optimized slow test",
				Status:      "completed",
			},
		},
		Regressions: []RegressionAlert{
			{
				TestName:          "TestRegressed",
				RegressionType:    "duration",
				RegressionPercent: 25.0,
				Severity:          "medium",
			},
		},
		Optimizations: []OptimizationResult{
			{
				RuleID:             "rule_1",
				RuleName:           "Parallel Execution",
				Applied:            true,
				ImprovementPercent: 30.0,
			},
		},
		Recommendations: []string{
			"Optimize slow tests",
			"Fix failing tests",
		},
		GenerationTime: time.Millisecond * 50,
	}
	
	// Add report
	generator.AddReport(report)
	
	// Verify report was added
	assert.Len(t, generator.reports, 1)
	assert.Equal(t, report.ID, generator.reports[0].ID)
	assert.Equal(t, report.OverallHealth, generator.reports[0].OverallHealth)
}

func TestHealthChecker(t *testing.T) {
	config := DefaultInfrastructureConfig()
	checker := NewHealthChecker(config)
	
	assert.NotNil(t, checker)
	assert.NotNil(t, checker.healthHistory)
	
	// Test with no health checks
	latest := checker.GetLatestHealthCheck()
	assert.Nil(t, latest)
	
	// Add health check
	healthCheck := HealthCheck{
		Timestamp:     time.Now(),
		OverallHealth: "healthy",
		Components: map[string]string{
			"tests":        "healthy",
			"performance":  "healthy",
			"dependencies": "healthy",
		},
		Issues: []HealthIssue{
			{
				Component:   "tests",
				Type:        "warning",
				Severity:    "low",
				Description: "Minor test issue",
				AutoFixable: true,
			},
		},
		Metrics: map[string]float64{
			"test_success_rate": 95.0,
			"avg_duration":      2.5,
		},
	}
	
	checker.AddHealthCheck(healthCheck)
	
	// Verify health check was added
	assert.Len(t, checker.healthHistory, 1)
	
	latest = checker.GetLatestHealthCheck()
	require.NotNil(t, latest)
	assert.Equal(t, healthCheck.OverallHealth, latest.OverallHealth)
	assert.Equal(t, len(healthCheck.Components), len(latest.Components))
	assert.Equal(t, len(healthCheck.Issues), len(latest.Issues))
	assert.Equal(t, len(healthCheck.Metrics), len(latest.Metrics))
}

func TestAutomationScheduler(t *testing.T) {
	config := DefaultInfrastructureConfig()
	scheduler := NewAutomationScheduler(config)
	
	assert.NotNil(t, scheduler)
	assert.Greater(t, len(scheduler.tasks), 0)
	assert.NotNil(t, scheduler.queue)
	
	// Test scheduled tasks
	for _, task := range scheduler.tasks {
		assert.NotEmpty(t, task.ID)
		assert.NotEmpty(t, task.Name)
		assert.NotEmpty(t, task.Type)
		assert.NotEmpty(t, task.Schedule)
		assert.Greater(t, task.Priority, 0)
		assert.Greater(t, task.MaxDuration, time.Duration(0))
	}
	
	// Test processing scheduled tasks
	initialQueueSize := len(scheduler.queue)
	
	// Set a task to run now
	scheduler.mu.Lock()
	if len(scheduler.tasks) > 0 {
		scheduler.tasks[0].NextRun = time.Now().Add(-time.Minute) // Past time
		scheduler.tasks[0].Enabled = true
	}
	scheduler.mu.Unlock()
	
	scheduler.ProcessScheduledTasks()
	
	// Should have queued a task
	assert.Greater(t, len(scheduler.queue), initialQueueSize)
	
	// Test next run calculation
	now := time.Now()
	nextHourly := scheduler.calculateNextRun("hourly", now)
	nextDaily := scheduler.calculateNextRun("daily", now)
	nextWeekly := scheduler.calculateNextRun("weekly", now)
	
	assert.True(t, nextHourly.After(now))
	assert.True(t, nextDaily.After(nextHourly))
	assert.True(t, nextWeekly.After(nextDaily))
}

func TestParseTestOutput(t *testing.T) {
	config := DefaultInfrastructureConfig()
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	testOutput := `=== RUN   TestExample
--- PASS: TestExample (0.12s)
=== RUN   TestAnother
--- PASS: TestAnother (0.05s)
=== RUN   TestSlow
--- PASS: TestSlow (2.45s)
=== RUN   TestFailing
--- FAIL: TestFailing (0.08s)
    test_file.go:25: assertion failed
PASS
ok  	example/package	2.7s`
	
	metrics := infra.parseTestOutput(testOutput)
	
	assert.Len(t, metrics, 4)
	
	// Check each parsed metric
	testNames := make(map[string]bool)
	for _, metric := range metrics {
		testNames[metric.TestName] = true
		assert.NotEmpty(t, metric.TestName)
		assert.NotEmpty(t, metric.TestSuite)
		assert.Greater(t, metric.Duration, time.Duration(0))
		assert.False(t, metric.Timestamp.IsZero())
	}
	
	assert.True(t, testNames["TestExample"])
	assert.True(t, testNames["TestAnother"])
	assert.True(t, testNames["TestSlow"])
	assert.True(t, testNames["TestFailing"])
	
	// Check specific test results
	for _, metric := range metrics {
		switch metric.TestName {
		case "TestExample":
			assert.Equal(t, time.Millisecond*120, metric.Duration)
			assert.True(t, metric.Passed)
		case "TestSlow":
			assert.Equal(t, time.Millisecond*2450, metric.Duration)
			assert.True(t, metric.Passed)
		case "TestFailing":
			assert.Equal(t, time.Millisecond*80, metric.Duration)
			assert.False(t, metric.Passed)
		}
	}
}

func TestExtractTestSuite(t *testing.T) {
	config := DefaultInfrastructureConfig()
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	testCases := []struct {
		testName     string
		expectedSuite string
	}{
		{"TestSimple", "default"},
		{"package/TestWithPath", "package"},
		{"complex/path/TestDeep", "complex"},
		{"TestNoPath", "default"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.testName, func(t *testing.T) {
			suite := infra.extractTestSuite(tc.testName)
			assert.Equal(t, tc.expectedSuite, suite)
		})
	}
}

func TestHealthChecking(t *testing.T) {
	config := DefaultInfrastructureConfig()
	config.MaxFailureRate = 10.0
	config.MaxTestDuration = time.Second
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	// Add some test metrics
	metrics := []PerformanceMetric{
		{TestName: "Test1", Duration: time.Millisecond * 100, Passed: true},
		{TestName: "Test2", Duration: time.Millisecond * 200, Passed: true},
		{TestName: "Test3", Duration: time.Millisecond * 150, Passed: false},
		{TestName: "Test4", Duration: time.Second * 2, Passed: true}, // Slow test
	}
	infra.performanceMonitor.AddMetrics(metrics)
	
	// Test test execution health
	testHealth := infra.checkTestExecutionHealth()
	assert.Contains(t, []string{"healthy", "warning", "critical"}, testHealth)
	
	// Test performance health
	perfHealth := infra.checkPerformanceHealth()
	assert.Contains(t, []string{"healthy", "warning", "critical"}, perfHealth)
	
	// Test dependency health
	depHealth := infra.checkDependencyHealth()
	assert.Contains(t, []string{"healthy", "warning", "critical"}, depHealth)
	
	// Test overall health calculation
	components := map[string]string{
		"test_execution": "healthy",
		"performance":    "warning",
		"dependencies":   "healthy",
	}
	
	overallHealth := infra.calculateHealthFromComponents(components)
	assert.Equal(t, "warning", overallHealth)
	
	components["performance"] = "critical"
	overallHealth = infra.calculateHealthFromComponents(components)
	assert.Equal(t, "critical", overallHealth)
	
	components["performance"] = "healthy"
	overallHealth = infra.calculateHealthFromComponents(components)
	assert.Equal(t, "healthy", overallHealth)
}

func TestRunMaintenance(t *testing.T) {
	// Create temporary project directory
	tempDir, err := os.MkdirTemp("", "maintenance-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create a simple test file
	testFile := filepath.Join(tempDir, "example_test.go")
	testContent := `package main

import "testing"

func TestExample(t *testing.T) {
	// Simple test
}
`
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)
	
	config := DefaultInfrastructureConfig()
	config.ProjectRoot = tempDir
	config.ContinuousMonitoring = false
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	err = infra.Start()
	require.NoError(t, err)
	defer func() {
		if err := infra.Stop(); err != nil {
			t.Logf("Warning: failed to stop infrastructure: %v", err)
		}
	}()
	
	// Run maintenance
	report, err := infra.RunMaintenance()
	require.NoError(t, err)
	assert.NotNil(t, report)
	
	// Verify report contents
	assert.NotEmpty(t, report.ID)
	assert.False(t, report.Timestamp.IsZero())
	assert.Equal(t, "maintenance_cycle", report.Period)
	assert.NotEmpty(t, report.OverallHealth)
	assert.NotNil(t, report.PerformanceSummary)
	assert.NotNil(t, report.MaintenanceActions)
	assert.NotNil(t, report.Regressions)
	assert.NotNil(t, report.Optimizations)
	assert.NotNil(t, report.Recommendations)
	assert.NotNil(t, report.TrendAnalysis)
	assert.Greater(t, report.GenerationTime, time.Duration(0))
}

func TestGetCurrentStatus(t *testing.T) {
	config := DefaultInfrastructureConfig()
	config.ContinuousMonitoring = false
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	// Test when not running
	status, err := infra.GetCurrentStatus()
	assert.Error(t, err)
	assert.Nil(t, status)
	
	// Start infrastructure
	err = infra.Start()
	require.NoError(t, err)
	defer func() {
		if err := infra.Stop(); err != nil {
			t.Logf("Warning: failed to stop infrastructure: %v", err)
		}
	}()
	
	// Perform a health check to have data
	err = infra.performHealthCheck()
	require.NoError(t, err)
	
	// Test when running
	status, err = infra.GetCurrentStatus()
	assert.NoError(t, err)
	assert.NotNil(t, status)
	assert.False(t, status.Timestamp.IsZero())
	assert.NotEmpty(t, status.OverallHealth)
}

func TestInfrastructureReportSerialization(t *testing.T) {
	report := InfrastructureReport{
		ID:            "test_report",
		Timestamp:     time.Now(),
		Period:        "test",
		OverallHealth: "healthy",
		PerformanceSummary: PerformanceSummary{
			TotalTests:      10,
			PassingTests:    9,
			FailingTests:    1,
			AverageDuration: time.Millisecond * 100,
			SlowestTests:    []string{"TestSlow"},
		},
		MaintenanceActions: []MaintenanceAction{
			{
				ID:          "action_1",
				Type:        "optimize",
				Target:      "TestTarget",
				Description: "Test action",
				Status:      "completed",
				StartTime:   time.Now(),
			},
		},
		Regressions: []RegressionAlert{
			{
				TestName:          "TestRegressed",
				RegressionType:    "duration",
				Current:           5.0,
				Baseline:          3.0,
				RegressionPercent: 66.7,
				Severity:          "high",
				DetectedAt:        time.Now(),
				Resolved:          false,
			},
		},
		Optimizations: []OptimizationResult{
			{
				RuleID:             "rule_1",
				RuleName:           "Test Optimization",
				Target:             "TestTarget",
				Applied:            true,
				ImprovementType:    "performance",
				BeforeValue:        10.0,
				AfterValue:         7.0,
				ImprovementPercent: 30.0,
				Timestamp:          time.Now(),
			},
		},
		Recommendations: []string{
			"Fix failing tests",
			"Optimize slow tests",
		},
		TrendAnalysis: TrendAnalysis{
			PerformanceTrend:     "improving",
			QualityTrend:         "stable",
			MaintenanceFrequency: "decreasing",
			PredictedIssues:      []string{"None predicted"},
			Confidence:           0.85,
		},
		GenerationTime: time.Millisecond * 150,
	}
	
	// Test JSON serialization
	data, err := json.Marshal(report)
	require.NoError(t, err)
	assert.Contains(t, string(data), "test_report")
	assert.Contains(t, string(data), "healthy")
	
	// Test JSON deserialization
	var decoded InfrastructureReport
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, report.ID, decoded.ID)
	assert.Equal(t, report.OverallHealth, decoded.OverallHealth)
	assert.Equal(t, report.PerformanceSummary.TotalTests, decoded.PerformanceSummary.TotalTests)
	assert.Equal(t, len(report.MaintenanceActions), len(decoded.MaintenanceActions))
	assert.Equal(t, len(report.Regressions), len(decoded.Regressions))
	assert.Equal(t, len(report.Optimizations), len(decoded.Optimizations))
	assert.Equal(t, len(report.Recommendations), len(decoded.Recommendations))
	assert.Equal(t, report.TrendAnalysis.PerformanceTrend, decoded.TrendAnalysis.PerformanceTrend)
}

func TestGenerateRecommendations(t *testing.T) {
	config := DefaultInfrastructureConfig()
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	// Add metrics with various issues
	metrics := []PerformanceMetric{
		{TestName: "TestFast", Duration: time.Millisecond * 50, Passed: true},
		{TestName: "TestSlow1", Duration: time.Second * 3, Passed: true},
		{TestName: "TestSlow2", Duration: time.Second * 4, Passed: true},
		{TestName: "TestFailing1", Duration: time.Millisecond * 100, Passed: false},
		{TestName: "TestFailing2", Duration: time.Millisecond * 150, Passed: false},
	}
	infra.performanceMonitor.AddMetrics(metrics)
	
	recommendations := infra.generateRecommendations()
	
	assert.Greater(t, len(recommendations), 0)
	
	// Should recommend fixing failing tests
	hasFailingRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "Fix") && strings.Contains(rec, "failing") {
			hasFailingRecommendation = true
			break
		}
	}
	assert.True(t, hasFailingRecommendation)
	
	// Should recommend optimizing slow tests
	hasSlowRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "Optimize slow tests") {
			hasSlowRecommendation = true
			break
		}
	}
	assert.True(t, hasSlowRecommendation)
}

func TestTrendAnalysis(t *testing.T) {
	config := DefaultInfrastructureConfig()
	infra := NewSelfMaintainingTestInfrastructure(config)
	
	trend := infra.analyzeTrends()
	
	assert.NotEmpty(t, trend.PerformanceTrend)
	assert.Contains(t, []string{"improving", "stable", "declining"}, trend.PerformanceTrend)
	
	assert.NotEmpty(t, trend.QualityTrend)
	assert.Contains(t, []string{"improving", "stable", "declining"}, trend.QualityTrend)
	
	assert.NotEmpty(t, trend.MaintenanceFrequency)
	assert.Contains(t, []string{"increasing", "stable", "decreasing"}, trend.MaintenanceFrequency)
	
	assert.NotNil(t, trend.PredictedIssues)
	assert.GreaterOrEqual(t, trend.Confidence, 0.0)
	assert.LessOrEqual(t, trend.Confidence, 1.0)
}