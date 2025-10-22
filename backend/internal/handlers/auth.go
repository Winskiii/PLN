package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"

	appcfg "backend/internal/config"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type loginResponse struct {
	Token string `json:"token"`
}

func Login(cfg *appcfg.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Ctx(r.Context()).Warn().Err(err).Msg("read login body")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if err := json.Unmarshal(b, &req); err != nil {
			log.Ctx(r.Context()).Warn().Err(err).RawJSON("body", b).Msg("decode login body")
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		if req.Username != cfg.DemoUser || req.Password != cfg.DemoPass {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		claims := jwt.MapClaims{
			"sub": req.Username,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(cfg.JWTTTL).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := token.SignedString([]byte(cfg.JWTSecret))
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		respondJSON(w, http.StatusOK, loginResponse{Token: signed})
	}
}

func Me() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Prefer header if present
		if user := r.Header.Get("X-User"); user != "" {
			respondJSON(w, http.StatusOK, map[string]string{"user": user})
			return
		}
		// Fallback to context
		if v := r.Context().Value(struct{ k string }{k: "user"}); v != nil {
			if user, ok := v.(string); ok && user != "" {
				respondJSON(w, http.StatusOK, map[string]string{"user": user})
				return
			}
		}
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	}
}
