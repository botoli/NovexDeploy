package models

import "time"

type GitHubAuthRequest struct {
	Code string `json:"code" binding:"required"`
}

type GitHubRepo struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	FullName    string    `json:"full_name"`
	Description string    `json:"description"`
	Private     bool      `json:"private"`
	HTMLURL     string    `json:"html_url"`
	CloneURL    string    `json:"clone_url"`
	SSHURL      string    `json:"ssh_url"`
	DefaultBranch string  `json:"default_branch"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	PushedAt    time.Time `json:"pushed_at"`
	Size        int       `json:"size"`
	Language    string    `json:"language"`
}

type GitHubWebhookPayload struct {
	Ref        string `json:"ref"`
	Before     string `json:"before"`
	After      string `json:"after"`
	Repository struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		FullName string `json:"full_name"`
		Private  bool   `json:"private"`
		HTMLURL  string `json:"html_url"`
		CloneURL string `json:"clone_url"`
		Owner    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
			Login string `json:"login"`
			ID    int    `json:"id"`
		} `json:"owner"`
	} `json:"repository"`
	Pusher struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
	} `json:"sender"`
	Commits []struct {
		ID      string `json:"id"`
		Message string `json:"message"`
		URL     string `json:"url"`
		Author  struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
	} `json:"commits"`
}

type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	Bio       string `json:"bio"`
}

type ConnectedGitHubAccount struct {
	UserID        string    `json:"user_id"`
	GitHubID      int       `json:"github_id"`
	GitHubLogin   string    `json:"github_login"`
	AccessToken   string    `json:"-"` // Не возвращаем в JSON
	ConnectedAt   time.Time `json:"connected_at"`
	Repositories  []GitHubRepo `json:"repositories,omitempty"`
}

type BuildConfig struct {
	ProjectID     string            `json:"project_id"`
	Branch        string            `json:"branch"`
	CommitSHA     string            `json:"commit_sha"`
	CommitMessage string            `json:"commit_message"`
	BuildCommand  string            `json:"build_command"`
	OutputDir     string            `json:"output_dir"`
	Environment   map[string]string `json:"environment"`
	Framework     string            `json:"framework"` // react, vue, angular, static, go, python и т.д.
}

type BuildResult struct {
	ID            string    `json:"id"`
	ProjectID     string    `json:"project_id"`
	Status        string    `json:"status"` // pending, building, success, failed
	CommitSHA     string    `json:"commit_sha"`
	CommitMessage string    `json:"commit_message"`
	Branch        string    `json:"branch"`
	Logs          string    `json:"logs"`
	OutputPath    string    `json:"output_path"`
	PreviewURL    string    `json:"preview_url"`
	StartedAt     time.Time `json:"started_at"`
	CompletedAt   time.Time `json:"completed_at,omitempty"`
	Duration      int       `json:"duration"` // in seconds
}

type WebhookConfig struct {
	ProjectID   string `json:"project_id"`
	GitHubRepo  string `json:"github_repo"`
	WebhookURL  string `json:"webhook_url"`
	WebhookID   int    `json:"webhook_id,omitempty"`
	Active      bool   `json:"active"`
	Secret      string `json:"secret"`
	Events      []string `json:"events"` // push, pull_request и т.д.
	Branch      string `json:"branch"`    // ветка для авто-деплоя
	BuildCommand string `json:"build_command"`
	OutputDir    string `json:"output_dir"`
}