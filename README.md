# EventHub — API Gateway для Event-платформы

**EventHub** — это **API Gateway**, построенный на **Gin** и **gRPC**, обеспечивающий безопасный и масштабируемый доступ к микросервисам:

- `event-service` — управление событиями
- `auth-service` — аутентификация и авторизация

---

## Особенности

- **gRPC → REST шлюз** — прозрачное проксирование REST-запросов в gRPC-микросервисы (`event`, `auth`)
- **JWT-аутентификация** — проверка `access_token` через `Authorization: Bearer <token>`
- **Rate Limiting** — глобальное ограничение (настраивается в `config`)
- **CORS включён** — `cors.Default()` для кросс-доменных запросов
- **Логирование с `request_id`** — каждый запрос получает уникальный UUID
- **Recovery от паник** — автоматический перехват и возврат `500` с логированием
- **Таймауты gRPC** — настраиваемый `timeout` из конфигурации
- **Чёткая обработка gRPC-ошибок** — `NotFound`, `InvalidArgument` → правильные HTTP-статусы
- **Graceful Shutdown** — в разработке

---

## API Endpoints

### Аутентификация (`/auth`)

| Метод | Путь | Описание |
|------|------|---------|
| `POST` | `/auth/register` | Регистрация пользователя |
| `POST` | `/auth/login` | Вход → `access_token`, `refresh_token` |
| `POST` | `/auth/refresh` | Обновление токенов |
| `POST` | `/auth/logout` | Выход |
| `POST` | `/auth/admin` | Проверка: админ ли? |

### События (`/events`)

> **Требуется `access_token` в заголовке `Authorization: Bearer <token>`**

| Метод | Путь | Описание |
|------|------|---------|
| `GET` | `/events` | Все события |
| `GET` | `/events/status/:status` | По статусу |
| `GET` | `/events/creator/:creator` | По создателю |
| `GET` | `/events/:id` | По ID |
| `POST` | `/events` | Создать событие |
| `PUT` | `/events` | Обновить (полное) (только владелец) |
| `DELETE` | `/events/:id` | Удалить (только владелец) |

---

## Шаги по запуску

1. **Клонируй репозиторий и перейдите в папку**:

   ```
   git clone https://github.com/Estriper0/AuthService.git
   cd AuthService
   ```

2. Настройте переменные окружения в `.env`:
   ```env
    APP_ENV=local

    DB_HOST=localhost
    DB_PORT=5432
    DB_NAME=db_name
    DB_USER=postgres
    DB_PASSWORD=12345

    REDIS_ADDR=redis:6379
    REDIS_PASSWORD=12345
   ```

3. **Запусти с помощью Docker Compose**:
   ```
   docker compose up --build -d
   ```
---