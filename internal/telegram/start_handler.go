package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewStartHandler(config config.Handler, bot *tgbotapi.BotAPI) Handler {
	return &startHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.AddButton),
				tgbotapi.NewKeyboardButton(config.ListButton),
				tgbotapi.NewKeyboardButton(config.HelpButton),
			),
		),
	}
}

type startHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *startHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "start"}, h.helper.sendMessage(
		update.Message.Chat.ID,
		"menu",
		h.markup,
	)
}
