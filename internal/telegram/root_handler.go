package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewRootHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &rootHandler{NewTextHandler(bot, db)}
}

type rootHandler struct {
	handler Handler
}

func (h *rootHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if update.Message != nil && update.Message.Chat != nil {
		return h.handler.HandleUpdate(update)
	}
	return log.Fields{"handler": "root"}, nil
}
