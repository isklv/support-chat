# Архитектура Support Chat

В этом документе подробно описана архитектура сервиса, поток передачи сообщений между веб-клиентом и Telegram, структура базы данных и механизмы фоновой синхронизации.

---

## 🏗 Общая схема системы

```mermaid
flowchart TD
    User["Посетитель сайта (Браузер)"]
    Widget["Клиентский виджет (widget.js / React)"]
    Backend["Go Backend (REST API / Relay)"]
    DB[("База данных (SQLite / PostgreSQL)")]
    TelegramAPI["Telegram Bot API"]
    Operators["Операторы поддержки (Telegram Supergroup Topics)"]
    SMTP["SMTP Сервер (Email уведомления)"]

    User <-->|Ввод сообщений, скриншоты| Widget
    Widget <-->|REST API + Polling| Backend
    Backend <-->|Хранение тредов и сообщений| DB
    Backend <-->|Relay Worker & Polling/Webhook| TelegramAPI
    TelegramAPI <-->|Топики, сообщения, фото, команды /close| Operators
    Backend -.->|Email при ответе оффлайн| SMTP -.-> User
```

---

## 🔄 Поток обработки сообщений (Sequence Diagrams)

### 1. Отправка первого сообщения посетителем

```mermaid
sequenceDiagram
    autonumber
    actor User as Посетитель
    participant W as Виджет
    participant API as Backend API
    participant DB as База данных
    participant TG as Telegram Bot API
    actor Op as Оператор (Telegram)

    User->>W: Вводит вопрос (+ email, имя, скриншот)
    W->>API: POST /api/v1/support/messages (Multipart Form)
    Note over API: Проверка Honeypot, rate-limit, валидация email и картинки
    API->>DB: Создание SupportThread, SupportMessage, SupportToken
    API-->>W: Возврат SupportView + token (сохраняется в localStorage)
    API->>API: Kick() фонового Relay Worker
    API->>TG: createForumTopic ("#1 Иван <ivan@example.com>")
    TG-->>API: message_thread_id (ID топика)
    API->>TG: sendMessage (Карточка обращения: страница, язык, устройство)
    API->>TG: sendPhoto / sendMessage (Текст и изображение)
    TG->>Op: Уведомление в теме операторов
```

---

### 2. Ответ оператора из Telegram

```mermaid
sequenceDiagram
    autonumber
    actor Op as Оператор (Telegram)
    participant TG as Telegram Bot API
    participant Poller as Poller / Webhook
    participant DB as База данных
    participant W as Виджет
    actor User as Посетитель

    Op->>TG: Отправляет сообщение в топик (текст / фото)
    alt Режим Long Polling
        Poller->>TG: getUpdates()
        TG-->>Poller: Update с новым сообщением
    else Режим Webhook
        TG->>Poller: POST /api/v1/support/telegram/webhook
    end
    Poller->>DB: Поиск SupportThread по telegram_topic_id
    Poller->>DB: Сохранение SupportMessage (author: "operator")
    opt Пользователь указал email и не открывал чат >15 мин
        Poller->>User: Отправка Email со ссылкой ?support=<token>
    end
    W->>Poller: GET /api/v1/support (Adaptive Polling: каждые 6 сек при открытом окне)
    Poller-->>W: Новые сообщения
    W->>User: Отображение баббла с ответом оператора
    W->>Poller: POST /api/v1/support/seen (сброс счетчика)
```

---

### 3. Закрытие и переоткрытие обращения

- Оператор отправляет команду `/close` прямо в топике Telegram:
  - Статус треда в БД меняется на `closed`.
  - В веб-чате появляется системное уведомление `Тема закрыта оператором`.
  - Топик в Telegram закрывается через `closeForumTopic`.
- Если посетитель пишет снова или оператор отправляет `/reopen`:
  - Статус треда меняется на `open`.
  - Топик в Telegram переоткрывается через `reopenForumTopic`.

---

## 🗄 Структура базы данных

### Таблица `support_threads`
| Колонка | Тип | Описание |
|---|---|---|
| `id` | VARCHAR(36) | Первичный ключ (UUID) |
| `number` | BIGINT | Порядковый номер тикета (#1, #2, ...) |
| `email` | VARCHAR(255) | Email автора обращения |
| `name` | VARCHAR(255) | Имя автора |
| `locale` | VARCHAR(10) | Язык (`ru`, `en`, `es`) |
| `status` | VARCHAR(20) | `open` или `closed` |
| `context` | TEXT | JSON с метаданными (page, user_agent, lang, custom) |
| `telegram_chat_id` | BIGINT | ID группы Telegram |
| `telegram_topic_id` | BIGINT | ID топика в группе |
| `user_seen_at` | TIMESTAMP | Время последнего прочтения пользователем |
| `emailed_at` | TIMESTAMP | Время последней отправки email-уведомления |
| `last_message_at` | TIMESTAMP | Время последнего сообщения в треде |
| `created_at` | TIMESTAMP | Дата создания |
| `closed_at` | TIMESTAMP | Дата закрытия |

### Таблица `support_messages`
| Колонка | Тип | Описание |
|---|---|---|
| `id` | VARCHAR(36) | Первичный ключ (UUID) |
| `thread_id` | VARCHAR(36) | Внешний ключ на `support_threads.id` |
| `author` | VARCHAR(20) | `user`, `operator`, `system` |
| `operator_name` | VARCHAR(100) | Отображаемое имя оператора |
| `text` | TEXT | Текст сообщения |
| `has_attachment` | BOOLEAN | Признак наличия скриншота/картинки |
| `telegram_message_id` | BIGINT | ID сообщения в Telegram |
| `relayed_at` | TIMESTAMP | Время успешной отправки в Telegram |
| `created_at` | TIMESTAMP | Дата создания |

### Таблица `support_attachments`
| Колонка | Тип | Описание |
|---|---|---|
| `message_id` | VARCHAR(36) | Внешний ключ на `support_messages.id` |
| `content_type` | VARCHAR(100) | MIME тип (`image/png`, `image/jpeg`, `image/webp`) |
| `data` | BLOB / BYTEA | Бинарные данные изображения |

### Таблица `support_tokens`
| Колонка | Тип | Описание |
|---|---|---|
| `token_hash` | VARCHAR(64) | SHA-256 хеш токена гостевой сессии |
| `thread_id` | VARCHAR(36) | Внешний ключ на тред |
| `created_at` | TIMESTAMP | Дата выдачи токена |
