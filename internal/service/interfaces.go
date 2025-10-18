package service

import (
	"context"

	"github.com/vedologic/task-manager/internal/domain"
	"github.com/vedologic/task-manager/internal/dto"
)

// AuthService defines the interface for authentication business logic
type AuthService interface {
	// Register registers a new user
	Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error)

	// Login authenticates a user and returns a JWT token
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
}

// TaskService defines the interface for task business logic
type TaskService interface {
	// Create creates a new task
	Create(ctx context.Context, userID string, req dto.CreateTaskRequest) (*domain.Task, error)

	// GetByID retrieves a task by ID
	GetByID(ctx context.Context, taskID string, userID string) (*domain.Task, error)

	// List retrieves all tasks for a user with pagination and filtering
	List(ctx context.Context, userID string, page, limit int, status string) (*dto.TaskListResponse, error)

	// Update updates a task
	Update(ctx context.Context, taskID string, userID string, req dto.UpdateTaskRequest) (*domain.Task, error)

	// Delete deletes a task
	Delete(ctx context.Context, taskID string, userID string) error

	// BulkComplete marks multiple tasks as done concurrently
	BulkComplete(ctx context.Context, userID string, req dto.BulkCompleteRequest) (*dto.BulkCompleteResponse, error)
}
