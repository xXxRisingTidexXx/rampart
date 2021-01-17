package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strconv"
)

func NewRoomNumberStatusHandler(
	config config.Handler,
	bot *tgbotapi.BotAPI,
	db *sql.DB,
) StatusHandler {
	return &roomNumberStatusHandler{
		&helper{bot},
		db,
		map[string]misc.RoomNumber{
			config.AnyRoomNumberButton:   misc.AnyRoomNumber,
			config.OneRoomNumberButton:   misc.OneRoomNumber,
			config.TwoRoomNumberButton:   misc.TwoRoomNumber,
			config.ThreeRoomNumberButton: misc.ThreeRoomNumber,
			config.ManyRoomNumberButton:  misc.ManyRoomNumber,
		},
		config.MaxRoomNumber,
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.LowFloorButton),
				tgbotapi.NewKeyboardButton(config.HighFloorButton),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.AnyFloorButton),
				tgbotapi.NewKeyboardButton(config.CancelButton),
			),
		),
	}
}

type roomNumberStatusHandler struct {
	helper        *helper
	db            *sql.DB
	mappings      map[string]misc.RoomNumber
	maxRoomNumber int64
	markup        tgbotapi.ReplyKeyboardMarkup
}

// TODO: long text handling.
func (h *roomNumberStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	var message tgbotapi.MessageConfig
	roomNumber, ok := h.mappings[update.Message.Text]
	if !ok {
		number, err := strconv.ParseInt(update.Message.Text, 10, 64)
		if err != nil || number < int64(misc.ManyRoomNumber) || number > h.maxRoomNumber {
			return h.helper.prepareMessage(update, "invalid_room_number", nil)
		}
		roomNumber = misc.ManyRoomNumber
	}
	_, err := tx.Exec(
		`update transients set room_number = $1, status = $2 where id = $3`,
		roomNumber.String(),
		misc.FloorStatus.String(),
		update.Message.Chat.ID,
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	return h.helper.prepareMessage(update, "valid_room_number", h.markup)
}
