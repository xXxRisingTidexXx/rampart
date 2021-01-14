package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewAddHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &addHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type addHandler struct {
	helper *helper
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

// TODO: move default city to config.
func (h *addHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "add"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
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
		return fields, fmt.Errorf("telegram: handler failed to create a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.sendMessage(update, "add", h.markup)
}
