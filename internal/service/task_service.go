package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/vedologic/task-manager/internal/domain"
	"github.com/vedologic/task-manager/internal/dto"
	"github.com/vedologic/task-manager/internal/repository"
)

// taskService implements TaskService interface with business logic
type taskService struct {
	taskRepo repository.TaskRepository
}

// NewTaskService creates a new task service
func NewTaskService(taskRepo repository.TaskRepository) TaskService {
	return &taskService{
		taskRepo: taskRepo,
	}
}

// Create creates a new task
func (s *taskService) Create(ctx context.Context, userID string, req dto.CreateTaskRequest) (*domain.Task, error) {
	// Validate status
	if !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid task status: %s", req.Status)
	}

	// Convert userID string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Create task entity
	task := &domain.Task{
		UserID:      userIDInt,
		Title:       req.Title,
		Description: req.Description,
		Status:      req.Status,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save to repository
	if err := s.taskRepo.Create(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

// GetByID retrieves a task by ID
func (s *taskService) GetByID(ctx context.Context, taskID string, userID string) (*domain.Task, error) {
	// Convert userID to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Find task
	task, err := s.taskRepo.FindByID(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to find task: %w", err)
	}

	// Verify ownership
	if task.UserID != userIDInt {
		return nil, fmt.Errorf("access denied: task does not belong to user")
	}

	return task, nil
}

// List retrieves all tasks for a user with pagination and filtering
func (s *taskService) List(ctx context.Context, userID string, page, limit int, status string) (*dto.TaskListResponse, error) {
	// Validate pagination
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get tasks from repository
	tasks, total, err := s.taskRepo.FindByUserID(ctx, userID, page, limit, status)
	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	// Convert to response DTOs
	taskResponses := make([]dto.TaskResponse, len(tasks))
	for i, task := range tasks {
		taskResponses[i] = dto.TaskResponse{
			ID:          fmt.Sprintf("%d", task.ID),
			UserID:      fmt.Sprintf("%d", task.UserID),
			Title:       task.Title,
			Description: task.Description,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt.String(),
			UpdatedAt:   task.UpdatedAt.String(),
		}
	}

	// Calculate total pages
	totalPages := int((total + int64(limit) - 1) / int64(limit))

	return &dto.TaskListResponse{
		Tasks:      taskResponses,
		TotalCount: total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}, nil
}

// Update updates a task
func (s *taskService) Update(ctx context.Context, taskID string, userID string, req dto.UpdateTaskRequest) (*domain.Task, error) {
	// Validate status
	if !req.Status.IsValid() {
		return nil, fmt.Errorf("invalid task status: %s", req.Status)
	}

	// Get existing task
	task, err := s.GetByID(ctx, taskID, userID)
	if err != nil {
		return nil, err
	}

	// Update fields
	task.Title = req.Title
	task.Description = req.Description
	task.Status = req.Status
	task.UpdatedAt = time.Now()

	// Save to repository
	if err := s.taskRepo.Update(ctx, task); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return task, nil
}

// Delete deletes a task
func (s *taskService) Delete(ctx context.Context, taskID string, userID string) error {
	// Verify task exists and belongs to user
	if _, err := s.GetByID(ctx, taskID, userID); err != nil {
		return err
	}

	// Delete from repository
	if err := s.taskRepo.Delete(ctx, taskID, userID); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}

// BulkComplete marks multiple tasks as done concurrently using goroutines and channels
func (s *taskService) BulkComplete(ctx context.Context, userID string, req dto.BulkCompleteRequest) (*dto.BulkCompleteResponse, error) {
	if len(req.TaskIDs) == 0 {
		return nil, fmt.Errorf("no task IDs provided")
	}

	// Convert userID string to int
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// Number of concurrent workers
	numWorkers := 5
	taskIDsChan := make(chan string, numWorkers)
	resultsChan := make(chan error, len(req.TaskIDs))

	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for taskID := range taskIDsChan {
				// Verify ownership
				existingTask, err := s.GetByID(ctx, taskID, userID)
				if err != nil {
					resultsChan <- fmt.Errorf("task %s: %w", taskID, err)
					continue
				}

				// Convert taskID to int
				taskIDInt, err := strconv.Atoi(taskID)
				if err != nil {
					resultsChan <- fmt.Errorf("invalid task ID %s: %w", taskID, err)
					continue
				}

				// Update only status while preserving title and description
				task := &domain.Task{
					ID:          taskIDInt,
					UserID:      userIDInt,
					Title:       existingTask.Title,
					Description: existingTask.Description,
					Status:      domain.TaskStatusDone,
					UpdatedAt:   time.Now(),
				}

				if err := s.taskRepo.Update(ctx, task); err != nil {
					resultsChan <- fmt.Errorf("failed to update task %s: %w", taskID, err)
				} else {
					resultsChan <- nil
				}
			}
		}()
	}

	// Send task IDs to channel
	go func() {
		for _, taskID := range req.TaskIDs {
			taskIDsChan <- taskID
		}
		close(taskIDsChan)
	}()

	// Wait for all workers to complete
	wg.Wait()
	close(resultsChan)

	// Collect results
	successCount := 0
	failedIDs := []string{}

	for i, taskID := range req.TaskIDs {
		err := <-resultsChan
		if err != nil {
			failedIDs = append(failedIDs, taskID)
		} else {
			successCount++
		}
		// Ensure we process all results
		if i+1 < len(req.TaskIDs) {
			<-resultsChan
		}
	}

	return &dto.BulkCompleteResponse{
		SuccessCount: successCount,
		FailedCount:  len(failedIDs),
		FailedIDs:    failedIDs,
	}, nil
}
