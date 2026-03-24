package handlers

import (
	"encoding/json"
	"localVercel/db"
	"localVercel/models"
	"localVercel/utils"
	"net/http"
	"time"
)

type ProjectHandler struct {
	*Handler
}

func NewProjectHandler(base *Handler) *ProjectHandler {
	return &ProjectHandler{
		Handler: base,
	}
}

// HandleListProjects возвращает список проектов пользователя
func (h *ProjectHandler) HandleListProjects(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var session models.Session
	if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
		return
	}

	var projects []models.Project
	db.DB.Where("user_id = ?", session.UserID).Find(&projects)

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Projects retrieved", projects))
}

// HandleCreateProject создает новый проект
func (h *ProjectHandler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	// Получаем пользователя из сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var session models.Session
	if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
		return
	}

	var req struct {
		Name          string `json:"name"`
		Description   string `json:"description"`
		Framework     string `json:"framework"`
		BuildCommand  string `json:"build_command"`
		OutputDir     string `json:"output_dir"`
		// Repository info can be passed here to preload, but connecting happens via separate endpoint usually
		// but we can support setting it here if needed. For now let's keep it minimal.
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid request", nil))
		return
	}

	if req.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Project name is required", nil))
		return
	}

	project := models.Project{
		UserID:      session.UserID,
		Name:        req.Name,
		Description: req.Description,
		Framework:   req.Framework,
		BuildCmd:    req.BuildCommand,
		OutputDir:   req.OutputDir,
	}

	// Important: We need to handle the ID generation if GORM doesn't correctly use the default function 
	// or if the struct definition with gorm.Model is problematic. 
	// Assuming the `default:gen_random_uuid()` works in Postgres.
	
	if result := db.DB.Create(&project); result.Error != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to create project: " + result.Error.Error(), nil))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, h.jsonResponse(true, "Project created", project))
}

// HandleGetProject получает проект по ID
func (h *ProjectHandler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")

	cookie, err := r.Cookie("session_token")
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var session models.Session
	if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
		return
	}

	var project models.Project
	if err := db.DB.Where("id = ? AND user_id = ?", projectID, session.UserID).First(&project).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Project retrieved", project))
}
