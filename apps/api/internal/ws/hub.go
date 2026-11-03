package ws

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/sde/collabflow/internal/model"
)

// Hub maintains the set of active clients and broadcasts messages to rooms
type Hub struct {
	mu          sync.RWMutex
	rooms       map[string]map[*Client]bool
	Register    chan *Client
	Unregister  chan *Client
	Broadcast   chan model.WSMessage
	redisClient *redis.Client
	presence    *PresenceTracker
}

// NewHub creates a new Hub instance
func NewHub(redisClient *redis.Client) *Hub {
	h := &Hub{
		rooms:       make(map[string]map[*Client]bool),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan model.WSMessage, 1024),
		redisClient: redisClient,
		presence:    NewPresenceTracker(redisClient),
	}

	if redisClient != nil {
		go h.subscribeRedisPubSub()
	}

	return h
}

// Run starts the event loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

			// Update presence
			_ = h.presence.AddUser(client.RoomID, model.PresenceUser{
				ID:        client.UserID,
				Name:      client.UserName,
				AvatarURL: client.AvatarURL,
			})

			// Notify room that user joined
			h.broadcastLocal(model.WSMessage{
				Type:      "user_joined",
				RoomID:    client.RoomID,
				UserID:    client.UserID,
				UserName:  client.UserName,
				AvatarURL: client.AvatarURL,
			})

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()

			// Remove presence
			_ = h.presence.RemoveUser(client.RoomID, client.UserID.String())

			// Notify room that user left
			h.broadcastLocal(model.WSMessage{
				Type:     "user_left",
				RoomID:   client.RoomID,
				UserID:   client.UserID,
				UserName: client.UserName,
			})

		case msg := <-h.Broadcast:
			// 1. If Redis is active, publish to multi-pod channel
			if h.redisClient != nil {
				data, err := json.Marshal(msg)
				if err == nil {
					channel := "collabflow:room:" + msg.RoomID
					_ = h.redisClient.Publish(context.Background(), channel, data).Err()
				}
			} else {
				// Local fallback
				h.broadcastLocal(msg)
			}
		}
	}
}

// broadcastLocal pushes message to all WebSocket clients connected to THIS pod in that room
func (h *Hub) broadcastLocal(msg model.WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	clients, ok := h.rooms[msg.RoomID]
	if !ok {
		return
	}

	for client := range clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}

// subscribeRedisPubSub listens to multi-pod room broadcasts
func (h *Hub) subscribeRedisPubSub() {
	ctx := context.Background()
	pubsub := h.redisClient.PSubscribe(ctx, "collabflow:room:*")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var wsMsg model.WSMessage
		if err := json.Unmarshal([]byte(msg.Payload), &wsMsg); err != nil {
			continue
		}
		// Distribute to local clients connected to this pod
		h.broadcastLocal(wsMsg)
	}
}

// GetPresenceTracker returns the presence tracker
func (h *Hub) GetPresenceTracker() *PresenceTracker {
	return h.presence
}
