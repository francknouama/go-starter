package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestMultiBinaryPerformanceATDD validates performance characteristics of the multi-binary structure
func TestMultiBinaryPerformanceATDD(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance tests in short mode")
	}

	tmpDir := t.TempDir()
	projectRoot := findProjectRoot(t)

	t.Run("build_performance_acceptable", func(t *testing.T) {
		// GIVEN: Multi-binary structure
		binaries := map[string]string{
			"go-starter":     "./cmd/go-starter",
			"go-starter-dev": "./cmd/go-starter-dev",
			"go-starter-web": "./web/cmd/web-server",
			"legacy":         ".",
		}

		buildTimes := make(map[string]time.Duration)
		binarySizes := make(map[string]int64)

		for name, buildPath := range binaries {
			t.Run("build_performance_"+name, func(t *testing.T) {
				// WHEN: Building binary and measuring time
				binaryName := name
				if runtime.GOOS == "windows" {
					binaryName += ".exe"
				}
				binaryPath := filepath.Join(tmpDir, binaryName)

				start := time.Now()
				cmd := exec.Command("go", "build", "-o", binaryPath, buildPath)
				cmd.Dir = projectRoot
				output, err := cmd.CombinedOutput()
				elapsed := time.Since(start)

				// THEN: Build should complete in reasonable time
				if err != nil {
					t.Logf("Build failed for %s: %s", name, string(output))
				}
				require.NoError(t, err, "Build should succeed")

				buildTimes[name] = elapsed
				t.Logf("Build time for %s: %v", name, elapsed)

				// Build time should be under 60 seconds (generous for CI environments)
				assert.Less(t, elapsed, 60*time.Second, 
					"Build time for %s should be under 60 seconds", name)

				// Check binary size
				info, err := os.Stat(binaryPath)
				require.NoError(t, err, "Should get binary info")
				
				size := info.Size()
				binarySizes[name] = size
				t.Logf("Binary size for %s: %.2f MB", name, float64(size)/(1024*1024))

				// Binary size should be reasonable
				minSize := int64(5 * 1024 * 1024)  // 5MB
				maxSize := int64(60 * 1024 * 1024) // 60MB (generous for embedded assets)
				
				assert.GreaterOrEqual(t, size, minSize, 
					"Binary %s should be at least 5MB", name)
				assert.LessOrEqual(t, size, maxSize, 
					"Binary %s should be under 60MB", name)
			})
		}

		// Compare build times - they should be similar
		t.Run("build_times_consistent", func(t *testing.T) {
			if len(buildTimes) < 2 {
				t.Skip("Need at least 2 binaries to compare build times")
			}

			var times []time.Duration
			for _, duration := range buildTimes {
				times = append(times, duration)
			}

			// Find min and max build times
			minTime := times[0]
			maxTime := times[0]
			for _, duration := range times {
				if duration < minTime {
					minTime = duration
				}
				if duration > maxTime {
					maxTime = duration
				}
			}

			// Max build time shouldn't be more than 3x the min time
			ratio := float64(maxTime) / float64(minTime)
			assert.Less(t, ratio, 3.0, 
				"Build time variation should be reasonable (max: %v, min: %v, ratio: %.2f)", 
				maxTime, minTime, ratio)
		})
	})

	t.Run("cli_startup_performance", func(t *testing.T) {
		// GIVEN: Built CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)
		
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// Test different command startup times
		commands := []struct {
			name string
			args []string
		}{
			{"version", []string{"version"}},
			{"help", []string{"--help"}},
			{"list", []string{"list"}},
			{"new_help", []string{"new", "--help"}},
		}

		for _, cmd := range commands {
			t.Run("startup_"+cmd.name, func(t *testing.T) {
				// WHEN: Running command and measuring startup time
				start := time.Now()
				execCmd := exec.Command(binaryPath, cmd.args...)
				output, err := execCmd.CombinedOutput()
				elapsed := time.Since(start)

				// THEN: Command should complete quickly
				if err != nil {
					t.Logf("Command %v failed: %s", cmd.args, string(output))
				}
				require.NoError(t, err, "Command should succeed")

				t.Logf("Startup time for %s: %v", cmd.name, elapsed)
				
				// Startup should be under 5 seconds (generous for CI)
				assert.Less(t, elapsed, 5*time.Second, 
					"Startup time for %s should be under 5 seconds", cmd.name)
				
				// For quick commands, should be under 2 seconds ideally
				if cmd.name == "version" || cmd.name == "help" {
					if elapsed > 2*time.Second {
						t.Logf("Warning: %s took longer than ideal (2s): %v", cmd.name, elapsed)
					}
				}
			})
		}

		// Test multiple rapid invocations (no memory leaks)
		t.Run("rapid_invocations", func(t *testing.T) {
			iterations := 5
			times := make([]time.Duration, iterations)
			
			for i := 0; i < iterations; i++ {
				start := time.Now()
				cmd := exec.Command(binaryPath, "version")
				err := cmd.Run()
				times[i] = time.Since(start)
				
				require.NoError(t, err, "Iteration %d should succeed", i)
			}

			// Calculate average time
			var total time.Duration
			for _, t := range times {
				total += t
			}
			avg := total / time.Duration(iterations)
			
			t.Logf("Average startup time over %d iterations: %v", iterations, avg)
			
			// Performance shouldn't degrade significantly
			for i, elapsed := range times {
				if elapsed > avg*2 {
					t.Logf("Warning: Iteration %d was significantly slower: %v vs avg %v", 
						i, elapsed, avg)
				}
			}
		})
	})

	t.Run("memory_usage_reasonable", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)
		
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Running memory-intensive operations
		memoryTestCases := []struct {
			name string
			args []string
		}{
			{"list_all_blueprints", []string{"list"}},
			{"dry_run_generation", []string{"new", "test", "--type=web-api", "--dry-run"}},
			{"help_with_all_flags", []string{"new", "--advanced", "--help"}},
		}

		for _, tc := range memoryTestCases {
			t.Run("memory_"+tc.name, func(t *testing.T) {
				// Note: Direct memory measurement requires platform-specific tools
				// Here we test indirectly by ensuring operations complete without hanging
				
				start := time.Now()
				cmd := exec.Command(binaryPath, tc.args...)
				output, err := cmd.CombinedOutput()
				elapsed := time.Since(start)

				// THEN: Should complete without excessive memory usage (indicated by hanging)
				require.NoError(t, err, "Memory test %s should succeed", tc.name)
				
				// If it takes too long, might indicate memory issues
				assert.Less(t, elapsed, 10*time.Second, 
					"Memory test %s should complete promptly", tc.name)
				
				// Output should be reasonable size (not excessive)
				outputSize := len(output)
				assert.Less(t, outputSize, 1024*1024, // 1MB
					"Output size should be reasonable for %s", tc.name)
				
				t.Logf("Memory test %s: %v, output size: %d bytes", 
					tc.name, elapsed, outputSize)
			})
		}
	})

	t.Run("project_generation_performance", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)
		
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// Test generation performance for different project types
		generationTests := []struct {
			name       string
			projectType string
			complexity string
			maxTime    time.Duration
		}{
			{"simple_cli", "cli", "simple", 10 * time.Second},
			{"standard_cli", "cli", "standard", 15 * time.Second},
			{"web_api", "web-api", "standard", 20 * time.Second},
			{"library", "library", "standard", 10 * time.Second},
		}

		for _, gt := range generationTests {
			t.Run("generation_"+gt.name, func(t *testing.T) {
				// WHEN: Generating project
				projectName := "perf-" + gt.name
				start := time.Now()
				
				cmd := exec.Command(binaryPath, "new", projectName,
					"--type="+gt.projectType,
					"--complexity="+gt.complexity,
					"--module=github.com/test/"+projectName)
				cmd.Dir = tmpDir
				output, err := cmd.CombinedOutput()
				elapsed := time.Since(start)

				// THEN: Generation should complete in reasonable time
				if err != nil {
					t.Logf("Generation failed for %s: %s", gt.name, string(output))
				}
				require.NoError(t, err, "Generation should succeed")

				t.Logf("Generation time for %s: %v", gt.name, elapsed)
				assert.Less(t, elapsed, gt.maxTime, 
					"Generation time for %s should be under %v", gt.name, gt.maxTime)

				// Verify project was actually created
				projectDir := filepath.Join(tmpDir, projectName)
				assert.DirExists(t, projectDir, "Project directory should exist")
				
				// Count generated files
				files, err := filepath.Glob(filepath.Join(projectDir, "**"))
				require.NoError(t, err, "Should list generated files")
				t.Logf("Generated %d files for %s", len(files), gt.name)
			})
		}
	})

	t.Run("concurrent_operations", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)
		
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Running multiple operations concurrently
		concurrency := 3
		results := make(chan error, concurrency)
		
		start := time.Now()
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				projectName := "concurrent-test-" + string(rune('a'+id))
				cmd := exec.Command(binaryPath, "new", projectName,
					"--type=cli", "--complexity=simple",
					"--module=github.com/test/"+projectName)
				cmd.Dir = tmpDir
				_, err := cmd.CombinedOutput()
				results <- err
			}(i)
		}

		// Wait for all operations to complete
		for i := 0; i < concurrency; i++ {
			err := <-results
			assert.NoError(t, err, "Concurrent operation %d should succeed", i)
		}
		elapsed := time.Since(start)

		// THEN: All operations should complete reasonably quickly
		t.Logf("Concurrent operations completed in: %v", elapsed)
		assert.Less(t, elapsed, 30*time.Second, 
			"Concurrent operations should complete in under 30 seconds")

		// Verify all projects were created
		for i := 0; i < concurrency; i++ {
			projectName := "concurrent-test-" + string(rune('a'+i))
			projectDir := filepath.Join(tmpDir, projectName)
			assert.DirExists(t, projectDir, "Concurrent project %d should exist", i)
		}
	})

	t.Run("resource_cleanup", func(t *testing.T) {
		// GIVEN: CLI binary
		binaryName := "go-starter"
		if runtime.GOOS == "windows" {
			binaryName += ".exe"
		}
		binaryPath := filepath.Join(tmpDir, binaryName)
		
		buildCmd := exec.Command("go", "build", "-o", binaryPath, "./cmd/go-starter")
		buildCmd.Dir = projectRoot
		err := buildCmd.Run()
		require.NoError(t, err, "CLI should build")

		// WHEN: Running multiple operations that could leave resources
		operations := []struct {
			name string
			args []string
		}{
			{"help", []string{"--help"}},
			{"list", []string{"list"}},
			{"dry_run", []string{"new", "test", "--type=cli", "--dry-run"}},
			{"version", []string{"version"}},
		}

		initialFiles, err := filepath.Glob(filepath.Join(tmpDir, "*"))
		require.NoError(t, err, "Should list initial files")
		initialCount := len(initialFiles)

		for _, op := range operations {
			cmd := exec.Command(binaryPath, op.args...)
			cmd.Dir = tmpDir
			_, err := cmd.CombinedOutput()
			require.NoError(t, err, "Operation %s should succeed", op.name)
		}

		// THEN: No temporary files should be left behind
		finalFiles, err := filepath.Glob(filepath.Join(tmpDir, "*"))
		require.NoError(t, err, "Should list final files")
		finalCount := len(finalFiles)

		// File count should not have increased (no temp files left)
		assert.Equal(t, initialCount, finalCount, 
			"No temporary files should be left behind after operations")
	})
}