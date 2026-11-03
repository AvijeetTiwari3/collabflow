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
	ErrProjectNotFound = errors.New("project not found")
)

// ProjectRepository handles project persistence
type ProjectRepository interface {
	Create(ctx context.Context, p *model.Project) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	ListAll(ctx context.Context) ([]model.Project, error)
}

type postgresProjectRepo struct {
	db *sql.DB
}

// NewProjectRepository creates project repository
func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &postgresProjectRepo{db: db}
}

func (r *postgresProjectRepo) Create(ctx context.Context, p *model.Project) error {
	query := `
		INSERT INTO projects (id, workspace_id, name, description, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	p.ID = uuid.New()
	p.CreatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, p.ID, p.WorkspaceID, p.Name, p.Description, p.CreatedBy, p.CreatedAt)
	return err
}

func (r *postgresProjectRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, created_by, created_at
		FROM projects
		WHERE id = $1
		LIMIT 1
	`
	var p model.Project
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProjectRepo) ListAll(ctx context.Context) ([]model.Project, error) {
	query := `
		SELECT id, workspace_id, name, description, created_by, created_at
		FROM projects
		ORDER BY created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	projects := make([]model.Project, 0)
	for rows.Next() {
		var p model.Project
		if err := rows.Scan(&p.ID, &p.WorkspaceID, &p.Name, &p.Description, &p.CreatedBy, &p.CreatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
}
