package websocket

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// Allow connections from localhost (development)
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:5173" || origin == "http://localhost:3000"
	},
}

// Client represents a WebSocket client
type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
	Hub  *Hub
	Type string // "generate" or "preview"
}

// Hub maintains active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from clients
	broadcast chan []byte

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			slog.Info("WebSocket client connected", "client_id", client.ID, "type", client.Type)

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send)
				slog.Info("WebSocket client disconnected", "client_id", client.ID, "type", client.Type)
			}

		case message := <-h.broadcast:
			// Broadcast to all clients
			for client := range h.clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.clients, client)
				}
			}
		}
	}
}

// BroadcastToType sends a message to all clients of a specific type
func (h *Hub) BroadcastToType(messageType string, data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal WebSocket message", "error", err)
		return
	}

	for client := range h.clients {
		if client.Type == messageType {
			select {
			case client.Send <- message:
			default:
				close(client.Send)
				delete(h.clients, client)
			}
		}
	}
}

// Broadcast sends a message to all connected clients
func (h *Hub) Broadcast(data interface{}) {
	message, err := json.Marshal(data)
	if err != nil {
		slog.Error("Failed to marshal WebSocket message", "error", err)
		return
	}

	select {
	case h.broadcast <- message:
	default:
		slog.Warn("Broadcast channel is full")
	}
}

// UpgradeConnection upgrades an HTTP connection to WebSocket
func (h *Hub) UpgradeConnection(c *gin.Context, clientType string) (*Client, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	// Generate client ID
	clientID := generateClientID()

	client := &Client{
		ID:   clientID,
		Conn: conn,
		Send: make(chan []byte, 256),
		Hub:  h,
		Type: clientType,
	}

	// Register client
	h.register <- client

	// Start goroutines for reading and writing
	go client.readPump()
	go client.writePump()

	return client, nil
}

// readPump handles reading from the WebSocket connection
func (c *Client) readPump() {
	defer func() {
		c.Hub.unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				slog.Error("WebSocket read error", "error", err, "client_id", c.ID)
			}
			break
		}

		// Handle incoming messages (for future interactive features)
		slog.Debug("Received WebSocket message", "client_id", c.ID, "message", string(message))
	}
}

// writePump handles writing to the WebSocket connection
func (c *Client) writePump() {
	defer c.Conn.Close()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				slog.Error("WebSocket write error", "error", err, "client_id", c.ID)
				return
			}
		}
	}
}

// generateClientID generates a unique client identifier
func generateClientID() string {
	// Simple implementation - in production, use UUID or similar
	return "client_" + randString(8)
}

// randString generates a random string of specified length
func randString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[len(charset)/2] // Simplified for this example
	}
	return string(b)
}