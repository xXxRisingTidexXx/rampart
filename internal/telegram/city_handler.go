package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewCityHandler(db *sql.DB) Handler {
	return &cityHandler{
		db,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type cityHandler struct {
	db     *sql.DB
	markup tgbotapi.ReplyKeyboardMarkup
}

func (handler *cityHandler) Name() string {
	return "city"
}

// TODO: add message randomization.
// TODO: add fuzzy string matching.
// TODO: move to config min city flat count.
// TODO: add branch/option "залишити як є".
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
	var count int
	row := tx.QueryRow(
		`select count(*) from transients where id = $1 and status = 'city'`,
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
	row = tx.QueryRow(`select count(*) from flats where city = $1`, update.Message.Text)
	if err := row.Scan(&count); err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to read a city, %v", err)
	}
	if count < 5 {
		if err := tx.Commit(); err != nil {
			return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
		}
		message := tgbotapi.NewMessage(update.Message.Chat.ID, "Я не знаю помешкань у цьому місті.")
		message.ParseMode = tgbotapi.ModeHTML
		message.ReplyMarkup = handler.markup
		if _, err := bot.Send(message); err != nil {
			return true, fmt.Errorf("telegram: handler failed to send a message, %v", err)
		}
		return true, nil
	}
	_, err = tx.Exec(
		`update transients set status = 'price', city = $1 where id = $2`,
		update.Message.Text,
		update.Message.Chat.ID,
	)
	if err != nil {
		_ = tx.Rollback()
		return true, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	if err := tx.Commit(); err != nil {
		return true, fmt.Errorf("telegram: handler failed to commit a transaction, %v", err)
	}
	message := tgbotapi.NewMessage(update.Message.Chat.ID, "Всьо заєбісь!")
	message.ParseMode = tgbotapi.ModeHTML
	message.ReplyMarkup = handler.markup
	if _, err := bot.Send(message); err != nil {
		return true, fmt.Errorf("telegram: handler failed to send a message, %v", err)
	}
	return true, nil
}
