package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rishuishind/rivalTaskAssestment/backend/config"
	"github.com/rishuishind/rivalTaskAssestment/backend/models"
	"github.com/rishuishind/rivalTaskAssestment/backend/utils"
	"github.com/rishuishind/rivalTaskAssestment/backend/validators"
	"golang.org/x/crypto/bcrypt"
)

// AuthResponse is the response returned after successful auth.
type AuthResponse struct {
	Token string      `json:"token"`
	User  UserProfile `json:"user"`
}

// UserProfile is the public user profile (no password hash).
type UserProfile struct {
	ID        uuid.UUID   `json:"id"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	Role      models.Role `json:"role"`
	CreatedAt string      `json:"created_at"`
}

// toUserProfile converts a User model to a public profile.
func toUserProfile(user models.User) UserProfile {
	return UserProfile{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

// Signup handles POST /api/auth/signup
func Signup(c *gin.Context) {
	var req validators.SignupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if validationErrors := validators.ValidateSignup(req); len(validationErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    utils.ErrCodeValidation,
				"message": "Validation failed",
				"details": validationErrors,
			},
		})
		return
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Check if user already exists
	var existingUser models.User
	if result := config.DB.Where("email = ?", email).First(&existingUser); result.Error == nil {
		utils.ErrorResponse(c, http.StatusConflict, utils.ErrCodeConflict, "A user with this email already exists")
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to process password")
		return
	}

	// Create user
	user := models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         models.RoleUser,
	}

	if result := config.DB.Create(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, AuthResponse{
		Token: token,
		User:  toUserProfile(user),
	})
}

// Login handles POST /api/auth/login
func Login(c *gin.Context) {
	var req validators.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if validationErrors := validators.ValidateLogin(req); len(validationErrors) > 0 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"error": gin.H{
				"code":    utils.ErrCodeValidation,
				"message": "Validation failed",
				"details": validationErrors,
			},
		})
		return
	}

	// Normalize email
	email := strings.TrimSpace(strings.ToLower(req.Email))

	// Find user by email
	var user models.User
	if result := config.DB.Where("email = ?", email).First(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "Invalid email or password")
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, string(user.Role))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to generate token")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, AuthResponse{
		Token: token,
		User:  toUserProfile(user),
	})
}

// GetMe handles GET /api/auth/me
func GetMe(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, utils.ErrCodeUnauthorized, "User not authenticated")
		return
	}

	var user models.User
	if result := config.DB.Where("id = ?", userID).First(&user); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, utils.ErrCodeNotFound, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, toUserProfile(user))
}
