# Novex Deploy Backend

Платформа деплоя для backend/telegram проектов с автодеплоем из GitHub.

## Что теперь работает

- Подключение GitHub и импорт репозитория
- Автодеплой по push через GitHub webhook
- Выбор `root_dir` (подпапка в репозитории для сборки)
- Запуск runtime каждого проекта в изолированном Docker контейнере
- Отдельная БД PostgreSQL на проект в Docker контейнере
- Автоматическая запись `DATABASE_URL` в env проекта после provisioning БД

## Ключевые endpoint'ы

- `POST /v1/projects/{projectId}/repo/connect`
- `POST /v1/projects/{projectId}/deployments`
- `GET /v1/projects/{projectId}/deployments`
- `GET /v1/deployments/{deploymentId}/logs`
- `POST /v1/projects/{projectId}/runtime/start`
- `POST /v1/projects/{projectId}/runtime/stop`
- `POST /v1/projects/{projectId}/runtime/restart`
- `POST /v1/projects/{projectId}/database/provision`
- `GET /v1/projects/{projectId}/database/status`
- `POST /v1/projects/{projectId}/database/stop`

## Запуск

Требуется:
- Go 1.23+
- Redis
- Docker

API:

```bash
go run ./cmd/api
```

Worker:

```bash
go run ./cmd/worker
```

## Быстрый flow

1. Авторизация GitHub: `GET /v1/auth/github/login`
2. Создать проект: `POST /v1/projects`
3. Подключить репозиторий:
   - `repo_full_name`
   - `branch`
   - `build_command`
   - `root_dir`
   - `output_dir`
4. Дождаться деплоя в `GET /v1/projects/{projectId}/deployments`
5. Поднять БД: `POST /v1/projects/{projectId}/database/provision`
6. Запустить runtime: `POST /v1/projects/{projectId}/runtime/start`
