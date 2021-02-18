package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewAssistantTextHandler(
	config config.AssistantHandler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) Handler {
	handlers := make(map[string]Handler)
	handlers[config.StartCommand] = NewAssistantStartHandler(config, bot)
	handlers[config.StartButton] = handlers[config.StartCommand]
	handlers[config.HelpCommand] = NewHelpHandler(config, bot)
	handlers[config.HelpButton] = handlers[config.HelpCommand]
	handlers[config.CancelButton] = NewCancelHandler(config, bot, db)
	handlers[config.AddButton] = NewAddHandler(config, bot, db)
	handlers[config.ListButton] = NewListHandler(config, bot, db)
	return &assistantTextHandler{handlers, NewDialogHandler(config, bot, db)}
}

type assistantTextHandler struct {
	commandHandlers map[string]Handler
	dialogHandler   Handler
}

func (h *assistantTextHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if handler, ok := h.commandHandlers[update.Message.Text]; ok {
		return handler.HandleUpdate(update)
	}
	return h.dialogHandler.HandleUpdate(update)
}
