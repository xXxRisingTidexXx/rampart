package main

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "moderator")
	if err := tgbotapi.SetLogger(entry); err != nil {
		entry.Fatalf("main: moderator failed to set the logger, %v", err)
	}
	c, err := config.NewConfig()
	if err != nil {
		entry.Fatal(err)
	}
	db, err := sql.Open("postgres", c.Moderator.DSN)
	if err != nil {
		entry.Fatalf("main: moderator failed to open the db, %v", err)
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		entry.Fatalf("main: moderator failed to ping the db, %v", err)
	}
	_, err = tgbotapi.NewBotAPI(c.Moderator.Token)
	if err != nil {
		_ = db.Close()
		entry.Fatalf("main: moderator failed to create the bot, %v", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals
	if err = db.Close(); err != nil {
		entry.Fatalf("main: moderator failed to close the db, %v", err)
	}
}
