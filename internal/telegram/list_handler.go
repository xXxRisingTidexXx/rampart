package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
)

func NewListHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) Handler {
	return &listHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
		),
	}
}

type listHandler struct {
	helper *helper
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (h *listHandler) HandleUpdate(update tgbotapi.Update) (log.Fields, error) {
	fields := log.Fields{"handler": "list"}
	tx, err := h.db.Begin()
	if err != nil {
		return fields, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
	}
	rows, err := tx.Query(
		`select id, city, price, room_number, floor
		from subscriptions
		where status = 'actual'
			and chat_id = $1`,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return fields, fmt.Errorf("telegram: handler failed to read subscriptions, %v", err)
	}
	subscriptions := make([]subscription, 0)
	for rows.Next() {
		var s subscription
		err := rows.Scan(&s.ID, &s.City, &s.Price, &s.RoomNumber, &s.Floor)
		if err != nil {
			_ = rows.Close()
			_ = tx.Rollback()
			fields["subscription_id"] = s.ID
			return fields, fmt.Errorf("telegram: handler failed to scan a row, %v", err)
		}
		subscriptions = append(subscriptions, s)
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
	if len(subscriptions) == 0 {
		return fields, h.helper.sendMessage(update, "empty_list", h.markup)
	}
	for _, s := range subscriptions {
		if err := h.helper.sendTemplate(update, "full_list", s, h.markup); err != nil {
			fields["subscription_id"] = s.ID
			return fields, err
		}
	}
	return fields, nil
}
