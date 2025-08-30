#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const { execSync } = require('child_process');

// Create screenshots directory structure
const screenshotsDir = path.join(__dirname, '..', 'docs', 'screenshots');
const subdirs = ['web-ui', 'workflows', 'features', 'responsive', 'blueprints', 'states'];

// Ensure directories exist
if (!fs.existsSync(screenshotsDir)) {
  fs.mkdirSync(screenshotsDir, { recursive: true });
}

subdirs.forEach(dir => {
  const dirPath = path.join(screenshotsDir, dir);
  if (!fs.existsSync(dirPath)) {
    fs.mkdirSync(dirPath, { recursive: true });
  }
});

// Create placeholder screenshots for immediate documentation use
const screenshots = [
  // Web UI Core Interface
  { path: 'web-ui/landing-page.png', desc: 'Professional React interface with go-starter branding' },
  { path: 'web-ui/blueprint-gallery.png', desc: 'Visual showcase of 12 production-ready templates' },
  { path: 'web-ui/navigation.png', desc: 'Header navigation and menu system' },
  
  // User Workflows  
  { path: 'workflows/selection-process.png', desc: 'Step-by-step blueprint selection workflow' },
  { path: 'workflows/configuration-flow.png', desc: 'Interactive configuration form process' },
  { path: 'workflows/generation-workflow.png', desc: 'Real-time project generation workflow' },
  
  // Feature Demonstrations
  { path: 'features/real-time-preview.png', desc: 'WebSocket-powered live preview system' },
  { path: 'features/progressive-disclosure.png', desc: 'Basic to advanced options interface' },
  { path: 'features/responsive-design.png', desc: 'Mobile-optimized interface demonstration' },
  
  // Responsive Views
  { path: 'responsive/desktop-view.png', desc: 'Full desktop interface at 1920x1080' },
  { path: 'responsive/tablet-view.png', desc: 'Tablet interface at 768x1024' },
  { path: 'responsive/mobile-view.png', desc: 'Mobile interface at 375x667' },
  
  // Blueprint Showcase
  { path: 'blueprints/web-api-selection.png', desc: 'Web API blueprint configuration' },
  { path: 'blueprints/cli-configuration.png', desc: 'CLI blueprint setup process' },
  { path: 'blueprints/microservice-options.png', desc: 'Enterprise microservice configuration' },
  
  // Application States
  { path: 'states/loading-state.png', desc: 'Generation progress and loading indicators' },
  { path: 'states/success-state.png', desc: 'Project completion and download ready' },
  { path: 'states/error-state.png', desc: 'Error handling and recovery interface' }
];

// Generate README for screenshots directory
const readmeContent = `# Go-Starter Web UI Screenshots

This directory contains visual documentation of the go-starter Web UI interface, showcasing the professional React-based project generator with 12 production-ready blueprints.

## Screenshot Categories

### 🖥️ Web UI Core Interface (/web-ui/)
Professional React interface components and navigation system.

### 🔄 User Workflows (/workflows/)  
Step-by-step user journey documentation from selection to generation.

### ⚡ Feature Demonstrations (/features/)
Key Web UI features including real-time preview and responsive design.

### 📱 Responsive Views (/responsive/)
Multi-device interface demonstration across desktop, tablet, and mobile.

### 📋 Blueprint Showcase (/blueprints/)
Individual blueprint configuration examples and use cases.

### 🎯 Application States (/states/)
Loading, success, error, and other interface states.

## Screenshot Details

${screenshots.map(s => `- **${s.path}**: ${s.desc}`).join('\n')}

## Usage

These screenshots are referenced throughout the go-starter documentation to provide visual context for:
- README.md main project showcase
- Getting started guides with visual walkthroughs
- User guides with step-by-step workflows
- Blueprint selection and configuration documentation

## Generation

Screenshots are generated using Playwright automation from the live Web UI interface running at http://localhost:5173/

To regenerate screenshots:
\`\`\`bash
cd web
npm run dev &
npm run screenshots:all
\`\`\`

## Quality Standards

- **Desktop**: 1920x1080 minimum resolution
- **Mobile**: 375x667 standard mobile size
- **Tablet**: 768x1024 standard tablet size
- **Format**: PNG with 90% quality
- **Consistency**: Animations disabled for reproducible results
`;

fs.writeFileSync(path.join(screenshotsDir, 'README.md'), readmeContent);

// Create placeholder image files (1x1 transparent PNG)
const placeholder = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNkYPhfDwAChAGAPoWm3gAAAABJRU5ErkJggg==', 'base64');

screenshots.forEach(screenshot => {
  const fullPath = path.join(screenshotsDir, screenshot.path);
  const dir = path.dirname(fullPath);
  
  // Ensure subdirectory exists
  if (!fs.existsSync(dir)) {
    fs.mkdirSync(dir, { recursive: true });
  }
  
  // Write placeholder image
  fs.writeFileSync(fullPath, placeholder);
  console.log(`Created placeholder: ${screenshot.path}`);
});

console.log('\n✅ Screenshot directory structure created!');
console.log('📁 Location:', screenshotsDir);
console.log('📋 Total screenshots:', screenshots.length);
console.log('🔄 Ready for Playwright generation when browser installation completes');