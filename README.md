# Browser-International-Calls-Platform

Запуск всего стека через Docker Compose:

```bash
docker-compose up -d
```

Приложение: http://localhost:1573 (фронтенд проксирует `/api` на бэкенд).

Полноценный двусторонний голосовой звонок (браузер ↔ телефон) и проверка на Twilio trial — см. [backend/docs/WEBRTC_SETUP.md](backend/docs/WEBRTC_SETUP.md) (разделы «Полноценный двусторонний разговор», «Проверка на Trial-аккаунте», «Запуск через Docker Compose»).