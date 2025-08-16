#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

console.log('🧪 Validating E2E Test Setup...\n');

// Check if all required files exist
const requiredFiles = [
  'playwright.config.ts',
  'tests/e2e/utils/test-helpers.ts',
  'tests/e2e/app-navigation.spec.ts',
  'tests/e2e/project-configuration.spec.ts',
  'tests/e2e/real-time-preview.spec.ts',
  'tests/e2e/project-generation.spec.ts',
  'tests/e2e/ai-advisor.spec.ts',
  'tests/e2e/performance-accessibility.spec.ts',
  'tests/e2e/test-setup.ts'
];

let allFilesExist = true;

console.log('📁 Checking required files:');
requiredFiles.forEach(file => {
  const filePath = path.join(__dirname, file);
  const exists = fs.existsSync(filePath);
  console.log(`   ${exists ? '✅' : '❌'} ${file}`);
  if (!exists) allFilesExist = false;
});

// Check package.json for required dependencies and scripts
console.log('\n📦 Checking package.json:');
const packagePath = path.join(__dirname, 'package.json');
if (fs.existsSync(packagePath)) {
  const packageJson = JSON.parse(fs.readFileSync(packagePath, 'utf8'));
  
  // Check for Playwright dependency
  const hasPlaywright = packageJson.devDependencies && packageJson.devDependencies['@playwright/test'];
  console.log(`   ${hasPlaywright ? '✅' : '❌'} @playwright/test dependency`);
  
  // Check for E2E scripts
  const requiredScripts = ['test:e2e', 'test:e2e:ui', 'test:e2e:headed', 'test:e2e:debug'];
  requiredScripts.forEach(script => {
    const hasScript = packageJson.scripts && packageJson.scripts[script];
    console.log(`   ${hasScript ? '✅' : '❌'} ${script} script`);
  });
} else {
  console.log('   ❌ package.json not found');
  allFilesExist = false;
}

// Check CI workflow
console.log('\n🔄 Checking CI configuration:');
const ciPath = path.join(__dirname, '../.github/workflows/e2e-tests.yml');
const hasCIWorkflow = fs.existsSync(ciPath);
console.log(`   ${hasCIWorkflow ? '✅' : '❌'} E2E CI workflow`);

// Check documentation
console.log('\n📚 Checking documentation:');
const docsPath = path.join(__dirname, 'README-E2E-Testing.md');
const hasDocs = fs.existsSync(docsPath);
console.log(`   ${hasDocs ? '✅' : '❌'} E2E documentation`);

// Summary
console.log('\n' + '='.repeat(50));
if (allFilesExist && hasCIWorkflow && hasDocs) {
  console.log('🎉 E2E Test Setup Complete!');
  console.log('\nNext steps:');
  console.log('1. Install Playwright browsers: npx playwright install');
  console.log('2. Start the web server: go run cmd/web-server/main.go');
  console.log('3. Start the frontend: npm run dev');
  console.log('4. Run tests: npm run test:e2e');
  process.exit(0);
} else {
  console.log('❌ E2E Test Setup Incomplete');
  console.log('\nSome required files or configurations are missing.');
  process.exit(1);
}