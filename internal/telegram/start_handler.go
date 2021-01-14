package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewStartHandler(bot *tgbotapi.BotAPI) Handler {
	return &startHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Підписка \U0001F49C"),
				tgbotapi.NewKeyboardButton("Довідка \U0001F64B"),
			),
		),
	}
}

type startHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *startHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "start"}, h.helper.sendMessage(update, "start", h.markup)
}
