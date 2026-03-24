package handlers

import (
	"localVercel/utils"
	"net/http"
	"time"
)

// HandleHealth проверка здоровья сервиса
// @Summary Проверка здоровья
// @Description Возвращает статус сервиса и время работы
// @Tags System
// @Produce json
// @Success 200 {object} models.APIResponse "Сервис работает"
// @Router /health [get]
func (h *Handler) HandleHealth(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "healthy", map[string]interface{}{
		"uptime_seconds": int(time.Since(h.startedAt).Seconds()),
	}))
}
