package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

type helper struct {
	bot *tgbotapi.BotAPI
}

// TODO: message randomization.
func (h *helper) sendMessage(update tgbotapi.Update, template string, markup interface{}) error {
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + template + ".html"))
	if err != nil {
		return fmt.Errorf("telegram: helper failed to read a file, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send a message, %v", err)
	}
	return nil
}
