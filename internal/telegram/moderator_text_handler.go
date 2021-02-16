package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorTextHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	handlers := make(map[string]Handler)
	handlers[config.StartCommand] = NewModeratorStartHandler(config, bot)
	handlers[config.ImageMarkupButton] = NewImageMarkupHandler(config, bot, db)
	return &moderatorTextHandler{handlers}
}

type moderatorTextHandler struct {
	handlers map[string]Handler
}

func (h *moderatorTextHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if handler, ok := h.handlers[update.Message.Text]; ok {
		return handler.HandleUpdate(update)
	}
	return log.Fields{"handler": "moderator-text"}, nil
}
