package model

import (
	"time"

	"github.com/google/uuid"
)

// Task represents a kanban task or collaborative card
type Task struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	ProjectID   uuid.UUID  `json:"project_id" db:"project_id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"` // CRDT doc representation or markdown
	Status      string     `json:"status" db:"status"`           // todo, in_progress, in_review, done
	Priority    string     `json:"priority" db:"priority"`       // low, medium, high, urgent
	Position    float64    `json:"position" db:"position"`       // Fractional index for O(1) drag-and-drop
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty" db:"assignee_id"`
	DueDate     *time.Time `json:"due_date,omitempty" db:"due_date"`
	CreatedBy   uuid.UUID  `json:"created_by" db:"created_by"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

// CreateTaskRequest payload
type CreateTaskRequest struct {
	Title       string     `json:"title" binding:"required,min=1,max=500"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Priority    string     `json:"priority"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// UpdateTaskRequest payload
type UpdateTaskRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	Status      *string    `json:"status,omitempty"`
	Priority    *string    `json:"priority,omitempty"`
	Position    *float64   `json:"position,omitempty"`
	AssigneeID  *uuid.UUID `json:"assignee_id,omitempty"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

// ReorderTaskRequest uses fractional indexing
type ReorderTaskRequest struct {
	PreviousPosition *float64 `json:"prev_pos"`
	NextPosition     *float64 `json:"next_pos"`
	Status           string   `json:"status" binding:"required"`
}
