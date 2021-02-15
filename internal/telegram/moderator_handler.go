package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewModeratorHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &moderatorHandler{}
}

type moderatorHandler struct {

}

func (h *moderatorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return nil, nil
}
