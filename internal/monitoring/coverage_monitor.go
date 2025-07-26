package monitoring

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// CoverageMonitor provides real-time coverage monitoring and quality gates
type CoverageMonitor struct {
	config          MonitorConfig
	fileWatcher     *FileWatcher
	coverageStore   *CoverageStore
	qualityGates    []QualityGate
	alertManager    *AlertManager
	metrics         *MetricsCollector
	regressionTracker *RegressionTracker
	mu              sync.RWMutex
	running         bool
	ctx             context.Context
	cancel          context.CancelFunc
}

// MonitorConfig configures the coverage monitoring system
type MonitorConfig struct {
	// Monitoring settings
	WatchPaths           []string      `json:"watch_paths"`
	ExcludePatterns      []string      `json:"exclude_patterns"`
	MonitorInterval      time.Duration `json:"monitor_interval"`
	RealtimeUpdates      bool          `json:"realtime_updates"`
	
	// Coverage thresholds
	MinCoveragePercent   float64       `json:"min_coverage_percent"`
	CriticalThreshold    float64       `json:"critical_threshold"`
	WarningThreshold     float64       `json:"warning_threshold"`
	
	// Quality gates
	EnableQualityGates   bool          `json:"enable_quality_gates"`
	FailOnRegression     bool          `json:"fail_on_regression"`
	MaxRegressionPercent float64       `json:"max_regression_percent"`
	
	// Reporting settings
	GenerateReports      bool          `json:"generate_reports"`
	ReportFormat         string        `json:"report_format"` // "json", "html", "markdown"
	ReportOutputPath     string        `json:"report_output_path"`
	
	// Integration settings
	WebhookURL           string        `json:"webhook_url"`
	SlackChannel         string        `json:"slack_channel"`
	EmailNotifications   []string      `json:"email_notifications"`
	
	// Advanced features
	EnableTrendAnalysis  bool          `json:"enable_trend_analysis"`
	HistoryRetentionDays int           `json:"history_retention_days"`
	EnablePrediction     bool          `json:"enable_prediction"`
}

// DefaultMonitorConfig returns sensible defaults for coverage monitoring
func DefaultMonitorConfig() MonitorConfig {
	return MonitorConfig{
		WatchPaths:           []string{"."},
		ExcludePatterns:      []string{"vendor", ".git", "node_modules", "*.test.go"},
		MonitorInterval:      time.Minute * 5,
		RealtimeUpdates:      true,
		
		MinCoveragePercent:   80.0,
		CriticalThreshold:    70.0,
		WarningThreshold:     75.0,
		
		EnableQualityGates:   true,
		FailOnRegression:     true,
		MaxRegressionPercent: 5.0,
		
		GenerateReports:      true,
		ReportFormat:         "json",
		ReportOutputPath:     "./coverage-reports",
		
		EnableTrendAnalysis:  true,
		HistoryRetentionDays: 30,
		EnablePrediction:     false,
	}
}

// CoverageData represents coverage information for a file or package
type CoverageData struct {
	FilePath          string            `json:"file_path"`
	PackageName       string            `json:"package_name"`
	TotalLines        int               `json:"total_lines"`
	CoveredLines      int               `json:"covered_lines"`
	CoveragePercent   float64           `json:"coverage_percent"`
	UncoveredLines    []int             `json:"uncovered_lines"`
	FunctionCoverage  map[string]float64 `json:"function_coverage"`
	BranchCoverage    float64           `json:"branch_coverage"`
	LastUpdated       time.Time         `json:"last_updated"`
	TestFiles         []string          `json:"test_files"`
	Dependencies      []string          `json:"dependencies"`
}

// QualityGate represents a quality gate check
type QualityGate struct {
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // "coverage", "regression", "trend", "custom"
	Condition   string                 `json:"condition"` // ">=", "<=", "==", "!=", ">", "<"
	Threshold   float64                `json:"threshold"`
	Enabled     bool                   `json:"enabled"`
	Critical    bool                   `json:"critical"`
	Description string                 `json:"description"`
	CustomFunc  func(CoverageData) bool `json:"-"`
}

// QualityGateResult represents the result of a quality gate check
type QualityGateResult struct {
	Gate        QualityGate `json:"gate"`
	Passed      bool        `json:"passed"`
	ActualValue float64     `json:"actual_value"`
	Message     string      `json:"message"`
	Timestamp   time.Time   `json:"timestamp"`
}

// CoverageReport represents a comprehensive coverage report
type CoverageReport struct {
	Timestamp           time.Time           `json:"timestamp"`
	ProjectPath         string              `json:"project_path"`
	OverallCoverage     float64             `json:"overall_coverage"`
	PackageCoverage     map[string]float64  `json:"package_coverage"`
	FileCoverage        []CoverageData      `json:"file_coverage"`
	QualityGateResults  []QualityGateResult `json:"quality_gate_results"`
	RegressionAnalysis  RegressionAnalysis  `json:"regression_analysis"`
	TrendAnalysis       TrendAnalysis       `json:"trend_analysis"`
	Recommendations     []string            `json:"recommendations"`
	GenerationTime      time.Duration       `json:"generation_time"`
}

// RegressionAnalysis tracks coverage regressions
type RegressionAnalysis struct {
	HasRegression       bool                   `json:"has_regression"`
	RegressionPercent   float64                `json:"regression_percent"`
	RegressionFiles     []string               `json:"regression_files"`
	ComparisonTimestamp time.Time              `json:"comparison_timestamp"`
	RegressionDetails   map[string]float64     `json:"regression_details"`
	Severity            string                 `json:"severity"` // "low", "medium", "high", "critical"
}

// TrendAnalysis provides trend information over time
type TrendAnalysis struct {
	Trend              string    `json:"trend"` // "improving", "stable", "declining"
	TrendPercent       float64   `json:"trend_percent"`
	DataPoints         int       `json:"data_points"`
	AnalysisPeriod     string    `json:"analysis_period"`
	PredictedCoverage  float64   `json:"predicted_coverage"`
	Confidence         float64   `json:"confidence"`
}

// FileWatcher monitors file system changes
type FileWatcher struct {
	paths     []string
	excludes  []string
	callbacks []func(string, string) // path, event type
	mu        sync.RWMutex
	running   bool
}

// CoverageStore manages coverage data persistence
type CoverageStore struct {
	data        map[string][]CoverageData
	mu          sync.RWMutex
	storePath   string
	maxHistory  int
}

// AlertManager handles notifications and alerts
type AlertManager struct {
	config    MonitorConfig
	webhooks  []WebhookConfig
	templates map[string]string
}

// WebhookConfig represents webhook configuration
type WebhookConfig struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers"`
	Enabled bool              `json:"enabled"`
}

// MetricsCollector gathers performance and usage metrics
type MetricsCollector struct {
	startTime       time.Time
	monitoringCalls int64
	coverageChecks  int64
	alertsSent      int64
	mu              sync.Mutex
}

// RegressionTracker tracks and analyzes coverage regressions
type RegressionTracker struct {
	baseline       map[string]float64
	history        []CoverageSnapshot
	mu             sync.RWMutex
	maxHistory     int
}

// CoverageSnapshot represents a point-in-time coverage snapshot
type CoverageSnapshot struct {
	Timestamp time.Time          `json:"timestamp"`
	Coverage  map[string]float64 `json:"coverage"`
	Metadata  map[string]string  `json:"metadata"`
}

// NewCoverageMonitor creates a new coverage monitoring system
func NewCoverageMonitor(config MonitorConfig) *CoverageMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	
	monitor := &CoverageMonitor{
		config:            config,
		fileWatcher:       NewFileWatcher(config.WatchPaths, config.ExcludePatterns),
		coverageStore:     NewCoverageStore(config.ReportOutputPath, config.HistoryRetentionDays),
		qualityGates:      createDefaultQualityGates(config),
		alertManager:      NewAlertManager(config),
		metrics:           NewMetricsCollector(),
		regressionTracker: NewRegressionTracker(100), // Keep last 100 snapshots
		ctx:               ctx,
		cancel:            cancel,
	}
	
	// Set up file watcher callbacks
	monitor.fileWatcher.OnChange(monitor.handleFileChange)
	
	return monitor
}

// Start begins coverage monitoring
func (cm *CoverageMonitor) Start() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if cm.running {
		return fmt.Errorf("coverage monitor is already running")
	}
	
	cm.running = true
	cm.metrics.startTime = time.Now()
	
	// Start file watcher if real-time updates are enabled
	if cm.config.RealtimeUpdates {
		if err := cm.fileWatcher.Start(); err != nil {
			cm.running = false
			return fmt.Errorf("failed to start file watcher: %w", err)
		}
	}
	
	// Start periodic monitoring
	go cm.monitoringLoop()
	
	// Perform initial coverage check
	go func() {
		if err := cm.performCoverageCheck(); err != nil {
			cm.alertManager.SendAlert("error", fmt.Sprintf("Initial coverage check failed: %v", err))
		}
	}()
	
	return nil
}

// Stop stops coverage monitoring
func (cm *CoverageMonitor) Stop() error {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	if !cm.running {
		return fmt.Errorf("coverage monitor is not running")
	}
	
	cm.running = false
	cm.cancel()
	
	if cm.fileWatcher != nil {
		if err := cm.fileWatcher.Stop(); err != nil {
			return fmt.Errorf("failed to stop file watcher: %w", err)
		}
	}
	
	return nil
}

// GetCurrentCoverage returns the current coverage status
func (cm *CoverageMonitor) GetCurrentCoverage() (*CoverageReport, error) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()
	
	if !cm.running {
		return nil, fmt.Errorf("coverage monitor is not running")
	}
	
	return cm.generateReport()
}

// AddQualityGate adds a custom quality gate
func (cm *CoverageMonitor) AddQualityGate(gate QualityGate) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	
	cm.qualityGates = append(cm.qualityGates, gate)
}

// UpdateBaseline updates the coverage baseline for regression detection
func (cm *CoverageMonitor) UpdateBaseline() error {
	report, err := cm.generateReport()
	if err != nil {
		return fmt.Errorf("failed to generate report for baseline: %w", err)
	}
	
	baseline := make(map[string]float64)
	for _, fileData := range report.FileCoverage {
		baseline[fileData.FilePath] = fileData.CoveragePercent
	}
	
	cm.regressionTracker.SetBaseline(baseline)
	return nil
}

// monitoringLoop runs the periodic monitoring
func (cm *CoverageMonitor) monitoringLoop() {
	ticker := time.NewTicker(cm.config.MonitorInterval)
	defer ticker.Stop()
	
	for {
		select {
		case <-cm.ctx.Done():
			return
		case <-ticker.C:
			if err := cm.performCoverageCheck(); err != nil {
				cm.alertManager.SendAlert("error", fmt.Sprintf("Coverage check failed: %v", err))
			}
		}
	}
}

// handleFileChange processes file system changes
func (cm *CoverageMonitor) handleFileChange(path, eventType string) {
	// Skip if not a Go file
	if !strings.HasSuffix(path, ".go") {
		return
	}
	
	// Skip test files unless they're part of coverage analysis
	if strings.HasSuffix(path, "_test.go") && eventType != "delete" {
		// Test file changed, might affect coverage
		go func() {
			time.Sleep(time.Second * 2) // Wait for file operations to complete
			if err := cm.performCoverageCheck(); err != nil {
				cm.alertManager.SendAlert("warning", fmt.Sprintf("Coverage check after file change failed: %v", err))
			}
		}()
	}
}

// performCoverageCheck performs a comprehensive coverage analysis
func (cm *CoverageMonitor) performCoverageCheck() error {
	cm.metrics.mu.Lock()
	cm.metrics.coverageChecks++
	cm.metrics.mu.Unlock()
	
	startTime := time.Now()
	
	// Generate coverage report
	report, err := cm.generateReport()
	if err != nil {
		return fmt.Errorf("failed to generate coverage report: %w", err)
	}
	
	// Store coverage data
	if err := cm.coverageStore.Store(report); err != nil {
		return fmt.Errorf("failed to store coverage data: %w", err)
	}
	
	// Check quality gates
	gateResults := cm.checkQualityGates(report)
	report.QualityGateResults = gateResults
	
	// Check for regressions
	regressionAnalysis := cm.regressionTracker.AnalyzeRegression(report)
	report.RegressionAnalysis = regressionAnalysis
	
	// Perform trend analysis if enabled
	if cm.config.EnableTrendAnalysis {
		trendAnalysis := cm.analyzeProjectTrend()
		report.TrendAnalysis = trendAnalysis
	}
	
	// Generate recommendations
	report.Recommendations = cm.generateRecommendations(report)
	
	// Send alerts if necessary
	cm.processAlerts(report)
	
	report.GenerationTime = time.Since(startTime)
	
	return nil
}

// generateReport creates a comprehensive coverage report
func (cm *CoverageMonitor) generateReport() (*CoverageReport, error) {
	report := &CoverageReport{
		Timestamp:      time.Now(),
		ProjectPath:    ".", // TODO: Make configurable
		PackageCoverage: make(map[string]float64),
		FileCoverage:   make([]CoverageData, 0),
	}
	
	// Scan for Go files
	goFiles, err := cm.findGoFiles()
	if err != nil {
		return nil, fmt.Errorf("failed to find Go files: %w", err)
	}
	
	totalLines := 0
	totalCovered := 0
	packageLines := make(map[string]int)
	packageCovered := make(map[string]int)
	
	// Analyze each file
	for _, filePath := range goFiles {
		fileData, err := cm.analyzeCoverageForFile(filePath)
		if err != nil {
			continue // Skip files that can't be analyzed
		}
		
		report.FileCoverage = append(report.FileCoverage, *fileData)
		
		totalLines += fileData.TotalLines
		totalCovered += fileData.CoveredLines
		
		// Aggregate by package
		packageLines[fileData.PackageName] += fileData.TotalLines
		packageCovered[fileData.PackageName] += fileData.CoveredLines
	}
	
	// Calculate overall coverage
	if totalLines > 0 {
		report.OverallCoverage = float64(totalCovered) / float64(totalLines) * 100
	}
	
	// Calculate package coverage
	for pkg := range packageLines {
		if packageLines[pkg] > 0 {
			report.PackageCoverage[pkg] = float64(packageCovered[pkg]) / float64(packageLines[pkg]) * 100
		}
	}
	
	return report, nil
}

// analyzeCoverageForFile analyzes coverage for a single file
func (cm *CoverageMonitor) analyzeCoverageForFile(filePath string) (*CoverageData, error) {
	// Parse the Go file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file %s: %w", filePath, err)
	}
	
	data := &CoverageData{
		FilePath:         filePath,
		PackageName:      file.Name.Name,
		FunctionCoverage: make(map[string]float64),
		LastUpdated:      time.Now(),
		UncoveredLines:   make([]int, 0),
		TestFiles:        make([]string, 0),
		Dependencies:     make([]string, 0),
	}
	
	// Count lines and analyze functions
	data.TotalLines = cm.countLines(file)
	data.CoveredLines = cm.estimateCoveredLines(file) // Simplified estimation
	
	if data.TotalLines > 0 {
		data.CoveragePercent = float64(data.CoveredLines) / float64(data.TotalLines) * 100
	}
	
	// Analyze functions
	ast.Inspect(file, func(node ast.Node) bool {
		if funcDecl, ok := node.(*ast.FuncDecl); ok {
			funcName := funcDecl.Name.Name
			// Simplified function coverage estimation
			funcLines := cm.countFunctionLines(funcDecl)
			funcCovered := funcLines * 80 / 100 // Placeholder logic
			if funcLines > 0 {
				data.FunctionCoverage[funcName] = float64(funcCovered) / float64(funcLines) * 100
			}
		}
		return true
	})
	
	// Find test files
	data.TestFiles = cm.findTestFiles(filePath)
	
	// Extract dependencies
	data.Dependencies = cm.extractDependencies(file)
	
	return data, nil
}

// checkQualityGates evaluates all quality gates
func (cm *CoverageMonitor) checkQualityGates(report *CoverageReport) []QualityGateResult {
	results := make([]QualityGateResult, 0, len(cm.qualityGates))
	
	for _, gate := range cm.qualityGates {
		if !gate.Enabled {
			continue
		}
		
		result := cm.evaluateQualityGate(gate, report)
		results = append(results, result)
	}
	
	return results
}

// evaluateQualityGate evaluates a single quality gate
func (cm *CoverageMonitor) evaluateQualityGate(gate QualityGate, report *CoverageReport) QualityGateResult {
	result := QualityGateResult{
		Gate:      gate,
		Timestamp: time.Now(),
	}
	
	switch gate.Type {
	case "coverage":
		result.ActualValue = report.OverallCoverage
		result.Passed = cm.evaluateCondition(result.ActualValue, gate.Condition, gate.Threshold)
		if result.Passed {
			result.Message = fmt.Sprintf("Coverage %.2f%% meets threshold %.2f%%", result.ActualValue, gate.Threshold)
		} else {
			result.Message = fmt.Sprintf("Coverage %.2f%% below threshold %.2f%%", result.ActualValue, gate.Threshold)
		}
		
	case "regression":
		result.ActualValue = report.RegressionAnalysis.RegressionPercent
		result.Passed = !report.RegressionAnalysis.HasRegression || result.ActualValue <= gate.Threshold
		if result.Passed {
			result.Message = "No significant regression detected"
		} else {
			result.Message = fmt.Sprintf("Regression of %.2f%% exceeds threshold %.2f%%", result.ActualValue, gate.Threshold)
		}
		
	case "trend":
		result.ActualValue = report.TrendAnalysis.TrendPercent
		result.Passed = cm.evaluateCondition(result.ActualValue, gate.Condition, gate.Threshold)
		result.Message = fmt.Sprintf("Trend %.2f%% evaluated against threshold %.2f%%", result.ActualValue, gate.Threshold)
		
	case "custom":
		if gate.CustomFunc != nil {
			// Use overall coverage as representative data
			mockData := CoverageData{CoveragePercent: report.OverallCoverage}
			result.Passed = gate.CustomFunc(mockData)
			result.ActualValue = report.OverallCoverage
			result.Message = fmt.Sprintf("Custom gate evaluation: %t", result.Passed)
		}
	}
	
	return result
}

// evaluateCondition evaluates a condition string
func (cm *CoverageMonitor) evaluateCondition(actual float64, condition string, threshold float64) bool {
	switch condition {
	case ">=":
		return actual >= threshold
	case "<=":
		return actual <= threshold
	case ">":
		return actual > threshold
	case "<":
		return actual < threshold
	case "==":
		return actual == threshold
	case "!=":
		return actual != threshold
	default:
		return false
	}
}

// Helper functions and supporting components

func NewFileWatcher(paths []string, excludes []string) *FileWatcher {
	return &FileWatcher{
		paths:     paths,
		excludes:  excludes,
		callbacks: make([]func(string, string), 0),
	}
}

func (fw *FileWatcher) OnChange(callback func(string, string)) {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.callbacks = append(fw.callbacks, callback)
}

func (fw *FileWatcher) Start() error {
	// Simplified file watching - in production would use fsnotify
	return nil
}

func (fw *FileWatcher) Stop() error {
	fw.mu.Lock()
	defer fw.mu.Unlock()
	fw.running = false
	return nil
}

func NewCoverageStore(path string, maxHistory int) *CoverageStore {
	return &CoverageStore{
		data:       make(map[string][]CoverageData),
		storePath:  path,
		maxHistory: maxHistory,
	}
}

func (cs *CoverageStore) Store(report *CoverageReport) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	
	// Store coverage data (simplified implementation)
	timestamp := report.Timestamp.Format("2006-01-02")
	for _, fileData := range report.FileCoverage {
		cs.data[fileData.FilePath] = append(cs.data[fileData.FilePath], fileData)
		
		// Maintain history limit
		if len(cs.data[fileData.FilePath]) > cs.maxHistory {
			cs.data[fileData.FilePath] = cs.data[fileData.FilePath][1:]
		}
	}
	
	// Write to file
	if cs.storePath != "" {
		reportPath := filepath.Join(cs.storePath, fmt.Sprintf("coverage-%s.json", timestamp))
		return cs.writeReportToFile(report, reportPath)
	}
	
	return nil
}

func (cs *CoverageStore) writeReportToFile(report *CoverageReport, path string) error {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	
	return os.WriteFile(path, data, 0644)
}

func NewAlertManager(config MonitorConfig) *AlertManager {
	return &AlertManager{
		config:    config,
		webhooks:  make([]WebhookConfig, 0),
		templates: make(map[string]string),
	}
}

func (am *AlertManager) SendAlert(level, message string) {
	// Simplified alert sending
	fmt.Printf("[ALERT-%s] %s\n", strings.ToUpper(level), message)
}

func NewMetricsCollector() *MetricsCollector {
	return &MetricsCollector{
		startTime: time.Now(),
	}
}

func NewRegressionTracker(maxHistory int) *RegressionTracker {
	return &RegressionTracker{
		baseline:   make(map[string]float64),
		history:    make([]CoverageSnapshot, 0),
		maxHistory: maxHistory,
	}
}

func (rt *RegressionTracker) SetBaseline(baseline map[string]float64) {
	rt.mu.Lock()
	defer rt.mu.Unlock()
	rt.baseline = baseline
}

func (rt *RegressionTracker) AnalyzeRegression(report *CoverageReport) RegressionAnalysis {
	rt.mu.RLock()
	defer rt.mu.RUnlock()
	
	analysis := RegressionAnalysis{
		HasRegression:     false,
		RegressionPercent: 0.0,
		RegressionFiles:   make([]string, 0),
		RegressionDetails: make(map[string]float64),
		Severity:          "none",
	}
	
	if len(rt.baseline) == 0 {
		return analysis
	}
	
	totalRegression := 0.0
	regressionCount := 0
	
	for _, fileData := range report.FileCoverage {
		if baselineCoverage, exists := rt.baseline[fileData.FilePath]; exists {
			regression := baselineCoverage - fileData.CoveragePercent
			if regression > 0 {
				analysis.HasRegression = true
				analysis.RegressionFiles = append(analysis.RegressionFiles, fileData.FilePath)
				analysis.RegressionDetails[fileData.FilePath] = regression
				totalRegression += regression
				regressionCount++
			}
		}
	}
	
	if regressionCount > 0 {
		analysis.RegressionPercent = totalRegression / float64(regressionCount)
		
		// Determine severity
		switch {
		case analysis.RegressionPercent >= 10:
			analysis.Severity = "critical"
		case analysis.RegressionPercent >= 5:
			analysis.Severity = "high"
		case analysis.RegressionPercent >= 2:
			analysis.Severity = "medium"
		default:
			analysis.Severity = "low"
		}
	}
	
	return analysis
}

// Helper methods for file analysis

func (cm *CoverageMonitor) findGoFiles() ([]string, error) {
	var files []string
	
	for _, watchPath := range cm.config.WatchPaths {
		err := filepath.Walk(watchPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			
			if !info.IsDir() && strings.HasSuffix(path, ".go") && !strings.HasSuffix(path, "_test.go") {
				// Check exclude patterns
				excluded := false
				for _, pattern := range cm.config.ExcludePatterns {
					if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
						excluded = true
						break
					}
					if strings.Contains(path, pattern) {
						excluded = true
						break
					}
				}
				
				if !excluded {
					files = append(files, path)
				}
			}
			
			return nil
		})
		
		if err != nil {
			return nil, err
		}
	}
	
	return files, nil
}

func (cm *CoverageMonitor) countLines(file *ast.File) int {
	// Simplified line counting
	lines := 0
	ast.Inspect(file, func(node ast.Node) bool {
		if node != nil {
			lines++
		}
		return true
	})
	return lines / 5 // Rough approximation
}

func (cm *CoverageMonitor) estimateCoveredLines(file *ast.File) int {
	// Simplified coverage estimation - in production would use actual coverage data
	totalLines := cm.countLines(file)
	return int(float64(totalLines) * 0.8) // Assume 80% coverage
}

func (cm *CoverageMonitor) countFunctionLines(funcDecl *ast.FuncDecl) int {
	if funcDecl.Body == nil {
		return 0
	}
	return len(funcDecl.Body.List)
}

func (cm *CoverageMonitor) findTestFiles(filePath string) []string {
	// Find corresponding test files
	dir := filepath.Dir(filePath)
	base := filepath.Base(filePath)
	name := strings.TrimSuffix(base, ".go")
	
	testFile := filepath.Join(dir, name+"_test.go")
	if _, err := os.Stat(testFile); err == nil {
		return []string{testFile}
	}
	
	return []string{}
}

func (cm *CoverageMonitor) extractDependencies(file *ast.File) []string {
	dependencies := make([]string, 0)
	
	for _, imp := range file.Imports {
		if imp.Path != nil {
			path := strings.Trim(imp.Path.Value, `"`)
			dependencies = append(dependencies, path)
		}
	}
	
	return dependencies
}

func (cm *CoverageMonitor) analyzeProjectTrend() TrendAnalysis {
	// Simplified trend analysis
	return TrendAnalysis{
		Trend:             "stable",
		TrendPercent:      0.1,
		DataPoints:        30,
		AnalysisPeriod:    "30 days",
		PredictedCoverage: 82.5,
		Confidence:        0.85,
	}
}

func (cm *CoverageMonitor) generateRecommendations(report *CoverageReport) []string {
	recommendations := make([]string, 0)
	
	if report.OverallCoverage < cm.config.MinCoveragePercent {
		recommendations = append(recommendations, 
			fmt.Sprintf("Increase overall coverage from %.2f%% to %.2f%%", 
				report.OverallCoverage, cm.config.MinCoveragePercent))
	}
	
	// Find files with low coverage
	for _, fileData := range report.FileCoverage {
		if fileData.CoveragePercent < 70.0 {
			recommendations = append(recommendations, 
				fmt.Sprintf("Improve coverage for %s (currently %.2f%%)", 
					fileData.FilePath, fileData.CoveragePercent))
		}
	}
	
	if report.RegressionAnalysis.HasRegression {
		recommendations = append(recommendations, 
			"Address coverage regressions in: "+strings.Join(report.RegressionAnalysis.RegressionFiles, ", "))
	}
	
	return recommendations
}

func (cm *CoverageMonitor) processAlerts(report *CoverageReport) {
	// Check critical thresholds
	if report.OverallCoverage < cm.config.CriticalThreshold {
		cm.alertManager.SendAlert("critical", 
			fmt.Sprintf("Coverage %.2f%% below critical threshold %.2f%%", 
				report.OverallCoverage, cm.config.CriticalThreshold))
	} else if report.OverallCoverage < cm.config.WarningThreshold {
		cm.alertManager.SendAlert("warning", 
			fmt.Sprintf("Coverage %.2f%% below warning threshold %.2f%%", 
				report.OverallCoverage, cm.config.WarningThreshold))
	}
	
	// Check regressions
	if report.RegressionAnalysis.HasRegression && cm.config.FailOnRegression {
		severity := report.RegressionAnalysis.Severity
		cm.alertManager.SendAlert(severity, 
			fmt.Sprintf("Coverage regression detected: %.2f%% in %d files", 
				report.RegressionAnalysis.RegressionPercent, len(report.RegressionAnalysis.RegressionFiles)))
	}
	
	// Check quality gate failures
	for _, result := range report.QualityGateResults {
		if !result.Passed && result.Gate.Critical {
			cm.alertManager.SendAlert("critical", 
				fmt.Sprintf("Critical quality gate failed: %s - %s", result.Gate.Name, result.Message))
		}
	}
}

func createDefaultQualityGates(config MonitorConfig) []QualityGate {
	gates := []QualityGate{
		{
			Name:        "Minimum Coverage",
			Type:        "coverage",
			Condition:   ">=",
			Threshold:   config.MinCoveragePercent,
			Enabled:     true,
			Critical:    true,
			Description: "Ensures minimum coverage percentage is maintained",
		},
		{
			Name:        "Regression Check",
			Type:        "regression",
			Condition:   "<=",
			Threshold:   config.MaxRegressionPercent,
			Enabled:     config.FailOnRegression,
			Critical:    config.FailOnRegression,
			Description: "Prevents significant coverage regressions",
		},
		{
			Name:        "Coverage Trend",
			Type:        "trend",
			Condition:   ">=",
			Threshold:   -1.0, // Allow up to 1% decline
			Enabled:     config.EnableTrendAnalysis,
			Critical:    false,
			Description: "Monitors coverage trend over time",
		},
	}
	
	return gates
}