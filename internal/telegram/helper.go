package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"html/template"
	"io/ioutil"
	"strings"
)

type helper struct {
	bot *tgbotapi.BotAPI
}

func (h *helper) sendMessage(update tgbotapi.Update, file string, markup interface{}) error {
	message, err := h.prepareMessage(update, file, markup)
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
	file string,
	markup interface{},
) (tgbotapi.MessageConfig, error) {
	var message tgbotapi.MessageConfig
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + file + ".html"))
	if err != nil {
		return message, fmt.Errorf("telegram: helper failed to read a file, %v", err)
	}
	message = tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	return message, nil
}

func (h *helper) sendTemplate(
	update tgbotapi.Update,
	file string,
	data interface{},
	markup interface{},
) error {
	t, err := template.ParseFiles(misc.ResolvePath("templates/" + file + ".html"))
	if err != nil {
		return fmt.Errorf("telegram: helper failed to parse a template, %v", err)
	}
	var builder strings.Builder
	if err := t.Execute(&builder, data); err != nil {
		return fmt.Errorf("telegram: helper failed to execute a template, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, builder.String())
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send a template, %v", err)
	}
	return nil
}

func (h *helper) answerCallback(update tgbotapi.Update) error {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "Відсмокчи мені!")
	callback.ShowAlert = true
	if _, err := h.bot.AnswerCallbackQuery(callback); err != nil {
		return fmt.Errorf("telegram: helper failed to answer a callback query, %v", err)
	}
	return nil
}
