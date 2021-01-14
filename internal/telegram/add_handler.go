package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewAddHandler(db *sql.DB) XHandler {
	return &addHandler{
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type addHandler struct {
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *addHandler) Name() string {
	return "add"
}

// TODO: move default city to config.
// TODO: add message randomization.
func (handler *addHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil ||
		update.Message.Chat == nil ||
		update.Message.Text != "Підписка \U0001F49C" {
		return false, nil
	}
	tx, err := handler.db.Begin()
	if err != nil {
		return true, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(
		`insert into transients
		(id, status, city, price, room_number, floor)
		values
		($1, 'city', 'Київ', 0, 'any', 'any')`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to create a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return true, sendMessage(bot, update, "add", handler.markup)
}
