package performance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/francknouama/go-starter/pkg/types"
)

// BenchmarkResult holds performance metrics for a single benchmark run
type BenchmarkResult struct {
	Blueprint     string        `json:"blueprint"`
	Duration      time.Duration `json:"duration"`
	FilesCreated  int           `json:"files_created"`
	MemoryUsed    int64         `json:"memory_used_bytes"`
	TempDirSize   int64         `json:"temp_dir_size_bytes"`
	Platform      string        `json:"platform"`
	GoVersion     string        `json:"go_version"`
	Timestamp     time.Time     `json:"timestamp"`
	Success       bool          `json:"success"`
	ErrorMessage  string        `json:"error_message,omitempty"`
}

// BenchmarkSuite holds results for multiple benchmarks
type BenchmarkSuite struct {
	Results     []BenchmarkResult `json:"results"`
	Summary     BenchmarkSummary  `json:"summary"`
	Environment SystemInfo        `json:"environment"`
}

// BenchmarkSummary provides aggregated metrics
type BenchmarkSummary struct {
	TotalRuns         int           `json:"total_runs"`
	SuccessfulRuns    int           `json:"successful_runs"`
	FailedRuns        int           `json:"failed_runs"`
	AverageDuration   time.Duration `json:"average_duration"`
	MinDuration       time.Duration `json:"min_duration"`
	MaxDuration       time.Duration `json:"max_duration"`
	TotalFilesCreated int           `json:"total_files_created"`
	TotalMemoryUsed   int64         `json:"total_memory_used_bytes"`
	PerformanceGrade  string        `json:"performance_grade"`
}

// SystemInfo captures environment details
type SystemInfo struct {
	OS               string    `json:"os"`
	Architecture     string    `json:"architecture"`
	NumCPU           int       `json:"num_cpu"`
	GoVersion        string    `json:"go_version"`
	MemoryAvailable  int64     `json:"memory_available_bytes"`
	TestStartTime    time.Time `json:"test_start_time"`
}

var (
	benchmarkSuite = &BenchmarkSuite{
		Results: make([]BenchmarkResult, 0),
		Environment: SystemInfo{
			OS:            runtime.GOOS,
			Architecture:  runtime.GOARCH,
			NumCPU:        runtime.NumCPU(),
			GoVersion:     runtime.Version(),
			TestStartTime: time.Now(),
		},
	}
	resultsMutex sync.Mutex
)

// getBlueprintID determines the blueprint ID from config
func getBlueprintID(config types.ProjectConfig) string {
	if config.Variables != nil {
		if blueprintID, exists := config.Variables["blueprint_id"]; exists && blueprintID != "" {
			return blueprintID
		}
	}
	
	if config.Architecture != "" && config.Architecture != "standard" {
		return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
	}
	return config.Type
}

// generateJSONReport creates a comprehensive JSON performance report
func generateJSONReport() error {
	reportPath := filepath.Join("performance_reports", fmt.Sprintf("performance_report_%s.json", 
		time.Now().Format("20060102_150405")))
	
	if err := os.MkdirAll(filepath.Dir(reportPath), 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	data, err := json.MarshalIndent(benchmarkSuite, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal benchmark results: %w", err)
	}

	if err := os.WriteFile(reportPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write JSON report: %w", err)
	}

	fmt.Printf("JSON performance report generated: %s\n", reportPath)
	return nil
}

// generateMarkdownReport creates a human-readable Markdown performance report
func generateMarkdownReport() error {
	reportPath := filepath.Join("performance_reports", fmt.Sprintf("performance_report_%s.md", 
		time.Now().Format("20060102_150405")))
	
	if err := os.MkdirAll(filepath.Dir(reportPath), 0755); err != nil {
		return fmt.Errorf("failed to create report directory: %w", err)
	}

	var report strings.Builder
	
	// Header
	report.WriteString("# Go-Starter Performance Report\n\n")
	report.WriteString(fmt.Sprintf("**Generated:** %s\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	report.WriteString(fmt.Sprintf("**Platform:** %s/%s\n", benchmarkSuite.Environment.OS, benchmarkSuite.Environment.Architecture))
	report.WriteString(fmt.Sprintf("**Go Version:** %s\n", benchmarkSuite.Environment.GoVersion))
	report.WriteString(fmt.Sprintf("**CPU Cores:** %d\n\n", benchmarkSuite.Environment.NumCPU))

	// Executive Summary
	report.WriteString("## Executive Summary\n\n")
	summary := benchmarkSuite.Summary
	
	report.WriteString(fmt.Sprintf("- **Performance Grade:** %s\n", summary.PerformanceGrade))
	report.WriteString(fmt.Sprintf("- **Average Generation Time:** %v\n", summary.AverageDuration))
	report.WriteString(fmt.Sprintf("- **Success Rate:** %.1f%% (%d/%d)\n", 
		float64(summary.SuccessfulRuns)/float64(summary.TotalRuns)*100, 
		summary.SuccessfulRuns, summary.TotalRuns))
	report.WriteString(fmt.Sprintf("- **Total Files Generated:** %d\n", summary.TotalFilesCreated))
	report.WriteString(fmt.Sprintf("- **Total Memory Used:** %.2f MB\n\n", 
		float64(summary.TotalMemoryUsed)/1024/1024))

	// Performance Goals Assessment
	report.WriteString("## Performance Goals Assessment\n\n")
	
	// Goal 1: <2s generation time
	if summary.AverageDuration <= 2*time.Second {
		report.WriteString("✅ **Generation Time Goal:** PASSED - Average generation time %.3fs is under 2s target\n")
	} else {
		report.WriteString("❌ **Generation Time Goal:** FAILED - Average generation time %.3fs exceeds 2s target\n")
	}
	
	// Goal 2: 95% success rate
	successRate := float64(summary.SuccessfulRuns) / float64(summary.TotalRuns) * 100
	if successRate >= 95.0 {
		report.WriteString(fmt.Sprintf("✅ **Reliability Goal:** PASSED - Success rate %.1f%% meets 95%% target\n", successRate))
	} else {
		report.WriteString(fmt.Sprintf("❌ **Reliability Goal:** FAILED - Success rate %.1f%% below 95%% target\n", successRate))
	}
	
	// Goal 3: Memory usage
	avgMemoryMB := float64(summary.TotalMemoryUsed) / float64(summary.SuccessfulRuns) / 1024 / 1024
	if avgMemoryMB <= 50.0 {
		report.WriteString(fmt.Sprintf("✅ **Memory Goal:** PASSED - Average memory usage %.1fMB is under 50MB target\n\n", avgMemoryMB))
	} else {
		report.WriteString(fmt.Sprintf("❌ **Memory Goal:** FAILED - Average memory usage %.1fMB exceeds 50MB target\n\n", avgMemoryMB))
	}

	// Detailed Results Table
	report.WriteString("## Detailed Results by Blueprint\n\n")
	report.WriteString("| Blueprint | Duration | Files | Memory (MB) | Platform | Status |\n")
	report.WriteString("|-----------|----------|-------|-------------|----------|--------|\n")
	
	for _, result := range benchmarkSuite.Results {
		status := "✅ Success"
		if !result.Success {
			status = "❌ Failed"
		}
		
		memoryMB := float64(result.MemoryUsed) / 1024 / 1024
		report.WriteString(fmt.Sprintf("| %s | %v | %d | %.2f | %s | %s |\n",
			result.Blueprint,
			result.Duration.Round(time.Millisecond),
			result.FilesCreated,
			memoryMB,
			result.Platform,
			status))
	}
	report.WriteString("\n")

	// Performance Analysis by Blueprint Type
	report.WriteString("## Performance Analysis by Blueprint Type\n\n")
	blueprintStats := calculateBlueprintStats()
	
	for blueprintType, stats := range blueprintStats {
		report.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(blueprintType)))
		report.WriteString(fmt.Sprintf("- **Average Duration:** %v\n", stats.AvgDuration))
		report.WriteString(fmt.Sprintf("- **Min/Max Duration:** %v / %v\n", stats.MinDuration, stats.MaxDuration))
		report.WriteString(fmt.Sprintf("- **Average Files:** %.1f\n", stats.AvgFiles))
		report.WriteString(fmt.Sprintf("- **Average Memory:** %.2f MB\n", stats.AvgMemoryMB))
		report.WriteString(fmt.Sprintf("- **Success Rate:** %.1f%%\n\n", stats.SuccessRate))
	}

	// Performance Trends and Insights
	report.WriteString("## Performance Insights\n\n")
	generateInsights(&report, blueprintStats)

	// Optimization Recommendations
	report.WriteString("## Optimization Recommendations\n\n")
	generateOptimizationRecommendations(&report, benchmarkSuite.Summary, blueprintStats)

	// Cross-Platform Compatibility
	report.WriteString("## Cross-Platform Performance\n\n")
	generateCrossPlatformAnalysis(&report)

	// Technical Details
	report.WriteString("## Technical Details\n\n")
	report.WriteString("### Test Environment\n\n")
	report.WriteString(fmt.Sprintf("- **Operating System:** %s\n", benchmarkSuite.Environment.OS))
	report.WriteString(fmt.Sprintf("- **Architecture:** %s\n", benchmarkSuite.Environment.Architecture))
	report.WriteString(fmt.Sprintf("- **CPU Cores:** %d\n", benchmarkSuite.Environment.NumCPU))
	report.WriteString(fmt.Sprintf("- **Go Version:** %s\n", benchmarkSuite.Environment.GoVersion))
	report.WriteString(fmt.Sprintf("- **Test Duration:** %v\n\n", time.Since(benchmarkSuite.Environment.TestStartTime)))

	// Footer
	report.WriteString("---\n\n")
	report.WriteString("*Generated by go-starter Performance & Reliability Testing Suite*\n")

	if err := os.WriteFile(reportPath, []byte(report.String()), 0644); err != nil {
		return fmt.Errorf("failed to write Markdown report: %w", err)
	}

	fmt.Printf("Markdown performance report generated: %s\n", reportPath)
	return nil
}

// BlueprintStats holds aggregated statistics for a blueprint type
type BlueprintStats struct {
	Type         string
	Count        int
	AvgDuration  time.Duration
	MinDuration  time.Duration
	MaxDuration  time.Duration
	AvgFiles     float64
	AvgMemoryMB  float64
	SuccessRate  float64
	TotalMemory  int64
	TotalFiles   int
	Successes    int
}

// calculateBlueprintStats aggregates statistics by blueprint type
func calculateBlueprintStats() map[string]*BlueprintStats {
	stats := make(map[string]*BlueprintStats)
	
	for _, result := range benchmarkSuite.Results {
		blueprintType := extractBlueprintType(result.Blueprint)
		
		if _, exists := stats[blueprintType]; !exists {
			stats[blueprintType] = &BlueprintStats{
				Type:        blueprintType,
				MinDuration: time.Hour, // Start with large value
			}
		}
		
		stat := stats[blueprintType]
		stat.Count++
		
		if result.Success {
			stat.Successes++
			stat.AvgDuration += result.Duration
			stat.TotalMemory += result.MemoryUsed
			stat.TotalFiles += result.FilesCreated
			
			if result.Duration < stat.MinDuration {
				stat.MinDuration = result.Duration
			}
			if result.Duration > stat.MaxDuration {
				stat.MaxDuration = result.Duration
			}
		}
	}
	
	// Calculate averages
	for _, stat := range stats {
		if stat.Successes > 0 {
			stat.AvgDuration = time.Duration(int64(stat.AvgDuration) / int64(stat.Successes))
			stat.AvgFiles = float64(stat.TotalFiles) / float64(stat.Successes)
			stat.AvgMemoryMB = float64(stat.TotalMemory) / float64(stat.Successes) / 1024 / 1024
		}
		stat.SuccessRate = float64(stat.Successes) / float64(stat.Count) * 100
	}
	
	return stats
}

// extractBlueprintType extracts the base type from a blueprint name
func extractBlueprintType(blueprint string) string {
	// Remove suffixes like -memory, -fileio, -parallel
	parts := strings.Split(blueprint, "-")
	if len(parts) >= 2 {
		// Handle cases like "cli-simple", "web-api-clean"
		if parts[0] == "cli" || parts[0] == "web" || parts[0] == "lambda" {
			if parts[0] == "web" && len(parts) > 1 && parts[1] == "api" {
				return "web-api"
			}
			return parts[0]
		}
	}
	return blueprint
}

// generateInsights creates performance insights based on the results
func generateInsights(report *strings.Builder, stats map[string]*BlueprintStats) {
	// Find fastest and slowest blueprints
	var fastest, slowest *BlueprintStats
	for _, stat := range stats {
		if stat.Successes == 0 {
			continue
		}
		if fastest == nil || stat.AvgDuration < fastest.AvgDuration {
			fastest = stat
		}
		if slowest == nil || stat.AvgDuration > slowest.AvgDuration {
			slowest = stat
		}
	}
	
	if fastest != nil && slowest != nil {
		report.WriteString(fmt.Sprintf("- **Fastest Blueprint:** %s (avg: %v)\n", fastest.Type, fastest.AvgDuration))
		report.WriteString(fmt.Sprintf("- **Slowest Blueprint:** %s (avg: %v)\n", slowest.Type, slowest.AvgDuration))
		
		if slowest.AvgDuration > fastest.AvgDuration {
			ratio := float64(slowest.AvgDuration) / float64(fastest.AvgDuration)
			report.WriteString(fmt.Sprintf("- **Performance Ratio:** %s is %.1fx slower than %s\n", 
				slowest.Type, ratio, fastest.Type))
		}
	}
	
	// Memory usage insights
	var highestMemory, lowestMemory *BlueprintStats
	for _, stat := range stats {
		if stat.Successes == 0 {
			continue
		}
		if highestMemory == nil || stat.AvgMemoryMB > highestMemory.AvgMemoryMB {
			highestMemory = stat
		}
		if lowestMemory == nil || stat.AvgMemoryMB < lowestMemory.AvgMemoryMB {
			lowestMemory = stat
		}
	}
	
	if highestMemory != nil && lowestMemory != nil {
		report.WriteString(fmt.Sprintf("- **Highest Memory Usage:** %s (%.2f MB)\n", 
			highestMemory.Type, highestMemory.AvgMemoryMB))
		report.WriteString(fmt.Sprintf("- **Lowest Memory Usage:** %s (%.2f MB)\n", 
			lowestMemory.Type, lowestMemory.AvgMemoryMB))
	}
	
	// File complexity insights
	var mostFiles, leastFiles *BlueprintStats
	for _, stat := range stats {
		if stat.Successes == 0 {
			continue
		}
		if mostFiles == nil || stat.AvgFiles > mostFiles.AvgFiles {
			mostFiles = stat
		}
		if leastFiles == nil || stat.AvgFiles < leastFiles.AvgFiles {
			leastFiles = stat
		}
	}
	
	if mostFiles != nil && leastFiles != nil {
		report.WriteString(fmt.Sprintf("- **Most Complex:** %s (%.1f files on average)\n", 
			mostFiles.Type, mostFiles.AvgFiles))
		report.WriteString(fmt.Sprintf("- **Simplest:** %s (%.1f files on average)\n", 
			leastFiles.Type, leastFiles.AvgFiles))
	}
	
	report.WriteString("\n")
}

// generateOptimizationRecommendations provides specific optimization recommendations
func generateOptimizationRecommendations(report *strings.Builder, summary BenchmarkSummary, stats map[string]*BlueprintStats) {
	recommendations := []string{}
	
	// Performance-based recommendations
	if summary.AverageDuration > 2*time.Second {
		recommendations = append(recommendations, 
			"**High Priority:** Average generation time exceeds 2s target. Consider:")
		recommendations = append(recommendations, 
			"  - Optimize template parsing and execution pipeline")
		recommendations = append(recommendations, 
			"  - Implement template caching for repeated generations")
		recommendations = append(recommendations, 
			"  - Parallelize file generation operations")
	}
	
	// Memory-based recommendations
	avgMemoryMB := float64(summary.TotalMemoryUsed) / float64(summary.SuccessfulRuns) / 1024 / 1024
	if avgMemoryMB > 50 {
		recommendations = append(recommendations, 
			"**Medium Priority:** High memory usage detected. Consider:")
		recommendations = append(recommendations, 
			"  - Implement streaming template processing")
		recommendations = append(recommendations, 
			"  - Optimize string concatenation and buffer usage")
		recommendations = append(recommendations, 
			"  - Add memory pooling for template contexts")
	}
	
	// Blueprint-specific recommendations
	for blueprintType, stat := range stats {
		if stat.AvgDuration > 3*time.Second {
			recommendations = append(recommendations, 
				fmt.Sprintf("**Blueprint-Specific:** %s blueprint is particularly slow (%.2fs avg):", 
					blueprintType, stat.AvgDuration.Seconds()))
			recommendations = append(recommendations, 
				"  - Profile this specific blueprint for bottlenecks")
			recommendations = append(recommendations, 
				"  - Consider simplifying the template structure")
		}
		
		if stat.SuccessRate < 100 {
			recommendations = append(recommendations, 
				fmt.Sprintf("**Reliability:** %s blueprint has %.1f%% success rate:", 
					blueprintType, stat.SuccessRate))
			recommendations = append(recommendations, 
				"  - Investigate and fix generation failures")
			recommendations = append(recommendations, 
				"  - Add better error handling and recovery")
		}
	}
	
	// General recommendations
	if len(recommendations) == 0 {
		recommendations = append(recommendations, 
			"**Current Performance:** All metrics are within acceptable ranges.")
		recommendations = append(recommendations, 
			"**Proactive Optimizations:**")
		recommendations = append(recommendations, 
			"  - Implement comprehensive caching strategy")
		recommendations = append(recommendations, 
			"  - Add performance monitoring in production")
		recommendations = append(recommendations, 
			"  - Consider adding progress indicators for long operations")
	}
	
	for _, rec := range recommendations {
		report.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	report.WriteString("\n")
}

// generateCrossPlatformAnalysis analyzes cross-platform performance differences
func generateCrossPlatformAnalysis(report *strings.Builder) {
	platformStats := make(map[string]*BlueprintStats)
	
	for _, result := range benchmarkSuite.Results {
		if _, exists := platformStats[result.Platform]; !exists {
			platformStats[result.Platform] = &BlueprintStats{
				Type: result.Platform,
				MinDuration: time.Hour,
			}
		}
		
		stat := platformStats[result.Platform]
		stat.Count++
		
		if result.Success {
			stat.Successes++
			stat.AvgDuration += result.Duration
			stat.TotalMemory += result.MemoryUsed
			stat.TotalFiles += result.FilesCreated
			
			if result.Duration < stat.MinDuration {
				stat.MinDuration = result.Duration
			}
			if result.Duration > stat.MaxDuration {
				stat.MaxDuration = result.Duration
			}
		}
	}
	
	// Calculate averages
	for _, stat := range platformStats {
		if stat.Successes > 0 {
			stat.AvgDuration = time.Duration(int64(stat.AvgDuration) / int64(stat.Successes))
			stat.AvgMemoryMB = float64(stat.TotalMemory) / float64(stat.Successes) / 1024 / 1024
		}
		stat.SuccessRate = float64(stat.Successes) / float64(stat.Count) * 100
	}
	
	if len(platformStats) > 1 {
		report.WriteString("| Platform | Avg Duration | Success Rate | Avg Memory (MB) |\n")
		report.WriteString("|----------|--------------|--------------|----------------|\n")
		
		for platform, stat := range platformStats {
			report.WriteString(fmt.Sprintf("| %s | %v | %.1f%% | %.2f |\n",
				platform, stat.AvgDuration.Round(time.Millisecond), stat.SuccessRate, stat.AvgMemoryMB))
		}
		report.WriteString("\n")
		
		// Find performance differences
		var fastest, slowest string
		var fastestTime, slowestTime time.Duration
		
		for platform, stat := range platformStats {
			if stat.Successes == 0 {
				continue
			}
			if fastest == "" || stat.AvgDuration < fastestTime {
				fastest = platform
				fastestTime = stat.AvgDuration
			}
			if slowest == "" || stat.AvgDuration > slowestTime {
				slowest = platform
				slowestTime = stat.AvgDuration
			}
		}
		
		if fastest != "" && slowest != "" && fastest != slowest {
			ratio := float64(slowestTime) / float64(fastestTime)
			report.WriteString(fmt.Sprintf("**Platform Performance:** %s is %.1fx faster than %s\n\n", 
				fastest, ratio, slowest))
		}
	} else {
		report.WriteString("*Single platform testing - no cross-platform comparison available*\n\n")
	}
}