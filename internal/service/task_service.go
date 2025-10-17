package service

import (
	"context"

	"github.com/vedologic/task-manager/internal/domain"
	"github.com/vedologic/task-manager/internal/dto"
)

type TaskService interface {
	Create(ctx context.Context, userID string, req dto.CreateTaskRequest) (*domain.Task, error)
	GetByID(ctx context.Context, taskID string, userID string) (*domain.Task, error)
	List(ctx context.Context, userID string, page, limit int, status string) (*dto.TaskListResponse, error)
	Update(ctx context.Context, taskID string, userID string, req dto.UpdateTaskRequest) (*domain.Task, error)
	Delete(ctx context.Context, taskID string, userID string) error
	BulkComplete(ctx context.Context, userID string, req dto.BulkCompleteRequest) (*dto.BulkCompleteResponse, error)
}
