package middleware

import (
	"context"
	"net/http"
	"strings"

	"backend/internal/utils"
)

type ctxKey string

const userKey ctxKey = "user"

// AuthJWT validates a Bearer JWT using project utils and stores claims in context.
func AuthJWT(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			// Validate and parse into Claims
			claims, err := utils.ValidateToken(tokenStr, secret)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			// Store claims in context and expose a header for compatibility
			ctx := context.WithValue(r.Context(), userKey, claims)
			r = r.WithContext(ctx)
			if claims != nil && claims.UserID != "" {
				w.Header().Set("X-User", claims.UserID)
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserFromContext returns parsed token claims stored by AuthJWT.
func GetUserFromContext(ctx context.Context) (*utils.Claims, bool) {
	if ctx == nil {
		return nil, false
	}
	v := ctx.Value(userKey)
	if v == nil {
		return nil, false
	}
	if claims, ok := v.(*utils.Claims); ok {
		return claims, true
	}
	return nil, false
}
