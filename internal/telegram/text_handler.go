package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewTextHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	handlers := make(map[string]Handler)
	handlers["/start"] = NewStartHandler(bot)
	handlers["Зрозуміло \U0001F44D"] = handlers["/start"]
	handlers["/help"] = NewHelpHandler(bot)
	handlers["Довідка \U0001F64B"] = handlers["/help"]
	handlers["Головне меню \U00002B05"] = NewCancelHandler(bot, db)
	handlers["Підписка \U0001F49C"] = NewAddHandler(bot, db)
	return &textHandler{handlers, NewDialogHandler(bot, db)}
}

type textHandler struct {
	commandHandlers map[string]Handler
	dialogHandler   Handler
}

func (h *textHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	if handler, ok := h.commandHandlers[update.Message.Text]; ok {
		return handler.HandleUpdate(update)
	}
	return h.dialogHandler.HandleUpdate(update)
}
