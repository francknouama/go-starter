package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/francknouama/go-starter/internal/web/websocket"
)

type WebSocketHandler struct {
	hub *websocket.Hub
}

func NewWebSocketHandler(hub *websocket.Hub) *WebSocketHandler {
	return &WebSocketHandler{
		hub: hub,
	}
}

// HandleGenerateWS handles WebSocket connections for generation updates
func (h *WebSocketHandler) HandleGenerateWS(c *gin.Context) {
	client, err := h.hub.UpgradeConnection(c, "generate")
	if err != nil {
		slog.Error("Failed to upgrade WebSocket connection", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upgrade to WebSocket",
		})
		return
	}

	// Send welcome message
	welcomeMsg := map[string]interface{}{
		"type":    "connected",
		"message": "Connected to generation WebSocket",
		"client_id": client.ID,
	}
	
	h.hub.BroadcastToType("generate", welcomeMsg)
}

// HandlePreviewWS handles WebSocket connections for preview updates
func (h *WebSocketHandler) HandlePreviewWS(c *gin.Context) {
	client, err := h.hub.UpgradeConnection(c, "preview")
	if err != nil {
		slog.Error("Failed to upgrade WebSocket connection", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to upgrade to WebSocket",
		})
		return
	}

	// Send welcome message
	welcomeMsg := map[string]interface{}{
		"type":    "connected",
		"message": "Connected to preview WebSocket",
		"client_id": client.ID,
	}
	
	h.hub.BroadcastToType("preview", welcomeMsg)
}