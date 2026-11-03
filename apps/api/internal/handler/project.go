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

// ProjectHandler handles project endpoints
type ProjectHandler struct {
	service service.ProjectService
}

// NewProjectHandler creates a project handler
func NewProjectHandler(service service.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// Create handles POST /api/v1/projects
func (h *ProjectHandler) Create(c *gin.Context) {
	userID := c.MustGet(middleware.ContextUserIDKey).(uuid.UUID)

	var req model.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation_error", "message": err.Error()})
		return
	}

	project, err := h.service.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "project_create_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, project)
}

// List handles GET /api/v1/projects
func (h *ProjectHandler) List(c *gin.Context) {
	projects, err := h.service.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, projects)
}

// GetByID handles GET /api/v1/projects/:id
func (h *ProjectHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_project_id"})
		return
	}

	project, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrProjectNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "project_not_found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "fetch_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, project)
}
