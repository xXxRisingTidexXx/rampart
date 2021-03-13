package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
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

func (h *assistantHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	if update.Message != nil && update.Message.Chat != nil && update.Message.Text != "" {
		info, err := h.textHandler.HandleUpdate(update)
		info.Extras["chat_id"] = update.Message.Chat.ID
		return info, err
	}
	if update.CallbackQuery != nil &&
		update.CallbackQuery.Message != nil &&
		update.CallbackQuery.Message.Chat != nil &&
		update.CallbackQuery.Data != "" {
		info, err := h.callbackHandler.HandleUpdate(update)
		info.Extras["chat_id"] = update.CallbackQuery.Message.Chat.ID
		return info, err
	}
	return NewInfo("assistant"), nil
}
