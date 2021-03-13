package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewAssistantStartHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI) Handler {
	return &assistantStartHandler{
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

type assistantStartHandler struct {
	helper *helper
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *assistantStartHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	return NewInfo("assistant-start"), h.helper.sendMessage(
		update.Message.Chat.ID,
		"assistant_menu",
		h.markup,
	)
}
