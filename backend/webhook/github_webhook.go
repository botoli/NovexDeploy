package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"localVercel/db"
	"localVercel/internal/queue"
	"localVercel/models"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type GitHubWebhookManager struct {
	Queue        queue.Queue
	webhooks     map[string]*models.WebhookConfig
	projects     map[string]string
	deliveries   map[string]time.Time
	mu           sync.RWMutex
	githubClient *http.Client
}

func NewGitHubWebhookManager(q queue.Queue) *GitHubWebhookManager {
	return &GitHubWebhookManager{
		Queue:        q,
		webhooks:     make(map[string]*models.WebhookConfig),
		projects:     make(map[string]string),
		deliveries:   make(map[string]time.Time),
		githubClient: &http.Client{},
	}
}

// VerifyGitHubSignature verifies the signature from GitHub
func (m *GitHubWebhookManager) VerifyGitHubSignature(secret string, payload []byte, signatureHeader string) bool {
	if secret == "" {
		return false
	}
	parts := strings.SplitN(signatureHeader, "=", 2)
	if len(parts) != 2 {
		return false
	}
	signature, err := hex.DecodeString(parts[1])
	if err != nil {
		return false
	}
	var mac hashWithSecret
	switch parts[0] {
	case "sha1":
		mac = hmac.New(sha1.New, []byte(secret))
	case "sha256":
		mac = hmac.New(sha256.New, []byte(secret))
	default:
		return false
	}
	mac.Write(payload)
	expected := mac.Sum(nil)
	return hmac.Equal(signature, expected)
}

type hashWithSecret interface {
	Write([]byte) (int, error)
	Sum([]byte) []byte
}

// HandleWebhook processes incoming webhook
func (m *GitHubWebhookManager) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusBadRequest)
		return
	}

	projectID := r.PathValue("projectId")

	m.mu.RLock()
	webhookConfig, exists := m.webhooks[projectID]
	m.mu.RUnlock()

	if !exists {
		// Try to load from DB
		var project models.Project
		if err := db.DB.First(&project, "id = ?", projectID).Error; err != nil {
			log.Printf("Webhook not found for project %s", projectID)
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		// Reconstruct WebhookConfig
		webhookConfig = &models.WebhookConfig{
			ProjectID:    project.ID,
			GitHubRepo:   project.Repository,
			WebhookID:    project.WebhookID,
			WebhookURL:   fmt.Sprintf("/webhook/github/%s", project.ID), // Reconstruction
			Active:       true,
			Events:       []string{"push"},
			Branch:       project.Branch,
			BuildCommand: project.BuildCmd,
			OutputDir:    project.OutputDir,
			Secret:       project.WebhookSecret,
		}

		m.mu.Lock()
		m.webhooks[projectID] = webhookConfig
		m.mu.Unlock()
	}

	signature := r.Header.Get("X-Hub-Signature-256")
	if signature == "" {
		signature = r.Header.Get("X-Hub-Signature")
	}
	if !m.VerifyGitHubSignature(webhookConfig.Secret, payload, signature) {
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}
	deliveryID := strings.TrimSpace(r.Header.Get("X-GitHub-Delivery"))
	if deliveryID != "" && m.isDuplicateDelivery(deliveryID) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Duplicate ignored"))
		return
	}

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

func (m *GitHubWebhookManager) isDuplicateDelivery(deliveryID string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for key, ts := range m.deliveries {
		if now.Sub(ts) > 30*time.Minute {
			delete(m.deliveries, key)
		}
	}
	if _, exists := m.deliveries[deliveryID]; exists {
		return true
	}
	m.deliveries[deliveryID] = now
	return false
}

func (m *GitHubWebhookManager) handlePushEvent(projectID string, config *models.WebhookConfig, payload []byte) {
	var push models.GitHubWebhookPayload
	if err := json.Unmarshal(payload, &push); err != nil {
		log.Printf("Failed to parse push payload: %v", err)
		return
	}

	if config.Branch != "" && !strings.HasSuffix(push.Ref, config.Branch) {
		return
	}

	// In a real implementation: Create Deployment record in DB first to get ID.
	// For now, generate ID.
	// Create Deployment in DB
	deployment := models.Deployment{
		ID:        fmt.Sprintf("deploy_%d", time.Now().Unix()),
		ProjectID: projectID,
		Status:    "pending",
		Branch:    strings.TrimPrefix(push.Ref, "refs/heads/"),
		StartedAt: time.Now(),
	}

	if len(push.Commits) > 0 {
		deployment.CommitSHA = push.Commits[0].ID
		deployment.CommitMsg = push.Commits[0].Message
	}

	if err := db.DB.Create(&deployment).Error; err != nil {
		log.Printf("Failed to create deployment record: %v", err)
		return
	}

	type DeployPayload struct {
		DeploymentID string `json:"deployment_id"`
		RepoURL      string `json:"repo_url"`
		Branch       string `json:"branch"`
		ProjectID    string `json:"project_id"`
		BuildCmd     string `json:"build_cmd"`
		OutputDir    string `json:"output_dir"`
	}

	jobPayload := DeployPayload{
		DeploymentID: deployment.ID,
		ProjectID:    projectID,
		RepoURL:      push.Repository.CloneURL,
		Branch:       deployment.Branch,
		BuildCmd:     config.BuildCommand,
		OutputDir:    config.OutputDir,
	}

	payloadBytes, _ := json.Marshal(jobPayload)

	job := &queue.Job{
		ID:        fmt.Sprintf("job_%d", time.Now().UnixNano()),
		Type:      "deploy",
		Payload:   payloadBytes,
		CreatedAt: time.Now(),
		Status:    "pending",
	}

	if err := m.Queue.Enqueue(context.Background(), job); err != nil {
		log.Printf("Failed to enqueue job: %v", err)
	} else {
		log.Printf("Enqueued deployment job %s for project %s", job.ID, projectID)
	}
}

// SetupWebhook creates a webhook in GitHub
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

	// Parse response
	var response struct {
		ID int `json:"id"`
	}
	// Warning: ignoring explicit error check for brevity in snippet
	json.NewDecoder(resp.Body).Decode(&response)
	config.WebhookID = response.ID

	m.mu.Lock()
	m.webhooks[config.ProjectID] = config
	m.mu.Unlock()

	return nil
}
