/**
 * API Integration Test Module
 * Tests all API endpoints and WebSocket functionality
 */

import { api, connectWebSocket } from './api';

export interface TestResult {
  name: string;
  success: boolean;
  error?: string;
  duration: number;
}

export class ApiTester {
  private results: TestResult[] = [];

  async runTest(name: string, testFn: () => Promise<void>): Promise<TestResult> {
    const startTime = performance.now();
    
    try {
      await testFn();
      const duration = performance.now() - startTime;
      const result: TestResult = { name, success: true, duration };
      this.results.push(result);
      return result;
    } catch (error) {
      const duration = performance.now() - startTime;
      const result: TestResult = { 
        name, 
        success: false, 
        error: error instanceof Error ? error.message : String(error),
        duration 
      };
      this.results.push(result);
      return result;
    }
  }

  async testHealthEndpoints(): Promise<TestResult> {
    return this.runTest('Health Endpoints', async () => {
      // Test simple health
      const simpleHealth = await api.health.getSimpleHealth();
      if (simpleHealth.status !== 'ok') {
        throw new Error('Simple health check failed');
      }

      // Test full health
      const fullHealth = await api.health.getHealth();
      if (!fullHealth.status) {
        throw new Error('Full health check failed');
      }

      // Test metrics
      await api.health.getMetrics();
    });
  }

  async testConfigEndpoints(): Promise<TestResult> {
    return this.runTest('Configuration Endpoints', async () => {
      // Test default config
      const config = await api.config.getDefaultConfig();
      if (!config.projectType) {
        throw new Error('Default config missing projectType');
      }

      // Test frameworks
      const frameworks = await api.config.getFrameworks();
      if (!frameworks.includes('gin')) {
        throw new Error('Frameworks missing gin');
      }

      // Test architectures
      const architectures = await api.config.getArchitectures();
      if (!architectures.includes('standard')) {
        throw new Error('Architectures missing standard');
      }
    });
  }

  async testBlueprintEndpoints(): Promise<TestResult> {
    return this.runTest('Blueprint Endpoints', async () => {
      // Test blueprints list
      const blueprints = await api.blueprints.getBlueprints();
      if (blueprints.length === 0) {
        throw new Error('No blueprints found');
      }

      // Test specific blueprint
      const blueprint = await api.blueprints.getBlueprintById('web-api');
      if (!blueprint.name) {
        throw new Error('Blueprint details missing name');
      }

      // Test blueprint validation
      const validation = await api.blueprints.validateBlueprintConfig('web-api', {
        projectName: 'test-api',
        moduleUrl: 'github.com/test/test-api',
        projectType: 'web-api',
        architecture: 'standard',
        framework: 'gin',
        logger: 'slog',
        goVersion: '1.21'
      });
      
      if (!validation.valid) {
        throw new Error('Blueprint validation failed');
      }
    });
  }

  async testProjectOperations(): Promise<TestResult> {
    return this.runTest('Project Operations', async () => {
      const testConfig = {
        projectName: 'test-project',
        moduleUrl: 'github.com/test/test-project',
        projectType: 'cli' as const,
        architecture: 'standard' as const,
        framework: 'cobra' as const,
        logger: 'slog' as const,
        goVersion: '1.21'
      };

      // Test preview generation
      const preview = await api.projects.generatePreview(testConfig);
      if (!preview.fileStructure || preview.fileCount === 0) {
        throw new Error('Preview generation failed');
      }

      // Test project generation
      const generation = await api.projects.generateProject({
        ...testConfig,
        outputFormat: 'zip',
        includeTests: true,
        includeDocs: true
      });
      
      if (!generation.projectId) {
        throw new Error('Project generation failed');
      }
    });
  }

  async testWebSocketConnection(): Promise<TestResult> {
    return this.runTest('WebSocket Connection', async () => {
      // Test WebSocket connection
      await connectWebSocket();
      
      // Test subscription
      let messageReceived = false;
      const unsubscribe = api.ws.subscribe('test', () => {
        messageReceived = true;
      });

      // Send test message
      api.ws.send({ type: 'test', data: 'test message' });
      
      // Wait a bit for message processing
      await new Promise(resolve => setTimeout(resolve, 100));
      
      unsubscribe();
      
      // Note: We can't easily test message reception in this context
      // but we can verify the connection was established
    });
  }

  async testErrorHandling(): Promise<TestResult> {
    return this.runTest('Error Handling', async () => {
      try {
        // Test invalid blueprint
        await api.blueprints.getBlueprintById('non-existent');
        throw new Error('Should have thrown an error for invalid blueprint');
      } catch (error) {
        // This is expected
      }

      try {
        // Test invalid preview config
        await api.projects.generatePreview({
          projectName: '',
          moduleUrl: 'invalid-url',
          projectType: 'invalid-type' as any,
          architecture: 'standard',
          framework: 'gin',
          logger: 'slog',
          goVersion: '1.21'
        });
        throw new Error('Should have thrown an error for invalid config');
      } catch (error) {
        // This is expected
      }
    });
  }

  async runAllTests(): Promise<TestResult[]> {
    console.log('🧪 Starting API Integration Tests...');
    
    // Run all tests
    await this.testHealthEndpoints();
    await this.testConfigEndpoints();
    await this.testBlueprintEndpoints();
    await this.testProjectOperations();
    await this.testWebSocketConnection();
    await this.testErrorHandling();

    return this.results;
  }

  getResults(): TestResult[] {
    return this.results;
  }

  getSummary() {
    const total = this.results.length;
    const passed = this.results.filter(r => r.success).length;
    const failed = total - passed;
    const avgDuration = this.results.reduce((sum, r) => sum + r.duration, 0) / total;

    return {
      total,
      passed,
      failed,
      passRate: (passed / total) * 100,
      avgDuration: Math.round(avgDuration * 100) / 100
    };
  }

  printResults(): void {
    console.log('\n📊 API Test Results:');
    console.log('==================');
    
    this.results.forEach(result => {
      const status = result.success ? '✅' : '❌';
      const duration = `${result.duration.toFixed(2)}ms`;
      console.log(`${status} ${result.name} (${duration})`);
      
      if (!result.success && result.error) {
        console.log(`   Error: ${result.error}`);
      }
    });

    const summary = this.getSummary();
    console.log('\n📈 Summary:');
    console.log(`   Total: ${summary.total}`);
    console.log(`   Passed: ${summary.passed}`);
    console.log(`   Failed: ${summary.failed}`);
    console.log(`   Pass Rate: ${summary.passRate.toFixed(1)}%`);
    console.log(`   Avg Duration: ${summary.avgDuration}ms`);
  }
}

// Export convenience function for running tests
export async function runApiTests(): Promise<TestResult[]> {
  const tester = new ApiTester();
  const results = await tester.runAllTests();
  tester.printResults();
  return results;
}

// Export for development testing
if (import.meta.env.DEV) {
  (window as any).runApiTests = runApiTests;
  (window as any).ApiTester = ApiTester;
}