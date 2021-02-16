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
		update.Message.Chat.UserName == h.admin {
		log.Info("Hello, bitch!")
	}
	return log.Fields{"handler": "moderator"}, nil
}
