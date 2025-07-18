#!/bin/bash

set -e

# Define packages to include in coverage report
PACKAGES=(
    "./cmd/..."
    "./internal/config/..."
    "./internal/generator/..."
    "./internal/prompts/..."
    "./internal/templates/..."
    "./internal/utils/..."
    "./pkg/..."
    "./tests/acceptance/..."
)

# Minimum acceptable coverage percentage
MIN_COVERAGE=50.0

# Output file for coverage profile
COVER_PROFILE="coverage.out"

# Remove existing coverage profile
rm -f "${COVER_PROFILE}"

echo "Running tests and collecting coverage..."

# Run tests and generate coverage profile
go test -v -covermode=count -coverprofile="${COVER_PROFILE}" "${PACKAGES[@]}"

# Check if coverage profile was generated
if [ ! -f "${COVER_PROFILE}" ]; then
    echo "Error: Coverage profile '${COVER_PROFILE}' not generated."
    exit 1
fi

# Get overall coverage percentage
OVERALL_COVERAGE=$(go tool cover -func="${COVER_PROFILE}" | grep "total:" | awk '{print $3}' | sed 's/%//')

echo "\nOverall Test Coverage: ${OVERALL_COVERAGE}%"

# Check if overall coverage meets the minimum requirement
if (( $(echo "${OVERALL_COVERAGE} < ${MIN_COVERAGE}" | bc -l) )); then
    echo "❌ Coverage is below the minimum required ${MIN_COVERAGE}%"
    exit 1
else
    echo "✅ Coverage meets or exceeds the minimum required ${MIN_COVERAGE}%"
fi

echo "\nTo view detailed coverage report, run: go tool cover -html=${COVER_PROFILE}"
