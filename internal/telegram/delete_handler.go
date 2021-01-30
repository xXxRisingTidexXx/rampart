package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewDeleteHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &deleteHandler{}
}

type deleteHandler struct {

}

func (h *deleteHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "delete"}, nil
}
