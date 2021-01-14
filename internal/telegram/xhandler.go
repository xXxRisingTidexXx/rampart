package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

// TODO: should handler "know" about bot? I.e., accept it in constructor?
type XHandler interface {
	Name() string
	// TODO: replace multiple return with context struct with log fields.
	XHandleUpdate(*tgbotapi.BotAPI, tgbotapi.Update) (bool, error)
}
