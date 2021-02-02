package main

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "snitch")
	if err := tgbotapi.SetLogger(entry); err != nil {
		entry.Fatalf("main: snitch failed to set the logger, %v", err)
	}
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", c.Snitch.DSN)
	if err != nil {
		entry.Fatalf("main: snitch failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: snitch failed to ping the db, %v", err)
	}
	metrics.RunServer(c.Snitch.Server, entry)
	if err = db.Close(); err != nil {
		entry.Fatalf("main: snitch failed to close the db, %v", err)
	}
}
