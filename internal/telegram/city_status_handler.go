package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewCityStatusHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &cityStatusHandler{
		&helper{bot},
		db,
		config.MinFlatCount,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.AnyPriceButton),
				tgbotapi.NewKeyboardButton(config.CancelButton),
			),
		),
	}
}

type cityStatusHandler struct {
	helper       *helper
	db           *sql.DB
	minFlatCount int
	markup       tgbotapi.ReplyKeyboardMarkup
}

// TODO: fuzzy city matching.
// TODO: invalid city metric.
func (h *cityStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	var (
		city    string
		message tgbotapi.MessageConfig
	)
	err := tx.QueryRow(
		`select city from flats where lower(city) = lower($1)`,
		update.Message.Text,
	).Scan(&city)
	if err == sql.ErrNoRows {
		return h.helper.prepareMessage(update, "absent_city", nil)
	}
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to read a city, %v", err)
	}
	_, err = tx.Exec(
		`update transients set status = $1, city = $2 where id = $3`,
		misc.PriceStatus.String(),
		city,
		update.Message.Chat.ID,
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	return h.helper.prepareMessage(update, "present_city", h.markup)
}
