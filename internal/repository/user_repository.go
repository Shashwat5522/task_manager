package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/vedologic/task-manager/internal/domain"
)

// userRepository implements UserRepository interface using raw SQL
type userRepository struct {
	db *sqlx.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// SQL Queries
const (
	queryCreateUser = `
		INSERT INTO users (id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	queryFindUserByEmail = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	queryFindUserByID = `
		SELECT id, email, password_hash, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	queryUserExists = `
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`
)

// Create creates a new user in the database
func (r *userRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := r.db.ExecContext(
		ctx,
		queryCreateUser,
		user.ID,
		user.Email,
		user.PasswordHash,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("no rows affected when creating user")
	}

	return nil
}

// FindByEmail finds a user by email address
func (r *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.GetContext(ctx, user, queryFindUserByEmail, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with email: %s", email)
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}

	return user, nil
}

// FindByID finds a user by ID
func (r *userRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	user := &domain.User{}

	err := r.db.GetContext(ctx, user, queryFindUserByID, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found with id: %s", id)
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	return user, nil
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool

	err := r.db.GetContext(ctx, &exists, queryUserExists, email)
	if err != nil {
		return false, fmt.Errorf("failed to check if user exists: %w", err)
	}

	return exists, nil
}
