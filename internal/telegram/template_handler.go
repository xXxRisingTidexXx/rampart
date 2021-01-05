package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

func NewHTMLHandler(command, format string) Handler {
	return &htmlHandler{command, misc.ResolvePath(fmt.Sprintf(format, command))}
}

type htmlHandler struct {
	command string
	path    string
}

func (handler *htmlHandler) ShouldServe(update tgbotapi.Update) bool {
	return update.Message != nil &&
		update.Message.Chat != nil &&
		update.Message.Command() == handler.command
}

func (handler *htmlHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (log.Fields, error) {
	fields := log.Fields{"handler": "template", "command": handler.command}
	bytes, err := ioutil.ReadFile(handler.path)
	if err != nil {
		return fields, fmt.Errorf("telegram: html handler failed to read the path, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML

	return nil, nil
}
