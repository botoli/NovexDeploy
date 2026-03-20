package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localVercel/db"
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

func NewGitHubHandler(base *Handler) *GitHubHandler {
	return &GitHubHandler{
		Handler:        base,
		webhookManager: webhook.NewGitHubWebhookManager(),
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
		WebhookURL:   fmt.Sprintf("http://localhost:8888/webhook/github/%s", projectID),
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

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "Repository connected", project))
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

// HandleWebhook точка входа для GitHub webhook (обновленная версия с сохранением в БД)
func (h *GitHubHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	projectID := strings.TrimPrefix(r.URL.Path, "/webhook/github/")
	
	// Получаем проект из БД
	var project models.Project
	if err := db.DB.First(&project, "id = ?", projectID).Error; err != nil {
		http.Error(w, "Project not found", http.StatusNotFound)
		return
	}

	// Читаем тело запроса
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Проверяем подпись
	signature := r.Header.Get("X-Hub-Signature")
	if !h.webhookManager.VerifyGitHubSignature(project.WebhookSecret, payload, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	event := r.Header.Get("X-GitHub-Event")
	
	switch event {
	case "push":
		h.handlePushEvent(project, payload)
	case "ping":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func (h *GitHubHandler) handlePushEvent(project models.Project, payload []byte) {
	var push models.GitHubWebhookPayload
	if err := json.Unmarshal(payload, &push); err != nil {
		log.Printf("Failed to parse push payload: %v", err)
		return
	}

	// Проверяем ветку
	if project.Branch != "" && !strings.HasSuffix(push.Ref, project.Branch) {
		return
	}

	// Создаем запись о деплое в БД
	deployment := models.Deployment{
		ProjectID: project.ID,
		Status:    "pending",
		Branch:    push.Ref,
		StartedAt: time.Now(),
	}
	
	if len(push.Commits) > 0 {
		deployment.CommitSHA = push.Commits[0].ID
		deployment.CommitMsg = push.Commits[0].Message
	}
	
	db.DB.Create(&deployment)

	// Конфигурация билда
	buildConfig := models.BuildConfig{
		ProjectID:     project.ID,
		Branch:        push.Ref,
		CommitSHA:     deployment.CommitSHA,
		CommitMessage: deployment.CommitMsg,
		BuildCommand:  project.BuildCmd,
		OutputDir:     project.OutputDir,
	}

	// Запускаем билд асинхронно
	go func() {
		// Обновляем статус
		db.DB.Model(&deployment).Update("status", "building")

		result, err := h.webhookManager.Builder.BuildProject(buildConfig)
		
		if err != nil {
			db.DB.Model(&deployment).Updates(map[string]interface{}{
				"status":       "failed",
				"logs":         err.Error(),
				"completed_at": time.Now(),
			})
			return
		}

		// Сохраняем результат
		db.DB.Model(&deployment).Updates(map[string]interface{}{
			"status":       result.Status,
			"logs":         result.Logs,
			"preview_url":  result.PreviewURL,
			"build_time":   result.Duration,
			"completed_at": result.CompletedAt,
		})
	}()
}

// Вспомогательные функции
func generateSessionToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// Остальные вспомогательные функции (exchangeCodeForToken, getGitHubUser, getUserRepos) остаются без изменений

func (h *GitHubHandler) exchangeCodeForToken(code string) (string, error) {
	tokenURL := "https://github.com/login/oauth/access_token"
	
	req, err := http.NewRequest("POST", tokenURL, nil)
	if err != nil {
		return "", err
	}

	q := req.URL.Query()
	q.Add("client_id", h.clientID)
	q.Add("client_secret", h.clientSecret)
	q.Add("code", code)
	q.Add("redirect_uri", h.redirectURL)
	req.URL.RawQuery = q.Encode()
	
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResponse struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		Scope       string `json:"scope"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return "", err
	}

	return tokenResponse.AccessToken, nil
}

func (h *GitHubHandler) getGitHubUser(accessToken string) (*models.GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var user models.GitHubUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (h *GitHubHandler) getUserRepos(accessToken string) ([]models.GitHubRepo, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user/repos?sort=updated&per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var repos []models.GitHubRepo
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, err
	}

	return repos, nil
}

func generateSecret() string {
	b := make([]byte, 20)
	for i := range b {
		b[i] = byte('a' + i%26)
	}
	return string(b)
}