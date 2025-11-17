package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID             string         `json:"id"`
	ProjectID      string         `json:"project_id"`
	ProjectName    string         `json:"project_name,omitempty"`
	Title          string         `json:"title"`
	Description    string         `json:"description"`
	Status         string         `json:"status"`
	Priority       string         `json:"priority"`
	AssignedTo     sql.NullString `json:"assigned_to"`
	AssignedToName string         `json:"assigned_to_name,omitempty"`
	DueDate        sql.NullTime   `json:"due_date"`
	CreatedBy      string         `json:"created_by"`
	CreatedByName  string         `json:"created_by_name,omitempty"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      sql.NullTime   `json:"-"`
}

type CreateTaskRequest struct {
	ProjectID   string `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  string `json:"assigned_to"`
	DueDate     string `json:"due_date"`
}

type UpdateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	Priority    string `json:"priority"`
	AssignedTo  string `json:"assigned_to"`
	DueDate     string `json:"due_date"`
}

const (
	TaskStatusNotStarted = "NotStarted"
	TaskStatusInProgress = "InProgress"
	TaskStatusCompleted  = "Completed"
	TaskStatusOnHold     = "OnHold"
	TaskStatusCancelled  = "Cancelled"

	TaskPriorityLow      = "Low"
	TaskPriorityMedium   = "Medium"
	TaskPriorityHigh     = "High"
	TaskPriorityCritical = "Critical"
)
