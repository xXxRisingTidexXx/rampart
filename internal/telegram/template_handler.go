package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"io/ioutil"
)

// TODO: configure start template.
func NewTemplateHandler(command, format string) Handler {
	return &templateHandler{command, misc.ResolvePath(fmt.Sprintf(format, command))}
}

type templateHandler struct {
	command string
	path    string
}

func (handler *templateHandler) Name() string {
	return handler.command
}

func (handler *templateHandler) ShouldServe(update tgbotapi.Update) bool {
	return update.Message != nil &&
		update.Message.Chat != nil &&
		update.Message.Command() == handler.command
}

func (handler *templateHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (log.Fields, error) {
	fields := log.Fields{"handler": handler.command}
	bytes, err := ioutil.ReadFile(handler.path)
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to read the path, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, string(bytes))
	message.ParseMode = tgbotapi.ModeHTML
	if _, err := bot.Send(message); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return nil, nil
}
