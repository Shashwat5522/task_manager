package database

import (
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type MigrationManager struct {
	dbURL string
	log   *zap.Logger
}

// NewMigrationManager creates a new migration manager instance
func NewMigrationManager(dbURL string, log *zap.Logger) *MigrationManager {
	return &MigrationManager{
		dbURL: dbURL,
		log:   log,
	}
}

// RunMigrationsIfNeeded runs pending migrations if the schema doesn't exist or is incomplete
func (m *MigrationManager) RunMigrationsIfNeeded() error {
	m.log.Info("Starting migration check...")

	// Get the absolute path to migrations directory
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	migrationSourceURL := fmt.Sprintf("file://%s", migrationsPath)
	m.log.Debug("Migration source URL", zap.String("path", migrationSourceURL))

	migrator, err := migrate.New(migrationSourceURL, m.dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	// Get current version
	version, dirty, err := migrator.Version()
	if err != nil {
		if err == migrate.ErrNilVersion {
			m.log.Info("No migrations applied yet, running all migrations...")
		} else {
			return fmt.Errorf("failed to get migration version: %w", err)
		}
	} else {
		m.log.Info("Current migration version", zap.Uint("version", version), zap.Bool("dirty", dirty))
	}

	// Run pending migrations
	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	if err == migrate.ErrNoChange {
		m.log.Info("No new migrations to apply")
	} else {
		m.log.Info("Migrations applied successfully")
	}

	return nil
}

// GetMigrationStatus returns the current migration status
func (m *MigrationManager) GetMigrationStatus() (uint, bool, error) {
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return 0, false, fmt.Errorf("failed to get migrations path: %w", err)
	}

	migrationSourceURL := fmt.Sprintf("file://%s", migrationsPath)
	migrator, err := migrate.New(migrationSourceURL, m.dbURL)
	if err != nil {
		return 0, false, fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Version()
}

// DownMigrations rollback all migrations
func (m *MigrationManager) DownMigrations() error {
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	migrationSourceURL := fmt.Sprintf("file://%s", migrationsPath)
	migrator, err := migrate.New(migrationSourceURL, m.dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Down()
}

// ForceVersion sets the migration version without running migrations
func (m *MigrationManager) ForceVersion(version int) error {
	migrationsPath, err := filepath.Abs("migrations")
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	migrationSourceURL := fmt.Sprintf("file://%s", migrationsPath)
	migrator, err := migrate.New(migrationSourceURL, m.dbURL)
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	return migrator.Force(version)
}

// CheckTablesExist checks if required tables exist
func (m *MigrationManager) CheckTablesExist(db *sqlx.DB) (bool, error) {
	tables := []string{"users", "tasks"}
	for _, table := range tables {
		query := fmt.Sprintf(`SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = '%s'`, table)
		var exists int64
		err := db.QueryRow(query).Scan(&exists)
		if err != nil || exists == 0 {
			return false, nil
		}
	}
	return true, nil
}

// VerifySchema verifies database schema integrity
// Non-blocking implementation: Only critical tables trigger errors
// Indices are optional - they will be added in future migrations if needed
func (m *MigrationManager) VerifySchema(db *sqlx.DB, log *zap.Logger) error {
	log.Info("Verifying database schema integrity...")

	// Critical checks that must exist for the app to function
	criticalChecks := []struct {
		name  string
		query string
	}{
		{name: "users_table", query: `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'users'`},
		{name: "tasks_table", query: `SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = 'public' AND table_name = 'tasks'`},
		{name: "task_status_enum", query: `SELECT COUNT(*) FROM pg_type WHERE typname = 'task_status'`},
	}

	// Verify critical components (blocking)
	for _, check := range criticalChecks {
		var exists int64
		err := db.QueryRow(check.query).Scan(&exists)
		if err != nil {
			log.Error("Critical schema verification failed", zap.String("check", check.name), zap.Error(err))
			return fmt.Errorf("critical schema verification failed for %s: %w", check.name, err)
		}
		if exists == 0 {
			log.Error("Missing critical schema component", zap.String("component", check.name))
			return fmt.Errorf("missing critical schema component: %s", check.name)
		}
		log.Debug("Critical schema component verified", zap.String("component", check.name))
	}

	log.Info("Database schema verification completed successfully")
	return nil
}
