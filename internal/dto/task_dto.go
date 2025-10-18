package dto

import "github.com/vedologic/task-manager/internal/domain"

type CreateTaskRequest struct {
	Title       string            `json:"title" binding:"required,min=1,max=255"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       string            `json:"title" binding:"required,min=1,max=255"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status" binding:"required"`
}

type BulkCompleteRequest struct {
	TaskIDs []string `json:"task_ids" binding:"required,min=1"`
}

type TaskResponse struct {
	ID          string            `json:"id"`
	UserID      string            `json:"user_id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Status      domain.TaskStatus `json:"status"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks      []TaskResponse `json:"tasks"`
	TotalCount int64          `json:"total_count"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type BulkCompleteResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedIDs    []string `json:"failed_ids,omitempty"`
}
