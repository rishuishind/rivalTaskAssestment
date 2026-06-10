package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rishuishind/rivalTaskAssestment/backend/utils"
)

// AuthMiddleware validates the JWT token and sets user info in context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "Authorization header must be in format: Bearer <token>")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		// Set user info in context for downstream handlers
		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)
		c.Set("userRole", claims.Role)

		c.Next()
	}
}

// AdminMiddleware checks if the authenticated user has admin role.
// Must be used AFTER AuthMiddleware.
func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists || role.(string) != string("admin") {
			utils.ErrorResponse(c, http.StatusForbidden, utils.ErrCodeForbidden, "Admin access required")
			c.Abort()
			return
		}

		c.Next()
	}
}
