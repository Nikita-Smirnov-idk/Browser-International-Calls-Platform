# Backend — Browser International Calls Platform

REST API на Golang (Gin) по clean architecture.

## Структура

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа
├── internal/
│   ├── app/                     # Инициализация приложения и зависимостей
│   │   └── app.go
│   ├── domain/                  # Сущности и интерфейсы репозиториев
│   │   ├── user.go
│   │   ├── call.go
│   │   └── repositories.go
│   ├── use_cases/               # Бизнес-логика
│   │   ├── auth/
│   │   ├── calls/
│   │   └── history/
│   ├── infrastructure/          # Реализации (PostgreSQL и др.)
│   │   └── postgres/
│   └── transport/               # HTTP-слой
│       └── http/
│           ├── router.go
│           ├── handlers/
│           └── middleware/
├── go.mod
└── go.sum
```

## Запуск

```bash
cd backend
go mod tidy
go run ./cmd/server
```

Сервер слушает `http://localhost:8080`.

## API

Соответствует `api/openapi.yml`:

- `POST /auth/register` — регистрация
- `POST /auth/login` — вход
- `POST /auth/logout` — выход (Bearer)
- `POST /calls/start` — старт звонка (Bearer)
- `POST /calls/end` — завершение звонка (Bearer)
- `GET /calls/history` — история звонков (Bearer)
- `GET /system/health` — health check
