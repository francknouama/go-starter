# go-starter Design System Guide

A comprehensive guide to the visual design system for the go-starter web interface.

## 🎨 Design Principles

### Professional Trust
- Clean, consistent interfaces that build user confidence
- Subtle animations and micro-interactions
- Clear visual hierarchy and information architecture

### Progressive Disclosure
- Complexity revealed gradually as users need it
- Context-aware help and guidance
- Beginner-friendly with power user features

### Accessibility First
- WCAG 2.1 AA compliance throughout
- Keyboard navigation support
- High contrast and reduced motion options

### Performance Conscious
- Optimized animations (60fps)
- Efficient CSS and minimal bundle size
- Responsive design for all device sizes

## 🎯 Color System

### Primary Colors
```css
--color-primary-50: #eff6ff
--color-primary-500: #3b82f6  /* Main brand color */
--color-primary-600: #2563eb  /* Hover states */
--color-primary-700: #1d4ed8  /* Active states */
```

### Semantic Colors
```css
/* Success */
--color-success-500: #22c55e
--color-success-bg: #f0fdf4
--color-success-border: #bbf7d0

/* Warning */
--color-warning-500: #f59e0b
--color-warning-bg: #fffbeb
--color-warning-border: #fde68a

/* Error */
--color-error-500: #ef4444
--color-error-bg: #fef2f2
--color-error-border: #fecaca
```

### Usage Guidelines

#### When to Use Primary Colors
- Call-to-action buttons ("Generate Project")
- Active states and selected items
- Important links and interactive elements
- Focus states and accessibility indicators

#### When to Use Semantic Colors
- **Success**: Completed actions, valid form inputs, positive feedback
- **Warning**: Cautionary messages, potential issues, important notices
- **Error**: Form validation errors, failed operations, critical alerts

## 📝 Typography Scale

### Font Families
```css
/* Primary Interface */
font-family: 'Inter', system-ui, sans-serif;

/* Code and Technical Content */
font-family: 'JetBrains Mono', Monaco, Consolas, monospace;

/* Display Text */
font-family: 'Inter', system-ui, sans-serif;
```

### Typography Hierarchy
```css
/* Display Text (Hero sections) */
.text-display-xl: 48px / 1.1 / 700 weight
.text-display-lg: 36px / 1.1 / 700 weight
.text-display-md: 30px / 1.2 / 700 weight

/* Headings (Section titles) */
.text-heading-xl: 20px / 1.3 / 600 weight
.text-heading-lg: 18px / 1.4 / 600 weight
.text-heading-md: 16px / 1.4 / 600 weight

/* Body Text (Content) */
.text-body-lg: 16px / 1.6 / 400 weight
.text-body-md: 14px / 1.5 / 400 weight
.text-body-sm: 12px / 1.4 / 400 weight

/* Captions (Labels, metadata) */
.text-caption: 12px / 1.3 / 500 weight / uppercase
```

## 🔘 Button System

### Variants and Usage

#### Primary Button
```tsx
<Button variant="primary" size="md">
  Generate Project
</Button>
```
- **Use for**: Main actions, CTAs, form submissions
- **Visual**: Blue background, white text, shadow

#### Secondary Button
```tsx
<Button variant="secondary" size="md">
  Cancel
</Button>
```
- **Use for**: Secondary actions, alternative choices
- **Visual**: Gray background, dark text, border

#### Outline Button
```tsx
<Button variant="outline" size="md">
  Learn More
</Button>
```
- **Use for**: Less important actions, external links
- **Visual**: Transparent background, colored border and text

#### Ghost Button
```tsx
<Button variant="ghost" size="md">
  Skip
</Button>
```
- **Use for**: Minimal actions, close buttons, optional steps
- **Visual**: Transparent background, subtle hover state

### Button Sizes
```tsx
<Button size="xs">Extra Small</Button>   // 24px height
<Button size="sm">Small</Button>         // 32px height
<Button size="md">Medium</Button>        // 40px height (default)
<Button size="lg">Large</Button>         // 48px height
<Button size="xl">Extra Large</Button>   // 56px height
```

### Button States
- **Normal**: Default appearance
- **Hover**: Slightly darker background, smooth transition
- **Active**: Pressed state with scale(0.98) animation
- **Focus**: Visible focus ring for accessibility
- **Disabled**: 50% opacity, no pointer events
- **Loading**: Spinner animation, maintains width

## 🃏 Card System

### Card Types

#### Basic Card
```tsx
<div className="card p-6">
  Content here
</div>
```

#### Elevated Card
```tsx
<div className="card card-elevated p-6">
  Important content
</div>
```

#### Interactive Card
```tsx
<div className="card card-interactive p-6">
  Clickable content
</div>
```

#### Selected Card
```tsx
<div className="card card-selected p-6">
  Selected option
</div>
```

### Card Padding Guidelines
- **Small content**: `p-4` (16px)
- **Medium content**: `p-6` (24px) - Default
- **Large content**: `p-8` (32px)

## 📥 Input System

### Input Variants

#### Standard Input
```tsx
<input className="input" placeholder="Enter text..." />
```

#### Error State
```tsx
<input className="input input-error" placeholder="Fix this field" />
```

#### Success State
```tsx
<input className="input input-success" placeholder="Looks good!" />
```

### Input Sizes
```tsx
<input className="input input-sm" />  // 32px height
<input className="input" />           // 40px height (default)
<input className="input input-lg" />  // 48px height
```

### Label System
```tsx
<label className="label">Project Name</label>
<label className="label label-required">Required Field</label>
```

## 🎬 Loading States

### Loading Components

#### Spinner
```tsx
<Spinner size="md" color="primary" />
```

#### Progress Bar
```tsx
<ProgressBar value={65} showLabel label="Generating..." />
```

#### Skeleton Loading
```tsx
<Skeleton lines={3} />
<SkeletonCard showAvatar lines={4} />
```

#### Loading Overlay
```tsx
<LoadingOverlay isLoading={true} message="Generating project...">
  <ProjectContent />
</LoadingOverlay>
```

## 📊 Spacing System

### Spacing Scale (4px grid)
```css
--space-xs: 4px    // 0.25rem
--space-sm: 8px    // 0.5rem
--space-md: 16px   // 1rem
--space-lg: 24px   // 1.5rem
--space-xl: 32px   // 2rem
--space-2xl: 48px  // 3rem
```

### Common Spacing Patterns
- **Component padding**: `space-md` to `space-lg`
- **Section spacing**: `space-xl` to `space-2xl`
- **Element spacing**: `space-xs` to `space-sm`
- **Layout margins**: `space-lg` to `space-2xl`

## 🎭 Animation System

### Animation Durations
```css
--duration-fast: 150ms    // Quick interactions
--duration-normal: 200ms  // Default (buttons, hovers)
--duration-slow: 300ms    // Complex transitions
--duration-slower: 500ms  // Page transitions
```

### Easing Functions
```css
--easing-smooth: cubic-bezier(0.4, 0, 0.2, 1)    // Default
--easing-spring: cubic-bezier(0.68, -0.55, 0.265, 1.55)  // Bouncy
```

### Animation Usage
- **Buttons**: Hover transitions (200ms)
- **Modals**: Fade + slide animations (300ms)
- **Loading**: Continuous smooth animations
- **Micro-interactions**: Quick feedback (150ms)

## 🌙 Dark Mode Support

### Color Adaptation
```css
/* Light mode */
.card { @apply bg-white border-gray-200; }

/* Dark mode */
.dark .card { @apply bg-gray-800 border-gray-700; }
```

### Implementation
```tsx
// Toggle dark mode
<button onClick={() => setDarkMode(!darkMode)}>
  {darkMode ? <SunIcon /> : <MoonIcon />}
</button>
```

### Dark Mode Colors
- **Backgrounds**: Gray 900 → Gray 800 → Gray 700
- **Text**: White → Gray 200 → Gray 400
- **Borders**: Gray 700 → Gray 600

## ♿ Accessibility Guidelines

### Focus Management
- Visible focus rings on all interactive elements
- Focus trap in modals and overlays
- Logical tab order through components

### Screen Reader Support
```tsx
<button aria-label="Generate Go project">
  <GenerateIcon />
</button>

<div role="status" aria-live="polite">
  Project generated successfully!
</div>
```

### Color Contrast
- **Text on background**: 4.5:1 minimum ratio
- **Interactive elements**: 3:1 minimum ratio
- **Focus indicators**: High contrast borders

### Keyboard Navigation
- All functionality available via keyboard
- Escape key closes modals and overlays
- Arrow keys for option navigation
- Enter/Space for activation

## 📱 Responsive Design

### Breakpoints
```css
sm: 640px   // Small tablets
md: 768px   // Tablets
lg: 1024px  // Small laptops
xl: 1280px  // Laptops
2xl: 1536px // Large screens
```

### Mobile-First Approach
```tsx
<div className="flex flex-col lg:flex-row gap-4 lg:gap-6">
  {/* Stack on mobile, side-by-side on desktop */}
</div>
```

### Touch Targets
- Minimum 44px touch targets on mobile
- Adequate spacing between interactive elements
- Swipe gestures for mobile navigation

## 🎨 Icon System

### Icon Guidelines
- 16px for inline text icons
- 20px for button icons (default)
- 24px for toolbar and navigation icons
- 32px+ for feature illustrations

### Icon Usage
```tsx
import { ProjectTypeIcons } from '../icons/ProjectIcons'

<ProjectTypeCard 
  icon={<ProjectTypeIcons.WebAPI className="w-6 h-6" />}
  title="Web API"
/>
```

## 📏 Layout System

### Container Widths
```css
.container-narrow: max-width: 672px   // Reading content
.container-wide: max-width: 1280px    // Application layouts
```

### Grid System
```tsx
<div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
  {/* Responsive grid items */}
</div>
```

## 🔧 Implementation Guidelines

### Component Structure
```tsx
// Standard component pattern
interface ComponentProps {
  variant?: 'primary' | 'secondary'
  size?: 'sm' | 'md' | 'lg'
  className?: string
  children: ReactNode
}

export function Component({ 
  variant = 'primary', 
  size = 'md', 
  className = '',
  children,
  ...props 
}: ComponentProps) {
  const baseStyles = 'base component styles'
  const variantStyles = { primary: '...', secondary: '...' }
  const sizeStyles = { sm: '...', md: '...', lg: '...' }
  
  return (
    <div 
      className={`${baseStyles} ${variantStyles[variant]} ${sizeStyles[size]} ${className}`}
      {...props}
    >
      {children}
    </div>
  )
}
```

### Design Token Usage
```tsx
import { designTokens } from '../styles/design-tokens'

// Type-safe token access
const primaryColor = designTokens.colors.primary[500]
const spacing = designTokens.spacing.md
const radius = designTokens.radius.lg
```

### CSS Class Patterns
```css
/* Component base class */
.component { /* base styles */ }

/* Size variants */
.component-sm { /* small variant */ }
.component-md { /* medium variant (default) */ }
.component-lg { /* large variant */ }

/* State variants */
.component-primary { /* primary state */ }
.component-secondary { /* secondary state */ }
.component-disabled { /* disabled state */ }

/* Dark mode */
.dark .component { /* dark mode styles */ }
```

This design system ensures consistency, accessibility, and maintainability across the entire go-starter web interface while providing a professional, trustworthy experience for developers of all skill levels.