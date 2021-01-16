package telegram

import (
	"database/sql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

//func NewPriceHandler(db *sql.DB) XHandler {
//	return &priceHandler{
//		db,
//		strings.NewReplacer(",", ".", " ", "", "_", "", "\n", "", "\t", ""),
//		tgbotapi.NewReplyKeyboard(
//			tgbotapi.NewKeyboardButtonRow(
//				tgbotapi.NewKeyboardButton("1"),
//				tgbotapi.NewKeyboardButton("2"),
//				tgbotapi.NewKeyboardButton("3"),
//				tgbotapi.NewKeyboardButton("4+"),
//			),
//			tgbotapi.NewKeyboardButtonRow(
//				tgbotapi.NewKeyboardButton("Байдуже \U0001F612"),
//				tgbotapi.NewKeyboardButton("Головне меню \U00002B05"),
//			),
//		),
//		tgbotapi.NewReplyKeyboard(
//			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
//		),
//	}
//}

type priceStatusHandler struct {
	db            *sql.DB
	replacer      *strings.Replacer
	validMarkup   tgbotapi.ReplyKeyboardMarkup
	invalidMarkup tgbotapi.ReplyKeyboardMarkup
}

// TODO: invalid input metric.
// TODO: handle too long strings.
// TODO: handle negative price.
func (h *priceStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
//	if update.Message.Text == "Не знаю \U0001F615" {
//		_, err = tx.Exec(
//			`update transients set status = 'room-number' where id = $1`,
//			update.Message.Chat.ID,
//		)
//		if err != nil {
//			_ = tx.Rollback()
//			return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
//		}
//		if err := tx.Commit(); err != nil {
//			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//		}
//		return true, sendMessage(bot, update, "valid_price", handler.validMarkup)
//	}
//	price, err := strconv.ParseFloat(handler.replacer.Replace(update.Message.Text), 64)
//	if err != nil || price <= 0 {
//		if err := tx.Commit(); err != nil {
//			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//		}
//		return true, sendMessage(bot, update, "invalid_price", handler.invalidMarkup)
//	}
//	_, err = tx.Exec(
//		`update transients set price = $1, status = 'room-number' where id = $2`,
//		price,
//		update.Message.Chat.ID,
//	)
//	if err != nil {
//		_ = tx.Rollback()
//		return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
//	}
//	if err := tx.Commit(); err != nil {
//		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
//	}
//	return true, sendMessage(bot, update, "valid_price", handler.validMarkup)
	return tgbotapi.MessageConfig{}, nil
}
