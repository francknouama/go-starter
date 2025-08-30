const puppeteer = require('puppeteer');
const fs = require('fs');
const path = require('path');

async function generateScreenshots() {
  const browser = await puppeteer.launch({ headless: true });
  const baseUrl = 'http://localhost:5173';
  
  const screenshots = [
    // Desktop screenshots (1920x1080)
    {
      name: 'landing-page-desktop',
      path: '/docs/screenshots/desktop/landing-page.png',
      viewport: { width: 1920, height: 1080 },
      url: baseUrl
    },
    {
      name: 'blueprint-gallery-desktop',
      path: '/docs/screenshots/desktop/blueprint-gallery.png', 
      viewport: { width: 1920, height: 1080 },
      url: baseUrl
    },
    {
      name: 'configuration-form-desktop',
      path: '/docs/screenshots/desktop/configuration-form.png',
      viewport: { width: 1920, height: 1080 },
      url: baseUrl
    },
    // Mobile screenshots (375x667)
    {
      name: 'landing-page-mobile',
      path: '/docs/screenshots/mobile/landing-page.png',
      viewport: { width: 375, height: 667 },
      url: baseUrl
    },
    {
      name: 'blueprint-gallery-mobile',
      path: '/docs/screenshots/mobile/blueprint-gallery.png',
      viewport: { width: 375, height: 667 },
      url: baseUrl
    },
    // Tablet screenshots (768x1024)
    {
      name: 'landing-page-tablet',
      path: '/docs/screenshots/tablet/landing-page.png',
      viewport: { width: 768, height: 1024 },
      url: baseUrl
    }
  ];

  for (const screenshot of screenshots) {
    console.log(`Generating ${screenshot.name}...`);
    const page = await browser.newPage();
    await page.setViewport(screenshot.viewport);
    
    try {
      await page.goto(screenshot.url, { waitUntil: 'networkidle0', timeout: 10000 });
      
      // Wait for React components to render
      await page.waitForTimeout(2000);
      
      // Take screenshot
      const screenshotPath = path.join(__dirname, '../', screenshot.path);
      await page.screenshot({ 
        path: screenshotPath, 
        fullPage: false,
        quality: 90 
      });
      
      console.log(`✓ Generated ${screenshot.name}`);
    } catch (error) {
      console.error(`✗ Failed to generate ${screenshot.name}:`, error.message);
    }
    
    await page.close();
  }

  await browser.close();
  console.log('Screenshot generation completed!');
}

generateScreenshots().catch(console.error);