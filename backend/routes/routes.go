package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishuishind/rivalTaskAssestment/backend/handlers"
)

// SetupRoutes configures all API routes.
func SetupRoutes(router *gin.Engine) {
	// API group
	api := router.Group("/api")

	// Health check
	api.GET("/health", handlers.HealthCheck)

	// Task routes (will be protected by auth middleware in Phase 2)
	tasks := api.Group("/tasks")
	{
		tasks.POST("", handlers.CreateTask)
		tasks.GET("", handlers.GetTasks)
		tasks.GET("/:id", handlers.GetTask)
		tasks.PATCH("/:id", handlers.UpdateTask)
		tasks.DELETE("/:id", handlers.DeleteTask)
	}
}
