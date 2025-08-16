package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/francknouama/go-starter/web/internal/web/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from any origin in development
		// In production, you should implement proper origin checking
		return true
	},
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for thread-safe operations
	mutex sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client registration/unregistration and broadcasting
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mutex.Lock()
			h.clients[client] = true
			h.mutex.Unlock()
			
			log.Printf("Client connected: %s (Total clients: %d)", client.ID, len(h.clients))
			
			// Send welcome message
			welcomeMsg := models.WebSocketMessage{
				Type:      "welcome",
				Data:      map[string]interface{}{"clientId": client.ID},
				Timestamp: time.Now(),
			}
			h.SendToClient(client, welcomeMsg)

		case client := <-h.unregister:
			h.mutex.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.mutex.Unlock()
				log.Printf("Client disconnected: %s (Total clients: %d)", client.ID, len(h.clients))
			} else {
				h.mutex.Unlock()
			}

		case message := <-h.broadcast:
			h.mutex.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					// Client's send channel is full, close it
					delete(h.clients, client)
					close(client.send)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// RegisterClient adds a client to the hub
func (h *Hub) RegisterClient(client *Client) {
	h.register <- client
}

// UnregisterClient removes a client from the hub
func (h *Hub) UnregisterClient(client *Client) {
	h.unregister <- client
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(message models.WebSocketMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}
	
	h.broadcast <- data
}

// SendToClient sends a message to a specific client
func (h *Hub) SendToClient(client *Client, message models.WebSocketMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling WebSocket message: %v", err)
		return
	}
	
	select {
	case client.send <- data:
	default:
		// Client's send channel is full, close it
		h.mutex.Lock()
		delete(h.clients, client)
		close(client.send)
		h.mutex.Unlock()
	}
}

// SendProgress sends a progress update to all clients or a specific client
func (h *Hub) SendProgress(requestID string, progress models.ProgressData, clientID ...string) {
	message := models.WebSocketMessage{
		Type:      models.WSMessageTypeProgress,
		Data:      progress,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		// Send to specific client
		h.mutex.RLock()
		for client := range h.clients {
			if client.ID == clientID[0] {
				h.SendToClient(client, message)
				break
			}
		}
		h.mutex.RUnlock()
	} else {
		// Broadcast to all clients
		h.Broadcast(message)
	}
}

// SendError sends an error message to all clients or a specific client
func (h *Hub) SendError(requestID string, errorMsg string, clientID ...string) {
	message := models.WebSocketMessage{
		Type:      models.WSMessageTypeError,
		Data:      map[string]interface{}{"error": errorMsg},
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		// Send to specific client
		h.mutex.RLock()
		for client := range h.clients {
			if client.ID == clientID[0] {
				h.SendToClient(client, message)
				break
			}
		}
		h.mutex.RUnlock()
	} else {
		// Broadcast to all clients
		h.Broadcast(message)
	}
}

// SendComplete sends a completion message to all clients or a specific client
func (h *Hub) SendComplete(requestID string, data interface{}, clientID ...string) {
	message := models.WebSocketMessage{
		Type:      models.WSMessageTypeComplete,
		Data:      data,
		Timestamp: time.Now(),
		RequestID: requestID,
	}
	
	if len(clientID) > 0 && clientID[0] != "" {
		// Send to specific client
		h.mutex.RLock()
		for client := range h.clients {
			if client.ID == clientID[0] {
				h.SendToClient(client, message)
				break
			}
		}
		h.mutex.RUnlock()
	} else {
		// Broadcast to all clients
		h.Broadcast(message)
	}
}

// GetClientCount returns the number of connected clients
func (h *Hub) GetClientCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetClientIDs returns a list of connected client IDs
func (h *Hub) GetClientIDs() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	ids := make([]string, 0, len(h.clients))
	for client := range h.clients {
		ids = append(ids, client.ID)
	}
	return ids
}

// HandleWebSocket handles WebSocket upgrade and client management
func HandleWebSocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	
	client := NewClient(hub, conn)
	hub.RegisterClient(client)
	
	// Start client goroutines
	go client.writePump()
	go client.readPump()
}

// BroadcastSystemStatus sends system status to all clients
func (h *Hub) BroadcastSystemStatus(status map[string]interface{}) {
	message := models.WebSocketMessage{
		Type:      models.WSMessageTypeStatus,
		Data:      status,
		Timestamp: time.Now(),
	}
	
	h.Broadcast(message)
}

// Stats returns hub statistics
func (h *Hub) Stats() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	
	return map[string]interface{}{
		"connectedClients": len(h.clients),
		"clientIds":        h.GetClientIDs(),
		"uptime":           time.Since(time.Now()).String(), // This would be set at hub creation
	}
}