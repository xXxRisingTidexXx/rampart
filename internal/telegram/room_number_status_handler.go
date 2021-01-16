package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewRoomNumberStatusHandler(bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &roomNumberStatusHandler{
		&helper{bot},
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Ні"),
				tgbotapi.NewKeyboardButton("Так"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Байдуже \U0001F612"),
				tgbotapi.NewKeyboardButton("Головне меню \U00002B05"),
			),
		),
	}
}

type roomNumberStatusHandler struct {
	helper *helper
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

// TODO: number and other text handling.
// TODO: long text handling.
func (h *roomNumberStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
//	if update.Message == nil || update.Message.Chat == nil || len(update.Message.Text) < 1 {
//		return false, nil
//	}
//	tx, err := handler.db.Begin()
//	if err != nil {
//		return false, fmt.Errorf("telegram: handler failed to begin a transaction, %v", err)
//	}
//	var count int
//	row := tx.QueryRow(
//		`select count(*) from transients where id = $1 and status = 'room-number'`,
//		update.Message.Chat.ID,
//	)
//	if err := row.Scan(&count); err != nil {
//		_ = tx.Rollback()
//		return false, fmt.Errorf("telegram: handler failed to read a transient, %v", err)
//	}
//	if count == 0 {
//		if err := tx.Commit(); err != nil {
//			return false, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//		}
//		return false, nil
//	}
//	roomNumber, ok := handler.mappings[update.Message.Text]
//	if !ok {
//		if err := tx.Commit(); err != nil {
//			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//		}
//		return true, sendMessage(bot, update, "invalid_room_number", handler.invalidMarkup)
//	}
//	_, err = tx.Exec(
//		`update transients set status = 'floor', room_number = $1 where id = $2`,
//		roomNumber,
//		update.Message.Chat.ID,
//	)
//	if err != nil {
//		_ = tx.Rollback()
//		return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
//	}
//	if err := tx.Commit(); err != nil {
//		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//	}
//	return true, sendMessage(bot, update, "valid_room_number", handler.validMarkup)
	return tgbotapi.MessageConfig{}, nil
}
