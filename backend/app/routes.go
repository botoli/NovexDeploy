package app

import (
	"localVercel/handlers"
	"net/http"
)

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	h := handlers.New(a.startedAt, &a.upgrader)
	githubHandler := handlers.NewGitHubHandler(h, a.Queue)
	projectHandler := handlers.NewProjectHandler(h)
	v1 := handlers.NewV1Handler(h, a.Queue, a.Runtime)

	// Core auth + GitHub integration
	mux.HandleFunc("GET /auth/github/login", githubHandler.HandleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", githubHandler.HandleGitHubCallback)
	mux.HandleFunc("GET /auth/me", h.HandleGet("auth.me"))
	mux.HandleFunc("GET /git/repos", githubHandler.HandleListRepos)

	// Project + build management
	mux.HandleFunc("GET /projects", projectHandler.HandleListProjects)
	mux.HandleFunc("POST /projects", projectHandler.HandleCreateProject)
	mux.HandleFunc("GET /projects/{id}", projectHandler.HandleGetProject)
	mux.HandleFunc("POST /projects/{projectId}/github/repo", githubHandler.HandleConnectRepo)
	mux.HandleFunc("GET /projects/{projectId}/builds", githubHandler.HandleListBuilds)
	mux.HandleFunc("GET /projects/{projectId}/builds/{buildId}", githubHandler.HandleGetBuild)

	// Public GitHub webhook endpoint
	mux.HandleFunc("POST /webhook/github/{projectId}", githubHandler.HandleWebhook)

	// Realtime logs by deployment
	mux.HandleFunc("GET /ws/logs/{deploymentId}", h.HandleWSLogsByDeployment)

	// System
	mux.HandleFunc("GET /health", h.HandleHealth)

	// v1 API
	mux.HandleFunc("GET /v1/auth/github/login", githubHandler.HandleGitHubLogin)
	mux.HandleFunc("GET /v1/auth/github/callback", githubHandler.HandleGitHubCallback)
	mux.HandleFunc("GET /v1/auth/me", h.HandleGet("auth.me"))
	mux.HandleFunc("POST /v1/auth/logout", v1.HandleLogout)
	mux.HandleFunc("GET /v1/git/repos", githubHandler.HandleListRepos)

	mux.HandleFunc("GET /v1/projects", projectHandler.HandleListProjects)
	mux.HandleFunc("POST /v1/projects", projectHandler.HandleCreateProject)
	mux.HandleFunc("GET /v1/projects/{projectId}", projectHandler.HandleGetProject)
	mux.HandleFunc("PATCH /v1/projects/{projectId}", v1.HandlePatchProject)
	mux.HandleFunc("DELETE /v1/projects/{projectId}", v1.HandleDeleteProject)

	mux.HandleFunc("POST /v1/projects/{projectId}/repo/connect", githubHandler.HandleConnectRepo)
	mux.HandleFunc("POST /v1/projects/{projectId}/deployments", v1.HandleManualDeploy)
	mux.HandleFunc("GET /v1/projects/{projectId}/deployments", v1.HandleListDeployments)
	mux.HandleFunc("GET /v1/deployments/{deploymentId}", v1.HandleGetDeployment)
	mux.HandleFunc("POST /v1/deployments/{deploymentId}/cancel", v1.HandleCancelDeployment)
	mux.HandleFunc("GET /v1/deployments/{deploymentId}/logs", v1.HandleDeploymentLogs)

	mux.HandleFunc("GET /v1/projects/{projectId}/runtime", v1.HandleRuntimeStatus)
	mux.HandleFunc("POST /v1/projects/{projectId}/runtime/start", v1.HandleRuntimeStart)
	mux.HandleFunc("POST /v1/projects/{projectId}/runtime/stop", v1.HandleRuntimeStop)
	mux.HandleFunc("POST /v1/projects/{projectId}/runtime/restart", v1.HandleRuntimeRestart)

	mux.HandleFunc("POST /v1/projects/{projectId}/telegram/config", v1.HandleTelegramConfig)
	mux.HandleFunc("GET /v1/projects/{projectId}/telegram/status", v1.HandleTelegramStatus)
	mux.HandleFunc("POST /v1/projects/{projectId}/telegram/webhook/sync", v1.HandleTelegramWebhookSync)

	mux.HandleFunc("GET /v1/projects/{projectId}/env", v1.HandleListEnv)
	mux.HandleFunc("POST /v1/projects/{projectId}/env", v1.HandleUpsertEnv)
	mux.HandleFunc("DELETE /v1/projects/{projectId}/env/{key}", v1.HandleDeleteEnv)

	mux.HandleFunc("POST /v1/webhooks/github/{projectId}", githubHandler.HandleWebhook)
	mux.HandleFunc("GET /v1/ws/logs/{deploymentId}", h.HandleWSLogsByDeployment)
}
