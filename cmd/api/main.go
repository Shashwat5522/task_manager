package main

import (
	"fmt"
	stdlog "log"

	"github.com/vedologic/task-manager/config"
	"github.com/vedologic/task-manager/internal/repository"
	"github.com/vedologic/task-manager/pkg/database"
	"github.com/vedologic/task-manager/pkg/logger"
)

func main() {
	// Load configuration from environment variables and .env file
	cfg, err := config.Load()
	if err != nil {
		stdlog.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		stdlog.Fatalf("Failed to initialize logger: %v", err)
	}
	defer log.Sync()

	// Log startup information
	log.Info("Task Manager API starting up")

	log.Info("Configuration loaded")

	log.Info(fmt.Sprintf("Server configured to run on %s:%s", cfg.Server.Host, cfg.Server.Port))
	log.Info(fmt.Sprintf("Environment: %s", cfg.Server.Env))
	log.Info(fmt.Sprintf("Log Level: %s", cfg.Log.Level))

	// Initialize database connection
	log.Info("Initializing database connection...")
	dbConfig := database.Config{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		User:            cfg.Database.User,
		Password:        cfg.Database.Password,
		DBName:          cfg.Database.DBName,
		SSLMode:         cfg.Database.SSLMode,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		stdlog.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close(db)

	log.Info("Database connection established successfully")
	log.Info(fmt.Sprintf("Database: %s:%d/%s", cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName))
	log.Info(fmt.Sprintf("Connection Pool - Max Open: %d, Max Idle: %d", cfg.Database.MaxOpenConns, cfg.Database.MaxIdleConns))

	// Initialize repositories
	log.Info("Initializing repositories...")
	_ = repository.NewUserRepository(db)
	_ = repository.NewTaskRepository(db)
	log.Info("Repositories initialized successfully")
	log.Info("User Repository: ready")
	log.Info("Task Repository: ready")

	// TODO: Initialize services
	// TODO: Initialize handlers
	// TODO: Setup routes and middleware
	// TODO: Start HTTP server
	// TODO: Handle graceful shutdown

	log.Info("Application initialization complete")

	// Keep application running
	log.Info("Task Manager API is ready to serve requests")
	select {}
}
