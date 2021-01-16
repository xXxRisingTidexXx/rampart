package telegram

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/xXxRisingTidexXx/rampart/internal/misc"
	"strconv"
	"strings"
)

func NewPriceStatusHandler(bot *tgbotapi.BotAPI, db *sql.DB) StatusHandler {
	return &priceStatusHandler{
		&helper{bot},
		db,
		"Не знаю \U0001F615",
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
	}
}

type priceStatusHandler struct {
	helper   *helper
	db       *sql.DB
	anyPrice string
	replacer *strings.Replacer
	markup   tgbotapi.ReplyKeyboardMarkup
}

// TODO: invalid input metric.
// TODO: handle too long strings.
// TODO: handle negative price.
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
		return h.helper.prepareMessage(update, "valid_price", h.markup)
	}
	price, err := strconv.ParseFloat(h.replacer.Replace(update.Message.Text), 64)
	if err != nil || price <= 0 {
		return h.helper.prepareMessage(update, "invalid_price", nil)
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
	return h.helper.prepareMessage(update, "valid_price", h.markup)
}
