package monitoring

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCoverageMonitor(t *testing.T) {
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	assert.NotNil(t, monitor)
	assert.Equal(t, config, monitor.config)
	assert.NotNil(t, monitor.fileWatcher)
	assert.NotNil(t, monitor.coverageStore)
	assert.NotNil(t, monitor.alertManager)
	assert.NotNil(t, monitor.metrics)
	assert.NotNil(t, monitor.regressionTracker)
	assert.False(t, monitor.running)
	assert.Greater(t, len(monitor.qualityGates), 0)
}

func TestDefaultMonitorConfig(t *testing.T) {
	config := DefaultMonitorConfig()
	
	// Monitoring settings
	assert.Equal(t, []string{"."}, config.WatchPaths)
	assert.Contains(t, config.ExcludePatterns, "vendor")
	assert.Contains(t, config.ExcludePatterns, ".git")
	assert.Equal(t, time.Minute*5, config.MonitorInterval)
	assert.True(t, config.RealtimeUpdates)
	
	// Coverage thresholds
	assert.Equal(t, 80.0, config.MinCoveragePercent)
	assert.Equal(t, 70.0, config.CriticalThreshold)
	assert.Equal(t, 75.0, config.WarningThreshold)
	
	// Quality gates
	assert.True(t, config.EnableQualityGates)
	assert.True(t, config.FailOnRegression)
	assert.Equal(t, 5.0, config.MaxRegressionPercent)
	
	// Reporting
	assert.True(t, config.GenerateReports)
	assert.Equal(t, "json", config.ReportFormat)
	assert.Equal(t, "./coverage-reports", config.ReportOutputPath)
	
	// Advanced features
	assert.True(t, config.EnableTrendAnalysis)
	assert.Equal(t, 30, config.HistoryRetentionDays)
	assert.False(t, config.EnablePrediction)
}

func TestCoverageMonitorStartStop(t *testing.T) {
	config := DefaultMonitorConfig()
	config.RealtimeUpdates = false // Disable to avoid file watcher setup
	monitor := NewCoverageMonitor(config)
	
	// Test start
	err := monitor.Start()
	assert.NoError(t, err)
	assert.True(t, monitor.running)
	
	// Test start when already running
	err = monitor.Start()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already running")
	
	// Test stop
	err = monitor.Stop()
	assert.NoError(t, err)
	assert.False(t, monitor.running)
	
	// Test stop when not running
	err = monitor.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

func TestCoverageDataSerialization(t *testing.T) {
	data := CoverageData{
		FilePath:        "test.go",
		PackageName:     "main",
		TotalLines:      100,
		CoveredLines:    80,
		CoveragePercent: 80.0,
		UncoveredLines:  []int{10, 15, 25},
		FunctionCoverage: map[string]float64{
			"main":    90.0,
			"helper":  70.0,
		},
		BranchCoverage: 85.0,
		LastUpdated:    time.Now(),
		TestFiles:      []string{"test_test.go"},
		Dependencies:   []string{"fmt", "os"},
	}
	
	// Test JSON serialization
	jsonData, err := json.Marshal(data)
	require.NoError(t, err)
	assert.Contains(t, string(jsonData), "test.go")
	assert.Contains(t, string(jsonData), "main")
	
	// Test JSON deserialization
	var decoded CoverageData
	err = json.Unmarshal(jsonData, &decoded)
	require.NoError(t, err)
	assert.Equal(t, data.FilePath, decoded.FilePath)
	assert.Equal(t, data.PackageName, decoded.PackageName)
	assert.Equal(t, data.TotalLines, decoded.TotalLines)
	assert.Equal(t, data.CoveredLines, decoded.CoveredLines)
	assert.InDelta(t, data.CoveragePercent, decoded.CoveragePercent, 0.001)
}

func TestQualityGateEvaluation(t *testing.T) {
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	testCases := []struct {
		name      string
		gate      QualityGate
		report    *CoverageReport
		expected  bool
	}{
		{
			name: "coverage gate passes",
			gate: QualityGate{
				Name:      "Min Coverage",
				Type:      "coverage",
				Condition: ">=",
				Threshold: 80.0,
				Enabled:   true,
			},
			report: &CoverageReport{
				OverallCoverage: 85.0,
			},
			expected: true,
		},
		{
			name: "coverage gate fails",
			gate: QualityGate{
				Name:      "Min Coverage",
				Type:      "coverage",
				Condition: ">=",
				Threshold: 90.0,
				Enabled:   true,
			},
			report: &CoverageReport{
				OverallCoverage: 85.0,
			},
			expected: false,
		},
		{
			name: "regression gate passes",
			gate: QualityGate{
				Name:      "No Regression",
				Type:      "regression",
				Condition: "<=",
				Threshold: 5.0,
				Enabled:   true,
			},
			report: &CoverageReport{
				RegressionAnalysis: RegressionAnalysis{
					HasRegression:     true,
					RegressionPercent: 3.0,
				},
			},
			expected: true,
		},
		{
			name: "regression gate fails",
			gate: QualityGate{
				Name:      "No Regression",
				Type:      "regression",
				Condition: "<=",
				Threshold: 2.0,
				Enabled:   true,
			},
			report: &CoverageReport{
				RegressionAnalysis: RegressionAnalysis{
					HasRegression:     true,
					RegressionPercent: 5.0,
				},
			},
			expected: false,
		},
		{
			name: "disabled gate",
			gate: QualityGate{
				Name:      "Disabled Gate",
				Type:      "coverage",
				Condition: ">=",
				Threshold: 100.0,
				Enabled:   false,
			},
			report: &CoverageReport{
				OverallCoverage: 50.0,
			},
			expected: true, // Should not be evaluated
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.gate.Enabled {
				// Disabled gates should not appear in results
				results := monitor.checkQualityGates(tc.report)
				for _, result := range results {
					assert.NotEqual(t, tc.gate.Name, result.Gate.Name)
				}
				return
			}
			
			result := monitor.evaluateQualityGate(tc.gate, tc.report)
			assert.Equal(t, tc.expected, result.Passed)
			assert.Equal(t, tc.gate.Name, result.Gate.Name)
			assert.NotEmpty(t, result.Message)
			assert.False(t, result.Timestamp.IsZero())
		})
	}
}

func TestEvaluateCondition(t *testing.T) {
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	testCases := []struct {
		actual    float64
		condition string
		threshold float64
		expected  bool
	}{
		{85.0, ">=", 80.0, true},
		{75.0, ">=", 80.0, false},
		{85.0, "<=", 90.0, true},
		{95.0, "<=", 90.0, false},
		{85.0, ">", 80.0, true},
		{80.0, ">", 80.0, false},
		{75.0, "<", 80.0, true},
		{85.0, "<", 80.0, false},
		{80.0, "==", 80.0, true},
		{85.0, "==", 80.0, false},
		{85.0, "!=", 80.0, true},
		{80.0, "!=", 80.0, false},
		{80.0, "invalid", 80.0, false},
	}
	
	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			result := monitor.evaluateCondition(tc.actual, tc.condition, tc.threshold)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestCustomQualityGate(t *testing.T) {
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	// Add custom quality gate
	customGate := QualityGate{
		Name:        "Custom Gate",
		Type:        "custom",
		Enabled:     true,
		Critical:    false,
		Description: "Custom logic for testing",
		CustomFunc: func(data CoverageData) bool {
			return data.CoveragePercent >= 75.0
		},
	}
	
	monitor.AddQualityGate(customGate)
	
	// Verify gate was added
	found := false
	for _, gate := range monitor.qualityGates {
		if gate.Name == "Custom Gate" {
			found = true
			break
		}
	}
	assert.True(t, found, "Custom gate should be added")
	
	// Test evaluation
	report := &CoverageReport{
		OverallCoverage: 80.0,
	}
	
	result := monitor.evaluateQualityGate(customGate, report)
	assert.True(t, result.Passed)
	assert.Equal(t, 80.0, result.ActualValue)
	assert.Contains(t, result.Message, "Custom gate evaluation")
}

func TestRegressionAnalysis(t *testing.T) {
	tracker := NewRegressionTracker(10)
	
	// Set baseline
	baseline := map[string]float64{
		"file1.go": 90.0,
		"file2.go": 85.0,
		"file3.go": 80.0,
	}
	tracker.SetBaseline(baseline)
	
	// Test with no regression
	report1 := &CoverageReport{
		FileCoverage: []CoverageData{
			{FilePath: "file1.go", CoveragePercent: 92.0},
			{FilePath: "file2.go", CoveragePercent: 87.0},
			{FilePath: "file3.go", CoveragePercent: 82.0},
		},
	}
	
	analysis1 := tracker.AnalyzeRegression(report1)
	assert.False(t, analysis1.HasRegression)
	assert.Equal(t, 0.0, analysis1.RegressionPercent)
	assert.Equal(t, "none", analysis1.Severity)
	assert.Empty(t, analysis1.RegressionFiles)
	
	// Test with regression
	report2 := &CoverageReport{
		FileCoverage: []CoverageData{
			{FilePath: "file1.go", CoveragePercent: 85.0}, // 5% regression
			{FilePath: "file2.go", CoveragePercent: 80.0}, // 5% regression
			{FilePath: "file3.go", CoveragePercent: 82.0}, // No regression
		},
	}
	
	analysis2 := tracker.AnalyzeRegression(report2)
	assert.True(t, analysis2.HasRegression)
	assert.Equal(t, 5.0, analysis2.RegressionPercent)
	assert.Equal(t, "high", analysis2.Severity)
	assert.Contains(t, analysis2.RegressionFiles, "file1.go")
	assert.Contains(t, analysis2.RegressionFiles, "file2.go")
	assert.NotContains(t, analysis2.RegressionFiles, "file3.go")
	assert.Equal(t, 5.0, analysis2.RegressionDetails["file1.go"])
	assert.Equal(t, 5.0, analysis2.RegressionDetails["file2.go"])
	
	// Test severity levels
	testCases := []struct {
		regressionPercent float64
		expectedSeverity  string
	}{
		{15.0, "critical"},
		{8.0, "high"},
		{3.0, "medium"},
		{1.0, "low"},
	}
	
	for _, tc := range testCases {
		report := &CoverageReport{
			FileCoverage: []CoverageData{
				{FilePath: "file1.go", CoveragePercent: baseline["file1.go"] - tc.regressionPercent},
			},
		}
		
		analysis := tracker.AnalyzeRegression(report)
		assert.Equal(t, tc.expectedSeverity, analysis.Severity)
	}
}

func TestCoverageStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "coverage-store-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	store := NewCoverageStore(tempDir, 5)
	
	// Create test report
	report := &CoverageReport{
		Timestamp:       time.Now(),
		ProjectPath:     "test-project",
		OverallCoverage: 85.0,
		FileCoverage: []CoverageData{
			{
				FilePath:        "test1.go",
				PackageName:     "main",
				CoveragePercent: 90.0,
				LastUpdated:     time.Now(),
			},
			{
				FilePath:        "test2.go",
				PackageName:     "utils",
				CoveragePercent: 80.0,
				LastUpdated:     time.Now(),
			},
		},
	}
	
	// Store report
	err = store.Store(report)
	assert.NoError(t, err)
	
	// Verify file was created
	expectedFile := filepath.Join(tempDir, fmt.Sprintf("coverage-%s.json", report.Timestamp.Format("2006-01-02")))
	_, err = os.Stat(expectedFile)
	assert.NoError(t, err)
	
	// Verify file content
	data, err := os.ReadFile(expectedFile)
	require.NoError(t, err)
	
	var storedReport CoverageReport
	err = json.Unmarshal(data, &storedReport)
	require.NoError(t, err)
	
	assert.Equal(t, report.ProjectPath, storedReport.ProjectPath)
	assert.Equal(t, report.OverallCoverage, storedReport.OverallCoverage)
	assert.Len(t, storedReport.FileCoverage, 2)
	
	// Verify in-memory storage
	assert.Len(t, store.data, 2)
	assert.Contains(t, store.data, "test1.go")
	assert.Contains(t, store.data, "test2.go")
}

func TestFileWatcher(t *testing.T) {
	watcher := NewFileWatcher([]string{"."}, []string{"vendor"})
	
	// Test callback registration
	callbackCalled := false
	var receivedPath, receivedEvent string
	
	watcher.OnChange(func(path, event string) {
		callbackCalled = true
		receivedPath = path
		receivedEvent = event
	})
	
	assert.Len(t, watcher.callbacks, 1)
	
	// Test callback execution
	watcher.callbacks[0]("test.go", "modify")
	assert.True(t, callbackCalled)
	assert.Equal(t, "test.go", receivedPath)
	assert.Equal(t, "modify", receivedEvent)
	
	// Test start/stop
	err := watcher.Start()
	assert.NoError(t, err)
	
	err = watcher.Stop()
	assert.NoError(t, err)
	assert.False(t, watcher.running)
}

func TestAlertManager(t *testing.T) {
	config := DefaultMonitorConfig()
	alertManager := NewAlertManager(config)
	
	assert.NotNil(t, alertManager)
	assert.Equal(t, config, alertManager.config)
	
	// Test alert sending (simplified - just ensures no panic)
	alertManager.SendAlert("info", "Test message")
	alertManager.SendAlert("warning", "Warning message")
	alertManager.SendAlert("error", "Error message")
	alertManager.SendAlert("critical", "Critical message")
}

func TestMetricsCollector(t *testing.T) {
	metrics := NewMetricsCollector()
	
	assert.NotNil(t, metrics)
	assert.False(t, metrics.startTime.IsZero())
	assert.Equal(t, int64(0), metrics.monitoringCalls)
	assert.Equal(t, int64(0), metrics.coverageChecks)
	assert.Equal(t, int64(0), metrics.alertsSent)
	
	// Test thread-safe operations
	metrics.mu.Lock()
	metrics.coverageChecks++
	metrics.alertsSent++
	metrics.mu.Unlock()
	
	assert.Equal(t, int64(1), metrics.coverageChecks)
	assert.Equal(t, int64(1), metrics.alertsSent)
}

func TestAnalyzeCoverageForFile(t *testing.T) {
	// Create a temporary Go file
	tempDir, err := os.MkdirTemp("", "coverage-analysis-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.go")
	testContent := `package main

import "fmt"

// ExampleFunction demonstrates basic functionality
func ExampleFunction(x int) int {
	if x > 0 {
		return x * 2
	}
	return 0
}

func main() {
	result := ExampleFunction(5)
	fmt.Println(result)
}
`
	
	err = os.WriteFile(testFile, []byte(testContent), 0644)
	require.NoError(t, err)
	
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	// Analyze the file
	data, err := monitor.analyzeCoverageForFile(testFile)
	require.NoError(t, err)
	
	assert.Equal(t, testFile, data.FilePath)
	assert.Equal(t, "main", data.PackageName)
	assert.Greater(t, data.TotalLines, 0)
	assert.Greater(t, data.CoveredLines, 0)
	assert.Greater(t, data.CoveragePercent, 0.0)
	assert.Contains(t, data.FunctionCoverage, "ExampleFunction")
	assert.Contains(t, data.FunctionCoverage, "main")
	assert.NotNil(t, data.UncoveredLines)
	assert.NotNil(t, data.TestFiles)
	assert.NotNil(t, data.Dependencies)
	assert.False(t, data.LastUpdated.IsZero())
}

func TestFindGoFiles(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "find-go-files-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	// Create test files
	files := map[string]string{
		"main.go":          "package main",
		"utils.go":         "package utils",
		"main_test.go":     "package main", // Should be excluded
		"vendor/lib.go":    "package lib",  // Should be excluded
		"subdir/helper.go": "package helper",
		"README.md":        "# README", // Should be excluded
	}
	
	for filePath, content := range files {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)
		
		err := os.MkdirAll(dir, 0755)
		require.NoError(t, err)
		
		err = os.WriteFile(fullPath, []byte(content), 0644)
		require.NoError(t, err)
	}
	
	config := DefaultMonitorConfig()
	config.WatchPaths = []string{tempDir}
	monitor := NewCoverageMonitor(config)
	
	goFiles, err := monitor.findGoFiles()
	require.NoError(t, err)
	
	// Should find main.go, utils.go, and subdir/helper.go
	assert.Len(t, goFiles, 3)
	
	expectedFiles := []string{"main.go", "utils.go", "helper.go"}
	for _, expectedFile := range expectedFiles {
		found := false
		for _, actualFile := range goFiles {
			if strings.Contains(actualFile, expectedFile) {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find %s", expectedFile)
	}
	
	// Should not find test files, vendor files, or non-Go files
	excludedFiles := []string{"main_test.go", "lib.go", "README.md"}
	for _, excludedFile := range excludedFiles {
		found := false
		for _, actualFile := range goFiles {
			if strings.Contains(actualFile, excludedFile) {
				found = true
				break
			}
		}
		assert.False(t, found, "Should not find %s", excludedFile)
	}
}

func TestGenerateRecommendations(t *testing.T) {
	config := DefaultMonitorConfig()
	config.MinCoveragePercent = 80.0
	monitor := NewCoverageMonitor(config)
	
	report := &CoverageReport{
		OverallCoverage: 75.0, // Below minimum
		FileCoverage: []CoverageData{
			{FilePath: "good.go", CoveragePercent: 85.0},
			{FilePath: "bad.go", CoveragePercent: 60.0}, // Below 70%
			{FilePath: "poor.go", CoveragePercent: 50.0}, // Below 70%
		},
		RegressionAnalysis: RegressionAnalysis{
			HasRegression:   true,
			RegressionFiles: []string{"regressed.go", "another.go"},
		},
	}
	
	recommendations := monitor.generateRecommendations(report)
	
	assert.Greater(t, len(recommendations), 0)
	
	// Should recommend improving overall coverage
	foundOverallRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "overall coverage") && strings.Contains(rec, "75.00%") {
			foundOverallRecommendation = true
			break
		}
	}
	assert.True(t, foundOverallRecommendation)
	
	// Should recommend improving specific files
	foundFileRecommendations := 0
	for _, rec := range recommendations {
		if strings.Contains(rec, "bad.go") || strings.Contains(rec, "poor.go") {
			foundFileRecommendations++
		}
	}
	assert.Equal(t, 2, foundFileRecommendations)
	
	// Should recommend addressing regressions
	foundRegressionRecommendation := false
	for _, rec := range recommendations {
		if strings.Contains(rec, "regression") {
			foundRegressionRecommendation = true
			break
		}
	}
	assert.True(t, foundRegressionRecommendation)
}

func TestUpdateBaseline(t *testing.T) {
	// Create temporary directory with test files
	tempDir, err := os.MkdirTemp("", "baseline-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	testFile := filepath.Join(tempDir, "test.go")
	err = os.WriteFile(testFile, []byte("package main\nfunc main() {}"), 0644)
	require.NoError(t, err)
	
	config := DefaultMonitorConfig()
	config.WatchPaths = []string{tempDir}
	monitor := NewCoverageMonitor(config)
	
	// Update baseline
	err = monitor.UpdateBaseline()
	assert.NoError(t, err)
	
	// Verify baseline was set
	assert.Greater(t, len(monitor.regressionTracker.baseline), 0)
	assert.Contains(t, monitor.regressionTracker.baseline, testFile)
}

func TestTrendAnalysis(t *testing.T) {
	config := DefaultMonitorConfig()
	monitor := NewCoverageMonitor(config)
	
	trend := monitor.analyzeProjectTrend()
	
	assert.NotEmpty(t, trend.Trend)
	assert.Contains(t, []string{"improving", "stable", "declining"}, trend.Trend)
	assert.GreaterOrEqual(t, trend.TrendPercent, -100.0)
	assert.LessOrEqual(t, trend.TrendPercent, 100.0)
	assert.Greater(t, trend.DataPoints, 0)
	assert.NotEmpty(t, trend.AnalysisPeriod)
	assert.GreaterOrEqual(t, trend.PredictedCoverage, 0.0)
	assert.LessOrEqual(t, trend.PredictedCoverage, 100.0)
	assert.GreaterOrEqual(t, trend.Confidence, 0.0)
	assert.LessOrEqual(t, trend.Confidence, 1.0)
}

func TestCreateDefaultQualityGates(t *testing.T) {
	config := DefaultMonitorConfig()
	gates := createDefaultQualityGates(config)
	
	assert.Greater(t, len(gates), 0)
	
	// Check for expected default gates
	expectedGates := []string{"Minimum Coverage", "Regression Check", "Coverage Trend"}
	for _, expectedGate := range expectedGates {
		found := false
		for _, gate := range gates {
			if gate.Name == expectedGate {
				found = true
				assert.True(t, gate.Enabled || !config.EnableQualityGates)
				assert.NotEmpty(t, gate.Type)
				assert.NotEmpty(t, gate.Condition)
				assert.NotEmpty(t, gate.Description)
				break
			}
		}
		assert.True(t, found, "Should find default gate: %s", expectedGate)
	}
}

func TestCoverageReportSerialization(t *testing.T) {
	report := &CoverageReport{
		Timestamp:       time.Now(),
		ProjectPath:     "test-project",
		OverallCoverage: 85.5,
		PackageCoverage: map[string]float64{
			"main":  90.0,
			"utils": 80.0,
		},
		FileCoverage: []CoverageData{
			{
				FilePath:        "main.go",
				PackageName:     "main",
				CoveragePercent: 90.0,
			},
		},
		QualityGateResults: []QualityGateResult{
			{
				Gate: QualityGate{
					Name: "Test Gate",
					Type: "coverage",
				},
				Passed:      true,
				ActualValue: 85.5,
				Message:     "Test passed",
				Timestamp:   time.Now(),
			},
		},
		RegressionAnalysis: RegressionAnalysis{
			HasRegression:     false,
			RegressionPercent: 0.0,
			Severity:          "none",
		},
		TrendAnalysis: TrendAnalysis{
			Trend:             "stable",
			TrendPercent:      0.1,
			AnalysisPeriod:    "30 days",
			PredictedCoverage: 85.0,
			Confidence:        0.9,
		},
		Recommendations: []string{"Improve test coverage"},
		GenerationTime:  time.Millisecond * 100,
	}
	
	// Test JSON serialization
	data, err := json.Marshal(report)
	require.NoError(t, err)
	assert.Contains(t, string(data), "test-project")
	assert.Contains(t, string(data), "85.5")
	
	// Test JSON deserialization
	var decoded CoverageReport
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	
	assert.Equal(t, report.ProjectPath, decoded.ProjectPath)
	assert.Equal(t, report.OverallCoverage, decoded.OverallCoverage)
	assert.Equal(t, len(report.PackageCoverage), len(decoded.PackageCoverage))
	assert.Equal(t, len(report.FileCoverage), len(decoded.FileCoverage))
	assert.Equal(t, len(report.QualityGateResults), len(decoded.QualityGateResults))
	assert.Equal(t, report.RegressionAnalysis.HasRegression, decoded.RegressionAnalysis.HasRegression)
	assert.Equal(t, report.TrendAnalysis.Trend, decoded.TrendAnalysis.Trend)
	assert.Equal(t, len(report.Recommendations), len(decoded.Recommendations))
}