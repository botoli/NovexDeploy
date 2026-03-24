package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"localVercel/db"
	"localVercel/internal/queue"
	"localVercel/models"
	"localVercel/utils"
	"localVercel/webhook"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type GitHubHandler struct {
	*Handler
	webhookManager *webhook.GitHubWebhookManager
	clientID       string
	clientSecret   string
	redirectURL    string
}

func NewGitHubHandler(base *Handler, q queue.Queue) *GitHubHandler {
	return &GitHubHandler{
		Handler:        base,
		webhookManager: webhook.NewGitHubWebhookManager(q),
		clientID:       os.Getenv("GITHUB_CLIENT_ID"),
		clientSecret:   os.Getenv("GITHUB_CLIENT_SECRET"),
		redirectURL:    "http://localhost:8888/auth/github/callback",
	}
}

// HandleGitHubLogin перенаправляет на GitHub OAuth
func (h *GitHubHandler) HandleGitHubLogin(w http.ResponseWriter, r *http.Request) {
	authURL := fmt.Sprintf(
		"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s",
		h.clientID, h.redirectURL, "repo admin:repo_hook user:email",
	)
	http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)
}

// HandleGitHubCallback обрабатывает callback от GitHub
func (h *GitHubHandler) HandleGitHubCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Missing code", nil))
		return
	}

	log.Println("=== GitHub Callback Started ===")
	log.Println("Code received:", code)

	// Обмениваем код на токен
	token, err := h.exchangeCodeForToken(code)
	if err != nil {
		log.Println("Token exchange failed:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to exchange code: "+err.Error(), nil))
		return
	}
	log.Println("Token obtained successfully")

	// Получаем информацию о пользователе
	githubUser, err := h.getGitHubUser(token)
	if err != nil {
		log.Println("Failed to get user info:", err)
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to get user info: "+err.Error(), nil))
		return
	}
	log.Println("GitHub user:", githubUser.Login, githubUser.ID)

	// Получаем репозитории пользователя
	_, err = h.getUserRepos(token)
	if err != nil {
		log.Println("Failed to get repos:", err)
		// Не фатально, продолжаем
	}

	// Ищем или создаем пользователя в БД
	var user models.User
	result := db.DB.Where("git_hub_id = ?", githubUser.ID).First(&user)

	if result.Error != nil {
		log.Println("Creating new user")
		// Создаем нового пользователя
		user = models.User{
			Email:       githubUser.Email,
			Name:        githubUser.Name,
			AvatarURL:   githubUser.AvatarURL,
			GitHubID:    githubUser.ID,           // это сохранится в git_hub_id
			GitHubLogin: githubUser.Login,        // это сохранится в git_hub_login
			GitHubToken: token,                    // это сохранится в git_hub_token
			LastLoginAt: time.Now(),
		}
		db.DB.Create(&user)
	} else {
		log.Println("Updating existing user")
		// Обновляем существующего
		user.GitHubToken = token
		user.LastLoginAt = time.Now()
		db.DB.Save(&user)
	}
	log.Println("User saved with ID:", user.ID)

	// Создаем сессию
	sessionToken := generateSessionToken()
	session := models.Session{
		UserID:    user.ID,
		Token:     sessionToken,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	db.DB.Create(&session)
	log.Println("Session created with token:", sessionToken)

	// Устанавливаем cookie с токеном
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  session.ExpiresAt,
		HttpOnly: true,
		Path:     "/",
		Domain:   "localhost",
		SameSite: http.SameSiteLaxMode,
	})
	log.Println("Cookie set")

	// Вместо JSON ответа, делаем редирект на фронтенд
	log.Println("Redirecting to frontend")
	http.Redirect(w, r, "http://localhost:5432/projects", http.StatusTemporaryRedirect)
}

// HandleConnectRepo подключает репозиторий к проекту
func (h *GitHubHandler) HandleConnectRepo(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	
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

	var user models.User
	db.DB.First(&user, "id = ?", session.UserID)

	var req struct {
		RepoFullName string `json:"repo_full_name"`
		Branch       string `json:"branch"`
		BuildCommand string `json:"build_command"`
		OutputDir    string `json:"output_dir"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, h.jsonResponse(false, "Invalid request", nil))
		return
	}

	// Получаем проект из БД
	var project models.Project
	if err := db.DB.First(&project, "id = ? AND user_id = ?", projectID, user.ID).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Project not found", nil))
		return
	}

	// Создаем конфигурацию webhook
	webhookConfig := &models.WebhookConfig{
		ProjectID:    projectID,
		GitHubRepo:   req.RepoFullName,
		WebhookURL:   fmt.Sprintf("http://localhost:8888/webhook/github/%s", projectID), // В продакшене нужен публичный URL
		Active:       true,
		Events:       []string{"push"},
		Branch:       req.Branch,
		BuildCommand: req.BuildCommand,
		OutputDir:    req.OutputDir,
		Secret:       generateSecret(),
	}

	// Создаем webhook в GitHub
	if err := h.webhookManager.SetupWebhook(user.GitHubToken, req.RepoFullName, webhookConfig); err != nil {
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "Failed to setup webhook: "+err.Error(), nil))
		return
	}

	// Обновляем проект в БД
	project.Repository = req.RepoFullName
	project.Branch = req.Branch
	project.BuildCmd = req.BuildCommand
	project.OutputDir = req.OutputDir
	project.WebhookID = webhookConfig.WebhookID
	project.WebhookSecret = webhookConfig.Secret
	
	db.DB.Save(&project)

	// --- FEATURE: Trigger Immediate Build on Connect ---
	// Create Deployment record
	deployment := models.Deployment{
		ProjectID: project.ID,
		Status:    "pending",
		Branch:    project.Branch,
		StartedAt: time.Now(),
		// Commit info might be missing here until we fetch it, or wait for clone
	}
	db.DB.Create(&deployment)

	// Create Job
	jobPayload := map[string]interface{}{
		"deployment_id": deployment.ID,
		"project_id":    project.ID,
		"repo_url":      fmt.Sprintf("https://github.com/%s.git", req.RepoFullName), // Construct generic HTTPS clone URL
		"branch":        project.Branch,
		"build_cmd":     project.BuildCmd,
		"output_dir":    project.OutputDir,
	}
	payloadBytes, _ := json.Marshal(jobPayload)
	
	job := &queue.Job{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Type:      "deploy",
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	// Enqueue
	if err := h.webhookManager.Queue.Enqueue(r.Context(), job); err != nil {
		log.Printf("Failed to enqueue initial build: %v", err)
		// Don't fail the request, just log
	}
	// ---------------------------------------------------

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Repository connected and build started", project))
}

// HandleListRepos возвращает список репозиториев пользователя
func (h *GitHubHandler) HandleListRepos(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling /git/repos list request")
	
	// Получаем пользователя из сессии
	cookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("Missing session token cookie")
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Unauthorized", nil))
		return
	}
	log.Println("Session cookie found")

	var session models.Session
	if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
		log.Printf("Session lookup failed: %v", err)
		utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
		return
	}
	log.Printf("Session OK for user %s", session.UserID)

	var user models.User
	if err := db.DB.First(&user, "id = ?", session.UserID).Error; err != nil {
		log.Printf("User lookup failed: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "User not found", nil))
		return
	}
	log.Printf("User loaded: %s, token length: %d", user.GitHubLogin, len(user.GitHubToken))

	repos, err := h.getUserRepos(user.GitHubToken)
	if err != nil {
		log.Printf("Failed to fetch repos: %v", err)
		utils.WriteJSON(w, http.StatusBadGateway, h.jsonResponse(false, "Failed to fetch repos from GitHub", nil))
		return
	}
	
	log.Printf("Returning %d repos to client", len(repos))
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Repositories retrieved", repos))
}

// HandleListBuilds список билдов проекта
func (h *GitHubHandler) HandleListBuilds(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")
	
	var deployments []models.Deployment
	db.DB.Where("project_id = ?", projectID).Order("created_at desc").Find(&deployments)
	
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Builds retrieved", map[string]interface{}{
		"project_id": projectID,
		"builds":     deployments,
	}))
}

// HandleGetBuild статус билда
func (h *GitHubHandler) HandleGetBuild(w http.ResponseWriter, r *http.Request) {
	buildID := r.PathValue("buildId")
	
	var deployment models.Deployment
	if err := db.DB.First(&deployment, "id = ?", buildID).Error; err != nil {
		utils.WriteJSON(w, http.StatusNotFound, h.jsonResponse(false, "Build not found", nil))
		return
	}
	
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Build status", deployment))
}

// HandleWebhook точка входа для GitHub webhook
func (h *GitHubHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	h.webhookManager.HandleWebhook(w, r)
}

// Helper methods

func generateSessionToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return hex.EncodeToString(b)
}

func generateSecret() string {
	b := make([]byte, 20)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (h *GitHubHandler) exchangeCodeForToken(code string) (string, error) {
	tokenURL := "https://github.com/login/oauth/access_token"
	reqBody, _ := json.Marshal(map[string]string{
		"client_id":     h.clientID,
		"client_secret": h.clientSecret,
		"code":          code,
		"redirect_uri":  h.redirectURL,
	})

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(string(reqBody)))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		AccessToken string `json:"access_token"`
		Error       string `json:"error"`
		ErrorDesc   string `json:"error_description"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Error != "" {
		return "", fmt.Errorf("%s: %s", result.Error, result.ErrorDesc)
	}

	return result.AccessToken, nil
}

func (h *GitHubHandler) getGitHubUser(token string) (*models.GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github api error: %s", resp.Status)
	}

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (h *GitHubHandler) getUserRepos(token string) ([]models.GitHubRepo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&per_page=100", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "token "+token)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("GitHub API Error: %s, Body: %s", resp.Status, string(bodyBytes))
		return nil, fmt.Errorf("github api error: %s", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("GitHub Response (truncated): %.200s", string(bodyBytes)) // Log first 200 chars

	var repos []models.GitHubRepo
	if err := json.Unmarshal(bodyBytes, &repos); err != nil {
		log.Printf("Failed to unmarshal repos: %v", err)
		return nil, err
	}
	// Explicitly handle empty array vs nil to avoid "null" in JSON response
	if repos == nil {
		repos = []models.GitHubRepo{}
	}
	log.Printf("Parsed %d repos from GitHub", len(repos))
	
	// Если репозиториев 0, это странно для активного пользователя - логируем тело
	if len(repos) == 0 {
		log.Println("WARNING: 0 requested repos found. This might be due to auth scopes or empty account.")
	}

	return repos, nil
}
