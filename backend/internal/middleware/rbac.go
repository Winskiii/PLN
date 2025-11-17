package middleware

import (
	"net/http"

	"backend/internal/models"
	"backend/internal/utils"
)

func RequireRole(minRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := GetUserFromContext(r.Context())
			if !ok {
				utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			role := &models.Role{Name: user.RoleName}
			if !role.HasPermission(minRole) {
				utils.RespondError(w, http.StatusForbidden, "insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
