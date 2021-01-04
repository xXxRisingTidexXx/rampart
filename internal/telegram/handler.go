package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	ShouldServe(tgbotapi.Update) bool
	ServeUpdate(tgbotapi.Update) (log.Fields, error)
}
