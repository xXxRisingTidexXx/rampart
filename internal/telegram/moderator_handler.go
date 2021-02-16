package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &moderatorHandler{
		NewModeratorTextHandler(config, bot, db),
		NewModeratorCallbackHandler(config, bot, db),
		config.Admin,
	}
}

type moderatorHandler struct {
	textHandler     Handler
	callbackHandler Handler
	admin           string
}

func (h *moderatorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if update.Message != nil &&
		update.Message.Chat != nil &&
		update.Message.Text != "" &&
		update.Message.Chat.UserName == h.admin {
		fields, err := h.textHandler.HandleUpdate(update)
		fields["chat_id"] = update.Message.Chat.ID
		return fields, err
	}
	if update.CallbackQuery != nil &&
		update.CallbackQuery.Message != nil &&
		update.CallbackQuery.Message.Chat != nil &&
		update.CallbackQuery.Data != "" &&
		update.CallbackQuery.Message.Chat.UserName == h.admin {
		fields, err := h.callbackHandler.HandleUpdate(update)
		fields["chat_id"] = update.CallbackQuery.Message.Chat.ID
		return fields, err
	}
	return log.Fields{"handler": "moderator"}, nil
}
