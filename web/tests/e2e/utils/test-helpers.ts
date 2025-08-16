import { Page, expect } from '@playwright/test';

export class GoStarterPageObject {
  constructor(private page: Page) {}

  // Navigation helpers
  async goto() {
    await this.page.goto('/');
    await this.waitForPageLoad();
  }

  async waitForPageLoad() {
    await this.page.waitForLoadState('networkidle');
    await expect(this.page.getByRole('heading', { name: /go-starter/i })).toBeVisible();
  }

  // Header interactions
  async toggleDisclosureMode(mode: 'basic' | 'advanced') {
    const currentMode = await this.getCurrentDisclosureMode();
    if (currentMode !== mode) {
      await this.page.getByRole('button', { name: /basic|advanced/i }).click();
    }
    await expect(this.page.getByText(mode)).toBeVisible();
  }

  async getCurrentDisclosureMode(): Promise<'basic' | 'advanced'> {
    const modeText = await this.page.getByTestId('disclosure-mode-indicator').textContent();
    return modeText?.toLowerCase().includes('advanced') ? 'advanced' : 'basic';
  }

  // Configuration Panel interactions
  async fillProjectName(name: string) {
    await this.page.getByLabel(/project name/i).fill(name);
  }

  async selectProjectType(type: string) {
    await this.page.getByLabel(/project type/i).selectOption(type);
  }

  async selectFramework(framework: string) {
    await this.page.getByLabel(/framework/i).selectOption(framework);
  }

  async selectLogger(logger: string) {
    await this.page.getByLabel(/logger/i).selectOption(logger);
  }

  async selectArchitecture(architecture: string) {
    await this.page.getByLabel(/architecture/i).selectOption(architecture);
  }

  async fillModulePath(path: string) {
    await this.page.getByLabel(/module path/i).fill(path);
  }

  // Advanced options (only visible in advanced mode)
  async selectDatabaseDriver(driver: string) {
    await this.page.getByLabel(/database driver/i).selectOption(driver);
  }

  async selectAuthType(authType: string) {
    await this.page.getByLabel(/authentication/i).selectOption(authType);
  }

  // Preview Panel interactions
  async waitForPreviewUpdate() {
    await this.page.waitForTimeout(1000); // Wait for preview to update
    await this.page.waitForLoadState('networkidle');
  }

  async getPreviewContent(): Promise<string> {
    await this.waitForPreviewUpdate();
    return await this.page.getByTestId('preview-content').textContent() || '';
  }

  async getEstimatedFiles(): Promise<number> {
    const filesText = await this.page.getByTestId('estimated-files').textContent();
    return parseInt(filesText?.match(/\d+/)?.[0] || '0');
  }

  // File Explorer interactions
  async expandFolder(folderName: string) {
    await this.page.getByRole('button', { name: folderName }).click();
  }

  async selectFile(fileName: string) {
    await this.page.getByText(fileName).click();
  }

  async getFileList(): Promise<string[]> {
    const fileElements = this.page.getByTestId('file-item');
    return await fileElements.allTextContents();
  }

  // Generation actions
  async generateProject() {
    await this.page.getByRole('button', { name: /generate project/i }).click();
  }

  async downloadProject() {
    const downloadPromise = this.page.waitForDownload();
    await this.page.getByRole('button', { name: /download/i }).click();
    return await downloadPromise;
  }

  // WebSocket and real-time updates
  async waitForWebSocketConnection() {
    // Wait for WebSocket connection to be established
    await this.page.waitForFunction(() => {
      return window.WebSocket && 
             document.querySelector('[data-testid="websocket-status"]')?.textContent === 'connected';
    });
  }

  async waitForRealTimeUpdate() {
    // Wait for real-time preview updates via WebSocket
    await this.page.waitForFunction(() => {
      return document.querySelector('[data-testid="preview-content"]')?.textContent !== '';
    });
  }

  // Form validation helpers
  async getValidationErrors(): Promise<string[]> {
    const errorElements = this.page.getByTestId('validation-error');
    return await errorElements.allTextContents();
  }

  async hasValidationError(field: string): Promise<boolean> {
    return await this.page.getByTestId(`${field}-error`).isVisible();
  }

  // Responsive design helpers
  async setViewportSize(width: number, height: number) {
    await this.page.setViewportSize({ width, height });
  }

  async isMobileViewport(): Promise<boolean> {
    const viewport = this.page.viewportSize();
    return viewport ? viewport.width < 768 : false;
  }

  // Advanced configuration helpers
  async openAdvancedSettings() {
    await this.page.getByRole('button', { name: /advanced settings/i }).click();
  }

  async enableFeature(feature: string) {
    await this.page.getByLabel(feature).check();
  }

  async disableFeature(feature: string) {
    await this.page.getByLabel(feature).uncheck();
  }

  // AI Advisor integration (if available)
  async openAIAdvisor() {
    await this.page.getByRole('button', { name: /ai advisor/i }).click();
  }

  async fillAdvisorRequirements(requirements: {
    domain?: string;
    teamSize?: number;
    experience?: string;
    timeline?: string;
  }) {
    if (requirements.domain) {
      await this.page.getByLabel(/domain/i).selectOption(requirements.domain);
    }
    if (requirements.teamSize) {
      await this.page.getByLabel(/team size/i).fill(requirements.teamSize.toString());
    }
    if (requirements.experience) {
      await this.page.getByLabel(/experience/i).selectOption(requirements.experience);
    }
    if (requirements.timeline) {
      await this.page.getByLabel(/timeline/i).selectOption(requirements.timeline);
    }
  }

  async getAIRecommendation(): Promise<{
    blueprint: string;
    confidence: number;
    reasoning: string[];
  }> {
    await this.page.getByRole('button', { name: /get recommendation/i }).click();
    await this.page.waitForSelector('[data-testid="ai-recommendation"]');
    
    const blueprint = await this.page.getByTestId('recommended-blueprint').textContent() || '';
    const confidenceText = await this.page.getByTestId('recommendation-confidence').textContent() || '0%';
    const confidence = parseInt(confidenceText.replace('%', ''));
    const reasoningElements = this.page.getByTestId('recommendation-reason');
    const reasoning = await reasoningElements.allTextContents();
    
    return { blueprint, confidence, reasoning };
  }

  // Progressive disclosure helpers
  async getVisibleFields(): Promise<string[]> {
    const fields = this.page.getByRole('textbox, combobox, checkbox');
    const visibleFields = [];
    const count = await fields.count();
    
    for (let i = 0; i < count; i++) {
      const field = fields.nth(i);
      if (await field.isVisible()) {
        const label = await field.getAttribute('aria-label') || 
                     await field.getAttribute('placeholder') || 
                     await field.getAttribute('name') || '';
        if (label) {
          visibleFields.push(label);
        }
      }
    }
    
    return visibleFields;
  }

  // Performance helpers
  async measurePageLoadTime(): Promise<number> {
    const startTime = Date.now();
    await this.goto();
    const endTime = Date.now();
    return endTime - startTime;
  }

  async measurePreviewUpdateTime(): Promise<number> {
    const startTime = Date.now();
    await this.fillProjectName(`test-${Date.now()}`);
    await this.waitForPreviewUpdate();
    const endTime = Date.now();
    return endTime - startTime;
  }
}

// Utility functions for common test scenarios
export async function createStandardWebAPIProject(page: GoStarterPageObject) {
  await page.fillProjectName('my-web-api');
  await page.selectProjectType('web-api');
  await page.selectFramework('gin');
  await page.selectLogger('slog');
  await page.fillModulePath('github.com/user/my-web-api');
}

export async function createAdvancedProject(page: GoStarterPageObject) {
  await page.toggleDisclosureMode('advanced');
  await page.fillProjectName('enterprise-api');
  await page.selectProjectType('web-api');
  await page.selectArchitecture('clean');
  await page.selectFramework('gin');
  await page.selectLogger('zap');
  await page.selectDatabaseDriver('postgres');
  await page.selectAuthType('jwt');
  await page.fillModulePath('github.com/enterprise/api');
}

export async function createCLIProject(page: GoStarterPageObject) {
  await page.fillProjectName('my-cli-tool');
  await page.selectProjectType('cli');
  await page.selectFramework('cobra');
  await page.selectLogger('slog');
  await page.fillModulePath('github.com/user/my-cli-tool');
}

// Test data generators
export function generateRandomProjectName(): string {
  return `test-project-${Date.now()}-${Math.random().toString(36).substring(7)}`;
}

export function generateRandomModulePath(): string {
  const username = `user${Math.random().toString(36).substring(7)}`;
  const projectName = generateRandomProjectName();
  return `github.com/${username}/${projectName}`;
}