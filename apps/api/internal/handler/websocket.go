package handler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/sde/collabflow/internal/auth"
	"github.com/sde/collabflow/internal/ws"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024 * 64,
	WriteBufferSize: 1024 * 64,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for WebSockets
		return true
	},
}

// WebSocketHandler handles real-time collaboration connections
type WebSocketHandler struct {
	hub      *ws.Hub
	tokenMgr *auth.TokenManager
}

// NewWebSocketHandler creates a WebSocketHandler
func NewWebSocketHandler(hub *ws.Hub, tokenMgr *auth.TokenManager) *WebSocketHandler {
	return &WebSocketHandler{
		hub:      hub,
		tokenMgr: tokenMgr,
	}
}

// HandleWS handles GET /ws?token=...&room=...
func (h *WebSocketHandler) HandleWS(c *gin.Context) {
	roomID := c.Query("room")
	if roomID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing_room_parameter"})
		return
	}

	// Validate token from query param (WebSockets cannot send custom headers from browser API)
	tokenStr := c.Query("token")
	var userID uuid.UUID
	var userName string = "Anonymous Collaborator"
	var avatarURL string = "https://api.dicebear.com/7.x/initials/svg?seed=Anon"

	if tokenStr != "" {
		claims, err := h.tokenMgr.ValidateAccessToken(tokenStr)
		if err == nil {
			userID = claims.UserID
			userName = claims.Name
			avatarURL = "https://api.dicebear.com/7.x/initials/svg?seed=" + claims.Name
		} else {
			userID = uuid.New()
		}
	} else {
		userID = uuid.New()
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[ws-upgrade] error upgrading connection: %v", err)
		return
	}

	client := &ws.Client{
		Hub:       h.hub,
		Conn:      conn,
		Send:      make(chan []byte, 256),
		UserID:    userID,
		UserName:  userName,
		AvatarURL: avatarURL,
		RoomID:    roomID,
	}

	h.hub.Register <- client

	// Start goroutines for read & write pumps
	go client.WritePump()
	go client.ReadPump()
}

// GetPresence handles GET /api/v1/rooms/:id/presence
func (h *WebSocketHandler) GetPresence(c *gin.Context) {
	roomID := c.Param("id")
	users, err := h.hub.GetPresenceTracker().GetRoomUsers(roomID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed_to_fetch_presence"})
		return
	}

	c.JSON(http.StatusOK, users)
}
