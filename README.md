# EventHub — API Gateway для Event-платформы

**EventHub** — это **API Gateway**, построенный на **Gin** и **gRPC**, обеспечивающий безопасный и масштабируемый доступ к микросервисам:

- [`event-service`](https://github.com/Estriper0/EventService) — управление событиями
- [`auth-service`](https://github.com/Estriper0/AuthService) — аутентификация и авторизация

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
- **Graceful Shutdown** — безопасное завершение работы приложения при его остановке.

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

| Метод   | Путь                         | Описание                                          |
|---------|------------------------------|---------------------------------------------------|
| `GET`   | `/events/`                   | Получить все события                              |
| `GET`   | `/events/status/:status`     | Получить события по статусу                       |
| `GET`   | `/events/creator/:creator`   | Получить события по создателю (UUID)              |
| `GET`   | `/events/:id/users`          | Получить всех пользователей, зарегистрированных на событие |
| `GET`   | `/events/:id`                | Получить событие по ID                            |
| `POST`  | `/events/`                   | Создать новое событие                             |
| `DELETE`| `/events/:id`                | Удалить событие                                   |
| `PUT`   | `/events/`                   | Обновить событие (полное)                         |
| `GET`   | `/events/me`                 | Получить все события, на которые зарегистрирован текущий пользователь |
| `POST`  | `/events/:id/register`       | Зарегистрироваться на событие                     |
| `DELETE`| `/events/:id/register`        | Отменить регистрацию на событие                   |

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