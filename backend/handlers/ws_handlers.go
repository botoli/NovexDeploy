package handlers

import (
	"localVercel/utils"
	"log"
	"net/http"
	"time"
)

func (h *Handler) handleWS(w http.ResponseWriter, r *http.Request, channel, scope string) {
	conn, err := h.upgrader.Upgrade(w, r, nil)
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
				"timestamp": utils.NowISO(),
			}
			if err := conn.WriteJSON(msg); err != nil {
				return
			}
		}
	}
}

func (h *Handler) HandleWSDeployments(w http.ResponseWriter, r *http.Request) {
	h.handleWS(w, r, "deployments", "all")
}

func (h *Handler) HandleWSLogs(w http.ResponseWriter, r *http.Request) {
	h.handleWS(w, r, "logs", "all")
}

func (h *Handler) HandleWSProjects(w http.ResponseWriter, r *http.Request) {
	h.handleWS(w, r, "projects", "all")
}

func (h *Handler) HandleWSLogsByDeployment(w http.ResponseWriter, r *http.Request) {
	h.handleWS(w, r, "logs", r.PathValue("deploymentId"))
}