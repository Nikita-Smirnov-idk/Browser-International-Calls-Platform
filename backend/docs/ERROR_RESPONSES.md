# Формат ответов с ошибками

## Общий формат

Согласно техническому заданию (roman_tz.txt, раздел 4.1.3), все ошибки API возвращаются в следующем формате:

```json
{
  "error": "error_type",
  "message": "Human readable error message"
}
```

Где:
- `error` - тип ошибки (machine-readable, используется для программной обработки)
- `message` - описание ошибки для пользователя (human-readable)

## Типы ошибок по категориям

### Аутентификация и авторизация

#### unauthorized
HTTP Status: 401
```json
{
  "error": "unauthorized",
  "message": "Authentication required"
}
```

#### invalid_credentials
HTTP Status: 401
```json
{
  "error": "invalid_credentials",
  "message": "Invalid email or password"
}
```

#### user_already_exists
HTTP Status: 409
```json
{
  "error": "user_already_exists",
  "message": "User with this email already exists"
}
```

#### token_generation_error
HTTP Status: 500
```json
{
  "error": "token_generation_error",
  "message": "failed to generate token"
}
```

### Валидация данных

#### validation_error
HTTP Status: 400
```json
{
  "error": "validation_error",
  "message": "Описание ошибки валидации"
}
```

Примеры сообщений:
- "email is required"
- "password is required"
- "phone_number is required"
- "call_id is required"
- "Key: 'Email' Error:Field validation for 'Email' failed on the 'email' tag"

### Операции со звонками

#### call_creation_error
HTTP Status: 500
```json
{
  "error": "call_creation_error",
  "message": "Описание ошибки создания звонка"
}
```

#### call_update_error
HTTP Status: 500
```json
{
  "error": "call_update_error",
  "message": "Описание ошибки обновления звонка"
}
```

#### call_not_found
HTTP Status: 404
```json
{
  "error": "call_not_found",
  "message": "call not found"
}
```

#### call_initiation_failed
HTTP Status: 400, 503
```json
{
  "error": "call_initiation_failed",
  "message": "Описание ошибки инициации звонка"
}
```

#### call_termination_failed
HTTP Status: 400, 403, 404, 500
```json
{
  "error": "call_termination_failed",
  "message": "Описание ошибки завершения звонка"
}
```

### История звонков

#### history_fetch_error
HTTP Status: 500
```json
{
  "error": "history_fetch_error",
  "message": "Описание ошибки получения истории"
}
```

### Регистрация

#### registration_error
HTTP Status: 500
```json
{
  "error": "registration_error",
  "message": "Описание ошибки регистрации"
}
```

## Соответствие HTTP статусам

| HTTP Status | Описание | Типичные error types |
|-------------|----------|---------------------|
| 400 | Bad Request | validation_error, call_initiation_failed |
| 401 | Unauthorized | unauthorized, invalid_credentials |
| 403 | Forbidden | unauthorized (для ресурсов) |
| 404 | Not Found | call_not_found |
| 409 | Conflict | user_already_exists |
| 500 | Internal Server Error | token_generation_error, call_creation_error, history_fetch_error, registration_error |
| 503 | Service Unavailable | call_initiation_failed (VoIP недоступен) |

## Примеры использования

### Пример 1: Неудачная регистрация (email занят)

Request:
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "existing@example.com",
  "password": "password123"
}
```

Response:
```http
HTTP/1.1 409 Conflict
Content-Type: application/json

{
  "error": "user_already_exists",
  "message": "user with email existing@example.com already exists"
}
```

### Пример 2: Неудачный вход

Request:
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "wrongpassword"
}
```

Response:
```http
HTTP/1.1 401 Unauthorized
Content-Type: application/json

{
  "error": "invalid_credentials",
  "message": "Invalid email or password"
}
```

### Пример 3: Валидация при инициации звонка

Request:
```http
POST /api/calls/initiate
Authorization: Bearer <token>
Content-Type: application/json

{
  "phone_number": ""
}
```

Response:
```http
HTTP/1.1 400 Bad Request
Content-Type: application/json

{
  "error": "validation_error",
  "message": "phone_number is required"
}
```

### Пример 4: Попытка завершить чужой звонок

Request:
```http
POST /api/calls/terminate
Authorization: Bearer <token>
Content-Type: application/json

{
  "call_id": "some-uuid"
}
```

Response:
```http
HTTP/1.1 403 Forbidden
Content-Type: application/json

{
  "error": "call_termination_failed",
  "message": "unauthorized"
}
```

### Пример 5: VoIP сервис недоступен

Request:
```http
POST /api/calls/initiate
Authorization: Bearer <token>
Content-Type: application/json

{
  "phone_number": "+491512345678"
}
```

Response:
```http
HTTP/1.1 503 Service Unavailable
Content-Type: application/json

{
  "error": "call_initiation_failed",
  "message": "failed to initiate call"
}
```

## Рекомендации для клиентов

1. Всегда проверяйте HTTP статус код для определения категории ошибки
2. Используйте поле `error` для программной обработки специфичных случаев
3. Отображайте пользователю поле `message` для понимания причины ошибки
4. Обрабатывайте все возможные статус коды для каждого endpoint

## Изменения относительно исходной реализации

До исправления ошибки возвращались только с полем `error`:
```json
{
  "error": "error message"
}
```

После исправления (согласно ТЗ) добавлено поле `message`:
```json
{
  "error": "error_type",
  "message": "error description"
}
```

Это изменение применено ко всем handlers:
- auth_handler.go
- calls_handler.go
- history_handler.go
- webrtc_handler.go

