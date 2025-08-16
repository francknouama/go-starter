package websocket

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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
		},
		"clientCommands": []string{
			"ping",
			"subscribe",
			"unsubscribe",
			"heartbeat",
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