package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// TaskStatus represents the status of a task.
type TaskStatus string

const (
	StatusTodo       TaskStatus = "todo"
	StatusInProgress TaskStatus = "in_progress"
	StatusDone       TaskStatus = "done"
)

// TaskPriority represents the priority level of a task.
type TaskPriority string

const (
	PriorityLow    TaskPriority = "low"
	PriorityMedium TaskPriority = "medium"
	PriorityHigh   TaskPriority = "high"
	PriorityUrgent TaskPriority = "urgent"
)

// Task represents a task in the system.
type Task struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Description string         `gorm:"type:text" json:"description"`
	Status      TaskStatus     `gorm:"type:varchar(20);default:'todo';not null" json:"status"`
	Priority    TaskPriority   `gorm:"type:varchar(20);default:'medium';not null" json:"priority"`
	DueDate     *time.Time     `gorm:"type:timestamp" json:"due_date"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	User        User           `gorm:"foreignKey:UserID" json:"-"`
}

// TableName overrides the table name.
func (Task) TableName() string {
	return "tasks"
}

// ValidStatuses returns all valid task statuses.
func ValidStatuses() []TaskStatus {
	return []TaskStatus{StatusTodo, StatusInProgress, StatusDone}
}

// ValidPriorities returns all valid task priorities.
func ValidPriorities() []TaskPriority {
	return []TaskPriority{PriorityLow, PriorityMedium, PriorityHigh, PriorityUrgent}
}

// IsValidStatus checks if a given status string is valid.
func IsValidStatus(s string) bool {
	for _, status := range ValidStatuses() {
		if string(status) == s {
			return true
		}
	}
	return false
}

// IsValidPriority checks if a given priority string is valid.
func IsValidPriority(p string) bool {
	for _, priority := range ValidPriorities() {
		if string(priority) == p {
			return true
		}
	}
	return false
}
