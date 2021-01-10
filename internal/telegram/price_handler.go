package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

func NewPriceHandler(db *sql.DB) Handler {
	return &priceHandler{
		db,
		strings.NewReplacer(",", ".", " ", "", "_", "", "\n", "", "\t", ""),
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("1"),
				tgbotapi.NewKeyboardButton("2"),
				tgbotapi.NewKeyboardButton("3"),
				tgbotapi.NewKeyboardButton("4+"),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Байдуже \U0001F612"),
				tgbotapi.NewKeyboardButton("Головне меню \U00002B05"),
			),
		),
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type priceHandler struct {
	db            *sql.DB
	replacer      *strings.Replacer
	validMarkup   tgbotapi.ReplyKeyboardMarkup
	invalidMarkup tgbotapi.ReplyKeyboardMarkup
}

func (handler *priceHandler) Name() string {
	return "price"
}

// TODO: message randomization.
// TODO: invalid input metric.
// TODO: handle too long strings.
// TODO: handle negative price.
// TODO: two template price buttons.
func (handler *priceHandler) HandleUpdate(
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
		`select count(*) from transients where id = $1 and status = 'price'`,
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
	if update.Message.Text == "Не знаю \U0001F615" {
		_, err = tx.Exec(
			`update transients set status = 'room-number' where id = $1`,
			update.Message.Chat.ID,
		)
		if err != nil {
			_ = tx.Rollback()
			return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
		}
		if err := tx.Commit(); err != nil {
			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return true, sendMessage(bot, update, "valid_price", handler.validMarkup)
	}
	price, err := strconv.ParseFloat(handler.replacer.Replace(update.Message.Text), 64)
	if err != nil || price <= 0 {
		if err := tx.Commit(); err != nil {
			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		return true, sendMessage(bot, update, "invalid_price", handler.invalidMarkup)
	}
	_, err = tx.Exec(
		`update transients set price = $1, status = 'room-number' where id = $2`,
		price,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	return true, sendMessage(bot, update, "valid_price", handler.validMarkup)
}
