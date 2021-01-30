package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strings"
)

func NewCallbackHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &callbackHandler{
		map[string]Handler{config.DeleteAction: NewDeleteHandler(config, bot, db)},
		config.DataSeparator,
	}
}

type callbackHandler struct {
	handlers      map[string]Handler
	dataSeparator string
}

func (h *callbackHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	index := strings.Index(update.CallbackQuery.Data, h.dataSeparator)
	if index == -1 {
		return log.Fields{"handler": "callback", "action": "absent"}, nil
	}
	if handler, ok := h.handlers[update.CallbackQuery.Data[index:]]; ok {
		return handler.HandleUpdate(update)
	}
	return log.Fields{"handler": "callback", "action": "unknown"}, nil
}
