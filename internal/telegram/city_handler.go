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

func (handler *cityHandler) HandleUpdate(
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
	var id int
	row := tx.QueryRow(
		`select id from subscriptions where chat_id = $1 and status = 'city' limit 1`,
		update.Message.Chat.ID,
	)
	switch err := row.Scan(&id); err {
	case sql.ErrNoRows:
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return false, nil
	case nil:
		return true, handler.handleCity(tx, id, bot, update)
	default:
		_ = tx.Rollback()
		return false, fmt.Errorf("telegram: handler failed to read a subscription, %v", err)
	}
}

// TODO: add fuzzy string matching.
// TODO: fix min city flat count.
// TODO: add various response for "city not found".
func (handler *cityHandler) handleCity(
	tx *sql.Tx,
	id int,
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) error {
	var count int
	row := tx.QueryRow(`select count(*) from flats where city = $1`, update.Message.Text)
	if err := row.Scan(&count); err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("telegram: handler failed to read a city, %v", err)
	}
	if count < 5 {
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		_, err := bot.Send(
			tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю помешкань у цьому місті."),
		)
		if err != nil {
			return fmt.Errorf("telegram: handler failed to notify about an absent city, %v", err)
		}
		return nil
	}
	_, err := tx.Exec(
		`update subscriptions set status = 'price', city = $1 where id = $2`,
		update.Message.Text,
		id,
	)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("telegram: handler failed to update a subscription, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	_, _ = bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "Всьо заєбісь!"))
	return nil
}
