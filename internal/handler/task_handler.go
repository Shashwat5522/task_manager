package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vedologic/task-manager/internal/dto"
	"github.com/vedologic/task-manager/internal/service"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskService service.TaskService
	log         *zap.Logger
}

// NewTaskHandler creates a new task handler
func NewTaskHandler(taskService service.TaskService, log *zap.Logger) *TaskHandler {
	return &TaskHandler{
		taskService: taskService,
		log:         log,
	}
}

// Create godoc
// @Summary Create a new task
// @Description Create a new task for the authenticated user
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.CreateTaskRequest true "Create task request"
// @Success 201 {object} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks [post]
func (h *TaskHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid create task request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	task, err := h.taskService.Create(c.Request.Context(), userID.(string), req)
	if err != nil {
		h.log.Error("Failed to create task", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, task)
}

// GetByID godoc
// @Summary Get a task by ID
// @Description Get a specific task by its ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks/{id} [get]
func (h *TaskHandler) GetByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")

	task, err := h.taskService.GetByID(c.Request.Context(), taskID, userID.(string))
	if err != nil {
		h.log.Warn("Failed to get task", zap.Error(err))
		c.JSON(404, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(200, task)
}

// List godoc
// @Summary List user tasks
// @Description Get all tasks for the authenticated user with pagination and filtering
// @Tags tasks
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param status query string false "Filter by status"
// @Success 200 {object} dto.TaskListResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks [get]
func (h *TaskHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")

	tasks, err := h.taskService.List(c.Request.Context(), userID.(string), page, limit, status)
	if err != nil {
		h.log.Error("Failed to list tasks", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, tasks)
}

// Update godoc
// @Summary Update a task
// @Description Update an existing task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Param request body dto.UpdateTaskRequest true "Update task request"
// @Success 200 {object} domain.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks/{id} [put]
func (h *TaskHandler) Update(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")
	var req dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid update task request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	task, err := h.taskService.Update(c.Request.Context(), taskID, userID.(string), req)
	if err != nil {
		h.log.Error("Failed to update task", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, task)
}

// Delete godoc
// @Summary Delete a task
// @Description Delete a specific task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 204
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks/{id} [delete]
func (h *TaskHandler) Delete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	taskID := c.Param("id")

	err := h.taskService.Delete(c.Request.Context(), taskID, userID.(string))
	if err != nil {
		h.log.Error("Failed to delete task", zap.Error(err))
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.Status(204)
}

// BulkComplete godoc
// @Summary Mark multiple tasks as completed
// @Description Mark multiple tasks as completed concurrently
// @Tags tasks
// @Accept json
// @Produce json
// @Param request body dto.BulkCompleteRequest true "Bulk complete request"
// @Success 200 {object} dto.BulkCompleteResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Security BearerAuth
// @Router /api/v1/tasks/bulk-complete [patch]
func (h *TaskHandler) BulkComplete(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req dto.BulkCompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Warn("Invalid bulk complete request", zap.Error(err))
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	resp, err := h.taskService.BulkComplete(c.Request.Context(), userID.(string), req)
	if err != nil {
		h.log.Error("Failed to bulk complete tasks", zap.Error(err))
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, resp)
}
