package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewTextHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	handlers := make(map[string]Handler)
	handlers["/start"] = NewStartHandler(bot)
	handlers["Зрозуміло \U0001F44D"] = handlers["/start"]
	handlers["/help"] = NewHelpHandler(bot)
	handlers["Довідка \U0001F64B"] = handlers["/help"]
	return &textHandler{handlers}
}

type textHandler struct {
	handlers map[string]Handler
}

func (h *textHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if handler, ok := h.handlers[update.Message.Text]; ok {
		return handler.HandleUpdate(update)
	}
	return log.Fields{"handler": "text"}, nil
}
