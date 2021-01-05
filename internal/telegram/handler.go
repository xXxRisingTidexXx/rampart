package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	Name() string
	ShouldServe(tgbotapi.Update) bool
	ServeUpdate(*tgbotapi.BotAPI, tgbotapi.Update) (log.Fields, error)
}
