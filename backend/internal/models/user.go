package models

import (
	"database/sql"
	"time"
)

type User struct {
	ID                  string       `json:"id"`
	Username            string       `json:"username"`
	Email               string       `json:"email"`
	PasswordHash        string       `json:"-"`
	RoleID              string       `json:"role_id"`
	RoleName            string       `json:"role_name,omitempty"`
	FailedLoginAttempts int          `json:"-"`
	AccountLockedUntil  sql.NullTime `json:"-"`
	LastLoginAt         sql.NullTime `json:"last_login_at,omitempty"`
	CreatedAt           time.Time    `json:"created_at"`
	UpdatedAt           time.Time    `json:"updated_at"`
	DeletedAt           sql.NullTime `json:"-"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   string `json:"role_id"`
}

type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   string `json:"role_id"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserResponse struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	RoleID      string     `json:"role_id"`
	RoleName    string     `json:"role_name"`
	LastLoginAt *time.Time `json:"last_login_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

func (u *User) ToResponse() *UserResponse {
	resp := &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		RoleID:    u.RoleID,
		RoleName:  u.RoleName,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.LastLoginAt.Valid {
		resp.LastLoginAt = &u.LastLoginAt.Time
	}
	return resp
}
