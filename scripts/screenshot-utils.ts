/**
 * Utility functions and helpers for screenshot automation
 * 
 * This module provides reusable utilities for screenshot generation,
 * image processing, and documentation integration.
 */

import { Page, Browser } from '@playwright/test';
import { promises as fs } from 'fs';
import path from 'path';

// Screenshot metadata interface
export interface ScreenshotMetadata {
  filename: string;
  description: string;
  timestamp: string;
  viewport: { width: number; height: number };
  category: string;
  tags: string[];
  fileSize?: number;
  dimensions?: { width: number; height: number };
}

// Screenshot categories for organization
export const SCREENSHOT_CATEGORIES = {
  INTERFACE: 'interface',
  WORKFLOW: 'workflow', 
  FEATURE: 'feature',
  RESPONSIVE: 'responsive',
  BLUEPRINT: 'blueprint',
  STATE: 'state',
  ACCESSIBILITY: 'accessibility'
} as const;

// Common selectors for go-starter Web UI
export const UI_SELECTORS = {
  // Main panels
  CONFIGURATION_PANEL: '[data-testid="configuration-panel"], .configuration-panel, form',
  PREVIEW_PANEL: '[data-testid="preview-panel"], .preview-panel, .code-preview',
  FILE_EXPLORER_PANEL: '[data-testid="file-explorer-panel"], .file-explorer, .files',
  
  // Form elements
  PROJECT_NAME: 'input[name="name"], [data-testid="project-name"]',
  PROJECT_TYPE: 'select[name="type"], [data-testid="project-type"]',
  MODULE_PATH: 'input[name="module"], [data-testid="module-path"]',
  FRAMEWORK: 'select[name="framework"], [data-testid="framework"]',
  LOGGER: 'select[name="logger"], [data-testid="logger"]',
  
  // Advanced options
  DATABASE: 'select[name="database"], [data-testid="database"]',
  AUTH_TYPE: 'select[name="auth"], [data-testid="auth"]',
  ARCHITECTURE: 'select[name="architecture"], [data-testid="architecture"]',
  
  // Buttons and controls
  GENERATE_BUTTON: '[data-testid="generate-button"], button[type="submit"]',
  ADVANCED_TOGGLE: '[data-testid="advanced-toggle"], .advanced-toggle',
  THEME_TOGGLE: '[data-testid="theme-toggle"], .theme-toggle, .dark-mode-toggle',
  
  // Navigation
  HEADER: 'header, [data-testid="header"]',
  SIDEBAR: '.sidebar, [data-testid="sidebar"]',
  NAVIGATION: 'nav, [data-testid="navigation"]',
  
  // Blueprint gallery
  BLUEPRINT_GALLERY: '.blueprint-gallery, [data-testid="blueprint-selection"]',
  BLUEPRINT_CARD: '.blueprint-card, [data-testid="blueprint-card"]',
  
  // Loading and state indicators
  LOADING_SPINNER: '.loading, .spinner, [data-testid="loading"]',
  ERROR_MESSAGE: '.error, .error-message, [data-testid="error"]',
  SUCCESS_MESSAGE: '.success, .success-message, [data-testid="success"]'
} as const;

/**
 * Page interaction utilities for consistent screenshot setup
 */
export class PageInteractionUtils {
  constructor(private page: Page) {}

  /**
   * Wait for the go-starter app to be fully loaded and ready
   */
  async waitForAppReady(timeout: number = 10000): Promise<void> {
    // Wait for main panels to be visible
    await this.page.waitForFunction(() => {
      return document.querySelector('[data-testid="configuration-panel"]') !== null ||
             document.querySelector('.configuration-panel') !== null ||
             document.querySelector('form') !== null;
    }, { timeout });

    // Wait for React hydration and any animations to settle
    await this.page.waitForLoadState('networkidle');
    await this.page.waitForTimeout(1000);

    // Ensure no loading spinners are visible
    await this.page.waitForFunction(() => {
      const spinners = document.querySelectorAll('.loading, .spinner, [data-testid="loading"]');
      return spinners.length === 0 || Array.from(spinners).every(s => 
        (s as HTMLElement).style.display === 'none' ||
        !(s as HTMLElement).offsetParent
      );
    }, { timeout: 5000 }).catch(() => {
      console.warn('Loading spinners still visible, proceeding anyway');
    });
  }

  /**
   * Disable all animations for consistent screenshots
   */
  async disableAnimations(): Promise<void> {
    await this.page.addInitScript(() => {
      const style = document.createElement('style');
      style.textContent = `
        *, *::before, *::after {
          animation-duration: 0.01ms !important;
          animation-delay: 0.01ms !important;
          transition-duration: 0.01ms !important;
          transition-delay: 0.01ms !important;
          animation-fill-mode: both !important;
        }
        
        /* Disable specific animations */
        .animate-pulse, .animate-spin, .animate-ping, .animate-bounce {
          animation: none !important;
        }
        
        /* Disable Tailwind animations */
        [class*="animate-"] {
          animation: none !important;
        }
        
        /* Disable Framer Motion animations */
        [style*="transform"], [style*="opacity"] {
          transition: none !important;
        }
      `;
      document.head.appendChild(style);
    });
  }

  /**
   * Fill form field if it exists
   */
  async fillIfExists(selector: string, value: string, waitTime: number = 200): Promise<boolean> {
    try {
      const element = await this.page.$(selector);
      if (element) {
        await element.fill(value);
        await this.page.waitForTimeout(waitTime);
        return true;
      }
      return false;
    } catch (error) {
      console.warn(`Could not fill ${selector}: ${error}`);
      return false;
    }
  }

  /**
   * Select option if element exists
   */
  async selectIfExists(selector: string, value: string, waitTime: number = 200): Promise<boolean> {
    try {
      const element = await this.page.$(selector);
      if (element) {
        await element.selectOption(value);
        await this.page.waitForTimeout(waitTime);
        return true;
      }
      return false;
    } catch (error) {
      console.warn(`Could not select ${selector}: ${error}`);
      return false;
    }
  }

  /**
   * Click element if it exists
   */
  async clickIfExists(selector: string, waitTime: number = 500): Promise<boolean> {
    try {
      const element = await this.page.$(selector);
      if (element && await element.isVisible()) {
        await element.click();
        await this.page.waitForTimeout(waitTime);
        return true;
      }
      return false;
    } catch (error) {
      console.warn(`Could not click ${selector}: ${error}`);
      return false;
    }
  }

  /**
   * Scroll element into view if it exists
   */
  async scrollIntoViewIfExists(selector: string): Promise<boolean> {
    try {
      const element = await this.page.$(selector);
      if (element) {
        await element.scrollIntoViewIfNeeded();
        await this.page.waitForTimeout(300);
        return true;
      }
      return false;
    } catch (error) {
      console.warn(`Could not scroll ${selector}: ${error}`);
      return false;
    }
  }

  /**
   * Get element bounding box if it exists
   */
  async getBoundingBoxIfExists(selector: string) {
    try {
      const element = await this.page.$(selector);
      if (element) {
        return await element.boundingBox();
      }
      return null;
    } catch (error) {
      console.warn(`Could not get bounding box for ${selector}: ${error}`);
      return null;
    }
  }

  /**
   * Fill a complete project form with validation
   */
  async fillProjectForm(config: {
    name: string;
    type: string;
    module: string;
    framework?: string;
    logger?: string;
    database?: string;
    auth?: string;
    architecture?: string;
  }): Promise<void> {
    // Fill basic required fields
    await this.fillIfExists(UI_SELECTORS.PROJECT_NAME, config.name);
    await this.selectIfExists(UI_SELECTORS.PROJECT_TYPE, config.type);
    await this.fillIfExists(UI_SELECTORS.MODULE_PATH, config.module);

    // Fill optional fields if provided
    if (config.framework) {
      await this.selectIfExists(UI_SELECTORS.FRAMEWORK, config.framework);
    }
    if (config.logger) {
      await this.selectIfExists(UI_SELECTORS.LOGGER, config.logger);
    }
    if (config.database) {
      await this.selectIfExists(UI_SELECTORS.DATABASE, config.database);
    }
    if (config.auth) {
      await this.selectIfExists(UI_SELECTORS.AUTH_TYPE, config.auth);
    }
    if (config.architecture) {
      await this.selectIfExists(UI_SELECTORS.ARCHITECTURE, config.architecture);
    }

    // Wait for form to process changes
    await this.page.waitForTimeout(1000);
  }
}

/**
 * Screenshot capture utilities with metadata tracking
 */
export class ScreenshotCapture {
  private metadata: ScreenshotMetadata[] = [];
  
  constructor(private page: Page, private outputDir: string) {}

  /**
   * Capture full page screenshot with metadata
   */
  async captureFullPage(
    filename: string,
    description: string,
    category: keyof typeof SCREENSHOT_CATEGORIES,
    tags: string[] = []
  ): Promise<ScreenshotMetadata> {
    const filePath = path.join(this.outputDir, filename);
    const viewport = this.page.viewportSize() || { width: 1920, height: 1080 };
    
    await this.page.screenshot({
      path: filePath,
      fullPage: true,
      type: 'png',
      quality: 90
    });

    // Get file stats
    const stats = await fs.stat(filePath);
    
    const metadata: ScreenshotMetadata = {
      filename,
      description,
      timestamp: new Date().toISOString(),
      viewport,
      category: SCREENSHOT_CATEGORIES[category],
      tags,
      fileSize: stats.size
    };

    this.metadata.push(metadata);
    console.log(`  ✅ ${description} -> ${filename}`);
    
    return metadata;
  }

  /**
   * Capture element screenshot with metadata
   */
  async captureElement(
    selector: string,
    filename: string,
    description: string,
    category: keyof typeof SCREENSHOT_CATEGORIES,
    tags: string[] = []
  ): Promise<ScreenshotMetadata | null> {
    const element = await this.page.$(selector);
    if (!element) {
      console.warn(`  ⚠️  Element not found: ${selector}`);
      return null;
    }

    const filePath = path.join(this.outputDir, filename);
    const viewport = this.page.viewportSize() || { width: 1920, height: 1080 };
    
    await element.screenshot({
      path: filePath,
      type: 'png',
      quality: 90
    });

    const stats = await fs.stat(filePath);
    
    const metadata: ScreenshotMetadata = {
      filename,
      description,
      timestamp: new Date().toISOString(),
      viewport,
      category: SCREENSHOT_CATEGORIES[category],
      tags,
      fileSize: stats.size
    };

    this.metadata.push(metadata);
    console.log(`  ✅ ${description} -> ${filename}`);
    
    return metadata;
  }

  /**
   * Get all captured screenshot metadata
   */
  getMetadata(): ScreenshotMetadata[] {
    return [...this.metadata];
  }

  /**
   * Export metadata to JSON file
   */
  async exportMetadata(filename: string = 'screenshots-metadata.json'): Promise<void> {
    const filePath = path.join(this.outputDir, filename);
    await fs.writeFile(filePath, JSON.stringify(this.metadata, null, 2));
    console.log(`  📝 Metadata exported to ${filename}`);
  }

  /**
   * Generate HTML gallery of screenshots
   */
  async generateHtmlGallery(filename: string = 'gallery.html'): Promise<void> {
    const groupedMetadata = this.metadata.reduce((groups, meta) => {
      if (!groups[meta.category]) {
        groups[meta.category] = [];
      }
      groups[meta.category].push(meta);
      return groups;
    }, {} as Record<string, ScreenshotMetadata[]>);

    const htmlContent = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>go-starter Web UI Screenshot Gallery</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; margin: 0; padding: 20px; background: #f5f5f5; }
        .container { max-width: 1200px; margin: 0 auto; background: white; border-radius: 12px; padding: 30px; box-shadow: 0 4px 6px rgba(0,0,0,0.1); }
        h1 { color: #2563eb; margin-bottom: 30px; }
        h2 { color: #374151; border-bottom: 2px solid #e5e7eb; padding-bottom: 10px; }
        .category { margin-bottom: 40px; }
        .screenshots { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 20px; }
        .screenshot { border: 1px solid #e5e7eb; border-radius: 8px; overflow: hidden; background: white; }
        .screenshot img { width: 100%; height: auto; display: block; }
        .screenshot-info { padding: 15px; }
        .screenshot-title { font-weight: 600; margin-bottom: 5px; color: #111827; }
        .screenshot-desc { color: #6b7280; font-size: 14px; margin-bottom: 10px; }
        .screenshot-meta { font-size: 12px; color: #9ca3af; }
        .tags { display: flex; gap: 5px; flex-wrap: wrap; margin-top: 8px; }
        .tag { background: #eff6ff; color: #2563eb; padding: 2px 8px; border-radius: 12px; font-size: 11px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>go-starter Web UI Screenshot Gallery</h1>
        <p>Generated on: ${new Date().toLocaleString()}</p>
        
        ${Object.entries(groupedMetadata).map(([category, screenshots]) => `
        <div class="category">
            <h2>${category.charAt(0).toUpperCase() + category.slice(1)} Screenshots</h2>
            <div class="screenshots">
                ${screenshots.map(meta => `
                <div class="screenshot">
                    <img src="./${meta.filename}" alt="${meta.description}" loading="lazy">
                    <div class="screenshot-info">
                        <div class="screenshot-title">${path.basename(meta.filename, path.extname(meta.filename))}</div>
                        <div class="screenshot-desc">${meta.description}</div>
                        <div class="screenshot-meta">
                            ${meta.viewport.width}x${meta.viewport.height} • ${(meta.fileSize || 0 / 1024).toFixed(1)}KB
                        </div>
                        <div class="tags">
                            ${meta.tags.map(tag => `<span class="tag">${tag}</span>`).join('')}
                        </div>
                    </div>
                </div>
                `).join('')}
            </div>
        </div>
        `).join('')}
    </div>
</body>
</html>
    `.trim();

    const filePath = path.join(this.outputDir, filename);
    await fs.writeFile(filePath, htmlContent);
    console.log(`  🎨 HTML gallery generated: ${filename}`);
  }
}

/**
 * Cross-browser screenshot utilities
 */
export class CrossBrowserScreenshot {
  private browsers: { name: string; browser: Browser }[] = [];

  async addBrowser(name: string, browser: Browser): Promise<void> {
    this.browsers.push({ name, browser });
  }

  /**
   * Capture screenshot across all browsers for comparison
   */
  async captureAcrossBrowsers(
    url: string,
    filename: string,
    description: string,
    outputDir: string
  ): Promise<void> {
    console.log(`📷 Capturing cross-browser: ${description}`);

    for (const { name, browser } of this.browsers) {
      const context = await browser.newContext();
      const page = await context.newPage();
      
      await page.goto(url, { waitUntil: 'networkidle' });
      
      const utils = new PageInteractionUtils(page);
      await utils.waitForAppReady();
      
      const browserFilename = filename.replace('.png', `-${name}.png`);
      const filePath = path.join(outputDir, 'cross-browser', browserFilename);
      
      await page.screenshot({
        path: filePath,
        fullPage: true,
        type: 'png',
        quality: 90
      });

      console.log(`  ✅ ${name}: ${browserFilename}`);
      
      await context.close();
    }
  }

  async cleanup(): Promise<void> {
    for (const { browser } of this.browsers) {
      await browser.close();
    }
  }
}

/**
 * Performance monitoring utilities for screenshot generation
 */
export class ScreenshotPerformanceMonitor {
  private startTime?: number;
  private metrics: Array<{ name: string; duration: number; timestamp: string }> = [];

  startOperation(name: string): void {
    this.startTime = Date.now();
    console.log(`⏱️  Starting: ${name}`);
  }

  endOperation(name: string): number {
    if (!this.startTime) {
      console.warn('No start time recorded');
      return 0;
    }

    const duration = Date.now() - this.startTime;
    this.metrics.push({
      name,
      duration,
      timestamp: new Date().toISOString()
    });

    console.log(`✅ Completed: ${name} (${duration}ms)`);
    this.startTime = undefined;
    return duration;
  }

  getMetrics() {
    return [...this.metrics];
  }

  async exportMetrics(outputDir: string): Promise<void> {
    const filePath = path.join(outputDir, 'performance-metrics.json');
    await fs.writeFile(filePath, JSON.stringify({
      totalOperations: this.metrics.length,
      totalTime: this.metrics.reduce((sum, m) => sum + m.duration, 0),
      averageTime: this.metrics.reduce((sum, m) => sum + m.duration, 0) / this.metrics.length,
      operations: this.metrics
    }, null, 2));
    console.log('📊 Performance metrics exported');
  }
}