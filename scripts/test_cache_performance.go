package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type benchmarkTest struct {
	name      string
	blueprint string
	args      []string
}

func main() {
	// Change to playground directory
	if err := os.Chdir("playground/blueprint-validation"); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Clean up any existing test projects
	exec.Command("rm", "-rf", "perf-*").Run()

	tests := []benchmarkTest{
		{
			name:      "CLI Simple",
			blueprint: "cli",
			args:      []string{"--complexity=simple"},
		},
		{
			name:      "CLI Standard", 
			blueprint: "cli",
			args:      []string{"--framework=cobra", "--logger=slog"},
		},
		{
			name:      "Web API Standard",
			blueprint: "web-api",
			args:      []string{"--framework=gin", "--logger=slog"},
		},
		{
			name:      "Web API Clean Architecture",
			blueprint: "web-api-clean",
			args:      []string{"--framework=gin", "--logger=slog"},
		},
		{
			name:      "Microservice",
			blueprint: "microservice",
			args:      []string{"--logger=zap"},
		},
	}

	fmt.Println("=== Template Caching Performance Analysis ===")
	fmt.Println()

	for i, test := range tests {
		fmt.Printf("%d. %s\n", i+1, test.name)
		
		// Cold cache run
		coldTime := runBenchmark(fmt.Sprintf("perf-cold-%d", i), test.blueprint, test.args)
		fmt.Printf("   Cold cache: %v\n", coldTime)
		
		// Warm cache run
		warmTime := runBenchmark(fmt.Sprintf("perf-warm-%d", i), test.blueprint, test.args)
		fmt.Printf("   Warm cache: %v\n", warmTime)
		
		// Calculate improvement
		if coldTime > 0 && warmTime > 0 {
			improvement := ((coldTime.Seconds() - warmTime.Seconds()) / coldTime.Seconds()) * 100
			fmt.Printf("   Improvement: %.1f%% (%.3fs faster)\n", improvement, coldTime.Seconds()-warmTime.Seconds())
		}
		fmt.Println()
	}

	// Clean up
	exec.Command("rm", "-rf", "perf-*").Run()
}

func runBenchmark(projectName, blueprint string, extraArgs []string) time.Duration {
	// Build command
	args := []string{
		"new", projectName,
		"--type=" + blueprint,
		"--module=github.com/test/" + projectName,
		"--no-git",
		"--quiet",
	}
	args = append(args, extraArgs...)

	// Ensure clean state
	os.RemoveAll(projectName)
	
	// Run command and measure time
	start := time.Now()
	cmd := exec.Command("../../bin/go-starter-cached", args...)
	cmd.Run()
	duration := time.Since(start)
	
	return duration
}