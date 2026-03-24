package handlers

import (
	"encoding/json"
	"localVercel/db"
	"localVercel/internal/deployer"
	"localVercel/models"
	"localVercel/utils"
	"net/http"
	"strings"
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
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var projects []models.Project
	db.DB.Where("user_id = ?", user.ID).Find(&projects)

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Projects retrieved", projects))
}

// HandleCreateProject создает новый проект
func (h *ProjectHandler) HandleCreateProject(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var req struct {
		Name         string `json:"name"`
		Description  string `json:"description"`
		Framework    string `json:"framework"`
		ProjectType  string `json:"project_type"`
		BuildCommand string `json:"build_command"`
		StartCommand string `json:"start_command"`
		RootDir      string `json:"root_dir"`
		OutputDir    string `json:"output_dir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid request", nil))
		return
	}

	if req.Name == "" {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Project name is required", nil))
		return
	}
	projectType := normalizeProjectType(req.ProjectType)
	if !isAllowedProjectType(projectType) {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "project_type must be backend or telegram", nil))
		return
	}
	if err := deployer.ValidateRootDirInput(req.RootDir); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid root_dir: "+err.Error(), nil))
		return
	}
	if strings.TrimSpace(req.OutputDir) == "" {
		req.OutputDir = "dist"
	}

	project := models.Project{
		UserID:      user.ID,
		Name:        req.Name,
		Description: req.Description,
		Framework:   req.Framework,
		ProjectType: projectType,
		BuildCmd:    req.BuildCommand,
		StartCmd:    req.StartCommand,
		RootDir:     req.RootDir,
		OutputDir:   req.OutputDir,
	}
	if project.ProjectType == "" {
		project.ProjectType = ProjectTypeBackend
	}
	if project.RootDir == "" {
		project.RootDir = "."
	}

	// Important: We need to handle the ID generation if GORM doesn't correctly use the default function
	// or if the struct definition with gorm.Model is problematic.
	// Assuming the `default:gen_random_uuid()` works in Postgres.

	if result := db.DB.Create(&project); result.Error != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to create project: "+result.Error.Error(), nil))
		return
	}

	utils.WriteJSON(w, http.StatusCreated, h.jsonResponse(true, "Project created", project))
}

// HandleGetProject получает проект по ID
func (h *ProjectHandler) HandleGetProject(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("id")
	if projectID == "" {
		projectID = r.PathValue("projectId")
	}

	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}

	var project models.Project
	if err := db.DB.Where("id = ? AND user_id = ?", projectID, user.ID).First(&project).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Project retrieved", project))
}
