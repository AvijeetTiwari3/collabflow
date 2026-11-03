package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sde/collabflow/internal/middleware"
	"github.com/sde/collabflow/internal/model"
	"github.com/sde/collabflow/internal/repository"
	"github.com/sde/collabflow/internal/service"
)

// TaskHandler handles task management endpoints
type TaskHandler struct {
	service service.TaskService
}

// NewTaskHandler creates a task handler instance
func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// Create handles POST /api/v1/projects/:id/tasks
func (h *TaskHandler) Create(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(uuid.UUID)
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_project_id"})
		return
	}

	var req model.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	task, err := h.service.Create(c.Request.Context(), userID, projectID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "task_create_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

// ListByProject handles GET /api/v1/projects/:id/tasks
func (h *TaskHandler) ListByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_project_id"})
		return
	}

	tasks, err := h.service.ListByProject(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// Update handles PATCH /api/v1/tasks/:id
func (h *TaskHandler) Update(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_task_id"})
		return
	}

	var req model.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	task, err := h.service.Update(c.Request.Context(), taskID, req)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "update_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Reorder handles POST /api/v1/tasks/:id/reorder (Fractional Indexing)
func (h *TaskHandler) Reorder(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_task_id"})
		return
	}

	var req model.ReorderTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	task, err := h.service.Reorder(c.Request.Context(), taskID, req)
	if err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "reorder_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

// Delete handles DELETE /api/v1/tasks/:id
func (h *TaskHandler) Delete(c *gin.Context) {
	taskIDStr := c.Param("id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_task_id"})
		return
	}

	if err := h.service.Delete(c.Request.Context(), taskID); err != nil {
		if errors.Is(err, repository.ErrTaskNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "task_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "delete_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
