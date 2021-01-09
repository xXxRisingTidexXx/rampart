package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

type Handler interface {
	Name() string
	HandleUpdate(*tgbotapi.BotAPI, tgbotapi.Update) (bool, error)
}
