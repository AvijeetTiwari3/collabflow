package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// WSMessage represents the typed wire protocol for WebSocket events
type WSMessage struct {
	Type      string          `json:"type"`                // join_room, crdt_update, cursor_move, task_updated, presence_sync
	RoomID    string          `json:"room_id,omitempty"`   // project-id or task-id
	UserID    uuid.UUID       `json:"user_id,omitempty"`
	UserName  string          `json:"user_name,omitempty"`
	AvatarURL string          `json:"avatar_url,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`   // generic JSON or base64 CRDT update
	Timestamp time.Time       `json:"timestamp"`
}

// PresenceUser represents a currently active user in a collaborative document/board
type PresenceUser struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	AvatarURL string    `json:"avatar_url"`
	Color     string    `json:"color"`      // Distinct cursor color
	LastSeen  time.Time `json:"last_seen"`
}

// CursorPayload represents cursor positioning for collaborative text / board editing
type CursorPayload struct {
	TaskID   string  `json:"task_id"`
	Position int     `json:"position"`
	Anchor   int     `json:"anchor"`
	Head     int     `json:"head"`
}
