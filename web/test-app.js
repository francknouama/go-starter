#!/usr/bin/env node

// Simple test to check if the React app loads without infinite re-render loops
const puppeteer = require('puppeteer');

async function testApp() {
  let browser;
  try {
    console.log('Launching browser...');
    browser = await puppeteer.launch({ 
      headless: true,
      args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    const page = await browser.newPage();
    
    // Listen for console messages
    const consoleMessages = [];
    page.on('console', msg => {
      consoleMessages.push({
        type: msg.type(),
        text: msg.text()
      });
    });
    
    // Listen for errors
    const errors = [];
    page.on('error', err => {
      errors.push(err.message);
    });
    
    page.on('pageerror', err => {
      errors.push(err.message);
    });
    
    console.log('Navigating to app...');
    await page.goto('http://localhost:5173', { 
      waitUntil: 'networkidle0',
      timeout: 10000
    });
    
    // Wait for initial render
    await page.waitForTimeout(2000);
    
    // Check if the app is stuck on loading
    const content = await page.content();
    const isStuckLoading = content.includes('Loading...') && !content.includes('blueprint-selection');
    
    // Get current page title
    const title = await page.title();
    
    // Check for React errors
    const reactErrors = consoleMessages.filter(msg => 
      msg.text.includes('React') || 
      msg.text.includes('Maximum update depth') ||
      msg.text.includes('commitHookPassiveUnmountEffects') ||
      msg.type === 'error'
    );
    
    console.log('\n=== TEST RESULTS ===');
    console.log(`Page Title: ${title}`);
    console.log(`Stuck on Loading: ${isStuckLoading}`);
    console.log(`Total Console Messages: ${consoleMessages.length}`);
    console.log(`Errors: ${errors.length}`);
    console.log(`React Errors: ${reactErrors.length}`);
    
    if (reactErrors.length > 0) {
      console.log('\n=== REACT ERRORS ===');
      reactErrors.slice(0, 5).forEach((msg, i) => {
        console.log(`${i + 1}. [${msg.type}] ${msg.text}`);
      });
    }
    
    if (errors.length > 0) {
      console.log('\n=== PAGE ERRORS ===');
      errors.slice(0, 5).forEach((err, i) => {
        console.log(`${i + 1}. ${err}`);
      });
    }
    
    // Success criteria
    const success = !isStuckLoading && reactErrors.length === 0 && errors.length === 0;
    console.log(`\n=== OVERALL STATUS: ${success ? 'PASS' : 'FAIL'} ===`);
    
    return success;
    
  } catch (error) {
    console.error('Test failed:', error.message);
    return false;
  } finally {
    if (browser) {
      await browser.close();
    }
  }
}

testApp().then(success => {
  process.exit(success ? 0 : 1);
});