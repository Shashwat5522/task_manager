package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vedologic/task-manager/internal/domain"
)

// taskRepository implements TaskRepository interface using raw SQL
type taskRepository struct {
	db *sqlx.DB
}

// NewTaskRepository creates a new task repository instance
func NewTaskRepository(db *sqlx.DB) TaskRepository {
	return &taskRepository{
		db: db,
	}
}

// SQL Queries
const (
	queryCreateTask = `
		INSERT INTO tasks (user_id, title, description, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	queryFindTaskByID = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE id = $1
	`

	queryFindTasksByUserID = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1
	`

	queryFindTasksByUserIDWithStatus = `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tasks
		WHERE user_id = $1 AND status = $2
	`

	queryUpdateTask = `
		UPDATE tasks
		SET title = $1, description = $2, status = $3, updated_at = $4
		WHERE id = $5 AND user_id = $6
	`

	queryDeleteTask = `
		DELETE FROM tasks
		WHERE id = $1 AND user_id = $2
	`

	queryBulkUpdateStatus = `
		UPDATE tasks
		SET status = $1, updated_at = $2
		WHERE id = ANY($3) AND user_id = $4
	`

	queryTaskExists = `
		SELECT EXISTS(SELECT 1 FROM tasks WHERE id = $1 AND user_id = $2)
	`

	queryCountTasksByUserID = `
		SELECT COUNT(*) FROM tasks WHERE user_id = $1
	`

	queryCountTasksByUserIDWithStatus = `
		SELECT COUNT(*) FROM tasks WHERE user_id = $1 AND status = $2
	`
)

// Create creates a new task in the database
func (r *taskRepository) Create(ctx context.Context, task *domain.Task) error {
	err := r.db.QueryRowContext(
		ctx,
		queryCreateTask,
		task.UserID,
		task.Title,
		task.Description,
		task.Status,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)

	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}

// FindByID finds a task by ID
func (r *taskRepository) FindByID(ctx context.Context, id string) (*domain.Task, error) {
	task := &domain.Task{}

	err := r.db.GetContext(ctx, task, queryFindTaskByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("task not found with id: %s", id)
		}
		return nil, fmt.Errorf("failed to find task by id: %w", err)
	}

	return task, nil
}

// FindByUserID finds all tasks for a user with filtering and pagination
func (r *taskRepository) FindByUserID(ctx context.Context, userID string, page, limit int, status string) ([]domain.Task, int64, error) {
	offset := (page - 1) * limit

	// Build query based on status filter
	query := queryFindTasksByUserID + fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)
	var countQuery string
	var args []interface{}

	if status != "" {
		query = queryFindTasksByUserIDWithStatus + fmt.Sprintf(" ORDER BY created_at DESC LIMIT %d OFFSET %d", limit, offset)
		countQuery = queryCountTasksByUserIDWithStatus
		args = []interface{}{userID, status}
	} else {
		countQuery = queryCountTasksByUserID
		args = []interface{}{userID}
	}

	// Get total count
	var total int64
	err := r.db.GetContext(ctx, &total, countQuery, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count tasks: %w", err)
	}

	// Get tasks
	var tasks []domain.Task
	if status != "" {
		err = r.db.SelectContext(ctx, &tasks, query, userID, status)
	} else {
		err = r.db.SelectContext(ctx, &tasks, query, userID)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []domain.Task{}, total, nil
		}
		return nil, 0, fmt.Errorf("failed to find tasks by user id: %w", err)
	}

	return tasks, total, nil
}

// Update updates an existing task
func (r *taskRepository) Update(ctx context.Context, task *domain.Task) error {
	result, err := r.db.ExecContext(
		ctx,
		queryUpdateTask,
		task.Title,
		task.Description,
		task.Status,
		task.UpdatedAt,
		task.ID,
		task.UserID,
	)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found or not owned by user")
	}

	return nil
}

// Delete deletes a task (owned by user)
func (r *taskRepository) Delete(ctx context.Context, id string, userID string) error {
	result, err := r.db.ExecContext(ctx, queryDeleteTask, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("task not found or not owned by user")
	}

	return nil
}

// BulkUpdateStatus updates the status of multiple tasks
func (r *taskRepository) BulkUpdateStatus(ctx context.Context, taskIDs []string, userID string, status domain.TaskStatus) error {
	result, err := r.db.ExecContext(
		ctx,
		queryBulkUpdateStatus,
		status,
		// Current timestamp for updated_at
		"NOW()",
		pq.Array(taskIDs),
		userID,
	)

	if err != nil {
		return fmt.Errorf("failed to bulk update task status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no tasks updated")
	}

	return nil
}

// ExistsByID checks if a task exists and belongs to the user
func (r *taskRepository) ExistsByID(ctx context.Context, id string, userID string) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, queryTaskExists, id, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check if task exists: %w", err)
	}

	return exists, nil
}
