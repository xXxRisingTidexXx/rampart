package main

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"time"
)

// TODO: think about signal handling & graceful shutdown.
// TODO: think about viber integration, https://github.com/mileusna/viber .
// TODO: configure bot commands at BotFather.
func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "telegram")
	if err := tgbotapi.SetLogger(entry); err != nil {
		entry.Fatalf("main: telegram failed to set the logger, %v", err)
	}
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", c.Telegram.DSN)
	if err != nil {
		entry.Fatalf("main: telegram failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: telegram failed to ping the db, %v", err)
	}
	bot, err := tgbotapi.NewBotAPI(c.Telegram.Token)
	if err != nil {
		_ = db.Close()
		entry.Fatalf("main: telegram failed to instantiate, %v", err)
	}
	updates, _ := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: c.Telegram.Timeout})
	time.Sleep(200 * time.Millisecond)
	updates.Clear()
	metrics.RunServer(c.Telegram.Server, entry)
	for update := range updates {
		if update.Message != nil && update.Message.Chat != nil && update.Message.IsCommand() {
			command, text := update.Message.Command(), "Hello from default!"
			switch command {
			case "start":
				text = "Hello from start!"
			case "help":
				text = "Hello from help!"
			default:
				command = "default"
			}
			metrics.TelegramCommands.WithLabelValues(command).Inc()
			_, err := bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, text))
			if err != nil {
				entry.Errorf("main: telegram failed to send a message, %v", err)
			}
		}
	}
	bot.StopReceivingUpdates()
	if err = db.Close(); err != nil {
		entry.Fatalf("main: telegram failed to close the db, %v", err)
	}
}
