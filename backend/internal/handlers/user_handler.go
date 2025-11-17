package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
)

type UserHandler struct {
	db  *sql.DB
	cfg *config.Config
}

func NewUserHandler(db *sql.DB, cfg *config.Config) *UserHandler {
	return &UserHandler{db: db, cfg: cfg}
}

func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT u.id, u.username, u.email, u.role_id, r.name as role_name, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.deleted_at IS NULL
		ORDER BY u.created_at DESC
	`

	rows, err := h.db.Query(query)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to query users")
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch users")
		return
	}
	defer rows.Close()

	var users []models.UserResponse
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.RoleID, &user.RoleName,
			&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan user")
			continue
		}
		users = append(users, *user.ToResponse())
	}

	utils.RespondJSON(w, http.StatusOK, users)
}

func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var user models.User
	query := `
		SELECT u.id, u.username, u.email, u.role_id, r.name as role_name, u.last_login_at, u.created_at, u.updated_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.id = ? AND u.deleted_at IS NULL
	`

	err := h.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.RoleID, &user.RoleName,
		&user.LastLoginAt, &user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to get user")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, user.ToResponse())
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := middleware.GetUserFromContext(r.Context())

	var req models.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateUsername(req.Username); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidatePassword(req.Password); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Sanitize
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeString(req.Email)

	// Check if email exists
	var exists int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE email = ? AND deleted_at IS NULL", req.Email).Scan(&exists)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to check email")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if exists > 0 {
		utils.RespondError(w, http.StatusConflict, "email already exists")
		return
	}

	// Check if username exists
	err = h.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND deleted_at IS NULL", req.Username).Scan(&exists)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to check username")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	if exists > 0 {
		utils.RespondError(w, http.StatusConflict, "username already exists")
		return
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password, bcrypt.DefaultCost)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to hash password")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	// Insert user
	id := uuid.NewString()
	_, err = h.db.Exec(`
		INSERT INTO users (id, username, email, password_hash, role_id)
		VALUES (?, ?, ?, ?, ?)
	`, id, req.Username, req.Email, passwordHash, req.RoleID)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to create user")
		utils.RespondError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	// Audit log
	h.auditLog(currentUser.UserID, "user_created", "user", id, r.RemoteAddr)

	// Fetch created user
	var user models.User
	query := `
		SELECT u.id, u.username, u.email, u.role_id, r.name as role_name, u.created_at, u.updated_at
		FROM users u
		INNER JOIN roles r ON u.role_id = r.id
		WHERE u.id = ?
	`

	err = h.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.RoleID, &user.RoleName,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to fetch created user")
		utils.RespondError(w, http.StatusInternalServerError, "user created but failed to fetch")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, user.ToResponse())
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate input
	if err := utils.ValidateUsername(req.Username); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := utils.ValidateEmail(req.Email); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Sanitize
	req.Username = utils.SanitizeString(req.Username)
	req.Email = utils.SanitizeString(req.Email)

	// Check if user exists
	var exists int
	err := h.db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ? AND deleted_at IS NULL", id).Scan(&exists)
	if err != nil || exists == 0 {
		utils.RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	// Update user
	_, err = h.db.Exec(`
		UPDATE users
		SET username = ?, email = ?, role_id = ?, updated_at = NOW()
		WHERE id = ?
	`, req.Username, req.Email, req.RoleID, id)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to update user")
		utils.RespondError(w, http.StatusInternalServerError, "failed to update user")
		return
	}

	// Audit log
	h.auditLog(currentUser.UserID, "user_updated", "user", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "user updated successfully")
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	currentUser, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	// Soft delete
	result, err := h.db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to delete user")
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete user")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		utils.RespondError(w, http.StatusNotFound, "user not found")
		return
	}

	// Audit log
	h.auditLog(currentUser.UserID, "user_deleted", "user", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "user deleted successfully")
}

func (h *UserHandler) auditLog(userID, action, entityType, entityID, ipAddress string) {
	_, err := h.db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, userID, action, entityType, entityID, ipAddress)

	if err != nil {
		log.Error().Err(err).Msg("failed to create audit log")
	}
}
