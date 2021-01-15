package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewCityStatusHandler(bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &cityStatusHandler{
		&helper{bot},
		db,
		5,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Не знаю \U0001F615"),
				tgbotapi.NewKeyboardButton("Головне меню \U00002B05"),
			),
		),
	}
}

type cityStatusHandler struct {
	helper        *helper
	db            *sql.DB
	minFlatCount  int
	absentMarkup  tgbotapi.ReplyKeyboardMarkup
	presentMarkup tgbotapi.ReplyKeyboardMarkup
}

// TODO: fuzzy city matching.
// TODO: check back if city has no flats.
// TODO: move min city flat count to config.
func (h *cityStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	var (
		count   int
		message tgbotapi.MessageConfig
	)
	row := tx.QueryRow(
		`select count(*) from flats where lower(city) = lower($1)`,
		update.Message.Text,
	)
	if err := row.Scan(&count); err != nil {
		return message, fmt.Errorf("telegram: handler failed to read a city, %v", err)
	}
	if count < h.minFlatCount {
		return h.helper.prepareMessage(update, "absent_city", h.absentMarkup)
	}
	_, err := tx.Exec(
		`update transients set status = $1, city = $2 where id = $3`,
		PriceStatus.String(),
		update.Message.Text,
		update.Message.Chat.ID,
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	return h.helper.prepareMessage(update, "present_city", h.presentMarkup)
}
