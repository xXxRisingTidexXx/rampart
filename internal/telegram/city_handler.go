package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewCityHandler(db *sql.DB) Handler {
	return &cityHandler{db}
}

type cityHandler struct {
	db *sql.DB
}

func (handler *cityHandler) Name() string {
	return "city"
}

// TODO: add fuzzy string matching.
func (handler *cityHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	if update.Message == nil || update.Message.Chat == nil || len(update.Message.Text) < 1 {
		return false, nil
	}
	tx, err := handler.db.Begin()
	if err != nil {
		return false, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	var count int
	row := tx.QueryRow(
		`select count(*) from subscriptions where chat_id = $1 and status = 'city'`,
		update.Message.Chat.ID,
	)
	if err := row.Scan(&count); err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("telegram: handler failed to read subscriptions, %v", err)
	}
	if count == 0 {
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return false, nil
	}

	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}

	return true, nil
}
