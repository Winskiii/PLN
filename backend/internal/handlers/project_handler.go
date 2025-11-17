package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
)

type ProjectHandler struct {
	db  *sql.DB
	cfg *config.Config
}

func NewProjectHandler(db *sql.DB, cfg *config.Config) *ProjectHandler {
	return &ProjectHandler{db: db, cfg: cfg}
}

func (h *ProjectHandler) List(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	// Show all projects for ADMIN+, only owned projects for others
	query := `
		SELECT p.id, p.name, p.description, p.status, p.owner_id, u.username as owner_name, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN users u ON p.owner_id = u.id
		WHERE p.deleted_at IS NULL
	`

	// Filter for non-admin users
	role := &models.Role{Name: user.RoleName}
	if !role.HasPermission(models.RoleAdmin) {
		query += " AND p.owner_id = ?"
	}

	query += " ORDER BY p.created_at DESC"

	var rows *sql.Rows
	var err error

	if role.HasPermission(models.RoleAdmin) {
		rows, err = h.db.Query(query)
	} else {
		rows, err = h.db.Query(query, user.UserID)
	}

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to query projects")
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch projects")
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		err := rows.Scan(
			&project.ID, &project.Name, &project.Description, &project.Status,
			&project.OwnerID, &project.OwnerName, &project.CreatedAt, &project.UpdatedAt,
		)
		if err != nil {
			log.Error().Err(err).Msg("failed to scan project")
			continue
		}
		projects = append(projects, project)
	}

	utils.RespondJSON(w, http.StatusOK, projects)
}

func (h *ProjectHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var project models.Project
	query := `
		SELECT p.id, p.name, p.description, p.status, p.owner_id, u.username as owner_name, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN users u ON p.owner_id = u.id
		WHERE p.id = ? AND p.deleted_at IS NULL
	`

	err := h.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.OwnerID, &project.OwnerName, &project.CreatedAt, &project.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to get project")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	var req models.CreateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "project name is required")
		return
	}

	// Sanitize
	req.Name = utils.SanitizeString(req.Name)
	req.Description = utils.SanitizeString(req.Description)

	if req.Status == "" {
		req.Status = models.ProjectStatusActive
	}

	// Insert project
	id := uuid.NewString()
	_, err := h.db.Exec(`
		INSERT INTO projects (id, name, description, status, owner_id)
		VALUES (?, ?, ?, ?, ?)
	`, id, req.Name, req.Description, req.Status, user.UserID)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to create project")
		utils.RespondError(w, http.StatusInternalServerError, "failed to create project")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "project_created", "project", id, r.RemoteAddr)

	// Fetch created project
	var project models.Project
	query := `
		SELECT p.id, p.name, p.description, p.status, p.owner_id, u.username as owner_name, p.created_at, p.updated_at
		FROM projects p
		INNER JOIN users u ON p.owner_id = u.id
		WHERE p.id = ?
	`

	err = h.db.QueryRow(query, id).Scan(
		&project.ID, &project.Name, &project.Description, &project.Status,
		&project.OwnerID, &project.OwnerName, &project.CreatedAt, &project.UpdatedAt,
	)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to fetch created project")
		utils.RespondError(w, http.StatusInternalServerError, "project created but failed to fetch")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdateProjectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "project name is required")
		return
	}

	// Sanitize
	req.Name = utils.SanitizeString(req.Name)
	req.Description = utils.SanitizeString(req.Description)

	// Check permission (owner or admin)
	var ownerID string
	err := h.db.QueryRow("SELECT owner_id FROM projects WHERE id = ? AND deleted_at IS NULL", id).Scan(&ownerID)
	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	role := &models.Role{Name: user.RoleName}
	if ownerID != user.UserID && !role.HasPermission(models.RoleAdmin) {
		utils.RespondError(w, http.StatusForbidden, "you don't have permission to update this project")
		return
	}

	// Update project
	_, err = h.db.Exec(`
		UPDATE projects
		SET name = ?, description = ?, status = ?, updated_at = NOW()
		WHERE id = ?
	`, req.Name, req.Description, req.Status, id)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to update project")
		utils.RespondError(w, http.StatusInternalServerError, "failed to update project")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "project_updated", "project", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "project updated successfully")
}

func (h *ProjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	// Soft delete
	result, err := h.db.Exec("UPDATE projects SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to delete project")
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete project")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		utils.RespondError(w, http.StatusNotFound, "project not found")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "project_deleted", "project", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "project deleted successfully")
}

func (h *ProjectHandler) auditLog(userID, action, entityType, entityID, ipAddress string) {
	_, err := h.db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, userID, action, entityType, entityID, ipAddress)

	if err != nil {
		log.Error().Err(err).Msg("failed to create audit log")
	}
}
