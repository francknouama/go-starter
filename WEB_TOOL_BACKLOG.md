# Free Web Tool Backlog

**Vision:** Browser-based go-starter with exact CLI feature parity - free forever  
**Timeline:** 3-4 weeks for MVP  
**Priority:** HIGH - Community value and adoption driver

---

## 🎯 Free Web Tool Philosophy

### **Core Principles**
1. **100% Free Forever** - No premium features, no paywalls
2. **No Registration** - Anonymous usage, no user tracking
3. **Exact CLI Feature Parity** - Only what the CLI currently supports
4. **Open Source** - Same license as CLI (MIT)
5. **Community First** - Built for developers by developers

### **Current CLI Scope (v1.3.1)**
The web tool will support ONLY what the CLI currently offers:

#### **4 Templates Available:**
1. **Web API** (web-api)
   - Framework: Gin only
   - Database: PostgreSQL, MySQL, MongoDB, SQLite, Redis
   - ORM: GORM or raw SQL only
2. **CLI Application** (cli)
   - Framework: Cobra only
3. **Library** (library)
   - Standard Go library structure
4. **AWS Lambda** (lambda)
   - Basic Lambda function

#### **4 Logger Options:**
- slog (default)
- zap
- logrus  
- zerolog

#### **Go Version Selection (Enhancement):**
- Go 1.23 (latest)
- Go 1.22
- Go 1.21
- Auto-detect current version (default)

#### **Configuration Options:**
- Project name
- Module path
- Go version (selector for latest 3 versions)
- Author name & email
- License (MIT, Apache-2.0, GPL-3.0, BSD-3-Clause)

---

## 🏗️ Technical Architecture (Simplified)

### **Frontend (Static Site)**
```
web-tool/
├── index.html          # Single page application
├── src/
│   ├── App.tsx         # Main React component
│   ├── components/
│   │   ├── TemplateSelector.tsx    # 4 templates only
│   │   ├── ConfigurationForm.tsx   # Current CLI options only
│   │   ├── GoVersionSelector.tsx   # Latest 3 Go versions
│   │   ├── FileTreePreview.tsx
│   │   └── DownloadButton.tsx
│   └── api/
│       └── generator.ts # API client
└── public/
    └── assets/         # Static assets
```

### **Backend (Stateless API)**
```go
// Minimal API - reuses exact CLI generation logic
type GenerateRequest struct {
    Template     string `json:"template"`      // web-api, cli, library, lambda
    ProjectName  string `json:"projectName"`
    ModulePath   string `json:"modulePath"`
    GoVersion    string `json:"goVersion"`     // "1.23", "1.22", "1.21", "auto"
    Author       string `json:"author"`
    Email        string `json:"email"`
    License      string `json:"license"`
    Logger       string `json:"logger"`        // slog, zap, logrus, zerolog
    // Web API specific
    Database     string `json:"database,omitempty"`
    DatabaseORM  string `json:"databaseORM,omitempty"`
}

// Single endpoint needed
POST /api/v1/generate → Returns zip file
```

---

## 📋 Implementation Plan (3-4 Weeks)

### Week 1: Backend API

#### 🔧 Feature: Stateless Generation API
**Estimate:** 3-5 days  
**Priority:** HIGH

**Tasks:**
- [ ] **Extract Generation Logic**
  - [ ] Create shared package from CLI generator
  - [ ] Remove file system dependencies
  - [ ] Support in-memory zip creation
  - [ ] Ensure thread safety for concurrent requests

- [ ] **Go Version Management**
  ```go
  // Supported Go versions (update quarterly)
  var SupportedGoVersions = []string{
      "1.23", // Latest
      "1.22", 
      "1.21",
      "auto", // Detect from system
  }
  
  func validateGoVersion(version string) error {
      if !contains(SupportedGoVersions, version) {
          return fmt.Errorf("unsupported Go version: %s", version)
      }
      return nil
  }
  ```

- [ ] **API Endpoint**
  ```go
  func GenerateHandler(c *gin.Context) {
      var req GenerateRequest
      if err := c.ShouldBindJSON(&req); err != nil {
          c.JSON(400, gin.H{"error": "Invalid request"})
          return
      }
      
      // Validate template is one of the 4 supported
      validTemplates := []string{"web-api", "cli", "library", "lambda"}
      if !contains(validTemplates, req.Template) {
          c.JSON(400, gin.H{"error": "Invalid template"})
          return
      }
      
      // Validate Go version
      if err := validateGoVersion(req.GoVersion); err != nil {
          c.JSON(400, gin.H{"error": err.Error()})
          return
      }
      
      // Reuse exact CLI generation logic
      zipData, err := generator.GenerateProject(req)
      if err != nil {
          c.JSON(500, gin.H{"error": err.Error()})
          return
      }
      
      // Stream zip file
      c.Header("Content-Type", "application/zip")
      c.Header("Content-Disposition", "attachment; filename=" + req.ProjectName + ".zip")
      c.Data(200, "application/zip", zipData)
  }
  ```

### Week 2: Frontend Interface

#### 🎨 Feature: React Web Interface
**Estimate:** 5-7 days  
**Priority:** HIGH

**Tasks:**
- [ ] **Template Selection (4 templates only)**
  ```tsx
  const templates = [
    {
      id: 'web-api',
      icon: '🌐',
      title: 'Web API',
      description: 'REST API with Gin framework and database support'
    },
    {
      id: 'cli',
      icon: '💻',
      title: 'CLI Application',
      description: 'Command-line tool with Cobra framework'
    },
    {
      id: 'library',
      icon: '📚',
      title: 'Go Library',
      description: 'Reusable Go package with examples'
    },
    {
      id: 'lambda',
      icon: '⚡',
      title: 'AWS Lambda',
      description: 'Serverless function for AWS Lambda'
    }
  ];
  ```

- [ ] **Go Version Selector Component**
  ```tsx
  const GoVersionSelector = () => {
    const versions = [
      { value: 'auto', label: 'Auto-detect (recommended)' },
      { value: '1.23', label: 'Go 1.23 (latest)' },
      { value: '1.22', label: 'Go 1.22' },
      { value: '1.21', label: 'Go 1.21' }
    ];
    
    return (
      <select name="goVersion" defaultValue="auto">
        {versions.map(v => (
          <option key={v.value} value={v.value}>{v.label}</option>
        ))}
      </select>
    );
  };
  ```

- [ ] **Configuration Form (matching CLI exactly)**
  - [ ] Common fields for all templates:
    - Project name
    - Module path  
    - Go version selector (1.23, 1.22, 1.21, auto)
    - Logger selection (slog, zap, logrus, zerolog)
    - Author name
    - Email
    - License
  - [ ] Web API specific:
    - Database selection (PostgreSQL, MySQL, MongoDB, SQLite, Redis)
    - ORM selection (GORM, raw) - only if SQL database selected
  - [ ] No other framework/library options

- [ ] **Live Preview**
  - [ ] Show exact file structure that will be generated
  - [ ] Display selected Go version in go.mod preview
  - [ ] Match CLI output exactly
  - [ ] No additional features

### Week 3: Polish & Deployment

#### 🚀 Feature: Production Deployment
**Estimate:** 3-5 days  
**Priority:** HIGH

**Tasks:**
- [ ] **Version Management Strategy**
  - [ ] Quarterly updates to supported Go versions
  - [ ] Automated script to update version list
  - [ ] Backward compatibility for generated projects
  - [ ] Clear documentation on version support

- [ ] **Deployment**
  - [ ] Frontend on GitHub Pages
  - [ ] Backend on free tier (Railway/Render)
  - [ ] No database needed
  - [ ] Simple monitoring

---

## 🎨 UI Design (Minimal)

### **Simple Interface**
```
┌─────────────────────────────────────────────────┐
│  🚀 go-starter web                              │
│  Generate Go projects in your browser           │
├─────────────────────────────────────────────────┤
│                                                  │
│  Select Template:                                │
│  ○ Web API   ○ CLI   ○ Library   ○ Lambda      │
│                                                  │
│  Project Name: [____________________]           │
│  Module Path:  [____________________]           │
│  Go Version:   [Auto-detect      ▼]             │
│  Logger:       [slog             ▼]             │
│                                                  │
│  [Show Advanced Options]                         │
│                                                  │
│         [Generate Project →]                     │
│                                                  │
└─────────────────────────────────────────────────┘
```

---

## 📊 Go Version Support Strategy

### **Version Selection Logic**
```go
// Update quarterly when new Go versions are released
const (
    GoVersionLatest = "1.23"
    GoVersionMinus1 = "1.22" 
    GoVersionMinus2 = "1.21"
)

// CLI Enhancement (to be implemented)
func GetSupportedGoVersions() []string {
    return []string{
        "auto",           // Detect from system
        GoVersionLatest,  // Always latest stable
        GoVersionMinus1,  // Previous version
        GoVersionMinus2,  // Two versions back
    }
}
```

### **Update Process**
1. When Go 1.24 is released:
   - Update constants in code
   - Drop Go 1.21 support
   - Add Go 1.24 as latest
2. Maintain compatibility for 3 major versions
3. Update both CLI and web tool simultaneously

---

## 📊 Success Metrics

### **Simple Goals**
- Works exactly like the CLI
- Go version selector improves user experience
- No bugs or feature drift
- Fast generation (< 3 seconds)
- Available globally via CDN

---

## ⚠️ Scope Boundaries

### **ONLY includes:**
- 4 current templates (web-api, cli, library, lambda)
- 4 logger options (slog, zap, logrus, zerolog)
- Go version selector (latest 3 versions + auto)
- Current database options for web-api
- GORM or raw SQL only (no other ORMs)
- Gin framework only for web-api
- Cobra framework only for CLI

### **Does NOT include:**
- Echo, Fiber, Chi frameworks
- Additional ORMs (sqlx, sqlc, ent)
- Architecture patterns (these are future templates)
- Any features not in current CLI v1.3.1 (except Go version selector)

---

## 🔄 CLI Enhancement Required

The Go version selector should also be added to the CLI for consistency:
- Update CLI prompts to include Go version selection
- Default to "auto" (current behavior)
- Allow selection of latest 3 major versions
- Maintain backward compatibility

---

The web tool is a direct mirror of the CLI with the practical enhancement of Go version selection.