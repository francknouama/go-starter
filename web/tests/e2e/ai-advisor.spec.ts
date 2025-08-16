import { test, expect } from '@playwright/test';
import { GoStarterPageObject } from './utils/test-helpers';

test.describe('AI Architecture Advisor Integration', () => {
  let page: GoStarterPageObject;

  test.beforeEach(async ({ page: browserPage }) => {
    page = new GoStarterPageObject(browserPage);
    await page.goto();
  });

  test('should open AI advisor modal', async () => {
    await page.openAIAdvisor();
    
    // AI advisor modal should be visible
    await expect(page.page.getByTestId('ai-advisor-modal')).toBeVisible();
    await expect(page.page.getByRole('heading', { name: /ai.*advisor/i })).toBeVisible();
    
    // Should show requirement gathering form
    await expect(page.page.getByLabel(/domain/i)).toBeVisible();
    await expect(page.page.getByLabel(/team size/i)).toBeVisible();
    await expect(page.page.getByLabel(/experience/i)).toBeVisible();
  });

  test('should provide recommendation for e-commerce project', async () => {
    await page.openAIAdvisor();
    
    // Fill requirements for e-commerce project
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 5,
      experience: 'senior',
      timeline: 'standard'
    });
    
    // Get recommendation
    const recommendation = await page.getAIRecommendation();
    
    // Should recommend appropriate architecture for e-commerce
    expect(['web-api-ddd', 'web-api-clean', 'microservice']).toContain(recommendation.blueprint);
    expect(recommendation.confidence).toBeGreaterThan(70);
    expect(recommendation.reasoning.length).toBeGreaterThan(0);
    expect(recommendation.reasoning.some(r => r.includes('e-commerce'))).toBe(true);
  });

  test('should recommend simple architecture for junior team', async () => {
    await page.openAIAdvisor();
    
    // Fill requirements for junior team
    await page.fillAdvisorRequirements({
      domain: 'devtools',
      teamSize: 2,
      experience: 'junior',
      timeline: 'mvp'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Should recommend simpler architectures for junior teams
    expect(['cli-simple', 'library', 'web-api']).toContain(recommendation.blueprint);
    expect(recommendation.confidence).toBeGreaterThan(60);
    expect(recommendation.reasoning.some(r => r.includes('junior') || r.includes('simple'))).toBe(true);
  });

  test('should recommend complex architecture for expert team', async () => {
    await page.openAIAdvisor();
    
    // Fill requirements for expert team with complex needs
    await page.fillAdvisorRequirements({
      domain: 'fintech',
      teamSize: 8,
      experience: 'expert',
      timeline: 'thorough'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Should recommend complex architectures for expert teams
    expect(['web-api-ddd', 'web-api-hexagonal', 'microservice', 'event-driven']).toContain(recommendation.blueprint);
    expect(recommendation.confidence).toBeGreaterThan(75);
    expect(recommendation.reasoning.some(r => r.includes('expert') || r.includes('complex'))).toBe(true);
  });

  test('should apply recommendation to configuration form', async () => {
    await page.openAIAdvisor();
    
    // Get recommendation
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 4,
      experience: 'mixed',
      timeline: 'standard'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Apply recommendation
    await page.page.getByRole('button', { name: /apply recommendation/i }).click();
    
    // Modal should close
    await expect(page.page.getByTestId('ai-advisor-modal')).not.toBeVisible();
    
    // Configuration form should be updated with recommendation
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe(recommendation.blueprint);
    
    // Preview should update with recommended configuration
    await page.waitForPreviewUpdate();
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain(recommendation.blueprint);
  });

  test('should show multiple alternative recommendations', async () => {
    await page.openAIAdvisor();
    
    await page.fillAdvisorRequirements({
      domain: 'healthcare',
      teamSize: 6,
      experience: 'senior',
      timeline: 'standard'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Should show alternatives
    await expect(page.page.getByTestId('alternative-recommendations')).toBeVisible();
    
    const alternatives = page.page.getByTestId('alternative-option');
    const alternativeCount = await alternatives.count();
    expect(alternativeCount).toBeGreaterThan(0);
    expect(alternativeCount).toBeLessThanOrEqual(3);
    
    // Each alternative should have confidence score
    for (let i = 0; i < alternativeCount; i++) {
      const alternative = alternatives.nth(i);
      await expect(alternative.getByTestId('alternative-confidence')).toBeVisible();
    }
  });

  test('should handle advisor with incomplete requirements', async () => {
    await page.openAIAdvisor();
    
    // Fill only partial requirements
    await page.fillAdvisorRequirements({
      domain: 'iot'
      // Missing team size, experience, timeline
    });
    
    // Should still provide recommendation with defaults
    const recommendation = await page.getAIRecommendation();
    
    expect(recommendation.blueprint).toBeTruthy();
    expect(recommendation.confidence).toBeGreaterThan(40); // Lower confidence for incomplete info
  });

  test('should show reasoning for recommendations', async () => {
    await page.openAIAdvisor();
    
    await page.fillAdvisorRequirements({
      domain: 'fintech',
      teamSize: 5,
      experience: 'expert',
      timeline: 'thorough'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Should show detailed reasoning
    await expect(page.page.getByTestId('recommendation-reasoning')).toBeVisible();
    
    const reasoningItems = page.page.getByTestId('reasoning-item');
    const reasoningCount = await reasoningItems.count();
    expect(reasoningCount).toBeGreaterThan(0);
    
    // Reasoning should be specific to the requirements
    const allReasoningText = await reasoningItems.allTextContents();
    const fullReasoningText = allReasoningText.join(' ').toLowerCase();
    
    expect(fullReasoningText).toContain('fintech');
    expect(fullReasoningText).toContain('expert');
  });

  test('should allow comparing different scenarios', async () => {
    await page.openAIAdvisor();
    
    // Get recommendation for scenario 1
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 2,
      experience: 'junior',
      timeline: 'mvp'
    });
    
    const recommendation1 = await page.getAIRecommendation();
    
    // Change to scenario 2
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 8,
      experience: 'expert',
      timeline: 'thorough'
    });
    
    const recommendation2 = await page.getAIRecommendation();
    
    // Recommendations should be different for different scenarios
    expect(recommendation1.blueprint).not.toBe(recommendation2.blueprint);
    expect(recommendation2.confidence).toBeGreaterThan(recommendation1.confidence);
  });

  test('should validate advisor form inputs', async () => {
    await page.openAIAdvisor();
    
    // Try to get recommendation without required fields
    await page.page.getByRole('button', { name: /get recommendation/i }).click();
    
    // Should show validation errors
    const errors = await page.getValidationErrors();
    expect(errors.length).toBeGreaterThan(0);
    
    // Should not proceed without valid inputs
    await expect(page.page.getByTestId('ai-recommendation')).not.toBeVisible();
  });

  test('should close advisor modal without applying changes', async () => {
    // Set initial configuration
    await page.fillProjectName('initial-project');
    await page.selectProjectType('cli');
    
    await page.openAIAdvisor();
    
    // Get a different recommendation
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 5,
      experience: 'senior',
      timeline: 'standard'
    });
    
    await page.getAIRecommendation();
    
    // Close modal without applying
    await page.page.getByRole('button', { name: /close|cancel/i }).click();
    
    // Original configuration should be preserved
    expect(await page.page.getByLabel(/project name/i).inputValue()).toBe('initial-project');
    expect(await page.page.getByLabel(/project type/i).inputValue()).toBe('cli');
  });

  test('should integrate advisor with real-time preview', async () => {
    await page.waitForWebSocketConnection();
    await page.openAIAdvisor();
    
    await page.fillAdvisorRequirements({
      domain: 'microservice',
      teamSize: 6,
      experience: 'senior',
      timeline: 'standard'
    });
    
    const recommendation = await page.getAIRecommendation();
    
    // Apply recommendation
    await page.page.getByRole('button', { name: /apply recommendation/i }).click();
    
    // Preview should update in real-time
    await page.waitForRealTimeUpdate();
    
    const previewContent = await page.getPreviewContent();
    expect(previewContent).toContain(recommendation.blueprint);
    
    // File explorer should update with new structure
    const fileList = await page.getFileList();
    expect(fileList.length).toBeGreaterThan(0);
  });

  test('should show confidence levels visually', async () => {
    await page.openAIAdvisor();
    
    await page.fillAdvisorRequirements({
      domain: 'e-commerce',
      teamSize: 5,
      experience: 'expert',
      timeline: 'standard'
    });
    
    await page.getAIRecommendation();
    
    // Confidence should be displayed visually
    await expect(page.page.getByTestId('confidence-bar')).toBeVisible();
    await expect(page.page.getByTestId('confidence-percentage')).toBeVisible();
    
    // High confidence should be indicated (green color, high percentage)
    const confidenceText = await page.page.getByTestId('confidence-percentage').textContent();
    const confidence = parseInt(confidenceText?.replace('%', '') || '0');
    expect(confidence).toBeGreaterThan(70);
    
    // Visual indicator should reflect confidence level
    if (confidence > 80) {
      await expect(page.page.getByTestId('confidence-bar')).toHaveClass(/green|high/);
    } else if (confidence > 60) {
      await expect(page.page.getByTestId('confidence-bar')).toHaveClass(/yellow|medium/);
    }
  });
});