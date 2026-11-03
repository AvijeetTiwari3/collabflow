package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sde/collabflow/internal/service"
)

// SearchHandler handles full-text search queries
type SearchHandler struct {
	taskService service.TaskService
}

// NewSearchHandler creates search handler
func NewSearchHandler(taskService service.TaskService) *SearchHandler {
	return &SearchHandler{taskService: taskService}
}

// Search handles GET /api/v1/search?q=...
func (h *SearchHandler) Search(c *gin.Context) {
	q := c.Query("q")
	if q == "" {
		c.JSON(http.StatusOK, []interface{}{})
		return
	}

	results, err := h.taskService.Search(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "search_failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, results)
}
