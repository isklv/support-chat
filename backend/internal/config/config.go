package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port           string
	DBDriver       string
	DBDSN          string
	AppURL         string
	AllowedOrigins []string
	StaticDir      string

	TelegramBotToken string
	TelegramChatID   int64
	TelegramBaseURL  string
	TelegramProxy    string
	TelegramMode     string // "polling" or "webhook"
	WebhookURL       string

	SMTPHost string
	SMTPPort int
	SMTPUser string
	SMTPPass string
	SMTPFrom string
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func Load() *Config {
	port := getEnv("PORT", "8080")
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbDSN := getEnv("DB_DSN", "./data/support.db")
	appURL := strings.TrimRight(getEnv("APP_URL", "http://localhost:"+port), "/")

	originsRaw := getEnv("ALLOWED_ORIGINS", "*")
	var origins []string
	if originsRaw == "*" {
		origins = []string{"*"}
	} else {
		for _, o := range strings.Split(originsRaw, ",") {
			if s := strings.TrimSpace(o); s != "" {
				origins = append(origins, s)
			}
		}
	}

	chatID, _ := strconv.ParseInt(getEnv("TELEGRAM_CHAT_ID", "0"), 10, 64)
	botToken := getEnv("TELEGRAM_BOT_TOKEN", "")
	webhookURL := getEnv("WEBHOOK_URL", "")
	telegramMode := strings.ToLower(getEnv("TELEGRAM_MODE", ""))
	if telegramMode == "" {
		if webhookURL != "" {
			telegramMode = "webhook"
		} else {
			telegramMode = "polling"
		}
	}

	smtpPort, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))

	return &Config{
		Port:             port,
		DBDriver:         dbDriver,
		DBDSN:            dbDSN,
		AppURL:           appURL,
		AllowedOrigins:   origins,
		StaticDir:        getEnv("STATIC_DIR", "./static"),
		TelegramBotToken: botToken,
		TelegramChatID:   chatID,
		TelegramBaseURL:  getEnv("TELEGRAM_BASE_URL", "https://api.telegram.org"),
		TelegramProxy:    getEnv("TELEGRAM_PROXY", ""),
		TelegramMode:     telegramMode,
		WebhookURL:       webhookURL,
		SMTPHost:         getEnv("SMTP_HOST", ""),
		SMTPPort:         smtpPort,
		SMTPUser:         getEnv("SMTP_USER", ""),
		SMTPPass:         getEnv("SMTP_PASS", ""),
		SMTPFrom:         getEnv("SMTP_FROM", ""),
	}
}
