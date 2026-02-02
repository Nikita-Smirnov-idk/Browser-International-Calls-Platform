# Архитектура Backend

## Общее описание

Backend реализован на языке Go 1.22+ с использованием Clean Architecture для обеспечения разделения ответственности и тестируемости кода.

## Архитектурные слои

### 1. Domain Layer (internal/domain)

Содержит доменные модели и интерфейсы репозиториев. Не зависит от других слоев.

**Файлы:**
- `user.go` - модель пользователя
- `call.go` - модель звонка с константами статусов
- `repositories.go` - интерфейсы UserRepository и CallRepository

**Основные типы:**
```
User {
    ID, Email, PasswordHash, CreatedAt
}

Call {
    ID, UserID, PhoneNumber, StartTime, Duration, Status, CreatedAt
}
```

### 2. Use Cases Layer (internal/use_cases)

Реализует бизнес-логику приложения. Зависит только от domain layer.

**Модули:**
- `auth/` - регистрация, вход, выход
- `calls/` - создание и завершение звонков
- `history/` - получение истории звонков с фильтрацией и пагинацией

**Принципы:**
- Каждый use case имеет структуры Input и Output
- Валидация входных данных
- Логирование операций через slog
- Возврат доменных ошибок

### 3. Infrastructure Layer (internal/infrastructure)

Реализует интерфейсы, определенные в domain layer.

**Компоненты:**
- `postgres/` - реализация репозиториев через GORM
  - `connection.go` - подключение к БД с настройкой пула соединений
  - `migrations.go` - автоматическое применение SQL миграций
  - `user_repository.go` - CRUD операции для users
  - `call_repository.go` - CRUD операции для calls
- `jwt/` - генерация и валидация JWT токенов

**Параметры подключения к БД:**
- MaxOpenConns: 25
- MaxIdleConns: 5
- Logger: GORM с silent mode для миграций

### 4. Transport Layer (internal/transport/http)

HTTP API реализован через Gin framework.

**Структура:**
- `router.go` - регистрация маршрутов и middleware
- `handlers/` - HTTP handlers для endpoints
  - `auth_handler.go` - /api/auth/*
  - `calls_handler.go` - /api/calls (Create, Update)
  - `history_handler.go` - /api/calls/history
  - `health_handler.go` - /system/health
- `middleware/` - промежуточное ПО
  - `auth.go` - валидация JWT токенов
  - `cors.go` - настройка CORS
  - `recovery.go` - обработка паник

### 5. Application Layer (internal/app)

Инициализация и связывание компонентов приложения.

**Файл:** `app.go`

**Процесс инициализации:**
1. Создание подключения к БД
2. Применение миграций
3. Инициализация репозиториев
4. Создание JWT сервиса
5. Инициализация use cases
6. Создание handlers
7. Настройка роутера

### 6. Configuration Layer (internal/config)

Загрузка конфигурации из переменных окружения.

**Структуры:**
- `Config` - основная конфигурация
- `ServerConfig` - настройки сервера
- `DatabaseConfig` - параметры подключения к БД
- `JWTConfig` - секрет для токенов

## База данных

### Схема

**Таблица users:**
```sql
id UUID PRIMARY KEY
email VARCHAR(255) UNIQUE NOT NULL
password_hash VARCHAR(255) NOT NULL
created_at TIMESTAMP WITH TIME ZONE
```

**Таблица calls:**
```sql
id UUID PRIMARY KEY
user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
phone_number VARCHAR(50) NOT NULL
start_time TIMESTAMP WITH TIME ZONE NOT NULL
duration INTEGER DEFAULT 0
status VARCHAR(20) NOT NULL DEFAULT 'initiated'
created_at TIMESTAMP WITH TIME ZONE
```

### Индексы

- `idx_users_email` ON users(email)
- `idx_calls_user_id` ON calls(user_id)
- `idx_calls_start_time` ON calls(start_time)
- `idx_calls_user_start` ON calls(user_id, start_time DESC)

### Миграции

Миграции хранятся в директории `migrations/` как SQL файлы. Применяются автоматически при старте приложения.

Порядок применения: сортировка по имени файла (001_, 002_, ...).

## Аутентификация и авторизация

### JWT токены

**Алгоритм:** HS256

**Claims:**
- user_id (string)
- email (string)
- exp (expiration time)
- iat (issued at)

**Время жизни:** определяется через переменную окружения JWT_EXPIRES_IN (по умолчанию 60 минут)

**Передача:** Bearer токен в заголовке Authorization

### Хеширование паролей

**Алгоритм:** bcrypt

**Cost factor:** DefaultCost (10)

## API Endpoints

### Публичные
- POST /api/auth/register
- POST /api/auth/login

### Защищенные (требуют JWT)
- POST /api/auth/logout
- POST /api/calls
- PUT /api/calls/:id
- GET /api/calls/history

### Системные
- GET /system/health

## Обработка ошибок

### Уровни обработки

1. **Use case level:** валидация, бизнес-логика
2. **Repository level:** ошибки БД
3. **Handler level:** маппинг на HTTP статус коды

### HTTP статус коды

- 200 OK - успешная операция
- 201 Created - создан ресурс
- 204 No Content - успешно без тела ответа
- 400 Bad Request - ошибка валидации
- 401 Unauthorized - отсутствует или невалидный токен
- 404 Not Found - ресурс не найден
- 409 Conflict - конфликт (например, email уже существует)
- 500 Internal Server Error - внутренняя ошибка

## Логирование

**Библиотека:** log/slog (стандартная библиотека Go)

**Уровни:**
- Info - успешные операции
- Error - ошибки с контекстом

**Логируемые операции:**
- Регистрация/вход пользователя
- Создание/завершение звонка
- Ошибки БД и внутренние ошибки

## Зависимости

### Основные
- github.com/gin-gonic/gin - веб-фреймворк
- gorm.io/gorm - ORM
- gorm.io/driver/postgres - драйвер PostgreSQL
- github.com/golang-jwt/jwt/v5 - JWT токены
- golang.org/x/crypto/bcrypt - хеширование паролей

### Стандартная библиотека
- log/slog - логирование
- context - управление контекстом
- time - работа со временем
- net/http - HTTP сервер

## Принципы разработки

1. **Dependency Rule:** зависимости направлены внутрь (к domain)
2. **Separation of Concerns:** каждый слой решает свою задачу
3. **Interface Segregation:** интерфейсы определены в domain
4. **Single Responsibility:** один use case - одна задача
5. **Explicit Dependencies:** все зависимости передаются через конструкторы

## Ограничения текущей реализации

1. JWT токены stateless, logout не инвалидирует токен на сервере
2. Миграции применяются только вперед (без rollback)
3. Отсутствует rate limiting
4. Отсутствует кеширование
5. Пагинация реализована in-memory после получения всех записей

