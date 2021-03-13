package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewEnoughHandler(config config.ModeratorHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &enoughHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.MarkupButton),
				tgbotapi.NewKeyboardButton(config.HelpButton),
			),
		),
	}
}

type enoughHandler struct {
	helper *helper
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *enoughHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	info := NewInfo("enough")
	tx, err := h.db.Begin()
	if err != nil {
		return info, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from moderations where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to delete a moderation, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return info, h.helper.sendMessage(update.Message.Chat.ID, "moderator_menu", h.markup)
}
