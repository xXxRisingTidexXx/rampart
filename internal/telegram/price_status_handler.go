package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/config"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strconv"
	"strings"
)

func NewPriceStatusHandler(config config.Handler, bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &priceStatusHandler{
		&helper{bot},
		db,
		config.AnyPriceButton,
		config.MaxPriceLength,
		strings.NewReplacer(",", ".", " ", "", "\n", "", "\t", "", "$", ""),
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.OneRoomNumberButton),
				tgbotapi.NewKeyboardButton(config.TwoRoomNumberButton),
				tgbotapi.NewKeyboardButton(config.ThreeRoomNumberButton),
				tgbotapi.NewKeyboardButton(config.ManyRoomNumberButton),
			),
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(config.AnyRoomNumberButton),
				tgbotapi.NewKeyboardButton(config.CancelButton),
			),
		),
	}
}

type priceStatusHandler struct {
	helper         *helper
	db             *sql.DB
	anyPrice       string
	maxPriceLength int
	replacer       *strings.Replacer
	markup         tgbotapi.ReplyKeyboardMarkup
}

// TODO: invalid input metric.
func (h *priceStatusHandler) HandleStatusUpdate(
	update tgbotapi.Update,
	tx *sql.Tx,
) (tgbotapi.MessageConfig, error) {
	var message tgbotapi.MessageConfig
	if update.Message.Text == h.anyPrice {
		_, err := tx.Exec(
			`update transients set status = $1 where id = $2`,
			misc.RoomNumberStatus.String(),
			update.Message.Chat.ID,
		)
		if err != nil {
			return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
		}
		return h.helper.prepareMessage(update.Message.Chat.ID, "valid_price", h.markup)
	}
	if len(update.Message.Text) > h.maxPriceLength {
		return h.helper.prepareMessage(update.Message.Chat.ID, "invalid_price", nil)
	}
	price, err := strconv.ParseFloat(h.replacer.Replace(update.Message.Text), 64)
	if err != nil || price < 0 {
		return h.helper.prepareMessage(update.Message.Chat.ID, "invalid_price", nil)
	}
	_, err = tx.Exec(
		`update transients set price = $1, status = $2 where id = $3`,
		price,
		misc.RoomNumberStatus.String(),
		update.Message.Chat.ID,
	)
	if err != nil {
		return message, fmt.Errorf("telegram: handler failed to update a transient, %v", err)
	}
	return h.helper.prepareMessage(update.Message.Chat.ID, "valid_price", h.markup)
}
