package main

import (
	"flag"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func main() {
	isDebug := flag.Bool("debug", false, "Run the bot in the debug mode")
	flag.Parse()
	log.SetFormatter(&log.JSONFormatter{})
	log.SetReportCaller(true)
	entry := log.WithField("app", "tgbot")
	if err := tgbotapi.SetLogger(entry); err != nil {
		entry.Fatalf("main: tgbot failed to set the logger, %v", err)
	}
	token, err := misc.GetEnv("RAMPART_TGBOT_TOKEN")
	if err != nil {
		entry.Fatal(err)
	}
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		entry.Fatalf("main: tgbot failed to start, %v", err)
	}
	bot.Debug = *isDebug
	updatesChannel, _ := bot.GetUpdatesChan(tgbotapi.UpdateConfig{Timeout: 60})
	for update := range updatesChannel {
		chatID := update.Message.Chat.ID
		if _, err := bot.Send(tgbotapi.NewMessage(chatID, "Привіт, солоденький!")); err != nil {
			userID := 0
			if update.Message.From != nil {
				userID = update.Message.From.ID
			}
			entry.WithFields(
				log.Fields{"update_id": update.UpdateID, "chat_id": chatID, "user_id": userID},
			).Errorf("main: tgbot failed to send, %v", err)
		}
	}
}
