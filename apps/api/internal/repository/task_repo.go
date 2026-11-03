package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/model"
)

var (
	ErrTaskNotFound = errors.New("task not found")
)

// TaskRepository handles task persistence
type TaskRepository interface {
	Create(ctx context.Context, t *model.Task) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Task, error)
	Update(ctx context.Context, t *model.Task) error
	UpdatePosition(ctx context.Context, taskID uuid.UUID, newPosition float64, newStatus string) error
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, q string) ([]model.Task, error)
}

type postgresTaskRepo struct {
	db *sql.DB
}

// NewTaskRepository creates task repository
func NewTaskRepository(db *sql.DB) TaskRepository {
	return &postgresTaskRepo{db: db}
}

func (r *postgresTaskRepo) Create(ctx context.Context, t *model.Task) error {
	query := `
		INSERT INTO tasks (id, project_id, title, description, status, priority, position, assignee_id, due_date, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	t.ID = uuid.New()
	now := time.Now()
	t.CreatedAt = now
	t.UpdatedAt = now
	if t.Status == "" {
		t.Status = "todo"
	}
	if t.Priority == "" {
		t.Priority = "medium"
	}

	_, err := r.db.ExecContext(ctx, query,
		t.ID, t.ProjectID, t.Title, t.Description, t.Status, t.Priority, t.Position, t.AssigneeID, t.DueDate, t.CreatedBy, t.CreatedAt, t.UpdatedAt,
	)
	return err
}

func (r *postgresTaskRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, position, assignee_id, due_date, created_by, created_at, updated_at
		FROM tasks
		WHERE id = $1
		LIMIT 1
	`
	var t model.Task
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Position, &t.AssigneeID, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresTaskRepo) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, position, assignee_id, due_date, created_by, created_at, updated_at
		FROM tasks
		WHERE project_id = $1
		ORDER BY position ASC
	`
	rows, err := r.db.QueryContext(ctx, query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Position, &t.AssigneeID, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *postgresTaskRepo) Update(ctx context.Context, t *model.Task) error {
	query := `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = NOW()
		WHERE id = $7
	`
	res, err := r.db.ExecContext(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.AssigneeID, t.DueDate, t.ID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

// UpdatePosition executes an O(1) single-row update using fractional indexing
func (r *postgresTaskRepo) UpdatePosition(ctx context.Context, taskID uuid.UUID, newPosition float64, newStatus string) error {
	query := `
		UPDATE tasks
		SET position = $1, status = $2, updated_at = NOW()
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, newPosition, newStatus, taskID)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *postgresTaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return ErrTaskNotFound
	}
	return nil
}

func (r *postgresTaskRepo) Search(ctx context.Context, q string) ([]model.Task, error) {
	query := `
		SELECT id, project_id, title, description, status, priority, position, assignee_id, due_date, created_by, created_at, updated_at
		FROM tasks
		WHERE title ILIKE '%' || $1 || '%' OR description ILIKE '%' || $1 || '%'
		ORDER BY updated_at DESC
		LIMIT 50
	`
	rows, err := r.db.QueryContext(ctx, query, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]model.Task, 0)
	for rows.Next() {
		var t model.Task
		if err := rows.Scan(
			&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.Position, &t.AssigneeID, &t.DueDate, &t.CreatedBy, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}
