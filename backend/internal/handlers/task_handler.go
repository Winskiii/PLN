package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"backend/internal/config"
	"backend/internal/middleware"
	"backend/internal/models"
	"backend/internal/utils"
)

type TaskHandler struct {
	db  *sql.DB
	cfg *config.Config
}

func NewTaskHandler(db *sql.DB, cfg *config.Config) *TaskHandler {
	return &TaskHandler{db: db, cfg: cfg}
}

func (h *TaskHandler) List(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	// Show all tasks for ADMIN+, only assigned tasks for EMPLOYEE
	query := `
		SELECT t.id, t.project_id, p.name as project_name, t.title, t.description, t.status, t.priority,
		       t.assigned_to, u1.username as assigned_to_name, t.due_date,
		       t.created_by, u2.username as created_by_name, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN users u1 ON t.assigned_to = u1.id
		INNER JOIN users u2 ON t.created_by = u2.id
		WHERE t.deleted_at IS NULL
	`

	// Filter for non-admin users
	role := &models.Role{Name: user.RoleName}
	if !role.HasPermission(models.RoleAdmin) {
		query += " AND (t.assigned_to = ? OR p.owner_id = ?)"
	}

	query += " ORDER BY t.created_at DESC"

	var rows *sql.Rows
	var err error

	if role.HasPermission(models.RoleAdmin) {
		rows, err = h.db.Query(query)
	} else {
		rows, err = h.db.Query(query, user.UserID, user.UserID)
	}

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to query tasks")
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.ProjectName, &task.Title, &task.Description,
			&task.Status, &task.Priority, &task.AssignedTo, &task.AssignedToName, &task.DueDate,
			&task.CreatedBy, &task.CreatedByName, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			log.Ctx(r.Context()).Error().Err(err).Msg("failed to scan task")
			continue
		}
		tasks = append(tasks, task)
	}

	utils.RespondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetMyTasks(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	query := `
		SELECT t.id, t.project_id, p.name as project_name, t.title, t.description, t.status, t.priority,
		       t.assigned_to, u1.username as assigned_to_name, t.due_date,
		       t.created_by, u2.username as created_by_name, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN users u1 ON t.assigned_to = u1.id
		INNER JOIN users u2 ON t.created_by = u2.id
		WHERE t.deleted_at IS NULL AND t.assigned_to = ?
		ORDER BY t.due_date ASC, t.priority DESC
	`

	rows, err := h.db.Query(query, user.UserID)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to query my tasks")
		utils.RespondError(w, http.StatusInternalServerError, "failed to fetch tasks")
		return
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var task models.Task
		err := rows.Scan(
			&task.ID, &task.ProjectID, &task.ProjectName, &task.Title, &task.Description,
			&task.Status, &task.Priority, &task.AssignedTo, &task.AssignedToName, &task.DueDate,
			&task.CreatedBy, &task.CreatedByName, &task.CreatedAt, &task.UpdatedAt,
		)
		if err != nil {
			log.Ctx(r.Context()).Error().Err(err).Msg("failed to scan task")
			continue
		}
		tasks = append(tasks, task)
	}

	utils.RespondJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var task models.Task
	query := `
		SELECT t.id, t.project_id, p.name as project_name, t.title, t.description, t.status, t.priority,
		       t.assigned_to, u1.username as assigned_to_name, t.due_date,
		       t.created_by, u2.username as created_by_name, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN users u1 ON t.assigned_to = u1.id
		INNER JOIN users u2 ON t.created_by = u2.id
		WHERE t.id = ? AND t.deleted_at IS NULL
	`

	err := h.db.QueryRow(query, id).Scan(
		&task.ID, &task.ProjectID, &task.ProjectName, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.AssignedTo, &task.AssignedToName, &task.DueDate,
		&task.CreatedBy, &task.CreatedByName, &task.CreatedAt, &task.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		utils.RespondError(w, http.StatusNotFound, "task not found")
		return
	}

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to get task")
		utils.RespondError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	utils.RespondJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	var req models.CreateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Title == "" {
		utils.RespondError(w, http.StatusBadRequest, "task title is required")
		return
	}

	if req.ProjectID == "" {
		utils.RespondError(w, http.StatusBadRequest, "project_id is required")
		return
	}

	// Sanitize
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)

	if req.Status == "" {
		req.Status = models.TaskStatusNotStarted
	}

	if req.Priority == "" {
		req.Priority = models.TaskPriorityMedium
	}

	// Parse due date
	var dueDate sql.NullTime
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			dueDate = sql.NullTime{Time: parsedDate, Valid: true}
		}
	}

	// Handle assigned_to
	var assignedTo sql.NullString
	if req.AssignedTo != "" {
		assignedTo = sql.NullString{String: req.AssignedTo, Valid: true}
	}

	// Insert task
	id := uuid.NewString()
	_, err := h.db.Exec(`
		INSERT INTO tasks (id, project_id, title, description, status, priority, assigned_to, due_date, created_by)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`, id, req.ProjectID, req.Title, req.Description, req.Status, req.Priority, assignedTo, dueDate, user.UserID)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to create task")
		utils.RespondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "task_created", "task", id, r.RemoteAddr)

	// Fetch created task
	var task models.Task
	query := `
		SELECT t.id, t.project_id, p.name as project_name, t.title, t.description, t.status, t.priority,
		       t.assigned_to, u1.username as assigned_to_name, t.due_date,
		       t.created_by, u2.username as created_by_name, t.created_at, t.updated_at
		FROM tasks t
		INNER JOIN projects p ON t.project_id = p.id
		LEFT JOIN users u1 ON t.assigned_to = u1.id
		INNER JOIN users u2 ON t.created_by = u2.id
		WHERE t.id = ?
	`

	err = h.db.QueryRow(query, id).Scan(
		&task.ID, &task.ProjectID, &task.ProjectName, &task.Title, &task.Description,
		&task.Status, &task.Priority, &task.AssignedTo, &task.AssignedToName, &task.DueDate,
		&task.CreatedBy, &task.CreatedByName, &task.CreatedAt, &task.UpdatedAt,
	)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to fetch created task")
		utils.RespondError(w, http.StatusInternalServerError, "task created but failed to fetch")
		return
	}

	utils.RespondJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) Update(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	var req models.UpdateTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate
	if req.Title == "" {
		utils.RespondError(w, http.StatusBadRequest, "task title is required")
		return
	}

	// Sanitize
	req.Title = utils.SanitizeString(req.Title)
	req.Description = utils.SanitizeString(req.Description)

	// Parse due date
	var dueDate sql.NullTime
	if req.DueDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.DueDate)
		if err == nil {
			dueDate = sql.NullTime{Time: parsedDate, Valid: true}
		}
	}

	// Handle assigned_to
	var assignedTo sql.NullString
	if req.AssignedTo != "" {
		assignedTo = sql.NullString{String: req.AssignedTo, Valid: true}
	}

	// Check if task exists
	var exists int
	err := h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE id = ? AND deleted_at IS NULL", id).Scan(&exists)
	if err != nil || exists == 0 {
		utils.RespondError(w, http.StatusNotFound, "task not found")
		return
	}

	// Update task
	_, err = h.db.Exec(`
		UPDATE tasks
		SET title = ?, description = ?, status = ?, priority = ?, assigned_to = ?, due_date = ?, updated_at = NOW()
		WHERE id = ?
	`, req.Title, req.Description, req.Status, req.Priority, assignedTo, dueDate, id)

	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to update task")
		utils.RespondError(w, http.StatusInternalServerError, "failed to update task")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "task_updated", "task", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "task updated successfully")
}

func (h *TaskHandler) Delete(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())
	id := chi.URLParam(r, "id")

	// Soft delete
	result, err := h.db.Exec("UPDATE tasks SET deleted_at = NOW() WHERE id = ? AND deleted_at IS NULL", id)
	if err != nil {
		log.Ctx(r.Context()).Error().Err(err).Msg("failed to delete task")
		utils.RespondError(w, http.StatusInternalServerError, "failed to delete task")
		return
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		utils.RespondError(w, http.StatusNotFound, "task not found")
		return
	}

	// Audit log
	h.auditLog(user.UserID, "task_deleted", "task", id, r.RemoteAddr)

	utils.RespondSuccess(w, nil, "task deleted successfully")
}

func (h *TaskHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	user, _ := middleware.GetUserFromContext(r.Context())

	// Get task statistics
	stats := make(map[string]interface{})

	// Total tasks (all or assigned based on role)
	role := &models.Role{Name: user.RoleName}
	var totalTasks, completedTasks, inProgressTasks, totalProjects int

	if role.HasPermission(models.RoleAdmin) {
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE deleted_at IS NULL").Scan(&totalTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = ? AND deleted_at IS NULL", models.TaskStatusCompleted).Scan(&completedTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = ? AND deleted_at IS NULL", models.TaskStatusInProgress).Scan(&inProgressTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM projects WHERE deleted_at IS NULL").Scan(&totalProjects)
	} else {
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE assigned_to = ? AND deleted_at IS NULL", user.UserID).Scan(&totalTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE assigned_to = ? AND status = ? AND deleted_at IS NULL", user.UserID, models.TaskStatusCompleted).Scan(&completedTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM tasks WHERE assigned_to = ? AND status = ? AND deleted_at IS NULL", user.UserID, models.TaskStatusInProgress).Scan(&inProgressTasks)
		h.db.QueryRow("SELECT COUNT(*) FROM projects WHERE owner_id = ? AND deleted_at IS NULL", user.UserID).Scan(&totalProjects)
	}

	stats["total_tasks"] = totalTasks
	stats["completed_tasks"] = completedTasks
	stats["in_progress_tasks"] = inProgressTasks
	stats["total_projects"] = totalProjects

	utils.RespondJSON(w, http.StatusOK, stats)
}

func (h *TaskHandler) auditLog(userID, action, entityType, entityID, ipAddress string) {
	_, err := h.db.Exec(`
		INSERT INTO audit_logs (user_id, action, entity_type, entity_id, ip_address)
		VALUES (?, ?, ?, ?, ?)
	`, userID, action, entityType, entityID, ipAddress)

	if err != nil {
		log.Error().Err(err).Msg("failed to create audit log")
	}
}
