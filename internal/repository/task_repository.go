package repository

import (
	"context"

	"github.com/vedologic/task-manager/internal/domain"
)

type TaskRepository interface {
	Create(ctx context.Context, task *domain.Task) error
	FindByID(ctx context.Context, id string) (*domain.Task, error)
	FindByUserID(ctx context.Context, userID string, page, limit int, status string) ([]domain.Task, int64, error)
	Update(ctx context.Context, task *domain.Task) error
	Delete(ctx context.Context, id string, userID string) error
	BulkUpdateStatus(ctx context.Context, taskIDs []string, userID string, status domain.TaskStatus) error
}
