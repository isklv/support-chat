# 💬 Support Chat & Telegram Relay

Автономный сервис обратной связи и чата поддержки с Telegram-ботом. Позволяет встроить аккуратный виджет чата на любые сайты (лендинги, Tilda, WordPress, статические страницы, React/Next.js приложения), а отвечать клиентам прямо из Telegram-группы с темами (форум-топиками).

---

## ✨ Возможности

- **Интеграция с Telegram Forum Topics**:
  - Каждое новое обращение автоматически создает отдельную тему (топик) вида `#123 Иван <ivan@example.com>`.
  - При создании топика отправляется карточка обращения с контекстом: URL страницы, язык браузера, устройство/User-Agent, имя и email.
  - Операторы отвечают прямо в топике Telegram — сообщения мгновенно приходят посетителю в веб-чат.
  - Если группа без тем — бот отправляет сообщения в общий чат с поддержкой ответов через цитирование (Reply).
- **Команды оператора в Telegram**:
  - `/close` — закрыть обращение (помечает тикет закрытым в БД и закрывает топик в Telegram).
  - `/reopen` — открыть обращение заново.
- **Поддержка скриншотов и изображений**:
  - Посетитель может прикрепить файл или вставить скриншот из буфера обмена (`Ctrl+V` / `Cmd+V`).
  - Оператор может отправить картинку прямо в Telegram — она отобразится в чате у пользователя.
  - Увеличение картинок по клику (Lightbox).
- **Два режима работы Telegram**:
  - **Polling (по умолчанию)**: работает «из коробки» на локальном компьютере, на тестовом стенде или сервере без публичного домена и HTTPS!
  - **Webhook**: для продакшена за reverse proxy с HTTPS.
- **Два формата клиентского виджета**:
  - **Универсальный `<script>` (32 КБ)**: подключение одной строчкой на **любой** сайт (HTML, лендинг, Tilda, WordPress, PHP и т.д.) со встроенными стилями.
  - **React-компонент**: `import { SupportWidget } from 'support-chat-widget'` для современных фронтенд-проектов.
- **Хранилище**:
  - **SQLite** по умолчанию (zero-config, чистый Go без CGO, база данных в одном файле).
  - **PostgreSQL** при необходимости.
- **Email-уведомления (SMTP)**:
  - Если пользователь указал email и покинул страницу, при ответе оператора ему придет письмо с текстом ответа и прямой ссылкой на продолжение диалога (`?support=<token>`).
- **Защита от спама**:
  - Honeypot-поле от спам-ботов.
  - Встроенный IP rate limiter.

---

## 🚀 Быстрый старт через Docker

### 1. Клонируйте репозиторий и настройте `.env`
```bash
git clone https://github.com/your-username/support-chat.git
cd support-chat
cp .env.example .env
```

Заполните в `.env`:
```env
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrSTUvwxYZ
TELEGRAM_CHAT_ID=-1001234567890
```

### 2. Запустите контейнер
```bash
docker compose up -d
```

Сервер запустится на `http://localhost:8080`.
Проверить статус: `http://localhost:8080/healthz`.

Демо-страница доступна по адресу `http://localhost:8080/demo.html`.

---

## 🤖 Настройка Telegram бота

1. **Создайте бота**:
   - Напишите [@BotFather](https://t.me/BotFather) в Telegram и выполните команду `/newbot`.
   - Скопируйте полученный **Bot Token** в `TELEGRAM_BOT_TOKEN`.

2. **Создайте группу для операторов**:
   - Создайте новую группу в Telegram и добавьте в нее созданного бота.
   - В настройках группы включите **«Темы» (Topics / Forums)**:
     - *Настройки группы -> Управление группой -> Темы (Topics) -> Включить*.
   - Назначьте бота **Администратором** группы с правами управления темами и отправки сообщений.

3. **Узнайте Chat ID группы**:
   - Добавьте в группу бота [@userinfobot](https://t.me/userinfobot) или [@RawDataBot](https://t.me/RawDataBot).
   - Скопируйте ID группы (обычно начинается с `-100...`) в `TELEGRAM_CHAT_ID`.
   - Удалите вспомогательного бота из группы.

---

## 💻 Встраивание на сайт

### Вариант 1: Встраивание через `<script>` (Любые страницы, лендинги, Tilda, WordPress)

Добавьте одну строчку перед закрывающим тегом `</body>`:

```html
<script
  src="https://chat.yourdomain.com/widget.js"
  data-api="https://chat.yourdomain.com"
  data-title="Служба заботы"
  data-color="#2563eb"
  data-position="bottom-right"
  data-lang="ru">
</script>
```

#### Параметры data-атрибутов:

| Атрибут | Описание | По умолчанию |
|---|---|---|
| `data-api` | **Обязательно**. URL вашего бэкенда | — |
| `data-title` | Заголовок окна чата | `Поддержка` |
| `data-subtitle` | Подзаголовок | `Оператор ответит вам здесь` |
| `data-color` | Основной акцентный цвет (HEX или CSS цвет) | `#2563eb` |
| `data-position` | Положение кнопки: `bottom-right` или `bottom-left` | `bottom-right` |
| `data-lang` | Язык интерфейса: `ru`, `en`, `es`, `auto` | `auto` (язык браузера) |
| `data-auto-open` | Автоматически открыть окно при загрузке страницы | `false` |

---

### Вариант 2: Использование в React / Next.js / Vite

Установите пакет (или скопируйте папку `widget`):

```bash
npm install support-chat-widget
```

Используйте компонент в вашем приложении:

```tsx
import { SupportWidget } from 'support-chat-widget';

export function Layout() {
  return (
    <>
      {/* Ваш контент */}

      <SupportWidget
        apiUrl="https://chat.yourdomain.com"
        title="Служба поддержки"
        subtitle="Ответим в течение нескольких минут"
        primaryColor="#2563eb"
        position="bottom-right"
        lang="ru"
        // Необязательно: если пользователь авторизован в вашей системе:
        userId="user_12345"
        userEmail="alex@example.com"
        userName="Александр"
        metadata={{
          tariff: 'Business',
          registeredDays: 45,
        }}
      />
    </>
  );
}
```

---

### Вариант 3: Вызов через JavaScript API

Если скрипт уже загружен, можно инициализировать его программно:

```html
<script src="https://chat.yourdomain.com/widget.js"></script>
<script>
  window.SupportChat.init({
    apiUrl: 'https://chat.yourdomain.com',
    title: 'Поддержка',
    primaryColor: '#10b981',
    lang: 'ru',
    metadata: { source: 'landing_page_promo' }
  });
</script>
```

---

## 🛠 Конфигурация бэкенда (`.env`)

```env
# Сервер
PORT=8080
APP_URL=https://chat.yourdomain.com
ALLOWED_ORIGINS=*
STATIC_DIR=./static

# База данных
DB_DRIVER=sqlite
DB_DSN=./data/support.db
# Для Postgres:
# DB_DRIVER=postgres
# DB_DSN=postgres://user:password@localhost:5432/support_db?sslmode=disable

# Telegram
TELEGRAM_BOT_TOKEN=1234567890:ABCdefGHIjklMNOpqrSTUvwxYZ
TELEGRAM_CHAT_ID=-1001234567890
TELEGRAM_MODE=polling # или webhook
WEBHOOK_URL=https://chat.yourdomain.com/api/v1/support/telegram/webhook

# Email (SMTP) — опционально
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USER=support@example.com
SMTP_PASS=secretpassword
SMTP_FROM=support@example.com
```

---

## 📁 Структура проекта

```
support-chat/
├── backend/                  # Сервер на Go
│   ├── cmd/server/main.go    # Точка входа сервера
│   ├── internal/api/         # REST API, CORS, Rate Limiters
│   ├── internal/config/      # Загрузка конфигурации из env
│   ├── internal/db/          # Модели GORM (SQLite & PostgreSQL)
│   ├── internal/mailer/      # Отправка email через SMTP
│   └── internal/support/     # Логика диалогов, релей и Telegram клиент
├── widget/                   # Клиентский виджет (Preact / TypeScript / Vite)
│   ├── src/components/       # UI компоненты чата и иконки
│   ├── src/embed.tsx         # Точка входа для standalone <script>
│   ├── src/index.tsx         # Точка входа для React/npm библиотеки
│   ├── src/styles.css        # Изолированные CSS стили
│   └── src/i18n.ts           # Локализация (RU, EN, ES)
├── examples/                 # Примеры интеграции
│   ├── html-embed/           # Пример для HTML / лендингов
│   └── react-app/            # Пример для React
├── Dockerfile                # Мультистейдж сборка (виджет + сервер в один образ)
├── docker-compose.yml        # Готовый запуск в Docker
├── Makefile                  # Команды сборки и тестов
└── README.md
```

---

## 🧪 Локальная разработка

### Сборка и запуск бэкенда:
```bash
cd backend
go test -v ./internal/...
go run ./cmd/server
```

### Сборка виджета:
```bash
cd widget
pnpm install
pnpm dev       # Режим разработки с live-reload
pnpm build     # Сборка dist/widget.js и библиотеки
```

---

## 📄 Лицензия

MIT. Свободно используйте в любых коммерческих и личных проектах!
