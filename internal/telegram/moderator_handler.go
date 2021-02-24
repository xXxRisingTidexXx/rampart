package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewModeratorHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	return &moderatorHandler{NewModeratorTextHandler(config, bot, db), config.Admins}
}

type moderatorHandler struct {
	handler Handler
	admins  misc.Set
}

func (h *moderatorHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if update.Message != nil &&
		update.Message.Chat != nil &&
		update.Message.Text != "" &&
		h.admins.Has(update.Message.Chat.UserName) {
		fields, err := h.handler.HandleUpdate(update)
		fields["chat_id"] = update.Message.Chat.ID
		return fields, err
	}
	return log.Fields{"handler": "moderator"}, nil
}
