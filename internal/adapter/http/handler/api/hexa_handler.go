package api

import (
	"github.com/gin-gonic/gin"
	"hexagonal/internal/core/domain"
	logutil "hexagonal/internal/core/infrastructure/log"
	"net/http"
)

// RegisterHexaRoutes sets up the routes for taskPO-related operations.
func (h *Handler) RegisterHexaRoutes(router *gin.Engine) {
	group := router.Group("/tasks")
	{
		group.GET("/", h.Index)         // Use uppercase method
		group.POST("/", h.Store)        // Use uppercase method
		group.PATCH("/:id", h.Update)   // Use uppercase method
		group.DELETE("/:id", h.Destroy) // Use uppercase method
	}
}

// Store handles the saving of a taskPO.
func (h *Handler) Store(c *gin.Context) {
	correlationID := logutil.GenerateCorrelationID()
	var task domain.TaskRequest
	if err := c.ShouldBindJSON(&task); err != nil {
		//RespondWithError(c, http.StatusBadRequest, "Invalid taskPO input", err)
		return
	}

	ctx := c.Request.Context()
	if err := h.TaskService.CreateTask(ctx, task); err != nil {
		RespondWithError(c, http.StatusInternalServerError, correlationID, err, "Failed to save taskPO")
		return
	}

	RespondWithSuccess(c, http.StatusCreated, "Task created successfully")
}

// Index handles fetching tasks.
func (h *Handler) Index(c *gin.Context) {
	correlationID := logutil.GenerateCorrelationID()
	tasks, err := h.TaskService.GetAllTasks(c)
	if err != nil {
		RespondWithError(c, http.StatusInternalServerError, correlationID, err, "Failed to retrieve tasks")
		return
	}

	RespondWithSuccess(c, http.StatusOK, tasks)
}

// Update handles updating an existing task by ID.
func (h *Handler) Update(c *gin.Context) {
	correlationID := logutil.GenerateCorrelationID()
	taskID := c.Param("id")

	var updatedTask domain.TaskRequest
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		RespondWithError(c, http.StatusBadRequest, correlationID, err, "Invalid task data")
		return
	}

	if err := h.TaskService.UpdateTask(c, taskID, updatedTask); err != nil {
		RespondWithError(c, http.StatusInternalServerError, correlationID, err, "Failed to update task")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Task updated successfully")
}

// Destroy handles deleting a task by ID.
func (h *Handler) Destroy(c *gin.Context) {
	correlationID := logutil.GenerateCorrelationID()
	taskID := c.Param("id")

	if err := h.TaskService.DeleteTask(c, taskID); err != nil {
		RespondWithError(c, http.StatusInternalServerError, correlationID, err, "Failed to delete task")
		return
	}

	RespondWithSuccess(c, http.StatusOK, "Task deleted successfully")
}
