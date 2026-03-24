# Novex Deploy Frontend

Control panel для платформы деплоя backend-сервисов и Telegram-ботов.

## Функциональность

- GitHub OAuth вход
- Импорт репозитория и создание проекта
- Запуск deploy и просмотр логов
- Runtime control (start/stop/restart)
- Telegram bot configuration
- Управление env variables

## Переменные окружения

Скопируйте `.env.example` в `.env` при необходимости.

## Запуск

```bash
npm install
npm run dev
```

Сборка:

```bash
npm run build
```

## Важное

В dev режиме `vite.config.ts` проксирует `/v1/*` на backend `http://localhost:8888`.
