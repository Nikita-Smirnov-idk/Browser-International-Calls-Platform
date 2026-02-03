# Backend — Browser International Calls Platform

REST API на Golang (Gin) по Clean Architecture для платформы международных звонков.

## Реализованная функциональность

### Аутентификация и авторизация
- Регистрация пользователей с валидацией email и пароля
- Вход с генерацией JWT токенов
- Выход из системы
- Хеширование паролей через bcrypt
- JWT middleware для защищенных endpoints
- Унифицированный формат ошибок с полями error и message

### База данных
- PostgreSQL с GORM ORM
- Автоматические миграции для таблиц users, calls и WebRTC полей
- Индексы для оптимизации запросов (email, user_id, start_time, session_id)
- Репозитории с параметризованными запросами

### История звонков
- Получение истории с фильтрацией по датам
- Пагинация результатов
- Создание и обновление записей о звонках
- Корректное форматирование времени с учетом timezone

### WebRTC интеграция
- Инициация звонков через WebRTC
- Завершение активных звонков
- Управление VoIP сессиями
- Поддержка Mock и Twilio провайдеров

### Безопасность
- Хеширование паролей через bcrypt (cost=10)
- Защита от SQL-инъекций через GORM
- Валидация данных на сервере
- JWT токены для аутентификации
- Логирование критических операций через slog
- Проверка прав доступа к ресурсам пользователя

## Структура проекта

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа с graceful shutdown
├── internal/
│   ├── app/                     # Инициализация приложения
│   ├── config/                  # Конфигурация из переменных окружения
│   ├── domain/                  # Доменные модели и интерфейсы
│   │   ├── user.go
│   │   ├── call.go
│   │   └── repositories.go
│   ├── use_cases/               # Бизнес-логика
│   │   ├── auth/                # Register, Login, Logout
│   │   ├── calls/               # StartCall, EndCall
│   │   └── history/             # ListHistory с фильтрацией
│   ├── infrastructure/          # Инфраструктурный слой
│   │   ├── jwt/                 # JWT сервис
│   │   └── postgres/            # Репозитории и подключение к БД
│   └── transport/               # HTTP транспорт
│       └── http/
│           ├── router.go
│           ├── handlers/        # Auth, Calls, History handlers
│           └── middleware/      # Auth, CORS, Recovery
├── migrations/                  # SQL миграции
│   ├── 001_create_users_table.sql
│   ├── 002_create_calls_table.sql
│   └── 003_add_webrtc_fields_to_calls.sql
├── go.mod
├── go.sum
├── SETUP.md                     # Инструкции по настройке
└── README.md
```

## Быстрый старт

Подробные инструкции см. в [SETUP.md](./SETUP.md)

### С Docker Compose

```bash
docker-compose up -d postgres
docker-compose up backend
```

### Локально

1. Установите переменные окружения:
```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=calls
POSTGRES_PASSWORD=calls
POSTGRES_DB=calls
JWT_SECRET=your-secret-key
SERVER_PORT=8080
```

2. Запустите сервер:
```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

Сервер доступен на `http://localhost:8080`.

## API Endpoints

### Аутентификация
- `POST /api/auth/register` — регистрация пользователя
- `POST /api/auth/login` — вход пользователя
- `POST /api/auth/logout` — выход из системы (требуется Bearer токен)

### Звонки
- `POST /api/calls` — создание записи о звонке (требуется Bearer токен)
- `PUT /api/calls/:id` — обновление статуса и длительности звонка (требуется Bearer токен)
- `GET /api/calls/history` — получение истории звонков с пагинацией и фильтрацией (требуется Bearer токен)

### WebRTC
- `POST /api/calls/initiate` — инициация звонка через WebRTC (требуется Bearer токен)
- `POST /api/calls/terminate` — завершение звонка через WebRTC (требуется Bearer токен)

### Система
- `GET /system/health` — проверка состояния сервиса

Документация API: `../api/openapi.yml`

### Формат ошибок

Все ошибки возвращаются в унифицированном формате согласно техническому заданию:

```json
{
  "error": "error_type",
  "message": "Human readable error message"
}
```

## Технологический стек

- **Go 1.22** — язык программирования
- **Gin** — веб-фреймворк
- **GORM** — ORM для работы с БД
- **PostgreSQL** — база данных
- **JWT** — аутентификация
- **bcrypt** — хеширование паролей
- **slog** — структурированное логирование

## Архитектурные решения

### Clean Architecture
Проект следует принципам Clean Architecture с четким разделением слоев:
- **Domain** — доменные модели и интерфейсы
- **Use Cases** — бизнес-логика
- **Infrastructure** — реализация репозиториев и сервисов
- **Transport** — HTTP handlers и middleware

### Безопасность
- Параметризованные запросы через GORM
- Хеширование паролей с bcrypt (cost=10)
- JWT токены с временем жизни 24 часа
- Валидация входных данных
- Структурированное логирование операций

### Производительность
- Индексы БД для оптимизации запросов
- Connection pooling (max 25 open connections, 5 idle connections)
- Graceful shutdown с таймаутом 10 секунд

## Дополнительная документация

- [ARCHITECTURE.md](./docs/ARCHITECTURE.md) — подробное описание архитектуры
- [WEBRTC_ARCHITECTURE.md](./docs/WEBRTC_ARCHITECTURE.md) — архитектура WebRTC интеграции
- [DEPLOYMENT.md](./docs/DEPLOYMENT.md) — инструкции по развертыванию
- [WEBRTC_SETUP.md](./docs/WEBRTC_SETUP.md) — настройка WebRTC функциональности
- [ERROR_RESPONSES.md](./docs/ERROR_RESPONSES.md) — формат ответов с ошибками
- [SETUP.md](./SETUP.md) — быстрый старт и настройка окружения
