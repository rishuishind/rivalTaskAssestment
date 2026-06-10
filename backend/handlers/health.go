package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rishuishind/rivalTaskAssestment/backend/config"
	"github.com/rishuishind/rivalTaskAssestment/backend/utils"
)

// HealthCheck returns the health status of the API.
func HealthCheck(c *gin.Context) {
	// Check database connection
	sqlDB, err := config.DB.DB()
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Database connection error")
		return
	}

	if err := sqlDB.Ping(); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Database ping failed")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "task-manager-api",
		"version": "1.0.0",
	})
}
