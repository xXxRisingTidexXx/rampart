package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewStartHandler(db * sql.DB) Handler {
	return &startHandler{
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Підписка \U0001F49C"),
				tgbotapi.NewKeyboardButton("Довідка \U0001F64B"),
			),
		),
	}
}

type startHandler struct {
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *startHandler) Name() string {
	return "start"
}

// TODO: message randomization.
// TODO: remove database connection - don't drop here.
func (handler *startHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil ||
		update.Message.Chat == nil ||
		!(update.Message.Command() == handler.Name() ||
			update.Message.Text == "Головне меню \U00002B05") {
		return false, nil
	}
	tx, err := handler.db.Begin()
	if err != nil {
		return true, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return true, sendMessage(bot, update, "start", handler.markup)
}
