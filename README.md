# Marketplace (Go + PostgreSQL) — Prototype

## Что внутри
- Gin HTTP API
- PGXpool + Squirrel + PostgreSQL
- JWT для аутентификации
- Favorites, Offers, Services

## Быстрый старт (docker)
1. Склонируй/распакуй проект
2. Запусти `docker-compose up -d`
3. Зайди в VSCode и нажми F5, автоматически запустится сервис
4. Для тестирования можешь использовать любой CURL клиент

Примеры:
- Регистрация: `POST /auth/register` { email, password, role: "customer"|"performer", name }
- Логин: `POST /auth/login` -> получишь token
- Создать оффер (customer): `POST /api/offers` (Authorization: Bearer <token>)
- Создать услугу (performer): `POST /api/services`
- Добавить в избранное (customer): `POST /api/favorites` { service_id }
