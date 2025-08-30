package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	gorillaws "github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/francknouama/go-starter/internal/web/websocket"
)

func setupWebSocketTestRouter() (*gin.Engine, *websocket.Hub) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	
	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run() // Start hub in background
	
	handler := NewWebSocketHandler(hub)
	
	// WebSocket routes
	api := router.Group("/api/v1")
	{
		api.GET("/ws/generate", handler.HandleGenerateWS)
		api.GET("/ws/preview", handler.HandlePreviewWS)
	}
	
	return router, hub
}

func TestWebSocketHandler_HandleGenerateWS(t *testing.T) {
	router, hub := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/generate"
	
	t.Run("successful WebSocket connection", func(t *testing.T) {
		// Connect to WebSocket
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "Should be able to connect to WebSocket")
		defer conn.Close()
		
		// Set read deadline to avoid hanging tests
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read welcome message
		var message map[string]interface{}
		err = conn.ReadJSON(&message)
		require.NoError(t, err, "Should receive welcome message")
		
		// Verify welcome message structure
		assert.Equal(t, "connected", message["type"])
		assert.Equal(t, "Connected to generation WebSocket", message["message"])
		assert.NotEmpty(t, message["client_id"], "Should have client ID")
	})
	
	t.Run("multiple clients can connect", func(t *testing.T) {
		const numClients = 3
		connections := make([]*gorillaws.Conn, numClients)
		clientIDs := make([]string, numClients)
		
		// Connect multiple clients
		for i := 0; i < numClients; i++ {
			conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
			require.NoError(t, err, "Should be able to connect client %d", i)
			defer conn.Close()
			
			conn.SetReadDeadline(time.Now().Add(5 * time.Second))
			connections[i] = conn
			
			// Read welcome message
			var message map[string]interface{}
			err = conn.ReadJSON(&message)
			require.NoError(t, err, "Client %d should receive welcome message", i)
			
			clientID, ok := message["client_id"].(string)
			require.True(t, ok, "Client ID should be string")
			clientIDs[i] = clientID
		}
		
		// Verify all client IDs are unique
		idMap := make(map[string]bool)
		for _, id := range clientIDs {
			assert.False(t, idMap[id], "Client ID %s should be unique", id)
			idMap[id] = true
		}
	})
	
	t.Run("broadcast to generate clients", func(t *testing.T) {
		// Connect to generate WebSocket
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "Should be able to connect to WebSocket")
		defer conn.Close()
		
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read and discard welcome message
		var welcomeMsg map[string]interface{}
		err = conn.ReadJSON(&welcomeMsg)
		require.NoError(t, err, "Should receive welcome message")
		
		// Send a test message via hub
		testMessage := map[string]interface{}{
			"type":    "generation_progress",
			"message": "Test generation progress",
			"progress": 50,
		}
		
		// Use a goroutine to send the message after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			hub.BroadcastToType("generate", testMessage)
		}()
		
		// Read the broadcasted message
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var receivedMessage map[string]interface{}
		err = conn.ReadJSON(&receivedMessage)
		require.NoError(t, err, "Should receive broadcasted message")
		
		assert.Equal(t, "generation_progress", receivedMessage["type"])
		assert.Equal(t, "Test generation progress", receivedMessage["message"])
		assert.Equal(t, float64(50), receivedMessage["progress"])
	})
}

func TestWebSocketHandler_HandlePreviewWS(t *testing.T) {
	router, hub := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	// Convert HTTP URL to WebSocket URL
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/preview"
	
	t.Run("successful WebSocket connection", func(t *testing.T) {
		// Connect to WebSocket
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "Should be able to connect to WebSocket")
		defer conn.Close()
		
		// Set read deadline to avoid hanging tests
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read welcome message
		var message map[string]interface{}
		err = conn.ReadJSON(&message)
		require.NoError(t, err, "Should receive welcome message")
		
		// Verify welcome message structure
		assert.Equal(t, "connected", message["type"])
		assert.Equal(t, "Connected to preview WebSocket", message["message"])
		assert.NotEmpty(t, message["client_id"], "Should have client ID")
	})
	
	t.Run("broadcast to preview clients", func(t *testing.T) {
		// Connect to preview WebSocket
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "Should be able to connect to WebSocket")
		defer conn.Close()
		
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read and discard welcome message
		var welcomeMsg map[string]interface{}
		err = conn.ReadJSON(&welcomeMsg)
		require.NoError(t, err, "Should receive welcome message")
		
		// Send a test preview message via hub
		previewMessage := map[string]interface{}{
			"type":     "preview_update",
			"file":     "main.go",
			"content":  "package main\n\nfunc main() {\n\t// Generated code\n}",
			"line_count": 5,
		}
		
		// Use a goroutine to send the message after a short delay
		go func() {
			time.Sleep(100 * time.Millisecond)
			hub.BroadcastToType("preview", previewMessage)
		}()
		
		// Read the broadcasted message
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var receivedMessage map[string]interface{}
		err = conn.ReadJSON(&receivedMessage)
		require.NoError(t, err, "Should receive broadcasted message")
		
		assert.Equal(t, "preview_update", receivedMessage["type"])
		assert.Equal(t, "main.go", receivedMessage["file"])
		assert.Contains(t, receivedMessage["content"], "package main")
		assert.Equal(t, float64(5), receivedMessage["line_count"])
	})
}

func TestWebSocketHandler_ClientIsolation(t *testing.T) {
	router, hub := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	generateURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/generate"
	previewURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/preview"
	
	t.Run("generate and preview clients are isolated", func(t *testing.T) {
		// Connect to generate WebSocket
		generateConn, _, err := gorillaws.DefaultDialer.Dial(generateURL, nil)
		require.NoError(t, err, "Should be able to connect to generate WebSocket")
		defer generateConn.Close()
		
		// Connect to preview WebSocket
		previewConn, _, err := gorillaws.DefaultDialer.Dial(previewURL, nil)
		require.NoError(t, err, "Should be able to connect to preview WebSocket")
		defer previewConn.Close()
		
		generateConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		previewConn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read welcome messages
		var genWelcome, prevWelcome map[string]interface{}
		err = generateConn.ReadJSON(&genWelcome)
		require.NoError(t, err, "Should receive generate welcome")
		err = previewConn.ReadJSON(&prevWelcome)
		require.NoError(t, err, "Should receive preview welcome")
		
		// Send message to generate clients only
		generateMessage := map[string]interface{}{
			"type":    "generation_started",
			"message": "Generation started",
		}
		
		go func() {
			time.Sleep(100 * time.Millisecond)
			hub.BroadcastToType("generate", generateMessage)
		}()
		
		// Generate client should receive the message
		generateConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		var genMessage map[string]interface{}
		err = generateConn.ReadJSON(&genMessage)
		require.NoError(t, err, "Generate client should receive message")
		assert.Equal(t, "generation_started", genMessage["type"])
		
		// Preview client should NOT receive the message (should timeout)
		previewConn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		var prevMessage map[string]interface{}
		err = previewConn.ReadJSON(&prevMessage)
		assert.Error(t, err, "Preview client should not receive generate message")
		
		// Verify it's a timeout error (WebSocket read deadline)
		if netErr, ok := err.(interface{ Timeout() bool }); ok {
			assert.True(t, netErr.Timeout(), "Should be timeout error")
		}
	})
}

func TestWebSocketHandler_ErrorHandling(t *testing.T) {
	router, _ := setupWebSocketTestRouter()
	
	t.Run("invalid WebSocket upgrade request", func(t *testing.T) {
		// Make regular HTTP request instead of WebSocket upgrade
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/ws/generate", nil)
		// Don't set WebSocket headers
		
		router.ServeHTTP(w, req)
		
		// Should return error response for invalid upgrade
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		
		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")
		assert.Equal(t, "Failed to upgrade to WebSocket", response["error"])
	})
}

func TestWebSocketHandler_ConnectionLifecycle(t *testing.T) {
	router, hub := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/generate"
	
	t.Run("connection cleanup on close", func(t *testing.T) {
		// Connect to WebSocket
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		require.NoError(t, err, "Should be able to connect to WebSocket")
		
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		
		// Read welcome message
		var message map[string]interface{}
		err = conn.ReadJSON(&message)
		require.NoError(t, err, "Should receive welcome message")
		
		clientID := message["client_id"].(string)
		assert.NotEmpty(t, clientID, "Should have client ID")
		
		// Close connection
		conn.Close()
		
		// Give hub time to process disconnection
		time.Sleep(100 * time.Millisecond)
		
		// Try to send message to disconnected client
		// This should not cause panic or error
		testMessage := map[string]interface{}{
			"type":    "test",
			"message": "This should not reach disconnected client",
		}
		
		// Should not panic
		assert.NotPanics(t, func() {
			hub.BroadcastToType("generate", testMessage)
		})
	})
}

// Benchmark tests for WebSocket performance
func BenchmarkWebSocketHandler_ConnectionSetup(b *testing.B) {
	router, _ := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/generate"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Connect and immediately close
		conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
		if err != nil {
			b.Fatalf("Failed to connect: %v", err)
		}
		
		// Read welcome message
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		var message map[string]interface{}
		err = conn.ReadJSON(&message)
		if err != nil {
			conn.Close()
			b.Fatalf("Failed to read welcome message: %v", err)
		}
		
		conn.Close()
	}
}

func BenchmarkWebSocketHandler_MessageBroadcast(b *testing.B) {
	router, hub := setupWebSocketTestRouter()
	
	// Create test server
	server := httptest.NewServer(router)
	defer server.Close()
	
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/api/v1/ws/generate"
	
	// Setup persistent connection
	conn, _, err := gorillaws.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		b.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	
	// Read welcome message
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	var welcomeMsg map[string]interface{}
	err = conn.ReadJSON(&welcomeMsg)
	if err != nil {
		b.Fatalf("Failed to read welcome message: %v", err)
	}
	
	testMessage := map[string]interface{}{
		"type":    "benchmark",
		"message": "Benchmark message",
		"data":    fmt.Sprintf("Test data %d", 12345),
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Broadcast message
		hub.BroadcastToType("generate", testMessage)
		
		// Read the message
		conn.SetReadDeadline(time.Now().Add(1 * time.Second))
		var receivedMsg map[string]interface{}
		err = conn.ReadJSON(&receivedMsg)
		if err != nil {
			b.Fatalf("Failed to read message: %v", err)
		}
	}
}