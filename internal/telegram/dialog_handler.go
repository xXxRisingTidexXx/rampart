package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewDialogHandler(bot *tgbotapi.BotAPI, db * sql.DB) Handler {
	return &dialogHandler{&helper{bot}}
}

type dialogHandler struct {
	helper *helper
}

func (h *dialogHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{}, nil
}
