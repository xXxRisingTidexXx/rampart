package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewHelpHandler(bot *tgbotapi.BotAPI) Handler {
	return &helpHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Зрозуміло \U0001F44D")),
		),
	}
}

type helpHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *helpHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "help"}, h.helper.sendMessage(update, "help", h.markup)
}
