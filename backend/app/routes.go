package app

import (
	"localVercel/handlers"
	"net/http"
)

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	h := handlers.New(a.startedAt, &a.upgrader)
	githubHandler := handlers.NewGitHubHandler(h)

	// GitHub OAuth
	mux.HandleFunc("GET /auth/github/login", githubHandler.HandleGitHubLogin)
	mux.HandleFunc("GET /auth/github/callback", githubHandler.HandleGitHubCallback)

	// GitHub Webhook (публичный, без аутентификации)
	mux.HandleFunc("POST /webhook/github/{projectId}", githubHandler.HandleWebhook)

	// Управление репозиториями и билдами
	mux.HandleFunc("POST /projects/{projectId}/github/repo", githubHandler.HandleConnectRepo)
	mux.HandleFunc("GET /projects/{projectId}/builds", githubHandler.HandleListBuilds)
	mux.HandleFunc("GET /projects/{projectId}/builds/{buildId}", githubHandler.HandleGetBuild)

	// 1. Auth / User
	mux.HandleFunc("POST /auth/register", h.HandleCreate("auth.register"))
	mux.HandleFunc("POST /auth/login", h.HandleAction("auth.login"))
	mux.HandleFunc("POST /auth/logout", h.HandleAction("auth.logout"))
	mux.HandleFunc("POST /auth/refresh", h.HandleAction("auth.refresh"))
	mux.HandleFunc("GET /auth/me", h.HandleGet("auth.me"))

	// 2. Users
	mux.HandleFunc("GET /users/me", h.HandleGet("users.me"))
	mux.HandleFunc("PATCH /users/me", h.HandlePatch("users.me"))
	mux.HandleFunc("GET /users/{id}", h.HandleGetByID("users", "id"))

	// 3. Workspaces / Teams
	mux.HandleFunc("GET /workspaces", h.HandleList("workspaces"))
	mux.HandleFunc("POST /workspaces", h.HandleCreate("workspaces"))
	mux.HandleFunc("GET /workspaces/{id}", h.HandleGetByID("workspaces", "id"))
	mux.HandleFunc("PATCH /workspaces/{id}", h.HandlePatchByID("workspaces", "id"))
	mux.HandleFunc("DELETE /workspaces/{id}", h.HandleDeleteByID("workspaces", "id"))
	mux.HandleFunc("POST /workspaces/{id}/members", h.HandleActionWithPath("workspaces.members.add", []string{"id"}))
	mux.HandleFunc("DELETE /workspaces/{id}/members/{userId}", h.HandleActionWithPath("workspaces.members.remove", []string{"id", "userId"}))

	// 4. Projects
	mux.HandleFunc("GET /projects", h.HandleList("projects"))
	mux.HandleFunc("POST /projects", h.HandleCreate("projects"))
	mux.HandleFunc("GET /projects/{id}", h.HandleGetByID("projects", "id"))
	mux.HandleFunc("PATCH /projects/{id}", h.HandlePatchByID("projects", "id"))
	mux.HandleFunc("DELETE /projects/{id}", h.HandleDeleteByID("projects", "id"))
	mux.HandleFunc("GET /projects/{id}/overview", h.HandleActionWithPath("projects.overview", []string{"id"}))
	mux.HandleFunc("GET /projects/{id}/stats", h.HandleActionWithPath("projects.stats", []string{"id"}))

	// 5. Deployments
	mux.HandleFunc("GET /deployments", h.HandleList("deployments"))
	mux.HandleFunc("POST /deployments", h.HandleCreate("deployments"))
	mux.HandleFunc("GET /deployments/{id}", h.HandleGetByID("deployments", "id"))
	mux.HandleFunc("POST /deployments/{id}/redeploy", h.HandleActionWithPath("deployments.redeploy", []string{"id"}))
	mux.HandleFunc("POST /deployments/{id}/rollback", h.HandleActionWithPath("deployments.rollback", []string{"id"}))
	mux.HandleFunc("POST /deployments/{id}/cancel", h.HandleActionWithPath("deployments.cancel", []string{"id"}))
	mux.HandleFunc("GET /projects/{id}/deployments", h.HandleActionWithPath("projects.deployments", []string{"id"}))

	// 6. Logs
	mux.HandleFunc("GET /logs", h.HandleList("logs"))
	mux.HandleFunc("GET /logs/{deploymentId}", h.HandleGetByID("logs", "deploymentId"))
	mux.HandleFunc("GET /projects/{id}/logs", h.HandleActionWithPath("projects.logs", []string{"id"}))
	mux.HandleFunc("GET /ws/logs/{deploymentId}", h.HandleWSLogsByDeployment)

	// 7. Domains
	mux.HandleFunc("GET /domains", h.HandleList("domains"))
	mux.HandleFunc("POST /domains", h.HandleCreate("domains"))
	mux.HandleFunc("DELETE /domains/{id}", h.HandleDeleteByID("domains", "id"))
	mux.HandleFunc("POST /projects/{id}/domains", h.HandleActionWithPath("projects.domains.add", []string{"id"}))
	mux.HandleFunc("DELETE /projects/{id}/domains/{domainId}", h.HandleActionWithPath("projects.domains.remove", []string{"id", "domainId"}))

	// 8. Env Variables
	mux.HandleFunc("GET /projects/{id}/env", h.HandleActionWithPath("projects.env.list", []string{"id"}))
	mux.HandleFunc("POST /projects/{id}/env", h.HandleActionWithPath("projects.env.create", []string{"id"}))
	mux.HandleFunc("PATCH /projects/{id}/env/{key}", h.HandleActionWithPath("projects.env.update", []string{"id", "key"}))
	mux.HandleFunc("DELETE /projects/{id}/env/{key}", h.HandleActionWithPath("projects.env.delete", []string{"id", "key"}))

	// 9. Telegram Bots
	mux.HandleFunc("GET /bots", h.HandleList("bots"))
	mux.HandleFunc("POST /bots", h.HandleCreate("bots"))
	mux.HandleFunc("GET /bots/{id}", h.HandleGetByID("bots", "id"))
	mux.HandleFunc("DELETE /bots/{id}", h.HandleDeleteByID("bots", "id"))
	mux.HandleFunc("POST /bots/{id}/restart", h.HandleActionWithPath("bots.restart", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/stop", h.HandleActionWithPath("bots.stop", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/start", h.HandleActionWithPath("bots.start", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/webhook", h.HandleActionWithPath("bots.webhook", []string{"id"}))
	mux.HandleFunc("GET /bots/{id}/updates", h.HandleActionWithPath("bots.updates", []string{"id"}))

	// 10. Builds
	mux.HandleFunc("GET /builds", h.HandleList("builds"))
	mux.HandleFunc("GET /builds/{id}", h.HandleGetByID("builds", "id"))
	mux.HandleFunc("POST /builds", h.HandleCreate("builds"))

	// 11. Metrics / Monitoring
	mux.HandleFunc("GET /metrics/system", h.HandleMetricsSystem)
	mux.HandleFunc("GET /metrics/projects/{id}", h.HandleActionWithPath("metrics.project", []string{"id"}))
	mux.HandleFunc("GET /metrics/deployments/{id}", h.HandleActionWithPath("metrics.deployment", []string{"id"}))

	// 12. Notifications
	mux.HandleFunc("GET /notifications", h.HandleList("notifications"))
	mux.HandleFunc("POST /notifications/mark-read", h.HandleAction("notifications.mark_read"))
	mux.HandleFunc("DELETE /notifications/{id}", h.HandleDeleteByID("notifications", "id"))

	// 13. Search
	mux.HandleFunc("GET /search", h.HandleSearch)

	// 14. Quick Actions
	mux.HandleFunc("POST /actions/deploy", h.HandleAction("actions.deploy"))
	mux.HandleFunc("POST /actions/redeploy", h.HandleAction("actions.redeploy"))
	mux.HandleFunc("POST /actions/new-project", h.HandleAction("actions.new_project"))
	mux.HandleFunc("POST /actions/import-repo", h.HandleAction("actions.import_repo"))

	// 15. Git Integration
	mux.HandleFunc("GET /git/repos", h.HandleAction("git.repos"))
	mux.HandleFunc("POST /git/connect", h.HandleAction("git.connect"))
	mux.HandleFunc("POST /git/disconnect", h.HandleAction("git.disconnect"))
	mux.HandleFunc("GET /git/{projectId}/branches", h.HandleActionWithPath("git.branches", []string{"projectId"}))
	mux.HandleFunc("GET /git/{projectId}/commits", h.HandleActionWithPath("git.commits", []string{"projectId"}))

	// 16. Templates
	mux.HandleFunc("GET /templates", h.HandleList("templates"))
	mux.HandleFunc("POST /templates/{id}/deploy", h.HandleActionWithPath("templates.deploy", []string{"id"}))

	// 17. Realtime
	mux.HandleFunc("GET /ws/deployments", h.HandleWSDeployments)
	mux.HandleFunc("GET /ws/logs", h.HandleWSLogs)
	mux.HandleFunc("GET /ws/projects", h.HandleWSProjects)

	// 18. System
	mux.HandleFunc("GET /health", h.HandleHealth)
	mux.HandleFunc("GET /status", h.HandleStatus)
	mux.HandleFunc("GET /usage", h.HandleUsage)
}