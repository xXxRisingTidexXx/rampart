package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewModeratorTextHandler(
	config config.ModeratorHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	handlers := make(map[string]Handler)
	handlers[config.StartCommand] = NewModeratorStartHandler(config, bot)
	handlers[config.StartButton] = handlers[config.StartCommand]
	handlers[config.HelpCommand] = NewModeratorHelpHandler(config, bot)
	handlers[config.HelpButton] = handlers[config.HelpCommand]
	handlers[config.MarkupButton] = NewMarkupHandler(config, bot, db)
	handlers[config.LuxuryButton] = NewInteriorHandler(config, bot, db)
	handlers[config.ComfortButton] = handlers[config.LuxuryButton]
	handlers[config.JunkButton] = handlers[config.LuxuryButton]
	handlers[config.ConstructionButton] = handlers[config.LuxuryButton]
	handlers[config.ExcessButton] = handlers[config.LuxuryButton]
	handlers[config.EnoughButton] = NewEnoughHandler(config, bot, db)
	return &moderatorTextHandler{handlers}
}

type moderatorTextHandler struct {
	handlers map[string]Handler
}

func (h *moderatorTextHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	if handler, ok := h.handlers[update.Message.Text]; ok {
		return handler.HandleUpdate(update)
	}
	return NewInfo("moderator-text"), nil
}
