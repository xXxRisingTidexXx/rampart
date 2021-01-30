package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewCallbackHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &callbackHandler{map[string]Handler{}}
}

type callbackHandler struct {
	handlers map[string]Handler
}

func (h *callbackHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	return log.Fields{"handler": "callback"}, nil
}
