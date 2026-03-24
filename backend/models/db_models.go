package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `json:"-"` // Может быть пустым для OAuth пользователей
	Name         string         `json:"name"`
	AvatarURL    string         `json:"avatar_url,omitempty"`
	GitHubID     int            `gorm:"column:git_hub_id;uniqueIndex" json:"github_id"`
	GitHubLogin  string         `gorm:"column:git_hub_login" json:"github_login"`
	GitHubToken  string         `gorm:"column:git_hub_token" json:"-"`
	LastLoginAt  time.Time      `gorm:"column:last_login_at" json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	Projects []Project `json:"projects,omitempty" gorm:"foreignKey:UserID"`
}

type Session struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID    string         `gorm:"not null;index;type:uuid" json:"user_id"` // важно: type:uuid
	Token     string         `gorm:"uniqueIndex;not null" json:"-"`
	ExpiresAt time.Time      `json:"expires_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Project struct {
	gorm.Model
	ID               string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	UserID           string         `gorm:"not null;index" json:"user_id"`
	Name             string         `gorm:"not null" json:"name"`
	Description      string         `json:"description"`
	Repository       string         `json:"repository"`
	Branch           string         `json:"branch" default:"main"`
	Framework        string         `json:"framework"`
	ProjectType      string         `gorm:"default:backend" json:"project_type"` // backend | telegram
	BuildCmd         string         `json:"build_command"`
	StartCmd         string         `json:"start_command"`
	RootDir          string         `gorm:"default:." json:"root_dir"`
	OutputDir        string         `json:"output_dir" default:"dist"`
	RuntimePort      int            `json:"runtime_port"`
	RuntimeState     string         `gorm:"default:configured" json:"runtime_state"` // configured|deploying|running|failed|stopped
	RuntimeHost      string         `json:"runtime_host"`
	RuntimeContainer string         `json:"runtime_container"`
	DBState          string         `gorm:"default:not_configured" json:"db_state"`
	DBContainer      string         `json:"db_container"`
	DBName           string         `json:"db_name"`
	DBUser           string         `json:"db_user"`
	DBPassword       string         `json:"-"`
	DBPort           int            `json:"db_port"`
	EnvVars          []EnvVar       `json:"env_vars,omitempty" gorm:"foreignKey:ProjectID"`
	Deployments      []Deployment   `json:"deployments,omitempty" gorm:"foreignKey:ProjectID"`
	WebhookID        int            `json:"webhook_id,omitempty"`
	WebhookSecret    string         `json:"-"`
	Telegram         TelegramConfig `json:"telegram,omitempty" gorm:"foreignKey:ProjectID"`
}

type EnvVar struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProjectID string `gorm:"not null;index" json:"project_id"`
	Key       string `gorm:"not null" json:"key"`
	Value     string `json:"value"` // В проде нужно шифровать
}

type Deployment struct {
	gorm.Model
	ID           string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProjectID    string    `gorm:"not null;index" json:"project_id"`
	Status       string    `json:"status"` // deploying|building|ready|failed|cancelled
	CommitSHA    string    `json:"commit_sha"`
	CommitMsg    string    `json:"commit_message"`
	Branch       string    `json:"branch"`
	Logs         string    `json:"logs" gorm:"type:text"`
	ArtifactPath string    `json:"artifact_path"`
	BuildTime    int       `json:"build_time"`
	StartedAt    time.Time `json:"started_at"`
	CompletedAt  time.Time `json:"completed_at,omitempty"`
}

type RuntimeInstance struct {
	gorm.Model
	ID            string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProjectID     string    `gorm:"not null;index" json:"project_id"`
	Status        string    `json:"status"` // starting|running|stopped|failed
	PID           int       `json:"pid"`
	Host          string    `json:"host"`
	Command       string    `json:"command"`
	LastError     string    `json:"last_error"`
	LastStartedAt time.Time `json:"last_started_at"`
}

type TelegramConfig struct {
	gorm.Model
	ID         string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ProjectID  string `gorm:"uniqueIndex;not null" json:"project_id"`
	Mode       string `gorm:"default:polling" json:"mode"` // polling|webhook
	BotToken   string `json:"-"`
	WebhookURL string `json:"webhook_url"`
	IsActive   bool   `gorm:"default:false" json:"is_active"`
	LastError  string `json:"last_error"`
}
