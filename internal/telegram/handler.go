package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler interface {
	HandleUpdate(update tgbotapi.Update) (Fields, error)
}
