package webhook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localVercel/builder"
	"localVercel/models"
	"log"
	"net/http"
	"strings"
	"sync"
)

type GitHubWebhookManager struct {
	Builder      *builder.Builder // Сделаем поле публичным
	webhooks     map[string]*models.WebhookConfig
	builds       map[string]*models.BuildResult
	projects     map[string]string
	mu           sync.RWMutex
	githubClient *http.Client
}

func NewGitHubWebhookManager() *GitHubWebhookManager {
	return &GitHubWebhookManager{
		Builder:      builder.NewBuilder(),
		webhooks:     make(map[string]*models.WebhookConfig),
		builds:       make(map[string]*models.BuildResult),
		projects:     make(map[string]string),
		githubClient: &http.Client{},
	}
}

// VerifyGitHubSignature проверяет подпись от GitHub
func (m *GitHubWebhookManager) VerifyGitHubSignature(secret string, payload []byte, signatureHeader string) bool {
	if secret == "" {
		return true // Если секрет не установлен, пропускаем проверку
	}

	// GitHub отправляет подпись в формате sha1=...
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 || parts[0] != "sha1" {
		return false
	}

	signature, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}

	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write(payload)
	expected := mac.Sum(nil)

	return hmac.Equal(signature, expected)
}

// HandleWebhook обрабатывает входящий webhook от GitHub
func (m *GitHubWebhookManager) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Читаем тело запроса
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	// Получаем информацию о проекте из URL
	projectID := strings.TrimPrefix(r.URL.Path, "/webhook/github/")
	
	m.mu.RLock()
	webhookConfig, exists := m.webhooks[projectID]
	m.mu.RUnlock()

	if !exists {
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	// Проверяем подпись
	signature := r.Header.Get("X-Hub-Signature")
	if !m.VerifyGitHubSignature(webhookConfig.Secret, payload, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Парсим событие
	event := r.Header.Get("X-GitHub-Event")
	
	switch event {
	case "push":
		m.handlePushEvent(projectID, webhookConfig, payload)
	case "ping":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
		return
	default:
		log.Printf("Unsupported event: %s", event)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Webhook received"))
}

func (m *GitHubWebhookManager) handlePushEvent(projectID string, config *models.WebhookConfig, payload []byte) {
	var push models.GitHubWebhookPayload
	if err := json.Unmarshal(payload, &push); err != nil {
		log.Printf("Failed to parse push payload: %v", err)
		return
	}

	// Проверяем, что пуш в нужную ветку
	if config.Branch != "" && !strings.HasSuffix(push.Ref, config.Branch) {
		return
	}

	// Получаем последний коммит
	var commitMessage string
	var commitSHA string
	if len(push.Commits) > 0 {
		commitMessage = push.Commits[0].Message
		commitSHA = push.Commits[0].ID
	}

	// Создаем конфигурацию билда
	buildConfig := models.BuildConfig{
		ProjectID:     projectID,
		Branch:        push.Ref,
		CommitSHA:     commitSHA,
		CommitMessage: commitMessage,
		BuildCommand:  config.BuildCommand,
		OutputDir:     config.OutputDir,
		Environment:   make(map[string]string),
	}

	// Запускаем билд асинхронно
	go func() {
		result, err := m.Builder.BuildProject(buildConfig)
		if err != nil {
			log.Printf("Build failed for project %s: %v", projectID, err)
			return
		}
		if result == nil {
			log.Printf("Build failed for project %s: empty build result", projectID)
			return
		}
		
		m.mu.Lock()
		m.builds[result.ID] = result
		m.mu.Unlock()
		
		log.Printf("Build completed for project %s: %s", projectID, result.Status)
	}()
}

// SetupWebhook создает webhook в GitHub
func (m *GitHubWebhookManager) SetupWebhook(accessToken string, repoFullName string, config *models.WebhookConfig) error {
	webhookURL := fmt.Sprintf("https://api.github.com/repos/%s/hooks", repoFullName)
	
	webhookData := map[string]interface{}{
		"name":   "web",
		"active": true,
		"events": config.Events,
		"config": map[string]string{
			"url":          config.WebhookURL,
			"content_type": "json",
			"secret":       config.Secret,
		},
	}

	data, _ := json.Marshal(webhookData)
	
	req, err := http.NewRequest("POST", webhookURL, bytes.NewReader(data))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.githubClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create webhook: %s", body)
	}

	// Парсим ответ чтобы получить ID webhook
	var response struct {
		ID int `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		config.WebhookID = response.ID
	}

	return nil
}

// DeleteWebhook удаляет webhook из GitHub
func (m *GitHubWebhookManager) DeleteWebhook(accessToken string, repoFullName string, webhookID int) error {
	webhookURL := fmt.Sprintf("https://api.github.com/repos/%s/hooks/%d", repoFullName, webhookID)
	
	req, err := http.NewRequest("DELETE", webhookURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+accessToken)

	resp, err := m.githubClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete webhook: %s", body)
	}

	return nil
}

// GetBuildStatus возвращает статус билда
func (m *GitHubWebhookManager) GetBuildStatus(buildID string) (*models.BuildResult, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	build, exists := m.builds[buildID]
	if !exists {
		return nil, fmt.Errorf("build not found")
	}
	return build, nil
}

// ListBuilds возвращает все билды проекта
func (m *GitHubWebhookManager) ListBuilds(projectID string) []*models.BuildResult {
	m.mu.RLock()
	defer m.mu.RUnlock()
	
	var builds []*models.BuildResult
	for _, build := range m.builds {
		if build.ProjectID == projectID {
			builds = append(builds, build)
		}
	}
	return builds
}