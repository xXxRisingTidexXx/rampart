package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorStartHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI) Handler {
	return &moderatorStartHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.MarkupButton),
				tgbotapi.NewKeyboardButton(config.HelpButton),
			),
		),
	}
}

type moderatorStartHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *moderatorStartHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "moderator-start"}, h.helper.sendMessage(
		update.Message.Chat.ID,
		"moderator_menu",
		h.markup,
	)
}
