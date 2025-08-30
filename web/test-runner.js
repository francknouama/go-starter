#!/usr/bin/env node

/**
 * Comprehensive Test Runner for Go-Starter Web UI
 * 
 * This script orchestrates all testing phases:
 * - Unit tests (Jest + React Testing Library)
 * - Integration tests (API + WebSocket)
 * - End-to-end tests (Playwright)
 * - Coverage reporting and validation
 */

const { execSync } = require('child_process');
const fs = require('fs');
const path = require('path');

const colors = {
  reset: '\x1b[0m',
  bright: '\x1b[1m',
  red: '\x1b[31m',
  green: '\x1b[32m',
  yellow: '\x1b[33m',
  blue: '\x1b[34m',
  magenta: '\x1b[35m',
  cyan: '\x1b[36m'
};

function log(message, color = colors.reset) {
  console.log(`${color}${message}${colors.reset}`);
}

function logHeader(message) {
  log(`\n${'='.repeat(60)}`, colors.cyan);
  log(`${colors.bright}${message}`, colors.cyan);
  log(`${'='.repeat(60)}`, colors.cyan);
}

function logStep(message) {
  log(`\n${colors.bright}▶ ${message}${colors.reset}`, colors.blue);
}

function logSuccess(message) {
  log(`✅ ${message}`, colors.green);
}

function logError(message) {
  log(`❌ ${message}`, colors.red);
}

function logWarning(message) {
  log(`⚠️ ${message}`, colors.yellow);
}

class TestRunner {
  constructor() {
    this.results = {
      unit: { passed: false, coverage: 0, duration: 0 },
      integration: { passed: false, coverage: 0, duration: 0 },
      e2e: { passed: false, duration: 0 }
    };
    this.coverageThresholds = {
      statements: 80,
      branches: 80,
      functions: 80,
      lines: 80
    };
  }

  async run() {
    logHeader('Go-Starter Web UI Comprehensive Test Suite');
    
    try {
      await this.setupEnvironment();
      await this.runUnitTests();
      await this.runIntegrationTests();
      await this.runE2ETests();
      await this.generateCoverageReport();
      await this.validateCoverage();
      
      this.printSummary();
      
      if (this.allTestsPassed()) {
        logSuccess('All tests passed! ✨');
        process.exit(0);
      } else {
        logError('Some tests failed. Please review the results above.');
        process.exit(1);
      }
    } catch (error) {
      logError(`Test runner failed: ${error.message}`);
      console.error(error);
      process.exit(1);
    }
  }

  async setupEnvironment() {
    logStep('Setting up test environment...');
    
    // Ensure directories exist
    const dirs = ['coverage', 'test-results', 'playwright-report'];
    dirs.forEach(dir => {
      if (!fs.existsSync(dir)) {
        fs.mkdirSync(dir, { recursive: true });
      }
    });

    // Check if Go backend is running
    try {
      execSync('curl -f http://localhost:8080/api/v1/health', { 
        stdio: 'pipe',
        timeout: 5000 
      });
      logSuccess('Go backend is running');
    } catch (error) {
      logWarning('Go backend is not running. Some integration tests may fail.');
    }

    logSuccess('Test environment ready');
  }

  async runUnitTests() {
    logHeader('Running Unit Tests (Jest + React Testing Library)');
    
    const startTime = Date.now();
    
    try {
      const output = execSync('npm run test:ci', { 
        encoding: 'utf8',
        stdio: 'pipe'
      });
      
      this.results.unit.duration = Date.now() - startTime;
      this.results.unit.passed = true;
      
      // Extract coverage information
      const coverageMatch = output.match(/All files[^\n]*?(\d+\.?\d*)/);
      if (coverageMatch) {
        this.results.unit.coverage = parseFloat(coverageMatch[1]);
      }
      
      logSuccess(`Unit tests passed in ${this.results.unit.duration}ms`);
      
      // Show test summary
      const testSummary = this.extractTestSummary(output);
      if (testSummary.total > 0) {
        log(`Tests: ${testSummary.passed} passed, ${testSummary.total} total`);
        log(`Snapshots: ${testSummary.snapshots} total`);
      }
      
    } catch (error) {
      this.results.unit.duration = Date.now() - startTime;
      this.results.unit.passed = false;
      
      logError('Unit tests failed');
      console.error(error.stdout || error.message);
    }
  }

  async runIntegrationTests() {
    logHeader('Running Integration Tests');
    
    const startTime = Date.now();
    
    try {
      // Run WebSocket integration tests
      logStep('Testing WebSocket integration...');
      execSync('npm run test -- --testPathPattern=websocket --verbose', {
        stdio: 'inherit'
      });
      
      // Run API integration tests
      logStep('Testing API integration...');
      execSync('npm run test -- --testPathPattern=api --verbose', {
        stdio: 'inherit'
      });
      
      this.results.integration.duration = Date.now() - startTime;
      this.results.integration.passed = true;
      
      logSuccess(`Integration tests passed in ${this.results.integration.duration}ms`);
      
    } catch (error) {
      this.results.integration.duration = Date.now() - startTime;
      this.results.integration.passed = false;
      
      logError('Integration tests failed');
      console.error(error.message);
    }
  }

  async runE2ETests() {
    logHeader('Running End-to-End Tests (Playwright)');
    
    const startTime = Date.now();
    
    try {
      // Check if both frontend and backend are running
      const processes = this.checkRequiredProcesses();
      if (!processes.frontend || !processes.backend) {
        throw new Error('Frontend and backend servers must be running for E2E tests');
      }
      
      // Run Playwright tests
      execSync('npm run test:e2e', {
        stdio: 'inherit',
        env: { ...process.env, CI: 'true' }
      });
      
      this.results.e2e.duration = Date.now() - startTime;
      this.results.e2e.passed = true;
      
      logSuccess(`E2E tests passed in ${this.results.e2e.duration}ms`);
      
    } catch (error) {
      this.results.e2e.duration = Date.now() - startTime;
      this.results.e2e.passed = false;
      
      logError('E2E tests failed');
      console.error(error.message);
      
      // Generate E2E test report even on failure
      try {
        execSync('npm run test:e2e:report', { stdio: 'inherit' });
        log('E2E test report generated: playwright-report/index.html');
      } catch (reportError) {
        logWarning('Could not generate E2E test report');
      }
    }
  }

  async generateCoverageReport() {
    logHeader('Generating Coverage Reports');
    
    try {
      // Merge coverage from different test types
      logStep('Merging coverage reports...');
      
      const coverageFiles = [
        'coverage/coverage-final.json',
        'coverage/integration-coverage.json'
      ].filter(file => fs.existsSync(file));
      
      if (coverageFiles.length > 0) {
        // Use nyc to merge coverage reports
        try {
          execSync(`npx nyc merge coverage coverage/merged-coverage.json`, {
            stdio: 'pipe'
          });
          logSuccess('Coverage reports merged');
        } catch (error) {
          logWarning('Could not merge coverage reports');
        }
      }
      
      // Generate HTML coverage report
      logStep('Generating HTML coverage report...');
      execSync('npx nyc report --reporter=html --report-dir=coverage/html', {
        stdio: 'inherit'
      });
      
      // Generate lcov report for CI
      execSync('npx nyc report --reporter=lcov', {
        stdio: 'pipe'
      });
      
      logSuccess('Coverage reports generated: coverage/html/index.html');
      
    } catch (error) {
      logWarning('Could not generate complete coverage report');
      console.error(error.message);
    }
  }

  async validateCoverage() {
    logHeader('Validating Coverage Thresholds');
    
    try {
      const coverageSummaryPath = 'coverage/coverage-summary.json';
      
      if (!fs.existsSync(coverageSummaryPath)) {
        logWarning('Coverage summary not found, skipping validation');
        return;
      }
      
      const coverageSummary = JSON.parse(fs.readFileSync(coverageSummaryPath, 'utf8'));
      const totalCoverage = coverageSummary.total;
      
      logStep('Checking coverage thresholds...');
      
      const results = {};
      let allPassed = true;
      
      for (const [metric, threshold] of Object.entries(this.coverageThresholds)) {
        const actual = totalCoverage[metric].pct;
        const passed = actual >= threshold;
        
        results[metric] = { actual, threshold, passed };
        
        if (passed) {
          log(`  ✅ ${metric}: ${actual}% (threshold: ${threshold}%)`);
        } else {
          log(`  ❌ ${metric}: ${actual}% (threshold: ${threshold}%)`);
          allPassed = false;
        }
      }
      
      if (allPassed) {
        logSuccess('All coverage thresholds met');
      } else {
        logError('Some coverage thresholds not met');
        
        // Show top uncovered files
        this.showUncoveredFiles(coverageSummary);
      }
      
    } catch (error) {
      logWarning('Could not validate coverage thresholds');
      console.error(error.message);
    }
  }

  showUncoveredFiles(coverageSummary) {
    logStep('Files with low coverage:');
    
    const files = Object.entries(coverageSummary)
      .filter(([key]) => key !== 'total')
      .map(([path, coverage]) => ({
        path: path.replace(process.cwd(), '.'),
        coverage: coverage.statements.pct
      }))
      .filter(file => file.coverage < this.coverageThresholds.statements)
      .sort((a, b) => a.coverage - b.coverage)
      .slice(0, 10);
    
    files.forEach(file => {
      log(`  ${file.path}: ${file.coverage}%`, colors.yellow);
    });
    
    if (files.length === 0) {
      log('  No files with low coverage found');
    }
  }

  extractTestSummary(output) {
    const summary = {
      passed: 0,
      failed: 0,
      total: 0,
      snapshots: 0
    };
    
    // Extract test counts from Jest output
    const testMatch = output.match(/Tests:\s+(\d+)\s+passed,\s+(\d+)\s+total/);
    if (testMatch) {
      summary.passed = parseInt(testMatch[1]);
      summary.total = parseInt(testMatch[2]);
      summary.failed = summary.total - summary.passed;
    }
    
    const snapshotMatch = output.match(/Snapshots:\s+(\d+)\s+total/);
    if (snapshotMatch) {
      summary.snapshots = parseInt(snapshotMatch[1]);
    }
    
    return summary;
  }

  checkRequiredProcesses() {
    const processes = {
      frontend: false,
      backend: false
    };
    
    try {
      // Check frontend (port 3000)
      execSync('curl -f http://localhost:3000', { 
        stdio: 'pipe',
        timeout: 2000 
      });
      processes.frontend = true;
    } catch {}
    
    try {
      // Check backend (port 8080)
      execSync('curl -f http://localhost:8080/api/v1/health', { 
        stdio: 'pipe',
        timeout: 2000 
      });
      processes.backend = true;
    } catch {}
    
    return processes;
  }

  allTestsPassed() {
    return this.results.unit.passed && 
           this.results.integration.passed && 
           this.results.e2e.passed;
  }

  printSummary() {
    logHeader('Test Results Summary');
    
    const formatDuration = (ms) => {
      if (ms < 1000) return `${ms}ms`;
      return `${(ms / 1000).toFixed(2)}s`;
    };
    
    log('\n📊 Test Results:');
    log(`  Unit Tests:        ${this.results.unit.passed ? '✅ PASSED' : '❌ FAILED'} (${formatDuration(this.results.unit.duration)})`);
    log(`  Integration Tests: ${this.results.integration.passed ? '✅ PASSED' : '❌ FAILED'} (${formatDuration(this.results.integration.duration)})`);
    log(`  E2E Tests:         ${this.results.e2e.passed ? '✅ PASSED' : '❌ FAILED'} (${formatDuration(this.results.e2e.duration)})`);
    
    if (this.results.unit.coverage > 0) {
      log(`\n📈 Coverage: ${this.results.unit.coverage}%`);
    }
    
    const totalDuration = this.results.unit.duration + 
                         this.results.integration.duration + 
                         this.results.e2e.duration;
    
    log(`\n⏱️  Total Duration: ${formatDuration(totalDuration)}`);
    
    if (this.allTestsPassed()) {
      log('\n🎉 All tests passed successfully!', colors.green);
    } else {
      log('\n💥 Some tests failed. Check the logs above for details.', colors.red);
    }
  }
}

// CLI interface
const args = process.argv.slice(2);
const flags = {
  unit: args.includes('--unit'),
  integration: args.includes('--integration'),
  e2e: args.includes('--e2e'),
  coverage: args.includes('--coverage'),
  help: args.includes('--help') || args.includes('-h')
};

if (flags.help) {
  console.log(`
Go-Starter Web UI Test Runner

Usage:
  node test-runner.js [options]

Options:
  --unit          Run only unit tests
  --integration   Run only integration tests  
  --e2e           Run only E2E tests
  --coverage      Generate coverage report only
  --help, -h      Show this help message

Examples:
  node test-runner.js                    # Run all tests
  node test-runner.js --unit             # Run only unit tests
  node test-runner.js --e2e --coverage   # Run E2E tests and generate coverage
`);
  process.exit(0);
}

// Run the test runner
if (require.main === module) {
  const runner = new TestRunner();
  
  // If specific flags are provided, run only those test types
  if (flags.unit || flags.integration || flags.e2e || flags.coverage) {
    runner.run = async function() {
      logHeader('Go-Starter Web UI Selective Test Suite');
      
      try {
        await this.setupEnvironment();
        
        if (flags.unit) await this.runUnitTests();
        if (flags.integration) await this.runIntegrationTests();
        if (flags.e2e) await this.runE2ETests();
        if (flags.coverage) await this.generateCoverageReport();
        
        this.printSummary();
        
        if (this.allTestsPassed()) {
          logSuccess('Selected tests passed! ✨');
          process.exit(0);
        } else {
          logError('Some selected tests failed.');
          process.exit(1);
        }
      } catch (error) {
        logError(`Test runner failed: ${error.message}`);
        process.exit(1);
      }
    };
  }
  
  runner.run().catch(console.error);
}

module.exports = TestRunner;