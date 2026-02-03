# WebRTC Architecture Documentation

## Обзор

Данный документ описывает архитектуру интеграции WebRTC в backend сервис платформы для международных звонков.

## Зона ответственности

Реализовано в соответствии с техническим заданием **Г. Д. Воронина** (gleb_tz.txt).

## Компоненты системы

### 1. Domain Layer

#### `domain/voip.go`
Содержит основные интерфейсы и структуры для работы с VoIP:

- **VoIPService** - интерфейс для работы с внешним VoIP провайдером
- **CallSession** - структура сессии звонка с WebRTC данными
- **SessionStatus** - статусы сессии (initialized, connecting, active, completed, failed)
- **WebRTCConfig** - конфигурация ICE серверов для WebRTC

#### `domain/call.go`
Расширена модель Call новыми полями:
- `SessionID` - идентификатор VoIP сессии
- `SDPOffer` - SDP offer для установки WebRTC соединения
- `SDPAnswer` - SDP answer от клиента
- Новые статусы: `connecting`, `active`

### 2. Infrastructure Layer

#### `infrastructure/voip/`

**client.go** - Фабрика для создания VoIP клиентов:
- Поддержка нескольких провайдеров (Twilio, Mock)
- Общие ошибки VoIP сервиса

**twilio_client.go** - Реализация для Twilio:
- Инициализация звонков через Twilio API
- Генерация SDP offer
- Управление сессиями

**mock_client.go** - Mock реализация для разработки и тестирования:
- Имитация VoIP функциональности
- Не требует реальных credentials

**session_manager.go** - Управление активными сессиями:
- Хранение сессий в памяти
- Автоматическая очистка истекших сессий (каждую минуту)
- Thread-safe операции

### 3. Use Cases Layer

#### `use_cases/calls/initiate.go`
Инициализация звонка через WebRTC:

**Входные данные:**
- `UserID` - идентификатор пользователя
- `PhoneNumber` - номер телефона для звонка

**Процесс:**
1. Валидация входных данных
2. Вызов VoIP сервиса для создания сессии
3. Получение SDP offer
4. Создание записи звонка в БД со статусом "connecting"
5. Возврат данных для установки WebRTC соединения

**Выходные данные:**
- `CallID` - идентификатор звонка в БД
- `SessionID` - идентификатор VoIP сессии
- `SDPOffer` - SDP offer для WebRTC
- `Status` - текущий статус
- `StartTime` - время начала

#### `use_cases/calls/terminate.go`
Завершение звонка через WebRTC:

**Входные данные:**
- `UserID` - идентификатор пользователя
- `CallID` - идентификатор звонка

**Процесс:**
1. Валидация входных данных
2. Получение звонка из БД
3. Проверка прав доступа (UserID)
4. Завершение VoIP сессии
5. Расчет длительности
6. Обновление записи в БД со статусом "completed"

**Выходные данные:**
- `CallID` - идентификатор звонка
- `Duration` - длительность в секундах
- `Status` - финальный статус

### 4. Transport Layer

#### `handlers/webrtc_handler.go`

**POST /api/calls/initiate** - Инициация звонка:
- Аутентификация через JWT middleware
- Валидация `phone_number`
- Вызов InitiateCallUseCase
- Возврат данных для WebRTC

**POST /api/calls/terminate** - Завершение звонка:
- Аутентификация через JWT middleware
- Валидация `call_id`
- Вызов TerminateCallUseCase
- Возврат длительности и статуса

## Диаграмма последовательности

### Инициация звонка

```
Frontend -> API: POST /api/calls/initiate {phone_number}
API -> Auth Middleware: Verify JWT
Auth Middleware -> InitiateUC: Execute(userID, phoneNumber)
InitiateUC -> VoIP Client: InitiateCall(phoneNumber)
VoIP Client -> Twilio API: Create Call
Twilio API -> VoIP Client: Session + SDP
VoIP Client -> InitiateUC: CallSession
InitiateUC -> CallRepo: Create(call)
CallRepo -> Database: INSERT call
InitiateUC -> API: {call_id, session_id, sdp_offer}
API -> Frontend: 200 OK + WebRTC data
```

### Завершение звонка

```
Frontend -> API: POST /api/calls/terminate {call_id}
API -> Auth Middleware: Verify JWT
Auth Middleware -> TerminateUC: Execute(userID, callID)
TerminateUC -> CallRepo: GetByID(callID)
CallRepo -> Database: SELECT call
TerminateUC -> VoIP Client: TerminateCall(sessionID)
VoIP Client -> Session Manager: Remove session
TerminateUC -> CallRepo: Update(call)
CallRepo -> Database: UPDATE call
TerminateUC -> API: {call_id, duration, status}
API -> Frontend: 200 OK
```

## Конфигурация

### Переменные окружения

```env
# VoIP Provider Configuration
VOIP_PROVIDER=twilio              # или "mock" для разработки
VOIP_ACCOUNT_SID=your-account-sid
VOIP_AUTH_TOKEN=your-auth-token
VOIP_API_KEY=your-api-key         # опционально
VOIP_FROM_NUMBER=+1234567890      # номер от которого идут звонки
```

### Mock режим для разработки

Для локальной разработки без реального VoIP провайдера:

```env
VOIP_PROVIDER=mock
```

Mock клиент имитирует работу VoIP сервиса без реальных звонков.

## Миграции БД

### 003_add_webrtc_fields_to_calls.sql

Добавляет поля для WebRTC в таблицу `calls`:
- `session_id VARCHAR(255)` - идентификатор VoIP сессии
- `sdp_offer TEXT` - SDP offer для WebRTC
- `sdp_answer TEXT` - SDP answer от клиента
- Индекс на `session_id` для быстрого поиска

## Обработка ошибок

### Формат ошибок

Все ошибки возвращаются в унифицированном формате согласно техническому заданию:

```json
{
  "error": "error_type",
  "message": "Human readable error description"
}
```

### Специфичные ошибки VoIP

- ErrVoIPServiceUnavailable - VoIP сервис недоступен (HTTP 503)
  - error: "call_initiation_failed" или "call_termination_failed"
  - message: описание ошибки VoIP сервиса

- ErrInvalidPhoneNumber - некорректный номер телефона (HTTP 400)
  - error: "validation_error"
  - message: "phone_number is required" или детали валидации

- ErrSessionNotFound - сессия не найдена (HTTP 404)
  - error: "call_not_found"
  - message: "call not found"

- ErrCallAlreadyActive - звонок уже активен (HTTP 409)
  - error: "call_already_active"
  - message: описание конфликта

- ErrUnauthorized - нет прав доступа (HTTP 403)
  - error: "unauthorized"
  - message: "Authentication required" или "No access to this call"

### Стратегия обработки

1. **Инициация звонка**: при ошибке VoIP сервиса возвращается HTTP 503 с соответствующим сообщением
2. **Завершение звонка**: ошибка VoIP сервиса логируется, но не прерывает процесс (звонок завершается в БД)

## Тестирование

### Unit тесты

Созданы тесты для use cases:
- `initiate_test.go` - тесты InitiateCallUseCase
- `terminate_test.go` - тесты TerminateCallUseCase

Покрываются сценарии:
- Успешная инициация/завершение
- Валидация входных данных
- Ошибки VoIP сервиса
- Ошибки БД
- Проверка прав доступа
- Звонок не найден

### Запуск тестов

```bash
cd backend
go test ./internal/use_cases/calls/... -v
```

## Ограничения и будущие улучшения

### Текущие ограничения

1. SDP exchange пока упрощен (mock SDP offer)
2. ICE candidates не обрабатываются
3. Нет реальной передачи аудио через Twilio
4. Сессии хранятся только в памяти (потеряются при рестарте)

### Планируемые улучшения

1. Полная интеграция с Twilio Programmable Voice
2. WebSocket для реал-тайм обмена ICE candidates
3. Хранение активных сессий в Redis
4. Мониторинг качества звонков
5. Ретраи при ошибках VoIP сервиса

## Взаимодействие с другими компонентами

### Интеграция с работой Р. А. Николаева

- Использует `CallRepository` для сохранения звонков
- Использует существующую схему БД (таблица `calls`)
- Расширяет схему миграцией 003

### Интеграция с работой Н. Д. Смирнова

- WebRTC handler интегрирован в общий роутер
- Использует общий Auth middleware
- Endpoints добавлены в группу `/api/calls`

## Безопасность

1. **Аутентификация**: все WebRTC endpoints защищены JWT middleware
2. **Авторизация**: проверка UserID при завершении звонка
3. **Валидация**: входные данные валидируются на уровне handler и use case
4. **Credentials**: VoIP credentials хранятся в переменных окружения
5. **Логирование**: критические операции логируются с контекстом

## Мониторинг и логирование

### Ключевые события

- Инициация звонка (уровень INFO)
- Завершение звонка (уровень INFO)
- Ошибки VoIP сервиса (уровень ERROR)
- Неавторизованные попытки (уровень WARN)
- Операции с сессиями (уровень DEBUG)

### Метрики для мониторинга

- Количество инициированных звонков
- Количество завершенных звонков
- Средняя длительность звонков
- Процент ошибок VoIP сервиса
- Количество активных сессий

