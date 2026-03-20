package handlers

import (
	"localVercel/models"
	"localVercel/utils"
	"time"

	"github.com/gorilla/websocket"
)

type Handler struct {
	startedAt time.Time
	upgrader  *websocket.Upgrader
}

func New(startedAt time.Time, upgrader *websocket.Upgrader) *Handler {
	return &Handler{
		startedAt: startedAt,
		upgrader:  upgrader,
	}
}

// Base response helper
func (h *Handler) jsonResponse(ok bool, message string, data interface{}) models.APIResponse {
	return models.APIResponse{
		OK:        ok,
		Message:   message,
		Data:      data,
		Timestamp: utils.NowISO(),
	}
}