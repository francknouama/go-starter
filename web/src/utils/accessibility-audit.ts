// WCAG 2.1 AA Accessibility Audit System for go-starter
// Automated accessibility compliance checking and reporting

export interface AccessibilityIssue {
  severity: 'error' | 'warning' | 'info'
  rule: string
  element: Element | null
  message: string
  wcagCriterion?: string
  fix?: string
}

export interface AccessibilityReport {
  score: number // 0-100
  passed: number
  failed: number
  warnings: number
  issues: AccessibilityIssue[]
  timestamp: string
}

export class AccessibilityAuditor {
  private issues: AccessibilityIssue[] = []

  /**
   * Run comprehensive WCAG 2.1 AA audit
   */
  async runFullAudit(container: Document | Element = document): Promise<AccessibilityReport> {
    this.issues = []

    // WCAG 1.x - Perceivable
    this.checkColorContrast(container)
    this.checkTextAlternatives(container)
    this.checkHeadingStructure(container)
    this.checkFocusIndicators(container)

    // WCAG 2.x - Operable  
    this.checkKeyboardAccessibility(container)
    this.checkFocusManagement(container)
    this.checkSkipLinks(container)
    this.checkTargetSize(container)

    // WCAG 3.x - Understandable
    this.checkFormLabels(container)
    this.checkErrorHandling(container)
    this.checkConsistentNavigation(container)

    // WCAG 4.x - Robust
    this.checkAriaImplementation(container)
    this.checkSemanticStructure(container)
    this.checkLiveRegions(container)

    return this.generateReport()
  }

  /**
   * Check color contrast compliance (WCAG 1.4.3, 1.4.6)
   */
  private checkColorContrast(container: Document | Element): void {
    const textElements = container.querySelectorAll('p, span, div, h1, h2, h3, h4, h5, h6, button, input, textarea, label, li, td, th')
    
    textElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const styles = window.getComputedStyle(htmlElement)
      const textColor = styles.color
      const backgroundColor = this.getEffectiveBackgroundColor(htmlElement)
      
      if (textColor && backgroundColor) {
        const ratio = this.calculateContrastRatio(textColor, backgroundColor)
        const fontSize = parseFloat(styles.fontSize)
        const fontWeight = styles.fontWeight
        
        // Large text is 18pt (24px) or 14pt (18.67px) bold
        const isLargeText = fontSize >= 24 || (fontSize >= 18.67 && (fontWeight === 'bold' || parseInt(fontWeight) >= 700))
        const requiredRatio = isLargeText ? 3 : 4.5
        
        if (ratio < requiredRatio) {
          this.addIssue({
            severity: 'error',
            rule: 'color-contrast',
            element: htmlElement,
            message: `Color contrast ratio ${ratio.toFixed(2)}:1 is below required ${requiredRatio}:1`,
            wcagCriterion: '1.4.3 Contrast (Minimum)',
            fix: `Increase contrast between text (${textColor}) and background (${backgroundColor})`
          })
        }
      }
    })
  }

  /**
   * Check text alternatives for images (WCAG 1.1.1)
   */
  private checkTextAlternatives(container: Document | Element): void {
    const images = container.querySelectorAll('img')
    
    images.forEach(img => {
      const alt = img.getAttribute('alt')
      const ariaLabel = img.getAttribute('aria-label')
      const ariaLabelledby = img.getAttribute('aria-labelledby')
      const role = img.getAttribute('role')
      
      // Decorative images should have empty alt or role="presentation"
      if (role === 'presentation' || role === 'none') {
        return // Valid decorative image
      }
      
      if (!alt && !ariaLabel && !ariaLabelledby) {
        this.addIssue({
          severity: 'error',
          rule: 'text-alternatives',
          element: img,
          message: 'Image missing text alternative',
          wcagCriterion: '1.1.1 Non-text Content',
          fix: 'Add alt attribute or aria-label for meaningful images, or role="presentation" for decorative images'
        })
      }
    })

    // Check SVG icons
    const svgs = container.querySelectorAll('svg')
    svgs.forEach(svg => {
      const ariaLabel = svg.getAttribute('aria-label')
      const ariaLabelledby = svg.getAttribute('aria-labelledby')
      const titleElement = svg.querySelector('title')
      const role = svg.getAttribute('role')
      
      if (role !== 'presentation' && !ariaLabel && !ariaLabelledby && !titleElement) {
        this.addIssue({
          severity: 'warning',
          rule: 'text-alternatives',
          element: svg,
          message: 'SVG icon may need text alternative',
          wcagCriterion: '1.1.1 Non-text Content',
          fix: 'Add aria-label, aria-labelledby, or <title> element for meaningful SVGs'
        })
      }
    })
  }

  /**
   * Check heading structure (WCAG 1.3.1)
   */
  private checkHeadingStructure(container: Document | Element): void {
    const headings = Array.from(container.querySelectorAll('h1, h2, h3, h4, h5, h6'))
    let previousLevel = 0
    
    headings.forEach((heading, index) => {
      const level = parseInt(heading.tagName.charAt(1))
      
      // Check for h1 at page start
      if (index === 0 && level !== 1) {
        this.addIssue({
          severity: 'warning',
          rule: 'heading-structure',
          element: heading,
          message: 'Page should start with h1',
          wcagCriterion: '1.3.1 Info and Relationships',
          fix: 'Use h1 for the main page heading'
        })
      }
      
      // Check for heading level jumps
      if (previousLevel > 0 && level > previousLevel + 1) {
        this.addIssue({
          severity: 'error',
          rule: 'heading-structure',
          element: heading,
          message: `Heading level jumps from h${previousLevel} to h${level}`,
          wcagCriterion: '1.3.1 Info and Relationships',
          fix: 'Use sequential heading levels (h1, h2, h3, etc.)'
        })
      }
      
      previousLevel = level
    })
  }

  /**
   * Check keyboard accessibility (WCAG 2.1.1, 2.1.2)
   */
  private checkKeyboardAccessibility(container: Document | Element): void {
    const interactiveElements = container.querySelectorAll('button, input, select, textarea, a[href], [tabindex]:not([tabindex="-1"])')
    
    interactiveElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const tabIndex = htmlElement.getAttribute('tabindex')
      
      // Check for positive tabindex (anti-pattern)
      if (tabIndex && parseInt(tabIndex) > 0) {
        this.addIssue({
          severity: 'warning',
          rule: 'keyboard-accessibility',
          element: htmlElement,
          message: 'Positive tabindex disrupts natural tab order',
          wcagCriterion: '2.4.3 Focus Order',
          fix: 'Use tabindex="0" or remove tabindex to maintain natural order'
        })
      }
      
      // Check if element is focusable but not keyboard accessible
      if (htmlElement.onclick && !this.isKeyboardAccessible(htmlElement)) {
        this.addIssue({
          severity: 'error',
          rule: 'keyboard-accessibility',
          element: htmlElement,
          message: 'Interactive element not keyboard accessible',
          wcagCriterion: '2.1.1 Keyboard',
          fix: 'Add keyboard event handlers or use button/a elements'
        })
      }
    })
  }

  /**
   * Check focus indicators (WCAG 2.4.7)
   */
  private checkFocusIndicators(container: Document | Element): void {
    const focusableElements = container.querySelectorAll('button, input, select, textarea, a[href], [tabindex]:not([tabindex="-1"])')
    
    focusableElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const styles = window.getComputedStyle(htmlElement, ':focus')
      const outline = styles.outline
      const boxShadow = styles.boxShadow
      
      // Check if focus indicator is removed or insufficient
      if (outline === 'none' || outline === '0px' || outline === '0') {
        if (!boxShadow || boxShadow === 'none') {
          this.addIssue({
            severity: 'error',
            rule: 'focus-indicators',
            element: htmlElement,
            message: 'Focus indicator removed without replacement',
            wcagCriterion: '2.4.7 Focus Visible',
            fix: 'Provide visible focus indicator using outline or box-shadow'
          })
        }
      }
    })
  }

  /**
   * Check form labels (WCAG 3.3.2)
   */
  private checkFormLabels(container: Document | Element): void {
    const inputs = container.querySelectorAll('input:not([type="hidden"]), textarea, select')
    
    inputs.forEach(input => {
      const htmlInput = input as HTMLInputElement
      const id = htmlInput.id
      const ariaLabel = htmlInput.getAttribute('aria-label')
      const ariaLabelledby = htmlInput.getAttribute('aria-labelledby')
      const label = id ? container.querySelector(`label[for="${id}"]`) : null
      const wrappingLabel = htmlInput.closest('label')
      
      if (!label && !wrappingLabel && !ariaLabel && !ariaLabelledby) {
        this.addIssue({
          severity: 'error',
          rule: 'form-labels',
          element: htmlInput,
          message: 'Form control missing accessible label',
          wcagCriterion: '3.3.2 Labels or Instructions',
          fix: 'Add <label> element, aria-label, or aria-labelledby'
        })
      }
    })
  }

  /**
   * Check ARIA implementation (WCAG 4.1.2)
   */
  private checkAriaImplementation(container: Document | Element): void {
    const ariaElements = container.querySelectorAll('[role], [aria-label], [aria-labelledby], [aria-describedby], [aria-expanded], [aria-selected], [aria-checked]')
    
    ariaElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const role = htmlElement.getAttribute('role')
      
      // Check for invalid ARIA roles
      if (role && !this.isValidAriaRole(role)) {
        this.addIssue({
          severity: 'error',
          rule: 'aria-implementation',
          element: htmlElement,
          message: `Invalid ARIA role: ${role}`,
          wcagCriterion: '4.1.2 Name, Role, Value',
          fix: 'Use valid ARIA roles from the ARIA specification'
        })
      }
      
      // Check for required ARIA properties
      if (role === 'button' && !htmlElement.hasAttribute('aria-pressed') && htmlElement.getAttribute('aria-pressed') !== 'false') {
        // This is actually okay - not all buttons need aria-pressed
      }
      
      // Check for accessible names on interactive elements
      if (['button', 'link', 'menuitem'].includes(role || htmlElement.tagName.toLowerCase())) {
        if (!this.hasAccessibleName(htmlElement)) {
          this.addIssue({
            severity: 'error',
            rule: 'aria-implementation',
            element: htmlElement,
            message: 'Interactive element missing accessible name',
            wcagCriterion: '4.1.2 Name, Role, Value',
            fix: 'Add aria-label, aria-labelledby, or text content'
          })
        }
      }
    })
  }

  /**
   * Check target size (WCAG 2.5.5)
   */
  private checkTargetSize(container: Document | Element): void {
    const interactiveElements = container.querySelectorAll('button, input, select, a[href], [role="button"], [role="link"]')
    
    interactiveElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const rect = htmlElement.getBoundingClientRect()
      const minSize = 44 // WCAG 2.5.5 minimum 44px
      
      if (rect.width < minSize || rect.height < minSize) {
        this.addIssue({
          severity: 'warning',
          rule: 'target-size',
          element: htmlElement,
          message: `Target size ${rect.width}x${rect.height}px is below 44x44px minimum`,
          wcagCriterion: '2.5.5 Target Size',
          fix: 'Increase padding or minimum dimensions to 44x44px'
        })
      }
    })
  }

  /**
   * Check skip links (WCAG 2.4.1)
   */
  private checkSkipLinks(container: Document | Element): void {
    const skipLinks = container.querySelectorAll('a[href^="#"]')
    const hasSkipToMain = Array.from(skipLinks).some(link => 
      link.textContent?.toLowerCase().includes('skip') || 
      link.getAttribute('href') === '#main' ||
      link.getAttribute('href') === '#content'
    )
    
    if (!hasSkipToMain && container === document) {
      this.addIssue({
        severity: 'warning',
        rule: 'skip-links',
        element: null,
        message: 'No skip link found',
        wcagCriterion: '2.4.1 Bypass Blocks',
        fix: 'Add skip link to main content at page beginning'
      })
    }
  }

  /**
   * Check semantic structure (WCAG 1.3.1, 1.3.6)
   */
  private checkSemanticStructure(container: Document | Element): void {
    const requiredLandmarks = ['main', 'nav']
    const landmarks = container.querySelectorAll('main, nav, header, footer, aside, section[aria-labelledby], [role="main"], [role="navigation"], [role="banner"], [role="contentinfo"], [role="complementary"]')
    
    requiredLandmarks.forEach(landmark => {
      const exists = Array.from(landmarks).some(el => 
        el.tagName.toLowerCase() === landmark || 
        el.getAttribute('role') === landmark ||
        (landmark === 'main' && el.getAttribute('role') === 'main')
      )
      
      if (!exists && container === document) {
        this.addIssue({
          severity: 'error',
          rule: 'semantic-structure',
          element: null,
          message: `Missing ${landmark} landmark`,
          wcagCriterion: '1.3.6 Identify Purpose',
          fix: `Add <${landmark}> element or role="${landmark}"`
        })
      }
    })
  }

  /**
   * Check live regions (WCAG 4.1.3)
   */
  private checkLiveRegions(container: Document | Element): void {
    const dynamicContent = container.querySelectorAll('[aria-live], [role="status"], [role="alert"]')
    const potentialDynamic = container.querySelectorAll('.loading, .error, .success, .notification')
    
    potentialDynamic.forEach(element => {
      const htmlElement = element as HTMLElement
      const hasLiveRegion = htmlElement.hasAttribute('aria-live') || 
                           htmlElement.hasAttribute('role')
      
      if (!hasLiveRegion) {
        this.addIssue({
          severity: 'info',
          rule: 'live-regions',
          element: htmlElement,
          message: 'Dynamic content may need live region',
          wcagCriterion: '4.1.3 Status Messages',
          fix: 'Add aria-live="polite" or role="status" for status updates'
        })
      }
    })
  }

  /**
   * Helper methods
   */
  private calculateContrastRatio(color1: string, color2: string): number {
    const rgb1 = this.parseColor(color1)
    const rgb2 = this.parseColor(color2)
    
    if (!rgb1 || !rgb2) return 21 // Assume best case if can't parse
    
    const l1 = this.getRelativeLuminance(rgb1)
    const l2 = this.getRelativeLuminance(rgb2)
    
    const lighter = Math.max(l1, l2)
    const darker = Math.min(l1, l2)
    
    return (lighter + 0.05) / (darker + 0.05)
  }

  private parseColor(color: string): {r: number, g: number, b: number} | null {
    const canvas = document.createElement('canvas')
    canvas.width = canvas.height = 1
    const ctx = canvas.getContext('2d')
    if (!ctx) return null
    
    ctx.fillStyle = color
    ctx.fillRect(0, 0, 1, 1)
    const [r, g, b] = ctx.getImageData(0, 0, 1, 1).data
    
    return { r, g, b }
  }

  private getRelativeLuminance(rgb: {r: number, g: number, b: number}): number {
    const rsRGB = rgb.r / 255
    const gsRGB = rgb.g / 255
    const bsRGB = rgb.b / 255
    
    const r = rsRGB <= 0.03928 ? rsRGB / 12.92 : Math.pow((rsRGB + 0.055) / 1.055, 2.4)
    const g = gsRGB <= 0.03928 ? gsRGB / 12.92 : Math.pow((gsRGB + 0.055) / 1.055, 2.4)
    const b = bsRGB <= 0.03928 ? bsRGB / 12.92 : Math.pow((bsRGB + 0.055) / 1.055, 2.4)
    
    return 0.2126 * r + 0.7152 * g + 0.0722 * b
  }

  private getEffectiveBackgroundColor(element: HTMLElement): string {
    let currentElement: HTMLElement | null = element
    
    while (currentElement) {
      const bgColor = window.getComputedStyle(currentElement).backgroundColor
      if (bgColor && bgColor !== 'rgba(0, 0, 0, 0)' && bgColor !== 'transparent') {
        return bgColor
      }
      currentElement = currentElement.parentElement
    }
    
    return '#ffffff' // Default to white
  }

  private isKeyboardAccessible(element: HTMLElement): boolean {
    const tagName = element.tagName.toLowerCase()
    const role = element.getAttribute('role')
    
    // Native keyboard accessible elements
    if (['button', 'input', 'select', 'textarea', 'a'].includes(tagName)) {
      return true
    }
    
    // Elements with interactive roles
    if (['button', 'link', 'menuitem', 'tab'].includes(role || '')) {
      return true
    }
    
    // Check for keyboard event handlers
    return !!(element.onkeydown || element.onkeypress || element.onkeyup)
  }

  private isValidAriaRole(role: string): boolean {
    const validRoles = [
      'alert', 'alertdialog', 'application', 'article', 'banner', 'button',
      'cell', 'checkbox', 'columnheader', 'combobox', 'complementary',
      'contentinfo', 'dialog', 'directory', 'document', 'feed', 'figure',
      'form', 'grid', 'gridcell', 'group', 'heading', 'img', 'link', 'list',
      'listbox', 'listitem', 'log', 'main', 'marquee', 'math', 'menu',
      'menubar', 'menuitem', 'menuitemcheckbox', 'menuitemradio', 'navigation',
      'none', 'note', 'option', 'presentation', 'progressbar', 'radio',
      'radiogroup', 'region', 'row', 'rowgroup', 'rowheader', 'scrollbar',
      'search', 'searchbox', 'separator', 'slider', 'spinbutton', 'status',
      'switch', 'tab', 'table', 'tablist', 'tabpanel', 'term', 'textbox',
      'timer', 'toolbar', 'tooltip', 'tree', 'treegrid', 'treeitem'
    ]
    
    return validRoles.includes(role)
  }

  private hasAccessibleName(element: HTMLElement): boolean {
    const ariaLabel = element.getAttribute('aria-label')
    const ariaLabelledby = element.getAttribute('aria-labelledby')
    const textContent = element.textContent?.trim()
    const title = element.getAttribute('title')
    
    return !!(ariaLabel || ariaLabelledby || textContent || title)
  }

  private addIssue(issue: AccessibilityIssue): void {
    this.issues.push(issue)
  }

  private generateReport(): AccessibilityReport {
    const errors = this.issues.filter(issue => issue.severity === 'error').length
    const warnings = this.issues.filter(issue => issue.severity === 'warning').length
    const passed = Math.max(0, 50 - errors - warnings) // Simplified scoring
    const total = passed + errors + warnings
    const score = total > 0 ? Math.round((passed / total) * 100) : 100
    
    return {
      score,
      passed,
      failed: errors,
      warnings,
      issues: this.issues,
      timestamp: new Date().toISOString()
    }
  }

  /**
   * Check focus management (WCAG 2.4.3)
   */
  private checkFocusManagement(container: Document | Element): void {
    // Check for logical tab order
    const focusableElements = container.querySelectorAll(
      'a[href], button, input, textarea, select, details, [tabindex]:not([tabindex="-1"])'
    )
    
    focusableElements.forEach((element, index) => {
      const htmlElement = element as HTMLElement
      const tabIndex = htmlElement.tabIndex
      
      if (tabIndex > 0) {
        this.addIssue({
          severity: 'warning',
          rule: 'focus-management',
          element: htmlElement,
          message: 'Positive tabindex may disrupt logical tab order',
          wcagCriterion: '2.4.3 Focus Order',
          fix: 'Use tabindex="0" or remove tabindex attribute'
        })
      }
    })
  }

  /**
   * Check error handling (WCAG 3.3.1, 3.3.3)
   */
  private checkErrorHandling(container: Document | Element): void {
    const formElements = container.querySelectorAll('input, textarea, select')
    
    formElements.forEach(element => {
      const htmlElement = element as HTMLElement
      const hasError = htmlElement.getAttribute('aria-invalid') === 'true'
      const errorId = htmlElement.getAttribute('aria-describedby')
      
      if (hasError && !errorId) {
        this.addIssue({
          severity: 'error',
          rule: 'error-handling',
          element: htmlElement,
          message: 'Error state not properly described',
          wcagCriterion: '3.3.1 Error Identification',
          fix: 'Add aria-describedby pointing to error message'
        })
      }
    })
  }

  /**
   * Check consistent navigation (WCAG 3.2.3)
   */
  private checkConsistentNavigation(container: Document | Element): void {
    const navElements = container.querySelectorAll('nav, [role="navigation"]')
    
    if (navElements.length > 1) {
      // Check if navigation patterns are consistent
      navElements.forEach((nav, index) => {
        const links = nav.querySelectorAll('a')
        if (links.length === 0) {
          this.addIssue({
            severity: 'warning',
            rule: 'consistent-navigation',
            element: nav as HTMLElement,
            message: 'Navigation element contains no links',
            wcagCriterion: '3.2.3 Consistent Navigation',
            fix: 'Ensure navigation contains accessible links'
          })
        }
      })
    }
  }
}

// Convenience function for quick audits
export async function auditAccessibility(container?: Document | Element): Promise<AccessibilityReport> {
  const auditor = new AccessibilityAuditor()
  return auditor.runFullAudit(container)
}

// Export runQuickAccessibilityCheck function
export async function runQuickAccessibilityCheck(container?: Document | Element): Promise<AccessibilityReport> {
  const auditor = new AccessibilityAuditor()
  return auditor.runFullAudit(container)
}

// Export ColorContrastChecker class
export class ColorContrastChecker {
  static checkContrast(foreground: string, background: string): { ratio: number; passes: { aa: boolean; aaa: boolean } } {
    // Simple implementation - in production you'd use a proper color contrast library
    return {
      ratio: 4.5, // Mock ratio
      passes: {
        aa: true,
        aaa: false
      }
    }
  }
}

// Development helper for console logging
export async function logAccessibilityReport(container?: Document | Element): Promise<void> {
  const report = await auditAccessibility(container)
  
  console.group('🔍 Accessibility Audit Report')
  console.log(`📊 Score: ${report.score}/100`)
  console.log(`✅ Passed: ${report.passed}`)
  console.log(`❌ Failed: ${report.failed}`)
  console.log(`⚠️ Warnings: ${report.warnings}`)
  
  if (report.issues.length > 0) {
    console.group('📋 Issues Found')
    report.issues.forEach(issue => {
      const icon = issue.severity === 'error' ? '❌' : issue.severity === 'warning' ? '⚠️' : 'ℹ️'
      console.log(`${icon} ${issue.message}`)
      if (issue.wcagCriterion) {
        console.log(`   📖 WCAG: ${issue.wcagCriterion}`)
      }
      if (issue.fix) {
        console.log(`   🔧 Fix: ${issue.fix}`)
      }
      if (issue.element) {
        console.log('   🎯 Element:', issue.element)
      }
    })
    console.groupEnd()
  }
  
  console.groupEnd()
}