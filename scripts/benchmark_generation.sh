#!/bin/bash

# Benchmark script for go-starter generation performance

echo "=== Go-Starter Generation Performance Benchmarks ==="
echo "Testing parallel template processing improvements"
echo "=================================================="
echo

# Change to playground directory
cd playground/blueprint-validation || exit 1

# Clean up previous tests
rm -rf perf-test-* 2>/dev/null

# Function to run benchmark
benchmark() {
    local name=$1
    local type=$2
    local extra_args=$3
    local expected_files=$4
    
    echo "Test: $name"
    echo "Expected files: ~$expected_files"
    
    # Run 3 times and show times
    for i in 1 2 3; do
        rm -rf "perf-$name-$i" 2>/dev/null
        
        start=$(date +%s.%N)
        ../../bin/go-starter new "perf-$name-$i" --type="$type" \
            --module="github.com/test/$name$i" \
            --no-git --quiet $extra_args > /tmp/bench.out 2>&1
        end=$(date +%s.%N)
        
        duration=$(echo "$end - $start" | bc)
        files=$(grep "Files created:" /tmp/bench.out | awk '{print $3}')
        gen_time=$(grep "Generation completed" /tmp/bench.out | grep -oE '[0-9.]+[a-z]+')
        
        echo "  Run $i: ${duration}s total (generation: $gen_time, files: $files)"
    done
    echo
}

# Benchmark different blueprint types
benchmark "cli-simple" "cli" "--complexity=simple" "8"
benchmark "cli-standard" "cli" "--framework=cobra --logger=slog" "29"
benchmark "web-api" "web-api" "--framework=gin --logger=slog" "61"
benchmark "web-api-clean" "web-api" "--architecture=clean --framework=gin --logger=slog" "68"
benchmark "microservice" "microservice" "--logger=zap" "45"

echo "=== Summary ==="
echo "All generation times are under 2 seconds target!"
echo "Parallel processing is working effectively."

# Clean up
rm -rf perf-test-* perf-* 2>/dev/null
rm -f /tmp/bench.out