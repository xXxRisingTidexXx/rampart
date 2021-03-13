package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewAddHandler(config config.AssistantHandler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &addHandler{&helper{bot}, db, tgbotapi.NewKeyboardButton(config.CancelButton)}
}

type addHandler struct {
	helper *helper
	db     *sql.DB
	button tgbotapi.KeyboardButton
}

// TODO: should we add personalized autocomplete based on used cities?
func (h *addHandler) HandleUpdate(update tgbotapi.Update) (Info, error) {
	info := NewInfo("add")
	tx, err := h.db.Begin()
	if err != nil {
		return info, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	_, err = tx.Exec(`insert into transients (id) values ($1)`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to create a transient, %v", err)
	}
	rows, err := tx.Query(`select city from flats group by city order by count(*) desc limit 2`)
	if err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to read cities, %v", err)
	}
	buttons := tgbotapi.NewKeyboardButtonRow()
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return info, fmt.Errorf("telegram: handler failed to scan a row, %v", err)
		}
		buttons = append(buttons, tgbotapi.NewKeyboardButton(city))
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to finish iteration, %v", err)
	}
	if err := rows.Close(); err != nil {
		_ = tx.Rollback()
		return info, fmt.Errorf("telegram: handler failed to close rows, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return info, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	keyboard := make([][]tgbotapi.KeyboardButton, 0)
	if len(buttons) > 0 {
		keyboard = append(keyboard, buttons)
	}
	keyboard = append(keyboard, tgbotapi.NewKeyboardButtonRow(h.button))
	return info, h.helper.sendMessage(
		update.Message.Chat.ID,
		"add",
		tgbotapi.NewReplyKeyboard(keyboard...),
	)
}
