package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// HealthHandler manages health and service info endpoints
type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler creates a health handler
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Root handles GET / and HEAD /
func (h *HealthHandler) Root(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"service":      "CollabFlow Real-Time Collaboration Engine",
		"status":       "HEALTHY",
		"version":      "1.0.0",
		"architecture": "Next.js 14 + Go 1.22 + WebSockets + Redis Pub/Sub + CRDT (Yjs) + PostgreSQL 16",
		"endpoints": gin.H{
			"auth_register": "POST /api/v1/auth/register",
			"auth_login":    "POST /api/v1/auth/login",
			"auth_refresh":  "POST /api/v1/auth/refresh",
			"projects":      "GET/POST /api/v1/projects",
			"tasks":         "GET/POST /api/v1/projects/:id/tasks",
			"reorder_task":  "POST /api/v1/tasks/:id/reorder (O(1) Fractional Indexing)",
			"websocket":     "GET /ws?room=:id&token=:jwt",
			"health":        "GET /health",
			"ready":         "GET /ready",
			"metrics":       "GET /metrics",
		},
	})
}

// Liveness handles GET /health
func (h *HealthHandler) Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}

// Readiness handles GET /ready
func (h *HealthHandler) Readiness(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "DOWN",
			"error":  "database unavailable: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "READY"})
}
