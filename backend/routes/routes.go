package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/rishuishind/rivalTaskAssestment/backend/handlers"
	"github.com/rishuishind/rivalTaskAssestment/backend/middleware"
)

// SetupRoutes configures all API routes.
func SetupRoutes(router *gin.Engine) {
	// API group
	api := router.Group("/api")

	// Health check (public)
	api.GET("/health", handlers.HealthCheck)

	// Auth routes (public)
	auth := api.Group("/auth")
	{
		auth.POST("/signup", handlers.Signup)
		auth.POST("/login", handlers.Login)
		auth.GET("/me", middleware.AuthMiddleware(), handlers.GetMe)
	}

	// Task routes (protected — requires authentication)
	tasks := api.Group("/tasks")
	tasks.Use(middleware.AuthMiddleware())
	{
		tasks.POST("", handlers.CreateTask)
		tasks.GET("", handlers.GetTasks)
		tasks.GET("/:id", handlers.GetTask)
		tasks.PATCH("/:id", handlers.UpdateTask)
		tasks.DELETE("/:id", handlers.DeleteTask)
	}
}
