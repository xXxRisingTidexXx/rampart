package main

import (
	"database/sql"
	"flag"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/metrics"
	"github.com/xXxRisingTidexXx/rampart/internal/telegram"
)

// TODO: think about signal handling & graceful shutdown.
func main() {
	isDebug := flag.Bool("debug", false, "Execute a single workflow instead of the whole schedule")
	flag.Parse()
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
	publisher, err := telegram.NewPublisher(c.Snitch.Publisher, db, entry)
	if err != nil {
		_ = db.Close()
		entry.Fatal(err)
	}
	if *isDebug {
		publisher.Run()
	} else {
		scheduler := cron.New()
		if _, err := scheduler.AddJob(c.Snitch.Spec, publisher); err != nil {
			_ = db.Close()
			entry.Fatalf("main: snitch failed to run: %v", err)
		}
		metrics.RunServer(c.Snitch.Server, entry)
		scheduler.Run()
	}
	if err = db.Close(); err != nil {
		entry.Fatalf("main: snitch failed to close the db, %v", err)
	}
}
