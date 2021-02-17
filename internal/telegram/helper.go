package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type helper struct {
	bot *tgbotapi.BotAPI
}

func (h *helper) sendMessage(chatID int64, file string, markup interface{}) error {
	message, err := h.prepareMessage(chatID, file, markup)
	if err != nil {
		return err
	}
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send a message, %v", err)
	}
	return nil
}

// TODO: file randomization.
func (h *helper) prepareMessage(
	chatID int64,
	file string,
	markup interface{},
) (tgbotapi.MessageConfig, error) {
	var message tgbotapi.MessageConfig
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + file + ".html"))
	if err != nil {
		return message, fmt.Errorf("telegram: helper failed to read a file, %v", err)
	}
	message = tgbotapi.NewMessage(chatID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	return message, nil
}

func (h *helper) sendTemplate(chatID int64, file string, data, markup interface{}) error {
	t, err := template.ParseFiles(misc.ResolvePath("templates/" + file + ".html"))
	if err != nil {
		return fmt.Errorf("telegram: helper failed to parse a template, %v", err)
	}
	var builder strings.Builder
	if err := t.Execute(&builder, data); err != nil {
		return fmt.Errorf("telegram: helper failed to execute a template, %v", err)
	}
	message := tgbotapi.NewMessage(chatID, builder.String())
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = markup
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send a template, %v", err)
	}
	return nil
}

// TODO: file randomization.
func (h *helper) answerCallback(callbackID, file string) error {
	bytes, err := ioutil.ReadFile(misc.ResolvePath("templates/" + file + ".html"))
	if err != nil {
		return fmt.Errorf("telegram: helper failed to read a file, %v", err)
	}
	_, err = h.bot.AnswerCallbackQuery(tgbotapi.NewCallback(callbackID, string(bytes)))
	if err != nil {
		return fmt.Errorf("telegram: helper failed to answer a callback query, %v", err)
	}
	return nil
}

// TODO: should we just send url or separate http client?
func (h *helper) sendImage(chatID int64, url string, markup interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("telegram: helper failed to make a request, %v", err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		_ = response.Body.Close()
		return fmt.Errorf("telegram: helper failed to read a response body, %v", err)
	}
	if err := response.Body.Close(); err != nil {
		return fmt.Errorf("telegram: helper faield to close a response body, %v", err)
	}
	message := tgbotapi.NewPhotoUpload(chatID, tgbotapi.FileBytes{Bytes: bytes})
	message.ReplyMarkup = markup
	if _, err := h.bot.Send(message); err != nil {
		return fmt.Errorf("telegram: helper failed to send an image, %v", err)
	}
	return nil
}
