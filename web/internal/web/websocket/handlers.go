package websocket

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/francknouama/go-starter/internal/generator"
	"github.com/francknouama/go-starter/internal/templates"
	"github.com/francknouama/go-starter/pkg/types"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

// WSHandler provides WebSocket-related HTTP handlers
type WSHandler struct {
	hub *Hub
}

// NewWSHandler creates a new WebSocket handler
func NewWSHandler(hub *Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// HandleWebSocketUpgrade handles WebSocket upgrade requests
func (h *WSHandler) HandleWebSocketUpgrade(c *gin.Context) {
	HandleWebSocket(h.hub, c.Writer, c.Request)
}

// GetConnectedClients returns information about connected clients
func (h *WSHandler) GetConnectedClients(c *gin.Context) {
	clients := make([]map[string]interface{}, 0)
	
	h.hub.mutex.RLock()
	for client := range h.hub.clients {
		clients = append(clients, client.Info())
	}
	h.hub.mutex.RUnlock()
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"clients": clients,
		"count":   len(clients),
	})
	
	c.JSON(http.StatusOK, response)
}

// BroadcastMessage broadcasts a message to all connected clients
func (h *WSHandler) BroadcastMessage(c *gin.Context) {
	var request struct {
		Type    string      `json:"type" binding:"required"`
		Data    interface{} `json:"data"`
		Channel string      `json:"channel,omitempty"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	message := models.WebSocketMessage{
		Type:      request.Type,
		Data:      request.Data,
		Timestamp: time.Now(),
	}
	
	h.hub.Broadcast(message)
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"message": "Message broadcasted successfully",
		"clients": h.hub.GetClientCount(),
	})
	
	c.JSON(http.StatusOK, response)
}

// SendMessageToClient sends a message to a specific client
func (h *WSHandler) SendMessageToClient(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		response := models.NewValidationErrorResponse("Client ID is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	var request struct {
		Type string      `json:"type" binding:"required"`
		Data interface{} `json:"data"`
	}
	
	if err := c.ShouldBindJSON(&request); err != nil {
		response := models.NewValidationErrorResponse(err.Error())
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	message := models.WebSocketMessage{
		Type:      request.Type,
		Data:      request.Data,
		Timestamp: time.Now(),
	}
	
	// Find and send to specific client
	h.hub.mutex.RLock()
	var found bool
	for client := range h.hub.clients {
		if client.ID == clientID {
			h.hub.SendToClient(client, message)
			found = true
			break
		}
	}
	h.hub.mutex.RUnlock()
	
	if !found {
		response := models.NewNotFoundErrorResponse("Client not found")
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"message":  "Message sent successfully",
		"clientId": clientID,
	})
	
	c.JSON(http.StatusOK, response)
}

// GetHubStats returns WebSocket hub statistics
func (h *WSHandler) GetHubStats(c *gin.Context) {
	stats := h.hub.Stats()
	response := models.NewSuccessResponse(stats)
	c.JSON(http.StatusOK, response)
}

// NotifyProgress sends a progress notification to clients
func (h *WSHandler) NotifyProgress(requestID string, progress models.ProgressData, clientID ...string) {
	h.hub.SendProgress(requestID, progress, clientID...)
}

// NotifyError sends an error notification to clients
func (h *WSHandler) NotifyError(requestID string, errorMsg string, clientID ...string) {
	h.hub.SendError(requestID, errorMsg, clientID...)
}

// NotifyComplete sends a completion notification to clients
func (h *WSHandler) NotifyComplete(requestID string, data interface{}, clientID ...string) {
	h.hub.SendComplete(requestID, data, clientID...)
}

// HandleSystemEvent handles system-wide events and broadcasts them
func (h *WSHandler) HandleSystemEvent(eventType string, data interface{}) {
	message := models.WebSocketMessage{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
	}
	
	h.hub.Broadcast(message)
}

// CreateProgressChannel creates a progress channel for real-time updates
func (h *WSHandler) CreateProgressChannel(requestID string) chan models.ProgressData {
	progressChan := make(chan models.ProgressData, 100)
	
	go func() {
		defer close(progressChan)
		
		for progress := range progressChan {
			h.hub.SendProgress(requestID, progress)
		}
	}()
	
	return progressChan
}

// HealthCheck returns the health status of the WebSocket service
func (h *WSHandler) HealthCheck(c *gin.Context) {
	health := map[string]interface{}{
		"service":         "websocket",
		"status":          "healthy",
		"connectedClients": h.hub.GetClientCount(),
		"timestamp":       time.Now(),
	}
	
	response := models.NewSuccessResponse(health)
	c.JSON(http.StatusOK, response)
}

// ServeWebSocketInfo serves information about WebSocket endpoints
func (h *WSHandler) ServeWebSocketInfo(c *gin.Context) {
	info := map[string]interface{}{
		"websocketEndpoint": "/api/v1/ws",
		"supportedEvents": []string{
			models.WSMessageTypeProgress,
			models.WSMessageTypeComplete,
			models.WSMessageTypeError,
			models.WSMessageTypePreview,
			models.WSMessageTypeStatus,
			"file_tree",
			"file_content",
			"preview_start",
			"preview_complete",
		},
		"clientCommands": []string{
			"ping",
			"subscribe",
			"unsubscribe",
			"heartbeat",
			"preview_request",
			"file_request",
		},
		"connectedClients": h.hub.GetClientCount(),
		"documentation":    "Connect to the WebSocket endpoint to receive real-time updates",
	}
	
	response := models.NewSuccessResponse(info)
	c.JSON(http.StatusOK, response)
}

// DisconnectClient forcefully disconnects a specific client
func (h *WSHandler) DisconnectClient(c *gin.Context) {
	clientID := c.Param("clientId")
	if clientID == "" {
		response := models.NewValidationErrorResponse("Client ID is required")
		c.JSON(http.StatusBadRequest, response)
		return
	}
	
	h.hub.mutex.Lock()
	var found bool
	for client := range h.hub.clients {
		if client.ID == clientID {
			delete(h.hub.clients, client)
			close(client.send)
			client.conn.Close()
			found = true
			break
		}
	}
	h.hub.mutex.Unlock()
	
	if !found {
		response := models.NewNotFoundErrorResponse("Client not found")
		c.JSON(http.StatusNotFound, response)
		return
	}
	
	response := models.NewSuccessResponse(map[string]interface{}{
		"message":  "Client disconnected successfully",
		"clientId": clientID,
	})
	
	c.JSON(http.StatusOK, response)
}

// HandlePreviewRequest handles real-time preview generation requests
func (h *WSHandler) HandlePreviewRequest(requestID string, request models.PreviewRequest, clientID ...string) {
	// Send initial status
	h.NotifyStatus(requestID, "Initializing preview generation...", clientID...)
	
	// Start preview generation in goroutine
	go func() {
		// Validate request
		if err := request.Validate(); err != nil {
			h.NotifyError(requestID, fmt.Sprintf("Validation failed: %v", err), clientID...)
			return
		}
		
		// Set defaults
		request.SetDefaults()
		
		// Convert to internal config
		config := convertToProjectConfig(request.GenerationRequest)
		
		h.NotifyStatus(requestID, "Loading blueprint...", clientID...)
		
		// Initialize generator and template registry
		gen := generator.New()
		registry := templates.NewRegistry()
		
		// Determine blueprint ID
		blueprintID := determineBlueprintID(config)
		
		// Get template
		tmpl, err := registry.Get(blueprintID)
		if err != nil {
			// Try fallback: look for templates by type
			templatesByType := registry.GetByType(config.Type)
			if len(templatesByType) > 0 {
				tmpl = templatesByType[0]
			} else {
				h.NotifyError(requestID, fmt.Sprintf("Blueprint '%s' not found", blueprintID), clientID...)
				return
			}
		}
		
		h.NotifyStatus(requestID, "Generating file structure...", clientID...)
		
		// Generate files in memory
		filesMap, err := gen.GenerateInMemory(&config, tmpl.ID)
		if err != nil {
			h.NotifyError(requestID, fmt.Sprintf("Failed to generate preview: %v", err), clientID...)
			return
		}
		
		// Convert to preview format and send file tree
		var files []models.GeneratedFile
		totalFiles := len(filesMap)
		processed := 0
		
		// Track directories for file tree
		dirs := make(map[string]bool)
		
		for filePath, content := range filesMap {
			// Add parent directories
			dir := filepath.Dir(filePath)
			for dir != "." && dir != "/" {
				dirs[dir] = true
				dir = filepath.Dir(dir)
			}
			
			// Create file info
			file := models.GeneratedFile{
				Path:    filePath,
				Content: string(content),
				Size:    int64(len(content)),
				IsDir:   false,
				Mode:    "0644",
				ModTime: time.Now().Format(time.RFC3339),
			}
			
			files = append(files, file)
			processed++
			
			// Send progress update
			progressData := models.ProgressData{
				Stage:           "generating",
				Progress:        float64(processed) / float64(totalFiles),
				Message:         fmt.Sprintf("Processing %s", filepath.Base(filePath)),
				CurrentFile:     filePath,
				TotalFiles:      totalFiles,
				ProcessedFiles:  processed,
			}
			h.NotifyProgress(requestID, progressData, clientID...)
			
			// Send individual file content
			h.SendFileContent(requestID, file, clientID...)
		}
		
		// Add directory entries
		for dir := range dirs {
			files = append(files, models.GeneratedFile{
				Path:    dir,
				Content: "",
				Size:    0,
				IsDir:   true,
				Mode:    "0755",
				ModTime: time.Now().Format(time.RFC3339),
			})
		}
		
		// Build and send file tree
		fileTree := buildFileTree(files)
		h.SendFileTree(requestID, fileTree, clientID...)
		
		// Send completion notification
		previewData := map[string]interface{}{
			"projectName":  config.Name,
			"modulePath":   config.Module,
			"type":         config.Type,
			"architecture": config.Architecture,
			"framework":    config.Framework,
			"totalFiles":   len(files),
			"blueprintUsed": tmpl.ID,
		}
		
		h.NotifyComplete(requestID, previewData, clientID...)
	}()
}

// SendFileTree sends file tree structure via WebSocket
func (h *WSHandler) SendFileTree(requestID string, fileTree *models.FileTreeNode, clientID ...string) {
	message := models.WebSocketMessage{
		Type:      "file_tree",
		Data:      fileTree,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		h.sendToSpecificClient(clientID[0], message)
	} else {
		h.hub.Broadcast(message)
	}
}

// SendFileContent sends individual file content via WebSocket
func (h *WSHandler) SendFileContent(requestID string, file models.GeneratedFile, clientID ...string) {
	message := models.WebSocketMessage{
		Type:      "file_content",
		Data:      file,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		h.sendToSpecificClient(clientID[0], message)
	} else {
		h.hub.Broadcast(message)
	}
}

// NotifyStatus sends a status notification to clients
func (h *WSHandler) NotifyStatus(requestID string, message string, clientID ...string) {
	statusData := map[string]interface{}{
		"message": message,
		"timestamp": time.Now().Format(time.RFC3339),
	}
	
	wsMessage := models.WebSocketMessage{
		Type:      models.WSMessageTypeStatus,
		Data:      statusData,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		h.sendToSpecificClient(clientID[0], wsMessage)
	} else {
		h.hub.Broadcast(wsMessage)
	}
}

// sendToSpecificClient sends a message to a specific client
func (h *WSHandler) sendToSpecificClient(clientID string, message models.WebSocketMessage) {
	h.hub.mutex.RLock()
	defer h.hub.mutex.RUnlock()
	
	for client := range h.hub.clients {
		if client.ID == clientID {
			h.hub.SendToClient(client, message)
			break
		}
	}
}

// WebSocketTestPage serves a simple test page for WebSocket connections
func (h *WSHandler) WebSocketTestPage(c *gin.Context) {
	html := `
<!DOCTYPE html>
<html>
<head>
    <title>WebSocket Test</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        #messages { border: 1px solid #ccc; height: 300px; overflow-y: scroll; padding: 10px; margin: 10px 0; }
        input[type="text"] { width: 300px; padding: 5px; }
        button { padding: 5px 10px; margin: 5px; }
    </style>
</head>
<body>
    <h1>WebSocket Test Page</h1>
    <div>
        <button onclick="connect()">Connect</button>
        <button onclick="disconnect()">Disconnect</button>
        <span id="status">Disconnected</span>
    </div>
    
    <div>
        <input type="text" id="messageInput" placeholder="Enter message">
        <button onclick="sendMessage()">Send</button>
    </div>
    
    <div id="messages"></div>
    
    <script>
        let ws = null;
        const messages = document.getElementById('messages');
        const status = document.getElementById('status');
        
        function connect() {
            const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
            const wsUrl = protocol + '//' + window.location.host + '/api/v1/ws';
            
            ws = new WebSocket(wsUrl);
            
            ws.onopen = function() {
                status.textContent = 'Connected';
                addMessage('Connected to WebSocket');
            };
            
            ws.onmessage = function(event) {
                addMessage('Received: ' + event.data);
            };
            
            ws.onclose = function() {
                status.textContent = 'Disconnected';
                addMessage('WebSocket closed');
            };
            
            ws.onerror = function(error) {
                addMessage('Error: ' + error);
            };
        }
        
        function disconnect() {
            if (ws) {
                ws.close();
            }
        }
        
        function sendMessage() {
            const input = document.getElementById('messageInput');
            if (ws && input.value) {
                const message = {
                    type: 'ping',
                    data: input.value,
                    timestamp: new Date().toISOString()
                };
                ws.send(JSON.stringify(message));
                addMessage('Sent: ' + input.value);
                input.value = '';
            }
        }
        
        function addMessage(message) {
            const div = document.createElement('div');
            div.textContent = new Date().toLocaleTimeString() + ' - ' + message;
            messages.appendChild(div);
            messages.scrollTop = messages.scrollHeight;
        }
        
        // Auto-connect on page load
        connect();
    </script>
</body>
</html>
    `
	
	c.Header("Content-Type", "text/html")
	c.String(http.StatusOK, html)
}

// Helper functions for preview generation

// convertToProjectConfig converts a GenerationRequest to internal ProjectConfig
func convertToProjectConfig(req models.GenerationRequest) types.ProjectConfig {
	config := types.ProjectConfig{
		Name:         req.Name,
		Module:       req.ModulePath,
		Type:         req.Type,
		Architecture: req.Architecture,
		Framework:    req.Framework,
		GoVersion:    req.GoVersion,
		Logger:       req.Logger,
		Variables:    make(map[string]string),
	}

	// Add complexity to variables for blueprint selection
	if req.Complexity != "" {
		config.Variables["complexity"] = req.Complexity
	}

	// Add advanced flag
	if req.Advanced {
		config.Variables["advanced"] = "true"
	}

	// Convert features if present
	if req.Features != nil {
		config.Features = &types.Features{}

		if req.Features.Database != nil {
			config.Features.Database = types.DatabaseConfig{
				Driver: req.Features.Database.Driver,
				ORM:    req.Features.Database.ORM,
			}
		}

		if req.Features.Authentication != nil {
			config.Features.Authentication = types.AuthConfig{
				Type:      req.Features.Authentication.Type,
				Providers: req.Features.Authentication.Providers,
			}
		}

		if req.Features.Deployment != nil {
			config.Features.Deployment = types.DeployConfig{
				Targets: req.Features.Deployment.Targets,
			}
		}

		if req.Features.Testing != nil {
			config.Features.Testing = types.TestConfig{
				Framework: req.Features.Testing.Framework,
				Coverage:  req.Features.Testing.Coverage,
			}
		}

		if req.Features.Logging != nil {
			config.Features.Logging = types.LoggingConfig{
				Type:       req.Logger, // Use logger from request
				Level:      req.Features.Logging.Level,
				Format:     req.Features.Logging.Format,
				Structured: true, // Default to structured logging
			}
		}

		if req.Features.Monitoring != nil {
			config.Features.Monitoring = types.MonitorConfig{
				Metrics: req.Features.Monitoring.Metrics,
				Tracing: req.Features.Monitoring.Tracing,
			}
		}
	}

	return config
}

// determineBlueprintID determines the blueprint ID based on the configuration
func determineBlueprintID(config types.ProjectConfig) string {
	// Check for complexity-based blueprint selection (e.g., cli-simple)
	if complexity, exists := config.Variables["complexity"]; exists && complexity == "simple" && config.Type == "cli" {
		return "cli-simple"
	}

	// Use architecture-based selection for other types
	if config.Architecture != "" && config.Architecture != "standard" {
		return fmt.Sprintf("%s-%s", config.Type, config.Architecture)
	}

	return config.Type
}

// buildFileTree builds a hierarchical file tree structure
func buildFileTree(files []models.GeneratedFile) *models.FileTreeNode {
	root := &models.FileTreeNode{
		Name:     "root",
		Path:     "",
		IsDir:    true,
		Children: []*models.FileTreeNode{},
	}

	// Build tree structure
	for _, file := range files {
		if file.Path == "" {
			continue
		}

		parts := strings.Split(filepath.Clean(file.Path), string(filepath.Separator))
		current := root

		// Navigate/create path to file
		currentPath := ""
		for i, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			// Look for existing child
			var found *models.FileTreeNode
			for _, child := range current.Children {
				if child.Name == part {
					found = child
					break
				}
			}

			if found == nil {
				// Create new node
				isDir := i < len(parts)-1 || file.IsDir
				newNode := &models.FileTreeNode{
					Name:     part,
					Path:     currentPath,
					IsDir:    isDir,
					Children: []*models.FileTreeNode{},
				}

				if !isDir {
					newNode.Size = file.Size
				}

				current.Children = append(current.Children, newNode)
				current = newNode
			} else {
				current = found
			}
		}
	}

	return root
}