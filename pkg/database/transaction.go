package database

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// TxFunc is a function type that operates within a transaction
type TxFunc func(*sqlx.Tx) error

// WithTransaction executes a function within a database transaction
// It handles commit/rollback automatically based on the function's return value
func WithTransaction(ctx context.Context, db *sqlx.DB, fn TxFunc) error {
	// Begin transaction
	tx, err := db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Execute the function with deferred rollback/commit
	defer func() {
		if p := recover(); p != nil {
			// Panic occurred, rollback transaction
			tx.Rollback()
			panic(p)
		}
	}()

	// Execute the function
	if err := fn(tx); err != nil {
		// Error occurred, rollback transaction
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("failed to rollback transaction: %w, original error: %v", rbErr, err)
		}
		return err
	}

	// No error, commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
