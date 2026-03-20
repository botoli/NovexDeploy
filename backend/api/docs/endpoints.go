package docs

// Этот файл содержит документацию Swagger для всех эндпоинтов API
// Каждая функция представляет отдельный эндпоинт и не выполняет никакого кода

// ==================== AUTHENTICATION ====================

// Register doc
// @Summary Регистрация нового пользователя
// @Description Создает нового пользователя в системе
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object true "Данные для регистрации (email, password, name)"
// @Success 201 {object} models.APIResponse "Пользователь успешно создан"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Failure 409 {object} models.ErrorResponse "Пользователь уже существует"
// @Router /auth/register [post]
func RegisterEndpoint() {}

// Login doc
// @Summary Вход в систему
// @Description Аутентификация пользователя и получение JWT токена
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object true "Данные для входа (email, password)"
// @Success 200 {object} models.APIResponse "Успешный вход, токен в ответе"
// @Failure 401 {object} models.ErrorResponse "Неверные учетные данные"
// @Router /auth/login [post]
func LoginEndpoint() {}

// Logout doc
// @Summary Выход из системы
// @Description Инвалидация текущей сессии
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Успешный выход"
// @Router /auth/logout [post]
func LogoutEndpoint() {}

// RefreshToken doc
// @Summary Обновление токена
// @Description Получение нового access токена по refresh токену
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object true "Refresh token"
// @Success 200 {object} models.APIResponse "Новый токен"
// @Failure 401 {object} models.ErrorResponse "Недействительный refresh token"
// @Router /auth/refresh [post]
func RefreshTokenEndpoint() {}

// GetMe doc
// @Summary Информация о текущем пользователе
// @Description Возвращает данные авторизованного пользователя
// @Tags Authentication
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Данные пользователя"
// @Failure 401 {object} models.ErrorResponse "Не авторизован"
// @Router /auth/me [get]
func GetMeEndpoint() {}

// ==================== USERS ====================

// GetUserProfile doc
// @Summary Профиль текущего пользователя
// @Description Полная информация о текущем пользователе
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Профиль пользователя"
// @Failure 401 {object} models.ErrorResponse "Не авторизован"
// @Router /users/me [get]
func GetUserProfileEndpoint() {}

// UpdateUserProfile doc
// @Summary Обновление профиля
// @Description Обновляет данные текущего пользователя
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные для обновления (name, email и т.д.)"
// @Success 200 {object} models.APIResponse "Профиль обновлен"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Router /users/me [patch]
func UpdateUserProfileEndpoint() {}

// GetUserByID doc
// @Summary Получение пользователя по ID
// @Description Возвращает публичные данные пользователя
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID пользователя"
// @Success 200 {object} models.APIResponse "Данные пользователя"
// @Failure 404 {object} models.ErrorResponse "Пользователь не найден"
// @Router /users/{id} [get]
func GetUserByIDEndpoint() {}

// ==================== WORKSPACES ====================

// ListWorkspaces doc
// @Summary Список рабочих пространств
// @Description Возвращает все рабочие пространства пользователя
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список рабочих пространств"
// @Router /workspaces [get]
func ListWorkspacesEndpoint() {}

// CreateWorkspace doc
// @Summary Создание рабочего пространства
// @Description Создает новое рабочее пространство
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные рабочего пространства (name, description)"
// @Success 201 {object} models.APIResponse "Рабочее пространство создано"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Router /workspaces [post]
func CreateWorkspaceEndpoint() {}

// GetWorkspace doc
// @Summary Получение рабочего пространства
// @Description Возвращает данные рабочего пространства по ID
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID рабочего пространства"
// @Success 200 {object} models.APIResponse "Данные рабочего пространства"
// @Failure 404 {object} models.ErrorResponse "Не найдено"
// @Router /workspaces/{id} [get]
func GetWorkspaceEndpoint() {}

// UpdateWorkspace doc
// @Summary Обновление рабочего пространства
// @Description Обновляет данные рабочего пространства
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID рабочего пространства"
// @Param request body object true "Данные для обновления"
// @Success 200 {object} models.APIResponse "Обновлено"
// @Failure 404 {object} models.ErrorResponse "Не найдено"
// @Router /workspaces/{id} [patch]
func UpdateWorkspaceEndpoint() {}

// DeleteWorkspace doc
// @Summary Удаление рабочего пространства
// @Description Удаляет рабочее пространство
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID рабочего пространства"
// @Success 200 {object} models.APIResponse "Удалено"
// @Failure 404 {object} models.ErrorResponse "Не найдено"
// @Router /workspaces/{id} [delete]
func DeleteWorkspaceEndpoint() {}

// AddWorkspaceMember doc
// @Summary Добавление участника
// @Description Добавляет пользователя в рабочее пространство
// @Tags Workspaces
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID рабочего пространства"
// @Param request body object true "ID пользователя и роль"
// @Success 200 {object} models.APIResponse "Участник добавлен"
// @Failure 400 {object} models.ErrorResponse "Ошибка"
// @Router /workspaces/{id}/members [post]
func AddWorkspaceMemberEndpoint() {}

// RemoveWorkspaceMember doc
// @Summary Удаление участника
// @Description Удаляет пользователя из рабочего пространства
// @Tags Workspaces
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID рабочего пространства"
// @Param userId path string true "ID пользователя"
// @Success 200 {object} models.APIResponse "Участник удален"
// @Failure 404 {object} models.ErrorResponse "Не найден"
// @Router /workspaces/{id}/members/{userId} [delete]
func RemoveWorkspaceMemberEndpoint() {}

// ==================== PROJECTS ====================

// ListProjects doc
// @Summary Список проектов
// @Description Возвращает все проекты пользователя
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список проектов"
// @Router /projects [get]
func ListProjectsEndpoint() {}

// CreateProject doc
// @Summary Создание проекта
// @Description Создает новый проект
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные проекта (name, description, repository)"
// @Success 201 {object} models.APIResponse "Проект создан"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Router /projects [post]
func CreateProjectEndpoint() {}

// GetProject doc
// @Summary Получение проекта
// @Description Возвращает данные проекта по ID
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Данные проекта"
// @Failure 404 {object} models.ErrorResponse "Проект не найден"
// @Router /projects/{id} [get]
func GetProjectEndpoint() {}

// UpdateProject doc
// @Summary Обновление проекта
// @Description Обновляет данные проекта
// @Tags Projects
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param request body object true "Данные для обновления"
// @Success 200 {object} models.APIResponse "Проект обновлен"
// @Failure 404 {object} models.ErrorResponse "Не найден"
// @Router /projects/{id} [patch]
func UpdateProjectEndpoint() {}

// DeleteProject doc
// @Summary Удаление проекта
// @Description Удаляет проект
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Проект удален"
// @Failure 404 {object} models.ErrorResponse "Не найден"
// @Router /projects/{id} [delete]
func DeleteProjectEndpoint() {}

// GetProjectOverview doc
// @Summary Обзор проекта
// @Description Возвращает общую информацию о проекте
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Обзор проекта"
// @Router /projects/{id}/overview [get]
func GetProjectOverviewEndpoint() {}

// GetProjectStats doc
// @Summary Статистика проекта
// @Description Возвращает статистику проекта
// @Tags Projects
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Статистика"
// @Router /projects/{id}/stats [get]
func GetProjectStatsEndpoint() {}

// ==================== DEPLOYMENTS ====================

// ListDeployments doc
// @Summary Список деплоев
// @Description Возвращает все деплои
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список деплоев"
// @Router /deployments [get]
func ListDeploymentsEndpoint() {}

// CreateDeployment doc
// @Summary Создание деплоя
// @Description Создает новый деплой
// @Tags Deployments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные деплоя (project_id, branch, environment)"
// @Success 201 {object} models.APIResponse "Деплой создан"
// @Failure 400 {object} models.ErrorResponse "Неверные данные"
// @Router /deployments [post]
func CreateDeploymentEndpoint() {}

// GetDeployment doc
// @Summary Получение деплоя
// @Description Возвращает данные деплоя по ID
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID деплоя"
// @Success 200 {object} models.APIResponse "Данные деплоя"
// @Failure 404 {object} models.ErrorResponse "Не найден"
// @Router /deployments/{id} [get]
func GetDeploymentEndpoint() {}

// Redeploy doc
// @Summary Передеплой
// @Description Запускает повторный деплой
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID деплоя"
// @Success 200 {object} models.APIResponse "Передеплой запущен"
// @Router /deployments/{id}/redeploy [post]
func RedeployEndpoint() {}

// RollbackDeployment doc
// @Summary Откат деплоя
// @Description Откатывает деплой к предыдущей версии
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID деплоя"
// @Success 200 {object} models.APIResponse "Откат выполнен"
// @Router /deployments/{id}/rollback [post]
func RollbackDeploymentEndpoint() {}

// CancelDeployment doc
// @Summary Отмена деплоя
// @Description Отменяет текущий деплой
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID деплоя"
// @Success 200 {object} models.APIResponse "Деплой отменен"
// @Router /deployments/{id}/cancel [post]
func CancelDeploymentEndpoint() {}

// GetProjectDeployments doc
// @Summary Деплои проекта
// @Description Возвращает все деплои проекта
// @Tags Deployments
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Список деплоев проекта"
// @Router /projects/{id}/deployments [get]
func GetProjectDeploymentsEndpoint() {}

// ==================== LOGS ====================

// ListLogs doc
// @Summary Список логов
// @Description Возвращает все логи
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список логов"
// @Router /logs [get]
func ListLogsEndpoint() {}

// GetDeploymentLogs doc
// @Summary Логи деплоя
// @Description Возвращает логи конкретного деплоя
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Param deploymentId path string true "ID деплоя"
// @Success 200 {object} models.APIResponse "Логи деплоя"
// @Router /logs/{deploymentId} [get]
func GetDeploymentLogsEndpoint() {}

// GetProjectLogs doc
// @Summary Логи проекта
// @Description Возвращает все логи проекта
// @Tags Logs
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Логи проекта"
// @Router /projects/{id}/logs [get]
func GetProjectLogsEndpoint() {}

// ==================== DOMAINS ====================

// ListDomains doc
// @Summary Список доменов
// @Description Возвращает все домены
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список доменов"
// @Router /domains [get]
func ListDomainsEndpoint() {}

// CreateDomain doc
// @Summary Добавление домена
// @Description Добавляет новый домен
// @Tags Domains
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные домена (name, project_id)"
// @Success 201 {object} models.APIResponse "Домен добавлен"
// @Router /domains [post]
func CreateDomainEndpoint() {}

// DeleteDomain doc
// @Summary Удаление домена
// @Description Удаляет домен
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID домена"
// @Success 200 {object} models.APIResponse "Домен удален"
// @Router /domains/{id} [delete]
func DeleteDomainEndpoint() {}

// AddProjectDomain doc
// @Summary Добавление домена к проекту
// @Description Привязывает домен к проекту
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param request body object true "ID домена"
// @Success 200 {object} models.APIResponse "Домен привязан"
// @Router /projects/{id}/domains [post]
func AddProjectDomainEndpoint() {}

// RemoveProjectDomain doc
// @Summary Отвязка домена от проекта
// @Description Отвязывает домен от проекта
// @Tags Domains
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param domainId path string true "ID домена"
// @Success 200 {object} models.APIResponse "Домен отвязан"
// @Router /projects/{id}/domains/{domainId} [delete]
func RemoveProjectDomainEndpoint() {}

// ==================== ENV VARIABLES ====================

// ListEnvVars doc
// @Summary Список переменных окружения
// @Description Возвращает все переменные окружения проекта
// @Tags Environment
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Список переменных"
// @Router /projects/{id}/env [get]
func ListEnvVarsEndpoint() {}

// CreateEnvVar doc
// @Summary Добавление переменной
// @Description Добавляет переменную окружения в проект
// @Tags Environment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param request body object true "Ключ и значение"
// @Success 201 {object} models.APIResponse "Переменная добавлена"
// @Router /projects/{id}/env [post]
func CreateEnvVarEndpoint() {}

// UpdateEnvVar doc
// @Summary Обновление переменной
// @Description Обновляет переменную окружения
// @Tags Environment
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param key path string true "Ключ переменной"
// @Param request body object true "Новое значение"
// @Success 200 {object} models.APIResponse "Переменная обновлена"
// @Router /projects/{id}/env/{key} [patch]
func UpdateEnvVarEndpoint() {}

// DeleteEnvVar doc
// @Summary Удаление переменной
// @Description Удаляет переменную окружения
// @Tags Environment
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID проекта"
// @Param key path string true "Ключ переменной"
// @Success 200 {object} models.APIResponse "Переменная удалена"
// @Router /projects/{id}/env/{key} [delete]
func DeleteEnvVarEndpoint() {}

// ==================== BOTS ====================

// ListBots doc
// @Summary Список ботов
// @Description Возвращает всех Telegram ботов
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список ботов"
// @Router /bots [get]
func ListBotsEndpoint() {}

// CreateBot doc
// @Summary Создание бота
// @Description Создает нового Telegram бота
// @Tags Bots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные бота (token, name, project_id)"
// @Success 201 {object} models.APIResponse "Бот создан"
// @Router /bots [post]
func CreateBotEndpoint() {}

// GetBot doc
// @Summary Получение бота
// @Description Возвращает данные бота по ID
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Данные бота"
// @Router /bots/{id} [get]
func GetBotEndpoint() {}

// DeleteBot doc
// @Summary Удаление бота
// @Description Удаляет бота
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Бот удален"
// @Router /bots/{id} [delete]
func DeleteBotEndpoint() {}

// RestartBot doc
// @Summary Перезапуск бота
// @Description Перезапускает бота
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Бот перезапущен"
// @Router /bots/{id}/restart [post]
func RestartBotEndpoint() {}

// StopBot doc
// @Summary Остановка бота
// @Description Останавливает бота
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Бот остановлен"
// @Router /bots/{id}/stop [post]
func StopBotEndpoint() {}

// StartBot doc
// @Summary Запуск бота
// @Description Запускает бота
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Бот запущен"
// @Router /bots/{id}/start [post]
func StartBotEndpoint() {}

// SetBotWebhook doc
// @Summary Настройка webhook
// @Description Настраивает webhook для бота
// @Tags Bots
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Param request body object true "URL webhook"
// @Success 200 {object} models.APIResponse "Webhook настроен"
// @Router /bots/{id}/webhook [post]
func SetBotWebhookEndpoint() {}

// GetBotUpdates doc
// @Summary Обновления бота
// @Description Получает последние обновления бота (long polling)
// @Tags Bots
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID бота"
// @Success 200 {object} models.APIResponse "Обновления"
// @Router /bots/{id}/updates [get]
func GetBotUpdatesEndpoint() {}

// ==================== BUILDS ====================

// ListBuilds doc
// @Summary Список сборок
// @Description Возвращает все сборки
// @Tags Builds
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список сборок"
// @Router /builds [get]
func ListBuildsEndpoint() {}

// GetBuild doc
// @Summary Получение сборки
// @Description Возвращает данные сборки по ID
// @Tags Builds
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID сборки"
// @Success 200 {object} models.APIResponse "Данные сборки"
// @Router /builds/{id} [get]
func GetBuildEndpoint() {}

// CreateBuild doc
// @Summary Создание сборки
// @Description Создает новую сборку
// @Tags Builds
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные сборки"
// @Success 201 {object} models.APIResponse "Сборка создана"
// @Router /builds [post]
func CreateBuildEndpoint() {}

// ==================== GIT ====================

// ListGitRepos doc
// @Summary Список репозиториев
// @Description Возвращает подключенные Git репозитории
// @Tags Git
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список репозиториев"
// @Router /git/repos [get]
func ListGitReposEndpoint() {}

// ConnectGitRepo doc
// @Summary Подключение репозитория
// @Description Подключает Git репозиторий
// @Tags Git
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "URL репозитория"
// @Success 200 {object} models.APIResponse "Репозиторий подключен"
// @Router /git/connect [post]
func ConnectGitRepoEndpoint() {}

// DisconnectGitRepo doc
// @Summary Отключение репозитория
// @Description Отключает Git репозиторий
// @Tags Git
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Репозиторий отключен"
// @Router /git/disconnect [post]
func DisconnectGitRepoEndpoint() {}

// GetGitBranches doc
// @Summary Ветки репозитория
// @Description Возвращает ветки Git репозитория
// @Tags Git
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Список веток"
// @Router /git/{projectId}/branches [get]
func GetGitBranchesEndpoint() {}

// GetGitCommits doc
// @Summary Коммиты репозитория
// @Description Возвращает коммиты Git репозитория
// @Tags Git
// @Produce json
// @Security BearerAuth
// @Param projectId path string true "ID проекта"
// @Success 200 {object} models.APIResponse "Список коммитов"
// @Router /git/{projectId}/commits [get]
func GetGitCommitsEndpoint() {}

// ==================== TEMPLATES ====================

// ListTemplates doc
// @Summary Список шаблонов
// @Description Возвращает доступные шаблоны проектов
// @Tags Templates
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список шаблонов"
// @Router /templates [get]
func ListTemplatesEndpoint() {}

// DeployFromTemplate doc
// @Summary Деплой из шаблона
// @Description Создает проект из шаблона
// @Tags Templates
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID шаблона"
// @Success 200 {object} models.APIResponse "Проект создан"
// @Router /templates/{id}/deploy [post]
func DeployFromTemplateEndpoint() {}

// ==================== WEBSOCKET ====================

// WSDeployments doc
// @Summary WebSocket деплои
// @Description Подключается к WebSocket для получения обновлений деплоев
// @Tags WebSocket
// @Router /ws/deployments [get]
func WSDeploymentsEndpoint() {}

// WSLogs doc
// @Summary WebSocket логи
// @Description Подключается к WebSocket для получения всех логов
// @Tags WebSocket
// @Router /ws/logs [get]
func WSLogsEndpoint() {}

// WSProjects doc
// @Summary WebSocket проекты
// @Description Подключается к WebSocket для получения обновлений проектов
// @Tags WebSocket
// @Router /ws/projects [get]
func WSProjectsEndpoint() {}

// WSLogsByDeployment doc
// @Summary WebSocket логи деплоя
// @Description Подключается к WebSocket для получения логов конкретного деплоя
// @Tags WebSocket
// @Param deploymentId path string true "ID деплоя"
// @Router /ws/logs/{deploymentId} [get]
func WSLogsByDeploymentEndpoint() {}

// ==================== QUICK ACTIONS ====================

// QuickDeploy doc
// @Summary Быстрый деплой
// @Description Запускает быстрый деплой
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Параметры деплоя"
// @Success 200 {object} models.APIResponse "Деплой запущен"
// @Router /actions/deploy [post]
func QuickDeployEndpoint() {}

// QuickRedeploy doc
// @Summary Быстрый редеплой
// @Description Запускает быстрый редеплой
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Параметры редеплоя"
// @Success 200 {object} models.APIResponse "Редеплой запущен"
// @Router /actions/redeploy [post]
func QuickRedeployEndpoint() {}

// QuickNewProject doc
// @Summary Быстрый проект
// @Description Быстрое создание проекта
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Данные проекта"
// @Success 200 {object} models.APIResponse "Проект создан"
// @Router /actions/new-project [post]
func QuickNewProjectEndpoint() {}

// QuickImportRepo doc
// @Summary Импорт репозитория
// @Description Быстрый импорт репозитория
// @Tags Actions
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "URL репозитория"
// @Success 200 {object} models.APIResponse "Репозиторий импортирован"
// @Router /actions/import-repo [post]
func QuickImportRepoEndpoint() {}

// ==================== NOTIFICATIONS ====================

// ListNotifications doc
// @Summary Список уведомлений
// @Description Возвращает уведомления пользователя
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Список уведомлений"
// @Router /notifications [get]
func ListNotificationsEndpoint() {}

// MarkNotificationsRead doc
// @Summary Отметить как прочитанные
// @Description Отмечает все уведомления как прочитанные
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse "Уведомления отмечены"
// @Router /notifications/mark-read [post]
func MarkNotificationsReadEndpoint() {}

// DeleteNotification doc
// @Summary Удаление уведомления
// @Description Удаляет уведомление
// @Tags Notifications
// @Produce json
// @Security BearerAuth
// @Param id path string true "ID уведомления"
// @Success 200 {object} models.APIResponse "Уведомление удалено"
// @Router /notifications/{id} [delete]
func DeleteNotificationEndpoint() {}