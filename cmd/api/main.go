package main

import (
	"context"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vedologic/task-manager/config"
	"github.com/vedologic/task-manager/internal/handler"
	"github.com/vedologic/task-manager/internal/repository"
	"github.com/vedologic/task-manager/internal/service"
	"github.com/vedologic/task-manager/pkg/database"
	"github.com/vedologic/task-manager/pkg/logger"

	_ "github.com/vedologic/task-manager/docs"
)

// @title Task Manager API
// @version 1.0
// @description Task Manager API for managing tasks and users.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
// @schemes http https
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description JWT Bearer token authentication. Use "Bearer <token>"
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

	// Run automatic database migrations
	log.Info("Executing database migrations...")
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	migrationManager := database.NewMigrationManager(dbURL, log.Logger)

	// Check and run pending migrations
	if err := migrationManager.RunMigrationsIfNeeded(); err != nil {
		stdlog.Fatalf("Failed to run migrations: %v", err)
	}

	// Verify schema integrity after migrations
	if err := migrationManager.VerifySchema(db, log.Logger); err != nil {
		stdlog.Fatalf("Database schema verification failed: %v", err)
	}

	log.Info("All migrations completed and schema verified successfully")

	// Initialize repositories
	log.Info("Initializing repositories...")
	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	log.Info("Repositories initialized successfully")
	log.Info("User Repository: ready")
	log.Info("Task Repository: ready")

	// Initialize services
	log.Info("Initializing services...")
	authService := service.NewAuthService(userRepo, cfg.JWT.Secret, cfg.JWT.ExpiryHours)
	taskService := service.NewTaskService(taskRepo)
	log.Info("Services initialized successfully")
	log.Info("Auth Service: ready")
	log.Info("Task Service: ready")

	// Initialize handlers
	log.Info("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService, log.Logger)
	taskHandler := handler.NewTaskHandler(taskService, log.Logger)
	log.Info("Handlers initialized successfully")
	log.Info("Auth Handler: ready")
	log.Info("Task Handler: ready")

	// Setup router and routes
	log.Info("Setting up routes and middleware...")

	// Convert environment to Gin mode (Gin only accepts: debug, release, test)
	ginMode := cfg.Server.Env
	if ginMode == "development" {
		ginMode = "debug"
	}
	gin.SetMode(ginMode)

	router := gin.New()
	handler.SetupRoutes(router, authHandler, taskHandler, cfg.JWT.Secret, log.Logger)
	log.Info("Routes and middleware configured successfully")

	// Setup HTTP server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Info(fmt.Sprintf("Starting HTTP server on %s", addr))

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			stdlog.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Info("Task Manager API is ready to serve requests on http://" + addr)
	log.Info("Swagger UI available at http://" + addr + "/swagger/index.html")
	log.Info("Health check available at http://" + addr + "/health")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutdown signal received, shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("Error during graceful shutdown: %v", err))
	}

	log.Info("Task Manager API shut down successfully")
}
