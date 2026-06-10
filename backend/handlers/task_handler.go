package handlers

import (
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rishuishind/rivalTaskAssestment/backend/config"
	"github.com/rishuishind/rivalTaskAssestment/backend/models"
	"github.com/rishuishind/rivalTaskAssestment/backend/utils"
	"github.com/rishuishind/rivalTaskAssestment/backend/validators"
)

// CreateTask handles POST /api/tasks
func CreateTask(c *gin.Context) {
	var req validators.CreateTaskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if validationErrors := validators.ValidateCreateTask(req); len(validationErrors) > 0 {
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

	// Build the task
	task := models.Task{
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		Status:      models.StatusTodo,
		Priority:    models.PriorityMedium,
	}

	// Set status if provided
	if req.Status != "" {
		task.Status = models.TaskStatus(req.Status)
	}

	// Set priority if provided
	if req.Priority != "" {
		task.Priority = models.TaskPriority(req.Priority)
	}

	// Parse and set due date if provided
	if req.DueDate != "" {
		dueDate, _ := time.Parse(time.RFC3339, req.DueDate)
		task.DueDate = &dueDate
	}

	// Get user ID from context (set by auth middleware — for now use a placeholder)
	userID, exists := c.Get("userID")
	if exists {
		task.UserID = userID.(uuid.UUID)
	}

	// Create in database
	if result := config.DB.Create(&task); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to create task")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, task)
}

// GetTasks handles GET /api/tasks
func GetTasks(c *gin.Context) {
	// Parse pagination params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	// Build query
	query := config.DB.Model(&models.Task{})

	// Get user ID from context (set by auth middleware)
	userID, exists := c.Get("userID")
	if exists {
		query = query.Where("user_id = ?", userID)
	}

	// Filter by status
	if statusFilter := c.Query("status"); statusFilter != "" {
		statuses := strings.Split(statusFilter, ",")
		validStatuses := []string{}
		for _, s := range statuses {
			s = strings.TrimSpace(s)
			if models.IsValidStatus(s) {
				validStatuses = append(validStatuses, s)
			}
		}
		if len(validStatuses) > 0 {
			query = query.Where("status IN ?", validStatuses)
		}
	}

	// Filter by priority
	if priorityFilter := c.Query("priority"); priorityFilter != "" {
		priorities := strings.Split(priorityFilter, ",")
		validPriorities := []string{}
		for _, p := range priorities {
			p = strings.TrimSpace(p)
			if models.IsValidPriority(p) {
				validPriorities = append(validPriorities, p)
			}
		}
		if len(validPriorities) > 0 {
			query = query.Where("priority IN ?", validPriorities)
		}
	}

	// Search by title
	if search := c.Query("search"); search != "" {
		search = strings.TrimSpace(search)
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Sorting
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	// Validate sort field
	allowedSortFields := map[string]bool{
		"created_at": true,
		"updated_at": true,
		"due_date":   true,
		"priority":   true,
		"title":      true,
		"status":     true,
	}
	if !allowedSortFields[sortBy] {
		sortBy = "created_at"
	}
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// For priority sorting, use a CASE expression to sort by severity
	var orderClause string
	if sortBy == "priority" {
		if sortOrder == "asc" {
			orderClause = "CASE priority WHEN 'low' THEN 1 WHEN 'medium' THEN 2 WHEN 'high' THEN 3 WHEN 'urgent' THEN 4 END ASC"
		} else {
			orderClause = "CASE priority WHEN 'urgent' THEN 1 WHEN 'high' THEN 2 WHEN 'medium' THEN 3 WHEN 'low' THEN 4 END ASC"
		}
	} else {
		orderClause = sortBy + " " + sortOrder
	}

	// Fetch tasks
	var tasks []models.Task
	result := query.Order(orderClause).Offset(offset).Limit(limit).Find(&tasks)
	if result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to fetch tasks")
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	utils.PaginatedResponse(c, tasks, utils.APIMeta{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetTask handles GET /api/tasks/:id
func GetTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	query := config.DB.Where("id = ?", taskID)

	// Scope to current user
	userID, exists := c.Get("userID")
	if exists {
		query = query.Where("user_id = ?", userID)
	}

	if result := query.First(&task); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, utils.ErrCodeNotFound, "Task not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, task)
}

// UpdateTask handles PATCH /api/tasks/:id
func UpdateTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid task ID")
		return
	}

	// Find existing task
	var task models.Task
	query := config.DB.Where("id = ?", taskID)

	// Scope to current user
	userID, exists := c.Get("userID")
	if exists {
		query = query.Where("user_id = ?", userID)
	}

	if result := query.First(&task); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, utils.ErrCodeNotFound, "Task not found")
		return
	}

	// Parse request body
	var req validators.UpdateTaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid request body")
		return
	}

	// Validate input
	if validationErrors := validators.ValidateUpdateTask(req); len(validationErrors) > 0 {
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

	// Build updates map (only update provided fields)
	updates := map[string]interface{}{}

	if req.Title != nil {
		updates["title"] = strings.TrimSpace(*req.Title)
	}
	if req.Description != nil {
		updates["description"] = strings.TrimSpace(*req.Description)
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.DueDate != nil {
		if *req.DueDate == "" {
			updates["due_date"] = nil
		} else {
			dueDate, _ := time.Parse(time.RFC3339, *req.DueDate)
			updates["due_date"] = &dueDate
		}
	}

	if len(updates) == 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "No fields to update")
		return
	}

	// Update in database
	if result := config.DB.Model(&task).Updates(updates); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to update task")
		return
	}

	// Reload the task to get updated fields
	config.DB.First(&task, "id = ?", taskID)

	utils.SuccessResponse(c, http.StatusOK, task)
}

// DeleteTask handles DELETE /api/tasks/:id (soft delete)
func DeleteTask(c *gin.Context) {
	taskID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, utils.ErrCodeBadRequest, "Invalid task ID")
		return
	}

	// Find existing task
	var task models.Task
	query := config.DB.Where("id = ?", taskID)

	// Scope to current user
	userID, exists := c.Get("userID")
	if exists {
		query = query.Where("user_id = ?", userID)
	}

	if result := query.First(&task); result.Error != nil {
		utils.ErrorResponse(c, http.StatusNotFound, utils.ErrCodeNotFound, "Task not found")
		return
	}

	// Soft delete
	if result := config.DB.Delete(&task); result.Error != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, utils.ErrCodeInternal, "Failed to delete task")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{
		"message": "Task deleted successfully",
	})
}
