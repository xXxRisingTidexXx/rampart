package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewFloorHandler(db *sql.DB) XHandler {
	return &floorHandler{
		db,
		map[string]string{"Байдуже \U0001F612": "any", "Ні": "low", "Так": "high"},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type floorHandler struct {
	db       *sql.DB
	mappings map[string]string
	markup   tgbotapi.ReplyKeyboardMarkup
}

func (handler *floorHandler) Name() string {
	return "floor"
}

// TODO: unique subscription generation + final template rendering.
func (handler *floorHandler) HandleUpdate(
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
		`select count(*) from transients where id = $1 and status = 'floor'`,
		update.Message.Chat.ID,
	)
	if err := row.Scan(&count); err != nil {
		_ = tx.Rollback()
		return false, fmt.Errorf("telegram: handler failed to read a transient, %v", err)
	}
	if count == 0 {
		if err := tx.Commit(); err != nil {
			return false, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return false, nil
	}
	floor, ok := handler.mappings[update.Message.Text]
	if !ok {
		if err := tx.Commit(); err != nil {
			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return true, sendMessage(bot, update, "invalid_floor", handler.markup)
	}
	_, err = tx.Exec(
		`insert into subscriptions
		(chat_id, status, city, price, room_number, floor)
		values
		(
			$1,
		 	'actual',
		 	(select city from transients where id = $1),
		 	(select price from transients where id = $1),
		 	(select room_number from transients where id = $1),
		 	$2
		)`,
		update.Message.Chat.ID,
		floor,
	)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to create a subscription, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return true, sendMessage(bot, update, "valid_floor", handler.markup)
}
