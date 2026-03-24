package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"localVercel/db"
	"localVercel/internal/queue"
	rt "localVercel/internal/runtime"
	"localVercel/models"
	"localVercel/utils"
	"net/http"
	"strings"
	"time"
)

type V1Handler struct {
	*Handler
	Queue   queue.Queue
	Runtime *rt.Manager
}

func NewV1Handler(base *Handler, q queue.Queue, runtime *rt.Manager) *V1Handler {
	return &V1Handler{Handler: base, Queue: q, Runtime: runtime}
}

func (h *V1Handler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		db.DB.Where("token = ?", cookie.Value).Delete(&models.Session{})
	}
	http.SetCookie(w, &http.Cookie{Name: "session_token", Value: "", MaxAge: -1, Path: "/"})
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Logged out", nil))
}

func (h *V1Handler) HandlePatchProject(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		BuildCmd    string `json:"build_command"`
		StartCmd    string `json:"start_command"`
		OutputDir   string `json:"output_dir"`
		Branch      string `json:"branch"`
		ProjectType string `json:"project_type"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid body", nil))
		return
	}
	if req.Name != "" {
		project.Name = req.Name
	}
	if req.Description != "" {
		project.Description = req.Description
	}
	if req.BuildCmd != "" {
		project.BuildCmd = req.BuildCmd
	}
	if req.StartCmd != "" {
		project.StartCmd = req.StartCmd
	}
	if req.OutputDir != "" {
		project.OutputDir = req.OutputDir
	}
	if req.Branch != "" {
		project.Branch = req.Branch
	}
	if req.ProjectType != "" {
		project.ProjectType = req.ProjectType
	}
	db.DB.Save(project)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Project updated", project))
}

func (h *V1Handler) HandleDeleteProject(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	db.DB.Delete(project)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Project deleted", nil))
}

func (h *V1Handler) HandleManualDeploy(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}

	deployment := models.Deployment{
		ProjectID: project.ID,
		Status:    "pending",
		Branch:    project.Branch,
		StartedAt: time.Now(),
	}
	if err := db.DB.Create(&deployment).Error; err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to create deployment", nil))
		return
	}

	payload := map[string]interface{}{
		"deployment_id": deployment.ID,
		"project_id":    project.ID,
		"repo_url":      fmt.Sprintf("https://github.com/%s.git", project.Repository),
		"branch":        project.Branch,
		"build_cmd":     project.BuildCmd,
		"output_dir":    project.OutputDir,
	}
	b, _ := json.Marshal(payload)
	job := &queue.Job{ID: fmt.Sprintf("job_%d", time.Now().UnixNano()), Type: "deploy", Payload: b, CreatedAt: time.Now(), Status: "pending"}
	if err := h.Queue.Enqueue(context.Background(), job); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to enqueue deployment", nil))
		return
	}
	utils.WriteJSON(w, http.StatusCreated, h.jsonResponse(true, "Deployment started", deployment))
}

func (h *V1Handler) HandleListDeployments(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var deployments []models.Deployment
	db.DB.Where("project_id = ?", project.ID).Order("created_at desc").Find(&deployments)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Deployments retrieved", deployments))
}

func (h *V1Handler) HandleGetDeployment(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	deploymentID := r.PathValue("deploymentId")
	var deployment models.Deployment
	if err := db.DB.First(&deployment, "id = ?", deploymentID).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Deployment not found", nil))
		return
	}
	if _, ok := h.ownedProject(deployment.ProjectID, user.ID); !ok {
		utils.WriteJSON(w, http.StatusForbidden, h.jsonResponse(false, "Access denied", nil))
		return
	}
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Deployment retrieved", deployment))
}

func (h *V1Handler) HandleCancelDeployment(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	deploymentID := r.PathValue("deploymentId")
	var deployment models.Deployment
	if err := db.DB.First(&deployment, "id = ?", deploymentID).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Deployment not found", nil))
		return
	}
	if _, ok := h.ownedProject(deployment.ProjectID, user.ID); !ok {
		utils.WriteJSON(w, http.StatusForbidden, h.jsonResponse(false, "Access denied", nil))
		return
	}
	deployment.Status = "cancelled"
	deployment.CompletedAt = time.Now()
	db.DB.Save(&deployment)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Deployment cancelled", deployment))
}

func (h *V1Handler) HandleRuntimeStatus(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var runtime models.RuntimeInstance
	db.DB.Where("project_id = ?", project.ID).Order("created_at desc").First(&runtime)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Runtime status", map[string]interface{}{
		"project_id": project.ID,
		"state":      project.RuntimeState,
		"instance":   runtime,
	}))
}

func (h *V1Handler) runtimeAction(w http.ResponseWriter, r *http.Request, action string) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}

	artifactDir := fmt.Sprintf("./workspace/deployments/%s", project.ID)
	var envVars []models.EnvVar
	db.DB.Where("project_id = ?", project.ID).Find(&envVars)
	env := make([]string, 0, len(envVars))
	for _, item := range envVars {
		env = append(env, item.Key+"="+item.Value)
	}

	switch action {
	case "start":
		info, err := h.Runtime.Start(project.ID, artifactDir, project.StartCmd, env)
		if err != nil {
			project.RuntimeState = "failed"
			db.DB.Save(project)
			utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, err.Error(), nil))
			return
		}
		project.RuntimeState = "running"
		db.DB.Save(project)
		inst := models.RuntimeInstance{ProjectID: project.ID, Status: "running", PID: info.PID, Command: info.Command, Host: "local", LastStartedAt: time.Now()}
		db.DB.Create(&inst)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Runtime started", inst))
	case "stop":
		if err := h.Runtime.Stop(project.ID); err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, err.Error(), nil))
			return
		}
		project.RuntimeState = "stopped"
		db.DB.Save(project)
		db.DB.Model(&models.RuntimeInstance{}).Where("project_id = ?", project.ID).Updates(map[string]interface{}{"status": "stopped"})
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Runtime stopped", nil))
	case "restart":
		info, err := h.Runtime.Restart(project.ID, artifactDir, project.StartCmd, env)
		if err != nil {
			utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, err.Error(), nil))
			return
		}
		project.RuntimeState = "running"
		db.DB.Save(project)
		inst := models.RuntimeInstance{ProjectID: project.ID, Status: "running", PID: info.PID, Command: info.Command, Host: "local", LastStartedAt: time.Now()}
		db.DB.Create(&inst)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Runtime restarted", inst))
	}
}

func (h *V1Handler) HandleRuntimeStart(w http.ResponseWriter, r *http.Request) {
	h.runtimeAction(w, r, "start")
}
func (h *V1Handler) HandleRuntimeStop(w http.ResponseWriter, r *http.Request) {
	h.runtimeAction(w, r, "stop")
}
func (h *V1Handler) HandleRuntimeRestart(w http.ResponseWriter, r *http.Request) {
	h.runtimeAction(w, r, "restart")
}

func (h *V1Handler) HandleTelegramConfig(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var req struct {
		Mode       string `json:"mode"`
		BotToken   string `json:"bot_token"`
		WebhookURL string `json:"webhook_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid body", nil))
		return
	}
	cfg := models.TelegramConfig{ProjectID: project.ID, Mode: req.Mode, BotToken: req.BotToken, WebhookURL: req.WebhookURL, IsActive: true}
	db.DB.Where("project_id = ?", project.ID).Delete(&models.TelegramConfig{})
	db.DB.Create(&cfg)
	project.ProjectType = "telegram_bot"
	db.DB.Save(project)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Telegram config saved", map[string]interface{}{
		"project_id": project.ID,
		"mode":       cfg.Mode,
		"is_active":  cfg.IsActive,
	}))
}

func (h *V1Handler) HandleTelegramStatus(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var cfg models.TelegramConfig
	if err := db.DB.Where("project_id = ?", project.ID).First(&cfg).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Telegram config not found", nil))
		return
	}
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Telegram status", map[string]interface{}{
		"mode":       cfg.Mode,
		"is_active":  cfg.IsActive,
		"webhook":    cfg.WebhookURL,
		"last_error": cfg.LastError,
	}))
}

func (h *V1Handler) HandleTelegramWebhookSync(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var cfg models.TelegramConfig
	if err := db.DB.Where("project_id = ?", project.ID).First(&cfg).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Telegram config not found", nil))
		return
	}
	cfg.IsActive = true
	cfg.LastError = ""
	db.DB.Save(&cfg)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Telegram webhook synced", map[string]interface{}{"project_id": project.ID}))
}

func (h *V1Handler) HandleListEnv(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var vars []models.EnvVar
	db.DB.Where("project_id = ?", project.ID).Find(&vars)
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Env vars", vars))
}

func (h *V1Handler) HandleUpsertEnv(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Key == "" {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid body", nil))
		return
	}
	var variable models.EnvVar
	if err := db.DB.Where("project_id = ? AND key = ?", project.ID, req.Key).First(&variable).Error; err != nil {
		variable = models.EnvVar{ProjectID: project.ID, Key: req.Key, Value: req.Value}
		db.DB.Create(&variable)
	} else {
		variable.Value = req.Value
		db.DB.Save(&variable)
	}
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Env var saved", variable))
}

func (h *V1Handler) HandleDeleteEnv(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	project, ok := h.ownedProject(r.PathValue("projectId"), user.ID)
	if !ok {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}
	key := strings.TrimSpace(r.PathValue("key"))
	db.DB.Where("project_id = ? AND key = ?", project.ID, key).Delete(&models.EnvVar{})
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Env var deleted", nil))
}

func (h *V1Handler) HandleDeploymentLogs(w http.ResponseWriter, r *http.Request) {
	user, err := currentUserFromSession(r)
	if err != nil {
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	deploymentID := r.PathValue("deploymentId")
	var deployment models.Deployment
	if err := db.DB.First(&deployment, "id = ?", deploymentID).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Deployment not found", nil))
		return
	}
	if _, ok := h.ownedProject(deployment.ProjectID, user.ID); !ok {
		utils.WriteJSON(w, http.StatusForbidden, h.jsonResponse(false, "Access denied", nil))
		return
	}
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Logs retrieved", map[string]interface{}{
		"deployment_id": deployment.ID,
		"logs":          deployment.Logs,
	}))
}

func (h *V1Handler) ownedProject(projectID, userID string) (*models.Project, bool) {
	var project models.Project
	if err := db.DB.Where("id = ? AND user_id = ?", projectID, userID).First(&project).Error; err != nil {
		return nil, false
	}
	return &project, true
}
