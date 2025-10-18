package handler

import (
	"fmt"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vedologic/task-manager/internal/middleware"
	"go.uber.org/zap"
)

// SetupRoutes configures all API routes and middleware
func SetupRoutes(
	router *gin.Engine,
	authHandler *AuthHandler,
	taskHandler *TaskHandler,
	jwtSecret string,
	log *zap.Logger,
) {
	// Apply global middleware
	router.Use(middleware.LoggerMiddleware(log))
	router.Use(middleware.RecoveryMiddleware(log))

	// Swagger UI
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/docs", func(c *gin.Context) {
		c.Redirect(301, "/swagger/index.html")
	})

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Public routes - Auth
	authRoutes := router.Group("/api/v1/auth")
	{
		authRoutes.POST("/register", authHandler.Register)
		authRoutes.POST("/login", authHandler.Login)
	}

	// Protected routes - Tasks
	taskRoutes := router.Group("/api/v1/tasks")
	taskRoutes.Use(middleware.AuthMiddleware(jwtSecret))
	{
		taskRoutes.POST("", taskHandler.Create)
		taskRoutes.GET("", taskHandler.List)
		taskRoutes.GET("/:id", taskHandler.GetByID)
		taskRoutes.PUT("/:id", taskHandler.Update)
		taskRoutes.DELETE("/:id", taskHandler.Delete)
		taskRoutes.PATCH("/bulk-complete", taskHandler.BulkComplete)
	}

	log.Info("Routes configured successfully")
	printRegisteredRoutes(router, log)
}

// printRegisteredRoutes logs all registered routes
func printRegisteredRoutes(router *gin.Engine, log *zap.Logger) {
	log.Info("Registered Routes:")
	separator := "============================================================"
	fmt.Println("\n" + separator)
	fmt.Println("REGISTERED ROUTES")
	fmt.Println(separator)

	for _, route := range router.Routes() {
		method := route.Method
		path := route.Path
		fmt.Printf("[%-6s] %s\n", method, path)
		log.Debug("Route registered", zap.String("method", method), zap.String("path", path))
	}

	fmt.Println(separator + "\n")
	log.Info("All routes loaded and ready!")
}
