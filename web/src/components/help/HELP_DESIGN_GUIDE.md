# Help System Design Guide

## Visual Design Principles

### 1. Progressive Disclosure
- **Basic Mode**: Show only essential help hints (14 essential options)
- **Advanced Mode**: Reveal comprehensive help (18+ options)
- **Smart Suggestions**: Context-aware tips based on user selections
- **On-Demand Help**: Full help available but not intrusive

### 2. Help Indicator Design

#### Subtle Integration
```css
/* Help icons should be subtle but discoverable */
.help-icon {
  color: #9CA3AF; /* Gray-400 */
  transition: color 0.2s ease;
}

.help-icon:hover {
  color: #6B7280; /* Gray-500 */
}
```

#### Visual Hierarchy
- **Primary**: Question mark icons for complex concepts
- **Secondary**: Info icons for additional context
- **Tertiary**: Light bulb icons for tips and suggestions

### 3. Tooltip Design System

#### Positioning Strategy
- **Auto-positioning**: Smart viewport-aware positioning
- **Fallback Order**: Bottom → Top → Right → Left
- **Minimum Clearance**: 150px vertical, 200px horizontal

#### Visual Styling
```css
.help-tooltip {
  background: rgba(31, 41, 55, 0.95); /* Gray-800 with opacity */
  backdrop-filter: blur(8px);
  border: 1px solid rgba(75, 85, 99, 0.3);
  border-radius: 8px;
  padding: 12px;
  max-width: 280px;
  font-size: 12px;
  line-height: 1.4;
  z-index: 50;
}
```

### 4. Color System

#### Help System Colors
- **Info**: Blue-500 (#3B82F6) - General information
- **Tips**: Yellow-500 (#EAB308) - Helpful suggestions
- **Warnings**: Orange-500 (#F97316) - Important considerations
- **Success**: Green-500 (#10B981) - Positive reinforcement

#### Accessibility
- **Contrast Ratios**: Minimum 4.5:1 for normal text
- **Focus Indicators**: 2px blue-500 outline with 2px offset
- **Color Independence**: Never rely solely on color for meaning

### 5. Animation Guidelines

#### Micro-interactions
```css
/* Smooth entrance animations */
.help-enter {
  transition: all 200ms ease-out;
  transform: scale(0.95);
  opacity: 0;
}

.help-enter-active {
  transform: scale(1);
  opacity: 1;
}

/* Gentle exit animations */
.help-exit {
  transition: all 150ms ease-in;
  transform: scale(1);
  opacity: 1;
}

.help-exit-active {
  transform: scale(0.95);
  opacity: 0;
}
```

#### Delay Strategy
- **Hover Delay**: 200ms to prevent accidental triggers
- **Exit Delay**: 100ms to allow cursor movement
- **Modal Animations**: 300ms for large overlays

### 6. Responsive Design

#### Mobile Adaptations
- **Touch Targets**: Minimum 44px for touch interfaces
- **Modal Fullscreen**: Quick start guide fills viewport on mobile
- **Simplified Tooltips**: Reduced content on small screens
- **Gesture Support**: Swipe to dismiss on touch devices

#### Breakpoint Considerations
```css
/* Mobile First Approach */
@media (max-width: 640px) {
  .help-tooltip {
    max-width: calc(100vw - 32px);
    font-size: 14px;
  }
}

@media (min-width: 768px) {
  .help-system {
    /* Enhanced features for larger screens */
  }
}
```

### 7. Performance Optimization

#### Lazy Loading
- **Component Splitting**: Help overlays load on demand
- **Image Optimization**: SVG icons for crisp rendering
- **Bundle Splitting**: Separate help system from main bundle

#### Memory Management
- **Event Cleanup**: Remove event listeners on unmount
- **Timeout Management**: Clear timeouts to prevent memory leaks
- **DOM Optimization**: Minimize DOM manipulation

### 8. Accessibility Standards

#### WCAG 2.1 AA Compliance
- **Keyboard Navigation**: All help elements accessible via keyboard
- **Screen Reader Support**: Proper ARIA labels and descriptions
- **Focus Management**: Logical focus order and restoration
- **High Contrast Support**: Adapts to user preferences

#### Implementation
```tsx
// Accessible help tooltip
<HelpTooltip
  content="Detailed explanation"
  aria-label="Help for project name field"
  role="tooltip"
  tabIndex={0}
/>
```

### 9. Content Strategy

#### Writing Guidelines
- **Concise**: Maximum 2-3 sentences per tooltip
- **Actionable**: Tell users what to do, not just what something is
- **Progressive**: Layer information from basic to advanced
- **Contextual**: Adapt content based on user's current state

#### Content Hierarchy
1. **Essential**: Must-know information for task completion
2. **Helpful**: Nice-to-know tips that improve experience
3. **Advanced**: Expert-level information for power users

### 10. User Experience Flow

#### Discovery Pattern
1. **Visual Cues**: Subtle help icons indicate help availability
2. **Progressive Reveal**: Basic help leads to more detailed information
3. **Just-in-Time**: Help appears when and where needed
4. **Non-Blocking**: Help enhances but never prevents task completion

#### Interaction Model
- **Hover**: Quick tips and contextual help
- **Click**: Detailed explanations and guides
- **Keyboard**: Full accessibility with shortcuts
- **Touch**: Tap-friendly interactions on mobile

### 11. Error Prevention

#### Smart Validation
- **Real-time Feedback**: Immediate validation with helpful messages
- **Progressive Enhancement**: Validation improves as user types
- **Clear Recovery**: Easy correction paths for errors
- **Contextual Help**: Error-specific help suggestions

#### Example Implementation
```tsx
<ValidatedInput
  value={projectName}
  validation={{
    required: true,
    pattern: /^[a-z0-9-]+$/,
    message: "Use lowercase letters, numbers, and hyphens only"
  }}
  helpContent="This becomes your project directory name"
/>
```

### 12. Testing Strategy

#### Usability Testing
- **Task Completion**: Can users complete tasks without help?
- **Help Discovery**: Do users find help when they need it?
- **Information Architecture**: Is help content well-organized?
- **Performance**: Does help system impact app performance?

#### Accessibility Testing
- **Screen Readers**: Test with NVDA, JAWS, VoiceOver
- **Keyboard Only**: Complete workflows without mouse
- **High Contrast**: Verify visibility in high contrast mode
- **Zoom Testing**: Functionality at 200% zoom level

This design system ensures the help system is discoverable, accessible, and enhances rather than clutters the user interface while following progressive disclosure principles.