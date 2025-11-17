package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
)

type AuthHandler struct {
	db     *sql.DB
	cfg    *config.Config
	logger *zap.Logger
}

func NewAuthHandler(db *sql.DB, cfg *config.Config, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg, logger: logger}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Get user with role
	var user models.User
	query := `
		SELECT u.id, u.username, u.email, u.password_hash, u.role_id, r.name as role_name,
		       u.failed_login_attempts, u.account_locked_until
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.email = ? AND u.deleted_at IS NULL
	`

	err := h.db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.RoleID, &user.RoleName, &user.FailedLoginAttempts, &user.AccountLockedUntil,
	)

	if err == sql.ErrNoRows {
		h.auditLog("", "login_failed", "user", "", r.RemoteAddr)
		utils.RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err != nil {
		h.logger.Error("failed to query user", zap.Error(err))
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Check account lockout
	if user.AccountLockedUntil.Valid && user.AccountLockedUntil.Time.After(time.Now()) {
		remaining := time.Until(user.AccountLockedUntil.Time)
		utils.RespondError(w, http.StatusForbidden,
			"account is locked due to too many failed login attempts. Try again in "+remaining.Round(time.Minute).String())
		return
	}

	// Verify password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		// Increment failed attempts
		newAttempts := user.FailedLoginAttempts + 1
		var lockUntil sql.NullTime

		if newAttempts >= h.cfg.Security.MaxLoginAttempts {
			lockUntil = sql.NullTime{
				Time:  time.Now().Add(h.cfg.Security.LockoutDuration),
				Valid: true,
			}
		}

		_, _ = h.db.Exec(`
			UPDATE users 
			SET failed_login_attempts = ?, account_locked_until = ?
			WHERE id = ?
		`, newAttempts, lockUntil, user.ID)

		h.auditLog(user.ID, "login_failed", "user", user.ID, r.RemoteAddr)
		utils.RespondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	// Reset failed attempts and update last login
	_, err = h.db.Exec(`
		UPDATE users 
		SET failed_login_attempts = 0, account_locked_until = NULL, last_login_at = ?
		WHERE id = ?
	`, time.Now(), user.ID)

	if err != nil {
		h.logger.Error("failed to reset login attempts", zap.Error(err))
	}

	// Generate tokens
	accessToken, err := utils.GenerateAccessToken(
		user.ID, user.Username, user.Email, user.RoleID, user.RoleName,
		h.cfg.JWT.Secret, h.cfg.JWT.AccessTokenTTL,
	)
	if err != nil {
		h.logger.Error("failed to generate access token", zap.Error(err))
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, h.cfg.JWT.Secret, h.cfg.JWT.RefreshTokenTTL)
	if err != nil {
		h.logger.Error("failed to generate refresh token", zap.Error(err))
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	h.auditLog(user.ID, "login_success", "user", user.ID, r.RemoteAddr)

	utils.RespondJSON(w, http.StatusOK, models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user.ToResponse(),
	})
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		utils.RespondError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var dbUser models.User
	query := `
		SELECT u.id, u.username, u.email, u.role_id, r.name as role_name, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`

	err := h.db.QueryRow(query, user.UserID).Scan(
		&dbUser.ID, &dbUser.Username, &dbUser.Email, &dbUser.RoleID, &dbUser.RoleName,
		&dbUser.LastLoginAt, &dbUser.CreatedAt, &dbUser.UpdatedAt,
	)

	if err != nil {
		h.logger.Error("failed to get user", zap.Error(err))
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, dbUser.ToResponse())
}

func (h *AuthHandler) auditLog(userID, action, entityType, entityID, ipAddress string) {
	_, err := h.db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, userID, action, entityType, entityID, ipAddress)

	if err != nil {
		h.logger.Error("failed to create audit log", zap.Error(err))
	}
}
