package validators

import (
	"fmt"
	"strings"
	"time"

	"github.com/rishuishind/rivalTaskAssestment/backend/models"
)

// CreateTaskRequest represents the request body for creating a task.
type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	DueDate     string `json:"due_date"`
}

// UpdateTaskRequest represents the request body for updating a task.
type UpdateTaskRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	Priority    *string `json:"priority"`
	DueDate     *string `json:"due_date"`
}

// ValidationError holds field-level validation errors.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidateCreateTask validates the create task request.
func ValidateCreateTask(req CreateTaskRequest) []ValidationError {
	var errors []ValidationError

	// Title is required and must be between 1 and 255 characters
	title := strings.TrimSpace(req.Title)
	if title == "" {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title is required",
		})
	} else if len(title) > 255 {
		errors = append(errors, ValidationError{
			Field:   "title",
			Message: "Title must be at most 255 characters",
		})
	}

	// Description is optional but has a max length
	if len(req.Description) > 5000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be at most 5000 characters",
		})
	}

	// Status validation (optional, defaults to "todo")
	if req.Status != "" && !models.IsValidStatus(req.Status) {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("Status must be one of: %s", strings.Join(statusStrings(), ", ")),
		})
	}

	// Priority validation (optional, defaults to "medium")
	if req.Priority != "" && !models.IsValidPriority(req.Priority) {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: fmt.Sprintf("Priority must be one of: %s", strings.Join(priorityStrings(), ", ")),
		})
	}

	// Due date validation (optional)
	if req.DueDate != "" {
		_, err := time.Parse(time.RFC3339, req.DueDate)
		if err != nil {
			errors = append(errors, ValidationError{
				Field:   "due_date",
				Message: "Due date must be in RFC3339 format (e.g., 2024-12-31T23:59:59Z)",
			})
		}
	}

	return errors
}

// ValidateUpdateTask validates the update task request.
func ValidateUpdateTask(req UpdateTaskRequest) []ValidationError {
	var errors []ValidationError

	// Title validation (if provided)
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			errors = append(errors, ValidationError{
				Field:   "title",
				Message: "Title cannot be empty",
			})
		} else if len(title) > 255 {
			errors = append(errors, ValidationError{
				Field:   "title",
				Message: "Title must be at most 255 characters",
			})
		}
	}

	// Description validation (if provided)
	if req.Description != nil && len(*req.Description) > 5000 {
		errors = append(errors, ValidationError{
			Field:   "description",
			Message: "Description must be at most 5000 characters",
		})
	}

	// Status validation (if provided)
	if req.Status != nil && !models.IsValidStatus(*req.Status) {
		errors = append(errors, ValidationError{
			Field:   "status",
			Message: fmt.Sprintf("Status must be one of: %s", strings.Join(statusStrings(), ", ")),
		})
	}

	// Priority validation (if provided)
	if req.Priority != nil && !models.IsValidPriority(*req.Priority) {
		errors = append(errors, ValidationError{
			Field:   "priority",
			Message: fmt.Sprintf("Priority must be one of: %s", strings.Join(priorityStrings(), ", ")),
		})
	}

	// Due date validation (if provided)
	if req.DueDate != nil && *req.DueDate != "" {
		_, err := time.Parse(time.RFC3339, *req.DueDate)
		if err != nil {
			errors = append(errors, ValidationError{
				Field:   "due_date",
				Message: "Due date must be in RFC3339 format (e.g., 2024-12-31T23:59:59Z)",
			})
		}
	}

	return errors
}

func statusStrings() []string {
	statuses := models.ValidStatuses()
	result := make([]string, len(statuses))
	for i, s := range statuses {
		result[i] = string(s)
	}
	return result
}

func priorityStrings() []string {
	priorities := models.ValidPriorities()
	result := make([]string, len(priorities))
	for i, p := range priorities {
		result[i] = string(p)
	}
	return result
}
