package main

import (
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "tgbot")
	_, err := tgbotapi.NewBotAPI(os.Getenv("RAMPART_TGBOT_TOKEN"))
	if err != nil {
		entry.Fatalf("main: tgbot failed to start, %v", err)
	}
}
