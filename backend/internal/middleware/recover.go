package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// Recover converts panics to 500 responses and logs the error.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Ctx(r.Context()).Error().Interface("panic", rec).Msg("panic recovered")
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
