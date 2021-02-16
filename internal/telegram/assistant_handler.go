package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewAssistantHandler(
	config config.AssistantHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &assistantHandler{
		NewAssistantTextHandler(config, bot, db),
		NewCallbackHandler(config, bot, db),
	}
}

type assistantHandler struct {
	textHandler     Handler
	callbackHandler Handler
}

func (h *assistantHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if update.Message != nil && update.Message.Chat != nil && update.Message.Text != "" {
		fields, err := h.textHandler.HandleUpdate(update)
		fields["chat_id"] = update.Message.Chat.ID
		return fields, err
	}
	if update.CallbackQuery != nil &&
		update.CallbackQuery.Message != nil &&
		update.CallbackQuery.Message.Chat != nil &&
		update.CallbackQuery.Data != "" {
		fields, err := h.callbackHandler.HandleUpdate(update)
		fields["chat_id"] = update.CallbackQuery.Message.Chat.ID
		return fields, err
	}
	return log.Fields{"handler": "assistant"}, nil
}
