package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewHelpHandler() XHandler {
	return &helpHandler{
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type helpHandler struct {
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *helpHandler) Name() string {
	return "help"
}

// TODO: add randomized message texts.
func (handler *helpHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil ||
		update.Message.Chat == nil ||
		update.Message.Text != "Довідка \U0001F64B" {
		return false, nil
	}
	return true, sendMessage(bot, update, "help", handler.markup)
}
