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
	return &textHandler{handlers}
}

type textHandler struct {
	handlers map[string]Handler
}

func (h *textHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if handler, ok := h.handlers[update.Message.Text]; ok {
		fields, err := handler.HandleUpdate(update)
		fields["chat_id"] = update.Message.Chat.ID
		return fields, err
	}
	return log.Fields{"handler": "text", "chat_id": update.Message.Chat.ID}, nil
}
