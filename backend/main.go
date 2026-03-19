package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

type App struct {
	startedAt time.Time
	upgrader  websocket.Upgrader
}

type APIResponse struct {
	OK        bool        `json:"ok"`
	Message   string      `json:"message,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp string      `json:"timestamp"`
}

func main() {
	app := &App{
		startedAt: time.Now(),
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}

	mux := http.NewServeMux()
	app.registerRoutes(mux)

	addr := ":8888"
	log.Printf("backend listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, withLogging(mux)); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func (a *App) registerRoutes(mux *http.ServeMux) {
	// 1. Auth / User
	mux.HandleFunc("POST /auth/register", a.handleCreate("auth.register"))
	mux.HandleFunc("POST /auth/login", a.handleAction("auth.login"))
	mux.HandleFunc("POST /auth/logout", a.handleAction("auth.logout"))
	mux.HandleFunc("POST /auth/refresh", a.handleAction("auth.refresh"))
	mux.HandleFunc("GET /auth/me", a.handleGet("auth.me"))

	// 2. Users
	mux.HandleFunc("GET /users/me", a.handleGet("users.me"))
	mux.HandleFunc("PATCH /users/me", a.handlePatch("users.me"))
	mux.HandleFunc("GET /users/{id}", a.handleGetByID("users", "id"))

	// 3. Workspaces / Teams
	mux.HandleFunc("GET /workspaces", a.handleList("workspaces"))
	mux.HandleFunc("POST /workspaces", a.handleCreate("workspaces"))
	mux.HandleFunc("GET /workspaces/{id}", a.handleGetByID("workspaces", "id"))
	mux.HandleFunc("PATCH /workspaces/{id}", a.handlePatchByID("workspaces", "id"))
	mux.HandleFunc("DELETE /workspaces/{id}", a.handleDeleteByID("workspaces", "id"))
	mux.HandleFunc("POST /workspaces/{id}/members", a.handleActionWithPath("workspaces.members.add", []string{"id"}))
	mux.HandleFunc("DELETE /workspaces/{id}/members/{userId}", a.handleActionWithPath("workspaces.members.remove", []string{"id", "userId"}))

	// 4. Projects
	mux.HandleFunc("GET /projects", a.handleList("projects"))
	mux.HandleFunc("POST /projects", a.handleCreate("projects"))
	mux.HandleFunc("GET /projects/{id}", a.handleGetByID("projects", "id"))
	mux.HandleFunc("PATCH /projects/{id}", a.handlePatchByID("projects", "id"))
	mux.HandleFunc("DELETE /projects/{id}", a.handleDeleteByID("projects", "id"))
	mux.HandleFunc("GET /projects/{id}/overview", a.handleActionWithPath("projects.overview", []string{"id"}))
	mux.HandleFunc("GET /projects/{id}/stats", a.handleActionWithPath("projects.stats", []string{"id"}))

	// 5. Deployments
	mux.HandleFunc("GET /deployments", a.handleList("deployments"))
	mux.HandleFunc("POST /deployments", a.handleCreate("deployments"))
	mux.HandleFunc("GET /deployments/{id}", a.handleGetByID("deployments", "id"))
	mux.HandleFunc("POST /deployments/{id}/redeploy", a.handleActionWithPath("deployments.redeploy", []string{"id"}))
	mux.HandleFunc("POST /deployments/{id}/rollback", a.handleActionWithPath("deployments.rollback", []string{"id"}))
	mux.HandleFunc("POST /deployments/{id}/cancel", a.handleActionWithPath("deployments.cancel", []string{"id"}))
	mux.HandleFunc("GET /projects/{id}/deployments", a.handleActionWithPath("projects.deployments", []string{"id"}))

	// 6. Logs
	mux.HandleFunc("GET /logs", a.handleList("logs"))
	mux.HandleFunc("GET /logs/{deploymentId}", a.handleGetByID("logs", "deploymentId"))
	mux.HandleFunc("GET /projects/{id}/logs", a.handleActionWithPath("projects.logs", []string{"id"}))
	mux.HandleFunc("GET /ws/logs/{deploymentId}", a.handleWSLogsByDeployment)

	// 7. Domains
	mux.HandleFunc("GET /domains", a.handleList("domains"))
	mux.HandleFunc("POST /domains", a.handleCreate("domains"))
	mux.HandleFunc("DELETE /domains/{id}", a.handleDeleteByID("domains", "id"))
	mux.HandleFunc("POST /projects/{id}/domains", a.handleActionWithPath("projects.domains.add", []string{"id"}))
	mux.HandleFunc("DELETE /projects/{id}/domains/{domainId}", a.handleActionWithPath("projects.domains.remove", []string{"id", "domainId"}))

	// 8. Env Variables
	mux.HandleFunc("GET /projects/{id}/env", a.handleActionWithPath("projects.env.list", []string{"id"}))
	mux.HandleFunc("POST /projects/{id}/env", a.handleActionWithPath("projects.env.create", []string{"id"}))
	mux.HandleFunc("PATCH /projects/{id}/env/{key}", a.handleActionWithPath("projects.env.update", []string{"id", "key"}))
	mux.HandleFunc("DELETE /projects/{id}/env/{key}", a.handleActionWithPath("projects.env.delete", []string{"id", "key"}))

	// 9. Telegram Bots
	mux.HandleFunc("GET /bots", a.handleList("bots"))
	mux.HandleFunc("POST /bots", a.handleCreate("bots"))
	mux.HandleFunc("GET /bots/{id}", a.handleGetByID("bots", "id"))
	mux.HandleFunc("DELETE /bots/{id}", a.handleDeleteByID("bots", "id"))
	mux.HandleFunc("POST /bots/{id}/restart", a.handleActionWithPath("bots.restart", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/stop", a.handleActionWithPath("bots.stop", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/start", a.handleActionWithPath("bots.start", []string{"id"}))
	mux.HandleFunc("POST /bots/{id}/webhook", a.handleActionWithPath("bots.webhook", []string{"id"}))
	mux.HandleFunc("GET /bots/{id}/updates", a.handleActionWithPath("bots.updates", []string{"id"}))

	// 10. Builds
	mux.HandleFunc("GET /builds", a.handleList("builds"))
	mux.HandleFunc("GET /builds/{id}", a.handleGetByID("builds", "id"))
	mux.HandleFunc("POST /builds", a.handleCreate("builds"))

	// 11. Metrics / Monitoring
	mux.HandleFunc("GET /metrics/system", a.handleMetricsSystem)
	mux.HandleFunc("GET /metrics/projects/{id}", a.handleActionWithPath("metrics.project", []string{"id"}))
	mux.HandleFunc("GET /metrics/deployments/{id}", a.handleActionWithPath("metrics.deployment", []string{"id"}))

	// 12. Notifications
	mux.HandleFunc("GET /notifications", a.handleList("notifications"))
	mux.HandleFunc("POST /notifications/mark-read", a.handleAction("notifications.mark_read"))
	mux.HandleFunc("DELETE /notifications/{id}", a.handleDeleteByID("notifications", "id"))

	// 13. Search
	mux.HandleFunc("GET /search", a.handleSearch)

	// 14. Quick Actions
	mux.HandleFunc("POST /actions/deploy", a.handleAction("actions.deploy"))
	mux.HandleFunc("POST /actions/redeploy", a.handleAction("actions.redeploy"))
	mux.HandleFunc("POST /actions/new-project", a.handleAction("actions.new_project"))
	mux.HandleFunc("POST /actions/import-repo", a.handleAction("actions.import_repo"))

	// 15. Git Integration (without GitHub webhook)
	mux.HandleFunc("GET /git/repos", a.handleAction("git.repos"))
	mux.HandleFunc("POST /git/connect", a.handleAction("git.connect"))
	mux.HandleFunc("POST /git/disconnect", a.handleAction("git.disconnect"))
	mux.HandleFunc("GET /git/{projectId}/branches", a.handleActionWithPath("git.branches", []string{"projectId"}))
	mux.HandleFunc("GET /git/{projectId}/commits", a.handleActionWithPath("git.commits", []string{"projectId"}))

	// 16. Templates
	mux.HandleFunc("GET /templates", a.handleList("templates"))
	mux.HandleFunc("POST /templates/{id}/deploy", a.handleActionWithPath("templates.deploy", []string{"id"}))

	// 17. Realtime
	mux.HandleFunc("GET /ws/deployments", a.handleWSDeployments)
	mux.HandleFunc("GET /ws/logs", a.handleWSLogs)
	mux.HandleFunc("GET /ws/projects", a.handleWSProjects)

	// 18. System
	mux.HandleFunc("GET /health", a.handleHealth)
	mux.HandleFunc("GET /status", a.handleStatus)
	mux.HandleFunc("GET /usage", a.handleUsage)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s (%s)", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}

func (a *App) handleList(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " list",
			Data: map[string]interface{}{
				"resource": resource,
				"items":    []interface{}{},
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleCreate(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := readPayload(r)
		writeJSON(w, http.StatusCreated, APIResponse{
			OK:      true,
			Message: resource + " created",
			Data: map[string]interface{}{
				"resource": resource,
				"payload":  payload,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleGet(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " fetched",
			Data: map[string]interface{}{
				"resource": resource,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handlePatch(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := readPayload(r)
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " updated",
			Data: map[string]interface{}{
				"resource": resource,
				"payload":  payload,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleAction(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := readPayload(r)
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: action + " executed",
			Data: map[string]interface{}{
				"action":  action,
				"payload": payload,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleGetByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " fetched",
			Data: map[string]interface{}{
				"resource": resource,
				param:       id,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handlePatchByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		payload := readPayload(r)
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " updated",
			Data: map[string]interface{}{
				"resource": resource,
				param:       id,
				"payload":  payload,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleDeleteByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		writeJSON(w, http.StatusOK, APIResponse{
			OK:      true,
			Message: resource + " deleted",
			Data: map[string]interface{}{
				"resource": resource,
				param:       id,
			},
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleActionWithPath(action string, pathParams []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := readPayload(r)
		data := map[string]interface{}{
			"action":  action,
			"payload": payload,
		}
		for _, p := range pathParams {
			data[p] = r.PathValue(p)
		}
		writeJSON(w, http.StatusOK, APIResponse{
			OK:        true,
			Message:   action + " executed",
			Data:      data,
			Timestamp: nowISO(),
		})
	}
}

func (a *App) handleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	writeJSON(w, http.StatusOK, APIResponse{
		OK:      true,
		Message: "search completed",
		Data: map[string]interface{}{
			"query": q,
			"results": map[string]interface{}{
				"projects":    []interface{}{},
				"deployments": []interface{}{},
				"logs":        []interface{}{},
				"bots":        []interface{}{},
			},
		},
		Timestamp: nowISO(),
	})
}

func (a *App) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		OK:      true,
		Message: "healthy",
		Data: map[string]interface{}{
			"uptime_seconds": int(time.Since(a.startedAt).Seconds()),
		},
		Timestamp: nowISO(),
	})
}

func (a *App) handleStatus(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		OK:      true,
		Message: "status",
		Data: map[string]interface{}{
			"service": "localVercel-backend",
			"version": "v0.1.0",
			"uptime":  time.Since(a.startedAt).Round(time.Second).String(),
		},
		Timestamp: nowISO(),
	})
}

func (a *App) handleUsage(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		OK:      true,
		Message: "usage",
		Data: map[string]interface{}{
			"projects":    0,
			"deployments": 0,
			"bots":        0,
		},
		Timestamp: nowISO(),
	})
}

func (a *App) handleMetricsSystem(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, APIResponse{
		OK:      true,
		Message: "system metrics",
		Data: map[string]interface{}{
			"cpu_percent": 0.0,
			"ram_mb":      0,
			"disk_mb":     0,
		},
		Timestamp: nowISO(),
	})
}

func (a *App) handleWSDeployments(w http.ResponseWriter, r *http.Request) {
	a.handleWS(w, r, "deployments", "all")
}

func (a *App) handleWSLogs(w http.ResponseWriter, r *http.Request) {
	a.handleWS(w, r, "logs", "all")
}

func (a *App) handleWSProjects(w http.ResponseWriter, r *http.Request) {
	a.handleWS(w, r, "projects", "all")
}

func (a *App) handleWSLogsByDeployment(w http.ResponseWriter, r *http.Request) {
	a.handleWS(w, r, "logs", r.PathValue("deploymentId"))
}

func (a *App) handleWS(w http.ResponseWriter, r *http.Request, channel, scope string) {
	conn, err := a.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			msg := map[string]interface{}{
				"type":      "event",
				"channel":   channel,
				"scope":     scope,
				"timestamp": nowISO(),
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, payload APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func readPayload(r *http.Request) map[string]interface{} {
	defer r.Body.Close()
	if r.Body == nil {
		return map[string]interface{}{}
	}
	var payload map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return map[string]interface{}{}
	}
	return payload
}

func nowISO() string {
	return time.Now().UTC().Format(time.RFC3339)
}
