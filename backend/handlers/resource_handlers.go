package handlers

import (
	"localVercel/db"
	"localVercel/models"
	"localVercel/utils"
	"net/http"
	"time"
)

// HandleList возвращает список ресурсов
// @Summary Получить список ресурсов
// @Description Возвращает список всех ресурсов указанного типа
// @Tags Resources
// @Produce json
// @Success 200 {object} models.APIResponse "Список ресурсов"
// @Router /{resource} [get]
func (h *Handler) HandleList(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" list", map[string]interface{}{
			"resource": resource,
			"items":    []interface{}{},
		}))
	}
}

// HandleCreate создает новый ресурс
// @Summary Создать ресурс
// @Description Создает новый ресурс указанного типа
// @Tags Resources
// @Accept json
// @Produce json
// @Param request body object true "Данные для создания"
// @Success 201 {object} models.APIResponse "Ресурс создан"
// @Router /{resource} [post]
func (h *Handler) HandleCreate(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := utils.ReadPayload(r)
		utils.WriteJSON(w, http.StatusCreated, h.jsonResponse(true, resource+" created", map[string]interface{}{
			"resource": resource,
			"payload":  payload,
		}))
	}
}

// HandleGet получает ресурс текущего пользователя
// @Summary Получить ресурс
// @Description Возвращает ресурс текущего пользователя
// @Tags Resources
// @Produce json
// @Success 200 {object} models.APIResponse "Ресурс получен"
// @Router /{resource}/me [get]
func (h *Handler) HandleGet(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Для /auth/me нужно вернуть реальные данные пользователя из сессии/БД
		if resource == "auth.me" {
			// Получаем токен из cookie
			cookie, err := r.Cookie("session_token")
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Not authenticated", nil))
				return
			}

			// Ищем сессию в БД
			var session models.Session
			if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
				return
			}

			// Получаем пользователя
			var user models.User
			if err := db.DB.First(&user, "id = ?", session.UserID).Error; err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "User not found", nil))
				return
			}

			// Возвращаем данные пользователя
			utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "auth.me fetched", map[string]interface{}{
				"id":           user.ID,
				"email":        user.Email,
				"name":         user.Name,
				"avatar_url":   user.AvatarURL,
				"github_login": user.GitHubLogin,
				"github_id":    user.GitHubID,
				"last_login_at": user.LastLoginAt,
			}))
			return
		}
		
		// Для остальных ресурсов
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" fetched", map[string]interface{}{
			"resource": resource,
		}))
	}
}

// HandlePatch обновляет ресурс
// @Summary Обновить ресурс
// @Description Обновляет ресурс текущего пользователя
// @Tags Resources
// @Accept json
// @Produce json
// @Param request body object true "Данные для обновления"
// @Success 200 {object} models.APIResponse "Ресурс обновлен"
// @Router /{resource}/me [patch]
func (h *Handler) HandlePatch(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := utils.ReadPayload(r)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" updated", map[string]interface{}{
			"resource": resource,
			"payload":  payload,
		}))
	}
}

// HandleGetByID получает ресурс по ID
// @Summary Получить ресурс по ID
// @Description Возвращает ресурс по его идентификатору
// @Tags Resources
// @Produce json
// @Param id path string true "ID ресурса"
// @Success 200 {object} models.APIResponse "Ресурс получен"
// @Router /{resource}/{id} [get]
func (h *Handler) HandleGetByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" fetched", map[string]interface{}{
			"resource": resource,
			param:      id,
		}))
	}
}

// HandlePatchByID обновляет ресурс по ID
// @Summary Обновить ресурс по ID
// @Description Обновляет существующий ресурс по его идентификатору
// @Tags Resources
// @Accept json
// @Produce json
// @Param id path string true "ID ресурса"
// @Param request body object true "Данные для обновления"
// @Success 200 {object} models.APIResponse "Ресурс обновлен"
// @Router /{resource}/{id} [patch]
func (h *Handler) HandlePatchByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		payload := utils.ReadPayload(r)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" updated", map[string]interface{}{
			"resource": resource,
			param:      id,
			"payload":  payload,
		}))
	}
}

// HandleDeleteByID удаляет ресурс по ID
// @Summary Удалить ресурс по ID
// @Description Удаляет существующий ресурс по его идентификатору
// @Tags Resources
// @Produce json
// @Param id path string true "ID ресурса"
// @Success 200 {object} models.APIResponse "Ресурс удален"
// @Router /{resource}/{id} [delete]
func (h *Handler) HandleDeleteByID(resource, param string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue(param)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" deleted", map[string]interface{}{
			"resource": resource,
			param:      id,
		}))
	}
}