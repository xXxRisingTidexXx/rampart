package main

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/telegram"
)

// TODO: think about signal handling & graceful shutdown.
// TODO: think about viber integration, https://github.com/mileusna/viber .
// TODO: configure bot commands at BotFather.
// TODO: think about max connections, https://www.alexedwards.net/blog/configuring-sqldb .
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
	dispatcher, err := telegram.NewDispatcher(c.Telegram.Dispatcher, db, entry)
	if err != nil {
		_ = db.Close()
		entry.Fatal(err)
	}
	metrics.RunServer(c.Telegram.Server, entry)
	dispatcher.Dispatch()
	if err = db.Close(); err != nil {
		entry.Fatalf("main: telegram failed to close the db, %v", err)
	}
}
