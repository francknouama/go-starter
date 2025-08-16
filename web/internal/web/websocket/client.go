package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/google/uuid"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer
	maxMessageSize = 512
)

// Client is a middleman between the websocket connection and the hub
type Client struct {
	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound messages
	send chan []byte

	// Hub reference
	hub *Hub

	// Unique client identifier
	ID string

	// Client metadata
	UserAgent string
	IP        string
	ConnectedAt time.Time
}

// NewClient creates a new Client instance
func NewClient(hub *Hub, conn *websocket.Conn) *Client {
	return &Client{
		conn:        conn,
		send:        make(chan []byte, 256),
		hub:         hub,
		ID:          uuid.New().String(),
		UserAgent:   "", // Will be set during upgrade if needed
		IP:          "", // Will be set during upgrade if needed
		ConnectedAt: time.Now(),
	}
}

// readPump pumps messages from the websocket connection to the hub
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.UnregisterClient(c)
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle incoming messages from client
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current websocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// handleMessage processes incoming messages from the client
func (c *Client) handleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error unmarshaling client message: %v", err)
		return
	}

	// Handle different message types
	msgType, ok := msg["type"].(string)
	if !ok {
		log.Printf("Invalid message type from client %s", c.ID)
		return
	}

	switch msgType {
	case "ping":
		// Respond with pong
		response := map[string]interface{}{
			"type":      "pong",
			"timestamp": time.Now(),
			"clientId":  c.ID,
		}
		c.sendMessage(response)

	case "subscribe":
		// Handle subscription to specific events
		c.handleSubscription(msg)

	case "unsubscribe":
		// Handle unsubscription from events
		c.handleUnsubscription(msg)

	case "heartbeat":
		// Update last seen time
		log.Printf("Heartbeat from client %s", c.ID)

	default:
		log.Printf("Unknown message type '%s' from client %s", msgType, c.ID)
	}
}

// handleSubscription handles client subscription requests
func (c *Client) handleSubscription(msg map[string]interface{}) {
	channel, ok := msg["channel"].(string)
	if !ok {
		log.Printf("Invalid subscription channel from client %s", c.ID)
		return
	}

	log.Printf("Client %s subscribed to channel: %s", c.ID, channel)
	
	// Send confirmation
	response := map[string]interface{}{
		"type":    "subscribed",
		"channel": channel,
		"clientId": c.ID,
	}
	c.sendMessage(response)
}

// handleUnsubscription handles client unsubscription requests
func (c *Client) handleUnsubscription(msg map[string]interface{}) {
	channel, ok := msg["channel"].(string)
	if !ok {
		log.Printf("Invalid unsubscription channel from client %s", c.ID)
		return
	}

	log.Printf("Client %s unsubscribed from channel: %s", c.ID, channel)
	
	// Send confirmation
	response := map[string]interface{}{
		"type":    "unsubscribed",
		"channel": channel,
		"clientId": c.ID,
	}
	c.sendMessage(response)
}

// sendMessage sends a message to the client
func (c *Client) sendMessage(msg interface{}) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message for client %s: %v", c.ID, err)
		return
	}

	select {
	case c.send <- data:
	default:
		// Channel is full, close the client
		close(c.send)
		log.Printf("Client %s send channel full, closing connection", c.ID)
	}
}

// Info returns client information
func (c *Client) Info() map[string]interface{} {
	return map[string]interface{}{
		"id":          c.ID,
		"userAgent":   c.UserAgent,
		"ip":          c.IP,
		"connectedAt": c.ConnectedAt,
		"uptime":      time.Since(c.ConnectedAt).String(),
	}
}

// getClientIP extracts the client IP from the request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return xff
	}
	
	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	
	// Fall back to RemoteAddr
	return r.RemoteAddr
}