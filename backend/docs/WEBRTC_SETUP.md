# WebRTC Setup Guide

## Быстрый старт

### 1. Установка зависимостей

```bash
cd backend
go mod download
```

### 2. Настройка переменных окружения

Создайте файл `.env` в корне директории `backend/`:

```env
# Server
SERVER_PORT=8080

# Database
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=calls_platform
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-change-in-production

# VoIP - Mock режим (для разработки)
VOIP_PROVIDER=mock
```

### 3. Запуск PostgreSQL

С помощью Docker:

```bash
docker run -d \
  --name postgres-calls \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=calls_platform \
  -p 5432:5432 \
  postgres:16-alpine
```

Или используйте `docker-compose.yml` из корня проекта:

```bash
docker-compose up -d postgres
```

### 4. Применение миграций

Миграции применяются автоматически при запуске приложения.

Миграции для WebRTC:
- `003_add_webrtc_fields_to_calls.sql` - добавляет поля `session_id`, `sdp_offer`, `sdp_answer`

### 5. Запуск сервера

```bash
cd backend
go run cmd/server/main.go
```

Или с использованием скомпилированного бинарника:

```bash
go build -o main cmd/server/main.go
./main
```

## Режимы работы VoIP

### Mock режим (для разработки)

Не требует реальных VoIP credentials. Имитирует работу VoIP сервиса.

```env
VOIP_PROVIDER=mock
```

**Преимущества:**
- Не требует регистрации в Twilio
- Бесплатно
- Быстрая разработка
- Предсказуемое поведение

**Ограничения:**
- Нет реальных звонков
- Упрощенный SDP

### Twilio режим (для production)

Требует регистрации и настройки аккаунта в Twilio.

```env
VOIP_PROVIDER=twilio
VOIP_ACCOUNT_SID=ACxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
VOIP_AUTH_TOKEN=your_auth_token
VOIP_FROM_NUMBER=+1234567890
```

#### Получение Twilio credentials

1. Зарегистрируйтесь на https://www.twilio.com/
2. Перейдите в Console Dashboard
3. Скопируйте `Account SID` и `Auth Token`
4. Купите номер телефона в разделе "Phone Numbers"
5. Укажите купленный номер в `VOIP_FROM_NUMBER`

#### Полноценный двусторонний разговор (Voice SDK)

Чтобы в браузере и на телефоне слышать друг друга, нужны дополнительные настройки:

1. **API Key** — в Twilio Console: Account → API keys & tokens → Create API Key. Сохраните **SID** и **Secret** (Secret показывается один раз).
2. **TwiML App** — в Twilio Console: Develop → TwiML Apps → Create TwiML App. Укажите:
   - **Friendly Name**: например `Browser Calls`
   - **Voice Request URL**: публичный URL вашего бэкенда, например `https://ВАШ-ДОМЕН/api/voice/twiml` (для локальной разработки — см. раздел про ngrok ниже).
3. В `.env` (или в переменных окружения для docker-compose) задайте:

```env
VOIP_API_KEY_SID=SKxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
VOIP_API_KEY_SECRET=your_api_key_secret
VOIP_TWIML_APP_SID=APxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
VOICE_PUBLIC_BASE_URL=https://ваш-туннель.euw.devtunnels.ms
```

Переменная `VOICE_PUBLIC_BASE_URL` необязательна. Если задана, в TwiML для `<Dial>` добавляется `statusCallback`: Twilio будет вызывать ваш бэкенд при смене состояния дозвона (initiated, ringing, answered, completed). В логах появятся записи `voice status callback from Twilio` с полем `DialCallStatus` (completed, busy, no-answer, failed и т.д.) — это помогает понять, почему звонок не дошёл до телефона.

После этого при инициации звонка бэкенд вернёт `voice_token`, фронтенд использует Twilio Voice SDK: браузер подключается к Twilio, Twilio по вашему TwiML URL дозванивается до номера и соединяет аудио. Вы слышите абонента в браузере, абонент слышит вас.

#### Проверка на Trial-аккаунте

На бесплатном trial-аккаунте Twilio исходящие звонки разрешены только на **верифицированные** номера.

1. В Twilio Console откройте **Phone Numbers → Manage → Verified Caller IDs** (или **Verify** в боковом меню).
2. Добавьте номер, на который планируете звонить (например второй телефон или номер друга), и пройдите верификацию по SMS/звонку.
3. В приложении указывайте этот номер в формате E.164 (например `+79031234567`).

После этого при нажатии «Позвонить» Twilio реально позвонит на указанный номер. При настроенном Voice SDK (см. выше) вы и абонент будете слышать друг друга.

**Полноценный двусторонний разговор (голос из браузера к номеру и обратно)** требует настройки Voice SDK (API Key + TwiML App + публичный URL для TwiML). Без этого абонент услышит только демо TwiML, а голос из браузера не передаётся.

#### Запуск через Docker Compose и проверка с trial

1. В корне проекта создайте `.env` с переменными Twilio (в т.ч. `VOIP_API_KEY_SID`, `VOIP_API_KEY_SECRET`, `VOIP_TWIML_APP_SID` для двустороннего разговора).
2. Запустите стек:

```bash
docker-compose up -d
```

3. Twilio должен вызывать ваш бэкенд по адресу **Voice Request URL** TwiML App. При локальной разработке бэкенд недоступен из интернета, поэтому нужен туннель:
   - Установите [ngrok](https://ngrok.com/) и запустите туннель на порт, где доступно приложение (фронтенд проксирует `/api` на бэкенд). Например, если заходите на приложение по `http://localhost:1573`, запустите: `ngrok http 1573`.
   - Скопируйте HTTPS-URL ngrok (например `https://abc123.ngrok.io`) и в Twilio Console в TwiML App укажите **Voice Request URL**: `https://abc123.ngrok.io/api/voice/twiml`.
4. Откройте в браузере приложение (через ngrok-URL или `http://localhost:1573`), войдите, введите верифицированный номер и нажмите «Позвонить». Должен установиться полноценный голосовой звонок: вы слышите абонента в браузере, абонент слышит вас на телефоне.

## API Endpoints

### Инициация звонка

```http
POST /api/calls/initiate
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "phone_number": "+491512345678"
}
```

**Ответ (без Voice SDK):**

```json
{
  "call_id": "uuid",
  "session_id": "sess_123456789",
  "sdp_offer": "v=0\no=- 0 0 IN IP4 127.0.0.1\n...",
  "status": "connecting",
  "start_time": "2026-02-03T12:34:56Z"
}
```

**Ответ (с Voice SDK — при заданных VOIP_API_KEY_SID, VOIP_API_KEY_SECRET, VOIP_TWIML_APP_SID):**

```json
{
  "call_id": "uuid",
  "session_id": "voice_sdk",
  "sdp_offer": "",
  "status": "connecting",
  "start_time": "2026-02-03T12:34:56Z",
  "voice_token": "eyJ..."
}
```

### Завершение звонка

```http
POST /api/calls/terminate
Authorization: Bearer <JWT_TOKEN>
Content-Type: application/json

{
  "call_id": "uuid"
}
```

**Ответ:**

```json
{
  "call_id": "uuid",
  "duration": 42,
  "status": "completed"
}
```

## Тестирование WebRTC функциональности

### Unit тесты

```bash
cd backend
go test ./internal/use_cases/calls/... -v
```

### Интеграционные тесты

```bash
# Запустите PostgreSQL
docker-compose up -d postgres

# Запустите тесты
go test ./internal/... -v -tags=integration
```

### Ручное тестирование с curl

#### 1. Регистрация пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

#### 2. Вход

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

Сохраните `access_token` из ответа.

#### 3. Инициация звонка

```bash
export TOKEN="your-jwt-token"

curl -X POST http://localhost:8080/api/calls/initiate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone_number":"+491512345678"}'
```

#### 4. Завершение звонка

```bash
export CALL_ID="uuid-from-previous-response"

curl -X POST http://localhost:8080/api/calls/terminate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"call_id":"'$CALL_ID'"}'
```

#### 5. Просмотр истории

```bash
curl -X GET http://localhost:8080/api/calls/history \
  -H "Authorization: Bearer $TOKEN"
```

## Troubleshooting

### Проблема: VoIP сервис недоступен (503)

**Причина:** Неверные credentials Twilio или сервис недоступен.

**Решение:**
1. Проверьте `VOIP_ACCOUNT_SID` и `VOIP_AUTH_TOKEN`
2. Убедитесь, что у вас есть активный баланс в Twilio
3. Переключитесь на `VOIP_PROVIDER=mock` для разработки

### Проблема: Миграции не применяются

**Причина:** Директория `migrations/` не найдена.

**Решение:**
1. Убедитесь, что запускаете приложение из корня проекта
2. Или установите абсолютный путь к миграциям

### Проблема: Ошибка при инициации звонка

**Причина:** Некорректный формат номера телефона.

**Решение:**
Используйте международный формат: `+<country_code><number>`
- Правильно: `+491512345678`
- Неправильно: `491512345678`, `+49 151 2345678`

### Проблема: Unauthorized при завершении звонка

**Причина:** Попытка завершить чужой звонок.

**Решение:**
Убедитесь, что используете правильный JWT токен пользователя, который инициировал звонок.

### Проблема: Звонок инициируется, но на телефон не приходит (слышу только гудки «тун-дун»)

**Возможные причины и что проверить:**

1. **Voice Request URL недоступен для Twilio**  
   Twilio запрашивает TwiML по URL, указанному в TwiML App. Этот URL должен быть **публично доступен по HTTPS** (не `localhost`, не внутренний адрес Docker).  
   - При использовании devtunnels/ngrok укажите в TwiML App URL вида: `https://ВАШ-ТУННЕЛЬ/api/voice/twiml`.  
   - В логах бэкенда при звонке должны появляться строки: `twiml request from Twilio` и `twiml returned Dial`. Если их **нет** — Twilio не доходит до вашего бэкенда: проверьте URL в Twilio Console и что туннель проксирует `/api` на бэкенд.

2. **Trial-аккаунт Twilio: только верифицированные номера**  
   На бесплатном trial исходящие вызовы разрешены только на номера из **Verified Caller IDs**.  
   - Twilio Console → **Phone Numbers → Manage → Verified Caller IDs** (или **Verify**).  
   - Добавьте номер `+79653659793` (или тот, на который звоните) и пройдите верификацию (SMS/звонок).  
   - Звоните только на этот верифицированный номер в формате E.164 (`+79653659793`).

3. **География и права на исходящие**  
   В Twilio Console проверьте, что для вашего аккаунта разрешены исходящие звонки в нужную страну (например, Россия). Для trial могут быть ограничения по странам.

4. **Проверка в Twilio**  
   В Twilio Console → **Monitor → Logs → Voice** посмотрите попытки вызова: статус, ошибки, доходит ли запрос до вашего TwiML URL.

5. **Слышу «тун-дун» два раза, потом тишина**  
   Часто это значит: Twilio начал дозвон (вы слышите гудки в браузере), но исходящий звонок к телефону завершился с ошибкой (например, trial — номер не в Verified Caller IDs, или отказ оператора). Задайте `VOICE_PUBLIC_BASE_URL` (публичный URL бэкенда без слэша в конце) и перезапустите бэкенд. После звонка в логах появятся вызовы `voice status callback from Twilio` с полем **DialCallStatus**: `no-answer`, `busy`, `failed`, `completed` и т.д. По нему видно причину (например, `failed` при ограничениях trial).

## Логи и отладка

### Включение debug логов

```go
slog.SetLogLoggerLevel(slog.LevelDebug)
```

### Ключевые логи

- "call initiated with voice sdk" - звонок инициирован через Voice SDK (токен выдан, ожидается запрос TwiML от Twilio)
- "twiml request from Twilio" - Twilio запросил TwiML (значит Voice URL доступен)
- "twiml returned Dial" - бэкенд вернул TwiML с \<Dial\> на номер
- "twiml invalid or missing To" - Twilio вызвал URL без корректного параметра To
- "voice status callback from Twilio" - результат дозвона (DialCallStatus: completed, no-answer, busy, failed)
- "call initiated successfully" - звонок инициирован (без Voice SDK)
- "call terminated successfully" - звонок успешно завершен
- "failed to initiate voip call" - ошибка VoIP сервиса
- "unauthorized call termination attempt" - попытка завершить чужой звонок
- "session added" и "session removed" - операции с сессиями

## Production рекомендации

1. **Безопасность:**
   - Используйте сильный `JWT_SECRET` (минимум 32 символа)
   - Храните credentials в секретах (Kubernetes Secrets, AWS Secrets Manager)
   - Включите HTTPS

2. **Масштабирование:**
   - Используйте Redis для хранения сессий (вместо in-memory)
   - Настройте connection pooling для PostgreSQL
   - Используйте load balancer

3. **Мониторинг:**
   - Настройте алерты на ошибки VoIP сервиса
   - Мониторьте количество активных сессий
   - Отслеживайте среднюю длительность звонков

4. **Backup:**
   - Регулярные бэкапы PostgreSQL
   - Логирование в централизованную систему (ELK, Grafana Loki)

## Дополнительные ресурсы

- [Twilio Voice Documentation](https://www.twilio.com/docs/voice)
- [WebRTC Specification](https://webrtc.org/)
- [Backend Architecture](./ARCHITECTURE.md)
- [WebRTC Architecture](./WEBRTC_ARCHITECTURE.md)

