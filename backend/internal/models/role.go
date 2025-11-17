package models

import "time"

type Role struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

const (
	RoleSysAdmin = "SYSADMIN"
	RoleAdmin    = "ADMIN"
	RoleManager  = "MANAGER"
	RoleEmployee = "EMPLOYEE"
)

func (r *Role) HasPermission(requiredRole string) bool {
	roles := map[string]int{
		RoleEmployee: 1,
		RoleManager:  2,
		RoleAdmin:    3,
		RoleSysAdmin: 4,
	}

	userLevel := roles[r.Name]
	requiredLevel := roles[requiredRole]

	return userLevel >= requiredLevel
}
