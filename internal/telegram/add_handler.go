package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

func NewAddHandler(db *sql.DB) Handler {
	return &addHandler{db, "add"}
}

type addHandler struct {
	db      *sql.DB
	command string
}

func (handler *addHandler) Name() string {
	return handler.command
}

func (handler *addHandler) ShouldServe(update tgbotapi.Update) bool {
	return update.Message != nil &&
		update.Message.Chat != nil &&
		update.Message.Command() == handler.command
}

// TODO: move message to file.
func (handler *addHandler) ServeUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (log.Fields, error) {
	fields := log.Fields{"handler": handler.command}
	tx, err := handler.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	_, err = tx.Exec(
		`delete from subscriptions
		where chat_id = $1 and status in ('city', 'price', 'room-number', 'floor')`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to purge subscriptions, %v", err)
	}
	_, err = tx.Exec(
		`insert into subscriptions
		(chat_id, creation_time, status, city, price, room_number, floor)
		values
		($1, now() at time zone 'utc', 'city', 'Київ', 0, 'any', 'any')`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to create a subscription, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	_, err = bot.Send(
		tgbotapi.NewMessage(update.Message.Chat.ID, "Окей, в якому місті шукаємо житло?"),
	)
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return nil, nil
}
