package handlers

import (
	"localVercel/utils"
	"net/http"
	"strings"
)

// HandleAction выполняет действие
// @Summary Выполнить действие
// @Description Выполняет указанное действие без параметров пути
// @Tags Actions
// @Accept json
// @Produce json
// @Param request body object false "Параметры действия"
// @Success 200 {object} models.APIResponse "Действие выполнено"
// @Router /{action} [post]
func (h *Handler) HandleAction(action string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := utils.ReadPayload(r)
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, action+" executed", map[string]interface{}{
			"action":  action,
			"payload": payload,
		}))
	}
}

// HandleActionWithPath выполняет действие с параметрами пути
// @Summary Выполнить действие с параметрами
// @Description Выполняет указанное действие с параметрами из пути
// @Tags Actions
// @Accept json
// @Produce json
// @Param id path string true "ID ресурса"
// @Param request body object false "Дополнительные параметры"
// @Success 200 {object} models.APIResponse "Действие выполнено"
// @Router /{resource}/{id}/{action} [post]
func (h *Handler) HandleActionWithPath(action string, pathParams []string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		payload := utils.ReadPayload(r)
		data := map[string]interface{}{
			"action":  action,
			"payload": payload,
		}
		for _, p := range pathParams {
			data[p] = r.PathValue(p)
		}
		utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, action+" executed", data))
	}
}

// HandleSearch выполняет поиск
// @Summary Поиск по всем ресурсам
// @Description Ищет по проектам, деплоям, логам и ботам
// @Tags Search
// @Produce json
// @Param q query string true "Поисковый запрос"
// @Success 200 {object} models.APIResponse "Результаты поиска"
// @Router /search [get]
func (h *Handler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	q := strings.TrimSpace(r.URL.Query().Get("q"))

	utils.WriteJSON(w, http.StatusOK, h.jsonResponse(true, "search completed", map[string]interface{}{
		"query": q,
		"results": map[string]interface{}{
			"projects":    []interface{}{},
			"deployments": []interface{}{},
			"logs":        []interface{}{},
			"bots":        []interface{}{},
		},
	}))
}