package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewAssistantHelpHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI) Handler {
	return &assistantHelpHandler{
		&helper{bot},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
		),
	}
}

type assistantHelpHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *assistantHelpHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "assistant-help"}, h.helper.sendMessage(
		update.Message.Chat.ID,
		"assistant_help",
		h.markup,
	)
}
