package crossplatform

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

// PlatformTestResult represents test results for a single platform
type PlatformTestResult struct {
	Platform         string        `json:"platform"`
	Architecture     string        `json:"architecture"`
	Blueprint        string        `json:"blueprint"`
	GenerationTime   time.Duration `json:"generation_time"`
	CompileTime      time.Duration `json:"compile_time"`
	ExecuteTime      time.Duration `json:"execute_time"`
	FilesGenerated   int           `json:"files_generated"`
	BinarySize       int64         `json:"binary_size_bytes"`
	Success          bool          `json:"success"`
	GenerationError  string        `json:"generation_error,omitempty"`
	CompileError     string        `json:"compile_error,omitempty"`
	ExecuteError     string        `json:"execute_error,omitempty"`
	PathIssues       []string      `json:"path_issues,omitempty"`
	PermissionIssues []string      `json:"permission_issues,omitempty"`
	Timestamp        time.Time     `json:"timestamp"`
}

// PlatformCompatibilityReport holds comprehensive cross-platform test results
type PlatformCompatibilityReport struct {
	TestResults         []PlatformTestResult          `json:"test_results"`
	PlatformSummary     map[string]PlatformSummary    `json:"platform_summary"`
	CrossPlatformIssues []CrossPlatformIssue          `json:"cross_platform_issues"`
	CompatibilityScore  float64                       `json:"compatibility_score"`
	TestEnvironment     TestEnvironment               `json:"test_environment"`
	Timestamp           time.Time                     `json:"timestamp"`
}

// PlatformSummary provides aggregated results per platform
type PlatformSummary struct {
	Platform            string        `json:"platform"`
	TotalTests          int           `json:"total_tests"`
	SuccessfulTests     int           `json:"successful_tests"`
	FailedTests         int           `json:"failed_tests"`
	AvgGenerationTime   time.Duration `json:"avg_generation_time"`
	AvgCompileTime      time.Duration `json:"avg_compile_time"`
	AvgBinarySize       int64         `json:"avg_binary_size"`
	Issues              []string      `json:"issues"`
}

// CrossPlatformIssue represents platform-specific problems
type CrossPlatformIssue struct {
	Type        string   `json:"type"`
	Description string   `json:"description"`
	Platforms   []string `json:"affected_platforms"`
	Severity    string   `json:"severity"`
	Suggestions []string `json:"suggestions"`
}

// TestEnvironment captures the testing environment
type TestEnvironment struct {
	HostPlatform     string    `json:"host_platform"`
	HostArchitecture string    `json:"host_architecture"`
	GoVersion        string    `json:"go_version"`
	TestStartTime    time.Time `json:"test_start_time"`
	CrossCompilation bool      `json:"cross_compilation_enabled"`
}

// Global compatibility report instance
var compatibilityReport = &PlatformCompatibilityReport{
	TestResults:         make([]PlatformTestResult, 0),
	PlatformSummary:     make(map[string]PlatformSummary),
	CrossPlatformIssues: make([]CrossPlatformIssue, 0),
	TestEnvironment: TestEnvironment{
		HostPlatform:     runtime.GOOS,
		HostArchitecture: runtime.GOARCH,
		GoVersion:        runtime.Version(),
		TestStartTime:    time.Now(),
		CrossCompilation: true,
	},
}

// saveCrossPlatformJSONReport saves the compatibility report as JSON
func saveCrossPlatformJSONReport(filePath string) error {
	data, err := json.MarshalIndent(compatibilityReport, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal compatibility report: %w", err)
	}

	return os.WriteFile(filePath, data, 0644)
}

// saveCrossPlatformMarkdownReport saves the compatibility report as Markdown
func saveCrossPlatformMarkdownReport(filePath string) error {
	var report strings.Builder

	// Header
	report.WriteString("# Go-Starter Cross-Platform Compatibility Report\n\n")
	report.WriteString(fmt.Sprintf("**Generated:** %s\n", time.Now().Format("2006-01-02 15:04:05 MST")))
	report.WriteString(fmt.Sprintf("**Host Platform:** %s/%s\n", 
		compatibilityReport.TestEnvironment.HostPlatform, 
		compatibilityReport.TestEnvironment.HostArchitecture))
	report.WriteString(fmt.Sprintf("**Go Version:** %s\n\n", compatibilityReport.TestEnvironment.GoVersion))

	// Executive Summary
	report.WriteString("## Executive Summary\n\n")
	report.WriteString(fmt.Sprintf("- **Compatibility Score:** %.1f%%\n", compatibilityReport.CompatibilityScore))
	
	totalTests := len(compatibilityReport.TestResults)
	successfulTests := 0
	for _, result := range compatibilityReport.TestResults {
		if result.Success {
			successfulTests++
		}
	}
	
	report.WriteString(fmt.Sprintf("- **Test Results:** %d/%d tests passed\n", successfulTests, totalTests))
	report.WriteString(fmt.Sprintf("- **Platforms Tested:** %d\n", len(compatibilityReport.PlatformSummary)))
	report.WriteString(fmt.Sprintf("- **Critical Issues:** %d\n\n", countCriticalIssues()))

	// Compatibility Status
	report.WriteString("## Compatibility Status\n\n")
	if compatibilityReport.CompatibilityScore >= 95 {
		report.WriteString("✅ **EXCELLENT** - High cross-platform compatibility\n\n")
	} else if compatibilityReport.CompatibilityScore >= 85 {
		report.WriteString("⚠️ **GOOD** - Minor cross-platform issues detected\n\n")
	} else if compatibilityReport.CompatibilityScore >= 70 {
		report.WriteString("🔶 **FAIR** - Moderate cross-platform issues require attention\n\n")
	} else {
		report.WriteString("❌ **POOR** - Significant cross-platform issues detected\n\n")
	}

	// Platform Summary Table
	report.WriteString("## Platform Summary\n\n")
	report.WriteString("| Platform | Tests | Success Rate | Avg Generation | Avg Compile | Avg Binary Size | Issues |\n")
	report.WriteString("|----------|-------|--------------|----------------|-------------|-----------------|--------|\n")

	for platform, summary := range compatibilityReport.PlatformSummary {
		successRate := 0.0
		if summary.TotalTests > 0 {
			successRate = float64(summary.SuccessfulTests) / float64(summary.TotalTests) * 100
		}

		avgBinarySizeMB := float64(summary.AvgBinarySize) / 1024 / 1024

		issueCount := len(summary.Issues)
		issueStatus := "✅"
		if issueCount > 0 {
			issueStatus = fmt.Sprintf("⚠️ %d", issueCount)
		}

		report.WriteString(fmt.Sprintf("| %s | %d | %.1f%% | %v | %v | %.1f MB | %s |\n",
			platform,
			summary.TotalTests,
			successRate,
			summary.AvgGenerationTime.Round(time.Millisecond),
			summary.AvgCompileTime.Round(time.Millisecond),
			avgBinarySizeMB,
			issueStatus))
	}
	report.WriteString("\n")

	// Detailed Test Results
	report.WriteString("## Detailed Test Results\n\n")

	// Group results by blueprint
	blueprintResults := make(map[string][]PlatformTestResult)
	for _, result := range compatibilityReport.TestResults {
		blueprintResults[result.Blueprint] = append(blueprintResults[result.Blueprint], result)
	}

	for blueprint, results := range blueprintResults {
		report.WriteString(fmt.Sprintf("### %s Blueprint\n\n", strings.Title(blueprint)))
		report.WriteString("| Platform | Status | Generation Time | Compile Time | Binary Size | Issues |\n")
		report.WriteString("|----------|--------|-----------------|--------------|-------------|--------|\n")

		for _, result := range results {
			status := "✅ Success"
			if !result.Success {
				status = "❌ Failed"
			}

			binarySizeMB := float64(result.BinarySize) / 1024 / 1024
			totalIssues := len(result.PathIssues) + len(result.PermissionIssues)
			
			issueCount := ""
			if result.GenerationError != "" || result.CompileError != "" || result.ExecuteError != "" {
				issueCount = "❌ Error"
			} else if totalIssues > 0 {
				issueCount = fmt.Sprintf("⚠️ %d", totalIssues)
			} else {
				issueCount = "✅"
			}

			report.WriteString(fmt.Sprintf("| %s | %s | %v | %v | %.1f MB | %s |\n",
				result.Platform,
				status,
				result.GenerationTime.Round(time.Millisecond),
				result.CompileTime.Round(time.Millisecond),
				binarySizeMB,
				issueCount))
		}
		report.WriteString("\n")
	}

	// Cross-Platform Issues Analysis
	if len(compatibilityReport.CrossPlatformIssues) > 0 {
		report.WriteString("## Cross-Platform Issues\n\n")

		for i, issue := range compatibilityReport.CrossPlatformIssues {
			severityIcon := "🔶"
			switch issue.Severity {
			case "High":
				severityIcon = "🔴"
			case "Medium":
				severityIcon = "🟡"
			case "Low":
				severityIcon = "🟢"
			}

			report.WriteString(fmt.Sprintf("### %s Issue %d: %s\n\n", severityIcon, i+1, issue.Type))
			report.WriteString(fmt.Sprintf("**Severity:** %s\n\n", issue.Severity))
			report.WriteString(fmt.Sprintf("**Affected Platforms:** %s\n\n", strings.Join(issue.Platforms, ", ")))
			report.WriteString(fmt.Sprintf("**Description:** %s\n\n", issue.Description))

			if len(issue.Suggestions) > 0 {
				report.WriteString("**Suggestions:**\n")
				for _, suggestion := range issue.Suggestions {
					report.WriteString(fmt.Sprintf("- %s\n", suggestion))
				}
				report.WriteString("\n")
			}
		}
	} else {
		report.WriteString("## Cross-Platform Issues\n\n")
		report.WriteString("✅ No cross-platform issues detected!\n\n")
	}

	// Platform-Specific Details
	report.WriteString("## Platform-Specific Details\n\n")

	for platform, summary := range compatibilityReport.PlatformSummary {
		report.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(platform)))

		// Performance metrics
		report.WriteString("**Performance:**\n")
		report.WriteString(fmt.Sprintf("- Average Generation Time: %v\n", summary.AvgGenerationTime.Round(time.Millisecond)))
		report.WriteString(fmt.Sprintf("- Average Compile Time: %v\n", summary.AvgCompileTime.Round(time.Millisecond)))
		report.WriteString(fmt.Sprintf("- Average Binary Size: %.1f MB\n", float64(summary.AvgBinarySize)/1024/1024))

		// Success rate
		successRate := float64(summary.SuccessfulTests) / float64(summary.TotalTests) * 100
		report.WriteString(fmt.Sprintf("- Success Rate: %.1f%% (%d/%d)\n\n", 
			successRate, summary.SuccessfulTests, summary.TotalTests))

		// Issues
		if len(summary.Issues) > 0 {
			report.WriteString("**Issues:**\n")
			for _, issue := range summary.Issues {
				report.WriteString(fmt.Sprintf("- %s\n", issue))
			}
			report.WriteString("\n")
		} else {
			report.WriteString("**Issues:** None detected ✅\n\n")
		}
	}

	// Recommendations
	report.WriteString("## Recommendations\n\n")
	generateCompatibilityRecommendations(&report)

	// Technical Details
	report.WriteString("## Technical Details\n\n")
	report.WriteString("### Test Configuration\n\n")
	report.WriteString(fmt.Sprintf("- **Host Platform:** %s/%s\n", 
		compatibilityReport.TestEnvironment.HostPlatform, 
		compatibilityReport.TestEnvironment.HostArchitecture))
	report.WriteString(fmt.Sprintf("- **Go Version:** %s\n", compatibilityReport.TestEnvironment.GoVersion))
	report.WriteString(fmt.Sprintf("- **Cross-Compilation:** %t\n", compatibilityReport.TestEnvironment.CrossCompilation))
	report.WriteString(fmt.Sprintf("- **Test Duration:** %v\n\n", time.Since(compatibilityReport.TestEnvironment.TestStartTime)))

	// Test Matrix
	report.WriteString("### Test Matrix\n\n")
	report.WriteString("| Blueprint | Platforms Tested | Success Rate |\n")
	report.WriteString("|-----------|------------------|---------------|\n")

	blueprintStats := calculateBlueprintCompatibility()
	for blueprint, stats := range blueprintStats {
		report.WriteString(fmt.Sprintf("| %s | %d | %.1f%% |\n", 
			blueprint, stats.PlatformsTested, stats.SuccessRate))
	}
	report.WriteString("\n")

	// Footer
	report.WriteString("---\n\n")
	report.WriteString("*Generated by go-starter Cross-Platform Compatibility Testing Suite*\n")

	return os.WriteFile(filePath, []byte(report.String()), 0644)
}

// BlueprintCompatibilityStats holds compatibility statistics for a blueprint
type BlueprintCompatibilityStats struct {
	Blueprint        string
	PlatformsTested  int
	SuccessfulTests  int
	FailedTests      int
	SuccessRate      float64
}

// calculateBlueprintCompatibility calculates compatibility stats per blueprint
func calculateBlueprintCompatibility() map[string]BlueprintCompatibilityStats {
	stats := make(map[string]BlueprintCompatibilityStats)

	for _, result := range compatibilityReport.TestResults {
		blueprint := result.Blueprint
		
		if _, exists := stats[blueprint]; !exists {
			stats[blueprint] = BlueprintCompatibilityStats{
				Blueprint: blueprint,
			}
		}

		stat := stats[blueprint]
		stat.PlatformsTested++

		if result.Success {
			stat.SuccessfulTests++
		} else {
			stat.FailedTests++
		}

		stat.SuccessRate = float64(stat.SuccessfulTests) / float64(stat.PlatformsTested) * 100
		stats[blueprint] = stat
	}

	return stats
}

// countCriticalIssues counts high-severity cross-platform issues
func countCriticalIssues() int {
	count := 0
	for _, issue := range compatibilityReport.CrossPlatformIssues {
		if issue.Severity == "High" {
			count++
		}
	}
	return count
}

// generateCompatibilityRecommendations provides platform-specific recommendations
func generateCompatibilityRecommendations(report *strings.Builder) {
	recommendations := []string{}

	// Overall compatibility recommendations
	if compatibilityReport.CompatibilityScore < 95 {
		recommendations = append(recommendations, 
			"**High Priority:** Improve cross-platform compatibility to reach 95%+ target")
	}

	// Platform-specific recommendations
	for platform, summary := range compatibilityReport.PlatformSummary {
		successRate := float64(summary.SuccessfulTests) / float64(summary.TotalTests) * 100
		
		if successRate < 90 {
			recommendations = append(recommendations, 
				fmt.Sprintf("**%s:** Success rate %.1f%% below target - investigate platform-specific issues", 
					strings.Title(platform), successRate))
		}

		if len(summary.Issues) > 0 {
			recommendations = append(recommendations, 
				fmt.Sprintf("**%s:** Address %d identified issues", 
					strings.Title(platform), len(summary.Issues)))
		}
	}

	// Performance recommendations
	for platform, summary := range compatibilityReport.PlatformSummary {
		if summary.AvgGenerationTime > 5*time.Second {
			recommendations = append(recommendations, 
				fmt.Sprintf("**%s:** Generation time %.2fs is high - optimize for this platform", 
					strings.Title(platform), summary.AvgGenerationTime.Seconds()))
		}

		if summary.AvgCompileTime > 30*time.Second {
			recommendations = append(recommendations, 
				fmt.Sprintf("**%s:** Compile time %.2fs is high - consider build optimizations", 
					strings.Title(platform), summary.AvgCompileTime.Seconds()))
		}
	}

	// Issue-specific recommendations
	hasPathIssues := false
	hasPermissionIssues := false
	
	for _, result := range compatibilityReport.TestResults {
		if len(result.PathIssues) > 0 {
			hasPathIssues = true
		}
		if len(result.PermissionIssues) > 0 {
			hasPermissionIssues = true
		}
	}

	if hasPathIssues {
		recommendations = append(recommendations, 
			"**Path Handling:** Implement consistent cross-platform path handling using filepath.Join()")
	}

	if hasPermissionIssues {
		recommendations = append(recommendations, 
			"**Permissions:** Review and standardize file permission handling across platforms")
	}

	// Cross-platform issue recommendations
	for _, issue := range compatibilityReport.CrossPlatformIssues {
		if issue.Severity == "High" {
			recommendations = append(recommendations, 
				fmt.Sprintf("**Critical:** %s affects %s - immediate attention required", 
					issue.Description, strings.Join(issue.Platforms, ", ")))
		}
	}

	// General recommendations if no issues
	if len(recommendations) == 0 {
		recommendations = append(recommendations, 
			"**Excellent Compatibility:** All platforms are performing well")
		recommendations = append(recommendations, 
			"**Maintenance:** Continue monitoring cross-platform performance")
		recommendations = append(recommendations, 
			"**Future:** Consider adding more target platforms for testing")
	}

	// Output recommendations
	for _, rec := range recommendations {
		report.WriteString(fmt.Sprintf("- %s\n", rec))
	}
	report.WriteString("\n")
}