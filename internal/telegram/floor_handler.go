package telegram

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func NewFloorHandler(db *sql.DB) Handler {
	return &floorHandler{
		db,
		map[string]string{"Байдуже \uF612": "any", "Ні": "low", "Так": "high"},
		tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("Головне меню \U00002B05")),
		),
	}
}

type floorHandler struct {
	db       *sql.DB
	mappings map[string]string
	markup   tgbotapi.ReplyKeyboardMarkup
}

func (handler *floorHandler) Name() string {
	return "floor"
}

func (handler *floorHandler) HandleUpdate(
	bot *tgbotapi.BotAPI,
	update tgbotapi.Update,
) (bool, error) {
	return true, nil
}
