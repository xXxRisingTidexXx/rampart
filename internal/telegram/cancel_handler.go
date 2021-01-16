package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewCancelHandler(bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &cancelHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Підписка \U0001F49C"),
				tgbotapi.NewKeyboardButton("Довідка \U0001F64B"),
			),
		),
	}
}

type cancelHandler struct {
	helper *helper
	db *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *cancelHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "cancel"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, h.helper.sendMessage(update, "menu", h.markup)
}