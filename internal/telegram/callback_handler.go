package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"strings"
)

func NewCallbackHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &callbackHandler{
		map[string]Handler{
			config.DeleteAction: NewDeleteHandler(config, bot, db),
			config.LikeAction:   NewLikeHandler(config, bot, db),
		},
		config.Separator,
	}
}

type callbackHandler struct {
	handlers  map[string]Handler
	separator string
}

func (h *callbackHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "callback"}
	index := strings.Index(update.CallbackQuery.Data, h.separator)
	if index == -1 {
		return fields, nil
	}
	if handler, ok := h.handlers[update.CallbackQuery.Data[:index]]; ok {
		return handler.HandleUpdate(update)
	}
	return fields, nil
}
