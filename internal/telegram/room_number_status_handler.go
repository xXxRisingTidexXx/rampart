package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strconv"
)

func NewRoomNumberStatusHandler(bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &roomNumberStatusHandler{
		&helper{bot},
		db,
		map[string]misc.RoomNumber{
			"Байдуже \U0001F612": misc.AnyRoomNumber,
			"1":                  misc.OneRoomNumber,
			"2":                  misc.TwoRoomNumber,
			"3":                  misc.ThreeRoomNumber,
			"4+":                 misc.ManyRoomNumber,
		},
		20,
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
	helper        *helper
	db            *sql.DB
	mappings      map[string]misc.RoomNumber
	maxRoomNumber int64
	markup        tgbotapi.ReplyKeyboardMarkup
}

// TODO: move max room number to config.
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
		roomNumber,
		misc.FloorStatus.String(),
		update.Message.Chat.ID,
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	return h.helper.prepareMessage(update, "valid_room_number", h.markup)
}
