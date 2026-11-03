package model

import (
	"time"

	"github.com/google/uuid"
)

// Workspace represents multi-tenant container for projects and teams
type Workspace struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Slug      string    `json:"slug" db:"slug"`
	OwnerID   uuid.UUID `json:"owner_id" db:"owner_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Project represents a collaborative board or roadmap
type Project struct {
	ID          uuid.UUID `json:"id" db:"id"`
	WorkspaceID uuid.UUID `json:"workspace_id" db:"workspace_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedBy   uuid.UUID `json:"created_by" db:"created_by"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// CreateProjectRequest payload
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required,min=2,max=200"`
	Description string `json:"description"`
}
