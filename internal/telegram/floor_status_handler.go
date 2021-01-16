package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewFloorStatusHandler(bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &floorStatusHandler{
		&helper{bot},
		db,
		map[string]misc.Floor{
			"Байдуже \U0001F612": misc.AnyFloor,
			"Ні":                 misc.LowFloor,
			"Так":                misc.HighFloor,
		},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Зрозуміло \U0001F44D")),
		),
	}
}

type floorStatusHandler struct {
	helper   *helper
	db       *sql.DB
	mappings map[string]misc.Floor
	markup   tgbotapi.ReplyKeyboardMarkup
}

// TODO: subscription template rendering.
func (h *floorStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	floor, ok := h.mappings[update.Message.Text]
	if !ok {
		return h.helper.prepareMessage(update, "invalid_floor", nil)
	}
	var message tgbotapi.MessageConfig
	_, err := tx.Exec(
		`insert into subscriptions
		(chat_id, city, price, room_number, floor)
		values
		(
			$1,
			(select city from transients where id = $1),
			(select price from transients where id = $1),
			(select room_number from transients where id = $1),
			$2
		)`,
		update.Message.Chat.ID,
		floor.String(),
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to create a subscription, %v", err)
	}
	_, err = tx.Exec(`delete from transients where id = $1`, update.Message.Chat.ID)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to delete a transient, %v", err)
	}
	return h.helper.prepareMessage(update, "valid_floor", h.markup)
}
