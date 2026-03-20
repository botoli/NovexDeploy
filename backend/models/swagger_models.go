package models

// AuthRegisterRequest представляет запрос на регистрацию
// @Description Данные для регистрации нового пользователя
type AuthRegisterRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required"`
	Password string `json:"password" example:"secure123" binding:"required"`
	Name     string `json:"name" example:"John Doe" binding:"required"`
}

// AuthLoginRequest представляет запрос на вход
// @Description Данные для входа в систему
type AuthLoginRequest struct {
	Email    string `json:"email" example:"user@example.com" binding:"required"`
	Password string `json:"password" example:"secure123" binding:"required"`
}

// ErrorResponse представляет ответ с ошибкой
// @Description Стандартный ответ при ошибке
type ErrorResponse struct {
	OK        bool   `json:"ok" example:"false"`
	Message   string `json:"message" example:"error description"`
	Timestamp string `json:"timestamp" example:"2024-01-01T00:00:00Z"`
}