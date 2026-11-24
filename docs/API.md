# REST API Справочник

Базовый путь API: `/api/v1`

---

## Аутентификация и заголовки

Все запросы от имени посетителя сайта идентифицируются:
1. Заголовком `X-Support-Token`: токен сессии, выданный при первом сообщении гостя (хранится в `localStorage`).
2. Опционально заголовком `X-User-Id`: идентификатор зарегистрированного пользователя в вашей системе.

Для загрузки изображений в теге `<img>` токен можно передавать в query-параметре: `?t=<token>`.

---

## Эндпоинты посетителя

### 1. Получение текущего диалога
`GET /api/v1/support`

**Заголовки:**
- `X-Support-Token`: `<token>` *(опционально)*
- `X-User-Id`: `<userId>` *(опционально)*

**Ответ (200 OK):**
```json
{
  "thread": {
    "id": "7f8b9c10-1234-4567-890a-bcdef0123456",
    "number": 42,
    "email": "user@example.com",
    "name": "Иван",
    "locale": "ru",
    "status": "open",
    "last_message_at": "2026-10-06T11:20:00Z",
    "created_at": "2026-10-06T11:00:00Z"
  },
  "messages": [
    {
      "id": "message-uuid-1",
      "thread_id": "7f8b9c10-1234-4567-890a-bcdef0123456",
      "author": "user",
      "text": "Здравствуйте, как оформить заказ?",
      "has_attachment": false,
      "created_at": "2026-10-06T11:00:00Z"
    },
    {
      "id": "message-uuid-2",
      "thread_id": "7f8b9c10-1234-4567-890a-bcdef0123456",
      "author": "operator",
      "operator_name": "Елена",
      "text": "Добрый день! Нажмите кнопку в корзине...",
      "has_attachment": false,
      "created_at": "2026-10-06T11:05:00Z"
    }
  ]
}
```

---

### 2. Отправка нового сообщения
`POST /api/v1/support/messages`

**Формат:** `multipart/form-data`

**Поля формы:**
- `text` *(string)*: текст сообщения (до 4000 знаков).
- `image` *(file, опционально)*: файл скриншота (PNG, JPEG, WebP до 4 МБ).
- `email` *(string, обязательно для первого сообщения гостя)*: email автора.
- `name` *(string, опционально)*: имя автора (до 100 знаков).
- `page` *(string, опционально)*: URL страницы, на которой находится посетитель.
- `lang` *(string, опционально)*: язык браузера (например, `ru`).
- `metadata` *(JSON string, опционально)*: произвольные данные пользователя.
- `website` *(string, honeypot)*: поле должно оставаться пустым. Если бот заполняет его, сервер игнорирует запрос.

**Ответ (200 OK):**
```json
{
  "thread": { ... },
  "messages": [ ... ],
  "token": "4a7f9b01c3e..." // Возвращается при создании нового треда гостя
}
```

---

### 3. Отметка о прочтении
`POST /api/v1/support/seen`

Помечает сообщения в теме как прочитанные пользователем, сбрасывая счетчик непрочитанных.

**Ответ (200 OK):**
```json
{ "ok": true }
```

---

### 4. Количество непрочитанных сообщений
`GET /api/v1/support/unread`

**Ответ (200 OK):**
```json
{ "unread": 2 }
```

---

### 5. Получение вложенного скриншота
`GET /api/v1/support/attachments/{messageId}?t={token}`

Возвращает бинарные данные изображения с заголовками `Content-Type: image/png` (или jpeg/webp) и кэшированием.

---

## Системные эндпоинты

### Проверка работоспособности (Healthcheck)
`GET /healthz`

**Ответ (200 OK):**
```json
{
  "status": "ok",
  "telegram": true
}
```

---

## Административные эндпоинты управления ботами

Веб-интерфейс панели управления доступен по адресу: `/admin` (или `/admin/`).

### 1. Список ботов
`GET /api/v1/admin/bots`

**Ответ (200 OK):**
```json
{
  "bots": [
    {
      "id": "c1f7b4e2-...",
      "name": "Основной бот поддержки",
      "token": "1234567890:ABC...",
      "username": "support_bot",
      "first_name": "Support Bot",
      "chat_id": -1001234567890,
      "mode": "polling",
      "webhook_url": "",
      "base_url": "https://api.telegram.org",
      "proxy": "",
      "is_default": true,
      "is_active": true,
      "description": "Основной бот для сайта",
      "created_at": "2026-11-23T19:00:00Z",
      "updated_at": "2026-11-23T19:00:00Z"
    }
  ]
}
```

---

### 2. Добавление нового бота
`POST /api/v1/admin/bots`

**Тело запроса (JSON):**
```json
{
  "name": "Бот поддержки EN",
  "token": "7123456789:ABCdefGHI...",
  "chat_id": -1009876543210,
  "mode": "polling",
  "webhook_url": "",
  "base_url": "https://api.telegram.org",
  "proxy": "",
  "is_default": false,
  "is_active": true,
  "description": "Бот для англоязычных пользователей"
}
```

**Ответ (201 Created):**
```json
{
  "bot": { ... }
}
```

---

### 3. Получение информации о боте
`GET /api/v1/admin/bots/{id}`

---

### 4. Редактирование бота
`PUT /api/v1/admin/bots/{id}`

**Тело запроса (JSON):**
Любые изменяемые поля: `name`, `token`, `chat_id`, `mode`, `webhook_url`, `base_url`, `proxy`, `is_active`, `is_default`, `description`.

---

### 5. Удаление бота
`DELETE /api/v1/admin/bots/{id}`

**Ответ (200 OK):**
```json
{ "ok": true }
```

---

### 6. Назначение бота основным (Default)
`POST /api/v1/admin/bots/{id}/default`

---

### 7. Включение / выключение бота
`POST /api/v1/admin/bots/{id}/toggle`

---

### 8. Проверка токена без сохранения
`POST /api/v1/admin/bots/test`

**Тело запроса (JSON):**
```json
{
  "token": "1234567890:ABCdef...",
  "base_url": "https://api.telegram.org",
  "proxy": ""
}
```

**Ответ (200 OK):**
```json
{
  "ok": true,
  "user": {
    "id": 1234567890,
    "is_bot": true,
    "first_name": "Support Bot",
    "username": "support_bot",
    "can_join_groups": true,
    "can_read_all_group_messages": false,
    "supports_inline_queries": false
  }
}
```

---

### 9. Проверка соединения сохраненного бота
`POST /api/v1/admin/bots/{id}/test`
