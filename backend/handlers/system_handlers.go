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

// HandleStatus информация о сервисе
// @Summary Информация о сервисе
// @Description Возвращает информацию о версии и статусе сервиса
// @Tags System
// @Produce json
// @Success 200 {object} models.APIResponse "Информация о сервисе"
// @Router /status [get]
func (h *Handler) HandleStatus(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "status", map[string]interface{}{
		"service": "localVercel-backend",
		"version": "v0.1.0",
		"uptime":  time.Since(h.startedAt).Round(time.Second).String(),
	}))
}

// HandleUsage статистика использования
// @Summary Статистика использования
// @Description Возвращает статистику использования ресурсов
// @Tags System
// @Produce json
// @Success 200 {object} models.APIResponse "Статистика использования"
// @Router /usage [get]
func (h *Handler) HandleUsage(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "usage", map[string]interface{}{
		"projects":    0,
		"deployments": 0,
		"bots":        0,
	}))
}

// HandleMetricsSystem системные метрики
// @Summary Системные метрики
// @Description Возвращает метрики системы (CPU, RAM, Disk)
// @Tags Metrics
// @Produce json
// @Success 200 {object} models.APIResponse "Системные метрики"
// @Router /metrics/system [get]
func (h *Handler) HandleMetricsSystem(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "system metrics", map[string]interface{}{
		"cpu_percent": 0.0,
		"ram_mb":      0,
		"disk_mb":     0,
	}))
}