package telegram

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewHTMLHandler(command, format string) Handler {
	return &htmlHandler{command, misc.ResolvePath(fmt.Sprintf(format, command))}
}

type htmlHandler struct {
	command string
	path    string
}

func (handler *htmlHandler) ShouldServe(update tgbotapi.Update) bool {
	panic("implement me")
}

func (handler *htmlHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (log.Fields, error) {
	panic("implement me")
}
