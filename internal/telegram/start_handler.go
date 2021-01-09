package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

func NewStartHandler() Handler {
	return &startHandler{
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Підписка \U0001F49C"),
				tgbotapi.NewKeyboardButton("Довідка \U0001F64B"),
			),
		),
	}
}

type startHandler struct {
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *startHandler) Name() string {
	return "start"
}

// TODO: add randomized message texts.
func (handler *startHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil ||
		update.Message.Chat == nil ||
		!(update.Message.Command() == handler.Name() ||
			update.Message.Text == "Зрозуміло \U0001F44C") {
		return false, nil
	}
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/start.html"))
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
