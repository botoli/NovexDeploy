package handlers

import (
	"localVercel/db"
	"localVercel/models"
	"localVercel/utils"
	"net/http"
	"time"
)

func (h *Handler) HandleGet(resource string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if resource == "auth.me" {
			cookie, err := r.Cookie("session_token")
			if err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Not authenticated", nil))
				return
			}

			var session models.Session
			if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
				utils.WriteJSON(w, http.StatusUnauthorized, h.jsonResponse(false, "Invalid session", nil))
				return
			}

			var user models.User
			if err := db.DB.First(&user, "id = ?", session.UserID).Error; err != nil {
				utils.WriteJSON(w, http.StatusInternalServerError, h.jsonResponse(false, "User not found", nil))
				return
			}

			// Возвращаем данные пользователя
			utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "auth.me fetched", map[string]interface{}{
				"id":            user.ID,
				"email":         user.Email,
				"name":          user.Name,
				"avatar_url":    user.AvatarURL,
				"github_login":  user.GitHubLogin,
				"github_id":     user.GitHubID,
				"last_login_at": user.LastLoginAt,
			}))
			return
		}

		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, resource+" fetched", map[string]interface{}{
			"resource": resource,
		}))
	}
}
