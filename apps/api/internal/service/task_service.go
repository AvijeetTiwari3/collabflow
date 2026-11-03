package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/model"
	"github.com/sde/collabflow/internal/repository"
)

// TaskService handles task business logic including fractional indexing reordering
type TaskService interface {
	Create(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, req model.CreateTaskRequest) (*model.Task, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Task, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateTaskRequest) (*model.Task, error)
	Reorder(ctx context.Context, taskID uuid.UUID, req model.ReorderTaskRequest) (*model.Task, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Search(ctx context.Context, query string) ([]model.Task, error)
}

type taskService struct {
	repo repository.TaskRepository
}

// NewTaskService creates a task service instance
func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) Create(ctx context.Context, userID uuid.UUID, projectID uuid.UUID, req model.CreateTaskRequest) (*model.Task, error) {
	// Find current tasks to determine append position
	tasks, err := s.repo.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var position float64 = DefaultInitialPosition
	if len(tasks) > 0 {
		lastTaskPos := tasks[len(tasks)-1].Position
		position = lastTaskPos + DefaultStep
	}

	task := &model.Task{
		ProjectID:   projectID,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		Priority:    req.Priority,
		Position:    position,
		AssigneeID:  req.AssigneeID,
		DueDate:     req.DueDate,
		CreatedBy:   userID,
	}

	if err := s.repo.Create(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) GetByID(ctx context.Context, id uuid.UUID) (*model.Task, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *taskService) ListByProject(ctx context.Context, projectID uuid.UUID) ([]model.Task, error) {
	return s.repo.ListByProject(ctx, projectID)
}

func (s *taskService) Update(ctx context.Context, id uuid.UUID, req model.UpdateTaskRequest) (*model.Task, error) {
	task, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		task.Title = *req.Title
	}
	if req.Description != nil {
		task.Description = *req.Description
	}
	if req.Status != nil {
		task.Status = *req.Status
	}
	if req.Priority != nil {
		task.Priority = *req.Priority
	}
	if req.AssigneeID != nil {
		task.AssigneeID = req.AssigneeID
	}
	if req.DueDate != nil {
		task.DueDate = req.DueDate
	}

	if err := s.repo.Update(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// Reorder applies fractional indexing to update order in O(1)
func (s *taskService) Reorder(ctx context.Context, taskID uuid.UUID, req model.ReorderTaskRequest) (*model.Task, error) {
	newPos, _ := CalculatePosition(req.PreviousPosition, req.NextPosition)

	if err := s.repo.UpdatePosition(ctx, taskID, newPos, req.Status); err != nil {
		return nil, err
	}

	return s.repo.GetByID(ctx, taskID)
}

func (s *taskService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *taskService) Search(ctx context.Context, query string) ([]model.Task, error) {
	return s.repo.Search(ctx, query)
}
