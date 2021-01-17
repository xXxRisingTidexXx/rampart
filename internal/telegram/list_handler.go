package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewListHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &listHandler{&helper{bot}, db}
}

type listHandler struct {
	helper *helper
	db     *sql.DB
}

func (h *listHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "list"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	rows, err := tx.Query(
		`select city, price, room_number, floor
		from subscriptions
		where status = 'actual'
			and chat_id = $1`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read subscriptions, %v", err)
	}
	for rows.Next() {
		var subscription Subscription
		err := rows.Scan(
			&subscription.City,
			&subscription.Price,
			&subscription.RoomNumber,
			&subscription.Floor,
		)
		if err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			return fields, fmt.Errorf("telegram: handler failed to scan a row, %v", err)
		}
		log.Info(subscription)
	}
	if err := rows.Err(); err != nil {
		_ = rows.Close()
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to finish iteration, %v", err)
	}
	if err := rows.Close(); err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to close rows, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return fields, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return fields, nil
}
