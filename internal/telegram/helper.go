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

func (h *helper) sendMessage(update tgbotapi.Update, template string, markup interface{}) error {
	message, err := h.prepareMessage(update, template, markup)
	if err != nil {
		return err
	}
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send a message, %v", err)
	}
	return nil
}

// TODO: message randomization.
func (h *helper) prepareMessage(
	update tgbotapi.Update,
	template string,
	markup interface{},
) (tgbotapi.MessageConfig, error) {
	var message tgbotapi.MessageConfig
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + template + ".html"))
	if err != nil {
		return message, fmt.Errorf("telegram: helper failed to read a file, %v", err)
	}
	message = tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	return message, nil
}
