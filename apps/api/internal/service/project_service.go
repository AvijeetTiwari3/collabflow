package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/model"
	"github.com/sde/collabflow/internal/repository"
)

// ProjectService handles project operations
type ProjectService interface {
	Create(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error)
	ListAll(ctx context.Context) ([]model.Project, error)
}

type projectService struct {
	repo repository.ProjectRepository
}

// NewProjectService creates a project service
func NewProjectService(repo repository.ProjectRepository) ProjectService {
	return &projectService{repo: repo}
}

func (s *projectService) Create(ctx context.Context, userID uuid.UUID, req model.CreateProjectRequest) (*model.Project, error) {
	p := &model.Project{
		WorkspaceID: uuid.New(), // Default single workspace context
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *projectService) GetByID(ctx context.Context, id uuid.UUID) (*model.Project, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *projectService) ListAll(ctx context.Context) ([]model.Project, error) {
	return s.repo.ListAll(ctx)
}
