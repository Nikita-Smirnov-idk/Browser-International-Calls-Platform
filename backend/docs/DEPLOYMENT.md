# Инструкция по развертыванию Backend

## Требования

### Системные требования
- Go 1.22 или выше
- PostgreSQL 16 или выше
- Docker и Docker Compose (опционально)

### Минимальные ресурсы
- CPU: 1 core
- RAM: 512 MB
- Disk: 100 MB

## Переменные окружения

### Обязательные

```bash
POSTGRES_HOST=<хост базы данных>
POSTGRES_PORT=<порт базы данных>
POSTGRES_USER=<пользователь базы данных>
POSTGRES_PASSWORD=<пароль базы данных>
POSTGRES_DB=<имя базы данных>
JWT_SECRET=<секретный ключ для JWT>
```

### Опциональные

```bash
SERVER_PORT=8080                    # порт сервера (по умолчанию 8080)
POSTGRES_SSLMODE=disable            # режим SSL для PostgreSQL (по умолчанию disable)
VOIP_PROVIDER=mock                  # VoIP провайдер: mock или twilio (по умолчанию mock)
```

### Для WebRTC (при использовании Twilio)

```bash
VOIP_PROVIDER=twilio
VOIP_ACCOUNT_SID=<Twilio Account SID>
VOIP_AUTH_TOKEN=<Twilio Auth Token>
VOIP_FROM_NUMBER=<Номер телефона в формате +1234567890>
```

## Развертывание через Docker Compose

### Шаг 1: Клонирование репозитория

```bash
git clone <repository-url>
cd Browser-International-Calls-Platform
```

### Шаг 2: Настройка переменных окружения

Создать файл `.env` в корневой директории проекта:

```bash
# Database
POSTGRES_HOST=postgres
POSTGRES_PORT=5432
POSTGRES_USER=calls
POSTGRES_PASSWORD=calls
POSTGRES_DB=calls
POSTGRES_SSLMODE=disable

# JWT
JWT_SECRET=<сгенерировать случайную строку>

# Server
SERVER_PORT=8080
```

### Шаг 3: Запуск сервисов

```bash
docker-compose up -d postgres
docker-compose up --build backend
```

### Шаг 4: Проверка работоспособности

```bash
curl http://localhost:8080/system/health
```

Ожидаемый ответ: `{"status":"ok"}`

## Развертывание без Docker

### Шаг 1: Установка зависимостей

```bash
cd backend
go mod download
```

### Шаг 2: Настройка базы данных

Создать базу данных PostgreSQL:

```sql
CREATE DATABASE calls;
CREATE USER calls WITH PASSWORD 'calls';
GRANT ALL PRIVILEGES ON DATABASE calls TO calls;
```

### Шаг 3: Настройка переменных окружения

Создать файл `.env` в директории `backend/`:

```bash
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=calls
POSTGRES_PASSWORD=calls
POSTGRES_DB=calls
POSTGRES_SSLMODE=disable
JWT_SECRET=<сгенерировать случайную строку>
SERVER_PORT=8080
```

Загрузить переменные:

```bash
export $(cat .env | xargs)
```

### Шаг 4: Запуск приложения

```bash
go run cmd/server/main.go
```

Миграции применяются автоматически при старте.

### Шаг 5: Проверка работоспособности

```bash
curl http://localhost:8080/system/health
```

## Production развертывание

### Сборка бинарного файла

```bash
cd backend
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./cmd/server
```

### Настройка systemd (Linux)

Создать файл `/etc/systemd/system/calls-backend.service`:

```ini
[Unit]
Description=Browser International Calls Platform Backend
After=network.target postgresql.service

[Service]
Type=simple
User=calls
WorkingDirectory=/opt/calls-backend
ExecStart=/opt/calls-backend/server
Restart=on-failure
RestartSec=5s

Environment="POSTGRES_HOST=localhost"
Environment="POSTGRES_PORT=5432"
Environment="POSTGRES_USER=calls"
Environment="POSTGRES_PASSWORD=<password>"
Environment="POSTGRES_DB=calls"
Environment="JWT_SECRET=<secret>"
Environment="SERVER_PORT=8080"
Environment="VOIP_PROVIDER=mock"

[Install]
WantedBy=multi-user.target
```

Запуск сервиса:

```bash
sudo systemctl daemon-reload
sudo systemctl enable calls-backend
sudo systemctl start calls-backend
sudo systemctl status calls-backend
```

### Настройка Nginx (reverse proxy)

Создать файл `/etc/nginx/sites-available/calls-backend`:

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Активация конфигурации:

```bash
sudo ln -s /etc/nginx/sites-available/calls-backend /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### SSL сертификат (Let's Encrypt)

```bash
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d api.example.com
```

## Миграции базы данных

### Автоматическое применение

Миграции применяются автоматически при старте приложения. Файлы миграций находятся в `backend/migrations/`.

### Ручное применение

Для ручного применения миграций использовать SQL клиент:

```bash
psql -h localhost -U calls -d calls -f migrations/001_create_users_table.sql
psql -h localhost -U calls -d calls -f migrations/002_create_calls_table.sql
psql -h localhost -U calls -d calls -f migrations/003_add_webrtc_fields_to_calls.sql
```

## Мониторинг и логирование

### Логи приложения

Логи выводятся в stdout. При использовании Docker Compose:

```bash
docker-compose logs -f backend
```

При использовании systemd:

```bash
sudo journalctl -u calls-backend -f
```

### Health check endpoint

```bash
curl http://localhost:8080/system/health
```

### Мониторинг PostgreSQL

Проверка подключения:

```bash
psql -h localhost -U calls -d calls -c "SELECT version();"
```

Проверка количества записей:

```bash
psql -h localhost -U calls -d calls -c "SELECT COUNT(*) FROM users;"
psql -h localhost -U calls -d calls -c "SELECT COUNT(*) FROM calls;"
```

## Резервное копирование

### Backup базы данных

```bash
pg_dump -h localhost -U calls calls > backup_$(date +%Y%m%d_%H%M%S).sql
```

### Восстановление из backup

```bash
psql -h localhost -U calls calls < backup_20260203_120000.sql
```

## Безопасность

### Рекомендации для production

1. Использовать сильный JWT_SECRET (минимум 32 символа)
2. Включить SSL для PostgreSQL (POSTGRES_SSLMODE=require)
3. Использовать HTTPS для API (настроить SSL в Nginx)
4. Ограничить доступ к PostgreSQL по IP (pg_hba.conf)
5. Использовать firewall для ограничения портов
6. Регулярно обновлять зависимости: `go get -u && go mod tidy`
7. Настроитьротацию логов
8. Использовать secrets management (vault, AWS Secrets Manager)

### Генерация JWT_SECRET

```bash
openssl rand -base64 32
```

## Troubleshooting

### Ошибка подключения к БД

Проверить доступность PostgreSQL:

```bash
psql -h $POSTGRES_HOST -p $POSTGRES_PORT -U $POSTGRES_USER -d $POSTGRES_DB
```

### Ошибка применения миграций

Проверить лог приложения. Миграции должны применяться в порядке 001, 002.

Если миграция уже применена, она будет пропущена (CREATE TABLE IF NOT EXISTS).

### Порт уже занят

Проверить процессы на порту 8080:

```bash
lsof -i :8080
netstat -tulpn | grep 8080
```

### Проблемы с JWT токенами

Проверить что JWT_SECRET одинаковый при каждом запуске приложения.

## Масштабирование

### Горизонтальное масштабирование

Приложение stateless и может быть запущено в нескольких экземплярах за load balancer.

Требования:
- Единая база данных PostgreSQL
- Единый JWT_SECRET для всех инстансов

### Вертикальное масштабирование

Увеличить параметры подключения к БД в `internal/infrastructure/postgres/connection.go`:

```go
sqlDB.SetMaxOpenConns(50)  // увеличить с 25
sqlDB.SetMaxIdleConns(10)  // увеличить с 5
```

## Обновление приложения

### Zero-downtime deployment

1. Собрать новый бинарный файл
2. Запустить новый инстанс на другом порту
3. Переключить Nginx на новый порт
4. Дождаться завершения запросов на старом инстансе (graceful shutdown)
5. Остановить старый инстанс

### Rolling update с Docker Compose

```bash
docker-compose build backend
docker-compose up -d --no-deps backend
```

