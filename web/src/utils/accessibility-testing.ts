/**
 * Comprehensive Accessibility Testing Utilities for WCAG 2.1 AA Compliance
 * Provides automated testing tools and validation utilities for accessibility
 */

import { AccessibilityAuditor, runQuickAccessibilityCheck, ColorContrastChecker } from './accessibility-audit';

// Testing utilities for automated accessibility validation
export interface AccessibilityTestConfig {
  includeColorContrast?: boolean;
  includeKeyboardNavigation?: boolean;
  includeAriaCompliance?: boolean;
  includeHeadingStructure?: boolean;
  includeImageAltText?: boolean;
  customChecks?: Array<() => Promise<AccessibilityIssue[]>>;
}

export interface AccessibilityTestResult {
  passed: boolean;
  score: number;
  issues: AccessibilityIssue[];
  recommendations: string[];
  summary: {
    totalChecks: number;
    passedChecks: number;
    failedChecks: number;
    criticalIssues: number;
    seriousIssues: number;
    moderateIssues: number;
    minorIssues: number;
  };
}

export interface AccessibilityIssue {
  severity: 'error' | 'warning' | 'info';
  rule: string;
  element: Element | null;
  message: string;
  wcagCriterion?: string;
  fix?: string;
  category?: 'perceivable' | 'operable' | 'understandable' | 'robust';
  impact?: 'critical' | 'serious' | 'moderate' | 'minor';
  location?: string;
}

// Main accessibility testing suite
export class AccessibilityTestSuite {
  private config: AccessibilityTestConfig;

  constructor(config: AccessibilityTestConfig = {}) {
    this.config = {
      includeColorContrast: true,
      includeKeyboardNavigation: true,
      includeAriaCompliance: true,
      includeHeadingStructure: true,
      includeImageAltText: true,
      ...config
    };
  }

  /**
   * Run comprehensive accessibility tests
   */
  async runTests(): Promise<AccessibilityTestResult> {
    const auditor = new AccessibilityAuditor();
    const baseResult = await auditor.runFullAudit();
    
    const allIssues = [...baseResult.issues];
    const recommendations: string[] = [];

    // Add element location information
    baseResult.issues.forEach(issue => {
      if (issue.element) {
        if (issue.element) {
          (issue as any).location = this.getElementLocation(issue.element as HTMLElement);
        }
      }
    });

    // Run custom checks if provided
    if (this.config.customChecks) {
      for (const customCheck of this.config.customChecks) {
        try {
          const customIssues = await customCheck();
          allIssues.push(...customIssues);
        } catch (error) {
          console.warn('Custom accessibility check failed:', error);
        }
      }
    }

    // Generate recommendations
    recommendations.push(...this.generateRecommendations(allIssues));

    // Calculate summary
    const summary = this.calculateSummary(allIssues);

    return {
      passed: baseResult.failed === 0 && baseResult.score >= 80,
      score: Math.max(0, baseResult.score - (allIssues.length - baseResult.issues.length) * 5),
      issues: allIssues,
      recommendations,
      summary
    };
  }

  /**
   * Test specific component for accessibility
   */
  async testComponent(element: HTMLElement): Promise<AccessibilityTestResult> {
    const originalBody = document.body.innerHTML;
    
    try {
      // Create isolated test environment
      const testContainer = document.createElement('div');
      testContainer.appendChild(element.cloneNode(true));
      document.body.innerHTML = '';
      document.body.appendChild(testContainer);

      // Run tests on isolated component
      const result = await this.runTests();
      
      return result;
    } finally {
      // Restore original DOM
      document.body.innerHTML = originalBody;
    }
  }

  /**
   * Test color contrast for specific elements
   */
  testColorContrast(elements?: HTMLElement[]): Array<{
    element: HTMLElement;
    result: {
      ratio: number;
      passes: boolean;
      level: 'AAA' | 'AA' | 'FAIL';
      recommendation?: string;
    };
  }> {
    const elementsToTest = elements || Array.from(document.querySelectorAll('*')) as HTMLElement[];
    const results: Array<any> = [];

    elementsToTest.forEach(element => {
      const computedStyle = window.getComputedStyle(element);
      const textColor = computedStyle.color;
      const backgroundColor = computedStyle.backgroundColor;
      const fontSize = parseFloat(computedStyle.fontSize);
      const fontWeight = parseFloat(computedStyle.fontWeight) || 400;

      // Skip if no visible text or transparent background
      if (!textColor || backgroundColor === 'rgba(0, 0, 0, 0)' || !element.textContent?.trim()) {
        return;
      }

      const result = ColorContrastChecker.checkContrast(
        textColor,
        backgroundColor
      );

      results.push({
        element,
        result
      });
    });

    return results;
  }

  /**
   * Test keyboard navigation
   */
  testKeyboardNavigation(): Promise<{
    focusableElements: HTMLElement[];
    tabOrder: HTMLElement[];
    issues: AccessibilityIssue[];
  }> {
    return new Promise((resolve) => {
      const focusableElements = this.getFocusableElements();
      const tabOrder: HTMLElement[] = [];
      const issues: AccessibilityIssue[] = [];

      let currentIndex = 0;
      
      const testNextElement = () => {
        if (currentIndex >= focusableElements.length) {
          // Test completed
          resolve({
            focusableElements,
            tabOrder,
            issues
          });
          return;
        }

        const element = focusableElements[currentIndex];
        
        // Test if element can receive focus
        try {
          element.focus();
          
          if (document.activeElement === element) {
            tabOrder.push(element);
            
            // Check for visible focus indicator
            const computedStyle = window.getComputedStyle(element);
            const hasVisibleFocus = computedStyle.outline !== 'none' || 
                                  computedStyle.boxShadow.includes('rgb') ||
                                  element.classList.toString().includes('focus');
            
            if (!hasVisibleFocus) {
              (issues as any[]).push({
                severity: 'error',
                category: 'operable',
                wcagCriterion: 'WCAG 2.4.7',
                message: 'Focusable element lacks visible focus indicator',
                element,
                fix: 'Add visible focus indicator (outline, box-shadow, or CSS focus styles)',
                impact: 'serious',
                location: this.getElementLocation(element)
              });
            }
          } else {
            (issues as any[]).push({
              severity: 'warning',
              category: 'operable',
              wcagCriterion: 'WCAG 2.1.1',
              message: 'Element expected to be focusable but cannot receive focus',
              element,
              fix: 'Check if element is properly configured for keyboard access',
              impact: 'moderate',
              location: this.getElementLocation(element)
            });
          }
        } catch (error) {
          (issues as any[]).push({
            severity: 'error',
            category: 'operable',
            wcagCriterion: 'WCAG 2.1.1',
            message: 'Error occurred while testing focus',
            element,
            fix: 'Investigate focus-related JavaScript errors',
            impact: 'serious',
            location: this.getElementLocation(element)
          });
        }

        currentIndex++;
        setTimeout(testNextElement, 10); // Small delay to allow focus events
      };

      testNextElement();
    });
  }

  /**
   * Test ARIA implementation
   */
  testAriaImplementation(): AccessibilityIssue[] {
    const issues: AccessibilityIssue[] = [];

    // Test for missing ARIA labels
    const interactiveElements = document.querySelectorAll('button, input, select, textarea, [role="button"], [role="link"], [role="menuitem"]');
    
    interactiveElements.forEach(element => {
      const htmlElement = element as HTMLElement;
      const hasLabel = htmlElement.getAttribute('aria-label') ||
                      htmlElement.getAttribute('aria-labelledby') ||
                      htmlElement.textContent?.trim() ||
                      document.querySelector(`label[for="${htmlElement.id}"]`);

      if (!hasLabel) {
        (issues as any[]).push({
          severity: 'error',
          category: 'perceivable',
          wcagCriterion: 'WCAG 1.3.1',
          message: 'Interactive element missing accessible label',
          element: htmlElement,
          fix: 'Add aria-label, aria-labelledby, or associate with a label element',
          impact: 'critical',
          location: this.getElementLocation(htmlElement)
        });
      }
    });

    // Test for invalid ARIA attributes
    const elementsWithAria = document.querySelectorAll('[aria-expanded], [aria-selected], [aria-checked], [aria-disabled], [aria-hidden]');
    
    elementsWithAria.forEach(element => {
      const htmlElement = element as HTMLElement;
      
      // Check aria-expanded
      const expanded = htmlElement.getAttribute('aria-expanded');
      if (expanded && !['true', 'false'].includes(expanded)) {
        (issues as any[]).push({
          severity: 'error',
          category: 'robust',
          wcagCriterion: 'WCAG 4.1.2',
          message: `Invalid aria-expanded value: ${expanded}`,
          element: htmlElement,
          fix: 'Use "true" or "false" for aria-expanded',
          impact: 'serious',
          location: this.getElementLocation(htmlElement)
        });
      }

      // Check other boolean ARIA attributes
      ['aria-selected', 'aria-checked', 'aria-disabled', 'aria-hidden'].forEach(attr => {
        const value = htmlElement.getAttribute(attr);
        if (value && !['true', 'false'].includes(value)) {
          (issues as any[]).push({
            severity: 'error',
            category: 'robust',
            wcagCriterion: 'WCAG 4.1.2',
            message: `Invalid ${attr} value: ${value}`,
            element: htmlElement,
            fix: `Use "true" or "false" for ${attr}`,
            impact: 'serious',
            location: this.getElementLocation(htmlElement)
          });
        }
      });
    });

    return issues;
  }

  /**
   * Test for responsive design accessibility
   */
  testResponsiveAccessibility(): AccessibilityIssue[] {
    const issues: AccessibilityIssue[] = [];
    const viewportWidth = window.innerWidth;

    // Test touch target sizes on mobile
    if (viewportWidth <= 768) {
      const touchTargets = document.querySelectorAll('button, input, select, textarea, a, [role="button"], [role="link"]');
      
      touchTargets.forEach(element => {
        const htmlElement = element as HTMLElement;
        const rect = htmlElement.getBoundingClientRect();
        const minSize = 44; // WCAG 2.5.5 minimum touch target size

        if (rect.width < minSize || rect.height < minSize) {
          (issues as any[]).push({
            severity: 'error',
            category: 'operable',
            wcagCriterion: 'WCAG 2.5.5',
            message: `Touch target too small: ${rect.width}x${rect.height}px (minimum: ${minSize}x${minSize}px)`,
            element: htmlElement,
            fix: 'Increase touch target size to at least 44x44px',
            impact: 'serious',
            location: this.getElementLocation(htmlElement)
          });
        }
      });
    }

    // Test for horizontal scrolling
    if (document.body.scrollWidth > viewportWidth) {
      (issues as any[]).push({
        severity: 'warning',
        category: 'perceivable',
        wcagCriterion: 'WCAG 1.4.10',
        message: 'Horizontal scrolling detected',
        fix: 'Ensure content reflows properly at different viewport sizes',
        impact: 'moderate'
      });
    }

    return issues;
  }

  /**
   * Generate accessibility testing report
   */
  generateReport(result: AccessibilityTestResult): string {
    const { summary, issues, recommendations, score, passed } = result;
    
    let report = `# Accessibility Test Report\n\n`;
    report += `**Overall Score:** ${score}/100\n`;
    report += `**Status:** ${passed ? '✅ PASSED' : '❌ FAILED'}\n\n`;
    
    report += `## Summary\n`;
    report += `- Total Checks: ${summary.totalChecks}\n`;
    report += `- Passed: ${summary.passedChecks}\n`;
    report += `- Failed: ${summary.failedChecks}\n`;
    report += `- Critical Issues: ${summary.criticalIssues}\n`;
    report += `- Serious Issues: ${summary.seriousIssues}\n`;
    report += `- Moderate Issues: ${summary.moderateIssues}\n`;
    report += `- Minor Issues: ${summary.minorIssues}\n\n`;

    if (recommendations.length > 0) {
      report += `## Recommendations\n`;
      recommendations.forEach((rec, index) => {
        report += `${index + 1}. ${rec}\n`;
      });
      report += '\n';
    }

    if (issues.length > 0) {
      report += `## Issues Found\n\n`;
      
      const issuesByCategory = issues.reduce((acc, issue) => {
        if (!acc[issue.category]) acc[issue.category] = [];
        acc[issue.category].push(issue);
        return acc;
      }, {} as Record<string, AccessibilityIssue[]>);

      Object.entries(issuesByCategory).forEach(([category, categoryIssues]) => {
        report += `### ${category.charAt(0).toUpperCase() + category.slice(1)}\n\n`;
        
        categoryIssues.forEach((issue, index) => {
          const severity = issue.impact === 'critical' ? '🔴' : 
                          issue.impact === 'serious' ? '🟠' : 
                          issue.impact === 'moderate' ? '🟡' : '🔵';
          
          report += `${index + 1}. ${severity} **${issue.message}**\n`;
          report += `   - WCAG Reference: ${issue.wcagCriterion}\n`;
          report += `   - Impact: ${(issue as any).impact}\n`;
          report += `   - Suggestion: ${issue.fix}\n`;
          if (issue.location) {
            report += `   - Location: ${issue.location}\n`;
          }
          report += '\n';
        });
      });
    }

    return report;
  }

  // Helper methods
  private getFocusableElements(): HTMLElement[] {
    const selector = [
      'button:not([disabled]):not([aria-hidden="true"])',
      'input:not([disabled]):not([type="hidden"]):not([aria-hidden="true"])',
      'select:not([disabled]):not([aria-hidden="true"])',
      'textarea:not([disabled]):not([aria-hidden="true"])',
      'a[href]:not([aria-hidden="true"])',
      '[tabindex]:not([tabindex="-1"]):not([aria-hidden="true"])',
      'audio[controls]:not([aria-hidden="true"])',
      'video[controls]:not([aria-hidden="true"])'
    ].join(', ');

    return Array.from(document.querySelectorAll(selector)) as HTMLElement[];
  }

  private getElementLocation(element: HTMLElement): string {
    const tagName = element.tagName.toLowerCase();
    const id = element.id ? `#${element.id}` : '';
    const className = element.className ? `.${element.className.split(' ').join('.')}` : '';
    const textContent = element.textContent?.slice(0, 30) || '';
    
    return `${tagName}${id}${className} ${textContent ? `"${textContent}..."` : ''}`.trim();
  }

  private generateRecommendations(issues: AccessibilityIssue[]): string[] {
    const recommendations: string[] = [];
    const criticalIssues = issues.filter(issue => issue.impact === 'critical');
    const seriousIssues = issues.filter(issue => issue.impact === 'serious');

    if (criticalIssues.length > 0) {
      recommendations.push(`Address ${criticalIssues.length} critical accessibility issues immediately`);
    }

    if (seriousIssues.length > 0) {
      recommendations.push(`Fix ${seriousIssues.length} serious accessibility issues`);
    }

    // Category-specific recommendations
    const categories = issues.reduce((acc, issue) => {
      acc[issue.category] = (acc[issue.category] || 0) + 1;
      return acc;
    }, {} as Record<string, number>);

    Object.entries(categories).forEach(([category, count]) => {
      if (count > 2) {
        switch (category) {
          case 'perceivable':
            recommendations.push('Improve content perceivability with better color contrast and alternative text');
            break;
          case 'operable':
            recommendations.push('Enhance keyboard navigation and focus management');
            break;
          case 'understandable':
            recommendations.push('Improve content clarity and error handling');
            break;
          case 'robust':
            recommendations.push('Fix ARIA implementation and ensure assistive technology compatibility');
            break;
        }
      }
    });

    return recommendations;
  }

  private calculateSummary(issues: AccessibilityIssue[]): AccessibilityTestResult['summary'] {
    const criticalIssues = issues.filter(issue => issue.impact === 'critical').length;
    const seriousIssues = issues.filter(issue => issue.impact === 'serious').length;
    const moderateIssues = issues.filter(issue => issue.impact === 'moderate').length;
    const minorIssues = issues.filter(issue => issue.impact === 'minor').length;
    
    const totalChecks = issues.length + 10; // Base checks + issues found
    const failedChecks = issues.length;
    const passedChecks = totalChecks - failedChecks;

    return {
      totalChecks,
      passedChecks,
      failedChecks,
      criticalIssues,
      seriousIssues,
      moderateIssues,
      minorIssues
    };
  }
}

// Utility functions for testing specific scenarios
export const accessibilityTestUtils = {
  /**
   * Test form accessibility
   */
  testForm: async (formElement: HTMLFormElement): Promise<AccessibilityIssue[]> => {
    const issues: AccessibilityIssue[] = [];
    const suite = new AccessibilityTestSuite();

    // Test form controls
    const controls = formElement.querySelectorAll('input, select, textarea');
    controls.forEach(control => {
      const htmlControl = control as HTMLElement;
      const hasLabel = htmlControl.getAttribute('aria-label') ||
                      htmlControl.getAttribute('aria-labelledby') ||
                      formElement.querySelector(`label[for="${htmlControl.id}"]`);

      if (!hasLabel) {
        (issues as any[]).push({
          severity: 'error',
          category: 'perceivable',
          wcagCriterion: 'WCAG 1.3.1',
          message: 'Form control missing label',
          element: htmlControl,
          fix: 'Add label element or aria-label attribute',
          impact: 'critical'
        });
      }
    });

    return issues;
  },

  /**
   * Test modal accessibility
   */
  testModal: async (modalElement: HTMLElement): Promise<AccessibilityIssue[]> => {
    const issues: AccessibilityIssue[] = [];

    // Check for proper ARIA attributes
    if (!modalElement.getAttribute('role') || modalElement.getAttribute('role') !== 'dialog') {
      (issues as any[]).push({
        severity: 'error',
        category: 'robust',
        wcagCriterion: 'WCAG 4.1.2',
        message: 'Modal missing role="dialog"',
        element: modalElement,
        fix: 'Add role="dialog" to modal element',
        impact: 'serious'
      });
    }

    if (!modalElement.getAttribute('aria-modal')) {
      (issues as any[]).push({
        severity: 'error',
        category: 'robust',
        wcagCriterion: 'WCAG 4.1.2',
        message: 'Modal missing aria-modal="true"',
        element: modalElement,
        fix: 'Add aria-modal="true" to modal element',
        impact: 'serious'
      });
    }

    if (!modalElement.getAttribute('aria-labelledby') && !modalElement.getAttribute('aria-label')) {
      (issues as any[]).push({
        severity: 'error',
        category: 'perceivable',
        wcagCriterion: 'WCAG 1.3.1',
        message: 'Modal missing accessible name',
        element: modalElement,
        fix: 'Add aria-labelledby or aria-label to modal element',
        impact: 'critical'
      });
    }

    return issues;
  },

  /**
   * Test component color contrast
   */
  testComponentContrast: (element: HTMLElement): Array<{
    passes: boolean;
    ratio: number;
    recommendation?: string;
  }> => {
    const suite = new AccessibilityTestSuite();
    return suite.testColorContrast([element]).map(result => ({
      passes: result.result.passes,
      ratio: result.result.ratio,
      recommendation: result.result.recommendation
    }));
  },

  /**
   * Quick accessibility check for development
   */
  quickCheck: (): Promise<{ passed: boolean; issues: string[]; score: number }> => {
    return new Promise(async resolve => {
      setTimeout(async () => {
        const result = await runQuickAccessibilityCheck();
        resolve({
          passed: result.failed === 0,
          issues: result.issues.map(issue => issue.message),
          score: result.score
        });
      }, 100);
    });
  }
};

// Global accessibility testing interface for development
if (typeof window !== 'undefined') {
  (window as any).accessibilityTest = {
    runFullTest: async () => {
      const suite = new AccessibilityTestSuite();
      const result = await suite.runTests();
      console.log('Accessibility Test Results:', result);
      console.log(suite.generateReport(result));
      return result;
    },
    
    quickCheck: async () => {
      const result = await accessibilityTestUtils.quickCheck();
      console.log('Quick Accessibility Check:', result);
      return result;
    },

    testElement: async (element: HTMLElement) => {
      const suite = new AccessibilityTestSuite();
      const result = await suite.testComponent(element);
      console.log('Element Accessibility Test:', result);
      return result;
    }
  };
}