package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/lithammer/shortuuid"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
)

func NewFloorStatusHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &floorStatusHandler{
		&helper{bot},
		db,
		map[string]misc.Floor{
			config.AnyFloorButton:  misc.AnyFloor,
			config.LowFloorButton:  misc.LowFloor,
			config.HighFloorButton: misc.HighFloor,
		},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(config.StartButton)),
		),
	}
}

type floorStatusHandler struct {
	helper   *helper
	db       *sql.DB
	mappings map[string]misc.Floor
	markup   tgbotapi.ReplyKeyboardMarkup
}

// TODO: handle similar subscriptions.
// TODO: invalid input metric.
func (h *floorStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	floor, ok := h.mappings[update.Message.Text]
	if !ok {
		return h.helper.prepareMessage(update.Message.Chat.ID, "invalid_floor", nil)
	}
	var message tgbotapi.MessageConfig
	_, err := tx.Exec(
		`insert into subscriptions
		(uuid, chat_id, city, price, room_number, floor)
		values
		(
			$1,
		 	$2,
			(select city from transients where id = $2),
			(select price from transients where id = $2),
			(select room_number from transients where id = $2),
			$3
		)`,
		shortuuid.New(),
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
	return h.helper.prepareMessage(update.Message.Chat.ID, "valid_floor", h.markup)
}
