package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

func NewHelpHandler() Handler {
	return &helpHandler{
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Зрозуміло \U0001F44C")),
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
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/help.html"))
	if err != nil {
		return true, fmt.Errorf("telegram: handler failed to read a file, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = handler.markup
	if _, err := bot.Send(message); err != nil {
		return true, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return true, nil
}
