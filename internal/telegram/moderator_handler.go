package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &moderatorHandler{NewModeratorTextHandler(config, bot, db)}
}

type moderatorHandler struct {
	handler Handler
}

func (h *moderatorHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	if update.Message != nil && update.Message.Chat != nil && update.Message.Text != "" {
		info, err := h.handler.HandleUpdate(update)
		info.Extras["chat_id"] = update.Message.Chat.ID
		return info, err
	}
	return NewInfo("moderator"), nil
}
