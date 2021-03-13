package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorHelpHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI) Handler {
	return &moderatorHelpHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
		),
	}
}

type moderatorHelpHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *moderatorHelpHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	return NewInfo("moderator-help"), h.helper.sendMessageNoPreview(
		update.Message.Chat.ID,
		"moderator_help",
		h.markup,
	)
}
