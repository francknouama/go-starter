package performance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/pkg/types"
)

// ProfileConfig configures profiling options
type ProfileConfig struct {
	EnableCPU    bool   `json:"enable_cpu"`
	EnableMemory bool   `json:"enable_memory"`
	EnableBlock  bool   `json:"enable_block"`
	EnableMutex  bool   `json:"enable_mutex"`
	OutputDir    string `json:"output_dir"`
	Duration     time.Duration `json:"duration"`
}

// ProfileResult contains profiling results and analysis
type ProfileResult struct {
	Blueprint       string          `json:"blueprint"`
	Duration        time.Duration   `json:"duration"`
	CPUProfilePath  string          `json:"cpu_profile_path,omitempty"`
	MemProfilePath  string          `json:"mem_profile_path,omitempty"`
	BlockProfilePath string         `json:"block_profile_path,omitempty"`
	MutexProfilePath string         `json:"mutex_profile_path,omitempty"`
	MemoryStats     MemoryAnalysis  `json:"memory_stats"`
	CPUStats        CPUAnalysis     `json:"cpu_stats"`
	FileIOStats     FileIOAnalysis  `json:"file_io_stats"`
	Hotspots        []Hotspot       `json:"hotspots"`
	Recommendations []string        `json:"recommendations"`
	Timestamp       time.Time       `json:"timestamp"`
}

// MemoryAnalysis contains memory usage analysis
type MemoryAnalysis struct {
	AllocBytes      uint64  `json:"alloc_bytes"`
	TotalAllocBytes uint64  `json:"total_alloc_bytes"`
	SysBytes        uint64  `json:"sys_bytes"`
	NumGC           uint32  `json:"num_gc"`
	GCCPUFraction   float64 `json:"gc_cpu_fraction"`
	HeapInUse       uint64  `json:"heap_in_use"`
	HeapObjects     uint64  `json:"heap_objects"`
	StackInUse      uint64  `json:"stack_in_use"`
}

// CPUAnalysis contains CPU usage analysis
type CPUAnalysis struct {
	UserTime      time.Duration `json:"user_time"`
	SystemTime    time.Duration `json:"system_time"`
	TotalTime     time.Duration `json:"total_time"`
	CPUPercent    float64       `json:"cpu_percent"`
	Goroutines    int           `json:"goroutines"`
}

// FileIOAnalysis contains file I/O performance analysis
type FileIOAnalysis struct {
	FilesProcessed    int           `json:"files_processed"`
	BytesWritten      int64         `json:"bytes_written"`
	BytesRead         int64         `json:"bytes_read"`
	AvgFileSize       int64         `json:"avg_file_size"`
	IODuration        time.Duration `json:"io_duration"`
	IOThroughputMBps  float64       `json:"io_throughput_mbps"`
}

// Hotspot represents a performance hotspot
type Hotspot struct {
	Function    string        `json:"function"`
	Package     string        `json:"package"`
	Duration    time.Duration `json:"duration"`
	Percentage  float64       `json:"percentage"`
	CallCount   int64         `json:"call_count"`
	Description string        `json:"description"`
}

// DefaultProfileConfig returns a default profiling configuration
func DefaultProfileConfig() ProfileConfig {
	return ProfileConfig{
		EnableCPU:    true,
		EnableMemory: true,
		EnableBlock:  true,
		EnableMutex:  true,
		OutputDir:    "performance_profiles",
		Duration:     30 * time.Second,
	}
}

// ProfileGeneration profiles template generation performance
func ProfileGeneration(config types.ProjectConfig, profileConfig ProfileConfig) (*ProfileResult, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(profileConfig.OutputDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create profile output directory: %w", err)
	}

	result := &ProfileResult{
		Blueprint: getBlueprintID(config),
		Timestamp: time.Now(),
	}

	// Setup profiling
	if err := startProfiling(result, profileConfig); err != nil {
		return nil, fmt.Errorf("failed to start profiling: %w", err)
	}

	// Capture initial memory stats
	var memStart runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memStart)

	// Create temporary directory for generation
	tempDir, err := os.MkdirTemp("", "profile-generation-*")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	projectPath := filepath.Join(tempDir, config.Name)
	
	// Perform generation
	gen := generator.New()
	startTime := time.Now()
	
	genResult, err := gen.Generate(config, types.GenerationOptions{
		OutputPath: projectPath,
		DryRun:     false,
		NoGit:      true,
	})
	
	result.Duration = time.Since(startTime)

	// Capture final memory stats
	var memEnd runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&memEnd)

	// Stop profiling
	if err := stopProfiling(result, profileConfig); err != nil {
		return nil, fmt.Errorf("failed to stop profiling: %w", err)
	}

	// Analyze results
	result.MemoryStats = analyzeMemory(memStart, memEnd)
	result.CPUStats = analyzeCPU(result.Duration)
	result.FileIOStats = analyzeFileIO(genResult, result.Duration)
	result.Hotspots = identifyHotspots(result)
	result.Recommendations = generateRecommendations(result)

	if err != nil {
		return result, fmt.Errorf("generation failed: %w", err)
	}

	return result, nil
}

// startProfiling initializes profiling based on configuration
func startProfiling(result *ProfileResult, config ProfileConfig) error {
	timestamp := time.Now().Format("20060102_150405")
	
	if config.EnableCPU {
		cpuFile := filepath.Join(config.OutputDir, fmt.Sprintf("cpu_%s_%s.prof", result.Blueprint, timestamp))
		f, err := os.Create(cpuFile)
		if err != nil {
			return fmt.Errorf("failed to create CPU profile file: %w", err)
		}
		
		if err := pprof.StartCPUProfile(f); err != nil {
			f.Close()
			return fmt.Errorf("failed to start CPU profiling: %w", err)
		}
		
		result.CPUProfilePath = cpuFile
	}

	if config.EnableBlock {
		runtime.SetBlockProfileRate(1)
	}

	if config.EnableMutex {
		runtime.SetMutexProfileFraction(1)
	}

	return nil
}

// stopProfiling finalizes profiling and writes profile files
func stopProfiling(result *ProfileResult, config ProfileConfig) error {
	timestamp := time.Now().Format("20060102_150405")

	if config.EnableCPU {
		pprof.StopCPUProfile()
	}

	if config.EnableMemory {
		memFile := filepath.Join(config.OutputDir, fmt.Sprintf("mem_%s_%s.prof", result.Blueprint, timestamp))
		f, err := os.Create(memFile)
		if err != nil {
			return fmt.Errorf("failed to create memory profile file: %w", err)
		}
		defer f.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(f); err != nil {
			return fmt.Errorf("failed to write memory profile: %w", err)
		}
		
		result.MemProfilePath = memFile
	}

	if config.EnableBlock {
		blockFile := filepath.Join(config.OutputDir, fmt.Sprintf("block_%s_%s.prof", result.Blueprint, timestamp))
		f, err := os.Create(blockFile)
		if err != nil {
			return fmt.Errorf("failed to create block profile file: %w", err)
		}
		defer f.Close()

		if err := pprof.Lookup("block").WriteTo(f, 0); err != nil {
			return fmt.Errorf("failed to write block profile: %w", err)
		}
		
		result.BlockProfilePath = blockFile
	}

	if config.EnableMutex {
		mutexFile := filepath.Join(config.OutputDir, fmt.Sprintf("mutex_%s_%s.prof", result.Blueprint, timestamp))
		f, err := os.Create(mutexFile)
		if err != nil {
			return fmt.Errorf("failed to create mutex profile file: %w", err)
		}
		defer f.Close()

		if err := pprof.Lookup("mutex").WriteTo(f, 0); err != nil {
			return fmt.Errorf("failed to write mutex profile: %w", err)
		}
		
		result.MutexProfilePath = mutexFile
	}

	return nil
}

// analyzeMemory analyzes memory usage patterns
func analyzeMemory(start, end runtime.MemStats) MemoryAnalysis {
	return MemoryAnalysis{
		AllocBytes:      end.Alloc,
		TotalAllocBytes: end.TotalAlloc - start.TotalAlloc,
		SysBytes:        end.Sys,
		NumGC:           end.NumGC - start.NumGC,
		GCCPUFraction:   end.GCCPUFraction,
		HeapInUse:       end.HeapInuse,
		HeapObjects:     end.HeapObjects,
		StackInUse:      end.StackInuse,
	}
}

// analyzeCPU analyzes CPU usage patterns
func analyzeCPU(duration time.Duration) CPUAnalysis {
	return CPUAnalysis{
		TotalTime:  duration,
		Goroutines: runtime.NumGoroutine(),
		// Note: User/System time would need platform-specific implementation
		CPUPercent: calculateCPUPercent(duration),
	}
}

// analyzeFileIO analyzes file I/O performance
func analyzeFileIO(genResult *types.GenerationResult, duration time.Duration) FileIOAnalysis {
	if genResult == nil {
		return FileIOAnalysis{}
	}

	filesProcessed := len(genResult.FilesCreated)
	totalSize := calculateTotalFileSize(genResult.FilesCreated)
	
	throughput := float64(totalSize) / duration.Seconds() / 1024 / 1024 // MB/s

	analysis := FileIOAnalysis{
		FilesProcessed:   filesProcessed,
		BytesWritten:     totalSize,
		IODuration:       duration,
		IOThroughputMBps: throughput,
	}

	if filesProcessed > 0 {
		analysis.AvgFileSize = totalSize / int64(filesProcessed)
	}

	return analysis
}

// calculateTotalFileSize calculates the total size of all generated files
func calculateTotalFileSize(filePaths []string) int64 {
	var totalSize int64
	for _, filePath := range filePaths {
		if info, err := os.Stat(filePath); err == nil {
			totalSize += info.Size()
		}
	}
	return totalSize
}

// calculateCPUPercent estimates CPU usage percentage
func calculateCPUPercent(duration time.Duration) float64 {
	// This is a simplified estimation
	// In a real implementation, you would measure actual CPU time
	numCPU := runtime.NumCPU()
	return float64(duration.Nanoseconds()) / float64(time.Second.Nanoseconds()) / float64(numCPU) * 100
}

// identifyHotspots identifies performance hotspots from profiling data
func identifyHotspots(result *ProfileResult) []Hotspot {
	// This is a simplified implementation
	// In a real implementation, you would parse the profile files to identify hotspots
	hotspots := []Hotspot{}

	// Template processing hotspot
	if result.Duration > time.Second {
		hotspots = append(hotspots, Hotspot{
			Function:    "Generator.processTemplateFile",
			Package:     "internal/generator",
			Duration:    result.Duration / 3, // Estimated
			Percentage:  33.3,
			Description: "Template file processing and execution",
		})
	}

	// File I/O hotspot
	if result.FileIOStats.IOThroughputMBps < 10.0 {
		hotspots = append(hotspots, Hotspot{
			Function:    "os.WriteFile",
			Package:     "os",
			Duration:    result.Duration / 4, // Estimated
			Percentage:  25.0,
			Description: "File system write operations",
		})
	}

	// Memory allocation hotspot
	if result.MemoryStats.TotalAllocBytes > 50*1024*1024 { // 50MB
		hotspots = append(hotspots, Hotspot{
			Function:    "template.Execute",
			Package:     "text/template",
			Duration:    result.Duration / 5, // Estimated
			Percentage:  20.0,
			Description: "Template execution and memory allocation",
		})
	}

	return hotspots
}

// generateRecommendations provides optimization recommendations
func generateRecommendations(result *ProfileResult) []string {
	var recommendations []string

	// Performance recommendations
	if result.Duration > 2*time.Second {
		recommendations = append(recommendations, "Consider optimizing template processing pipeline for faster generation")
	}

	if result.MemoryStats.TotalAllocBytes > 100*1024*1024 { // 100MB
		recommendations = append(recommendations, "High memory usage detected - consider streaming template processing")
	}

	if result.FileIOStats.IOThroughputMBps < 5.0 {
		recommendations = append(recommendations, "Low I/O throughput - consider batch file operations or use buffered writes")
	}

	if result.MemoryStats.NumGC > 10 {
		recommendations = append(recommendations, "Frequent garbage collection detected - optimize memory allocations")
	}

	if len(result.Hotspots) > 3 {
		recommendations = append(recommendations, "Multiple performance hotspots identified - profile individual functions for optimization")
	}

	// Add default recommendation if no issues found
	if len(recommendations) == 0 {
		recommendations = append(recommendations, "Performance is within acceptable ranges - no immediate optimizations needed")
	}

	return recommendations
}

// SaveProfileResult saves profiling results to a JSON file
func SaveProfileResult(result *ProfileResult, filename string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal profile result: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// LoadProfileResult loads profiling results from a JSON file
func LoadProfileResult(filename string) (*ProfileResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read profile result file: %w", err)
	}

	var result ProfileResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal profile result: %w", err)
	}

	return &result, nil
}

// CompareProfiles compares two profile results and identifies changes
func CompareProfiles(baseline, current *ProfileResult) *ProfileComparison {
	return &ProfileComparison{
		BaselineBlueprint: baseline.Blueprint,
		CurrentBlueprint:  current.Blueprint,
		DurationChange:    current.Duration - baseline.Duration,
		MemoryChange:      int64(current.MemoryStats.TotalAllocBytes) - int64(baseline.MemoryStats.TotalAllocBytes),
		IOChange:          current.FileIOStats.IOThroughputMBps - baseline.FileIOStats.IOThroughputMBps,
		Timestamp:         time.Now(),
	}
}

// ProfileComparison represents a comparison between two profile results
type ProfileComparison struct {
	BaselineBlueprint string        `json:"baseline_blueprint"`
	CurrentBlueprint  string        `json:"current_blueprint"`
	DurationChange    time.Duration `json:"duration_change"`
	MemoryChange      int64         `json:"memory_change_bytes"`
	IOChange          float64       `json:"io_change_mbps"`
	Timestamp         time.Time     `json:"timestamp"`
}