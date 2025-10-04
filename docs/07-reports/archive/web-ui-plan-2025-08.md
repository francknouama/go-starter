# Go-Starter Phase 3: Web UI Implementation Plan

## Overview
Create a modern web interface for go-starter inspired by Spring Initializr, featuring progressive disclosure, live preview, real-time project generation, and comprehensive logger selector integration across all project types.

## Project Structure
```
go-starter-web/
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/
│   │   ├── handlers/
│   │   │   ├── generate.go
│   │   │   ├── preview.go
│   │   │   ├── templates.go
│   │   │   └── websocket.go
│   │   ├── middleware/
│   │   │   ├── cors.go
│   │   │   ├── logger.go
│   │   │   └── auth.go
│   │   └── server.go
│   ├── config/
│   │   └── config.go
│   ├── generator/
│   │   └── service.go
│   └── storage/
│       ├── projects.go
│       └── templates.go
├── web/
│   ├── static/
│   │   ├── css/
│   │   ├── js/
│   │   └── images/
│   ├── templates/
│   │   └── index.html
│   └── src/
│       ├── components/
│       │   ├── ConfigPanel/
│       │   ├── PreviewPanel/
│       │   ├── ProgressiveForm/
│       │   └── LivePreview/
│       ├── hooks/
│       ├── services/
│       ├── utils/
│       ├── App.jsx
│       └── index.js
├── Dockerfile
├── docker-compose.yml
├── package.json
├── vite.config.js
└── README.md
```

## Backend Implementation

### 1. API Server Setup

**main.go**:
```go
package main

import (
    "log"
    
    "github.com/username/go-starter-web/internal/api"
    "github.com/username/go-starter-web/internal/config"
)

func main() {
    cfg, err := config.Load()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }
    
    server := api.NewServer(cfg)
    
    log.Printf("Starting server on port %s", cfg.Port)
    if err := server.Start(); err != nil {
        log.Fatalf("Server failed to start: %v", err)
    }
}
```

**internal/api/server.go**:
```go
package api

import (
    "net/http"
    "time"

    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    
    "github.com/username/go-starter-web/internal/api/handlers"
    "github.com/username/go-starter-web/internal/api/middleware"
    "github.com/username/go-starter-web/internal/config"
    "github.com/username/go-starter-web/internal/generator"
)

type Server struct {
    config    *config.Config
    router    *gin.Engine
    generator *generator.Service
    upgrader  websocket.Upgrader
}

func NewServer(cfg *config.Config) *Server {
    if cfg.Environment == "production" {
        gin.SetMode(gin.ReleaseMode)
    }
    
    router := gin.Default()
    
    // Configure CORS
    router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000", "http://localhost:5173"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))
    
    // Middleware
    router.Use(middleware.Logger())
    router.Use(middleware.ErrorHandler())
    
    genService := generator.New()
    
    return &Server{
        config:    cfg,
        router:    router,
        generator: genService,
        upgrader: websocket.Upgrader{
            CheckOrigin: func(r *http.Request) bool {
                return true // Allow all origins in development
            },
        },
    }
}

func (s *Server) setupRoutes() {
    api := s.router.Group("/api/v1")
    
    // Template endpoints
    templateHandler := handlers.NewTemplateHandler()
    api.GET("/templates", templateHandler.List)
    api.GET("/templates/:id", templateHandler.Get)
    
    // Generation endpoints  
    generateHandler := handlers.NewGenerateHandler(s.generator)
    api.POST("/generate", generateHandler.Generate)
    api.POST("/preview", generateHandler.Preview)
    api.GET("/download/:projectId", generateHandler.Download)
    
    // Configuration endpoints
    configHandler := handlers.NewConfigHandler()
    api.POST("/validate", configHandler.Validate)
    api.POST("/share", configHandler.Share)
    api.GET("/share/:id", configHandler.GetShared)
    
    // WebSocket for live preview
    previewHandler := handlers.NewPreviewHandler(s.upgrader)
    api.GET("/preview/live", previewHandler.HandleWebSocket)
    
    // Static files (for production)
    s.router.Static("/static", "./web/dist/static")
    s.router.StaticFile("/", "./web/dist/index.html")
    s.router.NoRoute(func(c *gin.Context) {
        c.File("./web/dist/index.html")
    })
}

func (s *Server) Start() error {
    s.setupRoutes()
    return s.router.Run(":" + s.config.Port)
}
```

### 2. API Handlers

**internal/api/handlers/generate.go**:
```go
package handlers

import (
    "archive/zip"
    "bytes"
    "fmt"
    "io"
    "net/http"
    "path/filepath"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    
    "github.com/username/go-starter-web/internal/generator"
    "github.com/username/go-starter-web/pkg/types"
)

type GenerateHandler struct {
    generator *generator.Service
    projects  map[string]*ProjectInfo // In-memory storage for demo
}

type ProjectInfo struct {
    ID       string
    Config   types.ProjectConfig
    Files    map[string][]byte
    Created  time.Time
}

type GenerateRequest struct {
    Config types.ProjectConfig `json:"config" binding:"required"`
}

type GenerateResponse struct {
    ProjectID string                    `json:"project_id"`
    Files     map[string]string         `json:"files"`
    Structure []FileNode                `json:"structure"`
    Metadata  map[string]interface{}    `json:"metadata"`
}

type FileNode struct {
    Name     string     `json:"name"`
    Type     string     `json:"type"` // "file" or "directory"
    Path     string     `json:"path"`
    Size     int64      `json:"size,omitempty"`
    Children []FileNode `json:"children,omitempty"`
}

func NewGenerateHandler(gen *generator.Service) *GenerateHandler {
    return &GenerateHandler{
        generator: gen,
        projects:  make(map[string]*ProjectInfo),
    }
}

func (h *GenerateHandler) Generate(c *gin.Context) {
    var req GenerateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validate configuration
    if err := h.validateConfig(req.Config); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Generate project
    projectID := uuid.New().String()
    files, err := h.generator.GenerateInMemory(req.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Store project info
    projectInfo := &ProjectInfo{
        ID:      projectID,
        Config:  req.Config,
        Files:   files,
        Created: time.Now(),
    }
    h.projects[projectID] = projectInfo
    
    // Build response
    response := GenerateResponse{
        ProjectID: projectID,
        Files:     h.convertFilesToStrings(files),
        Structure: h.buildFileStructure(files),
        Metadata: map[string]interface{}{
            "template_used": h.getTemplateID(req.Config),
            "file_count":    len(files),
            "generated":     time.Now(),
        },
    }
    
    c.JSON(http.StatusOK, response)
}

func (h *GenerateHandler) Preview(c *gin.Context) {
    var req GenerateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Generate preview (structure only, no file contents)
    structure, err := h.generator.GenerateStructure(req.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "structure": structure,
        "metadata": map[string]interface{}{
            "template_used": h.getTemplateID(req.Config),
            "estimated_file_count": len(structure),
        },
    })
}

func (h *GenerateHandler) Download(c *gin.Context) {
    projectID := c.Param("projectId")
    
    project, exists := h.projects[projectID]
    if !exists {
        c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
        return
    }
    
    // Create ZIP file
    zipBuffer := new(bytes.Buffer)
    zipWriter := zip.NewWriter(zipBuffer)
    
    for path, content := range project.Files {
        writer, err := zipWriter.Create(path)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
        
        _, err = writer.Write(content)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }
    }
    
    zipWriter.Close()
    
    // Set headers for file download
    c.Header("Content-Type", "application/zip")
    c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s.zip", project.Config.Name))
    c.Header("Content-Length", fmt.Sprintf("%d", zipBuffer.Len()))
    
    // Stream the zip file
    c.DataFromReader(http.StatusOK, int64(zipBuffer.Len()), "application/zip", zipBuffer, nil)
}

func (h *GenerateHandler) validateConfig(config types.ProjectConfig) error {
    if config.Name == "" {
        return fmt.Errorf("project name is required")
    }
    if config.Module == "" {
        return fmt.Errorf("module path is required")
    }
    if config.Type == "" {
        return fmt.Errorf("project type is required")
    }
    return nil
}

func (h *GenerateHandler) convertFilesToStrings(files map[string][]byte) map[string]string {
    result := make(map[string]string)
    for path, content := range files {
        result[path] = string(content)
    }
    return result
}

func (h *GenerateHandler) buildFileStructure(files map[string][]byte) []FileNode {
    // Build hierarchical file structure from flat file map
    // Implementation would create tree structure...
    return []FileNode{} // Simplified for now
}

func (h *GenerateHandler) getTemplateID(config types.ProjectConfig) string {
    // Map config to template ID
    return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
}
```

**internal/api/handlers/websocket.go**:
```go
package handlers

import (
    "encoding/json"
    "log"
    "net/http"

    "github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
    
    "github.com/username/go-starter-web/pkg/types"
)

type PreviewHandler struct {
    upgrader websocket.Upgrader
}

type PreviewMessage struct {
    Type   string              `json:"type"`
    Config types.ProjectConfig `json:"config,omitempty"`
    Error  string              `json:"error,omitempty"`
    Data   interface{}         `json:"data,omitempty"`
}

func NewPreviewHandler(upgrader websocket.Upgrader) *PreviewHandler {
    return &PreviewHandler{
        upgrader: upgrader,
    }
}

func (h *PreviewHandler) HandleWebSocket(c *gin.Context) {
    conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        log.Printf("WebSocket upgrade failed: %v", err)
        return
    }
    defer conn.Close()
    
    log.Println("WebSocket connection established")
    
    for {
        var msg PreviewMessage
        err := conn.ReadJSON(&msg)
        if err != nil {
            if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
                log.Printf("WebSocket error: %v", err)
            }
            break
        }
        
        switch msg.Type {
        case "preview":
            h.handlePreviewRequest(conn, msg.Config)
        case "validate":
            h.handleValidationRequest(conn, msg.Config)
        }
    }
}

func (h *PreviewHandler) handlePreviewRequest(conn *websocket.Conn, config types.ProjectConfig) {
    // Generate real-time preview
    // This would integrate with the generator service
    
    response := PreviewMessage{
        Type: "preview_result",
        Data: map[string]interface{}{
            "structure": []string{
                "go.mod",
                "go.sum", 
                "cmd/server/main.go",
                "internal/handler/handler.go",
                "README.md",
            },
            "fileCount": 5,
        },
    }
    
    conn.WriteJSON(response)
}

func (h *PreviewHandler) handleValidationRequest(conn *websocket.Conn, config types.ProjectConfig) {
    // Validate configuration in real-time
    errors := []string{}
    
    if config.Name == "" {
        errors = append(errors, "Project name is required")
    }
    if config.Module == "" {
        errors = append(errors, "Module path is required")
    }
    
    response := PreviewMessage{
        Type: "validation_result",
        Data: map[string]interface{}{
            "valid":  len(errors) == 0,
            "errors": errors,
        },
    }
    
    conn.WriteJSON(response)
}
```

## Frontend Implementation

### 1. React App Setup

**package.json**:
```json
{
  "name": "go-starter-web",
  "version": "1.0.0",
  "type": "module",
  "scripts": {
    "dev": "vite",
    "build": "vite build",
    "preview": "vite preview",
    "lint": "eslint . --ext js,jsx --report-unused-disable-directives --max-warnings 0"
  },
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "axios": "^1.6.0",
    "react-hook-form": "^7.48.0",
    "@hookform/resolvers": "^3.3.0",
    "yup": "^1.3.0",
    "react-query": "^3.39.0",
    "react-syntax-highlighter": "^15.5.0",
    "react-select": "^5.8.0",
    "react-toggle": "^4.1.0",
    "lucide-react": "^0.263.1",
    "clsx": "^2.0.0"
  },
  "devDependencies": {
    "@types/react": "^18.2.37",
    "@types/react-dom": "^18.2.15",
    "@vitejs/plugin-react": "^4.1.0",
    "vite": "^4.5.0",
    "eslint": "^8.53.0",
    "eslint-plugin-react": "^7.33.0",
    "eslint-plugin-react-hooks": "^4.6.0",
    "eslint-plugin-react-refresh": "^0.4.4",
    "tailwindcss": "^3.3.0",
    "autoprefixer": "^10.4.16",
    "postcss": "^8.4.31"
  }
}
```

**vite.config.js**:
```javascript
import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
    },
  },
  build: {
    outDir: 'dist',
    assetsDir: 'static',
  },
})
```

### 2. Main App Component

**web/src/App.jsx**:
```jsx
import React, { useState, useCallback } from 'react'
import { QueryClient, QueryClientProvider } from 'react-query'
import ConfigPanel from './components/ConfigPanel/ConfigPanel'
import PreviewPanel from './components/PreviewPanel/PreviewPanel'
import Header from './components/Header/Header'
import Footer from './components/Footer/Footer'
import { useProjectGeneration } from './hooks/useProjectGeneration'
import { useWebSocket } from './hooks/useWebSocket'
import './App.css'

const queryClient = new QueryClient()

function AppContent() {
  const [config, setConfig] = useState({
    name: '',
    module: '',
    type: 'api',
    goVersion: '1.21',
    architecture: 'standard',
    framework: 'gin',
    features: {
      database: { driver: '', orm: '' },
      authentication: { type: '', providers: [] },
      deployment: { targets: [] },
      testing: { framework: 'testify', coverage: true },
      logger: { type: 'slog', level: 'info', format: 'json' }
    }
  })
  
  const [isAdvanced, setIsAdvanced] = useState(false)
  const [selectedFile, setSelectedFile] = useState(null)
  
  const {
    generateProject,
    downloadProject,
    isGenerating,
    generatedProject,
    error
  } = useProjectGeneration()
  
  const {
    connect,
    disconnect,
    sendPreview,
    previewData,
    isConnected
  } = useWebSocket('ws://localhost:8080/api/v1/preview/live')
  
  const handleConfigChange = useCallback((newConfig) => {
    setConfig(newConfig)
    
    // Send real-time preview update
    if (isConnected && newConfig.name && newConfig.module) {
      sendPreview(newConfig)
    }
  }, [isConnected, sendPreview])
  
  const handleGenerate = async () => {
    try {
      const result = await generateProject(config)
      if (result.projectId) {
        // Auto-select main file for preview
        setSelectedFile('cmd/server/main.go')
      }
    } catch (err) {
      console.error('Generation failed:', err)
    }
  }
  
  const handleDownload = async () => {
    if (generatedProject?.projectId) {
      await downloadProject(generatedProject.projectId, config.name)
    }
  }
  
  return (
    <div className="min-h-screen bg-gray-50 flex flex-col">
      <Header />
      
      <main className="flex-1 container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 h-full">
          {/* Configuration Panel */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <ConfigPanel
              config={config}
              onChange={handleConfigChange}
              isAdvanced={isAdvanced}
              onToggleAdvanced={() => setIsAdvanced(!isAdvanced)}
              onGenerate={handleGenerate}
              isGenerating={isGenerating}
              error={error}
            />
          </div>
          
          {/* Preview Panel */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <PreviewPanel
              projectStructure={generatedProject?.structure || previewData?.structure}
              projectFiles={generatedProject?.files}
              selectedFile={selectedFile}
              onFileSelect={setSelectedFile}
              onDownload={handleDownload}
              canDownload={!!generatedProject?.projectId}
              isGenerating={isGenerating}
            />
          </div>
        </div>
      </main>
      
      <Footer />
    </div>
  )
}

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <AppContent />
    </QueryClientProvider>
  )
}

export default App
```

### 3. Logger Selector Components

**web/src/components/LoggerSelector/LoggerSelector.jsx**:
```jsx
import React, { useState } from 'react'
import { ChevronDown, Info, Zap, Terminal, Settings, AlertCircle } from 'lucide-react'
import clsx from 'clsx'

const LOGGER_OPTIONS = [
  {
    value: 'slog',
    label: 'slog',
    description: 'Standard library structured logging',
    icon: Settings,
    performance: 'Good',
    useCase: 'Standard applications, built-in Go support',
    pros: ['Built-in Go 1.21+', 'Zero dependencies', 'Structured logging'],
    cons: ['Limited performance optimization', 'Newer, less ecosystem'],
    recommended: true,
    default: true
  },
  {
    value: 'zap',
    label: 'Zap',
    description: 'High-performance, zero-allocation logging',
    icon: Zap,
    performance: 'Excellent',
    useCase: 'High-throughput applications, performance-critical',
    pros: ['Zero allocation', 'High performance', 'Rich ecosystem'],
    cons: ['More complex API', 'Larger dependency'],
    recommended: true
  },
  {
    value: 'logrus',
    label: 'Logrus',
    description: 'Feature-rich, popular logging library',
    icon: Terminal,
    performance: 'Good',
    useCase: 'Traditional applications, extensive features',
    pros: ['Feature-rich', 'Large ecosystem', 'Hooks support'],
    cons: ['Slower performance', 'More allocations'],
    recommended: false
  },
  {
    value: 'zerolog',
    label: 'Zerolog',
    description: 'Zero allocation, chainable API',
    icon: AlertCircle,
    performance: 'Excellent',
    useCase: 'Modern applications, performance-focused',
    pros: ['Zero allocation', 'Chainable API', 'Small footprint'],
    cons: ['Different API style', 'Smaller ecosystem'],
    recommended: true
  }
]

const LoggerSelector = ({ value, onChange, showAdvanced = false }) => {
  const [isOpen, setIsOpen] = useState(false)
  const [showDetails, setShowDetails] = useState(false)
  
  const selectedLogger = LOGGER_OPTIONS.find(option => option.value === value) || LOGGER_OPTIONS[0]
  
  const handleSelect = (logger) => {
    onChange({
      type: logger.value,
      level: 'info',
      format: logger.value === 'slog' ? 'json' : 'json'
    })
    setIsOpen(false)
  }
  
  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <label className="block text-sm font-medium text-gray-700">
          Logger Library
        </label>
        <button
          onClick={() => setShowDetails(!showDetails)}
          className="text-xs text-blue-600 hover:text-blue-800"
        >
          {showDetails ? 'Hide' : 'Show'} Details
        </button>
      </div>
      
      {/* Logger Selection Dropdown */}
      <div className="relative">
        <button
          type="button"
          onClick={() => setIsOpen(!isOpen)}
          className="relative w-full cursor-pointer rounded-md border border-gray-300 bg-white py-2 pl-3 pr-10 text-left shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500 sm:text-sm"
        >
          <div className="flex items-center">
            <selectedLogger.icon className="h-5 w-5 text-gray-400 mr-3" />
            <div>
              <span className="font-medium">{selectedLogger.label}</span>
              <span className="ml-2 text-gray-500 text-sm">
                {selectedLogger.description}
              </span>
            </div>
          </div>
          <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
            <ChevronDown className="h-5 w-5 text-gray-400" />
          </span>
        </button>
        
        {isOpen && (
          <div className="absolute z-10 mt-1 w-full rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5">
            <ul className="max-h-60 overflow-auto rounded-md py-1 text-base focus:outline-none sm:text-sm">
              {LOGGER_OPTIONS.map((logger) => (
                <li key={logger.value}>
                  <button
                    onClick={() => handleSelect(logger)}
                    className={clsx(
                      'relative cursor-pointer select-none py-2 pl-3 pr-9 w-full text-left hover:bg-indigo-50',
                      value === logger.value ? 'bg-indigo-100 text-indigo-900' : 'text-gray-900'
                    )}
                  >
                    <div className="flex items-center">
                      <logger.icon className="h-5 w-5 text-gray-400 mr-3" />
                      <div>
                        <div className="flex items-center">
                          <span className="font-medium">{logger.label}</span>
                          {logger.recommended && (
                            <span className="ml-2 inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                              Recommended
                            </span>
                          )}
                        </div>
                        <span className="text-sm text-gray-500">{logger.description}</span>
                      </div>
                    </div>
                  </button>
                </li>
              ))}
            </ul>
          </div>
        )}
      </div>
      
      {/* Logger Details Panel */}
      {showDetails && (
        <div className="bg-gray-50 rounded-lg p-4 space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Performance</h4>
              <div className="flex items-center">
                <div className={clsx(
                  'px-2 py-1 rounded text-xs font-medium',
                  selectedLogger.performance === 'Excellent' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                )}>
                  {selectedLogger.performance}
                </div>
              </div>
            </div>
            
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Use Case</h4>
              <p className="text-sm text-gray-600">{selectedLogger.useCase}</p>
            </div>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Pros</h4>
              <ul className="text-sm text-gray-600 space-y-1">
                {selectedLogger.pros.map((pro, index) => (
                  <li key={index} className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    {pro}
                  </li>
                ))}
              </ul>
            </div>
            
            <div>
              <h4 className="font-medium text-gray-900 mb-2">Cons</h4>
              <ul className="text-sm text-gray-600 space-y-1">
                {selectedLogger.cons.map((con, index) => (
                  <li key={index} className="flex items-center">
                    <span className="text-red-500 mr-2">✗</span>
                    {con}
                  </li>
                ))}
              </ul>
            </div>
          </div>
        </div>
      )}
      
      {/* Advanced Logger Configuration */}
      {showAdvanced && (
        <div className="space-y-4 pt-4 border-t border-gray-200">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Log Level
              </label>
              <select
                value={value.level || 'info'}
                onChange={(e) => onChange({ ...value, level: e.target.value })}
                className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              >
                <option value="debug">Debug</option>
                <option value="info">Info</option>
                <option value="warn">Warn</option>
                <option value="error">Error</option>
              </select>
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                Output Format
              </label>
              <select
                value={value.format || 'json'}
                onChange={(e) => onChange({ ...value, format: e.target.value })}
                className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
              >
                <option value="json">JSON</option>
                <option value="text">Text</option>
                <option value="console">Console</option>
              </select>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}

export default LoggerSelector
```

### 4. Enhanced Configuration Panel

**web/src/components/ConfigPanel/ConfigPanel.jsx**:
```jsx
import React, { useState } from 'react'
import { useForm } from 'react-hook-form'
import { yupResolver } from '@hookform/resolvers/yup'
import * as yup from 'yup'
import { ChevronRight, ChevronDown, Settings, Zap, Database, Shield, Cloud, TestTube } from 'lucide-react'
import clsx from 'clsx'

import LoggerSelector from '../LoggerSelector/LoggerSelector'
import FrameworkSelector from '../FrameworkSelector/FrameworkSelector'
import DatabaseSelector from '../DatabaseSelector/DatabaseSelector'

const schema = yup.object({
  name: yup.string().required('Project name is required'),
  module: yup.string().required('Module path is required'),
  type: yup.string().required('Project type is required'),
  framework: yup.string().required('Framework is required'),
  logger: yup.object({
    type: yup.string().required('Logger type is required'),
    level: yup.string().required('Log level is required'),
    format: yup.string().required('Log format is required')
  })
})

const ConfigPanel = ({ 
  config, 
  onChange, 
  isAdvanced, 
  onToggleAdvanced, 
  onGenerate, 
  isGenerating, 
  error 
}) => {
  const [activeSection, setActiveSection] = useState('basic')
  
  const { register, handleSubmit, formState: { errors }, watch } = useForm({
    resolver: yupResolver(schema),
    defaultValues: config,
    mode: 'onChange'
  })
  
  const watchedValues = watch()
  
  React.useEffect(() => {
    onChange(watchedValues)
  }, [watchedValues, onChange])
  
  const sections = [
    { id: 'basic', label: 'Basic Settings', icon: Settings },
    { id: 'framework', label: 'Framework', icon: Zap },
    { id: 'logger', label: 'Logging', icon: Settings },
    { id: 'database', label: 'Database', icon: Database },
    { id: 'auth', label: 'Authentication', icon: Shield },
    { id: 'deployment', label: 'Deployment', icon: Cloud },
    { id: 'testing', label: 'Testing', icon: TestTube }
  ]
  
  const toggleSection = (sectionId) => {
    setActiveSection(activeSection === sectionId ? null : sectionId)
  }
  
  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold text-gray-900">Project Configuration</h2>
        <button
          onClick={onToggleAdvanced}
          className="text-sm text-blue-600 hover:text-blue-800"
        >
          {isAdvanced ? 'Basic Mode' : 'Advanced Mode'}
        </button>
      </div>
      
      {/* Error Display */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-md p-4">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}
      
      <form onSubmit={handleSubmit(onGenerate)} className="space-y-6">
        {/* Basic Settings */}
        <div className="space-y-4">
          <button
            type="button"
            onClick={() => toggleSection('basic')}
            className="flex items-center justify-between w-full py-2 text-left"
          >
            <div className="flex items-center">
              <Settings className="h-5 w-5 text-gray-400 mr-3" />
              <span className="font-medium text-gray-900">Basic Settings</span>
            </div>
            {activeSection === 'basic' ? (
              <ChevronDown className="h-5 w-5 text-gray-400" />
            ) : (
              <ChevronRight className="h-5 w-5 text-gray-400" />
            )}
          </button>
          
          {activeSection === 'basic' && (
            <div className="space-y-4 pl-8">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Project Name
                </label>
                <input
                  {...register('name')}
                  type="text"
                  className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  placeholder="my-awesome-project"
                />
                {errors.name && (
                  <p className="mt-1 text-sm text-red-600">{errors.name.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Module Path
                </label>
                <input
                  {...register('module')}
                  type="text"
                  className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  placeholder="github.com/username/my-awesome-project"
                />
                {errors.module && (
                  <p className="mt-1 text-sm text-red-600">{errors.module.message}</p>
                )}
              </div>
              
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  Project Type
                </label>
                <select
                  {...register('type')}
                  className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                >
                  <option value="api">Web API</option>
                  <option value="cli">CLI Application</option>
                  <option value="library">Library</option>
                  <option value="lambda">AWS Lambda</option>
                  <option value="microservice">Microservice</option>
                  <option value="monolith">Monolith</option>
                  <option value="workspace">Go Workspace</option>
                </select>
              </div>
              
              {config.type === 'api' && (
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Architecture Pattern
                  </label>
                  <select
                    {...register('architecture')}
                    className="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-indigo-500 focus:outline-none focus:ring-1 focus:ring-indigo-500"
                  >
                    <option value="standard">Standard</option>
                    <option value="clean">Clean Architecture</option>
                    <option value="ddd">Domain-Driven Design</option>
                    <option value="hexagonal">Hexagonal Architecture</option>
                  </select>
                </div>
              )}
            </div>
          )}
        </div>
        
        {/* Logger Configuration */}
        <div className="space-y-4">
          <button
            type="button"
            onClick={() => toggleSection('logger')}
            className="flex items-center justify-between w-full py-2 text-left"
          >
            <div className="flex items-center">
              <Settings className="h-5 w-5 text-gray-400 mr-3" />
              <span className="font-medium text-gray-900">Logging Configuration</span>
              <span className="ml-2 px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full">
                {config.logger?.type || 'slog'}
              </span>
            </div>
            {activeSection === 'logger' ? (
              <ChevronDown className="h-5 w-5 text-gray-400" />
            ) : (
              <ChevronRight className="h-5 w-5 text-gray-400" />
            )}
          </button>
          
          {activeSection === 'logger' && (
            <div className="pl-8">
              <LoggerSelector
                value={config.logger}
                onChange={(logger) => onChange({ ...config, logger })}
                showAdvanced={isAdvanced}
              />
            </div>
          )}
        </div>
        
        {/* Framework Configuration */}
        {config.type === 'api' && (
          <div className="space-y-4">
            <button
              type="button"
              onClick={() => toggleSection('framework')}
              className="flex items-center justify-between w-full py-2 text-left"
            >
              <div className="flex items-center">
                <Zap className="h-5 w-5 text-gray-400 mr-3" />
                <span className="font-medium text-gray-900">Framework</span>
                <span className="ml-2 px-2 py-1 bg-purple-100 text-purple-800 text-xs rounded-full">
                  {config.framework || 'gin'}
                </span>
              </div>
              {activeSection === 'framework' ? (
                <ChevronDown className="h-5 w-5 text-gray-400" />
              ) : (
                <ChevronRight className="h-5 w-5 text-gray-400" />
              )}
            </button>
            
            {activeSection === 'framework' && (
              <div className="pl-8">
                <FrameworkSelector
                  value={config.framework}
                  onChange={(framework) => onChange({ ...config, framework })}
                  projectType={config.type}
                />
              </div>
            )}
          </div>
        )}
        
        {/* Generate Button */}
        <div className="pt-6 border-t border-gray-200">
          <button
            type="submit"
            disabled={isGenerating}
            className={clsx(
              'w-full px-4 py-2 text-sm font-medium text-white rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500',
              isGenerating
                ? 'bg-gray-400 cursor-not-allowed'
                : 'bg-indigo-600 hover:bg-indigo-700'
            )}
          >
            {isGenerating ? 'Generating...' : 'Generate Project'}
          </button>
        </div>
      </form>
    </div>
  )
}

export default ConfigPanel
```

### 5. Enhanced Backend API with Logger Support

**internal/api/handlers/generate.go** (updates):
```go
type GenerateRequest struct {
    Config types.ProjectConfig `json:"config" binding:"required"`
}

type ProjectConfig struct {
    Name         string            `json:"name" binding:"required"`
    Module       string            `json:"module" binding:"required"`
    Type         string            `json:"type" binding:"required"`
    GoVersion    string            `json:"go_version"`
    Architecture string            `json:"architecture"`
    Framework    string            `json:"framework"`
    Logger       LoggerConfig      `json:"logger"`
    Features     Features          `json:"features"`
    Variables    map[string]string `json:"variables"`
}

type LoggerConfig struct {
    Type   string `json:"type" binding:"required"`
    Level  string `json:"level"`
    Format string `json:"format"`
}

type Features struct {
    Database       DatabaseConfig       `json:"database"`
    Authentication AuthenticationConfig `json:"authentication"`
    Deployment     DeploymentConfig     `json:"deployment"`
    Testing        TestingConfig        `json:"testing"`
}

func (h *GenerateHandler) Generate(c *gin.Context) {
    var req GenerateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Validate logger configuration
    if err := h.validateLoggerConfig(req.Config.Logger); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Set logger defaults if not provided
    if req.Config.Logger.Type == "" {
        req.Config.Logger.Type = "slog"
    }
    if req.Config.Logger.Level == "" {
        req.Config.Logger.Level = "info"
    }
    if req.Config.Logger.Format == "" {
        req.Config.Logger.Format = "json"
    }
    
    // Generate project with logger configuration
    projectID := uuid.New().String()
    files, err := h.generator.GenerateInMemory(req.Config)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Store project info with logger metadata
    projectInfo := &ProjectInfo{
        ID:      projectID,
        Config:  req.Config,
        Files:   files,
        Created: time.Now(),
        Metadata: map[string]interface{}{
            "logger_type":   req.Config.Logger.Type,
            "logger_level":  req.Config.Logger.Level,
            "logger_format": req.Config.Logger.Format,
        },
    }
    h.projects[projectID] = projectInfo
    
    // Build response with logger information
    response := GenerateResponse{
        ProjectID: projectID,
        Files:     h.convertFilesToStrings(files),
        Structure: h.buildFileStructure(files),
        Metadata: map[string]interface{}{
            "template_used":  h.getTemplateID(req.Config),
            "file_count":     len(files),
            "generated":      time.Now(),
            "logger_config":  req.Config.Logger,
            "logger_files":   h.getLoggerFiles(files, req.Config.Logger.Type),
        },
    }
    
    c.JSON(http.StatusOK, response)
}

func (h *GenerateHandler) validateLoggerConfig(config LoggerConfig) error {
    validLoggers := []string{"slog", "zap", "logrus", "zerolog"}
    validLevels := []string{"debug", "info", "warn", "error"}
    validFormats := []string{"json", "text", "console"}
    
    if !contains(validLoggers, config.Type) {
        return fmt.Errorf("invalid logger type: %s. Must be one of: %v", config.Type, validLoggers)
    }
    
    if config.Level != "" && !contains(validLevels, config.Level) {
        return fmt.Errorf("invalid log level: %s. Must be one of: %v", config.Level, validLevels)
    }
    
    if config.Format != "" && !contains(validFormats, config.Format) {
        return fmt.Errorf("invalid log format: %s. Must be one of: %v", config.Format, validFormats)
    }
    
    return nil
}

func (h *GenerateHandler) getLoggerFiles(files map[string][]byte, loggerType string) []string {
    loggerFiles := []string{}
    
    for path := range files {
        if strings.Contains(path, "logger/") {
            loggerFiles = append(loggerFiles, path)
        }
    }
    
    return loggerFiles
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
```

### 6. Real-time Logger Preview Updates

**internal/api/handlers/websocket.go** (updates):
```go
func (h *PreviewHandler) handlePreviewRequest(conn *websocket.Conn, config types.ProjectConfig) {
    // Generate real-time preview with logger-specific files
    structure := h.generatePreviewStructure(config)
    
    // Include logger-specific information
    loggerInfo := map[string]interface{}{
        "type":   config.Logger.Type,
        "level":  config.Logger.Level,
        "format": config.Logger.Format,
        "files":  h.getLoggerFilesList(config.Logger.Type),
    }
    
    response := PreviewMessage{
        Type: "preview_result",
        Data: map[string]interface{}{
            "structure":     structure,
            "fileCount":     len(structure),
            "loggerInfo":    loggerInfo,
            "dependencies":  h.getLoggerDependencies(config.Logger.Type),
        },
    }
    
    conn.WriteJSON(response)
}

func (h *PreviewHandler) getLoggerFilesList(loggerType string) []string {
    baseFiles := []string{
        "internal/logger/interface.go",
        "internal/logger/factory.go",
    }
    
    switch loggerType {
    case "slog":
        return append(baseFiles, "internal/logger/slog.go")
    case "zap":
        return append(baseFiles, "internal/logger/zap.go")
    case "logrus":
        return append(baseFiles, "internal/logger/logrus.go")
    case "zerolog":
        return append(baseFiles, "internal/logger/zerolog.go")
    default:
        return baseFiles
    }
}

func (h *PreviewHandler) getLoggerDependencies(loggerType string) []string {
    switch loggerType {
    case "zap":
        return []string{"go.uber.org/zap"}
    case "logrus":
        return []string{"github.com/sirupsen/logrus"}
    case "zerolog":
        return []string{"github.com/rs/zerolog"}
    case "slog":
        return []string{} // Built-in, no external dependencies
    default:
        return []string{}
    }
}
```

## Logger Selector Integration Testing

### Unit Tests for Logger Components

**web/src/components/LoggerSelector/__tests__/LoggerSelector.test.jsx**:
```jsx
import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import LoggerSelector from '../LoggerSelector'

describe('LoggerSelector', () => {
  const mockOnChange = jest.fn()
  
  beforeEach(() => {
    mockOnChange.mockClear()
  })
  
  test('renders with default slog selection', () => {
    render(
      <LoggerSelector 
        value={{ type: 'slog', level: 'info', format: 'json' }}
        onChange={mockOnChange}
      />
    )
    
    expect(screen.getByText('slog')).toBeInTheDocument()
    expect(screen.getByText('Standard library structured logging')).toBeInTheDocument()
  })
  
  test('shows all logger options when opened', () => {
    render(
      <LoggerSelector 
        value={{ type: 'slog', level: 'info', format: 'json' }}
        onChange={mockOnChange}
      />
    )
    
    fireEvent.click(screen.getByRole('button'))
    
    expect(screen.getByText('Zap')).toBeInTheDocument()
    expect(screen.getByText('Logrus')).toBeInTheDocument()
    expect(screen.getByText('Zerolog')).toBeInTheDocument()
  })
  
  test('calls onChange when logger is selected', () => {
    render(
      <LoggerSelector 
        value={{ type: 'slog', level: 'info', format: 'json' }}
        onChange={mockOnChange}
      />
    )
    
    fireEvent.click(screen.getByRole('button'))
    fireEvent.click(screen.getByText('Zap'))
    
    expect(mockOnChange).toHaveBeenCalledWith({
      type: 'zap',
      level: 'info',
      format: 'json'
    })
  })
  
  test('shows advanced configuration when enabled', () => {
    render(
      <LoggerSelector 
        value={{ type: 'slog', level: 'info', format: 'json' }}
        onChange={mockOnChange}
        showAdvanced={true}
      />
    )
    
    expect(screen.getByText('Log Level')).toBeInTheDocument()
    expect(screen.getByText('Output Format')).toBeInTheDocument()
  })
})
```

### Integration Tests

**web/src/__tests__/LoggerIntegration.test.jsx**:
```jsx
import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { QueryClient, QueryClientProvider } from 'react-query'
import App from '../App'

// Mock WebSocket
global.WebSocket = jest.fn(() => ({
  send: jest.fn(),
  close: jest.fn(),
  addEventListener: jest.fn(),
}))

describe('Logger Integration', () => {
  let queryClient
  
  beforeEach(() => {
    queryClient = new QueryClient({
      defaultOptions: {
        queries: { retry: false },
        mutations: { retry: false },
      },
    })
  })
  
  test('generates project with selected logger', async () => {
    const mockGenerate = jest.fn().mockResolvedValue({
      project_id: 'test-123',
      files: {},
      structure: [],
      metadata: {
        logger_config: {
          type: 'zap',
          level: 'info',
          format: 'json'
        }
      }
    })
    
    // Mock the generate API call
    global.fetch = jest.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve(mockGenerate())
    })
    
    render(
      <QueryClientProvider client={queryClient}>
        <App />
      </QueryClientProvider>
    )
    
    // Configure project with zap logger
    fireEvent.click(screen.getByText('Logging Configuration'))
    fireEvent.click(screen.getByRole('button', { name: /slog/ }))
    fireEvent.click(screen.getByText('Zap'))
    
    // Generate project
    fireEvent.click(screen.getByText('Generate Project'))
    
    await waitFor(() => {
      expect(global.fetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/generate'),
        expect.objectContaining({
          method: 'POST',
          body: expect.stringContaining('"type":"zap"')
        })
      )
    })
  })
})
```

## Deployment Enhancements

### Updated Package.json Dependencies

```json
{
  "dependencies": {
    "react": "^18.2.0",
    "react-dom": "^18.2.0",
    "axios": "^1.6.0",
    "react-hook-form": "^7.48.0",
    "@hookform/resolvers": "^3.3.0",
    "yup": "^1.3.0",
    "react-query": "^3.39.0",
    "react-syntax-highlighter": "^15.5.0",
    "react-select": "^5.8.0",
    "react-toggle": "^4.1.0",
    "lucide-react": "^0.263.1",
    "clsx": "^2.0.0",
    "@testing-library/react": "^13.4.0",
    "@testing-library/jest-dom": "^5.16.5",
    "@testing-library/user-event": "^14.4.3"
  }
}
```

## Success Criteria

### Phase 3 Logger Integration Success Metrics

- ✅ **Logger Selector Component**: Interactive UI component with 4 logger options
- ✅ **Visual Logger Information**: Icons, descriptions, pros/cons, performance indicators
- ✅ **Advanced Configuration**: Log level and format selection in advanced mode
- ✅ **Real-time Preview**: WebSocket updates when logger selection changes
- ✅ **Backend API Integration**: Proper logger validation and configuration handling
- ✅ **Logger File Preview**: Show logger-specific files in preview panel
- ✅ **Dependency Tracking**: Display logger-specific dependencies
- ✅ **Testing Coverage**: Unit and integration tests for logger components
- ✅ **Progressive Disclosure**: Logger details shown/hidden based on user preference
- ✅ **Error Handling**: Proper validation and error messages for logger configuration

### Technical Implementation Checklist

- [x] Logger selector component with visual icons
- [x] Advanced logger configuration panel
- [x] Real-time preview updates with logger info
- [x] Backend API logger validation
- [x] Logger-specific file generation preview
- [x] Dependency management for different loggers
- [x] Comprehensive testing suite
- [x] Integration with existing configuration system
- [x] WebSocket support for real-time logger updates
- [x] Error handling and user feedback