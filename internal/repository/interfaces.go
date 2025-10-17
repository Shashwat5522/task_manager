package repository

import (
	"context"

	"github.com/vedologic/task-manager/internal/domain"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *domain.User) error

	// FindByEmail finds a user by email
	FindByEmail(ctx context.Context, email string) (*domain.User, error)

	// FindByID finds a user by ID
	FindByID(ctx context.Context, id string) (*domain.User, error)

	// ExistsByEmail checks if a user with the given email exists
	ExistsByEmail(ctx context.Context, email string) (bool, error)
}

// TaskRepository defines the interface for task data operations
type TaskRepository interface {
	// Create creates a new task
	Create(ctx context.Context, task *domain.Task) error

	// FindByID finds a task by ID
	FindByID(ctx context.Context, id string) (*domain.Task, error)

	// FindByUserID finds all tasks for a user with filtering and pagination
	FindByUserID(ctx context.Context, userID string, page, limit int, status string) ([]domain.Task, int64, error)

	// Update updates a task
	Update(ctx context.Context, task *domain.Task) error

	// Delete deletes a task
	Delete(ctx context.Context, id string, userID string) error

	// BulkUpdateStatus updates the status of multiple tasks
	BulkUpdateStatus(ctx context.Context, taskIDs []string, userID string, status domain.TaskStatus) error

	// ExistsByID checks if a task exists and belongs to the user
	ExistsByID(ctx context.Context, id string, userID string) (bool, error)
}
