package handlers

import (
	"errors"
	"localVercel/db"
	"localVercel/models"
	"net/http"
	"time"
)

func currentUserFromSession(r *http.Request) (*models.User, error) {
	cookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, errors.New("missing session cookie")
	}

	var session models.Session
	if err := db.DB.Where("token = ? AND expires_at > ?", cookie.Value, time.Now()).First(&session).Error; err != nil {
		return nil, errors.New("invalid session")
	}

	var user models.User
	if err := db.DB.First(&user, "id = ?", session.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	return &user, nil
}
