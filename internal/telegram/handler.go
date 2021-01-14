package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	HandleUpdate(update tgbotapi.Update) (log.Fields, error)
}
