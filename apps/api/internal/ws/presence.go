package ws

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sde/collabflow/internal/model"
)

var userColors = []string{
	"#3B82F6", // blue
	"#10B981", // green
	"#F59E0B", // amber
	"#EF4444", // red
	"#8B5CF6", // purple
	"#EC4899", // pink
	"#06B6D4", // cyan
}

// PresenceTracker tracks active users per room
type PresenceTracker struct {
	mu          sync.RWMutex
	localRooms  map[string]map[string]model.PresenceUser
	redisClient *redis.Client
}

// NewPresenceTracker creates a presence tracker instance
func NewPresenceTracker(client *redis.Client) *PresenceTracker {
	return &PresenceTracker{
		localRooms:  make(map[string]map[string]model.PresenceUser),
		redisClient: client,
	}
}

// AddUser adds an active user to the room presence set
func (p *PresenceTracker) AddUser(roomID string, user model.PresenceUser) error {
	user.LastSeen = time.Now()
	// Assign deterministic color based on user ID
	colorIndex := int(user.ID[0]) % len(userColors)
	user.Color = userColors[colorIndex]

	if p.redisClient != nil {
		data, err := json.Marshal(user)
		if err != nil {
			return err
		}
		key := "presence:" + roomID
		pipe := p.redisClient.Pipeline()
		pipe.HSet(context.Background(), key, user.ID.String(), data)
		pipe.Expire(context.Background(), key, 5*time.Minute)
		_, err = pipe.Exec(context.Background())
		return err
	}

	// In-memory fallback
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.localRooms[roomID] == nil {
		p.localRooms[roomID] = make(map[string]model.PresenceUser)
	}
	p.localRooms[roomID][user.ID.String()] = user
	return nil
}

// RemoveUser removes user from presence set
func (p *PresenceTracker) RemoveUser(roomID string, userID string) error {
	if p.redisClient != nil {
		key := "presence:" + roomID
		return p.redisClient.HDel(context.Background(), key, userID).Err()
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	if users, ok := p.localRooms[roomID]; ok {
		delete(users, userID)
		if len(users) == 0 {
			delete(p.localRooms, roomID)
		}
	}
	return nil
}

// GetRoomUsers retrieves active users in a room
func (p *PresenceTracker) GetRoomUsers(roomID string) ([]model.PresenceUser, error) {
	if p.redisClient != nil {
		key := "presence:" + roomID
		vals, err := p.redisClient.HGetAll(context.Background(), key).Result()
		if err != nil {
			return nil, err
		}

		users := make([]model.PresenceUser, 0, len(vals))
		for _, v := range vals {
			var u model.PresenceUser
			if err := json.Unmarshal([]byte(v), &u); err == nil {
				users = append(users, u)
			}
		}
		return users, nil
	}

	p.mu.RLock()
	defer p.mu.RUnlock()

	users := make([]model.PresenceUser, 0)
	if roomUsers, ok := p.localRooms[roomID]; ok {
		for _, u := range roomUsers {
			users = append(users, u)
		}
	}
	return users, nil
}
