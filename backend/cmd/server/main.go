package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"support-chat/internal/api"
	"support-chat/internal/config"
	"support-chat/internal/db"
	"support-chat/internal/mailer"
	"support-chat/internal/support"
)

func main() {
	log.Println("support-chat starting up...")
	cfg := config.Load()

	database, err := db.Init(cfg.DBDriver, cfg.DBDSN)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	mailClient := mailer.New(cfg)

	var tg *support.Telegram
	if cfg.TelegramBotToken != "" && cfg.TelegramChatID != 0 {
		tg = support.NewTelegram(cfg.TelegramBotToken, cfg.TelegramChatID, cfg.TelegramBaseURL, cfg.TelegramProxy)
		log.Printf("telegram: bot initialized for group %d", cfg.TelegramChatID)
	} else {
		log.Println("telegram: credentials not provided (bot disabled)")
	}

	service := support.NewService(database, cfg, mailClient, tg)
	srv := api.NewServer(cfg, service)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// Launch Telegram background workers if bot is configured
	if tg != nil {
		go service.RunRelay(ctx)

		if cfg.TelegramMode == "webhook" && cfg.WebhookURL != "" {
			log.Printf("telegram: registering webhook at %s", cfg.WebhookURL)
			if err := tg.SetWebhook(ctx, cfg.WebhookURL); err != nil {
				log.Printf("telegram: setWebhook failed: %v", err)
			}
		} else {
			// Default to long-polling mode
			log.Println("telegram: running in polling mode")
			go service.RunPoller(ctx)
		}
	}

	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      srv.Routes(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("support-chat server listening on http://localhost:%s", cfg.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("support-chat shutting down gracefully...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("http server shutdown error: %v", err)
	}
	log.Println("support-chat stopped")
}
