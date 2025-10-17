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

type TaskListResponse struct {
	Tasks      []domain.Task `json:"tasks"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

type BulkCompleteResponse struct {
	SuccessCount int      `json:"success_count"`
	FailedCount  int      `json:"failed_count"`
	FailedIDs    []string `json:"failed_ids,omitempty"`
}
