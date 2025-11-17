package models

import (
	"database/sql"
	"time"
)

type Project struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Status      string       `json:"status"`
	OwnerID     string       `json:"owner_id"`
	OwnerName   string       `json:"owner_name,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	DeletedAt   sql.NullTime `json:"-"`
}

type CreateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

const (
	ProjectStatusActive   = "Active"
	ProjectStatusInactive = "Inactive"
	ProjectStatusArchived = "Archived"
)
