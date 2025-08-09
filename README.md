# Marketplace (Go + PostgreSQL) — Prototype

## Что внутри
- Gin HTTP API
- GORM + PostgreSQL
- JWT для аутентификации
- Favorites, Offers, Services
- Chats: create chat endpoint + authenticated WebSocket; incoming messages are stored in DB

## Быстрый старт (docker)
1. Склонируй/распакуй проект
2. Запусти `docker-compose up -d`
3. Приложение будет на http://localhost:8080

Примеры:
- Регистрация: `POST /auth/register` { email, password, role: "customer"|"performer", name }
- Логин: `POST /auth/login` -> получишь token
- Создать оффер (customer): `POST /api/offers` (Authorization: Bearer <token>)
- Создать услугу (performer): `POST /api/services`
- Добавить в избранное (customer): `POST /api/favorites` { service_id }

## Замечания
- Секрет для JWT в `internal/auth/auth.go` нужно заменить на secure secret и вынести в env
- Для продакшн: миграции, валидация, логирование, тесты, rate-limit, CORS, HTTPS и т.д.