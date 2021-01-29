package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewHelpHandler(config config.Handler, bot *tgbotapi.BotAPI) Handler {
	return &helpHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
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
